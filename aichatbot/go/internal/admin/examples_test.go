package admin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
)

const evalDir = "../../../eval"

type example struct {
	ID      int64
	Text    string
	Intent  string
	Builtin bool
}

// newExampleFixture：匯入森日茶屋（內建案例來自 examples.jsonl），題庫是 eval/
func newExampleFixture(t *testing.T) (*fixture, merchant.Config) {
	t.Helper()
	f := newFixture(t)
	c, err := merchant.FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := merchant.Save(context.Background(), f.db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	f.srv.EvalDir = evalDir
	return f, c
}

// listExamples 用 API 列出案例
func (f *fixture) listExamples(cookie *http.Cookie) []example {
	f.t.Helper()
	w := f.do("GET", "/api/intent-examples", "", cookie)
	if w.Code != http.StatusOK {
		f.t.Fatalf("列出案例：%d %s", w.Code, w.Body)
	}
	var out []example
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		f.t.Fatal(err)
	}
	return out
}

func exampleBody(text, intent string) string {
	b, _ := json.Marshal(map[string]string{"text": text, "intent": intent})
	return string(b)
}

// 商家新增的案例：存成非內建、列得出來；可以刪除，刪過再刪是 404
func TestAddAndDeleteExample(t *testing.T) {
	f, c := newExampleFixture(t)
	cookie := f.login()
	const text = "你們的貓叫什麼名字"
	tests, err := merchant.LoadTestSet(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	if sets := tests.Find(text); len(sets) > 0 {
		t.Fatalf("測試句子不能在題庫裡：%v", sets)
	}
	w := f.do("POST", "/api/intent-examples", exampleBody(" "+text+" ", "general"), cookie)
	if w.Code != http.StatusCreated {
		t.Fatalf("新增：%d %s", w.Code, w.Body)
	}
	var added example
	json.Unmarshal(w.Body.Bytes(), &added)
	if added.ID == 0 || added.Text != text || added.Intent != "general" || added.Builtin {
		t.Errorf("新增的案例：%+v", added)
	}
	list := f.listExamples(cookie)
	if len(list) != len(c.Examples)+1 || list[len(list)-1] != added {
		t.Errorf("列表最後一筆：%+v，want %+v", list[len(list)-1], added)
	}
	path := fmt.Sprintf("/api/intent-examples/%d", added.ID)
	if w := f.do("DELETE", path, "{}", cookie); w.Code != http.StatusNoContent {
		t.Fatalf("刪除：%d %s", w.Code, w.Body)
	}
	if n := len(f.listExamples(cookie)); n != len(c.Examples) {
		t.Errorf("刪除後 %d 筆，want %d", n, len(c.Examples))
	}
	if w := f.do("DELETE", path, "{}", cookie); w.Code != http.StatusNotFound {
		t.Errorf("再刪一次：%d %s", w.Code, w.Body)
	}
}

// 內建案例不能刪（403）、不能改（沒有修改的 API，PUT、PATCH 回 405）；之後資料不變
func TestBuiltinExampleReadOnly(t *testing.T) {
	f, _ := newExampleFixture(t)
	cookie := f.login()
	before := f.listExamples(cookie)
	builtin := before[0]
	path := fmt.Sprintf("/api/intent-examples/%d", builtin.ID)
	if w := f.do("DELETE", path, "{}", cookie); w.Code != http.StatusForbidden || !strings.Contains(w.Body.String(), "內建案例不能刪除") {
		t.Errorf("刪除內建案例：%d %s", w.Code, w.Body)
	}
	for _, method := range []string{"PUT", "PATCH"} {
		if w := f.do(method, path, exampleBody("改掉", "general"), cookie); w.Code != http.StatusMethodNotAllowed {
			t.Errorf("%s 內建案例：%d %s", method, w.Code, w.Body)
		}
	}
	if after := f.listExamples(cookie); !slices.Equal(before, after) {
		t.Error("內建案例被改到了")
	}
}

// firstSentences 從題庫各取第一句真實句子（驗證題組 holdout/ 不能打開來看，不取；faq_set 跳過 holdout 題）
func firstSentences(t *testing.T) map[string]string {
	t.Helper()
	read := func(name string, v any) {
		data, err := os.ReadFile(filepath.Join(evalDir, name))
		if err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal(data, v); err != nil {
			t.Fatal(err)
		}
	}
	var items struct{ Items []struct{ Text string } }
	read("intent_set.json", &items)
	var testset struct {
		Cases []struct {
			Q string `json:"q"`
		}
	}
	read("testset.json", &testset)
	var faqSet struct {
		Cases []struct {
			Q       string `json:"q"`
			Holdout bool
		}
	}
	read("faq_set.json", &faqSet)
	faq := ""
	for _, c := range faqSet.Cases {
		if !c.Holdout {
			faq = c.Q
			break
		}
	}
	var flows struct {
		Conversations []struct{ Turns []struct{ User string } }
	}
	read("flow_tests.json", &flows)
	return map[string]string{
		"intent_set": items.Items[0].Text,
		"testset":    testset.Cases[0].Q,
		"faq_set":    faq,
		"flow_tests": flows.Conversations[0].Turns[0].User,
	}
}

// 和題庫重複的句子存不進去（409，訊息寫出是哪個題庫）；只差標點、大小寫也算重複
func TestExampleOverlapsTestSet(t *testing.T) {
	f, c := newExampleFixture(t)
	cookie := f.login()
	for set, text := range firstSentences(t) {
		for _, s := range []string{text, " " + strings.ToUpper(text) + "！？", strings.ToLower(text) + "..."} {
			w := f.do("POST", "/api/intent-examples", exampleBody(s, "general"), cookie)
			if w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), set) {
				t.Errorf("%s 的 %q：%d %s", set, s, w.Code, w.Body)
			}
		}
	}
	if n := len(f.listExamples(cookie)); n != len(c.Examples) {
		t.Errorf("和題庫重複的句子存進去了：%d 筆，want %d", n, len(c.Examples))
	}
}

// 和案例庫裡已有的句子（內建或商家新增的）比對後相同也存不進去（409）
func TestExampleDuplicate(t *testing.T) {
	f, c := newExampleFixture(t)
	cookie := f.login()
	builtin := c.Examples[0].Text
	if w := f.do("POST", "/api/intent-examples", exampleBody(builtin+"！", "general"), cookie); w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "案例庫裡已經有") {
		t.Errorf("和內建案例重複：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", "/api/intent-examples", exampleBody("你們的貓叫什麼名字", "general"), cookie); w.Code != http.StatusCreated {
		t.Fatalf("新增：%d %s", w.Code, w.Body)
	}
	if w := f.do("POST", "/api/intent-examples", exampleBody("你們的貓叫什麼名字？", "general"), cookie); w.Code != http.StatusConflict {
		t.Errorf("和自己新增的重複：%d %s", w.Code, w.Body)
	}
	if n := len(f.listExamples(cookie)); n != len(c.Examples)+1 {
		t.Errorf("案例 %d 筆，want %d", n, len(c.Examples)+1)
	}
}

// 句子去掉頭尾空白後不能是空的、要有字母或數字、最多 200 字；意圖要在清單裡；不符合回 400，不存
func TestExampleInvalid(t *testing.T) {
	f, c := newExampleFixture(t)
	cookie := f.login()
	for _, body := range []string{
		exampleBody("你們的貓叫什麼名字", "no_such_intent"),
		exampleBody("你們的貓叫什麼名字", ""),
		exampleBody("   ", "general"),
		exampleBody("！？…", "general"),
		exampleBody(strings.Repeat("貓", 201), "general"),
		`{"text": 1}`,
	} {
		if w := f.do("POST", "/api/intent-examples", body, cookie); w.Code != http.StatusBadRequest {
			t.Errorf("%.60s：%d %s", body, w.Code, w.Body)
		}
	}
	if n := len(f.listExamples(cookie)); n != len(c.Examples) {
		t.Errorf("不合格的案例存進去了：%d 筆", n)
	}
	if w := f.do("POST", "/api/intent-examples", exampleBody(strings.Repeat("貓", 200), "general"), cookie); w.Code != http.StatusCreated {
		t.Errorf("剛好 200 字：%d %s", w.Code, w.Body)
	}
}

// 題庫讀不到：拒絕存檔（500），不能默默當作沒有重複
func TestExampleTestSetMissing(t *testing.T) {
	f, c := newExampleFixture(t)
	cookie := f.login()
	f.srv.EvalDir = filepath.Join(t.TempDir(), "missing")
	if w := f.do("POST", "/api/intent-examples", exampleBody("你們的貓叫什麼名字", "general"), cookie); w.Code != http.StatusInternalServerError || !strings.Contains(w.Body.String(), "題庫讀不到") {
		t.Errorf("%d %s", w.Code, w.Body)
	}
	if n := len(f.listExamples(cookie)); n != len(c.Examples) {
		t.Errorf("題庫讀不到還存進去了：%d 筆", n)
	}
}

// 新增、刪除也要先登入
func TestExampleWritesRequireLogin(t *testing.T) {
	f, c := newExampleFixture(t)
	if w := f.do("POST", "/api/intent-examples", exampleBody("你們的貓叫什麼名字", "general")); w.Code != http.StatusUnauthorized {
		t.Errorf("新增：%d", w.Code)
	}
	if w := f.do("DELETE", "/api/intent-examples/1", "{}"); w.Code != http.StatusUnauthorized {
		t.Errorf("刪除：%d", w.Code)
	}
	if n := len(f.listExamples(f.login())); n != len(c.Examples) {
		t.Errorf("沒登入還改到了：%d 筆", n)
	}
}

// 列出案例：要先登入；匯入的案例都是內建的，每筆有 id
func TestListExamples(t *testing.T) {
	f, c := newExampleFixture(t)
	if w := f.do("GET", "/api/intent-examples", ""); w.Code != http.StatusUnauthorized {
		t.Errorf("沒登入：%d", w.Code)
	}
	got := f.listExamples(f.login())
	if len(got) != len(c.Examples) {
		t.Fatalf("案例 %d 筆，want %d", len(got), len(c.Examples))
	}
	for i, e := range got {
		if e.ID == 0 || !e.Builtin || e.Text != c.Examples[i].Text || e.Intent != c.Examples[i].Intent {
			t.Errorf("第 %d 筆：%+v，want %+v", i, e, c.Examples[i])
		}
	}
}

// 意圖清單：要先登入；每個意圖都有繁中、英、日的白話名稱
func TestIntentsAPI(t *testing.T) {
	f := newFixture(t)
	if w := f.do("GET", "/api/intents", ""); w.Code != http.StatusUnauthorized {
		t.Errorf("沒登入：%d", w.Code)
	}
	w := f.do("GET", "/api/intents", "", f.login())
	if w.Code != http.StatusOK {
		t.Fatalf("%d %s", w.Code, w.Body)
	}
	var intents []struct {
		Code   string
		Labels map[string]string
	}
	if err := json.Unmarshal(w.Body.Bytes(), &intents); err != nil {
		t.Fatal(err)
	}
	if len(intents) != len(merchant.Intents) || len(intents) == 0 {
		t.Errorf("意圖 %d 個，want %d", len(intents), len(merchant.Intents))
	}
	for _, in := range intents {
		for _, lang := range []string{"zh-TW", "en-US", "ja-JP"} {
			if in.Labels[lang] == "" {
				t.Errorf("%s 沒有 %s 的名稱", in.Code, lang)
			}
		}
	}
}
