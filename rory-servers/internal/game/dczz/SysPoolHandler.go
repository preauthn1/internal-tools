package dczz

import (
	"encoding/json"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"strconv"
)

func (pGame *Logic) OnRecvRPCMessage(pMsg *msg.Message) {
	if pMsg == nil {
		return
	}

	rm, ok := pMsg.MsgData.(*rpcmsg.RPCMessage)
	if !ok {
		return
	}

	mmbyte := rm.MsgData.([]byte)
	if len(mmbyte) == 0 {
		return
	}

	switch pMsg.SClassID {
	case int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLoseRet):
		pGame.recvSysPoolWinLose(pMsg)
	}
}

// 要求 遊戲輸贏結果
func (pGame *Logic) GetSysPoolWinLose() {
	if !g_bEnablePoolControl {
		return
	}

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetSysPoolWinLose)

	data := new(netproto.SuperTool_Msg_GetSysPoolWinLose)
	FlagID, _ := strconv.Atoi(pGame.TableManager.Room.Flag)
	data.GameID = proto.Int32(pGame.TableManager.Room.GameInfo.GameID)
	data.FlagID = proto.Int32(int32(FlagID))

	pMsg.MsgData = data

	pGame.SendRPCMessageFunc(true, pMsg)
}

// 回傳 遊戲輸贏結果
func (pGame *Logic) recvSysPoolWinLose(pMsg *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(pMsg, func() proto.Message {
		return new(netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	})

	if !g_bEnablePoolControl || g_bGetWinLoseData {
		return
	}

	data := rmm.(*netproto.SuperTool_Msg_GetSysPoolWinLoseRet)

	if data.GetGameID() != pGame.TableManager.Room.GameInfo.GameID && strconv.Itoa(int(data.GetFlagID())) != pGame.TableManager.Room.Flag {
		return
	}

	// 水線
	i64PoolAmount := int64(data.GetPoolAmount())
	if i64PoolAmount <= 0 {
		g_pSysPoolWinLose.i64PoolAmount = g_i64WinLine
	} else {
		g_pSysPoolWinLose.i64PoolAmount = i64PoolAmount
	}

	// 遊戲輸贏結果
	pGameData := new(GameDataDB)
	err := json.Unmarshal([]byte(data.GetWinLoseData()), pGameData)
	if err == nil {
		g_pSysPoolWinLose.GameDBData.Copy(pGameData)
	} else {
		g_pSysPoolWinLose.GameDBData.ChangeCardRate = DefaultChangeCardRate	// 換牌機率，採百分比
		g_pSysPoolWinLose.GameDBData.DealerWeight = DealerWeightDefault	// 莊家權重
	}

	// 以下為 debug 用
	mlog.Debug("recvSysPoolWinLose()")
	mlog.Debug("i64PoolAmount[%d]", g_pSysPoolWinLose.i64PoolAmount)
	mlog.Debug("［GameData］%s", g_pSysPoolWinLose.GameDBData.LogString())

	g_bGetWinLoseData = true
}

// 更新 遊戲輸贏結果
func (pGame *Logic) UpdateSysPoolWinLose() {
	if !g_bEnablePoolControl {
		return
	}

	mlog.Debug("更新 WinLose DB 資料: WinLine[%d], PoolAmount[%d], GameData[%s]",
		g_i64WinLine, g_pSysPoolWinLose.i64PoolAmount, g_pSysPoolWinLose.GameDBData.LogString())

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)

	data := new(netproto.SuperTool_Msg_UpdateSysPoolWinLose)
	data.GameID = proto.Int32(pGame.TableManager.Room.GameInfo.GameID) // 遊戲編號
	FlagID, _ := strconv.Atoi(pGame.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))                       // 房間編號
	data.WinLine = proto.Int64(g_i64WinLine)                       // 贏錢線
	data.PoolAmount = proto.Int64(g_pSysPoolWinLose.i64PoolAmount) // 水線
	// 遊戲輸贏結果
	byGameData, _ := json.Marshal(g_pSysPoolWinLose.GameDBData)
	data.WinLoseData = proto.String(string(byGameData))

	pMsg.MsgData = data

	pGame.SendRPCMessageFunc(true, pMsg)
}
