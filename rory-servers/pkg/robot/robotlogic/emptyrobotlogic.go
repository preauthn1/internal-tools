package robotlogic

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

// EmptyRobotLogic 空的機器人邏輯實現
type EmptyRobotLogic struct{}

// OnGameMessage ?
func (rbglc *EmptyRobotLogic) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {}

// OnGameStart ?
func (rbglc *EmptyRobotLogic) OnGameStart(robot *robotcomm.Robot) {}

// OnGameEnd ?
func (rbglc *EmptyRobotLogic) OnGameEnd(robot *robotcomm.Robot) {}

var emptyRobotLogic *EmptyRobotLogic

func init() {
	emptyRobotLogic = new(EmptyRobotLogic)
}

// EmptyRobotLogicCreator ?
var EmptyRobotLogicCreator = func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) RobotLogic {
	return emptyRobotLogic
}
