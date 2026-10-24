package pachislot

import (
	"context"
	"fmt"

	humanize "example.invalid/dustin/go-humanize"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// Analyser 博青斯洛分析工具
type Analyser interface {
	Run(context.Context, *Game)
}

type analyser struct {
	deadline int
}

func (a *analyser) Run(ctx context.Context, game *Game) {

	var totalCost, totalWin int64
	var totalReleaseRT, total1G int
	resultMap := make(map[pb.Scenario]map[pb.ResultType]int)
	resultMap[pb.Scenario_NG] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_BB] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_RB] = make(map[pb.ResultType]int)
	resultMap[pb.Scenario_JAC] = make(map[pb.ResultType]int)

	for game.RoadMap.TotalCnt < a.deadline {
		cost, result, releaseRT := game.Spin()
		totalCost += cost
		resultMap[game.Scenario][result]++
		if releaseRT {
			totalReleaseRT++
		}
		payouts := game.Reward(result)
		for _, p := range payouts {
			totalWin += p.Odds()
			switch p.Type() {
			case pb.PayType_PT1G:
				total1G += p.(*Payout1G).Times
			}
		}
	}

	fmt.Printf("Total Spin: %s (Times)\n", humanize.Comma(int64(game.RoadMap.TotalCnt)))
	fmt.Printf("Total Cost: %s (Coin)\n", humanize.Comma(totalCost))
	fmt.Printf("Total Win: %s (Coin)\n", humanize.Comma(totalWin))
	fmt.Printf("Total ReleaseRT: %s (Times)\n", humanize.Comma(int64(totalReleaseRT)))
	fmt.Printf("Total 1G: %s (Times)\n", humanize.Comma(int64(total1G)))
	fmt.Printf("Final RoadMap: %v\n", game.RoadMap)
	fmt.Printf("RTP: %.3f%%\n", float64(totalWin)/float64(totalCost)*100)
	fmt.Println("=======================================================")
	fmt.Printf("Result in NG:\n")
	for t, cnt := range resultMap[pb.Scenario_NG] {
		fmt.Printf("  %v: %s\n", t, humanize.Comma(int64(cnt)))
	}
	fmt.Printf("Result in BB:\n")
	for t, cnt := range resultMap[pb.Scenario_BB] {
		fmt.Printf("  %v: %s\n", t, humanize.Comma(int64(cnt)))
	}
	fmt.Printf("Result in RB:\n")
	for t, cnt := range resultMap[pb.Scenario_RB] {
		fmt.Printf("  %v: %s\n", t, humanize.Comma(int64(cnt)))
	}
	fmt.Printf("Result in JAC:\n")
	for t, cnt := range resultMap[pb.Scenario_JAC] {
		fmt.Printf("  %v: %s\n", t, humanize.Comma(int64(cnt)))
	}
}
