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

func TestDeferredAwardHandler_IsHook(t *testing.T) {
	type fields struct {
		Service service.DeferredAwardService
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
			name: "Hooked Read",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID),
			},
			want: true,
		},
		{
			name: "Hooked Record",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID),
			},
			want: true,
		},
		{
			name: "Hooked Pop",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID),
			},
			want: true,
		},
		{
			name: "Hooked Reset",
			fields: fields{
				Service: nil,
			},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID),
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
			h := &DeferredAwardHandler{
				Service: tt.fields.Service,
			}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("DeferredAwardHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBDeferredAwardHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_DeferredAwardRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_DeferredAwardRecord")
	}

	t.Log("DeferredAward")
	t.Log(" - UserID", msg.GetUserID())
	t.Log(" - GameID", msg.GetGameID())
	t.Log(" - FlagID", msg.GetFlagID())
	t.Logf(" - Records(%d): \n", len(msg.GetRecords()))
	for _, row := range msg.GetRecords() {
		t.Logf("   - Record: %v\n", row)
	}
}

func TestDBDeferredAwardHandler_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRecord{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1),
		FlagID: proto.Int32(1),
		Records: []*netproto.Slot_DeferredAward{
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestDBDeferredAwardHandler_OnMessage_Pop(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestDBDeferredAwardHandler_OnMessage_Reset(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestAPIDeferredAwardHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_DeferredAwardRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_DeferredAwardRecord")
	}

	t.Log("DeferredAward")
	t.Log(" - UserID", msg.GetUserID())
	t.Log(" - GameID", msg.GetGameID())
	t.Log(" - FlagID", msg.GetFlagID())
	t.Logf(" - Records(%d): \n", len(msg.GetRecords()))
	for _, row := range msg.GetRecords() {
		t.Logf("   - Record: %v\n", row)
	}
}

func TestAPIDeferredAwardHandler_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRecord{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1),
		FlagID: proto.Int32(1),
		Records: []*netproto.Slot_DeferredAward{
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
			&netproto.Slot_DeferredAward{
				AwardType: proto.String("Test"),
				OriginID:  proto.Int32(1),
				Data:      proto.String("123"),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestAPIDeferredAwardHandler_OnMessage_Pop(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestAPIDeferredAwardHandler_OnMessage_Reset(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_DeferredAwardRead{
		UserID:    proto.Int32(1),
		GameID:    proto.Int32(1),
		FlagID:    proto.Int32(1),
		AwardType: proto.String("Test"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID),
		MsgData:  rpcMessage,
	}

	handler := NewDeferredAwardHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}
