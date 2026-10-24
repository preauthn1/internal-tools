package msghandler

import (
	"example.invalid/golang/protobuf/proto"
)

import (	
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"	
)

import (
	"sync"
	"testing"
)

func testRunRange (t *testing.T) {
	var wg sync.WaitGroup

	for i:=0; i<30000; i++ {
		wg.Add(1)
		go func() {
			testGuestLogonHandler(t)
			wg.Done()
		}()
	}

	wg.Wait()
	mlog.System("Done")
}


func testGuestLogonHandler(t *testing.T) {
	initDBConnection()

	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.GuestLogin)
	msgdata.ServerID = proto.Int32(2)
	msgdata.HDCode = proto.String("weihao23")
	msgdata.HDType = proto.Int32(2)
	msgdata.SiteID = proto.Int32(2)
	msgdata.Version = proto.String("123")
	msgdata.InviteCode = proto.String("123")
	msgdata.PlatformID = proto.Int32(2)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GuestLoginID), rpcMessage)
	handler := new(GuestLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_LoginRetID) {
		t.Error("LoginRetID OnMessage 錯誤")
	}
}

