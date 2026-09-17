package command

import (
	"github.com/MEMOxiiii/PortalDF/packet"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// portalClient is the subset of *portaldf.Portal used by these commands. It is declared here instead of
// depending on the root package directly, so that the root package can register these commands from a
// single Enable call without creating an import cycle (the root package already needs to import this one).
type portalClient interface {
	ServerName() string
	Connected() bool
	FindPlayer(playerUUID uuid.UUID, playerName string, callback packet.FindPlayerCallback) error
	RequestServerList(callback packet.ServerListCallback) error
	TransferPlayer(playerUUID uuid.UUID, server string, callback packet.TransferCallback) error
	SetDisconnectPlayerHandler(handler packet.DisconnectPlayerHandler)
}

var (
	portalRef portalClient
	serverRef *server.Server
)

// Register registers all Portal commands with the dragonfly command system.
// It requires a connected Portal instance and the dragonfly Server for player lookups.
func Register(p portalClient, srv *server.Server) {
	portalRef = p
	serverRef = srv

	cmd.Register(cmd.New("transfer", "Transfer a player to another server.", nil, TransferSelf{}, TransferOther{}))
	cmd.Register(cmd.New("server", "Check which server a player is on.", nil, ServerSelf{}, ServerOther{}))
	cmd.Register(cmd.New("servers", "List all servers connected to the proxy.", nil, Servers{}))

	// Disconnect any stale local session for a player the proxy is about to transfer here, so they don't
	// end up with two connections open at once.
	p.SetDisconnectPlayerHandler(func(playerName string) {
		if serverRef == nil {
			return
		}
		handle, ok := serverRef.PlayerByName(playerName)
		if !ok {
			return
		}
		handle.Do(func(tx *world.Tx, e world.Entity) {
			if pl, ok := e.(*player.Player); ok {
				pl.Disconnect("Connecting from another location")
			}
		})
	})
}
