package dczz

import (
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	Global "example.invalid/ServerTeam/Servers/pkg/global"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

const (
	DefaultWinLine         int64 = 1500000 // 預設贏錢線
	DefaultPersonTimeLimit int   = 1000    // 預設判定人次
	DefaultPlayerWeight    int   = 100     // 玩家預設權重
	DefaultChangeCardRate  int   = 0       // 預設換牌機率
	MaxChangeCardRate      int   = 90      // 最大換牌機率
)

// 莊家權重
const (
	DealerWeightDefault      int = 100 // 預設值
	DealerWeightAboveWinLine int = 90  // 水線高於基準水線
	DealerWeightBelowWinLine int = 110 // 水線低於基準水線
)

//---------------------------------------------------------------------------
// 遊戲輸贏結果 DB 儲存資料
type GameDataDB struct {
	DealerWeight    int // 莊家權重
	ChangeCardRate  int // 換牌機率，採百分比
	PersonTime      int // 已統計人次
	DealerWin       int // 莊家比牌勝利
	DealerWinWar    int // 莊家戰爭勝利
	PlayerWin       int // 玩家比牌勝利
	PlayerWinWarTie int // 玩家戰爭後平手勝利
}

func (pData *GameDataDB) Copy(pOther *GameDataDB) {
	if pOther == nil {
		return
	}

	pData.DealerWeight = pOther.DealerWeight
	pData.ChangeCardRate = pOther.ChangeCardRate
	pData.PersonTime = pOther.PersonTime
	pData.DealerWin = pOther.DealerWin
	pData.DealerWinWar = pOther.DealerWinWar
	pData.PlayerWin = pOther.PlayerWin
	pData.PlayerWinWarTie = pOther.PlayerWinWarTie
}

func (pData *GameDataDB) LogString() string {
	strLog := fmt.Sprintf("DealerWeight[%d], ChangeCardRate[%d], PersonTime[%d], DealerWin[%d], DealerWinWar[%d], PlayerWin[%d], PlayerWinWarTie[%d]",
		pData.DealerWeight, pData.ChangeCardRate, pData.PersonTime, pData.DealerWin, pData.DealerWinWar, pData.PlayerWin, pData.PlayerWinWarTie)

	return strLog
}

//---------------------------------------------------------------------------
// 遊戲輸贏結果
type SysPoolWinLose struct {
	i64PoolAmount int64      // 水線
	GameDBData    GameDataDB // 遊戲輸贏結果 DB 儲存資料
}

func (pWinLose *SysPoolWinLose) Copy(pOther *SysPoolWinLose) {
	if pOther == nil {
		return
	}

	pWinLose.i64PoolAmount = pOther.i64PoolAmount
	pWinLose.GameDBData.Copy(&pOther.GameDBData)
}

// bClearPool［是否清除水池］
func (pWinLose *SysPoolWinLose) Reset() {
	pWinLose.GameDBData.PersonTime = 0
	pWinLose.GameDBData.PlayerWin = 0
	pWinLose.GameDBData.PlayerWinWarTie = 0
	pWinLose.GameDBData.DealerWin = 0
	pWinLose.GameDBData.DealerWinWar = 0
}

// (莊家勝利場次+莊家戰爭勝利場次)/(玩家比牌+戰爭勝利+玩家戰爭+平點勝利)場次的比值A
func (pWinLose *SysPoolWinLose) ValueA() float32 {
	f32A := float32(0)

	i32Numerator := pWinLose.GameDBData.DealerWin + pWinLose.GameDBData.DealerWinWar
	i32Denominator := pWinLose.GameDBData.PlayerWin + pWinLose.GameDBData.PlayerWinWarTie
	if i32Denominator == 0 {
		i32Denominator = 1
	}

	f32A = float32(i32Numerator) / float32(i32Denominator)

	return f32A
}

//---------------------------------------------------------------------------
// 換牌資料
type ChangingCardsData struct {
	i32SeatNo int32                  // 座位，-1 為莊家，>= 0 為玩家
	pCardSet  *pokerset.PokerCardSet // 牌組資料
	i32Point  int32                  // 點數
}

func NewChangingCardsData() *ChangingCardsData {
	pData := new(ChangingCardsData)
	pData.pCardSet = pokerset.NewPokerCardSet()

	return pData
}

//---------------------------------------------------------------------------
var g_i64WinLine = DefaultWinLine               // 贏錢線
var g_iPersonTimeLimit = DefaultPersonTimeLimit // 判定人次
var g_bEnablePoolControl = false                // 是否開啟水線控制功能
var g_bGetWinLoseData = false                   // 是否已自 DB 取得遊戲輸贏統計資料
var g_bGetSysPoolConfig = false                 // 是否已取得水池 config

var g_pSysPoolWinLose *SysPoolWinLose // 房間總輸贏結果

//---------------------------------------------------------------------------
// 換牌
// 傳入 pGameCardBox［遊戲牌盒］、pGameCardBox［玩家資料］、iGamePhase［遊戲階段］
// 回傳 ［發完牌後的遊戲牌盒］、［換牌資料］
func changingCards(pGameCardBox *pokerset.PokerCardSet, ary_pUserData []*User, iGamePhase int) (*pokerset.PokerCardSet, []*ChangingCardsData) {
	if pGameCardBox == nil || len(ary_pUserData) <= 0 {
		return pGameCardBox, make([]*ChangingCardsData, 0)
	}

	// 牌盒內應有足夠的張數
	iUserCount := 0
	for i := 0; i < len(ary_pUserData); i++ {
		if ary_pUserData[i] == nil {
			continue
		}

		if ary_pUserData[i].iUserState == UserState_Normal {
			iUserCount++
		}
	}
	if pGameCardBox.CardCount() < Count_CardCount_War*(1+iUserCount) {
		return pGameCardBox, make([]*ChangingCardsData, 0)
	}

	if iGamePhase < GamePhase_Normal || iGamePhase >= GamePhase_Max {
		return pGameCardBox, make([]*ChangingCardsData, 0)
	}

	mlog.Debug("換牌，階段：%s", g_map_strGamePhase_debug[iGamePhase])
	str := ""
	for i := 0; i < len(ary_pUserData); i++ {
		str += fmt.Sprintf("［UserID[%d(%d)] Usertype[%d]］", ary_pUserData[i].i32UserID, ary_pUserData[i].i32SeatNo, ary_pUserData[i].iUserType)
	}
	mlog.Debug("玩家%s", str)

	ary_pRetData := make([]*ChangingCardsData, 0)

	pCardBox := pokerset.NewPokerCardSet()
	pCardBox.SetCardSet(pGameCardBox.CodeString())
	var ary_pCard []*pokerset.PokerCard
	var ary_pUsers []*User

	iDealCount := 0
	switch iGamePhase {
	case GamePhase_Normal: // 一般階段
		iDealCount = Count_CardCount_Normal
	case GamePhase_War: // 戰爭階段
		iDealCount = Count_CardCount_War - Count_CardCount_Normal
	}

	// 由於水線控制是針對真人的部分做控制，所以此處僅先發真人的牌來進行換牌
	// 換完牌後再補發機器人的部分，這樣在畫面上就不會讓莊家的牌看起來都是不是最大就是最小

	// 玩家
	for i := 0; i < len(ary_pUserData); i++ {
		if ary_pUserData[i] == nil {
			continue
		}

		if ary_pUserData[i].iUserState == UserState_Normal && ary_pUserData[i].iUserType == gamelogic.UserTypeNormal &&
			!ary_pUserData[i].bOperationDone {
			pSet := pCardBox.Deal(iDealCount)
			pCard := pSet.GetSingleCard(pSet.CardCount()-1, true)
			ary_pCard = append(ary_pCard, pCard)

			pUser := NewUser()
			pUser.i32SeatNo = ary_pUserData[i].i32SeatNo
			pUser.iWeight = ary_pUserData[i].iWeight
			ary_pUsers = append(ary_pUsers, pUser)
		}
	}

	// 莊家
	pSet := pCardBox.Deal(iDealCount)
	pCard := pSet.GetSingleCard(pSet.CardCount()-1, true)
	ary_pCard = append(ary_pCard, pCard)

	pDealer := NewUser()
	pDealer.i32SeatNo = -1
	pDealer.iWeight = g_pSysPoolWinLose.GameDBData.DealerWeight
	ary_pUsers = append(ary_pUsers, pDealer)

	str = "換牌前，ary_pCard"
	for i := 0; i < len(ary_pCard); i++ {
		str += fmt.Sprintf("［%s］", ary_pCard[i].CodeString())
	}
	str += ", ary_pUsers"
	for i := 0; i < len(ary_pUsers); i++ {
		str += fmt.Sprintf("［UserID[%d(%d)]Weight[%d]", ary_pUsers[i].i32UserID, ary_pUsers[i].i32SeatNo, ary_pUsers[i].iWeight)
	}

	// 依牌的點數大小排序，由小到大
	sort.Slice(ary_pCard, func(i, j int) bool {
		return ary_pCard[i].I32Number < ary_pCard[j].I32Number
	})

	// 依玩家權重排序，由小到大
	sort.Slice(ary_pUsers, func(i, j int) bool {
		return ary_pUsers[i].iWeight < ary_pUsers[j].iWeight
	})

	str = "排序，ary_pCard"
	for i := 0; i < len(ary_pCard); i++ {
		str += fmt.Sprintf("［%s］", ary_pCard[i].CodeString())
	}
	str += ", ary_pUsers"
	for i := 0; i < len(ary_pUsers); i++ {
		str += fmt.Sprintf("［UserID[%d(%d)]Weight[%d]", ary_pUsers[i].i32UserID, ary_pUsers[i].i32SeatNo, ary_pUsers[i].iWeight)
	}

	// 進行換牌
	for i := 0; i < len(ary_pUsers); i++ {
		if i >= len(ary_pCard) {
			break
		}

		pData := NewChangingCardsData()
		pData.i32SeatNo = ary_pUsers[i].i32SeatNo
		pData.i32Point = ary_pCard[i].I32Number
		if iGamePhase == GamePhase_War {
			for j := 0; j < Count_HideCardCount; j++ {
				pHideCard := pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
				pData.pCardSet.AddSingleCard(pHideCard.CodeString())
			}
		}
		pData.pCardSet.AddSingleCard(ary_pCard[i].CodeString())

		ary_pRetData = append(ary_pRetData, pData)
	}

	// 補發牌給機器人
	for i := 0; i < len(ary_pUserData); i++ {
		if ary_pUserData[i] == nil {
			continue
		}

		if ary_pUserData[i].iUserState == UserState_Normal && ary_pUserData[i].iUserType == gamelogic.UserTypeRebot &&
			!ary_pUserData[i].bOperationDone {
			pSet := pCardBox.Deal(iDealCount)
			pCard := pSet.GetSingleCard(pSet.CardCount()-1, true)

			pData := NewChangingCardsData()
			pData.i32SeatNo = ary_pUserData[i].i32SeatNo
			pData.i32Point = pCard.I32Number
			if iGamePhase == GamePhase_War {
				for j := 0; j < Count_HideCardCount; j++ {
					pHideCard := pokerset.NewPokerCard(pokerset.POKERCardColorOther, pokerset.POKERCardNoHide)
					pData.pCardSet.AddSingleCard(pHideCard.CodeString())
				}
			}
			pData.pCardSet.AddSingleCard(pCard.CodeString())

			ary_pRetData = append(ary_pRetData, pData)
		}
	}

	pGameCardBox.SetCardSet(pCardBox.CodeString())

	return pGameCardBox, ary_pRetData
}

// 水線判定
// 回傳［本次結果 log］、［達到判定人次結果 log］
func checkSysPool() (string, string) {
	strCheckLog, strUpdateLog := "", ""

	// 檢查莊家權重
	if g_pSysPoolWinLose.i64PoolAmount > g_i64WinLine {
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightAboveWinLine
	} else if g_pSysPoolWinLose.i64PoolAmount < g_i64WinLine {
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightBelowWinLine
	} else {
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightDefault
	}

	strCheckLog = fmt.Sprintf("第 %d 人次，判定人次 %d，水线 %s，庄家比牌胜利 %d 次，庄家战争胜利 %d 次，玩家比牌胜利 %d 次，玩家战争后平手胜利 %d 次，庄家权重 %d",
		g_pSysPoolWinLose.GameDBData.PersonTime, g_iPersonTimeLimit, Global.MoneyString(g_pSysPoolWinLose.i64PoolAmount),
		g_pSysPoolWinLose.GameDBData.DealerWin, g_pSysPoolWinLose.GameDBData.DealerWinWar,
		g_pSysPoolWinLose.GameDBData.PlayerWin, g_pSysPoolWinLose.GameDBData.PlayerWinWarTie, g_pSysPoolWinLose.GameDBData.DealerWeight)

	bWinLose := false
	if g_pSysPoolWinLose.GameDBData.PersonTime >= g_iPersonTimeLimit {
		bWinLose = true
	}

	if bWinLose {
		// 達到判定人次
		f32A := g_pSysPoolWinLose.ValueA()
		if f32A == 0 {
			f32A = 1
		}

		i64PoolAmount := g_pSysPoolWinLose.i64PoolAmount / 100
		i64WinLine := g_i64WinLine / 100
		f32ChangeRate := float32(i64PoolAmount-i64WinLine) / float32(g_iPersonTimeLimit) / f32A / float32(100)
		if f32ChangeRate < 0 {
			f32ChangeRate = -f32ChangeRate
		}

		g_pSysPoolWinLose.GameDBData.ChangeCardRate = int(f32ChangeRate*100 + 0.5) // +0.5 是為了四捨五入
		if g_pSysPoolWinLose.GameDBData.ChangeCardRate > MaxChangeCardRate {
			g_pSysPoolWinLose.GameDBData.ChangeCardRate = MaxChangeCardRate
		}
		mlog.Debug("checkSysPool(), i64PoolAmount[%d], ValueA[%v], f32ChangeRate[%v], ChangeCardRate[%d]",
			g_pSysPoolWinLose.i64PoolAmount, g_pSysPoolWinLose.ValueA(), f32ChangeRate, g_pSysPoolWinLose.GameDBData.ChangeCardRate)

		strUpdateLog = fmt.Sprintf("统计人次已达判定人次，换牌几率变更为 %d%%", g_pSysPoolWinLose.GameDBData.ChangeCardRate)

		g_pSysPoolWinLose.Reset()
	}

	return strCheckLog, strUpdateLog
}
