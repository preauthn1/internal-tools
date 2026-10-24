package cd21

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
)

type judgeRange struct {
	index int32
	value *judgeValue
}

//J=Px(define)xjudge
type judgeValue struct {
	param    int32
	numrange int32
	judge    int32
}

type zhuangBuDian struct {
	index   int32
	weight  []int32
	correct []int32
	dian    []int32
}

func (cd21 *CD21gameLogic) getConfigBysyspool(cfgs *config.Group) {
	//取得判定區域
	syscfgs := cfgs.GetConfig("syspoolconf")
	count := syscfgs.GetConfigInt("value", "key")
	for i := 0; i < count; i++ {
		value := new(judgeRange)
		value.index = int32(i + 1)
		key := strconv.Itoa(i + 1)

		valuestr := syscfgs.GetConfigValue("value", key)
		strarr := strings.Split(valuestr, ",")
		param, _ := strconv.Atoi(strarr[0])
		numrange, _ := strconv.Atoi(strarr[1])
		judge, _ := strconv.Atoi(strarr[2])
		jv := new(judgeValue)
		jv.param = int32(param)
		jv.numrange = int32(numrange)
		jv.judge = int32(judge)
		value.value = jv

		cd21.judgeRange = append(cd21.judgeRange, value)
	}

	//取得補點結果
	zhuangstr := syscfgs.GetConfigValue("budian", "zhuang")

	strarr := strings.Split(zhuangstr, ",")
	for i := 0; i < len(strarr); i++ {
		count, _ := strconv.ParseInt(strarr[i], 10, 32)
		cd21.dianJudge = append(cd21.dianJudge, int32(count))
	}
	//補點數值
	for j := 0; j < len(cd21.dianJudge); j++ {
		buDianType := new(zhuangBuDian)
		buDianType.index = int32(j + 1)
		for k := 0; k < 5; k++ {
			key := strconv.Itoa(j + 1)
			valuestr := syscfgs.GetConfigValue(key, "weight")
			weightstrarr := strings.Split(valuestr, ",")
			weight, _ := strconv.Atoi(weightstrarr[k])
			buDianType.weight = append(buDianType.weight, int32(weight))

			correctvaluestr := syscfgs.GetConfigValue(key, "correct")
			correctstrarr := strings.Split(correctvaluestr, ",")
			correct, _ := strconv.Atoi(correctstrarr[k])
			buDianType.correct = append(buDianType.correct, int32(correct))

			dianvaluestr := syscfgs.GetConfigValue(key, "dian")
			dianstrarr := strings.Split(dianvaluestr, ",")
			dian, _ := strconv.Atoi(dianstrarr[k])
			buDianType.dian = append(buDianType.dian, int32(dian))
		}
		cd21.zhuangDianJudge = append(cd21.zhuangDianJudge, buDianType)
	}

}
