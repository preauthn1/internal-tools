package inrm

import (
	"fmt"
	"sort"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// PState 玩家狀態
type PState uint8

// PStates
const (
	// 等待中
	PStateWait PState = iota
	// 遊戲中(只可抽牌，棄權)
	PStateGameV1
	// 遊戲中(只可棄牌，勝利)
	PStateGameV2
	// 棄牌(第一輪)
	PStateFirstGiveup
	// 棄牌
	PStateGiveup
	// 詐胡
	PStateBluff
	// 勝利
	PStateVictory
	// 失敗
	PStateLoss
)

func (s PState) String() string {
	switch s {
	default:
		return "?"
	case PStateWait:
		return "等待"
	case PStateGameV1:
		return "抽/棄"
	case PStateGameV2:
		return "丟/贏"
	case PStateFirstGiveup, PStateGiveup:
		return "棄牌"
	case PStateBluff:
		return "詐胡"
	case PStateVictory:
		return "勝利"
	case PStateLoss:
		return "失敗"
	}
}

// Player 玩家
type Player struct {
	ID   int32
	Deck *poker.Deck

	State PState
	Turn  int // 操作次數

	LossPoint int           // 輸的點數
	Miss      int           // 錯誤局數
	LastPop   poker.Card    // 上次抓牌
	LastThrow poker.Card    // 上次丟牌
	Hands     []*poker.Deck // 手中分牌，單純紀錄，與結算無關

	final bool // 最終更新
}

func (p *Player) String() string {
	return fmt.Sprintf("玩家[%d] 狀態[%v] 抓[%v] %v(%d) 棄[%v] %v", p.ID, p.State, p.LastPop, p.Deck, p.Deck.Len(), p.LastThrow, p.Hands)
}

// NewPlayer 建立玩家
func NewPlayer(uid int32) *Player {
	p := new(Player)
	p.ID = uid
	p.Deck = poker.NewEmptyDeck()
	return p
}

// Alive 是否還活著
func (p *Player) Alive() bool {
	return p.State == PStateGameV1 || p.State == PStateGameV2 || p.State == PStateWait || p.State == PStateVictory || p.State == PStateLoss
}

func (p *Player) initHands() {

	deck := poker.NewEmptyDeck(p.Deck.Cards...)

	for _, s := range poker.RegularSuits {
		d := poker.NewEmptyDeck(deck.PickSuit(s)...)
		if d.Len() > 0 {
			p.Hands = append(p.Hands, d)
		}
	}

	if deck.Len() > 0 {
		p.Hands = append(p.Hands, deck)
	}

	sort.Slice(p.Hands, func(i, j int) bool {
		return p.Hands[i].Len() > p.Hands[j].Len()
	})
}
