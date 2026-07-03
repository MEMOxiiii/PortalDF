package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// DisconnectPlayer is sent by the proxy to a server to request that it disconnects any existing session for
// the specified player. This is sent before transferring a player to ensure stale sessions are cleaned up.
type DisconnectPlayer struct {
	PlayerName string
}

func (*DisconnectPlayer) ID() uint16 { return IDDisconnectPlayer }

func (pk *DisconnectPlayer) Marshal(w *protocol.Writer) {
	w.String(&pk.PlayerName)
}

func (pk *DisconnectPlayer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.PlayerName)
}
