// bot 是 Go 版的執行檔：
//
//	bot serve                 啟動 Telegram webhook 伺服器（環境變數見 README）
//	bot import <資料夾或檔案>  匯入商家設定：資料夾當 eval/ 格式，檔案當 export 輸出的 JSON；整份取代（商家新增的意圖案例保留）
//	bot export <檔案>          匯出商家設定成 JSON（- 是標準輸出）
//	bot create-admin <帳號>    建立後台管理員（密碼從終端機輸入，不顯示）
//	bot bind-code              印出員工群組的綁定碼和目前綁定的群組
//
// 每個子指令都可以加 -db <路徑>，預設是環境變數 BOT_DB，沒設就是工作目錄的 data/bot.sqlite
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"example.invalid/x/term"

	"aichatbot/internal/admin"
	"aichatbot/internal/dialog"
	"aichatbot/internal/embed"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
	"aichatbot/internal/telegram"
)

const usage = `用法：
  bot serve [-db 路徑]
  bot import [-db 路徑] <eval 資料夾或 JSON 檔>
  bot export [-db 路徑] <檔案|->
  bot create-admin [-db 路徑] <帳號>
  bot bind-code [-db 路徑]`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	dbPath := fs.String("db", defaultDBPath(), "SQLite 資料庫檔案")
	fs.Parse(args)

	var err error
	switch cmd {
	case "serve":
		err = serve(*dbPath)
	case "import":
		err = importConfig(*dbPath, fs.Args(), os.Stdout)
	case "export":
		err = exportConfig(*dbPath, fs.Args())
	case "create-admin":
		err = createAdmin(*dbPath, fs.Args(), os.Stdin, os.Stdout)
	case "bind-code":
		err = bindCode(*dbPath, os.Stdout)
	default:
		fmt.Fprintln(os.Stderr, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

func defaultDBPath() string {
	if p := os.Getenv("BOT_DB"); p != "" {
		return p
	}
	return "data/bot.sqlite"
}

func importConfig(dbPath string, args []string, out io.Writer) error {
	if len(args) != 1 {
		return errors.New(usage)
	}
	c, err := merchant.ReadFile(args[0])
	if err != nil {
		return err
	}
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx := context.Background()
	removed, err := merchant.Save(ctx, db.SQL(), c)
	if err != nil {
		return err
	}
	examples, err := merchant.ListExamples(ctx, db.SQL())
	if err != nil {
		return err
	}
	builtin := 0
	for _, e := range examples {
		if e.Builtin {
			builtin++
		}
	}
	fmt.Fprintf(out, "已匯入「%s」：知識段落 %d 段、FAQ %d 題、意圖案例 %d 句（內建 %d 句、商家新增的 %d 句）、固定回覆 %d 則\n",
		c.Merchant.Name, len(c.Knowledge), len(c.FAQ), len(examples), builtin, len(examples)-builtin, len(c.Replies))
	if len(removed) > 0 {
		fmt.Fprintf(out, "和內建案例重複，刪掉商家新增的 %d 句：\n", len(removed))
		for _, e := range removed {
			fmt.Fprintf(out, "- %s（%s）\n", e.Text, e.Intent)
		}
	}
	return nil
}

func exportConfig(dbPath string, args []string) error {
	if len(args) != 1 {
		return errors.New(usage)
	}
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	c, err := merchant.Load(context.Background(), db.SQL())
	if err != nil {
		return err
	}
	data, err := merchant.Export(c)
	if err != nil {
		return err
	}
	if args[0] == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}
	return os.WriteFile(args[0], data, 0o644)
}

// bindCode 印出員工群組的綁定碼和目前綁定的群組（後台還沒有這個畫面）
func bindCode(dbPath string, out io.Writer) error {
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	s, err := db.HandoffSetting(context.Background())
	if err != nil {
		return err
	}
	fmt.Fprintln(out, bindStatus(s))
	return nil
}

// bindStatus 是綁定碼的說明和目前綁定的群組（serve 啟動時也記在日誌）
func bindStatus(s store.HandoffSetting) string {
	status := "還沒綁定員工群組"
	if s.StaffChatID != 0 {
		status = fmt.Sprintf("目前綁定的員工群組：%d", s.StaffChatID)
	}
	return fmt.Sprintf("%s。把機器人加進員工群組，在群組裡傳「/bind %s」就會綁定那個群組（綁定後換新碼）", status, s.BindCode)
}

// createAdmin 建立後台管理員。in 是終端機時密碼不顯示、要輸入兩次；不是終端機時讀一行（給腳本用）
func createAdmin(dbPath string, args []string, in io.Reader, out io.Writer) error {
	if len(args) != 1 {
		return errors.New(usage)
	}
	password, err := readPassword(in)
	if err != nil {
		return err
	}
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := admin.CreateUser(context.Background(), db, args[0], password); err != nil {
		return err
	}
	fmt.Fprintf(out, "已建立管理員 %s\n", args[0])
	return nil
}

func readPassword(in io.Reader) (string, error) {
	if f, ok := in.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		fmt.Fprint(os.Stderr, "密碼：")
		p1, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		fmt.Fprint(os.Stderr, "再輸入一次：")
		p2, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(os.Stderr)
		if err != nil {
			return "", err
		}
		if string(p1) != string(p2) {
			return "", errors.New("兩次輸入的密碼不一樣")
		}
		return string(p1), nil
	}
	line, err := bufio.NewReader(in).ReadString('\n')
	if err != nil && !(errors.Is(err, io.EOF) && line != "") {
		return "", fmt.Errorf("讀不到密碼：%w", err)
	}
	return strings.TrimRight(line, "\r\n"), nil
}

// serve 從環境變數讀設定啟動伺服器，收到 SIGINT／SIGTERM 時停止
func serve(dbPath string) error {
	cfg, err := telegram.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	mcfg, err := model.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	acfg, err := admin.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return run(ctx, cfg, mcfg, acfg, dbPath, telegram.NewLogger(os.Stderr, cfg.Token, cfg.Secret, mcfg.APIKey), nil)
}

// run 啟動 webhook 伺服器和後台伺服器（兩個 port；開始接受連線後呼叫 ready(webhook 的位址, 後台的位址)），
// ctx 結束時停止收新請求、等背景處理完再回傳。
// 啟動後呼叫 setWebhook（失敗只記錄，伺服器繼續跑）；啟動時和之後每 24 小時清一次舊資料。
// 有開 AI 回答快取時，結束前印「AI 快取：命中 N 次、實際呼叫 M 次」
func run(ctx context.Context, cfg telegram.Config, mcfg model.Config, acfg admin.Config, dbPath string, log *slog.Logger, ready func(addr, adminAddr string)) error {
	db, err := store.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	chat, cacheStats, err := model.Build(mcfg, db, log)
	if err != nil {
		return err
	}
	if cacheStats != nil {
		defer func() { example.invalid(cacheStats.Summary()) }()
	}
	if err := applySavedModelSetting(ctx, db, chat, log); err != nil {
		return err
	}
	m, budget := chat.Settings()
	example.invalid("模型層", "model", m, "monthly_budget_usd", budget, "cache", mcfg.Cache)

	// 背景處理不跟著結束訊號取消：已經回 200 的訊息要處理完
	api := telegram.NewClient(cfg.Token, cfg.APIBase, log)
	// 待補問題在背景翻成商家慣用語言；關資料庫之前要等翻完
	var translating sync.WaitGroup
	defer translating.Wait()
	conv := &dialog.Bot{DB: db, Model: chat, Embedder: loadEmbedder(log), Log: log, Async: translating.Go}
	bot := &telegram.Bot{API: api, Store: db, Log: log, Dialog: conv}
	chat.OnBudgetExceeded = func(ctx context.Context, spent, budget float64) {
		bot.NotifyStaff(ctx, fmt.Sprintf("⚠️ Haiku 本月費用 $%.2f 已達上限 $%.2f，已自動切回本機模型（下個月恢復；要繼續用 Haiku 請到後台調高上限）", spent, budget))
	}
	hs, err := db.HandoffSetting(ctx)
	if err != nil {
		return err
	}
	example.invalid(bindStatus(hs))
	hook := &telegram.Webhook{Secret: [REDACTED], Dedup: db, Process: bot.Process, Log: log}
	mux := http.NewServeMux()
	mux.Handle(telegram.WebhookPath, hook)
	srv := &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	adm := &admin.Server{DB: db, Dir: acfg.Dir, EvalDir: acfg.EvalDir, Model: chat, Chat: chat, SelfTester: conv, Handoff: conv, Now: time.Now, Log: log}
	adminSrv := &http.Server{Handler: adm.Handler(), ReadHeaderTimeout: 10 * time.Second}

	ln, err := net.Listen("tcp", net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)))
	if err != nil {
		return err
	}
	adminLn, err := net.Listen("tcp", net.JoinHostPort(acfg.Host, strconv.Itoa(acfg.Port)))
	if err != nil {
		ln.Close()
		return err
	}
	example.invalid("伺服器已啟動", "url", "http://"+ln.Addr().String()+telegram.WebhookPath, "db", dbPath)
	example.invalid("後台已啟動", "url", "http://"+adminLn.Addr().String()+"/", "dir", acfg.Dir)
	if ready != nil {
		ready(ln.Addr().String(), adminLn.Addr().String())
	}
	// 關資料庫之前要等 setWebhook 和清理結束
	var bg sync.WaitGroup
	bgCtx, stopBG := context.WithCancel(ctx)
	bg.Go(func() {
		if err := api.SetWebhook(bgCtx, cfg.WebhookURL(), cfg.Secret); err != nil {
			log.Error("setWebhook 失敗，伺服器繼續執行（Telegram 送不進來，請檢查 token 和對外網址）")
			return
		}
		example.invalid("已設定 webhook", "url", cfg.WebhookURL())
	})
	bg.Go(func() { cleanupLoop(bgCtx, db, log) })
	bg.Go(func() {
		start := time.Now()
		if err := conv.WarmUp(bgCtx); err != nil {
			log.Warn("先算意圖案例的向量失敗，收到訊息時再算", "error", err.Error())
			return
		}
		example.invalid("已算好意圖案例的向量", "seconds", time.Since(start).Seconds())
	})
	defer func() { stopBG(); bg.Wait() }()

	errc := make(chan error, 2)
	go func() { errc <- srv.Serve(ln) }()
	go func() { errc <- adminSrv.Serve(adminLn) }()
	// 其中一個停了（出錯）就兩個一起停
	var serveErr error
	select {
	case serveErr = <-errc:
	case <-ctx.Done():
		example.invalid("收到結束訊號，停止伺服器")
	}
	shutdown, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	for _, s := range []*http.Server{srv, adminSrv} {
		if err := s.Shutdown(shutdown); err != nil {
			log.Error("停止伺服器失敗", "error", err.Error())
		}
	}
	hook.Wait()
	return serveErr
}

// loadEmbedder 載入找相似意圖案例用的 e5 模型（BOT_EMBED_MODEL）；載入失敗時記錄警告，分類照常、只是不附相似案例。
// 變數是為了讓測試換掉（不用載入 470 MB 的模型）
var loadEmbedder = func(log *slog.Logger) dialog.Embedder {
	m, err := embed.Load(embed.ModelDir())
	if err != nil {
		log.Warn("句子向量模型載入失敗，分類時不附相似案例（模型路徑見 README「多語言句子向量」）", "dir", embed.ModelDir(), "error", err.Error())
		return nil
	}
	return m
}

// applySavedModelSetting：後台存過模型設定時以資料庫為準（環境變數只是初始值）。
// 存的是 haiku 但現在沒有 API key 時記警告、用 AFM（上限照存的），不讓伺服器起不來
func applySavedModelSetting(ctx context.Context, db *store.DB, sw *model.Switch, log *slog.Logger) error {
	m, budget, ok, err := db.ModelSetting(ctx)
	if err != nil || !ok {
		return err
	}
	if m == model.ModelHaiku && !sw.HaikuAvailable() {
		log.Warn("後台設定用 Haiku，但沒有 API key（ANTHROPIC_API_KEY），先用本機模型")
		m = model.AFM
	}
	if err := sw.Apply(m, budget); err != nil {
		return fmt.Errorf("資料庫裡的模型設定不正確：%w", err)
	}
	return nil
}

func cleanupLoop(ctx context.Context, db *store.DB, log *slog.Logger) {
	run := func() {
		logs, updates, sessions, err := db.Cleanup(ctx, time.Now())
		if err != nil {
			log.Error("清理舊資料失敗", "error", err.Error())
			return
		}
		example.invalid("已清理舊資料", "conversation_log", logs, "processed_update", updates, "admin_session", sessions)
	}
	run()
	t := time.NewTicker(24 * time.Hour)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			run()
		}
	}
}
