package telegram

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

const (
	testToken  = [REDACTED]"
	testSecret = [REDACTED]"
)

// fakeAPI 是假的 Telegram Bot API：記下收到的 method 和參數；fail 回傳 true 時那次回 500
type fakeAPI struct {
	mu    sync.Mutex
	calls []call
	fail  func(method string, n int) bool // n：這個 method 第幾次（從 1 開始）
	srv   *httptest.Server
}

type call struct {
	Method string
	Params map[string]any
}

func newFakeAPI(t *testing.T) *fakeAPI {
	f := &fakeAPI{}
	f.srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prefix := "/bot" + testToken + "/"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			http.NotFound(w, r)
			return
		}
		method := strings.TrimPrefix(r.URL.Path, prefix)
		var params map[string]any
		json.NewDecoder(r.Body).Decode(&params)
		f.mu.Lock()
		f.calls = append(f.calls, call{method, params})
		n := 0
		for _, c := range f.calls {
			if c.Method == method {
				n++
			}
		}
		fail := f.fail != nil && f.fail(method, n)
		seq := len(f.calls)
		if !fail {
			params["_message_id"] = float64(1000 * seq) // 測試用：送出的訊息編號
		}
		f.mu.Unlock()
		if fail {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"ok":false,"description":"Internal Server Error"}`))
			return
		}
		// 送出訊息的 method 回訊息編號：第 n 次呼叫是 1000*n（看得出是哪一次送的）
		if method == "sendMessage" || method == "copyMessage" {
			w.Write([]byte(`{"ok":true,"result":{"message_id":` + itoa(int64(1000*seq)) + `}}`))
			return
		}
		w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeAPI) methods() []string {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []string
	for _, c := range f.calls {
		out = append(out, c.Method)
	}
	return out
}

func newClient(base string, log *bytes.Buffer) *Client {
	c := NewClient(testToken, base, NewLogger(&syncBuffer{b: log}, testToken, testSecret))
	c.RetryDelays = []time.Duration{0, 0, 0}
	return c
}

// syncBuffer：背景處理和測試同時讀寫日誌
type syncBuffer struct {
	mu sync.Mutex
	b  *bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

type env struct {
	t          *testing.T
	nextUpdate int64
	api        *fakeAPI
	db         *store.DB
	hook       *Webhook
	log        *bytes.Buffer
	process    func(context.Context, Update) // nil：交給 Bot
	dialog     *dialog.Bot                   // 對話層：預設的假模型每句都分成一般對話
	bot        *Bot
}

func newEnv(t *testing.T) *env {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	e := &env{t: t, nextUpdate: 1000, api: newFakeAPI(t), db: db, log: &bytes.Buffer{}}
	client := newClient(e.api.srv.URL, e.log)
	general := map[string]model.Value{"intent": model.Text("general")}
	e.dialog = &dialog.Bot{DB: db, Model: model.NewScripted(map[string][]map[string]model.Value{"Intent": {general, general, general}}), Log: client.Log}
	bot := &Bot{API: client, Store: db, Log: client.Log, Dialog: e.dialog}
	e.bot = bot
	e.hook = &Webhook{Secret: [REDACTED], Dedup: db, Log: client.Log, Process: func(ctx context.Context, u Update) {
		if e.process != nil {
			e.process(ctx, u)
			return
		}
		bot.Process(ctx, u)
	}}
	return e
}

func (e *env) post(secret string, body string) int {
	req := httptest.NewRequest(http.MethodPost, WebhookPath, strings.NewReader(body))
	if secret != "" {
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
	}
	rec := httptest.NewRecorder()
	e.hook.ServeHTTP(rec, req)
	return rec.Code
}

func privateText(updateID int64, text string) string {
	return `{"update_id":` + itoa(updateID) + `,"message":{"message_id":1,"chat":{"id":555,"type":"private"},"from":{"id":555},"text":"` + text + `"}}`
}

func itoa(n int64) string { b, _ := json.Marshal(n); return string(b) }

func TestWrongSecretIs401(t *testing.T) {
	e := newEnv(t)
	processed := 0
	e.process = func(context.Context, Update) { processed++ }
	for _, secret := range []string{"", "wrong", testSecret + "x"} {
		if code := e.post(secret, privateText(1, "hi")); code != http.StatusUnauthorized {
			t.Errorf("secret %q：got %d, want 401", secret, code)
		}
	}
	e.hook.Wait()
	if processed != 0 {
		t.Errorf("secret 不符卻處理了 %d 次", processed)
	}
	// 被擋下的 update 沒有記成已處理，之後帶對的 secret 還能處理
	if code := e.post(testSecret, privateText(1, "hi")); code != http.StatusOK {
		t.Errorf("got %d, want 200", code)
	}
	e.hook.Wait()
	if processed != 1 {
		t.Errorf("處理了 %d 次，應該 1 次", processed)
	}
}

func TestDuplicateUpdateProcessedOnce(t *testing.T) {
	e := newEnv(t)
	var mu sync.Mutex
	processed := 0
	e.process = func(context.Context, Update) { mu.Lock(); processed++; mu.Unlock() }
	for i := 0; i < 3; i++ {
		if code := e.post(testSecret, privateText(10, "hi")); code != http.StatusOK {
			t.Errorf("第 %d 次 got %d, want 200", i+1, code)
		}
	}
	e.post(testSecret, privateText(11, "hi"))
	e.hook.Wait()
	if processed != 2 {
		t.Errorf("處理了 %d 次，應該 2 次（update 10、11 各一次）", processed)
	}
}

func TestRespondsBeforeProcessing(t *testing.T) {
	e := newEnv(t)
	release := make(chan struct{})
	done := make(chan struct{})
	e.process = func(context.Context, Update) { <-release; close(done) }

	got := make(chan int, 1)
	go func() { got <- e.post(testSecret, privateText(20, "hi")) }()
	select {
	case code := <-got:
		if code != http.StatusOK {
			t.Errorf("got %d, want 200", code)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("背景處理還卡著，webhook 就該先回 200")
	}
	select {
	case <-done:
		t.Fatal("還沒放行就處理完了")
	default:
	}
	close(release)
	e.hook.Wait()
	<-done
}

func textFrom(updateID, userID int64, text string) string {
	return `{"update_id":` + itoa(updateID) + `,"message":{"message_id":1,"chat":{"id":` + itoa(userID) + `,"type":"private"},"from":{"id":` + itoa(userID) + `},"text":"` + text + `"}}`
}

// 依顧客加鎖：同一位顧客連傳兩則，第一則處理完才處理第二則（順序和收到的一樣）；不同顧客不用等
func TestSameCustomerProcessedInOrder(t *testing.T) {
	e := newEnv(t)
	release := make(chan struct{})
	started := make(chan string, 3)
	var mu sync.Mutex
	var order []string
	e.process = func(_ context.Context, u Update) {
		started <- u.Message.Text
		if u.Message.Text == "第一則" {
			<-release
		}
		mu.Lock()
		order = append(order, u.Message.Text)
		mu.Unlock()
	}
	e.post(testSecret, textFrom(60, 555, "第一則"))
	if got := <-started; got != "第一則" {
		t.Fatalf("先開始的是 %s", got)
	}
	e.post(testSecret, textFrom(61, 555, "第二則"))
	e.post(testSecret, textFrom(62, 777, "別的顧客"))
	select {
	case got := <-started:
		if got != "別的顧客" {
			t.Fatalf("第一則還沒處理完就開始了 %s", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("別的顧客被同一把鎖擋住了")
	}
	select {
	case got := <-started:
		t.Fatalf("第一則還沒處理完就開始了 %s", got)
	case <-time.After(50 * time.Millisecond):
	}
	close(release)
	e.hook.Wait()
	mu.Lock()
	defer mu.Unlock()
	if got := strings.Join(order, ","); got != "別的顧客,第一則,第二則" {
		t.Errorf("處理完成的順序：%s", got)
	}
}

func TestUnparsableBodyIs200(t *testing.T) {
	e := newEnv(t)
	processed := 0
	e.process = func(context.Context, Update) { processed++ }
	for _, body := range []string{"not json", "{}"} {
		if code := e.post(testSecret, body); code != http.StatusOK {
			t.Errorf("%q：got %d, want 200", body, code)
		}
	}
	req := httptest.NewRequest(http.MethodGet, WebhookPath, nil)
	rec := httptest.NewRecorder()
	e.hook.ServeHTTP(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET：got %d, want 405", rec.Code)
	}
	e.hook.Wait()
	if processed != 0 {
		t.Errorf("讀不懂的內容不該處理")
	}
}

// 私訊交給對話層：回對話層的回覆，對話紀錄記下回覆種類和 AI 判斷
func TestPrivateMessageRepliedByDialog(t *testing.T) {
	e := newEnv(t)
	e.post(testSecret, privateText(30, "好的"))
	e.hook.Wait()
	if got := strings.Join(e.api.methods(), ","); got != "sendChatAction,sendMessage" {
		t.Fatalf("呼叫的 API：%s", got)
	}
	e.api.mu.Lock()
	action, msg := e.api.calls[0].Params, e.api.calls[1].Params
	e.api.mu.Unlock()
	if action["action"] != "typing" || action["chat_id"] != float64(555) {
		t.Errorf("sendChatAction 參數：%v", action)
	}
	if msg["text"] != "收到～" || msg["chat_id"] != float64(555) {
		t.Errorf("sendMessage 參數：%v", msg)
	}
	logs, err := e.db.Logs(context.Background(), "telegram", "555")
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[0].Direction != "customer" || logs[0].Text != "好的" || logs[1].Direction != "bot" || logs[1].Text != "收到～" || logs[1].ReplyKey != "general" {
		t.Errorf("對話紀錄：%+v", logs)
	}
	if string(logs[1].Traces) != `[{"step":"intent","result":"general"},{"step":"chat","result":"general"}]` || logs[1].SendFailed {
		t.Errorf("AI 判斷：%s、送出失敗：%v", logs[1].Traces, logs[1].SendFailed)
	}
}

// bot 指令（/start、/start@bot、/help）不送進分類：回打招呼、用商家慣用語言（英文字母的指令不改顧客語言），不呼叫模型（不會挑 FAQ）
func TestBotCommandGreets(t *testing.T) {
	for _, cmd := range []string{"/start", "/start@sunchabot", "/help"} {
		e := newEnv(t)
		m := model.NewScripted(nil)
		e.dialog.Model = m
		e.post(testSecret, privateText(40, cmd))
		e.hook.Wait()
		if len(m.Requests()) != 0 {
			t.Errorf("%s：呼叫了模型 %+v", cmd, m.Requests())
		}
		logs, _ := e.db.Logs(context.Background(), "telegram", "555")
		if len(logs) != 2 || logs[0].Text != cmd || logs[1].Text != "您好～有什麼可以幫您的嗎？" || logs[1].ReplyKey != "greeting" {
			t.Errorf("%s 對話紀錄：%+v", cmd, logs)
		}
		if c, _ := e.dialog.LoadCustomer(context.Background(), "telegram", "555"); c.Lang != "zh" {
			t.Errorf("%s 顧客語言：%q", cmd, c.Lang)
		}
	}
}

// 回覆送不出去（重試後放棄）：對話紀錄照樣記下原本要回的話和 AI 判斷，標成送出失敗
func TestFailedReplyStillLogged(t *testing.T) {
	e := newEnv(t)
	e.api.fail = func(method string, _ int) bool { return method == "sendMessage" }
	e.post(testSecret, privateText(33, "好的"))
	e.hook.Wait()
	logs, err := e.db.Logs(context.Background(), "telegram", "555")
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[1].Direction != "bot" || logs[1].Text != "收到～" || logs[1].ReplyKey != "general" || !logs[1].SendFailed {
		t.Fatalf("對話紀錄：%+v", logs)
	}
	if string(logs[1].Traces) != `[{"step":"intent","result":"general"},{"step":"chat","result":"general"}]` {
		t.Errorf("AI 判斷：%s", logs[1].Traces)
	}
	if !strings.Contains(e.log.String(), "回覆送出失敗") {
		t.Errorf("應該記錄送出失敗：\n%s", e.log.String())
	}
}

// 轉人工中對話層不回覆：不送訊息（不顯示正在輸入），紀錄照寫一筆空白的
func TestSilentWhileHandedOff(t *testing.T) {
	e := newEnv(t)
	e.bind() // 沒綁員工群組不會轉人工（#55）
	e.dialog.Model = model.NewScripted(map[string][]map[string]model.Value{"Intent": {{"intent": model.Text("human_request")}}})
	e.post(testSecret, privateText(31, "我要找真人"))
	e.hook.Wait()
	e.post(testSecret, privateText(32, "在嗎"))
	e.hook.Wait()
	// 回顧客、通知員工群組各一則，之後的「在嗎」不送
	if got := strings.Join(e.api.methods(), ","); got != "sendChatAction,sendMessage,sendMessage" {
		t.Fatalf("呼叫的 API：%s", got)
	}
	logs, _ := e.db.Logs(context.Background(), "telegram", "555")
	if len(logs) != 4 || logs[3].Direction != "bot" || logs[3].Text != "" || logs[3].ReplyKey != "silent" {
		t.Errorf("對話紀錄：%+v", logs)
	}
}

// 轉人工中的顧客傳文字、bot 指令、貼圖、圖片：都不呼叫 AI，也不回顧客、不顯示正在輸入（#57）
func TestHandedOffCustomerMakesNoAICalls(t *testing.T) {
	e := newEnv(t)
	e.bind()
	m := model.NewScripted(map[string][]map[string]model.Value{"Intent": {{"intent": model.Text("general")}}})
	e.dialog.Model = m
	data, _ := json.Marshal(dialog.Customer{Human: true, HumanSince: time.Now()})
	e.db.SaveCustomer(context.Background(), "telegram", "555", data, time.Now())
	sticker := privateMsg(555, "")
	delete(sticker, "text")
	sticker["sticker"] = map[string]any{"emoji": "👍"}
	for _, u := range []map[string]any{{"message": privateMsg(555, "請問店員什麼時候會回我")}, {"message": privateMsg(555, "/start")}, {"message": sticker}, photo(555, "這個有貨嗎")} {
		e.send(u)
	}
	for _, req := range m.Requests() {
		t.Errorf("轉人工中呼叫了 AI：%s %q", req.Name, req.Text)
	}
	if c := e.toChat(555); len(c) != 0 || len(e.api.params("sendChatAction")) != 0 {
		t.Errorf("轉人工中回了顧客：%v %v", c, e.api.methods())
	}
}

// 還沒綁定員工群組時：群組訊息、頻道貼文、私訊裡看不懂的訊息（位置）不理會；按鈕回呼只回應 Telegram
func TestGroupAndUnknownIgnored(t *testing.T) {
	e := newEnv(t)
	e.post(testSecret, `{"update_id":40,"message":{"message_id":1,"chat":{"id":-100,"type":"group"},"from":{"id":1},"text":"hi"}}`)
	e.post(testSecret, `{"update_id":41,"message":{"message_id":2,"chat":{"id":555,"type":"private"},"from":{"id":555},"location":{"latitude":25,"longitude":121}}}`)
	e.post(testSecret, `{"update_id":43,"message":{"message_id":3,"chat":{"id":-200,"type":"channel"},"text":"hi"}}`)
	e.post(testSecret, `{"update_id":42,"callback_query":{"id":"x"}}`)
	e.hook.Wait()
	if m := strings.Join(e.api.methods(), ","); m != "answerCallbackQuery" {
		t.Errorf("呼叫的 API：%v", m)
	}
}

func TestRetryUntilSuccess(t *testing.T) {
	api := newFakeAPI(t)
	api.fail = func(method string, n int) bool { return n <= 2 }
	var log bytes.Buffer
	c := newClient(api.srv.URL, &log)
	if err := c.SendMessage(context.Background(), 1, "hi"); err != nil {
		t.Fatalf("第三次應該成功：%v", err)
	}
	if n := len(api.methods()); n != 3 {
		t.Errorf("送了 %d 次，應該 3 次", n)
	}
	if strings.Count(log.String(), "稍後重試") != 2 {
		t.Errorf("應該記錄 2 次重試：\n%s", log.String())
	}
}

func TestRetryGivesUp(t *testing.T) {
	api := newFakeAPI(t)
	api.fail = func(string, int) bool { return true }
	var log bytes.Buffer
	c := newClient(api.srv.URL, &log)
	if err := c.SendMessage(context.Background(), 1, "hi"); err == nil {
		t.Fatal("一直失敗應該回傳錯誤")
	}
	if n := len(api.methods()); n != 4 {
		t.Errorf("送了 %d 次，應該 1+3 次", n)
	}
	if !strings.Contains(log.String(), "放棄") {
		t.Errorf("應該記錄放棄：\n%s", log.String())
	}
}

func TestChatActionNotRetried(t *testing.T) {
	api := newFakeAPI(t)
	api.fail = func(string, int) bool { return true }
	var log bytes.Buffer
	c := newClient(api.srv.URL, &log)
	if err := c.SendChatAction(context.Background(), 1); err == nil {
		t.Fatal("失敗應該回傳錯誤")
	}
	if n := len(api.methods()); n != 1 {
		t.Errorf("「正在輸入」送了 %d 次，不該重試", n)
	}
}

// 連線錯誤（*url.Error 帶完整網址）和 HTTP 錯誤時，日誌和錯誤訊息都不能有 token
func TestTokenNotInLogs(t *testing.T) {
	ln, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		t.Fatal(err)
	}
	closed := "http://" + ln.Addr().String()
	ln.Close()

	api := newFakeAPI(t)
	api.fail = func(string, int) bool { return true }
	for _, base := range []string{closed, api.srv.URL} {
		var log bytes.Buffer
		c := newClient(base, &log)
		// 日誌不經過遮罩，直接檢查 Client 本身有沒有把 token 寫出去
		c.Log = NewLogger(&log)
		err := c.SendMessage(context.Background(), 1, "hi")
		if err == nil {
			t.Fatalf("%s：應該失敗", base)
		}
		if strings.Contains(err.Error(), testToken) {
			t.Errorf("錯誤訊息含 token：%v", err)
		}
		if strings.Contains(log.String(), testToken) || log.Len() == 0 {
			t.Errorf("日誌含 token 或沒有記錄：\n%s", log.String())
		}
	}
}

func TestLoggerRedactsSecrets(t *testing.T) {
	var buf bytes.Buffer
	NewLogger(&buf, testToken, testSecret, "").Info("x", "a", "token=[REDACTED], "b", testSecret)
	if s := buf.String(); strings.Contains(s, testToken) || strings.Contains(s, testSecret) || !strings.Contains(s, "***") {
		t.Errorf("沒遮掉：%s", s)
	}
}

func TestConfigFromEnv(t *testing.T) {
	env := map[string]string{
		"TELEGRAM_BOT_TOKEN":      testToken,
		"TELEGRAM_WEBHOOK_SECRET": testSecret,
		"TELEGRAM_PUBLIC_URL":     "https://example.invalid/",
		"PORT":                    "9000",
	}
	c, err := ConfigFromEnv(func(k string) string { return env[k] })
	if err != nil {
		t.Fatal(err)
	}
	if c.WebhookURL() != "https://example.invalid/telegram/webhook" || c.Port != 9000 || c.APIBase != "https://example.invalid" {
		t.Errorf("%+v", c)
	}
	for k, v := range map[string]string{"TELEGRAM_WEBHOOK_SECRET": "有空白 x", "TELEGRAM_PUBLIC_URL": "http://x", "PORT": "abc", "TELEGRAM_BOT_TOKEN": ""} {
		bad := map[string]string{}
		for k2, v2 := range env {
			bad[k2] = v2
		}
		bad[k] = v
		_, err := ConfigFromEnv(func(k string) string { return bad[k] })
		if err == nil || !strings.Contains(err.Error(), k) {
			t.Errorf("%s=%q 應該報錯並指出變數：%v", k, v, err)
		}
		if err != nil && strings.Contains(err.Error(), testToken) {
			t.Errorf("錯誤訊息含 token：%v", err)
		}
	}
}
