package hhdz

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"sort"
)

const (
	DanPai = iota
	DuiZi
	SunZi
	JiaHua
	SunJin
	BaoZi
)

var special bool

//是否是豹子
func isBaozi(Cards *poker.CardGroup) bool {
	if Cards.Cards[0].CardNumber == Cards.Cards[1].CardNumber && Cards.Cards[1].CardNumber == Cards.Cards[2].CardNumber {
		return true
	}

	return false
}

//是否是金花
func isJinHua(Cards *poker.CardGroup) bool {
	if Cards.Cards[0].Color == Cards.Cards[1].Color && Cards.Cards[1].Color == Cards.Cards[2].Color {
		return true
	}
	return false
}

//是否是顺子
func isSunZi(Cards *poker.CardGroup) bool {
	cloneCards := poker.NewCardGroup(nil)
	cloneCards.AddCardGroup(Cards)
	sort.Slice(cloneCards.Cards, func(i, j int) bool { return cloneCards.Cards[i].CardNumber < cloneCards.Cards[j].CardNumber })
	if cloneCards.Cards[0].CardNumber+1 == cloneCards.Cards[1].CardNumber && cloneCards.Cards[1].CardNumber+1 == cloneCards.Cards[2].CardNumber {
		return true
	}

	//特殊123
	if cloneCards.Cards[1].CardNumber == 2 && cloneCards.Cards[2].CardNumber == 3 && cloneCards.Cards[0].CardNumber == 14 {
		special = true
		return true
	} else if cloneCards.Cards[0].CardNumber == 2 && cloneCards.Cards[1].CardNumber == 3 && cloneCards.Cards[2].CardNumber == 14 {
		special = true
		return true
	}
	//特殊AOK
	if cloneCards.Cards[0].CardNumber == 14 && cloneCards.Cards[1].CardNumber == 12 && cloneCards.Cards[2].CardNumber == 13 {
		return true
	}
	return false
}

//是否是对子
func isDuiZi(Cards *poker.CardGroup) bool {
	if Cards.Cards[0].CardNumber == Cards.Cards[1].CardNumber || Cards.Cards[1].CardNumber == Cards.Cards[2].CardNumber ||
		Cards.Cards[0].CardNumber == Cards.Cards[2].CardNumber {
		return true
	}

	return false
}

//获取牌型
func GetType(Cards *poker.CardGroup) int32 {
	//默认已经排号序
	if Cards == nil || len(Cards.Cards) <= 0 {
		mlog.Debug("牌是空的")
		return DanPai
	}

	if isJinHua(Cards) {
		if isSunZi(Cards) {
			return SunJin
		}
		return JiaHua
	}

	if isSunZi(Cards) {
		return SunZi
	}

	if isBaozi(Cards) {
		return BaoZi
	}

	if isDuiZi(Cards) {
		return DuiZi
	}

	return DanPai
}

//比较用户牌型
func comPareCardsByNoAndColor(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int32 {
	//遍历比较点数大小
	if Cards1.Cards[0].CardNumber == 14 && Cards2.Cards[0].CardNumber != 14 {
		return 1
	} else if Cards1.Cards[0].CardNumber != 14 && Cards2.Cards[0].CardNumber == 14 {
		return 2
	}
	for i := 2; i >= 0; i-- {
		if Cards1.Cards[i].CardNumber > Cards2.Cards[i].CardNumber {
			return 1
		}
		if Cards1.Cards[i].CardNumber < Cards2.Cards[i].CardNumber {
			return 2
		}
	}
	return comPareCardsByColor(Cards1, Cards2)
}

//最大那个牌的花色
func comPareCardsByColor(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int32 {
	//1會排在最前面
	if Cards1.Cards[0].CardNumber == 1 {
		if Cards1.Cards[0].Color > Cards2.Cards[0].Color {
			return 1
		}
	} else if Cards1.Cards[2].Color > Cards2.Cards[2].Color {
		return 1
	}
	return 2
}

//返回1为card1用户赢
func comPareCards(cg1 *poker.CardGroup, cg2 *poker.CardGroup) int32 {
	Cards1 := poker.NewCardGroup(nil)
	Cards1.AddCardGroup(cg1)
	Cards2 := poker.NewCardGroup(nil)
	Cards2.AddCardGroup(cg2)
	sort.Slice(Cards1.Cards, func(i, j int) bool { return Cards1.Cards[i].CardNumber < Cards1.Cards[j].CardNumber })
	sort.Slice(Cards2.Cards, func(i, j int) bool { return Cards2.Cards[i].CardNumber < Cards2.Cards[j].CardNumber })
	type1 := GetType(Cards1)
	type2 := GetType(Cards2)

	if type1 > type2 {
		return 1
	}

	if type1 < type2 {
		return 2
	}

	if type1 == SunJin && special {
		if Cards1.Cards[0].CardNumber == 2 && Cards1.Cards[1].CardNumber == 3 && Cards1.Cards[2].CardNumber == 14 && Cards2.Cards[0].CardNumber != 2 {
			special = false
			return 2
		} else if Cards2.Cards[0].CardNumber == 2 && Cards2.Cards[1].CardNumber == 3 && Cards2.Cards[2].CardNumber == 14 {
			special = false
			return 1
		}
	}

	if type1 == SunZi && special {
		if Cards1.Cards[0].CardNumber == 2 && Cards1.Cards[1].CardNumber == 3 && Cards1.Cards[2].CardNumber == 14 && Cards2.Cards[0].CardNumber != 2 {
			special = false
			return 2
		} else if Cards2.Cards[0].CardNumber == 2 && Cards2.Cards[1].CardNumber == 3 && Cards2.Cards[2].CardNumber == 14 {
			special = false
			return 1
		}
	}

	if type1 == DuiZi {
		card1_tmp := Cards1.Cards[0]
		if Cards1.Cards[1].CardNumber == Cards1.Cards[2].CardNumber {
			card1_tmp = Cards1.Cards[1]
		}

		card2_tmp := Cards2.Cards[0]
		if Cards2.Cards[1].CardNumber == Cards2.Cards[2].CardNumber {
			card2_tmp = Cards2.Cards[1]
		}

		if card1_tmp.CardNumber > card2_tmp.CardNumber {
			return 1
		}

		if card1_tmp.CardNumber < card2_tmp.CardNumber {
			return 2
		}
	}
	return comPareCardsByNoAndColor(Cards1, Cards2)
}
