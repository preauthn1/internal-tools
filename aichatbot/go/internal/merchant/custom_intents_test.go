package merchant

import (
	"context"
	"errors"
	"strings"
	"testing"
)

// customMerchant 是有一段知識的最小商家設定
func customMerchant() Config {
	c := Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "電話 02-1234-5678"
	c.Knowledge = []Section{{Title: "包場", Body: "包場請提前一週預約。"}}
	return c
}

func giftBox() CustomIntent {
	return CustomIntent{
		Name: "禮盒", Description: "詢問禮盒或送禮", Examples: []string{"有賣禮盒嗎", "想買來送人", "過年送禮推薦"},
		Action: ActionFixedText, Text: "禮盒請洽門市。",
	}
}

// 自訂意圖存進去讀得出來；代號自動產生 custom_N（第一個沒用過的），已有的代號保留
func TestCustomIntentsSaveAndLoad(t *testing.T) {
	ctx := context.Background()
	db := openDB(t)
	c := customMerchant()
	party := CustomIntent{Code: "custom_2", Name: "包場", Description: "詢問包場", Examples: []string{"可以包場嗎", "想包下整間店", "包場怎麼收費"},
		Action: ActionKnowledge, Knowledge: []string{"包場"}}
	c.CustomIntents = []CustomIntent{party, giftBox()}
	if _, err := Save(ctx, db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	if len(got.CustomIntents) != 2 || got.CustomIntents[0].Code != "custom_2" || got.CustomIntents[1].Code != "custom_1" {
		t.Fatalf("讀出：%+v", got.CustomIntents)
	}
	g := got.CustomIntents[1]
	if g.Name != "禮盒" || strings.Join(g.Examples, "/") != "有賣禮盒嗎/想買來送人/過年送禮推薦" || g.Action != ActionFixedText || g.Text != "禮盒請洽門市。" || len(g.Knowledge) != 0 {
		t.Errorf("禮盒：%+v", g)
	}
	if k := got.CustomIntents[0].Knowledge; len(k) != 1 || k[0] != "包場" {
		t.Errorf("包場的知識段落：%v", k)
	}
}

// 後台存自訂意圖：只換掉自訂意圖；例句和題庫重複回 *OverlapError，不存
func TestSaveCustomIntents(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	if err := SaveCustomIntents(ctx, db, TestSet{}, []CustomIntent{giftBox()}); err != ErrInfoFirst {
		t.Errorf("還沒有商家設定：%v", err)
	}
	c := customMerchant()
	c.FAQ = []FAQ{{Key: "a", Question: "營業時間", AltQuestions: []string{}, Answer: "11 點"}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	tests := TestSet{Normalize("想買來送人"): {"testset"}}
	var overlap *OverlapError
	if err := SaveCustomIntents(ctx, db, tests, []CustomIntent{giftBox()}); !errors.As(err, &overlap) || !strings.Contains(err.Error(), "想買來送人") {
		t.Errorf("和題庫重複：%v", err)
	}
	if err := SaveCustomIntents(ctx, db, TestSet{}, []CustomIntent{giftBox()}); err != nil {
		t.Fatal(err)
	}
	got, _ := Load(ctx, db)
	if len(got.CustomIntents) != 1 || got.CustomIntents[0].Code != "custom_1" || len(got.FAQ) != 1 {
		t.Errorf("存好後：%+v", got)
	}
	bad := giftBox()
	bad.Examples = nil
	var invalid *InvalidError
	if err := SaveCustomIntents(ctx, db, TestSet{}, []CustomIntent{bad}); !errors.As(err, &invalid) {
		t.Errorf("沒通過檢查：%v", err)
	}
}

// 案例庫可以新增屬於自訂意圖的案例；刪掉自訂意圖時，屬於它的商家案例一起刪掉（內建意圖的不動）
func TestExamplesOfCustomIntent(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	if _, err := Save(ctx, db, customMerchant()); err != nil {
		t.Fatal(err)
	}
	if _, err := AddExample(ctx, db, TestSet{}, "送禮用的", "custom_1"); !errors.Is(err, ErrInvalidExample) {
		t.Errorf("還沒有這個自訂意圖：%v", err)
	}
	if err := SaveCustomIntents(ctx, db, TestSet{}, []CustomIntent{giftBox()}); err != nil {
		t.Fatal(err)
	}
	if _, err := AddExample(ctx, db, TestSet{}, "有賣禮盒嗎？", "booking"); !errors.Is(err, ErrDuplicateExample) {
		t.Errorf("和自訂意圖的例句重複：%v", err)
	}
	for _, e := range [][2]string{{"送禮用的", "custom_1"}, {"我要訂位", "booking"}} {
		if _, err := AddExample(ctx, db, TestSet{}, e[0], e[1]); err != nil {
			t.Fatalf("%s：%v", e[0], err)
		}
	}
	// 同一次存檔刪掉禮盒（custom_1）、新增另一個：新的拿到 custom_1，但禮盒的案例不能變成它的
	other := giftBox()
	other.Name, other.Examples = "外送", []string{"可以外送嗎", "有送到府嗎", "外送範圍"}
	if err := SaveCustomIntents(ctx, db, TestSet{}, []CustomIntent{other}); err != nil {
		t.Fatal(err)
	}
	list, _ := ListExamples(ctx, db)
	if len(list) != 1 || list[0].Intent != "booking" {
		t.Errorf("刪掉自訂意圖後的案例：%+v", list)
	}
}

// 代號只能是空的（程式給）或 custom_1～custom_10，而且不能重複：不能搶走內建意圖的代號
func TestValidateCustomIntentCodes(t *testing.T) {
	for _, codes := range [][]string{{"general"}, {"custom_11"}, {"custom_1", "custom_1"}} {
		cfg := customMerchant()
		for i, code := range codes {
			ci := giftBox()
			ci.Code, ci.Name = code, ci.Name+strings.Repeat("！", i)
			for j := range ci.Examples {
				ci.Examples[j] += strings.Repeat("a", i)
			}
			cfg.CustomIntents = append(cfg.CustomIntents, ci)
		}
		if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "自訂意圖的代號不正確") {
			t.Errorf("%v：%v", codes, err)
		}
	}
}

// 檢查規則：最多 10 個；名稱、說明必填；例句去掉空白後至少 3 句、不重複、不能和內建案例相同；動作要認得；
// 固定文字必填；只根據知識回答要指定存在的段落；名稱不能和內建意圖或其他自訂意圖重複
func TestValidateCustomIntents(t *testing.T) {
	for _, c := range []struct {
		change func(*CustomIntent)
		want   string
	}{
		{func(ci *CustomIntent) { ci.Name = " " }, "第 1 個自訂意圖的名稱是必填"},
		{func(ci *CustomIntent) { ci.Description = "" }, "第 1 個自訂意圖的說明是必填"},
		{func(ci *CustomIntent) { ci.Examples = []string{"有賣禮盒嗎", " ", "想買來送人"} }, "「禮盒」至少要 3 句例句"},
		{func(ci *CustomIntent) {
			ci.Examples = []string{"有賣禮盒嗎", "有賣禮盒嗎？", "想買來送人"}
		}, "「禮盒」的例句重複：有賣禮盒嗎？"},
		{func(ci *CustomIntent) { ci.Examples[0] = "我要訂位" }, "「禮盒」的例句和內建案例重複：我要訂位"},
		{func(ci *CustomIntent) { ci.Action = "booking" }, "「禮盒」的動作不正確"},
		{func(ci *CustomIntent) { ci.Text = " " }, "「禮盒」要填回覆的文字"},
		{func(ci *CustomIntent) { ci.Action, ci.Knowledge = ActionKnowledge, nil }, "「禮盒」至少要指定一段知識"},
		{func(ci *CustomIntent) { ci.Action, ci.Knowledge = ActionKnowledge, []string{"停車"} }, "「禮盒」指定的知識段落不存在：停車"},
		{func(ci *CustomIntent) { ci.Name = "想預約" }, "自訂意圖的名稱重複：想預約"},
	} {
		cfg := customMerchant()
		cfg.Examples = []Example{{Text: "我要訂位", Intent: "booking", Builtin: true}}
		ci := giftBox()
		c.change(&ci)
		cfg.CustomIntents = []CustomIntent{ci}
		if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("got %v, want %s", err, c.want)
		}
	}
	cfg := customMerchant()
	for range 11 {
		ci := giftBox()
		ci.Name += strings.Repeat("！", len(cfg.CustomIntents))
		for i := range ci.Examples {
			ci.Examples[i] += strings.Repeat("a", len(cfg.CustomIntents)+1)
		}
		cfg.CustomIntents = append(cfg.CustomIntents, ci)
	}
	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "自訂意圖最多 10 個") {
		t.Errorf("11 個：%v", err)
	}
	cfg.CustomIntents = cfg.CustomIntents[:10]
	if err := cfg.Validate(); err != nil {
		t.Errorf("10 個：%v", err)
	}
	// 兩個自訂意圖的名稱相同、例句相同
	cfg.CustomIntents = []CustomIntent{giftBox(), giftBox()}
	err := cfg.Validate()
	if err == nil || !strings.Contains(err.Error(), "自訂意圖的名稱重複：禮盒") || !strings.Contains(err.Error(), "「禮盒」的例句重複：有賣禮盒嗎") {
		t.Errorf("兩個一樣的：%v", err)
	}
}
