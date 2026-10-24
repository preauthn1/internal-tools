package jdesydLogic

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

var g_pTestCardBox *pokerset.PokerCardSet	// 牌盒
var g_pTestDealerCard *CardSetInfo	// 莊家牌

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
		ary_strDPoints := strings.Split(str, ",")
		for i := 0; i < len(ary_strDPoints); i++ {
			i64Point, _ := strconv.ParseInt(ary_strDPoints[i], 10, 32)
			g_ary_i32DealerPoints = append(g_ary_i32DealerPoints, int32(i64Point))

			strOdds := poolCfgs.GetConfigValue(ary_strDPoints[i], "Odds")
			ary_strOdds := strings.Split(strOdds, ",")
			strPoints := poolCfgs.GetConfigValue(ary_strDPoints[i],"Points")
			ary_strPoints := strings.Split(strPoints,",")
			if len(ary_strOdds) != len(ary_strPoints) {
				mlog.Error("庄家补牌表点数[%d]设定参数数目不相同。Odds 数目：%d，Points 数目：%d", i64Point,
					len(ary_strOdds), len(ary_strPoints))
			}

			pSupplement := new(DealerSupplement)
			iRndRange := 0
			iOddsSum := 0
			for j := 0; j < len(ary_strOdds); j++ {
				pData := new(SupplementData)
				i64Odds, _ := strconv.ParseInt(ary_strOdds[j],10,32)
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
					i64Points, _ := strconv.ParseInt(ary_strPoints[j],10,32)
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
		ary_strX := strings.Split(str, ",")
		str = poolCfgs.GetConfigValue("DynamicSetting", "ConditionOdds")
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

// 水線自動測試
func TestJDESYDSysPool(t *testing.T) {
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
	ary_pTestResult := make([]*PoolTestData, 0)
	g_pSysPoolWinLose = new(SysPoolWinLose)
	g_pTestCardBox = pokerset.NewPokerCardSet()
	g_pTestDealerCard = NewCardSetInfo(0)

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

	strResultFilename := fmt.Sprintf("經典21點水池自動測試[%d][%d][%d]", iXVariation, iConditionOdds, iTestRound)

	updateSupplement(iXVariation)

	iBustedCount, iAssignedPointCount, iNaturalDealingCount := 0, 0, 0
	for {
		if g_pTestCardBox.CardCount() <= iCutCardIndex {
			mlog.Warn("TestJDESYDSysPool(Shuffle), iCutCardIndex[%d], CardCount[%d]", iCutCardIndex, g_pTestCardBox.CardCount())
			g_pTestCardBox.NewSet(JD21_PokerSetCount, 0)
			g_pTestCardBox.Shuffle()

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
		pResult.i32DealerFinalPoint = g_pTestDealerCard.MaxPoint()
		if g_pTestDealerCard.bIsBusted {
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

		ary_pTestResult = append(ary_pTestResult, pResult)

		if iRound >= iTestRound {
			break
		}
	}

	pResult := new(PoolTestData)
	pResult.iTestNo = -1
	pResult.iBustedCount = iBustedCount
	pResult.iAssignedPointCount = iAssignedPointCount
	pResult.iNaturalDealingCount = iNaturalDealingCount

	ary_pTestResult = append(ary_pTestResult, pResult)

	WriteCSV_PoolTesting(ary_pTestResult, strPath, strResultFilename)
}

func Exists(strName string) bool {
	if _, err := os.Stat(strName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func WriteCSV_PoolTesting(ary_pTestResult []*PoolTestData, strPath, strFilename string)  {
	if len(ary_pTestResult) <= 0 || strFilename == "" {
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

	var ary_strResult [][]string
	// 設定標題
	ary_strCaption := []string{
		"測試編號",
		"變數加減(X)",
		"最初兩張牌點數",
		"莊家補到點數",
		"指定點數",
		"自然發牌",
	}

	ary_strResult = append(ary_strResult, ary_strCaption)

	// 填入測試結果資料
	for i := 0; i < len(ary_pTestResult); i++ {
		if ary_pTestResult[i] != nil {
			var ary_strData []string
			str := ""

			if ary_pTestResult[i].iTestNo > 0 {
				ary_strData = append(ary_strData, fmt.Sprintf("%d", ary_pTestResult[i].iTestNo)) // 測試編號
				ary_strData = append(ary_strData, fmt.Sprintf("%d", ary_pTestResult[i].iXvariation)) // 變數加減（X）
				ary_strData = append(ary_strData, fmt.Sprintf("%d", ary_pTestResult[i].i32FirstPoint))	// 最初兩張牌的點數
				// 莊家補到點數
				if ary_pTestResult[i].bBusted {
					str = fmt.Sprintf("爆牌 %d", ary_pTestResult[i].i32DealerFinalPoint)
				} else {
					str = fmt.Sprintf("%d", ary_pTestResult[i].i32DealerFinalPoint)
				}
				ary_strData = append(ary_strData, fmt.Sprintf("%s", str))
				if ary_pTestResult[i].bNaturalDealing {
					ary_strData = append(ary_strData, "")  // 指定點數
					ary_strData = append(ary_strData, "V") // 自然發牌
				} else {
					ary_strData = append(ary_strData, "V") // 指定點數
					ary_strData = append(ary_strData, "")  // 自然發牌
				}
			} else {
				ary_strData = append(ary_strData, "合計")
				ary_strData = append(ary_strData, "")
				ary_strData = append(ary_strData, "")
				ary_strData = append(ary_strData, fmt.Sprintf("爆牌數 %d", ary_pTestResult[i].iBustedCount))
				ary_strData = append(ary_strData, fmt.Sprintf("%d", ary_pTestResult[i].iAssignedPointCount))
				ary_strData = append(ary_strData, fmt.Sprintf("%d", ary_pTestResult[i].iNaturalDealingCount))
			}

			ary_strResult = append(ary_strResult, ary_strData)
		}
	}

	// 寫成 csv 檔
	pResultFile.WriteString("\xEF\xBB\xBF")
	pCsvWriter := csv.NewWriter(pResultFile)
	for i := 0; i < len(ary_strResult); i++ {
		err := pCsvWriter.Write(ary_strResult[i])
		if err != nil {
			mlog.Error("csv 檔寫入失敗：%v", err.Error())
			panic(err)
		}
	}
	pCsvWriter.Flush()
}

// 發牌
func testDeal() (bool, bool, int32) {
	g_pTestDealerCard.ClearData()
	pSet := g_pTestCardBox.Deal(2)
	g_pTestDealerCard.AddCardSet(pSet)

	bNaturalDealing, bControlPoint := true, false
	i32Point := g_pTestDealerCard.MaxPoint()
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
			pCardSet := getCardSetByAssignedPoint(g_pTestCardBox, g_pTestDealerCard, i32TargetPoint)
			g_pTestDealerCard.AddCardSet(pCardSet)
			mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), g_pTestDealerCard.pCardSet.String(),
				g_pTestDealerCard.i32Point_1, g_pTestDealerCard.i32Point_2)
		} else {
			// 自然發牌
			if !g_pTestDealerCard.bIsBlackJack && !g_pTestDealerCard.bIsBusted && !g_pTestDealerCard.IsSoftReached() {
				for {
					pCardSet := g_pTestCardBox.Deal(1)
					g_pTestDealerCard.AddCardSet(pCardSet)
					mlog.Debug("莊家補牌：%s, 莊家牌：%s (%d, %d )", pCardSet.String(), g_pTestDealerCard.pCardSet.String(),
						g_pTestDealerCard.i32Point_1, g_pTestDealerCard.i32Point_2)
					if g_pTestDealerCard.IsSoftReached() {
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
