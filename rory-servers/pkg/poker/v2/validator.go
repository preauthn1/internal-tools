package poker

// Validator 牌組驗證器
type Validator interface {
	// Validate 驗證牌型，並回是否符合規則
	Validate([]Card) bool
}

type sizeValidator struct {
	accept []int
}

func (v *sizeValidator) Validate(cards []Card) bool {
	length := len(cards)
	for _, size := range v.accept {
		if size == length {
			return true
		}
	}
	return false
}

// SizeValidator 檢查長度
func SizeValidator(sizes ...int) Validator {
	v := new(sizeValidator)
	v.accept = sizes
	return v
}

type sizeLessValidator struct {
	min int
}

func (v *sizeLessValidator) Validate(cards []Card) bool {
	return len(cards) >= v.min
}

// SizeLessValidator 檢查最小長度
func SizeLessValidator(size int) Validator {
	v := new(sizeLessValidator)
	v.min = size
	return v
}

type sequenceValidator struct {
	order     RankOrder
	wildCards []Card
	wildSuits []Suit
	wildRanks []Rank
}

func (v *sequenceValidator) Validate(cards []Card) (ok bool) {
	deck := NewEmptyDeck()
	wilds := NewEmptyDeck()

	deck.AppendCard(cards...)

	for _, c := range v.wildCards {
		if _, err := deck.PickCard(c); err == nil {
			wilds.AppendCard(c)
		}
	}
	for _, s := range v.wildSuits {
		wilds.AppendCard(deck.PickSuit(s)...)
	}
	for _, r := range v.wildRanks {
		wilds.AppendCard(deck.PickRank(r)...)
	}

	var indexOf = func(order RankOrder, r Rank) int {
		for i, t := range order {
			if t == r {
				return i
			}
		}
		return -1
	}

	var idx int = -1
	ok = true
	deck.Range(func(i int, c Card) bool {
		_idx := indexOf(v.order, c.Rank())
		if idx == -1 {
			idx = _idx
			return true
		}
		if _idx != idx+1 {
			// 使用一次Wild
			if _, err := wilds.Pop(); err != nil {
				ok = false
				return false
			}
			idx++
			// 再檢查一次index是否與當前卡牌一致
			if _idx != idx+1 {
				ok = false
				return false
			}
		}
		idx++
		return true
	})

	return
}

// SequenceValidator 檢查每張牌是否皆為順序
func SequenceValidator(order RankOrder, wildCard []Card, wildSuit []Suit, wildRank []Rank) Validator {
	v := new(sequenceValidator)
	v.order = order
	v.wildCards = wildCard
	v.wildSuits = wildSuit
	v.wildRanks = wildRank
	return v
}

type sameSuitValidator struct {
	wildCards []Card
	wildSuits []Suit
	wildRanks []Rank
}

func (v *sameSuitValidator) Validate(cards []Card) (ok bool) {
	deck := NewEmptyDeck()
	wilds := NewEmptyDeck()

	deck.AppendCard(cards...)

	for _, c := range v.wildCards {
		if _, err := deck.PickCard(c); err == nil {
			wilds.AppendCard(c)
		}
	}
	for _, s := range v.wildSuits {
		wilds.AppendCard(deck.PickSuit(s)...)
	}
	for _, r := range v.wildRanks {
		wilds.AppendCard(deck.PickRank(r)...)
	}

	ok = true

	var suit Suit
	deck.Range(func(i int, c Card) bool {
		s := c.Suit()
		if suit == 0 {
			suit = s
			return true
		}
		if s != suit {
			ok = false
			return false
		}
		return true
	})

	return
}

// SameSuitValidator 檢查是否同花色
func SameSuitValidator(wildCard []Card, wildSuit []Suit, wildRank []Rank) Validator {
	v := new(sameSuitValidator)
	v.wildCards = wildCard
	v.wildSuits = wildSuit
	v.wildRanks = wildRank
	return v
}

type sameRankValidator struct {
	wildCards []Card
	wildSuits []Suit
	wildRanks []Rank
}

func (v *sameRankValidator) Validate(cards []Card) (ok bool) {
	deck := NewEmptyDeck()
	wilds := NewEmptyDeck()

	deck.AppendCard(cards...)

	for _, c := range v.wildCards {
		if _, err := deck.PickCard(c); err == nil {
			wilds.AppendCard(c)
		}
	}
	for _, s := range v.wildSuits {
		wilds.AppendCard(deck.PickSuit(s)...)
	}
	for _, r := range v.wildRanks {
		wilds.AppendCard(deck.PickRank(r)...)
	}

	ok = true

	var rank Rank
	deck.Range(func(i int, c Card) bool {
		switch c.Suit() {
		default:
			ok = false
			return false
		case Spades, Hearts, Diamonds, Clubs:
		}
		r := c.Rank()
		if rank == 0 {
			rank = r
			return true
		}
		if r != rank {
			ok = false
			return false
		}
		return true
	})

	return
}

// SameRankValidator 檢查是否同樣點數
func SameRankValidator(wildCard []Card, wildSuit []Suit, wildRank []Rank) Validator {
	v := new(sameRankValidator)
	v.wildCards = wildCard
	v.wildSuits = wildSuit
	v.wildRanks = wildRank
	return v
}
