package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// RequestLoginBonusHandler 領取登入獎勵
type RequestLoginBonusHandler struct{}

// IsHook ?
func (rh *RequestLoginBonusHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && (sclsid == int32(netproto.HallMsgClassID_RequestLoginBonusID) || sclsid == int32(netproto.HallMsgClassID_ApplyLoginBonusID))
}

// OnMessage ?
func (rh *RequestLoginBonusHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo := rpcmsg.GetRPCInfo(m)

	var retData *netproto.LoginBonusData
	userid := rpcinfo.GetUserID()
	ip := rpcinfo.GetIPAddress()
	if m.SClassID == int32(netproto.HallMsgClassID_RequestLoginBonusID) { //查詢
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddIntOutput("LoginTimes", 0)
		ps.AddIntOutput("GiftCash", 0)
		ps.AddVarcharOutput("chvRtnMsg", "")
		proc := db.NewProcedure("PrPs_Event_LoginBonusCheck", ps)
		ret, err := dbutil.DBUtil.ExecProc(proc)
		defer dbutil.HandleDBParseErr(proc, ret)
		var code int32 = 0
		var loginTimes int64 = 0
		var giftCash int64 = 0
		var retmessage = ""
		if err != nil {
			ctx.Error("执行存储过程PrPs_Event_LoginBonusCheck错误, %v", err)
			code = -99
			retmessage = "服务器错误"
		} else {
			code = int32(ret.GetReturnValue())
			if code == 1 { //查詢的 code 若等於1, 則回給Client 2
				code = 2
			}
			retmessage = ret.GetOutputParamValue("chvRtnMsg").(string)
			loginTimes = ret.GetOutputParamValue("LoginTimes").(int64)
			giftCash = ret.GetOutputParamValue("GiftCash").(int64)
		}
		retData = new(netproto.LoginBonusData)
		retData.RetCode = proto.Int32(code)
		retData.RetMsg = proto.String(retmessage)
		retData.LoginTimes = proto.Int32(int32(loginTimes))
		retData.GiftCash = proto.Int64(giftCash)

	} else if m.SClassID == int32(netproto.HallMsgClassID_ApplyLoginBonusID) { //領取
		ps := db.NewSQLParameters()
		ps.AddIntInput("intUserID", userid)
		ps.AddVarcharInput("chvIPAddress", ip, 15)
		ps.AddIntOutput("LoginTimes", 0)
		ps.AddIntOutput("lngCurrentAmount", 0)
		ps.AddIntOutput("GiftCash", 0)
		ps.AddVarcharOutput("chvRtnMsg", "")
		proc := db.NewProcedure("PrPs_Event_LoginBonus", ps)
		ret, err := dbutil.DBUtil.ExecProc(proc)
		defer dbutil.HandleDBParseErr(proc, ret)
		var code int32 = 0
		var loginTimes int64 = 0
		var currentAmount int64 = 0
		var giftCash int64 = 0
		var retmessage = ""
		if err != nil {
			ctx.Error("执行存储过程PrPs_Event_LoginBonus错误, %v", err)
			code = -99
			retmessage = "服务器错误"
		} else {
			code = int32(ret.GetReturnValue())
			retmessage = ret.GetOutputParamValue("chvRtnMsg").(string)
			loginTimes = ret.GetOutputParamValue("LoginTimes").(int64)
			currentAmount = ret.GetOutputParamValue("lngCurrentAmount").(int64)
			giftCash = ret.GetOutputParamValue("GiftCash").(int64)
		}
		retData = new(netproto.LoginBonusData)
		retData.RetCode = proto.Int32(code)
		retData.RetMsg = proto.String(retmessage)
		retData.LoginTimes = proto.Int32(int32(loginTimes))
		retData.CurrentAmount = proto.Int64(currentAmount)
		retData.GiftCash = proto.Int64(giftCash)
	}
	retData.LoginBonusConfigList = rh.getLoginBonusConfig(ctx, userid)
	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_RequestLoginBonusRetID)
	retmsg.MsgData = retData

	return []*msg.Message{retmsg}
}

func (rh *RequestLoginBonusHandler) getLoginBonusConfig(ctx *module.MouduleCtx, userID int32) (configList []*netproto.LoginBonusConfig) {
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userID)
	proc := db.NewProcedure("PrGs_Event_GetLoginEventSetting", ps)
	ret, err := dbutil.DBUtil.ExecProc(proc)
	defer dbutil.HandleDBParseErr(proc, ret)
	if err != nil {
		ctx.Error("执行存储过程 [%s] 错误, %v", proc.Name, err)
	}
	if ret.GetRetTableCount() > 0 {
		table := ret[0]
		if len(table.Rows) > 0 {
			table.ForEachRow(func(r db.RecordRow) {
				config := new(netproto.LoginBonusConfig)
				config.Day = proto.Int32(int32(r[0].(int64)))
				config.BonusAmount = proto.Int64(r[1].(int64))
				configList = append(configList, config)
			})
		}
	}
	return
}
