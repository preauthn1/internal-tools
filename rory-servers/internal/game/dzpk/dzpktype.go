//==================================================================
//文件名： dzpktype.go
//说明：   德州扑克牌型检测以及大小比较
//修改记录：
//==================================================================
package dzpk

import (
	"bytes"
	"encoding/gob"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// 牌型定义
const (
	HIGHCARD     = iota /*高牌*/
	ONEPAIR             /*一对*/
	TWOPAIR             /*两对*/
	THREEOFAKIND        /*三条*/
	STRIGHT             /*顺子*/
	FLUSH               /*同花*/
	FULLHOUSE           /*葫芦*/
	FOUROFAKIND         /*四条*/
	STRIGHTFLUSH        /*同花顺*/
	ROYALFLUSH          /*皇家同花顺*/
)

// 牌型装换为中文名称
func strCardType(Type int32) string {
	if Type == HIGHCARD {
		return "高牌"
	} else if Type == ONEPAIR {
		return "一对"
	} else if Type == TWOPAIR {
		return "两对"
	} else if Type == THREEOFAKIND {
		return "三条"
	} else if Type == STRIGHT {
		return "顺子"
	} else if Type == FLUSH {
		return "同花"
	} else if Type == FULLHOUSE {
		return "葫芦"
	} else if Type == FOUROFAKIND {
		return "四条"
	} else if Type == STRIGHTFLUSH {
		return "同花顺"
	} else if Type == ROYALFLUSH {
		return "皇家同花顺"
	}
	return "未知牌型"
}

//牌型判断 输入：2~7张牌

//判断获得牌型
func GetCardType(CardsSrc *poker.CardGroup) int32 {

	if CardsSrc == nil {
		mlog.Error("Card is nil.")
		return 0
	}

	Cards := copyCardGroup(CardsSrc)

	flag := false
	if isRoyalFlush(Cards) {
		return ROYALFLUSH
	}
	flag, _, _ = isStrightFlush(Cards)
	if flag {
		return STRIGHTFLUSH
	}
	flag, _ = isFourOfAKind(Cards)
	if flag {
		return FOUROFAKIND
	}
	flag, _, _ = isFullHouse(Cards)
	if flag {
		return FULLHOUSE
	}
	flag, _ = isFlush(Cards)
	if flag {
		return FLUSH
	}
	flag, _ = isStright(Cards)
	if flag {
		return STRIGHT
	}
	flag, _ = isThreeOfAKind(Cards)
	if flag {
		return THREEOFAKIND
	}
	flag, _, _ = isTwoPair(Cards)
	if flag {
		return TWOPAIR
	}
	flag, _ = isOnePair(Cards)
	if flag {
		return ONEPAIR
	}
	return HIGHCARD
}

//牌型大小比较,返回-1：小于；0：等于；1大于
func ComPareCards(CardsSrc1 *poker.CardGroup, CardsSrc2 *poker.CardGroup) int32 {
	if CardsSrc1 == nil || CardsSrc2 == nil {
		mlog.Error("ComPareCards,CardsSrc1 or CardsSrc2 is nil.")
		return 0
	}
	Cards1 := copyCardGroup(CardsSrc1)
	Cards2 := copyCardGroup(CardsSrc2)

	typeOne := GetCardType(Cards1)
	typeTwo := GetCardType(Cards2)
	if typeOne > typeTwo {
		return 1
	}
	if typeOne < typeTwo {
		return -1
	}
	if typeOne == typeTwo {
		return compareCardsForSameType(Cards1, Cards2, typeOne)
	}
	return 0
}

//皇家同花顺
func isRoyalFlush(Cards *poker.CardGroup) bool {
	//mlog.Debug("isRoyalFlush begin.")
	flag, max, _ := isStrightFlush(Cards)
	if flag && max == 14 {
		return true
	}

	return false
}

//同花顺,返回值：flag 是否是同花顺，max最大点数
func isStrightFlush(Cards *poker.CardGroup) (flag bool, max int32, color int32) {
	//mlog.Debug("isStrightFlush begin.")
	if Cards.Len() < 5 {
		return false, 0, 0
	}
	flagColor, color := isFlush(Cards)
	if !flagColor {
		return false, 0, 0
	}
	var cardSameColor = new(poker.CardGroup)
	for _, cardTmp := range Cards.Cards {
		if cardTmp.Color == color-1 {
			cardSameColor.Cards = append(cardSameColor.Cards, cardTmp)
		}
	}
	f, max := isStright(cardSameColor)
	return f, max, color
}

//四条
func isFourOfAKind(Cards *poker.CardGroup) (flag bool, max int32) {
	//mlog.Debug("isFourOfAKind begin.")
	if Cards.Len() < 4 {
		return false, 0
	}
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	for index, card := range cardTmp {
		if card == 4 {
			return true, index
		}
	}
	return false, 0
}

//葫芦
func isFullHouse(Cards *poker.CardGroup) (flag bool, max int32, two int32) {
	//mlog.Debug("isFullHouse begin.")
	if Cards.Len() < 5 {
		return false, 0, 0
	}
	//挑出五张检测最大的
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	var nThreeOne, nThreeTwo int32 = 0, 0
	for i, card := range cardTmp {
		if card == 3 {
			nThreeOne = i
			break
		}
	}
	if nThreeOne == 0 {
		return false, 0, 0
	}
	for i, card := range cardTmp {
		if i == nThreeOne {
			continue
		}
		if card == 3 {
			nThreeTwo = i
			if nThreeOne > nThreeTwo {
				return true, nThreeOne, nThreeTwo
			} else {
				return true, nThreeTwo, nThreeOne
			}

		}
		if card == 2 {
			return true, nThreeOne, i
		}
	}
	return false, 0, 0
}

//同花
func isFlush(Cards *poker.CardGroup) (flag bool, color int32) {
	//mlog.Debug("isFlush begin.")
	if Cards.Len() < 5 {
		return false, 0
	}
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.Color+1]++
	}
	for color, card := range cardTmp {
		if card >= 5 {
			return true, color
		}
	}
	return false, 0
}

//顺子
func isStright(Cards *poker.CardGroup) (flag bool, max int32) {
	//mlog.Debug("isStright begin.")
	if Cards.Len() < 5 {
		return false, 0
	}
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	tmp := 15
	for i := 14; i > 1; i-- {
		if cardTmp[int32(i)] == 0 {
			tmp = i
		}
		if tmp-i >= 5 {
			max := tmp - 1
			//min := max - 5
			return true, int32(max)
		}
	}
	//检测是否有12345
	if cardTmp[14] != 0 && cardTmp[2] != 0 && cardTmp[3] != 0 && cardTmp[4] != 0 && cardTmp[5] != 0 {
		return true, 5
	}
	return false, 0
}

//三条
func isThreeOfAKind(Cards *poker.CardGroup) (flag bool, max int32) {
	//mlog.Debug("isThreeOfAKind begin.")
	if Cards.Len() < 3 {
		return false, 0
	}
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	for index, card := range cardTmp {
		if card == 3 {
			return true, index
		}
	}
	return false, 0
}

//两对
func isTwoPair(Cards *poker.CardGroup) (flag bool, max int32, min int32) {
	//mlog.Debug("isTwoPair begin.")
	if Cards.Len() < 4 {
		return false, 0, 0
	}
	var cardTmp [15]int32
	//使用桶子排序演算法
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	var nIndex int32 = 0
	mapReturn := make(map[int32]int32)
	for index := range cardTmp {
		cardNo := 14 - index //從由大到小讀取
		if cardTmp[cardNo] == 2 {
			mapReturn[nIndex] = int32(cardNo)
			nIndex++
			if nIndex >= 2 {
				if mapReturn[0] > mapReturn[1] { //讀取最大兩個值
					return true, mapReturn[0], mapReturn[1]
				}
			}
		}
	}
	return false, 0, 0
}

//一对
func isOnePair(Cards *poker.CardGroup) (flag bool, max int32) {
	//mlog.Debug("isOnePair begin.")
	var cardTmp = make(map[int32]int32)
	for _, card := range Cards.Cards {
		cardTmp[card.CardNumber]++
	}
	for i, card := range cardTmp {
		if card == 2 {
			return true, i
		}
	}
	return false, 0
}

//相同牌型大小比较,返回-1：小于；0：等于；1大于
func compareCardsForSameType(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, typeOne int32) int32 {
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareCardsForSameType,Cards1 or Cards2 is nil.")
		return 0
	}
	switch typeOne {
	case ROYALFLUSH:
		return 0
	case STRIGHTFLUSH:
		_, maxA, _ := isStrightFlush(Cards1)
		_, maxB, _ := isStrightFlush(Cards2)
		if maxA > maxB {
			return 1
		}

		if maxA < maxB {
			return -1
		}
		return 0
	case FOUROFAKIND:
		_, maxA := isFourOfAKind(Cards1)
		_, maxB := isFourOfAKind(Cards2)
		if maxA > maxB {
			return 1
		}
		if maxA < maxB {
			return -1
		}
		return compareSameFourOfAKind(Cards1, Cards2, maxA)
	case FULLHOUSE:
		_, maxA, _ := isFullHouse(Cards1)
		_, maxB, _ := isFullHouse(Cards2)
		if maxA > maxB {
			return 1
		} else {
			return -1
		}
	case FLUSH:
		_, colorA := isFlush(Cards1)
		_, colorB := isFlush(Cards1)
		cardA := getSameColorCards(Cards1, colorA-1)
		cardB := getSameColorCards(Cards2, colorB-1)
		return compareHighCards(cardA, cardB, 5)

	case STRIGHT:
		_, maxA := isStright(Cards1)
		_, maxB := isStright(Cards2)
		if maxA > maxB {
			return 1
		}
		if maxA < maxB {
			return -1
		}
		return 0
	case THREEOFAKIND:
		_, maxA := isThreeOfAKind(Cards1)
		_, maxB := isThreeOfAKind(Cards2)
		if maxA > maxB {
			return 1
		}
		if maxA < maxB {
			return -1
		}
		//公共牌为3条，则还需要比较剩余牌
		return compareSameThreeOfAKind(Cards1, Cards2, maxA)
	case TWOPAIR:
		_, maxA, minA := isTwoPair(Cards1)
		_, maxB, minB := isTwoPair(Cards2)
		if maxA > maxB {
			return 1
		}
		if maxA < maxB {
			return -1
		}
		if minA > minB {
			return 1
		}
		if minA < minB {
			return -1
		}
		return compareSameTwoPairCard(Cards1, Cards2, maxA, minA)

	case ONEPAIR:
		_, maxA := isOnePair(Cards1)
		_, maxB := isOnePair(Cards2)
		if maxA > maxB {
			return 1
		}
		if maxA < maxB {
			return -1
		}
		return compareSameOnePairMax(Cards1, Cards2, maxA)
	case HIGHCARD:
		return compareHighCards(Cards1, Cards2, 5)
	}
	return 0
}

//根据花色提取牌
func getSameColorCards(CardsIn *poker.CardGroup, color int32) *poker.CardGroup {
	CardsOut := new(poker.CardGroup)
	for _, cards := range CardsIn.Cards {
		if cards.Color == color {
			CardsOut.Cards = append(CardsOut.Cards, cards)
		}
	}
	return CardsOut
}

//两队都相同比较
func compareSameTwoPairCard(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, max, min int32) int32 {
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareSameTwoPairCard,Cards1 or Cards2 is nil.")
		return 0
	}
	cardTmpA := getMaxCardsExcept(Cards1, max)
	cardA := getMaxCardsExcept(cardTmpA, min)
	cardTmpB := getMaxCardsExcept(Cards2, max)
	cardB := getMaxCardsExcept(cardTmpB, min)
	return compareHighCards(cardA, cardB, 1)
}

//比较相同三条大小
func compareSameThreeOfAKind(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, card int32) int32 {
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareSameThreeOfAKind,Cards1 or Cards2 is nil.")
		return 0
	}
	cardTmpA := getMaxCardsExcept(Cards1, card)
	cardTmpB := getMaxCardsExcept(Cards2, card)
	return compareHighCards(cardTmpA, cardTmpB, 2)
}

//比较相同4条大小
func compareSameFourOfAKind(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, card int32) int32 {
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareSameThreeOfAKind,Cards1 or Cards2 is nil.")
		return 0
	}
	cardTmpA := getMaxCardsExcept(Cards1, card)
	cardTmpB := getMaxCardsExcept(Cards2, card)
	return compareHighCards(cardTmpA, cardTmpB, 1)
}

//去处card中特定的牌
func getMaxCardsExcept(Cards1 *poker.CardGroup, card int32) *poker.CardGroup {
	CardsOut := new(poker.CardGroup)
	for _, cards := range Cards1.Cards {
		if cards.CardNumber == card {
			continue
		}
		CardsOut.Cards = append(CardsOut.Cards, cards)
	}
	return CardsOut
}

//对子相同比较
func compareSameOnePairMax(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, max int32) int32 {
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareSameOnePairMax,Cards1 or Cards2 is nil.")
		return 0
	}
	cardA := getMaxCardsExcept(Cards1, max)
	cardB := getMaxCardsExcept(Cards2, max)
	return compareHighCards(cardA, cardB, 3)
}

//比较高牌大小
func compareHighCards(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup, deep int32) int32 {
	if Cards1.Len() < int(deep) || Cards2.Len() < int(deep) {
		mlog.Error("compareHighCards,card size= %d,%d;deep=%d", Cards1.Len(), Cards2.Len(), deep)
		return 0
	}
	if Cards1 == nil || Cards2 == nil {
		mlog.Error("compareHighCards,Cards1 or Cards2 is nil.")
		return 0
	}
	CardSort(Cards1)
	CardSort(Cards2)
	for i := 0; i < int(deep); i++ {
		if Cards1.Cards[i].CardNumber > Cards2.Cards[i].CardNumber {
			return 1
		}
		if Cards1.Cards[i].CardNumber < Cards2.Cards[i].CardNumber {
			return -1
		}
	}
	return 0
}

//7张牌中取5张最大的牌
func GetCardsByCardType(cardGroupSrc *poker.CardGroup, nCardType int32) *poker.CardGroup {
	if cardGroupSrc == nil {
		mlog.Error("getCardsByCardType failed,cardGroup is nil")
		return nil
	}
	cardGroup := copyCardGroup(cardGroupSrc)

	switch nCardType {
	case ROYALFLUSH:
		return getRoyalflushCardsGroup(cardGroup)
	case STRIGHTFLUSH:
		return getStrightflushCardsGroup(cardGroup)
	case FOUROFAKIND:
		return getFourofakindCardsGroup(cardGroup)
	case FULLHOUSE:
		return getFullhouseCardsGroup(cardGroup)
	case FLUSH:
		return getFlushCardsGroup(cardGroup)
	case STRIGHT:
		return getStrightCardsGroup(cardGroup)
	case THREEOFAKIND:
		return getThreeofakindCardsGroup(cardGroup)
	case TWOPAIR:
		return getTwoPairCardsGroup(cardGroup)
	case ONEPAIR:
		return getOnePairCardsGroup(cardGroup)
	case HIGHCARD:
		return getHighCardsGroup(cardGroup)
	}
	return nil
}

// 获取皇家同花顺牌组
func getRoyalflushCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getRoyalflushCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, colorA := isFlush(cardGroup)
	if !flag {
		mlog.Error("it's not RoyalFlush")
		return nil
	}
	for i := 0; i < 5; i++ {
		cd := new(poker.Card)
		cd.Color = colorA - 1
		cd.CardNumber = int32(14 - i)
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
	}
	return returnCardsGroup
}

// 获取同花顺牌组
func getStrightflushCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getStrightflushCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA, colorA := isStrightFlush(cardGroup)
	if !flag {
		mlog.Error("it's not StrightFlush")
		return nil
	}
	for i := 0; i < 5; i++ {
		cd := new(poker.Card)
		cd.Color = colorA - 1
		cd.CardNumber = maxA - int32(i)
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
	}
	return returnCardsGroup
}

// 获取四条牌组
func getFourofakindCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getFourofakindCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA := isFourOfAKind(cardGroup)
	if !flag {
		mlog.Error("it's not FourOfAKing")
		return nil
	}
	for i := 0; i < 4; i++ {
		cd := new(poker.Card)
		cd.Color = int32(i)
		cd.CardNumber = maxA
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
	}
	CardSort(cardGroup)
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			continue
		}
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}
	return returnCardsGroup
}

// 获取葫芦牌组
func getFullhouseCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getFullhouseCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA, twoA := isFullHouse(cardGroup)
	if !flag {
		mlog.Error("it's not FullHourse")
		return nil
	}
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			cd := new(poker.Card)
			cd.Color = card.Color
			cd.CardNumber = card.CardNumber
			returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		}
		if card.CardNumber == twoA {
			cd := new(poker.Card)
			cd.Color = card.Color
			cd.CardNumber = card.CardNumber
			returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
			if returnCardsGroup.Len() >= 5 {
				return returnCardsGroup
			}
		}
	}
	return returnCardsGroup
}

// 获取同花牌组
func getFlushCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getFlushCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, colorA := isFlush(cardGroup)
	if !flag {
		mlog.Error("it's not Flush")
		return nil
	}
	cardSameColor := new(poker.CardGroup)
	for _, cardTmp := range cardGroup.Cards {
		if cardTmp.Color == colorA-1 {
			cardSameColor.Cards = append(cardSameColor.Cards, cardTmp)
		}
	}
	CardSort(cardSameColor)
	for _, card := range cardSameColor.Cards {
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}
	return returnCardsGroup
}

// 获取顺子牌组
func getStrightCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getStrightCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA := isStright(cardGroup)
	if !flag {
		mlog.Error("it's not Stright")
		return nil
	}
	for i := 0; i < 5; i++ {
		for _, card := range cardGroup.Cards {
			if card.CardNumber == maxA-int32(i) {
				cd := new(poker.Card)
				cd.Color = card.Color
				cd.CardNumber = card.CardNumber
				returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
				break
			}
		}
	}
	if maxA == 5 {
		for _, card := range cardGroup.Cards {
			if card.CardNumber == 14 {
				cd := new(poker.Card)
				cd.Color = card.Color
				cd.CardNumber = card.CardNumber
				returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
				break
			}
		}
	}
	return returnCardsGroup
}

// 获取三条牌组
func getThreeofakindCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getThreeofakindCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA := isThreeOfAKind(cardGroup)
	if !flag {
		mlog.Error("it's not Threeofakind")
		return nil
	}
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			cd := new(poker.Card)
			cd.Color = card.Color
			cd.CardNumber = card.CardNumber
			returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		}
	}
	CardSort(cardGroup)
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			continue
		}
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}
	return returnCardsGroup
}

// 获取两对牌组
func getTwoPairCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getTwoPairCardsGroup begin.")
	if cardGroup == nil {
		mlog.Error("cardGroup is nil.")
		return nil
	}
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA, minA := isTwoPair(cardGroup)
	if !flag {
		mlog.Error("it's not TwoPair.")
		return nil
	}
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA || card.CardNumber == minA {
			cd := new(poker.Card)
			cd.Color = card.Color
			cd.CardNumber = card.CardNumber
			returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		}
	}
	CardSort(cardGroup)
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA || card.CardNumber == minA {
			continue
		}
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}

	return returnCardsGroup
}

// 获取一对牌组
func getOnePairCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getOnePairCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	flag, maxA := isOnePair(cardGroup)
	if !flag {
		mlog.Error("it's not one pair.")
		return nil
	}
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			cd := new(poker.Card)
			cd.Color = card.Color
			cd.CardNumber = card.CardNumber
			returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		}
	}
	CardSort(cardGroup)
	for _, card := range cardGroup.Cards {
		if card.CardNumber == maxA {
			continue
		}
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}
	return returnCardsGroup
}

// 获取高牌牌组
func getHighCardsGroup(cardGroup *poker.CardGroup) *poker.CardGroup {
	mlog.Debug("getHighCardsGroup begin.")
	returnCardsGroup := new(poker.CardGroup)
	CardSort(cardGroup)
	for _, card := range cardGroup.Cards {
		cd := new(poker.Card)
		cd.Color = card.Color
		cd.CardNumber = card.CardNumber
		returnCardsGroup.Cards = append(returnCardsGroup.Cards, cd)
		if returnCardsGroup.Len() >= 5 {
			return returnCardsGroup
		}
	}
	return returnCardsGroup
}

// 牌组排序
func CardSort(cardGroup *poker.CardGroup) {
	length := cardGroup.Len()
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if cardGroup.Cards[i].CardNumber < cardGroup.Cards[j].CardNumber {
				cardGroup.Cards[i].CardNumber, cardGroup.Cards[j].CardNumber = cardGroup.Cards[j].CardNumber, cardGroup.Cards[i].CardNumber
				cardGroup.Cards[i].Color, cardGroup.Cards[j].Color = cardGroup.Cards[j].Color, cardGroup.Cards[i].Color
			}
		}
	}
}

// 深度拷贝牌组
func copyCardGroup(cards *poker.CardGroup) *poker.CardGroup {
	if cards == nil {
		return nil
	}
	newCards := new(poker.CardGroup)
	err := DeepCopy(newCards, cards)
	if err != nil {
		return nil
	}
	return newCards
}

// 深度拷贝
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
