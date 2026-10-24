package bjl

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func init() {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBaiJiaLeLogic(t, tfManager, cfgs)
	}
	roomserver.SuperToolCreator = func(room *gamelogic.Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) gamelogic.ISuperToolService {
		return gamelogic.NewBrGameSuperToolServiceBase(room, tfm, dbMsgQueue)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewBJLRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
}

func fapaiTool(bjl *BaiJialeLogic) (xiancards []*card, zhuangcards []*card) {
	defer util.PanicErrStack()
	limit := 1000
	var times = 0
	for {
		bjl.fapai()
		times++
		if times > limit || bjl.cards.Len() < 6 {
			bjl.newcard = true
			times = 0
			continue
		}
		var ok bool
		xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
		if !ok {
			continue
		}
		//t.Logf("闲家的扑克牌是%v", xian)
		zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
		if ok {
			break
		}
		cArr := []int32{0, 1, 2, 3}
		for _, c := range append(xiancards, zhuangcards...) {
			if ok, _ := util.InArray(c.color, cArr); !ok {
				mlog.Error("沒有花色的牌 ")
			}
		}
		time.Sleep(time.Duration(15) * time.Millisecond)
	}
	return
}

// todo 對子的判斷
func TestNatualPai(t *testing.T) {

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}

	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 5000
	conf.TestMode = true

	tgrs := fixtures.NewTestGameRoomServer()
	bjl := tgrs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)

	for i := 0; i < 600; i++ {
		if i%30 == 0 {
			bjl.newcard = true
		}
		bjl.fapai()
		bjl.getResultType(10000)
		bjl.getPokerByType(bjl.jiangtype)
		bjl.roundid = int32(i)
		//bjl.jiangtype = []int32{1}
		//bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		//bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		//bjl.zuhetype |= 1 << uint32(HE-1) //和
		//bjl.zuhetype |= 1 << uint32(6-1) //莊對
		//bjl.zuhetype |= 1 << uint32(7-1) //閒對
		//xiancards, zhuangcards := fapaiTool(bjl)
		//
		//xian := poker.NewEmptyCardGroup(0)
		//for i := 0; i < len(xiancards); i++ {
		//	card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
		//	if card.CardNumber == 1 {
		//		card.CardNumber = 14
		//	}
		//	xian.AddCard(card)
		//}
		//zhuang := poker.NewEmptyCardGroup(0)
		//
		//for i := 0; i < len(zhuangcards); i++ {
		//	card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
		//	if card.CardNumber == 1 {
		//		card.CardNumber = 14
		//	}
		//	zhuang.AddCard(card)
		//}
		//
		//if bjl.getCardsDian(zhuangcards) != bjl.getCardsDian(xiancards) {
		//	//t.Errorf("錯誤 : 开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		//}
		//cArr := []int32{0,1,2,3}
		//for _, c := range append(bjl.cardxian, bjl.cardzhuang...) {
		//	if ok, _ := util.InArray(c.color, cArr);!ok {
		//		t.Errorf("沒有花色的牌 ")
		//	}
		//}
		//
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		bjl.clearGameData()
		//time.Sleep(time.Duration(15) * time.Millisecond)
	}

}

// todo 對子的判斷
func TestHeXianDui(t *testing.T) {

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}

	conf.CLogLevel = mlog.ErrorLevel
	conf.TimerFastMode = 5000
	conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewBaiJiaLeLogic(t, tfManager, cfgs)
	}
	tgrs := fixtures.NewTestGameRoomServer()
	bjl := tgrs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)

	for i := 0; i < 30; i++ {
		//bjl.getResultType(10000)
		bjl.roundid = int32(i)
		bjl.jiangtype = []int32{1}
		//bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		//bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		bjl.zuhetype |= 1 << uint32(HE-1)        //和
		bjl.zuhetype |= 1 << uint32(ZHUANGDUI-1) //莊對
		//bjl.zuhetype |= 1 << uint32(XIANDUI-1) //閒對
		xiancards, zhuangcards := fapaiTool(bjl)

		xian := poker.NewEmptyCardGroup(0)
		for i := 0; i < len(xiancards); i++ {
			card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			xian.AddCard(card)
		}
		zhuang := poker.NewEmptyCardGroup(0)

		for i := 0; i < len(zhuangcards); i++ {
			card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			zhuang.AddCard(card)
		}
		cArr := []int32{0, 1, 2, 3}
		for _, c := range xiancards {
			if ok, _ := util.InArray(c.color, cArr); !ok {
				t.Errorf("錯誤 : 开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
			}
		}
		for _, c := range zhuangcards {
			if ok, _ := util.InArray(c.color, cArr); !ok {
				t.Errorf("錯誤 : 开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
			}
		}
		if bjl.getCardsDian(zhuangcards) != bjl.getCardsDian(xiancards) {
			t.Errorf("錯誤 : 开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		}
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		bjl.clearGameData()
		time.Sleep(time.Duration(1) * time.Millisecond)
	}

}

// todo 對子的判斷
func TestXianWinXianDui(t *testing.T) {
	bjl := new(BaiJialeLogic)
	for i := 0; i < 600; i++ {
		//bjl.jiangtype = []int32{1}
		//bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		//bjl.zuhetype |= 1 << uint32(6-1) //莊對
		bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		bjl.zuhetype |= 1 << uint32(7-1)    //閒對
		bjl.fapai()
		var xiancards []*card
		var zhuangcards []*card

		var ok bool
		xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
		if !ok {
			continue
		}
		xian := poker.NewEmptyCardGroup(0)
		for i := 0; i < len(xiancards); i++ {
			card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			xian.AddCard(card)
		}
		//t.Logf("闲家的扑克牌是%v", xian)
		zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
		if !ok {
			continue
		}

		zhuang := poker.NewEmptyCardGroup(0)

		for i := 0; i < len(zhuangcards); i++ {
			card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			zhuang.AddCard(card)
		}

		if bjl.getCardsDian(zhuangcards) >= bjl.getCardsDian(xiancards) {
			t.Errorf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		}
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
	}

}

// todo 對子的判斷
func TestXianWinZhuangDui(t *testing.T) {
	bjl := new(BaiJialeLogic)
	for i := 0; i < 600; i++ {
		//bjl.jiangtype = []int32{1}
		//bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		bjl.zuhetype |= 1 << uint32(6-1)    //莊對
		bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		//bjl.zuhetype |= 1 << uint32(7-1) //閒對
		bjl.fapai()
		var xiancards []*card
		var zhuangcards []*card

		var ok bool
		xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
		if !ok {
			continue
		}
		xian := poker.NewEmptyCardGroup(0)
		for i := 0; i < len(xiancards); i++ {
			card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			xian.AddCard(card)
		}
		//t.Logf("闲家的扑克牌是%v", xian)
		zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
		if !ok {
			continue
		}

		zhuang := poker.NewEmptyCardGroup(0)

		for i := 0; i < len(zhuangcards); i++ {
			card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			zhuang.AddCard(card)
		}

		if bjl.getCardsDian(zhuangcards) >= bjl.getCardsDian(xiancards) {
			t.Errorf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		}
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
	}

}

// todo 對子的判斷
func TestZhuangWinZhuangDui(t *testing.T) {
	bjl := new(BaiJialeLogic)
	for i := 0; i < 600; i++ {
		//bjl.jiangtype = []int32{1}
		bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		bjl.zuhetype |= 1 << uint32(6-1)      //莊對
		//bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		//bjl.zuhetype |= 1 << uint32(7-1) //閒對
		bjl.fapai()
		var xiancards []*card
		var zhuangcards []*card

		var ok bool
		xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
		if !ok {
			continue
		}
		xian := poker.NewEmptyCardGroup(0)
		for i := 0; i < len(xiancards); i++ {
			card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			xian.AddCard(card)
		}
		//t.Logf("闲家的扑克牌是%v", xian)
		zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
		if !ok {
			continue
		}

		zhuang := poker.NewEmptyCardGroup(0)

		for i := 0; i < len(zhuangcards); i++ {
			card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			zhuang.AddCard(card)
		}

		if bjl.getCardsDian(zhuangcards) <= bjl.getCardsDian(xiancards) {
			t.Errorf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		}
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
	}

}

// todo 對子的判斷
func TestZhuangWinXianDui(t *testing.T) {
	bjl := new(BaiJialeLogic)
	for i := 0; i < 600; i++ {
		//bjl.jiangtype = []int32{1}
		bjl.zuhetype |= 1 << uint32(ZHUANG-1) //庄贏
		//bjl.zuhetype |= 1 << uint32(6-1) //莊對
		//bjl.zuhetype |= 1 << uint32(XIAN-1) //閒贏
		bjl.zuhetype |= 1 << uint32(7-1) //閒對
		bjl.fapai()
		var xiancards []*card
		var zhuangcards []*card

		var ok bool
		xiancards, ok = bjl.getXDuiCards(bjl.zuhetype&(1<<1) != 0, bjl.zuhetype&(1<<6) != 0)
		if !ok {
			continue
		}
		xian := poker.NewEmptyCardGroup(0)
		for i := 0; i < len(xiancards); i++ {
			card := poker.NewPokerCard(xiancards[i].color, xiancards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			xian.AddCard(card)
		}
		//t.Logf("闲家的扑克牌是%v", xian)
		zhuangcards, ok = bjl.getZDuiCards(xiancards, bjl.zuhetype&1 != 0, bjl.zuhetype&(1<<5) != 0)
		if !ok {
			continue
		}

		zhuang := poker.NewEmptyCardGroup(0)

		for i := 0; i < len(zhuangcards); i++ {
			card := poker.NewPokerCard(zhuangcards[i].color, zhuangcards[i].number)
			if card.CardNumber == 1 {
				card.CardNumber = 14
			}
			zhuang.AddCard(card)
		}

		if bjl.getCardsDian(zhuangcards) <= bjl.getCardsDian(xiancards) {
			t.Errorf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
		}
		//t.Logf("开奖结果是%d,庄家的扑克牌是(%d)%v ,闲家的扑克牌是(%d)%v", bjl.jiangtype, bjl.getCardsDian(zhuangcards), zhuang, bjl.getCardsDian(xiancards), xian)
	}
}

func TestZuheType(t *testing.T) {
	var zuhetype uint32
	zuhetype |= 1 << uint32(XIAN-1)
	t.Logf("%b", zuhetype) //庄贏
	t.Logf("%v", zuhetype&(1<<1))
	//zuhetype |= 1 << uint32(4-1) //庄天王
	//t.Logf("%b", zuhetype) //庄贏 + 庄天王
	zuhetype |= 1 << uint32(5-1)
	t.Logf("%b", zuhetype) //庄贏 + 庄天王 + 閒天王
	zuhetype |= 1 << uint32(6-1)
	t.Logf("%b", zuhetype) //庄贏 + 庄天王 + 閒天王 + 庄對
	zuhetype |= 1 << uint32(7-1)
	t.Logf("%b", zuhetype)        //庄贏 + 庄天王 + 閒天王 + 庄對 + 閒對
	t.Logf("%b", zuhetype&(1<<3)) //庄天王

	t.Logf("%b", 1<<1)
	t.Logf("%b", 1<<2)

	t.Logf("%b", 2&1)

}

//获取上庄结果消息
func TestGetShangzhuangRetMsg(t *testing.T) {

	//roomserver.StartGameLogicServer()

	bjl := new(BaiJialeLogic)

	result := bjl.getShangzhuangRetMsg(0)

	//Unmarshal
	msgparser := msg.NewParserGroup()
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BJL_GameMessageClassID_BJLShangZhuangRetID), msg.NewProtoMsgParser())

	if result.BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("result [%v] expected [%v]", result.BClassID, int32(netproto.MessageBClassID_Game))
	}

	if result.SClassID != int32(netproto.BJL_GameMessageClassID_BJLShangZhuangRetID) {
		t.Errorf("result [%v] expected [%v]", result.SClassID, int32(netproto.BJL_GameMessageClassID_BJLShangZhuangRetID))
	}

	msgData := fmt.Sprint(result.MsgData)

	expect := "ret:0 "
	if msgData != expect {
		t.Errorf("result [%s] expected [%s]", result.MsgData, expect)
	}
}

//获取赔钱最少的开奖类型
//func TestGetMinJiangType(t *testing.T) {
//
//	successCount := map[int]int{}
//	for i := 0; i < 4; i++ { //測試機率的次數
//
//		alltype := make([]int32, 0)
//		for i := 0; i < 3; i++ {
//			alltype = append(alltype, int32(i+1))
//		}
//
//		glc := new(BaiJialeLogic)
//
//		poollinevalue := new(poollinearr)
//		poollinevalue.limitmoney = 20000
//		poollinevalue.rate = 100
//
//		cf := new(poollineconfig)
//		cf.pool = make([]*poollinearr, 0)
//		cf.pool = append(cf.pool, poollinevalue)
//
//		glc.poolline = cf
//		//glc.syspool = 100
//		glc.betamode = 1
//
//		bcf := new(betconfig)
//		bcf.BetValue = make([]*betarr, 0)
//		bcf.AllType = make([]int32, 0)
//		value1 := new(betarr)
//		value1.Baijialetype = 1
//		value1.Odd = 8
//		value1.Money = 2
//		value1.Maxmoney = 100
//		value1.Weight = 10000000
//		bcf.BetValue = append(bcf.BetValue, value1)
//		bcf.AllType = append(bcf.AllType, value1.Baijialetype)
//		value2 := new(betarr)
//		value2.Baijialetype = 2
//		value2.Odd = 8
//		value2.Money = 2
//		value2.Maxmoney = 100
//		value2.Weight = 10000000
//		bcf.BetValue = append(bcf.BetValue, value2)
//		bcf.AllType = append(bcf.AllType, value2.Baijialetype)
//		value3 := new(betarr)
//		value3.Baijialetype = 3
//		value3.Odd = 1
//		value3.Money = 9
//		value3.Maxmoney = 100
//		value3.Weight = 10000000
//		bcf.BetValue = append(bcf.BetValue, value3)
//		bcf.AllType = append(bcf.AllType, value3.Baijialetype)
//
//		glc.betvalue = bcf
//
//		glc.updateBet(1, 1000)
//
//		alltype = glc.getminjiangtype(alltype)
//
//		//if AllType[0] == 1 {
//		//	t.Log("最有利開獎:莊")
//		//}
//
//		if alltype[0] == 2 {
//			if _, ok := successCount[2]; ok {
//				continue
//			}
//			successCount[2] = 2
//			t.Log("最有利開獎:閒")
//		}
//
//		if alltype[0] == 2 {
//			if _, ok := successCount[3]; ok {
//				continue
//			}
//			successCount[3] = 3
//			t.Log("最有利開獎:和")
//		}
//	}
//
//	expect := 2
//	if len(successCount) != expect {
//		t.Errorf("最有利開獎未開出閒或和")
//	}
//
//}

//測試獲取莊信息
func TestGetZhuangInfoMsg(t *testing.T) {

	glc := NewBaiJiaLeLogicByTest()
	glc.zhuangID = 0
	result := glc.getZhuangInfoMsg()
	msg := result.MsgData.(*netproto.BJL_ZhuangInfo)
	if msg.GetUserID() != 0 {
		t.Errorf("獲取莊家信息錯誤!")
	}
}

func TestUseRobotShangZhuang(t *testing.T) {
	glc := NewBaiJiaLeLogicByTest()
	glc.zhuangID = 0
	users := fixtures.GetGameUser()

	seat := glc.GetTable().GetFreeSeat()
	seat.SeatUser = users[0]

	glc.useRobotShangZhuang()

	result := glc.getZhuangInfoMsg()
	msg := result.MsgData.(*netproto.BJL_ZhuangInfo)
	if msg.GetUserID() == 0 {
		t.Errorf("機器人上莊失敗!")
	}

}

////水池注水 系統當莊
//func TestZhuShui(t *testing.T){
//	bjl := new(BaiJialeLogic)
//	glcb := new(gamelogic.ComputeGameLogicBase)
//	glb := new(gamelogic.GameLogicBase)
//	glb.GetServerConfig().TaxScale = 5
//	glcb.GameLogicBase = glb
//	bjl.ComputeGameLogicBase = glcb
//
//	bjl.syspool = 1000
//	bjl.zhuShui()
//
//	bjl.bettotal = 50 //下注總額
//	bjl.kaiJiang()
//
//}

//天王->對子
func TestGetTianWangIsDui(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("tianWangIsDui")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	successCount := map[int]int{}
	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards
		isDui := true
		isWin := true
		isZhuang := true

		cardzhuang := bjl.getTianWang(isDui, isWin, isZhuang)
		for k, v := range cardzhuang {

			t.Logf("index[%v] number[%v] color[%v]", k, v.number, v.color)
			if v.number == 4 {
				if _, ok := successCount[4]; ok {
					successCount[4]++
					continue
				}

			}

			if v.number == 9 {
				if _, ok := successCount[9]; ok {
					successCount[9]++
					continue
				}

			}
		}

	}

	expect := 2
	if len(successCount) >= expect {
		t.Errorf("未開出天王->對子")
	}

	for _, v := range successCount {
		if v >= expect {
			t.Errorf("未開出天王->對子")
		}
	}

}

//天王->莊
func TestGetTianWangIsZhuang(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("tianWangIsZhuang")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := false
		isWin := true
		isZhuang := true
		bjl.zuhetype |= 1 << uint32(4)
		cardzhuang := bjl.getTianWang(isDui, isWin, isZhuang)

		expect := int32(9)
		card1 := int32(0)
		card2 := int32(0)
		if cardzhuang[0].number <= 9 {
			card1 = cardzhuang[0].number
		}
		if cardzhuang[1].number <= 9 {
			card2 = cardzhuang[1].number
		}

		actual := (card1 + card2) % 10
		if actual != expect {
			t.Errorf("未開出天王->莊 , actual[%v] expect[%v], cardzhuang1[%v] + cardzhuang2[%v]", actual, expect, cardzhuang[0].number, cardzhuang[1].number)
		}

	}

}

//天王->閒
func TestGetTianWangIsZhuang2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("tianWangIsZhuang2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := false
		isWin := true
		isZhuang := false
		bjl.zuhetype |= 1 << uint32(3)
		cardzhuang := bjl.getTianWang(isDui, isWin, isZhuang)

		expect8 := int32(9)
		expect9 := int32(9)
		card1 := int32(0)
		card2 := int32(0)
		if cardzhuang[0].number <= 9 {
			card1 = cardzhuang[0].number
		}
		if cardzhuang[1].number <= 9 {
			card2 = cardzhuang[1].number
		}

		actual := (card1 + card2) % 10
		if actual != expect9 && actual != expect8 {
			t.Errorf("未開出天王->閒 , actual[%v] expect9[%v], cardzhuang1[%v] + cardzhuang2[%v]", actual, expect9, cardzhuang[0].number, cardzhuang[1].number)
		}

	}

}

//获取赢的三张牌型： 對子, 當前兩張點數和的個位數是0~5就發第三張牌. 10~13為0
//第三張同第一張牌點數
func TestGetThreefirstCards(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("threefirstCards")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := true
		isShun := false
		isWin := true
		isZhuang := true
		bjl.zuhetype |= 1 << uint32(3)

		cardzhuang, ok := bjl.getThreefirstCards(isDui, isShun, isWin, isZhuang)

		if !ok {
			t.Logf("获取赢的三张牌型： 對子 莊 失敗")
		}

		actual1 := int32(0)
		if cardzhuang[0].number < 10 {
			actual1 = cardzhuang[0].number
		}
		actual2 := int32(0)
		if cardzhuang[1].number < 10 {
			actual2 = cardzhuang[0].number
		}

		if (actual1+actual2)%10 < 6 {
			if cardzhuang[0].number != cardzhuang[2].number &&
				cardzhuang[1].number != cardzhuang[2].number &&
				cardzhuang[0].number != cardzhuang[1].number {
				t.Errorf("取赢的三张牌型： 對子 莊 失敗, %v, %v, %v", cardzhuang[0], cardzhuang[1], cardzhuang[2])
			}
		}
	}

}

//获取赢的三张牌型： 順子, 當前兩張點數和的個位數是0~5就發第三張牌. 10~13為0
//第三張同第一張牌點數
func TestGetThreefirstCards2(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("threefirstCards2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := false //對子
		isShun := true //順子
		isWin := true
		isZhuang := true
		bjl.zuhetype |= 1 << uint32(3)

		cardzhuang, ok := bjl.getThreefirstCards(isDui, isShun, isWin, isZhuang)

		if !ok {
			t.Logf("获取赢的三张牌型： 順子 莊 失敗")
		}

		actual1 := int32(0)
		if cardzhuang[0].number < 10 {
			actual1 = cardzhuang[0].number
		}
		actual2 := int32(0)
		if cardzhuang[1].number < 10 {
			actual2 = cardzhuang[1].number
		}
		actual3 := int32(0)
		if cardzhuang[2].number < 10 {
			actual3 = cardzhuang[2].number
		}

		if actual1 != (actual2-1) && actual1 != (actual3-2) {
			t.Errorf("錯誤 => 获取赢的三张牌型： 莊順子")
		}

	}

}

//获取赢的三张牌型： 同點和為第三張是對子,
func TestGetThreefirstCards3(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("threefirstCards3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := true
		isShun := false
		isWin := true
		isZhuang := true
		bjl.zuhetype |= 1 << uint32(9)

		cardzhuang, ok := bjl.getThreefirstCards(isDui, isShun, isWin, isZhuang)

		if !ok {
			t.Logf("获取赢的三张牌型： 同點和 失敗")
		}

		actual1 := int32(0)
		if cardzhuang[0].number < 10 {
			actual1 = cardzhuang[0].number
		}
		actual2 := int32(0)
		if cardzhuang[1].number < 10 {
			actual2 = cardzhuang[0].number
		}

		if (actual1+actual2)%10 < 6 {
			if cardzhuang[0].number != cardzhuang[2].number &&
				cardzhuang[1].number != cardzhuang[2].number {
				t.Errorf("取赢的三张牌型： 同點和對子 莊 失敗")
			}
		}
	}

}

//获取赢的三张牌型： 同點和為第三張是順子,
func TestGetThreefirstCards4(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("threefirstCards4")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		bjl := new(BaiJialeLogic)
		bjl.cards = cards

		isDui := false
		isShun := true
		isWin := true
		isZhuang := true
		bjl.zuhetype |= 1 << uint32(9)

		cardzhuang, ok := bjl.getThreefirstCards(isDui, isShun, isWin, isZhuang)

		if !ok {
			t.Logf("获取赢的三张牌型： 同點和順子 失敗")
		}

		actual1 := int32(0)
		if cardzhuang[0].number < 10 {
			actual1 = cardzhuang[0].number
		}
		actual2 := int32(0)
		if cardzhuang[1].number < 10 {
			actual2 = cardzhuang[1].number
		}
		actual3 := int32(0)
		if cardzhuang[2].number < 10 {
			actual3 = cardzhuang[2].number
		}

		if actual1 != (actual2-1) && actual1 != (actual3-2) {
			t.Errorf("錯誤 => 获取赢的三张牌型： 同點和莊順子")
		}
	}

}

////系統當莊, 下注龍開龍, 遊戲流程入桌→下注→開獎→結果
//func TestGameRound4(t *testing.T) {
//	bjl := NewBaiJiaLeLogicByTest()
//	bjl.robotzhuang = 0
//	bjl.syspool = 0
//	bjl.SendMessgeFunc = func(userid int32, m *msg.Message) {
//	}
//	bjl.OnGameStart(func() {
//	})
//	bjl.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
//	}
//	bjl.GameEndFunc = func() {
//	}
//	bjl.GameStartFunc()
//	users := fixtures.GetGameUser()
//	users[0].UserType = 1
//
//	//安排用戶入座
//	bjl.UserEnter(users[0], bjl.GetTable().TableNo, bjl.GetTable().GetFreeSeat().SeatNo)
//	//bjl.UserEnter(users[1], bjl.GetTable().TableNo, bjl.GetTable().GetFreeSeat().SeatNo)
//
//	//下注
//	bjl.userBet(users[0], fixtures.GetUserBet(ZHUANG, 100))
//	//bjl.userBet(users[0], fixtures.GetUserBet(HE, 100))
//	//bjl.userBet(users[0], fixtures.GetUserBet(HU, 100))
//
//	//bjl.setNextJiangType(0, 1)
//	bjl.nextjiangtype = int32(ZHUANG)
//
//	bjl.kaiJiang()
//	//bjl.sendResult()
//	calWinLose(bjl)
//
//	if bjl.syspool != int64(-10000) {
//		t.Errorf("syspool:%v", bjl.syspool)
//	}
//
//	if bjl.systotalwin != int64(-10000) {
//		t.Errorf("systotalwin:%v", bjl.systotalwin)
//	}
//
//	if bjl.syslastwin != int64(-10000) {
//		t.Errorf("syslastwin:%v", bjl.syslastwin)
//	}
//
//	if bjl.systaxincome != int64(500) {
//		t.Errorf("systaxincome:%v", bjl.systaxincome)
//	}
//
//	if bjl.userdata[0].updatemoney != int64(9500) {
//		t.Errorf("updatemoney:%v", bjl.userdata[0].updatemoney)
//	}
//
//	if bjl.userdata[0].win != int64(19500) {
//		t.Errorf("win:%v", bjl.userdata[0].win)
//	}
//
//	if bjl.zhuangWin != int64(-10000) {
//		t.Errorf("zhuangWin:%v", bjl.zhuangWin)
//	}
//
//	bjl.gameEnd()
//	bjl.clearGameData()
//}
//
func NewBaiJiaLeLogicByTest() *BaiJialeLogic {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	return glc
}

func TestGameRound(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.ComputeGameLogicBase.LogicBase)
	glc.robotzhuang = 0
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
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座
	if glc.GetCurrentGameUserCount() != 0 {
		t.Errorf("用戶數量有問題, should be 0 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}

	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	if glc.GetCurrentGameUserCount() != 2 {
		t.Errorf("用戶數量有問題, should be 1 but GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}

	t.Log("shouldbetotal => ", float64(glc.bettotal+700)*0.95)
	glc.userBet(users[0], fixtures.GetUserBet(ZHUANG, 1))
	glc.userBet(users[0], fixtures.GetUserBet(HE, 1))
	glc.userBet(users[1], fixtures.GetUserBet(XIAN, 3))
	glc.userBet(users[1], fixtures.GetUserBet(HE, 2))
	glc.UserOffLine(users[1])
	if len(glc.AllOnLineTableUser()) != 1 {
		t.Error("在線人數計算有誤 ", len(glc.AllOnLineTableUser()))
	}
	glc.UserRecover(users[1])
	if len(glc.AllOnLineTableUser()) != 2 {
		t.Error("在線人數計算有誤 ", len(glc.AllOnLineTableUser()))
	}
	robotZhuang := glc.getRobotBetMoneyByType(ZHUANG)
	robotXian := glc.getRobotBetMoneyByType(XIAN)
	robotHe := glc.getRobotBetMoneyByType(HE)

	glc.clearUserBet(users[1])
	glc.userBet(users[1], fixtures.GetUserBet(XIAN, 3))
	glc.userBet(users[1], fixtures.GetUserBet(HE, 2))

	if robotZhuang != 1*100 || robotXian != 0 || robotHe != 1*100 {
		t.Error("機器人的下注金額不符")
	}

	totalZhuang := glc.getBetMoneyByType(ZHUANG)
	totalXian := glc.getBetMoneyByType(XIAN)
	totalHe := glc.getBetMoneyByType(HE)

	if totalZhuang != 1*100 || totalXian != 300 || totalHe != 1*300 {
		t.Error("總下注金額不符")
	}

	glc.betEnd()
	glc.kaiJiang()

	t.Log("betotal => ", glc.bettotal)
	glc.sendResult()

	//if ok, _ := util.InArray(int32(ZHUANG), glc.jiangtype); !ok {
	//	t.Error("開獎類型有問題, ", ZHUANG, glc.jiangtype)
	//}
	t.Log("final syspool => ", glc.SysPool.GetPool())

	glc.gameEnd()
	glc.clearGameData()

	if glc.GetCurrentGameUserCount() != 2 {
		t.Errorf("用戶數量有問題, GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserFree(users[0])
	glc.UserLeave(users[0])
	if glc.GetCurrentGameUserCount() != 1 {
		t.Errorf("用戶數量有問題, GetCurrentGameUserCount %d", glc.GetCurrentGameUserCount())
	}
	glc.UserFree(users[1])
	glc.UserLeave(users[1])

}

//func TestStartGameLogicServer(t *testing.T) {
//
//	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
//		return NewBaiJiaLeLogic(t, tfManager, cfgs)
//	}
//
//	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.ConfigGroup, glc gamelogic.GameLogic, evtChannel *event.EventDispatcher) robotlogic.RobotLogic {
//		return NewBJLRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
//	}
//
//	fixtures.StartGameLogicByTest()
//}

func enterGameUser() (*BaiJialeLogic, []*gamelogic.GameUser) {
	bjl := NewBaiJiaLeLogicByTest()
	bjl.robotzhuang = 0
	bjl.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	bjl.OnGameStart(func() {
	})
	bjl.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	bjl.GameEndFunc = func() {
	}
	bjl.GameStartFunc()
	users := fixtures.GetGameUserLhd()
	seat := bjl.GetTable().GetFreeSeat()
	seat.SeatUser = users[0]
	seat.SeatUser = users[1]
	//
	//u := new(longhudouuserdata)
	//u.user = users[0]
	//bjl.userdata = append(bjl.userdata, u)
	//u1 := new(longhudouuserdata)
	//u1.user = users[1]
	//bjl.userdata = append(bjl.userdata, u1)

	bjl.zhuangID = 1
	bjl.startmoney = users[1].Score
	//安排用戶入座
	bjl.UserEnter(users[0], bjl.GetTable().TableNo, bjl.GetTable().GetFreeSeat().SeatNo)
	bjl.UserEnter(users[1], bjl.GetTable().TableNo, bjl.GetTable().GetFreeSeat().SeatNo)
	bjl.UserEnter(users[2], bjl.GetTable().TableNo, bjl.GetTable().GetFreeSeat().SeatNo)
	return bjl, users
}

//玩家當莊
//玩家0, 下注龍和 開虎, 遊戲流程入桌→下注→開獎→結果
//玩家1, 下注虎和 開虎, 遊戲流程入桌→下注→開獎→結果
func TestGameRound2(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	bjl, users := enterGameUser()

	//player1 下注
	bjl.userBet(users[0], fixtures.GetUserBet(ZHUANG, 100))
	bjl.userBet(users[0], fixtures.GetUserBet(HE, 100))

	//player2 下注
	bjl.userBet(users[2], fixtures.GetUserBet(XIAN, 100)) //
	bjl.userBet(users[2], fixtures.GetUserBet(HE, 50))
	bjl.jiangtype = make([]int32, 0)
	//bjl.setNextJiangType(0, 1)
	bjl.nextjiangtype = int32(XIAN)

	bjl.kaiJiang()
	//bjl.sendResult()
	calWinLose(bjl)

	if bjl.SysPool.GetPool() != int64(1500000) {
		t.Errorf("syspool:%v", bjl.SysPool.GetPool())
	}

	if bjl.systotalwin != int64(14250) {
		t.Errorf("systotalwin:%v", bjl.systotalwin)
	}

	if bjl.systaxincome != int64(1000) {
		t.Errorf("systaxincome:%v", bjl.systaxincome)
	}

	//玩家1輸贏
	if bjl.userdata[0].updatemoney != int64(-20000) {
		t.Errorf("updatemoney:%v", bjl.userdata[0].updatemoney)
	}

	//玩家1輸贏
	if bjl.userdata[0].win != int64(0) {
		t.Errorf("win:%v", bjl.userdata[0].win)
	}

	//玩家2輸贏
	if bjl.userdata[1].updatemoney != int64(4750) {
		t.Errorf("updatemoney:%v", bjl.userdata[1].updatemoney)
	}

	//玩家2輸贏
	if bjl.userdata[1].win != int64(19750) {
		t.Errorf("win:%v", bjl.userdata[1].win)
	}

	//莊輸贏
	if bjl.zhuangWin != int64(14250) {
		t.Errorf("zhuangWin:%v", bjl.zhuangWin)
	}

	bjl.gameEnd()
	bjl.clearGameData()
}

//玩家1當莊
//玩家2, 下注莊順子及莊, 開莊赢及莊順子

func TestGameRound3(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	//glc.syspool = 0
	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.ComputeGameLogicBase.LogicBase)
	glc.SysPool = gamelogic.NewSysPoolHelper(glc.LogicBase, glc.Configs)
	println(glc.SysPool.ToString())
	glc.robotzhuang = 0
	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}
	glc.GameStartFunc()
	glc.betvalue = nil
	/*
		about JsonBetValue
		bytes2, _ := json.Marshal(glc.betvalue)
		s := base64.StdEncoding.EncodeToString(bytes2)
		fmt.Println(s)
	*/
	JsonBetValue := "eyJCZXRWYWx1ZSI6W3siQmFpamlhbGV0eXBlIjoxLCJPZGQiOjIsIk1vbmV5IjoxMDAsIk1heG1vbmV5IjoxMDAwMDAsIldlaWdodCI6MTd9LHsiQmFpamlhbGV0eXBlIjoyLCJPZGQiOjIsIk1vbmV5IjoxMDAsIk1heG1vbmV5IjoxMDAwMDAsIldlaWdodCI6MH0seyJCYWlqaWFsZXR5cGUiOjMsIk9kZCI6OSwiTW9uZXkiOjEwMCwiTWF4bW9uZXkiOjEwMDAwMCwiV2VpZ2h0IjowfSx7IkJhaWppYWxldHlwZSI6NCwiT2RkIjo0LCJNb25leSI6MTAwLCJNYXhtb25leSI6MTAwMDAwMDAwMDAsIldlaWdodCI6MH0seyJCYWlqaWFsZXR5cGUiOjUsIk9kZCI6NCwiTW9uZXkiOjEwMCwiTWF4bW9uZXkiOjEwMDAwMDAwMDAwLCJXZWlnaHQiOjB9LHsiQmFpamlhbGV0eXBlIjo2LCJPZGQiOjEyLCJNb25leSI6MTAwLCJNYXhtb25leSI6MTAwMDAwMDAwMDAsIldlaWdodCI6MH0seyJCYWlqaWFsZXR5cGUiOjcsIk9kZCI6MTIsIk1vbmV5IjoxMDAsIk1heG1vbmV5IjoxMDAwMDAwMDAwMCwiV2VpZ2h0IjowfSx7IkJhaWppYWxldHlwZSI6OCwiT2RkIjo1NiwiTW9uZXkiOjEwMCwiTWF4bW9uZXkiOjEwMDAwMDAwMDAwLCJXZWlnaHQiOjEwMDB9LHsiQmFpamlhbGV0eXBlIjo5LCJPZGQiOjU2LCJNb25leSI6MTAwLCJNYXhtb25leSI6MTAwMDAwMDAwMDAsIldlaWdodCI6MTAwMH0seyJCYWlqaWFsZXR5cGUiOjEwLCJPZGQiOjM0LCJNb25leSI6MTAwLCJNYXhtb25leSI6MTAwMDAwMDAwMDAsIldlaWdodCI6MH1dLCJBbGxUeXBlIjpbMSwyLDMsNCw1LDYsNyw4LDksMTBdfQ=="
	bytes, err := base64.StdEncoding.DecodeString(JsonBetValue)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(bytes, &glc.betvalue)

	//json.Unmarshal(s, glc.betvalue)

	//t.Log("syspool betotal => ", glc.syspool, glc.bettotal)
	users, _ := fixtures.GetGameUsers()

	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], tgrs.Room)
	users[1].UserState.OnUserQueue(users[1], tgrs.Room)
	//安排用戶入座

	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	//設定庄及庄額度
	users[0].Score = 103400
	users[0].LoginScore = 103400
	glc.zhuangID = 9100001
	glc.startmoney = users[0].Score

	//player1 下注
	glc.userBet(users[1], fixtures.GetUserBet(ZHUANGSHUN, 10))
	glc.userBet(users[1], fixtures.GetUserBet(ZHUANG, 460))

	glc.kaiJiang()
	glc.sendResult()
	//玩家1輸贏
	if glc.userdata[0].updatemoney != int64(95950) {
		t.Errorf("updatemoney:%v", glc.userdata[0].updatemoney)
	}

	//玩家1輸贏
	if glc.userdata[0].win != int64(142950) {
		t.Errorf("win:%v", glc.userdata[0].win)
	}

	//莊輸贏
	if glc.zhuangWin != int64(-101000) {
		t.Errorf("zhuangWin:%v", glc.zhuangWin)
	}

	glc.gameEnd()
	glc.clearGameData()

	//
	//if bjl.syspool != int64(0) {
	//	t.Errorf("syspool:%v", bjl.syspool)
	//}
	//
	//if bjl.systotalwin != int64(14250) {
	//	t.Errorf("systotalwin:%v", bjl.systotalwin)
	//}
	//
	//if bjl.systaxincome != int64(1000) {
	//	t.Errorf("systaxincome:%v", bjl.systaxincome)
	//}
	//
	////玩家1輸贏
	//if bjl.userdata[0].updatemoney != int64(-20000) {
	//	t.Errorf("updatemoney:%v", bjl.userdata[0].updatemoney)
	//}
	//
	////玩家1輸贏
	//if bjl.userdata[0].win != int64(0) {
	//	t.Errorf("win:%v", bjl.userdata[0].win)
	//}
	//
	////玩家2輸贏
	//if bjl.userdata[1].updatemoney != int64(4750) {
	//	t.Errorf("updatemoney:%v", bjl.userdata[1].updatemoney)
	//}
	//
	////玩家2輸贏
	//if bjl.userdata[1].win != int64(19750) {
	//	t.Errorf("win:%v", bjl.userdata[1].win)
	//}
	//
	////莊輸贏
	//if bjl.zhuangWin != int64(14250) {
	//	t.Errorf("zhuangWin:%v", bjl.zhuangWin)
	//}
	//莊輸贏
	//if bjl.zhuangWin != int64(14250) {
	//	t.Errorf("zhuangWin:%v", bjl.zhuangWin)
	//}
	//bjl.gameEnd()
	//bjl.clearGameData()
}
