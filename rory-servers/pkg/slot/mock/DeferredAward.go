package mock

import (
	"strconv"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/baselogic"
)

// EmptyDeferredAwardRPCMessage 空的 DeferredAward
func EmptyDeferredAwardRPCMessage(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *msg.Message {

	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)

	detail := new(netproto.Slot_MSG_DeferredAwardRecord)
	detail.UserID = proto.Int32(user.UserID)
	detail.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	detail.FlagID = proto.Int32(int32(f))

	b, _ := proto.Marshal(detail)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID)
	mm.MsgData = _m

	return mm
}

// EmptyDeferredAwardDetail 空的 DeferredAward Detail
func EmptyDeferredAwardDetail(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *netproto.Slot_MSG_DeferredAwardRecord {

	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)

	detail := new(netproto.Slot_MSG_DeferredAwardRecord)

	detail.UserID = proto.Int32(user.UserID)
	detail.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	detail.FlagID = proto.Int32(int32(f))

	return detail
}

// MarshalDeferredAward 包裝成 *msg.Message
func MarshalDeferredAward(detail *netproto.Slot_MSG_DeferredAwardRecord) *msg.Message {
	b, _ := proto.Marshal(detail)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID)
	mm.MsgData = _m

	return mm
}
