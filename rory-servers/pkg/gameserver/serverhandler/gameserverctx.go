package serverhandler

import (
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GameServerCtx Context
type GameServerCtx struct {
	PlayerManager    *player.Manager         // 用戶
	DBMsgQueue       *rpcmsg.Queue           // 數據庫消息隊列
	MsgServer        *msgserver.MsgServer    // 消息服務器
	ModuleServer     *module.Server          // 模塊服務器
	TimerFuncManager *event.TimerFuncManager // 延遲函數管理器
	Throttle         int                     // 並發數
}
