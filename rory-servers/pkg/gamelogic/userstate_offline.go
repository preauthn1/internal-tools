package gamelogic

// UserOfflineState 玩家斷線狀態
type UserOfflineState struct {
}

// OnUserFree 玩家站立
func (ust *UserOfflineState) OnUserFree(user *GameUser, room *Room) bool {
	return false
}

// OnUserSit 玩家坐下
func (ust *UserOfflineState) OnUserSit(ur *GameUser, tableno int32, seatno int32, room *Room) bool {
	return false
}

// OnUserQueue 玩家排隊
func (ust *UserOfflineState) OnUserQueue(user *GameUser, room *Room) bool {
	return false
}

// OnUserReady 玩家準備
func (ust *UserOfflineState) OnUserReady(user *GameUser, room *Room) bool {
	return false
}

// OnUserGame 玩家遊戲中
func (ust *UserOfflineState) OnUserGame(user *GameUser, room *Room) bool {
	user.UserState = gameState
	return false
}

// OnGameEnd 玩家遊戲結束
func (ust *UserOfflineState) OnGameEnd(user *GameUser, room *Room) bool {
	return true
}

// OnUserLeave 玩家離開房間
func (ust *UserOfflineState) OnUserLeave(user *GameUser, room *Room) bool {
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
	user.UserState = freeState
	user.TableNo = FreeTableNo
	user.SeatNo = FreeSeatNo

	return true
}

// OnUserOffLine 玩家斷線
func (ust *UserOfflineState) OnUserOffLine(user *GameUser, room *Room) bool {
	return false
}

// StateName 狀態名稱
func (ust *UserOfflineState) StateName() string {
	return StateNameOffline
}

// StateValue 狀態值
func (ust *UserOfflineState) StateValue() int32 {
	return UserStateOffline
}
