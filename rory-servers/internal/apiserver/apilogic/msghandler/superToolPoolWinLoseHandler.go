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

// SuperToolPoolWinLoseHandler 加載水池資料
type SuperToolPoolWinLoseHandler struct{}

// IsHook ?
func (h *SuperToolPoolWinLoseHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose) ||
		sclsid == int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)
}

// OnMessage ?
func (h *SuperToolPoolWinLoseHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.GameMessageClassID(m.SClassID) {
	case netproto.GameMessageClassID_SuperToolGetSysPoolWinLose:
		return h.getSysPoolWinLose(ctx, clt, m)
	case netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose:
		return h.updateSysPoolWinLose(ctx, clt, m)
	}

	return []*msg.Message{}
}

func (h *SuperToolPoolWinLoseHandler) getSysPoolWinLose(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperTool_Msg_GetSysPoolWinLose{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperTool_Msg_GetSysPoolWinLose)
	ok, poolBody := service.SuperToolPoolService.GetWinLoss(map[string]interface{}{
		"gameID": mm.GetGameID(),
		"flagID": mm.GetFlagID(),
	})

	if !ok && poolBody.Data.GameID != mm.GetGameID() {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}
	retData := new(netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	retData.GameID = proto.Int32(poolBody.Data.GameID)
	retData.FlagID = proto.Int32(poolBody.Data.FlagID)
	retData.PoolAmount = proto.Int64(util.MoneyToCent(poolBody.Data.PoolAmount))
	retData.WinLine = proto.Int64(util.MoneyToCent(poolBody.Data.WinLine))
	retData.WinLoseData = proto.String(poolBody.Data.GameData)

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet)
	retMsg.MsgData = retData

	return []*msg.Message{retMsg}
}
func (h *SuperToolPoolWinLoseHandler) updateSysPoolWinLose(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.SuperTool_Msg_UpdateSysPoolWinLose{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.SuperTool_Msg_UpdateSysPoolWinLose)
	ok, poolBody := service.SuperToolPoolService.PostWinLoss(map[string]interface{}{
		"gameID":           mm.GetGameID(),
		"flagID":           mm.GetFlagID(),
		"updatePoolAmount": util.CentToMoney(mm.GetPoolAmount()),
		"gameData":         mm.GetWinLoseData(),
	})

	if !ok && poolBody.Data.GameID != mm.GetGameID() {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	retData := new(netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	retData.GameID = proto.Int32(poolBody.Data.GameID)
	retData.FlagID = proto.Int32(poolBody.Data.FlagID)
	retData.PoolAmount = proto.Int64(util.MoneyToCent(poolBody.Data.PoolAmount))
	retData.WinLine = proto.Int64(util.MoneyToCent(poolBody.Data.WinLine))
	retData.WinLoseData = proto.String(poolBody.Data.GameData)

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet)
	retMsg.MsgData = retData

	return []*msg.Message{retMsg}
}

//
//func (h *LoadBlockIPListHandler) requestAPI(method string, path string, params map[string]interface{}, headers ...string ) (result *LoadBlockIPListBody, err error) {
//
//	var groupRequest = new(httputil.WaitGroupRequest)
//	groupRequest.AddTask("0", method, path, params, headers...)
//	var respnose = groupRequest.Send()
//	groupRequest.Output() // 印資訊
//
//	body, err := respnose["0"].Body, respnose["0"].Err
//	if err != nil {
//		return nil, err
//	}
//	var model LoadBlockIPListBody
//	err = json.Unmarshal([]byte(body), &model)
//	result = &model
//	return
//}
