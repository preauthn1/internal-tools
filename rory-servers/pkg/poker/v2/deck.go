package poker

import "fmt"

// Deck 牌堆
type Deck struct {
	Cards []Card
}

// NewEmptyDeck 建立空牌堆
func NewEmptyDeck(cards ...Card) *Deck {
	d := new(Deck)
	d.AppendCard(cards...)
	return d
}

// NewDeck 建立牌堆
//  @ param c int 需要幾副牌()
//  * 基本牌堆為一副牌
//  * 一副牌預設為 [黑桃 紅心 方塊 梅花] 四個花色
//  * 每個花色有 [A 2 3 4 5 6 7 8 9 10 J Q K]
//  * 有需要鬼牌的遊戲自行加入
func NewDeck(c int) *Deck {
	deck := new(Deck)
	for i := 0; i < c; i++ {
		deck.Cards = append(deck.Cards, simpleCards...)
	}
	return deck
}

func (d *Deck) Len() int {
	return len(d.Cards)
}

func (d *Deck) Swap(i, j int) {
	d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i]
}

func (d *Deck) Less(i, j int) bool {
	if d.Cards[i].Suit() == d.Cards[j].Suit() {
		return d.Cards[i].Rank() < d.Cards[j].Rank()
	}
	return d.Cards[i].Suit() < d.Cards[j].Suit()
}

func (d *Deck) String() string {
	return fmt.Sprintf("%v", d.Cards)
}

// Sort 排序
func (d *Deck) Sort(p SortPolicy) {
	p.Sort(d)
}

// Range 尋遍每張牌
func (d *Deck) Range(f func(int, Card) bool) {
	for i, c := range d.Cards {
		if !f(i, c) {
			break
		}
	}
}

// Clean 清空卡牌
func (d *Deck) Clean() {
	d.Cards = []Card(nil)
}

// Exist 檢查卡牌是否存在此牌堆中
func (d *Deck) Exist(s Suit, r Rank) (int, bool) {
	c := Card(uint8(s) + uint8(r))
	return d.ExistCard(c)
}

// ExistCard 檢查卡牌是否存在此牌堆中
func (d *Deck) ExistCard(c Card) (idx int, ok bool) {

	idx = -1
	ok = false

	d.Range(func(i int, card Card) bool {
		if card == c {
			idx = i
			ok = true
			return false
		}
		return true
	})

	return
}

// Append 以花色點數新增卡牌在最後面
func (d *Deck) Append(s Suit, r Rank) {
	d.AppendCard(Card(uint8(s) + uint8(r)))
}

// AppendCard 新增卡牌在最後面
func (d *Deck) AppendCard(cards ...Card) {
	if len(cards) == 0 {
		return
	}
	d.Cards = append(d.Cards, cards...)
}

// AppendDeck 新增卡牌在最後面
func (d *Deck) AppendDeck(decks ...*Deck) {
	if len(decks) == 0 {
		return
	}
	for _, deck := range decks {
		if deck.Len() == 0 {
			continue
		}
		d.Cards = append(d.Cards, deck.Cards...)
	}
}

// Stack 以花色點數新增卡牌在最前面
func (d *Deck) Stack(s Suit, r Rank) {
	d.AppendCard(Card(uint8(s) + uint8(r)))
}

// StackCard 新增卡牌在最前面
func (d *Deck) StackCard(c Card) {
	d.Cards = append([]Card{c}, d.Cards...)
}

// remove 移除卡牌
func (d *Deck) remove(i int) error {
	if d.Len() < i {
		return ErrCardNotEnough
	}

	d.Cards = append(d.Cards[:i], d.Cards[i+1:]...)

	return nil
}

// Pick 抽出一張指定牌
func (d *Deck) Pick(s Suit, r Rank) (Card, error) {
	c := Card(uint8(s) + uint8(r))
	return d.PickCard(c)
}

// PickCard 抽出一張指定牌
func (d *Deck) PickCard(c Card) (Card, error) {

	i, ok := d.ExistCard(c)
	if !ok {
		return 0, ErrCardNotFound
	}

	d.remove(i)

	return c, nil
}

// PickCards 抽出指定牌
func (d *Deck) PickCards(cards ...Card) []Card {

	var picked []Card

	for _, c := range cards {
		if _, err := d.PickCard(c); err == nil {
			picked = append(picked, c)
		}
	}

	return picked
}

// Pop 從最上面抽一張牌
func (d *Deck) Pop() (Card, error) {
	if d.Len() == 0 {
		return 0, ErrCardNotEnough
	}

	c := d.Cards[0]

	d.remove(0)

	return c, nil
}

// CountOf 查詢指定花色點數共有幾張牌
func (d *Deck) CountOf(s Suit, r Rank) int {
	c := Card(uint8(s) + uint8(r))
	return d.CountOfCard(c)
}

// CountOfCard 查詢指定卡牌有幾張牌
func (d *Deck) CountOfCard(c Card) int {
	var cnt int
	d.Range(func(i int, card Card) bool {
		if card == c {
			cnt++
		}
		return true
	})
	return cnt
}

// CountOfSuit 查詢指定花色共有幾張牌
func (d *Deck) CountOfSuit(s Suit) int {
	var cnt int
	d.Range(func(i int, card Card) bool {
		if card.Suit() == s {
			cnt++
		}
		return true
	})
	return cnt
}

// CountOfRank 查詢指定點數共有幾張牌
func (d *Deck) CountOfRank(r Rank) int {
	var cnt int
	d.Range(func(i int, card Card) bool {
		switch card.Suit() {
		default:
			return true
		case Spades, Hearts, Diamonds, Clubs:
		}
		if card.Rank() == r {
			cnt++
		}
		return true
	})
	return cnt
}

// PickSuit 抽出指定花色的所有牌
func (d *Deck) PickSuit(s Suit) []Card {
	var cards []Card

	// 先找出哪些同花色的牌
	d.Range(func(i int, c Card) bool {
		if c.Suit() == s {
			cards = append(cards, c)
		}
		return true
	})

	// 再把這些牌移除
	for _, c := range cards {
		d.PickCard(c)
	}

	return cards
}

// PickRank 抽出指定點數的所有牌
func (d *Deck) PickRank(r Rank) []Card {
	var cards []Card

	// 先找出哪些同點數的牌
	d.Range(func(i int, card Card) bool {
		switch card.Suit() {
		default:
			return true
		case Spades, Hearts, Diamonds, Clubs:
		}
		if card.Rank() == r {
			cards = append(cards, card)
		}
		return true
	})

	// 再把這些牌移除
	for _, c := range cards {
		d.PickCard(c)
	}

	return cards
}
