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

func TestDBUserAgentWinAmountListHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyReward", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "StatTime")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekMyFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekTeamFLows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekMemberFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Reward")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Status")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2016/03/04"), int64(1002), int64(1003), int64(1004), int64(1005), int64(1006)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2017/05/06"), int64(2002), int64(2003), int64(2004), int64(2005), int64(2006)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{string("2018/07/08"), int64(3002), int64(3003), int64(3004), int64(3005), int64(3006)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生資料庫的回傳值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult1)

	// 抽換原有的資料庫
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 興資料相符時抽換
	// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)

	// 之前產生的mock instance
	dbutil.DBUtil = mockDbUtil

	// 產生資料庫查詢封包
	process := new(UserInitBindProcess)
	msgData := new(netproto.UserAgentWinAmountListMsg)
	msgData.UserID = proto.Int32(56022519)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentWinAmountListID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UserAgentWinAmountListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserAgentWinAmountListIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentWinAmountListRetID) {
		t.Error("HallMsgClassID_UserAgentWinAmountListRetIDOnMessage 錯誤")
	}

	// 查看第1筆sendtime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetSendTime() != "2016/03/04" {
		t.Error("第1筆sendtime資料錯誤，應為2016/03/04，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetSendTime())
	}

	// 查看第1筆weekresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekresult() != 1002 {
		t.Error("第1筆weekresult資料錯誤，應為1002，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekresult())
	}

	// 查看第1筆weekgroupresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekgroupresult() != 1003 {
		t.Error("第1筆weekgroupresult資料錯誤，應為1003，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekgroupresult())
	}

	// 查看第1筆weekdirectresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekdirectresult() != 1004 {
		t.Error("第1筆weekdirectresult資料錯誤，應為1004，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetWeekdirectresult())
	}

	// 查看第1筆bonus回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetBonus() != 1005 {
		t.Error("第1筆bonus資料錯誤，應為1005，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetBonus())
	}

	// 查看第1筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetStatus() != 1006 {
		t.Error("第1筆status資料錯誤，應為1006，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[0].GetStatus())
	}

	// 查看第2筆sendtime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetSendTime() != "2017/05/06" {
		t.Error("第2筆sendtime資料錯誤，應為2001，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetSendTime())
	}

	// 查看第2筆weekresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekresult() != 2002 {
		t.Error("第2筆weekresult資料錯誤，應為2002，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekresult())
	}

	// 查看第2筆weekgroupresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekgroupresult() != 2003 {
		t.Error("第2筆weekgroupresult資料錯誤，應為2003，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekgroupresult())
	}

	// 查看第2筆weekdirectresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekdirectresult() != 2004 {
		t.Error("第2筆weekdirectresult資料錯誤，應為2004，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetWeekdirectresult())
	}

	// 查看第2筆bonus回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetBonus() != 2005 {
		t.Error("第2筆bonus資料錯誤，應為2005，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetBonus())
	}

	// 查看第2筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetStatus() != 2006 {
		t.Error("第2筆status資料錯誤，應為2006，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[1].GetStatus())
	}

	// 查看第3筆sendtime回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetSendTime() != "2018/07/08" {
		t.Error("第3筆sendtime資料錯誤，應為3001，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetSendTime())
	}

	// 查看第3筆weekresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekresult() != 3002 {
		t.Error("第3筆weekresult資料錯誤，應為3002，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekresult())
	}

	// 查看第3筆weekgroupresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekgroupresult() != 3003 {
		t.Error("第3筆weekgroupresult資料錯誤，應為3003，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekgroupresult())
	}

	// 查看第3筆weekdirectresult回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekdirectresult() != 3004 {
		t.Error("第3筆weekdirectresult資料錯誤，應為3004，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetWeekdirectresult())
	}

	// 查看第3筆bonus回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetBonus() != 3005 {
		t.Error("第3筆bonus資料錯誤，應為3005，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetBonus())
	}

	// 查看第3筆status回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetStatus() != 3006 {
		t.Error("第3筆status資料錯誤，應為3006，但為", ret[0].MsgData.(*netproto.UserAgentWinAmountListRet).GetUserAgentWinAmount()[2].GetStatus())
	}

}
