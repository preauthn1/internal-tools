package dchjk

import (
	"encoding/csv"
	"flag"
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

var gpTestCardBox *pokerset.PokerCardSet	// 牌盒
var gpTestDealerCard *CardSetInfo	// 莊家牌

// 水池測試結果資料
type PoolTestData struct {
	iTestNo              int   // 測試編號
	iXvariation          int   // 變數加減（X）
	i32DealerFinalPoint  int32 // 莊家補到點數
	bBusted              bool  // 是否爆牌
	bNaturalDealing      bool  // 是否為自然發牌
	iBustedCount         int   // 爆牌總次數
	iAssignedPointCount  int   // 指定點數次數
	iNaturalDealingCount int   // 自然發牌總次數
	i32FirstPoint        int32 // 最初兩張牌的點數
}

//---------------------------------------------------------------------------
// 讀取 config
func getSysPoolConfig() {
	g_map_pDealerSupplement_Base = make(map[int32]*DealerSupplement)
	g_map_pDealerSupplement_Now = make(map[int32]*DealerSupplement)
	g_ary_pSysPoolSetting = make([]*PoolDynamicSetting, 0)

	poolCfgs := cfg.GetConfig("gameserver")
	if poolCfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		g_i64WinLine = poolCfgs.GetConfigInt64("SystemPool", "WinLine")	// 贏錢線
		g_iRoundLimit = poolCfgs.GetConfigInt("SystemPool", "RoundLimit")	// 判定局數
		g_i64AverageBet = poolCfgs.GetConfigInt64("SystemPool", "AverageBet")	// 平均押注
		g_iFactor = poolCfgs.GetConfigInt("SystemPool", "Factor")	// 係數定值
		mlog.Debug("WinLine[%d], RoundLimit[%d], AverageBet[%d], Factor[%d]",
			g_i64WinLine, g_iRoundLimit, g_i64AverageBet, g_iFactor)

		// 取得預設莊家補牌表
		str := poolCfgs.GetConfigValue("DealerPoint", "Points")
		arystrDPoints := strings.Split(str, ",")
		for i := 0; i < len(arystrDPoints); i++ {
			i64Point, _ := strconv.ParseInt(arystrDPoints[i], 10, 32)
			g_ary_i32DealerPoints = append(g_ary_i32DealerPoints, int32(i64Point))

			strOdds := poolCfgs.GetConfigValue(arystrDPoints[i], "Odds")
			arystrOdds := strings.Split(strOdds, ",")
			strPoints := poolCfgs.GetConfigValue(arystrDPoints[i],"Points")
			arystrPoints := strings.Split(strPoints,",")
			if len(arystrOdds) != len(arystrPoints) {
				mlog.Error("庄家补牌表点数[%d]设定参数数目不相同。Odds 数目：%d，Points 数目：%d", i64Point,
					len(arystrOdds), len(arystrPoints))
			}

			pSupplement := new(DealerSupplement)
			iRndRange := 0
			iOddsSum := 0
			for j := 0; j < len(arystrOdds); j++ {
				pData := new(SupplementData)
				i64Odds, _ := strconv.ParseInt(arystrOdds[j],10,32)
				pData.iOdds = int(i64Odds)
				iOddsSum += pData.iOdds

				if pData.iOdds <= 0 {
					pData.iRndRange = 0
				} else {
					iRndRange += pData.iOdds
					pData.iRndRange = iRndRange
				}

				if j >= len(arystrPoints) {
					pData.i32Points = 0
				} else {
					i64Points, _ := strconv.ParseInt(arystrPoints[j],10,32)
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
		str = poolCfgs.GetConfigValue("DynamicSetting", "XVariation")
		arystrX := strings.Split(str, ",")
		str = poolCfgs.GetConfigValue("DynamicSetting", "ConditionOdds")
		arystrConditionOdds := strings.Split(str, ",")
		if len(arystrX) != len(arystrConditionOdds) {
			mlog.Error("动态水线区间表 XVariation 与 ConditionOdds 参数数目不相同。XVariation 数目：%d，ConditionOdds 数目：%d",
				len(arystrX), len(arystrConditionOdds))
		}

		for i := 0; i < len(arystrX); i++ {
			pData := new(PoolDynamicSetting)
			pData.iAreaNo = i + 1
			i64Value, _ := strconv.ParseInt(arystrX[i], 10, 32)
			pData.iXVariation = int(i64Value)
			if i >= len(arystrConditionOdds) {
				pData.iConditionOdds = 0
			} else {
				i64Value, _ := strconv.ParseInt(arystrConditionOdds[i], 10, 32)
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

// 水線自動測試
func TestDCHJKSysPool(t *testing.T) {
	//建立測試結果存放資料夾
	strPath := "./TestResult" // 測試結果存放路徑
	bExist := Exists(strPath)

	if !bExist {
		err := os.Mkdir(strPath, os.ModePerm)
		if err != nil {
			mlog.Error("建立資料夾 %s 失敗", strPath)
		}
	}

	getSysPoolConfig()

	iXVariationBase, iConditionOdds, iTestRound := 50, 0, 1000
	iXVariation, iRound, iCutCardIndex := 0, 0, JD21_MinCutCardIndex
	arypTestResult := make([]*PoolTestData, 0)
	g_pSysPoolWinLose = new(SysPoolWinLose)
	gpTestCardBox = pokerset.NewPokerCardSet()
	gpTestDealerCard = NewCardSetInfo(0)

	iArgLen := len(flag.Args())
	if iArgLen >= 1 {
		// 變數加減（X）
		i64Value, _ := strconv.ParseInt(flag.Arg(0), 10, 32)
		iXVariationBase = int(i64Value)
	}

	if iArgLen >= 2 {
		// 輸贏區間
		i64Value, _ := strconv.ParseInt(flag.Arg(1), 10, 32)
		iConditionOdds = int(i64Value)
	}

	iXVariation = iXVariationBase * iConditionOdds
	if iXVariation < 0 {
		iXVariation = -iXVariation
	}

	if iArgLen >= 3 {
		// 執行局數
		i64Value, _ := strconv.ParseInt(flag.Arg(2), 10, 32)
		iTestRound = int(i64Value)
	}

	strResultFilename := fmt.Sprintf("賭場黑傑克水線自動測試[%d][%d][%d]", iXVariation, iConditionOdds, iTestRound)

	updateSupplement(iXVariation)

	iBustedCount, iAssignedPointCount, iNaturalDealingCount := 0, 0, 0
	for {
		if gpTestCardBox.CardCount() <= iCutCardIndex {
			mlog.Warn("TestJDESYDSysPool(Shuffle), iCutCardIndex[%d], CardCount[%d]", iCutCardIndex, gpTestCardBox.CardCount())
			gpTestCardBox.NewSet(JD21_PokerSetCount, 0)
			gpTestCardBox.Shuffle()

			iCutCardIndex = util.RndInt(JD21_MinCutCardIndex, JD21_MaxCutCardIndex)
		}

		bNaturalDealing, bControlPoint, i32Point := testDeal()

		if !bControlPoint {
			continue
		}

		iRound++

		pResult := new(PoolTestData)
		pResult.iTestNo = iRound
		pResult.iXvariation = iXVariation
		pResult.i32DealerFinalPoint = gpTestDealerCard.MaxPoint()
		if gpTestDealerCard.bIsBusted {
			pResult.bBusted = true
			iBustedCount++
		} else {
			pResult.bBusted = false
		}
		pResult.bNaturalDealing = bNaturalDealing
		if pResult.bNaturalDealing {
			iNaturalDealingCount++
		} else {
			iAssignedPointCount++
		}
		pResult.i32FirstPoint = i32Point

		arypTestResult = append(arypTestResult, pResult)

		if iRound >= iTestRound {
			break
		}
	}

	pResult := new(PoolTestData)
	pResult.iTestNo = -1
	pResult.iBustedCount = iBustedCount
	pResult.iAssignedPointCount = iAssignedPointCount
	pResult.iNaturalDealingCount = iNaturalDealingCount

	arypTestResult = append(arypTestResult, pResult)

	writeCSVPoolTesting(arypTestResult, strPath, strResultFilename)
}

func Exists(strName string) bool {
	if _, err := os.Stat(strName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func writeCSVPoolTesting(arypTestResult []*PoolTestData, strPath, strFilename string)  {
	if len(arypTestResult) <= 0 || strFilename == "" {
		return
	}

	tCurrent := time.Now()
	strTime := fmt.Sprintf("_%04d%02d%02d_%02d%02d%02d",
		tCurrent.Year(), tCurrent.Month(), tCurrent.Day(), tCurrent.Hour(), tCurrent.Minute(), tCurrent.Second())
	strResultFile := strPath + "/" + strFilename + strTime + ".csv"
	if Exists(strResultFile) {
		err := os.Remove(strResultFile)
		if err != nil {
			mlog.Error("Delete %s Error", strResultFile)
			return
		}
	}

	pResultFile, err := os.OpenFile(strResultFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer pResultFile.Close()

	var arystrResult [][]string
	// 設定標題
	arystrCaption := []string{
		"測試編號",
		"變數加減(X)",
		"最初兩張牌點數",
		"莊家補到點數",
		"指定點數",
		"自然發牌",
	}

	arystrResult = append(arystrResult, arystrCaption)

	// 填入測試結果資料
	for i := 0; i < len(arypTestResult); i++ {
		if arypTestResult[i] != nil {
			var arystrData []string
			str := ""

			if arypTestResult[i].iTestNo > 0 {
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iTestNo)) // 測試編號
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iXvariation)) // 變數加減（X）
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32FirstPoint))	// 最初兩張牌的點數
				// 莊家補到點數
				if arypTestResult[i].bBusted {
					str = fmt.Sprintf("爆牌 %d", arypTestResult[i].i32DealerFinalPoint)
				} else {
					str = fmt.Sprintf("%d", arypTestResult[i].i32DealerFinalPoint)
				}
				arystrData = append(arystrData, fmt.Sprintf("%s", str))
				if arypTestResult[i].bNaturalDealing {
					arystrData = append(arystrData, "")  // 指定點數
					arystrData = append(arystrData, "V") // 自然發牌
				} else {
					arystrData = append(arystrData, "V") // 指定點數
					arystrData = append(arystrData, "")  // 自然發牌
				}
			} else {
				arystrData = append(arystrData, "合計")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, fmt.Sprintf("爆牌數 %d", arypTestResult[i].iBustedCount))
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iAssignedPointCount))
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iNaturalDealingCount))
			}

			arystrResult = append(arystrResult, arystrData)
		}
	}

	// 寫成 csv 檔
	pResultFile.WriteString("\xEF\xBB\xBF")
	pCsvWriter := csv.NewWriter(pResultFile)
	for i := 0; i < len(arystrResult); i++ {
		err := pCsvWriter.Write(arystrResult[i])
		if err != nil {
			mlog.Error("csv 檔寫入失敗：%v", err.Error())
			panic(err)
		}
	}
	pCsvWriter.Flush()
}

// 發牌
func testDeal() (bool, bool, int32) {
	gpTestDealerCard.ClearData()
	pSet := gpTestCardBox.Deal(2)
	gpTestDealerCard.AddCardSet(pSet)

	bNaturalDealing, bControlPoint := true, false
	i32Point := gpTestDealerCard.MaxPoint()
	i32TargetPoint := int32(0)

	// 檢核補牌表
	if hasDealerControlPoint(-1, i32Point) {
		bControlPoint = true

		iRnd := util.RndInt(0, g_map_pDealerSupplement_Now[i32Point].iOddsSum)
		for i := 0; i < len(g_map_pDealerSupplement_Now[i32Point].ary_pData); i++ {
			iRange1 := g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iRndRange
			if i == 0 {
				if iRnd <= iRange1 {
					// 自然發牌
					break
				}
			} else {
				iRange2 := g_map_pDealerSupplement_Now[i32Point].ary_pData[i-1].iRndRange
				if iRnd > iRange2 && iRnd <= iRange1 {
					i32TargetPoint = g_map_pDealerSupplement_Now[i32Point].ary_pData[i].i32Points
					break
				}
			}
		}

		mlog.Warn("Rnd[%d], TargetPoint[%d]", iRnd, i32TargetPoint)

		if i32TargetPoint > 0 {
			bNaturalDealing = false
		}

		if !bNaturalDealing {
			// 依莊家補牌表發牌
			pCardSet := getCardSetByAssignedPoint(gpTestCardBox, gpTestDealerCard, i32TargetPoint)
			gpTestDealerCard.AddCardSet(pCardSet)
			mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), gpTestDealerCard.pCardSet.String(),
				gpTestDealerCard.i32Point_1, gpTestDealerCard.i32Point_2)
		} else {
			// 自然發牌
			if !gpTestDealerCard.bIsBlackJack && !gpTestDealerCard.bIsBusted && !gpTestDealerCard.IsSoftReached() {
				for {
					pCardSet := gpTestCardBox.Deal(1)
					gpTestDealerCard.AddCardSet(pCardSet)
					mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), gpTestDealerCard.pCardSet.String(),
						gpTestDealerCard.i32Point_1, gpTestDealerCard.i32Point_2)
					if gpTestDealerCard.IsSoftReached() {
						break
					}
				}
			}
		}
	}

	return bNaturalDealing, bControlPoint, i32Point
}

func updateSupplement(iXvariation int) {
	iX := 0
	if iXvariation > 0 {
		iX = iXvariation
	}

	if iX == 0 {
		for i32point, pData := range g_map_pDealerSupplement_Base {
			if pData != nil {
				g_map_pDealerSupplement_Now[i32point].Copy(pData)
			}
		}
	} else {
		for i32Point, pData := range g_map_pDealerSupplement_Base {
			if pData != nil {
				for i := 1; i < len(pData.ary_pData); i++ {
					g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds =
						g_map_pDealerSupplement_Base[i32Point].ary_pData[i].iOdds + iX
				}
			}

			iRndRange := 0
			g_map_pDealerSupplement_Now[i32Point].iOddsSum = 0
			for i := 0; i < len(g_map_pDealerSupplement_Now[i32Point].ary_pData); i++ {
				iRndRange += g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds
				g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iRndRange = iRndRange
				g_map_pDealerSupplement_Now[i32Point].iOddsSum += g_map_pDealerSupplement_Now[i32Point].ary_pData[i].iOdds
			}
		}

		mlog.Error("g_map_pDealerSupplement_Now, iXvariation[%d], iX[%d]", iXvariation, iX)
		for i32Point, pData := range g_map_pDealerSupplement_Now {
			if pData != nil {
				mlog.Error("Point[%d], iOddsSum[%d]", i32Point, pData.iOddsSum)
				for i := 0; i < len(pData.ary_pData); i++ {
					mlog.Warn("iOdds[%d], i32Points[%d], iRndRange[%d]",
						pData.ary_pData[i].iOdds, pData.ary_pData[i].i32Points, pData.ary_pData[i].iRndRange)
				}
			}
		}
	}
}
