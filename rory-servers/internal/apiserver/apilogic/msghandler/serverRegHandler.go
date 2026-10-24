package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ServerRegHandler ?
type ServerRegHandler struct{}

// IsHook ?
func (h *ServerRegHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_ServerMgm) && sclsid == int32(netproto.ServerMgmClassID_GameServerRegister)
}

// OnMessage ?
func (h *ServerRegHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ServerRegInfo{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.ServerRegInfo)

	params := map[string]interface{}{
		"serverID":   mm.GetServerID(),
		"serverType": mm.GetServerType(),
	}

	if !service.ServerSettingService.PostReRegister(rpc, params) {
		mlog.Error("服务器注册失败")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}
	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
