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



func TestUpdateContestUserDataHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *UpdateContestUserDataHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &UpdateContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_UpdateContestUserDataID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			rh:   &UpdateContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.DBServerClassID_UpdateContestUserDataID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			rh:   &UpdateContestUserDataHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &UpdateContestUserDataHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("UpdateContestUserDataHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateContestUserDataHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7556530),
	}

	msgData, _ := proto.Marshal(&netproto.UpdateContestUser{
		GameID:      proto.Int32(14),
		UserID:      proto.Int32(7556530),
		ContestID:   proto.Int32(5),
		ContestNo:   proto.Int32(0),
		FlagID:      proto.Int32(0),
		TableNo:     proto.Int32(0),
		PlayTime:    proto.Int32(0),
		RemainScore: proto.Int32(0),
		Score:       proto.Int32(0),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.DBServerClassID_UpdateContestUserDataID),
		MsgData:  rpcMessage,
	}

	handler := new(UpdateContestUserDataHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_DBServer), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.DBServerClassID_DBRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.DBServerRet)

	if !ok {
		t.Errorf("Return Message is not *netproto.DBServerRet")
	}

	fmt.Printf("Code: %d\n", msg.GetCode())
}
