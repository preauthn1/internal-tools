package gamelogic

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
)

// FixedGameLogicBase 固定用戶數遊戲邏輯基類
//  用戶坐滿準備即開始遊戲，遊戲結束後用戶再次排隊
type FixedGameLogicBase struct {
	*LogicBase
}

// NewFixedGameLogicBase 創建固定用戶數遊戲邏輯基類對象
func NewFixedGameLogicBase(t *Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *FixedGameLogicBase {
	rglc := new(FixedGameLogicBase)

	rglc.LogicBase = NewGameLogicBase(t, tfManager, cfgs)

	return rglc
}

// UserEnter 用戶進入
func (glc *FixedGameLogicBase) UserEnter(gameuser *GameUser, tableno int32, seatno int32) {
	if glc.isgaming {
		return
	}

	glc.LogicBase.UserEnter(gameuser, tableno, seatno)
}

// UserReady 用戶舉手
func (glc *FixedGameLogicBase) UserReady(gameuser *GameUser) {
	if glc.isgaming {
		return
	}

	glc.LogicBase.UserReady(gameuser)
}

// UserFree 用戶站立
func (glc *FixedGameLogicBase) UserFree(gameuser *GameUser) {
	if glc.isgaming {
		return
	}

	glc.LogicBase.UserFree(gameuser)
}

// UserLeave 用戶離開
func (glc *FixedGameLogicBase) UserLeave(gameuser *GameUser) bool {
	if glc.isgaming {
		return false
	}

	return glc.LogicBase.UserLeave(gameuser)
}

// GetFreeSeat 取得空閒座位
func (glc *FixedGameLogicBase) GetFreeSeat() (int32, int32) {
	if glc.isgaming {
		return FreeTableNo, FreeSeatNo
	}

	return glc.LogicBase.GetFreeSeat()
}
