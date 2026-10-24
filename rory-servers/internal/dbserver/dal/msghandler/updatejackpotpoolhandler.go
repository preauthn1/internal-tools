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

// UpdateJackpotPoolHandler ?
type UpdateJackpotPoolHandler struct{}

// IsHook ?
func (rh *UpdateJackpotPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.DFDC_GameMessageClassID_DFDCUpdateJackpotPoolID)
}

// OnMessage ?
func (rh *UpdateJackpotPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.DFDC_GetJackpotRet{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}
	mm := rmm.(*netproto.DFDC_GetJackpotRet)
	ps := db.NewSQLParameters()
	ps.AddIntInput("sintServerID", mm.GetServerID())
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddIntInput("intTableNo", mm.GetTableNo())
	ps.AddIntInput("intJackPotID", mm.GetJackpotID())
	ps.AddBigIntInput("intSpins", mm.GetSpinStack())
	ps.AddBigIntInput("lngPoolAmount", mm.GetPoolAmount())
	ps.AddBigIntInput("lngJPAmount", mm.GetJPAmountOut())
	mlog.System("执行存储过程PrGs_DFDC_UpdJackpotPool")
	proc := db.NewProcedure("PrGs_DFDC_UpdJackpotPool", ps)
	_, err := dbutil.WildCatDBUtil.ExecProc(proc)
	var retcode int32 = 1
	if err != nil {
		mlog.System("执行存储过程PrGs_DFDC_UpdJackpotPool错误")
		strMsg := fmt.Sprintf("执行存储过程PrGs_DFDC_UpdJackpotPool错误, %v", err)
		ctx.Error(strMsg)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
