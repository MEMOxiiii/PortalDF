package packet

import "github.com/google/uuid"

// TransferCallback is called when a transfer response is received from the proxy.
// status is one of the TransferResponse constants. err is only set when status == TransferResponseError.
type TransferCallback func(playerUUID uuid.UUID, status byte, err string)

// PlayerInfoCallback is called when a player info response is received from the proxy.
type PlayerInfoCallback func(playerUUID uuid.UUID, status byte, xuid string, address string)

// ServerListCallback is called when a server list response is received from the proxy.
type ServerListCallback func(servers []ServerEntry)

// FindPlayerCallback is called when a find player response is received from the proxy.
type FindPlayerCallback func(playerUUID uuid.UUID, playerName string, online bool, server string)

// LatencyHandler is called when the proxy sends a player latency update.
type LatencyHandler func(playerUUID uuid.UUID, latency int64)

// DisconnectPlayerHandler is called when the proxy asks this server to disconnect any existing session for
// the named player, because the player is about to be transferred here and may have a stale session left
// over from a previous connection.
type DisconnectPlayerHandler func(playerName string)
