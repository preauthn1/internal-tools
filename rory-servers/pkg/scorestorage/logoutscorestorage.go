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

// LogoutUserScoreStorage 退房間寫分
type LogoutUserScoreStorage struct {
	room       *gamelogic.Room // 房間信息
	dbMsgQueue *rpcmsg.Queue   // 數據庫遠程消息隊列
}

// NewLogoutUserScoreStorage 創建退房間寫分對象
func NewLogoutUserScoreStorage(room *gamelogic.Room, dbMsgQueue *rpcmsg.Queue) *LogoutUserScoreStorage {
	uss := new(LogoutUserScoreStorage)
	uss.room = room
	uss.dbMsgQueue = dbMsgQueue

	return uss
}

// UpdateUserScore 寫分
func (uss *LogoutUserScoreStorage) UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	user.Score += score
	user.UpdateScore += score
	user.TaxAmount += tax
	// user.Coupon += coupon
	// user.UpdateCoupon += coupon
}

// UpdateExperienceScore 體驗金幣寫分
func (uss *LogoutUserScoreStorage) UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) {
	user.ExperienceCash += score
}

// UpdateGameJackpotPool 遊戲獎金池寫分
func (uss *LogoutUserScoreStorage) UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64) {
	if amount <= 0 {
		return
	}

}

// OnUserLeave 登出房間
func (uss *LogoutUserScoreStorage) OnUserLeave(user *gamelogic.GameUser) {
	logoutRoomMsg := new(msg.Message)
	logoutRoomMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	logoutRoomMsg.SClassID = int32(netproto.DBServerClassID_LogoutRoomID)

	logoutRoomData := new(netproto.LogoutGameRoomInfo)
	logoutRoomData.UserID = proto.Int32(user.UserID)
	logoutRoomData.UpdateAmount = proto.Int64(user.UpdateScore)
	logoutRoomData.TaxAmount = proto.Int64(user.TaxAmount)
	logoutRoomData.WinCount = proto.Int32(user.UpdateWinCount)
	logoutRoomData.LoseCount = proto.Int32(user.UpdateLoseCount)
	logoutRoomData.DrawCount = proto.Int32(user.UpdateDrawCount)
	logoutRoomData.LoginTime = proto.String(util.FormatTimeStr(user.LoginTime))
	logoutRoomData.LogoutTime = proto.String(util.FormatTimeStr(time.Now()))
	logoutRoomData.IPAddress = proto.String(user.ConnInfo.IPAddress)
	logoutRoomData.HDCode = proto.String(user.HDCode)
	logoutRoomData.HDType = proto.Int32(user.HDType)
	logoutRoomData.GameTime = proto.Int32(user.GameTime)
	logoutRoomData.GameID = proto.Int32(uss.room.GameInfo.GameID)
	logoutRoomData.ServerID = proto.Int32(uss.room.ServerID)
	logoutRoomData.LoginMoney = proto.Int64(user.LoginScore)
	logoutRoomData.ScoreUpdateMode = proto.Int32(gamelogic.ScoreUpdateModeLogout)
	logoutRoomData.TradeId = proto.Int64(user.TradeID)
	// logoutRoomData.CouponAmount = proto.Int64(user.UpdateCoupon)

	logoutRoomMsg.MsgData = logoutRoomData

	uss.dbMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(uss.room.ServerID), true, logoutRoomMsg)

}

// OnUserRecover ?
func (uss *LogoutUserScoreStorage) OnUserRecover(user *gamelogic.GameUser, dbscore int64) {
	user.Score = dbscore + user.UpdateScore
}

// OnUserMoneyAdd ?
func (uss *LogoutUserScoreStorage) OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64) {
	user.InGamingMoneyChange += changeScore
}

// OnUserReady 在遊戲結束時，每個用戶都會做的事
func (uss *LogoutUserScoreStorage) OnUserReady(user *gamelogic.GameUser) {}
