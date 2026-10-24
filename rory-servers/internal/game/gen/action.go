package gen

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// IsGaming 是否遊戲中
func (game *Game) IsGaming() bool {
	return game.isgaming
}

// OnGameStart 遊戲開始
func (game *Game) OnGameStart(f func()) {
	game.GameStartFunc = func() {
		if game.IsGaming() {
			return
		}

		// 遊戲開始前調用函數
		for _, f := range game.BeforeGameStartFunc {
			f()
		}

		game.isgaming = true
		game.GameStartTime = time.Now()
		game.PlayerActionManager.Clear()

		// 通知用戶狀態改變
		for _, usr := range game.AllTableUser() {
			usr.UserState.OnUserGame(usr, game.TableManager.Room)
		}

		// 廣播遊戲開始
		game.BoardcastGameStart()
		mlog.Debug("遊戲開始...")

		// 牌局記錄
		game.GameLog.GameBegin()

		f()

		// 遊戲開始後調用函數
		for _, f := range game.AfterGameStartFunc {
			f()
		}
	}
}

// OnGameEnd 遊戲結束
func (game *Game) OnGameEnd(f func(int32)) {
	game.GameEndFunc = func() {
		// 遊戲結束前調用函數
		for _, f := range game.BeforeGameEndFunc {
			f()
		}

		game.isgaming = false
		game.LastGameEndTime = time.Now()
		game.PlayerActionManager.Clear()

		// 廣播遊戲結束
		game.BoardcastGameEnd()

		for _, usr := range game.AllGameUser() {
			usr.UserState.OnGameEnd(usr, game.TableManager.Room)
		}

		gametime := int32(time.Now().Sub(game.GameStartTime).Seconds())

		f(gametime)

		// 遊戲結束後調用函數
		for _, f := range game.AfterGameEndFunc {
			f()
		}

		mlog.Debug("遊戲結束...")

		for _, usr := range game.AllTableUser() {
			usr.UserState.OnUserReady(usr, game.TableManager.Room)
		}

		game.GameStartFunc()
	}
}

// UserEnter 玩家進入
func (game *Game) UserEnter(user *gamelogic.GameUser, tableno, seatno int32) {
	// 用戶坐下
	user.UserState.OnUserSit(user, tableno, seatno, game.TableManager.Room)

	// 通知用戶進入
	game.BoardcastUserEnter(user)

	// 用戶準備
	user.UserState.OnUserReady(user, game.TableManager.Room)

	{ // 傳送初始封包
		m := game.mGameSetting()
		game.SendMessgeFunc(user.UserID, m)
	}

	if game.IsGaming() {
		user.UserState.OnUserGame(user, game.TableManager.Room)
		return
	}

	game.GameStartFunc()
}

// UserLeave 用戶離開
func (game *Game) UserLeave(user *gamelogic.GameUser) bool {

	// 通知用戶離開
	user.UserState.OnUserLeave(user, game.TableManager.Room)

	if game.IsGaming() {
		game.isgaming = false
	}

	// 廣播用戶離開
	game.BoardcastUserLeave(user)

	return true
}
