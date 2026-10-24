package wlzb

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
	"example.invalid/ServerTeam/Servers/pkg/roomserver"
	"example.invalid/ServerTeam/Servers/pkg/scorestorage"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestWLZBNewLogic(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	slot_test.PrintLogic(t, logic.SlotLogicBase)
}

func TestWLZBLogic_UserEnter(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		userEnter(logic, user, t, safe, m)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)

		printSendMessgeFunc(t, userID, m)

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
}

func TestWLZBLogic_Spin_NoWin(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	user.UserState.OnUserQueue(user, grs.Room)

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("NoWin Case1", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_Spin_NormalWin(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)
	user.UserState.OnUserQueue(user, grs.Room)

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
		// slot_test.PrintSendMessgeFunc(tt, userID, m)
	}

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Two Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_Spin_FreeGameWin(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)
	user.UserState.OnUserQueue(user, grs.Room)

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			}

			logic.FreeSpin = slot.NewFreeSpin(logic.Proto.FreeSpinTable[0], logic.Proto.SlotSymbolDistribution, 10, 1)
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_Spin_FreeGameWin2(t *testing.T) {

	roomserver.GameLogicCreator = func(t *gamelogic.Table, tfManager *event.TimerFuncManager, cfgs *config.Group) gamelogic.GameLogic {
		return NewLogic(t, tfManager, cfgs)
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)
	user.UserState.OnUserQueue(user, grs.Room)

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			}

			logic.FreeSpin = slot.NewFreeSpin(logic.Proto.FreeSpinTable[0], logic.Proto.SlotSymbolDistribution, 10, 1)
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_Spin_LastFreeGame(t *testing.T) {

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

	logic.FreeSpin = slot.NewFreeSpin(logic.Proto.FreeSpinTable[0], logic.Proto.SlotSymbolDistribution, 10, 1)
	logic.FreeSpin.Count = 2

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("FreeGame One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("FreeGame One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("NormalGame Two Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_Spin_WinSelectChance(t *testing.T) {

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("Win", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("User Choose", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.WLZB_GameMessageClassID_WLZB_ChooseGameID)

		m := new(netproto.WLZB_MSG_ChooseGame)
		m.ID = proto.Int32(9)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("FreeGame One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("w2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Win Again", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Win Again", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Win Again", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("s1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	logic.FreeSpin.Count = 1

	t.Run("FreeGame Last Time Two Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(502), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func TestWLZBLogic_UserAddMoney(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)
	user.TimeSeq = 20190909225900000
	user.UserState.OnUserQueue(user, grs.Room)

	grs.NotifyMsgHandler.AddHandler(roomlogic.NewNotifyPayOrderHandler(grs.Room.GameInfo.GameID, grs.Room.ServerID, grs.Room.RoomUserManager, new(scorestorage.ConsoleUserScoreStorage)))

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

	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {}

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("Add Money", func(tt *testing.T) {
		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_NotifyServer)
		_m.SClassID = int32(netproto.NotifyServerClassID_PayOrderID)

		m := new(netproto.PayOrderList)
		m.ServerID = proto.Int32(logic.TableManager.Room.ServerID)
		m.PayOrderList = []*netproto.PayOrderInfo{
			{
				UserID:  proto.Int32(user.UserID),
				Amount:  proto.Int64(1000),
				TimeSeq: proto.Int64(20190909225959000),
			},
		}
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		grs.NotifyMsgHandler.OnNotifyMessage(_m, grs.ModuleCtx, grs.PlayerManager)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})

	t.Run("One Symbol", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(1), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(2), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(3), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(4), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(5), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(6), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(7), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(8), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(9), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(10), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(11), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(12), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(13), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(14), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{Id: proto.Uint32(15), GroupId: proto.Uint32(501), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			}
		}

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		time.Sleep(10 * time.Millisecond)
	})
}

func userEnter(logic *Logic, user *gamelogic.GameUser, t *testing.T, safe bool, m *msg.Message) {
	handlers1 := make(map[int32]func())

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID)] = func() {
		mm := mock.EmptyFreeSpinRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID)] = func() {
		mm := mock.EmptySelectionRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID)] = func() {
		mm := mock.EmptyUserAssignedModeRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	slot_test.PrintSendRPCMessgeFunc(t, safe, m, handlers1)
}

func printSendMessgeFunc(t *testing.T, userID int32, m *msg.Message) {

	if m.SClassID == int32(netproto.WLZB_GameMessageClassID_WLZB_GameOptionsID) {
		t.Log("[GameOptionsID]")
		t.Log(" - Times:", m.MsgData.(*netproto.WLZB_MSG_ChooseGameOptions).GetTimes())
	}
}

func printSendRPCMessgeFunc(t *testing.T, safe bool, m *msg.Message, handlers map[int32]func()) {

	if handlers != nil {
		if f, ok := handlers[m.SClassID]; ok {
			f()
		}
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameRecordID) {
		t.Log("[SelectGameRecordID]")

		data := m.MsgData.(*netproto.Slot_MSG_SelectGameRecord)

		t.Log(" - GameID:", example.invalid.GetGameID())
		t.Log(" - FlagID:", example.invalid.GetFlagID())
		t.Log(" - UserID:", example.invalid.GetUserID())
		t.Log(" - OriginID:", example.invalid.GetOriginID())
		t.Log(" - Times:", example.invalid.GetTimes())
		t.Log(" - Data:", example.invalid.GetData())

		return
	}
}
