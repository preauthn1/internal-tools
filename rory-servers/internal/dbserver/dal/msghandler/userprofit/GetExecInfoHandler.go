package userprofit

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

// GetExecInfoHandler ?
type GetExecInfoHandler struct{}

// IsHook ?
func (rh *GetExecInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserProfitGetExecInfoID)
}

// OnMessage ?
func (rh *GetExecInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo := rpcmsg.GetRPCInfo(m)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", rpcInfo.GetUserID())
	proc := db.NewProcedure("PrPs_UserProfit_GetExecInfo", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetExecInfoRetID), nil)

	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程PrPs_UserProfit_GetExecInfo錯誤,%v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retcode := int32(dbret.GetReturnValue())
	if retcode != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.UserProfitExecInfo)
	tb := dbret[0]
	if len(tb.Rows) == 0 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata.AgentTotalUser = proto.Int32(int32(tb.GetValueByColName(0, "AgentTotalUser").(int64)))
	retdata.UserProfit_SharingTotal = proto.Int64(int64(tb.GetValueByColName(0, "UserProfit_SharingTotal").(int64)))
	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
