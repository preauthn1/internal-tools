package jlbpk

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

func TestJLBPKLogic_SameUserEnter(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, _ := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	u2, _ := fixtures.NewGameUser(1001, "阿花", 10000000)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	u3, _ := fixtures.NewGameUser(1000, "天兵", 10000000)
	u3.UserState.OnUserQueue(u3, game.GetTable().Room)
	game.UserEnter(u3, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	if game.iUserCount != 2 {
		t.Errorf("User Count Want %d, but %d", 2, game.iUserCount)
	}

	if game.arypUserData[0].User != nil {
		t.Errorf("arypUserData[0].User Want nil, but %v", game.arypUserData[0])
	}

	if game.arypUserData[1].User != u2 {
		t.Errorf("arypUserData[1].User Want %v, but %v", u2, game.arypUserData[1].User)
	}

	if game.arypUserData[2].User != u3 {
		t.Errorf("arypUserData[2].User Want %v, but %v", u3, game.arypUserData[2].User)
	}
}

func TestJLBPKLogic_BetSideOverLimit(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted int64

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

		switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
		case netproto.JLBPK_GameMessageClassID_BetResult:
			data := m.MsgData.(*netproto.JLBPK_BetResult)
			if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

				{ // 下 玩家邊注 10.00
					var bet int64 = 1000
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

					data := new(netproto.JLBPK_Bet)
					data.BetType = proto.Int32(JLB_BetType_Side)
					data.Bet = proto.Int64(bet)

					m.MsgData, _ = proto.Marshal(data)

					beted += bet

					game.GameMessage(u1, m)
				}
			}

			if data.GetBetType() == int32(JLB_BetType_Side) && data.GetBet() == 1000 {

				if data.GetResult() != JLB_BetResult_BaseBetOverUpperLimit {
					t.Errorf("User Bet Side Result want %d, but %d", JLB_BetResult_BaseBetOverUpperLimit, data.GetResult())
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
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

func TestJLBPKLogic_BetZhuangSideOverLimit(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	g_iZhuangBetSide = ZhuangBetSideEnable

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted int64

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

		switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
		case netproto.JLBPK_GameMessageClassID_BetResult:
			data := m.MsgData.(*netproto.JLBPK_BetResult)
			if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

				{ // 下 玩家邊注 10.00
					var bet int64 = 1000
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

					data := new(netproto.JLBPK_Bet)
					data.BetType = proto.Int32(JLB_BetType_ZhuangSide)
					data.Bet = proto.Int64(bet)

					m.MsgData, _ = proto.Marshal(data)

					beted += bet

					game.GameMessage(u1, m)
				}
			}

			if data.GetBetType() == int32(JLB_BetType_ZhuangSide) && data.GetBet() == 1000 {

				if data.GetResult() != JLB_BetResult_BaseBetOverUpperLimit {
					t.Errorf("User Bet ZhuangSide Result want %d, but %d", JLB_BetResult_BaseBetOverUpperLimit, data.GetResult())
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
	case <-time.After(5 * time.Second):
		t.Errorf("Timeout")
	}
}

// 自動底注: 1.00
// 手動底注: 10.00
// 階段一: 葫蘆
// 階段二: 棄牌
// 總:
//   有效注: 11.00[底] = 11.00
//   贏: 0.00
//   稅: 0.00
//   返還: 0.00
//   流水: 11.00
//   更新金幣: 0.00 - 11.00 = -11.00
func TestJLBPKLogic_PlayerWin_FullHouse_Drop(t *testing.T) {

	pokerCfg := cfg.GetConfigSection("poker", "Poker")
	pokerCfg["JPTestEnable"] = "2"    // 彩金測試，0:不測試，自然發牌（預設值），1:隨機發彩金牌型，2:指定玩家手牌
	pokerCfg["Dealer"] = "s5c2s2s1s3" // 莊家手牌
	pokerCfg["User1"] = "h5s5s6h6d6"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted int64

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

		switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
		case netproto.JLBPK_GameMessageClassID_BetResult:
			data := m.MsgData.(*netproto.JLBPK_BetResult)
			if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

				{ // 下 玩家邊注 10.00
					var bet int64 = 1000
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

					data := new(netproto.JLBPK_Bet)
					data.BetType = proto.Int32(JLB_BetType_Base)
					data.Bet = proto.Int64(bet)

					m.MsgData, _ = proto.Marshal(data)

					beted += bet

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
// 手動底注: 10.00
// 階段一: 葫蘆
// 階段二: 加注 贏 9倍
// 總:
//   有效注: 11.00[底] + 22.00[加倍] = 33.00
//   贏: 11.00[底][贏] + 22.00[加倍][贏] * 9[特殊牌] = 11.00 + 198.00 = 209.00
//   稅: 209.00[總贏] * 0.05[稅率] = 10.45
//   返還: 11.00[底][退] + 22.00[加倍][退] + 209.00[總贏] - 10.45[稅] = 231.55
//   流水: 0.00
//   更新金幣: 231.55 - 33.00 = 198.55
func TestJLBPKLogic_PlayerWin_FullHouse_Double(t *testing.T) {

	pokerCfg := cfg.GetConfigSection("poker", "Poker")
	pokerCfg["JPTestEnable"] = "2"    // 彩金測試，0:不測試，自然發牌（預設值），1:隨機發彩金牌型，2:指定玩家手牌
	pokerCfg["Dealer"] = "s5c2s2s1s3" // 莊家手牌
	pokerCfg["User1"] = "h5s5s6h6d6"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	c1, done := make(chan bool), make(chan bool)
	go func() {
		done <- <-c1
	}()

	var beted int64

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

		switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
		case netproto.JLBPK_GameMessageClassID_BetResult:
			data := m.MsgData.(*netproto.JLBPK_BetResult)
			if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

				{ // 下 玩家邊注 10.00
					var bet int64 = 1000
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

					data := new(netproto.JLBPK_Bet)
					data.BetType = proto.Int32(JLB_BetType_Base)
					data.Bet = proto.Int64(bet)

					m.MsgData, _ = proto.Marshal(data)

					beted += bet

					game.GameMessage(u1, m)
				}
			}
		case netproto.JLBPK_GameMessageClassID_DoubleBetStart:

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.JLBPK_GameMessageClassID_DoubleBet)

			game.GameMessage(u1, m)
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

		expUpdateScore = 19855
		expBetAmount = 3300
		expPayAmount = 23155
		expFlowAmount = 0
		expTaxAmount = 1045

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

func TestJLBPKLogic_TwoRound_PlayerWin_FullHouse_Double(t *testing.T) {

	pokerCfg := cfg.GetConfigSection("poker", "Poker")
	pokerCfg["JPTestEnable"] = "2"    // 彩金測試，0:不測試，自然發牌（預設值），1:隨機發彩金牌型，2:指定玩家手牌
	pokerCfg["Dealer"] = "s5c2s2s1s3" // 莊家手牌
	pokerCfg["User1"] = "h5s5s6h6d6"  // User1~N 為玩家手牌

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	u1, p1 := fixtures.NewGameUser(1000, "UnitTest", 10000000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	grs.PlayerManager.SetPlayer(p1)
	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	t.Run("Round 1", func(tt *testing.T) {
		c1, done := make(chan bool), make(chan bool)
		go func() {
			done <- <-c1
		}()

		var beted int64

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

			switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
			case netproto.JLBPK_GameMessageClassID_BetResult:
				data := m.MsgData.(*netproto.JLBPK_BetResult)
				if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

					{ // 下 玩家邊注 10.00
						var bet int64 = 1000
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

						data := new(netproto.JLBPK_Bet)
						data.BetType = proto.Int32(JLB_BetType_Base)
						data.Bet = proto.Int64(bet)

						m.MsgData, _ = proto.Marshal(data)

						beted += bet

						game.GameMessage(u1, m)
					}
				}
			case netproto.JLBPK_GameMessageClassID_DoubleBetStart:

				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.JLBPK_GameMessageClassID_DoubleBet)

				game.GameMessage(u1, m)
			}
		}

		select {
		case <-done:
			var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

			expUpdateScore = 19855
			expBetAmount = 3300
			expPayAmount = 23155
			expFlowAmount = 0
			expTaxAmount = 1045

			if u1.UpdateScore != expUpdateScore {
				tt.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
			}

			if u1.BetAmount != expBetAmount {
				tt.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
			}

			if u1.PayAmount != expPayAmount {
				tt.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
			}

			if u1.FlowAmount != expFlowAmount {
				tt.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
			}

			if u1.TaxAmount != expTaxAmount {
				tt.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
			}

			ud := game.GameLog.Context().Users[u1.UserID]

			if ud.UpdateScore != expUpdateScore {
				tt.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
			}

			if ud.BetScore != expBetAmount {
				tt.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount, ud.BetScore)
			}

			if ud.PayoutScore != expPayAmount+expTaxAmount {
				tt.Errorf("User PerGameLog Context PayoutScore want %d, but %d", expPayAmount+expTaxAmount, ud.PayoutScore)
			}

			if ud.ValidScore != expFlowAmount {
				tt.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount, ud.ValidScore)
			}

			if ud.TaxAmount != expTaxAmount {
				tt.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount, ud.TaxAmount)
			}
		case <-time.After(5 * time.Second):
			tt.Errorf("Timeout")
		}
	})

	t.Run("Round 2", func(tt *testing.T) {
		c1, done := make(chan bool), make(chan bool)
		go func() {
			done <- <-c1
		}()

		var beted int64

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

			switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
			case netproto.JLBPK_GameMessageClassID_BetResult:
				data := m.MsgData.(*netproto.JLBPK_BetResult)
				if data.GetResult() == int32(JLB_BetResult_Success) && data.GetBetType() == int32(JLB_BetType_Base) && data.GetBet() == 100 {

					{ // 下 玩家邊注 10.00
						var bet int64 = 1000
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.JLBPK_GameMessageClassID_Bet)

						data := new(netproto.JLBPK_Bet)
						data.BetType = proto.Int32(JLB_BetType_Base)
						data.Bet = proto.Int64(bet)

						m.MsgData, _ = proto.Marshal(data)

						beted += bet

						game.GameMessage(u1, m)
					}
				}
			case netproto.JLBPK_GameMessageClassID_DoubleBetStart:

				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.JLBPK_GameMessageClassID_DoubleBet)

				game.GameMessage(u1, m)
			}
		}

		select {
		case <-done:
			var expUpdateScore, expBetAmount, expPayAmount, expFlowAmount, expTaxAmount int64

			expUpdateScore = 19855 * 2
			expBetAmount = 3300 * 2
			expPayAmount = 23155 * 2
			expFlowAmount = 0 * 2
			expTaxAmount = 1045 * 2

			if u1.UpdateScore != expUpdateScore {
				tt.Errorf("User UpdateScore want %d, but %d", expUpdateScore, u1.UpdateScore)
			}

			if u1.BetAmount != expBetAmount {
				tt.Errorf("User BetAmount want %d, but %d", expBetAmount, u1.BetAmount)
			}

			if u1.PayAmount != expPayAmount {
				tt.Errorf("User PayAmount want %d, but %d", expPayAmount, u1.PayAmount)
			}

			if u1.FlowAmount != expFlowAmount {
				tt.Errorf("User FlowAmount want %d, but %d", expFlowAmount, u1.FlowAmount)
			}

			if u1.TaxAmount != expTaxAmount {
				tt.Errorf("User TaxAmount want %d, but %d", expTaxAmount, u1.TaxAmount)
			}

			ud := game.GameLog.Context().Users[u1.UserID]

			if ud.UpdateScore != expUpdateScore/2 {
				tt.Errorf("User PerGameLog Context UpdateScore want %d, but %d", expUpdateScore/2, u1.UpdateScore)
			}

			if ud.BetScore != expBetAmount/2 {
				tt.Errorf("User PerGameLog Context BetScore want %d, but %d", expBetAmount/2, ud.BetScore)
			}

			if ud.PayoutScore != (expPayAmount+expTaxAmount)/2 {
				tt.Errorf("User PerGameLog Context PayoutScore want %d, but %d", (expPayAmount+expTaxAmount)/2, ud.PayoutScore)
			}

			if ud.ValidScore != expFlowAmount/2 {
				tt.Errorf("User PerGameLog Context ValidScore want %d, but %d", expFlowAmount/2, ud.ValidScore)
			}

			if ud.TaxAmount != expTaxAmount/2 {
				tt.Errorf("User PerGameLog Context TaxAmount want %d, but %d", expTaxAmount/2, ud.TaxAmount)
			}
		case <-time.After(5 * time.Second):
			tt.Errorf("Timeout")
		}
	})
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
					c1 <- true
				}()
			}
		case u2.UserID:
			switch netproto.GameMessageClassID(m.SClassID) {
			case netproto.GameMessageClassID_GameEndID:
				go func() {
					c2 <- true
				}()
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
