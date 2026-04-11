package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// UpdatePlayerLatency is sent by the proxy to update a player's latency.
type UpdatePlayerLatency struct {
	PlayerUUID uuid.UUID
	Latency    int64
}

func (*UpdatePlayerLatency) ID() uint16 { return IDUpdatePlayerLatency }

func (pk *UpdatePlayerLatency) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Int64(&pk.Latency)
}

func (pk *UpdatePlayerLatency) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Int64(&pk.Latency)
}
