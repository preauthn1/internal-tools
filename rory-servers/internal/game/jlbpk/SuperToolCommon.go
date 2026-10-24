package jlbpk

import (
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 莊家權重表編號
const (
	DealerWeightNo_A int = iota // A 情況
	DealerWeightNo_B            // B 情況
	DealerWeightNo_C            // C 情況
	DealerWeightNo_D            // D 情況
	DealerWeightNo_Count
)

//---------------------------------------------------------------------------
// 莊家牌型權重表
type DealerWeightData struct {
	i32CardPattern int32 // 牌型編號
	i32Weight      int32 // 權重值
	i32RndRangeMin int32 // 亂數值範圍最小值
	i32RndRangeMax int32 // 亂數值範圍最大值
}

type DealerWeightTable struct {
	ary_pWeightData []*DealerWeightData // 權重資料，依牌型降冪排序
	i32TotalWeight  int32               // 權重總和
}

// 遊戲輸贏結果 DB 儲存資料
type GameDataDB struct {
	WeightTableNoNow int   // 目前所使用的莊家權重表
	DealerNotOKCount int32 // 莊家不成局
	DealerWinCount   int32 // 莊家勝利
	DealerLoseCount  int32 // 玩家比牌勝利（莊家成局）
	GameRound        int32 // 進行局數
}

func (pData *GameDataDB) Copy(pOther *GameDataDB) {
	if pOther == nil {
		return
	}

	pData.WeightTableNoNow = pOther.WeightTableNoNow
	pData.DealerNotOKCount = pOther.DealerNotOKCount
	pData.DealerWinCount = pOther.DealerWinCount
	pData.DealerLoseCount = pOther.DealerLoseCount
	pData.GameRound = pOther.GameRound
}

func (pData *GameDataDB) LogString() string {
	strLog := fmt.Sprintf("WeightTableNoNow[%d], DealerNotOKCount[%d], DealerWinCount[%d], DealerLoseCount[%d], GameRound[%d]",
		pData.WeightTableNoNow, pData.DealerNotOKCount, pData.DealerWinCount, pData.DealerLoseCount, pData.GameRound)

	return strLog
}

// 遊戲輸贏結果
type SysPoolWinLose struct {
	i64PoolAmount int64      // 目前水池量
	bDealerNotOK  bool       // 莊家是否不成局
	GameDBData    GameDataDB // 遊戲輸贏結果 DB 儲存資料
}

// 清除遊戲輸贏結果
func (pWinLose *SysPoolWinLose) Reset() {
	pWinLose.i64PoolAmount = 0
	pWinLose.bDealerNotOK = false
	pWinLose.GameDBData.DealerNotOKCount = 0
	pWinLose.GameDBData.DealerWinCount = 0
	pWinLose.GameDBData.DealerLoseCount = 0
	pWinLose.GameDBData.GameRound = 0
}

// 換牌資料
type ChangingCardsData struct {
	i32seatNo         int32  // 座位
	strCardCodeString string // 手牌代碼
}

//---------------------------------------------------------------------------

var g_map_pDealerWeightTable map[int]*DealerWeightTable // 莊家權重表，以［莊家權重表編號］為 key

var g_i32PlayerWeightDefault = int32(0) // 預設玩家拿牌權重
var g_i32DealerWeightDefault = int32(0) // 預設莊家拿牌權重
var g_iCheatRate = 0                    // 換牌機率

var g_iRoundLimit = 0 // 莊家配牌機率計算場次

var g_GameData *GameDataDB = &GameDataDB{} // 系統水池-勝負資料
var g_GameDataInitialed bool               // 系統水池-勝負資料-初始化標記

var g_i64WinLine = int64(0)      // 贏錢線
var g_i32DealerWeight = int32(0) // 莊家權重

//var g_pSysPoolWinLose *SysPoolWinLose		// 房間總輸贏結果

//---------------------------------------------------------------------------
//  ----- 後台記錄用 -----
// 莊家權重表編號
var g_map_strDealerWeightNo = map[int]string{
	DealerWeightNo_A: "A 情况",
	DealerWeightNo_B: "B 情况",
	DealerWeightNo_C: "C 情况",
	DealerWeightNo_D: "D 情况",
}

// 初始化遊戲數據, 適用在遊戲開始前只會執行一次, 之後再調用都不會再執行了, 會放在遊戲開始前才進行是因為要給系統水池同步一個緩沖的時間.
func (pGame *Logic) initialGameDataOnce() {
	if !g_GameDataInitialed {
		weightTableNo := DealerWeightNo_B
		if pGame.SysPool.IsCtrlSystemWin() {
			weightTableNo = DealerWeightNo_C
		}
		g_GameDataInitialed = true
		g_GameData.WeightTableNoNow = weightTableNo
		strLog := fmt.Sprintf("初始化权重表变更为[%s]", g_map_strDealerWeightNo[weightTableNo])
		mlog.System(strLog)
	}
}

// 重置遊戲資料
func (pGame *Logic) resetGameData() {
	g_GameData.DealerWinCount = 0
	g_GameData.DealerLoseCount = 0
	g_GameData.DealerNotOKCount = 0
	g_GameData.GameRound = 0
}

// 確認是否調整權重, init 是否為桌子的初始化流程
func (pGame *Logic) changeDealerWeightTable() {
	if int(g_GameData.GameRound) >= g_iRoundLimit {
		weightTableNo := DealerWeightNo_B
		DealerWinRate := float32(0)
		if g_GameData.DealerLoseCount+g_GameData.DealerNotOKCount <= 0 {
			DealerWinRate = 2
		} else {
			DealerWinRate = float32(g_GameData.DealerWinCount) / float32(g_GameData.DealerLoseCount+g_GameData.DealerNotOKCount)
		}

		if pGame.SysPool.IsCtrlSystemWin() {
			if DealerWinRate > 1.5 { //莊家勝利/(玩家勝利+莊家不成局)>1.5
				weightTableNo = DealerWeightNo_C
			} else { //莊家勝利/(玩家勝利+莊家不成局)<=1.5
				weightTableNo = DealerWeightNo_D
			}
		} else {
			if DealerWinRate > 1.5 { //莊家勝利/(玩家勝利+莊家不成局)>1.5 (大贏, 放水表)
				weightTableNo = DealerWeightNo_A
			} else { //莊家勝利/(玩家勝利+莊家不成局)<1.5 (正常表)
				weightTableNo = DealerWeightNo_B
			}
		}
		g_GameData.WeightTableNoNow = weightTableNo
		strLog := fmt.Sprintf("游戏局数已达判定局数，庄家胜利/(玩家胜利+庄家不成局) = %f，权重表变更为[%s]",
			DealerWinRate, g_map_strDealerWeightNo[weightTableNo])
		example.invalid(strLog)
		pGame.GameLog.AddGameStep(strLog)
		pGame.resetGameData()
	}
}

//---------------------------------------------------------------------------
// 發莊家牌，pCardBox［牌盒］，iCardCount［要發的張數］, pTable［目前使用的莊家權重表］
func dealDealerHandCard(pCardBox *JlbpkPokerSet, iCardCount int, pTable *DealerWeightTable) *pokerset.PokerCardSet {
	if pCardBox == nil || pTable == nil {
		return new(pokerset.PokerCardSet)
	}

	if iCardCount <= 0 || iCardCount > pCardBox.CardCount() {
		return new(pokerset.PokerCardSet)
	}

	pRetSet := new(pokerset.PokerCardSet)

	i32Pattern := pokerset.POKERPatternHighCard
	i32Rnd := int32(util.RndInt(1, int(pTable.i32TotalWeight)+1))
	for i := 0; i < len(pTable.ary_pWeightData); i++ {
		if i32Rnd >= pTable.ary_pWeightData[i].i32RndRangeMin && i32Rnd <= pTable.ary_pWeightData[i].i32RndRangeMax {
			i32Pattern = pTable.ary_pWeightData[i].i32CardPattern
		}
	}

	pSet := pCardBox.DealByPattern(i32Pattern, iCardCount, false)
	pRetSet.SetCardSet(pSet.CodeString())
	pRetSet.Shuffle()

	mlog.Debug("dealDealerHandCard(), i32Rnd[%d], i32Pattern[%s(%d)], Cards[%s]",
		i32Rnd, pokerset.PatternString(i32Pattern), i32Pattern, pRetSet.String())

	return pRetSet
}

// 檢查水池
// 傳入 pRoomWinLose［房間總輸贏］, iRoundLimit［判定局數］, i64WinLine［贏錢線］，iWeightNo［權重表編號］
// 回傳［權重表編號］，［log］
//func changeDealerWeightTable(pRoomWinLose *SysPoolWinLose, iRoundLimit, iWeightNo int, i64WinLine int64) (int, string) {
//	if pRoomWinLose == nil || iRoundLimit <= 0 || i64WinLine <= 0 {
//		return -1, ""
//	}
//
//	iWeightTableNo := iWeightNo
//	strLog := ""
//
//	bResetWinLose := false
//	if pRoomWinLose.GameDBData.GameRound >= int32(iRoundLimit) {
//		bResetWinLose = true
//	}
//
//	f32Base, f32Cal := float32(1.5), float32(0)
//	if bResetWinLose {
//		// 依照莊家勝利場次/(玩家勝利+莊家不成局)場次的比值,和水線的變化讀取不同的莊家牌型表單
//		i32Numerator := pRoomWinLose.GameDBData.DealerWinCount
//		i32Denominator := pRoomWinLose.GameDBData.DealerLoseCount + pRoomWinLose.GameDBData.DealerNotOKCount
//		f32Cal = float32(i32Numerator) / float32(i32Denominator)
//		if i32Denominator > 0 {
//			if pRoomWinLose.i64PoolAmount > i64WinLine {
//				// 高於水線
//				if f32Cal > f32Base {
//					// 莊家勝利/(玩家勝利+莊家不成局)>1.5 (大贏, 放水表)
//					iWeightTableNo = DealerWeightNo_A
//				} else {
//					// 莊家勝利/(玩家勝利+莊家不成局)<=1.5
//					iWeightTableNo = DealerWeightNo_B
//				}
//			} else {
//				// 低於水線
//				if f32Cal > f32Base {
//					// 莊家勝利/(玩家勝利+莊家不成局)>1.5
//					iWeightTableNo = DealerWeightNo_C
//				} else {
//					// 莊家勝利/(玩家勝利+莊家不成局)<=1.5
//					iWeightTableNo = DealerWeightNo_D
//				}
//			}
//		}
//
//		strLog += fmt.Sprintf("游戏局数已达判定局数，庄家胜利/(玩家胜利+庄家不成局) = %f，权重表变更为[%s]",
//			f32Cal, g_map_strDealerWeightNo[iWeightTableNo])
//
//		pRoomWinLose.Reset(false)
//	}
//
//	return iWeightTableNo, strLog
//}

// 換牌
// 傳入 pGameCardBox［遊戲牌盒］, i32DealerWeight［莊家權重］, ary_pUserData［玩家資料］
// 回傳［換牌資料］
func changingCards(pGameCardBox *JlbpkPokerSet, i32DealerWeight int32, ary_pUserData []*JlbpkUser) []*ChangingCardsData {
	if pGameCardBox == nil || len(ary_pUserData) <= 0 {
		return make([]*ChangingCardsData, 0)
	}

	if pGameCardBox.CardCount() < JLBPK_HandCardCount*(1+len(ary_pUserData)) {
		// 牌盒內應有足夠的張數
		return make([]*ChangingCardsData, 0)
	}

	ary_pRetData := make([]*ChangingCardsData, 0)

	pCardBox := NewJlbpkPokerSet()
	pCardBox.SetCardSet(pGameCardBox.CodeString())
	var ary_pHandCards []*JlbpkPokerSet
	var ary_pUsers []*JlbpkUser

	// 莊家
	pSet := pCardBox.Deal(JLBPK_HandCardCount)
	pCardSet := NewJlbpkPokerSet()
	pCardSet.SetCardSet(pSet.CodeString())
	ary_pHandCards = append(ary_pHandCards, pCardSet)

	pUser := NewJlbpkUser()
	pUser.i32SeatNo = -1
	pUser.i32Weight = i32DealerWeight
	ary_pUsers = append(ary_pUsers, pUser)

	// 玩家
	for i := 0; i < len(ary_pUserData); i++ {
		if ary_pUserData[i] == nil {
			continue
		}

		if ary_pUserData[i].iUserState == JLB_UserState_Normal {
			pSet := pCardBox.Deal(JLBPK_HandCardCount)
			pCardSet := NewJlbpkPokerSet()
			pCardSet.SetCardSet(pSet.CodeString())
			ary_pHandCards = append(ary_pHandCards, pCardSet)

			pUser := NewJlbpkUser()
			pUser.i32SeatNo = ary_pUserData[i].i32SeatNo
			pUser.i32Weight = ary_pUserData[i].i32Weight
			ary_pUsers = append(ary_pUsers, pUser)
		}
	}

	str := "換牌前牌組"
	for i := 0; i < len(ary_pHandCards); i++ {
		str += fmt.Sprintf("【%s [%s(%d)]】", ary_pHandCards[i].CodeString(),
			ary_pHandCards[i].PatternString(false), ary_pHandCards[i].Pattern(false))
	}
	mlog.Debug("%s", str)

	// 依牌的大小排序，由大到小
	iLen := len(ary_pHandCards)
	for i := 0; i < iLen-1; i++ {
		for j := i + 1; j < iLen; j++ {
			if pCardBox.Compare(ary_pHandCards[i].GetPokerCardSet(), ary_pHandCards[j].GetPokerCardSet(), false, false) == pokerset.POKERCompareDealer {
				ary_pHandCards[i], ary_pHandCards[j] = ary_pHandCards[j], ary_pHandCards[i]
			}
		}
	}

	// 依莊家及玩家的權重排序，由小到大
	sort.Slice(ary_pUsers, func(i, j int) bool {
		return ary_pUsers[i].i32Weight < ary_pUsers[j].i32Weight
	})

	str = "權重"
	for i := 0; i < len(ary_pUsers); i++ {
		str += fmt.Sprintf("【Seat[%d], Weight[%d]】", ary_pUsers[i].i32SeatNo, ary_pUsers[i].i32Weight)
	}
	mlog.Debug("%s", str)

	for i := 0; i < len(ary_pUsers); i++ {
		if i >= len(ary_pHandCards) {
			break
		}

		pData := new(ChangingCardsData)
		pData.i32seatNo = ary_pUsers[i].i32SeatNo
		if ary_pUsers[i].i32SeatNo >= 0 {
			pData.strCardCodeString = ary_pHandCards[i].CodeString()
		} else {
			// 若莊家的牌型為［不成局］，則將原來給莊家的牌放回牌盒中，再找出可成局的高牌
			// 如果找不出可成局的高牌，則仍然發給莊家［不成局］的牌，而不換成最大的牌組
			// 這樣不至於讓莊家每把都贏監控玩家，甚至都拿到最大牌型，而不會有輸的時候
			pSet := NewJlbpkPokerSet()
			pSet.SetCardSet(ary_pHandCards[i].CodeString())
			if pSet.Pattern(true) != pokerset.POKERPatternNone {
				pData.strCardCodeString = ary_pHandCards[i].CodeString()
			} else {
				pCardBox.AddCardSet(pSet.GetPokerCardSet())
				pNewSet := pCardBox.DealByPattern(pokerset.POKERPatternHighCard, JLBPK_HandCardCount, true)
				pData.strCardCodeString = pNewSet.CodeString()
			}
		}

		ary_pRetData = append(ary_pRetData, pData)
	}

	return ary_pRetData
}
