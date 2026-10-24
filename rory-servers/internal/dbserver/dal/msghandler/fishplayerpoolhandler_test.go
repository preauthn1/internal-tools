package msghandler

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestFishPlayerPoolHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		fpp  *FishPlayerPoolHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			fpp:  &FishPlayerPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_FishPlayerPoolID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			fpp:  &FishPlayerPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.GameMessageClassID_FishPlayerPoolID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			fpp:  &FishPlayerPoolHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fpp := &FishPlayerPoolHandler{}
			if got := fpp.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("FishPlayerPoolHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFishPlayerPoolHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.FishPlayerPoolDetail{
		FlagID: proto.Int32(2),
		FishRewardList: []*netproto.FishPlayerPoolDetailList{
			{
				UserID:     proto.Int32(9936738),
				HitMoney:   proto.Int64(0),
				SpendMoney: proto.Int64(100),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.GameMessageClassID_FishPlayerPoolID),
		MsgData:  rpcMessage,
	}

	handler := new(FishPlayerPoolHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.GameMessageClassID_FishPlayerPoolRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.GameMessageClassID_FishPlayerPoolRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.FishPlayerPoolDetail)

	if !ok {
		t.Errorf("Return Message is not *netproto.FishPlayerPoolDetail")
	}

	fmt.Printf("FlagID: %d\n", msg.GetFlagID())
	fmt.Printf("FishRewardList length: %d\n", len(msg.GetFishRewardList()))
	for i, row := range msg.GetFishRewardList() {
		fmt.Printf("FishRewardList-%d: %+v\n", i, row)
	}

}
