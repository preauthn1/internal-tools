package inrm

import (
	"reflect"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestPopOpenDeckOpt(t *testing.T) {
	t.Run("過，遊戲中，抽牌一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)

		c := game.OpenDeck.Cards[0]

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", nil, err)
		}

		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}

		if p.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p.LastPop)
		}

		if p.State != PStateGameV2 {
			t.Errorf("Player State want %v, but %v", PStateGameV2, p.State)
		}

		last := p.Hands[len(p.Hands)-1]
		exp := []poker.Card{c}
		if !reflect.DeepEqual(exp, last.Cards) {
			t.Errorf("Player hands want append %v, but %v", exp, last)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)
		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrGameState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，抽牌一再抽一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)

		PopCloseDeckActionOpt().Do(game, p)

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrPlayerState, err)
		}
		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateWait

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateBluff

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家卡牌數不為13", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.Deck.Pop()

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != ErrCardsLenghtCantPop {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", ErrCardsLenghtCantPop, err)
		}
	})
	t.Run("不過，遊戲中，開牌區數量不足", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		game.OpenDeck.Cards = []poker.Card{}
		p, _ := game.Player(1000)

		err := PopOpenDeckActionOpt().Do(game, p)
		if err != poker.ErrCardNotEnough {
			t.Errorf("PopOpenDeckOpt err want %v, but %v", poker.ErrCardNotEnough, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
}

func TestPopCloseDeckOpt(t *testing.T) {
	t.Run("過，遊戲中，抽牌一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateGameV1

		c := game.CloseDeck.Cards[0]

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", nil, err)
		}

		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}

		if p.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p.LastPop)
		}

		if p.State != PStateGameV2 {
			t.Errorf("Player State want %v, but %v", PStateGameV2, p.State)
		}

		last := p.Hands[len(p.Hands)-1]
		exp := []poker.Card{c}
		if !reflect.DeepEqual(exp, last.Cards) {
			t.Errorf("Player hands want append %v, but %v", exp, last)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)
		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrGameState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，抽牌一再抽一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)

		PopCloseDeckActionOpt().Do(game, p)

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}
		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateWait

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateBluff

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家卡牌數不為13", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.Deck.Pop()

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != ErrCardsLenghtCantPop {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrCardsLenghtCantPop, err)
		}
	})
	t.Run("不過，遊戲中，蓋牌區數量不足", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		game.CloseDeck.Cards = []poker.Card{}
		p, _ := game.Player(1000)

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != poker.ErrCardNotEnough {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", poker.ErrCardNotEnough, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("過，遊戲中，蓋牌區最後一張用完", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		cards := game.CloseDeck.Cards[1:]
		game.CloseDeck.Cards = game.CloseDeck.Cards[:1]
		game.OpenDeck = poker.NewEmptyDeck(cards...)

		p, _ := game.Player(1000)

		var isreset bool
		game.OnDeckReset = func() { isreset = true }

		err := PopCloseDeckActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}
		if game.CloseDeck.Len() != 78 {
			t.Errorf("Close Deck length want %d, but %d", 78, game.CloseDeck.Len())
		}
		if game.OpenDeck.Len() != 0 {
			t.Errorf("Open Deck length want %d, but %d", 0, game.OpenDeck.Len())
		}
		if !isreset {
			t.Errorf("OnDeckReset want be called, but not")
		}
	})
}

func TestThrowCardOpt(t *testing.T) {
	t.Run("過，遊戲中，玩家抽一打一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := ThrowCardActionOpt(c).Do(game, p)
		if err != nil {
			t.Errorf("ThrowCardOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.LastThrow != c {
			t.Errorf("Player LastThrow want %d, but %d", c, p.LastThrow)
		}

		last := p.Hands[len(p.Hands)-1]
		exp := []poker.Card{c}
		if reflect.DeepEqual(exp, last.Cards) {
			t.Errorf("Player hands want remove %v, but %v", exp, last)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)
		err := ThrowCardActionOpt(poker.Card(0)).Do(game, p)
		if err != ErrGameState {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家直接打一", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("PopCloseDeckOpt err want %v, but %v", ErrPlayerState, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(2000)

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateBluff

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家卡牌數不為14", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.State = GStateGame
		p, _ := game.Player(1000)
		p.State = PStateGameV2

		err := ThrowCardActionOpt(p.Deck.Cards[0]).Do(game, p)
		if err != ErrCardsLenghtCantThrow {
			t.Errorf("ThrowCardOpt err want %v, but %v", ErrCardsLenghtCantThrow, err)
		}
	})
	t.Run("不過，遊戲中，玩家丟不存在卡牌", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		PopCloseDeckActionOpt().Do(game, p)

		err := ThrowCardActionOpt(0).Do(game, p)
		if err != poker.ErrCardNotFound {
			t.Errorf("ThrowCardOpt err want %v, but %v", poker.ErrCardNotFound, err)
		}

		if p.Deck.Len() != 14 {
			t.Errorf("Player Deck length want %d, but %d", 14, p.Deck.Len())
		}
	})
}

func TestGiveupCardOpt(t *testing.T) {
	t.Run("過，遊戲中，第一輪玩家放棄", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.GiveupPoint = 40
		p, _ := game.Player(1000)

		err := GiveupActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("GiveupOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateFirstGiveup {
			t.Errorf("Player State want %v, but %v", PStateFirstGiveup, p.State)
		}
		if p.LossPoint != 20 {
			t.Errorf("Player LossPoint want %d, but %d", 20, p.LossPoint)
		}
	})
	t.Run("過，遊戲中，第三輪玩家放棄", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 3
		game.GiveupPoint = 40
		p, _ := game.Player(1000)

		err := GiveupActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("GiveupOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateGiveup {
			t.Errorf("Player State want %v, but %v", PStateGiveup, p.State)
		}
		if p.LossPoint != 40 {
			t.Errorf("Player LossPoint want %d, but %d", 40, p.LossPoint)
		}
	})
	t.Run("過，遊戲中，玩家摸牌又放棄，自動打出摸的牌", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := GiveupActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("GiveupOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p.LastPop)
		}
		if game.OpenDeck.Cards[0] != c {
			t.Errorf("OpenDeck 1st card want %v, but %v", c, game.OpenDeck.Cards[0])
		}
		if p.State != PStateFirstGiveup {
			t.Errorf("Player State want %v, but %v", PStateFirstGiveup, p.State)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)

		err := GiveupActionOpt().Do(game, p)
		if err != ErrGameState {
			t.Errorf("GiveupOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(2000)

		err := GiveupActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("GiveupOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		err := GiveupActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("GiveupOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		err := GiveupActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("GiveupOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateBluff

		err := GiveupActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("GiveupOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
}

func TestMissOpt(t *testing.T) {
	t.Run("過，遊戲中，玩家還沒抽牌Miss", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		err := MissActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("MissOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateWait {
			t.Errorf("Player State want %v, but %v", PStateWait, p.State)
		}
		if p.Miss != 1 {
			t.Errorf("Player Miss want %d, but %d", 1, p.Miss)
		}
	})
	t.Run("過，遊戲中，玩家抽了牌Miss", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := MissActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("MissOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateWait {
			t.Errorf("Player State want %v, but %v", PStateWait, p.State)
		}
		if p.Miss != 1 {
			t.Errorf("Player Miss want %d, but %d", 1, p.Miss)
		}
		if p.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p.LastPop)
		}
	})
	t.Run("過，遊戲中，玩家還沒抽牌Miss第三次", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 3
		game.MissLimit = 3
		p, _ := game.Player(1000)
		p.Miss = 2

		err := MissActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("MissOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateGiveup {
			t.Errorf("Player State want %v, but %v", PStateGiveup, p.State)
		}
		if p.Miss != 3 {
			t.Errorf("Player Miss want %d, but %d", 3, p.Miss)
		}
	})
	t.Run("過，遊戲中，玩家抽了牌Miss第三次", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 3
		game.MissLimit = 3
		p, _ := game.Player(1000)
		p.Miss = 2

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := MissActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("MissOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateGiveup {
			t.Errorf("Player State want %v, but %v", PStateGiveup, p.State)
		}
		if p.Miss != 3 {
			t.Errorf("Player Miss want %d, but %d", 3, p.Miss)
		}
		if p.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p.LastPop)
		}
		if game.OpenDeck.Cards[0] != c {
			t.Errorf("OpenDeck 1st card want %v, but %v", c, game.OpenDeck.Cards[0])
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)

		err := MissActionOpt().Do(game, p)
		if err != ErrGameState {
			t.Errorf("MissOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(2000)

		err := MissActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("MissOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		err := MissActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("MissOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		err := MissActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("MissOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateBluff

		err := MissActionOpt().Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("MissOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
}

func TestVictoryOpt(t *testing.T) {
	t.Run("過，遊戲中，玩家手牌可胡", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		d1, _ := poker.ParseDeck("♣J ♣Q ♣K ♣A")
		d2, _ := poker.ParseDeck("♣6 ♦7 ♦8 ♥9 ♣10")
		d3, _ := poker.ParseDeck("♥4 ♥4 ♣4 ♠4")

		p.Hands = []*poker.Deck{d1, d2, d3}
		p.Deck.Clean()
		p.Deck.AppendDeck(d1, d2, d3)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := VictoryActionOpt(c, d1, d2, d3).Do(game, p)
		if err != nil {
			t.Errorf("VictoryOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateVictory {
			t.Errorf("Player State want %v, but %v", PStateVictory, p.State)
		}
		if game.State != GStateCalc {
			t.Errorf("Game State want %v, but %v", GStateCalc, game.State)
		}
	})
	t.Run("過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.MaxPoint = 80
		p, _ := game.Player(1000)

		d1, d2, d3 := poker.NewEmptyDeck(), poker.NewEmptyDeck(), poker.NewEmptyDeck()
		d1.AppendCard(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2.AppendCard(poker.SixClubs, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3.AppendCard(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades, poker.SevenDiamonds)

		p.Hands = []*poker.Deck{d1, d2, d3}
		p.Deck.Cards = []poker.Card{}
		p.Deck.AppendDeck(d1, d2, d3)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p)

		err := VictoryActionOpt(c, p.Hands...).Do(game, p)
		if err != nil {
			t.Errorf("VictoryOpt err want %v, but %v", nil, err)
		}
		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
		if p.State != PStateBluff {
			t.Errorf("Player State want %v, but %v", PStateBluff, p.State)
		}
		if p.LossPoint != 80 {
			t.Errorf("Player LossPoint want %d, but %d", 80, p.LossPoint)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)

		d1 := poker.NewEmptyDeck(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2 := poker.NewEmptyDeck(poker.SixClubs, poker.SevenDiamonds, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3 := poker.NewEmptyDeck(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades)

		err := VictoryActionOpt(poker.RedJoker, d1, d2, d3).Do(game, p)
		if err != ErrGameState {
			t.Errorf("VictoryOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家等待中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(2000)

		d1 := poker.NewEmptyDeck(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2 := poker.NewEmptyDeck(poker.SixClubs, poker.SevenDiamonds, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3 := poker.NewEmptyDeck(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades)

		err := VictoryActionOpt(poker.RedJoker, d1, d2, d3).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("VictoryOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄1", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateFirstGiveup

		d1 := poker.NewEmptyDeck(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2 := poker.NewEmptyDeck(poker.SixClubs, poker.SevenDiamonds, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3 := poker.NewEmptyDeck(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades)

		err := VictoryActionOpt(poker.RedJoker, d1, d2, d3).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("VictoryOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家已放棄2", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateGiveup

		d1 := poker.NewEmptyDeck(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2 := poker.NewEmptyDeck(poker.SixClubs, poker.SevenDiamonds, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3 := poker.NewEmptyDeck(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades)

		err := VictoryActionOpt(poker.RedJoker, d1, d2, d3).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("VictoryOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
	t.Run("不過，遊戲中，玩家詐糊", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)
		p.State = PStateBluff

		d1 := poker.NewEmptyDeck(poker.JackClubs, poker.QueenClubs, poker.KingClubs, poker.AceClubs)
		d2 := poker.NewEmptyDeck(poker.SixClubs, poker.SevenDiamonds, poker.EightDiamonds, poker.NineHearts, poker.TenClubs)
		d3 := poker.NewEmptyDeck(poker.FourHearts, poker.FourHearts, poker.FourClubs, poker.FourSpades)

		err := VictoryActionOpt(poker.RedJoker, d1, d2, d3).Do(game, p)
		if err != ErrPlayerState {
			t.Errorf("VictoryOpt err want %v, but %v", ErrPlayerState, err)
		}

		if p.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p.Deck.Len())
		}
	})
}

func TestPassOpt(t *testing.T) {
	t.Run("過，遊戲中，第一輪第一人，玩家直接跳過", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 1
		game.Idx = 0
		p1, _ := game.Player(1000)

		err := PassActionOpt().Do(game, p1)
		if err != nil {
			t.Errorf("PassOpt err want %v, but %v", nil, err)
		}
		if p1.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p1.Deck.Len())
		}
		if p1.State != PStateWait {
			t.Errorf("Player State want %v, but %v", PStateWait, p1.State)
		}
		p2, _ := game.Player(2000)
		if p2.State != PStateGameV1 {
			t.Errorf("Player2 State want %v, but %v", PStateGameV1, p2.State)
		}
		if game.Round != 1 {
			t.Errorf("Game Round want %d, but %d", 1, game.Round)
		}
		if game.Idx != 1 {
			t.Errorf("Game Idx want %d, but %d", 1, game.Idx)
		}
	})
	t.Run("過，遊戲中，第一輪第一人，玩家抽了牌又跳過", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 1
		game.Idx = 0
		p1, _ := game.Player(1000)

		c := game.CloseDeck.Cards[0]

		PopCloseDeckActionOpt().Do(game, p1)

		err := PassActionOpt().Do(game, p1)
		if err != nil {
			t.Errorf("PassOpt err want %v, but %v", nil, err)
		}
		if p1.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p1.Deck.Len())
		}
		if p1.LastPop != c {
			t.Errorf("Player LastPop want %v, but %v", c, p1.LastPop)
		}
		if p1.LastThrow != c {
			t.Errorf("Player LastThrow want %v, but %v", c, p1.LastThrow)
		}
		if p1.State != PStateWait {
			t.Errorf("Player State want %v, but %v", PStateWait, p1.State)
		}
		p2, _ := game.Player(2000)
		if p2.State != PStateGameV1 {
			t.Errorf("Player2 State want %v, but %v", PStateGameV1, p2.State)
		}
		if game.Round != 1 {
			t.Errorf("Game Round want %d, but %d", 1, game.Round)
		}
		if game.Idx != 1 {
			t.Errorf("Game Idx want %d, but %d", 1, game.Idx)
		}
	})
	t.Run("過，遊戲中，第一輪第二人，玩家直接跳過", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		game.Round = 1
		game.Idx = 1
		p1, _ := game.Player(2000)

		err := PassActionOpt().Do(game, p1)
		if err != nil {
			t.Errorf("PassOpt err want %v, but %v", nil, err)
		}
		if p1.Deck.Len() != 13 {
			t.Errorf("Player Deck length want %d, but %d", 13, p1.Deck.Len())
		}
		if p1.State != PStateWait {
			t.Errorf("Player State want %v, but %v", PStateWait, p1.State)
		}
		p2, _ := game.Player(1000)
		if p2.State != PStateGameV1 {
			t.Errorf("Player1 State want %v, but %v", PStateGameV1, p2.State)
		}
		if game.Round != 2 {
			t.Errorf("Game Round want %d, but %d", 2, game.Round)
		}
		if game.Idx != 0 {
			t.Errorf("Game Idx want %d, but %d", 0, game.Idx)
		}
	})
	t.Run("過，遊戲中，玩家剩一人活著", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		p.State = PStateBluff

		err := PassActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("PassOpt err want %v, but %v", nil, err)
		}
		if game.State != GStateOver {
			t.Errorf("Game State want %v, but %v", GStateOver, game.State)
		}
	})
	t.Run("過，蓋牌用完，不設定自動結束，繼續", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		game.Start()
		p, _ := game.Player(1000)

		game.CloseDeck.Clean()

		err := PassActionOpt().Do(game, p)
		if err != nil {
			t.Errorf("PassOpt err want %v, but %v", nil, err)
		}
		if game.State != GStateGame {
			t.Errorf("Game State want %v, but %v", GStateGame, game.State)
		}
		if game.Round != 1 {
			t.Errorf("Game Round want %d, but %d", 1, game.Round)
		}
		if game.Idx != 1 {
			t.Errorf("Game Idx want %d, but %d", 1, game.Idx)
		}
	})
	t.Run("過，玩家1死了，不再輪到他", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000, 3000})
		game.Start()
		p1, _ := game.Player(1000)
		p2, _ := game.Player(2000)
		p3, _ := game.Player(3000)

		GiveupActionOpt().Do(game, p1)

		// P2 動作
		PassActionOpt().Do(game, p2)

		// P3 動作
		PassActionOpt().Do(game, p3)

		if game.Idx != 1 {
			t.Errorf("Game Idx want %d, but %d", 1, game.Idx)
		}
	})
	t.Run("不過，不再遊戲中", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)

		err := PassActionOpt().Do(game, p)
		if err != ErrGameState {
			t.Errorf("PassOpt err want %v, but %v", ErrGameState, err)
		}
		if p.Deck.Len() != 0 {
			t.Errorf("Player Deck length want %d, but %d", 0, p.Deck.Len())
		}
	})
}

func TestRefreshHandsOpt(t *testing.T) {
	t.Run("過，牌組相同", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)
		p.Deck = poker.NewEmptyDeck(poker.AceSpades, poker.KingSpades, poker.QueenSpades, poker.JackSpades, poker.TenSpades)
		p.Hands = []*poker.Deck{}

		err := RefreshHandsActionOpt(
			poker.NewEmptyDeck(poker.TenSpades),
			poker.NewEmptyDeck(poker.JackSpades),
			poker.NewEmptyDeck(poker.QueenSpades),
			poker.NewEmptyDeck(poker.KingSpades),
			poker.NewEmptyDeck(poker.AceSpades),
		).Do(game, p)
		if err != nil {
			t.Errorf("RefreshHandsOpt err want %v, but %v", nil, err)
		}

		exp := []*poker.Deck{
			poker.NewEmptyDeck(poker.TenSpades),
			poker.NewEmptyDeck(poker.JackSpades),
			poker.NewEmptyDeck(poker.QueenSpades),
			poker.NewEmptyDeck(poker.KingSpades),
			poker.NewEmptyDeck(poker.AceSpades),
		}

		if !reflect.DeepEqual(exp, p.Hands) {
			t.Errorf("Player Hands want %v, but %v", exp, p.Hands)
		}
	})
	t.Run("不過，牌組不同", func(t *testing.T) {
		game, _ := NewGame([]int32{1000, 2000})
		p, _ := game.Player(1000)
		p.Deck = poker.NewEmptyDeck(poker.AceSpades, poker.KingSpades, poker.QueenSpades, poker.JackSpades, poker.TenSpades)
		p.Hands = []*poker.Deck{}

		err := RefreshHandsActionOpt(
			poker.NewEmptyDeck(poker.TenSpades),
			poker.NewEmptyDeck(poker.QueenSpades),
			poker.NewEmptyDeck(poker.KingSpades),
			poker.NewEmptyDeck(poker.AceSpades),
		).Do(game, p)
		if err != ErrHandsDiffDeck {
			t.Errorf("RefreshHandsOpt err want %v, but %v", ErrHandsDiffDeck, err)
		}
	})
}

// Error Case 玩家宣告勝利但最後分數是以點數計算
// 06:30:19Z   玩家[1015452] 宣告勝利 勝利, 手牌 [[♠6 ♠7 ♠8] [♣5 ♣6 ♣7] [♦7 ♦8 ♥9 ♣10] [♥K ♦K ♣K]]
// 06:30:24Z 牌局結束
// 06:30:24Z 玩家[1015452][失敗] 輸 15 分 最終 1,028.21 分，點數 30 牌型 [[♠6 ♠7 ♠8] [♣5 ♣6 ♣7] [♦7 ♦8 ♥9 ♣10] [♥K ♦K ♣K]]
// 06:30:24Z 玩家[1010743][失敗] 輸 40 分 最終 15,501.05 分，點數 80 牌型 [[♠9 ♠10 ♠J] [♥7 ♣7 ♥8 ♦Q ♣Q ♥K 🃟R ♠J ♥J ♣J]]
// 06:30:24Z 玩家[1005084][失敗] 輸 24 分 最終 3,969 分，點數 48 牌型 [[♠A ♠2 ♠3] [♦3 ♣4 ♥5 ♦6] [♣8 ♦10 ♥Q ♣Q ♣K 🃟R]]
// 06:30:24Z 玩家[1004060][失敗] 輸 24.5 分 最終 3,157.5 分，點數 49 牌型 [[♦2 ♦3 ♦4] [♦6 ♠7 ♦8 ♣9 ♣10] [♥9 ♠Q ♠K ♦K ♦J]]
// 06:30:24Z 玩家[1003790][勝利] 贏 103.5 分 稅收 5.17 分 最終 3,380.33 分，點數 0 牌型 [[♥2 ♥3 ♥4] [♣3 ♦4 ♦5] [♣8 ♣9 ♥10] [♠8 ♦9 ♠10] [♣J]]
func TestGameOverOpt_玩家宣告勝利但最後分數是以點數計算(t *testing.T) {

	game, _ := NewGame([]int32{1000, 2000, 3000, 4000, 5000})
	game.Start()

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)
	p3, _ := game.Player(3000)
	p4, _ := game.Player(4000)
	p5, _ := game.Player(5000)

	{ // P1 Cards
		p1.Deck, _ = poker.ParseDeck("♠6 ♠7 ♠8 ♣5 ♣6 ♣7 ♦7 ♦8 ♥9 ♣10 ♥K ♦K ♣K")
		p1.Deck.Sort(poker.SortSameRankAceLast)
		p1.State = PStateVictory

		// [♠6 ♠7 ♠8]
		d1, _ := poker.ParseDeck("♠6 ♠7 ♠8")
		// [♣5 ♣6 ♣7]
		d2, _ := poker.ParseDeck("♣5 ♣6 ♣7")
		// [♦7 ♦8 ♥9 ♣10]
		d3, _ := poker.ParseDeck("♦7 ♦8 ♥9 ♣10")
		// [♥K ♦K ♣K]
		d4, _ := poker.ParseDeck("♥K ♦K ♣K")

		p1.Hands = []*poker.Deck{d1, d2, d3, d4}
	}
	{ // P2 Cards
		p2.Deck, _ = poker.ParseDeck("♠9 ♠10 ♠J ♥7 ♣7 ♥8 ♦Q ♣Q ♥K 🃟R ♠J ♥J ♣J")
		p2.Deck.Sort(poker.SortSameRankAceLast)
		p2.State = PStateWait

		// [♠9 ♠10 ♠J]
		d1, _ := poker.ParseDeck("♠9 ♠10 ♠J")
		// [♥7 ♣7 ♥8 ♦Q ♣Q ♥K 🃟R ♠J ♥J ♣J]
		d2, _ := poker.ParseDeck("♥7 ♣7 ♥8 ♦Q ♣Q ♥K 🃟R ♠J ♥J ♣J")

		p2.Hands = []*poker.Deck{d1, d2}
	}
	{ // P3 Cards
		p3.Deck, _ = poker.ParseDeck("♠A ♠2 ♠3 ♦3 ♣4 ♥5 ♦6 ♣8 ♦10 ♥Q ♣Q ♣K 🃟R")
		p3.Deck.Sort(poker.SortSameRankAceLast)
		p3.State = PStateWait

		// [♠A ♠2 ♠3]
		d1, _ := poker.ParseDeck("♠A ♠2 ♠3")
		// [♦3 ♣4 ♥5 ♦6]
		d2, _ := poker.ParseDeck("♦3 ♣4 ♥5 ♦6")
		// [♣8 ♦10 ♥Q ♣Q ♣K 🃟R]
		d3, _ := poker.ParseDeck("♣8 ♦10 ♥Q ♣Q ♣K 🃟R")

		p3.Hands = []*poker.Deck{d1, d2, d3}
	}
	{ // P4 Cards
		p4.Deck, _ = poker.ParseDeck("♦2 ♦3 ♦4 ♦6 ♠7 ♦8 ♣9 ♣10 ♥9 ♠Q ♠K ♦K ♦J")
		p4.Deck.Sort(poker.SortSameRankAceLast)
		p4.State = PStateWait

		// [♦2 ♦3 ♦4]
		d1, _ := poker.ParseDeck("♦2 ♦3 ♦4")
		// [♦6 ♠7 ♦8 ♣9 ♣10]
		d2, _ := poker.ParseDeck("♦6 ♠7 ♦8 ♣9 ♣10")
		// [♥9 ♠Q ♠K ♦K ♦J]
		d3, _ := poker.ParseDeck("♥9 ♠Q ♠K ♦K ♦J")

		p4.Hands = []*poker.Deck{d1, d2, d3}
	}
	{ // P5 Cards
		p5.Deck, _ = poker.ParseDeck("♥2 ♥3 ♥4 ♣3 ♦4 ♦5 ♣8 ♣9 ♥10 ♠8 ♦9 ♠10 ♣J")
		p5.Deck.Sort(poker.SortSameRankAceLast)
		p5.State = PStateWait

		// [♥2 ♥3 ♥4]
		d1, _ := poker.ParseDeck("♥2 ♥3 ♥4")
		// [♣3 ♦4 ♦5]
		d2, _ := poker.ParseDeck("♣3 ♦4 ♦5")
		// [♣8 ♣9 ♥10]
		d3, _ := poker.ParseDeck("♣8 ♣9 ♥10")
		// [♠8 ♦9 ♠10]
		d4, _ := poker.ParseDeck("♠8 ♦9 ♠10")
		// [♣J]
		d5, _ := poker.ParseDeck("♣J")

		p5.Hands = []*poker.Deck{d1, d2, d3, d4, d5}
	}

	t.Log(game.Action(p1.ID, GameOverActionOpt()))

	t.Log(game)
	t.Log(p1.State, p1.LossPoint)
	t.Log(p2.State, p2.LossPoint)
	t.Log(p3.State, p3.LossPoint)
	t.Log(p4.State, p4.LossPoint)
	t.Log(p5.State, p5.LossPoint)
}
