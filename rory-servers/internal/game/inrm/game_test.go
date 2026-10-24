package inrm

import (
	"math/rand"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestGame_兩人_正常胡牌(t *testing.T) {

	pids := []int32{1000, 2000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)

	game.Action(p1.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
	))

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.NineClubs))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))

		t.Log(p2)

		game.Action(p2.ID, ThrowCardActionOpt(poker.ThreeDiamonds))
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 2 P1
		game.Action(p1.ID, PopCloseDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.TwoDiamonds))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 2 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))

		t.Log(p2)

		game.Action(p2.ID, ThrowCardActionOpt(poker.FiveSpades))
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 3 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, VictoryActionOpt(poker.SevenSpades,
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts, poker.RedJoker}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds, poker.FiveSpades}}))
	}
}

func TestGame_兩人_一人第一輪棄牌(t *testing.T) {

	pids := []int32{1000, 2000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)

	game.Action(p1.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
	))

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.NineClubs))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P2
		game.Action(p2.ID, GiveupActionOpt())
	}
}

func TestGame_兩人_一人第三輪棄牌(t *testing.T) {

	pids := []int32{1000, 2000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)

	game.Action(p1.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
	))

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.NineClubs))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))

		t.Log(p2)

		game.Action(p2.ID, ThrowCardActionOpt(poker.ThreeDiamonds))
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 2 P1
		game.Action(p1.ID, PopCloseDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.TwoDiamonds))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 2 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))

		t.Log(p2)

		game.Action(p2.ID, ThrowCardActionOpt(poker.FiveSpades))
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 3 P1
		game.Action(p1.ID, GiveupActionOpt())
	}
}

func TestGame_兩人_一人Miss三次(t *testing.T) {

	pids := []int32{1000, 2000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)

	game.Start()
	game.RandJoker = poker.Jack

	game.Action(p1.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
	))

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs, poker.RedJoker}},
		))

		game.Action(p1.ID, ThrowCardActionOpt(poker.NineClubs))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))
	}

	{ // Round 1 P2
		game.Action(p2.ID, MissActionOpt())
	}

	{ // Round 2 P1
		game.Action(p1.ID, PopCloseDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))

		game.Action(p1.ID, ThrowCardActionOpt(poker.SevenSpades))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.RedJoker}},
		))
	}

	{ // Round 2 P2
		game.Action(p2.ID, MissActionOpt())
	}

	{ // Round 3 P1
		game.Action(p1.ID, PopCloseDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.RedJoker}},
		))

		game.Action(p1.ID, ThrowCardActionOpt(poker.TwoDiamonds))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.AceClubs, poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.QueenClubs, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.RedJoker}},
		))
	}

	{ // Round 2 P2
		game.Action(p2.ID, MissActionOpt())
	}

	t.Log(p1)
	t.Log(p2)
}

func TestGame_兩人_一人第一輪詐胡(t *testing.T) {

	pids := []int32{1000, 2000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)

	game.Action(p1.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt(
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
	))

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.NineClubs, poker.QueenClubs, poker.RedJoker}},
		))

		t.Log(p1)

		game.Action(p1.ID, ThrowCardActionOpt(poker.NineClubs))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoClubs, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.KingDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeHearts}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.FiveDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TwoDiamonds, poker.SevenSpades, poker.QueenClubs, poker.RedJoker}},
		))
	}

	t.Log(p1)
	t.Log(p2)

	{ // Round 1 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}},
		))

		t.Log(p2)

		game.Action(p2.ID, VictoryActionOpt(poker.ThreeDiamonds,
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenDiamonds, poker.EightHearts, poker.EightDiamonds, poker.TenDiamonds, poker.TenClubs, poker.TenClubs, poker.JackClubs}}))
	}
}

func TestGame_四人_正常胡牌(t *testing.T) {

	pids := []int32{1000, 2000, 3000, 4000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)
	p3, _ := game.Player(3000)
	p4, _ := game.Player(4000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)

	game.Action(p1.ID, RefreshHandsActionOpt( // [[♥2 ♦2 ♣2 ♠3 ♣4 ♥5 ♦5 ♠6 ♣9 ♦10 ♦J ♥A ♣A]]
		&poker.Deck{Cards: []poker.Card{poker.NineClubs, poker.TenDiamonds, poker.JackDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.FourClubs, poker.FiveDiamonds, poker.SixSpades}},
		&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.TwoDiamonds, poker.TwoClubs, poker.FiveHearts, poker.AceHearts, poker.AceClubs}},
	))
	game.Action(p2.ID, RefreshHandsActionOpt( // [[♦3 ♠4 ♥4 ♠5 ♠5 ♠6 ♦7 ♣7 ♥8 ♦8 ♣10 ♠A 🃟]]
		&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades, poker.SixSpades}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourHearts, poker.FiveSpades}},
		&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenClubs}},
		&poker.Deck{Cards: []poker.Card{poker.EightHearts, poker.EightDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.TenClubs, poker.AceSpades, poker.BlackJoker}},
	))
	game.Action(p3.ID, RefreshHandsActionOpt( // [[🃟R ♥2 ♠3 ♥3 ♦3 ♣3 ♦4 ♠7 ♦8 ♠9 ♦9 ♣Q ♦K]]
		&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.ThreeHearts}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourDiamonds}},
		&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeClubs}},
		&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.EightDiamonds, poker.NineSpades}},
		&poker.Deck{Cards: []poker.Card{poker.NineDiamonds, poker.QueenClubs, poker.KingDiamonds, poker.RedJoker}},
	))
	game.Action(p4.ID, RefreshHandsActionOpt( // [[🃟R ♠2 ♣3 ♣4 ♣5 ♦7 ♠8 ♦10 ♣10 ♣J ♠K ♠K ♥K]]
		&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
		&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.EightSpades, poker.JackClubs}},
		&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.TenClubs, poker.RedJoker}},
		&poker.Deck{Cards: []poker.Card{poker.KingSpades, poker.KingSpades, poker.KingHearts}},
		&poker.Deck{Cards: []poker.Card{poker.TwoSpades}},
	))

	// t.Log(p1)
	// t.Log(p2)
	// t.Log(p3)
	// t.Log(p4)

	{ // Round 1 P1
		game.Action(p1.ID, PopOpenDeckActionOpt())
		game.Action(p1.ID, RefreshHandsActionOpt( // [[♣9 ♦10 ♦J] [♠3 ♣4 ♦5 ♠6] [♥2 ♦2 ♣2 ♥5 ♥A ♣A]] ♣8
			&poker.Deck{Cards: []poker.Card{poker.EightClubs, poker.NineClubs}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.FourClubs, poker.FiveDiamonds, poker.SixSpades}},
			&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.TwoDiamonds, poker.TwoClubs}},
			&poker.Deck{Cards: []poker.Card{poker.FiveHearts, poker.TenDiamonds, poker.JackDiamonds, poker.AceHearts, poker.AceClubs}},
		))

		game.Action(p1.ID, ThrowCardActionOpt(poker.FiveHearts))
		game.Action(p1.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.EightClubs, poker.NineClubs}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.FourClubs, poker.FiveDiamonds, poker.SixSpades}},
			&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.TwoDiamonds, poker.TwoClubs}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.JackDiamonds, poker.AceHearts, poker.AceClubs}},
		))
	}

	{ // Round 1 P2
		game.Action(p2.ID, PopCloseDeckActionOpt())
		game.Action(p2.ID, RefreshHandsActionOpt( // [[♠4 ♠5 ♠6] [♦3 ♥4 ♠5] [♦7 ♣7] [♥8 ♦8] [♣10 ♠A 🃟B]] ♥4
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades, poker.SixSpades}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourHearts, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenClubs}},
			&poker.Deck{Cards: []poker.Card{poker.EightHearts, poker.EightDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TenClubs, poker.AceSpades, poker.BlackJoker, poker.FourHearts}},
		))

		game.Action(p2.ID, ThrowCardActionOpt(poker.FourHearts))
		game.Action(p2.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.FourSpades, poker.FiveSpades, poker.SixSpades}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourHearts, poker.FiveSpades}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.SevenClubs}},
			&poker.Deck{Cards: []poker.Card{poker.EightHearts, poker.EightDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TenClubs, poker.AceSpades, poker.BlackJoker}},
		))
	}

	{ // Round 1 P3
		game.Action(p3.ID, PopOpenDeckActionOpt())
		game.Action(p3.ID, RefreshHandsActionOpt( // [[♥2 ♥3] [♦3 ♦4] [♠3 ♣3] [♠7 ♦8 ♠9] [♦9 ♣Q ♦K 🃟R]] ♥4
			&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.ThreeHearts, poker.FourHearts}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.EightDiamonds, poker.NineSpades}},
			&poker.Deck{Cards: []poker.Card{poker.NineDiamonds, poker.QueenClubs, poker.KingDiamonds, poker.RedJoker}},
		))

		game.Action(p3.ID, ThrowCardActionOpt(poker.NineDiamonds))
		game.Action(p3.ID, RefreshHandsActionOpt(
			&poker.Deck{Cards: []poker.Card{poker.TwoHearts, poker.ThreeHearts, poker.FourHearts}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeDiamonds, poker.FourDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.ThreeSpades, poker.ThreeClubs}},
			&poker.Deck{Cards: []poker.Card{poker.SevenSpades, poker.EightDiamonds, poker.NineSpades}},
			&poker.Deck{Cards: []poker.Card{poker.QueenClubs, poker.KingDiamonds, poker.RedJoker}},
		))
	}

	{ // Round 1 P4
		game.Action(p4.ID, PopOpenDeckActionOpt())
		game.Action(p4.ID, RefreshHandsActionOpt( // [[♣3 ♣4 ♣5] [♦7 ♠8 ♣J] [♦10 ♣10 🃟R] [♠K ♠K ♥K] [♠2]] ♦9
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.EightSpades, poker.NineDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.TenClubs, poker.RedJoker, poker.JackClubs}},
			&poker.Deck{Cards: []poker.Card{poker.KingSpades, poker.KingSpades, poker.KingHearts}},
			&poker.Deck{Cards: []poker.Card{poker.TwoSpades}},
		))

		game.Action(p4.ID, VictoryActionOpt(poker.TwoSpades,
			&poker.Deck{Cards: []poker.Card{poker.ThreeClubs, poker.FourClubs, poker.FiveClubs}},
			&poker.Deck{Cards: []poker.Card{poker.SevenDiamonds, poker.EightSpades, poker.NineDiamonds}},
			&poker.Deck{Cards: []poker.Card{poker.TenDiamonds, poker.TenClubs, poker.RedJoker, poker.JackClubs}},
			&poker.Deck{Cards: []poker.Card{poker.KingSpades, poker.KingSpades, poker.KingHearts}}))
	}

	game.Action(p4.ID, GameOverActionOpt())

	t.Log(game)
	t.Log(p1.LossPoint, p1)
	t.Log(p2.LossPoint, p2)
	t.Log(p3.LossPoint, p3)
	t.Log(p4.LossPoint, p4)
}

func TestGame_四人_初始牌堆(t *testing.T) {

	pids := []int32{1000, 2000, 3000, 4000}

	game, _ := NewGame(pids,
		Shuffler(poker.NewRandShuffler(rand.NewSource(1))),
		ShuffleTimes(1),
	)

	p1, _ := game.Player(1000)
	p2, _ := game.Player(2000)
	p3, _ := game.Player(3000)
	p4, _ := game.Player(4000)

	game.Start()
	game.RandJoker = poker.Jack

	t.Log(game)
	t.Log(p1.LossPoint, p1)
	t.Log(p2.LossPoint, p2)
	t.Log(p3.LossPoint, p3)
	t.Log(p4.LossPoint, p4)
}
