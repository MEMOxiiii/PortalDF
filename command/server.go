package command

import (
	"strings"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// ServerSelf shows which server the command sender is on.
type ServerSelf struct{}

// Run outputs the current server name from the Portal config.
func (ServerSelf) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	if _, ok := src.(*player.Player); !ok {
		o.Errorf("This command can only be used by players.")
		return
	}
	if portalRef == nil {
		o.Errorf("Not connected to the proxy.")
		return
	}
	o.Printf("§aYou are currently on §e%s", portalRef.ServerName())
}

// Allow restricts this command to players only.
func (ServerSelf) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}

// ServerOther shows which server another player is on.
type ServerOther struct {
	Player string `cmd:"player"`
}

// Run looks up the target player on the proxy and reports their server.
func (c ServerOther) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
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
	targetName := c.Player

	o.Printf("§eSearching for %s...", targetName)
	_ = portalRef.FindPlayer(uuid.Nil, targetName, func(_ uuid.UUID, name string, online bool, server string) {
		if !online {
			sendToPlayer(senderUUID, "§cPlayer '"+targetName+"' not found.")
			return
		}
		if strings.EqualFold(p.Name(), name) {
			sendToPlayer(senderUUID, "§aYou are currently on §e"+server)
		} else {
			sendToPlayer(senderUUID, "§a"+name+" is currently on §e"+server)
		}
	})
}

// Allow restricts this command to players only.
func (ServerOther) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}
