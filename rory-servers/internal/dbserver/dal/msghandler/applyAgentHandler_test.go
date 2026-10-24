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

func TestApplyAgentHandler(t *testing.T) {
	mockDbUtil := fixtures.MockDBUtil(t)
	mockRet := fixtures.GetExpectApplyAgentRet()
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockRet, nil)
	dbutil.DBUtil = mockDbUtil
	dbutil.DealLogDBUtil = mockDbUtil
	//initDBConnection()
	rpcInfo := netproto.CreateEmptyRPCInfo(1)
	rpcInfo.UserID = proto.Int32(1)
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo

	msgdata := new(netproto.ApplyAgent)
	msgdata.AreaID = proto.Int32(1)
	msgdata.Name = proto.String("test_name")
	msgdata.Tel = proto.String("test_tel")
	msgdata.QQ = proto.String("test_qq")
	msgdata.WXNo = proto.String("test_wxno")
	msgdata.Memo = proto.String("test_memo")
	msgdata.Alipay = proto.String("test_alipay")
	msgData, _ := proto.Marshal(msgdata)
	rpcMessage.MsgData = msgData
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ApplyAgentID), rpcMessage)
	handler := new(ApplyAgentHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("hook 錯誤")
	}
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_ApplyAgentRetID) {
		t.Error("OnMessage 錯誤")
	}
}
