package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"testing"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// recordPending 直接在資料庫記一次待補問題（文字相同才合併），回傳編號
func (f *fixture) recordPending(text, userID string, at time.Time) int64 {
	f.t.Helper()
	id, _, err := f.db.RecordPending(context.Background(), store.PendingHit{Time: at, Platform: "telegram", UserID: userID, ChatID: "c-" + userID, Text: text, Lang: "zh", Source: "unanswerable"},
		nil, func(other string, _ []float32) bool { return other == text })
	if err != nil {
		f.t.Fatal(err)
	}
	return id
}

type pendingView struct {
	ID          int64
	Text        string
	Count       int
	Status      string
	Resolution  string
	StaffAnswer string `json:"staff_answer"`
	Hits        []struct {
		UserID       string `json:"user_id"`
		Text         string
		Conversation []struct {
			Direction, Text string
			SendFailed      bool `json:"send_failed"`
		}
	}
}

func (f *fixture) getJSON(path string, cookie *http.Cookie, v any) {
	f.t.Helper()
	w := f.do("GET", path, "", cookie)
	if w.Code != http.StatusOK {
		f.t.Fatalf("GET %s：%d %s", path, w.Code, w.Body)
	}
	if err := json.Unmarshal(w.Body.Bytes(), v); err != nil {
		f.t.Fatal(err)
	}
}

// newPendingFixture：存好商家設定（有一題 FAQ、一段知識）的 fixture
func newPendingFixture(t *testing.T) (*fixture, *http.Cookie) {
	t.Helper()
	f := newFixture(t)
	c := merchant.Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "02-1234-5678"
	c.Knowledge = []merchant.Section{{Title: "交通", Body: "捷運站走路 5 分鐘"}}
	c.FAQ = []merchant.FAQ{{Key: "faq-1", Question: "可以刷卡嗎", AltQuestions: []string{}, Answer: "可以"}}
	if _, err := merchant.Save(context.Background(), f.db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	return f, f.login()
}

func (f *fixture) pending(id int64, cookie *http.Cookie) pendingView {
	f.t.Helper()
	var q pendingView
	f.getJSON(fmt.Sprintf("/api/pending-questions/%d", id), cookie, &q)
	return q
}

// 轉成 FAQ：在 FAQ 最後新增一題，待補問題改成已處理（faq）；答案空白回 400，處理過的再轉回 409
func TestPendingToFAQ(t *testing.T) {
	f, cookie := newPendingFixture(t)
	id := f.recordPending("有停車場嗎", "u1", f.now)
	path := fmt.Sprintf("/api/pending-questions/%d/faq", id)
	if w := f.do("POST", path, `{"question":"有停車場嗎","alt_questions":[],"answer":" "}`, cookie); w.Code != http.StatusBadRequest {
		t.Errorf("答案空白：%d %s", w.Code, w.Body)
	}
	w := f.do("POST", path, `{"question":"有停車場嗎","alt_questions":["哪裡停車"],"answer":"對面有付費停車場"}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("轉成 FAQ：%d %s", w.Code, w.Body)
	}
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	if len(c.FAQ) != 2 || c.FAQ[1].Question != "有停車場嗎" || c.FAQ[1].Answer != "對面有付費停車場" || c.FAQ[1].AltQuestions[0] != "哪裡停車" {
		t.Errorf("FAQ：%+v", c.FAQ)
	}
	if q := f.pending(id, cookie); q.Status != "resolved" || q.Resolution != "faq" {
		t.Errorf("待補問題：%+v", q)
	}
	if w := f.do("POST", path, `{"question":"有停車場嗎","alt_questions":[],"answer":"有"}`, cookie); w.Code != http.StatusConflict {
		t.Errorf("處理過的再轉：%d %s", w.Code, w.Body)
	}
}

// 加進知識段落：接在現有段落後面（或新增一段），待補問題改成已處理（knowledge）；內容空白回 400
func TestPendingToKnowledge(t *testing.T) {
	f, cookie := newPendingFixture(t)
	id := f.recordPending("有停車場嗎", "u1", f.now)
	path := fmt.Sprintf("/api/pending-questions/%d/knowledge", id)
	if w := f.do("POST", path, `{"title":"交通","body":" "}`, cookie); w.Code != http.StatusBadRequest {
		t.Errorf("內容空白：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", path, `{"title":"交通","body":"對面有付費停車場"}`, cookie); w.Code != http.StatusOK {
		t.Fatalf("加進知識：%d %s", w.Code, w.Body)
	}
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Knowledge) != 1 || c.Knowledge[0].Body != "捷運站走路 5 分鐘\n對面有付費停車場" {
		t.Errorf("知識段落：%+v", c.Knowledge)
	}
	if q := f.pending(id, cookie); q.Status != "resolved" || q.Resolution != "knowledge" {
		t.Errorf("待補問題：%+v", q)
	}
}

// 加成意圖案例：和新增案例一樣檢查題庫重複（409，待補問題不變），存進去後待補問題改成已處理（example）
func TestPendingToExample(t *testing.T) {
	f, _ := newExampleFixture(t)
	cookie := f.login()
	id := f.recordPending("你們的貓叫什麼名字", "u1", f.now)
	path := fmt.Sprintf("/api/pending-questions/%d/example", id)
	for _, text := range firstSentences(t) {
		if w := f.do("POST", path, exampleBody(text, "general"), cookie); w.Code != http.StatusConflict {
			t.Errorf("和題庫重複：%d %s", w.Code, w.Body)
		}
		break
	}
	if q := f.pending(id, cookie); q.Status != "open" {
		t.Errorf("存不進去時待補問題變了：%+v", q)
	}
	if w := f.do("POST", path, exampleBody("你們的貓叫什麼名字", "general"), cookie); w.Code != http.StatusOK {
		t.Fatalf("加成案例：%d %s", w.Code, w.Body)
	}
	if !slices.ContainsFunc(f.listExamples(cookie), func(e example) bool { return e.Text == "你們的貓叫什麼名字" && e.Intent == "general" }) {
		t.Error("案例庫裡沒有")
	}
	if q := f.pending(id, cookie); q.Status != "resolved" || q.Resolution != "example" {
		t.Errorf("待補問題：%+v", q)
	}
}

// 忽略：改成已忽略，不在待處理列表（?status=ignored 才列得出來）
func TestPendingIgnore(t *testing.T) {
	f, cookie := newPendingFixture(t)
	id := f.recordPending("今天天氣好嗎", "u1", f.now)
	if w := f.do("POST", fmt.Sprintf("/api/pending-questions/%d/ignore", id), "{}", cookie); w.Code != http.StatusOK {
		t.Fatalf("忽略：%d %s", w.Code, w.Body)
	}
	var open, ignored []pendingView
	f.getJSON("/api/pending-questions", cookie, &open)
	f.getJSON("/api/pending-questions?status=ignored", cookie, &ignored)
	if len(open) != 0 || len(ignored) != 1 || ignored[0].ID != id || ignored[0].Status != "ignored" {
		t.Errorf("待處理 %+v，已忽略 %+v", open, ignored)
	}
}

// 詳細頁：還沒翻成商家慣用語言的（背景翻譯失敗）這時用翻譯顧客訊息的提示詞補翻、存起來；慣用語言的不翻
func TestPendingDetailTranslates(t *testing.T) {
	f, cookie := newPendingFixture(t)
	chat := model.NewScripted(map[string][]map[string]model.Value{"TranslateCustomer": {{"translation": model.Text("哪裡可以停車？")}}})
	f.srv.Chat = chat
	en, _, err := f.db.RecordPending(context.Background(), store.PendingHit{Time: f.now, Platform: "telegram", UserID: "u1", Text: "Where can I park?", Lang: "en", Source: "unclear"},
		nil, func(string, []float32) bool { return false })
	if err != nil {
		t.Fatal(err)
	}
	zh := f.recordPending("有停車場嗎", "u2", f.now)
	var q struct{ Translated string }
	f.getJSON(fmt.Sprintf("/api/pending-questions/%d", en), cookie, &q)
	if q.Translated != "哪裡可以停車？" {
		t.Errorf("補翻：%+v", q)
	}
	f.getJSON(fmt.Sprintf("/api/pending-questions/%d", en), cookie, &q)
	f.getJSON(fmt.Sprintf("/api/pending-questions/%d", zh), cookie, &q)
	if reqs := chat.Requests(); len(reqs) != 1 {
		t.Errorf("翻譯了 %d 次", len(reqs))
	}
}

// 端到端（假模型）：答不出來的問題出現在待補問題，轉成 FAQ（問題預填顧客原話）後同一句再問，回 FAQ 的答案
func TestPendingFAQThenAnswered(t *testing.T) {
	f, cookie := newPendingFixture(t)
	general := map[string]model.Value{"intent": model.Text("general"), "self_name": model.Text("")}
	faq := func(key string) map[string]model.Value {
		return map[string]model.Value{"ask": model.Text("停車"), "faq": model.Text(key), "kind": model.Text("question")}
	}
	chat := model.NewScripted(map[string][]map[string]model.Value{
		"Intent": {general, general},
		"FAQ":    {faq("none"), faq("faq-2")},
		"Pick":   {{"ask": model.Text("停車"), "pick1": model.Text("none"), "pick2": model.Text("none"), "pick3": model.Text("none")}},
		"Ask":    {{"kind": model.Text("asks_about_store")}},
	})
	bot := &dialog.Bot{DB: f.db, Model: chat, Now: func() time.Time { return f.now }}
	ask := func() dialog.Reply {
		t.Helper()
		r, err := bot.Receive(context.Background(), dialog.Message{Platform: "telegram", UserID: "u1", ChatID: "c1", Text: "有停車場嗎"})
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	if r := ask(); r.Key == "answer" {
		t.Fatalf("第一次就答出來了：%+v", r)
	}
	var list []pendingView
	f.getJSON("/api/pending-questions", cookie, &list)
	if len(list) != 1 || list[0].Text != "有停車場嗎" {
		t.Fatalf("待補問題：%+v", list)
	}
	body, _ := json.Marshal(map[string]any{"question": list[0].Text, "alt_questions": []string{}, "answer": "對面有付費停車場"})
	if w := f.do("POST", fmt.Sprintf("/api/pending-questions/%d/faq", list[0].ID), string(body), cookie); w.Code != http.StatusOK {
		t.Fatalf("轉成 FAQ：%d %s", w.Code, w.Body)
	}
	// 轉人工中機器人不回覆：員工結束轉人工後再問
	if c, _ := bot.LoadCustomer(context.Background(), "telegram", "u1"); c.Human {
		if _, err := bot.EndHandoff(context.Background(), c.HandoffTicket); err != nil {
			t.Fatal(err)
		}
	}
	if r := ask(); r.Key != "answer" || r.Text != "對面有付費停車場" {
		t.Errorf("轉成 FAQ 後再問：%+v", r)
	}
}

// 列表依次數排序；詳細頁有每次出現和那位顧客前後的對話（最多前後各 pendingContext 則）
func TestPendingListAndDetail(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	ctx := context.Background()
	for i := range 2*pendingContext + 2 {
		at := f.now.Add(time.Duration(i-pendingContext) * time.Minute)
		e := store.LogEntry{Time: at, Platform: "telegram", UserID: "u1", ChatID: "c-u1", Direction: "customer", Text: fmt.Sprintf("第 %d 則", i)}
		if i == pendingContext+1 { // 送不出去的回覆照樣列出，附 send_failed
			e.Direction, e.SendFailed = "bot", true
		}
		if err := f.db.AppendLog(ctx, e); err != nil {
			t.Fatal(err)
		}
	}
	dog := f.recordPending("可以帶狗嗎", "u2", f.now)
	park := f.recordPending("有停車場嗎", "u1", f.now)
	f.recordPending("有停車場嗎", "u3", f.now.Add(time.Hour))

	var list []pendingView
	f.getJSON("/api/pending-questions", cookie, &list)
	if len(list) != 2 || list[0].ID != park || list[0].Count != 2 || list[1].ID != dog {
		t.Fatalf("列表：%+v", list)
	}
	var q pendingView
	f.getJSON(fmt.Sprintf("/api/pending-questions/%d", park), cookie, &q)
	if q.Text != "有停車場嗎" || len(q.Hits) != 2 || q.Hits[0].UserID != "u1" {
		t.Fatalf("詳細：%+v", q)
	}
	conv := q.Hits[0].Conversation
	if len(conv) != 2*pendingContext || conv[0].Text != "第 1 則" || conv[len(conv)-1].Text != fmt.Sprintf("第 %d 則", 2*pendingContext) {
		t.Errorf("前後的對話：%+v", conv)
	}
	if !conv[pendingContext].SendFailed || conv[pendingContext-1].SendFailed {
		t.Errorf("送出失敗：%+v", conv)
	}
	if w := f.do("GET", "/api/pending-questions/999", "", cookie); w.Code != http.StatusNotFound {
		t.Errorf("沒有這筆：%d", w.Code)
	}
}
