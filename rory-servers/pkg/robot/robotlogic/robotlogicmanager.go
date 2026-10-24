package robotlogic

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

// Manager 機器人邏輯管理器
type Manager struct {
	logics       map[int32]RobotLogic
	robotManager *robotcomm.RobotManager
	room         *gamelogic.Room
	tfManager    *event.TimerFuncManager
	gameLogics   *gamelogic.LogicManager
}

// NewRobotLogicManager 創建機器人邏輯管理器
func NewRobotLogicManager(r *gamelogic.Room, robotCfg *robotcomm.RobotConfig, cfgs *config.Group, tfManager *event.TimerFuncManager, rblgcf LogicFunc, glcs *gamelogic.LogicManager, evtChannel *event.Dispatcher, robotManager *robotcomm.RobotManager) *Manager {
	rblgcs := new(Manager)
	rblgcs.logics = make(map[int32]RobotLogic)
	rblgcs.robotManager = robotManager
	rblgcs.room = r
	rblgcs.tfManager = tfManager
	rblgcs.gameLogics = glcs

	var i int32
	for i = 0; i < r.TableCount; i++ {
		glc := glcs.GetGameLogic(i)
		rblgc := rblgcf(r.GetTableByNo(i), tfManager, robotCfg, cfgs, glc, evtChannel)

		rblgcs.logics[i] = rblgc
	}

	return rblgcs
}

// GetRobotLogic 獲取機器人邏輯
func (rblgc *Manager) GetRobotLogic(tableno int32) RobotLogic {
	lgc, ok := rblgc.logics[tableno]

	if ok {
		return lgc
	}

	return nil
}

// Range ?
func (rblgc *Manager) Range(f func(gl RobotLogic)) {
	for _, gl := range rblgc.logics {
		f(gl)
	}
}

// OnGameMessage 通知遊戲消息
func (rblgc *Manager) OnGameMessage(user *gamelogic.GameUser, m *msg.Message) {
	robotInfo := rblgc.robotManager.GetRobot(user.UserID)
	if robotInfo == nil {
		return
	}

	tableno := robotInfo.UserInfo.TableNo
	glc := rblgc.GetRobotLogic(tableno)
	if glc == nil {
		return
	}

	glc.OnGameMessage(robotInfo, m)

}
