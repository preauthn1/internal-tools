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

func TestGetMailDetailHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		rh   *GetMailDetailHandler
		args args
		want bool
	}{
		{
			name: "Hooked",
			rh:   &GetMailDetailHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_GetMailDetailID),
			},
			want: true,
		},
		{
			name: "Not Hooked BClass",
			rh:   &GetMailDetailHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Common),
				sclsid: int32(netproto.HallMsgClassID_GetMailDetailID),
			},
			want: false,
		},
		{
			name: "Not Hooked SClass",
			rh:   &GetMailDetailHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.HallMsgClassID_GameCompetioniInfoRuleID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rh := &GetMailDetailHandler{}
			if got := rh.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GetMailDetailHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetMailDetailHandler_OnMessage(t *testing.T) {
	initDBConnection()

	t.Run("Has Data", func(tt *testing.T) {
		rpcInfo := &netproto.RPCInfo{
			RouteServerID: proto.Int32(1),
			ConnID:        proto.String(""),
			IPAddress:     proto.String("local"),
			Cer:           proto.String(""),
			UserID:        proto.Int32(3748506),
		}

		msgData, _ := proto.Marshal(&netproto.RequestMailDetail{
			Area_ID: proto.Int32(2),
			ID:      proto.Int32(15157),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Hall),
			SClassID: int32(netproto.HallMsgClassID_GetMailDetailID),
			MsgData:  rpcMessage,
		}

		handler := new(GetMailDetailHandler)

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}

		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
			tt.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
		}

		if ret[0].SClassID != int32(netproto.HallMsgClassID_GetMailDetailRetID) {
			tt.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_GetMailDetailRetID), ret[0].SClassID)
		}

		msg, ok := ret[0].MsgData.(*netproto.MailDetail)

		if !ok {
			tt.Errorf("Return Message is not *netproto.MailDetail")
		}

		fmt.Printf("AreaID: %v ID: %v, Title: %v, IsRead: %v, SendTime: %v, Content: %v\n", msg.GetArea_ID(), msg.GetID(), msg.GetTitle(), msg.GetIsRead(), msg.GetSendTime(), msg.GetContent())

	})

	t.Run("No Data", func(tt *testing.T) {
		rpcInfo := &netproto.RPCInfo{
			RouteServerID: proto.Int32(1),
			ConnID:        proto.String(""),
			IPAddress:     proto.String("local"),
			Cer:           proto.String(""),
			UserID:        proto.Int32(1),
		}

		msgData, _ := proto.Marshal(&netproto.RequestMailDetail{
			Area_ID: proto.Int32(2),
			ID:      proto.Int32(1),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Hall),
			SClassID: int32(netproto.HallMsgClassID_GetMailDetailID),
			MsgData:  rpcMessage,
		}

		handler := new(GetMailDetailHandler)

		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}

		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
			tt.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
		}

		if ret[0].SClassID != int32(netproto.HallMsgClassID_GetMailDetailRetID) {
			tt.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_GetMailDetailRetID), ret[0].SClassID)
		}

		msg, ok := ret[0].MsgData.(*netproto.MailDetail)

		if !ok {
			tt.Errorf("Return Message is not *netproto.MailDetail")
		}

		fmt.Printf("AreaID: %v ID: %v, Title: %v, IsRead: %v, SendTime: %v, Content: %v\n", msg.GetArea_ID(), msg.GetID(), msg.GetTitle(), msg.GetIsRead(), msg.GetSendTime(), msg.GetContent())

	})

}
