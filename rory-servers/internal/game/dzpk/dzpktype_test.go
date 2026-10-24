package dzpk

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"testing"
)

//對子
func TestIsTwoPair(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("isTwoPair")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ok, maxA, minA := isTwoPair(cards)

		if !ok {
			expect := int32(14) //A牌
			if maxA != expect {
				t.Errorf("isTwoPair 最大值錯誤： 得到[%v] , 期望[%v]", maxA, expect)
			}
			expect = int32(10)
			if minA != expect {
				t.Errorf("isTwoPair 第二大值錯誤： 得到[%v] , 期望[%v]", minA, expect)
			}
		}
	}

}

//三條
func TestIsThreeOfAKind(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isThreeOfAKind")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, maxA := isThreeOfAKind(cards)

		if !ok {
			expect := int32(8)
			if maxA != expect {
				t.Errorf("TestIsThreeOfAKind 最大值錯誤： 得到[%v] , 期望[%v]", maxA, expect)
			}
		}
	}

}

//順子
func TestIsStright(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isStright")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, maxA := isStright(cards)

		if !ok {
			expect := int32(11)
			if maxA != expect {
				t.Errorf("TestIsStright 最大值錯誤： 得到[%v] , 期望[%v]", maxA, expect)
			}
		}
	}

}

//順子 1-5
func TestIsStright2(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("isStright2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, maxA := isStright(cards)

		if !ok {
			expect := int32(5)
			if maxA != expect {
				t.Errorf("TestIsStright2 最大值錯誤： 得到[%v] , 期望[%v]", maxA, expect)
			}
		}
	}
}

//同花
func TestIsFlush(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isFlush")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, maxA := isFlush(cards)

		if !ok {
			expect := int32(1)
			if maxA != expect {
				t.Errorf("TestIsFlush 最大值錯誤： 得到[%v] , 期望[%v]", maxA, expect)
			}
		}
	}
}

//葫芦
func TestIsFullHouse(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isFullHouse")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, max, two := isFullHouse(cards)

		if !ok {
			expect := int32(5)
			if max != expect {
				t.Errorf("TestIsFlush 最大值錯誤： 得到[%v] , 期望[%v]", max, expect)
			}
			expect = int32(2)
			if two != expect {
				t.Errorf("TestIsFullHouse 最大值錯誤： 得到[%v] , 期望[%v]", two, expect)
			}
		}
	}

}

//葫芦
func TestIsFullHouse2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("isFullHouse2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)
	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		ok, max, two := isFullHouse(cards)

		if !ok {
			expect := int32(12)
			if max != expect {
				t.Errorf("TestIsFlush 最大值錯誤： 得到[%v] , 期望[%v]", max, expect)
			}
			expect = int32(5)
			if two != expect {
				t.Errorf("TestIsFullHouse 最大值錯誤： 得到[%v] , 期望[%v]", two, expect)
			}
		}
	}

}
