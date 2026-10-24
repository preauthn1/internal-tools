package msghandler

import (
	"fmt"

	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// AllocGameServerHandler ?
type AllocGameServerHandler struct {
	svrProvider hallutil.GameServerProvider
}

// NewAllocGameServerHandler 創建分配遊戲服務器處理器對象
func NewAllocGameServerHandler(svrProvider hallutil.GameServerProvider) *AllocGameServerHandler {
	allochandler := new(AllocGameServerHandler)
	allochandler.svrProvider = svrProvider

	return allochandler
}

// IsHook ?
func (alloch *AllocGameServerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_RequestGameServerAddrID)
}

// OnMessage ?
func (alloch *AllocGameServerHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	dd := m.MsgData.([]byte)
	mm := &netproto.RequestGameServerAddrInfo{}

	err := proto.Unmarshal(dd, mm)

	if err != nil {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[AllocGameServerHandler]客戶端[%v]發送請求分配遊戲服務器數據包解析錯誤[%v] 將關閉客戶端連接conn:%v", clt, err, clt)).Error("HallServer壓測時")
		ctx.Error("客戶端%v發送請求分配遊戲服務器數據包解析錯誤,%v", clt, err)

		ctx.Server.CloseConn(clt)
		return nil
	}

	//TODO client改為直連遊戲port, 不再透過大廳二次取得遊戲port 位置

	return []*msg.Message{nil}
}
