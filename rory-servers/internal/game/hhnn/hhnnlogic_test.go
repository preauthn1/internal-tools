package hhnn

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
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
	"example.invalid/ServerTeam/Servers/pkg/util"
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
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
}

//
////廣播第一手三張牌
//func TestGetFirstCardGroup(t *testing.T) {
//	conf.CLogLevel = mlog.ErrorLevel
//	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewHHNNLogic(t, tfManager, cfgs)
//	}
//
//	tgrs := fixtures.NewTestGameRoomServer()
//	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)
//
//	glc.poker = poker.NewRandomPoker(true, true, 1)
//
//	glc.clearGameData()
//	glc.fapeiFirstCardGroup()
//	firstCardGroup := glc.getFirstCardGroup()
//
//	firstCardData := firstCardGroup.MsgData.(*netproto.HHNN_BrodcastPokerCard)
//
//	if len(firstCardData.UserPokerCard) != int(2) {
//		t.Error("發牌組數錯誤")
//	}
//
//	for _, v := range firstCardData.UserPokerCard {
//		if len(v.CardGroup) != int(3) {
//			t.Error("發牌數量不等於3")
//		}
//	}
//
//}
//
////廣播第二手兩張牌, 自然機率
//func TestGetSecondCardGroup(t *testing.T) {
//	conf.CLogLevel = mlog.DebugLevel
//	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewHHNNLogic(t, tfManager, cfgs)
//	}
//	//roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
//	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
//	//}
//
//	tgrs := fixtures.NewTestGameRoomServer()
//	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)
//
//	glc.poker = poker.NewRandomPoker(false, false, 1)
//
//	glc.winline = -1
//
//	glc.clearGameData()
//	glc.fapeiFirstCardGroup()
//	glc.getFirstCardGroup()
//	glc.getResultType()
//	glc.getPokerByType(NATURE)
//
//	glc.getGameResultPokerMsg()
//
//	if len(glc.cardhung.Cards) != 5 {
//		t.Errorf("紅區發牌未滿五張")
//	} else {
//		t.Logf("紅區發牌len[%v]", len(glc.cardhung.Cards))
//	}
//
//	if len(glc.cardhei.Cards) != 5 {
//		t.Errorf("黑區發牌未滿五張")
//	} else {
//		t.Logf("黑區發牌len[%v]", len(glc.cardhei.Cards))
//	}
//
//}

//流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *HHNNLogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	glc.robotzhuang = 0
	glc.syspool = 0
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

//系統當莊, 真人(輸)下注紅
func TestGameRoundBetAndKaijaing1(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card1")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(570) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(570) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(30) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-600) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-500) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(570) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 真人(贏)下注黑
func TestGameRoundBetAndKaijaing2(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card1")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-600) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-600) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(30) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(570) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(670) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-600) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 真人(贏)下注好運
func TestGameRoundBetAndKaijaing3(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card1")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(95) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(95) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(5) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-100) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(0) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(95) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開紅 五小牛 真人(贏)下注紅
func TestGameRoundBetAndKaijaing4(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(950) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1050) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開紅 五小牛 真人(輸)下注黑
func TestGameRoundBetAndKaijaing5(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(950) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(950) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-1000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-900) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(950) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開紅 五小牛 真人(輸)下注好運
func TestGameRoundBetAndKaijaing6(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-2000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-2000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(100) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(1900) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(2000) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-2000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 四詐 真人(輸)下注紅
func TestGameRoundBetAndKaijaing7(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(950) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(950) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-1000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-900) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(950) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 四詐 真人(贏)下注黑
func TestGameRoundBetAndKaijaing8(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(950) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1050) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 四詐 真人(贏)下注好運
func TestGameRoundBetAndKaijaing9(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1500) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1500) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(75) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(1425) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1525) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1500) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 五花牛 真人(輸)下注紅
func TestGameRoundBetAndKaijaing10(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(950) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(950) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-1000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-900) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(950) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 五花牛 真人(贏)下注黑
func TestGameRoundBetAndKaijaing11(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(950) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1050) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 五花牛 真人(贏)下注好運
func TestGameRoundBetAndKaijaing12(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(950) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1050) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 牛牛K 真人(輸)下注紅
func TestGameRoundBetAndKaijaing13(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card5")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(950) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(950) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-1000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-900) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(950) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 牛牛K 真人(贏)下注黑
func TestGameRoundBetAndKaijaing14(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card5")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-1000) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-1000) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(950) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(1050) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-1000) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 開黑 牛牛K 真人(贏)下注幸運
func TestGameRoundBetAndKaijaing15(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card5")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-200) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-200) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(10) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(190) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(290) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(-200) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

//系統當莊, 紅黑都散牌
func TestGameRoundBetAndKaijaing17(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card7")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-100) {
					t.Errorf("syspool:%v", glc.syspool)
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
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

func TestMinJiangType(t *testing.T) {
	conf.CLogLevel = mlog.DebugLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	//roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	//}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	glc.poker = poker.NewRandomPoker(false, false, 1)

	glc.winline = -1

	glc.clearGameData()
	glc.getFirstResultType()
	//glc.fapeiFirstCardGroup()
	//glc.getFirstCardGroup()
	glc.getResultType()
	glc.getPokerByType(HUNG)
	//hungNiuRet := getNiuNiuRet(glc.cardhung)
	//heiNiuRet := getNiuNiuRet(glc.cardhei)
	//betret := compareCards(hungNiuRet, heiNiuRet)
	//if betret < int32(0) {
	//	t.Errorf("開獎類型:",glc.jiangtype)
	//	t.Errorf("error:", glc.cardhung, glc.cardhei)
	//}

	glc.getGameResultPokerMsg()

	if len(glc.cardhung.Cards) != 5 {
		t.Errorf("紅區發牌未滿五張")
	} else {
		t.Logf("紅區發牌len[%v]", len(glc.cardhung.Cards))
	}

	if len(glc.cardhei.Cards) != 5 {
		t.Errorf("黑區發牌未滿五張")
	} else {
		t.Logf("黑區發牌len[%v]", len(glc.cardhei.Cards))
	}

}

//A炸彈 4炸彈 下紅(輸)
func TestBomb(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 600
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])

	var p poker.Poker
	p = poker.NewTestFilePoker("card6")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
			ust := m.MsgData.(*netproto.HHNN_Sate)

			if *ust.State == Bet && users[0].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0

				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
			}

			if *ust.State == KaiJiangAction {
				cards := fp.GetUserCards(int32(0))
				cards1 := fp.GetUserCards(int32(1))
				//竄改想要開獎的牌驗證資料
				glc.cardhung.Cards = cards.Cards
				glc.cardhei.Cards = cards1.Cards
			}

			if *ust.State == Wait {
				if glc.syspool != int64(950) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(950) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(50) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-1000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-900) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(950) {
					t.Errorf("zhuangWin:%v", glc.zhuangWin)
				}
			}
		}

	}

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterBeforeGameEndFunc(func() {

	})

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Info("brnn game end")
	}
}

func TestNiuNiu(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	for i := 0; i < 10000; i++ {
		glc.poker = poker.NewRandomPoker(false, false, 1)
		glc.poker.Shuffle()
		glc.clearGameData()
		takeCardGroup := getNiuNiu(glc, false)
		//t.Log(takeCardGroup)
		for _, u := range takeCardGroup.Cards {
			if u.CardNumber == 13 {
				t.Errorf("不能有K:%v", takeCardGroup)
			}
		}
	}
}

func TestWuNiu(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	for i := 0; i < 10000; i++ {
		glc.poker = poker.NewRandomPoker(false, false, 1)
		glc.poker.Shuffle()
		glc.clearGameData()
		takeCardGroup := getWuNiu(glc)
		ret := getNiuNiuRet(takeCardGroup)
		//t.Log(takeCardGroup,ret.number)
		for _, u := range takeCardGroup.Cards {
			if u.CardNumber == 0 || ret.number > 0 {
				t.Errorf("不能有牛:%v 牛%v", takeCardGroup, ret.number)
			}
		}
	}
}

func TestHasNiu(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	//for i := 0; i < 10000; i++ {
	glc.poker = poker.NewRandomPoker(false, false, 1)
	glc.poker.Shuffle()
	glc.clearGameData()
	number := util.RndInt(9, 18)
	takeCardGroup := getHasNiu(glc, 18-number)
	ret := getNiuNiuRet(takeCardGroup)
	//t.Log(takeCardGroup,ret.number, 18-number)
	for _, u := range takeCardGroup.Cards {
		if u.CardNumber == 0 || ret.number != 18-number {
			t.Errorf("牛給錯了:%v 牛%v ，應該要是牛%v", takeCardGroup, ret.number, 18-number)
		}
	}
	//}
}

func TestLoseCard(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	//for i := 0; i < 1000; i++ {
	glc.poker = poker.NewRandomPoker(false, false, 1)
	glc.poker.Shuffle()
	glc.clearGameData()
	number := util.RndInt(9, 18)
	takeCardGroup := getHasNiu(glc, 18-number)
	ret := getNiuNiuRet(takeCardGroup)
	//t.Log(takeCardGroup,ret.number, 18-number)
	for _, u := range takeCardGroup.Cards {
		if u.CardNumber == 0 || ret.number != 18-number {
			t.Errorf("牛給錯了:%v 牛%v ，應該要是牛%v", takeCardGroup, ret.number, 18-number)
		}
	}
	//LoseCard := getLoseCard(glc, takeCardGroup)
	//ret2 := getNiuNiuRet(LoseCard.cards)
	//if ret2.number > ret.number {
	//	t.Errorf("ret2應該要輸，贏家:%v 輸家:%v，ret1:%v，ret2:%v", takeCardGroup, LoseCard.cards, ret.number, ret2.number)
	//}
	//t.Log( takeCardGroup, LoseCard.cards, ret.number,ret2.number)
	//}
}

//func TestAutoFilePoker(t *testing.T) {
//	endChannel := make(chan bool)
//	conf.CLogLevel = mlog.DebugLevel
//	conf.TimerFastMode = 600
//	conf.TestMode = true
//	roomserver.GameLogicCreator = func(t *gamelogic.tableManager, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewHHNNLogic(t, tfManager, cfgs)
//	}
//	roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
//		return NewHHNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
//	}
//
//	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)
//
//	robotManager := tgrs.RobotService.GetRobotManager()
//
//	for i := 0; i < 3; i++ {
//		robot := new(robotcomm.Robot)
//		robot.UserInfo = robots[i]
//		robot.RebotLevel = 6
//		robot.UserInfo.UserType = 2
//		robotManager.AddRobot(robot)
//	}
//
//	users := fixtures.GetGameUser()
//
//	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
//	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
//
//	//安排用戶入座
//	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
//
//	glc.UserReady(robots[0])
//	glc.UserReady(users[0])
//
//	var p poker.Poker
//	p = poker.NewAutoTestFilePoker("autopoker")
//	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
//	fp := p.(*poker.FixedPoker)
//
//	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
//
//		if m.SClassID == int32(netproto.HHNN_GameMessageClassID_HHNNStateID) {
//			ust := m.MsgData.(*netproto.HHNN_Sate)
//
//			if *ust.State == Bet && users[0].UserID == userid { //下注
//				glc.syspool = 0
//				glc.syslastwin = 0
//				glc.systaxincome = 0
//				glc.zhuangWin = 0
//
//				glc.userBet(users[1], fixtures.GetUserBet(1, 100))
//			}
//
//			if *ust.State == KaiJiangAction {
//				cards := fp.GetUserCards(int32(0))
//				cards1 := fp.GetUserCards(int32(1))
//				//竄改想要開獎的牌驗證資料
//				glc.cardhung.Cards = cards.Cards
//				glc.cardhei.Cards = cards1.Cards
//			}
//
//			if *ust.State == Wait {
//				if glc.syspool != int64(950) {
//					t.Errorf("syspool:%v", glc.syspool)
//				}
//
//				if glc.syslastwin != int64(950) {
//					t.Errorf("syslastwin:%v", glc.syslastwin)
//				}
//
//				if glc.systaxincome != int64(50) {
//					t.Errorf("systaxincome:%v", glc.systaxincome)
//				}
//
//				if glc.userdata[0].updatemoney != int64(-1000) {
//					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
//				}
//
//				if glc.userdata[0].win != int64(-900) {
//					t.Errorf("win:%v", glc.userdata[0].win)
//				}
//
//				if glc.zhuangWin != int64(950) {
//					t.Errorf("zhuangWin:%v", glc.zhuangWin)
//				}
//			}
//		}
//
//	}
//
//	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
//	}
//
//	glc.RegisterBeforeGameEndFunc(func() {
//
//	})
//
//	select {
//	case _ = <-endChannel:
//		mlog.NewDefaultLogger().Info("brnn game end")
//	}
//}

func TestNatureKaiJiang(t *testing.T) {
	var kaiHung int = 0
	var kaiHei int = 0
	var nature int = 0
	var kaiHung1 int = 0
	var kaiHei1 int = 0
	var hungzuhetype int = 0
	var heizuhetype int = 0
	var hungzuhetype1 int = 0
	var heizuhetype1 int = 0
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewHHNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*HHNNLogic)

	//for i := 0; i < 10000; i++ {
	glc.cardhung = new(poker.CardGroup)
	glc.cardhei = new(poker.CardGroup)
	//glc.firstPoker = make(map[int]*poker.CardGroup)
	//glc.secondPokerGroup = make(map[int]*poker.CardGroup)
	glc.poker = poker.NewRandomPoker(false, false, 1) //初始化一副牌
	glc.poker.Shuffle()

	glc.getFirstResultType()
	if glc.jiangtype == 1 {
		if glc.zuhetype <= KNIUNIU && glc.zuhetype >= WUXIAONIU {
			hungzuhetype++
		}
		kaiHung1++
	} else if glc.jiangtype == 2 {
		if glc.zuhetype <= KNIUNIU && glc.zuhetype >= WUXIAONIU {
			heizuhetype++
		}
		kaiHei1++
	} else {
		nature++
	}
	glc.fapai()
	//glc.startFirstCard()

	glc.getPokerByType(glc.jiangtype)
	if glc.jiangtype != 0 {
		if glc.jiangtype == 1 {
			if glc.zuhetype <= KNIUNIU && glc.zuhetype >= WUXIAONIU {
				hungzuhetype1++
			}
			kaiHung++
		} else if glc.jiangtype == 2 {
			if glc.zuhetype <= KNIUNIU && glc.zuhetype >= WUXIAONIU {
				heizuhetype1++
			}
			kaiHei++
		}
	}
	//}
	//t.Log("組牌前紅贏:",kaiHung1,"黑贏:",kaiHei1,"自然:",nature,"紅特殊:",hungzuhetype ,"黑特殊:",heizuhetype)
	//t.Log("紅贏:",kaiHung,"黑贏:",kaiHei,"紅特殊:",hungzuhetype1 ,"黑特殊:",heizuhetype1)
}
