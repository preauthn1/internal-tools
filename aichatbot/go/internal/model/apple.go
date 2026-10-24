package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// AFM 是 Apple 本機模型在用量紀錄和快取 key 裡的名稱
const AFM = "afm"

// Identity 是模型的名稱和版本（AI 回答快取的 key 包含它）；OS 只有 Apple 本機模型有（系統更新可能換模型）
type Identity struct {
	Model   string
	Version string
	OS      string
}

// Apple 呼叫 Swift 模型服務（Sources/model-server），本機模型一次只能跑一個，要包在 Queued 裡面用
type Apple struct {
	base    string
	client  *http.Client
	Log     *slog.Logger
	Timeout time.Duration

	mu sync.Mutex
	id *Identity
}

// AppleModelTimeout 是模型服務那邊本機模型自己的逾時（Swift 的 AppleModel）
const AppleModelTimeout = 15 * time.Second

// 模型服務自己有 15 秒逾時，這裡多留一點傳輸時間
func NewApple(baseURL string) *Apple {
	return &Apple{base: baseURL, client: &http.Client{}, Log: slog.Default(), Timeout: 20 * time.Second}
}

func (a *Apple) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	ctx, cancel := context.WithTimeout(ctx, a.Timeout)
	defer cancel()
	body := map[string]any{"name": r.Name, "instructions": r.Instructions, "text": r.Text, "schema": Schema(r)}
	// 最晚開始時間：模型服務排隊排到這個時間還沒輪到，就算跑了 Go 這邊也等不到結果，模型服務直接放棄（回 504）
	startBy := time.Now().Add(a.Timeout - AppleModelTimeout)
	headers := map[string]string{"X-Start-By-Ms": strconv.FormatInt(startBy.UnixMilli(), 10)}
	data, answered, err := postJSON(ctx, a.client, a.base+"/respond", headers, body)
	var usage Usage
	if answered {
		usage = Usage{Model: AFM} // 有送到模型服務就記一筆；token 數成功時才有（舊版模型服務不回，是 0）
	}
	if err != nil {
		// 模型服務回 502 時 reason 是 context_exceeded：提示詞超過本機模型上下文；guardrail：被安全過濾擋下（#68）
		// （舊版模型服務不帶 reason，照一般失敗記）
		var failed struct {
			Reason string `json:"reason"`
		}
		if answered && json.Unmarshal(data, &failed) == nil {
			switch failed.Reason {
			case "context_exceeded":
				usage.Failure = FailureContextExceeded
				a.Log.Warn("提示詞超過本機模型上下文，請減少 FAQ 題數或縮短知識段落", "request", r.Name, "instructions_chars", len([]rune(r.Instructions)))
				return nil, usage
			case "guardrail":
				usage.Failure = FailureBlockedBySafety
				a.Log.Warn("被安全過濾擋下", "request", r.Name, "text", r.Text)
				return nil, usage
			}
		}
		a.fail(r, err)
		return nil, usage
	}
	var out struct {
		Values json.RawMessage `json:"values"`
		Usage  struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &out); err != nil {
		a.fail(r, fmt.Errorf("回覆不是 JSON：%w", err))
		return nil, usage
	}
	usage.InputTokens, usage.OutputTokens = out.Usage.InputTokens, out.Usage.OutputTokens
	values, ok := parseValues(r.Fields, out.Values)
	if !ok {
		a.fail(r, fmt.Errorf("讀不到欄位：%s", clip(data)))
		return nil, usage
	}
	return &Reply{Values: values}, usage
}

// 模型服務回 503（排隊滿了）、504（排太久沒執行）時錯誤訊息會帶 HTTP 狀態碼，log 看得出忙碌的頻率
func (a *Apple) fail(r Request, err error) {
	a.Log.Warn("Apple 本機模型呼叫失敗", "request", r.Name, "error", err.Error())
}

// Identity 問模型服務的模型版本和 macOS 版本；成功後記住，失敗下次再問
func (a *Apple) Identity(ctx context.Context) (Identity, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.id != nil {
		return *a.id, nil
	}
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.base+"/version", nil)
	if err != nil {
		return Identity{}, err
	}
	resp, err := a.client.Do(req)
	if err != nil {
		return Identity{}, err
	}
	defer resp.Body.Close()
	var v struct {
		Model string `json:"model"`
		OS    string `json:"os"`
	}
	if resp.StatusCode != http.StatusOK {
		return Identity{}, fmt.Errorf("模型服務 /version 回 HTTP %d", resp.StatusCode)
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil || v.Model == "" {
		return Identity{}, fmt.Errorf("模型服務 /version 的回覆不對")
	}
	a.id = &Identity{Model: AFM, Version: v.Model, OS: v.OS}
	return *a.id, nil
}
