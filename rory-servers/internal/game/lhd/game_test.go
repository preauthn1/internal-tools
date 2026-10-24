package lhd

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

func TestNewGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)

	fmt.Println(logic)
}

// 0: 虎
// 1: 龍
// 2: 和

// 押龍: 10.00
// 押和: 5.00
// [押龍][押和][開龍]
// 有效注: 10.00[押龍] + 5.00[押和] = 15.00
// 贏: 10.00[押龍][贏] = 10.00
// 稅: 10.00[總贏] * 0.05[稅率] = 0.50
// 返還: 10.00[押龍][退] + 10.00[總贏] - 0.50[稅] = 19.50
// 流水: 5.00[押和][輸] = 5.00
// 更新金幣: 19.50 - 15.00 = 4.50
func TestGame_OneRound_BetDragon_BetTie_DragonWin(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, LONG, "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 龍 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(LONG)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = 1950

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 450
		expBetAmount = 1500
		expPayAmount = 1950
		expFlowAmount = 500
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 押虎: 10.00
// 押和: 5.00
// [押虎][押和][開虎]
// 有效注: 10.00[押虎] + 5.00[押和] = 15.00
// 贏: 10.00[押虎][贏] = 10.00
// 稅: 10.00[總贏] * 0.05[稅率] = 0.50
// 返還: 10.00[押虎][退] + 10.00[總贏] - 0.50[稅] = 19.50
// 流水: 5.00[押和][輸] = 5.00
// 更新金幣: 19.50 - 15.00 = 4.50
func TestGame_OneRound_BetTiger_BetTie_TigerWin(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HU, "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 虎 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HU)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = 1950

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 450
		expBetAmount = 1500
		expPayAmount = 1950
		expFlowAmount = 500
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 押龍: 10.00
// 押和: 5.00
// [押龍][押和][開虎]
// 有效注: 10.00[押龍] + 5.00[押和] = 15.00
// 贏: 0.00
// 返還: 0.00
// 稅: 0.00
// 流水: 10.00[押龍][輸] + 5.00[押和][輸] = 15.00
// 更新金幣: 0.00 - 15.00 = -15.00
func TestGame_OneRound_BetDragon_BetTie_TigerWin(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HU, "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 龍 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(LONG)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = -1500

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -1500
		expBetAmount = 1500
		expPayAmount = 0
		expFlowAmount = 1500
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 押龍: 10.00
// 押和: 5.00
// [押龍][押和][開合][全退]
// 有效注: 10.00[押龍] + 5.00[押和] = 15.00
// 贏: 5.00[押和][贏] * 8 = 40.00
// 稅: 40.00[總贏] * 0.05[稅率] = 2.00
// 返還: 10.00[押龍][退] + 5.00[押和][退] + 40.00[總贏] - 2.00[稅] = 53.00
// 流水: 0.00
// 更新金幣: 53.00 - 15.00 = 38.00
func TestGame_OneRound_BetDragon_BetTie_Tie_Return(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HE, "UnitTest")
	game.HEandKill = HEandBack

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 龍 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(LONG)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = 5300

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 3800
		expBetAmount = 1500
		expPayAmount = 5300
		expFlowAmount = 0
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 押龍: 10.00
// 押和: 5.00
// [押龍][押和][開合][全殺]
// 有效注: 10.00[押龍] + 5.00[押和] = 15.00
// 贏: 5.00[押和][贏] * 8 = 40.00
// 稅: 40.00[總贏] * 0.05[稅率] = 2.00
// 返還: 5.00[押和][退] + 40.00[總贏] - 2.00[稅] = 43.00
// 流水: 10.00[押龍] = 10.00
// 更新金幣: 43.00 - 15.00 = 28.00
func TestGame_OneRound_BetDragon_BetTie_Tie_Kill(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HE, "UnitTest")
	game.HEandKill = HEandKill

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 龍 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(LONG)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = 4300

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 2800
		expBetAmount = 1500
		expPayAmount = 4300
		expFlowAmount = 1000
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 押龍: 10.00
// 押和: 5.00
// [押龍][押和][開合][退一半]
// 有效注: 10.00[押龍] + 5.00[押和] = 15.00
// 贏: 5.00[押和][贏] * 8 = 40.00
// 稅: 40.00[總贏] * 0.05[稅率] = 2.00
// 返還: 10.00[押龍][退] / 2 + 5.00[押和][退] + 40.00[總贏] - 2.00[稅] = 48.00
// 流水: 10.00[押龍] / 2 = 5.00
// 更新金幣: 48.00 - 15.00 = 33.00
func TestGame_OneRound_BetDragon_BetTie_Tie_Helf(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HE, "UnitTest")
	game.HEandKill = HEandKillHalf

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted bool

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

		switch netproto.LHD_GameMessageClassID(m.SClassID) {
		case netproto.LHD_GameMessageClassID_LHDStateID:
			data := m.MsgData.(*netproto.LHD_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 龍 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(LONG)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 和 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
					btdata := new(netproto.LHD_BetArr)
					betdata := new(netproto.LHD_Bet)
					betdata.BetType = proto.Int32(HE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.LHD_GameMessageClassID_LHDGameResultWinID:
			data := m.MsgData.(*netproto.LHD_GameResultUserMoney)

			var rtnAmount int64 = 4800

			if data.GetSelfWin() != rtnAmount {
				t.Errorf("Return Money want %d, but , %d", rtnAmount, data.GetSelfWin())
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 3300
		expBetAmount = 1500
		expPayAmount = 4800
		expFlowAmount = 500
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

		ud := game.GameLog.Context().Users[u1.UserID]

		if ud.UpdateScore != expUpdateScore {
			t.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
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

// 兩個玩家入場，檢查PerGameLog是否為兩個玩家
func TestGame_OneRound_TwoUser_OneEnterSlower(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*LongHudouLogic)
	game.setNextJiangType(0, HE, "UnitTest")
	game.HEandKill = HEandKillHalf

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
	go func() {
		done <- <-c1 && <-c2
	}()

	var u1beted, u2beted bool

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
					fmt.Println("c1 in")
					c1 <- true
					fmt.Println("c1 in done")
				}()
			}
			switch netproto.LHD_GameMessageClassID(m.SClassID) {
			case netproto.LHD_GameMessageClassID_LHDStateID:
				data := m.MsgData.(*netproto.LHD_Sate)
				if data.GetState() == Bet && !u1beted {
					u1beted = true

					{ // 下 龍 10.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
						btdata := new(netproto.LHD_BetArr)
						betdata := new(netproto.LHD_Bet)
						betdata.BetType = proto.Int32(LONG)
						betdata.Count = proto.Int32(1000)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u1, m)
					}

					{ // 下 和 5.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
						btdata := new(netproto.LHD_BetArr)
						betdata := new(netproto.LHD_Bet)
						betdata.BetType = proto.Int32(HE)
						betdata.Count = proto.Int32(500)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u1, m)
					}

				}
			}
		case u2.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					fmt.Println("c2 in")
					c2 <- true
					fmt.Println("c2 in done")
				}()
			}
			switch netproto.LHD_GameMessageClassID(m.SClassID) {
			case netproto.LHD_GameMessageClassID_LHDStateID:
				data := m.MsgData.(*netproto.LHD_Sate)
				if data.GetState() == Bet && !u2beted {
					u2beted = true

					{ // 下 龍 10.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
						btdata := new(netproto.LHD_BetArr)
						betdata := new(netproto.LHD_Bet)
						betdata.BetType = proto.Int32(LONG)
						betdata.Count = proto.Int32(1000)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u2, m)
					}

					{ // 下 和 5.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.LHD_GameMessageClassID_LHDBetID)
						btdata := new(netproto.LHD_BetArr)
						betdata := new(netproto.LHD_Bet)
						betdata.BetType = proto.Int32(HE)
						betdata.Count = proto.Int32(500)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u2, m)
					}

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

		fmt.Println("Done")

		var expUserCnt int = 2

		if len(game.GameLog.Context().Users) != expUserCnt {
			t.Errorf("PerGameLog Context Users Length want %d, but %d", expUserCnt, len(game.GameLog.Context().Users))
		}

	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}
