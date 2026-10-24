package ddz

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
)

//機器人跟牌 錯誤驗證 FollowCards 對子 連對 飛機

//拆王
func TestCombinationKing(t *testing.T) {
	lgc := new(DDZRobotLogic)

	var p poker.Poker
	p = poker.NewTestFilePoker("combinationKing")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)

		userId := int32(i)
		up := fp.GetUserCards(userId)
		if up == nil {
			mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		}

		mapCard := mapCardNumberCount(up)
		mapCard = lgc.combinationKing(mapCard, result) //王

		if i == 0 {
			if result.Ct[0].nMax != poker.PokerNumberBigKing {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_BigKing 錯誤 cards : %v", i, up.Cards)
			}
			if result.Ct[0].nMin != poker.PokerNumberSmallKing {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_SmallKing 錯誤 cards : %v", i, up.Cards)
			}
		}

		if i == 1 {
			if result.Ct[0].nMax != poker.PokerNumberSmallKing {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_BigKing 錯誤 cards : %v", i, up.Cards)
			}
			if result.Ct[0].nMin != 0 {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_SmallKing 錯誤 cards : %v", i, up.Cards)
			}
		}

		if i == 2 {
			if result.Ct[0].nMax != poker.PokerNumberBigKing {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_BigKing 錯誤 cards : %v", i, up.Cards)
			}
			if result.Ct[0].nMin != 0 {
				t.Errorf("index:%v , TestCombinationKing PokerNumber_SmallKing 錯誤 cards : %v", i, up.Cards)
			}
		}

	}
}

//拆二
func TestCombinationTwo(t *testing.T) {
	lgc := new(DDZRobotLogic)

	var p poker.Poker
	p = poker.NewTestFilePoker("combinationTwo")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)

		userId := int32(i)
		up := fp.GetUserCards(userId)
		if up == nil {
			mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		}

		mapCard := mapCardNumberCount(up)
		mapCard = lgc.combinationTwo(mapCard, result) //2

		if i == 0 && result.Ct[0].byteCardType != cardTypeSingle {
			t.Errorf("index:%v combinationTwo 錯誤 cards:%v", i, up.Cards)
		}

		if i == 1 && result.Ct[0].byteCardType != cardTypeDouble {
			t.Errorf("index:%v combinationTwo 錯誤 cards:%v", i, up.Cards)
		}

		if i == 2 && result.Ct[0].byteCardType != cardTypeThree {
			t.Errorf("index:%v combinationTwo 錯誤 cards:%v", i, up.Cards)
		}

		if i == 3 && result.Ct[0].byteCardType != cardTypeBomb {
			t.Errorf("index:%v combinationTwo 錯誤 cards:%v", i, up.Cards)
		}

	}
}

//提取对子和三张
func TestGetPairAndThree(t *testing.T) {
	lgc := new(DDZRobotLogic)

	var p poker.Poker
	p = poker.NewTestFilePoker("getPairAndThree")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)

		userId := int32(i)
		up := fp.GetUserCards(userId)
		if up == nil {
			mlog.Error("系统错误robotOutCards，获取玩家牌型失败")
		}

		mapCard := mapCardNumberCount(up)
		mapCard = lgc.combinationKing(mapCard, result)                               //王
		mapCard = lgc.combinationTwo(mapCard, result)                                //2
		mapCard = lgc.combinationCardsIndex(mapCard, cardTypeBomb, result, int32(i)) //剩下单牌，
		pair, three := lgc.getPairAndThree(result)

		if pair[0] != int32(5) && three[0] != int32(6) {
			t.Errorf("index:%v , GetPairAndThree 提取对子和三张 錯誤 cards : %v", i, up.Cards)
		}

	}
}

//單牌
func TestRobotSingle(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotSingle")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDanpai(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(2) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ok, cg = ddz.FollowCardsDanpai(robot[0])
	ckType = hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(5) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[0])
	ok, cg = ddz.FollowCardsDanpai(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(15) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//單牌
func TestRobotSingle2(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotSingle2")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDanpai(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(16) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ok, cg = ddz.FollowCardsDanpai(robot[0])
	ckType = hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(5) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[0])
	ok, cg = ddz.FollowCardsDanpai(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeSingle)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(15) {
		t.Errorf("機器人出單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//對子
func TestRobotFollowCardsDuizi(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotDuizi")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDuizi(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(5) {
		t.Errorf("機器人出對子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsDuizi(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(6) {
		t.Errorf("機器人出對子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[2])
	ok, cg = ddz.FollowCardsDuizi(robot[2])
	ckType = hisPushAndDeleteCardRobot(cg, robot[2], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[2], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(2) {
		t.Errorf("機器人出對子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//順子
func TestRobotFollowCardsShunzi(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotShunzi")
	ddz.outInfo.nMin, ddz.outInfo.nMax = 4, 8

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsShunzi(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeStraight)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeStraight || cg.Cards[4].CardNumber != int32(9) {
		t.Errorf("機器人出順子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsShunzi(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeStraight)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeStraight || cg.Cards[4].CardNumber != int32(9) {
		t.Errorf("機器人出順子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//連對
func TestRobotFollowCardsLiandui(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotLiandui")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsLiandui(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeThreeStraight)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThreeStraight || cg.Cards[5].CardNumber != int32(5) {
		t.Errorf("機器人出連對 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsLiandui(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeThreeStraight)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeThreeStraight || cg.Cards[5].CardNumber != int32(8) {
		t.Errorf("機器人出連對 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//飛機帶雙
func TestRobotFollowCardsFeiji(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotFeiji")
	ddz.outInfo.nMin, ddz.outInfo.nMax = 3, 6

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsFeiji(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeDoubleStraightd)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeDoubleStraights || cg.Cards[7].CardNumber != int32(8) {
		t.Errorf("機器人出飛機帶雙 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//飛機 3+3+2

//火箭
func TestRobotOutRocket(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotRocket")

	up := ddz.getGameLogic().gamePoker.GetUserPoker(robot[0].UserInfo.UserID)
	ddz.updateRoleInfo(robot[0])
	ddz.outInfo.RealSeatNo = 0
	cgArr := ddz.outRocket(robot[0].UserInfo.UserID)
	cg := ddz.list2cardGroup(up.Cards, cgArr)
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeDoubleStraightd)

	if ckType != cardTypeRocket {
		t.Errorf("機器人出火箭 錯誤 FollowCardsDanpai cg[%v]", cg)
	}
}

//炸彈
func TestFollowCardsBomb(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotBomb")
	ddz.outInfo.RealSeatNo = 0
	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsBomb(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeBomb)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeBomb || cg.Cards[3].CardNumber != int32(5) {
		t.Errorf("機器人出炸彈 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsBomb(robot[1])
	//ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeBomb)
	//ddz.ChupaiRecord(robot[0], cg, ckType)
	//if !ok || ckType != cardTypeBomb || cg.Cards[3].CardNumber != int32(6) {
	//	t.Errorf("機器人出炸彈 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	//}

}

//三條
func TestRobotFollowCardsThree(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotThree")
	ddz.outInfo.RealSeatNo = 0
	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsThree(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeThree)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThree || cg.Cards[2].CardNumber != int32(5) {
		t.Errorf("機器人出三條 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsThree(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeThree)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThree || cg.Cards[2].CardNumber != int32(8) {
		t.Errorf("機器人出三條 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//三带一
func TestRobotFollowCardsThree2(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotThree2")
	ddz.outInfo.RealSeatNo = 0
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsThree(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeThreeAndOne)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThreeAndOne || cg.Cards[3].CardNumber != int32(9) {
		t.Errorf("機器人出三带一 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsThree(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeThreeAndOne)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThreeAndOne || cg.Cards[2].CardNumber != int32(8) {
		t.Errorf("機器人出三带一 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//三带二 葫蘆
func TestRobotFollowCardsThree3(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotThree3")
	ddz.outInfo.RealSeatNo = 0
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))
	ddz.outInfo.cg.Cards = append(ddz.outInfo.cg.Cards, new(poker.Card))

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsThree(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeThreeAndOne)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeThreeAndOne || cg.Cards[4].CardNumber != int32(6) {
		t.Errorf("機器人出葫蘆 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsThree(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeThreeAndOne)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeThreeAndOne || cg.Cards[0].CardNumber != int32(10) {
		t.Errorf("機器人出葫蘆 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//拆对子为单牌
func TestRobotFollowCardsDanpaiUse2(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotDanpaiUse2")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDanpaiUse2(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(5) {
		t.Errorf("機器人出拆对子为单牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsDanpaiUse2(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(2) {
		t.Errorf("機器人出拆对子为单牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//拆3为单牌
func TestRobotFollowCardsDanpaiUse3(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotDanpaiUse3")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDanpaiUse3(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(3) {
		t.Errorf("機器人出拆3为单牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsDanpaiUse3(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(6) {
		t.Errorf("機器人出拆3为单牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[2])
	ok, cg = ddz.FollowCardsDanpaiUse3(robot[2])
	ckType = hisPushAndDeleteCardRobot(cg, robot[2], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[2], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(2) {
		t.Errorf("機器人出拆3为单牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

}

//拆火箭成單牌
func TestRobotOutRocketSplitOne(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotRocketSplitOne")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.outRocketSplitOne(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeSingle)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeSingle || cg.Cards[0].CardNumber != int32(15) {
		t.Errorf("機器人出拆火箭成單牌 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//拆连对配对子
func TestRobotFollowCardsDuiziLian(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotDuiziLian")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDuiziLian(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(3) {
		t.Errorf("機器人出拆连对配对子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

//拆三个的配对子
func TestRobotFollowCardsDuiziUse3(t *testing.T) {
	ddz, robot := NewDDZRobotLogicByTest("robotDuiziUse3")

	ddz.updateRoleInfo(robot[0])
	ok, cg := ddz.FollowCardsDuiziUse3(robot[0])
	ckType := hisPushAndDeleteCardRobot(cg, robot[0], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[0], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(5) {
		t.Errorf("機器人出拆连对配对子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}

	ddz.updateRoleInfo(robot[1])
	ok, cg = ddz.FollowCardsDuiziUse3(robot[1])
	ckType = hisPushAndDeleteCardRobot(cg, robot[1], ddz, cardTypeDouble)
	ddz.ChupaiRecord(robot[1], cg, ckType)
	if !ok || ckType != cardTypeDouble || cg.Cards[0].CardNumber != int32(2) {
		t.Errorf("機器人出拆连对配对子 錯誤 FollowCardsDanpai ok[%v] cg[%v]", ok, cg)
	}
}

func (ddz *DDZRobotLogic) ChupaiRecord(robot *robotcomm.Robot, cg *poker.CardGroup, ct cardType) {
	ddz.outInfo.RealSeatNo = 1
	ddz.outInfo.outType = 1
	ddz.outInfo.cg = new(poker.CardGroup)
	for _, cd := range cg.Cards {
		cdtmp := new(poker.Card)
		cdtmp.Color = cd.Color
		cdtmp.CardNumber = cd.CardNumber
		ddz.outInfo.cg.AddCard(cdtmp)
	}
	ddz.outInfo.nMin, ddz.outInfo.nMax = ddz.analysisCards(ddz.outInfo.cg, int32(ct))
	ddz.outInfo.SeatNo = robot.UserInfo.SeatNo
	ddz.outInfo.CardType = int32(ct)
}

func hisPushAndDeleteCardRobot(cg *poker.CardGroup, robot *robotcomm.Robot, ddz *DDZRobotLogic, ct cardType) cardType {
	action := new(userAction)
	action.ct = ct
	action.cg = cg
	action.userSeat = robot.UserInfo.SeatNo
	ckType := ddz.getGameLogic().ext.his.cardTypeCheck(cg, action)
	ddz.getGameLogic().ext.his.Push(action)
	up := ddz.getGameLogic().gamePoker.GetUserPoker(ddz.getGameLogic().ext.getNowSeat())
	up.Cards.DeleteGroup(cg)
	return ckType
}

//創建測試使用的GameLogic和建立撲克牌發牌器資料
func NewDDZRobotLogicByTest(filePoker string) (*DDZRobotLogic, []*robotcomm.Robot) {
	fixPoker := fixedPoker(filePoker) //讀取發牌資料 /testfixpoker/xx.dat
	table := enterTable(fixPoker.UserCardLen())
	tfManager := new(event.TimerFuncManager)

	currentPath, _ := os.Getwd()
	configDir := path.Join(currentPath, "config")
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)

	robotCfg := robotcomm.ParseRobotConfig(cfgs.GetConfig(gamelogic.ConfigFileNameGameServer))

	rblgc := new(DDZRobotLogic)

	glc := new(DDZGameLogic)
	glc.FixedGameLogicBase = gamelogic.NewFixedGameLogicBase(table, tfManager, cfgs)
	scfg := cfgs.GetConfig("gameserver")
	glc.ddt = ddzTimeOutInit()
	tmp := scfg.GetConfigInt("gameserver", "maxplayer")
	glc.maxUserCount = int32(tmp)
	glc.ext = new(userExtInfo)
	glc.oddsCfg = NewDdzOddsCfg(1, cfgs)
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

	robotArr := make([]*robotcomm.Robot, 0)
	for i, ur := range glc.AllGameUser() {
		upk := poker.NewUserPoker(ur, fixPoker.GetUserCards(ur.UserID))
		glc.gamePoker.AddUserPoker(upk)

		nur := gamelogic.NewGameUser()
		nur.TableNo = int32(i)
		nur.UserID = int32(i)
		nur.SeatNo = int32(i)
		robotitem := new(robotcomm.Robot)
		robotitem.UserInfo = nur
		robotitem.RebotLevel = int32(5)

		robotArr = append(robotArr, robotitem)
	}

	//游戏开始后调用函数
	for _, regf := range glc.AfterGameStartFunc {
		regf()
	}

	rblgc.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, cfgs, glc, nil)
	rblgc.parseRobotCha(cfgs)
	rblgc.robotUser = make(map[int32]*robotUserInfoLocal, 0)
	rblgc.sandBoxUser = make(map[int32]*robotUserInfoLocal, 0)
	rblgc.initOutInfo()
	rblgc.gameChas = make(map[int32]*RobotCha, 0)

	for i, _ := range glc.AllGameUser() {
		urPoker := rblgc.getGameLogic().gamePoker.GetUserPoker(robotArr[i].UserInfo.UserID)
		rblgc.recombination(robotArr[i].UserInfo.UserID, urPoker.Cards)
	}

	user := rblgc.getGameLogic().ext.getExtInfoByTurn(int32(1))
	user.isDiZhu = true

	return rblgc, robotArr
}

func TestA(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	r := rand.Int31n(10)
	fmt.Println(r)
}
