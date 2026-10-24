package sgj

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type sgjBetTypeConfig struct {
	betType  []int32 //下注类型
	money    []int32 //下注的钱
	maxmoney []int64 //下注的上限
}

//新建一个下注类型配置文件
func newBetTypeConfig(cfgs *config.Group) *sgjBetTypeConfig {
	betconfig := new(sgjBetTypeConfig)
	betconfig.betType = make([]int32, 0)
	betconfig.money = make([]int32, 0)
	scfgs := cfgs.GetConfig("gameserver")
	betconfig.getBetTypeConfig(scfgs)

	return betconfig
}

//获取下注配置文件
func (betconfig *sgjBetTypeConfig) getBetTypeConfig(scfgs *config.Config) {
	name := "bet"
	count := scfgs.GetConfigInt(name, typecount)
	for i := 1; i <= count; i++ {
		betconfig.betType = append(betconfig.betType, int32(i))
		keystr := strconv.Itoa(i)
		vstr := scfgs.GetConfigValue(name, keystr)
		arr := strings.Split(vstr, ",")
		money, _ := strconv.Atoi(arr[0])
		maxmoney, _ := strconv.Atoi(arr[1])
		betconfig.money = append(betconfig.money, int32(money))
		betconfig.maxmoney = append(betconfig.maxmoney, int64(maxmoney))
	}
}

//获取这个类型一次需要下注多少钱
func (betconfig *sgjBetTypeConfig) getBetMoneyByBetType(bettype int32) int32 {
	for i := 0; i < len(betconfig.betType); i++ {
		if bettype == betconfig.betType[i] {
			return betconfig.money[i]
		}
	}

	mlog.Debug("未找到类型，%v,所有类型为%v", bettype, betconfig.betType)
	return 0
}

//获取下注类型最多能下多少钱
func (betconfig *sgjBetTypeConfig) getMaxBetMoneyByBetType(bettype int32) int64 {
	for i := 0; i < len(betconfig.betType); i++ {
		if bettype == betconfig.betType[i] {
			return betconfig.maxmoney[i]
		}
	}

	mlog.Debug("未找到类型，%v,所有类型为%v", bettype, betconfig.betType)
	return 0
}
