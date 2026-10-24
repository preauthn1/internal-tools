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

func TestFishSysPoolHandler_IsHook_UnitTest(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		fsp  *FishSysPoolHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			fsp:  &FishSysPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_FishSysPoolID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			fsp:  &FishSysPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.GameMessageClassID_FishSysPoolID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			fsp:  &FishSysPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsp := &FishSysPoolHandler{}
			if got := fsp.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("FishSysPoolHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFishSysPoolHandler_OnMessage_IntegrationTest(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.FishSysPoolData{
		FlagID:        proto.Int32(1),
		SumBetMoney:   proto.Int64(10),
		SumBullet:     proto.Int64(5),
		PoolAmount:    proto.Int64(100),
		SysSumWinLoss: proto.Int64(0),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.GameMessageClassID_FishSysPoolID),
		MsgData:  rpcMessage,
	}

	handler := new(FishSysPoolHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.GameMessageClassID_FishSysPoolRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.GameMessageClassID_FishSysPoolRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.FishSysPoolData)

	if !ok {
		t.Errorf("Return Message is not *netproto.FishSysPoolData")
	}

	fmt.Printf("Data: %+v\n", msg)

}
