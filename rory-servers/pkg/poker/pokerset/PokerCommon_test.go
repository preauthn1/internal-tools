package pokerset

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// 將 CardID 轉換為 PokerCard
func TestCardID2PokerCard(t *testing.T) {
	for i := POKERCardIDC1; i < POKERCardIDMax; i++ {
		pCard := CardID2PokerCard(i)
		if pCard != nil {
			mlog.Debug("CardID[%d], [%s]", i, pCard.String())
		}
	}
}

// 將 CardCode 轉換為 PokerCardSet
func TestCardCode2PokerCard(t *testing.T) {
	strCodeString := ""
	iCount := util.RndInt(1, 21)
	for i := 0; i < iCount; i++ {
		i32Color := int32(util.RndInt(int(POKERCardColorClub), int(POKERCardColorMax)))
		strCodeString += gMapStrColorCode[i32Color]
		if i32Color == POKERCardColorOther {
			i32Number := int32(util.RndInt(int(POKERCardNoJoker1), int(POKERCardNoHide)))
			strCodeString += gMapStrNumberCode[i32Number]
		} else {
			i32Number := int32(util.RndInt(int(POKERCardNo1), int(POKERCardNoA)))
			strCodeString += gMapStrNumberCode[i32Number]
		}
	}
	mlog.Debug("strCodeString: Count[%d][%s]", iCount, strCodeString)

	pPokerSet := CardCode2PokerCardSet(strCodeString)
	mlog.Debug("牌組[%d張][%s]", pPokerSet.CardCount(), pPokerSet.String())
}

// 牌型字串
func TestPatternString(t *testing.T) {
	for i := POKERPatternNone - 1; i < POKERPatternMax+1; i++ {
		mlog.Debug("i = %d, 牌型[%s]", i, PatternString(i))
	}
}

// 各牌型張數
func TestPatternCardCount(t *testing.T) {
	for i := POKERPatternNone - 1; i < POKERPatternMax+1; i++ {
		mlog.Debug("i = %d, 牌型[%s], 張數[%d]", i, PatternString(i), PatternCardCount(i))
	}
}
