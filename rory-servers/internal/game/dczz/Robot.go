package dczz

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotlogic"
)

//机器人逻辑
type Robot struct {
	*robotlogic.LogicBase
}

//新建一个机器人逻辑
func NewRobot(table *gamelogic.Table, tfManager *event.TimerFuncManager, robotCfg *robotcomm.RobotConfig, configs *config.Group, glc gamelogic.GameLogic, evtChannel *event.Dispatcher) *Robot {
	pRobot := new(Robot)
	pRobot.LogicBase = robotlogic.NewRobotLogicBase(table, tfManager, robotCfg, configs, glc, evtChannel)

	return pRobot
}

//---------------------------------------------------------------------------
// 機器人押注參數
type RobotBetParam struct {
	iBetRange_Min    int // 最小下注籌碼索引值
	iBetRange_Max    int // 最大下注籌碼索引值
	iDefaultBetCount int // 預設下注數量
	iBetCount        int // 下注數量
	iBetTimes        int // 已下注次數
}

// 機器人 AI 設定
type RobotSetting struct {
	iRate_FoldAtWar int              // 進入戰爭後棄牌機率
	ary_pBetParam   []*RobotBetParam // 機器人押注參數
}

func NewRobotSetting() *RobotSetting {
	pSetting := new(RobotSetting)

	pSetting.ary_pBetParam = make([]*RobotBetParam, BetType_Max)
	for i := 0; i < len(pSetting.ary_pBetParam); i++ {
		pSetting.ary_pBetParam[i] = new(RobotBetParam)
	}

	return pSetting
}

func (pRobot *RobotSetting) Copy(pOther *RobotSetting) {
	if pOther == nil {
		return
	}

	pRobot.iRate_FoldAtWar = pOther.iRate_FoldAtWar
	for i := 0; i < len(pRobot.ary_pBetParam); i++ {
		pRobot.ary_pBetParam[i].iBetRange_Min = pOther.ary_pBetParam[i].iBetRange_Min
		pRobot.ary_pBetParam[i].iBetRange_Max = pOther.ary_pBetParam[i].iBetRange_Max
		pRobot.ary_pBetParam[i].iDefaultBetCount = pOther.ary_pBetParam[i].iDefaultBetCount
		pRobot.ary_pBetParam[i].iBetCount = pOther.ary_pBetParam[i].iBetCount
	}
}
