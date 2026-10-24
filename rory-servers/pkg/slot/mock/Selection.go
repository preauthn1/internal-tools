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

// EmptySelectionRPCMessage 空的 Selection
func EmptySelectionRPCMessage(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *msg.Message {

	detail := new(netproto.Slot_MSG_SelectGameRecord)
	example.invalid = new(netproto.Slot_SelectGameInfo)
	example.invalid.UserID = proto.Int32(user.UserID)
	example.invalid.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
	example.invalid.FlagID = proto.Int32(int32(f))
	example.invalid.Times = proto.Int32(0)

	b, _ := proto.Marshal(detail)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadRetID)
	mm.MsgData = _m

	return mm
}
