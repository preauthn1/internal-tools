// handler單元測試
package msghandler

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

//遊戲向服務器註冊
func TestServerRegHandler(t *testing.T) {
	InitTest()

	// 驗證 handler
	serverid := int32(365)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	// rpcInfo.Cer = proto.String(util.ToString(accessToken))
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	// 驗證資料
	msgdata := new(netproto.ServerRegInfo)
	msgdata.ServerID = proto.Int32(serverid)
	msgdata.ServerType = proto.Int32(0)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_ServerMgm), int32(netproto.ServerMgmClassID_GameServerRegister), rpcMessage)
	handler := new(ServerRegHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	handler.OnMessage(TestCtx, clt, m)

}
