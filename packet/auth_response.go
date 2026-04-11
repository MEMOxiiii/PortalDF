package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

const (
	AuthResponseSuccess byte = iota
	AuthResponseUnsupportedProtocol
	AuthResponseIncorrectSecret
	AuthResponseAlreadyConnected
	AuthResponseUnauthenticated
)

// AuthResponse is sent by the proxy in response to AuthRequest.
type AuthResponse struct {
	Protocol uint32
	Status   byte
}

func (*AuthResponse) ID() uint16 { return IDAuthResponse }

func (pk *AuthResponse) Marshal(w *protocol.Writer) {
	w.Uint32(&pk.Protocol)
	w.Uint8(&pk.Status)
}

func (pk *AuthResponse) Unmarshal(r *protocol.Reader) {
	r.Uint32(&pk.Protocol)
	r.Uint8(&pk.Status)
}
