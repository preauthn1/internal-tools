package inrm

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

// GameOpt 遊戲選項
type GameOpt interface {
	apply(game *Game)
}

type handCardLengthGameOpt struct{ length int }

func (opt *handCardLengthGameOpt) String() string {
	return fmt.Sprintf("玩家手牌: %d 張", opt.length)
}

func (opt *handCardLengthGameOpt) apply(game *Game) { game.HandCardLength = opt.length }

// HandCardLength 正常狀態下手牌數量 Default: 13
func HandCardLength(n int) GameOpt { return &handCardLengthGameOpt{length: n} }

type giveupPointGameOpt struct{ point int }

func (opt *giveupPointGameOpt) String() string {
	return fmt.Sprintf("放棄牌局點數: %d 點", opt.point)
}

func (opt *giveupPointGameOpt) apply(game *Game) { game.GiveupPoint = opt.point }

// GiveupPoint 放棄遊戲時，輸分點數 Default: 40
func GiveupPoint(n int) GameOpt { return &giveupPointGameOpt{point: n} }

type maxPointGameOpt struct{ point int }

func (opt *maxPointGameOpt) String() string {
	return fmt.Sprintf("遊戲輸分最大點數: %d 點", opt.point)
}

func (opt *maxPointGameOpt) apply(game *Game) { game.MaxPoint = opt.point }

// MaxPoint 遊戲輸分最大點數 Default: 80
func MaxPoint(n int) GameOpt { return &maxPointGameOpt{point: n} }

type missGameOpt struct{ n int }

func (opt *missGameOpt) String() string {
	return fmt.Sprintf("Miss上限: %d 次", opt.n)
}

func (opt *missGameOpt) apply(game *Game) { game.MissLimit = opt.n }

// MissLimit 錯過牌局次數上限 Default: 3
func MissLimit(n int) GameOpt { return &missGameOpt{n: n} }

type shuffleTimesGameOpt struct{ n int }

func (opt *shuffleTimesGameOpt) String() string {
	return fmt.Sprintf("每局洗牌次數: %d 次", opt.n)
}

func (opt *shuffleTimesGameOpt) apply(game *Game) { game.ShuffleTimes = opt.n }

// ShuffleTimes 每局洗牌次數 Default: 3
func ShuffleTimes(n int) GameOpt { return &shuffleTimesGameOpt{n: n} }

type usePokerNumGameOpt struct{ n int }

func (opt *usePokerNumGameOpt) String() string {
	return fmt.Sprintf("使用幾副牌: %d 副", opt.n)
}

func (opt *usePokerNumGameOpt) apply(game *Game) { game.UsePokerNum = opt.n }

// UsePokerNum 使用幾副牌 Default: 2
func UsePokerNum(n int) GameOpt { return &usePokerNumGameOpt{n: n} }

type shufflerGameOpt struct{ s poker.Shuffler }

func (opt *shufflerGameOpt) apply(game *Game) { game.Shuffler = opt.s }

// Shuffler 洗牌器 Default: poker.RandShuffler
//   亂數種子 Default: time.Now().Unitnano()
func Shuffler(s poker.Shuffler) GameOpt { return &shufflerGameOpt{s: s} }

type dealerGameOpt struct{ d poker.Dealer }

func (opt *dealerGameOpt) apply(game *Game) { game.Dealer = opt.d }

// Dealer 發牌器 Default: poker.NaturalDealer
func Dealer(d poker.Dealer) GameOpt { return &dealerGameOpt{d: d} }
