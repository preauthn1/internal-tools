// handler單元測試
package msghandler

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 遊戲進房
func TestRoomLogonHandler(t *testing.T) {
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
	msgdata := new(netproto.LoginGameRoomInfo)
	msgdata.Cer = proto.String(accessToken)
	msgdata.UserID = proto.Int32(1000000)
	msgdata.HDCode = proto.String("")
	msgdata.HDType = proto.Int32(0)

	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData

	m := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomID), rpcMessage)
	handler := new(RoomLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_GameRoom) || ret[0].SClassID != int32(netproto.GameRoomClassID_LoginRoomRetID) {
		t.Error("OnMessage 錯誤")
	}
	data := ret[0].MsgData.(*netproto.LoginGameRoomRet)
	fmt.Println(data.GetCode(), data.GetMessage())
	if data.GetCode() != 0 {
		t.Errorf("Code want %d, but %d", 0, data.GetCode())
	}
}
