//==================================================================
//文件名： ttzpoker.go
//说明：   推筒子发牌
//修改记录：
//==================================================================
package ttz

import (
	"bytes"
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/poker"
)

// color
const (
	TONGZI = 1 //筒子
	BAIBAN = 2 //白板
)

// 打印牌信息
func Print32(cardGroup *poker.CardGroup) {
	for index, card := range cardGroup.Cards {
		fmt.Printf("index=%02d: N=%02d, C=%02d\n", index, card.CardNumber, card.Color)
	}
}

// 转化牌值为string
func StringCardGroup(cg *poker.CardGroup) string {
	var buf bytes.Buffer
	for _, c := range cg.Cards {
		buf.WriteString(fmt.Sprintf("%s,", StringPokerCard(c)))
	}

	return buf.String()
}

// 转换牌值为string
func StringPokerCard(pc *poker.Card) string {
	color := ""

	if pc.Color == TONGZI {
		color = "筒子"
	}
	if pc.Color == BAIBAN {
		color = "白板"
	}

	if pc.CardNumber < 10 {
		return fmt.Sprintf("%s[%d]", color, pc.CardNumber)
	} else {
		return fmt.Sprintf("%s", color)
	}
}

// 生成牌池
func MakeRandomCards() *poker.RandomPoker {

	cards := new(poker.RandomPoker)
	for j := 0; j < 4; j++ {
		for i := 1; i < 10; i++ {
			card := new(poker.Card)
			card.CardNumber = int32(i)
			card.Color = TONGZI
			cards.Cards = append(cards.Cards, card)
		}
	}
	for i := 0; i < 4; i++ {
		card := new(poker.Card)
		card.CardNumber = 10
		card.Color = BAIBAN
		cards.Cards = append(cards.Cards, card)
	}

	cards.Shuffle()
	return cards
}
