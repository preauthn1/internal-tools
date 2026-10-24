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



func TestSuperToolFishUpdVIPUserHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.SuperToolfishUpdVIPUserRewardDetail)
	msgdata.GameID = proto.Int32(0)
	dataList := make([]*netproto.SuperToolFishUpdVIPUserRewardDataList, 0)
	data := new(netproto.SuperToolFishUpdVIPUserRewardDataList)
	data.UserID = proto.Int32(1021470)
	data.FinishAmount = proto.Int64(100)
	dataList = append(dataList, data)
	msgdata.FishUserRewardList = dataList
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardID), rpcMessage)
	handler := new(SuperToolFishUpdVIPUserHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) || ret[0].SClassID != int32(netproto.GameMessageClassID_SuperToolFishUpdVIPUserRewardRetID) {
		t.Error("OnMessage 錯誤")
	}
}
