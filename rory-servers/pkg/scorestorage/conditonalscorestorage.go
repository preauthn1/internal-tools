package scorestorage

import (
	"sync"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// Check Type
const (
	CheckByGameRound   = 1 // 遊戲局數
	CheckBySettleScore = 2 // 結算分數
	// Check_LastUpdateTime = 3 // 上次更新時間
)

// UpdateUserScoreData ?
type UpdateUserScoreData struct {
	UserID       int32
	UpdateAmount int64
	TaxAmount    int64
	GameRound    int
}

// ConditionalUserScoreStorage 滿足設定條件才寫分
type ConditionalUserScoreStorage struct {
	room       *gamelogic.Room // 房間信息
	dbMsgQueue *rpcmsg.Queue   // 數據庫遠程消息隊列
	rpcInfo    *netproto.RPCInfo
	serverCfg  *gamelogic.GameServerConfig
	// key: check type
	// value: check function
	dbUpdateChecker map[int](func(data *UpdateUserScoreData) bool)

	// check used data
	usrScoreData sync.Map // key: userid value:data
}

// NewConditionalUserScoreStorage 創建滿足設定條件寫分對象
func NewConditionalUserScoreStorage(room *gamelogic.Room, dbMsgQueue *rpcmsg.Queue, cfg *gamelogic.GameServerConfig) *ConditionalUserScoreStorage {
	uss := new(ConditionalUserScoreStorage)
	uss.room = room
	uss.dbMsgQueue = dbMsgQueue
	uss.serverCfg = cfg
	uss.dbUpdateChecker = map[int](func(data *UpdateUserScoreData) bool){
		CheckByGameRound:   uss.checkGameRound,
		CheckBySettleScore: uss.checkSettleScore,
	}
	uss.rpcInfo = netproto.CreateEmptyRPCInfo(uss.room.ServerID)

	return uss
}

// UpdateUserScore 寫分 : 用戶資料, 加減分數, 稅額
func (uss *ConditionalUserScoreStorage) UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	if score == 0 && tax == 0 {
		return
	}

	user.Score += score
	user.UpdateScore += score
	user.TaxAmount += tax

	// 是否滿足DB更新條件
	isUpdateGamelog, data := uss.processUserScore(user.UserID, score, tax)
	// fmt.Println("UpdateUserScore", user.UserID, user.TradeID, user.Score, score, user.UpdateScore, user.TaxAmount)
	if isUpdateGamelog {
		// fmt.Println("OnUserLeave", user.UserID, user.TradeID, user.Score, user.UpdateScore)
		uss.pushRPCUpdateUserScore(user, data.UpdateAmount, data.TaxAmount)
		// 不可以歸零 gameuser 的UpdateScore 那是要寫現金記錄用的
		//user.UpdateScore = 0
		//user.TaxAmount = 0
		uss.resetData(data)
	}

}

// UpdateExperienceScore ?
func (uss *ConditionalUserScoreStorage) UpdateExperienceScore(user *gamelogic.GameUser, score int64, tax int64) {
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
func (uss *ConditionalUserScoreStorage) UpdateGameJackpotPool(user *gamelogic.GameUser, amount int64) {
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

// OnUserLeave ?
func (uss *ConditionalUserScoreStorage) OnUserLeave(user *gamelogic.GameUser) {
	//UpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	// 1. 先更新用戶金幣
	if val, ok := uss.usrScoreData.Load(user.UserID); ok {
		data := val.(*UpdateUserScoreData)
		if data.GameRound > 0 {
			uss.pushRPCUpdateUserScore(user, data.UpdateAmount, data.TaxAmount)
		}
	}

	//fmt.Println("OnUserLeave", user.UserID, user.TradeID, user.Score, user.UpdateScore)
	// 2. 再更新現金記錄
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
	logoutRoomData.ScoreUpdateMode = proto.Int32(gamelogic.ScoreUpdateModeNow)
	logoutRoomData.TradeId = proto.Int64(user.TradeID)
	logoutRoomMsg.MsgData = logoutRoomData

	uss.dbMsgQueue.PushMessage(netproto.CreateEmptyRPCInfo(uss.room.ServerID), true, logoutRoomMsg)
	if val, ok := uss.usrScoreData.Load(user.UserID); ok {
		var data = val.(*UpdateUserScoreData)
		uss.resetData(data)
	}
}

// OnUserRecover 斷線重連時會呼叫此函式記錄斷線重連時的分數異動
func (uss *ConditionalUserScoreStorage) OnUserRecover(user *gamelogic.GameUser, dbscore int64) {
	user.Score = dbscore + user.UpdateScore
}

// OnUserMoneyAdd 即時到帳, 有可能一局中要Add兩次以上
func (uss *ConditionalUserScoreStorage) OnUserMoneyAdd(user *gamelogic.GameUser, changeScore int64) {
	user.InGamingMoneyChange += changeScore
}

// OnUserReady 在遊戲結束時，每個用戶都會做的事
func (uss *ConditionalUserScoreStorage) OnUserReady(user *gamelogic.GameUser) {
	uss.beforeGameStart(user)
}

// 通知用戶狀態改變
func (uss *ConditionalUserScoreStorage) notifyUserStateChange(user *gamelogic.GameUser) {
	if user.IsOffLine() {
		return
	}
	//如果用户是机器人
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
func (uss *ConditionalUserScoreStorage) beforeGameStart(user *gamelogic.GameUser) {
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

func (uss *ConditionalUserScoreStorage) processUserScore(userid int32, updateAmount int64, taxAmount int64) (bool, *UpdateUserScoreData) {
	var scoreData *UpdateUserScoreData
	if val, ok := uss.usrScoreData.Load(userid); ok {
		scoreData = val.(*UpdateUserScoreData)
		scoreData.UpdateAmount += updateAmount
		scoreData.TaxAmount += taxAmount
		scoreData.GameRound++
	} else {
		scoreData = new(UpdateUserScoreData)
		scoreData.UserID = userid
		scoreData.UpdateAmount = updateAmount
		scoreData.TaxAmount = taxAmount
		scoreData.GameRound = 1
	}
	uss.usrScoreData.Store(userid, scoreData)

	// fmt.Println(scoreData)

	if scoreData != nil && uss.isUpdateDB(scoreData) {
		return true, scoreData
	}
	return false, nil
}

func (uss *ConditionalUserScoreStorage) isUpdateDB(data *UpdateUserScoreData) bool {
	for _, checker := range uss.dbUpdateChecker {
		if checker(data) {
			return true
		}
	}
	return false
}

func (uss *ConditionalUserScoreStorage) resetData(data *UpdateUserScoreData) {
	data.UpdateAmount = 0
	data.TaxAmount = 0
	data.GameRound = 0
	uss.usrScoreData.Store(data.UserID, data)
}

// checker
func (uss *ConditionalUserScoreStorage) checkGameRound(data *UpdateUserScoreData) bool {
	if data.GameRound >= uss.serverCfg.SettleRound {
		return true
	}
	return false
}

func (uss *ConditionalUserScoreStorage) checkSettleScore(data *UpdateUserScoreData) bool {
	if data.UpdateAmount >= uss.serverCfg.SettleScore {
		return true
	}
	return false
}

//rpc update user score
func (uss *ConditionalUserScoreStorage) pushRPCUpdateUserScore(user *gamelogic.GameUser, score int64, tax int64) {
	updateScoreMsg := new(msg.Message)
	updateScoreMsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	updateScoreMsg.SClassID = int32(netproto.DBServerClassID_UpdateUserScoreID)

	updateScoreData := new(netproto.UpdateUserScore)
	updateScoreData.UserID = proto.Int32(user.UserID)
	updateScoreData.Score = proto.Int64(score)
	updateScoreData.ServerID = proto.Int32(uss.room.ServerID)
	updateScoreData.TaxAmount = proto.Int64(tax)

	updateScoreMsg.MsgData = updateScoreData
	uss.dbMsgQueue.PushMessage(uss.rpcInfo, true, updateScoreMsg)
}
