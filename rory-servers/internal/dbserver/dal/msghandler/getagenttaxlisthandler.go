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

// GetAgentTaxListHandler 查询全民推廣 收益規則
type GetAgentTaxListHandler struct{}

// IsHook ?
func (rh *GetAgentTaxListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetAgentTaxListID)
}

// OnMessage ?
func (rh *GetAgentTaxListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcinfo := rpcmsg.GetRPCInfo(m)

	if rpcinfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	userid := rpcinfo.GetUserID()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", userid)

	proc := db.NewProcedure("PrPs_UserAgent_GetAgentTaxList", ps)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetAgentTaxListRetID), nil)
	retdata := new(netproto.GetAgentTaxListRet)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程PrPs_UserAgent_GetAgentTaxList錯誤, %v", err)
		ctx.Error(strMsg)
		return []*msg.Message{retmsg}
	}

	if dbret.GetRetTableCount() > 0 {
		tb := dbret[0]
		if len(tb.Rows) > 0 {
			for i := 0; i < len(tb.Rows); i++ {
				pdata := new(netproto.AgentTaxList)
				pdata.Level = proto.Int32(int32(tb.GetValueByColName(i, "Level").(int64)))
				pdata.Stage1 = proto.Int32(int32(tb.GetValueByColName(i, "Stage1").(int64)))
				pdata.Stage2 = proto.Int32(int32(tb.GetValueByColName(i, "Stage2").(int64)))
				pdata.Stage3 = proto.Int32(int32(tb.GetValueByColName(i, "Stage3").(int64)))
				retdata.TaxList = append(retdata.TaxList, pdata)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
