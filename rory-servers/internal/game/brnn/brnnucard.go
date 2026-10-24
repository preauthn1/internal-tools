package brnn

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	nnCardNumber_PiZi      = 100000 //牛牛牌点,皮子
	nnCardNumber_BingKing  = 10000  //牛牛牌点,大王
	nnCardNumber_SmallKing = 9999   //牛牛牌点,小王
	nnCardNumber_A         = 1      //牛牛牌点,A
)

//比较牌型,返回玩家1 赢玩家2的倍数
func compareCards(ret1 *niuniuCardRet, ret2 *niuniuCardRet) int32 {
	if ret1.niuRet != ret2.niuRet {
		if ret1.niuRet > ret2.niuRet {
			return -1 * ret2.getBetTimes()
		} else {
			return ret1.getBetTimes()
		}
	}

	//有牛
	if ret1.niuRet == nnCard_Niu {
		if ret1.number > ret2.number {
			return ret1.getBetTimes()
		} else if ret1.number < ret2.number {
			return -1 * ret2.getBetTimes()
		}
	}

	//相同牌型,比较单牌牌点
	sc1 := newNiuNiuCardSorter(ret1.cards)
	sort.Sort(sc1)
	sc2 := newNiuNiuCardSorter(ret2.cards)
	sort.Sort(sc2)

	//fmt.Printf("sc1 %v\r\n", ret1.cards)
	//fmt.Printf("sc2 %v\r\n", ret2.cards)

	niuniu1, _ := isNiuNiu(sc1.cardGroup)
	niuniu2, _ := isNiuNiu(sc2.cardGroup)
	if niuniu1 && niuniu2 {
		for i := 0; i < len(sc1.cardGroup.Cards); i++ {
			_, special1 := getNiuNiuCardScore(sc1.cardGroup.Cards[i])
			_, special2 := getNiuNiuCardScore(sc2.cardGroup.Cards[i])

			if special1 && special2 {
				pk1 := getNiuNiuCardNumber(sc1.cardGroup.Cards[i])
				pk2 := getNiuNiuCardNumber(sc2.cardGroup.Cards[i])
				if pk1 < pk2 {
					return -1 * ret1.getBetTimes()
				} else {
					return ret1.getBetTimes()
				}
			}

			if special1 && !special2 {
				return -1 * ret1.getBetTimes()
			} else if !special1 && special2 {
				return ret1.getBetTimes()
			}
		}
	}

	for i, c1 := range sc1.cardGroup.Cards {
		c1n, c2n := int(0), int(0)

		c1n = getNiuNiuCardNumber(c1)
		c2 := sc2.cardGroup.Cards[i]
		c2n = getNiuNiuCardNumber(c2)

		if c1n > c2n {
			return ret1.getBetTimes()
		} else if c1n < c2n {
			return -1 * ret1.getBetTimes()
		} else {
			//皮子
			if c1.IsPiZi() {
				return ret1.getBetTimes()
			}
			if c2.IsPiZi() {
				return -1 * ret1.getBetTimes()
			}

			//大王
			if c1.IsBigKing() {
				return ret1.getBetTimes()
			}
			if c2.IsBigKing() {
				return -1 * ret1.getBetTimes()
			}

			//小王
			if c1.IsSmallKing() {
				return ret1.getBetTimes()
			}
			if c2.IsSmallKing() {
				return -1 * ret1.getBetTimes()
			}

		}
	}

	return 0
}

//获取牛牛结果
func getNiuNiuRet(cards *poker.CardGroup) *niuniuCardRet {
	if isWuXiaoNiu(cards) { //五小牛
		return newNiuNiuCardRet(nnCard_WuXiaoNiu, 0, cards)
	} else if ok, cs := is4Zha(cards); ok { //4炸
		nnret := newNiuNiuCardRet(nnCard_4Zha, 0, cards)
		nnret.cardStack = cs
		return nnret
	} else if isWuHuaNiu(cards) { //五花牛
		return newNiuNiuCardRet(nnCard_WuHuaNiu, 0, cards)
	} else if ok, cs := isNiuNiu(cards); ok { //牛牛
		nnret := newNiuNiuCardRet(nnCard_NiuNiu, 0, cards)
		nnret.cardStack = cs
		return nnret
	}

	nnumber, cs := getNiuNumber(cards)

	//无牛
	if nnumber == -1 {
		return newNiuNiuCardRet(nnCard_WuNiu, 0, cards)
	} else {
		nnret := newNiuNiuCardRet(nnCard_Niu, nnumber, cards)
		nnret.cardStack = cs

		return nnret
	}

}

//取得牛牛的賠率
func getNiuNiuOdds(niuRet *niuniuCardRet) int32 {
	odds := int32(0)
	if niuRet.niuRet >= nnCard_WuXiaoNiu && niuRet.niuRet <= nnCard_NiuNiu {
		odds = 11
	}
	if niuRet.niuRet == nnCard_WuNiu {
		odds = 2
	}
	if niuRet.niuRet == nnCard_Niu {
		odds = int32(niuRet.number) + 1
	}
	return odds
}

//是否五小牛
func isWuXiaoNiu(cards *poker.CardGroup) bool {
	//计算总点数
	totalnum := 0
	for _, c := range cards.Cards {
		totalnum += getNiuNiuCardPoint(c)
	}

	return totalnum <= 10
}

//是否四炸
func is4Zha(cards *poker.CardGroup) (bool, *poker.CardStack) {
	for i, c := range cards.Cards {
		samenumber := 0
		for j := 0; j < len(cards.Cards); j++ {
			if j == i {
				continue
			}

			if c.IsSameCardNumber(cards.Cards[j]) {
				samenumber++
			}
		}

		if samenumber >= 3 {
			cs := poker.NewCardStack()
			cg1 := poker.NewCardGroup(nil)
			cg2 := poker.NewCardGroup(nil)

			for _, c := range cards.Cards {
				if c.IsSameCardNumber(cards.Cards[i]) {
					cg2.AddCard(c)
				} else {
					cg1.AddCard(c)
				}
			}

			cs.AddCardGroup(cg1)
			cs.AddCardGroup(cg2)

			return true, cs
		}
	}

	return false, nil
}

//是否五花牛
func isWuHuaNiu(cards *poker.CardGroup) bool {
	for _, c := range cards.Cards {
		point, special := getNiuNiuCardScore(c)
		if point < 10 || c.CardNumber == 10 || special {
			return false
		}
	}

	return true
}

//是否牛牛
func isNiuNiu(cards *poker.CardGroup) (bool, *poker.CardStack) {
	for i := 0; i < len(cards.Cards); i++ {
		for j := i + 1; j < len(cards.Cards); j++ {
			var val int = 0

			for k := 0; k < len(cards.Cards); k++ {
				if k == i || k == j {
					continue
				}

				val += getNiuNiuCardPoint(cards.Cards[k])
			}

			if val%10 != 0 {
				continue
			}

			//湊出第一組牛
			val = 0
			val += getNiuNiuCardPoint(cards.Cards[i])
			val += getNiuNiuCardPoint(cards.Cards[j])

			//檢驗最後兩張是否為皮子或大小王
			if val%10 == 0 || cards.Cards[i].IsPiZi() || cards.Cards[i].IsBigKing() || cards.Cards[i].IsSmallKing() ||
				cards.Cards[j].IsPiZi() || cards.Cards[j].IsBigKing() || cards.Cards[j].IsSmallKing() {
				cs := poker.NewCardStack()
				cg1 := poker.NewCardGroup(nil)
				cg1.AddCard(cards.Cards[i])
				cg1.AddCard(cards.Cards[j])
				cs.AddCardGroup(cg1)

				cg2 := poker.NewCardGroup(nil)
				for k := 0; k < len(cards.Cards); k++ {
					if k == i || k == j {
						continue
					}

					cg2.AddCard(cards.Cards[k])
				}
				cs.AddCardGroup(cg2)

				return true, cs
			}
		}
	}

	return false, nil

}

//從三張牌計算出剩餘的兩張牌可以組合成牛牛的牌
func (brnn *BRNNLogic) takeNiuNiuByNumber(cards *poker.CardGroup) *poker.CardGroup {
	ok, _ := isNiuNiu(cards)

	if ok { //牛牛
		takeCardGroup := brnn.assignNewNiuNiuCardMoreThen(cards, nnCard_NiuNiu, 10)
		cardGroup := takeCardGroup.cards.Cards[3:5]

		return &poker.CardGroup{Cards: cardGroup}
	} else { //沒牛
		takeCardGroup := brnn.assignNewNiuNiuCardMoreThen(cards, nnCard_NiuNiu, 10)
		cardGroup := takeCardGroup.cards.Cards[3:5]

		return &poker.CardGroup{Cards: cardGroup}
	}
}

//前三張是牛牛的組合, 後兩張也要是牛牛
func (brnn *BRNNLogic) hasNiuNiu() *poker.CardGroup {
	takeCard1 := brnn.takeKingOrPiziCard()
	takeCard2 := brnn.takeKingOrPiziCard()
	takeCardGroup := new(poker.CardGroup)
	pokerCard1 := new(poker.Card)
	pokerCard1.CardNumber = takeCard1.CardNumber
	pokerCard1.Color = takeCard1.Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard1)
	pokerCard2 := new(poker.Card)
	pokerCard2.CardNumber = takeCard2.CardNumber
	pokerCard2.Color = takeCard2.Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard2)
	return takeCardGroup
}

//前三張無牛需由後面兩張組出牛牛
func nonNiuNiu(cards *poker.CardGroup, brnn *BRNNLogic) *poker.CardGroup {

	for k, v := range brnn.allKindPoker {
		mlog.Debug("%v 手牌: %v", k, v.Cards)
	}

	sc1 := newNiuNiuCardSorter(cards)
	sort.Sort(sc1)

	findCardNumber1 := 10 - (getNiuNiuCardPoint(sc1.cardGroup.Cards[0])+getNiuNiuCardPoint(sc1.cardGroup.Cards[1]))%10
	findCardNumber2 := 10 - getNiuNiuCardPoint(sc1.cardGroup.Cards[2])%10

	if findCardNumber1 == 1 {
		findCardNumber1 = 14
	}
	if findCardNumber2 == 1 {
		findCardNumber2 = 14
	}
	ok1, takeCard1 := brnn.poker.TakeCardNumber(0, int32(findCardNumber1))
	ok2, takeCard2 := brnn.poker.TakeCardNumber(0, int32(findCardNumber2))
	if !ok1 { //取不到牌就找小王到皮子的牌
		takeCard1 = brnn.takeKingOrPiziCard()
	}
	if !ok2 {
		takeCard2 = brnn.takeKingOrPiziCard()
	}

	//if takeCard1.CardNumber == int32(0) {
	//	takeCard1 = brnn.poker.TakeCardNumber(15)
	//}
	takeCardGroup := new(poker.CardGroup)
	pokerCard1 := new(poker.Card)
	pokerCard1.CardNumber = takeCard1.CardNumber
	pokerCard1.Color = takeCard1.Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard1)
	pokerCard2 := new(poker.Card)
	pokerCard2.CardNumber = takeCard2.CardNumber
	pokerCard2.Color = takeCard2.Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard2)
	return takeCardGroup
}

//指定取小王到皮子其中之一的牌
func (brnn *BRNNLogic) takeKingOrPiziCard() poker.Card {
	var takeCard1 poker.Card
	for p := 15; p <= 17; p++ {
		okKing, takeKingCard := brnn.poker.TakeCardNumber(0, int32(p))
		if okKing {
			takeCard1 = takeKingCard
			break
		}
	}
	return takeCard1
}

//是否有牛 -1表示无牛
func getNiuNumber(cards *poker.CardGroup) (int, *poker.CardStack) {

	for i := 0; i < len(cards.Cards); i++ {
		for j := i + 1; j < len(cards.Cards); j++ {
			var val int = 0

			for k := 0; k < len(cards.Cards); k++ {
				if k == i || k == j {
					continue
				}

				val += getNiuNiuCardPoint(cards.Cards[k])
			}

			if val%10 == 0 {
				cs := poker.NewCardStack()
				cg1 := poker.NewCardGroup(nil)
				cg1.AddCard(cards.Cards[i])
				cg1.AddCard(cards.Cards[j])
				cs.AddCardGroup(cg1)

				cg2 := poker.NewCardGroup(nil)
				for k := 0; k < len(cards.Cards); k++ {
					if k == i || k == j {
						continue
					}

					cg2.AddCard(cards.Cards[k])
				}
				cs.AddCardGroup(cg2)

				return int((getNiuNiuCardPoint(cards.Cards[i]) + getNiuNiuCardPoint(cards.Cards[j])) % 10), cs
			}
		}
	}

	return -1, nil
}

//判斷四張牌有沒有牛
func getNiuNumber4Poker(cards *poker.CardGroup) (int, *poker.CardStack) {
	for i := 0; i < len(cards.Cards); i++ {
		for j := i + 1; j < len(cards.Cards); j++ {
			var val int = 0

			cs := poker.NewCardStack()
			cg2 := poker.NewCardGroup(nil)
			for k := 0; k < len(cards.Cards); k++ {
				if i == j || k == j {
					continue
				}

				val += getNiuNiuCardPoint(cards.Cards[k])
				cg2.AddCard(cards.Cards[k])
			}
			cs.AddCardGroup(cg2)

			if val%10 == 0 {
				return int((getNiuNiuCardPoint(cards.Cards[i]) + getNiuNiuCardPoint(cards.Cards[j])) % 10), cs
			}
		}
	}
	return -1, nil
}

//获取牛牛牌点
func getNiuNiuCardPoint(pk *poker.Card) int {
	tenPoint := 10
	onePoint := 1

	if pk.IsPiZi() {
		return tenPoint
	} else if pk.IsBigKing() {
		return tenPoint
	} else if pk.IsSmallKing() {
		return tenPoint
	} else if pk.IsA() {
		return onePoint
	} else if pk.CardNumber >= 10 {
		return tenPoint
	} else {
		return int(pk.CardNumber)
	}
}

//获取牛牛牌点
func getNiuNiuCardScore(pk *poker.Card) (int, bool) {
	tenPoint := 10
	onePoint := 1

	if pk.IsPiZi() {
		return tenPoint, true
	} else if pk.IsBigKing() {
		return tenPoint, true
	} else if pk.IsSmallKing() {
		return tenPoint, true
	} else if pk.IsA() {
		return onePoint, false
	} else if pk.CardNumber >= 10 {
		return tenPoint, false
	} else {
		return int(pk.CardNumber), false
	}
}

//牛牛牌型排序
type niuniuCardSorter struct {
	cardGroup *poker.CardGroup
}

//创建牛牛牌型排序对象
func newNiuNiuCardSorter(cards *poker.CardGroup) *niuniuCardSorter {
	nncs := new(niuniuCardSorter)
	nncs.cardGroup = cards

	return nncs
}

func (nns *niuniuCardSorter) Len() int {
	return nns.cardGroup.Len()
}

func (nns *niuniuCardSorter) Swap(i, j int) {
	nns.cardGroup.Swap(i, j)
}

func (nns *niuniuCardSorter) Less(i, j int) bool {
	cardi := nns.cardGroup.Cards[i]
	cardj := nns.cardGroup.Cards[j]

	inum := getNiuNiuCardNumber(cardi)
	jnum := getNiuNiuCardNumber(cardj)

	//皮子
	if cardi.IsPiZi() {
		return true
	}
	if cardj.IsPiZi() {
		return false
	}

	//大王
	if cardi.IsBigKing() {
		return true
	}
	if cardj.IsBigKing() {
		return false
	}

	//小王
	if cardi.IsSmallKing() {
		return true
	}
	if cardj.IsSmallKing() {
		return false
	}

	return inum > jnum
}

//获取牛牛牌面大小
func getNiuNiuCardNumber(pk *poker.Card) int {
	if pk.IsPiZi() {
		return nnCardNumber_PiZi
	} else if pk.IsBigKing() {
		return nnCardNumber_BingKing
	} else if pk.IsSmallKing() {
		return nnCardNumber_SmallKing
	} else {
		num := pk.CardNumber
		if pk.IsA() {
			num = 1
		}

		str := fmt.Sprintf("%2d%d", num, pk.Color)

		str = strings.TrimSpace(str)
		cnum, _ := strconv.Atoi(str)

		return cnum
	}
}

//指定產生牛牛類型和點數之上的牌型
func (brnn *BRNNLogic) assignNewNiuNiuCardMoreThen(firstThreeCard *poker.CardGroup, nnRet int, number int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(9) {
		nnRet = nnRet - 1
		number = 0
	}
	//nnRet 牌型, 值愈小的愈大. number 值愈大就愈大
	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = brnn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet == nnRet && result.number > number {
			example.invalid("莊牌 : niuRet[%v] number[%v] card[%v]", result.niuRet, result.number, result.cards)
			break
		}

		if result != nil && result.niuRet < nnRet {
			example.invalid("莊牌 : niuRet[%v] number[%v] card[%v]", result.niuRet, result.number, result.cards)
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[3:5]
			brnn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	return result
}

//隨機產生指定牌的類型
func (brnn *BRNNLogic) naturlAssignNiuCard(firstThreeCard *poker.CardGroup) *niuniuCardRet {
	tmpCardGroup := copyPokerCards(firstThreeCard)
	brnn.poker.Shuffle()
	cg := brnn.poker.Deal(0, 2)

	if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil {
		mlog.Warn("沒取到兩張牌")
		return nil
	}

	tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
	result := getNiuNiuRet(tmpCardGroup)
	return result
}

//指定產生牛牛類型和點數之下的牌型
func (brnn *BRNNLogic) assignNewNiuNiuCardLessThen(firstThreeCard *poker.CardGroup, nnRet int, number int, runCount int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(1) {
		nnRet = nnCard_Niu
		number = 9
	}

	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = brnn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet == nnRet && result.number < number {
			break
		}

		if result != nil && result.niuRet > nnRet {
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[3:5]
			brnn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	if result == nil || len(result.cards.Cards) != int(5) {
		if runCount <= int(0) {
			return result
		}
		runCount--
		brnn.assignNewNiuNiuCardLessThen(firstThreeCard, nnCard_WuNiu, 0, runCount)
	}

	return result
}

////指定產生有牛點數牌型
func (brnn *BRNNLogic) assignNewNiuNiuCardNumber(firstThreeCard *poker.CardGroup, number int) (result *niuniuCardRet) {
	result = new(niuniuCardRet)
	// 2-14 隨機一個一個取
	pCardNums := []int{2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}
	iCardNums := util.IntArrayShuffle(pCardNums)
	for _, i := range iCardNums { //取第四張
		if ok, card4 := brnn.poker.TakeCardNumber(0, int32(i)); ok {
			jCardNums := util.IntArrayShuffle(pCardNums)
			for _, i := range jCardNums {
				if ok, card5 := brnn.poker.TakeCardNumber(0, int32(i)); ok {
					tmpCardGroup := copyPokerCards(firstThreeCard)
					tmpCardGroup.Cards = append(tmpCardGroup.Cards, &card4, &card5)
					result = getNiuNiuRet(tmpCardGroup)
					if result.niuRet == nnCard_Niu && result.number == number {
						return result
					}
					brnn.poker.AddCardGroup(0, poker.NewCardGroup([]*poker.Card{&card5}))
				}
			}
			brnn.poker.AddCardGroup(0, poker.NewCardGroup([]*poker.Card{&card4}))
		}
	}

	cg := brnn.poker.Deal(0, 2)
	tmpCardGroup := copyPokerCards(firstThreeCard)
	tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
	result = getNiuNiuRet(tmpCardGroup)
	return result
}

////指定無牛的牌型
//func (brnn *BRNNLogic) assignWuNiuCard(firstThreeCard *poker.CardGroup) *niuniuCardRet {
//	result := new(niuniuCardRet)
//	//rollbackCarpGroup := new(poker.CardGroup)
//
//	if result.niuRet != nnCard_WuNiu {
//		mlog.Error("指定類型結果錯誤 : ", result.cards)
//	}
//
//	return result
//}
//
////指定有牛的點數牌型, 點數之下的牌型
//func (brnn *BRNNLogic) assignHasNiuCardLessThen(firstThreeCard *poker.CardGroup, number int) *niuniuCardRet {
//	result := new(niuniuCardRet)
//	rollbackCarpGroup := new(poker.CardGroup)
//	for true {
//
//		tmpCardGroup := copyPokerCards(firstThreeCard)
//		var cg *poker.CardGroup
//		cg = brnn.poker.Deal(0, 1)
//
//		if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil { //找完牌了, 就隨機取牌組合
//			brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//			rollbackCarpGroup = new(poker.CardGroup)
//			cg = brnn.poker.Deal(0, 2)
//			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//			result = getNiuNiuRet(tmpCardGroup)
//			return result
//		}
//
//		tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//
//		nnumber, _ := getNiuNumber4Poker(tmpCardGroup)
//		if nnumber == int(-1) {
//			//rollback card
//			rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
//		} else {
//			cg = brnn.getCombineNiuCardLessThen(tmpCardGroup, number)
//
//			if cg == nil || cg.Cards == nil {
//				brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//				cg = brnn.poker.Deal(0, 2)
//				tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//				result = getNiuNiuRet(tmpCardGroup)
//				return result
//			}
//
//			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//			result = getNiuNiuRet(tmpCardGroup)
//			break
//		}
//
//		if len(tmpCardGroup.Cards) >= int(5) {
//			mlog.Error("超過四張牌")
//		}
//	}
//
//	if len(rollbackCarpGroup.Cards) > int(0) {
//		brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//	}
//
//	return result
//}
//
////指定有牛的點數牌型, 點數之大的牌型
//func (brnn *BRNNLogic) assignHasNiuCardMoreThen(firstThreeCard *poker.CardGroup, number int) *niuniuCardRet {
//	result := new(niuniuCardRet)
//	rollbackCarpGroup := new(poker.CardGroup)
//	for true {
//
//		tmpCardGroup := copyPokerCards(firstThreeCard)
//		var cg *poker.CardGroup
//		cg = brnn.poker.Deal(0, 1)
//
//		if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil { //找完牌了, 就隨機取牌組合
//			brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//			rollbackCarpGroup = new(poker.CardGroup)
//			cg = brnn.poker.Deal(0, 2)
//			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//			result = getNiuNiuRet(tmpCardGroup)
//			return result
//		}
//
//		tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//		nnumber, _ := getNiuNumber4Poker(tmpCardGroup)
//		if nnumber == int(-1) {
//			//rollback card
//			rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
//		} else {
//			cg = brnn.getCombineNiuCardMoreThen(tmpCardGroup, number)
//
//			if cg == nil || cg.Cards == nil {
//				brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//				cg = brnn.poker.Deal(0, 2)
//				tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//				result = getNiuNiuRet(tmpCardGroup)
//				return result
//			}
//
//			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
//			result = getNiuNiuRet(tmpCardGroup)
//			break
//		}
//	}
//
//	if len(rollbackCarpGroup.Cards) > int(0) {
//		brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//		rollbackCarpGroup = new(poker.CardGroup)
//	}
//
//	return result
//}
//
////取得一張牌配前四張的組合後大於指定的點數
//func (brnn *BRNNLogic) getCombineNiuCardMoreThen(tmpCardGroup *poker.CardGroup, number int) *poker.CardGroup {
//	var cg *poker.CardGroup
//	for true {
//		copyCardGroup := copyPokerCards(tmpCardGroup)
//		cg = brnn.poker.Deal(0, 1)
//
//		if cg == nil || cg.Cards == nil || copyCardGroup.Cards == nil {
//			return nil
//		}
//
//		copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
//		ret := getNiuNiuRet(copyCardGroup)
//		if ret.number > number && ret.niuRet == nnCard_Niu {
//			break
//		}
//	}
//	return cg
//}

////取得一張牌配前四張的組合後小於指定的點數
//func (brnn *BRNNLogic) getCombineNiuCardLessThen(tmpCardGroup *poker.CardGroup, number int) *poker.CardGroup {
//	var cg *poker.CardGroup
//	rollbackCarpGroup := new(poker.CardGroup)
//	for true {
//		copyCardGroup := copyPokerCards(tmpCardGroup)
//		cg = brnn.poker.Deal(0, 1)
//
//		if cg == nil || cg.Cards == nil { //找完牌了, 就隨機取牌組合
//			brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//			cg = brnn.poker.Deal(0, 1)
//
//			if cg == nil || cg.Cards == nil || copyCardGroup.Cards == nil {
//				return nil
//			}
//
//			copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
//			//result = getNiuNiuRet(tmpCardGroup)
//			return cg
//		}
//
//		copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
//		ret := getNiuNiuRet(copyCardGroup)
//		if ret.number < number && ret.niuRet == nnCard_Niu {
//			break
//		}
//		rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
//	}
//
//	if len(rollbackCarpGroup.Cards) > int(0) {
//		brnn.poker.AddCardGroup(0, rollbackCarpGroup)
//		rollbackCarpGroup = new(poker.CardGroup)
//	}
//
//	return cg
//}

//拷贝牌型
func copyPokerCards(Cards *poker.CardGroup) *poker.CardGroup {
	tmpCards := poker.NewEmptyCardGroup(len(Cards.Cards))
	for i := 0; i < len(Cards.Cards); i++ {
		tmpCards.Cards[i].Color = Cards.Cards[i].Color
		tmpCards.Cards[i].CardNumber = Cards.Cards[i].CardNumber
	}

	return tmpCards
}
