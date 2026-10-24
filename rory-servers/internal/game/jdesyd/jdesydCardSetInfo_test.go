package jdesydLogic

import (
	"fmt"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	Pokerset "example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 檢查是否為 BlackJack
func TestIsBlackJack(t *testing.T) {
	ary_strCardSet := []string{
		"h6",
		"h6c1",
		"s1hJc7",
		"s1hQ",
		"c0h1",
		"s7h2",
		"c1s1",
	}

	pCardSet := NewCardSetInfo(1)
	for i := 0; i < len(ary_strCardSet); i++ {
		pCardSet.pCardSet.SetCardSet(ary_strCardSet[i])
		mlog.Debug("牌組[%s], BlackJack[%t]", pCardSet.pCardSet.CodeString(), pCardSet.bIsBlackJack)
	}
}

// 撲克牌點數
func TestCardPoint(t *testing.T) {
	pCardSet := NewCardSetInfo(1)
	pPokerSet := new(Pokerset.PokerCardSet)

	ary_strCard := []string{
		"s5h2c9",
		"s5h2c9s8",
		"s1hK",
		"c1d7",
		"c1d7h9",
		"s1h1c9",
		"c1h1d1d0",
	}

	for i := 0; i < len(ary_strCard); i++ {
		pSet := Pokerset.NewPokerCardSet()
		pSet.SetCardSet(ary_strCard[i])
		pCardSet.ClearCardSet()
		pCardSet.AddCardSet(pSet)

		mlog.Debug("牌組[%s], 點數[%d, %d]", pSet.CodeString(), pCardSet.i32Point_1, pCardSet.i32Point_2)
	}

	mlog.Debug("---------------------------------------------------------------------------------------------------------------")

	for i := 0; i < 50; i++ {
		pPokerSet.NewSet(JD21_PokerSetCount, 0)
		pPokerSet.Shuffle()

		iCount := util.RndInt(2, 5)
		pSet := pPokerSet.Deal(iCount)
		pCardSet.ClearCardSet()
		pCardSet.AddCardSet(pSet)

		mlog.Debug("牌組[%s], 點數[%d, %d]", pSet.CodeString(), pCardSet.i32Point_1, pCardSet.i32Point_2)
	}
}

// 檢查牌組狀態
func TestCheckPoint(t *testing.T) {
	pCardSet := NewCardSetInfo(1)
	pPokerSet := new(Pokerset.PokerCardSet)
	for i := 0; i < 50; i++ {
		pPokerSet.NewSet(6,0)
		pPokerSet.Shuffle()
		iCount := util.RndInt(2,5)
		pCardSet.pCardSet = pPokerSet.Deal(iCount)
		pCardSet.checkStatus()
		mlog.Debug("牌組[%s], 點數[%d, %d], 爆牌[%t]", pCardSet.pCardSet.CodeString(), pCardSet.i32Point_1,
			pCardSet.i32Point_2, pCardSet.bIsBusted)
	}
}

// 加牌
func TestCardSetInfo_AddCardSet(t *testing.T) {
	pCardSet := NewCardSetInfo(1)
	pPokerSet := new(Pokerset.PokerCardSet)
	for i := 0; i < 20; i++ {
		pPokerSet.NewSet(6,0)
		pPokerSet.Shuffle()
		pCardSet.ClearData()

		iCount := 0
		iRnd := util.RndInt(0,2)
		if iRnd > 0 {
			iCount = util.RndInt(2, 3)
			pSet := pPokerSet.Deal(iCount)
			pCardSet.AddCardSet(pSet)
		}
		mlog.Debug("牌組[%s], 點數[%d, %d], BJ[%t], 爆牌[%t]", pCardSet.pCardSet.CodeString(), pCardSet.i32Point_1,
			pCardSet.i32Point_2, pCardSet.bIsBlackJack, pCardSet.bIsBusted)

		// 加牌
		pAddSet := pPokerSet.Deal(1)
		pCardSet.AddCardSet(pAddSet)
		mlog.Debug("牌組[%s], 加牌[%s], 點數[%d, %d], BJ[%t], 爆牌[%t]", pCardSet.pCardSet.CodeString(), pAddSet.CodeString(),
			pCardSet.i32Point_1, pCardSet.i32Point_2, pCardSet.bIsBlackJack, pCardSet.bIsBusted)
		mlog.Debug("---------------------------------------------------------------------------------------------------------------")
	}
}

// 是否已達到最低點數
func TestCardSetInfo_IsSoftReached(t *testing.T) {
	pPokerSet := Pokerset.NewPokerCardSet()
	pTestSet := NewCardSetInfo(JD21_IndexBase)

	bNatural := false
	iCount := 0
	if bNatural {
		iCount = 50
	} else {
		iCount = 10
	}

	for i := 0; i < iCount; i++ {
		pPokerSet.NewSet(JD21_PokerSetCount,0)
		pPokerSet.Shuffle()

		pTestSet.ClearData()
		pSet := Pokerset.NewPokerCardSet()
		if bNatural {
			pSet = pPokerSet.Deal(2)
		} else {
			pSet = Pokerset.NewPokerCardSet()
			pSet.SetCardSet(fmt.Sprintf("c1s%d", i))
		}
		pTestSet.AddCardSet(pSet)
		example.invalid("【牌組】%s (%d, %d)", pTestSet.pCardSet.String(), pTestSet.i32Point_1, pTestSet.i32Point_2)
		if !pTestSet.bIsBlackJack && !pTestSet.bIsBusted && !pTestSet.IsSoftReached() {
			for {
				pAddedSet := pPokerSet.Deal(1)
				pTestSet.AddCardSet(pAddedSet)
				mlog.Debug("補牌：%s【牌組】%s (%d, %d)", pAddedSet.String(), pTestSet.pCardSet.String(), pTestSet.i32Point_1, pTestSet.i32Point_2)
				if pTestSet.IsSoftReached() {
					break
				}
			}
		}
		example.invalid("【牌組】%s (%d, %d), BlackJack[%t], Busted[%t]", pTestSet.pCardSet.String(), pTestSet.i32Point_1, pTestSet.i32Point_2,
			pTestSet.bIsBlackJack, pTestSet.bIsBusted)
		mlog.Debug("---------------------------------------------------------------------------------------------------------------")
	}
}

// 最大點數
func TestCardSetInfo_MaxPoint(t *testing.T) {
	pPokerSet := new(Pokerset.PokerCardSet)
	pCardSet := NewCardSetInfo(1)

	for i := 0; i < 20; i++ {
		pPokerSet.NewSet(1, 0)
		pPokerSet.Shuffle()
		pCardSet.ClearData()

		pSet := pPokerSet.Deal(2)
		pCardSet.AddCardSet(pSet)

		if !pCardSet.bIsBlackJack && !pCardSet.bIsBusted && !pCardSet.IsSoftReached() {
			for {
				pSet := pPokerSet.Deal(1)
				pCardSet.AddCardSet(pSet)

				if pCardSet.IsSoftReached() {
					break
				}
			}
		}

		mlog.Debug("牌組[%s], 點數[%d, %d], 最大點數[%d]", pCardSet.pCardSet.CodeString(), pCardSet.i32Point_1,
			pCardSet.i32Point_2, pCardSet.MaxPoint())
	}
}

// 加一張牌
func TestCardSetInfo_AddSingleCard(t *testing.T) {
	pPokerSet := new(Pokerset.PokerCardSet)

	for i := 0; i < 20; i++ {
		pPokerSet.NewSet(1, 0)
		pPokerSet.Shuffle()

		pCardSet := NewCardSetInfo(1)
		pInitSet := pPokerSet.Deal(2)
		pCardSet.AddCardSet(pInitSet)

		str := fmt.Sprintf("初始牌組[%s](%d, %d), 最大點數[%d], ", pCardSet.pCardSet.CodeString(),
			pCardSet.i32Point_1, pCardSet.i32Point_2, pCardSet.MaxPoint())

		pCard := pPokerSet.GetSingleCard(0,false)
		pCardSet.AddSingleCard(pCard)

		str += fmt.Sprintf("加牌[%s], 新牌組[%s](%d, %d), 最大點數[%d]", pCard.CodeString(), pCardSet.pCardSet.CodeString(),
			pCardSet.i32Point_1, pCardSet.i32Point_2, pCardSet.MaxPoint())
		mlog.Debug("%s", str)
	}
}