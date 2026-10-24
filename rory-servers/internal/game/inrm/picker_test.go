package inrm

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestFirstLifePicker(t *testing.T) {
	deck := poker.NewEmptyDeck()
	deck.AppendCard(poker.AceClubs, poker.TwoClubs, poker.ThreeClubs)
	deck.AppendCard(poker.TenDiamonds, poker.JackDiamonds, poker.QueenDiamonds, poker.KingDiamonds)
	deck.AppendCard(poker.ThreeSpades, poker.ThreeHearts)
	deck.AppendCard(poker.AceClubs, poker.TwoClubs, poker.ThreeClubs)
	deck.AppendCard(poker.SevenSpades)
	deck.Sort(poker.SortSameRank)

	t.Log(deck.Len(), deck)

	ans, ok := FirstLifePicker().Pick(deck.Cards)

	t.Log(ok, ans)
}

func TestSecondLifePicker(t *testing.T) {
	deck := poker.NewEmptyDeck()
	deck.AppendCard(poker.AceClubs, poker.TwoClubs, poker.ThreeClubs)
	deck.AppendCard(poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)
	deck.AppendCard(poker.ThreeSpades, poker.ThreeHearts)
	deck.AppendCard(poker.FiveHearts, poker.FiveDiamonds)
	deck.AppendCard(poker.SevenSpades, poker.RedJoker)
	deck.Sort(poker.SortSameRank)

	t.Log(deck.Len(), deck)

	ans, ok := SecondLifePicker(poker.Rank(0)).Pick(deck.Cards)

	t.Log(ok, ans)
}

func TestDefaultRefresher(t *testing.T) {
	t.Run("Victory 2", func(t *testing.T) {
		r := DefaultRefresher()

		deck := poker.NewEmptyDeck()
		deck.AppendCard(poker.SevenClubs, poker.JackClubs, poker.FiveClubs)
		deck.AppendCard(poker.SevenDiamonds, poker.EightSpades, poker.NineDiamonds)
		deck.AppendCard(poker.NineClubs, poker.TenDiamonds, poker.TenClubs, poker.KingClubs, poker.JackClubs)
		deck.AppendCard(poker.KingSpades, poker.KingSpades, poker.KingHearts)

		ok, throw, hands := r.Refresh(13, poker.Jack, deck)

		t.Log(ok, throw, hands)

		if !ok {
			t.Errorf("Refresh() ok want %t, but %t", true, ok)
		}
	})
	t.Run("Throw 1", func(t *testing.T) {
		r := DefaultRefresher()

		deck := poker.NewEmptyDeck()
		deck.AppendCard(poker.SevenClubs, poker.JackClubs, poker.FiveClubs)
		deck.AppendCard(poker.SevenDiamonds, poker.EightSpades, poker.NineDiamonds)
		deck.AppendCard(poker.NineClubs, poker.TenDiamonds, poker.TenClubs, poker.KingClubs, poker.QueenHearts)
		deck.AppendCard(poker.KingSpades, poker.KingSpades, poker.KingHearts)

		t.Log(deck.Len(), deck)
		ok, throw, hands := r.Refresh(13, poker.Jack, deck)

		t.Log(ok, throw, hands)

		if ok {
			t.Errorf("Refresh() ok want %t, but %t", false, ok)
		}
	})
	t.Run("Refresh 1", func(t *testing.T) {
		r := DefaultRefresher()

		deck := poker.NewEmptyDeck()
		deck.AppendCard(poker.SevenClubs, poker.JackClubs, poker.FiveClubs)
		deck.AppendCard(poker.SevenDiamonds, poker.EightSpades, poker.NineDiamonds)
		deck.AppendCard(poker.NineClubs, poker.TenDiamonds, poker.TenClubs, poker.KingClubs, poker.QueenHearts)
		deck.AppendCard(poker.KingSpades, poker.KingHearts)

		t.Log(deck.Len(), deck)
		ok, throw, hands := r.Refresh(13, poker.Jack, deck)

		t.Log(ok, throw, hands)

		if ok {
			t.Errorf("Refresh() ok want %t, but %t", false, ok)
		}
	})
	t.Run("整理[♥4 ♥5 ♥6] [♣5 ♥8 ♥9 ♣9 ♥J ♦Q 🃟R 🃟R ♦7 ♣7]", func(t *testing.T) {
		r := DefaultRefresher()

		deck, _ := poker.ParseDeck("♥4 ♥5 ♥6 ♣5 ♥8 ♥9 ♣9 ♥J ♦Q 🃟R 🃟R ♦7 ♣7")
		ok, throw, hands := r.Refresh(13, poker.Seven, deck)
		t.Log(ok, throw, hands)
	})
	t.Run("整理 🃟R ♠2 ♥2 ♠3 ♦3 ♥4 ♥5 ♠8 ♠9 ♠10 ♠K ♥K ♣K ♥A", func(t *testing.T) {
		r := DefaultRefresher()

		deck, _ := poker.ParseDeck("🃟R ♠2 ♥2 ♠3 ♦3 ♥4 ♥5 ♠8 ♠9 ♠10 ♠K ♥K ♣K ♥A")
		ok, throw, hands := r.Refresh(13, poker.Seven, deck)
		t.Log(ok, throw, hands)
	})
}
