package notifyhandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// UserNotifyHandler ?
type UserNotifyHandler struct{}

// NewUserNotifyHandler ?
func NewUserNotifyHandler() *UserNotifyHandler {
	attr := new(UserNotifyHandler)

	return attr
}

// OnNotifyMessage ?
func (h *UserNotifyHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	defer util.PanicErrStack()

	mmbyte := m.MsgData.([]byte)
	mm := &netproto.UserNotifyList{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		mlog.Error("解析錯誤 UserNotifyHandler")
	}
	for _, info := range mm.InfoList {
		if pr := player.GetPlayerByUserID(info.GetUserID()); pr != nil && pr.IsOnline {
			msg := new(msg.Message)
			msg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
			msg.SClassID = int32(netproto.NotifyServerClassID_UserNotifyRetH2CID)
			msg.MsgData = info
			ctx.Server.SendMessageByConnID(pr.Conn.ConnID(), msg)
		}
	}
}

// IsHook ?
func (h *UserNotifyHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_UserNotifyRetN2HID)
}
