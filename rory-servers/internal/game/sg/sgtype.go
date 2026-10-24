package sg

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

const (
	DanPai      = 0
	TeShu       = 1
	HunSanGong  = 10
	XiaoSanGong = 11
	DaSanGong   = 12
)

//把A的点数变为1
func changeA(Cards *poker.CardGroup) bool {
	if Cards == nil || len(Cards.Cards) != 3 {
		mlog.Debug("changeA cards = nil")
		return false
	}
	for i := 0; i < 3; i++ {
		if Cards.Cards[i].IsA() {
			Cards.Cards[i].CardNumber = 1
		}
	}
	return true
}

//给扑克排序
func sortCards(Cards *poker.CardGroup) {
	for i := 0; i < 3; i++ {
		for j := i; j < 2; j++ {
			if Cards.Cards[j].CardNumber > Cards.Cards[j+1].CardNumber {
				Cards.Swap(j, j+1)
			}
		}
	}
}

//获取三公的类型，默认已排序并且将A替换为1
func getSanGongType(Cards *poker.CardGroup) int32 {
	if Cards.Cards[0].CardNumber == Cards.Cards[1].CardNumber && Cards.Cards[2].CardNumber == Cards.Cards[1].CardNumber {
		if Cards.Cards[0].CardNumber > 10 {
			return DaSanGong
		}

		return XiaoSanGong
	}

	if Cards.Cards[0].CardNumber > 10 && Cards.Cards[1].CardNumber > 10 && Cards.Cards[2].CardNumber > 10 {
		return HunSanGong
	}

	var dian int32 = 0
	for i := 0; i < 3; i++ {
		if Cards.Cards[i].CardNumber < 10 {
			dian += Cards.Cards[i].CardNumber
		}
	}

	dian = dian % 10
	return dian
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

//比较牌型
func comPareCards(Cards1 *poker.CardGroup, Cards2 *poker.CardGroup) int32 {
	tmpCards1 := copyPokerCards(Cards1)
	tmpCards2 := copyPokerCards(Cards2)

	ok := changeA(tmpCards1)
	if !ok {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("cards1 : %v", tmpCards1)).Error("撲克牌有少, 未滿三張")
		return 1
	}
	ok = changeA(tmpCards2)
	if !ok {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("cards2 : %v", tmpCards2)).Error("撲克牌有少, 未滿三張")
		return 1
	}
	sortCards(tmpCards1)
	sortCards(tmpCards2)
	type1 := getSanGongType(tmpCards1)
	type2 := getSanGongType(tmpCards2)

	if type1 > type2 {
		return 1
	}

	if type1 < type2 {
		return 2
	}

	for i := 2; i >= 0; i-- {
		if tmpCards1.Cards[i].CardNumber > tmpCards2.Cards[i].CardNumber {
			return 1
		}

		if tmpCards1.Cards[i].CardNumber < tmpCards2.Cards[i].CardNumber {
			return 2
		}
	}

	if tmpCards1.Cards[2].Color > tmpCards2.Cards[2].Color {
		return 1
	}

	if tmpCards1.Cards[2].Color < tmpCards2.Cards[2].Color {
		return 2
	}

	mlog.Debug("还没决出胜负:玩家1的牌%v，玩家2的牌%v", tmpCards1.Cards, tmpCards2.Cards)
	return 1
}
