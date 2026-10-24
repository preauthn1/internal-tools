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

func TestDBUserAgentBorrowListHandler(t *testing.T) {

	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)
	ps.AddIntInput("intPageIndex", 1)
	ps.AddIntInput("intPageSize", 2)
	ps.AddIntOutput("intPageCount", 0)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyPreConvert", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "bonusamount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "AmountLimit")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "PreConvertAmount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "RemainAmount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "withdrawlastweek")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(101), int64(102), int64(103), int64(104), int64(105)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table1 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "LogDate")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "ID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "Amount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "ConvertType")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "Account")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "Status")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/11/11"), int64(106), int64(107), string("108"), string("109"), string("110")})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/11/12"), int64(111), int64(112), string("113"), string("114"), string("115")})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/11/13"), int64(116), int64(117), string("118"), string("119"), string("120")})
	mockretset = append(mockretset, &dataSetResult1)

	// 產生資料庫的回傳值
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = append(dataSetResult2.Cols, "intPageCount")
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(121)})
	mockretset = append(mockretset, &dataSetResult2)

	// 產生資料庫的回傳值
	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult3)

	// 抽換原有的資料庫
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 興資料相符時抽換
	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)

	// 之前產生的mock instance
	dbutil.DBUtil = mockDbUtil

	// 產生資料庫查詢封包
	process := new(UserInitBindProcess)
	msgData := new(netproto.UserAgentBorrowListMsg)
	msgData.UserID = proto.Int32(56022519)
	msgData.PageIndex = proto.Int32(1)
	msgData.PageSize = proto.Int32(2)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentBorrowListID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UserAgentBorrowListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserAgentBorrowListIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentBorrowListRetID) {
		t.Error("HallMsgClassID_UserAgentBorrowListRetIDOnMessage 錯誤")
	}

	// 查看第1筆daytime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetLogDate() != "2019/11/11" {
		t.Error("第1筆daytime回傳訊息不對，應為2019/11/11，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetLogDate())
	}

	// 查看第1筆borrowid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetID() != 106 {
		t.Error("第1筆borrowid回傳訊息不對，應為106，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetID())
	}

	// 查看第1筆amount回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetAmount() != 107 {
		t.Error("第1筆amount回傳訊息不對，應為107，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetAccount())
	}

	// 查看第1筆type回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetConvertType() != "108" {
		t.Error("第1筆type回傳訊息不對，應為108，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetConvertType())
	}

	// 查看第1筆account回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetAccount() != "109" {
		t.Error("第1筆account回傳訊息不對，應為109，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetAccount())
	}

	// 查看第1筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetStatus() != "110" {
		t.Error("第1筆status回傳訊息不對，應為110，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[0].GetStatus())
	}

	// 查看第2筆daytime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetLogDate() != "2019/11/12" {
		t.Error("第2筆daytime回傳訊息不對，應為2019/11/12，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetLogDate())
	}

	// 查看第2筆borrowid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetID() != 111 {
		t.Error("第2筆borrowid回傳訊息不對，應為111，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetID())
	}

	// 查看第2筆amount回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetAmount() != 112 {
		t.Error("第2筆amount回傳訊息不對，應為107，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetAccount())
	}

	// 查看第2筆type回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetConvertType() != "113" {
		t.Error("第2筆type回傳訊息不對，應為113，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetConvertType())
	}

	// 查看第2筆account回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetAccount() != "114" {
		t.Error("第2筆account回傳訊息不對，應為114，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetAccount())
	}

	// 查看第2筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetStatus() != "115" {
		t.Error("第2筆status回傳訊息不對，應為115，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[1].GetStatus())
	}

	// 查看第3筆daytime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetLogDate() != "2019/11/13" {
		t.Error("第3筆daytime回傳訊息不對，應為2019/11/13，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetLogDate())
	}

	// 查看第3筆borrowid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetID() != 116 {
		t.Error("第3筆borrowid回傳訊息不對，應為116，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetID())
	}

	// 查看第3筆amount回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetAmount() != 117 {
		t.Error("第3筆amount回傳訊息不對，應為117，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetAccount())
	}

	// 查看第3筆type回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetConvertType() != "118" {
		t.Error("第3筆type回傳訊息不對，應為118，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetConvertType())
	}

	// 查看第3筆account回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetAccount() != "119" {
		t.Error("第1筆account回傳訊息不對，應為119，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetAccount())
	}

	// 查看第3筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetStatus() != "120" {
		t.Error("第3筆status回傳訊息不對，應為120，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetUserAgentBorrowList()[2].GetStatus())
	}

	// 查看資料內容 bonus amount
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetBonusamount() != 101 {
		t.Error("bonus amount資料回傳訊息不對，應為101，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetBonusamount())
	}

	// 查看資料內容 amount limit
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetAmountLimit() != 102 {
		t.Error("amount limit 資料回傳訊息不對，應為102，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetAmountLimit())
	}

	// 查看資料內容 pre convert amount
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetPreConvertAmount() != 103 {
		t.Error("pre convert amount 資料回傳訊息不對，應為103，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetPreConvertAmount())
	}

	// 查看資料內容 remain amount
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetRemainAmount() != 104 {
		t.Error("borrow資料回傳訊息不對，應為104，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetRemainAmount())
	}

	// 查看資料內容 withdraw last week
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetWithdrawlastweek() != 105 {
		t.Error("rest資料回傳訊息不對，應為105，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetWithdrawlastweek())
	}

	// 查看資料內容 page count
	if ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetPageCount() != 121 {
		t.Error("page count資料回傳訊息不對，應為121，卻為", ret[0].MsgData.(*netproto.UserAgentBorrowListRet).GetPageCount())
	}
}
