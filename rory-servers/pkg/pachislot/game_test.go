package pachislot

import (
	"reflect"
	"testing"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func TestNewGame(t *testing.T) {

	t.Run("Pattern Not Found", func(tt *testing.T) {

		Proto := DefaultProto.Clone()
		Proto.Pattern = []*pb.PatternOption{}

		game, err := NewGame(Proto, 1)

		if game != nil {
			tt.Errorf("game want %v, but %v", nil, game)
		}

		if err == nil {
			tt.Errorf("err doesn't want nil")
		}

	})

	t.Run("Mode Not Found", func(tt *testing.T) {

		Proto := DefaultProto.Clone()
		Proto.ModeConfigs = []*pb.Mode{}

		game, err := NewGame(Proto, DefaultPattern[0].ID)

		if game != nil {
			tt.Errorf("game want %v, but %v", nil, game)
		}

		if err == nil {
			tt.Errorf("err doesn't want nil")
		}

	})

	t.Run("Modules Not Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.ModuleConfigs = []*pb.Module{}

		game, err := NewGame(Proto, DefaultPattern[0].ID)

		if game != nil {
			tt.Errorf("game want %v, but %v", nil, game)
		}

		if err == nil {
			tt.Errorf("err doesn't want nil")
		}

	})

	t.Run("RT Release Distribution Not Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{}

		game, err := NewGame(Proto, DefaultPattern[0].ID)

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}

		tt.Logf("%+v\n", game)
	})

	t.Run("Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()

		game, err := NewGame(Proto, DefaultPattern[0].ID)

		if err != nil {
			tt.Errorf("err want %v, but %v", nil, err)
		}

		tt.Logf("%+v\n", game)
	})
}

func TestGame_Init(t *testing.T) {
	t.Run("Before Func", func(tt *testing.T) {
		var cnt int
		var f = func(*Game) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddBeforeInitFunc(f)

		game.Init()

		if cnt != 1 {
			tt.Errorf("Game.Init() execute before func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("After Func", func(tt *testing.T) {
		var cnt int
		var f = func(*Game, bool) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterInitFunc(f)

		game.Init()

		if cnt != 1 {
			tt.Errorf("Game.Init() execute after func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("Next Mode", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.NextMode = pb.ModeType_Heaven

		game.Init()

		if game.Mode != pb.ModeType_Heaven {
			tt.Errorf("Game.Init() Game.Mode want %v, but %v", pb.ModeType_Heaven, game.Mode)
		}
	})
	t.Run("Release Bonus Store", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.BonusStore = 3

		game.Init()

		if game.BonusSure != true {
			tt.Errorf("Game.Init() Game.BonusSure want %v, but %v", true, game.BonusSure)
		}
		if game.BonusStore != 2 {
			tt.Errorf("Game.Init() Game.BonusStore want %v, but %v", 2, game.BonusStore)
		}
	})
	t.Run("After Func, Release Bonus Store", func(tt *testing.T) {
		var cnt int
		var f = func(game *Game, releaseBonusStore bool) {
			cnt++
			if !releaseBonusStore {
				tt.Errorf("Game.Init() AfterInitFunc argument releaseBonusStore want %t, but %t", true, releaseBonusStore)
			}
		}

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterInitFunc(f)

		game.BonusStore = 3
		game.Init()

		if cnt != 1 {
			tt.Errorf("Game.Init() execute after func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("Next Mode Distribution Not Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.NextModeDistribution = []*pb.NextModeDistributionOption{}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.Init()

		if game.NextMode != DefaultMode {
			tt.Errorf("Game.Init() Game.NextMode want %v, but %v", DefaultMode, game.NextMode)
		}
	})
	t.Run("RT Distribution Not Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTDistribution = []*pb.RTDistributionOption{}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.Init()

		if game.RT != DefaultRT {
			tt.Errorf("Game.Init() Game.RT want %v, but %v", DefaultRT, game.RT)
		}
	})
	t.Run("BR Distribution Not Found", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.BRDistribution = []*pb.BRDistributionOption{}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.Init()

		if game.Bonus != DefaultBonus {
			tt.Errorf("Game.Init() Game.Bonus want %v, but %v", DefaultBonus, game.Bonus)
		}
	})
}

func TestGame_Payout(t *testing.T) {
	t.Run("Before Func", func(tt *testing.T) {
		var cnt int
		var f = func(*Game) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddBeforePayoutFunc(f)

		game.Payout(pb.ResultType_None)

		if cnt != 1 {
			tt.Errorf("Game.Payout() execute before func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("After Func", func(tt *testing.T) {
		var cnt int
		var f = func(*Game, []Payout) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterPayoutFunc(f)

		game.Payout(pb.ResultType_None)

		if cnt != 1 {
			tt.Errorf("Game.Payout() execute after func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("Not in UseType", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_NG

		payouts := game.Payout(pb.ResultType_JACNone)

		if len(payouts) != 0 {
			tt.Errorf("Payout length want %d, but %d", 0, len(payouts))
		}
	})
	t.Run("Replay", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_NG

		payouts := game.Payout(pb.ResultType_ReplayLink)

		if len(payouts) != 1 {
			tt.Errorf("Payout length want %d, but %d", 1, len(payouts))
		}

		if payouts[0].Type() != pb.PayType_PTReplay {
			tt.Errorf("Payout Type() want %v, but %v", pb.PayType_PTReplay, payouts[0].Type())
		}
	})
	t.Run("Multi Rule BB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_NG

		payouts := game.Payout(pb.ResultType_BlueSevenLink)

		if len(payouts) != 2 {
			tt.Errorf("Payout length want %d, but %d", 2, len(payouts))
		}

		if payouts[0].Type() != pb.PayType_PTBonus {
			tt.Errorf("Payout[0] Type() want %v, but %v", pb.PayType_PTBonus, payouts[0].Type())
		}

		if payouts[1].Type() != pb.PayType_PTChip {
			tt.Errorf("Payout[1] Type() want %v, but %v", pb.PayType_PTChip, payouts[1].Type())
		}
	})
	t.Run("Multi Rule 1G", func(tt *testing.T) {

		var f = func(game *Game, payouts []Payout) {
			if len(payouts) != 2 {
				tt.Errorf("Payout length want %d, but %d", 2, len(payouts))
			}

			if payouts[0].Type() != pb.PayType_PTChip {
				tt.Errorf("Payout[0] Type() want %v, but %v", pb.PayType_PTChip, payouts[0].Type())
			}

			if payouts[1].Type() != pb.PayType_PT1G {
				tt.Errorf("Payout[1] Type() want %v, but %v", pb.PayType_PT1G, payouts[1].Type())
			}
		}

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterPayoutFunc(f)

		game.Scenario = pb.Scenario_BB

		payouts := game.Payout(pb.ResultType_BlueSevenLink)

		if len(payouts) != 2 {
			tt.Errorf("Payout length want %d, but %d", 2, len(payouts))
		}

		if payouts[0].Type() != pb.PayType_PTChip {
			tt.Errorf("Payout[0] Type() want %v, but %v", pb.PayType_PTChip, payouts[0].Type())
		}

		if payouts[1].Type() != pb.PayType_PT1G {
			tt.Errorf("Payout[1] Type() want %v, but %v", pb.PayType_PT1G, payouts[1].Type())
		}
	})
	t.Run("Bonus Type Error", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.PayTable = []*pb.PayoutOption{
			{ID: 1, GroupID: 1, ResultType: pb.ResultType_BerryLeftAny, PayType: pb.PayType_PTBonus, Payout: 99, UseType: pb.Scenario_ALL},
		}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_BB

		payouts := game.Payout(pb.ResultType_BerryLeftAny)

		if len(payouts) != 0 {
			tt.Errorf("Payout length want %d, but %d", 0, len(payouts))
		}
	})
}

func TestGame_chkRTRelease(t *testing.T) {
	t.Run("RT Less Min Lock", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.RT = 2

		ok, cnt := game.chkRTRelease(DefaultProto.RTReleaseDistribution[0].ResultType)

		if ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", false, ok)
		}

		if cnt != 0 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 0, cnt)
		}
	})
	t.Run("Type Not Found", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		ok, cnt := game.chkRTRelease(pb.ResultType_BlueSevenLink)

		if ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", false, ok)
		}

		if cnt != 0 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 0, cnt)
		}
	})
	t.Run("Not in UseType", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
			{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_BarMiss, Times: []int32{5}, Rate: 100, Denominator: 100, UseType: pb.Scenario_NG},
		}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_BB

		ok, cnt := game.chkRTRelease(pb.ResultType_BarMiss)

		if ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", false, ok)
		}

		if cnt != 0 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 0, cnt)
		}
	})
	t.Run("Release Success", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
			{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_WaterMelonLink, Times: []int32{5}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL},
		}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)
		game.RT = 999
		ok, cnt := game.chkRTRelease(pb.ResultType_WaterMelonLink)

		if !ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", true, ok)
		}

		if cnt != 5 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 5, cnt)
		}
	})
	t.Run("Release Fail", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
			{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_WaterMelonLink, Times: []int32{5}, Rate: 0, Denominator: 100, UseType: pb.Scenario_ALL},
		}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		ok, cnt := game.chkRTRelease(pb.ResultType_WaterMelonLink)

		if ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", false, ok)
		}

		if cnt != 0 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 0, cnt)
		}
	})
	t.Run("Times none", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
			{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_WaterMelonLink, Times: []int32{}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL},
		}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		ok, cnt := game.chkRTRelease(pb.ResultType_WaterMelonLink)

		if ok {
			tt.Errorf("chkRTRelease ok want %t, but %t", false, ok)
		}

		if cnt != 0 {
			tt.Errorf("chkRTRelease cnt want %d, but %d", 0, cnt)
		}
	})
	t.Run("ReleaseRT Result", func(tt *testing.T) {
		Proto := DefaultProto.Clone()
		tt.Run("BB", func(ttt *testing.T) {
			Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
				{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_ReleaseRTonBB, Times: []int32{2}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL},
			}

			game, _ := NewGame(Proto, DefaultPattern[0].ID)
			game.RT = 999
			game.Bonus = pb.Scenario_BB

			ok, cnt := game.chkRTRelease(pb.ResultType_ReleaseRT)

			if !ok {
				ttt.Errorf("chkRTRelease ok want %t, but %t", true, ok)
			}

			if cnt != 2 {
				ttt.Errorf("chkRTRelease cnt want %d, but %d", 2, cnt)
			}
		})
		tt.Run("RB", func(ttt *testing.T) {
			Proto.RTReleaseDistribution = []*pb.RTReleaseOption{
				{ID: 1, GroupID: DefaultPattern[0].RTReleaseGroupID, ResultType: pb.ResultType_ReleaseRTonRB, Times: []int32{2}, Rate: 100, Denominator: 100, UseType: pb.Scenario_ALL},
			}

			game, _ := NewGame(Proto, DefaultPattern[0].ID)
			game.RT = 999
			game.Bonus = pb.Scenario_RB

			ok, cnt := game.chkRTRelease(pb.ResultType_ReleaseRT)

			if !ok {
				ttt.Errorf("chkRTRelease ok want %t, but %t", true, ok)
			}

			if cnt != 2 {
				ttt.Errorf("chkRTRelease cnt want %d, but %d", 2, cnt)
			}
		})
	})
}

func TestGame_chkFinished(t *testing.T) {
	t.Run("NG", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		finished := game.chkFinished()

		if finished {
			tt.Errorf("chkFinished() want %t, but %t", false, finished)
		}
	})
	t.Run("In BB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_BB
		game.BonusTimesLeft = 5

		finished := game.chkFinished()

		if finished {
			tt.Errorf("chkFinished() want %t, but %t", false, finished)
		}
	})
	t.Run("BB Finished", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_BB
		game.BonusTimesLeft = 0

		finished := game.chkFinished()

		if !finished {
			tt.Errorf("chkFinished() want %t, but %t", true, finished)
		}
	})
	t.Run("In RB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_RB
		game.BonusTimesLeft = 5

		finished := game.chkFinished()

		if finished {
			tt.Errorf("chkFinished() want %t, but %t", false, finished)
		}
	})
	t.Run("RB Finished", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_RB
		game.BonusTimesLeft = 0

		finished := game.chkFinished()

		if !finished {
			tt.Errorf("chkFinished() want %t, but %t", true, finished)
		}
	})
	t.Run("In JAC", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_JAC
		game.BonusTimesLeft = 5
		game.JACTimesLeft = 5

		finished := game.chkFinished()

		if finished {
			tt.Errorf("chkFinished() want %t, but %t", false, finished)
		}
	})
	t.Run("JAC Finished In BB", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_JAC
		game.BonusTimesLeft = 1
		game.JACTimesLeft = 0

		finished := game.chkFinished()

		if finished {
			tt.Errorf("chkFinished() want %t, but %t", false, finished)
		}
	})
	t.Run("JAC Finished BB Finished", func(tt *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

		game.Scenario = pb.Scenario_JAC
		game.BonusTimesLeft = 0
		game.JACTimesLeft = 0

		finished := game.chkFinished()

		if !finished {
			tt.Errorf("chkFinished() want %t, but %t", true, finished)
		}
	})
}

func TestGame_randResult(t *testing.T) {
	t.Run("Option Not Found", func(t *testing.T) {

		rt := RandResult(DefaultProto, 1)

		if rt != pb.ResultType_None {
			t.Errorf("ResultType want %v, but %v", pb.ResultType_None, rt)
		}
	})
	t.Run("Pick", func(t *testing.T) {

		Proto := DefaultProto.Clone()
		Proto.ResultDistribution = []*pb.ResultDistributionOption{
			{ID: 1, GroupID: 1, Result: pb.ResultType_BarLink, Weight: 100},
			{ID: 2, GroupID: 1, Result: pb.ResultType_BellLink, Weight: 0},
			{ID: 3, GroupID: 1, Result: pb.ResultType_BerryLink, Weight: 0},
		}

		rt := RandResult(Proto, 1)

		if rt != pb.ResultType_BarLink {
			t.Errorf("ResultType want %v, but %v", pb.ResultType_BarLink, rt)
		}
	})
}

func TestGame_BetCoin(t *testing.T) {
	t.Run("NG", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		coin := game.BetCoin()
		if coin != 3 {
			t.Errorf("Game.BetCoin() want %d, but %d", 3, coin)
		}
	})
	t.Run("NG Replay", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Replay = true
		coin := game.BetCoin()
		if coin != 0 {
			t.Errorf("Game.BetCoin() want %d, but %d", 0, coin)
		}
	})
	t.Run("BB", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Scenario = pb.Scenario_BB
		coin := game.BetCoin()
		if coin != 3 {
			t.Errorf("Game.BetCoin() want %d, but %d", 3, coin)
		}
	})
	t.Run("RB", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Scenario = pb.Scenario_RB
		coin := game.BetCoin()
		if coin != 1 {
			t.Errorf("Game.BetCoin() want %d, but %d", 1, coin)
		}
	})
	t.Run("JAC", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.Scenario = pb.Scenario_JAC
		coin := game.BetCoin()
		if coin != 1 {
			t.Errorf("Game.BetCoin() want %d, but %d", 1, coin)
		}
	})
}

func TestGame_Spin(t *testing.T) {
	t.Run("Before Func", func(t *testing.T) {
		var cnt int
		var f = func(*Game) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddBeforeSpinFunc(f)

		game.Spin()

		if cnt != 1 {
			t.Errorf("Game.Spin() execute before func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("After Func", func(t *testing.T) {
		var cnt int
		var f = func(*Game, pb.ResultType, bool) { cnt++ }

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterSpinFunc(f)

		game.Spin()

		if cnt != 1 {
			t.Errorf("Game.Spin() execute after func count want %d, but %d", 1, cnt)
		}
	})
	t.Run("Normal Play Assigned Result", func(t *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)

		game.RT = 55

		game.NextResult = pb.ResultType_BarMiss

		cost, result, rtrelesd := game.Spin()

		if cost != 3 {
			t.Errorf("Game.Spin() cost want %d, but %d", 3, cost)
		}

		if result != pb.ResultType_BarMiss {
			t.Errorf("Game.Spin() result want %v, but %v", pb.ResultType_BarMiss, result)
		}

		if rtrelesd {
			t.Errorf("Game.Spin() rtrelesd want %t, but %t", false, rtrelesd)
		}

		if game.NextResult != pb.ResultType_Empty {
			t.Errorf("Game.Spin() Game.NextResult want %v, but %v", pb.ResultType_Empty, game.NextResult)
		}

		if game.RT != 54 {
			t.Errorf("Game.Spin() Game.RT want %d, but %d", 54, game.RT)
		}

		if game.RoadMap.SpinCnt != 1 {
			t.Errorf("Game.Spin() Game.RoadMap.SpinCnt want %d, but %d", 1, game.RoadMap.SpinCnt)
		}

		if game.RoadMap.TotalCnt != 1 {
			t.Errorf("Game.Spin() Game.RoadMap.TotalCnt want %d, but %d", 1, game.RoadMap.TotalCnt)
		}
	})
	t.Run("In Replay", func(t *testing.T) {
		Proto := DefaultProto.Clone()
		Proto.RTReleaseDistribution = []*pb.RTReleaseOption{}

		game, _ := NewGame(Proto, DefaultPattern[0].ID)
		game.Replay = true

		cost, _, _ := game.Spin()

		if cost != 0 {
			t.Errorf("Game.Spin() cost want %d, but %d", 0, cost)
		}
	})
	t.Run("RT 0", func(t *testing.T) {
		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.RT = 0

		_, result, _ := game.Spin()

		if result != pb.ResultType_BonusSure {
			t.Errorf("Game.Spin() result want %v, but %v", pb.ResultType_BonusSure, result)
		}

		if !game.BonusSure {
			t.Errorf("Game.Spin() Game.BonusSure want %t, but %t", true, game.BonusSure)
		}

		if game.RT != 0 {
			t.Errorf("Game.Spin() Game.RT want %d, but %d", 0, game.RT)
		}
	})
	t.Run("RT Released", func(t *testing.T) {

		var f = func(g *Game, rt pb.ResultType, rtrelesd bool) {
			if !rtrelesd {
				t.Errorf("Game.Spin() after func rtrelesd want %t, but %t", true, rtrelesd)
			}
		}

		game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
		game.AddAfterSpinFunc(f)

		game.RT = 935
		game.NextResult = pb.ResultType_FeatureNone

		_, _, rtrelesd := game.Spin()

		if !rtrelesd {
			t.Errorf("Game.Spin() rtrelesd want %t, but %t", true, rtrelesd)
		}

		if game.RT != 33 {
			t.Errorf("Game.Spin() Game.RT want %d, but %d", 33, game.RT)
		}
	})
	t.Run("Bonus Sure", func(t *testing.T) {
		t.Run("Not Win", func(t *testing.T) {
			tmp := DefaultBonusSureWinRate
			defer func() {
				DefaultBonusSureWinRate = tmp
			}()

			DefaultBonusSureWinRate = 0

			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.BonusSure = true
			game.NextResult = pb.ResultType_BellLink

			_, result, _ := game.Spin()

			if result != pb.ResultType_BellLink {
				t.Errorf("Game.Spin() result want %v, but %v", pb.ResultType_BellLink, result)
			}

			if !game.BonusSure {
				t.Errorf("Game.Spin() Game.BonusSure want %t, but %t", true, game.BonusSure)
			}
		})
		t.Run("Win No Match Payout", func(t *testing.T) {
			tmpBonusSureWinRate := DefaultBonusSureWinRate

			defer func() {
				DefaultBonusSureWinRate = tmpBonusSureWinRate
			}()

			DefaultBonusSureWinRate = 100

			Proto := DefaultProto.Clone()
			Proto.PayTable = []*pb.PayoutOption{
				{ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTChip, Payout: 10, UseType: pb.Scenario_NG},
				{ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTBonus, Payout: int32(pb.Scenario_BB), UseType: pb.Scenario_NG},
				{ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTBonus, Payout: int32(pb.Scenario_JAC), UseType: pb.Scenario_BB},
			}

			game, _ := NewGame(Proto, DefaultPattern[0].ID)
			game.BonusSure = true
			game.Scenario = pb.Scenario_NG
			game.Bonus = pb.Scenario_RB
			game.NextResult = pb.ResultType_BellLink

			_, result, _ := game.Spin()

			if result != pb.ResultType_BellLink {
				t.Errorf("Game.Spin() result want %v, but %v", pb.ResultType_BellLink, result)
			}

			if !game.BonusSure {
				t.Errorf("Game.Spin() Game.BonusSure want %t, but %t", true, game.BonusSure)
			}
		})
		t.Run("Win", func(t *testing.T) {
			tmpBonusSureWinRate := DefaultBonusSureWinRate

			defer func() {
				DefaultBonusSureWinRate = tmpBonusSureWinRate
			}()

			DefaultBonusSureWinRate = 100

			Proto := DefaultProto.Clone()
			Proto.PayTable = []*pb.PayoutOption{
				{ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTChip, Payout: 10, UseType: pb.Scenario_NG},
				{ResultType: pb.ResultType_BarLink, PayType: pb.PayType_PTBonus, Payout: 6001001, UseType: pb.Scenario_NG},
				{ResultType: pb.ResultType_BlueSevenLink, PayType: pb.PayType_PTBonus, Payout: 6001002, UseType: pb.Scenario_NG},
				{ResultType: pb.ResultType_ReplayLink, PayType: pb.PayType_PTBonus, Payout: 6001003, UseType: pb.Scenario_BB},
			}

			game, _ := NewGame(Proto, DefaultPattern[0].ID)
			game.BonusSure = true
			game.Scenario = pb.Scenario_NG
			game.Bonus = pb.Scenario_RB
			game.NextResult = pb.ResultType_BellLink

			_, result, _ := game.Spin()

			if result != pb.ResultType_BlueSevenLink {
				t.Errorf("Game.Spin() result want %v, but %v", pb.ResultType_BlueSevenLink, result)
			}

			if game.BonusSure {
				t.Errorf("Game.Spin() Game.BonusSure want %t, but %t", false, game.BonusSure)
			}
		})
	})
	t.Run("In BB", func(t *testing.T) {
		t.Run("Spand Times", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_BB
			game.BonusTimesLeft = 21

			cost, _, _ := game.Spin()

			if cost != 3 {
				t.Errorf("Game.Spin() cost want %d, but %d", 3, cost)
			}

			if game.BonusTimesLeft != 20 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 20, game.BonusTimesLeft)
			}
		})
		t.Run("Times End", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_BB
			game.BonusTimesLeft = 1
			game.JACTimesLeft = 0
			game.JACTimes = 3

			game.Spin()

			if game.Scenario != pb.Scenario_BB {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_BB, game.Scenario)
			}

			if game.BonusTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 0, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 0, game.JACTimesLeft)
			}

			if game.JACTimes != 3 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 3, game.JACTimes)
			}
		})
		t.Run("Win JAC in Last Time", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_BB
			game.BonusTimesLeft = 1
			game.JACTimesLeft = 0
			game.JACTimes = 2
			game.NextResult = pb.ResultType_ReplayLink

			_, result, _ := game.Spin()
			game.Reward(result)

			if game.Scenario != pb.Scenario_JAC {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_JAC, game.Scenario)
			}

			if game.BonusTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 0, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 8 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 8, game.JACTimesLeft)
			}

			if game.JACTimes != 3 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 3, game.JACTimes)
			}
		})
	})
	t.Run("In RB", func(t *testing.T) {
		t.Run("Spand Times", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_RB
			game.BonusTimesLeft = 21

			cost, _, _ := game.Spin()

			if cost != 1 {
				t.Errorf("Game.Spin() cost want %d, but %d", 1, cost)
			}

			if game.BonusTimesLeft != 20 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 20, game.BonusTimesLeft)
			}
		})
		t.Run("Times End", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_RB
			game.BonusTimesLeft = 1
			game.JACTimesLeft = 0
			game.JACTimes = 3

			game.Spin()

			if game.Scenario != pb.Scenario_RB {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_RB, game.Scenario)
			}

			if game.BonusTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 0, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 0, game.JACTimesLeft)
			}

			if game.JACTimes != 3 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 3, game.JACTimes)
			}
		})
	})
	t.Run("In JAC", func(t *testing.T) {
		t.Run("Spand Times", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_JAC
			game.Bonus = pb.Scenario_BB
			game.BonusTimesLeft = 1
			game.JACTimesLeft = 5
			game.JACTimes = 2

			cost, _, _ := game.Spin()

			if cost != 1 {
				t.Errorf("Game.Spin() cost want %d, but %d", 1, cost)
			}

			if game.BonusTimesLeft != 1 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 1, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 4 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 4, game.JACTimesLeft)
			}

			if game.JACTimes != 2 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 2, game.JACTimes)
			}

			if game.Scenario != pb.Scenario_JAC {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_JAC, game.Scenario)
			}
		})
		t.Run("Times End", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_JAC
			game.Bonus = pb.Scenario_BB
			game.BonusTimesLeft = 1
			game.JACTimesLeft = 1
			game.JACTimes = 2

			game.Spin()

			if game.BonusTimesLeft != 1 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 1, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 0, game.JACTimesLeft)
			}

			if game.JACTimes != 2 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 2, game.JACTimes)
			}

			if game.Scenario != pb.Scenario_JAC {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_JAC, game.Scenario)
			}
		})
		t.Run("Times End at BB last", func(t *testing.T) {
			game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)
			game.Scenario = pb.Scenario_JAC
			game.Bonus = pb.Scenario_BB
			game.BonusTimesLeft = 0
			game.JACTimesLeft = 1
			game.JACTimes = 2

			game.Spin()

			if game.BonusTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.BonusTimesLeft want %d, but %d", 0, game.BonusTimesLeft)
			}

			if game.JACTimesLeft != 0 {
				t.Errorf("Game.Spin() Game.JACTimesLeft want %d, but %d", 0, game.JACTimesLeft)
			}

			if game.JACTimes != 2 {
				t.Errorf("Game.Spin() Game.JACTimes want %d, but %d", 0, game.JACTimes)
			}

			if game.Scenario != pb.Scenario_JAC {
				t.Errorf("Game.Spin() Game.Scenario want %v, but %v", pb.Scenario_JAC, game.Scenario)
			}
		})
	})
}

func TestGame_AddBeforeInitFunc(t *testing.T) {
	var f1 = func(*Game) { t.Log("f1") }
	var f2 = func(*Game) { t.Log("f2") }
	var f3 = func(*Game) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddBeforeInitFunc(f1)
	game.AddBeforeInitFunc(f2)
	game.AddBeforeInitFunc(f3)

	exp := []func(*Game){f1, f2, f3}

	if reflect.DeepEqual(exp, game.beforeInitFunc) {
		t.Errorf("beforeInitFunc want %v, but %v", exp, game.beforeInitFunc)
	}
}

func TestGame_AddAfterInitFunc(t *testing.T) {
	var f1 = func(*Game, bool) { t.Log("f1") }
	var f2 = func(*Game, bool) { t.Log("f2") }
	var f3 = func(*Game, bool) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddAfterInitFunc(f1)
	game.AddAfterInitFunc(f2)
	game.AddAfterInitFunc(f3)

	exp := []func(*Game, bool){f1, f2, f3}

	if reflect.DeepEqual(exp, game.afterInitFunc) {
		t.Errorf("afterInitFunc want %v, but %v", exp, game.afterInitFunc)
	}
}

func TestGame_AddBeforeSpinFunc(t *testing.T) {
	var f1 = func(*Game) { t.Log("f1") }
	var f2 = func(*Game) { t.Log("f2") }
	var f3 = func(*Game) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddBeforeSpinFunc(f1)
	game.AddBeforeSpinFunc(f2)
	game.AddBeforeSpinFunc(f3)

	exp := []func(*Game){f1, f2, f3}

	if reflect.DeepEqual(exp, game.beforeSpinFunc) {
		t.Errorf("beforeSpinFunc want %v, but %v", exp, game.beforeSpinFunc)
	}
}

func TestGame_AddAfterSpinFunc(t *testing.T) {
	var f1 = func(*Game, pb.ResultType, bool) { t.Log("f1") }
	var f2 = func(*Game, pb.ResultType, bool) { t.Log("f2") }
	var f3 = func(*Game, pb.ResultType, bool) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddAfterSpinFunc(f1)
	game.AddAfterSpinFunc(f2)
	game.AddAfterSpinFunc(f3)

	exp := []func(*Game, pb.ResultType, bool){f1, f2, f3}

	if reflect.DeepEqual(exp, game.afterSpinFunc) {
		t.Errorf("afterSpinFunc want %v, but %v", exp, game.afterSpinFunc)
	}
}

func TestGame_AddBeforePayoutFunc(t *testing.T) {
	var f1 = func(*Game) { t.Log("f1") }
	var f2 = func(*Game) { t.Log("f2") }
	var f3 = func(*Game) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddBeforePayoutFunc(f1)
	game.AddBeforePayoutFunc(f2)
	game.AddBeforePayoutFunc(f3)

	exp := []func(*Game){f1, f2, f3}

	if reflect.DeepEqual(exp, game.beforePayoutFunc) {
		t.Errorf("beforePayoutFunc want %v, but %v", exp, game.beforePayoutFunc)
	}
}

func TestGame_AddAfterPayoutFunc(t *testing.T) {
	var f1 = func(*Game, []Payout) { t.Log("f1") }
	var f2 = func(*Game, []Payout) { t.Log("f2") }
	var f3 = func(*Game, []Payout) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddAfterPayoutFunc(f1)
	game.AddAfterPayoutFunc(f2)
	game.AddAfterPayoutFunc(f3)

	exp := []func(*Game, []Payout){f1, f2, f3}

	if reflect.DeepEqual(exp, game.afterPayoutFunc) {
		t.Errorf("afterPayoutFunc want %v, but %v", exp, game.afterPayoutFunc)
	}
}

func TestGame_AddBeforeRewardFunc(t *testing.T) {
	var f1 = func(*Game) { t.Log("f1") }
	var f2 = func(*Game) { t.Log("f2") }
	var f3 = func(*Game) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddBeforeRewardFunc(f1)
	game.AddBeforeRewardFunc(f2)
	game.AddBeforeRewardFunc(f3)

	exp := []func(*Game){f1, f2, f3}

	if reflect.DeepEqual(exp, game.beforeRewardFunc) {
		t.Errorf("beforeRewardFunc want %v, but %v", exp, game.beforeRewardFunc)
	}
}

func TestGame_AddAfterRewardFunc(t *testing.T) {
	var f1 = func(*Game, []Payout) { t.Log("f1") }
	var f2 = func(*Game, []Payout) { t.Log("f2") }
	var f3 = func(*Game, []Payout) { t.Log("f3") }

	game, _ := NewGame(DefaultProto, DefaultPattern[0].ID)

	game.AddAfterRewardFunc(f1)
	game.AddAfterRewardFunc(f2)
	game.AddAfterRewardFunc(f3)

	exp := []func(*Game, []Payout){f1, f2, f3}

	if reflect.DeepEqual(exp, game.afterRewardFunc) {
		t.Errorf("afterRewardFunc want %v, but %v", exp, game.afterRewardFunc)
	}
}
