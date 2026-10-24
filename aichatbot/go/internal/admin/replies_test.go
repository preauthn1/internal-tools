package admin

import (
	"context"
	"maps"
	"net/http"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
)

// GET /api/fixed-replies：慣用語言、啟用的語言，和可以改的 4 則固定回覆（白話名稱、各語言的句子）
func TestFixedRepliesGet(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	w := f.do("GET", "/api/fixed-replies", "", cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("GET：%d %s", w.Code, w.Body)
	}
	for _, want := range []string{
		`"preferred_language":"zh"`, `"languages":["zh","en","ja"]`,
		`{"key":"greeting","labels":{"en-US":"Greeting","ja-JP":"あいさつ","zh-TW":"打招呼"},"texts":{"en":[],"ja":[],"zh":["您好～有什麼可以幫您的嗎？"]}}`,
	} {
		if !strings.Contains(w.Body.String(), want) {
			t.Errorf("少了 %s：%s", want, w.Body)
		}
	}
	if n := strings.Count(w.Body.String(), `"key"`); n != 5 {
		t.Errorf("應該有 5 則：%s", w.Body)
	}
}

// translator 是照順序回譯文的假模型
func translator(texts ...string) *model.Scripted {
	var list []map[string]model.Value
	for _, s := range texts {
		list = append(list, map[string]model.Value{"translation": model.Text(s)})
	}
	return model.NewScripted(map[string][]map[string]model.Value{"Translate": list})
}

// storedReplies 讀出資料庫裡某一則固定回覆：語言 → 存的文字
func (f *fixture) storedReplies(key string) map[string]string {
	f.t.Helper()
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if err != nil {
		f.t.Fatal(err)
	}
	out := map[string]string{}
	for _, r := range c.Replies {
		if r.Key == key {
			out[r.Lang] = r.Text
		}
	}
	return out
}

// 改了慣用語言版本：每一句用 AI 翻成其他啟用的語言（照送來的覆蓋掉），一起存進去，回傳翻好的內容給後台確認
func TestFixedRepliesTranslate(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	m := translator("One moment~", "Coming, {name}!", "少々お待ちください～", "{name}様、すぐ行きます！")
	f.srv.Chat = m
	w := f.do("PUT", "/api/fixed-replies/handoff", `{"texts":{"zh":["稍等一下喔～","馬上來，{name}！"],"en":["舊的"]}}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("PUT：%d %s", w.Code, w.Body)
	}
	if want := `"texts":{"en":["One moment~","Coming, {name}!"],"ja":["少々お待ちください～","{name}様、すぐ行きます！"],"zh":["稍等一下喔～","馬上來，{name}！"]}`; !strings.Contains(w.Body.String(), want) {
		t.Errorf("回傳：%s", w.Body)
	}
	want := map[string]string{"zh": "稍等一下喔～|馬上來，{name}！", "en": "One moment~|Coming, {name}!", "ja": "少々お待ちください～|{name}様、すぐ行きます！"}
	if got := f.storedReplies("handoff"); !maps.Equal(got, want) {
		t.Errorf("資料庫：%v", got)
	}
	if reqs := m.Requests(); len(reqs) != 4 || reqs[0].Text != "稍等一下喔～" || !strings.Contains(reqs[2].Instructions, "日本語") {
		t.Errorf("翻譯請求：%+v", reqs)
	}

	// 慣用語言沒改、只修正日文：不呼叫 AI（假模型沒準備答案，呼叫就會失敗），照送來的存
	f.srv.Chat = translator()
	w = f.do("PUT", "/api/fixed-replies/handoff", `{"texts":{"zh":["稍等一下喔～"," 馬上來，{name}！"],"en":["One moment~","Coming, {name}!"],"ja":["少々お待ちくださいませ","{name}様、ただいま参ります"]}}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("修正翻譯：%d %s", w.Code, w.Body)
	}
	want["ja"] = "少々お待ちくださいませ|{name}様、ただいま参ります"
	if got := f.storedReplies("handoff"); !maps.Equal(got, want) {
		t.Errorf("修正後的資料庫：%v", got)
	}
}

// 任何一句翻不出來（AI 失敗、變數被吃掉）：回 502、整次不存
func TestFixedRepliesTranslateFails(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	for _, m := range []*model.Scripted{translator("Hi {name}", "Hello"), translator("Hi {name}", "Hey {name}", "こんにちは")} {
		f.srv.Chat = m
		w := f.do("PUT", "/api/fixed-replies/greeting", `{"texts":{"zh":["{name}您好","嗨 {name}"]}}`, cookie)
		if w.Code != http.StatusBadGateway || !strings.Contains(w.Body.String(), "翻譯失敗") {
			t.Errorf("%d %s", w.Code, w.Body)
		}
	}
	if got := f.storedReplies("greeting"); len(got) != 0 {
		t.Errorf("不該存進去：%v", got)
	}
}

// 還沒有商家設定、不能改的固定回覆、沒過檢查的內容：不呼叫 AI 就擋下
func TestFixedRepliesRejects(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	f.srv.Chat = translator()
	if w := f.do("PUT", "/api/fixed-replies/handoff", `{"texts":{"zh":["好"]}}`, cookie); w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "請先填商家資訊") {
		t.Errorf("還沒有商家設定：%d %s", w.Code, w.Body)
	}
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	for _, c := range []struct {
		path, body string
		code       int
		want       string
	}{
		{"/api/fixed-replies/fallback", `{"texts":{"zh":["好"]}}`, http.StatusNotFound, "沒有這則可以改的固定回覆"},
		{"/api/fixed-replies/handoff", `{"texts":{"zh":["好", " "]}}`, http.StatusBadRequest, "zh 的第 2 句是空白"},
		{"/api/fixed-replies/handoff", `{"texts":{"zh":[]}}`, http.StatusBadRequest, "慣用語言（zh）至少要一句"},
		{"/api/fixed-replies/handoff", `{"texts":{"zh":["好"],"ko":["네"]}}`, http.StatusBadRequest, "語言 ko 沒有啟用"},
		{"/api/fixed-replies/handoff", `{"text":{"zh":["好"]}}`, http.StatusBadRequest, "格式不正確"},
	} {
		w := f.do("PUT", c.path, c.body, cookie)
		if w.Code != c.code || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("PUT %s %s：%d %s，want %d %s", c.path, c.body, w.Code, w.Body, c.code, c.want)
		}
	}
	if reqs := f.srv.Chat.(*model.Scripted).Requests(); len(reqs) != 0 {
		t.Errorf("不該呼叫 AI：%d 次", len(reqs))
	}
}
