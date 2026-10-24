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

func TestDBUserAgentGetWithAmountListHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)
	ps.AddIntInput("intTargetID", 56022519)
	ps.AddIntInput("intPageIndex", 1)
	ps.AddIntInput("intPageSize", 2)
	ps.AddIntOutput("intPageCount", 0)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyMember", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "userid")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "nickname")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekTeamFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "WeekMemberFlows")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "TeamMember")
	// dataSetResult0.Cols = append(dataSetResult0.Cols, "Agent")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "Member")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "MemberType")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(100), string("150"), int64(200), int64(300), int64(400), int64(600), int64(700)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(101), string("151"), int64(201), int64(301), int64(401), int64(601), int64(701)})
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(102), string("152"), int64(202), int64(302), int64(402), int64(602), int64(702)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table0 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "intPageCount")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(800)})
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
	msgData := new(netproto.UserAgentGetWithAmountListMsg)
	msgData.UserID = proto.Int32(56022519)
	msgData.TargetID = proto.Int32(56022519)
	msgData.PageIndex = proto.Int32(1)
	msgData.PageSize = proto.Int32(2)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountListID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UserAgentGetWithAmountListHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UserAgentGetWithAmountListIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)

	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentGetWithAmountListRetID) {
		t.Error("HallMsgClassID_UserAgentGetWithAmountListRetIDOnMessage 錯誤")
	}

	// 查看第1筆userid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetUserid() != 100 {
		t.Error("第1筆userid錯誤，非56022519，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetUserid())
	}

	// 查看第1筆groupoffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetGroupoffer() != 200 {
		t.Error("第1筆groupoffer錯誤，非200，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetGroupoffer())
	}

	// 查看第1筆playeroffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetPlayeroffer() != 300 {
		t.Error("第1筆playeroffer錯誤，非300，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetPlayeroffer())
	}

	// 查看第1筆membernumbers回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetMembernumbers() != 400 {
		t.Error("第1筆membernumbers錯誤，非400，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetMembernumbers())
	}

	/*
		// 查看第1筆agent回傳訊息內容
		if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetAgent() != 500 {
			t.Error("第1筆agent錯誤，非500，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetAgent())
		}
	*/

	// 查看第1筆belongs回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetBelongs() != 600 {
		t.Error("第1筆belongs錯誤，非600，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetBelongs())
	}

	// 查看第1筆types回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetType() != 700 {
		t.Error("第1筆type錯誤，非700，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetType())
	}

	// 查看第2筆userid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetUserid() != 101 {
		t.Error("第2筆userid錯誤，非56022519，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetUserid())
	}

	// 查看第2筆groupoffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetGroupoffer() != 201 {
		t.Error("第2筆groupoffer錯誤，非201，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetGroupoffer())
	}

	// 查看第2筆playeroffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetPlayeroffer() != 301 {
		t.Error("第2筆playeroffer錯誤，非301，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetPlayeroffer())
	}

	// 查看第2筆membernumbers回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetMembernumbers() != 401 {
		t.Error("第2筆membernumbers錯誤，非401，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetMembernumbers())
	}

	/*
		// 查看第2筆agent回傳訊息內容
		if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetAgent() != 501 {
			t.Error("第2筆agent錯誤，非501，為", ret[1].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[0].GetAgent())
		}
	*/

	// 查看第2筆belongs回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetBelongs() != 601 {
		t.Error("第2筆belongs錯誤，非601，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetBelongs())
	}

	// 查看第2筆types回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetType() != 701 {
		t.Error("第2筆type錯誤，非701，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[1].GetType())
	}

	// 查看第3筆userid回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetUserid() != 102 {
		t.Error("第3筆userid錯誤，非56022519，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetUserid())
	}

	// 查看第3筆groupoffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetGroupoffer() != 202 {
		t.Error("第3筆groupoffer錯誤，非202，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetGroupoffer())
	}

	// 查看第3筆playeroffer回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetPlayeroffer() != 302 {
		t.Error("第3筆playeroffer錯誤，非302，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetPlayeroffer())
	}

	// 查看第3筆membernumbers回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetMembernumbers() != 402 {
		t.Error("第3筆membernumbers錯誤，非402，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetMembernumbers())
	}

	/*
		// 查看第3筆agent回傳訊息內容
		if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetAgent() != 502 {
			t.Error("第3筆agent錯誤，非502，為", ret[2].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetAgent())
		}
	*/

	// 查看第3筆belongs回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetBelongs() != 602 {
		t.Error("第3筆belongs錯誤，非602，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetBelongs())
	}

	// 查看第3筆types回傳訊息內容
	if ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetType() != 702 {
		t.Error("第3筆type錯誤，非702，為", ret[0].MsgData.(*netproto.UserAgentGetWithAmountListRet).GetUserAgentGetWithAmountList()[2].GetType())
	}

}
