package qznn

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
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
	hasPokerKing           bool    //是否包含大小王
	niuniurate             int32   //玩家可以拿到牛牛的牌機率, 萬分比
	jackpotConfig          *JPConfig
}

//獎池相關設定
type JPConfig struct {
	totaljackpot   int32 //四種獎項
	roomamountrate int32 //房間彩金額度
	jackpotKind    []*JPKind
	maxbonus       int32 //最大獎金限制
}

//獎品細項幾率設定
type JPKind struct {
	kind           int32        //獎項
	syspoolrate    int32        //獎池比例
	getsyspoolrate int32        //獲得獎池概率
	bonusking      []*bonusKind //多個獎金概率
}

//獎金概率
type bonusKind struct {
	scale   int32 //獎金1比例
	getrate int32 //獲得獎金1概率
}

var nnconfig *nnConfig

//解析牛牛配置
func parseNNConfig(cfgs *config.Group) {
	cfg := cfgs.GetConfig(gamelogic.ConfigFileNameLogicConfig)
	nnconfig = new(nnConfig)
	nnconfig.qiangzhuangTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_qiangzhuang"))
	nnconfig.xianJiaBetTimesTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_xianjiabet"))
	nnconfig.cardGroupTimeout = int32(cfg.GetConfigInt(ConfigSectionName_GameLogic, "timeout_group"))
	nnconfig.baidakingConfig = cfg.GetConfigInt(ConfigSectionName_GameLogic, "baidaking")
	nnconfig.niuniurate = cfg.GetConfigInt32(ConfigSectionName_GameLogic, "niuniurate")

	nnconfig.hasPokerKing = false
	k := cfg.GetConfigInt(ConfigSectionName_GameLogic, "haspokerking")
	if k == 1 {
		nnconfig.hasPokerKing = true
		nnconfig.baidakingConfig = cfg.GetConfigInt(ConfigSectionName_GameLogic, "baidaking")
	}

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

	nnconfig.jackpotConfig = new(JPConfig)
	nnconfig.jackpotConfig.totaljackpot = cfg.GetConfigInt32("jackpolpool", "totaljackpot")
	nnconfig.jackpotConfig.roomamountrate = cfg.GetConfigInt32("jackpolpool", "roomamountrate")
	nnconfig.jackpotConfig.maxbonus = cfg.GetConfigInt32("jackpolpool", "maxbonus")

	for i := 1; i <= int(nnconfig.jackpotConfig.totaljackpot); i++ {
		jpKind := new(JPKind)
		name := fmt.Sprintf("jackpot%v", i)
		jpKind.kind = cfg.GetConfigInt32(name, "kind")
		jpKind.syspoolrate = cfg.GetConfigInt32(name, "syspoolrate")
		jpKind.getsyspoolrate = cfg.GetConfigInt32(name, "getsyspoolrate")

		bonus1 := new(bonusKind)
		bonus1.scale = cfg.GetConfigInt32(name, "bonusrate1")
		bonus1.getrate = cfg.GetConfigInt32(name, "getbonusrate1")
		bonus2 := new(bonusKind)
		bonus2.scale = cfg.GetConfigInt32(name, "bonusrate2")
		bonus2.getrate = cfg.GetConfigInt32(name, "getbonusrate2")
		bonus3 := new(bonusKind)
		bonus3.scale = cfg.GetConfigInt32(name, "bonusrate3")
		bonus3.getrate = cfg.GetConfigInt32(name, "getbonusrate3")

		jpKind.bonusking = append(jpKind.bonusking, bonus1, bonus2, bonus3)
		nnconfig.jackpotConfig.jackpotKind = append(nnconfig.jackpotConfig.jackpotKind, jpKind)
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
