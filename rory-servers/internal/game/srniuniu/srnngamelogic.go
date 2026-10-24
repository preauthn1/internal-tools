package srnn

import (
	"math"

	proto "example.invalid/golang/protobuf/proto"
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
	animationName_Card        = "card"        //发牌
	animationName_GameRet     = "gameret"     //结算
	animationName_QiangZhuang = "qiangzhuang" //抢庄

	animationTimeOut_GameStart   = 3
	animationTimeOut_Card        = 4
	animationTimeout_GameRet     = 4
	animationTimeout_QiangZhuang = 4
)

//2人牛牛
type SRNNGameLogic struct {
	*gamelogic.DynamicGameLogicBase
	poker                     poker.Poker            //游戏扑克牌
	gamePoker                 *poker.GamePoker       //玩家牌型
	zhuangInfo                *userBet               //庄家信息
	gameStep                  int                    //游戏阶段
	qiangZhuangBets           *zhuangJiaBets         //抢庄倍数
	xianjiaBets               *xianjiaBets           //闲家倍数
	gameCardStack             *poker.GameCardStack   //用户组牌
	gameRet                   *netproto.SRNN_GameRet //游戏结果
	actionName_QiangZhuang    string
	actionName_XianBetTimes   string
	actionName_CardGroup      string
	currentWatchAnimationName string
}

//创建2人牛牛游戏逻辑对象
func NewSRNNGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *SRNNGameLogic {
	glc := new(SRNNGameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)

	//解析配置
	parseNNConfig(cfgs.GetConfig(gamelogic.ConfigFileNameGameServer))

	glc.RegisterAfterGameStartFunc(func() {
		//通知用户抢庄倍数
		for _, ur := range glc.AllGameUser() {
			glc.notifyUserBetTimes(ur)
		}

		//游戏开始
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, func() {
			//启动确定庄家流程
			glc.actionName_QiangZhuang = glc.PlayerActionManager.NewAction(actionName_QiangZhuang, int32(glc.GetCurrentGameUserCount()), getNNConfig().getServerQiangZhuangTimeout(), glc.confirmZhuang)

			glc.boardcastGameStep()
		})
	})

	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})
	return glc
}

//清除游戏数据
func (qznn *SRNNGameLogic) clearGameData() {
	//随机发牌
	if qznn.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		qznn.poker = poker.NewRandomPoker(true, true, 1)
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
func (qznn *SRNNGameLogic) notifyUserBetTimes(ur *gamelogic.GameUser) {

	ubtmsg := new(msg.Message)
	ubtmsg.BClassID = int32(netproto.MessageBClassID_Game)
	ubtmsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_UserBetInfoID)
	ubtmsg.MsgData = qznn.getUserBetTimesMsgData(ur)

	qznn.SendMessageToUser(ur, ubtmsg)

}

//获取用户投注倍数消息数据
func (qznn *SRNNGameLogic) getUserBetTimesMsgData(ur *gamelogic.GameUser) *netproto.SRNN_UserBetInfo {
	ubtdata := new(netproto.SRNN_UserBetInfo)

	//庄家倍数
	zhuangBetTimes := qznn.computeUserZhuangBetTimes(ur)
	for i := 0; i < len(getNNConfig().zhuangBetTimes); i++ {
		bt := getNNConfig().zhuangBetTimes[i]

		if bt <= zhuangBetTimes {
			ubtdata.ZhuangBetTimes = append(ubtdata.ZhuangBetTimes, bt)
		}
	}

	//闲家倍数
	xianBetTimes := qznn.computeUserXianBetTimes(ur)
	for i := 0; i < len(getNNConfig().xianBetTimes); i++ {
		bt := getNNConfig().xianBetTimes[i]

		if bt <= xianBetTimes {
			ubtdata.XianBetTimes = append(ubtdata.XianBetTimes, bt)
		}
	}

	return ubtdata
}

//计算用户庄家底注倍数
func (qznn *SRNNGameLogic) computeUserZhuangBetTimes(ur *gamelogic.GameUser) int32 {
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
func (qznn *SRNNGameLogic) computeUserXianBetTimes(ur *gamelogic.GameUser) int32 {
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
func (qznn *SRNNGameLogic) fapai() {
	//洗牌
	qznn.poker.Shuffle()

	//给所有用户发牌,先发5张
	for _, ur := range qznn.AllGameUser() {
		cards := qznn.poker.Deal(ur.UserID, 5)
		upk := poker.NewUserPoker(ur, cards)

		qznn.gamePoker.AddUserPoker(upk)

		mlog.Debug("向玩家%v发送牌型", ur)
		//通知用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(upk))

		//牌局记录
		qznn.GameLog.Fapai(ur.UserID, cards.String())
	}

	//向所有旁观用户发送5张空牌
	for _, ur := range qznn.AllLookOnUser() {
		cards := poker.NewEmptyCardGroup(5)
		upk := poker.NewUserPoker(ur, cards)

		mlog.Debug("向旁观玩家%v发送空牌型", ur)
		//通知用户牌型
		qznn.SendMessageToUser(ur, qznn.getPokerMsg(upk))
	}

	mlog.Debug("发牌完成.%v", qznn.gamePoker)

	//设置游戏阶段
	qznn.gameStep = gameStep_CardGroup

	qznn.currentWatchAnimationName = qznn.RegisterAnimation(animationName_Card, animationTimeOut_Card, func() {
		//启动组牌超时检测
		qznn.actionName_CardGroup = qznn.PlayerActionManager.NewAction(actionName_CardGroup, int32(qznn.GetCurrentGameUserCount()), getNNConfig().getServerCardGroupTimeout(), qznn.cardGroupFinish)

		//广播游戏阶段
		qznn.boardcastGameStep()
	})
}

func (qznn *SRNNGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//检查非游戏玩家
	if !qznn.CheckGameUserMessage(user) {
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	//抢庄倍数
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SRNN_GameMessageClassID_SRNN_QiangZhuangID) {
		qznn.onUserQiangZhuang(user, gmsg.MsgData.([]byte))
	}
	//闲家倍数
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SRNN_GameMessageClassID_SRNN_XianjiaBetTimesID) {
		qznn.onXianJiaBetTimes(user, gmsg.MsgData.([]byte))
	}
	//用户组牌
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.SRNN_GameMessageClassID_SRNN_CardGroupFinishID) {
		qznn.onUserCardGroupCommited(user)
	}

	//检查用户动画播放
	qznn.CheckUserAnimationCommited(user, gmsg, qznn.currentWatchAnimationName, nil)

}

func (qznn *SRNNGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {

	qznn.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	qznn.sendGameSetting(gameuser)

	//如果非游戏玩家，发送当前场景
	if qznn.IsGaming() && !qznn.IsGameUser(gameuser) {
		qznn.sendGameScene(gameuser)
	}
}

func (qznn *SRNNGameLogic) UserRecover(gameuser *gamelogic.GameUser) {

	qznn.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	qznn.sendGameSetting(gameuser)

	//发送游戏场景
	if qznn.IsGaming() {
		qznn.sendGameScene(gameuser)
	}
}

//确定庄家
func (qznn *SRNNGameLogic) confirmZhuang() {
	//设置游戏阶段
	qznn.gameStep = gameStep_ZhuangConfirm

	//检查是否有玩家未设置抢庄倍数
	for _, ur := range qznn.AllGameUser() {
		if !qznn.qiangZhuangBets.checkUserBet(ur.UserID) {
			qznn.qiangZhuangBets.addMinUserBet(ur.UserID)
		}
	}

	//设置庄家
	qznn.zhuangInfo = qznn.qiangZhuangBets.getZhuang()

	mlog.Debug("确定庄家信息%v", qznn.zhuangInfo)

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
func (qznn *SRNNGameLogic) confirmXianBetTimes() {
	//闲家倍数确定阶段
	qznn.gameStep = gameStep_Xianbet

	for _, ur := range qznn.AllGameUser() {
		//如果是庄家
		if ur.UserID == qznn.zhuangInfo.userID {
			continue
		}

		if !qznn.xianjiaBets.checkUserBet(ur.UserID) {
			qznn.xianjiaBets.addMinUserBet(ur.UserID)
		}
	}

	//广播闲家倍数信息
	qznn.boardcastXianjiaInfo()

	//广播游戏阶段
	qznn.boardcastGameStep()

	//发牌
	qznn.fapai()
}

//用户组牌完成
func (qznn *SRNNGameLogic) onUserCardGroupCommited(user *gamelogic.GameUser) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_CardGroup {
		return
	}

	//广播用户组牌完成消息
	ugMsg := new(msg.Message)
	ugMsg.BClassID = int32(netproto.MessageBClassID_Game)
	ugMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_NotifyUserGroupFinishID)

	ugmsgData := new(netproto.SRNN_NotifyUserGroupFinish)
	ugmsgData.UserID = proto.Int32(user.UserID)
	ugMsg.MsgData = ugmsgData

	qznn.BoardcastMessage(ugMsg)

	//通知组牌动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_CardGroup, user.UserID)
}

//组牌完成
func (qznn *SRNNGameLogic) cardGroupFinish() {
	//进入亮牌阶段
	qznn.gameStep = gameStep_GameRet

	//计算牌局结果
	nnretLogic := newNNGameRet(qznn)

	qznn.gameRet = nnretLogic.getGameRet()

	nnret := qznn.gameRet

	//牌局记录
	uscores := make(map[int32]int64)

	//写分
	//庄家
	zhuangwinAmount, tax := qznn.GetUserWinAndTax(int64(nnret.ZhuangRet.GetWinLose()))
	qznn.UpdateUserScoreFunc(qznn.zhuangInfo.userID, zhuangwinAmount, tax)
	nnret.ZhuangRet.WinLose = proto.Int64(zhuangwinAmount)
	zhuangInfo := qznn.GetUserByUserID(qznn.zhuangInfo.userID)
	var flowAmount, betAmount, payAmount int64
	flowAmount = util.ABSInt64(zhuangwinAmount)
	betAmount = util.ABSInt64(zhuangwinAmount) + tax
	payAmount = betAmount + zhuangwinAmount
	zhuangInfo.FlowAmount += flowAmount
	zhuangInfo.BetAmount += betAmount
	zhuangInfo.PayAmount += payAmount
	retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
	if zhuangInfo.IsRobot() {
		retmsg = "[機器人]" + retmsg
	}
	qznn.GameLog.AddGameStep(retmsg, zhuangInfo.UserID, flowAmount, betAmount, payAmount)
	uscores[qznn.zhuangInfo.userID] = zhuangwinAmount

	//其他用户
	for _, uret := range nnret.GetUserRet() {
		userid := uret.GetUser().GetUserID()
		//如果是庄家,则不作处理
		if userid == qznn.zhuangInfo.userID {
			continue
		}
		var flowAmount, betAmount, payAmount int64
		user := qznn.GetUserByUserID(userid)
		userwinAmount, tax := qznn.GetUserWinAndTax(int64(uret.GetWinLose()))
		flowAmount = util.ABSInt64(userwinAmount)
		betAmount = util.ABSInt64(userwinAmount) + tax
		payAmount = betAmount + userwinAmount

		retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
		if user.IsRobot() {
			retmsg = "[機器人]" + retmsg
		}
		qznn.GameLog.AddGameStep(retmsg, user.UserID, flowAmount, betAmount, payAmount)

		uret.WinLose = proto.Int64(userwinAmount)
		qznn.UpdateUserScoreFunc(userid, userwinAmount, tax)

		user.FlowAmount += flowAmount
		user.BetAmount += betAmount
		user.PayAmount += payAmount

		uscores[userid] = userwinAmount

		qznn.GameLog.Context().UserUpdateScore(userid, userwinAmount)
		qznn.GameLog.Context().UserBetScore(userid, betAmount)
		qznn.GameLog.Context().UserPayoutScore(userid, payAmount)
		qznn.GameLog.Context().UserValidScore(userid, flowAmount)
		qznn.GameLog.Context().UserTaxAmount(userid, tax)

	}

	//广播牌局结果消息
	gameretMsg := new(msg.Message)
	gameretMsg.BClassID = int32(netproto.MessageBClassID_Game)
	gameretMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_GameRetID)
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
func (qznn *SRNNGameLogic) onUserQiangZhuang(user *gamelogic.GameUser, data []byte) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_Begin {
		return
	}

	ust := &netproto.SRNN_QiangZhuang{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v抢庄数据结构错误%v", user, err)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.Times = int32(math.Abs(float64(*ust.Times)))

	//检查庄家倍数
	uzbtimes := qznn.computeUserZhuangBetTimes(user)
	if uzbtimes < ust.GetTimes() {
		mlog.Error("用户%v发送了非法的抢庄倍数%d,不能超过%d", user, ust.GetTimes(), uzbtimes)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//检查庄家倍数
	if !getNNConfig().checkZhuangBetTimes(ust.GetTimes()) {
		mlog.Error("用户%v发送了非法的抢庄倍数%d", user, ust.GetTimes())
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	if qznn.qiangZhuangBets.addUserBet(user.UserID, ust.GetTimes(), false) {
		//广播抢庄倍数信息
		qznn.boardCastUserZhuangBet(user.UserID, ust.GetTimes())
	}

	//牌局记录
	qznn.GameLog.AddGameStep("用户%d抢庄.倍数:%d", user.UserID, ust.GetTimes())

	//通知用户动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_QiangZhuang, user.UserID)

}

//闲家倍数
func (qznn *SRNNGameLogic) onXianJiaBetTimes(user *gamelogic.GameUser, data []byte) {
	//检查游戏阶段
	if qznn.gameStep != gameStep_ZhuangConfirm {
		return
	}

	ust := &netproto.SRNN_XianJiaBetTimes{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v闲家倍数数据结构错误%v", user, err)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.Times = int32(math.Abs(float64(*ust.Times)))

	if user.UserID == qznn.zhuangInfo.userID {
		mlog.Error("庄家用户%v发送了闲家倍数信息", user)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return

	}

	//检查闲家倍数
	uxbtimes := qznn.computeUserXianBetTimes(user)
	if uxbtimes < ust.GetTimes() {
		mlog.Error("用户%v发送了非法的闲家倍数%d,不能超过%d", user, ust.GetTimes(), uxbtimes)
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//检查闲家倍数
	if !getNNConfig().checkXianBetTimes(ust.GetTimes()) {
		mlog.Error("用户%v发送了非法的闲家倍数%d", user, ust.GetTimes())
		qznn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//添加用户抢庄信息，并广播抢庄倍数消息
	if qznn.xianjiaBets.addUserBet(user.UserID, ust.GetTimes(), false) {
		//广播闲家倍数信息
		qznn.boardCastUserXianBet(user.UserID, ust.GetTimes())
	}

	//通知用户动作
	qznn.PlayerActionManager.OnUserAction(qznn.actionName_XianBetTimes, user.UserID)

	//牌局记录
	qznn.GameLog.AddGameStep("用户%d闲家倍数:%d", user.UserID, ust.GetTimes())
}

//广播游戏阶段
func (qznn *SRNNGameLogic) boardcastGameStep() {
	qznn.BoardcastMessage(qznn.createGameStepMsg())
}

//构造游戏阶段消息
func (qznn *SRNNGameLogic) createGameStepMsg() *msg.Message {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_GameStepID)

	gsdata := new(netproto.SRNN_GameStep)
	gsdata.Step = proto.Int32(int32(qznn.gameStep))
	gsmsg.MsgData = gsdata

	return gsmsg
}

//广播庄家信息
func (qznn *SRNNGameLogic) boardcastZhuangInfo() {
	//广播庄家信息
	zinfoMsg := qznn.createZhuangInfoMsg()

	qznn.BoardcastMessage(zinfoMsg)
}

//构造庄家信息消息
func (qznn *SRNNGameLogic) createZhuangInfoMsg() *msg.Message {
	zinfoMsg := new(msg.Message)
	zinfoMsg.BClassID = int32(netproto.MessageBClassID_Game)
	zinfoMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_NotifyQiangZhuangBetID)

	zinfodata := new(netproto.SRNN_ConfirmZhuang)
	zinfodata.UserBets = make([]*netproto.SRNN_UserBet, 0)

	if qznn.zhuangInfo != nil {
		zinfodata.ZhuangUserID = proto.Int32(qznn.zhuangInfo.userID)
		zinfodata.Times = proto.Int32(qznn.zhuangInfo.betTimes)
	}

	for _, bt := range qznn.qiangZhuangBets.userBets {
		ubzdata := new(netproto.SRNN_UserBet)
		ubzdata.UserID = proto.Int32(bt.userID)
		ubzdata.Times = proto.Int32(bt.betTimes)
		ubzdata.SystemSet = proto.Bool(bt.sysSet)

		zinfodata.UserBets = append(zinfodata.UserBets, ubzdata)
	}

	zinfoMsg.MsgData = zinfodata

	return zinfoMsg
}

//广播闲家信息
func (qznn *SRNNGameLogic) boardcastXianjiaInfo() {
	xjinfoMsg := qznn.createXianjiaInfoMsg()

	qznn.BoardcastMessage(xjinfoMsg)
}

//构造闲家信息消息
func (qznn *SRNNGameLogic) createXianjiaInfoMsg() *msg.Message {
	xjinfoMsg := new(msg.Message)
	xjinfoMsg.BClassID = int32(netproto.MessageBClassID_Game)
	xjinfoMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_NotifyXianjiaBetID)

	xjinfoData := new(netproto.SRNN_NotifyXianJiaBetTimes)
	xjinfoMsg.MsgData = xjinfoData
	xjinfoData.UserBets = make([]*netproto.SRNN_UserBet, 0)

	for _, xj := range qznn.xianjiaBets.userBets {
		ubzdata := new(netproto.SRNN_UserBet)
		ubzdata.UserID = proto.Int32(xj.userID)
		ubzdata.Times = proto.Int32(xj.betTimes)
		ubzdata.SystemSet = proto.Bool(xj.sysSet)

		xjinfoData.UserBets = append(xjinfoData.UserBets, ubzdata)
	}

	return xjinfoMsg
}

//创建牌型消息
func (qznn *SRNNGameLogic) getPokerMsg(upk *poker.UserPoker) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)

	pmsg.MsgData = upk.GetNetMessageData()

	return pmsg

}

//从网络消息中获取用户牌堆
func (qznn *SRNNGameLogic) getUserCardStack(user *gamelogic.GameUser, ps *netproto.SRNN_UserCardGroup) *poker.UserCardStack {
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
func (qznn *SRNNGameLogic) getCardGroup(npks []*netproto.PokerCard) *poker.CardGroup {
	cg := poker.NewCardGroup(nil)
	for _, npk := range npks {
		pk := poker.NewPokerCardFromNet(npk)

		cg.AddCard(pk)
	}

	return cg
}

//获取用户所有的牌
func (qznn *SRNNGameLogic) getAllUserCard(ps *netproto.SRNN_UserCardGroup) *poker.CardGroup {
	threecg := qznn.getCardGroup(ps.CardsThree)
	twocg := qznn.getCardGroup(ps.CardsTwo)

	threecg.AddCardGroup(twocg)

	return threecg
}

//发送游戏设置
func (qznn *SRNNGameLogic) sendGameSetting(user *gamelogic.GameUser) {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_GameSettingID)

	nnconfig := getNNConfig()

	gsdata := new(netproto.SRNN_GameSetting)
	for _, zbt := range nnconfig.zhuangBetTimes {
		gsdata.ZhuangBetTimes = append(gsdata.ZhuangBetTimes, zbt)
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

//发送当前游戏场景
func (qznn *SRNNGameLogic) sendGameScene(user *gamelogic.GameUser) {
	sceneMsg := new(msg.Message)
	sceneMsg.BClassID = int32(netproto.MessageBClassID_Game)
	sceneMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_GameSceneID)

	sceneData := new(netproto.SRNN_GameScene)
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
	sceneData.ZhuangInfo = zhuanginfoMsg.MsgData.(*netproto.SRNN_ConfirmZhuang)

	//闲家倍数阶段,发送闲家倍数信息
	if qznn.gameStep >= gameStep_ZhuangConfirm {
		xianjiaInfoMsg := qznn.createXianjiaInfoMsg()
		sceneData.XianInfo = xianjiaInfoMsg.MsgData.(*netproto.SRNN_NotifyXianJiaBetTimes)
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
func (qznn *SRNNGameLogic) boardCastUserZhuangBet(userid int32, times int32) {
	zbtMsg := new(msg.Message)
	zbtMsg.BClassID = int32(netproto.MessageBClassID_Game)
	zbtMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_OnUserZhuangBetID)

	zbtData := new(netproto.SRNN_UserBet)
	zbtData.UserID = proto.Int32(userid)
	zbtData.Times = proto.Int32(times)
	zbtData.SystemSet = proto.Bool(false)
	zbtMsg.MsgData = zbtData

	qznn.BoardcastMessage(zbtMsg)
}

//广播用户闲家倍数信息
func (qznn *SRNNGameLogic) boardCastUserXianBet(userid int32, times int32) {
	xbtMsg := new(msg.Message)
	xbtMsg.BClassID = int32(netproto.MessageBClassID_Game)
	xbtMsg.SClassID = int32(netproto.SRNN_GameMessageClassID_SRNN_OnUserXianBetID)

	xbtData := new(netproto.SRNN_UserBet)
	xbtData.UserID = proto.Int32(userid)
	xbtData.Times = proto.Int32(times)
	xbtData.SystemSet = proto.Bool(false)
	xbtMsg.MsgData = xbtData

	qznn.BoardcastMessage(xbtMsg)

}

//换最大牌
func (qznn *SRNNGameLogic) ExchangeMaxCard(userid int32) {
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
