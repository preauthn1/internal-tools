package dchjk

import (
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

type CardSetInfo struct {
	iIndex        int                    // 手中的第 N 組牌，Index 從 1 開始
	pCardSet      *Pokerset.PokerCardSet // 完整牌組
	i32Point_1    int32                  // 第一個點數，主要點數
	i32Point_2    int32                  // 第二個點數
	bIsBlackJack  bool                   // 是否為 BlackJack
	bIsBusted     bool                   // 是否爆牌
	bCanSplit     bool                   // 是否可以分牌
	bCanDouble    bool                   // 是否可以加倍
	i64BetBase    int64                  // 底注
	i64BetDouble  int64                  // 加倍的金額
	i64WinLose    int64                  // 輸贏，client 端顯示用
	bOerationDone bool                   // 此組牌操作是否已完成
	bCompensated  bool                   // 是否已賠付
	bShowCards    bool                   // 是否要顯示所有牌
}

func NewCardSetInfo(iIndex int) *CardSetInfo {
	pInfo := new(CardSetInfo)

	pInfo.iIndex = iIndex
	pInfo.pCardSet = new(Pokerset.PokerCardSet)

	pInfo.ClearData()

	return pInfo
}

// 清除資料
func (pInfo *CardSetInfo) ClearData() {
	pInfo.pCardSet.SetCardSet("")
	pInfo.i32Point_1 = -1
	pInfo.i32Point_2 = -1
	pInfo.bIsBlackJack = false
	pInfo.bIsBusted = false
	pInfo.bCanSplit = false
	pInfo.bCanDouble = false
	pInfo.i64BetBase = 0
	pInfo.i64BetDouble = 0
	pInfo.i64WinLose = 0
	pInfo.bOerationDone = false
	pInfo.bCompensated = false
	pInfo.bShowCards = true
}

// 複製資料
func (pInfo *CardSetInfo) Copy(pOther *CardSetInfo) {
	if pOther == nil {
		return
	}

	pInfo.pCardSet.SetCardSet(pOther.pCardSet.CodeString())
	pInfo.checkStatus()

	pInfo.i64BetBase = pOther.i64BetBase
	pInfo.i64BetDouble = pOther.i64BetDouble
	pInfo.i64WinLose = pOther.i64WinLose
	pInfo.bOerationDone = pOther.bOerationDone
	pInfo.bCompensated = pOther.bCompensated
	pInfo.bShowCards = pOther.bShowCards
}

// 加牌
func (pInfo *CardSetInfo) AddCardSet(pAddedSet *Pokerset.PokerCardSet) {
	if pAddedSet == nil {
		return
	}

	if pInfo.bOerationDone {
		return
	}

	pInfo.pCardSet.AddCardSet(pAddedSet)
	pInfo.checkStatus()
}

// 加一張牌
func (pInfo *CardSetInfo) AddSingleCard(pAddedCard *Pokerset.PokerCard) {
	if pAddedCard == nil {
		return
	}

	if pInfo.bOerationDone {
		return
	}

	pSet := new(Pokerset.PokerCardSet)
	pSet.AddSingleCard(pAddedCard.CodeString())
	pInfo.AddCardSet(pSet)
}

// 檢查是否為 BlackJack
func (pInfo *CardSetInfo) checkBlackJack() bool {
	if pInfo.pCardSet.CardCount() != JD21_BasicCardCount {
		return false
	}

	bHasA, bHas10 := false, false
	pSet := new(Pokerset.PokerCardSet)
	pSet.SetCardSet(pInfo.pCardSet.CodeString())
	for i := 0; i < pSet.CardCount(); i++ {
		pCard := pSet.GetSingleCard(i, true)
		if g_map_iCardNoValue[pCard.I32Number] == 10 {
			bHas10 = true
		} else if g_map_iCardNoValue[pCard.I32Number] == 1 {
			bHasA = true
		} else {
			return false
		}
	}

	if !bHasA || !bHas10 {
		return false
	}

	return true
}

// 撲克牌點數
func (pInfo *CardSetInfo) cardPoint(pCardSet *Pokerset.PokerCardSet) (int32, int32) {
	if pCardSet == nil {
		return -1, -1
	}

	i32Point1, i32Point2 := int32(0), int32(0)

	if pInfo.bIsBlackJack {
		// BlackJack
		for i := 0; i < pCardSet.CardCount(); i++ {
			pCard := pCardSet.GetSingleCard(i, false)
			if pCard != nil {
				if pCard.I32Number == Pokerset.POKERCardNo1 || pCard.I32Number == Pokerset.POKERCardNoA {
					i32Point1 += 11
				} else {
					i32Point1 += int32(g_map_iCardNoValue[pCard.I32Number])
				}
			}
		}
		i32Point2 = -1
	} else {
		// 點數 1，Ace 一律以 1 點計算
		for i := 0; i < pCardSet.CardCount(); i++ {
			pCard := pCardSet.GetSingleCard(i, false)
			if pCard != nil {
				i32Point1 += int32(g_map_iCardNoValue[pCard.I32Number])
			}
		}

		// 若有 Ace 時計算點數 2，若 Ace 有兩張（含）以上時，則其中一張以 11 點計算，其餘以 1 點計算
		if !pInfo.hasAce() {
			i32Point2 = -1
		} else {
			bBe11 := false
			for i := 0; i < pCardSet.CardCount(); i++ {
				pCard := pCardSet.GetSingleCard(i, false)
				if pCard != nil {
					if pCard.I32Number != Pokerset.POKERCardNo1 && pCard.I32Number != Pokerset.POKERCardNoA {
						i32Point2 += int32(g_map_iCardNoValue[pCard.I32Number])
					} else {
						if bBe11 {
							i32Point2 += int32(g_map_iCardNoValue[pCard.I32Number])
						} else {
							i32Point2 += 11
							bBe11 = true
						}
					}
				}
			}
		}
	}

	return i32Point1, i32Point2
}

// 檢查牌組狀態
func (pInfo *CardSetInfo) checkStatus() {
	if pInfo.pCardSet.CardCount() <= 0 {
		return
	}

	// 是否為 BlackJack
	pInfo.bIsBlackJack = pInfo.checkBlackJack()

	// 計算點數
	pInfo.i32Point_1, pInfo.i32Point_2 = pInfo.cardPoint(pInfo.pCardSet)

	// 是否爆掉
	bBusted := false
	if !pInfo.bIsBlackJack {
		if !pInfo.hasAce() {
			if pInfo.i32Point_1 > 21 {
				bBusted = true
			}
		} else {
			if pInfo.i32Point_1 > 21 && pInfo.i32Point_2 > 21 {
				bBusted = true
			}
		}
	}

	pInfo.bIsBusted = bBusted

	// 是否可分牌
	pInfo.bCanSplit = pInfo.canSplit()

	// 是否可加倍
	if pInfo.pCardSet.CardCount() == 2 {
		pInfo.bCanDouble = true
	} else {
		pInfo.bCanDouble = false
	}
}

// 牌組中是否有 A
func (pInfo *CardSetInfo) hasAce() bool {
	for i := 0; i < pInfo.pCardSet.CardCount(); i++ {
		pCard := pInfo.pCardSet.GetSingleCard(i, false)
		if pCard.I32Number == Pokerset.POKERCardNo1 || pCard.I32Number == Pokerset.POKERCardNoA {
			return true
		}
	}

	return false
}

// 是否可分牌
func (pInfo *CardSetInfo) canSplit() bool {
	if pInfo.pCardSet.CardCount() != JD21_BasicCardCount {
		return false
	}

	pCard1 := pInfo.pCardSet.GetSingleCard(0, false)
	iCardValue1 := g_map_iCardNoValue[pCard1.I32Number]

	pCard2 := pInfo.pCardSet.GetSingleCard(1, false)
	iCardValue2 := g_map_iCardNoValue[pCard2.I32Number]

	if iCardValue1 != iCardValue2 {
		return false
	}

	return true
}

// 是否已達到最低點數
func (pInfo *CardSetInfo) IsSoftReached() bool {
	if pInfo.bIsBlackJack || pInfo.bIsBusted {
		return true
	}

	if pInfo.i32Point_2 <= 0 {
		if pInfo.i32Point_1 < JD21_Soft {
			return false
		} else {
			return true
		}
	} else {
		if pInfo.i32Point_2 < JD21_Soft {
			return false
		} else {
			if pInfo.i32Point_2 <= 21 {
				return true
			} else {
				if pInfo.i32Point_1 < JD21_Soft {
					return false
				} else {
					return true
				}
			}
		}
	}
}

// 總押注額
func (pInfo *CardSetInfo) TotalBet() int64 {
	return pInfo.i64BetBase + pInfo.i64BetDouble
}

// 清除牌組
func (pInfo *CardSetInfo) ClearCardSet() {
	pInfo.pCardSet.SetCardSet("")
}

// 最大點數
func (pInfo *CardSetInfo) MaxPoint() int32 {
	if pInfo.i32Point_2 <= 0 {
		return pInfo.i32Point_1
	} else {
		if pInfo.i32Point_2 <= 21 {
			return pInfo.i32Point_2
		} else {
			if pInfo.i32Point_1 <= 21 {
				return pInfo.i32Point_1
			} else {
				return pInfo.i32Point_2
			}
		}
	}
}
