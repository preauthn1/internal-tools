package fgsb

import (
	"fmt"
)

// 遊戲設定
const (
	ConfigLabel         string = "GameSetting" // 設定檔標籤
	PokerSetCount       int    = 1             // 使用幾副牌
	HasKing             bool   = false         // 是否有鬼牌
	HasPizi             bool   = false         // 使否有皮子
	HandCardCount       int    = 3             // 手牌張數
	BetBasePayRate      int    = 1             // 莊家賠付玩家底注倍數
	ClientStartUpSeatNo int    = 1             // client 起始座位號
	CannotBeAGame       string = "不成局"         // 不成局
)

// 遊戲狀態
const (
	StateWaitingStart      int = iota // 等待遊戲開始
	StateBetting                      // 開始下注
	StateStopBetting                  // 停止下注
	StateDeal                         // 發牌
	StateBetSideSettlement            // 邊注、特殊牌型結算
	StateBetDouble                    // 加注
	StateSettlement                   // 底注、加注結算
	StateGameOver                     // 遊戲結束
)

// ActionName 遊戲狀態名稱
var ActionName = map[int]string{
	StateWaitingStart:      "ActionNameWaitingStart",      // 等待遊戲開始
	StateBetting:           "ActionNameBetting",           // 開始下注
	StateStopBetting:       "ActionNameStopBetting",       // 停止下注
	StateDeal:              "ActionNameDeal",              // 發牌
	StateBetSideSettlement: "ActionNameBetSideSettlement", // 边注、特殊牌型結算
	StateBetDouble:         "ActionNameBetDouble",         // 加注
	StateSettlement:        "ActionNameSettlement",        // 底注、加注結算
	StateGameOver:          "ActionNameGameOver",          // 遊戲結束
}

// TimeOut 類別
var TimeOut = map[int]string{
	StateWaitingStart:      "TimeoutWaitingStart",         // 等待遊戲開始
	StateBetting:           "TimeoutBetting",              // 開始下注
	StateStopBetting:       "TimeoutStopBettingAnimation", // 停止下注
	StateDeal:              "TimeoutDealAnimation",        // 發牌
	StateBetSideSettlement: "TimeoutBetSideSettlement",    // 边注、特殊牌型結算
	StateBetDouble:         "TimeoutDoubleBetting",        // 玩家加注
	StateSettlement:        "TimeoutSettlement",           // 底注結算
}

// 玩家狀態
const (
	UserStateEmptySeat int = iota // 空位
	UserStateNormal               // 正常
	UserStateWatch                // 观战
	UserStateFold                 // 弃牌
)

// 选位結果
const (
	SelectSeatResultSuccess      = 0  // 选位成功
	SelectSeatResultNoSeat       = -1 // 座位不存在
	SelectSeatResultHaveUser     = -2 // 座位已有玩家
	SelectSeatResultSelectFailed = -3 // 选位失敗
)

// 下注類別
const (
	BetTypeBase int32 = 1 // 底注
	BetTypeSide int32 = 2 // 边注
)

// 玩家下注結果
const (
	BetResultSuccess               = 0  // 下注成功
	BetResultNotEnoughMoney        = -1 // 玩家餘額不足
	BetResultBettingStopped        = -2 // 已停止下注
	BetResultBaseBetOverUpperLimit = -3 // 底注超過上限
	BetResultBetSmallThanZero      = -4 // 押注金額<=0
	BetResultNeverBetBase          = -5 // 尚未下過底注
	BetResultBetFailed             = -6 // 壓注失敗
)

// 玩家加注結果
const (
	DoubleBetResultSuccess          = 1  // 下注成功
	DoubleBetResultFold             = 0  // 棄牌
	DoubleBetResultNotEnoughMoney   = -1 // 玩家餘額不足
	DoublebetResultBetDoubleStopped = -2 // 已停止加注
	DoublebetResultBetDoubleHasDone = -3 // 已下過加注
	DoublebetResultBetDoubleFailed  = -4 // 加注失敗
)

// 發牌器發牌模式
const (
	DealTypeNatural  int = iota // 自然發牌（預設值）
	DealTypeRandomJP            // 隨機發彩金牌型
	DealTypeAssigned            // 指定玩家手牌
	DealTypeMax
)

// 玩家輸贏
const (
	WinStateNone int = -1 // 發牌前
	WinStateEven     = 0  // 和
	WinStateWin      = 1  // 贏
	WinStateLose     = 2  // 輸
)

// SidePattern 邊注中獎牌型
var SidePattern = []int{
	PatternPair,          // 對子
	PatternFlush,         // 同花
	PatternStraight,      // 順子
	PatternThreeOfAKind,  // 三條
	PatternStraightFlush, // 同花順
}

/*後台記錄用*/

// StrUserState 玩家狀態
var StrUserState = map[int]string{
	UserStateEmptySeat: "空位", // 空位
	UserStateNormal:    "正常", // 正常
	UserStateWatch:     "观战", // 觀戰
	UserStateFold:      "弃牌", // 棄牌
}

// StrBetType 下注類別
var StrBetType = map[int32]string{
	BetTypeBase: "底注", // 底注
	BetTypeSide: "边注", // 邊注
}

// StrSelectSeatResult 選位結果
var StrSelectSeatResult = map[int]string{
	SelectSeatResultSuccess:      "选位成功",   // 選位成功
	SelectSeatResultNoSeat:       "座位不存在",  // 座位不存在
	SelectSeatResultHaveUser:     "座位已有玩家", // 座位已有玩家
	SelectSeatResultSelectFailed: "选位失敗",   // 選位失敗
}

// StrBetResult 下注結果
var StrBetResult = map[int]string{
	BetResultSuccess:               "下注成功",     // 下注成功
	BetResultNotEnoughMoney:        "玩家余额不足",   // 玩家余额不足
	BetResultBettingStopped:        "已停止下注",    // 已停止下注
	BetResultBaseBetOverUpperLimit: "底注超过上限",   // 底注超过上限
	BetResultBetSmallThanZero:      "押注金额小於 0", // 押注金額<=0
	BetResultNeverBetBase:          "尚未下过底注",   // 尚未下过底注
	BetResultBetFailed:             "压注失败",     // 压注失败
}

// StrDoubleBetResult 加注結果
var StrDoubleBetResult = map[int]string{
	DoubleBetResultSuccess:          "加注成功",   // 加注成功
	DoubleBetResultFold:             "弃牌",     // 棄牌
	DoubleBetResultNotEnoughMoney:   "玩家余额不足", // 玩家餘額不足
	DoublebetResultBetDoubleStopped: "玩家余额不足", // 已停止加注
	DoublebetResultBetDoubleHasDone: "已下过加注",  // 已下過加注
	DoublebetResultBetDoubleFailed:  "加注失敗",   // 加注失敗
}

/*for debug*/

// StrGameStateDebug 遊戲狀態
var StrGameStateDebug = map[int]string{
	StateWaitingStart:      fmt.Sprintf("WaitingStart(%d)", StateWaitingStart),           // 等待遊戲開始
	StateBetting:           fmt.Sprintf("StartBetting(%d)", StateBetting),                // 開始下注
	StateStopBetting:       fmt.Sprintf("StopBetting(%d)", StateStopBetting),             // 停止下注
	StateDeal:              fmt.Sprintf("Deal(%d)", StateDeal),                           // 發牌
	StateBetSideSettlement: fmt.Sprintf("BetSideSettlement(%d)", StateBetSideSettlement), // 邊注結算
	StateBetDouble:         fmt.Sprintf("BetDouble(%d)", StateBetDouble),                 // 加注
	StateSettlement:        fmt.Sprintf("BetBaseSettlement(%d)", StateSettlement),        // 底注、加注結算
	StateGameOver:          fmt.Sprintf("GameOver(%d)", StateGameOver),                   // 遊戲結束
}

// StrUserStateDebug 玩家狀態
var StrUserStateDebug = map[int]string{
	UserStateEmptySeat: fmt.Sprintf("EmptySeat(%d)", UserStateEmptySeat), // 空位
	UserStateNormal:    fmt.Sprintf("Normal(%d)", UserStateNormal),       // 正常
	UserStateWatch:     fmt.Sprintf("Watch(%d)", UserStateWatch),         // 觀戰
	UserStateFold:      fmt.Sprintf("Fold(%d)", UserStateFold),           // 棄牌
}

// StrUserWinStateDebug 玩家輸贏
var StrUserWinStateDebug = map[int]string{
	WinStateNone: fmt.Sprintf("None(%d)", WinStateNone), // 發牌前
	WinStateEven: fmt.Sprintf("Even(%d)", WinStateEven), // 和
	WinStateWin:  fmt.Sprintf("Win(%d)", WinStateWin),   // 贏
	WinStateLose: fmt.Sprintf("Lose(%d)", WinStateLose), // 輸
}
