package ddz

import (
	"example.invalid/ServerTeam/Servers/pkg/poker"
	"testing"
)

//对子
func TestDoubleCheck(t *testing.T) {

	var p poker.Poker
	p = poker.NewTestFilePoker("double")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ctx := new(sortCtx)
		item := new(sItem)
		ctx.arr = append(ctx.arr, item)
		ua := new(userAction)
		ua.cg = cards

		cardType := doubleCheck(cards, ctx, ua)

		if cardType != cardTypeDouble && cardType != cardTypeRocket {
			t.Errorf("index:%v , 對子錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}
	}
}

//三张
func TestThreeCheck(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("three")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ctx := new(sortCtx)
		item := new(sItem)
		ctx.arr = append(ctx.arr, item)
		ua := new(userAction)
		ua.cg = cards

		cardType := threeCheck(cards, ctx, ua)

		if cardType != cardTypeThree {
			t.Errorf("index:%v , 三张錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}
	}
}

//炸弹
func TestFourCheck(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("fourBomb")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ctx := new(sortCtx)
		item := new(sItem)
		ctx.arr = append(ctx.arr, item)
		ua := new(userAction)
		ua.cg = cards

		cardType := fourCheck(cards, ctx, ua)

		if cardType != cardTypeBomb {
			t.Errorf("index:%v , 炸弹錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}
	}
}

//单顺:五张或更多的连续单牌五张或更多的连续单牌. 不包括 2 点和双王。
func TestFiveCheck(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("fivestraight")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ctx := new(sortCtx)
		item := new(sItem)
		ctx.arr = append(ctx.arr, item)
		ua := new(userAction)
		ua.cg = cards

		cardType := fiveCheck(cards, ctx, ua)

		if cardType != cardTypeStraight {
			t.Errorf("五張錯誤 cardType : %v , cards : %v", cardType, cards)
		}
	}
}

//三带一
func TestThreeAndOne(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("threeandone")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))
		ctx := new(sortCtx)
		item := new(sItem)
		ctx.arr = append(ctx.arr, item)
		ua := new(userAction)
		ua.cg = cards

		cardType := fourCheck(cards, ctx, ua)

		if cardType != cardTypeThreeAndOne {
			t.Errorf("index:%v , 三带一錯誤fourCheck cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

//飞机带翅膀 带单情况
func TestOtherCheckFeiJiDaiChibangDan(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("otherCheckFeiJiDaiChibangDan")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		cardType := otherCheckFeiJiDaiChibangDan(len(cards.Cards), ctx, ua)

		if cardType != cardTypeDoubleStraights {
			t.Errorf("index:%v, 飞机带翅膀 单 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

//飞机带翅膀 双
func TestOtherCheckFeiJiDaiChibangShuang(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("shuangShun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		cardType := otherCheckFeiJiDaiChibangShuang(len(cards.Cards), ctx, ua)

		if cardType != cardTypeDoubleStraightd {
			t.Errorf("index:%v , 飞机带翅膀 双 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

//連對
func TestOtherCheckShuangShun2(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("shuangShun2")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		cardType := otherCheckShuangShun(ctx, ua)

		if cardType != cardTypeThreeStraight {
			t.Errorf("index:%v , 飞机带翅膀 双 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

//三顺
func TestOtherCheckSanShun(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("sanShun")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		cardType := otherCheckSanShun(ctx, ua)

		if cardType != cardTypeThreeStraight {
			t.Errorf("index:%v , 三順 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

//四带二双
func TestOtherCheckSiDaiEr(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("siDaiEr")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		cardType := otherCheckSiDaiEr(len(cards.Cards), ctx, ua)

		if cardType != cardTypeFourAndTwod && cardType != cardTypeFourAndTwos {
			t.Errorf("index:%v , 四带二双 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}

func TestCardTypeCheck(t *testing.T) {
	var p poker.Poker
	p = poker.NewTestFilePoker("typeCheck")
	p = poker.NewFixedFilePoker(p.(*poker.FilePoker))
	fp := p.(*poker.FixedPoker)

	userLen := fp.UserCardLen()

	for i := 0; i < userLen; i++ {
		cards := fp.GetUserCards(int32(i))

		l := len(cards.Cards)
		ctx := new(sortCtx)
		ctx.arr = make([]*sItem, 0)
		m := make(map[int32]int32)
		for i := 0; i < l; i++ {
			m[cards.Cards[i].CardNumber] += 1
		}
		for k, v := range m {
			item := new(sItem)
			item.number = k
			item.count = v
			ctx.arr = append(ctx.arr, item)
		}

		ua := new(userAction)
		ua.cg = cards

		h := new(history)
		cardType := h.cardTypeCheck(cards, ua)

		if cardType != cardTypeSingle &&
			cardType != cardTypeDouble &&
			cardType != cardTypeRocket &&
			cardType != cardTypeThree &&
			cardType != cardTypeThreeStraight &&
			cardType != cardTypeBomb &&
			cardType != cardTypeDoubleStraightd &&
			cardType != cardTypeDoubleStraights &&
			cardType != cardTypeThreeAndOne &&
			cardType != cardTypeFourAndTwod &&
			cardType != cardTypeFourAndTwos &&
			cardType != cardTypeStraight {
			t.Errorf("index:%v , cardTypeCheck 錯誤 cardType : %v , cards : %v", i, cardType, cards)
		}

	}
}
