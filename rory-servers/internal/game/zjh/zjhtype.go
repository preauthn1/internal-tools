package zjh

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"	
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

const (
	SanPai = iota
	DuiZi
	SunZi
	JiaHua
	SunJin
	BaoZi
)

var typeStringMap = map[int]string{
	SanPai : "单牌",
	DuiZi : "对子",
	SunZi : "顺子",
	JiaHua : "金花",
	SunJin : "顺金",
	BaoZi : "豹子",
}

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
	if Cards.Cards[0].CardNumber+1 == Cards.Cards[1].CardNumber && Cards.Cards[1].CardNumber+1 == Cards.Cards[2].CardNumber {
		return true
	}

	//特殊123
	if Cards.Cards[0].CardNumber == 2 && Cards.Cards[1].CardNumber == 3 && Cards.Cards[2].CardNumber == 14 {
		special = true
		return true
	}
	//特殊AQK
	if Cards.Cards[0].CardNumber == 14 && Cards.Cards[1].CardNumber == 12 && Cards.Cards[2].CardNumber == 13 {
		return true
	}
	return false
}

//是否是对子
func isDuiZi(Cards *poker.CardGroup) bool {
	if Cards.Cards[0].CardNumber == Cards.Cards[1].CardNumber || Cards.Cards[1].CardNumber == Cards.Cards[2].CardNumber {
		return true
	}

	return false
}

//获取牌型
func GetType(Cards *poker.CardGroup) int32 {
	//默认已经排号序
	if Cards == nil || len(Cards.Cards) <= 0 {
		mlog.Debug("牌是空的")
		return SanPai
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

	return SanPai
}

//比较用户牌型
func comPareCardsByNoAndColor(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int {
	if Cards1.Cards[2].CardNumber == 14 && Cards2.Cards[2].CardNumber != 14 {
		return 1
	} else if Cards1.Cards[2].CardNumber != 14 && Cards2.Cards[2].CardNumber == 14 {
		return 2
	}
	//遍历比较点数大小
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
func comPareCardsByColor(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int {
	if Cards1.Cards[2].Color > Cards2.Cards[2].Color {
		return 1
	}
	return 2
}

//返回1为card1用户赢
func comPareCards(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int {
	type1 := GetType(Cards1)
	type2 := GetType(Cards2)

	//特殊情况235和豹子
	if type1 == BaoZi && type2 == SanPai {
		if Cards2.Cards[0].CardNumber == 2 && Cards2.Cards[1].CardNumber == 3 && Cards2.Cards[2].CardNumber == 5 {
			return 2
		}
	}

	if type2 == BaoZi && type1 == SanPai {
		if Cards1.Cards[0].CardNumber == 2 && Cards1.Cards[1].CardNumber == 3 && Cards1.Cards[2].CardNumber == 5 {
			return 1
		}
	}

	if type1 > type2 {
		return 1
	}

	if type1 < type2 {
		return 2
	}

	if (type1 == SunJin || type1 == SunZi) && special {
		if Cards1.Cards[0].CardNumber == 2 && Cards1.Cards[1].CardNumber == 3 && Cards1.Cards[2].CardNumber == 14 {
			//特殊123為最小牌型，非123的順子為贏家
			if !(Cards2.Cards[0].CardNumber == 2 && Cards2.Cards[1].CardNumber == 3 && Cards2.Cards[2].CardNumber == 14) {
				special = false
				return 2
			}
		} else if Cards2.Cards[0].CardNumber == 2 && Cards2.Cards[1].CardNumber == 3 && Cards2.Cards[2].CardNumber == 14 {
			if !(Cards1.Cards[0].CardNumber == 2 && Cards1.Cards[1].CardNumber == 3 && Cards1.Cards[2].CardNumber == 14) {
				special = false
				return 1
			}
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
