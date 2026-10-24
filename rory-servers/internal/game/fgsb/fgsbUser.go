package fgsb

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// User 玩家資料
type User struct {
	gameUser         *gamelogic.GameUser
	cards            *Poker // 用户手牌
	money            int64  // 玩家金錢
	state            int    // 玩家狀態
	noBetRound       int    // 無壓注回合
	betBase          int64  // 玩家底注
	betSide          int64  // 玩家邊注
	betDouble        int64  // 玩家加注
	winSide          int64  // 邊注輸贏金錢，若值為負數，則為玩家輸掉的金額
	winSpecial       int64  // 特殊牌型輸贏金錢，若值為負數，則為玩家輸掉的金額
	winBase          int64  // 底注輸贏金錢，若值為負數，則為玩家輸掉的金額
	winDouble        int64  // 加注輸贏金錢，若值為負數，則為玩家輸掉的金額
	moneyChange      int64  // 金錢變化
	totalWinLose     int64  // 總輸贏，client 端顯示用
	tempSideMoney    int64  // 暫存邊注贏多少，只有正數
	tempSpecialMoney int64  // 暫存特殊牌型贏多少，只有正數
	tempTax          int64  // 暫存稅收
	tempFlow         int64  // 暫存流水
	winState         int    // 玩家輸贏
	totalBet         int64  // 總押注額
}

// NewUser 創建玩家資料
func NewUser() *User {
	user := new(User)
	user.cards = NewPoker()
	user.Init()

	return user
}

// Init 初始化玩家資料
func (user *User) Init() {
	user.gameUser = nil
	user.money = 0
	user.state = UserStateEmptySeat
	user.noBetRound = 0

	user.ResetGameData()
}

// ResetGameData 重置遊戲資料
func (user *User) ResetGameData() {
	user.cards.CardGroup = poker.NewEmptyCardGroup(0)
	user.betBase = 0
	user.betSide = 0
	user.betDouble = 0
	user.winSide = 0
	user.winSpecial = 0
	user.winBase = 0
	user.winDouble = 0
	user.moneyChange = 0
	user.totalWinLose = 0
	user.tempSideMoney = 0
	user.tempSpecialMoney = 0
	user.tempTax = 0
	user.winState = WinStateNone
	user.totalBet = 0
}

/*-------------User Manager--------------*/

// CheckFreeUserSeat check user是否為空
func (glc *Logic) CheckFreeUserSeat(seat int) bool {
	guest := glc.guest[seat]
	if guest.user == nil {
		return true
	}

	return false
}

// GetUserData 取得玩家資料byID
func (glc *Logic) GetUserData(userID int32) *User {
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		user := glc.user[i]
		if user.state == UserStateEmptySeat {
			continue
		}

		if user.gameUser == nil {
			continue
		}

		if user.gameUser.UserID == userID {
			return user
		}
	}

	return nil
}

// GetUserSeat 取得玩家座位byID
func (glc *Logic) GetUserSeat(userID int32) int32 {
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		user := glc.user[i]
		if user.gameUser != nil && user.gameUser.UserID == userID {
			return i
		}
	}

	return gamelogic.FreeSeatNo
}

// GetNormalUserCount 取得正常遊戲人數
func (glc *Logic) GetNormalUserCount() int {
	count := 0
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		user := glc.user[i]
		if user.state == UserStateNormal && user.gameUser != nil {
			count++
		}
	}

	return count
}

// GetUserCount 取得玩家人數
func (glc *Logic) GetUserCount() int {
	count := 0
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		user := glc.user[i]
		if user.gameUser != nil {
			count++
		}
	}

	return count
}

// GetRealUserCount 取得玩家真人人數
func (glc *Logic) GetRealUserCount() int {
	count := 0
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		user := glc.user[i]
		if user.gameUser != nil && !user.gameUser.IsRobot() {
			count++
		}
	}

	return count
}
