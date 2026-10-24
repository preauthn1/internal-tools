package hhnn

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type betarr struct {
	hhnntype int32
	odd      int32
	money    int32
	maxmoney int64
	weight   int32
}
type betconfig struct {
	betvalue []*betarr
	alltype  []int32
}

type paytable struct {
	odd         int32
	criticalodd int32
}
type oddconfig struct {
	betvalue []*paytable
}

//新建一个下注配置实例
func newBetconfig(cfgs *config.Group, name string) *betconfig {
	cf := new(betconfig)
	cf.betvalue = make([]*betarr, 0)
	cf.alltype = make([]int32, 0)
	cf.getBetConfig(cfgs, name)

	return cf
}

//新建一个倍率配置实例
func newOddconfig(cfgs *config.Group, name string) *oddconfig {
	cf := new(oddconfig)
	cf.betvalue = make([]*paytable, 0)
	cf.getOddConfig(cfgs)

	return cf
}

//获取配置文件
func (cf *betconfig) getBetConfig(cfgs *config.Group, name string) {
	scfgs := cfgs.GetConfig("gameserver")
	if name == "hhnn" {
		count := scfgs.GetConfigInt(name, "count")
		for i := 0; i < count; i++ {
			value := new(betarr)
			value.hhnntype = int32(i + 1)
			key := strconv.Itoa(i + 1)

			valuestr := scfgs.GetConfigValue("hhnn", key)
			strarr := strings.Split(valuestr, ",")
			gallv, _ := strconv.Atoi(strarr[0])
			odd, _ := strconv.Atoi(strarr[1])
			betmoney := 1 //預設下注賠率1, 開獎後調整輸贏賠率
			maxmoney, _ := strconv.Atoi(strarr[3])
			value.odd = int32(odd)
			value.money = int32(betmoney)
			value.maxmoney = int64(maxmoney)
			value.weight = int32(gallv)

			cf.betvalue = append(cf.betvalue, value)
			cf.alltype = append(cf.alltype, value.hhnntype)
		}
	} else {
		count := scfgs.GetConfigInt(name, "count")
		for i := 0; i < count; i++ {
			value := new(betarr)
			value.hhnntype = int32(i + 1)
			key := strconv.Itoa(i + 1)

			valuestr := scfgs.GetConfigValue(name, key)
			strarr := strings.Split(valuestr, ",")
			gallv, _ := strconv.Atoi(strarr[0])
			value.odd = int32(1)
			value.money = int32(1)
			value.maxmoney = int64(1)
			value.weight = int32(gallv)

			cf.betvalue = append(cf.betvalue, value)
			cf.alltype = append(cf.alltype, value.hhnntype)
		}
	}
}

//获取配置文件
func (cf *oddconfig) getOddConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	count := scfgs.GetConfigInt("paytable", "count")
	for i := 0; i < count; i++ {
		value := new(paytable)
		key := strconv.Itoa(i + 1)
		valuestr := scfgs.GetConfigValue("paytable", key)
		strarr := strings.Split(valuestr, ",")
		odd, _ := strconv.Atoi(strarr[0])
		criticalodd, _ := strconv.Atoi(strarr[1])
		value.odd = int32(odd)
		value.criticalodd = int32(criticalodd)

		cf.betvalue = append(cf.betvalue, value)
	}
}

//获取下注需要多少钱
func (cf *betconfig) getMoneyByType(value int32) int32 {
	for _, v := range cf.betvalue {
		if v.hhnntype == value {
			return v.money
		}
	}
	return 0
}

//获取下注类型的倍率
func (cf *betconfig) getOddByType(value int32) int32 {
	for _, v := range cf.betvalue {
		if v.hhnntype == value {
			return v.odd
		}
	}
	return 0
}

//获取所有的类型
func (cf *betconfig) getAllConfigType() []int32 {
	alltype := make([]int32, 0)
	alltype = append(alltype, cf.alltype...)
	return alltype
}

//获取下注的最高金额
func (cf *betconfig) getMaxMoneyByType(value int32) int64 {
	for _, v := range cf.betvalue {
		if v.hhnntype == value {
			return v.maxmoney
		}
	}
	return 0
}

//获取权重
func (cf *betconfig) getWeightByType(betType int32) int32 {
	for _, v := range cf.betvalue {
		if v.hhnntype == betType {
			return v.weight
		}
	}
	return 0
}

//获取开奖类型
func (cf *betconfig) getJiangType(alltype []int32) int32 {
	totalweight := int32(0)
	for i := 0; i < len(alltype); i++ {
		totalweight += cf.getWeightByType(alltype[i])
	}

	if totalweight <= 0 {
		mlog.Debug("权重是0：alltype=%v", alltype)
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	v := rand.Int31n(totalweight) + 1

	totalweight = 0
	for i := 0; i < len(alltype); i++ {
		totalweight += cf.getWeightByType(alltype[i])
		if v <= totalweight {
			return alltype[i]
		}
	}

	return 0
}
