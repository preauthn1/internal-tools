package gamelogic

// UserFreeState 玩家站立狀態
type UserFreeState struct{}

// OnUserFree 玩家站立
func (ust *UserFreeState) OnUserFree(user *GameUser, room *Room) bool {
	return false
}

// OnUserSit 玩家坐下
func (ust *UserFreeState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
	// 只有機器人可以由站立狀態坐下
	if !ur.IsRobot() {
		return false
	}

	t := room.GetTableByNo(tableno)
	if t == nil {
		return false
	}

	st := t.GetSeatByNo(seatno)
	if st == nil {
		return false
	}

	if st.SeatUser != nil {
		return false
	}

	st.SeatUser = ur

	ur.TableNo = tableno
	ur.SeatNo = seatno

	ur.UserState = sitState

	return true
}

// OnUserQueue 玩家排隊
func (ust *UserFreeState) OnUserQueue(user *GameUser, room *Room) bool {
	user.UserState = queueState
	return true
}

// OnUserReady 玩家準備
func (ust *UserFreeState) OnUserReady(user *GameUser, room *Room) bool {
	return false
}

// OnUserGame 玩家遊戲中
func (ust *UserFreeState) OnUserGame(user *GameUser, room *Room) bool {
	return false
}

// OnGameEnd 玩家遊戲結束
func (ust *UserFreeState) OnGameEnd(user *GameUser, room *Room) bool {
	return false
}

// OnUserLeave 玩家離開房間
func (ust *UserFreeState) OnUserLeave(user *GameUser, room *Room) bool {
	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserFreeState) OnUserOffLine(user *GameUser, room *Room) bool {
	return false
}

// StateName 狀態名稱
func (ust *UserFreeState) StateName() string {
	return StateNameFree
}

// StateValue 狀態值
func (ust *UserFreeState) StateValue() int32 {
	return UserStateFree
}
