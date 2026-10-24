package inrm

import (
	"fmt"
	"math/rand"
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
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

func TestGameServer_玩家入桌(t *testing.T) {

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", 100000)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", 100000)
	u4, p4 := fixtures.NewGameUser(4000, "阿機", 100000)
	u5, p5 := fixtures.NewGameUser(5000, "三八", 100000)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)
	grs.PlayerManager.SetPlayer(p4)
	grs.PlayerManager.SetPlayer(p5)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)
	u4.UserState.OnUserQueue(u4, game.TableManager.Room)
	u5.UserState.OnUserQueue(u5, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u4, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u5, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	if u1.UserState.StateName() != gamelogic.StateNameReady {
		t.Errorf("User 1 State want %s, but %s", gamelogic.StateNameReady, u1.UserState.StateName())
	}
	if u2.UserState.StateName() != gamelogic.StateNameReady {
		t.Errorf("User 2 State want %s, but %s", gamelogic.StateNameReady, u2.UserState.StateName())
	}
	if u3.UserState.StateName() != gamelogic.StateNameReady {
		t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameReady, u3.UserState.StateName())
	}
	if u4.UserState.StateName() != gamelogic.StateNameReady {
		t.Errorf("User 4 State want %s, but %s", gamelogic.StateNameReady, u4.UserState.StateName())
	}
	if u5.UserState.StateName() != gamelogic.StateNameReady {
		t.Errorf("User 5 State want %s, but %s", gamelogic.StateNameReady, u5.UserState.StateName())
	}
	if game.State != CountDown {
		t.Errorf("Game State want %v, but %v", CountDown, game.State)
	}
}

func TestGameServer_遊戲啟動(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	u1, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", 100000)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", 100000)
	u4, p4 := fixtures.NewGameUser(4000, "阿機", 100000)
	u5, p5 := fixtures.NewGameUser(5000, "三八", 100000)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)
	grs.PlayerManager.SetPlayer(p4)
	grs.PlayerManager.SetPlayer(p5)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)
	u4.UserState.OnUserQueue(u4, game.TableManager.Room)
	u5.UserState.OnUserQueue(u5, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u4, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u5, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	select {
	case <-start:
		if u1.UserState.StateName() != gamelogic.StateNameGame {
			t.Errorf("User 1 State want %s, but %s", gamelogic.StateNameGame, u1.UserState.StateName())
		}
		if u2.UserState.StateName() != gamelogic.StateNameGame {
			t.Errorf("User 2 State want %s, but %s", gamelogic.StateNameGame, u2.UserState.StateName())
		}
		if u3.UserState.StateName() != gamelogic.StateNameGame {
			t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameGame, u3.UserState.StateName())
		}
		if u4.UserState.StateName() != gamelogic.StateNameGame {
			t.Errorf("User 4 State want %s, but %s", gamelogic.StateNameGame, u4.UserState.StateName())
		}
		if u5.UserState.StateName() != gamelogic.StateNameGame {
			t.Errorf("User 5 State want %s, but %s", gamelogic.StateNameGame, u5.UserState.StateName())
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

// 全部Miss
//  最後一個玩家獲勝，贏 40 * 4 * 50 = 8000分，稅金 8000 * 0.05 = 400
//  每人輸40點 2000分
func TestGameServer_全部Miss(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.ActionWaitTime = 1 * time.Second
	game.StartWaitTime = 30 * time.Second

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", orgScore)
	u4, p4 := fixtures.NewGameUser(4000, "阿機", orgScore)
	u5, p5 := fixtures.NewGameUser(5000, "三八", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)
	grs.PlayerManager.SetPlayer(p4)
	grs.PlayerManager.SetPlayer(p5)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)
	u4.UserState.OnUserQueue(u4, game.TableManager.Room)
	u5.UserState.OnUserQueue(u5, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u4, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u5, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	select {
	case <-end:
		if u1.UserState.StateName() != gamelogic.StateNameSit {
			t.Errorf("User 1 State want %s, but %s", gamelogic.StateNameSit, u1.UserState.StateName())
		}
		if u2.UserState.StateName() != gamelogic.StateNameSit {
			t.Errorf("User 2 State want %s, but %s", gamelogic.StateNameSit, u2.UserState.StateName())
		}
		if u3.UserState.StateName() != gamelogic.StateNameSit {
			t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameSit, u3.UserState.StateName())
		}
		if u4.UserState.StateName() != gamelogic.StateNameSit {
			t.Errorf("User 4 State want %s, but %s", gamelogic.StateNameSit, u4.UserState.StateName())
		}
		if u5.UserState.StateName() != gamelogic.StateNameSit {
			t.Errorf("User 5 State want %s, but %s", gamelogic.StateNameSit, u5.UserState.StateName())
		}

		if game.IsGaming() {
			t.Errorf("Game IsGaming() want %t, but %t", false, game.IsGaming())
		}

		if game.Game != nil {
			t.Errorf("Game innerGame want %v, but %v", nil, game)
		}

		if len(ctx.Users) != 5 {
			t.Errorf("GameLog Context users length want %d, but %d", 5, len(ctx.Users))
		}

		assertWinLoss(t, ctx, u1, orgScore, -2000, 2000, 2000, 0)
		assertWinLoss(t, ctx, u2, orgScore, -2000, 2000, 2000, 0)
		assertWinLoss(t, ctx, u3, orgScore, -2000, 2000, 2000, 0)
		assertWinLoss(t, ctx, u4, orgScore, -2000, 2000, 2000, 0)
		assertWinLoss(t, ctx, u5, orgScore, 8000, 0, 0, 400)
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_兩人_正常胡牌(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 15 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		game.Game.RandJoker = poker.Jack
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	{ // P1 更新手牌 [♦2 ♣2 ♠3 ♥3 ♣3 ♥5 ♦5 ♠7 ♣9 ♦10 ♦J ♣Q ♦K] -> [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♠7]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoDiamonds, poker.TwoClubs, poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.NineClubs, poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.FiveHearts, poker.FiveDiamonds, poker.SevenSpades)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u1, m)
	}
	{ // P2 更新手牌 [♦3 ♠4 ♣4 ♠5 ♣5 ♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] -> [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.FourSpades, poker.FiveSpades)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeDiamonds, poker.FourClubs, poker.FiveClubs)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u2, m)
	}

	{ // Round 1 P1
		{ // P1 抽開牌區 [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♠7] ♥J
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
			game.GameMessage(u1, m)
		}
		{ // P1 丟出 ♠7
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.SevenSpades)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
		{ // P1 更新手牌 [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♠7] -> [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♥J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoDiamonds, poker.TwoClubs, poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.NineClubs, poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.FiveHearts, poker.FiveDiamonds, poker.JackHearts)),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	{ // Round 1 P2
		{ // P2 抽開牌區 [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♣3
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
		{ // P2 丟出 ♦7
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.SevenDiamonds)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
		{ // P2 更新手牌 [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] -> [♦3 ♠4 ♠5][♣3 ♣4 ♣5][♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeClubs, poker.FourClubs, poker.FiveClubs)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeDiamonds, poker.FourSpades, poker.FiveSpades)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs)),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	{ // Round 2 P1
		{ // P1 抽開牌區 [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♥J] ♣A
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u1, m)
		}
		{ // P1 丟出 ♦2
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.TwoDiamonds)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
		{ // P1 更新手牌 [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♥J] -> [♣A ♣2 ♣3][♠3 ♥3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♥J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.AceClubs, poker.TwoClubs, poker.ThreeClubs)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeSpades, poker.ThreeHearts)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.NineClubs, poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.FiveHearts, poker.FiveDiamonds, poker.JackHearts)),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	{ // Round 2 P2
		{ // P2 抽開牌區 [♦3 ♠4 ♠5][♣3 ♣4 ♣5][♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♠5
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
		{ // P2 丟出 ♦7
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.FiveSpades)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	{ // Round 3 P1
		{ // P1 抽開牌區 [♣A ♣2 ♣3][♠3 ♥3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♥J] ♠5
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
			game.GameMessage(u1, m)
		}
		{ // P1 宣勝 丟出 ♣9 [♣A ♣2 ♣3][♦10 ♦J ♣Q ♦K][♠3 ♥3 ♥J][♠5 ♥5 ♦5]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeclareID)
			data := new(netproto.INRM_Declare)
			data.Throw = poker.TransePokerCard(poker.NineClubs)
			data.Hands = new(netproto.Poker_Hand)
			data.Hands.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.AceClubs, poker.TwoClubs, poker.ThreeClubs)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeSpades, poker.ThreeHearts, poker.JackHearts)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.FiveSpades, poker.FiveHearts, poker.FiveDiamonds)),
			}
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	{ // P2 最後一次更新手牌
		{ // P2 更新手牌 [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] -> [♦3 ♠4 ♠5][♣3 ♣4 ♣5][♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeClubs, poker.FourClubs, poker.FiveClubs)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeDiamonds, poker.FourSpades, poker.FiveSpades)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs)),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	select {
	case <-end:
		assertWinLoss(t, ctx, u1, orgScore, 2650, 0, 0, 132)
		assertWinLoss(t, ctx, u2, orgScore, -2650, 2650, 2650, 0)
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_兩人_一人斷線一回合_一人抽牌後Miss(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 15 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		game.Game.RandJoker = poker.Jack
		fmt.Println(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	r1p1, r1p2 := make(chan bool, 1), make(chan bool, 1)
	r2p1, r2p2 := make(chan bool, 1), make(chan bool, 1)
	r3p1, r3p2 := make(chan bool, 1), make(chan bool, 1)
	r4p1, r4p2 := make(chan bool, 1), make(chan bool, 1)

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		switch netproto.INRM_GameMessageClassID(m.SClassID) {
		case netproto.INRM_GameMessageClassID_INRM_GameStateID:
			data := m.MsgData.(*netproto.INRM_GameState)
			if userID == u1.UserID && data.PlayerID == userID && u2.UserState.StateName() == gamelogic.StateNameGame {
				return
			}
			if userID == u2.UserID && data.PlayerID == userID && u1.UserState.StateName() == gamelogic.StateNameGame {
				return
			}
			if data.Round == 1 {
				if data.PlayerID == u1.UserID {
					r1p1 <- true
				} else if data.PlayerID == u2.UserID {
					r1p2 <- true
				}
			} else if data.Round == 2 {
				if data.PlayerID == u1.UserID {
					r2p1 <- true
				} else if data.PlayerID == u2.UserID {
					r2p2 <- true
				}
			} else if data.Round == 3 {
				if data.PlayerID == u1.UserID {
					r3p1 <- true
				} else if data.PlayerID == u2.UserID {
					r3p2 <- true
				}
			} else if data.Round == 4 {
				if data.PlayerID == u1.UserID {
					r4p1 <- true
				} else if data.PlayerID == u2.UserID {
					r4p2 <- true
				}
			}
		}
	}

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	{ // P1 更新手牌 [♦2 ♣2 ♠3 ♥3 ♣3 ♥5 ♦5 ♠7 ♣9 ♦10 ♦J ♣Q ♦K] -> [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♠7]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoDiamonds, poker.TwoClubs, poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.NineClubs, poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.FiveHearts, poker.FiveDiamonds, poker.SevenSpades)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u1, m)
	}
	{ // P2 更新手牌 [♦3 ♠4 ♣4 ♠5 ♣5 ♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] -> [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.FourSpades, poker.FiveSpades)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeDiamonds, poker.FourClubs, poker.FiveClubs)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u2, m)
	}

	<-r1p1
	fmt.Println("Round 1 P1", game.Game)
	fmt.Println(game.Game.Player(u1.UserID))
	{ // Round 1 抽一張斷線
		{ // P1 抽開牌區 [♦2 ♣2 ♠3 ♥3 ♣3][♣9 ♦10 ♦J ♣Q ♦K][♥5 ♦5 ♠7] ♥J
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
			game.GameMessage(u1, m)
		}

		u1.UserState.OnUserOffLine(u1, game.TableManager.Room)
	}

	<-r1p2
	fmt.Println("Round 1 P2", game.Game)
	fmt.Println(game.Game.Player(u2.UserID))
	{ // Round 1 P2
		{ // P2 抽開牌區 [♠4 ♠5][♦3 ♣4 ♣5][♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♣3
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
	}

	<-r2p1
	fmt.Println("Round 2 P1", game.Game)
	{ // Round 2 P1 還沒回來
	}

	<-r2p2
	fmt.Println("Round 2 P2", game.Game)
	fmt.Println(game.Game.Player(u2.UserID))
	{ // Round 2 P2 (P1 回來)
		{ // P2 抽開牌區 [♠4 ♠5] [♦3 ♣4 ♣5] [♦7 ♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♣3
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
			game.GameMessage(u2, m)
		}
		game.UserRecover(u1)
		{ // P2 丟出 ♦7 [♦3 ♠4 ♠5] [♣3 ♣4 ♣5] [♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.SevenDiamonds)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
		{ // P2 更新手牌 [♦3 ♠4 ♠5] [♣3 ♣4 ♣5] [♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♣3 ♣4 ♣5")
			d2, _ := poker.ParseDeck("♦3 ♠4 ♠5")
			d3, _ := poker.ParseDeck("♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	<-r3p1
	fmt.Println("Round 3 P1", game.Game)
	fmt.Println(game.Game.Player(u1.UserID))
	{ // Round 3 P1
		{ // P1 抽開牌區 [♦2 ♣2 ♠3 ♥3 ♣3] [♣9 ♦10 ♦J ♣Q ♦K] [♥5 ♦5 ♠7] ♣A
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u1, m)
		}
		{ // P1 丟出 ♦K [♣A ♣2 ♣3][♦2 ♠3 ♥3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠7]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.KingDiamonds)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
		{ // P1 更新手牌 [♣A ♣2 ♣3][♦2 ♠3 ♥3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠7]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
			d2, _ := poker.ParseDeck("♦2 ♠3 ♥3")
			d3, _ := poker.ParseDeck("♣9 ♦10 ♦J ♣Q")
			d4, _ := poker.ParseDeck("♥5 ♦5 ♠7")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	<-r3p2
	fmt.Println("Round 3 P2", game.Game)
	fmt.Println(game.Game.Player(u2.UserID))
	{ // Round 3 P2
		{ // P2 抽開牌區 [♣3 ♣4 ♣5] [♦3 ♠4 ♠5] [♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♠5
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
		{ // P2 丟出 ♠5 [♣3 ♣4 ♣5] [♦3 ♠4 ♠5] [♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.AceClubs)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	<-r4p1
	fmt.Println("Round 4 P1", game.Game)
	fmt.Println(game.Game.Player(u1.UserID))
	{ // Round 4 P1
		{ // P1 抽開牌區 [♣A ♣2 ♣3] [♦2 ♠3 ♥3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠7] ♠5
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopOpenDeckID)
			game.GameMessage(u1, m)
		}
		{ // P1 丟出 ♠7 [♣A ♣2 ♣3] [♦2 ♠3 ♥3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠5]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_ThrowCardID)
			data := new(netproto.INRM_ThrowCard)
			data.Throw = poker.TransePokerCard(poker.SevenSpades)
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
		{ // P1 更新手牌 [♣A ♣2 ♣3] [♦2 ♠3 ♥3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠5]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
			d2, _ := poker.ParseDeck("♦2 ♠3 ♥3")
			d3, _ := poker.ParseDeck("♣9 ♦10 ♦J ♣Q")
			d4, _ := poker.ParseDeck("♥5 ♦5 ♠5")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	<-r4p2
	fmt.Println("Round 4 P2", game.Game)
	fmt.Println(game.Game.Player(u2.UserID))
	{ // Round 4 P2
		{ // P2 抽開牌區 [♣3 ♣4 ♣5] [♦3 ♠4 ♠5] [♦7 ♥8 ♦8 ♦10 ♣10 ♣10 ♣J] ♠9
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
		{ // P2 宣勝 丟出 ♦8 [♣3 ♣4 ♣5] [♦3 ♠4 ♠5] [♦7 ♥8 ♠9] [♦10 ♣10 ♣10 ♣J]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeclareID)
			data := new(netproto.INRM_Declare)
			data.Throw = poker.TransePokerCard(poker.EightDiamonds)
			data.Hands = new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♣3 ♣4 ♣5")
			d2, _ := poker.ParseDeck("♦3 ♠4 ♠5")
			d3, _ := poker.ParseDeck("♦7 ♥8 ♠9")
			d4, _ := poker.ParseDeck("♦10 ♣10 ♣10 ♣J")

			data.Hands.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	{ // P1 最後一次更新手牌
		{ // P1 更新手牌 [♣A ♣2 ♣3] [♣9 ♦10 ♦J ♣Q] [♥5 ♦5 ♠5] [♦2 ♠3 ♥3]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♣A ♣2 ♣3")
			d2, _ := poker.ParseDeck("♣9 ♦10 ♦J ♣Q")
			d3, _ := poker.ParseDeck("♥5 ♦5 ♠5")
			d4, _ := poker.ParseDeck("♦2 ♠3 ♥3")

			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
			}

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}
	}

	select {
	case <-end:
		// assertWinLoss(t, ctx, u1, orgScore, 2650, 0, 0, 132)
		// assertWinLoss(t, ctx, u2, orgScore, -2650, 2650, 2650, 0)
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_三人_一人詐胡一人放棄(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 30 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	{ // P1 更新手牌 [♣2 ♦3 ♣3 ♥4 ♣4 ♦5 ♠6 ♦8 ♠9 ♦J ♣J ♣Q ♦K] -> [♣2 ♣3 ♣4][♦3 ♥4 ♦5][♠6 ♦8 ♠9 ♦J ♣J ♣Q ♦K]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoClubs, poker.ThreeClubs, poker.FourClubs)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeDiamonds, poker.FourHearts, poker.FiveDiamonds)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.SixSpades, poker.EightDiamonds, poker.NineSpades, poker.JackDiamonds, poker.JackClubs, poker.QueenClubs, poker.KingDiamonds)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u1, m)
	}

	{ // P2 更新手牌 [♥2 ♥3 ♣3 ♠4 ♠5 ♥5 ♣5 ♠6 ♣9 ♦10 ♣10 ♠K ♣A] -> [♠4 ♠5 ♠6][♥2 ♥3 ♣3 ♥5 ♣5 ♣9 ♦10 ♣10 ♠K ♣A]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.FourSpades, poker.FiveSpades, poker.SixSpades)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoHearts, poker.ThreeHearts, poker.ThreeClubs, poker.FiveHearts, poker.FiveClubs, poker.NineClubs, poker.TenDiamonds, poker.TenClubs, poker.KingSpades, poker.AceClubs)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u2, m)
	}

	{ // P3 更新手牌 [🃟R ♦2 ♠3 ♣4 ♠5 ♠7 ♦7 ♦7 ♠8 ♥8 ♦8 ♦10 ♣10] -> [♦2 ♠3 ♣4 ♠5][♠7 ♦7 ♦7][♠8 ♥8 ♦8][🃟R ♦10 ♣10]
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_Game)
		m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
		data := new(netproto.Poker_Hand)

		data.Decks = []*netproto.Poker_Deck{
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.TwoDiamonds, poker.ThreeSpades, poker.FourClubs, poker.FiveSpades)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.SevenSpades, poker.SevenDiamonds, poker.SevenDiamonds)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.EightSpades, poker.EightHearts, poker.EightDiamonds)),
			poker.TransePokerDeck(poker.NewEmptyDeck(poker.RedJoker, poker.TenDiamonds, poker.TenClubs)),
		}

		m.MsgData, _ = proto.Marshal(data)
		game.GameMessage(u3, m)
	}

	{ // Round 1 P1
		{ // P1 放棄
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupID)
			game.GameMessage(u1, m)
		}
	}
	{ // Round 1 P2
		{ // P2 抽開牌區 [♠4 ♠5 ♠6][♥2 ♥3 ♣3 ♥5 ♣5 ♣9 ♦10 ♣10 ♠K ♣A] ♠3
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_PopCloseDeckID)
			game.GameMessage(u2, m)
		}
		{ // P2 丟出 ♥2
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_DeclareID)
			data := new(netproto.INRM_Declare)
			data.Throw = poker.TransePokerCard(poker.TwoHearts)
			data.Hands = new(netproto.Poker_Hand)
			data.Hands.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeSpades, poker.FourSpades, poker.FiveSpades, poker.SixSpades)),
				poker.TransePokerDeck(poker.NewEmptyDeck(poker.ThreeHearts, poker.ThreeClubs, poker.FiveHearts, poker.FiveClubs, poker.NineClubs, poker.TenDiamonds, poker.TenClubs, poker.KingSpades, poker.AceClubs)),
			}
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u2, m)
		}
	}

	select {
	case <-end:
		assertWinLoss(t, ctx, u1, orgScore, -1000, 1000, 1000, 0)
		assertWinLoss(t, ctx, u2, orgScore, -4000, 4000, 4000, 0)
		assertWinLoss(t, ctx, u3, orgScore, 5000, 0, 0, 250)
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_三人_下一輪_兩人準備_一人發呆(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 30 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v", userID, m)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	{ // Round 1 P1
		{ // P1 放棄
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupID)
			game.GameMessage(u1, m)
		}
	}
	{ // Round 1 P2
		{ // P2 放棄
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupID)
			game.GameMessage(u2, m)
		}
	}

	<-end

	{ // P1 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p1.Conn, m)
	}
	{ // P2 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p2.Conn, m)
	}

	select {
	case <-start:
		if game.TableUserCount() != 2 {
			t.Errorf("Table User Count want %d, but %d", 2, game.TableUserCount())
		}
		if u3.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameFree, u3.UserState.StateName())
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_三人_下一輪_兩人準備_一人斷線(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 30 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v", userID, m)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	// P3 斷線
	{
		t.Log(u3.UserState.OnUserOffLine(u3, game.TableManager.Room))
		t.Log(u3)
	}

	<-end

	{ // P1 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p1.Conn, m)
	}
	{ // P2 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p2.Conn, m)
	}

	select {
	case <-start:
		if game.TableUserCount() != 2 {
			t.Errorf("Table User Count want %d, but %d", 2, game.TableUserCount())
		}
		if u3.UserState.StateName() != gamelogic.StateNameFree {
			t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameFree, u3.UserState.StateName())
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_三人_下一輪_兩人準備_一人換桌(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 30 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v", userID, m)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)
	u3, p3 := fixtures.NewGameUser(3000, "阿飛", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)
	grs.PlayerManager.SetPlayer(p3)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)
	u3.UserState.OnUserQueue(u3, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u3, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	<-start

	{ // Round 1 P1
		{ // P1 放棄
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupID)
			game.GameMessage(u1, m)
		}
	}
	{ // Round 1 P2
		{ // P2 放棄
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_GiveupID)
			game.GameMessage(u2, m)
		}
	}

	<-end

	{ // P1 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p1.Conn, m)
	}
	{ // P2 準備
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_GameReadyID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p2.Conn, m)
	}
	{ // P3 換桌
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_UserQueueID)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, p3.Conn, m)
	}

	select {
	case <-start:
		if game.TableUserCount() != 2 {
			t.Errorf("Table User Count want %d, but %d", 2, game.TableUserCount())
		}
		if u3.UserState.StateName() != gamelogic.StateNameQueue {
			t.Errorf("User 3 State want %s, but %s", gamelogic.StateNameQueue, u3.UserState.StateName())
		}
	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_更新手牌_Case1(t *testing.T) {

	if testing.Short() {
		t.Skip("Skiped")
		return
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 15 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		game.Game.RandJoker = poker.Eight
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	select {
	case <-start:

		up1, _ := game.Game.Player(u1.UserID)

		up1.Deck, _ = poker.ParseDeck("♦8 ♦10 ♦Q ♦K ♥6 ♥7 ♦2 ♦3 ♣5 ♣9 ♣K ♠A 🃟R")

		{ // P1 更新手牌 ["♦8 ♦10 ♦Q ♦K"] ["♥6 ♥7"] ["♦2 ♦3"] ["♣5 ♣9 ♣K"] ["♠A 🃟R"]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♦8 ♦10 ♦Q ♦K")
			d2, _ := poker.ParseDeck("♥6 ♥7")
			d3, _ := poker.ParseDeck("♦2 ♦3")
			d4, _ := poker.ParseDeck("♣5 ♣9 ♣K")
			d5, _ := poker.ParseDeck("♠A 🃟R")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
				poker.TransePokerDeck(d5),
			}

			t.Log(data)

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}

		t.Log(up1)

	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_更新手牌_Case2(t *testing.T) {

	if testing.Short() {
		t.Skip("Skiped")
		return
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 15 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		game.Game.RandJoker = poker.Eight
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	select {
	case <-start:

		up1, _ := game.Game.Player(u1.UserID)

		up1.Deck, _ = poker.ParseDeck("♦8 ♦10 ♦Q ♥K ♦K ♣K ♥6 ♥7 ♦2 ♦3 ♣5 ♠A 🃟R")

		{ // P1 更新手牌 ["♦8 ♦10 ♦Q"] ["♥K ♦K ♣K"] ["♥6 ♥7"] ["♦2 ♦3"] ["♣5"] ["♠A 🃟R"]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♦8 ♦10 ♦Q")
			d2, _ := poker.ParseDeck("♥K ♦K ♣K")
			d3, _ := poker.ParseDeck("♥6 ♥7")
			d4, _ := poker.ParseDeck("♦2 ♦3")
			d5, _ := poker.ParseDeck("♣5")
			d6, _ := poker.ParseDeck("♠A 🃟R")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
				poker.TransePokerDeck(d5),
				poker.TransePokerDeck(d6),
			}

			t.Log(data)

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}

		t.Log(up1)

	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_更新手牌_Case3(t *testing.T) {

	if testing.Short() {
		t.Skip("Skiped")
		return
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(game.GetUserByUserID(userid), score, tax)
	}
	game.GameOptions = append(game.GameOptions, Shuffler(poker.NewRandShuffler(rand.NewSource(1))))
	game.GameOptions = append(game.GameOptions, ShuffleTimes(1))

	game.ActionWaitTime = 15 * time.Second
	game.StartWaitTime = 15 * time.Second
	var ctx *gamelog.RoundContext
	game.RegisterAfterGameStartFunc(func() {
		ctx = game.GameLog.Context()
		game.Game.RandJoker = poker.Ten
		t.Log(ctx)
	})

	start, end := make(chan bool, 1), make(chan bool, 1)
	game.RegisterAfterGameStartFunc(func() { start <- true })
	game.RegisterAfterGameEndFunc(func() { end <- true })

	if testing.Short() {
		mlog.GetCLogger().Enable(false)
		game.SetGameLog(new(gamelog.EmptyGameLog))
	} else {
		conf.CLogLevel = mlog.InfoLevel
		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			t.Logf("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}
	}

	var orgScore int64 = 100000
	u1, p1 := fixtures.NewGameUser(1000, "阿呆", orgScore)
	u2, p2 := fixtures.NewGameUser(2000, "阿瓜", orgScore)

	grs.PlayerManager.SetPlayer(p1)
	grs.PlayerManager.SetPlayer(p2)

	u1.UserState.OnUserQueue(u1, game.TableManager.Room)
	u2.UserState.OnUserQueue(u2, game.TableManager.Room)

	game.UserEnter(u1, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
	game.UserEnter(u2, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	select {
	case <-start:

		up1, _ := game.Game.Player(u1.UserID)

		up1.Deck, _ = poker.ParseDeck("♥2 ♥6 ♥9 ♣3 ♣3 ♣10 ♠3 ♠2 ♠4 🃟R 🃟R ♦2 ♦5 ♦4")

		{ // P1 更新手牌 ["♥2 ♥6 ♥9"] ["♣3 ♣3 ♣10"] ["♠3 ♠2 ♠4"] ["🃟R 🃟R"] ["♦2 ♦5 ♦4"]
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.INRM_GameMessageClassID_INRM_RefreshHandID)
			data := new(netproto.Poker_Hand)

			d1, _ := poker.ParseDeck("♥2 ♥6 ♥9")
			d2, _ := poker.ParseDeck("♣3 ♣3 ♣10")
			d3, _ := poker.ParseDeck("♠3 ♠2 ♠4")
			d4, _ := poker.ParseDeck("🃟R 🃟R")
			d5, _ := poker.ParseDeck("♦2 ♦5 ♦4")
			data.Decks = []*netproto.Poker_Deck{
				poker.TransePokerDeck(d1),
				poker.TransePokerDeck(d2),
				poker.TransePokerDeck(d3),
				poker.TransePokerDeck(d4),
				poker.TransePokerDeck(d5),
			}

			t.Log(data)

			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u1, m)
		}

		t.Log(up1)
		t.Log(game.Summary())

	case <-time.After(3 * time.Second):
		t.Errorf("timeout")
	}
}

func TestGameServer_TimerManager(t *testing.T) {

	if testing.Short() {
		t.Skip("Skiped")
		return
	}

	conf.TimerFastMode = 1

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	g0 := grs.GameLogicManager.GetGameLogic(0).(*GameServer)
	g1 := grs.GameLogicManager.GetGameLogic(1).(*GameServer)
	g2 := grs.GameLogicManager.GetGameLogic(2).(*GameServer)
	g3 := grs.GameLogicManager.GetGameLogic(3).(*GameServer)
	g4 := grs.GameLogicManager.GetGameLogic(4).(*GameServer)

	t.Logf("%p %p %p %p %p", g0.TimerFuncManager, g1.TimerFuncManager, g2.TimerFuncManager, g3.TimerFuncManager, g4.TimerFuncManager)

	g0.TimerFuncManager.NewNamedTicker(1, func() { fmt.Println("G0") }, TimerActionCD)
	g1.TimerFuncManager.NewNamedTicker(1, func() { fmt.Println("G1") }, TimerActionCD)
	g2.TimerFuncManager.NewNamedTicker(1, func() { fmt.Println("G2") }, TimerActionCD)
	g3.TimerFuncManager.NewNamedTicker(1, func() { fmt.Println("G3") }, TimerActionCD)
	g4.TimerFuncManager.NewNamedTicker(1, func() { fmt.Println("G4") }, TimerActionCD)

	time.Sleep(3 * time.Second)
	g0.TimerFuncManager.RemoveTimer(TimerActionCD)
	t.Log("Remove Timer")
	time.Sleep(5 * time.Second)
}

func assertWinLoss(t *testing.T, ctx *gamelog.RoundContext, u *gamelogic.GameUser, orgScore, win, bet, flow, tax int64) {
	if ctx == nil {
		t.Errorf("ctx can't be nil")
	}

	ud := ctx.Users[u.UserID]

	if ud == nil {
		t.Errorf("ctx dose'nt exist User[%d]", u.UserID)
	}

	if u.Score != orgScore+win-tax {
		t.Errorf("User[%d] Score want %d, but %d", u.UserID, orgScore+win-tax, u.Score)
	}

	if u.UpdateScore != win-tax {
		t.Errorf("User[%d] UpdateScore want %d, but %d", u.UserID, win-tax, u.UpdateScore)
	}

	if ud.UpdateScore != win-tax {
		t.Errorf("ctx User[%d] UpdateScore want %d, but %d", u.UserID, win-tax, ud.UpdateScore)
	}

	if u.BetAmount != bet {
		t.Errorf("User[%d] BetAmount want %d, but %d", u.UserID, bet, u.BetAmount)
	}
	if ud.BetScore != bet {
		t.Errorf("ctx User[%d] BetScore want %d, but %d", u.UserID, bet, ud.BetScore)
	}

	if u.FlowAmount != bet {
		t.Errorf("User[%d] FlowAmount want %d, but %d", u.UserID, bet, u.FlowAmount)
	}
	if ud.ValidScore != bet {
		t.Errorf("ctx User[%d] ValidScore want %d, but %d", u.UserID, bet, ud.ValidScore)
	}

	if u.TaxAmount != tax {
		t.Errorf("User[%d] TaxAmount want %d, but %d", u.UserID, tax, u.TaxAmount)
	}
	if ud.TaxAmount != tax {
		t.Errorf("ctx User[%d] TaxAmount want %d, but %d", u.UserID, tax, ud.TaxAmount)
	}

	if win > 0 {
		if u.PayAmount != win-tax {
			t.Errorf("User[%d] PayAmount want %d, but %d", u.UserID, win-tax, u.PayAmount)
		}
		if ud.PayoutScore != win {
			t.Errorf("ctx User[%d] PayoutScore want %d, but %d", u.UserID, win, ud.PayoutScore)
		}
	} else {
		if u.PayAmount != 0 {
			t.Errorf("User[%d] PayAmount want %d, but %d", u.UserID, 0, u.PayAmount)
		}
		if ud.PayoutScore != 0 {
			t.Errorf("ctx User[%d] PayoutScore want %d, but %d", u.UserID, 0, ud.PayoutScore)
		}
	}
}
