package baselogic

import (
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// AddSlotGameHistoryLog 新增歷程紀錄
func (logic *SlotLogicBase) AddSlotGameHistoryLog(user *gamelogic.GameUser, gameType int32, logType int32, gameData string) {
	// logic.SendAddSlotGameHistoryLog(user, gameType, logType, gameData) // 舊方法每筆log都送
	logic.SendAddBatchSlotGameHistoryLog(user, gameType, logType, gameData)
}

// AddSlotGameHistoryLogFull 新增歷程紀錄
func (logic *SlotLogicBase) AddSlotGameHistoryLogFull(user *gamelogic.GameUser, beginMoney, betMoney, winMoney int64, gameType int32, logType int32, gameData string) {
	data := new(netproto.Slot_MSG_GameLogAdd)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	data.ServerID = proto.Int32(logic.TableManager.Room.ServerID)

	data.CurrentAmount = proto.Int64(beginMoney)
	data.BetMoney = proto.Int64(betMoney)
	data.WinMoney = proto.Int64(winMoney)
	data.FinalScore = proto.Int64(beginMoney + winMoney - betMoney)
	data.GameType = proto.Int32(gameType)
	data.LogType = proto.Int32(logType)
	data.GameData = proto.String(gameData)

	logic.SendAddBatchSlotGameHistoryLogByMsg(user, data)
}

// AddSlotGameHistoryLogByMsg 新增歷程紀錄 by msg
func (logic *SlotLogicBase) AddSlotGameHistoryLogByMsg(user *gamelogic.GameUser, data *netproto.Slot_MSG_GameLogAdd) {
	logic.SendAddBatchSlotGameHistoryLogByMsg(user, data)
}

// SendAddSlotGameHistoryLog 舊方法每筆log都送
func (logic *SlotLogicBase) SendAddSlotGameHistoryLog(user *gamelogic.GameUser, gameType int32, logType int32, gameData string) {
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID)

	data := new(netproto.Slot_MSG_GameLogAdd)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	data.ServerID = proto.Int32(logic.TableManager.Room.ServerID)

	data.CurrentAmount = proto.Int64(user.Score)
	data.BetMoney = proto.Int64(logic.GamePayMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney)
	data.FinalScore = proto.Int64(user.Score + logic.GameWinMoney - logic.GamePayMoney)
	data.GameType = proto.Int32(gameType)
	data.LogType = proto.Int32(logType)
	data.GameData = proto.String(gameData)

	pmsg.MsgData = data
	logic.SendRPCMessageFunc(true, pmsg)
}

// SendAddBatchSlotGameHistoryLog 新增歷程紀錄(多筆)
func (logic *SlotLogicBase) SendAddBatchSlotGameHistoryLog(user *gamelogic.GameUser, gameType int32, logType int32, gameData string) {
	if logic.processAddLog(user, gameType, logType, gameData) {
		logic.SendCurrentSlotGameHistoryLogs()
	}
}

// SendAddBatchSlotGameHistoryLogByMsg 新增歷程紀錄(多筆)
func (logic *SlotLogicBase) SendAddBatchSlotGameHistoryLogByMsg(user *gamelogic.GameUser, data *netproto.Slot_MSG_GameLogAdd) {
	if logic.processAddLogByMsg(user, data) {
		logic.SendCurrentSlotGameHistoryLogs()
	}
}

func (logic *SlotLogicBase) processAddLog(user *gamelogic.GameUser, gameType int32, logType int32, gameData string) bool {
	var data = new(netproto.Slot_MSG_GameLogAdd)
	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	data.ServerID = proto.Int32(logic.TableManager.Room.ServerID)

	data.CurrentAmount = proto.Int64(user.Score)
	data.BetMoney = proto.Int64(logic.GamePayMoney)
	data.WinMoney = proto.Int64(logic.GameWinMoney)
	data.FinalScore = proto.Int64(user.Score + logic.GameWinMoney - logic.GamePayMoney)
	data.GameType = proto.Int32(gameType)
	data.LogType = proto.Int32(logType)
	data.GameData = proto.String(gameData)
	data.CreateTime = proto.String(time.Now().Format("2006-01-02 15:04:05.000"))

	logic.saveHistoryMutex.Lock()
	logic.cacheLog = append(logic.cacheLog, data)
	logic.saveHistoryMutex.Unlock()

	if logic.isUpdateDB() {
		return true
	}
	return false
}

func (logic *SlotLogicBase) processAddLogByMsg(user *gamelogic.GameUser, data *netproto.Slot_MSG_GameLogAdd) bool {
	if data == nil {
		return false
	}

	data.CreateTime = proto.String(time.Now().Format("2006-01-02 15:04:05.000"))

	logic.saveHistoryMutex.Lock()
	logic.cacheLog = append(logic.cacheLog, data)
	logic.saveHistoryMutex.Unlock()

	if logic.isUpdateDB() {
		return true
	}
	return false
}

// SendCurrentSlotGameHistoryLogs 如果有緩存發送SlotGameHistoryLogs
func (logic *SlotLogicBase) SendCurrentSlotGameHistoryLogs() {

	if logic.isCacheLogs() {
		return
	}

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID)
	batchData := new(netproto.Slot_MSG_BatchGameLogAdd)

	logic.saveHistoryMutex.RLock()
	for _, v := range logic.cacheLog {
		batchData.Logs = append(batchData.Logs, v)
	}
	logic.saveHistoryMutex.RUnlock()

	pmsg.MsgData = batchData
	logic.SendRPCMessageFunc(true, pmsg)
	logic.clearCacheLog()
}

func (logic *SlotLogicBase) isCacheLogs() bool {
	logic.saveHistoryMutex.RLock()
	defer logic.saveHistoryMutex.RUnlock()
	return len(logic.cacheLog) <= 0
}

func (logic *SlotLogicBase) isUpdateDB() bool {
	logic.saveHistoryMutex.RLock()
	defer logic.saveHistoryMutex.RUnlock()
	if len(logic.cacheLog) >= logic.LogicBase.GetServerConfig().SettleRound {
		return true
	}
	return false
}

func (logic *SlotLogicBase) clearCacheLog() {
	logic.saveHistoryMutex.Lock()
	defer logic.saveHistoryMutex.Unlock()
	logic.cacheLog = make([]*netproto.Slot_MSG_GameLogAdd, 0)
}

// ReadSlotGameHistoryLog 查詢歷程紀錄
func (logic *SlotLogicBase) ReadSlotGameHistoryLog(user *gamelogic.GameUser, page, size int32, start, end string) {

	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID)

	data := new(netproto.Slot_MSG_GameLogRead)

	data.UserID = proto.Int32(user.UserID)
	data.GameID = proto.Int32(logic.TableManager.Room.GameInfo.GameID)
	data.StartTime = proto.String(start)
	data.EndTime = proto.String(end)
	data.PageIndex = proto.Int32(page)
	data.PageSize = proto.Int32(size)

	pmsg.MsgData = data

	logic.SendRPCMessageFunc(true, pmsg)
}

// RecvSlotGameLog 接收db server回傳的SlotGameLog資訊，直接傳給Client
func (logic *SlotLogicBase) RecvSlotGameLog(m *msg.Message) {

	defer logic.CheckWaitGroup()
	mlog.Debug("RecvSlotGameLog")

	if logic.GameUser == nil {
		return
	}

	// 解析出來的detail
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_GameLogRecord{}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_GameLogRecord)

	if detail.GetUserID() != logic.GameUser.UserID || detail.GetGameID() != logic.TableManager.Room.GameInfo.GameID {
		return
	}

	logic.SendGameLogMessageToUser(logic.GameUser, detail.GetUserID(), detail.GetGameID(), detail.GetPageCount(), detail.GetLogs())
}
