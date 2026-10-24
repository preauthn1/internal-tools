package ttz

import (
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker"
)

//比较牌型大小
func TestCompareCards(t *testing.T) {
	player1CardGroup := new(poker.CardGroup)
	player1Card1 := new(poker.Card)
	player1Card1.CardNumber = 3
	player1Card1.Color = 0
	player1CardGroup.Cards = append(player1CardGroup.Cards, player1Card1)
	player1Card2 := new(poker.Card)
	player1Card2.CardNumber = 4
	player1Card2.Color = 0
	player1CardGroup.Cards = append(player1CardGroup.Cards, player1Card2)

	player2CardGroup := new(poker.CardGroup)
	player2Card1 := new(poker.Card)
	player2Card1.CardNumber = 6
	player2Card1.Color = 0
	player2CardGroup.Cards = append(player2CardGroup.Cards, player2Card1)
	player2Card2 := new(poker.Card)
	player2Card2.CardNumber = 1
	player2Card2.Color = 0
	player2CardGroup.Cards = append(player2CardGroup.Cards, player2Card2)

	nResult := compareCards(player1CardGroup, player2CardGroup)

	if nResult != -1 {
		t.Errorf("點數一樣但玩家2牌型大比較失敗. nResult : %v", nResult)
	}

	player3CardGroup := new(poker.CardGroup)
	player3Card1 := new(poker.Card)
	player3Card1.CardNumber = 3
	player3Card1.Color = 0
	player3CardGroup.Cards = append(player3CardGroup.Cards, player3Card1)
	player4Card2 := new(poker.Card)
	player4Card2.CardNumber = 4
	player4Card2.Color = 0
	player3CardGroup.Cards = append(player3CardGroup.Cards, player4Card2)

	player4CardGroup := new(poker.CardGroup)
	player4Card1 := new(poker.Card)
	player4Card1.CardNumber = 4
	player4Card1.Color = 0
	player4CardGroup.Cards = append(player4CardGroup.Cards, player4Card1)
	player5Card2 := new(poker.Card)
	player5Card2.CardNumber = 3
	player5Card2.Color = 0
	player4CardGroup.Cards = append(player4CardGroup.Cards, player5Card2)

	nResult2 := compareCards(player3CardGroup, player4CardGroup)

	if nResult2 != 0 {
		t.Errorf("點數一樣但牌型一樣比較失敗. nResult : %v", nResult)
	}

}

//是否是豹子
func TestIsLeopard(t *testing.T) {
	playerCardGroup := new(poker.CardGroup)
	player1Card1 := new(poker.Card)
	player1Card1.CardNumber = 0
	player1Card1.Color = 2
	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card1)
	player1Card2 := new(poker.Card)
	player1Card2.CardNumber = 0
	player1Card2.Color = 2
	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card2)

	ok := isLeopard(playerCardGroup)

	if !ok {
		t.Errorf("是否是豹子 錯誤")
	}
}

//是否是对子
func TestIsPair(t *testing.T) {
	playerCardGroup := new(poker.CardGroup)
	player1Card1 := new(poker.Card)
	player1Card1.CardNumber = 3
	player1Card1.Color = 1
	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card1)
	player1Card2 := new(poker.Card)
	player1Card2.CardNumber = 3
	player1Card2.Color = 1
	playerCardGroup.Cards = append(playerCardGroup.Cards, player1Card2)

	ok, result := isPair(playerCardGroup)

	if result != 3 {
		t.Errorf("是否是对子 錯誤")
	}
	if !ok {
		t.Errorf("是否是对子 錯誤")
	}
}
