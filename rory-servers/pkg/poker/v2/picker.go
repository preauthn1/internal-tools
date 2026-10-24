package poker

import (
	"sort"
)

// Picker 牌組選擇器
type Picker interface {
	// Pick 傳入卡牌查詢符合規則的牌組
	Pick([]Card) ([][]Card, bool)
}

type sequencePicker struct {
	minLen    int
	maxLen    int
	order     RankOrder
	wildCards []Card
	wildRanks []Rank
	wildSuits []Suit
}

func (p *sequencePicker) pick(cards []Card) [][]Card {

	if len(cards) < p.minLen {
		return nil
	}

	deck := NewEmptyDeck(cards...)

	wilds := NewEmptyDeck()
	wilds.AppendCard(deck.PickCards(p.wildCards...)...)
	for _, r := range p.wildRanks {
		wilds.AppendCard(deck.PickRank(r)...)
	}
	for _, s := range p.wildSuits {
		wilds.AppendCard(deck.PickSuit(s)...)
	}

	if deck.Len() == 0 {
		return nil
	}

	cards = deck.Cards

	var ans []Card
	var i, j int

	for ; i < len(p.order) && j < len(cards); i, j = i+1, j+1 {
		for ; (j+1) < len(cards) && cards[j].Rank() == cards[j+1].Rank(); j++ {
		}
		for ; p.order[i] < cards[j].Rank() && len(ans) == 0; i++ {
		}
		if p.order[i] != cards[j].Rank() {
			if len(ans) >= 1 {
				c, err := wilds.Pop()
				if err != nil {
					break
				}
				ans = append(ans, c)
				j--
				continue
			}
			if len(ans) == 0 {
				ans = append(ans, cards[j])
				continue
			}
			if len(ans) == 1 {
				j--
				continue
			}
		}

		if len(ans) == 1 && (p.order[i]+13-1)%13 != ans[0].Rank()%13 {
			ans[0] = cards[j]
			continue
		}

		ans = append(ans, cards[j])
		if len(ans) == p.maxLen {
			d := NewEmptyDeck(cards[j+1:]...)
			d.AppendDeck(wilds)
			return append([][]Card{ans}, p.pick(d.Cards)...)
		}
	}

	d := NewEmptyDeck(cards[j:]...)
	if d.Len() == 0 {
		ans = append(ans, wilds.Cards...)
	} else {
		d.AppendDeck(wilds)
	}

	if len(ans) < p.minLen {
		return p.pick(d.Cards)
	}

	return append([][]Card{ans}, p.pick(d.Cards)...)
}

func (p *sequencePicker) Pick(cards []Card) ([][]Card, bool) {

	deck := NewEmptyDeck(cards...)

	// 先過濾非正常牌
	// for _, c := range cards {
	// 	switch c.Suit() {
	// 	case Spades, Hearts, Diamonds, Clubs:
	// 		deck.AppendCard(c)
	// 	}
	// }

	var ans [][]Card
	for _ans := p.pick(deck.Cards); len(_ans) > 0; _ans = p.pick(deck.Cards) {
		deck.PickCards(_ans[0]...)
		ans = append(ans, _ans[0])

		// for _, a := range _ans {
		// 	deck.PickCards(a...)
		// 	ans = append(ans, a)
		// }
	}

	sort.Slice(ans, func(i, j int) bool {
		return ans[i][0].Rank() < ans[j][0].Rank()
	})

	return ans, len(ans) > 0
}

// SequencePicker 取得連續性(順)卡牌
func SequencePicker(min, max int, order RankOrder) Picker {
	p := new(sequencePicker)
	p.minLen = min
	p.maxLen = max
	p.order = order
	return p
}

// WildSequencePicker 取得連續性(順)卡牌
func WildSequencePicker(min, max int, order RankOrder, wildCards []Card, wildRanks []Rank, wildSuits []Suit) Picker {
	p := new(sequencePicker)
	p.minLen = min
	p.maxLen = max
	p.order = order
	p.wildCards = wildCards
	p.wildRanks = wildRanks
	p.wildSuits = wildSuits
	return p
}
