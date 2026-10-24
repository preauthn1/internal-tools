package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// 預約、對話紀錄、AI 費用都要先登入
func TestViewingRequiresLogin(t *testing.T) {
	f := newFixture(t)
	for _, c := range []struct{ method, path, body string }{
		{"GET", "/api/bookings", ""},
		{"GET", "/api/conversations", ""},
		{"GET", "/api/conversations/telegram/u1", ""},
		{"POST", "/api/conversations/logs/1/translate", "{}"},
		{"POST", "/api/conversations/telegram/u1/end-handoff", "{}"},
		{"GET", "/api/usage", ""},
	} {
		if w := f.do(c.method, c.path, c.body); w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s 沒登入：%d", c.method, c.path, w.Code)
		}
	}
}

// appendLog 記一則對話（traces 是 JSON；空字串是沒有）
func (f *fixture) appendLog(userID, direction, text, traces string, at time.Time) {
	f.t.Helper()
	e := store.LogEntry{Time: at, Platform: "telegram", UserID: userID, ChatID: "c-" + userID, Direction: direction, Text: text, Traces: json.RawMessage(traces)}
	if err := f.db.AppendLog(context.Background(), e); err != nil {
		f.t.Fatal(err)
	}
}

type customerView struct {
	Platform string
	UserID   string    `json:"user_id"`
	LastText string    `json:"last_text"`
	LastTime time.Time `json:"last_time"`
	Messages int
}

// 顧客列表：最近有對話的在前，附最後一則有文字的訊息和訊息數；?q= 比對顧客編號或訊息內容
func TestConversationCustomers(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	var list []customerView
	f.getJSON("/api/conversations", cookie, &list)
	if list == nil || len(list) != 0 {
		t.Errorf("沒有對話：%+v", list)
	}
	f.appendLog("u1", "customer", "有停車場嗎", "", f.now.Add(-2*time.Hour))
	f.appendLog("u1", "bot", "對面有付費停車場", "", f.now.Add(-2*time.Hour))
	f.appendLog("u2", "customer", "可以帶狗嗎", "", f.now.Add(-time.Hour))
	f.appendLog("u2", "bot", "", "", f.now.Add(-time.Hour)) // 轉人工中沒回覆

	f.getJSON("/api/conversations", cookie, &list)
	if len(list) != 2 || list[0].UserID != "u2" || list[0].LastText != "可以帶狗嗎" || list[0].Messages != 2 || !list[0].LastTime.Equal(f.now.Add(-time.Hour)) ||
		list[1].UserID != "u1" || list[1].LastText != "對面有付費停車場" || list[1].Platform != "telegram" {
		t.Errorf("顧客列表：%+v", list)
	}
	for q, want := range map[string]string{"停車": "u1", "u2": "u2"} {
		f.getJSON("/api/conversations?q="+url.QueryEscape(q), cookie, &list)
		if len(list) != 1 || list[0].UserID != want {
			t.Errorf("搜尋 %s：%+v", q, list)
		}
	}
}

type conversationView struct {
	Lang      string
	Preferred string `json:"preferred_language"`
	Messages  []struct {
		ID        int64
		Direction string
		Text      string
		Traces    []struct {
			Step, Result string
			Detail       map[string]string
		}
		Translatable bool
		SendFailed   bool `json:"send_failed"`
	}
}

// 後台結束轉人工（#55）：對話內容回 human；結束後狀態歸零、對話紀錄記一筆，下一句 AI 恢復回覆。
// 不在轉人工中回 409，沒有接上對話層回 503
func TestEndHandoffFromConversation(t *testing.T) {
	ctx := context.Background()
	f, cookie := newPendingFixture(t)
	path := "/api/conversations/telegram/u1/end-handoff"
	if w := f.do("POST", path, "{}", cookie); w.Code != http.StatusServiceUnavailable {
		t.Errorf("沒有對話層：%d %s", w.Code, w.Body)
	}
	s, _ := f.db.HandoffSetting(ctx)
	if ok, err := f.db.BindStaffChat(ctx, s.BindCode, -100); !ok || err != nil {
		t.Fatalf("綁員工群組：%v %v", ok, err)
	}
	intent := func(c string) map[string]model.Value { return map[string]model.Value{"intent": model.Text(c)} }
	bot := &dialog.Bot{DB: f.db, Model: model.NewScripted(map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("identity")}}), Now: f.srv.Now}
	f.srv.Handoff = bot
	say := func(text string) dialog.Reply {
		r, err := bot.Receive(ctx, dialog.Message{Platform: "telegram", UserID: "u1", ChatID: "c-u1", Text: text})
		if err != nil {
			t.Fatal(err)
		}
		return r
	}
	say("我要找真人")
	f.appendLog("u1", "customer", "我要找真人", "", f.now) // 對話紀錄由平台層寫
	var c struct {
		Human    bool
		Messages []struct{ Direction, Text string }
	}
	f.getJSON("/api/conversations/telegram/u1", cookie, &c)
	if !c.Human {
		t.Fatalf("轉人工中：%+v", c)
	}
	if w := f.do("POST", path, "{}", cookie); w.Code != http.StatusOK {
		t.Fatalf("結束轉人工：%d %s", w.Code, w.Body)
	}
	c.Human = false
	f.getJSON("/api/conversations/telegram/u1", cookie, &c)
	if l := c.Messages[len(c.Messages)-1]; c.Human || l.Direction != "staff" || l.Text != "（結束轉人工）" {
		t.Errorf("結束後：%+v", c)
	}
	if open, _ := f.db.OpenTickets(ctx); len(open) != 0 {
		t.Errorf("轉人工事項沒結束：%+v", open)
	}
	if r := say("你是誰"); r.Silent || r.Key != "identity" {
		t.Errorf("AI 沒恢復回覆：%+v", r)
	}
	if w := f.do("POST", path, "{}", cookie); w.Code != http.StatusConflict {
		t.Errorf("不在轉人工中：%d %s", w.Code, w.Body)
	}
}

// 一位顧客的對話：舊的在前，每則附編號和 AI 判斷；不是商家慣用語言的顧客、機器人訊息才能翻譯（看不出語言時用顧客的語言）。
// 沒有對話紀錄回 404
func TestConversationDetail(t *testing.T) {
	f, cookie := newPendingFixture(t)
	if err := f.db.SaveCustomer(context.Background(), "telegram", "u1", []byte(`{"lang":"en"}`), f.now); err != nil {
		t.Fatal(err)
	}
	f.appendLog("u1", "customer", "Do you have parking?", "", f.now)
	f.appendLog("u1", "bot", "Yes, across the street.", `[{"step":"intent","result":"question"},{"step":"verify","result":"number_mismatch"}]`, f.now)
	f.appendLog("u1", "customer", "👍", "", f.now)
	f.appendLog("u1", "staff", "有的喔", "", f.now)
	f.appendLog("u1", "customer", "謝謝", "", f.now)
	f.appendLog("u2", "customer", "別人", "", f.now)

	var c conversationView
	f.getJSON("/api/conversations/telegram/u1", cookie, &c)
	if c.Lang != "en" || c.Preferred != "zh" || len(c.Messages) != 5 {
		t.Fatalf("對話：%+v", c)
	}
	m := c.Messages
	if m[0].Text != "Do you have parking?" || m[0].ID == 0 || m[1].ID <= m[0].ID {
		t.Errorf("順序、編號：%+v", m)
	}
	if tr := m[1].Traces; len(tr) != 2 || tr[0].Step != "intent" || tr[0].Result != "question" || tr[1].Result != "number_mismatch" {
		t.Errorf("AI 判斷：%+v", tr)
	}
	if m[0].Traces == nil {
		t.Error("沒有 AI 判斷時要是空陣列")
	}
	for i, want := range []bool{true, true, true, false, false} {
		if m[i].Translatable != want {
			t.Errorf("%q 能不能翻譯：got %v", m[i].Text, m[i].Translatable)
		}
	}
	if w := f.do("GET", "/api/conversations/telegram/nobody", "", cookie); w.Code != http.StatusNotFound {
		t.Errorf("沒有紀錄：%d", w.Code)
	}
}

// 送不出去的回覆附 send_failed，後台標成「送出失敗」
func TestConversationSendFailed(t *testing.T) {
	f, cookie := newPendingFixture(t)
	f.appendLog("u1", "customer", "有停車場嗎", "", f.now)
	if err := f.db.AppendLog(context.Background(), store.LogEntry{Time: f.now, Platform: "telegram", UserID: "u1", ChatID: "c-u1", Direction: "bot", Text: "對面有", SendFailed: true}); err != nil {
		t.Fatal(err)
	}
	var c conversationView
	f.getJSON("/api/conversations/telegram/u1", cookie, &c)
	if len(c.Messages) != 2 || c.Messages[0].SendFailed || !c.Messages[1].SendFailed {
		t.Errorf("送出失敗：%+v", c.Messages)
	}
}

// 翻譯一則對話：從訊息的語言翻成商家慣用語言（不存起來）；顧客的訊息用翻譯顧客訊息的提示詞，機器人的回覆用翻譯店家回覆的。
// 已經是慣用語言回 400、模型沒翻出來回 502、沒有模型回 503、沒有這則回 404
func TestTranslateLog(t *testing.T) {
	f, cookie := newPendingFixture(t)
	f.appendLog("u1", "customer", "Do you have parking?", "", f.now)
	f.appendLog("u1", "customer", "有停車場嗎", "", f.now)
	f.appendLog("u1", "bot", "Yes, across the street.", "", f.now)
	var c conversationView
	f.getJSON("/api/conversations/telegram/u1", cookie, &c)
	en, zh, bot := c.Messages[0].ID, c.Messages[1].ID, c.Messages[2].ID
	path := func(id int64) string { return fmt.Sprintf("/api/conversations/logs/%d/translate", id) }

	if w := f.do("POST", path(en), "{}", cookie); w.Code != http.StatusServiceUnavailable {
		t.Errorf("沒有模型：%d %s", w.Code, w.Body)
	}
	chat := model.NewScripted(map[string][]map[string]model.Value{
		"TranslateCustomer": {{"translation": model.Text("有停車場嗎？")}},
		"Translate":         {{"translation": model.Text("有的，在對面。")}},
	})
	f.srv.Chat = chat
	w := f.do("POST", path(en), "{}", cookie)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"translation":"有停車場嗎？"`) {
		t.Errorf("翻譯：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", path(bot), "{}", cookie); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"translation":"有的，在對面。"`) {
		t.Errorf("翻譯機器人的回覆：%d %s", w.Code, w.Body)
	}
	if reqs := chat.Requests(); len(reqs) != 2 || reqs[0].Name != "TranslateCustomer" || reqs[0].Text != "Do you have parking?" ||
		!strings.Contains(reqs[0].Instructions, "從English翻成繁體中文") || reqs[1].Name != "Translate" || reqs[1].Text != "Yes, across the street." {
		t.Errorf("翻譯的請求：%+v", reqs)
	}
	if w := f.do("POST", path(en), "{}", cookie); w.Code != http.StatusBadGateway {
		t.Errorf("模型沒翻出來：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", path(zh), "{}", cookie); w.Code != http.StatusBadRequest {
		t.Errorf("已經是慣用語言：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", path(999), "{}", cookie); w.Code != http.StatusNotFound {
		t.Errorf("沒有這則：%d %s", w.Code, w.Body)
	}
}

// 本月（台北時間）的 AI 費用和呼叫次數，依模型分開、附失敗次數；附目前的模型和每月上限。沒有紀錄時都是 0
func TestUsage(t *testing.T) {
	f := newFixture(t)
	f.withSwitch(nil)
	cookie := f.login()
	want := `{"from":"2026-10-01T00:00:00+08:00","to":"2026-11-01T00:00:00+08:00","total_cost_usd":0,"total_calls":0,"by_model":[],"model":"afm","monthly_budget_usd":10}`
	if w := f.do("GET", "/api/usage", "", cookie); w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), want) {
		t.Errorf("沒有紀錄：%d %s", w.Code, w.Body)
	}
	octStart := time.Date(2026, 9, 30, 16, 0, 0, 0, time.UTC) // 台北 10/1 00:00
	novStart := time.Date(2026, 10, 31, 16, 0, 0, 0, time.UTC)
	for _, u := range []store.UsageRecord{
		{Time: octStart.Add(-time.Second), Model: "claude-haiku-5-5", Request: "Intent", InputTokens: 1, OutputTokens: 1, CostUSD: 1, OK: true},
		{Time: octStart, Model: "claude-haiku-5-5", Request: "Intent", InputTokens: 100, OutputTokens: 10, CostUSD: 0.25, OK: true},
		{Time: octStart.Add(time.Hour), Model: "claude-haiku-5-5", Request: "Pick", InputTokens: 50, OutputTokens: 300, CostUSD: 0.5, OK: false},
		{Time: octStart.Add(2 * time.Hour), Model: "afm", Request: "Intent", OK: true},
		{Time: octStart.Add(3 * time.Hour), Model: "afm", Request: "Answer", OK: false},
		{Time: octStart.Add(4 * time.Hour), Model: "afm", Request: "Answer", OK: true},
		{Time: novStart, Model: "claude-haiku-5-5", Request: "Intent", CostUSD: 4, OK: true},
	} {
		if err := f.db.RecordUsage(context.Background(), u); err != nil {
			t.Fatal(err)
		}
	}
	want = `{"from":"2026-10-01T00:00:00+08:00","to":"2026-11-01T00:00:00+08:00","total_cost_usd":0.75,"total_calls":5,"by_model":[
		{"model":"claude-haiku-5-5","calls":2,"failed":1,"input_tokens":150,"output_tokens":310,"cost_usd":0.75},
		{"model":"afm","calls":3,"failed":1,"input_tokens":0,"output_tokens":0,"cost_usd":0}
	],"model":"afm","monthly_budget_usd":10}`
	if w := f.do("GET", "/api/usage", "", cookie); w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), want) {
		t.Errorf("本月：%d %s", w.Code, w.Body)
	}
}

type bookingView struct {
	ID        int64
	Date      string
	Name      string
	Status    string
	HandledBy string    `json:"handled_by"`
	CreatedAt time.Time `json:"created_at"`
}

// addBooking 開一張預約申請（和對話層一樣先開待處理事項），status 不是 pending 時再確認或拒絕
func (f *fixture) addBooking(date, name, status string, at time.Time) int64 {
	f.t.Helper()
	ctx := context.Background()
	id, err := f.db.AddTicket(ctx, store.Ticket{Kind: "booking", Platform: "telegram", UserID: "u-" + name, ChatID: "c-" + name, Text: "📅"}, at)
	if err != nil {
		f.t.Fatal(err)
	}
	if err := f.db.AddBooking(ctx, store.Booking{ID: id, Platform: "telegram", UserID: "u-" + name, ChatID: "c-" + name, Date: date, Name: name, Summary: name + " " + date}, at); err != nil {
		f.t.Fatal(err)
	}
	if status != "pending" {
		if err := f.db.SetBookingStatus(ctx, id, status, "小明", at); err != nil {
			f.t.Fatal(err)
		}
	}
	return id
}

// 預約列表：新申請的在前；可以用狀態和預約日期（含頭尾）篩選；狀態、日期格式不對回 400
func TestBookings(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	var empty []bookingView
	f.getJSON("/api/bookings", cookie, &empty)
	if empty == nil || len(empty) != 0 {
		t.Errorf("沒有預約：%+v", empty)
	}
	f.addBooking("2026-10-20", "王", "confirmed", f.now.Add(-3*time.Hour))
	f.addBooking("2026-10-21", "李", "pending", f.now.Add(-2*time.Hour))
	f.addBooking("2026-10-25", "陳", "rejected", f.now.Add(-time.Hour))

	for _, c := range []struct{ query, want string }{
		{"", "陳,李,王"},
		{"?status=pending", "李"},
		{"?status=confirmed", "王"},
		{"?from=2026-10-21", "陳,李"},
		{"?from=2026-10-20&to=2026-10-21", "李,王"},
	} {
		var list []bookingView
		f.getJSON("/api/bookings"+c.query, cookie, &list)
		names := []string{}
		for _, b := range list {
			names = append(names, b.Name)
		}
		if got := strings.Join(names, ","); got != c.want {
			t.Errorf("%q：got %s, want %s", c.query, got, c.want)
		}
	}
	var list []bookingView
	f.getJSON("/api/bookings?status=confirmed", cookie, &list)
	if b := list[0]; b.HandledBy != "小明" || b.Date != "2026-10-20" || !b.CreatedAt.Equal(f.now.Add(-3*time.Hour)) {
		t.Errorf("預約內容：%+v", b)
	}
	for _, q := range []string{"?status=done", "?from=10/20", "?to=2026-13-01"} {
		if w := f.do("GET", "/api/bookings"+q, "", cookie); w.Code != http.StatusBadRequest {
			t.Errorf("%s：%d", q, w.Code)
		}
	}
}
