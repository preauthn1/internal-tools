package jdesydLogic

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	Global "example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

//---------------------------------------------------------------------------
type JdesydLogic struct {
	*gamelogic.DynamicGameLogicBase
	pGamePoker            *Pokerset.PokerCardSet   // 系統牌組
	iSeatCount            int                      // 座位數
	pDealerCard           *CardSetInfo             // 莊家手牌資料
	ary_pUserData         []*JdesydUser            // 玩家資料
	iGameState            int                      // 遊戲狀態
	i64TimeDuration       int64                    // Ticker 計數值
	i64Timeout            int64                    // 各階段 timeout
	i64LastMessageTime    int64                    // 上一次訊息時間
	map_i64Timeout        map[int]int64            // timeout 時間設定
	ary_strPokerManual    []string                 // 遊戲說明
	ary_i32Chip           []int32                  // 籌碼面額
	uiAllPlayerReady      uint                     // 各玩家均已 Ready，此值設定後就不再變更
	uiPlayerReady         uint                     // 標示各玩家是否 Ready
	iUserCount            int                      // 遊戲中玩家人數
	bIsLastGame           bool                     // 是否為最後一局
	iWhosTurn             int32                    // 輪到哪個玩家動作，如果值為 -1 表示所有玩家都可以動作
	iCutCardsUser         int32                    // 要切牌的玩家
	iLastGameCardIndex    int                      // 當牌盒裡撲克牌張數剩餘 iLastGameCardIndex 時，則判定為最後一局
	bTakeInsuranceOrNot   bool                     // 是否要買保險
	bHasUserTakeInsurance bool                     // 有玩家買保險
	iDealType             int                      // 發牌方式
	pTestDealerCards      *Pokerset.PokerCardSet   // 測試用莊家手牌
	ary_pTestUserCards    []*Pokerset.PokerCardSet // 測試用玩家手牌
	pTempDealerCardSet    *Pokerset.PokerCardSet   // 指定發牌用莊家牌盒
	ary_pTempUserCardSet  []*Pokerset.PokerCardSet // 指定發牌用玩家牌盒
	pRobotSetting         *RobotSetting            // 機器人設定
	bShuffleEveryGame     bool                     // 每一局發牌/切牌，測試用
	bRobotCutCards        bool                     // 機器人切牌
	pWinLose              *SysPoolWinLose          // 本局輸贏結果統計
}

//---------------------------------------------------------------------------
func NewJdesydLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, configs *config.Group) *JdesydLogic {
	pGame := new(JdesydLogic)
	pGame.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, configs)

	pGame.i64TimeDuration = 0
	tfManager.NewMilliSecTicker(100, func() {
		pGame.i64TimeDuration++
		pGame.gameTickerFunc()
	})

	pGame.setGameState(JD21_State_WaitingPlayer, 0)

	pGame.map_i64Timeout = make(map[int]int64)
	pGame.pRobotSetting = new(RobotSetting)
	pGame.iDealType = JD21_DealType_Natural
	pGame.bShuffleEveryGame = false

	g_pSysPoolWinLose = new(SysPoolWinLose)

	pGame.getGameConfig(configs)
	pGame.pWinLose = new(SysPoolWinLose)

	pGame.pGamePoker = Pokerset.NewPokerCardSet()
	pGame.pGamePoker.NewSet(JD21_PokerSetCount, 0)
	pGame.pGamePoker.Shuffle()

	pGame.pDealerCard = NewCardSetInfo(JD21_IndexBase)
	pGame.pTempDealerCardSet = Pokerset.NewPokerCardSet()

	pGame.ary_pUserData = make([]*JdesydUser, pGame.iSeatCount)
	pGame.ary_pTempUserCardSet = make([]*Pokerset.PokerCardSet, pGame.iSeatCount)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.ary_pUserData[iSeat] = NewJdesydUser()

		pGame.ary_pTempUserCardSet[iSeat] = Pokerset.NewPokerCardSet()
	}

	// 各玩家均已 Ready，此值設定後就不再變更
	uiReady := uint(1)
	for i := 0; i < pGame.iSeatCount; i++ {
		uiReady *= 2
	}
	pGame.uiAllPlayerReady = uiReady - 1

	// 遊戲說明
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "胜利：取得最接近「21点」的点数，而不能超过21点")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "分牌：玩家将与原注码金额相等的筹码加注到桌上，并将首两张牌拆分成两手。首两张牌必须「点数价值相符」才能拆分（例如一对8；或一张K和一张Q等）。在已拆分的牌中，一张A和一张价值10点的牌等于21点的点数，但并非是黑杰克。仅能击败20点或以下的庄家。另外，拆分一对A之后，玩家只能为每一张A多拿一张牌。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "保险（INSURE）：庄家明牌拿到A时，玩家可以买「保险」，保险金额为押注额的一半。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "如庄家拿到BlackJack，则可获得保险金额的「2倍赔偿」。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "如庄家不是BlackJack，则失去保险金。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "庄家拥有BlackJack，玩家将失去赌金，游戏在这一手即结束，除非玩家也有BlackJack，此时便为平手（Push）。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "如果庄家没有BlackJack，游戏照常进行。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "点数计算：K、Q、J及10点牌各算作10点；")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "A可算作1或11点")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "其他牌（2至9）则依其面值计算。")
	pGame.ary_strPokerManual = append(pGame.ary_strPokerManual, "庄补牌规则：若庄家手持总点数为「16或以下」，庄家必须拿牌；若庄家手持的总点数为「17或以上」，庄家则不能再拿牌。")

	pGame.bIsLastGame = true
	pGame.iLastGameCardIndex = JD21_MinCutCardIndex

	pGame.initData()

	pGame.RegisterBeforeGameStartFunc(func() {
	})
	pGame.RegisterAfterGameStartFunc(func() {
	})
	pGame.RegisterBeforeGameEndFunc(func() {
	})
	pGame.RegisterAfterGameEndFunc(func() {
	})

	return pGame
}

func (pGame *JdesydLogic) initData() {
	mlog.Debug("初始化遊戲資料")
	/**
	主要解決卡桌問題, 需要在用戶都離開的時候初始化遊戲資料
	若初始化過程中表示遊戲中則取消該局遊戲.
	*/
	if pGame.IsGaming() {
		mlog.Debug("initData 遊戲初始化時遊戲狀態為遊戲中")
		pGame.GameCancelFunc("初始化遊戲資料")
	}

	for i := 0; i < pGame.iSeatCount; i++ {
		pGame.ary_pUserData[i].ResetGameData()
	}
	pGame.pDealerCard.ClearData()
	pGame.pDealerCard.bShowCards = false // 遊戲一開始，不顯示莊家全部的牌

	pGame.pTempDealerCardSet.SetCardSet(pGame.pTestDealerCards.CodeString())

	pGame.pWinLose.i64PoolDifference = 0
	pGame.pWinLose.Reset(true)

	pGame.iUserCount = 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			pGame.ary_pUserData[iSeat].Init()
		} else {
			if pGame.ary_pUserData[iSeat].iUserState != JD21_UserState_EmptySeat {
				pGame.ary_pUserData[iSeat].iUserState = JD21_UserState_Normal
			}
			pGame.UserReady(pGameUser)
			pGame.iUserCount++
		}

		pGame.ary_pTempUserCardSet[iSeat].SetCardSet(pGame.ary_pTestUserCards[iSeat].CodeString())
	}

	if pGame.iUserCount <= 0 {
		pGame.setGameState(JD21_State_WaitingPlayer, 0)

		mlog.Debug("等待游戏开始")
		pGame.GameLog.AddGameStep("等待游戏开始")
	} else {
		pGame.setGameState(JD21_State_WaitingStart, pGame.map_i64Timeout[JD21_Timeout_WaitingStart])
		pGame.waitingGameBegin()
	}
	pGame.i64LastMessageTime = pGame.i64TimeDuration

	pGame.uiPlayerReady = 0
	pGame.iWhosTurn = -1
	pGame.iCutCardsUser = -2
	pGame.bTakeInsuranceOrNot = false
	pGame.bHasUserTakeInsurance = false
	pGame.bRobotCutCards = false
}

func (pGame *JdesydLogic) GameMessage(pGameUser *gamelogic.GameUser, pMsg *msg.Message) {
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
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_Bet): // C2S 玩家下注
		pGame.message_Bet(pUser, pMsg)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_StopBettingAnimationOK): // C2S 玩家停止下注動畫播放完成
		pGame.message_StopBettingAnimationOK(pUser)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_CutCards): // C2S 玩家切牌
		pGame.message_CutCards(pUser, pMsg)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_CutCardsAnimationOK): // C2S 切牌動畫播放完成
		pGame.message_CutCardsAnimationOK(pUser)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_DealAnimationOK): // C2S 發牌動畫播放完成
		pGame.message_DealAnimationOK(pUser)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_TakeInsurance): // C2S 玩家買保險
		pGame.message_TakeInsurance(pUser, pMsg)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_DealerPeeksCardsDone): // C2S 莊家瞇牌動畫播放完成
		pGame.message_DealerPeeksCardsDone(pUser)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_Settlement_Insurance_OK): // C2S 保險結算動畫播放完成
		pGame.message_Settlement_Insurance_OK(pUser)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_UserPlaying): // C2S 玩家動作
		pGame.message_UserPlaying(pUser, pMsg)
	case int32(netproto.JDESYD_GameMessageClassID_JDESYD_Settlement_Dealer_OK): // C2S 莊家開牌/結算動畫播放完成
		pGame.message_Settlement_Dealer_OK(pUser)
	}
}

func (pGame *JdesydLogic) PiSaiMessage(pUser *gamelogic.GameUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}
}

// 玩家進入遊戲
func (pGame *JdesydLogic) UserEnter(pGameUser *gamelogic.GameUser, i32TableNo int32, i32SeatNo int32) {
	if pGameUser == nil {
		return
	}

	// 自 DB 取得遊戲輸贏結果
	if g_bEnablePoolControl {
		if !g_bGetWinLoseData {
			pGame.GetSysPoolWinLose()
		}
	}

	mlog.Debug("TableNo[%d], UserEnter, User[%d], Seat[%d]", pGame.TableNo(), pGameUser.UserID, pGameUser.SeatNo)

	if pGameUser.IsRobot() {
		if pGame.ary_pUserData[i32SeatNo].i32UserID != pGameUser.UserID {
			for i := 0; i < len(pGame.ary_pUserData); i++ {
				if pGame.ary_pUserData[i].i32UserID == pGameUser.UserID {
					i32SeatNo = int32(i)
					break
				}
			}
		}
	}

	pGame.DynamicGameLogicBase.UserEnter(pGameUser, i32TableNo, i32SeatNo)
	pGame.DynamicGameLogicBase.UserReady(pGameUser)
	pGame.AddUser(pGameUser)
}

// 計時器處理程序
func (pGame *JdesydLogic) gameTickerFunc() {
	if pGame.i64Timeout > 0 && pGame.i64TimeDuration-pGame.i64LastMessageTime > pGame.i64Timeout {
		switch pGame.iGameState {
		case JD21_State_WaitingStart: // 等待遊戲開始
			pGame.GameStartFunc()
			pGame.setGameState(JD21_State_Betting, pGame.map_i64Timeout[JD21_Timeout_Betting])
			pGame.startBetting(false, nil)
		case JD21_State_Betting: // 押底注
			pGame.setGameState(JD21_State_StopBetting, pGame.map_i64Timeout[JD21_Timeout_StopBetting])
			pGame.stopBetting(false, nil)
		case JD21_State_StopBetting: // 停止下注
			if pGame.bIsLastGame {
				// 上一局為最後一局，切牌
				pGame.setGameState(JD21_State_CutCards, pGame.map_i64Timeout[JD21_Timeout_CutCards])
				pGame.cutCards()
			} else {
				// 發牌
				pGame.setGameState(JD21_State_Deal, pGame.map_i64Timeout[JD21_Timeout_DealAnimation])
				pGame.deal()
			}
		case JD21_State_CutCards: // 切牌
			pGame.setGameState(JD21_State_CutCardsAnimation, pGame.map_i64Timeout[JD21_Timeout_CutCardsAnimation])
			pGame.sendCutCardsResult(false, nil)
		case JD21_State_CutCardsAnimation: // 切牌動畫
			pGame.setGameState(JD21_State_Deal, pGame.map_i64Timeout[JD21_Timeout_DealAnimation])
			pGame.deal()
		case JD21_State_Deal: // 發牌
			// 是否要買保險
			if pGame.bTakeInsuranceOrNot {
				pGame.setGameState(JD21_State_Insurance, pGame.map_i64Timeout[JD21_Timeout_Insurance])
				pGame.sendInsuranceOrNot(false, nil)
			} else {
				// 進入玩家操作階段
				pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
				pGame.startPlaying()
			}
		case JD21_State_Insurance: // 保險
			for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
				pUser := pGame.ary_pUserData[iSeat]
				if pUser.iUserState == JD21_UserState_Normal && pUser.iUserType != gamelogic.UserTypeRebot {
					if pUser.i64Insurance <= 0 {
						pGame.userAction_TakeInsurance(pUser, 0)
					}
				}
			}

			if pGame.bHasUserTakeInsurance {
				// 有玩家買保險，莊家瞇牌
				pGame.setGameState(JD21_State_DealerPeeksCards, pGame.map_i64Timeout[JD21_Timeout_DealerPeeksCards])
				pGame.dealerPeeksCards()
			} else {
				// 進入玩家操作階段
				pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
				pGame.startPlaying()
			}
		case JD21_State_DealerPeeksCards: // 莊家瞇牌
			pGame.setGameState(JD21_State_Settlement_Insurance, pGame.map_i64Timeout[JD21_Timeout_Settlement_Insurance])
			pGame.settlement_Insurance()
		case JD21_State_Settlement_Insurance: // 保險結算
			if pGame.pDealerCard.bIsBlackJack {
				// 莊家牌是 BlackJack，結束本局
				pGame.setGameState(JD21_State_GameOver, 0)
				pGame.gameOver()
			} else {
				// 進入玩家操作階段
				pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
				pGame.startPlaying()
			}
		case JD21_State_Playing: // 玩家操作
			// 自動停牌
			pGame.ary_pUserData[pGame.iWhosTurn].Stand()
			pGame.sendUserPlayingResult(pGame.ary_pUserData[pGame.iWhosTurn], JD21_UserAction_Stand, JD21_Playing_Success)

			mlog.Debug("玩家[%d] 操作时间到，停牌", pGame.ary_pUserData[pGame.iWhosTurn].i32UserID)
			pGame.GameLog.AddGameStep("玩家[%d] 操作时间到，停牌", pGame.ary_pUserData[pGame.iWhosTurn].i32UserID)

			pGame.setGameState(JD21_State_PlayingDelay, pGame.map_i64Timeout[JD21_Timeout_PlayingDelay])
			//pGame.checkWhosTurn()
		case JD21_State_PlayingDelay: // 玩家操作結果傳送後的延遲時間
			pGame.checkWhosTurn()
		case JD21_State_Settlement_Dealer: // 莊家開牌/結算
			pGame.setGameState(JD21_State_GameOver, 0)
			pGame.gameOver()
		}
	}

	// 機器人動作
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == JD21_UserState_Normal {
			if pUser.i64Timeout > 0 && pGame.i64TimeDuration-pUser.i64LastMessageTime > pUser.i64Timeout {
				switch pGame.iGameState {
				case JD21_State_Betting: // 押底注
					if pUser.iBetTimes < pUser.pRobot.iBetCount {
						// 可下注籌碼
						ary_i64Chip := make([]int64, 0)
						for i := 0; i < len(pGame.ary_i32Chip); i++ {
							if i < pUser.pRobot.iBetRange_Min-1 || i > pUser.pRobot.iBetRange_Max-1 {
								continue
							}

							if pUser.i64Money >= int64(pGame.ary_i32Chip[i]) {
								ary_i64Chip = append(ary_i64Chip, int64(pGame.ary_i32Chip[i]))
							}
						}

						if len(ary_i64Chip) > 0 {
							iBetIdx := util.RndInt(0, len(ary_i64Chip))
							iResult := pGame.userAction_Bet(pUser, ary_i64Chip[iBetIdx])
							if iResult == JD21_BetResult_Success {
								pGame.sendBetResult(pUser, ary_i64Chip[iBetIdx], iResult)
							}
						}

						pUser.iBetTimes++

						if pUser.iBetTimes >= pUser.pRobot.iBetCount {
							pUser.i64Timeout = 0
						} else {
							pUser.i64Timeout = 5
							pUser.i64LastMessageTime = pGame.i64TimeDuration
						}
					}
				case JD21_State_Insurance: // 保險
					iRnd := util.RndInt(0, 100)
					if iRnd < pUser.pRobot.iRate_Insurrence {
						pGame.userAction_TakeInsurance(pUser, 1)
					} else {
						pGame.userAction_TakeInsurance(pUser, 0)
					}

					pUser.iPlayingState = JD21_UserPlayState_TakeInsuranceOK
					pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
					if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
						mlog.Debug("所有玩家都已決定是否買保險")
						if pGame.bHasUserTakeInsurance {
							// 有玩家買保險，莊家瞇牌
							pGame.setGameState(JD21_State_DealerPeeksCards, pGame.map_i64Timeout[JD21_Timeout_DealerPeeksCards])
							pGame.dealerPeeksCards()
						} else {
							// 進入玩家操作階段
							pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
							pGame.startPlaying()
						}
					}

					pUser.i64Timeout = 0
				case JD21_State_Playing: // 玩家操作
					pGame.robotPlaying(pUser)

					pUser.i64Timeout = 0
				}
			}
		}
	}
}

// 讀取遊戲設定
func (pGame *JdesydLogic) getGameConfig(configs *config.Group) {
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
		pGame.map_i64Timeout[JD21_Timeout_WaitingStart] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_WaitingStart")) * 10
		pGame.map_i64Timeout[JD21_Timeout_Betting] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Betting")) * 10
		pGame.map_i64Timeout[JD21_Timeout_StopBetting] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_StopBettingAnimation")) * 10
		pGame.map_i64Timeout[JD21_Timeout_CutCards] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_CutCards")) * 10
		pGame.map_i64Timeout[JD21_Timeout_CutCardsAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_CutCardsAnimation")) * 10
		pGame.map_i64Timeout[JD21_Timeout_DealAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_DealAnimation")) * 10
		pGame.map_i64Timeout[JD21_Timeout_Insurance] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Insurance")) * 10
		pGame.map_i64Timeout[JD21_Timeout_DealerPeeksCards] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_DealerPeeksCards")) * 10
		pGame.map_i64Timeout[JD21_Timeout_Settlement_Insurance] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Settlement_Insurance")) * 10
		pGame.map_i64Timeout[JD21_Timeout_Playing] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Playing")) * 10
		pGame.map_i64Timeout[JD21_Timeout_Settlement_Dealer] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Settlement_Dealer")) * 10

		// 以下為系統自訂 timeout，非自 config 中讀取的值
		pGame.map_i64Timeout[JD21_Timeout_PlayingDelay] = int64(5)  // 延遲 0.5 秒進行下一個動作
		pGame.map_i64Timeout[JD21_Timeout_RobotCutCard] = int64(10) // 若無可切牌玩家時，由機器人自動切牌的延遲時間

		// 籌碼面額
		str = cfgs.GetConfigValue("GameSetting", "ChipValue")
		ary_strChip := strings.Split(str, ",")
		for i := 0; i < len(ary_strChip); i++ {
			i64chip, _ := strconv.ParseInt(ary_strChip[i], 10, 32)
			pGame.ary_i32Chip = append(pGame.ary_i32Chip, int32(i64chip))
		}

		// 底注額上限
		g_i64MaxTotalBetBase = int64(cfgs.GetConfigInt64("GameSetting", "MaxTotalBetBase"))

		if len(pGame.ary_i32Chip) <= 0 {
			mlog.Fatal("讀取 GameSetting 失敗：沒有籌碼資料")
		}

		// 機器人參數
		pGame.pRobotSetting.iRate_Hit = cfgs.GetConfigInt("robot", "Rate_Hit")               // 小於17點時要牌機率
		pGame.pRobotSetting.iRate_Hit_ReachSoft = cfgs.GetConfigInt("robot", "Rate_Be17")    // 大於等於17點時要牌機率
		pGame.pRobotSetting.iRate_Split = cfgs.GetConfigInt("robot", "Rate_Split")           // 可分牌時分牌機率
		pGame.pRobotSetting.iRate_Double = cfgs.GetConfigInt("robot", "Rate_Double")         // 可加倍時加倍機率
		pGame.pRobotSetting.iRate_Insurrence = cfgs.GetConfigInt("robot", "Rate_Insurrence") // 可保險時保險機率
		pGame.pRobotSetting.iRate_Surrender = cfgs.GetConfigInt("robot", "Rate_Surrender")   // 莊家明牌有10或是A時投降的機率
		// 下注參數
		str = cfgs.GetConfigValue("robot", "BetRange")
		ary_strSetting := strings.Split(str, ",")
		var i64Param int64
		i64Param, _ = strconv.ParseInt(ary_strSetting[0], 10, 32)
		pGame.pRobotSetting.iBetRange_Min = int(i64Param)
		i64Param, _ = strconv.ParseInt(ary_strSetting[1], 10, 32)
		pGame.pRobotSetting.iBetRange_Max = int(i64Param)
		i64Param, _ = strconv.ParseInt(ary_strSetting[2], 10, 32)
		pGame.pRobotSetting.iDefaultBetCount = int(i64Param)

		// 是否開啟水線控制
		iPoolControl := cfgs.GetConfigInt("SystemPool", "EnablePoolControl")
		if iPoolControl == 0 {
			g_bEnablePoolControl = false
		} else {
			g_bEnablePoolControl = true
		}

		// 水線控制
		if g_bEnablePoolControl {
			if !g_bGetSysPoolConfig {
				g_map_pDealerSupplement_Base = make(map[int32]*DealerSupplement)
				g_map_pDealerSupplement_Now = make(map[int32]*DealerSupplement)
				g_ary_pSysPoolSetting = make([]*PoolDynamicSetting, 0)

				g_i64WinLine = cfgs.GetConfigInt64("SystemPool", "WinLine")       // 贏錢線
				g_iRoundLimit = cfgs.GetConfigInt("SystemPool", "RoundLimit")     // 判定局數
				g_i64AverageBet = cfgs.GetConfigInt64("SystemPool", "AverageBet") // 平均押注
				g_iFactor = cfgs.GetConfigInt("SystemPool", "Factor")             // 係數定值

				// 取得預設莊家補牌表
				str := cfgs.GetConfigValue("DealerPoint", "Points")
				ary_strDPoints := strings.Split(str, ",")
				for i := 0; i < len(ary_strDPoints); i++ {
					i64Point, _ := strconv.ParseInt(ary_strDPoints[i], 10, 32)
					g_ary_i32DealerPoints = append(g_ary_i32DealerPoints, int32(i64Point))

					strOdds := cfgs.GetConfigValue(ary_strDPoints[i], "Odds")
					ary_strOdds := strings.Split(strOdds, ",")
					strPoints := cfgs.GetConfigValue(ary_strDPoints[i], "Points")
					ary_strPoints := strings.Split(strPoints, ",")
					if len(ary_strOdds) != len(ary_strPoints) {
						mlog.Error("庄家补牌表点数[%d]设定参数数目不相同。Odds 数目：%d，Points 数目：%d", i64Point,
							len(ary_strOdds), len(ary_strPoints))
					}

					pSupplement := new(DealerSupplement)
					iRndRange := 0
					iOddsSum := 0
					for j := 0; j < len(ary_strOdds); j++ {
						pData := new(SupplementData)
						i64Odds, _ := strconv.ParseInt(ary_strOdds[j], 10, 32)
						pData.iOdds = int(i64Odds)
						iOddsSum += pData.iOdds

						if pData.iOdds <= 0 {
							pData.iRndRange = 0
						} else {
							iRndRange += pData.iOdds
							pData.iRndRange = iRndRange
						}

						if j >= len(ary_strPoints) {
							pData.i32Points = 0
						} else {
							i64Points, _ := strconv.ParseInt(ary_strPoints[j], 10, 32)
							pData.i32Points = int32(i64Points)
						}
						pSupplement.ary_pData = append(pSupplement.ary_pData, pData)
					}
					g_map_pDealerSupplement_Base[int32(i64Point)] = pSupplement
					g_map_pDealerSupplement_Base[int32(i64Point)].iOddsSum = iOddsSum
				}

				DealerSupplement_Base_Debug()

				// 設定目前莊家補牌表
				for i32Point, pData := range g_map_pDealerSupplement_Base {
					if pData != nil {
						pSupplement := new(DealerSupplement)
						pSupplement.Copy(pData)
						g_map_pDealerSupplement_Now[i32Point] = pSupplement
						g_map_pDealerSupplement_Now[i32Point].iOddsSum = pData.iOddsSum
					}
				}

				DealerSupplement_Now_Debug()

				// 取得動態水線區間表
				str = cfgs.GetConfigValue("DynamicSetting", "XVariation")
				ary_strX := strings.Split(str, ",")
				str = cfgs.GetConfigValue("DynamicSetting", "ConditionOdds")
				ary_strConditionOdds := strings.Split(str, ",")
				if len(ary_strX) != len(ary_strConditionOdds) {
					mlog.Error("动态水线区间表 XVariation 与 ConditionOdds 参数数目不相同。XVariation 数目：%d，ConditionOdds 数目：%d",
						len(ary_strX), len(ary_strConditionOdds))
				}

				for i := 0; i < len(ary_strX); i++ {
					pData := new(PoolDynamicSetting)
					pData.iAreaNo = i + 1
					i64Value, _ := strconv.ParseInt(ary_strX[i], 10, 32)
					pData.iXVariation = int(i64Value)
					if i >= len(ary_strConditionOdds) {
						pData.iConditionOdds = 0
					} else {
						i64Value, _ := strconv.ParseInt(ary_strConditionOdds[i], 10, 32)
						pData.iConditionOdds = int(i64Value)
					}
					pData.i64AreaValue = g_i64AverageBet * int64(g_iFactor) * int64(pData.iConditionOdds) * 100

					g_ary_pSysPoolSetting = append(g_ary_pSysPoolSetting, pData)
				}

				SysPoolSettingData_Debug()

				mlog.System("水线控制功能开启")
				g_bGetSysPoolConfig = true
			}
		}
	}

	// 撲克牌發牌器
	pGame.pTestDealerCards = Pokerset.NewPokerCardSet()
	pGame.ary_pTestUserCards = make([]*Pokerset.PokerCardSet, pGame.iSeatCount)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.ary_pTestUserCards[iSeat] = Pokerset.NewPokerCardSet()
	}

	pkcfgs := configs.GetConfig("poker")
	if pkcfgs == nil {
		pGame.iDealType = JD21_DealType_Natural
	} else {
		// 發牌方式
		iDealType := pkcfgs.GetConfigInt("Poker", "DealType")
		if iDealType < JD21_DealType_Natural || iDealType >= JD21_DealType_Max {
			iDealType = JD21_DealType_Natural
		}
		pGame.iDealType = iDealType

		// 測試用莊家手牌
		strCards := pkcfgs.GetConfigValue("Poker", "Dealer")
		pGame.pTestDealerCards.SetCardSet(strCards)
		mlog.Debug("測試用莊家手牌：%s", pGame.pTestDealerCards.CodeString())

		// 測試用玩家手牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			str = fmt.Sprintf("User%d", iSeat+1)
			strCards := pkcfgs.GetConfigValue("Poker", str)
			pGame.ary_pTestUserCards[iSeat].SetCardSet(strCards)
			mlog.Debug("測試用玩家手牌 %d：%s", iSeat+1, pGame.ary_pTestUserCards[iSeat].CodeString())
		}

		// 是否設定為每一局發牌/切牌
		iShuffle := pkcfgs.GetConfigInt("CutCard", "EveryGame")
		if iShuffle == 0 {
			pGame.bShuffleEveryGame = false
		} else {
			pGame.bShuffleEveryGame = true
		}
	}
}

// 加入新玩家
func (pGame *JdesydLogic) AddUser(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	iSeatNo := int(pGameUser.SeatNo)
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return
	}

	pUser := pGame.ary_pUserData[iSeatNo]
	if pUser.iUserState == JD21_UserState_EmptySeat {
		pUser.Init()
		pUser.User = pGameUser
		pUser.i32UserID = pGameUser.UserID
		pUser.i32SeatNo = pGameUser.SeatNo
		if pGameUser.IsRobot() {
			pUser.iUserType = gamelogic.UserTypeRebot
			pUser.pRobot.Copy(pGame.pRobotSetting)
		}
		pUser.i64Money = pGameUser.Score
		pUser.ResetGameData()

		if pGame.iGameState >= JD21_State_Betting && pGame.iGameState < JD21_State_GameOver {
			pUser.iUserState = JD21_UserState_Watch
		} else {
			pUser.iUserState = JD21_UserState_Normal
		}
	} else {
		mlog.Error("AddUser is not empty, %d", pUser.i32UserID)
	}
	pGame.iUserCount++

	if pGame.iUserCount > 0 && pGame.iGameState == JD21_State_WaitingPlayer {
		pGame.setGameState(JD21_State_WaitingStart, pGame.map_i64Timeout[JD21_Timeout_WaitingStart])
		pGame.waitingGameBegin()
	}
	pGame.sendGameInfo(pGameUser)
	pGame.sendAllUserInfo()

	example.invalid("加入新玩家[%d], 遊戲狀態[%s], 座位[%d], 類別[%d], 玩家狀態[%s], 玩家總人數[%d]", pUser.i32UserID,
		g_map_strGameState_debug[pGame.iGameState], iSeatNo, pUser.iUserType,
		g_map_strUserState_debug[pUser.iUserState], pGame.iUserCount)
	pGame.GameLog.AddGameStep("玩家[%d] 入座，状态：%s", pGameUser.UserID, g_map_strUserState[pUser.iUserState])
}

// 傳送牌局資訊
func (pGame *JdesydLogic) sendGameInfo(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_GameInfo)

	pData := new(netproto.JDESYD_Msg_GameInfo)
	//停止下注的狀態, 需要另外做處理, 剩餘時間需要為0
	if pGame.iGameState == JD21_State_StopBetting {
		pData.GameState = proto.Int32(int32(pGame.toClientGameState())) // 遊戲狀態
		pData.TimeRemain = proto.Int32(0)                               // 剩餘時間
	} else {
		pData.GameState = proto.Int32(int32(pGame.toClientGameState()))                                                   // 遊戲狀態
		pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間
	}

	pData.MaxTotalBetBase = proto.Int64(g_i64MaxTotalBetBase) // 底注押注總額上限
	// 籌碼面額
	for i := 0; i < len(pGame.ary_i32Chip); i++ {
		pData.ChipValue = append(pData.ChipValue, pGame.ary_i32Chip[i])
	}
	pData.DealerCards = pGame.cardSetInfoPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}
	// 執行切牌/遊戲操作的玩家
	i32SeatNo := int32(-1)
	switch pGame.toClientGameState() {
	case JD21_CState_CutCards: // 切牌
		i32SeatNo = pGame.iCutCardsUser + int32(JD21_IndexBase)
	case JD21_CState_Playing: // 玩家操作
		i32SeatNo = pGame.iWhosTurn + int32(JD21_IndexBase)
	default:
		i32SeatNo = -1
	}
	pData.SeatNo = proto.Int32(i32SeatNo)
	pData.MinCardIndex = proto.Int32(int32(JD21_MinCutCardIndex)) // 切牌區間的最小值
	pData.MaxCardIndex = proto.Int32(int32(JD21_MaxCutCardIndex)) // 切牌區間的最大值
	// 手上第幾組牌，從 1 開始
	if i32SeatNo >= int32(JD21_IndexBase) && i32SeatNo <= int32(pGame.iSeatCount) {
		pData.CardSetIndex = proto.Int32(int32(pGame.ary_pUserData[i32SeatNo-int32(JD21_IndexBase)].CardSetUsedNow().iIndex))
	} else {
		pData.CardSetIndex = proto.Int32(int32(-1))
	}
	// 遊戲說明
	for i := 0; i < len(pGame.ary_strPokerManual); i++ {
		pData.PokerManual = append(pData.PokerManual, pGame.ary_strPokerManual[i])
	}

	pMsg.MsgData = pData
	pGame.SendMessageToUser(pGameUser, pMsg)

	mlog.Debug("sendGameInfo(), GameState[%s], ToClientState[%s]", g_map_strGameState_debug[pGame.iGameState],
		g_map_strToClientState_debug[pGame.toClientGameState()])
}

// 送給 client 的遊戲狀態
func (pGame *JdesydLogic) toClientGameState() int {
	switch pGame.iGameState {
	case JD21_State_WaitingPlayer: // 等待玩家加入
		fallthrough
	case JD21_State_WaitingStart: // 等待遊戲開始
		return JD21_CState_WaitingStart
	case JD21_State_Betting: // 押底注
		fallthrough
	case JD21_State_StopBetting: // 停止下注
		return JD21_CState_Betting
	case JD21_State_CutCards: // 切牌
		fallthrough
	case JD21_State_CutCardsAnimation: // 切牌動畫
		return JD21_CState_CutCards
	case JD21_State_Deal: // 發牌
		return JD21_CState_Deal
	case JD21_State_Insurance: // 保險
		fallthrough
	case JD21_State_DealerPeeksCards: // 莊家瞇牌
		return JD21_CState_Insurance
	case JD21_State_Settlement_Insurance: // 保險結算
		return JD21_CState_Settlement_Insurance
	case JD21_State_Playing: // 玩家操作
		fallthrough
	case JD21_State_PlayingDelay: // 玩家操作結果傳送後的延遲時間
		return JD21_CState_Playing
	case JD21_State_Settlement_Dealer: // 莊家開牌/結算
		return JD21_CState_Settlement
	case JD21_State_GameOver: // 遊戲結束
		return JD21_CState_WaitingStart
	}

	return -1
}

// 傳送玩家資訊
func (pGame *JdesydLogic) sendAllUserInfo() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_AllUserInfo)

	pData := new(netproto.JDESYD_Msg_AllUserInfo)
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)
}

// 玩家資訊打包資料，iSeat：此座位的玩家資訊，iSendTo：要送給哪位玩家
func (pGame *JdesydLogic) userInfoPkg(iSeat int) *netproto.JDESYD_Msg_UserInfo {
	if iSeat < 0 || iSeat >= pGame.iSeatCount {
		return nil
	}

	pUser := pGame.ary_pUserData[iSeat]

	pUserInfo := new(netproto.JDESYD_Msg_UserInfo)

	pUserInfo.UserID = proto.Int32(pUser.i32UserID)               // 玩家ID，若是空位時，玩家ID = -1
	pUserInfo.SeatNo = proto.Int32(int32(iSeat + JD21_IndexBase)) // 玩家座位，client 座位號從 1 開始
	pUserInfo.State = proto.Int32(int32(pUser.iUserState))        // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
	pUserInfo.Money = proto.Int64(pUser.i64Money)                 // 玩家金錢
	pUserInfo.MaxBetBase = proto.Int64(pUser.i64MaxBetBase)       // 可以下的最大底注值
	// 玩家手牌
	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		pUserInfo.HandCards = append(pUserInfo.HandCards, pGame.cardSetInfoPkg(pUser.ary_pCardSet[i]))
	}
	pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose) // 總輸贏
	pUserInfo.IsSurrendered = proto.Bool(pUser.bSurrendered)    // 是否已投降
	pUserInfo.Insurance = proto.Int64(pUser.i64Insurance)       // 玩家買的保險金額

	return pUserInfo
}

// 等待遊戲開始
func (pGame *JdesydLogic) waitingGameBegin() {
	pGame.GameLog.GameBegin()
	mlog.Debug("等待开始下注")
	pGame.GameLog.AddGameStep("等待开始下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_WaitingStart)
	pData := new(netproto.JDESYD_Msg_WaitingStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	mlog.Debug("測試用莊家手牌：%s", pGame.pTestDealerCards.CodeString())

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		mlog.Debug("測試用玩家手牌 %d：%s", iSeat+1, pGame.ary_pTestUserCards[iSeat].CodeString())
	}
}

// 手牌牌組打包資料
func (pGame *JdesydLogic) cardSetInfoPkg(pCardSetInfo *CardSetInfo) *netproto.JDESYD_Msg_CardSetInfo {
	if pCardSetInfo == nil {
		return nil
	}

	pRetInfo := new(netproto.JDESYD_Msg_CardSetInfo)

	pRetInfo.CardSetIndex = proto.Int32(int32(pCardSetInfo.iIndex)) // 手中的第 N 組牌，Index 從 1 開始
	// 牌組
	for iIdx := 0; iIdx < pCardSetInfo.pCardSet.CardCount(); iIdx++ {
		pCard := new(netproto.PokerCard)
		pSCard := pCardSetInfo.pCardSet.GetSingleCard(iIdx, false)
		if pSCard != nil {
			if pCardSetInfo.bShowCards {
				pCard.Color = proto.Int32(pSCard.I32Color)
				pCard.CardNumber = proto.Int32(pSCard.I32Number)
			} else {
				if iIdx == 0 {
					pCard.Color = proto.Int32(pSCard.I32Color)
					pCard.CardNumber = proto.Int32(pSCard.I32Number)
				} else {
					pCard.Color = proto.Int32(Pokerset.POKERCardColorOther)
					pCard.CardNumber = proto.Int32(Pokerset.POKERCardNoHide)
				}
			}
		}

		pRetInfo.CardSet = append(pRetInfo.CardSet, pCard)
	}
	pRetInfo.Point1 = proto.Int32(pCardSetInfo.i32Point_1)       // 第一個點數
	pRetInfo.Point2 = proto.Int32(pCardSetInfo.i32Point_2)       // 第二個點數
	pRetInfo.IsBlackJack = proto.Bool(pCardSetInfo.bIsBlackJack) // 是否為 BlackJack
	pRetInfo.IsBusted = proto.Bool(pCardSetInfo.bIsBusted)       // 是否爆牌
	pRetInfo.BetBase = proto.Int64(pCardSetInfo.i64BetBase)      // 底注
	pRetInfo.BetDouble = proto.Int64(pCardSetInfo.i64BetDouble)  // 加倍的金額
	pRetInfo.CanSplit = proto.Bool(pCardSetInfo.bCanSplit)       // 是否可以分牌
	pRetInfo.CanDouble = proto.Bool(pCardSetInfo.bCanDouble)     // 是否可以加倍
	pRetInfo.WinLose = proto.Int64(pCardSetInfo.i64WinLose)      // 此組牌組的輸贏

	return pRetInfo
}

// 開始下注
func (pGame *JdesydLogic) startBetting(bUserRecover bool, pUser *gamelogic.GameUser) {
	mlog.Debug("開始下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_StartBetting)
	pData := new(netproto.JDESYD_Msg_WaitingStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData
	if bUserRecover {
		if pUser != nil {
			pGame.SendMessageToUser(pUser, pMsg)
		}
	} else {
		pGame.BoardcastMessage(pMsg)
	}

	pGame.GameLog.AddGameStep("开始下注")

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal {
			pUser.iPlayingState = JD21_UserPlayState_Betting
			if pUser.TotalBetBase() <= 0 {
				// 自動幫玩家下底注
				mlog.Debug("玩家[%d] 金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
				pUser.Bet(0, int64(pGame.ary_i32Chip[0]))
				pGame.sendBetResult(pUser, int64(pGame.ary_i32Chip[0]), JD21_BetResult_Auto)

				if pUser.iUserType != gamelogic.UserTypeRebot {
					pUser.iSumBetCount++
				}
			}
		}
	}

	if !bUserRecover {
		// 設定機器人 timeout
		for i := 0; i < pGame.iSeatCount; i++ {
			pUser := pGame.ary_pUserData[i]
			if pUser.iUserType == gamelogic.UserTypeRebot {
				pUser.SetRobotTimeout(pGame.map_i64Timeout[JD21_Timeout_Betting])
				if pUser.i64Timeout > 0 {
					pUser.i64LastMessageTime = pGame.i64TimeDuration
					if pUser.pRobot.iDefaultBetCount > 0 {
						pUser.pRobot.iBetCount = util.RndInt(0, pUser.pRobot.iDefaultBetCount+1)
					}
					pUser.iBetTimes = 0
				}
			}
		}
	}
}

// 訊息處理 - 玩家下注
func (pGame *JdesydLogic) message_Bet(pUser *JdesydUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	msgInfo := &netproto.JDESYD_Msg_Bet{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	i64Bet := msgInfo.GetBet()

	if pGame.iGameState != JD21_State_Betting {
		pGame.sendBetResult(pUser, i64Bet, JD21_BetResult_BettingStopped)
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_Betting {
		pGame.sendBetResult(pUser, i64Bet, JD21_BetResult_BettingStopped)
		return
	}

	iResult := pGame.userAction_Bet(pUser, i64Bet)
	pGame.sendBetResult(pUser, i64Bet, iResult)
}

// 玩家下注
func (pGame *JdesydLogic) userAction_Bet(pUser *JdesydUser, i64Bet int64) int {
	if pUser == nil {
		return JD21_BetResult_BetFailed
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return JD21_BetResult_BetFailed
	}

	iResult := JD21_BetResult_Success

	if i64Bet <= 0 {
		return JD21_BetResult_BetSmallThanZero
	}

	// 剩餘金錢是否足以下注
	if pUser.i64Money < i64Bet {
		return JD21_BetResult_NotEnoughMoney
	}

	// 總押注額是否超過上限
	if pUser.TotalBetBase()+i64Bet > g_i64MaxTotalBetBase {
		return JD21_BetResult_OverUpperLimit
	}

	pUser.Bet(pUser.CardSetUsedNow().iIndex-JD21_IndexBase, i64Bet)

	return iResult
}

// 傳送玩家下注結果
func (pGame *JdesydLogic) sendBetResult(pUser *JdesydUser, i64Bet int64, iResult int) {
	if pUser == nil {
		mlog.Debug("sendBetResult(), 無玩家資料")
		return
	}

	mlog.Debug("sendBetResult(), User[%d], SeatNo[%d], i64Bet[%d], iResult[%d]", pUser.User.UserID, pUser.User.SeatNo, i64Bet, iResult)

	iBetResult := iResult
	if iResult >= JD21_BetResult_Success {
		iBetResult = JD21_BetResult_Success
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_BetResult)

	pData := new(netproto.JDESYD_Msg_BetResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JD21_IndexBase)) // 下注玩家座位
	pData.Bet = proto.Int64(i64Bet)                                     // 下注金額
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                     // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(iBetResult))                       // 下注結果

	pMsg.MsgData = pData
	if iBetResult == JD21_BetResult_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}

	if iResult == JD21_BetResult_Auto {
		mlog.Debug("玩家[%d] 自动下注 %s, 总押注额 %s, 剩余金钱 %s", pUser.i32UserID,
			Global.MoneyString(i64Bet), Global.MoneyString(pUser.CardSetUsedNow().i64BetBase),
			Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 自动下注 %s, 总押注额 %s, 剩余金钱 %s", pUser.i32UserID,
			Global.MoneyString(i64Bet), Global.MoneyString(pUser.CardSetUsedNow().i64BetBase),
			Global.MoneyString(pUser.i64Money))
	} else {
		mlog.Debug("玩家[%d] 手牌 %d 下注 %s, 结果：%s, 总押注额 %s, 剩余金钱 %s", pUser.i32UserID,
			pUser.CardSetUsedNow().iIndex, Global.MoneyString(i64Bet), g_map_strBetResult[iResult],
			Global.MoneyString(pUser.CardSetUsedNow().i64BetBase), Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 手牌 %d 下注 %s, 结果：%s, 总押注额 %s, 剩余金钱 %s", pUser.i32UserID,
			pUser.CardSetUsedNow().iIndex, Global.MoneyString(i64Bet), g_map_strBetResult[iResult],
			Global.MoneyString(pUser.CardSetUsedNow().i64BetBase), Global.MoneyString(pUser.i64Money))
	}
}

// 播放停止下注動畫
func (pGame *JdesydLogic) stopBetting(bUserRecover bool, pUser *gamelogic.GameUser) {
	mlog.Debug("停止下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_StopBetting)

	if bUserRecover {
		if pUser != nil {
			pGame.SendMessageToUser(pUser, pMsg)
		}
	} else {
		pGame.BoardcastMessage(pMsg)

		pGame.uiPlayerReady = 0

		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
				pGame.ary_pUserData[iSeat].bOffline {
				pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_StopBettingOK
				pGame.uiPlayerReady |= 1 << uint(iSeat)
			}
		}
	}
}

// 訊息處理 - 停止下注動畫播放完成
func (pGame *JdesydLogic) message_StopBettingAnimationOK(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_StopBetting {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_Betting {
		return
	}

	mlog.Debug("User[%d] 停止下注動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JD21_UserPlayState_StopBettingOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家停止下注動畫播放完成")
		if pGame.bIsLastGame {
			// 上一局為最後一局，切牌
			pGame.setGameState(JD21_State_CutCards, pGame.map_i64Timeout[JD21_Timeout_CutCards])
			pGame.cutCards()
		} else {
			// 發牌
			pGame.setGameState(JD21_State_Deal, pGame.map_i64Timeout[JD21_Timeout_DealAnimation])
			pGame.deal()
		}
	}
}

// 現在輪到哪一個玩家
func (pGame *JdesydLogic) sendWhosTurn() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_WhosTurn)

	pData := new(netproto.JDESYD_Msg_WhosTurn)
	pData.SeatNo = proto.Int32(pGame.iWhosTurn + int32(JD21_IndexBase)) // 下注玩家座位

	pMsg.MsgData = pData

	pGame.BoardcastMessage(pMsg)

	if pGame.iWhosTurn >= 0 && pGame.iWhosTurn < int32(pGame.iSeatCount) {
		mlog.Debug("现在轮到玩家[%d] 操作", pGame.ary_pUserData[pGame.iWhosTurn].i32UserID)
		pGame.GameLog.AddGameStep("现在轮到玩家[%d] 操作", pGame.ary_pUserData[pGame.iWhosTurn].i32UserID)
	}
}

// 進入切牌階段
func (pGame *JdesydLogic) cutCards() {
	mlog.Debug("切牌")

	// 建立新的牌組
	pGame.pGamePoker.NewSet(JD21_PokerSetCount, 0)
	pGame.pGamePoker.Shuffle()

	pGame.bIsLastGame = false
	pGame.iLastGameCardIndex = JD21_MinCutCardIndex

	// 決定要切牌的玩家
	pGame.pickCutCardsuser()

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_UserCutCards)

	pData := new(netproto.JDESYD_Msg_UserCutCards)
	pData.SeatNo = proto.Int32(pGame.iCutCardsUser + int32(JD21_IndexBase))                                           // 執行切牌的玩家
	pData.MinCardIndex = proto.Int32(int32(JD21_MinCutCardIndex))                                                     // 切牌區間的最小值
	pData.MaxCardIndex = proto.Int32(int32(JD21_MaxCutCardIndex))                                                     // 切牌區間的最大值
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData

	pGame.BoardcastMessage(pMsg)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_CutCards
	}

	if pGame.bRobotCutCards {
		pGame.setGameState(JD21_State_CutCards, pGame.map_i64Timeout[JD21_Timeout_RobotCutCard])
	}
}

//  挑選切牌的玩家
func (pGame *JdesydLogic) pickCutCardsuser() {
	i32SeatNo := int32(-1)
	i64Money := int64(0)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal && pUser.iUserType != gamelogic.UserTypeRebot {
			if pUser.i64Money >= i64Money {
				i64Money = pUser.i64Money
				i32SeatNo = int32(iSeat)
			}
		}
	}

	if i32SeatNo >= 0 && i32SeatNo < int32(pGame.iSeatCount) {
		pGame.iCutCardsUser = i32SeatNo
		mlog.Debug("由玩家[%d(%d)] 切牌", pGame.ary_pUserData[pGame.iCutCardsUser].i32UserID, pGame.iCutCardsUser)
	} else {
		// 沒有可切牌的玩家，從機器人中挑選
		i64Money = 0
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.ary_pUserData[iSeat]
			if pUser.iUserState == JD21_UserState_Normal && pUser.iUserType == gamelogic.UserTypeRebot {
				if pUser.i64Money >= i64Money {
					i64Money = pUser.i64Money
					i32SeatNo = int32(iSeat)
				}
			}
		}
		pGame.iCutCardsUser = i32SeatNo
		pGame.bRobotCutCards = true

		mlog.Debug("沒有可切牌的玩家, 從機器人中挑選, SeatNo[%d]", pGame.iCutCardsUser)
		if pGame.iCutCardsUser >= 0 && pGame.iCutCardsUser < int32(pGame.iSeatCount) {
			mlog.Debug("由機器人[%d(%d)] 切牌", pGame.ary_pUserData[pGame.iCutCardsUser].i32UserID, pGame.iCutCardsUser)
		}

		// 隨機切牌
		pGame.iLastGameCardIndex = util.RndInt(JD21_MinCutCardIndex, JD21_MaxCutCardIndex)
		mlog.Debug("iLastGameCardIndex[%d]", pGame.iLastGameCardIndex)
	}
}

// 玩家切牌
func (pGame *JdesydLogic) message_CutCards(pUser *JdesydUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	mlog.Debug("TableNo[%d], message_CutCards(), iGameState[%s], User[%d], iPlayingState[%s]", pGame.TableNo(),
		g_map_strGameState_debug[pGame.iGameState], pUser.i32UserID, g_map_strPlayingState_debug[pUser.iPlayingState])

	if pUser.i32SeatNo != pGame.iCutCardsUser {
		return
	}

	if pGame.iGameState != JD21_State_CutCards {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_CutCards {
		return
	}

	msgInfo := &netproto.JDESYD_Msg_CutCards{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	iCardIndex := int(msgInfo.GetCardIndex())

	if iCardIndex >= JD21_MinCutCardIndex && iCardIndex <= JD21_MaxCutCardIndex {
		pGame.iLastGameCardIndex = iCardIndex
	}

	mlog.Debug("玩家[%d] 切牌，位置：%d", pGame.ary_pUserData[pGame.iCutCardsUser].i32UserID, pGame.iLastGameCardIndex)
	pGame.GameLog.AddGameStep("玩家[%d] 切牌，位置：%d", pGame.ary_pUserData[pGame.iCutCardsUser].i32UserID, pGame.iLastGameCardIndex)

	// 切牌
	pNewSet := new(Pokerset.PokerCardSet)
	pNewSet.SetCardSet(pGame.pGamePoker.CodeString())
	pSet := pNewSet.Deal(pGame.iLastGameCardIndex)
	pGame.pGamePoker.SetCardSet(pNewSet.CodeString() + pSet.CodeString())

	pGame.setGameState(JD21_State_CutCardsAnimation, pGame.map_i64Timeout[JD21_Timeout_CutCardsAnimation])
	pGame.sendCutCardsResult(false, nil)
}

// 玩家切牌結果
func (pGame *JdesydLogic) sendCutCardsResult(bUserRecover bool, pGameUser *gamelogic.GameUser) {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_CutCardsResult)

	pData := new(netproto.JDESYD_Msg_CutCardsResult)
	pData.SeatNo = proto.Int32(pGame.iCutCardsUser + int32(JD21_IndexBase)) // 執行切牌的玩家
	pData.CardIndex = proto.Int32(int32(pGame.iLastGameCardIndex))          // 切牌位置索引值

	pMsg.MsgData = pData

	if !bUserRecover {
		pGame.BoardcastMessage(pMsg)
	} else {
		if pGameUser != nil {
			pGame.SendMessageToUser(pGameUser, pMsg)
		}
	}

	pGame.iCutCardsUser = -2

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.ary_pUserData[iSeat].bOffline {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_CutCardsOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}
	}
}

// 切牌動畫播放完成
func (pGame *JdesydLogic) message_CutCardsAnimationOK(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_CutCardsAnimation {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_CutCards {
		return
	}

	mlog.Debug("User[%d] 切牌動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JD21_UserPlayState_CutCardsOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家切牌動畫播放完成")
		pGame.setGameState(JD21_State_Deal, pGame.map_i64Timeout[JD21_Timeout_DealAnimation])
		pGame.deal()
	}

}

// 發牌
func (pGame *JdesydLogic) deal() {
	mlog.Debug("發牌, iDealType[%d]", pGame.iDealType)

	switch pGame.iDealType {
	case JD21_DealType_Assigned: // 指定發牌
		// 發玩家牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_Normal {
				pCardSet := pGame.dealing(int32(iSeat), 2)
				pNewCardSet := NewCardSetInfo(len(pGame.ary_pUserData[iSeat].ary_pCardSet) + JD21_IndexBase)
				pNewCardSet.AddCardSet(pCardSet)
				if !pGame.ary_pUserData[iSeat].AddCardSet(pNewCardSet) {
					mlog.Error("玩家新增牌組失敗，牌組數量 %d", len(pGame.ary_pUserData[iSeat].ary_pCardSet))
				}

				if pGame.ary_pUserData[iSeat].CardSetUsedNow().bIsBlackJack {
					pGame.ary_pUserData[iSeat].CardSetUsedNow().bOerationDone = true
				}
			}
		}

		// 發莊家牌
		pCardSet := pGame.dealing(-1, 2)
		pGame.pDealerCard.AddCardSet(pCardSet)
	default:
		// 發玩家牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_Normal {
				pCardSet := pGame.pGamePoker.Deal(2)
				pNewCardSet := NewCardSetInfo(len(pGame.ary_pUserData[iSeat].ary_pCardSet) + JD21_IndexBase)
				pNewCardSet.AddCardSet(pCardSet)
				if !pGame.ary_pUserData[iSeat].AddCardSet(pNewCardSet) {
					mlog.Error("玩家新增牌組失敗，牌組數量 %d", len(pGame.ary_pUserData[iSeat].ary_pCardSet))
				}

				if pGame.ary_pUserData[iSeat].CardSetUsedNow().bIsBlackJack {
					pGame.ary_pUserData[iSeat].CardSetUsedNow().bOerationDone = true
				}
			}
		}

		// 發莊家牌
		pCardSet := pGame.pGamePoker.Deal(2)
		pGame.pDealerCard.AddCardSet(pCardSet)
	}

	// 是否要買保險
	pOpenCard := pGame.pDealerCard.pCardSet.GetSingleCard(0, false)
	if pOpenCard != nil {
		if pOpenCard.I32Number == Pokerset.POKERCardNo1 || pOpenCard.I32Number == Pokerset.POKERCardNoA {
			pGame.bTakeInsuranceOrNot = true
		}
	}

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_Deal)

	pData := new(netproto.JDESYD_Msg_Deal)
	pData.DealerCards = pGame.cardSetInfoPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}
	pData.IsLastGame = proto.Bool(pGame.bIsLastGame) // 是否為最後一局

	pMsg.MsgData = pData

	pGame.BoardcastMessage(pMsg)

	mlog.Debug("发牌")
	pGame.GameLog.AddGameStep("发牌")
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal {
			for i := 0; i < len(pUser.ary_pCardSet); i++ {
				if pUser.ary_pCardSet[i].pCardSet.CardCount() > 0 {
					mlog.Debug("玩家[%d] 手牌 %d：%s (%d, %d)", pUser.i32UserID, i+JD21_IndexBase, pUser.ary_pCardSet[i].pCardSet.String(),
						pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2)
					pGame.GameLog.AddGameStep("玩家[%d] 手牌 %d：%s (%d, %d)", pUser.i32UserID, i+JD21_IndexBase, pUser.ary_pCardSet[i].pCardSet.String(),
						pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2)
				}
			}
		}
	}
	mlog.Debug("庄家牌：%s (%d, %d), 明牌：%s", pGame.pDealerCard.pCardSet.String(), pGame.pDealerCard.i32Point_1,
		pGame.pDealerCard.i32Point_2, pOpenCard.String())
	pGame.GameLog.AddGameStep("庄家牌：%s (%d, %d), 明牌：%s", pGame.pDealerCard.pCardSet.String(), pGame.pDealerCard.i32Point_1,
		pGame.pDealerCard.i32Point_2, pOpenCard.String())

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal {
			for i := 0; i < len(pUser.ary_pCardSet); i++ {
				if pUser.ary_pCardSet[i].pCardSet.CardCount() > 0 {
				}
			}
		}
	}

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.ary_pUserData[iSeat].bOffline {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_DealOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}
	}

	// 檢查是否為最後一局
	if !pGame.bIsLastGame {
		if pGame.IsLastGame() {
			pGame.sendLastGame()
		}
	}
}

// 檢查是否為最後一局
func (pGame *JdesydLogic) IsLastGame() bool {
	if pGame.pGamePoker.CardCount() <= g_iCutCardIndex || pGame.bShuffleEveryGame {
		pGame.bIsLastGame = true
		return true
	}

	return false
}

// 最後一局
func (pGame *JdesydLogic) sendLastGame() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_LastGame)

	pGame.BoardcastMessage(pMsg)

	mlog.Debug("最后一局，切牌位置 %d，牌盒剩余张数 %d", pGame.iLastGameCardIndex, pGame.pGamePoker.CardCount())
	pGame.GameLog.AddGameStep("最后一局，切牌位置 %d，牌盒剩余张数 %d", pGame.iLastGameCardIndex, pGame.pGamePoker.CardCount())
}

// 發牌動畫播放完成
func (pGame *JdesydLogic) message_DealAnimationOK(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_Deal {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_CutCardsOK {
		return
	}

	mlog.Debug("User[%d] 發牌動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = JD21_UserPlayState_DealOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家發牌動畫播放完成")
		// 是否要買保險
		if pGame.bTakeInsuranceOrNot {
			pGame.setGameState(JD21_State_Insurance, pGame.map_i64Timeout[JD21_Timeout_Insurance])
			pGame.sendInsuranceOrNot(false, nil)
		} else {
			// 進入玩家操作階段
			pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
			pGame.startPlaying()
		}
	}
}

// 詢問玩家是否買保險
func (pGame *JdesydLogic) sendInsuranceOrNot(bUserRecover bool, pGameUser *gamelogic.GameUser) {
	mlog.Debug("买保险")

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_InsuranceOrNot)

	pData := new(netproto.JDESYD_Msg_InsuranceOrNot)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData

	if bUserRecover {
		if pGameUser != nil {
			pGame.SendMessageToUser(pGameUser, pMsg)
		}
	} else {
		pGame.BoardcastMessage(pMsg)

		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.ary_pUserData[iSeat]
			if pUser.iUserState == JD21_UserState_Normal {
				pUser.iPlayingState = JD21_UserPlayState_TakeInsurance

				// 設定機器人觸發時間
				if pUser.iUserType == gamelogic.UserTypeRebot {
					pUser.SetRobotTimeout(pGame.map_i64Timeout[JD21_Timeout_Insurance])
					if pUser.i64Timeout > 0 {
						pUser.i64LastMessageTime = pGame.i64TimeDuration
					}
				}
			}
		}

		pGame.uiPlayerReady = 0

		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.ary_pUserData[iSeat]
			if pUser.iUserState == JD21_UserState_EmptySeat || pUser.bOffline || pUser.CardSetUsedNow().bIsBlackJack {
				pUser.iPlayingState = JD21_UserPlayState_TakeInsuranceOK
				pGame.uiPlayerReady |= 1 << uint(iSeat)
			}
		}
	}
}

// 玩家買保險
func (pGame *JdesydLogic) message_TakeInsurance(pUser *JdesydUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	if pGame.iGameState != JD21_State_Insurance {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_TakeInsurance {
		return
	}

	if pUser.CardSetUsedNow().bIsBlackJack {
		return
	}

	msgInfo := &netproto.JDESYD_Msg_TakeInsurance{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	iAction := int(msgInfo.GetAction())

	// 保險金額
	pGame.userAction_TakeInsurance(pUser, iAction)

	pUser.iPlayingState = JD21_UserPlayState_TakeInsuranceOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家都已決定是否買保險")
		if pGame.bHasUserTakeInsurance {
			// 有玩家買保險，莊家瞇牌
			pGame.setGameState(JD21_State_DealerPeeksCards, pGame.map_i64Timeout[JD21_Timeout_DealerPeeksCards])
			pGame.dealerPeeksCards()
		} else {
			// 進入玩家操作階段
			pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
			pGame.startPlaying()
		}
	}
}

// 玩家行為 - 買保險
func (pGame *JdesydLogic) userAction_TakeInsurance(pUser *JdesydUser, iAction int) {
	if pUser == nil {
		return
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return
	}

	if iAction == 0 {
		mlog.Debug("玩家[%d] 不买保险", pUser.i32UserID)
		pGame.GameLog.AddGameStep("玩家[%d] 不买保险", pUser.i32UserID)
	} else {
		if pUser.i64Money < pUser.TotalBetBase()/2 {
			pGame.sendTakeInsuranceResult(pUser, JD21_Insurance_NotEnoughMoney)
		} else {
			pUser.i64Insurance = pUser.TotalBetBase() / 2
			pUser.i64Money -= pUser.i64Insurance
			pGame.bHasUserTakeInsurance = true
			pGame.sendTakeInsuranceResult(pUser, JD21_Insurance_Success)

			if pUser.iUserType != gamelogic.UserTypeRebot {
				pUser.iSumBetCount++
				pUser.i64SumBetMoney += pUser.i64Insurance
			}

			mlog.Debug("玩家[%d] 买保险 %s, 剩余金钱 %s", pUser.i32UserID,
				Global.MoneyString(pUser.i64Insurance), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 买保险 %s, 剩余金钱 %s", pUser.i32UserID,
				Global.MoneyString(pUser.i64Insurance), Global.MoneyString(pUser.i64Money))
		}
	}
}

// 玩家買保險結果
func (pGame *JdesydLogic) sendTakeInsuranceResult(pUser *JdesydUser, iResult int) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_TakeInsuranceResult)

	pData := new(netproto.JDESYD_Msg_TakeInsuranceResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JD21_IndexBase)) // 買保險的玩家
	pData.Insurance = proto.Int64(pUser.i64Insurance)                   // 保險金額
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                     // 玩家剩餘金錢
	pData.Result = proto.Int32(int32(iResult))                          // 買保險結果

	pMsg.MsgData = pData

	if iResult == JD21_Insurance_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}

	mlog.Debug("玩家[%d] 買保險結果：金額[%d], 結果[%d]", pUser.i32UserID, pUser.i64Insurance, iResult)
}

// 莊家瞇牌
func (pGame *JdesydLogic) dealerPeeksCards() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_DealerPeeksCards)

	pGame.BoardcastMessage(pMsg)

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.ary_pUserData[iSeat].bOffline {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_DealerOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}
	}
}

// 莊家瞇牌動畫播放完成
func (pGame *JdesydLogic) message_DealerPeeksCardsDone(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_DealerPeeksCards {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_TakeInsuranceOK {
		return
	}

	pUser.iPlayingState = JD21_UserPlayState_DealerOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家莊家瞇牌動畫都已播放完成")
		pGame.setGameState(JD21_State_Settlement_Insurance, pGame.map_i64Timeout[JD21_Timeout_Settlement_Insurance])
		pGame.settlement_Insurance()
	}
}

// 保險結算
func (pGame *JdesydLogic) settlement_Insurance() {
	mlog.Debug("保险结算")
	pGame.GameLog.AddGameStep("保险结算")

	// 保險結算
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_Normal {
			pUser := pGame.ary_pUserData[iSeat]
			i64MoneyChange := int64(0)
			if !pGame.pDealerCard.bIsBlackJack {
				// 莊家牌不是 BlackJack
				i64MoneyChange = -pUser.i64Insurance
				pUser.i64TotalWinLose = i64MoneyChange

				mlog.Debug("玩家[%d] 赔付保险额 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Insurance),
					Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 赔付保险额 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Insurance),
					Global.MoneyString(pUser.i64Money))
			} else if pUser.ary_pCardSet[0].bIsBlackJack {
				i64MoneyChange = pUser.TotalBetBase()
				pUser.i64TotalWinLose = 0
				pUser.i64Money += i64MoneyChange

				pGame.pDealerCard.bShowCards = true

				mlog.Debug("玩家[%d] 也是黑杰克, 退回底注 %s, 剩余金钱 %s",
					pUser.i32UserID,
					Global.MoneyString(pUser.TotalBetBase()),
					Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 也是黑杰克, 退回底注 %s, 剩余金钱 %s",
					pUser.i32UserID,
					Global.MoneyString(pUser.TotalBetBase()),
					Global.MoneyString(pUser.i64Money))
			} else {
				i64MoneyChange = pUser.i64Insurance*int64(JD21_InsuranceMultiple) - pUser.TotalBetBase()
				pUser.i64TotalWinLose = i64MoneyChange
				pUser.i64Money += pUser.i64Insurance + pUser.i64Insurance*int64(JD21_InsuranceMultiple)

				pGame.pDealerCard.bShowCards = true

				mlog.Debug("玩家[%d] 获得保险金 %s, 赔付底注 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Insurance*int64(JD21_InsuranceMultiple)),
					Global.MoneyString(pUser.TotalBetBase()), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 获得保险金 %s, 赔付底注 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Insurance*int64(JD21_InsuranceMultiple)),
					Global.MoneyString(pUser.TotalBetBase()), Global.MoneyString(pUser.i64Money))
			}
			pUser.i64MoneyChange += i64MoneyChange

			// 寫分
			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)
		}
	}

	//----- 傳送訊息給 client -----
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_Settlement_Insurance)

	pData := new(netproto.JDESYD_Msg_Settlement_Insurance)
	pData.DealerCards = pGame.cardSetInfoPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData

	pGame.BoardcastMessage(pMsg)

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.ary_pUserData[iSeat].bOffline {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_Settlement_InsuranceOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}

		pGame.ary_pUserData[iSeat].i64TotalWinLose = 0
	}
}

// 保險結算動畫播放完成
func (pGame *JdesydLogic) message_Settlement_Insurance_OK(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_Settlement_Insurance {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_DealerOK {
		return
	}

	pUser.iPlayingState = JD21_UserPlayState_Settlement_InsuranceOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家保險結算動畫都已播放完成")
		if pGame.pDealerCard.bIsBlackJack {
			// 莊家牌是 BlackJack，結束本局
			pGame.setGameState(JD21_State_GameOver, 0)
			pGame.gameOver()
		} else {
			// 進入玩家操作階段
			pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
			pGame.startPlaying()
		}
	}
}

// 進入玩家操作階段
func (pGame *JdesydLogic) startPlaying() {
	mlog.Debug("進入玩家操作階段")

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_Normal {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_Playing_Waiting
		} else {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_Playing_Done
		}
	}

	pGame.checkWhosTurn()
}

// 檢查現在輪到誰
func (pGame *JdesydLogic) checkWhosTurn() {
	i32WhosTurn := pGame.iWhosTurn
	if i32WhosTurn < 0 {
		i32WhosTurn = 0
	}

	bAllDone := true
	for iSeat := i32WhosTurn; iSeat < int32(pGame.iSeatCount); iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal {
			if pUser.iPlayingState != JD21_UserPlayState_Playing_Done {
				if !pUser.AllCardSetDone() {
					pUser.iPlayingState = JD21_UserPlayState_Playing_MyTurn
					i32WhosTurn = iSeat

					if i32WhosTurn != pGame.iWhosTurn {
						mlog.Debug("轮到玩家[%d] 操作", pUser.i32UserID)
						pGame.GameLog.AddGameStep("轮到玩家[%d] 操作", pUser.i32UserID)
					}
					pGame.iWhosTurn = i32WhosTurn
					bAllDone = false
					break
				}
			}
		}
	}

	if bAllDone {
		// 莊家開牌/結算
		pGame.setGameState(JD21_State_Settlement_Dealer, pGame.map_i64Timeout[JD21_Timeout_Settlement_Dealer])
		pGame.settlement_Dealer()
	} else {
		// 通知玩家操作
		pGame.setGameState(JD21_State_Playing, pGame.map_i64Timeout[JD21_Timeout_Playing])
		pGame.sendStartPlaying(false, nil)
	}
}

// 通知玩家操作
func (pGame *JdesydLogic) sendStartPlaying(bUserRecover bool, pGameUser *gamelogic.GameUser) {
	i32WhosTurn := pGame.iWhosTurn
	if i32WhosTurn < 0 || i32WhosTurn >= int32(pGame.iSeatCount) {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_StartPlaying)
	pData := new(netproto.JDESYD_Msg_StartPlaying)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間
	pData.SeatNo = proto.Int32(i32WhosTurn + int32(JD21_IndexBase))                                                   // 動作玩家的座位
	pData.CardSetIndex = proto.Int32(int32(pGame.ary_pUserData[i32WhosTurn].CardSetUsedNow().iIndex))                 // 手上第幾組牌，從 1 開始
	pData.CanSplit = proto.Bool(pGame.ary_pUserData[i32WhosTurn].CardSetUsedNow().bCanSplit)                          // 是否可以分牌
	pData.CanDouble = proto.Bool(pGame.ary_pUserData[i32WhosTurn].CardSetUsedNow().bCanDouble)                        // 是否可以加倍

	pMsg.MsgData = pData
	if !bUserRecover {
		pGame.BoardcastMessage(pMsg)

		// 設定機器人 timeout
		if pGame.ary_pUserData[i32WhosTurn].iUserType == gamelogic.UserTypeRebot {
			pGame.ary_pUserData[i32WhosTurn].SetRobotTimeout(pGame.map_i64Timeout[JD21_Timeout_Playing])
			if pGame.ary_pUserData[i32WhosTurn].i64Timeout > 0 {
				pGame.ary_pUserData[i32WhosTurn].i64LastMessageTime = pGame.i64TimeDuration
			}
		}
	} else {
		if pGameUser != nil {
			pGame.SendMessageToUser(pGameUser, pMsg)
		}
	}
}

// 玩家動作
func (pGame *JdesydLogic) message_UserPlaying(pUser *JdesydUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_Playing_MyTurn {
		return
	}

	msgInfo := &netproto.JDESYD_Msg_UserPlaying{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	iAction := int(msgInfo.GetAction())
	iCardSetIdx := int(msgInfo.GetCardSetIndex())

	if iCardSetIdx < JD21_IndexBase || iCardSetIdx > JD21_MaxSplitCount {
		return
	}

	if iCardSetIdx != pUser.CardSetUsedNow().iIndex {
		return
	}

	if pUser.CardSetUsedNow().bOerationDone {
		return
	}

	iResult := JD21_Playing_Success

	switch iAction {
	case JD21_UserAction_Stand: // 停牌
		pUser.Stand()

		mlog.Debug("玩家[%d] 停牌", pUser.i32UserID)
		pGame.GameLog.AddGameStep("玩家[%d] 停牌", pUser.i32UserID)
	case JD21_UserAction_Split: // 分牌
		iResult = pGame.userAction_Split(pUser)
	case JD21_UserAction_Double: // 加倍
		iResult = pGame.userAction_Double(pUser)
	case JD21_UserAction_Hit: // 要牌
		iResult = pGame.userAction_Hit(pUser)
	case JD21_UserAction_Surrender: // 投降
		iResult = pGame.userAction_Surrender(pUser)
	default:
		iResult = JD21_Playing_Failed
	}

	pGame.sendUserPlayingResult(pUser, iAction, iResult)

	if pUser.CardSetUsedNow().bIsBusted {
		pUser.CardSetUsedNow().bOerationDone = true
	}

	// 檢查是否為最後一局
	if !pGame.bIsLastGame {
		if pGame.IsLastGame() {
			pGame.sendLastGame()
		}
	}

	if iResult == JD21_Playing_Success {
		// 延遲 0.5 秒進行下一個動作
		pGame.setGameState(JD21_State_PlayingDelay, pGame.map_i64Timeout[JD21_Timeout_PlayingDelay])
	}
}

// 玩家操作 - 分牌
func (pGame *JdesydLogic) userAction_Split(pUser *JdesydUser) int {
	if pUser == nil {
		return JD21_Playing_Failed
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return JD21_Playing_Failed
	}

	iResult := JD21_Playing_Success

	if len(pUser.ary_pCardSet) >= JD21_MaxSplitCount {
		iResult = JD21_Playing_CannotSplit
	} else if pUser.CardSetUsedNow().bOerationDone {
		iResult = JD21_Playing_CardSetOperationDone
	} else if !pUser.CardSetUsedNow().bCanSplit {
		iResult = JD21_Playing_CannotSplit
	} else {
		if pUser.i64Money < pUser.CardSetUsedNow().i64BetBase {
			iResult = JD21_Playing_NotEnoughMoney

			mlog.Debug("玩家[%d] 分牌，余额不足，手牌 %d：%s", pUser.i32UserID, pUser.CardSetUsedNow().iIndex,
				pUser.CardSetUsedNow().pCardSet.String())
			pGame.GameLog.AddGameStep("玩家[%d] 分牌，余额不足，手牌 %d：%s", pUser.i32UserID, pUser.CardSetUsedNow().iIndex,
				pUser.CardSetUsedNow().pCardSet.String())
		} else {
			mlog.Debug("玩家[%d] 分牌，手牌 %d：%s", pUser.i32UserID, pUser.CardSetUsedNow().iIndex, pUser.CardSetUsedNow().pCardSet.String())
			pGame.GameLog.AddGameStep("玩家[%d] 分牌，手牌 %d：%s", pUser.i32UserID, pUser.CardSetUsedNow().iIndex, pUser.CardSetUsedNow().pCardSet.String())

			pCardSet := pGame.dealing(pUser.i32SeatNo, 2)
			pUser.Split(pCardSet)

			for i := 0; i < len(pUser.ary_pCardSet); i++ {
				if pUser.ary_pCardSet[i].pCardSet.CardCount() > 0 {
					mlog.Debug("玩家[%d] 手牌 %d：%s (%d, %d), 押注额 %s", pUser.i32UserID, i+JD21_IndexBase,
						pUser.ary_pCardSet[i].pCardSet.String(), pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2,
						Global.MoneyString(pUser.ary_pCardSet[i].i64BetBase))
					pGame.GameLog.AddGameStep("玩家[%d] 手牌 %d：%s (%d, %d), 押注额 %s", pUser.i32UserID, i+JD21_IndexBase,
						pUser.ary_pCardSet[i].pCardSet.String(), pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2,
						Global.MoneyString(pUser.ary_pCardSet[i].i64BetBase))
				}
			}
			mlog.Debug("玩家[%d] 剩余金额 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 剩余金额 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
		}
	}

	return iResult
}

// 玩家操作 - 加倍
func (pGame *JdesydLogic) userAction_Double(pUser *JdesydUser) int {
	if pUser == nil {
		return JD21_Playing_Failed
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return JD21_Playing_Failed
	}

	iResult := JD21_Playing_Success

	if pUser.CardSetUsedNow().bOerationDone {
		iResult = JD21_Playing_CardSetOperationDone
	} else if !pUser.CardSetUsedNow().bCanDouble {
		iResult = JD21_Playing_CannotDouble
	} else {
		if pUser.i64Money < pUser.CardSetUsedNow().i64BetBase {
			iResult = JD21_Playing_NotEnoughMoney

			mlog.Debug("玩家[%d] 加倍，余额不足", pUser.i32UserID)
			pGame.GameLog.AddGameStep("玩家[%d] 加倍，余额不足", pUser.i32UserID)
		} else {
			pCardSet := pGame.dealing(pUser.i32SeatNo, 1)
			pUser.Double(pCardSet)

			mlog.Debug("玩家[%d] 加倍 %s, 底注押注额 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.CardSetUsedNow().i64BetDouble),
				Global.MoneyString(pUser.TotalBetBase()), Global.MoneyString(pUser.i64Money))
			mlog.Debug("玩家[%d] 手牌 %d：%s (%d, %d)", pUser.i32UserID, pUser.iCardSetIndexNow+JD21_IndexBase, pUser.CardSetUsedNow().pCardSet.String(),
				pUser.CardSetUsedNow().i32Point_1, pUser.CardSetUsedNow().i32Point_2)
			pGame.GameLog.AddGameStep("玩家[%d] 加倍 %s, 底注押注额 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.CardSetUsedNow().i64BetDouble),
				Global.MoneyString(pUser.TotalBetBase()), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 手牌 %d：%s (%d, %d)", pUser.i32UserID, pUser.iCardSetIndexNow+JD21_IndexBase, pUser.CardSetUsedNow().pCardSet.String(),
				pUser.CardSetUsedNow().i32Point_1, pUser.CardSetUsedNow().i32Point_2)

			// 爆牌
			if pUser.CardSetUsedNow().bIsBusted {
				pUser.CardSetUsedNow().i64WinLose = -(pUser.CardSetUsedNow().i64BetBase + pUser.CardSetUsedNow().i64BetDouble)
				pUser.i64TotalWinLose += pUser.CardSetUsedNow().i64WinLose
				i64MoneyChange := pUser.CardSetUsedNow().i64WinLose

				pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)

				pUser.i64MoneyChange += i64MoneyChange
				pUser.CardSetUsedNow().bCompensated = true

				mlog.Debug("玩家[%d] 爆牌，输赢 %s，剩余金钱 %s", pUser.i32UserID, Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 爆牌，输赢 %s，剩余金钱 %s", pUser.i32UserID, Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
			}
		}
	}

	return iResult
}

// 玩家操作 - 要牌
func (pGame *JdesydLogic) userAction_Hit(pUser *JdesydUser) int {
	if pUser == nil {
		return JD21_Playing_Failed
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return JD21_Playing_Failed
	}

	iResult := JD21_Playing_Success

	if pUser.CardSetUsedNow().bOerationDone {
		iResult = JD21_Playing_CardSetOperationDone
	} else if pUser.CardSetUsedNow().bIsBusted {
		iResult = JD21_Playing_IsBusted
	} else {
		pCardSet := pGame.dealing(pUser.i32SeatNo, 1)
		strHit := pCardSet.String()
		pUser.Hit(pCardSet)

		mlog.Debug("玩家[%d] 要牌 %s, 手牌 %d：%s (%d, %d)", pUser.i32UserID, strHit, pUser.iCardSetIndexNow+JD21_IndexBase,
			pUser.CardSetUsedNow().pCardSet.String(), pUser.CardSetUsedNow().i32Point_1, pUser.CardSetUsedNow().i32Point_2)
		pGame.GameLog.AddGameStep("玩家[%d] 要牌 %s, 手牌 %d：%s (%d, %d)", pUser.i32UserID, strHit, pUser.iCardSetIndexNow+JD21_IndexBase,
			pUser.CardSetUsedNow().pCardSet.String(), pUser.CardSetUsedNow().i32Point_1, pUser.CardSetUsedNow().i32Point_2)

		// 爆牌
		if pUser.CardSetUsedNow().bIsBusted {
			pUser.CardSetUsedNow().i64WinLose = -pUser.CardSetUsedNow().i64BetBase
			pUser.i64TotalWinLose += pUser.CardSetUsedNow().i64WinLose
			i64MoneyChange := pUser.CardSetUsedNow().i64WinLose

			pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)
			pUser.i64MoneyChange += i64MoneyChange
			pUser.CardSetUsedNow().bCompensated = true

			mlog.Debug("玩家[%d] 爆牌，输赢 %s，剩余金钱 %s", pUser.i32UserID, Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 爆牌，输赢 %s，剩余金钱 %s", pUser.i32UserID, Global.MoneyString(i64MoneyChange), Global.MoneyString(pUser.i64Money))
		}
	}

	return iResult
}

// 玩家操作 - 投降
func (pGame *JdesydLogic) userAction_Surrender(pUser *JdesydUser) int {
	if pUser == nil {
		return JD21_Playing_Failed
	}

	if pUser.iUserState != JD21_UserState_Normal {
		return JD21_Playing_Failed
	}

	iResult := JD21_Playing_Success

	if pUser.CardSetUsedNow().bOerationDone {
		iResult = JD21_Playing_CardSetOperationDone
	} else if pUser.bSurrendered {
		iResult = JD21_Playing_CannotSurrender
	} else {
		pUser.Surrender()
		pUser.CardSetUsedNow().i64WinLose = pUser.TotalBetBase() / 2
		pUser.i64TotalWinLose += pUser.CardSetUsedNow().i64WinLose
		pUser.i64Money += pUser.i64TotalWinLose
		i64MoneyChange := -pUser.CardSetUsedNow().i64WinLose

		pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, 0)

		pUser.i64MoneyChange += i64MoneyChange
		pUser.CardSetUsedNow().bCompensated = true

		mlog.Debug("玩家[%d] 投降, 输赢 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64TotalWinLose),
			Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 投降, 输赢 %s, 剩余金钱 %s", pUser.i32UserID, Global.MoneyString(pUser.i64TotalWinLose),
			Global.MoneyString(pUser.i64Money))
	}

	return iResult
}

// 傳送玩家動作結果
func (pGame *JdesydLogic) sendUserPlayingResult(pUser *JdesydUser, iAction, iResult int) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_UserPlayingResult)

	pData := new(netproto.JDESYD_Msg_UserPlayingResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(JD21_IndexBase))                                               // 執行動作的玩家
	pData.CardSetIndex = proto.Int32(int32(pUser.iCardSetIndexNow + JD21_IndexBase))                                  // 手上第幾組牌，從 1 開始
	pData.Action = proto.Int32(int32(iAction))                                                                        // 執行的動作，0:停牌，1:分牌，2:加倍，3:要牌，4:投降
	pData.Result = proto.Int32(int32(iResult))                                                                        // 結果
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間
	pData.UserInfo = pGame.userInfoPkg(int(pUser.i32SeatNo))                                                          // 執行動作玩家資訊

	pMsg.MsgData = pData

	if iResult == JD21_Playing_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}
}

// 莊家開牌/結算
func (pGame *JdesydLogic) settlement_Dealer() {
	mlog.Debug("莊家開牌/結算")

	// 莊家開牌/結算
	pGame.pDealerCard.bShowCards = true

	// 檢查是否所有牌組都已賠付
	bAllCompensated := true
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_Normal && !pGame.ary_pUserData[iSeat].AllCompensated() {
			bAllCompensated = false
			break
		}
	}

	// 檢查尚未賠付的牌組是否全是黑傑克
	bAllBlackJack := true
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		for i := 0; i < len(pGame.ary_pUserData[iSeat].ary_pCardSet); i++ {
			if !pGame.ary_pUserData[iSeat].ary_pCardSet[i].bCompensated {
				if !pGame.ary_pUserData[iSeat].ary_pCardSet[i].bIsBlackJack {
					bAllBlackJack = false
					break
				}
			}
		}
	}

	mlog.Debug("莊家牌：%s (%d, %d), MaxPoint[%d]", pGame.pDealerCard.pCardSet.String(), pGame.pDealerCard.i32Point_1,
		pGame.pDealerCard.i32Point_2, pGame.pDealerCard.MaxPoint())
	mlog.Debug("bIsBlackJack[%t], bIsBusted[%t], bAllCompensated[%t], CheckSoftReached[%t], bAllBlackJack[%t]",
		pGame.pDealerCard.bIsBlackJack, pGame.pDealerCard.bIsBusted, bAllCompensated, pGame.pDealerCard.IsSoftReached(), bAllBlackJack)

	// 莊家補牌
	bNaturalDealing := true
	i32DealerPoint := pGame.pDealerCard.MaxPoint()
	i32TargetPoint := int32(0)

	if pGame.iDealType == JD21_DealType_Natural {
		// 若發牌方式為自然發牌，則檢核莊家補牌表
		if hasDealerControlPoint(-1, i32DealerPoint) {
			mlog.Debug("hasDealerControlPoint, i32DealerPoint[%d]", i32DealerPoint)
			DealerSupplement_Now_Debug()
			// 莊家點數在莊家補牌表控制範圍內，決定補到多少點或是自然發牌
			iRnd := util.RndInt(0, g_map_pDealerSupplement_Now[i32DealerPoint].iOddsSum)
			for i := 0; i < len(g_map_pDealerSupplement_Now[i32DealerPoint].ary_pData); i++ {
				iRange1 := g_map_pDealerSupplement_Now[i32DealerPoint].ary_pData[i].iRndRange
				if i == 0 {
					if iRnd <= iRange1 {
						// 自然發牌
						break
					}
				} else {
					iRange2 := g_map_pDealerSupplement_Now[i32DealerPoint].ary_pData[i-1].iRndRange
					if iRnd > iRange2 && iRnd <= iRange1 {
						i32TargetPoint = g_map_pDealerSupplement_Now[i32DealerPoint].ary_pData[i].i32Points
						break
					}
				}
			}

			mlog.Debug("Rnd[%d], TargetPoint[%d]", iRnd, i32TargetPoint)

			if i32TargetPoint > 0 {
				bNaturalDealing = false
			}
		}
	}

	mlog.Debug("庄家开牌结算")
	pGame.GameLog.AddGameStep("庄家开牌结算")
	str := fmt.Sprintf("庄家点数 %d，", i32DealerPoint)
	if bNaturalDealing {
		str += "自然发牌"
	} else {
		str += fmt.Sprintf("补到 %d 点", i32TargetPoint)
	}
	mlog.Debug("%s", str)
	pGame.GameLog.AddGameStep("%s", str)

	if !bNaturalDealing {
		// 依莊家補牌表發牌
		pCardSet := getCardSetByAssignedPoint(pGame.pGamePoker, pGame.pDealerCard, i32TargetPoint)
		pGame.pDealerCard.AddCardSet(pCardSet)
		mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), pGame.pDealerCard.pCardSet.String(),
			pGame.pDealerCard.i32Point_1, pGame.pDealerCard.i32Point_2)
	} else {
		// 自然發牌
		if !pGame.pDealerCard.bIsBlackJack && !pGame.pDealerCard.bIsBusted && !bAllCompensated && !pGame.pDealerCard.IsSoftReached() && !bAllBlackJack {
			for {
				pCardSet := pGame.dealing(-1, 1)
				pGame.pDealerCard.AddCardSet(pCardSet)
				mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), pGame.pDealerCard.pCardSet.String(),
					pGame.pDealerCard.i32Point_1, pGame.pDealerCard.i32Point_2)
				if pGame.pDealerCard.IsSoftReached() {
					break
				}
			}
		}
	}

	mlog.Debug("庄家牌：%s (%d, %d)", pGame.pDealerCard.pCardSet.String(), pGame.pDealerCard.i32Point_1, pGame.pDealerCard.i32Point_2)
	pGame.GameLog.AddGameStep("庄家牌：%s (%d, %d)", pGame.pDealerCard.pCardSet.String(), pGame.pDealerCard.i32Point_1, pGame.pDealerCard.i32Point_2)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.iUserState == JD21_UserState_Normal {
			mlog.Debug("User[%d], CardSetCount[%d]", pUser.i32UserID, len(pUser.ary_pCardSet))
			for i := 0; i < len(pUser.ary_pCardSet); i++ {
				mlog.Debug("Set %d: Cards[%s], BetBase[%d], BetDouble[%d], bCompensated[%t]", i+JD21_IndexBase, pUser.ary_pCardSet[i].pCardSet.CodeString(),
					pUser.ary_pCardSet[i].i64BetBase, pUser.ary_pCardSet[i].i64BetDouble, pUser.ary_pCardSet[i].bCompensated)
				if !pUser.ary_pCardSet[i].bCompensated {
					iCmpResult := CompareCards(pGame.pDealerCard, pUser.ary_pCardSet[i])
					i64WinLose := int64(0)
					switch iCmpResult {
					case JD21_Compare_Even: // 平局
						i64WinLose = pUser.ary_pCardSet[i].TotalBet()
					case JD21_Compare_Dealer: // 莊家贏
						i64WinLose = -pUser.ary_pCardSet[i].TotalBet()
					case JD21_Compare_Player: // 玩家贏
						if pUser.ary_pCardSet[i].bIsBlackJack {
							i64WinLose = pUser.ary_pCardSet[i].TotalBet() + pUser.ary_pCardSet[i].TotalBet()*3/2
						} else {
							i64WinLose = pUser.ary_pCardSet[i].TotalBet() * 2
						}
					}

					// 寫分
					if i64WinLose > 0 {
						i64WinLose -= pUser.ary_pCardSet[i].TotalBet()
					}
					i64MoneyChange, i64Tax := pGame.GetUserWinAndTax(i64WinLose)
					i64Flow := int64(0)
					if i64MoneyChange > 0 {
						i64Flow = i64MoneyChange
					} else {
						i64Flow = -i64MoneyChange
					}
					pGame.UpdateUserScoreFunc(pUser.User.UserID, i64MoneyChange, i64Tax)

					if i64MoneyChange > 0 {
						pUser.ary_pCardSet[i].i64WinLose = i64MoneyChange + pUser.ary_pCardSet[i].TotalBet()
					} else {
						pUser.ary_pCardSet[i].i64WinLose = i64MoneyChange
					}
					pUser.i64TotalWinLose += pUser.ary_pCardSet[i].i64WinLose

					if iCmpResult == JD21_Compare_Even {
						pUser.i64Money += pUser.ary_pCardSet[i].TotalBet()
					} else {
						if i64MoneyChange > 0 {
							pUser.i64Money += pUser.ary_pCardSet[i].TotalBet() + i64MoneyChange
						}
					}
					pUser.i64MoneyChange += i64MoneyChange
					pUser.i64Tax += i64Tax
					pUser.ary_pCardSet[i].bCompensated = true

					mlog.Debug("玩家[%d], 手牌 %d：%s (%d, %d), 输赢 %s, 税收 %s, 流水 %s, 结算后金钱 %s", pUser.i32UserID, i+JD21_IndexBase,
						pUser.ary_pCardSet[i].pCardSet.String(), pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2,
						Global.MoneyString(i64MoneyChange), Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
					pGame.GameLog.AddGameStep("玩家[%d], 手牌 %d：%s (%d, %d), 输赢 %s, 税收 %s, 流水 %s, 结算后金钱 %s", pUser.i32UserID, i+JD21_IndexBase,
						pUser.ary_pCardSet[i].pCardSet.String(), pUser.ary_pCardSet[i].i32Point_1, pUser.ary_pCardSet[i].i32Point_2,
						Global.MoneyString(i64MoneyChange), Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))

					pUser.ary_pCardSet[i].bCompensated = true
				}
			}
		}
	}

	//----- 傳送訊息給 client -----
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.JDESYD_GameMessageClassID_JDESYD_Settlement_Dealer)

	pData := new(netproto.JDESYD_Msg_Settlement_Dealer)
	pData.DealerCards = pGame.cardSetInfoPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData

	pGame.BoardcastMessage(pMsg)

	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.ary_pUserData[iSeat].iUserState == JD21_UserState_EmptySeat || pGame.ary_pUserData[iSeat].iUserType == gamelogic.UserTypeRebot ||
			pGame.ary_pUserData[iSeat].bOffline {
			pGame.ary_pUserData[iSeat].iPlayingState = JD21_UserPlayState_Settlement_DealerOK
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}
	}
}

// 莊家開牌/結算動畫播放完成
func (pGame *JdesydLogic) message_Settlement_Dealer_OK(pUser *JdesydUser) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != JD21_State_Settlement_Dealer {
		return
	}

	if pUser.iPlayingState != JD21_UserPlayState_Playing_Done {
		return
	}

	pUser.iPlayingState = JD21_UserPlayState_Settlement_DealerOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	if pGame.uiPlayerReady == pGame.uiAllPlayerReady {
		mlog.Debug("所有玩家莊家開牌結算動畫都已播放完成")
		pGame.setGameState(JD21_State_GameOver, 0)
		pGame.gameOver()
	}
}

// 遊戲結束
func (pGame *JdesydLogic) gameOver() {
	mlog.Debug("游戏结束")
	pGame.GameLog.AddGameStep("游戏结束")

	map_i64ScoreRet := make(map[int32]int64)
	map_i64CouponRet := make(map[int32]int64)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.ary_pUserData[iSeat]
		if pUser.User != nil {
			if pUser.i64MoneyChange >= 0 {
				pUser.User.FlowAmount += pUser.i64MoneyChange
			} else {
				pUser.User.FlowAmount += -pUser.i64MoneyChange
			}

			map_i64ScoreRet[pUser.i32UserID] = pUser.i64MoneyChange
			pGame.GameLog.Context().UserUpdateScore(pUser.i32UserID, pUser.i64MoneyChange)

			map_i64CouponRet[pUser.i32UserID] = 0

			if g_bEnablePoolControl {
				if !pUser.User.IsRobot() {
					pGame.pWinLose.GameDBData.BetCount += pUser.iSumBetCount
					pGame.pWinLose.GameDBData.BetMoney += pUser.i64SumBetMoney
					pGame.pWinLose.i64PoolDifference -= pUser.i64MoneyChange + pUser.i64Tax
				}
			}
		}
	}

	if g_bEnablePoolControl {
		g_pSysPoolWinLose.GameDBData.GameRound++
		g_pSysPoolWinLose.GameDBData.BetCount += pGame.pWinLose.GameDBData.BetCount
		g_pSysPoolWinLose.GameDBData.BetMoney += pGame.pWinLose.GameDBData.BetMoney
		strCheckLog, strUpdateLog, bInit := checkSysPool(pGame.iDealType, pGame.pWinLose)
		if strCheckLog != "" {
			mlog.Debug("%s", strCheckLog)
			pGame.GameLog.AddGameStep("%s", strCheckLog)
		}
		if strUpdateLog != "" {
			mlog.Debug("%s", strUpdateLog)
			if !bInit {
				pGame.GameLog.AddGameStep("%s", strUpdateLog)
			}
		}
		pGame.UpdateSysPoolWinLose()
	}

	pGame.GameLog.GameEnd()

	pGame.GameEndFunc()

	iUserCount := 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			pGame.ary_pUserData[iSeat].Init()
		} else {
			if pGameUser.IsOffLine() {
				pGame.UserLeave(pGameUser)
				pGame.ary_pUserData[iSeat].Init()
			} else {
				i32LoginMoney := pGame.GetTable().Room.LoginMoney
				if pGame.ary_pUserData[iSeat].i64Money < int64(i32LoginMoney) {
					pGame.LogicBase.KickUserFunc(pGameUser.UserID, "您的金币不足!请补充金币后再进行游戏.")
					pGame.ary_pUserData[iSeat].Init()
				} else {
					pGame.ary_pUserData[iSeat].ResetGameData()
					pGame.ary_pUserData[iSeat].iUserState = JD21_UserState_Normal
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

// 玩家站立
func (pJLB *JdesydLogic) UserFree(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}
	iSeat := int(pGameUser.SeatNo)
	if iSeat < 0 || iSeat >= pJLB.iSeatCount {
		return
	}
	pUser := pJLB.ary_pUserData[pGameUser.SeatNo]
	if pJLB.iGameState >= JD21_State_Betting && pJLB.iGameState <= JD21_State_Settlement_Dealer && pUser.iUserState != JD21_UserState_Watch {
		return
	}
	pJLB.ary_pUserData[iSeat].Init()
	if !pGameUser.IsRobot() {
		mlog.System("UserFree %s", pGameUser)
	}
	pJLB.sendAllUserInfo()
	defer func() {
		// check all user leave, close game ...
		var hasuser bool
		for _, user := range pJLB.ary_pUserData {
			if user.User != nil {
				hasuser = true
				break
			}
		}

		if !hasuser {
			pJLB.initData()
		}
	}()
	pJLB.LogicBase.UserFree(pGameUser)
}

// 玩家離線
func (pGame *JdesydLogic) UserLeave(pGameUser *gamelogic.GameUser) bool {
	if pGameUser == nil {
		return false
	}

	mlog.Debug("Table[%d], UserLeave, User[%d], Seat[%d], IsOffLine()[%t]", pGame.TableNo(), pGameUser.UserID, pGameUser.SeatNo, pGameUser.IsOffLine())

	if pGameUser.SeatNo >= 0 && pGameUser.SeatNo < int32(pGame.iSeatCount) {
		pGame.ary_pUserData[int(pGameUser.SeatNo)].bOffline = true
	}

	pUser := pGame.ary_pUserData[pGameUser.SeatNo]
	if pGame.iGameState >= JD21_State_Betting && pGame.iGameState <= JD21_State_Settlement_Dealer && pUser.iUserState != JD21_UserState_Watch {
		return false
	}

	if pGameUser.SeatNo >= 0 && pGameUser.SeatNo < int32(pGame.iSeatCount) {
		pGame.ary_pUserData[int(pGameUser.SeatNo)].Init()
		pGame.sendAllUserInfo()
	}

	defer func() {
		// check all user leave, close game ...
		bHasUser := false
		for _, pUser := range pGame.ary_pUserData {
			if pUser.User != nil {
				bHasUser = true
				break
			}
		}

		if !bHasUser {
			pGame.initData()
		}
	}()

	return pGame.LogicBase.UserLeave(pGameUser)
}

// 玩家斷線重連
func (pGame *JdesydLogic) UserRecover(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserRecover, User[%d], Seat[%d], IsOffLine()[%t]", pGame.TableNo(), pGameUser.UserID, pGameUser.SeatNo, pGameUser.IsOffLine())

	pGame.DynamicGameLogicBase.UserRecover(pGameUser)

	if pGameUser.SeatNo >= 0 && pGameUser.SeatNo < int32(pGame.iSeatCount) {
		pGame.ary_pUserData[int(pGameUser.SeatNo)].bOffline = false
	}

	pGame.sendGameInfo(pGameUser)
	pGame.sendAllUserInfo()

	switch pGame.iGameState {
	case JD21_State_Betting: // 押底注
		pGame.startBetting(true, pGameUser)
	case JD21_State_StopBetting: // 停止下注
		pGame.stopBetting(true, pGameUser)
	case JD21_State_CutCardsAnimation: // 切牌動畫
		pGame.sendCutCardsResult(true, pGameUser)
	case JD21_State_Insurance: // 保險
		pGame.sendInsuranceOrNot(true, pGameUser)
	case JD21_State_Playing: // 玩家操作
		pGame.sendStartPlaying(true, pGameUser)
	}
}

// 依設定發牌，i32Seat < 0 表示為莊家
func (pGame *JdesydLogic) dealing(i32Seat int32, iCardCount int) *Pokerset.PokerCardSet {
	pRetCardSet := Pokerset.NewPokerCardSet()

	if i32Seat >= int32(pGame.iSeatCount) {
		return pRetCardSet
	}

	if iCardCount <= 0 || iCardCount > pGame.pGamePoker.CardCount() {
		return pRetCardSet
	}

	switch pGame.iDealType {
	case JD21_DealType_Assigned: // 指定發牌
		if i32Seat < 0 {
			// 莊家
			if pGame.pTempDealerCardSet.CardCount() >= iCardCount {
				pSet := pGame.pTempDealerCardSet.Deal(iCardCount)
				pRetCardSet.SetCardSet(pSet.CodeString())
			} else {
				iCount := iCardCount - pGame.pTempDealerCardSet.CardCount()
				pRetCardSet.SetCardSet(pGame.pTempDealerCardSet.CodeString())
				pSet := pGame.pGamePoker.Deal(iCount)
				pRetCardSet.AddCardSet(pSet)
			}
		} else {
			// 玩家
			if pGame.ary_pTempUserCardSet[i32Seat].CardCount() >= iCardCount {
				pSet := pGame.ary_pTempUserCardSet[i32Seat].Deal(iCardCount)
				pRetCardSet.SetCardSet(pSet.CodeString())
			} else {
				iCount := iCardCount - pGame.ary_pTempUserCardSet[i32Seat].CardCount()
				pRetCardSet.SetCardSet(pGame.ary_pTempUserCardSet[i32Seat].CodeString())
				pSet := pGame.pGamePoker.Deal(iCount)
				pRetCardSet.AddCardSet(pSet)
			}
		}
	default:
		pSet := pGame.pGamePoker.Deal(iCardCount)
		pRetCardSet.SetCardSet(pSet.CodeString())
	}

	return pRetCardSet
}

// 機器人操作
func (pGame *JdesydLogic) robotPlaying(pUser *JdesydUser) {
	if pGame.iGameState != JD21_State_Playing && pGame.iGameState != JD21_State_PlayingDelay {
		return
	}

	if pUser == nil {
		return
	}

	if pUser.i32SeatNo != pGame.iWhosTurn {
		return
	}

	iAction, iResult := JD21_UserAction_Stand, JD21_Playing_Success
	iRnd := util.RndInt(0, 100)
	bOperationDone := false

	// 莊家明牌有10或是A時投降
	pCard := pGame.pDealerCard.pCardSet.GetSingleCard(0, false)
	if g_map_i32CardNoValue[pCard.I32Number] == 1 || g_map_i32CardNoValue[pCard.I32Number] == 10 {
		if !pUser.bSurrendered {
			if len(pUser.ary_pCardSet) == 1 && pUser.CardSetUsedNow().pCardSet.CardCount() == JD21_BasicCardCount {
				if iRnd < pUser.pRobot.iRate_Surrender && !pUser.CardSetUsedNow().IsSoftReached() {
					iAction = JD21_UserAction_Surrender
					iResult = pGame.userAction_Surrender(pUser)
					if iResult == JD21_Playing_Success {
						bOperationDone = true
					}
				}
			}
		}
	}

	if !bOperationDone {
		// 可分牌時分牌
		if pUser.CardSetUsedNow().bCanSplit {
			if iRnd < pUser.pRobot.iRate_Split {
				iAction = JD21_UserAction_Split
				iResult = pGame.userAction_Split(pUser)
				if iResult == JD21_Playing_Success {
					bOperationDone = true
				}
			}
		}
	}

	if !bOperationDone {
		// 可加倍時加倍
		if pUser.CardSetUsedNow().bCanDouble {
			if iRnd < pUser.pRobot.iRate_Double {
				iAction = JD21_UserAction_Double
				iResult = pGame.userAction_Double(pUser)
				if iResult == JD21_Playing_Success {
					bOperationDone = true
				}
			}
		}
	}

	if !bOperationDone {
		if !pUser.CardSetUsedNow().IsSoftReached() {
			// 小於17點時要牌
			if iRnd < pUser.pRobot.iRate_Hit {
				iAction = JD21_UserAction_Hit
				iResult = pGame.userAction_Hit(pUser)
				if iResult == JD21_Playing_Success {
					bOperationDone = true
				}
			}
		} else {
			// 大於等於17點時要牌機率
			if !pUser.CardSetUsedNow().bIsBlackJack && !pUser.CardSetUsedNow().bIsBusted {
				if iRnd < pUser.pRobot.iRate_Hit_ReachSoft {
					iAction = JD21_UserAction_Hit
					iResult = pGame.userAction_Hit(pUser)
					if iResult == JD21_Playing_Success {
						bOperationDone = true
					}
				}
			}
		}
	}

	// 停牌
	if !bOperationDone {
		pUser.Stand()
		iAction = JD21_UserAction_Stand
		iResult = JD21_Playing_Success

		mlog.Debug("玩家[%d] 停牌", pUser.i32UserID)
		pGame.GameLog.AddGameStep("玩家[%d] 停牌", pUser.i32UserID)
	}

	pGame.sendUserPlayingResult(pUser, iAction, iResult)

	if pUser.CardSetUsedNow().bIsBusted {
		pUser.CardSetUsedNow().i64BetBase = 0
		pUser.CardSetUsedNow().i64BetDouble = 0
		pUser.CardSetUsedNow().bOerationDone = true
	}

	// 檢查是否為最後一局
	if !pGame.bIsLastGame {
		if pGame.IsLastGame() {
			pGame.sendLastGame()
		}
	}

	if iResult == JD21_Playing_Success {
		// 延遲 0.5 秒進行下一個動作
		pGame.setGameState(JD21_State_PlayingDelay, pGame.map_i64Timeout[JD21_Timeout_PlayingDelay])
	}
}

// 設定遊戲狀態
func (pGame *JdesydLogic) setGameState(iGameState int, i64Timeout int64) {
	if iGameState < JD21_State_WaitingPlayer || iGameState > JD21_State_GameOver {
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

	mlog.Debug("TableNo[%d],【state】[%s], timeout[%d], duration[%d], lastmsg[%d]", pGame.TableNo(), g_map_strGameState_debug[pGame.iGameState],
		pGame.i64Timeout, pGame.i64TimeDuration, pGame.i64LastMessageTime)
}

// 依 UserID 取得 UserData
func (pGame *JdesydLogic) GetUserDataByUserID(i32UserID int32) *JdesydUser {
	iSeatNo := -1

	for i := 0; i < pGame.iSeatCount; i++ {
		if i32UserID == pGame.ary_pUserData[i].i32UserID {
			iSeatNo = i
			break
		}
	}

	if iSeatNo >= 0 && iSeatNo < pGame.iSeatCount {
		return pGame.ary_pUserData[iSeatNo]
	} else {
		return nil
	}
}
