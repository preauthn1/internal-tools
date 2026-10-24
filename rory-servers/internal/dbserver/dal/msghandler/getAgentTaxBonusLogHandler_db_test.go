package msghandler

import (
	"testing"
	"time"

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

func TestDBGetAgentTaxBonusLogHandler(t *testing.T) {

	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)

	// 設定觸發mock參數
	// mockProc1 := db.NewProcedure("PrPs_UserAgent_GetAgentTaxBonusLog", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "LogDate")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "B_TaxBonusAmount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "CD_TaxBonusAmount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "TotalTaxBonusAmount")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{time.Date(2019, 4, 21, 12, 34, 56, 59, time.UTC), int64(200), int64(300), int64(400)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{time.Date(2018, 5, 22, 3, 5, 6, 7, time.UTC), int64(500), int64(600), int64(700)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{time.Date(2017, 10, 25, 23, 30, 16, 75, time.UTC), int64(800), int64(900), int64(1000)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生資料庫的回傳值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult1)

	// 抽換原有的資料庫
	// mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 興資料相符時抽換
	mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)

	// 之前產生的mock instance
	dbutil.DBUtil = mockDbUtil

	// 產生資料庫查詢封包
	process := new(UserInitBindProcess)
	msgData := new(netproto.GetAgentTaxBonusLogMsg)
	msgData.UserID = proto.Int32(56022519)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxBonusLogID), rpcMessage)

	// 檢查協定HOOK
	handler := new(GetAgentTaxBonusLogHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogID hook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	// 檢查協定OnMessage
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetAgentTaxBonusLogRetID) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID OnMessage 錯誤")
	}

	// 檢查第一筆LogDate輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetLogDate() != "2019-04-21 12:34:56" {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID LogDate輸出 錯誤(非2019-04-21 12:34:56)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetLogDate())
	}

	// 檢查第一筆B_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetBTaxBonusAmount() != int32(200) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID B_TaxBonusAmount輸出 錯誤(非200)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetBTaxBonusAmount())
	}

	// 檢查第一筆CD_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetCDTaxBonusAmount() != int32(300) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID CD_TaxBonusAmount輸出 錯誤(非300)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetCDTaxBonusAmount())
	}

	// 檢查第一筆TotalTaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetTotalTaxBonusAmount() != int32(400) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID TotalTaxBonusAmount輸出 錯誤(非400)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetTotalTaxBonusAmount())
	}

	// 檢查第二筆LogDate輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetLogDate() != "2018-05-22 03:05:06" {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID LogDate輸出 錯誤(非2018-05-22 03:05:06)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[0].GetLogDate())
	}

	// 檢查第二筆B_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetBTaxBonusAmount() != int32(500) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID B_TaxBonusAmount輸出 錯誤(非500)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetBTaxBonusAmount())
	}

	// 檢查第二筆CD_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetCDTaxBonusAmount() != int32(600) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID CD_TaxBonusAmount輸出 錯誤(非600)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetCDTaxBonusAmount())
	}

	// 檢查第二筆TotalTaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetTotalTaxBonusAmount() != int32(700) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID TotalTaxBonusAmount輸出 錯誤(非700)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[1].GetTotalTaxBonusAmount())
	}

	// 檢查第三筆LogDate輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetLogDate() != "2017-10-25 23:30:16" {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID LogDate輸出 錯誤(非2017-10-25 23:30:16)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetLogDate())
	}

	// 檢查第三筆B_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetBTaxBonusAmount() != int32(800) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID B_TaxBonusAmount輸出 錯誤(非800)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetBTaxBonusAmount())
	}

	// 檢查第三筆CD_TaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetCDTaxBonusAmount() != int32(900) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID CD_TaxBonusAmount輸出 錯誤(非900)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetCDTaxBonusAmount())
	}

	// 檢查第三筆TotalTaxBonusAmount輸出是否正確
	if ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetTotalTaxBonusAmount() != int32(1000) {
		t.Error("HallMsgClassID_GetAgentTaxBonusLogRetID TotalTaxBonusAmount輸出 錯誤(非1000)", ret[0].MsgData.(*netproto.GetAgentTaxBonusLogRet).GetTaxBonusLogList()[2].GetTotalTaxBonusAmount())
	}
}
