package poker

import (
	"bytes"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// NoShufflePoker 撲克
type NoShufflePoker struct {
	Cards []*Card
}

// NewNoShufflePoker ?
func NewNoShufflePoker(hasKing bool, hasPizi bool, count int, groupPoker map[int]int) *NoShufflePoker {
	poker := new(NoShufflePoker)
	rand := util.RndInt(0, len(groupPoker))
	group := 4 - groupPoker[rand] //設定幾張一組
	arrCards := []*Card{}
	mcards := make(map[int]*Card)
	mapGroupCard := make(map[int][]*Card)

	for ip := 0; ip < count; ip++ {

		var idx int32 = 2
		for i := idx; i <= 14; i++ {
			var card *Card
			var groupCard []*Card
			//黑桃
			card = NewPokerCard(PokerColorHeiTao, i)
			if group < int(1) && i != int32(2) && i != int32(14) {
				groupCard = append(groupCard, card)
			} else {
				arrCards = append(arrCards, card)
			}

			//红桃
			card = NewPokerCard(PokerColorHongTao, i)
			if group < int(2) && i != int32(2) && i != int32(14) {
				groupCard = append(groupCard, card)
			} else {
				arrCards = append(arrCards, card)
			}

			//梅花
			card = NewPokerCard(PokerColorMeiHua, i)
			if group < int(3) && i != int32(2) && i != int32(14) {
				groupCard = append(groupCard, card)
			} else {
				arrCards = append(arrCards, card)
			}

			//方块
			card = NewPokerCard(PokerColorFangKuai, i)
			if group < int(4) && i != int32(2) && i != int32(14) {
				groupCard = append(groupCard, card)
			} else {
				arrCards = append(arrCards, card)
			}

			if len(groupCard) != int(0) {
				mapGroupCard[int(i)] = groupCard
			}
		}

		//是否有王
		if hasKing {
			skcard := NewPokerCard(PokerColorSmallKing, PokerNumberSmallKing)
			bkcard := NewPokerCard(PokerColorBigKing, PokerNumberBigKing)

			mapGroupCard[int(PokerNumberSmallKing)] = []*Card{skcard}
			mapGroupCard[int(PokerNumberBigKing)] = []*Card{bkcard}
		}

		//是否有皮子
		if hasPizi {
			pzcard := NewPokerCard(PokerColorPiZi, PokerNumberPiZi)
			mapGroupCard[int(PokerNumberBigKing)] = []*Card{pzcard}
		}

		//不在牌組的隨機發
		mapGroupLen := len(mapGroupCard) + 5
		for k, c := range arrCards {
			if _, ok := mapGroupCard[k+mapGroupLen]; ok {
				mlog.Error("洗牌位置已存在")
			}
			mapGroupCard[k+mapGroupLen] = []*Card{c}
		}

		for _, c := range mcards {
			poker.Cards = append(poker.Cards, c)
		}

		for _, c := range mapGroupCard {
			poker.Cards = append(poker.Cards, c...)
		}
	}

	if len(poker.Cards) != int(54) {
		mlog.Error("初始牌張數錯誤, 非54張, 實際:%v", len(poker.Cards))
	}

	return poker
}

// Shuffle 洗牌
func (p *NoShufflePoker) Shuffle() {

}

// UserCardsLen ?
func (p *NoShufflePoker) UserCardsLen(uid int32) int {
	return len(p.Cards)
}

// ShuffleOnce 洗牌一次
func (p *NoShufflePoker) ShuffleOnce() {

}

// ShuffleMore 洗多次牌
func (p *NoShufflePoker) ShuffleMore(times int) {

}

// Deal 發牌
func (p *NoShufflePoker) Deal(userid int32, cn int) *CardGroup {
	if len(p.Cards) < cn {
		return nil
	}

	ret := p.Cards[:cn]
	p.Cards = p.Cards[cn:]

	return NewCardGroup(ret)
}

// AddCardGroup 加牌
func (p *NoShufflePoker) AddCardGroup(userid int32, cg *CardGroup) {

}

// AddCard 加牌
func (p *NoShufflePoker) AddCard(userid int32, card *Card) {

}

// TakeCardNumber 指定取牌號碼
func (p *NoShufflePoker) TakeCardNumber(userid int32, cardNumber int32) (bool, Card) {
	return false, Card{}
}

// TakeDiffCardNumber 指定取非存在的號碼
func (p *NoShufflePoker) TakeDiffCardNumber(userid int32, diffNums []int32) (bool, Card) {
	return false, Card{}
}

// ExChange 換牌
func (p *NoShufflePoker) ExChange(userid int32, otheruserid int32) {

}

// 轉換為字符串
func (p *NoShufflePoker) String() string {
	var buf bytes.Buffer
	for _, c := range p.Cards {
		buf.WriteString(c.String() + ",")
	}

	return buf.String()
}

// TakeCardColor ?
func (p *NoShufflePoker) TakeCardColor(userid int32, color int32) (bool, Card) {
	for k, v := range p.Cards {
		if v.Color == color {
			p.Cards = RemoveIndex(p.Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}

// TakeCardNumberAndColor ?
func (p *NoShufflePoker) TakeCardNumberAndColor(userid int32, cardNumber int32, color int32) (bool, Card) {
	for k, v := range p.Cards {
		if v.CardNumber == cardNumber && v.Color == color {
			p.Cards = RemoveIndex(p.Cards, k)
			return true, *v
		}
	}

	return false, Card{}
}
