package msghandler

// 參考 C:\Users\20190314user\go\src\0.0.0.0\lioupinchen\GameServer\dbserver\src\dal\msghandler\netProfitAgentGetWithdrawDetail_db_test.go

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

func TestDBgetAgentReportLogHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 10)
	ps.AddVarcharInput("dtStattime", "2019/04/16", 32)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_NetProfitAgent_GetAgentReportLog", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "totallastdate")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "totalagentuser")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(8), int64(9)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table1 結果值，為LIST資料
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "userID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "CreateTime") // 這次時間是回傳字串
	dataSetResult1.Cols = append(dataSetResult1.Cols, "SumSubUsers")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "SumChargeAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "WithdrawAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "Counts")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(10), string("2019/04/17"), int64(11), int64(12), int64(13), int64(14)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(15), string("2019/04/18"), int64(16), int64(17), int64(18), int64(19)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(20), string("2019/04/19"), int64(21), int64(22), int64(23), int64(24)})
	mockretset = append(mockretset, &dataSetResult1)

	// 產生資料庫的回傳值
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult2)

	// 抽換原有的資料庫
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 興資料相符時抽換
	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)

	// 之前產生的mock instance
	dbutil.DBUtil = mockDbUtil

	// 產生資料庫查詢封包
	process := new(UserInitBindProcess)
	msgData := new(netproto.NetProfitAgentGetAgentReportLogMsg)
	msgData.UserID = proto.Int32(10)
	msgData.Datetime = proto.String("2019/04/16")
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogID), rpcMessage)

	// 檢查協定HOOK
	handler := new(NetProfitAgentGetAgentReportLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_NetProfitAgentGetAgentReportLogID hook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogRetID) {
		t.Error("HallMsgClassID_NetProfitAgentGetAgentReportLogRetID OnMessage 錯誤")
	}

	// 查看回傳訊息內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetTotallastdate())                            // 一定要回傳的資料
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetTotalagentuser())                           // 一定要回傳的資料
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[0].GetUserid())     // List資料內的資料內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[0].GetCreateTime()) // List資料內的資料內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[1].GetUserid())     // List資料內的資料內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[1].GetCreateTime()) // List資料內的資料內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[2].GetUserid())     // List資料內的資料內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentReportLogRet).GetGetWithdrawDetailList()[2].GetCreateTime()) // List資料內的資料內容
}
