package jdesydLogic

import (
	"fmt"
	Global "example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

//---------------------------------------------------------------------------
const (
	DefaultRoundLimit			int = 10		// 預設判定局數
	DefaultAverageBet			int64 = 60		// 預設平均押注值
	DefaultFactor				int = 20		// 預設係數定值
)

// 莊家補牌表
type SupplementData struct {
	iOdds     int   // 補牌機率，採千分比
	i32Points int32 // 最後補到多少點，若值為 0 表示為自然發牌
	iRndRange int   // 亂數值範圍
}

func (pData *SupplementData) Copy(pOther *SupplementData) {
	if pOther == nil {
		return
	}

	pData.iOdds = pOther.iOdds
	pData.i32Points = pOther.i32Points
	pData.iRndRange = pOther.iRndRange
}

type DealerSupplement struct {
	ary_pData []*SupplementData // 莊家補牌表
	iOddsSum  int               // 所有機率加總
}

func (pData *DealerSupplement) Copy(pOther *DealerSupplement) {
	if pOther == nil {
		return
	}

	iLen := len(pData.ary_pData)

	for i := 0; i < len(pOther.ary_pData); i++ {
		if i < iLen {
			pData.ary_pData[i].Copy(pOther.ary_pData[i])
		} else {
			pSupplement := new(SupplementData)
			pSupplement.Copy(pOther.ary_pData[i])
			pData.ary_pData = append(pData.ary_pData, pSupplement)
		}
	}
}

// 動態水線區間表
type PoolDynamicSetting struct {
	iAreaNo        int   // 區間編號
	iXVariation    int   // 變數加減，千分比，例如值 50 表示 0.050
	iConditionOdds int   // 輸贏區間
	i64AreaValue   int64 // 淨分區間值範圍
}

// 遊戲輸贏結果 DB 儲存資料
type GameDataDB struct {
	AverageBet int64 // 平均押注
	Factor     int   // 係數定值
	BetCount   int   // 押注次數
	BetMoney   int64 // 押注金額
	GameRound  int   // 已進行局數
	XVariation int   // 目前變數加減值（X）
}

func (pData *GameDataDB) Copy(pOther *GameDataDB) {
	if pOther == nil {
		return
	}

	pData.AverageBet = pOther.AverageBet
	pData.Factor = pOther.Factor
	pData.BetCount = pOther.BetCount
	pData.BetMoney = pOther.BetMoney
	pData.GameRound = pOther.GameRound
	pData.XVariation = pOther.XVariation
}

func (pData *GameDataDB) LogString() string {
	strLog := fmt.Sprintf("AverageBet[%d], Factor[%d], BetCount[%d], BetMoney[%d], GameRound[%d], XVariation[%d]",
		pData.AverageBet, pData.Factor, pData.BetCount, pData.BetMoney, pData.GameRound, pData.XVariation)

	return strLog
}

// 遊戲輸贏結果
type SysPoolWinLose struct {
	i64PoolDifference int64      // 水線變化量
	GameDBData        GameDataDB // 遊戲輸贏結果 DB 儲存資料
}

// bClearPool［是否清除水池］
func (pWinLose *SysPoolWinLose) Reset(bClearPool bool) {
	if bClearPool {
		pWinLose.i64PoolDifference = 0
	}
	pWinLose.GameDBData.BetMoney = 0
	pWinLose.GameDBData.BetCount = 0
}

//---------------------------------------------------------------------------
var g_i64WinLine = int64(0)				// 贏錢線
var g_iRoundLimit = DefaultRoundLimit	// 判定局數
var g_i64AverageBet = DefaultAverageBet	// 平均押注
var g_iFactor = DefaultFactor			// 係數定值

var g_map_pDealerSupplement_Base map[int32]*DealerSupplement	// 基準莊家補牌表
var g_map_pDealerSupplement_Now  map[int32]*DealerSupplement 	// 目前莊家補牌表
var g_ary_pSysPoolSetting []*PoolDynamicSetting					// 動態水線區間表

var g_ary_i32DealerPoints []int32

var g_bGetSysPoolConfig = false			// 是否已取得水池 config

var g_bGetWinLoseData = false				// 是否已自 DB 取得遊戲輸贏統計資料
var g_pSysPoolWinLose *SysPoolWinLose		// 房間總輸贏結果

var g_bEnablePoolControl = false		// 是否開啟水線控制功能

//---------------------------------------------------------------------------
// 是否有莊家管制點數，i32Seat < 0 表示為莊家
func hasDealerControlPoint(i32Seat, i32Point int32) bool {
	if i32Seat >= 0 {
		// 玩家
		return false
	}

	for i := 0; i < len(g_ary_i32DealerPoints); i++ {
		if g_ary_i32DealerPoints[i] == i32Point {
			return true
		}
	}

	return false
}

// 水池判定
func checkSysPool(iDealType int, pGameWinLose *SysPoolWinLose) (string, string, bool) {
	if iDealType != JD21_DealType_Natural {
		// 若發牌方式不是自然發牌，則不檢查水池
		return "", "", false
	}

	if pGameWinLose == nil {
		return "", "", false
	}

	strCheckLog, strUpdateLog := "", ""
	bRoundReset := false

	if g_pSysPoolWinLose.GameDBData.GameRound >= g_iRoundLimit {
		bRoundReset = true
	}

	g_pSysPoolWinLose.i64PoolDifference += pGameWinLose.i64PoolDifference

	strCheckLog = fmt.Sprintf("第 %d 局，判定局数 %d，变数加减（X）%d，本局水池变化量 %s，水池总变化量 %s",
		g_pSysPoolWinLose.GameDBData.GameRound, g_iRoundLimit, g_pSysPoolWinLose.GameDBData.XVariation, Global.MoneyString(pGameWinLose.i64PoolDifference),
		Global.MoneyString(g_pSysPoolWinLose.i64PoolDifference))

	if bRoundReset {
		strUpdateLog = updateDealerSupplement(false)
	}

	return strCheckLog, strUpdateLog, false
}

// 取得变数加减值（X）
func getXVariation() int {
	iIdx, iMinusIdx := -1, -1
	iLen := len(g_ary_pSysPoolSetting)

	mlog.Debug("getXVariation(), i64TotalPoolDifference[%d], iLen[%d]", g_pSysPoolWinLose.i64PoolDifference, iLen)
	SysPoolSettingData_Debug()
	// 輸贏區間 >= 0
	for i := 0; i < iLen; i++ {
		if g_ary_pSysPoolSetting[i].iConditionOdds < 0 {
			iMinusIdx = i
			break
		}

		i64Value1 := g_ary_pSysPoolSetting[i].i64AreaValue
		if i == 0 {
			if g_pSysPoolWinLose.i64PoolDifference >= i64Value1 {
				iIdx = i
				mlog.Debug("AreaNo[%d], iIdx[%d], i64TotalPoolDifference[%d], i64Value1[%d]",
					g_ary_pSysPoolSetting[iIdx].iAreaNo, iIdx, g_pSysPoolWinLose.i64PoolDifference, i64Value1)
				break
			}
		} else {
			i64Value2 := g_ary_pSysPoolSetting[i - 1].i64AreaValue
			if g_pSysPoolWinLose.i64PoolDifference < i64Value2 && g_pSysPoolWinLose.i64PoolDifference >= i64Value1 {
				iIdx = i
				mlog.Debug("AreaNo[%d], iIdx[%d], g_i64TotalPoolDifference[%d], i64Value1[%d], i64Value2[%d]",
					g_ary_pSysPoolSetting[iIdx].iAreaNo, iIdx, g_pSysPoolWinLose.i64PoolDifference, i64Value1, i64Value2)
				break
			}
		}
	}
	mlog.Debug("iMinusIdx[%d], iIdx[%d]", iMinusIdx, iIdx)

	// 輸贏區間 < 0
	if iMinusIdx > 0 {
		for i := iMinusIdx; i < iLen; i++ {
			i64Value1 := g_ary_pSysPoolSetting[i-1].i64AreaValue
			i64Value2 := g_ary_pSysPoolSetting[i].i64AreaValue
			mlog.Debug("i[%d], i64Value1[%d], i64Value2[%d]", i, i64Value1, i64Value2)
			if g_pSysPoolWinLose.i64PoolDifference > i64Value2 && g_pSysPoolWinLose.i64PoolDifference <= i64Value1 {
				iIdx = i - 1
				mlog.Debug("AreaNo[%d], iIdx[%d], i64TotalPoolDifference[%d], i64Value1[%d], i64Value2[%d]",
					g_ary_pSysPoolSetting[iIdx].iAreaNo, iIdx, g_pSysPoolWinLose.i64PoolDifference, i64Value1, i64Value2)
				break
			} else if i == iLen-1 && g_pSysPoolWinLose.i64PoolDifference <= i64Value2 {
				iIdx = i
				mlog.Debug("AreaNo[%d], iIdx[%d], g_i64TotalPoolDifference[%d], i64Value1[%d], i64Value2[%d]",
					g_ary_pSysPoolSetting[iIdx].iAreaNo, iIdx, g_pSysPoolWinLose.i64PoolDifference, i64Value1, i64Value2)
				break
			}
		}
	}

	if iIdx < 0 || iIdx >= iLen {
		return 0
	} else {
		return g_ary_pSysPoolSetting[iIdx].iXVariation
	}
}

// 更新目前莊家補牌表
func updateDealerSupplement(bInit bool) string {
	g_pSysPoolWinLose.GameDBData.XVariation = getXVariation()

	iInt := g_pSysPoolWinLose.GameDBData.XVariation / 1000
	iFraction := g_pSysPoolWinLose.GameDBData.XVariation % 1000
	strX := fmt.Sprintf("%d", iInt) + "." + fmt.Sprintf("%03d", iFraction)

	strLog := ""
	if bInit {
		strLog = fmt.Sprintf("初次设定，更新庄家补牌表，水池总变化量 %s，变数加减（X）更新为 %s", Global.MoneyString(g_pSysPoolWinLose.i64PoolDifference), strX)
	} else {
		strLog = fmt.Sprintf("进行局数已达判定局数 %d，更新庄家补牌表，变数加减（X）更新为 %s", g_iRoundLimit, strX)
	}

	if g_pSysPoolWinLose.GameDBData.XVariation == 0 {
		for i32Point, pData := range g_map_pDealerSupplement_Base {
			if pData != nil {
				g_map_pDealerSupplement_Now[i32Point].Copy(pData)
			}
		}
	} else {
		for i32Point, pData := range g_map_pDealerSupplement_Now {
			if pData != nil {
				for i := len(g_map_pDealerSupplement_Now[i32Point].ary_pData) - 1; i > 0; i-- {
					g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds =
						g_map_pDealerSupplement_Base[i32Point].ary_pData[i].iOdds + g_pSysPoolWinLose.GameDBData.XVariation
				}

				iRndRange := 0
				g_map_pDealerSupplement_Now[i32Point].iOddsSum = 0
				for i := 0; i < len(g_map_pDealerSupplement_Now[i32Point].ary_pData); i++ {
					iRndRange += g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds
					g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iRndRange = iRndRange
					g_map_pDealerSupplement_Now[i32Point].iOddsSum += g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds
				}
			}
		}
	}

	g_pSysPoolWinLose.GameDBData.GameRound = 0

	DealerSupplement_Now_Debug()

	return strLog
}

// 依指定點數取得牌組
func getCardSetByAssignedPoint(pCardBox *pokerset.PokerCardSet, pDealerCard *CardSetInfo, i32TargetPoint int32) *pokerset.PokerCardSet {
	if pCardBox == nil || pDealerCard == nil {
		return new(pokerset.PokerCardSet)
	}

	if pCardBox.CardCount() <= 0 {
		return new(pokerset.PokerCardSet)
	}

	if i32TargetPoint <= pDealerCard.MaxPoint() {
		return new(pokerset.PokerCardSet)
	}

	pRetSet := new(pokerset.PokerCardSet)
	pBaseSet := NewCardSetInfo(0)
	pBaseSet.AddCardSet(pDealerCard.pCardSet)

	mlog.Debug("getCardSetByAssignedPoint(), TargetPoint[%d]", i32TargetPoint)
	mlog.Debug("牌盒[%d][%s]", pCardBox.CardCount(), pCardBox.CodeString())
	mlog.Debug("莊家牌[%s](%d, %d), MaxPoint(%d)", pDealerCard.pCardSet.CodeString(),
		pDealerCard.i32Point_1, pDealerCard.i32Point_2, pDealerCard.MaxPoint())

	for i := 0; i < pCardBox.CardCount(); i++ {
		pTempSet := NewCardSetInfo(0)
		pTempSet.AddCardSet(pBaseSet.pCardSet)
		pCard := pCardBox.GetSingleCard(i,false)
		pTempSet.AddSingleCard(pCard)

		i32MaxPoint := pTempSet.MaxPoint()
		if i32MaxPoint < JD21_Soft {
			pBaseSet.AddSingleCard(pCard)
			pRetSet.AddSingleCard(pCard.CodeString())

			mlog.Debug("Idx[%d], pBaseSet[%s](%d, %d), MaxPoint(%d), pRetSet[%s]", i, pBaseSet.pCardSet.CodeString(),
				pBaseSet.i32Point_1, pBaseSet.i32Point_2, pBaseSet.MaxPoint(), pRetSet.CodeString())
		} else {
			if i32MaxPoint == i32TargetPoint {
				pBaseSet.AddSingleCard(pCard)
				pRetSet.AddSingleCard(pCard.CodeString())

				mlog.Debug("Idx[%d], pBaseSet[%s](%d, %d), MaxPoint(%d), pRetSet[%s]", i, pBaseSet.pCardSet.CodeString(),
					pBaseSet.i32Point_1, pBaseSet.i32Point_2, pBaseSet.MaxPoint(), pRetSet.CodeString())

				pCardBox.RemoveCardSet(pRetSet)
				mlog.Debug("剩餘牌盒[%d][%s]", pCardBox.CardCount(), pCardBox.CodeString())
				break
			}
		}
	}

	return pRetSet
}
//---------------------------------------------------------------------------
// debug 用
func SysPoolSettingData_Debug() {
	//mlog.Error("g_ary_pSysPoolSetting()")
	//for i := 0; i < len(g_ary_pSysPoolSetting); i++ {
	//	mlog.Warn("AreaNo[%d], iXVariation[%d], iConditionOdds[%d], i64AreaValue[%d]", g_ary_pSysPoolSetting[i].iAreaNo,
	//		g_ary_pSysPoolSetting[i].iXVariation, g_ary_pSysPoolSetting[i].iConditionOdds, g_ary_pSysPoolSetting[i].i64AreaValue)
	//}
}

func DealerSupplement_Base_Debug() {
	//mlog.Error("g_map_pDealerSupplement_Base")
	//for i32Point, pData := range g_map_pDealerSupplement_Base {
	//	if pData != nil {
	//		mlog.Error("Point[%d], iOddsSum[%d]", i32Point, pData.iOddsSum)
	//		for i := 0; i < len(pData.ary_pData); i++ {
	//			mlog.Warn("iOdds[%d], i32Points[%d], iRndRange[%d]",
	//				pData.ary_pData[i].iOdds, pData.ary_pData[i].i32Points, pData.ary_pData[i].iRndRange)
	//		}
	//	}
	//}
}

func DealerSupplement_Now_Debug() {
	//mlog.Error("g_map_pDealerSupplement_Now")
	//for i32Point, pData := range g_map_pDealerSupplement_Now {
	//	if pData != nil {
	//		mlog.Error("Point[%d], iOddsSum[%d]", i32Point, pData.iOddsSum)
	//		for i := 0; i < len(pData.ary_pData); i++ {
	//			mlog.Warn("iOdds[%d], i32Points[%d], iRndRange[%d]",
	//				pData.ary_pData[i].iOdds, pData.ary_pData[i].i32Points, pData.ary_pData[i].iRndRange)
	//		}
	//	}
	//}
}
