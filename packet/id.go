package packet

// ProtocolVersion is the protocol version supported by the proxy. It must match the proxy's version.
//
// v3 added the Transport field to RegisterServer.
const ProtocolVersion = 3

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
