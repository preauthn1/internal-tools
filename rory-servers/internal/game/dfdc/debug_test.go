package dfdc

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestDFDCDebug_Case1(t *testing.T) {

	conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	user.UserState.OnUserQueue(user, grs.Room)
	logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	logic.UserReady(user)

	time.Sleep(time.Millisecond)

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		printSendMessgeFunc(t, userID, m)
	}

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p14"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p8"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p14"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p13"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p13"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p12"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p9"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p13"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p11"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p14"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(9), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.DFDC_GameMessageClassID_DFDCGameStartID)

		m := new(netproto.DFDC_StartBet)
		m.Money = proto.Int64(10)
		m.GoldenSymbol = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)
	})
}
