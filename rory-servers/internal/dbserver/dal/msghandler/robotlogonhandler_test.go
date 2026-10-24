package msghandler

import (
	"testing"

	"example.invalid/golang/mock/gomock"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/msghandler/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestRobotLogonHandler(t *testing.T) {
	mockDbUtil := fixtures.MockDBUtil(t)
	mockRet := fixtures.GetExpectLogonRobotRet()
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil)
	dbutil.DBUtil = mockDbUtil

	serverid := int32(80)
	rpcInfo := netproto.CreateEmptyRPCInfo(serverid) //bjl 3
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
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_DBServer) || ret[0].SClassID != int32(netproto.DBServerClassID_LogonRobotRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)
}
