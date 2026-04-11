package portaldf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hexomc/portaldf/packet"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TransferCallback is called when a transfer response is received from the proxy.
// status is one of the TransferResponse constants. err is only set when status == TransferResponseError.
type TransferCallback func(playerUUID uuid.UUID, status byte, err string)

// PlayerInfoCallback is called when a player info response is received from the proxy.
type PlayerInfoCallback func(playerUUID uuid.UUID, status byte, xuid string, address string)

// ServerListCallback is called when a server list response is received from the proxy.
type ServerListCallback func(servers []packet.ServerEntry)

// FindPlayerCallback is called when a find player response is received from the proxy.
type FindPlayerCallback func(playerUUID uuid.UUID, playerName string, online bool, server string)

// LatencyHandler is called when the proxy sends a player latency update.
type LatencyHandler func(playerUUID uuid.UUID, latency int64)

// Portal represents a connection to a Portal proxy's socket server.
// It handles authentication, server registration, and provides methods
// for transferring players and querying information from the proxy.
type Portal struct {
	config Config
	log    *slog.Logger

	conn   net.Conn
	connMu sync.Mutex

	pool packet.Pool

	sendMu sync.Mutex
	hdr    *packet.Header
	buf    *bytes.Buffer

	closeCh chan struct{}
	closed  bool

	transferCallbacks   map[uuid.UUID]TransferCallback
	transferMu          sync.Mutex
	playerInfoCallbacks map[uuid.UUID]PlayerInfoCallback
	playerInfoMu        sync.Mutex
	serverListCallbacks []ServerListCallback
	serverListMu        sync.Mutex
	findPlayerCallbacks map[string]FindPlayerCallback
	findPlayerMu        sync.Mutex

	latencyHandler LatencyHandler
	latencyMu      sync.RWMutex

	connected   bool
	connectedMu sync.RWMutex
}

// New creates a new Portal client with the given configuration. Call Connect() to start the connection.
func New(config Config, log *slog.Logger) *Portal {
	if log == nil {
		log = slog.Default()
	}
	return &Portal{
		config:  config,
		log:     log,
		pool:    packet.NewPool(),
		hdr:     &packet.Header{},
		buf:     bytes.NewBuffer(make([]byte, 0, 4096)),
		closeCh: make(chan struct{}),

		transferCallbacks:   make(map[uuid.UUID]TransferCallback),
		playerInfoCallbacks: make(map[uuid.UUID]PlayerInfoCallback),
		findPlayerCallbacks: make(map[string]FindPlayerCallback),
	}
}

// Connect starts the connection to the Portal proxy and blocks until Close() is called.
// It will automatically reconnect if the connection is lost.
func (p *Portal) Connect() {
	for {
		select {
		case <-p.closeCh:
			return
		default:
		}

		err := p.connectOnce()
		if err != nil {
			p.log.Error("Portal socket connection error", "error", err)
		}

		p.setConnected(false)

		select {
		case <-p.closeCh:
			return
		case <-time.After(5 * time.Second):
			p.log.Info("Reconnecting to Portal proxy...")
		}
	}
}

// Close closes the connection to the Portal proxy.
func (p *Portal) Close() error {
	if p.closed {
		return nil
	}
	p.closed = true
	close(p.closeCh)

	p.connMu.Lock()
	defer p.connMu.Unlock()
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}

// Connected returns whether the client is currently connected and authenticated with the proxy.
func (p *Portal) Connected() bool {
	p.connectedMu.RLock()
	defer p.connectedMu.RUnlock()
	return p.connected
}

// ServerName returns the configured server name that this instance registered as on the proxy.
func (p *Portal) ServerName() string {
	return p.config.ServerName
}

func (p *Portal) setConnected(v bool) {
	p.connectedMu.Lock()
	p.connected = v
	p.connectedMu.Unlock()
}

// TransferPlayer sends a transfer request to the proxy for the given player UUID to the target server.
func (p *Portal) TransferPlayer(playerUUID uuid.UUID, server string, callback TransferCallback) error {
	if callback != nil {
		p.transferMu.Lock()
		p.transferCallbacks[playerUUID] = callback
		p.transferMu.Unlock()
	}
	return p.writePacket(&packet.TransferRequest{
		PlayerUUID: playerUUID,
		Server:     server,
	})
}

// RequestPlayerInfo requests information (XUID, IP) about a player from the proxy.
func (p *Portal) RequestPlayerInfo(playerUUID uuid.UUID, callback PlayerInfoCallback) error {
	if callback != nil {
		p.playerInfoMu.Lock()
		p.playerInfoCallbacks[playerUUID] = callback
		p.playerInfoMu.Unlock()
	}
	return p.writePacket(&packet.PlayerInfoRequest{
		PlayerUUID: playerUUID,
	})
}

// RequestServerList requests the list of all servers connected to the proxy.
func (p *Portal) RequestServerList(callback ServerListCallback) error {
	if callback != nil {
		p.serverListMu.Lock()
		sendPacket := len(p.serverListCallbacks) == 0
		p.serverListCallbacks = append(p.serverListCallbacks, callback)
		p.serverListMu.Unlock()

		if !sendPacket {
			return nil
		}
	}
	return p.writePacket(&packet.ServerListRequest{})
}

// FindPlayer searches for a player on the proxy by UUID and/or name.
// If uuid is uuid.Nil, the search is by name only.
func (p *Portal) FindPlayer(playerUUID uuid.UUID, playerName string, callback FindPlayerCallback) error {
	if callback != nil {
		p.findPlayerMu.Lock()
		if playerUUID != uuid.Nil {
			p.findPlayerCallbacks[playerUUID.String()] = callback
		} else {
			p.findPlayerCallbacks[playerName] = callback
		}
		p.findPlayerMu.Unlock()
	}
	return p.writePacket(&packet.FindPlayerRequest{
		PlayerUUID: playerUUID,
		PlayerName: playerName,
	})
}

// SetLatencyHandler sets the handler that will be called when the proxy sends player latency updates.
func (p *Portal) SetLatencyHandler(handler LatencyHandler) {
	p.latencyMu.Lock()
	p.latencyHandler = handler
	p.latencyMu.Unlock()
}

// connectOnce establishes a single connection to the proxy, authenticates, registers, and enters the read loop.
func (p *Portal) connectOnce() error {
	address := net.JoinHostPort(p.config.ProxyAddress, fmt.Sprintf("%d", p.config.SocketPort))
	p.log.Info("Connecting to Portal proxy", "address", address)

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	p.connMu.Lock()
	p.conn = conn
	p.connMu.Unlock()

	defer func() {
		_ = conn.Close()
		p.connMu.Lock()
		p.conn = nil
		p.connMu.Unlock()
	}()

	// Send auth request.
	if err := p.writePacket(&packet.AuthRequest{
		Protocol: packet.ProtocolVersion,
		Secret:   p.config.Secret,
		Name:     p.config.ServerName,
	}); err != nil {
		return fmt.Errorf("failed to send auth request: %w", err)
	}

	// Read loop.
	for {
		select {
		case <-p.closeCh:
			return nil
		default:
		}

		pk, err := p.readPacket()
		if err != nil {
			return fmt.Errorf("read error: %w", err)
		}

		p.handlePacket(pk)
	}
}

// handlePacket processes a received packet from the proxy.
func (p *Portal) handlePacket(pk packet.Packet) {
	switch pk := pk.(type) {
	case *packet.AuthResponse:
		p.handleAuthResponse(pk)
	case *packet.TransferResponse:
		p.handleTransferResponse(pk)
	case *packet.PlayerInfoResponse:
		p.handlePlayerInfoResponse(pk)
	case *packet.ServerListResponse:
		p.handleServerListResponse(pk)
	case *packet.FindPlayerResponse:
		p.handleFindPlayerResponse(pk)
	case *packet.UpdatePlayerLatency:
		p.handleUpdatePlayerLatency(pk)
	}
}

func (p *Portal) handleAuthResponse(pk *packet.AuthResponse) {
	if pk.Status != packet.AuthResponseSuccess {
		var reason string
		switch pk.Status {
		case packet.AuthResponseUnsupportedProtocol:
			reason = fmt.Sprintf("unsupported protocol version, proxy expects %d, we have %d", pk.Protocol, packet.ProtocolVersion)
		case packet.AuthResponseIncorrectSecret:
			reason = "incorrect secret provided"
		case packet.AuthResponseAlreadyConnected:
			reason = "client with this name already connected"
		case packet.AuthResponseUnauthenticated:
			reason = "attempted to send packets whilst not authenticated"
		default:
			reason = fmt.Sprintf("unknown auth error status: %d", pk.Status)
		}
		p.log.Error("Portal authentication failed", "reason", reason)
		return
	}

	p.log.Info("Authenticated with Portal proxy")
	p.setConnected(true)

	// Register our server address with the proxy.
	if err := p.writePacket(&packet.RegisterServer{
		Address: p.config.ServerAddress,
	}); err != nil {
		p.log.Error("Failed to send register server packet", "error", err)
	} else {
		p.log.Info("Registered server with Portal proxy", "name", p.config.ServerName, "address", p.config.ServerAddress)
	}
}

func (p *Portal) handleTransferResponse(pk *packet.TransferResponse) {
	p.transferMu.Lock()
	cb, ok := p.transferCallbacks[pk.PlayerUUID]
	if ok {
		delete(p.transferCallbacks, pk.PlayerUUID)
	}
	p.transferMu.Unlock()

	if ok && cb != nil {
		cb(pk.PlayerUUID, pk.Status, pk.Error)
	}
}

func (p *Portal) handlePlayerInfoResponse(pk *packet.PlayerInfoResponse) {
	p.playerInfoMu.Lock()
	cb, ok := p.playerInfoCallbacks[pk.PlayerUUID]
	if ok {
		delete(p.playerInfoCallbacks, pk.PlayerUUID)
	}
	p.playerInfoMu.Unlock()

	if ok && cb != nil {
		cb(pk.PlayerUUID, pk.Status, pk.XUID, pk.Address)
	}
}

func (p *Portal) handleServerListResponse(pk *packet.ServerListResponse) {
	p.serverListMu.Lock()
	callbacks := p.serverListCallbacks
	p.serverListCallbacks = nil
	p.serverListMu.Unlock()

	for _, cb := range callbacks {
		cb(pk.Servers)
	}
}

func (p *Portal) handleFindPlayerResponse(pk *packet.FindPlayerResponse) {
	server := ""
	if pk.Online {
		server = pk.Server
	}

	p.findPlayerMu.Lock()
	// Try by UUID first, then by name.
	key := pk.PlayerUUID.String()
	cb, ok := p.findPlayerCallbacks[key]
	if !ok {
		key = pk.PlayerName
		cb, ok = p.findPlayerCallbacks[key]
	}
	if ok {
		delete(p.findPlayerCallbacks, key)
	}
	p.findPlayerMu.Unlock()

	if ok && cb != nil {
		cb(pk.PlayerUUID, pk.PlayerName, pk.Online, server)
	}
}

func (p *Portal) handleUpdatePlayerLatency(pk *packet.UpdatePlayerLatency) {
	p.latencyMu.RLock()
	handler := p.latencyHandler
	p.latencyMu.RUnlock()

	if handler != nil {
		handler(pk.PlayerUUID, pk.Latency)
	}
}

// writePacket writes a packet to the proxy connection.
func (p *Portal) writePacket(pk packet.Packet) error {
	p.connMu.Lock()
	conn := p.conn
	p.connMu.Unlock()

	if conn == nil {
		return fmt.Errorf("not connected to proxy")
	}

	p.sendMu.Lock()
	defer p.sendMu.Unlock()

	p.hdr.PacketID = pk.ID()
	_ = p.hdr.Write(p.buf)
	pk.Marshal(protocol.NewWriter(p.buf, 0))

	data := make([]byte, p.buf.Len())
	copy(data, p.buf.Bytes())
	p.buf.Reset()

	out := bytes.NewBuffer(make([]byte, 0, 4+len(data)))
	if err := binary.Write(out, binary.LittleEndian, int32(len(data))); err != nil {
		return err
	}
	if _, err := out.Write(data); err != nil {
		return err
	}

	_, err := conn.Write(out.Bytes())
	return err
}

// readPacket reads a single packet from the proxy connection.
func (p *Portal) readPacket() (packet.Packet, error) {
	p.connMu.Lock()
	conn := p.conn
	p.connMu.Unlock()

	if conn == nil {
		return nil, fmt.Errorf("not connected to proxy")
	}

	// Read 4-byte length prefix.
	var length uint32
	if err := binary.Read(conn, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length > 1024*1024 {
		return nil, fmt.Errorf("packet too large: %d bytes", length)
	}

	// Read the full packet data.
	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer(data)
	header := &packet.Header{}
	if err := header.Read(buf); err != nil {
		return nil, err
	}

	pk, ok := p.pool[header.PacketID]
	if !ok {
		return nil, fmt.Errorf("unknown packet ID: 0x%02x", header.PacketID)
	}

	pk.Unmarshal(protocol.NewReader(buf, 0, false))
	return pk, nil
}
