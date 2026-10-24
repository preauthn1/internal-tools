package notifyhandler

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// AnnNotifyHandler 公告通知處理器
type AnnNotifyHandler struct{}

// NewAnnNotifyHandler ?
func NewAnnNotifyHandler() *AnnNotifyHandler {
	ann := new(AnnNotifyHandler)

	return ann
}

// OnNotifyMessage ?
func (ann *AnnNotifyHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	ctx.Server.SendAllMessage(m)

	example.invalid("廣播公告")
}

// IsHook ?
func (ann *AnnNotifyHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_AnnID)
}
