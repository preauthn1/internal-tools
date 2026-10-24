//==================================================================
//文件名： ttztype.go
//说明：   推筒子牌型检测
//修改记录：
//==================================================================
package ttz

import (	
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"	
)

import (
	"fmt"
	"sort"
)

const (
	HIGH    = 0 //点数牌
	PAIR    = 1 //对子
	LEOPARD = 2 //豹子
)

// 获取牌型
func strCardType(Type, Count int32) string {
	if Type == HIGH {
		return fmt.Sprintf("%d点", Count)
	} else if Type == PAIR {
		return fmt.Sprintf("对%d", Count)
	} else if Type == LEOPARD {
		return "豹子"
	}
	return "未知牌型"
}

//是否是豹子，如果是返回true和豹子点数
func isLeopard(cg *poker.CardGroup) bool {
	if cg == nil || cg.Len() != 2 {
		return false
	}
	if cg.Cards[0].CardNumber == cg.Cards[1].CardNumber && cg.Cards[0].Color == int32(BAIBAN) && cg.Cards[1].Color == int32(BAIBAN) {
		return true
	}
	return false
}

//是否是对子
func isPair(cg *poker.CardGroup) (bool, int32) {
	if cg == nil || cg.Len() != 2 {
		return false, 0
	}
	if cg.Cards[0].CardNumber == cg.Cards[1].CardNumber && cg.Cards[0].Color == int32(TONGZI) {
		return true, cg.Cards[0].CardNumber
	}
	return false, 0
}

//获取点数
func getCount(cg *poker.CardGroup) int32 {
	if cg == nil || cg.Len() != 2 {
		return 0
	}
	Count := (cg.Cards[0].CardNumber + cg.Cards[1].CardNumber) % 10
	return Count
}

// 获取牌型和点数
func getCardTypeAndCount(cg *poker.CardGroup) (cardType, count int32) {
	if cg.Len() < 2 {
		mlog.Error("CardGroup size[%d] error!", cg.Len())
		return 0, 0
	}
	defer func() {
		mlog.Debug("获取牌型,cardType[%d],count[%d]", cardType, count)
	}()
	if isLeopard(cg) {
		return LEOPARD, 0
	}
	flag, count := isPair(cg)
	if flag {
		return PAIR, count
	}
	count = getCount(cg)
	return HIGH, count
}

// 比较牌型大小
func compareCards(cgleft, cgright *poker.CardGroup) int32 {
	if cgleft == nil || cgright == nil {
		return 0
	}
	cardTypeLeft, countLeft := getCardTypeAndCount(cgleft)
	cardTypeRight, countRight := getCardTypeAndCount(cgright)
	if cardTypeLeft > cardTypeRight {
		return 1
	}
	if cardTypeLeft < cardTypeRight {
		return -1
	}
	switch cardTypeLeft {
	case LEOPARD:
		return 0
	case PAIR:
		if countLeft > countRight {
			return 1
		}
		if countLeft < countRight {
			return -1
		}
		return 0
	case HIGH:
		if countLeft > countRight {
			return 1
		}
		if countLeft < countRight {
			return -1
		}
		return compareSameHigh(cgleft, cgright)
	}
	return 0
}

// 比较相同牌型大小
func compareSameHigh(cgleft, cgright *poker.CardGroup) int32 {
	if cgleft == nil || cgright == nil {
		return 0
	}
	if cgleft.Len() != 2 || cgright.Len() != 2 {
		return 0
	}

	sort.Slice(cgleft.Cards[:], func(i, j int) bool {
		return cgleft.Cards[i].CardNumber > cgleft.Cards[j].CardNumber
	})
	sort.Slice(cgright.Cards[:], func(i, j int) bool {
		return cgright.Cards[i].CardNumber > cgright.Cards[j].CardNumber
	})

	if cgleft.Cards[0].CardNumber > cgright.Cards[0].CardNumber {
		return 1
	}
	if cgleft.Cards[0].CardNumber < cgright.Cards[0].CardNumber {
		return -1
	}
	return 0

}

// 获取牌型倍率
func getRateByCardType(cardType, count int32) int32 {
	switch cardType {
	case LEOPARD:
		return 4
	case PAIR:
		return 3
	case HIGH:
		if count >= 7 {
			return 2
		} else {
			return 1
		}
	default:
		return 1
	}	
}
