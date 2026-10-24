package cd21

import (
	"math"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"
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
	Wait         = 0 // 等待遊戲開始
	Bet          = 1 // 下注状态
	BetEnd       = 2 // 下注結束
	Cutpai       = 3 // 切牌階段
	BaoXiao      = 4 // 保險階段
	Fapai        = 5 // 發牌階段
	MainProgress = 6 // 主遊戲流程
	ZhuangFan    = 7 // 莊家翻牌
	End          = 8 // 結算結果
)

const (
	BetArea       = 1 // 用户下注區赢的钱
	UserBlackJack = 2 // 用户21點赢的钱
	TuoBao        = 3 // 用户保險或投降赢的钱
	ZhungBust     = 4 // 用户莊家爆牌區赢的钱
	MainGame      = 5 // 用户主遊戲輸贏
	TuoXian       = 6 // 用户投降輸錢
)

const (
	BetRecover = 1 // 用户下注區赢的钱
	AfterBet   = 2 // 用户21點赢的钱
)

// 下注區域
const (
	DUIZI     = 1 //對子
	ZHUANGBAO = 2 //莊家爆牌
	LUCKY20   = 3 //20點
	BIGGER13  = 4 //大於13點
	SMALLER13 = 5 //小於13點
	DIZHU     = 6 //玩家底注
	SPLIT     = 7 //玩家分牌
	BAOXIAN   = 8 //保險區域
	TUOXIAN   = 9 //投降區域
)

type card struct {
	color  int32
	number int32
}

type CD21gameLogic struct {
	*gamelogic.DynamicGameLogicBase
	opttimer           *event.TimerFunc
	optTime            int32                     //可操作时间
	bettotal           int64                     //下注总额
	bettotalrobot      int64                     //机器人下注总额
	bettypetotal       []*bettype                //下注分类汇总
	bettyperobot       []*bettype                //机器人下注分类汇总
	betvalue           *betconfig                //下注配置文件
	betodd             *betoddconfig             //下注配置文件
	userdata           []*cd21userdata           //用户数据
	cardresult         []*card                   //开奖的牌
	gamestate          int32                     //游戏状态
	bettime            int32                     //下注时间
	action_str         string                    //当前动画
	kaijiangactiontime int32                     //开奖动画超时
	endactiontime      int32                     //用户结算动画超时
	readytime          int32                     //用户准备时间
	betcount           []int32                   //下注数量
	poker              poker.Poker               //整副牌
	cutcardno          int                       //系統切牌的位置
	isCut              bool                      //這局是否要切牌
	cutcardtime        int32                     //切牌時間
	cutcardanimation   int32                     //切牌動畫
	firstuser          *cd21userdata             //起始玩家
	fapaiAction        int32                     //发牌动画超时
	currSeat           int32                     //当前操作的座位号
	nextSeat           int32                     //下個操作的座位号
	middleSettime      int32                     //中間結算动画超时
	betlimitmoney      int64                     //下注限制
	zhungCard          *poker.CardGroup          //莊家的牌
	zhungDian          int32                     //莊家牌點
	zhuangBust         bool                      //莊家爆牌
	zhungSeat          int32                     //莊家座位號
	first              bool                      //是否為第一局
	betAdd             []int32                   //加注数组
	roundidlimit       int32                     //判定局數
	defined            int32                     //係數定值
	judgeRange         []*judgeRange             //判定區域
	zhuangDianJudge    []*zhuangBuDian           //莊家補點結果
	dianJudge          []int32                   //莊家目前點數
	roundid            int32                     //本局局數(切牌用)
	redcard            int32                     //發紅牌
	basescore          int32                     //底分
	alreadyCut         bool                      //切牌了沒
	zhuangFanPaiTime   int32                     //翻牌時間
	tuobaotime         int32                     //託管倒數時間
	tuobaoflag         int                       //已操作投保玩家數
	retAction          int32                     //结算动画超时
	readyTime          int32                     //举手超时时间,单位秒
	uscoresRet         map[int32]int64           //結算結果
	syspoolData        *netproto.Cd21SysPoolData //水池資訊
	sumBetMoney        int64                     //總下注金額
	sumBetCount        int64                     //總下注次數
	poolAmount         int64                     //水池變化量
	zhuangbucontroll   int32                     //是否開啟莊家補牌表
	sysroundid         int32                     //本局局數(抓水線用)
	duiziconfig        int32                     //不開對子機率
	lucky20config      int32                     //不開Lucky20機率
	zhuangdianvariable int32                     //莊家點數補正值
	playervariable     int32                     //玩家機率補正值
}

//新建一个21点实例
func NewCD21gameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *CD21gameLogic {
	glc := new(CD21gameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)
	glc.getcd21config(cfgs)
	//載入取點設定
	glc.getConfigBysyspool(cfgs)
	glc.gamestate = Wait
	glc.RegisterAfterGameStartFunc(func() {
		//載入水池值(1)總注額/總壓注 (2)系統輸贏
		glc.syspoolcorrect()
		glc.initialfapai()
		glc.startAnimation()
	})
	glc.RegisterBeforeGameEndFunc(func() {

	})
	glc.RegisterBeforeGameStartFunc(func() {
		//初始化資料及獲得水池資訊
		glc.clearGameData()
	})

	glc.RegisterAfterGameEndFunc(func() {
		for _, u := range t.GetTableUser() {
			glc.UserLeave(u)
		}
	})
	return glc
}

//获取配置文件
func (cd21 *CD21gameLogic) getcd21config(cfgs *config.Group) {
	cd21.betvalue = newBetconfig(cfgs)
	cd21.betodd = newBetOddconfig(cfgs)
	//gameserver config
	scfgs := cfgs.GetConfig("gameserver")
	cd21.bettime = int32(scfgs.GetConfigInt("gameserver", "bettime"))
	cd21.cutcardtime = int32(scfgs.GetConfigInt("gameserver", "cutcardtime"))
	cd21.cutcardanimation = int32(scfgs.GetConfigInt("gameserver", "cutcardanimation"))
	cd21.fapaiAction = int32(scfgs.GetConfigInt("gameserver", "fapaiAction"))
	cd21.middleSettime = int32(scfgs.GetConfigInt("gameserver", "middleSettime"))
	cd21.betlimitmoney = int64(scfgs.GetConfigInt("gameserver", "betlimitmoney"))
	cd21.optTime = int32(scfgs.GetConfigInt("gameserver", "opttime"))
	cd21.basescore = int32(scfgs.GetConfigInt("gameserver", "basescore"))
	cd21.tuobaotime = int32(scfgs.GetConfigInt("gameserver", "tuobaotime"))
	cd21.retAction = int32(scfgs.GetConfigInt("gameserver", "retaction"))
	cd21.readyTime = int32(scfgs.GetConfigInt("gameserver", "readywaittime"))
	cd21.zhuangbucontroll = int32(scfgs.GetConfigInt("gameserver", "zhuangbucontroll"))

	betcountstr := scfgs.GetConfigValue("gameserver", "betcount")

	strarr := strings.Split(betcountstr, ",")
	for i := 0; i < len(strarr); i++ {
		count, _ := strconv.ParseInt(strarr[i], 10, 32)
		cd21.betcount = append(cd21.betcount, int32(count))
	}

	//syspool config
	syscfgs := cfgs.GetConfig("syspoolconf")
	cd21.roundidlimit = int32(syscfgs.GetConfigInt("property", "roundidlimit"))
	cd21.defined = int32(syscfgs.GetConfigInt("property", "defined"))
	cd21.duiziconfig = int32(syscfgs.GetConfigInt("property", "duiziconfig"))
	cd21.lucky20config = int32(syscfgs.GetConfigInt("property", "lucky20"))
}

//清理数据
func (cd21 *CD21gameLogic) clearGameData() {
	cd21.userdata = make([]*cd21userdata, 0)
	cd21.bettypetotal = make([]*bettype, 0)
	cd21.bettyperobot = make([]*bettype, 0)
	cd21.cardresult = make([]*card, 0)
	cd21.zhungCard = new(poker.CardGroup)
	cd21.bettotal = 0
	cd21.bettotalrobot = 0
	cd21.zhungSeat = 99
	cd21.first = true
	cd21.zhuangBust = false
	cd21.currSeat = 0
	cd21.nextSeat = 0
	cd21.zhungDian = 0
	cd21.firstuser = new(cd21userdata)
	cd21.alreadyCut = false
	cd21.zhuangFanPaiTime = 0
	cd21.tuobaoflag = 0
	cd21.sumBetMoney = 0
	cd21.sumBetCount = 0
	cd21.poolAmount = 0
	cd21.addUserData()
	cd21.uscoresRet = make(map[int32]int64)
	cd21.zhuangdianvariable = 0
	cd21.playervariable = 0
}

//檢查是否发牌/切牌
func (cd21 *CD21gameLogic) initialfapai() {
	//確認是否重新切牌或繼續發牌
	if cd21.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		mlog.System("tableno %d roundid %d", cd21.TableNo(), cd21.roundid)
		if cd21.roundid == 0 || cd21.redcard != 0 { //第一局或者識別到切牌記號則要重新洗牌
			cd21.redcard = 0
			cd21.roundid = 0
			cd21.isCut = true
			cd21.cutcardno = 6*52 - util.RndInt(260, 287)
			//正常发牌
			cd21.poker = poker.NewRandomPoker(false, false, 6) //初始化六副牌
			cd21.poker.Shuffle()
		} else {
			if cd21.poker.UserCardsLen(0) <= cd21.cutcardno {
				//用 redcard 做為切牌記號, 下局切牌
				redcard := util.RndInt(1, 5)
				cd21.redcard = int32(redcard)
			} else {
				cd21.poker.Shuffle()
			}
		}
		cd21.roundid += 1
		cd21.sysroundid += 1
	} else { //发牌器发牌
		cd21.poker = poker.NewFilePoker(cd21.Configs, cd21.AllTableUser())
		cd21.poker = poker.NewFixedFilePoker(cd21.poker.(*poker.FilePoker))
	}
}

//開始遊戲
func (cd21 *CD21gameLogic) startAnimation() {
	cd21.setState(Bet)
	cd21.startBet()
}

//开始下注
func (cd21 *CD21gameLogic) startBet() {
	cd21.opttimer = cd21.TimerFuncManager.NewTimer(cd21.bettime, cd21.betEnd)
	cd21.BoardcastMessage(cd21.getStateMsg())
	//遊戲開始強制下注
	for _, gameuser := range cd21.userdata {
		cd21.userHardBet(gameuser.user)
		//cd21.userTestHardBet1(gameuser.user) //測試用資料
		//cd21.userTestHardBet2(gameuser.user) //測試用資料
	}
}

//结束下注
func (cd21 *CD21gameLogic) betEnd() {
	cd21.setState(BetEnd)
	tableUser := cd21.AllTableUser()
	if tableUser == nil {
		cd21.GameEndFunc()
		return
	} else {
		for _, gameuser := range cd21.userdata {
			ud := cd21.getUserDataByID(gameuser.user.UserID)
			for _, bettype := range ud.userbet {
				if bettype.cd21type == DIZHU {
					ud.diZuBase = bettype.money
				}
			}
		}
	}

	if cd21.isCut {
		cd21.TimerFuncManager.NewTimer(int32(2), cd21.cutPai)
	} else {
		cd21.TimerFuncManager.NewTimer(int32(2), cd21.fapai)
	}
}

//设置状态
func (cd21 *CD21gameLogic) checkcut() {
	if !cd21.alreadyCut {
		cd21.hardCut()
	}
	cd21.TimerFuncManager.NewTimer(int32(4), cd21.fapai)
}

//设置状态
func (cd21 *CD21gameLogic) setState(state int32) {
	cd21.gamestate = state
}

//玩家強制下注
func (cd21 *CD21gameLogic) userHardBet(gameuser *gamelogic.GameUser) {
	cd21.GameLog.AddGameStep("用户%d, 遊戲開始 強制下注 %d", gameuser.UserID, cd21.basescore)
	user := cd21.getUserDataByID(gameuser.UserID)
	user.diZuBase = int64(cd21.basescore)
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetID)
	btdata := new(netproto.CD21_BetArr)
	betdata := new(netproto.CD21_Bet)
	betdata.BetType = proto.Int32(DIZHU)
	betdata.Count = proto.Int32(cd21.basescore)
	btdata.BetArr = append(btdata.BetArr, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	cd21.userBet(user.user, btmsg.MsgData.([]byte))
}

//測試用資料
func (cd21 *CD21gameLogic) userTestHardBet1(gameuser *gamelogic.GameUser) {
	user := cd21.getUserDataByID(gameuser.UserID)
	user.diZuBase = int64(cd21.basescore)
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetID)
	btdata := new(netproto.CD21_BetArr)
	betdata := new(netproto.CD21_Bet)
	betdata.BetType = proto.Int32(DUIZI)
	betdata.Count = proto.Int32(cd21.basescore)
	btdata.BetArr = append(btdata.BetArr, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	cd21.userBet(user.user, btmsg.MsgData.([]byte))
}

//測試用資料
func (cd21 *CD21gameLogic) userTestHardBet2(gameuser *gamelogic.GameUser) {
	user := cd21.getUserDataByID(gameuser.UserID)
	user.diZuBase = int64(cd21.basescore)
	btmsg := new(msg.Message)
	btmsg.BClassID = int32(netproto.MessageBClassID_Game)
	btmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetID)
	btdata := new(netproto.CD21_BetArr)
	betdata := new(netproto.CD21_Bet)
	betdata.BetType = proto.Int32(LUCKY20)
	betdata.Count = proto.Int32(cd21.basescore)
	btdata.BetArr = append(btdata.BetArr, betdata)
	btmsg.MsgData, _ = proto.Marshal(btdata)
	cd21.userBet(user.user, btmsg.MsgData.([]byte))
}

//设置状态
func (cd21 *CD21gameLogic) getStateMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StateID)

	pdata := new(netproto.CD21_Sate)
	pdata.State = proto.Int32(cd21.gamestate)
	pdata.Time = proto.Int32(cd21.getBetTime())
	pmsg.MsgData = pdata

	return pmsg
}

//獲取下注時間
func (cd21 *CD21gameLogic) getBetTime() int32 {
	time := int32(0)
	if cd21.opttimer != nil {
		time = cd21.opttimer.WaitSeconds - cd21.opttimer.ElapsedSeconds()

	}
	return time
}

//發牌
func (cd21 *CD21gameLogic) fapai() {
	cd21.setState(Fapai)
	//莊家自然發牌
	zhungCard := cd21.poker.Deal(0, 2)
	for _, dealcard := range zhungCard.Cards {
		dCard := new(poker.Card)
		dCard.CardNumber = dealcard.CardNumber
		dCard.Color = dealcard.Color
		cd21.zhungCard.Cards = append(cd21.zhungCard.Cards, dCard)
	}
	cd21.GameLog.Fapai(99999999, cd21.zhungCard.String())

	mlog.Debug("莊家牌:%v", cd21.zhungCard.String())
	if cd21.zhungCard.Cards[0].CardNumber <= 10 {
		cd21.zhungDian = cd21.zhungCard.Cards[0].CardNumber
	} else {
		if cd21.zhungCard.Cards[0].CardNumber == 14 {
			cd21.zhungDian = 11
		} else {
			cd21.zhungDian = 10
		}
	}

	//水線系統開啟時(使用發牌器時不操作)
	if cd21.GetServerConfig().PokerType == gamelogic.PokerTypeRandom && cd21.zhuangbucontroll == 1 {
		for _, userdata := range cd21.userdata {
			noBetBlock := true //沒下LUCKY20及DUIZI的不需要操控發牌
			for _, userbet := range userdata.userbet {
				if userbet.cd21type == LUCKY20 || userbet.cd21type == DUIZI {
					//自然發牌true, 已控制發完牌false
					noBetBlock = cd21.syspoolfapai(userdata)
					break
				}
			}
			if noBetBlock {
				//自然發牌
				natureCards := cd21.poker.Deal(userdata.user.UserID, 2)
				cd21.setUserCards(userdata, natureCards)
				mlog.Debug("cards=%v", natureCards)
				cd21.GameLog.AddGameStep("用户%d, 自然發牌: %v", userdata.user.UserID, natureCards.String())
			}

		}
	} else {
		//给所有用户发牌
		for _, ur := range cd21.userdata {
			cards := cd21.poker.Deal(ur.user.UserID, 2)
			cd21.setUserCards(ur, cards)
			mlog.Debug("cards=%v", cards)
			cd21.GameLog.AddGameStep("用户%d, 自然發牌: %v", ur.user.UserID, cards.String())
		}
	}

	//通知用户发牌
	cd21.BoardcastMessage(cd21.getFapaiMsg())

	//注册发牌动画
	cd21.TimerFuncManager.NewTimer(cd21.fapaiAction, cd21.middleSettlement)
}

//強制切牌
func (cd21 *CD21gameLogic) hardCut() {
	cd21.isCut = false
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastCutpaiRetID)
	btdata := new(netproto.CD21_BroadCastCutpaiRetID)
	btdata.SitNo = proto.Int32(cd21.firstuser.user.SeatNo)
	place := util.RndInt(53, 261)
	btdata.Place = proto.Int64(int64(place))
	gmsg.MsgData, _ = proto.Marshal(btdata)
	cd21.cutCard(cd21.firstuser, gmsg.MsgData.([]byte))
}

//通知切牌
func (cd21 *CD21gameLogic) cutPai() {
	//錢最多的玩家開始
	for i, user := range cd21.userdata {
		if i == 0 {
			cd21.firstuser = user
		} else {
			if user.money > cd21.firstuser.money {
				cd21.firstuser = user
			}
		}
	}
	//通知錢最多的切牌
	cd21.setState(Cutpai)
	cd21.BoardcastMessage(cd21.startCutMessage(cd21.firstuser))
	cd21.TimerFuncManager.NewTimer(cd21.cutcardtime, cd21.checkcut)
}

//切牌位置
func (cd21 *CD21gameLogic) cutCard(user *cd21userdata, data []byte) (bool, error) {
	//接收client的位置消息
	ust := &netproto.CD21_BroadCastCutpaiRetID{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v错误%v", user, err)
		return false, nil
	}
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastCutpaiID)
	pdata := new(netproto.CD21_BroadCastCutpaiRetID)
	pdata.SitNo = proto.Int32(*ust.SitNo)
	pdata.Place = proto.Int64(*ust.Place)
	pmsg.MsgData = pdata
	cd21.BoardcastMessage(pmsg)

	return true, nil
}

//用户进入
func (cd21 *CD21gameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	defer util.PanicErrStack()
	cd21.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)
	//如果非游戏玩家，发送当前场景
	cd21.sendScene(gameuser)

}

//发送场景
func (cd21 *CD21gameLogic) sendScene(gameuser *gamelogic.GameUser) {
	//发送设置
	if gameuser.IsRobot() {
		mlog.Debug("给机器人%d发送场景消息", gameuser.UserID)
	} else {
		mlog.Debug("给玩家%d发送场景消息", gameuser.UserID)
	}
	cd21.sendSetting(gameuser)
	ud := cd21.getUserDataByID(gameuser.UserID)
	if cd21.gamestate >= Bet && cd21.gamestate < End {
		//发送觀戰資訊
		cd21.BoardcastMessage(cd21.getWatchMsg())
	}
	//发送阶段
	cd21.SendMessageToUser(gameuser, cd21.getStateMsg())
	if cd21.gamestate == Bet {
		for _, user := range cd21.userdata {
			if user.user.UserID != gameuser.UserID {
				cd21.SendMessageToUser(gameuser, cd21.getRecoverBetTotalMsg(user, BetRecover))
			}
		}
		cd21.SendMessageToUser(gameuser, cd21.getRecoverWinLoseMsg(ud))
		cd21.sendUserBetInfo(gameuser)
	} else if cd21.gamestate > Bet && cd21.gamestate < End {
		if len(cd21.userdata) == 0 {
			return
		}
		for _, user := range cd21.userdata {
			cd21.SendMessageToUser(gameuser, cd21.getRecoverBetTotalMsg(user, AfterBet))
		}

		cd21.sendUserDiZuInfo(gameuser)

		if ud != nil {
			cd21.SendMessageToUser(gameuser, cd21.getRecoverWinLoseMsg(ud))
		}

	}
	if cd21.gamestate >= BaoXiao && cd21.gamestate < End {
		cd21.SendMessageToUser(gameuser, cd21.getPokerRecoverMsg())
	}

	if cd21.gamestate == BaoXiao {
		if cd21.zhungCard.Cards[0].CardNumber >= 10 { //跳保險階段
			msg := new(msg.Message)
			msg = cd21.getTuoBaoMsg(false, false)
			//檢查玩家錢夠不夠保險
			if ud != nil {
				if (ud.user.Score-ud.betmoney-ud.userDiZu) > (ud.diZuBase/2) && !ud.isBlackJack {
					msg = cd21.getTuoBaoMsg(false, true)
				}
			}

			cd21.SendMessageToUser(gameuser, msg)
		}
	}

	if cd21.gamestate == MainProgress {
		pmsg := new(msg.Message)
		pmsg.BClassID = int32(netproto.MessageBClassID_Game)
		pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21GameSceneID)
		pdata := new(netproto.CD21_UserInfoArr)
		pdata.CurrOptSeatno = proto.Int32(cd21.currSeat)
		opt := int32(0)
		if cd21.opttimer != nil {
			opt = cd21.opttimer.WaitSeconds - cd21.opttimer.ElapsedSeconds() - 1
		}
		pdata.OptTime = proto.Int32(opt)
		pdata.TotalTime = proto.Int32(cd21.optTime)
		pmsg.MsgData = pdata
		cd21.SendMessageToUser(gameuser, pmsg)
		ud := cd21.getUserDataByID(gameuser.UserID)
		if gameuser.SeatNo == cd21.currSeat {
			if ud.cansplitOpt() { //分牌處理
				cd21.sendSplitOptMsg(true)
			} else {
				cd21.sendOptMsg(true)
			}
		}
	}
}

//切牌信息
func (cd21 *CD21gameLogic) startCutMessage(user *cd21userdata) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21CutpaiID)
	cutMsg := new(netproto.CD21_BroadCastCutpaiID)
	cutMsg.SitNo = proto.Int32(user.user.SeatNo)
	pmsg.MsgData = cutMsg
	return pmsg
}

//用户消息处理逻辑
func (cd21 *CD21gameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	defer util.PanicErrStack()
	if !cd21.CheckGameUserMessage(user) {
		cd21.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	//Client確認時間
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21CheckTimeUpID) {
		cd21.SendMessageToUser(user, cd21.getStateMsg())
		return
	}

	//用戶下注
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21BetID) {
		if cd21.gamestate == Bet {
			cd21.userBet(user, gmsg.MsgData.([]byte))
		}
		return
	}

	//没有找到相关玩家数据
	urd := cd21.getUserBySeatNo(user.SeatNo)
	if urd == nil {
		mlog.Debug("用户未找到%v", user)
		return
	}

	//用户不可操作
	if !urd.canOpt() {
		mlog.Debug("用户已经停牌或投降%v", urd)
		return
	}

	//用戶切牌
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21BroadCastCutpaiRetID) {
		if cd21.gamestate == Cutpai && urd.user.UserID == cd21.firstuser.user.UserID {
			cd21.alreadyCut = true
			cd21.cutCard(urd, gmsg.MsgData.([]byte))
			cd21.isCut = false
		}
		return
	}

	////用户投降(1.10移除投降)
	//if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserSurrender) {
	//	cd21.tuobaoflag += 1
	//	cd21.surrender(urd)
	//	return
	//}

	//用户保險
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserInsurence) {
		cd21.tuobaoflag += 1
		cd21.insurance(urd)
		return
	}

	//用户取消
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserCancel) {
		cd21.tuobaoflag += 1
		cd21.passhand(urd)
		return
	}

	if user.SeatNo != cd21.currSeat {
		mlog.Error("用户座位号非当前操作用户座位号%v", user)
		return
	}

	//用户停牌
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserStand) {
		cd21.userStand(urd)
		return
	}

	//用户要牌
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserHit) {
		if !urd.isDouble {
			cd21.userHit(urd)
		}
		return
	}

	//用户分牌後要牌
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserSplitHit) {
		if urd.isSplit {
			cd21.userSplitHit(urd)
		}
		return
	}

	//用户分牌後停牌
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserSplitStand) {
		if urd.isSplit {
			cd21.userSplitStand(urd)
		}
		return
	}

	//用户分牌，再扣一倍底注
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserSplit) {
		//check spilt
		if urd.userCards.Cards[0].CardNumber == urd.userCards.Cards[1].CardNumber && len(urd.userCards.Cards) == 2 {
			cd21.splitCard(urd)
		}
		return
	}

	//用户加倍
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserDouble) {
		if len(urd.userCards.Cards) > 2 {
			return
		}
		cd21.doubleCard(urd)
		return
	}

	//用户分牌後加倍
	if gmsg.SClassID == int32(netproto.CD21_GameMessageClassID_CD21UserSplitDouble) {
		if urd.isSplit {
			cd21.splitdoubleCard(urd)
		}
		return
	}
}

//获取玩家数据通过座位号
func (cd21 *CD21gameLogic) getUserBySeatNo(seatno int32) *cd21userdata {
	for _, urd := range cd21.userdata {
		// 這裡很坑, TableNo 要一起檢查才行
		if cd21.TableNo() == urd.user.TableNo && urd.user.SeatNo == seatno {
			return urd
		}
	}

	return nil
}

//获取玩家数据通过用户ID
func (cd21 *CD21gameLogic) getUserDataByID(userid int32) *cd21userdata {
	for _, urd := range cd21.userdata {
		if urd.user.UserID == userid {
			return urd
		}
	}

	return nil
}

//将用户的牌放到用户数据里面去
func (cd21 *CD21gameLogic) setUserCards(ud *cd21userdata, cards *poker.CardGroup) {
	ud.userCards = poker.NewCardGroup(cards.Cards)
}

//将用户的牌放到用户数据里面去
func (cd21 *CD21gameLogic) setUserSplitCards(ud *cd21userdata, cards *poker.Card) {
	ud.splitCards.Cards = append(ud.splitCards.Cards, cards)
}

//获取发牌消息
func (cd21 *CD21gameLogic) getFapaiMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastFapaiID)
	fdata := new(netproto.CD21_Fapai)

	//各玩家牌型根據座位號塞入資訊
	for _, u := range cd21.userdata {
		pdata := new(netproto.CD21CardContent)
		pdata.SitNo = proto.Int32(u.user.SeatNo)
		for _, cardGroup := range u.userCards.Cards {
			npk := new(netproto.PokerCard)
			npk.Color = proto.Int32(cardGroup.Color)
			npk.CardNumber = proto.Int32(cardGroup.CardNumber)
			pdata.SitPoker = append(pdata.SitPoker, npk)
		}
		pdata.CardNumber = proto.Int32(cd21.getCardPoint(u.userCards))
		fdata.Content = append(fdata.Content, pdata)
	}

	//塞入莊家牌型
	zdata := new(netproto.CD21CardContent)
	zdata.SitNo = proto.Int32(99)
	for i := 0; i < 2; i++ {
		if i == 0 {
			zpk := new(netproto.PokerCard)
			zpk.CardNumber = proto.Int32(cd21.zhungCard.Cards[0].CardNumber)
			zpk.Color = proto.Int32(cd21.zhungCard.Cards[0].Color)
			zdata.SitPoker = append(zdata.SitPoker, zpk)
		} else {
			zpk := new(netproto.PokerCard)
			zpk.CardNumber = proto.Int32(0)
			zpk.Color = proto.Int32(0)
			zdata.SitPoker = append(zdata.SitPoker, zpk)
		}
	}
	zdata.CardNumber = proto.Int32(cd21.getCardPoint(cd21.zhungCard))
	fdata.Content = append(fdata.Content, zdata)

	fdata.LastTurn = proto.Int32(cd21.redcard)

	pmsg.MsgData = fdata
	return pmsg
}

//获取发牌消息
func (cd21 *CD21gameLogic) getFaZhuangPaiMsg(card *poker.CardGroup) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21ZhuangJiaYaoPaiID)
	fdata := new(netproto.CD21_ZhuangYaoPai)
	//塞入莊家牌型
	fdata.SitNo = proto.Int32(99)
	for _, cardGroup := range cd21.zhungCard.Cards {
		npk := new(netproto.PokerCard)
		npk.Color = proto.Int32(cardGroup.Color)
		npk.CardNumber = proto.Int32(cardGroup.CardNumber)
		fdata.SitPoker = append(fdata.SitPoker, npk)
	}
	fdata.CardTotal = proto.Int32(cd21.getCardPoint(cd21.zhungCard))
	pmsg.MsgData = fdata
	return pmsg
}

//获取翻牌消息
func (cd21 *CD21gameLogic) getFanPaiMsg(card *poker.CardGroup) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastFanPaiID)
	fdata := new(netproto.CD21_FanPai)
	//塞入莊家要翻的牌
	npk := new(netproto.PokerCard)
	npk.Color = proto.Int32(card.Cards[1].Color)
	npk.CardNumber = proto.Int32(card.Cards[1].CardNumber)
	fdata.SitPoker = append(fdata.SitPoker, npk)
	fdata.CardTotal = proto.Int32(cd21.getCardPoint(cd21.zhungCard))
	pmsg.MsgData = fdata
	return pmsg
}

func (cd21 *CD21gameLogic) getNextSeatUser() *cd21userdata {
	for _, seat := range cd21.GetTable().Seats {
		if seat.SeatNo >= cd21.nextSeat {
			if ud := cd21.getUserBySeatNo(seat.SeatNo); ud != nil {
				return ud
			}
		}
	}
	return nil
}

//获取下一个可操作的用户
func (cd21 *CD21gameLogic) waitOpt() {
	//停止保險倒數計時器
	cd21.setState(MainProgress)
	cd21.BoardcastMessage(cd21.getTimesUpStateMsg(0))

	curUrd := cd21.getUserBySeatNo(cd21.currSeat)
	if curUrd == nil { //代表沒有其他人需要操作了, 發莊家牌
		cd21.fanZhuangPoker()
		return
	} else {
		if curUrd.canOpt() && !curUrd.isBlackJack { //當前操作者還能操作
			cd21.nextSeat = cd21.currSeat
		} else {
			cd21.nextSeat = cd21.currSeat + 1
		}
	}

	//下一個操作者
	urd := cd21.getNextSeatUser()
	if urd == nil { //代表沒有其他人需要操作了, 發莊家牌
		cd21.fanZhuangPoker()
		return
	} else {
		cd21.currSeat = urd.user.SeatNo
		mlog.Debug("此次操作的座位号为：%v", cd21.currSeat)

		if urd.cansplitOpt() { //分牌處理
			cd21.sendSplitOptMsg(false)
			return
		} else if urd.canOpt() && !urd.isBlackJack {
			cd21.sendOptMsg(false)
			return
		} else {
			cd21.waitOpt()
			return
		}
	}
}

//中間結算
func (cd21 *CD21gameLogic) middleSettlement() {
	countBJ := 0
	cd21.sumBetCount += int64(len(cd21.userdata))
	for _, ud := range cd21.userdata {
		//给每个玩家发送他们的开奖输赢
		cd21.sendUserResult(ud)
		cd21.checkUserBJ(ud)
		if ud.isBlackJack {
			countBJ++
		}
	}

	//若要表演BJ，要多等一秒
	if countBJ > 0 {
		cd21.TimerFuncManager.NewTimer(int32(1), func() {
			cd21.BoardcastMessage(cd21.getUserWinMsg(BetArea))
		})
	} else {
		cd21.BoardcastMessage(cd21.getUserWinMsg(BetArea))
	}

	if countBJ > 0 && countBJ == len(cd21.userdata) {
		cd21.TimerFuncManager.NewTimer(int32(1), func() {
			msg := cd21.getFanPaiMsg(cd21.zhungCard)
			cd21.BoardcastMessage(msg)
			cd21.finalSettlement()
		})
		return
	}

	//是否需要投保階段
	insurence := false

	//投保流程(確認莊家是否有KQJA)
	biggerNo := cd21.zhungCard.Cards[0].CardNumber
	if biggerNo >= 10 { //跳保險階段
		NotEnoughFlag := int(0)
		//檢查玩家錢夠不夠保險
		for _, user := range cd21.userdata {
			if (user.user.Score - user.betmoney - user.userDiZu) < (user.diZuBase / 2) {
				NotEnoughFlag++
			}
		}
		if NotEnoughFlag == len(cd21.userdata) {
			insurence = false
		} else {
			insurence = true
		}
	}
	if insurence {
		cd21.TimerFuncManager.NewTimer(int32(1), func() {
			cd21.isTuoBao() //1.10移除投降
		})
	} else {
		cd21.waitOpt()
	}
}

func (cd21 *CD21gameLogic) sendUserResult(ud *cd21userdata) {
	var tax int64
	win := int64(0) //下注區塊贏錢(包含底注)
	//betmoney := int64(0) //下注底注
	for _, bettype := range ud.userbet {
		//計算玩家押注內容是否跟牌型相符
		if bettype.cd21type == DUIZI {
			cd21.sumBetCount += 1
			ud.userwinblock = append(ud.userwinblock, DUIZI)
			//check DUIZI
			if ud.userCards.Cards[0].CardNumber == ud.userCards.Cards[1].CardNumber {
				//玩家贏錢
				ud.winDUIZI = (cd21.betodd.getOddByType(rateDUIZI) - 1) * bettype.money
				_, taxDuiZi := cd21.GetUserWinAndTax(ud.winDUIZI)
				win += cd21.betodd.getOddByType(rateDUIZI)*bettype.money - taxDuiZi
				tax += taxDuiZi
			}
		} else if bettype.cd21type == LUCKY20 {
			cd21.sumBetCount += 1
			ud.userwinblock = append(ud.userwinblock, LUCKY20)
			if cd21.getCardPoint(ud.userCards) == 20 {
				if ud.userCards.Cards[0].CardNumber == 10 && ud.userCards.Cards[1].CardNumber == 10 && ud.userCards.Cards[0].Color == ud.userCards.Cards[1].Color {
					//rateLUCKY10TH:10同花
					ud.winLUCKY20 = (cd21.betodd.getOddByType(rateLUCKY10TH) - 1) * bettype.money
					_, taxLucky20 := cd21.GetUserWinAndTax(ud.winLUCKY20)
					tax += taxLucky20
					win += cd21.betodd.getOddByType(rateLUCKY10TH)*bettype.money - taxLucky20
					//當兩張牌都大於10小於K
				} else if ud.userCards.Cards[0].CardNumber >= 10 && ud.userCards.Cards[0].CardNumber <= 13 && ud.userCards.Cards[1].CardNumber >= 10 && ud.userCards.Cards[1].CardNumber <= 13 {
					if ud.userCards.Cards[0].CardNumber == ud.userCards.Cards[1].CardNumber && ud.userCards.Cards[0].Color == ud.userCards.Cards[1].Color {
						//rateLUCKYJQKTH:JQK同花對子(10同花對子在一開始就被挑出來了)
						ud.winLUCKY20 = (cd21.betodd.getOddByType(rateLUCKYJQKTH) - 1) * bettype.money
						_, taxLucky20 := cd21.GetUserWinAndTax(ud.winLUCKY20)
						tax += taxLucky20
						win += cd21.betodd.getOddByType(rateLUCKYJQKTH)*bettype.money - taxLucky20
					} else if ud.userCards.Cards[0].CardNumber == ud.userCards.Cards[1].CardNumber {
						//rateLUCKY10DZ:10JQK對子
						ud.winLUCKY20 = (cd21.betodd.getOddByType(rateLUCKY10DZ) - 1) * bettype.money
						_, taxLucky20 := cd21.GetUserWinAndTax(ud.winLUCKY20)
						tax += taxLucky20
						win += cd21.betodd.getOddByType(rateLUCKY10DZ)*bettype.money - taxLucky20
					} else {
						//rateLUCKY20:任意20點
						ud.winLUCKY20 = (cd21.betodd.getOddByType(rateLUCKY20) - 1) * bettype.money
						_, taxLucky20 := cd21.GetUserWinAndTax(ud.winLUCKY20)
						tax += taxLucky20
						win += cd21.betodd.getOddByType(rateLUCKY20)*bettype.money - taxLucky20
					}
				} else {
					//rateLUCKY20:任意20點(包括A9)
					ud.winLUCKY20 = (cd21.betodd.getOddByType(rateLUCKY20) - 1) * bettype.money
					_, taxLucky20 := cd21.GetUserWinAndTax(ud.winLUCKY20)
					tax += taxLucky20
					win += cd21.betodd.getOddByType(rateLUCKY20)*bettype.money - taxLucky20
				}
			}
		} else if bettype.cd21type == BIGGER13 {
			cd21.sumBetCount += 1
			ud.userwinblock = append(ud.userwinblock, BIGGER13)
			if cd21.getCardPoint(ud.userCards) > 13 {
				ud.winBIGGER13 += (cd21.betodd.getOddByType(rateBIGGER13) - 1) * bettype.money
				_, taxBigger13 := cd21.GetUserWinAndTax(ud.winBIGGER13)
				tax += taxBigger13
				win += cd21.betodd.getOddByType(rateBIGGER13)*bettype.money - taxBigger13
			}
		} else if bettype.cd21type == SMALLER13 {
			cd21.sumBetCount += 1
			ud.userwinblock = append(ud.userwinblock, SMALLER13)
			if cd21.getCardPoint(ud.userCards) < 13 {
				ud.winSMALLER13 += (cd21.betodd.getOddByType(rateSMALLER13) - 1) * bettype.money
				_, taxSmaller13 := cd21.GetUserWinAndTax(ud.winSMALLER13)
				tax += taxSmaller13
				win += cd21.betodd.getOddByType(rateSMALLER13)*bettype.money - taxSmaller13
			}
		}
	}
	ud.win = win
	ud.tax += tax
	cd21.GameLog.AddGameStep("中間結算, 用户%d, 下注區赢的钱 %d", ud.user.UserID, ud.win)

}

func (cd21 *CD21gameLogic) checkUserBJ(ud *cd21userdata) {
	var income, tax int64
	// check blackjack
	if cd21.getCardPoint(ud.userCards) == 21 {
		//玩家BJ
		ud.isBlackJack = true
		//BlackJack贏1.5倍
		income, tax = cd21.GetUserWinAndTax(int64(float64(ud.diZuBase) * 1.5))
		ud.winBLACKJACK = income
		income += ud.diZuBase
		ud.win1 = income
		ud.betmoney += ud.userDiZu
		ud.tax += tax
		ud.userDiZu = 0
		ud.diZuBase = 0
		mlog.Debug("玩家%v已经是21点", ud.user.UserID)
		cd21.GameLog.AddGameStep("中間結算, 用户21點赢的钱 %d", ud.win1)
	}
}

//最終結算
func (cd21 *CD21gameLogic) finalSettlement() {

	if cd21.gamestate == End {
		mlog.Warn("已結算過")
		return
	}

	cd21.setState(End)
	//發送下注user結算結果
	for _, ud := range cd21.userdata {
		//用户莊家爆牌區赢的钱
		cd21.sendUsersetResult(ud)
		//用户主遊戲輸贏
		cd21.compareWithBanker(ud)
		//分牌的也要比牌
		if ud.isSplit {
			cd21.splitWithBanker(ud)
		}
		cd21.finalScore(ud)
	}

	//將記錄丟進水池
	cd21.poolWater()
	cd21.BoardcastMessage(cd21.getUserWinMsg(ZhungBust))
	cd21.BoardcastMessage(cd21.getUserWinMsg(MainGame))
	cd21.GameLog.GameEnd()
	cd21.TimerFuncManager.NewTimer(cd21.retAction, cd21.sendReadyTime)
}

//发送设置
func (cd21 *CD21gameLogic) sendSetting(gameuser *gamelogic.GameUser) {
	mlog.NewDefaultLogger().Debug("sendSetting")
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21SettingID)

	pdata := new(netproto.CD21_Setting)
	pdata.BetCount = append(pdata.BetCount, cd21.betcount...)
	pdata.BetMoney = proto.Int32(int32(cd21.betvalue.getMoneyByType(1)))
	pdata.BetMoneyLimit = proto.Int64(cd21.betvalue.getMaxMoneyByType(1))

	pmsg.MsgData = pdata
	cd21.SendMessageToUser(gameuser, pmsg)
}

//获取总下注的消息
func (cd21 *CD21gameLogic) getBetTotalMsg(userid int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastBetID)

	pdata := new(netproto.CD21_BroadCastBetArr)
	pdata.UserID = proto.Int32(userid)
	for _, v := range cd21.betvalue.alltype {
		betinfo := new(netproto.CD21_BroadCastBet)
		betinfo.Type = proto.Int32(v)
		betinfo.Money = proto.Int64(cd21.getBetMoneyByType(v))
		pdata.BroadCastBetArr = append(pdata.BroadCastBetArr, betinfo)
	}

	pmsg.MsgData = pdata
	return pmsg
}

//获取下注的总钱数，通过游戏类型
func (cd21 *CD21gameLogic) getBetMoneyByType(t int32) int64 {
	for i := 0; i < len(cd21.bettypetotal); i++ {
		if cd21.bettypetotal[i].cd21type == t {
			return cd21.bettypetotal[i].money
		}
	}

	return 0
}

//发送用户下注信息
func (cd21 *CD21gameLogic) sendUserBetInfo(gameuser *gamelogic.GameUser) {
	ud := cd21.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetInfoID)

	pdata := new(netproto.CD21_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)

	for _, v := range ud.userbet {
		betinfo := new(netproto.CD21_BroadCastBet)
		betinfo.Type = proto.Int32(v.cd21type)
		betinfo.Money = proto.Int64(v.money)
		pdata.BetInfo = append(pdata.BetInfo, betinfo)
	}
	pmsg.MsgData = pdata
	cd21.SendMessageToUser(gameuser, pmsg)
}

//添加一个玩家数据
func (cd21 *CD21gameLogic) addUserData() {
	cd21.userdata = make([]*cd21userdata, 0)
	for _, ur := range cd21.AllTableUser() {
		userdata := newCD21userdata(ur)
		cd21.userdata = append(cd21.userdata, userdata)
	}
}

//删除一个玩家数据
func (cd21 *CD21gameLogic) delUserData(user *gamelogic.GameUser) {
	for i := 0; i < len(cd21.userdata); i++ {
		if cd21.userdata[i].user.UserID == user.UserID {
			cd21.userdata = append(cd21.userdata[:i], cd21.userdata[i+1:]...)
			break
		}
	}
}

//处理用户下注
func (cd21 *CD21gameLogic) userBet(user *gamelogic.GameUser, data []byte) {
	ud := cd21.getUserDataByID(user.UserID)
	ust := &netproto.CD21_BetArr{}
	err := proto.Unmarshal(data, ust)
	if err != nil {
		return
	}

	money := int32(0)
	for _, b := range ust.BetArr {
		*b.BetType = int32(math.Abs(float64(*b.BetType)))
		*b.Count = int32(math.Abs(float64(*b.Count)))

		//玩家身上金額小於下注限制
		if user.Score < cd21.betlimitmoney*100 {
			for i, u := range cd21.userdata {
				if u.user.UserID == ud.user.UserID {
					cd21.userdata = append(cd21.userdata[:i], cd21.userdata[i+1:]...)
				}
			}
			cd21.SendMessageToUser(user, cd21.getBetRetMsg(ud, -6, *b.BetType, ust.BetArr))
			return
		}

		//大小13不能同時下注
		if (*b.BetType == int32(BIGGER13) && ud.getUserBetMoneyByType(SMALLER13) > 0) ||
			(*b.BetType == int32(SMALLER13) && ud.getUserBetMoneyByType(BIGGER13) > 0) {
			cd21.SendMessageToUser(user, cd21.getBetRetMsg(ud, -5, *b.BetType, ust.BetArr))
			return
		}

		tmpmon := cd21.betvalue.getMoneyByType(*b.BetType) * (*b.Count)

		// 系統下注上限
		if (int64(tmpmon) + ud.getUserBetMoneyByType(*b.BetType)) > cd21.betvalue.getMaxMoneyByType(*b.BetType) {
			cd21.SendMessageToUser(user, cd21.getBetRetMsg(ud, -4, *b.BetType, ust.BetArr))
			return
		}

		money += tmpmon
	}

	ret := int32(0)
	if int64(money) > ud.money {
		ret = -1
	} else {
		for _, b := range ust.BetArr {
			betmon := cd21.betvalue.getMoneyByType(*b.BetType) * (*b.Count)
			if *b.BetType == int32(DIZHU) {
				ud.userDiZu += int64(betmon)
			} else if *b.BetType >= DUIZI && *b.BetType <= SMALLER13 {
				ud.betmoney += int64(betmon)
			}
			ud.setUserBet(*b.BetType, betmon)
			cd21.updateBet(*b.BetType, int64(betmon))
			if user.IsRobot() {
				cd21.updateRobotBet(*b.BetType, int64(betmon))
			}
		}
		cd21.GameLog.AddGameStep("用户%d下注%v", user.UserID, ust.BetArr)
		cd21.BoardcastMessage(cd21.getBetTotalMsg(user.UserID))
	}

	cd21.SendMessageToUser(user, cd21.getBetRetMsg(ud, ret, 0, ust.BetArr))
}

//获取下注的结果消息
func (cd21 *CD21gameLogic) getBetRetMsg(ud *cd21userdata, ret int32, t int32, betinfo []*netproto.CD21_Bet) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetRetID)

	pdata := new(netproto.CD21_BetRet)
	pdata.Money = proto.Int64(ud.money)
	pdata.Ret = proto.Int32(ret)
	pdata.Type = proto.Int32(t)
	pdata.BetArr = betinfo
	pmsg.MsgData = pdata
	return pmsg
}

//获取下注的结果消息
func (cd21 *CD21gameLogic) updateBet(cd21type int32, money int64) {
	cd21.bettotal += int64(money)
	//mlog.Debug("BetType=%v,Money=%v,bjl.bettotal=%v", bjlType, Money, bjl.bettotal)
	for _, v := range cd21.bettypetotal {
		if v.cd21type == cd21type {
			v.money += int64(money)
			if v.money < 0 {
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.cd21type = cd21type
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}
	cd21.bettypetotal = append(cd21.bettypetotal, b)
}

//获取机器人下注的结果消息
func (cd21 *CD21gameLogic) updateRobotBet(cd21type int32, money int64) {
	cd21.bettotalrobot += int64(money)
	//mlog.Debug("BetType=%v,Money=%v,bjl.bettotal=%v", bjlType, Money, bjl.bettotal)
	for _, v := range cd21.bettyperobot {
		if v.cd21type == cd21type {
			v.money += int64(money)
			if v.money < 0 {
				v.money = 0
			}
			return
		}
	}

	b := new(bettype)
	b.cd21type = cd21type
	b.money = int64(money)
	if b.money < 0 {
		b.money = 0
	}
	cd21.bettyperobot = append(cd21.bettyperobot, b)
}

//計算牌點數
func (cd21 *CD21gameLogic) getCardPoint(card *poker.CardGroup) int32 {
	number := int32(0)
	count := int32(0)
	for _, cardPoint := range card.Cards {
		if cardPoint.CardNumber == 14 {
			count++
		} else if cardPoint.CardNumber > 10 {
			number += 10
		} else {
			number += cardPoint.CardNumber
		}
	}

	if count != 0 {
		for i := int32(0); i < count; i++ {
			if (number + 11) <= 21-(count-i-1) {
				number += 11
			} else {
				number += 1
			}
		}
	}

	return number
}

//莊家要牌
func (cd21 *CD21gameLogic) fanZhuangPoker() {
	if cd21.gamestate == End {
		mlog.Warn("已結算過")
		return
	}

	cd21.setState(ZhuangFan)
	//先翻第二張牌
	msg := cd21.getFanPaiMsg(cd21.zhungCard)
	cd21.chesckZhuangBlackJack()
	cd21.BoardcastMessage(msg)
	cd21.TimerFuncManager.NewTimer(int32(1), cd21.dealZhuangPoker)
}

//莊家要牌
func (cd21 *CD21gameLogic) dealZhuangPoker() {
	for _, user := range cd21.userdata {
		//莊家爆牌、主遊戲流程結算
		if !user.isBust {
			if cd21.zhuangbucontroll == 1 && cd21.getCardPoint(cd21.zhungCard) >= 10 && cd21.getCardPoint(cd21.zhungCard) < 17 && cd21.syspoolData != nil {
				mlog.System("水池值:%d", cd21.syspoolData.PoolAmount)
				//根據水池補牌
				mlog.System("試圖要進水池補牌")
				cd21.syspoolbupai()
			} else {
				cd21.GameLog.AddGameStep("不控制莊家發牌")
				//自然發牌
				cd21.sendCardToBanker()
			}
			break
		}
	}

	cd21.opttimer = cd21.TimerFuncManager.NewTimer(cd21.zhuangFanPaiTime, cd21.finalSettlement)
}

//发送操作消息
func (cd21 *CD21gameLogic) sendOptMsg(bRecover bool) {
	ud := cd21.getUserBySeatNo(cd21.currSeat)
	if ud == nil {
		mlog.Error("玩家为空")
		return
	}
	opttime := cd21.optTime
	//如果不是场景恢复
	if !bRecover {
		msgall := cd21.getCurrOptMsg()
		cd21.BoardcastMessage(msgall)
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		//延迟一秒生效，防止网络卡顿造成弃牌的情况
		cd21.opttimer = cd21.TimerFuncManager.NewNamedTimer(cd21.optTime+1, cd21.userTimeOutStand, "")

	} else {
		if cd21.opttimer != nil {
			opttime -= cd21.opttimer.ElapsedSeconds()
		}
	}

	msg := ud.GetOptMsg(opttime, ud.canOpt())
	cd21.SendMessageToUser(ud.user, msg)
}

//发送分牌操作消息
func (cd21 *CD21gameLogic) sendSplitOptMsg(bRecover bool) {
	ud := cd21.getUserBySeatNo(cd21.currSeat)
	if ud == nil {
		mlog.Error("玩家为空")
		return
	}
	opttime := cd21.optTime
	//如果不是场景恢复
	if !bRecover {
		msgall := cd21.getSplitCurrOptMsg()
		cd21.BoardcastMessage(msgall)
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		//延迟一秒生效，防止网络卡顿造成弃牌的情况
		cd21.opttimer = cd21.TimerFuncManager.NewNamedTimer(cd21.optTime+1, cd21.userSplitTimeOutStand, "")

	} else {
		if cd21.opttimer != nil {
			opttime -= cd21.opttimer.ElapsedSeconds()
		}
	}

	msg := ud.GetSplitOptMsg(opttime, ud.cansplitOpt())
	cd21.SendMessageToUser(ud.user, msg)
}

//当前操作消息
func (cd21 *CD21gameLogic) getCurrOptMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21CurrOptID)

	pdata := new(netproto.CD21_WaitOptSeat)
	pdata.SitNo = proto.Int32(cd21.currSeat)
	pdata.OptTime = proto.Int32(cd21.optTime)
	pmsg.MsgData = pdata

	return pmsg
}

//当前操作消息
func (cd21 *CD21gameLogic) getSplitCurrOptMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21CurrOptID)

	pdata := new(netproto.CD21_WaitOptSeat)
	splitSitNo := int32(0)
	if cd21.currSeat == 1 {
		splitSitNo = 4
	} else if cd21.currSeat == 2 {
		splitSitNo = 5
	} else {
		splitSitNo = 3
	}
	pdata.SitNo = proto.Int32(splitSitNo)
	pdata.OptTime = proto.Int32(cd21.optTime)
	pmsg.MsgData = pdata

	return pmsg
}

//玩家停牌
func (cd21 *CD21gameLogic) userStand(urd *cd21userdata) {
	urd.SetStand()
	//移除操作超时Timer
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	cd21.GameLog.AddGameStep("用户%d主動停牌", urd.user.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StandBroadCastID)
	pdata := new(netproto.CD21_SitNo)
	pdata.SitNo = proto.Int32(urd.user.SeatNo)
	pmsg.MsgData = pdata

	cd21.BoardcastMessage(pmsg)
	cd21.waitOpt()
}

//玩家停牌
func (cd21 *CD21gameLogic) userSplitStand(urd *cd21userdata) {
	splitSitNo := int32(0)
	urd.SetSplitStand()
	//移除操作超时Timer
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	cd21.GameLog.AddGameStep("用户%d主動停牌", urd.user.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StandBroadCastID)
	pdata := new(netproto.CD21_SitNo)

	if urd.user.SeatNo == 1 {
		splitSitNo = 4
	} else if urd.user.SeatNo == 2 {
		splitSitNo = 5
	} else {
		splitSitNo = 3
	}
	pdata.SitNo = proto.Int32(splitSitNo)

	pmsg.MsgData = pdata

	cd21.BoardcastMessage(pmsg)

	if cd21.currSeat == urd.user.SeatNo { //當前操作不是自已，且用戶只剩一個人
		cd21.waitOpt()
	}
}

//玩家逾時停牌
func (cd21 *CD21gameLogic) userTimeOutStand() {
	urd := cd21.getUserBySeatNo(cd21.currSeat)
	if urd == nil {
		mlog.Error("用户逾時停牌")
		return
	}

	urd.SetStand()

	cd21.GameLog.AddGameStep("用户%d超時停牌", urd.user.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StandBroadCastID)
	pdata := new(netproto.CD21_SitNo)
	pdata.SitNo = proto.Int32(cd21.currSeat)
	pmsg.MsgData = pdata

	cd21.BoardcastMessage(pmsg)
	cd21.waitOpt()
}

//玩家分牌逾時停牌
func (cd21 *CD21gameLogic) userSplitTimeOutStand() {
	urd := cd21.getUserBySeatNo(cd21.currSeat)
	if urd == nil {
		mlog.Error("用户分牌逾時停牌")
		return
	}

	urd.SetSplitStand()

	cd21.GameLog.AddGameStep("用户%d分牌超時停牌", urd.user.UserID)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StandBroadCastID)
	pdata := new(netproto.CD21_SitNo)
	splitSitNo := int32(0)
	if cd21.currSeat == 1 {
		splitSitNo = 4
	} else if cd21.currSeat == 2 {
		splitSitNo = 5
	} else {
		splitSitNo = 3
	}
	pdata.SitNo = proto.Int32(splitSitNo)
	pmsg.MsgData = pdata

	cd21.BoardcastMessage(pmsg)
	cd21.waitOpt()
}

//玩家保險(21點贏2倍，非21扣保險)
func (cd21 *CD21gameLogic) insurance(urd *cd21userdata) {
	urd.isInsurance = true
	urd.betmoney += urd.diZuBase / 2
	msg := cd21.getSpecialBetMsg(urd, BAOXIAN, urd.diZuBase/2)
	cd21.GameLog.AddGameStep("用户%d保險%v", urd.user.UserID, urd.diZuBase/2)
	cd21.BoardcastMessage(msg)
	if cd21.tuobaoflag == len(cd21.userdata) {
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		cd21.waitOpt()
	}
}

//玩家投降及保險過程(1.10移除投降)
func (cd21 *CD21gameLogic) isTuoBao() {
	//傳送倒數計時
	cd21.setState(BaoXiao)
	cd21.BoardcastMessage(cd21.getTimesUpStateMsg(cd21.tuobaotime))

	for _, user := range cd21.AllTableUser() {
		msg := cd21.getTuoBaoMsg(false, false)
		ud := cd21.getUserDataByID(user.UserID)
		if ud != nil {
			//21點玩家不保險
			if ud.isBlackJack {
				cd21.tuobaoflag += 1
			}
			if (ud.user.Score-ud.betmoney-ud.userDiZu) > (ud.diZuBase/2) && !ud.isBlackJack {
				msg = cd21.getTuoBaoMsg(false, true)
			}
		}
		cd21.SendMessgeFunc(user.UserID, msg)
	}

	cd21.opttimer = cd21.TimerFuncManager.NewTimer(cd21.tuobaotime+1, cd21.waitOpt)
}

//获取用户的操作消息
func (cd21 *CD21gameLogic) getTuoBaoMsg(surrender bool, insurence bool) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21ZhuangTuoBaoID)

	pdata := new(netproto.CD21_UserCheck)
	pdata.Surrender = proto.Bool(surrender)
	pdata.Insurence = proto.Bool(insurence)
	if insurence {
		pdata.Cancel = proto.Bool(true)
	} else {
		pdata.Cancel = proto.Bool(false)
	}
	pmsg.MsgData = pdata
	return pmsg
}

//用户要牌
func (cd21 *CD21gameLogic) userHit(urd *cd21userdata) {
	cards := cd21.poker.Deal(urd.user.UserID, 1)
	for _, dealcard := range cards.Cards {
		dCard := new(poker.Card)
		dCard.CardNumber = dealcard.CardNumber
		dCard.Color = dealcard.Color
		urd.userCards.Cards = append(urd.userCards.Cards, dCard)
	}
	cd21.GameLog.AddGameStep("玩家%d, 玩家主動要牌 %v", urd.user.UserID, urd.userCards.String())
	msg := cd21.getHitpaiMsg(urd, false, false)
	cd21.BoardcastMessage(msg)
	if cd21.getCardPoint(urd.userCards) >= 21 {
		urd.SetStand()
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		if cd21.getCardPoint(urd.userCards) == 21 {
			cd21.GameLog.AddGameStep("玩家%d21點, 系統自動停牌 %d", urd.user.UserID, cd21.getCardPoint(urd.userCards))
		} else {
			urd.isBust = true
			cd21.GameLog.AddGameStep("玩家%d爆牌, 牌點總和 %d", urd.user.UserID, cd21.getCardPoint(urd.userCards))
		}
		cd21.TimerFuncManager.NewTimer(int32(1), cd21.waitOpt)
		return
	}
	cd21.sendOptMsg(false)
}

//用户分牌後要牌
func (cd21 *CD21gameLogic) userSplitHit(urd *cd21userdata) {
	cards := cd21.poker.Deal(urd.user.UserID, 1)
	cd21.setUserSplitCards(urd, cards.Cards[0])
	cd21.GameLog.AddGameStep("玩家%d, 玩家分牌後主動要牌 %v", urd.user.UserID, urd.splitCards.String())
	msg := cd21.getHitpaiMsg(urd, false, true)
	cd21.BoardcastMessage(msg)
	if cd21.getCardPoint(urd.splitCards) >= 21 {
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		if cd21.getCardPoint(urd.splitCards) == 21 {
			urd.isStand1 = true
			urd.isBust1 = true
			cd21.GameLog.AddGameStep("玩家%d分牌後21點, 系統自動停牌 %d", urd.user.UserID, cd21.getCardPoint(urd.userCards))
		} else {
			urd.isStand1 = true
			urd.isBust1 = true
			cd21.GameLog.AddGameStep("玩家%d分牌後爆牌, 牌點總和 %d", urd.user.UserID, cd21.getCardPoint(urd.splitCards))
		}
		cd21.TimerFuncManager.NewTimer(int32(1), cd21.waitOpt)
		return
	}
	cd21.sendSplitOptMsg(false)
}

//用户分牌
func (cd21 *CD21gameLogic) splitCard(urd *cd21userdata) {
	//移除操作超时Timer
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	cd21.GameLog.AddGameStep("用户%d分牌", urd.user.UserID)
	//廣播特殊下注
	urd.isSplit = true
	urd.userDiZu += urd.diZuBase
	spmsg := cd21.getSpecialBetMsg(urd, SPLIT, urd.diZuBase)
	cd21.BoardcastMessage(spmsg)
	//將原本的牌分成兩張
	cd21.setUserSplitCards(urd, urd.userCards.Cards[1])
	urd.userCards.Cards = append(urd.userCards.Cards[:1], urd.userCards.Cards[2:]...)
	//再各發一張牌
	card1 := cd21.poker.Deal(urd.user.UserID, 1)
	dCard := new(poker.Card)
	dCard.CardNumber = card1.Cards[0].CardNumber
	dCard.Color = card1.Cards[0].Color
	urd.userCards.Cards = append(urd.userCards.Cards, dCard)
	card2 := cd21.poker.Deal(urd.user.UserID, 1)
	cd21.setUserSplitCards(urd, card2.Cards[0])
	//廣播牌形變化
	msg := cd21.getSplitpaiMsg(urd)
	cd21.BoardcastMessage(msg)

	if cd21.getCardPoint(urd.userCards) == 21 && cd21.getCardPoint(urd.splitCards) == 21 {
		cd21.TimerFuncManager.NewTimer(int32(2), cd21.fanZhuangPoker)
	} else if cd21.getCardPoint(urd.userCards) != 21 {
		cd21.TimerFuncManager.NewTimer(int32(2), func() {
			cd21.sendOptMsg(false)
		})
	} else {
		cd21.TimerFuncManager.NewTimer(int32(2), func() {
			urd.isStand = true
			cd21.sendSplitOptMsg(false)
		})
	}
}

//用户加倍(再給一張牌 底注X2)
func (cd21 *CD21gameLogic) doubleCard(urd *cd21userdata) {
	urd.SetStand()
	//移除操作超时Timer
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	urd.userDiZu += urd.diZuBase
	spmsg := cd21.getSpecialBetMsg(urd, DIZHU, urd.diZuBase)
	cd21.BoardcastMessage(spmsg)
	urd.isDouble = true
	//底注X2
	cards := cd21.poker.Deal(urd.user.UserID, 1)
	for _, dealcard := range cards.Cards {
		dCard := new(poker.Card)
		dCard.CardNumber = dealcard.CardNumber
		dCard.Color = dealcard.Color
		urd.userCards.Cards = append(urd.userCards.Cards, dCard)
	}
	msg := cd21.getHitpaiMsg(urd, true, false)
	cd21.BoardcastMessage(msg)
	cd21.GameLog.AddGameStep("用户%d主動加倍, 牌型%v", urd.user.UserID, urd.userCards.String())
	if cd21.getCardPoint(urd.userCards) > 21 {
		urd.isBust = true
		cd21.GameLog.AddGameStep("玩家%d爆牌, 牌點總和 %d", urd.user.UserID, cd21.getCardPoint(urd.userCards))
	}
	cd21.waitOpt()
}

//用户加倍(再給一張牌 底注X2)
func (cd21 *CD21gameLogic) splitdoubleCard(urd *cd21userdata) {
	urd.SetSplitStand()
	//移除操作超时Timer
	if cd21.opttimer != nil {
		cd21.opttimer.Stop()
		cd21.opttimer = nil
	}
	urd.userDiZu += urd.diZuBase
	spmsg := cd21.getSpecialBetMsg(urd, SPLIT, urd.diZuBase)
	cd21.BoardcastMessage(spmsg)
	urd.isDouble1 = true
	//底注X2
	cards := cd21.poker.Deal(urd.user.UserID, 1)
	cd21.setUserSplitCards(urd, cards.Cards[0])
	msg := cd21.getHitpaiMsg(urd, true, true)
	cd21.GameLog.AddGameStep("用户%d分牌後主動加倍, 牌型%v", urd.user.UserID, urd.splitCards.String())
	cd21.BoardcastMessage(msg)
	if cd21.getCardPoint(urd.splitCards) > 21 {
		urd.isBust1 = true
		cd21.GameLog.AddGameStep("玩家%d分牌爆牌, 牌點總和 %d", urd.user.UserID, cd21.getCardPoint(urd.splitCards))
	}
	cd21.waitOpt()
}

//获取要牌消息
func (cd21 *CD21gameLogic) getHitpaiMsg(urd *cd21userdata, isDouble bool, isSplit bool) *msg.Message {
	splitSitNo := int32(0)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BroadCastYaopaiID)
	fdata := new(netproto.CD21_Yaopai)

	if isSplit {
		if urd.user.SeatNo == 1 {
			splitSitNo = 4
		} else if urd.user.SeatNo == 2 {
			splitSitNo = 5
		} else {
			splitSitNo = 3
		}
		fdata.SitNo = proto.Int32(splitSitNo)
		for _, cardGroup := range urd.splitCards.Cards {
			npk := new(netproto.PokerCard)
			npk.Color = proto.Int32(cardGroup.Color)
			npk.CardNumber = proto.Int32(cardGroup.CardNumber)
			fdata.SitPoker = append(fdata.SitPoker, npk)
		}
		fdata.CardTotal = proto.Int32(cd21.getCardPoint(urd.splitCards))
	} else {
		fdata.SitNo = proto.Int32(urd.user.SeatNo)
		for _, cardGroup := range urd.userCards.Cards {
			npk := new(netproto.PokerCard)
			npk.Color = proto.Int32(cardGroup.Color)
			npk.CardNumber = proto.Int32(cardGroup.CardNumber)
			fdata.SitPoker = append(fdata.SitPoker, npk)
		}
		fdata.CardTotal = proto.Int32(cd21.getCardPoint(urd.userCards))
	}

	fdata.Double = proto.Bool(isDouble)

	if cd21.isCut {
		fdata.LastTurn = proto.Bool(true)
	} else {
		fdata.LastTurn = proto.Bool(false)
	}
	pmsg.MsgData = fdata
	return pmsg
}

//获取要牌消息
func (cd21 *CD21gameLogic) getSplitpaiMsg(urd *cd21userdata) *msg.Message {
	splitSitNo := int32(0)
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21SplitCardBroacast)
	fdata := new(netproto.CD21_SplitCard)
	fdata.SitNo = proto.Int32(urd.user.SeatNo)
	for _, cardGroup := range urd.userCards.Cards {
		npk := new(netproto.PokerCard)
		npk.Color = proto.Int32(cardGroup.Color)
		npk.CardNumber = proto.Int32(cardGroup.CardNumber)
		fdata.SitPoker = append(fdata.SitPoker, npk)
	}
	fdata.CardTotal = proto.Int32(cd21.getCardPoint(urd.userCards))
	if urd.user.SeatNo == 1 {
		splitSitNo = 4
	} else if urd.user.SeatNo == 2 {
		splitSitNo = 5
	} else {
		splitSitNo = 3
	}
	for _, cardGroup := range urd.splitCards.Cards {
		npk := new(netproto.PokerCard)
		npk.Color = proto.Int32(cardGroup.Color)
		npk.CardNumber = proto.Int32(cardGroup.CardNumber)
		fdata.SplitSitPoker = append(fdata.SplitSitPoker, npk)
	}
	fdata.SplitCardTotal = proto.Int32(cd21.getCardPoint(urd.splitCards))
	fdata.SplitSitNo = proto.Int32(splitSitNo)
	pmsg.MsgData = fdata
	return pmsg
}

//获取赢家的消息
func (cd21 *CD21gameLogic) getUserWinMsg(winType int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21GameResultWinID)
	fdata := new(netproto.CD21_GameResultUserMoney)
	for _, user := range cd21.userdata {
		pdata := new(netproto.CD21_GameSettlement)
		pdata.SitNo = proto.Int32(user.user.SeatNo)
		if winType == BetArea {
			pdata.TotalWin = proto.Int64(user.win)
			for _, winBlock := range user.userwinblock {
				bdata := new(netproto.CD21_WinBlock)
				bdata.NoWin = proto.Bool(false)
				if winBlock == DUIZI {
					bdata.SelfWin = proto.Int64(user.winDUIZI)
					bdata.Block = proto.Int32(DUIZI)
				} else if winBlock == LUCKY20 {
					bdata.SelfWin = proto.Int64(user.winLUCKY20)
					bdata.Block = proto.Int32(LUCKY20)
				} else if winBlock == BIGGER13 {
					bdata.SelfWin = proto.Int64(user.winBIGGER13)
					bdata.Block = proto.Int32(BIGGER13)
				} else if winBlock == SMALLER13 {
					bdata.SelfWin = proto.Int64(user.winSMALLER13)
					bdata.Block = proto.Int32(SMALLER13)
				}
				pdata.Block = append(pdata.Block, bdata)
			}
		} else {
			if winType == MainGame {
				pdata.TotalWin = proto.Int64(user.win1 + user.win4 + user.win5)
				if user.isBlackJack {
					bjdata := new(netproto.CD21_WinBlock)
					bjdata.NoWin = proto.Bool(false)
					pdata.TotalWin = proto.Int64(user.win1)
					bjdata.SelfWin = proto.Int64(user.winBLACKJACK)
					bjdata.Block = proto.Int32(DIZHU)
					pdata.Block = append(pdata.Block, bjdata)
				} else {
					turn := 1
					if user.isSplit {
						turn = 2
					}
					for i := 0; i < turn; i++ {
						bdata := new(netproto.CD21_WinBlock)
						if i == 0 {
							bdata.NoWin = proto.Bool(user.userNoWin)
							bdata.SelfWin = proto.Int64(user.winMAINGAME)
							bdata.Block = proto.Int32(DIZHU)
						} else {
							bdata.NoWin = proto.Bool(user.userSplitNoWin)
							bdata.SelfWin = proto.Int64(user.winSPILTGAME)
							bdata.Block = proto.Int32(SPLIT)
						}
						pdata.Block = append(pdata.Block, bdata)
					}
				}

			} else {
				bdata := new(netproto.CD21_WinBlock)
				bdata.NoWin = proto.Bool(false)
				if winType == TuoBao {
					pdata.TotalWin = proto.Int64(user.win2)
					bdata.SelfWin = proto.Int64(user.win2)
					bdata.Block = proto.Int32(BAOXIAN)
				} else if winType == ZhungBust {
					//補上所有邊注遊戲金額
					pdata.TotalWin = proto.Int64(user.win + user.win3)
					bdata.SelfWin = proto.Int64(user.winZHUANGBAO)
					bdata.Block = proto.Int32(ZHUANGBAO)
				}
				pdata.Block = append(pdata.Block, bdata)
			}
		}
		fdata.ResultUserMoney = append(fdata.ResultUserMoney, pdata)
	}
	pmsg.MsgData = fdata
	return pmsg
}

//莊家21點
func (cd21 *CD21gameLogic) chesckZhuangBlackJack() {
	for _, ud := range cd21.userdata {
		if ud.isInsurance {
			if cd21.getCardPoint(cd21.zhungCard) == 21 {
				//賺保險金
				ud.win2 = ud.diZuBase + (ud.diZuBase / 2)
			}
			cd21.GameLog.AddGameStep("中間結算, 用户保險赢的钱 : %d", ud.win2-(ud.diZuBase/2))
			cd21.BoardcastMessage(cd21.getUserWinMsg(TuoBao))
		}
	}
}

//莊家要牌
func (cd21 *CD21gameLogic) sendCardToBanker() {
	//發莊家剩餘牌
	for {
		var number int32 = 0
		number = cd21.getCardPoint(cd21.zhungCard)
		if number < 17 {
			zhungCard := cd21.poker.Deal(0, 1)
			for _, dealcard := range zhungCard.Cards {
				dCard := new(poker.Card)
				dCard.CardNumber = dealcard.CardNumber
				dCard.Color = dealcard.Color
				cd21.zhungCard.Cards = append(cd21.zhungCard.Cards, dCard)
			}

		} else {
			break
		}
	}
	cd21.GameLog.AddGameStep("莊家要牌結果 %v", cd21.zhungCard.String())
	if cd21.getCardPoint(cd21.zhungCard) > 21 {
		cd21.zhuangBust = true
		cd21.GameLog.AddGameStep("莊家爆牌, 牌點總和 %d", cd21.getCardPoint(cd21.zhungCard))

	}
	time := math.Ceil(float64(len(cd21.zhungCard.Cards)) * 0.7)
	cd21.zhuangFanPaiTime = int32(time) + 1
	cd21.BoardcastMessage(cd21.getFaZhuangPaiMsg(cd21.zhungCard))
}

func (cd21 *CD21gameLogic) sendUsersetResult(ud *cd21userdata) {
	var income, tax int64
	for _, bettype := range ud.userbet {
		if bettype.cd21type == ZHUANGBAO {
			cd21.sumBetCount += 1
			if cd21.zhuangBust {
				if len(cd21.zhungCard.Cards) == 5 {
					ud.winZHUANGBAO += (cd21.betodd.getOddByType(rateZHUANGBAO5) - 1) * bettype.money
				} else if len(cd21.zhungCard.Cards) == 6 {
					ud.winZHUANGBAO += (cd21.betodd.getOddByType(rateZHUANGBAO6) - 1) * bettype.money
				} else if len(cd21.zhungCard.Cards) == 7 {
					ud.winZHUANGBAO += (cd21.betodd.getOddByType(rateZHUANGBAO7) - 1) * bettype.money
				} else if len(cd21.zhungCard.Cards) == 8 {
					ud.winZHUANGBAO += (cd21.betodd.getOddByType(rateZHUANGBAO8) - 1) * bettype.money
				} else if len(cd21.zhungCard.Cards) >= 9 {
					ud.winZHUANGBAO += (cd21.betodd.getOddByType(rateZHUANGBAO9) - 1) * bettype.money
				}
			}
		}
		income, tax = cd21.GetUserWinAndTax(ud.winZHUANGBAO)
		if income > 0 {
			income += bettype.money
		}
		ud.win3 = income
		ud.tax += tax
	}

	cd21.GameLog.AddGameStep("中間結算, 莊家爆牌區赢的钱: %d", ud.win3)

}

func (cd21 *CD21gameLogic) compareWithBanker(ud *cd21userdata) {
	var income, tax int64
	var winner string
	win := int64(0) //下注區塊贏錢
	odd := int64(1) //是否double倍率
	if ud.isDouble {
		odd = 2
	}
	ud.betmoney += ud.userDiZu //底注
	bpoint := cd21.getCardPoint(cd21.zhungCard)
	upoint := cd21.getCardPoint(ud.userCards)
	cd21.GameLog.AddGameStep("玩家 %d, 莊家牌點: %d,玩家牌點: %d", ud.user.UserID, bpoint, upoint)
	if bpoint > 21 {
		bpoint = 0
	}
	if upoint > 21 { //1.10移除投降
		upoint = 0
	}
	if bpoint > upoint {
		//莊贏
		winner = "莊贏"
	} else if bpoint < upoint {
		//閒贏
		winner = "玩家贏"
		win += ud.diZuBase * odd
		income, tax = cd21.GetUserWinAndTax(win)
		ud.winMAINGAME = income
		income += ud.diZuBase * odd
		ud.win4 += income
		ud.tax += tax
	} else {
		//莊閒同點，push (1.10移除投降)
		winner = "平手"
		win += ud.diZuBase * odd
		ud.userNoWin = true
		if upoint != 0 {
			ud.win4 += win
			ud.winMAINGAME = win
		}
	}
	mlog.Debug("贏家:%s", winner)
	cd21.GameLog.AddGameStep("主遊戲結算, 用户主遊戲贏錢: %d", ud.win4-ud.diZuBase)
}

func (cd21 *CD21gameLogic) splitWithBanker(ud *cd21userdata) {
	var income, tax int64
	var winner string
	win := int64(0) //下注區塊贏錢
	odd := int64(1) //是否double倍率
	if ud.isDouble1 {
		odd = 2
	}
	if ud.splitCards != nil {
		sbpoint := cd21.getCardPoint(cd21.zhungCard)
		supoint := cd21.getCardPoint(ud.splitCards)
		cd21.GameLog.AddGameStep("玩家 %d, 莊家牌點: %d,玩家分牌牌點: %d", ud.user.UserID, sbpoint, supoint)

		if sbpoint > 21 {
			sbpoint = 0
		}
		if supoint > 21 {
			supoint = 0
		}
		if sbpoint > supoint {
			//莊贏
			winner = "分牌莊贏"
		} else if sbpoint < supoint {
			//閒贏
			winner = "分牌玩家贏"
			win += ud.diZuBase * odd
			income, tax = cd21.GetUserWinAndTax(win)
			ud.winSPILTGAME = income
			income += ud.diZuBase * odd
			ud.win5 += income
			ud.tax += tax
		} else {
			//莊閒同點，push
			winner = "分牌平手"
			win += ud.diZuBase * odd
			ud.userSplitNoWin = true
			if supoint != 0 {
				ud.win5 += win
				ud.winSPILTGAME = win
			}
		}
		mlog.Debug("贏家:%s", winner)
		cd21.GameLog.AddGameStep("主遊戲結算, 用户分牌輸贏: %d", ud.win5)
	}
}

func (cd21 *CD21gameLogic) passhand(ud *cd21userdata) {
	if cd21.tuobaoflag == len(cd21.userdata) {
		//移除操作超时Timer
		if cd21.opttimer != nil {
			cd21.opttimer.Stop()
			cd21.opttimer = nil
		}
		cd21.waitOpt()
	}
}

//用户断线重连
func (cd21 *CD21gameLogic) UserRecover(gameuser *gamelogic.GameUser) {
	defer util.PanicErrStack()
	cd21.DynamicGameLogicBase.UserRecover(gameuser)
	//发送游戏场景
	cd21.sendScene(gameuser)
}

//全部分數結算
func (cd21 *CD21gameLogic) finalScore(ud *cd21userdata) {
	//ud.win含底注
	//ud.betmoney:底注
	cd21.GameLog.AddGameStep("用户%d, 底注 %d", ud.user.UserID, ud.betmoney)
	win := ud.win + ud.win1 + ud.win2 + ud.win3 + ud.win4 + ud.win5 + ud.win6 - ud.betmoney
	tax := ud.tax
	if win > 0 {
		ud.user.FlowAmount += win
	} else {
		ud.user.FlowAmount -= win
	}
	cd21.GameLog.AddGameStep("玩家 %d, 用户下注區赢的钱 %d, 用户21點赢的钱 %d, 用户保險赢的钱 %d, 用户莊家爆牌區赢的钱 %d, 用户主遊戲贏錢 %d, 用户主遊戲分牌贏錢 %d, 投降返還一半 %d", ud.user.UserID, ud.win, ud.win1, ud.win2, ud.win3, ud.win4, ud.win5, ud.win6)
	cd21.GameLog.AddGameStep("玩家 %d, 當前金幣 userMoney %d, 進房分數 %d, win %d (流水計算), tax %d", ud.user.UserID, ud.money, ud.user.Score, win, tax)
	if !ud.user.IsRobot() {
		cd21.sumBetMoney += ud.betmoney
		cd21.poolAmount -= win
	}
	cd21.UpdateUserScoreFunc(ud.user.UserID, win, tax)
	cd21.uscoresRet[ud.user.UserID] = win
	cd21.GameLog.Context().UserUpdateScore(ud.user.UserID, win)
}

//全部分數結算
func (cd21 *CD21gameLogic) setFirst() {
	cd21.first = false
}

//获取总下注的消息
func (cd21 *CD21gameLogic) getRecoverBetTotalMsg(user *cd21userdata, betstatus int32) *msg.Message {
	if len(cd21.userdata) == 0 {
		return nil
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserRecoverBetID)

	pdata := new(netproto.CD21_RecoverScene)

	pdata.UserID = proto.Int32(user.user.UserID)
	for _, v := range user.userbet {
		if betstatus == BetRecover {
			betinfo := new(netproto.CD21_RecoverBetInfo)
			betinfo.Type = proto.Int32(v.cd21type)
			betinfo.Money = proto.Int64(v.money)
			pdata.RecoverBetArr = append(pdata.RecoverBetArr, betinfo)
		} else {
			if v.cd21type >= DIZHU || v.cd21type == ZHUANGBAO {
				betinfo := new(netproto.CD21_RecoverBetInfo)
				betinfo.Type = proto.Int32(v.cd21type)
				betinfo.Money = proto.Int64(v.money)
				pdata.RecoverBetArr = append(pdata.RecoverBetArr, betinfo)
			}
		}

	}

	pmsg.MsgData = pdata
	return pmsg
}

//获取总下注的消息
func (cd21 *CD21gameLogic) sendUserDiZuInfo(gameuser *gamelogic.GameUser) {
	ud := cd21.getUserDataByID(gameuser.UserID)
	if ud == nil {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21BetInfoID)

	pdata := new(netproto.CD21_UserBetInfo)
	pdata.UserMoney = proto.Int64(ud.money)

	for _, v := range ud.userbet {
		if v.cd21type >= DIZHU {
			betinfo := new(netproto.CD21_BroadCastBet)
			betinfo.Type = proto.Int32(v.cd21type)
			betinfo.Money = proto.Int64(v.money)
			pdata.BetInfo = append(pdata.BetInfo, betinfo)
		}
	}
	pmsg.MsgData = pdata
	cd21.SendMessageToUser(gameuser, pmsg)
}

//当前操作消息
func (cd21 *CD21gameLogic) getSpecialBetMsg(urd *cd21userdata, betType int32, money int64) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserSpecialBetID)

	pdata := new(netproto.CD21_BroadCastSpecialBet)
	pdata.UserID = proto.Int32(urd.user.UserID)
	pdata.Type = proto.Int32(betType)
	pdata.Money = proto.Int64(money)

	pmsg.MsgData = pdata
	return pmsg
}

//设置状态
func (cd21 *CD21gameLogic) getTimesUpStateMsg(time int32) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21StateID)

	pdata := new(netproto.CD21_Sate)
	pdata.State = proto.Int32(cd21.gamestate)
	pdata.Time = proto.Int32(time)
	pmsg.MsgData = pdata
	return pmsg
}

//最終結算
func (cd21 *CD21gameLogic) sendReadyTime() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21GameEndReadyTimeID)

	pdata := new(netproto.CD21_ReadyTime)
	pdata.ReadyTime = proto.Int32(cd21.readyTime)
	pmsg.MsgData = pdata
	cd21.BoardcastMessage(pmsg)
	mlog.Debug("GameEndFunc 发送倒计时 readyTime[%v]", cd21.readyTime)
	cd21.GameEndFunc()
}

//获取发牌消息
func (cd21 *CD21gameLogic) getPokerRecoverMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserPokerRecoverID)
	fdata := new(netproto.CD21_UserPokerRecover)

	//各玩家牌型根據座位號塞入資訊
	for _, u := range cd21.userdata {
		pdata := new(netproto.CD21_PokerRecover)
		pdata.SitNo = proto.Int32(u.user.SeatNo)
		for _, cardGroup := range u.userCards.Cards {
			npk := new(netproto.PokerCard)
			npk.Color = proto.Int32(cardGroup.Color)
			npk.CardNumber = proto.Int32(cardGroup.CardNumber)
			pdata.SitPoker = append(pdata.SitPoker, npk)
		}
		pdata.CardTotal = proto.Int32(cd21.getCardPoint(u.userCards))
		pdata.Double = proto.Bool(u.isDouble)
		fdata.PokerRecover = append(fdata.PokerRecover, pdata)
		if u.isSplit {
			splitSitNo := int32(0)
			spdata := new(netproto.CD21_PokerRecover)
			if u.user.SeatNo == 1 {
				splitSitNo = 4
			} else if u.user.SeatNo == 2 {
				splitSitNo = 5
			} else {
				splitSitNo = 3
			}
			spdata.SitNo = proto.Int32(splitSitNo)
			for _, scardGroup := range u.splitCards.Cards {
				snpk := new(netproto.PokerCard)
				snpk.Color = proto.Int32(scardGroup.Color)
				snpk.CardNumber = proto.Int32(scardGroup.CardNumber)
				spdata.SitPoker = append(spdata.SitPoker, snpk)
			}
			spdata.CardTotal = proto.Int32(cd21.getCardPoint(u.splitCards))
			spdata.Double = proto.Bool(u.isDouble1)
			fdata.PokerRecover = append(fdata.PokerRecover, spdata)
		}
	}

	//塞入莊家牌型
	zdata := new(netproto.CD21_PokerRecover)
	zdata.SitNo = proto.Int32(99)
	if cd21.gamestate < ZhuangFan {
		cano := int32(0)
		for i := 0; i < 2; i++ {
			if i == 0 {
				zpk := new(netproto.PokerCard)
				cano = cd21.zhungCard.Cards[0].CardNumber
				zpk.CardNumber = proto.Int32(cd21.zhungCard.Cards[0].CardNumber)
				zpk.Color = proto.Int32(cd21.zhungCard.Cards[0].Color)
				zdata.SitPoker = append(zdata.SitPoker, zpk)
			} else {
				zpk := new(netproto.PokerCard)
				zpk.CardNumber = proto.Int32(0)
				zpk.Color = proto.Int32(0)
				zdata.SitPoker = append(zdata.SitPoker, zpk)
			}
		}
		if cano == 14 {
			cano = 11
		} else if cano > 10 {
			cano = 10
		}
		zdata.CardTotal = proto.Int32(cano)
	} else {
		for _, zcardGroup := range cd21.zhungCard.Cards {
			zpk := new(netproto.PokerCard)
			zpk.CardNumber = proto.Int32(zcardGroup.CardNumber)
			zpk.Color = proto.Int32(zcardGroup.Color)
			zdata.SitPoker = append(zdata.SitPoker, zpk)
		}
		zdata.CardTotal = proto.Int32(cd21.getCardPoint(cd21.zhungCard))
	}
	zdata.Double = proto.Bool(false)
	fdata.PokerRecover = append(fdata.PokerRecover, zdata)

	pmsg.MsgData = fdata
	return pmsg
}

//获取階段贏錢的消息
func (cd21 *CD21gameLogic) getRecoverWinLoseMsg(user *cd21userdata) *msg.Message {
	defer util.PanicErrStack()
	if len(cd21.userdata) == 0 || user == nil || user.user == nil {
		return nil
	}
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserWinLoseRecoverID)

	pdata := new(netproto.CD21_WinLoseRecover)
	pdata.UserId = proto.Int32(user.user.UserID)
	winLose := user.betmoney + user.userDiZu
	if user.win2 > 0 {
		winLose = user.betmoney + user.userDiZu - user.win2
	}
	pdata.WinLose = proto.Int64(winLose)
	for _, winBlock := range user.userwinblock {
		money := int64(0)
		bdata := new(netproto.CD21_WinLoseDetail)
		bdata.SelfWin = proto.Int64(0)
		bdata.Block = proto.Int32(0)
		if winBlock == DUIZI {
			if user.winDUIZI > 0 {
				money = user.getUserBetMoneyByType(DUIZI)
				bdata.SelfWin = proto.Int64(user.winDUIZI + money)
				bdata.Block = proto.Int32(DUIZI)
			}
		} else if winBlock == LUCKY20 {
			if user.winLUCKY20 > 0 {
				money = user.getUserBetMoneyByType(LUCKY20)
				bdata.SelfWin = proto.Int64(user.winLUCKY20 + money)
				bdata.Block = proto.Int32(LUCKY20)
			}
		} else if winBlock == BIGGER13 {
			if user.winBIGGER13 > 0 {
				money = user.getUserBetMoneyByType(BIGGER13)
				bdata.SelfWin = proto.Int64(user.winBIGGER13 + money)
				bdata.Block = proto.Int32(BIGGER13)
			}
		} else if winBlock == SMALLER13 {
			if user.winSMALLER13 > 0 {
				money = user.getUserBetMoneyByType(SMALLER13)
				bdata.SelfWin = proto.Int64(user.winSMALLER13 + money)
				bdata.Block = proto.Int32(SMALLER13)
			}
		}
		pdata.WinLoseBlock = append(pdata.WinLoseBlock, bdata)
	}

	pmsg.MsgData = pdata
	return pmsg

}

//设置觀戰資訊
func (cd21 *CD21gameLogic) getWatchMsg() *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.CD21_GameMessageClassID_CD21UserWatchDataID)

	pdata := new(netproto.CD21_WatchData)

	for _, user := range cd21.AllTableUser() {
		wdata := new(netproto.CD21_WatchDetail)
		wdata.UserID = proto.Int32(user.UserID)
		ud := cd21.getUserDataByID(user.UserID)
		if ud != nil {
			wdata.Watch = proto.Bool(false)
		} else {
			wdata.Watch = proto.Bool(true)
		}
		pdata.WatchDetail = append(pdata.WatchDetail, wdata)
	}
	pmsg.MsgData = pdata
	return pmsg
}
