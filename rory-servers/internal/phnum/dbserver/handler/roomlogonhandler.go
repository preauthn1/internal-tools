package handler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"

	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomLogonHandler 玩家登入
type RoomLogonHandler struct{}

// IsHook 是否攔截
func (h *RoomLogonHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_GameRoom) && sclsid == int32(netproto.GameRoomClassID_LoginRoomID)
}

// OnMessage 收到訊息
func (h *RoomLogonHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LoginGameRoomInfo{}
		return mm
	})

	mm := rmm.(*netproto.LoginGameRoomInfo)

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	example.invalid("user logon room.%v", mm)

	hdtype := mm.GetHDType()
	hdcode := mm.GetHDCode()
	ip := rpc.GetIPAddress()
	serverid := rpc.GetRouteServerID()
	userid := mm.GetUserID()
	cer := mm.GetCer()

	ps := db.NewSQLParameters()
	ps.AddIntInput("sintServerID", serverid)
	ps.AddIntInput("intUserID", userid)
	ps.AddIntInput("tnyHDType", hdtype)
	ps.AddVarcharInput("chvHDCode", hdcode, 64)
	ps.AddVarcharInput("chvIPAddress", ip, 15)
	ps.AddVarcharInput("chvCer", cer, 32)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrGs_UserLoginFish", ps)

	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomRetID), nil)

	if err != nil {
		mlog.Error("執行存儲過程PrGs_UserLogin出錯%v", err)
		logret := &netproto.LoginGameRoomRet{}
		logret.Code = proto.Int32(0)
		logret.Message = proto.String(fmt.Sprintf("服务器错误."))
		retmsg.MsgData = logret

	} else {
		var code int64 = 0
		if ret.GetReturnValue() == 1 {
			code = 1
		} else {
			code = ret.GetReturnValue()
		}

		logret := &netproto.LoginGameRoomRet{}
		logret.Code = proto.Int32(int32(code))
		logret.Message = proto.String(ret.GetOutputParamValueString("chvErrMsg"))
		logret.LoginRequestData = mm

		if ret.GetRetTableCount() > 1 {
			userdata := new(netproto.UserRoomLogonData)

			userdata.UserID = proto.Int32(ret[0].GetSingleValueInt32("UserID"))
			userdata.UserType = proto.Int32(ret[0].GetSingleValueInt32("UserType"))
			userdata.NickName = proto.String(ret[0].GetSingleValueString("NickName"))
			userdata.FaceID = proto.Int32(ret[0].GetSingleValueInt32("FaceID"))
			userdata.Sex = proto.Int32(ret[0].GetSingleValueInt32("Sex"))
			userdata.CashAmount = proto.Int64(ret[0].GetSingleValueInt64("CashAmount"))
			userdata.WinCount = proto.Int32(ret[0].GetSingleValueInt32("WinCount"))
			userdata.LoseCount = proto.Int32(ret[0].GetSingleValueInt32("LoseCount"))
			userdata.DrawCount = proto.Int32(ret[0].GetSingleValueInt32("DrawCount"))
			userdata.ServerAddr = proto.String(ret[0].GetSingleValueString("ServerAddr"))
			userdata.GameBuff = proto.String(ret[0].GetSingleValueString("GameBuff"))
			userdata.TradeId = proto.Int64(0)
			userdata.LastWinLose = proto.Int32(0)
			userdata.TimeSeq = proto.Int64(0)

			userdata.SumCharge = proto.Int64(0)
			userdata.HasBindPay = proto.Bool(false)
			userdata.SumConvert = proto.Int64(0)
			userdata.ExperienceCash = proto.Int64(0)
			userdata.LoginName = proto.String("")
			userdata.AlipayAccount = proto.String("")

			logret.UserData = userdata
			mlog.Debug("执行存储过程PrGs_UserLogin成功 %v", userdata.UserID)
		}

		retmsg.MsgData = logret
	}

	return []*msg.Message{retmsg}
}

//取得充值總金額和是否綁定
func getGetUserInfo(userid int32) (int64, bool, int64, int64, string, string, int32) {
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddVarcharOutput("chvErrMsg", "")
	proc := db.NewProcedure("PrGs_User_GetUserInfo", ps)

	ret, err := dbutil.DBUtil.ExecProc(proc)

	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().SetData(ret).Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()

	if err != nil {
		mlog.Error("执行存储过程PrGs_User_GetUserInfo出错%v", err)
		return int64(0), false, int64(0), int64(0), "", "", int32(0)
	}

	if ret.GetRetTableCount() > 1 {
		sumCharge := ret[0].GetSingleValueInt64("SumCharge")
		isBindZhifubao := ret[0].GetSingleValueInt32("IsBindZhifubao") == 1
		isBindUnionpayBank := ret[0].GetSingleValueInt32("IsBindUnionpayBank") == 1
		sumConvert := ret[0].GetSingleValueInt64("SumConvert")
		hasBindPay := (isBindZhifubao == true) || (isBindUnionpayBank == true)
		ExperienceCash := ret[0].GetSingleValueInt64("ExperienceCash")
		isBindGuest := ret[0].GetSingleValueInt32("IsBindZhifubao") == 1
		loginName := ""
		if isBindGuest {
			loginName = ret[0].GetSingleValueString("LoginName")
		}
		alipayAccount := ret[0].GetSingleValueString("AlipayAccount")
		i32AvatarFrameID := ret[0].GetSingleValueInt32("AvatarFrame")
		return sumCharge, hasBindPay, sumConvert, ExperienceCash, loginName, alipayAccount, i32AvatarFrameID
	}

	return int64(0), false, int64(0), int64(0), "", "", int32(0)
}
