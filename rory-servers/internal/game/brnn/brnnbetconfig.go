package brnn

import (
	"strconv"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/config"
)

type betarr struct {
	brnntype int32
	odd      int32
	money    int32
	maxmoney int64
	weight   winlossrate
}

type betconfig struct {
	betvalue                         []*betarr
	alltype                          []int32
	jiangWeight                      map[int]winlossrate //robot zhung win
	robotZhuangLoss                  map[int]winlossrate //robot zhung loss
	realUserZhuangSyspoolJiangWeight map[int]winlossrate //真人上莊 水線下
	realUserZhuangWinlineJiangWeight map[int]winlossrate //真人上莊 水線上
	allbetsyspool                    int32               //玩家全下吃水線表的機率
	winLossCardNumber                map[string]winLossCardNumberArr
}

type winLossCardNumberArr struct {
	cardNumber []int32
}

//開出項目的賺賠機率
type winlossrate struct {
	kaijianWeight int32 //開出此項的機率
	bigwin        int32 //大賺
	smallwin      int32 //小賺
	smallloss     int32 //小賠
	bigloss       int32 //大賠
}

//新建一个下注配置实例
func newBetconfig(cfgs *config.Group) *betconfig {
	cf := new(betconfig)
	cf.betvalue = make([]*betarr, 0)
	cf.alltype = make([]int32, 0)
	//cf.getBetConfig(cfgs)
	cf.getLogicConfig(cfgs)

	return cf
}

//遊戲邏輯配置文件, 相關的機率權重等
func (cf *betconfig) getLogicConfig(cfgs *config.Group) {
	cgf := cfgs.GetConfig("logicconfig")
	countbetoption := cgf.GetConfigInt("brnnbetoption", "count")
	cf.robotZhuangJiangWeight(cgf)
	cf.realUserSyspoolJiangWeight(cgf)
	cf.realUserWinlineJiangWeight(cgf)
	cf.setWinLossCardNumber(cgf)

	allbetsyspool := cgf.GetConfigInt32("logicsetting", "allbetsyspool")
	cf.allbetsyspool = allbetsyspool

	for i := 1; i <= countbetoption; i++ {
		cf.alltype = append(cf.alltype, int32(i))

		value := new(betarr)
		value.brnntype = int32(i)
		key := strconv.Itoa(i)

		valuestr := cgf.GetConfigValue("brnnbetoption", key)
		strarr := strings.Split(valuestr, ",")
		betmoney, _ := strconv.Atoi(strarr[0]) //一次下注的钱
		maxmoney, _ := strconv.Atoi(strarr[1]) //下注限制
		value.money = int32(betmoney)
		value.maxmoney = int64(maxmoney)

		cf.betvalue = append(cf.betvalue, value)
	}

}

func (cf *betconfig) robotZhuangJiangWeight(cgf *config.Config) {
	cf.jiangWeight = make(map[int]winlossrate)
	cf.jiangWeight[int(tongsha)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("jiangtyperate", "tongsha"),
		bigwin:        cgf.GetConfigInt32("tongsharate", "bigwin"),
		smallwin:      cgf.GetConfigInt32("tongsharate", "smallwin"),
		smallloss:     cgf.GetConfigInt32("tongsharate", "smallloss"),
		bigloss:       cgf.GetConfigInt32("tongsharate", "bigloss"),
	}
	cf.jiangWeight[int(tongpei)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("jiangtyperate", "tongpei"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.jiangWeight[int(win1loss3)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("jiangtyperate", "win1loss3"),
		bigwin:        cgf.GetConfigInt32("win1loss3rate", "bigwin"),
		smallwin:      cgf.GetConfigInt32("win1loss3rate", "smallwin"),
		smallloss:     cgf.GetConfigInt32("win1loss3rate", "smallloss"),
		bigloss:       cgf.GetConfigInt32("win1loss3rate", "bigloss"),
	}
	cf.jiangWeight[int(win2loss2)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("jiangtyperate", "win2loss2"),
		bigwin:        cgf.GetConfigInt32("win2loss2rate", "bigwin"),
		smallwin:      cgf.GetConfigInt32("win2loss2rate", "smallwin"),
		smallloss:     cgf.GetConfigInt32("win2loss2rate", "smallloss"),
		bigloss:       cgf.GetConfigInt32("win2loss2rate", "bigloss"),
	}
	cf.jiangWeight[int(win3loss1)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("jiangtyperate", "win3loss1"),
		bigwin:        cgf.GetConfigInt32("win3loss1rate", "bigwin"),
		smallwin:      cgf.GetConfigInt32("win3loss1rate", "smallwin"),
		smallloss:     cgf.GetConfigInt32("win3loss1rate", "smallloss"),
		bigloss:       cgf.GetConfigInt32("win3loss1rate", "bigloss"),
	}

	cf.robotZhuangLoss = make(map[int]winlossrate)
	cf.robotZhuangLoss[int(tongsha)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinline", "tongsha"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.robotZhuangLoss[int(tongpei)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinline", "tongpei"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.robotZhuangLoss[int(win1loss3)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinline", "win1loss3"),
		bigwin:        0,
		smallwin:      0,
		smallloss:     0,
		bigloss:       0,
	}
	cf.robotZhuangLoss[int(win2loss2)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinline", "win2loss2"),
		bigwin:        0,
		smallwin:      0,
		smallloss:     0,
		bigloss:       0,
	}
	cf.robotZhuangLoss[int(win3loss1)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinline", "win3loss1"),
		bigwin:        0,
		smallwin:      0,
		smallloss:     0,
		bigloss:       0,
	}
}

func (cf *betconfig) setWinLossCardNumber(cgf *config.Config) {
	cf.winLossCardNumber = make(map[string]winLossCardNumberArr)
	sbigwin := cgf.GetConfigValue("winlosscard", "bigwin")
	bigwinarr := strings.Split(sbigwin, ",")
	cf.setCardNumber(bigwinarr, "bigwin")

	sSmallwin := cgf.GetConfigValue("winlosscard", "smallwin")
	smallwinarr := strings.Split(sSmallwin, ",")
	cf.setCardNumber(smallwinarr, "smallwin")

	sSmallloss := cgf.GetConfigValue("winlosscard", "smallloss")
	smalllossarr := strings.Split(sSmallloss, ",")
	cf.setCardNumber(smalllossarr, "smallloss")

	sBigloss := cgf.GetConfigValue("winlosscard", "bigloss")
	biglossarr := strings.Split(sBigloss, ",")
	cf.setCardNumber(biglossarr, "bigloss")

}

func (cf *betconfig) setCardNumber(bigwinarr []string, mapName string) {
	cardNumberArr := winLossCardNumberArr{}
	for _, v := range bigwinarr {
		cardNumber, _ := strconv.Atoi(v)
		cardNumberArr.cardNumber = append(cardNumberArr.cardNumber, int32(cardNumber))
	}
	cf.winLossCardNumber[mapName] = cardNumberArr
}

func (cf *betconfig) realUserSyspoolJiangWeight(cgf *config.Config) {
	cf.realUserZhuangSyspoolJiangWeight = make(map[int]winlossrate)
	cf.realUserZhuangSyspoolJiangWeight[int(tongsha)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangsyspool", "tongsha"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangSyspoolJiangWeight[int(tongpei)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangsyspool", "tongpei"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangSyspoolJiangWeight[int(win1loss3)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangsyspool", "win1loss3"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangSyspoolJiangWeight[int(win2loss2)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangsyspool", "win2loss2"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangSyspoolJiangWeight[int(win3loss1)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangsyspool", "win3loss1"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
}

func (cf *betconfig) realUserWinlineJiangWeight(cgf *config.Config) {
	cf.realUserZhuangWinlineJiangWeight = make(map[int]winlossrate)
	cf.realUserZhuangWinlineJiangWeight[int(tongsha)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinpool", "tongsha"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangWinlineJiangWeight[int(tongpei)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinpool", "tongpei"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangWinlineJiangWeight[int(win1loss3)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinpool", "win1loss3"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangWinlineJiangWeight[int(win2loss2)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinpool", "win2loss2"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
	cf.realUserZhuangWinlineJiangWeight[int(win3loss1)] = winlossrate{
		kaijianWeight: cgf.GetConfigInt32("realuserzhuangwinpool", "win3loss1"),
		bigwin:        int32(0),
		smallwin:      int32(0),
		smallloss:     int32(0),
		bigloss:       int32(0),
	}
}

//获取配置文件
func (cf *betconfig) getBetConfig(cfgs *config.Group) {

}

//获取下注需要多少钱
func (cf *betconfig) getMoneyByType(value int32) int32 {
	for _, v := range cf.betvalue {
		if v.brnntype == value {
			return v.money
		}
	}
	return 0
}

//TODO　获取下注类型的倍率，賠率是依照贏家開的倍數
func (cf *betconfig) getOddByType(value int32) int32 {
	for _, v := range cf.betvalue {
		if v.brnntype == value {
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
		if v.brnntype == value {
			return v.maxmoney
		}
	}
	return 0
}

//获取权重
func (cf *betconfig) getWeightByType(betType int32) int32 {
	for _, v := range cf.betvalue {
		if v.brnntype == betType {
			return v.weight.bigwin //TODO 未完成機率
		}
	}
	return 0
}

//获取开奖类型
func (cf *betconfig) getJiangType(alltype []int32) int32 {
	//totalweight := int32(0)
	//for i := 0; i < len(alltype); i++ {
	//	totalweight += cf.getWeightByType(alltype[i])
	//}
	//
	//if totalweight <= 0 {
	//	mlog.Debug("权重是0：alltype=%v", alltype)
	//	return 0
	//}
	//
	//v := rand.Int31n(totalweight) + 1
	//totalweight = 0
	//for i := 0; i < len(alltype); i++ {
	//	totalweight += cf.getWeightByType(alltype[i])
	//	if v <= totalweight {
	//		return alltype[i]
	//	}
	//}

	return 0
}
