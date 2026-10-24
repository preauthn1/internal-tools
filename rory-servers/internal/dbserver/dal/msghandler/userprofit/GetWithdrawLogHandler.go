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

// GetWithdrawLogHandler ?
type GetWithdrawLogHandler struct{}

// IsHook ?
func (rh *GetWithdrawLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogID)
}

// OnMessage ?
func (rh *GetWithdrawLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserProfitDailySharingRequest{}
		return mm
	})

	if rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserProfitDailySharingRequest)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", rpcInfo.GetUserID())
	if mm.GetLimitDays() > 0 {
		ps.AddIntInput("intLimitDays", mm.GetLimitDays())
	}
	proc := db.NewProcedure("PrPs_UserProfit_GetWithdrawLog", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetWithdrawLogRetID), nil)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserProfit_GetWithdrawLog错误,%v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retcode := int32(dbret.GetReturnValue())
	if retcode != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	info := new(netproto.UserProfitWithdrawInfo)
	tb0 := dbret[0]
	info.UserProfit_SharingTotal = proto.Int64(tb0.GetValueByColName(0, "UserProfit_SharingTotal").(int64))
	info.UserProfit_WithdrawTotal = proto.Int64(tb0.GetValueByColName(0, "UserProfit_WithdrawTotal").(int64))
	info.UserProfit_Remaining = proto.Int64(tb0.GetValueByColName(0, "UserProfit_Remaining").(int64))

	tb := dbret[1]
	for i := 0; i < len(tb.Rows); i++ {
		log := new(netproto.UserProfitWithdrawLog)
		log.LogDate = proto.String(tb.GetValueByColName(i, "LogDate").(string))
		log.WithdrawAmount = proto.Int64(tb.GetValueByColName(i, "WithdrawAmount").(int64))
		log.Status = proto.String(tb.GetValueByColName(i, "Status").(string))

		info.UserProfitWithdrawLog = append(info.UserProfitWithdrawLog, log)
	}
	retmsg.MsgData = info

	return []*msg.Message{retmsg}
}
