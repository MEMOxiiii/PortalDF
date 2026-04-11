package command

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
)

// sendToPlayer sends a message to a player by UUID using the server's entity handle.
// This is safe to call from any goroutine (e.g. async Portal callbacks).
func sendToPlayer(uid uuid.UUID, msg string) {
	if serverRef == nil {
		return
	}
	handle, ok := serverRef.Player(uid)
	if !ok {
		return
	}
	handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
		if p, ok := e.(*player.Player); ok {
			p.Message(msg)
		}
	})
}
