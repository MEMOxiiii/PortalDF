package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// ServerEntry represents a server connected to the proxy.
type ServerEntry struct {
	Name        string
	PlayerCount int64
}

// ServerListResponse is sent by the proxy in response to ServerListRequest.
type ServerListResponse struct {
	Servers []ServerEntry
}

func (*ServerListResponse) ID() uint16 { return IDServerListResponse }

func (pk *ServerListResponse) Marshal(w *protocol.Writer) {
	l := uint32(len(pk.Servers))
	w.Uint32(&l)
	for _, s := range pk.Servers {
		w.String(&s.Name)
		w.Int64(&s.PlayerCount)
	}
}

func (pk *ServerListResponse) Unmarshal(r *protocol.Reader) {
	var l uint32
	r.Uint32(&l)
	pk.Servers = make([]ServerEntry, l)
	for i := uint32(0); i < l; i++ {
		r.String(&pk.Servers[i].Name)
		r.Int64(&pk.Servers[i].PlayerCount)
	}
}
