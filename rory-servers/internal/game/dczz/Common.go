package dczz

import "fmt"

const (
	IndexBase					int = 1			// 索引值起始值
	PokerSetCount				int32 = 1		// 使用幾副牌
	Rate_Normal					int64 = 1		// 一般比牌賠率
	Rate_War					int64 = 1		// 戰爭比牌賠率
	Rate_TieAward				int64 = 10		// 平手額賠率
	Count_CardCount_Normal		int = 1			// 一般階段手牌張數
	Count_CardCount_War			int = 5			// 戰爭階段手牌張數
	Count_HideCardCount			int = 3			// 戰爭階段蓋牌張數
	TimeDuration_GameOver		int64 = 3		// 遊戲結束 timeout
	Times_BetDouble				int64 = 1		// 加注額為底注額的幾倍
)

// 遊戲狀態
const (
	GameState_WaitingPlayer				int = iota	// 等待玩家加入
	GameState_WaitingStart							// 等待遊戲開始
	GameState_Betting								// 押注
	GameState_StopBetting							// 停止下注
	GameState_Shuffle								// 洗牌
	GameState_Deal_Normal							// 一般發牌
	GameState_Settlement_Normal						// 一般結算
	GameState_WarSelection							// 選擇是否進入戰爭
	GameState_WarPhaseAnimation						// 進入戰爭階段動畫
	GameState_Deal_War								// 戰爭發牌
	GameState_Settlement_War						// 戰爭結算
	GameState_GameOver								// 遊戲結束
)

// timeout 類別
const (
	Timeout_WaitingStart				int = iota	// 等待遊戲開始
	Timeout_Betting									// 下注
	Timeout_StopBettingAnimation 					// 停止下注動畫
	Timeout_Shuffle                       			// 洗牌動畫
	Timeout_Deal_Normal                 			// 一般發牌動畫
	Timeout_Settlement_Normal            			// 一般結算
	Timeout_WarSelection							// 選擇是否進入戰爭
	Timeout_WarPhaseBegin							// 進入戰爭階段動畫
	Timeout_Deal_War                 				// 戰爭發牌動畫
	Timeout_Settlement_War              			// 戰爭結算
	Timeout_GameOver								// 遊戲結束
)

// 玩家狀態
const (
	UserState_EmptySeat				int = iota	// 空位
	UserState_Normal							// 正常
	UserState_Watch								// 觀戰
)

// 玩家遊戲狀態
const (
	UserPlayState_WaitingStart				int = iota	// 等待遊戲開始
	UserPlayState_Betting								// 開始下注
	UserPlayState_StopBettingOK							// 停止下注動畫播放完成
	UserPlayState_ShuffleOK								// 洗牌動畫播放完成
	UserPlayState_DealOK_Normal							// 一般發牌動畫播放完成
	UserPlayState_Settlement_Normal_OK					// 一般結算動畫播放完成
	UserPlayState_WarSelection							// 選擇是否進入戰爭
	UserPlayState_WarSelectionOK						// 開始戰爭階段動畫播放完成
	UserPlayState_DealOK_War							// 戰爭發牌動畫播放完成
	UserPlayState_Settlement_War_OK						// 戰爭結算動畫播放完成
)

// 押注類別
const (
	BetType_Base					int = iota	// 底注
	BetType_Side								// 平手注
	BetType_Max
)

// 押注結果
const (
	BetResult_Auto						= 1		// 自動下底注
	BetResult_Success					= 0		// 下注成功
	BetResult_BetFailed					= -1	// 下注失敗
	BetResult_NotEnoughMoney			= -2	// 玩家餘額不足
	BetResult_BettingStopped			= -3	// 已停止下注
	BetResult_OverUpperLimit			= -4	// 總押注額超過上限
	BetResult_BetSmallThanZero			= -5	// 押注金額<=0
)

// 遊戲階段
const (
	GamePhase_Normal				int = iota	// 一般階段
	GamePhase_War								// 戰爭階段
	GamePhase_Max
)

// 戰爭階段選項
const (
	WarSelection_Fold				int = iota	// 棄牌
	WarSelection_War							// 戰爭
	WarSelection_Max
)

// 發牌器發牌方式
const (
	DealType_Natural			int = iota	// 自然發牌（預設值）
	DealType_Assigned						// 指定發牌
	DealType_Max
)

// 玩家輸贏
const (
	Settlement_PlayerLose		int = iota	// 玩家輸
	Settlement_PlayerWin					// 玩家贏
	Settlement_Tie							// 平手
)

//---------------------------------------------------------------------------
var g_i64BetLimit_Base = int64(0)	// 底注額上限
var g_i64BetLimit_Side = int64(0)	// 平手額上限
var g_iShuffleRestCardCount = 0		// 牌盒內張數低於此值時洗牌

//---------------------------------------------------------------------------
//  ----- 後台記錄用 -----
// 玩家狀態
var g_map_strUserState = map[int]string{
	UserState_EmptySeat: "空位", // 空位
	UserState_Normal:    "正常", // 正常
	UserState_Watch:     "观战", // 觀戰
}

// 押注類別
var g_map_strBetType = map[int]string{
	BetType_Base: "底注",  // 底注
	BetType_Side: "平手注", // 平手注
}

// 下注結果
var g_map_strBetResult = map[int]string{
	BetResult_Auto:             "自動下底注",    // 自動下底注
	BetResult_Success:          "下注成功",     // 下注成功
	BetResult_BetFailed:        "下注失败",     // 下注失敗
	BetResult_NotEnoughMoney:   "玩家余额不足",   // 玩家餘額不足
	BetResult_BettingStopped:   "已停止下注",    // 已停止下注
	BetResult_OverUpperLimit:   "总压住额超过上限", // 總押注額超過上限
	BetResult_BetSmallThanZero: "押注金额小於 0", // 押注金額<=0
}

//---------------------------------------------------------------------------
// ----- for debug -----
// 遊戲狀態
var g_map_strGameState_debug = map[int]string{
	GameState_WaitingPlayer:     fmt.Sprintf("WaitingPlayer(%d)", GameState_WaitingPlayer),         // 等待玩家加入
	GameState_WaitingStart:      fmt.Sprintf("WaitingStart(%d)", GameState_WaitingStart),           // 等待遊戲開始
	GameState_Betting:           fmt.Sprintf("Betting(%d)", GameState_Betting),                     // 押注
	GameState_StopBetting:       fmt.Sprintf("StopBetting(%d)", GameState_StopBetting),             // 停止下注
	GameState_Shuffle:           fmt.Sprintf("Shuffle(%d)", GameState_Shuffle),                     // 洗牌
	GameState_Deal_Normal:       fmt.Sprintf("Deal_Normal(%d)", GameState_Deal_Normal),             // 一般發牌
	GameState_Settlement_Normal: fmt.Sprintf("Settlement_Normal(%d)", GameState_Settlement_Normal), // 一般結算
	GameState_WarSelection:      fmt.Sprintf("WarSelection(%d)", GameState_WarSelection),           // 選擇是否進入戰爭
	GameState_WarPhaseAnimation: fmt.Sprintf("WarPhaseAnimation(%d)", GameState_WarPhaseAnimation), // 進入戰爭階段動畫
	GameState_Deal_War:          fmt.Sprintf("Deal_War(%d)", GameState_Deal_War),                   // 戰爭發牌
	GameState_Settlement_War:    fmt.Sprintf("Settlement_War(%d)", GameState_Settlement_War),       // 戰爭結算
	GameState_GameOver:          fmt.Sprintf("GameOver(%d)", GameState_GameOver),                   // 遊戲結束
}

// 玩家狀態
var g_map_strUserState_debug = map[int]string{
	UserState_EmptySeat: fmt.Sprintf("EmptySeat(%d)", UserState_EmptySeat), // 空位
	UserState_Normal:    fmt.Sprintf("Normal(%d)", UserState_Normal),       // 正常
	UserState_Watch:     fmt.Sprintf("Watch(%d)", UserState_Watch),         // 觀戰
}

// 玩家遊戲中狀態
var g_map_strPlayingState_debug = map[int]string{
	UserPlayState_WaitingStart:         fmt.Sprintf("WaitingStart(%d)", UserPlayState_WaitingStart),                 // 等待遊戲開始
	UserPlayState_Betting:              fmt.Sprintf("Betting(%d)", UserPlayState_Betting),                           // 開始下注
	UserPlayState_StopBettingOK:        fmt.Sprintf("StopBettingOK(%d)", UserPlayState_StopBettingOK),               // 停止下注動畫播放完成
	UserPlayState_ShuffleOK:            fmt.Sprintf("ShuffleOK(%d)", UserPlayState_ShuffleOK),                       // 洗牌動畫播放完成
	UserPlayState_DealOK_Normal:        fmt.Sprintf("DealOK_Normal(%d)", UserPlayState_DealOK_Normal),               // 一般發牌動畫播放完成
	UserPlayState_Settlement_Normal_OK: fmt.Sprintf("Settlement_Normal_OK(%d)", UserPlayState_Settlement_Normal_OK), // 一般結算動畫播放完成
	UserPlayState_WarSelection:         fmt.Sprintf("WarSelection(%d)", UserPlayState_WarSelection),                 // 選擇是否進入戰爭
	UserPlayState_WarSelectionOK:       fmt.Sprintf("WarSelectionOK(%d)", UserPlayState_WarSelectionOK),             // 開始戰爭階段動畫播放完成
	UserPlayState_DealOK_War:           fmt.Sprintf("DealOK_War(%d)", UserPlayState_DealOK_War),                     // 戰爭發牌動畫播放完成
	UserPlayState_Settlement_War_OK:    fmt.Sprintf("Settlement_War_OK(%d)", UserPlayState_Settlement_War_OK),       // 戰爭結算動畫播放完成
}

// 遊戲階段
var g_map_strGamePhase_debug = map[int]string{
	GamePhase_Normal: "一般階段", // 一般階段
	GamePhase_War:    "戰爭階段", // 戰爭階段
}

//---------------------------------------------------------------------------
