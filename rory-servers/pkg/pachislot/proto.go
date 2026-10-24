package pachislot

import (
	"encoding/json"

	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// Proto 機率表設定
type Proto struct {
	GameConfig            *pb.GameConfig                   // 遊戲設定
	PayLine               []*pb.PayLine                    // 賠線設定
	PayTable              []*pb.PayoutOption               // 賠付表設定
	Pattern               []*pb.PatternOption              // Pattern設定
	ResultDistribution    []*pb.ResultDistributionOption   // 牌型機率表設定
	ModuleConfigs         []*pb.Module                     // 遊戲模式設定
	ModeConfigs           []*pb.Mode                       // 模式設定
	NextModeDistribution  []*pb.NextModeDistributionOption // 下次模式機率表設定
	BRDistribution        []*pb.BRDistributionOption       // 大獎機率表設定
	RTDistribution        []*pb.RTDistributionOption       // RT機率表設定
	RTReleaseDistribution []*pb.RTReleaseOption            // RT解除機率表設定
	ReelStrips            []*pb.ReelStripSymbol            // 排列帶設定
}

// Clone 使用一樣的值複製一套出來
func (p *Proto) Clone() *Proto {
	dist := new(Proto)
	b, _ := json.Marshal(p)
	json.Unmarshal(b, dist)
	return dist
}
