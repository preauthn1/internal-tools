package merchant

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"
)

// valid 是一份最小、能通過檢查的設定
func valid() Config {
	c := Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "茶館", "02-1234-5678"
	return c
}

// FAQ 的問題、答案必填，錯誤訊息指出第幾題
func TestValidateFAQRequired(t *testing.T) {
	c := valid()
	c.FAQ = []FAQ{{Key: "a", Question: "營業時間", Answer: "11 點"}, {Key: "b", Question: " ", Answer: ""}}
	err := c.Validate()
	if err == nil || !strings.Contains(err.Error(), "第 2 題 FAQ 的問題是必填") || !strings.Contains(err.Error(), "第 2 題 FAQ 的答案是必填") {
		t.Errorf("got %v", err)
	}
}

// 知識段落的標題、內容必填；標題（去掉前後空白後）不能重複
func TestValidateKnowledge(t *testing.T) {
	c := valid()
	c.Knowledge = []Section{{Title: "營業時間", Body: "11 點開"}, {Title: "", Body: " "}, {Title: " 營業時間", Body: "週一公休"}}
	err := c.Validate()
	for _, want := range []string{"第 2 段知識的標題是必填", "第 2 段知識的內容是必填", "知識段落標題重複：營業時間"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("少了「%s」：%v", want, err)
		}
	}
}

// 顧客語言：代碼要在可勾選的清單裡、不能重複、至少勾一個
func TestValidateLanguages(t *testing.T) {
	for _, c := range []struct {
		langs []Language
		want  string
	}{
		{[]Language{{"zh", true}, {"xx", true}}, "不認得的語言代碼：xx"},
		{[]Language{{"zh", true}, {"zh", false}}, "語言重複：zh"},
		{[]Language{{"zh", false}}, "至少要勾一種顧客語言"},
	} {
		cfg := valid()
		cfg.Languages = c.langs
		if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%v：got %v, want %s", c.langs, err, c.want)
		}
	}
	cfg := valid()
	cfg.Languages = []Language{{"zh", true}, {"ko", true}, {"th", false}}
	if err := cfg.Validate(); err != nil {
		t.Errorf("中英日以外的語言也可以勾：%v", err)
	}
}

// 存商家資訊只換掉名稱、介紹、聯絡方式和知識段落，其他區塊（FAQ、意圖案例、固定回覆……）不動
func TestSaveInfoKeepsOtherBlocks(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	info := Info{Name: "森日茶屋 二號店", Intro: "新的介紹", Contact: "02-0000-0000", Knowledge: []Section{{Title: "停車", Body: "附近有停車場"}}}
	if err := SaveInfo(ctx, db, info); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	want := c
	want.Merchant.Name, want.Merchant.Intro, want.Merchant.Contact = info.Name, info.Intro, info.Contact
	want.Knowledge = info.Knowledge
	a, _ := Export(want)
	b, _ := Export(got)
	if !bytes.Equal(a, b) {
		t.Errorf("存商家資訊後不對\nwant:\n%s\ngot:\n%s", a, b)
	}
}

// 存 FAQ：空的代號由程式給一個還沒用過的 faq-N，原本的代號不變，順序照送來的
func TestSaveFAQGeneratesKeys(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	if err := SaveInfo(ctx, db, InfoOf(valid())); err != nil {
		t.Fatal(err)
	}
	faq := []FAQ{
		{Question: "新問題一", Answer: "答一"},
		{Key: "faq-1", Question: "舊問題", Answer: "舊答"},
		{Question: "新問題二", AltQuestions: []string{" 別的問法 ", " "}, Answer: "答二"},
	}
	if err := SaveFAQ(ctx, db, faq); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	var keys []string
	for _, f := range got.FAQ {
		keys = append(keys, f.Key+"="+f.Question)
	}
	if strings.Join(keys, ",") != "faq-2=新問題一,faq-1=舊問題,faq-3=新問題二" {
		t.Errorf("got %v", keys)
	}
	// 其他問法去掉前後空白，空白的項目拿掉（和後台頁面送的一致）
	if alt := got.FAQ[2].AltQuestions; len(alt) != 1 || alt[0] != "別的問法" {
		t.Errorf("其他問法：%q", alt)
	}
}

// 存功能開關：只改那一個開關（沒有那一列就新增），其他區塊不動；只能是 booking、handoff
func TestSaveFeature(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c := valid()
	c.Features = []Feature{{Name: "booking", Enabled: true}}
	c.Replies = []Reply{{Key: "greeting", Lang: "zh", Text: "嗨"}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	if err := SaveFeature(ctx, db, "handoff", false); err != nil {
		t.Fatal(err)
	}
	if err := SaveFeature(ctx, db, "booking", false); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	want := c
	want.Features = []Feature{{Name: "booking", Enabled: false}, {Name: "handoff", Enabled: false}}
	a, _ := Export(want)
	b, _ := Export(got)
	if !bytes.Equal(a, b) {
		t.Errorf("存功能開關後不對\nwant:\n%s\ngot:\n%s", a, b)
	}
	var invalid *InvalidError
	if err := SaveFeature(ctx, db, "pay", true); !errors.As(err, &invalid) {
		t.Errorf("不認得的開關：%v", err)
	}
}

// 還沒有商家設定時，FAQ、語言不能先存（商家資訊的必填欄位還是空的），要先填商家資訊
func TestSaveBeforeInfo(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	errs := []error{
		SaveFAQ(ctx, db, []FAQ{{Question: "營業時間", Answer: "11 點"}}),
		SaveLanguages(ctx, db, LanguageSettings{Preferred: "zh", Languages: []Language{{"zh", true}}}),
		SaveFeature(ctx, db, "booking", false),
	}
	for _, err := range errs {
		if err == nil || err.Error() != "請先填商家資訊" {
			t.Errorf("got %v", err)
		}
	}
	if _, err := Load(ctx, db); err != ErrNotConfigured {
		t.Errorf("不該存進去：%v", err)
	}
}

// FAQ 代號必填
func TestValidateFAQKeyRequired(t *testing.T) {
	c := valid()
	c.FAQ = []FAQ{{Key: " ", Question: "營業時間", Answer: "11 點"}}
	if err := c.Validate(); err == nil || !strings.Contains(err.Error(), "第 1 題 FAQ 的代號是必填") {
		t.Errorf("got %v", err)
	}
}
