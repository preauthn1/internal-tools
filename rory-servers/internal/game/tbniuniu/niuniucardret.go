package tbnn

import (	
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

import (
	"fmt"
)

const (
	nnCard_WuXiaoNiu = iota //五小牛
	nnCard_4Zha             //四炸
	nnCard_WuHuaNiu         //五花牛
	nnCard_NiuNiu           //牛牛
	nnCard_Niu              //有牛
	nnCard_WuNiu            //无牛
)

//牛牛牌型结果
type niuniuCardRet struct {
	niuRet    int              //牌型结果
	number    int              //牌点
	cardStack *poker.CardStack //组牌方式
	cards     *poker.CardGroup //牌组
}

//创建牛牛牌型结果对象
func newNiuNiuCardRet(nnRet int, number int, cards *poker.CardGroup) *niuniuCardRet {
	nncret := new(niuniuCardRet)
	nncret.niuRet = nnRet
	nncret.number = number
	nncret.cards = cards

	mlog.Debug("%s,%d,%d", cards.String(), nnRet, number)

	return nncret
}

//是否有牛
func (nncardret *niuniuCardRet) checkNiu() bool {
	return nncardret.niuRet < nnCard_WuNiu
}

//是否相同牌型
func (nncardret *niuniuCardRet) isSameRet(otherret *niuniuCardRet) bool {
	return nncardret.niuRet == otherret.niuRet
}

//获取牌型倍数
func (nncardret *niuniuCardRet) getBetTimes() int32 {
	//五小牛
	if nncardret.niuRet == nnCard_WuXiaoNiu {
		return 6
	} else if nncardret.niuRet == nnCard_4Zha { //4炸
		return 5
	} else if nncardret.niuRet == nnCard_WuHuaNiu { //五花牛
		return 4
	} else if nncardret.niuRet == nnCard_NiuNiu { //牛牛
		return 3
	} else if nncardret.niuRet == nnCard_Niu { //有牛
		if nncardret.number <= 6 {
			return 1
		} else {
			return 2
		}
	} else if nncardret.niuRet == nnCard_WuNiu { //无牛
		return 1
	}

	return 1
}

//比较两个牌型结果 1 大于 0 等于 -1 小于
func (nncardret *niuniuCardRet) compare(otherRet *niuniuCardRet) int32 {
	//牌型不等的情况
	if nncardret.niuRet < otherRet.niuRet {
		return 1
	} else if nncardret.niuRet > otherRet.niuRet {
		return -1
	}
	//牌型相等的情况
	return 0
}

func (nncardret *niuniuCardRet) getNetCardGroup() []*netproto.PokerCard {
	ncg := make([]*netproto.PokerCard, 0)

	if nncardret.cardStack != nil {
		for i := len(nncardret.cardStack.Groups) - 1; i >= 0; i-- {
			cgs := nncardret.cardStack.Groups[i]

			for _, cg := range cgs.Cards {
				ncg = append(ncg, cg.GetNetMessageData())
			}
		}
	} else {
		for _, cg := range nncardret.cards.Cards {
			ncg = append(ncg, cg.GetNetMessageData())
		}
	}

	return ncg
}

func (nncardret *niuniuCardRet) String() string {
	retstr := ""
	if nncardret.niuRet == nnCard_WuXiaoNiu {
		retstr = "五小牛"
	} else if nncardret.niuRet == nnCard_WuHuaNiu {
		retstr = "五花牛"
	} else if nncardret.niuRet == nnCard_NiuNiu {
		retstr = "牛牛"
	} else if nncardret.niuRet == nnCard_Niu {
		retstr = fmt.Sprintf("牛%d", nncardret.number)
	} else if nncardret.niuRet == nnCard_4Zha {
		retstr = "四炸"
	} else {
		retstr = "无牛"
	}

	cardstr := ""
	if nncardret.cardStack != nil {
		cardstr = nncardret.cardStack.GetCardGroup().String()
	} else {
		cardstr = nncardret.cards.String()
	}

	return fmt.Sprintf("%s,%s", retstr, cardstr)
}
