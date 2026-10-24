package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// HaikuVersion 是 Haiku 呼叫參數的版本：AI 回答快取的 key 包含它。改了下面的生成參數（思考、長度、逾時）要一起改
const HaikuVersion = "haiku/nothink/max300/timeout15"

// Haiku 直接用 net/http 呼叫 Anthropic Messages API（不引入 SDK）。可以同時呼叫，不排隊
type Haiku struct {
	key     string
	model   string
	base    string
	client  *http.Client
	Log     *slog.Logger
	Timeout time.Duration
}

func NewHaiku(apiKey, modelID, baseURL string) *Haiku {
	return &Haiku{key: apiKey, model: modelID, base: baseURL, client: &http.Client{}, Log: slog.Default(), Timeout: 15 * time.Second}
}

func (h *Haiku) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	ctx, cancel := context.WithTimeout(ctx, h.Timeout)
	defer cancel()
	body := map[string]any{
		"model":      h.model,
		"max_tokens": 300,
		// Haiku 5.5 預設會先思考，300 個 token 可能被思考用完；只做選擇題、抽欄位，不需要思考
		// 不送 temperature：Haiku 5.5 只接受預設值，送 0 會回 400（所以 Haiku 不是 greedy，見 README）
		"thinking": map[string]any{"type": "disabled"},
		"system":   r.Instructions,
		"messages": []map[string]string{{"role": "user", "content": r.Text}},
		// structured outputs：回答一定照 schema（GA，不用 beta header）
		"output_config": map[string]any{"format": map[string]any{"type": "json_schema", "schema": Schema(r)}},
	}
	headers := map[string]string{"x-api-key": h.key, "anthropic-version": "2023-06-01"}
	data, _, err := postJSON(ctx, h.client, h.base+"/v1/messages", headers, body)
	if err != nil {
		h.fail(r, err)
		return nil, Usage{}
	}
	var msg struct {
		StopReason string `json:"stop_reason"`
		Content    []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
		Usage struct {
			InputTokens  int `json:"input_tokens"`
			OutputTokens int `json:"output_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		h.fail(r, fmt.Errorf("回覆不是 JSON：%w", err))
		return nil, Usage{}
	}
	usage := Usage{Model: h.model, InputTokens: msg.Usage.InputTokens, OutputTokens: msg.Usage.OutputTokens}
	for _, c := range msg.Content {
		if c.Type == "text" {
			if values, ok := parseValues(r.Fields, []byte(c.Text)); ok {
				return &Reply{Values: values}, usage
			}
		}
	}
	h.fail(r, fmt.Errorf("讀不到欄位（stop_reason %s）", msg.StopReason))
	return nil, usage
}

// 記錄失敗原因；不印 API key 和顧客的話
func (h *Haiku) fail(r Request, err error) {
	h.Log.Warn("Haiku 呼叫失敗", "request", r.Name, "error", err.Error())
}
