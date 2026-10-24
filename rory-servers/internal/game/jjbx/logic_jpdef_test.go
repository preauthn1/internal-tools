package jjbx

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/defender"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJJBXLogic_JPDefender_JPGame(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
			data := m.MsgData.(*netproto.Slot_MSG_Status)
			if data.GetState() == 0 {
				c <- true
				close(c)
			}
			return
		}
	}

	user.UserState.OnUserQueue(user, logic.GetTable().Room)
	go logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	<-c

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	t.Run("Win Jackpot Game", func(tt *testing.T) {

		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(705), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		logic.Proto.SlotPayTable = []*protodata.SlotPayTable{
			{Id: proto.Int32(324), GroupId: proto.Int32(7), LinkType: proto.String("any"), Symbol: proto.String("w1"), PayType: proto.String("mysteryJPGR"), Count: proto.Int32(2), Payout: proto.String("8:100:0:1e4,18:200:0:1e4,38:300:0:1e4,68:400:0:1e4,88:1e4:0:1e4"), UseType: proto.String("all")},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(100)
		m.Rate = proto.Int32(88)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		defender.JPDefender = defender.NewJackpotDefender(100, slot.JackpotGrand, slot.JackpotMajor, slot.JackpotMini, slot.JackpotMinor)
		logic.GameMessage(logic.GameUser, _m)
		defender.JPDefender = nil

		time.Sleep(10 * time.Millisecond)

		if logic.JackpotGame != nil {
			t.Errorf("Jackpot Game should be defensed")
		}
	})
}
