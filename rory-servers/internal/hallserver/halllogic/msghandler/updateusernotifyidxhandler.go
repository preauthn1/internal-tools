package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// UpdateUserNotifyIdxHandler ?
type UpdateUserNotifyIdxHandler struct {
	server        *gameserver.GameServer
	PlayerManager *player.Manager
}

// NewUpdateUserNotifyIdxHandler ?
func NewUpdateUserNotifyIdxHandler(server *gameserver.GameServer, PlayerManager *player.Manager) *UpdateUserNotifyIdxHandler {
	handler := new(UpdateUserNotifyIdxHandler)
	handler.server = server
	handler.PlayerManager = PlayerManager
	return handler
}

// OnMessage ?
func (h *UpdateUserNotifyIdxHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ClientUpdateUserNotifyIdx{}
		return mm
	})
	msgData := rmm.(*netproto.ClientUpdateUserNotifyIdx)
	//mm := &netproto.ClientUpdateUserNotifyIdx{}
	//
	//err := proto.Unmarshal(dd, mm)
	//if err != nil {
	//	ctx.Error("ClientUpdateUserNotifyIdx 数据包解析错误,%v", err)
	//	return nil
	//}

	if pr := h.PlayerManager.GetPlayerByConnID(clt.ConnID()); pr != nil || !pr.IsGuest() {
		h.updateToNotify(pr.UserID, msgData.GetIdx())
	} else {
		ctx.Error("找不到連線用戶資料 %v", pr)
	}
	//
	//rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
	//	mm := &netproto.ClientUpdateUserNotifyIdx{}
	//	return mm
	//})
	//msgData := rmm.(*netproto.ClientUpdateUserNotifyIdx)
	//mlog.System("=========== Recieved ClientUpdateUserNotifyIdx UserID %d, Idx %d", rpcinfo.GetUserID(), msgData.GetIdx())
	//h.updateToNotify(rpcinfo.GetUserID(), msgData.GetIdx())
	return nil
}

// IsHook /
func (h *UpdateUserNotifyIdxHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_UpdateUserNotifyIdxC2HID)
}

func (h *UpdateUserNotifyIdxHandler) updateToNotify(userid int32, idx int64) {
	msg := new(msg.Message)
	msg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	msg.SClassID = int32(netproto.NotifyServerClassID_UpdateUserNotifyIdxH2NID)
	data := new(netproto.UpdateUserNotifyIdx)
	data.UserID = proto.Int32(userid)
	data.Idx = proto.Int64(idx)
	msg.MsgData = data
	h.server.NotifyHandler.TCPMsgClient.SendMessage(msg)
}
