package cd21

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
)

//各項倍率
const (
	rateDUIZI      = 1 //對子
	rateZHUANGBAO5 = 2 //莊家5張爆牌
	rateZHUANGBAO6 = 3
	rateZHUANGBAO7 = 4
	rateZHUANGBAO8 = 5
	rateZHUANGBAO9 = 6
	rateLUCKY10TH  = 7  //10同花
	rateLUCKYJQKTH = 8  //JQK同花
	rateLUCKY10DZ  = 9  //10JQK對子
	rateLUCKY20    = 10 //任意20點(包括A9)
	rateBIGGER13   = 11 //大於13點
	rateSMALLER13  = 12 //小於13點

)

type betarr struct {
	cd21type int32
	odd      int32
	money    int32
	maxmoney int64
}

type betconfig struct {
	betvalue []*betarr
	alltype  []int32
}

type betoddconfig struct {
	betvalue []*betarr
	alltype  []int32
}

func newBetconfig(cfgs *config.Group) *betconfig {
	cf := new(betconfig)
	cf.betvalue = make([]*betarr, 0)
	cf.alltype = make([]int32, 0)
	cf.getBetConfig(cfgs)

	return cf
}

func newBetOddconfig(cfgs *config.Group) *betoddconfig {
	cf := new(betoddconfig)
	cf.betvalue = make([]*betarr, 0)
	cf.alltype = make([]int32, 0)
	cf.getBetOddConfig(cfgs)

	return cf
}

//获取配置文件
func (cf *betconfig) getBetConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	count := scfgs.GetConfigInt("payrate", "count")
	for i := 0; i < count; i++ {
		value := new(betarr)
		value.cd21type = int32(i + 1)
		key := strconv.Itoa(i + 1)

		valuestr := scfgs.GetConfigValue("payrate", key)
		strarr := strings.Split(valuestr, ",")
		odd, _ := strconv.Atoi(strarr[0])
		betmoney, _ := strconv.Atoi(strarr[1])
		maxmoney, _ := strconv.Atoi(strarr[2])
		value.odd = int32(odd)
		value.money = int32(betmoney)
		value.maxmoney = int64(maxmoney)

		cf.betvalue = append(cf.betvalue, value)
		cf.alltype = append(cf.alltype, value.cd21type)
	}
}

//获取配置文件
func (cf *betoddconfig) getBetOddConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	count := scfgs.GetConfigInt("cody21", "count")
	for i := 0; i < count; i++ {
		value := new(betarr)
		value.cd21type = int32(i + 1)
		key := strconv.Itoa(i + 1)

		valuestr := scfgs.GetConfigValue("cody21", key)
		strarr := strings.Split(valuestr, ",")
		odd, _ := strconv.Atoi(strarr[0])
		betmoney, _ := strconv.Atoi(strarr[1])
		maxmoney, _ := strconv.Atoi(strarr[2])
		value.odd = int32(odd)
		value.money = int32(betmoney)
		value.maxmoney = int64(maxmoney)

		cf.betvalue = append(cf.betvalue, value)
		cf.alltype = append(cf.alltype, value.cd21type)
	}
}

//通过类型获取下注的钱
func (cf *betconfig) getMoneyByType(value int32) int32 {
	for _, v := range cf.betvalue {
		if v.cd21type == value {
			return v.money
		}
	}
	return 0
}

//获取此类型的倍率
func (cf *betoddconfig) getOddByType(value int32) int64 {
	for _, v := range cf.betvalue {
		if v.cd21type == value {
			return int64(v.odd)
		}
	}
	return 0
}

//获取最大下注金额
func (cf *betconfig) getMaxMoneyByType(value int32) int64 {
	for _, v := range cf.betvalue {
		if v.cd21type == value {
			return v.maxmoney
		}
	}
	return 0
}
