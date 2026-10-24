package qznn

import (
	"os"
	"path"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//func TestGameRound(t *testing.T) {
//	endChannel := make(chan bool)
//	conf.LogLevel = 0
//	conf.TimerFastMode = 10000
//
//	roomserver.GameLogicCreator = func(t *gamelogic.tableManager, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewQZNNGameLogic(t, tfManager, cfgs)
//	}
//	roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
//		return NewQZNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
//	}
//
//	tgrs := fixtures.NewTestGameRoomServer()
//	glc := tgrs.GameLogicManager.GetGameLogic(0).(*QZNNGameLogic)
//
//	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.GameLogicBase)
//	glc.RegisterAfterGameEndFunc(func() {
//		endChannel <- true
//	})
//	/*
//		待完善事項
//		1. 目前無法使用撮合(真實用戶)或配桌(機器人)逻辑, 需要主動安排入座
//	*/
//	_, robots := fixtures.GetGameUsers()
//
//	//Test Round 1 正常開始
//	robotManager := tgrs.RobotService.GetRobotManager()
//
//	if robotManager.RobotCount() != 0 {
//		t.Errorf("初始化, 機器人數量應為0, RobotCount=>%d", robotManager.RobotCount())
//	}
//
//	//機器人實例化
//	for i := 0; i < 3; i++ {
//		robot := new(robotcomm.Robot)
//		robot.UserInfo = robots[i]
//		robot.RebotLevel = 6
//		robotManager.AddRobot(robot)
//	}
//	if robotManager.RobotCount() != 3 {
//		t.Errorf("加載機器人後數量應為3, RobotCount=>%d", robotManager.RobotCount())
//	}
//
//	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
//	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
//	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)
//	robots[2].UserState.OnUserQueue(robots[2], tgrs.Room)
//
//	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//	glc.UserEnter(robots[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//
//	glc.UserReady(robots[0])
//	glc.UserReady(robots[1])
//	glc.UserReady(robots[2])
//	select {
//	case _ = <-endChannel:
//		mlog.NewDefaultLogger().Error("game end")
//	}
//
//}

func TestStartGameLogicServer(t *testing.T) {

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewQZNNGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewQZNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	fixtures.StartGameLogicByTest()
}

func NewQzniuniuLogicByTest() *QZNNGameLogic {
	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	rumanager := gamelogic.NewRoomUserManager(new(player.Manager))
	room := gamelogic.NewRoom(servercfg, rumanager)
	table := room.GetTableByNo(0)

	evtdispatcher := event.NewAsyncEventDispatcher(1)
	tfManager := event.NewTimerFuncManager(evtdispatcher)

	glc := new(QZNNGameLogic)

	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, cfgs)
	parseNNConfig(cfgs)
	glc.RegisterAfterGameEndFunc(func() {
		//通知用户抢庄倍数
		for _, ur := range glc.AllGameUser() {
			glc.notifyUserBetTimes(ur)
		}
		//发牌动画
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, glc.fapai)
	})
	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})
	return glc
}

func TestGameRound1(t *testing.T) {
	glc := NewQzniuniuLogicByTest()
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}
	glc.GameStartFunc()
	users := fixtures.GetGameUser()
	//玩家入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])

	gub := new(gameUserBet)
	gub.addUserBet(users[0].UserID, 0, true)
	glc.zhuangInfo = gub.userBets[0]

}
