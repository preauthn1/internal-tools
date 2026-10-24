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

func TestGetWithdrawListIHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	msgdata := new(netproto.GetWithdrawMsg)
	msgdata.UserID = proto.Int32(1)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawListID), rpcMessage)
	handler := new(GetWithdrawListIHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	handler.OnMessage(ctx, clt, m)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetWithdrawListRetID) {
		t.Error("OnMessage 錯誤")
	}

	/*
		mlog.System("檢查回傳資料 PerformedAmount %v", ret[0].MsgData.(*netproto.GetWithdrawRet).GetTotalTaxAmount())
		mlog.System("檢查回傳資料 ConvertAmount %v", ret[0].MsgData.(*netproto.GetWithdrawRet).GetConvertAmount())
		mlog.System("檢查回傳資料 PerformedAmount %v", ret[0].MsgData.(*netproto.GetWithdrawRet).GetPerformedAmount())
	*/
}
