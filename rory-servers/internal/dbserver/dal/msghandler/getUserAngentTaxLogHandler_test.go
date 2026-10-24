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

//取得全民推廣TaxLog
func TestGetUserAngentTaxLogHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.GetUserAgentTaxLogMsg)
	msgdata.UserID = proto.Int32(10)
	msgdata.DateTime = proto.String("2019/04/13")
	// msgdata.Area_ID = proto.Int32(0)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentTaxLogID), rpcMessage)
	handler := new(GetUserAngentTaxLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	handler.OnMessage(ctx, clt, m)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetUserAgentTaxLogRetID) {
		t.Error("OnMessage 錯誤")
	}

	/*
		mlog.System("查看回傳值的結果 GetTotalTax: %v", ret[0].MsgData.(*netproto.GetUserAgentTaxLogRetList).GetTotalTax())
		mlog.System("查看回傳值的結果 GetTotalPeople: %v", ret[0].MsgData.(*netproto.GetUserAgentTaxLogRetList).GetTotalPeople())
		mlog.System("查看回傳值的結果 TotalTaxBonux: %v", ret[0].MsgData.(*netproto.GetUserAgentTaxLogRetList).GetTotalTaxBonux())
		mlog.System("查看回傳值的結果 UserAgentTaxLog: %v", ret[0].MsgData.(*netproto.GetUserAgentTaxLogRetList).GetUserAgentTaxLog())
	*/
}
