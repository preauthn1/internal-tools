package roomlogic

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// NotifyRefreshVerHandler 刷新版本號處理器
type NotifyRefreshVerHandler struct {
	rpcmsgqueue *rpcmsg.Queue
	gameid      int32
	serverid    int32
}

// NewNotifyRefreshVerHandler 創建版本號刷新處理器
func NewNotifyRefreshVerHandler(rpcmsgqueue *rpcmsg.Queue, gameid int32, serverid int32) *NotifyRefreshVerHandler {

	refreshVerHandler := new(NotifyRefreshVerHandler)
	refreshVerHandler.rpcmsgqueue = rpcmsgqueue
	refreshVerHandler.gameid = gameid
	refreshVerHandler.serverid = serverid

	return refreshVerHandler
}

// OnNotifyMessage ?
func (ver *NotifyRefreshVerHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	mlog.System("收到更新版本號消息.開始拉取版本號...")
	svrmsg := new(msg.Message)
	svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	svrmsg.SClassID = int32(netproto.DBServerClassID_GameVerID)

	svrdata := new(netproto.GameVer)
	svrdata.GameID = proto.Int32(ver.gameid)
	svrdata.ServerID = proto.Int32(ver.serverid)

	svrmsg.MsgData = svrdata

	ver.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(ver.serverid), false, svrmsg)
}

// IsHook ?
func (ver *NotifyRefreshVerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_RefreshVersionID)
}
