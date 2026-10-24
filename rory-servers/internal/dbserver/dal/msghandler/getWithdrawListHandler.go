package msghandler

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

// GetWithdrawListIHandler ?
type GetWithdrawListIHandler struct{}

// IsHook ?
func (rh *GetWithdrawListIHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetWithdrawListID)
}

// OnMessage ?
func (rh *GetWithdrawListIHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetWithdrawMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetWithdrawMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_UserAgent_GetWithdrawList", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetWithdrawListRetID), nil)
	retdata := new(netproto.GetWithdrawRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程HallMsgClassID_GetWithdrawListID 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.TotalTaxAmount = proto.Int32(int32(tb.GetValueByColName(0, "TotalTaxAmount").(int64)))
				retdata.ConvertAmount = proto.Int32(int32(tb.GetValueByColName(0, "ConvertAmount").(int64)))
				retdata.PerformedAmount = proto.Int32(int32(tb.GetValueByColName(0, "PerformedAmount").(int64)))
				retdata.RemainAmount = proto.Int32(int32(tb.GetValueByColName(0, "RemainAmount").(int64)))
			} else { // 如果第一張表不存在，就回傳假資料
				retdata.TotalTaxAmount = proto.Int32(0)
				retdata.ConvertAmount = proto.Int32(0)
				retdata.PerformedAmount = proto.Int32(0)
				retdata.RemainAmount = proto.Int32(0)
			}
			tb = dbret[1] // 這裡是陣列
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.GetWithdrawList)
					subRetData.LogDate = proto.String(string(tb.GetValueByColName(i, "LogDate").(string)))
					subRetData.WithdrawAmount = proto.Int32(int32(tb.GetValueByColName(i, "WithdrawAmount").(int64)))
					subRetData.Status = proto.String(string(tb.GetValueByColName(i, "Status").(string)))
					retdata.WithdrawList = append(retdata.WithdrawList, subRetData)
				}
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
