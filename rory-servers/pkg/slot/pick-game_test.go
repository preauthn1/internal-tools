package slot

import (
	"fmt"
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestNewPickGame(t *testing.T) {

	bonusTbs := []*protodata.BonusTable{
		{Id: proto.Int32(1), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("10"), PerProb: proto.Int32(100000)},
		{Id: proto.Int32(2), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("25"), PerProb: proto.Int32(70000)},
		{Id: proto.Int32(3), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("30"), PerProb: proto.Int32(50000)},
		{Id: proto.Int32(4), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("45"), PerProb: proto.Int32(30000)},
		{Id: proto.Int32(5), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("50"), PerProb: proto.Int32(10000)},
		{Id: proto.Int32(6), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("65"), PerProb: proto.Int32(7000)},
		{Id: proto.Int32(7), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("70"), PerProb: proto.Int32(5000)},
		{Id: proto.Int32(8), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("85"), PerProb: proto.Int32(3000)},
		{Id: proto.Int32(9), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("100"), PerProb: proto.Int32(1000)},
	}

	game := NewPickGame(bonusTbs, 9, 1, 10, 10)

	for key, row := range game.Process {
		fmt.Println(key, row)
	}

	for _, row := range game.All {
		fmt.Println(row)
	}
}

func TestPickGame_Choose(t *testing.T) {
	game := new(PickGame)
	game.GameCount = 3
	game.GameBetMoney = 10
	game.WinMoney = 30
	game.Process = make(map[int32]*PickGameDetail)

	origin := &protodata.BonusTable{Id: proto.Int32(1), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("10"), PerProb: proto.Int32(100000)}

	for i := 0; i < 3; i++ {
		detail := new(PickGameDetail)
		detail.Origin = origin
		detail.ID = int32(i + 1)
		detail.Money = 100
		game.Process[detail.ID] = detail
	}

	game.Choose(1, 2)
	if game.Process[1].Location != 2 {
		t.Errorf("ID: %d, Location want %d, but %d", 1, 2, game.Process[1].Location)
	}
	game.Choose(2, 6)
	if game.Process[2].Location != 6 {
		t.Errorf("ID: %d, Location want %d, but %d", 2, 6, game.Process[2].Location)
	}
	game.Choose(3, 1)
	if game.Process[3].Location != 1 {
		t.Errorf("ID: %d, Location want %d, but %d", 3, 1, game.Process[3].Location)
	}
	game.Choose(3, 8)
	if game.Process[3].Location != 1 {
		t.Errorf("ID: %d, Location want %d, but %d", 3, 1, game.Process[3].Location)
	}
}

func TestPickGame_IsFinished(t *testing.T) {
	game := new(PickGame)
	game.GameCount = 3
	game.GameBetMoney = 10
	game.WinMoney = 30
	game.Process = make(map[int32]*PickGameDetail)

	origin := &protodata.BonusTable{Id: proto.Int32(1), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("10"), PerProb: proto.Int32(100000)}

	for i := 0; i < 3; i++ {
		detail := new(PickGameDetail)
		detail.Origin = origin
		detail.ID = int32(i + 1)
		detail.Money = 100
		game.Process[detail.ID] = detail
	}

	if game.IsFinished() {
		t.Errorf("Play 0 times, should not Finished")
	}
	game.Choose(1, 2)
	if game.IsFinished() {
		t.Errorf("Play 1 times, should not Finished")
	}
	game.Choose(2, 6)
	if game.IsFinished() {
		t.Errorf("Play 2 times, should not Finished")
	}
	game.Choose(3, 1)
	if !game.IsFinished() {
		t.Errorf("Play 3 times, should Finished")
	}
}

func TestPickGame_GetProcessData(t *testing.T) {
	bonusTbs := []*protodata.BonusTable{
		{Id: proto.Int32(1), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("10"), PerProb: proto.Int32(100000)},
		{Id: proto.Int32(2), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("25"), PerProb: proto.Int32(70000)},
		{Id: proto.Int32(3), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("30"), PerProb: proto.Int32(50000)},
		{Id: proto.Int32(4), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("45"), PerProb: proto.Int32(30000)},
		{Id: proto.Int32(5), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("50"), PerProb: proto.Int32(10000)},
		{Id: proto.Int32(6), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("65"), PerProb: proto.Int32(7000)},
		{Id: proto.Int32(7), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("70"), PerProb: proto.Int32(5000)},
		{Id: proto.Int32(8), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("85"), PerProb: proto.Int32(3000)},
		{Id: proto.Int32(9), GroupId: proto.Int32(1), RewardType: proto.String("chip"), RewardValue: proto.String("100"), PerProb: proto.Int32(1000)},
	}

	game := NewPickGame(bonusTbs, 12, 5, 10, 10)

	fmt.Println(game.GetProcessData())
}

func TestPickGame_Exist(t *testing.T) {

	game := new(PickGame)
	game.Process = make(map[int32]*PickGameDetail)

	fmt.Println(game.Exist(0))
	fmt.Println(game.Exist(1))
	fmt.Println(game.Exist(2))
}
