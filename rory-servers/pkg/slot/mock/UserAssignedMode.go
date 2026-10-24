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

// EmptyUserAssignedModeRPCMessage 空的 UserAssignedMode
func EmptyUserAssignedModeRPCMessage(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *msg.Message {

	data := new(netproto.Slot_MSG_UserAssignedModeRecord)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
	data.FlagID = proto.Int32(int32(f))

	b, _ := proto.Marshal(data)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadRetID)
	mm.MsgData = _m

	return mm
}
