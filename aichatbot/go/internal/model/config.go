package model

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ModelHaiku 是模型開關選 Haiku 時的值（另一個值是 AFM）
const ModelHaiku = "haiku"

// Config 從環境變數讀；模型開關和上限只是初始值，後台存過之後以資料庫為準（cmd/bot）。API key 不寫進日誌
type Config struct {
	Model         string // 模型開關：afm（預設）或 haiku
	APIKey        string
	HaikuModel    string
	MonthlyBudget float64 // Haiku 每月費用上限（美元），0 表示不限
	AppleURL      string  // Swift 模型服務
	AnthropicURL  string
	Cache         bool   // AI 回答快取（量測用，正式伺服器不開）
	CacheDir      string // 快取根目錄，底下再分命名空間
	CacheNS       string // dev 或 holdout
}

// ConfigFromEnv：BOT_MODEL、ANTHROPIC_API_KEY、BOT_HAIKU_MODEL、BOT_AI_MONTHLY_BUDGET_USD、BOT_APPLE_MODEL_URL、
// ANTHROPIC_BASE_URL（測試用）、BOT_AI_CACHE、BOT_AI_CACHE_DIR、BOT_AI_CACHE_NS，全部選填
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	c := Config{
		Model:         AFM,
		APIKey:        [REDACTED]"ANTHROPIC_API_KEY"),
		HaikuModel:    "claude-haiku-5-5",
		MonthlyBudget: 10,
		AppleURL:      "http://0.0.0.0:8765",
		AnthropicURL:  "https://example.invalid",
		Cache:         getenv("BOT_AI_CACHE") == "1",
		CacheNS:       "dev",
	}
	switch m := getenv("BOT_MODEL"); m {
	case "", AFM:
	case ModelHaiku:
		if c.APIKey == "" {
			return Config{}, errors.New("BOT_MODEL=haiku 但沒有設定 ANTHROPIC_API_KEY：先設好 API key，或改回 BOT_MODEL=afm（Apple 本機模型）")
		}
		c.Model = m
	default:
		return Config{}, fmt.Errorf("環境變數 BOT_MODEL 不正確：只能是 afm 或 haiku（現在是 %s）", m)
	}
	if v := getenv("BOT_HAIKU_MODEL"); v != "" {
		c.HaikuModel = v
	}
	if v := getenv("BOT_AI_MONTHLY_BUDGET_USD"); v != "" {
		b, err := strconv.ParseFloat(v, 64)
		if err != nil || b < 0 {
			return Config{}, errors.New("環境變數 BOT_AI_MONTHLY_BUDGET_USD 不正確：要是 0 以上的數字（美元，0 表示不限）")
		}
		c.MonthlyBudget = b
	}
	if v := getenv("BOT_APPLE_MODEL_URL"); v != "" {
		c.AppleURL = strings.TrimSuffix(v, "/")
	}
	if v := getenv("ANTHROPIC_BASE_URL"); v != "" {
		c.AnthropicURL = strings.TrimSuffix(v, "/")
	}
	if v := getenv("BOT_AI_CACHE_NS"); v != "" {
		c.CacheNS = v
	}
	if c.Cache && c.CacheNS != "dev" && c.CacheNS != "holdout" {
		return Config{}, fmt.Errorf("環境變數 BOT_AI_CACHE_NS 不正確：只能是 dev 或 holdout（現在是 %s）", c.CacheNS)
	}
	c.CacheDir = getenv("BOT_AI_CACHE_DIR")
	if c.CacheDir == "" {
		// Mac 上是 ~/Library/Caches/aichatbot/ai，和 Swift 版同一個位置
		dir, err := os.UserCacheDir()
		if err != nil {
			dir = os.TempDir()
		}
		c.CacheDir = filepath.Join(dir, "aichatbot", "ai")
	}
	return c, nil
}

// Build 組出對外用的模型：Switch 選 Cache(Queued(Apple)) 或 Cache(Haiku)，快取包在排隊外面（命中不用排隊）。
// 沒開快取時 stats 是 nil。回傳 *Switch：後台用 Apply 改模型開關和上限
func Build(c Config, ledger Ledger, log *slog.Logger) (*Switch, *CacheStats, error) {
	apple := NewApple(c.AppleURL)
	apple.Log = log
	var afm ChatModel = NewQueued(apple)
	var haiku ChatModel
	if c.APIKey != "" {
		h := NewHaiku(c.APIKey, c.HaikuModel, c.AnthropicURL)
		h.Log = log
		haiku = h
	}
	// Config 不一定來自 ConfigFromEnv，這裡再擋一次
	if c.Model == ModelHaiku && haiku == nil {
		return nil, nil, errors.New("模型開關是 haiku 但沒有 API key")
	}
	if _, ok := Cost(Usage{Model: c.HaikuModel}); haiku != nil && !ok {
		log.Warn("價格表沒有這個 Haiku model id，費用記 0、月上限不會生效", "model", c.HaikuModel)
	}
	var stats *CacheStats
	if c.Cache {
		stats = &CacheStats{}
		dir := filepath.Join(c.CacheDir, c.CacheNS)
		afm = NewCache(afm, dir, apple.Identity, stats)
		if haiku != nil {
			id := Identity{Model: c.HaikuModel, Version: HaikuVersion}
			haiku = NewCache(haiku, dir, func(context.Context) (Identity, error) { return id, nil }, stats)
		}
	}
	s := &Switch{Haiku: haiku, AFM: afm, UseHaiku: c.Model == ModelHaiku, Budget: c.MonthlyBudget, Ledger: ledger, Log: log, Now: time.Now}
	return s, stats, nil
}
