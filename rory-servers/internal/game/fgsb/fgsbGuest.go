package fgsb

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// Guest 觀戰者資料
type Guest struct {
	user      *gamelogic.GameUser
	userID    int32 // 玩家ID
	enterTime int64 // 入桌時間
	kickTip   bool  // 踢人倒數提醒
}

// NewGuest 創建觀戰者資料
func NewGuest() *Guest {
	guest := new(Guest)

	guest.Init()

	return guest
}

// Init 初始化旁觀者資料
func (user *Guest) Init() {
	user.user = nil
	user.userID = -1
	user.enterTime = time.Now().UnixNano()
	user.kickTip = false
}

/*-------------Guest Manager--------------*/

// GetFreeGuestSeat 取得觀戰座位空位
func (glc *Logic) GetFreeGuestSeat() int32 {
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		guest := glc.guest[i]
		if guest.user == nil {
			return i
		}
	}

	return -1
}

// GetGuestData 取得觀戰者資料byID
func (glc *Logic) GetGuestData(userID int32) *Guest {
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		guest := glc.guest[i]
		if guest.userID == userID {
			return guest
		}
	}

	return nil
}

// GetGuestSeat 取得觀戰者座位byID
func (glc *Logic) GetGuestSeat(userID int32) int32 {
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		guest := glc.guest[i]
		if guest.userID == userID {
			return i
		}
	}

	return gamelogic.FreeSeatNo
}

// GetGuestCount 取得觀戰人數
func (glc *Logic) GetGuestCount() int {
	count := 0
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		guest := glc.guest[i]
		if guest.user != nil {
			count++
		}
	}

	return count
}

// GetRealGuestCount 取得觀戰真人人數
func (glc *Logic) GetRealGuestCount() int {
	count := 0
	for i := int32(0); i < glc.TableManager.Room.SeatCount; i++ {
		guest := glc.guest[i]
		if guest.user != nil && !guest.user.IsRobot() {
			count++
		}
	}

	return count
}
