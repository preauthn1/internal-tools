package jlbpk

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/pokerset"
)

type JlbpkPokerSet struct {
	pPokerSet *pokerset.PokerCardSet
}

func NewJlbpkPokerSet() *JlbpkPokerSet {
	pPokerSet := new(JlbpkPokerSet)
	pPokerSet.pPokerSet = new(pokerset.PokerCardSet)

	return pPokerSet
}

// 建立新牌組，傳入參數［幾副撲克牌］、［鬼牌張數］
func (pPoker *JlbpkPokerSet) NewSet(i32SetCount, i32JokerCount int32) {
	pPoker.pPokerSet.NewSet(i32SetCount, i32JokerCount)
}

// 洗牌
func (pPoker *JlbpkPokerSet) Shuffle() {
	pPoker.pPokerSet.Shuffle()
}

// 加入一張牌
func (pPoker *JlbpkPokerSet) AddSingleCard(strCodeString string) {
	pPoker.pPokerSet.AddSingleCard(strCodeString)
}

// 加入一副牌組
func (pPoker *JlbpkPokerSet) AddCardSet(pCardSet *pokerset.PokerCardSet) {
	pPoker.pPokerSet.AddCardSet(pCardSet)
}

// 牌組張數
func (pPoker *JlbpkPokerSet) CardCount() int {
	return pPoker.pPokerSet.CardCount()
}

// 牌組文字字串
func (pPoker *JlbpkPokerSet) String() string {
	return pPoker.pPokerSet.String()
}

// 移除一張牌
func (pPoker *JlbpkPokerSet) RemoveSingleCard(pCard *pokerset.PokerCard) {
	pPoker.RemoveSingleCard(pCard)
}

// 移除牌組中的牌
func (pPoker *JlbpkPokerSet) RemoveCardSet(pCardSet *pokerset.PokerCardSet) {
	pPoker.pPokerSet.RemoveCardSet(pCardSet)
}

// 牌組代碼字串
func (pPoker *JlbpkPokerSet) CodeString() string {
	return pPoker.pPokerSet.CodeString()
}

// 設置牌組
func (pPoker *JlbpkPokerSet) SetCardSet(strCodeString string) {
	pPoker.pPokerSet.SetCardSet(strCodeString)
}

// 發牌
func (pPoker *JlbpkPokerSet) Deal(iCardCount int) *pokerset.PokerCardSet {
	return pPoker.pPokerSet.Deal(iCardCount)
}

// 將 strCodeString 中的牌從牌組中抽出
func (pPoker *JlbpkPokerSet) DrawCards(strCodeString string) *pokerset.PokerCardSet {
	return pPoker.pPokerSet.DrawCards(strCodeString)
}

// 取得牌組中的一張牌
func (pPoker *JlbpkPokerSet) GetSingleCard(iIndex int, bOneAsAce bool) *pokerset.PokerCard {
	return pPoker.pPokerSet.GetSingleCard(iIndex, bOneAsAce)
}

// 排序
// bByColor = true：依花色排序，bByColor = false：依數字排序
// bOneAsAce = true：數字 1 當作 Ace，bOneAsAce = false：數字 1 不變
// bAsc = true：由小到大排序，bAsc = false：由大到小排序
func (pPoker *JlbpkPokerSet) Sort(bByColor, bOneAsAce, bAsc bool) {
	pPoker.pPokerSet.Sort(bByColor, bOneAsAce, bAsc)
}

// 牌組牌型
func (pPoker *JlbpkPokerSet) Pattern(bIsDealer bool) int32 {
	iRetPattern := pokerset.POKERPatternNone

	if !bIsDealer {
		// 玩家
		iRetPattern = pPoker.pPokerSet.Pattern()
	} else {
		// 莊家
		bCheckedOK := true

		if pPoker.pPokerSet.Pattern() == pokerset.POKERPatternHighCard {
			bCheckedOK = pPoker.CheckIfHasAK()
		}

		if bCheckedOK {
			iRetPattern = pPoker.pPokerSet.Pattern()
		} else {
			iRetPattern = pokerset.POKERPatternNone
		}
	}

	return iRetPattern
}

// 依所需牌型發牌，iCardCount［牌組總張數］，bHasAK［若為高牌時是否必須有 A、K］
//func (pPoker *JlbpkPokerSet) DealByPattern(i32Pattern int32, iCardCount int) *pokerset.PokerCardSet {
//	return pPoker.pPokerSet.DealByPattern(i32Pattern, iCardCount)
//}
func (pPoker *JlbpkPokerSet) DealByPattern(i32Pattern int32, iCardCount int, bHasAK bool) *pokerset.PokerCardSet {
	pRetSet := pokerset.NewPokerCardSet()

	pCardSet := pPoker.pPokerSet.DealByPattern(i32Pattern, iCardCount)
	mlog.Debug("DealByPattern(), 原牌組[%s] %s", pCardSet.CodeString(), pokerset.PatternString(pCardSet.Pattern()))
	mlog.Debug("DealByPattern(), 牌盒(%d)[%s]", pPoker.CardCount(), pPoker.CodeString())
	if pCardSet.Pattern() == pokerset.POKERPatternHighCard {
		if bHasAK {
			// 若為高牌時必須有 A、K
			pCopy := pokerset.NewPokerCardSet()
			pCopy.SetCardSet(pCardSet.CodeString())

			// 記錄有哪些數字，並統計 A、K 缺的張數
			map_bCheck := make(map[int32]bool)
			for i := 0; i < pCopy.CardCount(); i++ {
				pCard := pCopy.GetSingleCard(i, false)
				if pCard.I32Number == pokerset.POKERCardNo1 || pCard.I32Number == pokerset.POKERCardNoK {
					map_bCheck[pCard.I32Number] = true
				}
			}

			// 補上缺的 A、K
			iCount := 0
			for i := 0; i < pPoker.CardCount(); i++ {
				pCard := pPoker.GetSingleCard(i, false)
				if pCard.I32Number == pokerset.POKERCardNo1 || pCard.I32Number == pokerset.POKERCardNoK {
					if !map_bCheck[pCard.I32Number] {
						pCopy.AddSingleCard(pCard.CodeString())
						pPoker.DrawCards(pCard.CodeString())
						map_bCheck[pCard.I32Number] = true
						iCount++

						if iCount >= 2 {
							break
						}
					}
				}
			}
			mlog.Debug("DealByPattern(), 補上AK[%s]", pCopy.CodeString())

			// 將多餘的非 A、K 的牌抽出來，並補回牌盒中
			for i := 0; i < iCount; i++ {
				for j := 0; j < pCopy.CardCount(); j++ {
					pCard := pCopy.GetSingleCard(j, false)
					if pCard.I32Number != pokerset.POKERCardNo1 && pCard.I32Number != pokerset.POKERCardNoK {
						pCopy.DrawCards(pCard.CodeString())
						pPoker.AddSingleCard(pCard.CodeString())
						break
					}
				}
			}
			pCardSet.SetCardSet(pCopy.CodeString())
		}
	}

	pRetSet.SetCardSet(pCardSet.CodeString())
	mlog.Debug("DealByPattern(), 回傳牌組[%s] %s", pRetSet.CodeString(), pokerset.PatternString(pRetSet.Pattern()))
	mlog.Debug("DealByPattern(), 牌盒(%d)[%s]", pPoker.CardCount(), pPoker.CodeString())

	return pRetSet
}

// 比牌，bCmpColor：是否要比花色,bCmpSingle:是否比較牌型後，比較單牌
func (pPoker *JlbpkPokerSet) Compare(pCardSet1, pCardSet2 *pokerset.PokerCardSet, bCmpColor bool, bCmpSingle bool) int {
	return pPoker.pPokerSet.Compare(pCardSet1, pCardSet2, bCmpColor, bCmpSingle)
}

// 檢查是哪種順子，回傳值：［是否有 A］、［是否為 1-5 順子］
func (pPoker *JlbpkPokerSet) CheckAStraightType(pStraightSet *pokerset.PokerCardSet) (bool, bool) {
	return pPoker.pPokerSet.CheckAStraightType(pStraightSet)
}

// 牌型文字
func (pPoker *JlbpkPokerSet) PatternString(bIsDealer bool) string {
	strRet := ""
	iPattern := pPoker.Pattern(bIsDealer)

	if !bIsDealer {
		strRet = pokerset.PatternString(iPattern)
	} else {
		if iPattern == pokerset.POKERPatternNone {
			strRet = JLBPK_STR_CannotBeAGame
		} else {
			strRet = pokerset.PatternString(iPattern)
		}
	}

	return strRet
}

// 取得 PokerCardSet
func (pPoker *JlbpkPokerSet) GetPokerCardSet() *pokerset.PokerCardSet {
	pSet := new(pokerset.PokerCardSet)
	pSet.SetCardSet(pPoker.CodeString())

	return pSet
}

// 檢查是否包含 A 及 K
func (pPoker *JlbpkPokerSet) CheckIfHasAK() bool {
	map_bAKCheck := map[int32]bool{
		pokerset.POKERCardNoK: false,
		pokerset.POKERCardNoA: false,
	}

	for i := 0; i < pPoker.CardCount(); i++ {
		pCard := pPoker.GetSingleCard(i, true)
		if pCard != nil {
			if _, ok := map_bAKCheck[pCard.I32Number]; ok {
				map_bAKCheck[pCard.I32Number] = true
			}
		}
	}

	for _, bValue := range map_bAKCheck {
		if !bValue {
			return false
		}
	}

	return true
}
