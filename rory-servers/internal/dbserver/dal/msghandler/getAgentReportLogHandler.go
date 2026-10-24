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

// NetProfitAgentGetAgentReportLogHandler USER 新增代理(淨利潤)
type NetProfitAgentGetAgentReportLogHandler struct{}

// IsHook ?
func (rh *NetProfitAgentGetAgentReportLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogID)
}

// OnMessage ?
func (rh *NetProfitAgentGetAgentReportLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetAgentReportLogMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetAgentReportLogMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharInput("dtStattime", mm.GetDatetime(), 32)

	//開始準備SP指令
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetAgentReportLog", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentReportLogRetID), nil)
	retdata := new(netproto.NetProfitAgentGetAgentReportLogRet)

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_NetProfitAgent_GetAgentReportLog錯誤 ,%v", err)
		ctx.Error(strMsg)
	} else {

		if retcode == 1 {
			tb := dbret[1]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentGetAgentReportLogList)
					subRetData.Userid = proto.Int32(int32(tb.GetValueByColName(i, "userID").(int64)))
					subRetData.CreateTime = proto.String(string(tb.GetValueByColName(i, "CreateTime").(string)))
					// subRetData.CreateTime = proto.String(tb.GetValueByColName(i, "CreateTime").(time.Time).Format("2006-01-02 15:04:05"))
					subRetData.SumSubUsers = proto.Int32(int32(tb.GetValueByColName(i, "SumSubUsers").(int64)))
					subRetData.SumChargeAmount = proto.Int32(int32(tb.GetValueByColName(i, "SumChargeAmount").(int64)))
					subRetData.WithdrawAmount = proto.Int32(int32(tb.GetValueByColName(i, "WithdrawAmount").(int64)))
					subRetData.Counts = proto.Int32(int32(tb.GetValueByColName(i, "Counts").(int64)))
					retdata.GetWithdrawDetailList = append(retdata.GetWithdrawDetailList, subRetData)
				}
			}
			tb = dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Totallastdate = proto.Int32(int32(tb.GetValueByColName(0, "totallastdate").(int64)))
				retdata.Totalagentuser = proto.Int32(int32(tb.GetValueByColName(0, "totalagentuser").(int64)))
			} else {
				retdata.Totallastdate = proto.Int32(0)
				retdata.Totalagentuser = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
