package s5ll

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestNewLogic(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	slot_test.PrintLogic(t, logic.SlotLogicBase)
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

	// printSendRPCMessgeFunc(t, safe, m, handlers2)
}
