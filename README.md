# PortalDF

A Go library for connecting [Dragonfly](https://github.com/df-mc/dragonfly) servers to the [Portal](https://github.com/MEMOxiiii/portal) proxy via its socket communication protocol.

📖 **[Full documentation is in the Wiki](https://github.com/MEMOxiiii/PortalDF/wiki)**

---

## Features

- Connects to the Portal proxy via TCP socket, with automatic reconnection
- Transfer players, query player info, list servers, find a player anywhere on the network
- Player latency reporting, stale-session cleanup, and draining
- Built-in Dragonfly commands: `/transfer`, `/server`, `/servers`
- Works over RakNet or NetherNet, chosen per server

## Installation

```
go get github.com/MEMOxiiii/PortalDF
```

## Quick Start

`Enable` connects to the proxy and registers the built-in commands in a single call:

```go
package main

import (
	"github.com/MEMOxiiii/PortalDF"
	"github.com/df-mc/dragonfly/server"
)

func main() {
	// ... set up dragonfly server config ...
	srv := conf.New()
	srv.CloseOnProgramEnd()
	srv.Listen()

	portaldf.Enable(srv, portaldf.Config{
		ProxyAddress:  "127.0.0.1",
		SocketPort:    19131,
		Secret:        "your-secret",
		ServerName:    "Hub1",
		ServerAddress: "127.0.0.1:19132",
	})

	for p := range srv.Accept() {
		// handle players...
	}
}
```

See the [Wiki](https://github.com/MEMOxiiii/PortalDF/wiki) for every `Config` field, manual (non-`Enable`) setup, running over NetherNet, the full API (`TransferPlayer`, `FindPlayer`, `SetDraining`, etc.), and the commands this registers.

## Issues

If you encounter any problems, please [open an issue](https://github.com/MEMOxiiii/PortalDF/issues).
