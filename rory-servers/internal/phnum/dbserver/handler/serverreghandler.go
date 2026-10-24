package handler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/netproto/phnum"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ServerRegHandler 註冊遊戲
type ServerRegHandler struct{}

// IsHook 是否攔截
func (rh *ServerRegHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_ServerMgm) && sclsid == int32(netproto.ServerMgmClassID_GameServerRegister)
}

// OnMessage 收到訊息
func (rh *ServerRegHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &phnum.PhnumServerRegInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*phnum.PhnumServerRegInfo)

	serverid := mm.GetServerID()
	servertype := mm.GetServerType()

	ps := db.NewSQLParameters()
	ps.AddIntInput("intServerID", serverid)
	ps.AddIntInput("intServerType", servertype)
	ps.AddIntInput("intGroupID", mm.GetLianyunID())

	proc := db.NewProcedure("PrGs_ReRegister", ps)
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
		ctx.Error("执行存储过程PrGs_ServerReg错误, %v", err)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
