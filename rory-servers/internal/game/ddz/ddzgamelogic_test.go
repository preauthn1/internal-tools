package ddz

import (
	"os"
	"path"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/mylog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//outCardByType 出牌牌型 手排最後是大王
func TestOutCardByType31(t *testing.T) {
	ddz := NewDDZGameLogicByTest("lastCardType4")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出火箭
	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	if ckType != cardTypeSingle || autoCardNumber != int32(16) {
		t.Errorf("出牌牌型 手排最後是大王, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//outCardByType 出牌牌型 手排最後是火箭
func TestOutCardByType(t *testing.T) {
	ddz := NewDDZGameLogicByTest("lastCardType1")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出火箭
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	if ckType != cardTypeRocket || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 手排最後是火箭, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//outCardByType 出牌牌型 手排最後是炸彈
func TestOutCardByType2(t *testing.T) {
	ddz := NewDDZGameLogicByTest("lastCardType2")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//手排最後是炸彈
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeBomb || ckType != cardTypeBomb || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 炸彈, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//outCardByType 出牌牌型 手排只剩火箭跟炸彈
func TestOutCardByType3(t *testing.T) {
	ddz := NewDDZGameLogicByTest("lastCardType3")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//user=1
	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if ckType != cardTypeRocket || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 炸彈, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//連對
func TestOutLiandui(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outLiandui")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThreeStraight || ckType != cardTypeThreeStraight || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 連對, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//user=1
	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThreeStraight || ckType != cardTypeThreeStraight || autoMainCard != int32(1) || autoCardNumber != int32(8) {
		t.Errorf("出牌牌型 連對, 跟牌 壓牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}
}

//outCardByType 出牌牌型 順子
func TestOutCardByTypeStraight(t *testing.T) {
	ddz := NewDDZGameLogicByTest("Straight01")

	//先出單牌
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeStraight || ckType != cardTypeStraight || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 順子, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//user=1
	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeStraight || ckType != cardTypeStraight || autoMainCard != int32(1) || autoCardNumber != int32(14) {
		t.Errorf("出牌牌型 順子, 跟牌 壓牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}
}

//三條
func TestOutThreeCard(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outThreeCards")

	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 三順, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//三條
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThree || ckType != cardTypeThree || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 三條, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++ //user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThree || ckType != cardTypeThree || autoMainCard != int32(1) || autoCardNumber != int32(6) {
		t.Errorf("出牌牌型 三條, 跟牌 壓牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++ //user=2
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThree || ckType != cardTypeThree || autoMainCard != int32(1) || autoCardNumber != int32(2) {
		t.Errorf("出牌牌型 三條, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//user=2
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeThree || ckType != cardTypeThree || autoMainCard != int32(1) || autoCardNumber != int32(4) {
		t.Errorf("出牌牌型 三條, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//outCardByType 出牌牌型 對子
func TestOutCardByTypeDouble(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outCardByDouble")

	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出對子 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 對子, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++ //user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(6) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	//TODO 要評估對手牌, 當對手牌很多就從3開始. 當對手牌很少, 要從2 -> 6 -> 3 出牌
	mlog.Error("要評估對手牌, 當對手牌很多就從3開始. 當對手牌很少, 要從2 -> 6 -> 3 出牌")
	//user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	//if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(2) {
	//	t.Errorf("出牌牌型 對子, 主動 手上牌<=2種對子就從特殊牌最小開始出牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
	//		autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	//}

	//user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	//if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(3) {
	//	t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
	//		autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	//}

	ddz = NewDDZGameLogicByTest("outCardByDouble")
	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出對子 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	ddz.ext.turn++ //user=1
	ddz.ext.turn++ //user=2
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) || autoCardNumber != int32(2) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

func TestOutCardByTypeDouble2(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outCardByDouble2")

	//先出單牌 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)

	//出對子 user=0
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDouble || ckType != cardTypeDouble || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 對子, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++ //user=1
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	if autoCardType != cardTypeNone || autoMainCard != int32(0) || autoCardNumber != int32(0) {
		t.Errorf("出牌牌型 對子, 跟牌 壓牌(0) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//outCardByType 出牌牌型 單牌, 主動出牌
func TestOutCardByTypeDanpai(t *testing.T) {

	ddz := NewDDZGameLogicByTest("outCardByDanpai")

	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()

	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 單牌, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) || autoCardNumber != int32(2) {
		t.Errorf("出牌牌型 單牌, 跟牌 壓牌(2) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) || autoCardNumber != int32(15) {
		t.Errorf("出牌牌型 單牌, 跟牌 手上牌小於３張就從特殊牌最小開始出牌(2,小王,大王) 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	ddz.ext.turn++
	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	if autoCardType != cardTypeSingle || autoMainCard != int32(1) || autoCardNumber != int32(16) {
		t.Errorf("出牌牌型 單牌, 跟牌 2不跟 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)

	}

}

//飛機
func TestOutFeiji(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outFeiji")

	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDoubleStraightd || ckType != cardTypeDoubleStraightd || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//飛機
func TestOutFeiji2(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outFeiji2")

	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDoubleStraights || ckType != cardTypeDoubleStraights || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//飛機
func TestOutFeiji3(t *testing.T) {
	ddz := NewDDZGameLogicByTest("outFeiji3")

	autoCg, autoCardType, autoCardNumber, autoMainCard := ddz.onAutoChuPaiByTuoGuan()
	ckType := hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeSingle || ckType != cardTypeSingle || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

	autoCg, autoCardType, autoCardNumber, autoMainCard = ddz.onAutoChuPaiByTuoGuan()
	ckType = hisPushAndDeleteCard(autoCg, autoMainCard, autoCardNumber, ddz, autoCardType)
	if autoCardType != cardTypeDoubleStraightd || ckType != cardTypeDoubleStraightd || autoMainCard != int32(1) {
		t.Errorf("出牌牌型 飛機, 主動出牌 錯誤！ zCg[%v] CardType[%v] CardNumber[%v] MainCard[%v] ckType[%v]",
			autoCg, autoCardType, autoCardNumber, autoMainCard, ckType)
	}

}

//增加到出牌歷史和刪除手牌出牌資料
func hisPushAndDeleteCard(autoCg *poker.CardGroup, autoMainCard int32, autoCardNumber int32, ddz *DDZGameLogic, pCardType cardType) cardType {
	action := new(userAction)
	action.ct = pCardType
	action.cg = autoCg
	action.mCard = autoMainCard
	action.maxCard = autoCardNumber
	action.userSeat = ddz.ext.getNowSeat()
	ckType := ddz.ext.his.cardTypeCheck(autoCg, action)
	ddz.ext.his.Push(action)
	up := ddz.gamePoker.GetUserPoker(ddz.ext.getNowSeat())
	up.Cards.DeleteGroup(autoCg)
	return ckType
}

func enterTable(userLen int) *gamelogic.Table {

	table := gamelogic.NewTable(3, 0)
	for i := 0; i < userLen; i++ {
		//userid := int32(i)
		//cards := fp.GetUserCards(userid)
		//ddz.recombination(userid, cards)

		user := new(gamelogic.GameUser)
		user.SeatNo = int32(i)
		user.UserID = int32(i)

		table.Seats[i].SeatUser = user
		gameState := new(gamelogic.UserGameState)
		user.UserState = gameState
	}

	return table
}
func fixedPoker(filePoker string) *poker.FixedPoker {
	var p poker.Poker
	p = poker.NewTestFilePoker(filePoker)
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	return fp
}

//創建測試使用的GameLogic和建立撲克牌發牌器資料
func NewDDZGameLogicByTest(filePoker string) *DDZGameLogic {
	fixPoker := fixedPoker(filePoker) //讀取發牌資料 /testfixpoker/xx.dat
	table := enterTable(fixPoker.UserCardLen())
	tfManager := new(event.TimerFuncManager)

	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)

	glc := new(DDZGameLogic)
	glc.FixedGameLogicBase = gamelogic.NewFixedGameLogicBase(table, tfManager, cfgs)
	scfg := cfgs.GetConfig("gameserver")
	glc.ddt = ddzTimeOutInit()
	tmp := scfg.GetConfigInt("gameserver", "maxplayer")
	glc.maxUserCount = int32(tmp)
	glc.ext = new(userExtInfo)
	glc.oddsCfg = NewDdzOddsCfg(1, cfgs)
	glc.exchangeCfg = NewDdzExchangeCardWeight(cfgs)
	glc.tuoGuanUser = make(map[int32]*robotUserInfoLocal, 0)
	glc.ext.userOpts = make([]*userDdzOpt, glc.maxUserCount)
	glc.ext.turn = 0
	glc.ext.difen = int32(scfg.GetConfigInt("gameserver", "basescore"))
	glc.ext.readytime = int32(scfg.GetConfigInt("gameserver", "readywaittime"))
	glc.ext.timeOut = int32(scfg.GetConfigInt("gameserver", "endactiontimeout"))
	glc.ext.opttime = int32(scfg.GetConfigInt("gameserver", "opttime"))
	glc.ext.maxodd = int32(scfg.GetConfigInt("gameserver", "maxodd"))

	glc.gamePoker = poker.NewGamePoker() //创建游戏牌型
	glc.gamestep = gameProgressNone
	glc.ext.baseScore = int32(scfg.GetConfigInt("gameserver", "basescore"))
	glc.ext.lastQiangDiZhu = -1
	glc.RegisterAfterGameStartFunc(func() {
		glc.ext.init(table)
		//glc.broadcastDeskCfg()
		//glc.oddsChange(0, "base", netproto.OddsChangeType_null)
		//glc.broadcastGameStep(gameProgressBegin) //设置游戏阶段
		//glc.broadcastGameStep(gameProgressFapai)
		//glc.faPai(true, 17) //发牌
	})
	glc.RegisterAfterGameEndFunc(func() {
		glc.clearGameData()
	})

	glc.poker = fixPoker

	for _, ur := range glc.AllGameUser() {
		upk := poker.NewUserPoker(ur, fixPoker.GetUserCards(ur.UserID))
		glc.gamePoker.AddUserPoker(upk)
	}

	//游戏开始后调用函数
	for _, regf := range glc.AfterGameStartFunc {
		regf()
	}

	return glc
}

func TestStartGameLogicServer(t *testing.T) {
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}
	fixtures.StartGameLogicByTest()
}

func TestGameRound(t *testing.T) {
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 8
	//conf.TestMode = true
	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewDDZGameLogic(t, tfManager, cfgs)
	}
	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewDDZRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	fixPoker := fixedPoker("ai2Poker") //讀取發牌資料 /testfixpoker/xx.dat
	table := enterTable(fixPoker.UserCardLen())
	glm, room := fixtures.StartGameLogicByTest()
	glc := glm.GetGameLogic(0).(*DDZGameLogic)

	glc.GameLog = gamelogic.NewGameLogHelper(gamelog.NewEmptyGameLog(), glc.FixedGameLogicBase.LogicBase)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
		ur := glc.GetUserByUserID(userid)
		if ur.IsRobot() {
			if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && glc.gamestep == gameProgressJiaoDiZhu {
				data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)

				u := glc.ext.getExtInfoByTurn(glc.ext.turn)
				if userid != u.userSeat {
					mlog.Error("userId[%v] 非法操作，不该您出牌 u.isTuoGuan[%v] seatno %v ", userid,
						u.isTuoGuan, u.userSeat)
					return
				}

				if *data.SeatNo == ur.SeatNo {
					btmsg := new(msg.Message)
					btmsg.BClassID = int32(netproto.MessageBClassID_Game)
					btmsg.SClassID = int32(netproto.DDZ_GameMessageClassID_JiaoDiZhu)
					btdata := new(netproto.DDZ_JiaoDiZhu)
					btdata.IsJiaoDiZhu = proto.Bool(true)
					btmsg.MsgData, _ = proto.Marshal(btdata)
					glc.GameMessage(ur, btmsg)
				}

				if m.SClassID == int32(netproto.DDZ_GameMessageClassID_NotifyGameTimeOut) && glc.gamestep == gameProgressStart {
					if m.MsgData == nil {
						mlog.Error("data is nil")
						return
					}
					data := m.MsgData.(*netproto.DDZ_NotifyGameTimeOut)
					if *data.SeatNo == ur.SeatNo {

					}
				}
			}

		}
	}
	glc.OnGameStart(func() {
		glc.ExchangeMaxCard(9000001)
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}
	glc.RegisterAfterGameStartFunc(func() {
		glc.ext.init(table)
		//glc.broadcastDeskCfg()
		//glc.oddsChange(0, "base", netproto.OddsChangeType_null)
		//glc.broadcastGameStep(gameProgressBegin) //设置游戏阶段
		//glc.broadcastGameStep(gameProgressFapai)
		//glc.faPai(true, 17) //发牌
	})
	glc.poker = fixedPoker("ai2Poker")

	users := fixtures.GetGameUser()
	//非機器人無法從free 到 sit, 故先統一狀態為QUEUE
	users[0].UserState.OnUserQueue(users[0], room)
	users[1].UserState.OnUserQueue(users[1], room)
	users[2].UserState.OnUserQueue(users[2], room)

	//用戶入座
	glc.UserEnter(users[0], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[1], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)
	glc.UserEnter(users[2], glc.GetTable().TableNo, glc.GetTable().GetFreeSeat().SeatNo)

	glc.UserReady(users[0])
	glc.UserReady(users[1])
	glc.UserReady(users[2])

	//time.Sleep(50* time.Second)
	//glc.GameStartFunc()
}

func confInit() {
	conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 400 //40
	conf.TestMode = true
	mylog.InitDatabaseQueue("", 0, 0, 0, mlog.NewDefaultLogger())
}

//流程整合測試共用初始化
func IntegrationTestInit(endChannel chan bool, t *testing.T) (*fixtures.TestGameRoomServer, *DDZGameLogic, *DDZRobotLogic, chan bool, []*gamelogic.GameUser) {
	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*DDZGameLogic)
	rbglc := tgrs.RobotService.GetRobotLogicManager().GetRobotLogic(0).(*DDZRobotLogic)

	glc.RegisterAfterGameEndFunc(func() {
		endChannel <- true
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.SendFishNotifyFunc = func(m *msg.Message) {

	}

	glc.OnGameStart(func() {
	})

	_, robots := fixtures.GetGameUsers()
	robotManager := tgrs.RobotService.GetRobotManager()
	if robotManager.RobotCount() != 0 {
		t.Errorf("初始化, 機器人數量應為0, RobotCount=>%d", robotManager.RobotCount())
	}
	for i := 0; i < 5; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 10
		robot.UserInfo.UserType = gamelogic.UserTypeRebot
		robotManager.AddRobot(robot)
	}

	glc.TableManager.Room.Seq = proto.Int64(0)

	return tgrs, glc, rbglc, endChannel, robots
}
