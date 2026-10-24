package msghandler

import "testing"

func TestDBUserAgentGetWithAmountHandler(t *testing.T) {
	// SP 準備中
	/*
		// 準備控制器
		mockCtrl := gomock.NewController(t) // mock生態控制器
		defer mockCtrl.Finish()

		// 模擬資料庫的輸入參數
		ps := db.NewSqlParameters()
		ps.AddIntInput("intUserID", 56022519)
		ps.AddIntOutput("intReward", 0)             // 截止本日零點預估現金獎勵
		ps.AddIntOutput("intConvert", 0)            // 可提現金額 : 當前現金獎勵x90% - 已預提過的金額
		ps.AddIntOutput("intSumConvert", 0)         // 已預提現金額 : 本周(周二~周日)內預提過的金額累計
		ps.AddIntOutput("intRemainConvert", 0)      // 剩餘可預提現金額 : 可預提金額 – 已預提金額
		ps.AddIntOutput("intLastWeekSumConvert", 0) // 上周已預提現总额
		ps.AddVarcharOutput("chvHelp", "")          // 预提规则

		// 設定觸發mock參數
		mockProc1 := db.NewProcedure("PrPs_UserAgentGetWithAmount_Get", ps)

		// 產生mock instance
		mockDbUtil := db.NewMockDatabaseEngine(mockCtrl)

		// 產生結果集
		mockretset := db.RecordSet{}

		// 產生 table0 結果值
		dataSetResult0 := db.RecordTable{}
		dataSetResult0.Cols = append(dataSetResult0.Cols, "Status")
		dataSetResult0.Rows = append(dataSetResult0.Rows, db.RecordRow{int64(1)})
		mockretset = append(mockretset, &dataSetResult0)

		// 產生資料庫的回傳值
		dataSetResult1 := db.RecordTable{}
		dataSetResult1.Rows = append(dataSetResult1.Rows, db.RecordRow{int64(1)})
		mockretset = append(mockretset, &dataSetResult1)

		// 抽換原有的資料庫
		mockDbUtil.EXPECT().ExecProc(gomock.Eq(mockProc1)).Return(mockretset, nil) // 興資料相符時抽換
		// mockDbUtil.EXPECT().ExecProc(gomock.Any()).Return(mockretset, nil)

		// 之前產生的mock instance
		DbUtil = mockDbUtil

		// 產生資料庫查詢封包
		process := new(UserInitBindProcess)
		msgData := new(netproto.UserAgentGetWithAmountMsg)
		msgData.UserID = proto.Int32(56022519)
		rpcInfo := process.getRPCInfo()
		rpcMessage := new(rpcmsg.RPCMessage)
		rpcMessage.RPCInfo = rpcInfo
		rpcMessage.MsgData, _ = proto.Marshal(msgData) // 在這裡把資料數值塞入封包

		// 產生新的協定訊息
		m := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserAgentGetWithAmountID), rpcMessage)

		// 檢查協定HOOK
		handler := new(UserAgentGetWithAmountHandler)
		if !handler.IsHook(m.BClassID, m.SClassID) {
			t.Error("HallMsgClassID_UserAgentGetWithAmountID 錯誤")
		}

		// 檢查協定回傳訊息
		ctx := new(module.MouduleCtx)
		clt := new(network.ConnInfo)
		ret := handler.OnMessage(ctx, clt, m)
		if ret[0].BClassID != int32(netproto.MessageBClassID_Hall) || ret[0].SClassID != int32(netproto.HallMsgClassID_UserAgentGetWithAmountRetID) {
			t.Error("HallMsgClassID_UserAgentGetWithAmountRetID 錯誤")
		}

		// 查看回傳訊息內容
		t.Log(ret[0].MsgData.(*netproto.UserAgentGetWithAmountRet).GetStatus())
	*/
}
