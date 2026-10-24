package admin

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"aichatbot/internal/model"
)

// withSwitch 給 fixture 一個模型開關：haiku 是 nil 表示沒有設定 API key
func (f *fixture) withSwitch(haiku model.ChatModel) *model.Switch {
	sw := &model.Switch{Haiku: haiku, AFM: model.NewScripted(nil), Budget: 10, Ledger: f.db, Log: f.srv.Log, Now: f.srv.Now}
	f.srv.Model = sw
	return sw
}

// GET /api/model：目前的模型開關、每月上限，和有沒有設定 API key
func TestModelSettingsGet(t *testing.T) {
	f := newFixture(t)
	f.withSwitch(nil)
	cookie := f.login()
	w := f.do("GET", "/api/model", "", cookie)
	if w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), `{"model":"afm","monthly_budget_usd":10,"haiku_available":false}`) {
		t.Errorf("GET /api/model：%d %s", w.Code, w.Body)
	}
}

// PUT /api/model：有 API key 時可以切到 Haiku、改上限；存進資料庫，模型開關馬上換掉，再切回 AFM 也可以
func TestModelSettingsPut(t *testing.T) {
	f := newFixture(t)
	sw := f.withSwitch(model.NewScripted(nil))
	cookie := f.login()
	w := f.do("PUT", "/api/model", `{"model":"haiku","monthly_budget_usd":3.5}`, cookie)
	if w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), `{"model":"haiku","monthly_budget_usd":3.5,"haiku_available":true}`) {
		t.Fatalf("PUT /api/model：%d %s", w.Code, w.Body)
	}
	if m, b := sw.Settings(); m != model.ModelHaiku || b != 3.5 {
		t.Errorf("模型開關沒換掉：%s %v", m, b)
	}
	if m, b, ok, _ := f.db.ModelSetting(context.Background()); !ok || m != "haiku" || b != 3.5 {
		t.Errorf("資料庫：%s %v %v", m, b, ok)
	}
	w = f.do("PUT", "/api/model", `{"model":"afm","monthly_budget_usd":0}`, cookie)
	if w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), `{"model":"afm","monthly_budget_usd":0,"haiku_available":true}`) {
		t.Errorf("切回 AFM：%d %s", w.Code, w.Body)
	}
}

// 沒有 API key 時不能切到 Haiku；上限負數或不是數字、模型不認得、多出欄位都回 400，設定和資料庫都不變
func TestModelSettingsRejects(t *testing.T) {
	f := newFixture(t)
	sw := f.withSwitch(nil)
	cookie := f.login()
	for _, c := range []struct{ body, want string }{
		{`{"model":"haiku","monthly_budget_usd":5}`, "沒有設定 ANTHROPIC_API_KEY"},
		{`{"model":"afm","monthly_budget_usd":-1}`, "0 以上"},
		{`{"model":"afm","monthly_budget_usd":"5"}`, "格式不正確"},
		{`{"model":"gpt","monthly_budget_usd":5}`, "afm 或 haiku"},
		{`{"model":"afm","monthly_budget_usd":5,"api_key":"x"}`, "格式不正確"},
	} {
		w := f.do("PUT", "/api/model", c.body, cookie)
		if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("PUT %s：%d %s，want %s", c.body, w.Code, w.Body, c.want)
		}
	}
	if m, b := sw.Settings(); m != model.AFM || b != 10 {
		t.Errorf("設定不該變：%s %v", m, b)
	}
	if _, _, ok, _ := f.db.ModelSetting(context.Background()); ok {
		t.Error("不該存進資料庫")
	}
}

// 存不進資料庫時回 500，模型開關也不換（記憶體和資料庫不會不一致）
func TestModelSettingsSaveFails(t *testing.T) {
	f := newFixture(t)
	sw := f.withSwitch(model.NewScripted(nil))
	cookie := f.login()
	if _, err := f.db.SQL().Exec("DROP TABLE model_setting"); err != nil {
		t.Fatal(err)
	}
	if w := f.do("PUT", "/api/model", `{"model":"haiku","monthly_budget_usd":3}`, cookie); w.Code != http.StatusInternalServerError {
		t.Errorf("%d %s", w.Code, w.Body)
	}
	if m, b := sw.Settings(); m != model.AFM || b != 10 {
		t.Errorf("模型開關不該換：%s %v", m, b)
	}
}
