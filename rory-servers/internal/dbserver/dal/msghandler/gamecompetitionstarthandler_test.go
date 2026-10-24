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

func TestGameCompetitionStartHandler_IsHook_UnitTest(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		fsu  *GameCompetitionStartHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			fsu:  &GameCompetitionStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_BuyTicketID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			fsu:  &GameCompetitionStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.HallMsgClassID_BuyTicketID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			fsu:  &GameCompetitionStartHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsu := &GameCompetitionStartHandler{}
			if got := fsu.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GameCompetitionStartHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGameCompetitionStartHandler_OnMessage_IntegrationTest(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.GetBuyTicketMsg{
		Area_ID: proto.Int32(2),
		FlagID:  proto.Int32(1),
		GameID:  proto.Int32(14),
		UserID:  proto.Int32(9936738),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Hall),
		SClassID: int32(netproto.HallMsgClassID_BuyTicketID),
		MsgData:  rpcMessage,
	}

	handler := new(GameCompetitionStartHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.HallMsgClassID_BuyTicketRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_BuyTicketRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.BuyTicketMsg)

	if !ok {
		t.Errorf("Return Message is not *netproto.BuyTicketMsg")
	}

	fmt.Printf("AreaID: %d, RetCode: %d, RetMsg: %s\n", msg.GetArea_ID(), msg.GetRetCode(), msg.GetRetMsg())
}

func TestGameCompetitionStartHandler_OnMessage_IntegrationTest2(t *testing.T) {
	initDBConnection()

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.GetBuyTicketMsg{
		Area_ID: proto.Int32(2),
		FlagID:  proto.Int32(1),
		GameID:  proto.Int32(25),
		UserID:  proto.Int32(9936738),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Hall),
		SClassID: int32(netproto.HallMsgClassID_BuyTicketID),
		MsgData:  rpcMessage,
	}

	handler := new(GameCompetitionStartHandler)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.HallMsgClassID_BuyTicketRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_BuyTicketRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.BuyTicketMsg)

	if !ok {
		t.Errorf("Return Message is not *netproto.BuyTicketMsg")
	}

	fmt.Printf("AreaID: %d, RetCode: %d, RetMsg: %s\n", msg.GetArea_ID(), msg.GetRetCode(), msg.GetRetMsg())
}
