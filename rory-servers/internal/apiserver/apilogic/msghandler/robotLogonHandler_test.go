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

func TestRobotLogonHandler(t *testing.T) {
	// mockDbUtil := fixtures.MockDBUtil(t)
	// mockRet := fixtures.GetExpectLogonRobotRet()
	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil)
	// dbutil.DbUtil = mockDbUtil
	InitTest()

	// 取得 token
	var accessToken = [REDACTED]

	serverid := int32(102300410)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid)
	rpcInfo.Cer = proto.String(util.ToString(accessToken))
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	logonRobotMsgData := new(netproto.LogonRobot)
	logonRobotMsgData.CurrentRobotCount = proto.Int32(1)
	logonRobotMsgData.MinMoneyAmount = proto.Int64(100000)
	logonRobotMsgData.MaxMoneyAmount = proto.Int64(500000000)
	logonRobotMsgData.ServerID = proto.Int32(serverid)
	msgData, _ := proto.Marshal(logonRobotMsgData)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_DBServer), int32(netproto.DBServerClassID_LogonRobotID), rpcMessage)

	handler := new(RobotLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}

	clt := new(network.ConnInfo)
	ret := handler.OnMessage(TestCtx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_LogonRobotRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
