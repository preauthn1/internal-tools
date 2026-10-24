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

// FindSetPwdByPhoneHandler 找回賬號密碼
type FindSetPwdByPhoneHandler struct{}

// IsHook ?
func (rh *FindSetPwdByPhoneHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_FindSetPwdByPhoneID)
}

// OnMessage ?
func (rh *FindSetPwdByPhoneHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FindSetPwdByPhone{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.FindSetPwdByPhone)

	tel := mm.GetTel()
	vcode := mm.GetVCode()
	passwd := mm.GetPassword()
	ip := rpcinfo.GetIPAddress()

	ps := db.NewSQLParameters()
	ps.AddVarcharInput("chvTel", tel, 32)
	ps.AddVarcharInput("chvCode", vcode, 32)
	ps.AddVarcharInput("chvNewPwd", passwd, 32)
	ps.AddVarcharInput("chvIPAddress", ip, 32)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_User_FindSetPwdByPhone", ps)
	ret, err := dbutil.DBUtil.ExecProc(proc)
	//var ret db.RecordSet
	//var err error
	//if userid < 150000000 {	//方舟
	//	ret, err = DbUtil.ExecProc(proc)
	//} else {	//其他公司
	//	ret, err = CompanyDbUtil.ExecProc(proc)
	//}

	var code int32 = 0
	retmessage := ""
	if err != nil {
		ctx.Error("執行存儲過程PrPs_User_FindSetPwdByPhone錯誤, %v", err)
		code = 0
		retmessage = "服務器錯誤"
	} else {
		if ret.GetReturnValue() == 1 {
			code = 1
		}
		retmessage = ret.GetOutputParamValue("chvErrMsg").(string)
	}

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_FindSetPwdByPhoneRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
