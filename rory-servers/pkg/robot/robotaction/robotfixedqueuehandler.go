package robotaction

import (
	"math/rand"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotFixedQueueHandler 一次數個機器人遊戲撮合管理
type RobotFixedQueueHandler struct {
}

// NewRobotFixedQueueHandler 創建機器人遊戲撮合管理器
func NewRobotFixedQueueHandler() *RobotFixedQueueHandler {
	rgm := new(RobotFixedQueueHandler)

	return rgm
}

// Start 啟動
func (rgm *RobotFixedQueueHandler) Start(ctx *HandlerCtx) {
	scope := util.NewIntScope(ctx.RobotCfg.QueueCheckWaitSpanMin, ctx.RobotCfg.QueueCheckWaitSpanMax)
	ctx.TfManager.NewNamedRandomTicker(scope, func() {
		if ctx.ServerStatus.IsServerRuning() {
			rgm.robotGameMatch(ctx)
		}
	}, "RobotMatchGameTicker")
}

// 一桌固定機器人數量上限
func (rgm *RobotFixedQueueHandler) robotGameMatch(ctx *HandlerCtx) {
	//example.invalid("机器人配桌逻辑开始检查...")

	tables := ctx.Room.GetRobotSortedTable()
	//需要再次匹配
	needReMatch := false

	for _, table := range tables {
		tno := table.TableNo
		if !ctx.ServerStatus.IsServerRuning() {
		}

		urcount := int32(table.GameLogic.TableRealUserCount())

		if urcount == 0 && ctx.RobotCfg.StressTest != 1 {
			//example.invalid("机器人配桌逻辑检查完成.")
			return
		}

		//SitGameing => 1 => false , 不可以觀戰
		//SitGameing => 0 => true , 可以觀戰
		if !ctx.RobotCfg.AgainstGameing && table.GameLogic.IsGaming() {
			continue
		}

		//检查桌子上的人数是否匹配撮合条件
		if urcount < ctx.RobotCfg.TableMinUserCount || urcount >= ctx.RobotCfg.TableMaxUserCount {
			continue
		}

		//隨機獲取放入機器人數量
		randNum := rand.Int31n(100)

		var robotCount int
		for k, v := range ctx.RobotCfg.TableMatchRobotRate {
			if randNum < int32(v) && int(table.TableRealRobotCount()) < k {
				robotCount = k
				break
			}
		}
		mlog.Debug("tableno %d,  addRobotCount %d TableRealRobotCount[%v]",
			table.TableNo, robotCount, table.TableRealRobotCount())

		var standbyReadyRobot []*gamelogic.GameUser
		for i := 1; i <= robotCount; i++ {
			//获取一个空闲的机器人，防止机器人连续2次配桌到同一桌号
			freerobot := ctx.RobotManager.GetFreeRobot(table.TableNo)

			if freerobot == nil {
				continue
			}

			//获取一个空闲的座位
			st := new(gamelogic.Seat)
			if ctx.Room.GameInfo.GameID == 14 { //BY
				st = table.GetBYFreeSeat()
			} else {
				st = table.GetFreeSeat()
			}

			if st == nil {
				continue
			}
			mlog.Debug("tableno %d 配桌, %v", table.TableNo, freerobot)
			ctx.ActionController.UserSit(freerobot.UserInfo, int32(tno), int32(st.SeatNo))
			standbyReadyRobot = append(standbyReadyRobot, freerobot.UserInfo)
			//ctx.ActionController.UserReady(freerobot.UserInfo)

			if !table.GameLogic.IsGaming() {
				needReMatch = true
			}
		}

		//清除机器人上局游戏桌号
		ctx.RobotManager.ClearFreeRobotLastTable(table.TableNo)

		for _, robotInfo := range standbyReadyRobot {
			ctx.ActionController.UserReady(robotInfo)
		}

	}

	//example.invalid("机器人配桌逻辑检查完成...")

	//如果还有游戏未开局,随机一定时间后再次配桌
	if needReMatch {
		ctx.TfManager.NewTimer(ctx.RobotCfg.QueueAgainWaitSpan(), func() {
			rgm.robotGameMatch(ctx)
		})
	}

}
