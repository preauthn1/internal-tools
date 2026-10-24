//==================================================================
//文件名： ttzlogic.go
//说明：   推筒子战斗逻辑
//修改记录：
//==================================================================
package ttz

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

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

// 推筒子逻辑数据结构
type TtzLogic struct {
	*gamelogic.DynamicGameLogicBase
	timerTimeOutAction *event.TimerFunc
	ttzCfgs            *ttzCfg
	strAinimationName  string //动画名

	nTimeGameEnd        int32 //结束游戏超时时间
	nTimeDealOne        int32 //发牌结束时间
	nTimeDealTwo        int32
	nRobTimeOut         int32  //抢庄超时时间
	nBetTimeOut         int32  //押注超时时间
	nLightTimeOut       int32  //亮牌超时时间
	nBeginGameTimeOut   int32  //开始游戏
	strAniNameDealOne   string //发牌动画
	strAniNameDealTwo   string
	strAniNameGameEnd   string //结束游戏动画
	strAniNameBeginGame string //开始游戏

	nGameType  int32 //游戏类型，新手局等
	nGameCount int32 //游戏人数

	nNumberDice      int32                  //骰子点数
	nGameState       int32                  //游戏阶段
	pokerGroup       poker.Poker            //公共牌池
	nButtonSeatNo    int32                  //庄家座位号
	nButtonValue     int32                  //抢庄值
	mapGameUserData  map[int32]*TtzUserData //玩家数据
	stSettlementInfo []*settleMentInfo      //玩家小局结算

	nCurActionSeatNo  int32 //当前操作玩家
	nLastActionSeatNo int32 //上次操作玩家
	nDealIndex        int32 //发牌序号
	nAutoBeginTime    int32 //自动操作开始时间，超时
	nGameIndex        int32 //游戏进行局数
	nLastDealSeat     int32 //发牌开始计数点
	nLastDealId       int32 //上次发牌userid
	nBaseScore        int64 //底分
	nPokerType        int32 //发牌类型
	mLock             *sync.Mutex
}

//清除数据，初始化状态
func (ttz *TtzLogic) clearGameData() {
	ttz.pokerGroup = MakeRandomCards() //初始化一副牌
	//
	//随机发牌
	if ttz.nPokerType == 0 {
		ttz.pokerGroup = MakeRandomCards()
		ttz.pokerGroup = poker.NewFixedPoker(ttz.AllTableUser(), ttz.pokerGroup.(*poker.RandomPoker), 2)
	} else { //发牌器发牌
		ttz.pokerGroup = poker.NewFilePoker(ttz.Configs, ttz.AllTableUser())
		ttz.pokerGroup = poker.NewFixedFilePoker(ttz.pokerGroup.(*poker.FilePoker))
	}

	mlog.Debug("初始化数据。")
	ttz.setGameUserData()
	ttz.nGameIndex++
	ttz.nNumberDice = 0
	ttz.nGameState = GAMEPROGRESS_DICE
	ttz.nButtonValue = 1
	ttz.nCurActionSeatNo = 0
	ttz.nLastActionSeatNo = 0
	ttz.timerTimeOutAction = nil
	ttz.strAinimationName = ""
	ttz.nDealIndex = 0
	ttz.nAutoBeginTime = int32(time.Now().Unix())
	ttz.mLock = new(sync.Mutex)
	ttz.updateDealSeatInit()

}

// 创建推筒子逻辑
func NewTtzLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *TtzLogic {
	//mlog.Debug("创建推筒子逻辑")
	ttz := new(TtzLogic)
	ttz.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)
	ttz.getGameConfig(cfgs)
	ttz.nGameIndex = 0
	ttz.nButtonSeatNo = 0
	ttz.RegisterAfterGameStartFunc(func() {
		ttz.GameLog.AddGameStep("开始游戏,桌子号%d,第%d局！", ttz.TableNo(), ttz.nGameIndex)
		mlog.Debug("开始游戏，桌子号%d,第%d局。", ttz.TableNo(), ttz.nGameIndex)
		ttz.strAinimationName = ttz.RegisterAnimation(ttz.strAniNameBeginGame, ttz.nBeginGameTimeOut, ttz.beginGame)
	})
	ttz.RegisterBeforeGameStartFunc(func() {
		//清除数据
		ttz.clearGameData() //清除数据
	})
	return ttz
}

// 推筒子消息管理
func (ttz *TtzLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	mlog.Debug("GameMessage收到消息，座位号：%d", user.SeatNo)
	ttz.CheckUserAnimationCommited(user, gmsg, ttz.strAinimationName, nil)
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		mlog.Debug("GameMessage,animation commite SeatNo=%d,Name=%s", user.SeatNo, ttz.strAinimationName)
		return
	}
	//操作合法性校验
	if !ttz.oprateCheckLegal(user, gmsg) {
		return
	}
	//战斗逻辑操作
	ttz.oprateBattleLogic(user, gmsg)
}

//获取当前游戏配置
func (ttz *TtzLogic) getGameConfig(cfgs *config.Group) {
	//mlog.Debug("获取配置文件数据")
	cfg := NewTtzCfg(cfgs)
	ttz.ttzCfgs = cfg
	ttz.nRobTimeOut = int32(cfg.GetValueInt("Rob"))
	ttz.nBetTimeOut = int32(cfg.GetValueInt("Bet"))
	ttz.nLightTimeOut = int32(cfg.GetValueInt("Light"))
	ttz.nBeginGameTimeOut = int32(cfg.GetValueInt("AniTimeBeginGame"))
	ttz.nTimeGameEnd = int32(cfg.GetValueInt("AniTimeEndGame"))
	ttz.nTimeDealOne = int32(cfg.GetValueInt("AniTimeDealOne"))
	ttz.nTimeDealTwo = int32(cfg.GetValueInt("AniTimeDealTwo"))
	ttz.nBaseScore = int64(cfg.GetValueInt("BaseScore"))
	ttz.nPokerType = int32(cfg.GetValueInt("PokerType"))

	ttz.strAniNameDealOne = cfg.GetValueStr("AniNameDealOne")
	ttz.strAniNameDealTwo = cfg.GetValueStr("AniNameDealTwo")
	ttz.strAniNameGameEnd = cfg.GetValueStr("AniNameGameEnd")
	ttz.strAniNameBeginGame = cfg.GetValueStr("AniNameBeginGame")

}

//开始游戏
func (ttz *TtzLogic) beginGame() {
	ttz.playDice()         //骰子
	ttz.calMaxUserCanRob() //抢庄倍率
	ttz.TtzDeal()          //发牌
}

//发牌，每人发1张底牌
func (ttz *TtzLogic) TtzDeal() {
	if ttz.nGameState != GAMEPROGRESS_DEALONE {
		mlog.Error("游戏阶段错误，发牌失败")
		return
	}

	//给所有用户发牌
	ttz.dealCardsOne(1, ttz.strAniNameDealOne, ttz.nTimeDealOne, ttz.autoDealOne)
	mlog.Debug("发牌结束")
	//ttz.nGameState = GAMEPROGRESS_ROBBUTTON
}

//第一次发牌
func (ttz *TtzLogic) dealCardsOne(number int32, strAniName string, timeOut int32, callback func()) {
	ttz.nDealIndex = 1
	var count int32 = 0
	nSeatNo := ttz.nButtonSeatNo - 1 + ((ttz.nNumberDice + ttz.nButtonSeatNo) % ttz.nGameCount) //丢骰子决定的顺序
	for i := 0; i < 10; i++ {
		nSeatNo++
		if nSeatNo > 3 {
			nSeatNo = 0
		}
		ud := ttz.mapGameUserData[nSeatNo]
		if ud == nil {
			//mlog.Error("用户数据错误，seatNo=%d", nSeatNo)
			continue
		}

		if ud.userCards.Len() >= 1 {
			count++
			if count >= ttz.nGameCount {
				break
			}
			continue
		}

		cards := ttz.pokerGroup.Deal(ud.User.UserID, int(number))
		if cards == nil {
			ttz.endGame()
		}
		ttz.setUserCards(ud.User, cards)
		//广播出去
		sendDealMsg(ttz, nSeatNo, cards)
		count++
		ttz.GameLog.AddGameStep("第%d次发牌,用户ID[%d],座位号[%d],牌[%s]", ttz.nDealIndex, ud.User.UserID, ud.User.SeatNo, StringCardGroup(cards))
		mlog.Debug("开始发牌[%d],SeatNo=%d", count, nSeatNo)
		if count >= ttz.nGameCount {
			break
		}
	}
	//ttz.updateCardType()
	ttz.sendDealWatch()
	mlog.Debug("dealCards, animation name=%s,timeout=%d", strAniName, timeOut)
	ttz.strAinimationName = ttz.RegisterAnimation(strAniName, timeOut, callback)
}

//第二次发牌
func (ttz *TtzLogic) dealCardsTwo(number int32, strAniName string, timeOut int32, callback func()) {
	ttz.nDealIndex = 2
	var count int32 = 0
	nSeatNo := ttz.nButtonSeatNo - 1 + ((ttz.nNumberDice + ttz.nButtonSeatNo) % ttz.nGameCount) //丢骰子决定的顺序
	for i := 0; i < 10; i++ {
		nSeatNo++
		if nSeatNo > 3 {
			nSeatNo = 0
		}
		ud := ttz.mapGameUserData[nSeatNo]
		if ud == nil {
			//mlog.Error("用户数据错误，seatNo=%d", nSeatNo)
			continue
		}
		if ud.userCards.Len() >= 2 {
			count++
			if count >= ttz.nGameCount {
				break
			}
			continue
		}

		cards := ttz.pokerGroup.Deal(ud.User.UserID, int(number))
		if cards == nil {
			ttz.endGame()
		}
		ttz.setUserCards(ud.User, cards)
		//广播出去
		sendDealMsg(ttz, nSeatNo, cards)
		count++
		ttz.GameLog.AddGameStep("第%d次发牌,用户ID[%d],座位号[%d],牌[%s]", ttz.nDealIndex, ud.User.UserID, ud.User.SeatNo, StringCardGroup(cards))
		mlog.Debug("开始发牌[%d],SeatNo=%d", count, nSeatNo)
		if count >= ttz.nGameCount {
			break
		}
	}
	ttz.updateCardType()
	ttz.sendDealWatch()
	mlog.Debug("dealCards, animation name=%s,timeout=%d", strAniName, timeOut)
	ttz.strAinimationName = ttz.RegisterAnimation(strAniName, timeOut, callback)
}

//发送观战者看信息
func (ttz *TtzLogic) sendDealWatch() {
	for _, user := range ttz.AllTableUser() {
		if ttz.mapGameUserData[user.SeatNo] != nil {
			continue
		}
		sendDealMsgWatch(ttz, user.SeatNo, 1)
	}
}

//第一次自动发牌
func (ttz *TtzLogic) autoDealOne() {
	mlog.Debug("第一次发牌动画播放完成，流程继续")
	ttz.nGameState = GAMEPROGRESS_ROBBUTTON
	BroadCastGameState(ttz, ttz.getOutTime())
	ttz.waitOpt() //第一个人操作
}

//更新牌型信息
func (ttz *TtzLogic) updateCardType() {
	for _, ud := range ttz.mapGameUserData {
		ud.nCardType, ud.nCardCount = getCardTypeAndCount(ud.userCards)
		ttz.GameLog.AddGameStep("玩家牌型信息：用户ID[%d],座位号[%d],手牌[%s],牌型[%s]", ud.User.UserID, ud.User.SeatNo, StringCardGroup(ud.userCards), strCardType(ud.nCardType, ud.nCardCount))
	}
}

//设置指定牌给指定玩家
func (ttz *TtzLogic) setUserCards(ur *gamelogic.GameUser, cards *poker.CardGroup) {
	for seatNo, ud := range ttz.mapGameUserData {
		if seatNo == ur.SeatNo {
			ud.userCards.Cards = append(ud.userCards.Cards, cards.Cards...)
		}
	}
}

//获取超时时间
func (ttz *TtzLogic) getOutTime() int32 {
	if ttz.nGameState == GAMEPROGRESS_DEALONE {
		return 0
	}
	if ttz.nGameState == GAMEPROGRESS_ROBBUTTON {
		return ttz.nRobTimeOut
	}
	if ttz.nGameState == GAMEPROGRESS_CHOOSEBET {
		return ttz.nBetTimeOut
	}
	if ttz.nGameState == GAMEPROGRESS_LIGHT {
		return ttz.nLightTimeOut
	}
	return 0
}

//发送初始化操作信息,获取下次操作数据
func (ttz *TtzLogic) waitOpt() {
	mlog.Debug("waitOpt...")
	if ttz.isEndCurState() {
		ttz.resetActionState()
		switch ttz.nGameState {
		case GAMEPROGRESS_ROBBUTTON: //抢庄结束,op 然后进入押注
			mlog.Debug("抢庄结束")
			ttz.setButtonInfo()      //设定庄家信息
			BroadCastButtonInfo(ttz) //广播庄家信息
			//注册动画
			ttz.nGameState = GAMEPROGRESS_ROBANIMATION
			ttz.strAinimationName = ttz.RegisterAnimation("AniNameRob", 2, ttz.autoRotAnimation)
			return
		case GAMEPROGRESS_ROBANIMATION:
			return
		case GAMEPROGRESS_CHOOSEBET: //选择倍数结束，op 然后进入亮牌
			mlog.Debug("押注结束")
			ttz.nGameState = GAMEPROGRESS_LIGHT
			ttz.dealCardsTwo(1, ttz.strAniNameDealTwo, ttz.nTimeDealTwo, ttz.autoDealTwo)
		case GAMEPROGRESS_LIGHT: //亮牌结束，op 然后进入结算
			mlog.Debug("亮牌结束")
			ttz.nGameState = GAMEPROGRESS_SETTLE
			ttz.gameSettlement() //结算

			//牌局记录
			uscoresRet := make(map[int32]int64)
			for _, ud := range ttz.mapGameUserData {
				uscoresRet[ud.User.UserID] = ud.nCurRoundWinMoney
			}
			ttz.GameLog.GameEnd()

			BroadCastTotalSettleRound(ttz) //广播结算信息
			ttz.strAinimationName = ttz.RegisterAnimation(ttz.strAniNameGameEnd, ttz.nTimeGameEnd, ttz.endGame)
			return
		default:
			mlog.Error("未知游戏阶段[%d]", ttz.nGameState)
		}
	}
	if ttz.nGameState == GAMEPROGRESS_ROBANIMATION {
		mlog.Debug("抢庄动画播放过程中，不做处理")
		return
	}
	ttz.timeOutAction()
}

//自动抢庄
func (ttz *TtzLogic) autoRotAnimation() {
	ttz.nGameState = GAMEPROGRESS_CHOOSEBET
	ttz.calMaxUserCanBet()
	BroadCastGameState(ttz, ttz.getOutTime())
	ttz.timeOutAction()
}

//自动发第二张牌
func (ttz *TtzLogic) autoDealTwo() {
	mlog.Debug("第二次发牌动画播放完成")
	//ttz.updateCardType()
	BroadCastUpdateCT(ttz) //广播牌型信息
	BroadCastGameState(ttz, ttz.getOutTime())
}

//结束游戏
func (ttz *TtzLogic) endGame() {
	ttz.GameEndFunc()
}

//设置超时自动操作
func (ttz *TtzLogic) timeOutAction() {
	if ttz.timerTimeOutAction == nil {
		ttz.nAutoBeginTime = int32(time.Now().Unix())
		str := fmt.Sprintf("autoTimer%d_%d", time.Now().UnixNano(), ttz.TableNo())
		ttz.timerTimeOutAction = ttz.TimerFuncManager.NewNamedTimer(ttz.getOutTime(), ttz.autoAction, str)
	}
}

//超时自动操作
func (ttz *TtzLogic) autoAction() {
	mlog.Debug("autoAction...")
	ttz.autoActionInfo()
	ttz.waitOpt()
}

//超时自动操作动作信息
func (ttz *TtzLogic) autoActionInfo() {
	for nSeatNo, ud := range ttz.mapGameUserData {
		switch ttz.nGameState {
		case GAMEPROGRESS_ROBBUTTON:
			if ud.getRob() {
				ud.setAction(true)
				continue
			}
			ud.setRob(true)
			ud.setAction(true)
			ud.nRateRob = 0
			BroadCastActionMsg(ttz, nSeatNo, ACTION_ROB, 0) //广播，不抢
			ttz.GameLog.AddGameStep("玩家操作[抢庄]：用户ID[%d],座位号[%d],超时自动操作[不抢]", ud.User.UserID, ud.User.SeatNo)
		case GAMEPROGRESS_CHOOSEBET:
			if ud.getBet() {
				ud.setAction(true)
				continue
			}
			if ud.User.SeatNo == ttz.nButtonSeatNo {
				continue
			}
			ud.setAction(true)
			ud.setBet(true)
			ud.nRateBet = 2
			BroadCastActionMsg(ttz, nSeatNo, ACTION_BET, 2) //默认，//两倍
			ttz.GameLog.AddGameStep("玩家操作[下注]：用户ID[%d],座位号[%d],超时自动操作[下注%d倍]", ud.User.UserID, ud.User.SeatNo, 2)
		case GAMEPROGRESS_LIGHT:
			if ud.getLight() {
				ud.setAction(true)
				continue
			}
			ud.setLight(true)
			ud.setAction(true)
			BroadCastLight(ttz, nSeatNo) //亮牌
			ttz.GameLog.AddGameStep("玩家操作[亮牌]：用户ID[%d],座位号[%d],超时自动操作[亮牌]", ud.User.UserID, ud.User.SeatNo)
		}
	}
}

//移除自动操作定时器
func (ttz *TtzLogic) removeTimeoutAction() {
	//移除操作超时Timer
	if ttz.timerTimeOutAction != nil {
		ttz.timerTimeOutAction.Stop()
		ttz.timerTimeOutAction = nil
	}
}

//玩家操作是否合法
func (ttz *TtzLogic) isValid(user *TtzUserData) bool {

	return true
}

//根据座位号获取用户信息
func (ttz *TtzLogic) getUserInfoBySeatNo(seatno int32) *TtzUserData {
	for SeatNo, urd := range ttz.mapGameUserData {
		if SeatNo == seatno {
			return urd
		}
	}
	return nil
}

//操作合法性校验
func (ttz *TtzLogic) oprateCheckLegal(user *gamelogic.GameUser, gmsg *msg.Message) bool {

	if !ttz.CheckGameUserMessage(user) {
		ttz.KickUserFunc(user.UserID, "非法操作")
		return false
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return false
	}
	//没有找到相关玩家数据
	urd := ttz.getUserInfoBySeatNo(user.SeatNo)
	if urd == nil {
		mlog.Error("用户未找到%v", user)
		return false
	}

	return true
}

//战斗逻辑操作
func (ttz *TtzLogic) oprateBattleLogic(user *gamelogic.GameUser, gmsg *msg.Message) {

	var nRet int32 = 0
	var strErrMessage string = ""
	switch gmsg.SClassID {

	case int32(netproto.TTZ_GameMessageClassID_ActionOprateID):
		//解析数据
		ust := &netproto.TTZ_ActionInfo{}
		err := proto.Unmarshal(gmsg.MsgData.([]byte), ust)
		//错误校验
		if err != nil {
			mlog.Error("解析用户%v用户错误%v", user, err)
			return
		}
		*ust.ActionType = int32(math.Abs(float64(*ust.ActionType)))
		*ust.ActionValue = int32(math.Abs(float64(*ust.ActionValue)))

		nRet = ttz.actionInfo(user, *(ust.ActionType), *(ust.ActionValue))
	//加个GM

	default:
		mlog.Error("unknow class id,gmsg.SClassID=%d", gmsg.SClassID)
		nRet = ERROR_CLASSID
		strErrMessage = "Unknow class id."
	}
	if nRet != 0 {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[TtzLogic]战斗逻辑操作 失敗, userid[%v] SeatNo[%v] 原因 [%v] [%v]", user.UserID, user.SeatNo, nRet, strErrMessage)).Error("HallServer壓測時")
		mlog.Error("oprateBattleLogic failed,SeatNo=%d.", user.SeatNo)
		BroadCastErrMsg(ttz, user.SeatNo, nRet, strErrMessage)
		//return
	}
	ttz.waitOpt()

}

//进入房间
func (ttz *TtzLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	mlog.Debug("UserEnter,SeatNo=%d enter room", seatno)
	ttz.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//发送游戏设置
	sendGameSettings(ttz, gameuser)
	//如果非游戏玩家，发送当前场景
	if ttz.IsGaming() && !ttz.IsGameUser(gameuser) {
		sendGameScene(ttz, gameuser)
	}
}

//用户断线重回
func (ttz *TtzLogic) UserRecover(gameuser *gamelogic.GameUser) {

	ttz.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	sendGameSettings(ttz, gameuser)

	//发送游戏场景
	if ttz.IsGaming() {
		sendGameScene(ttz, gameuser)
	}
}

//设置用户数据
func (ttz *TtzLogic) setGameUserData() {

	ttz.mapGameUserData = make(map[int32]*TtzUserData, 0)
	ttz.nGameCount = 0
	for _, ur := range ttz.AllTableUser() {
		userdata := NewTtzUserData(ur)
		if ur.SeatNo < 0 {
			mlog.Error("userid[%d],seatNo[%d] error!!!", ur.UserID, ur.SeatNo)
			ttz.endGame()
		}
		ttz.mapGameUserData[ur.SeatNo] = userdata
		ttz.nGameCount++
	}
	mlog.Debug("设置玩家数据，玩家人数：%d人。", ttz.nGameCount)
}

//骰子
func (ttz *TtzLogic) playDice() int32 {
	if ttz.nGameState != GAMEPROGRESS_DICE {
		mlog.Error("playDice,Game state is error.")
		return -1
	}
	nDice := ttz.randOneNumber(6)
	ttz.nNumberDice = int32(nDice)
	//ttz.GameLog.AddGameStep("丢骰子，点数[%d]", ttz.nNumberDice)
	ttz.nGameState = GAMEPROGRESS_DEALONE
	return 0
}

//玩家操作处理逻辑
func (ttz *TtzLogic) actionInfo(user *gamelogic.GameUser, nActionType, nActionValue int32) int32 {
	mlog.Debug("actionInfo,SeatNo=%d,nActionType=%d.nActionValue=%d", user.SeatNo, nActionType, nActionValue)
	var nRet int32 = 0
	switch nActionType {
	case ACTION_ROB:
		nRet = ttz.actionRob(user.SeatNo, nActionType, nActionValue)
	case ACTION_BET:
		nRet = ttz.actionBet(user.SeatNo, nActionType, nActionValue)
	case ACTION_LIGHT:
		nRet = ttz.actionLight(user.SeatNo)
	default:
		mlog.Error("actionInfo,error ActionType[%d]", nActionType)
		nRet = ERROR_ACTIONTYPE
	}
	return nRet
}

//抢庄
func (ttz *TtzLogic) actionRob(nSeatNo, nActionType, nActionValue int32) int32 {
	//mlog.Debug("actionRob,nSeatNo=%d,nActionType=%d.nActionValue=%d", nSeatNo, nActionType, nActionValue)
	//逻辑操作
	userInfo := ttz.getUserInfoBySeatNo(nSeatNo)
	if userInfo == nil {
		mlog.Error("userinfo is not find by seatno=%d", nSeatNo)
		return ERROR_USERNOTEXITST
	}
	if userInfo.getRob() {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[TtzLogic]機器人[%v] 重复操作抢庄,seat[%d]", userInfo.User.UserID, nSeatNo)).Error("HallServer壓測時")
		mlog.Error("重复操作抢庄,seat[%d]", nSeatNo)
		return ERROR_REPEATOP
	}
	if userInfo.nMaxRob < nActionValue {
		mlog.Error("非法抢庄操作，金币不够。 by seatno=%d", nSeatNo)
		return ERROR_ILL
	}
	userInfo.setRob(true)
	userInfo.setAction(true)
	userInfo.nRateRob = nActionValue
	//广播给所有玩家操作
	BroadCastActionMsg(ttz, nSeatNo, nActionType, nActionValue)
	ttz.GameLog.AddGameStep("玩家操作[抢庄]：用户ID[%d],座位号[%d],操作倍数[%d]", userInfo.User.UserID, userInfo.User.SeatNo, nActionValue)
	return 0
}

//是否结束抢庄阶段
func (ttz *TtzLogic) isEndRobButton() bool {
	if ttz.nGameState != GAMEPROGRESS_ROBBUTTON {
		return false
	}
	for _, ud := range ttz.mapGameUserData {
		if !ud.getAction() {
			return false
		}
	}
	return true
}

//是否达到结束游戏条件
func (ttz *TtzLogic) isEndGame() bool {
	return false
}

//是否结束当前操作状态
func (ttz *TtzLogic) isEndCurState() bool {
	bReturn := false
	defer func() {
		mlog.Debug("isEndCurState,当前场景[%d]是否结束[%v]", ttz.nGameState, bReturn)
	}()
	for _, ud := range ttz.mapGameUserData {
		if !ud.getAction() {
			return bReturn
		}
	}
	ttz.removeTimeoutAction()
	bReturn = true
	return bReturn
}

//重置操作标志
func (ttz *TtzLogic) resetActionState() {
	for _, ud := range ttz.mapGameUserData {
		ud.setAction(false)
	}
}

//结算
func (ttz *TtzLogic) gameSettlement() {
	mlog.Debug("settleMent begin")
	udButton := ttz.getUserInfoBySeatNo(ttz.nButtonSeatNo)
	if udButton == nil {
		mlog.Error("settleMent ud is nil")
		return
	}
	if ttz.nBaseScore < 1 {
		ttz.nBaseScore = 1
	}
	ttz.stSettlementInfo = make([]*settleMentInfo, 0)
	for seatNo, ud := range ttz.mapGameUserData {
		var winner, loser *TtzUserData
		if seatNo == ttz.nButtonSeatNo {
			continue
		}
		InfoData := new(settleMentInfo)
		InfoData.SeatNo = seatNo
		var winRate int32 = 1
		nResult := compareCards(udButton.userCards, ud.userCards)
		if nResult >= 0 {
			winner = udButton
			loser = ud
		} else {
			winner = ud
			loser = udButton
		}
		winRate = getRateByCardType(winner.nCardType, winner.nCardCount)
		amount := int64(winRate) * int64(ttz.nButtonValue) * int64(ud.nRateBet) * ttz.nBaseScore
		{ //Winner
			score, tax := ttz.GetUserWinAndTax(amount)
			var flowAmount, betAmount, payAmount int64
			ttz.UpdateUserScoreFunc(winner.User.UserID, score, tax)
			flowAmount = util.ABSInt64(score)
			betAmount = util.ABSInt64(score) + tax
			payAmount = betAmount + score
			ttz.GameLog.AddGameStep("用户ID %d, 增加 %d, 流水%d, 有效投注 %d, 返還金額 %d", winner.User.UserID, score, flowAmount, betAmount, payAmount)
			winner.User.FlowAmount += flowAmount
			winner.User.BetAmount += betAmount
			winner.User.PayAmount += payAmount
			winner.nCurRoundWinMoney += score
			winner.userTotalMoney += score
			winner.nCurGameWinMoney += score

			ttz.GameLog.Context().UserUpdateScore(winner.User.UserID, score)
			ttz.GameLog.Context().UserBetScore(winner.User.UserID, betAmount)
			ttz.GameLog.Context().UserPayoutScore(winner.User.UserID, payAmount)
			ttz.GameLog.Context().UserValidScore(winner.User.UserID, flowAmount)
			ttz.GameLog.Context().UserTaxAmount(winner.User.UserID, tax)
		}
		{ //Loser
			score, tax := ttz.GetUserWinAndTax(-amount)
			var flowAmount, betAmount, payAmount int64
			ttz.UpdateUserScoreFunc(winner.User.UserID, score, tax)
			flowAmount = util.ABSInt64(score)
			betAmount = util.ABSInt64(score) + tax
			payAmount = betAmount + score
			ttz.GameLog.AddGameStep("用户ID %d, 增加 %d, 流水%d, 有效投注 %d, 返還金額 %d", loser.User.UserID, score, flowAmount, betAmount, payAmount)
			loser.User.FlowAmount += flowAmount
			loser.User.BetAmount += betAmount
			loser.User.PayAmount += payAmount
			loser.nCurRoundWinMoney += score
			loser.userTotalMoney += score
			loser.nCurGameWinMoney += score

			ttz.GameLog.Context().UserUpdateScore(loser.User.UserID, score)
			ttz.GameLog.Context().UserBetScore(loser.User.UserID, betAmount)
			ttz.GameLog.Context().UserPayoutScore(loser.User.UserID, payAmount)
			ttz.GameLog.Context().UserValidScore(loser.User.UserID, flowAmount)
			ttz.GameLog.Context().UserTaxAmount(loser.User.UserID, tax)
		}
		amountAfterTax, _ := ttz.GetUserWinAndTax(amount)
		if winner.User.UserID == udButton.User.UserID {
			InfoData.Score = -amount
			InfoData.ButtonScore = amountAfterTax
		} else {
			InfoData.Score = amountAfterTax
			InfoData.ButtonScore = -amount
		}
		ttz.stSettlementInfo = append(ttz.stSettlementInfo, InfoData)
	}
}

/*
//结算
func (ttz *TtzLogic) settleMent() {
	mlog.Debug("settleMent begin")
	udButton := ttz.getUserInfoBySeatNo(ttz.nButtonSeatNo)
	if udButton == nil {
		mlog.Error("settleMent ud is nil")
		return
	}
	if ttz.nBaseScore < 1 {
		ttz.nBaseScore = 1
	}
	ttz.stSettlementInfo = make([]*settleMentInfo, 0)
	for seatNo, ud := range ttz.mapGameUserData {
		if seatNo == ttz.nButtonSeatNo {
			continue
		}
		InfoData := new(settleMentInfo)
		InfoData.SeatNo = seatNo
		var winRate int32 = 1
		nResult := compareCards(udButton.userCards, ud.userCards)
		if nResult >= 0 {
			winRate = getRateByCardType(udButton.nCardType, udButton.nCardCount)
		} else {
			winRate = -getRateByCardType(ud.nCardType, ud.nCardCount)
		}

		amount := int64(winRate) * int64(ttz.nButtonValue) * int64(ud.nRateBet) * ttz.nBaseScore
		if amount >= 0 {
			score, tax := ttz.GetUserWinAndTax(amount)
			ttz.UpdateUserScoreFunc(udButton.User.UserID, score, tax)
			ttz.GameLog.AddGameStep("莊家流水：用户ID[%d],增加[%d]", udButton.User.UserID, score)
			udButton.User.FlowAmount += score
			udButton.nCurRoundWinMoney += score
			udButton.userTotalMoney += score
			udButton.nCurGameWinMoney += score
			ttz.UpdateUserScoreFunc(ud.User.UserID, -amount, 0)
			ttz.GameLog.AddGameStep("玩家流水：用户ID[%d],增加[%d]", ud.User.UserID, amount)
			ud.User.FlowAmount += amount
			ud.nCurRoundWinMoney += -amount
			ud.userTotalMoney += -amount
			ud.nCurGameWinMoney += -amount
			InfoData.Score = -amount
			InfoData.ButtonScore = score
			ttz.stSettlementInfo = append(ttz.stSettlementInfo, InfoData)
			mlog.Debug("Settlement,SeatNo[%d],Win[%d];SeatNo[%d],Win[%d],tax[%d]", ttz.nButtonSeatNo, score, ud.User.SeatNo, -amount, tax)
		} else {
			score, tax := ttz.GetUserWinAndTax(-amount)
			ttz.UpdateUserScoreFunc(udButton.User.UserID, amount, 0)
			ttz.GameLog.AddGameStep("莊家流水：用户ID[%d],增加[%d]", udButton.User.UserID, -amount)
			udButton.User.FlowAmount += -amount
			udButton.nCurRoundWinMoney += amount
			udButton.userTotalMoney += amount
			udButton.nCurGameWinMoney += amount
			ttz.UpdateUserScoreFunc(ud.User.UserID, score, tax)
			ttz.GameLog.AddGameStep("玩家流水：用户ID[%d],增加[%d]", ud.User.UserID, score)
			ud.User.FlowAmount += score
			ud.nCurRoundWinMoney += score
			ud.userTotalMoney += score
			ud.nCurGameWinMoney += score
			InfoData.Score = score
			InfoData.ButtonScore = amount
			ttz.stSettlementInfo = append(ttz.stSettlementInfo, InfoData)
			mlog.Debug("Settlement,SeatNo[%d],Win[%d];SeatNo[%d],Win[%d];tax[%d]", ttz.nButtonSeatNo, amount, ud.User.SeatNo, score, tax)
		}
	}
}
*/
//玩家操作下注
func (ttz *TtzLogic) actionBet(nSeatNo, nActionType, nActionValue int32) int32 {
	//mlog.Debug("actionRob,nSeatNo=%d,nActionType=%d.nActionValue=%d", nSeatNo, nActionType, nActionValue)
	//逻辑操作
	if nSeatNo == ttz.nButtonSeatNo {
		mlog.Error("非法操作，不允许庄家下注 by seatno=%d", nSeatNo)
		return -1
	}
	userInfo := ttz.getUserInfoBySeatNo(nSeatNo)
	if userInfo == nil {
		mlog.Error("userinfo is not find by seatno=%d", nSeatNo)
		return ERROR_USERNOTEXITST
	}
	if userInfo.getBet() {
		mlog.Error("重复操作下注 by seatno=%d", nSeatNo)
		return ERROR_REPEATOP
	}
	if userInfo.nMaxBet < nActionValue {
		mlog.Error("非法下注操作，金币不够。 by seatno=%d", nSeatNo)
		return ERROR_ILL
	}
	userInfo.setBet(true)
	userInfo.setAction(true)
	userInfo.nRateBet = nActionValue
	//广播给所有玩家操作
	BroadCastActionMsg(ttz, nSeatNo, nActionType, nActionValue)
	ttz.GameLog.AddGameStep("玩家操作[下注]：用户ID[%d],座位号[%d],操作倍数[%d]", userInfo.User.UserID, userInfo.User.SeatNo, nActionValue)
	return 0
}

//玩家操作亮牌
func (ttz *TtzLogic) actionLight(nSeatNo int32) int32 {
	//mlog.Debug("actionRob,nSeatNo=%d,nActionType=%d.nActionValue=%d", nSeatNo, nActionType, nActionValue)
	//逻辑操作
	userInfo := ttz.getUserInfoBySeatNo(nSeatNo)
	if userInfo == nil {
		mlog.Error("userinfo is not find by seatno=%d", nSeatNo)
		return ERROR_USERNOTEXITST
	}
	if userInfo.getLight() {
		mlog.Error("重复操作亮牌 by seatno=%d", nSeatNo)
		return ERROR_REPEATOP
	}
	userInfo.setLight(true)
	userInfo.setAction(true)
	//广播给所有玩家操作
	BroadCastLight(ttz, nSeatNo)
	ttz.GameLog.AddGameStep("玩家操作[亮牌]：用户ID[%d],座位号[%d]", userInfo.User.UserID, userInfo.User.SeatNo)
	return 0
}

//设置庄家信息
func (ttz *TtzLogic) setButtonInfo() {
	defer func() {
		mlog.Debug("setButtonInfo 设置庄家信息，ButtonSeatNo=%d,ButtonValue=%d", ttz.nButtonSeatNo, ttz.nButtonValue)
	}()
	//获取最大的抢庄玩家倍率
	var maxRateRob int32 = 0
	for _, ud := range ttz.mapGameUserData {
		if maxRateRob < ud.nRateRob {
			maxRateRob = ud.nRateRob
		}
	}
	//获取所有最大抢庄玩家
	mapMaxSeatNo := make(map[int32]int32, 0)
	for seatNo, ud := range ttz.mapGameUserData {
		if maxRateRob == ud.nRateRob {
			mapMaxSeatNo[seatNo] = maxRateRob
		}
	}
	//随机庄家
	for seatNo, _ := range mapMaxSeatNo {
		ttz.nButtonSeatNo = seatNo
		break
	}
	if maxRateRob == 0 {
		ttz.nButtonValue = 1
	} else {
		ttz.nButtonValue = maxRateRob
	}
	ttz.GameLog.AddGameStep("抢庄结束，设定庄家信息：座位号[%d],庄家倍率[%d]", ttz.nButtonSeatNo, ttz.nButtonValue)
	user := ttz.getUserInfoBySeatNo(ttz.nButtonSeatNo)
	if user != nil {
		user.setAction(true)
		user.nRateRob = ttz.nButtonValue
		ttz.nLastDealId = user.User.UserID
	}
}

// 随机一个数
func (ttz *TtzLogic) randOneNumber(number int) int32 {
	nSeed := time.Now().UnixNano()
	nSource := rand.NewSource(nSeed)
	rRand := rand.New(nSource)
	Index := rRand.Intn(number)
	return int32(Index)
}

// 更新发牌玩家信息
func (ttz *TtzLogic) updateDealSeatInit() {
	if ttz.nGameIndex > 0 {
		for _, ud := range ttz.mapGameUserData {
			if ud.User.UserID == ttz.nLastDealId {
				ttz.nLastDealSeat++
				for i := 0; i < 10; i++ {
					user := ttz.getUserInfoBySeatNo(ttz.nLastDealSeat)
					if user != nil {
						ttz.nLastDealId = user.User.UserID
						break
					} else {
						ttz.nLastDealSeat++
						if ttz.nLastDealSeat > 4 {
							ttz.nLastDealSeat = 0
						}
					}
				}
				return
			}
		}
		ttz.nLastDealSeat = 0
		for i := 0; i < 10; i++ {
			user := ttz.getUserInfoBySeatNo(0)
			if user != nil {
				ttz.nLastDealId = user.User.UserID
				break
			} else {
				ttz.nLastDealSeat++
				if ttz.nLastDealSeat > 4 {
					ttz.nLastDealSeat = 0
				}
			}
		}

	} else {
		ttz.nLastDealSeat = 0
		for i := 0; i < 10; i++ {
			user := ttz.getUserInfoBySeatNo(0)
			if user != nil {
				ttz.nLastDealId = user.User.UserID
				break
			} else {
				ttz.nLastDealSeat++
				if ttz.nLastDealSeat > 4 {
					ttz.nLastDealSeat = 0
				}
			}
		}
	}
}

// 计算玩家最大抢庄倍率
func (ttz *TtzLogic) calMaxUserCanRob() {
	for _, ud := range ttz.mapGameUserData {
		ud.nMaxRob = int32(ud.userTotalMoney / (4 * 8 * int64(ttz.nGameCount-1) * ttz.nBaseScore))
		//广播
		BroadCastUserMaxOpRate(ttz, ud.User, ud.nMaxRob, 1)
	}
}

// 计算玩家最大下注倍率
func (ttz *TtzLogic) calMaxUserCanBet() {
	for _, ud := range ttz.mapGameUserData {
		ud.nMaxBet = int32(ud.userTotalMoney / (int64(ttz.nButtonValue) * 4 * 1 * ttz.nBaseScore))
		//广播
		BroadCastUserMaxOpRate(ttz, ud.User, ud.nMaxBet, 2)
	}
}

//换最大牌
func (ttz *TtzLogic) ExchangeMaxCard(userid int32) {
	ttz.LogicBase.ExchangeMaxCard(userid)
	var maxRetUserID int32 = 0
	var nType, nCount, nTypeMax, nCountMax int32 = 0, 0, 0, 0
	if ttz.pokerGroup == nil {
		mlog.Error("换牌失败，牌库为空")
		return
	}
	fpoker := ttz.pokerGroup.(*poker.FixedPoker)
	for _, ur := range ttz.AllGameUser() {
		if ur.IsRobot() || ur.UserID == userid {
			continue
		}
		cg := fpoker.GetUserCards(ur.UserID)
		nType, nCount = getCardTypeAndCount(cg)

		//比较牌型
		if nType > nTypeMax {
			nTypeMax = nType
			nCountMax = nCount
			maxRetUserID = ur.UserID
		} else if nType == nTypeMax {
			if nCount > nCountMax {
				nTypeMax = nType
				nCountMax = nCount
				maxRetUserID = ur.UserID
			}
		}
	}

	cgMax := fpoker.GetUserCards(maxRetUserID)
	nTypeMax1, nCountMax1 := getCardTypeAndCount(cgMax)

	cgur := fpoker.GetUserCards(userid)
	nTypeUr, nCountUr := getCardTypeAndCount(cgur)

	//比较牌型
	if nTypeUr > nTypeMax {
		return
	} else if nTypeMax1 == nTypeUr {
		if nCountUr > nCountMax1 {
			return
		}
	}

	ttz.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	fpoker.ExChange(userid, maxRetUserID)

}
