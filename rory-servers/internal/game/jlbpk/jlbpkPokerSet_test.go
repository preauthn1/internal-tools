package jlbpk

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

// 建立新牌組，傳入參數［幾副撲克牌］、［鬼牌張數］
func TestJlbpkPokerSet_NewSet(t *testing.T) {
	pPokerSet := NewJlbpkPokerSet()
	pPokerSet.NewSet(1, 2)

	mlog.Debug("建立新牌組：%d 張牌", pPokerSet.CardCount())
	mlog.Debug("牌組[%s]", pPokerSet.String())
	mlog.Debug("代碼[%s]", pPokerSet.CodeString())
}

// 比牌，bCmpColor：是否要比花色
func TestJlbpkPokerSet_Compare(t *testing.T) {
	pPoker := NewJlbpkPokerSet()
	for i := pokerset.POKERPatternHighCard; i <= pokerset.POKERPatternRoyalFlush; i++ {
		for j := pokerset.POKERPatternHighCard; j <= pokerset.POKERPatternRoyalFlush; j++ {
			pPoker.NewSet(1, 0)
			pPoker.Shuffle()

			pSet1 := new(pokerset.PokerCardSet)
			pSet2 := new(pokerset.PokerCardSet)

			pSet1 = pPoker.DealByPattern(i, 5, false)
			pSet2 = pPoker.DealByPattern(j, 5, false)

			// 比牌
			pSet1.Sort(false, true, true)
			pSet2.Sort(false, true, true)
			if pSet1.Pattern() == pSet2.Pattern() {
				mlog.Debug("i = %d, j = %d", i, j)
				mlog.Debug("牌組 1：%s [%s]", pSet1.CodeString(), pokerset.PatternString(pSet1.Pattern()))
				mlog.Debug("牌組 2：%s [%s]", pSet2.CodeString(), pokerset.PatternString(pSet2.Pattern()))
				iWinner := pPoker.Compare(pSet1, pSet2, false, false)
				if iWinner == pokerset.POKERCompareEven {
					mlog.Debug("結果：平手")
				} else {
					mlog.Debug("結果：牌組 %d 贏", iWinner)
				}
				mlog.Debug("-------------------------------------------------")
			}
		}
	}
}

// 依所需牌型發牌
func TestJlbpkPokerSet_DealByPattern(t *testing.T) {
	pCardBox := NewJlbpkPokerSet()
	pCardBox.NewSet(1, 0)
	pCardBox.Shuffle()

	for i := 0; i < 7; i++ {
		pSet := pCardBox.Deal(5)
		example.invalid("CardSet %d: %s", i+1, pSet.CodeString())
	}

	pNoAKSet := pCardBox.DealByPattern(pokerset.POKERPatternHighCard, 5, false)
	mlog.Warn("正常高牌：%s", pNoAKSet.CodeString())
	pAKSet := pCardBox.DealByPattern(pokerset.POKERPatternHighCard, 5, true)
	mlog.Warn("含 A、K 高牌：%s", pAKSet.CodeString())
}
