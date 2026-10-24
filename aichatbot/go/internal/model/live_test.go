package model

import (
	"context"
	"os"
	"testing"
)

// 實機驗收：Haiku 和 Apple 本機模型回答同一個選擇題。要設 BOT_LIVE_MODEL_TEST=1 才跑：
//
//	swift build -c release && .build/release/model-server &
//	BOT_LIVE_MODEL_TEST=1 ANTHROPIC_API_KEY=[REDACTED] go -C go test -run TestLive -v ./internal/model/
//
// 沒有 ANTHROPIC_API_KEY 時略過 Haiku；BOT_APPLE_MODEL_URL 預設 http://0.0.0.0:8765
func TestLiveSameChoiceQuestion(t *testing.T) {
	if os.Getenv("BOT_LIVE_MODEL_TEST") != "1" {
		t.Skip("要設 BOT_LIVE_MODEL_TEST=1 才跑")
	}
	cfg, err := ConfigFromEnv(os.Getenv)
	if err != nil {
		t.Fatal(err)
	}
	r := Request{
		Name: "Intent", Text: "我想訂明天晚上七點四位",
		Instructions: "你是餐廳客服。判斷顧客這句話的意圖：想預約訂位選 booking，其他選 other。",
		Fields:       []Field{{Kind: KindChoice, Name: "intent", Description: "顧客的意圖", Options: []string{"booking", "other"}}},
	}
	models := map[string]ChatModel{"afm": NewQueued(NewApple(cfg.AppleURL))}
	if cfg.APIKey != "" {
		models["haiku"] = NewHaiku(cfg.APIKey, cfg.HaikuModel, cfg.AnthropicURL)
	} else {
		t.Log("沒有 ANTHROPIC_API_KEY，略過 Haiku")
	}
	for name, m := range models {
		reply, usage := m.Respond(context.Background(), r)
		if reply == nil {
			t.Errorf("%s 沒有回答", name)
			continue
		}
		intent, _ := reply.String("intent")
		cost, _ := Cost(usage)
		t.Logf("%s：intent=%s，token %d／%d，$%.6f", name, intent, usage.InputTokens, usage.OutputTokens, cost)
		if intent != "booking" {
			t.Errorf("%s 回答 %q，應該是 booking", name, intent)
		}
	}
}
