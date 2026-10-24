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

// GetWithdrawAmountHandler 提現紀錄(淨利潤)
type GetWithdrawAmountHandler struct{}

// IsHook ?
func (rh *GetWithdrawAmountHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawAmountID)
}

// OnMessage ?
func (rh *GetWithdrawAmountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetWithdrawAmountMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetWithdrawAmountMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetAgentWithdrawAmount", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())
	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawAmountRetID), nil)
	retdata := new(netproto.NetProfitAgentGetWithdrawAmountRet)

	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_GetAgentWithdrawAmount 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				accumulateNumber := int32(0)
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentGetWithdrawAmountList)
					subRetData.CreateTime = proto.String(string(tb.GetValueByColName(i, "CreateTime").(string)))
					subRetData.WithDrawAmount = proto.Int32(int32(tb.GetValueByColName(i, "WithDrawAmount").(int64)))
					subRetData.Status = proto.Int32(int32(tb.GetValueByColName(i, "Status").(int64)))
					subRetData.Memo = proto.String(string(tb.GetValueByColName(i, "Memo").(string)))
					accumulateNumber += int32(tb.GetValueByColName(i, "WithDrawAmount").(int64))
					retdata.GetWithdrawAmountRet = append(retdata.GetWithdrawAmountRet, subRetData)
				}
				retdata.Sum = proto.Int32(accumulateNumber)
			} else {
				retdata.Sum = proto.Int32(int32(0))
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
