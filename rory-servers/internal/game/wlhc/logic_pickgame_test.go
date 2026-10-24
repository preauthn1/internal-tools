package wlhc

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestWLHCLogic_PickGame_Win(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

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

	var records []*netproto.Slot_BonusGameRecord

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameID) {
			data := m.MsgData.(*netproto.Slot_MSG_BonusGame)
			records = data.Records
		}
	}

	t.Run("Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(1), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(2), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(3), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(4), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(5), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(6), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(7), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(8), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(9), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(10), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(11), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(12), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(13), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(14), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{Id: proto.Uint32(15), GroupId: proto.Uint32(1101), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(5)
		m.ID = proto.Int32(22)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)

		if logic.PickGame == nil {
			t.Errorf("Game should has PickGame")
		}
	})

	if len(records) == 0 {
		t.Errorf("Pick Game should has record")
	}

	t.Run("Pick Item", func(tt *testing.T) {

		for i := 0; i < len(records); i++ {
			_m := new(msg.Message)
			_m.BClassID = int32(netproto.MessageBClassID_Game)
			_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

			m := new(netproto.Slot_MSG_PickBonus)
			m.Type = proto.String(string(slot.GameModePickGame))
			m.Record = new(netproto.Slot_BonusGameRecord)
			m.Record.ID = proto.Int32(records[i].GetID())
			m.Record.Location = proto.Int32(int32(i + 1))
			m.Record.RewareValue = proto.Int64(records[i].GetRewareValue())

			b, _ := proto.Marshal(m)

			_m.MsgData = b

			logic.GameMessage(logic.GameUser, _m)

			time.Sleep(10 * time.Millisecond)
		}
	})
}
