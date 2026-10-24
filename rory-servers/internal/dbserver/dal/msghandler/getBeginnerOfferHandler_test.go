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

func TestGetBeginnerOfferHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GetBeginnerOfferHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &GetBeginnerOfferHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_GetBeginnerOfferID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			rh:   &GetBeginnerOfferHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.GameMessageClassID_GetBeginnerOfferID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			rh:   &GetBeginnerOfferHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.GameMessageClassID_FishPlayerPoolID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GetBeginnerOfferHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GetBeginnerOfferHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBeginnerOfferHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.FindFishBeginnerOfferByUserID{
		UserID:                 proto.Int32(9936738),
		NewPlayerAdjustSyspool: proto.Int32(0),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.GameMessageClassID_GetBeginnerOfferID),
		MsgData:  rpcMessage,
	}

	handler := new(GetBeginnerOfferHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.GameMessageClassID_GetBeginnerOfferRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.GameMessageClassID_GetBeginnerOfferRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.FishBeginnerOfferDetail)

	if !ok {
		t.Errorf("Return Message is not *netproto.FishBeginnerOfferDetail")
	}

	fmt.Printf("NewPlayerAdjustSyspool: %d, NewPlayerList length: %d\n", msg.GetNewPlayerAdjustSyspool(), len(msg.GetNewPlayerList()))
	for _, row := range msg.GetNewPlayerList() {
		fmt.Printf("UserID: %d, Score: %d, Status: %d, Bullet: %d\n", row.GetUserID(), row.GetScore(), row.GetStatus(), row.GetBullet())
	}

}
