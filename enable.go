package portaldf

import (
	"log/slog"

	"github.com/MEMOxiiii/PortalDF/command"
	"github.com/df-mc/dragonfly/server"
)

// Enable connects srv to a Portal proxy using config and registers the /transfer, /server and /servers
// commands, along with stale-session cleanup for players the proxy transfers in. The connection and its
// automatic reconnect loop run in the background; call Close on the returned Portal to stop them.
//
// It is the fastest way to wire a Dragonfly server up to a Portal proxy:
//
//	portaldf.Enable(srv, portaldf.DefaultConfig())
func Enable(srv *server.Server, config Config) *Portal {
	return EnableWithLogger(srv, config, nil)
}

// EnableWithLogger behaves like Enable but logs through log instead of the default slog.Logger.
func EnableWithLogger(srv *server.Server, config Config, log *slog.Logger) *Portal {
	p := New(config, log)
	command.Register(p, srv)
	go p.Connect()
	return p
}
