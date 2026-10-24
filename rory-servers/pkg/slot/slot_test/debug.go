package slot_test

import (
	"os"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
	"example.invalid/ServerTeam/Servers/pkg/slot/helper"
)

// PrintLogic 印出Logic內參數
func PrintLogic(t *testing.T, logic *baselogic.SlotLogicBase) {
	helper.PrintLogicInfo(os.Stdout, logic)
	helper.PrintLogicProto(os.Stdout, logic.Proto, true)
}

// PrintSendMessgeFunc print information for SendMessgeFunc
func PrintSendMessgeFunc(t *testing.T, userID int32, m *msg.Message, skipIDs ...int32) {
	if !conf.PrintMessage {
		return
	}

	skipIDmap := make(map[int32]bool)
	for _, ID := range skipIDs {
		skipIDmap[ID] = true
	}

	if !skipIDmap[m.SClassID] {
		helper.PrintMessage(os.Stdout, m)
	}
}

// PrintSendRPCMessgeFunc print information for SendRPCMessgeFunc
func PrintSendRPCMessgeFunc(t *testing.T, safe bool, m *msg.Message, handlers map[int32]func()) {

	if handlers != nil {
		if f, ok := handlers[m.SClassID]; ok {
			f()
		}
	}

	if !conf.PrintMessage {
		return
	}

	t.Log("SendRPCMessgeFunc", safe, m)

	return

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID) {
		t.Log("[Read JackpotNoLink]")

		data := m.MsgData.(*netproto.Slot_MSG_JackpotRead)

		t.Log(" - GameID:", example.invalid.GetGameID())
		t.Log(" - FlagID:", example.invalid.GetFlagID())
		t.Log(" - TableNo:", example.invalid.GetTableNo())
		t.Log(" - Level:", example.invalid.GetLevel())
		t.Log(" - Contribution:", example.invalid.GetContribution())
		t.Log(" - TotalContribution:", example.invalid.GetTotalContribution())
		t.Log(" - AwardTotal:", example.invalid.GetAwardTotal())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID) {
		t.Log("[Record JackpotNoLink]")

		data := m.MsgData.(*netproto.Slot_MSG_JackpotRecord)

		t.Log(" - Records:")

		for _, row := range data.GetRecords() {
			t.Log("   - GameID:", row.GetGameID())
			t.Log("   - FlagID:", row.GetFlagID())
			t.Log("   - TableNo:", row.GetTableNo())
			t.Log("   - Level:", row.GetLevel())
			t.Log("   - Contribution:", row.GetContribution())
			t.Log("   - TotalContribution:", row.GetTotalContribution())
			t.Log("   - AwardTotal:", row.GetAwardTotal())
		}

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID) {
		t.Log("[Read FreeSpinTimes]")

		data := m.MsgData.(*netproto.Slot_MSG_FreeGameRead)

		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - UserID:", data.GetUserID())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_FreeGameRecordID) {
		t.Log("[Record FreeSpin]")

		data := m.MsgData.(*netproto.Slot_MSG_FreeGameRecord)

		t.Log(" - GameID:", example.invalid.GetGameID())
		t.Log(" - FlagID:", example.invalid.GetFlagID())
		t.Log(" - UserID:", example.invalid.GetUserID())
		t.Log(" - UserBet:", example.invalid.GetUserBet())
		t.Log(" - ChooseRate:", example.invalid.GetChooseRate())
		t.Log(" - FreeSpinGroup:", example.invalid.GetFreeSpinGroup())
		t.Log(" - Times:", example.invalid.GetTimes())
		t.Log(" - WinMoney:", example.invalid.GetWinMoney())
		// t.Log(" - IsCoupon:", example.invalid.GetIsCoupon())
		t.Log(" - Data:", example.invalid.GetData())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID) {
		t.Log("[Read BonusGame]")

		data := m.MsgData.(*netproto.Slot_MSG_BonusGameRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID) {
		t.Log("[Record BonusGame]")

		data := m.MsgData.(*netproto.Slot_MSG_BonusGameRecord)

		for _, row := range data.Detail {
			t.Logf(" - Detail - %v\n", row)
		}

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID) {
		t.Log("[Reset BonusGame]")

		data := m.MsgData.(*netproto.Slot_MSG_BonusGameRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID) {

	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID) {
		data := m.MsgData.(*netproto.Slot_MSG_GameLogAdd)

		t.Log("[Add GameLog]")

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - ServerID:", data.GetServerID())
		t.Log(" - GameType:", data.GetGameType())
		t.Log(" - LogType:", data.GetLogType())
		t.Log(" - CurrentAmount:", data.GetCurrentAmount())
		t.Log(" - BetMoney:", data.GetBetMoney())
		t.Log(" - WinMoney:", data.GetWinMoney())
		t.Log(" - FinalScore:", data.GetFinalScore())
		t.Log(" - GameData:", data.GetGameData())
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID) {
		data := m.MsgData.(*netproto.Slot_MSG_BatchGameLogAdd)

		t.Log("[Batch Add GameLog]")

		for _, log := range data.GetLogs() {
			t.Log(" - [GameLog]")
			t.Log("   - UserID:", log.GetUserID())
			t.Log("   - GameID:", log.GetGameID())
			t.Log("   - ServerID:", log.GetServerID())
			t.Log("   - GameType:", log.GetGameType())
			t.Log("   - LogType:", log.GetLogType())
			t.Log("   - CurrentAmount:", log.GetCurrentAmount())
			t.Log("   - BetMoney:", log.GetBetMoney())
			t.Log("   - WinMoney:", log.GetWinMoney())
			t.Log("   - FinalScore:", log.GetFinalScore())
			t.Log("   - GameData:", log.GetGameData())
		}

	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID) {
		t.Log("[Read SelectGame]")

		data := m.MsgData.(*netproto.Slot_MSG_SelectGameRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		return
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

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID) {
		t.Log("[Read LightningGame]")

		data := m.MsgData.(*netproto.Slot_MSG_LightningRead)

		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - UserID:", data.GetUserID())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID) {
		t.Log("[Rcord LightningGame]")

		data := m.MsgData.(*netproto.Slot_MSG_LightningRecord)

		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - Data:", data.GetData())

		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID) {
		t.Log("[Read DeferredAward]")

		data := m.MsgData.(*netproto.Slot_MSG_DeferredAwardRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - AwardType:", data.GetAwardType())
		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID) {
		t.Log("[Record DeferredAward]")

		data := m.MsgData.(*netproto.Slot_MSG_DeferredAwardRecord)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Logf(" - Records(%d):\n", len(data.GetRecords()))
		for _, row := range data.GetRecords() {
			t.Log("   - ", row)
		}
		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID) {
		t.Log("[Pop DeferredAward]")

		data := m.MsgData.(*netproto.Slot_MSG_DeferredAwardRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - AwardType:", data.GetAwardType())
		return
	}

	if m.SClassID == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID) {
		t.Log("[Reset DeferredAward]")

		data := m.MsgData.(*netproto.Slot_MSG_DeferredAwardRead)

		t.Log(" - UserID:", data.GetUserID())
		t.Log(" - GameID:", data.GetGameID())
		t.Log(" - FlagID:", data.GetFlagID())
		t.Log(" - AwardType:", data.GetAwardType())
		return
	}
}
