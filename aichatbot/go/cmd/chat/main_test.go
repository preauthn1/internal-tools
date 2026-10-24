package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"aichatbot/internal/admin"
	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// 對話模式：eval/run_flows.py 用的指令和 FLOW_TRACE 格式（每行一個 JSON：回覆種類、回覆、語言；/reset、/pending 各一行）
func TestConverseTrace(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	intent := func(c string) map[string]model.Value { return map[string]model.Value{"intent": model.Text(c)} }
	m := model.NewScripted(map[string][]map[string]model.Value{"Intent": {intent("general"), intent("claim"), intent("general")}})
	bot := &dialog.Bot{DB: db, Model: m}
	in := strings.NewReader("/reset\n/start\n你好\n上次說打五折\n/pending\n/reset\nThanks\n/quit\n不會跑到這句\n")
	var out, trace bytes.Buffer
	if err := converse(context.Background(), bot, db, in, &out, &trace); err != nil {
		t.Fatal(err)
	}
	// /start 照 Telegram 平台層的規則：不分類，回打招呼（#54）
	want := `{"key":"reset"}
{"key":"greeting","lang":"zh","reply":"您好～有什麼可以幫您的嗎？"}
{"key":"general","lang":"zh","reply":"您好～有什麼可以幫您的嗎？"}
{"key":"claim","lang":"zh","reply":"這部分我沒辦法確認，已幫您記錄下來，會請店員回覆您喔！"}
{"key":"pending","pending":["#1 ❓ 顧客聲稱／更正：上次說打五折"]}
{"key":"reset"}
{"key":"general","lang":"en","reply":"You're welcome~ Let me know if you need anything else!"}
`
	if trace.String() != want {
		t.Errorf("trace：\n%s", trace.String())
	}
}

// 回答題庫（BOT_TESTSET）：每題一位新顧客；expect=answer 要回答而且內容對、no_answer 不能回答、either 回答了內容就要對。
// 驗證題（holdout）不跑；分開列出答錯的、該答沒答的（需要計算的另外列）和被擋下的（驗證或譯文驗證沒過的原因取自 traces）
func TestTestsetEval(t *testing.T) {
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	c := merchant.Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "02-1234-5678"
	c.Knowledge = []merchant.Section{{Title: "營業時間", Body: "11:00–20:00 營業。"}, {Title: "菜單", Body: "桂花烏龍 110 元"}}
	c.FAQ = []merchant.FAQ{{Key: "card", Question: "可以刷卡嗎", Answer: "信用卡需消費滿 200 元。"}}
	if err := prepareStore(context.Background(), db, c); err != nil {
		t.Fatal(err)
	}
	general := map[string]model.Value{"intent": model.Text("general")}
	faq := func(id string) map[string]model.Value {
		return map[string]model.Value{"faq": model.Text(id), "kind": model.Text("question")}
	}
	pick := func(id string) map[string]model.Value {
		return map[string]model.Value{"pick1": model.Text(id)}
	}
	answer := func(id, quote, reply string) map[string]model.Value {
		return map[string]model.Value{"source": model.Text(id), "quote": model.Text(quote), "answerable": model.Text("yes"), "reply": model.Text(reply)}
	}
	m := model.NewScripted(map[string][]map[string]model.Value{
		"Intent":    {general, general, general, general, general, general, general},
		"FAQ":       {faq("none"), faq("none"), faq("none"), faq("card"), faq("card"), faq("none"), faq("card")},
		"Translate": {{"translation": model.Text("Card needs 300.")}},
		"Pick":      {pick("K1"), pick("K2"), pick("none"), pick("K2"), pick("K2")},
		"Answer":    {answer("K1", "11:00–20:00 營業", "我們 11 點開門～"), answer("K2", "桂花烏龍 110 元", "桂花烏龍 120 元"), answer("K2", "桂花烏龍 110 元", "Oolong is 310.")},
	})
	bot := &dialog.Bot{DB: db, Model: m}
	path := filepath.Join(t.TempDir(), "set.json")
	set := `{"cases": [
		{"id": "a1", "q": "你們幾點開", "expect": "answer", "any": ["11"]},
		{"id": "a2", "q": "桂花烏龍多少", "expect": "answer", "any": ["110"]},
		{"id": "n1", "q": "可以帶狗嗎", "expect": "no_answer"},
		{"id": "f1", "q": "能刷卡嗎", "expect": "answer", "any": ["300"]},
		{"id": "f2", "q": "Can I pay by card?", "expect": "answer", "lang": "en"},
		{"id": "calc-1", "q": "桂花烏龍外送總共多少", "expect": "answer", "any": ["220"]},
		{"id": "calc-2", "q": "100塊可以刷卡嗎", "expect": "answer", "any": ["不"]},
		{"id": "h1", "q": "驗證題", "expect": "answer", "holdout": true}
	]}`
	if err := os.WriteFile(path, []byte(set), 0o644); err != nil {
		t.Fatal(err)
	}
	var out bytes.Buffer
	if err := testsetEval(context.Background(), bot, []string{path}, &out, "", false); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"answer       1/6",
		"no_answer    1/1",
		"========== 答錯的（回答了但不對）==========\nf1 能刷卡嗎　沒有提到 300 → [answer] 信用卡需消費滿 200 元。（faq:card）\n\n",
		"========== 該答沒答==========\n" +
			"a2 桂花烏龍多少 → [handoff] 我先幫您轉給店員，請稍候喔～（faq:none pick:K2 answer:yes verify:number_not_in_quote:120）\n" +
			"f2 Can I pay by card? → [handoff] Let me connect you with our staff, please wait a moment~（faq:card translate:number_not_in_quote:300 pick:K2 answer:yes verify:number_not_in_quote:310）\n",
		"========== 需要計算（ADR 0003，沒答或只回通則）==========\n" +
			"calc-1 桂花烏龍外送總共多少　應回答 → [handoff] 我先幫您轉給店員，請稍候喔～（faq:none pick:K2 answer:ai_failed）\n" +
			"calc-2 100塊可以刷卡嗎　沒有提到 不 → [answer] 信用卡需消費滿 200 元。（faq:card）\n",
		"========== 被擋下的（驗證沒過）==========\n" +
			"a2 桂花烏龍多少 → 桂花烏龍 120 元（verify:number_not_in_quote:120）\n" +
			"f2 Can I pay by card? → Card needs 300.（translate:number_not_in_quote:300）\n" +
			"f2 Can I pay by card? → Oolong is 310.（verify:number_not_in_quote:310）\n",
		"總計：2/7（29%）　回答 3 題　答錯 1 題　該答沒答 2 題　需要計算 2 題（其中回答了 1 題）　被擋下 2 題　略過驗證題 1 題",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("少了 %q：\n%s", want, out.String())
		}
	}
}

const evalDir = "../../../eval"

// 商家設定：預設從 eval/ 匯入；BOT_CONFIG_FILE 設了就讀那份匯出檔（後台做出的設定，#39）；BOT_LANGUAGES 兩種都套用
func TestLoadConfig(t *testing.T) {
	env := map[string]string{"BOT_CONFIG_DIR": evalDir}
	getenv := func(k string) string { return env[k] }
	c, err := loadConfig(getenv)
	if err != nil || c.Merchant.Name != "森日茶屋" {
		t.Fatalf("eval/：%s %v", c.Merchant.Name, err)
	}
	c.Merchant.Name = "匯出的店"
	data, err := merchant.Export(c)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
	env["BOT_CONFIG_FILE"], env["BOT_LANGUAGES"] = path, "zh,ko"
	c, err = loadConfig(getenv)
	if err != nil || c.Merchant.Name != "匯出的店" || len(c.Languages) != 2 || c.Languages[1].Code != "ko" {
		t.Errorf("BOT_CONFIG_FILE：%s %v %v", c.Merchant.Name, c.Languages, err)
	}
}

// 第二家測試商家（#66）的匯出檔只有商家自己的案例：內建案例是產品的一部分，照 eval/examples.jsonl 補上
// （改進迴圈改了內建案例，兩家店量到的是同一份），商家自己的接在後面
func TestLoadConfigAddsBuiltinExamples(t *testing.T) {
	env := map[string]string{"BOT_CONFIG_DIR": evalDir, "BOT_CONFIG_FILE": evalDir + "/stores/salon/config.json"}
	c, err := loadConfig(func(k string) string { return env[k] })
	if err != nil {
		t.Fatal(err)
	}
	tea, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	n := len(tea.Examples)
	if c.Merchant.Name != "霧光髮廊" || len(c.Examples) <= n || !reflect.DeepEqual(c.Examples[:n], tea.Examples) || c.Examples[n].Builtin {
		t.Errorf("%s：案例 %d 句（內建 %d 句），第一句商家案例 %+v", c.Merchant.Name, len(c.Examples), n, c.Examples[min(n, len(c.Examples)-1)])
	}
}

// #39：用後台 API 做出的設定（go/cmd/adminparity 的匯出檔）當 BOT_CONFIG_FILE，讀進來的設定和直接匯入 eval/ 的完全相同，
// 量測時機器人拿到的是同一份設定
func TestAdminConfigFileMatchesEvalImport(t *testing.T) {
	ctx := context.Background()
	env := map[string]string{"BOT_CONFIG_DIR": evalDir}
	getenv := func(k string) string { return env[k] }
	imported, err := loadConfig(getenv)
	if err != nil {
		t.Fatal(err)
	}
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := admin.BuildViaAPI(ctx, db, evalDir, imported); err != nil {
		t.Fatal(err)
	}
	c, err := merchant.Load(ctx, db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	data, err := merchant.Export(c)
	if err != nil {
		t.Fatal(err)
	}
	env["BOT_CONFIG_FILE"] = filepath.Join(t.TempDir(), "admin.json")
	if err := os.WriteFile(env["BOT_CONFIG_FILE"], data, 0o644); err != nil {
		t.Fatal(err)
	}
	viaAPI, err := loadConfig(getenv)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(imported, viaAPI) {
		t.Errorf("讀進來的設定不同：\n匯入：%+v\n後台：%+v", imported, viaAPI)
	}
}

// 題庫的相對日期標記：和 eval/testdates.py 的 doctest 相同
func TestResolveTestDates(t *testing.T) {
	d := time.Date(2026, 10, 7, 15, 0, 0, 0, time.Local)
	got := resolveTestDates("{{date:+1:M/D}} {{today:+1:MM-DD}} {{weekday:next:2:M/D}} {{weekday:this:1:M月D日}} {{weekday:near:3:YYYY-MM-DD}}", d)
	if want := "10/20 10-08 10/13 10月5日 2026-10-07"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got := resolveTestDates("{{weekday:next:1:M/D}} {{weekday:near:1:M/D}} {{date:-1:M月D號}}", time.Date(2026, 10, 11, 0, 0, 0, 0, time.Local)); got != "10/12 10/12 10月25號" {
		t.Errorf("got %q", got)
	}
}
