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

// BindBankCardHandler ?
type BindBankCardHandler struct{}

// IsHook ?
func (rh *BindBankCardHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_BindBankCardID)
}

// OnMessage ?
func (rh *BindBankCardHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.BankCardInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.BankCardInfo)

	bankcardnum := mm.GetBankCardNum()
	realname := mm.GetRealName()
	bankcode := mm.GetBankCode()
	vcode := mm.GetVCode()
	codetype := mm.GetCodeType()

	userid := rpcinfo.GetUserID()
	cer := rpcinfo.GetCer()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddVarcharInput("chvCer", cer, 32)
	ps.AddVarcharInput("chvCode", vcode, 6)
	ps.AddIntInput("tintCodeType", codetype)
	ps.AddVarcharInput("chvRealName", realname, 64)
	ps.AddVarcharInput("chvUnionpayBank", bankcode, 5)
	ps.AddVarcharInput("chvUnionpayAccount", bankcardnum, 64)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_BindUnionPayAccount", ps)
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	retmessage := ""
	if err != nil {
		ctx.Error("PrPs_BindUnionPayAccount, %v", err)
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
	retmsg.SClassID = int32(netproto.HallMsgClassID_BindBankCardRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
