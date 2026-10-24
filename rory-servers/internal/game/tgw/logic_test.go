package tgw

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/mock"
	"example.invalid/ServerTeam/Servers/pkg/slot/slot_test"
)

func TestTGWNewLogic(t *testing.T) {

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)

	logic := grs.GameLogicManager.GetGameLogic(0).(*Logic)

	slot_test.PrintLogic(t, logic.SlotLogicBase)
}

func TestTGWLogic_UserEnter(t *testing.T) {

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
		slot_test.PrintSendMessgeFunc(t, userID, m, int32(netproto.DFDC_GameMessageClassID_DFDCSettingID))
		printSendMessgeFunc(t, userID, m)

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
}

func printSendMessgeFunc(t *testing.T, userID int32, m *msg.Message) {

	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_SettingID) {
		data := m.MsgData.(*netproto.TGW_MSG_Setting)

		t.Logf("[Game Setting]\n")

		t.Logf(" - BetCount: %d\n", data.GetBetCount())
		t.Logf(" - RowCounts: %d\n", data.GetRowCounts())
		t.Logf(" - TotalPayLine: %d\n", data.GetTotalPayLine())
		t.Logf(" - BigWinStep: %d\n", data.GetBigWinStep())
		t.Logf(" - PayLineMode: %d\n", data.GetPayLineMode())

		t.Log(" - ReelSetting:")
		for _, row := range data.GetReelSetting() {
			t.Logf("   - %v\n", row)
		}
		t.Log(" - SymbolPayTable:")
		for _, row := range data.GetSymbolPayTable() {
			t.Logf("   - %v\n", row)
		}

		t.Log(" - UnitPrefix:")
		for _, row := range data.GetUnitPrefix() {
			t.Logf("   - %v\n", row)
		}
	}

	if m.SClassID == int32(netproto.DFDC_GameMessageClassID_DFDCStateID) {
		data := m.MsgData.(*netproto.Slot_MSG_Status)

		t.Logf("[Status] State: %d\n", data.GetState())
	}

	if m.SClassID == int32(netproto.TGW_GameMessageClassID_TGW_BonusGameID) {
		data := m.MsgData.(*netproto.TGW_MSG_BonusGame)

		t.Log("[BonusGame]")
		t.Logf("  - WinMoney: %d\n", data.GetWinMoney())
		t.Logf("  - Count: %d\n", data.GetCount())
		t.Logf("  - All: %v\n", data.GetAll())
		t.Logf("  - Records: \n")
		for _, row := range data.GetRecords() {
			t.Logf("    - %v\n", row)
		}
	}
}

func printSendRPCMessgeFunc(t *testing.T, safe bool, m *msg.Message, handlers map[int32]func()) {
	if handlers != nil {
		if f, ok := handlers[m.SClassID]; ok {
			f()
		}
	}
}

func userEnter(logic *Logic, user *gamelogic.GameUser, t *testing.T, safe bool, m *msg.Message) {
	handlers1 := make(map[int32]func())

	handlers1[int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID)] = func() {
		mm := mock.EmptyJackpotPoolRPCMessage(logic.SlotLogicBase, user)
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
