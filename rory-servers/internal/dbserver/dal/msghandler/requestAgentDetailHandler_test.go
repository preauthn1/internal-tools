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



func TestRequestAgentDetailHandler(t *testing.T) {
	initDBConnection()

	serverid := int32(999)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestAgentListID), rpcMessage)
	handler := new(RequestAgentDetailHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_RequestAgentListRetID) {
		t.Error("OnMessage 錯誤")
	}
}

func TestJoinAgentHandler(t *testing.T) {
	initDBConnection()

	serverid := int32(999)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	data := new(netproto.JoinAgentList)
	data.UserID = proto.Int32(0)
	data.InviteCode = proto.String("1qaz2wsx")
	data.RetCode = proto.Int32(0)
	data.RetMsg = proto.String("")
	msgData, _ := proto.Marshal(data)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_JoinAgentListID), rpcMessage)
	handler := new(RequestAgentDetailHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_JoinAgentListRetID) {
		t.Error("OnMessage 錯誤")
	}
}

func TestApplyAgentListHandler(t *testing.T) {
	initDBConnection()

	serverid := int32(999)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	data := new(netproto.RequestApplyAgent)
	data.UserID = proto.Int32(0)
	data.Tel = proto.String("1qaz2wsx")
	data.Weixin = proto.String("1qaz2wsx")
	data.QQ = proto.String("")
	data.Email = proto.String("")
	data.VCode = proto.String("1qaz2wsx")
	msgData, _ := proto.Marshal(data)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentListID), rpcMessage)
	handler := new(RequestAgentDetailHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ApplyAgentListRetID) {
		t.Error("OnMessage 錯誤")
	}
}

func TestCheckAgentHandler(t *testing.T) {
	initDBConnection()

	serverid := int32(999)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CheckAgentID), rpcMessage)
	handler := new(RequestAgentDetailHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_CheckAgentRetID) {
		t.Error("OnMessage 錯誤")
	}
}
