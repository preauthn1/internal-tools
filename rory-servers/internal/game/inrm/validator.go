package inrm

import (
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

type firstLifeValidator struct{}

func (v *firstLifeValidator) Validate(cards []poker.Card) bool {
	// 檢查長度
	if ok := poker.SizeLessValidator(3).Validate(cards); !ok {
		return false
	}

	// 檢查是否同花色，不可百搭
	if ok := poker.SameSuitValidator(nil, nil, nil).Validate(cards); !ok {
		return false
	}

	var ok1, ok2 bool
	{ // 檢查是否順子(A頭)，不可百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRank)
		ok1 = poker.SequenceValidator(poker.RegularRanks, nil, nil, nil).Validate(deck.Cards)
	}

	{ // 檢查是否順子(A最後)，不可百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRankAceLast)
		ok2 = poker.SequenceValidator(poker.RegularRanksAceLast, nil, nil, nil).Validate(deck.Cards)
	}

	return ok1 || ok2
}

// FirstLifeValidator 檢驗是否為 FirstLife
//  FirstLife 為純色同花順
var FirstLifeValidator = func() poker.Validator {
	return new(firstLifeValidator)
}

type secondLifeValidator struct {
	wildRanks []poker.Rank
}

func (v *secondLifeValidator) Validate(cards []poker.Card) bool {
	// 檢查長度
	if ok := poker.SizeLessValidator(3).Validate(cards); !ok {
		return false
	}

	var ok1, ok2 bool
	{ // 檢查是否順子(A頭)，鬼牌百搭，指定Rank百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRank)
		ok1 = poker.SequenceValidator(poker.RegularRanks, nil, []poker.Suit{poker.Jokers}, v.wildRanks).Validate(deck.Cards)
	}

	{ // 檢查是否順子(A最後)，鬼牌百搭，指定Rank百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRankAceLast)
		ok2 = poker.SequenceValidator(poker.RegularRanksAceLast, nil, []poker.Suit{poker.Jokers}, v.wildRanks).Validate(deck.Cards)
	}

	return ok1 || ok2
}

// SecondLifeValidator 檢驗是否為 SecondLife
//  SecondLife 為順子或同花順 (可包含鬼牌)
var SecondLifeValidator = func(ranks ...poker.Rank) poker.Validator {
	v := new(secondLifeValidator)
	v.wildRanks = ranks
	return v
}

type setValidator struct {
	wildRanks []poker.Rank
}

func (v *setValidator) Validate(cards []poker.Card) bool {
	// 檢查長度
	if ok := poker.SizeValidator(3, 4).Validate(cards); !ok {
		return false
	}

	// 檢查是否同樣點數，鬼牌百搭，指定Rank百搭
	deck := poker.NewEmptyDeck()
	deck.AppendCard(cards...)
	deck.Sort(poker.SortSameRank)
	return poker.SameRankValidator(nil, []poker.Suit{poker.Jokers}, v.wildRanks).Validate(deck.Cards)
}

// SetValidator 檢查是否為對牌 (可包含鬼牌)
var SetValidator = func(ranks ...poker.Rank) poker.Validator {
	v := new(setValidator)
	v.wildRanks = ranks
	return v
}

type pureSequenceValidator struct{}

func (v *pureSequenceValidator) Validate(cards []poker.Card) bool {
	// 檢查長度
	if ok := poker.SizeLessValidator(3).Validate(cards); !ok {
		return false
	}

	// 檢查是否同花色，不可百搭
	if ok := poker.SameSuitValidator(nil, nil, nil).Validate(cards); !ok {
		return false
	}

	var ok1, ok2 bool
	{ // 檢查是否順子(A頭)，不可百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRank)
		ok1 = poker.SequenceValidator(poker.RegularRanks, nil, nil, nil).Validate(deck.Cards)
	}

	{ // 檢查是否順子(A最後)，不可百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRankAceLast)
		ok2 = poker.SequenceValidator(poker.RegularRanksAceLast, nil, nil, nil).Validate(deck.Cards)
	}

	return ok1 || ok2
}

// PureSequenceValidator 檢驗是否為 同花順
var PureSequenceValidator = func() poker.Validator {
	return new(pureSequenceValidator)
}

type impureSequenceValidator struct {
	wildRanks []poker.Rank
}

func (v *impureSequenceValidator) Validate(cards []poker.Card) bool {
	// 檢查長度
	if ok := poker.SizeLessValidator(3).Validate(cards); !ok {
		return false
	}

	var ok1, ok2 bool
	{ // 檢查是否順子(A頭)，鬼牌百搭，指定Rank百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRank)
		ok1 = poker.SequenceValidator(poker.RegularRanks, nil, []poker.Suit{poker.Jokers}, v.wildRanks).Validate(deck.Cards)
	}

	{ // 檢查是否順子(A最後)，鬼牌百搭，指定Rank百搭
		deck := poker.NewEmptyDeck()
		deck.AppendCard(cards...)
		deck.Sort(poker.SortSameRankAceLast)
		ok2 = poker.SequenceValidator(poker.RegularRanksAceLast, nil, []poker.Suit{poker.Jokers}, v.wildRanks).Validate(deck.Cards)
	}

	return ok1 || ok2
}

// ImpureSequenceValidator 檢驗是否為 順
var ImpureSequenceValidator = func(ranks ...poker.Rank) poker.Validator {
	v := new(impureSequenceValidator)
	v.wildRanks = ranks
	return v
}
