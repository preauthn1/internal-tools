package pachislot

import (
	"testing"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func TestPayoutChip(t *testing.T) {

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	odds := int64(10)

	p := NewPayoutChip(odds)

	if p.Type() != pb.PayType_PTChip {
		t.Errorf("PayoutChip Type() want %v, but %v", pb.PayType_PTChip, p.Type())
	}

	if p.Odds() != 10 {
		t.Errorf("PayoutChip Reward() odds want %d, but %d", 10, p.Odds())
	}

	winodds, ok := p.Reward(game)

	if winodds != odds {
		t.Errorf("PayoutChip Reward() win odds want %d, but %d", odds, winodds)
	}

	if !ok {
		t.Errorf("PayoutChip Reward() ok want %v, but %v", true, ok)
	}

	t.Log(p)
}

func TestPayoutReplay(t *testing.T) {

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	p := NewPayoutReplay()

	if p.Type() != pb.PayType_PTReplay {
		t.Errorf("PayoutReplay Type() want %v, but %v", pb.PayType_PTReplay, p.Type())
	}

	if p.Odds() != 0 {
		t.Errorf("PayoutReplay Reward() odds want %d, but %d", 0, p.Odds())
	}

	winodds, ok := p.Reward(game)

	if winodds != 0 {
		t.Errorf("PayoutReplay Reward() win odds want %d, but %d", 0, winodds)
	}

	if !ok {
		t.Errorf("PayoutReplay Reward() ok want %v, but %v", true, ok)
	}

	if !game.Replay {
		t.Errorf("PayoutReplay Reward() Game.Replay want %t, but %t", true, game.Replay)
	}

	t.Log(p)
}

func TestPayoutBonus(t *testing.T) {

	t.Run("BB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Bonus = pb.Scenario_BB

		p := NewPayoutBonus(pb.Scenario_BB)

		if p.Type() != pb.PayType_PTBonus {
			tt.Errorf("PayoutBonus Type() want %v, but %v", pb.PayType_PTBonus, p.Type())
		}

		if p.Module != pb.Scenario_BB {
			tt.Errorf("PayoutBonus Module want %v, but %v", pb.Scenario_BB, p.Module)
		}

		odds := p.Odds()
		if odds != 0 {
			tt.Errorf("PayoutBonus Reward() odds want %d, but %d", 0, odds)
		}

		winodds, ok := p.Reward(game)

		if winodds != 0 {
			tt.Errorf("PayoutBonus Reward() win odds want %d, but %d", 0, winodds)
		}

		if !ok {
			t.Errorf("PayoutBonus Reward() ok want %v, but %v", true, ok)
		}

		if game.Scenario != pb.Scenario_BB {
			tt.Errorf("PayoutBonus Reward() Game.Scenario want %v, but %v", pb.Scenario_BB, game.Scenario)
		}

		if game.BonusTimesLeft != 30 {
			tt.Errorf("PayoutBonus Reward() Game.BonusTimesLeft want %d, but %d", 30, game.BonusTimesLeft)
		}

		t.Log(p)
	})

	t.Run("RB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Bonus = pb.Scenario_RB

		p := NewPayoutBonus(pb.Scenario_RB)

		if p.Type() != pb.PayType_PTBonus {
			tt.Errorf("PayoutBonus Type() want %v, but %v", pb.PayType_PTBonus, p.Type())
		}

		if p.Module != pb.Scenario_RB {
			tt.Errorf("PayoutBonus Module want %v, but %v", pb.Scenario_RB, p.Module)
		}

		odds := p.Odds()
		if odds != 0 {
			tt.Errorf("PayoutBonus Reward() odds want %d, but %d", 0, odds)
		}

		winodds, ok := p.Reward(game)

		if winodds != 0 {
			tt.Errorf("PayoutBonus Reward() win odds want %d, but %d", 0, winodds)
		}

		if !ok {
			t.Errorf("PayoutBonus Reward() ok want %v, but %v", true, ok)
		}

		if game.Scenario != pb.Scenario_RB {
			tt.Errorf("PayoutBonus Reward() Game.Scenario want %v, but %v", pb.Scenario_RB, game.Scenario)
		}

		if game.BonusTimesLeft != 8 {
			tt.Errorf("PayoutBonus Reward() Game.BonusTimesLeft want %d, but %d", 8, game.BonusTimesLeft)
		}

		t.Log(p)
	})

	t.Run("JAC in Times", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Scenario = pb.Scenario_BB
		game.Bonus = pb.Scenario_BB
		game.BonusTimesLeft = 11

		p := NewPayoutBonus(pb.Scenario_JAC)

		if p.Type() != pb.PayType_PTBonus {
			tt.Errorf("PayoutBonus Type() want %v, but %v", pb.PayType_PTBonus, p.Type())
		}

		if p.Module != pb.Scenario_JAC {
			tt.Errorf("PayoutBonus Module want %v, but %v", pb.Scenario_JAC, p.Module)
		}

		odds := p.Odds()
		if odds != 0 {
			tt.Errorf("PayoutBonus Reward() odds want %d, but %d", 0, odds)
		}

		winodds, ok := p.Reward(game)

		if winodds != 0 {
			tt.Errorf("PayoutBonus Reward() win odds want %d, but %d", 0, winodds)
		}

		if !ok {
			t.Errorf("PayoutBonus Reward() ok want %v, but %v", true, ok)
		}

		if game.Scenario != pb.Scenario_JAC {
			tt.Errorf("PayoutBonus Reward() Game.Scenario want %v, but %v", pb.Scenario_JAC, game.Scenario)
		}

		if game.BonusTimesLeft != 11 {
			tt.Errorf("PayoutBonus Reward() Game.BonusTimesLeft want %d, but %d", 11, game.BonusTimesLeft)
		}

		if game.JACTimesLeft != 8 {
			tt.Errorf("PayoutBonus Reward() Game.JACTimesLeft want %d, but %d", 8, game.JACTimesLeft)
		}

		t.Log(p)
	})

	t.Run("JAC over Times", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Scenario = pb.Scenario_BB
		game.Bonus = pb.Scenario_BB
		game.BonusTimesLeft = 11

		for i := 0; i < 3; i++ {
			p := NewPayoutBonus(pb.Scenario_JAC)
			p.Reward(game)

			if game.JACTimes != i+1 {
				tt.Errorf("PayoutBonus Reward() Game.JACTimes want %d, but %d", i+1, game.JACTimes)
			}
		}

		p := NewPayoutBonus(pb.Scenario_JAC)
		_, ok := p.Reward(game)

		if ok {
			t.Errorf("PayoutBonus Reward() ok want %v, but %v", false, ok)
		}

		if game.JACTimes != 3 {
			tt.Errorf("PayoutBonus Reward() Game.JACTimes want %d, but %d", 3, game.JACTimes)
		}
	})
}

func TestPayout1G(t *testing.T) {
	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
	game.Scenario = pb.Scenario_BB
	game.Bonus = pb.Scenario_BB

	p := NewPayout1G(1)

	if p.Type() != pb.PayType_PT1G {
		t.Errorf("Payout1G Type() want %v, but %v", pb.PayType_PT1G, p.Type())
	}

	odds := p.Odds()
	if odds != 0 {
		t.Errorf("Payout1G Reward() odds want %d, but %d", 0, odds)
	}

	winodds, ok := p.Reward(game)

	if winodds != 0 {
		t.Errorf("Payout1G Reward() win odds want %d, but %d", 0, winodds)
	}

	if !ok {
		t.Errorf("Payout1G Reward() ok want %v, but %v", true, ok)
	}

	if game.Scenario != pb.Scenario_BB {
		t.Errorf("Payout1G Reward() Game.Scenario want %v, but %v", pb.Scenario_BB, game.Scenario)
	}

	if game.BonusStore != 1 {
		t.Errorf("Payout1G Reward() Game.BonusStore want %d, but %d", 1, game.BonusStore)
	}

	t.Log(p)
}
