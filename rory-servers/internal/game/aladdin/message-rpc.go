package aladdin

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// AddSlotGameHistoryLogByStep 新增歷程紀錄
func (logic *Logic) AddSlotGameHistoryLogByStep(user *gamelogic.GameUser, step *StepResult) {

	data := new(netproto.Slot_MSG_GameLogAdd)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	data.ServerID = proto.Int32(logic.TableManager.Room.ServerID)

	data.CurrentAmount = proto.Int64(user.Score)
	data.BetMoney = proto.Int64(step.GamePayMoney)
	data.WinMoney = proto.Int64(step.GameWinMoney)
	data.FinalScore = proto.Int64(user.Score + step.GameWinMoney - step.GamePayMoney)

	var historyType int32
	if !step.IsRespin {
		historyType = slot.HistoryTypeMAIN
	} else {
		historyType = slot.HistoryTypeRESPIN
	}

	data.GameType = proto.Int32(historyType)
	data.LogType = proto.Int32(slot.LogTypeREEL)

	step.ExpandResult.ReelMode()
	data.GameData = proto.String(step.ExpandResult.ToSymbolString())

	logic.AddSlotGameHistoryLogByMsg(user, data)
}
