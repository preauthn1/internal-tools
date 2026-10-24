package hljbp

import (
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// User 玩家資料
type User struct {
	gameUser     *gamelogic.GameUser
	money        int64 // 玩家金錢
	state        int   // 玩家狀態
	noBetRound   int   // 無壓注回合
	totalBet     int64
	moneyChange  int64 // 金錢變化，最後寫入輸贏流水用
	winLoseMoney int64 // 玩家輸贏金錢(不扣稅收)，若值為負數，則為玩家輸掉的金額，client 端顯示用
	totalWinLose int64 // 總輸贏(已扣稅收)，client 端顯示用
	winState     int   // 玩家輸贏
	betAmount    int64
	winAmount    int64
	flowAmount   int64
	taxAmount 	 int64
}

// NewUser 創建玩家資料
func NewUser() *User {
	user := new(User)
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
	user.totalBet = 0
	user.winLoseMoney = 0
	user.moneyChange = 0
	user.totalWinLose = 0
	user.winState = WinStateNone
	user.betAmount = 0
	user.winAmount = 0
	user.flowAmount = 0
	user.taxAmount = 0
}

/*-------------User Manager--------------*/

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
