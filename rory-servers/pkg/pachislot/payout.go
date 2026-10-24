package pachislot

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	pb "example.invalid/ServerTeam/Servers/pkg/pachislot/protobuf-config"
)

// Payout 獎勵
type Payout interface {
	Type() pb.PayType           // 取得獎勵類型
	Odds() int64                // 贏得倍數
	Reward(*Game) (int64, bool) // 給獎，回傳贏得倍數
}

// PayoutChip 給錢獎勵
type PayoutChip struct {
	odds int64
}

func (p *PayoutChip) String() string {
	return fmt.Sprintf("[獎勵]%d倍金幣", p.odds)
}

// NewPayoutChip 建立給錢獎勵
func NewPayoutChip(odds int64) *PayoutChip {
	p := new(PayoutChip)
	p.odds = odds
	return p
}

// Type 取得獎勵類型
func (p *PayoutChip) Type() pb.PayType {
	return pb.PayType_PTChip
}

// Odds 取得獎勵類型
func (p *PayoutChip) Odds() int64 {
	return p.odds
}

// Reward 給獎，回傳贏得倍數
func (p *PayoutChip) Reward(game *Game) (int64, bool) {
	return p.odds, true
}

// PayoutReplay 免費遊戲獎勵
type PayoutReplay struct{}

func (p *PayoutReplay) String() string {
	return fmt.Sprintf("[獎勵]免費一次")
}

// NewPayoutReplay 建立免費遊戲獎勵
func NewPayoutReplay() *PayoutReplay {
	p := new(PayoutReplay)
	return p
}

// Type 取得獎勵類型
func (p *PayoutReplay) Type() pb.PayType {
	return pb.PayType_PTReplay
}

// Odds 取得獎勵類型
func (p *PayoutReplay) Odds() int64 {
	return 0
}

// Reward 給獎，回傳贏得倍數
func (p *PayoutReplay) Reward(game *Game) (int64, bool) {
	game.Replay = true
	return 0, true
}

// PayoutBonus 大獎獎勵
type PayoutBonus struct {
	Module pb.Scenario // 大獎類型
}

func (p *PayoutBonus) String() string {
	return fmt.Sprintf("[獎勵]%v", p.Module)
}

// NewPayoutBonus 建立大獎獎勵
func NewPayoutBonus(m pb.Scenario) *PayoutBonus {
	p := new(PayoutBonus)
	p.Module = m
	return p
}

// Type 取得獎勵類型
func (p *PayoutBonus) Type() pb.PayType {
	return pb.PayType_PTBonus
}

// Odds 取得獎勵類型
func (p *PayoutBonus) Odds() int64 {
	return 0
}

// Reward 給獎，回傳贏得倍數
func (p *PayoutBonus) Reward(game *Game) (int64, bool) {
	m := game.ModuleConfigs[p.Module]

	if p.Module == pb.Scenario_JAC {
		if game.JACTimes >= DefaultJACLimit {
			example.invalid("JAC次數達上限")
			return 0, false
		}
		game.Scenario = pb.Scenario_JAC
		game.JACTimes++
		game.JACTimesLeft = int(m.Count)
		return 0, true
	}

	game.Scenario = game.Bonus
	game.BonusTimesLeft = int(m.Count)
	game.RoadMap.Win(game.Bonus)

	return 0, true
}

// Payout1G 1G連獎勵
type Payout1G struct {
	Times int // 次數
}

func (p *Payout1G) String() string {
	return fmt.Sprintf("[獎勵]1G連[%d]次", p.Times)
}

// NewPayout1G 建立1G連獎勵
func NewPayout1G(times int) *Payout1G {
	p := new(Payout1G)
	p.Times = times
	return p
}

// Type 取得獎勵類型
func (p *Payout1G) Type() pb.PayType {
	return pb.PayType_PT1G
}

// Odds 取得獎勵類型
func (p *Payout1G) Odds() int64 {
	return 0
}

// Reward 給獎，回傳贏得倍數
func (p *Payout1G) Reward(game *Game) (int64, bool) {
	game.BonusStore += p.Times
	return 0, true
}
