package jdesydLogic

import (
	"encoding/json"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"strconv"
)

func (pGame *JdesydLogic) OnRecvRPCMessage(pMsg *msg.Message) {
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
func (pGame *JdesydLogic) GetSysPoolWinLose() {
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
func (pGame *JdesydLogic) recvSysPoolWinLose(pMsg *msg.Message) {
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
		g_pSysPoolWinLose.i64PoolDifference = 0
	} else {
		g_pSysPoolWinLose.i64PoolDifference = i64PoolAmount - g_i64WinLine
	}

	// 遊戲輸贏結果
	pGameData := new(GameDataDB)
	err := json.Unmarshal([]byte(data.GetWinLoseData()), pGameData)
	if err == nil {
		g_pSysPoolWinLose.GameDBData.Copy(pGameData)
	} else {
		g_pSysPoolWinLose.GameDBData.AverageBet = g_i64AverageBet
		g_pSysPoolWinLose.GameDBData.Factor = g_iFactor
	}

	updateDealerSupplement(true)

	// 以下為 debug 用
	mlog.Error("recvSysPoolWinLose()")
	mlog.Warn("i64PoolAmount[%d], i64PoolDifference[%d]", i64PoolAmount, g_pSysPoolWinLose.i64PoolDifference)
	mlog.Warn("［GameData］%s", g_pSysPoolWinLose.GameDBData.LogString())

	g_bGetWinLoseData = true
}

// 更新 遊戲輸贏結果
func (pGame *JdesydLogic) UpdateSysPoolWinLose() {
	if !g_bEnablePoolControl {
		return
	}

	i64PoolAmount := g_i64WinLine + g_pSysPoolWinLose.i64PoolDifference
	mlog.Debug("更新 WinLose DB 資料: WinLine[%d], PoolAmount[%d], i64PoolDifference[%d], GameData[%s]",
		g_i64WinLine, i64PoolAmount, g_pSysPoolWinLose.i64PoolDifference, g_pSysPoolWinLose.GameDBData.LogString())

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)

	data := new(netproto.SuperTool_Msg_UpdateSysPoolWinLose)
	data.GameID = proto.Int32(pGame.TableManager.Room.GameInfo.GameID) // 遊戲編號
	FlagID, _ := strconv.Atoi(pGame.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))                       // 房間編號
	data.PoolAmount = proto.Int64(i64PoolAmount) // 水線
	data.WinLine = proto.Int64(g_i64WinLine)                       // 贏錢線
	// 遊戲輸贏結果
	byGameData, _ := json.Marshal(g_pSysPoolWinLose.GameDBData)
	data.WinLoseData = proto.String(string(byGameData))

	pMsg.MsgData = data

	pGame.SendRPCMessageFunc(true, pMsg)
}
