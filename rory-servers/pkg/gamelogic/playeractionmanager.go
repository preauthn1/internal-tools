package gamelogic

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// PlayerActionManager 玩家動作管理器
type PlayerActionManager struct {
	actions    map[string]*PlayerAction
	tfm        *event.TimerFuncManager
	actionTurn int
	gameLogic  *LogicBase
}

// NewPlayerActionManager 創建動作管理器對象
func NewPlayerActionManager(tfm *event.TimerFuncManager, gamelogic *LogicBase) *PlayerActionManager {
	pam := new(PlayerActionManager)
	pam.tfm = tfm
	pam.actionTurn = 0
	pam.actions = make(map[string]*PlayerAction)
	pam.gameLogic = gamelogic

	return pam
}

// Clear 清理動作管理器
func (pam *PlayerActionManager) Clear() {
	pam.actionTurn = 0
	pam.actions = make(map[string]*PlayerAction)
}

// GetActionTimer 獲取動作的timer
func (pam *PlayerActionManager) GetActionTimer(actname string) *event.TimerFunc {
	timername := pam.getActionTimerName(actname)

	return pam.tfm.GetTimerFunc(timername)
}

// GetActionTimerElapsedSeconds 獲取動作timer消耗的秒數
func (pam *PlayerActionManager) GetActionTimerElapsedSeconds(actname string) int32 {
	tm := pam.GetActionTimer(actname)
	if tm == nil {
		return 0
	}

	return tm.ElapsedSeconds()
}

// NewAction 創建動作
func (pam *PlayerActionManager) NewAction(actbasename string, usercount int32, timeout int32, finishFunc func()) string {
	pam.actionTurn++
	actname := fmt.Sprintf("%s_%d_%d_%d", actbasename, pam.gameLogic.TableNo(), pam.gameLogic.gameCounter, pam.actionTurn)

	act := NewPlayerAction(usercount, actname, finishFunc, timeout)

	pam.actions[actname] = act

	// 啟動超時timer
	pam.tfm.NewNamedTimer(timeout, func() {
		pam.actionFinishedCallBack(actname)
	}, pam.getActionTimerName(actname))

	//example.invalid("创建动作监听%s,%d,%d", actbasename, usercount, timeout)

	return actname
}

// CancelAction 取消動作
func (pam *PlayerActionManager) CancelAction(actname string) {
	timername := pam.getActionTimerName(actname)
	pam.tfm.RemoveTimer(timername)
}

// 獲取動作timer 名稱
func (pam *PlayerActionManager) getActionTimerName(actname string) string {
	return fmt.Sprintf("acttimer_%s", actname)
}

// OnUserAction 接收到用戶動作
func (pam *PlayerActionManager) OnUserAction(actname string, userid int32) {
	if act, ok := pam.actions[actname]; ok {
		act.recvUserAction(userid)

		mlog.Debug("收到用戶%d動作%s", userid, actname)

		if act.isAllUserFinished() {
			pam.actionFinishedCallBack(actname)
		}
	}
}

// GetFinishedActionUsers 獲取已經提交動作的用戶
func (pam *PlayerActionManager) GetFinishedActionUsers(actname string) []int32 {
	if act, ok := pam.actions[actname]; ok {
		return act.finishedUsers
	}

	return nil
}

// 動作完成毀掉函數
func (pam *PlayerActionManager) actionFinishedCallBack(actname string) {
	if theact, ok := pam.actions[actname]; ok {
		mlog.Debug("執行動作%s完成回調", actname)

		// 移除timer
		timername := pam.getActionTimerName(actname)
		pam.tfm.RemoveTimer(timername)

		// 執行回調
		theact.finishFunc()

		// 移除動作
		delete(pam.actions, actname)
	}
}

// PlayerAction 玩家動作
type PlayerAction struct {
	finishedUsers []int32
	userCount     int32
	actionName    string
	finishFunc    func()
	timeout       int32
}

// NewPlayerAction 創建用戶動作對象
func NewPlayerAction(usercount int32, actname string, finishFunc func(), timeout int32) *PlayerAction {
	pa := new(PlayerAction)
	pa.userCount = usercount
	pa.actionName = actname
	pa.finishFunc = finishFunc
	pa.timeout = timeout

	return pa
}

// 接收到用戶動作
func (pa *PlayerAction) recvUserAction(userID int32) {
	if pa.checkUserAction(userID) {
		return
	}

	pa.finishedUsers = append(pa.finishedUsers, userID)
}

// 檢查用戶是否已經完成了動作
func (pa *PlayerAction) checkUserAction(userID int32) bool {
	for _, uid := range pa.finishedUsers {
		if uid == userID {
			return true
		}
	}

	return false
}

// 是否所有用戶都完成了動作
func (pa *PlayerAction) isAllUserFinished() bool {
	if int32(len(pa.finishedUsers)) == pa.userCount {
		return true
	}

	return false
}

//已经完成动作的用户数
func (pa *PlayerAction) finishedUserCount() int32 {
	return int32(len(pa.finishedUsers))
}

// 未完成動作的用戶數
func (pa *PlayerAction) unFinishedUserCount() int32 {
	return pa.userCount - int32(len(pa.finishedUsers))
}
