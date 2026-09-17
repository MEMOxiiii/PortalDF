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
	// server's HTTP(S) signaling endpoint (e.g. "http://127.0.0.1:19132").
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
	// ServerAddress is this server's own address, in the format Transport requires (see its docs).
	ServerAddress string
	// Transport is the network transport the proxy should use to dial this server. Empty means
	// TransportRakNet.
	Transport Transport
	// Group is the load-balancer group this server belongs to. Leave empty for no group.
	Group string
	// Weight controls how large a share of new players this server gets relative to others in Group. 0 is
	// treated as 1 (an even split).
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

// validate checks ServerAddress is well-formed for Transport locally, before ever contacting the proxy.
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
	if u.Port() == "" {
		return fmt.Errorf("ServerAddress must include an explicit port, got %q", c.ServerAddress)
	}
	if u.Path != "" {
		return fmt.Errorf("ServerAddress must not have a path, got %q", c.ServerAddress)
	}
	return nil
}
