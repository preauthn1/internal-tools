package s888

import (
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func Test888Logic_WinLineMultiSymbol(t *testing.T) {
	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user := fixtures.GetGameUser()[0]
	user.Score = 10000
	user.UserType = gamelogic.UserTypeNormal

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
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

	t.Run("Win", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("s6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("s6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10005001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(1)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)

	})
}
