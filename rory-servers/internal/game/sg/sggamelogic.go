package sg

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"math"

	proto "example.invalid/golang/protobuf/proto"
)

const (
	actionName_FaPai       = "fapai"       //发牌动画
	actionName_KaiPai      = "kaipai"      //开牌动画
	actionName_QiangZhuang = "qiangzhuang" //抢庄
	actionName_End         = "end"         //结束动画
)

const (
	QiangZhuang = iota
	ActionQiangZhuang
	Bet
	FaPai
	KaiPai
	GameResult
	WaitStartGame
)

type SGGameLogic struct {
	*gamelogic.DynamicGameLogicBase
	stateTimer        *event.TimerFunc
	poker             poker.Poker
	pokerType         int32         //1为随机牌，2为发牌器
	userData          []*SgUserData //用户数据
	state             int32         //游戏阶段
	userCount         int32         //正在游戏的人数
	betTotal          int32         //下注总额
	zhuangSeat        int32         //庄家座位号
	doneCount         int32         //已经进行操作的人总数
	betbase           int32         //下注的底分
	betArr            []int32       //下注倍数数组
	qiangZhuangTime   int32         //抢庄超时时间
	betTime           int32         //下注超时时间
	daSanGongOdd      int32         //大三公倍率
	xiaoSanGongOdd    int32         //小三公倍率
	hunSanGongOdd     int32         //混三公倍率
	teShuOdd          int32         //特殊牌倍率
	readyTime         int32         //准备时间
	actionQiangZhuang int32         //抢庄动画
	retAction         int32         //结算动画超时
	actionName_curr   string        //当前动画名称
	maxplayer         int32         //最多好多个人
	fapai             int32         //发完牌等待开牌时间
	kaipai            int32         //开牌动画等待时间
	resultmsg         *msg.Message  //结束消息给场景恢复用
	totalbet          int32         //总下注额
	isfapai           bool          //是否发牌，防止重复发牌
}

func NewSGGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *SGGameLogic {
	glc := new(SGGameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)
	glc.getsgConfig(cfgs)
	glc.state = WaitStartGame

	glc.RegisterAfterGameStartFunc(func() {
		//glc.sendSetting()
		glc.qiangZhuang()
	})

	glc.RegisterBeforeGameStartFunc(func() {
		glc.clearGameData() //清除数据
	})

	return glc
}

//获取配置文件
func (sg *SGGameLogic) getsgConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	sg.qiangZhuangTime = int32(scfgs.GetConfigInt("gameserver", "qiangzhuang"))
	sg.pokerType = int32(scfgs.GetConfigInt("gameserver", "pokertype"))
	sg.betTime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	sg.daSanGongOdd = int32(scfgs.GetConfigInt("gameserver", "dasangong"))
	sg.xiaoSanGongOdd = int32(scfgs.GetConfigInt("gameserver", "xiaosangong"))
	sg.hunSanGongOdd = int32(scfgs.GetConfigInt("gameserver", "hunsangong"))
	sg.teShuOdd = int32(scfgs.GetConfigInt("gameserver", "teshu"))
	sg.readyTime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	sg.retAction = int32(scfgs.GetConfigInt("gameserver", "retaction"))
	sg.maxplayer = int32(scfgs.GetConfigInt("gameserver", "maxplayer"))
	sg.fapai = int32(scfgs.GetConfigInt("gameserver", "fapai"))
	sg.kaipai = int32(scfgs.GetConfigInt("gameserver", "kaipai"))
	sg.betbase = int32(scfgs.GetConfigInt("gameserver", "basescore"))
	sg.actionQiangZhuang = int32(scfgs.GetConfigInt("gameserver", "actionqiangzhuang"))
	betstr := scfgs.GetConfigValue("gameserver", "bet")
	lastindex := 0
	for index, char := range betstr {
		if char == ',' {
			substr := betstr[lastindex:index]
			lastindex = index + 1
			bet, _ := strconv.ParseInt(substr, 10, 32)
			sg.betArr = append(sg.betArr, int32(bet))
		}
	}
}

//清理数据
func (sg *SGGameLogic) clearGameData() {
	sg.zhuangSeat = -1
	sg.stateTimer = nil
	sg.resultmsg = nil
	sg.totalbet = 0
	sg.isfapai = false
	sg.setUserData()

	if sg.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		sg.poker = poker.NewRandomPoker(false, false, 1)
		sg.poker = poker.NewFixedPoker(sg.AllTableUser(), sg.poker.(*poker.RandomPoker), 3)
	} else { //发牌器发牌
		sg.poker = poker.NewFilePoker(sg.Configs, sg.AllTableUser())
		sg.poker = poker.NewFixedFilePoker(sg.poker.(*poker.FilePoker))
	}
}

//通过用户座位号获取用户数据
func (sg *SGGameLogic) getUserBySeatNO(seatno int32) *SgUserData {
	for _, ud := range sg.userData {
		if seatno == ud.User.SeatNo {
			return ud
		}
	}

	return nil
}

//设置用户的牌
func (sg *SGGameLogic) setUserCards(user *gamelogic.GameUser, cards *poker.CardGroup) {
	ud := sg.getUserBySeatNO(user.SeatNo)
	if ud != nil {
		ud.userCards = poker.NewCardGroup(cards.Cards)
	}
}

//抢庄
func (sg *SGGameLogic) qiangZhuang() {
	sg.state = QiangZhuang
	sg.stateTimer = sg.TimerFuncManager.NewNamedTimer(sg.qiangZhuangTime+1, sg.qiangZhuangResult, "")
	sg.broadCastState()
}

//下注
func (sg *SGGameLogic) bet() {
	sg.state = Bet
	sg.stateTimer = sg.TimerFuncManager.NewNamedTimer(sg.betTime+1, sg.betResult, "")
	sg.broadCastState()
}

//发牌
func (sg *SGGameLogic) faPai() {
	//给所有用户发牌
	if sg.isfapai {
		mlog.Debug("已经发过牌了，还发牌？")
		return
	}

	sg.isfapai = true
	mlog.Debug("有多少人：%d", len(sg.AllGameUser()))
	mlog.Debug("实际的人数:%d", len(sg.userData))
	for _, ur := range sg.AllGameUser() {
		mlog.Debug("UserID=%v", ur.UserID)
		cards := sg.poker.Deal(ur.UserID, 3)
		sg.setUserCards(ur, cards)
		upk := poker.NewUserPoker(ur, cards)
		sg.GameLog.Fapai(ur.UserID, cards.String())
		//通知所有用户牌型
		//sg.BoardcastMessage(sg.getPokerMsg(upk))
		sg.SendMessageToUser(ur, sg.getPokerMsg(upk))
	}

	sg.stopStateTimer()
	//sg.stateTimer = sg.TimerFuncManager.NewNamedTimer(sg.fapai+1, sg.showCard, "")
	sg.actionName_curr = sg.RegisterAnimation(actionName_FaPai, sg.fapai, sg.startActionKaiPai)
}

//开始开牌动画
func (sg *SGGameLogic) startActionKaiPai() {
	sg.state = KaiPai
	sg.broadCastState()
	sg.actionName_curr = sg.RegisterAnimation(actionName_KaiPai, sg.kaipai, sg.showCard)
}

//开牌阶段
func (sg *SGGameLogic) showCard() {
	sg.sendGameResult()
	sg.actionName_curr = sg.RegisterAnimation(actionName_End, sg.retAction, sg.sendReadyTime)
}

//发送比牌结果
func (sg *SGGameLogic) sendGameResult() {
	sg.state = GameResult
	sg.broadCastState()

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_SGGameResultID)

	pdata := new(netproto.SG_GameResult)
	pdata.BankerSeatNO = proto.Int32(sg.zhuangSeat)
	var userArr []*netproto.SG_UserInfo

	userArr, uscoresRet, allwin, alllose, money, banker, bankerodd, bankertype := sg.calGameMoney()

	pdata.UserArr = userArr

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("紀錄遊戲結果 uscoresRet[%v]", uscoresRet)).Debug("結算")

	if alllose != 0 {
		pdata.TongChi = proto.Int32(alllose)
	} else if allwin != 0 {
		pdata.TongChi = proto.Int32(allwin)
	} else {
		pdata.TongChi = proto.Int32(0)
	}

	pdata.UserArr = append(pdata.UserArr, sg.getUserResult(banker, money, bankertype, bankerodd))
	pmsg.MsgData = pdata
	sg.resultmsg = pmsg
	sg.BoardcastMessage(pmsg)
	sg.GameLog.GameEnd()
}

func (sg *SGGameLogic) calGameMoney() (userArr []*netproto.SG_UserInfo, uscoresRet map[int32]int64, allwin int32, alllose int32, money int64, banker *SgUserData, bankerodd int32, bankertype int32) {
	uscoresRet = make(map[int32]int64)
	banker = sg.getUserBySeatNO(sg.zhuangSeat)
	if banker == nil {
		mlog.Debug("没找到庄家")
		return
	}

	bankerodd, bankertype = sg.getCardsOdd(banker.userCards) //取得莊家牌型
	allwin = int32(1)
	alllose = int32(2)

	for _, ud := range sg.userData {
		if ud.User.SeatNo == sg.zhuangSeat {
			continue
		}

		udodd, udtype := sg.getCardsOdd(ud.userCards) //取得閒家牌型
		ret := comPareCards(banker.userCards, ud.userCards)
		odd := udodd
		var udmoney int64 = 0
		if ret == 1 { //莊勝
			udmoney = int64(ud.bet * bankerodd)
			odd = bankerodd
			alllose = 0
		} else { //閒勝
			udmoney = int64(-ud.bet * udodd)
			allwin = 0
		}

		money += udmoney
		newmoney, _ := sg.GetUserWinAndTax(money)
		newudmoney, _ := sg.GetUserWinAndTax(-udmoney)
		uscoresRet[banker.User.UserID] = int64(newmoney)
		uscoresRet[ud.User.UserID] = int64(newudmoney)
		sg.GameLog.Context().UserUpdateScore(banker.User.UserID, newmoney)
		sg.GameLog.Context().UserUpdateScore(ud.User.UserID, newudmoney)

		userArr = append(userArr, sg.getUserResult(ud, -udmoney, udtype, odd))
	}

	return
}

//获取用户结果
func (sg *SGGameLogic) getUserResult(ud *SgUserData, money int64, udtype int32, odd int32) *netproto.SG_UserInfo {
	win, tax := sg.GetUserWinAndTax(money)
	if win > 0 {
		win -= ud.User.Score - ud.userMoney
	}
	sg.UpdateUserScoreFunc(ud.User.UserID, win, tax)

	var flowAmount, betAmount, payAmount int64
	flowAmount = util.ABSInt64(win)      //流水=派彩
	betAmount = util.ABSInt64(win) + tax //有效投注=輸贏+稅
	payAmount = betAmount + win          //返邏金額=有效投注+派彩
	step := fmt.Sprintf("用戶%v 輸贏%v 稅收%v, 流水%d, 有效投注 %d, 返還金額 %d", ud.User.UserID, win, tax, flowAmount, betAmount, payAmount)
	sg.GameLog.AddGameStep(step)
	fmt.Println(step)
	ud.User.FlowAmount += flowAmount
	ud.User.BetAmount += betAmount
	ud.User.PayAmount += payAmount

	udinfo := new(netproto.SG_UserInfo)
	udinfo.SeatNO = proto.Int32(ud.User.SeatNo)
	udinfo.PokerType = proto.Int32(udtype)
	udinfo.Money = proto.Int64(win)
	udinfo.Odd = proto.Int32(odd)
	for _, c := range ud.userCards.Cards {
		udinfo.Userpokergroup = append(udinfo.Userpokergroup, c.GetNetMessageData())
	}
	mlog.Debug("个人信息：%v", udinfo)
	return udinfo
}

//设置本轮用户已操作
func (sg *SGGameLogic) setUserData() {
	sg.userData = make([]*SgUserData, 0)
	sg.userCount = 0
	sg.doneCount = 0
	for _, ur := range sg.AllTableUser() {
		userdata := NewSgUerData(ur)
		sg.userData = append(sg.userData, userdata)
		sg.userCount++
	}
}

func (sg *SGGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	if !sg.CheckGameUserMessage(user) {
		sg.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	ud := sg.getUserBySeatNO(user.SeatNo)
	if ud == nil {
		return
	}

	if gmsg.SClassID == int32(netproto.SG_GameMessageClassID_SGQiangZhuangID) {
		if sg.state != QiangZhuang {
			return
		}

		sg.userQiangZhuang(ud, gmsg.MsgData.([]byte))

		return
	}

	if gmsg.SClassID == int32(netproto.SG_GameMessageClassID_BetID) {
		if sg.state != Bet {
			return
		}
		sg.userBet(ud, gmsg.MsgData.([]byte))
		return
	}

	sg.CheckUserAnimationCommited(user, gmsg, sg.actionName_curr, nil)
}

//获取扑克消息
func (sg *SGGameLogic) getPokerMsg(upk *poker.UserPoker) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)

	pmsg.MsgData = upk.GetNetMessageData()

	return pmsg
}

//获取阶段消息
func (sg *SGGameLogic) getSateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_StateID)

	pdata := new(netproto.SG_Sate)
	pdata.State = proto.Int32(sg.state)

	waitTime := int32(0)
	if sg.stateTimer != nil {
		waitTime = sg.stateTimer.WaitSeconds - sg.stateTimer.ElapsedSeconds() - 1
	}
	pdata.Time = proto.Int32(waitTime)
	pmsg.MsgData = pdata
	mlog.Debug("状态消息：%v", sg.state)
	return pmsg
}

//广播阶段
func (sg *SGGameLogic) broadCastState() {
	//等待游戏开始，不发游戏状态
	if sg.state == WaitStartGame {
		return
	}
	sg.doneCount = 0
	pmsg := sg.getSateMsg()
	sg.BoardcastMessage(pmsg)
}

//抢庄结果
func (sg *SGGameLogic) qiangZhuangResult() {
	sg.stopStateTimer()
	//找到抢庄中最有钱的人当庄。
	money := int64(0)
	rand.Seed(time.Now().Unix())
	for _, ud := range sg.userData {
		if ud.isQiangZhuang {
			if ud.userMoney > money {
				money = ud.userMoney
				sg.zhuangSeat = ud.User.SeatNo
			} else if ud.userMoney == money {
				i := rand.Int31n(2)
				if i == 1 {
					sg.zhuangSeat = ud.User.SeatNo
				}
			}
		} else {
			ud.CanQiangZhuangOpt = true
		}
	}

	//没人抢庄就随机庄
	if sg.zhuangSeat == -1 || sg.getUserBySeatNO(sg.zhuangSeat) == nil {
		rand.Seed(time.Now().Unix())
		sg.zhuangSeat = rand.Int31n(sg.maxplayer)
		for i := 0; i < 10; i++ {
			if sg.getUserBySeatNO(sg.zhuangSeat) != nil {
				break
			}

			sg.zhuangSeat = rand.Int31n(sg.maxplayer)
		}

		//10次都没找到庄家就用第一家作为庄家
		if sg.getUserBySeatNO(sg.zhuangSeat) == nil {
			sg.zhuangSeat = sg.userData[0].User.SeatNo
		}
	}

	//通知用户庄
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_BroadCastZhuangID)

	pdata := new(netproto.SG_SeatNO)
	pdata.SeatNO = proto.Int32(sg.zhuangSeat)
	pmsg.MsgData = pdata

	sg.GameLog.AddGameStep("产生庄家%d", sg.zhuangSeat)
	sg.BoardcastMessage(pmsg)

	sg.state = ActionQiangZhuang
	sg.broadCastState()
	sg.actionName_curr = sg.RegisterAnimation(actionName_QiangZhuang, sg.actionQiangZhuang, sg.bet)
}

//广播用户抢庄
func (sg *SGGameLogic) broadCastUserQiangZhuang(ud *SgUserData) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_SGUserQiangZhuang)

	pdata := new(netproto.SG_BroadCastQiangZhuang)
	pdata.SeatNO = proto.Int32(ud.User.SeatNo)
	pdata.QiangZhuang = proto.Bool(ud.isQiangZhuang)
	pmsg.MsgData = pdata

	sg.BoardcastMessage(pmsg)
}

//处理用户抢庄
func (sg *SGGameLogic) userQiangZhuang(ud *SgUserData, data []byte) {
	ust := &netproto.SG_QiangZhuang{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Debug("解析用户%v用户座位号错误%v", ud.User, err)
		return
	}

	if ud.isQaingZhuang() {
		mlog.Debug("已经抢过了")
		return
	}

	sg.GameLog.AddGameStep("用户%d抢庄%v", ud.User.UserID, *ust.QiangZhuang)
	ud.setUserQiangZhuang(*ust.QiangZhuang)
	sg.broadCastUserQiangZhuang(ud)
	sg.doneCount++
	if sg.doneCount >= sg.userCount {
		sg.qiangZhuangResult()
	}
}

//处理用户下注
func (sg *SGGameLogic) userBet(ud *SgUserData, data []byte) {
	ust := &netproto.SG_Bet{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v用户座位号错误%v", ud.User, err)
		return
	}
	*ust.Bet = int32(math.Abs(float64(*ust.Bet)))

	if ud.isBet {
		mlog.Debug("已经下过注了")
		return
	}
	//下注不包括庄家
	if sg.zhuangSeat == ud.User.SeatNo {
		return
	}
	bFind := false
	for i := 0; i < len(sg.betArr); i++ {
		if *ust.Bet == sg.betArr[i] {
			bFind = true
		}
	}

	if !bFind {
		mlog.Debug("没找到相应的下注倍数")
		return
	}

	sg.GameLog.AddGameStep("用户%d下注%d", ud.User.UserID, *ust.Bet*sg.betbase)
	ud.setUserBet(*ust.Bet, sg.betbase)
	sg.totalbet += ud.bet
	sg.sendUserBet(ud)
	sg.doneCount++
	//下注不包括庄家
	if sg.doneCount >= sg.userCount-1 {
		sg.betResult()
	}
}

//下注结果
func (sg *SGGameLogic) betResult() {
	sg.stopStateTimer()

	for _, ud := range sg.userData {
		if ud.bet == 0 && ud.User.SeatNo != sg.zhuangSeat {
			sg.GameLog.AddGameStep("用户下注超时%d自动下注%d", ud.User.UserID, sg.betArr[0]*sg.betbase)
			ud.setUserBet(sg.betArr[0], sg.betbase)
			sg.totalbet += ud.bet
			sg.sendUserBet(ud)
		}
	}

	//先改变状态再发牌
	sg.state = FaPai
	sg.broadCastState()
	sg.TimerFuncManager.NewNamedTimer(1, sg.faPai, "")
	//sg.faPai()
}

//获取设置消息
func (sg *SGGameLogic) getSettingMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_SettingID)

	pdata := new(netproto.SG_Setting)
	pdata.BetBase = proto.Int32(sg.betbase)
	pdata.ReadyTime = proto.Int32(sg.readyTime)
	pdata.BetOdd = append(pdata.BetOdd, sg.betArr...)
	pmsg.MsgData = pdata

	return pmsg
}

//发送设置消息
func (sg *SGGameLogic) sendSetting() {
	pmsg := sg.getSettingMsg()
	sg.BoardcastMessage(pmsg)
}

//拷贝牌
func (sg *SGGameLogic) copyPokerCards(Cards *poker.CardGroup) *poker.CardGroup {
	tmpCards := poker.NewEmptyCardGroup(len(Cards.Cards))
	for i := 0; i < len(Cards.Cards); i++ {
		tmpCards.Cards[i].Color = Cards.Cards[i].Color
		tmpCards.Cards[i].CardNumber = Cards.Cards[i].CardNumber
	}

	return tmpCards
}

//获取牌型的倍数
func (sg *SGGameLogic) getCardsOdd(Cards *poker.CardGroup) (int32, int32) {
	tmpCards := copyPokerCards(Cards)
	changeA(tmpCards)
	sortCards(tmpCards)

	cardsType := getSanGongType(tmpCards)

	if cardsType == DaSanGong {
		return sg.daSanGongOdd, cardsType
	}

	if cardsType == XiaoSanGong {
		return sg.xiaoSanGongOdd, cardsType
	}

	if cardsType == HunSanGong {
		return sg.hunSanGongOdd, cardsType
	}

	if cardsType >= 8 {
		return sg.teShuOdd, TeShu
	}

	return 1, DanPai
}

//广播玩家下注
func (sg *SGGameLogic) sendUserBet(ud *SgUserData) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_SGBroadCastBetID)

	pdata := new(netproto.SG_BroadCastBet)
	pdata.SeatNO = proto.Int32(ud.User.SeatNo)
	pdata.Bet = proto.Int32(ud.betodd)
	pdata.TotalBet = proto.Int32(sg.totalbet)
	pmsg.MsgData = pdata

	sg.BoardcastMessage(pmsg)
}

//用户进入
func (sg *SGGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	sg.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)
	setmsg := sg.getSettingMsg()
	sg.SendMessageToUser(gameuser, setmsg)

	//如果非游戏玩家，发送当前场景
	if !sg.IsGameUser(gameuser) {
		sg.sendGameScene(gameuser)
	}
}

//断线重连
func (sg *SGGameLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("发送场景数据")
	sg.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏场景
	sg.sendGameScene(gameuser)
}

//发送当前游戏场景
func (sg *SGGameLogic) sendGameScene(user *gamelogic.GameUser) {
	if sg.state == WaitStartGame {
		mlog.Debug("游戏结束所以不发场景了")
		return
	}
	setmsg := sg.getSettingMsg()
	sg.SendMessageToUser(user, setmsg)
	stateMsg := sg.getOptMsg(user)
	sg.SendMessageToUser(user, stateMsg)
	if sg.state == FaPai {
		for _, ud := range sg.userData {
			if ud.userCards == nil {
				break
			}
			upk := poker.NewUserPoker(ud.User, ud.userCards)
			//通知所有用户牌型
			sg.SendMessageToUser(user, sg.getPokerMsg(upk))
			//sg.BoardcastMessage(sg.getPokerMsg(upk))
			//sg.SendMessageToUser(ur, sg.getPokerMsg(upk))
		}
	}
}

//游戏结束
func (sg *SGGameLogic) sendReadyTime() {
	sg.state = WaitStartGame
	sg.GameEndFunc()
}

//获取操作消息
func (sg *SGGameLogic) getOptMsg(user *gamelogic.GameUser) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.SG_GameMessageClassID_SGGameSceneID)

	pdata := new(netproto.SG_Scene)
	pdata.State = proto.Int32(sg.state)
	pdata.ZhuangSeatNo = proto.Int32(sg.zhuangSeat)
	ud := sg.getUserBySeatNO(user.SeatNo)
	opt := false
	if ud != nil {
		if sg.state == QiangZhuang {
			opt = !ud.CanQiangZhuangOpt
		} else if sg.state == Bet {
			opt = (ud.bet == 0)
		}
	}
	pdata.CanOpt = proto.Bool(opt)

	waitTime := int32(0)
	if sg.stateTimer != nil {
		waitTime = sg.stateTimer.WaitSeconds - sg.stateTimer.ElapsedSeconds() - 1
	}
	pdata.Time = proto.Int32(waitTime)

	for _, ud := range sg.userData {
		udinfo := new(netproto.SG_UserSceneInfo)
		udinfo.SeatNO = proto.Int32(ud.User.SeatNo)
		qiangzhuangopt := int32(0)
		if ud.CanQiangZhuangOpt == true && ud.isQiangZhuang == false {
			qiangzhuangopt = 2
		} else if ud.CanQiangZhuangOpt == true {
			qiangzhuangopt = 1
		} else {
			qiangzhuangopt = 0
		}

		udinfo.QiangZhuang = proto.Int32(qiangzhuangopt)

		udinfo.BetOdd = proto.Int32(ud.betodd)
		pdata.UserArr = append(pdata.UserArr, udinfo)
	}

	pdata.TotalBet = proto.Int32(sg.totalbet)

	if sg.resultmsg != nil {
		pdata.GameRet = sg.resultmsg.MsgData.(*netproto.SG_GameResult)
	}

	pmsg.MsgData = pdata

	return pmsg
}

//停止当前时间管理器
func (sg *SGGameLogic) stopStateTimer() {
	if sg.stateTimer != nil {
		sg.stateTimer.Stop()
		sg.stateTimer = nil
	}
}

//换牌
func (sg *SGGameLogic) ExchangeMaxCard(userid int32) {
	sg.LogicBase.ExchangeMaxCard(userid)
	var maxRetUserID int32

	var maxcg *poker.CardGroup
	fpoker := sg.poker.(*poker.FixedPoker)
	for _, ur := range sg.AllTableUser() {
		if ur.IsRobot() && ur.UserID != userid {
			continue
		}

		cg := fpoker.GetUserCards(ur.UserID)
		if maxcg == nil {
			maxcg = cg
			maxRetUserID = ur.UserID
			continue
		}
		//比较牌型
		if comPareCards(maxcg, cg) != 1 {
			maxcg = cg
			maxRetUserID = ur.UserID
		}
	}

	if maxRetUserID == 0 {
		return
	}

	sg.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	fpoker.ExChange(userid, maxRetUserID)
}
