package bjl

import (
	"fmt"
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

func TestNewGame(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)

	fmt.Println(logic)
}

// 測試莊對閒對對押，是否都不會開
func TestGame_OneRound_CheckBetBankPairAndPlayerPair(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 20000

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	// game.setNextJiangType(0, "2,6,7", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)

	var u1cnt, u2cnt int
	var total int = 1

	go func() {
		for {
			select {
			default:
				if u1cnt == total && u2cnt == total {
					done <- true
				}
			case <-c1:
				u1cnt++
			case <-c2:
				u2cnt++
			}
		}
	}()

	var betcnt1, betcnt2 int

	var BankerCnt, PlayerCnt, TieCnt, BankerPairCnt, PlayerPairCnt int

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
			case u2.UserID:
				c2 <- true
			}
		}

		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			switch userID {
			case u1.UserID:
				data := m.MsgData.(*netproto.BJL_Sate)
				if data.GetState() == Bet && betcnt1 <= u1cnt+1 {
					betcnt1++
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
			case u2.UserID:
				data := m.MsgData.(*netproto.BJL_Sate)
				if data.GetState() == Bet && betcnt2 <= u1cnt+1 {
					betcnt2++
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

						game.GameMessage(u2, m)
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

						game.GameMessage(u2, m)
					}
				}
			}
		case netproto.BJL_GameMessageClassID_BJLGameResultPokerID:
			if userID != u1.UserID {
				return
			}
			data := m.MsgData.(*netproto.BJL_GameResultPoker)

			// Banker Win
			if data.GetZuHe()&(1<<(ZHUANG-1)) != 0 {
				BankerCnt++
			}
			// Player Win
			if data.GetZuHe()&(1<<(XIAN-1)) != 0 {
				PlayerCnt++
			}
			// Tie Win
			if data.GetZuHe()&(1<<(HE-1)) != 0 {
				TieCnt++
			}
			// Banker Pair
			if data.GetZuHe()&(1<<(ZHUANGDUI-1)) != 0 {
				BankerPairCnt++
			}
			// Player Pair
			if data.GetZuHe()&(1<<(XIANDUI-1)) != 0 {
				PlayerPairCnt++
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
		fmt.Printf("Banker: %d\n", BankerCnt)
		fmt.Printf("Player: %d\n", PlayerCnt)
		fmt.Printf("Tie: %d\n", TieCnt)
		fmt.Printf("Banker Pair: %d\n", BankerPairCnt)
		fmt.Printf("Player Pair: %d\n", PlayerPairCnt)

	case <-time.After(5 * time.Minute):
		t.Errorf("Timeout")
	}
}

// 兩個玩家入場，檢查PerGameLog是否為兩個玩家
func TestGame_OneRound_TwoUser_OneEnterSlower(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// conf.CLogLevel = mlog.DebugLevel
	conf.TimerFastMode = 1000

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*BaiJialeLogic)
	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	// game.setNextJiangType(0, "2,6,7", "UnitTest")

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 10000000)
	u2, p2 := fixtures.NewGameUser(1001, "阿瓜", 10000000)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)

	var u1cnt, u2cnt int
	var total int = 1

	go func() {
		for {
			select {
			default:
				if u1cnt == total && u2cnt == total {
					done <- true
				}
			case <-c1:
				u1cnt++
			case <-c2:
				u2cnt++
			}
		}
	}()

	var betcnt1, betcnt2 int

	var BankerCnt, PlayerCnt, TieCnt, BankerPairCnt, PlayerPairCnt int

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
			case u2.UserID:
				c2 <- true
			}
		}

		switch netproto.BJL_GameMessageClassID(m.SClassID) {
		case netproto.BJL_GameMessageClassID_BJLStateID:
			switch userID {
			case u1.UserID:
				data := m.MsgData.(*netproto.BJL_Sate)
				if data.GetState() == Bet && betcnt1 <= u1cnt+1 {
					betcnt1++
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
			case u2.UserID:
				data := m.MsgData.(*netproto.BJL_Sate)
				if data.GetState() == Bet && betcnt2 <= u1cnt+1 {
					betcnt2++
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

						game.GameMessage(u2, m)
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

						game.GameMessage(u2, m)
					}
				}
			}
		case netproto.BJL_GameMessageClassID_BJLGameResultPokerID:
			if userID != u1.UserID {
				return
			}
			data := m.MsgData.(*netproto.BJL_GameResultPoker)

			// Banker Win
			if data.GetZuHe()&(1<<(ZHUANG-1)) != 0 {
				BankerCnt++
			}
			// Player Win
			if data.GetZuHe()&(1<<(XIAN-1)) != 0 {
				PlayerCnt++
			}
			// Tie Win
			if data.GetZuHe()&(1<<(HE-1)) != 0 {
				TieCnt++
			}
			// Banker Pair
			if data.GetZuHe()&(1<<(ZHUANGDUI-1)) != 0 {
				BankerPairCnt++
			}
			// Player Pair
			if data.GetZuHe()&(1<<(XIANDUI-1)) != 0 {
				PlayerPairCnt++
			}
		}
	}

	go func() { // User1 Enter
		u1.UserState.OnUserQueue(u1, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p1)
		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	go func() { // User2 Enter
		// time.Sleep(1 * time.Second)
		u2.UserState.OnUserQueue(u2, game.GetTable().Room)
		grs.PlayerManager.SetPlayer(p2)
		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	}()

	select {
	case <-done:

		fmt.Println("Done")
		fmt.Println("PerGameLog Context Users")
		for _, u := range game.GameLog.Context().Users {
			fmt.Printf("Player: %+v\n", u)
		}

		var expUserCnt int = 2

		if len(game.GameLog.Context().Users) != expUserCnt {
			t.Errorf("PerGameLog Context Users Length want %d, but %d", expUserCnt, len(game.GameLog.Context().Users))
		}

	case <-time.After(5 * time.Minute):
		t.Errorf("Timeout")
	}
}
