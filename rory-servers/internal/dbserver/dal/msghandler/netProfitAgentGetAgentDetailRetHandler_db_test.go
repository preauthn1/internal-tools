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

func TestDBNetProfitAgentGetAgentDetailRetHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 10)
	ps.AddVarcharOutput("chvMsg", "")

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_NetProfitAgent_GetAgentDetail", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table0 結果值
	dataSetResult0 := db.RecordTable{}
	dataSetResult0.Cols = append(dataSetResult0.Cols, "rate")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "sumbonusAmount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "poolamount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "newagenttoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "newagentlast")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "totalagentuser")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "newusertoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "chargelast")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "convertlast")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "bonusamountlast")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "chargetoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "converttoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "bonusamounttoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "payusertoday")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "payuserlast")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "storageamount")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "qrurl")
	dataSetResult0.Cols = append(dataSetResult0.Cols, "withdrawsecs")

	dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(101), int64(102), int64(103), int64(104), int64(105), int64(106), int64(107), int64(108), int64(109), int64(110), int64(111), int64(112), int64(113), int64(114), int64(115), int64(116), string("example.invalid"), int64(118)})
	mockretset = append(mockretset, &dataSetResult0)

	// 產生 table1 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "chvMsg")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("執行完成")})
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
	msgData := new(netproto.NetProfitAgentGetAgentDetailMsg)
	msgData.UserID = proto.Int32(10)
	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailID), rpcMessage)

	// 檢查協定HOOK
	handler := new(NetProfitAgentGetAgentDetailRetHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_NetProfitAgentGetAgentDetailID hook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_NetProfitAgentGetAgentDetailRetID) {
		t.Error("HallMsgClassID_NetProfitAgentGetAgentDetailRetID OnMessage 錯誤")
	}

	// 查看回傳訊息內容
	t.Log(ret[0].MsgData.(*netproto.NetProfitAgentGetAgentDetailRet).GetChvMsg())
}
