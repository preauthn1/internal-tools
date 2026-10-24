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
// 	"testing"
// )

// func TestCampaignGetListHandler(t *testing.T) {
// 	initDBConnection()
// 	rpcInfo := netproto.CreateEmptyRPCInfo(1)
// 	rpcMessage := new(rpcmsg.RPCMessage)
// 	rpcMessage.RPCInfo = rpcInfo

// 	msgdata := new(netproto.CampaignGetList)
// 	msgdata.UserID = proto.Int32(1021470)
// 	msgData, _ := proto.Marshal(msgdata)
// 	rpcMessage.MsgData = msgData
// 	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_CampaignGetListID), rpcMessage)
// 	handler := new(CampaignGetListHandler)
// 	if !handler.IsHook(m.BClassID, m.SClassID) {
// 		t.Error("hook 錯誤")
// 	}
// 	ctx := new(module.MouduleCtx)
// 	clt := new(network.ConnInfo)
// 	ret := handler.OnMessage(ctx, clt, m)
// 	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_CampaignGetListRetID) {
// 		t.Error("OnMessage 錯誤")
// 	}
// }

//test parse to time
//func TestParseCustomTimeString(t *testing.T) {
//	time := parseCustomTimeString("")
//	t.Log("time:", time)
//}
