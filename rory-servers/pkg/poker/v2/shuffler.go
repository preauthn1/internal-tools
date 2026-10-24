package poker

import (
	"math/rand"
)

// Shuffler 洗牌器
type Shuffler interface {
	Shuffle(*Deck, int) // 洗牌
}

// RandShuffler 亂數種子洗牌器
type RandShuffler struct {
	r *rand.Rand
}

// NewRandShuffler 建立亂數種子洗牌器
func NewRandShuffler(src rand.Source) *RandShuffler {
	s := new(RandShuffler)
	s.r = rand.New(src)
	return s
}

// Shuffle 洗牌
func (s *RandShuffler) shuffle(d *Deck) {
	perm := s.r.Perm(d.Len())
	var dest []Card
	for _, v := range perm {
		dest = append(dest, d.Cards[v])
	}
	d.Cards = dest
	return
}

// Shuffle 洗牌
func (s *RandShuffler) Shuffle(d *Deck, times int) {
	for i := 0; i < times; i++ {
		s.shuffle(d)
	}
}

// EmptyShuffler 不會洗牌的洗牌器
type emptyShuffler struct{}

// Shuffle 洗牌
func (s *emptyShuffler) Shuffle(d *Deck, times int) {}

// EmptyShuffler 不會洗牌的洗牌器
var EmptyShuffler Shuffler = new(emptyShuffler)
