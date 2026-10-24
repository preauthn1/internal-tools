package msghandler

// 參考 C:\Users\20190314user\go\src\0.0.0.0\lioupinchen\GameServer\dbserver\src\dal\msghandler\getuseragentlisthandler_test.go
// 參考 C:\Users\20190314user\go\src\0.0.0.0\lioupinchen\GameServer\dbserver\src\dal\msghandler\fixtures\mockDbUtil.go
// 參考 C:\Users\20190314user\go\src\0.0.0.0\lioupinchen\GameServer\dbserver\src\dal\msghandler\robotlogonhandler_test.go

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

func TestDBNetProfitAgentGetWithdrawDetailHandler(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 10)
	ps.AddVarcharInput("dtStattime", "2019/04/16", 32)

	mockProc1 := db.NewProcedure("PrPs_NetProfitAgent_GetWithdrawDetail", ps) // 測試與資料庫的資料相符

	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)
	mockretset := db.RecordSet{}

	// 準備table0的資炓
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "chargeamountlast")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(100)})
	mockretset = append(mockretset, &dataSetResult0)

	//準備table1的資炓
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "LogDate")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "UserID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "ChargeAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "ConvertAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "StorageAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "BonusAmount")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/04/16"), int64(10), int64(11), int64(12), int64(13), int64(14)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/04/17"), int64(20), int64(21), int64(22), int64(23), int64(24)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/04/18"), int64(30), int64(31), int64(32), int64(33), int64(34)})
	mockretset = append(mockretset, &dataSetResult1)

	//準備DB回傳值
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult2)

	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil) // 測試與資料庫的資料不相符
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 測試與資料庫的資料相符

	dbutil.DBUtil = mockDbUtil
	process := new(UserInitBindProcess)
	msgData := new(netproto.NetProfitAgentGetWithdrawDetailMsg)
	msgData.UserID = proto.Int32(10)
	msgData.Datetime = proto.String("2019/04/16")

	rpcInfo := process.getRPCInfo()

	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData)

	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailID), rpcMessage)

	handler := new(NetProfitAgentGetWithdrawDetailHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_NetProfitAgentGetWithdrawDetailID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailRetID) {
		t.Error("OnMessage 錯誤")
	}

	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetWithdrawDetailRet).GetGetWithdrawDetailList()[0])
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetWithdrawDetailRet).GetGetWithdrawDetailList()[1])
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetWithdrawDetailRet).GetGetWithdrawDetailList()[2])

}
