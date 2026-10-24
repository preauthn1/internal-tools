package pokerset

// CardID2PokerCard 將 CardID 轉換為 PokerCard
func CardID2PokerCard(i32CardID int32) *PokerCard {
	if i32CardID <= POKERCardIDNone || i32CardID >= POKERCardIDHide {
		return nil
	}

	i32Color := int32(-1)
	i32Number := int32(-1)

	if i32CardID >= POKERCardIDC1 && i32CardID <= POKERCardIDCK {
		// 梅花
		i32Color = POKERCardColorClub
		i32Number = i32CardID - POKEROneColorCardsCount*i32Color
	} else if i32CardID >= POKERCardIDD1 && i32CardID <= POKERCardIDDK {
		// 方塊
		i32Color = POKERCardColorDiamond
		i32Number = i32CardID - POKEROneColorCardsCount*i32Color
	} else if i32CardID >= POKERCardIDH1 && i32CardID <= POKERCardIDHK {
		// 紅桃
		i32Color = POKERCardColorHeart
		i32Number = i32CardID - POKEROneColorCardsCount*i32Color
	} else if i32CardID >= POKERCardIDS1 && i32CardID <= POKERCardIDSK {
		// 黑桃
		i32Color = POKERCardColorSpade
		i32Number = i32CardID - POKEROneColorCardsCount*i32Color
	} else {
		i32Color = POKERCardColorOther
		switch i32CardID {
		case POKERCardIDJ1: // 小王
			i32Number = POKERCardNoJoker1
		case POKERCardIDJ2: // 大王
			i32Number = POKERCardNoJoker2
		}
	}

	if i32Color == int32(-1) || i32Number == int32(-1) {
		return nil
	}

	pRetCard := NewPokerCard(i32Color, i32Number)

	return pRetCard
}

// CardCode2PokerCardSet 將 CardCode 轉換為 PokerCardSet
func CardCode2PokerCardSet(strCodeString string) *PokerCardSet {
	i32Len := int32(len(strCodeString))
	if (i32Len%int32(LENSingleCardCode) != 0) ||
		(i32Len/int32(LENSingleCardCode) > POKERSetCardsCount*POKERMaxSet) {
		return new(PokerCardSet)
	}

	pRetCardSet := new(PokerCardSet)
	var strColor, strNumber string
	for i := int32(0); i < i32Len; i += 2 {
		strColor = strCodeString[i : i+1]
		if _, ok := gMapI32ColorByCode[strColor]; !ok {
			continue
		} else {
			if i == i32Len-2 {
				strNumber = strCodeString[i+1:]
			} else {
				strNumber = strCodeString[i+1 : i+2]
			}

			if _, ok := gMapI32NumberByCode[strNumber]; !ok {
				continue
			}
		}

		pCard := new(PokerCard)
		pCard.I32Color = gMapI32ColorByCode[strColor]
		pCard.I32Number = gMapI32NumberByCode[strNumber]

		pRetCardSet.AddSingleCard(pCard.CodeString())
	}

	return pRetCardSet
}

// PatternString 牌型字串
func PatternString(i32Pattern int32) string {
	if i32Pattern < POKERPatternNone && i32Pattern >= POKERPatternMax {
		return ""
	}

	return gMapStrPattern[i32Pattern]
}

// PatternCardCount 各牌型張數
func PatternCardCount(i32Pattern int32) int {
	if i32Pattern < POKERPatternNone && i32Pattern >= POKERPatternMax {
		return 0
	}

	return gMapIPatternCardCount[i32Pattern]
}
