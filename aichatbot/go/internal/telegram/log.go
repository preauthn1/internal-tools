package telegram

import (
	"io"
	"log/slog"
	"strings"
	"sync"
)

// NewLogger 建立伺服器日誌：一行一筆，寫出去之前把 token、secret、API key 這類秘密換成 ***
// （最後一道防線；呼叫的地方本來就不該把它們放進日誌）
func NewLogger(w io.Writer, secrets ...string) *slog.Logger {
	return slog.New(slog.NewTextHandler(&redactWriter{w: w, secrets: nonEmpty(secrets)}, nil))
}

type redactWriter struct {
	mu      sync.Mutex
	w       io.Writer
	secrets []string
}

// slog 一筆紀錄呼叫一次 Write，所以 secret 不會被切在兩次 Write 之間
func (r *redactWriter) Write(p []byte) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, err := io.WriteString(r.w, redact(string(p), r.secrets...)); err != nil {
		return 0, err
	}
	return len(p), nil
}

func redact(s string, secrets ...string) string {
	for _, secret := range secrets {
		if secret != "" {
			s = strings.ReplaceAll(s, secret, "***")
		}
	}
	return s
}

func nonEmpty(list []string) []string {
	var out []string
	for _, s := range list {
		if s != "" {
			out = append(out, s)
		}
	}
	return out
}
