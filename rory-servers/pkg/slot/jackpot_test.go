package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJackpot_GetCurrentPool_NoLimit(t *testing.T) {

	jpTable := &protodata.JackpotTable{
		Id:           proto.Int32(1),
		GroupId:      proto.Int32(1),
		IsProgress:   proto.Int32(1),
		IsLinked:     proto.Int32(0),
		Increment:    proto.Int32(5),
		StartUp:      proto.Int64(20000000),
		IsMultiLevel: proto.Int32(1),
		Level:        proto.String("Grand"),
		IsCeling:     proto.Int32(0),
		PerProb:      proto.Int32(1),
	}

	jp := &Jackpot{
		Origin:            jpTable,
		TotalContribution: 50000,
	}

	rate := 88
	money := 1000

	if award := jp.GetCurrentPool(); award != 20000005 {
		t.Errorf("Award want %d, but %d", 20000005, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000049 {
		t.Errorf("Award want %d, but %d", 20000049, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000093 {
		t.Errorf("Award want %d, but %d", 20000093, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000137 {
		t.Errorf("Award want %d, but %d", 20000137, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000181 {
		t.Errorf("Award want %d, but %d", 20000181, award)
	}
}

func TestJackpot_GetCurrentPool_Limit(t *testing.T) {

	jpTable := &protodata.JackpotTable{
		Id:           proto.Int32(1),
		GroupId:      proto.Int32(1),
		IsProgress:   proto.Int32(1),
		IsLinked:     proto.Int32(0),
		Increment:    proto.Int32(5),
		StartUp:      proto.Int64(20000000),
		IsMultiLevel: proto.Int32(1),
		Level:        proto.String("Grand"),
		IsCeling:     proto.Int32(1),
		CelingValue:  proto.Int64(20000150),
		PerProb:      proto.Int32(1),
	}

	jp := &Jackpot{
		Origin:            jpTable,
		TotalContribution: 50000,
	}

	rate := 88
	money := 1000

	if award := jp.GetCurrentPool(); award != 20000005 {
		t.Errorf("Award want %d, but %d", 20000005, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000049 {
		t.Errorf("Award want %d, but %d", 20000049, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000093 {
		t.Errorf("Award want %d, but %d", 20000093, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000137 {
		t.Errorf("Award want %d, but %d", 20000137, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000150 {
		t.Errorf("Award want %d, but %d", 20000150, award)
	}

	jp.Increase(int64(money * rate))

	if award := jp.GetCurrentPool(); award != 20000150 {
		t.Errorf("Award want %d, but %d", 20000150, award)
	}
}

func TestJackpot_Increase_Progress_NoLimit(t *testing.T) {

	jpTable := &protodata.JackpotTable{
		Id:           proto.Int32(1),
		GroupId:      proto.Int32(1),
		IsProgress:   proto.Int32(1),
		IsLinked:     proto.Int32(0),
		Increment:    proto.Int32(5),
		StartUp:      proto.Int64(20000000),
		IsMultiLevel: proto.Int32(1),
		Level:        proto.String("Grand"),
		IsCeling:     proto.Int32(0),
		CelingValue:  proto.Int64(500000),
		PerProb:      proto.Int32(1),
	}

	jp := &Jackpot{
		Origin:            jpTable,
		TotalContribution: 50000,
	}

	rate := 88
	money := 1000

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 490000 {
		t.Errorf("TotalContribution want %d, but %d", 490000, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 930000 {
		t.Errorf("TotalContribution want %d, but %d", 930000, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 1370000 {
		t.Errorf("TotalContribution want %d, but %d", 1370000, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 1810000 {
		t.Errorf("TotalContribution want %d, but %d", 1810000, jp.TotalContribution)
	}
}

func TestJackpot_Increase_NoProgress(t *testing.T) {

	jpTable := &protodata.JackpotTable{
		Id:           proto.Int32(1),
		GroupId:      proto.Int32(1),
		IsProgress:   proto.Int32(0),
		IsLinked:     proto.Int32(0),
		Increment:    proto.Int32(5),
		StartUp:      proto.Int64(20000000),
		IsMultiLevel: proto.Int32(1),
		Level:        proto.String("Grand"),
		IsCeling:     proto.Int32(1),
		CelingValue:  proto.Int64(500000),
		PerProb:      proto.Int32(1),
	}

	jp := &Jackpot{
		Origin: jpTable,
	}

	rate := 88
	money := 1000

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 0 {
		t.Errorf("TotalContribution want %d, but %d", 0, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 0 {
		t.Errorf("TotalContribution want %d, but %d", 0, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 0 {
		t.Errorf("TotalContribution want %d, but %d", 0, jp.TotalContribution)
	}

	jp.Increase(int64(money * rate))

	if jp.TotalContribution != 0 {
		t.Errorf("TotalContribution want %d, but %d", 0, jp.TotalContribution)
	}
}

func TestJackpot_GetCurrentPoolByBet(t *testing.T) {

	t.Run("Link Pool", func(tt *testing.T) {
		jpTable := &protodata.JackpotTable{
			Id:           proto.Int32(1),
			GroupId:      proto.Int32(1),
			IsProgress:   proto.Int32(1),
			IsLinked:     proto.Int32(1),
			Increment:    proto.Int32(5),
			StartUp:      proto.Int64(10000),
			IsMultiLevel: proto.Int32(1),
			Level:        proto.String("Grand"),
			IsCeling:     proto.Int32(0),
			PerProb:      proto.Int32(1),
		}

		jp := &Jackpot{
			Origin:            jpTable,
			TotalContribution: 50000,
		}

		ans := jp.GetCurrentPoolByBet(10, 100)
		if ans != 10005 {
			tt.Errorf("Want %d, but %d", 10005, ans)
		}
	})

	t.Run("NoLink Pool", func(tt *testing.T) {
		jpTable := &protodata.JackpotTable{
			Id:           proto.Int32(1),
			GroupId:      proto.Int32(1),
			IsProgress:   proto.Int32(1),
			IsLinked:     proto.Int32(0),
			Increment:    proto.Int32(5),
			StartUp:      proto.Int64(10000),
			IsMultiLevel: proto.Int32(1),
			Level:        proto.String("Grand"),
			IsCeling:     proto.Int32(0),
			PerProb:      proto.Int32(1),
		}

		jp := &Jackpot{
			Origin:            jpTable,
			TotalContribution: 50000,
		}

		ans := jp.GetCurrentPoolByBet(10, 100)
		if ans != 1000 {
			tt.Errorf("Want %d, but %d", 1000, ans)
		}
	})

	t.Run("NoLink Pool Max Bet", func(tt *testing.T) {
		jpTable := &protodata.JackpotTable{
			Id:           proto.Int32(1),
			GroupId:      proto.Int32(1),
			IsProgress:   proto.Int32(1),
			IsLinked:     proto.Int32(0),
			Increment:    proto.Int32(5),
			StartUp:      proto.Int64(10000),
			IsMultiLevel: proto.Int32(1),
			Level:        proto.String("Grand"),
			IsCeling:     proto.Int32(0),
			PerProb:      proto.Int32(1),
		}

		jp := &Jackpot{
			Origin:            jpTable,
			TotalContribution: 50000,
		}

		ans := jp.GetCurrentPoolByBet(100, 100)
		if ans != 10005 {
			tt.Errorf("Want %d, but %d", 10005, ans)
		}
	})
}
