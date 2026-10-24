package poker

import (
	"fmt"
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	PokerColorHeiTao    int32 = 3 // 黑桃
	PokerColorHongTao   int32 = 2 // 红桃
	PokerColorMeiHua    int32 = 1 // 梅花
	PokerColorFangKuai  int32 = 0 // 方塊
	PokerColorSmallKing int32 = 4 // 小王
	PokerColorBigKing   int32 = 5 // 大王
	PokerColorPiZi      int32 = 6 // 皮子
	PokerColorUnKnown   int32 = 7 // 未知
)

// ?
const (
	PokerNumberTwo       int32 = 2  // 2點
	PokerNumberA         int32 = 14 // A點數
	PokerNumberSmallKing int32 = 15 // 小王點數
	PokerNumberBigKing   int32 = 16 // 大王點數
	PokerNumberPiZi      int32 = 17 // 皮子點數
)

// Card 撲克牌
type Card struct {
	Color      int32 // 花色
	CardNumber int32 // 點數
}

// NewPokerCard 創建撲克牌對象
func NewPokerCard(color int32, number int32) *Card {
	pc := new(Card)
	pc.Color = color
	pc.CardNumber = number

	return pc
}

// NewEmptyPokerCard 創建空的撲克牌
func NewEmptyPokerCard() *Card {
	pc := new(Card)
	pc.Color = PokerColorUnKnown
	pc.CardNumber = 0

	return pc
}

// NewPokerCardByStr 根據字符串創建撲克牌
func NewPokerCardByStr(str string, sep string) *Card {
	str = util.TrimAllSpace(str)
	sep = util.TrimAllSpace(sep)

	sarr := strings.Split(str, sep)

	if len(sarr) != 2 {
		return nil
	}

	icolor, err := strconv.Atoi(sarr[0])
	if err != nil {
		return nil
	}

	inum, err := strconv.Atoi(sarr[1])
	if err != nil {
		return nil
	}

	return NewPokerCard(int32(icolor), int32(inum))

}

// NewPokerCardFromNet 根據網絡數據創建撲克牌對象
func NewPokerCardFromNet(npk *netproto.PokerCard) *Card {
	pk := new(Card)
	pk.SetNetMessageData(npk)

	return pk
}

// IsKing 是否大小王
func (pc *Card) IsKing() bool {
	return pc.IsSmallKing() || pc.IsBigKing()
}

// IsSmallKing 是否小王
func (pc *Card) IsSmallKing() bool {
	return pc.Color == PokerColorSmallKing && pc.CardNumber == PokerNumberSmallKing
}

// IsBigKing 是否大王
func (pc *Card) IsBigKing() bool {
	return pc.Color == PokerColorBigKing && pc.CardNumber == PokerNumberBigKing
}

// IsHeiTao 是否黑桃
func (pc *Card) IsHeiTao() bool {
	return pc.Color == PokerColorHeiTao
}

// IsHongTao 是否红桃
func (pc *Card) IsHongTao() bool {
	return pc.Color == PokerColorHongTao
}

// IsMeiHua 是否梅花
func (pc *Card) IsMeiHua() bool {
	return pc.Color == PokerColorMeiHua
}

// IsFangKuai 是否方块
func (pc *Card) IsFangKuai() bool {
	return pc.Color == PokerColorFangKuai
}

// IsPiZi 是否皮子
func (pc *Card) IsPiZi() bool {
	return pc.Color == PokerColorPiZi && pc.CardNumber == PokerNumberPiZi
}

// IsA ?
func (pc *Card) IsA() bool {
	return pc.CardNumber == PokerNumberA
}

// IsEqual 是否同一個牌
func (pc *Card) IsEqual(other *Card) bool {
	return pc.CardNumber == other.CardNumber && pc.Color == other.Color
}

// IsSameCardNumber 是否同一個牌點
func (pc *Card) IsSameCardNumber(other *Card) bool {
	return pc.CardNumber == other.CardNumber && !pc.IsKing() && !pc.IsPiZi()
}

// IsValid 是否正確的牌
func (pc *Card) IsValid() bool {
	if pc.Color < 0 || pc.Color > 6 {
		return false
	}

	if pc.CardNumber < 2 || pc.CardNumber > PokerNumberPiZi {
		return false
	}

	return true
}

func (pc *Card) String() string {
	color := ""
	//♠ ♥ ♣ ♦ ♤ ♡ ♧ ♢
	if pc.IsHeiTao() {
		color = "黑桃"
	} else if pc.IsHongTao() {
		color = "红桃"
	} else if pc.IsMeiHua() {
		color = "梅花"
	} else if pc.IsFangKuai() {
		color = "方块"
	} else if pc.IsSmallKing() {
		return "小王"
	} else if pc.IsBigKing() {
		return "大王"
	} else if pc.IsPiZi() {
		return "皮子"
	}

	if pc.CardNumber <= 10 {
		return fmt.Sprintf("%s%d", color, pc.CardNumber)
	}

	if pc.CardNumber == 11 {
		return fmt.Sprintf("%sJ", color)
	} else if pc.CardNumber == 12 {
		return fmt.Sprintf("%sQ", color)
	} else if pc.CardNumber == 13 {
		return fmt.Sprintf("%sK", color)
	} else if pc.CardNumber == 14 {
		return fmt.Sprintf("%sA", color)
	} else {
		return "invalid card"
	}

}

// GetNetMessageData 獲取牌型網絡數據對象
func (pc *Card) GetNetMessageData() *netproto.PokerCard {
	npk := new(netproto.PokerCard)
	npk.Color = proto.Int32(pc.Color)
	npk.CardNumber = proto.Int32(pc.CardNumber)

	return npk
}

// SetNetMessageData 設置網絡牌型數據對象
func (pc *Card) SetNetMessageData(npk *netproto.PokerCard) {
	pc.Color = npk.GetColor()
	pc.CardNumber = npk.GetCardNumber()
}
