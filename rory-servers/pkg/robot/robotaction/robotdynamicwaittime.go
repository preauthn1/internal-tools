package robotaction

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RobotDynamicWaitTimeByRealUserHandler 當真人超過某數量時調整機器人配桌時間增加
type RobotDynamicWaitTimeByRealUserHandler struct {
}

// NewRobotDynamicWaitTimeByRealUserHandler 創建機器人遊戲撮合管理器
func NewRobotDynamicWaitTimeByRealUserHandler() *RobotDynamicWaitTimeByRealUserHandler {
	rgm := new(RobotDynamicWaitTimeByRealUserHandler)

	return rgm
}

// Start 啟動
func (rgm *RobotDynamicWaitTimeByRealUserHandler) Start(ctx *HandlerCtx) {
	scope := util.NewIntScope(ctx.RobotCfg.QueueCheckWaitSpanMin, ctx.RobotCfg.QueueCheckWaitSpanMax)
	ctx.TfManager.NewNamedRandomTicker(scope, func() {
		if ctx.ServerStatus.IsServerRuning() && ctx.Room.GetRoomRealUserCount() <= ctx.RobotCfg.TableMatchRobotsWiatSpan.UserCount {
			rbsqueue := NewRobotQueueHandler()
			rbsqueue.robotGameMatch(ctx)
		}
	}, "RobotDynamicWaitTimeByRealUserTicker")

	//超過人數增加搓合時間
	scope2 := util.NewIntScope(ctx.RobotCfg.QueueCheckWaitSpanMin+ctx.RobotCfg.TableMatchRobotsWiatSpan.WaitSpan, ctx.RobotCfg.QueueCheckWaitSpanMax+ctx.RobotCfg.TableMatchRobotsWiatSpan.WaitSpan)
	ctx.TfManager.NewNamedRandomTicker(scope2, func() {
		if ctx.ServerStatus.IsServerRuning() && ctx.Room.GetRoomRealUserCount() > ctx.RobotCfg.TableMatchRobotsWiatSpan.UserCount {
			rgm.robotGameMatch(ctx)
		}
	}, "RobotDynamicWaitTimeByRealUserTicker")
}

func (rgm *RobotDynamicWaitTimeByRealUserHandler) robotGameMatch(ctx *HandlerCtx) {
	//example.invalid("机器人配桌逻辑开始检查...")

	tables := ctx.Room.GetRobotSortedTable()
	//需要再次匹配
	needReMatch := false

	for _, table := range tables {
		tno := table.TableNo
		if !ctx.ServerStatus.IsServerRuning() {
			example.invalid("服务器停服,终止配桌.")
			return
		}

		urcount := int32(table.TableRealUserCount())

		if urcount == 0 && ctx.RobotCfg.StressTest != 1 {
			//example.invalid("机器人配桌逻辑检查完成.")
			return
		}
		//AgainstGameing => 1 => false , 不可以觀戰
		//AgainstGameing => 0 => true , 可以觀戰
		if !ctx.RobotCfg.AgainstGameing && table.GameLogic.IsGaming() {
			continue
		}

		//检查桌子上的人数是否匹配撮合条件
		if urcount < ctx.RobotCfg.TableMinUserCount || urcount >= ctx.RobotCfg.TableMaxUserCount {
			continue
		}

		//获取一个空闲的机器人，防止机器人连续2次配桌到同一桌号
		freerobot := ctx.RobotManager.GetFreeRobot(table.TableNo)
		//清除机器人上局游戏桌号
		ctx.RobotManager.ClearFreeRobotLastTable(table.TableNo)

		if freerobot == nil {
			continue
		}

		//获取一个空闲的座位
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

	//example.invalid("机器人配桌逻辑检查完成...")

	//如果还有游戏未开局,随机一定时间后再次配桌
	if needReMatch {
		ctx.TfManager.NewTimer(ctx.RobotCfg.QueueAgainRealUserWaitSpan(), func() {
			rgm.robotGameMatch(ctx)
		})
	}

}
