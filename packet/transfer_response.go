package packet

import (
	"github.com/google/uuid"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	TransferResponseSuccess byte = iota
	TransferResponseServerNotFound
	TransferResponseAlreadyOnServer
	TransferResponsePlayerNotFound
	TransferResponseError
)

// TransferResponse is sent by the proxy in response to a transfer request.
type TransferResponse struct {
	PlayerUUID uuid.UUID
	Status     byte
	Error      string
}

func (*TransferResponse) ID() uint16 { return IDTransferResponse }

func (pk *TransferResponse) Marshal(w *protocol.Writer) {
	w.UUID(&pk.PlayerUUID)
	w.Uint8(&pk.Status)
	if pk.Status == TransferResponseError {
		w.String(&pk.Error)
	}
}

func (pk *TransferResponse) Unmarshal(r *protocol.Reader) {
	r.UUID(&pk.PlayerUUID)
	r.Uint8(&pk.Status)
	if pk.Status == TransferResponseError {
		r.String(&pk.Error)
	}
}
