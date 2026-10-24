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

func TestRoomLogonHandler(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	ps := db.NewSQLParameters()
	ps.AddIntInput("sintServerID", 1)
	ps.AddIntInput("intUserID", 1021234)
	ps.AddIntInput("tnyHDType", 1)
	ps.AddVarcharInput("chvHDCode", "", 64)
	ps.AddVarcharInput("chvIPAddress", "0.0.0.0", 15)
	ps.AddVarcharInput("chvCer", "", 32)
	ps.AddVarcharOutput("chvErrMsg", "")
	ps.AddBigIntOutput("lngTradeID", 0)
	mockProc1 := db.NewProcedure("PrGs_UserLogin", ps)

	mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)
	mockretset := db.RecordSet{}

	dataSetResult1 := db.RecordTable{}
	dataSetResult1.Cols = append(dataSetResult1.Cols, "UserID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "UserType")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "NickName")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "FaceID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "Sex")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "CashAmount")
	// dataSetResult1.Cols = append(dataSetResult1.Cols, "CouponAmount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "WinCount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "LoseCount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "DrawCount")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "IfOffLine")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "ServerAddr")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "GameBuff")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "TradeID")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "LastWinLose")
	dataSetResult1.Cols = append(dataSetResult1.Cols, "TimeSeq")
	dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(1021234), int64(1), "來自火星", int64(13), int64(1), int64(2600), int64(0), int64(0), int64(0), int64(1), "", "", int64(2018050000000753), int64(0), int64(0)})
	mockretset = append(mockretset, &dataSetResult1)

	dataSetResult2 := db.RecordTable{}
	dataSetResult2.Cols = append(dataSetResult2.Cols, "chvErrMsg")
	dataSetResult2.Cols = append(dataSetResult2.Cols, "lngTradeID")
	dataSetResult2.Rows = append(dataSetResult2.Rows, db.RecordRow{"", int64(2018050000000753)})
	mockretset = append(mockretset, &dataSetResult2)
	dataSetResult3 := db.RecordTable{}
	dataSetResult3.Rows = append(dataSetResult3.Rows, db.RecordRow{int64(1)})
	mockretset = append(mockretset, &dataSetResult3)
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil)

	ps2 := db.NewSQLParameters()
	ps2.AddIntInput("intUserID", 1021234)
	ps2.AddVarcharOutput("chvErrMsg", "")
	mockProc2 := db.NewProcedure("PrGs_User_GetUserInfo", ps2)
	dataSetResult4 := db.RecordTable{}
	dataSetResult4.Cols = append(dataSetResult4.Cols, "LoginName")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "NickName")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UserType")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "Sex")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "HeadID")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "CashAmount")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "BankAmount")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "SumCharge")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "SumConvert")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "ShowConvertLimit")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "ShowAllGameLimit")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "IsGaming")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "IsBindGuest")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "IsBindZhifubao")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "AlipayAccount")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "RealName")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "IsBindUnionpayBank")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UnionpayBank")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UnionpayAccount")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UnionpayRealName")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UnionpayBankName")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "UserLevel")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "ExperienceCash")
	dataSetResult4.Cols = append(dataSetResult4.Cols, "AvatarFrame")
	dataSetResult4.Rows = append(dataSetResult4.Rows, db.RecordRow{"", "", 3, 0, 1, 3000, 0, int64(0), int64(0), int64(0), 0, 0, 0, int64(0), "", "", int64(0), "", "", "", "", 0, int64(10000), int64(1)})
	mockretset2 := db.RecordSet{}
	mockretset2 = append(mockretset2, &dataSetResult4)
	dataSetResult5 := db.RecordTable{}
	dataSetResult5.Cols = append(dataSetResult5.Cols, "chvErrMsg")
	dataSetResult5.Rows = append(dataSetResult5.Rows, db.RecordRow{""})
	mockretset2 = append(mockretset2, &dataSetResult5)
	mockretset2 = append(mockretset2, &dataSetResult3)
	mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc2)).Return(mockretset2, nil)

	dbutil.DBUtil = mockDbUtil
	//initDBConnection()
	process := new(UserInitBindProcess)
	//Client傳入資料的協議號
	msgData := new(netproto.LoginGameRoomInfo)
	msgData.UserID = proto.Int32(1021234)
	msgData.HDCode = proto.String("")
	msgData.HDType = proto.Int32(1)
	msgData.Cer = proto.String("")

	//rpc資料必填欄位
	rpcInfo := process.getRPCInfo()

	//rpc詳細資料, MsgData要放入Marshal後的資料
	rpcMessage := new(rpcmsg.RPCMessage)
	rpcMessage.RPCInfo = rpcInfo
	rpcMessage.MsgData, _ = proto.Marshal(msgData) //MsgData要放入Marshal後的資料

	m := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomID), rpcMessage)

	handler := new(RoomLogonHandler)
	if !handler.IsHook(m.BClassID, m.SClassID) {
		t.Error("GameRoomClassID_LoginRoomID hook 錯誤")
	}

	ctx := new(module.MouduleCtx)
	clt := new(network.ConnInfo)
	ret := handler.OnMessage(ctx, clt, m)
	if ret[0].BClassID != int32(netproto.MessageBClassID_GameRoom) || ret[0].SClassID != int32(netproto.GameRoomClassID_LoginRoomRetID) {
		t.Error("OnMessage 錯誤")
	}
	t.Log(ret[0].MsgData)

}
