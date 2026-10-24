package wzp

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

//是否同花大顺
func isTonghuaDashun(cg *poker.CardGroup) bool {
	return isTonghuaShun(cg) && cg.GetMinCardNumber() == 10
}

//是否同花顺
func isTonghuaShun(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 5 {
		return false
	}

	return cg.IsTongHuaShun()
}

//是否同花
func isTongHua(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 5 {
		return false
	}

	return cg.IsTonghua()
}

//是否顺子
func isShunzi(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 5 {
		return false
	}

	return cg.IsShunzi()
}

//是否四条
func isSitiao(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 4 {
		return false
	}

	min := cg.GetMinCardNumber()
	max := cg.GetMaxCardNumber()

	if len(cg.GetPokerCardByNumber(min)) == 4 {
		return true
	}

	if len(cg.GetPokerCardByNumber(max)) == 4 {
		return true
	}

	return false
}

//是否三条
func isSanTiao(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 3 {
		return false
	}

	min := cg.GetMinCardNumber()
	max := cg.GetMaxCardNumber()

	for i := min; i <= max; i++ {
		if len(cg.GetPokerCardByNumber(i)) == 3 {
			return true
		}
	}

	return false
}

//是否葫芦
func isHulu(cg *poker.CardGroup) bool {
	if len(cg.Cards) != 5 {
		return false
	}

	cnm := getPokerNumber(cg)

	if len(cnm) != 2 {
		return false
	}

	for _, cn := range cnm {
		if cn != 3 && cn != 2 {
			return false
		}
	}

	return true
}

//是否两对
func isTwoPair(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 4 {
		return false
	}

	cnm := getPokerNumber(cg)

	if len(cnm) < 2 {
		return false
	}

	twopaircn := 0
	for _, cn := range cnm {
		if cn != 1 && cn != 2 {
			return false
		}

		if cn == 2 {
			twopaircn++
		}
	}

	return twopaircn >= 2
}

//是否一对
func isOnePair(cg *poker.CardGroup) bool {
	if len(cg.Cards) < 2 {
		return false
	}

	cnm := getPokerNumber(cg)

	if len(cnm) != len(cg.Cards)-1 {
		return false
	}

	return true
}

//是否散牌
func isSanpai(cg *poker.CardGroup) bool {
	cnm := getPokerNumber(cg)

	return len(cnm) == len(cg.Cards)
}

//获取扑克点数
func getPokerNumber(cg *poker.CardGroup) map[int32]int32 {
	cnm := make(map[int32]int32)

	for _, c := range cg.Cards {
		if _, ok := cnm[c.CardNumber]; ok {
			cnm[c.CardNumber]++
		} else {
			cnm[c.CardNumber] = 1
		}
	}

	return cnm
}

//根据重复盘点数获取牌点
func getCardNumberBySameCount(cg *poker.CardGroup, n int32) int32 {
	s := getPokerNumber(cg)

	for cno, cn := range s {
		if cn == n {
			return cno
		}
	}

	return 0
}

//获取两对牌型的点数
func getTwoPairNo(cg *poker.CardGroup) (int32, int32) {
	cgn := getPokerNumber(cg)

	pr_cno := make([]int32, 0)

	for cno, cn := range cgn {
		if cn == 2 {
			pr_cno = append(pr_cno, cno)
		}
	}

	if pr_cno[0] > pr_cno[1] {
		return pr_cno[0], pr_cno[1]
	} else {
		return pr_cno[1], pr_cno[0]
	}
}

//是否含有黑桃
func containsHeitao(cg []*poker.Card) bool {
	for _, c := range cg {
		if c.Color == poker.PokerColorHeiTao {

			return true
		}
	}

	return false
}

//获取牌型
//同花大顺 1 同花顺 2 四条 3 葫芦 4 同花 5 顺子 6 三条 7 两对 8 一对 9 散牌 10 弃牌 11
func getCardSize(cg *poker.CardGroup) int32 {
	if isTonghuaDashun(cg) {
		return 1
	}
	if isTonghuaShun(cg) {
		return 2
	}
	if isSitiao(cg) {
		return 3
	}
	if isHulu(cg) {
		return 4
	}
	if isTongHua(cg) {
		return 5
	}
	if isShunzi(cg) {
		return 6
	}
	if isSanTiao(cg) {
		return 7
	}
	if isTwoPair(cg) {
		return 8
	}
	if isOnePair(cg) {
		return 9
	}

	return 10

}

//获取牌型string
func getCardTypeString(cg *poker.CardGroup) string {
	if isTonghuaDashun(cg) {
		return "同花大顺"
	}
	if isTonghuaShun(cg) {
		return "同花顺"
	}
	if isSitiao(cg) {
		return "四条"
	}
	if isHulu(cg) {
		return "葫芦"
	}
	if isTongHua(cg) {
		return "同花"
	}
	if isShunzi(cg) {
		return "顺子"
	}
	if isSanTiao(cg) {
		return "三条"
	}
	if isTwoPair(cg) {
		return "两对"
	}
	if isOnePair(cg) {
		return "一对"
	}

	return "散牌"

}

//比较顺子大小
func compareShunzi(cg1 *poker.CardGroup, cg2 *poker.CardGroup) bool {
	max1 := cg1.GetMaxCard()
	max2 := cg2.GetMaxCard()

	if cg1.IsShunziFromOne() {
		max1 = cg1.GetPokerCardByNumber(5)[0]
	}
	if cg2.IsShunziFromOne() {
		max2 = cg2.GetPokerCardByNumber(5)[0]
	}

	if max1.CardNumber != max2.CardNumber {
		return max1.CardNumber > max2.CardNumber
	}

	return max1.Color > max2.Color
}

//比较牌型
func compareCard(cg1 *poker.CardGroup, cg2 *poker.CardGroup) (bool, int32, int32) {
	sz1 := getCardSize(cg1)
	sz2 := getCardSize(cg2)

	mlog.Debug("cmp ret:%s:%d,%s:%d", cg1, sz1, cg2, sz2)

	ret := false
	if sz1 != sz2 {
		ret = sz1 < sz2
	} else {
		//同花大顺
		if sz1 == 1 {
			ret = cg1.GetMaxCard().Color > cg2.GetMaxCard().Color
		} else if sz1 == 2 { //同花顺
			c1 := cg1.GetMaxCard()
			c2 := cg2.GetMaxCard()

			if c1.CardNumber != c2.CardNumber {
				ret = c1.CardNumber > c2.CardNumber
			} else {
				ret = c1.Color > c2.Color
			}
		} else if sz1 == 3 { //四条
			cn1 := getCardNumberBySameCount(cg1, 4)
			cn2 := getCardNumberBySameCount(cg2, 4)

			ret = cn1 > cn2
		} else if sz1 == 4 { //葫芦
			cn1 := getCardNumberBySameCount(cg1, 3)
			cn2 := getCardNumberBySameCount(cg2, 3)

			ret = cn1 > cn2
		} else if sz1 == 5 { //同花
			c1 := cg1.GetMaxCard()
			c2 := cg2.GetMaxCard()

			if c1.CardNumber != c2.CardNumber {
				ret = c1.CardNumber > c2.CardNumber
			}

			ret = c1.Color > c2.Color
		} else if sz1 == 6 { //顺子
			ret = compareShunzi(cg1, cg2)
		} else if sz1 == 7 { //三条
			cn1 := getCardNumberBySameCount(cg1, 3)
			cn2 := getCardNumberBySameCount(cg2, 3)

			ret = cn1 > cn2
		} else if sz1 == 8 { //两对
			p1_1, p1_2 := getTwoPairNo(cg1)
			p2_1, p2_2 := getTwoPairNo(cg2)

			if p1_1 == p2_1 {
				if p1_2 == p2_2 {
					//完全一样的点数,比较花色
					p1 := cg1.GetPokerCardByNumber(p1_1)

					if containsHeitao(p1) {
						ret = true
					} else {
						ret = false
					}
				} else {
					ret = p1_2 > p2_2
				}
			} else {
				ret = p1_1 > p2_1
			}
		} else if sz1 == 9 { //一对
			cn1 := getCardNumberBySameCount(cg1, 2)
			cn2 := getCardNumberBySameCount(cg2, 2)

			if cn1 != cn2 {
				ret = cn1 > cn2
			} else {
				cnp1 := cg1.GetPokerCardByNumber(cn1)

				if containsHeitao(cnp1) {
					ret = true
				} else {
					ret = false
				}
			}
		} else { //散牌
			c1 := cg1.GetMaxCard()
			c2 := cg2.GetMaxCard()

			if c1.CardNumber != c2.CardNumber {
				ret = c1.CardNumber > c2.CardNumber
				mlog.Debug("cmp ret 9 num:%d,%d", c1.CardNumber, c2.CardNumber)
			} else {
				ret = c1.Color > c2.Color
				mlog.Debug("cmp ret 9 color:%d,%d", c1.Color, c2.Color)
			}
		}
	}

	return ret, sz1, sz2
}

//获取当前游戏最大的用户牌型
func getGameMaxCard(pk *poker.GamePoker, trunno int32, glc *WZPGameLogic) (*poker.UserPoker, map[int32]int32) {
	var max *poker.UserPoker
	cardtype := make(map[int32]int32, 0)

	for _, up := range pk.UserPokers {
		if glc.turnManager.isGiveUp(up.User.UserID) {
			continue
		}

		if max == nil {
			max = up
			continue
		}

		if max == up {
			continue
		}

		cmpcards := poker.NewCardGroup(nil)
		cmpcards.AddCardGroup(up.Cards)

		maxcards := poker.NewCardGroup(nil)
		maxcards.AddCardGroup(max.Cards)

		if trunno <= 4 {
			cmpcards.Cards = cmpcards.Cards[1:]
			maxcards.Cards = maxcards.Cards[1:]
		}

		ret, sz1, sz2 := compareCard(maxcards, cmpcards)

		cardtype[max.User.UserID] = sz1
		cardtype[up.User.UserID] = sz2

		if !ret {
			max = up
		}

	}

	return max, cardtype
}
