package jgyf

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func userEnter(logic *Logic, user *gamelogic.GameUser, t *testing.T, safe bool, m *msg.Message) {
	handlers1 := make(map[int32]func())

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID)] = func() {
		mm := mock.EmptyFreeSpinRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)] = func() {
		mm := mock.EmptyJackpotPoolRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID)] = func() {
		mm := mock.EmptySelectionRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID)] = func() {
		detail := mock.EmptyBonusGameDetail(logic.SlotLogicBase, user)
		mm := mock.MarshalBonusGame(detail)
		logic.OnRecvRPCMessage(mm)
	}

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID)] = func() {
		mm := mock.EmptyUserAssignedModeRPCMessage(logic.SlotLogicBase, user)
		logic.OnRecvRPCMessage(mm)
	}

	slot_test.PrintSendRPCMessgeFunc(t, safe, m, handlers1)
}

func TestJGYFLogic_UserEnter_Empty(t *testing.T) {

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

func TestJGYFLogic_UserEnter_JackpotGame_Only(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	user, p1 := fixtures.NewGameUser(1000, "阿呆", 100000)
	grs.PlayerManager.SetPlayer(p1)

	logic.SendRPCMessageFunc = func(safe bool, m *msg.Message) {
		handlers1 := make(map[int32]func())

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID)] = func() {
			mm := mock.EmptyFreeSpinRPCMessage(logic.SlotLogicBase, user)
			logic.OnRecvRPCMessage(mm)
		}

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)] = func() {
			mm := mock.EmptyJackpotPoolRPCMessage(logic.SlotLogicBase, user)
			logic.OnRecvRPCMessage(mm)
		}

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID)] = func() {
			mm := mock.EmptySelectionRPCMessage(logic.SlotLogicBase, user)
			logic.OnRecvRPCMessage(mm)
		}

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID)] = func() {
			detail := mock.EmptyBonusGameDetail(logic.SlotLogicBase, user)
			mock.AppendJackpotGame(logic.SlotLogicBase, user, detail)
			mm := mock.MarshalBonusGame(detail)
			logic.OnRecvRPCMessage(mm)
		}

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID)] = func() {
			detail := mock.EmptyDeferredAwardDetail(logic.SlotLogicBase, user)
			mm := mock.MarshalDeferredAward(detail)
			logic.OnRecvRPCMessage(mm)
		}

		handlers1[int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID)] = func() {
			mm := mock.EmptyUserAssignedModeRPCMessage(logic.SlotLogicBase, user)
			logic.OnRecvRPCMessage(mm)
		}

		slot_test.PrintSendRPCMessgeFunc(t, safe, m, handlers1)
	}

	c := make(chan bool)
	logic.SendMessgeFunc = func(userID int32, m *msg.Message) {
		slot_test.PrintSendMessgeFunc(t, userID, m)

		if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_StatusID) {
			data := m.MsgData.(*netproto.Slot_MSG_Status)
			if data.GetState() == 5 {
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
