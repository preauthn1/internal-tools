package gamelogic

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
)

// DynamicLeaveGameLogicBase 不固定人數的遊戲邏輯基類
//  遊戲過程中，其他玩家可以入座準備。遊戲未結束前用戶可換桌
//  遍歷遊戲用戶時應該用此基類中提供的方法AllGameUser
type DynamicLeaveGameLogicBase struct {
	*LogicBase
}

// NewDynamicLeaveGameLogicBase 創建固定用戶數遊戲邏輯基類對象
func NewDynamicLeaveGameLogicBase(t *Table, tfManager *event.TimerFuncManager, cfgs *config.Group) *DynamicLeaveGameLogicBase {
	rglc := new(DynamicLeaveGameLogicBase)

	rglc.LogicBase = NewGameLogicBase(t, tfManager, cfgs)

	return rglc
}

// UserEnter 用戶進入
func (glc *DynamicLeaveGameLogicBase) UserEnter(gameuser *GameUser, tableno int32, seatno int32) {
	glc.LogicBase.UserEnter(gameuser, tableno, seatno)
}

// UserReady 用戶舉手
func (glc *DynamicLeaveGameLogicBase) UserReady(gameuser *GameUser) {
	glc.LogicBase.UserReady(gameuser)
}

// UserFree 用戶站立
func (glc *DynamicLeaveGameLogicBase) UserFree(gameuser *GameUser) {
	// 如果是正在遊戲中的用戶
	// if (glc.isgaming) && glc.IsGameUser(gameuser) {
	// 	return
	// }

	glc.LogicBase.UserFree(gameuser)
}

// UserLeave 用戶離開
func (glc *DynamicLeaveGameLogicBase) UserLeave(gameuser *GameUser) bool {
	// 如果是正在遊戲中的用戶
	if glc.isgaming && glc.IsGameUser(gameuser) {
		return false
	}

	return glc.LogicBase.UserLeave(gameuser)
}
