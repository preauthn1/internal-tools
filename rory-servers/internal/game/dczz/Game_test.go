package dczz

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

func TestNewGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	fmt.Println(logic)
}

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 輸
//   有效注: 101.00[底] + 10.00[和] = 111.00
//   返還: 0.00
//   稅: 0.00
//   流水: 111.00
// 階段二: 無
// 總:
//   有效注: 111.00[階段一]
//   返還: 0.00
//   稅: 0.00
//   流水: 111.00[階段一]
//   更新金幣: 0.00 - 111.00 = -111.00
func TestGame_OneRound_BankerWin(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"        // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s5c2h3s2c2" // 莊家手牌
	pokerCfg["User1"] = "h2c2s3c2s2"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -11100
		expBetAmount = 11100
		expPayAmount = 0
		expFlowAmount = 11100
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

		if ud.PayoutScore != expPayAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount, ud.PayoutScore)
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

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 贏
//   有效注: 101.00[底] + 10.00[和] = 111.00
//   返還: 101.00[底] + 101.00[底贏] - 101.00[總贏] * 0.05[稅率] = 196.95
//   稅:  5.05
//   流水: 0.00
// 階段二: 無
// 總:
//   有效注: 111.00[階段一]
//   返還: 196.95[階段一]
//   稅:  5.05[階段一]
//   流水: 0.00[階段一]
//   更新金幣: 196.95 - 111.00 = 85.95
func TestGame_OneRound_PlayerWin(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"        // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s2c2h3s2c2" // 莊家手牌
	pokerCfg["User1"] = "h5c2s3c2s2"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 8595
		expBetAmount = 11100
		expPayAmount = 19695
		expFlowAmount = 0
		expTaxAmount = 505

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

// 自動底注: 1.00
// 手動底注: 100.00
// 階段一: 和
//   有效注: 0.00
//   返還: 0.00
//   稅:  0.00
//   流水: 0.00
// 階段二: 棄牌(投降輸一半)
//   有效注: 101.00[底] = 101.00
//   返還: 101.00[底] / 2 = 50.50
//   稅:  0.00
//   流水: 50.50
// 總:
//   有效注: 0.00[階段一] + 101.00[階段二] = 101.00
//   返還: 0.00[階段一] + 50.50[階段二] = 50.50
//   稅:  0.00[階段一] + 0.00[階段二] = 0.00
//   流水: 0.00[階段一] + 50.50[階段二] = 50.50
//   更新金幣: 50.50 - 101.00 = -50.50
func TestGame_OneRound_Tie_Fold_NoBet(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"        // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s1c2h3s2c2" // 莊家手牌
	pokerCfg["User1"] = "h1c2s3c2s2"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -5050
		expBetAmount = 10100
		expPayAmount = 5050
		expFlowAmount = 5050
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

		if ud.PayoutScore != expPayAmount {
			t.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount, ud.PayoutScore)
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

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 和
//   有效注: 10.00[和] = 10.00
//   返還: 10.00[和] + 100.00[和贏] - 100.00[總贏] * 0.05[稅率]) = 110.00 - 5.00 = 105.00
//   稅:  5.00
//   流水: 0.00
// 階段二: 棄牌(投降輸一半)
//   有效注: 101.00[底] = 101.00
//   返還: 101.00[底] / 2 = 50.50
//   稅:  0.00
//   流水: 50.50
// 總:
//   有效注: 10.00[階段一] + 101.00[階段二] = 111.00
//   返還: 105.00[階段一] + 50.50[階段二] = 155.50
//   稅:  5.00[階段一] + 0.00[階段二] = 5.00
//   流水: 0.00[階段一] + 50.50[階段二] = 50.50 --> 因為沒輸所以為 `0`
//   更新金幣: 155.50 - 111.00 = 44.50
func TestGame_OneRound_Tie_Fold_Bet(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"        // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s1c2h3s2c2" // 莊家手牌
	pokerCfg["User1"] = "h1c2s3c2s2"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 4450
		expBetAmount = 11100
		expPayAmount = 15550
		expFlowAmount = 0
		expTaxAmount = 500

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

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 和
//   有效注: 10.00[和] = 10.00
//   返還: 10.00[和] + 100.00[和贏] - 100.00[總贏] * 0.05[稅率]) = 105.00
//   稅:  5.00
//   流水: 0.00
// 階段二: 戰 贏
//   有效注: 101.00[底] + 101.00[加倍] = 202.00
//   返還: 101.00[底] + 101.00[加倍] + 101.00[戰贏] - 101.00[總贏] * 0.05[稅率]) = 303.00 - 5.05 = 297.95
//   稅: 5.05
//   流水: 0.00
// 總:
//   有效注: 10.00[階段一] + 202.00[階段二] = 212.00
//   返還: 105.00[階段一] + 297.95[階段二] = 402.95
//   稅: 5.00[階段一] + 5.05[階段二] = 10.05
//   流水: 0.00[階段一] + 0.00[階段二] = 0.00
//   更新金幣: 402.95 - 212.00 = 190.95
func TestGame_OneRound_Tie_War_Win(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"  // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s5c2" // 莊家手牌
	pokerCfg["User1"] = "h5s4"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
			}
		case netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase:
			fmt.Println("WarPhase")
			{ // 戰
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection)
				data := new(netproto.DCZZ_Msg_WarPhaseSelection)
				data.Selection = proto.Int32(int32(WarSelection_War))
				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 19095
		expBetAmount = 21200
		expPayAmount = 40295
		expFlowAmount = 0
		expTaxAmount = 1005

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

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 和
//   有效注: 10.00[和] = 10.00
//   返還: 10.00[和] + 100.00[和贏] - 100.00[總贏] * 0.05[稅率]) = 105.00
//   稅:  5.00
//   流水: 0.00
// 階段二: 戰 輸
//   有效注: 101.00[底] + 101.00[加倍] = 202.00
//   返還: 0.00
//   稅: 0.00
//   流水: 202.00
// 總:
//   有效注: 10.00[階段一] + 202.00[階段二] = 212.00
//   返還: 105.00[階段一] + 0.00[階段二] = 105.00
//   稅: 5.00[階段一] + 0.00[階段二] = 5.00
//   流水: 0.00[階段一] + 202.00[階段二] = 202.00  --> 因為輸的絕對值所以為 `107`
//   更新金幣: 105.00 - 212.00 = -107.00
func TestGame_OneRound_Tie_War_Lose(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"  // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s5c5" // 莊家手牌
	pokerCfg["User1"] = "h5s4"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
			}
		case netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase:
			fmt.Println("WarPhase")
			{ // 戰
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection)
				data := new(netproto.DCZZ_Msg_WarPhaseSelection)
				data.Selection = proto.Int32(int32(WarSelection_War))
				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = -10700
		expBetAmount = 21200
		expPayAmount = 10500
		expFlowAmount = 10700
		expTaxAmount = 500

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

// 自動底注: 1.00
// 手動底注: 100.00
// 手動和牌: 10.00
// 階段一: 和
//   有效注: 10.00[和] = 10.00
//   返還: 10.00[和] + 100.00[和贏] - 100.00[總贏] * 0.05[稅率]) = 105.00
//   稅:  5.00
//   流水: 0.00
// 階段二: 戰 和
//   有效注: 101.00[底] + 101.00[加倍] = 202.00
//   返還: 101.00[底] + 101.00[加倍] + 101.00[戰和] * 2 - 202.00[總贏] * 0.05[稅率]) = 404.00 - 10.10 = 393.90
//   稅: 10.10
//   流水: 0.00
// 總:
//   有效注: 10.00[階段一] + 202.00[階段二] = 212.00
//   返還: 105.00[階段一] + 393.90[階段二] = 498.90
//   稅: 5.00[階段一] + 10.10[階段二] = 15.10
//   流水: 0.00[階段一] + 0.00[階段二] = 0.00
//   更新金幣: 498.90 - 212.00 = 286.90
func TestGame_OneRound_Tie_War_Tie(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	pokerCfg := cfg.GetConfigSection("poker", "poker")
	pokerCfg["DealType"] = "1"  // 發牌方式，0:自然發牌（預設值），1:指定發牌
	pokerCfg["Dealer"] = "s5c2" // 莊家手牌
	pokerCfg["User1"] = "h5s2"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()
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

		switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
		case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
			data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
			if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
				{ // 下底注100
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(0)
					data.Bet = proto.Int64(10000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
				{ // 下平手10
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
					data := new(netproto.DCZZ_Msg_Bet)
					data.BetType = proto.Int32(1)
					data.Bet = proto.Int64(1000)
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
				}
			}
		case netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase:
			fmt.Println("WarPhase")
			{ // 戰
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection)
				data := new(netproto.DCZZ_Msg_WarPhaseSelection)
				data.Selection = proto.Int32(int32(WarSelection_War))
				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u1, m)
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
		game.UserLeave(u1)

		var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

		expUpdateScore = 28690
		expBetAmount = 21200
		expPayAmount = 49890
		expFlowAmount = 0
		expTaxAmount = 1510

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
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
	go func() {
		done <- <-c1 && <-c2
	}()
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
			switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
			case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
				data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
				if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
					{ // 下底注100
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
						data := new(netproto.DCZZ_Msg_Bet)
						data.BetType = proto.Int32(0)
						data.Bet = proto.Int64(10000)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					}
					{ // 下平手10
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
						data := new(netproto.DCZZ_Msg_Bet)
						data.BetType = proto.Int32(1)
						data.Bet = proto.Int64(1000)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					}
				}
			case netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase:
				fmt.Println("WarPhase")
				{ // 戰
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection)
					data := new(netproto.DCZZ_Msg_WarPhaseSelection)
					data.Selection = proto.Int32(int32(WarSelection_War))
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u1, m)
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
			switch netproto.DCZZ_GameMessageClassID(m.SClassID) {
			case netproto.DCZZ_GameMessageClassID_DCZZ_BetResult:
				data := m.MsgData.(*netproto.DCZZ_Msg_BetResult)
				if data.GetResult() == int32(BetResult_Success) && data.GetBetType() == int32(BetType_Base) && data.GetBet() == 100 {
					{ // 下底注100
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
						data := new(netproto.DCZZ_Msg_Bet)
						data.BetType = proto.Int32(0)
						data.Bet = proto.Int64(10000)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					}
					{ // 下平手10
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_Bet)
						data := new(netproto.DCZZ_Msg_Bet)
						data.BetType = proto.Int32(1)
						data.Bet = proto.Int64(1000)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					}
				}
			case netproto.DCZZ_GameMessageClassID_DCZZ_WarPhase:
				fmt.Println("WarPhase")
				{ // 戰
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.DCZZ_GameMessageClassID_DCZZ_WarPhaseSelection)
					data := new(netproto.DCZZ_Msg_WarPhaseSelection)
					data.Selection = proto.Int32(int32(WarSelection_War))
					m.MsgData, _ = proto.Marshal(data)

					game.GameMessage(u2, m)
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
