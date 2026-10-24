package poker

import (
	"reflect"
	"testing"
)

func TestSequencePicker(t *testing.T) {

	t.Run("pick Regular Order", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.AppendCard(TwoSpades, TwoHearts, TwoDiamonds, TwoClubs, ThreeClubs, FourDiamonds, FourClubs, FiveSpades, FiveHearts, FiveClubs, SixHearts, SixDiamonds, SixClubs, EightDiamonds, NineHearts, JackSpades, QueenClubs, KingSpades, AceHearts)
		deck.Sort(SortSameRank)

		p := new(sequencePicker)
		p.minLen = 3
		p.maxLen = 3

		p.order = RegularRanks

		ans := p.pick(deck.Cards)

		t.Log(p.order)
		t.Log(deck)
		t.Log(ans)

		if len(ans) != 3 {
			t.Errorf("pick() ans legnth want %d, but %d", 3, len(ans))
			return
		}
		exp0 := []Card{AceHearts, TwoClubs, ThreeClubs}
		if !reflect.DeepEqual(exp0, ans[0]) {
			t.Errorf("pick() ans[0] want %v, but %v", exp0, ans[0])
			return
		}
		exp1 := []Card{FourClubs, FiveClubs, SixClubs}
		if !reflect.DeepEqual(exp1, ans[1]) {
			t.Errorf("pick() ans[1] want %v, but %v", exp1, ans[1])
			return
		}
		exp2 := []Card{JackSpades, QueenClubs, KingSpades}
		if !reflect.DeepEqual(exp2, ans[2]) {
			t.Errorf("pick() ans[2] want %v, but %v", exp2, ans[2])
			return
		}

	})

	t.Run("pick AceLast Order", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.AppendCard(TwoSpades, TwoHearts, TwoDiamonds, TwoClubs, ThreeClubs, FourDiamonds, FourClubs, FiveSpades, FiveHearts, FiveClubs, SixHearts, SixDiamonds, SixClubs, EightDiamonds, NineHearts, JackSpades, QueenClubs, KingSpades, AceHearts)
		deck.Sort(SortSameRankAceLast)

		p := new(sequencePicker)
		p.minLen = 3
		p.maxLen = 5

		p.order = RegularRanksAceLast

		ans := p.pick(deck.Cards)

		t.Log(p.order)
		t.Log(deck)
		t.Log(ans)

		if len(ans) != 2 {
			t.Errorf("pick() ans legnth want %d, but %d", 2, len(ans))
			return
		}
		exp0 := []Card{TwoClubs, ThreeClubs, FourClubs, FiveClubs, SixClubs}
		if !reflect.DeepEqual(exp0, ans[0]) {
			t.Errorf("pick() ans[0] want %v, but %v", exp0, ans[0])
			return
		}
		exp1 := []Card{JackSpades, QueenClubs, KingSpades, AceHearts}
		if !reflect.DeepEqual(exp1, ans[1]) {
			t.Errorf("pick() ans[1] want %v, but %v", exp1, ans[1])
			return
		}

	})

	t.Run("Pick Regular Order", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.AppendCard(TwoSpades, TwoHearts, TwoDiamonds, TwoClubs, ThreeClubs, FourDiamonds, FourClubs, FiveSpades, FiveHearts, FiveClubs, SixHearts, SixDiamonds, SixClubs, EightDiamonds, NineHearts, JackSpades, QueenClubs, KingSpades, AceHearts)
		deck.Sort(SortSameRank)

		p := SequencePicker(3, 3, RegularRanks)
		ans, ok := p.Pick(deck.Cards)

		t.Log(RegularRanks)
		t.Log(deck)
		t.Log(ok, ans)

		if len(ans) != 4 {
			t.Errorf("pick() ans legnth want %d, but %d", 4, len(ans))
			return
		}

		exp0 := []Card{AceHearts, TwoClubs, ThreeClubs}
		if !reflect.DeepEqual(exp0, ans[0]) {
			t.Errorf("pick() ans[0] want %v, but %v", exp0, ans[0])
			return
		}
		exp1 := []Card{FourClubs, FiveClubs, SixClubs}
		if !reflect.DeepEqual(exp1, ans[1]) {
			t.Errorf("pick() ans[1] want %v, but %v", exp1, ans[1])
			return
		}
		exp2 := []Card{FourDiamonds, FiveHearts, SixDiamonds}
		if !reflect.DeepEqual(exp2, ans[2]) {
			t.Errorf("pick() ans[2] want %v, but %v", exp2, ans[2])
			return
		}
		exp3 := []Card{JackSpades, QueenClubs, KingSpades}
		if !reflect.DeepEqual(exp3, ans[3]) {
			t.Errorf("pick() ans[3] want %v, but %v", exp3, ans[3])
			return
		}

	})

	t.Run("Pick AceLast Order", func(t *testing.T) {
		deck := NewEmptyDeck()
		deck.AppendCard(TwoSpades, TwoHearts, TwoDiamonds, TwoClubs, ThreeClubs, FourDiamonds, FourClubs, FiveSpades, FiveHearts, FiveClubs, SixHearts, SixDiamonds, SixClubs, EightDiamonds, NineHearts, JackSpades, QueenClubs, KingSpades, AceHearts)
		deck.Sort(SortSameRankAceLast)

		p := SequencePicker(3, 5, RegularRanksAceLast)
		ans, ok := p.Pick(deck.Cards)

		t.Log(RegularRanksAceLast)
		t.Log(deck)
		t.Log(ok, ans)

		if len(ans) != 3 {
			t.Errorf("pick() ans legnth want %d, but %d", 3, len(ans))
			return
		}
		exp0 := []Card{TwoClubs, ThreeClubs, FourClubs, FiveClubs, SixClubs}
		if !reflect.DeepEqual(exp0, ans[0]) {
			t.Errorf("Pick() ans[0] want %v, but %v", exp0, ans[0])
			return
		}
		exp1 := []Card{FourDiamonds, FiveHearts, SixDiamonds}
		if !reflect.DeepEqual(exp1, ans[1]) {
			t.Errorf("Pick() ans[1] want %v, but %v", exp1, ans[1])
			return
		}
		exp2 := []Card{JackSpades, QueenClubs, KingSpades, AceHearts}
		if !reflect.DeepEqual(exp2, ans[2]) {
			t.Errorf("Pick() ans[2] want %v, but %v", exp2, ans[2])
			return
		}
	})

}

func TestWildSequencePicker(t *testing.T) {
	t.Run("Two", func(t *testing.T) {
		deck, _ := ParseDeck("♣5 ♥8 ♥9 ♣9 ♥J ♦Q 🃟R 🃟R ♦7 ♣7")
		deck.Sort(SortSameRankAceLast)

		p := WildSequencePicker(3, 5, RegularRanksAceLast, nil, []Rank{Seven}, []Suit{Jokers})
		ans, ok := p.Pick(deck.Cards)
		t.Log(ans, ok)

		if !ok {
			t.Errorf("Wild Pick() ok want %t, but %t", true, ok)
		}
	})

	t.Run("One", func(t *testing.T) {
		deck, _ := ParseDeck("♥9 ♥J ♦Q 🃟R 🃟R")
		deck.Sort(SortSameRankAceLast)

		p := WildSequencePicker(3, 5, RegularRanksAceLast, nil, []Rank{Seven}, []Suit{Jokers})
		ans, ok := p.Pick(deck.Cards)
		t.Log(ans, ok)

		if !ok {
			t.Errorf("Wild Pick() ok want %t, but %t", true, ok)
		}
	})

	t.Run("OneAce", func(t *testing.T) {
		deck, _ := ParseDeck("♥J ♥K ♦A 🃟R 🃟R")
		deck.Sort(SortSameRankAceLast)

		p := WildSequencePicker(3, 5, RegularRanksAceLast, nil, []Rank{Seven}, []Suit{Jokers})
		ans, ok := p.Pick(deck.Cards)
		t.Log(ans, ok)

		if !ok {
			t.Errorf("Wild Pick() ok want %t, but %t", true, ok)
		}
	})
}
