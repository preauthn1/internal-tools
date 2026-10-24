package gamelogic

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// ComputeGameLogicBase 人機對戰類游戲基類
type ComputeGameLogicBase struct {
	*LogicBase
}

// NewComputeGameLogicBase Create ComputeGameLogicBase
func NewComputeGameLogicBase(t *Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *ComputeGameLogicBase {
	rglc := new(ComputeGameLogicBase)

	rglc.LogicBase = NewGameLogicBase(t, tfManager, cfgs)

	return rglc
}

// UserEnter 玩家進入
func (glc *ComputeGameLogicBase) UserEnter(gameuser *GameUser, tableno int32, seatno int32) {
	glc.LogicBase.UserEnter(gameuser, tableno, seatno)
	glc.UserReady(gameuser)
	//if !glc.isgaming {
	//gameuser.UserState.OnUserReady(gameuser, glc.TableManager.Room)
	//glc.checkGameStart()
	//} else {
	//	gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)
	//}
}

// ChangeUserStateEnd 將玩家的遊戲狀態修改為遊戲結束狀態
func (glc *ComputeGameLogicBase) ChangeUserStateEnd(gameuser *GameUser) {
	gameuser.UserState.OnGameEnd(gameuser, glc.TableManager.Room)
}

// OnGameEnd 遊戲結束
func (glc *ComputeGameLogicBase) OnGameEnd(f func(int32)) {
	glc.GameEndFunc = func() {
		// 遊戲結束前調用函數
		for _, regf := range glc.BeforeGameEndFunc {
			regf()
		}

		glc.isgaming = false
		glc.LastGameEndTime = time.Now()
		glc.PlayerActionManager.Clear()

		// 廣播遊戲結束
		glc.BoardcastGameEnd()

		// 通知用戶狀態改變
		for _, ur := range glc.AllTableUser() {
			ur.UserState.OnGameEnd(ur, glc.TableManager.Room)
			// 廣播用戶狀態改變
			glc.BoardcastUserStateChange(ur)
			// 檢查準備超時
			if ur.IsSit() {
				glc.checkUserReadyTimer(ur)
			}
		}

		gametime := int32(time.Now().Sub(glc.GameStartTime).Seconds())
		f(gametime)

		// 遊戲結束後調用函數
		for _, regf := range glc.AfterGameEndFunc {
			regf()
		}

		mlog.Debug("遊戲結束...")
		for _, ur := range glc.AllTableUser() {
			mlog.Debug("用戶ID=%v,%v", ur.UserID, ur.UserState.StateName())
			// 用戶狀態變為準備狀態
			ur.UserState.OnUserReady(ur, glc.TableManager.Room)
		}

		// 遊戲結束後又開始遊戲
		glc.checkGameStart()
	}
}

// UserFree 用户站立
func (glc *ComputeGameLogicBase) UserFree(gameuser *GameUser) {
	mlog.Debug("user free %v", gameuser)

	// 廣播用户離開
	glc.BoardcastUserLeave(gameuser)

	// 通知用户站立
	gameuser.UserState.OnUserFree(gameuser, glc.TableManager.Room)
}

// UserLeave 用户離開
func (glc *ComputeGameLogicBase) UserLeave(gameuser *GameUser) bool {
	mlog.Debug("user leave %v", gameuser)

	glc.UserFree(gameuser)

	// 通知用户離開
	gameuser.UserState.OnUserLeave(gameuser, glc.TableManager.Room)

	glc.GameLog.RemoveGameUsers(gameuser.UserID)

	return true
}

// UserReady 用戶準備
func (glc *ComputeGameLogicBase) UserReady(gameuser *GameUser) {
	glc.LogicBase.UserReady(gameuser)
	// 如果遊戲已經開始，用戶設為遊戲狀態
	if glc.isgaming {
		gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)
	}
}

// UserRecover 玩家斷線返回
func (glc *ComputeGameLogicBase) UserRecover(gameuser *GameUser) {
	glc.LogicBase.UserRecover(gameuser)
	gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)
}
