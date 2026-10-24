package tbnn

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
	gameStep_ANI_Start = iota //开始动画
	gameStep_BaseScore        //选择底注
	gameStep_ANI_Fapai        //发牌动画
	gameStep_CardGroup        //组牌
	gameStep_GameRet          //亮牌
)

const (
	actionName_BaseScore = "actionName_basescore" //选择底注
	actionName_CardGroup = "actionName_CardGroup" //组牌
	actionName_GameRet   = "actionName_GameRet"   //亮牌

)

const (
	animationName_GameStart = "gamestart" //开始动画
	animationName_Fapai     = "fapai"     //发牌动画
	animationName_GameRet   = "gameret"

	animationTimeOut_GameStart = 2 //开始动画
	animationTimeout_Fapai     = 4 //发牌动画
	animationTimeout_GameRet   = 2 //亮牌
)

//抢庄牛牛
type TBNNGameLogic struct {
	*gamelogic.DynamicGameLogicBase
	poker                  poker.Poker            //游戏扑克牌
	gamePoker              *poker.GamePoker       //玩家牌型
	gameStep               int                    //游戏阶段
	gameBet                *gameBet               //用户底注
	cardGroupFinishedUsers []int32                //用户组牌
	gameRet                *netproto.TBNN_GameRet //游戏结果

	currentWatchAnimationName string //当前动画名称

	actionName_BaseScore string //选择底注
	actionName_CardGroup string //组牌
}

//创建抢庄牛牛游戏逻辑对象
func NewTBNNGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *TBNNGameLogic {
	glc := new(TBNNGameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)

	//解析牛牛配置
	parseNNConfig(cfgs.GetConfig(gamelogic.ConfigFileNameGameServer))

	glc.RegisterAfterGameStartFunc(func() {

		//通知用户倍数
		for _, ur := range glc.AllGameUser() {
			glc.notifyUserBetTimes(ur)
		}

		//发牌动画
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, glc.baseScoreSetting)
	})

	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})
	return glc
}

func (tbnn *TBNNGameLogic) clearGameData() {
	//随机发牌
	if tbnn.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		tbnn.poker = poker.NewRandomPoker(true, true, 1)
		tbnn.poker = poker.NewFixedPoker(tbnn.AllTableUser(), tbnn.poker.(*poker.RandomPoker), 5)
	} else { //发牌器发牌
		tbnn.poker = poker.NewFilePoker(tbnn.Configs, tbnn.AllTableUser())
		tbnn.poker = poker.NewFixedFilePoker(tbnn.poker.(*poker.FilePoker))
	}

	tbnn.gamePoker = poker.NewGamePoker() //创建游戏牌型
	tbnn.gameStep = gameStep_ANI_Start    //设置游戏阶段
	tbnn.gameBet = newGameBet()           //底注
	tbnn.cardGroupFinishedUsers = make([]int32, 0)
	tbnn.gameRet = nil

	tbnn.actionName_BaseScore = ""
	tbnn.actionName_CardGroup = ""

}

//通知用户牛牛倍数
func (tbnn *TBNNGameLogic) notifyUserBetTimes(ur *gamelogic.GameUser) {
	//计算用户底注倍数
	userMaxBetTimes := tbnn.computeUserBetTimes(ur)

	ubtmsg := new(msg.Message)
	ubtmsg.BClassID = int32(netproto.MessageBClassID_Game)
	ubtmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_UserBetInfoID)

	ubtdata := new(netproto.TBNN_UserBetInfo)
	ubtmsg.MsgData = ubtdata

	for i := 0; i < len(getNNConfig().baseScoreTimes); i++ {
		bt := getNNConfig().baseScoreTimes[i]

		if bt <= userMaxBetTimes {
			ubtdata.BetTimes = append(ubtdata.BetTimes, bt)
		}
	}

	tbnn.SendMessageToUser(ur, ubtmsg)

}

//计算用户底注倍数
func (tbnn *TBNNGameLogic) computeUserBetTimes(ur *gamelogic.GameUser) int32 {
	//牛牛倍数
	var nnTimes int64 = 6

	//用户最大底注备注
	userMaxTimes := int64(getNNConfig().maxBaseScoreTimes())

	//计算用户最大倍数
	userMaxBetTimes := int64(ur.Score / (nnTimes * userMaxTimes * int64(tbnn.GetServerConfig().BaseScore)))

	if userMaxBetTimes >= int64(getNNConfig().maxBaseScoreTimes()) {
		userMaxBetTimes = int64(getNNConfig().maxBaseScoreTimes())
	}

	return int32(userMaxBetTimes)
}

//用户进入桌子
func (tbnn *TBNNGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {

	tbnn.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	tbnn.sendGameSetting(gameuser)

	//如果非游戏玩家，发送当前场景
	if tbnn.IsGaming() && !tbnn.IsGameUser(gameuser) {
		tbnn.sendGameScene(gameuser)
	}
}

//用户断线重连
func (tbnn *TBNNGameLogic) UserRecover(gameuser *gamelogic.GameUser) {

	tbnn.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	tbnn.sendGameSetting(gameuser)

	//发送游戏场景
	if tbnn.IsGaming() {
		tbnn.sendGameScene(gameuser)
	}
}

func (tbnn *TBNNGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//检查非游戏玩家
	if !tbnn.CheckGameUserMessage(user) {
		tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//底注
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.TBNN_GameMessageClassID_TBNN_UserBaseScoreID) {
		tbnn.onUserBaseScore(user, gmsg.MsgData.([]byte))
	}
	//组牌
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.TBNN_GameMessageClassID_TBNN_UserCardGroupID) {
		tbnn.onUserCardGroup(user, gmsg.MsgData.([]byte))
	}

	//检查用户动画播放
	tbnn.CheckUserAnimationCommited(user, gmsg, tbnn.currentWatchAnimationName, nil)

}

//选择底注
func (tbnn *TBNNGameLogic) baseScoreSetting() {
	//设置游戏阶段
	tbnn.gameStep = gameStep_BaseScore

	tbnn.actionName_BaseScore = tbnn.PlayerActionManager.NewAction(actionName_BaseScore, int32(tbnn.GetCurrentGameUserCount()), getNNConfig().baseScoreTimeout, tbnn.baseScoreConfirm)

	//广播游戏阶段
	tbnn.boardcastGameStep()
}

//用户底注
func (tbnn *TBNNGameLogic) onUserBaseScore(user *gamelogic.GameUser, data []byte) {
	if tbnn.gameStep != gameStep_BaseScore {
		mlog.Debug("用户%v发送底注消息,但当前游戏阶段是%d", user, tbnn.gameStep)
		//tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	ust := &netproto.TBNN_UserBaseScore{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v底注数据结构错误%v", user, err)
		tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.BaseScore = int32(math.Abs(float64(*ust.BaseScore)))

	//检查用户底注倍数
	ubtimes := tbnn.computeUserBetTimes(user)
	if ubtimes < ust.GetBaseScore() {
		mlog.Error("用户%v提交了非法的底注%d。不能超过%d", user, ust.GetBaseScore(), ubtimes)
		tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	if !tbnn.gameBet.checkBetTimes(ust.GetBaseScore()) {
		mlog.Error("用户%v提交了非法的底注%d", user, ust.GetBaseScore())
		tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//添加用户底注
	if tbnn.gameBet.addUserBet(user.UserID, ust.GetBaseScore()) {
		tbnn.boardCastUserBaseScore(user.UserID, ust.GetBaseScore())
	}

	//通知用户动作
	tbnn.PlayerActionManager.OnUserAction(tbnn.actionName_BaseScore, user.UserID)

	//牌局记录
	tbnn.GameLog.AddGameStep("用户%d选择底注%d", user.UserID, ust.GetBaseScore())

}

//底注确定
func (tbnn *TBNNGameLogic) baseScoreConfirm() {
	for _, ur := range tbnn.AllGameUser() {
		tbnn.gameBet.ensureUserBet(ur.UserID)
	}

	//广播游戏底注
	tbnn.boardCastGameBaseScore()

	//设置游戏阶段
	tbnn.gameStep = gameStep_ANI_Fapai
	//广播游戏阶段
	tbnn.boardcastGameStep()

	//发牌
	tbnn.fapai()
}

//发牌
func (tbnn *TBNNGameLogic) fapai() {

	//洗牌
	tbnn.poker.Shuffle()

	//给所有游戏用户发牌
	for _, ur := range tbnn.AllGameUser() {
		cards := tbnn.poker.Deal(ur.UserID, 5)
		upk := poker.NewUserPoker(ur, cards)

		tbnn.gamePoker.AddUserPoker(upk)

		mlog.Debug("向玩家%v发送牌型", ur)
		//通知用户牌型
		tbnn.SendMessageToUser(ur, tbnn.getPokerMsg(upk))

		//牌局记录
		tbnn.GameLog.Fapai(ur.UserID, cards.String())
	}

	//向所有旁观用户发送5张空牌
	for _, ur := range tbnn.AllLookOnUser() {
		cards := poker.NewEmptyCardGroup(5)
		upk := poker.NewUserPoker(ur, cards)

		mlog.Debug("向旁观玩家%v发送空牌型", ur)
		//通知用户牌型
		tbnn.SendMessageToUser(ur, tbnn.getPokerMsg(upk))
	}

	mlog.Debug("发牌完成.%v", tbnn.gamePoker)

	//设置游戏阶段
	tbnn.gameStep = gameStep_CardGroup

	tbnn.currentWatchAnimationName = tbnn.RegisterAnimation(animationName_Fapai, animationTimeout_Fapai, func() {
		//组牌
		tbnn.actionName_CardGroup = tbnn.PlayerActionManager.NewAction(actionName_CardGroup, int32(tbnn.GetCurrentGameUserCount()), getNNConfig().cardGroupTimeout, tbnn.cardGroupFinish)

		//广播游戏阶段
		tbnn.boardcastGameStep()
	})
}

//用户组牌
func (tbnn *TBNNGameLogic) onUserCardGroup(user *gamelogic.GameUser, data []byte) {
	if tbnn.gameStep != gameStep_CardGroup {
		mlog.Debug("用户%v发送组牌消息,但当前游戏阶段是%d", user, tbnn.gameStep)
		//tbnn.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//检查用户是否已经提交组牌完成
	for _, uid := range tbnn.cardGroupFinishedUsers {
		if uid == user.UserID {
			return
		}
	}

	//记录组牌完成用户
	tbnn.cardGroupFinishedUsers = append(tbnn.cardGroupFinishedUsers, user.UserID)

	//广播用户组牌完成
	tbnn.boardCastUserCardGroup(user)

	//收到组牌动作
	tbnn.PlayerActionManager.OnUserAction(tbnn.actionName_CardGroup, user.UserID)
}

//广播用户组牌
func (tbnn *TBNNGameLogic) boardCastUserCardGroup(user *gamelogic.GameUser) {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_NotifyUserCardGroupID)

	ugData := new(netproto.TBNN_NotifyUserCardGroup)
	ugData.UserID = proto.Int32(user.UserID)
	gmsg.MsgData = ugData

	tbnn.BoardcastMessage(gmsg)

}

//组牌完成
func (tbnn *TBNNGameLogic) cardGroupFinish() {
	//亮牌阶段
	tbnn.gameStep = gameStep_GameRet
	//广播游戏阶段
	tbnn.boardcastGameStep()

	//结算
	gameret := newNNGameRet(tbnn).getGameRet()

	//牌局结果
	uscores := make(map[int32]int64)
	for _, ret := range gameret.UserRet {
		var flowAmount, betAmount, payAmount int64
		//计算税收
		amount, tax := tbnn.GetUserWinAndTax(ret.GetWinLose())
		flowAmount = util.ABSInt64(amount)
		betAmount = util.ABSInt64(amount) + tax
		payAmount = betAmount + amount

		user := tbnn.GetUserByUserID(ret.GetUserID())

		retmsg := "用戶%d, 流水%d, 有效投注 %d, 返還金額 %d"
		if user.IsRobot() {
			retmsg = "[機器人]" + retmsg
		}
		tbnn.GameLog.AddGameStep(retmsg, user.UserID, flowAmount, betAmount, payAmount)
		//写分
		tbnn.UpdateUserScoreFunc(ret.GetUserID(), int64(amount), int64(tax))

		tbnn.GameLog.Context().UserUpdateScore(ret.GetUserID(), amount)
		tbnn.GameLog.Context().UserBetScore(ret.GetUserID(), betAmount)
		tbnn.GameLog.Context().UserPayoutScore(ret.GetUserID(), payAmount)
		tbnn.GameLog.Context().UserValidScore(ret.GetUserID(), flowAmount)
		tbnn.GameLog.Context().UserTaxAmount(ret.GetUserID(), tax)

		ret.WinLose = proto.Int64(amount)

		//牌局结果
		uscores[ret.GetUserID()] = amount

	}

	//设置游戏结果
	tbnn.gameRet = gameret

	//广播游戏结果
	tbnn.boardCastGameRet()

	//结算动画
	tbnn.currentWatchAnimationName = tbnn.RegisterAnimation(animationName_GameRet, animationTimeout_GameRet, func() {
		tbnn.GameEndFunc()
		tbnn.clearGameData()
	})

	//牌局记录
	tbnn.GameLog.GameEnd()

}

//广播游戏结果
func (tbnn *TBNNGameLogic) boardCastGameRet() {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_GameRetID)

	gmsg.MsgData = tbnn.gameRet
	tbnn.BoardcastMessage(gmsg)
}

//广播用户底注信息
func (tbnn *TBNNGameLogic) boardCastUserBaseScore(userid int32, baseScore int32) {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_NotifyUserBaseScoreID)

	ubtData := new(netproto.TBNN_NotifyUserBaseScore)
	ubtData.UserID = proto.Int32(userid)
	ubtData.BaseScore = proto.Int32(baseScore)
	gmsg.MsgData = ubtData

	tbnn.BoardcastMessage(gmsg)
}

//广播游戏底注信息
func (tbnn *TBNNGameLogic) boardCastGameBaseScore() {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_NotifyGameBaseScoreID)
	gmsg.MsgData = tbnn.getGameBaseScoreData()

	tbnn.BoardcastMessage(gmsg)
}

//构造游戏底注数据
func (tbnn *TBNNGameLogic) getGameBaseScoreData() *netproto.TBNN_NotifyGameBaseScore {
	gbtData := new(netproto.TBNN_NotifyGameBaseScore)

	for _, ubt := range tbnn.gameBet.userbets {
		ubtData := new(netproto.TBNN_NotifyUserBaseScore)
		ubtData.UserID = proto.Int32(ubt.userID)
		ubtData.BaseScore = proto.Int32(ubt.betTimes)

		gbtData.GameBet = append(gbtData.GameBet, ubtData)
	}

	return gbtData
}

//广播游戏阶段
func (tbnn *TBNNGameLogic) boardcastGameStep() {
	tbnn.BoardcastMessage(tbnn.createGameStepMsg())
}

//构造游戏阶段消息
func (tbnn *TBNNGameLogic) createGameStepMsg() *msg.Message {
	gsmsg := new(msg.Message)
	gsmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gsmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_GameStepID)

	gsdata := new(netproto.TBNN_GameStep)
	gsdata.Step = proto.Int32(int32(tbnn.gameStep))
	gsmsg.MsgData = gsdata

	return gsmsg
}

//创建牌型消息
func (tbnn *TBNNGameLogic) getPokerMsg(upk *poker.UserPoker) *msg.Message {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)

	pmsg.MsgData = upk.GetNetMessageData()

	return pmsg
}

//发送游戏配置
func (tbnn *TBNNGameLogic) sendGameSetting(ur *gamelogic.GameUser) {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_GameConfigID)

	stData := new(netproto.TBNN_GameSetting)
	stData.BaseScoreTimeout = proto.Int32(getNNConfig().baseScoreTimeout)
	stData.CardGroupTimeout = proto.Int32(getNNConfig().cardGroupTimeout)
	stData.BaseScore = proto.Int32(tbnn.GetServerConfig().BaseScore)
	stData.TableNo = proto.Int32(tbnn.TableNo())
	stData.GameID = proto.Int32(tbnn.GetServerConfig().GameID)
	stData.ReadyTimeout = proto.Int32(tbnn.GetServerConfig().ReadyWaitTime)
	stData.BaseScoreTimes = getNNConfig().baseScoreTimes

	gmsg.MsgData = stData

	tbnn.BoardcastMessage(gmsg)
}

//发送当前游戏场景
func (tbnn *TBNNGameLogic) sendGameScene(user *gamelogic.GameUser) {
	sceneMsg := new(msg.Message)
	sceneMsg.BClassID = int32(netproto.MessageBClassID_Game)
	sceneMsg.SClassID = int32(netproto.TBNN_GameMessageClassID_TBNN_GameSceneID)

	sceneData := new(netproto.TBNN_GameScene)
	sceneMsg.MsgData = sceneData

	//游戏阶段
	sceneData.GameStep = proto.Int32(int32(tbnn.gameStep))

	//设置倒计时数字
	var elapsedSeconds int32 = 0
	//选择底注阶段
	if tbnn.gameStep == gameStep_BaseScore {
		elapsedSeconds = getNNConfig().baseScoreTimeout - tbnn.PlayerActionManager.GetActionTimerElapsedSeconds(tbnn.actionName_BaseScore)

	} else if tbnn.gameStep == gameStep_CardGroup {
		elapsedSeconds = getNNConfig().cardGroupTimeout - tbnn.PlayerActionManager.GetActionTimerElapsedSeconds(tbnn.actionName_CardGroup)
	}
	sceneData.ElapsedSeconds = proto.Int32(elapsedSeconds)

	//设置用户底注倍数
	userMaxBetTimes := tbnn.computeUserBetTimes(user)
	for i := 0; i < len(getNNConfig().baseScoreTimes); i++ {
		bt := getNNConfig().baseScoreTimes[i]

		if bt <= userMaxBetTimes {
			sceneData.BetTimes = append(sceneData.BetTimes, bt)
		}
	}

	//底注
	if tbnn.gameStep >= gameStep_BaseScore {
		sceneData.GameBaseScore = tbnn.getGameBaseScoreData()
	}

	//发牌
	if tbnn.gameStep >= gameStep_ANI_Fapai {
		//游戏中的用户,发送用户牌型
		if tbnn.IsGameUser(user) {
			upk := tbnn.gamePoker.GetUserPoker(user.UserID)
			if upk != nil {
				sceneData.UserPoker = upk.GetNetMessageData()
			}
		} else { //旁观用户,发送空的牌型
			if len(tbnn.gamePoker.UserPokers) > 0 {
				pklen := len(tbnn.gamePoker.UserPokers[0].Cards.Cards)

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

	//组牌阶段
	if tbnn.gameStep == gameStep_CardGroup {
		sceneData.FinishedGroupUsers = tbnn.PlayerActionManager.GetFinishedActionUsers(tbnn.actionName_CardGroup)
	}

	//亮牌阶段
	if tbnn.gameStep == gameStep_GameRet {
		sceneData.GameRet = tbnn.gameRet
	}

	tbnn.SendMessageToUser(user, sceneMsg)
}

//换最大牌
func (tbnn *TBNNGameLogic) ExchangeMaxCard(userid int32) {
	tbnn.LogicBase.ExchangeMaxCard(userid)
	var maxRet *niuniuCardRet
	var maxRetUserID int32

	fpoker := tbnn.poker.(*poker.FixedPoker)
	for _, ur := range tbnn.AllGameUser() {
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
	mlog.Debug("交换用户%d和用户%d牌型", userid, maxRetUserID)
	tbnn.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	fpoker.ExChange(userid, maxRetUserID)

}
