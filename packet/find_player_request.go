package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// FindPlayerRequest is sent by the client to find the server a player is currently on.
type FindPlayerRequest struct {
	PlayerUUID uuid.UUID
	PlayerName string
}

func (*FindPlayerRequest) ID() uint16 { return IDFindPlayerRequest }

func (pk *FindPlayerRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.String(&pk.PlayerName)
}

func (pk *FindPlayerRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.String(&pk.PlayerName)
}
