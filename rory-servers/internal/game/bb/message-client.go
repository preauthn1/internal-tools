package bb

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
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
	data := new(netproto.Slot_MSG_GameResult)

	data.BetMoney = proto.Int64(logic.GamePayMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney - logic.GameWinJackpot)
	data.JackpotMoney = proto.Int64(logic.GameWinJackpot)
	data.FinalMoney = proto.Int64(logic.GameUser.Score)
	data.ReelResult = logic.BuildReelResult(logic.GameResultManager, logic.GamePayResult)
	data.PayGroups = logic.BuildPayGroups(logic.GamePayResult)
	if logic.FreeSpin != nil {
		data.FreeGameWinMoney = proto.Int64(logic.FreeSpin.WinMoney)
	}

	{ // 整理加倍符號
		if logic.GameMultiDetail != nil {
			for _, s := range data.ReelResult {
				if m, ok := logic.GameMultiDetail[slot.Position{Reel: s.Position.GetReel(), Row: s.Position.GetRow()}]; ok {
					s.Symbol.SymbolName = proto.String(fmt.Sprintf("w%d", m))
				}
			}
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}
