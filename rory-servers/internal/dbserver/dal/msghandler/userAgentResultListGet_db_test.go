// handler資料庫單元測試
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

func TestDBUserAgentResultListHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)
	ps.AddIntInput("intPageIndex", 1) //第幾頁
	ps.AddIntInput("intPageSize", 2)  //共幾頁
	ps.AddIntOutput("intPageCount", 0)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyPerformance", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekMyFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekTeamFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekMemberFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Reward")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(1000), int64(900), int64(800), int64(700)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table1 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "StatTime")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "WeekMyFlows")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "WeekTeamFlows")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "WeekMemberFlows")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2019/10/10"), int64(600), int64(550), int64(500)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2018/10/10"), int64(450), int64(400), int64(350)})
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("2017/10/10"), int64(300), int64(250), int64(200)})
	mockretset = append(mockretset, &dataSetResult1)

	// 產生 table2 結果值
	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = append(dataSetResult2.Cols, "intPageCount")
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{int64(25)})
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
	msgData := new(netproto.UserAgentResultListMsg)
	msgData.UserID = proto.Int32(56022519)
	msgData.PageIndex = proto.Int32(1)
	msgData.PageSize = proto.Int32(2)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentResultListID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UserAgentResultListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserAgentResultListIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentResultListRetID) {
		t.Error("HallMsgClassID_UserAgentResultListRetIDOnMessage 錯誤")
	}

	// 查看第1筆回傳訊息內容time
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetTime() != "2019/10/10" {
		t.Error("第1筆 time資料錯誤，應為2019/10/10，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetTime())
	}

	// 查看第1筆回傳訊息內容dayresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDayresult() != 600 {
		t.Error("第1筆 dayresult資料錯誤，應為600，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDayresult())
	}

	// 查看第1筆回傳訊息內容daygroupresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDaygroupresult() != 550 {
		t.Error("第1筆 daygroupresult資料錯誤，應為550，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDaygroupresult())
	}

	// 查看第1筆回傳訊息內容daydirectresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDaydirectresult() != 500 {
		t.Error("第1筆 daydirectresult資料錯誤，應為550，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[0].GetDaydirectresult())
	}

	// 查看第2筆回傳訊息內容time
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetTime() != "2018/10/10" {
		t.Error("第2筆 time資料錯誤，應為2018/10/10，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetTime())
	}

	// 查看第2筆回傳訊息內容dayresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDayresult() != 450 {
		t.Error("第2筆 dayresult資料錯誤，應為450，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDayresult())
	}

	// 查看第2筆回傳訊息內容daygroupresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDaygroupresult() != 400 {
		t.Error("第2筆 daygroupresult資料錯誤，應為400，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDaygroupresult())
	}

	// 查看第2筆回傳訊息內容daydirectresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDaydirectresult() != 350 {
		t.Error("第2筆 daydirectresult資料錯誤，應為350，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[1].GetDaydirectresult())
	}
	// 查看第3筆回傳訊息內容time
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetTime() != "2017/10/10" {
		t.Error("第3筆 time資料錯誤，應為2017/10/10，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetTime())
	}

	// 查看第3筆回傳訊息內容dayresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDayresult() != 300 {
		t.Error("第3筆 dayresult資料錯誤，應為300，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDayresult())
	}

	// 查看第3筆回傳訊息內容daygroupresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDaygroupresult() != 250 {
		t.Error("第3筆 daygroupresult資料錯誤，應為250，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDaygroupresult())
	}

	// 查看第3筆回傳訊息內容daydirectresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDaydirectresult() != 200 {
		t.Error("第3筆 daydirectresult資料錯誤，應為200，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetUserAgentResultList()[2].GetDaydirectresult())
	}

	// 查看資料result
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetResult() != 1000 {
		t.Error("資料 Result，應為100，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetResult())
	}

	// 查看資料groupresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetGroupresult() != 900 {
		t.Error("資料 groupresult，應為900，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetGroupresult())
	}

	// 查看資料directresult
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetDirectresult() != 800 {
		t.Error("資料 directresult，應為800，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetDirectresult())
	}

	// 查看資料bonus
	if ret[0].MsgData.(*netproto.UserAgentResultListRet).GetBonus() != 700 {
		t.Error("資料 bonus，應為700，卻為", ret[0].MsgData.(*netproto.UserAgentResultListRet).GetBonus())
	}
}
