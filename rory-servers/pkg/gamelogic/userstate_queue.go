package gamelogic

// UserQueueState 玩家排隊狀態
type UserQueueState struct{}

// OnUserFree 玩家站立
func (ust *UserQueueState) OnUserFree(user *GameUser, room *Room) bool {
	return false
}

// OnUserSit 玩家坐下
func (ust *UserQueueState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
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
func (ust *UserQueueState) OnUserQueue(user *GameUser, room *Room) bool {
	return false
}

// OnUserReady 玩家準備
func (ust *UserQueueState) OnUserReady(user *GameUser, room *Room) bool {
	return false
}

// OnUserGame 玩家遊戲中
func (ust *UserQueueState) OnUserGame(user *GameUser, room *Room) bool {
	return false
}

// OnGameEnd 玩家遊戲結束
func (ust *UserQueueState) OnGameEnd(user *GameUser, room *Room) bool {
	return false
}

// OnUserLeave 玩家離開房間
func (ust *UserQueueState) OnUserLeave(user *GameUser, room *Room) bool {
	user.UserState = freeState
	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserQueueState) OnUserOffLine(user *GameUser, room *Room) bool {
	return false
}

// StateName 狀態名稱
func (ust *UserQueueState) StateName() string {
	return StateNameQueue
}

// StateValue 狀態值
func (ust *UserQueueState) StateValue() int32 {
	return UserStateQueue
}
