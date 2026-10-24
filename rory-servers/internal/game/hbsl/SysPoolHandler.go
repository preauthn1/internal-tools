package hbsl

import (
	"encoding/json"
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"strconv"
)

// OnRecvRPCMessage 接收 RPC message
func (pGame *HbslLogic) OnRecvRPCMessage(pMsg *msg.Message) {
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

// GetSysPoolWinLose 要求 遊戲輸贏結果
func (pGame *HbslLogic) GetSysPoolWinLose() {
	if !gbEnablePoolControl {
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
func (pGame *HbslLogic) recvSysPoolWinLose(pMsg *msg.Message) {
	_, rmm := rpcmsg.ParseRPCMessage(pMsg, func() proto.Message {
		return new(netproto.SuperTool_Msg_GetSysPoolWinLoseRet)
	})

	if !gbEnablePoolControl || gbGetWinLoseData {
		return
	}

	data := rmm.(*netproto.SuperTool_Msg_GetSysPoolWinLoseRet)

	if data.GetGameID() != pGame.TableManager.Room.GameInfo.GameID && strconv.Itoa(int(data.GetFlagID())) != pGame.TableManager.Room.Flag {
		return
	}

	// 水線
	i64PoolAmount := int64(data.GetPoolAmount())
	if i64PoolAmount <= 0 {
		gpSysPoolWinLose.i64PoolAmount = gi64WinLine
	} else {
		gpSysPoolWinLose.i64PoolAmount = i64PoolAmount
	}

	// 遊戲輸贏結果
	pGameData := new(GameDataDB)
	err := json.Unmarshal([]byte(data.GetWinLoseData()), pGameData)
	if err == nil {
		gpSysPoolWinLose.GameDBData.Copy(pGameData)
	}

	// 以下為 debug 用
	mlog.Debug("recvSysPoolWinLose()")
	mlog.Debug("i64PoolAmount[%d]", gpSysPoolWinLose.i64PoolAmount)
	mlog.Debug("［GameData］%s", gpSysPoolWinLose.GameDBData.LogString())

	gbGetWinLoseData = true
}

// UpdateSysPoolWinLose 更新 遊戲輸贏結果
func (pGame *HbslLogic) UpdateSysPoolWinLose() {
	if !gbEnablePoolControl {
		return
	}

	mlog.Debug("更新 WinLose DB 資料: WinLine[%d], PoolAmount[%d], GameData[%s]",
		gi64WinLine, gpSysPoolWinLose.i64PoolAmount, gpSysPoolWinLose.GameDBData.LogString())

	pMsg := new(msg.Message)
	pMsg.BClassID = int32(netproto.MessageBClassID_Game)
	pMsg.SClassID = int32(netproto.GameMessageClassID_SuperToolUpdateSysPoolWinLose)

	data := new(netproto.SuperTool_Msg_UpdateSysPoolWinLose)
	data.GameID = proto.Int32(pGame.TableManager.Room.GameInfo.GameID) // 遊戲編號
	FlagID, _ := strconv.Atoi(pGame.TableManager.Room.Flag)
	data.FlagID = proto.Int32(int32(FlagID))                       // 房間編號
	data.WinLine = proto.Int64(gi64WinLine)                       // 贏錢線
	data.PoolAmount = proto.Int64(gpSysPoolWinLose.i64PoolAmount) // 水線
	// 遊戲輸贏結果
	byGameData, _ := json.Marshal(gpSysPoolWinLose.GameDBData)
	data.WinLoseData = proto.String(string(byGameData))

	pMsg.MsgData = data

	pGame.SendRPCMessageFunc(true, pMsg)
}
