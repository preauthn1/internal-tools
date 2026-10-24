package jlbpk

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"

	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

const (
	JLBPK_PokerSetCount       int32 = 1 // 使用幾副牌
	JLBPK_HandCardCount       int   = 5 // 手牌張數
	JLBPK_BetBasePayRate      int   = 1 // 莊家賠付玩家底注倍數
	JLBPK_ClientStartUpSeatNo int   = 1 // client 起始座位號
)

const (
	JLBPK_STR_CannotBeAGame string = "不成局"
)

// 遊戲狀態
const (
	JLB_State_WaitingPlayer         int = iota // 等待玩家加入
	JLB_State_WaitingStart                     // 等待遊戲開始
	JLB_State_StartBettingAnimation            // 開始下注動畫
	JLB_State_StartBetting                     // 開始下注
	JLB_State_StopBettingAnimation             // 停止下注動畫
	JLB_State_Deal                             // 發牌
	JLB_State_BetDouble                        // 加注
	JLB_State_BetSideSettlement                // 邊注結算
	JLB_State_BetBaseSettlement                // 底注結算
	JLB_State_BetDoubleSettlement              // 加注結算
	JLB_State_GameOver                         // 遊戲結束
	DHC_State_Settlement                       // 下注結算
)

// 送給 client 的遊戲狀態
const (
	JLB_CState_WaitingStart        int = iota // 等待遊戲開始
	JLB_CState_StartBetting                   // 開始下注
	JLB_CState_Deal                           // 發牌
	JLB_CState_BetDouble                      // 加注
	JLB_CState_BetSideSettlement              // 邊注結算
	JLB_CState_BetBaseSettlement              // 底注結算
	JLB_CState_BetDoubleSettlement            // 加注結算
	DHC_CState_Settlement                     // 下注結算
)

// timeout 類別
const (
	JLB_Timeout_WaitingStart          int = iota // 等待遊戲開始
	JLB_Timeout_StartBettingAnimation            // 開始下注動畫
	JLB_Timeout_StartBetting                     // 開始下注
	JLB_Timeout_StopBettingAnimation             // 停止下注動畫
	JLB_Timeout_DealAnimation                    // 發牌
	JLB_Timeout_DoubleBetting                    // 玩家加注
	JLB_Timeout_BetSideSettlement                // 邊注結算
	JLB_Timeout_BetBaseSettlement                // 底注結算
	JLB_Timeout_BetDoubleSettlement              // 加注結算
	DHC_Timeout_Settlement                       // 下注結算
)

// 玩家狀態
const (
	JLB_UserState_EmptySeat int = iota // 空位
	JLB_UserState_Normal               // 正常
	JLB_UserState_Watch                // 觀戰
	JLB_UserState_Fold                 // 棄牌
)

// 玩家遊戲狀態
const (
	JLB_UserPlayState_WaitingStart           int = iota // 等待遊戲開始
	JLB_UserPlayState_StartBettingOK                    // 開始下注動畫播放完成
	JLB_UserPlayState_Betting                           // 開始下注
	JLB_UserPlayState_StopBettingOK                     // 停止下注動畫播放完成
	JLB_UserPlayState_DealOK                            // 發牌動畫播放完成
	JLB_UserPlayState_BetDouble                         // 玩家加注
	JLB_UserPlayState_ShowCardsOK                       // 開牌，邊注結算效果完成
	JLB_UserPlayState_Compensate_BetBaseOK              // 賠付，底注結算效果完成
	JLB_UserPlayState_Compensate_DoubleBetOK            // 賠付，加注結算效果完成
	DHC_UserPlayState_SettlementOK                      // 結算表演，下注結算效果完成
)

// 下注類別
const (
	JLB_BetType_Base       int32 = iota // 底注
	JLB_BetType_Side                    // 邊注
	JLB_BetType_ZhuangSide              // 莊家邊注
)

// 玩家下注結果
const (
	JLB_BetResult_Auto                  = 1  // 自動下底注
	JLB_BetResult_Success               = 0  // 下注成功
	JLB_BetResult_NotEnoughMoney        = -1 // 玩家餘額不足
	JLB_BetResult_BettingStopped        = -2 // 已停止下注
	JLB_BetResult_BetSideHasDone        = -3 // 已下過邊注
	JLB_BetResult_BaseBetOverUpperLimit = -4 // 底注超過上限
	JLB_BetResult_BetSideOverLowerLimit = -5 // 邊注小於下限
	JLB_BetResult_BetSmallThanZero      = -6 // 押注金額<=0
	JLB_BetResult_NotOpenZhuangSide     = -7 // 沒開啟莊家邊注
)

// 玩家加注結果
const (
	JLB_DoubleBetResult_Success          = 0  // 下注成功
	JLB_DoubleBetResult_NotEnoughMoney   = -1 // 玩家餘額不足
	JLB_DoublebetResult_BetDoubleHasDone = -2 // 已下過加注
)

// 發牌器發牌模式
const (
	DealType_Natural  int = iota // 自然發牌（預設值）
	DealType_RandomJP            // 隨機發彩金牌型
	DealType_Assigned            // 指定玩家手牌
	DealType_Max
)

// 玩家輸贏
const (
	JLB_WinState_None int = -1 // 發牌前
	JLB_WinState_Even     = 0  // 和
	JLB_WinState_Win      = 1  // 贏
	JLB_WinState_Lose     = 2  // 輸
)

// 是否開啟莊家邊注
const (
	ZhuangBetSideDisable int = iota // 關閉莊家彩金
	ZhuangBetSideEnable             //開啟
)

// 是否開啟複數邊注
const (
	PluralBetSideDisable int = iota // 關閉複數邊注
	PluralBetSideEnable             //開啟
)

// 海盜船模式
const (
	HDCModeDisable int = iota // 關閉海盜船模式
	HDCModeEnable             //開啟
)

//---------------------------------------------------------------------------
// 彩金
var g_i64JackpotBase = int64(0)       // 獎池回補基準值
var g_i64Jackpot = int64(0)           // 目前獎池彩金
var g_i64Contribution = int64(0)      //本次貢獻值
var g_i64TotalContribution = int64(0) // 累計貢獻值
var g_i64AwardTotal = int64(0)        // 累計出獎金額
var g_bIsReadJackpotFromDB = false    // 是否從DB讀出Jackpot資料

// 機器人
var g_iFoldByNoAKOdds_Robot = 0 // 機器人拿到高牌且無A、K時的棄牌機率
var g_iBetSideOdds_Robot = 0    // 機器人下邊注機率

// 彩金中獎牌型
var g_ary_i32JPPattern = []int32{
	pokerset.POKERPatternRoyalFlush,    // 皇家同花順
	pokerset.POKERPatternStraightFlush, // 同花順
	pokerset.POKERPatternFourOfAKind,   // 四條
	pokerset.POKERPatternFullHouse,     // 葫蘆
	pokerset.POKERPatternFlush,         // 同花
}

// 邊注
var g_i64DefaultBetSide = int64(0) // 預設邊注金額
//---------------------------------------------------------------------------
//  ----- 後台記錄用 -----
// 玩家狀態
var map_strUserState = map[int]string{
	JLB_UserState_EmptySeat: "空位", // 空位
	JLB_UserState_Normal:    "正常", // 正常
	JLB_UserState_Watch:     "观战", // 觀戰
	JLB_UserState_Fold:      "弃牌", // 棄牌
}

// 下注類別
var map_strBetType = map[int32]string{
	JLB_BetType_Base: "底注", // 底注
	JLB_BetType_Side: "边注", // 邊注
}

// 下注結果
var map_strBetResult = map[int]string{
	JLB_BetResult_Auto:                  "自動下底注",    // 自動下底注
	JLB_BetResult_Success:               "下注成功",     // 下注成功
	JLB_BetResult_NotEnoughMoney:        "玩家余额不足",   // 玩家餘額不足
	JLB_BetResult_BettingStopped:        "已停止下注",    // 已停止下注
	JLB_BetResult_BetSideHasDone:        "已下过边注",    // 已下過邊注
	JLB_BetResult_BaseBetOverUpperLimit: "底注超过上限",   // 底注超過上限
	JLB_BetResult_BetSideOverLowerLimit: "边注小于下限",   // 邊注小於下限
	JLB_BetResult_BetSmallThanZero:      "押注金额小於 0", // 押注金額<=0
}

// 加注結果
var map_strDoubleBetResult = map[int]string{
	JLB_DoubleBetResult_Success:        "加注成功",   // 加注成功
	JLB_DoubleBetResult_NotEnoughMoney: "玩家余额不足", // 玩家餘額不足
}

//---------------------------------------------------------------------------
// ----- for debug -----
// 遊戲狀態
var map_strGameState_debug = map[int]string{
	JLB_State_WaitingPlayer:         fmt.Sprintf("WaitingPlayer(%d)", JLB_State_WaitingPlayer),                 // 等待玩家加入
	JLB_State_WaitingStart:          fmt.Sprintf("WaitingStart(%d)", JLB_State_WaitingStart),                   // 等待遊戲開始
	JLB_State_StartBettingAnimation: fmt.Sprintf("StartBettingAnimation(%d)", JLB_State_StartBettingAnimation), // 開始下注動畫
	JLB_State_StartBetting:          fmt.Sprintf("StartBetting(%d)", JLB_State_StartBetting),                   // 開始下注
	JLB_State_StopBettingAnimation:  fmt.Sprintf("StopBettingAnimation(%d)", JLB_State_StopBettingAnimation),   // 停止下注動畫
	JLB_State_Deal:                  fmt.Sprintf("Deal(%d)", JLB_State_Deal),                                   // 發牌
	JLB_State_BetDouble:             fmt.Sprintf("BetDouble(%d)", JLB_State_BetDouble),                         // 加注
	JLB_State_BetSideSettlement:     fmt.Sprintf("BetSideSettlement(%d)", JLB_State_BetSideSettlement),         // 邊注結算
	JLB_State_BetBaseSettlement:     fmt.Sprintf("BetBaseSettlement(%d)", JLB_State_BetBaseSettlement),         // 底注結算
	JLB_State_BetDoubleSettlement:   fmt.Sprintf("BetDoubleSettlement(%d)", JLB_State_BetDoubleSettlement),     // 加注結算
	JLB_State_GameOver:              fmt.Sprintf("GameOver(%d)", JLB_State_GameOver),                           // 遊戲結束
	DHC_State_Settlement:            fmt.Sprintf("Settlement(%d)", DHC_State_Settlement),                       // 下注結算
}

// 玩家狀態
var map_strUserState_debug = map[int]string{
	JLB_UserState_EmptySeat: fmt.Sprintf("EmptySeat(%d)", JLB_UserState_EmptySeat), // 空位
	JLB_UserState_Normal:    fmt.Sprintf("Normal(%d)", JLB_UserState_Normal),       // 正常
	JLB_UserState_Watch:     fmt.Sprintf("Watch(%d)", JLB_UserState_Watch),         // 觀戰
	JLB_UserState_Fold:      fmt.Sprintf("Fold(%d)", JLB_UserState_Fold),           // 棄牌
}

// 玩家遊戲中狀態
var map_strPlayingState_debug = map[int]string{
	JLB_UserPlayState_WaitingStart:           fmt.Sprintf("WaitingStart(%d)", JLB_UserPlayState_WaitingStart),                     // 等待遊戲開始
	JLB_UserPlayState_StartBettingOK:         fmt.Sprintf("StartBettingOK(%d)", JLB_UserPlayState_StartBettingOK),                 // 開始下注動畫播放完成
	JLB_UserPlayState_Betting:                fmt.Sprintf("Betting(%d)", JLB_UserPlayState_Betting),                               // 開始下注
	JLB_UserPlayState_StopBettingOK:          fmt.Sprintf("StopBettingOK(%d)", JLB_UserPlayState_StopBettingOK),                   // 停止下注動畫播放完成
	JLB_UserPlayState_DealOK:                 fmt.Sprintf("DealOK(%d)", JLB_UserPlayState_DealOK),                                 // 發牌動畫播放完成
	JLB_UserPlayState_BetDouble:              fmt.Sprintf("BetDouble(%d)", JLB_UserPlayState_BetDouble),                           // 玩家加注
	JLB_UserPlayState_ShowCardsOK:            fmt.Sprintf("ShowCardsOK(%d)", JLB_UserPlayState_ShowCardsOK),                       // 開牌，邊注結算效果完成
	JLB_UserPlayState_Compensate_BetBaseOK:   fmt.Sprintf("Compensate_BetBaseOK(%d)", JLB_UserPlayState_Compensate_BetBaseOK),     // 賠付，底注結算效果完成
	JLB_UserPlayState_Compensate_DoubleBetOK: fmt.Sprintf("Compensate_DoubleBetOK(%d)", JLB_UserPlayState_Compensate_DoubleBetOK), // 賠付，加注結算效果完成
	DHC_UserPlayState_SettlementOK:           fmt.Sprintf("SettlementOK(%d)", DHC_UserPlayState_SettlementOK),                     // 結算表演，下注結算效果完成
}

// 送給 client 的遊戲狀態
var map_strClientGameState_debug = map[int]string{
	JLB_CState_WaitingStart:        fmt.Sprintf("WaitingStart(%d)", JLB_CState_WaitingStart),               // 等待遊戲開始
	JLB_CState_StartBetting:        fmt.Sprintf("StartBetting(%d)", JLB_CState_StartBetting),               // 開始下注
	JLB_CState_Deal:                fmt.Sprintf("Deal(%d)", JLB_CState_Deal),                               // 發牌
	JLB_CState_BetDouble:           fmt.Sprintf("BetDouble(%d)", JLB_CState_BetDouble),                     // 加注
	JLB_CState_BetSideSettlement:   fmt.Sprintf("BetSideSettlement(%d)", JLB_CState_BetSideSettlement),     // 邊注結算
	JLB_CState_BetDoubleSettlement: fmt.Sprintf("BetDoubleSettlement(%d)", JLB_CState_BetDoubleSettlement), // 加注結算
	JLB_CState_BetBaseSettlement:   fmt.Sprintf("BetBaseSettlement(%d)", JLB_CState_BetBaseSettlement),     // 底注結算
	DHC_CState_Settlement:          fmt.Sprintf("Settlement(%d)", DHC_CState_Settlement),                   // 下注結算
}

// 玩家輸贏
var map_strUserWinState_debug = map[int]string{
	JLB_WinState_None: fmt.Sprintf("None(%d)", JLB_WinState_None), // 發牌前
	JLB_WinState_Even: fmt.Sprintf("Even(%d)", JLB_WinState_Even), // 和
	JLB_WinState_Win:  fmt.Sprintf("Win(%d)", JLB_WinState_Win),   // 贏
	JLB_WinState_Lose: fmt.Sprintf("Lose(%d)", JLB_WinState_Lose), // 輸
}

//---------------------------------------------------------------------------

// 開啟莊家邊注 0關閉1開啟
var g_iZhuangBetSide = 0

// 開啟複數邊注 0關閉1開啟
var g_iPluralBetSide = 0

// 複數邊注上限
var g_iPluralBetSideLimit = 0

// 海盜船模式 0關閉1開啟
var g_iHDCMode = 0

func (pGame *Logic) getDealerWeightTable(configs *config.Group) {
	cfgs := configs.GetConfig("gameserver")
	g_map_pDealerWeightTable = make(map[int]*DealerWeightTable)
	g_iRoundLimit = cfgs.GetConfigInt("SystemPool", "RoundLimit") // 判定局數
	// 機率表
	iTableCount := cfgs.GetConfigInt("WeightTable", "TableCount") // 權重表張數
	str := ""
	for iNo := 0; iNo < iTableCount; iNo++ {
		if _, ok := g_map_pDealerWeightTable[iNo]; !ok {
			g_map_pDealerWeightTable[iNo] = new(DealerWeightTable)
		}

		strName := fmt.Sprintf("%d", iNo)
		str = cfgs.GetConfigValue(strName, "Pattern")
		arystrPattern := strings.Split(str, ",")
		str = cfgs.GetConfigValue(strName, "Weight")
		arystrWeight := strings.Split(str, ",")

		for i := 0; i < len(arystrPattern); i++ {
			pData := new(DealerWeightData)
			i64Pattern, _ := strconv.ParseInt(arystrPattern[i], 10, 32)
			pData.i32CardPattern = int32(i64Pattern)
			if i < len(arystrWeight) {
				i64Weight, _ := strconv.ParseInt(arystrWeight[i], 10, 32)
				pData.i32Weight = int32(i64Weight)
			}

			g_map_pDealerWeightTable[iNo].ary_pWeightData = append(g_map_pDealerWeightTable[iNo].ary_pWeightData, pData)
		}
	}

	// 依牌型大小降冪排序並計算亂數範圍
	for _, pTable := range g_map_pDealerWeightTable {
		if pTable != nil {
			// 排序
			iLen := len(pTable.ary_pWeightData)
			for i := 0; i < iLen; i++ {
				for j := i; j < iLen; j++ {
					if pTable.ary_pWeightData[i].i32CardPattern < pTable.ary_pWeightData[j].i32CardPattern {
						pTable.ary_pWeightData[i], pTable.ary_pWeightData[j] = pTable.ary_pWeightData[j], pTable.ary_pWeightData[i]
					}
				}
			}

			// 計算亂數範圍
			i32Range := int32(0)
			for i := 0; i < len(pTable.ary_pWeightData); i++ {
				if pTable.ary_pWeightData[i].i32Weight <= 0 {
					pTable.ary_pWeightData[i].i32RndRangeMin = -1
					pTable.ary_pWeightData[i].i32RndRangeMax = -1
					continue
				}

				pTable.ary_pWeightData[i].i32RndRangeMin = i32Range + 1
				i32Range += pTable.ary_pWeightData[i].i32Weight
				pTable.ary_pWeightData[i].i32RndRangeMax = i32Range
			}
			pTable.i32TotalWeight = i32Range
		}
	}
}
