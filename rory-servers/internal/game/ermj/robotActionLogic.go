package ermj

import (
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

// RobotActionManager 機器人動作管理器
// 因為遊戲邏輯會同時發送操作封包及事件封包, 這樣做可以有效避免重複操作
type RobotActionManager struct {
	robot *robotcomm.Robot
	logic *SCMJRobotLogic
	eventMsg *msg.Message
	waitOptMsg *msg.Message
	dingqueMsg *msg.Message
}

// NewRobotActionManager 創建機器人動作管理器
func NewRobotActionManager(robotLogic *SCMJRobotLogic, robot *robotcomm.Robot) *RobotActionManager {
	ral := &RobotActionManager{robot, robotLogic, nil, nil, nil}
	robotLogic.TFManager.NewTimer(robotLogic.RobotConfig.GameActionRndWaitSpan(), ral.Exec)
	return ral
}

// OnEvent 事件封包
func(ral *RobotActionManager) OnEvent(m *msg.Message) {
	ral.eventMsg= m
}

// OnWaitOp 操作封包
func(ral *RobotActionManager) OnWaitOp(m *msg.Message) {
	ral.waitOptMsg = m
}

// OnDingque 定缺封包
func(ral *RobotActionManager) OnDingque(m *msg.Message) {
	ral.dingqueMsg = m
}

// Exec 判斷處理消息的優先序
func(ral *RobotActionManager) Exec() {
	ral.logic.deleteRobotActionManager(ral.robot)

	if ral.dingqueMsg != nil {
		ral.logic.actionDingque(ral.robot, ral.dingqueMsg)
	} else if ral.eventMsg != nil {
		ral.logic.actionEvent(ral.robot, ral.eventMsg)
	} else if ral.waitOptMsg != nil {
		ral.logic.actionOutCard(ral.robot, ral.waitOptMsg)
	}
	ral.Reset()
}

// Reset 結束時清理暫存消息
func(ral *RobotActionManager) Reset() {
	ral.eventMsg= nil
	ral.waitOptMsg = nil
	ral.dingqueMsg = nil
}