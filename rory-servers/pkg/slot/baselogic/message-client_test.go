package baselogic

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestBuildReelResult(t *testing.T) {

	logic := new(SlotLogicBase)
	logic.PayLineMode = slot.PayLineModeAll

	var manager *slot.GameResultManager
	{
		result := slot.GameResult{
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
		_manager, err := slot.NewGameResultManager(5, 6)
		if err != nil {
			t.Error(err)
			return
		}
		_manager.GameResult = result

		manager = _manager
	}

	ans := logic.BuildReelResult(manager, nil)
	for _, row := range ans {
		t.Log(row)
	}
}
