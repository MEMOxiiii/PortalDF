package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// FindPlayerResponse is sent by the proxy in response to FindPlayerRequest.
type FindPlayerResponse struct {
	PlayerUUID uuid.UUID
	PlayerName string
	Online     bool
	Server     string
}

func (*FindPlayerResponse) ID() uint16 { return IDFindPlayerResponse }

func (pk *FindPlayerResponse) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.String(&pk.PlayerName)
	w.Bool(&pk.Online)
	if pk.Online {
		w.String(&pk.Server)
	}
}

func (pk *FindPlayerResponse) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.String(&pk.PlayerName)
	r.Bool(&pk.Online)
	if pk.Online {
		r.String(&pk.Server)
	}
}
