package model

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"sync"
	"time"

	"aichatbot/internal/store"
)

// prices 是每百萬 token 的價格（美元）：輸入、輸出。本機模型不在表上，費用記 0
var prices = map[string][2]float64{
	"claude-haiku-5-5": {0.10, 0.50},
}

// Cost 算一次呼叫的費用（美元）；model 不在價格表上時回 false
func Cost(u Usage) (float64, bool) {
	if u.Model == AFM {
		return 0, true
	}
	p, ok := prices[u.Model]
	if !ok {
		return 0, false
	}
	return (float64(u.InputTokens)*p[0] + float64(u.OutputTokens)*p[1]) / 1e6, true
}

// Ledger 是用量紀錄（store.DB）
type Ledger interface {
	RecordUsage(ctx context.Context, u store.UsageRecord) error
	CostBetween(ctx context.Context, from, to time.Time) (float64, error)
}

// Switch 是對外用的模型：依模型開關選 Haiku 或 Apple 本機模型（AFM），每次呼叫記一筆用量。
// UseHaiku、Budget 是建立時的初始值，之後（後台存檔）改用 Apply，不要直接改欄位
type Switch struct {
	Haiku    ChatModel // 沒有 API key 時是 nil
	AFM      ChatModel
	UseHaiku bool
	Budget   float64 // Haiku 每月費用上限（美元），0 表示不限
	Ledger   Ledger
	Log      *slog.Logger
	Now      func() time.Time
	// OnBudgetExceeded：本月 Haiku 費用到上限、自動切回本機模型時呼叫（每個月一次；通知員工群組）。nil 不通知
	OnBudgetExceeded func(ctx context.Context, spent, budget float64)

	setMu      sync.RWMutex // 保護 UseHaiku、Budget
	mu         sync.Mutex
	warnedFrom time.Time // 已經記過「到上限」的月份（月初）
}

// Settings 回目前的模型開關（afm 或 haiku）和每月上限
func (s *Switch) Settings() (string, float64) {
	s.setMu.RLock()
	defer s.setMu.RUnlock()
	if s.UseHaiku {
		return ModelHaiku, s.Budget
	}
	return AFM, s.Budget
}

// HaikuAvailable：有設定 API key（才能切到 Haiku）
func (s *Switch) HaikuAvailable() bool { return s.Haiku != nil }

// ErrNoAPIKey：沒有設定 API key，不能切到 Haiku
var ErrNoAPIKey = [REDACTED]"沒有設定 ANTHROPIC_API_KEY，不能切到 Haiku")

// CheckSetting 檢查模型開關和每月費用上限：model 要是 afm 或 haiku（haiku 要有 API key），上限要是 0 以上的數字
func CheckSetting(model string, budget float64, haikuAvailable bool) error {
	switch {
	case model != AFM && model != ModelHaiku:
		return fmt.Errorf("模型只能是 afm 或 haiku（現在是 %s）", model)
	case model == ModelHaiku && !haikuAvailable:
		return ErrNoAPIKey
	case budget < 0 || math.IsNaN(budget) || math.IsInf(budget, 0):
		return errors.New("每月費用上限要是 0 以上的數字")
	}
	return nil
}

// Apply 換掉模型開關和每月上限，下一次呼叫就生效。沒過 CheckSetting 時回錯誤，設定不變
func (s *Switch) Apply(model string, budget float64) error {
	if err := CheckSetting(model, budget, s.HaikuAvailable()); err != nil {
		return err
	}
	s.setMu.Lock()
	defer s.setMu.Unlock()
	s.UseHaiku, s.Budget = model == ModelHaiku, budget
	return nil
}

// 費用上限的月份照台北時間算（台灣沒有日光節約時間，用固定時區，開發容器不用裝時區資料）
var taipei = time.FixedZone("Asia/Taipei", 8*60*60)

// MonthRange 是 now 所在月份（台北時間）的 [月初, 下個月初)：費用上限和後台的本月費用都照這個算
func MonthRange(now time.Time) (from, to time.Time) {
	now = now.In(taipei)
	from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, taipei)
	return from, from.AddDate(0, 1, 0)
}

func (s *Switch) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	m := s.AFM
	s.setMu.RLock()
	useHaiku, budget := s.UseHaiku, s.Budget
	s.setMu.RUnlock()
	if useHaiku && s.Haiku != nil && s.withinBudget(ctx, budget) {
		m = s.Haiku
	}
	reply, usage := m.Respond(ctx, r)
	s.record(ctx, r, reply, usage)
	return reply, usage
}

// withinBudget 查本月的 Haiku 累計費用；到上限時每個月記一次 log、通知一次。
// 查不到（資料庫錯誤）時照樣用 Haiku，只記錯誤：不讓資料庫問題擋住回覆
func (s *Switch) withinBudget(ctx context.Context, budget float64) bool {
	if budget <= 0 {
		return true
	}
	from, to := MonthRange(s.Now())
	spent, err := s.Ledger.CostBetween(ctx, from, to)
	if err != nil {
		s.Log.Error("查本月 AI 費用失敗，照樣用 Haiku", "error", err.Error())
		return true
	}
	if spent < budget {
		return true
	}
	s.mu.Lock()
	first := !s.warnedFrom.Equal(from)
	s.warnedFrom = from
	s.mu.Unlock()
	if first {
		s.Log.Warn(fmt.Sprintf("Haiku 本月費用 $%.2f 已達上限 $%.2f，自動切回本機模型", spent, budget))
		if s.OnBudgetExceeded != nil {
			s.OnBudgetExceeded(ctx, spent, budget)
		}
	}
	return false
}

// record 記一筆用量；沒有實際呼叫模型（快取命中、連不上）不記。model id 不在價格表上時費用記 0（Build 啟動時警告過），不擋呼叫
func (s *Switch) record(ctx context.Context, r Request, reply *Reply, u Usage) {
	if u.Model == "" {
		return
	}
	cost, _ := Cost(u)
	err := s.Ledger.RecordUsage(ctx, store.UsageRecord{
		Time: s.Now(), Model: u.Model, Request: r.Name,
		InputTokens: u.InputTokens, OutputTokens: u.OutputTokens, CostUSD: cost, OK: reply != nil,
	})
	if err != nil {
		s.Log.Error("記錄 AI 用量失敗", "error", err.Error())
	}
}
