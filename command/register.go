package command

import (
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/hexomc/portaldf"
)

var (
	portalRef *portaldf.Portal
	serverRef *server.Server
)

// Register registers all Portal commands with the dragonfly command system.
// It requires a connected Portal instance and the dragonfly Server for player lookups.
func Register(p *portaldf.Portal, srv *server.Server) {
	portalRef = p
	serverRef = srv

	cmd.Register(cmd.New("transfer", "Transfer a player to another server.", nil, TransferSelf{}, TransferOther{}))
	cmd.Register(cmd.New("server", "Check which server a player is on.", nil, ServerSelf{}, ServerOther{}))
	cmd.Register(cmd.New("servers", "List all servers connected to the proxy.", nil, Servers{}))
}
