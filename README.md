# PortalDF

A Go library for connecting [Dragonfly](https://github.com/df-mc/dragonfly) servers to the [Portal](https://github.com/paroxity/portal) proxy via its socket communication protocol.

This is the Dragonfly equivalent of [PortalPM](https://github.com/MEMOxiiii/PortalPM) for PocketMine-MP.

## Features

- Connects to the Portal proxy via TCP socket
- Automatic authentication and server registration
- Transfer players between servers
- Query player information (XUID, IP address)
- List all connected servers with player counts
- Find which server a player is on
- Receive player latency updates from the proxy
- Automatic reconnection on disconnect
- **Built-in dragonfly commands**: `/transfer`, `/server`, `/servers`

## Installation

```
go get github.com/MEMOxiiii/PortalDF
```

## Quick Start

```go
package main

import (
	"log/slog"

	"github.com/df-mc/dragonfly/server"
	"github.com/MEMOxiiii/PortalDF"
	portalcmd "github.com/MEMOxiiii/PortalDF/command"
)

func main() {
	// ... set up dragonfly server config ...
	srv := conf.New()
	srv.CloseOnProgramEnd()

	portal := portaldf.New(portaldf.Config{
		ProxyAddress:  "127.0.0.1",
		SocketPort:    19131,
		Secret:        "your-secret",
		ServerName:    "Hub1",
		ServerAddress: "127.0.0.1:19132",
	}, slog.Default())

	go portal.Connect()

	// Register /transfer, /server, /servers commands.
	portalcmd.Register(portal, srv)

	srv.Listen()
	for p := range srv.Accept() {
		// handle players...
	}
}
```

## Commands

The `command` sub-package provides ready-to-use dragonfly commands:

| Command | Description |
|---|---|
| `/transfer <server>` | Transfer yourself to another server |
| `/transfer <server> <player>` | Transfer another player to a server |
| `/server` | Check which server you are on |
| `/server <player>` | Check which server another player is on |
| `/servers` | List all servers connected to the proxy |

Register all commands with one call:

```go
import portalcmd "github.com/MEMOxiiii/PortalDF/command"

portalcmd.Register(portal, srv)
```

## API Usage

```go
import (
	"github.com/google/uuid"
	"github.com/MEMOxiiii/PortalDF"
	"github.com/MEMOxiiii/PortalDF/packet"
)

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

// Request player info (XUID, IP).
portal.RequestPlayerInfo(playerUUID, func(uid uuid.UUID, status byte, xuid string, address string) {
	log.Info("Player info", "xuid", xuid, "address", address)
})

// Handle latency updates from the proxy.
portal.SetLatencyHandler(func(uid uuid.UUID, latency int64) {
	log.Info("Player latency update", "uuid", uid, "latency_ms", latency)
})
```

## Configuration

| Field | Description | Default |
|---|---|---|
| `ProxyAddress` | IP address of the Portal proxy | `127.0.0.1` |
| `SocketPort` | Communication socket port | `19131` |
| `Secret` | Authentication secret (must match proxy) | `""` |
| `ServerName` | Server identifier on the proxy | `Server1` |
| `ServerAddress` | Address for proxy to connect players to | `127.0.0.1:19132` |

## Transfer Response Statuses

| Constant | Value | Meaning |
|---|---|---|
| `TransferResponseSuccess` | 0 | Player transferred successfully |
| `TransferResponseServerNotFound` | 1 | Target server not found on proxy |
| `TransferResponseAlreadyOnServer` | 2 | Player is already on that server |
| `TransferResponsePlayerNotFound` | 3 | Player could not be found |
| `TransferResponseError` | 4 | An error occurred (check error string) |

## Protocol

This library implements the Portal proxy's binary TCP socket protocol:
- 4-byte little-endian length prefix
- 2-byte little-endian packet ID header
- Payload serialized using gophertunnel's protocol.Reader/Writer

## Issues

If you encounter any problems, please [open an issue](https://github.com/MEMOxiiii/PortalDF/issues).
