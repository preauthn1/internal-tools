package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/stretchr/testify/assert"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"testing"
)

func TestLoadBlockIPListHandler(t *testing.T) {
	InitTest()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LoadBlockIPID), rpcMessage)
	handler := new(LoadBlockIPListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_BlockIPListID) {
		t.Error("HallServerStart OnMessage 錯誤")
	}

	assert.NotNil(t, ret[0].MsgData)
}
