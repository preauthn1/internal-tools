package wlzb

import (
	"math"
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestLogic_RemoveFirstSelectChance(t *testing.T) {

	arr := []*SelectChance(nil)

	paytb := &protodata.SlotPayTable{Payout: proto.String("501")}

	opts := []*protodata.FreeSpinTable{
		{Id: proto.Uint32(1), GroupId: proto.Int32(501), FreeSpinCount: proto.Uint32(5)},
		{Id: proto.Uint32(2), GroupId: proto.Int32(501), FreeSpinCount: proto.Uint32(10)},
		{Id: proto.Uint32(3), GroupId: proto.Int32(501), FreeSpinCount: proto.Uint32(15)},
		{Id: proto.Uint32(4), GroupId: proto.Int32(502), FreeSpinCount: proto.Uint32(15)},
	}

	var betMoney int64 = 10
	var betRate int64 = 1

	for i := 0; i < 1; i++ {
		sc, err := NewSelectChance(paytb, opts, betMoney, betRate)
		if err != nil {
			t.Error(err)
		}

		t.Log(sc)

		arr = append(arr, sc)
	}

	sc := arr[0]
	t.Log(sc)
	t.Log(sc.ChooseByID(2))

	t.Log(arr)

	arr = arr[1:]

	t.Log(sc)
	t.Log(arr)

	t.Log(math.Pow(float64(3), float64(5)))
}
