package robotcomm

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	ConfigSectionNameRobot = "robot"
)

// RobotConfig 機器人配置
type RobotConfig struct {
	Enable                   bool    // 是否開啟機器人
	LogonInterval            int32   // 機器人登錄頻率
	MinMoneyAmount           int64   // 最小金幣數
	MaxMoneyAmount           int64   // 最大金幣數
	QueueCheckWaitSpanMin    int32   // 配桌等待時間最小值
	QueueCheckWaitSpanMax    int32   // 配桌等待時間最大值
	QueueAgainWaitSpanMin    int32   // 繼續配桌等待時間最小值
	QueueAgainWaitSpanMax    int32   // 繼續配桌等待時間最大值
	TableMaxUserCount        int32   // 桌子上最大人數
	TableMinUserCount        int32   // 桌子上最小人數
	GameEndLeaveOdds         float32 // 一局遊戲結束機器人離開概率
	MaxUserGameEndLeaveOdds  float32 // 桌子上人數達到上限時機器人離開概率
	GameActionWaitSpanMin    int32   // 機器人遊戲動作等待時間最小值
	GameActionWaitSpanMax    int32   // 機器人動作等待時間最大值
	GameReadyWaitSpanMin     int32   // 機器人準備等待時間最小值
	GameReadyWaitSpanMax     int32   // 機器人準備等待時間最大值
	MaxCardsRndLvl           int32   // 機器人換牌等級隨機值
	AutoSitWaitSpanMin       int32   // 機器人主動入座等待時間最小值
	AutoSitWaitSpanMax       int32   // 機器人主動入座等待時間最大值
	StressTest               int32
	AgainstGameing           bool                                 // 進入正在遊戲中的(觀戰) 0否 1是
	MatchTableMode           int32                                // 配桌方式, 0輪詢桌子每次放一個機器人 1輪詢桌子每次放N個機率機器人
	TableMatchRobotRate      [300]int                             // index為機器人人數, value為機率%. 桌子隨機配N個機器人.
	TableMatchRobotsWiatSpan TableMatchRobotByRealUserAndWaitTime // 當真人超過某數量時調整機器人配桌時間增加
}

// TableMatchRobotByRealUserAndWaitTime 當真人超過某數量時調整機器人配桌時間增加
type TableMatchRobotByRealUserAndWaitTime struct {
	UserCount int32
	WaitSpan  int32
}

// ParseRobotConfig 解析機器人配置
func ParseRobotConfig(cfg *config.Config) *RobotConfig {
	rcfg := new(RobotConfig)

	enable := (cfg.GetConfigInt(ConfigSectionNameRobot, "enable") == 1)
	if !enable {
		return rcfg
	}

	// 登錄頻率
	robotLogonInterval := cfg.GetConfigInt32(ConfigSectionNameRobot, "logoninterval")
	// 最小金幣數
	minMoneyAmount := cfg.GetConfigInt64(ConfigSectionNameRobot, "minmneyamount")
	// 最大金幣數
	maxMoneyAmount := cfg.GetConfigInt64(ConfigSectionNameRobot, "maxmoneyamount")

	// 機器人狀態
	rcfg.Enable = enable
	// 登錄頻率
	rcfg.LogonInterval = robotLogonInterval
	// 最小金幣數
	rcfg.MinMoneyAmount = minMoneyAmount
	//最大金币数
	rcfg.MaxMoneyAmount = maxMoneyAmount

	// 配桌
	queuecheckwaitspan := cfg.GetConfigValue(ConfigSectionNameRobot, "queuecheckwaitspan")
	queuecheckSpanArr := strings.Split(queuecheckwaitspan, ",")
	queuecheckWaitSpanMin, _ := strconv.Atoi(queuecheckSpanArr[0])
	queuecheckWaitSpanMax, _ := strconv.Atoi(queuecheckSpanArr[1])
	rcfg.QueueCheckWaitSpanMin = int32(queuecheckWaitSpanMin)
	rcfg.QueueCheckWaitSpanMax = int32(queuecheckWaitSpanMax)

	// 再次配桌
	queueagainwaitspan := cfg.GetConfigValue(ConfigSectionNameRobot, "queueagainwaitspan")
	queueagainSpanArr := strings.Split(queueagainwaitspan, ",")
	queueagainWaitSpanMin, _ := strconv.Atoi(queueagainSpanArr[0])
	queueagainWaitSpanMax, _ := strconv.Atoi(queueagainSpanArr[1])
	rcfg.QueueAgainWaitSpanMin = int32(queueagainWaitSpanMin)
	rcfg.QueueAgainWaitSpanMax = int32(queueagainWaitSpanMax)

	// 配桌最大人数
	rcfg.TableMaxUserCount = cfg.GetConfigInt32(ConfigSectionNameRobot, "tablemaxusercount")
	// 配桌最小人数
	rcfg.TableMinUserCount = cfg.GetConfigInt32(ConfigSectionNameRobot, "tableminusercount")
	// 遊戲結束後離開機率
	rcfg.GameEndLeaveOdds = cfg.GetConfigFloat32(ConfigSectionNameRobot, "gameendleaveodds")
	// 桌子人數達到人數上限後離開機率
	rcfg.MaxUserGameEndLeaveOdds = cfg.GetConfigFloat32(ConfigSectionNameRobot, "maxusergameendleaveodds")

	// 遊戲動作
	gameactionwaitspan := cfg.GetConfigValue(ConfigSectionNameRobot, "gameactionwaitspan")
	gameactionSpanArr := strings.Split(gameactionwaitspan, ",")
	gameactionWaitSpanMin, _ := strconv.Atoi(gameactionSpanArr[0])
	gameactionWaitSpanMax, _ := strconv.Atoi(gameactionSpanArr[1])

	rcfg.GameActionWaitSpanMin = int32(gameactionWaitSpanMin)
	rcfg.GameActionWaitSpanMax = int32(gameactionWaitSpanMax)

	// 準備
	gamereadywaitspan := cfg.GetConfigValue(ConfigSectionNameRobot, "gamereadywaitspan")
	gamereadySpanArr := strings.Split(gamereadywaitspan, ",")
	gamereadyWaitSpanMin, _ := strconv.Atoi(gamereadySpanArr[0])
	gamereadyWaitSpanMax, _ := strconv.Atoi(gamereadySpanArr[1])

	rcfg.GameReadyWaitSpanMin = int32(gamereadyWaitSpanMin)
	rcfg.GameReadyWaitSpanMax = int32(gamereadyWaitSpanMax)

	// 主動入座
	autositywaitspan := cfg.GetConfigValue(ConfigSectionNameRobot, "autositwaitspan")
	autositySpanArr := strings.Split(autositywaitspan, ",")
	autosityWaitSpanMin, _ := strconv.Atoi(autositySpanArr[0])
	autosityWaitSpanMax, _ := strconv.Atoi(autositySpanArr[1])

	rcfg.AutoSitWaitSpanMin = int32(autosityWaitSpanMin)
	rcfg.AutoSitWaitSpanMax = int32(autosityWaitSpanMax)

	// 換牌概率
	rcfg.MaxCardsRndLvl = cfg.GetConfigInt32(ConfigSectionNameRobot, "maxcardsrndlvl")
	rcfg.StressTest = cfg.GetConfigInt32(ConfigSectionNameRobot, "stresstest")
	rcfg.AgainstGameing = (cfg.GetConfigInt32(ConfigSectionNameRobot, "againstgameing") == 0)
	rcfg.MatchTableMode = cfg.GetConfigInt32(ConfigSectionNameRobot, "matchtablemode")

	tableMatchRobotRate := cfg.GetConfigValue(ConfigSectionNameRobot, "tablematchrobotrate")
	if tableMatchRobotRate != "" {
		matchRobotArray := strings.Split(tableMatchRobotRate, "|")
		var sumRate int
		for _, v := range matchRobotArray {
			robotRateArr := strings.Split(v, ",")
			index, _ := strconv.Atoi(robotRateArr[0])
			rate, _ := strconv.Atoi(robotRateArr[1])
			sumRate += int(rate)
			rcfg.TableMatchRobotRate[index] = sumRate
		}
	}

	tableMatchRobotsWiatSpan := cfg.GetConfigValue(ConfigSectionNameRobot, "tablematchrobotswiatspan")
	if tableMatchRobotsWiatSpan != "" {
		matchRobotArray := strings.Split(tableMatchRobotsWiatSpan, ",")
		userCount, _ := strconv.Atoi(matchRobotArray[0])
		spanTime, _ := strconv.Atoi(matchRobotArray[1])
		rcfg.TableMatchRobotsWiatSpan.UserCount = int32(userCount)
		rcfg.TableMatchRobotsWiatSpan.WaitSpan = int32(spanTime)
	}

	return rcfg
}

// GameActionRndWaitSpan 生成一個隨機的遊戲動作等待時間
func (rcfg *RobotConfig) GameActionRndWaitSpan() int32 {
	return util.NewIntScope(rcfg.GameActionWaitSpanMin, rcfg.GameActionWaitSpanMax+1).GetRandomValue()
}

// GameReadyRndWaitSpan 生成一個隨機的遊戲準備等待時間
func (rcfg *RobotConfig) GameReadyRndWaitSpan() int32 {
	return util.NewIntScope(rcfg.GameReadyWaitSpanMin, rcfg.GameReadyWaitSpanMax+1).GetRandomValue()
}

// QueueAgainWaitSpan 生成一個再次配桌等待時間
func (rcfg *RobotConfig) QueueAgainWaitSpan() int32 {
	return util.NewIntScope(rcfg.QueueAgainWaitSpanMin, rcfg.QueueAgainWaitSpanMax+1).GetRandomValue()
}

// QueueAgainRealUserWaitSpan 當真人超過某數量時調整機器人配桌時間增加
func (rcfg *RobotConfig) QueueAgainRealUserWaitSpan() int32 {
	return util.NewIntScope(rcfg.QueueAgainWaitSpanMin+rcfg.TableMatchRobotsWiatSpan.WaitSpan,
		rcfg.QueueAgainWaitSpanMax+rcfg.TableMatchRobotsWiatSpan.WaitSpan).GetRandomValue()
}
