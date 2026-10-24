package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestPayLine_GeneratePosition(t *testing.T) {

	origin := &protodata.SlotPayLine{
		Id:      proto.Uint32(1),
		GroupId: proto.Uint32(1),
		LineId:  proto.Uint32(1),
		Content: proto.String("1_2,2_2,3_2,4_2,5_2"),
	}

	line := NewPayLine(origin)

	line.GeneratePosition()

	t.Log(len(line.Positions))

	for _, position := range line.Positions {
		t.Log(position)
	}
}

func TestPayLine_Reverse(t *testing.T) {
	origin := &protodata.SlotPayLine{
		Id:      proto.Uint32(1),
		GroupId: proto.Uint32(1),
		LineId:  proto.Uint32(1),
		Content: proto.String("1_1,2_1,3_3,4_2,5_2"),
	}

	line := NewPayLine(origin)

	line.GeneratePosition()

	line.Reverse()

	t.Log(len(line.Positions))

	for _, position := range line.Positions {
		t.Log(position)
	}

}
