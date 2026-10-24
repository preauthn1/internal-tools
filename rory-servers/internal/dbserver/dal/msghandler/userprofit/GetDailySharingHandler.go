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

// GetDailySharingHandler ?
type GetDailySharingHandler struct{}

// IsHook ?
func (rh *GetDailySharingHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserProfitGetDailySharingID)
}

// OnMessage ?
func (rh *GetDailySharingHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserProfitDailySharingRequest{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserProfitDailySharingRequest)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", rpcInfo.GetUserID())
	if mm.GetLimitDays() > 0 {
		ps.AddIntInput("intLimitDays", mm.GetLimitDays())
	}
	proc := db.NewProcedure("PrPs_UserProfit_GetDailySharing", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetDailySharingRetID), nil)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserProfit_GetDailySharing错误,%v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}
	retcode := int32(dbret.GetReturnValue())
	if retcode != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	info := new(netproto.UserProfitDailySharingInfo)
	tb := dbret[0]
	for i := 0; i < len(tb.Rows); i++ {
		log := new(netproto.UserProfitDailySharingLog)
		log.LogDate = proto.String(tb.GetValueByColName(i, "LogDate").(string))
		log.Level1Sharing = proto.Int64(tb.GetValueByColName(i, "Level1Sharing").(int64))
		log.Level2Sharing = proto.Int64(tb.GetValueByColName(i, "Level2Sharing").(int64))
		log.Level3Sharing = proto.Int64(tb.GetValueByColName(i, "Level3Sharing").(int64))

		info.UserProfitDailySharingLog = append(info.UserProfitDailySharingLog, log)
	}
	retmsg.MsgData = info

	return []*msg.Message{retmsg}
}
