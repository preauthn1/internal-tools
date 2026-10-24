package tbnn

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/poker"
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
		c1n := getNiuNiuCardNumber(c1)
		c2 := sc2.cardGroup.Cards[i]
		c2n := getNiuNiuCardNumber(c2)

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

			//是否開啟百搭鬼牌 0 關閉 1 開啟
			t := 0
			if nnconfig.baidakingConfig == 0 {
				if val%10 == 0 {
					t = 1
				}
			} else {
				//若百搭，檢驗最後兩張是否為皮子或大小王
				if val%10 == 0 || cards.Cards[i].IsPiZi() || cards.Cards[i].IsBigKing() || cards.Cards[i].IsSmallKing() ||
					cards.Cards[j].IsPiZi() || cards.Cards[j].IsBigKing() || cards.Cards[j].IsSmallKing() {
					t = 1
				}
			}
			if t == 1 {
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
