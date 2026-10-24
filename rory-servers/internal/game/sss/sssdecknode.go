package sss

import (
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// 13水牌組結構
type SSSDeck struct {
	hc   *netproto.HandCards
	foot *poker.CardGroup //尾墩
	body *poker.CardGroup //中墩
	head *poker.CardGroup //頭墩
}

// sortHCs 排序 netproto.HandCards 結構牌組
func (d *SSSDeck) sortHCs(hc *netproto.HandCards) *netproto.HandCards {
	if hc.Head != nil && hc.Body != nil && hc.Foot != nil {
		// 排序 牌多的放前面 小的放前面
		sort.SliceStable(hc.Head.Cards, func(i, j int) bool {
			var cg = new(poker.CardGroup)
			for _, v := range hc.Head.Cards {
				cg.AddCard(&poker.Card{Color: v.GetColor(), CardNumber: v.GetNumber()})
			}
			var number1 = hc.Head.Cards[i].GetNumber()
			var number2 = hc.Head.Cards[j].GetNumber()
			var count1 = len(cg.GetPokerCardByNumber(number1))
			var count2 = len(cg.GetPokerCardByNumber(number2))
			if count1 > count2 { // 牌多的牌前面
				return true
			} else if count1 < count2 { // 牌少的牌後面
				return false
			} else {
				return number1 < number2 // 相同數量牌點數小的牌前面
			}
		})

		sort.SliceStable(hc.Body.Cards, func(i, j int) bool {
			var cg = new(poker.CardGroup)
			for _, v := range hc.Body.Cards {
				cg.AddCard(&poker.Card{Color: v.GetColor(), CardNumber: v.GetNumber()})
			}
			var number1 = hc.Body.Cards[i].GetNumber()
			var number2 = hc.Body.Cards[j].GetNumber()
			var count1 = len(cg.GetPokerCardByNumber(number1))
			var count2 = len(cg.GetPokerCardByNumber(number2))
			if count1 > count2 {
				return true
			} else if count1 < count2 {
				return false
			} else {
				return number1 < number2
			}
		})

		sort.SliceStable(hc.Foot.Cards, func(i, j int) bool {
			var cg = new(poker.CardGroup)
			for _, v := range hc.Foot.Cards {
				cg.AddCard(&poker.Card{Color: v.GetColor(), CardNumber: v.GetNumber()})
			}
			var number1 = hc.Foot.Cards[i].GetNumber()
			var number2 = hc.Foot.Cards[j].GetNumber()
			var count1 = len(cg.GetPokerCardByNumber(number1))
			var count2 = len(cg.GetPokerCardByNumber(number2))
			if count1 > count2 {
				return true
			} else if count1 < count2 {
				return false
			} else {
				return number1 < number2
			}
		})
	}

	return hc
}

// SetMaxDeck 最大牌型
func (d *SSSDeck) SetMaxDeck(s *SSSTable, cards *poker.CardGroup) {

	d.hc = s.cardTypeGroup(cards, nil)
	d.hc = d.sortHCs(d.hc) // 排序 牌多的放前面 小的放前面

	d.foot = poker.NewCardGroup(nil)
	for _, v := range d.hc.GetFoot().GetCards() {
		d.foot.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	d.body = poker.NewCardGroup(nil)
	for _, v := range d.hc.GetBody().GetCards() {
		d.body.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	d.head = poker.NewCardGroup(nil)
	for _, v := range d.hc.GetHead().GetCards() {
		d.head.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}
}

// IsEqualDeck
func (d *SSSDeck) IsEqualDeck(data *SSSDeck) bool {
	if d.foot.IsEqual(data.foot) && d.body.IsEqual(data.body) && d.head.IsEqual(data.head) {
		return true
	}
	return false
}

// String
func (d *SSSDeck) String() string {
	var headType = netproto.SSSDunCardsNomalType(d.GetHeadCardType())
	var bodyType = netproto.SSSDunCardsNomalType(d.GetBodyCardType())
	var footType = netproto.SSSDunCardsNomalType(d.GetFootCardType())
	return fmt.Sprintf("[%v head:%v] [%v body:%v] [%v foot:%v] ", headType, d.head, bodyType, d.body, footType, d.foot)
}

// GetHeadCardType
func (d *SSSDeck) GetHeadCardType() int32 {
	return d.hc.GetHead().GetCt()
}

// GetBodyCardType
func (d *SSSDeck) GetBodyCardType() int32 {
	return d.hc.GetBody().GetCt()
}

// GetFootCardType
func (d *SSSDeck) GetFootCardType() int32 {
	return d.hc.GetFoot().GetCt()
}

// NextDeckByFoot 尾墩同牌型不同牌，中墩牌型較大或尾墩牌型較大
func (d *SSSDeck) NextDeckByExtra(s *SSSTable, cards *poker.CardGroup) *SSSDeck {
	var hc = s.findNextDeckExtra(cards, d.hc)
	if hc == nil {
		return nil
	}
	hc = d.sortHCs(hc) // 排序 牌多的放前面 小的放前面
	var foot = poker.NewCardGroup(nil)
	for _, v := range hc.GetFoot().GetCards() {
		foot.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var body = poker.NewCardGroup(nil)
	for _, v := range hc.GetBody().GetCards() {
		body.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var head = poker.NewCardGroup(nil)
	for _, v := range hc.GetHead().GetCards() {
		head.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}
	return &SSSDeck{hc: hc, head: head, body: body, foot: foot}
}

// NextDeckByFoot 根據尾墩產生下一組可能牌組
func (d *SSSDeck) NextDeckByFoot(s *SSSTable, cards *poker.CardGroup) *SSSDeck {
	var hc = s.findNextDeck1(cards, d.hc)
	if hc == nil {
		return nil
	}
	hc = d.sortHCs(hc) // 排序 牌多的放前面 小的放前面

	var foot = poker.NewCardGroup(nil)
	for _, v := range hc.GetFoot().GetCards() {
		foot.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var body = poker.NewCardGroup(nil)
	for _, v := range hc.GetBody().GetCards() {
		body.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var head = poker.NewCardGroup(nil)
	for _, v := range hc.GetHead().GetCards() {
		head.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}
	return &SSSDeck{hc: hc, head: head, body: body, foot: foot}
}

// NextDeckByBody 根據中墩產生下一組可能牌組
func (d *SSSDeck) NextDeckByBody(s *SSSTable, cards *poker.CardGroup) *SSSDeck {
	var hc = s.findNextDeck2(cards, d.hc)
	if hc == nil {
		return nil
	}
	hc = d.sortHCs(hc) // 排序 牌多的放前面 小的放前面
	var foot = poker.NewCardGroup(nil)
	for _, v := range hc.GetFoot().GetCards() {
		foot.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var body = poker.NewCardGroup(nil)
	for _, v := range hc.GetBody().GetCards() {
		body.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var head = poker.NewCardGroup(nil)
	for _, v := range hc.GetHead().GetCards() {
		head.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}
	return &SSSDeck{hc: hc, head: head, body: body, foot: foot}
}

// NextDeckByHead 根據頭墩產生下一組可能牌組
func (d *SSSDeck) NextDeckByHead(s *SSSTable, cards *poker.CardGroup) *SSSDeck {
	var hc = s.findNextDeck3(cards, d.hc)
	if hc == nil {
		return nil
	}
	hc = d.sortHCs(hc) // 排序 牌多的放前面 小的放前面
	var foot = poker.NewCardGroup(nil)
	for _, v := range hc.GetFoot().GetCards() {
		foot.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var body = poker.NewCardGroup(nil)
	for _, v := range hc.GetBody().GetCards() {
		body.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}

	var head = poker.NewCardGroup(nil)
	for _, v := range hc.GetHead().GetCards() {
		head.AddCard(&poker.Card{
			Color:      v.GetColor(),
			CardNumber: v.GetNumber(),
		})
	}
	return &SSSDeck{hc: hc, head: head, body: body, foot: foot}
}

type SSSDeckNodeType string

const (
	SSSDeckNodeType_Root   SSSDeckNodeType = "root"
	SSSDeckNodeType_Extra  SSSDeckNodeType = "extra"
	SSSDeckNodeType_Right  SSSDeckNodeType = "right"
	SSSDeckNodeType_Middle SSSDeckNodeType = "middle"
	SSSDeckNodeType_Left   SSSDeckNodeType = "left"
)

type SSSDeckNode struct {
	position int
	nodeType SSSDeckNodeType //node type

	// 尋找牌組策略
	left   *SSSDeckNode // 比頭墩牌型小的牌
	middle *SSSDeckNode // 比中墩牌型小的牌
	right  *SSSDeckNode // 比尾墩牌型小的牌
	extra  *SSSDeckNode // 與尾墩同牌型不同牌，並比中墩牌型小的牌

	data  *SSSDeck // 目前的牌
	layer int
}

func (n *SSSDeckNode) bindData(data *SSSDeck, layer int, nodeType SSSDeckNodeType, position int) *SSSDeckNode {
	n = &SSSDeckNode{data: data, left: nil, middle: nil, right: nil, extra: nil, layer: layer, nodeType: nodeType, position: position}
	return n
}

// // sortHCs 排序
// func (n *SSSDeckNode) sortHCs() {
// 	footCards := n.data.hc.GetHead().GetCards()
// 	sort.SliceStable(footCards, func(i, j int) bool {
// 		return footCards[i].GetNumber() < footCards[j].GetNumber()
// 	})

// 	bodyCards := n.data.hc.GetHead().GetCards()
// 	sort.SliceStable(bodyCards, func(i, j int) bool {
// 		return bodyCards[i].GetNumber() < bodyCards[j].GetNumber()
// 	})

// 	headCards := n.data.hc.GetHead().GetCards()
// 	sort.SliceStable(headCards, func(i, j int) bool {
// 		return headCards[i].GetNumber() < headCards[j].GetNumber()
// 	})
// }

// GetHC 取得目前 node 儲存 netproto.HandCards 的結構
func (n *SSSDeckNode) GetHC() *netproto.HandCards {
	// n.sortHCs()
	return n.data.hc
}

func (n *SSSDeckNode) NextMiddle() *SSSDeckNode {
	return n.middle
}

func (n *SSSDeckNode) NextExtra() *SSSDeckNode {
	return n.extra
}

func (n *SSSDeckNode) NextRight() *SSSDeckNode {
	return n.right
}

func (n *SSSDeckNode) NextLeft() *SSSDeckNode {
	return n.left
}

func (n *SSSDeckNode) String() string {
	if n.data == nil {
		return ""
	}
	return fmt.Sprintf("layer:%d, type:%v, pos:%v %v", n.layer, n.nodeType, n.position, n.data)
}
