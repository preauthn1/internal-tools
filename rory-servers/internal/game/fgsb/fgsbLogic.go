package fgsb

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

// Logic 富貴三寶主要結構
type Logic struct {
	*gamelogic.DynamicGameLogicBase
	poker            poker.Poker   // 遊戲牌組
	dealerHandCards  *Poker        // 莊家手牌
	dealerIsBad      bool          // 莊家不成局
	user             []*User       // 玩家資料
	guest            []*Guest      // 觀戰者資料
	gameState        int           // 遊戲狀態
	nowActionName    string        // 遊戲狀態名稱
	round            int           // 回合數
	manual           []string      // 遊戲說明
	timeout          map[int]int32 // 狀態運行時間
	chip             []int64       // 籌碼面額
	sideSheet        map[int]int64 // 邊注表
	specialSheet     map[int]int64 // 特殊牌型表
	strSideSheet     []string      // 邊注表說明
	strSpecialSheet  []string      // 特殊牌型表說明
	BetLimit         int64         // 底注上限
	BetSideLimit     int64         // 邊注上限
	kickGuestTime    int64         // 多久沒下注將玩家踢出遊戲，單位秒
	kickGuestTipTime int64         // 距離踢出房剩餘10秒先提示訊息
	kickUserRound    int           // 多少回合沒下注將玩家踢出遊戲，單位回合
}

// NewLogic 創建富貴三寶邏輯
func NewLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, configs *config.Group) *Logic {
	glc := new(Logic)
	glc.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, configs)

	glc.GetGameConfig(configs)

	glc.poker = nil
	glc.dealerHandCards = NewPoker()
	glc.dealerIsBad = false
	glc.user = make([]*User, glc.TableManager.Room.SeatCount)
	glc.guest = make([]*Guest, glc.TableManager.Room.SeatCount)
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		glc.user[i] = NewUser()
		glc.guest[i] = NewGuest()
	}

	glc.gameState = StateWaitingStart
	glc.nowActionName = ActionName[StateWaitingStart]
	glc.round = 1

	// 遊戲說明
	glc.manual = append(glc.manual, "一副标准的52张牌中，每名玩家发到3张牌，之后就不再发任何牌了。这3张牌都是暗牌（面朝下）。")
	glc.manual = append(glc.manual, "三张牌扑克中，顺子比同花大。所有花色不分大小，牌型排序依次为：A（最大）、K、Q、J、10、9、8、7、6、5、4、3 和 2。")
	glc.manual = append(glc.manual, "同花A-K-Q最大，有时也称为皇家同花顺。2-3-4最小。如果闲家和庄家牌型相同，则打成平手")
	glc.manual = append(glc.manual, "对牌以上下注：")
	glc.manual = append(glc.manual, "只要玩家的牌大于等于一对，玩家就可赢钱。如果玩家的牌不成对，则玩家会输掉赌注。")
	glc.manual = append(glc.manual, "底注红利：")
	glc.manual = append(glc.manual, "无论庄家成手与否，无论闲家牌大还是牌小，始终会支付底注红利")
	glc.manual = append(glc.manual, "顺子（底注 1 赔 1）、三条（底注 1 赔 4）和同花顺（底注 1 赔 5）")

	//----- 邊注 賠率表 -----
	// 同花順
	str := fmt.Sprintf("%s,1赔%d", PatternString(PatternStraightFlush),
		glc.sideSheet[PatternStraightFlush])
	glc.strSideSheet = append(glc.strSideSheet, str)

	// 三條
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternThreeOfAKind),
		glc.sideSheet[PatternThreeOfAKind])
	glc.strSideSheet = append(glc.strSideSheet, str)

	// 順子
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternStraight),
		glc.sideSheet[PatternStraight])
	glc.strSideSheet = append(glc.strSideSheet, str)

	// 同花
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternFlush),
		glc.sideSheet[PatternFlush])
	glc.strSideSheet = append(glc.strSideSheet, str)

	// 對子
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternPair),
		glc.sideSheet[PatternPair])
	glc.strSideSheet = append(glc.strSideSheet, str)

	//----- 特殊牌型 賠率表 -----
	// 同花順
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternStraightFlush),
		glc.specialSheet[PatternStraightFlush])
	glc.strSpecialSheet = append(glc.strSpecialSheet, str)

	// 三條
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternThreeOfAKind),
		glc.specialSheet[PatternThreeOfAKind])
	glc.strSpecialSheet = append(glc.strSpecialSheet, str)

	// 順子
	str = fmt.Sprintf("%s,1赔%d", PatternString(PatternStraight),
		glc.specialSheet[PatternStraight])
	glc.strSpecialSheet = append(glc.strSpecialSheet, str)

	tfManager.NewMilliSecTicker(1000, func() {
		glc.KickGuest()
	})

	glc.RegisterBeforeGameStartFunc(func() {
		//fmt.Printf("BeforeGameStart\n")
		glc.SendWatchUserCount()
		glc.InitData()
	})
	glc.RegisterAfterGameStartFunc(func() {
		//fmt.Printf("AfterGameStart\n")
		glc.Betting(false, nil)
	})
	glc.RegisterBeforeGameEndFunc(func() {
		//fmt.Printf("BeforeGameEnd\n")
	})
	glc.RegisterAfterGameEndFunc(func() {
		//fmt.Printf("AfterGameEnd\n")

		for _, ur := range glc.AllTableUser() {
			// 用戶狀態變為準備狀態
			ur.UserState.OnUserReady(ur, glc.TableManager.Room)
		}

		if !glc.IsGaming() && glc.GetUserCount()+glc.GetGuestCount() > 0 {
			glc.GameStartFunc()
		}
	})

	return glc
}

/*interface*/

// UserEnter d用戶進入
func (glc *Logic) UserEnter(gameUser *gamelogic.GameUser, tableno int32, seatno int32) {
	//fmt.Printf("UserEnter %d\n", gameUser.UserID)

	if gameUser == nil {
		return
	}

	mlog.Debug("TableNo[%d], UserEnter Guest, User[%d]", glc.TableNo(), gameUser.UserID)
	// 改變為GameState 避免再次抓去進入房間
	gameUser.UserState = new(gamelogic.UserGameState)
	// 改變桌號為目前桌號 避免斷線無法呼叫Leave
	gameUser.TableNo = glc.TableNo()

	glc.AddGuest(gameUser)

	if !glc.IsGaming() && glc.GetUserCount()+glc.GetGuestCount() > 0 {
		glc.GameStartFunc()
	}

	//fmt.Printf("UserEnter OK %d\n", gameUser.UserID)
}

// UserFree 用戶站立
func (glc *Logic) UserFree(gameUser *gamelogic.GameUser) {
	//fmt.Printf("UserFree %d\n", gameuser.UserID)
	mlog.Debug("UserFree%d", gameUser.UserID)

	if gameUser == nil {
		return
	}

	guestSeatNo := int32(gamelogic.FreeSeatNo)
	if gameUser.SeatNo == gamelogic.FreeSeatNo {
		guestSeatNo = glc.GetGuestSeat(gameUser.UserID)
	}

	mlog.Debug("Table[%d], UserFree, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	if guestSeatNo != gamelogic.FreeSeatNo {
		if guestSeatNo >= 0 && guestSeatNo < glc.TableManager.Room.SeatCount {
			glc.guest[guestSeatNo].Init()
		}
	} else if gameUser.SeatNo != gamelogic.FreeSeatNo {
		user := glc.user[gameUser.SeatNo]

		if user.totalBet > 0 && glc.IsGaming() {
			return
		}

		gameUser.UserState.OnGameEnd(gameUser, glc.TableManager.Room)

		if gameUser.SeatNo >= 0 && gameUser.SeatNo < glc.TableManager.Room.SeatCount {
			glc.user[int(gameUser.SeatNo)].Init()
		}
	}

	if !gameUser.IsRobot() {
		mlog.System("UserFree %s", gameUser)
	}

	glc.SendAllUserInfo()
	glc.SendWatchUserCount()

	glc.LogicBase.UserFree(gameUser)
	//fmt.Printf("UserFree OK %d\n", gameuser.UserID)
}

// UserRecover 用户斷線重連
func (glc *Logic) UserRecover(gameUser *gamelogic.GameUser) {
	//fmt.Printf("UserRecover %d\n", gameUser.UserID)
	mlog.Debug("UserRecover %d", gameUser.UserID)

	if gameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserRecover, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	glc.DynamicGameLogicBase.UserRecover(gameUser)

	glc.SendGameInfo(gameUser)
	glc.SendAllUserInfo()
	glc.SendWatchUserCount()
	//fmt.Printf("UserRecover OK %d\n", gameUser.UserID)
}

// UserLeave 用戶離開
func (glc *Logic) UserLeave(gameUser *gamelogic.GameUser) bool {
	//fmt.Printf("UserLeave %d\n", gameUser.UserID)
	mlog.Debug("UserLeave %d", gameUser.UserID)
	if gameUser == nil {
		return false
	}

	guestSeatNo := int32(gamelogic.FreeSeatNo)
	if gameUser.SeatNo == gamelogic.FreeSeatNo {
		guestSeatNo = glc.GetGuestSeat(gameUser.UserID)
	}

	mlog.Debug("Table[%d], UserLeave, User[%d], Seat[%d], IsOffLine()[%t]", glc.TableNo(), gameUser.UserID, gameUser.SeatNo, gameUser.IsOffLine())

	if guestSeatNo != gamelogic.FreeSeatNo {
		gameUser.UserState.OnGameEnd(gameUser, glc.TableManager.Room)
		if guestSeatNo >= 0 && guestSeatNo < glc.TableManager.Room.SeatCount {
			glc.guest[guestSeatNo].Init()
		}
	} else if gameUser.SeatNo != gamelogic.FreeSeatNo {
		pUser := glc.user[gameUser.SeatNo]

		if pUser.totalBet > 0 && glc.IsGaming() {
			return false
		}

		gameUser.UserState.OnGameEnd(gameUser, glc.TableManager.Room)

		if gameUser.SeatNo >= 0 && gameUser.SeatNo < glc.TableManager.Room.SeatCount {
			glc.user[int(gameUser.SeatNo)].Init()
		}
	}

	glc.SendAllUserInfo()
	glc.SendWatchUserCount()
	//fmt.Printf("UserLeave OK %d\n", gameUser.UserID)
	return glc.DynamicGameLogicBase.UserLeave(gameUser)
}

// GameMessage 遊戲封包
func (glc *Logic) GameMessage(gameUser *gamelogic.GameUser, msg *msg.Message) {
	if gameUser == nil || msg == nil {
		return
	}

	if msg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	glc.CheckUserAnimationCommited(gameUser, msg, glc.nowActionName, nil)

	user := glc.GetUserData(gameUser.UserID)

	// C2S 玩家選座位
	if user == nil && msg.SClassID == int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat) {
		glc.MessageSelectSeat(gameUser, msg)
		return
	}

	if user == nil {
		return
	}

	switch msg.SClassID {
	case int32(netproto.FGSB_GameMessageClassID_FGSB_Bet): // C2S 玩家下注，包含底注及邊注
		glc.MessageBet(user, msg)
	case int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet): // C2S 玩家加注
		fallthrough
	case int32(netproto.FGSB_GameMessageClassID_FGSB_FoldCard): // C2S 棄牌
		glc.MessageDoubleOrFold(user, msg)
	}
}

// GetFreeSeat 獲取一個空閒的座位
func (glc *Logic) GetFreeSeat() (int32, int32) {

	if int32(glc.GetUserCount()+glc.GetGuestCount()) >= glc.TableManager.Room.SeatCount {
		return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
	}

	st := glc.TableManager.GetFreeSeat()

	if st != nil {
		return glc.TableNo(), st.SeatNo
	}

	return gamelogic.FreeTableNo, gamelogic.FreeSeatNo
}

// TableRealUserCount 遊戲中真實用户數
func (glc *Logic) TableRealUserCount() int {
	return glc.GetRealUserCount() + glc.GetRealGuestCount()
}

// BoardcastMessage 廣播封包給所有人
func (glc *Logic) BoardcastMessage(m *msg.Message) {
	glc.DynamicGameLogicBase.BoardcastMessage(m)

	for _, u := range glc.guest {
		if u.user != nil {
			glc.SendMessgeFunc(u.userID, m)
		}
	}
}

/*Interface*/

/*GameLogic*/

/**遊戲邏輯**/

// GetGameConfig 讀取遊戲設定
func (glc *Logic) GetGameConfig(configs *config.Group) {
	mlog.Debug("GetGameConfig")
	if configs == nil {
		mlog.Error("桌號[%d] 讀取 config 失敗", glc.TableNo())
		return
	}

	var str string = ""

	// 讀取 gameserver 設定值
	cfgs := configs.GetConfig(gamelogic.ConfigFileNameGameServer)
	if cfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		//----- timeout -----
		glc.timeout = make(map[int]int32)
		glc.timeout[StateWaitingStart] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateWaitingStart])
		glc.timeout[StateBetting] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateBetting])
		glc.timeout[StateStopBetting] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateStopBetting])
		glc.timeout[StateDeal] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateDeal])
		glc.timeout[StateBetSideSettlement] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateBetSideSettlement])
		glc.timeout[StateBetDouble] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateBetDouble])
		glc.timeout[StateSettlement] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateSettlement])
		glc.timeout[StateGameOver] = cfgs.GetConfigInt32(ConfigLabel, TimeOut[StateGameOver])

		// 籌碼面額
		str = cfgs.GetConfigValue("GameSetting", "ChipValue")
		strChip := strings.Split(str, ",")
		for i := 0; i < len(strChip); i++ {
			i64chip, _ := strconv.ParseInt(strChip[i], 10, 32)
			glc.chip = append(glc.chip, i64chip)
		}

		//----- 特殊邊注 -----
		glc.sideSheet = make(map[int]int64)

		// 皇家同花順，賠率
		glc.sideSheet[PatternRoyalFlush] = cfgs.GetConfigInt64(ConfigLabel, "Side_StraightFlush")
		// 同花順，賠率
		glc.sideSheet[PatternStraightFlush] = cfgs.GetConfigInt64(ConfigLabel, "Side_StraightFlush")
		// 三條，賠率
		glc.sideSheet[PatternThreeOfAKind] = cfgs.GetConfigInt64(ConfigLabel, "Side_ThreeOfAKind")
		// 順子，賠率
		glc.sideSheet[PatternStraight] = cfgs.GetConfigInt64(ConfigLabel, "Side_Straight")
		// 同花，賠率
		glc.sideSheet[PatternFlush] = cfgs.GetConfigInt64(ConfigLabel, "Side_Flush")
		// 對子，賠率
		glc.sideSheet[PatternPair] = cfgs.GetConfigInt64(ConfigLabel, "Side_Pair")

		//----- 特殊牌型 ----
		glc.specialSheet = make(map[int]int64)

		// 皇家同花順，賠率
		glc.specialSheet[PatternRoyalFlush] = cfgs.GetConfigInt64(ConfigLabel, "Special_StraightFlush")
		// 同花順，賠率
		glc.specialSheet[PatternStraightFlush] = cfgs.GetConfigInt64(ConfigLabel, "Special_StraightFlush")
		// 三條，賠率
		glc.specialSheet[PatternThreeOfAKind] = cfgs.GetConfigInt64(ConfigLabel, "Special_ThreeOfKind")
		// 順子，賠率
		glc.specialSheet[PatternStraight] = cfgs.GetConfigInt64(ConfigLabel, "Special_Straight")

		// 底注上限
		glc.BetLimit = cfgs.GetConfigInt64(ConfigLabel, "BetLimit")

		// 邊注上限
		glc.BetSideLimit = cfgs.GetConfigInt64(ConfigLabel, "BetSideLimit")

		// 踢人參數
		glc.kickGuestTime = int64(cfgs.GetConfigInt(ConfigLabel, "KickGuestTime")) * int64(time.Second)
		glc.kickGuestTipTime = int64(cfgs.GetConfigInt(ConfigLabel, "KickGuestTipTime")) * int64(time.Second)
		glc.kickUserRound = cfgs.GetConfigInt(ConfigLabel, "KickUserRound")
	}
}

// InitData 初始化遊戲參數
func (glc *Logic) InitData() {
	mlog.Debug("InitData")

	if glc.GetServerConfig().PokerType == gamelogic.PokerTypeRandom {
		glc.poker = poker.NewRandomPoker(HasKing, HasPizi, PokerSetCount)
		glc.poker.Shuffle()
	} else { //發牌器發牌
		glc.poker = poker.NewFilePoker(glc.Configs, glc.AllTableUser())
		glc.poker = poker.NewFixedFilePoker(glc.poker.(*poker.FilePoker))
		//glc.poker = poker.NewFixedFileFillPoker(glc.poker.(*poker.FilePoker), glc.AllTableUser())
	}

	glc.dealerHandCards.CardGroup = poker.NewEmptyCardGroup(0)
	glc.dealerIsBad = false
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		gameUser := glc.GetUserBySeatNo(i)
		if gameUser == nil {
			glc.user[i].Init()
		} else {
			glc.user[i].ResetGameData()
			glc.user[i].state = UserStateNormal
		}
	}

	glc.gameState = StateWaitingStart
	glc.nowActionName = ActionName[StateWaitingStart]
}

// AddGuest 加入觀戰
func (glc *Logic) AddGuest(gameUser *gamelogic.GameUser) {
	if gameUser == nil {
		return
	}

	seatNo := int(gameUser.SeatNo)
	if seatNo != gamelogic.FreeSeatNo {
		return
	}

	seat := glc.GetFreeGuestSeat()
	guest := glc.guest[seat]
	if guest.user == nil {
		guest.Init()
		guest.user = gameUser
		guest.userID = gameUser.UserID
	}

	glc.SendGameInfo(gameUser)
	glc.SendAllUserInfo()
	glc.SendWatchUserCount()
	example.invalid("加入新玩家[%d], 玩家總人數[%d]", guest.userID, glc.GetUserCount()+glc.GetGuestCount())
	glc.GameLog.AddGameStep("玩家[%d] 旁觀", guest.userID)
}

// AddUser 加入新玩家
func (glc *Logic) AddUser(gameUser *gamelogic.GameUser) {
	if gameUser == nil {
		return
	}

	seatNo := gameUser.SeatNo
	if seatNo < 0 || seatNo >= glc.TableManager.Room.SeatCount {
		return
	}

	guestSeat := glc.GetGuestSeat(gameUser.UserID)
	if guestSeat < 0 || guestSeat >= glc.TableManager.Room.SeatCount {
		return
	}

	glc.guest[guestSeat].Init()

	user := glc.user[seatNo]
	if user.state == UserStateEmptySeat {
		user.Init()
		user.gameUser = gameUser
		user.money = gameUser.Score
		user.ResetGameData()

		if glc.gameState >= StateStopBetting && glc.gameState < StateGameOver {
			user.state = UserStateWatch
		} else {
			user.state = UserStateNormal
		}
	} else {
		mlog.Error("AddUser is not empty, %+v", user)
	}

	//glc.sendGameInfo(gameUser)
	glc.SendAllUserInfo()
	glc.SendWatchUserCount()
	example.invalid("玩家入座[%d], 遊戲狀態[%s], 座位[%d], 類別[%d], 玩家狀態[%s], 玩家總人數[%d]",
		user.gameUser.UserID,
		StrGameStateDebug[glc.gameState], seatNo,
		user.gameUser.UserType,
		StrUserStateDebug[user.state],
		glc.TableRealUserCount())

	glc.GameLog.AddGameStep("玩家[%d] 入座，狀態：%s", gameUser.UserID, StrUserStateDebug[user.state])
}

// SendGameInfo 傳送牌局資訊
func (glc *Logic) SendGameInfo(gameUser *gamelogic.GameUser) {
	if gameUser == nil {
		return
	}

	mlog.Debug("sendGameInfo, GameState[%s], DealerCard[%s]",
		StrGameStateDebug[glc.gameState], glc.dealerHandCards.String())

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_GameInfo)

	data := new(netproto.FGSB_Msg_GameInfo)

	// 邊注表說明
	for i := 0; i < len(glc.strSideSheet); i++ {
		data.SideSheet = append(data.SideSheet, glc.strSideSheet[i])
	}

	// 特殊牌型表說明
	for i := 0; i < len(glc.strSpecialSheet); i++ {
		data.SpecialSheet = append(data.SpecialSheet, glc.strSpecialSheet[i])
	}

	// 遊戲說明
	for i := 0; i < len(glc.manual); i++ {
		data.Manual = append(data.Manual, glc.manual[i])
	}

	// 所有玩家資訊
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		userInfo := glc.UserInfoPkg(i)
		data.UserInfo = append(data.UserInfo, userInfo)
	}

	// 籌碼面額
	for i := 0; i < len(glc.chip); i++ {
		data.ChipValue = append(data.ChipValue, glc.chip[i])
	}

	// 莊家手牌
	dealerCard := new(netproto.FGSB_Msg_UserCards)

	for idx := 0; idx < glc.dealerHandCards.Len(); idx++ {
		card := new(netproto.FGSB_Msg_PokerCard)
		if glc.gameState >= StateSettlement {
			sCard := glc.dealerHandCards.Cards[idx]
			card.Color = proto.Int32(sCard.Color)
			card.CardNumber = proto.Int32(sCard.CardNumber)
		} else {
			card.Color = proto.Int32(poker.PokerColorUnKnown)
			card.CardNumber = proto.Int32(0)
		}
		dealerCard.Cards = append(dealerCard.Cards, card)
	}
	dealerCard.Pattern = proto.Int32(int32(glc.dealerHandCards.Pattern(true)))
	data.BankerCards = dealerCard

	data.Bet_Limit = proto.Int64(glc.BetLimit)
	data.Bet_SideLimit = proto.Int64(glc.BetSideLimit)
	data.GameState = proto.Int32(int32(glc.gameState))
	data.TimeRemain = proto.Int32(glc.timeout[glc.gameState] - glc.PlayerActionManager.GetActionTimerElapsedSeconds(glc.nowActionName))
	data.Round = proto.Int32(int32(glc.round))
	data.SelfInfo = glc.SelfInfoPkg(gameUser)
	data.TableNo = proto.Int32(glc.TableNo())

	message.MsgData = data
	glc.SendMessageToUser(gameUser, message)

	switch glc.gameState {
	case StateBetting: // 押底注
	case StateStopBetting: // 停止下注
	case StateDeal: // 發牌
	case StateBetSideSettlement: // 邊注、特殊牌型結算
	case StateBetDouble: // 加注
	case StateSettlement: // 底注結算
	}
}

// SendAllUserInfo 傳送玩家資訊
func (glc *Logic) SendAllUserInfo() {
	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo)

	data := new(netproto.FGSB_Msg_AllUserInfo)

	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		pUserInfo := glc.UserInfoPkg(seat)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	message.MsgData = data

	glc.BoardcastMessage(message)
}

// SendWatchUserCount 觀戰人數
func (glc *Logic) SendWatchUserCount() {
	var count = int32(0)
	for _, user := range glc.guest {
		if user.user != nil {
			count++
		}
	}

	msg := new(msg.Message)
	msg.BClassID = int32(netproto.MessageBClassID_Game)
	msg.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_WatchUserCountID)
	data := &netproto.FGSB_Msg_WatchUserCount{}
	data.Count = proto.Int32(count)
	msg.MsgData = data
	glc.BoardcastMessage(msg)
}

// UserInfoPkg 玩家資訊打包資料，seat：此座位的玩家資訊
func (glc *Logic) UserInfoPkg(seat int32) *netproto.FGSB_Msg_UserInfo {
	if seat < 0 || seat >= glc.TableManager.Room.SeatCount {
		return nil
	}

	user := glc.user[seat]

	userInfo := new(netproto.FGSB_Msg_UserInfo)

	if user.gameUser != nil {
		userInfo.UserID = proto.Int32(user.gameUser.UserID) // 玩家ID，若是空位時，玩家ID = -1
		userInfo.FaceID = proto.Int32(user.gameUser.FaceID) // 玩家頭像ID
	} else {
		userInfo.UserID = proto.Int32(-1)
		userInfo.FaceID = proto.Int32(0) // 玩家頭像ID
	}
	userInfo.State = proto.Int32(int32(user.state))                  // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
	userInfo.SeatNo = proto.Int32(seat + int32(ClientStartUpSeatNo)) // 玩家座位，client 座位號從 1 開始
	userInfo.Money = proto.Int64(user.money)                         // 玩家金錢
	userInfo.Bet_Base = proto.Int64(user.betBase)                    // 玩家底注
	userInfo.Bet_Side = proto.Int64(user.betSide)                    // 玩家邊注
	userInfo.Bet_Double = proto.Int64(user.betDouble)                // 玩家加注

	// 玩家手牌
	pHandCard := new(netproto.FGSB_Msg_UserCards)
	for idx := 0; idx < glc.user[seat].cards.Len(); idx++ {
		card := new(netproto.FGSB_Msg_PokerCard)
		sCard := glc.user[seat].cards.Cards[idx]
		card.Color = proto.Int32(sCard.Color)
		card.CardNumber = proto.Int32(sCard.CardNumber)

		pHandCard.Cards = append(pHandCard.Cards, card)
	}
	pHandCard.Pattern = proto.Int32(int32(glc.user[seat].cards.Pattern(false)))
	userInfo.UserCards = pHandCard

	userInfo.Winner = proto.Int32(int32(user.winState))    // 贏家
	userInfo.Win_Side = proto.Int64(user.winSide)          // 邊注輸贏金錢，若值為負數，則為玩家輸掉的金額
	userInfo.Win_Special = proto.Int64(user.winSpecial)    // 特殊牌型輸贏金錢，若值為負數，則為玩家輸掉的金額
	userInfo.Win_Base = proto.Int64(user.winBase)          // 底注輸贏金錢，若值為負數，則為玩家輸掉的金額
	userInfo.Win_Double = proto.Int64(user.winDouble)      // 加注輸贏金錢，若值為負數，則為玩家輸掉的金額
	userInfo.TotalWinLose = proto.Int64(user.totalWinLose) // 總輸贏，client 端顯示用

	return userInfo
}

// SelfInfoPkg 玩家自己資料
func (glc *Logic) SelfInfoPkg(gameUser *gamelogic.GameUser) *netproto.FGSB_Msg_UserInfo {

	guestSeat := glc.GetGuestSeat(gameUser.UserID)
	seat := glc.GetUserSeat(gameUser.UserID)
	userInfo := new(netproto.FGSB_Msg_UserInfo)

	if guestSeat >= 0 || guestSeat < glc.TableManager.Room.SeatCount {
		userInfo.UserID = proto.Int32(gameUser.UserID) // 玩家ID，若是空位時，玩家ID = -1
		userInfo.FaceID = proto.Int32(gameUser.FaceID) // 玩家頭像ID
		userInfo.State = proto.Int32(-1)               // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
		userInfo.SeatNo = proto.Int32(-1)              // 玩家座位，client 座位號從 1 開始
		userInfo.Money = proto.Int64(gameUser.Score)   // 玩家金錢
		userInfo.Bet_Base = proto.Int64(0)             // 玩家底注
		userInfo.Bet_Side = proto.Int64(0)             // 玩家邊注
		userInfo.Bet_Double = proto.Int64(0)           // 玩家加注
		userInfo.UserCards = new(netproto.FGSB_Msg_UserCards)
		userInfo.UserCards.Pattern = proto.Int32(0)
		userInfo.Winner = proto.Int32(int32(WinStateNone)) // 贏家
		userInfo.Win_Side = proto.Int64(0)                 // 邊注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Special = proto.Int64(0)              // 特殊牌型輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Base = proto.Int64(0)                 // 底注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Double = proto.Int64(0)               // 加注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.TotalWinLose = proto.Int64(0)             // 總輸贏，client 端顯示用

		return userInfo
	} else if seat >= 0 || seat < glc.TableManager.Room.SeatCount {
		user := glc.user[seat]

		if user.gameUser != nil {
			userInfo.UserID = proto.Int32(user.gameUser.UserID) // 玩家ID，若是空位時，玩家ID = -1
			userInfo.FaceID = proto.Int32(user.gameUser.FaceID) // 玩家頭像ID
		} else {
			userInfo.UserID = proto.Int32(-1)
			userInfo.FaceID = proto.Int32(0) // 玩家頭像ID
		}
		userInfo.State = proto.Int32(int32(user.state))                  // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
		userInfo.SeatNo = proto.Int32(seat + int32(ClientStartUpSeatNo)) // 玩家座位，client 座位號從 1 開始
		userInfo.Money = proto.Int64(user.money)                         // 玩家金錢
		userInfo.Bet_Base = proto.Int64(user.betBase)                    // 玩家底注
		userInfo.Bet_Side = proto.Int64(user.betSide)                    // 玩家邊注
		userInfo.Bet_Double = proto.Int64(user.betDouble)                // 玩家加注

		// 玩家手牌
		pHandCard := new(netproto.FGSB_Msg_UserCards)
		for idx := 0; idx < glc.user[seat].cards.Len(); idx++ {
			card := new(netproto.FGSB_Msg_PokerCard)
			sCard := glc.user[seat].cards.Cards[idx]
			card.Color = proto.Int32(sCard.Color)
			card.CardNumber = proto.Int32(sCard.CardNumber)

			pHandCard.Cards = append(pHandCard.Cards, card)
		}
		pHandCard.Pattern = proto.Int32(int32(glc.user[seat].cards.Pattern(false)))
		userInfo.UserCards = pHandCard

		userInfo.Winner = proto.Int32(int32(user.winState))    // 贏家
		userInfo.Win_Side = proto.Int64(user.winSide)          // 邊注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Special = proto.Int64(user.winSpecial)    // 特殊牌型輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Base = proto.Int64(user.winBase)          // 底注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.Win_Double = proto.Int64(user.winDouble)      // 加注輸贏金錢，若值為負數，則為玩家輸掉的金額
		userInfo.TotalWinLose = proto.Int64(user.totalWinLose) // 總輸贏，client 端顯示用

		return userInfo
	}

	return nil
}

// CheckGetSide 檢查是否獲得邊注金額，回傳值［所中邊注金額］
func (glc *Logic) CheckGetSide(seatNo int32) int64 {
	if seatNo < 0 || seatNo >= glc.TableManager.Room.SeatCount {
		return 0
	}

	if glc.user[seatNo].betSide <= 0 {
		return 0
	}

	if glc.user[seatNo].state != UserStateNormal {
		return 0
	}

	sideMoney := int64(0)

	pattern := glc.user[seatNo].cards.Pattern(false)
	if rate, ok := glc.sideSheet[pattern]; ok {
		sideMoney = rate * glc.user[seatNo].betSide
	}

	return sideMoney
}

// CheckGetSpecial 檢查是否獲得特殊牌型金額，回傳值［所中特殊牌型金額］
func (glc *Logic) CheckGetSpecial(seatNo int32) int64 {
	if seatNo < 0 || seatNo >= glc.TableManager.Room.SeatCount {
		return 0
	}

	if glc.user[seatNo].betBase <= 0 {
		return 0
	}

	if glc.user[seatNo].state != UserStateNormal {
		return 0
	}

	specialMoney := int64(0)

	pattern := glc.user[seatNo].cards.Pattern(false)
	if rate, ok := glc.specialSheet[pattern]; ok {
		specialMoney = rate * glc.user[seatNo].betBase
	}

	return specialMoney
}

// IsNoBet 確認是否無人下注
func (glc *Logic) IsNoBet() bool {
	TotalBetBase := int64(0)
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]

		if user == nil {
			continue
		}

		if user.state == UserStateNormal {
			TotalBetBase += user.betBase
		}
	}
	if TotalBetBase > 0 {
		return false
	}

	return true
}

// KickGuest 踢觀戰
func (glc *Logic) KickGuest() {
	startTime := time.Now().UnixNano()
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		gameUser := glc.guest[seat].user
		if gameUser != nil && !gameUser.IsRobot() {
			if !glc.guest[seat].kickTip && startTime-glc.guest[seat].enterTime > int64(glc.kickGuestTime)-int64(glc.kickGuestTipTime) {
				glc.guest[seat].kickTip = true
				t := glc.kickGuestTipTime / int64(time.Second)
				//踢人-未操作時間到 提示剩餘時間
				glc.SendMessageToUser(gameUser, glc.KickTipMessage(int32(t), 0))
			}
			if startTime-glc.guest[seat].enterTime > int64(glc.kickGuestTime) {
				glc.KickUserFunc(gameUser.UserID, "您即将返回大厅")
			}
		}
	}
}

// KickTipMessage 踢觀戰提示訊息
func (glc *Logic) KickTipMessage(kickTime, kickRound int32) *msg.Message {
	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_KickTip)

	data := new(netproto.FGSB_Msg_KickTip)
	data.Time = proto.Int32(kickTime)
	data.Round = proto.Int32(kickRound)
	message.MsgData = data

	return message
}

/**遊戲流程**/

// Betting 開始下注
func (glc *Logic) Betting(forUser bool, gameUser *gamelogic.GameUser) {
	//fmt.Printf("Betting\n")
	if !forUser {
		mlog.Debug("Betting")
		glc.GameLog.AddGameStep("開始下注")
		glc.gameState = StateBetting
		glc.nowActionName = glc.PlayerActionManager.NewAction(ActionName[glc.gameState], 0, glc.timeout[glc.gameState], func() {
			glc.StopBetting(false, nil)
		})

		glc.SendAllUserInfo()
	}

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart)
	pData := new(netproto.FGSB_Msg_UserBettingStart)
	pData.TimeRemain = proto.Int32(glc.timeout[glc.gameState] - glc.PlayerActionManager.GetActionTimerElapsedSeconds(glc.nowActionName))
	pData.Round = proto.Int32(int32(glc.round))

	message.MsgData = pData
	if forUser {
		if gameUser != nil {
			glc.SendMessageToUser(gameUser, message)
		}
	} else {
		glc.BoardcastMessage(message)
	}
}

// StopBetting 停止下注
func (glc *Logic) StopBetting(forUser bool, gameUser *gamelogic.GameUser) {
	//fmt.Printf("StopBetting\n")
	if !forUser {
		mlog.Debug("StopBetting")
		glc.GameLog.AddGameStep("停止下注")
		glc.gameState = StateStopBetting
		// 無人下注，重頭開始
		if glc.IsNoBet() {
			glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
				glc.GameOver()
			})
		} else {
			glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
				glc.Deal()
			})
		}

		for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
			user := glc.user[seat]
			if user.state == UserStateNormal {
				if user.totalBet <= 0 {
					user.state = UserStateWatch
					user.noBetRound++ // 踢人-未操作回合增加
				} else {
					user.noBetRound = 0 // 踢人-有操作回合重置
				}

				if user.noBetRound >= glc.kickUserRound && !user.gameUser.IsRobot() {
					// 踢人-未操作回合到
					glc.KickUserFunc(user.gameUser.UserID, "您即将返回大厅")
				} else if user.noBetRound >= glc.kickUserRound-1 && !user.gameUser.IsRobot() {
					// 踢人-未操作回合到 提示剩一回合
					glc.SendMessageToUser(user.gameUser, glc.KickTipMessage(0, 1))
				}

			}
		}
	}

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_StopBettingAnimation)

	data := new(netproto.FGSB_Msg_StopBettingAnimation)

	// 動畫名稱
	data.ActionName = proto.String(ActionName[glc.gameState])
	message.MsgData = data

	if forUser {
		if gameUser != nil {
			glc.SendMessageToUser(gameUser, message)
		}
	} else {
		glc.BoardcastMessage(message)
	}
}

// Deal 發牌
func (glc *Logic) Deal() {
	//fmt.Printf("Deal\n")
	mlog.Debug("Deal")
	glc.GameLog.AddGameStep("開始發牌")
	glc.gameState = StateDeal
	glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
		glc.BetSideSettlement()
	})

	// 開始發玩家牌
	for _, user := range glc.user {
		if user.state != UserStateNormal {
			continue
		}

		if user.gameUser == nil {
			continue
		}

		cards := glc.poker.Deal(user.gameUser.UserID, HandCardCount)
		user.cards.SetCards(cards)
		mlog.Debug("玩家[%d], %s, 牌型[%s]",
			user.gameUser.UserID,
			cards.String(),
			user.cards.PatternString(false))

		glc.GameLog.AddGameStep("玩家[%d], %s, 牌型[%s]",
			user.gameUser.UserID,
			cards.String(),
			user.cards.PatternString(false))
	}

	// 發莊家牌
	cards := glc.poker.Deal(0, HandCardCount)
	glc.dealerHandCards.SetCards(cards)
	mlog.Debug("莊家, %s, 牌型[%s]",
		cards.String(),
		glc.dealerHandCards.PatternString(true))

	glc.GameLog.AddGameStep("莊家, %s, 牌型[%s]",
		cards.String(),
		glc.dealerHandCards.PatternString(true))

	// 莊家牌是否成局
	if glc.dealerHandCards.Pattern(true) == PatternNone {
		glc.dealerIsBad = true
	} else {
		glc.dealerIsBad = false
	}

	// 傳送訊息給 client
	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Deal)

	data := new(netproto.FGSB_Msg_Deal)

	// 莊家手牌
	dealerCard := new(netproto.FGSB_Msg_UserCards)
	for idx := 0; idx < glc.dealerHandCards.Len(); idx++ {
		card := new(netproto.FGSB_Msg_PokerCard)
		card.Color = proto.Int32(poker.PokerColorUnKnown)
		card.CardNumber = proto.Int32(0)
		dealerCard.Cards = append(dealerCard.Cards, card)
	}
	dealerCard.Pattern = proto.Int32(int32(glc.dealerHandCards.Pattern(true)))
	data.BankerCards = dealerCard

	// 玩家手牌
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		pUserInfo := glc.UserInfoPkg(seat)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	// 動畫名稱
	data.ActionName = proto.String(ActionName[glc.gameState])

	message.MsgData = data

	glc.BoardcastMessage(message)
}

// BetSideSettlement 結算邊注、特殊牌型
func (glc *Logic) BetSideSettlement() {
	//fmt.Printf("BetSideSettlement\n")
	mlog.Debug("BetSideSettlement")
	glc.GameLog.AddGameStep("結算邊注，特殊牌型")
	glc.gameState = StateBetSideSettlement
	glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
		glc.BetDouble()
	})

	// 結算邊注
	glc.GameLog.AddGameStep("開始結算邊注")
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		if user.state == UserStateNormal {
			mlog.Debug("玩家[%d] 原来金錢 %s", user.gameUser.UserID, global.MoneyString(user.money))

			sideMoney := glc.CheckGetSide(seat)
			sideMoney, tax := glc.GetUserWinAndTax(sideMoney)

			if sideMoney > 0 {
				sideMoney += user.betSide
			}

			moneyChange := sideMoney - user.betSide
			user.winSide = sideMoney + tax - user.betSide
			user.totalWinLose += sideMoney

			// 暫存起來第二階段再結算
			//user.money += sideMoney
			user.tempSideMoney = sideMoney
			user.tempTax += tax
			user.moneyChange += moneyChange

			// 暫存起來第二階段再結算
			// 寫分
			//if user.gameUser.UserType != gamelogic.UserType_Rebot {
			//	glc.UpdateUserScoreFunc(user.gameUser.UserID, moneyChange, tax)
			//}

			str := fmt.Sprintf("玩家[%d] 邊注 %s, %s[%s], 中邊注 %s, 税收 %s, 輸贏 %s, 結算輸後金錢 %s",
				user.gameUser.UserID,
				global.MoneyString(user.betSide),
				glc.user[seat].cards.String(),
				PatternString(glc.user[seat].cards.Pattern(false)),
				global.MoneyString(sideMoney+tax),
				global.MoneyString(tax),
				global.MoneyString(user.winSide),
				global.MoneyString(user.money))

			mlog.Debug("%s", str)
			glc.GameLog.AddGameStep("%s", str)
		}
	}

	// 結算特殊牌型
	glc.GameLog.AddGameStep("開始結算特殊牌型")
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		if user.state == UserStateNormal {
			mlog.Debug("玩家[%d] 原来金錢 %s", user.gameUser.UserID, global.MoneyString(user.money))

			specialMoney := glc.CheckGetSpecial(seat)
			specialMoney, tax := glc.GetUserWinAndTax(specialMoney)

			user.winSpecial = specialMoney + tax
			user.totalWinLose += specialMoney

			// 暫存起來第二階段再結算
			//user.money += specialMoney
			user.tempSpecialMoney = specialMoney
			user.tempTax += tax
			user.moneyChange += specialMoney

			// 暫存起來第二階段再結算
			// 寫分
			//if user.gameUser.UserType != gamelogic.UserType_Rebot {
			//	glc.UpdateUserScoreFunc(user.gameUser.UserID, specialMoney, tax)
			//}

			str := fmt.Sprintf("玩家[%d] 底注 %s, %s[%s], 中特殊牌型 %s, 税收 %s, 輸贏 %s, 結算輸後金錢 %s",
				user.gameUser.UserID,
				global.MoneyString(user.betBase),
				glc.user[seat].cards.String(),
				PatternString(glc.user[seat].cards.Pattern(false)),
				global.MoneyString(specialMoney+tax),
				global.MoneyString(tax),
				global.MoneyString(user.winSide),
				global.MoneyString(user.money))

			mlog.Debug("%s", str)
			glc.GameLog.AddGameStep("%s", str)
		}
	}

	// 邊注結算流水
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]

		if user.tempSideMoney+user.tempSpecialMoney == 0 {
			user.tempFlow = user.betSide
		}
	}

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SideSettlement)

	data := new(netproto.FGSB_Msg_SideSettlement)

	// 玩家手牌
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		pUserInfo := glc.UserInfoPkg(seat)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	// 動畫名稱
	data.ActionName = proto.String(ActionName[glc.gameState])

	message.MsgData = data

	glc.BoardcastMessage(message)
}

// BetDouble 玩家開始加注
func (glc *Logic) BetDouble() {
	//fmt.Printf("BetDouble\n")
	glc.GameLog.AddGameStep("玩家開始加注")
	glc.gameState = StateBetDouble
	// 這邊要取用的是當局玩家的數量, by @Rory
	glc.nowActionName = glc.PlayerActionManager.NewAction(ActionName[glc.gameState], int32(glc.GetNormalUserCount()), glc.timeout[glc.gameState], func() {
		glc.Settlement()
	})

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart)

	pData := new(netproto.FGSB_Msg_DoubleBetStart)
	pData.TimeRemain = proto.Int32(glc.timeout[glc.gameState] - glc.PlayerActionManager.GetActionTimerElapsedSeconds(glc.nowActionName))

	pMsg.MsgData = pData
	glc.BoardcastMessage(pMsg)
}

// Settlement 總結算
func (glc *Logic) Settlement() {
	//fmt.Printf("Settlement\n")
	mlog.Debug("Settlement")
	glc.GameLog.AddGameStep("總結算")
	glc.gameState = StateSettlement
	glc.nowActionName = glc.RegisterAnimationByUserCount(ActionName[glc.gameState], int32(glc.TableRealUserCount()), glc.timeout[glc.gameState], func() {
		glc.GameOver()
	})

	// 未加注玩家棄牌
	for iSeat := int32(0); iSeat < glc.TableManager.Room.SeatCount; iSeat++ {
		user := glc.user[iSeat]
		if user.state == UserStateNormal && user.betDouble <= 0 {
			user.state = UserStateFold
			glc.SendDoubleOrFoldResult(user, DoubleBetResultFold)
		}
	}

	mlog.Debug("莊家：%s[%s]", glc.dealerHandCards.String(), glc.dealerHandCards.PatternString(true))
	glc.GameLog.AddGameStep("莊家：%s[%s]", glc.dealerHandCards.String(), glc.dealerHandCards.PatternString(true))

	winLose := int64(0)
	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		switch user.state {
		case UserStateNormal: // 正常
			// 莊家不成局
			if glc.dealerIsBad {
				winLose = user.betBase*int64(1+BetBasePayRate) + user.betDouble
				//user.winBase = winLose - user.betBase - user.betDouble
				user.winBase = winLose - user.betBase - user.betDouble
				user.totalWinLose += winLose
				user.money += winLose
				user.winState = WinStateWin
			} else {
				UserCards := glc.user[seat].cards
				compareResult := glc.dealerHandCards.Compare(UserCards, false, false)
				switch compareResult {
				case CompareEven: // 平手
					winLose = user.betBase + user.betDouble
					user.winState = WinStateEven
				case CompareDealer: // 莊家贏
					winLose = -user.betBase - user.betDouble
					user.winState = WinStateLose
				case ComparePlayer: // 玩家贏
					winLose = (user.betBase + user.betDouble) * int64(1+BetBasePayRate)
					user.winState = WinStateWin
				}

				if compareResult == pokerset.POKERCompareDealer {
					user.winBase = winLose
					//user.winBase = 0
				} else {
					user.winBase = winLose - user.betBase - user.betDouble
					//user.winBase = winLose
					user.money += winLose
					user.totalWinLose += winLose
				}
			}

			// 寫分
			if winLose > 0 {
				winLose -= user.betBase + user.betDouble
			}
			moneyChange, tax := glc.GetUserWinAndTax(winLose)
			//user.winBase -= tax
			user.totalWinLose -= tax
			user.money -= tax
			user.moneyChange += moneyChange
			user.tempTax += tax

			// 取回第一階段暫存
			user.money += user.tempSideMoney + user.tempSpecialMoney

			flow := int64(0)
			if moneyChange < 0 {
				flow = user.betBase + user.betDouble
			} else {
				flow = 0
			}
			user.tempFlow += flow

			//if user.gameUser.UserType != gamelogic.UserType_Rebot {
			//	glc.UpdateUserScoreFunc(user.gameUser.UserID, moneyChange, tax)
			//}

			//fmt.Printf("更新前 玩家[%d] money[%d] Score[%d] moneyChange[%d]\n", user.gameUser.UserID, user.money, user.gameUser.Score, user.moneyChange)

			// 寫回總變動值
			//if user.gameUser.UserType != gamelogic.UserType_Rebot {
			mlog.Debug("結算寫分")
			glc.UpdateUserScoreFunc(user.gameUser.UserID, user.moneyChange, user.tempTax)
			//}

			//fmt.Printf("更新後 玩家[%d] money[%d] Score[%d]\n", user.gameUser.UserID, user.money, user.gameUser.Score)

			mlog.Debug("玩家[%d], %s[%s], 底注 %s, 加注 %s, 輸贏 %s, 税收 %s, 流水 %s, 結算後金錢 %s",
				user.gameUser.UserID,
				glc.user[seat].cards.String(),
				glc.user[seat].cards.PatternString(false),
				global.MoneyString(user.betBase),
				global.MoneyString(user.betDouble),
				global.MoneyString(winLose),
				global.MoneyString(tax),
				global.MoneyString(flow),
				global.MoneyString(user.money))
			glc.GameLog.AddGameStep("玩家[%d], %s[%s], 底注 %s, 加注 %s, 輸贏 %s, 税收 %s, 流水 %s, 結算輸後金錢 %s",
				user.gameUser.UserID,
				glc.user[seat].cards.String(),
				glc.user[seat].cards.PatternString(false),
				global.MoneyString(user.betBase),
				global.MoneyString(user.betDouble),
				global.MoneyString(winLose),
				global.MoneyString(tax),
				global.MoneyString(flow),
				global.MoneyString(user.money))
		case UserStateFold: // 棄牌
			winLose = -user.betBase
			user.winBase = winLose
			//user.winBase = 0
			user.moneyChange += winLose

			// 取回第一階段暫存
			user.money += user.tempSideMoney + user.tempSpecialMoney

			//glc.UpdateUserScoreFunc(user.gameUser.UserID, winLose, 0)
			// 寫回總變動值
			//if user.gameUser.UserType != gamelogic.UserType_Rebot {
			mlog.Debug("棄牌寫分")
			glc.UpdateUserScoreFunc(user.gameUser.UserID, user.moneyChange, user.tempTax)
			//}

			user.winState = WinStateLose

			flow := int64(0)
			if winLose < 0 {
				flow = user.betBase
			} else {
				flow = 0
			}
			user.tempFlow += flow

			mlog.Debug("玩家[%d] 棄牌, %s[%s], 底注 %s, 加注 %s, 輸贏 %s, 流水 %s, 結算輸後金錢 %s",
				user.gameUser.UserID,
				glc.user[seat].cards.String(),
				glc.user[seat].cards.PatternString(false),
				global.MoneyString(user.betBase),
				global.MoneyString(user.betDouble),
				global.MoneyString(winLose),
				global.MoneyString(flow),
				global.MoneyString(user.money))
			glc.GameLog.AddGameStep("玩家[%d] 棄牌, %s[%s], 底注 %s, 加注 %s, 輸贏 %s, 流水 %s, 結算輸後金錢 %s",
				user.gameUser.UserID,
				glc.user[seat].cards.String(),
				glc.user[seat].cards.PatternString(false),
				global.MoneyString(user.betBase),
				global.MoneyString(user.betDouble),
				global.MoneyString(winLose),
				global.MoneyString(flow),
				global.MoneyString(user.money))
		}
	}

	//----- 傳送訊息給 client -----

	message := new(msg.Message)
	message.BClassID = int32(netproto.MessageBClassID_Game)
	message.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Settlement)

	data := new(netproto.FGSB_Msg_Settlement)
	data.ActionName = proto.String(ActionName[glc.gameState])

	// 莊家手牌
	dealerCard := new(netproto.FGSB_Msg_UserCards)
	for idx := 0; idx < glc.dealerHandCards.Len(); idx++ {
		card := new(netproto.FGSB_Msg_PokerCard)
		sCard := glc.dealerHandCards.Cards[idx]
		card.Color = proto.Int32(sCard.Color)
		card.CardNumber = proto.Int32(sCard.CardNumber)

		dealerCard.Cards = append(dealerCard.Cards, card)
	}
	dealerCard.Pattern = proto.Int32(int32(glc.dealerHandCards.Pattern(true)))
	data.BankerCards = dealerCard

	// 所有玩家資訊
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		pUserInfo := glc.UserInfoPkg(i)
		data.UserInfo = append(data.UserInfo, pUserInfo)
	}

	message.MsgData = data
	glc.BoardcastMessage(message)
}

// GameOver 遊戲結束
func (glc *Logic) GameOver() {
	//fmt.Printf("GameOver\n")
	mlog.Debug("GameOver")
	glc.GameLog.AddGameStep("遊戲結束")

	if glc.gameState == StateSettlement {
		glc.round++
	}

	glc.gameState = StateGameOver

	for seat := int32(0); seat < glc.TableManager.Room.SeatCount; seat++ {
		user := glc.user[seat]
		if user.gameUser != nil {
			var flowAmount, betAmount, payAmount int64
			betAmount = user.totalBet
			payAmount = betAmount + user.moneyChange
			glc.GameLog.AddGameStep("用戶[%d] 結算流水 %d, 有效投注 %d, 返還金額 %d", user.gameUser.UserID, flowAmount, betAmount, payAmount)
			user.gameUser.FlowAmount += user.tempFlow
			user.gameUser.BetAmount += betAmount
			user.gameUser.PayAmount += payAmount

			glc.GameLog.Context().UserUpdateScore(user.gameUser.UserID, user.moneyChange)
			glc.GameLog.Context().UserBetScore(user.gameUser.UserID, betAmount)
			glc.GameLog.Context().UserPayoutScore(user.gameUser.UserID, payAmount+user.tempTax)
			glc.GameLog.Context().UserValidScore(user.gameUser.UserID, user.tempFlow)
			glc.GameLog.Context().UserTaxAmount(user.gameUser.UserID, user.tempTax)

		}
	}

	glc.GameLog.GameEnd()

	glc.GameEndFunc()
}

/**封包處理**/

// MessageSelectSeat 訊息處理 - 玩家選位
func (glc *Logic) MessageSelectSeat(gameUser *gamelogic.GameUser, message *msg.Message) {
	if gameUser == nil || message == nil {
		return
	}

	guestSeat := glc.GetGuestSeat(gameUser.UserID)
	if guestSeat == gamelogic.FreeSeatNo {
		return
	}

	msgInfo := &netproto.DCHJK_Msg_SelectSeat{}
	proto.Unmarshal(message.MsgData.([]byte), msgInfo)

	// Clent端座位從1開始
	iSeat := msgInfo.GetSeatNo() - int32(ClientStartUpSeatNo)

	iResult := glc.UserActionSelectSeat(gameUser, iSeat)

	user := glc.GetUserData(gameUser.UserID)
	glc.SendSelectSeatResult(user, iSeat, iResult)
}

// UserActionSelectSeat 玩家選位
func (glc *Logic) UserActionSelectSeat(gameUser *gamelogic.GameUser, seat int32) int {
	if gameUser == nil {
		return SelectSeatResultSelectFailed
	}

	st := glc.TableManager.GetSeatByNo(seat)
	if st == nil {
		return SelectSeatResultNoSeat
	}

	if st.SeatUser != nil {
		return SelectSeatResultHaveUser
	}
	gameUser.UserState = new(gamelogic.UserQueueState)
	glc.DynamicGameLogicBase.UserEnter(gameUser, glc.TableNo(), seat)
	glc.DynamicGameLogicBase.UserReady(gameUser)
	if glc.IsGaming() {
		gameUser.UserState.OnUserGame(gameUser, glc.TableManager.Room)
	}
	glc.AddUser(gameUser)

	return SelectSeatResultSuccess
}

// SendSelectSeatResult 傳送玩家選位結果
func (glc *Logic) SendSelectSeatResult(user *User, seat int32, result int) {
	if user == nil {
		mlog.Debug("sendSelectSeatResult(), 無玩家資料")
		return
	}

	mlog.Debug("sendSelectSeatResult(), User[%d], SeatNo[%d], i64Bet[%d], iResult[%d]", user.gameUser.UserID, user.gameUser.SeatNo, seat, result)

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat_OK)

	pData := new(netproto.FGSB_Msg_SelectSeat_OK)
	pData.SeatNo = proto.Int32(seat + int32(ClientStartUpSeatNo)) // 玩家座位
	pData.Result = proto.Int32(int32(result))                     // 結果

	pMsg.MsgData = pData
	glc.SendMessageToUser(user.gameUser, pMsg)

	mlog.Debug("玩家[%d] 選位 %d, 結果：%s", user.gameUser.UserID, seat, StrSelectSeatResult[result])
	glc.GameLog.AddGameStep("玩家[%d] 選位 %d, 結果：%s", user.gameUser.UserID, seat, StrSelectSeatResult[result])
}

// MessageBet 訊息處理 - 玩家下注
func (glc *Logic) MessageBet(user *User, message *msg.Message) {
	if user == nil || message == nil {
		return
	}

	msgInfo := &netproto.FGSB_Msg_Bet{}
	proto.Unmarshal(message.MsgData.([]byte), msgInfo)
	betType := msgInfo.GetBetType()
	bet := msgInfo.GetBet()

	example.invalid("玩家[%d] 下%s %s", user.gameUser.UserID, StrBetType[betType], global.MoneyString(bet))
	glc.GameLog.AddGameStep("玩家[%d] 下%s %s", user.gameUser.UserID, StrBetType[betType], global.MoneyString(bet))

	if glc.gameState != StateBetting {
		glc.SendBetResult(user, betType, bet, BetResultBettingStopped)
		return
	}

	result := glc.UserActionBet(user, betType, bet)
	glc.SendBetResult(user, betType, bet, result)
}

// UserActionBet 玩家下注
func (glc *Logic) UserActionBet(user *User, betType int32, bet int64) int {
	if user == nil {
		return BetResultBetFailed
	}

	result := BetResultSuccess

	if bet <= 0 {
		return BetResultBetSmallThanZero
	}

	// 剩餘金錢是否足以下注
	if user.money < bet {
		return BetResultNotEnoughMoney
	}

	switch betType {
	case BetTypeBase: // 底注
		// 下注是否超過上限
		if user.betBase+bet > (user.money+user.betBase)/2 {
			return BetResultBaseBetOverUpperLimit
		}

		// 下注是否超過Config上限
		if user.betBase+bet > glc.BetLimit && glc.BetLimit > 0 {
			return BetResultBaseBetOverUpperLimit
		}

		user.betBase += bet
		user.totalBet += bet
		user.totalWinLose -= bet
	case BetTypeSide: // 邊注
		// 是否下過底注
		if user.betBase <= 0 {
			return BetResultNeverBetBase
		}

		// 邊注是否超過上限
		if user.money-bet < user.betBase {
			return BetResultBaseBetOverUpperLimit
		}

		// 下注是否超過Config上限
		if user.betSide+bet > glc.BetSideLimit && glc.BetSideLimit > 0 {
			return BetResultBaseBetOverUpperLimit
		}

		user.betSide += bet
		user.totalBet += bet
		user.totalWinLose -= bet
	default:
		return BetResultBetFailed
	}

	user.money -= bet
	return result
}

// SendBetResult 傳送玩家下注結果
func (glc *Logic) SendBetResult(user *User, betType int32, bet int64, result int) {
	if user == nil {
		mlog.Debug("sendBetResult(), 無玩家資料")
		return
	}

	mlog.Debug("sendBetResult(), User[%d], SeatNo[%d], bet[%d], result[%d]", user.gameUser.UserID, user.gameUser.SeatNo, bet, result)

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_BetResult)

	pData := new(netproto.FGSB_Msg_BetResult)
	pData.SeatNo = proto.Int32(user.gameUser.SeatNo + int32(ClientStartUpSeatNo)) // 下注玩家座位
	pData.BetType = proto.Int32(betType)                                          // 下注類別
	pData.Bet = proto.Int64(bet)                                                  // 下注金額
	pData.MoneyRemain = proto.Int64(user.money)                                   // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(result))                                     // 下注結果

	pMsg.MsgData = pData
	if result == BetResultSuccess {
		glc.BoardcastMessage(pMsg)
	} else {
		glc.SendMessageToUser(user.gameUser, pMsg)
	}

	mlog.Debug("玩家[%d] 下注%s %s, 結果：%s, 底注額 %s, 邊注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		StrBetType[betType],
		global.MoneyString(bet),
		StrBetResult[result],
		global.MoneyString(user.betBase),
		global.MoneyString(user.betSide),
		global.MoneyString(user.money))
	glc.GameLog.AddGameStep("玩家[%d] 下注%s %s, 結果：%s, 底注額 %s, 邊注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		StrBetType[betType],
		global.MoneyString(bet),
		StrBetResult[result],
		global.MoneyString(user.betBase),
		global.MoneyString(user.betSide),
		global.MoneyString(user.money))

}

// MessageDoubleOrFold 訊息處理 - 玩家加注
func (glc *Logic) MessageDoubleOrFold(user *User, message *msg.Message) {
	if user == nil || message == nil || user.state != UserStateNormal {
		return
	}

	switch message.SClassID {
	case int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet): // 加注
		example.invalid("玩家[%d] 加注", user.gameUser.UserID)
		glc.GameLog.AddGameStep("玩家[%d] 加注", user.gameUser.UserID)
	case int32(netproto.FGSB_GameMessageClassID_FGSB_FoldCard): // 棄牌
		example.invalid("玩家[%d] 棄牌", user.gameUser.UserID)
		glc.GameLog.AddGameStep("玩家[%d] 棄牌", user.gameUser.UserID)
	}

	if glc.gameState != StateBetDouble {
		glc.SendDoubleOrFoldResult(user, DoublebetResultBetDoubleStopped)
		return
	}

	result := BetResultBetFailed
	switch message.SClassID {
	case int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet): // 加注
		result = glc.UserActionDoubleBet(user)
	case int32(netproto.FGSB_GameMessageClassID_FGSB_FoldCard): // 棄牌
		result = glc.UserActionFold(user)
	}

	glc.SendDoubleOrFoldResult(user, result)

	if result == DoubleBetResultSuccess || result == DoubleBetResultFold {
		glc.PlayerActionManager.OnUserAction(glc.nowActionName, user.gameUser.UserID)
	}
}

// UserActionDoubleBet 玩家加注
func (glc *Logic) UserActionDoubleBet(user *User) int {
	if user == nil {
		return DoublebetResultBetDoubleFailed
	}

	result := DoubleBetResultSuccess

	if user.betDouble > 0 {
		return DoublebetResultBetDoubleHasDone
	}

	doubleBet := user.betBase

	if user.money < doubleBet {
		glc.UserActionFold(user)
		return DoubleBetResultNotEnoughMoney
	}

	user.betDouble += doubleBet
	user.totalBet += doubleBet
	user.money -= doubleBet
	user.totalWinLose -= doubleBet

	return result
}

// UserActionFold 玩家棄牌
func (glc *Logic) UserActionFold(user *User) int {
	if user == nil {
		return DoublebetResultBetDoubleFailed
	}

	if user.state == UserStateFold {
		return DoubleBetResultFold
	}

	result := DoubleBetResultFold

	user.state = UserStateFold

	return result
}

// SendDoubleOrFoldResult 傳送玩家加注結果
func (glc *Logic) SendDoubleOrFoldResult(user *User, result int) {
	if user == nil {
		mlog.Debug("SendDoubleBetResult(), 無玩家資料")
		return
	}

	mlog.Debug("SendDoubleBetResult(), User[%d], SeatNo[%d], bet[%d], result[%d]", user.gameUser.UserID, user.gameUser.SeatNo, user.betDouble, result)

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBetResult)

	pData := new(netproto.FGSB_Msg_DoubleBetResult)
	pData.SeatNo = proto.Int32(user.gameUser.SeatNo + int32(ClientStartUpSeatNo)) // 下注玩家座位
	pData.MoneyRemain = proto.Int64(user.money)                                   // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(result))                                     // 下注結果
	pData.Bet = proto.Int64(user.betDouble)

	pMsg.MsgData = pData
	if result == DoubleBetResultSuccess || result == DoubleBetResultFold {
		glc.BoardcastMessage(pMsg)
	} else {
		glc.SendMessageToUser(user.gameUser, pMsg)
	}

	mlog.Debug("玩家[%d] 加注, 結果：%s, 加注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		StrDoubleBetResult[result],
		global.MoneyString(user.betDouble),
		global.MoneyString(user.money))
	glc.GameLog.AddGameStep("玩家[%d] 加注, 結果：%s, 加注額 %s, 剩餘金錢 %s",
		user.gameUser.UserID,
		StrDoubleBetResult[result],
		global.MoneyString(user.betDouble),
		global.MoneyString(user.money))

}

/*GameLogic*/
