// handler單元測試
package msghandler

import (
	"testing"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 登出房间
func TestRoomLogoutHandler(t *testing.T) {
	InitTest()

	// 取得 token
	var accessToken = [REDACTED]

	// 驗證 handler
	serverid := int32(365)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcInfo.Cer = proto.String(util.ToString(accessToken))
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	// 驗證資料
	msgdata := new(netproto.LogoutGameRoomInfo)
	msgdata.UserID = proto.Int32(1000000)
	msgdata.UpdateAmount = proto.Int64(40)
	msgdata.TaxAmount = proto.Int64(0)
	msgdata.WinCount = proto.Int32(0)
	msgdata.LoseCount = proto.Int32(0)

	msgdata.DrawCount = proto.Int32(0)
	msgdata.LoginTime = proto.String(time.Now().Format("2006-01-02T15:04:05Z"))
	msgdata.LogoutTime = proto.String(time.Now().Format("2006-01-02T15:04:05Z"))
	msgdata.IPAddress = proto.String("0.0.0.0")
	msgdata.HDCode = proto.String("")

	msgdata.HDType = proto.Int32(0)
	msgdata.GameTime = proto.Int32(0)
	msgdata.GameID = proto.Int32(1023)
	msgdata.ServerID = proto.Int32(serverid)
	msgdata.LoginMoney = proto.Int64(0)

	msgdata.ScoreUpdateMode = proto.Int32(0)
	msgdata.TradeId = proto.Int64(20200100000064745)
	msgdata.LastWinLose = proto.Int32(0)
	msgdata.FlowAmount = proto.Int64(0)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LogoutRoomID), rpcMessage)
	handler := new(RoomLogoutHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Error("HallServerStart OnMessage 錯誤")
	}

}
