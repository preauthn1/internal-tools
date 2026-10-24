package infc

import (
	"fmt"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
)

func TestNewFixedTimeGame(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewFixedTimeGame(t, tfManager, cfgs)
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*FixedTimeGame)

	fmt.Printf("%+v\n", logic)
}

// func TestFixedTimeGame_PlayOneRound_Finish(t *testing.T) {

// 	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
// 		return NewFixedTimeGame(t, tfManager, cfgs)
// 	}

// 	conf.CLogLevel = mlog.InfoLevel

// 	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
// 	game := grs.GameLogicManager.GetGameLogic(0).(*FixedTimeGame)

// 	u1 := fixtures.GetGameUser()[0]
// 	u2 := fixtures.GetGameUser()[1]

// 	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
// 	go func() {
// 		done <- <-c1 && <-c2
// 	}()
// 	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
// 		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
// 		if m.BClassID != int32(netproto.MessageBClassID_Game) {
// 			return
// 		}
// 		switch netproto.GameMessageClassID(m.SClassID) {
// 		case netproto.GameMessageClassID_GameEndID:
// 			switch userID {
// 			case u1.UserID:
// 				c1 <- true
// 			case u2.UserID:
// 				c2 <- true
// 			}
// 		}
// 		switch GameMessageID(m.SClassID) {
// 		case GameMessageID_ReadyCheck:
// 			m := new(msg.Message)
// 			m.BClassID = int32(netproto.MessageBClassID_Game)
// 			m.SClassID = int32(GameMessageID_Ready)

// 			switch userID {
// 			case u1.UserID:
// 				game.GameMessage(u1, m)
// 			case u2.UserID:
// 				game.UserLeave(u2)
// 			}
// 		case GameMessageID_Start:
// 			switch userID {
// 			case u1.UserID:

// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1500
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(300 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u1, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()

// 				go func() {
// 					select {
// 					case <-time.After(4 * time.Second):
// 						var miss int64 = 1
// 						data := make(map[string]interface{})
// 						data["id"] = userID
// 						data["miss"] = miss
// 						data["type"] = "add"

// 						m := new(msg.Message)
// 						m.BClassID = int32(netproto.MessageBClassID_Game)
// 						m.SClassID = int32(GameMessageID_Miss)

// 						m.MsgData = data

// 						game.GameMessage(u1, m)
// 					}
// 				}()
// 			case u2.UserID:
// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1000
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(600 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u2, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()
// 			}
// 		}
// 	}

// 	go func() { // User1 Enter
// 		u1.Score = 10000
// 		u1.NickName = "WTF"
// 		u1.UserType = gamelogic.UserType_Normal
// 		u1.UserState.OnUserQueue(u1, game.GetTable().Room)

// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p1 := player.NewGuest(conninfo)
// 		p1.UserID = u1.UserID
// 		p1.UserData = u1
// 		grs.PlayerManager.SetPlayer(p1)

// 		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	go func() { // User2 Enter
// 		time.Sleep(10 * time.Millisecond)

// 		u2.Score = 10000
// 		u2.NickName = "Rory87"
// 		u2.UserType = gamelogic.UserType_Normal
// 		u2.UserState.OnUserQueue(u2, game.GetTable().Room)

// 		// tcpconn := network.NewTcpConn(nil, 2048)
// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p2 := player.NewGuest(conninfo)
// 		p2.UserID = u2.UserID
// 		p2.UserData = u2
// 		grs.PlayerManager.SetPlayer(p2)

// 		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	select {
// 	case <-done:
// 		// for _, u := range game.TableManager.GetTableUser() {
// 		// 	example.invalid("(退房總紀錄) User[%s] Score: %10d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, u.Score, u.BetAmount, u.PayAmount, u.TaxAmount)
// 		// 	example.invalid("(單場紀錄) User[%s] UpdateAmount: %5d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, game.ctx.UpdateAmount[u.UserID], game.ctx.BetAmount[u.UserID], game.ctx.PayAmount[u.UserID], game.ctx.ValidAmount[u.UserID])
// 		// }

// 		if u1.Score != 10450 {
// 			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
// 		}
// 		if u2.Score != 9500 {
// 			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
// 		}

// 		time.Sleep(500 * time.Millisecond)

// 		for _, u := range game.TableManager.GetTableUser() {
// 			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
// 		}

// 	case <-time.After(30 * time.Second):
// 		t.Error("timeout")
// 	}
// }

// // 玩家斷線
// func TestFixedTimeGame_PlayOneRound_OnePlayerOffLine(t *testing.T) {

// 	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
// 		return NewFixedTimeGame(t, tfManager, cfgs)
// 	}

// 	conf.CLogLevel = mlog.InfoLevel

// 	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
// 	game := grs.GameLogicManager.GetGameLogic(0).(*FixedTimeGame)

// 	// watch table users
// 	// go func() {
// 	// 	for {
// 	// 		example.invalid("Watch Table Users")
// 	// 		for _, u := range game.TableManager.GetTableUser() {
// 	// 			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
// 	// 		}
// 	// 		time.Sleep(100 * time.Millisecond)
// 	// 	}
// 	// }()

// 	u1 := fixtures.GetGameUser()[0]
// 	u2 := fixtures.GetGameUser()[1]

// 	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
// 	go func() {
// 		done <- <-c1
// 	}()
// 	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
// 		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
// 		if m.BClassID != int32(netproto.MessageBClassID_Game) {
// 			return
// 		}
// 		switch netproto.GameMessageClassID(m.SClassID) {
// 		case netproto.GameMessageClassID_GameEndID:
// 			switch userID {
// 			case u1.UserID:
// 				c1 <- true
// 			case u2.UserID:
// 				c2 <- true
// 			}
// 		}
// 		switch GameMessageID(m.SClassID) {
// 		case GameMessageID_ReadyCheck:
// 			m := new(msg.Message)
// 			m.BClassID = int32(netproto.MessageBClassID_Game)
// 			m.SClassID = int32(GameMessageID_Ready)

// 			switch userID {
// 			case u1.UserID:
// 				game.GameMessage(u1, m)
// 			case u2.UserID:
// 				// game.GameMessage(u2, m)
// 			}
// 		case GameMessageID_Start:
// 			switch userID {
// 			case u1.UserID:

// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1500
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(300 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u1, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()
// 			case u2.UserID:
// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1000
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(600 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u2, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()

// 				go func() {
// 					select {
// 					case <-time.After(3 * time.Second):
// 						example.invalid("玩家[%d] 離線", userID)
// 						u2.UserState.OnUserOffLine(u2, game.TableManager.Room)
// 						cancel()
// 					}
// 				}()
// 			}
// 		}
// 	}

// 	go func() { // User1 Enter
// 		u1.Score = 10000
// 		u1.NickName = "WTF"
// 		u1.UserType = gamelogic.UserType_Normal
// 		u1.UserState.OnUserQueue(u1, game.GetTable().Room)

// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p1 := player.NewGuest(conninfo)
// 		p1.UserID = u1.UserID
// 		p1.UserData = u1
// 		grs.PlayerManager.SetPlayer(p1)

// 		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	go func() { // User2 Enter
// 		time.Sleep(10 * time.Millisecond)

// 		u2.Score = 10000
// 		u2.NickName = "Rory87"
// 		u2.UserType = gamelogic.UserType_Normal
// 		u2.UserState.OnUserQueue(u2, game.GetTable().Room)

// 		// tcpconn := network.NewTcpConn(nil, 2048)
// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p2 := player.NewGuest(conninfo)
// 		p2.UserID = u2.UserID
// 		p2.UserData = u2
// 		grs.PlayerManager.SetPlayer(p2)

// 		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	select {
// 	case <-done:
// 		// for _, u := range game.TableManager.GetTableUser() {
// 		// 	example.invalid("(退房總紀錄) User[%s] Score: %10d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, u.Score, u.BetAmount, u.PayAmount, u.TaxAmount)
// 		// 	example.invalid("(單場紀錄) User[%s] UpdateAmount: %5d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, game.ctx.UpdateAmount[u.UserID], game.ctx.BetAmount[u.UserID], game.ctx.PayAmount[u.UserID], game.ctx.ValidAmount[u.UserID])
// 		// }

// 		if u1.Score != 10450 {
// 			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
// 		}
// 		if u2.Score != 9500 {
// 			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
// 		}

// 		time.Sleep(500 * time.Millisecond)

// 		if len(game.TableManager.GetTableUser()) != 1 {
// 			t.Errorf("Table User Count should be %d, but %d", 1, len(game.TableManager.GetTableUser()))
// 		}

// 	case <-time.After(30 * time.Second):
// 		t.Error("timeout")
// 	}
// }

// func TestFixedTimeGame_PlayOneRound_OnePlayerRecover(t *testing.T) {

// 	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.ConfigGroup) gamelogic.GameLogic {
// 		return NewFixedTimeGame(t, tfManager, cfgs)
// 	}

// 	conf.CLogLevel = mlog.InfoLevel

// 	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
// 	game := grs.GameLogicManager.GetGameLogic(0).(*FixedTimeGame)

// 	// watch table users
// 	// go func() {
// 	// 	for {
// 	// 		example.invalid("Watch Table Users")
// 	// 		for _, u := range game.TableManager.GetTableUser() {
// 	// 			example.invalid("User[%s] 狀態 %s", u.NickName, u.UserState.StateName())
// 	// 		}
// 	// 		time.Sleep(100 * time.Millisecond)
// 	// 	}
// 	// }()

// 	u1 := fixtures.GetGameUser()[0]
// 	u2 := fixtures.GetGameUser()[1]

// 	c1, c2, done := make(chan bool), make(chan bool), make(chan bool)
// 	go func() {
// 		done <- <-c1 && <-c2
// 	}()
// 	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
// 		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
// 		if m.BClassID != int32(netproto.MessageBClassID_Game) {
// 			return
// 		}
// 		switch netproto.GameMessageClassID(m.SClassID) {
// 		case netproto.GameMessageClassID_GameEndID:
// 			switch userID {
// 			case u1.UserID:
// 				c1 <- true
// 			case u2.UserID:
// 				c2 <- true
// 			}
// 		}
// 		switch GameMessageID(m.SClassID) {
// 		case GameMessageID_ReadyCheck:
// 			m := new(msg.Message)
// 			m.BClassID = int32(netproto.MessageBClassID_Game)
// 			m.SClassID = int32(GameMessageID_Ready)

// 			switch userID {
// 			case u1.UserID:
// 				game.GameMessage(u1, m)
// 			case u2.UserID:
// 				// game.GameMessage(u2, m)
// 			}
// 		case GameMessageID_Start:
// 			switch userID {
// 			case u1.UserID:

// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1500
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(300 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u1, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()
// 			case u2.UserID:
// 				ctx, cancel := context.WithCancel(context.TODO())

// 				go func() {

// 					var total int64 = 1000
// 					var send int64

// 					for {
// 						select {
// 						case <-time.After(600 * time.Millisecond):

// 							left := total - send

// 							s := rand.Int63n(left)

// 							if left < 5 {
// 								s = left
// 							}

// 							data := make(map[string]interface{})
// 							data["id"] = userID
// 							data["score"] = s
// 							data["type"] = "add"

// 							m := new(msg.Message)
// 							m.BClassID = int32(netproto.MessageBClassID_Game)
// 							m.SClassID = int32(GameMessageID_Score)

// 							m.MsgData = data

// 							game.GameMessage(u2, m)

// 							send += s

// 							if send == total {
// 								cancel()
// 							}
// 						case <-ctx.Done():
// 							return
// 						}
// 					}
// 				}()

// 				go func() {
// 					select {
// 					case <-time.After(3 * time.Second):
// 						example.invalid("玩家[%d] 離線", userID)
// 						u2.UserState.OnUserOffLine(u2, game.TableManager.Room)
// 						// cancel()

// 						time.Sleep(1 * time.Second)
// 						example.invalid("玩家[%d] 返回", userID)
// 						// u2.UserState.OnUserGame(u2, game.TableManager.Room)
// 						game.UserRecover(u2)
// 					}
// 				}()
// 			}
// 		}
// 	}

// 	go func() { // User1 Enter
// 		u1.Score = 10000
// 		u1.NickName = "WTF"
// 		u1.UserType = gamelogic.UserType_Normal
// 		u1.UserState.OnUserQueue(u1, game.GetTable().Room)

// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p1 := player.NewGuest(conninfo)
// 		p1.UserID = u1.UserID
// 		p1.UserData = u1
// 		grs.PlayerManager.SetPlayer(p1)

// 		game.UserEnter(u1, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	go func() { // User2 Enter
// 		time.Sleep(10 * time.Millisecond)

// 		u2.Score = 10000
// 		u2.NickName = "Rory87"
// 		u2.UserType = gamelogic.UserType_Normal
// 		u2.UserState.OnUserQueue(u2, game.GetTable().Room)

// 		// tcpconn := network.NewTcpConn(nil, 2048)
// 		conninfo := network.NewConnInfo(nil, fmt.Sprintf("%d", grs.GameServer.ServerID), 0)

// 		p2 := player.NewGuest(conninfo)
// 		p2.UserID = u2.UserID
// 		p2.UserData = u2
// 		grs.PlayerManager.SetPlayer(p2)

// 		game.UserEnter(u2, game.GetTable().TableNo, game.GetTable().GetFreeSeat().SeatNo)
// 	}()

// 	select {
// 	case <-done:
// 		// for _, u := range game.TableManager.GetTableUser() {
// 		// 	example.invalid("(退房總紀錄) User[%s] Score: %10d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, u.Score, u.BetAmount, u.PayAmount, u.TaxAmount)
// 		// 	example.invalid("(單場紀錄) User[%s] UpdateAmount: %5d, BetAmount: %5d, PayAmount: %5d, TaxAmount: %5d", u.NickName, game.ctx.UpdateAmount[u.UserID], game.ctx.BetAmount[u.UserID], game.ctx.PayAmount[u.UserID], game.ctx.ValidAmount[u.UserID])
// 		// }

// 		if u1.Score != 10450 {
// 			t.Errorf("User1 Score should be %d, but %d", 10450, u1.Score)
// 		}
// 		if u2.Score != 9500 {
// 			t.Errorf("User2 Score should be %d, but %d", 9500, u2.Score)
// 		}

// 		time.Sleep(500 * time.Millisecond)

// 		if len(game.TableManager.GetTableUser()) != 2 {
// 			t.Errorf("Table User Count should be %d, but %d", 2, len(game.TableManager.GetTableUser()))
// 		}

// 	case <-time.After(30 * time.Second):
// 		t.Error("timeout")
// 	}
// }
