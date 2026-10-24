package pokerset

// PokerCard 牌
type PokerCard struct {
	I32Color  int32 // 花色
	I32Number int32 // 數字
	I32State  int32 // 牌的狀態
}

// NewPokerCard -----------------------------------------------------------------------------
func NewPokerCard(i32Color, i32Number int32) *PokerCard {
	if (i32Color < POKERCardColorClub || i32Color >= POKERCardColorMax) ||
		(i32Number <= POKERCardNoNone || i32Number >= POKERCardNoMax) {
		return nil
	}

	pCard := new(PokerCard)
	pCard.I32Color = i32Color
	pCard.I32Number = i32Number

	return pCard
}

// 文字字串
func (pCard *PokerCard) String() string {
	return gMapStrCardColor[pCard.I32Color] + gMapStrCardNumber[pCard.I32Number]
}

// CardID ?
func (pCard *PokerCard) CardID() int32 {
	i32CardID := POKERCardIDNone

	if pCard.I32Color >= POKERCardColorClub && pCard.I32Color <= POKERCardColorSpade {
		i32CardID = pCard.I32Color*POKEROneColorCardsCount + pCard.I32Number
	} else if pCard.I32Color == POKERCardColorOther {
		switch pCard.I32Number {
		case POKERCardNoJoker1: // 小王
			i32CardID = POKERCardIDJ1
		case POKERCardNoJoker2: // 大王
			i32CardID = POKERCardIDJ2
		case POKERCardNoHide: // 蓋牌
			i32CardID = POKERCardIDHide
		}
	}

	return i32CardID
}

// CodeString 代碼字串
func (pCard *PokerCard) CodeString() string {
	return gMapStrColorCode[pCard.I32Color] + gMapStrNumberCode[pCard.I32Number]
}

// Copy ?
func (pCard *PokerCard) Copy(pOther *PokerCard) {
	if pOther == nil {
		return
	}

	pCard.I32Number = pOther.I32Number
	pCard.I32Color = pOther.I32Color
	pCard.I32State = pOther.I32State
}
