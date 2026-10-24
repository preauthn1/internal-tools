package zjh

import (	
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

import (
	"strconv"
	"strings"
	"testing"
	"fmt"
)

/**
牌型解釋 :
1. 基礎牌型獲取函式為 GetType 定義在 hhdztype.go 從0開始依序遞增
const (
	DanPai = iota
	DuiZi
	SunZi
	JiaHua
	SunJin
	BaoZi
)
但是因為需要區分大小對子, 固實作一個新的method 用來取分大小對子.
2. 遊戲邏輯內獲取牌型函數為 hhdzPokerStore.getZuheType 函數定義在 hhdzpokerhelper.go, 組合牌型定義在 hhdzlogic.go 裡
const (
	BAOZI       = 1  //豹子
	SHUNJIN     = 2  //順金
	JINHUA      = 3  //金花
	SHUNZI      = 4  //順子
	BIGDUIZI    = 5  //對子(9-A)
	SMALLDUIZI  = 6  //對子(2-8)
	DANPAI      = 7  //單張
	BIGJINHUA   = 11 //A/K金花
	SMALLJINHUA = 12 //非A/K金花
	DUIZI       = 13 //對子
)
*/

//測試案例由測試撰寫
func TestCardType1(t *testing.T) {
	p := poker.NewTestFilePoker("testcardtype1") //牌型定義在 hhdz/testcardtype1.dat
	validateFunction(t, p)
	compareFunction(t, p)

}

//驗證函數由開發實作
func validateFunction(t *testing.T, poker *poker.FilePoker) {
	cfgs := poker.GetConfig()
	userCards := poker.GetUserCards()
	total := 0
	for idx, userPoker := range userCards {
		store := new(zjhPokerStore)
		store.cards = userPoker.Cards
		cardType := store.getZuheType()
		configString := cfgs.GetConfigValue("userpoker", fmt.Sprintf("%d", idx))
		tmpArr := strings.Split(configString, "||")
		if len(tmpArr) < 2 {
			t.Errorf("牌型配置無法驗證缺少了正確牌型代號, 請檢查配置 %d=%s", idx, configString)
			continue
		}
		validCardType, _ := strconv.Atoi(tmpArr[1])
		if cardType == validCardType {
			//t.Logf("%d : %v = %d", idx, store.cards, validCardType)
		} else {
			t.Errorf("驗證錯誤 : %d : %v != %d is %d", idx, store.cards, validCardType, cardType)
		}
		//尚未有賠率計算
		if len(tmpArr) == 3 {
			//betvalue := getBetValue()
			//validOdd, _ := strconv.Atoi(tmpArr[2])
			//odd := betvalue.getOddByType(int32(cardType))
			//if odd != int32(validOdd) {
			//	t.Errorf("驗證錯誤 : %d : %v odd not %d is %d", idx, store.cards, validOdd, odd)
			//}
		}
		total++
	}
	t.Logf("己驗證 %d 個牌型", total)
}

func compareFunction(t *testing.T, p *poker.FilePoker) {
	splitToInt := func(s string, sep string) (int, int) {
		if len(s) == 0 {
			return 0, 0
		}
		uids := strings.Split(s, ",")
		uid0, _ := strconv.Atoi(uids[0])
		uid1, _ := strconv.Atoi(uids[1])
		return uid0, uid1
	}
	cfgs := p.GetConfig()
	userCards := p.GetUserCards()
	total := 0
	sections := cfgs.GetConfigSection("usercompare")
	for idx, configString := range sections {
		tmpArr := strings.Split(configString, "||") //0 as usercard idx, 1 as win or lose, 2 as odd
		uid0, uid1 := splitToInt(tmpArr[0], ",")    //0 紅, 1 黑
		ph := newPokerHelper(nil)
		cards0 := poker.NewCardGroup(userCards[int32(uid0)].Cards)
		cards1 := poker.NewCardGroup(userCards[int32(uid1)].Cards)
		if cards0 == nil || cards1 == nil {
			t.Errorf("牌型錯誤 idx %s cg0 %v cg1 %v", idx, cards0, cards1)
			continue
		}

		if len(cards0.Cards) < 3 || len(cards1.Cards) < 3 {
			t.Errorf("牌型錯誤 idx %s cg0 %v cg1 %v", idx, cards0, cards1)
			continue
		}

		ph.TestInitStorePoker(cards0, cards1)
		jiangType := ph.getJiangType()
		if (jiangType == USER1 && tmpArr[1] != "1") || (jiangType == USER2 && tmpArr[1] != "0") {
			t.Errorf("驗證錯誤 idx %s %v vs %v != %s is %d", idx, ph.user1Store.cards, ph.user2Store.cards, tmpArr[1], jiangType)
		}
		total++
	}
	t.Logf("己驗證 %d 個牌型比對", total)
}
