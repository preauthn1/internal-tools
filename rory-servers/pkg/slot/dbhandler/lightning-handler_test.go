package dbhandler

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/dbhandler/service"
)

func TestLightningHandler_IsHook(t *testing.T) {
	type fields struct {
		Service service.LightningService
	}
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Hooked",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID),
			},
			want: true,
		},
		{
			name: "Hooked",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID),
			},
			want: true,
		},
		{
			name: "Not Hooked, Wrong BClassID",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Hall),
				sclsid: int32(netproto.TGW_GameMessageClassID_TGW_RecordBonusGameID),
			},
			want: false,
		},
		{
			name: "Not Hooked, Wrong SClassID",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.DBServerClassID_BlockIPID),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &LightningHandler{
				Service: tt.fields.Service,
			}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("LightningHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBLightningHandler_OnMessage_Record(t *testing.T) {

	handler := NewLightningHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	t.Run("Add New One", func(tt *testing.T) {
		msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRecord{
			UserID: proto.Int32(1),
			GameID: proto.Int32(1003),
			FlagID: proto.Int32(5),
			Data:   proto.String("123122"),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Game),
			SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID),
			MsgData:  rpcMessage,
		}

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}
	})

	t.Run("Empty Data to Reset", func(tt *testing.T) {
		msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRecord{
			UserID: proto.Int32(1),
			GameID: proto.Int32(1003),
			FlagID: proto.Int32(5),
			Data:   proto.String(""),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Game),
			SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID),
			MsgData:  rpcMessage,
		}

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}
	})
}

func TestDBLightningHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRead{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1003),
		FlagID: proto.Int32(5),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID),
		MsgData:  rpcMessage,
	}

	handler := NewLightningHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_LightningRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_LightningRecord")
	}

	t.Logf("UserID: %v\n", msg.GetUserID())
	t.Logf("GameID: %v\n", msg.GetGameID())
	t.Logf("FlagID: %v\n", msg.GetFlagID())
	t.Logf("Data: %v\n", msg.GetData())
}

func TestAPILightningHandler_OnMessage_Record(t *testing.T) {

	handler := NewLightningHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	t.Run("Add New One", func(tt *testing.T) {
		msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRecord{
			UserID: proto.Int32(1),
			GameID: proto.Int32(1003),
			FlagID: proto.Int32(5),
			Data:   proto.String("123122"),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Game),
			SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID),
			MsgData:  rpcMessage,
		}

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}
	})

	t.Run("Empty Data to Reset", func(tt *testing.T) {
		msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRecord{
			UserID: proto.Int32(1),
			GameID: proto.Int32(1003),
			FlagID: proto.Int32(5),
			Data:   proto.String(""),
		})

		rpcMessage := &rpcmsg.RPCMessage{
			RPCInfo: rpcInfo,
			MsgData: msgData,
		}

		m := &msg.Message{
			BClassID: int32(netproto.MessageBClassID_Game),
			SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID),
			MsgData:  rpcMessage,
		}

		ret := handler.OnMessage(ctx, clt, m)

		if len(ret) == 0 {
			tt.Errorf("No Message Return")
		}
	})
}

func TestAPILightningHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_LightningRead{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1003),
		FlagID: proto.Int32(5),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID),
		MsgData:  rpcMessage,
	}

	handler := NewLightningHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_LightningRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_LightningRecord")
	}

	t.Logf("UserID: %v\n", msg.GetUserID())
	t.Logf("GameID: %v\n", msg.GetGameID())
	t.Logf("FlagID: %v\n", msg.GetFlagID())
	t.Logf("Data: %v\n", msg.GetData())
}
