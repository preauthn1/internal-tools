package msghandler

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"testing"
)



func TestSuperToolFishUpdSystemRewardHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.SuperToolFishUpdSystemRewardDetail)
	msgdata.FlagID = proto.Int32(1)
	fishRewardList := make([]*netproto.SuperToolFishUpdSystemRewardDataList, 0)
	fishReward := new(netproto.SuperToolFishUpdSystemRewardDataList)
	fishReward.FishID = proto.Int32(1)
	fishReward.FinishAmount = proto.Int64(100)
	fishRewardList = append(fishRewardList, fishReward)
	msgdata.FishRewardList = fishRewardList
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardID), rpcMessage)
	handler := new(SuperToolFishUpdSystemRewardHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) || ret[0].SClassID != int32(netproto.GameMessageClassID_SuperToolFishUpdSystemRewardRetID) {
		t.Error("OnMessage 錯誤")
	}
}
