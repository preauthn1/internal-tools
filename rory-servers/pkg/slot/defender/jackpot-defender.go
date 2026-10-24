package defender

import (
	"fmt"
	"io"

	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// JPDefender Jackpot 防護器
var JPDefender *JackpotDefender

// JackpotDefender Jackpot 防護器
type JackpotDefender struct {
	contribution int64
	award        int64
	increment    int64
	target       []slot.JackpotLevel
}

// NewJackpotDefender 建立 Jackpot 防護器
func NewJackpotDefender(increment int64, target ...slot.JackpotLevel) *JackpotDefender {
	def := new(JackpotDefender)
	def.increment = increment
	def.target = target
	return def
}

// Pool 取得目前水池金額
func (def *JackpotDefender) Pool() int64 {
	return def.contribution/1e4 - def.award
}

// Increase 增加防護池
func (def *JackpotDefender) Increase(pay int64) int64 {
	inc := pay * def.increment
	def.contribution += inc
	return inc
}

// Award 被贏得金額
func (def *JackpotDefender) Award(award int64) {
	def.award += award
}

// Eligible 是否可以出獎
func (def *JackpotDefender) Eligible(jackpot *slot.Jackpot, poolFn func() int64) bool {
	var isTarget bool

	for _, t := range def.target {
		if jackpot.Level == t {
			isTarget = true
			break
		}
	}

	if !isTarget {
		return true
	}

	return def.Pool() > poolFn()
}

// Info 對 writer 寫入資訊
func (def *JackpotDefender) Info(w io.Writer) {
	fmt.Fprintln(w, "[Jackpot Defender]")
	fmt.Fprintln(w, " ==========  Setting  ==========")
	fmt.Fprintf(w, "   - Increment: %d\n", def.increment)
	fmt.Fprintf(w, "   - Target: %v\n", def.target)
	fmt.Fprintln(w)
	fmt.Fprintln(w, " ============  Pool  ============")
	fmt.Fprintf(w, "   - Contribution(1/10000): %d\n", def.contribution)
	fmt.Fprintf(w, "   - Award: %d\n", def.award)
	fmt.Fprintf(w, "   - Pool: %d\n", def.Pool())
}
