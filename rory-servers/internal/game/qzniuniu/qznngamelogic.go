package qznn

import (
	"fmt"
	"math"
	"net/http"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	gameStep_Begin         = iota //游戏开始
	gameStep_ZhuangConfirm        //庄家及倍数确定
	gameStep_Xianbet              //闲家倍数确定
	gameStep_CardGroup            //组牌完成
	gameStep_GameRet              //亮牌
)

//动作名称
const (
	actionName_QiangZhuang     = "qiangzhuang"     //抢庄动作名称
	actionName_XianjiaBetTimes = "xianjiabettimes" //闲家倍数动作名称
	actionName_CardGroup       = "cardgroup"       //组牌动作名称
)

//动画名称
const (
	animationName_GameStart   = "gamestart"   //游戏开始
	animationName_FirstCard   = "firstcard"   //第一手牌
	animationName_SecondCard  = "secondcard"  //第二手牌
	animationName_GameRet     = "gameret"     //结算
	animationName_QiangZhuang = "qiangzhuang" //抢庄

	animationTimeOut_GameStart   = 3
	animationTimeOut_FirstCard   = 4
	animationTimeout_SecondCard  = 3
	animationTimeout_GameRet     = 4
	animationTimeout_QiangZhuang = 4
)

//抢庄牛牛
type QZNNGameLogic struct {
	*gamelogic.DynamicGameLogicBase
	poker                     poker.Poker            //游戏扑克牌
	gamePoker                 *poker.GamePoker       //玩家牌型
	zhuangInfo                *userBet               //庄家信息
	gameStep                  int                    //游戏阶段
	qiangZhuangBets           *zhuangJiaBets         //抢庄倍数
	xianjiaBets               *xianjiaBets           //闲家倍数
	gameCardStack             *poker.GameCardStack   //用户组牌
	gameRet                   *netproto.QZNN_GameRet //游戏结果
	actionName_QiangZhuang    string
	actionName_XianBetTimes   string
	actionName_CardGroup      string
	currentWatchAnimationName string
	jackpotTotalPool          int64
	SuperToolService          *gamelogic.SuperToolBrGameServiceBase //超端服務
}

//创建抢庄牛牛游戏逻辑对象
func NewQZNNGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *QZNNGameLogic {
	glc := new(QZNNGameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)

	//解析配置
	parseNNConfig(cfgs)
	if glc.GetSuperToolService() != nil {
		glc.SuperToolService = glc.GetSuperToolService().(*gamelogic.SuperToolBrGameServiceBase)
	}

	glc.RegisterAfterGameStartFunc(func() {
		//通知用户抢庄倍数
		for _, ur := range glc.AllGameUser() {
			glc.sendGameSetting(ur)
			glc.notifyUserBetZhuangTimes(ur)
		}

		//发牌动画
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, glc.fapai)
	})

	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})
	return glc
}

//清除游戏数据
func (qznn *QZNNGameLogic) clearGameData() {
	//随机发牌
	if qznn.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		qznn.poker = poker.NewRandomPoker(nnconfig.hasPokerKing, false, 1)
		qznn.poker = poker.NewFixedPoker(qznn.AllTableUser(), qznn.poker.(*poker.RandomPoker), 5)
	} else { //发牌器发牌
		qznn.poker = poker.NewFilePoker(qznn.Configs, qznn.AllTableUser())
		qznn.poker = poker.NewFixedFilePoker(qznn.poker.(*poker.FilePoker))
	}

	qznn.gamePoker = poker.NewGamePoker()         //创建游戏牌型
	qznn.zhuangInfo = nil                         //庄家信息
	qznn.gameStep = gameStep_Begin                //设置游戏阶段
	qznn.qiangZhuangBets = newZhuangJiaBets()     //初始化抢庄信息
	qznn.xianjiaBets = newXianjiaBets()           //初始化闲家倍数信息
	qznn.gameCardStack = poker.NewGameCardStack() //初始化用户组牌信息
	qznn.gameRet = nil

	qznn.actionName_QiangZhuang = ""
	qznn.actionName_XianBetTimes = ""
	qznn.actionName_CardGroup = ""

	//当前监视动画名称
	qznn.currentWatchAnimationName = ""
}

//通知用户牛牛倍数
func (qznn *QZNNGameLogic) notifyUserBetTimes(ur *gamelogic.GameUser) {

	ubtmsg := new(msg.Message)
	ubtmsg.BClassID = int32(netproto.MessageBClassID_Game)
	ubtmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_UserBetInfoID)
	ubtmsg.MsgData = qznn.getUserBetTimesMsgData(ur)

	qznn.SendMessageToUser(ur, ubtmsg)

}

//通知用户抢庄倍数
func (qznn *QZNNGameLogic) notifyUserBetZhuangTimes(ur *gamelogic.GameUser) {

	ubtmsg := new(msg.Message)
	ubtmsg.BClassID = int32(netproto.MessageBClassID_Game)
	ubtmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_UserBetInfoID)
	ubtmsg.MsgData = qznn.getUserBetZhuangTimesMsgData(ur)

	qznn.SendMessageToUser(ur, ubtmsg)

}

//閒家下注按鈕條件
//莊家身上金額需>=下注金額x人數(4)x特殊牌型(3)x莊家所選倍數,
//玩家身上金額需>=下注金額x3x莊家所選倍數(兩個條件需同時成立才可按)
//X=身上金額/莊家所選倍數/特殊牌倍數(最高3) 取整數
//Y=莊家金額/莊家所選倍數/特殊牌倍數(最高3)/4 取整數
//X,Y選擇最小值當可下注數
//低於可下注數的下注按鈕都可以按
func (qznn *QZNNGameLogic) getUserBetTimesMsgData(ur *gamelogic.GameUser) *netproto.QZNN_UserBetInfo {
	ubtdata := new(netproto.QZNN_UserBetInfo)

	//庄家倍数
	//zhuangBetTimes := qznn.computeUserZhuangBetTimes(ur)
	for i := 0; i < len(getNNConfig().zhuangBetTimes); i++ {
		bt := getNNConfig().zhuangBetTimes[i]
		ubtdata.ZhuangBetTimes = append(ubtdata.ZhuangBetTimes, bt)
		//if bt <= zhuangBetTimes {
		//	ubtdata.ZhuangBetTimes = append(ubtdata.ZhuangBetTimes, bt)
		//}
	}

	if qznn.zhuangInfo == nil {
		return ubtdata
	}

	//闲家倍数
	//X=身上金額/莊家所選倍數/特殊牌倍數(最高3) 取整數
	x := int(int(ur.Score) / int(qznn.zhuangInfo.betTimes) / 3)
	//Y=莊家金額/莊家所選倍數/特殊牌倍數(最高3)/4 取整數
	zuser := qznn.GetUserByUserID(qznn.zhuangInfo.userID)
	y := int(int(zuser.Score) / int(qznn.zhuangInfo.betTimes) / 3 / 4)
	min := x
	if y < min {
		min = y
	}

	//xianBetTimes := qznn.computeUserXianBetTimes(ur)
	for i := 0; i < len(getNNConfig().xianBetTimes); i++ {
		bt := getNNConfig().xianBetTimes[i]
		if int(bt) <= min {
			ubtdata.XianBetTimes = append(ubtdata.XianBetTimes, int64(bt))
		}

		//if bt <= xianBetTimes {
		//	ubtdata.XianBetTimes = append(ubtdata.XianBetTimes, int64(bt))
		//}
	}

	return ubtdata
}

func (qznn *QZNNGameLogic) getUserBetZhuangTimesMsgData(ur *gamelogic.GameUser) *netproto.QZNN_UserBetInfo {
	ubtdata := new(netproto.QZNN_UserBetInfo)

	//庄家倍数
	//zhuangBetTimes := qznn.computeUserZhuangBetTimes(ur)
	for i := 0; i < len(getNNConfig().zhuangBetTimes); i++ {
		bt := getNNConfig().zhuangBetTimes[i]
		ubtdata.ZhuangBetTimes = append(ubtdata.ZhuangBetTimes, bt)
		//if bt <= zhuangBetTimes {
		//	ubtdata.ZhuangBetTimes = append(ubtdata.ZhuangBetTimes, bt)
		//}
	}

	//闲家倍数
	//xianBetTimes := qznn.computeUserXianBetTimes(ur)
	for i := 0; i < len(getNNConfig().xianBetTimes); i++ {
		bt := getNNConfig().xianBetTimes[i]
		ubtdata.XianBetTimes = append(ubtdata.XianBetTimes, int64(bt))

		//if bt <= xianBetTimes {
		//	ubtdata.XianBetTimes = append(ubtdata.XianBetTimes, int64(bt))
		//}
	}

	return ubtdata
}

//计算用户庄家底注倍数
func (qznn *QZNNGameLogic) computeUserZhuangBetTimes(ur *gamelogic.GameUser) int32 {
	if qznn.zhuangInfo != nil {
		return 0
	}

	//牛牛倍数
	var nnTimes int64 = 6

	//用户闲家最大底注倍数
	userMaxTimes := int64(getNNConfig().maxXianTimes())

	//计算用户最大倍数
	userMaxBetTimes := int64(ur.Score / (nnTimes * userMaxTimes * int64(qznn.GetServerConfig().BaseScore) * int64(len(qznn.AllGameUser())-1)))

	if userMaxBetTimes >= int64(getNNConfig().maxZhuangTimes()) {
		userMaxBetTimes = int64(getNNConfig().maxZhuangTimes())
	}

	return int32(userMaxBetTimes)
}

//计算用户闲家底注倍数
func (qznn *QZNNGameLogic) computeUserXianBetTimes(ur *gamelogic.GameUser) int32 {
	if qznn.zhuangInfo == nil {
		return 0
	}

	//牛牛倍数
	var nnTimes int64 = 6

	//庄家底注倍数
	zhuanBetTimes := qznn.zhuangInfo.betTimes

	//计算用户最大倍数
	userMaxBetTimes := int64(ur.Score / (nnTimes * int64(zhuanBetTimes) * int64(qznn.GetServerConfig().BaseScore)))

	if userMaxBetTimes >= int64(getNNConfig().maxXianTimes()) {
		userMaxBetTimes = int64(getNNConfig().maxXianTimes())
	}

	return int32(userMaxBetTimes)
}

//发牌
func (qznn *QZNNGameLogic) fapai() {
	//洗牌
	qznn.poker.Shuffle()

	//给所有用户发牌,先发4张
	for _, ur := range qznn.AllGameUser() {
		cards := qznn.poker.Deal(ur.UserID, 4)
		upk := poker.NewUserPoker(ur, cards)
		qznn.gamePoker.AddUserPoker(upk)

		mlog.Debug("向玩家%v发送牌型%v", ur, upk.Cards)
		//通知用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(upk))

		//牌局记录
		qznn.GameLog.Fapai(ur.UserID, cards.String())
	}
	//其他人的都空牌
	for _, ur := range qznn.AllGameUser() {
		upk2 := poker.NewUserEmptyPoker(ur)
		qznn.BoardcastMessage(qznn.getPokerMsg(upk2))
	}

	//向所有旁观用户发送4张空牌
	for _, ur := range qznn.AllLookOnUser() {
		cards := poker.NewEmptyCardGroup(4)
		upk := poker.NewUserPoker(ur, cards)

		mlog.Debug("向旁观玩家%v发送空牌型", ur)
		//通知用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(upk))
	}

	mlog.Debug("发牌完成.%v", qznn.gamePoker)

	qznn.currentWatchAnimationName = qznn.RegisterAnimation(animationName_FirstCard, animationTimeOut_FirstCard, func() {
		//启动确定庄家流程
		qznn.actionName_QiangZhuang = qznn.PlayerActionManager.NewAction(actionName_QiangZhuang, int32(qznn.GetCurrentGameUserCount()), getNNConfig().getServerQiangZhuangTimeout(), qznn.confirmZhuang)

		//广播游戏阶段
		qznn.boardcastGameStep()
	})
}

func (qznn *QZNNGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//检查非游戏玩家
	if !qznn.CheckGameUserMessage(user) {
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//抢庄倍数
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.QZNN_GameMessageClassID_QiangZhuangID) {
		qznn.onUserQiangZhuang(user, gmsg.MsgData.([]byte))
	}
	//闲家倍数
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.QZNN_GameMessageClassID_XianjiaBetTimesID) {
		qznn.onXianJiaBetTimes(user, gmsg.MsgData.([]byte))
	}
	//用户组牌
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.QZNN_GameMessageClassID_CardGroupFinishID) {
		qznn.onUserCardGroupCommited(user)
	}

	//检查用户动画播放
	qznn.CheckUserAnimationCommited(user, gmsg, qznn.currentWatchAnimationName, nil)

}

func (qznn *QZNNGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {

	qznn.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	qznn.sendGameSetting(gameuser)

	//如果非游戏玩家，发送当前场景
	if qznn.IsGaming() && !qznn.IsGameUser(gameuser) {
		qznn.sendGameScene(gameuser)
	}
}

func (qznn *QZNNGameLogic) UserRecover(gameuser *gamelogic.GameUser) {

	qznn.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	qznn.sendGameSetting(gameuser)

	//发送游戏场景
	if qznn.IsGaming() {
		qznn.sendGameScene(gameuser)
	}
}

//确定庄家
func (qznn *QZNNGameLogic) confirmZhuang() {
	//设置游戏阶段
	qznn.gameStep = gameStep_ZhuangConfirm

	//检查是否有玩家未设置抢庄倍数
	for _, ur := range qznn.AllGameUser() {
		if !qznn.qiangZhuangBets.checkUserBet(ur.UserID) {
			qznn.qiangZhuangBets.addMinUserBet(ur.UserID)
			qznn.GameLog.AddGameStep("用户%d自动设置抢庄.倍数:%d", ur.UserID, getNNConfig().minZhuangTimes())
		}
	}

	//设置庄家
	qznn.zhuangInfo = qznn.qiangZhuangBets.getZhuang()

	//广播庄家信息
	qznn.boardcastZhuangInfo()

	qznn.currentWatchAnimationName = qznn.RegisterAnimation(actionName_QiangZhuang, animationTimeout_QiangZhuang, func() {
		//开始闲家倍数超时检测
		qznn.actionName_XianBetTimes = qznn.PlayerActionManager.NewAction(actionName_XianjiaBetTimes, int32(qznn.GetCurrentGameUserCount())-1, getNNConfig().getServerXianJiaBetTimesTimeout(), qznn.confirmXianBetTimes)
		//广播游戏阶段
		qznn.boardcastGameStep()
	})

	//牌局记录
	qznn.GameLog.AddGameStep("庄家确定,%d", qznn.zhuangInfo.userID)

	//通知用户闲家倍数
	for _, ur := range qznn.AllGameUser() {
		qznn.notifyUserBetTimes(ur)
	}
}

//确定闲家倍数
func (qznn *QZNNGameLogic) confirmXianBetTimes() {
	//闲家倍数确定阶段
	qznn.gameStep = gameStep_Xianbet

	for _, ur := range qznn.AllGameUser() {
		//如果是庄家
		if ur.UserID == qznn.zhuangInfo.userID {
			continue
		}

		if !qznn.xianjiaBets.checkUserBet(ur.UserID) {
			qznn.xianjiaBets.addMinUserBet(ur.UserID)
			qznn.GameLog.AddGameStep("用户%d添加闲家最小金额.金额:%d", ur.UserID, getNNConfig().minXianTimes())
		}
	}

	//广播闲家倍数信息
	qznn.boardcastXianjiaInfo()

	//广播游戏阶段
	qznn.boardcastGameStep()

	//发最后兩张牌
	qznn.sendLastCard()
}

//发最后兩张牌
func (qznn *QZNNGameLogic) sendLastCard() {
	for k, ur := range qznn.AllGameUser() {
		//upk := new(poker.UserPoker)

		bowBool, upk := qznn.bowJiangType(ur, k)

		if ur.IsRobot() {
			pk := qznn.poker.Deal(ur.UserID, 1)
			upk = poker.NewUserPoker(ur, pk)
			qznn.gamePoker.AddUserPoker(upk)
		} else if !bowBool {
			pk := make([]*poker.Card, 0)

			niuniurate := 10000 - nnconfig.niuniurate
			rnd := util.RndInt(1, 10001)

			if rnd >= int(niuniurate) {
				niuCardRet := qznn.assignNewNiuNiuCard(qznn.gamePoker.UserPokers[k].Cards, nnCard_NiuNiu, 0, 30)

				pk = niuCardRet.cards.Cards[4:5]
				cg := new(poker.CardGroup)
				cg.Cards = append(cg.Cards, pk...)
				upk = poker.NewUserPoker(ur, cg)
				qznn.gamePoker.AddUserPoker(upk)
			} else {
				niuCardRet := qznn.assignNewNiuNiuCardLessThen(qznn.gamePoker.UserPokers[k].Cards, nnCard_Niu, 9, 30)

				pk = niuCardRet.cards.Cards[4:5]
				cg := new(poker.CardGroup)
				cg.Cards = append(cg.Cards, pk...)
				upk = poker.NewUserPoker(ur, cg)
				qznn.gamePoker.AddUserPoker(upk)
			}
		}

		//发送用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(qznn.gamePoker.GetUserPoker(ur.UserID)))

		//牌局记录
		qznn.GameLog.Fapai(ur.UserID, upk.Cards.String())
	}

	//向所有旁观用户发送5张空牌
	for _, ur := range qznn.AllLookOnUser() {
		cards := poker.NewEmptyCardGroup(5)
		upk := poker.NewUserPoker(ur, cards)

		mlog.Debug("向旁观玩家%v发送空牌型", ur)
		//通知用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(upk))
	}

	//进入组牌阶段
	qznn.gameStep = gameStep_CardGroup

	qznn.getJackpotPool()

	qznn.currentWatchAnimationName = qznn.RegisterAnimation(animationName_SecondCard, animationTimeout_SecondCard, func() {
		//启动组牌超时检测
		qznn.actionName_CardGroup = qznn.PlayerActionManager.NewAction(actionName_CardGroup, int32(qznn.GetCurrentGameUserCount()), getNNConfig().getServerCardGroupTimeout(), qznn.cardGroupFinish)

		//广播游戏阶段
		qznn.boardcastGameStep()

	})
}

//指定產生牛牛類型牌型
func (qznn *QZNNGameLogic) assignNewNiuNiuCard(firstThreeCard *poker.CardGroup, nnRet int, number int, runCount int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(1) {
		nnRet = nnCard_Niu
		number = 9
	}

	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = qznn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet == nnRet && result.number == number {
			break
		}

		if result != nil && result.niuRet <= nnRet {
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[4:5]
			qznn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	if result == nil || len(result.cards.Cards) != int(5) {
		if runCount <= int(0) {
			return result
		}
		runCount--
		qznn.assignNewNiuNiuCard(firstThreeCard, nnCard_WuNiu, 0, runCount)
	}

	return result
}

//指定產生牛牛類型和點數之下的牌型
func (brnn *QZNNGameLogic) assignNewNiuNiuCardLessThen(firstThreeCard *poker.CardGroup, nnRet int, number int, runCount int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(1) {
		nnRet = nnCard_Niu
		number = 9
	}

	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = brnn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet == nnRet && result.number < number {
			break
		}

		if result != nil && result.niuRet > nnRet {
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[4:5]
			brnn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	if result == nil || len(result.cards.Cards) != int(5) {
		if runCount <= int(0) {
			return result
		}
		runCount--
		brnn.assignNewNiuNiuCardLessThen(firstThreeCard, nnCard_WuNiu, 0, runCount)
	}

	return result
}

//隨機產生指定牌的類型
func (qznn *QZNNGameLogic) naturlAssignNiuCard(firstThreeCard *poker.CardGroup) *niuniuCardRet {
	tmpCardGroup := copyPokerCards(firstThreeCard)
	qznn.poker.Shuffle()
	cg := qznn.poker.Deal(0, 1)

	if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil {
		mlog.Warn("沒取到1張牌")
		return nil
	}

	tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
	result := getNiuNiuRet(tmpCardGroup)
	return result
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

//取得資料庫最新JackpotPool遊戲獎池
func (qznn *QZNNGameLogic) getJackpotPool() {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.HallMsgClassID_GetJackpotPoolID)
	qznn.SendRPCMessageFunc(true, pmsg)
}

//用户组牌完成
func (qznn *QZNNGameLogic) onUserCardGroupCommited(user *gamelogic.GameUser) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_CardGroup {
		return
	}

	//广播用户组牌完成消息
	ugMsg := new(msg.Message)
	ugMsg.BClassID = int32(netproto.MessageBClassID_Game)
	ugMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_NotifyUserGroupFinishID)

	ugmsgData := new(netproto.QZNN_NotifyUserGroupFinish)
	ugmsgData.UserID = proto.Int32(user.UserID)
	ugMsg.MsgData = ugmsgData

	qznn.BoardcastMessage(ugMsg)

	//通知组牌动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_CardGroup, user.UserID)
}

//组牌完成
func (qznn *QZNNGameLogic) cardGroupFinish() {
	//进入亮牌阶段
	qznn.gameStep = gameStep_GameRet

	//计算牌局结果
	nnretLogic := newNNGameRet(qznn)

	qznn.gameRet = nnretLogic.getGameRet(qznn)

	nnret := qznn.gameRet

	//牌局记录
	uscores := make(map[int32]int64)

	//写分
	//庄家
	zhuangwinAmount, tax := qznn.GetUserWinAndTax(int64(nnret.ZhuangRet.GetWinLose()))
	zuser := qznn.GetUserByUserID(qznn.zhuangInfo.userID)

	logmsg3 := fmt.Sprintf("userid[%v] 原始輸贏金額[%v] 稅後輸贏[%v] 輸贏的稅[%v] ",
		zuser.UserID, nnret.ZhuangRet.GetWinLose(), zhuangwinAmount, tax)

	mlog.Debug(logmsg3)
	qznn.GameLog.AddGameStep(logmsg3)

	tax, addjackpot := qznn.GetUserTaxAddJackpotPool(tax, zuser.IsRobot())

	logmsg4 := fmt.Sprintf("userid[%v]是否機器人[%v] 輸贏金額[%v] 系統稅收拆 一般稅收增加[%v] 獎金池稅收增加[%v]",
		zuser.UserID, zuser.IsRobot(), zhuangwinAmount, tax, addjackpot)

	mlog.Debug(logmsg4)
	qznn.GameLog.AddGameStep(logmsg4)

	qznn.UpdateGameJackpotPoolFunc(zuser.UserID, addjackpot)
	qznn.UpdateUserScoreFunc(qznn.zhuangInfo.userID, zhuangwinAmount, tax)
	nnret.ZhuangRet.WinLose = proto.Int64(zhuangwinAmount)

	uscores[qznn.zhuangInfo.userID] = zhuangwinAmount
	qznn.GameLog.Context().UserUpdateScore(qznn.zhuangInfo.userID, zhuangwinAmount)

	//其他用户
	for _, uret := range nnret.GetUserRet() {
		userid := uret.GetUser().GetUserID()

		//如果是庄家,则不作处理
		if userid == qznn.zhuangInfo.userID {
			continue
		}

		userwinAmount, tax := qznn.GetUserWinAndTax(int64(uret.GetWinLose()))
		zuser := qznn.GetUserByUserID(userid)

		logmsg1 := fmt.Sprintf("userid[%v] 原始輸贏金額[%v] 稅後輸贏[%v] 輸贏的稅[%v] ",
			zuser.UserID, uret.GetWinLose(), userwinAmount, tax)

		qznn.GameLog.AddGameStep(logmsg1)
		mlog.System(logmsg1)

		tax, addjackpot = qznn.GetUserTaxAddJackpotPool(tax, zuser.IsRobot())
		logmsg2 := fmt.Sprintf("userid[%v]是否機器人[%v] 輸贏金額[%v] 系統稅收拆 一般稅收增加[%v] 獎金池稅收增加[%v]",
			zuser.UserID, zuser.IsRobot(), userwinAmount, tax, addjackpot)
		mlog.Debug(logmsg2)
		qznn.GameLog.AddGameStep(logmsg2)

		qznn.UpdateGameJackpotPoolFunc(userid, addjackpot)
		uret.WinLose = proto.Int64(userwinAmount)
		qznn.UpdateUserScoreFunc(userid, userwinAmount, tax)
		uscores[userid] = userwinAmount
		qznn.GameLog.Context().UserUpdateScore(userid, userwinAmount)
		qznn.GameLog.Context().UserTaxAmount(userid, tax)
	}

	//广播牌局结果消息
	gameretMsg := new(msg.Message)
	gameretMsg.BClassID = int32(netproto.MessageBClassID_Game)
	gameretMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_GameRetID)
	gameretMsg.MsgData = qznn.gameRet

	qznn.BoardcastMessage(gameretMsg)

	//结算动画
	qznn.currentWatchAnimationName = qznn.RegisterAnimation(animationName_GameRet, animationTimeout_GameRet, func() {
		qznn.GameEndFunc()
		qznn.clearGameData()
	})

	//牌局记录
	qznn.GameLog.GameEnd()

}

//用户抢庄
func (qznn *QZNNGameLogic) onUserQiangZhuang(user *gamelogic.GameUser, data []byte) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_Begin {
		return
	}

	ust := &netproto.QZNN_QiangZhuang{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v抢庄数据结构错误%v", user, err)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.Times = int32(math.Abs(float64(*ust.Times)))

	//检查闲家倍数
	check := false
	for _, betTimes := range nnconfig.zhuangBetTimes {
		if *ust.Times == betTimes {
			check = true
		} else if *ust.Times == 0 {
			check = true
		}
	}
	if !check {
		mlog.Warn("用户%v发送了非法的抢庄倍数%d", user, ust.GetTimes())
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//检查庄家倍数
	//uzbtimes := qznn.computeUserZhuangBetTimes(user)
	//if uzbtimes < ust.GetTimes() {
	//	mlog.Error("用户%v发送了非法的抢庄倍数%d,不能超过%d", user, ust.GetTimes(), uzbtimes)
	//	qznn.KickUserFunc(user.UserID, "非法操作")
	//	return
	//}
	//
	////检查庄家倍数
	//if !getNNConfig().checkZhuangBetTimes(ust.GetTimes()) {
	//	mlog.Error("用户%v发送了非法的抢庄倍数%d", user, ust.GetTimes())
	//	qznn.KickUserFunc(user.UserID, "非法操作")
	//	return
	//}

	times := int32(1)
	if ust.GetTimes() > 0 {
		times = ust.GetTimes()
	}

	if qznn.qiangZhuangBets.addUserBet(user.UserID, times, false) {
		//广播抢庄倍数信息
		qznn.boardCastUserZhuangBet(user.UserID, times)
	}

	//牌局记录
	qznn.GameLog.AddGameStep("用户%d抢庄.倍数:%d", user.UserID, times)

	//通知用户动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_QiangZhuang, user.UserID)

}

//闲家倍数
func (qznn *QZNNGameLogic) onXianJiaBetTimes(user *gamelogic.GameUser, data []byte) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_ZhuangConfirm {
		return
	}

	ust := &netproto.QZNN_XianJiaBetTimes{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v闲家倍数数据结构错误%v", user, err)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.Times = int64(math.Abs(float64(*ust.Times)))

	if user.UserID == qznn.zhuangInfo.userID {
		mlog.Error("庄家用户%v发送了闲家倍数信息", user)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return

	}

	//检查闲家倍数
	check := false
	for _, betTimes := range nnconfig.xianBetTimes {
		if *ust.Times == int64(betTimes) {
			check = true
		}
	}
	if !check {
		mlog.Error("用户%v发送了非法的闲家倍数%d", user, ust.GetTimes())
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	////检查闲家倍数
	//uxbtimes := qznn.computeUserXianBetTimes(user)
	//if int64(uxbtimes) < ust.GetTimes() {
	//	mlog.Error("用户%v发送了非法的闲家倍数%d,不能超过%d", user, ust.GetTimes(), uxbtimes)
	//	qznn.KickUserFunc(user.UserID, "非法操作")
	//	return
	//}
	//
	////检查闲家倍数
	//if !getNNConfig().checkXianBetTimes(int32(ust.GetTimes())) {
	//	mlog.Error("用户%v发送了非法的闲家倍数%d", user, ust.GetTimes())
	//	qznn.KickUserFunc(user.UserID, "非法操作")
	//	return
	//}

	//添加用户抢庄信息，并广播抢庄倍数消息
	if qznn.xianjiaBets.addUserBet(user.UserID, int32(ust.GetTimes()), false) {
		//广播闲家倍数信息
		qznn.boardCastUserXianBet(user.UserID, int32(ust.GetTimes()))
	}

	//通知用户动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_XianBetTimes, user.UserID)

	//牌局记录
	qznn.GameLog.AddGameStep("用户%d闲家金额:%d", user.UserID, ust.GetTimes())
}

////用户组牌
//func (qznn *QZNNGameLogic) onUserCardGroup(user *gamelogic.GameUser, data []byte) {
//	//检查游戏阶段
//	if qznn.gameStep != gameStep_Xianbet {
//		return
//	}

//	ust := &netproto.QZNN_UserCardGroup{}

//	err := proto.Unmarshal(data, ust)
//	if err != nil {
//		mlog.Error("解析用户%v组牌数据结构错误%v", user, err)
//		qznn.KickUserFunc(user.UserID, "非法操作")
//		return
//	}

//	//检查用户牌
//	userPoker := qznn.gamePoker.GetUserPoker(user.UserID)
//	if userPoker == nil {
//		mlog.Error("服务器未获取到用户牌型", user)
//		qznn.KickUserFunc(user.UserID, "非法操作")
//		return
//	}
//	userCards := userPoker.Cards
//	netUserCards := qznn.getAllUserCard(ust)

//	if !userCards.IsEqual(netUserCards) {
//		mlog.Error("检查用户%v牌型和服务器不一致,用户:%v,服务器:%v", user, netUserCards, userCards)
//		qznn.KickUserFunc(user.UserID, "非法操作")
//		return
//	}

//	//解析用户牌堆
//	ucs := qznn.getUserCardStack(user, ust)

//	if ucs == nil {
//		return
//	}

//	nret := qznn.checkUserStackNumber(ucs)

//	//添加用户牌堆
//	//qznn.gameCardStack.AddUserCardStack(ucs)

//}

//广播游戏阶段
func (qznn *QZNNGameLogic) boardcastGameStep() {
	qznn.BoardcastMessage(qznn.createGameStepMsg())
}

//构造游戏阶段消息
func (qznn *QZNNGameLogic) createGameStepMsg() *msg.Message {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_GameStepID)

	gsdata := new(netproto.QZNN_GameStep)
	gsdata.Step = proto.Int32(int32(qznn.gameStep))
	gsmsg.MsgData = gsdata

	return gsmsg
}

//广播庄家信息
func (qznn *QZNNGameLogic) boardcastZhuangInfo() {
	//广播庄家信息
	zinfoMsg := qznn.createZhuangInfoMsg()

	qznn.BoardcastMessage(zinfoMsg)
}

//构造庄家信息消息
func (qznn *QZNNGameLogic) createZhuangInfoMsg() *msg.Message {
	zinfoMsg := new(msg.Message)
	zinfoMsg.BClassID = int32(netproto.MessageBClassID_Game)
	zinfoMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_NotifyQiangZhuangBetID)

	zinfodata := new(netproto.QZNN_ConfirmZhuang)
	zinfodata.UserBets = make([]*netproto.QZNN_UserBet, 0)

	if qznn.zhuangInfo != nil {
		zinfodata.ZhuangUserID = proto.Int32(qznn.zhuangInfo.userID)
		zinfodata.Times = proto.Int32(qznn.zhuangInfo.betTimes)
	}

	for _, bt := range qznn.qiangZhuangBets.userBets {
		ubzdata := new(netproto.QZNN_UserBet)
		ubzdata.UserID = proto.Int32(bt.userID)
		ubzdata.Times = proto.Int64(int64(bt.betTimes))
		ubzdata.SystemSet = proto.Bool(bt.sysSet)

		zinfodata.UserBets = append(zinfodata.UserBets, ubzdata)
	}

	zinfoMsg.MsgData = zinfodata

	return zinfoMsg
}

//广播闲家信息
func (qznn *QZNNGameLogic) boardcastXianjiaInfo() {
	xjinfoMsg := qznn.createXianjiaInfoMsg()

	qznn.BoardcastMessage(xjinfoMsg)
}

//构造闲家信息消息
func (qznn *QZNNGameLogic) createXianjiaInfoMsg() *msg.Message {
	xjinfoMsg := new(msg.Message)
	xjinfoMsg.BClassID = int32(netproto.MessageBClassID_Game)
	xjinfoMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_NotifyXianjiaBetID)

	xjinfoData := new(netproto.QZNN_NotifyXianJiaBetTimes)
	xjinfoMsg.MsgData = xjinfoData
	xjinfoData.UserBets = make([]*netproto.QZNN_UserBet, 0)

	for _, xj := range qznn.xianjiaBets.userBets {
		ubzdata := new(netproto.QZNN_UserBet)
		ubzdata.UserID = proto.Int32(xj.userID)
		ubzdata.Times = proto.Int64(int64(xj.betTimes))
		ubzdata.SystemSet = proto.Bool(xj.sysSet)

		xjinfoData.UserBets = append(xjinfoData.UserBets, ubzdata)
	}

	return xjinfoMsg
}

//创建牌型消息
func (qznn *QZNNGameLogic) getPokerMsg(upk *poker.UserPoker) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)

	pmsg.MsgData = upk.GetNetMessageData()

	return pmsg

}

//从网络消息中获取用户牌堆
func (qznn *QZNNGameLogic) getUserCardStack(user *gamelogic.GameUser, ps *netproto.QZNN_UserCardGroup) *poker.UserCardStack {
	cs := poker.NewCardStack()

	if len(ps.CardsThree) != 3 {
		return nil
	}
	if len(ps.CardsTwo) != 2 {
		return nil
	}

	threecg := qznn.getCardGroup(ps.CardsThree)
	twocg := qznn.getCardGroup(ps.CardsTwo)
	cs.AddCardGroup(threecg)
	cs.AddCardGroup(twocg)

	ucs := poker.NewUserCardStack(user, cs)

	return ucs
}

//从网络消息中获取用户牌组
func (qznn *QZNNGameLogic) getCardGroup(npks []*netproto.PokerCard) *poker.CardGroup {
	cg := poker.NewCardGroup(nil)
	for _, npk := range npks {
		pk := poker.NewPokerCardFromNet(npk)

		cg.AddCard(pk)
	}

	return cg
}

//获取用户所有的牌
func (qznn *QZNNGameLogic) getAllUserCard(ps *netproto.QZNN_UserCardGroup) *poker.CardGroup {
	threecg := qznn.getCardGroup(ps.CardsThree)
	twocg := qznn.getCardGroup(ps.CardsTwo)

	threecg.AddCardGroup(twocg)

	return threecg
}

//发送游戏设置
func (qznn *QZNNGameLogic) sendGameSetting(user *gamelogic.GameUser) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.QZNN_GameMessageClassID_GameSettingID)

	nnconfig := getNNConfig()

	gsdata := new(netproto.QZNN_GameSetting)
	for _, zbt := range nnconfig.zhuangBetTimes {
		maxBet := qznn.getZhuangBetBtn(user)
		if zbt <= maxBet {
			gsdata.ZhuangBetTimes = append(gsdata.ZhuangBetTimes, zbt)
		}
	}
	for _, xbt := range nnconfig.xianBetTimes {
		gsdata.XianBetTimes = append(gsdata.XianBetTimes, xbt)
	}
	gsdata.QiangZhuangTimeout = proto.Int32(nnconfig.qiangzhuangTimeout)
	gsdata.XianBetTimeout = proto.Int32(nnconfig.xianJiaBetTimesTimeout)
	gsdata.CardGroupTimeout = proto.Int32(nnconfig.cardGroupTimeout)
	gsdata.BaseScore = proto.Int32(qznn.GetServerConfig().BaseScore)
	gsdata.TableNo = proto.Int32(qznn.TableNo())
	gsdata.GameID = proto.Int32(qznn.GetServerConfig().GameID)
	gsdata.ReadyTimeout = proto.Int32(qznn.GetServerConfig().ReadyWaitTime)

	gsmsg.MsgData = gsdata

	qznn.SendMessageToUser(user, gsmsg)

}

//莊家搶倍按鈕條件
//玩家金額/玩家最低注金額/特殊牌倍數(最高3)/閒家人數(4) 之後取整數 就是最高能搶倍數
func (qznn *QZNNGameLogic) getZhuangBetBtn(user *gamelogic.GameUser) int32 {
	maxBetOption := int32(int32(user.Score) / int32(getNNConfig().xianBetTimes[0]) / 3 / 4)
	return maxBetOption
}

//发送当前游戏场景
func (qznn *QZNNGameLogic) sendGameScene(user *gamelogic.GameUser) {
	sceneMsg := new(msg.Message)
	sceneMsg.BClassID = int32(netproto.MessageBClassID_Game)
	sceneMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_GameSceneID)

	sceneData := new(netproto.QZNN_GameScene)
	sceneMsg.MsgData = sceneData

	//游戏阶段
	sceneData.GameStep = proto.Int32(int32(qznn.gameStep))

	//设置倒计时数字
	var elapsedSeconds int32 = 0
	if qznn.gameStep == gameStep_Begin {
		elapsedSeconds = getNNConfig().qiangzhuangTimeout - qznn.PlayerActionManager.GetActionTimerElapsedSeconds(qznn.actionName_QiangZhuang)

		//设置用户庄家倍数
		sceneData.UserBetInfo = qznn.getUserBetTimesMsgData(user)
	} else if qznn.gameStep == gameStep_ZhuangConfirm {
		elapsedSeconds = getNNConfig().xianJiaBetTimesTimeout - qznn.PlayerActionManager.GetActionTimerElapsedSeconds(qznn.actionName_XianBetTimes)

		//设置用户庄家倍数
		sceneData.UserBetInfo = qznn.getUserBetTimesMsgData(user)
	} else if qznn.gameStep == gameStep_CardGroup {
		elapsedSeconds = getNNConfig().cardGroupTimeout - qznn.PlayerActionManager.GetActionTimerElapsedSeconds(qznn.actionName_CardGroup)
	}
	sceneData.ElapsedSeconds = proto.Int32(elapsedSeconds)

	//游戏开始阶段
	if qznn.gameStep >= gameStep_Begin {
		//游戏中的用户,发送用户牌型
		if qznn.IsGameUser(user) {
			upk := qznn.gamePoker.GetUserPoker(user.UserID)
			if upk != nil {
				sceneData.UserPoker = upk.GetNetMessageData()
			}
		} else { //旁观用户,发送空的牌型
			if len(qznn.gamePoker.UserPokers) > 0 {
				pklen := len(qznn.gamePoker.UserPokers[0].Cards.Cards)

				//构造空的牌型
				cg := poker.NewCardGroup(nil)
				for i := 0; i < pklen; i++ {
					cg.AddCard(poker.NewEmptyPokerCard())
				}

				upk := poker.NewUserPoker(user, cg)
				sceneData.UserPoker = upk.GetNetMessageData()
			}
		}
	}

	//抢庄倍数
	zhuanginfoMsg := qznn.createZhuangInfoMsg()
	sceneData.ZhuangInfo = zhuanginfoMsg.MsgData.(*netproto.QZNN_ConfirmZhuang)

	//闲家倍数阶段,发送闲家倍数信息
	if qznn.gameStep >= gameStep_ZhuangConfirm {
		xianjiaInfoMsg := qznn.createXianjiaInfoMsg()
		sceneData.XianInfo = xianjiaInfoMsg.MsgData.(*netproto.QZNN_NotifyXianJiaBetTimes)
	}

	//组牌阶段
	if qznn.gameStep == gameStep_CardGroup {
		sceneData.FinishedGroupUsers = qznn.PlayerActionManager.GetFinishedActionUsers(qznn.actionName_CardGroup)
	}

	//亮牌阶段
	if qznn.gameStep == gameStep_GameRet {
		sceneData.GameRet = qznn.gameRet
	}

	qznn.SendMessageToUser(user, sceneMsg)
}

//广播用户抢庄倍数信息
func (qznn *QZNNGameLogic) boardCastUserZhuangBet(userid int32, times int32) {
	zbtMsg := new(msg.Message)
	zbtMsg.BClassID = int32(netproto.MessageBClassID_Game)
	zbtMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_OnUserZhuangBetID)

	zbtData := new(netproto.QZNN_UserBet)
	zbtData.UserID = proto.Int32(userid)
	zbtData.Times = proto.Int64(int64(times))
	zbtData.SystemSet = proto.Bool(false)
	zbtMsg.MsgData = zbtData

	qznn.BoardcastMessage(zbtMsg)
}

//广播用户闲家倍数信息
func (qznn *QZNNGameLogic) boardCastUserXianBet(userid int32, times int32) {
	xbtMsg := new(msg.Message)
	xbtMsg.BClassID = int32(netproto.MessageBClassID_Game)
	xbtMsg.SClassID = int32(netproto.QZNN_GameMessageClassID_OnUserXianBetID)

	xbtData := new(netproto.QZNN_UserBet)
	xbtData.UserID = proto.Int32(userid)
	xbtData.Times = proto.Int64(int64(times))
	xbtData.SystemSet = proto.Bool(false)
	xbtMsg.MsgData = xbtData

	qznn.BoardcastMessage(xbtMsg)

}

//换最大牌
func (qznn *QZNNGameLogic) ExchangeMaxCard(userid int32) {
	qznn.LogicBase.ExchangeMaxCard(userid)
	var maxRet *niuniuCardRet
	var maxRetUserID int32

	fpoker := qznn.poker.(*poker.FixedPoker)
	for _, ur := range qznn.AllGameUser() {
		if ur.IsRobot() && ur.UserID != userid {
			continue
		}

		cg := fpoker.GetUserCards(ur.UserID)
		nret := getNiuNiuRet(cg)
		if maxRet == nil {
			maxRet = nret
			maxRetUserID = ur.UserID

			continue
		}

		//比较牌型
		cret := compareCards(maxRet, nret)
		if cret < 0 {
			maxRet = nret
			maxRetUserID = ur.UserID
		}
	}

	//如果没有真实用户
	if maxRetUserID == 0 {
		mlog.Debug("交换机器人%d牌型,但当前没有真实用户参与游戏", userid)
		return
	}

	mlog.Debug("交换用户%d和用户%d牌型", userid, maxRetUserID)
	qznn.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	fpoker.ExChange(userid, maxRetUserID)

}

func (glc *QZNNGameLogic) OnRecvHTTPMessage(req *http.Request) proto.Message {
	data := new(netproto.TipMessage)
	data.Message = proto.String("test")
	data.MsgType = proto.Int(0)
	data.Param = proto.Int(0)

	return data
}

//收到RPC消息解析處理
func (qznn *QZNNGameLogic) OnRecvRPCMessage(m *msg.Message) {
	rm, ok := m.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}
	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	if m.SClassID == int32(netproto.HallMsgClassID_GetJackpotPoolRetID) {
		qznn.setJackpotPoolAmount(m)
	}
}

func (qznn *QZNNGameLogic) setJackpotPoolAmount(m *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.JackpotPoolRet{}
		return mm
	})
	detail := rmm.(*netproto.JackpotPoolRet)
	qznn.jackpotTotalPool = detail.GetAmount()
}
