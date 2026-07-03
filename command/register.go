package command

import (
	portaldf "github.com/MEMOxiiii/PortalDF"
	"github.com/df-mc/dragonfly/server"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
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

	// Disconnect any stale local session for a player the proxy is about to transfer here, so they don't
	// end up with two connections open at once.
	p.SetDisconnectPlayerHandler(func(playerName string) {
		if serverRef == nil {
			return
		}
		handle, ok := serverRef.PlayerByName(playerName)
		if !ok {
			return
		}
		handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
			if pl, ok := e.(*player.Player); ok {
				pl.Disconnect("Connecting from another location")
			}
		})
	})
}
