package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/service"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
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
func (h *UpdateUserScoreHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_UpdateUserScoreID)
}

// OnMessage ?
func (h *UpdateUserScoreHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UpdateUserScore{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UpdateUserScore)
	example.invalid("update score %v", mm)

	params := map[string]interface{}{
		"playerID":     mm.GetUserID(),
		"serverID":     mm.GetServerID(),
		"tradeID":      0,
		"updateAmount": util.CentToMoney(mm.GetScore()),
		"IPAddress":    rpc.GetIPAddress(),
	}

	if !service.GameResultService.PostSettlement(rpc, params) {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
