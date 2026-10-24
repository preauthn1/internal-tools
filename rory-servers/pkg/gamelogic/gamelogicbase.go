package gamelogic

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// LogicBase 遊戲邏輯基類
type LogicBase struct {
	BeforeGameStartFunc       []func()                                         // 遊戲開始前調用函數組
	AfterGameStartFunc        []func()                                         // 遊戲開始後調用函數組
	GameStartFunc             func()                                           // 遊戲開始調用函數
	BeforeGameEndFunc         []func()                                         // 遊戲結束前調用函數組
	AfterGameEndFunc          []func()                                         // 遊戲結束後調用函數組
	GameEndFunc               func()                                           // 遊戲結束調用函數
	GameCancelFunc            func(str string)                                 // 遊戲取消調用函數
	BeforePiSaiStartFunc      []func()                                         // 比賽開始前調用函數
	AfterPiSaiStartFunc       []func()                                         // 比賽開始後調用函數
	PiSaiStartFunc            func()                                           // 比賽開始調用函數
	BeforePiSaiEndFunc        []func()                                         // 比賽結束前調用函數
	AfterPiSaiEndFunc         []func()                                         // 比賽結束後調用函數
	PiSaiEndFunc              func()                                           // 比賽結束調用函數
	PiSaiCancelFunc           func(str string)                                 // 比賽取消調用函數
	UpdateUserScoreFunc       func(userid int32, score int64, taxamount int64) // 更新分数
	UpdateExperienceScoreFunc func(userid int32, score int64, taxamount int64) // 體驗金幣
	UpdateGameJackpotPoolFunc func(userid int32, amount int64)                 // 更新遊戲獎金池
	KickUserFunc              func(userid int32, msg string)                   // 踢人
	SendMessgeFunc            func(userid int32, m *msg.Message)               // 發送消息
	SendRPCMessageFunc        func(safe bool, m *msg.Message)                  // 發送RPC消息
	SendFishNotifyFunc        func(m *msg.Message)                             // 發送捕魚公告
	TableManager              *Table                                           // 桌子
	isgaming                  bool                                             // 是否遊戲中
	LastGameEndTime           time.Time                                        // 上一次遊戲結束時間
	GameStartTime             time.Time                                        // 遊戲開始時間
	Configs                   *config.Group                                    // 配置信息
	TimerFuncManager          *event.TimerFuncManager                          // 定時器
	PlayerActionManager       *PlayerActionManager                             // 遊戲動作管理器
	gameCounter               int32                                            // 遊戲計數器,每次遊戲開始後累加
	GameLog                   *GameLogHelper                                   // 牌局記錄
	AsyncCall                 AsyncCallFunc                                    // 異步調用函數
	UserEarlyChangeTable      func(user *GameUser, early string)               // 提前離桌
	GameTimeOutStopSig        chan bool
	SysPool                   *SysPoolHelper // 系統水池工具

	//VipUserProvider           *VipUserProvider //VIP用戶
	//VipUserService            *VipUserService
}

// 遊戲類型
const (
	ChenChien = 0 //真錢
	PiSai     = 1 //比賽活動
	HsiuChien = 2 //休閒
)

// NewGameLogicBase 創建遊戲邏輯基類對象
func NewGameLogicBase(t *Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *LogicBase {
	glcb := new(LogicBase)
	glcb.TableManager = t
	glcb.Configs = cfgs
	glcb.isgaming = false
	glcb.TimerFuncManager = tfManager
	glcb.LastGameEndTime = time.Now()
	glcb.GameStartTime = time.Now()
	glcb.PlayerActionManager = NewPlayerActionManager(tfManager, glcb)
	glcb.gameCounter = 0
	glcb.GameLog = NewGameLogHelper(gamelog.NewEmptyGameLog(), glcb)
	glcb.GameTimeOutStopSig = make(chan bool, 1)
	glcb.SysPool = NewSysPoolHelper(glcb, cfgs)
	return glcb
}

// Summary 此遊戲邏輯目前的狀態描述
func (glc *LogicBase) Summary() string {
	return ""
}

// OnGameStart 遊戲開始
func (glc *LogicBase) OnGameStart(f func()) {
	glc.GameStartFunc = func() {
		if glc.isgaming {
			return
		}

		// 真錢
		if glc.TableManager.Room.GameType == ChenChien {
			// 遊戲開始前調用函數
			for _, regf := range glc.BeforeGameStartFunc {
				regf()
			}

			glc.isgaming = true
			glc.GameStartTime = time.Now()
			glc.PlayerActionManager.Clear()

			// 通知用戶狀態改變
			for _, ur := range glc.AllTableUser() {
				ur.UserState.OnUserGame(ur, glc.TableManager.Room)
				// 廣播用戶狀態改變
				glc.BoardcastUserStateChange(ur)
			}

			// 廣播遊戲開始
			glc.BoardcastGameStart()
			mlog.Debug("遊戲開始...")

			// 牌局記錄
			glc.GameLog.GameBegin()

			f()

			// 遊戲開始後調用函數
			for _, regf := range glc.AfterGameStartFunc {
				regf()
			}

			// 更新遊戲邏輯計數器
			glc.gameCounter++
		}

		// 是比賽
		if glc.TableManager.Room.GameType == PiSai {
			// 遊戲開始前調用函數
			for _, regf := range glc.BeforePiSaiStartFunc {
				regf()
			}

			glc.isgaming = true
			glc.GameStartTime = time.Now()
			glc.PlayerActionManager.Clear()

			// 通知用戶狀態改變
			for _, ur := range glc.AllTableUser() {
				ur.UserState.OnUserGame(ur, glc.TableManager.Room)
				// 廣播用戶狀態改變
				glc.BoardcastUserStateChange(ur)
			}

			// 廣播遊戲開始
			glc.BoardcastGameStart()
			mlog.Debug("遊戲開始...")

			// 牌局記錄
			glc.GameLog.GameBegin()

			f()

			// 遊戲開始後調用函數
			for _, regf := range glc.AfterPiSaiStartFunc {
				regf()
			}

			// 更新遊戲邏輯計數器
			glc.gameCounter++
		}

	}
}

// GameTimeOutKickUser 當遊戲開始後超過30分鐘未結束該局就踢出此張桌子所有人
func (glc *LogicBase) GameTimeOutKickUser() {

	for {

		select {
		case <-glc.GameTimeOutStopSig:
			mlog.Debug("遊戲正常結束, 離開超時線程")
			return
		case <-time.After(time.Duration(glc.GetServerConfig().GameTimeOut) * time.Minute):
			//case <- time.After(time.Duration(glc.GetServerConfig().GameTimeOut) * time.Second):
			mlog.Debug("xxxxxxxx")
			if glc.IsGaming() {
				mlog.Error("遊戲超時異常 取消遊戲 tableno %d", glc.TableNo())
				glc.GameCancelFunc("遊戲超時異常")

				//for _, ur := range glc.AllTableUser() {
				//	mlog.System("遊戲超時異常 踢出userid:%v", ur.UserID)
				//	glc.KickUserFunc(ur.UserID, "游戏超时异常")
				//}

			}
			return
		}
	}
}

// OnGameEnd 遊戲結束
func (glc *LogicBase) OnGameEnd(f func(int32)) {
	glc.GameEndFunc = func() {
		// 遊戲結束前調用函數
		for _, regf := range glc.BeforeGameEndFunc {
			regf()
		}

		glc.isgaming = false
		glc.LastGameEndTime = time.Now()
		glc.PlayerActionManager.Clear()

		// 廣播遊戲結束
		glc.BoardcastGameEnd()

		// 通知用戶狀態改變
		for _, ur := range glc.AllGameUser() {
			ur.UserState.OnGameEnd(ur, glc.TableManager.Room)
			// 廣播用戶狀態改變
			glc.BoardcastUserStateChange(ur)
			//检查准备超时
			if ur.IsSit() {
				glc.checkUserReadyTimer(ur)
			}
		}

		gametime := int32(time.Now().Sub(glc.GameStartTime).Seconds())

		f(gametime)

		// 遊戲結束後調用函數
		for _, regf := range glc.AfterGameEndFunc {
			regf()
		}

		// 廣播旁觀用戶
		glc.boardcastLookOnUser()

		mlog.Debug("遊戲結束...")
	}
}

// OnGameCancel 遊戲取消
func (glc *LogicBase) OnGameCancel(f func(str string)) {
	glc.GameCancelFunc = func(str string) {
		// 遊戲結束前調用函數
		for _, regf := range glc.BeforeGameEndFunc {
			regf()
		}

		glc.isgaming = false
		glc.LastGameEndTime = time.Now()
		glc.PlayerActionManager.Clear()

		// 廣播遊戲結束
		glc.BoardcastGameEnd()

		// 通知用戶狀態改變
		for _, ur := range glc.AllTableUser() {
			ur.UserState.OnGameEnd(ur, glc.TableManager.Room)
			// 廣播用戶狀態改變
			glc.BoardcastUserStateChange(ur)
			// 檢查準備超時
			if ur.IsSit() {
				glc.checkUserReadyTimer(ur)
			}
		}

		f(str)

		// 遊戲結束後調用函數
		for _, regf := range glc.AfterGameEndFunc {
			regf()
		}

		// 廣播旁觀用戶
		glc.boardcastLookOnUser()
		example.invalid("遊戲取消...")
	}
}

// OnUserEarlyChangeTable 提前換桌
func (glc *LogicBase) OnUserEarlyChangeTable(f func(user *GameUser, early string)) {
	glc.UserEarlyChangeTable = f
}

// OnUpdateUserScore 寫分
func (glc *LogicBase) OnUpdateUserScore(f func(userid int32, score int64, taxamount int64)) {
	glc.UpdateUserScoreFunc = f
}

// OnUpdateExperienceScore 體驗金幣寫分
func (glc *LogicBase) OnUpdateExperienceScore(f func(userid int32, score int64, taxamount int64)) {
	glc.UpdateExperienceScoreFunc = f
}

// OnUpdateGameJackpotPool 遊戲獎金池寫分
func (glc *LogicBase) OnUpdateGameJackpotPool(f func(userid int32, amount int64)) {
	glc.UpdateGameJackpotPoolFunc = f
}

// OnKickUser 踢人
func (glc *LogicBase) OnKickUser(f func(userid int32, msg string)) {
	glc.KickUserFunc = f
}

// OnSendMessage 發送消息
func (glc *LogicBase) OnSendMessage(f func(userid int32, m *msg.Message)) {
	glc.SendMessgeFunc = f
}

// GetSendMessageFunc 發送消息 function
func (glc *LogicBase) GetSendMessageFunc() func(userid int32, m *msg.Message) {
	return glc.SendMessgeFunc
}

// OnSendRPCMessage 發送RPC消息
func (glc *LogicBase) OnSendRPCMessage(f func(safe bool, m *msg.Message)) {
	glc.SendRPCMessageFunc = f
}

// GetSendRPCMessageFunc 發送RPC消息 function
func (glc *LogicBase) GetSendRPCMessageFunc() func(safe bool, m *msg.Message) {
	return glc.SendRPCMessageFunc
}

// OnNotifyFishMessage 發送捕魚公告
func (glc *LogicBase) OnNotifyFishMessage(f func(m *msg.Message)) {
	glc.SendFishNotifyFunc = f
}

// OnRecvRPCMessage 收到RPC消息
func (glc *LogicBase) OnRecvRPCMessage(m *msg.Message) {}

// OnRecvHTTPMessage 收到Http消息
func (glc *LogicBase) OnRecvHTTPMessage(req *http.Request) proto.Message {
	return nil
}

// OnAsyncCall 異步調用
func (glc *LogicBase) OnAsyncCall(f AsyncCallFunc) {
	glc.AsyncCall = f
}

// GetGameReadyUserCount 準備好的人數
func (glc *LogicBase) GetGameReadyUserCount() int32 {
	return int32(len(glc.TableManager.ReadyUser()))
}

// TableNo 桌號
func (glc *LogicBase) TableNo() int32 {
	return glc.TableManager.TableNo
}

// IsGaming 是否遊戲中
func (glc *LogicBase) IsGaming() bool {
	return glc.isgaming
}

// IsAllUserReady 是否所有人都準備
func (glc *LogicBase) IsAllUserReady() bool {
	for _, u := range glc.TableManager.GetTableUser() {
		if !u.IsReady() {
			return false
		}
	}
	return true
}

// BoardcastMessage 廣播消息
func (glc *LogicBase) BoardcastMessage(m *msg.Message) {
	for _, u := range glc.TableManager.GetTableUser() {
		if !u.IsOffLine() {
			glc.SendMessgeFunc(u.UserID, m)
		}
	}
}

// NotifyFishMessage 捕魚公告
func (glc *LogicBase) NotifyFishMessage(m *msg.Message) {
	glc.SendFishNotifyFunc(m)
}

// SendMessgeToSeat 向座位上的用戶發送消息
func (glc *LogicBase) SendMessgeToSeat(seatno int32, m *msg.Message) {
	seat := glc.TableManager.GetSeatByNo(seatno)

	if seat == nil {
		return
	}

	ur := seat.SeatUser
	if ur != nil {
		glc.SendMessgeFunc(ur.UserID, m)
	}
}

// SendMessageToUserID 向指定用戶ID發送消息
func (glc *LogicBase) SendMessageToUserID(userid int32, m *msg.Message) {
	glc.SendMessgeFunc(userid, m)
}

// SendMessageToUser 向指定用戶發送消息
func (glc *LogicBase) SendMessageToUser(user *GameUser, m *msg.Message) {
	glc.SendMessgeFunc(user.UserID, m)
}

// SendMessageToUserIDs 向一組指定用戶ID發送消息
func (glc *LogicBase) SendMessageToUserIDs(userids []int32, m *msg.Message) {
	for _, uid := range userids {
		glc.SendMessageToUserID(uid, m)
	}
}

// SendMessageToUsers 向一組指定用戶發送消息
func (glc *LogicBase) SendMessageToUsers(users []*GameUser, m *msg.Message) {
	for _, ur := range users {
		glc.SendMessageToUser(ur, m)
	}
}

// CheckSeatEmpty 座位是否空閒
func (glc *LogicBase) CheckSeatEmpty(seatno int32) bool {
	seat := glc.TableManager.GetSeatByNo(seatno)

	if seat == nil {
		return false
	}

	return seat.SeatUser == nil
}

// TableUserCount 當前桌子上用戶數
func (glc *LogicBase) TableUserCount() int32 {
	return int32(len(glc.TableManager.GetTableUser()))
}

// GetUserByUserID 根據用戶ID獲取用戶信息
func (glc *LogicBase) GetUserByUserID(userid int32) *GameUser {
	urs := glc.TableManager.GetTableUser()

	for _, ur := range urs {
		if ur.UserID == userid {
			return ur
		}
	}

	return nil
}

// GetUserBySeatNo 根據座位號獲取用戶
func (glc *LogicBase) GetUserBySeatNo(seatno int32) *GameUser {
	urs := glc.TableManager.GetTableUser()

	for _, ur := range urs {
		if ur.SeatNo == seatno {
			return ur
		}
	}

	return nil
}

// UserEnter 用戶進入
func (glc *LogicBase) UserEnter(gameuser *GameUser, tableno int32, seatno int32) {
	st := glc.TableManager.GetSeatByNo(seatno)
	if st == nil {
		mlog.Warn("用戶%v坐下,%d,%d,但座位不存在", gameuser, tableno, seatno)
		return
	}
	if st.SeatUser != nil {
		mlog.Warn("用戶%v坐下,%d,%d,但座位上已經有用戶%v", gameuser, tableno, seatno, st.SeatUser)
		return
	}

	mlog.Debug("user enter %v", gameuser)

	// 通知用戶坐下
	gameuser.UserState.OnUserSit(gameuser, tableno, seatno, glc.TableManager.Room)

	// 通知用戶進入
	glc.BoardcastUserEnter(gameuser)

	// 檢查用戶準備超時
	glc.checkUserReadyTimer(gameuser)

	// 廣播旁觀用戶
	glc.boardcastLookOnUser()

	// gamelog 初始化用戶
	urs := make(map[int32]*gamelog.UserData)
	urs[gameuser.UserID] = new(gamelog.UserData)
	urs[gameuser.UserID].Score = gameuser.Score
	urs[gameuser.UserID].TradeID = gameuser.TradeID
	glc.GameLog.UpdateGameUsers(urs)
	if glc.IsGaming() && glc.GameLog.Context() != nil {
		glc.GameLog.Context().AddUser(gameuser.UserID, gameuser.TradeID)
		glc.GameLog.Context().UserScore(gameuser.UserID, gameuser.Score)
	}
}

// checkGameStart 檢查遊戲開始
func (glc *LogicBase) checkGameStart() {
	if !glc.isgaming {
		if glc.IsAllUserReady() && glc.TableUserCount() >= glc.GetServerConfig().MinPlayer {
			glc.GameStartFunc()
		}
	}

}

// UserReady 用戶舉手
func (glc *LogicBase) UserReady(gameuser *GameUser) {

	mlog.Debug("user ready %v", gameuser)

	// 通知用戶舉手
	gameuser.UserState.OnUserReady(gameuser, glc.TableManager.Room)

	// 清除準備超時檢查
	glc.clearUserReadyTimer(gameuser)

	// 廣播用戶狀態改變
	glc.BoardcastUserStateChange(gameuser)

	// 檢查遊戲開始
	glc.checkGameStart()

	glc.checkMatchmaking(gameuser)
}

// UserFree 用戶站立
func (glc *LogicBase) UserFree(gameuser *GameUser) {
	mlog.Debug("user free %v", gameuser)

	// 清除用戶準備超時
	glc.clearUserReadyTimer(gameuser)

	// 廣播用戶離開
	glc.BoardcastUserLeave(gameuser)

	// 通知用戶站立
	gameuser.UserState.OnUserFree(gameuser, glc.TableManager.Room)

	// 檢查遊戲開始
	glc.checkGameStart()

	glc.checkMatchmaking(gameuser)

}

// UserLeave 用戶離開
func (glc *LogicBase) UserLeave(gameuser *GameUser) bool {
	mlog.Debug("user leave %v", gameuser)

	glc.UserFree(gameuser)

	// 通知用戶離開
	gameuser.UserState.OnUserLeave(gameuser, glc.TableManager.Room)

	// 檢查遊戲開始
	glc.checkGameStart()

	//glc.checkMatchmaking(gameuser)

	return true
}

// UserOffLine 用戶斷線
func (glc *LogicBase) UserOffLine(gameuser *GameUser) {
	// 如果用戶已經是斷線狀態
	if gameuser.IsOffLine() {
		return
	}

	mlog.Debug("user offline %v", gameuser)

	// 牌局記錄
	glc.GameLog.AddGameStep("用戶%d斷線", gameuser.UserID)

	// 通知用戶斷線
	gameuser.UserState.OnUserOffLine(gameuser, glc.TableManager.Room)

	glc.BoardcastUserStateChange(gameuser)
}

// UserRecover 用戶斷線重回
func (glc *LogicBase) UserRecover(gameuser *GameUser) {
	mlog.Debug("user recover %v", gameuser)

	// 牌局記錄
	glc.GameLog.AddGameStep("用戶%d斷線重連", gameuser.UserID)

	glc.notifyUserRecover(gameuser)

	// 通知用戶斷線重連
	gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)

	// 發送自己進入桌子消息
	glc.BoardcastUserSelfEnter(gameuser)

	// 發送其他用戶進入桌子消息
	glc.BoardcastOhterUserOnTable(gameuser)

	// 廣播用戶狀態改變
	glc.BoardcastUserStateChange(gameuser)

	// 廣播旁觀用戶
	glc.boardcastLookOnUser()
}

// GetFreeSeat 獲取一個空閒的座位
func (glc *LogicBase) GetFreeSeat() (int32, int32) {

	st := glc.TableManager.GetFreeSeat()

	if st != nil {
		return glc.TableNo(), st.SeatNo
	}

	return FreeTableNo, FreeSeatNo
}

// 通知用戶斷線重連
func (glc *LogicBase) notifyUserRecover(user *GameUser) {
	userRecoverMsg := new(msg.Message)
	userRecoverMsg.BClassID = int32(netproto.MessageBClassID_Game)
	userRecoverMsg.SClassID = int32(netproto.GameMessageClassID_UserRecoverID)

	glc.SendMessageToUser(user, userRecoverMsg)
}

// BoardcastUserEnter 通知用戶進入
func (glc *LogicBase) BoardcastUserEnter(user *GameUser) {
	// 向桌子上的所有玩家廣播用戶進入
	uenterMsg := new(msg.Message)
	uenterMsg.BClassID = int32(netproto.MessageBClassID_Game)
	uenterMsg.SClassID = int32(netproto.GameMessageClassID_UserEnterID)

	uenterdata := new(netproto.UserEnter)
	uenterdata.User = getProtobufGameUser(user)
	uenterMsg.MsgData = uenterdata

	glc.BoardcastMessage(uenterMsg)
	mlog.Debug("廣播用戶%v進入桌子消息", user)

	// 向用戶自己發送其他用戶進入消息
	glc.BoardcastOhterUserOnTable(user)
}

// BoardcastUserSelfEnter 廣播用戶自己進入
func (glc *LogicBase) BoardcastUserSelfEnter(user *GameUser) {
	tuenterMsg := new(msg.Message)
	tuenterMsg.BClassID = int32(netproto.MessageBClassID_Game)
	tuenterMsg.SClassID = int32(netproto.GameMessageClassID_UserEnterID)

	tuenterdata := new(netproto.UserEnter)
	tuenterdata.User = getProtobufGameUser(user)
	tuenterMsg.MsgData = tuenterdata
	glc.SendMessageToUser(user, tuenterMsg)

	mlog.Debug("向用戶%v發送自己進入桌子消息", user)
}

// BoardcastOhterUserOnTable 向用戶廣播桌子上的其他用戶
func (glc *LogicBase) BoardcastOhterUserOnTable(user *GameUser) {
	// 向用戶廣播所有桌子上的玩家
	for _, tur := range glc.AllTableUser() {
		tuenterMsg := new(msg.Message)
		tuenterMsg.BClassID = int32(netproto.MessageBClassID_Game)
		tuenterMsg.SClassID = int32(netproto.GameMessageClassID_UserEnterID)

		tuenterdata := new(netproto.UserEnter)
		tuenterdata.User = getProtobufGameUser(tur)
		tuenterMsg.MsgData = tuenterdata

		if user.UserID == tur.UserID {
			continue
		}

		mlog.Debug("向用戶%v發送%v進入桌子消息", user, tur)
		glc.SendMessageToUser(user, tuenterMsg)
	}
}

// BoardcastOther 向用戶廣播桌子上的其他用戶
func (glc *LogicBase) BoardcastOther(user *GameUser, m *msg.Message) {
	// 向用戶廣播所有桌子上的玩家
	for _, tur := range glc.AllTableUser() {
		if user != nil && user.UserID == tur.UserID {
			continue
		}

		glc.SendMessageToUser(tur, m)
	}
}

// BoardcastUserStateChange 廣播用戶狀態改變
func (glc *LogicBase) BoardcastUserStateChange(user *GameUser) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_UserStateChangeID)

	ust := new(netproto.UserStateChange)
	ust.User = getProtobufGameUser(user)

	m.MsgData = ust

	glc.BoardcastMessage(m)
}

// BoardcastUserLeave 廣播用戶離開
func (glc *LogicBase) BoardcastUserLeave(user *GameUser) {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_UserLeaveID)

	ulv := new(netproto.UserLeave)
	ulv.User = getProtobufGameUser(user)
	m.MsgData = ulv

	glc.BoardcastMessage(m)
}

// BoardcastGameStart 廣播遊戲開始
func (glc *LogicBase) BoardcastGameStart() {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_GameStartID)

	glc.BoardcastMessage(m)
}

// BoardcastGameEnd 廣播遊戲結束
func (glc *LogicBase) BoardcastGameEnd() {
	m := new(msg.Message)
	m.BClassID = int32(netproto.MessageBClassID_Game)
	m.SClassID = int32(netproto.GameMessageClassID_GameEndID)

	glc.BoardcastMessage(m)
}

func getProtobufGameUser(user *GameUser) *netproto.GameUser {
	return user.GetProtoBufGameUser()

}

// GetServerConfig 獲取服務器配置
func (glc *LogicBase) GetServerConfig() *GameServerConfig {
	return GetGameServerConfig(glc.Configs.GetConfig(ConfigFileNameServerInfo), glc.Configs.GetConfig(ConfigFileNameGameServer))
}

// 檢查用戶準備超時時間
func (glc *LogicBase) checkUserReadyTimer(user *GameUser) {
	readyWaitTime := glc.GetServerConfig().ReadyWaitTime + 5 // 增加5秒,確保客戶端倒計時結束後主動退出
	timername := glc.getTimerName(user)

	// 檢查timer是否已經啟動
	extimer := glc.TimerFuncManager.GetTimerFunc(timername)
	if extimer != nil && !extimer.IsStop() {
		extimer.Stop()
	}

	glc.TimerFuncManager.NewNamedTimer(readyWaitTime, func() {
		if user.IsSit() && user.TableNo == glc.TableNo() {
			mlog.Debug("用戶%v超時未準備,被踢出", user)
			glc.KickUserFunc(user.UserID, "由于您超时未准备,请尝试重新进入房间.")
		}
	}, timername)
}

// 清除用戶準備超時timer
func (glc *LogicBase) clearUserReadyTimer(user *GameUser) {
	timername := glc.getTimerName(user)
	mlog.Debug("移除準備超時 timer %s", timername)
	glc.TimerFuncManager.RemoveTimer(timername)
}

// CheckUserReadyTimer 檢查用戶準備超時時間
func (glc *LogicBase) CheckUserReadyTimer(user *GameUser) {
	glc.checkUserReadyTimer(user)
}

// ClearUserReadyTimer 清除用戶準備超時timer
func (glc *LogicBase) ClearUserReadyTimer(user *GameUser) {
	glc.clearUserReadyTimer(user)
}

// 獲取timer名稱
func (glc *LogicBase) getTimerName(user *GameUser) string {
	return fmt.Sprintf("gamelogic_readywait_timer_%d_%d_%d_%d", glc.TableNo(), glc.gameCounter, user.UserID, user.SeatNo)
}

// RegisterAfterGameStartFunc 註冊遊戲開始後處理函數
func (glc *LogicBase) RegisterAfterGameStartFunc(f func()) {
	glc.AfterGameStartFunc = append(glc.AfterGameStartFunc, f)
}

// RegisterAfterPiSaiStartFunc 註冊遊戲開始後處理函數
func (glc *LogicBase) RegisterAfterPiSaiStartFunc(f func()) {
	glc.AfterPiSaiStartFunc = append(glc.AfterPiSaiStartFunc, f)
}

// RegisterBeforeGameStartFunc 註冊遊戲開始前處理函數
func (glc *LogicBase) RegisterBeforeGameStartFunc(f func()) {
	glc.BeforeGameStartFunc = append(glc.BeforeGameStartFunc, f)
}

// RegisterBeforePiSaiStartFunc 註冊遊戲開始前處理函數
func (glc *LogicBase) RegisterBeforePiSaiStartFunc(f func()) {
	glc.BeforePiSaiStartFunc = append(glc.BeforePiSaiStartFunc, f)
}

// RegisterAfterGameEndFunc 註冊遊戲結束後處理函數
func (glc *LogicBase) RegisterAfterGameEndFunc(f func()) {
	glc.AfterGameEndFunc = append(glc.AfterGameEndFunc, f)
}

// RegisterAfterPiSaiEndFunc 註冊遊戲結束後處理函數
func (glc *LogicBase) RegisterAfterPiSaiEndFunc(f func()) {
	glc.AfterPiSaiEndFunc = append(glc.AfterPiSaiEndFunc, f)
}

// RegisterBeforeGameEndFunc 註冊遊戲結束前處理函數
func (glc *LogicBase) RegisterBeforeGameEndFunc(f func()) {
	glc.BeforeGameEndFunc = append(glc.BeforeGameEndFunc, f)
}

// RegisterBeforePiSaiEndFunc 註冊遊戲結束前處理函數
func (glc *LogicBase) RegisterBeforePiSaiEndFunc(f func()) {
	glc.BeforePiSaiEndFunc = append(glc.BeforePiSaiEndFunc, f)
}

// RegisterAnimation 註冊動畫播放監聽
func (glc *LogicBase) RegisterAnimation(actname string, timeout int32, callback func()) string {
	gameusercount := int32(glc.GetGameStateRealUserCount())
	return glc.RegisterAnimationByUserCount(actname, gameusercount, timeout, callback)
}

// RegisterAnimationByUserCount 註冊用戶動畫播放監聽
func (glc *LogicBase) RegisterAnimationByUserCount(actname string, usercount int32, timeout int32, callback func()) string {
	animationCallback := func() {
		// 廣播通知所有用戶動畫都已經播放完成
		glc.BoardcastAllUserAnimationCommited(actname)
		if callback != nil {
			callback()
		}
	}
	return glc.PlayerActionManager.NewAction(actname, usercount, timeout, animationCallback)
}

// CheckUserAnimationCommited 響應用戶提交動畫完成
func (glc *LogicBase) CheckUserAnimationCommited(user *GameUser, m *msg.Message, currentactname string, f func(string)) {
	if m.BClassID == int32(netproto.MessageBClassID_Game) && m.SClassID == int32(netproto.GameMessageClassID_UserActionCommitID) {
		data := m.MsgData.([]byte)

		ust := &netproto.UserActionCommited{}

		err := proto.Unmarshal(data, ust)
		if err != nil {
			mlog.Error("解析用戶%v動畫播放完成數據結構錯誤%v", user, err)
			glc.KickUserFunc(user.UserID, "非法操作")
			return
		}
		actname := *ust.ActionName
		mlog.Debug("收到用戶%d提交動畫%s播放完成.", user.UserID, actname)

		if f != nil {
			f(actname)
		}

		if strings.Index(currentactname, actname) == 0 {
			glc.PlayerActionManager.OnUserAction(currentactname, user.UserID)
		}
	}
}

// BoardcastAllUserAnimationCommited 廣播所有用戶動畫播放完成
func (glc *LogicBase) BoardcastAllUserAnimationCommited(actname string) {
	amsg := new(msg.Message)
	amsg.BClassID = int32(netproto.MessageBClassID_Game)
	amsg.SClassID = int32(netproto.GameMessageClassID_AllUserActionCommitID)

	acdata := new(netproto.AllUserActionCommited)
	acdata.ActionName = proto.String(actname)

	amsg.MsgData = acdata

	glc.BoardcastMessage(amsg)
}

// GetGameStateUserCount 當前遊戲狀態中的用戶數，排除離線用戶
func (glc *LogicBase) GetGameStateUserCount() int {
	return len(glc.AllOnLineGameUser())
}

// GetGameStateRealUserCount 當前遊戲狀態中的真實用戶數，排除機器人,離線用戶
func (glc *LogicBase) GetGameStateRealUserCount() int {
	urs := glc.AllOnLineGameUser()

	cn := 0
	for _, ur := range urs {
		if !ur.IsRobot() {
			cn++
		}
	}

	return cn
}

// GetCurrentGameUserCount 獲取參與當前牌局遊戲的用戶數,包括離線用戶
func (glc *LogicBase) GetCurrentGameUserCount() int {
	return len(glc.AllGameUser())
}

// AllGameUser 獲取所有遊戲用戶
func (glc *LogicBase) AllGameUser() []*GameUser {
	urs := make([]*GameUser, 0)
	for _, ur := range glc.AllTableUser() {
		if ur.IsGameing() || ur.IsOffLine() {
			urs = append(urs, ur)
		}
	}

	return urs
}

// AllLookOnUser 獲取所有旁觀用戶
func (glc *LogicBase) AllLookOnUser() []*GameUser {
	urs := make([]*GameUser, 0)
	for _, ur := range glc.AllTableUser() {
		if !ur.IsGameing() && !ur.IsOffLine() {
			urs = append(urs, ur)
		}
	}

	return urs
}

// AllOnLineGameUser 獲取在線遊戲的用戶
func (glc *LogicBase) AllOnLineGameUser() []*GameUser {
	urs := make([]*GameUser, 0)
	for _, ur := range glc.AllTableUser() {
		if ur.IsGameing() {
			urs = append(urs, ur)
		}
	}

	return urs
}

// AllTableUser 桌子上的所有玩家
func (glc *LogicBase) AllTableUser() []*GameUser {
	return glc.TableManager.GetTableUser()
}

// AllOnLineTableUser 桌子的所有在線玩家
func (glc *LogicBase) AllOnLineTableUser() []*GameUser {
	urs := glc.AllTableUser()

	onurs := make([]*GameUser, 0)
	for _, ur := range urs {
		if ur.IsOffLine() {
			continue
		}

		onurs = append(onurs, ur)
	}

	return onurs
}

// IsGameUser 是否遊戲中的用戶
func (glc *LogicBase) IsGameUser(user *GameUser) bool {
	return user.IsGameing() || user.IsOffLine()
}

// BoardcastGameUserMessage 向桌子上的遊戲玩家廣播消息
func (glc *LogicBase) BoardcastGameUserMessage(m *msg.Message) {
	for _, u := range glc.AllOnLineGameUser() {
		if !u.IsOffLine() {
			glc.SendMessgeFunc(u.UserID, m)
		}
	}
}

// CheckGameUserMessage 檢查用戶消息
func (glc *LogicBase) CheckGameUserMessage(user *GameUser) bool {
	if !glc.IsGameUser(user) {
		glc.KickUserFunc(user.UserID, "非法操作")
		return false
	}

	return true
}

// GetUserWinAndTax 獲取用戶扣除稅收後的輸贏和稅收金額
func (glc *LogicBase) GetUserWinAndTax(amount int64) (int64, int64) {
	if amount <= 0 {
		return amount, 0
	}

	taxscale := glc.GetServerConfig().TaxScale

	taxamount := int64(amount * int64(taxscale) / 100)
	return amount - taxamount, taxamount
}

// GetUserTaxAddJackpotPool 將玩家稅收抽一部分放入獎金池
func (glc *LogicBase) GetUserTaxAddJackpotPool(taxamount int64, isrobot bool) (int64, int64) {
	if taxamount <= 0 || isrobot {
		return 0, 0
	}

	taxscale := glc.GetServerConfig().TaxAddPool
	addjackpottax := int64(taxamount * int64(taxscale) / 100)
	newtaxamount := taxamount - addjackpottax

	return newtaxamount, addjackpottax
}

// GetTax 取得稅額
func (glc *LogicBase) GetTax() float64 {
	return (float64(100) - float64(glc.GetServerConfig().TaxScale)) / float64(100)
}

// GetZhushuimoneyTax 注水的錢要包含計算收稅比例後無條件捨去
func (glc *LogicBase) GetZhushuimoneyTax(zhushuimoney int64) int64 {
	taxscale := glc.GetServerConfig().TaxScale
	t1 := float32(1) - float32(taxscale)/float32(100)
	t2 := math.Ceil(float64(float32(zhushuimoney) * t1))
	return int64(t2)
}

// GetTable 獲取桌子
func (glc *LogicBase) GetTable() *Table {
	return glc.TableManager
}

// 廣播遊戲旁觀用戶
func (glc *LogicBase) boardcastLookOnUser() {
	lookonMsg := new(msg.Message)
	lookonMsg.BClassID = int32(netproto.MessageBClassID_Game)
	lookonMsg.SClassID = int32(netproto.GameMessageClassID_LookOnUserNotifyID)

	lookOnData := new(netproto.LookOnUserNotify)
	lookonMsg.MsgData = lookOnData

	urs := glc.AllTableUser()
	for _, ur := range urs {
		if glc.isgaming && !ur.IsGameing() && !ur.IsOffLine() {
			lookOnData.UserID = append(lookOnData.UserID, ur.UserID)
		}
	}

	mlog.Debug("發送旁觀用戶,%v", lookOnData)
	glc.BoardcastMessage(lookonMsg)

}

// SetGameLog 設置牌局記錄器
func (glc *LogicBase) SetGameLog(glog gamelog.GameLog) {
	glc.GameLog = NewGameLogHelper(glog, glc)
}

// GetSuperToolService 取得超端服務
func (glc *LogicBase) GetSuperToolService() ISuperToolService {
	return glc.GetTable().Room.SuperToolService
}

//func (glc *GameLogicBase) SetVipUserProvider(provider *VipUserProvider) {
//	//glc.VipUserProvider = provider
//}
//
//func (glc *GameLogicBase) SetVipUserService(service *VipUserService) {
//	//glc.VipUserService = service
//}

// ExchangeMaxCard 交换最大牌型
func (glc *LogicBase) ExchangeMaxCard(userid int32) {
	mlog.Debug("用戶%d交换最大牌型", userid)
	//牌局记录
	glc.GameLog.AddGameStep("用戶%d交换最大牌型", userid)
}

// ReadProtoConfig 讀取proto配置文件
func (glc *LogicBase) ReadProtoConfig(fName string, unmarshFunc func(data []byte) (proto.Message, error)) (proto.Message, error) {
	protoConfigFolderName := "protoconfig"

	if conf.TestMode {
		protoConfigFolderName = "./protoconfig"
	}

	currentPath, _ := os.Getwd()
	fpath := path.Join(currentPath, protoConfigFolderName, fName)

	fi, err := os.Open(fpath)

	if err != nil {
		return nil, err
	}
	defer fi.Close()

	fd, err := ioutil.ReadAll(fi)
	if err != nil {
		return nil, err
	}

	return unmarshFunc(fd)
}

// RefreshConfigGroup 更新設定檔
func (glc *LogicBase) RefreshConfigGroup(congifs *config.Group) {
	glc.Configs = congifs
	//for _, c := range glc.Configs.Configs {
	//	for _, s := range c.Sections {
	//		for k, v := range s {
	//			//mlog.Debug("reload k[%v] v[%v]", k, v)
	//		}
	//	}
	//}
}

// TableRealUserCount 桌子上的用戶數
func (glc *LogicBase) TableRealUserCount() int {
	return glc.TableManager.TableRealUserCount()
}

func (glc *LogicBase) checkMatchmaking(gameuser *GameUser) {
	if glc.GetServerConfig().MatchMakingNumber == 0 {
		return
	}
	var gameStart = func() {
		if glc.IsAllUserReady() && glc.TableUserCount() >= glc.GetServerConfig().MatchMakingNumber {
			glc.GameStartFunc()
		}
	}

	// X秒配對，到達Y人開始倒數，低於Y人重設計時器
	if !glc.isgaming {
		timername := fmt.Sprintf("CountDown_%d", glc.TableNo())
		if glc.IsAllUserReady() && glc.TableUserCount() >= glc.GetServerConfig().MatchMakingNumber {
			// 檢查是否有正在倒數的計時器
			if tf := glc.TimerFuncManager.GetTimerFunc(timername); tf == nil {
				// 開始倒數時間，倒數完成開始遊戲
				// mlog.System("開始倒數")
				glc.TimerFuncManager.NewNamedTimer(glc.GetServerConfig().MatchMakingTime, gameStart, timername)
				// glc.TimerFuncManager.NewNamedTimer(glc.GetServerConfig().MatchMakingTime, glc.GameStartFunc, timername)
				pmsg := new(msg.Message)
				pmsg.BClassID = int32(netproto.MessageBClassID_Game)
				pmsg.SClassID = int32(netproto.GameMessageClassID_GameBeginCountDownID)
				pdata := new(netproto.GameBeginCountDown)
				pdata.Time = proto.Int32(glc.GetServerConfig().MatchMakingTime)
				pmsg.MsgData = pdata

				glc.BoardcastMessage(pmsg)
			} else {
				time := glc.GetServerConfig().MatchMakingTime - tf.ElapsedSeconds()
				// mlog.System("玩家進入")
				pmsg := new(msg.Message)
				pmsg.BClassID = int32(netproto.MessageBClassID_Game)
				pmsg.SClassID = int32(netproto.GameMessageClassID_GameBeginCountDownID)
				pdata := new(netproto.GameBeginCountDown)
				pdata.Time = proto.Int32(time)
				pmsg.MsgData = pdata

				glc.SendMessageToUser(gameuser, pmsg)
			}

		} else {
			if tf := glc.TimerFuncManager.GetTimerFunc(timername); tf != nil {
				// 低於兩個人 暫停計時器
				// mlog.System("計時器關閉")
				glc.TimerFuncManager.RemoveTimer(timername)
				pmsg := new(msg.Message)
				pmsg.BClassID = int32(netproto.MessageBClassID_Game)
				pmsg.SClassID = int32(netproto.GameMessageClassID_GameBeginCountDownStopID)
				glc.BoardcastMessage(pmsg)
			}
		}
	}
}
