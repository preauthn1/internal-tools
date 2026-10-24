package logic

import (
	"example.invalid/gogo/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

func (logic *Logic) srvSpin(money float64, rate int32) {
	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionSpinID)

	_data := new(netproto.Slot_MSG_Spin)
	_data.Money = proto.Int64(int64(logic.MoneyTransfer.C2S(money)))

	if rate > 0 {
		_data.Rate = proto.Int32(rate)
	}

	_m.MsgData = _data

	logic.sw.WriteServer(_m)

	logic.Round++
	logic.ServeState = ServeStateWaitResult
	logic.lastspin = _data

	logic.clean()
}

func (logic *Logic) srvPickGameItem(id, loc int32, money int64) {
	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

	_data := new(netproto.Slot_MSG_PickBonus)
	_data.Type = proto.String(string(slot.GameModePickGame))

	_data.Record = &netproto.Slot_BonusGameRecord{
		ID:          proto.Int32(id),
		Location:    proto.Int32(loc),
		RewareValue: proto.Int64(money),
	}

	_m.MsgData = _data

	logic.sw.WriteServer(_m)
}

func (logic *Logic) srvJackpotGameItem(id, loc int32, level int64) {
	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_PickBonusID)

	_data := new(netproto.Slot_MSG_PickBonus)
	_data.Type = proto.String(string(slot.GameModeJackpotGame))

	_data.Record = &netproto.Slot_BonusGameRecord{
		ID:          proto.Int32(id),
		Location:    proto.Int32(loc),
		RewareValue: proto.Int64(level),
	}

	_m.MsgData = _data

	logic.sw.WriteServer(_m)
}

func (logic *Logic) srvJackpotPool(ids ...int32) {
	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotPoolID)

	_data := new(netproto.Slot_MSG_GetJackpot)
	_data.ID = ids

	_m.MsgData = _data

	logic.sw.WriteServer(_m)
}

func (logic *Logic) srvSelectionChoose(t int) {

	if logic.SelectGameOptions == nil {
		return
	}

	opt := logic.SelectGameOptions.Options[t]

	_m := new(msg.Message)
	_m.BClassID = int32(netproto.MessageBClassID_Game)
	_m.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_ActionChooseGameID)

	_data := new(netproto.Slot_MSG_ChooseGame)
	_data.Option = opt

	_m.MsgData = _data

	logic.sw.WriteServer(_m)
}
