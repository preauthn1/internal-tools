package msghandler

import (
	"testing"

	"example.invalid/golang/mock/gomock"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

func TestGetuseragentlisthandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 10)

	mockProc1 := db.NewProcedure("PrPs_UserAgent_GetUserAgentList", ps) // 測試與資料庫的資料相符

	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)
	mockretset := db.RecordSet{}

	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "totaluser")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "temp_buser")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "todaynewuser")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "todaylogin")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(7), int64(8), int64(9), int64(1)})
	mockretset = append(mockretset, &dataSetResult1)

	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult2)

	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 測試與資料庫的資料相符
	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)         // 測試與資料庫的資料 不 相符

	dbutil.DBUtil = mockDbUtil
	process := new(UserInitBindProcess)
	msgData := new(netproto.GetUserAgentListMsg)
	msgData.UserID = proto.Int32(10)

	rpcInfo := process.getRPCInfo()

	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData)

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentListID), rpcMessage)

	handler := new(GetUserAgentListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_GetUserAgentListID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetUserAgentListRetID) {
		t.Error("OnMessage 錯誤")
	}

	if ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTotalUser() != int32(7) {
		t.Errorf("非7的TotalUser假資料: %v", ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTotalUser())
	}

	if ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTempBUser() != int32(8) {
		t.Errorf("非8的TempBUser假資料: %v", ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTempBUser())
	}

	if ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTodayNewUser() != int32(9) {
		t.Errorf("非9的TodayNewUser假資料: %v", ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTodayNewUser())
	}

	if ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTodayLogin() != int32(1) {
		t.Errorf("非9的TodayLogin假資料: %v", ret[0].MsgData.(*netproto.GetUserAgentListRet).GetTodayLogin())
	}

	t.Log(ret[0].MsgData)
}
