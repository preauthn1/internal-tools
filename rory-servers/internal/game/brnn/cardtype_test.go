package brnn

import (
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/config"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

/**
牌型解釋 :
1. 牌型定義在 brnnucardret.go 從0開始依序遞增
const (
	nnCard_WuXiaoNiu = iota //五小牛
	nnCard_4Zha             //四炸
	nnCard_WuHuaNiu         //五花牛
	nnCard_NiuNiu           //牛牛
	nnCard_Niu              //有牛
	nnCard_WuNiu            //无牛
)
2. 遊戲邏輯內牌型獲取函式為 getNiuNiuRet 定義在 brnnucard.go, 牛牛牌型结果, -1 代表無牛, 有牛的話點數為 1-9
type niuniuCardRet struct {
	niuRet    int              //牌型结果
	number    int              //牌点
	cardStack *poker.CardStack //组牌方式
	cards     *poker.CardGroup //牌组
}
*/

//測試案例由測試撰寫
func TestCardType1(t *testing.T) {
	poker := poker.NewTestFilePoker("testcardtype1") //牌型定義在 brnn/testcardtype1.dat
	validateFunction(t, poker)
	compareFunction(t, poker)
}

//驗證函數由開發實作
func validateFunction(t *testing.T, poker *poker.FilePoker) {
	cfgs := poker.GetConfig()
	userCards := poker.GetUserCards()
	total := 0
	for idx, userPoker := range userCards {
		ret := getNiuNiuRet(userPoker)
		configString := cfgs.GetConfigValue("userpoker", fmt.Sprintf("%d", idx))
		tmpArr := strings.Split(configString, "||")
		if len(tmpArr) < 2 {
			t.Errorf("牌型配置無法驗證缺少了正確牌型代號, 請檢查配置 %d=%s", idx, configString)
			continue
		}
		arr := strings.Split(tmpArr[1], ",")
		validCardType, _ := strconv.Atoi(arr[0])
		validCardNumber, _ := strconv.Atoi(arr[1])

		if ret.niuRet == validCardType && ret.number == validCardNumber {
			//t.Logf("%d : %v = %d,%d", idx, ret.cards, validCardType, validCardNumber)
		} else {
			t.Errorf("驗證錯誤 %d : %v != %d,%d is %d,%d", idx, ret.cards, validCardType, validCardNumber, ret.niuRet, ret.number)
		}
		if len(tmpArr) == 3 {
			odd := getNiuNiuOdds(ret)
			validOdd, _ := strconv.Atoi(tmpArr[2])
			if odd != int32(validOdd) {
				t.Errorf("驗證錯誤 : %d : %v odd not %d is %d", idx, ret.cards, validOdd, odd)
			}
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
		tmpArr := strings.Split(configString, "||") //0 as usercard idx, 1 as win or lose
		uid0, uid1 := splitToInt(tmpArr[0], ",")
		cards0 := poker.NewCardGroup(userCards[int32(uid0)].Cards)
		cards1 := poker.NewCardGroup(userCards[int32(uid1)].Cards)
		if cards0 == nil || cards1 == nil {
			t.Errorf("牌型錯誤 idx %s cg0 %v cg1 %v", idx, cards0, cards1)
			continue
		}

		if len(cards0.Cards) < 5 || len(cards1.Cards) < 5 {
			t.Errorf("牌型錯誤 idx %s cg0 %v cg1 %v", idx, cards0, cards1)
			continue
		}
		ret0 := getNiuNiuRet(cards0)
		ret1 := getNiuNiuRet(cards1)
		betret := compareCards(ret0, ret1) //小於0是輸, 大於0是贏
		compare := 2
		if betret < 0 {
			compare = 0
		} else if betret > 0 {
			compare = 1
		}
		if compare == 1 && tmpArr[1] != "1" || compare == 0 && tmpArr[1] != "0" {
			t.Errorf("驗證錯誤 idx %s %v vs %v != %s is %d", idx, cards0, cards1, tmpArr[1], compare)
		}
		total++
	}
	t.Logf("己驗證 %d 個牌型比對", total)
}

func getBetValue() *betconfig {
	_, filename, _, _ := runtime.Caller(1)
	configDir := path.Dir(filename) + "/../../config"
	cfgs := new(config.Group)
	cfgs.LoadConfig(configDir, 0)
	betvalue := newBetconfig(cfgs)
	return betvalue
}
