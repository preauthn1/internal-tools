package msghandler

import (
	"math"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ConvertMoneyHandler 聯繫客服
type ConvertMoneyHandler struct{}

// IsHook ?
func (rh *ConvertMoneyHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_ConvertMoneyID)
}

// OnMessage ?
func (rh *ConvertMoneyHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ConvertMoney{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.ConvertMoney)
	amount := int64(math.Abs(float64(mm.GetAmount())))
	convertType := mm.GetConvertType()

	userid := rpcinfo.GetUserID()
	cer := rpcinfo.GetCer()
	ip := rpcinfo.GetIPAddress()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)
	ps.AddVarcharInput("chvCer", cer, 32)
	ps.AddBigIntInput("lngAmount", amount)
	ps.AddTinyIntInput("tnyConvertType", byte(convertType))
	ps.AddVarcharInput("chvIP", ip, 15)
	ps.AddVarcharOutput("chvErrMsg", "")

	proc := db.NewProcedure("PrPs_AddConvertLog", ps)
	//ret, err := DbUtil.ExecProc(proc)
	var ret db.RecordSet
	var err error
	ret, err = dbutil.DBUtil.ExecProc(proc)

	var code int32 = 0
	retmessage := ""
	if err != nil {
		ctx.Error("执行存储过程PrPs_AddConvertLog错误, %v", err)
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
	retmsg.SClassID = int32(netproto.HallMsgClassID_ConvertMoneyRetID)
	retmsg.MsgData = dbutil.GetRetMsg(code, retmessage)

	return []*msg.Message{retmsg}
}
