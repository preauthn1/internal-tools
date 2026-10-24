package msghandler

import (
	"fmt"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// GetUserAngentExecInfoHandler 查询推廣員資料
type GetUserAngentExecInfoHandler struct{}

// IsHook ?
func (rh *GetUserAngentExecInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetUserAgentExecID)
}

// OnMessage ?
func (rh *GetUserAngentExecInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetUserAgentExecMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetUserAgentExecMsg)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_UserAgent_GetAgentExecInfo", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentExecRetID), nil)
	retdata := new(netproto.GetUserAgentExecRet)

	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_UserAgent_GetAgentExecInfo 錯誤,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.AgentAmount = proto.Int32(int32(tb.GetValueByColName(0, "AgentAmount").(int64)))
				retdata.AgentUserAmount = proto.Int32(int32(tb.GetValueByColName(0, "AgentUserAmount").(int64)))
				retdata.AgentLevel = proto.Int32(int32(tb.GetValueByColName(0, "AgentLevel").(int64)))
				retdata.AgentTaxAmount = proto.Int32(int32(tb.GetValueByColName(0, "AgentTaxAmount").(int64)))
				retdata.AgentTaxBonus = proto.Int32(int32(tb.GetValueByColName(0, "AgentTaxBonus").(int64)))
				retdata.AgentTotalUser = proto.Int32(int32(tb.GetValueByColName(0, "AgentTotalUser").(int64)))
			} else {
				retdata.AgentAmount = proto.Int32(0)
				retdata.AgentUserAmount = proto.Int32(0)
				retdata.AgentLevel = proto.Int32(1)
				retdata.AgentTaxAmount = proto.Int32(0)
				retdata.AgentTaxBonus = proto.Int32(0)
				retdata.AgentTotalUser = proto.Int32(0)
			}
		} else {
			// retcode == 0 時
			mlog.System("SQL 發生錯誤 %v", proc.GetSQLString())
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
