package msgserver

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// LogicHandler 服務器處邏輯理器
type LogicHandler interface {
	// 客戶端連接
	OnClientConnect(clt *network.ConnInfo)

	// 客戶端斷開
	OnClientDisconnect(clt *network.ConnInfo)

	// 收到數據
	OnClientData(data *network.ConnData)

	// 收到消息
	OnClientMessage(clt *network.ConnInfo, msg *msg.Message)

	// 斷開客戶端接口
	CloseClientFunc(f func(clt *network.ConnInfo))
}

// LogicHandlerGroup 服務器邏輯處理器組
type LogicHandlerGroup struct {
	innerLogicHandlers []LogicHandler
	closCltFunc        func(clt *network.ConnInfo)
}

// NewMsgServerLogicGroup 創建服務器組
func NewMsgServerLogicGroup() *LogicHandlerGroup {
	return new(LogicHandlerGroup)
}

// RegisterHandler 註冊邏輯處理器
func (slg *LogicHandlerGroup) RegisterHandler(h LogicHandler) *LogicHandlerGroup {
	slg.innerLogicHandlers = append(slg.innerLogicHandlers, h)
	h.CloseClientFunc(slg.closCltFunc)
	return slg
}

// OnClientConnect ?
func (slg *LogicHandlerGroup) OnClientConnect(clt *network.ConnInfo) {
	for _, lg := range slg.innerLogicHandlers {
		lg.OnClientConnect(clt)
	}
}

// OnClientDisconnect ?
func (slg *LogicHandlerGroup) OnClientDisconnect(clt *network.ConnInfo) {
	for _, lg := range slg.innerLogicHandlers {
		lg.OnClientDisconnect(clt)
	}
}

// OnClientData ?
func (slg *LogicHandlerGroup) OnClientData(data *network.ConnData) {
	for _, lg := range slg.innerLogicHandlers {
		lg.OnClientData(data)
	}
}

// OnClientMessage ?
func (slg *LogicHandlerGroup) OnClientMessage(clt *network.ConnInfo, cmsg *msg.Message) {
	for _, lg := range slg.innerLogicHandlers {
		lg.OnClientMessage(clt, cmsg)
	}
}

// CloseClientFunc ?
func (slg *LogicHandlerGroup) CloseClientFunc(f func(clt *network.ConnInfo)) {
	slg.closCltFunc = f
}
