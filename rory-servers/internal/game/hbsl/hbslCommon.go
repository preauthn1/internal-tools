package hbsl

// 設定檔檔名
const ConfigName_Game = "GameSetting"

// timeout 類別
const (
	Timeout_Grab       int = iota // 開始搶紅包
	Timeout_Settlement            // 結算
)

// 遊戲狀態
const (
	State_Grab       int = iota // 開始搶紅包
	State_Settlement            // 結算
)

// 遊戲狀態名稱
const (
	ActionName_Grab       string = "ActionName_Start"
	ActionName_Settlement string = "ActionName_End"
)

// UserInfo種類
const (
	UserInfoKind_Grab  int = iota // 搶紅包玩家
	UserInfoKind_All              // 線上所有玩家
	UserInfoKind_Empty            // 空玩家資訊
)

// 搶紅包結果
const (
	GrabResult_Success        = 0  // 搶紅包成功
	GrabResult_NotEnoughMoney = -1 // 玩家餘額不足
	GrabResult_GrabingStopped = -2 // 已停止搶紅包
	GrabResult_AlreadyGrab    = -3 // 已搶過紅包
	GrabResult_BetFailed      = -4 // 搶紅包失敗
)

// 發紅包結果
const (
	FaHBResult_Success         = 0  // 發紅包成功
	FaHBResult_NotEnoughMoney  = -1 // 玩家餘額不足
	FaHBResult_AlreadyFaHB     = -2 // 已發過紅包
	FaHBResult_HBMoneyError    = -3 // 紅包金額錯誤
	FaHBResult_BombNumberError = -4 // 雷號錯誤
	FaHBResult_BetFailed       = -5 // 發紅包失敗
)

// 線上玩家數量
const AllUserInfoMax = 20

/* 后台记录用 */

// 搶紅包結果
var g_map_strGrabResult = map[int]string{
	GrabResult_Success:        "抢红包成功",  // 搶紅包成功
	GrabResult_NotEnoughMoney: "玩家剩余不足", // 玩家餘額不足
	GrabResult_GrabingStopped: "已停止抢红包", // 已停止搶紅包
	GrabResult_AlreadyGrab:    "已抢过红包",  // 已搶過紅包
	GrabResult_BetFailed:      "抢红包失败",  // 搶紅包失敗
}

// 搶紅包結果
var g_map_strFaHBResult = map[int]string{
	FaHBResult_Success:         "发红包成功",  // 發紅包成功
	FaHBResult_NotEnoughMoney:  "玩家剩余不足", // 玩家餘額不足
	FaHBResult_AlreadyFaHB:     "已发过红包",  // 已發過紅包
	FaHBResult_HBMoneyError:    "红包金额错误", // 紅包金額錯誤
	FaHBResult_BombNumberError: "雷号错误",   // 雷號錯誤
	FaHBResult_BetFailed:       "发红包失败",  // 發紅包失敗
}
