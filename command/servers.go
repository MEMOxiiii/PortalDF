package command

import (
	"fmt"
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/MEMOxiiii/PortalDF/packet"
)

// Servers lists all servers connected to the proxy.
type Servers struct{}

// Run requests the server list from the proxy and displays it.
func (Servers) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("This command can only be used by players.")
		return
	}
	if portalRef == nil || !portalRef.Connected() {
		o.Errorf("Not connected to the proxy.")
		return
	}

	senderUUID := p.UUID()

	_ = portalRef.RequestServerList(func(servers []packet.ServerEntry) {
		var parts []string
		for _, s := range servers {
			parts = append(parts, fmt.Sprintf("§e%s §7(%d players)", s.Name, s.PlayerCount))
		}
		msg := fmt.Sprintf("§aThere are §e%d §aservers connected to the proxy:\n%s", len(servers), strings.Join(parts, "\n"))
		sendToPlayer(senderUUID, msg)
	})
}

// Allow restricts this command to players only.
func (Servers) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}
