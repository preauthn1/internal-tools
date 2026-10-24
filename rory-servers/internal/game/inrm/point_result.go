package inrm

import "example.invalid/ServerTeam/Servers/pkg/poker/v2"

// PointResult 點數結果
type PointResult struct {
	WildJoker      poker.Rank
	Hands          []*poker.Deck
	done           bool
	FirstLife      *poker.Deck
	SecondLife     *poker.Deck
	PureSequence   []*poker.Deck
	ImpureSequence []*poker.Deck
	Set            []*poker.Deck
	None           []*poker.Deck
}

// NewPointResult 點數結果
func NewPointResult(wildJoker poker.Rank, hands ...*poker.Deck) *PointResult {
	r := new(PointResult)
	r.WildJoker = wildJoker
	r.Hands = hands

	for _, d := range r.Hands {
		// Pure Sequence
		if PureSequenceValidator().Validate(d.Cards) {
			r.PureSequence = append(r.PureSequence, d)
			continue
		}
		// Impure Sequence
		if ImpureSequenceValidator(r.WildJoker).Validate(d.Cards) {
			r.ImpureSequence = append(r.ImpureSequence, d)
			continue
		}
		// Set
		if SetValidator(r.WildJoker).Validate(d.Cards) {
			r.Set = append(r.Set, d)
			continue
		}
		// None
		r.None = append(r.None, d)
	}

	return r
}

// Valid 驗證
func (r *PointResult) Valid() {
	if r.done {
		return
	}

	r.done = true

	// 挑出 First Life，沒有就不挑了
	if len(r.PureSequence) == 0 {
		return
	}

	r.FirstLife = r.PureSequence[0]
	r.PureSequence = r.PureSequence[1:]

	// 挑出 Second Life，沒有就不挑了
	if len(r.PureSequence) > 0 {
		r.SecondLife = r.PureSequence[0]
		r.PureSequence = r.PureSequence[1:]
	} else if len(r.ImpureSequence) == 0 {
		return
	} else {
		r.SecondLife = r.ImpureSequence[0]
		r.ImpureSequence = r.ImpureSequence[1:]
	}
}

// IsValidDeclaration 是否為合法的宣告
func (r *PointResult) IsValidDeclaration() bool {
	r.Valid()
	return r.FirstLife != nil && r.SecondLife != nil && len(r.None) == 0
}

// Points 牌型點數
//  A K Q J 算 10 分
//  2 - 10 依牌面上點數計分
//  鬼牌不算分
//  * 沒有 First Life 全算
//  * 有 First Life, 沒有 Second Life, 不算 First Life, 其他全算
func (r *PointResult) Points() int {
	r.Valid()

	deck := poker.NewEmptyDeck()
	var point int

	if r.FirstLife == nil || r.SecondLife == nil {
		deck.AppendDeck(r.PureSequence...)
		deck.AppendDeck(r.ImpureSequence...)
		deck.AppendDeck(r.Set...)
	}

	deck.AppendDeck(r.None...)

	deck.Range(func(i int, c poker.Card) bool {
		// 不是點數卡牌不算分
		if c.Suit() != poker.Spades && c.Suit() != poker.Hearts && c.Suit() != poker.Diamonds && c.Suit() != poker.Clubs {
			return true
		}
		// 鬼牌不算分
		if c.Rank() == r.WildJoker {
			return true
		}

		switch c.Rank() {
		default: // 2 - 10 依牌面上點數計分
			point += int(c.Rank())
		case poker.Ace, poker.King, poker.Queen, poker.Jack: // A K Q J 算 10 分
			point += 10
		}
		return true
	})

	return point
}
