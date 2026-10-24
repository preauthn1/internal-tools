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

// GetAgentWithdrawAmountHandler 代理提現紀錄(淨利潤)
type GetAgentWithdrawAmountHandler struct{}

// IsHook ?
func (rh *GetAgentWithdrawAmountHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetAgentWithdrawAmountID)
}

// OnMessage ?
func (rh *GetAgentWithdrawAmountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetAgentWithdrawAmountMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetAgentWithdrawAmountMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetAgentWithdrawAmount", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetAgentWithdrawAmountRetID), nil)
	retdata := new(netproto.NetProfitAgentGetAgentWithdrawAmountRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程PrPs_NetProfitAgent_GetAgentWithdrawAmount错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			retdata.Sum = proto.Int32(0)
			accumulateNumber := int32(0)
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentGetAgentWithdrawAmountList)
					subRetData.UserID = proto.Int32(int32(tb.GetValueByColName(i, "UserID").(int64)))
					subRetData.CreateTime = proto.String(string(tb.GetValueByColName(i, "CreateTime").(string)))
					subRetData.WithDrawAmount = proto.Int32(int32(tb.GetValueByColName(i, "WithDrawAmount").(int64)))
					subRetData.Status = proto.String(string(tb.GetValueByColName(i, "Status").(string)))
					retdata.GetAgentWithdrawAmountList = append(retdata.GetAgentWithdrawAmountList, subRetData)
					accumulateNumber += int32(tb.GetValueByColName(i, "WithDrawAmount").(int64))

				}
				retdata.Sum = proto.Int32(accumulateNumber)
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
