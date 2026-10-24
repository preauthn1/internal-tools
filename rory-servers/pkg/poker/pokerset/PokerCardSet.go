package pokerset

import (
	"math/rand"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// PokerCardSet ?
type PokerCardSet struct {
	aryPCards []*PokerCard // 牌組
}

// NewPokerCardSet ?
func NewPokerCardSet() *PokerCardSet {
	pPoker := new(PokerCardSet)
	pPoker.aryPCards = make([]*PokerCard, 0)

	return pPoker
}

// NewSet 建立新牌組，傳入參數［幾副撲克牌］、［鬼牌張數］
func (pPoker *PokerCardSet) NewSet(i32SetCount, i32JokerCount int32) {
	if (i32SetCount <= 0 || i32SetCount > POKERMaxSet) || (i32JokerCount < 0 || i32JokerCount > POKERJokerCardsCount) {
		mlog.Debug("牌組設定錯誤，SetCount[%d], JokerCount[%d]", i32SetCount, i32JokerCount)
		return
	}

	pPoker.aryPCards = nil
	pPoker.aryPCards = make([]*PokerCard, 0)

	for i := int32(0); i < i32SetCount; i++ {
		for i32Color := POKERCardColorClub; i32Color <= POKERCardColorSpade; i32Color++ {
			for i32No := POKERCardNo1; i32No < POKERCardNoA; i32No++ {
				pCard := new(PokerCard)
				pCard.I32Color = i32Color
				pCard.I32Number = i32No

				pPoker.aryPCards = append(pPoker.aryPCards, pCard)
			}
		}

		// 鬼牌
		if i32JokerCount == 1 || i32JokerCount == 2 {
			pCard := new(PokerCard)
			pCard.I32Color = POKERCardColorOther
			pCard.I32Number = POKERCardNoJoker1

			pPoker.aryPCards = append(pPoker.aryPCards, pCard)
		}

		if i32JokerCount == 2 {
			pCard := new(PokerCard)
			pCard.I32Color = POKERCardColorOther
			pCard.I32Number = POKERCardNoJoker2

			pPoker.aryPCards = append(pPoker.aryPCards, pCard)
		}
	}
}

// Shuffle 洗牌
func (pPoker *PokerCardSet) Shuffle() {
	rand.Seed(time.Now().UnixNano()) // 設置種子

	for iCount := 0; iCount < 5; iCount++ {
		rand.Shuffle(len(pPoker.aryPCards), func(i, j int) { // 調用算法
			pPoker.aryPCards[i], pPoker.aryPCards[j] = pPoker.aryPCards[j], pPoker.aryPCards[i]
		})
	}
}

// AddSingleCard 加入一張牌
func (pPoker *PokerCardSet) AddSingleCard(strCodeString string) {
	if len(strCodeString) != LENSingleCardCode {
		return
	}

	if _, ok := gMapI32ColorByCode[strCodeString[0:1]]; !ok {
		return
	}
	if _, ok := gMapI32NumberByCode[strCodeString[1:]]; !ok {
		return
	}

	pCard := new(PokerCard)
	pCard.I32Color = gMapI32ColorByCode[strCodeString[0:1]]
	pCard.I32Number = gMapI32NumberByCode[strCodeString[1:]]

	pPoker.aryPCards = append(pPoker.aryPCards, pCard)
}

// AddCardSet 加入一副牌組
func (pPoker *PokerCardSet) AddCardSet(pCardSet *PokerCardSet) {
	if pCardSet == nil {
		return
	}

	for _, pCard := range pCardSet.aryPCards {
		if pCard != nil {
			pPoker.AddSingleCard(pCard.CodeString())
		}
	}
}

// FindCard 尋找一張牌
func (pPoker *PokerCardSet) FindCard(pCard *PokerCard) bool {
	if pCard == nil {
		return false
	}

	bGet := false
	for _, pData := range pPoker.aryPCards {
		if pData != nil {
			if pData.I32Color == pCard.I32Color && pData.I32Number == pCard.I32Number {
				bGet = true
				break
			}
		}
	}

	return bGet
}

// RemoveSingleCard 移除一張牌
func (pPoker *PokerCardSet) RemoveSingleCard(pCard *PokerCard) {
	if pCard == nil {
		return
	}

	pTemp := make([]*PokerCard, 0)
	bGet := false
	for _, pData := range pPoker.aryPCards {
		if pData != nil {
			if !bGet && (pData.I32Color == pCard.I32Color && pData.I32Number == pCard.I32Number) {
				bGet = true
				continue
			}

			pTemp = append(pTemp, pData)
		}
	}

	pPoker.aryPCards = nil
	pPoker.aryPCards = pTemp
}

// RemoveCardSet 移除牌組中的牌
func (pPoker *PokerCardSet) RemoveCardSet(pCardSet *PokerCardSet) {
	if pCardSet == nil {
		return
	}

	if pCardSet.CardCount() > pPoker.CardCount() {
		return
	}

	for _, pCard := range pCardSet.aryPCards {
		if pCard != nil {
			pPoker.RemoveSingleCard(pCard)
		}
	}
}

// CardCount 牌組張數
func (pPoker *PokerCardSet) CardCount() int {
	return len(pPoker.aryPCards)
}

// 牌組文字字串
func (pPoker *PokerCardSet) String() string {
	strCards := ""
	for i := 0; i < len(pPoker.aryPCards); i++ {
		str := ""
		if pPoker.aryPCards[i] != nil {
			str = pPoker.aryPCards[i].String()

			if i < len(pPoker.aryPCards)-1 {
				str += ", "
			}
		}

		strCards += str
	}

	return strCards
}

// CodeString 牌組代碼字串
func (pPoker *PokerCardSet) CodeString() string {
	strCards := ""
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i] != nil {
			strCards += pPoker.aryPCards[i].CodeString()
		}
	}

	return strCards
}

// SetCardSet 設置牌組
func (pPoker *PokerCardSet) SetCardSet(strCodeString string) {
	if strCodeString == "" {
		pPoker.aryPCards = nil
		pPoker.aryPCards = make([]*PokerCard, 0)
		return
	}

	pNewSet := CardCode2PokerCardSet(strCodeString)
	if pNewSet != nil {
		pPoker.aryPCards = nil
		pPoker.aryPCards = pNewSet.aryPCards
	}
}

// Deal 發牌
func (pPoker *PokerCardSet) Deal(iCardCount int) *PokerCardSet {
	if iCardCount <= 0 || iCardCount > pPoker.CardCount() {
		return NewPokerCardSet()
	}

	strCodeString := pPoker.CodeString()[0 : LENSingleCardCode*iCardCount]
	pRetSet := pPoker.DrawCards(strCodeString)

	return pRetSet
}

// DrawCards 將 strCodeString 中的牌從牌組中抽出
func (pPoker *PokerCardSet) DrawCards(strCodeString string) *PokerCardSet {
	if strCodeString == "" {
		return NewPokerCardSet()
	}

	pRetSet := CardCode2PokerCardSet(strCodeString)
	if pRetSet == nil {
		return NewPokerCardSet()
	}

	for _, pCard := range pRetSet.aryPCards {
		if !pPoker.FindCard(pCard) {
			return NewPokerCardSet()
		}
	}

	pPoker.RemoveCardSet(pRetSet)

	return pRetSet
}

// Sort 排序
// bByColor = true：依花色排序，bByColor = false：依數字排序
// bOneAsAce = true：數字 1 當作 Ace，bOneAsAce = false：數字 1 不變
// bAsc = true：由小到大排序，bAsc = false：由大到小排序
func (pPoker *PokerCardSet) Sort(bByColor, bOneAsAce, bAsc bool) {
	iLen := int(pPoker.CardCount())

	pSort := new(PokerCardSet)
	pSort.SetCardSet(pPoker.CodeString())
	if bOneAsAce {
		for i := 0; i < len(pSort.aryPCards); i++ {
			if pSort.aryPCards[i].I32Number == POKERCardNo1 {
				pSort.aryPCards[i].I32Number = POKERCardNoA
			}
		}
	}

	if bByColor {
		for i := 0; i < iLen; i++ {
			for j := i; j < iLen; j++ {
				i32Card1 := pSort.aryPCards[i].I32Color*100 + pSort.aryPCards[i].I32Number
				i32Card2 := pSort.aryPCards[j].I32Color*100 + pSort.aryPCards[j].I32Number
				if bAsc {
					if i32Card1 > i32Card2 {
						pSort.aryPCards[i], pSort.aryPCards[j] = pSort.aryPCards[j], pSort.aryPCards[i]
						pPoker.aryPCards[i], pPoker.aryPCards[j] = pPoker.aryPCards[j], pPoker.aryPCards[i]
					}
				} else {
					if i32Card1 < i32Card2 {
						pSort.aryPCards[i], pSort.aryPCards[j] = pSort.aryPCards[j], pSort.aryPCards[i]
						pPoker.aryPCards[i], pPoker.aryPCards[j] = pPoker.aryPCards[j], pPoker.aryPCards[i]
					}
				}
			}
		}
	} else {
		// 依數字排序
		for i := 0; i < iLen; i++ {
			for j := i; j < iLen; j++ {
				i32Card1 := pSort.aryPCards[i].I32Number*10 + pSort.aryPCards[i].I32Color
				i32Card2 := pSort.aryPCards[j].I32Number*10 + pSort.aryPCards[j].I32Color
				if bAsc {
					if i32Card1 > i32Card2 {
						pSort.aryPCards[i], pSort.aryPCards[j] = pSort.aryPCards[j], pSort.aryPCards[i]
						pPoker.aryPCards[i], pPoker.aryPCards[j] = pPoker.aryPCards[j], pPoker.aryPCards[i]
					}
				} else {
					if i32Card1 < i32Card2 {
						pSort.aryPCards[i], pSort.aryPCards[j] = pSort.aryPCards[j], pSort.aryPCards[i]
						pPoker.aryPCards[i], pPoker.aryPCards[j] = pPoker.aryPCards[j], pPoker.aryPCards[i]
					}
				}
			}
		}
	}
}

// Pattern 牌組牌型
func (pPoker *PokerCardSet) Pattern() int32 {
	if pPoker.CardCount() <= 0 {
		return POKERPatternNone
	}

	//POKER_Pattern_DragonFlush						// 同花一條龍
	// todo

	//POKER_Pattern_Dragon							// 一條龍
	// todo

	// 皇家同花順
	if pPoker.isRoyalFlush() {
		return POKERPatternRoyalFlush
	}

	//POKER_Pattern_StraightFlush						// 同花順
	if pPoker.isStraightFlush() {
		return POKERPatternStraightFlush
	}

	// 四條
	if pPoker.isFourOfAKind() {
		return POKERPatternFourOfAKind
	}

	// 葫蘆
	if pPoker.isFullHouse() {
		return POKERPatternFullHouse
	}

	// 同花
	if pPoker.isFlush() {
		return POKERPatternFlush
	}

	// 順子
	if pPoker.isStraight() {
		return POKERPatternStraight
	}

	// 三條
	if pPoker.isThreeOfAKind() {
		return POKERPatternThreeOfAKind
	}

	// 兩對
	if pPoker.isTwoPairs() {
		return POKERPatternTwoPairs
	}

	// 對子
	if pPoker.isPair() {
		return POKERPatternPair
	}

	return POKERPatternHighCard
}

// 是否為對子
func (pPoker *PokerCardSet) isPair() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternPair] {
		return false
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	// 只能有一個對子
	iPairCount := 0
	for _, iCount := range mapICardCount {
		if iCount > gMapIPatternCardCount[POKERPatternPair] {
			return false
		} else if iCount == gMapIPatternCardCount[POKERPatternPair] {
			iPairCount++
		}
	}

	if iPairCount != 1 {
		return false
	}

	return true
}

// 是否為兩對
func (pPoker *PokerCardSet) isTwoPairs() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternTwoPairs] {
		return false
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	// 只能有兩個不同對子
	iPairCount := 0
	for _, iCount := range mapICardCount {
		if iCount > gMapIPatternCardCount[POKERPatternPair] {
			return false
		} else if iCount == gMapIPatternCardCount[POKERPatternPair] {
			iPairCount++
		}
	}

	if iPairCount != 2 {
		return false
	}

	return true
}

// 是否為三條
func (pPoker *PokerCardSet) isThreeOfAKind() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternThreeOfAKind] {
		return false
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	iThreeCount := 0
	for _, iCount := range mapICardCount {
		if iCount > gMapIPatternCardCount[POKERPatternThreeOfAKind] || iCount == gMapIPatternCardCount[POKERPatternPair] {
			return false
		} else if iCount == gMapIPatternCardCount[POKERPatternThreeOfAKind] {
			iThreeCount++
		}
	}

	// 只能有一個三條
	if iThreeCount != 1 {
		return false
	}

	return true
}

// 是否為順子
func (pPoker *PokerCardSet) isStraight() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternStraight] {
		return false
	}

	iCheckCount := gMapIPatternCardCount[POKERPatternStraight]
	if pPoker.CardCount() >= POKERPatternBasicCardCount {
		iCheckCount = POKERPatternBasicCardCount
	}

	if !pPoker.HasNumberInOrder(iCheckCount) {
		return false
	}

	return true
}

// 是否為同花
func (pPoker *PokerCardSet) isFlush() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFlush] {
		return false
	}

	i32Color := pPoker.aryPCards[0].I32Color
	for i := 1; i < pPoker.CardCount(); i++ {
		if pPoker.aryPCards[i].I32Color != i32Color {
			return false
		}
	}

	iCheckCount := gMapIPatternCardCount[POKERPatternFlush]
	if pPoker.CardCount() >= POKERPatternBasicCardCount {
		iCheckCount = POKERPatternBasicCardCount
	}

	if pPoker.HasNumberInOrder(iCheckCount) {
		return false
	}

	return true
}

// 是否為葫蘆
func (pPoker *PokerCardSet) isFullHouse() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFullHouse] {
		return false
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	iTwoCount := 0
	iThreeCount := 0
	for _, iCount := range mapICardCount {
		if iCount == gMapIPatternCardCount[POKERPatternPair] {
			iTwoCount++
		} else if iCount == gMapIPatternCardCount[POKERPatternThreeOfAKind] {
			iThreeCount++
		}
	}

	if iTwoCount != 1 || iThreeCount != 1 {
		return false
	}

	return true
}

// 是否為四條
func (pPoker *PokerCardSet) isFourOfAKind() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFourOfAKind] {
		return false
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	iFourCount := 0
	for _, iCount := range mapICardCount {
		if iCount == gMapIPatternCardCount[POKERPatternFourOfAKind] {
			iFourCount++
		}
	}

	// 只能有一個四條
	if iFourCount != 1 {
		return false
	}

	return true
}

// 是否為同花順
func (pPoker *PokerCardSet) isStraightFlush() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternStraightFlush] {
		return false
	}

	if !pPoker.IsFlushType() {
		return false
	}

	iCheckCount := gMapIPatternCardCount[POKERPatternStraightFlush]
	if pPoker.CardCount() >= POKERPatternBasicCardCount {
		iCheckCount = POKERPatternBasicCardCount
	}

	if !pPoker.HasNumberInOrder(iCheckCount) {
		return false
	}

	return true
}

// 是否為皇家同花順
func (pPoker *PokerCardSet) isRoyalFlush() bool {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternRoyalFlush] {
		return false
	}

	aryPCards := make([]*PokerCardSet, POKERColorsCount)
	for i := 0; i < len(aryPCards); i++ {
		aryPCards[i] = new(PokerCardSet)
	}

	pTemp := new(PokerCardSet)
	pTemp.SetCardSet(pPoker.CodeString())
	pTemp.Sort(false, true, true)

	for i := 0; i < pTemp.CardCount(); i++ {
		if pTemp.aryPCards[i].I32Color >= POKERCardColorClub && pTemp.aryPCards[i].I32Color <= POKERCardColorSpade {
			if pTemp.aryPCards[i].I32Number == POKERCardNo1 {
				pTemp.aryPCards[i].I32Number = POKERCardNoA
			}
			aryPCards[pTemp.aryPCards[i].I32Color].AddSingleCard(pTemp.aryPCards[i].CodeString())
		}
	}

	for i := POKERCardColorClub; i <= POKERCardColorSpade; i++ {
		iGetCard := 0
		for j := 0; j < aryPCards[i].CardCount(); j++ {
			if aryPCards[i].aryPCards[j].I32Number >= POKERCardNo10 && aryPCards[i].aryPCards[j].I32Number <= POKERCardNoA {
				iGetCard |= 1 << uint(aryPCards[i].aryPCards[j].I32Number-POKERCardNo10)
			}
		}

		if iGetCard == 31 {
			return true
		}
	}

	return false
}

//POKER_Pattern_Dragon							// 是否為一條龍
// todo

//POKER_Pattern_DragonFlush						// 是否為同花一條龍
// todo

// HasNumberInOrder 檢查是否為順子
// iCheckCount：要檢查是否為順子的牌組張數，例如 7 張牌裡要檢查是否有 5 張的順面，則 iCheckCount = 5
func (pPoker *PokerCardSet) HasNumberInOrder(iCheckCount int) bool {
	if iCheckCount < POKERMinStraightCount || iCheckCount > pPoker.CardCount() {
		return false
	}

	aryBGetCard := make([]bool, POKEROneColorCardsCount)
	bGet := true

	pTemp := new(PokerCardSet)
	pTemp.SetCardSet(pPoker.CodeString())
	// 數字1
	for i := 0; i < len(pTemp.aryPCards); i++ {
		if pTemp.aryPCards[i].I32Number >= POKERCardNo1 && pTemp.aryPCards[i].I32Number <= POKERCardNoK {
			aryBGetCard[pTemp.aryPCards[i].I32Number-POKERCardNo1] = true
		}
	}

	for i := 0; i <= len(aryBGetCard)-iCheckCount; i++ {
		bGet = true
		for j := i; j < i+iCheckCount; j++ {
			if !aryBGetCard[j] {
				bGet = false
				break
			}
		}

		if bGet {
			break
		}
	}

	// A
	if !bGet {
		aryBGetCard = make([]bool, POKEROneColorCardsCount)
		for i := 0; i < len(pTemp.aryPCards); i++ {
			if pTemp.aryPCards[i].I32Number == POKERCardNo1 {
				pTemp.aryPCards[i].I32Number = POKERCardNoA
			}

			if pTemp.aryPCards[i].I32Number >= POKERCardNo2 && pTemp.aryPCards[i].I32Number <= POKERCardNoA {
				aryBGetCard[pTemp.aryPCards[i].I32Number-POKERCardNo2] = true
			}
		}

		for i := 0; i <= len(aryBGetCard)-iCheckCount; i++ {
			bGet = true
			for j := i; j < i+iCheckCount; j++ {
				if !aryBGetCard[j] {
					bGet = false
					break
				}
			}

			if bGet {
				break
			}
		}
	}

	return bGet
}

// 依牌型取得牌組，［i32Pattern］牌型，［iCardCount］0: 依該牌型預設張數
func (pPoker *PokerCardSet) getCardSetsByPattern(i32Pattern int32, iCardCount int) []*PokerCardSet {
	aryPCardSet := make([]*PokerCardSet, 0)

	b3Straight := true
	if iCardCount >= POKERPatternBasicCardCount {
		b3Straight = false
	}

	switch i32Pattern {
	case POKERPatternHighCard: // 高牌
		aryPCardSet = pPoker.getHighCard()
	case POKERPatternPair: // 對子
		aryPCardSet = pPoker.getPair()
	case POKERPatternTwoPairs: // 兩對
		aryPCardSet = pPoker.getTwoPairs()
	case POKERPatternThreeOfAKind: // 三條
		aryPCardSet = pPoker.getThreeOfAKind()
	case POKERPatternStraight: // 順子
		aryPCardSet = pPoker.getStraight(b3Straight)
	case POKERPatternFlush: // 同花
		aryPCardSet = pPoker.getFlush(iCardCount)
	case POKERPatternFullHouse: // 葫蘆
		aryPCardSet = pPoker.getFullHouse()
	case POKERPatternFourOfAKind: // 四條
		aryPCardSet = pPoker.getFourOfAKind()
	case POKERPatternStraightFlush: // 同花順
		aryPCardSet = pPoker.getStraightFlush(b3Straight)
	case POKERPatternRoyalFlush: // 皇家同花順
		aryPCardSet = pPoker.getRoyalFlush()
		//case POKER_Pattern_Dragon:        // 一條龍	todo
		//case POKER_Pattern_DragonFlush:   // 同花一條龍	todo
	}

	return aryPCardSet
}

// DealByPattern 依所需牌型發牌，iCardCount：牌組總張數
func (pPoker *PokerCardSet) DealByPattern(i32Pattern int32, iCardCount int) *PokerCardSet {
	if i32Pattern < POKERPatternHighCard || i32Pattern >= POKERPatternMax {
		return NewPokerCardSet()
	}

	if iCardCount > pPoker.CardCount() {
		return NewPokerCardSet()
	}

	aryPCardSet := pPoker.getCardSetsByPattern(i32Pattern, iCardCount)

	// 若牌盒中已取不出所要的牌型就發高牌
	if len(aryPCardSet) <= 0 {
		aryPCardSet = pPoker.getHighCard()
	}

	pRetSet := new(PokerCardSet)

	// 隨機取得牌組
	iIdx := util.RndInt(0, len(aryPCardSet))
	pRetSet.SetCardSet(aryPCardSet[iIdx].CodeString())
	if pRetSet.CardCount() <= 0 {
		return NewPokerCardSet()
	}

	// 補足張數
	pSet := pPoker.supplementCards(pRetSet, iCardCount-pRetSet.CardCount())

	pRetSet.AddCardSet(pSet)
	pPoker.RemoveCardSet(pRetSet)

	return pRetSet
}

// 補牌，［pBaseSet］原已發出的牌組，［iSupplementCount］要補的張數
// 回傳［補牌的牌組］
func (pPoker *PokerCardSet) supplementCards(pBaseSet *PokerCardSet, iSupplementCount int) *PokerCardSet {
	if pBaseSet == nil {
		return NewPokerCardSet()
	}

	if pBaseSet.CardCount() <= 0 || iSupplementCount <= 0 {
		return NewPokerCardSet()
	}

	mlog.Debug("pBaseSet(%d)%s [%s], iSupplementCount[%d]", pBaseSet.CardCount(), pBaseSet.CodeString(), PatternString(pBaseSet.Pattern()), iSupplementCount)

	pRetset := NewPokerCardSet()

	pTempSet := NewPokerCardSet()
	pCheckSet := NewPokerCardSet()
	i32BasePattern := pBaseSet.Pattern()
	for i := 0; i < pPoker.CardCount(); i++ {
		pCard := pPoker.GetSingleCard(i, false)
		pCheckSet.SetCardSet(pBaseSet.CodeString())
		pCheckSet.AddCardSet(pTempSet)

		if pCheckSet.IsNumberExisted(pCard.I32Number) {
			continue
		}

		pCheckSet.AddSingleCard(pCard.CodeString())
		if pCheckSet.Pattern() != i32BasePattern {
			continue
		}

		pTempSet.AddSingleCard(pPoker.aryPCards[i].CodeString())
		if pTempSet.CardCount() >= iSupplementCount {
			break
		}
	}
	mlog.Debug("pTempSet(%d)%s", pTempSet.CardCount(), pTempSet.CodeString())

	pRetset.SetCardSet(pTempSet.CodeString())

	return pRetset
}

// 取得高牌
func (pPoker *PokerCardSet) getHighCard() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternHighCard] {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	for i := 0; i < pPoker.CardCount()-gMapIPatternCardCount[POKERPatternHighCard]; i++ {
		pSet := new(PokerCardSet)
		pSet.SetCardSet(pPoker.aryPCards[i].CodeString())
		for j := i + 1; j < pPoker.CardCount(); j++ {
			if pSet.CardCount() >= gMapIPatternCardCount[POKERPatternHighCard] {
				if pSet.Pattern() == POKERPatternHighCard {
					aryPRetSet = append(aryPRetSet, pSet)
				}
				break
			}

			if pSet.IsNumberExisted(pPoker.aryPCards[j].I32Number) {
				continue
			}
			pSet.AddSingleCard(pPoker.aryPCards[j].CodeString())
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得對子
func (pPoker *PokerCardSet) getPair() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternPair] {
		return make([]*PokerCardSet, 0)
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	// 取出對子
	for i32No, iCount := range mapICardCount {
		if iCount >= gMapIPatternCardCount[POKERPatternPair] {
			for i := 0; i < len(pPoker.aryPCards); i++ {
				if pPoker.aryPCards[i].I32Number == i32No {
					pSet := new(PokerCardSet)
					for j := i; j < len(pPoker.aryPCards); j++ {
						if pPoker.aryPCards[j].I32Number == i32No {
							pSet.AddSingleCard(pPoker.aryPCards[j].CodeString())
						}

						if pSet.CardCount() >= gMapIPatternCardCount[POKERPatternPair] {
							pSet.Sort(true, false, true)
							aryPRetSet = append(aryPRetSet, pSet)
							break
						}
					}
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得兩對
func (pPoker *PokerCardSet) getTwoPairs() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternTwoPairs] {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)
	aryPPairs := pPoker.getPair()
	if len(aryPPairs) < 2 {
		return make([]*PokerCardSet, 0)
	}
	for i := 0; i < len(aryPPairs)-1; i++ {
		for j := i + 1; j < len(aryPPairs); j++ {
			if aryPPairs[i].aryPCards[0].I32Number != aryPPairs[j].aryPCards[0].I32Number {
				pSet := new(PokerCardSet)
				pSet.AddCardSet(aryPPairs[i])
				pSet.AddCardSet(aryPPairs[j])
				pSet.Sort(false, false, true)
				aryPRetSet = append(aryPRetSet, pSet)
			}
		}
	}

	return aryPRetSet
}

// 取得三條
func (pPoker *PokerCardSet) getThreeOfAKind() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternThreeOfAKind] {
		return make([]*PokerCardSet, 0)
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	// 取出三條
	for i32No, iCount := range mapICardCount {
		if iCount >= gMapIPatternCardCount[POKERPatternThreeOfAKind] {
			pSet := new(PokerCardSet)
			for i := 0; i < len(pPoker.aryPCards); i++ {
				if pPoker.aryPCards[i].I32Number == i32No {
					pSet.AddSingleCard(pPoker.aryPCards[i].CodeString())
				}

				if pSet.CardCount() >= gMapIPatternCardCount[POKERPatternThreeOfAKind] {
					aryPRetSet = append(aryPRetSet, pSet)
					break
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得順子，［b3Straight］true: 3 張順子，false: 5 張順子
func (pPoker *PokerCardSet) getStraight(b3Straight bool) []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternStraight] {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	iCardCount := gMapIPatternCardCount[POKERPatternStraight]
	if !b3Straight {
		iCardCount = POKERPatternBasicCardCount
	}

	// 取得各數字的張數
	aryPCards := make([]*PokerCardSet, POKEROneColorCardsCount)
	for i := 0; i < len(aryPCards); i++ {
		aryPCards[i] = NewPokerCardSet()
	}

	for i := 0; i < pPoker.CardCount(); i++ {
		pCard := pPoker.GetSingleCard(i, false)
		if pCard.I32Number >= POKERCardNo1 && pCard.I32Number <= POKERCardNoK {
			aryPCards[pCard.I32Number-POKERCardNo1].AddSingleCard(pCard.CodeString())
		}
	}

	// 取出順子
	for i := 0; i <= len(aryPCards)-iCardCount+1; i++ {
		bStraight := true
		for j := i; j < i+iCardCount; j++ {
			if j < len(aryPCards) {
				if aryPCards[j].CardCount() <= 0 {
					bStraight = false
					break
				}
			} else {
				// 如果後面的牌是 QK 或是 10JQK，則判斷 A（數字 1）的張數
				if aryPCards[0].CardCount() <= 0 {
					bStraight = false
					break
				}
			}
		}

		if bStraight {
			for j := 0; j < int(POKERColorsCount); j++ {
				pSet := NewPokerCardSet()
				for k := i; k < i+iCardCount; k++ {
					iIdx := 0
					if k < len(aryPCards) {
						if j < aryPCards[k].CardCount() {
							iIdx = j
						} else {
							iIdx = aryPCards[k].CardCount() - 1
						}
						pSet.AddSingleCard(aryPCards[k].aryPCards[iIdx].CodeString())
					} else {
						// 如果後面的牌是 QK 或是 10JQK，則把 A（數字 1） 加進來
						if j < aryPCards[0].CardCount() {
							iIdx = j
						} else {
							iIdx = aryPCards[0].CardCount() - 1
						}
						pSet.AddSingleCard(aryPCards[0].aryPCards[iIdx].CodeString())
					}
				}
				pSet.Shuffle()
				if pSet.Pattern() == POKERPatternStraight {
					aryPRetSet = append(aryPRetSet, pSet)
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得同花
func (pPoker *PokerCardSet) getFlush(iCardCount int) []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFlush] {
		return make([]*PokerCardSet, 0)
	}

	if iCardCount < gMapIPatternCardCount[POKERPatternFlush] || iCardCount > pPoker.CardCount() {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	aryPCards := make([]*PokerCardSet, POKERColorsCount)
	for i := 0; i < len(aryPCards); i++ {
		aryPCards[i] = new(PokerCardSet)
	}

	for i := 0; i < pPoker.CardCount(); i++ {
		if pPoker.aryPCards[i].I32Color >= POKERCardColorClub && pPoker.aryPCards[i].I32Color <= POKERCardColorSpade {
			aryPCards[pPoker.aryPCards[i].I32Color].AddSingleCard(pPoker.aryPCards[i].CodeString())
		}
	}

	for i32Pattern := POKERCardColorClub; i32Pattern <= POKERCardColorSpade; i32Pattern++ {
		if aryPCards[i32Pattern].CardCount() >= iCardCount {
			for i := 0; i < pPoker.CardCount()-iCardCount; i++ {
				pSet := new(PokerCardSet)
				pSet.AddSingleCard(pPoker.aryPCards[i].CodeString())
				i32Color := pPoker.aryPCards[i].I32Color
				for j := i + 1; j < pPoker.CardCount(); j++ {
					if pSet.CardCount() >= iCardCount {
						if pSet.Pattern() == POKERPatternFlush {
							aryPRetSet = append(aryPRetSet, pSet)
						}
						break
					}

					pCard := pPoker.GetSingleCard(j, false)
					if pCard.I32Color != i32Color || pSet.IsNumberExisted(pCard.I32Number) {
						continue
					}

					pSet.AddSingleCard(pCard.CodeString())
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得葫蘆
func (pPoker *PokerCardSet) getFullHouse() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFullHouse] {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	aryPSet3 := pPoker.getThreeOfAKind()
	aryPSet2 := pPoker.getPair()

	if len(aryPSet3) <= 0 || len(aryPSet2) <= 0 {
		return make([]*PokerCardSet, 0)
	}
	for i := 0; i < len(aryPSet3); i++ {
		for j := 0; j < len(aryPSet2); j++ {
			if aryPSet3[i].aryPCards[0].I32Number != aryPSet2[j].aryPCards[0].I32Number {
				pSet := new(PokerCardSet)
				pSet.AddCardSet(aryPSet3[i])
				pSet.AddCardSet(aryPSet2[j])

				aryPRetSet = append(aryPRetSet, pSet)
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得四條
func (pPoker *PokerCardSet) getFourOfAKind() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternFourOfAKind] {
		return make([]*PokerCardSet, 0)
	}

	// 統計各數字張數
	mapICardCount := make(map[int32]int)
	for i := 0; i < len(pPoker.aryPCards); i++ {
		if pPoker.aryPCards[i].I32Number >= POKERCardNo1 && pPoker.aryPCards[i].I32Number <= POKERCardNoK {
			mapICardCount[pPoker.aryPCards[i].I32Number]++
		}
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	// 取出四條
	for i32No, iCount := range mapICardCount {
		if iCount >= gMapIPatternCardCount[POKERPatternFourOfAKind] {
			pSet := new(PokerCardSet)
			for i := 0; i < len(pPoker.aryPCards); i++ {
				if pPoker.aryPCards[i].I32Number == i32No {
					pSet.AddSingleCard(pPoker.aryPCards[i].CodeString())
				}

				if pSet.CardCount() >= gMapIPatternCardCount[POKERPatternFourOfAKind] {
					aryPRetSet = append(aryPRetSet, pSet)
					break
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得同花順
func (pPoker *PokerCardSet) getStraightFlush(b3Straight bool) []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternStraightFlush] {
		return make([]*PokerCardSet, 0)
	}

	iCardCount := gMapIPatternCardCount[POKERPatternStraight]
	if !b3Straight {
		iCardCount = POKERPatternBasicCardCount
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	aryPCards := make([]*PokerCardSet, POKERColorsCount)
	for i := 0; i < len(aryPCards); i++ {
		aryPCards[i] = new(PokerCardSet)
	}

	pTemp := new(PokerCardSet)
	pTemp.SetCardSet(pPoker.CodeString())

	for i := 0; i < pTemp.CardCount(); i++ {
		pCard := pTemp.aryPCards[i]
		if pCard.I32Color >= POKERCardColorClub && pCard.I32Color <= POKERCardColorSpade {
			aryPCards[pCard.I32Color].AddSingleCard(pCard.CodeString())
		}
	}

	for i := POKERCardColorClub; i <= POKERCardColorSpade; i++ {
		aryPCards[i].Sort(false, false, true)

		if aryPCards[i].CardCount() >= iCardCount {
			for a := 0; a <= aryPCards[i].CardCount()-iCardCount; a++ {
				bStraight := true
				i32No := aryPCards[i].aryPCards[a].I32Number
				for b := a; b < a+iCardCount; b++ {
					if aryPCards[i].aryPCards[b].I32Number != i32No {
						bStraight = false
						break
					}

					i32No++
				}

				if bStraight {
					pSet := new(PokerCardSet)
					for b := a; b < a+iCardCount; b++ {
						pSet.AddSingleCard(aryPCards[i].aryPCards[b].CodeString())
					}
					aryPRetSet = append(aryPRetSet, pSet)
				}
			}
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

// 取得皇家同花順
func (pPoker *PokerCardSet) getRoyalFlush() []*PokerCardSet {
	if pPoker.CardCount() < gMapIPatternCardCount[POKERPatternRoyalFlush] {
		return make([]*PokerCardSet, 0)
	}

	aryPRetSet := make([]*PokerCardSet, 0)

	aryPCards := make([]*PokerCardSet, POKERColorsCount)
	for i := 0; i < len(aryPCards); i++ {
		aryPCards[i] = new(PokerCardSet)
	}

	pTemp := new(PokerCardSet)
	pTemp.SetCardSet(pPoker.CodeString())
	pTemp.Sort(false, true, true)

	for i := 0; i < pTemp.CardCount(); i++ {
		if pTemp.aryPCards[i].I32Color >= POKERCardColorClub && pTemp.aryPCards[i].I32Color <= POKERCardColorSpade {
			aryPCards[pTemp.aryPCards[i].I32Color].AddSingleCard(pTemp.aryPCards[i].CodeString())
		}
	}

	for i := POKERCardColorClub; i <= POKERCardColorSpade; i++ {
		mapBCheck := map[int32]bool{
			POKERCardNo1:  false,
			POKERCardNo10: false,
			POKERCardNoJ:  false,
			POKERCardNoQ:  false,
			POKERCardNoK:  false,
		}

		pSet := new(PokerCardSet)
		for j := 0; j < aryPCards[i].CardCount(); j++ {
			if bValue, ok := mapBCheck[aryPCards[i].aryPCards[j].I32Number]; ok {
				if !bValue {
					pSet.AddSingleCard(aryPCards[i].aryPCards[j].CodeString())
					mapBCheck[aryPCards[i].aryPCards[j].I32Number] = true
				}
			}
		}

		if pSet.CardCount() >= gMapIPatternCardCount[POKERPatternRoyalFlush] {
			aryPRetSet = append(aryPRetSet, pSet)
		}
	}

	if len(aryPRetSet) <= 0 {
		return make([]*PokerCardSet, 0)
	}

	return aryPRetSet
}

//POKER_Pattern_Dragon							// 取得一條龍
// todo

//POKER_Pattern_DragonFlush						// 取得同花一條龍
// todo

// Compare 比牌
func (pPoker *PokerCardSet) Compare(pCardSet1, pCardSet2 *PokerCardSet, bCmpColor bool, bCmpSingle bool) int {
	if pCardSet1 == nil || pCardSet2 == nil {
		return POKERCompareEven
	}

	if pCardSet1.CardCount() <= 0 || pCardSet2.CardCount() <= 0 {
		return POKERCompareEven
	}

	if pCardSet1.CardCount() != pCardSet2.CardCount() {
		return POKERCompareEven
	}

	// 先比牌型
	if pCardSet1.Pattern() > pCardSet2.Pattern() {
		return POKERCompareDealer
	} else if pCardSet1.Pattern() < pCardSet2.Pattern() {
		return POKERComparePlayer
	}

	// 牌型若相同則比數字
	aryPSet1 := pCardSet1.getCardSetsByPattern(pCardSet1.Pattern(), pCardSet1.CardCount())
	aryPSet2 := pCardSet2.getCardSetsByPattern(pCardSet2.Pattern(), pCardSet2.CardCount())
	if len(aryPSet1) <= 0 || len(aryPSet2) <= 0 {
		return POKERCompareEven
	}

	pTempSet1, pTempSet2 := new(PokerCardSet), new(PokerCardSet)
	pTempSet1.SetCardSet(pCardSet1.CodeString())
	pTempSet2.SetCardSet(pCardSet2.CodeString())

	switch pCardSet1.Pattern() {
	case POKERPatternHighCard: // 高牌
		fallthrough
	case POKERPatternFlush: // 同花
		pSet1, pSet2 := new(PokerCardSet), new(PokerCardSet)
		pSet1.SetCardSet(pCardSet1.CodeString())
		pSet1.Sort(false, true, false)
		pSet2.SetCardSet(pCardSet2.CodeString())
		pSet2.Sort(false, true, false)

		for i := 0; i < pSet1.CardCount(); i++ {
			pCard1 := pSet1.GetSingleCard(i, true)
			pCard2 := pSet2.GetSingleCard(i, true)
			if pCard1.I32Number > pCard2.I32Number {
				return POKERCompareDealer
			} else if pCard1.I32Number < pCard2.I32Number {
				return POKERComparePlayer
			}
		}
	case POKERPatternPair: // 對子
		pSet1 := pTempSet1.DrawCards(aryPSet1[0].CodeString())
		pSet2 := pTempSet2.DrawCards(aryPSet2[0].CodeString())

		// 先比對子數字
		pCard1 := pSet1.GetSingleCard(0, true)
		pCard2 := pSet2.GetSingleCard(0, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		} else if !bCmpSingle && pCard1.I32Number == pCard2.I32Number {
			return POKERCompareEven
		}

		// 對子數字相同，比剩下的牌
		pTempSet1.Sort(false, true, false)
		pTempSet2.Sort(false, true, false)

		for i := 0; i < pTempSet1.CardCount(); i++ {
			pCard1 := pTempSet1.GetSingleCard(i, true)
			pCard2 := pTempSet2.GetSingleCard(i, true)
			if pCard1.I32Number > pCard2.I32Number {
				return POKERCompareDealer
			} else if pCard1.I32Number < pCard2.I32Number {
				return POKERComparePlayer
			}
		}
	case POKERPatternTwoPairs: // 兩對
		pSet1 := pTempSet1.DrawCards(aryPSet1[0].CodeString())
		pSet2 := pTempSet2.DrawCards(aryPSet2[0].CodeString())

		pSet1.Sort(false, true, false)
		pSet2.Sort(false, true, false)

		// 先比第一組對子數字
		pCard1 := pSet1.GetSingleCard(0, true)
		pCard2 := pSet2.GetSingleCard(0, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		}

		// 再比第二組對子數字
		pCard1 = pSet1.GetSingleCard(pSet1.CardCount()-1, true)
		pCard2 = pSet2.GetSingleCard(pSet2.CardCount()-1, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		} else if !bCmpSingle && pCard1.I32Number == pCard2.I32Number {
			return POKERCompareEven
		}

		// 對子數字相同，比剩下的牌
		pTempSet1.Sort(false, true, false)
		pTempSet2.Sort(false, true, false)

		for i := 0; i < pTempSet1.CardCount(); i++ {
			pCard1 := pTempSet1.GetSingleCard(i, true)
			pCard2 := pTempSet2.GetSingleCard(i, true)
			if pCard1.I32Number > pCard2.I32Number {
				return POKERCompareDealer
			} else if pCard1.I32Number < pCard2.I32Number {
				return POKERComparePlayer
			}
		}
	case POKERPatternThreeOfAKind: // 三條
		fallthrough
	case POKERPatternFourOfAKind: // 四條
		pSet1 := pTempSet1.DrawCards(aryPSet1[0].CodeString())
		pSet2 := pTempSet2.DrawCards(aryPSet2[0].CodeString())

		// 先比三條/四條數字
		pCard1 := pSet1.GetSingleCard(0, true)
		pCard2 := pSet2.GetSingleCard(0, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		}

		// 三條/四條數字相同，比剩下的牌
		pTempSet1.Sort(false, true, false)
		pTempSet2.Sort(false, true, false)

		for i := 0; i < pTempSet1.CardCount(); i++ {
			pCard1 := pTempSet1.GetSingleCard(i, true)
			pCard2 := pTempSet2.GetSingleCard(i, true)
			if pCard1.I32Number > pCard2.I32Number {
				return POKERCompareDealer
			} else if pCard1.I32Number < pCard2.I32Number {
				return POKERComparePlayer
			}
		}
	case POKERPatternStraight: // 順子
		fallthrough
	case POKERPatternStraightFlush: // 同花順
		pSet1 := pTempSet1.DrawCards(aryPSet1[0].CodeString())
		pSet2 := pTempSet2.DrawCards(aryPSet2[0].CodeString())

		// 比最大的數字
		bHasA1, bIs15Straight1 := pPoker.CheckAStraightType(pSet1)
		pCard1, pCard2 := new(PokerCard), new(PokerCard)
		if !bHasA1 {
			pSet1.Sort(false, true, false)
			pCard1 = pSet1.GetSingleCard(0, true)
		} else {
			if bIs15Straight1 {
				pSet1.Sort(false, false, false)
				pCard1 = pSet1.GetSingleCard(0, false)
			} else {
				pSet1.Sort(false, true, false)
				pCard1 = pSet1.GetSingleCard(0, true)
			}
		}
		bHasA2, bIs15Straight2 := pPoker.CheckAStraightType(pSet2)
		if !bHasA2 {
			pSet2.Sort(false, true, false)
			pCard2 = pSet2.GetSingleCard(0, true)
		} else {
			if bIs15Straight2 {
				pSet2.Sort(false, false, false)
				pCard2 = pSet2.GetSingleCard(0, false)
			} else {
				pSet2.Sort(false, true, false)
				pCard2 = pSet2.GetSingleCard(0, true)
			}
		}

		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		}
	case POKERPatternFullHouse: // 葫蘆
		// 先取出三條的部分比
		aryPThree1 := pTempSet1.getCardSetsByPattern(POKERPatternThreeOfAKind, 0)
		aryPThree2 := pTempSet2.getCardSetsByPattern(POKERPatternThreeOfAKind, 0)
		if len(aryPThree1) <= 0 || len(aryPThree2) <= 0 {
			return POKERCompareEven
		}

		pSet1 := pTempSet1.DrawCards(aryPThree1[0].CodeString())
		pSet2 := pTempSet2.DrawCards(aryPThree2[0].CodeString())

		pCard1 := pSet1.GetSingleCard(0, true)
		pCard2 := pSet2.GetSingleCard(0, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		}

		// 後比對子的部分
		pCard1 = pTempSet1.GetSingleCard(0, true)
		pCard2 = pTempSet2.GetSingleCard(0, true)
		if pCard1.I32Number > pCard2.I32Number {
			return POKERCompareDealer
		} else if pCard1.I32Number < pCard2.I32Number {
			return POKERComparePlayer
		}
	case POKERPatternRoyalFlush: // 皇家同花順
		return POKERCompareEven
	}

	return POKERCompareEven
}

// CheckAStraightType 檢查是哪種順子，回傳值：［是否有 A］、［是否為 1-5 順子］
func (pPoker *PokerCardSet) CheckAStraightType(pStraightSet *PokerCardSet) (bool, bool) {
	if pStraightSet == nil {
		return false, false
	}

	if pStraightSet.CardCount() != gMapIPatternCardCount[POKERPatternStraight] &&
		pStraightSet.CardCount() != gMapIPatternCardCount[POKERPatternStraightFlush] {
		return false, false
	}

	// 檢查是否為順子
	if pStraightSet.Pattern() != POKERPatternStraight && pStraightSet.Pattern() != POKERPatternStraightFlush {
		return false, false
	}

	bHasA, bIs15Straight := false, false

	// 檢查是否有 A
	for i := 0; i < pStraightSet.CardCount(); i++ {
		if pStraightSet.aryPCards[i].I32Number == POKERCardNo1 {
			bHasA = true
			break
		}
	}

	if !bHasA {
		return bHasA, bIs15Straight
	}

	// 檢查是否為 1-5 順子
	i15Straight := 0
	for i := 0; i < pStraightSet.CardCount(); i++ {
		if pStraightSet.aryPCards[i].I32Number >= POKERCardNo1 && pStraightSet.aryPCards[i].I32Number <= POKERCardNo5 {
			i15Straight |= 1 << uint(pStraightSet.aryPCards[i].I32Number-POKERCardNo1)
		}

		if i15Straight == 31 {
			bIs15Straight = true
			return bHasA, bIs15Straight
		}
	}

	return bHasA, bIs15Straight
}

// GetSingleCard 取得牌組中的一張牌
func (pPoker *PokerCardSet) GetSingleCard(iIndex int, bOneAsAce bool) *PokerCard {
	if iIndex < 0 || iIndex >= pPoker.CardCount() {
		return NewPokerCard(POKERCardColorOther, POKERCardNoHide)
	}

	pRetCard := NewPokerCard(pPoker.aryPCards[iIndex].I32Color, pPoker.aryPCards[iIndex].I32Number)
	if bOneAsAce && pRetCard.I32Number == POKERCardNo1 {
		pRetCard.I32Number = POKERCardNoA
	}

	return pRetCard
}

// GetPokerCardSet 取得 PokerCardSet
func (pPoker *PokerCardSet) GetPokerCardSet() *PokerCardSet {
	pRetSet := new(PokerCardSet)
	pRetSet.SetCardSet(pPoker.CodeString())

	return pRetSet
}

// IsFlushType 檢查是否為同花類型
func (pPoker *PokerCardSet) IsFlushType() bool {
	if pPoker.CardCount() <= 1 {
		return false
	}

	i32Color := pPoker.aryPCards[0].I32Color
	for i := 1; i < pPoker.CardCount(); i++ {
		if pPoker.aryPCards[i].I32Color != i32Color {
			return false
		}
	}

	return true
}

// IsNumberExisted 牌組內是否有此數字
func (pPoker *PokerCardSet) IsNumberExisted(i32CheckNo int32) bool {
	for i := 0; i < pPoker.CardCount(); i++ {
		if pPoker.aryPCards[i].I32Number == i32CheckNo {
			return true
		}
	}

	return false
}
