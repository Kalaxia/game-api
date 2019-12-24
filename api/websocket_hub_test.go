package api

type HubMock struct {
	Hub
}

func InitWsHubMock() {
	WsHub = &HubMock{}
}


func (h *HubMock) Run() {
	return
}

func (h *HubMock) sendTo(player *Player, message *WsMessage) {
	return
}

func (h *HubMock) sendBroadcast(message *WsMessage) {
	return
}