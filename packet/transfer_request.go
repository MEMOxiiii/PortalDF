package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// TransferRequest is sent by the client to request a player transfer.
type TransferRequest struct {
	PlayerUUID uuid.UUID
	Server     string
}

func (*TransferRequest) ID() uint16 { return IDTransferRequest }

func (pk *TransferRequest) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.String(&pk.Server)
}

func (pk *TransferRequest) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.String(&pk.Server)
}
