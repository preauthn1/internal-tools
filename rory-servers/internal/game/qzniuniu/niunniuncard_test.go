package qznn

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"testing"
)

//获取牛牛牌点
func TestGetNiuNiuCardPoint(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardPoint")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		cardPoint := getNiuNiuCardPoint(cards.Cards[0])

		if i == 0 && cardPoint != 10 {
			t.Errorf("小王錯誤")
		}

		if i == 1 && cardPoint != 1 {
			t.Errorf("A錯誤")
		}
	}
}

//是否有牛 -1表示无牛
func TestGetNiuNumber(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("niuNumber")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		card, _ := getNiuNumber(cards)

		if card != -1 {
			t.Logf("無牛錯誤")
		}
	}
}

//是否有牛 表示有牛
func TestGetNiuNumber2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("cardPoint2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		card, cardStack := getNiuNumber(cards)

		if card != 8 {
			t.Errorf("有牛錯誤 => 點數 : %v, cardStack : %v", card, cardStack)
		}
	}
}

//是否牛牛
func TestIsNiuNiu(t *testing.T) {
	NewQzniuniuLogicByTest()
	var p poker.Poker
	p = poker.NewTestFilePoker("cardPoint3")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok, cardStack := isNiuNiu(cards)
		if !ok {
			t.Errorf("牛牛錯誤 => 點數 : %v, cardStack : %v", ok, cardStack)
		}
	}
}

//是否五花牛
func TestIsWuHuaNiu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isWuHuaNiu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isWuHuaNiu(cards)
		if !ok {
			t.Errorf("是否五花牛 => 點數 : %v, cardStack : %v", ok, cards)
		}
	}

}

//是否四炸
func TestIs4Zha(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("is4Zha")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok, cardStack := is4Zha(cards)
		if !ok {
			t.Errorf("是否四炸 => 點數 : %v, cardStack : %v", ok, cardStack)
		}
	}
}

//是否五小牛
func TestIsWuXiaoNiu(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isWuXiaoNiu")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok := isWuXiaoNiu(cards)
		if !ok {
			t.Errorf("是否五小牛 => 點數 : %v, cardStack : %v", ok, cards)
		}
	}
}
