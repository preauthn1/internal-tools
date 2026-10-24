package ddz

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strconv"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	ExperienceCashFlag = 6
)

const (
	actionName_JiaoDiZhu = "fapai"  //发牌动画
	actionName_ZhaDan    = "zhadan" //炸弹动画
	actionName_End       = "end"    //结束动画
)

//斗地主
type DDZGameLogic struct {
	*gamelogic.FixedGameLogicBase
	poker             poker.Poker      //游戏扑克牌
	gamePoker         *poker.GamePoker //玩家牌型
	gamestep          gameProgress     //游戏进度
	maxUserCount      int32
	ext               *userExtInfo
	ddt               *ddzTimeOutCfg //阶段超时配置
	oddsCfg           *ddzOddsCfg
	exchangeCfg       *ddzExchangeCardWeight
	noShufflePokerCfg *ddzNoShufflePoker
	actionName_curr   string
	jiaoDiZhuTimer    *event.TimerFunc
	qiangDiZhuTimer   *event.TimerFunc
	mingPaiTimer      *event.TimerFunc
	chuPaiTimer       *event.TimerFunc
	autoChuPai        *event.TimerFunc
	winuser           *userDdzOpt
	tuoGuanUser       map[int32]*robotUserInfoLocal //机器人拆牌结果
	diZhuChuPaiCount  int32                         //地主出牌次數
	jiaoDiZhuRound    int32                         //第幾輪叫地主
	maxJiaoDiZhuRound int32                         //最多叫幾輪
	jiaoDizhuUserRate int32                         //達到n輪沒叫地主時, 真人當地主的機率
	record            *gameRecord
}

//创建抢庄牛牛游戏逻辑对象
func NewDDZGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *DDZGameLogic {
	glc := new(DDZGameLogic)
	glc.FixedGameLogicBase = gamelogic.NewFixedGameLogicBase(t, tfManager, cfgs)
	scfg := cfgs.GetConfig("gameserver")
	glc.ddt = ddzTimeOutInit()
	tmp := scfg.GetConfigInt("gameserver", "maxplayer")
	glc.maxUserCount = int32(tmp)
	glc.maxJiaoDiZhuRound = scfg.GetConfigInt32("gameserver", "maxjiaodizhuround")
	glc.jiaoDizhuUserRate = scfg.GetConfigInt32("gameserver", "jiaodizhuuserrate")
	glc.ext = new(userExtInfo)
	glc.oddsCfg = NewDdzOddsCfg(glc.GetServerConfig().PokerType, cfgs)
	glc.exchangeCfg = NewDdzExchangeCardWeight(cfgs)

	if glc.GetServerConfig().PokerType == gamelogic.PokerTypeShuffle {
		glc.noShufflePokerCfg = NewDdzNoShufflePoker(cfgs)
	}

	//glc.calculation = NewDdzCalculation(cfgs)

	glc.tuoGuanUser = make(map[int32]*robotUserInfoLocal, 0)
	glc.ext.userOpts = make([]*userDdzOpt, glc.maxUserCount)
	glc.ext.turn = 0
	glc.ext.difen = int32(scfg.GetConfigInt("gameserver", "basescore"))
	glc.ext.readytime = int32(scfg.GetConfigInt("gameserver", "readywaittime"))
	glc.ext.timeOut = int32(scfg.GetConfigInt("gameserver", "endactiontimeout"))
	glc.ext.opttime = int32(scfg.GetConfigInt("gameserver", "opttime"))
	glc.ext.maxodd = int32(scfg.GetConfigInt("gameserver", "maxodd"))
	glc.ext.bindaccount = bool(scfg.GetConfigInt("chipaichi", "bindaccount") == 1)
	glc.ext.payminmoney = scfg.GetConfigInt32("chipaichi", "payminmoney")

	glc.gamePoker = poker.NewGamePoker() //创建游戏牌型
	glc.gamestep = gameProgressNone
	glc.ext.baseScore = int32(scfg.GetConfigInt("gameserver", "basescore"))
	glc.ext.lastQiangDiZhu = -1
	//牌局紀錄
	glc.record = new(gameRecord)
	glc.RegisterAfterGameStartFunc(func() {
		glc.ext.init(t)
		//牌局紀錄 設定玩家
		glc.record.setUser(t)
		//glc.broadcastDeskCfg()
		glc.oddsChange(0, "base", netproto.OddsChangeType_null)
		glc.broadcastGameStep(gameProgressBegin) //设置游戏阶段
		glc.broadcastGameStep(gameProgressFapai)
		glc.faPai(true, 17) //发牌
	})
	glc.RegisterAfterGameEndFunc(func() {
		glc.clearGameData()
		glc.jiaoDiZhuRound = 0
	})
	glc.RegisterBeforeGameStartFunc(func() {
		glc.initPai()
		//牌局紀錄 初始化
		glc.record.init()
	})

	glc.RegisterAfterGameEndFunc(func() {
		for _, u := range t.GetTableUser() {
			glc.UserFree(u)
		}
	})
	return glc
}

func (d *DDZGameLogic) getUserOpt(userId int32) *userDdzOpt {
	for _, v := range d.ext.userOpts {
		if v.userId == userId {
			return v
		}
	}
	return nil
}

//判断当前阶段是否所有用户都已操作
func (d *DDZGameLogic) allUserDo() bool {
	for _, v := range d.ext.userOpts {
		if v.isDo == false {
			return false
		}
	}
	return true
}

//重置
func (d *DDZGameLogic) resetDo() {
	for _, v := range d.ext.userOpts {
		v.isDo = false
	}
}

//清除游戏数据
func (d *DDZGameLogic) clearGameData() {
	d.poker = poker.NewRandomPoker(true, false, 1)
	if d.chuPaiTimer != nil {
		d.chuPaiTimer.Stop()
	}
	if d.mingPaiTimer != nil {
		d.mingPaiTimer.Stop()
	}
	if d.jiaoDiZhuTimer != nil {
		d.jiaoDiZhuTimer.Stop()
	}
	if d.autoChuPai != nil {
		d.autoChuPai.Stop()
	}
	if d.qiangDiZhuTimer != nil {
		d.qiangDiZhuTimer.Stop()
	}

	d.winuser = nil
	d.gamePoker = poker.NewGamePoker()
	d.gamestep = gameProgressNone
	d.tuoGuanUser = make(map[int32]*robotUserInfoLocal, 0)
	d.diZhuChuPaiCount = 0
	d.ext.reset()
}

//广播游戏阶段
func (d *DDZGameLogic) broadcastGameStep(p gameProgress) {
	d.gamestep = p
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyGameStep)
	gsdata := new(netproto.DDZ_NotifyGameStep)
	gsdata.Step = proto.Int32(int32(d.gamestep))
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

//广播指定用户进入等待超时
func (d *DDZGameLogic) broadcastTimeOpt(seat int32, timeOut int32) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut)
	gsdata := new(netproto.DDZ_NotifyGameTimeOut)
	gsdata.SeatNo = proto.Int32(seat)
	gsdata.TimeOut = proto.Int32(timeOut)
	gsmsg.MsgData = gsdata
	//mlog.Debug("当前该出牌的玩家座位号%d", d.ext.getNowSeat())
	for _, v := range d.ext.userOpts {
		if v.action != nil {
			if v.action.ct == cardTypeNone {
				mlog.Debug("座位号为%d id=%d的玩家出牌%v isTuoGuan[%v]", v.userSeat, v.userId, v.isTuoGuan, "不出")
			}
			if v.action.cg != nil {
				mlog.Debug("座位号为%d id=%d的玩家出牌%s isTuoGuan[%v]", v.userSeat, v.userId, v.action.cg.String(), v.isTuoGuan)
			}
		}
	}
	d.BoardcastMessage(gsmsg)
}

func (d *DDZGameLogic) broadcastDiZhuInfo(seat int32, cg *poker.CardGroup) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyDiZhuInfo)
	gsdata := new(netproto.DDZ_NotifyDizhuInfo)
	gsdata.Calculation = proto.Int32(int32(d.oddsCfg.calMethod))
	gsdata.SeatNo = proto.Int32(seat)
	gsdata.Cards = new(netproto.CardArry)
	gsdata.Cards.Cards = make([]*netproto.Card, 0)
	for _, v := range cg.Cards {
		c := new(netproto.Card)
		c.Number = &v.CardNumber
		c.Color = &v.Color
		gsdata.Cards.Cards = append(gsdata.Cards.Cards, c)
	}
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

//广播用户叫地主操作
func (d *DDZGameLogic) broadcastJiaoDiZhu(seat int32, isJiaoDiZhu bool) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyJiaoDiZhu)
	gsdata := new(netproto.DDZ_NotifyJiaoDiZhu)
	gsdata.SeatNo = proto.Int32(seat)
	gsdata.IsJiaoDiZhu = proto.Bool(isJiaoDiZhu)
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

//广播用户抢地主操作
func (d *DDZGameLogic) broadcastQiangDiZhu(seat int32, isQiang bool) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyQiangDiZhu)
	gsdata := new(netproto.DDZ_NotifyQiangDiZhu)
	gsdata.SeatNo = proto.Int32(seat)
	gsdata.IsQiang = proto.Bool(isQiang)
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

//广播用户抢地主操作
func (d *DDZGameLogic) broadcastMingPai(u *gamelogic.GameUser, isMingPai bool) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyMingPai)
	gsdata := new(netproto.DDZ_NotifyMingPai)
	gsdata.SeatNo = proto.Int32(u.SeatNo)
	gsdata.IsMing = proto.Bool(isMingPai)
	if isMingPai {
		user := d.gamePoker.GetUserPoker(u.UserID)
		if user != nil {
			gsdata.Mpcs = new(netproto.CardArry)
			for _, v := range user.Cards.Cards {
				card := new(netproto.Card)
				card.Number = &v.CardNumber
				card.Color = &v.Color
				gsdata.Mpcs.Cards = append(gsdata.Mpcs.Cards, card)
			}
		}
	}
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

func (d *DDZGameLogic) sendDeskCfg(user *gamelogic.GameUser) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyDeskCfg)
	gsdata := new(netproto.DDZ_NotifyDeskCfg)
	gsdata.Difen = proto.Int32(d.ext.difen)
	gsdata.TimeOut = proto.Int32(d.ext.readytime)
	gsdata.Odds = proto.Int32(1)
	gsdata.MaxOdd = proto.Int32(d.ext.maxodd)
	gsmsg.MsgData = gsdata
	d.SendMessageToUser(user, gsmsg)
}

//广播用户出牌操作
//userId 用户id
//cg 牌组
//cardNumber 玩家出牌中 主牌的点数
//mainCardCnt 玩家出牌中 xx副牌
func (d *DDZGameLogic) broadcastChuPai(seat int32, cg *poker.CardGroup, ct cardType, cardNumber int32, mainCardCnt int32) {
	uz := d.ext.getExtInfoByTurn(seat)
	upk := d.gamePoker.GetUserPoker(uz.userId)
	tuoGuanUser := d.tuoGuanUser[uz.userId]
	var w bytes.Buffer

	if tuoGuanUser != nil {
		for _, v := range tuoGuanUser.ctg.Ct {
			w.WriteString(fmt.Sprintf(" {cardType[%v] max[%v], min[%v] }; ", v.byteCardType, v.nMax, v.nMin))
		}
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("廣播出牌, userId[%v] cg[%v] ct[%v] cardNumber[%v] 剩下手牌[%v] 牌型[%v]",
		uz.userId, cg, ct, cardNumber, upk.Cards.Cards, w.String())).Debug("託管")

	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyChuPai)
	gsdata := new(netproto.DDZ_NotifyChuPai)
	gsdata.SeatNo = proto.Int32(seat)
	gsdata.Cards = new(netproto.CardArry)
	gsdata.Cards.Cards = make([]*netproto.Card, 0)
	var odds int32 = 0
	if ct == cardTypeRocket || ct == cardTypeBomb {
		odds = int32(d.oddsCfg.odds[oddsZhaDan])
	}
	gsdata.CurOdds = proto.Int32(odds)
	for _, c := range cg.Cards {
		card := new(netproto.Card)
		card.Color = &c.Color
		card.Number = &c.CardNumber
		gsdata.Cards.Cards = append(gsdata.Cards.Cards, card)
	}
	gsdata.MainCards = &mainCardCnt
	gsdata.MainCardsNumber = &cardNumber
	tmp := int32(ct)
	gsdata.CardType = &tmp
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)

	//gameuser := d.GetUserByUserID(uz.userId)
	//d.SendChiPaiChiStatus(gameuser)
	//d.SendChiPaiChiInfo(gameuser)
}

//广播房间信息
func (d *DDZGameLogic) broadcastRoomInfo(odds int32, ct netproto.OddsChangeType) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyRoomInfo)
	gsdata := new(netproto.DDZ_NotifyRoomInfo)
	gsdata.Odds = proto.Int32(odds)
	gsdata.NowOdds = proto.Int32(d.ext.odds)
	c := int32(ct)
	gsdata.Oct = &c
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

func (d *DDZGameLogic) broadcastSettlement(u *userDdzOpt) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifySettlement)
	gsdata := new(netproto.DDZ_NotifySettlement)
	m := d.ext.or.CaculateOdds(d.oddsCfg.calMethod)
	gsdata.BaseScore = proto.Int32(d.ext.baseScore)
	gsdata.JiaoDiZhuOdds = proto.Int32(m[DdzOddsTypeJdz])
	gsdata.QiangDiZhuOdds = proto.Int32(m[DdzOddsTypeQdz])
	gsdata.DiZhuCardTypeOdds = proto.Int32(m[DdzOddsTypeDiZhuPai])
	gsdata.ZhaDanOdds = proto.Int32(m[DdzOddsTypeZhaDan])
	gsdata.ChunTianOdds = proto.Int32(m[DdzOddsTypeChunTian])
	gsdata.FanChunTianOdds = proto.Int32(m[DdzOddsTypeFanChunTian])
	gsdata.MingPaiOdds = proto.Int32(m[DdzOddsTypeMingPai])
	gsdata.Rocket = proto.Int32(m[DdzOddsTypeRocket])
	gsdata.TotalOdds = proto.Int32(d.ext.odds)
	gsdata.WinSeat = proto.Int32(u.userSeat)
	gsdata.Calculation = proto.Int32(int32(d.oddsCfg.calMethod))
	for _, v := range d.ext.userOpts {
		s := new(netproto.DDZ_ResultRecord)
		s.UserId = proto.Int32(v.userId)
		s.Money = proto.Int32(v.money)
		s.IsDiZhu = proto.Bool(v.isDiZhu)
		user := d.GetUserByUserID(v.userId)
		s.NickName = proto.String(user.NickName)
		if v.lt == loseEnumBankruptcy {
			s.Bankruptcy = proto.Bool(true)
		} else {
			s.Bankruptcy = proto.Bool(false)
		}
		s.Cards = new(netproto.CardArry)
		upk := d.gamePoker.GetUserPoker(v.userId)
		for _, cc := range upk.Cards.Cards {
			c := new(netproto.Card)
			c.Number = &cc.CardNumber
			c.Color = &cc.Color
			s.Cards.Cards = append(s.Cards.Cards, c)
		}
		gsdata.Settlement = append(gsdata.Settlement, s)
	}
	gsmsg.MsgData = gsdata
	d.BoardcastMessage(gsmsg)
}

func (d *DDZGameLogic) initPai() {
	//随机发牌
	if d.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		d.poker = poker.NewRandomPoker(true, false, 1)
		d.poker = poker.NewFixedPoker(d.AllTableUser(), d.poker.(*poker.RandomPoker), 17)
	} else if d.GetServerConfig().PokerType == gamelogic.PokerTypeShuffle {
		d.poker = poker.NewNoShufflePoker(true, false, 1, d.noShufflePokerCfg.groupPoker)
		d.poker = poker.NewFixedNoShufflePoker(d.AllTableUser(), d.poker.(*poker.NoShufflePoker), 17, d.noShufflePokerCfg.alloctake)
	} else { //发牌器发牌
		d.poker = poker.NewFilePoker(d.Configs, d.AllTableUser())
		d.poker = poker.NewFixedFilePoker(d.poker.(*poker.FilePoker))
		d.poker.(*poker.FixedPoker).FixedFileAlloc(d.AllTableUser())
		//for _, user := range d.AllTableUser() {
		//	d.poker.Deal(user.UserID, 17)
		//}
	}
}

//发牌
func (d *DDZGameLogic) faPai(shuffle bool, count int) {
	if shuffle {
		d.initPai()
		//通知機器人要發牌了
		for _, ur := range d.AllGameUser() {
			if ur.IsRobot() {
				//mlog.Debug("通知機器人要發牌了")
				d.beforeFapai(ur.UserID)
			}
		}
	}

	//给所有用户发牌,先发13张
	for _, ur := range d.AllGameUser() {
		cards := d.poker.Deal(ur.UserID, count)
		upk := poker.NewUserPoker(ur, cards)
		d.gamePoker.AddUserPoker(upk)
		//通知用户牌型
		gsmsg := new(msg.Message)
		gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
		gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyFaPai)
		gsdata := new(netproto.DDZ_NotifyFaPai)
		gsdata.Cards = new(netproto.CardArry)
		gsdata.Cards.Cards = make([]*netproto.Card, 0)
		for _, c := range cards.Cards {
			card := new(netproto.Card)
			card.Color = proto.Int32(c.Color)
			card.Number = proto.Int32(c.CardNumber)
			gsdata.Cards.Cards = append(gsdata.Cards.Cards, card)
		}
		tmp := sortCardGroup(*upk.Cards)
		sort.Sort(&tmp)
		d.GameLog.Fapai(ur.UserID, fmt.Sprintf("%s", tmp.Cards))
		gsdata.SeatNo = proto.Int32(ur.SeatNo)
		gsmsg.MsgData = gsdata
		d.SendMessageToUser(ur, gsmsg)
		d.ext.turn++
	}
	//启动确定庄家流程
	//d.actionName_curr = d.PlayerActionManager.NewAction("fapai", d.GetGameReadyUserCount(), 3, d.jiaoDiZhuStart)
	d.actionName_curr = d.RegisterAnimation(actionName_JiaoDiZhu, 3, d.jiaoDiZhuStart)
}

//房间倍率发生变化
func (d *DDZGameLogic) oddsChange(userId int32, oddsName string, ot netproto.OddsChangeType) {
	r := new(record)
	r.odds = int32(d.oddsCfg.odds[oddsName])
	r.ot = ot
	r.userId = userId

	if oddsName == "base" {
		d.ext.odds = r.odds
		d.ext.or.push(r)
		d.GameLog.AddGameStep("%s_%d 倍率发生变化:(%v)", oddsName, ot, d.ext.odds)
		d.broadcastRoomInfo(r.odds, r.ot)
		return
	}

	if d.oddsCfg.calMethod == 0 { //加法計算
		if r.odds > 0 {
			d.ext.odds += r.odds
			d.ext.or.push(r)
		} else {
			return
		}
	} else { //乘法計算
		if r.odds > 1 {
			d.ext.or.push(r)
			d.ext.odds *= r.odds
		} else {
			return
		}
	}

	if d.ext.odds > d.ext.maxodd {
		d.ext.odds = d.ext.maxodd
	}
	//mlog.Debug("倍率发生变化:%s %d", oddsName, d.ext.odds)
	d.GameLog.AddGameStep("%s_%d 倍率发生变化:(%v)", oddsName, ot, d.ext.odds)
	d.broadcastRoomInfo(r.odds, r.ot)
}

func (d *DDZGameLogic) caculateResult(u *userDdzOpt) {
	var winActMoney int32 = 0
	u.isWin = true

	uscoresRet := make(map[int32]int64)
	flowsRet := make(map[int32]int64)
	cscoresRet := make(map[int32]int64)
	if u.isDiZhu { //地主贏
		for _, v := range d.ext.userOpts {
			if v.userId != u.userId {
				v.isWin = false
			} else {
				continue
			}
			v.money = d.ext.odds * d.ext.baseScore
			lm := d.GetUserByUserID(v.userId)
			if lm.Score <= int64(v.money) {
				v.lt = loseEnumBankruptcy //破产
				winActMoney += int32(lm.Score)
				v.money = int32(lm.Score)
			} else {
				winActMoney += v.money
			}
			flowsRet[v.userId] = int64(v.money)
			v.money = 0 - v.money
			amount, tax := d.GetUserWinAndTax(int64(v.money))
			//d.UpdateUserScoreFunc(v.userId, amount, tax)
			d.updateUserScoreOrExperienceCash(v.userId, amount, tax)
			v.money = int32(amount)
			uscoresRet[v.userId] = amount
			d.GameLog.Context().UserUpdateScore(v.userId, amount)
			d.GameLog.Context().UserTaxAmount(v.userId, tax)

			cscoresRet[v.userId] = 0
		}
		u.money = winActMoney
		flowsRet[u.userId] = int64(u.money)
		//mlog.Debug("计算输赢信息:")
		amount, tax := d.GetUserWinAndTax(int64(u.money))
		mlog.Debug("钱：%v，税收：%v", amount, tax)
		//d.UpdateUserScoreFunc(u.userId, amount, tax)
		d.updateUserScoreOrExperienceCash(u.userId, amount, tax)
		uscoresRet[u.userId] = amount
		d.GameLog.Context().UserUpdateScore(u.userId, amount)
		d.GameLog.Context().UserTaxAmount(u.userId, tax)
		cscoresRet[u.userId] = 0
		u.money = int32(amount)
	} else { //地主輸
		var dz *gamelogic.GameUser = nil
		var dzext *userDdzOpt = nil
		for _, v := range d.ext.userOpts {
			if !v.isDiZhu {
				v.isWin = true
				v.money = d.ext.odds * d.ext.baseScore
				winActMoney += v.money
			} else {
				dz = d.GetUserByUserID(v.userId)
				dzext = v
				v.isWin = false
			}
		}
		if dz.Score <= int64(winActMoney) {
			winActMoney = int32(dz.Score)
			dzext.lt = loseEnumBankruptcy //破产
			dzext.money = 0 - winActMoney
			amount, tax := d.GetUserWinAndTax(int64(dzext.money))
			//d.UpdateUserScoreFunc(dz.UserID, amount, tax)
			d.updateUserScoreOrExperienceCash(dz.UserID, amount, tax)
			uscoresRet[dz.UserID] = amount
			d.GameLog.Context().UserUpdateScore(dz.UserID, amount)
			d.GameLog.Context().UserTaxAmount(dz.UserID, tax)
			cscoresRet[dz.UserID] = 0
			dzext.money = int32(amount)
		}
		for _, v := range d.ext.userOpts {
			if v.isDiZhu {
				flowsRet[v.userId] = int64(winActMoney)
				v.money = 0 - winActMoney
				amount, tax := d.GetUserWinAndTax(int64(v.money))
				//d.UpdateUserScoreFunc(v.userId, amount, tax)
				d.updateUserScoreOrExperienceCash(v.userId, amount, tax)
				uscoresRet[v.userId] = amount
				d.GameLog.Context().UserUpdateScore(v.userId, amount)
				d.GameLog.Context().UserTaxAmount(v.userId, tax)
				cscoresRet[v.userId] = 0
				continue
			}
			v.money = winActMoney / 2
			flowsRet[v.userId] = int64(v.money)
			amount, tax := d.GetUserWinAndTax(int64(v.money))
			//d.UpdateUserScoreFunc(v.userId, amount, tax)
			d.updateUserScoreOrExperienceCash(v.userId, amount, tax)
			uscoresRet[v.userId] = amount
			d.GameLog.Context().UserUpdateScore(v.userId, amount)
			d.GameLog.Context().UserTaxAmount(v.userId, tax)
			cscoresRet[v.userId] = 0
			v.money = int32(amount)
		}
	}
	d.GameLog.AddGameStep("計算流水: %v", flowsRet)
	d.GameLog.GameEnd()
	//牌局紀錄 結算
	d.record.checkResult(u.isDiZhu, u.userSeat, d.ext.odds)
	d.broadcastSettlement(u)
}

func (d *DDZGameLogic) updateUserScoreOrExperienceCash(userId int32, amount int64, tax int64) {
	ur := d.GetUserByUserID(userId)
	flag, _ := strconv.ParseInt(d.TableManager.Room.Flag, 10, 32)
	if flag == ExperienceCashFlag && !ur.IsRobot() {
		ur.ExperienceCash += amount
		d.UpdateExperienceScoreFunc(userId, amount, tax)
		ur.Score = ur.ExperienceCash
		mlog.Debug("體驗更新金幣 ud.money[%v] Score[%v] ExperienceCash[%v]",
			amount, ur.Score, ur.ExperienceCash)
	} else {
		var flowAmount, betAmount, payAmount int64
		flowAmount = util.ABSInt64(amount)      //流水=派彩
		betAmount = util.ABSInt64(amount) + tax //有效投注=輸贏+稅
		payAmount = betAmount + amount          //返邏金額=有效投注+派彩
		ur.FlowAmount += flowAmount
		ur.BetAmount += betAmount
		ur.PayAmount += payAmount
		d.GameLog.AddGameStep("用戶[%d] 結算金額, win %d, 稅收 %d, 流水 %d, 有效投注 %d, 返還金額 %d", ur.UserID, amount, tax, flowAmount, betAmount, payAmount)
		d.UpdateUserScoreFunc(userId, amount, tax)
	}
}

func (d *DDZGameLogic) sendSettlement() {
	if d.winuser != nil {
		d.caculateResult(d.winuser)
		d.actionName_curr = d.RegisterAnimation(actionName_End, d.ext.timeOut, d.test)
	}
}

func (d *DDZGameLogic) test() {
	mlog.Debug("game end!!!!!!!!!!!!!!")
	d.GameEndFunc()
}

func (d *DDZGameLogic) settlement(winUserExtInfo *userDdzOpt) {
	d.broadcastGameStep(gameProgressGameEnd)
	d.winuser = winUserExtInfo
	d.TimerFuncManager.NewNamedTimer(1, d.sendSettlement, "")
}

func (d *DDZGameLogic) show() {
	for _, v := range d.AllTableUser() {
		mlog.Debug("=====playerId=%d,seat=%d", v.UserID, v.SeatNo)
		mlog.Debug("=====%v", d.gamePoker.GetUserPoker(v.UserID).Cards.Cards)
	}
}

//检测地主牌牌型
func (d *DDZGameLogic) checkDiZhuCardType(cg *poker.CardGroup) (string, netproto.OddsChangeType) {
	l := len(cg.Cards)
	if l < 3 {
		return "", netproto.OddsChangeType_null
	}
	tmp := sortCardGroup(*cg)
	sort.Sort(&tmp)
	if cg.Cards[0].IsSameCardNumber(cg.Cards[1]) && cg.Cards[1].IsSameCardNumber(cg.Cards[2]) {
		return oddsBaoZi, netproto.OddsChangeType_baozi
	}
	if cg.Cards[0].IsBigKing() || cg.Cards[1].IsBigKing() || cg.Cards[2].IsBigKing() {
		return oddsDaWang, netproto.OddsChangeType_dawang
	}
	if cg.Cards[0].IsSmallKing() || cg.Cards[1].IsSmallKing() || cg.Cards[2].IsSmallKing() {
		return oddsXiaoWang, netproto.OddsChangeType_xiaowang
	}
	if cg.Cards[0].Color == cg.Cards[1].Color && cg.Cards[1].Color == cg.Cards[2].Color {
		return oddsTongHua, netproto.OddsChangeType_tonghua
	}
	if math.Abs(float64(cg.Cards[0].CardNumber-cg.Cards[1].CardNumber)) == 1 && math.Abs(float64(cg.Cards[1].CardNumber-cg.Cards[2].CardNumber)) == 1 {
		return oddsShunZi, netproto.OddsChangeType_shunzi
	}
	if cg.Cards[0].IsSameCardNumber(cg.Cards[1]) || cg.Cards[0].IsSameCardNumber(cg.Cards[2]) || cg.Cards[1].IsSameCardNumber(cg.Cards[2]) {
		return oddsDuiZi, netproto.OddsChangeType_duizi
	}
	return "sanpai", netproto.OddsChangeType_null
}

//準備發牌了
func (d *DDZGameLogic) beforeFapai(userId int32) {
	ur := d.GetUserByUserID(userId)
	if ur == nil {
		mlog.Error("系统错误,準備發牌時用戶信息異常")
		return
	}
	if !ur.IsRobot() {
		mlog.Error("系统错误,用戶不是機器人不用通知準備要發牌了")
		return
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_BeforeFapai)
	pmsg.MsgData = nil
	d.SendMessageToUserID(userId, pmsg)
}

//发最后三张牌给地主
func (d *DDZGameLogic) faPaiLastCards(userId int32) {
	ur := d.GetUserByUserID(userId)
	if ur == nil {
		mlog.Error("系统错误,发地主牌时获取地主信息异常")
		return
	}
	//牌局紀錄 設定地主
	d.record.setDiZhu(userId)

	//牌局紀錄 設定初始牌
	d.record.setInitCards(d.gamePoker)

	cards := d.poker.Deal(poker.RestCardsUserID, 3)
	upk := poker.NewUserPoker(ur, cards)
	d.gamePoker.AddUserPoker(upk)
	str, ct := d.checkDiZhuCardType(cards)
	d.oddsChange(userId, str, ct)
	d.ext.dizhuCards.Cards = append(d.ext.dizhuCards.Cards, cards.Cards...)

	//牌局紀錄 設定最後三張牌
	d.record.setLastCards(cards.Cards)

	//通知用户牌型
	d.GameLog.Fapai(userId, cards.String())
	d.broadcastDiZhuInfo(ur.SeatNo, cards)
}

func (d *DDZGameLogic) passJiaoDiZhu() {
	user := d.GetUserBySeatNo(d.ext.getNowSeat())
	dzindex, _ := d.ext.getLastDiZhu()
	if (d.ext.turn+1)%d.maxUserCount == int32(dzindex) { //没有人叫地主 从新发牌 进入新一轮游戏
		if d.maxJiaoDiZhuRound != 0 && d.jiaoDiZhuRound >= d.maxJiaoDiZhuRound { //超過直接指定第一個為地主
			d.randomDizZhu()
			d.GameLog.AddGameStep("第 %d 輪沒人叫地主，指定第 %d 座位為地主", d.jiaoDiZhuRound, d.ext.lastQiangDiZhu)
			d.makesureDiZhu()
			return
		}
		d.GameLog.AddGameStep("用户%d不叫地主", user.UserID)
		d.GameLog.AddGameStep("没有人叫地主 从新发牌 进入新一轮游戏")
		d.ext.turn = d.ext.turn % d.maxUserCount
		d.clearGameData()
		for _, ur := range d.AllGameUser() {
			if ur.IsOffLine() {
				//有游戏玩家离线，直接结束游戏
				d.test()
				return
			}
		}
		d.oddsChange(0, "base", netproto.OddsChangeType_null)
		d.broadcastGameStep(gameProgressFapai)
		d.faPai(true, 17)
		return
	}

	d.GameLog.AddGameStep("用户%d不叫地主", user.UserID)
	d.broadcastJiaoDiZhu(user.SeatNo, false) //不叫
	uop := d.ext.getExtInfoByTurn(d.ext.turn)
	uop.JiaoDiZhu = JDZBuJiao
	d.ext.turn++
	d.broadcastTimeOpt(d.ext.getNowSeat(), d.ddt.GetTime(gameProgressJiaoDiZhu))
	d.jiaoDiZhuTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressJiaoDiZhu), d.jiaoDiZhuTimeOut, "")
}

//達到n輪沒當地主時, 真人當地主
func (d *DDZGameLogic) realUserDizZhu() {
	rndJiaoDizhuUser := util.RndInt(0, 100)
	if rndJiaoDizhuUser < int(d.jiaoDizhuUserRate) {
		for _, u := range d.AllTableUser() {
			if !u.IsRobot() {
				d.ext.lastQiangDiZhu = u.SeatNo
				break
			}
		}
	}
}

//達到n輪沒當地主時, 隨機一個當地主，不只針對真實用，
func (d *DDZGameLogic) randomDizZhu() {
	rndJiaoDizhuUser := util.RndInt(0, 100)
	if rndJiaoDizhuUser < int(d.jiaoDizhuUserRate) {
		var rndUserIndex = util.RndInt(0, len(d.AllTableUser()))
		for i, u := range d.AllTableUser() {
			if i == rndUserIndex {
				d.ext.lastQiangDiZhu = u.SeatNo
				break
			}
		}
	}
}

func (d *DDZGameLogic) jiaoDiZhuStart() {
	d.jiaoDiZhuRound++
	//启动叫牌流程
	var seat int32 = 0
	_, v := d.ext.getLastDiZhu()
	if d.ext.curTurn <= 1 && v == nil { //第一把 随机叫牌
		rand.Seed(time.Now().UnixNano())
		i := rand.Int31() % int32(len(d.ext.userOpts))
		d.ext.userOpts[i].JiaoDiZhu = JDZJiao
		d.ext.turn = i
		d.ext.userOpts[i].isLastDiZhu = true
		seat = d.ext.userOpts[i].userSeat
		d.ext.lastQiangDiZhu = i
	} else {
		d.ext.resetLastDiZhu()
		seat = int32(v.userSeat)
		d.ext.turn = seat
	}
	d.broadcastGameStep(gameProgressJiaoDiZhu)
	d.broadcastTimeOpt(seat, d.ddt.GetTime(gameProgressJiaoDiZhu))
	d.jiaoDiZhuTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressJiaoDiZhu), d.jiaoDiZhuTimeOut, "")
}

//叫牌阶段超时 做出相应处理进入下一阶段
func (d *DDZGameLogic) jiaoDiZhuTimeOut() {
	//启动叫牌流程
	if d.jiaoDiZhuTimer != nil {
		d.jiaoDiZhuTimer.Stop()
	}
	d.passJiaoDiZhu()
}

func (d *DDZGameLogic) qiangZhuangStart() {
	//启动抢地主流程
	if d.qiangDiZhuTimer != nil {
		d.qiangDiZhuTimer.Stop()
	}
	var xt *userDdzOpt
	for {
		xt = d.ext.getExtInfoByTurn(d.ext.turn)
		if xt.JiaoDiZhu == JDZBuJiao {
			d.broadcastQiangDiZhu(xt.userSeat, false)
			d.ext.turn++
			continue
		} else {
			break
		}
	}
	if d.checkBuqiangCount() == 2 {
		d.makesureDiZhu()
	} else {
		d.broadcastTimeOpt(xt.userSeat, d.ddt.GetTime(gameProgressRob))
		d.qiangDiZhuTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressRob), d.qiangDiZhuTimeOut, "")
	}
}

func (d *DDZGameLogic) makesureDiZhu() {
	dz := d.ext.getExtInfoByTurn(d.ext.lastQiangDiZhu)
	if dz == nil {
		mlog.Error("系统错误,获取地主信息失败")
		d.GameEndFunc()
		return
	}
	urp := d.gamePoker.GetUserPoker(dz.userId)
	if urp == nil {
		mlog.Error("系统错误,获取抢地主后最终地主人选失败")
		d.GameEndFunc()
		return
	}
	d.ext.resetLastDiZhu()
	d.faPaiLastCards(dz.userId)
	dz.isDiZhu = true
	dz.isLastDiZhu = true
	d.ext.lastQiangDiZhu = dz.userSeat
	d.ext.turn = d.ext.lastQiangDiZhu
	d.broadcastGameStep(gameProgressShow)
	d.broadcastTimeOpt(dz.userSeat, d.ddt.GetTime(gameProgressShow))
	d.mingPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressShow), d.mingPaiTimeOut, "")
}

//抢地主阶段超时 做出相应处理进入下一阶段
func (d *DDZGameLogic) qiangDiZhuTimeOut() {
	idx, jd := d.ext.getJiaoDiZhu()
	if jd == nil {
		mlog.Error("系统错误,获取叫地主用户信息失败")
		return
	}
	xt := d.ext.getExtInfoByTurn(d.ext.turn)
	xt.isQiangDizhu = 1
	d.broadcastQiangDiZhu(xt.userSeat, false)
	if (d.ext.turn)%d.maxUserCount == int32(idx) || d.checkBuqiangCount() == 2 { //确定地主后进入明牌阶段
		d.makesureDiZhu()
		return
	}
	d.ext.turn++
	d.broadcastTimeOpt(d.ext.getNowSeat(), d.ddt.GetTime(gameProgressRob))
	d.qiangDiZhuTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressRob), d.qiangDiZhuTimeOut, "")
}

//明牌超时处理
func (d *DDZGameLogic) mingPaiTimeOut() {
	mlog.Debug("託管 明牌")
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	if d.mingPaiTimer != nil {
		d.mingPaiTimer.Stop()
	}
	user := d.GetUserByUserID(u.userId)
	d.broadcastMingPai(user, false)
	d.broadcastGameStep(gameProgressStart)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("明牌 broadcastTimeOpt u.userId[%v], d.ext.opttime[%v]", u.userId, d.ext.opttime)).Debug("託管")
	d.broadcastTimeOpt(u.userSeat, d.ext.opttime)
	d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
}

type sortCardGroup poker.CardGroup

func (scg *sortCardGroup) Len() int {
	return len(scg.Cards)
}

func (scg *sortCardGroup) Swap(i, j int) {
	scg.Cards[i], scg.Cards[j] = scg.Cards[j], scg.Cards[i]
}

func (scg *sortCardGroup) Less(i, j int) bool {
	return cardMap[scg.Cards[i].CardNumber] < cardMap[scg.Cards[j].CardNumber]
}

//等待出牌超时
func (d *DDZGameLogic) chupaiTimeOut() {
	d.AutoChuPai()
	mlog.Debug("chupaiTimeOut exit.....")
}

//叫地主响应
func (d *DDZGameLogic) responseUserJiaoDiZhu(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_JiaoDiZhuRet)
	gsdata := new(netproto.DDZ_JiaoDiZhuRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

//抢地主响应
func (d *DDZGameLogic) responseUserQiangDiZhu(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_QiangDiZhuRet)
	gsdata := new(netproto.DDZ_QiangDiZhuRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

//明牌主响应
func (d *DDZGameLogic) responseUserMingPai(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_MingPaiRet)
	gsdata := new(netproto.DDZ_MingPaiRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

//出牌响应
func (d *DDZGameLogic) responseUserChuPai(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_ChuPaiRet)
	gsdata := new(netproto.DDZ_ChuPaiRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

func (d *DDZGameLogic) responseUserPass(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_PassRet)
	gsdata := new(netproto.DDZ_PassRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

func (d *DDZGameLogic) responseTuoGuan(userId int32, errocode netproto.DDZ_ErrorCode, msgstr string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_TuoGuanRet)
	gsdata := new(netproto.DDZ_TuoGuanRet)
	gsdata.Code = &errocode
	gsdata.Msg = proto.String(msgstr)
	pmsg.MsgData = gsdata
	d.SendMessageToUserID(userId, pmsg)
}

func (d *DDZGameLogic) broadcastTuoGuanState(seat int32, state bool) {
	uz := d.ext.getExtInfoByTurn(seat)
	if uz == nil {
		mlog.Error("Got user exit info faild")
		return
	}
	uz.isTuoGuan = state
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_NotifyTuoGuan)
	gsdata := new(netproto.DDZ_NotifyTuoGuan)
	gsdata.UserSeat = proto.Int32(uz.userSeat)
	gsdata.State = proto.Bool(state)
	pmsg.MsgData = gsdata
	d.BoardcastMessage(pmsg)
}

func (d *DDZGameLogic) sendGameInfoToUser(userId int32) {
	uz := d.ext.getExtInfoByTurn(userId)
	if uz == nil {
		mlog.Error("Got user exit info faild")
		return
	}
	rUser := d.GetUserByUserID(userId)
	if rUser == nil {
		mlog.Error("Got user info faild")
		return
	}
	var dizhuseat int32 = 0
	_, dz := d.ext.getDiZhu()
	if dz == nil {
		dizhuseat = -1
	} else {
		dizhuseat = dz.userSeat
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_UserRecoverRet)
	gsdata := new(netproto.DDZ_UserRecoverRet)
	gsdata.UserId = proto.Int32(userId)
	gsdata.GameStep = proto.Int32(int32(d.gamestep))
	gsdata.Odds = proto.Int32(d.ext.odds)
	gsdata.OptSeat = proto.Int32(d.ext.getNowSeat())
	if len(d.ext.his.actions) > 0 {
		gsdata.IsSendCard = proto.Bool(true)
	} else {
		gsdata.IsSendCard = proto.Bool(false)
	}
	if d.gamestep == gameProgressJiaoDiZhu && d.jiaoDiZhuTimer != nil {
		gsdata.OptLastTime = proto.Int32(d.ddt.GetTime(gameProgressJiaoDiZhu) - d.jiaoDiZhuTimer.ElapsedSeconds())
	} else if d.gamestep == gameProgressRob && d.qiangDiZhuTimer != nil {
		gsdata.OptLastTime = proto.Int32(d.ddt.GetTime(gameProgressRob) - d.qiangDiZhuTimer.ElapsedSeconds())
	} else if d.gamestep == gameProgressShow && d.mingPaiTimer != nil {
		gsdata.OptLastTime = proto.Int32(d.ddt.GetTime(gameProgressShow) - d.mingPaiTimer.ElapsedSeconds())
	} else if d.gamestep == gameProgressStart && d.chuPaiTimer != nil {
		gsdata.OptLastTime = proto.Int32(d.ext.opttime - d.chuPaiTimer.ElapsedSeconds())
	} else {
		gsdata.OptLastTime = proto.Int32(0)
	}
	gsdata.DiZhuSeat = proto.Int32(dizhuseat)
	gsdata.Opi = make([]*netproto.DDZ_PlayerInfo, 0)
	gsdata.Sci = make([]*netproto.DDZ_SendCardInfo, 0)
	for _, v := range d.ext.userOpts {
		user := d.GetUserByUserID(v.userId)
		if user == nil {
			continue
		}
		oui := new(netproto.DDZ_PlayerInfo)
		oui.UserId = proto.Int32(user.UserID)
		oui.SeatNo = proto.Int32(user.SeatNo)
		oui.IsDiZhu = proto.Bool(v.isDiZhu)
		oui.IsMingPai = proto.Bool(v.isMingPai)
		oui.IsTuoGuan = proto.Bool(v.isTuoGuan)
		if v.JiaoDiZhu == JDZJiao {
			oui.IsJiaoDiZhu = netproto.UserActionDo_UADO.Enum()
		} else if v.JiaoDiZhu == JDZBuJiao {
			oui.IsJiaoDiZhu = netproto.UserActionDo_UANOTDO.Enum()
		} else {
			oui.IsJiaoDiZhu = netproto.UserActionDo_UANONE.Enum()
		}
		if v.isQiangDizhu == 1 {
			oui.IsQiangDiZhu = netproto.UserActionDo_UANOTDO.Enum()
		} else if v.isQiangDizhu == 2 {
			oui.IsQiangDiZhu = netproto.UserActionDo_UADO.Enum()
		} else {
			oui.IsQiangDiZhu = netproto.UserActionDo_UANONE.Enum()
		}
		upk := d.gamePoker.GetUserPoker(v.userId)
		if upk != nil {
			oui.CardsCount = proto.Int32(int32(upk.Cards.Len()))
		}
		if v.isMingPai || v.userId == userId {
			oui.Cards = new(netproto.CardArry)
			for _, c := range upk.Cards.Cards {
				cc := new(netproto.Card)
				cc.Color = proto.Int32(c.Color)
				cc.Number = proto.Int32(c.CardNumber)
				oui.Cards.Cards = append(oui.Cards.Cards, cc)
			}
		}
		si := new(netproto.DDZ_SendCardInfo)
		si.SeatNo = proto.Int32(v.userSeat)
		if v.action == nil {
			si.SendType = proto.Int32(2)
		} else {
			if v.action.ct == cardTypeNone {
				si.SendType = proto.Int32(0)
			} else {
				si.SendType = proto.Int32(1)
				si.SendOutCards = new(netproto.CardArry)
				for _, c := range v.action.cg.Cards {
					cc := new(netproto.Card)
					cc.Color = proto.Int32(c.Color)
					cc.Number = proto.Int32(c.CardNumber)
					si.SendOutCards.Cards = append(si.SendOutCards.Cards, cc)
				}
			}
		}
		gsdata.Opi = append(gsdata.Opi, oui)
		gsdata.Sci = append(gsdata.Sci, si)
	}
	gsdata.DiZhuCards = new(netproto.CardArry)
	for _, v := range d.ext.dizhuCards.Cards {
		c := new(netproto.Card)
		c.Color = proto.Int32(v.Color)
		c.Number = proto.Int32(v.CardNumber)
		gsdata.DiZhuCards.Cards = append(gsdata.DiZhuCards.Cards, c)
	}
	_, t := d.checkDiZhuCardType(d.ext.dizhuCards)
	gsdata.DiZhuCardType = proto.Int32(int32(t))
	gsdata.Calculation = proto.Int32(int32(d.oddsCfg.calMethod))
	pmsg.MsgData = gsdata
	d.SendMessageToUser(rUser, pmsg)
}

//叫地主
func (d *DDZGameLogic) onUserJiaoDiZhu(user *gamelogic.GameUser, data []byte) {
	//检查游戏阶段
	if d.gamestep != gameProgressJiaoDiZhu {
		mlog.Error("非法操作，当前不是叫地主阶段(%d)", d.gamestep)
		d.responseUserJiaoDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前不是叫地主阶段")
		return
	}
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	if u == nil {
		mlog.Error("系统错误，获取用户扩展信息失败(%d)", d.gamestep)
		d.responseUserJiaoDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取用户信息失败")
		return
	}
	if user.SeatNo != u.userSeat {
		mlog.Error("非法操作，当前玩游戏的并非同一个人(step=%d,turn=%d,userid=%d,userseat=%d)", d.gamestep, d.ext.turn, u.userId, u.userSeat)
		d.responseUserJiaoDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前玩游戏的并非同一个人")
		return
	}
	ust := &netproto.DDZ_JiaoDiZhu{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v抢庄数据结构错误%v", user, err)
		d.responseUserJiaoDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
	if d.jiaoDiZhuTimer != nil {
		d.jiaoDiZhuTimer.Stop()
	}

	d.responseUserJiaoDiZhu(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
	if !ust.GetIsJiaoDiZhu() {
		d.passJiaoDiZhu()
	} else {
		d.GameLog.AddGameStep("用户%d叫地主", user.UserID)
		d.oddsChange(user.UserID, oddsJiaoDiZhu, netproto.OddsChangeType_jiaodizhu)
		d.ext.lastQiangDiZhu = user.SeatNo //叫
		d.broadcastJiaoDiZhu(user.SeatNo, true)
		u.JiaoDiZhu = JDZJiao
		d.ext.turn++ //下一个玩家
		//进入抢地主阶段
		d.broadcastGameStep(gameProgressRob)
		mlog.Debug("进入抢地主流程...")
		d.qiangZhuangStart()
	}
}

func (d *DDZGameLogic) checkBuqiangCount() int32 {
	var bqc int32 = 0
	for _, v := range d.ext.userOpts {
		if v.isQiangDizhu == 1 || v.JiaoDiZhu == JDZBuJiao {
			bqc++
		}
	}
	return bqc
}

//抢地主
func (d *DDZGameLogic) onQiangDiZhu(user *gamelogic.GameUser, data []byte) {
	if d.gamestep != gameProgressRob {
		mlog.Error("非法操作，当前不是抢地主阶段(%d)", d.gamestep)
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前不是抢地主阶段")
		return
	}
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	if u == nil {
		mlog.Error("系统错误，获取用户扩展信息失败(%d)", d.gamestep)
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取用户信息失败")
		return
	}
	if user.SeatNo != u.userSeat {
		mlog.Error("非法操作，当前玩游戏的并非同一个人(%d) trun=%d,userseat=%d", d.gamestep, d.ext.turn, user.SeatNo)
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前玩游戏的并非同一个人")
		return
	}
	if u.JiaoDiZhu == JDZBuJiao {
		mlog.Error("没有叫地主的玩家不能抢地主 %v", d.gamestep)
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，没有叫地主的玩家不能抢地主")
		return
	}
	ust := &netproto.DDZ_QiangDiZhu{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v抢庄数据结构错误%v", user, err)
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
	if ust.GetIsRob() {
		d.GameLog.AddGameStep("用户%d抢地主", user.UserID)
		u.isQiangDizhu = 2
	} else {
		d.GameLog.AddGameStep("用户%d不抢地主", user.UserID)
		u.isQiangDizhu = 1
	}
	if u.isQiangDizhu == 2 {
		d.ext.lastQiangDiZhu = user.SeatNo
		d.oddsChange(user.UserID, oddsQiangDiZhu, netproto.OddsChangeType_qiangdizhu)
	}
	idx, jd := d.ext.getJiaoDiZhu()
	if jd == nil {
		mlog.Error("系统错误，没有叫地主的玩家")
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，没有叫地主的玩家")
		return
	}
	if d.qiangDiZhuTimer != nil {
		d.qiangDiZhuTimer.Stop()
	}

	d.broadcastQiangDiZhu(user.SeatNo, ust.GetIsRob())
	if (d.ext.turn)%d.maxUserCount == int32(idx) || d.checkBuqiangCount() == 2 { //确定地主后进入明牌阶段
		dz := d.ext.getExtInfoByTurn(d.ext.lastQiangDiZhu)
		if dz == nil {
			mlog.Error("系统错误，没有地主玩家")
			d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，没有地主玩家")
			return
		}
		d.ext.resetLastDiZhu()
		dz.isDiZhu = true
		dz.isLastDiZhu = true
		d.ext.turn = dz.userSeat
		d.faPaiLastCards(dz.userId)
		d.broadcastGameStep(gameProgressShow)
		d.broadcastTimeOpt(dz.userSeat, d.ddt.GetTime(gameProgressShow))
		d.mingPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ddt.GetTime(gameProgressShow), d.mingPaiTimeOut, "")
	} else {
		d.responseUserQiangDiZhu(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
		d.ext.turn++
		d.qiangZhuangStart()
	}
}

//明牌
func (d *DDZGameLogic) onMingPai(user *gamelogic.GameUser, data []byte) {
	if d.gamestep != gameProgressShow {
		mlog.Error("非法操作，当前不是明牌(%d)", d.gamestep)
		d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前不是明牌主阶段")
		return
	}
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	if u == nil {
		mlog.Error("系统错误，获取用户扩展信息失败(%d)", d.gamestep)
		d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取用户信息失败")
		return
	}
	if user.SeatNo != u.userSeat {
		mlog.Error("非法操作，当前玩游戏的并非同一个人(%d)", d.gamestep)
		d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，当前玩游戏的并非同一个人")
		return
	}
	if !u.isDiZhu {
		mlog.Error("非法操作，非地主无法进行明牌操作(%d)", d.gamestep)
		d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，非地主无法进行明牌操作")
		return
	}
	ust := &netproto.DDZ_MingPai{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v抢庄数据结构错误%v", user, err)
		d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
	if d.mingPaiTimer != nil {
		d.mingPaiTimer.Stop()
	}
	if d.chuPaiTimer != nil {
		d.chuPaiTimer.Stop()
	}

	u.isMingPai = ust.GetIsMingPai()
	if u.isMingPai {
		d.GameLog.AddGameStep("用户%d明牌", user.UserID)
		d.oddsChange(user.UserID, oddsMingpai, netproto.OddsChangeType_mingpai)
	}
	d.responseUserMingPai(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
	d.broadcastMingPai(user, u.isMingPai)
	u.isDo = true
	//进入出牌阶段
	d.broadcastGameStep(gameProgressStart)
	d.broadcastTimeOpt(d.ext.getNowSeat(), d.ext.opttime)
	d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
}

func (d *DDZGameLogic) OnPass(user *gamelogic.GameUser, data []byte) {
	//检测牌是否合法 并步进座位
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	mlog.Debug("pass =============================用戶 %d %d", user.UserID, d.ext.turn)
	if u == nil {
		mlog.Error("系统错误，获取用户扩展信息失败(%d)", d.ext.turn)
		d.responseUserPass(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取用户信息失败")
		return
	}
	if u.userId != user.UserID {
		mlog.Warn("非法操作, 不该您出牌 userid %d, curr_id %d", u.userId, user.UserID)
		d.responseUserPass(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，不该您出牌")
		return
	}
	if u.isDiZhu && len(d.ext.his.actions) <= 0 {
		mlog.Error("非法操作，地主第一把必须出牌")
		d.responseUserPass(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，地主第一把必须出牌")
		return
	}
	ust := &netproto.DDZ_ChuPai{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析玩家出牌消息错误")
		d.responseUserPass(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
	if d.qiangDiZhuTimer != nil {
		d.qiangDiZhuTimer.Stop()
	}
	if d.chuPaiTimer != nil {
		d.chuPaiTimer.Stop()
	}
	d.GameLog.AddGameStep("用户%d选择了不出", user.UserID)

	//牌局紀錄 出牌過程
	//up := d.gamePoker.GetUserPoker(user.UserID)
	//d.record.addGameStep(user.SeatNo, true, nil, up.Cards)

	u.action = new(userAction)
	u.action.ct = cardTypeNone
	u.action.userSeat = u.userSeat
	d.responseUserPass(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
	d.broadcastChuPai(user.SeatNo, new(poker.CardGroup), cardTypeNone, 0, 0)
	d.ext.turn++
	tmp := d.ext.getExtInfoByTurn(d.ext.turn)
	if tmp.isTuoGuan {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[OnPass] userid[%v] isTuoGuan[%v]", user.UserID, u.isTuoGuan)).Debug("託管")
		d.broadcastTimeOpt(d.ext.turn, 1)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(1, d.chupaiTimeOut, "")
	} else {
		d.broadcastTimeOpt(d.ext.getNowSeat(), d.ext.opttime)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
	}
}

func (d *DDZGameLogic) RecordChuPai(u *userDdzOpt, action *userAction) {
	u.action = new(userAction)
	u.action.ct = action.ct
	u.action.userSeat = u.userSeat
	u.action.cg = new(poker.CardGroup)
	u.action.cg.Cards = make([]*poker.Card, 0)
	u.action.cg.Cards = append(u.action.cg.Cards, action.cg.Cards...)
	u.action.maxCard = action.maxCard
	u.action.mCard = action.mCard
}

//出牌
func (d *DDZGameLogic) onChuPai(user *gamelogic.GameUser, data []byte) {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("自動出牌ing userid[%v] btmsg.MsgData", user.UserID)).Debug("託管")
	if d.gamestep != gameProgressStart {
		if d.gamestep == gameProgressShow {
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("自動出牌ing userId[%v] gameProgressStart出牌", user.UserID)).Debug("託管")
			d.gamestep = gameProgressStart
			d.broadcastGameStep(d.gamestep)
			d.mingPaiTimer.Stop()
		} else {
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("自動出牌ing userId[%v] 系统错误，当前阶段不能出牌(%d)", user.UserID, d.gamestep)).Debug("託管")
			mlog.Error("系统错误，当前阶段不能出牌(%d)", d.gamestep)
			d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，当前阶段不能出牌")
			return
		}
	}
	//检测牌是否合法 并步进座位
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	mlog.Debug("chupai =============================用戶 %d %d", user.UserID, d.ext.turn)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[onChuPai] userid[%v] isTuoGuan[%v]", user.UserID, u.isTuoGuan)).Debug("託管")
	if u == nil {
		mlog.Error("userId[%v] 系统错误，获取用户扩展信息失败(%d)", user.UserID, d.ext.turn)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取用户信息失败")
		return
	}

	if user.SeatNo != u.userSeat {
		example.invalid("userId[%v] 非法操作，不该您出牌 u.isTuoGuan[%v] seatno %v, userseat %v ", user.UserID, u.isTuoGuan, u.userSeat, user.SeatNo)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作，不该您出牌")
		return
	}
	ust := &netproto.DDZ_ChuPai{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("userId[%v] 解析玩家出牌消息错误", user.UserID)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
	up := d.gamePoker.GetUserPoker(user.UserID)
	if up == nil {
		mlog.Error("userId[%v] 系统错误，获取玩家牌型失败", user.UserID)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "系统错误，获取玩家牌型失败")
		return
	}
	tmp := make([]*poker.Card, 0)
	if len(ust.Cards) > 20 {
		mlog.Error("userId[%v] 非法操作", user.UserID)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作")
		return
	}
	cg := new(poker.CardGroup)
	if len(ust.Cards) == 0 {
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
		d.chupaiTimeOut()
		return
	}
	for _, v := range ust.Cards {
		pp := new(poker.Card)
		pp.CardNumber = v.GetNumber()
		pp.Color = v.GetColor()
		cg.AddCard(pp)
		tmp = append(tmp, pp)
	}
	if !up.Cards.IsContent(tmp) {
		mlog.Error("userId[%v] 非法操作,玩家发送的牌不在手牌中", user.UserID)
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作,玩家发送的牌不在手牌中")
		return
	}
	action := new(userAction)
	action.cg = cg
	action.userSeat = user.SeatNo
	action.ct = d.ext.his.cardTypeCheck(cg, action)
	//d.GameLog.AddGameStep("用户%d出牌%s", user.UserID, cg.String())
	his := d.ext.his.GetBack()
	if action.ct != cardTypeNone {
		if his != nil {
			if action.userSeat == his.userSeat {
				if action.ct == cardTypeBomb {
					d.oddsChange(user.UserID, oddsZhaDan, netproto.OddsChangeType_zhadan)
				}
				if action.ct == cardTypeRocket {
					d.oddsChange(user.UserID, oddsRocket, netproto.OddsChangeType_rocket)
				}
				d.ext.his.Push(action)
			} else {
				if d.followCheckCardType(user.UserID, action, his) {
					d.ext.his.Push(action)
				} else {
					mlog.NewDefaultLogger().SetCallDepLen(10).SetData(fmt.Sprintf("userid[%v] 非法操作,出牌错误,action.ct=%d,action.maxCard=%d,his.ct=%d,his.maxCard[%v], actionCards[%v], his.cg.Cards[%v]", user.UserID, action.ct, action.maxCard, his.ct, his.maxCard, action.cg.Cards, his.cg.Cards)).Error("出牌")
					d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作,出牌错误")

					if conf.TestMode {
						mlog.Fatal("stop, 開發階段有錯誤要修好.")
					}

					return
				}
			}
		} else {
			if action.ct == cardTypeBomb {
				d.oddsChange(user.UserID, oddsZhaDan, netproto.OddsChangeType_zhadan)
			}
			if action.ct == cardTypeRocket {
				d.oddsChange(user.UserID, oddsRocket, netproto.OddsChangeType_rocket)
			}
			d.ext.his.Push(action)
		}
		if d.qiangDiZhuTimer != nil {
			d.qiangDiZhuTimer.Stop()
		}
		if d.chuPaiTimer != nil {
			d.chuPaiTimer.Stop()
		}

		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[RecordChuPai] userid[%v] isTuoGuan[%v]", user.UserID, u.isTuoGuan)).Debug("託管")
		d.RecordChuPai(u, action)

		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[responseUserChuPai] userid[%v] isTuoGuan[%v]", user.UserID, u.isTuoGuan)).Debug("託管")
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")

		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[broadcastChuPai] userid[%v] isTuoGuan[%v]", user.UserID, u.isTuoGuan)).Debug("託管")
		d.broadcastChuPai(user.SeatNo, cg, action.ct, action.maxCard, action.mCard)

		up.Cards.DeleteGroup(cg)
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[DeleteGroup] userid[%v] cg[%v]", user.UserID, cg)).Debug("託管")
		//農民春天判斷
		if !u.isDiZhu {
			d.ext.isChunTian = false
		} else { //地主春天判斷
			d.diZhuChuPaiCount++
			if d.diZhuChuPaiCount > 1 {
				d.ext.isFanChunTian = false
			}
		}

		if up.Cards.Len() <= 0 { //出牌结束 进入结算流程
			if d.ext.isChunTian {
				d.oddsChange(u.userId, oddsChunTian, netproto.OddsChangeType_chuntian)
			} else if d.ext.isFanChunTian {
				d.oddsChange(u.userId, oddsFanChunTian, netproto.OddsChangeType_fanchuntian)
			}
			d.settlement(u)
			mlog.Debug("玩家%d,牌已出完进入结算流程", u.userId)
		} else {
			d.ext.turn++
			if action.ct == cardTypeBomb || action.ct == cardTypeRocket {
				d.actionName_curr = d.RegisterAnimation(actionName_ZhaDan, 1, d.sendNextOptMsg)
			} else {
				d.sendNextOptMsg()
			}
		}
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("用户[%d]出牌%s 手牌還有 %v", user.UserID, cg.String(), up.Cards)).Warn("託管")
		d.GameLog.AddGameStep("用户[%d]出牌%s 手牌還有 %v", user.UserID, cg.String(), up.Cards)

		//牌局紀錄 出牌過程
		//d.record.addGameStep(user.SeatNo, false, cg, up.Cards)

	} else {
		mlog.Error("非法操作,牌型异常")
		d.responseUserChuPai(user.UserID, netproto.DDZ_ErrorCode_Fail, "非法操作,牌型异常")
		return
	}
}

//跟牌比對是否符合跟牌條件
func (d *DDZGameLogic) followCheckCardType(userId int32, action *userAction, his *userAction) bool {
	if (action.ct == cardTypeRocket || //火箭最大
		(action.ct == cardTypeBomb && action.maxCard == 2 && his.ct != cardTypeRocket) || //如果是2的炸弹且上手牌 不是火箭
		(action.ct == cardTypeBomb && action.ct == his.ct && action.IsBig(his)) || //手上的炸弹比桌面上的炸弹大时
		(action.ct == cardTypeBomb && his.ct != cardTypeRocket && his.ct != cardTypeBomb)) || //手上的牌是炸弹 桌面的牌不是火箭 也不是炸弹时
		(action.ct == his.ct && action.IsBig(his) && len(his.cg.Cards) == len(action.cg.Cards)) {
		if action.ct == cardTypeBomb {
			d.oddsChange(userId, oddsZhaDan, netproto.OddsChangeType_zhadan)
		}
		if action.ct == cardTypeRocket {
			d.oddsChange(userId, oddsRocket, netproto.OddsChangeType_rocket)
		}
		return true
	}
	return false
}

// 跟牌
func (d *DDZGameLogic) FollowCards(u *gamelogic.GameUser, outCardType cardType, more int32) (*poker.CardGroup, cardType, int32, int32) {
	up := d.gamePoker.GetUserPoker(u.UserID)
	if up == nil {
		mlog.Error("系统错误FollowCards，获取玩家牌型失败")
		p := new(poker.CardGroup)
		return p, cardTypeNone, int32(0), int32(0)
	}
	d.recombination(u.UserID, up.Cards)

	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌前 outCardType[%v] more[%v]", u.UserID, outCardType, more)).Debug("託管")

	//增加主動出牌時, 其他人pass
	var outCard []int32
	his := d.ext.his.GetBack()

	if his != nil && u.SeatNo == his.userSeat { //自己主動出牌時
		for i := int(13); i >= int(0); i-- {

			upk := d.gamePoker.GetUserPoker(u.UserID)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌前其他人pass cardType[%v] 剩下手牌[%v]", u.UserID, i, upk.Cards.Cards)).Debug("託管")
			outCard = d.outCardByType(cardType(i), u.UserID, 0)
			if len(outCard) > 0 {
				outCardType = cardType(i)
				mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌前其他人pass 出牌成功cardType[%v] outCard[%v]", u.UserID, i, outCard)).Debug("託管")
				break
			}
		}
	} else { //壓其他玩家牌時
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌前壓其他玩家牌時 outCardType[%v] more[%v]", u.UserID, outCardType, more)).Debug("託管")
		outCard = d.outCardByType(outCardType, u.UserID, more)
	}

	tmpCardGroup := copyPokerCards(up.Cards)
	uPoker := new(poker.CardGroup)

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌後 outCard[%v] outCardType[%v] tmpCardGroup[%v]", u.UserID, outCard, outCardType, tmpCardGroup)).Debug("託管")

	outCardNumber := int32(0)
	//組合出牌
	for _, n := range outCard {
		for _, cd := range tmpCardGroup.Cards {
			if cd.CardNumber == n {

				pp := new(poker.Card)
				pp.CardNumber = cd.CardNumber
				pp.Color = cd.Color
				uPoker.AddCard(pp)
				mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合廣播自動出牌ing uPoker.AddCard(pp) number[%v] color[%v]", u.UserID, pp.CardNumber, pp.Color)).Debug("託管")

				if pp.CardNumber > outCardNumber {
					outCardNumber = pp.CardNumber
				}
				cd.CardNumber = int32(0)

				break
			}
			if len(uPoker.Cards) == len(outCard) {
				break
			}
		}
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("自動出牌前 userid[%v] outCard[%v] lenOutCard[%v] outCardType[%v]", u.UserID, outCard, len(outCard), outCardType)).Debug("託管")

	if len(outCard) == int(0) {
		outCardNumber = int32(0)
		outCardType = cardTypeNone
		return uPoker, outCardType, outCardNumber, int32(0)
	}

	//check card type
	outCardType = checkOutCardType(uPoker, outCardType)

	return uPoker, outCardType, outCardNumber, int32(1)
}
func checkOutCardType(uPoker *poker.CardGroup, outCardType cardType) cardType {
	l := len(uPoker.Cards)
	if l == 2 {
		ctx := sortCtxByCardGroup(l, uPoker)
		nMin := ctx.GetValueByIndex(0)
		nMax := ctx.GetValueByIndex(1)

		if nMin == 15 && nMax == 16 {
			outCardType = cardTypeRocket
		}

	}
	if l == 4 {
		ctx := sortCtxByCardGroup(l, uPoker)
		ua := new(userAction)
		outCardType = fourCheck(uPoker, ctx, ua)
	}
	if l > 4 {
		ctx := sortCtxByCardGroup(l, uPoker)
		ua := new(userAction)
		outCardType = otherCheck(uPoker, ctx, ua)
	}
	return outCardType
}

func sortCtxByCardGroup(l int, uPoker *poker.CardGroup) *sortCtx {
	ctx := new(sortCtx)
	ctx.arr = make([]*sItem, 0)
	m := make(map[int32]int32)
	for i := 0; i < l; i++ {
		m[uPoker.Cards[i].CardNumber] += 1
	}
	for k, v := range m {
		item := new(sItem)
		item.number = k
		item.count = v
		ctx.arr = append(ctx.arr, item)
	}
	sort.Sort(ctx)
	return ctx
}

//拷贝牌型
func copyPokerCards(Cards *poker.CardGroup) *poker.CardGroup {
	tmpCards := poker.NewEmptyCardGroup(len(Cards.Cards))
	for i := 0; i < len(Cards.Cards); i++ {
		tmpCards.Cards[i].Color = Cards.Cards[i].Color
		tmpCards.Cards[i].CardNumber = Cards.Cards[i].CardNumber
	}

	return tmpCards
}

//转化list为牌组
func (lgc *DDZGameLogic) list2cardGroup(cg *poker.CardGroup, cdlist []int32) *poker.CardGroup {
	mlog.Debug("出牌[%v]", cg.Cards)
	cg2 := new(poker.CardGroup)
	cgcal := new(poker.CardGroup)
	for _, cgcd := range cg.Cards {
		cd := new(poker.Card)
		cd.Color = cgcd.Color
		cd.CardNumber = cgcd.CardNumber
		cgcal.Cards = append(cgcal.Cards, cd)
	}
	for _, cd := range cdlist {
		for _, cgcd := range cgcal.Cards {
			if cd == cgcd.CardNumber {
				cg2.AddCard(cgcd)
				cgcal.DeleteCard(cgcd)
				break
			}
		}
	}
	return cg2
}

//减少牌型
func (lgc *DDZGameLogic) reduceCardType(robotUser *robotUserInfoLocal, bct cardType, num int32) {
	// mlog.Debug("减取牌型[%d],数量[%d]", int(bct), num)
	robotUser.cardTypesCount[bct] -= num
}

//重新拆牌, 整理完牌型後就可以指定牌型出牌
func (d *DDZGameLogic) recombination(userid int32, cg *poker.CardGroup) {
	//重新拆牌
	result := d.Combination(cg, userid)
	ctg := new(robotUserInfoLocal)
	ctg.cardTypesCount = make(map[cardType]int32, 0)
	ctg.ctg = result

	for _, data := range result.Ct {
		ctg.cardTypesCount[data.byteCardType]++
	}
	ctg.cardTypesCount[cardTypeSingle] += int32(len(result.SingleCards))
	//mlog.Debug("牌组大小[%d]", len(ctg.ctg.Ct))
	//d.robotUser[userid] = ctg
	d.tuoGuanUser[userid] = ctg
}

//拆分手牌
func (d *DDZGameLogic) Combination(cg *poker.CardGroup, userid int32) *CtGroup {
	//mlog.Debug("开始拆分手牌")
	cMin := int32(0)
	resultBest := new(CtGroup)
	resultBest.Ct = make([]*cardTypeGroup, 0)
	singleBest := make([]int32, 0)
	for i := 1; i <= 7; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)
		single := make([]int32, 0)
		mapCard := mapCardNumberCount(cg)
		mapCard = d.combinationKing(mapCard, result)                               //王
		mapCard = d.combinationTwo(mapCard, result)                                //2
		mapCard = d.combinationCardsIndex(mapCard, cardTypeBomb, result, int32(i)) //剩下单牌，

		erasePair := make([]int32, 0)
		eraseThree := make([]int32, 0)
		//mlog.Debug("补张开始")
		//把对子或者三放入单牌中看能否消掉一些，这里单牌还没有带
		pairList, threeList := d.getPairAndThree(result)

		for _, cd := range pairList {
			mapCard[cd] = 2
		}
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result) //顺子

		for _, cd := range threeList {
			mapCard[cd] = 3
			mapCard = d.combinationCards(mapCard, cardTypeStraight, result) //顺子
			if mapCard[cd] != 3 {
				//标记，并且去除结果中的三条
				eraseThree = append(eraseThree, cd)
			} else {
				mapCard[cd] = 0 //如果没有匹配到，从重置为0
			}
		}

		for _, cd := range pairList {
			if mapCard[cd] != 2 {
				//标记，并且去除结果中的对子
				erasePair = append(erasePair, cd)
			} else {
				mapCard[cd] = 0
			}
		}
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result) //对子

		//跟新result,erasePair,eraseThree
		result = d.updataResult(result, erasePair, eraseThree)

		mapCard, single = d.combinationSingle(mapCard) //单牌
		for number, count := range mapCard {
			if count > 0 {
				mlog.Error("错误数据 number[%d],count[%d]", number, count)
			}
		}

		singleBest = single //託管預設
		resultBest = result
	}

	resultBest.SingleCards = singleBest
	resultBest.hands = cMin

	return resultBest
}

//获取手数最小的牌 ，如果手数一样获取权值最大的牌
func (d *DDZGameLogic) combinationSingle(mapCard map[int32]int32) (mapReturn map[int32]int32, sg []int32) {
	single := make([]int32, 0)
	for number, count := range mapCard {
		if count == 1 {
			single = append(single, number)
			mapCard[number] = 0
		}
	}
	single = d.sort(single)
	return mapCard, single
}

//排序
func (d *DDZGameLogic) sort(single []int32) []int32 {
	for i := 0; i < len(single); i++ {
		for j := i + 1; j < len(single); j++ {
			if single[i] > single[j] {
				single[i], single[j] = single[j], single[i]
			}
		}
	}
	return single
}

//跟新result,erasePair,eraseThree
func (d *DDZGameLogic) updataResult(result *CtGroup, erasePair, eraseThree []int32) *CtGroup {
	for _, cd := range erasePair {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeDouble && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	for _, cd := range eraseThree {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeThree && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	resultNew := new(CtGroup)
	resultNew.Ct = make([]*cardTypeGroup, 0)
	for _, ctg := range result.Ct {
		if ctg.byteCardType != cardTypeNone {
			ctg2 := new(cardTypeGroup)
			ctg2.byteCardType = ctg.byteCardType
			ctg2.nMax = ctg.nMax
			ctg2.nMin = ctg.nMin
			resultNew.Ct = append(resultNew.Ct, ctg2)
		}
	}
	return resultNew

}

//提取对子和三张
func (d *DDZGameLogic) getPairAndThree(result *CtGroup) (pair, three []int32) {

	pairList := make([]int32, 0)
	threeList := make([]int32, 0)
	for _, ctg := range result.Ct {
		//三带
		if ctg.nMax == 2 {
			continue
		}
		if ctg.byteCardType == cardTypeThree {
			threeList = append(threeList, ctg.nMax)
			continue
		}
		//对子
		if ctg.byteCardType == cardTypeDouble {
			pairList = append(pairList, ctg.nMax)
			continue
		}
	}
	return pairList, threeList

}

//按类型组牌
func (d *DDZGameLogic) combinationCards(mapCard map[int32]int32, label cardType, result *CtGroup) map[int32]int32 {
	if result == nil {
		return mapCard
	}
	if label == cardTypeDouble {
		for number, count := range mapCard {
			if count == 2 {
				//mlog.Debug("找到对子[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeDouble
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeThree {
		for number, count := range mapCard {
			if count == 3 {
				//mlog.Debug("找到三条[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeThree
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeBomb {
		for number, count := range mapCard {
			if count == 4 {
				//mlog.Debug("找到炸弹[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeBomb
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeStraight { //顺子
		d.checkStraightLong(mapCard, 1, 5, cardTypeStraight, result)
	} else if label == cardTypeThreeStraight { //连队
		d.checkStraightLong(mapCard, 2, 3, cardTypeThreeStraight, result)
	} else if label == cardTypePlaneAndWing { //飞机
		d.checkStraightLong(mapCard, 3, 2, cardTypePlaneAndWing, result)
	}

	return mapCard
}

//连对,顺子，飞机,size表示至少要几连
func (d *DDZGameLogic) checkStraightLong(mapCard map[int32]int32, zhang, size int32, bType cardType, result *CtGroup) map[int32]int32 {
	count := int32(0)
	min := int32(2)
	max := int32(2)
	for i := int32(2); i <= 15; i++ {
		//mlog.Debug("前count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
		if mapCard[i] >= zhang {
			count++
			max = i
		} else {
			if count >= size {
				//mlog.Debug("找到[%d,%d]一个[%d~%d]", zhang, size, min, max)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = bType
				ctg.nMax = max
				ctg.nMin = min
				result.Ct = append(result.Ct, ctg)
				for j := min; j <= max; j++ {
					mapCard[j] -= zhang
				}
				mapCard = d.checkStraightLong(mapCard, zhang, size, bType, result)
			}
			//mlog.Debug("重置count")
			count = 0
			min = i + 1
		}
		//mlog.Debug("后count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
	}
	return mapCard
}

// 随机拆牌
func (d *DDZGameLogic) combinationCardsIndex(mapCard map[int32]int32, label cardType, result *CtGroup, index int32) map[int32]int32 {

	if index == 1 {
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 2 {
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 3 {
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 4 {
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 5 {
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 6 {
		mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子
		mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		return mapCard
	}
	mapCard = d.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
	mapCard = d.combinationCards(mapCard, cardTypeStraight, result)      //顺子
	mapCard = d.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
	mapCard = d.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
	mapCard = d.combinationCards(mapCard, cardTypeThree, result)         //三条
	mapCard = d.combinationCards(mapCard, cardTypeDouble, result)        //对子

	return mapCard

}

//拆2
func (d *DDZGameLogic) combinationTwo(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {
	//mlog.Debug("拆分2")
	if mapCard[2] > 0 {
		ctg := new(cardTypeGroup)
		ctg.nMax = 2
		if mapCard[2] == 1 {
			//mlog.Debug("找到一个2")
			ctg.byteCardType = cardTypeSingle
		} else if mapCard[2] == 2 {
			//mlog.Debug("找到一对2")
			ctg.byteCardType = cardTypeDouble
		} else if mapCard[2] == 3 {
			//mlog.Debug("找到三个2")
			ctg.byteCardType = cardTypeThree
		} else if mapCard[2] == 4 {
			//mlog.Debug("找到炸弹2")
			ctg.byteCardType = cardTypeBomb
		}
		result.Ct = append(result.Ct, ctg)
	}
	mapCard[2] = 0
	return mapCard
}

//拆王
func (d *DDZGameLogic) combinationKing(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {

	if mapCard[poker.PokerNumberSmallKing]+mapCard[poker.PokerNumberBigKing] == 2 {
		//mlog.Debug("找到火箭")
		ctg := new(cardTypeGroup)
		ctg.byteCardType = cardTypeRocket
		ctg.nMax = poker.PokerNumberBigKing
		ctg.nMin = poker.PokerNumberSmallKing
		result.Ct = append(result.Ct, ctg)
	} else {
		if mapCard[poker.PokerNumberSmallKing] > 0 {
			//mlog.Debug("找到小王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberSmallKing
			result.Ct = append(result.Ct, ctg)
		}
		if mapCard[poker.PokerNumberBigKing] > 0 {
			//mlog.Debug("找到大王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberBigKing
			result.Ct = append(result.Ct, ctg)
		}
	}
	mapCard[poker.PokerNumberSmallKing] = 0
	mapCard[poker.PokerNumberBigKing] = 0
	return mapCard
}

//託管自動出牌(跟牌)
func (d *DDZGameLogic) onAutoChuPaiByTuoGuan() (*poker.CardGroup, cardType, int32, int32) {
	u := d.ext.getExtInfoByTurn(d.ext.getNowSeat())
	user := d.GetUserByUserID(u.userId)

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管自動出牌前 userid[%v] isTuoGuan[%v] user.UserID[%v] != u.userId[%v] turn[%v] seat[%v]",
		user.UserID, u.isTuoGuan, user.UserID, u.userId, d.ext.turn, u.userSeat)).Debug("託管")

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管自動出牌前 userid[%v] isTuoGuan[%v]  ",
		user.UserID, u.isTuoGuan)).Debug("託管")

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管自動出牌後 userid[%v] isTuoGuan[%v]", user.UserID, !u.isTuoGuan)).Debug("託管")

	//跟牌
	mlog.Debug("UserId[%d]自動跟牌", user.UserID)

	//牌桌上最後出的牌
	his := d.ext.his.GetBack()

	var autoCg *poker.CardGroup
	var autoCardType cardType
	var autoCardNumber int32
	var autoMainCard int32

	if his == nil {
		autoCg, autoCardType, autoCardNumber, autoMainCard = d.FollowCards(user, cardTypeSingle, 2)
	} else {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管自動出牌前FollowCards userid[%v]  his.ct[%v], his.maxCard[%v]", user.UserID, his.ct, his.maxCard)).Debug("託管")
		autoCg, autoCardType, autoCardNumber, autoMainCard = d.FollowCards(user, his.ct, his.maxCard)
	}

	return autoCg, autoCardType, autoCardNumber, autoMainCard
}

func (d *DDZGameLogic) sendNextOptMsg() {
	u := d.ext.getExtInfoByTurn(d.ext.turn)
	if u.isTuoGuan {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[sendNextOptMsg] userid[%v] isTuoGuan[%v]", u.userId, u.isTuoGuan)).Debug("託管")
		d.broadcastTimeOpt(d.ext.turn, 1)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(1, d.chupaiTimeOut, "")
	} else {
		d.broadcastTimeOpt(d.ext.getNowSeat(), d.ext.opttime)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
	}
}

func (d *DDZGameLogic) OnWaitOpt(user *gamelogic.GameUser, data []byte) {
	ust := &netproto.UserActionCommited{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析玩家出牌消息错误")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}
}

func (d *DDZGameLogic) AutoChuPai() {
	if d.gamestep == gameProgressGameEnd {
		return
	}
	mlog.Debug("chupaiTimeOut enter.....")
	if d.qiangDiZhuTimer != nil {
		d.qiangDiZhuTimer.Stop()
	}
	if d.chuPaiTimer != nil {
		d.chuPaiTimer.Stop()
	}

	u := d.ext.getExtInfoByTurn(d.ext.getNowSeat())
	his := d.ext.his.GetBack()
	upk := d.gamePoker.GetUserPoker(u.userId)
	if his == nil || his.userSeat == u.userSeat { //主動出牌時
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] AutoChuPai() 最後出牌是自己", u.userId)).Debug("託管")
		if upk == nil || len(upk.Cards.Cards) <= 0 {
			mlog.Error("获取玩家手牌信息出错！（%v）", upk)
			return
		}

	} else { //跟牌時
		logUser := d.ext.getExtInfoByTurn(d.ext.getNowSeat())
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] AutoChuPai() 最後出牌是別人[%v]", u.userId, logUser.userId)).Debug("託管")

		if !u.isTuoGuan {
			u.isTuoGuan = true
			d.broadcastTuoGuanState(u.userSeat, u.isTuoGuan)
		}

		mlog.Debug("auto chupai ok buchu!")
	}

	//cardsQty := len(upk.Cards.Cards)
	//if cardsQty >= 10 {
	//	mlog.Debug("玩家手牌太多, 託管直接不出")
	//	d.broadcastChuPai(d.ext.getNowSeat(), new(poker.CardGroup), cardTypeNone, 0, 0)
	//} else {
	//	d.autoChuPaiStart(u)
	//}
	d.autoChuPaiStart(u)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管前... UserId[%v] SeatNo[%v]", u.userId, u.userSeat)).Debug("託管")
	d.ext.turn++
	nu := d.ext.getExtInfoByTurn(d.ext.turn)
	if nu.isTuoGuan {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管中... 下個位置 UserId[%v] SeatNo[%v]", nu.userId, nu.userSeat)).Debug("託管")
		d.broadcastTimeOpt(d.ext.getNowSeat(), 1)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(1, d.chupaiTimeOut, "")
		//d.broadcastTimeOpt(d.ext.getNowSeat(), d.ext.opttime)
		//d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
		mlog.Debug("broadcastTimeOpt ok 1!")
	} else {
		d.broadcastTimeOpt(d.ext.getNowSeat(), d.ext.opttime)
		d.chuPaiTimer = d.TimerFuncManager.NewNamedTimer(d.ext.opttime, d.chupaiTimeOut, "")
		mlog.Debug("broadcastTimeOpt ok 2!")
	}
}

//開始自動出牌與檢查出牌牌行 , TODO 可以與出牌整合onChuPai或共用部分程式
func (d *DDZGameLogic) autoChuPaiStart(u *userDdzOpt) {
	autoCg, autoCardType, autoCardNumber, autoMainCard := d.onAutoChuPaiByTuoGuan()

	u.action = new(userAction)
	u.action.cg = autoCg
	u.action.ct = autoCardType
	u.action.maxCard = autoCardNumber
	u.action.mCard = autoMainCard
	u.action.userSeat = u.userSeat

	//d.GameLog.AddGameStep("用户%d出牌%s", u.userId, autoCg.String())
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("託管廣播前... UserId[%v] SeatNo[%v] cg[%v] autoCardType[%v]", u.userId, d.ext.getNowSeat(), autoCg, autoCardType)).Debug("託管")
	if len(autoCg.Cards) == int(0) { //不出時
		d.broadcastChuPai(d.ext.getNowSeat(), autoCg, autoCardType, autoCardNumber, autoMainCard)
	} else { //出牌
		up := d.gamePoker.GetUserPoker(u.userId)
		if up == nil {
			mlog.Error("userId[%v] 系统错误，获取玩家牌型失败", u.userId)
			d.responseUserChuPai(u.userId, netproto.DDZ_ErrorCode_Fail, "系统错误，获取玩家牌型失败")
			return
		}

		//檢查出牌內容並廣撥出牌
		action := new(userAction)
		action.cg = autoCg
		action.userSeat = u.userSeat
		action.ct = d.ext.his.cardTypeCheck(autoCg, action)
		//d.GameLog.AddGameStep("用户%d出牌%s", u.userId, autoCg.String())

		his := d.ext.his.GetBack()
		if action.ct != cardTypeNone {
			if his != nil {
				if action.userSeat == his.userSeat {
					if action.ct == cardTypeBomb {
						d.oddsChange(u.userId, oddsZhaDan, netproto.OddsChangeType_zhadan)
					}
					if action.ct == cardTypeRocket {
						d.oddsChange(u.userId, oddsRocket, netproto.OddsChangeType_rocket)
					}
					d.ext.his.Push(action)
				} else {
					if d.followCheckCardType(u.userId, action, his) {
						d.ext.his.Push(action)
					} else {
						mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 非法操作,出牌错误,action.ct=%d,action.maxCard=%d,his.ct=%d,his.maxCard[%v], actionCards[%v], his.cg.Cards[%v]", u.userId, action.ct, action.maxCard, his.ct, his.maxCard, action.cg.Cards, his.cg.Cards)).Error("託管")
						d.responseUserChuPai(u.userId, netproto.DDZ_ErrorCode_Fail, "非法操作,出牌错误")

						if conf.TestMode {
							mlog.Fatal("stop")
						}

						return
					}
				}
			} else {
				if action.ct == cardTypeBomb {
					d.oddsChange(u.userId, oddsZhaDan, netproto.OddsChangeType_zhadan)
				}
				if action.ct == cardTypeRocket {
					d.oddsChange(u.userId, oddsRocket, netproto.OddsChangeType_rocket)
				}
				d.ext.his.Push(action)
			}
			if d.qiangDiZhuTimer != nil {
				d.qiangDiZhuTimer.Stop()
			}
			if d.chuPaiTimer != nil {
				d.chuPaiTimer.Stop()
			}
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[RecordChuPai] userid[%v] isTuoGuan[%v]", u.userId, u.isTuoGuan)).Debug("託管")
			d.RecordChuPai(u, action)

			mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[responseUserChuPai] userid[%v] isTuoGuan[%v]", u.userId, u.isTuoGuan)).Debug("託管")
			d.responseUserChuPai(u.userId, netproto.DDZ_ErrorCode_OK, "OK")

			mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[broadcastChuPai] userid[%v] isTuoGuan[%v]", u.userId, u.isTuoGuan)).Debug("託管")
			d.broadcastChuPai(u.userSeat, autoCg, action.ct, action.maxCard, action.mCard)

			up.Cards.DeleteGroup(autoCg)

			mlog.NewDefaultLogger().SetData(fmt.Sprintf("出牌[DeleteGroup] userid[%v] cg[%v]", u.userId, autoCg)).Debug("託管")

			//農民春天判斷
			if !u.isDiZhu {
				d.ext.isChunTian = false
			} else { //地主春天判斷
				d.diZhuChuPaiCount++
				if d.diZhuChuPaiCount > 1 {
					d.ext.isFanChunTian = false
				}
			}

			if up.Cards.Len() <= 0 { //出牌结束 进入结算流程
				if d.ext.isChunTian {
					d.oddsChange(u.userId, oddsChunTian, netproto.OddsChangeType_chuntian)
				} else if d.ext.isFanChunTian {
					d.oddsChange(u.userId, oddsFanChunTian, netproto.OddsChangeType_fanchuntian)
				}
				d.settlement(u)
				mlog.Debug("玩家%d,牌已出完进入结算流程", u.userId)
			}
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("用户%d出牌%s", u.userId, autoCg.String())).Debug("託管")
			mlog.Debug("用户[%d] 出牌[%s] _ 手牌還有 [%v]", u.userId, autoCg.String(), up.Cards)
			d.GameLog.AddGameStep("用户[%d] 出牌[%s] _ 手牌還有 [%v]", u.userId, autoCg.String(), up.Cards)

			//牌局紀錄 出牌過程
			//d.record.addGameStep(u.userSeat, false, autoCg, up.Cards)

		} else {
			mlog.Error("非法操作,牌型异常")
			d.responseUserChuPai(u.userId, netproto.DDZ_ErrorCode_Fail, "非法操作,牌型异常")
			return
		}
	}

}

//断线重连
func (d *DDZGameLogic) UserRecover(gameuser *gamelogic.GameUser) {
	d.setExperienceCash(gameuser)
	d.FixedGameLogicBase.UserRecover(gameuser)
	d.broadcastTuoGuanState(gameuser.SeatNo, false)
	d.sendDeskCfg(gameuser)
	d.sendGameInfoToUser(gameuser.UserID)
	d.SendChiPaiChiStatus(gameuser)
	d.SendChiPaiChiInfo(gameuser)
}

//設定體驗幣
func (d *DDZGameLogic) setExperienceCash(gameuser *gamelogic.GameUser) {
	flag, _ := strconv.ParseInt(d.TableManager.Room.Flag, 10, 32)
	if flag == ExperienceCashFlag && !gameuser.IsRobot() {
		gameuser.Score = gameuser.ExperienceCash
	}
}

func (d *DDZGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	d.FixedGameLogicBase.UserEnter(gameuser, tableno, seatno)

	d.setExperienceCash(gameuser)
	d.sendDeskCfg(gameuser)
	d.SendChiPaiChiStatus(gameuser)
}

//計牌器資料:剩餘牌紀錄, 上下家出牌紀錄
func (d *DDZGameLogic) SendChiPaiChiInfo(gameuser *gamelogic.GameUser) {

	if len(d.ext.his.actions) == 0 {
		mlog.Warn("no his action")
		return
	}

	if !enableChiPaiChi(gameuser, d) {
		mlog.Warn("non-enable")
		return
	}

	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_ChiPaiChiInfo)
	gsdata := new(netproto.DDZ_ChiPaiChiInfo)
	gsdata.ShengYuPai = d.getShengYuPai(gameuser) //剩餘牌

	seatArr := make([]int32, 0)

	for i := int32(0); i < 3; i++ {
		if gameuser.SeatNo != i {
			seatArr = append(seatArr, i)
		}
	}

	gsdata.ShangChia = d.getShangChia(seatArr[0]) //上家牌
	gsdata.HsiaChia = d.getHsiaChia(seatArr[1])   //下家牌
	gsmsg.MsgData = gsdata

	d.SendMessageToUser(gameuser, gsmsg)
}

//上家出過的牌
func (d *DDZGameLogic) getShangChia(seatno int32) []*netproto.DDZ_ChuPaiChiLu {
	shangChia := make([]*netproto.DDZ_ChuPaiChiLu, 0)

	shangChia = d.getTurnChia(seatno)
	for _, v := range shangChia {
		for _, p := range v.Cards {
			mlog.Debug("上家 取得該位置[%v]的出牌紀錄 cardNumber[%v]", v.GetSeatNo(), p.GetNumber())
		}
	}

	return shangChia
}

//下家出過的牌
func (d *DDZGameLogic) getHsiaChia(seatno int32) []*netproto.DDZ_ChuPaiChiLu {
	shangChia := make([]*netproto.DDZ_ChuPaiChiLu, 0)

	shangChia = d.getTurnChia(seatno)

	for _, v := range shangChia {
		for _, p := range v.Cards {
			mlog.Debug("下家 取得該位置[%v]的出牌紀錄 cardNumber[%v]", v.GetSeatNo(), p.GetNumber())
		}
	}
	return shangChia
}

//取得該位置的出牌紀錄
func (d *DDZGameLogic) getTurnChia(seatno int32) []*netproto.DDZ_ChuPaiChiLu {
	hisPai := make([]*netproto.DDZ_ChuPaiChiLu, 0)
	nu := d.ext.getExtInfoByTurn(seatno)

	for _, v := range d.ext.his.actions {
		if seatno == v.userSeat {
			mlog.Debug("取得該位置的出牌紀錄 seatno[%v] userseat[%v] userid[%v]", seatno, v.userSeat, nu.userId)
			chuPaiChiLu := new(netproto.DDZ_ChuPaiChiLu)
			cardGroup := make([]*netproto.Card, 0)
			for _, p := range v.cg.Cards {
				card := new(netproto.Card)
				card.Color = proto.Int32(p.Color)
				card.Number = proto.Int32(p.CardNumber)
				cardGroup = append(cardGroup, card)
			}
			chuPaiChiLu.Cards = cardGroup
			chuPaiChiLu.SeatNo = proto.Int32(v.userSeat)
			hisPai = append(hisPai, chuPaiChiLu)
			mlog.Debug("hispai userid[%v] seat[%v] hispai[%v]", nu.userId, nu.userSeat, hisPai)
		}
	}

	return hisPai
}

//剩餘牌
func (d *DDZGameLogic) getShengYuPai(gameuser *gamelogic.GameUser) []*netproto.Card {
	shengYuPai := make([]*netproto.Card, 0)

	if len(d.ext.his.actions) == 0 {
		return shengYuPai
	}

	if !enableChiPaiChi(gameuser, d) {
		return shengYuPai
	}

	for _, u := range d.ext.userOpts {
		if gameuser.UserID != u.userId {
			upk := d.gamePoker.GetUserPoker(u.userId)
			for _, p := range upk.Cards.Cards {
				card := new(netproto.Card)
				card.Color = proto.Int32(p.Color)
				card.Number = proto.Int32(p.CardNumber)
				shengYuPai = append(shengYuPai, card)
			}
		}
	}
	mlog.Debug("剩餘牌 [%v]", shengYuPai)
	return shengYuPai
}

//是否有計牌器功能
func (d *DDZGameLogic) SendChiPaiChiStatus(gameuser *gamelogic.GameUser) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_EnableUserChiPaiChi)
	gsdata := new(netproto.DDZ_EnableUserChiPaiChi)
	status := false
	if enableChiPaiChi(gameuser, d) {
		status = true
	}
	mlog.Debug("DDZ_EnableUserChiPaiChi [%v]", status)

	gsdata.Status = proto.Bool(status)
	gsmsg.MsgData = gsdata
	d.SendMessageToUser(gameuser, gsmsg)
}

//開啟計牌器
func enableChiPaiChi(gameuser *gamelogic.GameUser, d *DDZGameLogic) bool {
	// enable := false

	if d.ext.bindaccount {
		return gameuser.SumCharge >= int64(d.ext.payminmoney) && gameuser.UserType == gamelogic.UserTypeNormal
	} else { //不用綁定也能使用
		return true
	}

	// return enable
}

func (d *DDZGameLogic) OnTuoGuan(user *gamelogic.GameUser, data []byte) {
	if d.gamestep != gameProgressStart {
		mlog.Error("当前阶段不能托管 step%d", d.gamestep)
		return
	}
	ust := &netproto.DDZ_TuoGuan{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析玩家出牌消息错误")
		d.KickUserFunc(user.UserID, "非法操作")
		return
	}

	d.GameLog.AddGameStep("用户%d选择了托管", user.UserID)
	uz := d.ext.getExtInfoByTurn(user.SeatNo)
	if uz != nil {
		uz.isTuoGuan = !uz.isTuoGuan
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[OnTuoGuan] userId[%v] isTuoGuan[%v]", user.UserID, uz.isTuoGuan)).Debug("託管 ")
	}
	tmp := d.ext.getExtInfoByTurn(d.ext.getNowSeat())
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[OnTuoGuan] uz.userSeat[%v] == tmp.userSeat[%v]", uz.userSeat, tmp.userSeat)).Debug("託管 ")
	if uz.userSeat == tmp.userSeat {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[OnTuoGuan][開始自動出牌] uz.userSeat[%v] == tmp.userSeat[%v]", uz.userSeat, tmp.userSeat)).Debug("託管 ")

		d.AutoChuPai()
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[responseTuoGuan] userId[%v] isTuoGuan[%v]", user.UserID, uz.isTuoGuan)).Debug("託管 ")
	d.responseTuoGuan(user.UserID, netproto.DDZ_ErrorCode_OK, "OK")
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[broadcastTuoGuanState] userId[%v] isTuoGuan[%v]", user.UserID, uz.isTuoGuan)).Debug("託管 ")
	d.broadcastTuoGuanState(user.SeatNo, uz.isTuoGuan)
}

func (d *DDZGameLogic) UserOffLine(gameuser *gamelogic.GameUser) {
	d.setExperienceCash(gameuser)
	uz := d.ext.getExtInfoByTurn(gameuser.SeatNo)
	if uz != nil {
		uz.isTuoGuan = true
	}
	d.broadcastTuoGuanState(gameuser.SeatNo, uz.isTuoGuan)
	d.FixedGameLogicBase.UserOffLine(gameuser)
}

func (d *DDZGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//叫地主
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_JiaoDiZhu) {
		d.onUserJiaoDiZhu(user, gmsg.MsgData.([]byte))
	}
	//抢地主
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_QiangDiZhu) {
		d.onQiangDiZhu(user, gmsg.MsgData.([]byte))
	}
	//明牌
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_MingPai) {
		d.onMingPai(user, gmsg.MsgData.([]byte))
	}
	//出牌
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_ChuPai) {
		d.onChuPai(user, gmsg.MsgData.([]byte))
	}
	//Pass
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_Pass) {
		d.OnPass(user, gmsg.MsgData.([]byte))
	}
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.DDZ_GameMessageClassID_TuoGuan) {
		d.OnTuoGuan(user, gmsg.MsgData.([]byte))
	}
	//动画等待
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		d.CheckUserAnimationCommited(user, gmsg, d.actionName_curr, nil)
	}
}

type changeUserInfo struct {
	userId int32
	weight int
	hands  int
	bigs   int
	isUser bool
}

//换最大牌
func (d *DDZGameLogic) ExchangeMaxCard(userid int32) {
	ur := d.GetUserByUserID(userid)
	if !ur.IsRobot() {
		mlog.Error("不是機器人不可以換牌")
		return
	}
	d.LogicBase.ExchangeMaxCard(userid)

	var maxRetUserID int32
	newCg := new(poker.CardGroup)
	poker := d.poker.(*poker.FixedPoker)
	userCardWeight := []changeUserInfo{} //玩家手牌的權重(越大代表好牌)
	var selfWeight changeUserInfo

	poker3 := poker.GetUserCards(0).Cards

	for _, user := range d.AllTableUser() {
		userPoker := poker.GetUserCards(user.UserID)
		newCg.Cards = append(userPoker.Cards, poker3...)

		dc := NewDefineCard(d)
		//重新拆牌
		userCard := dc.Combination(newCg)
		weight, hands, bigs := dc.CalWeight(userCard, userCard.SingleCards)

		w := changeUserInfo{
			userId: user.UserID,
			weight: int(weight),
			hands:  int(hands),
			bigs:   int(bigs),
			isUser: !user.IsRobot(),
		}
		if userid == user.UserID {
			selfWeight = w
		}
		userCardWeight = append(userCardWeight, w)

	}
	mlog.Debug("userCardWeight : %v", userCardWeight)

	bestuc := selfWeight
	for _, r := range userCardWeight {
		if userid == r.userId || !r.isUser || (r.weight < 19 && bestuc.weight > 18) {
			continue
		}
		bvalue := bestuc.hands - (bestuc.bigs * 2)
		rvalue := r.hands - (r.bigs * 2)
		if rvalue < bvalue {
			bestuc = r
		} else if rvalue == bvalue && r.weight > bestuc.weight {
			bestuc = r
		}
	}
	d.GameLog.AddGameStep("玩家牌型權重總分:%v", userCardWeight)
	if bestuc.userId == userid {
		return
	}
	maxRetUserID = bestuc.userId

	mlog.Debug("交换用户%d和用户%d牌型", userid, maxRetUserID)
	d.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	userPoker := poker.GetUserCards(userid)
	tmp := sortCardGroup(*userPoker)
	sort.Sort(&tmp)
	d.GameLog.AddGameStep("用户%v換牌前牌型 : %v", userid, userPoker)
	userPoker = poker.GetUserCards(maxRetUserID)
	tmp2 := sortCardGroup(*userPoker)
	sort.Sort(&tmp2)
	d.GameLog.AddGameStep("用户%v換牌前牌型 : %v", maxRetUserID, userPoker)

	d.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)

	poker.ExChange(userid, maxRetUserID)

	userPoker = poker.GetUserCards(userid)
	tmp3 := sortCardGroup(*userPoker)
	sort.Sort(&tmp3)
	d.GameLog.AddGameStep("用户%v還牌後牌型 : %v", userid, userPoker)
	mlog.Debug("用户%v換牌後牌型 : %v", userid, userPoker)

	userPoker = poker.GetUserCards(maxRetUserID)
	tmp4 := sortCardGroup(*userPoker)
	sort.Sort(&tmp4)
	d.GameLog.AddGameStep("用户%v換牌後牌型 : %v", maxRetUserID, userPoker)
	mlog.Debug("用户%v換牌後牌型 : %v", maxRetUserID, userPoker)

}

//根据牌型出牌
func (d *DDZGameLogic) outCardByType(outCardType cardType, userid int32, more int32) (a []int32) {
	defer func() {
		if a == nil || len(a) == 0 {
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 根据牌型出牌，选定的牌型为outCardType[%d],牌%v", userid, int(outCardType), a)).Debug("託管")
		}
		mlog.Debug("根据牌型出牌，选定的牌型为[%d],牌%v", int(outCardType), a)

	}()

	var result []int32

	tuoGuanUser := d.tuoGuanUser[userid]
	if (len(tuoGuanUser.ctg.Ct) <= int(2) && len(tuoGuanUser.ctg.SingleCards) == int(0)) ||
		(len(tuoGuanUser.ctg.Ct) == int(1) && len(tuoGuanUser.ctg.SingleCards) == int(1)) {
		//手牌只剩最大牌型時
		his := d.ext.his.GetBack()
		if his != nil && his.ct != cardTypeRocket {
			result = d.lastCardType(userid, more)
		}
	}

	if len(result) != int(0) {
		return result
	}

	//三带
	if outCardType == cardTypeThree { //三带
		return d.outThreeCards(userid, more)
	}
	//飞机
	if outCardType == cardTypePlaneAndWing {
		return d.outFeiji(userid, 5, more)
	}
	//连对
	if outCardType == cardTypeThreeStraight {
		return d.outLiandui(userid, more)
	}
	//顺子
	if outCardType == cardTypeStraight {
		return d.outShunzi(userid, more)
	}
	//对子
	if outCardType == cardTypeDouble {
		return d.outDuizi(userid, more)
	}
	//单牌
	if outCardType == cardTypeSingle {
		return d.outDanpai(userid, more)
	}
	//炸弹
	if outCardType == cardTypeBomb {
		return d.outZhadan(userid, more)
	}
	//火箭
	if outCardType == cardTypeRocket {
		return d.outRocket(userid)
	}

	return nil
}

//最後手牌只剩一種類型：火箭或炸彈就直接出牌
func (d *DDZGameLogic) lastCardType(userid int32, more int32) []int32 {
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("userid[%v] 获取火箭牌型%v , 最後手牌只剩一種類型：火箭或炸彈就直接出牌, more[%v]", userid, outCards, more)
		}
	}()
	tuoGuanUser := d.tuoGuanUser[userid]

	if tuoGuanUser == nil {
		mlog.Debug("获取机器人数据失败, 最後手牌只剩一種類型")
		return outCards
	}
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 要出最後手牌了...len(Ct)[%v] , len(SingleCards)[%v] ",
		userid, len(tuoGuanUser.ctg.Ct), len(tuoGuanUser.ctg.SingleCards))).Debug("託管")

	iCardType := 0
	for k, v := range tuoGuanUser.ctg.Ct {
		if v.byteCardType != cardTypeNone {
			iCardType = k
			break
		}
	}

	if tuoGuanUser.ctg.Ct[iCardType].byteCardType == cardTypeNone {
		if conf.TestMode {
			mlog.Fatal("error stop outCards[%v]", outCards)
		}
	}

	if tuoGuanUser.ctg.Ct[iCardType].byteCardType == cardTypeRocket ||
		tuoGuanUser.ctg.Ct[iCardType].byteCardType == cardTypeBomb {

		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 要出最後手牌了... Ct[%v]", userid, &tuoGuanUser.ctg.Ct)).Debug("託管")

		outCardTypeGroup := d.getMinCardType(tuoGuanUser.ctg.Ct[iCardType].byteCardType, 0, tuoGuanUser.ctg)
		if outCardTypeGroup == nil {
			mlog.Error("获取机器人数据失败，火箭 [%v]", outCards)
			return outCards
		}

		if tuoGuanUser.ctg.Ct[iCardType].byteCardType == cardTypeRocket {
			outCards = append(outCards, outCardTypeGroup.nMin)
			outCards = append(outCards, outCardTypeGroup.nMax)
		}

		his := d.ext.his.GetBack()
		if more <= poker.PokerNumberA && outCardTypeGroup.nMax > poker.PokerNumberTwo && outCardTypeGroup.nMax < more {
			mlog.Debug("robot[%v] 最後手牌炸彈比牌失敗, nMax[%v] more[%v]", userid, outCardTypeGroup.nMax, more)
			return outCards
		}

		user := d.GetUserByUserID(userid)
		if his.maxCard == poker.PokerNumberTwo && his.ct == cardTypeBomb && his.userSeat != user.SeatNo {
			mlog.Debug("最後出牌是2炸彈就不壓 maxCard[%v] hisCt[%v] ", his.maxCard, his.ct)
			return outCards
		}

		if tuoGuanUser.ctg.Ct[iCardType].byteCardType == cardTypeBomb {
			outCards = append(outCards, outCardTypeGroup.nMax)
			outCards = append(outCards, outCardTypeGroup.nMax)
			outCards = append(outCards, outCardTypeGroup.nMax)
			outCards = append(outCards, outCardTypeGroup.nMax)
		}

		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, outCardTypeGroup.byteCardType, 1)
	}

	return outCards
}

//火箭
func (d *DDZGameLogic) outRocket(userid int32) []int32 {
	outCardType := cardTypeRocket
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("获取火箭牌型%v", outCards)
	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("获取机器人数据失败，火箭")
		return outCards
	}
	outCardTypeGroup := d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，火箭")
		return outCards
	}
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMin)
	d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	d.reduceCardType(tuoGuanUser, cardTypeRocket, 1)
	return outCards
}

//出牌 三带
func (d *DDZGameLogic) outThreeCards(userid int32, more int32) []int32 {
	//三带
	outCardType := cardTypeThree
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("获取三带牌型%v , more[%v]", outCards, more)
		}

	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("机器人数据获取失败，三带")
		return outCards
	}

	if more == poker.PokerNumberTwo || more >= poker.PokerNumberSmallKing {
		return outCards
	}

	outCardTypeGroup := d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		if d.getHand(userid) <= 3 {
			outCardTypeGroup = d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
			if outCardTypeGroup == nil {
				//mlog.Error("机器人数据获取失败，三带")
				return outCards
			}
		} else {
			//mlog.Error("机器人数据获取失败，三带")
			return outCards
		}
	}

	if (outCardTypeGroup.nMax > more && more != poker.PokerNumberTwo) || outCardTypeGroup.nMax == poker.PokerNumberTwo {
		outCards = append(outCards, outCardTypeGroup.nMax)
		outCards = append(outCards, outCardTypeGroup.nMax)
		outCards = append(outCards, outCardTypeGroup.nMax)
		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, cardTypeThree, 1)
	}

	return outCards
}

// 获取手牌手数
func (lgc *DDZGameLogic) getHand(userid int32) int32 {
	totalHand := int32(0)
	ctg := lgc.tuoGuanUser[userid]
	if ctg == nil {
		return 20
	}
	for _, count := range ctg.cardTypesCount {
		if count <= 0 {
			continue
		}
		totalHand += count
	}
	return totalHand
}

//出牌 飞机
func (d *DDZGameLogic) outFeiji(userid, sum, more int32) []int32 {
	outCardType := cardTypePlaneAndWing
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("获取飞机牌型%v", outCards)
	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		// mlog.Error("机器人数据获取失败，飞机")
		return outCards
	}
	outCardTypeGroup := d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("机器人数据获取失败，飞机")
		return outCards
	}
	if sum < 2 {
		return outCards //
	}
	size := outCardTypeGroup.nMax - outCardTypeGroup.nMin + 1
	if size > sum {
		size = sum
	}
	for i := outCardTypeGroup.nMax; i > outCardTypeGroup.nMax-size; i-- {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}
	d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	d.reduceCardType(tuoGuanUser, cardTypePlaneAndWing, 1)
	//有size单牌
	bNBeedReset := false
	if tuoGuanUser.cardTypesCount[cardTypeSingle] >= size {
		singleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			min := int32(0)
			tuoGuanUser.ctg.SingleCards, min = d.minOneSingle(tuoGuanUser.ctg.SingleCards)
			if min <= 0 {
				if d.getHand(userid) <= size+3 {
					singOne := d.getMinCardType(cardTypeSingle, 0, tuoGuanUser.ctg)
					if singOne == nil {
						bNBeedReset = true
						break
					}
					singleList = append(singleList, singOne.nMax)
					d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
					d.eraseResult(tuoGuanUser.ctg, singOne)
					continue
				} else {
					bNBeedReset = true
					break
				}

			}
			singleList = append(singleList, min)
			d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)

		}
		if len(singleList) == int(size) {
			outCards = append(outCards, singleList...)
			return outCards
		}
	}
	//有两对的情况
	if tuoGuanUser.cardTypesCount[cardTypeDouble] >= size {
		doubleList := make([]int32, 0)
		for i := int32(0); i < size; i++ {
			double := new(cardTypeGroup)
			if d.getHand(userid) <= size+3 {
				double = d.getMinCardType(cardTypeDouble, 0, tuoGuanUser.ctg)
			} else {
				double = d.getMinCardType(cardTypeDouble, 2, tuoGuanUser.ctg)
			}
			if double == nil {
				// mlog.Error("机器人数据获取失败，飞机")
				bNBeedReset = true
				break
			}
			doubleList = append(doubleList, double.nMax)
			doubleList = append(doubleList, double.nMax)
			d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
			d.eraseResult(tuoGuanUser.ctg, double)
		}
		if len(doubleList) == int(size)*2 {
			outCards = append(outCards, doubleList...)
			return outCards
		}

	}
	//检测手牌时候除了三条之外还有1~2张2 或一个王
	up := d.gamePoker.GetUserPoker(userid)
	if up == nil {
		mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		return outCards
	}
	//有1对的情况
	if bNBeedReset {
		bNBeedReset = false
		d.recombination(userid, up.Cards)
		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, cardTypePlaneAndWing, 1)
	}
	if tuoGuanUser.cardTypesCount[cardTypeDouble]*2 >= size {
		doubleList := make([]int32, 0)
		x := size / 2 * 2
		if x < size {
			for i := int32(0); i < x/2; i++ {
				double := new(cardTypeGroup)
				if d.getHand(userid) <= size+3 {
					double = d.getMinCardType(cardTypeDouble, 0, tuoGuanUser.ctg)
				} else {
					double = d.getMinCardType(cardTypeDouble, 2, tuoGuanUser.ctg)
				}
				if double == nil {
					// mlog.Error("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
				d.eraseResult(tuoGuanUser.ctg, double)
			}
			double := new(cardTypeGroup)
			if d.getHand(userid) <= size+3 {
				double = d.getMinCardType(cardTypeDouble, 0, tuoGuanUser.ctg)
			} else {
				double = d.getMinCardType(cardTypeDouble, 2, tuoGuanUser.ctg)
			}

			if double == nil {
				// mlog.Error("机器人数据获取失败，飞机")
				bNBeedReset = true
			}
			doubleList = append(doubleList, double.nMax)
			d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
			d.eraseResult(tuoGuanUser.ctg, double)
		} else {
			for i := int32(0); i < size/2; i++ {
				double := new(cardTypeGroup)
				if d.getHand(userid) <= size+3 {
					double = d.getMinCardType(cardTypeDouble, 0, tuoGuanUser.ctg)
				} else {
					double = d.getMinCardType(cardTypeDouble, 2, tuoGuanUser.ctg)
				}
				if double == nil {
					// mlog.Error("机器人数据获取失败，飞机")
					bNBeedReset = true
					break
				}
				doubleList = append(doubleList, double.nMax)
				doubleList = append(doubleList, double.nMax)
				d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
				d.eraseResult(tuoGuanUser.ctg, double)
			}
		}
		if len(doubleList) == int(size) {
			outCards = append(outCards, doubleList...)
			return outCards
		}
	}

	return outCards
}

//连对
func (d *DDZGameLogic) outLiandui(userid int32, more int32) []int32 {
	outCardType := cardTypeThreeStraight
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("获取连对牌型%v", outCards)
	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("机器人数据获取失败，连对")
		return outCards
	}
	outCardTypeGroup := d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("机器人数据获取失败，连对")
		return outCards
	}

	for i := outCardTypeGroup.nMin; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
		outCards = append(outCards, i)
	}

	his := d.ext.his.GetBack()
	user := d.GetUserByUserID(userid)

	if user.SeatNo != his.userSeat && len(outCards) != len(his.cg.Cards) {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("出連對 失敗 outCards[%v] vs his.cg.Cards[%v]", outCards, his.cg.Cards)).Debug("託管")
		return nil
	}

	d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	d.reduceCardType(tuoGuanUser, cardTypeThreeStraight, 1)
	return outCards
}

//顺子
func (d *DDZGameLogic) outShunzi(userid int32, more int32) []int32 {
	outCardType := cardTypeStraight
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("获取顺子牌型%v", outCards)
	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("机器人数据获取失败，顺子")
		return outCards
	}
	outCardTypeGroup := d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("机器人数据获取失败，顺子")
		return outCards
	}
	his := d.ext.his.GetBack()

	var limit int32

	user := d.GetUserByUserID(userid)
	if user.SeatNo != his.userSeat {
		limit = (outCardTypeGroup.nMax - int32(len(his.cg.Cards)) + 1)
		if limit < outCardTypeGroup.nMin {
			return outCards
		}
	} else {
		limit = outCardTypeGroup.nMin
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("出順子 outCardTypeGroup.nMin[%v] nMax[%v] limit[%v]", outCardTypeGroup.nMin, outCardTypeGroup.nMax, limit)).Debug("託管")
	for i := limit; i <= outCardTypeGroup.nMax; i++ {
		outCards = append(outCards, i)
	}

	if user.SeatNo != his.userSeat {
		if len(outCards) < int(5) && len(outCards) != len(his.cg.Cards) {
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("出順子 失敗 outCards[%v] vs his.cg.Cards[%v]", outCards, his.cg.Cards)).Debug("託管")
			return nil
		}
	}

	d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	d.reduceCardType(tuoGuanUser, cardTypeStraight, 1)

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("出順子 %v", outCards)).Debug("託管")
	return outCards
}

//对子
func (d *DDZGameLogic) outDuizi(userid int32, more int32) []int32 {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] more[%v]", userid, more)).Debug("託管")
	outCardType := cardTypeDouble
	outCards := make([]int32, 0)
	defer func() {
		if len(outCards) > 0 {
			mlog.Debug("userid[%v] 获取对子牌型%v , more[%v]", userid, outCards, more)
		}
	}()

	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("获取机器人数据失败，对子")
		return outCards
	}

	his := d.ext.his.GetBack()
	user := d.GetUserByUserID(userid)

	if user.SeatNo != his.userSeat && more == 2 { //牌桌最後出2就不壓
		return outCards
	}

	outCardTypeGroup := new(cardTypeGroup)
	if more == 0 {
		outCardTypeGroup = d.getMinCardType(outCardType, 2, tuoGuanUser.ctg)
		if outCardTypeGroup != nil {
			more = 2
		}
	}
	var handDoubleCount int32 //手牌對子的數量
	for _, v := range tuoGuanUser.ctg.Ct {
		if v.byteCardType == cardTypeDouble {
			handDoubleCount++
		}
	}

	if his != nil && (user.SeatNo == his.userSeat || handDoubleCount <= int32(2)) {
		outCardTypeGroup = d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
		if outCardTypeGroup != nil && (outCardTypeGroup.nMax == int32(2) || outCardTypeGroup.nMax > 13) && outCardTypeGroup.nMax != more {
			outCards = append(outCards, outCardTypeGroup.nMax)
			outCards = append(outCards, outCardTypeGroup.nMax)
			d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
			d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 outCards[%v]", userid, outCards)).Debug("託管")
			return outCards
		}
	}

	outCardTypeGroup = d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] more[%v] allCardCtg[%v]",
		userid, more, tuoGuanUser.ctg.Ct)).Debug("託管")
	if outCardTypeGroup == nil {
		outCardTypeGroup = d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
		if outCardTypeGroup == nil || (outCardTypeGroup.nMax > int32(2) && outCardTypeGroup.nMax < more) {
			return outCards
		}
	}

	if outCardTypeGroup.nMax > more && more != poker.PokerNumberTwo || outCardTypeGroup.nMax == poker.PokerNumberTwo {
		outCards = append(outCards, outCardTypeGroup.nMax)
		outCards = append(outCards, outCardTypeGroup.nMax)
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] outCards[%v] more[%v]", userid, outCards, more)).Debug("託管")
		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
	} else if outCardTypeGroup != nil && outCardTypeGroup.nMax > more {
		outCards = append(outCards, outCardTypeGroup.nMax)
		outCards = append(outCards, outCardTypeGroup.nMax)
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-對子 userid[%v] outCards[%v] more[%v]", userid, outCards, more)).Debug("託管")
		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, cardTypeDouble, 1)
	}

	return outCards
}

//炸弹
func (d *DDZGameLogic) outZhadan(userid int32, more int32) []int32 {
	outCardType := cardTypeBomb
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("获取炸弹牌型%v", outCards)
	}()
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("获取机器人数据失败，炸弹")
		return outCards
	}

	if more == int32(2) {
		return outCards
	}

	outCardTypeGroup := d.getMinCardType(outCardType, more, tuoGuanUser.ctg)
	if outCardTypeGroup == nil {
		// mlog.Error("获取机器人数据失败，炸弹")
		return outCards
	}
	if outCardTypeGroup.nMax <= more {
		return outCards
	}

	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	outCards = append(outCards, outCardTypeGroup.nMax)
	d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
	d.reduceCardType(tuoGuanUser, cardTypeBomb, 1)
	return outCards
}

//单牌, 获取不小于more的最小牌型
func (d *DDZGameLogic) outDanpai(userid int32, more int32) []int32 {
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌 more[%v]", userid, more)).Debug("託管")
	outCardType := cardTypeSingle
	outCards := make([]int32, 0)
	defer func() {
		mlog.Debug("userid[%v]获取单牌牌型%v , more[%v]", userid, outCards, more)
	}()
	if more == poker.PokerNumberBigKing {
		return outCards
	}
	tuoGuanUser := d.tuoGuanUser[userid]
	if tuoGuanUser == nil {
		mlog.Error("获取託管人数据失败，单牌")
		return outCards
	}

	user := d.GetUserByUserID(userid)
	his := d.ext.his.GetBack()

	if d.getHand(userid) > 7 {
		if his != nil && user.SeatNo != his.userSeat && (more == int32(2) || more == int32(15) || more == int32(16)) {
			return outCards
		}
	}

	var outCardTypeGroup *cardTypeGroup
	//主動出牌優先時，手上牌小於３張就從特殊牌最小開始出牌(2,小王,大王)

	if his != nil && (user.SeatNo == his.userSeat || d.getHand(userid) <= 3) {
		outCardTypeGroup = d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)
		if outCardTypeGroup != nil && outCardTypeGroup.nMax != more && more < poker.PokerNumberSmallKing {
			outCards = append(outCards, outCardTypeGroup.nMax)
			d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
			d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 outCards[%v]", userid, outCards)).Debug("託管")
			return outCards
		}

		outCardTypeGroup = d.getMinCardType(outCardType, poker.PokerNumberA, tuoGuanUser.ctg)
		if outCardTypeGroup != nil && outCardTypeGroup.nMax != more {
			outCards = append(outCards, outCardTypeGroup.nMax)
			d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
			d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 outCards[%v]", userid, outCards)).Debug("託管")
			return outCards
		}
	}

	outCardTypeGroup = d.getMinCardSingleType(outCardType, more, tuoGuanUser.ctg.SingleCards)
	//mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌[%v] more[%v]", userid, outCardTypeGroup, more)).Debug("託管")
	if outCardTypeGroup == nil { //手牌點數沒有比牌桌大時, 出手牌的2(如果有)
		outCardTypeGroup = d.getMinCardType(outCardType, 0, tuoGuanUser.ctg)

		if outCardTypeGroup != nil && outCardTypeGroup.nMax == poker.PokerNumberTwo && more > poker.PokerNumberTwo && more < poker.PokerNumberSmallKing {
			outCards = append(outCards, outCardTypeGroup.nMax)
			d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
			d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
			mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌[%v] more[%v]", userid, outCards, more)).Debug("託管")
			return outCards
		}

	} else if more != poker.PokerNumberTwo {
		outCards = append(outCards, outCardTypeGroup.nMax)
		d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
		d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌[%v] more[%v]", userid, outCards, more)).Debug("託管")
		return outCards
	}

	if his == nil {
		outCardTypeGroup = d.getMinCardSingleType(outCardType, 0, tuoGuanUser.ctg.SingleCards)
		if outCardTypeGroup != nil {
			outCards = append(outCards, outCardTypeGroup.nMax)
			d.eraseResult(tuoGuanUser.ctg, outCardTypeGroup)
			d.reduceCardType(tuoGuanUser, cardTypeSingle, 1)
			//mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌 outCards[%v]", userid, outCards)).Debug("託管")
			return outCards
		}
	}
	mlog.NewDefaultLogger().SetData(fmt.Sprintf("userid[%v] 組合自動出牌中-單牌出牌[%v] more[%v]", userid, outCards, more)).Debug("託管")
	return outCards
}

// 最小单牌
func (lgc *DDZGameLogic) minOneSingle(a []int32) ([]int32, int32) {
	// mlog.Debug("单牌数量[%d],减少[1]", len(a))
	if len(a) <= 0 {
		return a, 0
	}
	a = lgc.sort(a)
	min := a[:1][0]
	a = a[1:]
	return a, min
}

func (d *DDZGameLogic) getMinCardSingleType(findCardType cardType, more int32, singleCards []int32) (p *cardTypeGroup) {
	minCt := new(cardTypeGroup)
	minCt.nMax = 20
	for _, val := range singleCards {
		if minCt.nMax > val && val > more {
			minCt.byteCardType = cardTypeSingle
			minCt.nMax = val
			minCt.nMin = val
			break
		}
	}

	if minCt.nMax == 20 {
		return nil
	}
	return minCt
}

//获取不小于more的最小牌型
func (lgc *DDZGameLogic) getMinCardType(findCardType cardType, more int32, ctg *CtGroup) (p *cardTypeGroup) {
	// defer func() {
	// 	if p != nil {
	// 		mlog.Debug("找到对应牌型[%d],more[%d],Info%v", int(findCardType), more, *p)
	// 	} else {
	// 		mlog.Debug("没有找到对应牌型[%d],more[%d]", int(findCardType), more)
	// 	}
	// }()
	minCt := new(cardTypeGroup)
	minCt.nMax = 20
	for _, ct := range ctg.Ct {
		//mlog.NewDefaultLogger().SetData(fmt.Sprintf("組合自動出牌中-获取不小于more的最小牌型 outCardTypeGroup[%v] max[%v] min[%v]", ct.byteCardType, ct.nMax, ct.nMin)).Debug("託管")
		if findCardType == ct.byteCardType {
			if minCt.nMax > ct.nMax && ct.nMax > more {
				minCt.byteCardType = ct.byteCardType
				minCt.nMax = ct.nMax
				minCt.nMin = ct.nMin
			}
		}
	}

	if minCt.nMax == 20 {
		return nil
	}
	return minCt
}

//移除对应牌组
func (lgc *DDZGameLogic) eraseResult(ctg *CtGroup, minCt *cardTypeGroup) {
	for _, ct := range ctg.Ct {
		if minCt.byteCardType == ct.byteCardType && minCt.nMax == ct.nMax && minCt.nMin == ct.nMin {
			ct.byteCardType = cardTypeNone
			return
		}
	}
	return
}

//是不是場上最大的牌型
func (lgc *DDZGameLogic) isBiggestCT(ct *cardTypeGroup, mineUserId int32) bool {
	if ct.byteCardType == cardTypeRocket {
		return true
	}
	for _, upk := range lgc.gamePoker.UserPokers {
		//場上的牌要排除自己
		if upk.User.UserID == mineUserId {
			continue
		}
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)
		mapCard := mapCardNumberCount(upk.Cards)
		mapCard = lgc.combinationKing(mapCard, result)
		mapCard = lgc.combinationCards(mapCard, cardTypeBomb, result)
		mapCard = lgc.combinationCards(mapCard, ct.byteCardType, result)

		for _, uct := range result.Ct {
			if uct.byteCardType == cardTypeRocket {
				//mlog.Warn("不是場上最大的牌型false cardTypeRocket")
				return false
			}
			if ct.byteCardType != cardTypeBomb && uct.byteCardType == cardTypeBomb {
				//mlog.Warn("不是場上最大的牌型false cardTypeBomb")
				return false
			}
			if uct.byteCardType == ct.byteCardType && cardMap[uct.nMax] > cardMap[ct.nMax] {
				//mlog.Debug("不是場上最大的牌型false 牌型[%v] cardMap[%v](max[%v]) c2[%v](max[%v])",
				//	uct.byteCardType, cardMap[uct.nMax], uct.nMax, cardMap[ct.nMax], ct.nMax)
				return false
			}
		}
	}
	return true
}

//是不是場上最大的單牌
func (lgc *DDZGameLogic) isBiggestSingle(card int32) bool {
	for _, upk := range lgc.gamePoker.UserPokers {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)
		mapCard := mapCardNumberCount(upk.Cards)
		for i, c := range mapCard {
			if c > 3 { //炸彈
				return false
			}
			if c > 0 {
				if cardMap[i] > cardMap[card] {
					return false
				}
			}
		}
	}
	return true
}

//取用戶手牌手數
func (lgc *DDZGameLogic) getUserPokerHands(userId int32) int32 {
	his := lgc.ext.his.GetBack()
	if his == nil {
		return 0
	}
	//seat := lgc.GetTable().GetSeatByNo(his.userSeat)
	//enemy := lgc.gamePoker.GetUserPoker(seat.SeatUser.UserID)
	//if enemy != nil {
	up := lgc.gamePoker.GetUserPoker(userId)
	if up != nil {
		dc := NewDefineCard(lgc)
		//重新拆牌
		userCard := dc.Combination(up.Cards)

		mlog.Debug("userid[%v] 取用戶手牌數[%v]", up.User.UserID, userCard.hands)
		return userCard.hands
	}
	mlog.Error("userid[%v]取不到用戶手牌 %v", userId, up)
	//}

	return 0
}
