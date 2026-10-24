package gen

import (
	"testing"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/fixtures"
	"example.invalid/ServerTeam/Servers/pkg/pachislot"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

func Test_mDataPayLine(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	lines := game.mDataPayLine()
	for _, line := range lines {
		t.Log(line)
	}
}

func Test_mDataReelStrips(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	strips := game.mDataReelStrips()
	for _, strip := range strips {
		t.Log(strip)
	}
}

func Test_mDataRoadMapBonus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	rm := new(pachislot.RoadMap)
	rm.TotalCnt = 9939
	rm.SpinCnt = 387
	rm.FeaturesCnt = map[pb.Scenario]int{
		pb.Scenario_BB: 4,
		pb.Scenario_RB: 8,
	}

	data := game.mDataRoadMapBonus(rm)
	for _, row := range data {
		t.Log(row)
	}
}

func Test_mDataRoadMapSnapshots(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	rm := new(pachislot.RoadMap)
	rm.TotalCnt = 9939
	rm.SpinCnt = 387
	rm.FeaturesCnt = map[pb.Scenario]int{
		pb.Scenario_BB: 4,
		pb.Scenario_RB: 8,
	}
	rm.Snapshots = []*pachislot.RoadMapSnapshot{
		{Feature: pb.Scenario_BB, Count: 332, Time: time.Now()},
		{Feature: pb.Scenario_BB, Count: 222, Time: time.Now()},
		{Feature: pb.Scenario_RB, Count: 131, Time: time.Now()},
		{Feature: pb.Scenario_BB, Count: 455, Time: time.Now()},
	}

	data := game.mDataRoadMapSnapshots(rm)
	for _, row := range data {
		t.Log(row)
	}
}

func Test_mGameSetting(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	m := game.mGameSetting()
	t.Log(m, m.MsgData)
}

func Test_mRoadMap(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	m := game.mRoadMap()
	t.Log(m, m.MsgData)
}

func Test_mGameStatus(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	grs := fixtures.NewTestGameRoomServerWithConfig(cfg)
	game := grs.GameLogicManager.GetGameLogic(0).(*Game)

	m := game.mGameStatus()
	t.Log(m, m.MsgData)
}
