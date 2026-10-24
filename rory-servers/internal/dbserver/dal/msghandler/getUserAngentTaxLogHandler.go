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

// GetUserAngentTaxLogHandler 查询我的推廣碼
type GetUserAngentTaxLogHandler struct{}

// IsHook ?
func (rh *GetUserAngentTaxLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_GetUserAgentTaxLogID)
}

// OnMessage ?
func (rh *GetUserAngentTaxLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GetUserAgentTaxLogMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.GetUserAgentTaxLogMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharInput("dtStattime", mm.GetDateTime(), 32)
	// ps.AddIntInput("intArea_ID", mm.GetArea_ID())
	proc := db.NewProcedure("PrPs_UserAgent_GetAgentTaxLog", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GetUserAgentTaxLogRetID), nil)
	retdata := new(netproto.GetUserAgentTaxLogRetList)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_UserAgent_GetUserAgentTaxLog错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.TotalTax = proto.Int32(int32(tb.GetValueByColName(0, "totalTaxAmount").(int64)))     //稅收總額
				retdata.TotalTaxBonux = proto.Int32(int32(tb.GetValueByColName(0, "totalTaxBonus").(int64))) //稅收總額獎勵
				retdata.TotalPeople = proto.Int32(int32(tb.GetValueByColName(0, "temp_bpeople").(int64)))    ///總人數
			} else {
				retdata.TotalTax = proto.Int32(0)      //稅收總額
				retdata.TotalTaxBonux = proto.Int32(0) //稅收總額獎勵
				retdata.TotalPeople = proto.Int32(0)   ///總人數
			}
			tb = dbret[1]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.GetUserAgentTaxLogRet)
					subRetData.UserID = proto.Int32(int32(tb.GetValueByColName(i, "userID").(int64)))                 // 玩家ID // 小寫
					subRetData.Nickname = proto.String(string(tb.GetValueByColName(i, "nickname").(string)))          // 玩家昵稱 // 小寫
					subRetData.Btaxamount = proto.Int32(int32(tb.GetValueByColName(i, "B_taxamount").(int64)))        // 產生稅收 // 大寫
					subRetData.Cdtaxamount = proto.Int32(int32(tb.GetValueByColName(i, "CD_taxamount").(int64)))      // 一級+二級產生稅收 // 大寫
					subRetData.TotalTaxAmount = proto.Int32(int32(tb.GetValueByColName(i, "TotalTaxAmount").(int64))) // 稅收收益 // 大寫
					retdata.UserAgentTaxLog = append(retdata.UserAgentTaxLog, subRetData)
				}
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
