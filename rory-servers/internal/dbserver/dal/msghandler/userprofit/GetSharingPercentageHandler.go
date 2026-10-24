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

// GetSharingPercentageHandler ?
type GetSharingPercentageHandler struct{}

// IsHook ?
func (rh *GetSharingPercentageHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) &&
		sclsid == int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageID)
}

// OnMessage ?
func (rh *GetSharingPercentageHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo := rpcmsg.GetRPCInfo(m)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", rpcInfo.GetUserID())
	proc := db.NewProcedure("PrPs_UserProfit_GetSharingPercentage ", ps)
	dbret, err := dbutil.DBUtil.ExecProc(proc)

	retmsg := msg.NewMessage(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_UserProfitGetSharingPercentageRetID), nil)

	if err != nil {
		strMsg := fmt.Sprintf("執行存儲過程PrPs_UserAgent_GetUserAgentList錯誤,%v", err)
		ctx.Error(strMsg)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retcode := int32(dbret.GetReturnValue())
	if retcode != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.UserProfitSharingPercentage)
	tb := dbret[0]
	if len(tb.Rows) == 0 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata.Level1 = proto.Int32(int32(tb.GetValueByColName(0, "Level1").(int64)))
	retdata.Level2 = proto.Int32(int32(tb.GetValueByColName(0, "Level2").(int64)))
	retdata.Level3 = proto.Int32(int32(tb.GetValueByColName(0, "Level3").(int64)))
	retmsg.MsgData = retdata

	return []*msg.Message{retmsg}
}
