package msghandler

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetAgentTaxBonusLogHandler 查询全民推廣 每日收益
type GetAgentTaxBonusLogHandler struct{}

// IsHook ?
func (rh *GetAgentTaxBonusLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetAgentTaxBonusLogID)
}

// OnMessage ?
func (rh *GetAgentTaxBonusLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetAgentTaxBonusLogMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetAgentTaxBonusLogMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_UserAgent_GetAgentTaxBonusLog", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxBonusLogRetID), nil)
	retdata := new(netproto.GetAgentTaxBonusLogRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_UserAgent_GetAgentTaxBonusLog 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					pdata := new(netproto.AgentTaxBonusLogList)
					pdata.LogDate = proto.String(tb.GetValueByColName(i, "LogDate").(time.Time).Format("2006-01-02 15:04:05"))
					pdata.BTaxBonusAmount = proto.Int32(int32(tb.GetValueByColName(i, "B_TaxBonusAmount").(int64)))
					pdata.CDTaxBonusAmount = proto.Int32(int32(tb.GetValueByColName(i, "CD_TaxBonusAmount").(int64)))
					pdata.TotalTaxBonusAmount = proto.Int32(int32(tb.GetValueByColName(i, "TotalTaxBonusAmount").(int64)))
					retdata.TaxBonusLogList = append(retdata.TaxBonusLogList, pdata)
				}
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
