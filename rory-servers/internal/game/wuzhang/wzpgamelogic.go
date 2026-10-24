package wzp

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"fmt"
	"math"

	proto "example.invalid/golang/protobuf/proto"
)

const (
	actionName_GameTurn = "wzpGameTurn"
)

const (
	animationName_GameStart = "gamestart" //游戏开始
	animationName_FirstCard = "firstcard" //第一手发牌
	animationName_Card      = "card"      //发牌
	animationName_GameRet   = "gameret"   //亮牌

	animationTimeOut_GameStart = 1
	animationTimeOut_FirstCard = 2
	animationTimeOut_Card      = 3
	animationTimeout_GameRet   = 10
)

const (
	actionID_BaseBet   int32 = 1 //底注
	actionID_DoBet     int32 = 2 //下注
	actionID_Followbet int32 = 3 //跟注
	actionID_Giveup    int32 = 4 //弃牌
	actionID_AddBet    int32 = 5 //加注
	actionID_Suoha     int32 = 6 //梭哈

)

//五张牌游戏逻辑
type WZPGameLogic struct {
	*gamelogic.DynamicGameLogicBase
	poker                     poker.Poker           //游戏扑克牌
	gamePoker                 *poker.GamePoker      //玩家牌型
	gameBetLog                *gameBetLog           //投注记录
	turnManager               *gameTurnManager      //轮次管理
	currentActionName         string                //当前动作名称
	isShowHand                bool                  //是否梭哈
	gameret                   *netproto.WZP_GameRet //游戏结果
	currentWatchAnimationName string                //当前播放动画名称
}

//创建抢五张牌游戏逻辑对象
func NewWZPGameLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *WZPGameLogic {
	glc := new(WZPGameLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)

	//解析配置
	parseGameConfig(cfgs.GetConfig(gamelogic.ConfigFileNameGameServer))

	glc.RegisterAfterGameStartFunc(func() {

		//投入底注
		for _, ur := range glc.AllGameUser() {
			glc.gameBetLog.addBaseScoreBetLog(ur)

			//广播用户投入底注
			glc.boardcastUserGameAction(ur, actionID_BaseBet, int64(gameconfig.baseAmount))

			//牌局记录
			glc.GameLog.AddGameStep("用户%d投入底注%d", ur.UserID, gameconfig.baseAmount)
		}

		glc.turnManager.newTurnFunc = func(turnNo int32) {
			//最后一轮 结束发牌
			if glc.turnManager.isGameEnd() {
				return
			}

			//牌局记录
			glc.GameLog.AddGameStep("进入第%d轮", turnNo)

			//发下一张牌
			for _, ur := range glc.AllGameUser() {
				//如果用户已经弃牌
				if glc.turnManager.isGiveUp(ur.UserID) {
					continue
				}

				//发一张牌
				cards := glc.poker.Deal(ur.UserID, 1)

				upk := poker.NewUserPoker(ur, cards)

				glc.gamePoker.AddUserPoker(upk)

				//牌局记录
				glc.GameLog.Fapai(ur.UserID, cards.String())
			}

			//向用户发送牌型数据
			for _, ur := range glc.AllTableUser() {
				glc.sendUserGamePoker(ur)

			}

			//通知用户牌型
			glc.notifyAllUserCardType()

			mlog.Debug("发第%d轮牌完成.", glc.turnManager.turnNo)

		}

		//开始动画
		glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_GameStart, animationTimeOut_GameStart, func() {
			glc.fapai() //发牌

			//计算超时时间
			firstCardTimeout := animationTimeOut_FirstCard * glc.GetCurrentGameUserCount()
			glc.currentWatchAnimationName = glc.RegisterAnimation(animationName_FirstCard, int32(firstCardTimeout), func() {
				glc.turnManager.gameBegin() //游戏轮次开始
				glc.nextGameTurnAction()    //进入游戏轮次
			})
		})

		//		glc.AsyncCall(func() {
		//			time.Sleep(time.Duration(10) * time.Second)

		//		}, func() {
		//			print("async call back exec...")
		//		})

	})

	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData()
	})

	//	//解析protobuf 配置
	//	_, pcfgErr := glc.ReadProtoConfig("GOODS_INFO.bytes", func(data []byte) (proto.Message, error) {
	//		goodsinfo := &DataCfg.GOODS_INFO_ARRAY{}
	//		err := proto.Unmarshal(data, goodsinfo)

	//		return goodsinfo, err

	//	})
	//	if pcfgErr != nil {
	//		mlog.Error("解析proto配置错误,%v", pcfgErr)
	//	}

	return glc
}

//清除游戏数据
func (wzp *WZPGameLogic) clearGameData() {
	//初始化一副牌型

	//随机发牌
	if wzp.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		wzp.poker = poker.NewRandomPoker(false, false, 1)
		wzp.poker = poker.NewFixedPoker(wzp.AllTableUser(), wzp.poker.(*poker.RandomPoker), 5)
	} else {
		wzp.poker = poker.NewFilePoker(wzp.Configs, wzp.AllTableUser())
		wzp.poker = poker.NewFixedFilePoker(wzp.poker.(*poker.FilePoker))
	}

	wzp.gamePoker = poker.NewGamePoker() //创建游戏牌型
	wzp.gameBetLog = newGameBetLog(wzp)
	wzp.turnManager = newGameTurnManager(wzp.gamePoker, wzp)
	wzp.currentActionName = ""
	wzp.isShowHand = false
	wzp.gameret = nil
}

//发牌
func (wzp *WZPGameLogic) fapai() {
	//洗牌
	wzp.poker.Shuffle()

	for _, ur := range wzp.AllGameUser() {
		//发两张牌
		cards := wzp.poker.Deal(ur.UserID, 2)

		upk := poker.NewUserPoker(ur, cards)

		wzp.gamePoker.AddUserPoker(upk)

		//牌局记录
		wzp.GameLog.Fapai(ur.UserID, cards.String())
	}

	//向用户发送牌型数据
	for _, ur := range wzp.AllTableUser() {
		wzp.sendUserGamePoker(ur)
	}

	//通知用户牌型
	wzp.notifyAllUserCardType()

	mlog.Debug("发第一轮牌完成.")

}

//开始游戏轮次
func (wzp *WZPGameLogic) nextGameTurnAction() {

	//第五轮,游戏进行结算
	if wzp.turnManager.isGameEnd() {
		// wzp.GameLog.AddGameStep("游戏进行结算")
		wzp.gameret = newWZPGameRet(wzp).getGameRet()

		var desc = make(map[int32]string, 0)

		//写分
		for _, uret := range wzp.gameret.UserRets {
			user := wzp.GetUserByUserID(uret.GetUserID())
			cards := wzp.gamePoker.GetUserPoker(user.UserID).Cards
			cardType := getCardTypeString(cards)
			desc[user.UserID] = fmt.Sprintf("用户%d 结算金额, 牌型:%v, 手牌:%v ", user.UserID, cardType, cards)
			var flowAmount, betAmount, payAmount int64
			amount, tax := wzp.GetUserWinAndTax(uret.GetWinLose())
			flowAmount = util.ABSInt64(amount)
			betAmount = util.ABSInt64(amount) + tax
			payAmount = betAmount + amount
			uret.WinLose = proto.Int64(amount)
			wzp.UpdateUserScoreFunc(uret.GetUserID(), amount, tax)
			user.FlowAmount += flowAmount
			user.BetAmount += betAmount
			user.PayAmount += payAmount
			desc[user.UserID] += fmt.Sprintf("win %d, 税收 %d, 流水 %d, 有效投注 %d, 返還金額 %d", amount, tax, flowAmount, betAmount, payAmount)

			wzp.GameLog.Context().UserUpdateScore(user.UserID, amount)
			wzp.GameLog.Context().UserBetScore(user.UserID, betAmount)
			wzp.GameLog.Context().UserPayoutScore(user.UserID, payAmount)
			wzp.GameLog.Context().UserValidScore(user.UserID, flowAmount)
			wzp.GameLog.Context().UserTaxAmount(user.UserID, tax)
		}

		funcRet := func() {
			//广播游戏结果
			wzp.boardcastGameRet()

			wzp.currentWatchAnimationName = wzp.RegisterAnimation(animationName_GameRet, animationTimeout_GameRet, func() {
				wzp.GameEndFunc()
				//清理游戏数据
				wzp.clearGameData()
			})
		}

		if wzp.turnManager.isAllGiveUp() {
			funcRet()
		} else {
			wzp.TimerFuncManager.NewTimer(1, funcRet)
		}

		//牌局记录
		for _, v := range desc {
			//
			//if wzp.turnManager.isGiveUp(userID) {
			//	continue
			//}
			wzp.GameLog.AddGameStep(v)
		}

		uscoresRet := make(map[int32]int64)
		for _, uret := range wzp.gameret.UserRets {
			uscoresRet[uret.GetUserID()] = uret.GetWinLose()
		}
		wzp.GameLog.GameEnd()

		return
	}

	//当前轮次
	currentTurnNo := wzp.turnManager.turnNo
	//当前用户
	curentTurnUser := wzp.turnManager.currentTurnUser

	wzp.currentActionName = wzp.PlayerActionManager.NewAction(actionName_GameTurn, 1, gameconfig.betTimeOut+5, func() {
		if currentTurnNo == wzp.turnManager.turnNo && curentTurnUser == wzp.turnManager.currentTurnUser {
			//广播用户操作
			wzp.boardcastUserGameAction(wzp.turnManager.currentTurnUser, actionID_Giveup, 0)

			//设置当前用户弃牌
			wzp.turnManager.onUserAction(wzp.turnManager.currentTurnUser, true, actionID_Giveup)

			//开始下一轮
			wzp.nextGameTurnAction()
		}
	})

	//广播游戏轮次
	wzp.boardcastGameTrun()

	//牌局记录
	wzp.GameLog.AddGameStep("轮到用户%d操作,当前轮次%d", curentTurnUser.UserID, currentTurnNo)
}

//广播游戏结果
func (wzp *WZPGameLogic) boardcastGameRet() {
	gretmsg := new(msg.Message)
	gretmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gretmsg.SClassID = int32(netproto.WZP_GameMessageClassID_WZPGameRetID)

	gretmsg.MsgData = wzp.gameret

	wzp.BoardcastMessage(gretmsg)
}

//通知所有游戏用户牌型
func (wzp *WZPGameLogic) notifyAllUserCardType() {
	for _, ur := range wzp.AllGameUser() {
		wzp.notifyUserCardType(ur)
	}
}

//通知用户牌型
func (wzp *WZPGameLogic) notifyUserCardType(ur *gamelogic.GameUser) {
	ctmsg := new(msg.Message)
	ctmsg.BClassID = int32(netproto.MessageBClassID_Game)
	ctmsg.SClassID = int32(netproto.WZP_GameMessageClassID_WZPCardTypeID)

	ctdata := new(netproto.WZP_CardType)
	ctype := getCardSize(wzp.gamePoker.GetUserPoker(ur.UserID).Cards)
	ctdata.TypeID = proto.Int32(ctype)
	ctmsg.MsgData = ctdata

	wzp.SendMessageToUser(ur, ctmsg)

	mlog.Debug("发送用户%v牌型%d", ur, ctype)
}

func (wzp *WZPGameLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	//检查非游戏玩家
	if !wzp.CheckGameUserMessage(user) {
		wzp.KickUserFunc(user.UserID, "非法操作")
		return
	}

	//用户投注
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.WZP_GameMessageClassID_UserGameActionID) {
		wzp.onUserGameAction(user, gmsg.MsgData.([]byte))
	}

	//检查用户动画播放
	wzp.CheckUserAnimationCommited(user, gmsg, wzp.currentWatchAnimationName, nil)
}

//用户游戏动作
func (wzp *WZPGameLogic) onUserGameAction(user *gamelogic.GameUser, data []byte) {

	//检查是否用户操作
	if !wzp.turnManager.checkUserAction(user) {
		mlog.Debug("收到用户%v游戏动作,但当前不是该用户操作", user)
		return
	}

	//检查游戏是否已经结束
	if wzp.turnManager.isGameEnd() {
		mlog.Debug("收到用户%v游戏动作,但游戏已经结束", user)
		return
	}

	ust := &netproto.WZP_UserGameAction{}

	err := proto.Unmarshal(data, ust)
	if err != nil {
		mlog.Error("解析用户%v游戏动作错误,%v", user, err)
		wzp.KickUserFunc(user.UserID, "非法操作")
		return
	}
	*ust.ActionID = int32(math.Abs(float64(*ust.ActionID)))
	*ust.BetAmount = int64(math.Abs(float64(*ust.BetAmount)))

	mlog.Debug("收到用户%v游戏动作%v", user, ust)

	//是否弃牌
	giveup := false

	var usrTotalBet int64 = 0 // 玩家總下注
	for _, v := range wzp.gameBetLog.betlogs {
		if v.userID == user.UserID {
			usrTotalBet += v.amount
		}
	}

	if ust.GetActionID() == actionID_DoBet {
		if !wzp.gameBetLog.checkTurnBetAmount(ust.GetBetAmount(), wzp.turnManager.turnNo, user) {
			mlog.Error("用户%v,提交了非法的金额%d", user, ust.GetBetAmount())
			wzp.KickUserFunc(user.UserID, "非法操作")
			return
		}
	}

	if ust.GetActionID() == actionID_AddBet {
		if !wzp.gameBetLog.checkTurnBetAmount(ust.GetBetAmount(), wzp.turnManager.turnNo, user) {
			ust.ActionID = proto.Int32(actionID_Followbet)
		}
	}

	//跟注
	if ust.GetActionID() == actionID_Followbet {
		// 2020.02.11 玩家跟注的時候只補這回合目前最高下注的差額 by Alex

		var usrCurrentTurnTotalBet int64 = 0    // 當前回合玩家下注總量
		var usrCurrentHighestTotalBet int64 = 0 // 當前回合最高下注玩家的下注金額總量
		for _, v := range wzp.gameBetLog.betlogs {

			if v.userID == user.UserID && v.trun == wzp.turnManager.turnNo {
				usrCurrentTurnTotalBet += v.amount
				// fmt.Println("@@@", v.trun, v.userID, v.amount)
			}

			if v.userID == wzp.turnManager.firstUser.UserID && v.trun == wzp.turnManager.turnNo {
				usrCurrentHighestTotalBet += v.amount
			}
		}
		if wzp.isShowHand {
			amt, minuser := wzp.gameBetLog.getSuohaAmount()
			if minuser.UserID == user.UserID {
				ust.BetAmount = proto.Int64(amt)
			} else {
				ust.BetAmount = proto.Int64(usrCurrentHighestTotalBet)
			}
		} else {
			ust.BetAmount = proto.Int64(usrCurrentHighestTotalBet - usrCurrentTurnTotalBet)
		}
		// fmt.Println("跟注", wzp.turnManager.turnNo, user.UserID, wzp.turnManager.firstUser.UserID, usrCurrentHighestTotalBet, usrCurrentTurnTotalBet, ust.GetBetAmount())
	} else {
		// if ust.GetActionID() == actionID_AddBet {
		// 	fmt.Println("加注", wzp.turnManager.turnNo, user.UserID, wzp.turnManager.firstUser.UserID, ust.GetBetAmount())
		// } else if ust.GetActionID() == actionID_Suoha {
		// 	fmt.Println("梭哈", wzp.turnManager.turnNo, user.UserID, wzp.turnManager.firstUser.UserID, ust.GetBetAmount())
		// }
	}

	usrTotalBet += ust.GetBetAmount() // 用户预计总下注金额
	if ust.GetActionID() == actionID_Followbet || ust.GetActionID() == actionID_AddBet {
		if user.Score < usrTotalBet { // 钱不够下注设为弃牌
			ust.ActionID = proto.Int32(actionID_Giveup)
		}
	}

	//弃牌
	if ust.GetActionID() == actionID_Giveup {
		giveup = true

		//牌局记录
		wzp.GameLog.AddGameStep("用户%d弃牌", user.UserID)
	}

	//梭哈
	if ust.GetActionID() == actionID_Suoha {

		amt, minuser := wzp.gameBetLog.getSuohaAmount()
		var usrCurrentTurnTotalBet int64 = 0    // 當前回合玩家下注總量
		var usrCurrentHighestTotalBet int64 = 0 // 當前回合最高下注玩家的下注金額總量
		for _, v := range wzp.gameBetLog.betlogs {

			if v.userID == user.UserID && v.trun == wzp.turnManager.turnNo {
				usrCurrentTurnTotalBet += v.amount
			}

			if v.userID == wzp.turnManager.firstUser.UserID && v.trun == wzp.turnManager.turnNo {
				usrCurrentHighestTotalBet += v.amount
			}
		}

		if minuser.UserID == user.UserID {
			ust.BetAmount = proto.Int64(amt)
		} else {
			ust.BetAmount = proto.Int64(amt + usrCurrentHighestTotalBet)
		}

		// fmt.Println("梭哈", user.UserID, amt, usrCurrentHighestTotalBet)

		wzp.isShowHand = true

	}

	// fmt.Println("####", user.UserID, ust.GetActionID(), ust.GetBetAmount())

	//下注,跟注,加注
	if ust.GetActionID() == actionID_DoBet || ust.GetActionID() == actionID_Followbet || ust.GetActionID() == actionID_AddBet || ust.GetActionID() == actionID_Suoha {

		// 2020.02.11 為了能跟注籌碼以外的金額，所以先註解這裡 by Alex
		//检查本轮下注金额
		// if !wzp.gameBetLog.checkTurnBetAmount(ust.GetBetAmount(), wzp.turnManager.turnNo, user) {
		// 	mlog.Error("用户%v,提交了非法的金额%d", user, ust.GetBetAmount())
		// 	wzp.KickUserFunc(user.UserID, "非法操作")
		// 	return
		// }

		mlog.Debug("用户%v投注%d", user, ust.GetBetAmount())

		//添加下注记录
		wzp.gameBetLog.addBetLog(user, ust.GetBetAmount(), wzp.turnManager.turnNo)

		//牌局记录
		stepdesc := "下注"
		if ust.GetActionID() == actionID_Followbet {
			stepdesc = "跟注"
		} else if ust.GetActionID() == actionID_AddBet {
			stepdesc = "加注"
		} else if ust.GetActionID() == actionID_Suoha {
			stepdesc = "梭哈"
		}
		wzp.GameLog.AddGameStep("用户%d%s,金额:%d,轮次:%d", user.UserID, stepdesc, ust.GetBetAmount(), wzp.turnManager.turnNo)

	}

	//广播用户操作
	wzp.boardcastUserGameAction(user, ust.GetActionID(), ust.GetBetAmount())

	//通知用户操作
	wzp.turnManager.onUserAction(user, giveup, ust.GetActionID())

	//撤销本轮动作监听
	wzp.PlayerActionManager.CancelAction(wzp.currentActionName)

	//下一轮
	wzp.nextGameTurnAction()

}

//向用户发送游戏牌型
func (wzp *WZPGameLogic) sendUserGamePoker(user *gamelogic.GameUser) {

	for _, ur := range wzp.AllGameUser() {
		//如果用户已经弃牌
		if wzp.turnManager.isGiveUp(ur.UserID) {
			continue
		}

		gpmsg := new(msg.Message)
		gpmsg.BClassID = int32(netproto.MessageBClassID_Game)
		gpmsg.SClassID = int32(netproto.GameMessageClassID_UserPokerID)
		gpmsgData := wzp.getUserPokerMsgData(ur, user)
		gpmsg.MsgData = gpmsgData

		wzp.SendMessageToUser(user, gpmsg)

		mlog.Debug("向用户%v发送牌型数据,%v", user, gpmsg.MsgData)
	}
}

//获取用户牌型消息
func (wzp *WZPGameLogic) getUserPokerMsgData(ur *gamelogic.GameUser, recvUser *gamelogic.GameUser) *netproto.UserPoker {
	upk := wzp.gamePoker.GetUserPoker(ur.UserID)

	//向自己发牌,发送暗牌
	if ur.UserID == recvUser.UserID {

	} else {
		//向他人发牌
		cards := poker.NewCardGroup(nil)
		if upk != nil && upk.Cards != nil && upk.Cards.Cards != nil {
			for i := 0; i < len(upk.Cards.Cards); i++ {
				//不发暗牌
				if i == 0 {
					cards.AddCard(poker.NewEmptyPokerCard())
				} else {
					cards.AddCard(upk.Cards.Cards[i])
				}
			}
		}
		upk = poker.NewUserPoker(ur, cards)
	}

	return upk.GetNetMessageData()
}

func (wzp *WZPGameLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {

	wzp.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	wzp.sendGameConfig(gameuser)

	//如果非游戏玩家，发送当前场景
	if wzp.IsGaming() && !wzp.IsGameUser(gameuser) {
		// 发送游戏场景
		wzp.sendGameScene(gameuser)
	}
}

//用户断线重连
func (wzp *WZPGameLogic) UserRecover(gameuser *gamelogic.GameUser) {

	wzp.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	wzp.sendGameConfig(gameuser)

	//发送游戏场景
	if wzp.IsGaming() {
		// 发送游戏场景
		wzp.sendGameScene(gameuser)
	}
}

//发送游戏配置
func (wzp *WZPGameLogic) sendGameConfig(user *gamelogic.GameUser) {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.WZP_GameMessageClassID_GameConfigID)

	gcdata := new(netproto.WZP_GameConfig)
	gcdata.BaseAmount = proto.Int32(int32(gameconfig.baseAmount))
	gcdata.BetTimeout = proto.Int32(int32(gameconfig.betTimeOut))
	gcdata.ReadyTimeout = proto.Int32(int32(wzp.GetServerConfig().ReadyWaitTime))

	for _, ba := range gameconfig.betAmount {
		gcdata.BetAmount = append(gcdata.BetAmount, int32(ba))
	}
	gmsg.MsgData = gcdata

	wzp.SendMessageToUser(user, gmsg)
}

//广播当前轮次
func (wzp *WZPGameLogic) boardcastGameTrun() {
	gmsg := wzp.getGameTurnMsg()
	wzp.BoardcastMessage(gmsg)

	mlog.Debug("广播用户轮次%v", gmsg.MsgData)
}

//获取游戏轮次消息
func (wzp *WZPGameLogic) getGameTurnMsg() *msg.Message {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.WZP_GameMessageClassID_UserGameTurnID)

	gmsg.MsgData = wzp.getGameTurnMsgData()

	return gmsg
}

func (wzp *WZPGameLogic) getGameTurnMsgData() *netproto.WZP_NotifyUserGameTrun {
	cur := wzp.turnManager.currentTurnUser

	var isFirst int32
	var isLast int32
	var isSuoha int32

	isFirst = 0
	isLast = 0
	isSuoha = 0

	if wzp.turnManager.isFirstTurn() {
		isFirst = 1
	}

	// 2020.02.11 為了讓最這回合最後一個下注用戶可以做加注動作，所以注解這裡 by Alex
	// if wzp.turnManager.isLastTurn() {
	// 	isLast = 1
	// }

	if wzp.isShowHand {
		isSuoha = 1
	}

	timeout := gameconfig.betTimeOut
	hasExTime := wzp.PlayerActionManager.GetActionTimerElapsedSeconds(wzp.currentActionName)
	if hasExTime > 0 {
		hasExTime += 1
	}

	timeout = timeout - hasExTime

	gtdata := new(netproto.WZP_NotifyUserGameTrun)
	gtdata.UserID = proto.Int32(cur.UserID)
	gtdata.TrunNo = proto.Int32(wzp.turnManager.turnNo)
	gtdata.Timeout = proto.Int32(timeout)
	gtdata.IsFirst = proto.Int32(isFirst)
	gtdata.IsLast = proto.Int32(isLast)
	gtdata.IsSuoha = proto.Int32(isSuoha)
	//第一轮
	if isFirst == 1 {
		gtdata.TurnBetAmount = proto.Int64(0)
	} else {
		gtdata.TurnBetAmount = proto.Int64(wzp.gameBetLog.currentTurnBetAmount)
	}

	return gtdata
}

//广播用户游戏操作
func (wzp *WZPGameLogic) boardcastUserGameAction(user *gamelogic.GameUser, actionID int32, betamount int64) {
	gmsg := new(msg.Message)
	gmsg.BClassID = int32(netproto.MessageBClassID_Game)
	gmsg.SClassID = int32(netproto.WZP_GameMessageClassID_NotifyUserGameActionID)

	//是否梭哈
	var issuoha int32 = 0
	if wzp.isShowHand {
		issuoha = 1
	}

	var userBetAmount int64 = 0
	if ubtamt, ok := wzp.gameBetLog.userBets[user.UserID]; ok {
		userBetAmount = ubtamt
	}

	actdata := new(netproto.WZP_NotifyUserGameAction)
	actdata.ActionID = proto.Int32(actionID)
	actdata.BetAmount = proto.Int64(betamount)
	actdata.UserID = proto.Int32(user.UserID)
	actdata.TotalBetAmount = proto.Int64(wzp.gameBetLog.totalAmount)
	actdata.IsSuoha = proto.Int32(issuoha)
	actdata.GameUser = user.GetProtoBufGameUser()
	actdata.UserBetAmount = proto.Int64(userBetAmount)

	gmsg.MsgData = actdata

	wzp.BoardcastMessage(gmsg)

	mlog.Debug("广播用户动作%v", actdata)
}

//发送游戏场景
func (wzp *WZPGameLogic) sendGameScene(ur *gamelogic.GameUser) {
	sceneMsg := new(msg.Message)
	sceneMsg.BClassID = int32(netproto.MessageBClassID_Game)
	sceneMsg.SClassID = int32(netproto.WZP_GameMessageClassID_WZPGameSceneID)

	sceneData := new(netproto.WZP_GameScene)
	//牌局结果
	sceneData.GameRet = wzp.gameret
	//当前轮次
	sceneData.TurnNo = proto.Int32(wzp.turnManager.turnNo)

	//弃牌用户
	for _, uid := range wzp.turnManager.giveupUsers {
		sceneData.GiveUpUsers = append(sceneData.GiveUpUsers, uid)
	}

	//总投注金额
	sceneData.TotalBetAmount = proto.Int64(wzp.gameBetLog.totalAmount)

	//用户牌型
	for _, gur := range wzp.AllGameUser() {
		upk := wzp.getUserPokerMsgData(gur, ur)

		sceneData.UserPoker = append(sceneData.UserPoker, upk)
	}

	//当前轮次
	if sceneData.GameRet == nil {
		if wzp.turnManager.currentTurnUser != nil {
			sceneData.GameTurnUser = wzp.getGameTurnMsgData()
		}
	}

	//投注筹码
	for _, btlog := range wzp.gameBetLog.betlogs {
		sceneData.BetLogs = append(sceneData.BetLogs, btlog.amount)
	}

	sceneMsg.MsgData = sceneData

	//发送场景数据包
	wzp.SendMessageToUser(ur, sceneMsg)

	//发送用户牌型
	if ur.IsGameing() {
		wzp.notifyUserCardType(ur)
	}

}

//换最大牌
func (wzp *WZPGameLogic) ExchangeMaxCard(userid int32) {
	wzp.LogicBase.ExchangeMaxCard(userid)

	fpoker := wzp.poker.(*poker.FixedPoker)

	//计算最大牌型
	var maxcards *poker.CardGroup
	var maxCardsUserID int32

	for _, ur := range wzp.AllGameUser() {
		ucards := fpoker.GetUserCards(ur.UserID)
		if maxcards == nil {
			maxcards = ucards
			maxCardsUserID = ur.UserID
		} else {
			ret, _, _ := compareCard(maxcards, ucards)

			if !ret {
				maxcards = ucards
				maxCardsUserID = ur.UserID
			}
		}
	}

	if maxCardsUserID != userid {
		mlog.Debug("交换用户%d和用户%d牌型", userid, maxCardsUserID)
		wzp.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxCardsUserID)
		fpoker.ExChange(userid, maxCardsUserID)
	}
}
