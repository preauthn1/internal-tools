package gamelogic

// UserReadyState 玩家準備狀態
type UserReadyState struct{}

// OnUserFree 玩家站立
func (ust *UserReadyState) OnUserFree(user *GameUser, room *Room) bool {
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
	user.GameCountOnTheTable = 0

	return true
}

// OnUserSit 玩家坐下
func (ust *UserReadyState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
	return false
}

// OnUserQueue 玩家排隊
func (ust *UserReadyState) OnUserQueue(user *GameUser, room *Room) bool {
	ust.OnUserFree(user, room)
	return user.UserState.OnUserQueue(user, room)
}

// OnUserReady 玩家準備
func (ust *UserReadyState) OnUserReady(user *GameUser, room *Room) bool {
	return false
}

// OnUserGame 玩家遊戲中
func (ust *UserReadyState) OnUserGame(user *GameUser, room *Room) bool {
	user.UserState = gameState
	//设置用户在同一桌游戏的局数
	user.GameCountOnTheTable++

	//设置用户最后一局桌号
	user.LastTableNo = user.TableNo
	user.LastSeatNo = user.SeatNo

	return true
}

// OnGameEnd 玩家遊戲結束
func (ust *UserReadyState) OnGameEnd(user *GameUser, room *Room) bool {
	return false
}

// OnUserLeave 玩家離開房間
func (ust *UserReadyState) OnUserLeave(user *GameUser, room *Room) bool {
	ust.OnUserFree(user, room)
	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserReadyState) OnUserOffLine(user *GameUser, room *Room) bool {
	return false
}

// StateName 狀態名稱
func (ust *UserReadyState) StateName() string {
	return StateNameReady
}

// StateValue 狀態值
func (ust *UserReadyState) StateValue() int32 {
	return UserStateReady
}
