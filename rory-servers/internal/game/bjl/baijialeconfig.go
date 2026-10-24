package bjl

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type poollinearr struct {
	limitmoney int64 //水線金額
	rate       int32 //最有利開獎概率
}

type poollineconfig struct {
	pool []*poollinearr
}

type betarr struct {
	Baijialetype int32
	Odd          int32
	Money        int32
	Maxmoney     int64
	Weight       int32
}
type betconfig struct {
	BetValue []*betarr
	AllType  []int32
}

//創建水線配置
func newPoollineConfig(cfgs *config.Group) *poollineconfig {
	cf := new(poollineconfig)
	cf.pool = make([]*poollinearr, 0)
	cf.getPoollineconfig(cfgs)
	return cf
}

//取得水線配置
func (cf *poollineconfig) getPoollineconfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	count := scfgs.GetConfigInt("poolline", "count")
	for i := 1; i <= count; i++ {
		poollinevalue := new(poollinearr)

		key := strconv.Itoa(i)
		valuestr := scfgs.GetConfigValue("poolline", key)
		strarr := strings.Split(valuestr, ",")
		limitmoney, _ := strconv.Atoi(strarr[0])
		rate, _ := strconv.Atoi(strarr[1])

		poollinevalue.limitmoney = int64(limitmoney)
		poollinevalue.rate = int32(rate)

		cf.pool = append(cf.pool, poollinevalue)
	}
}

//创建一个实例
func newBetconfig(cfgs *config.Group) *betconfig {
	cf := new(betconfig)
	cf.BetValue = make([]*betarr, 0)
	cf.AllType = make([]int32, 0)
	cf.getBetConfig(cfgs)

	return cf
}

//获取配置
func (cf *betconfig) getBetConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	count := scfgs.GetConfigInt("baijiale", "count")
	for i := 0; i < count; i++ {
		value := new(betarr)
		value.Baijialetype = int32(i + 1)
		key := strconv.Itoa(i + 1)

		valuestr := scfgs.GetConfigValue("baijiale", key)
		strarr := strings.Split(valuestr, ",")
		gallv, _ := strconv.Atoi(strarr[0])
		odd, _ := strconv.Atoi(strarr[1])
		betmoney, _ := strconv.Atoi(strarr[2])
		maxmoney, _ := strconv.Atoi(strarr[3])
		value.Odd = int32(odd)
		value.Money = int32(betmoney)
		value.Maxmoney = int64(maxmoney)
		value.Weight = int32(gallv)

		cf.BetValue = append(cf.BetValue, value)
		cf.AllType = append(cf.AllType, value.Baijialetype)
	}
}

//通过类型获取下注的钱
func (cf *betconfig) getMoneyByType(value int32) int32 {
	for _, v := range cf.BetValue {
		if v.Baijialetype == value {
			return v.Money
		}
	}
	return 0
}

//获取此类型的倍率
func (cf *betconfig) getOddByType(value int32) int32 {
	for _, v := range cf.BetValue {
		if v.Baijialetype == value {
			return v.Odd
		}
	}
	return 0
}

//获取所有的类型
func (cf *betconfig) getAllConfigType() []int32 {
	alltype := make([]int32, 0)
	alltype = append(alltype, cf.AllType...)
	return alltype
}

//获取最大下注金额
func (cf *betconfig) getMaxMoneyByType(value int32) int64 {
	for _, v := range cf.BetValue {
		if v.Baijialetype == value {
			return v.Maxmoney
		}
	}
	return 0
}

//获取此类型的权重值
func (cf *betconfig) getWeightByType(betType int32) int32 {
	for _, v := range cf.BetValue {
		if v.Baijialetype == betType {
			return v.Weight
		}
	}
	return 0
}

//获取所有的开奖类型
func (cf *betconfig) getJiangType(alltype []int32) int32 {
	totalweight := int32(0)
	for i := 0; i < len(alltype); i++ {
		totalweight += cf.getWeightByType(alltype[i])
	}

	if totalweight <= 0 {
		mlog.Debug("权重是0：AllType=%v", alltype)
		return 0
	}

	rand.Seed(time.Now().UnixNano())
	v := rand.Int31n(totalweight) + 1
	totalweight = 0
	for i := 0; i < len(alltype); i++ {
		totalweight += cf.getWeightByType(alltype[i])
		// 原本是 v < totalweight, 但覺得有問題
		if v <= totalweight {
			return alltype[i]
		}
	}

	return 0
}

//获取是否能开出这个类型, 執行類型的開獎機率骰子, 機率為0~1000
func (cf *betconfig) getIsKai(t int32) bool {
	totalweight := int32(1000)
	v := rand.Int31n(totalweight)
	if v < cf.getWeightByType(t) {
		return true
	}

	return false
}
