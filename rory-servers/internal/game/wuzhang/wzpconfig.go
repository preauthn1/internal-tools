package wzp

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"strconv"
	"strings"
)

//五张牌游戏配置
type wzpGameConfig struct {
	baseAmount int32   //底注
	betAmount  []int32 //投注金额
	betTimeOut int32   //投注超时时间
}

//游戏配置
var gameconfig *wzpGameConfig

const (
	GameLogicSectionName = "gamelogic"
)

//解析游戏配置
func parseGameConfig(cfg *config.Config) *wzpGameConfig {
	gameconfig = new(wzpGameConfig)
	gameconfig.baseAmount = int32(cfg.GetConfigInt(GameLogicSectionName, "baseAmount"))
	gameconfig.betTimeOut = int32(cfg.GetConfigInt(GameLogicSectionName, "betTimeout"))

	betamounts := cfg.GetConfigValue(GameLogicSectionName, "betAmount")
	betAmountArr := strings.Split(betamounts, ",")

	for _, ba := range betAmountArr {
		iba, _ := strconv.Atoi(ba)
		gameconfig.betAmount = append(gameconfig.betAmount, int32(iba))
	}

	return gameconfig
}

func (cfg *wzpGameConfig) minBetAmount() int32 {
	return cfg.betAmount[0]
}

func (cfg *wzpGameConfig) maxBetAmount() int32 {
	return cfg.betAmount[len(cfg.betAmount)-1]
}
