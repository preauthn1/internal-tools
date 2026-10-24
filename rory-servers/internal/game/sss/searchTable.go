package sss

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// // 查詢牌型用表(只給同花順, 四條, 葫蘆, 同花, 順子 的5張牌型用)
// type HandsTable struct {
// 	searchTable map[netproto.SSSDunCardsNomalType][]*WeightTable
// }
var (
	searchTable            map[netproto.SSSDunCardsNomalType][]*WeightTable
	twoTongHuaShunPatterns []*poker.CardGroup // 同時兩個同花順牌型
	twoTongHuaPatterns     []*poker.CardGroup // 同時兩個同花牌型
	twoShunZiPatterns      []*poker.CardGroup // 同時兩個順子牌型
)

func init() {
	InitSearchTable()
}

// https://example.invalid/wiki/%E6%92%B2%E5%85%8B%E7%89%8C%E5%9E%8B
// 牌型		組合數		機率
// 同花順	40			0.00154%
// 四條		624			0.0240%
// 葫蘆		3,744		0.144%
// 同花		5,108		0.197%
// 順子		10,200		0.392%
// 三條		54,912		2.11%
// 兩對		123,552		4.75%
// 一對		1,098,240	42.26%
// 散牌		1,302,540	50.12%
// 總計		2,598,960	100.00%

// 同花順 	cardTypeWeight uint32 = 1 << (16+5)
// 四條		cardTypeWeight uint32 = 1 << (16+4)
// 葫蘆		cardTypeWeight uint32 = 1 << (16+3)
// 同花		cardTypeWeight uint32 = 1 << (16+2)
// 順子		cardTypeWeight uint32 = 1 << (16+1)
// 同花順所有組合和比較權重 (拿 16 bit以上做牌型比較) cardTypeWeight + cardValueWeight
func TongHuaShunPatterns() []*WeightTable {
	if len(searchTable[netproto.SSSDunCardsNomalType_TongHuaShun]) > 0 {
		return searchTable[netproto.SSSDunCardsNomalType_TongHuaShun]
	}

	tables := make([]*WeightTable, 0)
	var cardTypeWeight uint32 = 1 << (16 + 5)
	var cardValueWeight uint32 = 0
	// fmt.Println(fmt.Sprintf("%b", cardTypeWeight))

	// 產生其他可能
	for headCardNumber := 1; headCardNumber <= 10; headCardNumber++ {

		//梅花
		cards := new(poker.CardGroup)
		for i := 0; i < 5; i++ {
			number := headCardNumber + i
			if number == 1 {
				number = 14
			}
			cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
		}
		cardValueWeight++
		tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

		// 方块
		cards = new(poker.CardGroup)
		for i := 0; i < 5; i++ {
			number := headCardNumber + i
			if number == 1 {
				number = 14
			}
			cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
		}
		cardValueWeight++
		tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

		// 红桃
		cards = new(poker.CardGroup)
		for i := 0; i < 5; i++ {
			number := headCardNumber + i
			if number == 1 {
				number = 14
			}
			cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
		}
		cardValueWeight++
		tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

		// 黑桃
		cards = new(poker.CardGroup)
		for i := 0; i < 5; i++ {
			number := headCardNumber + i
			if number == 1 {
				number = 14
			}
			cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
		}
		cardValueWeight++
		tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})
	}

	return tables
}

// 同花順 	cardTypeWeight uint32 = 1 << (16+5)
// 四條		cardTypeWeight uint32 = 1 << (16+4)
// 葫蘆		cardTypeWeight uint32 = 1 << (16+3)
// 同花		cardTypeWeight uint32 = 1 << (16+2)
// 順子		cardTypeWeight uint32 = 1 << (16+1)
// 四條所有組合和比較權重 (拿 16 bit以上做牌型比較) cardTypeWeight + cardValueWeight
func ZhaDanPatterns() []*WeightTable {
	if len(searchTable[netproto.SSSDunCardsNomalType_ZhaDan]) > 0 {
		return searchTable[netproto.SSSDunCardsNomalType_ZhaDan]
	}

	tables := make([]*WeightTable, 0)
	var cardTypeWeight uint32 = 1 << (16 + 4)
	var cardValueWeight uint32 = 0
	// fmt.Println(fmt.Sprintf("%b", cardTypeWeight))

	// 產生可能
	for number := 2; number <= 14; number++ {

		// remainCards
		oriCards := poker.NewCardGroup(poker.NewRandomPoker(false, false, 1).Cards)

		oriCards.DeleteCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
		oriCards.DeleteCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
		oriCards.DeleteCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
		oriCards.DeleteCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))

		// 從小排到大
		sort.SliceStable(oriCards.Cards, func(i, j int) bool {
			if oriCards.Cards[i].CardNumber < oriCards.Cards[j].CardNumber {
				return true
			} else if oriCards.Cards[i].CardNumber == oriCards.Cards[j].CardNumber {
				// 相同點數比較花色
				if oriCards.Cards[i].IsMeiHua() {
					return true
				}

				if oriCards.Cards[i].IsFangKuai() && !oriCards.Cards[j].IsMeiHua() {
					return true
				}

				if oriCards.Cards[i].IsHongTao() && oriCards.Cards[j].IsHeiTao() {
					return true
				}
			}

			return false
		})

		// fmt.Println(oriCards)
		for _, v := range oriCards.Cards {
			cards := new(poker.CardGroup)
			cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
			cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
			cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
			cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
			cards.AddCard(v)
			cardValueWeight++
			tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})
		}
	}

	return tables
}

// 同花順 	cardTypeWeight uint32 = 1 << (16+5)
// 四條		cardTypeWeight uint32 = 1 << (16+4)
// 葫蘆		cardTypeWeight uint32 = 1 << (16+3)
// 同花		cardTypeWeight uint32 = 1 << (16+2)
// 順子		cardTypeWeight uint32 = 1 << (16+1)
// 葫蘆所有組合和比較權重 (拿 16 bit以上做牌型比較) cardTypeWeight + cardValueWeight
func HuluPatterns() []*WeightTable {
	if len(searchTable[netproto.SSSDunCardsNomalType_HuLu]) > 0 {
		return searchTable[netproto.SSSDunCardsNomalType_HuLu]
	}

	tables := make([]*WeightTable, 0)
	var cardTypeWeight uint32 = 1 << (16 + 3)
	var cardValueWeight uint32 = 0
	// fmt.Println(fmt.Sprintf("%b", cardTypeWeight))

	// 產生可能
	for number := 2; number <= 14; number++ {

		for twoPairNumber := 2; twoPairNumber <= 14; twoPairNumber++ { // 取一個數字做對
			if twoPairNumber != number {
				// 三條不取黑桃
				// 一對取梅花方塊
				cards := new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取紅心黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 三條不取紅心
				// 一對取梅花方塊
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取紅心黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 三條不取方塊
				// 一對取梅花方塊
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取紅心黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 三條不取梅花
				// 一對取梅花方塊
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取梅花黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorMeiHua, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊紅心
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取方塊黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

				// 一對取紅心黑桃
				cards = new(poker.CardGroup)
				cards.AddCard(poker.NewPokerCard(poker.PokerColorFangKuai, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(number)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHongTao, int32(twoPairNumber)))
				cards.AddCard(poker.NewPokerCard(poker.PokerColorHeiTao, int32(twoPairNumber)))
				cardValueWeight++
				tables = append(tables, &WeightTable{cards, cardTypeWeight + cardValueWeight})

			}
		}
	}

	return tables
}

// 同花順 	cardTypeWeight uint32 = 1 << (16+5)
// 四條		cardTypeWeight uint32 = 1 << (16+4)
// 葫蘆		cardTypeWeight uint32 = 1 << (16+3)
// 同花		cardTypeWeight uint32 = 1 << (16+2)
// 順子		cardTypeWeight uint32 = 1 << (16+1)
// 同花所有組合和比較權重 (拿 16 bit以上做牌型比較) cardTypeWeight + cardValueWeight
func TongHuaPatterns() []*WeightTable {
	if len(searchTable[netproto.SSSDunCardsNomalType_TongHua]) > 0 {
		return searchTable[netproto.SSSDunCardsNomalType_TongHua]
	}
	tables := make([]*WeightTable, 0)
	var cardTypeWeight uint32 = 1 << (16 + 2)
	var cardValueWeight uint32 = 0
	// fmt.Println(fmt.Sprintf("%b", cardTypeWeight))

	for number := 2; number <= 14; number++ {
		for number2 := number + 1; number2 <= 14; number2++ {
			for number3 := number2 + 1; number3 <= 14; number3++ {
				for number4 := number3 + 1; number4 <= 14; number4++ {
					for number5 := number4 + 1; number5 <= 14; number5++ {
						// 梅花
						color := poker.PokerColorMeiHua
						cards := new(poker.CardGroup)
						cards.AddCard(poker.NewPokerCard(color, int32(number)))
						cards.AddCard(poker.NewPokerCard(color, int32(number2)))
						cards.AddCard(poker.NewPokerCard(color, int32(number3)))
						cards.AddCard(poker.NewPokerCard(color, int32(number4)))
						cards.AddCard(poker.NewPokerCard(color, int32(number5)))
						if !cards.IsShunzi() {
							// cardValueWeight++
							tables = append(tables, &WeightTable{cards, 0})

							// 方塊
							color = poker.PokerColorFangKuai
							cards = new(poker.CardGroup)
							cards.AddCard(poker.NewPokerCard(color, int32(number)))
							cards.AddCard(poker.NewPokerCard(color, int32(number2)))
							cards.AddCard(poker.NewPokerCard(color, int32(number3)))
							cards.AddCard(poker.NewPokerCard(color, int32(number4)))
							cards.AddCard(poker.NewPokerCard(color, int32(number5)))
							// cardValueWeight++
							tables = append(tables, &WeightTable{cards, 0})

							// 紅心
							color = poker.PokerColorHongTao
							cards = new(poker.CardGroup)
							cards.AddCard(poker.NewPokerCard(color, int32(number)))
							cards.AddCard(poker.NewPokerCard(color, int32(number2)))
							cards.AddCard(poker.NewPokerCard(color, int32(number3)))
							cards.AddCard(poker.NewPokerCard(color, int32(number4)))
							cards.AddCard(poker.NewPokerCard(color, int32(number5)))
							// cardValueWeight++
							tables = append(tables, &WeightTable{cards, 0})

							// 黑桃
							color = poker.PokerColorHeiTao
							cards = new(poker.CardGroup)
							cards.AddCard(poker.NewPokerCard(color, int32(number)))
							cards.AddCard(poker.NewPokerCard(color, int32(number2)))
							cards.AddCard(poker.NewPokerCard(color, int32(number3)))
							cards.AddCard(poker.NewPokerCard(color, int32(number4)))
							cards.AddCard(poker.NewPokerCard(color, int32(number5)))
							// cardValueWeight++
							tables = append(tables, &WeightTable{cards, 0})
						}
					}
				}
			}
		}
	}

	// 同花牌組 由小到大排序
	sort.SliceStable(tables, func(i, j int) bool {
		if tables[i].Cards.Cards[4].CardNumber < tables[j].Cards.Cards[4].CardNumber {
			return true
		} else if tables[i].Cards.Cards[4].CardNumber == tables[j].Cards.Cards[4].CardNumber {
			if tables[i].Cards.Cards[3].CardNumber < tables[j].Cards.Cards[3].CardNumber {
				return true
			} else if tables[i].Cards.Cards[3].CardNumber == tables[j].Cards.Cards[3].CardNumber {
				if tables[i].Cards.Cards[2].CardNumber < tables[j].Cards.Cards[2].CardNumber {
					return true
				} else if tables[i].Cards.Cards[2].CardNumber == tables[j].Cards.Cards[2].CardNumber {
					if tables[i].Cards.Cards[1].CardNumber < tables[j].Cards.Cards[1].CardNumber {
						return true
					} else if tables[i].Cards.Cards[1].CardNumber == tables[j].Cards.Cards[1].CardNumber {
						if tables[i].Cards.Cards[0].CardNumber < tables[j].Cards.Cards[0].CardNumber {
							return true
						}
					}
				}
			}
		}
		return false
	})

	// 給予比較權重
	cardValueWeight = 0
	for _, v := range tables {
		cardTypeWeight++
		v.Weight = cardTypeWeight + cardValueWeight
	}

	return tables
}

// 同花順 	cardTypeWeight uint32 = 1 << (16+5)
// 四條		cardTypeWeight uint32 = 1 << (16+4)
// 葫蘆		cardTypeWeight uint32 = 1 << (16+3)
// 同花		cardTypeWeight uint32 = 1 << (16+2)
// 順子		cardTypeWeight uint32 = 1 << (16+1)
// 順子所有組合和比較權重 (拿 16 bit以上做牌型比較) cardTypeWeight + cardValueWeight
func ShunZiPatterns() []*WeightTable { //沒想到好作法，先這樣
	if len(searchTable[netproto.SSSDunCardsNomalType_ShunZi]) > 0 {
		return searchTable[netproto.SSSDunCardsNomalType_ShunZi]
	}
	tables := make([]*WeightTable, 0)
	var cardTypeWeight uint32 = 1 << (16 + 1)
	var cardValueWeight uint32 = 0
	// fmt.Println(fmt.Sprintf("%b", cardTypeWeight))

	// 產生可能
	for number := 2; number <= 14; number++ {
		for number2 := number + 1; number2 <= 14; number2++ {
			for number3 := number2 + 1; number3 <= 14; number3++ {
				for number4 := number3 + 1; number4 <= 14; number4++ {
					for number5 := number4 + 1; number5 <= 14; number5++ {
						// 決定五張牌的花色
						colors := []int32{
							poker.PokerColorMeiHua,
							poker.PokerColorFangKuai,
							poker.PokerColorHongTao,
							poker.PokerColorHeiTao,
						}

						for c1 := 0; c1 < len(colors); c1++ {
							for c2 := 0; c2 < len(colors); c2++ {
								for c3 := 0; c3 < len(colors); c3++ {
									for c4 := 0; c4 < len(colors); c4++ {
										for c5 := 0; c5 < len(colors); c5++ {
											cards := new(poker.CardGroup)
											cards.AddCard(poker.NewPokerCard(colors[c1], int32(number)))
											cards.AddCard(poker.NewPokerCard(colors[c2], int32(number2)))
											cards.AddCard(poker.NewPokerCard(colors[c3], int32(number3)))
											cards.AddCard(poker.NewPokerCard(colors[c4], int32(number4)))
											cards.AddCard(poker.NewPokerCard(colors[c5], int32(number5)))
											if !cards.IsTonghua() && cards.IsShunzi() {
												tables = append(tables, &WeightTable{cards, 0})
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	sort.SliceStable(tables, func(i, j int) bool {
		group1 := tables[i].Cards
		group2 := tables[j].Cards
		return group1.GetMaxCardNumber() < group2.GetMaxCardNumber()
		// group1.SortMinFirst()
		// group2.SortMinFirst()
		// for _, v1 := range group1.Cards {
		// 	for _, v2 := range group2.Cards {
		// 		if v1.CardNumber < v2.CardNumber {
		// 			return true
		// 		} else if v1.CardNumber > v2.CardNumber {
		// 			return false
		// 		}
		// 	}
		// }
		// return false
	})

	// 給予比較權重
	cardValueWeight = 0
	for _, v := range tables {
		cardTypeWeight++
		v.Weight = cardTypeWeight + cardValueWeight
	}

	return tables
}

// CheckSanTongHuaShunPatterns 確認是否為三同花順組合(窮舉確認所有組合)
func CheckSanTongHuaShunPatterns(cg *poker.CardGroup) bool {

	for _, v := range twoTongHuaShunPatterns {
		oriCards := poker.NewCardGroup(cg.Cards)
		oriCards.DeleteGroup(v)

		if oriCards.IsShunzi() && oriCards.IsTonghua() {
			return true
		}
	}

	return false
}

// CheckSanTongHuaPatterns 確認是否為三同花組合(窮舉確認所有組合)
func CheckSanTongHuaPatterns(cg *poker.CardGroup) bool {

	for _, v := range twoTongHuaPatterns {
		oriCards := poker.NewCardGroup(cg.Cards)
		oriCards.DeleteGroup(v)

		if oriCards.IsTonghua() {
			return true
		}
	}

	return false
}

// CheckSanShunZiPatterns 確認是否為三順子組合(窮舉確認所有組合)
func CheckSanShunZiPatterns(cg *poker.CardGroup) bool {

	for _, v := range twoShunZiPatterns {
		oriCards := poker.NewCardGroup(cg.Cards)
		oriCards.DeleteGroup(v)

		if oriCards.IsShunzi() {
			return true
		}
	}

	return false
}

func GetWeightByDun(dun *netproto.Dun) uint32 {
	cards := new(poker.CardGroup)
	for _, v := range dun.Cards {
		c := new(poker.Card)
		c.CardNumber = v.GetNumber()
		c.Color = v.GetColor()
		cards.AddCard(c)
	}
	return GetWeight(cards)
}

// GetWeight 取得5張牌的牌型跟權重(目前只做同花順, 四條, 葫蘆, 同花, 順子 的5張牌型)
func GetWeight(cards *poker.CardGroup) uint32 {
	if len(cards.Cards) != 5 {
		return 0
	}

	// 同花順
	for _, v := range searchTable[netproto.SSSDunCardsNomalType_TongHuaShun] {
		if cards.IsEqual(v.Cards) {
			return v.Weight
		}
	}

	// 四條
	for _, v := range searchTable[netproto.SSSDunCardsNomalType_ZhaDan] {
		if cards.IsEqual(v.Cards) {
			return v.Weight
		}
	}

	// 葫蘆
	for _, v := range searchTable[netproto.SSSDunCardsNomalType_HuLu] {
		if cards.IsEqual(v.Cards) {
			return v.Weight
		}
	}

	// 同花
	for _, v := range searchTable[netproto.SSSDunCardsNomalType_TongHua] {
		if cards.IsEqual(v.Cards) {
			return v.Weight
		}
	}

	// 順子
	for _, v := range searchTable[netproto.SSSDunCardsNomalType_ShunZi] {
		if cards.IsEqual(v.Cards) {
			return v.Weight
		}
	}
	return 0
}

// 同花順	40			0.00154%
// 四條		624			0.0240%
// 葫蘆		3,744		0.144%
// 同花		5,108		0.197%
// 順子		10,200		0.392%
func InitSearchTable() {
	mlog.System("等待初始化牌组比较权重表")

	searchTable = make(map[netproto.SSSDunCardsNomalType][]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_TongHuaShun] = make([]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_TongHuaShun] = TongHuaShunPatterns() // 40 種組合
	mlog.System("同花順組合表")

	searchTable[netproto.SSSDunCardsNomalType_ZhaDan] = make([]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_ZhaDan] = ZhaDanPatterns() // 624 種組合
	mlog.System("炸彈組合表")

	searchTable[netproto.SSSDunCardsNomalType_HuLu] = make([]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_HuLu] = HuluPatterns() // 3744 種組合
	mlog.System("葫蘆組合表")

	searchTable[netproto.SSSDunCardsNomalType_TongHua] = make([]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_TongHua] = TongHuaPatterns() // 5108 種組合
	mlog.System("同花組合表")

	searchTable[netproto.SSSDunCardsNomalType_ShunZi] = make([]*WeightTable, 0)
	searchTable[netproto.SSSDunCardsNomalType_ShunZi] = ShunZiPatterns() // 10200 種組合
	mlog.System("順子組合表")

	// 特殊報到牌計算用
	twoTongHuaShunPatterns = make([]*poker.CardGroup, 0) // 同時兩個同花順牌型
	tongHuaShunPatterns := TongHuaShunPatterns()
	patternLen := len(tongHuaShunPatterns)
	for i := 0; i < patternLen; i++ {
		for j := i + 1; j < patternLen; j++ {
			group1 := tongHuaShunPatterns[i].Cards
			group2 := tongHuaShunPatterns[j].Cards
			if !group1.HaveSameCards(group2) {
				combineGroup := poker.NewCardGroup(group1.Cards)
				combineGroup.AddCardGroup(group2)
				twoTongHuaShunPatterns = append(twoTongHuaShunPatterns, combineGroup)
			}
		}
	}
	mlog.System("同時兩個同花順牌型表")

	// twoTongHuaPatterns = make([]*poker.CardGroup, 0) // 同時兩個同花牌型
	// tongHuaPatterns := TongHuaPatterns()
	// patternLen = len(tongHuaPatterns)
	// for i := 0; i < patternLen; i++ {
	// 	for j := i + 1; j < patternLen; j++ {
	// 		group1 := tongHuaPatterns[i].Cards
	// 		group2 := tongHuaPatterns[j].Cards
	// 		if !group1.HaveSameCards(group2) {
	// 			// combineGroup := poker.NewCardGroup(group1.Cards)
	// 			// combineGroup.AddCardGroup(group2)
	// 			twoTongHuaPatterns = append(twoTongHuaPatterns, group1)
	// 			twoTongHuaPatterns = append(twoTongHuaPatterns, group2)
	// 		}
	// 	}
	// }
	// mlog.System("同時兩個同花牌型")

	// twoShunZiPatterns = make([]*poker.CardGroup, 0) // 同時兩個順子牌型
	// shunZiPatterns := ShunZiPatterns()
	// patternLen = len(shunZiPatterns)
	// for i := 0; i < patternLen; i++ {
	// 	for j := i + 1; j < patternLen; j++ {
	// 		group1 := shunZiPatterns[i].Cards
	// 		group2 := shunZiPatterns[j].Cards
	// 		if !group1.HaveSameCards(group2) {
	// 			// combineGroup := poker.NewCardGroup(group1.Cards)
	// 			// combineGroup.AddCardGroup(group2)
	// 			twoShunZiPatterns = append(twoShunZiPatterns, group1)
	// 			twoShunZiPatterns = append(twoShunZiPatterns, group2)
	// 		}
	// 	}
	// }
	// mlog.System("同時兩個順子牌型")

	mlog.System("初始化权重表完成")
}

// 權重表
type WeightTable struct {
	// 牌組(由小到大)
	Cards *poker.CardGroup
	// 比較權重
	Weight uint32
}
