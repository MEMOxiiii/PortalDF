package portaldf

import (
	"fmt"
	"net/url"
)

// Transport identifies the network transport the proxy should use to dial this server.
type Transport string

const (
	// TransportRakNet is the classic UDP transport. ServerAddress must be a "host:port" pair.
	TransportRakNet Transport = "raknet"
	// TransportNetherNet is Bedrock's WebRTC-based transport. ServerAddress must be the URL of this
	// server's HTTP(S) signaling endpoint (e.g. "http://127.0.0.1:19132"), matching whatever address
	// Dragonfly's own NetherNet listener is bound to.
	TransportNetherNet Transport = "nethernet"
)

// Config holds the configuration for connecting to a Portal proxy.
type Config struct {
	// ProxyAddress is the IP address of the Portal proxy (e.g. "127.0.0.1").
	ProxyAddress string
	// SocketPort is the port of the Portal proxy's communication socket (e.g. 19131).
	SocketPort int
	// Secret is the authentication secret. Must match the proxy's configured secret.
	Secret string
	// ServerName is the name this server will be identified as on the proxy (e.g. "Hub1", "SkyWars1").
	ServerName string
	// ServerAddress is the address of this server that the proxy should connect players to. Its format
	// depends on Transport: for TransportRakNet it is a "host:port" pair (e.g. "127.0.0.1:19132"); for
	// TransportNetherNet it is the URL of this server's HTTP(S) signaling endpoint.
	ServerAddress string
	// Transport is the network transport the proxy should use to dial this server. If empty,
	// TransportRakNet is used, matching every server registered before this field existed.
	Transport Transport
	// Group is the name of the group this server belongs to, used by group-aware load balancers on the
	// proxy to route players to the correct set of servers. Leave empty if the server does not belong to
	// a group.
	Group string
	// Weight controls how large a share of new players this server should receive relative to others in
	// the same group. A weight of 0 is treated by the proxy as 1, giving all servers in a group without an
	// explicit weight an even split.
	Weight uint32
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		ProxyAddress:  "127.0.0.1",
		SocketPort:    19131,
		Secret:        "",
		ServerName:    "Server1",
		ServerAddress: "127.0.0.1:19132",
		Transport:     TransportRakNet,
	}
}

// validate checks that ServerAddress is well-formed for Transport, catching the most common
// misconfiguration -- a "host:port" pair left over from TransportRakNet after switching Transport to
// TransportNetherNet without updating ServerAddress to match -- locally, instead of it being silently
// rejected by the proxy with a confusing URL-parse error minutes later, deep inside a health check or a
// player transfer.
func (c Config) validate() error {
	transport := c.Transport
	if transport == "" {
		transport = TransportRakNet
	}
	if c.ServerAddress == "" {
		return fmt.Errorf("ServerAddress must not be empty")
	}
	if transport != TransportNetherNet {
		return nil
	}
	u, err := url.Parse(c.ServerAddress)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("ServerAddress must be the full URL of this server's NetherNet signaling endpoint (e.g. %q), got %q", "http://host:port", c.ServerAddress)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("ServerAddress must use the \"http\" or \"https\" scheme for TransportNetherNet, got %q", c.ServerAddress)
	}
	return nil
}
