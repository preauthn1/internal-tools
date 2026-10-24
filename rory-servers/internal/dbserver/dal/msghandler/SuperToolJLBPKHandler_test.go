package msghandler

import (
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"testing"
)

func TestSuperToolJLBPKHandler_IsHook(t *testing.T) {
	type args struct {
		bclsid int32
		sclsid int32
	}
	tests := []struct {
		name string
		h    *SuperToolJLBPKHandler
		args args
		want bool
	}{
		{
			name: "Hooked GetVIPUserList",
			h:    &SuperToolJLBPKHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserList),
			},
			want: true,
		},
		{
			name: "Hooked UpdateVIPUser",
			h:    &SuperToolJLBPKHandler{},
			args: args{
				bclsid: int32(netproto.MessageBClassID_Game),
				sclsid: int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_UpdateVIPUser),
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SuperToolJLBPKHandler{}
			if got := h.IsHook(tt.args.bclsid, tt.args.sclsid); got != tt.want {
				t.Errorf("SuperToolJLBPKHandler.IsHook() = %v, want %v", got, tt.want)
			}
		})
	}
}

// 要求 監控玩家名單
func TestSuperToolJLBPKHandler_OnMessage_GetVIPUserList(t *testing.T) {
	//initDBConnection()
	//
	//rpcInfo := &netproto.RPCInfo{
	//	RouteServerID: proto.Int32(1),
	//	ConnID:        proto.String(""),
	//	IPAddress:     proto.String("local"),
	//	Cer:           proto.String(""),
	//	UserID:        proto.Int32(9936738),
	//}
	//
	//rpcMessage := &rpcmsg.RPCMessage{
	//	RPCInfo: rpcInfo,
	//}
	//
	//m := &msg.Message{
	//	BClassID: int32(netproto.MessageBClassID_Game),
	//	SClassID: int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserList),
	//	MsgData:  rpcMessage,
	//}
	//
	//handler := new(SuperToolJLBPKHandler)
	//
	//ctx := new(module.MouduleCtx)
	//clt := new(network.ConnInfo)
	//
	//ret := handler.OnMessage(ctx, clt, m)
	//
	//if len(ret) == 0 {
	//	t.Errorf("No Message Return")
	//}
	//
	//if ret[0].BClassID != int32(netproto.MessageBClassID_Game) {
	//	t.Errorf("Message BClassID want %d, but %d", int32(netproto.MessageBClassID_Game), ret[0].BClassID)
	//}
	//
	//if ret[0].SClassID != int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserListRet) {
	//	t.Errorf("Message SClassID want %d, but %d", int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserListRet), ret[0].SClassID)
	//}
	//
	//msg, ok := ret[0].MsgData.(*netproto.JLBPK_SuperTool_Msg_GetVIPUserListRet)
	//
	//if !ok {
	//	t.Errorf("Return Message is not *netproto.JLBPK_SuperTool_Msg_GetVIPUserListRet")
	//}
	//
	//fmt.Printf("Records length: %d\n", len(msg.GetUserList()))
	//for _, row := range msg.GetUserList() {
	//	fmt.Println(row)
	//}
}

// 更新 監控玩家資料
func TestSuperToolJLBPKHandler_OnMessage_UpdateVIPUser(t *testing.T) {
	//initDBConnection()
	//
	//rpcInfo := &netproto.RPCInfo{
	//	RouteServerID: proto.Int32(1),
	//	ConnID:        proto.String(""),
	//	IPAddress:     proto.String("local"),
	//	Cer:           proto.String(""),
	//	UserID:        proto.Int32(9936738),
	//}
	//
	//msgData, _ := proto.Marshal(&netproto.JLBPK_SuperTool_Msg_UpdateVIPUser{
	//	UserData: &netproto.JLBPK_SuperTool_Msg_VIPUserData{
	//		UserID:           proto.Int32(15266910),
	//		Finished:         proto.Int32(0),
	//		Weight:           proto.Int32(100),
	//		TargetPoint:      proto.Int64(0),
	//		AccumulatedPoint: proto.Int64(0),
	//		Status:           proto.Int32(1),
	//	},
	//})
	//
	//rpcMessage := &rpcmsg.RPCMessage{
	//	RPCInfo: rpcInfo,
	//	MsgData: msgData,
	//}
	//
	//m := &msg.Message{
	//	BClassID: int32(netproto.MessageBClassID_Game),
	//	SClassID: int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_UpdateVIPUser),
	//	MsgData:  rpcMessage,
	//}
	//
	//handler := new(SuperToolJLBPKHandler)
	//
	//ctx := new(module.MouduleCtx)
	//clt := new(network.ConnInfo)
	//
	//ret := handler.OnMessage(ctx, clt, m)
	//
	//if len(ret) == 0 {
	//	t.Errorf("No Message Return")
	//}
}