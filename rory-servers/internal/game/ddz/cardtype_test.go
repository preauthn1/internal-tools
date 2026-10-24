package ddz

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"strconv"
	"strings"
	"testing"
)

/**
牌型解釋 :
1. 牌型定義在 ddzdefine.go 從0開始依序遞增
const (
	cardTypeNone            cardType = iota //无牌型
	cardTypeRocket                          //火箭1
	cardTypeBomb                            //炸弹2
	cardTypeSingle                          //单牌3
	cardTypeDouble                          //对子4
	cardTypeThree                           //三张5
	cardTypeThreeAndOne                     //三带一6
	cardTypeStraight                        //单顺（顺子）7
	cardTypeDoubleStraights                 //飞机带翅膀 单8
	cardTypeDoubleStraightd                 //飞机带翅膀 双9
	cardTypeThreeStraight                   //連對(對子順)
	//三順 (三條順子)
	cardTypePlaneAndWing //飞机带翅膀
	cardTypeFourAndTwod  //四带二双
	cardTypeFourAndTwos  //四带二单
	cardTypeThreeAndTwo  //三帶二
)
2. 遊戲邏輯內牌型獲取函式為 history.cardTypeCheck 定義在 ddzUserExt.go
 */

//測試案例由測試撰寫
func TestCardType1(t *testing.T) {
	poker := poker.NewTestFilePoker("testcardtype1") //牌型定義在 brnn/testcardtype1.dat
	validateFunction(t, poker)
}

//驗證函數由開發實作
func validateFunction(t *testing.T, poker *poker.FilePoker) {
	cfgs := poker.GetConfig()
	userCards := poker.GetUserCards()
	total := 0
	for idx, userPoker := range userCards {
		ua := new(userAction)
		ua.cg = userPoker

		h := new(history)
		cardType := h.cardTypeCheck(userPoker, ua)
		configString := cfgs.GetConfigValue("userpoker", fmt.Sprintf("%d", idx))
		tmpArr := strings.Split(configString, "||")
		if len(tmpArr) < 2 {
			t.Errorf("牌型配置無法驗證缺少了正確牌型代號, 請檢查配置 %d=%s", idx, configString)
			continue
		}
		validCardType, _ := strconv.Atoi(tmpArr[1])
		if int(cardType) == validCardType {
			//t.Logf("%d : %v = %d and %d", idx, userPoker.Cards, validCardType, cardType)
		} else {
			t.Errorf("%d : %v != %d is %d", idx, userPoker.Cards, validCardType, cardType)
		}
		total++
	}
	//t.Logf("己驗證 %d 個牌型", total)
}

