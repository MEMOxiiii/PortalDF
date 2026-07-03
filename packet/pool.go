package packet

// Pool is a map holding packets indexed by a packet ID.
type Pool map[uint16]Packet

// NewPool returns a new pool with all supported packets.
func NewPool() Pool {
	return Pool{
		IDAuthRequest:         &AuthRequest{},
		IDAuthResponse:        &AuthResponse{},
		IDRegisterServer:      &RegisterServer{},
		IDTransferRequest:     &TransferRequest{},
		IDTransferResponse:    &TransferResponse{},
		IDPlayerInfoRequest:   &PlayerInfoRequest{},
		IDPlayerInfoResponse:  &PlayerInfoResponse{},
		IDServerListRequest:   &ServerListRequest{},
		IDServerListResponse:  &ServerListResponse{},
		IDFindPlayerRequest:   &FindPlayerRequest{},
		IDFindPlayerResponse:  &FindPlayerResponse{},
		IDUpdatePlayerLatency: &UpdatePlayerLatency{},
		IDDisconnectPlayer:    &DisconnectPlayer{},
		IDSetServerDraining:   &SetServerDraining{},
	}
}
