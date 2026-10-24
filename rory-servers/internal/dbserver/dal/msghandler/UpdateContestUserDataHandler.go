package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/db"
	"example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// UpdateContestUserDataHandler 更新比賽活動用戶資料
type UpdateContestUserDataHandler struct{}

// IsHook ?
func (rh *UpdateContestUserDataHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.DBServerClassID_UpdateContestUserDataID)
}

// OnMessage ?
func (rh *UpdateContestUserDataHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.UpdateContestUser{}
		return mm
	})

	if rmm == nil || rpcInfo == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.UpdateContestUser)

	ps := db.NewSQLParameters()
	ps.AddIntInput("intGameID", mm.GetGameID())
	ps.AddIntInput("intUserID", mm.GetUserID())
	ps.AddIntInput("intContestID", mm.GetContestID())
	ps.AddIntInput("intContestNo", mm.GetContestNo())
	ps.AddIntInput("intFlagID", mm.GetFlagID())
	ps.AddIntInput("intTableNo", mm.GetTableNo())
	ps.AddIntInput("intPlaytime", mm.GetPlayTime())
	ps.AddIntInput("lngRemainScore", mm.GetRemainScore())
	ps.AddIntInput("lngScore", mm.GetScore())
	proc := db.NewProcedure("PrGs_Event_UpdateContestUserData", ps)
	_, err := dbutil.DBUtil.ExecProc(proc)

	var retcode int32 = 1
	if err != nil {
		ctx.Error("执行存储过程PrGs_Event_UpdateContestUserData错误, %v", err)
		retcode = 0
	}

	return []*msg.Message{dbutil.GetDBRetMsg(retcode)}
}
