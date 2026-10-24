package msghandler

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

func TestRoomLogoutHandler(t *testing.T) {
	initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.LogoutGameRoomInfo)
	msgdata.UserID = proto.Int32(0)
	msgdata.UpdateAmount = proto.Int64(0)
	msgdata.TaxAmount = proto.Int64(0)
	msgdata.WinCount = proto.Int32(0)
	msgdata.LoseCount = proto.Int32(0)

	msgdata.DrawCount = proto.Int32(0)
	msgdata.LoginTime = proto.String(util.FormatTimeStr(time.Now()))
	msgdata.LogoutTime = proto.String(util.FormatTimeStr(time.Now()))
	msgdata.IPAddress = proto.String("")
	msgdata.HDCode = proto.String("")

	msgdata.HDType = proto.Int32(0)
	msgdata.GameTime = proto.Int32(0)
	msgdata.GameID = proto.Int32(0)
	msgdata.ServerID = proto.Int32(0)
	msgdata.LoginMoney = proto.Int64(0)

	msgdata.ScoreUpdateMode = proto.Int32(0)
	msgdata.TradeId = proto.Int64(0)
	msgdata.LastWinLose = proto.Int32(0)
	msgdata.FlowAmount = proto.Int64(0)
	// msgdata.CouponAmount = proto.Int64(0)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LogoutRoomID), rpcMessage)
	handler := new(RoomLogoutHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	handler.OnMessage(ctx, clt, m)
}
