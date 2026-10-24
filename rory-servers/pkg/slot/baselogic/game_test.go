package baselogic

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestSlotLogicBase_CalcGameResult_FixedMatcher_Filter_OneLineTwoRule(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(nil, nil, nil)

	logic.Proto = new(protodata.ProtoData)

	logic.Proto.SymbolAttributes = []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String(""), Multiplier: proto.Int32(1)},
	}

	rule1 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(5), Payout: proto.String("150")}
	rule2 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("100")}
	rule3 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(3), Payout: proto.String("50")}
	rule4 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("200")}
	rule5 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(3), Payout: proto.String("100")}
	rule6 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(2), Payout: proto.String("50")}
	logic.Proto.SlotPayTable = []*protodata.SlotPayTable{rule1, rule2, rule3, rule4, rule5, rule6}

	logic.Proto.SlotPayLine = []*protodata.SlotPayLine{
		{Id: proto.Uint32(100001), GroupId: proto.Uint32(1), LineId: proto.Uint32(1), Content: proto.String("1_1,2_1,3_1,4_1,5_1")},
	}

	logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p2")}},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p1")},
		},
	}

	manager, err := slot.NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	logic.GameResultManager = manager

	results := logic.CalcGameResult(logic.GameResultManager)

	if len(results) != 1 || results[0].Rule != rule2 {
		t.Errorf("PayResult want only %v", rule2)
		t.Logf("But PayResults has")
		for _, row := range results {
			t.Log(row)
		}
	}

}

func TestSlotLogicBase_CalcGameResult_FixedMatcher_Filter_Rule_MultipleSymbol(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(nil, nil, nil)

	logic.Proto = new(protodata.ProtoData)

	logic.Proto.SymbolAttributes = []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w2"), Multiplier: proto.Int32(1)},
		{Symbol: proto.String("w2"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(2)},
		{Symbol: proto.String("w3"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(3)},
		{Symbol: proto.String("w4"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(4)},
	}

	rule1 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1,w2,w3,w4"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("150")}
	rule2 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("150")}
	logic.Proto.SlotPayTable = []*protodata.SlotPayTable{rule1, rule2}

	logic.Proto.SlotPayLine = []*protodata.SlotPayLine{
		{Id: proto.Uint32(100001), GroupId: proto.Uint32(1), LineId: proto.Uint32(1), Content: proto.String("1_1,2_1,3_1,4_1,5_1")},
	}

	logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p1")}},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p1")},
		},
	}

	manager, err := slot.NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	logic.GameResultManager = manager

	results := logic.CalcGameResult(logic.GameResultManager)

	if len(results) != 1 || results[0].Rule != rule2 {
		t.Errorf("PayResult want only %v", rule2)
		t.Logf("But PayResults has")
		for _, row := range results {
			t.Log(row)
		}
	}

}

func TestSlotLogicBase_CalcGameResult_FixedMatcher_Filter_Rule_MultipleSymbol_Jackpot(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(nil, nil, nil)

	logic.Proto = new(protodata.ProtoData)

	logic.Proto.SymbolAttributes = []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w2"), Multiplier: proto.Int32(1)},
		{Symbol: proto.String("w2"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(2)},
		{Symbol: proto.String("w3"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(3)},
		{Symbol: proto.String("w4"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(4)},
	}

	rule1 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("jackpot"), Count: proto.Int32(5), Payout: proto.String("5")}
	rule2 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1,w2,w3,w4"), PayType: proto.String("chip"), Count: proto.Int32(5), Payout: proto.String("200")}
	logic.Proto.SlotPayTable = []*protodata.SlotPayTable{rule1, rule2}

	logic.Proto.SlotPayLine = []*protodata.SlotPayLine{
		{Id: proto.Uint32(100001), GroupId: proto.Uint32(1), LineId: proto.Uint32(1), Content: proto.String("1_1,2_1,3_1,4_1,5_1")},
	}

	logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("w1")}},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p1")},
		},
	}

	manager, err := slot.NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	logic.GameResultManager = manager

	results := logic.CalcGameResult(logic.GameResultManager)

	if len(results) != 1 || results[0].Rule != rule1 {
		t.Errorf("PayResult want only %v", rule1)
		t.Logf("But PayResults has")
		for _, row := range results {
			t.Log(row)
		}
	}

}

func TestSlotLogicBase_CalcGameAward_FixedMatcher_Filter_OneLineTwoRule(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(nil, nil, nil)

	logic.Proto = new(protodata.ProtoData)

	logic.Proto.SymbolAttributes = []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String(""), Multiplier: proto.Int32(1)},
	}

	rule1 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(5), Payout: proto.String("150")}
	rule2 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("100")}
	rule3 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(3), Payout: proto.String("50")}
	rule4 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("200")}
	rule5 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(3), Payout: proto.String("100")}
	rule6 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1"), PayType: proto.String("chip"), Count: proto.Int32(2), Payout: proto.String("50")}
	logic.Proto.SlotPayTable = []*protodata.SlotPayTable{rule1, rule2, rule3, rule4, rule5, rule6}

	logic.Proto.SlotPayLine = []*protodata.SlotPayLine{
		{Id: proto.Uint32(100001), GroupId: proto.Uint32(1), LineId: proto.Uint32(1), Content: proto.String("1_1,2_1,3_1,4_1,5_1")},
	}

	logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)

	result := slot.GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p2")}},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p1")},
		},
	}

	manager, err := slot.NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	logic.GameResultManager = manager

	results := logic.CalcGameResult(logic.GameResultManager)

	win := logic.CalcGameAward(results, 10, 1)

	if win != 1000 {
		t.Errorf("WinMoney want %d, but %d", 1000, win)
	}
}

func TestSlotLogicBase_CalcGameAward_FixedMatcher_Filter_Rule_MultipleSymbol(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.SlotGameLogicBase = gamelogic.NewSlotGameLogicBase(nil, nil, nil)

	logic.Proto = new(protodata.ProtoData)

	logic.Proto.SymbolAttributes = []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w2"), Multiplier: proto.Int32(1)},
		{Symbol: proto.String("w2"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(2)},
		{Symbol: proto.String("w3"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(3)},
		{Symbol: proto.String("w4"), IsWild: proto.Int32(1), WildEx: proto.String("s1,w1"), Multiplier: proto.Int32(4)},
	}

	rule1 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("w1,w2,w3,w4"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("100")}
	rule2 := &protodata.SlotPayTable{LinkType: proto.String("left"), Symbol: proto.String("p1"), PayType: proto.String("chip"), Count: proto.Int32(4), Payout: proto.String("10")}
	logic.Proto.SlotPayTable = []*protodata.SlotPayTable{rule1, rule2}

	logic.Proto.SlotPayLine = []*protodata.SlotPayLine{
		{Id: proto.Uint32(100001), GroupId: proto.Uint32(1), LineId: proto.Uint32(1), Content: proto.String("1_1,2_1,3_1,4_1,5_1")},
		{Id: proto.Uint32(100002), GroupId: proto.Uint32(1), LineId: proto.Uint32(2), Content: proto.String("1_2,2_3,3_2,4_3,5_2")},
	}

	logic.PaylineMatcher = slot.NewPaylineFixedMatcher(logic.Proto.SymbolAttributes, logic.Proto.SlotPayLine)
	logic.BackgroundMatcher = new(slot.BasicBackgroundMatcher)
	logic.MultiMode = slot.MultiModeEachWildSum

	t.Run("Case 1", func(tt *testing.T) {
		result := slot.GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w2")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w3")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w4")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p1")}},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p5")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p2")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p1")},
			},
		}

		manager, err := slot.NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		logic.GameResultManager = manager

		results := logic.CalcGameResult(logic.GameResultManager)

		// 10 * 100 * (2+3+4)
		win := logic.CalcGameAward(results, 10, 1)

		if win != 9000 {
			tt.Errorf("WinMoney want %d, but %d", 9000, win)
		}
	})

	t.Run("Case 2", func(tt *testing.T) {
		result := slot.GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p2")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")}},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w2")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w3")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w4")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("w1")},
			},
		}

		manager, err := slot.NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		logic.GameResultManager = manager

		results := logic.CalcGameResult(logic.GameResultManager)

		// 10 * 10 * (2 + 4)
		win := logic.CalcGameAward(results, 10, 1)

		if win != 600 {
			tt.Errorf("WinMoney want %d, but %d", 600, win)
		}
	})

}
