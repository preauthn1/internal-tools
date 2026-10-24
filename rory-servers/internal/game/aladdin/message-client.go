package aladdin

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// SendGameResultMessage 發送遊戲結果
//
//   * BClassID = 10001 (MessageBClassID_Game)
//   * SClassID = 20004 (Slot_GameMessageClassID_Slot_GameResultID)
//   * Message = { BetMoney: 100, WinMoney: 0, FinalMoney: 300, ReelResult: [{},{},...] }
func (logic *Logic) SendGameResultMessage(user *gamelogic.GameUser) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameResultID)
	data := new(netproto.Slot_MSG_StepGameResult)

	var currentMoney int64

	data.BetMoney = proto.Int64(logic.GamePayMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney)
	data.FinalMoney = proto.Int64(logic.GameUser.Score)

	currentMoney = logic.GameUser.Score - logic.GameWinMoney

	for _, step := range logic.StepResults {
		currentMoney += step.GameWinMoney
		result := new(netproto.Slot_StepGameResult)

		result.IsRespin = proto.Bool(step.IsRespin)
		result.IsExpand = proto.Bool(step.IsExpand)
		result.ExpandReel = step.ExpandReel
		result.ExpandOrder = proto.String(string(step.ExpandOrder))

		result.InitResult = logic.BuildReelResult(step.UnexpandResult, step.GamePayResult)

		result.SlotResult = new(netproto.Slot_MSG_GameResult)

		result.SlotResult.BetMoney = proto.Int64(step.GamePayMoney)
		result.SlotResult.WinMoney = proto.Int64(step.GameWinMoney)
		result.SlotResult.FinalMoney = proto.Int64(currentMoney)
		result.SlotResult.ReelResult = logic.BuildReelResult(step.ExpandResult, step.GamePayResult)
		result.SlotResult.PayGroups = logic.BuildPayGroups(step.GamePayResult)

		data.SlotResult = append(data.SlotResult, result)
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}
