package serverhandler

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// TimerNameServerReg
const (
	TimerNameServerReg = "gameserver_serverhandler_gameserverreg"
)

// GameServerRegHandler gameserver 定時註冊
type GameServerRegHandler struct {
	reginterval  int
	servertype   int32
	serverid     int32
	svrRegTicker *event.TimerFunc
}

// NewGameServerRegHandler 創建gameserver定時註冊處理器
func NewGameServerRegHandler(interval int, serverid int32, servertype int32) *GameServerRegHandler {
	regh := new(GameServerRegHandler)
	regh.reginterval = interval
	regh.serverid = serverid
	regh.servertype = servertype

	return regh
}

// Name 處理器名稱
func (rh *GameServerRegHandler) Name() string {
	return "IntervalReg"
}

// Start 啟動處理器
func (rh *GameServerRegHandler) Start(ctx *GameServerCtx) {
	defer util.PanicErrStack()

	regFunc := func() {
		serverinfo := new(netproto.ServerRegInfo)
		serverinfo.ServerID = proto.Int32(rh.serverid)
		serverinfo.ServerType = proto.Int32(rh.servertype)

		rpcinfo := netproto.CreateEmptyRPCInfo(rh.serverid)

		msgReg := new(msg.Message)
		msgReg.BClassID = int32(netproto.MessageBClassID_ServerMgm)
		msgReg.SClassID = int32(netproto.ServerMgmClassID_GameServerRegister)
		msgReg.MsgData = serverinfo

		ctx.DBMsgQueue.PushMessage(rpcinfo, false, msgReg)
	}

	rh.svrRegTicker = ctx.TimerFuncManager.NewNamedTicker(int32(rh.reginterval), regFunc, TimerNameServerReg)
}

// Stop 停止
func (rh *GameServerRegHandler) Stop() {
	if rh.svrRegTicker != nil {
		rh.svrRegTicker.Stop()
	}
}
