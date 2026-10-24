package mock

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
)

// EmptyBonusGameDetail 空的 BonusGame Detail
func EmptyBonusGameDetail(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *netproto.Slot_MSG_BonusGameRecord {

	detail := new(netproto.Slot_MSG_BonusGameRecord)

	return detail
}

// MarshalBonusGame 包裝成 *msg.Message
func MarshalBonusGame(detail *netproto.Slot_MSG_BonusGameRecord) *msg.Message {
	b, _ := proto.Marshal(detail)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID)
	mm.MsgData = _m

	return mm
}
