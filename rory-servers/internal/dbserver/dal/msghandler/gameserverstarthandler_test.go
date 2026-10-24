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
	"fmt"
	"testing"
)

func TestGameServerStartHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GameServerStartHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &GameServerStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_DBServer),
				sclsid: int32(netproto.DBServerClassID_GameServerStartID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			rh:   &GameServerStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_GameServerStartID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			rh:   &GameServerStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_DBServer),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GameServerStartHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GameServerStartHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameServerStartHandler_OnMessage_IntegrationTest(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.GameServerStart{
		ServerID:   proto.Int32(0),
		ServerName: proto.String(""),
		GameID:     proto.Int32(0),
		ShowName:   proto.String(""),
		LoginMoney: proto.Int32(0),
		HttpAddr:   proto.String(""),
		ClientAddr: proto.String(""),
		NotifyAddr: proto.String(""),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_DBServer),
		SClassID: int32(netproto.DBServerClassID_GameServerStartID),
		MsgData:  rpcMessage,
	}

	handler := new(GameServerStartHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_DBServer), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.DBServerClassID_GameServerStartRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.DBServerClassID_GameServerStartRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.GameServerStartRet)

	if !ok {
		t.Errorf("Return Message is not *netproto.GameServerStartRet")
	}

	fmt.Printf("Code: %d, Message: %s, AndroidVersion: %s, IOSVersion: %s, Status: %d, Status Message: %s\n", msg.GetCode(), msg.GetMessage(), msg.GetAndroidVersion(), msg.GetIOSVersion(), msg.GetStatus(), msg.GetStatusMessage())
}
