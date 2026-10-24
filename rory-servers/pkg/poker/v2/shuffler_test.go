package poker

import (
	"math/rand"
	"testing"
)

func TestRandShuffler(t *testing.T) {
	deck := NewDeck(1)
	t.Logf("%p %v", deck, deck)

	src := rand.NewSource(1)
	s := NewRandShuffler(src)

	s.Shuffle(deck, 1)
	t.Logf("%p %v", deck, deck)
}

func TestEmptyShuffler(t *testing.T) {
	deck := NewDeck(1)
	t.Logf("%p %v", deck, deck)

	s := EmptyShuffler

	s.Shuffle(deck, 1)
	t.Logf("%p %v", deck, deck)
}
