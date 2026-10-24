package hhnn

import (
	"fmt"
	"math/rand"
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

	if ret1.niuRet == nnCard_4Zha && ret2.niuRet == nnCard_4Zha {
		card1 := ret1.cardStack.Groups[1].Cards
		card2 := ret2.cardStack.Groups[1].Cards
		if card1[0].CardNumber > card2[0].CardNumber {
			if card1[0].CardNumber == 14 { //1最小
				return -1 * ret1.getBetTimes()
			}
			return ret1.getBetTimes()
		} else {
			if card2[0].CardNumber == 14 { //1最小
				return 1 * ret1.getBetTimes()
			}
			return -1 * ret1.getBetTimes()
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

//[致命一擊]特殊牌的賠率
func getSpecialNiuNiuOdds(niuRet *niuniuCardRet, hhnn *HHNNLogic) int32 {
	odds := int32(0)

	if niuRet.niuRet == nnCard_NiuNiu {
		for _, v := range niuRet.cards.Cards {
			if v.CardNumber == int32(13) {
				odds = hhnn.paytable[0].betvalue[nnCard_NiuNiu].criticalodd
				break
			}
		}
	}
	if niuRet.niuRet == nnCard_WuXiaoNiu {
		odds = hhnn.paytable[0].betvalue[nnCard_WuXiaoNiu].criticalodd
	}
	if niuRet.niuRet == nnCard_4Zha {
		odds = hhnn.paytable[0].betvalue[nnCard_4Zha].criticalodd
	}
	if niuRet.niuRet == nnCard_WuHuaNiu {
		odds = hhnn.paytable[0].betvalue[nnCard_WuHuaNiu].criticalodd
	}

	return odds
}

//取得牛牛的賠率
func getNiuNiuOdds(niuRet *niuniuCardRet, hhnn *HHNNLogic) int32 {
	odds := int32(0)
	if niuRet.niuRet >= nnCard_WuXiaoNiu && niuRet.niuRet <= nnCard_NiuNiu {
		odds = hhnn.paytable[0].betvalue[nnCard_NiuNiu].odd
	}
	if niuRet.niuRet == nnCard_WuNiu {
		odds = hhnn.paytable[0].betvalue[14].odd
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
func (hhnn *HHNNLogic) takeNiuNiuByNumber(cards *poker.CardGroup) *poker.CardGroup {
	ok, _ := isNiuNiu(cards)

	if ok { //牛牛
		takeCardGroup := hhnn.hasNiuNiu()
		return takeCardGroup
	} else { //沒牛
		takeCardGroup := nonNiuNiu(cards, hhnn)
		return takeCardGroup
	}
}

//前三張是牛牛的組合, 後兩張也要是牛牛
func (hhnn *HHNNLogic) hasNiuNiu() *poker.CardGroup {
	takeCard1 := hhnn.takeKingOrPiziCard()
	takeCard2 := hhnn.takeKingOrPiziCard()
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
func nonNiuNiu(cards *poker.CardGroup, HHNN *HHNNLogic) *poker.CardGroup {

	//for k, v := range HHNN.firstPoker {
	//	mlog.Debug("%v 手牌: %v", k, v.Cards)
	//}

	sc1 := newNiuNiuCardSorter(cards)
	sort.Sort(sc1)

	findCardNumber1 := 10 - (getNiuNiuCardPoint(sc1.cardGroup.Cards[0])+getNiuNiuCardPoint(sc1.cardGroup.Cards[1]))%10
	findCardNumber2 := 10 - getNiuNiuCardPoint(sc1.cardGroup.Cards[2])

	if findCardNumber1 == 1 {
		findCardNumber1 = 14
	}
	if findCardNumber2 == 1 {
		findCardNumber2 = 14
	}
	ok1, takeCard1 := HHNN.poker.TakeCardNumber(0, int32(findCardNumber1))
	ok2, takeCard2 := HHNN.poker.TakeCardNumber(0, int32(findCardNumber2))
	if !ok1 { //取不到牌就找小王到皮子的牌
		takeCard1 = HHNN.takeKingOrPiziCard()
	}
	if !ok2 {
		takeCard2 = HHNN.takeKingOrPiziCard()
	}

	//if takeCard1.CardNumber == int32(0) {
	//	takeCard1 = HHNN.poker.TakeCardNumber(15)
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
func (hhnn *HHNNLogic) takeKingOrPiziCard() poker.Card {
	var takeCard1 poker.Card
	for p := 15; p <= 17; p++ {
		okKing, takeKingCard := hhnn.poker.TakeCardNumber(0, int32(p))
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
func (hhnn *HHNNLogic) assignNewNiuNiuCardMoreThen(firstThreeCard *poker.CardGroup, nnRet int, number int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(9) {
		nnRet = nnRet - 1
		number = 0
	}

	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = hhnn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet <= nnRet && result.number >= number {
			mlog.Debug("莊牌 : niuRet[%v] number[%v] card[%v]", result.niuRet, result.number, result.cards)
			break
		}

		if result != nil && result.niuRet <= nnRet && result.niuRet != nnCard_Niu {
			mlog.Debug("莊牌 : niuRet[%v] number[%v] card[%v]", result.niuRet, result.number, result.cards)
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[3:5]
			hhnn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	return result
}

//隨機產生指定牌的類型
func (hhnn *HHNNLogic) naturlAssignNiuCard(firstThreeCard *poker.CardGroup) *niuniuCardRet {
	tmpCardGroup := copyPokerCards(firstThreeCard)
	hhnn.poker.Shuffle()
	cg := hhnn.poker.Deal(0, 2)

	if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil {
		mlog.Warn("沒取到兩張牌")
		return nil
	}

	tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
	result := getNiuNiuRet(tmpCardGroup)
	return result
}

//指定產生牛牛類型和點數之下的牌型
func (hhnn *HHNNLogic) assignNewNiuNiuCardLessThen(firstThreeCard *poker.CardGroup, nnRet int, number int, runCount int) *niuniuCardRet {
	result := new(niuniuCardRet)

	if number == int(0) {
		nnRet = nnCard_Niu
		number = 9
	}

	limit := 1000
	for runLimit := 0; runLimit < limit; runLimit++ {
		result = hhnn.naturlAssignNiuCard(firstThreeCard)

		if result != nil && result.niuRet == nnRet && result.number < number {
			break
		}

		if result != nil && result.niuRet >= nnRet && result.niuRet != nnCard_Niu {
			break
		}

		if result != nil && runLimit < limit-1 {
			result.cards.Cards = result.cards.Cards[3:5]
			hhnn.poker.AddCardGroup(0, result.cards) //把牌放回去
		}
	}

	if result == nil || len(result.cards.Cards) != int(5) {
		if runCount <= int(0) {
			return result
		}
		runCount--
		hhnn.assignNewNiuNiuCardLessThen(firstThreeCard, nnCard_WuNiu, 0, runCount)
	}

	return result
}

//指定無牛的牌型
func (hhnn *HHNNLogic) assignWuNiuCard(firstThreeCard *poker.CardGroup) *niuniuCardRet {
	result := new(niuniuCardRet)
	//rollbackCarpGroup := new(poker.CardGroup)

	if result.niuRet != nnCard_WuNiu {
		mlog.Error("指定類型結果錯誤 : %v", result.cards)
	}

	return result
}

//指定有牛的點數牌型, 點數之下的牌型
func (hhnn *HHNNLogic) assignHasNiuCardLessThen(firstThreeCard *poker.CardGroup, number int) *niuniuCardRet {
	result := new(niuniuCardRet)
	rollbackCarpGroup := new(poker.CardGroup)
	for true {

		tmpCardGroup := copyPokerCards(firstThreeCard)
		var cg *poker.CardGroup
		cg = hhnn.poker.Deal(0, 1)

		if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil { //找完牌了, 就隨機取牌組合
			hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
			rollbackCarpGroup = new(poker.CardGroup)
			cg = hhnn.poker.Deal(0, 2)
			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
			result = getNiuNiuRet(tmpCardGroup)
			return result
		}

		tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)

		nnumber, _ := getNiuNumber4Poker(tmpCardGroup)
		if nnumber == int(-1) {
			//rollback card
			rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
		} else {
			cg = hhnn.getCombineNiuCardLessThen(tmpCardGroup, number)

			if cg == nil || cg.Cards == nil {
				hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
				cg = hhnn.poker.Deal(0, 2)
				tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
				result = getNiuNiuRet(tmpCardGroup)
				return result
			}

			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
			result = getNiuNiuRet(tmpCardGroup)
			break
		}

		if len(tmpCardGroup.Cards) >= int(5) {
			mlog.Error("超過四張牌")
		}
	}

	if len(rollbackCarpGroup.Cards) > int(0) {
		hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
	}

	return result
}

//指定有牛的點數牌型, 點數之大的牌型
func (hhnn *HHNNLogic) assignHasNiuCardMoreThen(firstThreeCard *poker.CardGroup, number int) *niuniuCardRet {
	result := new(niuniuCardRet)
	rollbackCarpGroup := new(poker.CardGroup)
	for true {

		tmpCardGroup := copyPokerCards(firstThreeCard)
		var cg *poker.CardGroup
		cg = hhnn.poker.Deal(0, 1)

		if cg == nil || cg.Cards == nil || tmpCardGroup.Cards == nil { //找完牌了, 就隨機取牌組合
			hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
			rollbackCarpGroup = new(poker.CardGroup)
			cg = hhnn.poker.Deal(0, 2)
			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
			result = getNiuNiuRet(tmpCardGroup)
			return result
		}

		tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
		nnumber, _ := getNiuNumber4Poker(tmpCardGroup)
		if nnumber == int(-1) {
			//rollback card
			rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
		} else {
			cg = hhnn.getCombineNiuCardMoreThen(tmpCardGroup, number)

			if cg == nil || cg.Cards == nil {
				hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
				cg = hhnn.poker.Deal(0, 2)
				tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
				result = getNiuNiuRet(tmpCardGroup)
				return result
			}

			tmpCardGroup.Cards = append(tmpCardGroup.Cards, cg.Cards...)
			result = getNiuNiuRet(tmpCardGroup)
			break
		}
	}

	if len(rollbackCarpGroup.Cards) > int(0) {
		hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
		rollbackCarpGroup = new(poker.CardGroup)
	}

	return result
}

//取得一張牌配前四張的組合後大於指定的點數
func (hhnn *HHNNLogic) getCombineNiuCardMoreThen(tmpCardGroup *poker.CardGroup, number int) *poker.CardGroup {
	var cg *poker.CardGroup
	for true {
		copyCardGroup := copyPokerCards(tmpCardGroup)
		cg = hhnn.poker.Deal(0, 1)

		if cg == nil || cg.Cards == nil || copyCardGroup.Cards == nil {
			return nil
		}

		copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
		ret := getNiuNiuRet(copyCardGroup)
		if ret.number > number && ret.niuRet == nnCard_Niu {
			break
		}
	}
	return cg
}

//取得一張牌配前四張的組合後小於指定的點數
func (hhnn *HHNNLogic) getCombineNiuCardLessThen(tmpCardGroup *poker.CardGroup, number int) *poker.CardGroup {
	var cg *poker.CardGroup
	rollbackCarpGroup := new(poker.CardGroup)
	for true {
		copyCardGroup := copyPokerCards(tmpCardGroup)
		cg = hhnn.poker.Deal(0, 1)

		if cg == nil || cg.Cards == nil { //找完牌了, 就隨機取牌組合
			hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
			cg = hhnn.poker.Deal(0, 1)

			if cg == nil || cg.Cards == nil || copyCardGroup.Cards == nil {
				return nil
			}

			copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
			//result = getNiuNiuRet(tmpCardGroup)
			return cg
		}

		copyCardGroup.Cards = append(copyCardGroup.Cards, cg.Cards...)
		ret := getNiuNiuRet(copyCardGroup)
		if ret.number < number && ret.niuRet == nnCard_Niu {
			break
		}
		rollbackCarpGroup.Cards = append(rollbackCarpGroup.Cards, cg.Cards...)
	}

	if len(rollbackCarpGroup.Cards) > int(0) {
		hhnn.poker.AddCardGroup(0, rollbackCarpGroup)
		rollbackCarpGroup = new(poker.CardGroup)
	}

	return cg
}

//拷贝牌型
func copyPokerCards(Cards *poker.CardGroup) *poker.CardGroup {
	tmpCards := poker.NewEmptyCardGroup(len(Cards.Cards))
	for i := 0; i < len(Cards.Cards); i++ {
		tmpCards.Cards[i].Color = Cards.Cards[i].Color
		tmpCards.Cards[i].CardNumber = Cards.Cards[i].CardNumber
	}

	return tmpCards
}

//獲得某種牌型
func getCardsbyType(cardtype int, HHNN *HHNNLogic) *poker.CardGroup {
	if cardtype == WUXIAONIU {
		return getWuXiaoNiu(HHNN)
	} else if cardtype == ZHADAN {
		return getZhaDan(HHNN)
	} else if cardtype == WUHUANIU {
		return getWuHuaNiu(HHNN)
	} else if cardtype == KNIUNIU {
		return getNiuNiu(HHNN, true)
	} else if cardtype == NIUNIU {
		return getNiuNiu(HHNN, false)
	} else if cardtype == WUNIU {
		return getWuNiu(HHNN)
	}
	return getHasNiu(HHNN, 18-cardtype) //牛? = cardtype + 18
}

func getZhaDan(HHNN *HHNNLogic) *poker.CardGroup {
	r := util.RndInt(2, 10)
	takeCardGroup := new(poker.CardGroup)
	for i := 0; i < 4; i++ {
		_, takeCard := HHNN.poker.TakeCardNumber(0, int32(r))
		pokerCard := new(poker.Card)
		pokerCard.CardNumber = takeCard.CardNumber
		pokerCard.Color = takeCard.Color
		takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard)
	}
	//第五張牌隨便發
	cg := HHNN.poker.Deal(0, 1)
	if r == 2 && cg.Cards[0].CardNumber == 14 {
		HHNN.poker.AddCardGroup(0, cg)
		cg = HHNN.poker.Deal(0, 1)
	}
	pokerCard5 := new(poker.Card)
	pokerCard5.CardNumber = cg.Cards[0].CardNumber
	pokerCard5.Color = cg.Cards[0].Color
	takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard5)
	return takeCardGroup
}

func getWuXiaoNiu(HHNN *HHNNLogic) *poker.CardGroup {
	wxn := wuXiaoNiuType()
	takeCardGroup := new(poker.CardGroup)
	for _, v := range wxn {
		if v == 1 {
			v = 14
		}
		_, takeCard := HHNN.poker.TakeCardNumber(0, int32(v))
		pokerCard := new(poker.Card)
		pokerCard.CardNumber = takeCard.CardNumber
		pokerCard.Color = takeCard.Color
		takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard)
	}
	return takeCardGroup
}

func wuXiaoNiuType() []int {
	r := rand.Intn(16)
	switch {
	case r == 0:
		return []int{1, 1, 1, 1, 2}
	case r == 1:
		return []int{1, 1, 1, 1, 3}
	case r == 2:
		return []int{1, 1, 1, 2, 2}
	case r == 3:
		return []int{1, 1, 1, 1, 4}
	case r == 4:
		return []int{1, 1, 2, 2, 2}
	case r == 5:
		return []int{1, 1, 1, 1, 5}
	case r == 6:
		return []int{1, 1, 1, 3, 3}
	case r == 7:
		return []int{1, 1, 1, 2, 4}
	case r == 8:
		return []int{1, 1, 2, 2, 3}
	case r == 9:
		return []int{1, 2, 2, 2, 3}
	case r == 10:
		return []int{1, 1, 1, 2, 5}
	case r == 11:
		return []int{1, 1, 1, 3, 4}
	case r == 12:
		return []int{1, 1, 2, 3, 3}
	case r == 13:
		return []int{1, 1, 2, 2, 4}
	case r == 14:
		return []int{1, 1, 1, 2, 3}
	case r == 15:
		return []int{1, 2, 2, 2, 2}
	}
	return nil
}

func getWuHuaNiu(HHNN *HHNNLogic) *poker.CardGroup {
	takeCardGroup := new(poker.CardGroup)
	for i := 0; i < 5; i++ {
		r := util.RndInt(11, 14)
		ok, takeCard := HHNN.poker.TakeCardNumber(0, int32(r))
		if !ok { //取不到牌就再取一次
			if r > 11 {
				_, takeCard = HHNN.poker.TakeCardNumber(0, int32(r-1))
			} else {
				_, takeCard = HHNN.poker.TakeCardNumber(0, int32(r+1))
			}
		}
		pokerCard := new(poker.Card)
		pokerCard.CardNumber = takeCard.CardNumber
		pokerCard.Color = takeCard.Color
		takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard)
	}
	return takeCardGroup
}

func getNiuNiu(HHNN *HHNNLogic, hasK bool) *poker.CardGroup {
	takeCardGroup := new(poker.CardGroup)
	if hasK { //第一張牌就取K
		_, takeCardK := HHNN.poker.TakeCardNumber(0, int32(13))
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCardK)
		cg := HHNN.poker.Deal(0, 2)
		for _, u := range cg.Cards {
			pokerCard2 := new(poker.Card)
			pokerCard2.CardNumber = u.CardNumber
			pokerCard2.Color = u.Color
			takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard2)
		}
	} else { //取2~Q中間的牌
		for i := 0; i < 3; i++ {
			p := util.RndInt(2, 13)
			_, takeCard1 := HHNN.poker.TakeCardNumber(0, int32(p))
			takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard1)
		}
	}
	sc1 := newNiuNiuCardSorter(takeCardGroup)
	sort.Sort(sc1)
	findCardNumber1 := 10 - (getNiuNiuCardPoint(sc1.cardGroup.Cards[0])+getNiuNiuCardPoint(sc1.cardGroup.Cards[1]))%10
	findCardNumber2 := 10 - getNiuNiuCardPoint(sc1.cardGroup.Cards[2])%10

	cardPoint := []int{findCardNumber1, findCardNumber2}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := HHNN.poker.TakeCardNumber(0, int32(point))
		if !ok {
			_, takeCard = HHNN.poker.TakeCardNumber(0, int32(11))
		}
		pokerCard := new(poker.Card)
		pokerCard.CardNumber = takeCard.CardNumber
		pokerCard.Color = takeCard.Color
		takeCardGroup.Cards = append(takeCardGroup.Cards, pokerCard)
	}
	return takeCardGroup
}

func getHasNiu(HHNN *HHNNLogic, number int) *poker.CardGroup {
	takeCardGroup := new(poker.CardGroup)
	for i := 0; i < 2; i++ {
		p := util.RndInt(2, 10)
		_, takeCard1 := HHNN.poker.TakeCardNumber(0, int32(p))
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard1)
	}
	sc1 := newNiuNiuCardSorter(takeCardGroup)
	sort.Sort(sc1)
	findCardNumber1 := 10 - (getNiuNiuCardPoint(sc1.cardGroup.Cards[0])+getNiuNiuCardPoint(sc1.cardGroup.Cards[1]))%10
	findCardNumber2 := util.RndInt(2, 10)
	if findCardNumber2 == findCardNumber1 || int32(findCardNumber2) == sc1.cardGroup.Cards[0].CardNumber {
		findCardNumber2 = 10
	}
	findCardNumber3 := (10 + number - findCardNumber2) % 10

	cardPoint := []int{findCardNumber3, findCardNumber2, findCardNumber1}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := HHNN.poker.TakeCardNumber(0, int32(point))
		if !ok {
			takeCard = getCardOverTen(HHNN)
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	return takeCardGroup
}

//扣除3張成10
func getWuNiu(HHNN *HHNNLogic) *poker.CardGroup {
	cardnumber := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	cardPoint := make([]int, 0)
	takeCardGroup := HHNN.poker.Deal(0, 2)

	removeNumber := 10 - (getNiuNiuCardPoint(takeCardGroup.Cards[0])+getNiuNiuCardPoint(takeCardGroup.Cards[1]))%10
	cardnumber = remove(removeNumber, cardnumber)
	index := rand.Intn(len(cardnumber))
	findCardNumber1 := cardnumber[index]
	cardPoint = append(cardPoint, findCardNumber1)
	removeNumberSlice := checkRemoveNumber(takeCardGroup, cardPoint)
	for _, number := range removeNumberSlice {
		cardnumber = remove(number, cardnumber)
	}

	index = rand.Intn(len(cardnumber))
	findCardNumber2 := cardnumber[index]
	cardPoint = append(cardPoint, findCardNumber2)
	removeNumberSlice = checkRemoveNumber(takeCardGroup, cardPoint)
	for _, u := range removeNumberSlice {
		cardnumber = remove(u, cardnumber)
	}

	index = rand.Intn(len(cardnumber))
	findCardNumber3 := cardnumber[index]

	cardPoint = append(cardPoint, findCardNumber3)
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := HHNN.poker.TakeCardNumber(0, int32(point))
		if !ok { //若取不到，隨機取10~K，再取不到，隨便發一張
			takeCard = getCardOverTen(HHNN)
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	return takeCardGroup
}

func remove(number int, numberSlice []int) []int {
	for n := 0; n < len(numberSlice); n++ {
		if numberSlice[n] == number {
			numberSlice = append(numberSlice[:n], numberSlice[n+1:]...)
		}
	}
	return numberSlice
}

func checkRemoveNumber(sc *poker.CardGroup, findCardNumber []int) []int {
	cardnumber1 := getNiuNiuCardPoint(sc.Cards[0])
	cardnumber2 := getNiuNiuCardPoint(sc.Cards[1])
	switch {
	case len(findCardNumber) == 1:
		num1 := 10 - (cardnumber1+findCardNumber[0])%10
		num2 := 10 - (cardnumber2+findCardNumber[0])%10
		return []int{num1, num2}
	case len(findCardNumber) == 2:
		num1 := 10 - (cardnumber1+findCardNumber[1])%10
		num2 := 10 - (cardnumber2+findCardNumber[1])%10
		num3 := 10 - (findCardNumber[0]+findCardNumber[1])%10
		return []int{num1, num2, num3}
	}
	return nil
}

func getCardOverTen(HHNN *HHNNLogic) poker.Card {
	var takeCard poker.Card
	for p := 11; p <= 13; p++ {
		okKing, takeKingCard := HHNN.poker.TakeCardNumber(0, int32(p))
		if okKing {
			takeCard = takeKingCard
			break
		}
	}
	return takeCard
}

func getLoseCard(HHNN *HHNNLogic, wincard *poker.CardGroup) (*niuniuCardRet, *niuniuCardRet) {
	winNiuRet := getNiuNiuRet(wincard)
	takeCardGroup := HHNN.poker.Deal(0, 3)
	loseCard := HHNN.assignNewNiuNiuCardLessThen(takeCardGroup, winNiuRet.niuRet, winNiuRet.number, 5)
	judge := compareCards(winNiuRet, loseCard)
	if judge < 0 {
		loseCard = getNiuNiuRet(getWuNiu(HHNN))
		judge = compareCards(winNiuRet, loseCard)
		if judge < 0 {
			winNiuRet = loseCard
			loseCard = getNiuNiuRet(wincard)
		}
	}
	return winNiuRet, loseCard
}

//指定產生類型和點數之上的牌型
func (hhnn *HHNNLogic) assignNewCardMoreThen(firstThreeCard *poker.CardGroup, nnRet int, number int, runCount int) *niuniuCardRet {
	result := new(niuniuCardRet)
	takeCardGroup := new(poker.CardGroup)
	takeCardGroup.Cards = append(takeCardGroup.Cards, firstThreeCard.Cards...)
	//前三張 ret 希望牛幾
	first := getNiuNiuCardPoint(firstThreeCard.Cards[0])
	second := getNiuNiuCardPoint(firstThreeCard.Cards[1])
	third := getNiuNiuCardPoint(firstThreeCard.Cards[2])
	var (
		findCardNumber1 int
		findCardNumber2 int
	)

	if (first+second+third)%10 == 0 {
		findCardNumber1 = util.RndInt(2, 10)
		findCardNumber2 = (10 + number - findCardNumber1) % 10
	} else {
		findCardNumber1 = 10 - (first+second)%10
		findCardNumber2 = (10 + number - third) % 10
	}

	cardPoint := []int{findCardNumber1, findCardNumber2}
	for _, point := range cardPoint {
		if point == 1 {
			point = 14
		}
		ok, takeCard := hhnn.poker.TakeCardNumber(0, int32(point))
		if !ok {
			if point < 14 {
				ok, takeCard = hhnn.poker.TakeCardNumber(0, int32(point+1))
			}
			if !ok {
				takeCard = getCardOverTen(hhnn)
			}
		}
		takeCardGroup.Cards = append(takeCardGroup.Cards, &takeCard)
	}
	result = getNiuNiuRet(takeCardGroup)
	if result.niuRet == nnCard_WuNiu && result.number < 2 {
		if runCount <= int(0) {
			return result
		}
		runCount--
		hhnn.assignNewCardMoreThen(firstThreeCard, nnRet, 0, runCount)
	}

	return result
}
