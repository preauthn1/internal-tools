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

func TestGameCompetitionRuleHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GameCompetitionRuleHandler
		args args
		want bool
	}{
		{
			name: "Hooked Hall",
			rh:   &GameCompetitionRuleHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: true,
		},
		{
			name: "Hooked Game",
			rh:   &GameCompetitionRuleHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: true,
		},
		{
			name: "Not Hooked",
			rh:   &GameCompetitionRuleHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Common),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GameCompetitionRuleHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GameCompetitionRuleHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameCompetitionRuleHandler_OnMessage(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(77874156),
	}

	msgData, _ := proto.Marshal(&netproto.GetGameCompetitionRuleMsg{
		Area_ID: proto.Int32(1),
		GameID:  proto.Int32(14),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Hall),
		SClassID: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
		MsgData:  rpcMessage,
	}

	handler := new(GameCompetitionRuleHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.HallMsgClassID_GameCompetioniInfoRuleRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_GameCompetioniInfoRuleRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.GameCompetitionRuleMsg)

	if !ok {
		t.Errorf("Return Message is not *netproto.GameCompetitionRuleMsg")
	}

	fmt.Printf("%+v\n", msg)

}
