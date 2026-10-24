package bjl

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
)

// 1: 莊
// 2: 閒
// 3: 和
// 6: 莊對
// 7: 閒對

// 押莊: 10.00
// 押和: 5.00
// 押莊對: 1.00
// 押閒對: 1.00
// [開莊]
// 有效注: 10.00[押莊] + 5.00[押和] + 1.00[押莊對] + 1.00[押閒對] = 17.00
// 贏: 10.00[押莊][贏] = 10.00
// 稅: 10.00[總贏] * 0.05[稅率] = 0.50
// 返還: 10.00[押莊][退] + 10.00[總贏] - 0.50[稅] = 19.50
// 流水: 5.00[押和][輸] + 1.00[押莊對][輸] + 1.00[押閒對][輸] = 7.00
// 更新金幣: 19.50 - 17.00 = 2.50
func TestGame_OneRound_Case1(t *testing.T) {
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "1", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}

		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				if len(startBet) == 0 {
					startBet <- true
				}
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // User1 Bet
		{ // 下 莊家 10.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(1)
			betdata.Count = proto.Int32(10)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 和 5.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(3)
			betdata.Count = proto.Int32(5)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 莊對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(6)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 閒對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(7)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 250
		expBetAmount = 1700
		expPayAmount = 1950
		expFlowAmount = 700
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

// 押莊: 10.00
// 押和: 5.00
// 押莊對: 1.00
// 押閒對: 1.00
// [開閒]
// 有效注: 10.00[押莊] + 5.00[押和] + 1.00[押莊對] + 1.00[押閒對] = 17.00
// 贏: 0.00
// 返還: 0.00
// 稅: 0.00
// 流水: 10.00[押莊][輸] + 5.00[押和][輸] + 1.00[押莊對][輸] + 1.00[押閒對][輸] = 17.00
// 更新金幣: 0.00 - 17.00 = -17.00
func TestGame_OneRound_Case2(t *testing.T) {
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "2", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				startBet <- true
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // User1 Bet
		{ // 下 莊家 10.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(1)
			betdata.Count = proto.Int32(10)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 和 5.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(3)
			betdata.Count = proto.Int32(5)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 莊對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(6)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 閒對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(7)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -1700
		expBetAmount = 1700
		expPayAmount = 0
		expFlowAmount = 1700
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

		game.UserLeave(u1)

	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 押莊: 10.00
// 押和: 5.00
// 押莊對: 1.00
// 押閒對: 1.00
// [開合]
// 有效注: 10.00[押莊] + 5.00[押和] + 1.00[押莊對] + 1.00[押閒對] = 17.00
// 贏: 5.00[押和][贏] * 8 = 40.00
// 稅: 40.00[總贏] * 0.05[稅率] = 2.00
// 返還: 10.00[押莊][退] + 5.00[押和][退] + 40.00[總贏] - 2.00[稅] = 53.00
// 流水: 1.00[押莊對][輸] + 1.00[押閒對][輸] = 2.00
// 更新金幣: 53.00 - 17.00 = 36.00
func TestGame_OneRound_Case3(t *testing.T) {
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "3", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				startBet <- true
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // u1 bet
		{ // 下 莊家 10.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(1)
			betdata.Count = proto.Int32(10)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 和 5.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(3)
			betdata.Count = proto.Int32(5)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 莊對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(6)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 閒對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(7)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 3600
		expBetAmount = 1700
		expPayAmount = 5300
		expFlowAmount = 200
		expTaxAmount = 200

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

// 押莊: 10.00
// 押和: 5.00
// 押莊對: 1.00
// 押閒對: 1.00
// [開閒 莊對]
// 有效注: 10.00[押莊] + 5.00[押和] + 1.00[押莊對] + 1.00[押閒對] = 17.00
// 贏: 1.00[押莊對][贏] * 11 = 11.00
// 稅: 11.00[總贏] * 0.05[稅率] = 0.55
// 返還: 1.00[押莊對][退] + 11.00[總贏] - 0.55[稅] = 11.45
// 流水: 10.00[押莊][輸] + 5.00[押和][輸] + 1.00[押閒對][輸] = 16.00
// 更新金幣: 11.45 - 17.00 = -5.55
func TestGame_OneRound_Case4(t *testing.T) {
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "2,6", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				startBet <- true
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // u1 bet
		{ // 下 莊家 10.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(1)
			betdata.Count = proto.Int32(10)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 和 5.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(3)
			betdata.Count = proto.Int32(5)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 莊對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(6)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 閒對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(7)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -555
		expBetAmount = 1700
		expPayAmount = 1145
		expFlowAmount = 1600
		expTaxAmount = 55

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

// 押莊: 10.00
// 押和: 5.00
// 押莊對: 1.00
// 押閒對: 1.00
// [開閒 莊對 閒對]
// 有效注: 10.00[押莊] + 5.00[押和] + 1.00[押莊對] + 1.00[押閒對] = 17.00
// 贏: 1.00[押莊對][贏] * 11 + 1.00[押閒對][贏] * 11 = 22.00
// 稅: 22.00[總贏] * 0.05[稅率] = 1.10
// 返還: 1.00[押莊對][退] + 1.00[押閒對][退] + 22.00[總贏] - 1.10[稅率] = 22.90
// 流水: 10.00[押莊][輸] + 5.00[押和][輸] = 15.00
// 更新金幣: 22.90 - 17.00 = 5.90
func TestGame_OneRound_Case5(t *testing.T) {

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "2,6,7", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}

		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				startBet <- true
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // u1 bet
		{ // 下 莊家 10.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(1)
			betdata.Count = proto.Int32(10)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 和 5.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(3)
			betdata.Count = proto.Int32(5)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 莊對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(6)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}

		{ // 下 閒對 1.00
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
			btdata := new(netproto.BJL_BetArr)
			betdata := new(netproto.BJL_Bet)
			betdata.BetType = proto.Int32(7)
			betdata.Count = proto.Int32(1)
			btdata.BetArr = append(btdata.BetArr, betdata)
			m.MsgData, _ = proto.Marshal(btdata)

			game.GameMessage(u1, m)
		}
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 590
		expBetAmount = 1700
		expPayAmount = 2290
		expFlowAmount = 1500
		expTaxAmount = 110

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

// 押閒: 1.00
// [開閒]
// 有效注: 1.00[押閒] = 1.00
// 贏: 1.00[押閒][贏] = 1.00
// 稅: 1.00[總贏] * 0.05[稅率] = 0.05
// 返還: 1.00[押閒][退] + 1.00[總贏] - 0.05[稅] = 1.95
// 流水: 0
// 更新金幣: 1.95 - 1.00 = 0.95
func TestGame_OneRound_Case6(t *testing.T) {
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}
	game.setNextJiangType(0, "2", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	done, startBet := make(chan bool, 1), make(chan bool, 1)

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})
	game.RegisterAfterGameEndFunc(func() {
		done <- true
	})

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			data := m.MsgData.(*netproto.BJL_Sate)
			if data.GetState() == Bet {
				if len(startBet) == 0 {
					startBet <- true
				}
			}
		}
	}

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	<-startBet

	{ // 下 閒家 1.00
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.BJL_GameMessageClassID_BJLBetID)
		btdata := new(netproto.BJL_BetArr)
		betdata := new(netproto.BJL_Bet)
		betdata.BetType = proto.Int32(XIAN)
		betdata.Count = proto.Int32(1)
		btdata.BetArr = append(btdata.BetArr, betdata)
		m.MsgData, _ = proto.Marshal(btdata)

		game.GameMessage(u1, m)
	}

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 95
		expBetAmount = 100
		expPayAmount = 195
		expFlowAmount = 0
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
