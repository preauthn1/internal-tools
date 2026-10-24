package msghandler

import (
	"example.invalid/golang/protobuf/proto"
)

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

import (
	"testing"
)



// 無限級預提按鈕
func TestUnlimitConvertOrderHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.UnlimitConvertOrderMsg)
	msgdata.UserID = proto.Int32(56022519)
	msgdata.Type = proto.Int32(0)
	msgdata.ChvIPaddress = proto.String("0.0.0.0")
	msgdata.LngAmount = proto.Int32(-10)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitConvertOrderID), rpcMessage)
	handler := new(UnlimitConvertOrderHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	handler.OnMessage(ctx, clt, m)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UnlimitConvertOrderRetID) {
		t.Error("OnMessage 錯誤")
	}

	mlog.System("查看UnlimitConvertOrderHandler回傳值的結果 : %v", ret[0].MsgData.(*netproto.UnlimitConvertOrderRet).GetChvMsg())
}
