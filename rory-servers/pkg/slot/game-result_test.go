package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestGameResult_String(t *testing.T) {

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p1")},
			{Symbol: proto.String("p2")},
			{Symbol: proto.String("p3")},
			{Symbol: proto.String("p4")},
			{Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p6")},
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
			{Symbol: proto.String("p10")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p11")},
			{Symbol: proto.String("p12")},
			{Symbol: proto.String("p13")},
			{Symbol: proto.String("p14")},
			{Symbol: proto.String("p15")},
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

	t.Log(result)
}

func TestGameResult_Transfor(t *testing.T) {

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p1")},
			{Symbol: proto.String("p2")},
			{Symbol: proto.String("p3")},
			{Symbol: proto.String("p4")},
			{Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p6")},
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
			{Symbol: proto.String("p10")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p11")},
			{Symbol: proto.String("p12")},
			{Symbol: proto.String("p13")},
			{Symbol: proto.String("p14")},
			{Symbol: proto.String("p15")},
		},
	}

	t.Log(result)

	result = result.Transfor()
	t.Log(result)

	result = result.Transfor()
	t.Log(result)
}

func TestGameResult_Reverse(t *testing.T) {

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p1")},
			{Symbol: proto.String("p2")},
			{Symbol: proto.String("p3")},
			{Symbol: proto.String("p4")},
			{Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p6")},
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
			{Symbol: proto.String("p10")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p11")},
			{Symbol: proto.String("p12")},
			{Symbol: proto.String("p13")},
			{Symbol: proto.String("p14")},
			{Symbol: proto.String("p15")},
		},
	}

	t.Log(result)

	result = result.Transfor()
	t.Log(result)

	result = result.Reverse()
	t.Log(result)

	result = result.Transfor()
	t.Log(result)
}

func TestGameResult_GetSymbolIDByRow(t *testing.T) {

	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p1")},
			{Symbol: proto.String("p2")},
			{Symbol: proto.String("p3")},
			{Symbol: proto.String("p4")},
			{Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p6")},
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
			{Symbol: proto.String("p10")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p11")},
			{Symbol: proto.String("p12")},
			{Symbol: proto.String("p13")},
			{Symbol: proto.String("p14")},
			{Symbol: proto.String("p15")},
		},
	}

	t.Log(result.GetSymbolIDByRow(0))
	t.Log(result.GetSymbolIDByRow(1))
	t.Log(result.GetSymbolIDByRow(2))
	t.Log(result.GetSymbolIDByRow(3))
	t.Log(result.GetSymbolIDByRow(4))

}

func TestGameResult_GetReelSymbol(t *testing.T) {
	result := GameResult{
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p1")},
			{Symbol: proto.String("p2")},
			{Symbol: proto.String("p3")},
			{Symbol: proto.String("p4")},
			{Symbol: proto.String("p5")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p6")},
			{Symbol: proto.String("p7")},
			{Symbol: proto.String("p8")},
			{Symbol: proto.String("p9")},
			{Symbol: proto.String("p10")},
		},
		[]*protodata.SlotSymbolDistribution{
			{Symbol: proto.String("p11")},
			{Symbol: proto.String("p12")},
			{Symbol: proto.String("p13")},
			{Symbol: proto.String("p14")},
			{Symbol: proto.String("p15")},
		},
	}

	t.Log(result)

	data := result.GetReelSymbol()

	for reelID, reel := range data {
		for rowID, row := range reel {
			t.Logf("Reel: %d, Row: %d, Symbol: %s\n", reelID, rowID, row[0].GetSymbol())
		}
	}
}
