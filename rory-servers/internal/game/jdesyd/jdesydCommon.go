package jdesydLogic

import (
	"fmt"

	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

const (
	JD21_IndexBase         int   = 1   // 索引值起始值
	JD21_PokerSetCount     int32 = 6   // 使用幾副牌
	JD21_MinCutCardIndex   int   = 53  // 切牌區間的最小值
	JD21_MaxCutCardIndex   int   = 260 // 切牌區間的最大值
	JD21_MaxSplitCount     int   = 4   // 最大分牌組數
	JD21_Soft              int32 = 17  // 最低點數
	JD21_InsuranceMultiple int   = 2   // 保險賠付倍數
	JD21_BasicCardCount    int   = 2   // 基本張數
)

// 遊戲狀態
const (
	JD21_State_WaitingPlayer        int = iota // 等待玩家加入
	JD21_State_WaitingStart                    // 等待遊戲開始
	JD21_State_Betting                         // 押底注
	JD21_State_StopBetting                     // 停止下注
	JD21_State_CutCards                        // 切牌
	JD21_State_CutCardsAnimation               // 切牌動畫
	JD21_State_Deal                            // 發牌
	JD21_State_Insurance                       // 保險
	JD21_State_DealerPeeksCards                // 莊家瞇牌
	JD21_State_Settlement_Insurance            // 保險結算
	JD21_State_Playing                         // 玩家操作
	JD21_State_PlayingDelay                    // 玩家操作結果傳送後的延遲時間
	JD21_State_Settlement_Dealer               // 莊家開牌/結算
	JD21_State_GameOver                        // 遊戲結束
)

// 送給 client 的遊戲狀態
const (
	JD21_CState_WaitingStart         int = iota // 等待遊戲開始
	JD21_CState_Betting                         // 押底注
	JD21_CState_CutCards                        // 切牌
	JD21_CState_Deal                            // 發牌
	JD21_CState_Insurance                       // 保險
	JD21_CState_Settlement_Insurance            // 保險結算
	JD21_CState_Playing                         // 玩家操作
	JD21_CState_Settlement                      // 莊家開牌/結算
)

// timeout 類別
const (
	JD21_Timeout_WaitingStart         int = iota // 等待遊戲開始
	JD21_Timeout_Betting                         // 押底注
	JD21_Timeout_StopBetting                     // 停止下注
	JD21_Timeout_CutCards                        // 切牌
	JD21_Timeout_CutCardsAnimation               // 切牌動畫
	JD21_Timeout_DealAnimation                   // 發牌動畫
	JD21_Timeout_Insurance                       // 保險
	JD21_Timeout_DealerPeeksCards                // 莊家瞇牌
	JD21_Timeout_Settlement_Insurance            // 保險結算
	JD21_Timeout_Playing                         // 玩家操作
	JD21_Timeout_Settlement_Dealer               // 莊家開牌/結算

	// 以下為系統自訂 timeout，非自 config 中讀取的值
	JD21_Timeout_PlayingDelay // 玩家操作結果傳送後的延遲時間
	JD21_Timeout_RobotCutCard // 若無可切牌玩家時，由機器人自動切牌的延遲時間
)

// 玩家狀態
const (
	JD21_UserState_EmptySeat int = iota // 空位
	JD21_UserState_Normal               // 正常
	JD21_UserState_Watch                // 觀戰
)

// 玩家遊戲狀態
const (
	JD21_UserPlayState_WaitingStart           int = iota // 等待遊戲開始
	JD21_UserPlayState_Betting                           // 開始下注
	JD21_UserPlayState_StopBettingOK                     // 停止下注動畫播放完成
	JD21_UserPlayState_CutCards                          // 切牌
	JD21_UserPlayState_CutCardsOK                        // 切牌動畫播放完成
	JD21_UserPlayState_DealOK                            // 發牌動畫播放完成
	JD21_UserPlayState_TakeInsurance                     // 玩家是否買保險
	JD21_UserPlayState_TakeInsuranceOK                   // 玩家已買保險
	JD21_UserPlayState_DealerOK                          // 莊家瞇牌動畫播放完成
	JD21_UserPlayState_Settlement_InsuranceOK            // 保險結算動畫播放完成
	JD21_UserPlayState_Playing_Waiting                   // 玩家操作 - 等待中
	JD21_UserPlayState_Playing_MyTurn                    // 玩家操作 - 輪到我了
	JD21_UserPlayState_Playing_Done                      // 玩家操作 - 完成操作
	JD21_UserPlayState_Settlement_DealerOK               // 莊家開牌/結算
)

// 押注結果
const (
	JD21_BetResult_Auto             = 1  // 自動下底注
	JD21_BetResult_Success          = 0  // 下注成功
	JD21_BetResult_NotEnoughMoney   = -1 // 玩家餘額不足
	JD21_BetResult_BettingStopped   = -2 // 已停止下注
	JD21_BetResult_OverUpperLimit   = -3 // 總押注額超過上限
	JD21_BetResult_BetSmallThanZero = -4 // 押注金額<=0
	JD21_BetResult_BetFailed        = -5 // 下注失敗
)

// 玩家買保險結果
const (
	JD21_Insurance_Success        = 0  // 成功
	JD21_Insurance_NotEnoughMoney = -1 // 餘額不足
)

// 玩家動作結果
const (
	JD21_Playing_Success              = 0  // 成功
	JD21_Playing_Failed               = -1 // 失敗
	JD21_Playing_NotEnoughMoney       = -2 // 餘額不足
	JD21_Playing_CardSetOperationDone = -3 // 此牌組已不可操作
	JD21_Playing_CannotSplit          = -4 // 不可分牌
	JD21_Playing_CannotDouble         = -5 // 不可加倍
	JD21_Playing_IsBusted             = -6 // 已經爆牌
	JD21_Playing_CannotSurrender      = -7 // 不可投降
)

// 玩家動作
const (
	JD21_UserAction_Stand     = 0 // 停牌
	JD21_UserAction_Split     = 1 // 分牌
	JD21_UserAction_Double    = 2 // 加倍
	JD21_UserAction_Hit       = 3 // 要牌
	JD21_UserAction_Surrender = 4 // 投降
)

// 比牌結果
const (
	JD21_Compare_Even   int = iota // 平局
	JD21_Compare_Dealer            // 莊家贏
	JD21_Compare_Player            // 玩家贏
)

// 發牌器發牌方式
const (
	JD21_DealType_Natural  int = iota // 自然發牌（預設值）
	JD21_DealType_Assigned            // 指定發牌
	JD21_DealType_Max
)

//---------------------------------------------------------------------------
var g_i64MaxTotalBetBase = int64(0)            // 押注總額上限
var g_iCutCardIndex = JD21_MinCutCardIndex - 1 // 切牌位置索引值
// 點數值
var g_map_i32CardNoValue = map[int32]int32{
	Pokerset.POKERCardNo1:  1,
	Pokerset.POKERCardNo2:  2,
	Pokerset.POKERCardNo3:  3,
	Pokerset.POKERCardNo4:  4,
	Pokerset.POKERCardNo5:  5,
	Pokerset.POKERCardNo6:  6,
	Pokerset.POKERCardNo7:  7,
	Pokerset.POKERCardNo8:  8,
	Pokerset.POKERCardNo9:  9,
	Pokerset.POKERCardNo10: 10,
	Pokerset.POKERCardNoJ:  10,
	Pokerset.POKERCardNoQ:  10,
	Pokerset.POKERCardNoK:  10,
	Pokerset.POKERCardNoA:  1,
}

//---------------------------------------------------------------------------
//  ----- 後台記錄用 -----
// 玩家狀態
var g_map_strUserState = map[int]string{
	JD21_UserState_EmptySeat: "空位", // 空位
	JD21_UserState_Normal:    "正常", // 正常
	JD21_UserState_Watch:     "观战", // 觀戰
}

// 下注結果
var g_map_strBetResult = map[int]string{
	JD21_BetResult_Auto:             "自動下底注",    // 自動下底注
	JD21_BetResult_Success:          "下注成功",     // 下注成功
	JD21_BetResult_NotEnoughMoney:   "玩家余额不足",   // 玩家餘額不足
	JD21_BetResult_BettingStopped:   "已停止下注",    // 已停止下注
	JD21_BetResult_OverUpperLimit:   "总压住额超过上限", // 總押注額超過上限
	JD21_BetResult_BetSmallThanZero: "押注金额小於 0", // 押注金額<=0
}

// 玩家動作
var g_map_strUserAction = map[int]string{
	JD21_UserAction_Stand:     "停牌", // 停牌
	JD21_UserAction_Split:     "分牌", // 分牌
	JD21_UserAction_Double:    "加倍", // 加倍
	JD21_UserAction_Hit:       "要牌", // 要牌
	JD21_UserAction_Surrender: "投降", // 投降
}

//---------------------------------------------------------------------------
// ----- for debug -----
// 遊戲狀態
var g_map_strGameState_debug = map[int]string{
	JD21_State_WaitingPlayer:        fmt.Sprintf("WaitingPlayer(%d)", JD21_State_WaitingPlayer),               // 等待玩家加入
	JD21_State_WaitingStart:         fmt.Sprintf("WaitingStart(%d)", JD21_State_WaitingStart),                 // 等待遊戲開始
	JD21_State_Betting:              fmt.Sprintf("Betting(%d)", JD21_State_Betting),                           // 押底注
	JD21_State_StopBetting:          fmt.Sprintf("StopBetting(%d)", JD21_State_StopBetting),                   // 停止下注
	JD21_State_CutCards:             fmt.Sprintf("CutCard(%d)", JD21_State_CutCards),                          // 切牌
	JD21_State_CutCardsAnimation:    fmt.Sprintf("CutCardsAnimation(%d)", JD21_State_CutCardsAnimation),       // 切牌動畫
	JD21_State_Deal:                 fmt.Sprintf("Deal(%d)", JD21_State_Deal),                                 // 發牌
	JD21_State_Insurance:            fmt.Sprintf("Insurance(%d)", JD21_State_Insurance),                       // 保險
	JD21_State_DealerPeeksCards:     fmt.Sprintf("DealerPeeksCards(%d)", JD21_State_DealerPeeksCards),         // 莊家瞇牌
	JD21_State_Settlement_Insurance: fmt.Sprintf("Settlement_Insurance(%d)", JD21_State_Settlement_Insurance), // 保險結算
	JD21_State_Playing:              fmt.Sprintf("Playing(%d)", JD21_State_Playing),                           // 玩家操作
	JD21_State_Settlement_Dealer:    fmt.Sprintf("Settlement(%d)", JD21_State_Settlement_Dealer),              // 莊家開牌/結算
	JD21_State_PlayingDelay:         fmt.Sprintf("PlayingDelay(%d)", JD21_State_PlayingDelay),                 // 玩家操作結果傳送後的延遲時間
	JD21_State_GameOver:             fmt.Sprintf("GameOver(%d)", JD21_State_GameOver),                         // 遊戲結束
}

// 送給 client 的遊戲狀態
var g_map_strToClientState_debug = map[int]string{
	JD21_CState_WaitingStart:         fmt.Sprintf("WaitingStart(%d)", JD21_CState_WaitingStart),                 // 等待遊戲開始
	JD21_CState_Betting:              fmt.Sprintf("Betting(%d)", JD21_CState_Betting),                           // 押底注
	JD21_CState_CutCards:             fmt.Sprintf("CutCards(%d)", JD21_CState_CutCards),                         // 切牌
	JD21_CState_Deal:                 fmt.Sprintf("Deal(%d)", JD21_CState_Deal),                                 // 發牌
	JD21_CState_Insurance:            fmt.Sprintf("Insurance(%d)", JD21_CState_Insurance),                       // 保險
	JD21_CState_Settlement_Insurance: fmt.Sprintf("Settlement_Insurance(%d)", JD21_CState_Settlement_Insurance), // 保險結算
	JD21_CState_Playing:              fmt.Sprintf("Playing(%d)", JD21_CState_Playing),                           // 玩家操作
	JD21_CState_Settlement:           fmt.Sprintf("Settlement(%d)", JD21_CState_Settlement),                     // 莊家開牌/結算
}

// 玩家狀態
var g_map_strUserState_debug = map[int]string{
	JD21_UserState_EmptySeat: fmt.Sprintf("EmptySeat(%d)", JD21_UserState_EmptySeat), // 空位
	JD21_UserState_Normal:    fmt.Sprintf("Normal(%d)", JD21_UserState_Normal),       // 正常
	JD21_UserState_Watch:     fmt.Sprintf("Watch(%d)", JD21_UserState_Watch),         // 觀戰
}

// 玩家遊戲中狀態
var g_map_strPlayingState_debug = map[int]string{
	JD21_UserPlayState_WaitingStart:           fmt.Sprintf("WaitingStart(%d)", JD21_UserPlayState_WaitingStart),                     // 等待遊戲開始
	JD21_UserPlayState_Betting:                fmt.Sprintf("Betting(%d)", JD21_UserPlayState_Betting),                               // 開始下注
	JD21_UserPlayState_StopBettingOK:          fmt.Sprintf("StopBettingOK(%d)", JD21_UserPlayState_StopBettingOK),                   // 停止下注動畫播放完成
	JD21_UserPlayState_CutCards:               fmt.Sprintf("CutCards(%d)", JD21_UserPlayState_CutCards),                             // 切牌
	JD21_UserPlayState_CutCardsOK:             fmt.Sprintf("CutCardsOK(%d)", JD21_UserPlayState_CutCardsOK),                         // 切牌動畫播放完成
	JD21_UserPlayState_DealOK:                 fmt.Sprintf("DealOK(%d)", JD21_UserPlayState_DealOK),                                 // 發牌動畫播放完成
	JD21_UserPlayState_TakeInsurance:          fmt.Sprintf("TakeInsurance(%d)", JD21_UserPlayState_TakeInsurance),                   // 玩家是否買保險
	JD21_UserPlayState_TakeInsuranceOK:        fmt.Sprintf("TakeInsuranceOK(%d)", JD21_UserPlayState_TakeInsuranceOK),               // 玩家已買保險
	JD21_UserPlayState_DealerOK:               fmt.Sprintf("DealerOK(%d)", JD21_UserPlayState_DealerOK),                             // 莊家瞇牌動畫播放完成
	JD21_UserPlayState_Settlement_InsuranceOK: fmt.Sprintf("Settlement_InsuranceOK(%d)", JD21_UserPlayState_Settlement_InsuranceOK), // 保險結算動畫播放完成
	JD21_UserPlayState_Playing_Waiting:        fmt.Sprintf("Playing_Waiting(%d)", JD21_UserPlayState_Playing_Waiting),               // 玩家操作 - 等待中
	JD21_UserPlayState_Playing_MyTurn:         fmt.Sprintf("Playing_MyTurn(%d)", JD21_UserPlayState_Playing_MyTurn),                 // 玩家操作 - 輪到我了
	JD21_UserPlayState_Playing_Done:           fmt.Sprintf("Playing_Done(%d)", JD21_UserPlayState_Playing_Done),                     // 玩家操作 - 完成操作
	JD21_UserPlayState_Settlement_DealerOK:    fmt.Sprintf("Settlement_DealerOK(%d)", JD21_UserPlayState_Settlement_DealerOK),       // 莊家開牌/結算
}

//---------------------------------------------------------------------------
// 莊家與玩家比牌
func CompareCards(pDealerSet *CardSetInfo, pPlayerSet *CardSetInfo) int {
	if pDealerSet == nil || pPlayerSet == nil {
		return JD21_Compare_Even
	}

	// 莊家是 BlackJack
	if pDealerSet.bIsBlackJack {
		if pPlayerSet.bIsBlackJack {
			return JD21_Compare_Even
		} else {
			return JD21_Compare_Dealer
		}
	} else {
		if pPlayerSet.bIsBlackJack {
			return JD21_Compare_Player
		}
	}

	// 莊家爆牌
	if pDealerSet.bIsBusted {
		if pPlayerSet.bIsBusted {
			return JD21_Compare_Even
		} else {
			return JD21_Compare_Player
		}
	} else {
		if pPlayerSet.bIsBusted {
			return JD21_Compare_Dealer
		}
	}

	// 一般點數
	i32DealerPoint, i32PlayerPoint := int32(-1), int32(-1)
	// 莊家點數
	if pDealerSet.i32Point_2 > 0 {
		if pDealerSet.i32Point_2 > 21 {
			i32DealerPoint = pDealerSet.i32Point_1
		} else {
			i32DealerPoint = pDealerSet.i32Point_2
		}
	} else {
		i32DealerPoint = pDealerSet.i32Point_1
	}

	// 玩家點數
	if pPlayerSet.i32Point_2 > 0 {
		if pPlayerSet.i32Point_2 > 21 {
			i32PlayerPoint = pPlayerSet.i32Point_1
		} else {
			i32PlayerPoint = pPlayerSet.i32Point_2
		}
	} else {
		i32PlayerPoint = pPlayerSet.i32Point_1
	}

	// 比點數
	if i32DealerPoint > i32PlayerPoint {
		return JD21_Compare_Dealer
	} else if i32DealerPoint < i32PlayerPoint {
		return JD21_Compare_Player
	} else {
		return JD21_Compare_Even
	}
}

//---------------------------------------------------------------------------
