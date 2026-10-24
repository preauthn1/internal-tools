package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
)

func customBody(list ...merchant.CustomIntent) string {
	b, _ := json.Marshal(map[string]any{"custom_intents": list})
	return string(b)
}

func giftBox() merchant.CustomIntent {
	return merchant.CustomIntent{
		Name: "禮盒", Description: "詢問禮盒或送禮", Examples: []string{"茶葉禮盒有哪些組合", "想挑一盒茶送長輩", "公司要大量訂購送禮"},
		Action: merchant.ActionFixedText, Text: "禮盒請洽門市。",
	}
}

type customIntentsResponse struct {
	CustomIntents   []merchant.CustomIntent `json:"custom_intents"`
	KnowledgeTitles []string                `json:"knowledge_titles"`
}

func (f *fixture) getCustomIntents(cookie *http.Cookie) customIntentsResponse {
	f.t.Helper()
	w := f.do("GET", "/api/custom-intents", "", cookie)
	if w.Code != http.StatusOK {
		f.t.Fatalf("讀自訂意圖：%d %s", w.Code, w.Body)
	}
	var out customIntentsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		f.t.Fatal(err)
	}
	return out
}

// 自訂意圖整份讀寫：附上知識段落標題給後台選；存好回存進去的內容（程式給的代號）；
// 沒通過檢查回 400，例句和題庫重複回 409，都不存
func TestCustomIntentsAPI(t *testing.T) {
	f, _ := newExampleFixture(t)
	cookie := f.login()
	got := f.getCustomIntents(cookie)
	if len(got.CustomIntents) != 0 || strings.Join(got.KnowledgeTitles, ",") != "營業時間,菜單與價格,付款,自取優惠,人數上限" {
		t.Errorf("一開始：%+v", got)
	}

	w := f.do("PUT", "/api/custom-intents", customBody(giftBox()), cookie)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"code":"custom_1"`) {
		t.Fatalf("存檔：%d %s", w.Code, w.Body)
	}
	if got := f.getCustomIntents(cookie).CustomIntents; len(got) != 1 || got[0].Name != "禮盒" {
		t.Errorf("存好後：%+v", got)
	}

	bad := giftBox()
	bad.Examples = bad.Examples[:2]
	if w := f.do("PUT", "/api/custom-intents", customBody(bad), cookie); w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), "至少要 3 句例句") {
		t.Errorf("例句不夠：%d %s", w.Code, w.Body)
	}
	overlap := giftBox()
	overlap.Examples[0] = firstSentences(t)["testset"]
	if w := f.do("PUT", "/api/custom-intents", customBody(overlap), cookie); w.Code != http.StatusConflict || !strings.Contains(w.Body.String(), "testset") {
		t.Errorf("和題庫重複：%d %s", w.Code, w.Body)
	}
	if got := f.getCustomIntents(cookie).CustomIntents; len(got) != 1 || got[0].Examples[0] != "茶葉禮盒有哪些組合" {
		t.Errorf("沒過的不該存：%+v", got)
	}
}

// 意圖清單接著列出自訂意圖（名稱就是商家取的，三種介面語言都一樣），案例庫可以選
func TestIntentsListCustomIntents(t *testing.T) {
	f, _ := newExampleFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/custom-intents", customBody(giftBox()), cookie); w.Code != http.StatusOK {
		t.Fatalf("存檔：%d %s", w.Code, w.Body)
	}
	var intents []merchant.Intent
	json.Unmarshal(f.do("GET", "/api/intents", "", cookie).Body.Bytes(), &intents)
	if len(intents) != len(merchant.Intents)+1 {
		t.Fatalf("意圖：%+v", intents)
	}
	if last := intents[len(intents)-1]; last.Code != "custom_1" || last.Labels["zh-TW"] != "禮盒" || last.Labels["ja-JP"] != "禮盒" {
		t.Errorf("自訂意圖：%+v", last)
	}
	if w := f.do("POST", "/api/intent-examples", exampleBody("年節送禮有推薦的嗎", "custom_1"), cookie); w.Code != http.StatusCreated {
		t.Errorf("新增自訂意圖的案例：%d %s", w.Code, w.Body)
	}
}

// fakeSelfTester：每次自我測試都等 release 放行；回報第幾次執行、ctx 有沒有被取消
type fakeSelfTester struct {
	release  chan struct{}
	runs     atomic.Int32
	canceled atomic.Int32
}

func (f *fakeSelfTester) SelfTest(ctx context.Context, progress func(done, total int)) (dialog.SelfTestReport, error) {
	n := f.runs.Add(1)
	progress(1, 2)
	select {
	case <-f.release:
	case <-ctx.Done():
		f.canceled.Add(1)
		return dialog.SelfTestReport{}, ctx.Err()
	}
	progress(2, 2)
	return dialog.SelfTestReport{Total: 2, Conflicts: []dialog.SelfTestItem{{Text: fmt.Sprintf("第 %d 次", n), Expected: "custom_1", Got: "general"}}, Failed: []dialog.SelfTestItem{}}, nil
}

type selfTestStatus struct {
	Status string
	Done   int
	Total  int
	Report *dialog.SelfTestReport
}

func (f *fixture) selfTestStatus(cookie *http.Cookie) selfTestStatus {
	f.t.Helper()
	w := f.do("GET", "/api/custom-intents/self-test", "", cookie)
	if w.Code != http.StatusOK {
		f.t.Fatalf("自我測試狀態：%d %s", w.Code, w.Body)
	}
	var s selfTestStatus
	if err := json.Unmarshal(w.Body.Bytes(), &s); err != nil {
		f.t.Fatal(err)
	}
	return s
}

// waitStatus 等到自我測試變成 status、做完 done 句（最多 2 秒）
func (f *fixture) waitStatus(cookie *http.Cookie, status string, done int) selfTestStatus {
	f.t.Helper()
	for range 200 {
		if s := f.selfTestStatus(cookie); s.Status == status && s.Done == done {
			return s
		}
		time.Sleep(10 * time.Millisecond)
	}
	f.t.Fatalf("等不到 %s %d：%+v", status, done, f.selfTestStatus(cookie))
	return selfTestStatus{}
}

// 分類自我測試：還沒跑過是 idle；存自訂意圖後在背景跑（running 附進度），跑完是 done 附報告；
// 可以手動重跑；新的一次會取消還沒跑完的那次（同時間只跑一組）
func TestSelfTestAPI(t *testing.T) {
	f, _ := newExampleFixture(t)
	tester := &fakeSelfTester{release: make(chan struct{})}
	f.srv.SelfTester = tester
	cookie := f.login()
	if s := f.selfTestStatus(cookie); s.Status != "idle" || s.Report != nil {
		t.Errorf("一開始：%+v", s)
	}

	if w := f.do("PUT", "/api/custom-intents", customBody(giftBox()), cookie); w.Code != http.StatusOK {
		t.Fatalf("存檔：%d %s", w.Code, w.Body)
	}
	if s := f.waitStatus(cookie, "running", 1); s.Total != 2 || s.Report != nil {
		t.Errorf("執行中：%+v", s)
	}
	// 還沒跑完又手動重跑：取消上一組
	if w := f.do("POST", "/api/custom-intents/self-test", "{}", cookie); w.Code != http.StatusAccepted {
		t.Fatalf("重跑：%d %s", w.Code, w.Body)
	}
	for range 200 {
		if tester.runs.Load() == 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	tester.release <- struct{}{}
	s := f.waitStatus(cookie, "done", 2)
	if tester.canceled.Load() != 1 || s.Report == nil || len(s.Report.Conflicts) != 1 || s.Report.Conflicts[0].Text != "第 2 次" {
		t.Errorf("跑完：%+v（取消 %d 次）", s, tester.canceled.Load())
	}
}

// 自我測試在背景呼叫本機模型（分類 request=Intent），開始和跑完都寫進服務的 log，
// 才看得出 log 裡的分類失敗是自我測試、不是顧客訊息（#57）
func TestSelfTestLogged(t *testing.T) {
	f, _ := newExampleFixture(t)
	var buf bytes.Buffer
	f.srv.Log = slog.New(slog.NewTextHandler(&buf, nil))
	tester := &fakeSelfTester{release: make(chan struct{})}
	f.srv.SelfTester = tester
	cookie := f.login()
	if w := f.do("POST", "/api/custom-intents/self-test", "{}", cookie); w.Code != http.StatusAccepted {
		t.Fatalf("開始自我測試：%d %s", w.Code, w.Body)
	}
	tester.release <- struct{}{}
	f.waitStatus(cookie, "done", 2)
	for _, want := range []string{"msg=分類自我測試開始", "msg=分類自我測試完成 total=2 conflicts=1 failed=0"} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("log 沒有 %q：\n%s", want, buf.String())
		}
	}
}
