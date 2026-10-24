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

// EmptyJackpotPoolRPCMessage 空的 Jackpoot Pool
func EmptyJackpotPoolRPCMessage(logic *baselogic.SlotLogicBase, user *gamelogic.GameUser) *msg.Message {

	detail := new(netproto.Slot_MSG_JackpotRecord)

	jp1 := new(netproto.Slot_JackpotData)
	jp1.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	f, _ := strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
	jp1.FlagID = proto.Int32(int32(f))
	jp1.TableNo = proto.Int32(logic.TableManager.TableNo)
	jp1.Contribution = proto.Int64(0)
	jp1.TotalContribution = proto.Int64(0)
	jp1.AwardTotal = proto.Int64(0)
	jp1.Level = proto.Int32(1)

	jp2 := new(netproto.Slot_JackpotData)
	jp2.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	f, _ = strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
	jp2.FlagID = proto.Int32(int32(f))
	jp2.TableNo = proto.Int32(logic.TableManager.TableNo)
	jp2.Contribution = proto.Int64(0)
	jp2.TotalContribution = proto.Int64(0)
	jp2.AwardTotal = proto.Int64(0)
	jp2.Level = proto.Int32(2)

	jp3 := new(netproto.Slot_JackpotData)
	jp3.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	f, _ = strconv.ParseFloat(logic.TableManager.Room.Flag, 64)
	jp3.FlagID = proto.Int32(int32(f))
	jp3.TableNo = proto.Int32(logic.TableManager.TableNo)
	jp3.Contribution = proto.Int64(0)
	jp3.TotalContribution = proto.Int64(0)
	jp3.AwardTotal = proto.Int64(0)
	jp3.Level = proto.Int32(3)

	detail.Records = append(detail.Records, jp1)
	detail.Records = append(detail.Records, jp2)
	detail.Records = append(detail.Records, jp3)

	b, _ := proto.Marshal(detail)

	_m := new(rpcmsg.RPCMessage)
	_m.MsgData = b

	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Game)
	mm.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID)
	mm.MsgData = _m

	return mm
}
