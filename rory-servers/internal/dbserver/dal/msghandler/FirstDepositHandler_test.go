package msghandler

// import (
// 	"example.invalid/golang/protobuf/proto"
// )

// import (
// 	"example.invalid/ServerTeam/Servers/pkg/module"
// 	"example.invalid/ServerTeam/Servers/pkg/msg"
// 	"example.invalid/ServerTeam/Servers/pkg/netproto"
// 	"example.invalid/ServerTeam/Servers/pkg/network"
// 	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
// )

// import (
// 	"fmt"
// 	"testing"
// )

// func TestFirstDepositHandler_IsHook(t *testing.T) {
// 	type args struct {
// 		bclsid int32
// 		sclsid int32
// 	}
// 	tests := []struct {
// 		name string
// 		h    *FirstDepositHandler
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "Hooked GetInfo",
// 			h:    &FirstDepositHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Hall),
// 				sclsid: int32(netproto.HallMsgClassID_FirstDepositGetInfoID),
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "Hooked Collect",
// 			h:    &FirstDepositHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Hall),
// 				sclsid: int32(netproto.HallMsgClassID_FirstDepositCollectID),
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "Not Hooked, Wrong BClassID",
// 			h:    &FirstDepositHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Game),
// 				sclsid: int32(netproto.HallMsgClassID_FirstDepositGetInfoID),
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "Not Hooked, Wrong SClassID",
// 			h:    &FirstDepositHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Hall),
// 				sclsid: int32(netproto.DBServerClassID_BlockIPID),
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &FirstDepositHandler{}
// 			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
// 				t.Errorf("FirstDepositHandler.IsHook() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestFirstDepositHandler_OnMessage_GetInfo(t *testing.T) {
// 	initDBConnection()

// 	rpcInfo := &netproto.RPCInfo{
// 		RouteServerID: proto.Int32(1),
// 		ConnID:        proto.String(""),
// 		IPAddress:     proto.String("local"),
// 		Cer:           proto.String(""),
// 		UserID:        proto.Int32(9936738),
// 	}

// 	msgData, _ := proto.Marshal(&netproto.FirstDepositGetInfoMsg{
// 		UserID: proto.Int32(9936738),
// 	})

// 	rpcMessage := &rpcmsg.RPCMessage{
// 		RPCInfo: rpcInfo,
// 		MsgData: msgData,
// 	}

// 	m := &msg.Message{
// 		BClassID: int32(netproto.MessageBClassID_Hall),
// 		SClassID: int32(netproto.HallMsgClassID_FirstDepositGetInfoID),
// 		MsgData:  rpcMessage,
// 	}

// 	handler := new(FirstDepositHandler)

// 	ctx := new(module.MouduleCtx)
// 	clt := new(network.ConnInfo)

// 	ret := handler.OnMessage(ctx, clt, m)

// 	if len(ret) == 0 {
// 		t.Errorf("No Message Return")
// 	}

// 	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 	}

// 	if ret[0].SClassID != int32(netproto.HallMsgClassID_FirstDepositGetInfoRetID) {
// 		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_FirstDepositGetInfoRetID), ret[0].SClassID)
// 	}

// 	msg, ok := ret[0].MsgData.(*netproto.FirstDepositGetInfoRetMsg)

// 	if !ok {
// 		t.Errorf("Return Message is not *netproto.FirstDepositGetInfoRetMsg")
// 	}

// 	fmt.Printf("RetCode: %d\n", msg.GetRetCode())

// 	if msg.GetFirstDepositInfo() != nil {
// 		info := msg.GetFirstDepositInfo()
// 		fmt.Printf("StartTime: %v\n", info.GetStartTime())
// 		fmt.Printf("EndTime: %v\n", info.GetEndTime())
// 		fmt.Printf("DepositLimit: %v\n", info.GetDepositLimit())
// 		fmt.Printf("WinAmount: %v\n", info.GetWinAmount())
// 		fmt.Printf("ReturnAmount: %v\n", info.GetReturnAmount())
// 		fmt.Printf("Status: %v\n", info.GetStatus())
// 		fmt.Printf("CurrentUserCount: %v\n", info.GetCurrentUserCount())
// 	}

// 	if msg.GetFirstDepositUserRecord() != nil {
// 		record := msg.GetFirstDepositUserRecord()
// 		fmt.Printf("UserID: %v\n", record.GetUserID())
// 		fmt.Printf("EventID: %v\n", record.GetEventID())
// 		fmt.Printf("JoinTime: %v\n", record.GetJoinTime())
// 		fmt.Printf("DepositAmount: %v\n", record.GetDepositAmount())
// 		fmt.Printf("WinAmount: %v\n", record.GetWinAmount())
// 		fmt.Printf("Flag: %v\n", record.GetFlag())
// 		fmt.Printf("CollectTime: %v\n", record.GetCollectTime())
// 		fmt.Printf("CollectAmount: %v\n", record.GetCollectAmount())
// 		fmt.Printf("UserType: %v\n", record.GetUserType())
// 	}
// }

// func TestFirstDepositHandler_OnMessage_Collect(t *testing.T) {
// 	initDBConnection()

// 	rpcInfo := &netproto.RPCInfo{
// 		RouteServerID: proto.Int32(1),
// 		ConnID:        proto.String(""),
// 		IPAddress:     proto.String("local"),
// 		Cer:           proto.String(""),
// 		UserID:        proto.Int32(9936738),
// 	}

// 	msgData, _ := proto.Marshal(&netproto.FirstDepositCollectMsg{
// 		UserID:    proto.Int32(9936738),
// 		EventID:   proto.Int64(1),
// 		IPAddress: proto.String("0.0.0.0"),
// 	})

// 	rpcMessage := &rpcmsg.RPCMessage{
// 		RPCInfo: rpcInfo,
// 		MsgData: msgData,
// 	}

// 	m := &msg.Message{
// 		BClassID: int32(netproto.MessageBClassID_Hall),
// 		SClassID: int32(netproto.HallMsgClassID_FirstDepositCollectID),
// 		MsgData:  rpcMessage,
// 	}

// 	handler := new(FirstDepositHandler)

// 	ctx := new(module.MouduleCtx)
// 	clt := new(network.ConnInfo)

// 	ret := handler.OnMessage(ctx, clt, m)

// 	if len(ret) == 0 {
// 		t.Errorf("No Message Return")
// 	}

// 	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 	}

// 	if ret[0].SClassID != int32(netproto.HallMsgClassID_FirstDepositCollectRetID) {
// 		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_FirstDepositCollectRetID), ret[0].SClassID)
// 	}

// 	msg, ok := ret[0].MsgData.(*netproto.FirstDepositCollectRetMsg)

// 	if !ok {
// 		t.Errorf("Return Message is not *netproto.FirstDepositCollectRetMsg")
// 	}

// 	fmt.Printf("RetCode: %d, RetMessage: %s\n", msg.GetRetCode(), msg.GetRetMessage())

// }
