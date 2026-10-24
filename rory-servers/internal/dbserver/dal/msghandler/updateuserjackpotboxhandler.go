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

// UpdateUserJackpotBoxHandler ?
type UpdateUserJackpotBoxHandler struct{}

// IsHook ?
func (rh *UpdateUserJackpotBoxHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.GameMessageClassID_GameUpdateUserJackpotBox)
}

// OnMessage ?
func (rh *UpdateUserJackpotBoxHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UserJackpotToBox{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UserJackpotToBox)
	ps := db.NewSQLParameters()
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddBigIntInput("lngUpdateAmount", mm.GetUpdateAmount())
	ps.AddVarcharOutput("chvErrMsg", "")
	proc := db.NewProcedure("PrGs_UpdateQZNN_UserAmount", ps)
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
		ctx.Error("执行存储过程PrGs_UpdateQZNN_UserAmount错误, %v", err)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
