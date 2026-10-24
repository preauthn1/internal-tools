package userprofit

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// WithdrawRequestHandler ?
type WithdrawRequestHandler struct{}

// IsHook ?
func (rh *WithdrawRequestHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserProfitWithdrawLogID)
}

// OnMessage ?
func (rh *WithdrawRequestHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserProfitWithdrawRequest{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserProfitWithdrawRequest)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", rpcInfo.GetUserID())
	ps.AddBigIntInput("intWithdrawAmount", mm.GetAmount())
	ps.AddVarcharOutput("chvMsg ", "")
	proc := db.NewProcedure("PrPs_UserProfit_WithdrawRequest", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitWithdrawLogRetID), nil)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserProfit_WithdrawRequest错误,%v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	chvMsg := dbret.GetOutputParamValue("chvMsg").(string)
	retcode := int32(dbret.GetReturnValue())
	retmsg.MsgData = dbutil.GetRetMsg(retcode, chvMsg)

	return []*msg.Message{retmsg}
}
