package gamelogic

// UserSitState 玩家坐下狀態
type UserSitState struct{}

// OnUserFree 玩家站立
func (ust *UserSitState) OnUserFree(user *GameUser, room *Room) bool {
	tno := user.TableNo
	sno := user.SeatNo

	t := room.GetTableByNo(tno)
	if t == nil {
		return false
	}

	st := t.GetSeatByNo(sno)
	if st == nil {
		return false
	}

	st.SeatUser = nil
	user.TableNo = FreeTableNo
	user.SeatNo = FreeSeatNo

	user.LastTableNo = tno
	user.LastSeatNo = sno
	user.GameCountOnTheTable = 0

	user.UserState = freeState

	return true
}

// OnUserSit 玩家坐下
func (ust *UserSitState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
	return false
}

// OnUserQueue 玩家排隊
func (ust *UserSitState) OnUserQueue(user *GameUser, room *Room) bool {
	ust.OnUserFree(user, room)

	return user.UserState.OnUserQueue(user, room)
}

// OnUserReady 玩家準備
func (ust *UserSitState) OnUserReady(user *GameUser, room *Room) bool {
	user.UserState = readyState
	return true
}

// OnUserGame 玩家遊戲中
func (ust *UserSitState) OnUserGame(user *GameUser, room *Room) bool {
	return false
}

// OnGameEnd 玩家遊戲結束
func (ust *UserSitState) OnGameEnd(user *GameUser, room *Room) bool {
	return false
}

// OnUserLeave 玩家離開房間
func (ust *UserSitState) OnUserLeave(user *GameUser, room *Room) bool {
	ust.OnUserFree(user, room)
	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserSitState) OnUserOffLine(user *GameUser, room *Room) bool {
	return false
}

// StateName 狀態名稱
func (ust *UserSitState) StateName() string {
	return StateNameSit
}

// StateValue 狀態值
func (ust *UserSitState) StateValue() int32 {
	return UserStateSit
}
