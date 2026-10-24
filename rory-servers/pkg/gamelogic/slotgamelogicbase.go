package gamelogic

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelog"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// SlotGameLogicBase 老虎機基本邏輯
type SlotGameLogicBase struct {
	*LogicBase
}

// NewSlotGameLogicBase Create SlotGameLogicBase
func NewSlotGameLogicBase(t *Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *SlotGameLogicBase {
	rglc := new(SlotGameLogicBase)

	rglc.LogicBase = NewGameLogicBase(t, tfManager, cfgs)

	return rglc
}

// OnGameStart 遊戲開始
func (glc *SlotGameLogicBase) OnGameStart(f func()) {
	glc.GameStartFunc = func() {
		if glc.isgaming {
			return
		}

		for _, regf := range glc.BeforeGameStartFunc {
			regf()
		}

		glc.isgaming = true
		glc.GameStartTime = time.Now()
		glc.PlayerActionManager.Clear()

		for _, ur := range glc.AllTableUser() {
			ur.UserState.OnUserGame(ur, glc.TableManager.Room)
		}

		glc.BoardcastGameStart()
		mlog.Debug("遊戲開始...")

		glc.GameLog.GameBegin()

		f()

		for _, regf := range glc.AfterGameStartFunc {
			regf()
		}

		glc.gameCounter++
	}
}

// OnGameEnd 遊戲結束
func (glc *SlotGameLogicBase) OnGameEnd(f func(int32)) {
	glc.GameEndFunc = func() {

		for _, regf := range glc.BeforeGameEndFunc {
			regf()
		}

		glc.isgaming = false
		glc.LastGameEndTime = time.Now()
		glc.PlayerActionManager.Clear()

		glc.BoardcastGameEnd()

		for _, ur := range glc.AllTableUser() {
			ur.UserState.OnGameEnd(ur, glc.TableManager.Room)
		}

		gametime := int32(time.Now().Sub(glc.GameStartTime).Seconds())

		f(gametime)

		for _, regf := range glc.AfterGameEndFunc {
			regf()
		}

		mlog.Debug("遊戲結束...")
		for _, ur := range glc.AllTableUser() {
			ur.UserState.OnUserReady(ur, glc.TableManager.Room)
		}

		glc.checkGameStart()
	}
}

// UserEnter 玩家進入
func (glc *SlotGameLogicBase) UserEnter(gameuser *GameUser, tableno int32, seatno int32) {
	// glc.LogicBase.UserEnter(gameuser, tableno, seatno)

	st := glc.TableManager.GetSeatByNo(seatno)
	if st == nil {
		mlog.Warn("用户%v坐下,%d,%d,但座位不存在", gameuser, tableno, seatno)
		return
	}
	if st.SeatUser != nil {
		mlog.Warn("用户%v坐下,%d,%d,但座位上已经有用户%v", gameuser, tableno, seatno, st.SeatUser)
		return
	}

	// 通知用户坐下
	gameuser.UserState.OnUserSit(gameuser, tableno, seatno, glc.TableManager.Room)

	// 通知用户进入
	glc.BoardcastUserEnter(gameuser)

	//gamelog 初始化用戶
	urs := make(map[int32]*gamelog.UserData)
	urs[gameuser.UserID] = new(gamelog.UserData)
	urs[gameuser.UserID].Score = gameuser.Score
	urs[gameuser.UserID].TradeID = gameuser.TradeID
	glc.GameLog.UpdateGameUsers(urs)

	glc.UserReady(gameuser)
}

// UserFree 用户站立
func (glc *SlotGameLogicBase) UserFree(gameuser *GameUser) {
	mlog.Trace("user free %v", gameuser)

	// 廣播用户離開
	glc.BoardcastUserLeave(gameuser)

	// 通知用户站立
	gameuser.UserState.OnUserFree(gameuser, glc.TableManager.Room)
}

// UserLeave 用户離開
func (glc *SlotGameLogicBase) UserLeave(gameuser *GameUser) bool {
	mlog.Trace("user leave %v", gameuser)

	glc.UserFree(gameuser)

	// 通知用户離開
	gameuser.UserState.OnUserLeave(gameuser, glc.TableManager.Room)

	glc.GameLog.RemoveGameUsers(gameuser.UserID)

	return true
}

// UserReady 用戶準備
func (glc *SlotGameLogicBase) UserReady(gameuser *GameUser) {
	// glc.LogicBase.UserReady(gameuser)
	gameuser.UserState.OnUserReady(gameuser, glc.TableManager.Room)

	glc.checkGameStart()

	if glc.isgaming {
		gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)
	}
}

// UserRecover 玩家斷線返回
func (glc *SlotGameLogicBase) UserRecover(gameuser *GameUser) {
	// glc.LogicBase.UserRecover(gameuser)

	glc.notifyUserRecover(gameuser)

	gameuser.UserState.OnUserGame(gameuser, glc.TableManager.Room)
}
