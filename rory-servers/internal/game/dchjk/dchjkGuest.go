package dchjk

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

type DchjkGuest struct {
	User         *gamelogic.GameUser
	i32UserID    int32 // 玩家ID
	i64EnterTime int64 //入桌時間
	bKickTip     bool  //踢人倒數提醒
}

func NewDchjkGuest() *DchjkGuest {
	pGuest := new(DchjkGuest)

	pGuest.Init()

	return pGuest
}

// 初始化旁觀者資料
func (pUser *DchjkGuest) Init() {
	pUser.User = nil
	pUser.i32UserID = -1
	pUser.i64EnterTime = time.Now().UnixNano()
	pUser.bKickTip = false
}

/*-------------Guest Manager--------------*/

func (pGame *DchjkLogic) GetFreeGuestSeat() int {
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGuest := pGame.ary_pGuestData[iSeat]
		if pGuest.User == nil {
			return iSeat
		}
	}

	return -1
}

func (pGame *DchjkLogic) GetGuestSeat(pGameUser *gamelogic.GameUser) int {
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGuest := pGame.ary_pGuestData[iSeat]
		if pGuest.i32UserID == pGameUser.UserID {
			return iSeat
		}
	}

	return gamelogic.FreeSeatNo
}

func (pGame *DchjkLogic) GetGuestCount() int {
	count := 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGuest := pGame.ary_pGuestData[iSeat]
		if pGuest.User != nil {
			count++
		}
	}

	return count
}

func (pGame *DchjkLogic) GetRealGuestCount() int {
	count := 0
	for iSeat := 0; iSeat < pGame.iSeatCount; iSeat++ {
		pGuest := pGame.ary_pGuestData[iSeat]
		if pGuest.User != nil && !pGuest.User.IsRobot() {
			count++
		}
	}

	return count
}
