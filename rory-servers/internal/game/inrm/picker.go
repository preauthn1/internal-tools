package inrm

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// FirstLifePicker 抓取FistLife
var FirstLifePicker = func() poker.Picker { return new(firstLifePicker) }

type firstLifePicker struct{}

func (p *firstLifePicker) Pick(cards []poker.Card) ([][]poker.Card, bool) {

	var ans [][]poker.Card

	deck := poker.NewEmptyDeck()
	deck.AppendCard(cards...)

	var pick = func(s poker.Suit) {
		var tmp [][]poker.Card
		cards := deck.PickSuit(s)
		if set, ok := poker.SequencePicker(3, 3, poker.RegularRanks).Pick(cards); ok {
			tmp = append(tmp, set...)
		}
		if set, ok := poker.SequencePicker(3, 3, poker.RegularRanksAceLast).Pick(cards); ok {
			tmp = append(tmp, set...)
		}
		sort.Slice(tmp, func(i, j int) bool {
			return len(tmp[i]) > len(tmp[j])
		})
		if len(tmp) > 0 {
			ans = append(ans, tmp[0])
			deck.PickCards(tmp[0]...)

			for sets, ok := p.Pick(deck.Cards); ok; sets, ok = p.Pick(deck.Cards) {
				ans = append(ans, sets[0])
				deck.PickCards(sets[0]...)
			}
		}
	}

	for _, s := range poker.RegularSuits {
		pick(s)
	}

	return ans, len(ans) > 0
}

// SecondLifePicker 抓取SecondLife
var SecondLifePicker = func(wildRank poker.Rank) poker.Picker {
	p := new(secondLifePicker)
	p.wildRank = wildRank
	return p
}

type secondLifePicker struct {
	wildRank poker.Rank
}

func (p *secondLifePicker) Pick(cards []poker.Card) ([][]poker.Card, bool) {

	var ans [][]poker.Card

	deck := poker.NewEmptyDeck()
	deck.AppendCard(cards...)

	var tmp [][]poker.Card
	if set, ok := poker.WildSequencePicker(3, 3, poker.RegularRanks, nil, []poker.Rank{p.wildRank}, []poker.Suit{poker.Jokers}).Pick(cards); ok {
		tmp = append(tmp, set...)
	}
	if set, ok := poker.WildSequencePicker(3, 3, poker.RegularRanksAceLast, nil, []poker.Rank{p.wildRank}, []poker.Suit{poker.Jokers}).Pick(cards); ok {
		tmp = append(tmp, set...)
	}
	sort.Slice(tmp, func(i, j int) bool {
		return len(tmp[i]) > len(tmp[j])
	})
	if len(tmp) > 0 {
		ans = append(ans, tmp[0])
		deck.PickCards(tmp[0]...)

		for sets, ok := p.Pick(deck.Cards); ok; sets, ok = p.Pick(deck.Cards) {
			ans = append(ans, sets[0])
			deck.PickCards(sets[0]...)
		}
	}

	return ans, len(ans) > 0
}

// SetPicker 抓取Set
var SetPicker = func(wildRank poker.Rank) poker.Picker {
	p := new(setPicker)
	p.wildRank = wildRank
	return p
}

type setPicker struct {
	wildRank poker.Rank
}

func (p *setPicker) Pick(cards []poker.Card) ([][]poker.Card, bool) {

	deck := poker.NewEmptyDeck()
	deck.AppendCard(cards...)

	wilds := poker.NewEmptyDeck()
	wilds.AppendCard(deck.PickSuit(poker.Jokers)...)
	wilds.AppendCard(deck.PickRank(p.wildRank)...)

	rMap := make(map[poker.Rank]int)

	for _, c := range deck.Cards {
		s := c.Suit()
		if s != poker.Spades && s != poker.Hearts && s != poker.Diamonds && s != poker.Clubs {
			continue
		}
		rMap[c.Rank()]++
	}

	var ans [][]poker.Card

	for r, cnt := range rMap {
		if cnt > 2 {
			ans = append(ans, deck.PickRank(r))
		}
		if cnt == 2 {
			c, err := wilds.Pop()
			if err != nil {
				continue
			}
			tmp := deck.PickRank(r)
			tmp = append(tmp, c)
			ans = append(ans, tmp)
		}
	}

	return ans, len(ans) > 0
}

// Refresher 整牌器 (整理手牌)
type Refresher interface {
	// 整理手牌
	//  @param handLength 規定手牌數
	//  @param jokerR 指定為鬼牌的點數
	//  @param hands 手牌
	Refresh(handLength int, jokerR poker.Rank, hands ...*poker.Deck) (bool, poker.Card, []*poker.Deck)
}

type refresher struct{}

func (r *refresher) Refresh(handLength int, jokerR poker.Rank, hands ...*poker.Deck) (bool, poker.Card, []*poker.Deck) {
	wilds := poker.NewEmptyDeck()

	deck := poker.NewEmptyDeck()
	deck.AppendDeck(hands...)
	length := deck.Len()

	// 先把百搭牌全挑出
	wilds.AppendCard(deck.PickSuit(poker.Jokers)...)

	deck.Sort(poker.SortSameRank)

	var newhands []*poker.Deck

	{ // 1st Life
		if ans, ok := FirstLifePicker().Pick(deck.Cards); ok {
			for _, a := range ans {
				d := poker.NewEmptyDeck()
				d.AppendCard(a...)
				newhands = append(newhands, d)
				deck.PickCards(a...)
			}
		}
	}

	// wilds.AppendCard(deck.PickRank(jokerR)...)
	deck.AppendDeck(wilds)

	{ // 2st Life
		if ans, ok := SecondLifePicker(jokerR).Pick(deck.Cards); ok {
			for _, a := range ans {
				d := poker.NewEmptyDeck()
				d.AppendCard(a...)
				newhands = append(newhands, d)
				deck.PickCards(a...)
			}
		}
	}

	{ // Set
		if ans, ok := SetPicker(jokerR).Pick(deck.Cards); ok {
			for _, a := range ans {
				d := poker.NewEmptyDeck()
				d.AppendCard(a...)
				newhands = append(newhands, d)
				deck.PickCards(a...)
			}
		}
	}

	if deck.Len() > 0 {
		newhands = append(newhands, deck)
	}

	if length == handLength {
		return false, 0, newhands
	}

	// 檢查可否胡牌
	cMap := make(map[Category][]*poker.Deck)
	{ // 整理每個deck的組合牌型
		for _, d := range newhands {
			if d.Len() == 0 {
				continue
			}
			if ok := FirstLifeValidator().Validate(d.Cards); ok {
				cMap[FirstLife] = append(cMap[FirstLife], d)
				continue
			}
			if ok := SecondLifeValidator(jokerR).Validate(d.Cards); ok {
				cMap[SecondLife] = append(cMap[SecondLife], d)
				continue
			}
			if ok := SetValidator(jokerR).Validate(d.Cards); ok {
				cMap[Set] = append(cMap[Set], d)
				continue
			}
			cMap[None] = append(cMap[None], d)
		}
	}

	var throw poker.Card

	var findOver3 = func(decks []*poker.Deck) poker.Card {
		for _, d := range decks {
			if d.Len() > 3 {
				throw, _ = d.Pop()
				return throw
			}
		}
		return 0
	}

	var clean = func(decks []*poker.Deck) []*poker.Deck {
		var ans []*poker.Deck
		for _, d := range decks {
			if d.Len() > 0 {
				ans = append(ans, d)
			}
		}
		return ans
	}

	if len(cMap[None]) == 0 { //  沒有亂牌 挑一組超過三張的牌組打
		throw = findOver3(newhands)
	} else {
		d := cMap[None][0]
		noneLength := d.Len()
		if cards := d.PickSuit(poker.Jokers); len(cards) > 0 {
			if len(cMap[SecondLife]) > 0 {
				cMap[SecondLife][0].AppendCard(cards...)
			} else if len(cMap[Set]) > 0 {
				cMap[Set][0].AppendCard(cards...)
			}
		}
		if cards := d.PickRank(jokerR); len(cards) > 0 {
			if len(cMap[SecondLife]) > 0 {
				cMap[SecondLife][0].AppendCard(cards...)
			} else if len(cMap[Set]) > 0 {
				cMap[Set][0].AppendCard(cards...)
			}
		}

		if noneLength >= 3 && d.Len() == 0 {
			throw = findOver3(newhands)
		} else {
			if d.Len() > 0 {
				throw, _ = d.Pop()
			} else {
				throw = findOver3(newhands)
			}
		}
	}

	newhands = clean(newhands)

	pr := NewPointResult(jokerR, newhands...)

	return pr.IsValidDeclaration(), throw, newhands
}

// DefaultRefresher 預設整牌器
//  先把 [鬼牌][隨機鬼牌] 抽出
//  整理出 1st Life
//  整理出 2st Life
//  整理出 Set
//  剩餘的 牌堆加入 [鬼牌][隨機鬼牌]
//  檢查是否可胡牌
//  可胡牌則打出超過3張牌的牌堆
//  不可胡牌則打出亂牌的第一張 (排除[鬼牌][隨機鬼牌])
func DefaultRefresher() Refresher { return new(refresher) }
