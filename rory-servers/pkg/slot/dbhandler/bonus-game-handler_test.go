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

func TestBonusGameHandler_IsHook(t *testing.T) {
	type fields struct {
		Service service.BonusGameService
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID),
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID),
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID),
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
			h := &BonusGameHandler{
				Service: tt.fields.Service,
			}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("BonusGameHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBonusGameHandler_DB_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRead{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1),
		FlagID: proto.Int32(1),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_BonusGameRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_BonusGameRecord")
	}

	t.Logf("BonusGameRecord Lenght: %d\n", len(msg.GetDetail()))
	for _, row := range msg.GetDetail() {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestBonusGameHandler_DB_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRecord{
		Detail: []*netproto.Slot_BonusGameDetail{
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(1), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(2), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(3), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(4), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(5), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(6), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestBonusGameHandler_DB_OnMessage_Reset(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRead{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1),
		FlagID: proto.Int32(1),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestBonusGameHandler_API_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRead{
		UserID: proto.Int32(1),
		GameID: proto.Int32(1),
		FlagID: proto.Int32(1),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_BonusGameRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_BonusGameRecord")
	}

	t.Logf("BonusGameRecord Lenght: %d\n", len(msg.GetDetail()))
	for _, row := range msg.GetDetail() {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestBonusGameHandler_API_OnMessage_Record(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRecord{
		Detail: []*netproto.Slot_BonusGameDetail{
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(1), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(2), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(3), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(4), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(5), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
			&netproto.Slot_BonusGameDetail{UserID: proto.Int32(1), GameID: proto.Int32(1), FlagID: proto.Int32(1), OriginID: proto.Int32(1), TurnID: proto.Int32(6), Location: proto.Int32(0), BetMoney: proto.Int32(1), BetRate: proto.Int32(1), Data: proto.String("")},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}

func TestBonusGameHandler_API_OnMessage_Reset(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(1),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BonusGameRead{
		UserID:   proto.Int32(1),
		GameID:   proto.Int32(1),
		FlagID:   proto.Int32(1),
		OriginID: proto.Int32(1),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID),
		MsgData:  rpcMessage,
	}

	handler := NewBonusGameHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}
}
