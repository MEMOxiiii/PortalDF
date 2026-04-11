package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/google/uuid"
	"github.com/hexomc/portaldf/packet"
)

// TransferSelf transfers the command sender to the specified server.
type TransferSelf struct {
	Server string `cmd:"server"`
}

// Run executes the transfer command for the sender.
func (c TransferSelf) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
	p, ok := src.(*player.Player)
	if !ok {
		o.Errorf("This command can only be used by players.")
		return
	}
	if portalRef == nil || !portalRef.Connected() {
		o.Errorf("Not connected to the proxy.")
		return
	}

	playerUUID := p.UUID()
	server := c.Server
	o.Printf("§eTransferring to %s...", server)

	_ = portalRef.TransferPlayer(playerUUID, server, func(_ uuid.UUID, status byte, errMsg string) {
		switch status {
		case packet.TransferResponseSuccess:
			sendToPlayer(playerUUID, "§aTransferred to "+server)
		case packet.TransferResponseServerNotFound:
			sendToPlayer(playerUUID, "§cServer '"+server+"' not found.")
		case packet.TransferResponseAlreadyOnServer:
			sendToPlayer(playerUUID, "§cYou are already on that server.")
		case packet.TransferResponsePlayerNotFound:
			sendToPlayer(playerUUID, "§cPlayer not found.")
		case packet.TransferResponseError:
			sendToPlayer(playerUUID, "§cTransfer error: "+errMsg)
		}
	})
}

// Allow restricts this command to players only.
func (TransferSelf) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}

// TransferOther transfers another player (by name) to the specified server.
type TransferOther struct {
	Server string `cmd:"server"`
	Player string `cmd:"player"`
}

// Run executes the transfer command targeting another player.
func (c TransferOther) Run(src cmd.Source, o *cmd.Output, tx *world.Tx) {
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
	server := c.Server
	targetName := c.Player

	// Try local player first.
	if serverRef != nil {
		if handle, found := serverRef.PlayerByName(targetName); found {
			var targetUUID uuid.UUID
			handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
				if tp, ok := e.(*player.Player); ok {
					targetUUID = tp.UUID()
				}
			})
			if targetUUID != uuid.Nil {
				doTransfer(senderUUID, targetUUID, server)
				o.Printf("§eTransferring %s to %s...", targetName, server)
				return
			}
		}
	}

	// Not local — search across the proxy network.
	o.Printf("§eSearching for %s...", targetName)
	_ = portalRef.FindPlayer(uuid.Nil, targetName, func(uid uuid.UUID, name string, online bool, _ string) {
		if !online {
			sendToPlayer(senderUUID, "§cPlayer '"+targetName+"' not found.")
			return
		}
		doTransfer(senderUUID, uid, server)
	})
}

// Allow restricts this command to players only.
func (TransferOther) Allow(src cmd.Source) bool {
	_, ok := src.(*player.Player)
	return ok
}

// doTransfer issues a transfer request and notifies the sender of the result.
func doTransfer(senderUUID, targetUUID uuid.UUID, server string) {
	_ = portalRef.TransferPlayer(targetUUID, server, func(_ uuid.UUID, status byte, errMsg string) {
		switch status {
		case packet.TransferResponseSuccess:
			sendToPlayer(senderUUID, "§aPlayer transferred to "+server)
		case packet.TransferResponseServerNotFound:
			sendToPlayer(senderUUID, "§cServer '"+server+"' not found.")
		case packet.TransferResponseAlreadyOnServer:
			sendToPlayer(senderUUID, "§cPlayer is already on that server.")
		case packet.TransferResponsePlayerNotFound:
			sendToPlayer(senderUUID, "§cPlayer not found.")
		case packet.TransferResponseError:
			sendToPlayer(senderUUID, "§cTransfer error: "+errMsg)
		}
	})
}
