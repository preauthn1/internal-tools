package inrm

import (
	"math/rand"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestHandCardLengthGameOpt(t *testing.T) {
	opt := HandCardLength(17)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.HandCardLength != 17 {
		t.Errorf("Game.HandCardLength want %d, but %d", 17, game.HandCardLength)
	}
	t.Log(opt)
}

func TestGiveupPointGameOpt(t *testing.T) {
	opt := GiveupPoint(99)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.GiveupPoint != 99 {
		t.Errorf("Game.GiveupPoint want %d, but %d", 99, game.GiveupPoint)
	}
	t.Log(opt)
}

func TestMaxPointGameOpt(t *testing.T) {
	opt := MaxPoint(99)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.MaxPoint != 99 {
		t.Errorf("Game.MaxPoint want %d, but %d", 99, game.MaxPoint)
	}
	t.Log(opt)
}

func TestMissLimitGameOpt(t *testing.T) {
	opt := MissLimit(99)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.MissLimit != 99 {
		t.Errorf("Game.MissLimit want %d, but %d", 99, game.MissLimit)
	}
	t.Log(opt)
}

func TestShuffleTimesGameOpt(t *testing.T) {
	opt := ShuffleTimes(99)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.ShuffleTimes != 99 {
		t.Errorf("Game.ShuffleTimes want %d, but %d", 99, game.ShuffleTimes)
	}
	t.Log(opt)
}

func TestUsePokerNumGameOpt(t *testing.T) {
	opt := UsePokerNum(99)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.UsePokerNum != 99 {
		t.Errorf("Game.UsePokerNum want %d, but %d", 99, game.UsePokerNum)
	}
	t.Log(opt)
}

func TestShufflerGameOpt(t *testing.T) {
	s := poker.NewRandShuffler(rand.NewSource(1))
	opt := Shuffler(s)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.Shuffler != s {
		t.Errorf("Game.Shuffler want %v, but %v", s, game.Shuffler)
	}
	t.Logf("%p %p", s, game.Shuffler)
	t.Log(opt)
}

func TestDealerGameOpt(t *testing.T) {
	d := poker.NaturalDealer
	opt := Dealer(d)
	game, _ := NewGame([]int32{1000, 2000}, opt)
	if game.Dealer != d {
		t.Errorf("Game.Dealer want %v, but %v", d, game.Dealer)
	}
	t.Logf("%p %p", d, game.Dealer)
	t.Log(opt)
}
