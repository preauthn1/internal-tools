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

func TestJackpotNoLinkHandler_IsHook(t *testing.T) {
	type fields struct {
		Service service.JackpotNoLinkService
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID),
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID),
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
			h := &JackpotNoLinkHandler{
				Service: tt.fields.Service,
			}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("JackpotNoLinkHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBJackpotNoLinkHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_JackpotRead{
		Info: &netproto.Slot_JackpotData{
			GameID:            proto.Int32(25),
			FlagID:            proto.Int32(1),
			TableNo:           proto.Int32(0),
			Level:             proto.Int32(1),
			Contribution:      proto.Int64(0),
			TotalContribution: proto.Int64(0),
			AwardTotal:        proto.Int64(0),
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID),
		MsgData:  rpcMessage,
	}

	handler := NewJackpotNoLinkHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_JackpotRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_JackpotRecord")
	}

	t.Logf("Jackpot Records Length: %d\n", len(msg.Records))
	for _, row := range msg.Records {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestDBJackpotNoLinkHandler_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_JackpotRecord{
		Records: []*netproto.Slot_JackpotData{
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(25),
				FlagID:            proto.Int32(1),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(1),
				Contribution:      proto.Int64(1),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(25),
				FlagID:            proto.Int32(1),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(2),
				Contribution:      proto.Int64(2),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(25),
				FlagID:            proto.Int32(1),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(3),
				Contribution:      proto.Int64(3),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(25),
				FlagID:            proto.Int32(1),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(4),
				Contribution:      proto.Int64(4),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewJackpotNoLinkHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_JackpotRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_JackpotRecord")
	}

	t.Logf("Jackpot Records Length: %d\n", len(msg.Records))
	for _, row := range msg.Records {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestAPIJackpotNoLinkHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(9936738),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_JackpotRead{
		Info: &netproto.Slot_JackpotData{
			GameID:            proto.Int32(1005),
			FlagID:            proto.Int32(41),
			TableNo:           proto.Int32(0),
			Level:             proto.Int32(1),
			Contribution:      proto.Int64(0),
			TotalContribution: proto.Int64(0),
			AwardTotal:        proto.Int64(0),
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID),
		MsgData:  rpcMessage,
	}

	handler := NewJackpotNoLinkHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_JackpotRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_JackpotRecord")
	}

	t.Logf("Jackpot Records Length: %d\n", len(msg.Records))
	for _, row := range msg.Records {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestAPIJackpotNoLinkHandler_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_JackpotRecord{
		Records: []*netproto.Slot_JackpotData{
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(1023),
				FlagID:            proto.Int32(41),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(1),
				Contribution:      proto.Int64(1),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(1023),
				FlagID:            proto.Int32(41),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(2),
				Contribution:      proto.Int64(2),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(1023),
				FlagID:            proto.Int32(41),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(3),
				Contribution:      proto.Int64(3),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
			&netproto.Slot_JackpotData{
				GameID:            proto.Int32(1023),
				FlagID:            proto.Int32(41),
				TableNo:           proto.Int32(0),
				Level:             proto.Int32(4),
				Contribution:      proto.Int64(4),
				TotalContribution: proto.Int64(0),
				AwardTotal:        proto.Int64(0),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewJackpotNoLinkHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_JackpotRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_JackpotRecord")
	}

	t.Logf("Jackpot Records Length: %d\n", len(msg.Records))
	for _, row := range msg.Records {
		t.Logf("  - Record: %v\n", row)
	}
}
