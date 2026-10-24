package msghandler

import (
	proto "example.invalid/golang/protobuf/proto"
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
	"fmt"
)

func TestFindSetBankPwdByPhoneHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *FindSetBankPwdByPhoneHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &FindSetBankPwdByPhoneHandler{},
			args: args{
				bclsid: 2,
				sclsid: 1019,
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			rh:   &FindSetBankPwdByPhoneHandler{},
			args: args{
				bclsid: 3,
				sclsid: 1019,
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			rh:   &FindSetBankPwdByPhoneHandler{},
			args: args{
				bclsid: 2,
				sclsid: 1039,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &FindSetBankPwdByPhoneHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("FindSetBankPwdByPhoneHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindSetBankPwdByPhoneHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(77874156),
	}

	msgData, _ := proto.Marshal(&netproto.FindSetBankPwdByPhone{
		Area_ID:  proto.Int32(1),
		Password: [REDACTED]"pwd"),
		VCode:    proto.String("vcode"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Hall),
		SClassID: int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneID),
		MsgData:  rpcMessage,
	}

	handler := new(FindSetBankPwdByPhoneHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.RetMessage)

	if !ok {
		t.Errorf("Return Message is not *netproto.RetMessage")
	}

	fmt.Printf("Code: %d, Message: %s\n", msg.GetCode(), msg.GetMessage())

}
