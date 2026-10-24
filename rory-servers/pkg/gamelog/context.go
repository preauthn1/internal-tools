package gamelog

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

// RoundContext 單局紀錄
type RoundContext struct {
	sync.Mutex
	Users     map[int32]*UserData
	BeginTime time.Time
	EndTime   time.Time
	Steps     []*GameStep
}

// AddUser 新增玩家
func (ctx *RoundContext) AddUser(uid int32, tradeid int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if _, ok := ctx.Users[uid]; !ok {
		ud := new(UserData)
		ud.TradeID = tradeid
		ctx.Users[uid] = ud
	}
}

// RemoveUser 移除玩家
func (ctx *RoundContext) RemoveUser(uid int32) {
	ctx.Lock()
	defer ctx.Unlock()
	delete(ctx.Users, uid)
}

// AddGameStep 添加遊戲步驟
func (ctx *RoundContext) AddGameStep(step *GameStep) {
	ctx.Steps = append(ctx.Steps, step)
}

// ClearGameStep 清除牌局記錄步驟
func (ctx *RoundContext) ClearGameStep(step *GameStep) {
	ctx.Steps = make([]*GameStep, 0)
}

// GetGameLogDesc 獲取牌局記錄描述
func (ctx *RoundContext) GetGameLogDesc() string {
	var buf bytes.Buffer

	buf.WriteString("\r\n")
	for _, st := range ctx.Steps {
		buf.WriteString(fmt.Sprintf("%s %s\r\n", st.StepTime, st.StepDesc))
	}

	return buf.String()
}

// UserScore 更新 Score
func (ctx *RoundContext) UserScore(uid int32, score int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.Score = score
	}
}

// UserUpdateScore 更新 UpdateScore
func (ctx *RoundContext) UserUpdateScore(uid int32, score int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.UpdateScore += score
	}
}

// UserBetScore 更新 BetScore
func (ctx *RoundContext) UserBetScore(uid int32, score int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.BetScore += score
	}
}

// UserPayoutScore 更新 PayoutScore
func (ctx *RoundContext) UserPayoutScore(uid int32, score int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.PayoutScore += score
	}
}

// UserValidScore 更新 ValidScore
func (ctx *RoundContext) UserValidScore(uid int32, score int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.ValidScore += score
	}
}

// UserTaxAmount 更新 TaxAmount
func (ctx *RoundContext) UserTaxAmount(uid int32, tax int64) {
	ctx.Lock()
	defer ctx.Unlock()
	if u, ok := ctx.Users[uid]; ok {
		u.TaxAmount += tax
	}
}
