package poker

import "sort"

// SortPolicy 排序策略
type SortPolicy interface {
	Sort(*Deck)
}

// SortSameSuit 同花色的牌一起
var SortSameSuit SortPolicy = new(sortSameSuit)

type sortSameSuit struct{}

func (s *sortSameSuit) Sort(d *Deck) {
	sort.Slice(d.Cards, func(i, j int) bool {
		if d.Cards[i].Suit() == d.Cards[j].Suit() {
			return d.Cards[i].Rank() < d.Cards[j].Rank()
		}
		return d.Cards[i].Suit() < d.Cards[j].Suit()
	})
}

// SortSameSuitAceLast 同花色的牌一起A最後
var SortSameSuitAceLast SortPolicy = new(sortSameSuitAceLast)

type sortSameSuitAceLast struct{}

func (s *sortSameSuitAceLast) Sort(d *Deck) {
	sort.Slice(d.Cards, func(i, j int) bool {
		if d.Cards[i].Suit() == d.Cards[j].Suit() {
			return d.Cards[i].Rank().aceLastIndexOf() < d.Cards[j].Rank().aceLastIndexOf()
		}
		return d.Cards[i].Suit() < d.Cards[j].Suit()
	})
}

// SortSameRank 同花色的牌一起
var SortSameRank SortPolicy = new(sortSameRank)

type sortSameRank struct{}

func (s *sortSameRank) Sort(d *Deck) {
	sort.Slice(d.Cards, func(i, j int) bool {
		if d.Cards[i].Rank() == d.Cards[j].Rank() {
			return d.Cards[i].Suit() < d.Cards[j].Suit()
		}
		return d.Cards[i].Rank() < d.Cards[j].Rank()
	})
}

// SortSameRankAceLast 同花色的牌一起A最後
var SortSameRankAceLast SortPolicy = new(sortSameRankAceLast)

type sortSameRankAceLast struct{}

func (s *sortSameRankAceLast) Sort(d *Deck) {
	sort.Slice(d.Cards, func(i, j int) bool {
		if d.Cards[i].Rank() == d.Cards[j].Rank() {
			return d.Cards[i].Suit() < d.Cards[j].Suit()
		}
		return d.Cards[i].Rank().aceLastIndexOf() < d.Cards[j].Rank().aceLastIndexOf()
	})
}
