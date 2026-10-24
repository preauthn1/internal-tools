package hhf

import (
	"testing"
	"time"

	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestHHFLogic_WinFreeSpin12ToEnd(t *testing.T) {
	// conf.CLogLevel = mlog.TraceLevel
	// conf.PrintMessage = true

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	t.Run("Win Free Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})

	t.Run("Free Spin 1, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("ss1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 2, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 3, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 4, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 5, New Link", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 6, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 2 {
			tt.Errorf("FreeSpin lock should be 2")
		}
	})

	t.Run("Free Spin 7, Lock New Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 3 {
			tt.Errorf("FreeSpin lock should be 3")
		}
	})

	t.Run("Free Spin 8, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 3 {
			tt.Errorf("FreeSpin lock should be 3")
		}
	})

	t.Run("Free Spin 9, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 3 {
			tt.Errorf("FreeSpin lock should be 3")
		}
	})

	t.Run("Free Spin 10, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 3 {
			tt.Errorf("FreeSpin lock should be 3")
		}
	})

	t.Run("Free Spin 11, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if len(logic.locks) != 3 {
			tt.Errorf("FreeSpin lock should be 3")
		}
	})

	t.Run("Free Spin 12, Lock Symbol", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010002), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin != nil {
			tt.Errorf("FreeSpin should be finished")
		}

		if len(logic.locks) != 0 {
			tt.Errorf("FreeSpin lock should be 0")
		}
	})

	t.Run("Normal Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin != nil {
			tt.Errorf("FreeSpin should be finished")
		}

		if len(logic.locks) != 0 {
			tt.Errorf("FreeSpin lock should be 0")
		}
	})
}

func TestHHFLogic_WinFreeSpinTwoW2(t *testing.T) {
	// conf.CLogLevel = mlog.TraceLevel
	// conf.PrintMessage = true

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	t.Run("Win Free Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}
	})
}

func TestHHFLogic_WinFreeSpinLeftRightTwoRule(t *testing.T) {
	// conf.CLogLevel = mlog.TraceLevel
	// conf.PrintMessage = true

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	// w2,s4,s5,w1,s3|s4,s4,s3,w2,w2|s3,s5,s2,w1,s1

	t.Run("Win Free Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if logic.FreeSpin.Count != 12 {
			tt.Errorf("FreeSpin count want %d, but %d", 12, logic.FreeSpin.Count)
		}
	})
}

func TestHHFLogic_WinFreeSpinTwoRule(t *testing.T) {
	// conf.CLogLevel = mlog.TraceLevel
	// conf.PrintMessage = true

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	// w2,s4,s5,w1,s3|s4,s4,s3,w2,w2|s3,s5,s2,w1,s1

	t.Run("Win Free Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("s4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("s3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("s2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if logic.FreeSpin.Count != 6 {
			tt.Errorf("FreeSpin count want %d, but %d", 6, logic.FreeSpin.Count)
		}
	})
}

func TestHHFLogic_WinFreeSpinSingleRule(t *testing.T) {
	// conf.CLogLevel = mlog.TraceLevel
	// conf.PrintMessage = true

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)
	}

	// w2,w1,s5,s5,s5|s4,s4,s3,w2,w2|s3,s5,s2,w1,s1

	t.Run("Win Free Spin", func(tt *testing.T) {
		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("s5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(10010001), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("a1"), Proportion: proto.Uint32(1)},
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

		if logic.FreeSpin == nil {
			tt.Errorf("FreeSpin should not be nil")
		}

		if logic.FreeSpin.Count != 6 {
			tt.Errorf("FreeSpin count want %d, but %d", 6, logic.FreeSpin.Count)
		}
	})
}
