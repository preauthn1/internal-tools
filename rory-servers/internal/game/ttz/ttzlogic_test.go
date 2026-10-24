package ttz

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestGameRound(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 10000
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewTtzLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewTTZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*TtzLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	/*
		待完善事項
		1. 目前無法使用撮合(真實用戶)或配桌(機器人)逻辑, 需要主動安排入座
	*/
	_, robots := fixtures.GetGameUsers()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	//Test Round 1 正常開始
	robotManager := tgrs.RobotService.GetRobotManager()

	if robotManager.RobotCount() != 0 {
		t.Errorf("初始化, 機器人數量應為0, RobotCount=>%d", robotManager.RobotCount())
	}

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robotManager.AddRobot(robot)
	}
	if robotManager.RobotCount() != 3 {
		t.Errorf("加載機器人後數量應為3, RobotCount=>%d", robotManager.RobotCount())
	}

	//users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)
	robots[2].UserState.OnUserQueue(robots[2], tgrs.Room)

	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(robots[1])
	glc.UserReady(robots[2])

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Error("ddz game end")
	}

}
