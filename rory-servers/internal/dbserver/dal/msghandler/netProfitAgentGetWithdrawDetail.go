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

// NetProfitAgentGetWithdrawDetailHandler USER 新增代理(淨利潤)
type NetProfitAgentGetWithdrawDetailHandler struct{}

// IsHook ?
func (rh *NetProfitAgentGetWithdrawDetailHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailID)
}

// OnMessage ?
func (rh *NetProfitAgentGetWithdrawDetailHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetWithdrawDetailMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetWithdrawDetailMsg)
	ps := db.NewSQLParameters()

	//在這裡設定輸入參數
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharInput("dtStattime", mm.GetDatetime(), 32)

	//開始準備SP指令
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetWithdrawDetail", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetWithdrawDetailRetID), nil)
	retdata := new(netproto.NetProfitAgentGetWithdrawDetailRet)

	if err != nil || retcode != 1 {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_GetWithdrawDetail 错误,%v", err)
		ctx.Error(strMsg)
	} else {

		if retcode == 1 {
			tb := dbret[1]
			accumulateNumber := int32(0)
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentGetWithdrawDetailList)
					subRetData.LogDate = proto.String(string(tb.GetValueByColName(i, "LogDate").(string)))
					// subRetData.LogDate = proto.String(tb.GetValueByColName(i, "LogDate").(time.Time).Format("2006-01-02 15:04:05"))
					subRetData.UserID = proto.Int32(int32(tb.GetValueByColName(i, "UserID").(int64)))
					subRetData.ChargeAmount = proto.Int32(int32(tb.GetValueByColName(i, "ChargeAmount").(int64)))
					subRetData.ConvertAmount = proto.Int32(int32(tb.GetValueByColName(i, "ConvertAmount").(int64)))
					subRetData.StorageAmount = proto.Int32(int32(tb.GetValueByColName(i, "StorageAmount").(int64)))
					subRetData.BonusAmount = proto.Int32(int32(tb.GetValueByColName(i, "BonusAmount").(int64)))
					retdata.GetWithdrawDetailList = append(retdata.GetWithdrawDetailList, subRetData)
					accumulateNumber += int32(tb.GetValueByColName(i, "BonusAmount").(int64))
				}
				retdata.Sum = proto.Int32(accumulateNumber) // 如果LIST有資料，就立刻更新SUM值
			} else {
				retdata.Sum = proto.Int32(0) // 如果LIST沒有資料，SUM值就為0
			}
			tb = dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Chargeamountlast = proto.Int32(int32(tb.GetValueByColName(0, "chargeamountlast").(int64)))
				// retdata.Sum = proto.Int32(accumulateNumber) // SUM值不應在這裡除理
			} else {
				retdata.Chargeamountlast = proto.Int32(int32(0))
				retdata.Sum = proto.Int32(int32(0))
			}
		}
	}

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
