# PortalDF

A Go library for connecting [Dragonfly](https://github.com/df-mc/dragonfly) servers to the [Portal](https://github.com/paroxity/portal) proxy via its socket communication protocol.

This is the Dragonfly equivalent of [PortalPM](https://github.com/paroxity/portal-pm) for PocketMine-MP.

## Features

- Connects to the Portal proxy via TCP socket
- Automatic authentication and server registration
- Transfer players between servers
- Query player information (XUID, IP address)
- List all connected servers with player counts
- Find which server a player is on
- Receive player latency updates from the proxy
- Automatic reconnection on disconnect

## Usage

```go
package main

import (
	"log/slog"
	"os"

	"github.com/google/uuid"
	"github.com/hexomc/portaldf"
	"github.com/hexomc/portaldf/packet"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))

	portal := portaldf.New(portaldf.Config{
		ProxyAddress:  "127.0.0.1",
		SocketPort:    19131,
		Secret:        "your-secret",
		ServerName:    "Hub1",
		ServerAddress: "127.0.0.1:19132",
	}, log)

	// Start connection in a goroutine (it blocks and auto-reconnects).
	go portal.Connect()

	// Transfer a player to another server.
	portal.TransferPlayer(playerUUID, "SkyWars1", func(uid uuid.UUID, status byte, err string) {
		if status == packet.TransferResponseSuccess {
			log.Info("Player transferred successfully")
		}
	})

	// Get a list of all servers.
	portal.RequestServerList(func(servers []packet.ServerEntry) {
		for _, s := range servers {
			log.Info("Server", "name", s.Name, "players", s.PlayerCount)
		}
	})

	// Find a player across all servers.
	portal.FindPlayer(uuid.Nil, "PlayerName", func(uid uuid.UUID, name string, online bool, server string) {
		if online {
			log.Info("Player found", "name", name, "server", server)
		}
	})

	// Handle latency updates from the proxy.
	portal.SetLatencyHandler(func(uid uuid.UUID, latency int64) {
		log.Info("Player latency update", "uuid", uid, "latency_ms", latency)
	})
}
```

## Configuration

| Field | Description | Default |
|---|---|---|
| `ProxyAddress` | IP address of the Portal proxy | `127.0.0.1` |
| `SocketPort` | Communication socket port | `19131` |
| `Secret` | Authentication secret (must match proxy) | `""` |
| `ServerName` | Server identifier on the proxy | `Server1` |
| `ServerAddress` | Address for proxy to connect players to | `127.0.0.1:19132` |

## Protocol

This library implements the Portal proxy's binary TCP socket protocol:
- 4-byte little-endian length prefix
- 2-byte little-endian packet ID header
- Payload serialized using gophertunnel's protocol.Reader/Writer
