package infc

import (
	"context"
	"sort"
	"sync"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
)

// RoundContext 一局內資訊
type RoundContext struct {
	Pool int64

	mu sync.Mutex

	ctx    context.Context
	cancel context.CancelFunc

	Rank      Rank
	GameState GameState
}

// NewRoundContext Create RoundContext
func NewRoundContext(users []*gamelogic.GameUser) *RoundContext {
	if len(users) == 0 {
		return nil
	}

	ctx := new(RoundContext)

	ctx.ctx, ctx.cancel = context.WithCancel(context.TODO())

	for _, u := range users {
		p := NewPlayer(u)
		p.State = StateOnline
		ctx.Rank = ctx.Rank.Append(p)
	}

	return ctx
}

// Player 取得玩家
func (ctx *RoundContext) Player(uid int32) *Player {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.Rank.Player(uid)
}

// RefreshRank 更新排行榜
func (ctx *RoundContext) RefreshRank() {
	sort.Sort(ctx.Rank)
}

// Winner 取得分數最高的玩家
func (ctx *RoundContext) Winner() []*Player {
	return ctx.Rank.Winner()
}

// GetScore 取得分數
func (ctx *RoundContext) GetScore(uid int32) int64 {
	if p := ctx.Player(uid); p != nil {
		return p.Score
	}
	return 0
}

// SetScore 設定分數
func (ctx *RoundContext) SetScore(uid int32, score int64) {
	if p := ctx.Player(uid); p != nil && p.State == StateGaming {
		p.Score = score
	}
}

// AddScore 異動分數
func (ctx *RoundContext) AddScore(uid int32, score int64) {
	if p := ctx.Player(uid); p != nil && p.State == StateGaming {
		p.Score += score
	}
}

// GetMiss 取得Miss次數
func (ctx *RoundContext) GetMiss(uid int32) int64 {
	if p := ctx.Player(uid); p != nil {
		return p.Miss
	}
	return 0
}

// SetMiss 設定Miss次數
func (ctx *RoundContext) SetMiss(uid int32, miss int64) {
	if p := ctx.Player(uid); p != nil && p.State == StateGaming {
		p.Miss = miss
	}
}

// AddMiss 異動Miss次數
func (ctx *RoundContext) AddMiss(uid int32, miss int64) {
	if p := ctx.Player(uid); p != nil && p.State == StateGaming {
		p.Miss += miss
	}
}
