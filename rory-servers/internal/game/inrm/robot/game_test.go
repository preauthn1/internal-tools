package robot

import (
	"fmt"
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/internal/game/inrm"
	"example.invalid/ServerTeam/Servers/pkg/poker/v2"
)

func TestGame_四人_自動組牌(t *testing.T) {

	pids := []int32{1000, 2000, 3000, 4000}

	game, _ := inrm.NewGame(pids)
	game.Start()

	done := make(chan bool, 1)

	r := inrm.DefaultRefresher()

	fmt.Println(game)

	game.OnGameOver = func() {
		fmt.Println(game)
		for _, p := range game.Players {
			fmt.Println(p)
		}
	}
	game.OnPlayerVictory = func(p *inrm.Player, c poker.Card, hand []*poker.Deck) {
		fmt.Println(game)
		for _, p := range game.Players {
			game.Action(p.ID, inrm.RefreshHandsActionOpt(p.Hands...))
		}
	}

	go func() {
		for game.State == inrm.GStateGame {
			for _, p := range game.Players {
				if game.State == inrm.GStateOver {
					break
				}
				current, _ := game.CurrentPlayer()
				if current != p {
					continue
				}
				game.Action(p.ID, inrm.PopCloseDeckActionOpt())

				ok, throw, hands := r.Refresh(game.HandCardLength, game.RandJoker, p.Deck)
				if ok {
					game.Action(p.ID, inrm.VictoryActionOpt(throw, hands...))
				} else {
					game.Action(p.ID, inrm.ThrowCardActionOpt(throw))
					game.Action(p.ID, inrm.RefreshHandsActionOpt(hands...))
				}
			}
		}
		done <- true
	}()

	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("timeout")
	}
}
