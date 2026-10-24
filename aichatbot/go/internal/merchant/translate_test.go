package merchant

import (
	"context"
	"errors"
	"strings"
	"testing"

	"aichatbot/internal/model"
)

func translations(texts ...string) *model.Scripted {
	var list []map[string]model.Value
	for _, s := range texts {
		list = append(list, map[string]model.Value{"translation": model.Text(s)})
	}
	return model.NewScripted(map[string][]map[string]model.Value{"Translate": list})
}

// 翻譯成功：回譯文（去掉前後空白）；請求帶原文，指示裡寫清楚從什麼語言翻成什麼語言
func TestTranslate(t *testing.T) {
	m := translations(" Hello {name}~ \n")
	got, err := Translate(context.Background(), m, "zh", "en", "{name}您好～")
	if err != nil || got != "Hello {name}~" {
		t.Fatalf("got %q, %v", got, err)
	}
	r := m.Requests()[0]
	if r.Text != "{name}您好～" || !strings.Contains(r.Instructions, "中文") || !strings.Contains(r.Instructions, "English") {
		t.Errorf("請求：%+v", r)
	}
}

// 翻譯顧客訊息用自己的題目和指示（照原意翻、不加表情符號），不同於翻店家回覆的；
// 譯文含 | 或 {…} 照樣接受（不是固定回覆），AI 失敗或空白才算翻譯失敗
func TestTranslateCustomer(t *testing.T) {
	m := model.NewScripted(map[string][]map[string]model.Value{"TranslateCustomer": {
		{"translation": model.Text(" 有停車場嗎？ ")}, {"translation": model.Text("A|B {x}")}, {"translation": model.Text(" ")},
	}})
	got, err := TranslateCustomer(context.Background(), m, "ja", "zh", "駐車場はありますか？")
	if err != nil || got != "有停車場嗎？" {
		t.Fatalf("got %q, %v", got, err)
	}
	r := m.Requests()[0]
	if r.Text != "駐車場はありますか？" || !strings.Contains(r.Instructions, "從日本語翻成繁體中文") || !strings.Contains(r.Instructions, "表情符號") ||
		strings.Contains(r.Instructions, "店家客服的回覆") {
		t.Errorf("請求：%+v", r)
	}
	if got, err := TranslateCustomer(context.Background(), m, "en", "zh", "A|B {x}"); err != nil || got != "A|B {x}" {
		t.Errorf("含 | 和 {…}：got %q, %v", got, err)
	}
	for name, chat := range map[string]*model.Scripted{"空白": m, "AI 失敗": model.NewScripted(nil)} {
		if got, err := TranslateCustomer(context.Background(), chat, "en", "zh", "hi"); !errors.Is(err, ErrTranslate) {
			t.Errorf("%s：got %q, %v", name, got, err)
		}
	}
}

// AI 沒回答、回空白、{…} 變數被吃掉或多出來、譯文含 | 都算翻譯失敗
func TestTranslateFailures(t *testing.T) {
	for _, c := range []struct {
		name  string
		model *model.Scripted
	}{
		{"AI 失敗", model.NewScripted(nil)},
		{"空白", translations("  ")},
		{"變數被吃掉", translations("Hello~")},
		{"變數被翻譯", translations("Hello {名字}~")},
		{"多出變數", translations("Hello {name} {total}~")},
		{"含 |（存的時候是分隔多句用的）", translations("Hello {name}|Hi {name}")},
	} {
		if got, err := Translate(context.Background(), c.model, "zh", "en", "{name}您好～"); !errors.Is(err, ErrTranslate) {
			t.Errorf("%s：got %q, %v", c.name, got, err)
		}
	}
}
