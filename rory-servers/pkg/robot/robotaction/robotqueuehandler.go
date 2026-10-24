package robotaction

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotQueueHandler 機器人遊戲撮合管理
type RobotQueueHandler struct{}

// NewRobotQueueHandler 創建機器人遊戲撮合管理器
func NewRobotQueueHandler() *RobotQueueHandler {
	rgm := new(RobotQueueHandler)

	return rgm
}

// Start 啟動
func (rgm *RobotQueueHandler) Start(ctx *HandlerCtx) {
	scope := util.NewIntScope(ctx.RobotCfg.QueueCheckWaitSpanMin, ctx.RobotCfg.QueueCheckWaitSpanMax)
	ctx.TfManager.NewNamedRandomTicker(scope, func() {
		if ctx.ServerStatus.IsServerRuning() {
			rgm.robotGameMatch(ctx)
		}
	}, "RobotMatchGameTicker")
}

func (rgm *RobotQueueHandler) robotGameMatch(ctx *HandlerCtx) {
	//example.invalid("機器人配桌邏輯開始檢查...")

	tables := ctx.Room.GetRobotSortedTable()
	// 需要再次匹配
	needReMatch := false

	for _, table := range tables {
		tno := table.TableNo
		if !ctx.ServerStatus.IsServerRuning() {
			example.invalid("服務器停服,終止配桌.")
			return
		}

		urcount := int32(table.TableRealUserCount())

		if urcount == 0 && ctx.RobotCfg.StressTest != 1 {
			//example.invalid("機器人配桌邏輯檢查完成.")
			return
		}
		//AgainstGameing => 1 => false , 不可以觀戰
		//AgainstGameing => 0 => true , 可以觀戰
		if !ctx.RobotCfg.AgainstGameing && table.GameLogic.IsGaming() {
			continue
		}

		// 檢查桌子上的人數是否匹配撮合條件
		if urcount < ctx.RobotCfg.TableMinUserCount || urcount >= ctx.RobotCfg.TableMaxUserCount {
			continue
		}

		// 獲取一個空閒的機器人，防止機器人連續2次配桌到同一桌號
		freerobot := ctx.RobotManager.GetFreeRobot(table.TableNo)
		// 清除機器人上局遊戲桌號
		ctx.RobotManager.ClearFreeRobotLastTable(table.TableNo)

		if freerobot == nil {
			continue
		}

		// 獲取一個空閒的座位
		st := table.GetFreeSeat()
		if st == nil {
			continue
		}

		ctx.ActionController.UserSit(freerobot.UserInfo, int32(tno), int32(st.SeatNo))
		ctx.ActionController.UserReady(freerobot.UserInfo)

		if !table.GameLogic.IsGaming() {
			needReMatch = true
		}
	}

	//example.invalid("機器人配桌邏輯檢查完成...")

	// 如果還有遊戲未開局,隨機一定時間後再次配桌
	if needReMatch {
		ctx.TfManager.NewTimer(ctx.RobotCfg.QueueAgainWaitSpan(), func() {
			rgm.robotGameMatch(ctx)
		})
	}

}
