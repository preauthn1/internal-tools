package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestGameResultManager_String(t *testing.T) {

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			nil,
			nil,
			{Symbol: proto.String("p3")},
			nil,
			nil,
		},
		[]*protodata.SlotSymbolDistribution{
			nil,
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
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
			{Symbol: proto.String("p16")},
			{Symbol: proto.String("p17")},
			{Symbol: proto.String("p18")},
			{Symbol: proto.String("p19")},
			{Symbol: proto.String("p20")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p21")},
			{Symbol: proto.String("p22")},
			{Symbol: proto.String("p23")},
			{Symbol: proto.String("p24")},
			{Symbol: proto.String("p25")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p26")},
			{Symbol: proto.String("p27")},
			{Symbol: proto.String("p28")},
			{Symbol: proto.String("p29")},
			{Symbol: proto.String("p30")},
		},
	}

	manager, err := NewGameResultManager(5, 6)
	if err != nil {
		t.Error(err)
		return
	}

	manager.GameResult = result

	t.Log(manager)
}

func TestGameResultManager_Generate_GameResult(t *testing.T) {

	conf.CLogLevel = mlog.DebugLevel
	conf.ClearLog = true
	conf.TestMode = false

	var SymbolTable []*protodata.SlotSymbolDistribution
	{
		SymbolTable = []*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(1), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(2), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(3), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(4), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(5), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(6), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(7), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(8), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(9), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(10), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(11), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(12), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(13), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(14), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(15), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(16), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(17), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(18), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},

			{Id: proto.Uint32(19), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(20), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(21), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(22), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(23), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(24), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(25), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(26), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(27), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(28), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(29), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(30), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(31), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(32), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(33), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(34), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(35), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(36), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(31), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(32), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(33), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(34), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(35), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(36), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(31), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(32), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(33), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(34), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(35), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(36), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},

			{Id: proto.Uint32(37), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(38), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(39), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(40), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(41), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(42), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(43), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(44), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(45), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(46), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(47), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(48), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(49), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(50), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(51), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(52), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(53), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(54), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(43), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(44), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(45), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(46), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(47), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(48), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(43), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(44), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(45), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(46), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(47), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(48), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(43), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{Id: proto.Uint32(44), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{Id: proto.Uint32(45), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{Id: proto.Uint32(46), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{Id: proto.Uint32(47), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{Id: proto.Uint32(48), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},

			{Id: proto.Uint32(55), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(56), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(57), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(58), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(59), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(60), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(61), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(62), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(63), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(64), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(65), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(66), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(67), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(68), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(69), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(70), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(71), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(72), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(61), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(62), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(63), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(64), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(65), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(66), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{Id: proto.Uint32(61), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(62), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(63), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(64), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(65), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{Id: proto.Uint32(66), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},

			{Id: proto.Uint32(73), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(74), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(75), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(76), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(77), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(78), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{Id: proto.Uint32(79), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(80), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(81), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(82), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(83), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(84), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{Id: proto.Uint32(85), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(86), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(87), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(88), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(89), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{Id: proto.Uint32(90), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
		}
	}

	manager, err := NewGameResultManager(5, 6)
	if err != nil {
		t.Error(err)
		return
	}

	manager.InvalidPosition = []Position{
		{Reel: 1, Row: 1},
		{Reel: 1, Row: 2},
		{Reel: 1, Row: 3},
		{Reel: 2, Row: 1},
		{Reel: 4, Row: 1},
		{Reel: 5, Row: 1},
		{Reel: 5, Row: 2},
		{Reel: 5, Row: 3},
	}

	manager.GenerateGameResult(1, SymbolTable)

	t.Log(manager)
}

func TestGameResultManager_Generate_Background(t *testing.T) {
	dists := []*protodata.ReelBgDistribution{
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(1), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(2), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(3), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg1"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg2"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg3"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(4), BgEnum: proto.String("bg4"), Proportion: proto.Uint32(1)},
		{GroupId: proto.Uint32(801), ReelNum: proto.Uint32(5), BgEnum: proto.String("bg0"), Proportion: proto.Uint32(1)},
	}

	manager, err := NewGameResultManager(5, 3)
	if err != nil {
		t.Error(err)
		return
	}

	manager.GenerateBackground(801, dists)

	t.Log(manager.BackgroundResult)
}

func TestGameResultManager_GenerateGameResult_New(t *testing.T) {

	conf.CLogLevel = mlog.FatalLevel
	conf.TimerFastMode = 1000 //40
	conf.TestMode = true
	conf.ClearLog = true
	conf.AnalysisMode = false

	t.Run("3 Row", func(tt *testing.T) {
		manager, err := NewGameResultManager(5, 7)
		if err != nil {
			t.Error(err)
			return
		}

		SymbolTable := []*protodata.SlotSymbolDistribution{
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},

			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},

			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},

			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(4), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},

			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(1)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(2)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(3)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(4)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(5)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(6)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p1"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p2"), Proportion: proto.Uint32(100), Order: proto.Uint32(7)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p3"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p4"), Proportion: proto.Uint32(100), Order: proto.Uint32(8)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("p5"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
			{GroupId: proto.Uint32(1), ReelNum: proto.Uint32(5), Symbol: proto.String("w1"), Proportion: proto.Uint32(100), Order: proto.Uint32(9)},
		}

		manager.GenerateGameResult(1, SymbolTable)

		t.Log(manager)
	})

}

func TestGameResultManager_Generate_FixedReel(t *testing.T) {

	conf.CLogLevel = mlog.DebugLevel
	conf.ClearLog = true
	conf.TestMode = false

	var SymbolTable []*protodata.SlotSymbolDistribution
	{
		SymbolTable = []*protodata.SlotSymbolDistribution{
			{Id: proto.Uint32(1), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(1), Symbol: proto.String("a1"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(2), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(2), Symbol: proto.String("a2"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(3), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(3), Symbol: proto.String("a3"), Proportion: proto.Uint32(100000)},
			{Id: proto.Uint32(4), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(4), Symbol: proto.String("a4"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(5), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(5), Symbol: proto.String("a5"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(6), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(6), Symbol: proto.String("a6"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(7), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(7), Symbol: proto.String("a7"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(8), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(8), Symbol: proto.String("a8"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(9), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(9), Symbol: proto.String("a9"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(10), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(10), Symbol: proto.String("a10"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(11), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(11), Symbol: proto.String("a11"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(12), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(12), Symbol: proto.String("a12"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(13), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(13), Symbol: proto.String("a13"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(14), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(14), Symbol: proto.String("a14"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(15), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(15), Symbol: proto.String("a15"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(16), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(16), Symbol: proto.String("a16"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(17), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(17), Symbol: proto.String("a17"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(18), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(18), Symbol: proto.String("a18"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(19), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(19), Symbol: proto.String("a19"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(20), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(20), Symbol: proto.String("a20"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(21), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(1), Order: proto.Uint32(21), Symbol: proto.String("a21"), Proportion: proto.Uint32(100)},

			{Id: proto.Uint32(22), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(1), Symbol: proto.String("b1"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(23), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(2), Symbol: proto.String("b2"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(24), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(3), Symbol: proto.String("b3"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(25), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(4), Symbol: proto.String("b4"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(26), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(5), Symbol: proto.String("b5"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(27), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(6), Symbol: proto.String("b6"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(28), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(7), Symbol: proto.String("b7"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(29), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(8), Symbol: proto.String("b8"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(30), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(9), Symbol: proto.String("b9"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(31), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(10), Symbol: proto.String("b10"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(32), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(11), Symbol: proto.String("b11"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(33), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(12), Symbol: proto.String("b12"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(34), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(13), Symbol: proto.String("b13"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(35), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(14), Symbol: proto.String("b14"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(36), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(15), Symbol: proto.String("b15"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(37), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(16), Symbol: proto.String("b16"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(38), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(17), Symbol: proto.String("b17"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(39), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(18), Symbol: proto.String("b18"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(40), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(19), Symbol: proto.String("b19"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(41), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(20), Symbol: proto.String("b20"), Proportion: proto.Uint32(100000)},
			{Id: proto.Uint32(42), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(2), Order: proto.Uint32(21), Symbol: proto.String("b21"), Proportion: proto.Uint32(100)},

			{Id: proto.Uint32(43), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(1), Symbol: proto.String("c1"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(44), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(2), Symbol: proto.String("c2"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(45), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(3), Symbol: proto.String("c3"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(46), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(4), Symbol: proto.String("c4"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(47), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(5), Symbol: proto.String("c5"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(48), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(6), Symbol: proto.String("c6"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(49), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(7), Symbol: proto.String("c7"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(50), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(8), Symbol: proto.String("c8"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(51), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(9), Symbol: proto.String("c9"), Proportion: proto.Uint32(100000)},
			{Id: proto.Uint32(52), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(10), Symbol: proto.String("c10"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(53), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(11), Symbol: proto.String("c11"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(54), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(12), Symbol: proto.String("c12"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(55), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(13), Symbol: proto.String("c13"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(56), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(14), Symbol: proto.String("c14"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(57), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(15), Symbol: proto.String("c15"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(58), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(16), Symbol: proto.String("c16"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(59), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(17), Symbol: proto.String("c17"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(60), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(18), Symbol: proto.String("c18"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(61), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(19), Symbol: proto.String("c19"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(62), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(20), Symbol: proto.String("c20"), Proportion: proto.Uint32(100)},
			{Id: proto.Uint32(63), GroupId: proto.Uint32(1), ReelNum: proto.Uint32(3), Order: proto.Uint32(21), Symbol: proto.String("c21"), Proportion: proto.Uint32(100)},
		}
	}

	manager, err := NewGameResultManager(3, 3)
	if err != nil {
		t.Error(err)
		return
	}

	manager.GenerateGameResult(1, SymbolTable)

	t.Log(manager)
}
