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

// NowUserScoreStorage 即時寫分
type NowUserScoreStorage struct {
	room       *gamelogic.Room // 房間信息
	dbMsgQueue *rpcmsg.Queue   // 數據庫遠程消息隊列
	rpcInfo    *netproto.RPCInfo
}

// NewNowUserScoreStorage 創建即時寫分對象
func NewNowUserScoreStorage(room *gamelogic.Room, dbMsgQueue *rpcmsg.Queue) *NowUserScoreStorage {
	uss := new(NowUserScoreStorage)
	uss.room = room
	uss.dbMsgQueue = dbMsgQueue
	uss.rpcInfo = netproto.CreateEmptyRPCInfo(uss.room.ServerID)

	return uss
}

// UpdateUserScore 寫分 : 用戶資料, 加減分數, 稅額
func (uss *NowUserScoreStorage) UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	if score == 0 && tax == 0 {
		return
	}

	user.Score += score
	user.UpdateScore += score
	user.TaxAmount += tax

	updateScoreMsg := new(msg.Message)
	updateScoreMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	updateScoreMsg.SClassID = int32(netproto.DBServerClassID_UpdateUserScoreID)

	updateScoreData := new(netproto.UpdateUserScore)
	updateScoreData.UserID = proto.Int32(user.UserID)
	updateScoreData.Score = proto.Int64(score)
	updateScoreData.ServerID = proto.Int32(uss.room.ServerID)
	updateScoreData.TaxAmount = proto.Int64(tax)
	updateScoreData.TradeID = proto.Int64(user.TradeID)

	updateScoreMsg.MsgData = updateScoreData

	uss.dbMsgQueue.PushMessage(uss.rpcInfo, true, updateScoreMsg)
	//uss.afterUpdateUserScore(user)
}

// UpdateExperienceScore 寫分 : 用戶資料, 加減分數, 稅額
func (uss *NowUserScoreStorage) UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) {
	updateScoreMsg := new(msg.Message)
	updateScoreMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	updateScoreMsg.SClassID = int32(netproto.DBServerClassID_UpdateExperenceScore)

	updateScoreData := new(netproto.UpdateUserScore)
	updateScoreData.UserID = proto.Int32(user.UserID)
	updateScoreData.Score = proto.Int64(score)
	updateScoreData.ServerID = proto.Int32(uss.room.ServerID)
	updateScoreData.TaxAmount = proto.Int64(tax)
	// updateScoreData.CouponAmount = proto.Int64(coupon)

	updateScoreMsg.MsgData = updateScoreData

	uss.dbMsgQueue.PushMessage(uss.rpcInfo, true, updateScoreMsg)
}

// UpdateGameJackpotPool 遊戲獎金池寫分
func (uss *NowUserScoreStorage) UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64) {
	if amount <= 0 {
		return
	}
	updateScoreMsg := new(msg.Message)
	updateScoreMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	updateScoreMsg.SClassID = int32(netproto.DBServerClassID_UpdateGameJackpotScore)

	updateScoreData := new(netproto.UpdateUserScore)
	updateScoreData.UserID = proto.Int32(user.UserID)
	updateScoreData.Score = proto.Int64(amount)
	updateScoreData.ServerID = proto.Int32(uss.room.ServerID)
	updateScoreData.TaxAmount = proto.Int64(0)

	updateScoreMsg.MsgData = updateScoreData

	uss.dbMsgQueue.PushMessage(uss.rpcInfo, true, updateScoreMsg)
}

// OnUserLeave 退房
func (uss *NowUserScoreStorage) OnUserLeave(user *gamelogic.GameUser) {
	logoutRoomMsg := new(msg.Message)
	logoutRoomMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	logoutRoomMsg.SClassID = int32(netproto.DBServerClassID_LogoutRoomID)

	logoutRoomData := new(netproto.LogoutGameRoomInfo)
	logoutRoomData.UserID = proto.Int32(user.UserID)
	logoutRoomData.UpdateAmount = proto.Int64(user.UpdateScore)
	logoutRoomData.TaxAmount = proto.Int64(user.TaxAmount)
	logoutRoomData.WinCount = proto.Int32(user.UpdateWinCount)
	logoutRoomData.LoseCount = proto.Int32(user.UpdateLoseCount)
	logoutRoomData.LastWinLose = proto.Int32(user.LastWinLose)
	logoutRoomData.DrawCount = proto.Int32(user.UpdateDrawCount)
	logoutRoomData.LoginTime = proto.String(util.FormatTimeStr(user.LoginTime))
	logoutRoomData.LogoutTime = proto.String(util.FormatTimeStr(time.Now()))
	logoutRoomData.FlowAmount = proto.Int64(user.FlowAmount)
	logoutRoomData.BetAmount = proto.Int64(user.BetAmount)
	logoutRoomData.PayAmount = proto.Int64(user.PayAmount)

	// logoutRoomData.CouponAmount = proto.Int64(user.UpdateCoupon) // 應該要寫進入到離開使用的Coupon量
	//登录成功后马上断开的情况
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
	logoutRoomData.ScoreUpdateMode = proto.Int32(gamelogic.ScoreUpdateModeNow)
	logoutRoomData.TradeId = proto.Int64(user.TradeID)
	logoutRoomMsg.MsgData = logoutRoomData

	uss.dbMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(uss.room.ServerID), true, logoutRoomMsg)
}

// OnUserRecover 斷線重連時會呼叫此函式記錄斷線重連時的分數異動
func (uss *NowUserScoreStorage) OnUserRecover(user *gamelogic.GameUser, dbscore int64) {
	//changeScore := dbscore - user.Score
	//user.InGamingMoneyChange = changeScore
}

// OnUserMoneyAdd 即時到帳, 有可能一局中要Add兩次以上
func (uss *NowUserScoreStorage) OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64) {
	user.InGamingMoneyChange += changeScore
}

// OnUserReady 在遊戲結束時，每個用戶都會做的事
func (uss *NowUserScoreStorage) OnUserReady(user *gamelogic.GameUser) {
	uss.beforeGameStart(user)
}

// notifyUserStateChange 通知用戶狀態改變
func (uss *NowUserScoreStorage) notifyUserStateChange(user *gamelogic.GameUser) {
	if user.IsOffLine() {
		return
	}
	// 如果用戶是機器人
	if user.IsRobot() {
		return
	}

	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_UserStateChangeID)

	ust := new(netproto.UserStateChange)
	ust.User = user.GetProtoBufGameUser()

	m.MsgData = ust

	uss.room.RoomUserManager.SendMessage(user.ConnInfo, m)
}

// 正在遊戲中有需要異動金額的部份等到下局開始時才生效
func (uss *NowUserScoreStorage) beforeGameStart(user *gamelogic.GameUser) {
	if user.InGamingMoneyChange != 0 {
		//msgMoney := user.InGamingMoneyChange
		user.Score += user.InGamingMoneyChange
		user.LoginScore += user.InGamingMoneyChange
		user.InGamingMoneyChange = 0
		if user.IsOffLine() {
			//用戶不在線，回寫資料表?
			return
		}
		//msg, err := roomutil.GetTipMessage(fmt.Sprintf("亲爱的玩家:\n您充值的金币:%d已到帐\n( 如有任何疑问请联系客服)", msgMoney / 100), roomutil.MsgTypeID_TIP)
		//if err != nil {
		//	mlog.Error("打包客戶端消息失敗 %v", err)
		//	return
		//}
		//uss.room.RoomUserManager.SendMessage(user.ConnInfo, msg)

		uss.notifyUserStateChange(user)
	}
}
