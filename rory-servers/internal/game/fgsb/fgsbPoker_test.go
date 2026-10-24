package fgsb

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// TestPoker_Compare 測試比牌100次
func TestPoker_Compare(t *testing.T) {
	for i := 0; i < 100; i++ {
		pokerSet := poker.NewRandomPoker(HasKing, HasPizi, PokerSetCount)
		pokerSet.Shuffle()

		cards1 := new(Poker)
		cards1.SetCards(pokerSet.Deal(0, HandCardCount))
		cards2 := new(Poker)
		cards2.SetCards(pokerSet.Deal(0, HandCardCount))

		t.Logf("牌組 1：%s [%s]", cards1.String(), StrPattern[cards1.Pattern(false)])
		t.Logf("牌組 2：%s [%s]", cards2.String(), StrPattern[cards2.Pattern(false)])

		winner := cards1.Compare(cards2, false, false)
		if winner == CompareEven {
			t.Log("結果：平手")
		} else {
			t.Logf("結果：牌組 %d 贏", winner)
		}
	}

}
