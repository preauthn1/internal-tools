package robotaction

import (
	"math/rand"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotRateQueueHandler 一次數個機器人遊戲撮合管理
type RobotRateQueueHandler struct{}

// NewRobotRateQueueHandler 創建機器人遊戲撮合管理器
func NewRobotRateQueueHandler() *RobotRateQueueHandler {
	rgm := new(RobotRateQueueHandler)

	return rgm
}

// Start 啟動
func (rgm *RobotRateQueueHandler) Start(ctx *HandlerCtx) {
	scope := util.NewIntScope(ctx.RobotCfg.QueueCheckWaitSpanMin, ctx.RobotCfg.QueueCheckWaitSpanMax)
	ctx.TfManager.NewNamedRandomTicker(scope, func() {
		if ctx.ServerStatus.IsServerRuning() {
			rgm.robotGameMatch(ctx)
		}
	}, "RobotMatchGameTicker")
}

// 一次數個機器人遊戲撮合
func (rgm *RobotRateQueueHandler) robotGameMatch(ctx *HandlerCtx) {
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

		//SitGameing => 1 => false , 不可以觀戰
		//SitGameing => 0 => true , 可以觀戰
		if !ctx.RobotCfg.AgainstGameing && table.GameLogic.IsGaming() {
			continue
		}

		// 檢查桌子上的人數是否匹配撮合條件
		if urcount < ctx.RobotCfg.TableMinUserCount || urcount >= ctx.RobotCfg.TableMaxUserCount {
			continue
		}

		// 隨機獲取放入機器人數量
		randNum := rand.Int31n(100)

		var robotCount int
		for k, v := range ctx.RobotCfg.TableMatchRobotRate {
			if randNum <= int32(v) {
				robotCount = k
				break
			}
		}
		//mlog.System("tableno %d, randNum %d, robotCount %d", table.TableNo, randNum, robotCount)

		var standbyReadyRobot []*gamelogic.GameUser
		for i := 1; i <= robotCount; i++ {
			// 獲取一個空閒的機器人，防止機器人連續2次配桌到同一桌號
			freerobot := ctx.RobotManager.GetFreeRobot(table.TableNo)

			if freerobot == nil {
				continue
			}

			// 獲取一個空閒的座位
			st := table.GetFreeSeat()
			if st == nil {
				continue
			}
			//mlog.System("tableno %d 配桌, %v", table.TableNo, freerobot)
			ctx.ActionController.UserSit(freerobot.UserInfo, int32(tno), int32(st.SeatNo))
			standbyReadyRobot = append(standbyReadyRobot, freerobot.UserInfo)
			//ctx.ActionController.UserReady(freerobot.UserInfo)

			if !table.GameLogic.IsGaming() {
				needReMatch = true
			}
		}

		// 清除機器人上局遊戲桌號
		ctx.RobotManager.ClearFreeRobotLastTable(table.TableNo)

		for _, robotInfo := range standbyReadyRobot {
			ctx.ActionController.UserReady(robotInfo)
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
