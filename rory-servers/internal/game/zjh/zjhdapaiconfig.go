package zjh

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type weightarr struct {
	dapaitype   int32
	dapaiweight int32
}

type weightconfig struct {
	betvalue []*weightarr
	alltype  []int32
}

func newNormalWeightConfig(cfgs *config.Group) *weightconfig {
	cf := new(weightconfig)
	cf.betvalue = make([]*weightarr, 0)
	cf.alltype = make([]int32, 0)
	cf.getNormalConfig(cfgs)
	return cf
}

//新建一个權重配置实例
func newDapaiWeightConfig(cfgs *config.Group) *weightconfig {
	cf := new(weightconfig)
	cf.betvalue = make([]*weightarr, 0)
	cf.alltype = make([]int32, 0)
	cf.getDapaiConfig(cfgs)
	return cf
}

//获取配置文件
func (cf *weightconfig) getDapaiConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("logicconfig")
	count := scfgs.GetConfigInt("dapaiweight", "count")
	for i := 0; i < count; i++ {
		value := new(weightarr)
		key := strconv.Itoa(i + 1)
		valuestr := scfgs.GetConfigValue("dapaiweight", key)
		strarr := strings.Split(valuestr, ",")
		dapaitype, _ := strconv.Atoi(key)
		dapaiweight, _ := strconv.Atoi(strarr[0])
		value.dapaitype = int32(dapaitype)
		value.dapaiweight = int32(dapaiweight)

		cf.betvalue = append(cf.betvalue, value)
		cf.alltype = append(cf.alltype, value.dapaitype)
	}
}

func (cf *weightconfig) getNormalConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("logicconfig")
	count := scfgs.GetConfigInt("normalweight", "count")
	for i := 0; i < count; i++ {
		value := new(weightarr)
		key := strconv.Itoa(i + 1)
		valuestr := scfgs.GetConfigValue("normalweight", key)
		strarr := strings.Split(valuestr, ",")
		dapaitype, _ := strconv.Atoi(strarr[0])
		dapaiweight, _ := strconv.Atoi(strarr[1])
		value.dapaitype = int32(dapaitype)
		value.dapaiweight = int32(dapaiweight)

		cf.betvalue = append(cf.betvalue, value)
		cf.alltype = append(cf.alltype, value.dapaitype)
	}
}

//获取开牌类型
func (cf *weightconfig) getPaiType(alltype []int32) int32 {
	totalweight := int32(0)
	for i := 0; i < len(alltype); i++ {
		totalweight += cf.getWeightByType(alltype[i])
	}

	if totalweight <= 0 {
		mlog.Debug("权重是0：alltype=%v", alltype)
		return 0
	}
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

//获取权重
func (cf *weightconfig) getWeightByType(betType int32) int32 {
	for _, v := range cf.betvalue {
		if v.dapaitype == betType {
			return v.dapaiweight
		}
	}
	return 0
}

func (cf *weightconfig) getWeightString() string {
	s := ""
	for i := 0; i < len(cf.alltype); i++ {
		thisType := cf.alltype[i]
		s += fmt.Sprintf("%d:%d,", thisType, cf.getWeightByType(thisType))
	}
	return s[:len(s)-1]
}
