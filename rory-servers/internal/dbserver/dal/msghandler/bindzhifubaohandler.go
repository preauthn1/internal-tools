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

// BindZhifubaoHandler ?
type BindZhifubaoHandler struct{}

// IsHook ?
func (rh *BindZhifubaoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_BindZhifubaoID)
}

// OnMessage ?
func (rh *BindZhifubaoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ZhifubaoInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.ZhifubaoInfo)
	zhifubao := mm.GetZhifubao()
	realname := mm.GetRealName()
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
	ps.AddVarcharInput("chvAlipayAccount", zhifubao, 64)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_BindAlipayAccount", ps)
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	retmessage := ""
	if err != nil {
		ctx.Error("执行存储过程PrPs_BindAlipayAccount错误, %v", err)
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
	retmsg.SClassID = int32(netproto.HallMsgClassID_BindZhifubaoRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
