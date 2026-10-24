package gen

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gamequeue"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
)

func TestRoom_FackMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	grs.RoomActionManager.SetGameQueue(gamequeue.NewUserChooseGameQueue())

	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	u, p := fixtures.NewGameUser(10000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p)
	u.UserState.OnUserQueue(u, game.TableManager.Room)

	var send = func() {
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_ROOM_Read_TableInfoID)
		ans := grs.RoomActionManager.OnMessage(grs.ModuleCtx, nil, m)
		for _, m := range ans {
			fmt.Println(m)
			data := m.MsgData.(*netproto.GEN_RoomInfo)
			fmt.Println(data.Odds)
			for _, rm := range data.Tables {
				fmt.Println(rm)
			}
		}
	}

	t.Run("UserEnter", func(t *testing.T) {
		game.UserEnter(u, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)
		send()
	})

	t.Run("UserLeave", func(t *testing.T) {
		game.UserLeave(u)
		send()
	})
}

func TestRoom_RealMessage(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// conf.CLogLevel = mlog.InfoLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	grs.RoomActionManager.SetGameQueue(gamequeue.NewUserChooseGameQueue())
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	clt := network.NewConnInfo(nil, "UnitTest", 0)
	p1 := player.NewGuest(clt)
	grs.PlayerManager.SetPlayer(p1)
	{ // 收到 p1 登入成功封包
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_LoginRoomRetID)

		rpcm := new(rpcmsg.RPCMessage)
		rpcm.RPCInfo = new(netproto.RPCInfo)
		rpcm.RPCInfo.ConnID = proto.String(p1.Conn.ConnID())

		data := new(netproto.LoginGameRoomRet)
		data.Code = proto.Int32(1)
		data.Message = proto.String("OK")

		data.LoginRequestData = new(netproto.LoginGameRoomInfo)
		data.LoginRequestData.UserID = proto.Int32(1000)
		data.LoginRequestData.Cer = proto.String("Cer")
		data.LoginRequestData.HDCode = proto.String("HDCode")
		data.LoginRequestData.HDType = proto.Int32(0)

		data.UserData = new(netproto.UserRoomLogonData)
		data.UserData.UserID = proto.Int32(1000)
		data.UserData.UserType = proto.Int32(gamelogic.UserTypeNormal)
		data.UserData.NickName = proto.String("阿呆")
		data.UserData.FaceID = proto.Int32(1)
		data.UserData.Sex = proto.Int32(1)
		data.UserData.CashAmount = proto.Int64(10000000)
		data.UserData.WinCount = proto.Int32(0)
		data.UserData.LoseCount = proto.Int32(0)
		data.UserData.DrawCount = proto.Int32(0)
		data.UserData.TradeId = proto.Int64(10000000)
		data.UserData.SumCharge = proto.Int64(0)
		data.UserData.HasBindPay = proto.Bool(false)
		data.UserData.SumConvert = proto.Int64(0)

		rpcm.MsgData, _ = proto.Marshal(data)

		m.MsgData = rpcm

		grs.RoomActionManager.OnRPCMessage(m, grs.ModuleCtx)
	}

	{
		m := new(msg.Message)
		m.BClassID = int32(netproto.MessageBClassID_GameRoom)
		m.SClassID = int32(netproto.GameRoomClassID_ROOM_ChooseTableID)
		data := new(netproto.UserChooseTable)
		data.TableNo = proto.Int32(0)
		data.SeatNo = proto.Int32(0)
		m.MsgData, _ = proto.Marshal(data)
		grs.RoomActionManager.OnMessage(grs.ModuleCtx, clt, m)
	}
}

func TestNewGame(t *testing.T) {
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	expRangeOpt := map[pb.ModeType][][2]int{
		pb.ModeType_A: {
			{1, 193},
			{502, 565},
			{902, 965},
			{1302, 1365},
			{1666, 1921},
		},
		pb.ModeType_B: {
			{1, 193},
			{402, 465},
			{802, 865},
			{1202, 1265},
			{1502, 1565},
			{1794, 1921},
		},
	}

	expCntOpt := [][2]int{
		{20, 30},
		{30, 20},
	}

	if !reflect.DeepEqual(expRangeOpt, game.ChanceRangeOptMap) {
		t.Errorf("NewGame() ChanceRangeOptMap want %v, but %v", expRangeOpt, game.ChanceRangeOptMap)
	}
	if !reflect.DeepEqual(expCntOpt, game.ChanceCntOpt) {
		t.Errorf("NewGame() ChanceCntOpt want %v, but %v", expCntOpt, game.ChanceCntOpt)
	}
}

func TestGame_Play_None(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	var score int64 = 100000

	u, _ := fixtures.NewGameUser(1000, "阿呆", score)

	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u, score, tax)
	}
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	game.PachiSlotGame.RT = 888
	game.PachiSlotGame.NextResult = pb.ResultType_None

	var odds int64 = 1000
	var bet int64 = 3000

	game.Play(odds, u)

	if u.Score != score-bet {
		t.Errorf("User Score want %d, but %d", score-bet, u.Score)
	}

	if u.BetAmount != bet {
		t.Errorf("User BetAmount want %d, but %d", bet, u.BetAmount)
	}

	if u.FlowAmount != bet {
		t.Errorf("User FlowAmount want %d, but %d", bet, u.FlowAmount)
	}
}

func TestGame_Play_3松解除RT(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	var score int64 = 100000

	u, _ := fixtures.NewGameUser(1000, "阿呆", score)

	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
		s := new(scorestorage.ConsoleUserScoreStorage)
		s.UpdateUserScore(u, score, tax)
	}
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	Proto := game.PachiSlotGame.Proto.Clone()
	Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
		{GroupID: pachislot.DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_WaterMelonLink, Rate: 100, Denominator: 100, Times: []int32{2}, UseType: pb.Scenario_ALL},
	}
	game.PachiSlotGame, _ = pachislot.NewGame(Proto, game.PachiSlotGame.Pattern.ID)

	game.SendMessgeFunc = func(userID int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
	}

	game.PachiSlotGame.RT = 888
	game.PachiSlotGame.NextResult = pb.ResultType_WaterMelonLink

	var odds int64 = 1000
	var bet int64 = 3000
	var win int64 = 5000

	game.Play(odds, u)

	if u.Score != score+win-bet {
		t.Errorf("User Score want %d, but %d", score+win-bet, u.Score)
	}

	if u.PayAmount != win {
		t.Errorf("User PayAmount want %d, but %d", win, u.PayAmount)
	}

	if u.BetAmount != bet {
		t.Errorf("User BetAmount want %d, but %d", bet, u.BetAmount)
	}

	if u.FlowAmount != bet {
		t.Errorf("User FlowAmount want %d, but %d", bet, u.FlowAmount)
	}

	if game.PachiSlotGame.RT != 2 {
		t.Log(game.PachiSlotGame.RTReleseConfigs)
		t.Errorf("PachiGame RT want %d, but %d", 2, game.PachiSlotGame.RT)
	}

	if game.PachiSlotGame.NextMode != pb.ModeType_Heaven {
		t.Log(game.PachiSlotGame.RTReleseConfigs)
		t.Errorf("PachiGame NextMode want %v, but %v", pb.ModeType_Heaven, game.PachiSlotGame.NextMode)
	}
}

func TestGame_Play_BB贏得JAC(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	t.Run("0次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 1
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0
		game.PachiSlotGame.NextResult = pb.ResultType_None

		game.Play(odds, u)

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 2 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 2, game.PachiSlotGame.BonusStore)
		}
	})

	t.Run("1次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 1
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink

		game.Play(odds, u)

		for game.PachiSlotGame.Scenario == pb.Scenario_JAC {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
		}

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 1 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 1, game.PachiSlotGame.BonusStore)
		}
	})

	t.Run("2次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 2
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		for game.PachiSlotGame.Scenario == pb.Scenario_BB {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			for game.PachiSlotGame.Scenario == pb.Scenario_JAC {
				game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
				game.Play(odds, u)
			}
		}

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 0 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 0, game.PachiSlotGame.BonusStore)
		}
	})

	t.Run("3次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 3
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		for game.PachiSlotGame.Scenario == pb.Scenario_BB {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			for game.PachiSlotGame.Scenario == pb.Scenario_JAC {
				game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
				game.Play(odds, u)
			}
		}

		if game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", false, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 0 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 0, game.PachiSlotGame.BonusStore)
		}
	})

	t.Run("4次", func(t *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 4
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		for game.PachiSlotGame.Scenario == pb.Scenario_BB {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			for game.PachiSlotGame.Scenario == pb.Scenario_JAC {
				game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
				game.Play(odds, u)
			}
		}

		if game.PachiSlotGame.BonusSure {
			t.Errorf("PachiGame BonusSure want %v, but %v", false, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 0 {
			t.Errorf("PachiGame BonusStore want %v, but %v", 0, game.PachiSlotGame.BonusStore)
		}
	})
}

func TestGame_Play_JAC連續8次銅鐘(t *testing.T) {

	// conf.CLogLevel = mlog.InfoLevel
	t.Run("1次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 1
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		t.Log("In BB", game.BellLinkInJAC, game.JACInBB)
		for j := 0; j < 1; j++ {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			t.Log("Win JAC", game.BellLinkInJAC, game.JACInBB)
			for i := 0; i < 8; i++ {
				game.PachiSlotGame.NextResult = pb.ResultType_BellLink
				game.Play(odds, u)
				t.Log("Each JAC", game.BellLinkInJAC, game.JACInBB)
			}
		}
		t.Log("After BB", game.BellLinkInJAC, game.JACInBB)

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 2 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 2, game.PachiSlotGame.BonusStore)
		}
	})
	t.Run("2次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 2
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		for j := 0; j < 2; j++ {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			for i := 0; i < 8; i++ {
				game.PachiSlotGame.NextResult = pb.ResultType_BellLink
				game.Play(odds, u)
			}
		}

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 2 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 2, game.PachiSlotGame.BonusStore)
		}
	})
	t.Run("3次", func(tt *testing.T) {
		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		var odds int64 = 1000

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.Bonus = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 3
		game.PachiSlotGame.JACTimes = 0
		game.PachiSlotGame.JACTimesLeft = 0

		for j := 0; j < 3; j++ {
			game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
			game.Play(odds, u)
			for i := 0; i < 8; i++ {
				game.PachiSlotGame.NextResult = pb.ResultType_BellLink
				game.Play(odds, u)
			}
		}

		if !game.PachiSlotGame.BonusSure {
			tt.Errorf("PachiGame BonusSure want %v, but %v", true, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.BonusStore != 2 {
			tt.Errorf("PachiGame BonusStore want %v, but %v", 2, game.PachiSlotGame.BonusStore)
		}
	})
}

func TestGame_Play_機會模式(t *testing.T) {

	t.Run("Win in NG", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
			},
			pb.ModeType_B: {
				{1, 100},
			},
			pb.ModeType_Heaven: {
				{1, 100},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.Scenario = pb.Scenario_NG
		game.PachiSlotGame.Mode = pb.ModeType_A
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		game.PachiSlotGame.RT = 999
		game.PachiSlotGame.RoadMap.SpinCnt = 55
		game.ChanceAlways = true
		game.ChanceRate = 1

		var odds int64 = 1000

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", false, game.ChanceBuff)
		}

		game.Play(odds, u)

		if !game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", true, game.ChanceBuff)
		}
	})
	t.Run("In BB", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		game.PachiSlotGame.RoadMap.SpinCnt = 55
		game.ChanceAlways = true

		var odds int64 = 1000

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", false, game.ChanceBuff)
		}

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}
	})
	t.Run("In RB", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.Scenario = pb.Scenario_RB
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		game.PachiSlotGame.RoadMap.SpinCnt = 55
		game.ChanceAlways = true

		var odds int64 = 1000

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", false, game.ChanceBuff)
		}

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}
	})
	t.Run("In JAC", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.Scenario = pb.Scenario_JAC
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		game.PachiSlotGame.RoadMap.SpinCnt = 55
		game.ChanceAlways = true

		var odds int64 = 1000

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", false, game.ChanceBuff)
		}

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}
	})
	t.Run("Not In Range", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		game.PachiSlotGame.RoadMap.SpinCnt = 101
		game.ChanceAlways = true

		var odds int64 = 1000

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", false, game.ChanceBuff)
		}

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}
	})
	t.Run("解除", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.NextResult = pb.ResultType_None
		game.PachiSlotGame.RoadMap.SpinCnt = 100
		game.ChanceBuff = true
		game.ChanceDeadLine = 100

		var odds int64 = 1000

		if !game.ChanceBuff {
			tt.Errorf("ChanceBuff before spin want %t, but %t", true, game.ChanceBuff)
		}

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}
	})
	t.Run("Assign", func(tt *testing.T) {
		// conf.CLogLevel = mlog.DebugLevel

		grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

		var score int64 = 100000

		u, _ := fixtures.NewGameUser(1000, "阿呆", score)

		game := grs.GameLogicManager.GetGameLogic(0).(*Game)
		game.UpdateUserScoreFunc = func(userid int32, score, tax int64) {
			s := new(scorestorage.ConsoleUserScoreStorage)
			s.UpdateUserScore(u, score, tax)
		}
		if testing.Short() {
			game.SetGameLog(new(gamelog.EmptyGameLog))
		}

		game.SendMessgeFunc = func(userID int32, m *msg.Message) {
			example.invalid("玩家[%d] 收到封包 %v, %v", userID, m, m.MsgData)
		}

		game.ChanceRangeOptMap = map[pb.ModeType][][2]int{
			pb.ModeType_A: {
				{1, 100},
				{201, 300},
				{401, 500},
				{801, 1000},
			},
		}
		game.ChanceCntOpt = [][2]int{
			{20, 50},
			{30, 50},
			{40, 50},
			{50, 50},
		}

		game.PachiSlotGame.Scenario = pb.Scenario_BB
		game.PachiSlotGame.BonusTimesLeft = 1
		game.PachiSlotGame.NextResult = pb.ResultType_None
		game.ChanceBuff = true
		game.ChanceDeadLine = 1000

		var odds int64 = 1000

		game.Play(odds, u)

		if game.ChanceBuff {
			tt.Errorf("ChanceBuff after spin want %t, but %t", false, game.ChanceBuff)
		}

		if game.PachiSlotGame.Mode != pb.ModeType_Heaven {
			tt.Errorf("Mode want %v, but %v", pb.ModeType_Heaven, game.PachiSlotGame.Mode)
		}
	})
}

func TestGame_封包一連串(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}

	u, p := fixtures.NewGameUser(10000, "阿呆", 100000)

	var printMessage = func(uid int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包%v, %v", uid, m, m.MsgData)
	}

	game.PachiSlotGame.Bonus = pb.Scenario_RB
	pachislot.DefaultBonusSureWinRate = 100
	game.PachiSlotGame.RT = 999

	// game.OnSendMessage(printMessage)

	grs.PlayerManager.SetPlayer(p)
	u.UserState.OnUserQueue(u, game.TableManager.Room)

	game.UserEnter(u, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	t.Run("查詢遊戲設定", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameSettingID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameSettingID), m.SClassID)
			}
			_, ok := m.MsgData.(*netproto.GEN_GameSetting)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameSetting, but %v", reflect.TypeOf(m.MsgData))
			}
			done <- true
		})

		{ // 發送查詢封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_Read_GameSettingID)
			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("查詢轉數表", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID), m.SClassID)
			}
			_, ok := m.MsgData.(*netproto.GEN_RoadMap)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_RoadMap, but %v", reflect.TypeOf(m.MsgData))
			}
			done <- true
		})

		{ // 發送查詢封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_Read_RoadMapID)
			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("查詢遊戲狀態", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			_, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			done <- true
		})

		{ // 發送查詢封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_Read_GameStatusID)
			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("下注失敗。倍率錯誤", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_ErrorID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_ErrorID), m.SClassID)
			}
			err, ok := m.MsgData.(*netproto.GEN_Error)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_Error, but %v", reflect.TypeOf(m.MsgData))
			}
			if err.Type != netproto.GEN_ErrorType_GEN_ErrBet_IllegalOdds {
				t.Errorf("Error.Type want %v, but %v", netproto.GEN_ErrorType_GEN_ErrBet_IllegalOdds, err.Type)
			}
			if err.Source != netproto.GEN_GameMessageClassID_GEN_BetID {
				t.Errorf("Error.Source want %v, but %v", netproto.GEN_GameMessageClassID_GEN_BetID, err.Source)
			}
			t.Log(err.GetMessage())
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = 99
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("下注失敗。金幣不足", func(t *testing.T) {
		tmpScore := u.Score
		defer func() {
			u.Score = tmpScore
		}()
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		u.Score = 299
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_ErrorID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_ErrorID), m.SClassID)
			}
			err, ok := m.MsgData.(*netproto.GEN_Error)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_Error, but %v", reflect.TypeOf(m.MsgData))
			}
			if err.Type != netproto.GEN_ErrorType_GEN_ErrBet_MoneyNotEnough {
				t.Errorf("Error.Type want %v, but %v", netproto.GEN_ErrorType_GEN_ErrBet_MoneyNotEnough, err.Type)
			}
			if err.Source != netproto.GEN_GameMessageClassID_GEN_BetID {
				t.Errorf("Error.Source want %v, but %v", netproto.GEN_GameMessageClassID_GEN_BetID, err.Source)
			}
			t.Log(err.GetMessage())
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = 1000
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(999)，槓龜", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != odds*3 {
				t.Errorf("BetMoney want %d, but %d", odds*3, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 97000
			var expUpdate int64 = -3000
			var expBet int64 = 3000
			var expWin int64 = 0
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != expBet {
				t.Errorf("User BetAmount want %d, but %d", expBet, u.BetAmount)
			}
			if u.PayAmount != expWin {
				t.Errorf("User PayAmount want %d, but %d", expWin, u.PayAmount)
			}
			if u.FlowAmount != expFlow {
				t.Errorf("User FlowAmount want %d, but %d", expFlow, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(998)，贏得Replay", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if !status.Replay {
				t.Errorf("Status.Replay want %t, but %t", true, status.Replay)
			}
			if status.RequireCoin != 0 {
				t.Errorf("Status.RequireCoin want %d, but %d", 0, status.RequireCoin)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 94000
			var expUpdate int64 = -3000
			var expBet int64 = 3000
			var expWin int64 = 0
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 6000 {
				t.Errorf("User BetAmount want %d, but %d", 6000, u.BetAmount)
			}
			if u.PayAmount != expWin {
				t.Errorf("User PayAmount want %d, but %d", expWin, u.PayAmount)
			}
			if u.FlowAmount != 6000 {
				t.Errorf("User FlowAmount want %d, but %d", 6000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(997)，免錢，純槓解除RT", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_FeatureNone
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 0 {
				t.Errorf("BetMoney want %d, but %d", 0, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if !result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", true, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 94000
			var expUpdate int64 = 0
			var expBet int64 = 0
			var expWin int64 = 0
			var expFlow int64 = 0
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 6000 {
				t.Errorf("User BetAmount want %d, but %d", 6000, u.BetAmount)
			}
			if u.PayAmount != expWin {
				t.Errorf("User PayAmount want %d, but %d", expWin, u.PayAmount)
			}
			if u.FlowAmount != 6000 {
				t.Errorf("User FlowAmount want %d, but %d", 6000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 33 {
				t.Errorf("RT Release By [%v], Left RT want %d but %d", pb.ResultType_FeatureNone, 33, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(33)，銅鐘連線10倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 10000 {
				t.Errorf("WinMoney want %d, but %d", 10000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 101000
			var expUpdate int64 = 7000
			var expBet int64 = 3000
			var expWin int64 = 10000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 9000 {
				t.Errorf("User BetAmount want %d, but %d", 9000, u.BetAmount)
			}
			if u.PayAmount != expWin {
				t.Errorf("User PayAmount want %d, but %d", expWin, u.PayAmount)
			}
			if u.FlowAmount != 9000 {
				t.Errorf("User FlowAmount want %d, but %d", 9000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 32 {
				t.Errorf("RT Release By [%v], Left RT want %d but %d", pb.ResultType_FeatureNone, 32, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(32)，西瓜連線5倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_WaterMelonLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 5000 {
				t.Errorf("WinMoney want %d, but %d", 5000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 103000
			var expUpdate int64 = 2000
			var expBet int64 = 3000
			var expWin int64 = 5000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 12000 {
				t.Errorf("User BetAmount want %d, but %d", 12000, u.BetAmount)
			}
			if u.PayAmount != 15000 {
				t.Errorf("User PayAmount want %d, but %d", 15000, u.PayAmount)
			}
			if u.FlowAmount != 12000 {
				t.Errorf("User FlowAmount want %d, but %d", 12000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 31 {
				t.Errorf("RT Release By [%v], Left RT want %d but %d", pb.ResultType_FeatureNone, 31, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(31)，櫻桃出現左欄2倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BerryLeftAny
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 2000 {
				t.Errorf("WinMoney want %d, but %d", 2000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 102000
			var expUpdate int64 = -1000
			var expBet int64 = 3000
			var expWin int64 = 2000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 15000 {
				t.Errorf("User BetAmount want %d, but %d", 15000, u.BetAmount)
			}
			if u.PayAmount != 17000 {
				t.Errorf("User PayAmount want %d, but %d", 17000, u.PayAmount)
			}
			if u.FlowAmount != 15000 {
				t.Errorf("User FlowAmount want %d, but %d", 15000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 30 {
				t.Errorf("RT Release By [%v], Left RT want %d but %d", pb.ResultType_FeatureNone, 30, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(30-1)，跑30場槓龜", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		var odds int64 = 1000
		for i := 0; i < 30; i++ {
			game.PachiSlotGame.NextResult = pb.ResultType_None

			{ // 發送封包
				m := new(msg.Message)
				m.BClassID = int32(netproto.MessageBClassID_Game)
				m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
				data := new(netproto.GEN_Bet)
				data.Odds = odds
				m.MsgData, _ = proto.Marshal(data)

				game.GameMessage(u, m)
			}
		}

		bet := 30 * odds * 3

		if u.Score != 102000-bet {
			t.Errorf("User Score want %d, but %d", 102000-bet, u.Score)
		}
		if u.BetAmount != 15000+bet {
			t.Errorf("User BetAmount want %d, but %d", 15000+bet, u.BetAmount)
		}
		if u.PayAmount != 17000 {
			t.Errorf("User PayAmount want %d, but %d", 17000, u.PayAmount)
		}
		if u.FlowAmount != 15000+bet {
			t.Errorf("User FlowAmount want %d, but %d", 15000+bet, u.FlowAmount)
		}
		if u.TaxAmount != 0 {
			t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
		}

		if game.PachiSlotGame.RT != 0 {
			t.Errorf("RT Release By [%v], Left RT want %d but %d", pb.ResultType_FeatureNone, 0, game.PachiSlotGame.RT)
		}
	})
	t.Run("NG(0)，RT達成", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.DebugLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 9000
			var expUpdate int64 = -3000
			var expBet int64 = 3000
			var expWin int64 = 0
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 108000 {
				t.Errorf("User BetAmount want %d, but %d", 108000, u.BetAmount)
			}
			if u.PayAmount != 17000 {
				t.Errorf("User PayAmount want %d, but %d", 17000, u.PayAmount)
			}
			if u.FlowAmount != 108000 {
				t.Errorf("User FlowAmount want %d, but %d", 108000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if !game.PachiSlotGame.BonusSure {
				t.Errorf("PachiSlotGame.BonusSure want %t but %t", true, game.PachiSlotGame.BonusSure)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(?)，獲得RB", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 15000 {
				t.Errorf("WinMoney want %d, but %d", 15000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 21000
			var expUpdate int64 = 12000
			var expBet int64 = 3000
			var expWin int64 = 15000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 111000 {
				t.Errorf("User BetAmount want %d, but %d", 111000, u.BetAmount)
			}
			if u.PayAmount != 32000 {
				t.Errorf("User PayAmount want %d, but %d", 32000, u.PayAmount)
			}
			if u.FlowAmount != 111000 {
				t.Errorf("User FlowAmount want %d, but %d", 111000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusSure {
				t.Errorf("PachiSlotGame.BonusSure want %t but %t", false, game.PachiSlotGame.BonusSure)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_RB {
				t.Errorf("PachiSlotGame.Scenario want %v but %v", pb.Scenario_RB, game.PachiSlotGame.Scenario)
			}

		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("查詢轉數表", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_RoadMapID), m.SClassID)
			}
			rm, ok := m.MsgData.(*netproto.GEN_RoadMap)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_RoadMap, but %v", reflect.TypeOf(m.MsgData))
			}

			if rm.Total != 38 {
				t.Errorf("Total want %d, but %d", 38, rm.Total)
			}
			if rm.Spin != 0 {
				t.Errorf("Spin want %d, but %d", 0, rm.Spin)
			}
			if len(rm.Bonus) != 2 {
				t.Errorf("Bonus count want %d, but %d", 2, len(rm.Bonus))
			}
			var rmb *netproto.GEN_RoadMapBonus
			for _, b := range rm.Bonus {
				if b.Type == netproto.GEN_Scenario_GEN_Scenario_RB {
					rmb = b
				}
			}
			if rmb.Type != netproto.GEN_Scenario_GEN_Scenario_RB {
				t.Errorf("Bonus[RB] Type want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_RB, rmb.Type)
			}
			if rmb.Count != 1 {
				t.Errorf("Bonus[RB] Count want %d, but %d", 1, rmb.Count)
			}
			if len(rm.Snapshots) != 1 {
				t.Errorf("Snapshots count want %d, but %d", 1, len(rm.Snapshots))
			}
			if rm.Snapshots[0].Type != netproto.GEN_Scenario_GEN_Scenario_RB {
				t.Errorf("Snapshots[0] Type want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_RB, rm.Snapshots[0].Type)
			}
			if rm.Snapshots[0].Count != 38 {
				t.Errorf("Snapshots[0] Count want %d, but %d", 38, rm.Snapshots[0].Count)
			}
			done <- true
		})

		{ // 發送查詢封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_Read_RoadMapID)
			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("查詢遊戲狀態", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)
		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if status.Scenario != netproto.GEN_Scenario_GEN_Scenario_RB {
				t.Errorf("Scenario want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_RB, status.Scenario)
			}
			if status.Bonus != netproto.GEN_Scenario_GEN_Scenario_RB {
				t.Errorf("Bonus want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_RB, status.Bonus)
			}
			if status.BonusLeft != 8 {
				t.Errorf("BonusLeft want %v, but %v", 8, status.BonusLeft)
			}
			done <- true
		})

		{ // 發送查詢封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_Read_GameStatusID)
			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("RB(8)，槓龜", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != odds*1 {
				t.Errorf("BetMoney want %d, but %d", odds*1, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 20000
			var expUpdate int64 = -1000
			var expBet int64 = 1000
			var expWin int64 = 0
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 112000 {
				t.Errorf("User BetAmount want %d, but %d", 112000, u.BetAmount)
			}
			if u.PayAmount != 32000 {
				t.Errorf("User PayAmount want %d, but %d", 32000, u.PayAmount)
			}
			if u.FlowAmount != 112000 {
				t.Errorf("User FlowAmount want %d, but %d", 112000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusTimesLeft != 7 {
				t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 7, game.PachiSlotGame.BonusTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("RB(7)，Replay連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != odds*1 {
				t.Errorf("BetMoney want %d, but %d", odds*1, result.BetMoney)
			}
			if result.WinMoney != odds*15 {
				t.Errorf("WinMoney want %d, but %d", odds*15, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 34000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 113000 {
				t.Errorf("User BetAmount want %d, but %d", 113000, u.BetAmount)
			}
			if u.PayAmount != 47000 {
				t.Errorf("User PayAmount want %d, but %d", 47000, u.PayAmount)
			}
			if u.FlowAmount != 113000 {
				t.Errorf("User FlowAmount want %d, but %d", 113000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusTimesLeft != 6 {
				t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 6, game.PachiSlotGame.BonusTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("RB(6)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != odds*1 {
				t.Errorf("BetMoney want %d, but %d", odds*1, result.BetMoney)
			}
			if result.WinMoney != odds*15 {
				t.Errorf("WinMoney want %d, but %d", odds*15, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 48000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 114000 {
				t.Errorf("User BetAmount want %d, but %d", 114000, u.BetAmount)
			}
			if u.PayAmount != 62000 {
				t.Errorf("User PayAmount want %d, but %d", 62000, u.PayAmount)
			}
			if u.FlowAmount != 114000 {
				t.Errorf("User FlowAmount want %d, but %d", 114000, u.FlowAmount)
			}
			if u.TaxAmount != expTax {
				t.Errorf("User TaxAmount want %d, but %d", expTax, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusTimesLeft != 5 {
				t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 5, game.PachiSlotGame.BonusTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("RB(5-1)，槓龜5場", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		var odds int64 = 1000

		for i := 0; i < 5; i++ {
			game.PachiSlotGame.NextResult = pb.ResultType_None

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		if u.Score != 43000 {
			t.Errorf("User Score want %d, but %d", 43000, u.Score)
		}
		if u.BetAmount != 119000 {
			t.Errorf("User BetAmount want %d, but %d", 119000, u.BetAmount)
		}
		if u.PayAmount != 62000 {
			t.Errorf("User PayAmount want %d, but %d", 62000, u.PayAmount)
		}
		if u.FlowAmount != 119000 {
			t.Errorf("User FlowAmount want %d, but %d", 119000, u.FlowAmount)
		}
		if u.TaxAmount != 0 {
			t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
		}

		if game.PachiSlotGame.BonusTimesLeft != 0 {
			t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 0, game.PachiSlotGame.BonusTimesLeft)
		}
		if game.PachiSlotGame.Scenario != 0 {
			t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_NG, game.PachiSlotGame.Scenario)
		}
	})
	game.PachiSlotGame.Bonus = pb.Scenario_BB
	game.PachiSlotGame.RT = 1234
	t.Run("NG(1234)，西瓜連線5倍", func(t *testing.T) {

		tmp := game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink]
		defer func() {
			game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink] = tmp
		}()
		game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink] = &pb.RTReleaseOption{ResultType: pb.ResultType_WaterMelonLink, Times: []int32{2}, Rate: 0, Denominator: 100, UseType: pb.Scenario_ALL}

		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_WaterMelonLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				printMessage(uid, m)
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 5000 {
				t.Errorf("WinMoney want %d, but %d", 5000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 45000
			var expUpdate int64 = 2000
			var expBet int64 = 3000
			var expWin int64 = 5000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 122000 {
				t.Errorf("User BetAmount want %d, but %d", 122000, u.BetAmount)
			}
			if u.PayAmount != 67000 {
				t.Errorf("User PayAmount want %d, but %d", 67000, u.PayAmount)
			}
			if u.FlowAmount != 122000 {
				t.Errorf("User FlowAmount want %d, but %d", 122000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(1233)，ReleaseRT", func(t *testing.T) {

		tmp := game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonBB]
		defer func() {
			game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonBB] = tmp
		}()
		game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonBB] = &pb.RTReleaseOption{ResultType: pb.ResultType_ReleaseRTonBB, Times: []int32{2}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL}

		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReleaseRT
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if !result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", true, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 42000
			var expUpdate int64 = -3000
			var expBet int64 = 3000
			var expWin int64 = 0
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 125000 {
				t.Errorf("User BetAmount want %d, but %d", 125000, u.BetAmount)
			}
			if u.PayAmount != 67000 {
				t.Errorf("User PayAmount want %d, but %d", 67000, u.PayAmount)
			}
			if u.FlowAmount != 125000 {
				t.Errorf("User FlowAmount want %d, but %d", 125000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 2 {
				t.Errorf("PachiSlotGame.RT want %d, but %d", 2, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(2)，櫻桃出現左欄2倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BerryLeftAny
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 2000 {
				t.Errorf("WinMoney want %d, but %d", 2000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 41000
			var expUpdate int64 = -1000
			var expBet int64 = 3000
			var expWin int64 = 2000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 128000 {
				t.Errorf("User BetAmount want %d, but %d", 128000, u.BetAmount)
			}
			if u.PayAmount != 69000 {
				t.Errorf("User PayAmount want %d, but %d", 69000, u.PayAmount)
			}
			if u.FlowAmount != 128000 {
				t.Errorf("User FlowAmount want %d, but %d", 128000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 1 {
				t.Errorf("PachiSlotGame.RT want %d but %d", 1, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(1)，贏得Replay", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if !status.Replay {
				t.Errorf("Status.Replay want %t, but %t", true, status.Replay)
			}
			if status.RequireCoin != 0 {
				t.Errorf("Status.RequireCoin want %d, but %d", 0, status.RequireCoin)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 38000
			var expUpdate int64 = -3000
			var expBet int64 = 3000
			var expWin int64 = 0
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 131000 {
				t.Errorf("User BetAmount want %d, but %d", 131000, u.BetAmount)
			}
			if u.PayAmount != 69000 {
				t.Errorf("User PayAmount want %d, but %d", 69000, u.PayAmount)
			}
			if u.FlowAmount != 131000 {
				t.Errorf("User FlowAmount want %d, but %d", 131000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.RT != 0 {
				t.Errorf("PachiSlotGame.RT want %d but %d", 0, game.PachiSlotGame.RT)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(0)，免錢，RT達成", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 0 {
				t.Errorf("BetMoney want %d, but %d", 0, result.BetMoney)
			}
			if result.WinMoney != 0 {
				t.Errorf("WinMoney want %d, but %d", 0, result.WinMoney)
			}
			if result.Result != netproto.GEN_ResultType_GEN_RT_BonusSure {
				t.Errorf("Result want %v, but %v", netproto.GEN_ResultType_GEN_RT_BonusSure, result.Result)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 38000
			var expUpdate int64 = 0
			var expBet int64 = 0
			var expWin int64 = 0
			var expFlow int64 = 0
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 131000 {
				t.Errorf("User BetAmount want %d, but %d", 131000, u.BetAmount)
			}
			if u.PayAmount != 69000 {
				t.Errorf("User PayAmount want %d, but %d", 69000, u.PayAmount)
			}
			if u.FlowAmount != 131000 {
				t.Errorf("User FlowAmount want %d, but %d", 131000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if !game.PachiSlotGame.BonusSure {
				t.Errorf("PachiSlotGame.BonusSure want %t but %t", true, game.PachiSlotGame.BonusSure)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_NG {
				t.Errorf("PachiSlotGame.Scenario want %v but %v", pb.Scenario_NG, game.PachiSlotGame.Scenario)
			}

		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(?)，獲得BB", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != 3000 {
				t.Errorf("BetMoney want %d, but %d", 3000, result.BetMoney)
			}
			if result.WinMoney != 15000 {
				t.Errorf("WinMoney want %d, but %d", 15000, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 50000
			var expUpdate int64 = 12000
			var expBet int64 = 3000
			var expWin int64 = 15000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 134000 {
				t.Errorf("User BetAmount want %d, but %d", 134000, u.BetAmount)
			}
			if u.PayAmount != 84000 {
				t.Errorf("User PayAmount want %d, but %d", 84000, u.PayAmount)
			}
			if u.FlowAmount != 134000 {
				t.Errorf("User FlowAmount want %d, but %d", 134000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusSure {
				t.Errorf("PachiSlotGame.BonusSure want %t but %t", false, game.PachiSlotGame.BonusSure)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_BB {
				t.Errorf("PachiSlotGame.Scenario want %v but %v", pb.Scenario_BB, game.PachiSlotGame.Scenario)
			}

		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("BB(30)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameResultID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameResultID), m.SClassID)
			}
			result, ok := m.MsgData.(*netproto.GEN_GameResult)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameResult, but %v", reflect.TypeOf(m.MsgData))
			}
			if result.BetMoney != odds*3 {
				t.Errorf("BetMoney want %d, but %d", odds*3, result.BetMoney)
			}
			if result.WinMoney != odds*15 {
				t.Errorf("WinMoney want %d, but %d", odds*15, result.WinMoney)
			}
			if result.ReleaseRT {
				t.Errorf("ReleaseRT want %t, but %t", false, result.ReleaseRT)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 62000
			var expUpdate int64 = 12000
			var expBet int64 = 3000
			var expWin int64 = 15000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 137000 {
				t.Errorf("User BetAmount want %d, but %d", 137000, u.BetAmount)
			}
			if u.PayAmount != 99000 {
				t.Errorf("User PayAmount want %d, but %d", 99000, u.PayAmount)
			}
			if u.FlowAmount != 137000 {
				t.Errorf("User FlowAmount want %d, but %d", 137000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.BonusTimesLeft != 29 {
				t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 29, game.PachiSlotGame.BonusTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("BB(29)，Replay連線15倍，中JAC", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if status.Replay {
				t.Errorf("Status.Replay want %t, but %t", false, status.Replay)
			}
			if status.RequireCoin != 1 {
				t.Errorf("Status.RequireCoin want %d, but %d", 1, status.RequireCoin)
			}
			if status.Scenario != netproto.GEN_Scenario_GEN_Scenario_JAC {
				t.Errorf("Status.Scenario want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_JAC, status.Scenario)
			}
			if status.BonusLeft != 28 {
				t.Errorf("Status.BonusLeft want %v, but %v", 28, status.BonusLeft)
			}
			if status.JACLeft != 8 {
				t.Errorf("Status.JACLeft want %v, but %v", 8, status.JACLeft)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 74000
			var expUpdate int64 = 12000
			var expBet int64 = 3000
			var expWin int64 = 15000
			var expFlow int64 = 3000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 140000 {
				t.Errorf("User BetAmount want %d, but %d", 140000, u.BetAmount)
			}
			if u.PayAmount != 114000 {
				t.Errorf("User PayAmount want %d, but %d", 114000, u.PayAmount)
			}
			if u.FlowAmount != 140000 {
				t.Errorf("User FlowAmount want %d, but %d", 140000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 8 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 8, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(8)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 88000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 141000 {
				t.Errorf("User BetAmount want %d, but %d", 141000, u.BetAmount)
			}
			if u.PayAmount != 129000 {
				t.Errorf("User PayAmount want %d, but %d", 129000, u.PayAmount)
			}
			if u.FlowAmount != 141000 {
				t.Errorf("User FlowAmount want %d, but %d", 141000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 7 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 7, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(7)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 102000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 142000 {
				t.Errorf("User BetAmount want %d, but %d", 142000, u.BetAmount)
			}
			if u.PayAmount != 144000 {
				t.Errorf("User PayAmount want %d, but %d", 144000, u.PayAmount)
			}
			if u.FlowAmount != 142000 {
				t.Errorf("User FlowAmount want %d, but %d", 142000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 6 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 6, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(6)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 116000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 143000 {
				t.Errorf("User BetAmount want %d, but %d", 143000, u.BetAmount)
			}
			if u.PayAmount != 159000 {
				t.Errorf("User PayAmount want %d, but %d", 159000, u.PayAmount)
			}
			if u.FlowAmount != 143000 {
				t.Errorf("User FlowAmount want %d, but %d", 143000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 5 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 5, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(5)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 130000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 144000 {
				t.Errorf("User BetAmount want %d, but %d", 144000, u.BetAmount)
			}
			if u.PayAmount != 174000 {
				t.Errorf("User PayAmount want %d, but %d", 174000, u.PayAmount)
			}
			if u.FlowAmount != 144000 {
				t.Errorf("User FlowAmount want %d, but %d", 144000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 4 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 4, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(4)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 144000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 145000 {
				t.Errorf("User BetAmount want %d, but %d", 145000, u.BetAmount)
			}
			if u.PayAmount != 189000 {
				t.Errorf("User PayAmount want %d, but %d", 189000, u.PayAmount)
			}
			if u.FlowAmount != 145000 {
				t.Errorf("User FlowAmount want %d, but %d", 145000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 3 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 3, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(3)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 158000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 146000 {
				t.Errorf("User BetAmount want %d, but %d", 146000, u.BetAmount)
			}
			if u.PayAmount != 204000 {
				t.Errorf("User PayAmount want %d, but %d", 204000, u.PayAmount)
			}
			if u.FlowAmount != 146000 {
				t.Errorf("User FlowAmount want %d, but %d", 146000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 2 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 2, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(2)，銅鐘連線15倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 172000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 147000 {
				t.Errorf("User BetAmount want %d, but %d", 147000, u.BetAmount)
			}
			if u.PayAmount != 219000 {
				t.Errorf("User PayAmount want %d, but %d", 219000, u.PayAmount)
			}
			if u.FlowAmount != 147000 {
				t.Errorf("User FlowAmount want %d, but %d", 147000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 1 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 1, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(1)，銅鐘連線15倍，1G連", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BellLink
		ctx := game.GameLog.Context()

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			var expScore int64 = 186000
			var expUpdate int64 = 14000
			var expBet int64 = 1000
			var expWin int64 = 15000
			var expFlow int64 = 1000
			var expTax int64 = 0

			if u.Score != expScore {
				t.Errorf("User Score want %d, but %d", expScore, u.Score)
			}
			if u.BetAmount != 148000 {
				t.Errorf("User BetAmount want %d, but %d", 148000, u.BetAmount)
			}
			if u.PayAmount != 234000 {
				t.Errorf("User PayAmount want %d, but %d", 234000, u.PayAmount)
			}
			if u.FlowAmount != 148000 {
				t.Errorf("User FlowAmount want %d, but %d", 148000, u.FlowAmount)
			}
			if u.TaxAmount != 0 {
				t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
			}

			ud := ctx.Users[u.UserID]

			if ud.Score != expScore {
				t.Errorf("GameLog.Context() Score want %d, but %d", expScore, u.Score)
			}
			if ud.UpdateScore != expUpdate {
				t.Errorf("GameLog.Context() UpdateScore want %d, but %d", expUpdate, ud.UpdateScore)
			}
			if ud.BetScore != expBet {
				t.Errorf("GameLog.Context() BetScore want %d, but %d", expBet, ud.BetScore)
			}
			if ud.PayoutScore != expWin {
				t.Errorf("GameLog.Context() PayoutScore want %d, but %d", expWin, ud.PayoutScore)
			}
			if ud.ValidScore != expFlow {
				t.Errorf("GameLog.Context() ValidScore want %d, but %d", expFlow, ud.ValidScore)
			}
			if ud.TaxAmount != expTax {
				t.Errorf("GameLog.Context() TaxAmount want %d, but %d", expTax, ud.TaxAmount)
			}

			if game.PachiSlotGame.Scenario != pb.Scenario_BB {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_BB, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.JACTimesLeft != 0 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 0, game.PachiSlotGame.JACTimesLeft)
			}
			if game.PachiSlotGame.BonusStore != 1 {
				t.Errorf("PachiSlotGame.BonusStore want %d, but %d", 0, game.PachiSlotGame.BonusStore)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("BB(28-2)，銅鐘連線15倍，27場", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		var odds int64 = 1000

		for i := 0; i < 27; i++ { // 發送封包
			game.PachiSlotGame.NextResult = pb.ResultType_BellLink
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		if u.Score != 510000 {
			t.Errorf("User Score want %d, but %d", 510000, u.Score)
		}
		if u.BetAmount != 229000 {
			t.Errorf("User BetAmount want %d, but %d", 229000, u.BetAmount)
		}
		if u.PayAmount != 639000 {
			t.Errorf("User PayAmount want %d, but %d", 639000, u.PayAmount)
		}
		if u.FlowAmount != 229000 {
			t.Errorf("User FlowAmount want %d, but %d", 229000, u.FlowAmount)
		}
		if u.TaxAmount != 0 {
			t.Errorf("User TaxAmount want %d, but %d", 0, u.TaxAmount)
		}

		if game.PachiSlotGame.BonusTimesLeft != 1 {
			t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 1, game.PachiSlotGame.BonusTimesLeft)
		}
	})
	t.Run("BB(1)，Replay連線15倍，中JAC", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if status.Replay {
				t.Errorf("Status.Replay want %t, but %t", false, status.Replay)
			}
			if status.RequireCoin != 1 {
				t.Errorf("Status.RequireCoin want %d, but %d", 1, status.RequireCoin)
			}
			if status.Scenario != netproto.GEN_Scenario_GEN_Scenario_JAC {
				t.Errorf("Status.Scenario want %v, but %v", netproto.GEN_Scenario_GEN_Scenario_JAC, status.Scenario)
			}
			if status.BonusLeft != 0 {
				t.Errorf("Status.BonusLeft want %v, but %v", 0, status.BonusLeft)
			}
			if status.JACLeft != 8 {
				t.Errorf("Status.JACLeft want %v, but %v", 8, status.JACLeft)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
			if game.PachiSlotGame.Scenario != pb.Scenario_JAC {
				t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_JAC, game.PachiSlotGame.Scenario)
			}
			if game.PachiSlotGame.BonusTimesLeft != 0 {
				t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 0, game.PachiSlotGame.BonusTimesLeft)
			}
			if game.PachiSlotGame.JACTimesLeft != 8 {
				t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 8, game.PachiSlotGame.JACTimesLeft)
			}
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("JAC(8-1)，銅鐘連線15倍，8場，1G連，BB兩次JAC(1G連*1)，自動解開", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		var odds int64 = 1000

		for i := 0; i < 8; i++ { // 發送封包
			game.PachiSlotGame.NextResult = pb.ResultType_BellLink
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		if game.PachiSlotGame.Scenario != pb.Scenario_NG {
			t.Errorf("PachiSlotGame.Scenario want %v, but %v", pb.Scenario_NG, game.PachiSlotGame.Scenario)
		}
		if game.PachiSlotGame.BonusTimesLeft != 0 {
			t.Errorf("PachiSlotGame.BonusTimesLeft want %d, but %d", 0, game.PachiSlotGame.BonusTimesLeft)
		}
		if game.PachiSlotGame.JACTimesLeft != 0 {
			t.Errorf("PachiSlotGame.JACTimesLeft want %d, but %d", 0, game.PachiSlotGame.JACTimesLeft)
		}
		if game.PachiSlotGame.BonusStore != 2 {
			t.Errorf("PachiSlotGame.BonusStore want %d, but %d", 2, game.PachiSlotGame.BonusStore)
		}
		if !game.PachiSlotGame.BonusSure {
			t.Errorf("PachiSlotGame.BonusSure want %t, but %t", true, game.PachiSlotGame.BonusSure)
		}
	})
	game.PachiSlotGame.Bonus = pb.Scenario_BB
	game.PachiSlotGame.RT = 1234
	t.Run("NG(?)，獲得BB", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		if game.PachiSlotGame.BonusSure {
			t.Errorf("PachiSlotGame.BonusSure want %t but %t", false, game.PachiSlotGame.BonusSure)
		}

		if game.PachiSlotGame.Scenario != pb.Scenario_BB {
			t.Errorf("PachiSlotGame.Scenario want %v but %v", pb.Scenario_BB, game.PachiSlotGame.Scenario)
		}
	})
	t.Run("BB(30-1)，純槓1G連，30場，BB 0次JAC(1G連*3) (2+93-1)，自動解開", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		var odds int64 = 1000

		for i := 0; i < 30; i++ { // 發送封包
			game.PachiSlotGame.NextResult = pb.ResultType_FeatureNone
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		if game.PachiSlotGame.BonusStore != 94 {
			t.Errorf("PachiSlotGame.BonusStore want %d, but %d", 94, game.PachiSlotGame.BonusStore)
		}
	})
}

func TestGame_封包再一連串(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)
	if testing.Short() {
		game.SetGameLog(new(gamelog.EmptyGameLog))
	}
	game.ChanceRate = 1

	u, p := fixtures.NewGameUser(10000, "阿呆", 100000)

	var printMessage = func(uid int32, m *msg.Message) {
		example.invalid("玩家[%d] 收到封包%v, %v", uid, m, m.MsgData)
	}

	game.OnSendMessage(printMessage)

	grs.PlayerManager.SetPlayer(p)
	u.UserState.OnUserQueue(u, game.TableManager.Room)

	game.UserEnter(u, game.TableManager.TableNo, game.TableManager.GetFreeSeat().SeatNo)

	game.PachiSlotGame.Bonus = pb.Scenario_RB
	pachislot.DefaultBonusSureWinRate = 100
	game.PachiSlotGame.RoadMap.TotalCnt = 124
	game.PachiSlotGame.RoadMap.SpinCnt = 124
	game.PachiSlotGame.RT = 1234

	t.Run("NG(1234)，西瓜連線5倍，獲得機會模式", func(t *testing.T) {

		tmp := game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink]
		defer func() {
			game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink] = tmp
		}()
		game.PachiSlotGame.RTReleseConfigs[pb.ResultType_WaterMelonLink] = &pb.RTReleaseOption{ResultType: pb.ResultType_WaterMelonLink, Times: []int32{2}, Rate: 0, Denominator: 100, UseType: pb.Scenario_ALL}

		defer game.OnSendMessage(printMessage)
		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_WaterMelonLink
		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}
	})
	t.Run("NG(1233)，ReleaseRT", func(t *testing.T) {

		tmp := game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonRB]
		defer func() {
			game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonRB] = tmp
		}()
		game.PachiSlotGame.RTReleseConfigs[pb.ResultType_ReleaseRTonRB] = &pb.RTReleaseOption{ResultType: pb.ResultType_ReleaseRTonRB, Times: []int32{2}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL}

		defer game.OnSendMessage(printMessage)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReleaseRT

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

	})
	t.Run("NG(2)，櫻桃出現左欄2倍", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		// conf.CLogLevel = mlog.InfoLevel

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_BerryLeftAny
		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}
	})
	t.Run("NG(1)，贏得Replay", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		done := make(chan bool, 10)

		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_ReplayLink

		game.OnSendMessage(func(uid int32, m *msg.Message) {
			printMessage(uid, m)
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				return
			}
			if uid != u.UserID {
				t.Errorf("uid want %d, but %d", u.UserID, uid)
			}
			if m.BClassID != int32(netproto.MessageBClassID_Game) {
				t.Errorf("BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), m.BClassID)
			}
			if m.SClassID != int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID) {
				t.Errorf("SClassID want %d, but %d", int32(netproto.GEN_GameMessageClassID_GEN_GameStatusID), m.SClassID)
			}
			status, ok := m.MsgData.(*netproto.GEN_GameStatus)
			if !ok {
				t.Errorf("MsgData want *netproto.GEN_GameStatus, but %v", reflect.TypeOf(m.MsgData))
			}
			if !status.Replay {
				t.Errorf("Status.Replay want %t, but %t", true, status.Replay)
			}
			if status.RequireCoin != 0 {
				t.Errorf("Status.RequireCoin want %d, but %d", 0, status.RequireCoin)
			}
			if !status.ChanceBuff {
				t.Errorf("Status.ChanceBuff want %t, but %t", true, status.ChanceBuff)
			}
			done <- true
		})

		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}

		select {
		case <-done:
		case <-time.After(1 * time.Second):
			t.Errorf("timeout")
		}
	})
	t.Run("NG(0)，免錢，RT達成", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		// conf.CLogLevel = mlog.InfoLevel
		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)
			game.GameMessage(u, m)
		}
	})
	t.Run("NG(?)，獲得RB", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)
		// conf.CLogLevel = mlog.InfoLevel
		var odds int64 = 1000
		game.PachiSlotGame.NextResult = pb.ResultType_None
		{ // 發送封包
			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}
	})
	t.Run("RB(8-1)，槓龜8場", func(t *testing.T) {
		defer game.OnSendMessage(printMessage)

		var odds int64 = 1000

		// conf.CLogLevel = mlog.InfoLevel

		for i := 0; i < 8; i++ {
			game.PachiSlotGame.NextResult = pb.ResultType_None

			m := new(msg.Message)
			m.BClassID = int32(netproto.MessageBClassID_Game)
			m.SClassID = int32(netproto.GEN_GameMessageClassID_GEN_BetID)
			data := new(netproto.GEN_Bet)
			data.Odds = odds
			m.MsgData, _ = proto.Marshal(data)

			game.GameMessage(u, m)
		}
	})
}
