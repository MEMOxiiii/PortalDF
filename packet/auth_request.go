package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// AuthRequest is sent by the client to authenticate with the proxy.
type AuthRequest struct {
	Protocol uint32
	Secret   string
	Name     string
}

func (*AuthRequest) ID() uint16 { return IDAuthRequest }

func (pk *AuthRequest) Marshal(w *protocol.Writer) {
	w.Uint32(&pk.Protocol)
	w.String(&pk.Secret)
	w.String(&pk.Name)
}

func (pk *AuthRequest) Unmarshal(r *protocol.Reader) {
	r.Uint32(&pk.Protocol)
	r.String(&pk.Secret)
	r.String(&pk.Name)
}
