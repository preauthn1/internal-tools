package jdesydLogic

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 機器人 AI 設定
type RobotSetting struct {
	iRate_Hit           int // 小於17點時要牌機率
	iRate_Hit_ReachSoft int // 大於等於17點時要牌機率
	iRate_Split         int // 可分牌時分牌機率
	iRate_Double        int // 可加倍時加倍機率
	iRate_Insurrence    int // 可保險時保險機率
	iRate_Surrender     int // 莊家明牌有10或是A時投降的機率
	iBetRange_Min       int // 最小下注籌碼索引值
	iBetRange_Max       int // 最大下注籌碼索引值
	iDefaultBetCount    int // 預設下注數量
	iBetCount           int // 下注數量
}

func (pRobot *RobotSetting) Copy(pOther *RobotSetting) {
	if pOther == nil {
		return
	}

	pRobot.iRate_Hit = pOther.iRate_Hit
	pRobot.iRate_Hit_ReachSoft = pOther.iRate_Hit_ReachSoft
	pRobot.iRate_Split = pOther.iRate_Split
	pRobot.iRate_Double = pOther.iRate_Double
	pRobot.iRate_Insurrence = pOther.iRate_Insurrence
	pRobot.iRate_Surrender = pOther.iRate_Surrender
	pRobot.iBetRange_Min = pOther.iBetRange_Min
	pRobot.iBetRange_Max = pOther.iBetRange_Max
	pRobot.iDefaultBetCount = pOther.iDefaultBetCount
	pRobot.iBetCount = pOther.iBetCount
}

//---------------------------------------------------------------------------
type JdesydUser struct {
	User               *gamelogic.GameUser
	i32UserID          int32          // 玩家ID
	i32SeatNo          int32          // 座位
	iUserType          int            // 玩家類別
	i64Money           int64          // 玩家金錢
	iUserState         int            // 玩家狀態
	i64MaxBetBase      int64          // 可以下的最大底注值
	i64Timeout         int64          // 各階段 tumeout
	i64LastMessageTime int64          // 上一次訊息時間
	ary_pCardSet       []*CardSetInfo // 玩家手牌
	i64MoneyChange     int64          // 金錢變化
	i64TotalWinLose    int64          // 總輸贏，client 端顯示用
	bSurrendered       bool           // 是否已投降
	i64Insurance       int64          // 玩家買的保險金額
	iPlayingState      int            // 遊戲中狀態
	bOffline           bool           // 玩家是否為離線
	iCardSetIndexNow   int            // 現在要玩的牌組編號
	pRobot             *RobotSetting  // 機器人設定
	iBetTimes          int            // 已下注次數，用於機器人
	i64Tax             int64          // 被抽的稅
	iSumBetCount       int            // 押注次數
	i64SumBetMoney     int64          // 押注金額
}

func NewJdesydUser() *JdesydUser {
	pUser := new(JdesydUser)

	pUser.ary_pCardSet = make([]*CardSetInfo, 0)
	pUser.pRobot = new(RobotSetting)

	pUser.Init()

	return pUser
}

// 初始化玩家資料
func (pUser *JdesydUser) Init() {
	pUser.User = nil
	pUser.iUserType = gamelogic.UserTypeNormal
	pUser.i32UserID = -1
	pUser.i64Money = 0
	pUser.iUserState = JD21_UserState_EmptySeat

	pUser.ResetGameData()
}

// 重置遊戲資料
func (pUser *JdesydUser) ResetGameData() {
	pUser.ary_pCardSet = nil
	pUser.ary_pCardSet = make([]*CardSetInfo, 0)
	pCardSet := NewCardSetInfo(JD21_IndexBase)
	pUser.AddCardSet(pCardSet)

	pUser.i64MaxBetBase = g_i64MaxTotalBetBase
	pUser.i64Timeout = 0
	pUser.i64LastMessageTime = 0
	pUser.i64MoneyChange = 0
	pUser.i64TotalWinLose = 0
	pUser.bSurrendered = false
	pUser.i64Insurance = 0
	pUser.iPlayingState = JD21_UserPlayState_WaitingStart
	pUser.bOffline = false
	pUser.iCardSetIndexNow = 0
	pUser.iBetTimes = 0
	pUser.i64Tax = 0
	pUser.iSumBetCount = 0
	pUser.i64SumBetMoney = 0

	if pUser.i32UserID == -1 {
		for i := 0; i < len(pUser.ary_pCardSet); i++ {
			pUser.ary_pCardSet[i].bOerationDone = true
			pUser.ary_pCardSet[i].bCompensated = true
		}
	}
}

// 押注
func (pUser *JdesydUser) Bet(iCardIndex int, i64Bet int64) {
	if iCardIndex < 0 || iCardIndex >= len(pUser.ary_pCardSet) {
		return
	}

	if i64Bet <= 0 || i64Bet > pUser.i64Money {
		return
	}
	if len(pUser.ary_pCardSet) <= 0 {
		pNewCardSet := NewCardSetInfo(JD21_IndexBase)
		pNewCardSet.bOerationDone = false
		pUser.AddCardSet(pNewCardSet)
	}

	pUser.ary_pCardSet[iCardIndex].i64BetBase += i64Bet
	pUser.i64Money -= i64Bet

	if pUser.iUserType != gamelogic.UserTypeRebot {
		pUser.i64SumBetMoney += i64Bet
	}
}

// 底注總押注額
func (pUser *JdesydUser) TotalBetBase() int64 {
	i64TotalBet := int64(0)
	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		i64TotalBet += pUser.ary_pCardSet[i].i64BetBase
	}

	return i64TotalBet
}

// 總押注額
func (pUser *JdesydUser) TotalBet() int64 {
	i64TotalBet := int64(0)
	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		i64TotalBet += pUser.ary_pCardSet[i].i64BetBase
		i64TotalBet += pUser.ary_pCardSet[i].i64BetDouble
	}

	return i64TotalBet
}

// 新增牌組
func (pUser *JdesydUser) AddCardSet(pCardSet *CardSetInfo) bool {
	if pCardSet == nil {
		return false
	}

	if len(pUser.ary_pCardSet) >= JD21_MaxSplitCount {
		return false
	}

	bEmptySet := false
	iEmptyIdx := 0
	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		if pUser.ary_pCardSet[i].pCardSet.CardCount() <= 0 {
			bEmptySet = true
			iEmptyIdx = i
			break
		}
	}

	if bEmptySet {
		pUser.ary_pCardSet[iEmptyIdx].AddCardSet(pCardSet.pCardSet)
	} else {
		pUser.ary_pCardSet = append(pUser.ary_pCardSet, pCardSet)
	}

	return true
}

// 所有牌組是否已操作完成
func (pUser *JdesydUser) AllCardSetDone() bool {
	bAllDone := true
	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		if !pUser.ary_pCardSet[i].bOerationDone {
			pUser.iCardSetIndexNow = i
			bAllDone = false
			break
		}
	}

	if bAllDone {
		pUser.iPlayingState = JD21_UserPlayState_Playing_Done
	}

	return bAllDone
}

// 目前操作的牌組
func (pUser *JdesydUser) CardSetUsedNow() *CardSetInfo {
	if pUser.iCardSetIndexNow < 0 || pUser.iCardSetIndexNow >= len(pUser.ary_pCardSet) {
		return NewCardSetInfo(JD21_IndexBase)
	}

	return pUser.ary_pCardSet[pUser.iCardSetIndexNow]
}

// 停牌
func (pUser *JdesydUser) Stand() {
	pUser.CardSetUsedNow().bOerationDone = true
}

// 分牌
func (pUser *JdesydUser) Split(pCardSet *Pokerset.PokerCardSet) {
	if pCardSet == nil {
		return
	}

	if pCardSet.CardCount() != JD21_BasicCardCount {
		return
	}

	if len(pUser.ary_pCardSet) <= 0 || len(pUser.ary_pCardSet) >= JD21_MaxSplitCount {
		return
	}

	if !pUser.CardSetUsedNow().bCanSplit {
		return
	}

	// 檢查是否為兩張 Ace，如果是，分牌後 bOerationDone 設為 true，不能再操作
	bDoubleAce := true
	for i := 0; i < pUser.CardSetUsedNow().pCardSet.CardCount(); i++ {
		pCard := pUser.CardSetUsedNow().pCardSet.GetSingleCard(0, false)
		if pCard.I32Number != Pokerset.POKERCardNo1 && pCard.I32Number != Pokerset.POKERCardNoA {
			bDoubleAce = false
			break
		}
	}

	// 分牌方式：【c0h0】【h4sK】->【c0h1】【h4sK】[h0s5】
	// 建立分牌後的兩個牌組
	var ary_pSet [JD21_BasicCardCount]*Pokerset.PokerCardSet
	var ary_pCard [JD21_BasicCardCount]*Pokerset.PokerCard
	for i := 0; i < JD21_BasicCardCount; i++ {
		ary_pSet[i] = Pokerset.NewPokerCardSet()
		ary_pCard[i] = pUser.CardSetUsedNow().pCardSet.GetSingleCard(i, false)
		ary_pSet[i].AddSingleCard(ary_pCard[i].CodeString())

		// 補牌
		pSet := pCardSet.Deal(1)
		ary_pSet[i].AddCardSet(pSet)
	}

	// 建立新的牌組資訊
	pNewCardInfo := NewCardSetInfo(len(pUser.ary_pCardSet) + JD21_IndexBase)
	pNewCardInfo.Copy(pUser.CardSetUsedNow())
	pNewCardInfo.ClearCardSet()
	pNewCardInfo.AddCardSet(ary_pSet[1])
	if pNewCardInfo.bIsBlackJack {
		pNewCardInfo.bOerationDone = true
		pNewCardInfo.bIsBlackJack = false
	}
	if bDoubleAce {
		pNewCardInfo.bOerationDone = true
	}

	pUser.ary_pCardSet = append(pUser.ary_pCardSet, pNewCardInfo)
	pUser.i64Money -= pNewCardInfo.i64BetBase

	if pUser.iUserType != gamelogic.UserTypeRebot {
		pUser.iSumBetCount++
		pUser.i64SumBetMoney += pNewCardInfo.i64BetBase
	}

	// 更新原牌組資訊
	pUser.CardSetUsedNow().ClearCardSet()
	pUser.CardSetUsedNow().AddCardSet(ary_pSet[0])
	if pUser.CardSetUsedNow().bIsBlackJack {
		pUser.CardSetUsedNow().bOerationDone = true
		pUser.CardSetUsedNow().bIsBlackJack = false
	}
	if bDoubleAce {
		pUser.CardSetUsedNow().bOerationDone = true
	}

	// 如果已達最大分牌組數，則所有牌組不可再分牌
	if len(pUser.ary_pCardSet) >= JD21_MaxSplitCount {
		for i := 0; i < len(pUser.ary_pCardSet); i++ {
			pUser.ary_pCardSet[i].bCanSplit = false
		}
	}
}

// 加倍
func (pUser *JdesydUser) Double(pCardSet *Pokerset.PokerCardSet) {
	if pCardSet == nil {
		return
	}

	if pCardSet.CardCount() < 1 {
		return
	}

	if !pUser.CardSetUsedNow().bCanDouble {
		return
	}

	pSet := pCardSet.Deal(1)

	pUser.CardSetUsedNow().AddCardSet(pSet)
	pUser.CardSetUsedNow().i64BetDouble = pUser.CardSetUsedNow().i64BetBase
	pUser.i64Money -= pUser.CardSetUsedNow().i64BetDouble
	pUser.CardSetUsedNow().bOerationDone = true

	if pUser.iUserType != gamelogic.UserTypeRebot {
		pUser.iSumBetCount++
		pUser.i64SumBetMoney += pUser.CardSetUsedNow().i64BetDouble
	}
}

// 要牌
func (pUser *JdesydUser) Hit(pCardSet *Pokerset.PokerCardSet) {
	if pCardSet == nil {
		return
	}

	if pCardSet.CardCount() < 1 {
		return
	}

	pSet := pCardSet.Deal(1)

	pUser.CardSetUsedNow().AddCardSet(pSet)

	// 如果玩家不是 BlackJack，而點數剛好為 21 點，就算操作完成
	if !pUser.CardSetUsedNow().bIsBlackJack {
		if pUser.CardSetUsedNow().i32Point_1 == 21 || pUser.CardSetUsedNow().i32Point_2 == 21 {
			pUser.CardSetUsedNow().bOerationDone = true
		}
	}
}

// 投降
func (pUser *JdesydUser) Surrender() {
	if pUser.bSurrendered {
		return
	}

	if len(pUser.ary_pCardSet) != 1 {
		return
	}

	if pUser.CardSetUsedNow().pCardSet.CardCount() != JD21_BasicCardCount {
		return
	}

	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		pUser.ary_pCardSet[i].bOerationDone = true
	}
	pUser.bSurrendered = true
}

// 所有牌組是否均已賠付
func (pUser *JdesydUser) AllCompensated() bool {
	bAllDone := true

	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		if !pUser.ary_pCardSet[i].bCompensated {
			bAllDone = false
			break
		}
	}

	return bAllDone
}

// 玩家手上所有的牌
func (pUser *JdesydUser) WholePokerCardSet() *Pokerset.PokerCardSet {
	pRetCardSet := Pokerset.NewPokerCardSet()

	for i := 0; i < len(pUser.ary_pCardSet); i++ {
		pRetCardSet.AddCardSet(pUser.ary_pCardSet[i].pCardSet)
	}

	return pRetCardSet
}

// 設定機器人 timeout
func (pUser *JdesydUser) SetRobotTimeout(i64Timeout int64) {
	pUser.i64Timeout = 0

	if i64Timeout > 0 {
		if pUser.iUserType == gamelogic.UserTypeRebot && pUser.iUserState == JD21_UserState_Normal {
			iRnd := util.RndInt(10, int(i64Timeout/2))
			pUser.i64Timeout = int64(iRnd)
		}
	}

	mlog.Debug("機器人[%d] i64Timeout[%d], timeout %d, PlayingState[%s]", pUser.i32UserID, i64Timeout, pUser.i64Timeout, g_map_strPlayingState_debug[pUser.iPlayingState])
}
