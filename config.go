package portaldf

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
	// ServerAddress is the address of this server that the proxy should connect players to,
	// in the format "ip:port" (e.g. "127.0.0.1:19132").
	ServerAddress string
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
	}
}
