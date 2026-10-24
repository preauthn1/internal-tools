package model

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var intentRequest = Request{
	Name: "Intent", Text: "我要訂位", Instructions: "判斷顧客想做什麼",
	Fields: []Field{
		{Kind: KindChoice, Name: "intent", Description: "意圖", Options: []string{"booking", "other"}},
		{Kind: KindInteger, Name: "people", Description: "人數"},
	},
}

// Haiku：指示放 system、顧客的話放 user，不送 temperature（Haiku 5.5 只接受預設值，送了會 400）、關閉思考（max_tokens 只有 300），用 structured outputs 強制照 schema 回答；讀出各欄位和 token 數
func TestHaikuRequestAndReply(t *testing.T) {
	var got map[string]any
	var header http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/messages" {
			t.Errorf("got %s %s", r.Method, r.URL.Path)
		}
		header = r.Header
		body, _ := io.ReadAll(r.Body)
		json.Unmarshal(body, &got)
		w.Write([]byte(`{"type":"message","stop_reason":"end_turn",
			"content":[{"type":"text","text":"{\"intent\":\"booking\",\"people\":4}"}],
			"usage":{"input_tokens":120,"output_tokens":15}}`))
	}))
	defer srv.Close()

	h := NewHaiku("sk-test", "claude-haiku-5-5", srv.URL)
	reply, usage := h.Respond(context.Background(), intentRequest)
	if reply == nil {
		t.Fatal("reply 是 nil")
	}
	if s, _ := reply.String("intent"); s != "booking" {
		t.Errorf("intent：got %q", s)
	}
	if n, _ := reply.Int("people"); n != 4 {
		t.Errorf("people：got %d", n)
	}
	if usage != (Usage{Model: "claude-haiku-5-5", InputTokens: 120, OutputTokens: 15}) {
		t.Errorf("usage：got %+v", usage)
	}
	if header.Get("x-api-key") != "sk-test" || header.Get("anthropic-version") != "2023-06-01" {
		t.Errorf("header 不對：%v", header)
	}
	if got["model"] != "claude-haiku-5-5" || got["system"] != "判斷顧客想做什麼" {
		t.Errorf("請求內容不對：%v", got)
	}
	if _, ok := got["temperature"]; ok {
		t.Errorf("不能送 temperature（Haiku 5.5 只接受預設值）：%v", got["temperature"])
	}
	if th, _ := got["thinking"].(map[string]any); th["type"] != "disabled" {
		t.Errorf("要關閉思考：%v", got["thinking"])
	}
	msgs, _ := got["messages"].([]any)
	if len(msgs) != 1 || msgs[0].(map[string]any)["role"] != "user" || msgs[0].(map[string]any)["content"] != "我要訂位" {
		t.Errorf("messages 不對：%v", got["messages"])
	}
	format := got["output_config"].(map[string]any)["format"].(map[string]any)
	schema, _ := json.Marshal(format["schema"])
	var wantSchema any
	json.Unmarshal(Schema(intentRequest), &wantSchema)
	want, _ := json.Marshal(wantSchema)
	if format["type"] != "json_schema" || string(schema) != string(want) {
		t.Errorf("output_config 不對：%v", format)
	}
}

// HTTP 錯誤、回覆壞掉、逾時都回 nil；回 200 但欄位讀不到時照樣回報 token（有計費）
func TestHaikuFailures(t *testing.T) {
	cases := []struct {
		name   string
		handle func(w http.ResponseWriter)
		usage  Usage
	}{
		{"HTTP 500", func(w http.ResponseWriter) { w.WriteHeader(500); w.Write([]byte(`{"type":"error"}`)) }, Usage{}},
		{"不是 JSON", func(w http.ResponseWriter) { w.Write([]byte(`<html>`)) }, Usage{}},
		{"欄位讀不到", func(w http.ResponseWriter) {
			w.Write([]byte(`{"stop_reason":"max_tokens","content":[{"type":"text","text":"{\"intent\":"}],"usage":{"input_tokens":5,"output_tokens":300}}`))
		}, Usage{Model: "m", InputTokens: 5, OutputTokens: 300}},
		{"逾時", func(w http.ResponseWriter) { time.Sleep(300 * time.Millisecond) }, Usage{}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { c.handle(w) }))
			defer srv.Close()
			h := NewHaiku("k", "m", srv.URL)
			h.Log = slog.New(slog.DiscardHandler)
			h.Timeout = 100 * time.Millisecond
			reply, usage := h.Respond(context.Background(), intentRequest)
			if reply != nil {
				t.Errorf("應該回 nil，got %+v", reply)
			}
			if usage != c.usage {
				t.Errorf("usage：got %+v, want %+v", usage, c.usage)
			}
		})
	}
}
