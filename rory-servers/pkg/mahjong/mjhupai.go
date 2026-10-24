package mahjong

import "example.invalid/ServerTeam/Servers/pkg/mlog"

// 牌型
const (
	CardTypeNormal int32 = 0  // 無
	CardTypePH     int32 = 1  // 屁胡1/1
	CardTypeDDZ    int32 = 2  // 大對子1/2
	CardTypeJGD    int32 = 3  // 金勾吊3
	CardTypeQDZ    int32 = 4  // 七對子2/4
	CardTypeLQD    int32 = 5  // 龍七對4
	CardTypeQYS    int32 = 6  // 清一色3/8
	CardTypeDYJ    int32 = 7  // 帶么九3/8
	CardTypeTH     int32 = 8  // 天胡6/64
	CardTypeDH     int32 = 9  // 地胡6/64
	CardTypeHD     int32 = 10 // 海底

)

// GetCardType 獲取牌型
func (cg *CardGroup) GetCardType(pcg, gcg []int32, isDYJ int32) map[int32]int32 {
	mapCardType := make(map[int32]int32)
	// 清一色
	if cg.IsSameColor() {
		color := cg.Cards[0] / 100
		bFlag := true
		for _, pcd := range pcg {
			if color != pcd/100 {
				bFlag = false
			}
		}
		for _, gcd := range gcg {
			if color != gcd/100 {
				bFlag = false
			}
		}
		if bFlag {
			mapCardType[CardTypeQYS] = 1
		}
	}
	// 七對
	if cg.IsSevenPair() {
		mapCardType[CardTypeQDZ] = 1
	}
	// 大對子
	if cg.IsBigThree() {
		mapCardType[CardTypeDDZ] = 1
	}
	// 金勾吊
	if cg.IsJinGouDiao() {
		mapCardType[CardTypeJGD] = 1
		//mapCardType[CARD_TYPE_DDZ] = 0	// Remarked by Joecool Wang at 2019-06-05
		mapCardType[CardTypeDDZ] = 1
	}
	// 龍7對
	if cg.IsLongSevenPair() {
		mapCardType[CardTypeLQD] = 1
	}
	// 么九
	if cg.IsOneNine(pcg, gcg) && isDYJ == 1 {
		mapCardType[CardTypeDYJ] = 1
	}

	// 屁胡
	if cg.CanHuBase(0) && len(mapCardType) == 0 {
		mapCardType[CardTypePH] = 1
	}

	return mapCardType
}

// CanHuBase 是否可以胡牌,簡單胡牌（a*AAA+b*ABC+AA）
//  只檢測手牌，碰槓吃不列入檢測範圍
func (cg *CardGroup) CanHuBase(cd int32) bool {
	if cg == nil {
		return false
	}
	cgtmp := cg.CopyCardGroup() //深度拷贝
	if cgtmp == nil {
		return false
	}
	cgtmp.AddCard(cd)
	if cgtmp.Len()%3 != 2 {
		return false
	}
	// cgtmp.Cards = append(cgtmp.Cards, cd)
	// mlog.Debug("檢測胡牌，%s", cgtmp.String())
	mapCard := make(map[int32]int32, 0)
	mapCardOther := make(map[int32]int32, 0)
	// 對序牌和非序牌分類
	for _, card := range cgtmp.Cards {
		if card/100 <= ColorWan {
			mapCard[card/10]++
		} else {
			mapCardOther[card/10]++
		}
	}
	// 提取序牌可以做將牌的麻將
	pair := make([]int32, 0)
	for card, count := range mapCard {
		if count >= 2 {
			pair = append(pair, card)
		}
	}
	// 序牌沒有對子，檢測字牌等其他牌
	nPair := 0
	for card, count := range mapCardOther {
		if count == 3 {
			continue
		} else if count == 2 {
			nPair++
			pair = append(pair, card)
		} else {
			return false // 如果有單牌，直接false
		}
	}
	if nPair > 1 {
		return false
	}
	// 遍歷所有將牌，檢查是否可以組成胡牌
	if len(pair) > 0 {
		// 檢測序牌
		for _, card := range pair {
			mapTmp := make(map[int32]int32, 0)
			err := DeepCopy(&mapTmp, &mapCard)
			if err != nil {
				mlog.Error("%v", err)
				return false
			}
			// 剔除所有能組成順子的牌
			if card < 40 {
				mapTmp[card] -= 2
			}
			for i := int32(10); i < 40; i++ {
				mapTmp, _ = checkShunzi(mapTmp, i)
			}
			bCanhu := true
			for _, check := range mapTmp {
				if check == 3 {
					continue
				}
				if check != 0 {
					bCanhu = false
					break
				}
			}
			if bCanhu {
				return true
			}
		}
	}
	return false
}

// IsSevenPair 是否可以胡7對
func (cg *CardGroup) IsSevenPair() bool {
	if cg == nil {
		return false
	}
	if cg.Len() != 14 {
		return false
	}
	mapCard := make(map[int32]int32, 0)
	for _, card := range cg.Cards {
		mapCard[card/10]++
	}
	for _, count := range mapCard {
		// Modified by Joecool Wang at 2019-06-05
		// 若有四張牌是相同的，應該是龍七對，而非七對
		//if count != 2 && count != 4 {
		//
		//	return false
		//}
		if count != 2 {
			return false
		}
		// Modified by Joecool Wang at 2019-06-05 end
	}
	return true
}

// IsLongSevenPair 龍7對
func (cg *CardGroup) IsLongSevenPair() bool {
	if cg == nil {
		return false
	}
	if cg.Len() != 14 {
		return false
	}
	mapCard := make(map[int32]int32, 0)
	for _, card := range cg.Cards {
		mapCard[card/10]++
	}
	bLQD := false
	for _, count := range mapCard {
		if count != 2 && count != 4 {
			return false
		}
		if count == 4 {
			bLQD = true
		}
	}
	return bLQD
}

// IsSameColor 是否清一色
func (cg *CardGroup) IsSameColor() bool {
	if cg == nil {
		return false
	}
	if cg.Len() < 1 {
		return false
	}
	color := cg.Cards[0] / 100
	for _, card := range cg.Cards {
		if card/100 != color {
			return false
		}
	}
	return true
}

// IsBigThree 是否為大對子
func (cg *CardGroup) IsBigThree() bool {
	if cg == nil {
		return false
	}
	if cg.Len() < 2 {
		return false
	}
	mapCard := make(map[int32]int32, 0)
	for _, card := range cg.Cards {
		mapCard[card/10]++
	}
	pair := 0
	for _, count := range mapCard {
		if count == 2 {
			pair++
			continue
		}
		if count != 3 {
			return false
		}
	}

	if pair != 1 {
		return false
	}
	return true
}

// IsJinGouDiao 金勾吊,判定手牌為兩張且為對子
func (cg *CardGroup) IsJinGouDiao() bool {
	if cg.Len() == 2 {
		if cg.Cards[0]/10 == cg.Cards[1]/10 {
			return true
		}
		return false
	}
	return false
}

// IsOneNine 么九(包括碰pcg,槓gcg)
func (cg *CardGroup) IsOneNine(pcg, gcg []int32) bool {
	// 碰
	for _, cd := range pcg {
		//if GetCardColor(cd) > COLOR_WAN {	// Modified by Joecool Wang at 2019-06-05
		if GetCardColor(cd) < ColorTong || GetCardColor(cd) > ColorWan {
			return false
		}
		if GetCardNumber(cd) != 1 && GetCardNumber(cd) != 9 {
			return false
		}
	}
	// 槓
	for _, cd := range gcg {
		//if GetCardColor(cd) > COLOR_WAN {	// Modified by Joecool Wang at 2019-06-05
		if GetCardColor(cd) < ColorTong || GetCardColor(cd) > ColorWan {
			return false
		}
		if GetCardNumber(cd) != 1 && GetCardNumber(cd) != 9 {
			return false
		}
	}
	// 過濾不能和九組成順子的牌
	for _, cd := range cg.Cards {
		//if GetCardColor(cd) > COLOR_WAN {	// Modified by Joecool Wang at 2019-06-05
		if GetCardColor(cd) < ColorTong || GetCardColor(cd) > ColorWan {
			return false
		}
		if GetCardNumber(cd) > 3 && GetCardNumber(cd) < 7 {
			return false
		}
	}
	mapCard := make(map[int32]int32, 0)
	for _, card := range cg.Cards {
		mapCard[GetCardValue(card)]++
	}

	mapTmp := make(map[int32]int32, 0)
	err := DeepCopy(&mapTmp, &mapCard)
	if err != nil {
		mlog.Error("%v", err)
		return false
	}
	for i := int32(10); i < 40; i++ {
		mapTmp, _ = deleteShunzi(mapTmp, i)
	}
	jp := make([]int32, 0)
	for card, count := range mapTmp {
		if GetNumberByValue(card) == 1 || GetNumberByValue(card) == 9 {
			if count == 3 || count == 0 {
				continue
			} else if count == 2 {
				jp = append(jp, card)
			} else if count == 4 {
				jp = append(jp, card)
				jp = append(jp, card)
			} else {
				return false
			}
		} else {
			if count != 0 {
				return false
			}
		}
	}
	if len(jp) == 1 || len(jp) == 7 {
		return true
	}
	return false
}

// CanTingPai 是否可以聽牌,7對單獨加
func (cg *CardGroup) CanTingPai() bool {
	if cg == nil {
		return false
	}
	for i := int32(110); i < 900; {
		if i/10%10 == 0 {
			i += 10
			continue
		}
		cgcopy := new(CardGroup)
		for _, cd := range cg.Cards {
			cgcopy.AddCard(cd)
		}
		if cg.CanHuBase(i) {
			return true
		}
		i += 10
	}
	return false
}

// CanTingSevenPair 是否可以聽7對
func (cg *CardGroup) CanTingSevenPair() bool {
	if cg == nil {
		return false
	}
	if cg.Len() < 13 {
		return false
	}
	mapCard := make(map[int32]int32, 0)
	for _, card := range cg.Cards {
		mapCard[card/10]++
	}
	Sum := 0
	for _, count := range mapCard {
		if count == 2 {
			Sum++
		} else if count == 2 {
			Sum += 2
		}
	}
	if Sum < 6 {
		return false
	}
	return true
}

// NeedLaizi 計算不能組成順子的剩餘牌需要幾張癩子可以補齊
func NeedLaizi(mapSurplus map[int32]int32) int32 {
	//计算需要几张癞子
	sum := int32(0)
	for i := int32(1); i < 40; i++ {
		if mapSurplus[i] <= 0 {
			continue
		}
		if mapSurplus[i] == 2 {
			sum++
			continue
		}
		if mapSurplus[i] == 1 {
			mapSurplus[i]--
			if (i+1)%10 == 0 {
				sum += 2
				continue
			}
			if mapSurplus[i+1] > 0 {
				mapSurplus[i+1]--
				if (i+2)%10 == 0 {
					sum++
					continue
				}
				if mapSurplus[i+2] > 0 {
					mapSurplus[i+2]--
				} else {
					sum++
				}
			} else {
				sum++
				if mapSurplus[i+2] > 0 {
					mapSurplus[i+2]--
				} else {
					sum++
				}
			}
		}
	}
	return sum
}

// 計算組成順子後 ，去處一坎，剩餘牌
func calSurPlusMap(mapTmp map[int32]int32) map[int32]int32 {
	mapSurplus := make(map[int32]int32, 0)
	for card, count := range mapTmp {
		if count == 3 {
			continue
		}
		if count == 4 {
			mapSurplus[card] = 1
			continue
		}
		if count != 0 {
			mapSurplus[card] = count
		}
	}
	return mapSurplus
}

// 檢測順子並剔除
func checkShunzi(cd map[int32]int32, i int32) (map[int32]int32, int32) {
	if cd[i] >= 3 {
		cd[i] -= 3
	}
	if cd[i] > 0 && cd[i+1] > 0 && cd[i+2] > 0 {
		cd[i]--
		cd[i+1]--
		cd[i+2]--
		cd, i = checkShunzi(cd, i)
		return cd, i
	}
	return cd, i
}

// 去掉順子
func deleteShunzi(cd map[int32]int32, i int32) (map[int32]int32, int32) {
	if cd[i] > 0 && cd[i+1] > 0 && cd[i+2] > 0 {
		cd[i]--
		cd[i+1]--
		cd[i+2]--
		cd, i = checkShunzi(cd, i)
		return cd, i
	}
	return cd, i
}
