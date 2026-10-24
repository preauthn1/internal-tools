package dbhandler

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/dbhandler/service"
)

func TestGameLogHandler_IsHook(t *testing.T) {
	type fields struct {
		Service service.GameLogService
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID),
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
				sclsid: int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID),
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
			h := &GameLogHandler{
				Service: tt.fields.Service,
			}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("GameLogHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBGameLogHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_GameLogRead{
		UserID:    proto.Int32(15267553),
		GameID:    proto.Int32(1004),
		PageIndex: proto.Int32(1),
		PageSize:  proto.Int32(10),
		StartTime: proto.String("2020-03-20"),
		EndTime:   proto.String("2020-03-20"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_GameLogRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_GameLogRecord")
	}

	t.Logf("GameLog UserID: %d\n", msg.GetUserID())
	t.Logf("GameLog GameID: %d\n", msg.GetGameID())
	t.Logf("GameLog Total Record: %d\n", msg.GetPageCount())
	t.Logf("GameLog Records Length: %d\n", len(msg.Logs))
	for _, row := range msg.Logs {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestDBGameLogHandler_OnMessage_Add(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_GameLogAdd{
		UserID:        proto.Int32(15267553),
		GameID:        proto.Int32(1004),
		ServerID:      proto.Int32(414),
		CurrentAmount: proto.Int64(100),
		BetMoney:      proto.Int64(10),
		WinMoney:      proto.Int64(30),
		FinalScore:    proto.Int64(120),
		GameType:      proto.Int32(1),
		LogType:       proto.Int32(1),
		GameData:      proto.String("[test add]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
		CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	t.Log(ret[0])
}

func TestDBGameLogHandler_OnMessage_BatchAdd(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BatchGameLogAdd{
		Logs: []*netproto.Slot_MSG_GameLogAdd{
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(15267553),
				GameID:        proto.Int32(1004),
				ServerID:      proto.Int32(414),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 1]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(15267553),
				GameID:        proto.Int32(1004),
				ServerID:      proto.Int32(414),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 2]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(15267553),
				GameID:        proto.Int32(1004),
				ServerID:      proto.Int32(414),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 3]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(DB)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	t.Log(ret[0])
}

func TestAPIGameLogHandler_OnMessage_Read(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_GameLogRead{
		UserID:    proto.Int32(1000174),
		GameID:    proto.Int32(1011),
		PageIndex: proto.Int32(1),
		PageSize:  proto.Int32(10),
		StartTime: proto.String("2020-03-20"),
		EndTime:   proto.String("2020-03-20"),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	}

	if ret[0].SClassID != int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID) {
		t.Errorf("Message SClassID want %d, but %d", int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID), ret[0].SClassID)
	}

	msg, ok := ret[0].MsgData.(*netproto.Slot_MSG_GameLogRecord)

	if !ok {
		t.Errorf("Return Message is not *netproto.Slot_MSG_GameLogRecord")
	}

	t.Logf("GameLog UserID: %d\n", msg.GetUserID())
	t.Logf("GameLog GameID: %d\n", msg.GetGameID())
	t.Logf("GameLog Total Record: %d\n", msg.GetPageCount())
	t.Logf("GameLog Records Length: %d\n", len(msg.Logs))
	for _, row := range msg.Logs {
		t.Logf("  - Record: %v\n", row)
	}
}

func TestAPIGameLogHandler_OnMessage_Add(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_GameLogAdd{
		UserID:        proto.Int32(1000174),
		GameID:        proto.Int32(1011),
		ServerID:      proto.Int32(475),
		CurrentAmount: proto.Int64(100),
		BetMoney:      proto.Int64(10),
		WinMoney:      proto.Int64(30),
		FinalScore:    proto.Int64(120),
		GameType:      proto.Int32(1),
		LogType:       proto.Int32(1),
		GameData:      proto.String("[test add]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
		CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	t.Log(ret[0])
}

func TestAPIGameLogHandler_OnMessage_BatchAdd(t *testing.T) {

	rpcInfo := &netproto.RPCInfo{
		RouteServerID: proto.Int32(1),
		ConnID:        proto.String(""),
		IPAddress:     proto.String("local"),
		Cer:           proto.String(""),
		UserID:        proto.Int32(7276429),
	}

	msgData, _ := proto.Marshal(&netproto.Slot_MSG_BatchGameLogAdd{
		Logs: []*netproto.Slot_MSG_GameLogAdd{
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(1000174),
				GameID:        proto.Int32(1011),
				ServerID:      proto.Int32(475),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 1]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(1000174),
				GameID:        proto.Int32(1011),
				ServerID:      proto.Int32(475),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 2]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
			&netproto.Slot_MSG_GameLogAdd{
				UserID:        proto.Int32(1000174),
				GameID:        proto.Int32(1011),
				ServerID:      proto.Int32(475),
				CurrentAmount: proto.Int64(100),
				BetMoney:      proto.Int64(10),
				WinMoney:      proto.Int64(30),
				FinalScore:    proto.Int64(120),
				GameType:      proto.Int32(1),
				LogType:       proto.Int32(1),
				GameData:      proto.String("[test batch 3]101#p5,p5,p7|p3,p1,p7|p5,p1,p6|p2,p1,p7|p6,p7,p6"),
				CreateTime:    proto.String(time.Now().Format("2006-01-02 15:04:05.000")),
			},
		},
	})

	rpcMessage := &rpcmsg.RPCMessage{
		RPCInfo: rpcInfo,
		MsgData: msgData,
	}

	m := &msg.Message{
		BClassID: int32(netproto.MessageBClassID_Game),
		SClassID: int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID),
		MsgData:  rpcMessage,
	}

	handler := NewGameLogHandler(API)

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)

	ret := handler.OnMessage(ctx, clt, m)

	if len(ret) == 0 {
		t.Errorf("No Message Return")
	}

	t.Log(ret[0])
}
