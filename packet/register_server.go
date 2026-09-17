package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// RegisterServer is sent by the client to register itself with an address on the proxy.
type RegisterServer struct {
	// Address is the address of this server. For Transport "raknet" this is a "host:port" pair. For
	// Transport "nethernet" this is the URL of this server's HTTP(S) signaling endpoint.
	Address string
	// Transport is the network transport the proxy should use to dial this server: "raknet" or
	// "nethernet". An empty value is treated by the proxy as "raknet".
	Transport  string
	LegacyAuth bool
	// Group is the name of the group this server belongs to, used by group-aware load balancers on the
	// proxy to route players to the correct set of servers. It may be left empty if the server does not
	// belong to a group.
	Group string
	// Weight controls how large a share of new players this server should receive relative to others in
	// the same group. A weight of 0 is treated by the proxy as 1.
	Weight uint32
}

func (*RegisterServer) ID() uint16 { return IDRegisterServer }

func (pk *RegisterServer) Marshal(w *protocol.Writer) {
	w.String(&pk.Address)
	w.String(&pk.Transport)
	w.Bool(&pk.LegacyAuth)
	w.String(&pk.Group)
	w.Varuint32(&pk.Weight)
}

func (pk *RegisterServer) Unmarshal(r *protocol.Reader) {
	r.String(&pk.Address)
	r.String(&pk.Transport)
	r.Bool(&pk.LegacyAuth)
	r.String(&pk.Group)
	r.Varuint32(&pk.Weight)
}
