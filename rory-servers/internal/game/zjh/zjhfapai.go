package zjh

import (	
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"	
	"example.invalid/ServerTeam/Servers/pkg/mlog"	
	"example.invalid/ServerTeam/Servers/pkg/util"
)

import (
	"fmt"
	"math/rand"
	"sort"
)

type normalWeightCal struct {
	typeWeight map[int32]int32 //邏輯處理過的權重
	weightCfg  *weightconfig   //原始的權重
}

func newNormalWeightCal(weight *weightconfig) *normalWeightCal {
	cal := new(normalWeightCal)
	cal.weightCfg = weight
	cal.typeWeight = make(map[int32]int32)
	for _, t := range weight.alltype {
		w := weight.getWeightByType(t)
		cal.typeWeight[t] = w
	}
	return cal
}

func (cal *normalWeightCal) getPaiType() int32 {
	natureType := int32(7) //自然機率牌型
	totalweight := int32(0)
	for _, w := range cal.typeWeight {
		totalweight += w
	}
	if totalweight <= 0 {
		example.invalid("权重是0：alltype=%v", cal.typeWeight)
		return 0
	}
	v := rand.Int31n(totalweight) + 1
	totalweight = 0
	for i := 0; i < len(cal.weightCfg.alltype); i++ {
		thisType := cal.weightCfg.alltype[i]
		totalweight += cal.typeWeight[thisType]
		if v <= totalweight {
			//除了自然機率外還沒翻倍過的權重需要翻倍
			if thisType != natureType && cal.typeWeight[thisType] == cal.weightCfg.getWeightByType(thisType) {
				diff := cal.typeWeight[thisType]
				cal.typeWeight[natureType] -= diff
				cal.typeWeight[thisType] += diff
			}
			return thisType
		}
	}
	return 0
}

func (cal *normalWeightCal) getWeightString() string {
	s := ""
	for i := 0; i < len(cal.weightCfg.alltype); i++ {
		thisType := cal.weightCfg.alltype[i]
		s += fmt.Sprintf("%d:%d,", thisType, cal.typeWeight[thisType])
	}
	return s[:len(s)-1]
}

func (zjh *ZhaJinHuaLogic) fapaiNormalWeight() map[int32]int {
	mapSeatUser := make(map[int32]*gamelogic.GameUser) //座位號對應的用戶
	arrSeat := make([]int32, 0)                        //座位號陣列, 處理發牌順序

	for _, user := range zjh.AllTableUser() {
		arrSeat = append(arrSeat, user.SeatNo)
		mapSeatUser[user.SeatNo] = user
	}
	r := util.RndInt(0, 2)
	if r == 1 {
		//順時針從莊開始發
		sort.Slice(arrSeat, func(i, j int) bool { return arrSeat[i] < arrSeat[j] })
	} else {
		//逆時針從莊開始發
		sort.Slice(arrSeat, func(i, j int) bool { return arrSeat[i] > arrSeat[j] })
	}

	if exists, idx := util.InArray(zjh.zhuangSeat, arrSeat); exists {
		arrSeat = append(arrSeat[idx:], arrSeat[0:idx]...)
	} else {
		mlog.Error("注意, 發牌期間找不到莊座位!!")
	}
	fapaiRecord := "[一般牌型規則] "
	result := make(map[int32]int)
	cal := newNormalWeightCal(zjh.normalweight)
	for _, seatNo := range arrSeat {
		//牌型記錄
		if ur, ok := mapSeatUser[seatNo]; ok {
			weightString := cal.getWeightString()
			cardtype := int(cal.getPaiType())
			fapaiRecord += fmt.Sprintf("[user %d_%d:%d, weight %s], ", ur.UserID, seatNo, cardtype, weightString)
			result[ur.UserID] = cardtype
		} else {
			mlog.Error("找不到座位對應的用戶")
		}
	}
	example.invalid(fapaiRecord)
	zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fapaiRecord)
	return result
}

func (zjh *ZhaJinHuaLogic) fapaiDapaiWeight() map[int32]int {
	result := make(map[int32]int)
	var uids []int32
	for _, user := range zjh.userData {
		uids = append(uids, user.User.UserID)
	}
	zjh.dapaiNumber = util.RndInt(int(zjh.dapairandom[0]), int(zjh.dapairandom[1]+1))
	//產生跟用戶數量一樣的牌型數量, 大牌數量不夠的用散牌補上
	for i := 0; i < len(uids); i++ {
		if i < zjh.dapaiNumber {
			zuhetype := zjh.dapaiweight.getPaiType(zjh.dapaiweight.alltype)
			zjh.cardType = append(zjh.cardType, int(zuhetype))
		} else {
			zjh.cardType = append(zjh.cardType, 7)
		}
	}
	sort.Ints(zjh.cardType)

	cardType := util.IntArrayShuffle(zjh.cardType)
	for i, cardtype := range cardType {
		result[uids[i]] = cardtype
	}
	fapaiRecord := fmt.Sprintf("[大牌牌型規則] %v, weight %s", result, zjh.dapaiweight.getWeightString())
	example.invalid(fapaiRecord)
	zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fmt.Sprintf("大牌份數: %d", zjh.dapaiNumber))
	zjh.beforeGameStartStep = append(zjh.beforeGameStartStep, fapaiRecord)

	return result
}
