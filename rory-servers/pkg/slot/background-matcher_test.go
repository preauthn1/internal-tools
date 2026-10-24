package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestBackgroundMatcher_MatchAny(t *testing.T) {
	dists := []*protodata.ReelBgDistribution{
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.BackgroundResult = new(BackgroundResult)
	manager.BackgroundResult.Distributions = dists

	t.Log(manager.BackgroundResult)

	// debug.go:54: id:411 groupId:8 linkType:"left" symbol:"bg4" payType:"jackpot" count:5 payout:"25.0" useType:"all"
	// debug.go:54: id:412 groupId:8 linkType:"left" symbol:"bg4" payType:"mysteryJPGR" count:4 payout:"28:58837,38:64000,58:70000,88:67000" useType:"all"
	// debug.go:54: id:413 groupId:8 linkType:"left" symbol:"bg4" payType:"goldreelbonus" count:3 payout:"800.0" useType:"all"

	rule := &protodata.SlotPayTable{LinkType: proto.String("any"), Symbol: proto.String("bg1"), Count: proto.Int32(3)}

	matcher := new(BasicBackgroundMatcher)
	if ok, result := matcher.IsMatch(manager, rule); !ok {
		t.Errorf("Background should be matched")
	} else {
		t.Log(result)
	}
}

func TestBackgroundMatcher_MatchLeft(t *testing.T) {

	t.Run("Normal Left", func(tt *testing.T) {

		dists := []*protodata.ReelBgDistribution{
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			tt.Error(err)
			return
		}
		manager.BackgroundResult = new(BackgroundResult)
		manager.BackgroundResult.Distributions = dists

		tt.Log(manager.BackgroundResult)

		rule := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("bg1"), Count: proto.Int32(3)}

		matcher := new(BasicBackgroundMatcher)
		if ok, result := matcher.IsMatch(manager, rule); !ok {
			tt.Errorf("Background should be matched")
		} else {
			tt.Log(result)
		}
	})

	t.Run("Right Only", func(tt *testing.T) {

		dists := []*protodata.ReelBgDistribution{
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			tt.Error(err)
			return
		}
		manager.BackgroundResult = new(BackgroundResult)
		manager.BackgroundResult.Distributions = dists

		tt.Log(manager.BackgroundResult)

		rule := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("bg2"), Count: proto.Int32(3)}

		matcher := new(BasicBackgroundMatcher)
		if ok, result := matcher.IsMatch(manager, rule); ok {
			tt.Errorf("Background should not be matched")
		} else {
			tt.Log(result)
		}
	})
}

func TestBackgroundMatcher_MatchRight(t *testing.T) {

	t.Run("Normal Right", func(tt *testing.T) {
		dists := []*protodata.ReelBgDistribution{
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			tt.Error(err)
			return
		}
		manager.BackgroundResult = new(BackgroundResult)
		manager.BackgroundResult.Distributions = dists

		tt.Log(manager.BackgroundResult)

		rule := &protodata.SlotPayTable{LinkType: proto.String("right"), Symbol: proto.String("bg0"), Count: proto.Int32(3)}

		matcher := new(BasicBackgroundMatcher)
		if ok, result := matcher.IsMatch(manager, rule); !ok {
			tt.Errorf("Background should be matched")
		} else {
			tt.Log(result)
		}
	})

	t.Run("Left Only", func(tt *testing.T) {
		dists := []*protodata.ReelBgDistribution{
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			tt.Error(err)
			return
		}
		manager.BackgroundResult = new(BackgroundResult)
		manager.BackgroundResult.Distributions = dists

		tt.Log(manager.BackgroundResult)

		rule := &protodata.SlotPayTable{LinkType: proto.String("right"), Symbol: proto.String("bg1"), Count: proto.Int32(3)}

		matcher := new(BasicBackgroundMatcher)
		if ok, result := matcher.IsMatch(manager, rule); ok {
			tt.Errorf("Background should not be matched")
		} else {
			tt.Log(result)
		}
	})

	t.Run("Error Case 1", func(tt *testing.T) {
		dists := []*protodata.ReelBgDistribution{
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
			{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			tt.Error(err)
			return
		}
		manager.BackgroundResult = new(BackgroundResult)
		manager.BackgroundResult.Distributions = dists

		tt.Log(manager.BackgroundResult)

		rule := &protodata.SlotPayTable{LinkType: proto.String("right"), Symbol: proto.String("bg4"), Count: proto.Int32(4)}

		matcher := new(BasicBackgroundMatcher)
		if ok, result := matcher.IsMatch(manager, rule); ok {
			tt.Errorf("Background should not be matched")
		} else {
			tt.Log(result)
		}
	})

}
