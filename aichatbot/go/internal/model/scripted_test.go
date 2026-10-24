package model

import (
	"context"
	"testing"
)

// 假模型：同一題依序回答，腳本用完或沒寫到的題目回 nil，收到的請求依順序記下來
func TestScriptedAnswersInOrder(t *testing.T) {
	ctx := context.Background()
	m := NewScripted(map[string][]map[string]Value{
		"Intent": {{"intent": Text("booking")}, {"intent": Text("other"), "people": Integer(4)}},
	})
	r, _ := m.Respond(ctx, Request{Name: "Intent", Text: "我要訂位"})
	if got, _ := r.String("intent"); got != "booking" {
		t.Errorf("第 1 次：got %q", got)
	}
	r, _ = m.Respond(ctx, Request{Name: "Intent", Text: "四位"})
	if got, _ := r.String("intent"); got != "other" {
		t.Errorf("第 2 次：got %q", got)
	}
	if n, ok := r.Int("people"); !ok || n != 4 {
		t.Errorf("people：got %d %v", n, ok)
	}
	if _, ok := r.Int("intent"); ok {
		t.Error("文字欄位不該讀成整數")
	}
	if r, _ := m.Respond(ctx, Request{Name: "Intent"}); r != nil {
		t.Error("腳本用完應該回 nil")
	}
	if r, _ := m.Respond(ctx, Request{Name: "Booking"}); r != nil {
		t.Error("沒寫到的題目應該回 nil")
	}
	reqs := m.Requests()
	if len(reqs) != 4 || reqs[0].Text != "我要訂位" || reqs[3].Name != "Booking" {
		t.Errorf("記錄的請求不對：%+v", reqs)
	}
}
