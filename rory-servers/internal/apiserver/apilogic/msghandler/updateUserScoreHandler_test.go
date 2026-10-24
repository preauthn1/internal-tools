// handler單元測試
package msghandler

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 寫分
func TestUpdateUserScoreHandler(t *testing.T) {
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
	msgdata := new(netproto.UpdateUserScore)
	msgdata.UserID = proto.Int32(1000014)
	msgdata.ServerID = proto.Int32(serverid)

	msgdata.Score = proto.Int64(100)
	msgdata.TaxAmount = proto.Int64(0)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_UpdateUserScoreID), rpcMessage)
	handler := new(UpdateUserScoreHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_DBRetID) {
		t.Error("OnMessage 錯誤")
	}

}
