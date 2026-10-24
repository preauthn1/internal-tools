package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomLogonHandler ?
type RoomLogonHandler struct{}

// IsHook ?
func (h *RoomLogonHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_GameRoom) && sclsid == int32(netproto.GameRoomClassID_LoginRoomID)
}

// OnMessage ?
func (h *RoomLogonHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LoginGameRoomInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.LoginGameRoomInfo)

	var params = map[string]interface{}{
		"serverID":     rpc.GetRouteServerID(),
		"isFishHunter": false,
		"ipAddress":    rpc.GetIPAddress(),
	}

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomRetID), nil)

	ok, gameRoomEnterBody, playerWalletBody := service.GamePlayService.PostRoomIn(rpc, params)
	if !ok || gameRoomEnterBody.Data.PlayerInfo.PlayerID == 0 {
		logret := new(netproto.LoginGameRoomRet)

		if gameRoomEnterBody == nil && playerWalletBody != nil {
			mlog.Error("呼叫 API [%s] 失敗", "/players/wallet")
			logret.Code = proto.Int32(0)
			logret.Message = proto.String(playerWalletBody.Message)
		} else if gameRoomEnterBody != nil {
			mlog.Error("呼叫 API [%s] 失敗", "/game/room/")
			logret.Code = proto.Int32(0)
			logret.Message = proto.String(gameRoomEnterBody.Message)
		} else {
			logret.Code = proto.Int32(0)
			logret.Message = proto.String(fmt.Sprintf("服务器错误."))
		}

		//@todo 400 Bad Request {"code":403203,"message":"金币不足，进入房间需要 1000.00 金币,您的金币数为 999.40 金币","data":{"ColumnName":null,"ReplaceMessage":"金币不足，进入房间需要 1000.00 金币,您的金币数为 999.40 金币"}}

		retmsg.MsgData = logret
	} else {
		logret := &netproto.LoginGameRoomRet{}
		logret.Code = proto.Int32(1)
		logret.Message = proto.String(fmt.Sprintf("%v", gameRoomEnterBody.Data.ReconnectInfo.ReConnectMsg))
		logret.LoginRequestData = mm

		// js := util.ToJSON(result.Data)
		// var player = util.ToMap([]byte(js))

		userdata := new(netproto.UserRoomLogonData)

		userdata.UserID = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.PlayerID)
		userdata.UserType = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.PlayerType)
		userdata.NickName = proto.String(gameRoomEnterBody.Data.PlayerInfo.DisplayName)
		userdata.FaceID = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.FaceID)
		userdata.Sex = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.Sex)
		userdata.CashAmount = proto.Int64(util.MoneyToCent(playerWalletBody.Data.Amount)) // userdata.CashAmount = proto.Int64(ret[0].GetSingleValueInt64("CashAmount"))
		userdata.WinCount = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.WinCount)
		userdata.LoseCount = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.LoseCount)
		userdata.DrawCount = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.DrawCount)
		userdata.ServerAddr = proto.String(rpc.GetIPAddress()) // userdata.ServerAddr = proto.String(ret[0].GetSingleValue("ServerAddr").(string))
		userdata.GameBuff = proto.String("")                   // userdata.GameBuff = proto.String(ret[0].GetSingleValue("GameBuff").(string))
		userdata.TradeId = proto.Int64(gameRoomEnterBody.Data.TradeID)
		userdata.LastWinLose = proto.Int32(gameRoomEnterBody.Data.PlayerInfo.LastWinLose)
		userdata.TimeSeq = proto.Int64(gameRoomEnterBody.Data.PlayerInfo.TimeSeq)
		// sumCharge, hasBindPay, sumConvert, experienceCash, loginName, alipayAccount := getGetUserInfo(userid)
		userdata.SumCharge = proto.Int64(0)                                             // userdata.SumCharge = proto.Int64(sumCharge)
		userdata.HasBindPay = proto.Bool(false)                                         // userdata.HasBindPay = proto.Bool(hasBindPay)
		userdata.SumConvert = proto.Int64(0)                                            // userdata.SumConvert = proto.Int64(sumConvert)
		userdata.ExperienceCash = proto.Int64(0)                                        // userdata.ExperienceCash = proto.Int64(experienceCash)
		userdata.LoginName = proto.String(gameRoomEnterBody.Data.PlayerInfo.PlayerName) // userdata.LoginName = proto.String(loginName)
		userdata.AlipayAccount = proto.String("")                                       // userdata.AlipayAccount = proto.String(alipayAccount)
		logret.UserData = userdata
		mlog.Debug("执行存储过程PrGs_UserLogin成功 %v", userdata.UserID)
		retmsg.MsgData = logret
	}

	return []*msg.Message{retmsg}
}

//
//func (h *RoomLogonHandler) requestAPI(method string, api string, params map[string]interface{}, cer string, apiCer string) (result1 *GameRoomEnterBody, result2 *PlayerWalletBody, err error) {
//
//	var groupRequest = new(httputil.WaitGroupRequest)
//	var accessTokenHeader = fmt.Sprintf("Access_Token:%s", cer)
//	var certificateHeader = fmt.Sprintf("X-API-Certificate:%s", apiCer)
//
//	var roominAPI = fmt.Sprintf("http://%s", path.Join(api, "/game/room/in"))
//	groupRequest.AddTask("0", httputil.POST, roominAPI, params, certificateHeader, accessTokenHeader)
//
//	var walletAPI = fmt.Sprintf("http://%s", path.Join(api, "/players/wallet"))
//	groupRequest.AddTask("1", httputil.GET, walletAPI, nil, certificateHeader, accessTokenHeader)
//
//	var respnose = groupRequest.Send()
//	groupRequest.Output() // 印資訊
//	// result = make(map[string]interface{})
//	// for _, v := range respnose {
//	// 	// fmt.Println(v.Uid, v.Status)
//	// 	// fmt.Println(v.Body)
//	// 	body, err := v.Body, v.Err
//	// 	if err != nil {
//	// 		fmt.Println(err.Error())
//	// 	}
//	// 	var data = util.ToMapMap([]byte(body)) // 把 api 回傳 json 轉成 map 型別
//	// 	for k, v2 := range data["data"] {
//	// 		result[k] = v2
//	// 	}
//	// }
//	body, err := respnose["0"].Body, respnose["0"].Err
//	if err != nil {
//		var message = fmt.Sprintf("%s %s 错误 %s, %s", httputil.POST, roominAPI, util.ToJSON(params), err.Error())
//		mlog.Error("%s", message)
//		return nil, nil, err
//	}
//
//	var model1 GameRoomEnterBody
//	err = json.Unmarshal([]byte(body), &model1)
//	result1 = &model1
//
//	body2, err := respnose["1"].Body, respnose["1"].Err
//	if err != nil {
//		var message = fmt.Sprintf("%s %s 错误 %s, %s", httputil.POST, roominAPI, util.ToJSON(params), err.Error())
//		mlog.Error("%s", message)
//		return nil, nil, err
//	}
//
//	var model2 PlayerWalletBody
//	err = json.Unmarshal([]byte(body2), &model2)
//	result2 = &model2
//
//	return result1, result2, nil
//}

// //取得充值總金額和是否綁定
// func getGetUserInfo(userid int32) (int64, bool, int64, int64, string, string) {
// 	ps := db.NewSqlParameters()
// 	ps.AddIntInput("intUserID", userid)
// 	ps.AddVarcharOutput("chvErrMsg", "")
// 	proc := db.NewProcedure("PrGs_User_GetUserInfo", ps)

// 	ret, err := dbutil.DbUtil.ExecProc(proc)

// 	defer func() {
// 		if err := recover(); err != nil {
// 			errorStr := util.GetStackInfo()
// 			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
// 			panic(err)
// 		}
// 	}()

// 	if err != nil {
// 		mlog.Error("执行存储过程PrGs_User_GetUserInfo出错%v", err)
// 		return int64(0), false, int64(0), int64(0), "", ""
// 	} else {
// 		if ret.GetRetTableCount() > 1 {
// 			sumCharge := ret[0].GetSingleValueInt64("SumCharge")
// 			isBindZhifubao := ret[0].GetSingleValueInt32("IsBindZhifubao") == 1
// 			isBindUnionpayBank := ret[0].GetSingleValueInt32("IsBindUnionpayBank") == 1
// 			sumConvert := ret[0].GetSingleValueInt64("SumConvert")
// 			hasBindPay := (isBindZhifubao == true) || (isBindUnionpayBank == true)
// 			ExperienceCash := ret[0].GetSingleValueInt64("ExperienceCash")
// 			isBindGuest := ret[0].GetSingleValueInt32("IsBindZhifubao") == 1
// 			loginName := ""
// 			if isBindGuest {
// 				loginName = ret[0].GetSingleValueString("LoginName")
// 			}
// 			alipayAccount := ret[0].GetSingleValueString("AlipayAccount")
// 			return sumCharge, hasBindPay, sumConvert, ExperienceCash, loginName, alipayAccount
// 		}
// 	}

// 	return int64(0), false, int64(0), int64(0), "", ""
// }
