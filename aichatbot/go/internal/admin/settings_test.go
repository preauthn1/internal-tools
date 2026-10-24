package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
)

// 商家設定、固定回覆、模型設定的 API 都要先登入
func TestSettingsRequireLogin(t *testing.T) {
	f := newFixture(t)
	for _, path := range []string{"/api/merchant", "/api/faq", "/api/languages", "/api/fixed-replies", "/api/fixed-replies/greeting", "/api/model"} {
		if w := f.do("GET", path, ""); w.Code != http.StatusUnauthorized {
			t.Errorf("GET %s：%d", path, w.Code)
		}
		if w := f.do("PUT", path, "{}"); w.Code != http.StatusUnauthorized {
			t.Errorf("PUT %s：%d", path, w.Code)
		}
	}
}

// 防 CSRF：PUT 不是 application/json 回 415，不會存進去
func TestSettingsRequireJSON(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	for _, path := range []string{"/api/merchant", "/api/faq", "/api/languages"} {
		req := httptest.NewRequest("PUT", path, strings.NewReader(infoBody))
		req.Header.Set("Content-Type", "text/plain")
		req.AddCookie(cookie)
		w := httptest.NewRecorder()
		f.h.ServeHTTP(w, req)
		if w.Code != http.StatusUnsupportedMediaType {
			t.Errorf("PUT %s：%d", path, w.Code)
		}
	}
	if _, err := merchant.Load(context.Background(), f.db.SQL()); err != merchant.ErrNotConfigured {
		t.Errorf("不該存進去：%v", err)
	}
}

// 還沒有商家設定時，GET 回空白的商家資訊、空的 FAQ、預設語言（慣用中文、顧客語言中英日），加上可勾選和量過的語言
func TestSettingsDefaults(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	for _, c := range []struct{ path, want string }{
		{"/api/merchant", `{"name":"","intro":"","contact":"","knowledge":[],"warnings":[]}`},
		{"/api/faq", `{"faq":[],"warnings":[]}`},
		{"/api/languages", `{"preferred_language":"zh","languages":[{"code":"zh","enabled":true},{"code":"en","enabled":true},{"code":"ja","enabled":true}],` +
			`"supported":["zh","en","ja","ko","th","vi","id","ms","fr","de","es"],"measured":["zh","en","ja"]}`},
	} {
		w := f.do("GET", c.path, "", cookie)
		if w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), c.want) {
			t.Errorf("GET %s：%d %s", c.path, w.Code, w.Body)
		}
	}
}

const infoBody = `{"name":"森日茶屋","intro":"茶館","contact":"02-1234-5678","knowledge":[{"title":"營業時間","body":"11 點開"}]}`

// 欄位驗證沒過、格式不對都回 400 和錯在哪裡，不會存進去
func TestSettingsValidation(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	for _, c := range []struct{ path, body, want string }{
		{"/api/merchant", `{"name":" ","intro":"茶館","contact":"02","knowledge":[]}`, "商家名稱是必填"},
		{"/api/merchant", `{"name":"x","intro":"x","contact":"x","knowledge":[{"title":"a","body":"1"},{"title":"a","body":"2"}]}`, "知識段落標題重複：a"},
		{"/api/merchant", `{"name":"x","intro":"x","contact":"x","knowledge":[{"title":"","body":"1"}]}`, "第 1 段知識的標題是必填"},
		{"/api/merchant", `{"nmae":"x"}`, "格式不正確"},
		{"/api/faq", `{"faq":[{"key":"a","question":"q","answer":"a"},{"key":"a","question":"q2","answer":"a2"}]}`, "FAQ 代號重複：a"},
		{"/api/faq", `{"faq":[{"key":"","question":"q","alt_questions":[],"answer":" "}]}`, "第 1 題 FAQ 的答案是必填"},
		{"/api/faq", `not json`, "格式不正確"},
		{"/api/languages", `{"preferred_language":"ja","languages":[{"code":"zh","enabled":true},{"code":"ja","enabled":false}]}`, "慣用語言 ja 不在啟用的語言裡"},
		{"/api/languages", `{"preferred_language":"zh","languages":[{"code":"zh","enabled":true},{"code":"xx","enabled":true}]}`, "不認得的語言代碼：xx"},
	} {
		w := f.do("PUT", c.path, c.body, cookie)
		if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("PUT %s %s：%d %s，want %s", c.path, c.body, w.Code, w.Body, c.want)
		}
	}
	w := f.do("GET", "/api/merchant", "", cookie)
	if !jsonEqual(t, w.Body.String(), strings.Replace(infoBody, "}]}", `}],"warnings":[]}`, 1)) {
		t.Errorf("沒過的不該存進去：%s", w.Body)
	}
}

// 存了馬上生效，不用重開服務：同一個 Server 上，介面語言跟著新的慣用語言變，機器人 Load 拿到的也是新的
func TestSettingsTakeEffectImmediately(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	w := f.do("PUT", "/api/languages", `{"preferred_language":"en","languages":[{"code":"zh","enabled":true},{"code":"en","enabled":true}]}`, cookie)
	if w.Code != http.StatusOK {
		t.Fatalf("存語言：%d %s", w.Code, w.Body)
	}
	if w := f.do("GET", "/api/auth/me", "", cookie); !strings.Contains(w.Body.String(), `"language":"en-US"`) {
		t.Errorf("介面語言沒跟著變：%s", w.Body)
	}
	if w := f.do("PUT", "/api/merchant", strings.Replace(infoBody, "森日茶屋", "森日茶屋 二號店", 1), cookie); w.Code != http.StatusOK {
		t.Fatalf("改名稱：%d %s", w.Code, w.Body)
	}
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	if c.Merchant.Name != "森日茶屋 二號店" || c.Merchant.PreferredLanguage != "en" || len(c.Languages) != 2 {
		t.Errorf("Load 拿到的不是新的：%+v %v", c.Merchant, c.Languages)
	}
}

// FAQ 超過門檻題數（#65）：存得進去（不擋），GET 和 PUT 的回應都帶警告，說明每則訊息會先挑最像的幾題；沒超過時是空陣列
func TestFAQWarnings(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	var faq []merchant.FAQ
	for i := range merchant.FAQFilterAbove + 1 {
		faq = append(faq, merchant.FAQ{Question: fmt.Sprintf("問題%d", i), Answer: "答案"})
	}
	body, _ := json.Marshal(map[string]any{"faq": faq})
	want := fmt.Sprintf(`[{"kind":"faq_filtered","count":%d,"top":%d}]`, merchant.FAQFilterAbove+1, merchant.FAQFilterTop)
	for _, w := range []*httptest.ResponseRecorder{f.do("PUT", "/api/faq", string(body), cookie), f.do("GET", "/api/faq", "", cookie)} {
		var got struct{ Warnings json.RawMessage }
		if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &got) != nil || !jsonEqual(t, string(got.Warnings), want) {
			t.Errorf("%d %s", w.Code, got.Warnings)
		}
	}
	if w := f.do("PUT", "/api/faq", `{"faq":[{"question":"q","answer":"a"}]}`, cookie); !strings.Contains(w.Body.String(), `"warnings":[]`) {
		t.Errorf("沒超過門檻：%s", w.Body)
	}
}

// 知識段落太長（#65）：存得進去（不擋），GET 和 PUT 的回應都帶警告；PUT 的 body 不收 warnings
func TestKnowledgeWarnings(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	info := example.invalid{Name: "店", Intro: "介紹", Contact: "02", Knowledge: []merchant.Section{{Title: "菜單", Body: strings.Repeat("字", merchant.KnowledgeSectionLimit)}}}
	body, _ := json.Marshal(info)
	want := fmt.Sprintf(`[{"kind":"knowledge_section_long","limit":%d,"targets":["菜單"]}]`, merchant.KnowledgeSectionLimit)
	for _, w := range []*httptest.ResponseRecorder{f.do("PUT", "/api/merchant", string(body), cookie), f.do("GET", "/api/merchant", "", cookie)} {
		var got struct{ Warnings json.RawMessage }
		if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &got) != nil || !jsonEqual(t, string(got.Warnings), want) {
			t.Errorf("%d %s", w.Code, got.Warnings)
		}
	}
	if w := f.do("PUT", "/api/merchant", strings.Replace(infoBody, "}]}", `}],"warnings":[]}`, 1), cookie); w.Code != http.StatusBadRequest {
		t.Errorf("PUT 帶 warnings：%d %s", w.Code, w.Body)
	}
}

// 從森日茶屋匯入後，用後台 API 改三個區塊再匯出，和直接改設定、Save 再匯出的結果逐字相同；
// 新增的 FAQ 代號留空，由程式給 faq-1
func TestSettingsEditMatchesDirectEdit(t *testing.T) {
	ctx := context.Background()
	orig, err := merchant.FromEvalDir("../../../eval")
	if err != nil {
		t.Fatal(err)
	}
	f := newFixture(t)
	if _, err := merchant.Save(ctx, f.db.SQL(), orig); err != nil {
		t.Fatal(err)
	}
	cookie := f.login()

	// 直接改設定
	want, _ := merchant.ParseJSON(mustExport(t, orig)) // 深拷貝
	want.Merchant.Intro = "改過的介紹"
	want.Knowledge = append([]merchant.Section{{Title: "停車", Body: "附近有停車場"}}, want.Knowledge[1:]...)
	want.FAQ = append(want.FAQ[1:], merchant.FAQ{Key: "faq-1", Question: "可以帶寵物嗎", AltQuestions: []string{"狗可以進去嗎"}, Answer: "可以"})
	want.Merchant.PreferredLanguage = "ja"
	want.Languages = []merchant.Language{{Code: "ja", Enabled: true}, {Code: "zh", Enabled: true}, {Code: "ko", Enabled: false}}

	// 用後台 API 改
	for _, c := range []struct {
		path string
		body any
	}{
		{"/api/merchant", merchant.InfoOf(want)},
		{"/api/faq", map[string]any{"faq": append(want.FAQ[:len(want.FAQ)-1:len(want.FAQ)-1],
			merchant.FAQ{Question: "可以帶寵物嗎", AltQuestions: []string{"狗可以進去嗎", " "}, Answer: "可以"})}},
		{"/api/languages", merchant.LanguagesOf(want)},
	} {
		body, _ := json.Marshal(c.body)
		if w := f.do("PUT", c.path, string(body), cookie); w.Code != http.StatusOK {
			t.Fatalf("PUT %s：%d %s", c.path, w.Code, w.Body)
		}
	}
	got, err := merchant.Load(ctx, f.db.SQL())
	if err != nil {
		t.Fatal(err)
	}

	direct := openDB(t)
	if _, err := merchant.Save(ctx, direct.SQL(), want); err != nil {
		t.Fatal(err)
	}
	wantLoaded, err := merchant.Load(ctx, direct.SQL())
	if err != nil {
		t.Fatal(err)
	}
	if a, b := mustExport(t, wantLoaded), mustExport(t, got); string(a) != string(b) {
		t.Errorf("後台編輯和直接改設定匯出不同\n直接改:\n%s\n後台:\n%s", a, b)
	}
}

func mustExport(t *testing.T, c merchant.Config) []byte {
	t.Helper()
	b, err := merchant.Export(c)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func jsonEqual(t *testing.T, a, b string) bool {
	t.Helper()
	var x, y any
	if err := json.Unmarshal([]byte(a), &x); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(b), &y); err != nil {
		t.Fatal(err)
	}
	ja, _ := json.Marshal(x)
	jb, _ := json.Marshal(y)
	return string(ja) == string(jb)
}
