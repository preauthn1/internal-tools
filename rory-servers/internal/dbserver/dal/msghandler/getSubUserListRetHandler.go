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

// GetSubUserListRetHandler ?
type GetSubUserListRetHandler struct{}

// IsHook ?
func (rh *GetSubUserListRetHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetSubUserListID)
}

// OnMessage ?
func (rh *GetSubUserListRetHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetSubUserListMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetSubUserListMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddVarcharInput("dtStatDate", mm.GetDatetime(), 32)
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetSubUserList", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetSubUserListRetID), nil)
	retdata := new(netproto.NetProfitAgentGetSubUserListRet)

	accumulateSumCharge := int32(0)
	accumulateSumWithDraw := int32(0)
	accumulateSumBonus := int32(0)
	if err != nil {
		strMsg := fmt.Sprintf("执行存储过程 PrPs_NetProfitAgent_GetSubUserList 错误,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[1]
			if len(tb.Rows) > 0 {
				for i := 0; i < len(tb.Rows); i++ {
					subRetData := new(netproto.NetProfitAgentGetSubUserList)
					subRetData.UserID = proto.Int32(int32(tb.GetValueByColName(i, "userID").(int64)))
					subRetData.RegDate = proto.String(string(tb.GetValueByColName(i, "RegDate").(string)))
					subRetData.SumChargeAmount = proto.Int32(int32(tb.GetValueByColName(i, "SumChargeAmount").(int64)))
					subRetData.SumWithDrawAmount = proto.Int32(int32(tb.GetValueByColName(i, "SumWithDrawAmount").(int64)))
					subRetData.BonusAmount = proto.Int32(int32(tb.GetValueByColName(i, "BonusAmount").(int64)))
					retdata.GetSubUserList = append(retdata.GetSubUserList, subRetData)
					accumulateSumCharge += int32(tb.GetValueByColName(i, "SumChargeAmount").(int64))
					accumulateSumWithDraw += int32(tb.GetValueByColName(i, "SumWithDrawAmount").(int64))
					accumulateSumBonus += int32(tb.GetValueByColName(i, "BonusAmount").(int64))
				}
			}
			tb = dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Totaluserlast = proto.Int32(int32(tb.GetValueByColName(0, "totaluserlast").(int64)))
			} else { //如果第一張表沒有資料，就填假資料
				retdata.Totaluserlast = proto.Int32(0) //如果第一張表沒有資料，就填假資料
				retdata.SumCharge = proto.Int32(0)     //如果第一張表沒有資料，就填假資料
				retdata.SumWithDraw = proto.Int32(0)   //如果第一張表沒有資料，就填假資料
				retdata.SumBonus = proto.Int32(0)      //如果第一張表沒有資料，就填假資料
			}
		}
	}
	retdata.SumCharge = proto.Int32(accumulateSumCharge)
	retdata.SumWithDraw = proto.Int32(accumulateSumWithDraw)
	retdata.SumBonus = proto.Int32(accumulateSumBonus)

	retmsg.MsgData = retdata
	return []*msg.Message{retmsg}
}
