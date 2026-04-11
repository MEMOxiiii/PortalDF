package packet

import "github.com/sandertv/gophertunnel/minecraft/protocol"

// ServerListRequest is sent by the client to request a list of all servers.
type ServerListRequest struct{}

func (*ServerListRequest) ID() uint16 { return IDServerListRequest }

func (pk *ServerListRequest) Marshal(*protocol.Writer) {}

func (pk *ServerListRequest) Unmarshal(*protocol.Reader) {}
