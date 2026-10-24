package mahjong

import (
	"bytes"
)

// Mahjong ?
type Mahjong interface {
	Shuffle()
	ShuffleMore(times int)
	Deal(cn int) *CardGroup
	String() string
	Len() int
	GetCardGroup() *CardGroup
	SetCardGroup(cg *CardGroup) bool
	AddCardGroup(cg *CardGroup)
}

// RandomMahjong 撲克
type RandomMahjong struct {
	Cards []int32
}

// NewRandomSimple 創建麻將,nCopy幾副牌
func NewRandomSimple(nCopy int32) *RandomMahjong {
	mahjong := new(RandomMahjong)
	// 筒子,條子,萬子
	var begin int32 = 1
	for n := begin - 1; n < nCopy; n++ {
		for z := 0; z < 4; z++ {
			for number := begin; number <= 9; number++ {
				for color := begin; color <= ColorWan; color++ {
					cd := 100*color + 10*number + n
					mahjong.Cards = append(mahjong.Cards, cd)
				} //
			}
		}

	}
	return mahjong
}

// NewRandomEmpty 生成空牌
func NewRandomEmpty() *RandomMahjong {
	mahjong := new(RandomMahjong)
	return mahjong
}

// AddCardGroup 生成牌
func (p *RandomMahjong) AddCardGroup(cg *CardGroup) {
	if cg != nil {
		p.Cards = append(p.Cards, cg.Cards...)
	}
}

// TestNewRandom test用
func TestNewRandom() *RandomMahjong {
	mahjong := new(RandomMahjong)
	//筒子,条子,万子
	var begin int32 = 2

	for z := 0; z < 4; z++ {
		for number := begin; number <= 9; number++ {
			for color := begin; color <= ColorWan; color++ {
				cd := 100*color + 10*number + 0
				mahjong.Cards = append(mahjong.Cards, cd)
			} //
		}
	}

	return mahjong
}

// Shuffle 洗牌
func (p *RandomMahjong) Shuffle() {
	p.ShuffleMore(3)
}

// ShuffleOnce 洗牌一次
func (p *RandomMahjong) ShuffleOnce() {
	mcards := make(map[int]int32)

	for idx, c := range p.Cards {
		mcards[idx] = c
	}

	p.Cards = make([]int32, 0)

	for _, c := range mcards {
		p.Cards = append(p.Cards, c)
	}
}

// ShuffleMore 洗多次牌
func (p *RandomMahjong) ShuffleMore(times int) {
	for i := 0; i < times; i++ {
		p.ShuffleOnce()
	}
}

// Deal 發牌
func (p *RandomMahjong) Deal(cn int) *CardGroup {
	if len(p.Cards) < cn {
		return nil
	}

	ret := p.Cards[:cn]
	p.Cards = p.Cards[cn:]
	//fmt.Printf("發牌 ：%v\r\n", ret)
	return NewCardGroup(ret)
}

// GetCardGroup 獲取牌
func (p *RandomMahjong) GetCardGroup() *CardGroup {
	return NewCardGroup(p.Cards)
}

// SetCardGroup 設置牌
func (p *RandomMahjong) SetCardGroup(cg *CardGroup) bool {
	if cg == nil {
		return false
	}
	if cg.Len() != p.Len() {
		return false
	}
	p.Cards = make([]int32, 0)
	p.Cards = append(p.Cards, cg.Cards...)
	return true
}

// String 轉換為字符串
func (p *RandomMahjong) String() string {
	var buf bytes.Buffer
	for _, c := range p.Cards {
		buf.WriteString(CardString(c) + ",")
	}

	return buf.String()
}

// Len ?
func (p *RandomMahjong) Len() int {
	return len(p.Cards)
}
