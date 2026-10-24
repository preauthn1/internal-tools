package model

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"testing"
	"time"

	"aichatbot/internal/store"
)

// memLedger 是記憶體版的用量紀錄
type memLedger struct {
	mu      sync.Mutex
	records []store.UsageRecord
}

func (l *memLedger) RecordUsage(ctx context.Context, u store.UsageRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, u)
	return nil
}

func (l *memLedger) CostBetween(ctx context.Context, from, to time.Time) (float64, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var sum float64
	for _, u := range l.records {
		if !u.Time.Before(from) && u.Time.Before(to) {
			sum += u.CostUSD
		}
	}
	return sum, nil
}

func answers(n int, intent string) map[string][]map[string]Value {
	var list []map[string]Value
	for range n {
		list = append(list, map[string]Value{"intent": Text(intent)})
	}
	return map[string][]map[string]Value{"Intent": list}
}

// 設成 Haiku、還沒到上限：用 Haiku，記下 token 和費用（claude-haiku-5-5 每百萬 token 輸入 $0.10、輸出 $0.50）
func TestSwitchUsesHaikuAndRecordsCost(t *testing.T) {
	haiku := NewScripted(answers(1, "booking"))
	haiku.Usage = Usage{Model: "claude-haiku-5-5", InputTokens: 2_000_000, OutputTokens: 1_000_000}
	afm := NewScripted(answers(1, "other"))
	ledger := &memLedger{}
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	s := &Switch{Haiku: haiku, AFM: afm, UseHaiku: true, Budget: 10, Ledger: ledger, Now: func() time.Time { return now }}

	reply, _ := s.Respond(context.Background(), intentRequest)
	if v, _ := reply.String("intent"); v != "booking" {
		t.Errorf("應該用 Haiku 回答：got %q", v)
	}
	if len(ledger.records) != 1 {
		t.Fatalf("記了 %d 筆", len(ledger.records))
	}
	want := store.UsageRecord{Time: now, Model: "claude-haiku-5-5", Request: "Intent", InputTokens: 2_000_000, OutputTokens: 1_000_000, CostUSD: 0.7, OK: true}
	if got := ledger.records[0]; got != want {
		t.Errorf("got  %+v\nwant %+v", got, want)
	}
}

// 沒有實際呼叫模型（快取命中）不記用量；上限 0 表示不限
func TestSwitchSkipsCacheHitsAndZeroBudgetIsUnlimited(t *testing.T) {
	ctx := context.Background()
	haiku := NewScripted(answers(1, "booking"))
	haiku.Usage = Usage{Model: "claude-haiku-5-5", InputTokens: 10}
	ledger := &memLedger{records: []store.UsageRecord{{Time: time.Now(), Model: "claude-haiku-5-5", CostUSD: 1000}}}
	s := &Switch{Haiku: NewCache(haiku, t.TempDir(), fixedIdentity(haikuID), &CacheStats{}), AFM: NewScripted(nil),
		UseHaiku: true, Budget: 0, Ledger: ledger, Now: time.Now}
	for i := range 2 {
		if reply, _ := s.Respond(ctx, intentRequest); reply == nil {
			t.Fatalf("第 %d 次：上限 0 應該照樣用 Haiku", i+1)
		}
	}
	if len(ledger.records) != 2 {
		t.Errorf("應該只記實際呼叫的那一次，記了 %d 筆：%+v", len(ledger.records)-1, ledger.records[1:])
	}
}

// 本月（台北時間）Haiku 費用到上限：自動切回本機模型，記一次 log、通知一次（員工群組）；下個月恢復用 Haiku
func TestSwitchOverBudgetFallsBackToAFM(t *testing.T) {
	haiku := NewScripted(answers(2, "booking"))
	haiku.Usage = Usage{Model: "claude-haiku-5-5"}
	afm := NewScripted(answers(2, "other"))
	afm.Usage = Usage{Model: AFM}
	// 台北 10/1 00:30 是 UTC 9/30 16:30：算十月
	ledger := &memLedger{records: []store.UsageRecord{
		{Time: time.Date(2026, 9, 30, 16, 30, 0, 0, time.UTC), Model: "claude-haiku-5-5", CostUSD: 10},
		{Time: time.Date(2026, 9, 30, 15, 59, 0, 0, time.UTC), Model: "claude-haiku-5-5", CostUSD: 100}, // 台北九月
	}}
	var logs bytes.Buffer
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	s := &Switch{Haiku: haiku, AFM: afm, UseHaiku: true, Budget: 10, Ledger: ledger,
		Log: slog.New(slog.NewTextHandler(&logs, nil)), Now: func() time.Time { return now }}
	var notified []string
	s.OnBudgetExceeded = func(_ context.Context, spent, budget float64) {
		notified = append(notified, fmt.Sprint(spent, "/", budget))
	}

	for range 2 {
		reply, _ := s.Respond(context.Background(), intentRequest)
		if v, _ := reply.String("intent"); v != "other" {
			t.Errorf("超過上限應該用本機模型：got %q", v)
		}
	}
	if n := strings.Count(logs.String(), "已達上限"); n != 1 {
		t.Errorf("上限通知記了 %d 次：%s", n, logs.String())
	}
	if fmt.Sprint(notified) != "[10/10]" {
		t.Errorf("通知：%v", notified)
	}
	if last := ledger.records[len(ledger.records)-1]; last.Model != AFM || last.CostUSD != 0 {
		t.Errorf("本機模型的用量：%+v", last)
	}

	// 台北十一月
	now = time.Date(2026, 10, 31, 16, 0, 0, 0, time.UTC)
	reply, _ := s.Respond(context.Background(), intentRequest)
	if v, _ := reply.String("intent"); v != "booking" {
		t.Errorf("換月後應該恢復用 Haiku：got %q", v)
	}
}

// 後台改模型開關和上限：馬上生效；選 haiku 但沒有 API key、或不認得的值時回錯誤，設定不變
func TestSwitchApply(t *testing.T) {
	haiku := NewScripted(answers(1, "booking"))
	s := &Switch{Haiku: haiku, AFM: NewScripted(answers(1, "other")), Budget: 10, Ledger: &memLedger{}, Now: time.Now}
	if !s.HaikuAvailable() {
		t.Error("有 Haiku 時應該可以切")
	}
	if err := s.Apply(ModelHaiku, 3.5); err != nil {
		t.Fatal(err)
	}
	if m, b := s.Settings(); m != ModelHaiku || b != 3.5 {
		t.Errorf("Settings：%s %v", m, b)
	}
	if reply, _ := s.Respond(context.Background(), intentRequest); reply == nil {
		t.Error("切到 Haiku 後應該用 Haiku 回答")
	}
	if err := s.Apply("gpt", 1); err == nil {
		t.Error("不認得的模型應該失敗")
	}

	noKey := &Switch{AFM: NewScripted(nil), Budget: 10}
	if noKey.HaikuAvailable() {
		t.Error("沒有 Haiku 時不能切")
	}
	if err := noKey.Apply(ModelHaiku, 1); err == nil {
		t.Error("沒有 API key 切到 haiku 應該失敗")
	}
	if m, b := noKey.Settings(); m != AFM || b != 10 {
		t.Errorf("失敗時設定不該變：%s %v", m, b)
	}
}

// 後台存檔（Apply）和機器人回覆（Respond）同時進行不會互相干擾（go test -race）
func TestSwitchApplyWhileResponding(t *testing.T) {
	s := &Switch{Haiku: NewScripted(nil), AFM: NewScripted(nil), Budget: 10, Ledger: &memLedger{}, Now: time.Now}
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() { s.Respond(context.Background(), intentRequest) })
		wg.Go(func() { s.Apply([]string{AFM, ModelHaiku}[i%2], float64(i)) })
	}
	wg.Wait()
}
