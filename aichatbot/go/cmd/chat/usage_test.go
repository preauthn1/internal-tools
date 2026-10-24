package main

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// AI 用量：加總實際呼叫模型的 token 和估算費用（快取命中沒有用量，不算），量測結束時印出（eval/verify.py 讀）
func TestUsageTally(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	haiku := model.NewScripted(map[string][]map[string]model.Value{"Intent": {generalIntent, generalIntent}})
	haiku.Usage = model.Usage{Model: "claude-haiku-5-5", InputTokens: 1000, OutputTokens: 200}
	tally := &usageTally{Ledger: db}
	s := &model.Switch{Haiku: haiku, AFM: haiku, UseHaiku: true, Ledger: tally, Now: time.Now}
	for range 2 {
		s.Respond(context.Background(), model.Request{Name: "Intent"})
	}
	if got, want := tally.Summary(), "AI 用量：輸入 2000 token、輸出 400 token、估算 US$0.0004"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
