package gamelogic

import (
	"fmt"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/player"
)

// 性别
const (
	UserSexGirl = iota // 性别，女
	UserSexBoy
)

// 用戶類型
const (
	UserTypeNormal = 1 // 用戶類型，普通用戶
	UserTypeRebot  = 2 // 機器人
	UserTypeGuest  = 3 // 遊客
)

// ?
const (
	FreeTableNo = -1
	FreeSeatNo  = -1
)

// GameUser 用户信息
type GameUser struct {
	TableNo             int32             // 桌子號
	SeatNo              int32             // 椅子號
	UserID              int32             // 用戶ID
	Score               int64             // 當前積分 (當局積分)
	InGamingMoneyChange int64             // 遊戲中金額異動 (用於斷線重連時充值金額到帳或遊戲中充值)
	LoginScore          int64             // 進入時的積分
	UpdateScore         int64             // 變化的積分
	TaxAmount           int64             // 稅收金額
	NickName            string            // 暱稱
	UserState           UserState         // 用戶狀態
	ConnInfo            *network.ConnInfo // 網絡連接信息
	WinCount            int32             // 贏
	LoseCount           int32             // 輸
	DrawCount           int32             // 平
	UpdateWinCount      int32             // 贏變化
	UpdateLoseCount     int32             // 輸變化
	UpdateDrawCount     int32             // 平變化
	HDCode              string            // 設備碼
	HDType              int32             // 設備類型
	UserSex             int32             // 性別
	UserType            int32             // 用戶類型
	FaceID              int32             // 頭像ID
	GameTime            int32             // 遊戲時間
	LoginTime           time.Time         // 進入時間
	LastTableNo         int32             // 上一次桌子號
	LastSeatNo          int32             // 上一次椅子號
	GameBuff            string            // 遊戲控制參數
	GameCountOnTheTable int32             // 當前桌上游戲次數
	TradeID             int64             // 每次進入遊戲房間會獲取一個唯一的值
	LastWinLose         int32             // 上一次輸贏結果: 0輸 1贏
	SumCharge           int64             // 充值金額
	SumConvert          int64             // 總兌換
	HasBindPay          bool              // 是否綁定支付寶或銀行卡之一
	ExperienceCash      int64             // 體驗金幣
	FlowAmount          int64             // 流水加總
	BetAmount           int64             // RTP-有效投注金額加總
	PayAmount           int64             // RTP-返還金額加總
	LoginName           string            // ?
	AlipayAccount       string            // ?
	ContestID           int32             // 比賽活動ID
	TimeSeq             int64             // 用戶登入遊戲時間戳
}

// NewGameUser 創建用戶對象
func NewGameUser() *GameUser {
	ur := new(GameUser)
	ur.UserState = freeState
	ur.TableNo = FreeTableNo
	ur.WinCount = 0
	ur.LoseCount = 0
	ur.DrawCount = 0
	ur.UpdateWinCount = 0
	ur.UpdateLoseCount = 0
	ur.UpdateDrawCount = 0
	ur.UpdateScore = 0
	ur.Score = 0
	ur.InGamingMoneyChange = 0
	ur.LoginScore = 0
	ur.TableNo = FreeTableNo
	ur.SeatNo = FreeSeatNo
	ur.TaxAmount = 0
	ur.GameTime = 0
	ur.LoginTime = time.Now()
	ur.LastTableNo = FreeTableNo
	ur.LastSeatNo = FreeSeatNo
	ur.GameBuff = ""
	ur.GameCountOnTheTable = 0
	ur.TradeID = 0
	ur.LastWinLose = 0
	ur.FlowAmount = 0
	ur.BetAmount = 0
	ur.PayAmount = 0
	ur.TimeSeq = 0

	return ur
}

// IsFree 是否空閒
func (ur *GameUser) IsFree() bool {
	return ur.UserState.StateName() == StateNameFree
}

// IsQueue 是否排隊
func (ur *GameUser) IsQueue() bool {
	return ur.UserState.StateName() == StateNameQueue
}

// IsSit 是否坐下
func (ur *GameUser) IsSit() bool {
	return ur.UserState.StateName() == StateNameSit
}

// IsReady 是否舉手
func (ur *GameUser) IsReady() bool {
	return ur.UserState.StateName() == StateNameReady
}

// IsGameing 是否遊戲中
func (ur *GameUser) IsGameing() bool {
	return ur.UserState.StateName() == StateNameGame
}

// IsOffLine 是否斷線
func (ur *GameUser) IsOffLine() bool {
	return ur == nil || ur.UserState.StateName() == StateNameOffline
}

// GetProtoBufGameUser 獲取網絡傳輸的gameuser 對象
func (ur *GameUser) GetProtoBufGameUser() *netproto.GameUser {
	nUser := new(netproto.GameUser)
	nUser.UserID = proto.Int32(ur.UserID)
	nUser.NickName = proto.String(ur.NickName)
	nUser.Sex = proto.Int32(ur.UserSex)
	nUser.FaceID = proto.Int32(ur.FaceID)
	nUser.Score = proto.Int64(ur.Score) // 給Client的分數是Score加Coupon
	nUser.WinCount = proto.Int32(ur.WinCount)
	nUser.LoseCount = proto.Int32(ur.LoseCount)
	nUser.DrawCount = proto.Int32(ur.DrawCount)
	nUser.TableNo = proto.Int32(ur.TableNo)
	nUser.SeatNo = proto.Int32(ur.SeatNo)
	nUser.StateValue = proto.Int32(ur.UserState.StateValue())
	nUser.LastWinLose = proto.Int32(ur.LastWinLose)

	return nUser
}

func (ur *GameUser) String() string {
	prefix := ""
	if ur.IsRobot() {
		prefix = "[Robot] "
	}
	return fmt.Sprintf(prefix+"userid:%d nickname:%s,score:%d,updatescore:%d,tax:%d,gametime:%d,logintime:%s,state:%s,tableno:%d,seatno:%d", ur.UserID, ur.NickName, ur.Score, ur.UpdateScore, ur.TaxAmount, ur.GameTime, ur.LoginTime, ur.UserState.StateName(), ur.TableNo, ur.SeatNo)
}

// IsRobot 是否機器人
func (ur *GameUser) IsRobot() bool {
	return ur == nil || ur.UserType == UserTypeRebot
}

// RoomUserManager 房間用戶管理
type RoomUserManager struct {
	*player.Manager
}

// NewRoomUserManager 創建房間用戶管理器
func NewRoomUserManager(playerManager *player.Manager) *RoomUserManager {
	rm := new(RoomUserManager)
	rm.Manager = playerManager

	return rm
}

// GetGameUserByConnID 根據連接ID獲取用戶信息
func (rm *RoomUserManager) GetGameUserByConnID(connid string) *GameUser {
	clt := rm.Manager.GetPlayerByConnID(connid)

	if clt == nil {
		return nil
	}

	if clt.IsGuest() {
		return nil
	}

	return rm.GetGameUserByUserID(clt.UserID)
}

// GetGameUserByUserID 根據用戶ID獲取用戶信息
func (rm *RoomUserManager) GetGameUserByUserID(userid int32) *GameUser {
	pr := rm.Manager.GetPlayerByUserID(userid)

	if pr == nil {
		return nil
	}

	ur := pr.UserData.(*GameUser)
	ur.ConnInfo = pr.Conn

	return ur
}

// GetGameUserByTableNo 根據桌子號獲取用戶
func (rm *RoomUserManager) GetGameUserByTableNo(tableno int32) []*GameUser {
	var urs []*GameUser

	for _, pr := range rm.Manager.LogonPlayer() {
		ur := pr.UserData.(*GameUser)

		if ur.TableNo == tableno {
			urs = append(urs, ur)
		}
	}

	return urs
}

// GetGameUserByTableAndSeat 根據桌子號椅子號獲取用戶
func (rm *RoomUserManager) GetGameUserByTableAndSeat(tableno int32, seatno int32) *GameUser {

	for _, pr := range rm.Manager.LogonPlayer() {
		ur := pr.UserData.(*GameUser)

		if ur.TableNo == tableno && ur.SeatNo == seatno {
			return ur
		}
	}

	return nil
}

// GetAllUser 獲取所有的房間用戶
func (rm *RoomUserManager) GetAllUser() []*GameUser {
	var urs []*GameUser

	for _, pr := range rm.Manager.LogonPlayer() {
		ur := pr.UserData.(*GameUser)

		urs = append(urs, ur)
	}

	return urs
}

// GetFreeUser 獲取房間所有空閒用戶
func (rm *RoomUserManager) GetFreeUser() []*GameUser {
	return rm.GetStateUser(UserStateFree)
}

// GetQueueUser 獲取房間所有排隊用戶
func (rm *RoomUserManager) GetQueueUser() []*GameUser {
	return rm.GetStateUser(UserStateQueue)
}

// GetGameingUser 獲取房間所有遊戲用戶
func (rm *RoomUserManager) GetGameingUser() []*GameUser {
	return rm.GetStateUser(UserStateGame)
}

// GetOffLineUser 獲取房間所有遊戲用戶
func (rm *RoomUserManager) GetOffLineUser() []*GameUser {
	return rm.GetStateUser(UserStateOffline)
}

// GetSitUser 獲取房間所有坐下用戶
func (rm *RoomUserManager) GetSitUser() []*GameUser {
	return rm.GetStateUser(UserStateSit)
}

// GetReadyUser 獲取房間所有準備用戶
func (rm *RoomUserManager) GetReadyUser() []*GameUser {
	return rm.GetStateUser(UserStateReady)
}

// GetStateUser 獲取房間指定狀態的用戶
func (rm *RoomUserManager) GetStateUser(stateValue int32) []*GameUser {
	var urs []*GameUser
	for _, pr := range rm.Manager.LogonPlayer() {
		ur := pr.UserData.(*GameUser)

		if ur.UserState.StateValue() == stateValue {
			urs = append(urs, ur)
		}
	}

	return urs
}

// UserCount 用戶數
func (rm *RoomUserManager) UserCount() int32 {
	return int32(rm.Manager.UserCount())
}
