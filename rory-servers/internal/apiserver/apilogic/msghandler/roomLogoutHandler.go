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
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// RoomLogoutHandler 玩家退房
type RoomLogoutHandler struct{}

// IsHook 是否匹配
func (h *RoomLogoutHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_LogoutRoomID)
}

// OnMessage 收到訊息處理
func (h *RoomLogoutHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.LogoutGameRoomInfo{}
		return mm
	})

	if rpc == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.LogoutGameRoomInfo)

	example.invalid("RoomLogoutHandler-OnMessage:UserId[%d],UpdateAmount[%d],LoginMoney[%d],TradeID[%d]", mm.GetUserID(), mm.GetUpdateAmount(), mm.GetLoginMoney(), mm.GetTradeId())

	params := map[string]interface{}{
		"playerID":     mm.GetUserID(),
		"serverID":     mm.GetServerID(),
		"gameID":       mm.GetGameID(),
		"tradeID":      mm.GetTradeId(),
		"updateAmount": util.CentToMoney(mm.GetUpdateAmount()),
		"betAmount":    util.CentToMoney(mm.GetBetAmount()),
		"payoutAmount": util.CentToMoney(mm.GetPayAmount()),
		"taxAmount":    util.CentToMoney(mm.GetTaxAmount()),
		"flowAmount":   util.CentToMoney(mm.GetFlowAmount()),
		"winCount":     mm.GetWinCount(),
		"loseCount":    mm.GetLoseCount(),
		"drawCount":    mm.GetDrawCount(),
		"ipAddress":    mm.GetIPAddress(),
		"lastWinLose":  mm.GetLastWinLose(),
		"loginTime":    mm.GetLoginTime(),
		"logoutTime":   mm.GetLogoutTime(),
	}

	if !service.GamePlayService.PostRoomOut(rpc, params) {
		mlog.Error("游戏退房失败, playerID:%v, serverID:%v, gameID:%v, tradeID%v", mm.GetUserID(), mm.GetServerID(), mm.GetGameID(), mm.GetTradeId())
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
