package infc

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestNewDeathGame(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)

	fmt.Printf("%+v\n", logic)

	for _, row := range logic.Proto.Patterns {
		fmt.Println(row)
	}
	for _, row := range logic.Proto.FruitAttributes {
		fmt.Println(row)
	}
	for _, row := range logic.Proto.FruitEmitters {
		fmt.Println(row)
	}
	for _, row := range logic.Proto.StageGroups {
		fmt.Println(row)
	}
}

func TestDeathGame_PlayOneRound_兩個玩家都中離(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
	go func() {
		done <- <-c1 && <-c2
	}()
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
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
		switch netproto.INFC_GameMessageClassID(m.SClassID) {
		case netproto.INFC_GameMessageClassID_INFC_GameStartID:
			switch userID {
			case u1.UserID:
				go func() {
					time.Sleep(100 * time.Second)
					game.UserLeave(u1)
					c1 <- true
				}()
			case u2.UserID:
				go func() {
					game.UserLeave(u2)
					c2 <- true
				}()
			}
		}
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-done:

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		time.Sleep(10 * time.Millisecond)

		if len(game.TableManager.GetTableUser()) > 0 {
			t.Errorf("Table want clean user, but %d in there", len(game.TableManager.GetTableUser()))
		}

		if u1.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User1 State want %s, but %s", gamelogic.StateNameFree, u1.UserState.StateName())
		}

		if u2.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User2 State want %s, but %s", gamelogic.StateNameFree, u2.UserState.StateName())
		}

	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P2_MISS3(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			ctx, cancel := context.WithCancel(mainctx)
			go func() {
				var total int64 = 1500
				var send int64
				for {
					select {
					case <-time.After(10 * time.Millisecond):

						left := total - send

						s := rand.Int63n(left)

						if left < 5 {
							s = left
						}

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)

						send += s

						if send == total {
							cancel()
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			ctx, cancel := context.WithCancel(mainctx)
			go func() {

				var total int64 = 1000
				var send int64

				for {
					select {
					case <-time.After(50 * time.Millisecond):

						left := total - send

						s := rand.Int63n(left)

						if left < 5 {
							s = left
						}

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)

						send += s

						if send == total {
							cancel()
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Miss
			go func() {
				for {
					select {
					case <-time.After(100 * time.Millisecond):
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_MissID)
						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
	})

	game.RegisterAfterGameEndFunc(func() {
		maincancel()
	})
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		if m.BClassID != int32(netproto.MessageBClassID_Game) {
			return
		}
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		for _, u := range game.TableManager.GetTableUser() {
			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
		}

		if len(game.TableManager.GetTableUser()) > 0 {
			t.Errorf("Table want clean user, but %d in there", len(game.TableManager.GetTableUser()))
		}

		if u1.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User1 State want %s, but %s", gamelogic.StateNameFree, u1.UserState.StateName())
		}

		if u2.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User2 State want %s, but %s", gamelogic.StateNameFree, u2.UserState.StateName())
		}

	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P2_BOOM(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			ctx, cancel := context.WithCancel(mainctx)

			go func() {
				var total int64 = 1500
				var send int64
				for {
					select {
					case <-time.After(10 * time.Millisecond):

						left := total - send

						s := rand.Int63n(left)

						if left < 5 {
							s = left
						}

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)

						send += s

						if send == total {
							cancel()
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			ctx, cancel := context.WithCancel(mainctx)

			go func() {
				var total int64 = 1000
				var send int64
				for {
					select {
					case <-time.After(50 * time.Millisecond):

						left := total - send

						s := rand.Int63n(left)

						if left < 5 {
							s = left
						}

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)

						send += s

						if send == total {
							cancel()
						}
					case <-ctx.Done():
						return
					}
				}
			}()
		}
		{ // User 2 Boom
			go func() {
				for {
					select {
					case <-time.After(100 * time.Millisecond):
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BoomID)
						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		for _, u := range game.TableManager.GetTableUser() {
			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P1_OVERSCORE_P2_BOOM(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			go func() {
				for {
					select {
					case <-time.After(10 * time.Millisecond):

						s := rand.Int63n(100)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(200)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Boom
			go func() {
				select {
				case <-time.After(100 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BoomID)
					game.GameMessage(u2, m)
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		for _, u := range game.TableManager.GetTableUser() {
			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P1_Giveup_After_P2_BOOM(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(5)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(200)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Boom
			go func() {
				select {
				case <-time.After(100 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BoomID)
					game.GameMessage(u2, m)
				}
			}()
		}
		{ // User1 Giveup
			go func() {
				select {
				case <-time.After(150 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GiveupID)
					game.GameMessage(u1, m)
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 9500 {
			t.Errorf("User1 Score should be %d, but %d", 9500, u1.Score)
		}
		if u2.Score != 10450 {
			t.Errorf("User2 Score should be %d, but %d", 10450, u2.Score)
		}

		for _, u := range game.TableManager.GetTableUser() {
			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P2_Giveup_After_P2_BOOM(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(5)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(200)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Boom
			go func() {
				select {
				case <-time.After(100 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_BoomID)
					game.GameMessage(u2, m)
				}
			}()
		}
		{ // User2 Giveup
			go func() {
				select {
				case <-time.After(150 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GiveupID)
					game.GameMessage(u2, m)
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	mCntMap := make(map[int32]int)
	var mu sync.Mutex
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		mu.Lock()
		mCntMap[m.SClassID]++
		mu.Unlock()
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		if mCntMap[10040] != 2 {
			t.Errorf("GameResult message (10001:10040) want %d, but %d", 2, mCntMap[10040])
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}

func TestDeathGame_PlayOneRound_P2_Giveup_After_P2_MISS3(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(5)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(200)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 MISS
			go func() {
				for {
					select {
					case <-time.After(50 * time.Millisecond):
						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_MissID)
						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Giveup
			go func() {
				select {
				case <-time.After(200 * time.Millisecond):
					m := new(msg.Message)
					m.BClassID = int32(netproto.MessageBClassID_Game)
					m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_GiveupID)
					game.GameMessage(u2, m)
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	mCntMap := make(map[int32]int)
	var mu sync.Mutex
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		mu.Lock()
		mCntMap[m.SClassID]++
		mu.Unlock()
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		if mCntMap[10040] != 2 {
			t.Errorf("GameResult message (10001:10040) want %d, but %d", 2, mCntMap[10040])
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}
func TestDeathGame_PlayOneRound_P2_OffLine_Recover(t *testing.T) {

	roomserver.GameLogicCreator = DeathGameCreater

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*DeathGame)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	u1, p1 := fixtures.NewGameUser(1000, "咖哩", 10000)
	u2, p2 := fixtures.NewGameUser(2000, "辣雞", 10000)

	u1.UserState.OnUserQueue(u1, game.GetTable().Room)
	u2.UserState.OnUserQueue(u2, game.GetTable().Room)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	mainctx, maincancel := context.WithCancel(context.Background())

	game.RegisterAfterGameStartFunc(func() {
		{ // User1 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(50)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u1, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Score
			go func() {
				for {
					select {
					case <-time.After(5 * time.Millisecond):

						s := rand.Int63n(200)

						data := new(netproto.INFC_IncScore)
						data.Score = s

						m := new(msg.Message)
						m.BClassID = int32(netproto.MessageBClassID_Game)
						m.SClassID = int32(netproto.INFC_GameMessageClassID_INFC_IncScoreID)
						m.MsgData, _ = proto.Marshal(data)

						game.GameMessage(u2, m)
					case <-mainctx.Done():
						return
					}
				}
			}()
		}
		{ // User2 Leave
			go func() {
				select {
				case <-time.After(100 * time.Millisecond):
					game.UserLeave(u2)
				}
			}()
		}
		{ // User2 Recover
			go func() {
				select {
				case <-time.After(200 * time.Millisecond):
					game.UserRecover(u2)
				}
			}()
		}
	})
	game.RegisterAfterGameEndFunc(func() { maincancel() })

	mCntMap := make(map[int32]int)
	var mu sync.Mutex
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		mu.Lock()
		mCntMap[m.SClassID]++
		mu.Unlock()
	}

	game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)

	select {
	case <-mainctx.Done():

		if u1.Score != 10450 {
			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
		}
		if u2.Score != 9500 {
			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
		}

		if mCntMap[10040] != 2 {
			t.Errorf("GameResult message (10001:10040) want %d, but %d", 2, mCntMap[10040])
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout")
	}
}
