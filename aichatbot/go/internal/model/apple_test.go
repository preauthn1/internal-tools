package model

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// 假的 Swift 模型服務：/respond 收 {name, instructions, text, schema}、回 {values}；/version 回模型版本和 macOS 版本
func fakeModelServer(t *testing.T, respond http.HandlerFunc) (*httptest.Server, *int) {
	versions := 0
	mux := http.NewServeMux()
	mux.HandleFunc("POST /respond", respond)
	mux.HandleFunc("GET /version", func(w http.ResponseWriter, r *http.Request) {
		versions++
		w.Write([]byte(`{"model":"apple-fm/greedy/max150/timeout15","os":"Version 26.1 (Build 25B78)"}`))
	})
	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)
	return srv, &versions
}

func TestAppleRespond(t *testing.T) {
	var got map[string]json.RawMessage
	var startBy string
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		startBy = r.Header.Get("X-Start-By-Ms")
		json.NewDecoder(r.Body).Decode(&got)
		w.Write([]byte(`{"values":{"intent":"booking","people":2}}`))
	})
	a := NewApple(srv.URL)
	sent := time.Now()
	reply, usage := a.Respond(context.Background(), intentRequest)
	if reply == nil {
		t.Fatal("reply 是 nil")
	}
	// 最晚開始時間＝送出時間＋（Go 端逾時 − 模型服務自己的逾時）：在模型服務排隊超過這個時間，Go 已經等不到結果
	ms, err := strconv.ParseInt(startBy, 10, 64)
	want := sent.Add(a.Timeout - AppleModelTimeout)
	if err != nil || time.UnixMilli(ms).Sub(want).Abs() > time.Second {
		t.Errorf("X-Start-By-Ms：got %q，應該約 %d", startBy, want.UnixMilli())
	}
	if s, _ := reply.String("intent"); s != "booking" {
		t.Errorf("intent：got %q", s)
	}
	if n, _ := reply.Int("people"); n != 2 {
		t.Errorf("people：got %d", n)
	}
	if usage != (Usage{Model: AFM}) {
		t.Errorf("usage：got %+v", usage)
	}
	if string(got["name"]) != `"Intent"` || string(got["text"]) != `"我要訂位"` || string(got["instructions"]) != `"判斷顧客想做什麼"` ||
		string(got["schema"]) != string(Schema(intentRequest)) {
		t.Errorf("請求內容不對：%s", got)
	}
}

// 模型服務回了 token 數就記進用量（費用照舊 0）；沒回（舊版模型服務）時是 0，見 TestAppleRespond
func TestAppleRespondUsage(t *testing.T) {
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"values":{"intent":"booking"},"usage":{"input_tokens":312,"output_tokens":9}}`))
	})
	reply, usage := NewApple(srv.URL).Respond(context.Background(), intentRequest)
	if reply == nil {
		t.Fatal("reply 是 nil")
	}
	if usage != (Usage{Model: AFM, InputTokens: 312, OutputTokens: 9}) {
		t.Errorf("usage：got %+v", usage)
	}
}

// 模型服務排隊滿了（503）、排太久沒執行（504）都當失敗：回 nil，算有送到模型服務
func TestAppleBusyOrExpired(t *testing.T) {
	for _, status := range []int{http.StatusServiceUnavailable, http.StatusGatewayTimeout} {
		srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(status)
			w.Write([]byte(`{"error":"忙碌"}`))
		})
		a := NewApple(srv.URL)
		a.Log = slog.New(slog.DiscardHandler)
		if reply, usage := a.Respond(context.Background(), intentRequest); reply != nil || usage != (Usage{Model: AFM}) {
			t.Errorf("HTTP %d：got %+v %+v", status, reply, usage)
		}
	}
}

// 模型服務回錯誤（本機模型失敗或逾時）、連不上都回 nil
func TestAppleFailures(t *testing.T) {
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"本機模型沒有回答"}`))
	})
	a := NewApple(srv.URL)
	a.Log = slog.New(slog.DiscardHandler)
	if reply, usage := a.Respond(context.Background(), intentRequest); reply != nil || usage != (Usage{Model: AFM}) {
		t.Errorf("502：got %+v %+v", reply, usage)
	}
	down := NewApple("http://0.0.0.0:1")
	down.Log = slog.New(slog.DiscardHandler)
	if reply, usage := down.Respond(context.Background(), intentRequest); reply != nil || usage != (Usage{}) {
		t.Errorf("連不上：got %+v %+v", reply, usage)
	}
}

// 提示詞超過本機模型上下文（模型服務的 502 帶 reason context_exceeded）：回 nil，Failure 是 too_long，看得出不是一般的 AI 失敗
func TestAppleContextExceeded(t *testing.T) {
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"提示詞超過本機模型上下文","reason":"context_exceeded"}`))
	})
	a := NewApple(srv.URL)
	a.Log = slog.New(slog.DiscardHandler)
	if reply, usage := a.Respond(context.Background(), intentRequest); reply != nil || usage != (Usage{Model: AFM, Failure: FailureContextExceeded}) {
		t.Errorf("got %+v %+v", reply, usage)
	}
}

// 被安全過濾擋下（模型服務的 502 帶 reason guardrail）：回 nil，Failure 是 blocked_by_safety，和逾時、太長、一般失敗分得開（#68）
func TestAppleBlockedBySafety(t *testing.T) {
	srv, _ := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(`{"error":"被安全過濾擋下","reason":"guardrail"}`))
	})
	a := NewApple(srv.URL)
	a.Log = slog.New(slog.DiscardHandler)
	if reply, usage := a.Respond(context.Background(), intentRequest); reply != nil || usage != (Usage{Model: AFM, Failure: FailureBlockedBySafety}) {
		t.Errorf("got %+v %+v", reply, usage)
	}
}

// 版本（快取 key 用）只問一次
func TestAppleIdentity(t *testing.T) {
	srv, versions := fakeModelServer(t, func(w http.ResponseWriter, r *http.Request) {})
	a := NewApple(srv.URL)
	for range 2 {
		id, err := a.Identity(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if id != (Identity{Model: AFM, Version: "apple-fm/greedy/max150/timeout15", OS: "Version 26.1 (Build 25B78)"}) {
			t.Errorf("got %+v", id)
		}
	}
	if *versions != 1 {
		t.Errorf("問了 %d 次版本", *versions)
	}
}
