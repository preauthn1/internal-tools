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

func TestDBUnlimitConvertOrderHandler(t *testing.T) {
	// 準備控制器
	mockCtrl := gomock.NewController(t) // mock生態控制器
	defer mockCtrl.Finish()

	// 模擬資料庫的輸入參數
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", 56022519)
	ps.AddIntInput("intType", 3)
	ps.AddVarcharInput("chvIPaddress", "0.0.0.0", 32)
	ps.AddIntInput("lngAmount", 9)
	ps.AddVarcharOutput("chvMsg", "") //chvMsg的輸出

	// 設定觸發mock參數
	mockProc1 := db.NewProcedure("PrPs_Unlimit_ConvertOrder", ps)

	// 產生mock instance
	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

	// 產生結果集
	mockretset := db.RecordSet{}

	// 產生 table1 結果值
	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "chvMsg")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{string("finish")})
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
	msgData := new(netproto.UnlimitConvertOrderMsg)
	msgData.UserID = proto.Int32(56022519)
	msgData.Type = proto.Int32(3)
	msgData.ChvIPaddress = proto.String("0.0.0.0")
	msgData.LngAmount = proto.Int32(9)

	rpcInfo := process.getRPCInfo()
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

	// 產生新的協定訊息
	m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UnlimitConvertOrderID), rpcMessage)

	// 檢查協定HOOK
	handler := new(UnlimitConvertOrderHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("HallMsgClassID_UnlimitConvertOrderIDhook 錯誤")
	}

	// 檢查協定回傳訊息
	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UnlimitConvertOrderRetID) {
		t.Error("HallMsgClassID_UnlimitConvertOrderRetIDOnMessage 錯誤")
	}

	// 查看回傳訊息內容
	if ret[0].MsgData.(*netproto.UnlimitConvertOrderRet).GetChvMsg() != "finish" {
		t.Error("資料讀取錯誤")
	}

}
