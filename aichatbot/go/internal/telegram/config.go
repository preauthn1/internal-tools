// Package telegram 是 Telegram 平台層：webhook、Bot API 用戶端（含重試）、私訊的處理。
// 做法照搬 Swift 版（tag swift-final 的 Sources/Telegram）。
package telegram

import (
	"fmt"
	"strconv"
	"strings"
)

// Config 從環境變數讀；token、secret 不進 git、不寫進日誌
type Config struct {
	Token     string
	Secret    string // setWebhook 的 secret_token，Telegram 每個請求都會帶回來
	PublicURL string // 對外網址（Cloudflare Tunnel），webhook 是 <PublicURL>/telegram/webhook
	APIBase   string
	Host      string // 只聽本機，對外由 cloudflared 轉進來
	Port      int
}

const WebhookPath = "/telegram/webhook"

func (c Config) WebhookURL() string { return c.PublicURL + WebhookPath }

// ConfigFromEnv：TELEGRAM_BOT_TOKEN、TELEGRAM_WEBHOOK_SECRET、TELEGRAM_PUBLIC_URL 必填；
// TELEGRAM_API_BASE（測試用）、PORT 選填。錯誤訊息只說是哪個變數，不印值
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	var missing []string
	for _, name := range []string{"TELEGRAM_BOT_TOKEN", "TELEGRAM_WEBHOOK_SECRET", "TELEGRAM_PUBLIC_URL"} {
		if getenv(name) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return Config{}, fmt.Errorf("缺少環境變數：%s", strings.Join(missing, "、"))
	}
	c := Config{
		Token:     [REDACTED]"TELEGRAM_BOT_TOKEN"),
		Secret:    [REDACTED]"TELEGRAM_WEBHOOK_SECRET"),
		PublicURL: strings.TrimSuffix(getenv("TELEGRAM_PUBLIC_URL"), "/"),
		APIBase:   "https://example.invalid",
		Host:      "0.0.0.0",
		Port:      8080,
	}
	// Telegram 的規定：1～256 個字，只能用 A-Z、a-z、0-9、_、-
	if len(c.Secret) > 256 || strings.IndexFunc(c.Secret, func(r rune) bool {
		return !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '_' || r == '-')
	}) >= 0 {
		return Config{}, fmt.Errorf("環境變數 TELEGRAM_WEBHOOK_SECRET 不正確：只能用英數字、_、-，最多 256 個字")
	}
	if !strings.HasPrefix(c.PublicURL, "https://") {
		return Config{}, fmt.Errorf("環境變數 TELEGRAM_PUBLIC_URL 不正確：要用 https:// 開頭")
	}
	if base := getenv("TELEGRAM_API_BASE"); base != "" {
		c.APIBase = base
	}
	if p := getenv("PORT"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil {
			return Config{}, fmt.Errorf("環境變數 PORT 不正確：要是數字")
		}
		c.Port = port
	}
	return c, nil
}
