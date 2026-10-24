package wzp

import (
	"fmt"
	"testing"

	"example.invalid/magiconair/properties/assert"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

//是否一对
func TestIsOnePair(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isOnePair")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok := isOnePair(cards)

		if !ok {
			t.Errorf("是否一对 錯誤")
		}
	}
}

//是否两对
func TestIsTwoPair(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isTwoPair")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isTwoPair(cards)

		if !ok {
			t.Errorf("是否两对 錯誤")
		}
	}
}

//是否葫芦
func TestIsHulu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isHulu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isHulu(cards)

		if !ok {
			t.Errorf("是否葫芦 錯誤")
		}
	}
}

//是否三条
func TestIsSanTiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isSanTiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok := isSanTiao(cards)

		if !ok {
			t.Errorf("是否三条 錯誤")
		}
	}
}

//是否四条
func TestIsSitiao(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isSitiao")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok := isSitiao(cards)

		if !ok {
			t.Errorf("是否四条 錯誤")
		}
	}
}

//是否顺子
func TestIsShunzi(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isShunzi")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok := isShunzi(cards)

		if !ok {
			t.Errorf("是否顺子 錯誤")
		}
	}
}

//是否同花
func TestIsTongHua(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isTongHua")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isTongHua(cards)

		if !ok {
			t.Errorf("是否同花 錯誤")
		}
	}
}

//是否同花顺
func TestIsTonghuaShun(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isTonghuaShun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isTonghuaShun(cards)

		if !ok {
			t.Errorf("是否同花顺 錯誤")
		}
	}
}

//是否同花大顺
func TestIsTonghuaDashun(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isTonghuaDashun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isTonghuaDashun(cards)

		if !ok {
			t.Errorf("是否同花大顺 錯誤")
		}
	}
}

func TestGameBegin(t *testing.T) {

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewWZPGameLogic(t, tfManager, cfgs)
	}

	roomserver.RobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) robotlogic.RobotLogic {
		return NewWZPRobotLogic(table, tfManager, robotCfg, configs, glc, evtChannel)
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*WZPGameLogic)

	glc.SendMessgeFunc = func(userid int32, m *msg.Message) {
	}
	glc.OnGameStart(func() {
	})
	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}
	glc.GameEndFunc = func() {
	}

	glc.GameStartFunc()

	users := fixtures.GetGameUser()
	allTableUser := make([]*gamelogic.GameUser, 0)
	allTableUser = append(allTableUser, users[0])
	allTableUser = append(allTableUser, users[1])

	for i := 0; i < 10000; i++ {
		glc.gamePoker = poker.NewGamePoker() //创建游戏牌型

		glc.poker = poker.NewRandomPoker(false, false, 1)
		glc.poker = poker.NewFixedPoker(allTableUser, glc.poker.(*poker.RandomPoker), 5)

		glc.poker.Shuffle()

		ur := users[0]
		ur.SeatNo = 0
		//发两张牌
		cards := glc.poker.Deal(ur.UserID, 2)
		upk1 := poker.NewUserPoker(ur, cards)
		glc.gamePoker.AddUserPoker(upk1)

		ur = users[1]
		ur.SeatNo = 1
		//发两张牌
		cards = glc.poker.Deal(ur.UserID, 2)
		upk2 := poker.NewUserPoker(ur, cards)
		glc.gamePoker.AddUserPoker(upk2)

		// for i:= 0 , ur := range glc.AllGameUser() {
		// 	//发两张牌
		// 	cards := glc.poker.Deal(ur.UserID, 2)

		// 	upk := poker.NewUserPoker(ur, cards)

		// 	glc.gamePoker.AddUserPoker(upk)
		// }

		up, cardtype := getGameMaxCard(glc.gamePoker, 1, glc)
		// assert.Equal(t, upk1.Cards[1:], true)

		cards1 := glc.gamePoker.GetUserPoker(users[0].UserID)
		showCard1 := cards1.Cards.Cards[1]
		cardtype1 := cardtype[cards1.User.UserID]
		compare1 := new(poker.CardGroup)
		compare1.AddCard(showCard1)
		max1 := compare1.GetMaxCardNumber()

		cards2 := glc.gamePoker.GetUserPoker(users[1].UserID)
		showCard2 := cards2.Cards.Cards[1]
		cardtype2 := cardtype[cards2.User.UserID]
		compare2 := new(poker.CardGroup)
		compare2.AddCard(showCard2)
		max2 := compare2.GetMaxCardNumber()

		if cardtype1 == cardtype2 && cardtype1 == 10 {
			if max1 == max2 {
				// assert.Equal(t, users[0].UserID, users[1].UserID, fmt.Sprintf("result:%v user1:%v user2:%v", up.User.UserID, cards1.Cards, cards2.Cards))
				// assert.Equal(t, up.User.UserID, users[0].UserID, fmt.Sprintf("%v > %v result:%v user1:%v user2:%v", max1, max2, up.User.UserID, cards1.Cards, cards2.Cards))
			} else {
				if max1 > max2 {
					assert.Equal(t, up.User.UserID, users[0].UserID, fmt.Sprintf("%v > %v result:%v user1:%v user2:%v", max1, max2, up.User.UserID, cards1.Cards, cards2.Cards))
				} else {
					assert.Equal(t, up.User.UserID, users[1].UserID, fmt.Sprintf("%v < %v result:%v user1:%v user2:%v", max1, max2, up.User.UserID, cards1.Cards, cards2.Cards))
				}
			}
		}

	}
	fmt.Println("finish")
}
