package model

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

func env(vars map[string]string) func(string) string {
	return func(k string) string { return vars[k] }
}

// 沒設定時用 AFM；Haiku 的 model id 預設 claude-haiku-5-5、每月上限 10 美元
func TestConfigDefaults(t *testing.T) {
	c, err := ConfigFromEnv(env(nil))
	if err != nil {
		t.Fatal(err)
	}
	want := Config{Model: AFM, HaikuModel: "claude-haiku-5-5", MonthlyBudget: 10,
		AppleURL: "http://0.0.0.0:8765", AnthropicURL: "https://example.invalid", CacheNS: "dev"}
	if c.CacheDir == "" {
		t.Error("快取資料夾沒有預設值")
	}
	c.CacheDir = ""
	if c != want {
		t.Errorf("got  %+v\nwant %+v", c, want)
	}
}

func TestConfigHaiku(t *testing.T) {
	c, err := ConfigFromEnv(env(map[string]string{
		"BOT_MODEL": "haiku", "ANTHROPIC_API_KEY": "sk-x", "BOT_HAIKU_MODEL": "claude-haiku-4-5",
		"BOT_AI_MONTHLY_BUDGET_USD": "0", "BOT_APPLE_MODEL_URL": "http://host.docker.internal:8765/",
		"BOT_AI_CACHE": "1", "BOT_AI_CACHE_DIR": "/tmp/ai", "BOT_AI_CACHE_NS": "holdout",
	}))
	if err != nil {
		t.Fatal(err)
	}
	want := Config{Model: "haiku", APIKey: [REDACTED]", HaikuModel: "claude-haiku-4-5", MonthlyBudget: 0,
		AppleURL: "http://host.docker.internal:8765", AnthropicURL: "https://example.invalid",
		Cache: true, CacheDir: "/tmp/ai", CacheNS: "holdout"}
	if c != want {
		t.Errorf("got  %+v\nwant %+v", c, want)
	}
}

// 組出整串（預設 AFM＋快取）：第一次呼叫模型服務並記用量，第二次命中快取、不呼叫也不記
func TestBuildAFMWithCache(t *testing.T) {
	calls := 0
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Write([]byte(`{"values":{"intent":"booking","people":2}}`))
	})
	cfg, _ := ConfigFromEnv(env(map[string]string{"BOT_APPLE_MODEL_URL": srv.URL, "BOT_AI_CACHE": "1", "BOT_AI_CACHE_DIR": t.TempDir()}))
	ledger := &memLedger{}
	m, stats, err := Build(cfg, ledger, slog.New(slog.DiscardHandler))
	if err != nil {
		t.Fatal(err)
	}
	for range 2 {
		reply, _ := m.Respond(context.Background(), intentRequest)
		if s, _ := reply.String("intent"); s != "booking" {
			t.Errorf("got %+v", reply)
		}
	}
	if calls != 1 || len(ledger.records) != 1 || ledger.records[0].Model != AFM {
		t.Errorf("模型服務呼叫 %d 次，用量 %+v", calls, ledger.records)
	}
	if stats == nil || stats.Summary() != "AI 快取：命中 1 次、實際呼叫 1 次" {
		t.Errorf("stats：%v", stats)
	}
	if _, stats, _ := Build(Config{Model: AFM, AppleURL: srv.URL}, ledger, slog.Default()); stats != nil {
		t.Error("沒開快取時 stats 應該是 nil")
	}
}

// dev、holdout 分開存，互不命中
func TestBuildCacheNamespaces(t *testing.T) {
	calls := 0
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.Write([]byte(`{"values":{"intent":"booking"}}`))
	})
	root := t.TempDir()
	for _, ns := range []string{"dev", "holdout", "dev"} {
		cfg, _ := ConfigFromEnv(env(map[string]string{"BOT_APPLE_MODEL_URL": srv.URL, "BOT_AI_CACHE": "1", "BOT_AI_CACHE_DIR": root, "BOT_AI_CACHE_NS": ns}))
		m, _, _ := Build(cfg, &memLedger{}, slog.New(slog.DiscardHandler))
		m.Respond(context.Background(), intentRequest)
	}
	if calls != 2 {
		t.Errorf("模型服務呼叫 %d 次，應該 2 次（dev 第二次命中）", calls)
	}
}

// 設定錯了啟動就失敗，錯誤訊息說明怎麼改；不印 API key
func TestConfigErrors(t *testing.T) {
	cases := map[string]struct {
		vars map[string]string
		want string
	}{
		"Haiku 沒有 key": {map[string]string{"BOT_MODEL": "haiku"}, "ANTHROPIC_API_KEY"},
		"不認得的模型":       {map[string]string{"BOT_MODEL": "gpt"}, "BOT_MODEL"},
		"上限不是數字":       {map[string]string{"BOT_AI_MONTHLY_BUDGET_USD": "十塊"}, "BOT_AI_MONTHLY_BUDGET_USD"},
		"上限是負數":        {map[string]string{"BOT_AI_MONTHLY_BUDGET_USD": "-1"}, "BOT_AI_MONTHLY_BUDGET_USD"},
		"快取命名空間不對":     {map[string]string{"BOT_AI_CACHE": "1", "BOT_AI_CACHE_NS": "prod"}, "BOT_AI_CACHE_NS"},
	}
	for name, c := range cases {
		_, err := ConfigFromEnv(env(c.vars))
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s：got %v", name, err)
		}
	}
}
