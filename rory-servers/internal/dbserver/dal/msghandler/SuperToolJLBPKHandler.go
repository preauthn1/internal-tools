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
)

// SuperToolJLBPKHandler ?
type SuperToolJLBPKHandler struct{}

// IsHook ?
func (h *SuperToolJLBPKHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserList) ||
		sclsid == int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_UpdateVIPUser)
}

// OnMessage ?
func (h *SuperToolJLBPKHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.JLBPK_GameMessageClassID(m.SClassID) {
	case netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserList:
		return h.getVIPUserList(ctx, clt, m)
	case netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_UpdateVIPUser:
		return h.updateVIPUser(ctx, clt, m)
	}

	return []*msg.Message{}
}

// 要求 監控玩家名單
func (h *SuperToolJLBPKHandler) getVIPUserList(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	ps := db.NewSQLParameters()
	proc := db.NewProcedure("PrGs_SuperTool_JLBPK_GetVIPUser", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetReturnValue() != 1 {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retdata := new(netproto.JLBPK_SuperTool_Msg_GetVIPUserListRet)
	if err != nil {
		ctx.Error("执行存储过程 PrGs_SuperTool_JLBPK_GetVIPUser 错误, %v, %s", err, proc.GetSQLString())
	} else {
		if dbret.GetRetTableCount() >= 1 {
			for i := 0; i < len(dbret[0].Rows); i++ {
				pData := new(netproto.JLBPK_SuperTool_Msg_VIPUserData)
				pData.UserID = proto.Int32(dbret[0].GetValueByColNameInt32(i, "UserID"))                     // 玩家 ID
				pData.Finished = proto.Int32(dbret[0].GetValueByColNameInt32(i, "Finished"))                 // 是否已完成監控
				pData.Weight = proto.Int32(dbret[0].GetValueByColNameInt32(i, "Weight"))                     // 玩家拿牌權重，［可後台設定］
				pData.TargetMoney = proto.Int64(dbret[0].GetValueByColNameInt64(i, "TargetPoint"))           // 玩家目標金額，［可後台設定］
				pData.AccumulatedMoney = proto.Int64(dbret[0].GetValueByColNameInt64(i, "AccumulatedPoint")) // 目前累積金額
				pData.Status = proto.Int32(dbret[0].GetValueByColNameInt32(i, "Status"))

				retdata.UserList = append(retdata.UserList, pData)
			}
		}
	}

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.JLBPK_GameMessageClassID_JLBPK_SuperTool_GetVIPUserListRet)
	retMsg.MsgData = retdata

	return []*msg.Message{retMsg}
}

// 更新 監控玩家資料
func (h *SuperToolJLBPKHandler) updateVIPUser(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.JLBPK_SuperTool_Msg_VIPUserData{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.JLBPK_SuperTool_Msg_VIPUserData)

	ps := db.NewSQLParameters()

	ps.AddIntInput("intUserID", mm.GetUserID())                        // 玩家 ID
	ps.AddIntInput("intFinished", mm.GetFinished())                    // 是否已完成監控
	ps.AddIntInput("intWeight", mm.GetWeight())                        // 玩家拿牌權重，［可後台設定］
	ps.AddBigIntInput("intTargetPoint", mm.GetTargetMoney())           // 玩家目標金額，［可後台設定］
	ps.AddBigIntInput("intAccumulatedPoint", mm.GetAccumulatedMoney()) // 目前累積金額

	proc := db.NewProcedure("PrGs_SuperTool_JLBPK_UpdVIPUser", ps)

	dbret, err := dbutil.DBUtil.ExecProc(proc)
	if err != nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if dbret.GetReturnValue() != 1 {
		mlog.Error("Procedure PrGs_SuperTool_JLBPK_UpdVIPUser execute error.")
		return []*msg.Message{dbutil.GetDBRetMsg(int32(dbret.GetReturnValue()))}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
