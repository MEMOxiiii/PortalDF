package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// RegisterServer is sent by the client to register itself with an address on the proxy.
type RegisterServer struct {
	Address string
}

func (*RegisterServer) ID() uint16 { return IDRegisterServer }

func (pk *RegisterServer) Marshal(w *protocol.Writer) {
	w.String(&pk.Address)
}

func (pk *RegisterServer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Address)
}
