package zjh

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestGameRound(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 10000
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewZhaJinHuaLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewZJHRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return NewZjhSuperToolService(room, tfm, dbMsgQueue)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*ZhaJinHuaLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.LogicBase)
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
		mlog.NewDefaultLogger().Error("zjh game end")
	default:
	}

}

//userBetAdd处理用户加注 負面列表測試
func TestIntelliUserBetAdd(t *testing.T) {
	//endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 100
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewZhaJinHuaLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewZJHRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*ZhaJinHuaLogic)

	users := fixtures.GetGameUser()

	enterUserAndGameStart(users, tgrs, glc)

	for _, testdata := range fixtures.IntellitestInt32() {
		glc.canOpt = true
		gmsg := new(msg.Message)
		gmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_BetAddID)
		gmsg.BClassID = int32(netproto.MessageBClassID_Game)
		ust := new(netproto.ZJH_BetAdd)
		ust.Betadd = proto.Int32(testdata)
		mData, err := proto.Marshal(ust)
		if err != nil {
			t.Errorf("parse proto error")
		}
		gmsg.MsgData = mData

		urd := glc.getUserBySeatNo(users[1].SeatNo)
		_, err = glc.userBetAdd(urd, gmsg.MsgData.([]byte))

		if err != nil {
			t.Logf("处理用户加注檢查正確! %v", err)
		}
	}
}

func enterUserAndGameStart(users []*gamelogic.GameUser, tgrs *fixtures.TestGameRoomServer, glc *ZhaJinHuaLogic) {
	for i := 0; i < 3; i++ {
		users[i].UserState.OnUserQueue(users[i], tgrs.Room)
		glc.UserEnter(users[i], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
		glc.UserReady(users[i])
	}
	glc.GameStartFunc()
}

func TestNickName(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 10000
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewZhaJinHuaLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewZJHRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*ZhaJinHuaLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.LogicBase)
	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	//for i:=0;i<100;i++ {
	nickName := glc.getNickName()
	//t.Log("nickName:", nickName)
	if nickName == "" {
		t.Errorf("nickName:%v", nickName)
	}
	//}

}

//指定發牌遊戲流程
func TestPokerGameRound(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 10000
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewZhaJinHuaLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewZJHRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return NewZjhSuperToolService(room, tfm, dbMsgQueue)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*ZhaJinHuaLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.LogicBase)
	glc.RegisterBeforeGameStartFunc(func() {
		//因為機器人要進行換牌所以，牌型必需要在 beforeGameStart 定義好
		glc.clearGameData()
		glc.defineCardType()
	})

	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	//if glc.GetCurrentGameUserCount() != 0 {
	//	t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	//}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := NewZjhUserData(users[0])
	glc.userData = append(glc.userData, userdata)
	userdata1 := NewZjhUserData(users[1])
	glc.userData = append(glc.userData, userdata1)

	//開始遊戲，丟底注
	glc.betDi()
	card0, card1 := getfilePoker()
	glc.setUserCards(userdata, card0)
	glc.setUserCards(userdata1, card1)

	glc.canOpt = true
	msg := getAddMsg(glc)
	glc.userBetAdd(userdata, msg.MsgData.([]byte))
	glc.userBetAdd(userdata1, msg.MsgData.([]byte))

	cmsg := getCompareMsg(userdata1.User.SeatNo)
	glc.compareCardsByplayer(userdata1, cmsg.MsgData.([]byte))
	glc.calWinMoney(userdata)
	if userdata.userMoney == 99990 && userdata1.userMoney == 100000 {
		//沒毛病

	} else {
		t.Errorf("玩家0贏錢錯誤:%v", userdata.userMoney)
		t.Errorf("玩家1贏錢錯誤:%v", userdata1.userMoney)
	}
}

func getfilePoker() (*poker.CardGroup, *poker.CardGroup) {
	var p poker.Poker
	p = poker.NewTestFilePoker("testsettlement")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	return fp.GetUserCards(int32(0)), fp.GetUserCards(int32(1))
}

//加注消息
func getAddMsg(glc *ZhaJinHuaLogic) *msg.Message {
	mlog.Debug("加注")
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_BetAddID)
	btdata := new(netproto.ZJH_BetAdd)
	btdata.Betadd = proto.Int32(glc.betAdd[0])
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}

//比牌消息
func getCompareMsg(seat int32) *msg.Message {
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.ZJH_GameMessageClassID_CompareID)

	btdata := new(netproto.ZJH_SeatNO)
	btdata.Seatno = proto.Int32(seat)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	return btmsg
}
