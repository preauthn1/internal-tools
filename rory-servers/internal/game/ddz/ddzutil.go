package ddz

import (
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
)

type DefineCard struct {
	GameLogic *DDZGameLogic
}

func NewDefineCard(glc *DDZGameLogic) *DefineCard {
	dc := new(DefineCard)
	dc.GameLogic = glc
	return dc
}

//拆2
func (dc *DefineCard) combinationTwo(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {
	//mlog.Debug("拆分2")
	if mapCard[2] > 0 {
		ctg := new(cardTypeGroup)
		ctg.nMax = 2
		if mapCard[2] == 1 {
			//mlog.Debug("找到一个2")
			ctg.byteCardType = cardTypeSingle
		} else if mapCard[2] == 2 {
			//mlog.Debug("找到一对2")
			ctg.byteCardType = cardTypeDouble
		} else if mapCard[2] == 3 {
			//mlog.Debug("找到三个2")
			ctg.byteCardType = cardTypeThree
		} else if mapCard[2] == 4 {
			//mlog.Debug("找到炸弹2")
			ctg.byteCardType = cardTypeBomb
		}
		result.Ct = append(result.Ct, ctg)
	}
	mapCard[2] = 0
	return mapCard
}

//拆王
func (dc *DefineCard) combinationKing(mapCard map[int32]int32, result *CtGroup) map[int32]int32 {

	if mapCard[poker.PokerNumberSmallKing]+mapCard[poker.PokerNumberBigKing] == 2 {
		//mlog.Debug("找到火箭")
		ctg := new(cardTypeGroup)
		ctg.byteCardType = cardTypeRocket
		ctg.nMax = poker.PokerNumberBigKing
		ctg.nMin = poker.PokerNumberSmallKing
		result.Ct = append(result.Ct, ctg)
	} else {
		if mapCard[poker.PokerNumberSmallKing] > 0 {
			//mlog.Debug("找到小王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberSmallKing
			result.Ct = append(result.Ct, ctg)
		}
		if mapCard[poker.PokerNumberBigKing] > 0 {
			//mlog.Debug("找到大王")
			ctg := new(cardTypeGroup)
			ctg.byteCardType = cardTypeSingle
			ctg.nMax = poker.PokerNumberBigKing
			result.Ct = append(result.Ct, ctg)
		}
	}
	mapCard[poker.PokerNumberSmallKing] = 0
	mapCard[poker.PokerNumberBigKing] = 0
	return mapCard
}

// 随机拆牌
func (dc *DefineCard) combinationCardsIndex(mapCard map[int32]int32, label cardType, result *CtGroup, index int32) map[int32]int32 {

	if index == 1 {
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 2 {
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 3 {
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 4 {
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 5 {
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		return mapCard
	} else if index == 6 {
		mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
		mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
		mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
		mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
		return mapCard
	}
	mapCard = dc.combinationCards(mapCard, cardTypePlaneAndWing, result)  //飞机
	mapCard = dc.combinationCards(mapCard, cardTypeStraight, result)      //顺子
	mapCard = dc.combinationCards(mapCard, cardTypeBomb, result)          //炸弹
	mapCard = dc.combinationCards(mapCard, cardTypeThreeStraight, result) //连对
	mapCard = dc.combinationCards(mapCard, cardTypeThree, result)         //三条
	mapCard = dc.combinationCards(mapCard, cardTypeDouble, result)        //对子

	return mapCard

}

//按类型组牌
func (dc *DefineCard) combinationCards(mapCard map[int32]int32, label cardType, result *CtGroup) map[int32]int32 {
	if result == nil {
		return mapCard
	}
	if label == cardTypeDouble {
		for number, count := range mapCard {
			if count == 2 {
				//mlog.Debug("找到对子[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeDouble
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeThree {
		for number, count := range mapCard {
			if count == 3 {
				//mlog.Debug("找到三条[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeThree
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeBomb {
		for number, count := range mapCard {
			if count == 4 {
				//mlog.Debug("找到炸弹[%d]", number)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = cardTypeBomb
				ctg.nMax = number
				result.Ct = append(result.Ct, ctg)
				mapCard[number] = 0
			}
		}
	} else if label == cardTypeStraight { //顺子
		dc.checkStraightLong(mapCard, 1, 5, cardTypeStraight, result)
	} else if label == cardTypeThreeStraight { //连队
		dc.checkStraightLong(mapCard, 2, 3, cardTypeThreeStraight, result)
	} else if label == cardTypePlaneAndWing { //飞机
		dc.checkStraightLong(mapCard, 3, 2, cardTypePlaneAndWing, result)
	}

	return mapCard
}

//连对,顺子，飞机,size表示至少要几连
func (dc *DefineCard) checkStraightLong(mapCard map[int32]int32, zhang, size int32, bType cardType, result *CtGroup) map[int32]int32 {
	count := int32(0)
	min := int32(2)
	max := int32(2)
	for i := int32(2); i <= 15; i++ {
		//mlog.Debug("前count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
		if mapCard[i] >= zhang {
			count++
			max = i
		} else {
			if count >= size {
				//mlog.Debug("找到[%d,%d]一个[%d~%d]", zhang, size, min, max)
				ctg := new(cardTypeGroup)
				ctg.byteCardType = bType
				ctg.nMax = max
				ctg.nMin = min
				result.Ct = append(result.Ct, ctg)
				for j := min; j <= max; j++ {
					mapCard[j] -= zhang
				}
				mapCard = dc.checkStraightLong(mapCard, zhang, size, bType, result)
			}
			//mlog.Debug("重置count")
			count = 0
			min = i + 1
		}
		//mlog.Debug("后count[%d],min[%d],max[%d],i[%d]", count, min, max, i)
	}
	return mapCard
}

//提取对子和三张
func (dc *DefineCard) getPairAndThree(result *CtGroup) (pair, three []int32) {

	pairList := make([]int32, 0)
	threeList := make([]int32, 0)
	for _, ctg := range result.Ct {
		//三带
		if ctg.nMax == poker.PokerNumberTwo {
			continue
		}
		if ctg.byteCardType == cardTypeThree {
			threeList = append(threeList, ctg.nMax)
			continue
		}
		//对子
		if ctg.byteCardType == cardTypeDouble {
			pairList = append(pairList, ctg.nMax)
			continue
		}
	}
	return pairList, threeList
}

//提取四帶需要用到的資訊
func (dc *DefineCard) isChuFourAndTwo(result *CtGroup) cardType {
	singles, pairs, fours := make([]int32, 0), make([]int32, 0), make([]int32, 0)

	for _, s := range result.SingleCards {
		if s <= 10 {
			singles = append(singles, s)
		}
	}
	for _, ctg := range result.Ct {
		if ctg.byteCardType == cardTypeDouble && ctg.nMax <= 10 && ctg.nMax != 2 {
			pairs = append(pairs, ctg.nMax)
			continue
		}
		if ctg.byteCardType == cardTypeBomb {
			fours = append(fours, ctg.nMax)
			continue
		}
	}

	if len(fours) > 1 && len(singles) > 1 {
		mlog.Debug("決定出四帶二單")
		return cardTypeFourAndTwos
	}

	if len(fours) > 1 && len(pairs) > 1 {
		mlog.Debug("決定出四帶二對")
		return cardTypeFourAndTwod
	}
	return cardTypeNone
}

func (dc *DefineCard) getAllSingles(result *CtGroup) (allSingles []int32) {

	allSingles = make([]int32, 0)
	for _, ctg := range result.Ct {
		if ctg.byteCardType == cardTypeSingle {
			allSingles = append(allSingles, ctg.nMax)
			continue
		}
	}
	allSingles = append(allSingles, result.SingleCards...)
	sort.Slice(allSingles, func(i, j int) bool { return cardMap[allSingles[i]] < cardMap[allSingles[j]] })
	return
}

//跟新result,erasePair,eraseThree
func (dc *DefineCard) updataResult(result *CtGroup, erasePair, eraseThree []int32) *CtGroup {
	for _, cd := range erasePair {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeDouble && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	for _, cd := range eraseThree {
		for _, ctg := range result.Ct {
			if ctg.byteCardType == cardTypeThree && ctg.nMax == cd {
				ctg.byteCardType = cardTypeNone
			}
		}
	}
	resultNew := new(CtGroup)
	resultNew.Ct = make([]*cardTypeGroup, 0)
	for _, ctg := range result.Ct {
		if ctg.byteCardType != cardTypeNone {
			ctg2 := new(cardTypeGroup)
			ctg2.byteCardType = ctg.byteCardType
			ctg2.nMax = ctg.nMax
			ctg2.nMin = ctg.nMin
			resultNew.Ct = append(resultNew.Ct, ctg2)
		}
	}
	return resultNew

}

//獲取單牌
func (dc *DefineCard) combinationSingle(mapCard map[int32]int32) (mapReturn map[int32]int32, sg []int32) {
	single := make([]int32, 0)
	for number, count := range mapCard {
		if count == 1 {
			single = append(single, number)
			mapCard[number] = 0
		}
	}
	single = dc.sort(single)
	return mapCard, single
}

//排序
func (dc *DefineCard) sort(single []int32) []int32 {
	for i := 0; i < len(single); i++ {
		for j := i + 1; j < len(single); j++ {
			if single[i] > single[j] {
				single[i], single[j] = single[j], single[i]
			}
		}
	}
	return single
}

//获取手数最小的牌 ，如果手数一样获取权值最大的牌
func (dc *DefineCard) Combination(cg *poker.CardGroup) *CtGroup {
	//mlog.Debug("开始拆分手牌")
	wMax, cMin := int32(0), int32(21)
	cHands, cBigs := int32(21), int32(0)
	resultBest := new(CtGroup)
	resultBest.Ct = make([]*cardTypeGroup, 0)
	singleBest := make([]int32, 0)
	for i := 1; i <= 7; i++ {
		result := new(CtGroup)
		result.Ct = make([]*cardTypeGroup, 0)
		single := make([]int32, 0)
		mapCard := mapCardNumberCount(cg)
		mapCard = dc.combinationKing(mapCard, result)                               //王
		mapCard = dc.combinationTwo(mapCard, result)                                //2
		mapCard = dc.combinationCardsIndex(mapCard, cardTypeBomb, result, int32(i)) //剩下单牌，

		erasePair := make([]int32, 0)
		eraseThree := make([]int32, 0)
		//mlog.Debug("补张开始")
		//把对子或者三放入单牌中看能否消掉一些，这里单牌还没有带
		pairList, threeList := dc.getPairAndThree(result)

		for _, cd := range pairList {
			mapCard[cd] = 2
		}
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result) //顺子

		for _, cd := range threeList {
			mapCard[cd] = 3
			mapCard = dc.combinationCards(mapCard, cardTypeStraight, result) //顺子
			if mapCard[cd] != 3 {
				//标记，并且去除结果中的三条
				eraseThree = append(eraseThree, cd)
			} else {
				mapCard[cd] = 0 //如果没有匹配到，从重置为0
			}
		}

		for _, cd := range pairList {
			if mapCard[cd] != 2 {
				//标记，并且去除结果中的对子
				erasePair = append(erasePair, cd)
			} else {
				mapCard[cd] = 0
			}
		}
		mapCard = dc.combinationCards(mapCard, cardTypeDouble, result) //对子

		//跟新result,erasePair,eraseThree
		result = dc.updataResult(result, erasePair, eraseThree)

		mapCard, single = dc.combinationSingle(mapCard) //单牌
		for number, count := range mapCard {
			if count > 0 {
				mlog.Error("错误数据 number[%d],count[%d]", number, count)
			}
		}
		//mlog.Debug("single%v", single)

		//mlog.Debug("牌组大小[%d]", len(result.Ct))
		w, c, b := dc.CalWeight(result, single)
		//手數價值 c - (b * 2)
		cvalue := c - (b * 2)
		if cvalue < 1 {
			cvalue = 1
		}
		//mlog.Debug("计算所得[w=%d,c=%d]", w, c)
		//mlog.Debug("------------------")
		if cMin > cvalue {
			wMax, cHands, cBigs = w, c, b
			cMin = cvalue
			resultBest = result
			singleBest = single
		} else if cvalue == cMin {
			if wMax < w {
				wMax, cHands, cBigs = w, c, b
				cMin = cvalue
				resultBest = result
				singleBest = single
				//singleBest = append(singleBest, single...)
			}
		}
	}
	resultBest.SingleCards = singleBest
	resultBest.hands = cHands
	resultBest.bigs = cBigs
	//mlog.Debug("牌组大小[%d]", len(resultBest.Ct))
	//mlog.Debug("计算所得[w=%d,c=%d]", wMax, cMin)
	//mlog.Debug("single%v", resultBest.SingleCards)
	// for _, prinData := range resultBest.Ct {
	// 	mlog.Debug("resultBest%v", prinData)
	// }
	return resultBest
}

//计算权值 return 權值, 手數, 絕對大牌數
func (dc *DefineCard) CalWeight(result *CtGroup, single []int32) (w, c, b int32) {

	totalweight := int32(0)
	need := int32(0)
	pairSize := int32(0)
	count := int32(0)
	big := int32(0)
	//
	huojian := int32(7)
	zhadan := int32(7)
	feiji := int32(6)
	feijiplus := int32(3)
	liandui := int32(5)
	lianduiplus := int32(2)

	shunzi := int32(4)
	shunziplus := int32(1)
	santiao := int32(3)
	duizi := int32(2)
	danpai := int32(1)
	dadanpai := int32(1)
	if dc.GameLogic.GetServerConfig().PokerType == gamelogic.PokerTypeShuffle {
		dadanpai = int32(1)
	} else {
		dadanpai = int32(2)
	}

	for _, ctg := range result.Ct {
		//火箭
		if ctg.byteCardType == cardTypeRocket {
			totalweight += huojian
			count++
			big++
			continue
		}
		//炸弹
		if ctg.byteCardType == cardTypeBomb {
			totalweight += zhadan
			count++
			big++
			continue
		}
		//飞机
		if ctg.byteCardType == cardTypePlaneAndWing {
			size := ctg.nMax - ctg.nMin + 1
			s := feiji + (size-2)*feijiplus
			totalweight += s
			need += 2
			count++
			continue
		}
		//连对
		if ctg.byteCardType == cardTypeThreeStraight {
			size := ctg.nMax - ctg.nMin + 1
			s := liandui + (size-3)*lianduiplus
			totalweight += s
			count++
			continue
		}
		//顺子
		if ctg.byteCardType == cardTypeStraight {
			size := ctg.nMax - ctg.nMin + 1
			s := shunzi + (size-5)*shunziplus
			totalweight += s
			count++
			continue
		}
		//三带
		if ctg.byteCardType == cardTypeThree {
			totalweight += santiao
			need += 1
			count++
			continue
		}
		//对子
		if ctg.byteCardType == cardTypeDouble {
			pairSize += 1
			continue
		}
		//2 王
		if ctg.byteCardType == cardTypeSingle {
			totalweight += dadanpai
			count++
			continue
		}
	}
	singleSize := int32(len(single))
	if singleSize-need >= 0 {
		totalweight += (singleSize - need) * danpai
		count += (singleSize - need)
		need = 0
	} else {
		need = need - singleSize
	}
	if pairSize-need >= 0 {
		totalweight += (pairSize - need) * duizi
		count += (pairSize - need)
	}
	return totalweight, count, big
}

func (dc *DefineCard) SortPoker(cg *poker.CardGroup) []*poker.Card {
	sorter := sortCardGroup(*cg)
	sort.Sort(&sorter)
	return sorter.Cards
}

func (dc *DefineCard) ChuPai(cg *poker.CardGroup) {

	//return dc.ActiveChuPai(cg)
}

//主動出牌邏輯
func (dc *DefineCard) ActiveChuPai(cg *poker.CardGroup, userId int32, robot *robotcomm.Robot) (ctResult cardType, ctg *robotUserInfoLocal) {

	uc := dc.Combination(cg)
	ctg = new(robotUserInfoLocal)
	ctg.cardTypesCount = make(map[cardType]int32, 0)
	ctg.ctg = uc

	for _, data := range uc.Ct {
		ctg.cardTypesCount[data.byteCardType]++
	}
	ctg.cardTypesCount[cardTypeSingle] += int32(len(uc.SingleCards))

	//剪枝：如果能出去最后一手牌直接出
	if checkAIWithFlag(robot, RobotAIFlag2) {
		if uc.hands < 2 {
			mlog.Debug("最後一手的判斷, hands[%v], ct %s", uc.hands, uc.String())

			//三帶2
			if uc.hands == 1 {
				for _, v := range uc.Ct {
					if v.byteCardType == cardTypeThree {
						ctResult = cardTypeThree
						return
					}
				}
			}

			if uc.hands == 1 && len(uc.Ct) > 0 {
				ctResult = uc.Ct[len(uc.Ct)-1].byteCardType
			} else if len(uc.Ct) > 0 {
				ctResult = uc.Ct[0].byteCardType
			} else {
				ctResult = cardTypeSingle
			}
			return
		}
	}
	//倒数第二手時的出牌邏輯
	if checkAIWithFlag(robot, RobotAIFlag2) {
		if uc.hands < 3 {
			mlog.Debug("倒數第二手時的判斷")
			if uc.bigs > 0 { //有大牌優先出大牌
				mlog.Debug("倒數第二手 %d, 大牌數 %d, %s", uc.hands, uc.bigs, uc.String())
				for _, ct := range uc.Ct { //只出王炸
					if cardTypeRocket == ct.byteCardType {
						ctResult = cardTypeRocket
						return
					}
				}
			} else { //判斷是否有場上最大的牌
				if checkAIWithFlag(robot, RobotAIFlag4) {
					for _, ct := range uc.Ct {
						if dc.GameLogic.isBiggestCT(ct, userId) {
							mlog.Debug("倒數第二手時有場上最大的組合牌 %d, %s", ct.byteCardType, uc.String())
							ctResult = ct.byteCardType
							return
						}
					}

					for _, card := range uc.SingleCards {
						if dc.GameLogic.isBiggestSingle(card) {
							mlog.Debug("倒數第二手時有場上最大的單牌 %d, %s", card, uc.String())
							ctResult = cardTypeSingle
							return
						}
					}
				}
			}
		}
	}

	//接下來重新拆牌準備決定出牌順序
	mapCard := mapCardNumberCount(cg) //hand card list
	result := new(CtGroup)
	result.Ct = make([]*cardTypeGroup, 0)
	single := make([]int32, 0)
	mapCard = dc.combinationKing(mapCard, result)                               //王
	mapCard = dc.combinationTwo(mapCard, result)                                //2
	mapCard = dc.combinationCardsIndex(mapCard, cardTypeBomb, result, int32(5)) //炸彈, 飞机, 连对, 三条, 顺子, 对子
	mlog.Debug("userid[%v] 主動出牌牌型 %v", userId, result)

	erasePair := make([]int32, 0)
	eraseThree := make([]int32, 0)
	//把对子或者三放入单牌中看能否消掉一些，这里单牌还没有带
	pairList, threeList := dc.getPairAndThree(result)

	for _, cd := range pairList {
		mapCard[cd] = 2
	}
	mapCard = dc.combinationCards(mapCard, cardTypeStraight, result) //顺子

	for _, cd := range threeList {
		mapCard[cd] = 3
		mapCard = dc.combinationCards(mapCard, cardTypeStraight, result) //顺子
		if mapCard[cd] != 3 {
			//标记，并且去除结果中的三条
			eraseThree = append(eraseThree, cd)
		} else {
			mapCard[cd] = 0 //如果没有匹配到，从重置为0
		}
	}

	for _, cd := range pairList {
		if mapCard[cd] != 2 {
			//标记，并且去除结果中的对子
			erasePair = append(erasePair, cd)
		} else {
			mapCard[cd] = 0
		}
	}
	mapCard = dc.combinationCards(mapCard, cardTypeDouble, result) //对子

	//跟新result,erasePair,eraseThree
	result = dc.updataResult(result, erasePair, eraseThree)

	mapCard, single = dc.combinationSingle(mapCard) //单牌
	for number, count := range mapCard {
		if count > 0 {
			mlog.Error("错误数据 number[%d],count[%d]", number, count)
		}
	}
	result.SingleCards = single
	_, c, b := dc.CalWeight(result, single)
	result.hands = c
	result.bigs = b
	ctg = new(robotUserInfoLocal)
	ctg.cardTypesCount = make(map[cardType]int32, 0)
	ctg.ctg = result

	for _, data := range result.Ct {
		ctg.cardTypesCount[data.byteCardType]++
	}
	ctg.cardTypesCount[cardTypeSingle] += int32(len(result.SingleCards))
	//優先出沒有2的飛機
	for _, ct := range result.Ct {
		if ct.byteCardType == cardTypePlaneAndWing && ct.nMax != 2 {
			ctResult = cardTypePlaneAndWing
			return
		}
	}
	//優先出連隊
	for _, ct := range result.Ct {
		if ct.byteCardType == cardTypeThreeStraight && ct.nMax < 13 {
			ctResult = cardTypeThreeStraight
			return
		}
	}
	//優先出三帶
	pairs, threes := dc.getPairAndThree(result)

	if len(threes) > 0 && len(result.SingleCards) > 0 {
		mlog.Debug("決定出三帶一")
		ctResult = cardTypeThree
		return
	} else if len(threes) > 0 && len(pairs) > 0 {
		mlog.Debug("決定出三帶二")
		ctResult = cardTypeThree
		return
	} else if len(threes) > 0 {
		mlog.Debug("決定出三張")
		ctResult = cardTypeThree
		return
	}
	//優先出順子
	for _, ct := range result.Ct {
		if ct.byteCardType == cardTypeStraight {
			ctResult = cardTypeThreeStraight
			return
		}
	}

	//優先出四帶
	if ct := dc.isChuFourAndTwo(result); ct != cardTypeNone {
		ctResult = ct
		return
	}

	//對子多於單張優先出對子
	if len(pairs) > len(result.SingleCards) {
		ctResult = cardTypeDouble
		return
	}

	allSingles := dc.getAllSingles(result)
	if len(allSingles) > 0 {
		ctResult = cardTypeSingle
		return
	}
	//以上都沒有的話愈晚拆出的牌型優先出
	if len(result.Ct) > 0 {
		ctResult = result.Ct[len(result.Ct)-1].byteCardType
		return
	}
	//什麼都沒有就出單張
	ctResult = cardTypeSingle
	//mlog.Debug("牌组大小[%d]", len(ctg.ctg.Ct))
	return
}
