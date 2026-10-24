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

//全民推廣員
func TestGetUserAngentExecInfoHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.GetUserAgentExecMsg)
	msgdata.UserID = proto.Int32(1)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentExecID), rpcMessage)
	handler := new(GetUserAngentExecInfoHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetUserAgentExecRetID) {
		t.Error("OnMessage 錯誤")
	}

	/*
		mlog.System("查看回傳值的結果 AgentTaxBonus: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentTaxBonus())
		mlog.System("查看回傳值的結果 AgentLevel: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentLevel())
		mlog.System("查看回傳值的結果 AgentUserAmount: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentUserAmount())
		mlog.System("查看回傳值的結果 AgentAmount: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentAmount())
		mlog.System("查看回傳值的結果 AgentTotalUser: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentTotalUser())
		mlog.System("查看回傳值的結果 AgentTaxAmount: %v", ret[0].MsgData.(*netproto.GetUserAgentExecRet).GetAgentTaxAmount())
	*/
}
