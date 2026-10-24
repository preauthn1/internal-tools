package sg

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//牌型驗證
func TestGetSanGongPokerType(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("getSanGongType")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		actual := getSanGongType(cards)

		expect := int32(5)
		if i == 0 && actual != expect {
			t.Errorf("index:%v 單牌錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = 9
		if i == 1 && actual != expect {
			t.Errorf("index:%v 特殊牌9錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = 8
		if i == 2 && actual != expect {
			t.Errorf("index:%v 特殊牌8錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = HunSanGong
		if i == 3 && actual != expect {
			t.Errorf("index:%v 混三公錯誤 actual:%v expect:%v", i, actual, expect)
		}

		if i == 4 && actual != expect {
			t.Errorf("index:%v 混三公錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = XiaoSanGong
		if i == 5 && actual != expect {
			t.Errorf("index:%v 小三公錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = XiaoSanGong
		if i == 6 && actual != expect {
			t.Errorf("index:%v 小三公錯誤 actual:%v expect:%v", i, actual, expect)
		}

		expect = DaSanGong
		if i == 7 && actual != expect {
			t.Errorf("index:%v 大s三公錯誤 actual:%v expect:%v", i, actual, expect)
		}

	}
}

//比较牌型
func TestCompareCards(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("comPareCards")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	cards1 := fp.GetUserCards(int32(0))
	cards2 := fp.GetUserCards(int32(1))
	actual := comPareCards(cards1, cards2)
	expect := int32(2)
	if actual != expect {
		t.Errorf("牌型比較 card1 < cards2 錯誤 acutal : %v", actual)
	}

	cards3 := fp.GetUserCards(int32(3))
	cards4 := fp.GetUserCards(int32(4))
	actual2 := comPareCards(cards3, cards4)
	expect2 := int32(2)
	if actual2 != expect2 {
		t.Errorf("牌型比較 card3 < cards4 錯誤 acutal : %v", actual2)
	}

	cards5 := fp.GetUserCards(int32(5))
	cards6 := fp.GetUserCards(int32(6))
	actual3 := comPareCards(cards5, cards6)
	expect3 := int32(1)
	if actual3 != expect3 {
		t.Errorf("牌型比較 card5 > cards6 錯誤 acutal : %v", actual3)
	}

	cards7 := fp.GetUserCards(int32(7))
	cards8 := fp.GetUserCards(int32(8))
	actual4 := comPareCards(cards7, cards8)
	expect4 := int32(1)
	if actual4 != expect4 {
		t.Errorf("牌型比較 card7 > cards8 錯誤 acutal : %v", actual4)
	}

	cards9 := fp.GetUserCards(int32(9))
	cards10 := fp.GetUserCards(int32(10))
	actual5 := comPareCards(cards9, cards10)
	expect5 := int32(2)
	if actual5 != expect5 {
		t.Errorf("牌型比較 card9 > cards10 錯誤 acutal : %v", actual5)
	}

	cards11 := fp.GetUserCards(int32(11))
	actual6 := comPareCards(cards11, cards10)
	expect6 := int32(2)
	if actual6 != expect6 {
		t.Errorf("牌型比較 card11 > cards10 錯誤 acutal : %v", actual6)
	}

	cards12 := fp.GetUserCards(int32(12))
	actual7 := comPareCards(cards10, cards12)
	expect7 := int32(1)
	if actual7 != expect7 {
		t.Errorf("牌型比較 card10 > cards12 錯誤 acutal : %v", actual7)
	}

	actual7 = comPareCards(cards12, cards12)
	expect7 = int32(1)
	if actual7 != expect7 {
		t.Errorf("牌型比較 card10 == cards12 錯誤 acutal : %v", actual7)
	}

	cards13 := fp.GetUserCards(int32(13))
	actual7 = comPareCards(cards13, cards13)
	expect7 = int32(1)
	if actual7 != expect7 {
		t.Errorf("牌型比較 card10 == cards12 錯誤 acutal : %v", actual7)
	}
}

func TestGameRound1(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 10000
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewSGGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewSGRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*SGGameLogic)

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

	users := fixtures.GetGameUser()

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)
	robots[2].UserState.OnUserQueue(robots[2], tgrs.Room)

	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(robots[1])
	glc.UserReady(robots[2])

	//金額驗證未完成
	glc.qiangZhuangResult()
	glc.betResult()

	select {
	case _ = <-endChannel:
		mlog.NewDefaultLogger().Error("sg game end")
	}

}
