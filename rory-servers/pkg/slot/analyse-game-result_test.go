package slot

import (
	"fmt"
	"runtime"
	"testing"

	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

func TestAnalyse_GameResult(t *testing.T) {
	data, _ := protodata.NewProtoData("./protodata/mock", protodata.OptSlotSymbolDistribution)
	data.Reload()

	grm, _ := NewGameResultManager(5, 3)
	grm.GenerateGameResult(1, data.SlotSymbolDistribution)

	numCPUs := runtime.NumCPU()

	var WorkerCount = 4 * numCPUs
	var TotalCount = 100
	var SymbolGroupIDs = []int32{301, 302, 303, 304, 305, 306, 307, 308}

	for _, SymbolGroupID := range SymbolGroupIDs {
		t.Run(fmt.Sprintf("SymbolGroupID: %d", SymbolGroupID), func(tt *testing.T) {
			var JobChan = make(chan int32, TotalCount)
			var ResultChan = make(chan *GameResultManager, TotalCount)

			var Worker = func(id int, jobs chan int32, results chan *GameResultManager) {
				for job := range jobs {
					grm, _ := NewGameResultManager(5, 3)
					grm = grm.GenerateGameResult(job, data.SlotSymbolDistribution)

					results <- grm
				}
			}

			for i := 0; i < WorkerCount; i++ {
				go Worker(i, JobChan, ResultChan)
			}

			for i := 0; i < TotalCount; i++ {
				JobChan <- SymbolGroupID
			}
			close(JobChan)

			var CreateSymbolMap = func() map[string]int {
				m := make(map[string]int)
				m["p1"] = 0
				m["p2"] = 0
				m["p3"] = 0
				m["p4"] = 0
				m["p5"] = 0
				m["p6"] = 0
				m["p7"] = 0
				m["w1"] = 0
				return m
			}

			ResultMap := make(map[int]map[int]map[string]int)

			for i := 1; i <= 5; i++ {
				ResultMap[i] = make(map[int]map[string]int)
				for j := 1; j <= 3; j++ {
					ResultMap[i][j] = CreateSymbolMap()
				}
			}

			for i := 0; i < TotalCount; i++ {
				grm := <-ResultChan

				grm.ReelMode()

				for reelNum, reel := range grm.GameResult {
					for rowNum, row := range reel {
						ResultMap[reelNum+1][rowNum+1][row.GetSymbol()]++
					}
				}
			}

			t.Log("=================================================================")
			t.Logf("TotalCount: %d, SymbolGroupID: %d\n", TotalCount, SymbolGroupID)
			t.Log("=================================================================")
			for i := 1; i <= 5; i++ {
				for j := 1; j <= 3; j++ {
					t.Logf("Reel: %d, Row: %d\n", i, j)
					// t.Logf("p1: %d, p2: %d, p3: %d, p4: %d, p5: %d, p6: %d, p7: %d, w1: %d\n", ResultMap[i][j]["p1"], ResultMap[i][j]["p2"], ResultMap[i][j]["p3"], ResultMap[i][j]["p4"], ResultMap[i][j]["p5"], ResultMap[i][j]["p6"], ResultMap[i][j]["p7"], ResultMap[i][j]["w1"])
					t.Logf("p1: %.2f%%, p2: %.2f%%, p3: %.2f%%, p4: %.2f%%, p5: %.2f%%, p6: %.2f%%, p7: %.2f%%, w1: %.2f%%\n",
						float64(ResultMap[i][j]["p1"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p2"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p3"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p4"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p5"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p6"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["p7"])/float64(TotalCount)*100,
						float64(ResultMap[i][j]["w1"])/float64(TotalCount)*100,
					)
				}
			}
			t.Log("=================================================================")
		})
	}
}
