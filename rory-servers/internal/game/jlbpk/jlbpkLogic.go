package jlbpk

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	Global "example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Logic 遊戲邏輯
type Logic struct {
	*gamelogic.DynamicGameLogicBase
	pPokerSet          *JlbpkPokerSet   // 系統牌組
	iSeatCount         int              // 座位數
	pDealerHandCards   *JlbpkPokerSet   // 莊家手牌
	arypHandCards      []*JlbpkPokerSet // 玩家手牌
	arypUserData       []*JlbpkUser     // 玩家資料
	pGameTicker        *time.Ticker     // 遊戲計時器
	iGameState         int              // 遊戲狀態
	i64TimeDuration    int64            // Ticker 計數值
	i64Timeout         int64            // 各階段 timeout
	i64LastMessageTime int64            // 上一次訊息時間
	mapi64Timeout      map[int]int64    // timeout 時間設定
	arystrJPSheet      []string         // Jackpot 賠率表
	arystrSpecialSheet []string         // 特殊牌賠率表
	arystrPokerManual  []string         // 遊戲說明
	aryi32Chip         []int32          // 籌碼面額
	uiAllPlayerReady   uint             // 各玩家均已 Ready，此值設定後就不再變更
	uiPlayerReady      uint             // 標示各玩家是否 Ready
	iUserCount         int              // 遊戲中玩家人數
	mapi64JPSheet      map[int32]int64  // 特殊邊注表
	mapi64SpecialSheet map[int32]int64  // 賠注表
	iDealType          int              // 發牌器發牌模式
	pTestDealerCards   *JlbpkPokerSet   // 測試用莊家手牌
	arypTestUserCards  []*JlbpkPokerSet // 測試用玩家手牌
	bDealerIsGood      bool             // 莊家不成局
	pWinLose           *SysPoolWinLose  // 本局輸贏結果統計
	bHasVIPUser        bool             // 是否有監控玩家在遊戲中
	bDealerCheat       bool             // 莊家設定目標金額
	//pSuperToolService   *JlbpkSuperToolService // 超端服務
}

func NewJlbpkLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, configs *config.Group) *Logic {
	pGame := new(Logic)
	pGame.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, configs)

	pGame.i64TimeDuration = 0
	tfManager.NewMilliSecTicker(100, func() {
		pGame.i64TimeDuration++
		pGame.gameTickerFunc()
	})

	pGame.setGameState(JLB_State_WaitingPlayer, 0)

	pGame.mapi64Timeout = make(map[int]int64)
	pGame.mapi64JPSheet = make(map[int32]int64)
	pGame.mapi64SpecialSheet = make(map[int32]int64)
	pGame.iDealType = DealType_Natural

	// g_pSysPoolWinLose = new(SysPoolWinLose)
	// g_pSysPoolWinLose.GameDBData.WeightTableNoNow = DealerWeightNo_B

	pGame.getGameConfig(configs)
	//pGame.pSuperToolService = pGame.GetTable().Room.SuperToolService.(*JlbpkSuperToolService)

	pGame.pPokerSet = NewJlbpkPokerSet()
	pGame.pDealerHandCards = NewJlbpkPokerSet()
	pGame.arypHandCards = make([]*JlbpkPokerSet, pGame.iSeatCount)
	pGame.arypUserData = make([]*JlbpkUser, pGame.iSeatCount)
	for i := 0; i < pGame.iSeatCount; i++ {
		pGame.arypHandCards[i] = NewJlbpkPokerSet()
		pGame.arypUserData[i] = NewJlbpkUser()
		pGame.arypUserData[i].Init()
	}

	// 各玩家均已 Ready，此值設定後就不再變更
	uiReady := uint(1)
	for i := 0; i < pGame.iSeatCount; i++ {
		uiReady *= 2
	}
	pGame.uiAllPlayerReady = uiReady - 1

	pGame.iUserCount = 0

	//----- JP 賠率表 -----
	// 皇家同花順
	str := fmt.Sprintf("%s %d%s奖池", pokerset.PatternString(pokerset.POKERPatternRoyalFlush),
		pGame.mapi64JPSheet[pokerset.POKERPatternRoyalFlush], "%")
	pGame.arystrJPSheet = append(pGame.arystrJPSheet, str)
	// 同花順
	str = fmt.Sprintf("%s %d%s奖池", pokerset.PatternString(pokerset.POKERPatternStraightFlush),
		pGame.mapi64JPSheet[pokerset.POKERPatternStraightFlush], "%")
	pGame.arystrJPSheet = append(pGame.arystrJPSheet, str)
	// 四條
	str = fmt.Sprintf("%s %s元", pokerset.PatternString(pokerset.POKERPatternFourOfAKind),
		Global.MoneyString(pGame.mapi64JPSheet[pokerset.POKERPatternFourOfAKind]))
	pGame.arystrJPSheet = append(pGame.arystrJPSheet, str)
	//葫蘆
	str = fmt.Sprintf("%s %s元", pokerset.PatternString(pokerset.POKERPatternFullHouse),
		Global.MoneyString(pGame.mapi64JPSheet[pokerset.POKERPatternFullHouse]))
	pGame.arystrJPSheet = append(pGame.arystrJPSheet, str)
	// 同花
	str = fmt.Sprintf("%s %s元", pokerset.PatternString(pokerset.POKERPatternFlush),
		Global.MoneyString(pGame.mapi64JPSheet[pokerset.POKERPatternFlush]))
	pGame.arystrJPSheet = append(pGame.arystrJPSheet, str)

	// 特殊賠率表
	for i := pokerset.POKERPatternRoyalFlush; i >= pokerset.POKERPatternTwoPairs-1; i-- {
		if i == pokerset.POKERPatternTwoPairs-1 {
			str = fmt.Sprintf("其他 1赔%d", pGame.mapi64SpecialSheet[i])
		} else {
			str = fmt.Sprintf("%s 1赔%d", pokerset.PatternString(i), pGame.mapi64SpecialSheet[i])
		}
		pGame.arystrSpecialSheet = append(pGame.arystrSpecialSheet, str)
	}

	// 遊戲說明
	if g_iHDCMode == HDCModeEnable {
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "游戏使用一副普通52张牌扑克没有王牌，每轮游戏前重新洗牌。")
	} else {
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "加勒比扑克使用一副普通52张牌扑克没有王牌，每轮游戏前重新洗牌。")
	}
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "玩家先下底注。玩家、庄家各发五张牌，玩家五张牌面朝下，庄家四张牌面朝下，一张牌面朝上。然后玩家评估自己手中的牌，并据此作出以下选择。\n")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "1. 弃牌 - 玩家将输掉所下的底注。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "2. 加注 - 玩家必须在注金区中另外再下底注两倍的赌金。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "3. 若身上现金不足以加注(底注的两倍),则无法下彩金注。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "庄家将翻开手中的牌与玩家的牌进行比较。但是庄家手中必须同时有A和K或者更好的牌，才有资格继续玩。如果庄家没有资格继续玩则“不成局”，玩家便可赢得与底注相等的赌金。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "如果庄家有A和K或更好的牌，并且手里的牌面高于玩家的牌面，玩家就会输掉底注和加注赌金。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "如果玩家的牌面高于庄家的牌面，则会赢得与底注相等的赌金以及按照赔率规则计算得出的红利赌金。")
	pGame.arystrPokerManual = append(pGame.arystrPokerManual, "如果客牌与庄牌点数相同，则不考虑花色，双方打成和局。\n")
	if g_iHDCMode == HDCModeEnable {
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "在开始下注的时候，您可以额外下注彩金注。")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "无论庄牌是否成局，累进大奖彩金注都可以获得彩金，若身上现金不足以加注(底注的两倍),则无法下彩金注")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "赔率如下:")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "皇家同花顺 - 累进彩池的100%")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "同花顺 - 累进彩池的 10%")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "四条 - $500 奖金")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "葫芦 - $100 奖金")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "同花 - $50 奖金")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "賠注表如下:")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "皇家同花順 1 賠 50")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "同花順 1 賠 30")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "四條 1賠20")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "葫蘆 1賠9")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "同花 1賠7")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "順子 1賠4")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "三條 1賠3")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "兩對 1賠2")
		pGame.arystrPokerManual = append(pGame.arystrPokerManual, "其他 1賠1")
	}

	pGame.RegisterBeforeGameStartFunc(func() {
		// 初始化遊戲數據, 適用在遊戲開始前只會執行一次, 之後再調用都不會再執行了
		pGame.initialGameDataOnce()
		pGame.pWinLose = new(SysPoolWinLose)
		pGame.initData()
	})
	pGame.RegisterAfterGameStartFunc(func() {
		pGame.setGameState(JLB_State_StartBettingAnimation, pGame.mapi64Timeout[JLB_Timeout_StartBettingAnimation])
		pGame.playStartBettingAnimation()
	})
	pGame.RegisterBeforeGameEndFunc(func() {
	})
	pGame.RegisterAfterGameEndFunc(func() {
	})

	return pGame
}

func (pGame *Logic) initData() {
	mlog.Debug("初始化遊戲資料")
	/**
	主要解決卡桌問題, 需要在用戶都離開的時候初始化遊戲資料
	若初始化過程中表示遊戲中則取消該局遊戲.
	*/
	if pGame.IsGaming() {
		mlog.Debug("initData 遊戲初始化時遊戲狀態為遊戲中")
		pGame.GameCancelFunc("初始化遊戲資料")
	}
	pGame.pPokerSet.NewSet(JLBPK_PokerSetCount, 0)
	pGame.pPokerSet.Shuffle()
	pGame.pDealerHandCards.SetCardSet("")
	pGame.iUserCount = 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.arypHandCards[iSeat].SetCardSet("")

		pGameUser := pGame.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			pGame.arypUserData[iSeat].Init()
		} else {
			if pGame.arypUserData[iSeat].iUserState != JLB_UserState_EmptySeat {
				pGame.arypUserData[iSeat].iUserState = JLB_UserState_Normal
			}
			pGame.UserReady(pGameUser)
			pGame.iUserCount++
		}
	}
	pGame.bDealerIsGood = true
	pGame.pWinLose.Reset()
	pGame.bDealerCheat = false

	if pGame.iUserCount <= 0 {
		pGame.setGameState(JLB_State_WaitingPlayer, 0)

		mlog.Debug("等待玩家进入")
		pGame.GameLog.AddGameStep("等待玩家进入")
	} else {
		pGame.setGameState(JLB_State_WaitingStart, pGame.mapi64Timeout[JLB_Timeout_WaitingStart])
		pGame.waitingGameBegin()
	}
	pGame.i64LastMessageTime = pGame.i64TimeDuration
	mlog.Debug("initData, 玩家人數：%d", pGame.iUserCount)

	pGame.uiPlayerReady = 0
}

func (pGame *Logic) GameMessage(pGameUser *gamelogic.GameUser, pMsg *msg.Message) {
	if pGameUser == nil || pMsg == nil {
		return
	}

	if pMsg.BClassID != int32(netproto.MessageBClassID_Game) {
		return
	}

	pUser := pGame.GetUserDataByUserID(pGameUser.UserID)
	if pUser == nil {
		return
	}

	switch pMsg.SClassID {
	case int32(netproto.JLBPK_GameMessageClassID_StartBettingAnimationOK): // C2S 開始下注動畫播放完成
		pGame.messageStartBettingAnimationOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_Bet): // C2S 玩家下注，包含底注及邊注
		pGame.messageBet(pUser, pMsg)
	case int32(netproto.JLBPK_GameMessageClassID_StopBettingAnimationOK): // C2S 停止下注動畫播放完成
		pGame.messageStopBettingAnimationOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_DealOk): // C2S 發牌動畫播放完成
		pGame.messageDealOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_DoubleBet): // C2S 玩家加注
		pGame.messageDoubleBet(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_FoldCard): // C2S 棄牌
		pGame.messageFold(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_ShowCardsOK): // C2S 開牌，邊注結算效果完成
		pGame.messageShowCardsOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_Compensate_BetBase_OK): // C2S 賠付，底注結算效果完成
		pGame.messageBetBaseSettlementOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_Compensate_DoubleBet_OK): // C2S 賠付，加注結算效果完成
		pGame.messageBetDoubleSettlementOK(pUser)
	case int32(netproto.JLBPK_GameMessageClassID_SettlementOK): // C2S 總結算效果完成
		pGame.messageSettlementOK(pUser)
	}
}

func (pGame *Logic) PiSaiMessage(pUser *gamelogic.GameUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}
}

// 玩家進入遊戲
func (pGame *Logic) UserEnter(pGameUser *gamelogic.GameUser, i32TableNo int32, i32SeatNo int32) {
	if pGameUser == nil {
		return
	}

	if !g_bIsReadJackpotFromDB {
		g_bIsReadJackpotFromDB = true
		pGame.ReadJackpotPool()
	}

	// 自 DB 取得遊戲輸贏結果
	//if g_bEnablePoolControl {
	//	if !g_bGetWinLoseData {
	//		pGame.GetSysPoolWinLose()
	//	}
	//}

	// 先檢查玩家是否有殘留資料在桌內
	for _, u := range pGame.arypUserData {
		if u.User == nil {
			continue
		}
		if u.User.UserID == pGameUser.UserID {
			if !pGame.UserLeave(u.User) {
				return
			}
			// pGame.KickUserFunc(pGameUser.UserID, "Same UserID login room")
		}
	}

	mlog.Debug("UserEnter, User[%d], Seat[%d]", pGameUser.UserID, pGameUser.SeatNo)

	pGame.DynamicGameLogicBase.UserEnter(pGameUser, i32TableNo, i32SeatNo)

	pGame.DynamicGameLogicBase.UserReady(pGameUser)

	pGame.AddUser(pGameUser)
}

// 計時器處理程序
func (pGame *Logic) gameTickerFunc() {
	if pGame.i64Timeout > 0 && pGame.i64TimeDuration-pGame.i64LastMessageTime > pGame.i64Timeout {
		switch pGame.iGameState {
		//case JLB_State_WaitingStart: // 等待遊戲開始
		//	pGame.GameStartFunc()
		//	pGame.setGameState(JLB_State_StartBettingAnimation, pGame.mapi64Timeout[JLB_Timeout_StartBettingAnimation])
		//	pGame.playStartBettingAnimation()
		case JLB_State_StartBettingAnimation: // 開始下注動畫
			pGame.setGameState(JLB_State_StartBetting, pGame.mapi64Timeout[JLB_Timeout_StartBetting])
			pGame.startBetting()
		case JLB_State_StartBetting: // 開始下注
			pGame.setGameState(JLB_State_StopBettingAnimation, pGame.mapi64Timeout[JLB_Timeout_StopBettingAnimation])
			pGame.playStopBettingAnimation()
		case JLB_State_StopBettingAnimation: // 停止下注動畫
			pGame.setGameState(JLB_State_Deal, pGame.mapi64Timeout[JLB_Timeout_DealAnimation])
			pGame.deal()
		case JLB_State_Deal: // 發牌
			pGame.setGameState(JLB_State_BetDouble, pGame.mapi64Timeout[JLB_Timeout_DoubleBetting])
			pGame.betDouble()
		case JLB_State_BetDouble: // 加注
			if g_iHDCMode == HDCModeDisable {
				pGame.setGameState(JLB_State_BetSideSettlement, pGame.mapi64Timeout[JLB_Timeout_BetSideSettlement])
				pGame.showCards()
			} else if g_iHDCMode == HDCModeEnable {
				pGame.setGameState(DHC_State_Settlement, pGame.mapi64Timeout[DHC_Timeout_Settlement])
				pGame.settlement()
			}
		case JLB_State_BetSideSettlement: // 邊注結算
			for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
				pGame.arypUserData[iSeat].iPlayingState = JLB_UserPlayState_ShowCardsOK
			}
			pGame.setGameState(JLB_State_BetBaseSettlement, pGame.mapi64Timeout[JLB_Timeout_BetBaseSettlement])
			pGame.settlementBetBase()
		case JLB_State_BetBaseSettlement: // 底注結算
			pGame.setGameState(JLB_State_BetDoubleSettlement, pGame.mapi64Timeout[JLB_Timeout_BetDoubleSettlement])
			pGame.settlementBetDouble()
		case JLB_State_BetDoubleSettlement: // 加注結算
			fallthrough
		case DHC_State_Settlement:
			pGame.setGameState(JLB_State_GameOver, 0)
			pGame.gameOver()
		}
	}

	// 機器人動作
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == JLB_UserState_Normal {
			if pUser.i64Timeout > 0 && pGame.i64TimeDuration-pUser.i64LastMessageTime > pUser.i64Timeout {
				switch pGame.iGameState {
				case JLB_State_StartBetting: // 開始下注
					// 機器人隨機下邊注
					i64BetSide := g_i64DefaultBetSide
					if pGame.arypUserData[iSeat].i64Money-i64BetSide >= pGame.arypUserData[iSeat].i64BetBase*2 {
						pUser.i64BetSide = i64BetSide
						pUser.i64Money -= pUser.i64BetSide
						pGame.sendBetResult(pUser, JLB_BetType_Side, pUser.i64BetSide, JLB_BetResult_Success)
					}
				case JLB_State_BetDouble: // 加注
					// 有對子就加注，無AK就有機率棄牌
					if pGame.arypHandCards[iSeat].Pattern(false) >= pokerset.POKERPatternPair {
						// 有對子（含）以上的牌
						i64BetDouble := pUser.i64BetBase * 2
						if pUser.i64Money >= i64BetDouble {
							pUser.i64BetDouble = i64BetDouble
							pUser.i64Money -= pUser.i64BetDouble
							pGame.sendDoubleBetResult(pUser, JLB_DoubleBetResult_Success)
						}
					} else if pGame.arypHandCards[iSeat].Pattern(false) == pokerset.POKERPatternHighCard {
						// 高牌
						if !pGame.arypHandCards[iSeat].CheckIfHasAK() {
							iRnd := util.RndInt(0, 100)
							if iRnd < g_iFoldByNoAKOdds_Robot {
								// 棄牌
								pUser.iUserState = JLB_UserState_Fold
								pUser.i64MoneyChange += -pUser.i64BetBase
								pGame.sendUserFoldCard(pUser)
							}
						}
					}
				}

				pUser.i64Timeout = 0
			}
		}
	}
}

// 讀取遊戲設定
func (pGame *Logic) getGameConfig(configs *config.Group) {
	if configs == nil {
		mlog.Error("桌號[%d] 讀取 config 失敗", pGame.TableNo())
		return
	}

	str := ""

	// 讀取 gameserver 設定值
	cfgs := configs.GetConfig("gameserver")
	if cfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		pGame.iSeatCount = cfgs.GetConfigInt("gameserver", "seatcount")

		//----- timeout -----
		pGame.mapi64Timeout[JLB_Timeout_WaitingStart] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_WaitingStart")) * 10
		pGame.mapi64Timeout[JLB_Timeout_StartBettingAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_StartBettingAnimation")) * 10
		pGame.mapi64Timeout[JLB_Timeout_StartBetting] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_StartBetting")) * 10
		pGame.mapi64Timeout[JLB_Timeout_StopBettingAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_StopBettingAnimation")) * 10
		pGame.mapi64Timeout[JLB_Timeout_DealAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_DealAnimation")) * 10
		pGame.mapi64Timeout[JLB_Timeout_DoubleBetting] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_DoubleBetting")) * 10
		pGame.mapi64Timeout[JLB_Timeout_BetSideSettlement] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_BetSideSettlement")) * 10
		pGame.mapi64Timeout[JLB_Timeout_BetBaseSettlement] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_BetBaseSettlement")) * 10
		pGame.mapi64Timeout[JLB_Timeout_BetDoubleSettlement] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_BetDoubleSettlement")) * 10
		pGame.mapi64Timeout[DHC_Timeout_Settlement] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Settlement")) * 10

		// 籌碼面額
		str = cfgs.GetConfigValue("GameSetting", "ChipValue")
		arystrChip := strings.Split(str, ",")
		for i := 0; i < len(arystrChip); i++ {
			i64chip, _ := strconv.ParseInt(arystrChip[i], 10, 32)
			pGame.aryi32Chip = append(pGame.aryi32Chip, int32(i64chip))
		}

		g_i64DefaultBetSide = int64(cfgs.GetConfigInt64("GameSetting", "BetSide"))

		// 獎池回補基準值
		g_i64JackpotBase = int64(cfgs.GetConfigInt64("GameSetting", "JackpotBase"))
		g_i64Jackpot = g_i64JackpotBase

		//----- 特殊邊注 -----
		// 皇家同花順，此值為獎池的百分比比例
		pGame.mapi64JPSheet[pokerset.POKERPatternRoyalFlush] = int64(cfgs.GetConfigInt64("GameSetting", "JP_RoyalFlush"))
		// 同花順，此值為獎池的百分比比例
		pGame.mapi64JPSheet[pokerset.POKERPatternStraightFlush] = int64(cfgs.GetConfigInt64("GameSetting", "JP_StraightFlush"))
		// 四條，中獎金額
		pGame.mapi64JPSheet[pokerset.POKERPatternFourOfAKind] = int64(cfgs.GetConfigInt64("GameSetting", "JP_FourOfKind"))
		// 葫蘆，中獎金額
		pGame.mapi64JPSheet[pokerset.POKERPatternFullHouse] = int64(cfgs.GetConfigInt64("GameSetting", "JP_FullHouse"))
		// 同花，中獎金額
		pGame.mapi64JPSheet[pokerset.POKERPatternFlush] = int64(cfgs.GetConfigInt64("GameSetting", "JP_Flush"))

		//----- 賠注表 ----
		// 皇家同花順
		pGame.mapi64SpecialSheet[pokerset.POKERPatternRoyalFlush] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_RoyalFlush"))
		// 同花順
		pGame.mapi64SpecialSheet[pokerset.POKERPatternStraightFlush] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_StraightFlush"))
		// 四條
		pGame.mapi64SpecialSheet[pokerset.POKERPatternFourOfAKind] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_FourOfKind"))
		// 葫蘆
		pGame.mapi64SpecialSheet[pokerset.POKERPatternFullHouse] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_FullHouse"))
		// 同花
		pGame.mapi64SpecialSheet[pokerset.POKERPatternFlush] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_Flush"))
		// 順子
		pGame.mapi64SpecialSheet[pokerset.POKERPatternStraight] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_Straight"))
		// 三條
		pGame.mapi64SpecialSheet[pokerset.POKERPatternThreeOfAKind] = int64(cfgs.GetConfigInt64("GameSetting", "Pay_ThreeOfKind"))
		// 兩對
		pGame.mapi64SpecialSheet[pokerset.POKERPatternTwoPairs] = cfgs.GetConfigInt64("GameSetting", "Pay_TwoPairs")
		// 對子
		pGame.mapi64SpecialSheet[pokerset.POKERPatternPair] = cfgs.GetConfigInt64("GameSetting", "Pay_Other")
		// 高牌
		pGame.mapi64SpecialSheet[pokerset.POKERPatternHighCard] = cfgs.GetConfigInt64("GameSetting", "Pay_Other")

		g_iZhuangBetSide = cfgs.GetConfigInt("GameSetting", "OpenZhuangBetSide")
		g_iPluralBetSide = cfgs.GetConfigInt("GameSetting", "PluralBetSide")
		g_iPluralBetSideLimit = cfgs.GetConfigInt("GameSetting", "PluralBetSideLimit")
		g_iHDCMode = cfgs.GetConfigInt("GameSetting", "HDCMode")

		// 機器人參數
		g_iFoldByNoAKOdds_Robot = cfgs.GetConfigInt("robot", "Odds_FoldByNoAK")
		g_iBetSideOdds_Robot = cfgs.GetConfigInt("robot", "Odds_BetSide")

		// 是否開啟水線控制
		//iPoolControl := cfgs.GetConfigInt("SystemPool", "EnablePoolControl")
		//if iPoolControl == 0 {
		//	g_bEnablePoolControl = false
		//} else {
		//	g_bEnablePoolControl = true
		//}
		pGame.getDealerWeightTable(configs)
		// 水線控制
		//if g_bEnablePoolControl {
		//	if !g_bGetPoolSetting {
		//		g_map_pDealerWeightTable = make(map[int]*DealerWeightTable)
		//
		//		g_i64WinLine = cfgs.GetConfigInt64("SystemPool", "WinLine")                  // 贏錢線
		//		g_i32PlayerWeightDefault = cfgs.GetConfigInt32("SystemPool", "PlayerWeight") // 預設玩家拿牌權重
		//		g_i32DealerWeightDefault = cfgs.GetConfigInt32("SystemPool", "DealerWeight") // 預設莊家權重
		//		g_iCheatRate = cfgs.GetConfigInt("SystemPool", "CheatRate")                  // 換牌率
		//		g_iRoundLimit = cfgs.GetConfigInt("SystemPool", "RoundLimit")                // 判定局數
		//
		//		// 機率表
		//		iTableCount := cfgs.GetConfigInt("WeightTable", "TableCount") // 權重表張數
		//		for iNo := 0; iNo < iTableCount; iNo++ {
		//			if _, ok := g_map_pDealerWeightTable[iNo]; !ok {
		//				g_map_pDealerWeightTable[iNo] = new(DealerWeightTable)
		//			}
		//
		//			strName := fmt.Sprintf("%d", iNo)
		//			str = cfgs.GetConfigValue(strName, "Pattern")
		//			arystrPattern := strings.Split(str, ",")
		//			str = cfgs.GetConfigValue(strName, "Weight")
		//			arystrWeight := strings.Split(str, ",")
		//
		//			for i := 0; i < len(arystrPattern); i++ {
		//				pData := new(DealerWeightData)
		//				i64Pattern, _ := strconv.ParseInt(arystrPattern[i], 10, 32)
		//				pData.i32CardPattern = int32(i64Pattern)
		//				if i < len(arystrWeight) {
		//					i64Weight, _ := strconv.ParseInt(arystrWeight[i], 10, 32)
		//					pData.i32Weight = int32(i64Weight)
		//				}
		//
		//				g_map_pDealerWeightTable[iNo].ary_pWeightData = append(g_map_pDealerWeightTable[iNo].ary_pWeightData, pData)
		//			}
		//		}
		//
		//		// 依牌型大小降冪排序並計算亂數範圍
		//		for _, pTable := range g_map_pDealerWeightTable {
		//			if pTable != nil {
		//				// 排序
		//				iLen := len(pTable.ary_pWeightData)
		//				for i := 0; i < iLen; i++ {
		//					for j := i; j < iLen; j++ {
		//						if pTable.ary_pWeightData[i].i32CardPattern < pTable.ary_pWeightData[j].i32CardPattern {
		//							pTable.ary_pWeightData[i], pTable.ary_pWeightData[j] = pTable.ary_pWeightData[j], pTable.ary_pWeightData[i]
		//						}
		//					}
		//				}
		//
		//				// 計算亂數範圍
		//				i32Range := int32(0)
		//				for i := 0; i < len(pTable.ary_pWeightData); i++ {
		//					if pTable.ary_pWeightData[i].i32Weight <= 0 {
		//						pTable.ary_pWeightData[i].i32RndRangeMin = -1
		//						pTable.ary_pWeightData[i].i32RndRangeMax = -1
		//						continue
		//					}
		//
		//					pTable.ary_pWeightData[i].i32RndRangeMin = i32Range + 1
		//					i32Range += pTable.ary_pWeightData[i].i32Weight
		//					pTable.ary_pWeightData[i].i32RndRangeMax = i32Range
		//				}
		//				pTable.i32TotalWeight = i32Range
		//			}
		//		}
		//
		//		mlog.System("水线控制功能开启")
		//		g_bGetPoolSetting = true
		//	}
		//}
	}

	// 撲克牌發牌器
	pGame.pTestDealerCards = NewJlbpkPokerSet()
	pGame.arypTestUserCards = make([]*JlbpkPokerSet, pGame.iSeatCount)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.arypTestUserCards[iSeat] = NewJlbpkPokerSet()
	}

	pkcfgs := configs.GetConfig("poker")
	if pkcfgs == nil {
		pGame.iDealType = DealType_Natural
	} else {
		// 發牌器發牌模式
		iDealType := pkcfgs.GetConfigInt("Poker", "JPTestEnable")
		if iDealType < DealType_Natural || iDealType >= DealType_Max {
			iDealType = DealType_Natural
		}
		pGame.iDealType = iDealType

		// 測試用莊家手牌
		strCards := pkcfgs.GetConfigValue("Poker", "Dealer")
		pSet := new(pokerset.PokerCardSet)
		pSet.SetCardSet(strCards)
		if pSet.CardCount() < JLBPK_HandCardCount {
			pSet.SetCardSet("")
		} else if pSet.CardCount() > JLBPK_HandCardCount {
			pTemp := new(pokerset.PokerCardSet)
			for i := 0; i < JLBPK_HandCardCount; i++ {
				pCard := pSet.GetSingleCard(i, false)
				pTemp.AddSingleCard(pCard.CodeString())
			}
			pSet.SetCardSet(pTemp.CodeString())
		}

		pGame.pTestDealerCards.SetCardSet(pSet.CodeString())

		// 測試用玩家手牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			str = fmt.Sprintf("User%d", iSeat+1)
			strCards := pkcfgs.GetConfigValue("Poker", str)
			pSet := new(pokerset.PokerCardSet)
			pSet.SetCardSet(strCards)
			if pSet.CardCount() < JLBPK_HandCardCount {
				pSet.SetCardSet("")
			} else if pSet.CardCount() > JLBPK_HandCardCount {
				pTemp := new(pokerset.PokerCardSet)
				for i := 0; i < JLBPK_HandCardCount; i++ {
					pCard := pSet.GetSingleCard(i, false)
					pTemp.AddSingleCard(pCard.CodeString())
				}
				pSet.SetCardSet(pTemp.CodeString())
			}

			pGame.arypTestUserCards[iSeat].SetCardSet(pSet.CodeString())
		}
	}
}

// 加入新玩家
func (pGame *Logic) AddUser(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	iSeatNo := int(pGameUser.SeatNo)
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return
	}
	pUser := pGame.arypUserData[iSeatNo]

	if pUser.iUserState == JLB_UserState_EmptySeat {
		pUser.Init()
		pUser.User = pGameUser
		pUser.i32UserID = pGameUser.UserID
		pUser.i32SeatNo = pGameUser.SeatNo
		if pGameUser.IsRobot() {
			pUser.iUserType = gamelogic.UserTypeRebot
		}
		pUser.i64Money = pGameUser.Score
		pUser.ResetGameData()

		// 檢查玩家是否為監控玩家，如果是則設定相關資料
		//pHandler := pGame.pSuperToolService.GetJlbpkSuperToolHandler()
		//if ok, _ := util.InArray(pUser.i32UserID, pHandler.Proiority); ok {
		//	if pVIPUser := pHandler.GetVipUser(pUser.i32UserID); pVIPUser != nil {
		//		if pVIPUser.GetStatus() == 1 && pVIPUser.GetFinished() == 0 {
		//			pUser.bIsVIPFinished = false
		//		} else {
		//			pUser.bIsVIPFinished = true
		//		}
		//		pUser.i32Weight = pVIPUser.GetWeight()
		//		pUser.i64TargetPoint = pVIPUser.GetTargetPoint()
		//		pUser.i64AccumulatedPoint = pVIPUser.GetAccumulatedPoint()
		//	}
		//}

		if pGame.iGameState >= JLB_State_StartBetting {
			pUser.iUserState = JLB_UserState_Watch
		} else {
			pUser.iUserState = JLB_UserState_Normal

			switch pGame.iGameState {
			case JLB_State_WaitingPlayer: // 等待玩家加入
				fallthrough
			case JLB_State_WaitingStart: // 等待遊戲開始
				pUser.iPlayingState = JLB_UserPlayState_WaitingStart
			case JLB_State_StartBettingAnimation: // 開始下注動畫
				fallthrough
			case JLB_State_StartBetting: // 開始下注
				pUser.iPlayingState = JLB_UserPlayState_StartBettingOK
				if pUser.i64BetBase <= 0 {
					pUser.i64BetBase += int64(pGame.aryi32Chip[0])
					pUser.i64TotalBet += pUser.i64BetBase
					pUser.i64Money -= pUser.i64BetBase
					pUser.i64TotalWinLose -= pUser.i64BetBase
				}
			}
		}
	}

	pGame.iUserCount++
	mlog.Debug("Table[%d], AddUser, 玩家人數：%d", pGame.TableNo(), pGame.iUserCount)
	//if pGame.iUserCount > 0 && pGame.iGameState == JLB_State_WaitingPlayer {
	//
	//	pGame.setGameState(JLB_State_WaitingStart, pGame.mapi64Timeout[JLB_Timeout_WaitingStart])
	//	pGame.waitingGameBegin()
	//}
	pGame.sendGameInfo(pGameUser)
	pGame.sendAllUserInfo()

	example.invalid("加入新玩家[%d], 遊戲狀態[%s], 座位[%d], 類別[%d], 玩家狀態[%s], 玩家總人數[%d]", pUser.i32UserID,
		map_strGameState_debug[pGame.iGameState], iSeatNo, pUser.iUserType,
		map_strUserState_debug[pUser.iUserState], pGame.iUserCount)
	pGame.GameLog.AddGameStep("玩家[%d] 入座，状态：%s", pGameUser.UserID, map_strUserState[pUser.iUserState])
}

// 傳送牌局資訊
func (pGame *Logic) sendGameInfo(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	mlog.Debug("sendGameInfo, GameState[%s], ToClientState[%s], JP[%d], Time[%d], DealerCard[%s]",
		map_strGameState_debug[pGame.iGameState], map_strClientGameState_debug[pGame.toClientGameState()], g_i64Jackpot,
		pGame.i64Timeout-pGame.i64TimeDuration+pGame.i64LastMessageTime, pGame.pDealerHandCards.String())

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_GameInfo)

	pData := new(netproto.JLBPK_GameInfo)
	pData.GameState = proto.Int32(int32(pGame.toClientGameState())) // 遊戲狀態
	// Jackpot 賠率表
	for i := 0; i < len(pGame.arystrJPSheet); i++ {
		pData.JPSheet = append(pData.JPSheet, pGame.arystrJPSheet[i])
	}
	// 特殊牌賠率表
	for i := 0; i < len(pGame.arystrSpecialSheet); i++ {
		pData.SpecialSheet = append(pData.SpecialSheet, pGame.arystrSpecialSheet[i])
	}
	// 遊戲說明
	for i := 0; i < len(pGame.arystrPokerManual); i++ {
		pData.PokerManual = append(pData.PokerManual, pGame.arystrPokerManual[i])
	}
	pData.JPMoney = proto.Int64(g_i64Jackpot) // Jackpot 金額
	// 剩餘時間
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10))
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i, int(pGameUser.SeatNo))
		pData.UserInfo = append(pData.UserInfo, pUserInfo)
	}
	// 籌碼面額
	for i := 0; i < len(pGame.aryi32Chip); i++ {
		pData.ChipValue = append(pData.ChipValue, pGame.aryi32Chip[i])
	}
	pData.Bet_Side = proto.Int64(g_i64DefaultBetSide) // 邊注金額
	// 莊家手牌
	pDealerCard := new(netproto.JLBPK_UserCards)
	for iIdx := 0; iIdx < pGame.pDealerHandCards.CardCount(); iIdx++ {
		pCard := new(netproto.JLBPK_PokerCard)
		if pGame.iGameState >= JLB_State_BetSideSettlement {
			pSCard := pGame.pDealerHandCards.GetSingleCard(iIdx, false)
			pCard.Color = proto.Int32(pSCard.I32Color)
			pCard.CardNumber = proto.Int32(pSCard.I32Number)
		} else {
			if iIdx == 0 {
				pSCard := pGame.pDealerHandCards.GetSingleCard(iIdx, false)
				pCard.Color = proto.Int32(pSCard.I32Color)
				pCard.CardNumber = proto.Int32(pSCard.I32Number)
			} else {
				pCard.Color = proto.Int32(pokerset.POKERCardColorOther)
				pCard.CardNumber = proto.Int32(pokerset.POKERCardNoHide)
			}
		}
		pDealerCard.Cards = append(pDealerCard.Cards, pCard)
	}
	pDealerCard.Pattern = proto.Int32(pGame.pDealerHandCards.Pattern(true))
	pData.BankerCards = pDealerCard

	pMsg.MsgData = pData
	pGame.SendMessageToUser(pGameUser, pMsg)
}

// 傳送玩家資訊
func (pGame *Logic) sendAllUserInfo() {
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_AllUserInfo)

		pData := new(netproto.JLBPK_AllUserInfo)
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}
}

// 等待遊戲開始
func (pGame *Logic) waitingGameBegin() {
	pGame.GameLog.GameBegin()
	mlog.Debug("等待开始下注")
	pGame.GameLog.AddGameStep("等待开始下注")

	for i := 0; i < pGame.iSeatCount; i++ {
		pGame.arypUserData[i].ResetGameData()
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_WaitingStart)
	pData := new(netproto.JLBPK_WaitingStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10))

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	mlog.Debug("state[%s], timeout[%d], duration[%d], lastmsg[%d]", map_strGameState_debug[pGame.iGameState],
		pGame.i64Timeout, pGame.i64TimeDuration, pGame.i64LastMessageTime)
}

// 播放開始下注動畫
func (pGame *Logic) playStartBettingAnimation() {
	mlog.Debug("播放開始下注動畫")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_StartBettingAnimation)

	pGame.BoardcastMessage(pMsg)

	pGame.uiPlayerReady = 0

	// 改為播放開始下注動畫時就可以下注
	for i := 0; i < pGame.iSeatCount; i++ {
		//if pGame.arypUserData[i].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[i].iUserType == gamelogic.UserType_Rebot ||
		//	pGame.arypUserData[i].bOffline {
		//	pGame.arypUserData[i].iPlayingState = JLB_UserPlayState_StartBettingOK
		//	pGame.uiPlayerReady |= 1 << uint(i)
		//}
		pGame.arypUserData[i].iPlayingState = JLB_UserPlayState_StartBettingOK
	}

	pGame.setGameState(JLB_State_StartBetting, pGame.mapi64Timeout[JLB_Timeout_StartBetting])
	pGame.startBetting()
}

// 訊息處理 - 開始下注動畫播放完成
func (pGame *Logic) messageStartBettingAnimationOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_StartBettingAnimation {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_WaitingStart {
		return
	}

	mlog.Debug("User[%d] 開始下注動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_StartBettingOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家開始下注動畫播放完成")
		//pGame.startBetting()
	}
}

// 開始下注
func (pGame *Logic) startBetting() {
	mlog.Debug("開始下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_UserBettingStart)
	pData := new(netproto.JLBPK_UserBettingStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10))

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	pGame.GameLog.AddGameStep("开始下注")
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		pUser.iPlayingState = JLB_UserPlayState_Betting
		if pUser.iUserState == JLB_UserState_Normal {
			// 自動幫玩家下底注
			mlog.Debug("玩家[%d] 金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
			pUser.i64BetBase += int64(pGame.aryi32Chip[0])
			pUser.i64TotalBet += pUser.i64BetBase
			pUser.i64Money -= pUser.i64BetBase
			pUser.i64TotalWinLose -= pUser.i64BetBase
			pGame.sendBetResult(pUser, JLB_BetType_Base, pUser.i64BetBase, JLB_BetResult_Auto)

			// 隨機決定機器人是否要下邊注
			if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == JLB_UserState_Normal {
				iRnd := util.RndInt(0, 100)
				if iRnd < g_iBetSideOdds_Robot {
					// 要下邊注，設定觸發時間
					pUser.SetRobotTimeout(pGame.mapi64Timeout[JLB_Timeout_StartBetting] - 1)
					if pUser.i64Timeout > 0 {
						pUser.i64LastMessageTime = pGame.i64TimeDuration
					}
				}
			}
		}
	}
}

// 送給 client 的遊戲狀態
func (pGame *Logic) toClientGameState() int {
	switch pGame.iGameState {
	case JLB_State_WaitingPlayer: // 等待玩家加入
		fallthrough
	case JLB_State_WaitingStart: // 等待遊戲開始
		return JLB_CState_WaitingStart
	case JLB_State_StartBettingAnimation: // 開始下注動畫
		fallthrough
	case JLB_State_StartBetting: // 開始下注
		return JLB_CState_StartBetting
	case JLB_State_StopBettingAnimation: // 停止下注動畫
		fallthrough
	case JLB_State_Deal: // 發牌
		return JLB_CState_Deal
	case JLB_State_BetDouble: // 加注
		return JLB_CState_BetDouble
	case JLB_State_BetSideSettlement: // 邊注結算
		return JLB_CState_BetSideSettlement
	case JLB_State_BetDoubleSettlement: // 加注結算
		return JLB_CState_BetDoubleSettlement
	case JLB_State_BetBaseSettlement: // 底注結算
		return JLB_CState_BetBaseSettlement
	case JLB_State_GameOver: // 遊戲結束
		return JLB_CState_WaitingStart
	case DHC_State_Settlement: // 遊戲結束
		return DHC_CState_Settlement
	}

	return -1
}

// 訊息處理 - 玩家下注，包含底注及邊注
func (pGame *Logic) messageBet(pUser *JlbpkUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	msgInfo := &netproto.JLBPK_Bet{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	i32BetType := msgInfo.GetBetType()
	i64Bet := msgInfo.GetBet()

	example.invalid("玩家[%d] 下%s %s", pUser.i32UserID, map_strBetType[i32BetType], Global.MoneyString(i64Bet))
	pGame.GameLog.AddGameStep("玩家[%d] 下%s %s", pUser.i32UserID, map_strBetType[i32BetType], Global.MoneyString(i64Bet))

	if pGame.iGameState != JLB_State_StartBetting {
		pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BettingStopped)
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_Betting {
		pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BettingStopped)
		return
	}

	if i64Bet <= 0 {
		pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BetSmallThanZero)
		return
	}

	// 剩餘金錢是否足以下注
	if pUser.i64Money < i64Bet {
		pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_NotEnoughMoney)
		return
	}

	switch i32BetType {
	case JLB_BetType_Base: // 底注
		// 下注是否超過上限
		if pUser.i64BetBase+i64Bet > (pUser.i64Money+pUser.i64BetBase)/3 {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		pUser.i64BetBase += i64Bet
		pUser.i64TotalBet += i64Bet
		pUser.i64TotalWinLose -= i64Bet
	case JLB_BetType_Side: // 邊注
		// 是否已下過邊注
		if pUser.i64BetSide > 0 && g_iPluralBetSide == PluralBetSideDisable {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BetSideHasDone)
			return
		}

		// 邊注值是否正確
		if i64Bet != g_i64DefaultBetSide {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		// 邊注是否超過上限
		if pUser.i64Money-i64Bet < pUser.i64BetBase*2 {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		// 複數邊注是否超過限制
		if g_iPluralBetSide == PluralBetSideEnable && (pUser.i64BetSide+i64Bet > int64(g_iPluralBetSideLimit) && g_iPluralBetSideLimit > 0) {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		pUser.i64BetSide += i64Bet
		pUser.i64TotalBet += i64Bet
		pUser.i64TotalWinLose -= i64Bet
	case JLB_BetType_ZhuangSide:
		// 是否開啟莊家邊注
		if g_iZhuangBetSide == ZhuangBetSideDisable {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_NotOpenZhuangSide)
			return
		}

		// 是否已下過邊注
		if pUser.i64ZhuangBetSide > 0 && g_iPluralBetSide == PluralBetSideDisable {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BetSideHasDone)
			return
		}

		// 邊注值是否正確
		if i64Bet != g_i64DefaultBetSide {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		// 邊注是否超過上限
		if pUser.i64Money-i64Bet < pUser.i64BetBase*2 {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		// 複數邊注是否超過限制
		if g_iPluralBetSide == PluralBetSideEnable && (pUser.i64ZhuangBetSide+i64Bet > int64(g_iPluralBetSideLimit) && g_iPluralBetSideLimit > 0) {
			pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_BaseBetOverUpperLimit)
			return
		}

		pUser.i64ZhuangBetSide += i64Bet
		pUser.i64TotalBet += i64Bet
		pUser.i64TotalWinLose -= i64Bet
	default:
		return
	}

	pUser.i64Money -= i64Bet
	pGame.sendBetResult(pUser, i32BetType, i64Bet, JLB_BetResult_Success)
}

// 傳送玩家下注結果
func (pGame *Logic) sendBetResult(pUser *JlbpkUser, i32BetType int32, i64Bet int64, iResult int) {
	if pUser == nil {
		return
	}

	iBetResult := iResult
	if iResult >= JLB_BetResult_Success {
		iBetResult = JLB_BetResult_Success
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_BetResult)

	pData := new(netproto.JLBPK_BetResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JLBPK_ClientStartUpSeatNo)) // 下注玩家座位
	pData.BetType = proto.Int32(i32BetType)                                        // 下注類別
	pData.Bet = proto.Int64(i64Bet)                                                // 下注金額
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                                // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(iBetResult))                                  // 下注結果

	pMsg.MsgData = pData
	if iBetResult == JLB_BetResult_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}

	if iResult == JLB_BetResult_Auto {
		mlog.Debug("玩家[%d] 自动下底注：%s, 剩余金钱 %s", pUser.i32UserID,
			Global.MoneyString(pUser.i64BetBase), Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 自动下底注：%s, 剩余金钱 %s", pUser.i32UserID,
			Global.MoneyString(pUser.i64BetBase), Global.MoneyString(pUser.i64Money))
	} else {
		example.invalid("玩家[%d] 下注结果：%s, 剩余金钱 %s", pUser.i32UserID, map_strBetResult[iResult],
			Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 下注结果：%s, 剩余金钱 %s", pUser.i32UserID, map_strBetResult[iResult],
			Global.MoneyString(pUser.i64Money))
	}
}

// 播放停止下注動畫
func (pGame *Logic) playStopBettingAnimation() {
	mlog.Debug("播放停止下注動畫")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_StopBettingAnimation)

	pGame.BoardcastMessage(pMsg)

	pGame.uiPlayerReady = 0

	for i := 0; i < pGame.iSeatCount; i++ {
		if pGame.arypUserData[i].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[i].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[i].bOffline {
			pGame.arypUserData[i].iPlayingState = JLB_UserPlayState_StopBettingOK
			pGame.uiPlayerReady |= 1 << uint(i)
		}
	}
}

// 訊息處理 - 停止下注動畫播放完成
func (pGame *Logic) messageStopBettingAnimationOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_StopBettingAnimation {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_Betting {
		return
	}

	mlog.Debug("User[%d] 停止下注動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_StopBettingOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家停止下注動畫播放完成")
		pGame.setGameState(JLB_State_Deal, pGame.mapi64Timeout[JLB_Timeout_DealAnimation])
		pGame.deal()
	}
}

// 發牌
func (pGame *Logic) deal() {
	mlog.Debug("發牌")

	// 發牌
	mlog.Debug("iJPTest = %d", pGame.iDealType)
	switch pGame.iDealType {
	case DealType_RandomJP: // 隨機發彩金牌型
		mlog.Debug("JLB_JPTest_Random")
		// 發玩家手牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
				iRnd := util.RndInt(0, len(g_ary_i32JPPattern))
				i32Pattern := g_ary_i32JPPattern[iRnd]
				pSet := pGame.pPokerSet.DealByPattern(i32Pattern, JLBPK_HandCardCount, false)
				pGame.arypHandCards[iSeat].SetCardSet(pSet.CodeString())
			}
		}

		// 發莊家牌，莊家牌不排序
		pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
		pGame.pDealerHandCards.SetCardSet(pSet.CodeString())
	case DealType_Assigned: // 指定玩家手牌
		mlog.Debug("JLB_JPTest_Assigned")
		// 先將指定手牌從遊戲牌組中抽出
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
				if pGame.arypTestUserCards[iSeat].CardCount() > 0 {
					pSet := new(pokerset.PokerCardSet)
					pSet.SetCardSet(pGame.arypTestUserCards[iSeat].CodeString())
					pGame.pPokerSet.RemoveCardSet(pSet)
				}
			}
		}

		if pGame.pTestDealerCards.CardCount() > 0 {
			pSet := new(pokerset.PokerCardSet)
			pSet.SetCardSet(pGame.pTestDealerCards.CodeString())
			pGame.pPokerSet.RemoveCardSet(pSet)
		}

		// 發牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
				if pGame.arypTestUserCards[iSeat].CardCount() > 0 {
					pGame.arypHandCards[iSeat].SetCardSet(pGame.arypTestUserCards[iSeat].CodeString())
				} else {
					pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
					pGame.arypHandCards[iSeat].SetCardSet(pSet.CodeString())
				}
			}
		}

		if pGame.pTestDealerCards.CardCount() > 0 {
			pGame.pDealerHandCards.SetCardSet(pGame.pTestDealerCards.CodeString())
		} else {
			pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
			pGame.pDealerHandCards.SetCardSet(pSet.CodeString())
		}
	default:
		mlog.Debug("default")
		if !pGame.checkIfChangeToCheatMode() {
			if pGame.SysPool.GetMode() == gamelogic.SysPoolModeNone {
				// 自然發牌
				// 發玩家手牌
				for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
					if pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
						pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
						pGame.arypHandCards[iSeat].SetCardSet(pSet.CodeString())
					}
				}

				// 莊家牌
				pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
				pGame.pDealerHandCards.SetCardSet(pSet.CodeString())
			} else {
				// 水線控制
				// 先發莊家牌，以利取得要指定給莊家的牌型，莊家牌不排序
				mlog.System("目前使用水線控制表 %s", g_map_strDealerWeightNo[g_GameData.WeightTableNoNow])
				pSet := dealDealerHandCard(pGame.pPokerSet, JLBPK_HandCardCount, g_map_pDealerWeightTable[g_GameData.WeightTableNoNow])
				pGame.pDealerHandCards.SetCardSet(pSet.CodeString())

				// 發玩家手牌
				for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
					if pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
						pSet := pGame.pPokerSet.Deal(JLBPK_HandCardCount)
						pGame.arypHandCards[iSeat].SetCardSet(pSet.CodeString())
					}
				}
			}
		} else {
			// 換牌模式
			arypChangingData := changingCards(pGame.pPokerSet, g_i32DealerWeight, pGame.arypUserData)
			for i := 0; i < len(arypChangingData); i++ {
				if arypChangingData[i].i32seatNo >= 0 && arypChangingData[i].i32seatNo < int32(pGame.iSeatCount) {
					// 玩家
					pGame.arypHandCards[arypChangingData[i].i32seatNo].SetCardSet(arypChangingData[i].strCardCodeString)
				} else if arypChangingData[i].i32seatNo < 0 {
					// 莊家
					pGame.pDealerHandCards.SetCardSet(arypChangingData[i].strCardCodeString)
				}
			}

			mlog.Debug("莊家【權重(%d), %s [%s(%d)]】", g_i32DealerWeight, pGame.pDealerHandCards.CodeString(),
				pGame.pDealerHandCards.PatternString(true), pGame.pDealerHandCards.Pattern(true))
			for i := 0; i < len(pGame.arypUserData); i++ {
				mlog.Debug("玩家(%d)【權重(%d), %s [%s(%d)]】", i, pGame.arypUserData[i].i32Weight, pGame.arypHandCards[i].CodeString(),
					pGame.arypHandCards[i].PatternString(false), pGame.arypHandCards[i].Pattern(false))
			}

			// 將發出的牌自遊戲牌盒中抽出
			pGame.pPokerSet.DrawCards(pGame.pDealerHandCards.CodeString())
			for iSeat := 0; iSeat < len(pGame.arypHandCards); iSeat++ {
				pGame.pPokerSet.DrawCards(pGame.arypHandCards[iSeat].CodeString())
			}
		}
	}

	// 排序玩家手牌
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.arypHandCards[iSeat].Sort(false, true, true)
	}

	// 莊家牌是否成局
	if pGame.pDealerHandCards.Pattern(true) == pokerset.POKERPatternNone {
		pGame.bDealerIsGood = false
	} else {
		pGame.bDealerIsGood = true
	}

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("Seat[%d], User[%d], %s[%s]", i, pGame.arypUserData[i].i32UserID,
			pGame.arypHandCards[i].String(), pGame.arypHandCards[i].PatternString(false))
	}
	mlog.Debug("莊家: %s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))
	mlog.Debug("PoketSet left[%d]", pGame.pPokerSet.CardCount())

	// 傳送訊息給 client
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_Deal)

		pData := new(netproto.JLBPK_Deal)
		// 莊家手牌
		pDealerCard := new(netproto.JLBPK_UserCards)
		for iIdx := 0; iIdx < pGame.pDealerHandCards.CardCount(); iIdx++ {
			pCard := new(netproto.JLBPK_PokerCard)
			if iIdx == 0 {
				pSCard := pGame.pDealerHandCards.GetSingleCard(iIdx, false)
				pCard.Color = proto.Int32(pSCard.I32Color)
				pCard.CardNumber = proto.Int32(pSCard.I32Number)
			} else {
				pCard.Color = proto.Int32(pokerset.POKERCardColorOther)
				pCard.CardNumber = proto.Int32(pokerset.POKERCardNoHide)
			}
			pDealerCard.Cards = append(pDealerCard.Cards, pCard)
		}
		pDealerCard.Pattern = proto.Int32(pGame.pDealerHandCards.Pattern(true))
		pData.BankerCards = pDealerCard
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}

	pGame.GameLog.AddGameStep("发牌")
	for i := 0; i < pGame.iSeatCount; i++ {
		if pGame.arypUserData[i].iUserState != JLB_UserState_EmptySeat {
			pGame.GameLog.AddGameStep("玩家[%d], 状态[%s]：%s, 牌型[%s]", pGame.arypUserData[i].i32UserID,
				map_strUserState[pGame.arypUserData[i].iUserState], pGame.arypHandCards[i].String(),
				pGame.arypHandCards[i].PatternString(false))
		}
	}
	pGame.GameLog.AddGameStep("庄家：%s, 牌型[%s]",
		pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))

	pGame.uiPlayerReady = 0

	for i := 0; i < pGame.iSeatCount; i++ {
		if pGame.arypUserData[i].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[i].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[i].bOffline {
			pGame.arypUserData[i].iPlayingState = JLB_UserPlayState_DealOK
			pGame.uiPlayerReady |= 1 << uint(i)
		}
	}
}

// 訊息處理 - 發牌動畫播放完成
func (pGame *Logic) messageDealOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_Deal {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_StopBettingOK {
		return
	}

	mlog.Debug("User[%d] 發牌動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_DealOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家發牌動畫播放完成")
		pGame.setGameState(JLB_State_BetDouble, pGame.mapi64Timeout[JLB_Timeout_DoubleBetting])
		pGame.betDouble()
	}
}

// 玩家開始加注
func (pGame *Logic) betDouble() {
	mlog.Debug("玩家開始加注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_DoubleBetStart)

	pData := new(netproto.JLBPK_DoubleBetStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData

	pGame.GameLog.AddGameStep("玩家开始加注")

	// 機器人操作：有對子就加注，無AK就有機率棄牌
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserType == gamelogic.UserTypeRebot && pGame.arypUserData[iSeat].iUserState == JLB_UserState_Normal {
			// 設定觸發時間
			pGame.arypUserData[iSeat].SetRobotTimeout(pGame.mapi64Timeout[JLB_Timeout_DoubleBetting] - 1)
			if pGame.arypUserData[iSeat].i64Timeout > 0 {
				pGame.arypUserData[iSeat].i64LastMessageTime = pGame.i64TimeDuration
			}
		}
	}

	pGame.uiPlayerReady = 0

	for i := 0; i < pGame.iSeatCount; i++ {
		pGame.arypUserData[i].iPlayingState = JLB_UserPlayState_BetDouble
	}

	pGame.BoardcastMessage(pMsg)
}

// 訊息處理 - 玩家加注
func (pGame *Logic) messageDoubleBet(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_BetDouble {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_BetDouble {
		return
	}

	if pUser.i64BetDouble > 0 {
		pGame.sendDoubleBetResult(pUser, JLB_DoublebetResult_BetDoubleHasDone)
		return
	}

	i64Doublebet := pUser.i64BetBase * 2

	if pUser.i64Money < i64Doublebet {
		pGame.sendDoubleBetResult(pUser, JLB_DoubleBetResult_NotEnoughMoney)
		pGame.sendUserFoldCard(pUser)
		return
	}

	pUser.i64BetDouble += i64Doublebet
	pUser.i64TotalBet += i64Doublebet
	pUser.i64Money -= i64Doublebet
	pUser.i64TotalWinLose -= i64Doublebet
	pGame.sendDoubleBetResult(pUser, JLB_DoubleBetResult_Success)
}

// 傳送玩家加注結果
func (pGame *Logic) sendDoubleBetResult(pUser *JlbpkUser, iResult int) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_DoubleBetResult)

	pData := new(netproto.JLBPK_DoubleBetResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JLBPK_ClientStartUpSeatNo)) // 加注玩家座位
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                                // 加注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(iResult))                                     // 加注結果
	pData.Bet = proto.Int64(pUser.i64BetDouble)                                    // 加注金額

	pMsg.MsgData = pData
	if iResult == JLB_DoubleBetResult_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}

	example.invalid("玩家[%d] 加注结果：加注 %s, %s, 剩余金钱 %s", pUser.i32UserID,
		Global.MoneyString(pUser.i64BetDouble), map_strBetResult[iResult],
		Global.MoneyString(pUser.i64Money))
	pGame.GameLog.AddGameStep("玩家[%d] 加注结果：加注 %s, %s, 剩余金钱 %s", pUser.i32UserID,
		Global.MoneyString(pUser.i64BetDouble), map_strBetResult[iResult],
		Global.MoneyString(pUser.i64Money))
}

// 訊息處理 - 棄牌
func (pGame *Logic) messageFold(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_BetDouble {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_BetDouble {
		return
	}

	pUser.iUserState = JLB_UserState_Fold
	pUser.i64MoneyChange += -pUser.i64BetBase
	pGame.sendUserFoldCard(pUser)
}

// 傳送玩家棄牌訊息
func (pGame *Logic) sendUserFoldCard(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_UserFoldCard)

	pData := new(netproto.JLBPK_UserFoldCard)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JLBPK_ClientStartUpSeatNo)) // 棄牌玩家座位

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	mlog.Debug("玩家[%d] 弃牌", pUser.i32UserID)
	pGame.GameLog.AddGameStep("玩家[%d] 弃牌", pUser.i32UserID)
}

// 更新 Jackpot 金額
func (pGame *Logic) sendUpdateJackpot() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_UpdateJackpot)

	pData := new(netproto.JLBPK_UpdateJackpot)
	pData.JPMoney = proto.Int64(g_i64Jackpot) // Jackpot 金額

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)
}

// 開牌，邊注結算
func (pGame *Logic) showCards() {
	mlog.Debug("開牌，邊注結算")
	pGame.GameLog.AddGameStep("开牌，边注结算")

	// 未加注玩家棄牌
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.iUserState == JLB_UserState_Normal &&
			pUser.i64BetDouble <= 0 {
			pUser.iUserState = JLB_UserState_Fold
			pUser.i64MoneyChange += -pUser.i64BetBase
			pGame.sendUserFoldCard(pUser)
		}
	}

	// 計算彩金
	mlog.Debug("原彩金金额 %s", Global.MoneyString(g_i64Jackpot))
	pGame.GameLog.AddGameStep("原彩金金额 %s", Global.MoneyString(g_i64Jackpot))

	allJPMoney := int64(0)
	baseJackpotGap := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.iUserState == JLB_UserState_Normal || pUser.iUserState == JLB_UserState_Fold {
			mlog.Debug("玩家[%d] 原来金钱 %s", pUser.i32SeatNo, Global.MoneyString(pUser.i64Money))

			bGetJP, i64JPMoney := pGame.checkIfGetJackpot(iSeat)
			pUser.bGetJackpot = bGetJP
			pUser.i64JPMoney = i64JPMoney
			g_i64Jackpot -= pUser.i64JPMoney
			allJPMoney += pUser.i64JPMoney
			pGame.sendUpdateJackpot() // 更新獎池金額

			i64MoneyChange := pUser.i64JPMoney - pUser.i64BetSide
			pUser.i64WinLoseMoney = pUser.i64JPMoney - pUser.i64BetSide
			pUser.i64TotalWinLose += pUser.i64JPMoney
			pUser.i64Money += pUser.i64JPMoney
			pUser.i64MoneyChange += i64MoneyChange

			// 寫分
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)

			str := ""
			switch pUser.iUserState {
			case JLB_UserState_Fold:
				str = fmt.Sprintf("玩家[%d] 弃牌, 边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64BetSide),
					pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
			default:
				str = fmt.Sprintf("玩家[%d] 边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64BetSide),
					pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
			}

			if !bGetJP {
				str += fmt.Sprintf("未拉中彩金, 输赢 %s, 结算后金钱 %s", Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
			} else {
				str += fmt.Sprintf("拉中彩金 %s, 输赢 %s, 结算后金钱 %s", Global.MoneyString(pUser.i64JPMoney),
					Global.MoneyString(pUser.i64WinLoseMoney), Global.MoneyString(pUser.i64Money))
			}

			mlog.Debug("%s", str)
			mlog.Debug("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))
			pGame.GameLog.AddGameStep("%s", str)
			pGame.GameLog.AddGameStep("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))

			if g_i64Jackpot < g_i64JackpotBase {
				baseJackpotGap += g_i64JackpotBase - g_i64Jackpot
				g_i64Jackpot = g_i64JackpotBase
				pGame.sendUpdateJackpot() // 更新獎池金額
				mlog.Debug("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
				pGame.GameLog.AddGameStep("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
			}
		}
	}

	// 收邊注
	i64BetSide := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserType != gamelogic.UserTypeRebot {
			i64BetSide += pGame.arypUserData[iSeat].i64BetSide
		}
	}
	mlog.Debug("原彩金金额 %s，共计收取边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64BetSide))
	pGame.GameLog.AddGameStep("原彩金金额 %s，共计收取边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64BetSide))

	g_i64Jackpot += i64BetSide

	pGame.sendUpdateJackpot() // 更新獎池金額

	// 更新DB彩金池
	pGame.UpdateJackpotPool(0, i64BetSide, allJPMoney)

	mlog.Debug("收取边注后彩金金额 %s", Global.MoneyString(g_i64Jackpot))
	pGame.GameLog.AddGameStep("收取边注后彩金金额 %s", Global.MoneyString(g_i64Jackpot))

	// 傳送訊息給 client
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_ShowCards)

		pData := new(netproto.JLBPK_ShowCards)
		// 莊家手牌
		pDealerCard := new(netproto.JLBPK_UserCards)
		for iIdx := 0; iIdx < pGame.pDealerHandCards.CardCount(); iIdx++ {
			pCard := new(netproto.JLBPK_PokerCard)
			pSCard := pGame.pDealerHandCards.GetSingleCard(iIdx, false)
			pCard.Color = proto.Int32(pSCard.I32Color)
			pCard.CardNumber = proto.Int32(pSCard.I32Number)

			pDealerCard.Cards = append(pDealerCard.Cards, pCard)
		}
		pDealerCard.Pattern = proto.Int32(pGame.pDealerHandCards.Pattern(true))
		pData.BankerCards = pDealerCard
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[iSeat].bOffline {
			pGame.arypUserData[iSeat].iPlayingState = JLB_UserPlayState_ShowCardsOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}

		// 邊注結算完，將邊注值歸零
		//pGame.arypUserData[iSeat].i64BetSide = 0
	}
}

// 訊息處理 - 開牌，邊注結算效果完成
func (pGame *Logic) messageShowCardsOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_BetSideSettlement {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_BetDouble {
		return
	}

	mlog.Debug("User[%d] 邊注結算動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_ShowCardsOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家邊注結算效果完成")
		pGame.setGameState(JLB_State_BetBaseSettlement, pGame.mapi64Timeout[JLB_Timeout_BetBaseSettlement])
		pGame.settlementBetBase()
	}
}

// 檢查是否獲得彩金金額，回傳值［是否中彩金］、［所中彩金金額］
func (pGame *Logic) checkIfGetJackpot(iSeatNo int) (bool, int64) {
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return false, 0
	}

	if pGame.arypUserData[iSeatNo].i64BetSide <= 0 {
		return false, 0
	}

	if pGame.arypUserData[iSeatNo].iUserState != JLB_UserState_Normal {
		return false, 0
	}

	bGetJP, i64JPMoney, i64Jackpot := false, int64(0), g_i64Jackpot

	i32Pattern := pGame.arypHandCards[iSeatNo].Pattern(false)
	if i64Value, ok := pGame.mapi64JPSheet[i32Pattern]; ok {
		switch i32Pattern {
		case pokerset.POKERPatternRoyalFlush: // 皇家同花順
			fallthrough
		case pokerset.POKERPatternStraightFlush: // 同花順
			bGetJP = true
			i64JPMoney = i64Jackpot * i64Value / 100
		default:
			bGetJP = true
			i64JPMoney = i64Value
			// 若是複數邊注
			if g_iPluralBetSide == PluralBetSideEnable {
				i64JPMoney *= pGame.arypUserData[iSeatNo].i64BetSide / g_i64DefaultBetSide
			}
		}
	}

	return bGetJP, i64JPMoney
}

// 檢查是否獲得莊家彩金金額，回傳值［是否中彩金］、［所中彩金金額］
func (pGame *Logic) checkIfGetZhuangJackpot(iSeatNo int) (bool, int64) {
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return false, 0
	}

	if pGame.arypUserData[iSeatNo].i64ZhuangBetSide <= 0 {
		return false, 0
	}

	if pGame.arypUserData[iSeatNo].iUserState != JLB_UserState_Normal {
		return false, 0
	}

	bGetJP, i64JPMoney, i64Jackpot := false, int64(0), g_i64Jackpot

	i32Pattern := pGame.pDealerHandCards.Pattern(false)
	if i64Value, ok := pGame.mapi64JPSheet[i32Pattern]; ok {
		switch i32Pattern {
		case pokerset.POKERPatternRoyalFlush: // 皇家同花順
			fallthrough
		case pokerset.POKERPatternStraightFlush: // 同花順
			bGetJP = true
			i64JPMoney = i64Jackpot * i64Value / 100
		default:
			bGetJP = true
			i64JPMoney = i64Value
			// 若是複數邊注
			if g_iPluralBetSide == PluralBetSideEnable {
				i64JPMoney *= pGame.arypUserData[iSeatNo].i64ZhuangBetSide / g_i64DefaultBetSide
			}
		}
	}

	return bGetJP, i64JPMoney
}

// 賠付，底注結算 - ［加勒比海撲克］流程
func (pGame *Logic) settlementBetBase() {
	mlog.Debug("赔付，底注结算")
	pGame.GameLog.AddGameStep("赔付，底注结算")

	//----- 底注結算 -----
	pDealerSet := pGame.pDealerHandCards.GetPokerCardSet()
	mlog.Debug("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))
	pGame.GameLog.AddGameStep("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))

	i64BetBaseWinLose := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		switch pUser.iUserState {
		case JLB_UserState_Normal: // 正常
			pPlayerSet := pGame.arypHandCards[iSeat].GetPokerCardSet()
			if !pGame.bDealerIsGood {
				// 莊家不成局
				i64BetBaseWinLose = pUser.i64BetBase * int64(1+JLBPK_BetBasePayRate)
				pUser.i64WinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
				pUser.i64TotalWinLose += i64BetBaseWinLose
				pUser.i64Money += i64BetBaseWinLose
				pUser.iWinState = JLB_WinState_Win

				pGame.pWinLose.bDealerNotOK = true
				if !pUser.User.IsRobot() {
					pGame.pWinLose.i64PoolAmount -= pUser.i64BetBase * int64(JLBPK_BetBasePayRate)
				}
			} else {
				iCmpResult := pGame.pPokerSet.Compare(pDealerSet, pPlayerSet, false, false)
				switch iCmpResult {
				case pokerset.POKERCompareEven: // 平手
					i64BetBaseWinLose = pUser.i64BetBase
					pUser.iWinState = JLB_WinState_Even
				case pokerset.POKERCompareDealer: // 莊家贏
					i64BetBaseWinLose = -pUser.i64BetBase
					pUser.iWinState = JLB_WinState_Lose
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount += pUser.i64BetBase
					}
				case pokerset.POKERComparePlayer: // 玩家贏
					i64BetBaseWinLose = pUser.i64BetBase * int64(1+JLBPK_BetBasePayRate)
					pUser.iWinState = JLB_WinState_Win
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount -= pUser.i64BetBase * int64(JLBPK_BetBasePayRate)
					}
				}

				if iCmpResult == pokerset.POKERCompareDealer {
					pUser.i64WinLoseMoney = i64BetBaseWinLose
					pUser.i64TotalBet -= pUser.i64BetBase
				} else {
					pUser.i64WinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
					pUser.i64Money += i64BetBaseWinLose
					pUser.i64TotalWinLose += i64BetBaseWinLose
				}
			}

			// 寫分
			i64WinLose := i64BetBaseWinLose
			if i64BetBaseWinLose > 0 {
				i64WinLose = i64BetBaseWinLose - pUser.i64BetBase
			}
			i64MoneyChange, i64Tax := pGame.GetUserWinAndTax(i64WinLose)
			pUser.i64TotalWinLose -= i64Tax
			pUser.i64Money -= i64Tax
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, i64Tax)

			mlog.Debug("玩家[%d], %s[%s], 底注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 底注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
		case JLB_UserState_Fold: // 棄牌
			i64BetBaseWinLose = -pUser.i64BetBase
			pUser.i64WinLoseMoney = i64BetBaseWinLose
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64BetBaseWinLose, 0)
			pUser.iWinState = JLB_WinState_Lose
			if !pUser.User.IsRobot() {
				pGame.pWinLose.i64PoolAmount += pUser.i64BetBase
			}

			mlog.Debug("玩家[%d] 弃牌, %s[%s], 底注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64BetBaseWinLose), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 弃牌, %s[%s], 底注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64BetBaseWinLose), Global.MoneyString(pUser.i64Money))
		}
	}

	//----- 傳送訊息給 client -----
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_Compensate_BetBase)

		pData := new(netproto.JLBPK_Compensate_BetBase)
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[iSeat].bOffline {
			pGame.arypUserData[iSeat].iPlayingState = JLB_UserPlayState_Compensate_BetBaseOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}

		// 底注結算完，將底注值歸零
		//pGame.arypUserData[iSeat].i64BetBase = 0
	}
}

// 訊息處理 - 底注結算效果完成
func (pGame *Logic) messageBetBaseSettlementOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_BetBaseSettlement {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_ShowCardsOK {
		return
	}

	mlog.Debug("User[%d] 底注結算動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_Compensate_BetBaseOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家底注結算效果完成")
		pGame.setGameState(JLB_State_BetDoubleSettlement, pGame.mapi64Timeout[JLB_Timeout_BetDoubleSettlement])
		pGame.settlementBetDouble()
	}
}

// 賠付，加注結算 - ［加勒比海撲克］流程
func (pGame *Logic) settlementBetDouble() {
	mlog.Debug("賠付，加注結算")
	pGame.GameLog.AddGameStep("赔付，加注结算")

	//----- 加注結算 -----
	pDealerSet := pGame.pDealerHandCards.GetPokerCardSet()
	// 莊家是否成局
	mlog.Debug("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))
	pGame.GameLog.AddGameStep("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))

	i64DoubleWinLose := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		switch pUser.iUserState {
		case JLB_UserState_Normal: // 正常
			pPlayerSet := pGame.arypHandCards[iSeat].GetPokerCardSet()
			if !pGame.bDealerIsGood {
				// 莊家不成局
				i64DoubleWinLose = pUser.i64BetDouble
				pUser.i64WinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
				pUser.i64TotalWinLose += i64DoubleWinLose
				pUser.i64Money += i64DoubleWinLose

			} else {
				iCmpResult := pGame.pPokerSet.Compare(pDealerSet, pPlayerSet, false, false)
				switch iCmpResult {
				case pokerset.POKERCompareEven: // 平手
					i64DoubleWinLose = pUser.i64BetDouble
				case pokerset.POKERCompareDealer: // 第一組牌
					i64DoubleWinLose = -pUser.i64BetDouble
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount += pUser.i64BetDouble
					}
				case pokerset.POKERComparePlayer: // 第二組牌
					i64DoubleWinLose = pUser.i64BetDouble * (1 + pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)])
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount -= pUser.i64BetDouble * pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)]
					}
				}

				if iCmpResult == pokerset.POKERCompareDealer {
					pUser.i64WinLoseMoney = i64DoubleWinLose
					pUser.i64TotalBet -= pUser.i64BetDouble
				} else {
					pUser.i64WinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
					pUser.i64Money += i64DoubleWinLose
					pUser.i64TotalWinLose += i64DoubleWinLose
				}
			}

			// 寫分
			i64WinLose := i64DoubleWinLose
			if i64DoubleWinLose > 0 {
				i64WinLose = i64DoubleWinLose - pUser.i64BetDouble
			}
			i64MoneyChange, i64Tax := pGame.GetUserWinAndTax(i64WinLose)
			pUser.i64TotalWinLose -= i64Tax
			pUser.i64Money -= i64Tax
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			// i64Flow := int64(0)
			// if i64MoneyChange > 0 {
			// 	i64Flow = i64MoneyChange
			// } else {
			// 	i64Flow = -i64MoneyChange
			// }

			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, i64Tax)

			if pGame.bDealerIsGood {
				mlog.Debug("玩家[%d], %s[%s], 加注 %s, 赔注[1赔%d], 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)],
					Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 加注 %s, 赔注[1赔%d], 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)],
					Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
			} else {
				mlog.Debug("玩家[%d], %s[%s], 加注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax),
					Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 加注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax),
					Global.MoneyString(pUser.i64Money))
			}
		case JLB_UserState_Fold: // 棄牌
			i64DoubleWinLose = -pUser.i64BetDouble
			pUser.i64WinLoseMoney = i64DoubleWinLose
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64DoubleWinLose, 0)
			if !pUser.User.IsRobot() {
				pGame.pWinLose.i64PoolAmount += pUser.i64BetDouble
			}

			mlog.Debug("玩家[%d] 弃牌, %s[%s], 加注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64DoubleWinLose),
				Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 弃牌, %s[%s], 加注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64DoubleWinLose),
				Global.MoneyString(pUser.i64Money))
		}
	}

	//----- 傳送訊息給 client -----
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_Compensate_DoubleBet)

		pData := new(netproto.JLBPK_Compensate_DoubleBet)
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[iSeat].bOffline {
			pGame.arypUserData[iSeat].iPlayingState = JLB_UserPlayState_Compensate_DoubleBetOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}

		// 加注結算完，將加注值歸零
		//pGame.arypUserData[iSeat].i64JPMoney = 0
	}
}

// 訊息處理 - 加注結算效果完成
func (pGame *Logic) messageBetDoubleSettlementOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JLB_State_BetDoubleSettlement {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_Compensate_BetBaseOK {
		return
	}

	mlog.Debug("User[%d] 加注結算動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JLB_UserPlayState_Compensate_DoubleBetOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家加注結算效果完成")
		pGame.setGameState(JLB_State_GameOver, 0)
		pGame.gameOver()
	}
}

// 遊戲結束
func (pGame *Logic) gameOver() {
	mlog.Debug("游戏结束")
	pGame.GameLog.AddGameStep("游戏结束")

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].User != nil {
			var flowAmount, betAmount, payAmount int64
			betAmount = pGame.arypUserData[iSeat].i64BetBase + pGame.arypUserData[iSeat].i64BetDouble + pGame.arypUserData[iSeat].i64BetSide + pGame.arypUserData[iSeat].i64ZhuangBetSide
			payAmount = betAmount + pGame.arypUserData[iSeat].i64MoneyChange

			if pGame.arypUserData[iSeat].i64MoneyChange < 0 {
				flowAmount = betAmount
			} else {
				flowAmount = 0
			}

			pGame.arypUserData[iSeat].User.FlowAmount += flowAmount
			pGame.arypUserData[iSeat].User.BetAmount += betAmount
			pGame.arypUserData[iSeat].User.PayAmount += payAmount

			pGame.GameLog.Context().UserBetScore(pGame.arypUserData[iSeat].User.UserID, betAmount)
			pGame.GameLog.Context().UserPayoutScore(pGame.arypUserData[iSeat].User.UserID, payAmount+pGame.arypUserData[iSeat].i64Tax)
			pGame.GameLog.Context().UserValidScore(pGame.arypUserData[iSeat].User.UserID, flowAmount)

			pGame.GameLog.AddGameStep("用戶 %d : 本局流水 %s, 有效投注 %s, 返還金額 %s", pGame.arypUserData[iSeat].User.UserID,
				Global.MoneyString(flowAmount), Global.MoneyString(betAmount), Global.MoneyString(payAmount))
			//if pGame.arypUserData[iSeat].i64MoneyChange >= 0 {
			//	pGame.arypUserData[iSeat].User.FlowAmount = pGame.arypUserData[iSeat].i64MoneyChange
			//} else {
			//	pGame.arypUserData[iSeat].User.FlowAmount = -pGame.arypUserData[iSeat].i64MoneyChange
			//}
		}
	}

	// 結算本局輸贏次數
	if pGame.pWinLose.bDealerNotOK {
		pGame.pWinLose.GameDBData.DealerNotOKCount++
	} else {
		if pGame.pWinLose.i64PoolAmount >= 0 {
			pGame.pWinLose.GameDBData.DealerWinCount++
		} else {
			pGame.pWinLose.GameDBData.DealerLoseCount++
		}
	}

	mlog.Debug("【本局】庄家胜利[%d]，玩家胜利[%d]，庄家不成局[%d]，水量变化[%s]",
		pGame.pWinLose.GameDBData.DealerWinCount, pGame.pWinLose.GameDBData.DealerLoseCount, pGame.pWinLose.GameDBData.DealerNotOKCount,
		Global.MoneyString(pGame.pWinLose.i64PoolAmount))
	pGame.GameLog.AddGameStep("【本局】庄家胜利[%d]，玩家胜利[%d]，庄家不成局[%d]，水量变化[%s]",
		pGame.pWinLose.GameDBData.DealerWinCount, pGame.pWinLose.GameDBData.DealerLoseCount, pGame.pWinLose.GameDBData.DealerNotOKCount,
		Global.MoneyString(pGame.pWinLose.i64PoolAmount))

	// 將本局輸贏統計至房間輸贏
	pGame.SysPool.UpdateSysPool(pGame.pWinLose.i64PoolAmount)
	g_GameData.DealerWinCount += pGame.pWinLose.GameDBData.DealerWinCount
	g_GameData.DealerNotOKCount += pGame.pWinLose.GameDBData.DealerNotOKCount
	g_GameData.DealerLoseCount += pGame.pWinLose.GameDBData.DealerLoseCount
	g_GameData.GameRound++

	// g_pSysPoolWinLose.i64PoolAmount += pGame.pWinLose.i64PoolAmount
	// g_pSysPoolWinLose.GameDBData.DealerLoseCount += pGame.pWinLose.GameDBData.DealerLoseCount
	// g_pSysPoolWinLose.GameDBData.DealerNotOKCount += pGame.pWinLose.GameDBData.DealerNotOKCount
	// g_pSysPoolWinLose.GameDBData.DealerWinCount += pGame.pWinLose.GameDBData.DealerWinCount
	// g_pSysPoolWinLose.GameDBData.GameRound++
	desc := fmt.Sprintf("判定局数[%d]，目前权重表[%s]。【結算後的勝利局數】游戏局数[%d]，庄家胜利[%d]，玩家胜利[%d]，庄家不成局[%d]，水线[%s]，赢钱线[%s]",
		g_iRoundLimit,
		g_map_strDealerWeightNo[g_GameData.WeightTableNoNow],
		g_GameData.GameRound,
		g_GameData.DealerWinCount,
		g_GameData.DealerLoseCount,
		g_GameData.DealerNotOKCount,
		Global.MoneyString(pGame.SysPool.GetPool()),
		Global.MoneyString(pGame.SysPool.GetWinLine()))
	mlog.Debug(desc)
	pGame.GameLog.AddGameStep(desc)

	pGame.changeDealerWeightTable()
	//iTableNo, strLog := changeDealerWeightTable(g_pSysPoolWinLose, g_iRoundLimit, g_pSysPoolWinLose.GameDBData.WeightTableNoNow, g_i64WinLine)
	//if iTableNo >= DealerWeightNo_A && iTableNo < DealerWeightNo_Count {
	//	g_pSysPoolWinLose.GameDBData.WeightTableNoNow = iTableNo
	//}
	//if strLog != "" {
	//	mlog.Debug("%s", strLog)
	//	pGame.GameLog.AddGameStep("%s", strLog)
	//}

	//pGame.UpdateSysPoolWinLose()

	mapi64ScoreRet := make(map[int32]int64)
	mapi64CouponRet := make(map[int32]int64)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.User != nil {
			mapi64ScoreRet[pUser.i32UserID] = pUser.i64MoneyChange
			pGame.GameLog.Context().UserUpdateScore(pUser.i32UserID, pUser.i64MoneyChange)
			pGame.GameLog.Context().UserTaxAmount(pUser.i32UserID, pUser.i64Tax)
			mapi64CouponRet[pUser.i32UserID] = 0

			// 檢查玩家是否為監控玩家，且是否已完成監控
			pUser.i64AccumulatedPoint -= pUser.i64MoneyChange
			pGame.checkIfVIPUserDone(pUser)
		}
	}

	pGame.GameLog.GameEnd()

	pGame.GameEndFunc()

	iUserCount := 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			pGame.arypUserData[iSeat].Init()
		} else {
			if pGameUser.IsOffLine() {
				pGame.UserLeave(pGameUser)
				pGame.arypUserData[iSeat].Init()
			} else {
				i32LoginMoney := pGame.GetTable().Room.LoginMoney
				if pGame.arypUserData[iSeat].i64Money < int64(i32LoginMoney) {
					pGame.LogicBase.KickUserFunc(pGameUser.UserID, "您的金币不足!请补充金币后再进行游戏.")
					pGame.arypUserData[iSeat].Init()
				} else {
					pGame.arypUserData[iSeat].ResetGameData()
					pGame.arypUserData[iSeat].iUserState = JLB_UserState_Normal
					iUserCount++
				}
			}
		}
	}
	pGame.iUserCount = iUserCount
	mlog.Debug("gameOver, 玩家人數：%d", pGame.iUserCount)

	pGame.initData()
	pGame.sendAllUserInfo()
}

// 總結算 - ［海盜船］流程
func (pGame *Logic) settlement() {
	// 未加注玩家棄牌
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.iUserState == JLB_UserState_Normal &&
			pUser.i64BetDouble <= 0 {
			pUser.iUserState = JLB_UserState_Fold
			pUser.i64MoneyChange += -pUser.i64BetBase
			pGame.sendUserFoldCard(pUser)
		}
	}

	//----- 邊注結算 -----
	mlog.Debug("開牌，邊注結算")
	pGame.GameLog.AddGameStep("开牌，边注结算")

	// 計算彩金
	mlog.Debug("原彩金金额 %s", Global.MoneyString(g_i64Jackpot))
	pGame.GameLog.AddGameStep("原彩金金额 %s", Global.MoneyString(g_i64Jackpot))

	allJPMoney := int64(0)
	baseJackpotGap := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		if pUser.iUserState == JLB_UserState_Normal || pUser.iUserState == JLB_UserState_Fold {
			mlog.Debug("玩家[%d] 原来金钱 %s", pUser.i32SeatNo, Global.MoneyString(pUser.i64Money))

			bGetJP, i64JPMoney := pGame.checkIfGetJackpot(iSeat)
			pUser.bGetJackpot = bGetJP
			pUser.i64JPMoney = i64JPMoney
			g_i64Jackpot -= pUser.i64JPMoney
			allJPMoney += pUser.i64JPMoney
			pGame.sendUpdateJackpot() // 更新獎池金額

			i64MoneyChange := pUser.i64JPMoney - pUser.i64BetSide
			pUser.i64WinLoseMoney = pUser.i64JPMoney - pUser.i64BetSide
			pUser.i64TotalWinLose += pUser.i64JPMoney
			pUser.i64Money += pUser.i64JPMoney
			pUser.i64MoneyChange += i64MoneyChange

			// 寫分
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)

			str := ""
			switch pUser.iUserState {
			case JLB_UserState_Fold:
				str = fmt.Sprintf("玩家[%d] 弃牌, 边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64BetSide),
					pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
			default:
				str = fmt.Sprintf("玩家[%d] 边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64BetSide),
					pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
			}

			if !bGetJP {
				str += fmt.Sprintf("未拉中彩金, 输赢 %s, 结算后金钱 %s", Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
			} else {
				str += fmt.Sprintf("拉中彩金 %s, 输赢 %s, 结算后金钱 %s", Global.MoneyString(pUser.i64JPMoney),
					Global.MoneyString(pUser.i64WinLoseMoney), Global.MoneyString(pUser.i64Money))
			}

			mlog.Debug("%s", str)
			mlog.Debug("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))
			pGame.GameLog.AddGameStep("%s", str)
			pGame.GameLog.AddGameStep("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))

			if g_i64Jackpot < g_i64JackpotBase {
				baseJackpotGap += g_i64JackpotBase - g_i64Jackpot
				g_i64Jackpot = g_i64JackpotBase
				pGame.sendUpdateJackpot() // 更新獎池金額
				mlog.Debug("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
				pGame.GameLog.AddGameStep("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
			}
		}
	}

	// 收邊注
	i64BetSide := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserType != gamelogic.UserTypeRebot {
			i64BetSide += pGame.arypUserData[iSeat].i64BetSide
		}
	}
	mlog.Debug("原彩金金额 %s，共计收取边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64BetSide))
	pGame.GameLog.AddGameStep("原彩金金额 %s，共计收取边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64BetSide))

	g_i64Jackpot += i64BetSide

	// 莊家邊注彩池計算
	i64ZhuangBetSide := int64(0)
	if g_iZhuangBetSide == ZhuangBetSideEnable {
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.arypUserData[iSeat]
			if pUser.iUserState == JLB_UserState_Normal || pUser.iUserState == JLB_UserState_Fold {
				mlog.Debug("玩家[%d] 原来金钱 %s", pUser.i32SeatNo, Global.MoneyString(pUser.i64Money))

				bGetJP, i64JPMoney := pGame.checkIfGetZhuangJackpot(iSeat)
				pUser.bGetZhuangJackpot = bGetJP
				pUser.i64ZhuangJPMoney = i64JPMoney
				g_i64Jackpot -= pUser.i64ZhuangJPMoney
				allJPMoney += pUser.i64ZhuangJPMoney
				pGame.sendUpdateJackpot() // 更新獎池金額

				i64MoneyChange := pUser.i64ZhuangJPMoney - pUser.i64ZhuangBetSide
				pUser.i64WinLoseMoney = pUser.i64ZhuangJPMoney - pUser.i64ZhuangBetSide
				pUser.i64TotalWinLose += pUser.i64ZhuangJPMoney
				pUser.i64Money += pUser.i64ZhuangJPMoney
				pUser.i64MoneyChange += i64MoneyChange

				// 寫分
				pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)

				str := ""
				switch pUser.iUserState {
				case JLB_UserState_Fold:
					str = fmt.Sprintf("玩家[%d] 弃牌, 莊家边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64ZhuangBetSide),
						pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
				default:
					str = fmt.Sprintf("玩家[%d] 莊家边注 %s, %s[%s], ", pUser.i32UserID, Global.MoneyString(pUser.i64ZhuangBetSide),
						pGame.arypHandCards[iSeat].String(), pokerset.PatternString(pGame.arypHandCards[iSeat].Pattern(false)))
				}

				if !bGetJP {
					str += fmt.Sprintf("未拉中莊家彩金, 输赢 %s, 结算后金钱 %s", Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
				} else {
					str += fmt.Sprintf("拉中莊家彩金 %s, 输赢 %s, 结算后金钱 %s", Global.MoneyString(pUser.i64ZhuangJPMoney),
						Global.MoneyString(pUser.i64WinLoseMoney), Global.MoneyString(pUser.i64Money))
				}

				mlog.Debug("%s", str)
				mlog.Debug("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))
				pGame.GameLog.AddGameStep("%s", str)
				pGame.GameLog.AddGameStep("奖池剩余彩金 %s", Global.MoneyString(g_i64Jackpot))

				if g_i64Jackpot < g_i64JackpotBase {
					baseJackpotGap += g_i64JackpotBase - g_i64Jackpot
					g_i64Jackpot = g_i64JackpotBase
					pGame.sendUpdateJackpot() // 更新獎池金額
					mlog.Debug("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
					pGame.GameLog.AddGameStep("彩池金额不足基准值 %s, 补足至基准值", Global.MoneyString(g_i64JackpotBase))
				}
			}
		}

		// 收邊注
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.arypUserData[iSeat].iUserType != gamelogic.UserTypeRebot {
				i64ZhuangBetSide += pGame.arypUserData[iSeat].i64ZhuangBetSide
			}
		}
		mlog.Debug("原彩金金额 %s，共计收取莊家边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64ZhuangBetSide))
		pGame.GameLog.AddGameStep("原彩金金额 %s，共计收取莊家边注 %s", Global.MoneyString(g_i64Jackpot), Global.MoneyString(i64ZhuangBetSide))

		g_i64Jackpot += i64ZhuangBetSide
	}

	pGame.sendUpdateJackpot() // 更新獎池金額

	// 更新DB彩金池
	i64AllBetSide := i64BetSide
	if g_iZhuangBetSide == ZhuangBetSideEnable {
		i64AllBetSide += i64ZhuangBetSide
	}
	pGame.UpdateJackpotPool(0, i64AllBetSide, allJPMoney)

	mlog.Debug("收取边注后彩金金额 %s", Global.MoneyString(g_i64Jackpot))
	pGame.GameLog.AddGameStep("收取边注后彩金金额 %s", Global.MoneyString(g_i64Jackpot))

	//----- 底注結算 -----
	mlog.Debug("赔付，底注结算")
	pGame.GameLog.AddGameStep("赔付，底注结算")

	pDealerSet := pGame.pDealerHandCards.GetPokerCardSet()
	// 莊家是否成局
	mlog.Debug("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))
	pGame.GameLog.AddGameStep("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))

	i64BetBaseWinLose := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		switch pUser.iUserState {
		case JLB_UserState_Normal: // 正常
			pPlayerSet := pGame.arypHandCards[iSeat].GetPokerCardSet()
			if !pGame.bDealerIsGood {
				// 莊家不成局
				i64BetBaseWinLose = pUser.i64BetBase * int64(1+JLBPK_BetBasePayRate)
				pUser.i64WinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
				pUser.i64BetBaseWinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
				pUser.i64TotalWinLose += i64BetBaseWinLose
				pUser.i64Money += i64BetBaseWinLose
				pUser.iWinState = JLB_WinState_Win

				pGame.pWinLose.bDealerNotOK = true
				if !pUser.User.IsRobot() {
					pGame.pWinLose.i64PoolAmount -= pUser.i64BetBase * int64(JLBPK_BetBasePayRate)
				}
			} else {
				iCmpResult := pGame.pPokerSet.Compare(pDealerSet, pPlayerSet, false, false)
				switch iCmpResult {
				case pokerset.POKERCompareEven: // 平手
					i64BetBaseWinLose = pUser.i64BetBase
					pUser.iWinState = JLB_WinState_Even
				case pokerset.POKERCompareDealer: // 莊家贏
					i64BetBaseWinLose = -pUser.i64BetBase
					pUser.iWinState = JLB_WinState_Lose
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount += pUser.i64BetBase
					}
				case pokerset.POKERComparePlayer: // 玩家贏
					i64BetBaseWinLose = pUser.i64BetBase * int64(1+JLBPK_BetBasePayRate)
					pUser.iWinState = JLB_WinState_Win
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount -= pUser.i64BetBase * int64(JLBPK_BetBasePayRate)
					}
				}

				if iCmpResult == pokerset.POKERCompareDealer {
					pUser.i64WinLoseMoney = i64BetBaseWinLose
					pUser.i64BetBaseWinLoseMoney = i64BetBaseWinLose
					pUser.i64TotalBet -= pUser.i64BetBase
				} else {
					pUser.i64WinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
					pUser.i64BetBaseWinLoseMoney = i64BetBaseWinLose - pUser.i64BetBase
					pUser.i64Money += i64BetBaseWinLose
					pUser.i64TotalWinLose += i64BetBaseWinLose
				}
			}

			// 寫分
			i64WinLose := i64BetBaseWinLose
			if i64BetBaseWinLose > 0 {
				i64WinLose = i64BetBaseWinLose - pUser.i64BetBase
			}
			i64MoneyChange, i64Tax := pGame.GetUserWinAndTax(i64WinLose)
			pUser.i64TotalWinLose -= i64Tax
			pUser.i64Money -= i64Tax
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, i64Tax)

			mlog.Debug("玩家[%d], %s[%s], 底注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 底注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
		case JLB_UserState_Fold: // 棄牌
			i64BetBaseWinLose = -pUser.i64BetBase
			pUser.i64WinLoseMoney = i64BetBaseWinLose
			pUser.i64BetBaseWinLoseMoney = i64BetBaseWinLose
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64BetBaseWinLose, 0)
			pUser.iWinState = JLB_WinState_Lose
			if !pUser.User.IsRobot() {
				pGame.pWinLose.i64PoolAmount += pUser.i64BetBase
			}

			mlog.Debug("玩家[%d] 弃牌, %s[%s], 底注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64BetBaseWinLose), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 弃牌, %s[%s], 底注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetBase), Global.MoneyString(i64BetBaseWinLose), Global.MoneyString(pUser.i64Money))
		}
	}

	//----- 加注結算 -----
	mlog.Debug("賠付，加注結算")
	pGame.GameLog.AddGameStep("赔付，加注结算")

	pDealerSet = pGame.pDealerHandCards.GetPokerCardSet()
	// 莊家是否成局
	mlog.Debug("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))
	pGame.GameLog.AddGameStep("庄家：%s[%s]", pGame.pDealerHandCards.String(), pGame.pDealerHandCards.PatternString(true))

	i64DoubleWinLose := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.arypUserData[iSeat]
		switch pUser.iUserState {
		case JLB_UserState_Normal: // 正常
			pPlayerSet := pGame.arypHandCards[iSeat].GetPokerCardSet()
			if !pGame.bDealerIsGood {
				// 莊家不成局
				i64DoubleWinLose = pUser.i64BetDouble
				pUser.i64WinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
				pUser.i64BetDoubleWinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
				pUser.i64TotalWinLose += i64DoubleWinLose
				pUser.i64Money += i64DoubleWinLose

				pGame.pWinLose.bDealerNotOK = true
			} else {
				iCmpResult := pGame.pPokerSet.Compare(pDealerSet, pPlayerSet, false, false)
				switch iCmpResult {
				case pokerset.POKERCompareEven: // 平手
					i64DoubleWinLose = pUser.i64BetDouble
				case pokerset.POKERCompareDealer: // 第一組牌
					i64DoubleWinLose = -pUser.i64BetDouble
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount += pUser.i64BetDouble
					}
				case pokerset.POKERComparePlayer: // 第二組牌
					i64DoubleWinLose = pUser.i64BetDouble * (1 + pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)])
					if !pUser.User.IsRobot() {
						pGame.pWinLose.i64PoolAmount -= pUser.i64BetDouble * pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)]
					}
				}

				if iCmpResult == pokerset.POKERCompareDealer {
					pUser.i64WinLoseMoney = i64DoubleWinLose
					pUser.i64BetDoubleWinLoseMoney = i64DoubleWinLose
					pUser.i64TotalBet -= pUser.i64BetDouble
				} else {
					pUser.i64WinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
					pUser.i64BetDoubleWinLoseMoney = i64DoubleWinLose - pUser.i64BetDouble
					pUser.i64Money += i64DoubleWinLose
					pUser.i64TotalWinLose += i64DoubleWinLose
				}
			}

			// 寫分
			i64WinLose := i64DoubleWinLose
			if i64DoubleWinLose > 0 {
				i64WinLose = i64DoubleWinLose - pUser.i64BetDouble
			}
			i64MoneyChange, i64Tax := pGame.GetUserWinAndTax(i64WinLose)
			pUser.i64TotalWinLose -= i64Tax
			pUser.i64Money -= i64Tax
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, i64Tax)

			if pGame.bDealerIsGood {
				mlog.Debug("玩家[%d], %s[%s], 加注 %s, 赔注[1赔%d], 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)],
					Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 加注 %s, 赔注[1赔%d], 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), pGame.mapi64SpecialSheet[pGame.arypHandCards[iSeat].Pattern(false)],
					Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
			} else {
				mlog.Debug("玩家[%d], %s[%s], 加注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d], %s[%s], 加注 %s, 输赢 %s, 税收 %s, 结算后金钱 %s", pUser.i32UserID,
					pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
					Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64WinLose), Global.MoneyString(i64Tax), Global.MoneyString(pUser.i64Money))
			}
		case JLB_UserState_Fold: // 棄牌
			i64DoubleWinLose = -pUser.i64BetDouble
			pUser.i64WinLoseMoney = i64DoubleWinLose
			pUser.i64BetDoubleWinLoseMoney = i64DoubleWinLose
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64DoubleWinLose, 0)
			if !pUser.User.IsRobot() {
				pGame.pWinLose.i64PoolAmount += pUser.i64BetDouble
			}

			mlog.Debug("玩家[%d] 弃牌, %s[%s], 加注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64DoubleWinLose), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 弃牌, %s[%s], 加注 %s, 输赢 %s, 结算后金钱 %s", pUser.i32UserID,
				pGame.arypHandCards[iSeat].String(), pGame.arypHandCards[iSeat].PatternString(false),
				Global.MoneyString(pUser.i64BetDouble), Global.MoneyString(i64DoubleWinLose), Global.MoneyString(pUser.i64Money))
		}
	}

	// 傳送訊息給 client
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			continue
		}

		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_Settlement)

		pData := new(netproto.JLBPK_Settlement)
		// 莊家手牌
		pDealerCard := new(netproto.JLBPK_UserCards)
		for iIdx := 0; iIdx < pGame.pDealerHandCards.CardCount(); iIdx++ {
			pCard := new(netproto.JLBPK_PokerCard)
			pSCard := pGame.pDealerHandCards.GetSingleCard(iIdx, false)
			pCard.Color = proto.Int32(pSCard.I32Color)
			pCard.CardNumber = proto.Int32(pSCard.I32Number)

			pDealerCard.Cards = append(pDealerCard.Cards, pCard)
		}
		pDealerCard.Pattern = proto.Int32(pGame.pDealerHandCards.Pattern(true))
		pData.BankerCards = pDealerCard
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, iSeat)
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	}

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.arypUserData[iSeat].iUserState == JLB_UserState_EmptySeat || pGame.arypUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.arypUserData[iSeat].bOffline {
			pGame.arypUserData[iSeat].iPlayingState = DHC_UserPlayState_SettlementOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}

		// 邊注結算完，將邊注值歸零
		//pGame.arypUserData[iSeat].i64BetSide = 0
	}
}

// 訊息處理 - 總結算效果完成
func (pGame *Logic) messageSettlementOK(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != DHC_State_Settlement {
		return
	}

	if pUser.iPlayingState != JLB_UserPlayState_BetDouble {
		return
	}

	mlog.Debug("User[%d] 總結算動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = DHC_UserPlayState_SettlementOK

	for i := 0; i < pGame.iSeatCount; i++ {
		mlog.Debug("User[%d], seat[%d], UserState[%s], PlayingState[%s]", pGame.arypUserData[i].i32UserID,
			i, map_strUserState_debug[pGame.arypUserData[i].iUserState], map_strPlayingState_debug[pGame.arypUserData[i].iPlayingState])
	}

	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家總結算效果完成")
		pGame.setGameState(JLB_State_GameOver, 0)
		pGame.gameOver()
	}
}

// 玩家斷線重連
func (pGame *Logic) UserRecover(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserRecover, User[%d], Seat[%d]", pGame.TableNo(), pGameUser.UserID, pGameUser.SeatNo)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pUser != nil {
			mlog.Debug("Seat[%d], pUserData.bIsOffline[%t], pUser[%t]", iSeat, pGame.arypUserData[iSeat].bOffline, pUser.IsOffLine())
		} else {
			mlog.Debug("Seat[%d], pUserData.bIsOffline[%t], pUser nil", iSeat, pGame.arypUserData[iSeat].bOffline)
		}
	}

	iSeat := int(pGameUser.SeatNo)
	if iSeat < 0 || iSeat >= pGame.iSeatCount {
		return
	}

	pGame.DynamicGameLogicBase.UserRecover(pGameUser)

	pUser := pGame.arypUserData[iSeat]
	pUser.bOffline = false

	pGame.sendGameInfo(pGameUser)
	pGame.sendAllUserInfo()

	switch pGame.iGameState {
	case JLB_State_StartBettingAnimation: // 開始下注動畫
		// 玩家於此階段離線，未收到 Compensate_DoubleBet，所以重送
		pMsg := new(msg.Message)
		pMsg.BClassID = int32(netproto.MessageBClassID_Game)
		pMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_Compensate_DoubleBet)

		pData := new(netproto.JLBPK_Compensate_DoubleBet)
		// 所有玩家資訊
		for i := 0; i < pGame.iSeatCount; i++ {
			pUserInfo := pGame.userInfoPkg(i, int(pGameUser.SeatNo))
			pData.UserInfo = append(pData.UserInfo, pUserInfo)
		}

		pMsg.MsgData = pData
		pGame.SendMessageToUser(pGameUser, pMsg)
	case JLB_State_StartBetting: // 開始下注
		fallthrough
	case DHC_State_Settlement:
		fallthrough
	case JLB_State_BetDoubleSettlement: // 加注結算
		if pUser.i64BetBase <= 0 {
			pUser.i64BetBase += int64(pGame.aryi32Chip[0])
			pUser.i64TotalBet += pUser.i64BetBase
			pUser.i64Money -= pUser.i64BetBase
			pUser.i64TotalWinLose -= pUser.i64BetBase
			pGame.sendBetResult(pUser, JLB_BetType_Base, pUser.i64BetBase, JLB_BetResult_Auto)
		}
	}
}

// 玩家站立
func (pGame *Logic) UserFree(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}
	iSeat := int(pGameUser.SeatNo)
	if iSeat < 0 || iSeat >= pGame.iSeatCount {
		return
	}

	if pGame.iGameState >= JLB_State_StartBetting && pGame.iGameState <= JLB_State_BetDoubleSettlement {
		return
	}

	if g_iHDCMode == HDCModeEnable && pGame.iGameState == DHC_State_Settlement {
		return
	}

	pGame.arypUserData[iSeat].Init()
	pGame.sendAllUserInfo()
	defer func() {
		// check all user leave, close game ...
		var hasuser bool
		for _, user := range pGame.arypUserData {
			if user.User != nil {
				hasuser = true
				break
			}
		}

		if !hasuser {
			pGame.iGameState = JLB_State_WaitingPlayer
		}
	}()
	pGame.LogicBase.UserFree(pGameUser)
}

// 玩家離線
func (pGame *Logic) UserLeave(pGameUser *gamelogic.GameUser) bool {
	if pGameUser == nil {
		return false
	}

	mlog.Debug("Table[%d], UserLeave, User[%d], Seat[%d]", pGame.TableNo(), pGameUser.UserID, pGameUser.SeatNo)

	iSeat := int(pGameUser.SeatNo)
	if iSeat < 0 || iSeat >= pGame.iSeatCount {
		return false
	}

	pGame.arypUserData[iSeat].bOffline = true

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pUser != nil {
			mlog.Debug("Seat[%d], pUserData.bIsOffline[%t], pUser[%t]", iSeat, pGame.arypUserData[iSeat].bOffline, pUser.IsOffLine())
		} else {
			mlog.Debug("Seat[%d], pUserData.bIsOffline[%t], pUser nil", iSeat, pGame.arypUserData[iSeat].bOffline)
		}
	}

	if pGame.iGameState >= JLB_State_StartBetting && pGame.iGameState <= JLB_State_BetDoubleSettlement && pGame.arypUserData[iSeat].iUserState != JLB_UserState_Watch {
		return false
	}

	if g_iHDCMode == HDCModeEnable && pGame.iGameState == DHC_State_Settlement && pGame.arypUserData[iSeat].iUserState != JLB_UserState_Watch {
		return false
	}

	pGame.arypUserData[iSeat].Init()
	pGame.sendAllUserInfo()

	defer func() {
		pGame.iUserCount--

		// check all user leave, close game ...
		var hasuser bool
		for _, user := range pGame.arypUserData {
			if user.User != nil {
				hasuser = true
				break
			}
		}

		if !hasuser {
			pGame.initData()
		}
	}()

	return pGame.LogicBase.UserLeave(pGameUser)
}

// 玩家資訊打包資料，iSeat：此座位的玩家資訊，iSendTo：要送給哪位玩家
func (pGame *Logic) userInfoPkg(iSeat, iSendTo int) *netproto.JLBPK_UserInfo {
	if (iSeat < 0 || iSeat >= pGame.iSeatCount) || (iSendTo < 0 || iSendTo >= pGame.iSeatCount) {
		return nil
	}

	pUser := pGame.arypUserData[iSeat]

	pUserInfo := new(netproto.JLBPK_UserInfo)

	pUserInfo.UserID = proto.Int32(pUser.i32UserID)                          // 玩家ID，若是空位時，玩家ID = -1
	pUserInfo.SeatNo = proto.Int32(int32(iSeat + JLBPK_ClientStartUpSeatNo)) // 玩家座位，client 座位號從 1 開始
	pUserInfo.State = proto.Int32(int32(pUser.iUserState))                   // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
	pUserInfo.Money = proto.Int64(pUser.i64Money)                            // 玩家金錢
	pUserInfo.Bet_Base = proto.Int64(pUser.i64BetBase)                       // 玩家底注
	pUserInfo.Bet_Side = proto.Int64(pUser.i64BetSide)                       // 玩家邊注
	pUserInfo.Bet_Double = proto.Int64(pUser.i64BetDouble)                   // 玩家加注
	pUserInfo.GetJackpot = proto.Bool(pUser.bGetJackpot)                     // 玩家是否中 Jackpot
	pUserInfo.JPMoney = proto.Int64(pUser.i64JPMoney)                        // 玩家所中 Jackpot 金額
	pUserInfo.WinLoseMoney = proto.Int64(pUser.i64WinLoseMoney)              // 玩家輸贏金錢，若值為負數，則為玩家輸掉的金額
	pUserInfo.ZhuangBet_Side = proto.Int64(pUser.i64ZhuangBetSide)           // 玩家莊家邊注
	pUserInfo.ZhuangGetJackpot = proto.Bool(pUser.bGetZhuangJackpot)         // 玩家是否中 莊家Jackpot
	pUserInfo.ZhuangJPMoney = proto.Int64(pUser.i64ZhuangJPMoney)            // 玩家所中 莊家Jackpot 金額

	if pUser.User != nil {
		pUserInfo.FaceID = proto.Int32(pUser.User.FaceID % 2) // 玩家頭像ID
	} else {
		pUserInfo.FaceID = proto.Int32(0) // 玩家頭像ID
	}

	pUserInfo.BetBaseWinLoseMoney = proto.Int64(pUser.i64BetBaseWinLoseMoney)     // 玩家底注輸贏金錢，若值為負數，則為玩家輸掉的金額
	pUserInfo.BetDoubleWinLoseMoney = proto.Int64(pUser.i64BetDoubleWinLoseMoney) // 玩家加注輸贏金錢，若值為負數，則為玩家輸掉的金額

	// 玩家手牌
	pHandCard := new(netproto.JLBPK_UserCards)
	for iIdx := 0; iIdx < pGame.arypHandCards[iSeat].CardCount(); iIdx++ {
		pCard := new(netproto.JLBPK_PokerCard)
		if pGame.iGameState >= JLB_State_BetSideSettlement {
			// 結算階段
			if pUser.iUserState == JLB_UserState_Normal ||
				pUser.iUserState == JLB_UserState_Fold {
				pSCard := pGame.arypHandCards[iSeat].GetSingleCard(iIdx, false)
				pCard.Color = proto.Int32(pSCard.I32Color)
				pCard.CardNumber = proto.Int32(pSCard.I32Number)
			} else {
				pCard.Color = proto.Int32(pokerset.POKERCardColorOther)
				pCard.CardNumber = proto.Int32(pokerset.POKERCardNoHide)
			}
		} else {
			if iSeat != iSendTo {
				pCard.Color = proto.Int32(pokerset.POKERCardColorOther)
				pCard.CardNumber = proto.Int32(pokerset.POKERCardNoHide)
			} else {
				pSCard := pGame.arypHandCards[iSeat].GetSingleCard(iIdx, false)
				pCard.Color = proto.Int32(pSCard.I32Color)
				pCard.CardNumber = proto.Int32(pSCard.I32Number)
			}
		}
		pHandCard.Cards = append(pHandCard.Cards, pCard)
	}
	pHandCard.Pattern = proto.Int32(pGame.arypHandCards[iSeat].Pattern(false))
	pUserInfo.UserCards = pHandCard
	// 玩家總輸贏金錢
	if pUser.i64TotalWinLose > 0 {
		pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose + pUser.i64TotalBet)
	} else {
		if !pGame.bDealerIsGood {
			if pUser.iUserState == JLB_UserState_Fold {
				pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose)
			} else {
				pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose + pUser.i64TotalBet)
			}
		} else {
			pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose)
		}
	}
	pUserInfo.Winner = proto.Int32(int32(pUser.iWinState))  // 贏家
	pUserInfo.MaxBetBase = proto.Int64(pUser.i64MaxBetBase) // 可以下的最大底注值

	return pUserInfo
}

// 設定遊戲狀態
func (pGame *Logic) setGameState(iGameState int, i64Timeout int64) {
	if (iGameState < JLB_State_WaitingPlayer || iGameState > JLB_State_GameOver) && iGameState != DHC_State_Settlement {
		return
	}

	i64Time := int64(0)
	if i64Timeout < 0 {
		i64Time = 0
	} else {
		i64Time = i64Timeout
	}

	pGame.iGameState = iGameState
	pGame.i64Timeout = i64Time
	pGame.i64LastMessageTime = pGame.i64TimeDuration

	mlog.Debug("SetGameState[%s], timeout[%d], duration[%d], lastmsg[%d]", map_strGameState_debug[pGame.iGameState],
		pGame.i64Timeout, pGame.i64TimeDuration, pGame.i64LastMessageTime)
}

// 依 UserID 取得 UserData
func (pGame *Logic) GetUserDataByUserID(i32UserID int32) *JlbpkUser {
	iSeatNo := -1

	for i := 0; i < pGame.iSeatCount; i++ {
		if i32UserID == pGame.arypUserData[i].i32UserID {
			iSeatNo = i
			break
		}
	}

	if iSeatNo >= 0 && iSeatNo < pGame.iSeatCount {
		return pGame.arypUserData[iSeatNo]
	}

	return nil
}

// 檢查是否開啟換牌模式
func (pGame *Logic) checkIfChangeToCheatMode() bool {
	//bChanged := false
	//mlog.Warn("checkIfChangeToCheatMode()")
	//
	//// 檢查莊家目標金額
	//if g_i64DealerTargetPoint > 0 && g_pSysPoolWinLose.i64DealerAccumulatedPoint < g_i64DealerTargetPoint {
	//	pGame.bDealerCheat = true
	//	bChanged = true
	//} else {
	//	pGame.bDealerCheat = false
	//}
	//mlog.Debug("檢查莊家目標金額[%t]", bChanged)
	//
	//// 檢查是否有玩家被監控
	//for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
	//	pUser := pGame.arypUserData[iSeat]
	//	if pUser.User.IsRobot() || pUser.bIsVIPFinished {
	//		continue
	//	}
	//
	//	if pUser.i64TargetPoint > 0 && pUser.i64AccumulatedPoint < pUser.i64TargetPoint {
	//		bChanged = true
	//		mlog.Debug("監控中玩家, UserID[%d], 權重[%d], 目標金額[%d], 已累計金額[%d]", pUser.i32UserID,
	//			pUser.i32Weight, pUser.i64TargetPoint, pUser.i64AccumulatedPoint)
	//	}
	//}
	//
	//if bChanged {
	//	iRnd := util.RndInt(0,100)
	//	mlog.Debug("CheatRate[%d], iRnd[%d]", g_iCheatRate, iRnd)
	//	if iRnd < g_iCheatRate {
	//		mlog.Debug("checkIfChangeToCheatMode() return true")
	//		return true
	//	}
	//}

	return false
}

// 檢查監控玩家是否已完成監控，並更新 DB 資料
func (pGame *Logic) checkIfVIPUserDone(pUser *JlbpkUser) {
	if pUser == nil {
		return
	}

	//pHandler := pGame.pSuperToolService.GetJlbpkSuperToolHandler()
	//if ok, _ := util.InArray(pUser.i32UserID, pHandler.Proiority); ok {
	//	if pVIPUser := pHandler.GetVipUser(pUser.i32UserID); pVIPUser != nil {
	//		if pVIPUser.GetFinished() == 0 { // 監控尚未完成
	//			mlog.Debug("玩家[%d]在监控中，权重[%d]，目标金额[%s]，已累计金额[%s]", pUser.i32UserID, pUser.i32Weight,
	//				Global.MoneyString(pUser.i64TargetPoint), Global.MoneyString(pUser.i64AccumulatedPoint))
	//			pGame.GameLog.AddGameStep("玩家[%d]在监控中，权重[%d]，目标金额[%s]，已累计金额[%s]", pUser.i32UserID, pUser.i32Weight,
	//				Global.MoneyString(pUser.i64TargetPoint), Global.MoneyString(pUser.i64AccumulatedPoint))
	//
	//			if pUser.i64AccumulatedPoint >= pUser.i64TargetPoint {
	//				// 監控完成，恢復預設值
	//				pUser.bIsVIPFinished = true
	//				pUser.i32Weight = g_i32PlayerWeightDefault
	//				pUser.i64TargetPoint = 0
	//				pUser.i64AccumulatedPoint = 0
	//
	//				mlog.Debug("玩家[%d]监控已完成", pUser.i32UserID)
	//				pGame.GameLog.AddGameStep("玩家[%d]监控已完成", pUser.i32UserID)
	//			}
	//
	//			// 更新 DB 資料
	//			if pUser.bIsVIPFinished {
	//				pVIPUser.Finished = proto.Int32(1)
	//			} else {
	//				pVIPUser.Finished = proto.Int32(0)
	//			}
	//			pVIPUser.Weight = proto.Int32(pUser.i32Weight)
	//			pVIPUser.TargetPoint = proto.Int64(pUser.i64TargetPoint)
	//			pVIPUser.AccumulatedPoint = proto.Int64(pUser.i64AccumulatedPoint)
	//
	//			pHandler.SaveVipUser(pUser.i32UserID)
	//		}
	//	}
	//}
}
