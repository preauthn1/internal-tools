package tgw

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestTGWLogic_SpinPickGame(t *testing.T) {

	// conf.CLogLevel = mlog.DebugLevel

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_StatusID) {
			data := m.MsgData.(*netproto.DFDC_GameStatus)
			if data.GetState() == 0 {
				c <- true
				close(c)
			}
			return
		}
	}

	user.UserState.OnUserQueue(user, grs.Room)
	go logic.UserEnter(user, logic.GetTable().TableNo, logic.GetTable().GetFreeSeat().SeatNo)
	<-c

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		// slot_test.PrintSendMessgeFunc(t, userID, m)
		// printSendMessgeFunc(t, userID, m)
	}

	t.Run("Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p9"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p8"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p8"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(201), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("b1"), Proportion: proto.Uint32(1)},
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionSpinID)

		m := new(netproto.TGW_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		if logic.PickGame == nil {
			tt.Errorf("PickGame must have value")
			return
		}

		if len(logic.PickGame.All) != 12 {
			tt.Errorf("PickGame ALL should have 12 items")
		}
	})

	t.Run("Pick1", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionPickID)

		var data = &netproto.TGW_MSG_Pick{}
		data.Record = &netproto.TGW_BonusGameRecord{}
		data.Record.ID = proto.Int32(1)
		data.Record.Location = proto.Int32(3)
		data.Record.RewareValue = proto.Int64(-1)
		var detail = &netproto.Slot_MSG_PickBonus{
			Record: &netproto.Slot_BonusGameRecord{
				ID:          proto.Int32(data.GetRecord().GetID()),
				Location:    proto.Int32(data.GetRecord().GetLocation()),
				RewareValue: proto.Int64(int64(data.GetRecord().GetRewareValue())),
			},
		}

		b, _ := proto.Marshal(detail)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)
	})

	t.Run("Pick2", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionPickID)

		var data = &netproto.TGW_MSG_Pick{}
		data.Record = &netproto.TGW_BonusGameRecord{}
		data.Record.ID = proto.Int32(2)
		data.Record.Location = proto.Int32(8)
		data.Record.RewareValue = proto.Int64(-1)
		var detail = &netproto.Slot_MSG_PickBonus{
			Record: &netproto.Slot_BonusGameRecord{
				ID:          proto.Int32(data.GetRecord().GetID()),
				Location:    proto.Int32(data.GetRecord().GetLocation()),
				RewareValue: proto.Int64(int64(data.GetRecord().GetRewareValue())),
			},
		}

		b, _ := proto.Marshal(detail)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)
	})

	t.Run("Pick3", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.TGW_GameMessageClassID_TGW_ActionPickID)

		var data = &netproto.TGW_MSG_Pick{}
		data.Record = &netproto.TGW_BonusGameRecord{}
		data.Record.ID = proto.Int32(3)
		data.Record.Location = proto.Int32(6)
		data.Record.RewareValue = proto.Int64(-1)
		var detail = &netproto.Slot_MSG_PickBonus{
			Record: &netproto.Slot_BonusGameRecord{
				ID:          proto.Int32(data.GetRecord().GetID()),
				Location:    proto.Int32(data.GetRecord().GetLocation()),
				RewareValue: proto.Int64(int64(data.GetRecord().GetRewareValue())),
			},
		}

		b, _ := proto.Marshal(detail)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		// time.Sleep(1 * time.Millisecond)
		if logic.PickGame != nil {
			tt.Errorf("PickGame must be nil")
		}
	})
}
