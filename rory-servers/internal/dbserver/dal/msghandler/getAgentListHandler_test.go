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

func TestGetAgentListHandler(t *testing.T) {
	mockDbUtil := fixtures.MockDBUtil(t)
	mockRet := fixtures.GetExpectAgentListRet()
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil)
	dbutil.DBUtil = mockDbUtil
	dbutil.DealLogDBUtil = mockDbUtil
	//initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.RequestAgentList)
	msgdata.AreaID = proto.Int32(1)
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentListID), rpcMessage)
	handler := new(GetAgentListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetAgentListRetID) {
		t.Error("OnMessage 錯誤")
	}
}
