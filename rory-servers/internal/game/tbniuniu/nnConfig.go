package tbnn

import (
	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"	
)

import (
	"strconv"
	"strings"
)

const (
	ConfigSectionName_GameLogic = "gamelogic"
)

//配置信息
type nnConfig struct {
	baseScoreTimeout int32   //底分超时时间
	cardGroupTimeout int32   //组牌超时时间
	baseScoreTimes   []int32 //庄家倍数
	baidakingConfig  int     //百搭配置
}

var nnconfig *nnConfig

//解析牛牛配置
func parseNNConfig(cfg *config.Config) {
	nnconfig = new(nnConfig)
	nnconfig.baseScoreTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_basescore"))
	nnconfig.cardGroupTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_group"))
	nnconfig.baidakingConfig = cfg.GetConfigInt(ConfigSectionName_GameLogic, "baidaking")

	baseScoreTimes := cfg.GetConfigValue(ConfigSectionName_GameLogic, "baseScoreTimes")
	baseScoreArr := strings.Split(baseScoreTimes, ",")
	for _, szbt := range baseScoreArr {
		zbt, err := strconv.Atoi(szbt)

		if err != nil {
			mlog.Error("解析游戏逻辑配置底注倍数发生错误,%v", err)
			return
		}

		nnconfig.baseScoreTimes = append(nnconfig.baseScoreTimes, int32(zbt))
	}

}

//获取牛牛配置
func getNNConfig() *nnConfig {
	//	nnc := new(nnConfig)
	//	nnc.baseScoreTimeout = 10
	//	nnc.cardGroupTimeout = 10
	//	nnc.baseScoreTimes = []int32{1, 2, 3, 4}

	return nnconfig
}

//最小底注
func (nnc *nnConfig) minBaseScoreTimes() int32 {
	return nnc.baseScoreTimes[0]
}

//最大底注
func (nnc *nnConfig) maxBaseScoreTimes() int32 {
	return nnc.baseScoreTimes[len(nnc.baseScoreTimes)-1]
}
