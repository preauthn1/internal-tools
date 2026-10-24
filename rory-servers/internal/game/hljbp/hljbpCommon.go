package hljbp

import (
	"fmt"
)

// 遊戲設定
const (
	ConfigLabel         string = "GameSetting" // 設定檔標籤
	ClientStartUpSeatNo int    = 1             // client 起始座位號
)

// 遊戲狀態
const (
	StateWaitingStart int = iota // 等待遊戲開始
	StateBetting                 // 開始下注
	StateSettlement              // 結算
	StateGameOver                // 遊戲結束
)

// ActionName 遊戲狀態名稱
var ActionName = map[int]string{
	StateWaitingStart: "ActionNameWaitingStart", // 等待遊戲開始
	StateBetting:      "ActionNameBetting",      // 開始下注
	StateSettlement:   "ActionNameSettlement",   // 結算
	StateGameOver:     "ActionNameGameOver",     // 遊戲結束
}

// TimeOut 類別
var TimeOut = map[int]string{
	StateBetting:    "TimeoutBetting",    // 開始下注
	StateSettlement: "TimeoutSettlement", // 底注結算
}

// 玩家狀態
const (
	UserStateEmptySeat int = iota // 空位
	UserStateNormal               // 正常
	UserStateWatch                // 觀戰
)

// 玩家下注結果
const (
	BetResultSuccess           = 0  // 下注成功
	BetResultNotEnoughMoney    = -1 // 玩家餘額不足
	BetResultBettingStopped    = -2 // 已停止下注
	BetResultBetOverUpperLimit = -3 // 押注超過上限
	BetResultBetSmallThanZero  = -4 // 押注金額<=0
	BetResultBetFailed         = -5 // 壓注失敗
)

// 玩家輸贏
const (
	WinStateNone int = -1 // 發牌前
	WinStateEven     = 0  // 和
	WinStateWin      = 1  // 贏
	WinStateLose     = 2  // 輸
)

/*後台記錄用*/

// StrUserState 玩家狀態
var StrUserState = map[int]string{
	UserStateEmptySeat: "空位", // 空位
	UserStateNormal:    "正常", // 正常
	UserStateWatch:     "观战", // 觀戰
}

// StrBetResult 下注結果
var StrBetResult = map[int]string{
	BetResultSuccess:           "下注成功",     // 下注成功
	BetResultNotEnoughMoney:    "玩家余额不足",   // 玩家餘額不足
	BetResultBettingStopped:    "已停止下注",    // 已停止下注
	BetResultBetOverUpperLimit: "押注超过上限",   // 押注超過上限
	BetResultBetSmallThanZero:  "押注金额小於 0", // 押注金額<=0
	BetResultBetFailed:         "压注失败",     // 壓注失敗
}

/*for debug*/

// StrGameStateDebug 遊戲狀態
var StrGameStateDebug = map[int]string{
	StateWaitingStart: fmt.Sprintf("WaitingStart(%d)", StateWaitingStart),    // 等待遊戲開始
	StateBetting:      fmt.Sprintf("StartBetting(%d)", StateBetting),         // 開始下注
	StateSettlement:   fmt.Sprintf("BetBaseSettlement(%d)", StateSettlement), // 結算
	StateGameOver:     fmt.Sprintf("GameOver(%d)", StateGameOver),            // 遊戲結束
}

// StrUserStateDebug 玩家狀態
var StrUserStateDebug = map[int]string{
	UserStateEmptySeat: fmt.Sprintf("EmptySeat(%d)", UserStateEmptySeat), // 空位
	UserStateNormal:    fmt.Sprintf("Normal(%d)", UserStateNormal),       // 正常
	UserStateWatch:     fmt.Sprintf("Watch(%d)", UserStateWatch),         // 觀戰
}

// StrUserWinStateDebug 玩家輸贏
var StrUserWinStateDebug = map[int]string{
	WinStateNone: fmt.Sprintf("None(%d)", WinStateNone), // 發牌前
	WinStateEven: fmt.Sprintf("Even(%d)", WinStateEven), // 和
	WinStateWin:  fmt.Sprintf("Win(%d)", WinStateWin),   // 贏
	WinStateLose: fmt.Sprintf("Lose(%d)", WinStateLose), // 輸
}
