package srnn

import (
	"os"
	"path"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestGameRound(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 20

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewSRNNGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewSRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*SRNNGameLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})

	_, robots := fixtures.GetGameUsers()
	robotManager := tgrs.RobotService.GetRobotManager()

	if robotManager.RobotCount() != 0 {
		t.Errorf("初始化機器人數量應該為0, RobotCount => %d", robotManager.RobotCount())
	}
	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robotManager.AddRobot(robot)
	}
	if robotManager.RobotCount() != 3 {
		t.Errorf("加載後機器人數量應為3, RobotCount=>%d", robotManager.RobotCount())
	}
	//非機器人無法從free狀態入座，先統一狀態為queue
	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)

	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(robots[1])

}

func TestStartGameLogicServer(t *testing.T) {
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewSRNNGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewSRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
	fixtures.StartGameLogicByTest()
}

func TestRound(t *testing.T) {
	glc := NewSrNiuNiuLogicByTest()
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}
	glc.GameStartFunc()
	users, robot := fixtures.GetGameUsers()

	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robot[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(robot[0])
}

func NewSrNiuNiuLogicByTest() *SRNNGameLogic {
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

	glc := new(SRNNGameLogic)

	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, cfgs)
	parseNNConfig(cfg)
	glc.RegisterAfterGameEndFunc(func() {
		for _, ur := range glc.AllGameUser() {
			glc.notifyUserBetTimes(ur)
		}
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_Card, glc.fapai)
	})
	glc.RegisterBeforeGameStartFunc(func() {
		glc.clearGameData()
	})

	return glc
}
