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

func TestGetContestUserDataHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GetContestUserDataHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &GetContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_GetContestUserDataID),
			},
			want: true,
		},
		{
			name: "Not Hooked BClass",
			rh:   &GetContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Common),
				sclsid: int32(netproto.DBServerClassID_GetContestUserDataID),
			},
			want: false,
		},
		{
			name: "Not Hooked SClass",
			rh:   &GetContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GetContestUserDataHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GetContestUserDataHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetContestUserDataHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(77874156),
	}

	msgData, _ := proto.Marshal(&netproto.GetConestUser{
		GameID: proto.Int32(14),
		UserID: proto.Int32(89342501),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.DBServerClassID_GetContestUserDataID),
		MsgData:  rpcMessage,
	}

	handler := new(GetContestUserDataHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.DBServerClassID_GetContestUserDataRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.DBServerClassID_GetContestUserDataRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.UpdateContestUser)

	if !ok {
		t.Errorf("Return Message is not *netproto.UpdateContestUser")
	}

	fmt.Printf("%+v\n", msg)
}
