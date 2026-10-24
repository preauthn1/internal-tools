package sgj

import (
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

const (
	shuiguoji = "shuiguoji"
)

const typecount = "typecount"

const (
	pingguo    = "pingguo"
	chengzi    = "chengzi"
	mugua      = "mugua"
	lingdang   = "lingdang"
	xigua      = "xigua"
	shuangxing = "shuangxing"
	qiqi       = "qiqi"
	bar        = "bar"
)

type sgjTypeValue struct {
	sgjType       int32 //主类型
	sgjSecondType int32 //次类型
	money         int32 //钱
	oddValue      int32 //倍率
}

type sgjShowTypeConfig struct {
	configData []*sgjTypeValue
	allType    []int32
	fiveOdd    []*sgjTypeValue //倍率为5的水果，限定为苹果，用于开奖
	tenOdd     []*sgjTypeValue //倍率为10的水果
	twentyOdd  []*sgjTypeValue //倍率为20-40的水果
}

func newSgjShowTyeConfig(cfgs *config.Group) *sgjShowTypeConfig {
	config := new(sgjShowTypeConfig)
	config.configData = make([]*sgjTypeValue, 0)
	config.allType = make([]int32, 0)
	config.fiveOdd = make([]*sgjTypeValue, 0)
	config.tenOdd = make([]*sgjTypeValue, 0)
	config.twentyOdd = make([]*sgjTypeValue, 0)
	config.getConfig(cfgs)
	return config
}

func (cf *sgjShowTypeConfig) getConfig(cfgs *config.Group) {
	scfgs := cfgs.GetConfig("gameserver")
	cf.getConfigTypeValue(scfgs, shuiguoji)
}

func (cf *sgjShowTypeConfig) getConfigTypeValue(scfgs *config.Config, sgjName string) {
	count := scfgs.GetConfigInt(sgjName, typecount)
	for i := 1; i <= count; i++ {
		value := new(sgjTypeValue)
		value.sgjType = int32(i)

		keystr := strconv.Itoa(i)
		valuestr := scfgs.GetConfigValue(sgjName, keystr)
		lastindex := 0
		count := 0
		valuestr += ","
		for index, char := range valuestr {
			if char == ',' {
				substr := valuestr[lastindex:index]
				lastindex = index + 1
				v, _ := strconv.ParseInt(substr, 10, 32)
				if count == 0 {
					value.sgjSecondType = int32(v)
				} else if count == 1 {
					value.money = int32(v)
				} else {
					value.oddValue = int32(v)
				}
				count++
			}
		}

		if value.oddValue == 5 && value.sgjSecondType == 1 {
			cf.fiveOdd = append(cf.fiveOdd, value)
		} else if value.oddValue == 10 {
			cf.tenOdd = append(cf.tenOdd, value)
		} else if value.oddValue == 20 || value.oddValue == 30 || value.oddValue == 40 {
			cf.twentyOdd = append(cf.twentyOdd, value)
		}

		cf.configData = append(cf.configData, value)
		if value.money != 0 {
			cf.allType = append(cf.allType, value.sgjType)
		}
	}
}

func (cf *sgjShowTypeConfig) getMoneyByType(sgjType int32) int32 {
	for _, cfdata := range cf.configData {
		if cfdata.sgjType == sgjType {
			return cfdata.money
		}
	}
	return 0
}

func (cf *sgjShowTypeConfig) getOddByType(sgjType int32) int32 {
	for _, cfdata := range cf.configData {
		if cfdata.sgjType == sgjType {
			return cfdata.oddValue
		}
	}

	return 0
}

func (cf *sgjShowTypeConfig) getShowTypeDataByType(sgjType int32) *sgjTypeValue {
	for _, cfdata := range cf.configData {
		if cfdata.sgjType == sgjType {
			return cfdata
		}
	}

	mlog.Debug("没有找到对应的类型%v", sgjType)
	return nil
}

func (cf *sgjShowTypeConfig) getSecondTypeByType(sgjType int32) int32 {
	for _, cfdata := range cf.configData {
		if cfdata.sgjType == sgjType {
			return cfdata.sgjSecondType
		}
	}

	return 0
}
