package slot

import (
	"testing"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestJackpotGame_GenerateProcess(t *testing.T) {

	for i := 0; i < 1; i++ {
		game := NewJackpotGame([]*Jackpot{
			{Level: JackpotGrand, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
			{Level: JackpotMajor, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
			{Level: JackpotMinor, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
			{Level: JackpotMini, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
			{Level: JackpotGoldReel, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
			{Level: JackpotSingle, Origin: &protodata.JackpotTable{PerProb: proto.Int32(1)}},
		}, 1, 1, 1, 1)

		game.WinJackpot = game.RandomJackpot()

		process := game.GenerateProcess()

		cntMap := make(map[JackpotLevel]int)

		for _, row := range process {
			cntMap[NewJackpotLevel(row.Icon)]++
		}

		t.Log(game.WinJackpot, cntMap)

		for level, cnt := range cntMap {
			if cnt == 3 {
				if level != game.WinJackpot.Level {
					t.Errorf("Has 3 not win icon, win: %s, process: %v", game.WinJackpot.Level.String(), process)
				}
			}
		}
	}
}

func TestJackpotGame_GetProcessData(t *testing.T) {

	game := new(JackpotGame)

	jackpot1 := &Jackpot{Level: JackpotGrand}
	jackpot2 := &Jackpot{Level: JackpotMajor}
	jackpot3 := &Jackpot{Level: JackpotMinor}
	jackpot4 := &Jackpot{Level: JackpotMini}

	game.Jackpots = []*Jackpot{jackpot1, jackpot2, jackpot3, jackpot4}
	game.WinJackpot = jackpot1

	process := game.GenerateProcess()

	for _, row := range process {
		t.Log(row)
	}

	game.Process = process

	t.Log(game.GetProcessData())

}
