package brnn

import (
	"fmt"
	"path"
	"runtime"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

}

func TestRandomFapai(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	//roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	//}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)
	stats := make(map[int]int)
	total := 0
	for i := 0; i < 2000; i++ {
		glc.clearGameData()
		glc.poker = poker.NewRandomPoker(false, false, 1)
		glc.poker.Shuffle()

		glc.fapeiFirstCardGroupNew()
		glc.natureResultType()
		for i := int32(0); i <= int32(len(glc.betvalue.alltype)); i++ {
			cg := glc.getCalSecondPokerGroup(i)
			glc.allKindPoker[int(i)].Cards = append(glc.allKindPoker[int(i)].Cards, cg.Cards...)

			result := getNiuNiuRet(glc.allKindPoker[int(i)])
			if _, ok := stats[result.niuRet]; ok {
				stats[result.niuRet]++
			} else {
				stats[result.niuRet] = 1
			}
			//t.Logf("%v", result)
			total++
		}
	}
	showStats := func(stats map[int]int, total int) {
		s := ""
		for k, v := range stats {
			switch k {
			case nnCard_WuNiu:
				s += fmt.Sprintf("無牛 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			case nnCard_Niu:
				s += fmt.Sprintf("有牛 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			case nnCard_NiuNiu:
				s += fmt.Sprintf("牛牛 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			case nnCard_WuHuaNiu:
				s += fmt.Sprintf("無花牛 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			case nnCard_WuXiaoNiu:
				s += fmt.Sprintf("無小牛 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			case nnCard_4Zha:
				s += fmt.Sprintf("四炸 (%d/%d) %.4f%% ", v, total, (float32(v)/float32(total))*100)
				break
			}
		}
		t.Log(s)
	}
	showStats(stats, total)

	stats = make(map[int]int)
	total = 0
	for i := 0; i < 10000; i++ {
		glc.clearGameData()
		glc.poker = poker.NewRandomPoker(false, false, 1)
		glc.poker.Shuffle()

		cg := glc.poker.Deal(0, 5)
		result := getNiuNiuRet(cg)
		if _, ok := stats[result.niuRet]; ok {
			stats[result.niuRet]++
		} else {
			stats[result.niuRet] = 1
		}
		//t.Logf("%v", result)
		total++
	}
	showStats(stats, total)
	//firstCardGroup := glc.getFirstCardGroup()
}

func TestRandomSlice(t *testing.T) {
	arr := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	sarr := util.IntArrayShuffle(arr)
	t.Logf("%v, %v", arr, sarr)

}

//指定通殺_
func TestSetNextJiangTypeTongsha(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel

	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 5000
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs, glc, endChannel, robots := IntegrationTestInit(endChannel, t)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.LogicBase)

	/*
		待完善事項
		1. 目前無法使用撮合(真實用戶)或配桌(機器人)逻辑, 需要主動安排入座
	*/

	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	//Test Round 1 正常開始
	robotManager := tgrs.RobotService.GetRobotManager()

	for i := 0; i < 3; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robotManager.AddRobot(robot)
	}
	if robotManager.RobotCount() != 3 {
		t.Errorf("加載機器人後數量應為3, RobotCount=>%d", robotManager.RobotCount())
	}

	//glc.poker = poker.NewRandomPoker(true, true, 1)
	//glc.poker = poker.NewFixedPoker(glc.AllTableUser(), glc.poker.(*poker.RandomPoker), 3)
	glc.setNextJiangType(0, int32(tongsha), "", "admin")

	users := fixtures.GetGameUser()

	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)
	//robots[2].UserState.OnUserQueue(robots[2], tgrs.Room)
	users[3].UserState.OnUserQueue(users[3], tgrs.Room)

	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	//glc.UserEnter(robots[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[3], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[1])
	//glc.UserReady(robots[2])
	glc.UserReady(users[3])

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNStateID) {
			ust := m.MsgData.(*netproto.BRNN_Sate)

			if *ust.State == Bet { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0
				//下注
				glc.zhuangID = robots[1].UserID

			}

			if *ust.State == KaiJiangAction {
				//userLen := fp.UserCardLen()
				//for i := 0; i < userLen; i++ {
				//	cards := fp.GetUserCards(int32(i))
				//	glc.allKindPoker[i] = cards //竄改想要開獎的牌驗證資料
				//}
			}

			if *ust.State == Wait {
				cards := new(netproto.BRNN_Poker)
				zhuangcardret := getNiuNiuRet(glc.allKindPoker[0])

				for i := 1; i < len(glc.brnnpokerresult); i++ {
					xiancardret := getNiuNiuRet(glc.allKindPoker[i])

					betret := compareCards(zhuangcardret, xiancardret)

					if betret < int32(0) { //莊輸
						win := int32(1)
						cards.RetType = append(cards.RetType, win)
						//	i, cards.RetType[i-1], zhuangcardret.niuRet, zhuangcardret.number, xiancardret.niuRet, xiancardret.number, zhuangcardret.cards, xiancardret.cards))
					}
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

//廣播第一手三張牌
func TestGetFirstCardGroup(t *testing.T) {
	conf.CLogLevel = mlog.ErrorLevel
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	//roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	//}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

	glc.poker = poker.NewRandomPoker(true, true, 1)

	glc.clearGameData()
	glc.fapeiFirstCardGroupNew()
	firstCardGroup := glc.getFirstCardGroup()

	firstCardData := firstCardGroup.MsgData.(*netproto.BRNN_BrodcastPokerCard)

	if len(firstCardData.UserPokerCard) != int(5) {
		t.Error("發牌組數錯誤")
	}

	for _, v := range firstCardData.UserPokerCard {
		if len(v.CardGroup) != int(3) {
			t.Error("發牌數量不等於3")
		}
	}

}

//廣播第二手兩張牌
func TestGetSecondCardGroup(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	//roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	//}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

	glc.poker = poker.NewRandomPoker(true, true, 1)

	glc.clearGameData()
	glc.fapeiFirstCardGroupNew()
	glc.getFirstCardGroup()
	glc.getResultType()
	glc.getKaiJiangPokerByType()
	firstCardGroup := glc.getSecondCardGroup()

	firstCardData := firstCardGroup.MsgData.(*netproto.BRNN_BrodcastPokerCard)

	if len(firstCardData.UserPokerCard) != int(5) {
		t.Error("發牌組數錯誤")
	}

	for _, v := range firstCardData.UserPokerCard {
		if len(v.CardGroup) != int(2) {
			t.Error("發牌數量不等於2")
		}
	}

	for _, v := range glc.allKindPoker {
		if len(v.Cards) != int(5) {
			t.Error("發牌數量不等於5")
		}
	}

}

//指定第二手牌型
func TestTakeNiuNiuByNumber(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	//roomserver.RobotLogicCreator = func(table *gamelogic.tableManager, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
	//	return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	//}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

	glc.poker = poker.NewRandomPoker(true, true, 1)

	glc.clearGameData()
	glc.fapeiFirstCardGroupNew()
	firstCardGroup := glc.getFirstCardGroup()

	secondPoker := glc.takeNiuNiuByNumber(glc.allKindPoker[0])

	if secondPoker == nil || len(secondPoker.Cards) != int(2) {
		t.Errorf("指定第二手牌型 失敗 , 前三張牌 %v", glc.allKindPoker[0].Cards)
	}

	firstCardData := firstCardGroup.MsgData.(*netproto.BRNN_BrodcastPokerCard)

	if len(firstCardData.UserPokerCard) != int(5) {
		t.Error("發牌組數錯誤")
	}

	for _, v := range firstCardData.UserPokerCard {
		if len(v.CardGroup) != int(3) {
			t.Error("發牌數量不等於3")
		}
	}
}

//比牌莊贏
func TestGetNiuNiuRet(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("brnncard")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret < int32(0) {
		t.Error("期望莊贏, 實際莊輸")
	}
}

//比牌
func TestGetNiuNiuRet3(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("card3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	getNiuNiuRet(fp.GetUserCards(int32(0)))

}

//比牌閒贏 皮子贏
func TestGetNiuNiuRetPiZi(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardpizi")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret > int32(0) {
		t.Error("皮子牌　期望閒贏, 實際閒輸")
	}
}

//比牌閒輸 皮子輸
func TestGetNiuNiuRetPiZi2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("card8")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret > int32(0) {
		t.Error("皮子牌　期望閒輸, 實際閒贏")
	}
}

//百搭 比牌牛牛閒輸, 皮子輸
func TestGetNiuNiuRetPiZi3(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("card9")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret < int32(0) {
		t.Error("皮子牌　期望閒輸, 實際閒贏")
	}
}

func TestGetNiuNiuRetPiZi4(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("card10")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret > int32(0) {
		t.Error("皮子牌　期望閒贏, 實際閒輸")
	}
}

//比牌莊贏 小王贏
func TestGetNiuNiuRetSmallKing(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardsmallking")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret < int32(0) {
		t.Error("小王牌　期望莊贏, 實際莊輸")
	}
}

//比牌閒贏兩張鬼牌
func TestGetNiuNiuRetTwoKing(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardtwoking")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret > int32(0) {
		t.Error("two王牌　期望閒贏, 實際閒輸")
	}
}

//是否有牛 -1表示无牛
func TestGetNiuNumber(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("niuNumber")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		card, _ := getNiuNumber(cards)

		if card != -1 {
			t.Error("無牛錯誤")
		}
	}
}

//是否有牛 表示有牛
func TestGetNiuNumber2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardPoint2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		card, cardStack := getNiuNumber(cards)

		if card != 8 {
			t.Errorf("有牛錯誤 => 點數 : %v, cardStack : %v", card, cardStack)
		}
	}
}

//是否牛牛
func TestIsNiuNiu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardPoint3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok, cardStack := isNiuNiu(cards)
		if !ok {
			t.Errorf("牛牛錯誤 => 點數 : %v, cardStack : %v", ok, cardStack)
		}
	}
}

//是否五花牛
func TestIsWuHuaNiu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isWuHuaNiu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isWuHuaNiu(cards)
		if !ok {
			t.Errorf("是否五花牛 => 點數 : %v, cardStack : %v", ok, cards)
		}
	}

}

//是否四炸
func TestIs4Zha(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("is4Zha")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok, cardStack := is4Zha(cards)
		if !ok {
			t.Errorf("是否四炸 => 點數 : %v, cardStack : %v", ok, cardStack)
		}
	}
}

//是否五小牛
func TestIsWuXiaoNiu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isWuXiaoNiu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isWuXiaoNiu(cards)
		if !ok {
			t.Errorf("是否五小牛 => 點數 : %v, cardStack : %v", ok, cards)
		}
	}
}

//單牌　閒贏
func TestGetNiuNiuRet2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("card2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	zhuangcardret := getNiuNiuRet(fp.GetUserCards(int32(0)))
	xiancardret := getNiuNiuRet(fp.GetUserCards(int32(1)))

	betret := compareCards(zhuangcardret, xiancardret)

	if betret > int32(0) {
		t.Error("小王牌　期望閒贏, 實際閒輸")
	}
}

//產生指定類型和大於某個點數的牛牛牌
func TestAsignNewNiuNiuCardMoreThen(t *testing.T) {
	conf.TestMode = true
	for i := 0; i < 500; i++ {
		roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
			return NewBRNNLogic(t, tfManager, cfgs)
		}

		tgrs := fixtures.NewTestGameRoomServer()
		glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

		glc.poker = poker.NewRandomPoker(true, true, 1)

		glc.clearGameData()
		glc.fapeiFirstCardGroupNew()
		glc.getFirstCardGroup()

		niuCardRet := glc.assignNewNiuNiuCardMoreThen(glc.allKindPoker[0], nnCard_Niu, 0)

		result := getNiuNiuRet(niuCardRet.cards)

		if result.niuRet == nnCard_WuNiu {
			t.Error("指定類型結果錯誤 : ", niuCardRet)
		}

	}
}

//產生指定類型和小於某個點數的牛牛牌
func TestAsignNewNiuNiuCardLessThen(t *testing.T) {
	conf.TestMode = true
	for i := 0; i < 500; i++ {
		roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
			return NewBRNNLogic(t, tfManager, cfgs)
		}

		tgrs := fixtures.NewTestGameRoomServer()
		glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

		glc.poker = poker.NewRandomPoker(true, true, 1)

		glc.clearGameData()
		glc.fapeiFirstCardGroupNew()
		glc.getFirstCardGroup()
		niuCardRet := glc.assignNewNiuNiuCardLessThen(glc.allKindPoker[0], nnCard_Niu, 7, 3)

		result := getNiuNiuRet(niuCardRet.cards)

		if result.niuRet < nnCard_Niu {
			t.Error("指定類型結果錯誤 : ", niuCardRet)
		}
	}
}

//真人當莊, 真人(贏)和機器人(輸)都下注
func TestGameRoundBetAndKaijaing12(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 500
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
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
	users[0].UserType = 1
	users[1].UserType = 1

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])
	glc.UserReady(users[1])

	var p poker.Poker
	p = poker.NewTestFilePoker("card4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.zhuangID = 9000001

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNStateID) {
			ust := m.MsgData.(*netproto.BRNN_Sate)

			if *ust.State == Bet { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0
				//下注
				glc.zhuangID = 9000001
				glc.userBet(users[1], fixtures.GetUserBet(2, 100))
				glc.userBet(robots[0], fixtures.GetUserBet(3, 100))
			}

			if *ust.State == KaiJiangAction {
				userLen := fp.UserCardLen()
				for i := 0; i < userLen; i++ {
					cards := fp.GetUserCards(int32(i))
					glc.allKindPoker[i] = cards //竄改想要開獎的牌驗證資料
				}
			}

			if *ust.State == Wait {
				if glc.syspool != int64(-300) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(-300) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(105) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(1995) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(2295) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.userdata[1].updatemoney != int64(-300) {
					t.Errorf("updatemoney:%v", glc.userdata[1].updatemoney)
				}

				if glc.userdata[1].win != int64(0) {
					t.Errorf("win:%v", glc.userdata[1].win)
				}

				if glc.zhuangWin != int64(-1800) {
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

//機器人當莊(通殺), 真人(輸)和機器人(輸)都下注
func TestGameRoundBetAndKaijaing15(t *testing.T) {
	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 500
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
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
	users[0].UserType = 1
	users[1].UserType = 1

	robots[0].UserState.OnUserQueue(robots[0], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)

	//安排用戶入座
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(robots[0])
	glc.UserReady(users[0])
	glc.UserReady(users[1])

	var p poker.Poker
	p = poker.NewTestFilePoker("card4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {

		if m.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNStateID) {
			ust := m.MsgData.(*netproto.BRNN_Sate)

			if *ust.State == Bet && users[1].UserID == userid { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0
				//下注
				glc.zhuangID = 9200011
				glc.userBet(users[1], fixtures.GetUserBet(3, 100))

			}

			if *ust.State == KaiJiangAction {
				userLen := fp.UserCardLen()
				for i := 0; i < userLen; i++ {
					cards := fp.GetUserCards(int32(i))
					glc.allKindPoker[i] = cards //竄改想要開獎的牌驗證資料
				}
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

				//if glc.userdata[1].updatemoney != int64(0) {
				//	t.Errorf("updatemoney:%v", glc.userdata[1].updatemoney)
				//}

				//if glc.userdata[1].win != int64(0) {
				//	t.Errorf("win:%v", glc.userdata[1].win)
				//}

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

//機器人當莊(贏) 機器人下注, 指定牌,下注→開獎→結果
func TestGameRoundBetAndKaijaing19(t *testing.T) {

	endChannel := make(chan bool)
	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 500
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBRNNRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
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
	robots[1].UserState.OnUserQueue(robots[1], tgrs.Room)
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	//安排用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(robots[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	var p poker.Poker
	p = poker.NewTestFilePoker("card5")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	glc.GameStartFunc = func() {

	}

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		if m.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNZhuangInfoID) {
			t.Log("莊")
		}

		if m.SClassID == int32(netproto.BRNN_GameMessageClassID_BRNNStateID) {
			ust := m.MsgData.(*netproto.BRNN_Sate)

			if *ust.State == Bet { //下注
				glc.syspool = 0
				glc.syslastwin = 0
				glc.systaxincome = 0
				glc.zhuangWin = 0
				//下注
				glc.zhuangID = 9200011
				glc.userBet(robots[1], fixtures.GetUserBet(4, 1000))
				glc.userBet(robots[1], fixtures.GetUserBet(1, 700))
				glc.userBet(robots[1], fixtures.GetUserBet(4, 100))
			}

			if *ust.State == KaiJiangAction {
				userLen := fp.UserCardLen()
				for i := 0; i < userLen; i++ {
					cards := fp.GetUserCards(int32(i))
					glc.allKindPoker[i] = cards //竄改想要開獎的牌驗證資料
				}
			}

			if *ust.State == Wait {
				if glc.syspool != int64(0) {
					t.Errorf("syspool:%v", glc.syspool)
				}

				if glc.syslastwin != int64(0) {
					t.Errorf("syslastwin:%v", glc.syslastwin)
				}

				if glc.systaxincome != int64(0) {
					t.Errorf("systaxincome:%v", glc.systaxincome)
				}

				if glc.userdata[0].updatemoney != int64(-54000) {
					t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
				}

				if glc.userdata[0].win != int64(-48600) {
					t.Errorf("win:%v", glc.userdata[0].win)
				}

				if glc.zhuangWin != int64(51300) {
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

//指定無牛牌型
func TestAssignWuNiuCard(t *testing.T) {
	//roomserver.GameLogicCreator = func(t *gamelogic.tableManager, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
	//	return NewBRNNLogic(t, tfManager, cfgs)
	//}
	//
	//tgrs := fixtures.NewTestGameRoomServer()
	//glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)
	//
	//for i := 0; i < 1; i++ {
	//	glc.poker = poker.NewRandomPoker(false, false, 1)
	//
	//	glc.clearGameData()
	//	glc.fapeiFirstCardGroup()
	//	glc.getFirstCardGroup()
	//
	//	niuCardRet := glc.assignWuNiuCard(glc.allKindPoker[0])
	//
	//	if niuCardRet == nil || niuCardRet.cards == nil {
	//		t.Logf("niuCardRet card : %v", niuCardRet.cards)
	//	}
	//
	//	result := getNiuNiuRet(niuCardRet.cards)
	//
	//	t.Logf("card : %v", result.cards)
	//	//
	//	//if result.niuRet != nnCard_WuNiu {
	//	//	t.Error("指定類型結果錯誤 : ", niuCardRet)
	//	//}
	//}

}

//指定牛幾牌行, 點數以下的牌行
func TestAssignHasNiuCardLessThen(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

	for i := 0; i < 100; i++ {
		glc.poker = poker.NewRandomPoker(false, false, 1)

		glc.clearGameData()
		glc.fapeiFirstCardGroupNew()
		glc.getFirstCardGroup()

		number := 7
		niuCardRet := glc.assignNewNiuNiuCardLessThen(glc.allKindPoker[0], nnCard_Niu, number, 1000)

		result := getNiuNiuRet(niuCardRet.cards)

		if result.niuRet < nnCard_Niu && result.number < number {
			t.Error("指定類型結果錯誤 : ", niuCardRet)
		}
	}

}

//指定牛幾牌行, 點數以上的牌行
func TestAssignHasNiuCardMoreThen(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBRNNLogic(t, tfManager, cfgs)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)

	for i := 0; i < 100; i++ {

		glc.poker = poker.NewRandomPoker(false, false, 1)

		glc.clearGameData()
		glc.fapeiFirstCardGroupNew()
		glc.getFirstCardGroup()

		number := 3
		niuCardRet := glc.assignNewNiuNiuCardMoreThen(glc.allKindPoker[0], nnCard_Niu, number)

		result := getNiuNiuRet(niuCardRet.cards)

		if result.niuRet >= nnCard_Niu && result.number < number {
			t.Error("指定類型結果錯誤 : ", niuCardRet)
		}
	}

}

//流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *BRNNLogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BRNNLogic)
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

func NewBRNNLogicByTest() *BRNNLogic {
	_, filename, _, _ := runtime.Caller(1)
	configDir := path.Dir(filename) + "/../../config"
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameGameServer)
	serverInfo := cfgs.GetConfig(gamelogic.ConfigFileNameServerInfo)

	//获取配置信息
	servercfg := gamelogic.GetGameServerConfig(serverInfo, cfg)
	rumanager := gamelogic.NewRoomUserManager(new(player.Manager))
	room := gamelogic.NewRoom(servercfg, rumanager)
	table := room.GetTableByNo(0)

	evtdispatcher := event.NewAsyncEventDispatcher(1)
	tfManager := event.NewTimerFuncManager(evtdispatcher)

	glc := new(BRNNLogic)

	glc.ComputeGameLogicBase = gamelogic.NewComputeGameLogicBase(table, tfManager, cfgs)
	glc.getbrnnconfig(cfgs)
	glc.gamestate = Wait
	glc.zhuangID = 0
	glc.nextjiangtype = -1
	glc.syslastwin = 0
	glc.resultcardlist = new(netproto.BRNN_ResultCardList)
	glc.startsyspool = glc.syspool

	//设置桌子对应的逻辑实例
	for _, t := range room.Tables {
		t.GameLogic = new(BRNNLogic)
	}

	glc.RegisterAfterGameStartFunc(func() {
		glc.clearGameData() //清除数据
		glc.huangZhuang()
		glc.GameLog.AddGameStep("庄家%d，剩余%d", glc.zhuangID, glc.lastturn)
		glc.zhuShui()
		glc.broadCastZhuangInfo()
		//glc.startBet()
	})

	return glc
}
