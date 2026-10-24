package notifyhandler

import (
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// RefreshVerHandler 刷新版本號處理器
type RefreshVerHandler struct {
	verLoader hallutil.HallVerLoader
}

// NewRefreshVerHandler 創建版本號刷新處理器
func NewRefreshVerHandler(verLoader hallutil.HallVerLoader) *RefreshVerHandler {
	refreshVerHandler := new(RefreshVerHandler)
	refreshVerHandler.verLoader = verLoader

	return refreshVerHandler
}

// OnNotifyMessage ?
func (ver *RefreshVerHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, player *player.Manager) {
	mlog.System("收到更新平台版本號消息.開始拉取平台版本號...")
	ver.verLoader.LoadHallVer()
}

// IsHook ?
func (ver *RefreshVerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_RefreshVersionID)
}
