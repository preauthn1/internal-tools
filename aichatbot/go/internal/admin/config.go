package admin

import (
	"fmt"
	"strconv"
)

// Config 是後台伺服器的設定：另開一個 port，只聽本機，對外由 cloudflared 用另一個網址轉進來（前面加 Cloudflare Access）
type Config struct {
	Host    string
	Port    int
	Dir     string // 建好的後台頁面
	EvalDir string // 題庫：新增意圖案例時檢查重複
}

// DefaultDir 是建好的後台頁面（相對工作目錄，從 repo 根目錄啟動）
const DefaultDir = "web/apps/web-antd/dist"

// DefaultEvalDir 是題庫（相對工作目錄，從 repo 根目錄啟動）
const DefaultEvalDir = "eval"

// ConfigFromEnv：BOT_ADMIN_PORT（預設 8081）、BOT_ADMIN_DIR（預設 DefaultDir）、BOT_EVAL_DIR（預設 DefaultEvalDir）都選填
func ConfigFromEnv(getenv func(string) string) (Config, error) {
	c := Config{Host: "0.0.0.0", Port: 8081, Dir: DefaultDir, EvalDir: DefaultEvalDir}
	if p := getenv("BOT_ADMIN_PORT"); p != "" {
		port, err := strconv.Atoi(p)
		if err != nil {
			return Config{}, fmt.Errorf("環境變數 BOT_ADMIN_PORT 不正確：要是數字")
		}
		c.Port = port
	}
	if d := getenv("BOT_ADMIN_DIR"); d != "" {
		c.Dir = d
	}
	if d := getenv("BOT_EVAL_DIR"); d != "" {
		c.EvalDir = d
	}
	return c, nil
}
