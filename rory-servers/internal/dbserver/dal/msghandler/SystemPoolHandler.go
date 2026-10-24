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

// SystemPoolHandler ?
type SystemPoolHandler struct{}

// IsHook ?
func (h *SystemPoolHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose) ||
		sclsid == int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)
}

// OnMessage ?
func (h *SystemPoolHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.GameMessageClassID(m.SClassID) {
	case netproto.GameMessageClassID_SuperToolGetSysPoolWinLose:
		return h.getSysPoolWinLose(ctx, clt, m)
	case netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose:
		return h.updateSysPoolWinLose(ctx, clt, m)
	}

	return []*msg.Message{}
}

// 要求 水線控制遊戲輸贏結果
func (h *SystemPoolHandler) getSysPoolWinLose(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperTool_Msg_GetSysPoolWinLose{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperTool_Msg_GetSysPoolWinLose)

	ps := db.NewSQLParameters()

	ps.AddIntInput("intGameID", mm.GetGameID()) // 遊戲編號
	ps.AddIntInput("intFlagID", mm.GetFlagID()) // 房間編號

	proc := db.NewProcedure("PrGs_SuperTool_GetSystemPoolWinLost", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetReturnValue() != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	if err != nil {
		ctx.Error("执行存储过程 PrGs_SuperTool_GetSystemPoolWinLost 错误, %v, %s", err, proc.GetSQLString())
	} else {
		if dbret.GetRetTableCount() >= 1 {
			retdata.GameID = proto.Int32(dbret[0].GetSingleValueInt32("GameID"))          // 遊戲編號
			retdata.FlagID = proto.Int32(dbret[0].GetSingleValueInt32("FlagID"))          // 房間編號
			retdata.PoolAmount = proto.Int64(dbret[0].GetSingleValueInt64("PoolAmount"))  // 水線
			retdata.WinLine = proto.Int64(dbret[0].GetSingleValueInt64("WinLine"))        // 贏錢線
			retdata.WinLoseData = proto.String(dbret[0].GetSingleValueString("GameData")) // 遊戲輸贏結果
		}
	}

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet)
	retMsg.MsgData = retdata

	return []*msg.Message{retMsg}
}

// 更新 水線控制遊戲輸贏結果
func (h *SystemPoolHandler) updateSysPoolWinLose(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperTool_Msg_UpdateSysPoolWinLose{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperTool_Msg_UpdateSysPoolWinLose)

	ps := db.NewSQLParameters()

	ps.AddIntInput("intGameID", mm.GetGameID())                  // 遊戲編號
	ps.AddIntInput("intFlagID", mm.GetFlagID())                  // 房間編號
	ps.AddBigIntInput("lngPoolAmount", mm.GetPoolAmount())       // 水線
	ps.AddBigIntInput("lngWinLine", mm.GetWinLine())             // 贏錢線
	ps.AddVarcharInput("chvGameData", mm.GetWinLoseData(), 4000) // 遊戲輸贏結果

	proc := db.NewProcedure("PrGs_SuperTool_UpdSystemPoolWinLost", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetReturnValue() != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetReturnValue() != 1 {
		ctx.Error("执行存储过程 PrGs_SuperTool_UpdSystemPoolWinLost 错误")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
