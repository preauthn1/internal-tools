package main

import (
	"aichatbot/internal/admin"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
	"aichatbot/internal/telegram"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"aichatbot/internal/dialog"
)

const (
	token   = [REDACTED]"
	secret  = [REDACTED]"
	evalDir = "../../../eval"
)

// 整支伺服器：setWebhook、私訊回固定的話；重開（同一個資料庫）後同一個 update 不會再處理，對話紀錄還在
func TestServeAndRestart(t *testing.T) {
	var mu sync.Mutex
	var methods []string
	var webhookParams map[string]any
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method := strings.TrimPrefix(r.URL.Path, "/bot"+token+"/")
		var params map[string]any
		json.NewDecoder(r.Body).Decode(&params)
		mu.Lock()
		methods = append(methods, method)
		if method == "setWebhook" {
			webhookParams = params
		}
		mu.Unlock()
		w.Write([]byte(`{"ok":true,"result":true}`))
	}))
	defer api.Close()

	cfg := telegram.Config{Token: [REDACTED], Secret: [REDACTED], PublicURL: "https://example.invalid", APIBase: api.URL, Host: "0.0.0.0", Port: 0}
	dbPath := filepath.Join(t.TempDir(), "bot.sqlite")
	// 模型層：預設 AFM，開 AI 回答快取（結束時要印統計）
	mcfg, err := model.ConfigFromEnv(func(k string) string {
		return map[string]string{"BOT_AI_CACHE": "1", "BOT_AI_CACHE_DIR": t.TempDir()}[k]
	})
	if err != nil {
		t.Fatal(err)
	}
	mcfg.AppleURL = api.URL // 不碰開發機上可能開著的模型服務
	loadEmbedder = func(*slog.Logger) dialog.Embedder { return nil }
	pages := t.TempDir()
	os.WriteFile(filepath.Join(pages, "index.html"), []byte("<html></html>"), 0o644)
	acfg := admin.Config{Host: "0.0.0.0", Port: 0, Dir: pages}
	update := `{"update_id":7,"message":{"message_id":1,"chat":{"id":9,"type":"private"},"from":{"id":9},"text":"你們幾點開"}}`
	var logs bytes.Buffer
	var logMu sync.Mutex

	start := func() {
		ctx, cancel := context.WithCancel(context.Background())
		addr := make(chan string, 1)
		done := make(chan error, 1)
		log := telegram.NewLogger(writerFunc(func(p []byte) (int, error) {
			logMu.Lock()
			defer logMu.Unlock()
			return logs.Write(p)
		}), token, secret)
		adminAddr := make(chan string, 1)
		go func() {
			done <- run(ctx, cfg, mcfg, acfg, dbPath, log, func(a, adm string) { addr <- a; adminAddr <- adm })
		}()
		var base string
		select {
		case base = <-addr:
		case err := <-done:
			t.Fatal(err)
		}
		// 後台另一個 port：頁面不用登入，API 要先登入；webhook 的 port 上沒有後台
		adm := <-adminAddr
		if adm == base {
			t.Errorf("後台和 webhook 用同一個位址 %s", adm)
		}
		for _, c := range []struct {
			url  string
			want int
		}{{"http://" + adm + "/", http.StatusOK}, {"http://" + adm + "/api/auth/me", http.StatusUnauthorized}, {"http://" + base + "/api/auth/me", http.StatusNotFound}} {
			resp, err := http.Get(c.url)
			if err != nil {
				t.Fatal(err)
			}
			resp.Body.Close()
			if resp.StatusCode != c.want {
				t.Errorf("GET %s：%d, want %d", c.url, resp.StatusCode, c.want)
			}
		}
		req, _ := http.NewRequest(http.MethodPost, "http://"+base+telegram.WebhookPath, strings.NewReader(update))
		req.Header.Set("X-Telegram-Bot-Api-Secret-Token", secret)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("got %d, want 200", resp.StatusCode)
		}
		// 等 setWebhook 送出再停（不然停機時會被取消）
		deadline := time.Now().Add(5 * time.Second)
		for {
			mu.Lock()
			ok := webhookParams != nil
			mu.Unlock()
			if ok || time.Now().After(deadline) {
				break
			}
			time.Sleep(10 * time.Millisecond)
		}
		cancel()
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}

	start()
	start()

	mu.Lock()
	defer mu.Unlock()
	if webhookParams["url"] != "https://example.invalid/telegram/webhook" || webhookParams["secret_token"] != secret {
		t.Errorf("setWebhook 參數：%v", webhookParams)
	}
	var sends int
	for _, m := range methods {
		if m == "sendMessage" {
			sends++
		}
	}
	if sends != 1 {
		t.Errorf("重開後同一個 update 又處理了：%v", methods)
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	entries, err := db.Logs(context.Background(), "telegram", "9")
	if err != nil {
		t.Fatal(err)
	}
	// 模型服務是假的 Telegram API（回不出欄位）：AI 失敗，對話層當成看不懂
	if len(entries) != 2 || entries[1].ReplyKey != "unclear" || entries[1].Text == "" {
		t.Errorf("對話紀錄：%+v", entries)
	}
	logMu.Lock()
	defer logMu.Unlock()
	if strings.Contains(logs.String(), token) {
		t.Errorf("日誌含 token：\n%s", logs.String())
	}
	// 啟動時印出員工群組的綁定碼
	if s, _ := db.HandoffSetting(context.Background()); !strings.Contains(logs.String(), "/bind "+s.BindCode) {
		t.Errorf("沒印綁定碼：\n%s", logs.String())
	}
	// 第一次處理了那則私訊（分類呼叫一次模型），重開後同一個 update 不再處理
	if !strings.Contains(logs.String(), "AI 快取：命中 0 次、實際呼叫 1 次") || !strings.Contains(logs.String(), "AI 快取：命中 0 次、實際呼叫 0 次") {
		t.Errorf("每次結束都要印 AI 快取統計：\n%s", logs.String())
	}
}

// create-admin：stdin 不是終端機時讀一行當密碼；建好就能登入；帳號重複、密碼太短會失敗
func TestCreateAdmin(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "bot.sqlite")
	var out bytes.Buffer
	if err := createAdmin(dbPath, []string{"owner"}, strings.NewReader("correct horse\n"), &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "owner") || strings.Contains(out.String(), "correct horse") {
		t.Errorf("輸出：%q", out.String())
	}
	if err := createAdmin(dbPath, []string{"owner"}, strings.NewReader("another password\n"), &out); !errors.Is(err, admin.ErrUserExists) {
		t.Errorf("帳號重複：%v", err)
	}
	if err := createAdmin(dbPath, []string{"owner2"}, strings.NewReader("short\n"), &out); err == nil {
		t.Error("密碼太短應該失敗")
	}
	if err := createAdmin(dbPath, nil, strings.NewReader("correct horse\n"), &out); err == nil {
		t.Error("沒給帳號應該失敗")
	}

	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	srv := &admin.Server{DB: db, Dir: t.TempDir(), Now: time.Now, Log: slog.New(slog.DiscardHandler)}
	for _, c := range []struct {
		password string
		want     int
	}{{"correct horse", http.StatusOK}, {"another password", http.StatusUnauthorized}} {
		body, _ := json.Marshal(map[string]string{"username": "owner", "password": c.password})
		req := httptest.NewRequest("POST", "/api/auth/login", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		srv.Handler().ServeHTTP(w, req)
		if w.Code != c.want {
			t.Errorf("用 %q 登入：%d, want %d", c.password, w.Code, c.want)
		}
	}
}

// bind-code：印出綁定碼和目前綁定的員工群組
func TestBindCode(t *testing.T) {
	ctx := context.Background()
	dbPath := filepath.Join(t.TempDir(), "bot.sqlite")
	var out bytes.Buffer
	if err := bindCode(dbPath, &out); err != nil {
		t.Fatal(err)
	}
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	s, _ := db.HandoffSetting(ctx)
	if !strings.Contains(out.String(), "/bind "+s.BindCode) || !strings.Contains(out.String(), "還沒綁定") {
		t.Errorf("輸出：%s", out.String())
	}
	db.BindStaffChat(ctx, s.BindCode, -1001234)
	db.Close()
	out.Reset()
	if err := bindCode(dbPath, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "-1001234") {
		t.Errorf("輸出：%s", out.String())
	}
}

// import：重新匯入時印出保留的商家案例數，和因為和內建案例重複而刪掉的句子
func TestImportPrintsRemovedExamples(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "bot.sqlite")
	var out bytes.Buffer
	if err := importConfig(dbPath, []string{evalDir}, &out); err != nil {
		t.Fatal(err)
	}
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	tests, err := merchant.LoadTestSet(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"你們的貓叫什麼名字", "店裡可以帶狗嗎"} {
		if _, err := merchant.AddExample(ctx, db.SQL(), tests, text, "general"); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()

	c, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	c.Examples = append(c.Examples, merchant.Example{Text: "你們的貓叫什麼名字！", Intent: "off_topic", Builtin: true})
	data, err := merchant.Export(c)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "merchant.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := importConfig(dbPath, []string{path}, &out); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"商家新增的 1 句", "刪掉商家新增的 1 句", "- 你們的貓叫什麼名字（general）"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("輸出少了「%s」：\n%s", want, out.String())
		}
	}
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }

// 啟動時用後台存在資料庫的模型設定（環境變數只是初始值）；存的是 haiku 但現在沒有 API key 時記警告、用 AFM，不會啟動失敗
func TestServeUsesSavedModelSetting(t *testing.T) {
	ctx := context.Background()
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.Write([]byte(`{"ok":true,"result":true}`)) }))
	defer api.Close()
	dbPath := filepath.Join(t.TempDir(), "bot.sqlite")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.SaveModelSetting(ctx, model.ModelHaiku, 7); err != nil {
		t.Fatal(err)
	}
	if err := admin.CreateUser(ctx, db, "owner", "correct horse"); err != nil {
		t.Fatal(err)
	}
	db.Close()

	cfg := telegram.Config{Token: [REDACTED], Secret: [REDACTED], PublicURL: "https://example.invalid", APIBase: api.URL, Host: "0.0.0.0", Port: 0}
	mcfg, err := model.ConfigFromEnv(func(k string) string { return map[string]string{"BOT_AI_MONTHLY_BUDGET_USD": "10"}[k] })
	if err != nil {
		t.Fatal(err)
	}
	var logs bytes.Buffer
	var logMu sync.Mutex
	log := slog.New(slog.NewTextHandler(writerFunc(func(p []byte) (int, error) {
		logMu.Lock()
		defer logMu.Unlock()
		return logs.Write(p)
	}), nil))
	runCtx, cancel := context.WithCancel(ctx)
	adminAddr := make(chan string, 1)
	done := make(chan error, 1)
	go func() {
		done <- run(runCtx, cfg, mcfg, admin.Config{Host: "0.0.0.0", Port: 0, Dir: t.TempDir()}, dbPath, log, func(_, adm string) { adminAddr <- adm })
	}()
	var base string
	select {
	case base = <-adminAddr:
	case err := <-done:
		t.Fatal(err)
	}
	defer func() {
		cancel()
		if err := <-done; err != nil {
			t.Error(err)
		}
	}()

	login, _ := json.Marshal(map[string]string{"username": "owner", "password": "correct horse"})
	resp, err := http.Post("http://"+base+"/api/auth/login", "application/json", bytes.NewReader(login))
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	req, _ := http.NewRequest("GET", "http://"+base+"/api/model", nil)
	for _, c := range resp.Cookies() {
		req.AddCookie(c)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var got map[string]any
	json.NewDecoder(resp.Body).Decode(&got)
	if got["model"] != "afm" || got["monthly_budget_usd"] != 7.0 || got["haiku_available"] != false {
		t.Errorf("GET /api/model：%d %v", resp.StatusCode, got)
	}
	logMu.Lock()
	defer logMu.Unlock()
	if !strings.Contains(logs.String(), "沒有 API key") {
		t.Errorf("應該記警告：\n%s", logs.String())
	}
}
