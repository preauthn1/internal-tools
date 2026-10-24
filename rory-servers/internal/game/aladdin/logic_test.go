package aladdin

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestALDNewLogic(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	slot_test.PrintLogic(t, logic.SlotLogicBase)
}

func TestALDLogic_UserEnter(t *testing.T) {

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

func TestALDLogic_GameMessage_UseCase_All(t *testing.T) {

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

	{
		logic.SlotPatternID = uint32(1)
		logic.Proto.SlotPattern = append(logic.Proto.SlotPattern, &protodata.SlotPattern{Id: proto.Uint32(1), MainGameTableGroupId: proto.Uint32(1)})

		logic.Proto.MainGameTable = []*protodata.MainGameTable{
			{GroupId: proto.Uint32(1), SymbolGroupId: proto.Uint32(1)},
		}

		logic.DefaultSymbolGroupID = 1
		logic.DefaultReelNum = 5
		logic.DefaultRowNum = 3

		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		}

		logic.Proto.SlotPayTable = []*protodata.SlotPayTable{
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(1), PayType: proto.String("respin"), Payout: proto.String("1")},
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(2), PayType: proto.String("respin"), Payout: proto.String("2")},
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(3), PayType: proto.String("respin"), Payout: proto.String("3")},
			{Symbol: proto.String("p1"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p2"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p3"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p4"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p1"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p2"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p3"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p4"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
		}

		logic.Proto.RespinTable = []*protodata.RespinTable{
			{GroupId: proto.Int32(1), RespinCount: proto.Uint32(1), AppearReel: proto.String("2"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(2)},
			{GroupId: proto.Int32(2), RespinCount: proto.Uint32(2), AppearReel: proto.String("2,4"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(3)},
			{GroupId: proto.Int32(3), RespinCount: proto.Uint32(3), AppearReel: proto.String("2,3,4"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(4)},
		}
	}

	t.Run("User Bet & Spin", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		tt.Log(logic.GameBetMoney, logic.GamePayMoney, logic.GameWinMoney, logic.GameUser.FlowAmount)
	})
}

func TestALDLogic_GameMessage_UseCase1(t *testing.T) {

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

	{
		logic.SlotPatternID = uint32(1)
		logic.Proto.SlotPattern = append(logic.Proto.SlotPattern, &protodata.SlotPattern{Id: proto.Uint32(1), MainGameTableGroupId: proto.Uint32(1)})

		logic.Proto.MainGameTable = []*protodata.MainGameTable{
			{GroupId: proto.Uint32(1), SymbolGroupId: proto.Uint32(1)},
		}

		logic.DefaultSymbolGroupID = 1
		logic.DefaultReelNum = 5
		logic.DefaultRowNum = 3

		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
		}

		logic.Proto.SlotPayTable = []*protodata.SlotPayTable{
			{Symbol: proto.String("p4"), LinkType: proto.String("left"), Count: proto.Int32(3), PayType: proto.String("chip"), Payout: proto.String("8")},
			{Symbol: proto.String("p4"), LinkType: proto.String("right"), Count: proto.Int32(3), PayType: proto.String("chip"), Payout: proto.String("8")},
		}
	}

	t.Run("User Bet & Spin", func(tt *testing.T) {
		logic.GameState = slot.GameStateBet

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(6)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		t.Log(logic.GameBetMoney, logic.GamePayMoney, logic.GameWinMoney, logic.GameUser.FlowAmount)
	})
}

func TestALDLogic_GameMessage_UseCase_ExpandTwo_OneTime(t *testing.T) {

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

	{
		logic.SlotPatternID = uint32(1)
		logic.Proto.SlotPattern = append(logic.Proto.SlotPattern, &protodata.SlotPattern{Id: proto.Uint32(1), MainGameTableGroupId: proto.Uint32(1)})

		logic.Proto.MainGameTable = []*protodata.MainGameTable{
			{GroupId: proto.Uint32(1), SymbolGroupId: proto.Uint32(1)},
		}

		logic.DefaultSymbolGroupID = 1
		logic.DefaultReelNum = 5
		logic.DefaultRowNum = 3

		logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(2), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(3), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},

			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(4), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
		}

		logic.Proto.SlotPayTable = []*protodata.SlotPayTable{
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(1), PayType: proto.String("respin"), Payout: proto.String("1")},
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(2), PayType: proto.String("respin"), Payout: proto.String("2")},
			{Symbol: proto.String("w1"), LinkType: proto.String("straight"), Count: proto.Int32(3), PayType: proto.String("respin"), Payout: proto.String("3")},
			{Symbol: proto.String("p1"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p2"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p3"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p4"), LinkType: proto.String("left"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p1"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p2"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p3"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
			{Symbol: proto.String("p4"), LinkType: proto.String("right"), Count: proto.Int32(5), PayType: proto.String("chip"), Payout: proto.String("50")},
		}

		logic.Proto.RespinTable = []*protodata.RespinTable{
			{GroupId: proto.Int32(1), RespinCount: proto.Uint32(1), AppearReel: proto.String("2"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(2)},
			{GroupId: proto.Int32(2), RespinCount: proto.Uint32(2), AppearReel: proto.String("2,4"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(3)},
			{GroupId: proto.Int32(3), RespinCount: proto.Uint32(3), AppearReel: proto.String("2,3,4"), RespinType: proto.String("Expand"), RespinValue: proto.Uint32(4)},
		}
	}

	t.Run("User Bet & Spin", func(tt *testing.T) {

		_m := new(msg.Message)
		_m.BClassID = int32(netproto.MessageBClassID_Game)
		_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

		m := new(netproto.Slot_MSG_Spin)
		m.Money = proto.Int64(10)
		b, _ := proto.Marshal(m)

		_m.MsgData = b

		logic.GameMessage(logic.GameUser, _m)

		fmt.Println(logic.GameBetMoney, logic.GamePayMoney, logic.GameWinMoney, logic.GameUser.FlowAmount)
	})
}

func TestALDLogic_Bugfix_Issue8(t *testing.T) {

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

	t.Run("Win Right Once", func(tt *testing.T) {
		{
			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
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

		time.Sleep(100 * time.Millisecond)
	})

	t.Run("Win Left Only", func(tt *testing.T) {
		{
			logic.SendMessgeFunc = func(userID int32, m *msg.Message) {

				if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {
					data := m.MsgData.(*netproto.Slot_MSG_StepGameResult)
					if len(data.SlotResult) != 1 {
						tt.Errorf("Slot Result length want %d, but %d", 1, len(data.SlotResult))
					}
					if len(data.SlotResult[0].SlotResult.PayGroups) != 1 {
						tt.Errorf("Slot Result PayGroups length want %d, but %d", 1, len(data.SlotResult[0].SlotResult.PayGroups))
					}
				}
			}

			logic.Proto.SlotSymbolDistribution = []*protodata.SlotSymbolDistribution{
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(1), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(2), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("p6"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(4), Order: proto.Uint32(3), Symbol: proto.String("p7"), Proportion: proto.Uint32(1)},
				{GroupId: proto.Uint32(301), ReelNum: proto.Uint32(5), Order: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(1)},
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

		time.Sleep(100 * time.Millisecond)
	})
}

func printSendMessgeFunc(t *testing.T, userID int32, m *msg.Message) {

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameResultID) {
		t.Log("[GameResult]")

		data := m.MsgData.(*netproto.Slot_MSG_StepGameResult)

		t.Log("  BetMoney:", data.GetBetMoney())
		t.Log("  FinalMoney:", data.GetFinalMoney())
		t.Log("  WinMoney:", data.GetWinMoney())
		t.Logf("  SlotResult: (%d)\n", len(data.GetSlotResult()))
		for _, row := range data.GetSlotResult() {
			t.Log("   - Slot_StepGameResult")
			t.Logf("     - IsExpand: %v\n", row.GetIsExpand())
			t.Logf("     - IsRespin: %v\n", row.GetIsRespin())
			t.Logf("     - ExpandReel: %v\n", row.GetExpandReel())
			t.Logf("     - ExpandOrder: %v\n", row.GetExpandOrder())
			t.Log("     - InitResult:")
			for _, _row := range row.GetInitResult() {
				t.Logf("       - %v\n", _row)
			}
			t.Log("     - SlotResult:")
			t.Log("       - BetMoney:", row.GetSlotResult().GetBetMoney())
			t.Log("       - WinMoney:", row.GetSlotResult().GetWinMoney())
			t.Log("       - FinalMoney:", row.GetSlotResult().GetFinalMoney())
			t.Log("       - JackpotMoney:", row.GetSlotResult().GetJackpotMoney())
			t.Log("       - FreeGameWinMoney:", row.GetSlotResult().GetFreeGameWinMoney())
			t.Log("       - ReelResult")
			for _, _row := range row.GetSlotResult().GetReelResult() {
				t.Logf("       - %v\n", _row)
			}
			t.Log("       - PayGroups")
			for _, _row := range row.GetSlotResult().GetPayGroups() {
				t.Logf("         - WinMoney: %d, Count: %d, Way: %d, LinkType: %s, Multi: %d, \n", _row.GetWinMoney(), _row.GetCount(), _row.GetWay(), _row.GetLinkType(), _row.GetMulti())
				t.Log("         - PayLines:", _row.GetPayLines())
				t.Log("         - BlinkPositions:")
				for _, pos := range _row.GetBlinkPositions() {
					t.Logf("           - %v\n", pos)
				}
			}
		}
		return
	}

}

func printSendRPCMessgeFunc(t *testing.T, safe bool, m *msg.Message, handlers map[int32]func()) {
	if handlers != nil {
		if f, ok := handlers[m.SClassID]; ok {
			f()
		}
	}
}
