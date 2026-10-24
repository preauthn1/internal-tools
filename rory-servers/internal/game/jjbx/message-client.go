package jjbx

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
	data := new(netproto.Slot_MSG_GameResult)

	data.BetMoney = proto.Int64(logic.GamePayMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney - logic.GameWinJackpot)
	data.JackpotMoney = proto.Int64(logic.GameWinJackpot)
	data.FinalMoney = proto.Int64(logic.GameUser.Score)
	data.ReelResult = logic.BuildReelResult(logic.GameResultManager, logic.GamePayResult)
	data.PayGroups = logic.BuildPayGroups(logic.GamePayResult)
	if logic.FreeSpin != nil {
		data.FreeGameWinMoney = proto.Int64(logic.FreeSpin.WinMoney)

		if logic.FreeSpin.FreeSpinExtAward != nil {

			if logic.FreeSpin.FreeSpinExtAward.IsFinished {
				data.WinMoney = proto.Int64(logic.GameWinMoney - logic.GameWinJackpot - logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)
			}

			data.FreeGameExtAward = new(netproto.Slot_FreeGameExtAward)
			data.FreeGameExtAward.TotalCount = proto.Int32(logic.FreeSpin.FreeSpinExtAward.TotalCount)
			data.FreeGameExtAward.TotalWinMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.TotalWinMoney)
			data.FreeGameExtAward.Count = proto.Int32(logic.FreeSpin.FreeSpinExtAward.Count)
			data.FreeGameExtAward.WinMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.WinMoney)
			data.FreeGameExtAward.EachCountMoney = proto.Int64(logic.FreeSpin.FreeSpinExtAward.EachCountMoney)
			data.FreeGameExtAward.IsFinished = proto.Bool(logic.FreeSpin.FreeSpinExtAward.IsFinished)

			for pos, symbol := range logic.FreeSpin.FreeSpinExtAward.Symbols {
				s := new(netproto.Slot_Symbol)
				s.SymbolID = proto.Int32(int32(symbol.GetId()))
				s.SymbolName = proto.String(symbol.GetSymbol())
				s.Position = &netproto.Slot_ReelPosition{
					Reel: proto.Int32(pos.Reel),
					Row:  proto.Int32(pos.Row),
				}
				data.FreeGameExtAward.Symbols = append(data.FreeGameExtAward.Symbols, s)
			}
		}
	}

	pmsg.MsgData = data

	logic.SendMessageToUser(user, pmsg)
}
