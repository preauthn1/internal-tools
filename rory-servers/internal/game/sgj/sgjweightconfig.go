package sgj

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

type sgjDaJiangWeight struct {
	djType   int32 //大奖类型
	djWeight int32 //大奖权重
}

func (dw *sgjDaJiangWeight) String() string {
	return fmt.Sprintf("(%d_%d)", dw.djType, dw.djWeight)
}

type sgjWeightConfig struct {
	sgjType  []int32 //类型
	weight   []int32 //权重
	chiShui  int32   //吃水的百分比
	lucky    int32   //开大奖的百分比
	djConfig []*sgjDaJiangWeight
	djType   []int32
}

func newSgjWeightConfig(cfgs *config.Group, name string) *sgjWeightConfig {
	cf := new(sgjWeightConfig)
	scfgs := cfgs.GetConfig("gameserver")
	mlog.Debug("配置分数段为:%v", name)
	cf.getWeightConfig(scfgs, name)
	return cf
}

//获取权重文件
func (w *sgjWeightConfig) getWeightConfig(scfgs *config.Config, TypeName string) {
	w.sgjType = make([]int32, 0)
	w.weight = make([]int32, 0)
	w.djConfig = make([]*sgjDaJiangWeight, 0)
	w.djType = make([]int32, 0)
	count := scfgs.GetConfigInt(TypeName, typecount)
	//一般类型
	for i := 1; i <= count-6; i++ {
		w.sgjType = append(w.sgjType, int32(i))
		keystr := strconv.Itoa(i)
		w.weight = append(w.weight, int32(scfgs.GetConfigInt(TypeName, keystr)))
	}
	//大奖类型
	for i := count - 5; i <= count; i++ {
		dj := new(sgjDaJiangWeight)
		dj.djType = int32(i)
		keystr := strconv.Itoa(i)
		dj.djWeight = int32(scfgs.GetConfigInt(TypeName, keystr))
		w.djConfig = append(w.djConfig, dj)
		w.djType = append(w.djType, dj.djType)
	}

	w.chiShui = int32(scfgs.GetConfigInt(TypeName, "chishui"))
	w.lucky = int32(scfgs.GetConfigInt(TypeName, "lucky"))
}

//获取权重值通过类型
func (w *sgjWeightConfig) getWeightBySgjType(sgjType int32) int32 {
	for i := 0; i < len(w.sgjType); i++ {
		if w.sgjType[i] == sgjType {
			return w.weight[i]
		}
	}
	mlog.Debug("未找到相应的类型,%v", sgjType)
	return 0
}

//获取大奖权重
func (w *sgjWeightConfig) getDJWeightBySgjType(sgjType int32) int32 {
	for i := 0; i < len(w.djConfig); i++ {
		if w.djConfig[i].djType == sgjType {
			return w.djConfig[i].djWeight
		}
	}

	return 0
}

//吃水 1 或Lucky 2
func (w *sgjWeightConfig) getLuckyOrLose() int {
	v := rand.Int31n(100) + 1
	example.invalid("getLuckyOrLose param, chishui 1-%d, lucky %d-%d, (rand %d)", w.chiShui, w.chiShui+1, w.chiShui+w.lucky, v)
	if v <= w.chiShui {
		return 1
	}
	if v <= w.chiShui+w.lucky {
		return 2
	}
	return 0
}

//是否吃水
func (w *sgjWeightConfig) isChiShui() bool {
	v := rand.Int31n(100) + 1
	example.invalid("本轮吃水概率为：%v, rand(%d)", w.chiShui, v)
	if v <= w.chiShui {
		return true
	}

	return false
}

//是否大奖
func (w *sgjWeightConfig) isLucky() bool {
	v := rand.Int31n(100) + 1
	example.invalid("本轮大奖概率为：%v, rand(%d)", w.lucky, v)
	if v <= w.lucky {
		return true
	}

	return false
}

//获取开奖的类型
func (w *sgjWeightConfig) getJiangType(allType []int32) int32 {
	total := int32(0)
	for i := 0; i < len(allType); i++ {
		total += w.getWeightBySgjType(allType[i])
	}

	if total <= 0 {
		mlog.Debug("权重都为0：%v", allType)
		return 0
	}
	rand.Seed(time.Now().UnixNano())
	v := rand.Int31n(total) + 1
	total = 0
	jiang := int32(0)
	for i := 0; i < len(allType); i++ {
		total += w.getWeightBySgjType(allType[i])
		if v <= total {
			jiang = allType[i]
			break
		}
	}

	return jiang
}

//获取开大奖的类型
func (w *sgjWeightConfig) getDaJiangType(allType []int32) int32 {
	total := int32(0)
	for i := 0; i < len(allType); i++ {
		total += w.getDJWeightBySgjType(allType[i])
	}

	if total <= 0 {
		mlog.Debug("权重都为0：%v", allType)
		return 0
	}
	v := rand.Int31n(total) + 1
	total = 0
	jiang := int32(0)
	for i := 0; i < len(allType); i++ {
		total += w.getDJWeightBySgjType(allType[i])
		if v < total {
			jiang = allType[i]
			break
		}
	}

	return jiang
}

//获取大奖的水果列表
func (w *sgjWeightConfig) getDaJiangTypeList(allType *[][]int32) []int32 {
	alltotal := int32(0)
	totalarr := make([]int32, 0)
	for i := 0; i < len(*allType); i++ {
		tmptotal := int32(0)
		for j := 0; j < len((*allType)[i]); j++ {
			tmptotal += w.getWeightBySgjType((*allType)[i][j])

		}
		alltotal += tmptotal
		totalarr = append(totalarr, tmptotal)
	}
	if alltotal <= 0 {
		mlog.Debug("权重都为0：%v", allType)
		return nil
	}
	v := rand.Int31n(alltotal) + 1
	alltotal = 0
	jiang := (*allType)[0]
	for i := 0; i < len(*allType); i++ {
		alltotal += totalarr[i]
		if v <= alltotal {
			jiang = (*allType)[i]
			break
		}
	}

	return jiang
}

// weight string
func (w *sgjWeightConfig) String() string {
	sf := "chishui:%d,lucky:%d,sgjtype:%s,djconfig:%s\n"
	len := len(w.sgjType)
	sgjtype := ""
	for i := 0; i < len; i++ {
		sgjtype += fmt.Sprintf("(%d_%d),", w.sgjType[i], w.weight[i])
	}

	return fmt.Sprintf(sf, w.chiShui, w.lucky, sgjtype, w.djConfig)
}
