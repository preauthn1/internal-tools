package poker

import "testing"

func TestSortPolicy(t *testing.T) {
	deck := NewDeck(1)
	t.Log(deck)

	deck.Sort(SortSameRank)
	t.Log(deck)

	deck.Sort(SortSameRankAceLast)
	t.Log(deck)

	deck.Sort(SortSameSuit)
	t.Log(deck)

	deck.Sort(SortSameSuitAceLast)
	t.Log(deck)
}
