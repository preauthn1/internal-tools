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



// 申請提現例表(全民代理)
func TestUserAgentGetWithAmountListHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.UserAgentGetWithAmountListMsg)
	msgdata.UserID = proto.Int32(36613013)
	msgdata.TargetID = proto.Int32(0)
	msgdata.PageIndex = proto.Int32(1)
	msgdata.PageSize = proto.Int32(7)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountListID), rpcMessage)
	handler := new(UserAgentGetWithAmountListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	handler.OnMessage(ctx, clt, m)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentGetWithAmountListRetID) {
		t.Error("OnMessage 錯誤")
	}

	// t.Log(ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).UserAgentGetWithAmountList[0].Userid)

}
