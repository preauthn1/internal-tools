package jlbpk

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	CSVSum                 int32 = -1 // 合計
	CSVPatternCountCaption int32 = -2 // 牌型統計標題
	CSVPatternCount        int32 = -3 // 牌型統計
	CSVSumPatternCount     int32 = -4 // 牌型統計合計
)

// 水池測試結果資料
type PoolTestData struct {
	i32GameRound           int32    // 測試局數
	iWeightTableNo         int      // 權重表
	i32DealerWinCount      int32    // 莊家勝利局數
	i32DealerLoseCount     int32    // 玩家勝利局數
	i32DealerNotOkCount    int32    // 莊家不成局局數
	i64PoolAmount          int64    // 水線
	f32GameRoundRatio      float32  // 場次比值
	strDealerHandCards     string   // 莊家手牌
	arystrUserHandCards    []string // 玩家手牌
	iUserCount             int      // 玩家人數
	strPattern             string   // 牌型
	aryiPatternCountByUser []int    // 各牌型出牌數
}

// 牌型出牌數
type PatternCount struct {
	aryiPatternCount []int // 各牌型出牌數
	iTotalCount      int   // 總數
}

func NewPatternCount() *PatternCount {
	pPattern := new(PatternCount)
	for i := pokerset.POKERPatternNone; i <= pokerset.POKERPatternRoyalFlush; i++ {
		pPattern.aryiPatternCount = append(pPattern.aryiPatternCount, 0)
	}
	pPattern.iTotalCount = 0

	return pPattern
}

func (pPattern *PatternCount) AddPattern(i32Pattern int32) {
	if i32Pattern < pokerset.POKERPatternNone || i32Pattern >= int32(len(pPattern.aryiPatternCount)) {
		return
	}

	pPattern.aryiPatternCount[i32Pattern]++
	pPattern.iTotalCount++
}

var gmappPatternCount map[int]*PatternCount // map[座位]牌型統計，座位 = -1 為莊家，>=0 為玩家
var gbGetSystemPoolConfig = false

func getSystemPoolConfig() {
	g_map_pDealerWeightTable = make(map[int]*DealerWeightTable)
	poolcfgs := cfg.GetConfig("gameserver")
	if poolcfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		g_i64WinLine = poolcfgs.GetConfigInt64("SystemPool", "WinLine")                  // 贏錢線
		g_i32PlayerWeightDefault = poolcfgs.GetConfigInt32("SystemPool", "PlayerWeight") // 預設玩家拿牌權重
		g_i32DealerWeightDefault = poolcfgs.GetConfigInt32("SystemPool", "DealerWeight") // 預設莊家權重
		g_iCheatRate = poolcfgs.GetConfigInt("SystemPool", "CheatRate")                  // 換牌率
		g_iRoundLimit = poolcfgs.GetConfigInt("SystemPool", "RoundLimit")                // 判定局數

		mlog.Warn("g_i64WinLine[%d], g_i32PlayerWeightDefault[%d], g_i32DealerWeightDefault[%d], g_iCheatRate[%d], g_iRoundLimit[%d]",
			g_i64WinLine, g_i32PlayerWeightDefault, g_i32DealerWeightDefault, g_iCheatRate, g_iRoundLimit)

		// 機率表
		str := ""
		iTableCount := poolcfgs.GetConfigInt("WeightTable", "TableCount") // 權重表張數
		for iNo := 0; iNo < iTableCount; iNo++ {
			if _, ok := g_map_pDealerWeightTable[iNo]; !ok {
				g_map_pDealerWeightTable[iNo] = new(DealerWeightTable)
			}

			strName := fmt.Sprintf("%d", iNo)
			str = poolcfgs.GetConfigValue(strName, "Pattern")
			arystrPattern := strings.Split(str, ",")
			str = poolcfgs.GetConfigValue(strName, "Weight")
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

	mlog.Error("setDealerWeightTable()")
	mlog.Error("g_map_pDealerWeightTable")
	for iTableNo, pTable := range g_map_pDealerWeightTable {
		if pTable != nil {
			mlog.Error("Table[%d], i32TotalWeight[%d]", iTableNo, pTable.i32TotalWeight)
			for i := 0; i < len(pTable.ary_pWeightData); i++ {
				str := ""
				str += fmt.Sprintf("CardPattern[%s(%d)] ",
					pokerset.PatternString(pTable.ary_pWeightData[i].i32CardPattern), pTable.ary_pWeightData[i].i32CardPattern)
				str += fmt.Sprintf("Weight[%d] ", pTable.ary_pWeightData[i].i32Weight)
				str += fmt.Sprintf("RndRangeMin[%d] ", pTable.ary_pWeightData[i].i32RndRangeMin)
				str += fmt.Sprintf("RndRangeMax[%d] ", pTable.ary_pWeightData[i].i32RndRangeMax)
				mlog.Warn("%s", str)
			}
		}
	}
}

func getWeightTable(iTableNo int) *DealerWeightTable {
	if iTableNo < DealerWeightNo_A || iTableNo >= DealerWeightNo_Count {
		return new(DealerWeightTable)
	}

	if !gbGetSystemPoolConfig {
		getSystemPoolConfig()
		gbGetSystemPoolConfig = true
	}

	return g_map_pDealerWeightTable[iTableNo]
}

// 發莊家牌
func TestDealDealerHandCard(t *testing.T) {
	iTableNo := util.RndInt(0, DealerWeightNo_Count)
	mlog.Debug("莊家權重表編號：%d", iTableNo)
	pTable := getWeightTable(iTableNo)
	for i := 0; i < len(pTable.ary_pWeightData); i++ {
		example.invalid("Pattern[%s(%d)], RndRangeMin[%d], RndRangeMax[%d]", pokerset.PatternString(pTable.ary_pWeightData[i].i32CardPattern),
			pTable.ary_pWeightData[i].i32CardPattern, pTable.ary_pWeightData[i].i32RndRangeMin, pTable.ary_pWeightData[i].i32RndRangeMax)
	}
	mlog.Debug("-------------------------------------------------------------")

	pCardBox := NewJlbpkPokerSet()
	pCardBox.NewSet(JLBPK_PokerSetCount, 0)
	pCardBox.Shuffle()

	iCount, iTestCount := 0, 100

	for {
		mlog.Debug("CardBox, Count[%d], %s", pCardBox.CardCount(), pCardBox.CodeString())
		pCards := dealDealerHandCard(pCardBox, JLBPK_HandCardCount, pTable)
		mlog.Warn("pCards: %s", pCards.CodeString())
		mlog.Debug("-------------------------------------------------------------")
		iCount++

		if pCardBox.CardCount() < JLBPK_HandCardCount || iCount >= iTestCount {
			break
		}
	}
}

// 特殊賠率表
func specialSheet() map[int32]int64 {
	mapi64Sheet := make(map[int32]int64)

	aryi64Pay := []int64{100, 50, 20, 7, 5, 4, 3, 2, 1, 1}
	for i := pokerset.POKERPatternRoyalFlush; i >= pokerset.POKERPatternHighCard; i-- {
		mapi64Sheet[i] = aryi64Pay[pokerset.POKERPatternRoyalFlush-i]
	}

	for i32Pattern, i64Pay := range mapi64Sheet {
		mlog.Debug("Pattern[%s(%d)], Pay[%d]", pokerset.PatternString(i32Pattern), i32Pattern, i64Pay)
	}

	return mapi64Sheet
}

// 檢查水池
//func TestCheckSysPool(t *testing.T) {
//	iRoundLimit, iTestRound, iWeightNoNow, iUserCount := 1000, 1, DealerWeightNo_B, 1
//	i64WinLine, i64BetBase := int64(1500000), int64(10)
//	pRoomWinLose, pGameWinLose := new(SysPoolWinLose), new(SysPoolWinLose)
//	pRoomWinLose.i64PoolAmount = 1500000
//	pTable := getWeightTable(iWeightNoNow)
//	bDealerIsGood := true
//	mapi64SpecialSheet := specialSheet()
//
//	pCardBox := NewJlbpkPokerSet()
//	pDealerHandCards := NewJlbpkPokerSet() // 莊家手牌
//	var arypHandCards []*JlbpkPokerSet     // 玩家手牌
//	var arypUser []*JlbpkUser
//	for i := 0; i < iUserCount; i++ {
//		pCards := NewJlbpkPokerSet()
//		arypHandCards = append(arypHandCards, pCards)
//		pUser := NewJlbpkUser()
//		arypUser = append(arypUser, pUser)
//	}
//
//	for i := 0; i < iTestRound; i++ {
//		pCardBox.NewSet(JLBPK_PokerSetCount, 0)
//		pCardBox.Shuffle()
//		pGameWinLose.Reset(true)
//		pDealerHandCards.SetCardSet("")
//		for iSeat := 0; iSeat < iUserCount; iSeat++ {
//			arypHandCards[iSeat].SetCardSet("")
//			arypUser[iSeat].ResetGameData()
//		}
//
//		// 發牌
//		str := "【玩家牌】"
//		for iSeat := 0; iSeat < iUserCount; iSeat++ {
//			pSet := pCardBox.Deal(JLBPK_HandCardCount)
//			arypHandCards[iSeat].SetCardSet(pSet.CodeString())
//			str += fmt.Sprintf("【User[%d]%s[%s(%d)]】", iSeat+1, arypHandCards[iSeat].String(),
//				pokerset.PatternString(arypHandCards[iSeat].Pattern(false)), arypHandCards[iSeat].Pattern(false))
//		}
//		mlog.Debug("%s", str)
//
//		str = "【庄家牌】"
//		pSet := dealDealerHandCard(pCardBox, JLBPK_HandCardCount, pTable)
//		pDealerHandCards.SetCardSet(pSet.CodeString())
//		str += fmt.Sprintf("%s[%s(%d)]", pDealerHandCards.String(),
//			pokerset.PatternString(pDealerHandCards.Pattern(false)), pDealerHandCards.Pattern(false))
//		if pDealerHandCards.Pattern(true) == pokerset.POKER_Pattern_None {
//			bDealerIsGood = false
//			str += "（不成局）"
//		} else {
//			bDealerIsGood = true
//		}
//		mlog.Debug("%s", str)
//
//		// 押注
//		for iSeat := 0; iSeat < iUserCount; iSeat++ {
//			arypUser[iSeat].i64BetBase = i64BetBase
//			arypUser[iSeat].i64BetDouble = i64BetBase * 2
//		}
//
//		// 底注結算
//		for iSeat := 0; iSeat < iUserCount; iSeat++ {
//			if !bDealerIsGood {
//				// 莊家不成局
//				pGameWinLose.bDealerNotOK = true
//				pGameWinLose.i64PoolAmount -= arypUser[iSeat].i64BetBase * int64(1+JLBPK_BetBasePayRate)
//			} else {
//				iCmpResult := pCardBox.Compare(pDealerHandCards.GetPokerCardSet(), arypHandCards[iSeat].GetPokerCardSet(), false, false)
//				i64BetBaseWinLose := int64(0)
//				switch iCmpResult {
//				case pokerset.POKER_Compare_Even: // 平手
//					i64BetBaseWinLose = arypUser[iSeat].i64BetBase
//				case pokerset.POKER_Compare_Dealer: // 莊家贏
//					i64BetBaseWinLose = -arypUser[iSeat].i64BetBase
//				case pokerset.POKER_Compare_Player: // 玩家贏
//					i64BetBaseWinLose = arypUser[iSeat].i64BetBase * int64(1+JLBPK_BetBasePayRate)
//				}
//
//				if iCmpResult != pokerset.POKER_Compare_Even {
//					pGameWinLose.i64PoolAmount -= i64BetBaseWinLose
//				}
//			}
//		}
//
//		// 加注結算
//		for iSeat := 0; iSeat < iUserCount; iSeat++ {
//			if !bDealerIsGood {
//				// 莊家不成局
//				pGameWinLose.i64PoolAmount -= arypUser[iSeat].i64BetDouble
//			} else {
//				iCmpResult := pCardBox.Compare(pDealerHandCards.GetPokerCardSet(), arypHandCards[iSeat].GetPokerCardSet(), false, false)
//				i64DoubleWinLose := int64(0)
//				switch iCmpResult {
//				case pokerset.POKER_Compare_Even: // 平手
//					i64DoubleWinLose = arypUser[iSeat].i64BetDouble
//				case pokerset.POKER_Compare_Dealer: // 第一組牌
//					i64DoubleWinLose = -arypUser[iSeat].i64BetDouble
//				case pokerset.POKER_Compare_Player: // 第二組牌
//					i64DoubleWinLose = arypUser[iSeat].i64BetDouble * (1 + mapi64SpecialSheet[arypHandCards[iSeat].Pattern(false)])
//				}
//
//				if iCmpResult != pokerset.POKER_Compare_Even {
//					pGameWinLose.i64PoolAmount -= i64DoubleWinLose
//				}
//			}
//		}
//
//		// 遊戲結束
//		if pGameWinLose.bDealerNotOK {
//			pGameWinLose.GameDBData.DealerNotOKCount++
//		} else {
//			if pGameWinLose.i64PoolAmount >= 0 {
//				pGameWinLose.GameDBData.DealerWinCount++
//			} else {
//				pGameWinLose.GameDBData.DealerLoseCount++
//			}
//		}
//
//		mlog.Warn("【本局】庄家胜利[%d]，玩家胜利[%d]，庄家不成局[%d]，水量变化[%s]",
//			pGameWinLose.GameDBData.DealerWinCount, pGameWinLose.GameDBData.DealerLoseCount, pGameWinLose.GameDBData.DealerNotOKCount,
//			Global.MoneyString(pGameWinLose.i64PoolAmount))
//
//		// 將本局輸贏統計至房間輸贏
//		pRoomWinLose.i64PoolAmount += pGameWinLose.i64PoolAmount
//		pRoomWinLose.GameDBData.DealerLoseCount += pGameWinLose.GameDBData.DealerLoseCount
//		pRoomWinLose.GameDBData.DealerNotOKCount += pGameWinLose.GameDBData.DealerNotOKCount
//		pRoomWinLose.GameDBData.DealerWinCount += pGameWinLose.GameDBData.DealerWinCount
//		pRoomWinLose.GameDBData.GameRound++
//
//		mlog.Warn("判定局数[%d]，赢钱线[%s]，权重表[%s]，【房间】游戏局数[%d]，庄家胜利[%d]，玩家胜利[%d]，庄家不成局[%d]，水线[%s]",
//			iRoundLimit, Global.MoneyString(i64WinLine), g_map_strDealerWeightNo[iWeightNoNow],
//			pRoomWinLose.GameDBData.GameRound, pRoomWinLose.GameDBData.DealerWinCount, pRoomWinLose.GameDBData.DealerLoseCount,
//			pRoomWinLose.GameDBData.DealerNotOKCount, Global.MoneyString(pRoomWinLose.i64PoolAmount))
//
//		iWeightNo, strLog := changeDealerWeightTable(pRoomWinLose, iRoundLimit, iWeightNoNow, i64WinLine)
//		if iWeightNo >= DealerWeightNo_A && iWeightNo < DealerWeightNo_Count {
//			iWeightNoNow = iWeightNo
//			pTable = getWeightTable(iWeightNoNow)
//		}
//		if strLog != "" {
//			mlog.Warn("%s", strLog)
//		}
//	}
//}

// 換牌
func TestChangingCards(t *testing.T) {
	i32DealerWeight := int32(100)
	iSeatCount := 7
	arypUserData := make([]*JlbpkUser, iSeatCount)
	for i := 0; i < iSeatCount; i++ {
		arypUserData[i] = NewJlbpkUser()
		arypUserData[i].i32SeatNo = int32(i)
		arypUserData[i].iUserState = JLB_UserState_Normal
		iRnd := util.RndInt(0, 2)
		if iRnd == 0 {
			arypUserData[i].i32Weight = int32(util.RndInt(80, 100))
		} else {
			arypUserData[i].i32Weight = 100
		}
	}

	pCardBox := NewJlbpkPokerSet()
	pCardBox.NewSet(JLBPK_PokerSetCount, 0)
	pCardBox.Shuffle()

	arypRetData := changingCards(pCardBox, i32DealerWeight, arypUserData)
	for i := 0; i < len(arypRetData); i++ {
		mlog.Debug("RetData: i32seatNo[%d], strCardCodeString[%s]",
			arypRetData[i].i32seatNo, arypRetData[i].strCardCodeString)
		pSet := NewJlbpkPokerSet()
		pSet.SetCardSet(arypRetData[i].strCardCodeString)
		if arypRetData[i].i32seatNo < 0 {
			mlog.Debug("牌組：%s [%s(%d)]", pSet.CodeString(), pSet.PatternString(true), pSet.Pattern(true))
		} else {
			mlog.Debug("牌組：%s [%s(%d)]", pSet.CodeString(), pSet.PatternString(false), pSet.Pattern(false))
		}
	}
}

// 水線自動測試
//func TestHDCSysPool(t *testing.T) {
//	//建立測試結果存放資料夾
//	strPath := "./TestResult" // 測試結果存放路徑
//	bExist := Exists(strPath)
//
//	if !bExist {
//		err := os.Mkdir(strPath, os.ModePerm)
//		if err != nil {
//			mlog.Error("建立資料夾 %s 失敗", strPath)
//		}
//	}
//
//	iWeightNoNow, iUserCount := DealerWeightNo_B, 1
//	i64BetBase := int64(100)
//
//	pRoomWinLose, pGameWinLose := new(SysPoolWinLose), new(SysPoolWinLose)
//	pRoomWinLose.i64PoolAmount = g_i64WinLine
//	pTable := getWeightTable(iWeightNoNow)
//	bDealerIsGood := true
//	mapi64SpecialSheet := specialSheet()
//	arypTestResult := make([]*PoolTestData, 0)
//	gmappPatternCount = make(map[int]*PatternCount)
//
//	tCurrent := time.Now()
//	strTime := fmt.Sprintf("%04d%02d%02d_%02d%02d%02d",
//		tCurrent.Year(), tCurrent.Month(), tCurrent.Day(), tCurrent.Hour(), tCurrent.Minute(), tCurrent.Second())
//	strResultFilename := fmt.Sprintf("海盜船水線自動測試[%d]_%s", g_iRoundLimit, strTime)
//	mlog.Warn("TestHDCSysPool, time: %s", strTime)
//
//	pCardBox := NewJlbpkPokerSet()
//	pDealerHandCards := NewJlbpkPokerSet() // 莊家手牌
//	var arypHandCards []*JlbpkPokerSet     // 玩家手牌
//	var arypUser []*JlbpkUser
//	for i := 0; i < iUserCount; i++ {
//		pCards := NewJlbpkPokerSet()
//		arypHandCards = append(arypHandCards, pCards)
//		pUser := NewJlbpkUser()
//		arypUser = append(arypUser, pUser)
//	}
//
//	iGameRound, i32Pattern := 0, pokerset.POKER_Pattern_None
//	pTicker := time.NewTicker(10 * time.Millisecond)
//	bTestDone := false
//	for {
//		select {
//		case <-pTicker.C:
//			pCardBox.NewSet(JLBPK_PokerSetCount, 0)
//			pCardBox.Shuffle()
//			pGameWinLose.Reset(true)
//			pDealerHandCards.SetCardSet("")
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				arypHandCards[iSeat].SetCardSet("")
//				arypUser[iSeat].ResetGameData()
//			}
//
//			// 發牌
//			pSet := dealDealerHandCard(pCardBox, JLBPK_HandCardCount, pTable)
//			pDealerHandCards.SetCardSet(pSet.CodeString())
//			i32Pattern = pDealerHandCards.Pattern(true)
//			if i32Pattern == pokerset.POKER_Pattern_None {
//				bDealerIsGood = false
//			} else {
//				bDealerIsGood = true
//			}
//			if gmappPatternCount[-1] == nil {
//				gmappPatternCount[-1] = NewPatternCount()
//			}
//			gmappPatternCount[-1].AddPattern(i32Pattern)
//
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				pSet := pCardBox.Deal(JLBPK_HandCardCount)
//				arypHandCards[iSeat].SetCardSet(pSet.CodeString())
//
//				i32Pattern = arypHandCards[iSeat].Pattern(false)
//				if gmappPatternCount[iSeat] == nil {
//					gmappPatternCount[iSeat] = NewPatternCount()
//				}
//				gmappPatternCount[iSeat].AddPattern(i32Pattern)
//			}
//
//			// 押注
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				arypUser[iSeat].i64BetBase = i64BetBase
//				arypUser[iSeat].i64BetDouble = i64BetBase * 2
//			}
//
//			// 底注結算
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				if !bDealerIsGood {
//					// 莊家不成局
//					pGameWinLose.bDealerNotOK = true
//					pGameWinLose.i64PoolAmount -= arypUser[iSeat].i64BetBase * int64(JLBPK_BetBasePayRate)
//				} else {
//					iCmpResult := pCardBox.Compare(pDealerHandCards.GetPokerCardSet(), arypHandCards[iSeat].GetPokerCardSet(), false, false)
//					i64BetBaseWinLose := int64(0)
//					switch iCmpResult {
//					case pokerset.POKER_Compare_Even: // 平手
//						i64BetBaseWinLose = arypUser[iSeat].i64BetBase
//					case pokerset.POKER_Compare_Dealer: // 莊家贏
//						i64BetBaseWinLose = -arypUser[iSeat].i64BetBase
//					case pokerset.POKER_Compare_Player: // 玩家贏
//						i64BetBaseWinLose = arypUser[iSeat].i64BetBase * int64(JLBPK_BetBasePayRate)
//					}
//
//					if iCmpResult != pokerset.POKER_Compare_Even {
//						pGameWinLose.i64PoolAmount -= i64BetBaseWinLose
//					}
//				}
//			}
//
//			// 加注結算
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				if bDealerIsGood {
//					iCmpResult := pCardBox.Compare(pDealerHandCards.GetPokerCardSet(), arypHandCards[iSeat].GetPokerCardSet(), false, false)
//					i64DoubleWinLose := int64(0)
//					switch iCmpResult {
//					case pokerset.POKER_Compare_Even: // 平手
//						i64DoubleWinLose = arypUser[iSeat].i64BetDouble
//					case pokerset.POKER_Compare_Dealer: // 第一組牌
//						i64DoubleWinLose = -arypUser[iSeat].i64BetDouble
//					case pokerset.POKER_Compare_Player: // 第二組牌
//						i64DoubleWinLose = arypUser[iSeat].i64BetDouble * mapi64SpecialSheet[arypHandCards[iSeat].Pattern(false)]
//					}
//
//					if iCmpResult != pokerset.POKER_Compare_Even {
//						pGameWinLose.i64PoolAmount -= i64DoubleWinLose
//					}
//				}
//			}
//
//			// 遊戲結束
//			if pGameWinLose.bDealerNotOK {
//				pGameWinLose.GameDBData.DealerNotOKCount++
//			} else {
//				if pGameWinLose.i64PoolAmount >= 0 {
//					pGameWinLose.GameDBData.DealerWinCount++
//				} else {
//					pGameWinLose.GameDBData.DealerLoseCount++
//				}
//			}
//
//			// 將本局輸贏統計至房間輸贏
//			pRoomWinLose.i64PoolAmount += pGameWinLose.i64PoolAmount
//			pRoomWinLose.GameDBData.DealerLoseCount += pGameWinLose.GameDBData.DealerLoseCount
//			pRoomWinLose.GameDBData.DealerNotOKCount += pGameWinLose.GameDBData.DealerNotOKCount
//			pRoomWinLose.GameDBData.DealerWinCount += pGameWinLose.GameDBData.DealerWinCount
//			pRoomWinLose.GameDBData.GameRound++
//
//			pResult := new(PoolTestData)
//			pResult.i32GameRound = pRoomWinLose.GameDBData.GameRound
//			pResult.iWeightTableNo = iWeightNoNow
//			pResult.i32DealerWinCount = pGameWinLose.GameDBData.DealerWinCount
//			pResult.i32DealerLoseCount = pGameWinLose.GameDBData.DealerLoseCount
//			pResult.i32DealerNotOkCount = pGameWinLose.GameDBData.DealerNotOKCount
//			pResult.i64PoolAmount = pRoomWinLose.i64PoolAmount
//			pResult.strDealerHandCards = fmt.Sprintf("%s(%s)", pDealerHandCards.String(), pDealerHandCards.PatternString(true))
//			i32Pattern = pDealerHandCards.Pattern(true)
//			for iSeat := 0; iSeat < iUserCount; iSeat++ {
//				str := fmt.Sprintf("%s(%s)", arypHandCards[iSeat].String(), arypHandCards[iSeat].PatternString(false))
//				pResult.arystrUserHandCards = append(pResult.arystrUserHandCards, str)
//			}
//
//			arypTestResult = append(arypTestResult, pResult)
//
//			iGameRound++
//			mlog.Warn("測試局數：%d", iGameRound)
//
//			if pRoomWinLose.GameDBData.GameRound >= int32(g_iRoundLimit) {
//				pResult := new(PoolTestData)
//				pResult.i32GameRound = CSVSum
//				pResult.i32DealerWinCount = pRoomWinLose.GameDBData.DealerWinCount
//				pResult.i32DealerLoseCount = pRoomWinLose.GameDBData.DealerLoseCount
//				pResult.i32DealerNotOkCount = pRoomWinLose.GameDBData.DealerNotOKCount
//				i32Numerator := pRoomWinLose.GameDBData.DealerWinCount
//				i32Denominator := pRoomWinLose.GameDBData.DealerLoseCount + pRoomWinLose.GameDBData.DealerNotOKCount
//				pResult.f32GameRoundRatio = float32(i32Numerator) / float32(i32Denominator)
//
//				iWeightNo, strLog := changeDealerWeightTable(pRoomWinLose, g_iRoundLimit, iWeightNoNow, g_i64WinLine)
//				if iWeightNo >= DealerWeightNo_A && iWeightNo < DealerWeightNo_Count {
//					iWeightNoNow = iWeightNo
//					pTable = getWeightTable(iWeightNoNow)
//				}
//				pResult.iWeightTableNo = iWeightNoNow
//
//				if strLog != "" {
//					mlog.Warn("%s", strLog)
//				}
//
//				arypTestResult = append(arypTestResult, pResult)
//
//				// 牌型統計標題
//				pResult1 := new(PoolTestData)
//				pResult1.i32GameRound = CSVPatternCountCaption
//				pResult1.iUserCount = iUserCount
//
//				arypTestResult = append(arypTestResult, pResult1)
//
//				// 牌型統計
//				for i32Type := pokerset.POKER_Pattern_None; i32Type <= pokerset.POKER_Pattern_RoyalFlush; i32Type++ {
//					pResult2 := new(PoolTestData)
//					pResult2.i32GameRound = CSVPatternCount
//					if i32Type == pokerset.POKER_Pattern_None {
//						pResult2.strPattern = "不成局"
//					} else {
//						pResult2.strPattern = pokerset.PatternString(i32Type)
//					}
//					pResult2.aryiPatternCountByUser = append(pResult2.aryiPatternCountByUser, gmappPatternCount[-1].aryiPatternCount[i32Type]) // 莊家
//					// 玩家
//					for iSeat := 0; iSeat < iUserCount; iSeat++ {
//						pResult2.aryiPatternCountByUser = append(pResult2.aryiPatternCountByUser, gmappPatternCount[iSeat].aryiPatternCount[i32Type])
//					}
//
//					arypTestResult = append(arypTestResult, pResult2)
//				}
//
//				// 牌型統計合計
//				pResult3 := new(PoolTestData)
//				pResult3.i32GameRound = CSVSumPatternCount
//
//				arypTestResult = append(arypTestResult, pResult3)
//				pResult3.aryiPatternCountByUser = append(pResult3.aryiPatternCountByUser, gmappPatternCount[-1].iTotalCount) // 莊家
//				// 玩家
//				for iSeat := 0; iSeat < iUserCount; iSeat++ {
//					pResult3.aryiPatternCountByUser = append(pResult3.aryiPatternCountByUser, gmappPatternCount[iSeat].iTotalCount)
//				}
//
//				gmappPatternCount = make(map[int]*PatternCount)
//			}
//
//			if iGameRound >= g_iRoundLimit {
//				bTestDone = true
//				break
//			}
//		}
//
//		if bTestDone {
//			writeCSVPoolTesting(arypTestResult, strPath, strResultFilename)
//			break
//		}
//	}
//
//	pTicker.Stop()
//}

func Exists(strName string) bool {
	if _, err := os.Stat(strName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

// 輸出 csv 檔
func writeCSVPoolTesting(arypTestResult []*PoolTestData, strPath, strFilename string) {
	if len(arypTestResult) <= 0 || strFilename == "" {
		return
	}

	strResultFile := strPath + "/" + strFilename + ".csv"
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

	// 權重表
	mapstrWeightTable := map[int]string{
		DealerWeightNo_A: "A 情況",
		DealerWeightNo_B: "B 情況",
		DealerWeightNo_C: "C 情況",
		DealerWeightNo_D: "D 情況",
	}

	var arystrResult [][]string
	// 設定標題
	arystrCaption := []string{
		"測試編號",
		"權重表",
		"莊家勝利",
		"玩家勝利",
		"莊家不成局",
		"水線",
		"莊家手牌",
		"玩家手牌",
	}

	arystrResult = append(arystrResult, arystrCaption)

	// 填入測試結果資料
	for i := 0; i < len(arypTestResult); i++ {
		if arypTestResult[i] != nil {
			var arystrData []string
			str := ""

			switch arypTestResult[i].i32GameRound {
			case CSVSum: // 合計
				arystrData = append(arystrData, "合計")
				arystrData = append(arystrData, "")                                                          // 權重表
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerWinCount))      // 莊家勝利局數
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerLoseCount))     // 玩家勝利局數
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerNotOkCount))    // 莊家不成局局數
				arystrData = append(arystrData, fmt.Sprintf("場次比值：%v", arypTestResult[i].f32GameRoundRatio)) // 場次比值
			case CSVPatternCountCaption: // 牌型統計標題
				arystrData = append(arystrData, "牌型") // 牌型
				arystrData = append(arystrData, "莊家") // 莊家
				for j := 0; j < arypTestResult[i].iUserCount; j++ {
					arystrData = append(arystrData, fmt.Sprintf("玩家 %d", j+1))
				}
			case CSVPatternCount: // 牌型統計
				arystrData = append(arystrData, fmt.Sprintf("%s", arypTestResult[i].strPattern))
				for j := 0; j < len(arypTestResult[i].aryiPatternCountByUser); j++ {
					arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].aryiPatternCountByUser[j]))
				}
			case CSVSumPatternCount: // 牌型統計合計
				arystrData = append(arystrData, "合計")
				for j := 0; j < len(arypTestResult[i].aryiPatternCountByUser); j++ {
					arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].aryiPatternCountByUser[j]))
				}
			default:
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32GameRound))                      // 測試局數
				arystrData = append(arystrData, fmt.Sprintf("%s", mapstrWeightTable[arypTestResult[i].iWeightTableNo])) // 權重表
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerWinCount))                 // 莊家勝利局數
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerLoseCount))                // 玩家勝利局數
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i32DealerNotOkCount))               // 莊家不成局局數
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].i64PoolAmount))                     // 水線
				arystrData = append(arystrData, fmt.Sprintf("%s", arypTestResult[i].strDealerHandCards))                // 莊家手牌
				str = ""
				for j := 0; j < len(arypTestResult[i].arystrUserHandCards); j++ {
					str += fmt.Sprintf("【%s】", arypTestResult[i].arystrUserHandCards[j])
				}
				arystrData = append(arystrData, fmt.Sprintf("%s", str)) // 玩家手牌
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
