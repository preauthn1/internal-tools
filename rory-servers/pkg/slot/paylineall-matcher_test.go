package slot

import (
	"reflect"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestPaylineAllMatcher_ANY(t *testing.T) {
	matcher := NewPaylineAllMatcher(nil)

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

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("b1"),
		Count:    proto.Int32(int32(3)),
		LinkType: proto.String("any"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

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
}

func TestPaylineAllMatcher_ANY_GOLADBALL(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("s1,g1")},
		{Symbol: proto.String("g1"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("jackpot"), GoldballValue: proto.String("7.0")},
		{Symbol: proto.String("g2"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("jackpot"), GoldballValue: proto.String("8.0")},
		{Symbol: proto.String("g3"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("1.0")},
		{Symbol: proto.String("g4"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("2.0")},
		{Symbol: proto.String("g5"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("3.0")},
		{Symbol: proto.String("g6"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("4.0")},
		{Symbol: proto.String("g7"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("5.0")},
		{Symbol: proto.String("g8"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("6.0")},
		{Symbol: proto.String("g9"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("7.0")},
		{Symbol: proto.String("g10"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("8.0")},
		{Symbol: proto.String("g11"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("10.0")},
		{Symbol: proto.String("g12"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("14.0")},
		{Symbol: proto.String("g13"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("16.0")},
		{Symbol: proto.String("g14"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("18.0")},
		{Symbol: proto.String("g15"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("20.0")},
		{Symbol: proto.String("g16"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("24.0")},
		{Symbol: proto.String("g17"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("30.0")},
		{Symbol: proto.String("g18"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("100.0")},
	}

	matcher := NewPaylineAllMatcher(attrs)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("g6")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("g3")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("g14")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("g5")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("g8")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p7")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p4")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("g18")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("g10")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("g8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p2")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("g1"),
		Count:    proto.Int32(int32(8)),
		LinkType: proto.String("any"),
		PayType:  proto.String("lightning"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.PayCount != 8 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 8, paylineres.PayCount)
	}

	if paylineres.WayCount != 1 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
	}

	if paylineres.SymbolName != "g1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "g1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{2, 3, 5, 6, 7, 11, 12, 13}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{2, 3, 5, 6, 7, 11, 12, 13}, paylineres.SymbolIDs)
	}

}

func TestPaylineAllMatcher_ANY_JACKPOT_GOLADBALL(t *testing.T) {

	attrs := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("s1,g1")},
		{Symbol: proto.String("g1"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("jackpot"), GoldballValue: proto.String("7.0")},
		{Symbol: proto.String("g2"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("jackpot"), GoldballValue: proto.String("8.0")},
		{Symbol: proto.String("g3"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("1.0")},
		{Symbol: proto.String("g4"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("2.0")},
		{Symbol: proto.String("g5"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("3.0")},
		{Symbol: proto.String("g6"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("4.0")},
		{Symbol: proto.String("g7"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("5.0")},
		{Symbol: proto.String("g8"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("6.0")},
		{Symbol: proto.String("g9"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("7.0")},
		{Symbol: proto.String("g10"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("8.0")},
		{Symbol: proto.String("g11"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("10.0")},
		{Symbol: proto.String("g12"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("14.0")},
		{Symbol: proto.String("g13"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("16.0")},
		{Symbol: proto.String("g14"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("18.0")},
		{Symbol: proto.String("g15"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("20.0")},
		{Symbol: proto.String("g16"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("24.0")},
		{Symbol: proto.String("g17"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("30.0")},
		{Symbol: proto.String("g18"), IsWild: proto.Int32(0), IsGoldBall: proto.Int32(1), GoldballType: proto.String("chipXrate"), GoldballValue: proto.String("100.0")},
	}

	matcher := NewPaylineAllMatcher(attrs)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("g1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("g6")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("g3")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("g4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("g14")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("g5")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("g8")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("g7")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("g9")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("g2")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("g18")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("g10")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("g8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("g3")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("g4")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("g1"),
		Count:    proto.Int32(int32(15)),
		LinkType: proto.String("any"),
		PayType:  proto.String("jackpot"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.PayCount != 15 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 15, paylineres.PayCount)
	}

	if paylineres.WayCount != 1 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
	}

	if paylineres.SymbolName != "g1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "g1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, paylineres.SymbolIDs)
	}
}

func TestPaylineAllMatcher_IsMatch_Left_NoWild(t *testing.T) {

	matcher := NewPaylineAllMatcher(nil)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("p1"),
		Count:    proto.Int32(int32(4)),
		LinkType: proto.String("left"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	t.Logf("%+v\n", manager)

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.HasWild {
		t.Errorf("PaylineResult.HasWild want %v, but %v", false, paylineres.HasWild)
	}

	if paylineres.PayCount != 4 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 4, paylineres.PayCount)
	}

	if paylineres.WayCount != 4 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 4, paylineres.WayCount)
	}

	if paylineres.SymbolName != "p1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{1, 2, 3, 7, 9, 14}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{1, 2, 3, 3, 9, 14}, paylineres.SymbolIDs)
	}
}

func TestPaylineAllMatcher_IsMatch_Left_HasWild_NoExt(t *testing.T) {

	wilds := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("")},
	}

	matcher := NewPaylineAllMatcher(wilds)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("p1"),
		Count:    proto.Int32(int32(4)),
		LinkType: proto.String("left"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if !paylineres.HasWild {
		t.Errorf("PaylineResult.HasWild want %v, but %v", true, paylineres.HasWild)
	}

	if paylineres.PayCount != 4 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 4, paylineres.PayCount)
	}

	if paylineres.WayCount != 2 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 2, paylineres.WayCount)
	}

	if paylineres.SymbolName != "p1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p1", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{1, 2, 3, 9, 14}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{1, 2, 3, 9, 14}, paylineres.SymbolIDs)
	}
}

func TestPaylineAllMatcher_IsMatch_Left_HasWild_HasExt(t *testing.T) {

	wilds := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("p1,p2,p3")},
	}

	matcher := NewPaylineAllMatcher(wilds)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("p1"),
		Count:    proto.Int32(int32(4)),
		LinkType: proto.String("left"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if match {
		t.Errorf("Match want %v, but %v", false, match)
	}

	if paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", false, paylineres.IsMatch)
	}
}

func TestPaylineAllMatcher_IsMatch_Right_NoWild(t *testing.T) {

	matcher := NewPaylineAllMatcher(nil)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p2")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p4")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("p4"),
		Count:    proto.Int32(int32(3)),
		LinkType: proto.String("right"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
	}

	if paylineres.HasWild {
		t.Errorf("PaylineResult.HasWild want %v, but %v", false, paylineres.HasWild)
	}

	if paylineres.PayCount != 3 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 3, paylineres.PayCount)
	}

	if paylineres.WayCount != 1 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 1, paylineres.WayCount)
	}

	if paylineres.SymbolName != "p4" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p4", paylineres.SymbolName)
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{4, 8, 10}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{4, 8, 10}, paylineres.SymbolIDs)
	}
}

func TestPaylineAllMatcher_IsMatch_Right_HasWild_NoExt(t *testing.T) {

	wilds := []*protodata.SymbolAttributes{
		{Symbol: proto.String("w1"), IsWild: proto.Int32(1), WildEx: proto.String("")},
	}

	matcher := NewPaylineAllMatcher(wilds)

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(1)), Symbol: proto.String("p5")},
			{Id: proto.Uint32(uint32(2)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(3)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(4)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(5)), Symbol: proto.String("p1")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(6)), Symbol: proto.String("p6")},
			{Id: proto.Uint32(uint32(7)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(8)), Symbol: proto.String("p4")},
			{Id: proto.Uint32(uint32(9)), Symbol: proto.String("w1")},
			{Id: proto.Uint32(uint32(10)), Symbol: proto.String("p3")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(uint32(11)), Symbol: proto.String("p3")},
			{Id: proto.Uint32(uint32(12)), Symbol: proto.String("p12")},
			{Id: proto.Uint32(uint32(13)), Symbol: proto.String("p8")},
			{Id: proto.Uint32(uint32(14)), Symbol: proto.String("p1")},
			{Id: proto.Uint32(uint32(15)), Symbol: proto.String("p7")},
		},
	}

	rule := &protodata.SlotPayTable{
		Symbol:   proto.String("p1"),
		Count:    proto.Int32(int32(4)),
		LinkType: proto.String("right"),
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}
	manager.GameResult = result

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
		return
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		return
	}

	if !paylineres.HasWild {
		t.Errorf("PaylineResult.HasWild want %v, but %v", true, paylineres.HasWild)
		return
	}

	if paylineres.PayCount != 4 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 4, paylineres.PayCount)
		return
	}

	if paylineres.WayCount != 3 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 3, paylineres.WayCount)
		return
	}

	if paylineres.SymbolName != "p1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p1", paylineres.SymbolName)
		return
	}

	if !reflect.DeepEqual(paylineres.SymbolIDs, []int32{3, 4, 5, 7, 9, 14}) {
		t.Errorf("PaylineResult.SymbolIDs want %v, but %v", []int32{3, 4, 5, 7, 9, 14}, paylineres.SymbolIDs)
		return
	}
}

func TestPaylineAllMatcher_IsMatch_Left_Irregular(t *testing.T) {

	matcher := NewPaylineAllMatcher(nil)
	rule := &protodata.SlotPayTable{Symbol: proto.String("p1"), Count: proto.Int32(4), LinkType: proto.String("left")}
	manager, err := NewGameResultManager(5, 6)
	if err != nil {
		t.Error(err)
		return
	}
	{
		manager.GameResult = GameResult{
			[]*protodata.SlotSymbolDistribution{
				nil,
				nil,
				{Symbol: proto.String("p1")},
				nil,
				nil,
			},
			[]*protodata.SlotSymbolDistribution{
				nil,
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p8")},
				{Symbol: proto.String("p1")},
				nil,
			},
			[]*protodata.SlotSymbolDistribution{
				nil,
				{Symbol: proto.String("p12")},
				{Symbol: proto.String("p13")},
				{Symbol: proto.String("p14")},
				nil,
			},
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p17")},
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p19")},
				{Symbol: proto.String("p20")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p22")},
				{Symbol: proto.String("p23")},
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p25")},
			},
			[]*protodata.SlotSymbolDistribution{
				{Symbol: proto.String("p26")},
				{Symbol: proto.String("p27")},
				{Symbol: proto.String("p1")},
				{Symbol: proto.String("p29")},
				{Symbol: proto.String("p30")},
			},
		}
	}

	match, paylineres := matcher.IsMatch(manager, rule)

	if !match {
		t.Errorf("Match want %v, but %v", true, match)
		return
	}

	if !paylineres.IsMatch {
		t.Errorf("PaylineResult.IsMatch want %v, but %v", true, paylineres.IsMatch)
		return
	}

	if paylineres.HasWild {
		t.Errorf("PaylineResult.HasWild want %v, but %v", false, paylineres.HasWild)
		return
	}

	if paylineres.PayCount != 4 {
		t.Errorf("PaylineResult.PayCount want %v, but %v", 4, paylineres.PayCount)
		return
	}

	if paylineres.WayCount != 12 {
		t.Errorf("PaylineResult.WayCount want %v, but %v", 12, paylineres.WayCount)
		return
	}

	if paylineres.WildWayCount != 0 {
		t.Errorf("PaylineResult.WildWayCount want %v, but %v", 0, paylineres.WildWayCount)
		return
	}

	if paylineres.NormalWayCount != 12 {
		t.Errorf("PaylineResult.NormalWayCount want %v, but %v", 12, paylineres.NormalWayCount)
		return
	}

	if paylineres.SymbolName != "p1" {
		t.Errorf("PaylineResult.SymbolName want %v, but %v", "p1", paylineres.SymbolName)
		return
	}
}
