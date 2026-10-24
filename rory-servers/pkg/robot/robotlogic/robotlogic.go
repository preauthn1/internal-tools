package robotlogic

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// LogicFunc ?
type LogicFunc func(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) RobotLogic

// RobotLogic 機器人邏輯
type RobotLogic interface {
	// 響應遊戲消息
	OnGameMessage(robot *robotcomm.Robot, m *msg.Message)

	// 遊戲開始
	OnGameStart(robot *robotcomm.Robot)

	// 遊戲結束
	OnGameEnd(robot *robotcomm.Robot)
}

// LogicBase 機器人邏輯基類
type LogicBase struct {
	Table        *gamelogic.Table
	TFManager    *event.TimerFuncManager
	Configs      *config.Group
	GameLogic    gamelogic.GameLogic
	EventChannel *event.Dispatcher
	RobotConfig  *robotcomm.RobotConfig
}

// NewRobotLogicBase 創建機器人邏輯基類對象實例
func NewRobotLogicBase(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *LogicBase {
	robotGlc := new(LogicBase)
	robotGlc.Table = table
	robotGlc.TFManager = tfManager
	robotGlc.Configs = configs
	robotGlc.GameLogic = glc
	robotGlc.EventChannel = evtChannel

	//解析机器人配置
	robotGlc.RobotConfig = robotCfg

	return robotGlc
}

// OnGameMessage 遊戲消息
func (lgc *LogicBase) OnGameMessage(robot *robotcomm.Robot, m *msg.Message) {}

// OnGameStart 遊戲開始
func (lgc *LogicBase) OnGameStart(robot *robotcomm.Robot) {}

// OnGameEnd 遊戲結束
func (lgc *LogicBase) OnGameEnd(robot *robotcomm.Robot) {}

// DelayRun 延遲執行
func (lgc *LogicBase) DelayRun(delay int32, f func()) {
	lgc.TFManager.NewTimer(delay, f)
}

// SendGameMessage 發送遊戲消息
func (lgc *LogicBase) SendGameMessage(ur *gamelogic.GameUser, m *msg.Message) {
	defer util.PanicErrStack()
	lgc.DelayRun(lgc.RobotConfig.GameActionRndWaitSpan(), func() {
		lgc.GameLogic.GameMessage(ur, m)
	})
}

// SendGameMessageNow 發送遊戲消息
func (lgc *LogicBase) SendGameMessageNow(ur *gamelogic.GameUser, m *msg.Message) {
	lgc.GameLogic.GameMessage(ur, m)
}
