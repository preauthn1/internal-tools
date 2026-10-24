package fgsb

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

type TestRoundLimit struct {
	CurrentRound int
	RoundLimit   int
	EndSignal    chan bool
}

// TestGameRound1 測試遊戲流程, 十局
func TestGameRound1(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	conf.TestMode = true
	conf.TimerFastMode = 100
	roundLimit := &TestRoundLimit{
		1,
		10,
		make(chan bool),
	}

	tgrs := fixtures.NewTestGameRoomServer()
	glc := tgrs.GameLogicManager.GetGameLogic(0).(*Logic)

	glc.RegisterAfterGameStartFunc(func() {
		mlog.System("game start ...")
	})

	glc.UpdateUserScoreFunc = func(userid int32, score int64, taxamount int64) {
	}

	glc.RegisterAfterGameEndFunc(func() {
		mlog.System("game end ...")
		roundLimit.EndSignal <- true
	})
	robotManager := tgrs.RobotService.GetRobotManager()
	_, robots := fixtures.GetGameUsers()

	for i := 0; i < 2; i++ {
		robot := new(robotcomm.Robot)
		robot.UserInfo = robots[i]
		robot.RebotLevel = 6
		robot.UserInfo.UserType = 2
		robotManager.AddRobot(robot)
	}

	glc.UserEnter(robots[0], 0, 0)
	glc.UserEnter(robots[1], 0, 1)

End:
	for {
		select {
		case <-roundLimit.EndSignal:
			if roundLimit.CurrentRound >= roundLimit.RoundLimit {
				glc.GameStartFunc = func() {

				}
				break End
			}
			roundLimit.CurrentRound++
		}

	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 6 [亂]
// 玩家 1, 5, 6 [亂]
// 棄牌
// 總:
//   有效注: 10.00[底] + 1.00[邊] = 11.00
//   贏: 0.00
//   稅: 0.00
//   返還: 0.00
//   流水: 11.00
//   更新金幣: 0.00 - 11.00 = -11.00
func TestGame_Shit_Drop(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,6"    // 莊家手牌
	pokerCfg["1000"] = "2,1|3,5|4,6" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_FoldCard)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -1100
		expBetAmount = 1100
		expPayAmount = 0
		expFlowAmount = 1100
		expTaxAmount = 0

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 6 [亂]
// 玩家 1, 5, 6 [亂]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[加注] = 10.00
//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
//   返還: 10.00[底][返] + 10.00[加注][返] + 10.00[加注][贏] - 0.50[稅] = 29.50
//   流水: 1.00[邊] = 1.00
//   更新金幣: 29.50 - 21.00 = 8.50
func TestGame_BankerPlayerShit_PlayerBigger_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,6"    // 莊家手牌
	pokerCfg["1000"] = "2,1|3,5|4,6" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[加注] = 10.00
		//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
		//   返還: 10.00[底][返] + 10.00[加注][返] + 10.00[加注][贏] - 0.50[稅] = 29.50
		//   流水: 1.00[邊] = 1.00
		//   更新金幣: 29.50 - 21.00 = 8.50

		expUpdateScore = 850
		expBetAmount = 2100
		expPayAmount = 2950
		expFlowAmount = 100
		expTaxAmount = 50

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 8, 6 [亂]
// 玩家 1, 5, 6 [亂]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[加注] = 10.00
//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
//   返還: 10.00[底][返] + 10.00[加注][返] + 10.00[加注][贏] - 0.50[稅] = 29.50
//   流水: 1.00[邊] = 1.00
//   更新金幣: 29.50 - 21.00 = 8.50
func TestGame_BankerPlayerShit_BankerBigger_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,8|2,6"    // 莊家手牌
	pokerCfg["1000"] = "2,1|3,5|4,6" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[加注] = 10.00
		//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
		//   返還: 10.00[底][返] + 10.00[加注][返] + 10.00[加注][贏] - 0.50[稅] = 29.50
		//   流水: 1.00[邊] = 1.00
		//   更新金幣: 29.50 - 21.00 = 8.50

		expUpdateScore = 850
		expBetAmount = 2100
		expPayAmount = 2950
		expFlowAmount = 100
		expTaxAmount = 50

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 4 [對]
// 玩家 1, 5, 5 [對]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[底] + 10.00[加注] + 1.00[邊] * 1 = 21.00
//   稅: 21.00[總贏] * 0.05[稅率] = 1.05
//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 21.00[贏] - 1.05[稅] = 40.95
//   流水: 0.00
//   更新金幣: 40.95 - 21.00 = 19.95
func TestGame_Pair_Double_Win(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,4"    // 莊家手牌
	pokerCfg["1000"] = "2,1|3,5|4,5" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[底] + 10.00[加注] + 1.00[邊] * 1 = 21.00
		//   稅: 21.00[總贏] * 0.05[稅率] = 1.05
		//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 21.00[贏] - 1.05[稅] = 40.95
		//   流水: 0.00
		//   更新金幣: 40.95 - 21.00 = 19.95

		expUpdateScore = 1995
		expBetAmount = 2100
		expPayAmount = 4095
		expFlowAmount = 0
		expTaxAmount = 105

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 7, 7 [對]
// 玩家 1, 5, 5 [對]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 1.00[邊] * 1 = 1.00
//   稅: 1.00[總贏] * 0.05[稅率] = 0.05
//   返還: 1.00[邊][返]+ 1.00[贏] - 0.05[稅] = 1.95
//   流水: 10.00[底] + 10.00[加注] = 20.00
//   更新金幣: 1.95 - 21.00 = -19.05
func TestGame_Pair_Double_Lose(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,7|2,7"    // 莊家手牌
	pokerCfg["1000"] = "2,1|3,5|4,5" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 1.00[邊] * 1 = 1.00
		//   稅: 1.00[總贏] * 0.05[稅率] = 0.05
		//   返還: 1.00[邊][返]+ 1.00[贏] - 0.05[稅] = 1.95
		//   流水: 10.00[底] + 10.00[加注] = 20.00
		//   更新金幣: 1.95 - 21.00 = -19.05

		expUpdateScore = -1905
		expBetAmount = 2100
		expPayAmount = 195
		expFlowAmount = 2000
		expTaxAmount = 5

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 4 [對]
// 玩家 1, 2, 6 [同]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[底] + 10.00[加注] + 1.00[邊] * 4 = 24.00
//   稅: 24.00[總贏] * 0.05[稅率] = 1.20
//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 24.00[贏] - 1.20[稅] = 43.80
//   流水: 0.00
//   更新金幣: 43.80 - 21.00 = 22.80
func TestGame_Flush_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,4"    // 莊家手牌
	pokerCfg["1000"] = "2,1|2,2|2,6" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[底] + 10.00[加注] + 1.00[邊] * 4 = 24.00
		//   稅: 24.00[總贏] * 0.05[稅率] = 1.20
		//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 24.00[贏] - 1.20[稅] = 43.80
		//   流水: 0.00
		//   更新金幣: 43.80 - 21.00 = 22.80

		expUpdateScore = 2280
		expBetAmount = 2100
		expPayAmount = 4380
		expFlowAmount = 0
		expTaxAmount = 120

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 4 [對]
// 玩家 2, 3, 4 [順]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[底] + 10.00[底][特] * 1 + 10.00[加注] + 1.00[邊] * 6 = 36.00
//   稅: 36.00[總贏] * 0.05[稅率] = 1.80
//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 36.00[贏] - 1.80[稅] = 55.20
//   流水: 0.00
//   更新金幣: 55.20 - 21.00 = 34.20
func TestGame_Straight_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,4"    // 莊家手牌
	pokerCfg["1000"] = "2,2|3,3|4,4" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[底] + 10.00[底][特] * 1 + 10.00[加注] + 1.00[邊] * 6 = 36.00
		//   稅: 36.00[總贏] * 0.05[稅率] = 1.80
		//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 36.00[贏] - 1.80[稅] = 55.20
		//   流水: 0.00
		//   更新金幣: 55.20 - 21.00 = 34.20

		expUpdateScore = 3420
		expBetAmount = 2100
		expPayAmount = 5520
		expFlowAmount = 0
		expTaxAmount = 180

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 4 [對]
// 玩家 7, 7, 7 [三]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[底] + 10.00[底][特] * 4 + 10.00[加注] + 1.00[邊] * 30 = 90.00
//   稅: 90.00[總贏] * 0.05[稅率] = 4.50
//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 90.00[贏] - 4.50[稅] = 106.50
//   流水: 0.00
//   更新金幣: 106.50 - 21.00 = 85.50
func TestGame_ThreeOfKind_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,4"    // 莊家手牌
	pokerCfg["1000"] = "2,7|3,7|4,7" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[底] + 10.00[底][特] * 4 + 10.00[加注] + 1.00[邊] * 30 = 90.00
		//   稅: 90.00[總贏] * 0.05[稅率] = 4.50
		//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 90.00[贏] - 4.50[稅] = 106.50
		//   流水: 0.00
		//   更新金幣: 106.50 - 21.00 = 85.50

		expUpdateScore = 8550
		expBetAmount = 2100
		expPayAmount = 10650
		expFlowAmount = 0
		expTaxAmount = 450

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 手動底注: 10.00
// 手動邊注: 1.00
// 莊 1, 4, 4 [對]
// 玩家 5, 6, 7 [同花順]
// 加注
// 總:
//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
//   贏: 10.00[底] + 10.00[底][特] * 5 + 10.00[加注] + 1.00[邊] * 40 = 110.00
//   稅: 110.00[總贏] * 0.05[稅率] = 5.50
//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 110.00[贏] - 5.50[稅] = 125.50
//   流水: 0.00
//   更新金幣: 125.50 - 21.00 = 104.50
func TestGame_StraightFlush_Double(t *testing.T) {

	gameCfg := cfg.GetConfigSection("gameserver", "gameserver")
	gameCfg["pokertype"] = "2"

	pokerCfg := cfg.GetConfigSection("poker", "userpoker")
	pokerCfg["0"] = "1,1|1,4|2,4"    // 莊家手牌
	pokerCfg["1000"] = "4,5|4,6|4,7" // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	game.UpdateUserScoreFunc = func(uid int32, score int64, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u1, score, tax)
	}

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var sited bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.GameMessageClassID(m.SClassID) {
		case netproto.GameMessageClassID_GameEndID:
			switch userID {
			case u1.UserID:
				c1 <- true
			}
		}

		switch netproto.FGSB_GameMessageClassID(m.SClassID) {
		case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo:
			if sited {
				return
			}
			sited = true

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)

			data := new(netproto.DCHJK_Msg_SelectSeat)
			data.SeatNo = proto.Int32(1)

			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u1, m)
		case netproto.FGSB_GameMessageClassID_FGSB_UserBettingStart:
			{ // 底注 10.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeBase)
				data.Bet = proto.Int64(1000)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}

			{ // 邊注 1.00
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_Bet)

				data := new(netproto.FGSB_Msg_Bet)
				data.BetType = proto.Int32(BetTypeSide)
				data.Bet = proto.Int64(100)

				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
			}
		case netproto.FGSB_GameMessageClassID_FGSB_DoubleBetStart:
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_DoubleBet)

			game.GameMessage(u1, m)
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	ctx := game.GameLog.Context()

	select {
	case <-done:
		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		//   有效注: 10.00[底] + 1.00[邊] + 10.00[加注] = 21.00
		//   贏: 10.00[底] + 10.00[底][特] * 5 + 10.00[加注] + 1.00[邊] * 40 = 110.00
		//   稅: 110.00[總贏] * 0.05[稅率] = 5.50
		//   返還: 10.00[底][返] + 1.00[邊][返] + 10.00[加注][返] + 110.00[贏] - 5.50[稅] = 125.50
		//   流水: 0.00
		//   更新金幣: 125.50 - 21.00 = 104.50

		expUpdateScore = 10450
		expBetAmount = 2100
		expPayAmount = 12550
		expFlowAmount = 0
		expTaxAmount = 550

		if u1.UpdateScore != expUpdateScore {
			t.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
		}

		if u1.BetAmount != expBetAmount {
			t.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
		}

		if u1.PayAmount != expPayAmount {
			t.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
		}

		if u1.FlowAmount != expFlowAmount {
			t.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
		}

		if u1.TaxAmount != expTaxAmount {
			t.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
		}

		ud := ctx.Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, ud.UpdateScore)
		}

		if ud.BetScore != expBetAmount {
			t.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
		}

		if ud.PayoutScore != expPayAmount+expTaxAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
		}

		if ud.ValidScore != expFlowAmount {
			t.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
		}

		if ud.TaxAmount != expTaxAmount {
			t.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
		}
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestGame_OneRound_TwoUser_OneEnterSlower(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	// conf.TimerFastMode = 10

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
	go func() {
		done <- <-c1 && <-c2
	}()

	var u1Sit, u2Sit bool

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}

		switch userID {
		case u1.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					c1 <- true
				}()
			}
			switch netproto.FGSB_GameMessageClassID(m.SClassID) {
			case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo: // 所有玩家座位資訊
				if !u1Sit {
					u1Sit = true
					_m := new(msg.Message)
					_m.BClassID = int32(netproto.MessageBClassID_Game)
					_m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)
					data := new(netproto.DCHJK_Msg_SelectSeat)
					data.SeatNo = proto.Int32(1)
					_m.MsgData, _ = proto.Marshal(data)
					game.GameMessage(u1, _m)
				}

			}
		case u2.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					c2 <- true
				}()
			}
			switch netproto.FGSB_GameMessageClassID(m.SClassID) {
			case netproto.FGSB_GameMessageClassID_FGSB_AllUserInfo: // 所有玩家座位資訊
				if !u2Sit {
					u2Sit = true
					_m := new(msg.Message)
					_m.BClassID = int32(netproto.MessageBClassID_Game)
					_m.SClassID = int32(netproto.FGSB_GameMessageClassID_FGSB_SelectSeat)
					data := new(netproto.DCHJK_Msg_SelectSeat)
					data.SeatNo = proto.Int32(2)
					_m.MsgData, _ = proto.Marshal(data)
					game.GameMessage(u2, _m)
				}
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	go func() { // User2 Enter
		u2.UserState.OnUserQueue(u2, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p2)
		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUserCnt int = 2

		if len(game.GameLog.Context().Users) != expUserCnt {
			t.Errorf("PerGameLog Context Users Length want %d, but %d", expUserCnt, len(game.GameLog.Context().Users))
		}

	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}
