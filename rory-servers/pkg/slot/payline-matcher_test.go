package slot

import (
	"reflect"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestBasicPayLineMatcher_MatchAny(t *testing.T) {
	matcher := new(BasicPayLineMatcher)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("b1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("b1"),
		Count:    proto.Int32(int32(3)),
		LinkType: proto.String("any"),
	}

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.PayCount != 3 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 3, paylineres.PayCount)
	}

	if paylineres.WayCount != 1 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
	}

	if paylineres.SymbolName != "b1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{3, 9, 14}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{3, 9, 14}, paylineres.SymbolIDs)
	}

	exceptPos := []Position{
		{3, 1},
		{4, 2},
		{4, 3},
	}

	if !reflect.DeepEqual(paylineres.BlinkPositions, exceptPos) {
		t.Errorf("PaylineResult.BlinkPositions want %v, but %v", exceptPos, paylineres.BlinkPositions)
	}
}

func TestBasicPayLineMatcher_MatchAnyZero(t *testing.T) {
	matcher := new(BasicPayLineMatcher)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p5")},
		},
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("g1"),
		Count:    proto.Int32(0),
		LinkType: proto.String("any"),
	}

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.PayCount != 0 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 0, paylineres.PayCount)
	}

	if paylineres.WayCount != 1 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
	}

	if paylineres.SymbolName != "g1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "g1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32(nil)) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{}, paylineres.SymbolIDs)
	}
}

func TestBasicPayLineMatcher_MatchStraight(t *testing.T) {
	matcher := new(BasicPayLineMatcher)

	t.Run("One Way", func(tt *testing.T) {
		result := GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
			},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		rule := &protodata.SlotPayTable{
			Symbol:   proto.String("w1"),
			Count:    proto.Int32(int32(1)),
			LinkType: proto.String("straight"),
		}

		match, paylineres := matcher.IsMatch(manager, rule)

		if !match {
			tt.Errorf("Match want %v, but %v", true, match)
		}

		if !paylineres.IsMatch {
			tt.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		}

		if paylineres.PayCount != 1 {
			tt.Errorf("PaylineResult.PayCount want %v, but %v", 1, paylineres.PayCount)
		}

		if paylineres.WayCount != 1 {
			tt.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
		}

		if paylineres.SymbolName != "w1" {
			tt.Errorf("PaylineResult.SymbolName want %v, but %v", "w1", paylineres.SymbolName)
		}

		if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 7, 12}) {
			tt.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 7, 12}, paylineres.SymbolIDs)
		}
	})

	t.Run("One Way Two Symbol", func(tt *testing.T) {
		result := GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
			},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		rule := &protodata.SlotPayTable{
			Symbol:   proto.String("p1,w1"),
			Count:    proto.Int32(int32(1)),
			LinkType: proto.String("straight"),
		}

		match, paylineres := matcher.IsMatch(manager, rule)

		if !match {
			tt.Errorf("Match want %v, but %v", true, match)
		}

		if !paylineres.IsMatch {
			tt.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		}

		if paylineres.PayCount != 1 {
			tt.Errorf("PaylineResult.PayCount want %v, but %v", 1, paylineres.PayCount)
		}

		if paylineres.WayCount != 1 {
			tt.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
		}

		if paylineres.SymbolName != "p1,w1" {
			tt.Errorf("PaylineResult.SymbolName want %v, but %v", "p1,w1", paylineres.SymbolName)
		}

		if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 7, 12}) {
			tt.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 7, 12}, paylineres.SymbolIDs)
		}
	})

	t.Run("Two Way", func(tt *testing.T) {
		result := GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
			},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		rule := &protodata.SlotPayTable{
			Symbol:   proto.String("w1"),
			Count:    proto.Int32(int32(2)),
			LinkType: proto.String("straight"),
		}

		match, paylineres := matcher.IsMatch(manager, rule)

		if !match {
			tt.Errorf("Match want %v, but %v", true, match)
		}

		if !paylineres.IsMatch {
			tt.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		}

		if paylineres.PayCount != 2 {
			tt.Errorf("PaylineResult.PayCount want %v, but %v", 2, paylineres.PayCount)
		}

		if paylineres.WayCount != 1 {
			tt.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
		}

		if paylineres.SymbolName != "w1" {
			tt.Errorf("PaylineResult.SymbolName want %v, but %v", "w1", paylineres.SymbolName)
		}

		if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 4, 7, 9, 12, 14}) {
			tt.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 4, 7, 9, 12, 14}, paylineres.SymbolIDs)
		}
	})

	t.Run("Two Way Two Symbol", func(tt *testing.T) {
		result := GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("b1")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w2")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w2")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w2")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
			},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		rule := &protodata.SlotPayTable{
			Symbol:   proto.String("w1,w2"),
			Count:    proto.Int32(int32(2)),
			LinkType: proto.String("straight"),
		}

		match, paylineres := matcher.IsMatch(manager, rule)

		if !match {
			tt.Errorf("Match want %v, but %v", true, match)
		}

		if !paylineres.IsMatch {
			tt.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		}

		if paylineres.PayCount != 2 {
			tt.Errorf("PaylineResult.PayCount want %v, but %v", 2, paylineres.PayCount)
		}

		if paylineres.WayCount != 1 {
			tt.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
		}

		if paylineres.SymbolName != "w1,w2" {
			tt.Errorf("PaylineResult.SymbolName want %v, but %v", "w1,w2", paylineres.SymbolName)
		}

		if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 4, 7, 9, 12, 14}) {
			tt.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 4, 7, 9, 12, 14}, paylineres.SymbolIDs)
		}
	})

	t.Run("Three Way, One Symbol", func(tt *testing.T) {
		result := GameResult{
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
				{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(3)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
				{Id: proto.Uint32(uint32(7)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(8)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
				{Id: proto.Uint32(uint32(12)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(13)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(14)), Symbol: proto.String("w1")},
				{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
			},
		}

		manager, err := NewGameResultManager(5, 3)
		if err != nil {
			t.Error(err)
			return
		}
		manager.GameResult = result

		rule := &protodata.SlotPayTable{
			Symbol:   proto.String("w1"),
			Count:    proto.Int32(int32(3)),
			LinkType: proto.String("straight"),
		}

		match, paylineres := matcher.IsMatch(manager, rule)

		if !match {
			tt.Errorf("Match want %v, but %v", true, match)
		}

		if !paylineres.IsMatch {
			tt.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		}

		if paylineres.PayCount != 3 {
			tt.Errorf("PaylineResult.PayCount want %v, but %v", 3, paylineres.PayCount)
		}

		if paylineres.WayCount != 1 {
			tt.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
		}

		if paylineres.SymbolName != "w1" {
			tt.Errorf("PaylineResult.SymbolName want %v, but %v", "w1", paylineres.SymbolName)
		}

		if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 3, 4, 7, 8, 9, 12, 13, 14}) {
			tt.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 7, 12}, paylineres.SymbolIDs)
		}
	})
}
