package packet

// ProtocolVersion is the protocol version supported by the proxy. It must match the proxy's version.
const ProtocolVersion = 2

const (
	IDAuthRequest uint16 = iota
	IDAuthResponse
	IDRegisterServer
	IDTransferRequest
	IDTransferResponse
	IDPlayerInfoRequest
	IDPlayerInfoResponse
	IDServerListRequest
	IDServerListResponse
	IDFindPlayerRequest
	IDFindPlayerResponse
	IDUpdatePlayerLatency
	IDDisconnectPlayer
	IDSetServerDraining
)
