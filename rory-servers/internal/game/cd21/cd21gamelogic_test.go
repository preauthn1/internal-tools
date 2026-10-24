package cd21

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewCD21gameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewCD21RobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
}

//隨機發牌遊戲流程
func TestGameRound(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

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
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)
	userdata1 := newCD21userdata(users[1])
	glc.userdata = append(glc.userdata, userdata1)

	//t.Log("shouldbetotal => ", float64(glc.bettotal+700)*0.95)
	glc.userBet(users[0], fixtures.GetUserBet(BIGGER13, 100))
	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))

	glc.userBet(users[1], fixtures.GetUserBet(BIGGER13, 100))
	glc.userBet(users[1], fixtures.GetUserBet(DIZHU, 100))

	glc.betEnd()
	glc.fapai()
	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, cardPoint:%v", userdata.userCards.Cards, glc.getCardPoint(userdata.userCards))
	//t.Logf("user2 poker => %v, cardPoint:%v", userdata1.userCards.Cards, glc.getCardPoint(userdata1.userCards))
	glc.middleSettlement()
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//t.Logf("user2 poker => %v, win:%v", userdata1.userCards.Cards, userdata1.win)

	glc.finalSettlement()
	if glc.getCardPoint(glc.zhungCard) > glc.getCardPoint(userdata.userCards) {
		if userdata.win4 != 0 {
			t.Errorf("user point => %v, zhuang point => %v, 玩家不應該贏錢:%v", glc.getCardPoint(userdata.userCards), glc.getCardPoint(glc.zhungCard), userdata.win4)
		}
	} else if glc.getCardPoint(glc.zhungCard) < glc.getCardPoint(userdata.userCards) {
		if userdata.win4 == 0 && userdata.win1 == 0 {
			t.Errorf("user point => %v, zhuang point => %v, 玩家應該要贏錢:%v", glc.getCardPoint(userdata.userCards), glc.getCardPoint(glc.zhungCard), userdata.win4)
		}
	}
	//t.Logf("玩家1 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)
	//t.Logf("玩家2 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata1.user.UserID, userdata1.win, userdata1.win1, userdata1.win2, userdata1.win3, userdata1.win4, userdata1.win5, userdata1.tax)

}

//指定發牌遊戲流程
func TestFilePokerGameRound(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

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
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)
	userdata1 := newCD21userdata(users[1])
	glc.userdata = append(glc.userdata, userdata1)

	//t.Log("shouldbetotal => ", float64(glc.bettotal+700)*0.95)
	glc.userBet(users[0], fixtures.GetUserBet(BIGGER13, 100))
	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))

	glc.userBet(users[1], fixtures.GetUserBet(BIGGER13, 100))
	glc.userBet(users[1], fixtures.GetUserBet(DIZHU, 100))

	glc.betEnd()
	//發牌
	glc.zhungCard, userdata.userCards, userdata1.userCards = getfilePoker()

	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, cardPoint:%v", userdata.userCards.Cards, glc.getCardPoint(userdata.userCards))
	//t.Logf("user2 poker => %v, cardPoint:%v", userdata1.userCards.Cards, glc.getCardPoint(userdata1.userCards))
	glc.middleSettlement()
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//t.Logf("user2 poker => %v, win:%v", userdata1.userCards.Cards, userdata1.win)

	glc.finalSettlement()
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if glc.getCardPoint(glc.zhungCard) > glc.getCardPoint(userdata.userCards) {
		if userdata.win4 != 0 {
			t.Errorf("user point => %v, zhuang point => %v, 玩家不應該贏錢:%v", glc.getCardPoint(userdata.userCards), glc.getCardPoint(glc.zhungCard), userdata.win4)
		}
	} else if glc.getCardPoint(glc.zhungCard) < glc.getCardPoint(userdata.userCards) {
		if userdata.win4 == 0 && userdata.win1 == 0 {
			t.Errorf("user point => %v, zhuang point => %v, 玩家應該要贏錢:%v", glc.getCardPoint(userdata.userCards), glc.getCardPoint(glc.zhungCard), userdata.win4)
		}
	}
}

//輸入需要返回的牌副數
func getfilePoker() (*poker.CardGroup, *poker.CardGroup, *poker.CardGroup) {
	var p poker.Poker
	p = poker.NewTestFilePoker("filepoker")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	return fp.GetUserCards(int32(0)), fp.GetUserCards(int32(1)), fp.GetUserCards(int32(2))
}

//單玩家隨機發牌分牌流程
func TestSplitSettlements(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	//glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)
	//userdata1 := newCD21userdata(users[1])
	//glc.userdata = append(glc.userdata, userdata1)

	//t.Log("shouldbetotal => ", float64(glc.bettotal+700)*0.95)
	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))

	//glc.userBet(users[1], fixtures.GetUserBet(BIGGER13, 100))
	//glc.userBet(users[1], fixtures.GetUserBet(DIZHU, 100))

	glc.betEnd()
	userdata.isSplit = true
	faSplitPai(glc, userdata)

	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user poker => %v, cardPoint:%v", userdata.userCards.Cards, glc.getCardPoint(userdata.userCards))
	//t.Logf("split poker => %v, cardPoint:%v", userdata.splitCards.Cards, glc.getCardPoint(userdata.splitCards))
	glc.middleSettlement()
	//glc.dealZhuangPoker()
	glc.finalSettlement()
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)
	//t.Logf("玩家2 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata1.user.UserID, userdata1.win, userdata1.win1, userdata1.win2, userdata1.win3, userdata1.win4, userdata1.win5, userdata1.tax)

}

func faSplitPai(glc *CD21gameLogic, userdata *cd21userdata) {
	zhungCard := glc.poker.Deal(0, 2)
	for _, dealcard := range zhungCard.Cards {
		dCard := new(poker.Card)
		dCard.CardNumber = dealcard.CardNumber
		dCard.Color = dealcard.Color
		glc.zhungCard.Cards = append(glc.zhungCard.Cards, dCard)
	}

	user1Card := glc.poker.Deal(userdata.user.UserID, 2)
	for _, usercard := range user1Card.Cards {
		uCard := new(poker.Card)
		uCard.CardNumber = usercard.CardNumber
		uCard.Color = usercard.Color
		userdata.userCards.Cards = append(userdata.userCards.Cards, uCard)
	}

	user2Card := glc.poker.Deal(userdata.user.UserID, 2)
	for _, usercard1 := range user2Card.Cards {
		uCard1 := new(poker.Card)
		uCard1.CardNumber = usercard1.CardNumber
		uCard1.Color = usercard1.Color
		userdata.splitCards.Cards = append(userdata.splitCards.Cards, uCard1)
	}
}

//保險結算
func TestBaoXianGameRound(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))

	glc.betEnd()

	var p poker.Poker
	p = poker.NewTestFilePoker("baoxian")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.insurance(userdata)
	glc.fanZhuangPoker()
	glc.finalSettlement()
	if userdata.win2 != 150 {
		t.Errorf("保險贏錢錯誤 %d", userdata.win2)
	}
}

//對子結算
func TestDuiZiSettlement(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))
	glc.userBet(users[0], fixtures.GetUserBet(DUIZI, 100))

	glc.betEnd()

	var p poker.Poker
	p = poker.NewTestFilePoker("duizi")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.fanZhuangPoker()
	glc.finalSettlement()

	t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if userdata.win != 1145 {
		t.Errorf("對子贏錢錯誤 %d", userdata.win)
	}
}

//莊爆結算
func TestZhuangBaoSettlement(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))
	glc.userBet(users[0], fixtures.GetUserBet(ZHUANGBAO, 100))

	glc.betEnd()

	var p poker.Poker
	p = poker.NewTestFilePoker("zhuangbao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.fanZhuangPoker()
	glc.dealZhuangPoker()
	glc.finalSettlement()
	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if userdata.win3 != 1525 {
		t.Errorf("莊爆贏錢錯誤 %d", userdata.win3)
	}
}

//Lucky20結算
func TestLucky20Settlement(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))
	glc.userBet(users[0], fixtures.GetUserBet(LUCKY20, 100))

	glc.betEnd()

	var p poker.Poker
	p = poker.NewTestFilePoker("lucky20")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.fanZhuangPoker()
	glc.finalSettlement()
	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if userdata.win != 575 {
		t.Errorf("幸運20贏錢錯誤 %d", userdata.win)
	}
}

//大小13結算
func Test13Settlement(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))
	glc.userBet(users[0], fixtures.GetUserBet(BIGGER13, 100))

	glc.betEnd()

	var p poker.Poker
	p = poker.NewTestFilePoker("bigger13")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.fanZhuangPoker()
	glc.finalSettlement()
	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if userdata.win != 195 {
		t.Errorf("幸運20贏錢錯誤 %d", userdata.win)
	}
}

//加倍測試
func TestDouble(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()
	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	glc.userBet(users[0], fixtures.GetUserBet(DIZHU, 100))

	glc.betEnd()
	var p poker.Poker
	p = poker.NewTestFilePoker("bigger13")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	glc.zhungCard = fp.GetUserCards(int32(0))
	userdata.userCards = fp.GetUserCards(int32(1))

	glc.middleSettlement()
	glc.doubleCard(userdata)
	glc.fanZhuangPoker()
	glc.finalSettlement()
	//t.Logf("zhuang poker => %v, cardPoint:%v", glc.zhungCard.Cards, glc.getCardPoint(glc.zhungCard))
	//t.Logf("user1 poker => %v, win:%v", userdata.userCards.Cards, userdata.win)
	//
	//t.Logf("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 稅 %d", userdata.user.UserID, userdata.win, userdata.win1, userdata.win2, userdata.win3, userdata.win4, userdata.win5, userdata.tax)

	if userdata.win4 != 390 {
		t.Errorf("用户double贏錢錯誤 %d", userdata.win4)
	}
}

//13測試
func TestThirteen(t *testing.T) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()

	var bigger, smaller, middle, bigger2, smaller2, middle2 int32
	for i := 0; i < 156; i++ {
		poker := glc.poker.Deal(0, 2)
		point1 := getPoint(poker)
		point2 := glc.getCardPoint(poker)
		if point1 > 13 {
			bigger++
		} else if point1 < 13 {
			smaller++
		} else {
			middle++
		}

		if point2 > 13 {
			bigger2++
		} else if point2 < 13 {
			smaller2++
		} else {
			middle2++
		}
	}
	if smaller > bigger || smaller2 > bigger2 {
		t.Errorf("這不合理阿")
		t.Logf("更新版大於13 %d, 小於13 %d, 等於13 %d, 總場次 %d ", bigger, smaller, middle, 156)
		t.Logf("原版大於13 %d, 小於13 %d, 等於13 %d, 總場次 %d", bigger2, smaller2, middle2, 156)
	}
}

//計算牌點數
func getPoint(card *poker.CardGroup) int32 {
	number := int32(0)
	//count := int32(0)
	for _, cardPoint := range card.Cards {
		if cardPoint.CardNumber == 14 {
			if card.Cards[0].CardNumber == 14 && card.Cards[1].CardNumber >= 10 && card.Cards[1].CardNumber < 14 {
				number += 11
			} else if card.Cards[1].CardNumber == 14 && card.Cards[0].CardNumber >= 10 && card.Cards[0].CardNumber < 14 {
				number += 11
			} else {
				number += 1
			}
		} else if cardPoint.CardNumber > 10 {
			number += 10
		} else {
			number += cardPoint.CardNumber
		}
	}

	return number
}

//syspool test +100
func TestSyspoolTest(t *testing.T) {
	conf.TestMode = true
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*CD21gameLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.DynamicGameLogicBase.LogicBase)

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
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	userdata := newCD21userdata(users[0])
	glc.userdata = append(glc.userdata, userdata)

	serverid := int32(999)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	syspool := new(netproto.Cd21SysPoolData)
	syspool.FlagID = proto.Int32(2)
	syspool.PoolAmount = proto.Int64(0)
	syspool.SumBetCount = proto.Int64(0)
	syspool.SumBetMoney = proto.Int64(0)

	glc.setSyspoolData(syspool)

	//mlog.System("PoolAmount: ", glc.syspoolData.PoolAmount)
	//mlog.System("SumBetCount: ", glc.syspoolData.SumBetCount)
	//mlog.System("SumBetMoney: ", glc.syspoolData.SumBetMoney)
	//mlog.System("sys: ", glc.syspoolData)
}

//看有沒有什麼需要補足的測試
