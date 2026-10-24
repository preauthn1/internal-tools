package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// newTestsetBot：只有一題 FAQ（刷卡）的店，模型照 script 回答
func newTestsetBot(t *testing.T, script map[string][]map[string]model.Value) *dialog.Bot {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	c := merchant.Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "02-1234-5678"
	c.FAQ = []merchant.FAQ{{Key: "card", Question: "可以刷卡嗎", Answer: "信用卡需消費滿 200 元。"}}
	if err := prepareStore(context.Background(), db, c); err != nil {
		t.Fatal(err)
	}
	return &dialog.Bot{DB: db, Model: model.NewScripted(script)}
}

func writeSet(t *testing.T, set string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "set.json")
	if err := os.WriteFile(path, []byte(set), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

var (
	generalIntent = map[string]model.Value{"intent": model.Text("general")}
	noFAQ         = map[string]model.Value{"faq": model.Text("none"), "kind": model.Text("question")}
	cardFAQ       = map[string]model.Value{"faq": model.Text("card"), "kind": model.Text("question")}
	noPick        = map[string]model.Value{"pick1": model.Text("none")}
)

// 誘答題（group=trap）：資料裡沒有答案，回答了（answer）而且判錯就是編造；轉人工或說不知道（不是 answer）都算對
func TestTestsetEvalTrap(t *testing.T) {
	bot := newTestsetBot(t, map[string][]map[string]model.Value{
		"Intent": {generalIntent, generalIntent, generalIntent},
		"FAQ":    {noFAQ, cardFAQ, cardFAQ},
		"Pick":   {noPick},
	})
	path := writeSet(t, `{"cases": [
		{"id": "t1", "q": "可以帶狗嗎", "expect": "no_answer", "group": "trap"},
		{"id": "t2", "q": "可以用街口嗎", "expect": "no_answer", "group": "trap"},
		{"id": "a1", "q": "能刷卡嗎", "expect": "answer", "any": ["200"]}
	]}`)
	var out bytes.Buffer
	if err := testsetEval(context.Background(), bot, []string{path}, &out, "", false); err != nil {
		t.Fatal(err)
	}
	if want := "誘答題 編造 1 題（共 2 題）"; !strings.Contains(out.String(), want) {
		t.Errorf("少了 %q：\n%s", want, out.String())
	}
}

// 每題結果（eval/snapshot.py 逐題比對用）：結果（ok／wrong／missed）、回覆種類、挑到的 FAQ，寫成 JSON 並印出路徑
func TestTestsetEvalResults(t *testing.T) {
	bot := newTestsetBot(t, map[string][]map[string]model.Value{
		"Intent": {generalIntent, generalIntent, generalIntent},
		"FAQ":    {noFAQ, cardFAQ, cardFAQ},
		"Pick":   {noPick},
	})
	path := writeSet(t, `{"cases": [
		{"id": "m1", "q": "幾點開", "expect": "answer"},
		{"id": "w1", "q": "可以用街口嗎", "expect": "no_answer"},
		{"id": "o1", "q": "能刷卡嗎", "expect": "answer", "any": ["200"]}
	]}`)
	result := filepath.Join(t.TempDir(), "results.json")
	var out bytes.Buffer
	if err := testsetEval(context.Background(), bot, []string{path}, &out, result, false); err != nil {
		t.Fatal(err)
	}
	if want := "結果：" + result + "\n"; !strings.Contains(out.String(), want) {
		t.Errorf("少了 %q：\n%s", want, out.String())
	}
	data, err := os.ReadFile(result)
	if err != nil {
		t.Fatal(err)
	}
	want := `[{"id":"m1","result":"missed","key":"handoff","pick":"none"},{"id":"w1","result":"wrong","key":"answer","pick":"card"},{"id":"o1","result":"ok","key":"answer","pick":"card"}]`
	if strings.TrimSpace(string(data)) != want {
		t.Errorf("結果：\n%s", data)
	}
}

// 跨語言題（group=xlang）：照一般的判分（內容和語言），另外印答對幾題；店家沒開的語言（韓文）不會回答
func TestTestsetEvalXlang(t *testing.T) {
	bot := newTestsetBot(t, map[string][]map[string]model.Value{
		"Intent":    {generalIntent},
		"FAQ":       {cardFAQ},
		"Translate": {{"translation": model.Text("Credit cards need NT$200 or more.")}},
	})
	path := writeSet(t, `{"cases": [
		{"id": "x1", "q": "Can I pay with a card?", "expect": "answer", "any": ["200"], "lang": "en", "group": "xlang"},
		{"id": "x2", "q": "카드 결제 되나요?", "expect": "answer", "any": ["200"], "lang": "ko", "group": "xlang"}
	]}`)
	var out bytes.Buffer
	if err := testsetEval(context.Background(), bot, []string{path}, &out, "", false); err != nil {
		t.Fatal(err)
	}
	if want := "跨語言 答對 1/2"; !strings.Contains(out.String(), want) {
		t.Errorf("少了 %q：\n%s", want, out.String())
	}
	if strings.Contains(out.String(), "誘答題") {
		t.Errorf("沒有誘答題就不印：\n%s", out.String())
	}
	// 每則平均時間（eval/verify.py 讀）：時間不固定，只看格式和則數
	if !regexp.MustCompile(`每則平均 \d+\.\d\d 秒（2 則）`).MatchString(out.String()) {
		t.Errorf("少了每則平均時間：\n%s", out.String())
	}
}

// 分類題組（BOT_INTENT_EVAL）也印每則平均時間（每輪每句算一則）
func TestIntentEvalAverageTime(t *testing.T) {
	bot := newTestsetBot(t, map[string][]map[string]model.Value{"Intent": {generalIntent, generalIntent, generalIntent, generalIntent}})
	path := writeSet(t, `{"items": [{"text": "你好", "accept": ["greeting"]}, {"text": "我要找老闆", "accept": ["human_request"]}]}`)
	var out bytes.Buffer
	if err := intentEval(context.Background(), bot, path, 2, &out); err != nil {
		t.Fatal(err)
	}
	if !regexp.MustCompile(`每則平均 \d+\.\d\d 秒（4 則）`).MatchString(out.String()) {
		t.Errorf("少了每則平均時間：\n%s", out.String())
	}
}

// HOLDOUT=1 只跑驗證題（eval/snapshot.py 比對驗證題組用），其他題略過
func TestTestsetEvalHoldoutOnly(t *testing.T) {
	bot := newTestsetBot(t, map[string][]map[string]model.Value{"Intent": {generalIntent}, "FAQ": {cardFAQ}})
	path := writeSet(t, `{"cases": [
		{"id": "d1", "q": "幾點開", "expect": "answer"},
		{"id": "h1", "q": "能刷卡嗎", "expect": "answer", "any": ["200"], "holdout": true}
	]}`)
	var out bytes.Buffer
	if err := testsetEval(context.Background(), bot, []string{path}, &out, "", true); err != nil {
		t.Fatal(err)
	}
	if want := "總計：1/1（100%）　回答 1 題　答錯 0 題　該答沒答 0 題　需要計算 0 題（其中回答了 0 題）　被擋下 0 題　略過驗證題 0 題"; !strings.Contains(out.String(), want) {
		t.Errorf("少了 %q：\n%s", want, out.String())
	}
}

// BOT_LANGUAGES 換掉店家啟用的語言（跨語言題要開韓文；流程題的「不支援的語言」要維持中英日，所以不改預設值）
func TestWithLanguages(t *testing.T) {
	c := merchant.Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "02-1234-5678"
	got := withLanguages(c, " zh, ko ")
	want := []merchant.Language{{Code: "zh", Enabled: true}, {Code: "ko", Enabled: true}}
	if !slices.Equal(got.Languages, want) {
		t.Errorf("got %v", got.Languages)
	}
	if err := got.Validate(); err != nil {
		t.Error(err)
	}
	if got := withLanguages(c, ""); !slices.Equal(got.Languages, c.Languages) {
		t.Errorf("沒設就不換：%v", got.Languages)
	}
}
