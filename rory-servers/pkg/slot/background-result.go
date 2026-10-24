package slot

import (
	"sort"
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/prob"
	"example.invalid/ServerTeam/Servers/pkg/slot/protodata"
)

// Background 背景盤面
type Background string

// Backgrounds
const (
	BackgroundNone Background = "none"
	BackgroundGold Background = "gold"
)

// BackgroundResult 整體背景盤面
type BackgroundResult struct {
	Reel int32
	Row  int32

	Distributions []*protodata.ReelBgDistribution

	Result map[Position]Background `json:"-"`
}

func (result *BackgroundResult) String() string {

	sort.Slice(result.Distributions, func(i, j int) bool {
		return result.Distributions[i].GetReelNum() < result.Distributions[j].GetReelNum()
	})

	var bgs []string
	for _, row := range result.Distributions {
		bgs = append(bgs, row.GetBgEnum())
	}
	return strings.Join(bgs, ",")
}

// AddDistribution 加入 Distribution
func (result *BackgroundResult) AddDistribution(dist *protodata.ReelBgDistribution) {
	result.Distributions = append(result.Distributions, dist)

	if result.Result == nil {
		result.Result = make(map[Position]Background)
	}

	// 建立實際地圖, 產生背景, BgEnum 決定有幾格金
	var cnt int
	var ans []Background

	switch dist.GetBgEnum() {
	case "bg0":
		cnt = 0
	case "bg1":
		cnt = 1
	case "bg2":
		cnt = 2
	case "bg3":
		cnt = 3
	case "bg4":
		cnt = 4
	}

	for i := 0; i < int(result.Row); i++ {
		if i < cnt {
			ans = append(ans, BackgroundGold)
		} else {
			ans = append(ans, BackgroundNone)
		}
	}

	type order int
	const (
		asc order = iota
		desc
	)
	r := new(prob.Random)
	r.Add(asc, 1)
	r.Add(desc, 1)

	t := r.Choose().(order)

	switch t {
	case desc:
		sort.Slice(ans, func(i, j int) bool {
			return ans[i] == BackgroundNone
		})
	}

	for i, row := range ans {
		result.Result[Position{Reel: int32(dist.GetReelNum()), Row: int32(i + 1)}] = row
	}

}
