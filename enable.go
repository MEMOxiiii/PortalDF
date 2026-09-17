package portaldf

import (
	"log/slog"

	"github.com/MEMOxiiii/PortalDF/command"
	"github.com/df-mc/dragonfly/server"
)

// Enable connects srv to a Portal proxy and registers /transfer, /server and /servers, running in the
// background. Call Close on the returned Portal to stop it.
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
