package pachislot

import (
	"fmt"
	"reflect"
	"testing"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func TestSymbolGroup(t *testing.T) {
	t.Run("Match", func(tt *testing.T) {
		sg := SymbolGroup{
			{pb.Symbol_Bar, pb.Symbol_Bar, pb.Symbol_Bar},
			{pb.Symbol_Bar, pb.Symbol_Bell, pb.Symbol_Bar},
			{pb.Symbol_Bar, pb.Symbol_Bar, pb.Symbol_Bar},
		}
		t.Log(sg)
		tt.Run("No Line", func(ttt *testing.T) {
			match := sg.Match(nil, pb.Symbol_Bar)
			if match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", false, match)
			}
		})
		tt.Run("Line 1-1-1", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 1},
				{Reel: 3, Row: 1},
			}}

			match := sg.Match(line, pb.Symbol_Bar)
			if !match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", true, match)
			}
		})
		tt.Run("Line 2-2-2", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 2},
				{Reel: 2, Row: 2},
				{Reel: 3, Row: 2},
			}}

			match := sg.Match(line, pb.Symbol_Bar)
			if match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", false, match)
			}
		})
		tt.Run("Line 3-3-3", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 3},
				{Reel: 2, Row: 3},
				{Reel: 3, Row: 3},
			}}

			match := sg.Match(line, pb.Symbol_Bar)
			if !match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", true, match)
			}
		})
		tt.Run("Line 1-2-3", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 2},
				{Reel: 3, Row: 3},
			}}

			match := sg.Match(line, pb.Symbol_Bar)
			if match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", false, match)
			}
		})
		tt.Run("Line 3-2-1", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 3},
				{Reel: 2, Row: 2},
				{Reel: 3, Row: 1},
			}}

			match := sg.Match(line, pb.Symbol_Bar)
			if match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", false, match)
			}
		})
		tt.Run("Diff Symbol", func(ttt *testing.T) {
			line := &pb.PayLine{Positions: []*pb.Position{
				{Reel: 1, Row: 1},
				{Reel: 2, Row: 3},
				{Reel: 3, Row: 3},
			}}

			match := sg.Match(line, pb.Symbol_Bell)
			if match {
				ttt.Errorf("SymbolGroup.Match() want %t, but %t", false, match)
			}
		})
	})
	t.Run("Symbol", func(t *testing.T) {
		sg := SymbolGroup{
			{pb.Symbol_Bar, pb.Symbol_Bell, pb.Symbol_WaterMelon},
			{pb.Symbol_BlueSeven, pb.Symbol_RedSeven, pb.Symbol_Replay},
			{pb.Symbol_Berry, pb.Symbol_Bar, pb.Symbol_Bell},
		}
		t.Run("1-1", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 1, Row: 1})
			if s != pb.Symbol_Bar {
				t.Errorf("Symbol() 1-1 want %v, but %v", pb.Symbol_Bar, s)
			}
		})
		t.Run("1-2", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 1, Row: 2})
			if s != pb.Symbol_Bell {
				t.Errorf("Symbol() 1-2 want %v, but %v", pb.Symbol_Bell, s)
			}
		})
		t.Run("1-3", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 1, Row: 3})
			if s != pb.Symbol_WaterMelon {
				t.Errorf("Symbol() 1-3 want %v, but %v", pb.Symbol_WaterMelon, s)
			}
		})
		t.Run("2-1", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 2, Row: 1})
			if s != pb.Symbol_BlueSeven {
				t.Errorf("Symbol() 2-1 want %v, but %v", pb.Symbol_BlueSeven, s)
			}
		})
		t.Run("2-2", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 2, Row: 2})
			if s != pb.Symbol_RedSeven {
				t.Errorf("Symbol() 2-2 want %v, but %v", pb.Symbol_RedSeven, s)
			}
		})
		t.Run("2-3", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 2, Row: 3})
			if s != pb.Symbol_Replay {
				t.Errorf("Symbol() 2-3 want %v, but %v", pb.Symbol_Replay, s)
			}
		})
		t.Run("3-1", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 3, Row: 1})
			if s != pb.Symbol_Berry {
				t.Errorf("Symbol() 3-1 want %v, but %v", pb.Symbol_Berry, s)
			}
		})
		t.Run("3-2", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 3, Row: 2})
			if s != pb.Symbol_Bar {
				t.Errorf("Symbol() 3-2 want %v, but %v", pb.Symbol_Bar, s)
			}
		})
		t.Run("3-3", func(t *testing.T) {
			s := sg.Symbol(&pb.Position{Reel: 3, Row: 3})
			if s != pb.Symbol_Bell {
				t.Errorf("Symbol() 3-3 want %v, but %v", pb.Symbol_Bell, s)
			}
		})
	})
}

func TestReelStripGroup(t *testing.T) {
	t.Run("Valid No Reel1", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}
		g := NewReelStripGroup(ReelStrips)
		valid := g.Valid()
		if valid {
			tt.Errorf("ReelStripGroup.Valid() want %t, but %t", false, valid)
		}
	})
	t.Run("Valid No Reel2", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}
		g := NewReelStripGroup(ReelStrips)
		valid := g.Valid()
		if valid {
			tt.Errorf("ReelStripGroup.Valid() want %t, but %t", false, valid)
		}
	})
	t.Run("Valid No Reel3", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},
		}
		g := NewReelStripGroup(ReelStrips)
		valid := g.Valid()
		if valid {
			tt.Errorf("ReelStripGroup.Valid() want %t, but %t", false, valid)
		}
	})
	t.Run("Valid", func(tt *testing.T) {
		g := NewReelStripGroup(DefaultReelStrips)
		valid := g.Valid()
		if !valid {
			tt.Errorf("ReelStripGroup.Valid() want %t, but %t", true, valid)
		}
	})
	t.Run("Valid Symbols Less 3", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}
		g := NewReelStripGroup(ReelStrips)
		valid := g.Valid()
		if valid {
			tt.Errorf("ReelStripGroup.Valid() want %t, but %t", false, valid)
		}
	})
	t.Run("SymbolIndex", func(tt *testing.T) {
		g := NewReelStripGroup(DefaultReelStrips)

		exp := map[int][]int{
			1: {0, 4, 9, 14, 19},
			2: {1, 5, 10, 14, 18},
			3: {0, 1, 5, 8, 13, 16},
		}

		idxMap := g.SymbolIndex(pb.Symbol_Bell)

		if !reflect.DeepEqual(exp, idxMap) {
			tt.Errorf("ReelStripGroup.SymbolIndex() want %v, but %v", exp, idxMap)
		}
	})
}

func TestReelStripMap(t *testing.T) {
	m := ReelStripMap{2, 0, 20}

	idx1, idx2, idx3 := m.StopIndex()
	if idx1 != 2 {
		t.Errorf("idx1 want %d, but %d", 2, idx1)
	}
	if idx2 != 0 {
		t.Errorf("idx2 want %d, but %d", 0, idx2)
	}
	if idx3 != 20 {
		t.Errorf("idx3 want %d, but %d", 20, idx3)
	}

	g := NewReelStripGroup(DefaultProto.ReelStrips)
	sg := m.SymbolGroup(g)

	if sg[0][0] != pb.Symbol_Replay {
		t.Errorf("Reel: 1, Row: 1 want %v, but %v", pb.Symbol_Replay, sg[0][0])
	}
	if sg[0][1] != pb.Symbol_RedSeven {
		t.Errorf("Reel: 1, Row: 2 want %v, but %v", pb.Symbol_RedSeven, sg[0][1])
	}
	if sg[0][2] != pb.Symbol_Berry {
		t.Errorf("Reel: 1, Row: 3 want %v, but %v", pb.Symbol_Berry, sg[0][2])
	}
	if sg[1][0] != pb.Symbol_WaterMelon {
		t.Errorf("Reel: 2, Row: 1 want %v, but %v", pb.Symbol_WaterMelon, sg[1][0])
	}
	if sg[1][1] != pb.Symbol_Berry {
		t.Errorf("Reel: 2, Row: 2 want %v, but %v", pb.Symbol_Berry, sg[1][1])
	}
	if sg[1][2] != pb.Symbol_Bell {
		t.Errorf("Reel: 2, Row: 3 want %v, but %v", pb.Symbol_Bell, sg[1][2])
	}
	if sg[2][0] != pb.Symbol_Bar {
		t.Errorf("Reel: 3, Row: 1 want %v, but %v", pb.Symbol_Bar, sg[2][0])
	}
	if sg[2][1] != pb.Symbol_Replay {
		t.Errorf("Reel: 3, Row: 2 want %v, but %v", pb.Symbol_Replay, sg[2][1])
	}
	if sg[2][2] != pb.Symbol_Bell {
		t.Errorf("Reel: 3, Row: 3 want %v, but %v", pb.Symbol_Bell, sg[2][2])
	}
}

func TestBellLinkResultSet(t *testing.T) {
	r := new(BellLinkResultSet)

	if r.Type() != pb.ResultType_BellLink {
		t.Errorf("BellLinkResult Type() want %v, but %v", pb.ResultType_BellLink, r.Type())
	}

	if r.Symbol() != pb.Symbol_Bell {
		t.Errorf("BellLinkResult Symbol() want %v, but %v", pb.Symbol_Bell, r.Symbol())
	}

	t.Run("Match", func(tt *testing.T) {
		sg := SymbolGroup{
			{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Berry},
			{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Berry},
			{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Berry},
		}
		match := r.Match(DefaultPayLine, sg)
		if !match {
			tt.Errorf("Match() want %t, but %t", true, match)
		}
	})

	t.Run("Not Match", func(tt *testing.T) {
		sg := SymbolGroup{
			{pb.Symbol_Bar, pb.Symbol_Bar, pb.Symbol_Berry},
			{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Berry},
			{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Bar},
		}
		match := r.Match(DefaultPayLine, sg)
		if match {
			tt.Errorf("Match() want %t, but %t, %v", false, match, sg)
		}
	})

	t.Run("Over Range", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Replay},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Replay},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Replay},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Replay},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_Bell},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Replay},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Replay},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}
		g := NewReelStripGroup(ReelStrips)
		expLine := &pb.PayLine{Positions: []*pb.Position{{Reel: 1, Row: 1}, {Reel: 2, Row: 2}, {Reel: 3, Row: 3}}}

		lines := []*pb.PayLine{expLine}

		pos, m, err := r.Generete(lines, g)

		expMap := ReelStripMap{0, 14, 20}

		if !reflect.DeepEqual(expLine.Positions, pos) {
			tt.Errorf("Generete() pos want %v, but %v", expLine, pos)
		}
		if !reflect.DeepEqual(expMap, m) {
			tt.Errorf("Generete() ReelStripMap want %v, but %v", expMap, m)
		}
		if err != nil {
			tt.Errorf("Generete() err want %v, but %v", nil, err)
		}
	})

	t.Run("One Reel Has no Passibale Symbol", func(tt *testing.T) {
		var ReelStrips []*pb.ReelStripSymbol = []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Replay},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Replay},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Replay},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Replay},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_Bell},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Berry},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Replay},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Replay},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}
		g := NewReelStripGroup(ReelStrips)
		_, m, err := r.Generete(DefaultPayLine, g)

		expMap := ReelStripMap{}

		if !reflect.DeepEqual(expMap, m) {
			tt.Errorf("Generete() ReelStripMap want %v, but %v", expMap, m)
		}
		if err != ErrResultSetNotGeneratable {
			tt.Errorf("Generete() err want %v, but %v", ErrResultSetNotGeneratable, err)
		}
	})

	t.Run("PayLine Not Found", func(tt *testing.T) {
		PayLine := []*pb.PayLine{}
		g := NewReelStripGroup(DefaultReelStrips)
		_, m, err := r.Generete(PayLine, g)

		expMap := ReelStripMap{}

		if !reflect.DeepEqual(expMap, m) {
			tt.Errorf("Generete() ReelStripMap want %v, but %v", expMap, m)
		}
		if err != ErrPayLineNotFound {
			tt.Errorf("Generete() err want %v, but %v", ErrPayLineNotFound, err)
		}
	})
}

func TestNoneResultSet(t *testing.T) {
	r := new(NoneResultSet)

	if r.Type() != pb.ResultType_None {
		t.Errorf("NoneResultSet Type() want %v, but %v", pb.ResultType_None, r.Type())
	}

	if r.Symbol() != pb.Symbol(-1) {
		t.Errorf("NoneResultSet Symbol() want %v, but %v", pb.Symbol(-1), r.Symbol())
	}

	g := NewReelStripGroup(DefaultProto.ReelStrips)
	line, m, err := r.Generete(DefaultProto.PayLine, g)
	t.Log(line, m, err)
	t.Log(m.SymbolGroup(g))
}

func TestBerryLeftAnyResultSet(t *testing.T) {
	r := new(BerryLeftAnyResultSet)

	if r.Type() != pb.ResultType_BerryLeftAny {
		t.Errorf("BerryLeftAnyResultSet Type() want %v, but %v", pb.ResultType_BerryLeftAny, r.Type())
	}

	if r.Symbol() != pb.Symbol_Berry {
		t.Errorf("BerryLeftAnyResultSet Symbol() want %v, but %v", pb.Symbol_Berry, r.Symbol())
	}

	g := NewReelStripGroup(DefaultProto.ReelStrips)
	pos, m, err := r.Generete(DefaultProto.PayLine, g)

	if err != nil {
		t.Errorf("Generete() err want %v, but %v", nil, err)
	}

	sg := m.SymbolGroup(g)
	t.Log(pos, sg)
}

func TestDoubleSevenResultSet(t *testing.T) {
	r := new(DoubleSevenLinkResultSet)

	if r.Type() != pb.ResultType_DoubleSevenLink {
		t.Errorf("DoubleSevenLinkResultSet Type() want %v, but %v", pb.ResultType_DoubleSevenLink, r.Type())
	}

	if r.Symbol() != pb.Symbol_BlueSeven {
		t.Errorf("DoubleSevenLinkResultSet Symbol() want %v, but %v", pb.Symbol_BlueSeven, r.Symbol())
	}

	g := NewReelStripGroup(DefaultProto.ReelStrips)
	_, m, err := r.Generete(DefaultProto.PayLine, g)

	if err != nil {
		t.Errorf("Generete() err want %v, but %v", nil, err)
	}

	expsg1 := SymbolGroup{
		{pb.Symbol_BlueSeven, pb.Symbol_BlueSeven, pb.Symbol_Bell},
		{pb.Symbol_Bell, pb.Symbol_BlueSeven, pb.Symbol_Replay},
		{pb.Symbol_Bell, pb.Symbol_BlueSeven, pb.Symbol_BlueSeven},
	}
	expsg2 := SymbolGroup{
		{pb.Symbol_WaterMelon, pb.Symbol_BlueSeven, pb.Symbol_BlueSeven},
		{pb.Symbol_Bell, pb.Symbol_BlueSeven, pb.Symbol_Replay},
		{pb.Symbol_BlueSeven, pb.Symbol_BlueSeven, pb.Symbol_WaterMelon},
	}

	sg := m.SymbolGroup(g)
	t.Log(sg)

	if !reflect.DeepEqual(expsg1, sg) && !reflect.DeepEqual(expsg2, sg) {
		t.Errorf("DoubleSevenLinkResultSet Generete() SymbolGroup want %v or %v, but %v", expsg1, expsg2, sg)
	}

	t.Run("Reel1 No Continue B7", func(t *testing.T) {
		strips := []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_Replay},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}

		g := NewReelStripGroup(strips)

		_, _, err := r.Generete(DefaultProto.PayLine, g)
		if err != ErrResultSetNotGeneratable {
			t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
		}
	})
	t.Run("Reel2 No B7", func(t *testing.T) {
		strips := []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_Bell},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
		}

		g := NewReelStripGroup(strips)

		_, _, err := r.Generete(DefaultProto.PayLine, g)
		if err != ErrResultSetNotGeneratable {
			t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
		}
	})
	t.Run("Reel3 No Continue B7", func(t *testing.T) {
		strips := []*pb.ReelStripSymbol{
			{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
			{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
			{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
			{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
			{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
			{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
			{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
			{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
			{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
			{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
			{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
			{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
			{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
			{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
			{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
			{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
			{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
			{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
			{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
			{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
			{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
			{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
			{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
			{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
			{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
			{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

			{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
			{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
			{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
			{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
			{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
			{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
			{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
			{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
			{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
			{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_Bar},
			{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
			{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
			{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
			{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
			{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
			{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
			{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
			{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
			{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
			{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_BlueSeven},
		}

		g := NewReelStripGroup(strips)

		_, _, err := r.Generete(DefaultProto.PayLine, g)
		if err != ErrResultSetNotGeneratable {
			t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
		}
	})
}

func TestGenerateResultSet(t *testing.T) {

	g := NewReelStripGroup(DefaultProto.ReelStrips)

	t.Run("Empty", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_Empty)
		if set != nil {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_Empty, nil, set)
		}
	})
	t.Run("None", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_None)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_None, expType, reflect.TypeOf(set))
		}
	})
	t.Run("Feature None", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_FeatureNone)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_FeatureNone, expType, reflect.TypeOf(set))
		}
	})
	t.Run("JAC None", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_JACNone)
		expType := reflect.TypeOf(new(JACNoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_JACNone, expType, reflect.TypeOf(set))
		}
		if set.Symbol() != pb.Symbol(-1) {
			t.Errorf("Symbol() by %v, want %v, but %v", pb.ResultType_JACNone, pb.Symbol(-1), set.Symbol())
		}
		t.Run("Match", func(t *testing.T) {
			sg := SymbolGroup{
				{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Replay},
				{pb.Symbol_WaterMelon, pb.Symbol_Replay, pb.Symbol_WaterMelon},
				{pb.Symbol_Bell, pb.Symbol_BlueSeven, pb.Symbol_RedSeven},
			}
			match := set.Match(DefaultProto.PayLine, sg)
			if !match {
				t.Errorf("Match() by %v, want %t, but %t", pb.ResultType_JACNone, true, match)
			}
		})
		t.Run("Not Match", func(t *testing.T) {
			sg := SymbolGroup{
				{pb.Symbol_Bell, pb.Symbol_Bar, pb.Symbol_Replay},
				{pb.Symbol_WaterMelon, pb.Symbol_Replay, pb.Symbol_WaterMelon},
				{pb.Symbol_WaterMelon, pb.Symbol_BlueSeven, pb.Symbol_RedSeven},
			}
			match := set.Match(DefaultProto.PayLine, sg)
			if match {
				t.Errorf("Match() by %v, want %t, but %t", pb.ResultType_JACNone, false, match)
			}
		})
		t.Run("Generate", func(t *testing.T) {
			_, m, err := set.Generete(DefaultProto.PayLine, g)
			if err != nil {
				t.Errorf("err want %v, but %v", nil, err)
			}
			t.Log(m)

			t.Run("No PayLine", func(t *testing.T) {
				_, _, err := set.Generete([]*pb.PayLine{}, g)
				if err != ErrPayLineNotFound {
					t.Errorf("err want %v, but %v", ErrPayLineNotFound, err)
				}
			})
			t.Run("Reel1 No Replay", func(t *testing.T) {
				strips := []*pb.ReelStripSymbol{
					{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
					{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
					{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
					{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
					{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Bell},
					{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
					{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
					{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Bell},
					{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
					{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Bell},
					{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
					{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
					{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
					{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
					{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
					{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
					{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
					{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
					{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
					{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
					{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
					{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
					{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
					{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
					{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
					{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
					{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
					{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
					{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
					{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
					{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
					{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
					{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
					{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
					{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
					{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
					{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
				}

				g := NewReelStripGroup(strips)

				_, _, err := set.Generete(DefaultProto.PayLine, g)
				if err != ErrResultSetNotGeneratable {
					t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
				}
			})
			t.Run("Reel2 No Replay", func(t *testing.T) {
				strips := []*pb.ReelStripSymbol{
					{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
					{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
					{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
					{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
					{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
					{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
					{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
					{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
					{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
					{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
					{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
					{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
					{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
					{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
					{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
					{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Bell},
					{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Bell},
					{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
					{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
					{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Bell},
					{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
					{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
					{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Bell},
					{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
					{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
					{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Bell},
					{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_Bell},
					{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
					{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
					{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
					{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_Bell},
					{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
					{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_Bell},
					{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
					{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
					{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_Bell},
					{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
					{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
					{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
				}

				g := NewReelStripGroup(strips)

				_, _, err := set.Generete(DefaultProto.PayLine, g)
				if err != ErrResultSetNotGeneratable {
					t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
				}
			})
			t.Run("Reel3 No Bell", func(t *testing.T) {
				strips := []*pb.ReelStripSymbol{
					{ID: 60010001, GroupID: 6001, Reel: 1, Row: 1, Symbol: pb.Symbol_Bell},
					{ID: 60010002, GroupID: 6001, Reel: 1, Row: 2, Symbol: pb.Symbol_Replay},
					{ID: 60010003, GroupID: 6001, Reel: 1, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010004, GroupID: 6001, Reel: 1, Row: 4, Symbol: pb.Symbol_Berry},
					{ID: 60010005, GroupID: 6001, Reel: 1, Row: 5, Symbol: pb.Symbol_Bell},
					{ID: 60010006, GroupID: 6001, Reel: 1, Row: 6, Symbol: pb.Symbol_Replay},
					{ID: 60010007, GroupID: 6001, Reel: 1, Row: 7, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010008, GroupID: 6001, Reel: 1, Row: 8, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010009, GroupID: 6001, Reel: 1, Row: 9, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010010, GroupID: 6001, Reel: 1, Row: 10, Symbol: pb.Symbol_Bell},
					{ID: 60010011, GroupID: 6001, Reel: 1, Row: 11, Symbol: pb.Symbol_Replay},
					{ID: 60010012, GroupID: 6001, Reel: 1, Row: 12, Symbol: pb.Symbol_Bar},
					{ID: 60010013, GroupID: 6001, Reel: 1, Row: 13, Symbol: pb.Symbol_Berry},
					{ID: 60010014, GroupID: 6001, Reel: 1, Row: 14, Symbol: pb.Symbol_Replay},
					{ID: 60010015, GroupID: 6001, Reel: 1, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010016, GroupID: 6001, Reel: 1, Row: 16, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010017, GroupID: 6001, Reel: 1, Row: 17, Symbol: pb.Symbol_Bar},
					{ID: 60010018, GroupID: 6001, Reel: 1, Row: 18, Symbol: pb.Symbol_Replay},
					{ID: 60010019, GroupID: 6001, Reel: 1, Row: 19, Symbol: pb.Symbol_RedSeven},
					{ID: 60010020, GroupID: 6001, Reel: 1, Row: 20, Symbol: pb.Symbol_Bell},
					{ID: 60010021, GroupID: 6001, Reel: 1, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010022, GroupID: 6001, Reel: 2, Row: 1, Symbol: pb.Symbol_Berry},
					{ID: 60010023, GroupID: 6001, Reel: 2, Row: 2, Symbol: pb.Symbol_Bell},
					{ID: 60010024, GroupID: 6001, Reel: 2, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010025, GroupID: 6001, Reel: 2, Row: 4, Symbol: pb.Symbol_Replay},
					{ID: 60010026, GroupID: 6001, Reel: 2, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010027, GroupID: 6001, Reel: 2, Row: 6, Symbol: pb.Symbol_Bell},
					{ID: 60010028, GroupID: 6001, Reel: 2, Row: 7, Symbol: pb.Symbol_Replay},
					{ID: 60010029, GroupID: 6001, Reel: 2, Row: 8, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010030, GroupID: 6001, Reel: 2, Row: 9, Symbol: pb.Symbol_Berry},
					{ID: 60010031, GroupID: 6001, Reel: 2, Row: 10, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010032, GroupID: 6001, Reel: 2, Row: 11, Symbol: pb.Symbol_Bell},
					{ID: 60010033, GroupID: 6001, Reel: 2, Row: 12, Symbol: pb.Symbol_Replay},
					{ID: 60010034, GroupID: 6001, Reel: 2, Row: 13, Symbol: pb.Symbol_Bar},
					{ID: 60010035, GroupID: 6001, Reel: 2, Row: 14, Symbol: pb.Symbol_Berry},
					{ID: 60010036, GroupID: 6001, Reel: 2, Row: 15, Symbol: pb.Symbol_Bell},
					{ID: 60010037, GroupID: 6001, Reel: 2, Row: 16, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010038, GroupID: 6001, Reel: 2, Row: 17, Symbol: pb.Symbol_Replay},
					{ID: 60010039, GroupID: 6001, Reel: 2, Row: 18, Symbol: pb.Symbol_Berry},
					{ID: 60010040, GroupID: 6001, Reel: 2, Row: 19, Symbol: pb.Symbol_Bell},
					{ID: 60010041, GroupID: 6001, Reel: 2, Row: 20, Symbol: pb.Symbol_Replay},
					{ID: 60010042, GroupID: 6001, Reel: 2, Row: 21, Symbol: pb.Symbol_WaterMelon},

					{ID: 60010043, GroupID: 6001, Reel: 3, Row: 1, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010044, GroupID: 6001, Reel: 3, Row: 2, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010045, GroupID: 6001, Reel: 3, Row: 3, Symbol: pb.Symbol_RedSeven},
					{ID: 60010046, GroupID: 6001, Reel: 3, Row: 4, Symbol: pb.Symbol_Replay},
					{ID: 60010047, GroupID: 6001, Reel: 3, Row: 5, Symbol: pb.Symbol_Berry},
					{ID: 60010048, GroupID: 6001, Reel: 3, Row: 6, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010049, GroupID: 6001, Reel: 3, Row: 7, Symbol: pb.Symbol_Bar},
					{ID: 60010050, GroupID: 6001, Reel: 3, Row: 8, Symbol: pb.Symbol_Replay},
					{ID: 60010051, GroupID: 6001, Reel: 3, Row: 9, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010052, GroupID: 6001, Reel: 3, Row: 10, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010053, GroupID: 6001, Reel: 3, Row: 11, Symbol: pb.Symbol_BlueSeven},
					{ID: 60010054, GroupID: 6001, Reel: 3, Row: 12, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010055, GroupID: 6001, Reel: 3, Row: 13, Symbol: pb.Symbol_Replay},
					{ID: 60010056, GroupID: 6001, Reel: 3, Row: 14, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010057, GroupID: 6001, Reel: 3, Row: 15, Symbol: pb.Symbol_Bar},
					{ID: 60010058, GroupID: 6001, Reel: 3, Row: 16, Symbol: pb.Symbol_Replay},
					{ID: 60010059, GroupID: 6001, Reel: 3, Row: 17, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010060, GroupID: 6001, Reel: 3, Row: 18, Symbol: pb.Symbol_WaterMelon},
					{ID: 60010061, GroupID: 6001, Reel: 3, Row: 19, Symbol: pb.Symbol_Berry},
					{ID: 60010062, GroupID: 6001, Reel: 3, Row: 20, Symbol: pb.Symbol_Bar},
					{ID: 60010063, GroupID: 6001, Reel: 3, Row: 21, Symbol: pb.Symbol_Replay},
				}

				g := NewReelStripGroup(strips)

				_, _, err := set.Generete(DefaultProto.PayLine, g)
				if err != ErrResultSetNotGeneratable {
					t.Errorf("err want %v, but %v", ErrResultSetNotGeneratable, err)
				}
			})
		})
	})
	t.Run("BlueSevenLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BlueSevenLink)
		expType := reflect.TypeOf(new(BlueSevenLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BlueSevenLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_BlueSevenLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_BlueSevenLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_BlueSeven {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_BlueSeven, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("RedSevenLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_RedSevenLink)
		expType := reflect.TypeOf(new(RedSevenLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_RedSevenLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_RedSevenLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_RedSevenLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_RedSeven {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_RedSeven, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("BarLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BarLink)
		expType := reflect.TypeOf(new(BarLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BarLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_BarLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_BarLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_Bar {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_Bar, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("BellLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BellLink)
		expType := reflect.TypeOf(new(BellLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BellLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_BellLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_BellLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_Bell {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_Bell, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("WaterMelonLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_WaterMelonLink)
		expType := reflect.TypeOf(new(WaterMelonLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_WaterMelonLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_WaterMelonLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_WaterMelonLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_WaterMelon {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_WaterMelon, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("BerryLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BerryLink)
		expType := reflect.TypeOf(new(BerryLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BerryLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_BerryLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_BerryLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_Berry {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_Berry, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("ReplayLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_ReplayLink)
		expType := reflect.TypeOf(new(ReplayLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_ReplayLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_ReplayLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_ReplayLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_Replay {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_Replay, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("DoubleSevenLink", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_DoubleSevenLink)
		expType := reflect.TypeOf(new(DoubleSevenLinkResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_DoubleSevenLink, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_DoubleSevenLink {
			t.Errorf("Type() want %v, but %v", pb.ResultType_DoubleSevenLink, set.Type())
		}
		if set.Symbol() != pb.Symbol_BlueSeven {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_BlueSeven, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("BerryLeftAny", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BerryLeftAny)
		expType := reflect.TypeOf(new(BerryLeftAnyResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BerryLeftAny, expType, reflect.TypeOf(set))
		}
		if set.Type() != pb.ResultType_BerryLeftAny {
			t.Errorf("Type() want %v, but %v", pb.ResultType_BerryLeftAny, set.Type())
		}
		if set.Symbol() != pb.Symbol_Berry {
			t.Errorf("Symbol() want %v, but %v", pb.Symbol_Berry, set.Symbol())
		}
		_, m, _ := set.Generete(DefaultProto.PayLine, g)
		match := set.Match(DefaultProto.PayLine, m.SymbolGroup(g))
		if !match {
			t.Errorf("Match() want %v, but %v", true, match)
		}
	})
	t.Run("ReleaseRT", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_ReleaseRT)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_ReleaseRT, expType, reflect.TypeOf(set))
		}
	})
	t.Run("ReleaseRTonBB", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_ReleaseRTonBB)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_ReleaseRTonBB, expType, reflect.TypeOf(set))
		}
	})
	t.Run("ReleaseRTonRB", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_ReleaseRTonRB)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_ReleaseRTonRB, expType, reflect.TypeOf(set))
		}
	})
	t.Run("BonusSure", func(t *testing.T) {
		set := GenereteResultSet(pb.ResultType_BonusSure)
		expType := reflect.TypeOf(new(NoneResultSet))
		if reflect.TypeOf(set) != expType {
			t.Errorf("GenereteResultSet() by %v, want %v, but %v", pb.ResultType_BonusSure, expType, reflect.TypeOf(set))
		}
	})

	// switch rt {
	// case pb.ResultType_BlueSevenMiss: // 藍7 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_RedSevenMiss: // 紅7 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_BarMiss: // Bar 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_BellMiss: // 銅鐘 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_WaterMelonMiss: // 西瓜 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_BerryMiss: // 櫻桃 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_ReplayMiss: // Replay 前兩格連線，第三格出現但沒連線
	// 	return nil
	// case pb.ResultType_DoubleSevenMiss: // 雙7(藍7紅7混合) 前兩格連線，第三格出現但沒連線
	// 	return nil
	// }
}

func TestResultSetValidator(t *testing.T) {

	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	g := NewReelStripGroup(DefaultProto.ReelStrips)

	var types []pb.ResultType = UniqResultTypes(DefaultProto.PayTable)
	fmt.Println(types)

	v := NewResultSetValidator(types)

	r := new(NoneResultSet)
	_, m, _ := r.Generete(DefaultProto.PayLine, g)
	sg := m.SymbolGroup(g)

	var cnt int
	for v.MatchOthers(r.Type(), sg, DefaultProto.PayLine) {
		cnt++
		fmt.Println("regen", sg)
		_, m, _ := r.Generete(DefaultProto.PayLine, g)
		sg = m.SymbolGroup(g)
	}
	fmt.Println(cnt)
	fmt.Println(sg)
}
