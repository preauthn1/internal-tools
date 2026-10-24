//==================================================================
//文件名： dzpklogic.go
//说明：   德州扑克游戏战斗逻辑
//修改记录：
//==================================================================
package dzpk

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"

	"math"

	proto "example.invalid/golang/protobuf/proto"
)

// 德州扑克逻辑数据结构
type DzpkLogic struct {
	*gamelogic.DynamicGameLogicBase
	opttimer           *event.TimerFunc
	dzpkCfgs           *dzpkCfg
	poker              poker.Poker        //游戏扑克牌
	gameUserData       []*DzpkUserData    //游戏用户数据
	watchUserData      []*DzpkUserData    //观战用户数据
	gameUserCount      int32              //游戏用户人数
	watchUserCount     int32              //观战人数
	buttonSeat         int32              //庄家座位号
	nGun               int32              //枪口位置
	currActionSeat     int32              //当前操作的座位号
	opLastType         int32              //上一次操作类型
	opLastSeat         int32              //上一次操作位置
	optTime            int32              //可操作时间
	beginOutTime       int32              //操作开始时间
	betBase            int64              //底注
	betCurr            int64              //当前的下注额
	betCurrMinCall     int64              //当前跟注额
	betCurrMinRaise    int64              //当前最小加注额
	betTotal           int64              //当前下注总额
	betCurrMax         int64              //当前最大下注额
	potArray           *potPool           //彩池
	gameState          gameProgress       //游戏状态
	publicCards        *poker.CardGroup   //用户公共牌
	gameType           playType           //玩法类型
	nCountOpNumber     int32              //当前操作的玩家数量
	nFoldNumber        int32              //弃牌人数
	nAllinNumber       int32              //allin 人数
	nMaxCardTypeSeatNo int32              //牌型最大玩家座位号
	settlePotInfo      map[int32]*potInfo //小局结算
	actionName_curr    string             //当前动画
	actionGameEndTime  int32              //结算动画超时时间
	newGameBeginTime   int32              //新一局游戏开始时间
	nPokerType         int32              //发牌类型
	nGameIndex         int32              //游戏局数
	nAnimationDy       int32              //动态结束时间
	nGM                int32              //是否打开测试命令
	savePublic         *poker.CardGroup   //用户公共牌，不做改动
}

//清除数据，初始化状态
func (dzpk *DzpkLogic) clearGameData() {
	dzpk.setGameUserData()
	//dzpk.optTime = 60
	//dzpk.buttonSeat = 0
	dzpk.currActionSeat = -1
	dzpk.opLastType = 0
	dzpk.opLastSeat = 0
	dzpk.betCurr = 0
	dzpk.betCurrMinCall = 0
	dzpk.betCurrMinRaise = 0
	dzpk.betTotal = 0
	dzpk.gameState = GAMEPROGRESS_BEGIN
	dzpk.nCountOpNumber = 0
	dzpk.nFoldNumber = 0
	dzpk.nAllinNumber = 0
	dzpk.nMaxCardTypeSeatNo = -1
	dzpk.betCurrMax = 0
	dzpk.actionName_curr = ""
	dzpk.beginOutTime = 0
	dzpk.nGun = -1
	dzpk.nGameIndex++
	dzpk.publicCards = nil
	dzpk.updateButtonInit()

	if 0 == dzpk.nPokerType {
		dzpk.poker = poker.NewRandomPoker(false, false, 1) //初始化一副牌
		dzpk.poker = poker.NewFixedPoker(dzpk.AllTableUser(), dzpk.poker.(*poker.RandomPoker), 2)
		dzpk.savePublic = dzpk.poker.Deal(poker.RestCardsUserID, 5)
	} else {
		dzpk.poker = poker.NewFilePoker(dzpk.Configs, dzpk.AllTableUser())
		dzpk.poker = poker.NewFixedFilePoker(dzpk.poker.(*poker.FilePoker))
		dzpk.savePublic = dzpk.poker.Deal(poker.RestCardsUserID, 5)
		if dzpk.savePublic == nil {
			dzpk.endGame()
			return
		}
		if dzpk.savePublic.Len() < 5 {
			dzpk.endGame()
			return
		}

	}
	dzpk.poker.Shuffle()

}

// 重置游戏逻辑数据
func (dzpk *DzpkLogic) resetRoundData() {
	mlog.Debug("开始重置小局数据")
	//dzpk.setGameUserData()
	//重置牌局部分信息
	dzpk.currActionSeat = 0
	dzpk.opLastType = 0
	dzpk.opLastSeat = 0
	dzpk.betCurr = 0
	dzpk.betCurrMinCall = 0
	dzpk.betCurrMinRaise = 0
	dzpk.betTotal = 0
	dzpk.gameState = GAMEPROGRESS_BEGIN
	//重置用户部分数据
	dzpk.nCountOpNumber = 0
	dzpk.nFoldNumber = 0
	dzpk.nAllinNumber = 0
	dzpk.nMaxCardTypeSeatNo = -1
	dzpk.betCurrMax = 0
	dzpk.actionName_curr = ""
	dzpk.nGun = -1
	dzpk.publicCards = nil
	for _, userdata := range dzpk.gameUserData {
		userdata.userResetData()
	}
}

// 创建德州扑克逻辑
func NewDzpkLogic(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *DzpkLogic {
	//mlog.Debug("创建德州扑克逻辑")
	glc := new(DzpkLogic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(t, tfManager, cfgs)

	tmpflag, _ := strconv.Atoi(t.Room.Flag)
	glc.gameType = playType(tmpflag)
	glc.getGameConfig(cfgs)
	glc.buttonSeat = -1
	glc.nGameIndex = 0
	//------
	glc.RegisterAfterGameStartFunc(func() {
		glc.GameLog.AddGameStep("开始游戏,桌子号%d,第%d局！", glc.TableNo(), glc.nGameIndex)
		mlog.Debug("开始游戏,桌子号%d,第%d局！", glc.TableNo(), glc.nGameIndex)
		//更新牌桌数据
		glc.updateRoleSeatNo(true)
		glc.betInit() //初始化筹码
		glc.deal()    //发牌
		glc.waitOpt() //第一个人操作
	})
	glc.RegisterBeforeGameStartFunc(func() {
		//清除数据
		glc.clearGameData() //清除数据
	})
	return glc
}

// 德州扑克消息管理
func (dzpk *DzpkLogic) GameMessage(user *gamelogic.GameUser, gmsg *msg.Message) {
	mlog.Debug("GameMessage, receive message，SeatNo：%d", user.SeatNo)
	dzpk.CheckUserAnimationCommited(user, gmsg, dzpk.actionName_curr, nil)
	if gmsg.BClassID == int32(netproto.MessageBClassID_Game) && gmsg.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		mlog.Debug("GameMessage,animation commite SeatNo=%d,Name=%s", user.SeatNo, dzpk.actionName_curr)
		return
	}
	//检测是不是测试命令
	if dzpk.nGM == 1 {
		nRet := dzpk.setTestCommond(gmsg)
		if nRet <= 0 {
			if nRet != 0 {
				mlog.Error("user input error commond.")
				return
			}
			return
		}
	}
	//操作合法性校验
	if !dzpk.oprateCheckLegal(user, gmsg) {
		return
	}
	//战斗逻辑操作
	dzpk.oprateBattleLogic(user, gmsg)
}

// 更新庄家座位号
func (dzpk *DzpkLogic) updateButtonInit() {
	dzpk.buttonSeat++
	for i := 0; i < 20; i++ {
		if dzpk.buttonSeat > 9 {
			dzpk.buttonSeat = 0
		}
		user := dzpk.getUserInfoBySeatNo(dzpk.buttonSeat)
		if user != nil {
			return
		}
		dzpk.buttonSeat++
	}
}

//初始化下注
func (dzpk *DzpkLogic) betInit() {
	for _, urd := range dzpk.gameUserData {
		money := dzpk.getInitJetton(urd)
		urd.SetMoney(-money)
		//dzpk.UpdateUserScoreFunc(urd.User.UserID, -money, 0)
		dzpk.betTotal += money
		dzpk.GameLog.AddGameStep("用户下底：用户ID[%d],座位号[%d],筹码[%d].", urd.User.UserID, urd.seatNo, money)
	}

	dzpk.updateCurBetMaxJetton()
	sendBetBaseMsg(dzpk)
	dzpk.gameState = GAMEPROGRESS_DEAL
}

//获取当前游戏配置
func (dzpk *DzpkLogic) getGameConfig(cfgs *config.Group) {
	cfg := NewDzpkCfg(cfgs)
	dzpk.dzpkCfgs = cfg
	dzpk.actionGameEndTime = int32(cfg.GetValueInt("GameEndTime"))
	dzpk.optTime = int32(cfg.GetValueInt("ActionTime"))
	dzpk.nPokerType = int32(cfg.GetValueInt("PokerType"))
	dzpk.nGM = int32(cfg.GetValueInt("GM"))
	dzpk.betBase = int64(cfg.GetValueInt("BaseScore"))
	if dzpk.betBase < 1 {
		dzpk.betBase = 1
	}

}

//发牌，每人发两张底牌
func (dzpk *DzpkLogic) deal() {
	if dzpk.gameState != GAMEPROGRESS_DEAL {
		mlog.Error("Game state error, deal card failed.")
		return
	}

	//给所有用户发牌
	for _, ur := range dzpk.AllTableUser() {
		cards := dzpk.poker.Deal(ur.UserID, 2)
		if cards != nil {
			sort.Sort(cards) //从大到小
		} else {
			dzpk.endGame()
		}

		mlog.Debug("Deal cards begin,SeatNo=%d", ur.SeatNo)
		dzpk.setUserCards(ur, cards)
	}
	//获取牌型
	for i, ur := range dzpk.gameUserData {
		dzpk.gameUserData[i].cardType = GetCardType(ur.userCards)
		//牌局记录
		dzpk.GameLog.AddGameStep("发放手牌：用户ID[%d],座位号[%d],手牌[%s],牌型[%s]", ur.User.UserID, ur.seatNo, ur.userCards.String(), strCardType(ur.cardType))
	}

	//广播出去
	sendDealMsg(dzpk)
	//dzpk.updateCardType()
	mlog.Debug("Deal cards end.")
	dzpk.gameState = GAMEPROGRESS_PREFLOP
}

// 更新牌型
func (dzpk *DzpkLogic) updateCardType() {
	for _, ur := range dzpk.gameUserData {
		cardGroup := new(poker.CardGroup)

		if ur.userCards != nil {
			cardGroup.Cards = append(cardGroup.Cards, ur.userCards.Cards...)
		}

		if dzpk.publicCards != nil {
			cardGroup.Cards = append(cardGroup.Cards, dzpk.publicCards.Cards...)
		}
		ur.cardType = GetCardType(cardGroup)
		BroadCastUserCardType(dzpk, ur)
		//mlog.Debug("更新牌型信息，SeatNo=%d,CardSize=%d,CardType=%d", ur.seatNo, cardGroup.Len(), ur.cardType)
	}
}

//设置指定牌给指定玩家
func (dzpk *DzpkLogic) setUserCards(ur *gamelogic.GameUser, cards *poker.CardGroup) {
	for _, urd := range dzpk.gameUserData {
		if urd.seatNo == ur.SeatNo {
			urd.userCards = poker.NewCardGroup(cards.Cards)
			break
		}
	}

}

//allin 和弃牌造成的提前结算是否需要继续发牌
//allin人数大于1需要发牌
//all人数等于1，还有未弃牌的玩家时候需要发牌，allin+未弃牌>1 == allin+弃牌<游戏人数
func (dzpk *DzpkLogic) continueDeal() bool {
	if dzpk.nAllinNumber > 1 {
		return true
	}
	if dzpk.nAllinNumber == 1 && dzpk.nAllinNumber+dzpk.nFoldNumber < dzpk.gameUserCount {
		return true
	}
	return false
}

//发送初始化操作信息,获取下次操作数据
func (dzpk *DzpkLogic) waitOpt() {
	mlog.Debug("waitOpt loop,gameState=%d,opration number=%d。", dzpk.gameState, dzpk.nCountOpNumber)
	dzpk.updateCurBetMaxJetton()

	if dzpk.isEndTurn() { /*检测是否结束本轮*/
		dzpk.nCountOpNumber = 0
		dzpk.resetUserOpState()
		dzpk.updatePotData()
		mlog.Debug("结束本轮")
		if dzpk.isEndRound() { /*检测是否结束小局游戏*/
			mlog.Debug("游戏小局结束")

			if dzpk.continueDeal() {

				if dzpk.publicCards == nil {
					dzpk.publicCards = new(poker.CardGroup)
				}
				number := 5 - dzpk.publicCards.Len()
				if number > 0 {
					dzpk.publicDeal(number)
				}
			}
			dzpk.gameState = GAMEPROGRESS_END
			if dzpk.isEndGame() { /*检测是否结束游戏*/
				mlog.Debug("游戏结束")
				//游戏结束
				return
			}
			//结算
			dzpk.settlementRound()
			mlog.Debug("Register animation %s,outtime=%ds", ACTION_NAME_GAMEEND, dzpk.nAnimationDy)

			//牌局记录
			uscoresRet := make(map[int32]int64)
			for _, ud := range dzpk.gameUserData {
				uscoresRet[ud.User.UserID] = ud.curRoundSettle
				dzpk.GameLog.Context().UserUpdateScore(ud.User.UserID, ud.curRoundSettle)
			}
			dzpk.GameLog.GameEnd()

			dzpk.actionName_curr = dzpk.RegisterAnimation(ACTION_NAME_GAMEEND, dzpk.nAnimationDy, dzpk.endGame)
			return

		} else {
			dzpk.updateRoundData() //更新数据
		}

	}

	if dzpk.gameState == GAMEPROGRESS_PREFLOP && dzpk.nCountOpNumber == 0 {
		dzpk.calFirstOpData() //下一次操作玩家数据
	} else {
		dzpk.calNextOpData() //下一次操作玩家数据
	}
	dzpk.preOpration() //预操作信息
	//广播给游戏中的玩家,下一个操作玩家位置
	var waitTimeAutoWaitMsg int32 = 2
	dzpk.TimerFuncManager.NewNamedTimer(waitTimeAutoWaitMsg, dzpk.autoBroadCastWaitOptMsg, opttimerWait)
	//BroadCastWaitOptMsg(dzpk)
	//广播给观战中的玩家

	//设置超时自动处理
	dzpk.beginOutTime = int32(time.Now().Unix())
	if dzpk.opttimer != nil {
		mlog.Debug("移除超时timer")
		dzpk.opttimer.Stop()
		dzpk.opttimer = nil
	}

	mlog.Debug("注册超时timer,autoFold timeout=%ds", waitTimeAutoWaitMsg+dzpk.optTime)
	dzpk.opttimer = dzpk.TimerFuncManager.NewNamedTimer(waitTimeAutoWaitMsg+dzpk.optTime, dzpk.autoFold, "")
}

// 自动操作 等待
func (dzpk *DzpkLogic) autoBroadCastWaitOptMsg() {
	//广播给游戏中的玩家,下一个操作玩家位置
	BroadCastWaitOptMsg(dzpk)
}

//是否满足结算条件
func (dzpk *DzpkLogic) settlementState() bool {
	var count int32 = 0
	for _, user := range dzpk.gameUserData {
		if user.isAllin {
			count++
		}
		if user.isFold {
			count++
		}
	}
	if count+1 == dzpk.gameUserCount {
		return true
	}
	if dzpk.gameState == GAMEPROGRESS_END {
		return true
	}
	return false
}

//计算第一次操作玩家数据
func (dzpk *DzpkLogic) calFirstOpData() {
	mlog.Debug("calFirstOpData begin.")
	dzpk.currActionSeat = dzpk.calNextOpSeat()
	user := dzpk.getUserInfoBySeatNo(dzpk.currActionSeat)
	if user != nil {
		user.updateUserState(dzpk.betCurrMax, dzpk.betBase/2)
		user.setBeginOprateForBet(dzpk.betCurrMax)
	}
}

//计算下一次操作玩家数据
func (dzpk *DzpkLogic) calNextOpData() {
	mlog.Debug("calNextOpData begin.")
	dzpk.currActionSeat = dzpk.calNextOpSeat()
	user := dzpk.getUserInfoBySeatNo(dzpk.currActionSeat)
	if user != nil {
		user.updateUserState(dzpk.betCurrMax, dzpk.betBase/2)
		if dzpk.nCountOpNumber == 0 {
			user.setBeginOprateForBet(dzpk.betCurrMax)
		} else {
			user.updateUserOpInfo(dzpk.betCurrMax, dzpk.betBase/2)
		}

	}
}

//计算下一次操作玩家数据
func (dzpk *DzpkLogic) calNextOpSeat() int32 {
	if dzpk.nGun >= 0 {
		seat := dzpk.nGun
		dzpk.nGun = -1
		mlog.Debug("遍历座位号gun nextSeat:%v", seat)
		return seat
	}
	var nextSeat int32 = -1
	for _, ud := range dzpk.gameUserData {
		if ud.seatNo > dzpk.currActionSeat {
			if dzpk.isValid(ud) {
				mlog.Debug("遍历座位号nextSeat:%v", ud.seatNo)
				return ud.seatNo
			}
		}
	}
	for _, ud := range dzpk.gameUserData {
		if ud.seatNo > nextSeat {
			if dzpk.isValid(ud) {
				mlog.Debug("遍历座位号nextSeat:%v", ud.seatNo)
				return ud.seatNo
			}
		}
	}
	mlog.Error("遍历座位号nextSeat:%v", -1)
	return -1
}

//玩家操作是否合法
func (dzpk *DzpkLogic) isValid(user *DzpkUserData) bool {
	if user.isAllin == true {
		return false
	}
	if user.isFold == true {
		return false
	}
	return true
}

//根据座位号获取用户信息
func (dzpk *DzpkLogic) getUserInfoBySeatNo(seatno int32) *DzpkUserData {
	for _, urd := range dzpk.gameUserData {
		if urd.User.SeatNo == seatno {
			return urd
		}
	}
	return nil
}

//操作合法性校验
func (dzpk *DzpkLogic) oprateCheckLegal(user *gamelogic.GameUser, gmsg *msg.Message) bool {

	if !dzpk.CheckGameUserMessage(user) {
		dzpk.KickUserFunc(user.UserID, "非法操作")
		return false
	}

	//不是给游戏逻辑服务器的消息
	if gmsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return false
	}

	if user.SeatNo != dzpk.currActionSeat {
		mlog.Error("用户座位号非当前操作用户座位号%v", user)
		return false
	}

	//没有找到相关玩家数据
	urd := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if urd == nil {
		mlog.Error("用户未找到%v", user)
		return false
	}

	//用户不可操作
	if !urd.userCanOprate() {
		mlog.Error("用户已经弃牌或Allin%v", user)
		return false
	}
	return true
}

//战斗逻辑操作
func (dzpk *DzpkLogic) oprateBattleLogic(user *gamelogic.GameUser, gmsg *msg.Message) {

	var nRet int32 = 0
	switch gmsg.SClassID {
	case int32(netproto.DZPK_GameMessageClassID_ActionCallID):
		//用户跟注
		nRet = dzpk.actionCall(user, gmsg.MsgData.([]byte))
	case int32(netproto.DZPK_GameMessageClassID_ActionRaiseID):
		//用户加注
		nRet = dzpk.actionRaise(user, gmsg.MsgData.([]byte))
	case int32(netproto.DZPK_GameMessageClassID_ActionCheckID):
		//用户看牌，过
		nRet = dzpk.actionCheck(user, gmsg.MsgData.([]byte))
	case int32(netproto.DZPK_GameMessageClassID_ActionAllinID):
		//ALLIN
		nRet = dzpk.actionAllin(user, gmsg.MsgData.([]byte))
	case int32(netproto.DZPK_GameMessageClassID_ActionFoldID):
		//用户弃牌
		nRet = dzpk.actionFold(user)
	default:
		mlog.Error("unknow class id,gmsg.SClassID=%d", gmsg.SClassID)
		nRet = ERROR_CLASSID
	}

	if nRet != 0 {
		mlog.Error("oprateBattleLogic failed,SeatNo=%d.", user.SeatNo)
		BroadCastErrMsg(dzpk, user.SeatNo, nRet, "action error.")
		return
	}
	//检测是否回合结束，如果不是直接返回，如果是（更新pot，进入下一阶段）
	dzpk.nCountOpNumber++
	dzpk.waitOpt()

}

//跟注操作
func (dzpk *DzpkLogic) actionCall(user *gamelogic.GameUser, data []byte) int32 {
	mlog.Debug("action actionCall")
	//解析数据
	ust := &netproto.DZPK_Call{}
	err := proto.Unmarshal(data, ust)
	//错误校验
	if err != nil {
		mlog.Error("解析用户%v用户座位号错误%v", user, err)
		return ERROR_OPRATE
	}
	*ust.Jetton = int64(math.Abs(float64(*ust.Jetton)))

	//逻辑操作
	userInfo := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if userInfo == nil {
		return ERROR_OPRATE
	}
	if !userInfo.isCanCall {
		mlog.Error("actionCall, seatNo[%d] can't use op Call！", user.SeatNo)
		return ERROR_OPRATE
	}
	jetton := dzpk.betCurrMax - userInfo.currBet
	if jetton > userInfo.userMoney {
		mlog.Error("actionCall illegal op, gold[%d] inadequate", userInfo.userMoney)
		return ERROR_NOMONEY
	}
	//移除操作超时Timer
	if dzpk.opttimer != nil {
		dzpk.opttimer.Stop()
	}
	dzpk.opLastType = int32(CALL)
	userInfo.isAction = true
	userInfo.curOpration = int32(CALL)
	userInfo.SetMoney(-jetton)
	dzpk.betCurr += jetton
	dzpk.betTotal += jetton
	BroadCastCallMsg(dzpk, user.SeatNo, jetton)
	dzpk.GameLog.AddGameStep("玩家操作[跟注],游戏阶段[%s],用户ID[%d],座位号[%d],筹码[%d]", getState(dzpk.gameState), userInfo.User.UserID, userInfo.seatNo, jetton)
	return 0
}

//加注
func (dzpk *DzpkLogic) actionRaise(user *gamelogic.GameUser, data []byte) int32 {
	mlog.Debug("action actionRaise")
	//解析数据
	ust := &netproto.DZPK_Raise{}
	err := proto.Unmarshal(data, ust)
	//错误校验
	if err != nil {
		mlog.Error("解析用户%v用户座位号错误%v", user, err)
		return ERROR_OPRATE
	}
	*ust.Jetton = int64(math.Abs(float64(*ust.Jetton)))

	//逻辑操作
	userInfo := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if userInfo == nil {
		return ERROR_OPRATE
	}
	if !userInfo.isCanRaise {
		mlog.Error("actionRaise, seatNo[%d] can't use op Raise!", user.SeatNo)
		return ERROR_OPRATE
	}

	jetton := *ust.Jetton

	//有限注調整
	if true && !user.IsRobot() {
		nPool := dzpk.betTotal                      //底池
		curMax := dzpk.betCurrMax                   //当前牌局最大下注额
		ud := dzpk.getUserInfoBySeatNo(user.SeatNo) //玩家下注额
		switch jetton {
		case 50:
			jetton = (nPool+curMax-ud.currBet)/2 + curMax - ud.currBet
		case 66:
			jetton = (nPool+curMax-ud.currBet)*2/3 + curMax - ud.currBet
		case 100:
			jetton = (nPool + curMax - ud.currBet) + curMax - ud.currBet
		default:
			jetton = (nPool+curMax-ud.currBet)/2 + curMax - ud.currBet
		}
	}

	if userInfo.userMoney < jetton {
		mlog.Error("actionRaise,be short of money.")
		return ERROR_NOMONEY
	}
	if jetton+userInfo.currBet < dzpk.betCurrMax+dzpk.betBase/2 {
		mlog.Error("actionRaise illegal op，gold must more than %d！", dzpk.betCurrMax-userInfo.currBet+dzpk.betBase/2)
		return ERROR_NOMONEY
	}
	//移除操作超时Timer
	if dzpk.opttimer != nil {
		dzpk.opttimer.Stop()
	}
	dzpk.opLastType = int32(RAISE)
	userInfo.curOpration = int32(RAISE)
	userInfo.isAction = true

	userInfo.SetMoney(-jetton)
	dzpk.betCurr += jetton
	dzpk.betTotal += jetton
	//广播给所有玩家操作
	BroadCastRaiseMsg(dzpk, user.SeatNo, jetton)
	dzpk.GameLog.AddGameStep("玩家操作[加注],游戏阶段[%s],用户ID[%d],座位号[%d],筹码[%d]", getState(dzpk.gameState), userInfo.User.UserID, userInfo.seatNo, jetton)
	return 0
}

//过，看牌
func (dzpk *DzpkLogic) actionCheck(user *gamelogic.GameUser, data []byte) int32 {
	mlog.Debug("action actionCheck")
	//逻辑，设置已操作状态
	userInfo := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if userInfo == nil {
		return ERROR_OPRATE
	}
	if !userInfo.isCanCheck {
		mlog.Error("actionCheck, seatNo[%d] can't use op Check!", user.SeatNo)
		return ERROR_OPRATE
	}
	//移除操作超时Timer
	if dzpk.opttimer != nil {
		dzpk.opttimer.Stop()
	}
	dzpk.opLastType = int32(CHECK)
	userInfo.curOpration = int32(CHECK)
	userInfo.isAction = true
	//广播给所有玩家操作
	BroadCastCheckMsg(dzpk, user.SeatNo)
	dzpk.GameLog.AddGameStep("玩家操作[看牌],游戏阶段[%s],用户ID[%d],座位号[%d],筹码[%d]", getState(dzpk.gameState), userInfo.User.UserID, userInfo.seatNo, 0)
	return 0
}

//Allin
func (dzpk *DzpkLogic) actionAllin(user *gamelogic.GameUser, data []byte) int32 {
	mlog.Debug("action Allin")
	//解析数据
	userInfo := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if userInfo == nil {
		mlog.Error("actionAllin error user！")
		return ERROR_OPRATE
	}
	if !(userInfo.isCanAllIn || userInfo.isCanRaise) {
		mlog.Error("actionAllin, seatNo[%d] can't use op Allin!", user.SeatNo)
		return ERROR_OPRATE
	}

	jetton := userInfo.userMoney
	if jetton <= 0 {
		mlog.Error("actionAllin非法操作,this user[seatNo=%d] have no money！", user.SeatNo)
		return ERROR_OPRATE
	}
	//移除操作超时Timer
	if dzpk.opttimer != nil {
		dzpk.opttimer.Stop()
	}
	dzpk.opLastType = int32(ALLIN)
	userInfo.curOpration = int32(ALLIN)
	dzpk.nAllinNumber++
	userInfo.isAllin = true
	userInfo.isAction = true
	userInfo.SetMoney(-jetton)
	dzpk.betCurr += jetton
	dzpk.betTotal += jetton
	//广播给所有玩家操作
	BroadCastAllinMsg(dzpk, user.SeatNo, jetton)
	dzpk.GameLog.AddGameStep("玩家操作[全押],游戏阶段[%s],用户ID[%d],座位号[%d],筹码[%d]", getState(dzpk.gameState), userInfo.User.UserID, userInfo.seatNo, jetton)
	return 0
}

//弃牌操作
func (dzpk *DzpkLogic) actionFold(user *gamelogic.GameUser) int32 {
	mlog.Debug("action actionFold")
	//逻辑操作
	userInfo := dzpk.getUserInfoBySeatNo(user.SeatNo)
	if userInfo == nil {
		mlog.Error("actionFold,get the user[seatNo=%d] error！", user.SeatNo)
		return ERROR_OPRATE
	}
	if dzpk.opttimer != nil {
		dzpk.opttimer.Stop()
	}
	userInfo.curOpration = int32(FOLD)
	dzpk.nFoldNumber++
	userInfo.isAction = true
	userInfo.isFold = true
	//广播给所有玩家操作
	BroadCastFoldMsg(dzpk, user.SeatNo)
	dzpk.GameLog.AddGameStep("玩家操作[弃牌],游戏阶段[%s],用户ID[%d],座位号[%d],筹码[%d]", getState(dzpk.gameState), userInfo.User.UserID, userInfo.seatNo, 0)
	return 0
}

// 自动弃牌
func (dzpk *DzpkLogic) autoFold() {
	//自动弃牌
	user := dzpk.getUserInfoBySeatNo(dzpk.currActionSeat)
	if user != nil && user.isFold != true {
		dzpk.actionFold(user.User)
		dzpk.nCountOpNumber++
		dzpk.waitOpt()
	}
}

// 游戏结束是否是因为allin或者弃牌
func (dzpk *DzpkLogic) isEndByAllinAndFold() bool {
	if dzpk.nAllinNumber+dzpk.nFoldNumber+1 >= dzpk.gameUserCount {
		return true
	}
	return false
}

//判断回合是否结束
func (dzpk *DzpkLogic) isEndTurn() bool {

	for _, user := range dzpk.gameUserData {
		if user.isFold || user.isAllin {
			continue
		}
		if !user.isAction {
			if dzpk.isEndByAllinAndFold() && user.currBet == dzpk.betCurrMax {
				return true
			} else {
				return false
			}

		}
		if user.currBet != dzpk.betCurrMax {
			return false
		}
		//return false
	}
	return true
}

//判断小局游戏是否结束
func (dzpk *DzpkLogic) isEndRound() bool {

	if dzpk.isEndByAllinAndFold() {
		return true
	}

	if dzpk.gameState == GAMEPROGRESS_RIVER {
		dzpk.gameState++
		return true
	}
	if dzpk.gameState != GAMEPROGRESS_END {
		return false
	}
	return true
}

//更新小局数据
func (dzpk *DzpkLogic) updateRoundData() {
	mlog.Debug("更新小局数据")
	//重置某些数据
	if dzpk.gameState == GAMEPROGRESS_PREFLOP {
		dzpk.publicCards = new(poker.CardGroup)
		dzpk.publicDeal(3)
	} else {
		//发牌一张
		dzpk.publicDeal(1)
	}
	for _, ud := range dzpk.gameUserData {
		if ud.isAllin || ud.isFold {
			continue
		}
		ud.isAction = false
	}
	//游戏状态加1
	dzpk.gameState++

}

//发count张公共牌
func (dzpk *DzpkLogic) publicDeal(count int) int32 {
	//发牌
	mlog.Debug("发%d张公共牌", count)
	if count > 5 {
		mlog.Error("参数错误，没有这么多牌")
		return -1
	}
	//cards := dzpk.poker.Deal(count)

	cg := new(poker.CardGroup)
	size := dzpk.publicCards.Len()
	size1 := size + count
	cg.Cards = dzpk.savePublic.Cards[size:size1]
	if cg == nil {
		mlog.Error("publicDeal error,cards is nil,deal count=%d.", count)
		return -1
	}
	dzpk.publicCards.AddCardGroup(cg)
	dzpk.GameLog.AddGameStep("发公共牌，游戏阶段[%s],本次新牌[%s],当前公牌[%s]", getState(dzpk.gameState), cg.String(), dzpk.publicCards.String())
	//广播
	BroadCastDealPublicCards(dzpk, cg)
	dzpk.updateCardType()
	return 0
}

//更新彩池信息
func (dzpk *DzpkLogic) updatePotData() {
	mlog.Debug("更新彩池信息")
	dzpk.potArray = calPot(dzpk.gameUserData)
	dzpk.GameLog.AddGameStep("更新奖池信息,游戏阶段[%s]：%s", getState(dzpk.gameState), dzpk.potArray.String())
	BroadCastPotMsg(dzpk)
}

//更新桌子数据
func (dzpk *DzpkLogic) updateTableData() {

	//清除小局数据
	dzpk.resetRoundData()
	//更新牌桌数据
	dzpk.updateRoleSeatNo(false)
}

//更新最大下注筹码
func (dzpk *DzpkLogic) updateCurBetMaxJetton() int64 {

	for _, user := range dzpk.gameUserData {
		if dzpk.betCurrMax < user.currBet {
			dzpk.betCurrMax = user.currBet
		}
	}
	jetton := dzpk.betCurrMax
	mlog.Debug("更新最大下注筹码，jetton=%d", jetton)
	dzpk.GameLog.AddGameStep("当前下注最大筹码，游戏阶段[%s],筹码[%d]", getState(dzpk.gameState), jetton)
	return jetton
}

// 是否结束游戏
func (dzpk *DzpkLogic) isEndGame() bool {
	return false
}

//小局结算
func (dzpk *DzpkLogic) settlementRound() {
	mlog.Debug("settlementRound,开始小局结算")
	BroadCastHandCards(dzpk)
	settlement(dzpk)
	BroadCastRoundSettleMent(dzpk)
}

// 获取玩家角色信息
func (dzpk *DzpkLogic) getRole() (s, b, g int32) {
	var seatSmallBlind, seatBigBlind, seatGun int32 = 0, 0, 0
	defer func() {
		mlog.Debug("seatSmallBlind[%d], seatBigBlind[%d], seatGun[%d]", seatSmallBlind, seatBigBlind, seatGun)
	}()
	for i := 0; i < 20; i++ {
		ud := dzpk.getUserInfoBySeatNo(dzpk.buttonSeat)
		if ud == nil {
			dzpk.buttonSeat++
			if dzpk.buttonSeat > 9 {
				dzpk.buttonSeat = 0
			}
			continue
		}
		break
	}
	seatSmallBlind = dzpk.buttonSeat + 1
	for i := 0; i < 20; i++ {
		ud := dzpk.getUserInfoBySeatNo(seatSmallBlind)
		if ud == nil {
			seatSmallBlind++
			if seatSmallBlind > 9 {
				seatSmallBlind = 0
			}
			continue
		}
		break
	}
	seatBigBlind = seatSmallBlind + 1
	for i := 0; i < 20; i++ {
		ud := dzpk.getUserInfoBySeatNo(seatBigBlind)
		if ud == nil {
			seatBigBlind++
			if seatBigBlind > 9 {
				seatBigBlind = 0
			}
			continue
		}
		break
	}
	seatGun = seatBigBlind + 1
	for i := 0; i < 20; i++ {
		ud := dzpk.getUserInfoBySeatNo(seatGun)
		if ud == nil {
			seatGun++
			if seatGun > 9 {
				seatGun = 0
			}
			continue
		}
		break
	}
	return seatSmallBlind, seatBigBlind, seatGun
}

//更新角色位置信息
func (dzpk *DzpkLogic) updateRoleSeatNo(bIsFirstUpdate bool) {
	var seatButton, seatSmallBlind, seatBigBlind, seatGun int32 = 0, 0, 0, 0
	//设置button
	seatButton = dzpk.buttonSeat
	//bButton, bSmall, bBig, bGun := false, false, false, false
	seatSmallBlind, seatBigBlind, seatGun = dzpk.getRole()
	mlog.Debug("seatSmallBlind=%d,seatSmallBlind=%d,seatGun=%d", seatSmallBlind, seatSmallBlind, seatGun)
	//两人
	if dzpk.gameUserCount == 2 {
		seatBigBlind = seatSmallBlind
		seatSmallBlind = seatButton
		dzpk.nGun = seatButton
	} else if dzpk.gameUserCount == 3 {
		dzpk.nGun = seatButton
	} else {
		dzpk.nGun = seatGun
	}
	mlog.Debug("角色信息座位号，Button=%d,SmallBlind=%d,BigBlind=%d", seatButton, seatSmallBlind, seatBigBlind)
	dzpk.GameLog.AddGameStep("角色信息座位号，庄家座位号[%d],小盲座位号[%d],大盲座位号[%d]", seatButton, seatSmallBlind, seatBigBlind)
	userButton := dzpk.getUserInfoBySeatNo(seatButton)
	if userButton != nil {
		userButton.userRole = ROLE_BUTTON
	}
	userSmallBlind := dzpk.getUserInfoBySeatNo(seatSmallBlind)
	if userSmallBlind != nil {
		if seatSmallBlind == seatButton {
			userSmallBlind.userRole = ROLE_BUTTONSMALL
		} else {
			userSmallBlind.userRole = ROLE_SMALLBLIND
		}

	}
	userBigBlind := dzpk.getUserInfoBySeatNo(seatBigBlind)
	if userBigBlind != nil {
		userBigBlind.userRole = ROLE_BIGBLIND
	}
}

//进入房间
func (dzpk *DzpkLogic) UserEnter(gameuser *gamelogic.GameUser, tableno int32, seatno int32) {
	mlog.Debug("dzpk,进入房间,SeatNo=%d,userId=%d", gameuser.SeatNo, gameuser.UserID)
	dzpk.DynamicGameLogicBase.UserEnter(gameuser, tableno, seatno)

	//更新用户钱币
	//dzpk.updateUserScore()
	//发送游戏设置
	sendGameSettings(dzpk, gameuser)
	//如果非游戏玩家，发送当前场景
	if dzpk.IsGaming() && !dzpk.IsGameUser(gameuser) {
		mlog.Debug("UserEnter,userId=%d is not game user,send cur sence.", gameuser.UserID)
		sendGameScene(dzpk, gameuser)
	}
}

// 更新玩家分数
func (dzpk *DzpkLogic) updateUserScore() {
	for _, ud := range dzpk.gameUserData {
		ud.userMoney = ud.User.Score
	}
}

//用户断线重回
func (dzpk *DzpkLogic) UserRecover(gameuser *gamelogic.GameUser) {
	mlog.Debug("UserRecover，用户断线重回")
	dzpk.DynamicGameLogicBase.UserRecover(gameuser)

	//发送游戏设置
	sendGameSettings(dzpk, gameuser)

	//发送游戏场景
	if dzpk.IsGaming() {
		sendGameScene(dzpk, gameuser)
	}
}

//获取对应角色初始化筹码
func (dzpk *DzpkLogic) getInitJetton(ud *DzpkUserData) int64 {
	role := ud.getUserRole()
	if role == ROLE_SMALLBLIND || role == ROLE_BUTTONSMALL {
		return dzpk.betBase / 2
	}
	if role == ROLE_BIGBLIND {
		return dzpk.betBase
	}
	return 0
}

//设置用户数据
func (dzpk *DzpkLogic) setGameUserData() {

	dzpk.gameUserData = make([]*DzpkUserData, 0)
	dzpk.gameUserCount = 0
	for _, ur := range dzpk.AllTableUser() {
		userdata := NewDzpkUserData(ur)
		dzpk.gameUserData = append(dzpk.gameUserData, userdata)
		dzpk.gameUserCount++
	}
	mlog.Debug("设置玩家数据，玩家人数：%d人。", dzpk.gameUserCount)
}

// 结束游戏
func (dzpk *DzpkLogic) endGame() {
	mlog.Debug("游戏结束")
	dzpk.GameEndFunc()
}

//预操作
func (dzpk *DzpkLogic) preOpration() {
	mlog.Debug("preOpration begin.")
	//跟新预操作信息
	for _, ud := range dzpk.gameUserData {
		if ud.seatNo == dzpk.currActionSeat {
			continue
		}
		if ud.isFold || ud.isAllin {
			continue
		}
		ud.updatePreOpration(dzpk)
	}
	//广播出去
	BroadCastPreOpration(dzpk)
}

// 重置玩家操作状态
func (dzpk *DzpkLogic) resetUserOpState() {
	for _, ud := range dzpk.gameUserData {
		ud.curTurnBet = 0
		if ud.isAllin || ud.isFold {
			continue
		}
		ud.resetCurOpration()
	}
}

// 设置测试命令
func (dzpk *DzpkLogic) setTestCommond(gmsg *msg.Message) (nRet int32) {
	defer func() {
		mlog.Debug("setTestCommond return[%d]", nRet)
	}()
	if gmsg.SClassID == int32(netproto.DZPK_GameMessageClassID_ActionGMCommond) {
		//解析数据
		data := gmsg.MsgData.([]byte)
		ust := &netproto.DZPK_GM{}
		err := proto.Unmarshal(data, ust)
		//错误校验
		if err != nil {
			mlog.Error("解析数据错误%v", err)
			return ERROR_OPRATE
		}

		str := *ust.Commond
		str = strings.TrimLeft(str, " ")
		str = strings.TrimRight(str, " ")
		commond := strings.Split(str, " ")
		if len(commond) < 2 {
			return ERROR_TESTDATAERROR
		}
		return dzpk.commondInfo(commond)
	}
	return 1
}

func (dzpk *DzpkLogic) commondInfo(commond []string) int32 {
	/*
		if len(commond) < 2 {
			mlog.Error("commondInfo,commond size<2.")
			return ERROR_TESTDATAERROR
		}
		mlog.Debug("commond type[%s]", commond[0])
		if commond[0] == "SetOutTime" {
			nTime, err := strconv.Atoi(commond[1])
			if err != nil {
				mlog.Error("commondInfo,SetOutTime error.")
				return ERROR_TESTDATAERROR
			}
			if nTime > 5 {
				dzpk.optTime = int32(nTime)
			}
			return 0
		}
		randomPoker := new(poker.RandomPoker)
		for index, tmp := range commond[1:] {
			if tmp == "" || tmp == " " {
				continue
			}
			cd, err := strconv.Atoi(tmp)
			if err != nil {
				mlog.Error("error,parameter[%d] is not number[%s].", index, tmp)
				return ERROR_TESTDATAERROR
			}
			if cd > 150 || cd < 20 {
				mlog.Error("error,cd > 150 || cd < 20.")
				return ERROR_TESTDATAERROR
			}
			card := new(poker.PokerCard)
			card.Color = int32(cd) % 10
			card.CardNumber = int32(cd) / 10
			randomPoker.Cards = append(randomPoker.Cards, card)

		}
		if commond[0] == "ChangeAllCards" {
			if dzpk.publicCards.Len() < 5 {
				mlog.Error("commondInfo,ChangeAllCards,publicCards size<5.")
				return ERROR_TESTDATAERROR
			}
			if len(randomPoker.Cards) < int(dzpk.gameUserCount)*2+5 {
				mlog.Error("commondInfo,ChangeAllCards parameter number not enough.")
				return ERROR_TESTDATAERROR
			}
			dzpk.publicCards = randomPoker.Deal(dzpk.publicCards.Len())
			for _, ud := range dzpk.gameUserData {
				ud.userCards = randomPoker.Deal(ud.userCards.Len())
			}
			dzpk.updateCardType()
			for _, ud := range dzpk.gameUserData {
				sendGameScene(dzpk, ud.User)
			}
			return 0
		}
		if commond[0] == "ChangePubCards" {
			if len(randomPoker.Cards) < dzpk.publicCards.Len() {
				mlog.Error("commondInfo,ChangePubCards parameter size < public cards size.")
				return ERROR_TESTDATAERROR
			}
			dzpk.publicCards = randomPoker.Deal(dzpk.publicCards.Len())
			dzpk.updateCardType()
			for _, ud := range dzpk.gameUserData {
				sendGameScene(dzpk, ud.User)
			}
			return 0
		}
		if commond[0] == "ChangeUserCards" {
			if len(randomPoker.Cards) < int(dzpk.gameUserCount)*2 {
				mlog.Error("commondInfo,ChangeUserCards parameter size is not enough.")
				return ERROR_TESTDATAERROR
			}
			for _, ud := range dzpk.gameUserData {
				ud.userCards = randomPoker.Deal(ud.userCards.Len())
			}
			dzpk.updateCardType()
			for _, ud := range dzpk.gameUserData {
				sendGameScene(dzpk, ud.User)
			}
			return 0
		}
	*/
	return 0
}

//换最大牌
func (dzpk *DzpkLogic) ExchangeMaxCard(userid int32) {

	dzpk.LogicBase.ExchangeMaxCard(userid)
	var nType, nTypeMax, maxRetUserID int32 = 0, -1, 0
	fpoker := dzpk.poker.(*poker.FixedPoker)
	for _, ur := range dzpk.AllGameUser() {
		if ur.IsRobot() || ur.UserID == userid {
			continue
		}
		cg := fpoker.GetUserCards(ur.UserID)
		cp := new(poker.CardGroup)
		cp.AddCardGroup(cg)
		cp.AddCardGroup(dzpk.savePublic)
		nType = GetCardType(cp)
		//

		//比较牌型
		if nType > nTypeMax {
			nTypeMax = nType
			maxRetUserID = ur.UserID
		}
	}
	for _, ur := range dzpk.AllGameUser() {
		if ur.IsRobot() || ur.UserID == userid {
			continue
		}

		cgmax := fpoker.GetUserCards(maxRetUserID)
		cpmax := new(poker.CardGroup)
		cpmax.AddCardGroup(cgmax)
		cpmax.AddCardGroup(dzpk.savePublic)

		cgur := fpoker.GetUserCards(ur.UserID)
		cpur := new(poker.CardGroup)
		cpur.AddCardGroup(cgur)
		cpur.AddCardGroup(dzpk.savePublic)

		nRet := ComPareCards(cpur, cpmax)
		if nRet > 0 {
			maxRetUserID = ur.UserID
		}

	}

	//该牌小于机器人牌才换
	cgmax := fpoker.GetUserCards(maxRetUserID)
	cpmax := new(poker.CardGroup)
	cpmax.AddCardGroup(cgmax)
	cpmax.AddCardGroup(dzpk.savePublic)

	cgur := fpoker.GetUserCards(userid)
	cpur := new(poker.CardGroup)
	cpur.AddCardGroup(cgur)
	cpur.AddCardGroup(dzpk.savePublic)

	nRet := ComPareCards(cpur, cpmax)
	if nRet >= 0 {
		return
	}

	dzpk.GameLog.AddGameStep("交换用户%d和用户%d牌型", userid, maxRetUserID)
	fpoker.ExChange(userid, maxRetUserID)

}
