package msghandler

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestSlotJackpotLinkHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.Slot_JackpotData)
	msgdata.GameID = proto.Int32(25)
	msgdata.FlagID = proto.Int32(2)
	msgdata.TableNo = proto.Int32(1)
	msgdata.Contribution = proto.Int64(0)
	msgdata.TotalContribution = proto.Int64(0)
	msgdata.AwardTotal = proto.Int64(0)
	msgdata.Level = proto.Int32(1)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SlotJackpotLinkID), rpcMessage)
	handler := new(SlotJackpotLinkHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) || ret[0].SClassID != int32(netproto.GameMessageClassID_SlotJackpotLinkRetID) {
		t.Error("OnMessage 錯誤")
	}

	fmt.Println(ret[0].MsgData)
}
