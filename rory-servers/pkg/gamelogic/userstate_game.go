package gamelogic

// UserGameState 玩家遊戲狀態
type UserGameState struct{}

// OnUserFree 玩家站立
func (ust *UserGameState) OnUserFree(user *GameUser, room *Room) bool {
	return false
}

// OnUserSit 玩家坐下
func (ust *UserGameState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
	return false
}

// OnUserQueue 玩家排隊
func (ust *UserGameState) OnUserQueue(user *GameUser, room *Room) bool {
	return false
}

// OnUserReady 玩家準備
func (ust *UserGameState) OnUserReady(user *GameUser, room *Room) bool {
	return false
}

// OnUserGame 玩家遊戲中
func (ust *UserGameState) OnUserGame(user *GameUser, room *Room) bool {
	return false
}

// OnGameEnd 玩家遊戲結束
func (ust *UserGameState) OnGameEnd(user *GameUser, room *Room) bool {
	user.UserState = sitState
	return true
}

// OnUserLeave 玩家離開房間
func (ust *UserGameState) OnUserLeave(user *GameUser, room *Room) bool {
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
	user.UserState = freeState

	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserGameState) OnUserOffLine(user *GameUser, room *Room) bool {
	user.UserState = offlineState
	return true
}

// StateName 狀態名稱
func (ust *UserGameState) StateName() string {
	return StateNameGame
}

// StateValue 狀態值
func (ust *UserGameState) StateValue() int32 {
	return UserStateGame
}
