package dczz

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	Global "example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

var gtdTimeInterval = time.Duration(10)

var gmapstrUserType = map[int]string{
	gamelogic.UserTypeNormal: "普通用戶",
	gamelogic.UserTypeRebot:  "機器人",
	gamelogic.UserTypeGuest:  "遊客",
}

// 水池測試結果資料
type PoolTestData struct {
	iPersonTime           int    // 測試人次
	iChangeCardRate       int    // 換牌機率，採百分比
	iDealerWeight         int    // 莊家權重
	bChangeCards          bool   // 換牌
	strDealerHandCards    string // 莊家手牌
	strUserHandCards      string // 玩家手牌
	i64PoolAmount         int64  // 結算後水線
	bDealerWin            bool   // 莊家比牌勝利
	bDealerWinWar         bool   // 莊家戰爭勝利
	bPlayerWin            bool   // 玩家比牌勝利
	bPlayerWinWarTie      bool   // 玩家戰爭後平手勝利
	iDealerWinCount       int    // 莊家比牌勝利次數
	iDealerWinWarCount    int    // 莊家戰爭勝利次數
	iPlayerWinCount       int    // 玩家比牌勝利次數
	iPlayerWinWarTieCount int    // 玩家戰爭後平手勝利次數
}

const (
	CSVSum int = -1 // 合計
)

// 換牌
func TestChangingCards(t *testing.T) {
	pCardBox := pokerset.NewPokerCardSet()
	pCardBox.NewSet(PokerSetCount, 0)
	pCardBox.Shuffle()

	pTicker := time.NewTicker(gtdTimeInterval * time.Millisecond)
	iUserCount, iGamePhase, iRound, iTestRound := 3, GamePhase_Normal, 0, 1

	g_pSysPoolWinLose = new(SysPoolWinLose)
	setDealerWeightRandomly()

	pDealerHandCard := pokerset.NewPokerCardSet()
	i32DealerPoint := int32(0)

	// 設定玩家
	arypUser := make([]*User, iUserCount)
	iPlayer := util.RndInt(0, iUserCount)
	for i := 0; i < len(arypUser); i++ {
		arypUser[i] = NewUser()
		arypUser[i].i32SeatNo = int32(i)
		arypUser[i].i32UserID = int32(1001 + i)
		arypUser[i].iUserState = UserState_Normal
		arypUser[i].iWeight = DefaultPlayerWeight
		if i == iPlayer {
			arypUser[i].iUserType = gamelogic.UserTypeNormal
		} else {
			arypUser[i].iUserType = util.RndInt(gamelogic.UserTypeNormal, gamelogic.UserTypeGuest)
		}
	}

	bTestDone := false

	for {
		select {
		case <-pTicker.C:
			if iGamePhase >= GamePhase_Max {
				iGamePhase = GamePhase_Normal

				// 重設莊家
				pDealerHandCard.SetCardSet("")
				i32DealerPoint = 0
				setDealerWeightRandomly()

				// 重設玩家
				iPlayer = util.RndInt(0, iUserCount)
				for i := 0; i < iUserCount; i++ {
					arypUser[i].ResetGameData()
					if i == iPlayer {
						arypUser[i].iUserType = gamelogic.UserTypeNormal
					} else {
						arypUser[i].iUserType = util.RndInt(gamelogic.UserTypeNormal, gamelogic.UserTypeGuest)
					}
				}

				// 檢查是否需使用新的牌組
				if pCardBox.CardCount() < Count_CardCount_War*(1+iUserCount) {
					pCardBox.NewSet(PokerSetCount, 0)
					pCardBox.Shuffle()
				}

				iRound++
				mlog.Warn("第 %d 局完成", iRound)
				mlog.Error("===========================================================")

				if iRound >= iTestRound {
					bTestDone = true
					break
				}
			}

			// 進行測試
			if iGamePhase == GamePhase_War {
				mlog.Error("-----------------------------------------------------------")
			}
			var arypChangingData []*ChangingCardsData

			mlog.Error("遊戲階段：%s，玩家人數：%d", g_map_strGamePhase_debug[iGamePhase], iUserCount)
			mlog.Warn("換牌前 pCardBox[%d][%s]", pCardBox.CardCount(), pCardBox.CodeString())

			// 換牌
			pCardBox, arypChangingData = changingCards(pCardBox, arypUser, iGamePhase)
			for i := 0; i < len(arypChangingData); i++ {
				pData := arypChangingData[i]
				if pData != nil {
					if pData.i32SeatNo >= 0 && pData.i32SeatNo < int32(iUserCount) {
						// 玩家
						arypUser[pData.i32SeatNo].pHandCard.AddCardSet(pData.pCardSet)
						arypUser[pData.i32SeatNo].i32Point = pData.i32Point
					} else if pData.i32SeatNo == -1 {
						// 莊家
						pDealerHandCard.AddCardSet(pData.pCardSet)
						i32DealerPoint = pData.i32Point
					}
				}
			}

			mlog.Warn("換牌後 pCardBox[%d][%s]", pCardBox.CardCount(), pCardBox.CodeString())
			mlog.Warn("換牌結果：")
			mlog.Warn("莊家：權重[%d], [%s], 點數[%d]", g_pSysPoolWinLose.GameDBData.DealerWeight, pDealerHandCard.CodeString(), i32DealerPoint)
			for i := 0; i < len(arypUser); i++ {
				pUser := arypUser[i]
				if pUser != nil {
					mlog.Warn("玩家 %d [%d][%s] [%s], 點數[%d]", i, pUser.i32UserID, gmapstrUserType[pUser.iUserType],
						pUser.pHandCard.CodeString(), pUser.i32Point)
				}
			}

			iGamePhase++
		}

		if bTestDone {
			break
		}
	}
}

// 亂數設定莊家權重
func setDealerWeightRandomly() {
	iRnd := util.RndInt(0, 3)
	switch iRnd {
	case 0:
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightDefault
	case 1:
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightAboveWinLine
	case 2:
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightBelowWinLine
	}
}

// 驗證換牌機率計算公式
func TestChangeRate(t *testing.T) {
	iPersonTimeLimit := 1000
	i64WinLine := int64(15000)
	i64PoolAmount := int64(150000)
	pWinLose := new(SysPoolWinLose)
	pWinLose.GameDBData.DealerWin = 487
	pWinLose.GameDBData.DealerWinWar = 0
	pWinLose.GameDBData.PlayerWin = 513
	pWinLose.GameDBData.PlayerWinWarTie = 0
	mlog.Debug("i64PoolAmount[%d], i64WinLine[%d], ", i64PoolAmount, i64WinLine)
	mlog.Debug("DealerWin[%d], DealerWinWar[%d], PlayerWin[%d], PlayerWinWarTie[%d]",
		pWinLose.GameDBData.DealerWin, pWinLose.GameDBData.DealerWinWar, pWinLose.GameDBData.PlayerWin, pWinLose.GameDBData.PlayerWinWarTie)

	f32A := pWinLose.ValueA()
	mlog.Debug("ValueA = %v", f32A)
	if f32A == 0 {
		f32A = 1
	}

	f32ChangeRate := float32(i64PoolAmount-i64WinLine) / float32(iPersonTimeLimit) / f32A / float32(100)
	if f32ChangeRate < 0 {
		f32ChangeRate = -f32ChangeRate
	}

	iChangeRate := int(f32ChangeRate*100 + 0.5) // +0.5 是為了四捨五入
	if iChangeRate > MaxChangeCardRate {
		iChangeRate = MaxChangeCardRate
	}
	mlog.Debug("f32ChangeRate = %v, iChangeRate = %d%%", f32ChangeRate, iChangeRate)
}

// 讀取水線設定
func getSysPoolConfig() {
	poolCfgs := cfg.GetConfig("gameserver")
	if poolCfgs == nil {
		mlog.Error("獲取配置 [gameserver] 失敗")
	} else {
		g_i64WinLine = poolCfgs.GetConfigInt64("SystemPool", "WinLine")             // 贏錢線
		g_iPersonTimeLimit = poolCfgs.GetConfigInt("SystemPool", "PersonTimeLimit") // 判定人次

		mlog.Debug("g_i64WinLine[%d], g_iPersonTimeLimit[%d]", g_i64WinLine, g_iPersonTimeLimit)
	}
}

// 水線自動測試
func TestDCZZSysPool(t *testing.T) {
	//建立測試結果存放資料夾
	strPath := "./TestResult" // 測試結果存放路徑
	bExist := exists(strPath)

	if !bExist {
		err := os.Mkdir(strPath, os.ModePerm)
		if err != nil {
			mlog.Error("建立資料夾 %s 失敗", strPath)
		}
	}

	getSysPoolConfig()

	iRound, iTestRoundLimit, i32DealerPoint, iUserCount := 0, g_iPersonTimeLimit, int32(0), 1
	i64BetBase, iMaxBet := int64(100), 5000
	pCardBox := pokerset.NewPokerCardSet()
	pDealerCard := pokerset.NewPokerCardSet()
	arypUser := make([]*User, iUserCount)
	for i := 0; i < iUserCount; i++ {
		arypUser[i] = NewUser()
		arypUser[i].i32SeatNo = int32(i)
	}

	arypTestResult := make([]*PoolTestData, 0)
	g_iShuffleRestCardCount = Count_CardCount_War * 2
	g_pSysPoolWinLose = new(SysPoolWinLose)
	g_pSysPoolWinLose.i64PoolAmount = g_i64WinLine
	pWinLose := new(SysPoolWinLose)

	strResultFilename := fmt.Sprintf("賭場戰爭水池自動測試[%d]", g_iPersonTimeLimit)

	pTicker := time.NewTicker(gtdTimeInterval * time.Millisecond)
	bTestDone := false
	for {
		select {
		case <-pTicker.C:
			if pCardBox.CardCount() < g_iShuffleRestCardCount {
				pCardBox.NewSet(PokerSetCount, 0)
				pCardBox.Shuffle()
			}

			pDealerCard.SetCardSet("")
			for i := 0; i < iUserCount; i++ {
				arypUser[i].ResetGameData()
			}

			if g_pSysPoolWinLose.i64PoolAmount < g_i64WinLine {
				g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightBelowWinLine
			} else if g_pSysPoolWinLose.i64PoolAmount > g_i64WinLine {
				g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightAboveWinLine
			} else if g_pSysPoolWinLose.i64PoolAmount == g_i64WinLine {
				g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightDefault
			}

			bChangeCard := false
			if g_pSysPoolWinLose.GameDBData.ChangeCardRate > 0 {
				iRnd := util.RndInt(0, 100)
				if iRnd < g_pSysPoolWinLose.GameDBData.ChangeCardRate {
					bChangeCard = true
				}
			}

			// 一般階段
			// 押注
			for i := 0; i < iUserCount; i++ {
				iRnd := util.RndInt(1, iMaxBet+1)
				arypUser[i].i64Bet_Base = i64BetBase * int64(iRnd)
			}

			// 發牌
			var pSet *pokerset.PokerCardSet
			var pCard *pokerset.PokerCard
			if !bChangeCard {
				// 自然發牌
				// 玩家
				for i := 0; i < iUserCount; i++ {
					pSet = pCardBox.Deal(Count_CardCount_Normal)
					arypUser[i].pHandCard.AddCardSet(pSet)
					pCard = arypUser[i].pHandCard.GetSingleCard(arypUser[i].pHandCard.CardCount()-1, true)
					arypUser[i].i32Point = pCard.I32Number
				}

				// 莊家
				pSet = pCardBox.Deal(Count_CardCount_Normal)
				pDealerCard.AddCardSet(pSet)
				pCard = pDealerCard.GetSingleCard(pDealerCard.CardCount()-1, true)
				i32DealerPoint = pCard.I32Number
			} else {
				// 換牌
				var arypChangingData []*ChangingCardsData
				pCardBox, arypChangingData = changingCards(pCardBox, arypUser, GamePhase_Normal)
				for i := 0; i < len(arypChangingData); i++ {
					if arypChangingData[i] == nil {
						continue
					}

					if arypChangingData[i].pCardSet != nil {
						if arypChangingData[i].i32SeatNo >= 0 && arypChangingData[i].i32SeatNo < int32(iUserCount) {
							// 玩家
							arypUser[arypChangingData[i].i32SeatNo].pHandCard.AddCardSet(arypChangingData[i].pCardSet)
							arypUser[arypChangingData[i].i32SeatNo].i32Point = arypChangingData[i].i32Point
						} else if arypChangingData[i].i32SeatNo == -1 {
							// 莊家
							pDealerCard.AddCardSet(arypChangingData[i].pCardSet)
							i32DealerPoint = arypChangingData[i].i32Point
						}
					}
				}
			}

			// 結算
			for i := 0; i < iUserCount; i++ {
				pUser := arypUser[i]
				if !pUser.bOperationDone {
					if pUser.i32Point != i32DealerPoint {
						pUser.bOperationDone = true

						pResult := new(PoolTestData)
						pResult.iChangeCardRate = g_pSysPoolWinLose.GameDBData.ChangeCardRate
						pResult.iDealerWeight = g_pSysPoolWinLose.GameDBData.DealerWeight
						pResult.bChangeCards = bChangeCard
						pResult.strDealerHandCards = fmt.Sprintf("%s (%d)", pDealerCard.String(), i32DealerPoint)
						pResult.strUserHandCards = fmt.Sprintf("%s (%d)", pUser.pHandCard.String(), pUser.i32Point)

						if pUser.i32Point > i32DealerPoint {
							// 玩家贏
							g_pSysPoolWinLose.i64PoolAmount -= pUser.i64Bet_Base * Rate_Normal
							g_pSysPoolWinLose.GameDBData.PlayerWin++

							pResult.bPlayerWin = true
						} else if pUser.i32Point < i32DealerPoint {
							// 莊家贏
							g_pSysPoolWinLose.i64PoolAmount += pUser.i64Bet_Base + pUser.i64Bet_Side
							g_pSysPoolWinLose.GameDBData.DealerWin++

							pResult.bDealerWin = true
						}
						pResult.i64PoolAmount = g_pSysPoolWinLose.i64PoolAmount

						g_pSysPoolWinLose.GameDBData.PersonTime++
						iRound++
						pResult.iPersonTime = g_pSysPoolWinLose.GameDBData.PersonTime

						arypTestResult = append(arypTestResult, pResult)

						pWinLose.Copy(g_pSysPoolWinLose)
						strCheckLog, strUpdateLog := checkSysPool()
						if strCheckLog != "" {
							mlog.Debug("%s", strCheckLog)
						}

						if strUpdateLog != "" {
							mlog.Warn("%s", strUpdateLog)
						}
						if strUpdateLog != "" {
							// 已達到判定人數
							pResult := new(PoolTestData)
							pResult.iPersonTime = CSVSum
							pResult.iChangeCardRate = g_pSysPoolWinLose.GameDBData.ChangeCardRate
							pResult.iDealerWinCount = pWinLose.GameDBData.DealerWin
							pResult.iDealerWinWarCount = pWinLose.GameDBData.DealerWinWar
							pResult.iPlayerWinCount = pWinLose.GameDBData.PlayerWin
							pResult.iPlayerWinWarTieCount = pWinLose.GameDBData.PlayerWinWarTie

							arypTestResult = append(arypTestResult, pResult)
						}
					}
				}
			}

			// 戰爭階段
			bEnterWar := true
			for i := 0; i < iUserCount; i++ {
				if !arypUser[i].bOperationDone {
					bEnterWar = false
					break
				}
			}

			if bEnterWar {
				// 押注
				for i := 0; i < iUserCount; i++ {
					arypUser[i].i64Bet_Double = arypUser[i].i64Bet_Base * Times_BetDouble
				}

				// 發牌
				iDealCount := Count_CardCount_War - Count_CardCount_Normal
				if !bChangeCard {
					// 自然發牌
					// 玩家
					for i := 0; i < iUserCount; i++ {
						pSet = pCardBox.Deal(iDealCount)
						for j := 0; j < pSet.CardCount(); j++ {
							if j < Count_HideCardCount {
								pCard = pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
							} else {
								pCard = pSet.GetSingleCard(j, true)
							}
							arypUser[i].pHandCard.AddSingleCard(pCard.CodeString())
						}
						pCard = arypUser[i].pHandCard.GetSingleCard(arypUser[i].pHandCard.CardCount()-1, true)
						arypUser[i].i32Point = pCard.I32Number
					}

					// 莊家
					pSet = pCardBox.Deal(iDealCount)
					for j := 0; j < pSet.CardCount(); j++ {
						if j < Count_HideCardCount {
							pCard = pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
						} else {
							pCard = pSet.GetSingleCard(j, true)
						}
						pDealerCard.AddSingleCard(pCard.CodeString())
					}
					pCard = pDealerCard.GetSingleCard(pDealerCard.CardCount()-1, true)
					i32DealerPoint = pCard.I32Number
				} else {
					// 換牌
					var arypChangingData []*ChangingCardsData
					pCardBox, arypChangingData = changingCards(pCardBox, arypUser, GamePhase_War)
					for i := 0; i < len(arypChangingData); i++ {
						if arypChangingData[i] == nil {
							continue
						}

						if arypChangingData[i].pCardSet != nil {
							if arypChangingData[i].i32SeatNo >= 0 && arypChangingData[i].i32SeatNo < int32(iUserCount) {
								// 玩家
								arypUser[arypChangingData[i].i32SeatNo].pHandCard.AddCardSet(arypChangingData[i].pCardSet)
								arypUser[arypChangingData[i].i32SeatNo].i32Point = arypChangingData[i].i32Point
							} else if arypChangingData[i].i32SeatNo == -1 {
								// 莊家
								pDealerCard.AddCardSet(arypChangingData[i].pCardSet)
								i32DealerPoint = arypChangingData[i].i32Point
							}
						}
					}
				}

				// 結算
				for i := 0; i < iUserCount; i++ {
					pUser := arypUser[i]
					if !pUser.bOperationDone {
						if pUser.i32Point != i32DealerPoint {
							pUser.bOperationDone = true

							pResult := new(PoolTestData)
							pResult.iChangeCardRate = g_pSysPoolWinLose.GameDBData.ChangeCardRate
							pResult.iDealerWeight = g_pSysPoolWinLose.GameDBData.DealerWeight
							pResult.bChangeCards = bChangeCard
							pResult.strDealerHandCards = fmt.Sprintf("%s (%d)", pDealerCard.String(), i32DealerPoint)
							pResult.strUserHandCards = fmt.Sprintf("%s (%d)", pUser.pHandCard.String(), pUser.i32Point)

							if pUser.i32Point > i32DealerPoint {
								// 玩家贏
								g_pSysPoolWinLose.i64PoolAmount -= pUser.i64Bet_Double * Rate_War
								g_pSysPoolWinLose.GameDBData.PlayerWin++

								pResult.bPlayerWin = true
							} else if pUser.i32Point < i32DealerPoint {
								// 莊家贏
								g_pSysPoolWinLose.i64PoolAmount += pUser.i64Bet_Base + pUser.i64Bet_Double
								g_pSysPoolWinLose.GameDBData.DealerWinWar++

								pResult.bDealerWinWar = true
							} else {
								// 平手
								g_pSysPoolWinLose.i64PoolAmount -= (pUser.i64Bet_Double + pUser.i64Bet_Base) * Rate_War
								g_pSysPoolWinLose.GameDBData.PlayerWinWarTie++

								pResult.bPlayerWinWarTie = true
							}
							pResult.i64PoolAmount = g_pSysPoolWinLose.i64PoolAmount

							g_pSysPoolWinLose.GameDBData.PersonTime++
							iRound++
							pResult.iPersonTime = g_pSysPoolWinLose.GameDBData.PersonTime

							arypTestResult = append(arypTestResult, pResult)

							pWinLose.Copy(g_pSysPoolWinLose)
							strCheckLog, strUpdateLog := checkSysPool()
							if strCheckLog != "" {
								mlog.Debug("%s", strCheckLog)
							}

							if strUpdateLog != "" {
								mlog.Warn("%s", strUpdateLog)
							}
							if strUpdateLog != "" {
								// 已達到判定人數
								pResult := new(PoolTestData)
								pResult.iPersonTime = CSVSum
								pResult.iChangeCardRate = g_pSysPoolWinLose.GameDBData.ChangeCardRate
								pResult.iDealerWinCount = pWinLose.GameDBData.DealerWin
								pResult.iDealerWinWarCount = pWinLose.GameDBData.DealerWinWar
								pResult.iPlayerWinCount = pWinLose.GameDBData.PlayerWin
								pResult.iPlayerWinWarTieCount = pWinLose.GameDBData.PlayerWinWarTie

								arypTestResult = append(arypTestResult, pResult)
							}
						}
					}
				}
			}

			if iRound >= iTestRoundLimit {
				bTestDone = true
				break
			}
		}

		if bTestDone {
			writeCSVPoolTesting(arypTestResult, strPath, strResultFilename)
			break
		}
	}

	pTicker.Stop()
}

func exists(strName string) bool {
	if _, err := os.Stat(strName); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}

func writeCSVPoolTesting(arypTestResult []*PoolTestData, strPath, strFilename string) {
	if len(arypTestResult) <= 0 || strFilename == "" {
		return
	}

	tCurrent := time.Now()
	strTime := fmt.Sprintf("_%04d%02d%02d_%02d%02d%02d",
		tCurrent.Year(), tCurrent.Month(), tCurrent.Day(), tCurrent.Hour(), tCurrent.Minute(), tCurrent.Second())
	strResultFile := strPath + "/" + strFilename + strTime + ".csv"
	if exists(strResultFile) {
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
		"測試人次",
		"換牌機率",
		"莊家權重",
		"換牌",
		"莊家手牌",
		"玩家手牌",
		"結算後水線",
		"莊家比牌勝利",
		"莊家戰爭勝利",
		"玩家比牌勝利",
		"玩家戰爭後平手勝利",
	}

	arystrResult = append(arystrResult, arystrCaption)

	// 填入測試結果資料
	for i := 0; i < len(arypTestResult); i++ {
		if arypTestResult[i] != nil {
			var arystrData []string
			str := ""

			switch arypTestResult[i].iPersonTime {
			case CSVSum: // 合計
				arystrData = append(arystrData, "合計")
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iChangeCardRate)) // 換牌機率
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, "")
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iDealerWinCount))       // 莊家比牌勝利
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iDealerWinWarCount))    // 莊家戰爭勝利
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iPlayerWinCount))       // 玩家比牌勝利
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iPlayerWinWarTieCount)) // 玩家戰爭後平手勝利
			default:
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iPersonTime))     // 測試人次
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iChangeCardRate)) // 換牌機率
				arystrData = append(arystrData, fmt.Sprintf("%d", arypTestResult[i].iDealerWeight))   // 莊家權重
				// 換牌
				if arypTestResult[i].bChangeCards {
					str = "V"
				} else {
					str = ""
				}
				arystrData = append(arystrData, str)
				// 莊家手牌
				arystrData = append(arystrData, arypTestResult[i].strDealerHandCards)
				// 玩家手牌
				arystrData = append(arystrData, arypTestResult[i].strUserHandCards)
				// 結算後水線
				arystrData = append(arystrData, fmt.Sprintf("%s", Global.MoneyString(arypTestResult[i].i64PoolAmount)))
				// 莊家比牌勝利
				if arypTestResult[i].bDealerWin {
					str = "V"
				} else {
					str = ""
				}
				arystrData = append(arystrData, str)
				// 莊家戰爭勝利
				if arypTestResult[i].bDealerWinWar {
					str = "V"
				} else {
					str = ""
				}
				arystrData = append(arystrData, str)
				// 玩家比牌勝利
				if arypTestResult[i].bPlayerWin {
					str = "V"
				} else {
					str = ""
				}
				arystrData = append(arystrData, str)
				// 玩家戰爭後平手勝利
				if arypTestResult[i].bPlayerWinWarTie {
					str = "V"
				} else {
					str = ""
				}
				arystrData = append(arystrData, str)
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
