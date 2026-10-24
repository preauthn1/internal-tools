package serverhandler

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// TimerNameLogonCheck
const (
	TimerNameLogonCheck = "gameserver_serverhandler_logoncheck"
)

// LogonCheckHandler 用戶登錄超時處理器
type LogonCheckHandler struct {
	timespan         int
	cgx              *GameServerCtx
	logonCheckTicker *event.TimerFunc
}

// NewLogonCheckHandler 創建登錄超時檢查處理器
func NewLogonCheckHandler(timespan int) *LogonCheckHandler {
	lh := new(LogonCheckHandler)
	lh.timespan = timespan

	return lh
}

// Name 處理器名稱
func (lh *LogonCheckHandler) Name() string {
	return "LogonCheck"
}

// Start 啟動處理器
func (lh *LogonCheckHandler) Start(cgx *GameServerCtx) {
	lh.cgx = cgx

	lh.logonCheckTicker = cgx.TimerFuncManager.NewNamedTicker(int32(lh.timespan), func() {
		kickguestusers := make([]*player.Player, 0)

		guestusers := cgx.PlayerManager.UnLogonPlayer()
		for _, u := range guestusers {
			if int(time.Now().Sub(u.Conn.ConnTime).Seconds()) >= lh.timespan {
				kickguestusers = append(kickguestusers, u)
			}
		}

		for _, ku := range kickguestusers {
			example.invalid("客戶端%v超時未登錄，斷開連接", ku.Conn)
			cgx.MsgServer.CloseConn(ku.Conn)
		}

	}, TimerNameLogonCheck)

}

// Stop 停止
func (lh *LogonCheckHandler) Stop() {
	if lh.logonCheckTicker != nil {
		lh.logonCheckTicker.Stop()
	}
}
