// chat 是終端機對話工具（照搬 Swift 版的 bot 執行檔），題庫用它對 Go 版量測：
//
//	go -C go build -o bin/chat ./cmd/chat && go/bin/chat          對話模式（標準輸入一行一句）
//	BOT_INTENT_EVAL=eval/intent_set.json go/bin/chat             只測分類（ROUNDS=輪數，預設 1）
//	BOT_TESTSET=eval/testset.json,eval/faq_set.json go/bin/chat  回答題庫（每題一位新顧客，列出答錯和被擋下的；不跑驗證題）
//	python3 eval/run_flows.py --bot go-chat --exe go/bin/chat    流程題（對話模式＋FLOW_TRACE）
//
// 回答題庫另外把每題結果寫到 eval/results/testset-*.json（印「結果：路徑」，eval/snapshot.py 讀）；HOLDOUT=1 只跑驗證題。
// 分類和回答題庫印每則平均時間；結束時在 stderr 印 AI 用量（token、估算費用）。BOT_LANGUAGES=zh,en,ja,ko 換掉店家啟用的語言。
//
// 商家設定每次啟動時從 BOT_CONFIG_DIR（沒設就是 eval/ 或 ../eval/）的森日茶屋檔案匯入到暫存的資料庫，不寫入正式資料；
// 設了 BOT_CONFIG_FILE 時改讀那份匯出檔（例如 go/cmd/adminparity 用後台 API 做出的設定）。
// 模型和 AI 回答快取的環境變數同 bot serve（BOT_MODEL、BOT_AI_CACHE…，見 README）；句子向量模型見 BOT_EMBED_MODEL。
// 今天可用 BOT_TODAY=YYYY-MM-DD 指定（題庫的相對日期和對話的「今天」都用它）。
//
// 對話模式的指令：/reset 換一位新顧客  /as 使用者ID 切換顧客  /pending 待處理清單  /quit 離開。
// 設了 FLOW_TRACE=檔案 時，每句的回覆種類、回覆、語言（/reset、/pending 也各一行）以 JSON 寫進檔案（eval/run_flows.py 讀）
package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/embed"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "錯誤：", err)
		os.Exit(1)
	}
}

// prepareStore 存好商家設定，綁一個假的員工群組：沒綁群組不會轉人工（#55），量測要和已經綁好群組的商家一樣
func prepareStore(ctx context.Context, db *store.DB, c merchant.Config) error {
	if _, err := merchant.Save(ctx, db.SQL(), c); err != nil {
		return err
	}
	s, err := db.HandoffSetting(ctx)
	if err != nil {
		return err
	}
	if ok, err := db.BindStaffChat(ctx, s.BindCode, -1); err != nil || !ok {
		return fmt.Errorf("綁定假的員工群組失敗：%v", err)
	}
	return nil
}

func run() error {
	ctx := context.Background()
	log := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))
	tmp, err := os.MkdirTemp("", "chat-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	db, err := store.Open(filepath.Join(tmp, "bot.sqlite"))
	if err != nil {
		return err
	}
	defer db.Close()
	c, err := loadConfig(os.Getenv)
	if err != nil {
		return err
	}
	if err := prepareStore(ctx, db, c); err != nil {
		return err
	}
	mcfg, err := model.ConfigFromEnv(os.Getenv)
	if err != nil {
		return err
	}
	usage := &usageTally{Ledger: db}
	chat, stats, err := model.Build(mcfg, usage, log)
	if err != nil {
		return err
	}
	defer func() { fmt.Fprintln(os.Stderr, usage.Summary()) }()
	bot := &dialog.Bot{DB: db, Model: chat, Log: log, Now: now}
	if m, err := embed.Load(embed.ModelDir()); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️ 句子向量模型載入失敗，分類時不附相似案例（%s）：%v\n", embed.ModelDir(), err)
	} else {
		defer m.Close()
		bot.Embedder = m
	}
	if k, err := strconv.Atoi(os.Getenv("BOT_EXAMPLES_K")); err == nil {
		bot.K = k
	}
	if err := bot.WarmUp(ctx); err != nil {
		return err
	}
	printStats := func(w io.Writer) {
		if stats != nil {
			fmt.Fprintln(w, stats.Summary())
		}
	}

	if path := os.Getenv("BOT_INTENT_EVAL"); path != "" {
		rounds, _ := strconv.Atoi(os.Getenv("ROUNDS"))
		if err := intentEval(ctx, bot, path, max(rounds, 1), os.Stdout); err != nil {
			return err
		}
		printStats(os.Stdout)
		return nil
	}
	if paths := os.Getenv("BOT_TESTSET"); paths != "" {
		dir, err := filepath.Abs(filepath.Join(configDir(os.Getenv), "results"))
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		result := filepath.Join(dir, fmt.Sprintf("testset-%d.json", time.Now().UnixNano()))
		if err := testsetEval(ctx, bot, strings.Split(paths, ","), os.Stdout, result, os.Getenv("HOLDOUT") == "1"); err != nil {
			return err
		}
		printStats(os.Stdout)
		return nil
	}

	var trace io.Writer = io.Discard
	if p := os.Getenv("FLOW_TRACE"); p != "" {
		f, err := os.OpenFile(p, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
		if err != nil {
			return err
		}
		defer f.Close()
		trace = f
	}
	fmt.Printf("🍵 %s客服（Go 版）\n   /reset 換新顧客  /as 使用者ID 切換顧客  /pending 待處理清單  /quit 離開\n\n", c.Merchant.Name)
	err = converse(ctx, bot, db, os.Stdin, os.Stdout, trace)
	// 對話模式的 stdout 被 eval/run_flows.py 丟掉，快取統計印到 stderr
	printStats(os.Stderr)
	return err
}

// loadConfig 是這次用的商家設定：BOT_CONFIG_FILE 設了就讀那份匯出檔，否則從 configDir 的森日茶屋檔案匯入；再套用 BOT_LANGUAGES。
// 匯出檔沒有內建案例時（第二家測試商家 eval/stores/salon/，#66），內建案例照 configDir 的 examples.jsonl 補在前面
func loadConfig(getenv func(string) string) (merchant.Config, error) {
	var c merchant.Config
	var err error
	if path := getenv("BOT_CONFIG_FILE"); path != "" {
		c, err = merchant.ReadFile(path)
		if err == nil && !slices.ContainsFunc(c.Examples, func(e merchant.Example) bool { return e.Builtin }) {
			var tea merchant.Config
			if tea, err = merchant.FromEvalDir(configDir(getenv)); err == nil {
				builtin := slices.DeleteFunc(tea.Examples, func(e merchant.Example) bool { return !e.Builtin })
				c.Examples = append(builtin, c.Examples...)
			}
		}
	} else {
		c, err = merchant.FromEvalDir(configDir(getenv))
	}
	if err != nil {
		return c, err
	}
	return withLanguages(c, getenv("BOT_LANGUAGES")), nil
}

// configDir：BOT_CONFIG_DIR，沒設就找工作目錄的 eval/，再找 ../eval/（在 go/ 底下執行時）
func configDir(getenv func(string) string) string {
	if d := getenv("BOT_CONFIG_DIR"); d != "" {
		return d
	}
	if _, err := os.Stat("eval/store.json"); err == nil {
		return "eval"
	}
	return filepath.Join("..", "eval")
}

// withLanguages：codes（BOT_LANGUAGES，逗號分隔）不是空字串時，店家只啟用這些語言（eval/verify.py 量跨語言題時開韓文）
func withLanguages(c merchant.Config, codes string) merchant.Config {
	if strings.TrimSpace(codes) == "" {
		return c
	}
	c.Languages = nil
	for _, code := range strings.Split(codes, ",") {
		c.Languages = append(c.Languages, merchant.Language{Code: strings.TrimSpace(code), Enabled: true})
	}
	return c
}

// converse 是對話模式：一行一句，/ 開頭是指令；trace 每句寫一行 JSON
func converse(ctx context.Context, bot *dialog.Bot, db *store.DB, in io.Reader, out io.Writer, trace io.Writer) error {
	const platform = "terminal"
	userID := "user1"
	writeTrace := func(v map[string]any) error {
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		_, err = trace.Write(append(data, '\n'))
		return err
	}
	sc := bufio.NewScanner(in)
	for {
		fmt.Fprintf(out, "👤 %s：", userID)
		if !sc.Scan() {
			break
		}
		line := strings.TrimSpace(sc.Text())
		switch {
		case line == "":
			continue
		case line == "/quit" || line == "/exit":
			fmt.Fprintln(out, "👋")
			return nil
		case line == "/reset":
			if err := bot.ResetCustomer(ctx, platform, userID); err != nil {
				return err
			}
			if err := writeTrace(map[string]any{"key": "reset"}); err != nil {
				return err
			}
			fmt.Fprintln(out, "🔄 已換一位新顧客")
			continue
		case line == "/pending":
			tickets, err := db.OpenTickets(ctx)
			if err != nil {
				return err
			}
			items := []string{}
			for _, t := range tickets {
				items = append(items, fmt.Sprintf("#%d %s", t.ID, t.Text))
			}
			if err := writeTrace(map[string]any{"key": "pending", "pending": items}); err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Fprintln(out, "📭 沒有待處理項目")
			} else {
				fmt.Fprintln(out, "📬 待員工處理：\n   "+strings.Join(items, "\n   "))
			}
			continue
		case strings.HasPrefix(line, "/as "):
			userID = strings.TrimSpace(strings.TrimPrefix(line, "/as "))
			fmt.Fprintf(out, "👥 切換到顧客 %s\n", userID)
			continue
		}
		start := time.Now()
		receive := bot.Receive
		if dialog.IsBotCommand(line) {
			receive = bot.Greet // 和 Telegram 平台層一樣：/start、/help 不送進分類，回打招呼（#54）
		}
		r, err := receive(ctx, dialog.Message{Platform: platform, UserID: userID, ChatID: userID, Text: line})
		if err != nil {
			return err
		}
		c, err := bot.LoadCustomer(ctx, platform, userID)
		if err != nil {
			return err
		}
		if err := writeTrace(map[string]any{"key": r.Key, "reply": r.Text, "lang": c.Lang}); err != nil {
			return err
		}
		if r.Silent {
			fmt.Fprintln(out, "🙋 （轉人工中，AI 不回覆）")
			continue
		}
		fmt.Fprintf(out, "🤖 %s  \x1b[90m[%s·%s] %.2fs\x1b[0m\n", r.Text, r.Key, c.Lang, time.Since(start).Seconds())
	}
	return sc.Err()
}

// intentEval 只測分類（BOT_INTENT_EVAL）：和正式流程同一個分類，每題給一位新顧客，依題目設定狀態。
// 題組的舊代號（Swift 版的 14 類）合併成 general 再比對；另外印「10 類換算」（不含想找真人的題目），和 Swift 版的基準比
func intentEval(ctx context.Context, bot *dialog.Bot, path string, rounds int, out io.Writer) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var set struct {
		Items []struct {
			Text   string   `json:"text"`
			Accept []string `json:"accept"`
			State  string   `json:"state"`
		} `json:"items"`
	}
	if err := json.Unmarshal([]byte(resolveTestDates(string(data), now())), &set); err != nil {
		return fmt.Errorf("%s：%w", path, err)
	}
	if len(set.Items) == 0 {
		return errors.New(path + "：沒有題目")
	}
	for i, item := range set.Items {
		if len(item.Accept) == 0 {
			return fmt.Errorf("%s：第 %d 題（%s）沒有 accept", path, i+1, item.Text)
		}
	}
	merged := map[string]bool{"greeting": true, "thanks": true, "goodbye": true, "small_talk": true, "store_question": true}
	label := func(s string) string {
		if merged[s] {
			return "general"
		}
		return s
	}
	results := make([][]string, len(set.Items))
	var elapsed time.Duration // 每次分類的時間加總（開快取時含命中）
	for range rounds {
		for i, item := range set.Items {
			var c dialog.Customer
			switch {
			case item.State == "ticket":
				c.OpenTicket = true
			case strings.HasPrefix(item.State, "booking:"):
				c.AskingField = strings.TrimPrefix(item.State, "booking:")
			}
			start := time.Now()
			got, err := bot.Classify(ctx, item.Text, c)
			if err != nil {
				return err
			}
			elapsed += time.Since(start)
			results[i] = append(results[i], got)
		}
	}
	type tally struct{ ok, total int }
	byLabel := map[string]*tally{}
	var all, ten tally
	unstable := 0
	var wrong []string
	for i, item := range set.Items {
		accept := map[string]bool{}
		for _, a := range item.Accept {
			accept[label(a)] = true
		}
		first := label(item.Accept[0])
		if byLabel[first] == nil {
			byLabel[first] = &tally{}
		}
		ok, seen := 0, map[string]bool{}
		for _, r := range results[i] {
			seen[r] = true
			if accept[r] {
				ok++
			}
		}
		if len(seen) > 1 {
			unstable++
		}
		byLabel[first].ok += ok
		byLabel[first].total += rounds
		all.ok, all.total = all.ok+ok, all.total+rounds
		if first != "human_request" {
			ten.ok, ten.total = ten.ok+ok, ten.total+rounds
		}
		if ok < rounds {
			state := ""
			if item.State != "" {
				state = "［" + item.State + "］"
			}
			wrong = append(wrong, fmt.Sprintf("%s%s　應為 %s　→ %s", state, item.Text, strings.Join(item.Accept, "/"), strings.Join(results[i], "、")))
		}
	}
	labels := make([]string, 0, len(byLabel))
	for l := range byLabel {
		labels = append(labels, l)
	}
	sort.Strings(labels)
	fmt.Fprintf(out, "========== 各類別（%d 輪）==========\n", rounds)
	for _, l := range labels {
		fmt.Fprintf(out, "%-16s %d/%d\n", l, byLabel[l].ok, byLabel[l].total)
	}
	fmt.Fprintln(out, "\n========== 答錯的（各輪結果）==========")
	for _, w := range wrong {
		fmt.Fprintln(out, w)
	}
	pct := func(t tally) float64 { return float64(t.ok) / float64(max(t.total, 1)) * 100 }
	fmt.Fprintf(out, "\n總計：%d/%d（%.0f%%）　不穩定：%d 題\n", all.ok, all.total, pct(all), unstable)
	fmt.Fprintf(out, "10 類換算（不含想找真人）：%d/%d（%.0f%%）\n", ten.ok, ten.total, pct(ten))
	fmt.Fprintf(out, "每則平均 %.2f 秒（%d 則）\n", elapsed.Seconds()/float64(max(all.total, 1)), all.total)
	return nil
}
