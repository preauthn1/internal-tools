package gamelog

import (
	"time"

	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// EmptyGameLog 空的遊戲牌局記錄實現
type EmptyGameLog struct {
	ctx             *RoundContext
	nextJiangRecord *netproto.NextJiangRecord
}

// NewEmptyGameLog 創建空的牌局記錄實現
func NewEmptyGameLog() *EmptyGameLog {
	return new(EmptyGameLog)
}

// Context 當下牌局context
func (l *EmptyGameLog) Context() *RoundContext {
	return l.ctx
}

// BeginGameLog 開始遊戲記錄
func (l *EmptyGameLog) BeginGameLog() *RoundContext {
	ctx := new(RoundContext)
	l.ctx = ctx
	l.ctx.Users = make(map[int32]*UserData)
	l.ctx.BeginTime = time.Now()
	return ctx
}

// AddGameStep 添加遊戲步驟
func (l *EmptyGameLog) AddGameStep(step *GameStep) {
	if l.ctx == nil {
		ctx := new(RoundContext)
		l.ctx = ctx
		l.ctx.BeginTime = time.Now()
	}
	l.ctx.AddGameStep(step)
}

// EndGameLog 遊戲步驟結束
func (l *EmptyGameLog) EndGameLog() {
	l.ctx.EndTime = time.Now()
}

// SetNextJiangRecord ?
func (l *EmptyGameLog) SetNextJiangRecord(record *netproto.NextJiangRecord) {
	l.nextJiangRecord = record
}
