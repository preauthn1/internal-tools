package serverdefense

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// BlackIPDefense IP黑名單防護
type BlackIPDefense struct {
	*DefenserBase
}

// NewBlackIPDefense 創建IP黑名單防護
func NewBlackIPDefense() *BlackIPDefense {
	bd := new(BlackIPDefense)
	bd.DefenserBase = NewDefenserBase()

	return bd
}

// OnClientConnect 客戶端連接
func (dd *BlackIPDefense) OnClientConnect(clt *network.ConnInfo) {
	if dd.BlackIPList.CheckIP(clt.IPAddress) && !dd.checkWhiteIPList(clt.IPAddress) {
		example.invalid("客戶端%v連接,但被IP黑名單斷開", clt)
		dd.CloseConn(clt)
	}
}

// OnClientDisconnect 客戶端斷開
func (dd *BlackIPDefense) OnClientDisconnect(clt *network.ConnInfo) {}

// OnClientData 收到數據
func (dd *BlackIPDefense) OnClientData(data *network.ConnData) {}

// OnClientMessage 收到消息
func (dd *BlackIPDefense) OnClientMessage(clt *network.ConnInfo, msg *msg.Message) {
	if dd.BlackIPList.CheckIP(clt.IPAddress) && !dd.checkWhiteIPList(clt.IPAddress) {
		example.invalid("客戶端%v發送消息,但被IP黑名單斷開", clt)
		dd.CloseConn(clt)
	}
}
