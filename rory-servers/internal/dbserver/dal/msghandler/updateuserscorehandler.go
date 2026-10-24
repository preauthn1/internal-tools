package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// UpdateUserScoreHandler ?
type UpdateUserScoreHandler struct{}

// IsHook ?
func (rh *UpdateUserScoreHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_UpdateUserScoreID)
}

// OnMessage ?
func (rh *UpdateUserScoreHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UpdateUserScore{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UpdateUserScore)
	example.invalid("update score %v", mm)

	msgid := rpc.GetQueueID()
	userid := mm.GetUserID()
	score := mm.GetScore()
	serverid := mm.GetServerID()
	taxamount := mm.GetTaxAmount()
	// couponAmount := mm.GetCouponAmount()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intMsgID", msgid)
	ps.AddIntInput("intUserID", userid)
	ps.AddBigIntInput("lngUpdateAmount", score)
	ps.AddIntInput("intServerid", serverid)
	ps.AddBigIntInput("lngTaxAmount", taxamount)
	// ps.AddBigIntInput("lngCouponAmount", couponAmount)

	proc := db.NewProcedure("PrGs_UpdateUserMoney", ps)
	defer func() {
		if err := recover(); err != nil {
			errorStr := util.GetStackInfo()
			mlog.NewDefaultLogger().Error("发生错误%s, 執行存儲過程%s", errorStr, proc.GetSQLString())
			panic(err)
		}
	}()
	_, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("方舟 执行存储过程PrGs_UpdateUserMoney错误, %v", err)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
