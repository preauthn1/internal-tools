package roomlogic

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/serverhandler"
)

// RefreshVerHandler 刷新版本號處理器
type RefreshVerHandler struct {
	verprovider    *RPCGameVerProvider
	regHallHandler *serverhandler.HallServerRegHandler
}

// NewRefreshVerHandler 創建刷新版本號處理器
func NewRefreshVerHandler(verprovider *RPCGameVerProvider, regHallHandler *serverhandler.HallServerRegHandler) *RefreshVerHandler {
	gameverHandler := new(RefreshVerHandler)
	gameverHandler.verprovider = verprovider
	gameverHandler.regHallHandler = regHallHandler

	return gameverHandler
}

// IsHook 是否攔截
func (h *RefreshVerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GameVerRetID)
}

// OnRPCMessage 收到RPC封包
func (h *RefreshVerHandler) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameVerInfo{}
		return mm
	})

	mm := rmm.(*netproto.GameVerInfo)
	h.verprovider.gsreg.gameServerRegRet.AndroidVersion = mm.AndroidVersion
	h.verprovider.gsreg.gameServerRegRet.IOSVersion = mm.IOSVersion
	// 當收到更新版本信息時向大廳更新資料
	h.regHallHandler.RefreshGameVer(mm)
	mlog.NewDefaultLogger().System("收到更新版本信息.android ver:%s,ios ver:%s", mm.GetAndroidVersion(), mm.GetIOSVersion())
}
