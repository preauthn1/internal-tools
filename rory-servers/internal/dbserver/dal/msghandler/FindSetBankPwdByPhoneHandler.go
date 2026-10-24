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

// FindSetBankPwdByPhoneHandler 找回银行密码
type FindSetBankPwdByPhoneHandler struct{}

// IsHook ?
func (rh *FindSetBankPwdByPhoneHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneID)
}

// OnMessage ?
func (rh *FindSetBankPwdByPhoneHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.FindSetBankPwdByPhone{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.FindSetBankPwdByPhone)

	vcode := mm.GetVCode()
	passwd := mm.GetPassword()
	ip := rpcinfo.GetIPAddress()
	userid := rpcinfo.GetUserID()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddVarcharInput("chvCode", vcode, 32)
	ps.AddVarcharInput("chvNewPwd", passwd, 32)
	ps.AddVarcharInput("chvIPAddress", ip, 32)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_User_FindSetBankPwdByPhone", ps)
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	retmessage := ""
	if err != nil {
		ctx.Error("执行存储过程PrPs_User_FindSetBankPwdByPhone错误, %v", err)
		code = 0
		retmessage = "服务器错误"
	} else {
		if ret.GetReturnValue() == 1 {
			code = 1
		}
		retmessage = ret.GetOutputParamValue("chvErrMsg").(string)
	}

	retmsg := new(msg.Message)
	retmsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retmsg.SClassID = int32(netproto.HallMsgClassID_FindSetBankPwdByPhoneRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
