package fgsb

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// 牌型
const (
	PatternNone          int = iota // 無
	PatternHighCard                 // 高牌
	PatternPair                     // 對子
	PatternFlush                    // 同花
	PatternStraight                 // 順子
	PatternThreeOfAKind             // 三條
	PatternStraightFlush            // 同花順
	PatternRoyalFlush               // 皇家同花順
	PatternMax
)

// StrPattern 牌型字串
var StrPattern = map[int]string{
	PatternNone:          "无",
	PatternHighCard:      "高牌",
	PatternPair:          "对子",
	PatternFlush:         "同花",
	PatternStraight:      "顺子",
	PatternThreeOfAKind:  "三条",
	PatternStraightFlush: "同花顺",
	PatternRoyalFlush:    "皇家同花顺",
}

// 比牌結果
const (
	CompareEven   int = iota // 平手
	CompareDealer            // 莊家贏
	ComparePlayer            // 玩家贏
)

// Poker 撲克牌組
type Poker struct {
	*poker.CardGroup
}

// NewPoker 創建撲克牌組資料
func NewPoker() *Poker {
	nPoker := new(Poker)
	nPoker.CardGroup = poker.NewEmptyCardGroup(0)

	return nPoker
}

// SetCards 設定用戶牌組
func (pokerSet *Poker) SetCards(cards *poker.CardGroup) {
	pokerSet.CardGroup = poker.NewCardGroup(cards.Cards)
	pokerSet.Sort(true)
}

// Sort 给牌排序
// asc = true：由小到大排序，asc = false：由大到小排序
func (pokerSet *Poker) Sort(asc bool) {
	arrlen := len(pokerSet.Cards)
	for i := 0; i < arrlen; i++ {
		for j := 0; j+1 < arrlen-i; j++ {
			card1 := pokerSet.Cards[j].CardNumber*10 + pokerSet.Cards[j].Color
			card2 := pokerSet.Cards[j+1].CardNumber*10 + pokerSet.Cards[j+1].Color
			if asc {
				if card1 > card2 {
					pokerSet.Swap(j, j+1)
				}
			} else {
				if card1 < card2 {
					pokerSet.Swap(j, j+1)
				}
			}

		}
	}
}

// Pattern 取得牌組牌型
func (pokerSet *Poker) Pattern(bIsDealer bool) int {
	iRetPattern := PatternNone

	if !bIsDealer {
		// 玩家
		iRetPattern = pokerSet.CheckPattern()
	} else {
		// 莊家
		bCheckedOK := true

		if pokerSet.CheckPattern() == PatternHighCard {
			bCheckedOK = pokerSet.CheckIfHasAKQ()
		}

		if bCheckedOK {
			iRetPattern = pokerSet.CheckPattern()
		} else {
			iRetPattern = PatternNone
		}
	}

	return iRetPattern
}

// PatternString 牌型文字
func (pokerSet *Poker) PatternString(bIsDealer bool) string {
	strRet := ""
	iPattern := pokerSet.Pattern(bIsDealer)

	if !bIsDealer {
		strRet = PatternString(iPattern)
	} else {
		if iPattern == PatternNone {
			strRet = CannotBeAGame
		} else {
			strRet = PatternString(iPattern)
		}
	}

	return strRet
}

// CheckPattern check 牌組牌型
func (pokerSet *Poker) CheckPattern() int {
	if pokerSet.Len() <= 0 {
		return PatternNone
	}

	// 皇家同花順
	if pokerSet.isRoyalFlush() {
		return PatternRoyalFlush
	}

	// 同花順
	if pokerSet.isStraightFlush() {
		return PatternStraightFlush
	}

	// 三條
	if pokerSet.isThreeOfAKind() {
		return PatternThreeOfAKind
	}

	// 順子
	if pokerSet.isStraight() {
		return PatternStraight
	}

	// 同花
	if pokerSet.isFlush() {
		return PatternFlush
	}

	// 對子
	if pokerSet.isPair() {
		return PatternPair
	}

	return PatternHighCard
}

// CheckIfHasAKQ 檢查是否包含 A 或 K 或 Q
func (pokerSet *Poker) CheckIfHasAKQ() bool {
	AKCheck := map[int32]bool{
		12:                 false,
		13:                 false,
		poker.PokerNumberA: false,
	}

	for i := 0; i < pokerSet.Len(); i++ {
		card := pokerSet.Cards[i]
		if card != nil {
			if _, ok := AKCheck[card.CardNumber]; ok {
				return true
			}
		}
	}

	return false
}

// Compare 比牌，compareColor：是否要比花色,compareSingle:是否比較牌型後，比較單牌
func (pokerSet *Poker) Compare(Cards *Poker, compareColor bool, compareSingle bool) int {
	if pokerSet == nil || Cards == nil {
		return CompareEven
	}

	if pokerSet.Len() <= 0 || Cards.Len() <= 0 {
		return CompareEven
	}

	if pokerSet.Len() != Cards.Len() {
		return CompareEven
	}

	// 先比牌型
	if pokerSet.CheckPattern() > Cards.CheckPattern() {
		return CompareDealer
	} else if pokerSet.CheckPattern() < Cards.CheckPattern() {
		return ComparePlayer
	}

	tempCard1 := new(Poker)
	tempCard1.SetCards(pokerSet.CardGroup)
	tempCard1.Sort(false)

	tempCard2 := new(Poker)
	tempCard2.SetCards(Cards.CardGroup)
	tempCard2.Sort(false)

	// 再比單張
	switch tempCard1.CheckPattern() {
	case PatternHighCard: // 高牌
		fallthrough
	case PatternFlush: // 同花
		for i := 0; i < tempCard1.Len(); i++ {
			card1 := tempCard1.Cards[i]
			card2 := tempCard2.Cards[i]
			if card1.CardNumber > card2.CardNumber {
				return CompareDealer
			} else if card1.CardNumber < card2.CardNumber {
				return ComparePlayer
			}
		}
	case PatternPair: // 對子
		card1 := tempCard1.Cards[0]
		card1_2 := tempCard1.Cards[2]
		if tempCard1.Cards[1].CardNumber == tempCard1.Cards[2].CardNumber {
			card1 = tempCard1.Cards[1]
			card1_2 = tempCard1.Cards[0]
		}

		card2 := tempCard2.Cards[0]
		card2_2 := tempCard2.Cards[2]
		if tempCard2.Cards[1].CardNumber == tempCard2.Cards[2].CardNumber {
			card2 = tempCard2.Cards[1]
			card2_2 = tempCard2.Cards[0]
		}

		if card1.CardNumber > card2.CardNumber {
			return CompareDealer
		} else if card1.CardNumber < card2.CardNumber {
			return ComparePlayer
		} else if card1_2.CardNumber > card2_2.CardNumber {
			return CompareDealer
		} else if card1_2.CardNumber < card2_2.CardNumber {
			return ComparePlayer
		}
	case PatternThreeOfAKind: // 三條
		fallthrough
	case PatternStraight: // 順子
		fallthrough
	case PatternStraightFlush: // 同花順
		card1 := tempCard1.Cards[0]
		card2 := tempCard2.Cards[0]
		if card1.CardNumber > card2.CardNumber {
			return CompareDealer
		} else if card1.CardNumber < card2.CardNumber {
			return ComparePlayer
		}
	}

	return CompareEven
}

/**確認牌型**/

// 是否為對子
func (pokerSet *Poker) isPair() bool {
	tempCard := new(Poker)
	tempCard.SetCards(pokerSet.CardGroup)
	tempCard.Sort(false)

	if tempCard.Cards[0].CardNumber == tempCard.Cards[1].CardNumber && tempCard.Cards[1].CardNumber != tempCard.Cards[2].CardNumber {
		return true
	}

	if tempCard.Cards[0].CardNumber != tempCard.Cards[1].CardNumber && tempCard.Cards[1].CardNumber == tempCard.Cards[2].CardNumber {
		return true
	}

	return false
}

// 是否為同花
func (pokerSet *Poker) isFlush() bool {
	if pokerSet.Cards[0].Color == pokerSet.Cards[1].Color && pokerSet.Cards[1].Color == pokerSet.Cards[2].Color {
		return true
	}

	return false
}

// 是否為順子
func (pokerSet *Poker) isStraight() bool {
	tempCard := new(Poker)
	tempCard.SetCards(pokerSet.CardGroup)
	tempCard.Sort(true)

	if tempCard.Cards[0].CardNumber+1 == tempCard.Cards[1].CardNumber && tempCard.Cards[1].CardNumber+1 == tempCard.Cards[2].CardNumber {
		return true
	}

	return false
}

// 是否為三條
func (pokerSet *Poker) isThreeOfAKind() bool {
	if pokerSet.Cards[0].CardNumber == pokerSet.Cards[1].CardNumber && pokerSet.Cards[1].CardNumber == pokerSet.Cards[2].CardNumber {
		return true
	}

	return false
}

// 是否為同花順
func (pokerSet *Poker) isStraightFlush() bool {
	if pokerSet.isFlush() && pokerSet.isStraight() {
		return true
	}

	return false
}

// 是否為皇家同花順
func (pokerSet *Poker) isRoyalFlush() bool {
	tempCard := new(Poker)
	tempCard.SetCards(pokerSet.CardGroup)
	tempCard.Sort(true)

	if pokerSet.isStraightFlush() && (tempCard.Cards[0].CardNumber == 12 && tempCard.Cards[1].CardNumber == 13 && tempCard.Cards[2].CardNumber == 14) {
		return true
	}

	return false
}

/*-------------Poker Manager--------------*/

// PatternString 牌型字串轉換
func PatternString(pattern int) string {
	if pattern < PatternNone && pattern >= PatternMax {
		return ""
	}

	return StrPattern[pattern]
}
