package dczz

import (
	"fmt"
	"math"
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
	pGamePoker                 *pokerset.PokerCardSet   // 系統牌盒
	iSeatCount                 int                      // 座位數
	pDealerCard                *pokerset.PokerCardSet   // 莊家手牌資料
	aryPUserData               []*User                  // 玩家資料
	pGameTicker                *time.Ticker             // 遊戲計時器
	iGameState                 int                      // 遊戲狀態
	i64TimeDuration            int64                    // Ticker 計數值
	i64Timeout                 int64                    // 各階段 timeout
	i64LastMessageTime         int64                    // 上一次訊息時間
	mapI64Timeout              map[int]int64            // timeout 時間設定
	aryStrPokerManual          []string                 // 遊戲說明
	aryStrJPManual             []string                 // 彩金說明
	aryI32Chip                 []int32                  // 籌碼面額
	uiAllPlayerReady           uint                     // 各玩家均已 Ready，此值設定後就不再變更
	uiPlayerReady              uint                     // 標示各玩家是否 Ready
	iUserCount                 int                      // 遊戲中玩家人數
	pRobotSetting              *RobotSetting            // 機器人設定
	i32DealerCardPoint         int32                    // 莊家點數
	bGameOver                  bool                     // 遊戲結束
	iDealType                  int                      // 發牌方式
	pTestDealerCards           *pokerset.PokerCardSet   // 測試用莊家手牌
	aryPTestUserCards          []*pokerset.PokerCardSet // 測試用玩家手牌
	pTestDealerCardBox         *pokerset.PokerCardSet   // 測試用莊家手牌牌盒
	aryPTestUserCardBox        []*pokerset.PokerCardSet // 測試用玩家手牌牌盒
	bAllUserDoneAtWarSelection bool                     // 所有玩家是否都已在戰爭階段完成動作
}

// NewLogic 建立遊戲
func NewLogic(table *gamelogic.Table, tfManager *event.TimerFuncManager, configs *config.Group) *Logic {
	pGame := new(Logic)
	pGame.DynamicGameLogicBase = gamelogic.NewDynamicGameLogicBase(table, tfManager, configs)

	pGame.i64TimeDuration = 0
	tfManager.NewMilliSecTicker(100, func() {
		pGame.i64TimeDuration++
		pGame.gameTickerFunc()
	})

	pGame.setGameState(GameState_WaitingPlayer, 0)

	pGame.mapI64Timeout = make(map[int]int64)
	pGame.pRobotSetting = NewRobotSetting()

	g_pSysPoolWinLose = new(SysPoolWinLose)

	pGame.getGameConfig(configs)

	pGame.pGamePoker = pokerset.NewPokerCardSet()
	pGame.pGamePoker.NewSet(PokerSetCount, 0)
	pGame.pGamePoker.Shuffle()

	pGame.pDealerCard = pokerset.NewPokerCardSet()
	pGame.pTestDealerCardBox = pokerset.NewPokerCardSet()

	pGame.aryPUserData = make([]*User, pGame.iSeatCount)
	pGame.aryPTestUserCardBox = make([]*pokerset.PokerCardSet, pGame.iSeatCount)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.aryPUserData[iSeat] = NewUser()
		pGame.aryPTestUserCardBox[iSeat] = pokerset.NewPokerCardSet()
	}

	// 各玩家均已 Ready，此值設定後就不再變更
	uiReady := uint(1)
	for i := 0; i < pGame.iSeatCount; i++ {
		uiReady *= 2
	}
	pGame.uiAllPlayerReady = uiReady - 1

	// 遊戲說明
	pGame.aryStrPokerManual = append(pGame.aryStrPokerManual, "庄家和每位下注的客人会获发一张牌，单凭比较庄家和玩家的点数大小(A最大)分胜负。如果玩家的点数大于庄家，则可获相等于下注金额的赔彩。牌的花色于此游戏中无效，故当庄家和玩家的点数相同时，玩家可作以下两个选择：")
	pGame.aryStrPokerManual = append(pGame.aryStrPokerManual, "弃牌：玩家可取回下注金额的一半。")
	pGame.aryStrPokerManual = append(pGame.aryStrPokerManual, "进入赌场战争：玩家需额外再加上相等于原注的金额进入赌场战争，然后庄荷会消去三张牌，再发牌给客人，及后再消去三张牌，发牌给自己，如果庄家点数较大，玩家会输掉原注和额外加注。如玩家点数较大，则只会赢取与额外加注相等的金额。如两者点数再次相同，玩家则可赢取原注与额外加注的两份赔彩。")

	// 彩金說明
	pGame.aryStrJPManual = append(pGame.aryStrJPManual, "平手1賠10")

	pGame.initData()

	pGame.RegisterBeforeGameStartFunc(func() {})
	pGame.RegisterAfterGameStartFunc(func() {})
	pGame.RegisterBeforeGameEndFunc(func() {})
	pGame.RegisterAfterGameEndFunc(func() {})

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

	for i := 0; i < pGame.iSeatCount; i++ {
		pGame.aryPUserData[i].ResetGameData()
	}
	pGame.pDealerCard.SetCardSet("")
	pGame.i32DealerCardPoint = 0
	pGame.bGameOver = false
	pGame.pTestDealerCardBox.SetCardSet(pGame.pTestDealerCards.CodeString())
	pGame.bAllUserDoneAtWarSelection = false

	pGame.iUserCount = 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGameUser := pGame.GetUserBySeatNo(int32(iSeat))
		if pGameUser == nil {
			pGame.aryPUserData[iSeat].Init()
		} else {
			if pGame.aryPUserData[iSeat].iUserState != UserState_EmptySeat {
				pGame.aryPUserData[iSeat].iUserState = UserState_Normal
			}
			pGame.UserReady(pGameUser)
			pGame.iUserCount++
		}

		pGame.aryPTestUserCardBox[iSeat].SetCardSet(pGame.aryPTestUserCards[iSeat].CodeString())
	}

	if pGame.iUserCount <= 0 {
		pGame.setGameState(GameState_WaitingPlayer, 0)

		mlog.Debug("等待游戏开始")
		pGame.GameLog.AddGameStep("等待游戏开始")
	} else {
		pGame.setGameState(GameState_WaitingStart, pGame.mapI64Timeout[Timeout_WaitingStart])
		pGame.waitingGameBegin()
	}

	pGame.uiPlayerReady = 0
}

// GameMessage 遊戲訊息
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
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_ShuffleOK): // C2S 洗牌動畫播放完成
		pGame.messageShuffleOK(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet): // C2S 玩家下注
		pGame.messageBet(pUser, pMsg)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_StopBettingAnimationOK): // C2S 玩家停止下注動畫播放完成
		pGame.messageStopBettingAnimationOK(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_DealAnimationOK_Normal): // C2S 一般發牌動畫播放完成
		pGame.messageDealOKNormal(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_Settlement_Normal_OK): // C2S 一般結算動畫播放完成
		pGame.messageSettlementNormalOK(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection): // C2S 戰爭階段選項
		pGame.messageWarSelection(pUser, pMsg)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseAnimationOK): // C2S 進入戰爭階段動畫播放完成
		pGame.messageWarPhaseAnimationOK(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_DealAnimationOK_War): // C2S 戰爭階段發牌動畫播放完成
		pGame.messageDealOKWar(pUser)
	case int32(netproto.DCZZ_GameMessageClassID_DCZZ_Settlement_War_OK): // C2S 戰爭階段結算動畫播放完成
		pGame.messageSettlementWarOK(pUser)
	}
}

// PiSaiMessage 比賽訊息
func (pGame *Logic) PiSaiMessage(pUser *gamelogic.GameUser, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}
}

// UserEnter 玩家進入遊戲
func (pGame *Logic) UserEnter(pGameUser *gamelogic.GameUser, i32TableNo int32, i32SeatNo int32) {
	if pGameUser == nil {
		return
	}

	// 自 DB 取得遊戲輸贏結果
	if g_bEnablePoolControl {
		if !g_bGetWinLoseData {
			pGame.GetSysPoolWinLose()
		}
	}

	mlog.Debug("Table[%d], UserEnter, GameState[%s], User[%d], Seat[%d]",
		pGame.TableNo(), g_map_strGameState_debug[pGame.iGameState], pGameUser.UserID, pGameUser.SeatNo)

	if pGameUser.IsRobot() {
		if pGame.aryPUserData[i32SeatNo].i32UserID != pGameUser.UserID {
			for i := 0; i < len(pGame.aryPUserData); i++ {
				if pGame.aryPUserData[i].i32UserID == pGameUser.UserID {
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
func (pGame *Logic) gameTickerFunc() {
	if pGame.i64Timeout > 0 && pGame.i64TimeDuration-pGame.i64LastMessageTime > pGame.i64Timeout {
		switch pGame.iGameState {
		case GameState_WaitingStart: // 等待遊戲開始
			pGame.GameStartFunc()
			pGame.setGameState(GameState_Betting, pGame.mapI64Timeout[Timeout_Betting])
			pGame.startBetting(false, nil)
		case GameState_Betting: // 押注
			pGame.setGameState(GameState_StopBetting, pGame.mapI64Timeout[Timeout_StopBettingAnimation])
			pGame.stopBetting()
		case GameState_StopBetting: // 停止下注
			if pGame.pGamePoker.CardCount() <= g_iShuffleRestCardCount {
				pGame.setGameState(GameState_Shuffle, pGame.mapI64Timeout[Timeout_Shuffle])
				pGame.shuffle(false, nil)
			} else {
				pGame.setGameState(GameState_Deal_Normal, pGame.mapI64Timeout[Timeout_Deal_Normal])
				pGame.dealNormal()
			}
		case GameState_Shuffle: // 洗牌
			pGame.setGameState(GameState_Deal_Normal, pGame.mapI64Timeout[Timeout_Deal_Normal])
			pGame.dealNormal()
		case GameState_Deal_Normal: // 一般發牌
			pGame.setGameState(GameState_Settlement_Normal, pGame.mapI64Timeout[Timeout_Settlement_Normal])
			pGame.settlementNormal()
		case GameState_Settlement_Normal: // 一般結算
			pGame.checkIfGameOver()
			if pGame.bGameOver {
				// 遊戲結束
				pGame.sendGameOver(false)
			} else {
				// 進入戰爭階段
				pGame.setGameState(GameState_WarSelection, pGame.mapI64Timeout[Timeout_WarSelection])
				pGame.enterWarPhase(false, nil)
			}
		case GameState_WarSelection: // 選擇是否進入戰爭
			for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
				if pGame.iGameState != GameState_WarSelection {
					break
				}

				pUser := pGame.aryPUserData[iSeat]
				if pUser.iUserState == UserState_Normal && !pUser.bOperationDone {
					if !pUser.bGoForWar {
						pGame.userActionWarSelection(pUser, WarSelection_Fold)
					}
				}
			}
		case GameState_WarPhaseAnimation: // 進入戰爭階段動畫
			pGame.setGameState(GameState_Deal_War, pGame.mapI64Timeout[Timeout_Deal_War])
			pGame.dealWar()
		case GameState_Deal_War: // 戰爭發牌
			pGame.setGameState(GameState_Settlement_War, pGame.mapI64Timeout[Timeout_Settlement_War])
			pGame.settlementWar()
		case GameState_Settlement_War: // 戰爭結算
			pGame.sendGameOver(false)
		case GameState_GameOver: // 遊戲結束
			if pGame.bAllUserDoneAtWarSelection {
				pGame.bAllUserDoneAtWarSelection = false
				pGame.sendGameOver(false)
			} else {
				pGame.gameOver()
			}
		}
	}

	// 機器人動作
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == UserState_Normal {
			if pUser.i64Timeout > 0 && pGame.i64TimeDuration-pUser.i64LastMessageTime > pUser.i64Timeout {
				switch pGame.iGameState {
				case GameState_Betting: // 押注
					// 下注類別
					iBetType := util.RndInt(0, BetType_Max)
					if pUser.pRobot.ary_pBetParam[iBetType].iBetTimes >= pUser.pRobot.ary_pBetParam[iBetType].iBetCount {
						pUser.i64Timeout = 5
						pUser.i64LastMessageTime = pGame.i64TimeDuration
					} else {
						// 可下注籌碼
						aryI64Chip := make([]int64, 0)
						for i := 0; i < len(pGame.aryI32Chip); i++ {
							if i < pUser.pRobot.ary_pBetParam[iBetType].iBetRange_Min-1 || i > pUser.pRobot.ary_pBetParam[iBetType].iBetRange_Max-1 {
								continue
							}

							if pUser.i64Money >= int64(pGame.aryI32Chip[i]) {
								aryI64Chip = append(aryI64Chip, int64(pGame.aryI32Chip[i]))
							}
						}

						if len(aryI64Chip) > 0 {
							iBetIdx := util.RndInt(0, len(aryI64Chip))
							iResult := pGame.userActionBet(pUser, iBetType, aryI64Chip[iBetIdx])
							if iResult == BetResult_Success {
								pGame.sendBetResult(pUser, iBetType, aryI64Chip[iBetIdx], iResult)
							}
						}

						pUser.pRobot.ary_pBetParam[iBetType].iBetTimes++

						if pUser.pRobot.ary_pBetParam[iBetType].iBetTimes >= pUser.pRobot.ary_pBetParam[iBetType].iBetCount {
							pUser.i64Timeout = 0
						} else {
							pUser.i64Timeout = 5
							pUser.i64LastMessageTime = pGame.i64TimeDuration
						}
					}
				case GameState_WarSelection: // 選擇是否進入戰爭
					if !pUser.bCanEnterWar || pUser.bIsFolded {
						continue
					}

					iRnd := util.RndInt(0, 100)
					if iRnd < pUser.pRobot.iRate_FoldAtWar {
						pGame.userActionWarSelection(pUser, WarSelection_Fold)
					} else {
						pGame.userActionWarSelection(pUser, WarSelection_War)
					}

					pUser.i64Timeout = 0
				}
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
		// 牌盒內張數低於此值時洗牌
		g_iShuffleRestCardCount = Count_CardCount_War * (pGame.iSeatCount + 1)

		//----- timeout -----
		pGame.mapI64Timeout[Timeout_WaitingStart] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_WaitingStart")) * 10
		pGame.mapI64Timeout[Timeout_Betting] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Betting")) * 10
		pGame.mapI64Timeout[Timeout_StopBettingAnimation] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_StopBettingAnimation")) * 10
		pGame.mapI64Timeout[Timeout_Shuffle] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Shuffle")) * 10
		pGame.mapI64Timeout[Timeout_Deal_Normal] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Deal_Normal")) * 10
		pGame.mapI64Timeout[Timeout_Settlement_Normal] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Settlement_Normal")) * 10
		pGame.mapI64Timeout[Timeout_WarSelection] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_WarSelection")) * 10
		pGame.mapI64Timeout[Timeout_WarPhaseBegin] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_WarPhaseBegin")) * 10
		pGame.mapI64Timeout[Timeout_Deal_War] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Deal_War")) * 10
		pGame.mapI64Timeout[Timeout_Settlement_War] = int64(cfgs.GetConfigInt("GameSetting", "Timeout_Settlement_War")) * 10

		// 系統自訂，不是從 gameserver 讀取的項目
		pGame.mapI64Timeout[Timeout_GameOver] = TimeDuration_GameOver * 10

		// 籌碼面額
		str = cfgs.GetConfigValue("GameSetting", "ChipValue")
		aryStrChip := strings.Split(str, ",")
		for i := 0; i < len(aryStrChip); i++ {
			i64chip, _ := strconv.ParseInt(aryStrChip[i], 10, 32)
			pGame.aryI32Chip = append(pGame.aryI32Chip, int32(i64chip))
		}

		if len(pGame.aryI32Chip) <= 0 {
			mlog.Fatal("讀取 GameSetting 失敗：沒有籌碼資料")
		}

		g_i64BetLimit_Base = int64(cfgs.GetConfigInt64("GameSetting", "BetLimit_Base")) // 底注額上限
		g_i64BetLimit_Side = int64(cfgs.GetConfigInt64("GameSetting", "BetLimit_Side")) // 平手額上限

		// 機器人參數
		pGame.pRobotSetting.iRate_FoldAtWar = cfgs.GetConfigInt("robot", "Rate_FoldAtWar") // 進入戰爭後棄牌機率

		// 下注參數
		// 底注
		var i64Param int64
		str = cfgs.GetConfigValue("robot", "BetRange_Base")
		aryStrSetting := strings.Split(str, ",")
		if len(aryStrSetting) < 3 {
			mlog.Error("[robot] 参数 BetRange_Base 设定错误：%s", str)
		} else {
			i64Param, _ = strconv.ParseInt(aryStrSetting[0], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Base].iBetRange_Min = int(i64Param)
			i64Param, _ = strconv.ParseInt(aryStrSetting[1], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Base].iBetRange_Max = int(i64Param)
			i64Param, _ = strconv.ParseInt(aryStrSetting[2], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Base].iDefaultBetCount = int(i64Param)
		}

		// 平手注
		str = cfgs.GetConfigValue("robot", "BetRange_Side")
		aryStrSetting = strings.Split(str, ",")
		if len(aryStrSetting) < 3 {
			mlog.Error("[robot] 参数 BetRange_Side 设定错误：%s", str)
		} else {
			i64Param, _ = strconv.ParseInt(aryStrSetting[0], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Side].iBetRange_Min = int(i64Param)
			i64Param, _ = strconv.ParseInt(aryStrSetting[1], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Side].iBetRange_Max = int(i64Param)
			i64Param, _ = strconv.ParseInt(aryStrSetting[2], 10, 32)
			pGame.pRobotSetting.ary_pBetParam[BetType_Side].iDefaultBetCount = int(i64Param)
		}

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
				g_i64WinLine = cfgs.GetConfigInt64("SystemPool", "WinLine")             // 贏錢線
				g_iPersonTimeLimit = cfgs.GetConfigInt("SystemPool", "PersonTimeLimit") // 判定人次

				mlog.System("水线控制功能开启")
				g_bGetSysPoolConfig = true
			}
		}
	}

	// 發牌器
	pGame.pTestDealerCards = pokerset.NewPokerCardSet()
	pGame.aryPTestUserCards = make([]*pokerset.PokerCardSet, pGame.iSeatCount)
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.aryPTestUserCards[iSeat] = pokerset.NewPokerCardSet()
	}

	pkcfgs := configs.GetConfig("poker")
	if pkcfgs == nil {
		pGame.iDealType = DealType_Natural
	} else {
		// 發牌方式
		iDealType := pkcfgs.GetConfigInt("poker", "DealType")
		if iDealType < DealType_Natural || iDealType >= DealType_Max {
			iDealType = DealType_Natural
		}
		pGame.iDealType = iDealType

		// 測試用莊家手牌
		strCards := pkcfgs.GetConfigValue("poker", "Dealer")
		pGame.pTestDealerCards.SetCardSet(strCards)
		example.invalid("測試用莊家手牌：%s", pGame.pTestDealerCards.CodeString())

		// 測試用玩家手牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			str = fmt.Sprintf("User%d", iSeat+1)
			strCards := pkcfgs.GetConfigValue("poker", str)
			pGame.aryPTestUserCards[iSeat].SetCardSet(strCards)
			example.invalid("測試用玩家手牌 %d：%s", iSeat+1, pGame.aryPTestUserCards[iSeat].CodeString())
		}
	}
}

// AddUser 加入新玩家
func (pGame *Logic) AddUser(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	iSeatNo := int(pGameUser.SeatNo)
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return
	}

	pUser := pGame.aryPUserData[iSeatNo]

	if pUser.iUserState == UserState_EmptySeat {
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

		if pGame.iGameState >= GameState_Betting && pGame.iGameState < GameState_GameOver {
			pUser.iUserState = UserState_Watch
		} else {
			pUser.iUserState = UserState_Normal
		}
	}
	pGame.iUserCount++

	if pGame.iUserCount > 0 && pGame.iGameState == GameState_WaitingPlayer {
		pGame.setGameState(GameState_WaitingStart, pGame.mapI64Timeout[Timeout_WaitingStart])
		pGame.waitingGameBegin()
	} else {
		if pGame.iGameState < GameState_Betting {
			pGame.waitingGameBegin()
		}
	}

	pGame.sendGameInfo(pUser)
	pGame.sendAllUserInfo()

	example.invalid("加入新玩家[%d], 遊戲狀態[%s], 座位[%d], 類別[%d], 玩家狀態[%s], 玩家總人數[%d]", pUser.i32UserID,
		g_map_strGameState_debug[pGame.iGameState], pUser.i32SeatNo, pUser.iUserType,
		g_map_strUserState_debug[pUser.iUserState], pGame.iUserCount)
	pGame.GameLog.AddGameStep("玩家[%d] 入座，狀態：%s", pGameUser.UserID, g_map_strUserState[pUser.iUserState])
}

// 等待遊戲開始
func (pGame *Logic) waitingGameBegin() {
	pGame.GameLog.GameBegin()
	mlog.Debug("等待開始下注")
	pGame.GameLog.AddGameStep("等待開始下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WaitingStart)
	pData := new(netproto.JDESYD_Msg_WaitingStart)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)
}

// GetUserDataByUserID 依 UserID 取得 UserData
func (pGame *Logic) GetUserDataByUserID(i32UserID int32) *User {
	iSeatNo := -1

	for i := 0; i < pGame.iSeatCount; i++ {
		if i32UserID == pGame.aryPUserData[i].i32UserID {
			iSeatNo = i
			break
		}
	}

	if iSeatNo >= 0 && iSeatNo < pGame.iSeatCount {
		return pGame.aryPUserData[iSeatNo]
	}

	return nil
}

// 傳送牌局資訊
func (pGame *Logic) sendGameInfo(pUser *User) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_GameInfo)

	pData := new(netproto.DCZZ_Msg_GameInfo)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間
	// 遊戲說明
	for i := 0; i < len(pGame.aryStrPokerManual); i++ {
		pData.PokerManual = append(pData.PokerManual, pGame.aryStrPokerManual[i])
	}
	// 籌碼面額
	for i := 0; i < len(pGame.aryI32Chip); i++ {
		pData.ChipValue = append(pData.ChipValue, pGame.aryI32Chip[i])
	}
	pData.DealerCards = pGame.cardSetPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}
	pData.BetBaseLimit = proto.Int64(g_i64BetLimit_Base) // 底注額上限
	pData.BetSideLimit = proto.Int64(g_i64BetLimit_Side) // 平手額上限

	pMsg.MsgData = pData
	pGame.SendMessageToUser(pUser.User, pMsg)
}

// 手牌牌組打包資料
func (pGame *Logic) cardSetPkg(pCardSet *pokerset.PokerCardSet) []*netproto.PokerCard {
	if pCardSet == nil {
		return nil
	}

	pRetSet := make([]*netproto.PokerCard, 0)

	// 牌組
	for iIdx := 0; iIdx < pCardSet.CardCount(); iIdx++ {
		pCard := new(netproto.PokerCard)
		pSCard := pCardSet.GetSingleCard(iIdx, true)
		if pSCard != nil {
			pCard.Color = proto.Int32(pSCard.I32Color)
			i32Number := pSCard.I32Number
			if pSCard.I32Number == pokerset.POKERCardNoA {
				i32Number = pokerset.POKERCardNo1
			}
			pCard.CardNumber = proto.Int32(i32Number)
		}

		pRetSet = append(pRetSet, pCard)
	}

	return pRetSet
}

// 玩家資訊打包資料，iSeat：此座位的玩家資訊，iSendTo：要送給哪位玩家
func (pGame *Logic) userInfoPkg(iSeat int) *netproto.DCZZ_Msg_UserInfo {
	if iSeat < 0 || iSeat >= pGame.iSeatCount {
		return nil
	}

	pUser := pGame.aryPUserData[iSeat]

	pUserInfo := new(netproto.DCZZ_Msg_UserInfo)

	pUserInfo.UserID = proto.Int32(pUser.i32UserID)                 // 玩家ID，若是空位時，玩家ID = -1
	pUserInfo.SeatNo = proto.Int32(int32(iSeat + IndexBase))        // 玩家座位，client 座位號從 1 開始
	pUserInfo.State = proto.Int32(int32(pUser.iUserState))          // 玩家狀態，0:空位，1:正常，2:觀戰，3:棄牌
	pUserInfo.Money = proto.Int64(pUser.i64Money)                   // 玩家金錢
	pUserInfo.HandCards = pGame.cardSetPkg(pUser.pHandCard)         // 玩家手牌
	pUserInfo.Bet_Base = proto.Int64(pUser.i64Bet_Base)             // 玩家底注
	pUserInfo.Bet_Side = proto.Int64(pUser.i64Bet_Side)             // 玩家平手注
	pUserInfo.Bet_Double = proto.Int64(pUser.i64Bet_Double)         // 玩家戰爭階段加注
	pUserInfo.TotalWinLose = proto.Int64(pUser.i64TotalWinLose)     // 總輸贏
	pUserInfo.IsTied = proto.Bool(pUser.bIsTied)                    // 是否跟莊家平手
	pUserInfo.GetTieAward = proto.Bool(pUser.bGetTieAward)          // 玩家是否於一般階段贏得平手賠付
	pUserInfo.CanEnterWar = proto.Bool(pUser.bCanEnterWar)          // 可否進入戰爭階段
	pUserInfo.IsFolded = proto.Bool(pUser.bIsFolded)                // 玩家是否已棄牌
	pUserInfo.FoldRecover = proto.Int64(pUser.i64FoldRecover)       // 玩家於戰爭階段選擇棄牌收回的注碼
	pUserInfo.Winner = proto.Int32(int32(pUser.iSettlementWinLose)) // 玩家輸贏

	return pUserInfo
}

// 傳送玩家資訊
func (pGame *Logic) sendAllUserInfo() {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_AllUserInfo)

	pData := new(netproto.DCZZ_Msg_AllUserInfo)
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)
}

// 開始下注
func (pGame *Logic) startBetting(bUserRecover bool, pUser *gamelogic.GameUser) {
	mlog.Debug("開始下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_StartBetting)
	pData := new(netproto.DCZZ_Msg_StartBetting)
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData
	if !bUserRecover {
		pGame.BoardcastMessage(pMsg)
	} else {
		if pUser != nil {
			pGame.SendMessageToUser(pUser, pMsg)
		}
	}

	pGame.GameLog.AddGameStep("開始下注")

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		pUser.iPlayingState = UserPlayState_Betting
		if pUser.iUserState == UserState_Normal {
			// 自動幫玩家下底注
			mlog.Debug("玩家[%d] 金錢 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] 金錢 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Money))
			mlog.Debug("pUser.i64Bet_Base = %d", pUser.i64Bet_Base)

			if pUser.i64Bet_Base <= 0 {
				pUser.Bet(BetType_Base, int64(pGame.aryI32Chip[0]))
				pGame.sendBetResult(pUser, BetType_Base, int64(pGame.aryI32Chip[0]), BetResult_Auto)
			}
		}
	}

	if !bUserRecover {
		// 機器人隨機下注
		for i := 0; i < pGame.iSeatCount; i++ {
			pUser := pGame.aryPUserData[i]
			if pUser.iUserType == gamelogic.UserTypeRebot {
				pUser.SetRobotTimeout(pGame.mapI64Timeout[Timeout_Betting])
				if pUser.i64Timeout > 0 {
					pUser.i64LastMessageTime = pGame.i64TimeDuration

					for i := 0; i < len(pUser.pRobot.ary_pBetParam); i++ {
						pUser.pRobot.ary_pBetParam[i].iBetTimes = 0
						if pUser.pRobot.ary_pBetParam[i].iDefaultBetCount > 0 {
							pUser.pRobot.ary_pBetParam[i].iBetCount = util.RndInt(0, pUser.pRobot.ary_pBetParam[i].iDefaultBetCount+1)
						}
					}
				}
			}
		}
	}
}

// 玩家下注結果
func (pGame *Logic) sendBetResult(pUser *User, iBetType int, i64Bet int64, iResult int) {
	if pUser == nil {
		return
	}

	iBetResult := iResult
	if iResult >= BetResult_Success {
		iBetResult = BetResult_Success
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_BetResult)

	pData := new(netproto.DCZZ_Msg_BetResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(IndexBase)) // 下注玩家座位
	pData.BetType = proto.Int32(int32(iBetType))                   // 下注類別
	pData.Bet = proto.Int64(i64Bet)                                // 下注金額
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                // 下注後玩家剩餘金錢
	pData.Result = proto.Int32(int32(iBetResult))                  // 下注結果

	pMsg.MsgData = pData
	if iBetResult == BetResult_Success {
		pGame.BoardcastMessage(pMsg)
	} else {
		pGame.SendMessageToUser(pUser.User, pMsg)
	}

	str := ""
	switch iBetType {
	case BetType_Base: // 底注
		str = fmt.Sprintf("%s總額 %s", g_map_strBetType[iBetType], Global.MoneyString(pUser.i64Bet_Base))
	case BetType_Side: // 平手注
		str = fmt.Sprintf("%s總額 %s", g_map_strBetType[iBetType], Global.MoneyString(pUser.i64Bet_Side))
	}

	if iResult == BetResult_Auto {
		mlog.Debug("玩家[%d] 自動下%s %s, %s, 剩餘金錢 %s", pUser.i32UserID, g_map_strBetType[iBetType],
			Global.MoneyString(i64Bet), str, Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 自動下%s %s, %s, 剩餘金錢 %s", pUser.i32UserID, g_map_strBetType[iBetType],
			Global.MoneyString(i64Bet), str, Global.MoneyString(pUser.i64Money))
	} else {
		mlog.Debug("玩家[%d] 下%s %s, 結果：%s, %s, 剩餘金錢 %s", pUser.i32UserID, g_map_strBetType[iBetType],
			Global.MoneyString(i64Bet), g_map_strBetResult[iResult], str, Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 下%s %s, 結果：%s, %s, 剩餘金錢 %s", pUser.i32UserID, g_map_strBetType[iBetType],
			Global.MoneyString(i64Bet), g_map_strBetResult[iResult], str, Global.MoneyString(pUser.i64Money))
	}
}

// 訊息處理 - 玩家下注
func (pGame *Logic) messageBet(pUser *User, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	msgInfo := &netproto.DCZZ_Msg_Bet{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	iBetType := int(msgInfo.GetBetType())
	i64Bet := msgInfo.GetBet()

	if pGame.iGameState != GameState_Betting {
		pGame.sendBetResult(pUser, iBetType, i64Bet, BetResult_BettingStopped)
		return
	}

	if pUser.iPlayingState != UserPlayState_Betting {
		pGame.sendBetResult(pUser, iBetType, i64Bet, BetResult_BettingStopped)
		return
	}

	iResult := pGame.userActionBet(pUser, iBetType, i64Bet)
	pGame.sendBetResult(pUser, iBetType, i64Bet, iResult)
}

// 玩家動作 - 下注
func (pGame *Logic) userActionBet(pUser *User, iBetType int, i64Bet int64) int {
	if pUser == nil {
		return BetResult_BetFailed
	}

	if iBetType < BetType_Base || iBetType >= BetType_Max {
		return BetResult_BetFailed
	}

	iResult := BetResult_Success

	// 餘額不足
	if pUser.i64Money < i64Bet {
		return BetResult_NotEnoughMoney
	}

	// 押注金額<=0
	if i64Bet <= 0 {
		return BetResult_BetSmallThanZero
	}

	// 總押注額超過上限
	switch iBetType {
	case BetType_Base: // 底注
		if pUser.i64Bet_Base+i64Bet > g_i64BetLimit_Base {
			return BetResult_OverUpperLimit
		}
	case BetType_Side: // 平手注
		if pUser.i64Bet_Side+i64Bet > g_i64BetLimit_Side {
			return BetResult_OverUpperLimit
		}
	}

	// 押注後剩餘金錢不足底注的兩倍
	if pUser.i64Money-i64Bet < pUser.i64Bet_Base*Times_BetDouble {
		return BetResult_OverUpperLimit
	}

	pUser.Bet(iBetType, i64Bet)

	return iResult
}

// 播放停止下注動畫
func (pGame *Logic) stopBetting() {
	mlog.Debug("停止下注")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_StopBetting)

	pGame.BoardcastMessage(pMsg)

	pGame.setUserReadyPlayingState(UserPlayState_StopBettingOK)
	pGame.checkAllUserReady(UserPlayState_StopBettingOK)
}

// 訊息處理 - 停止下注動畫播放完成
func (pGame *Logic) messageStopBettingAnimationOK(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_StopBetting {
		return
	}

	if pUser.iPlayingState != UserPlayState_Betting {
		return
	}

	mlog.Debug("User[%d] 停止下注動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_StopBettingOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_StopBettingOK)
}

// 設定玩家回傳訊息準備好狀態
func (pGame *Logic) setUserReadyPlayingState(iPlayState int) {
	pGame.uiPlayerReady = 0

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState != UserState_Normal || pUser.iUserType == gamelogic.UserTypeRebot || pUser.bOffline {
			pUser.iPlayingState = iPlayState
			pGame.uiPlayerReady |= 1 << uint(iSeat)
		}
	}

	//mlog.Debug("setUserReadyPlayingState(%s), uiPlayerReady[%d]", g_map_strPlayingState_debug[iPlayState], pGame.uiPlayerReady)
	//
	//for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
	//	pUser := pGame.ary_pUserData[iSeat]
	//	mlog.Debug("User[%d], iPlayingState[%s]", pUser.i32UserID, g_map_strPlayingState_debug[pUser.iPlayingState])
	//}
}

// 檢查是否所有玩家都已發送 raedy
func (pGame *Logic) checkAllUserReady(iPlayState int) {
	//mlog.Debug("checkAllUserReady(%s), uiPlayerReady[%d], uiAllPlayerReady[%d]", g_map_strPlayingState_debug[iPlayState],
	//	pGame.uiPlayerReady, pGame.uiAllPlayerReady)

	if pGame.uiPlayerReady != pGame.uiAllPlayerReady {
		return
	}

	switch iPlayState {
	case UserPlayState_StopBettingOK: // 停止下注動畫播放完成
		mlog.Debug("所有玩家停止下注動畫播放完成")
		if pGame.pGamePoker.CardCount() <= g_iShuffleRestCardCount {
			pGame.setGameState(GameState_Shuffle, pGame.mapI64Timeout[Timeout_Shuffle])
			pGame.shuffle(false, nil)
		} else {
			pGame.setGameState(GameState_Deal_Normal, pGame.mapI64Timeout[Timeout_Deal_Normal])
			pGame.dealNormal()
		}
	case UserPlayState_ShuffleOK: // 洗牌動畫播放完成
		mlog.Debug("所有玩家洗牌動畫播放完成")
		pGame.setGameState(GameState_Deal_Normal, pGame.mapI64Timeout[Timeout_Deal_Normal])
		pGame.dealNormal()
	case UserPlayState_DealOK_Normal: // 一般發牌動畫播放完成
		mlog.Debug("所有玩家一般發牌動畫播放完成")
		pGame.setGameState(GameState_Settlement_Normal, pGame.mapI64Timeout[Timeout_Settlement_Normal])
		pGame.settlementNormal()
	case UserPlayState_Settlement_Normal_OK: // 一般結算動畫播放完成
		mlog.Debug("所有玩家一般結算動畫播放完成")
		pGame.checkIfGameOver()
		if pGame.bGameOver {
			// 遊戲結束
			pGame.sendGameOver(false)
		} else {
			// 進入戰爭階段
			pGame.setGameState(GameState_WarSelection, pGame.mapI64Timeout[Timeout_WarSelection])
			pGame.enterWarPhase(false, nil)
		}
	case UserPlayState_WarSelectionOK: // 開始戰爭階段動畫播放完成
		mlog.Debug("所有玩家進入戰爭階段動畫播放完成")
		pGame.setGameState(GameState_Deal_War, pGame.mapI64Timeout[Timeout_Deal_War])
		pGame.dealWar()
	case UserPlayState_DealOK_War: // 戰爭發牌動畫播放完成
		mlog.Debug("所有玩家戰爭發牌動畫播放完成")
		pGame.setGameState(GameState_Settlement_War, pGame.mapI64Timeout[Timeout_Settlement_War])
		pGame.settlementWar()
	case UserPlayState_Settlement_War_OK: // 戰爭結算動畫播放完成
		mlog.Debug("所有玩家戰爭結算動畫播放完成")
		pGame.sendGameOver(false)
	}
}

// 洗牌
func (pGame *Logic) shuffle(bUserRecover bool, pUser *gamelogic.GameUser) {
	mlog.Debug("洗牌")

	// 建立新的牌組
	if !bUserRecover {
		pGame.pGamePoker.NewSet(PokerSetCount, 0)
		pGame.pGamePoker.Shuffle()
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Shuffle)

	if !bUserRecover {
		pGame.BoardcastMessage(pMsg)
	} else {
		if pUser != nil {
			pGame.SendMessageToUser(pUser, pMsg)
		}
	}

	pGame.setUserReadyPlayingState(UserPlayState_ShuffleOK)
	pGame.checkAllUserReady(UserPlayState_ShuffleOK)
}

// 訊息處理 - 洗牌動畫播放完成
func (pGame *Logic) messageShuffleOK(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_Shuffle {
		return
	}

	if pUser.iPlayingState != UserPlayState_StopBettingOK {
		return
	}

	mlog.Debug("User[%d] 洗牌播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_ShuffleOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_ShuffleOK)
}

// 一般發牌
func (pGame *Logic) dealNormal() {
	bChangeCard := false
	if g_bEnablePoolControl && g_pSysPoolWinLose.GameDBData.ChangeCardRate > 0 {
		iRnd := util.RndInt(0, 100)
		if iRnd < g_pSysPoolWinLose.GameDBData.ChangeCardRate {
			bChangeCard = true
		}
	}

	str := ""
	if g_bEnablePoolControl {
		str = fmt.Sprintf("一般發牌，換牌機率 %d%%", g_pSysPoolWinLose.GameDBData.ChangeCardRate)
	} else {
		str = "一般發牌"
	}
	mlog.Debug("%s", str)
	pGame.GameLog.AddGameStep("%s", str)

	if bChangeCard {
		// 換牌
		if g_bEnablePoolControl {
			mlog.Debug("換牌")
			pGame.GameLog.AddGameStep("換牌")
		}

		var arypChangingData []*ChangingCardsData
		pGame.pGamePoker, arypChangingData = changingCards(pGame.pGamePoker, pGame.aryPUserData, GamePhase_Normal)
		for i := 0; i < len(arypChangingData); i++ {
			if arypChangingData[i] == nil {
				continue
			}

			if arypChangingData[i].pCardSet != nil {
				if arypChangingData[i].i32SeatNo >= 0 && arypChangingData[i].i32SeatNo < int32(pGame.iSeatCount) {
					// 玩家
					pGame.aryPUserData[arypChangingData[i].i32SeatNo].pHandCard.AddCardSet(arypChangingData[i].pCardSet)
					pGame.aryPUserData[arypChangingData[i].i32SeatNo].i32Point = arypChangingData[i].i32Point
				} else if arypChangingData[i].i32SeatNo == -1 {
					// 莊家
					pGame.pDealerCard.AddCardSet(arypChangingData[i].pCardSet)
					pGame.i32DealerCardPoint = arypChangingData[i].i32Point
				}
			}
		}

		mlog.Debug("一般發牌")
		mlog.Debug("莊家, 權重[%d], 牌組[%s], 點數[%d]", g_pSysPoolWinLose.GameDBData.DealerWeight,
			pGame.pDealerCard.CodeString(), pGame.i32DealerCardPoint)
		for i := 0; i < len(pGame.aryPUserData); i++ {
			pUser := pGame.aryPUserData[i]
			if pUser != nil {
				mlog.Debug("玩家[%d(%d)], 權重[%d], 牌組[%s], 點數[%d]", pUser.i32UserID, pUser.i32SeatNo,
					pUser.iWeight, pUser.pHandCard.CodeString(), pUser.i32Point)
			}
		}
	} else {
		// 自然發牌
		if g_bEnablePoolControl {
			mlog.Debug("自然發牌")
			pGame.GameLog.AddGameStep("自然發牌")
		}

		// 發玩家牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			if pGame.aryPUserData[iSeat].iUserState == UserState_Normal {
				pSet, i32Point := pGame.dealing(GamePhase_Normal, iSeat)
				pGame.aryPUserData[iSeat].pHandCard.AddCardSet(pSet)
				pGame.aryPUserData[iSeat].i32Point = i32Point
			}
		}

		// 發莊家牌
		pSet, i32Point := pGame.dealing(GamePhase_Normal, -1)
		pGame.pDealerCard.AddCardSet(pSet)
		pGame.i32DealerCardPoint = i32Point
	}

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Deal_Normal)

	pData := new(netproto.DCZZ_Msg_Deal_Normal)
	pData.DealerCards = pGame.cardSetPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.aryPUserData[iSeat].pHandCard.CardCount() > 0 {
			mlog.Debug("玩家[%d] 手牌：%s，點數 %d", pGame.aryPUserData[iSeat].i32UserID,
				pGame.aryPUserData[iSeat].pHandCard.String(), pGame.aryPUserData[iSeat].i32Point)
			pGame.GameLog.AddGameStep("玩家[%d] 手牌：%s，點數 %d", pGame.aryPUserData[iSeat].i32UserID,
				pGame.aryPUserData[iSeat].pHandCard.String(), pGame.aryPUserData[iSeat].i32Point)
		}
	}
	mlog.Debug("莊家牌：%s，點數 %d", pGame.pDealerCard.String(), pGame.i32DealerCardPoint)
	pGame.GameLog.AddGameStep("莊家牌：%s，點數 %d", pGame.pDealerCard.String(), pGame.i32DealerCardPoint)

	pGame.setUserReadyPlayingState(UserPlayState_DealOK_Normal)
	pGame.checkAllUserReady(UserPlayState_DealOK_Normal)
}

// 發牌
// 傳入［遊戲階段］、［座位］
// 回傳［牌組］、［點數］
func (pGame *Logic) dealing(iGamePhase, iSeatNo int) (*pokerset.PokerCardSet, int32) {
	pRetSet := pokerset.NewPokerCardSet()
	mlog.Debug("dealing(), iGamePhase[%d]", iGamePhase)

	if iGamePhase < GamePhase_Normal || iGamePhase >= GamePhase_Max {
		return pRetSet, 0
	}

	bDealer := true
	if iSeatNo >= 0 && iSeatNo < pGame.iSeatCount {
		bDealer = false
		if pGame.aryPUserData[iSeatNo].bOperationDone {
			return pRetSet, 0
		}
	}
	var pCardBox *pokerset.PokerCardSet

	i32Point := int32(0)

	switch iGamePhase {
	case GamePhase_Normal: // 一般階段
		iDealCount := Count_CardCount_Normal
		switch pGame.iDealType {
		case DealType_Assigned: // 指定發牌
			if bDealer {
				pCardBox = pGame.pTestDealerCardBox
			} else {
				pCardBox = pGame.aryPTestUserCardBox[iSeatNo]
			}

			if pCardBox.CardCount() >= iDealCount {
				pSet := pCardBox.Deal(iDealCount)
				pRetSet.SetCardSet(pSet.CodeString())
			} else {
				iCount := iDealCount - pCardBox.CardCount()
				pRetSet.SetCardSet(pCardBox.CodeString())
				pSet := pGame.pGamePoker.Deal(iCount)
				pRetSet.AddCardSet(pSet)
			}
		default:
			pSet := pGame.pGamePoker.Deal(iDealCount)
			pRetSet.SetCardSet(pSet.CodeString())
		}
	case GamePhase_War: // 戰爭階段
		switch pGame.iDealType {
		case DealType_Assigned: // 指定發牌
			iDealCount := Count_CardCount_War - Count_HideCardCount - Count_CardCount_Normal
			if bDealer {
				pCardBox = pGame.pTestDealerCardBox
			} else {
				pCardBox = pGame.aryPTestUserCardBox[iSeatNo]
			}

			pHide := pokerset.NewPokerCardSet()
			for i := 0; i < Count_HideCardCount; i++ {
				pCard := pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
				pHide.AddSingleCard(pCard.CodeString())
			}
			pRetSet.AddCardSet(pHide)

			if pCardBox.CardCount() >= iDealCount {
				pSet := pCardBox.Deal(iDealCount)
				pRetSet.AddCardSet(pSet)
			} else {
				iCount := iDealCount - pCardBox.CardCount()
				pSet := pGame.pGamePoker.Deal(iCount)
				pRetSet.AddCardSet(pSet)
			}
		default:
			iDealCount := Count_CardCount_War - Count_CardCount_Normal
			pSet := pGame.pGamePoker.Deal(iDealCount)

			for i := 0; i < pSet.CardCount(); i++ {
				var pCard *pokerset.PokerCard
				if i < Count_HideCardCount {
					pCard = pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
				} else {
					pCard = pSet.GetSingleCard(i, true)
				}
				pRetSet.AddSingleCard(pCard.CodeString())
			}
		}
	}

	pLastCard := pRetSet.GetSingleCard(pRetSet.CardCount()-1, true)
	i32Point = pLastCard.I32Number

	if bDealer {
		mlog.Debug("dealing(%s), 莊家牌：%s, 點數[%d]", g_map_strGamePhase_debug[iGamePhase], pRetSet.String(), i32Point)
	} else {
		mlog.Debug("dealing(%s), 玩家[%d], 牌組：%s, 點數[%d]", g_map_strGamePhase_debug[iGamePhase],
			pGame.aryPUserData[iSeatNo].i32UserID, pRetSet.String(), i32Point)
	}

	return pRetSet, i32Point
}

// 一般發牌動畫播放完成
func (pGame *Logic) messageDealOKNormal(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_Deal_Normal {
		return
	}

	if pUser.iPlayingState != UserPlayState_StopBettingOK && pUser.iPlayingState != UserPlayState_ShuffleOK {
		return
	}

	mlog.Debug("User[%d] 一般發牌播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_DealOK_Normal
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_DealOK_Normal)
}

// 一般結算
func (pGame *Logic) settlementNormal() {
	mlog.Debug("一般階段結算")
	pGame.GameLog.AddGameStep("一般階段結算")

	strWinLose := ""
	// 一般階段結算
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState == UserState_Normal && !pUser.bOperationDone {
			i64WinLose, i64MoneyChange, i64Tax, i64Flow := int64(0), int64(0), int64(0), int64(0)
			if pUser.i32Point > pGame.i32DealerCardPoint { // 玩家贏
				i64MoneyChange, i64Tax = pGame.GetUserWinAndTax(pUser.i64Bet_Base * Rate_Normal)
				i64WinLose = pUser.i64Bet_Base + i64MoneyChange
				i64MoneyChange -= pUser.i64Bet_Side

				pUser.iSettlementWinLose = Settlement_PlayerWin
				pUser.bOperationDone = true
				pUser.i64Money += i64WinLose

				strWinLose = "贏"

				if g_bEnablePoolControl {
					if pUser.iUserType != gamelogic.UserTypeRebot {
						g_pSysPoolWinLose.i64PoolAmount -= pUser.i64Bet_Base * Rate_Normal
						g_pSysPoolWinLose.GameDBData.PlayerWin++
						pGame.checkSysPool()
					}
				}
			} else if pUser.i32Point < pGame.i32DealerCardPoint { // 莊家贏

				i64MoneyChange = -(pUser.i64Bet_Base + pUser.i64Bet_Side)
				i64WinLose = i64MoneyChange

				pUser.iSettlementWinLose = Settlement_PlayerLose
				pUser.bOperationDone = true

				strWinLose = "輸"

				if g_bEnablePoolControl {
					if pUser.iUserType != gamelogic.UserTypeRebot {
						g_pSysPoolWinLose.i64PoolAmount += pUser.i64Bet_Base + pUser.i64Bet_Side
						g_pSysPoolWinLose.GameDBData.DealerWin++
						pGame.checkSysPool()
					}
				}

				i64Flow = pUser.i64Bet_Base + pUser.i64Bet_Side
			} else { // 平手

				pUser.i64TieAward = pUser.i64Bet_Side * Rate_TieAward
				i64MoneyChange, i64Tax = pGame.GetUserWinAndTax(pUser.i64TieAward)
				i64WinLose = pUser.i64Bet_Side + i64MoneyChange

				pUser.iSettlementWinLose = Settlement_Tie
				pUser.bIsTied = true
				if pUser.i64TieAward > 0 {
					pUser.bGetTieAward = true
				} else {
					pUser.bGetTieAward = false
				}
				pUser.bCanEnterWar = true
				pUser.i64Money += i64WinLose
			}

			pUser.i64TotalWinLose = i64WinLose
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			if pUser.bIsTied {
				// 玩家與莊家平手
				mlog.Debug("玩家[%d] 與莊家平手，平手注 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID,
					Global.MoneyString(pUser.i64Bet_Side), Global.MoneyString(i64WinLose),
					Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] 與莊家平手，平手注 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID,
					Global.MoneyString(pUser.i64Bet_Side), Global.MoneyString(i64WinLose),
					Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
			} else {
				// 玩家輸贏已定，寫分
				mlog.Debug("User[%d] 寫分, i64MoneyChange[%s], i64Tax[%s], Money[%s]", pUser.i32UserID,
					Global.MoneyString(pUser.i64MoneyChange), Global.MoneyString(pUser.i64Tax), Global.MoneyString(pUser.i64Money))
				pGame.UpdateUserScoreFunc(pUser.User.UserID, pUser.i64MoneyChange, pUser.i64Tax)

				mlog.Debug("玩家[%d] %s 底注 %s，平手注 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID, strWinLose,
					Global.MoneyString(pUser.i64Bet_Base), Global.MoneyString(pUser.i64Bet_Side), Global.MoneyString(i64WinLose),
					Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
				pGame.GameLog.AddGameStep("玩家[%d] %s 底注 %s，平手注 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID, strWinLose,
					Global.MoneyString(pUser.i64Bet_Base), Global.MoneyString(pUser.i64Bet_Side), Global.MoneyString(i64WinLose),
					Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
			}
		}
	}
	// 記錄下來總投注金額, 以便結束時更新有效投注
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState == UserState_Normal {
			pUser.totalBet = pUser.i64Bet_Base + pUser.i64Bet_Side + pUser.i64Bet_Double
		}
	}

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Settlement_Normal)

	pData := new(netproto.DCZZ_Msg_Settlement_Normal)
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	pGame.setUserReadyPlayingState(UserPlayState_Settlement_Normal_OK)
	pGame.checkAllUserReady(UserPlayState_Settlement_Normal_OK)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState == UserState_Normal {
			pUser.i64Bet_Side = 0
			if !pUser.bIsTied {
				pUser.i64Bet_Base = 0
			}
		}
	}
}

// 一般結算動畫播放完成
func (pGame *Logic) messageSettlementNormalOK(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_Settlement_Normal {
		return
	}

	if pUser.iPlayingState != UserPlayState_DealOK_Normal {
		return
	}

	mlog.Debug("User[%d] 一般結算播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_Settlement_Normal_OK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_Settlement_Normal_OK)
}

// 檢查遊戲是否已結束
func (pGame *Logic) checkIfGameOver() {
	bGameOver := true

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.aryPUserData[iSeat].iUserState == UserState_Normal {
			if !pGame.aryPUserData[iSeat].bOperationDone {
				bGameOver = false
				break
			}
		}
	}

	pGame.bGameOver = bGameOver
	mlog.Debug("checkIfGameOver[%t]", pGame.bGameOver)
}

// 本局結束
func (pGame *Logic) sendGameOver(bUserRecover bool) {
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_GameOver)

	pGame.BoardcastMessage(pMsg)

	if !bUserRecover {
		pGame.setGameState(GameState_GameOver, pGame.mapI64Timeout[Timeout_GameOver])

		mlog.Debug("遊戲結束")
		pGame.GameLog.AddGameStep("遊戲結束")

		mapi64ScoreRet := make(map[int32]int64)
		betRet := make(map[int32]int64)
		winRet := make(map[int32]int64)
		validRet := make(map[int32]int64)
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.aryPUserData[iSeat]
			if pUser.User != nil {
				var flowAmount, betAmount, payAmount int64
				betAmount = pUser.totalBet
				payAmount = betAmount + pUser.i64MoneyChange
				// 總流水為輸錢絕對值
				if pUser.i64MoneyChange < 0 {
					flowAmount = -pUser.i64MoneyChange
				} else {
					flowAmount = 0
				}
				pGame.GameLog.AddGameStep("用戶 %d : 有效投注 %s, 返還金額 %s, 本局流水 %s, 本局稅金: %s, 更新金額: %s", pUser.User.UserID,
					Global.MoneyString(betAmount), Global.MoneyString(payAmount), Global.MoneyString(flowAmount), Global.MoneyString(pUser.i64Tax), Global.MoneyString(pUser.i64MoneyChange))
				pUser.User.FlowAmount += flowAmount
				pUser.User.BetAmount += betAmount
				pUser.User.PayAmount += payAmount
				mapi64ScoreRet[pUser.i32UserID] = pUser.i64MoneyChange
				betRet[pUser.i32UserID] = betAmount
				winRet[pUser.i32UserID] = payAmount + pUser.i64Tax
				validRet[pUser.i32UserID] = flowAmount

				pGame.GameLog.Context().UserUpdateScore(pUser.i32UserID, pUser.i64MoneyChange)
				pGame.GameLog.Context().UserBetScore(pUser.i32UserID, betAmount)
				pGame.GameLog.Context().UserPayoutScore(pUser.i32UserID, winRet[pUser.i32UserID])
				pGame.GameLog.Context().UserValidScore(pUser.i32UserID, flowAmount)
				pGame.GameLog.Context().UserTaxAmount(pUser.i32UserID, pUser.i64Tax)
			}
		}

		if g_bEnablePoolControl {
			pGame.UpdateSysPoolWinLose()
		}

		pGame.GameLog.GameEnd()

		pGame.GameEndFunc()

		mlog.Debug("TableNo[%d], sendGameOver(), IsGaming[%t]", pGame.TableNo(), pGame.IsGaming())

		// 清除資料
		pGame.pDealerCard.SetCardSet("")
		pGame.i32DealerCardPoint = 0

		iUserCount := 0
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pGameUser := pGame.LogicBase.GetUserBySeatNo(int32(iSeat))
			if pGameUser == nil {
				pGame.aryPUserData[iSeat].Init()
			} else {
				if pGameUser.IsOffLine() {
					pGame.UserLeave(pGameUser)
					pGame.aryPUserData[iSeat].Init()
				} else {
					i32LoginMoney := pGame.GetTable().Room.LoginMoney
					if pGame.aryPUserData[iSeat].i64Money < int64(i32LoginMoney) {
						pGame.LogicBase.KickUserFunc(pGameUser.UserID, "您的金币不足！请补充金币后再进行游戏。")
						pGame.aryPUserData[iSeat].Init()
					} else {
						pGame.aryPUserData[iSeat].ResetGameData()
						pGame.aryPUserData[iSeat].iUserState = UserState_Normal
						iUserCount++
					}
				}
			}
		}

		pGame.iUserCount = iUserCount
		mlog.Debug("gameOver, 玩家人數：%d", pGame.iUserCount)
	}
}

// 進入戰爭階段
func (pGame *Logic) enterWarPhase(bUserRecover bool, pUser *gamelogic.GameUser) {
	mlog.Debug("進入戰爭階段")

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase)

	pData := new(netproto.DCZZ_Msg_WarPhase)
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}
	pData.TimeRemain = proto.Int32(int32((pGame.i64Timeout - pGame.i64TimeDuration + pGame.i64LastMessageTime) / 10)) // 剩餘時間

	pMsg.MsgData = pData
	if bUserRecover {
		if pUser != nil {
			pGame.SendMessageToUser(pUser, pMsg)
		}
	} else {

		pGame.uiPlayerReady = 0

		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pGame.aryPUserData[iSeat].iPlayingState = UserPlayState_WarSelection
		}

		pGame.BoardcastMessage(pMsg)
	}
}

// 戰爭階段選項
func (pGame *Logic) messageWarSelection(pUser *User, pMsg *msg.Message) {
	if pUser == nil || pMsg == nil {
		return
	}

	if !pUser.bCanEnterWar || pUser.bIsFolded {
		return
	}

	if pGame.iGameState != GameState_WarSelection {
		return
	}

	if pUser.iPlayingState != UserPlayState_WarSelection {
		return
	}

	msgInfo := &netproto.DCZZ_Msg_WarPhaseSelection{}
	proto.Unmarshal(pMsg.MsgData.([]byte), msgInfo)
	iSelection := int(msgInfo.GetSelection())

	pGame.userActionWarSelection(pUser, iSelection)
}

// 玩家動作 - 戰爭階段選項
func (pGame *Logic) userActionWarSelection(pUser *User, iSelection int) {
	if pUser == nil {
		return
	}

	if iSelection < WarSelection_Fold || iSelection >= WarSelection_Max {
		return
	}

	switch iSelection {
	case WarSelection_Fold: // 棄牌
		i64WinLose, i64MoneyChange, i64Flow := int64(0), int64(0), int64(0)
		i64WinLose = pUser.i64Bet_Base / 2
		i64MoneyChange = -i64WinLose
		i64Flow = i64WinLose

		pUser.bIsFolded = true
		pUser.bOperationDone = true

		pUser.i64TotalWinLose = i64WinLose
		pUser.i64MoneyChange += i64MoneyChange
		pUser.i64Money += pUser.i64TotalWinLose

		if g_bEnablePoolControl {
			if pUser.iUserType != gamelogic.UserTypeRebot {
				g_pSysPoolWinLose.i64PoolAmount += pUser.i64Bet_Base / 2
				g_pSysPoolWinLose.GameDBData.DealerWin++
				pGame.checkSysPool()
			}
		}

		mlog.Debug("User[%d] 寫分, i64MoneyChange[%s], i64Tax[%s], Money[%s]", pUser.i32UserID,
			Global.MoneyString(pUser.i64MoneyChange), Global.MoneyString(pUser.i64Tax), Global.MoneyString(pUser.i64Money))
		pGame.UpdateUserScoreFunc(pUser.User.UserID, pUser.i64MoneyChange, pUser.i64Tax)

		mlog.Debug("玩家[%d] 棄牌，收回 %s，流水 %s，剩餘金錢 %s", pUser.i32UserID, Global.MoneyString(i64WinLose),
			Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 棄牌，收回 %s，流水 %s，剩餘金錢 %s", pUser.i32UserID, Global.MoneyString(i64WinLose),
			Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
	case WarSelection_War: // 戰爭
		pUser.bGoForWar = true
		pUser.i64Bet_Double = pUser.i64Bet_Base * Times_BetDouble
		pUser.i64Money -= pUser.i64Bet_Double

		mlog.Debug("玩家[%d] 加入戰爭，底注 %s，加倍 %s，剩餘金錢 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Bet_Base),
			Global.MoneyString(pUser.i64Bet_Double), Global.MoneyString(pUser.i64Money))
		pGame.GameLog.AddGameStep("玩家[%d] 加入戰爭，底注 %s，加倍 %s，剩餘金錢 %s", pUser.i32UserID, Global.MoneyString(pUser.i64Bet_Base),
			Global.MoneyString(pUser.i64Bet_Double), Global.MoneyString(pUser.i64Money))
	}

	pGame.sendWarSelectionResult(pUser, iSelection)

	pGame.checkIfGameOver()
	if pGame.bGameOver {
		// 遊戲結束
		//pGame.sendGameOver(false)
		pGame.bAllUserDoneAtWarSelection = true
		pGame.setGameState(GameState_GameOver, pGame.mapI64Timeout[Timeout_GameOver])
	} else {
		// 玩家是否都選擇完成
		bAllSelected := true
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.aryPUserData[iSeat]
			if pUser.iUserState == UserState_Normal && !pUser.bOperationDone {
				if !pUser.bGoForWar && !pUser.bIsFolded {
					bAllSelected = false
				}
			}
		}

		if bAllSelected {
			// 開始戰爭階段
			pGame.setGameState(GameState_WarPhaseAnimation, pGame.mapI64Timeout[Timeout_WarPhaseBegin])
			pGame.warPhaseBegin()
		}
	}
}

// 戰爭階段選擇結果
func (pGame *Logic) sendWarSelectionResult(pUser *User, iSelection int) {
	if pUser == nil {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseselectionResult)

	pData := new(netproto.DCZZ_Msg_WarPhaseselectionResult)
	pData.SeatNo = proto.Int32(pUser.i32SeatNo + int32(IndexBase)) // 做選擇的玩家
	pData.Selection = proto.Int32(int32(iSelection))               // 玩家選擇
	pData.MoneyRemain = proto.Int64(pUser.i64Money)                // 下注後玩家剩餘金錢

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)
}

// 開始戰爭階段
func (pGame *Logic) warPhaseBegin() {
	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseBegin)

	pGame.BoardcastMessage(pMsg)

	pGame.setUserReadyPlayingState(UserPlayState_WarSelectionOK)
	pGame.checkAllUserReady(UserPlayState_WarSelectionOK)
}

// 進入戰爭階段動畫播放完成
func (pGame *Logic) messageWarPhaseAnimationOK(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_WarPhaseAnimation {
		return
	}

	if pUser.iPlayingState != UserPlayState_WarSelection {
		return
	}

	mlog.Debug("User[%d] 進入戰爭階段動畫播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_WarSelectionOK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_WarSelectionOK)
}

// 戰爭發牌
func (pGame *Logic) dealWar() {
	bChangeCard := false
	if g_bEnablePoolControl && g_pSysPoolWinLose.GameDBData.ChangeCardRate > 0 {
		iRnd := util.RndInt(0, 100)
		if iRnd < g_pSysPoolWinLose.GameDBData.ChangeCardRate {
			bChangeCard = true
		}
	}

	str := ""
	if g_bEnablePoolControl {
		str = fmt.Sprintf("戰爭發牌，換牌機率 %d%%", g_pSysPoolWinLose.GameDBData.ChangeCardRate)
	} else {
		str = "戰爭發牌"
	}
	mlog.Debug("%s", str)
	pGame.GameLog.AddGameStep("%s", str)

	if bChangeCard {
		// 換牌
		if g_bEnablePoolControl {
			mlog.Debug("換牌")
			pGame.GameLog.AddGameStep("換牌")
		}

		var aryPChangingData []*ChangingCardsData
		pGame.pGamePoker, aryPChangingData = changingCards(pGame.pGamePoker, pGame.aryPUserData, GamePhase_War)
		for i := 0; i < len(aryPChangingData); i++ {
			if aryPChangingData[i] == nil {
				continue
			}

			if aryPChangingData[i].pCardSet != nil {
				if aryPChangingData[i].i32SeatNo >= 0 && aryPChangingData[i].i32SeatNo < int32(pGame.iSeatCount) {
					// 玩家
					pGame.aryPUserData[aryPChangingData[i].i32SeatNo].pHandCard.AddCardSet(aryPChangingData[i].pCardSet)
					pGame.aryPUserData[aryPChangingData[i].i32SeatNo].i32Point = aryPChangingData[i].i32Point
				} else if aryPChangingData[i].i32SeatNo == -1 {
					// 莊家
					pGame.pDealerCard.AddCardSet(aryPChangingData[i].pCardSet)
					pGame.i32DealerCardPoint = aryPChangingData[i].i32Point
				}
			}
		}

		mlog.Debug("戰爭發牌")
		mlog.Debug("莊家, 權重[%d], 牌組[%s], 點數[%d]", g_pSysPoolWinLose.GameDBData.DealerWeight,
			pGame.pDealerCard.CodeString(), pGame.i32DealerCardPoint)
		for i := 0; i < len(pGame.aryPUserData); i++ {
			pUser := pGame.aryPUserData[i]
			if pUser != nil {
				mlog.Debug("玩家[%d(%d)], 權重[%d], 牌組[%s], 點數[%d]", pUser.i32UserID, pUser.i32SeatNo,
					pUser.iWeight, pUser.pHandCard.CodeString(), pUser.i32Point)
			}
		}
	} else {
		// 自然發牌
		if g_bEnablePoolControl {
			mlog.Debug("自然發牌")
			pGame.GameLog.AddGameStep("自然發牌")
		}

		// 發玩家牌
		for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
			pUser := pGame.aryPUserData[iSeat]
			if pUser.iUserState == UserState_Normal {
				pSet, i32Point := pGame.dealing(GamePhase_War, iSeat)
				pUser.pHandCard.AddCardSet(pSet)
				pUser.i32Point = i32Point
			}
		}

		// 發莊家牌
		pSet, i32Point := pGame.dealing(GamePhase_War, -1)
		pGame.pDealerCard.AddCardSet(pSet)
		pGame.i32DealerCardPoint = i32Point
	}

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Deal_War)

	pData := new(netproto.DCZZ_Msg_Deal_War)
	pData.DealerCards = pGame.cardSetPkg(pGame.pDealerCard) // 莊家手牌
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if pGame.aryPUserData[iSeat].pHandCard.CardCount() > 0 {
			mlog.Debug("玩家[%d] 手牌：%s，點數 %d", pGame.aryPUserData[iSeat].i32UserID,
				pGame.aryPUserData[iSeat].pHandCard.String(), pGame.aryPUserData[iSeat].i32Point)
			pGame.GameLog.AddGameStep("玩家[%d] 手牌：%s，點數 %d", pGame.aryPUserData[iSeat].i32UserID,
				pGame.aryPUserData[iSeat].pHandCard.String(), pGame.aryPUserData[iSeat].i32Point)
		}
	}
	mlog.Debug("莊家牌：%s，點數 %d", pGame.pDealerCard.String(), pGame.i32DealerCardPoint)
	pGame.GameLog.AddGameStep("莊家牌：%s，點數 %d", pGame.pDealerCard.String(), pGame.i32DealerCardPoint)

	pGame.setUserReadyPlayingState(UserPlayState_DealOK_War)
	pGame.checkAllUserReady(UserPlayState_DealOK_War)
}

// 戰爭發牌動畫播放完成
func (pGame *Logic) messageDealOKWar(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_Deal_War {
		return
	}

	if pUser.iPlayingState != UserPlayState_WarSelectionOK {
		return
	}

	mlog.Debug("User[%d] 戰爭發牌播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_DealOK_War
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_DealOK_War)
}

// 戰爭結算
func (pGame *Logic) settlementWar() {
	mlog.Debug("戰爭階段結算")
	pGame.GameLog.AddGameStep("戰爭階段結算")

	strWinLose := ""
	// 戰爭階段結算
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState == UserState_Normal && !pUser.bOperationDone {
			var i64WinLose, i64MoneyChange, i64Tax, i64Flow int64
			if pUser.i32Point > pGame.i32DealerCardPoint {
				// 玩家贏
				i64MoneyChange, i64Tax = pGame.GetUserWinAndTax(pUser.i64Bet_Double * Rate_War)
				i64WinLose = pUser.i64Bet_Base + pUser.i64Bet_Double + i64MoneyChange

				pUser.iSettlementWinLose = Settlement_PlayerWin
				pUser.i64Money += i64WinLose

				strWinLose = "贏"

				if g_bEnablePoolControl {
					if pUser.iUserType != gamelogic.UserTypeRebot {
						g_pSysPoolWinLose.i64PoolAmount -= pUser.i64Bet_Double * Rate_War
						g_pSysPoolWinLose.GameDBData.PlayerWin++
						pGame.checkSysPool()
					}
				}
			} else if pUser.i32Point < pGame.i32DealerCardPoint {
				// 莊家贏
				i64MoneyChange = -(pUser.i64Bet_Base + pUser.i64Bet_Double)
				i64WinLose = i64MoneyChange
				i64Flow = int64(math.Abs(float64(i64MoneyChange)))

				pUser.iSettlementWinLose = Settlement_PlayerLose

				strWinLose = "輸"

				if g_bEnablePoolControl {
					if pUser.iUserType != gamelogic.UserTypeRebot {
						g_pSysPoolWinLose.i64PoolAmount += pUser.i64Bet_Base + pUser.i64Bet_Double
						g_pSysPoolWinLose.GameDBData.DealerWinWar++
						pGame.checkSysPool()
					}
				}
			} else {
				// 平手
				i64MoneyChange, i64Tax = pGame.GetUserWinAndTax((pUser.i64Bet_Double + pUser.i64Bet_Base) * Rate_War)
				i64WinLose = pUser.i64Bet_Double + pUser.i64Bet_Base + i64MoneyChange

				pUser.iSettlementWinLose = Settlement_PlayerWin
				pUser.i64Money += i64WinLose

				if g_bEnablePoolControl {
					if pUser.iUserType != gamelogic.UserTypeRebot {
						g_pSysPoolWinLose.i64PoolAmount -= (pUser.i64Bet_Double + pUser.i64Bet_Base) * Rate_War
						g_pSysPoolWinLose.GameDBData.PlayerWinWarTie++
						pGame.checkSysPool()
					}
				}
			}

			pUser.bOperationDone = true
			pUser.i64TotalWinLose = i64WinLose
			pUser.i64MoneyChange += i64MoneyChange
			pUser.i64Tax += i64Tax

			mlog.Debug("User[%d] 寫分, i64MoneyChange[%s], i64Tax[%s], Money[%s]", pUser.i32UserID,
				Global.MoneyString(pUser.i64MoneyChange), Global.MoneyString(pUser.i64Tax), Global.MoneyString(pUser.i64Money))
			pGame.UpdateUserScoreFunc(pUser.User.UserID, pUser.i64MoneyChange, pUser.i64Tax)

			mlog.Debug("玩家[%d] %s 底注 %s，加倍 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID, strWinLose,
				Global.MoneyString(pUser.i64Bet_Base), Global.MoneyString(pUser.i64Bet_Double), Global.MoneyString(i64WinLose),
				Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
			pGame.GameLog.AddGameStep("玩家[%d] %s 底注 %s，加倍 %s，輸贏 %s，稅收 %s，流水 %s，剩餘金額 %s", pUser.i32UserID, strWinLose,
				Global.MoneyString(pUser.i64Bet_Base), Global.MoneyString(pUser.i64Bet_Double), Global.MoneyString(i64WinLose),
				Global.MoneyString(i64Tax), Global.MoneyString(i64Flow), Global.MoneyString(pUser.i64Money))
		}
	}

	// 記錄下來總投注金額, 以便結束時更新有效投注
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pUser := pGame.aryPUserData[iSeat]
		if pUser.iUserState == UserState_Normal {
			pUser.totalBet += pUser.i64Bet_Double
		}
	}

	// 傳送訊息給 client
	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Settlement_War)

	pData := new(netproto.DCZZ_Msg_Settlement_War)
	// 所有玩家資訊
	for i := 0; i < pGame.iSeatCount; i++ {
		pUserInfo := pGame.userInfoPkg(i)
		pData.AllUserInfo = append(pData.AllUserInfo, pUserInfo)
	}

	pMsg.MsgData = pData
	pGame.BoardcastMessage(pMsg)

	pGame.setUserReadyPlayingState(UserPlayState_Settlement_War_OK)
	pGame.checkAllUserReady(UserPlayState_Settlement_War_OK)

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGame.aryPUserData[iSeat].i64Bet_Base = 0
		pGame.aryPUserData[iSeat].i64Bet_Double = 0
	}
}

// 戰爭結算動畫播放完成
func (pGame *Logic) messageSettlementWarOK(pUser *User) {
	if pUser == nil {
		return
	}

	if pGame.iGameState != GameState_Settlement_War {
		return
	}

	if pUser.iPlayingState != UserPlayState_DealOK_War {
		return
	}

	mlog.Debug("User[%d] 戰爭結算播放完成", pUser.i32UserID)

	pUser.iPlayingState = UserPlayState_Settlement_War_OK
	pGame.uiPlayerReady |= 1 << uint(pUser.i32SeatNo)
	pGame.checkAllUserReady(UserPlayState_Settlement_War_OK)
}

// 遊戲結束
func (pGame *Logic) gameOver() {
	pGame.initData()
	pGame.sendAllUserInfo()
}

// UserLeave 玩家離線
func (pGame *Logic) UserLeave(pGameUser *gamelogic.GameUser) bool {
	if pGameUser == nil {
		return false
	}

	mlog.Debug("Table[%d], UserLeave, GameState[%s], User[%d], Seat[%d], IsOffLine()[%t]",
		pGame.TableNo(), g_map_strGameState_debug[pGame.iGameState], pGameUser.UserID, pGameUser.SeatNo, pGameUser.IsOffLine())

	iSeatNo := int(pGameUser.SeatNo)
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return false
	}

	pUser := pGame.aryPUserData[iSeatNo]
	pUser.bOffline = true

	if pGame.iGameState >= GameState_Betting && pGame.iGameState <= GameState_Settlement_War && pUser.iUserState != UserState_Watch {
		return false
	}

	pUser.Init()
	pGame.iUserCount--
	pGame.sendAllUserInfo()

	defer func() {
		// check all user leave, close game ...
		bHasUser := false
		for _, pUser := range pGame.aryPUserData {
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

// UserRecover 玩家斷線重連
func (pGame *Logic) UserRecover(pGameUser *gamelogic.GameUser) {
	if pGameUser == nil {
		return
	}

	mlog.Debug("Table[%d], UserRecover, GameState[%s], User[%d], Seat[%d], IsOffLine()[%t]",
		pGame.TableNo(), g_map_strGameState_debug[pGame.iGameState], pGameUser.UserID, pGameUser.SeatNo, pGameUser.IsOffLine())

	pGame.DynamicGameLogicBase.UserRecover(pGameUser)

	iSeatNo := int(pGameUser.SeatNo)
	if iSeatNo < 0 || iSeatNo >= pGame.iSeatCount {
		return
	}

	pUser := pGame.aryPUserData[iSeatNo]
	pUser.bOffline = false

	mlog.Debug("GameState[%s], User[%d], BetBase[%d], BetSide[%d], BetDouble[%d], HandCard[%s]", g_map_strGameState_debug[pGame.iGameState],
		pUser.i32UserID, pUser.i64Bet_Base, pUser.i64Bet_Side, pUser.i64Bet_Double, pUser.pHandCard.String())

	pGame.sendGameInfo(pUser)
	pGame.sendAllUserInfo()

	switch pGame.iGameState {
	case GameState_Betting: // 押注
		pGame.startBetting(true, pGameUser)
	case GameState_Shuffle: // 洗牌
		pGame.shuffle(true, pGameUser)
	case GameState_WarSelection: // 選擇是否進入戰爭
		pGame.enterWarPhase(true, pGameUser)
	case GameState_GameOver: // 遊戲結束
		pGame.sendGameOver(true)
	}
}

// 設定遊戲狀態
func (pGame *Logic) setGameState(iGameState int, i64Timeout int64) {
	if iGameState < GameState_WaitingPlayer || iGameState > GameState_GameOver {
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

	mlog.Debug("TableNo[%d],【state】[%s], timeout[%d], duration[%d], lastmsg[%d], IsGaming[%t]", pGame.TableNo(), g_map_strGameState_debug[pGame.iGameState],
		pGame.i64Timeout, pGame.i64TimeDuration, pGame.i64LastMessageTime, pGame.IsGaming())
}

// 是否所有玩家都完成操作
func (pGame *Logic) allUserOperationDone() bool {
	bAllDone := true

	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		if !pGame.aryPUserData[iSeat].bOperationDone {
			bAllDone = false
			break
		}
	}

	return bAllDone
}

// 水線判定
func (pGame *Logic) checkSysPool() {
	g_pSysPoolWinLose.GameDBData.PersonTime++
	strCheckLog, strUpdateLog := checkSysPool()

	if strCheckLog != "" {
		mlog.Debug("%s", strCheckLog)
		pGame.GameLog.AddGameStep("%s", strCheckLog)
	}

	if strUpdateLog != "" {
		mlog.Debug("%s", strUpdateLog)
		pGame.GameLog.AddGameStep("%s", strUpdateLog)
	}
}
