package hhdz

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

type TestGameRoomServer struct {
	*logicmodule.GameServerLogicModule
	GameServer           *gameserver.GameServer
	GameLogicManager     *gamelogic.LogicManager
	Room                 *gamelogic.Room
	RoomActionController *roomlogic.RoomActionController
	RoomActionManager    *roomlogic.RoomActionManager
	RobotService         *robot.Service
	sstorage             scorestorage.UserScoreStorage
}

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHDZLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
}

//驗證指定紅贏開牌是否有重複牌型
func TestGetPokerByType(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHDZLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)

	//for i := 0; i < 1000; i++ {
	glc.clearGameData()
	glc.fapai()
	glc.getResultType()
	glc.getPokerByType(glc.jiangtype)

	for _, vhung := range glc.pokerHelper.hungStore.cards {
		for _, vhei := range glc.pokerHelper.heiStore.cards {
			if vhung.Color == vhei.Color && vhung.CardNumber == vhei.CardNumber {
				t.Errorf("指定紅贏, 重複牌型 紅color[%v]number[%v] 黑color[%v]number[%v]", vhung.Color, vhung.CardNumber, vhei.Color, vhei.CardNumber)
			}
		}
	}
	//}
}

//驗證指定順子開牌是否正確
func TestGetPokerByShunZi(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHDZLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)

	//for i := 0; i < 10000; i++ {

	glc.clearGameData()
	glc.fapai()
	glc.getResultType()
	glc.getPokerByType(glc.jiangtype)
	//glc.pokerHelper.getPokerCardByWinUserID(_HungUserID)
	if len(glc.pokerHelper.hungStore.cards) != 3 || len(glc.pokerHelper.heiStore.cards) != 3 {
		t.Errorf("空牌型 紅card[%v] 黑card[%v]", glc.pokerHelper.hungStore.cards, glc.pokerHelper.heiStore.cards)
	}
	//}
}
func TestOneTimeIntegration(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	glc.RegisterBeforeGameEndFunc(func() {

	})

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	select {
	case _ = <-endChannel:
		t.Log("hhdz game end")
	}
}

//驗證開獎
func TestKaiJiang(t *testing.T) {
	var kaiHung int = 0
	var kaiHei int = 0
	var nature int = 0

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)

	for i := 0; i < 1000; i++ {
		glc.clearGameData()
		glc.fapai()
		glc.getResultType()
		glc.getPokerByType(glc.jiangtype)
		if glc.jiangtype != 0 {
			if glc.jiangtype == 1 {
				kaiHung++
			} else if glc.jiangtype == 2 {
				kaiHei++
			} else {
				nature++
			}
		}
	}
	t.Log("紅贏:", kaiHung, "黑贏:", kaiHei, "自然:", nature)
}

//流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *HHDZLogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	glc.robotzhuang = 0
	glc.OnGameStart(func() {
	})

	_, robots := fixtures.GetGameUsers()
	robotManager := tgrs.RobotService.GetRobotManager()
	if robotManager.RobotCount() != 0 {
		t.Errorf("初始化, 機器人數量應為0, RobotCount=>%d", robotManager.RobotCount())
	}
	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}
	return tgrs, glc, endChannel, robots
}

func TestNatureKaiJiang(t *testing.T) {
	var kaiHung int = 0
	var kaiHei int = 0
	var hungzuhetype1 int = 0
	var heizuhetype1 int = 0
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHDZLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	mlog.System("pokertype %d", glc.GetServerConfig().PokerType)
	//rand.Seed(time.Now().UnixNano())
	for i := 0; i < 1000; i++ {
		glc.clearGameData()
		glc.fapai()
		glc.getResultType()
		glc.getPokerByType(glc.jiangtype)
		//glc.jiangtype = int32(glc.pokerHelper.getJiangType())
		//_, _, z := glc.pokerHelper.getZuheTypes()
		//glc.zuhetype = int32(z)
		if glc.jiangtype != 0 {
			if glc.jiangtype == 1 {
				if glc.zuhetype <= 8 && glc.zuhetype >= 4 {
					hungzuhetype1++
				}
				kaiHung++
			} else if glc.jiangtype == 2 {
				if glc.zuhetype <= 8 && glc.zuhetype >= 4 {
					heizuhetype1++
				}
				kaiHei++
			}
		}
	}
	//t.Log("紅贏:", kaiHung, "黑贏:", kaiHei, "紅特殊:", hungzuhetype1, "黑特殊:", heizuhetype1)
}

//下注紅開紅, 遊戲流程入桌→下注→開獎→結果
func TestGameRound(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 5000
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	glc.SysPool = gamelogic.NewSysPoolHelper(glc.LogicBase, glc.Configs)
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
	users[0].UserType = 1
	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.fapai()
	//下注
	glc.userBet(users[0], fixtures.GetUserBet(HUNG, 100))
	glc.nextjiangtype = int32(HUNG)

	glc.kaiJiang()
	//lhd.sendResult()
	glc.calWinLose()

	if glc.SysPool.GetPool() != int64(1499900) {
		t.Errorf("syspool:%v", glc.SysPool.GetPool())
	}

	if glc.systotalwin != int64(0) {
		t.Errorf("systotalwin:%v", glc.systotalwin)
	}

	if glc.syslastwin != int64(-100) {
		t.Errorf("syslastwin:%v", glc.syslastwin)
	}

	if glc.systaxincome != int64(5) {
		t.Errorf("systaxincome:%v", glc.systaxincome)
	}

	if glc.userdata[0].updatemoney != int64(95) {
		t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
	}

	if glc.userdata[0].win != int64(195) {
		t.Errorf("win:%v", glc.userdata[0].win)
	}

	if glc.zhuangWin != int64(-100) {
		t.Errorf("zhuangWin:%v", glc.zhuangWin)
	}

	glc.gameEnd()
	glc.clearGameData()
}
