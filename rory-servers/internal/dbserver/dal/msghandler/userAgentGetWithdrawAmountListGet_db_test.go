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

func TestDBUserAgentGetWithdrawAmountListHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)
	ps.AddIntInput("intPageIndex", 1) //第幾頁
	ps.AddIntInput("intPageSize", 2)  //共幾頁
	ps.AddIntOutput("intPageCount", 0)
	ps.AddIntOutput("intSumReward", 0)
	ps.AddIntOutput("intSumConvert", 0)
	ps.AddIntOutput("intRemainConvert", 0)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyConvert", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "StatTime")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "OrderId")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Amount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Type")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Account")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Status")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2019/04/26"), int64(127), int64(128), string("129"), "130", int64(131)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2018/04/26"), int64(132), int64(133), string("134"), "135", int64(136)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2017/04/26"), int64(137), int64(138), string("139"), "140", int64(141)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table1 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intPageCount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intSumReward")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intSumConvert")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intRemainConvert")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(123), int64(124), int64(125), int64(126)})
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
	msgData := new(netproto.UserAgentGetWithdrawAmountListMsg)
	msgData.UserID = proto.Int32(56022519)
	msgData.PageIndex = proto.Int32(1)
	msgData.PageSize = proto.Int32(2)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UserAgentGetWithdrawAmountListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserAgentGetWithdrawAmountListIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentGetWithdrawAmountListRetID) {
		t.Error("HallMsgClassID_UserAgentGetWithdrawAmountListRetIDOnMessage 錯誤")
	}

	// 查看第1筆 Datetime 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetDatetime() != "2019/04/26" {
		t.Error("第1筆 Datetime 資料錯誤，應為2019/04/26，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetDatetime())
	}

	// 查看第1筆 Noid 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetNoid() != 127 {
		t.Error("第1筆 Noid 資料錯誤，應為127，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetNoid())
	}

	// 查看第1筆 Amount 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetAmount() != 128 {
		t.Error("第1筆 Amount 資料錯誤，應為128，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetAmount())
	}

	// 查看第1筆 Type 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetType() != "129" {
		t.Error("第1筆 Type 資料錯誤，應為129，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetType())
	}

	// 查看第1筆 Account 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetAccount() != "130" {
		t.Error("第1筆 Account 資料錯誤，應為130，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetAccount())
	}

	// 查看第1筆 Status 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetStatus() != 131 {
		t.Error("第1筆 Status 資料錯誤，應為131，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[0].GetStatus())
	}

	// 查看第2筆 Datetime 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetDatetime() != "2018/04/26" {
		t.Error("第2筆 Datetime 資料錯誤，應為2018/04/26，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetDatetime())
	}

	// 查看第2筆 Noid 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetNoid() != 132 {
		t.Error("第2筆 Noid 資料錯誤，應為132，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetNoid())
	}

	// 查看第2筆 Amount 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetAmount() != 133 {
		t.Error("第2筆 Amount 資料錯誤，應為133，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetAmount())
	}

	// 查看第2筆 Type 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetType() != "134" {
		t.Error("第2筆 Type 資料錯誤，應為134，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetType())
	}

	// 查看第2筆 Account 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetAccount() != "135" {
		t.Error("第2筆 Account 資料錯誤，應為135，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetAccount())
	}

	// 查看第2筆 Status 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetStatus() != 136 {
		t.Error("第2筆 Status 資料錯誤，應為136，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[1].GetStatus())
	}

	// 查看第3筆 Datetime 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetDatetime() != "2017/04/26" {
		t.Error("第3筆 Datetime 資料錯誤，應為2017/04/26，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetDatetime())
	}

	// 查看第3筆 Noid 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetNoid() != 137 {
		t.Error("第3筆 Noid 資料錯誤，應為137，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetNoid())
	}

	// 查看第3筆 Amount 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetAmount() != 138 {
		t.Error("第3筆 Amount 資料錯誤，應為138，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetAmount())
	}

	// 查看第3筆 Type 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetType() != "139" {
		t.Error("第3筆 Type 資料錯誤，應為139，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetType())
	}

	// 查看第3筆 Account 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetAccount() != "140" {
		t.Error("第3筆 Account 資料錯誤，應為140，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetAccount())
	}

	// 查看第3筆 Status 資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetStatus() != 141 {
		t.Error("第3筆 Status 資料錯誤，應為141，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserAgentGetWithdrawAmountList()[2].GetStatus())
	}

	// 查看userid資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserid() != 56022519 {
		t.Error("userid資料錯誤，應為56022519，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetUserid())
	}

	// 查看pagecount資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetPagecount() != 123 {
		t.Error("pageeount資料錯誤，應為123，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetPagecount())
	}

	// 查看bonus資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetBonus() != 124 {
		t.Error("bonus資料錯誤，應為124，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetBonus())
	}

	// 查看rest資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetRest() != 125 {
		t.Error("rest資料錯誤，應為125，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetRest())
	}

	// 查看allow資料
	if ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetAllow() != 126 {
		t.Error("allow資料錯誤，應為126，卻為", ret[0].MsgData.(*netproto.UserAgentGetWithdrawAmountListRet).GetAllow())
	}

}
