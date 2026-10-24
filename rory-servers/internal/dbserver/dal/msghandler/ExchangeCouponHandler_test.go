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
// 	"time"
// )

// func TestExchangeCouponHandler_IsHook(t *testing.T) {
// 	type args struct {
// 		bclsid int32
// 		sclsid int32
// 	}
// 	tests := []struct {
// 		name string
// 		h    *ExchangeCouponHandler
// 		args args
// 		want bool
// 	}{
// 		{
// 			name: "Hooked",
// 			h:    &ExchangeCouponHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Hall),
// 				sclsid: int32(netproto.HallMsgClassID_ExchangeCouponID),
// 			},
// 			want: true,
// 		},
// 		{
// 			name: "Not Hooked, Wrong BClassID",
// 			h:    &ExchangeCouponHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Game),
// 				sclsid: int32(netproto.HallMsgClassID_ExchangeCouponID),
// 			},
// 			want: false,
// 		},
// 		{
// 			name: "Not Hooked, Wrong SClassID",
// 			h:    &ExchangeCouponHandler{},
// 			args: args{
// 				bclsid: int32(netproto.MessageBClassID_Hall),
// 				sclsid: int32(netproto.DBServerClassID_BlockIPID),
// 			},
// 			want: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			h := &ExchangeCouponHandler{}
// 			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
// 				t.Errorf("ExchangeCouponHandler.IsHook() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

// func TestExchangeCouponHandler_OnMessage(t *testing.T) {
// 	initDBConnection()

// 	rpcInfo := &netproto.RPCInfo{
// 		RouteServerID: proto.Int32(1),
// 		ConnID:        proto.String(""),
// 		IPAddress:     proto.String("local"),
// 		Cer:           proto.String(""),
// 		UserID:        proto.Int32(9936738),
// 	}

// 	msgData, _ := proto.Marshal(&netproto.ExchangeCouponMsg{
// 		ExchangeSequence: proto.String("tROxHvPwq6jz"),
// 	})

// 	rpcMessage := &rpcmsg.RPCMessage{
// 		RPCInfo: rpcInfo,
// 		MsgData: msgData,
// 	}

// 	m := &msg.Message{
// 		BClassID: int32(netproto.MessageBClassID_Hall),
// 		SClassID: int32(netproto.HallMsgClassID_ExchangeCouponID),
// 		MsgData:  rpcMessage,
// 	}

// 	handler := new(ExchangeCouponHandler)

// 	ctx := new(module.MouduleCtx)
// 	clt := new(network.ConnInfo)

// 	ret := handler.OnMessage(ctx, clt, m)

// 	if len(ret) == 0 {
// 		t.Errorf("No Message Return")
// 	}

// 	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 		t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 	}

// 	if ret[0].SClassID != int32(netproto.HallMsgClassID_ExchangeCouponRetID) {
// 		t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_ExchangeCouponRetID), ret[0].SClassID)
// 	}

// 	msg, ok := ret[0].MsgData.(*netproto.ExchangeCouponRetMsg)

// 	if !ok {
// 		t.Errorf("Return Message is not *netproto.ExchangeCouponRetMsg")
// 	}

// 	fmt.Printf("Code: %d, Message: %s\n", msg.GetRetCode(), msg.GetChvMsg())
// }

// func TestExchangeCouponHandler_BlockUser(t *testing.T) {
// 	h := new(ExchangeCouponHandler)

// 	var userID int32 = 978
// 	var userID2 int32 = 987

// 	if h.IsBlocked(userID) {
// 		t.Errorf("User should not be blocked, but blocked")
// 	}

// 	h.BlockUser(userID, 3*time.Second)

// 	// user1 被 block
// 	if !h.IsBlocked(userID) {
// 		t.Errorf("User should be blocked, but not blocked")
// 	}
// 	// user2 不可被 block
// 	if h.IsBlocked(userID2) {
// 		t.Errorf("User should not be blocked, but blocked")
// 	}

// 	for i := 0; i < 10; i++ {
// 		time.Sleep(1 * time.Second)
// 		fmt.Printf("Time after %d sec, IsBlocked() is %t\n", i+1, h.IsBlocked(userID))
// 		if i < 2 {
// 			if !h.IsBlocked(userID) {
// 				t.Errorf("User should be blocked, but not blocked")
// 			}
// 		} else if i > 3 {
// 			if h.IsBlocked(userID) {
// 				t.Errorf("User should not be blocked, but blocked")
// 			}
// 		}
// 	}

// }

// func BenchmarkExchangeCouponHandler_BlockUser(b *testing.B) {
// 	h := new(ExchangeCouponHandler)

// 	for i := 0; i < b.N; i++ {
// 		go h.BlockUser(int32(i), time.Second)
// 	}
// }

// func TestExchangeCouponHandler_OnMessage_UserBlocked(t *testing.T) {
// 	initDBConnection()

// 	rpcInfo := &netproto.RPCInfo{
// 		RouteServerID: proto.Int32(1),
// 		ConnID:        proto.String(""),
// 		IPAddress:     proto.String("local"),
// 		Cer:           proto.String(""),
// 		UserID:        proto.Int32(9936738),
// 	}

// 	msgData, _ := proto.Marshal(&netproto.ExchangeCouponMsg{
// 		ExchangeSequence: proto.String("tROxHvPwq6jz"),
// 	})

// 	rpcMessage := &rpcmsg.RPCMessage{
// 		RPCInfo: rpcInfo,
// 		MsgData: msgData,
// 	}

// 	m := &msg.Message{
// 		BClassID: int32(netproto.MessageBClassID_Hall),
// 		SClassID: int32(netproto.HallMsgClassID_ExchangeCouponID),
// 		MsgData:  rpcMessage,
// 	}

// 	handler := new(ExchangeCouponHandler)

// 	ctx := new(module.MouduleCtx)
// 	clt := new(network.ConnInfo)

// 	// 第一次
// 	{
// 		ret := handler.OnMessage(ctx, clt, m)

// 		if len(ret) == 0 {
// 			t.Errorf("No Message Return")
// 		}

// 		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 			t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 		}

// 		if ret[0].SClassID != int32(netproto.HallMsgClassID_ExchangeCouponRetID) {
// 			t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_ExchangeCouponRetID), ret[0].SClassID)
// 		}

// 		msg, ok := ret[0].MsgData.(*netproto.ExchangeCouponRetMsg)

// 		if !ok {
// 			t.Errorf("Return Message is not *netproto.ExchangeCouponRetMsg")
// 		}

// 		fmt.Printf("Code: %d, Message: %s\n", msg.GetRetCode(), msg.GetChvMsg())
// 	}

// 	// 第二次
// 	{
// 		ret := handler.OnMessage(ctx, clt, m)

// 		if len(ret) == 0 {
// 			t.Errorf("No Message Return")
// 		}

// 		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 			t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 		}

// 		if ret[0].SClassID != int32(netproto.HallMsgClassID_ExchangeCouponRetID) {
// 			t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_ExchangeCouponRetID), ret[0].SClassID)
// 		}

// 		msg, ok := ret[0].MsgData.(*netproto.ExchangeCouponRetMsg)

// 		if !ok {
// 			t.Errorf("Return Message is not *netproto.ExchangeCouponRetMsg")
// 		}

// 		expectedCode := int32(0)
// 		expectedMsg := "请勿短时间连续输入，稍后再试"

// 		fmt.Printf("Code: %d, Message: %s\n", msg.GetRetCode(), msg.GetChvMsg())

// 		if expectedCode != msg.GetRetCode() {
// 			t.Errorf("Return Code want %d, but %d", expectedCode, msg.GetRetCode())
// 		}

// 		if expectedMsg != msg.GetChvMsg() {
// 			t.Errorf("Return Message want %s, but %s", expectedMsg, msg.GetChvMsg())
// 		}
// 	}

// 	// 第三次
// 	{
// 		ret := handler.OnMessage(ctx, clt, m)

// 		if len(ret) == 0 {
// 			t.Errorf("No Message Return")
// 		}

// 		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 			t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 		}

// 		if ret[0].SClassID != int32(netproto.HallMsgClassID_ExchangeCouponRetID) {
// 			t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_ExchangeCouponRetID), ret[0].SClassID)
// 		}

// 		msg, ok := ret[0].MsgData.(*netproto.ExchangeCouponRetMsg)

// 		if !ok {
// 			t.Errorf("Return Message is not *netproto.ExchangeCouponRetMsg")
// 		}

// 		expectedCode := int32(0)
// 		expectedMsg := "请勿短时间连续输入，稍后再试"

// 		fmt.Printf("Code: %d, Message: %s\n", msg.GetRetCode(), msg.GetChvMsg())

// 		if expectedCode != msg.GetRetCode() {
// 			t.Errorf("Return Code want %d, but %d", expectedCode, msg.GetRetCode())
// 		}

// 		if expectedMsg != msg.GetChvMsg() {
// 			t.Errorf("Return Message want %s, but %s", expectedMsg, msg.GetChvMsg())
// 		}
// 	}

// 	// 睡一下
// 	time.Sleep(11 * time.Second)

// 	// 超過時間要過
// 	{
// 		ret := handler.OnMessage(ctx, clt, m)

// 		if len(ret) == 0 {
// 			t.Errorf("No Message Return")
// 		}

// 		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) {
// 			t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Hall), ret[0].BClassID)
// 		}

// 		if ret[0].SClassID != int32(netproto.HallMsgClassID_ExchangeCouponRetID) {
// 			t.Errorf("Message SClassID want %d, but %d", int32(netproto.HallMsgClassID_ExchangeCouponRetID), ret[0].SClassID)
// 		}

// 		msg, ok := ret[0].MsgData.(*netproto.ExchangeCouponRetMsg)

// 		if !ok {
// 			t.Errorf("Return Message is not *netproto.ExchangeCouponRetMsg")
// 		}

// 		expectedMsg := "请勿短时间连续输入，稍后再试"

// 		fmt.Printf("Code: %d, Message: %s\n", msg.GetRetCode(), msg.GetChvMsg())

// 		if expectedMsg == msg.GetChvMsg() {
// 			t.Errorf("Return Message dont want %s, but %s", expectedMsg, msg.GetChvMsg())
// 		}
// 	}
// }
