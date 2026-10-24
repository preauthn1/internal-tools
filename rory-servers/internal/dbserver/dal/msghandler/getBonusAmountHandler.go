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

// GetBonusAmountHandler 查詢我的推廣碼
type GetBonusAmountHandler struct{}

// IsHook ?
func (rh *GetBonusAmountHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_NetProfitAgentGetBonusAmountID)
}

// OnMessage ?
func (rh *GetBonusAmountHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NetProfitAgentGetBonusAmountMsg{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.NetProfitAgentGetBonusAmountMsg)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	proc := db.NewProcedure("PrPs_NetProfitAgent_GetBonusAmount", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)
	retcode := int32(dbret.GetReturnValue())

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_NetProfitAgentGetBonusAmountRetID), nil)
	retdata := new(netproto.NetProfitAgentGetBonusAmountRet)

	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程 PrPs_NetProfitAgent_GetBonusAmount 錯誤,%v", err)
		ctx.Error(strMsg)
	} else {
		if retcode == 1 {
			tb := dbret[0]
			if len(tb.Rows) > 0 {
				retdata.Poolamount = proto.Int32(int32(tb.GetValueByColName(0, "poolamount").(int64))) // 玩家ID // 小寫
			} else {
				retdata.Poolamount = proto.Int32(0)
			}
		}
	}

	retmsg.MsgData = retdata
	mlog.System("準備回傳結果結果 USERID: %v, VALUE: %v", mm.GetUserID(), retmsg.MsgData.(*netproto.NetProfitAgentGetBonusAmountRet).GetPoolamount())
	return []*msg.Message{retmsg}
}
