package srnn

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
	zhuangBetTimes         []int32 //庄家倍数
	xianBetTimes           []int32 //闲家倍数
	qiangzhuangTimeout     int32   //抢庄超时时间
	xianJiaBetTimesTimeout int32   //闲家倍数超时时间
	cardGroupTimeout       int32   //组牌超时时间
	baidakingConfig        int     //百搭配置
}

var nnconfig *nnConfig

//解析牛牛配置
func parseNNConfig(cfg *config.Config) {
	nnconfig = new(nnConfig)
	nnconfig.qiangzhuangTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_qiangzhuang"))
	nnconfig.xianJiaBetTimesTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_xianjiabet"))
	nnconfig.cardGroupTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_group"))
	nnconfig.baidakingConfig = cfg.GetConfigInt(ConfigSectionName_GameLogic, "baidaking")

	zhuangbetTimes := cfg.GetConfigValue(ConfigSectionName_GameLogic, "zhuangbetTimes")
	zhuangbetArr := strings.Split(zhuangbetTimes, ",")
	for _, szbt := range zhuangbetArr {
		zbt, err := strconv.Atoi(szbt)

		if err != nil {
			mlog.Error("解析游戏逻辑配置庄家倍数发生错误,%v", err)
			return
		}

		nnconfig.zhuangBetTimes = append(nnconfig.zhuangBetTimes, int32(zbt))
	}

	xianbetTimes := cfg.GetConfigValue(ConfigSectionName_GameLogic, "xianbetTimes")
	xianbetArr := strings.Split(xianbetTimes, ",")
	for _, sxbt := range xianbetArr {
		xbt, err := strconv.Atoi(sxbt)

		if err != nil {
			mlog.Error("解析游戏逻辑配置闲家倍数发生错误,%v", err)
			return
		}

		nnconfig.xianBetTimes = append(nnconfig.xianBetTimes, int32(xbt))
	}

}

//获取配置信息
func getNNConfig() *nnConfig {
	//	nnc := new(nnConfig)
	//	nnc.zhuangBetTimes = []int32{2, 3, 4}
	//	nnc.xianBetTimes = []int32{1, 2, 3, 4}
	//	nnc.qiangzhuangTimeout = 10
	//	nnc.cardGroupTimeout = 10
	//	nnc.xianJiaBetTimesTimeout = 10

	return nnconfig
}

//庄家最小倍数
func (nnc *nnConfig) minZhuangTimes() int32 {
	return 1 //nnc.zhuangBetTimes[0]
}

//庄家最大倍数
func (nnc *nnConfig) maxZhuangTimes() int32 {
	return nnc.zhuangBetTimes[len(nnc.zhuangBetTimes)-1]
}

//闲家最小倍数
func (nnc *nnConfig) minXianTimes() int32 {
	return nnc.xianBetTimes[0]
}

//闲家最大倍数
func (nnc *nnConfig) maxXianTimes() int32 {
	return nnc.xianBetTimes[len(nnc.xianBetTimes)-1]
}

//获取服务器抢庄超时时间
func (nnc *nnConfig) getServerQiangZhuangTimeout() int32 {
	return nnc.qiangzhuangTimeout + 2
}

//获取服务器牛牛闲家倍数超时时间
func (nnc *nnConfig) getServerXianJiaBetTimesTimeout() int32 {
	return nnc.xianJiaBetTimesTimeout + 2
}

//获取服务器牛牛组牌超时时间
func (nnc *nnConfig) getServerCardGroupTimeout() int32 {
	return nnc.cardGroupTimeout + 2
}

//检查庄家倍数
func (nnc *nnConfig) checkZhuangBetTimes(times int32) bool {
	if times == 0 {
		return true
	}

	for _, t := range nnc.zhuangBetTimes {
		if t == times {
			return true
		}
	}

	return false
}

//检查闲家倍数
func (nnc *nnConfig) checkXianBetTimes(times int32) bool {

	for _, t := range nnc.xianBetTimes {
		if t == times {
			return true
		}
	}

	return false
}
