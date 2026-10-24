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

func TestDBGetMyPromotionHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_GetMyPromotion", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "BossID")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "TeamMember")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "MyMember")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "TodayNewMember")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "TodayActiveMember")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "QrURL")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "invitecode")
	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(321), int64(320), int64(319), int64(318), int64(317), string("example.invalid"), string("")})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table0 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "BossID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "TeamMember")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "MyMember")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "TodayNewMember")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "TodayActiveMember")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "QrURL")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "invitecode")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(321), int64(320), int64(319), int64(318), int64(317), string("example.invalid"), string("")})
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
	msgData := new(netproto.GetMyPromotionMsg)
	msgData.UserID = proto.Int32(56022519)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetMyPromotionID), rpcMessage)

	// 檢查協定HOOK
	handler := new(GetMyPromotionHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_GetMyPromotionIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_GetMyPromotionRetID) {
		t.Error("HallMsgClassID_GetMyPromotionRetIDOnMessage 錯誤")
	}

	// 查看回傳訊息內容
	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetBossID() != 321 {
		t.Error("BossID值應為321卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetBossID())
	}

	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTeamMember() != 320 {
		t.Error("TeamMember值應為320卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTeamMember())
	}

	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetMyMember() != 319 {
		t.Error("MyMember值應為319卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetMyMember())
	}

	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTodayNewMember() != 318 {
		t.Error("TodayNewMember值應為318卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTodayNewMember())
	}

	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTodayActiveMember() != 317 {
		t.Error("TodayActiveMember值應為317卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTodayActiveMember())
	}

	if ret[0].MsgData.(*netproto.GetMyPromotionRet).GetQrURL() != "example.invalid" {
		t.Error("QrURL值應為www.123.com卻為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetTodayActiveMember())
	}

	t.Log("現在系統間為", ret[0].MsgData.(*netproto.GetMyPromotionRet).GetSystemTime())
}
