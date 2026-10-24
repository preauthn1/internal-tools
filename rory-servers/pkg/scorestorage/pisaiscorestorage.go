package scorestorage

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// PiSaiScoreStorage 比賽寫積分
type PiSaiScoreStorage struct {
	room       *gamelogic.Room // 房間信息
	dbMsgQueue *rpcmsg.Queue   // 數據庫遠程消息隊列
	rpcInfo    *netproto.RPCInfo
}

// NewPiSaiScoreStorage 創建比賽對象
func NewPiSaiScoreStorage(room *gamelogic.Room, dbMsgQueue *rpcmsg.Queue) *PiSaiScoreStorage {
	uss := new(PiSaiScoreStorage)
	uss.room = room
	uss.dbMsgQueue = dbMsgQueue
	uss.rpcInfo = netproto.CreateEmptyRPCInfo(uss.room.ServerID)
	return uss
}

// UpdateUserScore ?
func (uss *PiSaiScoreStorage) UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	if score == 0 && tax == 0 {
		return
	}

	//user.Score += score
	//user.UpdateScore += score
	//user.TaxAmount += tax

	updateScoreMsg := new(msg.Message)
	updateScoreMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	updateScoreMsg.SClassID = int32(netproto.DBServerClassID_UpdatePiSaiScoreID)

	updateScoreData := new(netproto.UpdateUserScore)
	updateScoreData.ContestID = proto.Int32(user.ContestID)
	updateScoreData.UserID = proto.Int32(user.UserID)
	updateScoreData.Score = proto.Int64(score)
	updateScoreData.ServerID = proto.Int32(uss.room.ServerID)
	updateScoreData.TaxAmount = proto.Int64(tax)
	// updateScoreData.CouponAmount = proto.Int64(0)

	updateScoreMsg.MsgData = updateScoreData

	uss.dbMsgQueue.PushMessage(uss.rpcInfo, true, updateScoreMsg)
}

// UpdateGameJackpotPool ?
func (uss *PiSaiScoreStorage) UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64) {}

// OnUserLeave ?
func (uss *PiSaiScoreStorage) OnUserLeave(user *gamelogic.GameUser) {
	logoutRoomMsg := new(msg.Message)
	logoutRoomMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	logoutRoomMsg.SClassID = int32(netproto.DBServerClassID_LogoutRoomID)

	logoutRoomData := new(netproto.LogoutGameRoomInfo)
	logoutRoomData.UserID = proto.Int32(user.UserID)
	logoutRoomData.UpdateAmount = proto.Int64(0)
	logoutRoomData.TaxAmount = proto.Int64(0)
	logoutRoomData.WinCount = proto.Int32(0)
	logoutRoomData.LoseCount = proto.Int32(0)
	logoutRoomData.LastWinLose = proto.Int32(0)
	logoutRoomData.DrawCount = proto.Int32(user.UpdateDrawCount)
	logoutRoomData.LoginTime = proto.String(util.FormatTimeStr(user.LoginTime))
	logoutRoomData.LogoutTime = proto.String(util.FormatTimeStr(time.Now()))
	logoutRoomData.FlowAmount = proto.Int64(0)
	// logoutRoomData.CouponAmount = proto.Int64(0)
	// 登錄成功後馬上斷開的情況
	if user.ConnInfo != nil {
		logoutRoomData.IPAddress = proto.String(user.ConnInfo.IPAddress)
	} else {
		logoutRoomData.IPAddress = proto.String("unknown")
	}
	logoutRoomData.HDCode = proto.String(user.HDCode)
	logoutRoomData.HDType = proto.Int32(user.HDType)
	logoutRoomData.GameTime = proto.Int32(user.GameTime)
	logoutRoomData.GameID = proto.Int32(uss.room.GameInfo.GameID)
	logoutRoomData.ServerID = proto.Int32(uss.room.ServerID)
	logoutRoomData.LoginMoney = proto.Int64(user.LoginScore)
	logoutRoomData.ScoreUpdateMode = proto.Int32(gamelogic.ScoreUpdateModeLogout)
	logoutRoomData.TradeId = proto.Int64(user.TradeID)
	logoutRoomMsg.MsgData = logoutRoomData

	uss.dbMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(uss.room.ServerID), true, logoutRoomMsg)
}

// OnUserRecover ?
func (uss *PiSaiScoreStorage) OnUserRecover(user *gamelogic.GameUser, dbscore int64) {}

// OnUserMoneyAdd ?
func (uss *PiSaiScoreStorage) OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64) {}

// OnUserReady ?
func (uss *PiSaiScoreStorage) OnUserReady(user *gamelogic.GameUser) {}

// UpdateExperienceScore ?
func (uss *PiSaiScoreStorage) UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) {
}
