package hhdz

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
	logic := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)

	fmt.Println(logic)
}

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [押黑方][押幸運一擊][開黑方][亂]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 10.00[押黑方][贏] = 10.00
//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
//   返還: 10.00[押黑方][退] + 10.00[總贏] - 0.50[稅] = 19.50
//   流水: 5.00[押幸運一擊][輸] = 5.00
//   更新金幣: 19.50 - 15.00 = 4.50
func TestGame_OneRound_BetBlack_BetLucky_BlackWin_None(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HEI
	game.nextzuhetype = DANPAI

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開紅方][亂]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 0.00
//   稅: 0.00
//   返還: 0.00
//   流水: 10.00[押黑方][輸] + 5.00[押幸運一擊][輸] = 15.00
//   更新金幣: 0.00 - 15.00 = -15.00
func TestGame_OneRound_BetBlack_BetLucky_RedWin_None(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HUNG
	game.nextzuhetype = DANPAI

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開黑方][對子]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 10.00[押黑方][贏] + 5.00[押幸運一擊][贏] * 1[對子] = 15.00
//   稅: 15.00[總贏] * 0.05[稅率] = 0.75
//   返還: 10.00[押黑方][退] + 5.00[押幸運一擊][退] + 15.00[總贏] - 0.75[稅] = 29.25
//   流水: 0.00
//   更新金幣: 29.25 - 15.00 = 14.25
func TestGame_OneRound_BetBlack_BetLucky_BlackWin_Pair(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HEI
	game.nextzuhetype = BIGDUIZI

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

			var rtnAmount int64 = 2925

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

		expUpdateScore = 1425
		expBetAmount = 1500
		expPayAmount = 2925
		expFlowAmount = 0
		expTaxAmount = 75

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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開紅方][順子]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 5.00[押幸運一擊][贏] * 2[順子] = 10.00
//   稅: 10.00[總贏] * 0.05[稅率] = 0.50
//   返還: 5.00[押幸運一擊][退] + 10.00[總贏] - 0.50[稅] = 14.50
//   流水: 10.00[押黑方][輸] = 10.00
//   更新金幣: 14.50 - 15.00 = -0.50
func TestGame_OneRound_BetBlack_BetLucky_RedWin_Straight(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HUNG
	game.nextzuhetype = SHUNZI

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

			var rtnAmount int64 = 1450

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

		expUpdateScore = -50
		expBetAmount = 1500
		expPayAmount = 1450
		expFlowAmount = 1000
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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開黑方][金花]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 10.00[押黑方][贏] + 5.00[押幸運一擊][贏] * 3[金花] = 25.00
//   稅: 25.00[總贏] * 0.05[稅率] = 1.25
//   返還: 10.00[押黑方][退] + 5.00[押幸運一擊][退] + 25.00[總贏] - 1.25[稅] = 38.75
//   流水: 0.00
//   更新金幣: 38.75 - 15.00 = 23.75
func TestGame_OneRound_BetBlack_BetLucky_BlackWin_Flush(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HEI
	game.nextzuhetype = JINHUA

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

			var rtnAmount int64 = 3875

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

		expUpdateScore = 2375
		expBetAmount = 1500
		expPayAmount = 3875
		expFlowAmount = 0
		expTaxAmount = 125

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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開紅方][順金]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 5.00[押幸運一擊][贏] * 5[順金] = 25.00
//   稅: 25.00[總贏] * 0.05[稅率] = 1.25
//   返還: 5.00[押幸運一擊][退] + 25.00[總贏] - 1.25[稅] = 28.75
//   流水: 10.00[押黑方][輸] = 10.00
//   更新金幣: 28.75 - 15.00 = 13.75
func TestGame_OneRound_BetBlack_BetLucky_RedWin_StraightFlush(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HUNG
	game.nextzuhetype = SHUNJIN

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

			var rtnAmount int64 = 2875

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

		expUpdateScore = 1375
		expBetAmount = 1500
		expPayAmount = 2875
		expFlowAmount = 1000
		expTaxAmount = 125

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

// 押黑方: 10.00
// 押幸運一擊: 5.00
// [開黑方][豹子]
//   有效注: 10.00[押黑方] + 5.00[押幸運一擊] = 15.00
//   贏: 10.00[押黑方][贏] + 5.00[押幸運一擊][贏] * 10[豹子] = 60.00
//   稅: 60.00[總贏] * 0.05[稅率] = 3.00
//   返還: 10.00[押黑方][退] + 5.00[押幸運一擊][退] + 60.00[總贏] - 3.00[稅] = 72.00
//   流水: 0.00
//   更新金幣: 72.00 - 15.00 = 67.00
func TestGame_OneRound_BetBlack_BetLucky_BlackWin_Three(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)
	game.nextjiangtype = HEI
	game.nextzuhetype = BAOZI

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

		switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
		case netproto.HHDZ_GameMessageClassID_HHDZStateID:
			data := m.MsgData.(*netproto.HHDZ_Sate)
			if data.GetState() == Bet && !beted {
				beted = true

				{ // 下 黑方 10.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(HEI)
					betdata.Count = proto.Int32(1000)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

				{ // 下 幸運一擊 5.00
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
					btdata := new(netproto.HHDZ_BetArr)
					betdata := new(netproto.HHDZ_Bet)
					betdata.BetType = proto.Int32(NATURE)
					betdata.Count = proto.Int32(500)
					btdata.BetArr = append(btdata.BetArr, betdata)
					m.MsgData, _ = proto.Marshal(btdata)

					game.GameMessage(u1, m)
				}

			}
		case netproto.HHDZ_GameMessageClassID_HHDZGameResultWinID:
			data := m.MsgData.(*netproto.HHDZ_GameResultUserMoney)

			var rtnAmount int64 = 7200

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

		expUpdateScore = 5700
		expBetAmount = 1500
		expPayAmount = 7200
		expFlowAmount = 0
		expTaxAmount = 300

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
	game := grs.GameLogicManager.GetGameLogic(0).(*HHDZLogic)

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
					c1 <- true
				}()
			}
			switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
			case netproto.HHDZ_GameMessageClassID_HHDZStateID:
				data := m.MsgData.(*netproto.HHDZ_Sate)
				if data.GetState() == Bet && !u1beted {
					u1beted = true
					{ // 下 黑方 10.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
						btdata := new(netproto.HHDZ_BetArr)
						betdata := new(netproto.HHDZ_Bet)
						betdata.BetType = proto.Int32(HEI)
						betdata.Count = proto.Int32(1000)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u1, m)
					}
					{ // 下 幸運一擊 5.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
						btdata := new(netproto.HHDZ_BetArr)
						betdata := new(netproto.HHDZ_Bet)
						betdata.BetType = proto.Int32(NATURE)
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
					c2 <- true
				}()
			}
			switch netproto.HHDZ_GameMessageClassID(m.SClassID) {
			case netproto.HHDZ_GameMessageClassID_HHDZStateID:
				data := m.MsgData.(*netproto.HHDZ_Sate)
				if data.GetState() == Bet && !u2beted {
					u2beted = true
					{ // 下 黑方 10.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
						btdata := new(netproto.HHDZ_BetArr)
						betdata := new(netproto.HHDZ_Bet)
						betdata.BetType = proto.Int32(HEI)
						betdata.Count = proto.Int32(1000)
						btdata.BetArr = append(btdata.BetArr, betdata)
						m.MsgData, _ = proto.Marshal(btdata)

						game.GameMessage(u1, m)
					}
					{ // 下 幸運一擊 5.00
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.HHDZ_GameMessageClassID_HHDZBetID)
						btdata := new(netproto.HHDZ_BetArr)
						betdata := new(netproto.HHDZ_Bet)
						betdata.BetType = proto.Int32(NATURE)
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
