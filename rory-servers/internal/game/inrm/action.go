package inrm

import (
	"reflect"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// ActionOpt 動作類型
type ActionOpt interface {
	Do(game *Game, p *Player) error
}

// PopOpenDeckActionOpt 從亮牌區取第一張牌
var PopOpenDeckActionOpt = func() ActionOpt { return new(popOpenDeckActionOpt) }

type popOpenDeckActionOpt struct{}

func (opt *popOpenDeckActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	if p.State != PStateGameV1 {
		return ErrPlayerState
	}
	if p.Deck.Len() != game.HandCardLength {
		return ErrCardsLenghtCantPop
	}

	c, err := game.OpenDeck.Pop()
	if err != nil {
		return err
	}

	mlog.Debug("玩家[%d] 抽取亮牌區, 抽到 %v", p.ID, c)

	p.Deck.AppendCard(c)
	p.Deck.Sort(poker.SortSameRankAceLast)

	p.LastPop = c

	p.State = PStateGameV2

	// 自動更新玩家手牌
	p.Hands = append(p.Hands, poker.NewEmptyDeck(c))

	game.OnPlayerPopOpen(c, p)

	return nil
}

// PopCloseDeckActionOpt 從蓋牌區取第一張牌
var PopCloseDeckActionOpt = func() ActionOpt { return new(popCloseDeckActionOpt) }

type popCloseDeckActionOpt struct{}

func (opt *popCloseDeckActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	if p.State != PStateGameV1 {
		return ErrPlayerState
	}
	if p.Deck.Len() != game.HandCardLength {
		return ErrCardsLenghtCantPop
	}

	c, err := game.CloseDeck.Pop()
	if err != nil {
		return err
	}

	mlog.Debug("玩家[%d] 抽取蓋牌區, 抽到 %v", p.ID, c)

	p.Deck.AppendCard(c)
	p.Deck.Sort(poker.SortSameRankAceLast)

	p.LastPop = c

	p.State = PStateGameV2

	// 自動更新玩家手牌
	p.Hands = append(p.Hands, poker.NewEmptyDeck(c))

	game.OnPlayerPopClose(c, p)

	// What if the closed deck runs out of cards and no player has declared yet?
	// Once the last card from the closed deck has been picked up, the discard pile is reshuffled and put face down to become the new closed deck. The card you discard is used to form the new discard pile.
	if game.CloseDeck.Len() == 0 {
		game.CloseDeck, game.OpenDeck = game.OpenDeck, game.CloseDeck
		game.Shuffler.Shuffle(game.CloseDeck, game.ShuffleTimes)
		game.OnDeckReset()
	}

	return nil
}

// ThrowCardActionOpt 丟一張牌到亮牌區
//  @param c poker.Card 要丟棄的牌
var ThrowCardActionOpt = func(c poker.Card) ActionOpt {
	opt := new(thorwCardActionOpt)
	opt.card = c
	return opt
}

type thorwCardActionOpt struct {
	card     poker.Card
	autoPass bool
}

// 不檢查遊戲狀態，玩家狀態，直接丟牌到開牌區
func (opt *thorwCardActionOpt) do(game *Game, p *Player) error {
	if p.Deck.Len() != game.HandCardLength+1 {
		return ErrCardsLenghtCantThrow
	}

	c, err := p.Deck.PickCard(opt.card)
	if err != nil {
		return err
	}

	mlog.Debug("玩家[%d] 丟一張 %v 到亮牌區", p.ID, c)

	game.OpenDeck.StackCard(c)

	p.Deck.Sort(poker.SortSameRankAceLast)

	p.LastThrow = c

	// 自動從玩家手牌中移除
	// 由最後一敦找牌，Client 需要再發送更新手牌的動作確保是使用者的排法
	for i := len(p.Hands) - 1; i >= 0; i-- {
		if _, err := p.Hands[i].PickCard(c); err == nil {
			break
		}
	}

	game.OnPlayerThrow(c, p)

	return nil
}

func (opt *thorwCardActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	if p.State != PStateGameV2 {
		return ErrPlayerState
	}

	if err := opt.do(game, p); err != nil {
		return err
	}

	return PassActionOpt().Do(game, p)
}

// GiveupActionOpt 放棄牌局
//  斷線視同棄牌
//  棄牌後才可離桌或換桌
//  第一局棄牌算20點
//  之後棄牌都算40點
var GiveupActionOpt = func() ActionOpt { return new(giveupActionOpt) }

type giveupActionOpt struct{}

func (opt *giveupActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	switch p.State {
	default:
		return ErrPlayerState
	case PStateGameV1, PStateGameV2:
		mlog.Debug("玩家[%d] 放棄牌局", p.ID)

		if game.Round == 1 {
			p.State = PStateFirstGiveup
			p.LossPoint = game.GiveupPoint / 2

			game.OnPlayerGiveup(p)
			return PassActionOpt().Do(game, p)
		}

		p.State = PStateGiveup
		p.LossPoint = game.GiveupPoint

		game.OnPlayerGiveup(p)
		return PassActionOpt().Do(game, p)
	}
}

// MissActionOpt 錯過牌局
//  若玩家連續錯過三局(可由設置檔更動)沒操作則自動棄牌
var MissActionOpt = func() ActionOpt { return new(missActionOpt) }

type missActionOpt struct{}

func (opt *missActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	switch p.State {
	default:
		return ErrPlayerState
	case PStateGameV1, PStateGameV2:
		p.Miss++
		mlog.Debug("玩家[%d] 錯過牌局", p.ID)

		game.OnPlayerMiss(p)

		if p.Miss >= game.MissLimit {
			return GiveupActionOpt().Do(game, p)
		}
		return PassActionOpt().Do(game, p)
	}
}

// VictoryActionOpt 宣告勝利
//  若牌型違規會算點數,其他人繼續牌局
//  @param c poker.Card 要丟出的卡牌
//  @param hands []*poker.Deck 要判斷的牌型
var VictoryActionOpt = func(c poker.Card, hands ...*poker.Deck) ActionOpt {
	opt := new(victoryActionOpt)
	opt.card = c
	opt.hands = hands
	return opt
}

type victoryActionOpt struct {
	card  poker.Card
	hands []*poker.Deck
}

func (opt *victoryActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}
	if p.State != PStateGameV2 {
		return ErrPlayerState
	}

	if p.Deck.Len() != game.HandCardLength+1 {
		return ErrCardsLenghtCantThrow
	}

	if err := ThrowCardActionOpt(opt.card).(*thorwCardActionOpt).do(game, p); err != nil {
		return err
	}

	// 更新手牌
	if err := RefreshHandsActionOpt(opt.hands...).Do(game, p); err != nil {
		return err
	}

	// 檢查牌型
	if !game.ValidateVectory(p.Hands) {
		mlog.Debug("玩家[%d] 丟棄 %v 宣告勝利, 但牌型不可胡牌, 宣告詐胡", p.ID, p.LastThrow)
		// 違規，詐胡，點數直接為80點(最大點數)
		p.State = PStateBluff
		p.LossPoint = game.MaxPoint

		game.OnPlayerBluff(p, opt.card, opt.hands)

		return PassActionOpt().Do(game, p)
	}

	mlog.Debug("玩家[%d] 丟棄 %v 宣告勝利 牌型 %v", p.ID, p.LastThrow, p.Hands)
	p.State = PStateVictory
	p.final = true
	game.State = GStateCalc

	game.OnPlayerVictory(p, opt.card, opt.hands)

	return nil
}

// PassActionOpt 完成動作跳過，讓指針指向下一個玩家
var PassActionOpt = func() ActionOpt { return new(passActionOpt) }

type passActionOpt struct{}

func (opt *passActionOpt) Do(game *Game, p *Player) error {
	if game.State != GStateGame {
		return ErrGameState
	}

	if p.Deck.Len() == game.HandCardLength+1 { // 檢查玩家存在14張牌，丟一張牌到亮牌區 再放棄
		if err := ThrowCardActionOpt(p.LastPop).(*thorwCardActionOpt).do(game, p); err != nil {
			return err
		}
	}

	// 檢查所有玩家狀態
	if game.AlivePlayerCount() == 1 {
		mlog.Debug("場上僅剩1人, 結束遊戲")
		return GameOverActionOpt().Do(game, p)
	}

	if p.Alive() {
		p.State = PStateWait
	}
	p.Turn++
	n, _ := game.NextPlayer()
	n.State = PStateGameV1

	game.OnPlayerPass(p, n)

	mlog.Debug("%v", game)

	return nil
}

// RefreshHandsActionOpt 更新手牌
var RefreshHandsActionOpt = func(hands ...*poker.Deck) ActionOpt {
	opt := new(refreshHandsActionOpt)
	opt.hands = hands
	return opt
}

type refreshHandsActionOpt struct {
	hands []*poker.Deck
}

func (opt *refreshHandsActionOpt) Do(game *Game, p *Player) error {
	if !p.Alive() {
		return ErrPlayerState
	}

	chk := poker.NewEmptyDeck()
	for _, h := range opt.hands {
		chk.AppendDeck(h)
	}
	deck := poker.NewEmptyDeck(p.Deck.Cards...)
	deck.Sort(poker.SortSameRankAceLast)
	chk.Sort(poker.SortSameRankAceLast)

	if !reflect.DeepEqual(deck, chk) {
		mlog.Debug("玩家[%d] 更新手牌錯誤, %v\n%v\n%v", p.ID, ErrHandsDiffDeck, deck, chk)
		return ErrHandsDiffDeck
	}

	p.Hands = opt.hands

	mlog.Debug("玩家[%d] 更新手牌 %v", p.ID, p.Hands)

	game.OnRefreshHands(p)

	if game.State == GStateCalc {
		p.final = true
		// 如果遊戲已經在結算，所有人都更新完最後一手，就結束遊戲
		for _, p := range game.Players {
			if p.final == false {
				return nil
			}
		}
		mlog.Debug("有玩家都更新完成最後一手牌, 使遊戲結束")
		game.Action(p.ID, GameOverActionOpt())
	}

	return nil
}

// GameOverActionOpt 遊戲結束
var GameOverActionOpt = func() ActionOpt { return new(gameOverActionOpt) }

type gameOverActionOpt struct{}

func (opt *gameOverActionOpt) Do(game *Game, p *Player) error {

	// 計算所有玩家點數
	var min int = game.MaxPoint
	var winner *Player
	for _, p := range game.Players {
		game.CalcPlayerLossPoint(p)
		// If the losing player(s) did not have their first turn before the opponent has made a valid declaration, then :
		// * The losing player gets points equal to half of his/her hand score. For example, if the hand score for a losing player is 30 points, then the player gets 15 points as his score.
		// * The minimum points a losing player can get in this case is 2 (for scores 2, 3, 4, & 5).
		// * The maximum points for a losing player in this case is 40 (for scores 80 or more).
		if p.Alive() && p.Turn == 0 {
			p.LossPoint = p.LossPoint / 2
			switch p.LossPoint {
			case 2, 3, 4, 5:
				p.LossPoint = 2
			}
		}
		if p.LossPoint < min && p.Alive() {
			min = p.LossPoint
		}
		if p.State == PStateVictory {
			winner = p
		}
	}

	for _, p := range game.Players {
		if p.Alive() && winner == nil {
			if p.LossPoint == min {
				p.State = PStateVictory
				p.LossPoint = 0
			}
		}
		if p.State == PStateLoss && p.LossPoint == 0 {
			p.LossPoint = 2
		}
		mlog.Debug("玩家[%d][%v] 點數 %d 牌型 %v", p.ID, p.State, p.LossPoint, p.Hands)
	}

	game.State = GStateOver

	mlog.Debug("%v", game)

	game.OnGameOver()

	return nil
}
