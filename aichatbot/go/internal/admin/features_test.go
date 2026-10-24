package admin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
)

// 預約設定：要先登入；還沒有商家設定時是開著的，要先填商家資訊才能存；存了 GET 讀得回來、對話層讀到的開關跟著改
func TestBookingSetting(t *testing.T) {
	f := newFixture(t)
	for _, method := range []string{"GET", "PUT"} {
		if w := f.do(method, "/api/booking-setting", `{"enabled":false}`); w.Code != http.StatusUnauthorized {
			t.Errorf("%s 沒登入：%d", method, w.Code)
		}
	}
	cookie := f.login()
	if w := f.do("GET", "/api/booking-setting", "", cookie); w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), `{"enabled":true,"people_enabled":true,"fields":[]}`) {
		t.Errorf("預設：%d %s", w.Code, w.Body)
	}
	if w := f.do("PUT", "/api/booking-setting", `{"enabled":false}`, cookie); w.Code != http.StatusBadRequest || !jsonEqual(t, w.Body.String(), `{"error":"請先填商家資訊"}`) {
		t.Errorf("還沒有商家資訊：%d %s", w.Code, w.Body)
	}
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	for _, body := range []string{`{"enabled":false,"off_text":"請來電"}`, `{"enabled":"no"}`, `not json`} {
		if w := f.do("PUT", "/api/booking-setting", body, cookie); w.Code != http.StatusBadRequest {
			t.Errorf("%s：%d %s", body, w.Code, w.Body)
		}
	}
	if w := f.do("PUT", "/api/booking-setting", `{"enabled":false}`, cookie); w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), `{"enabled":false,"people_enabled":true,"fields":[]}`) {
		t.Fatalf("存：%d %s", w.Code, w.Body)
	}
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if err != nil || merchant.FeatureEnabled(c, "booking") || !merchant.FeatureEnabled(c, "handoff") {
		t.Errorf("存了之後：%+v %v", c.Features, err)
	}
	if w := f.do("GET", "/api/booking-setting", "", cookie); !jsonEqual(t, w.Body.String(), `{"enabled":false,"people_enabled":true,"fields":[]}`) {
		t.Errorf("讀回來：%s", w.Body)
	}
}

// #67：預約欄位跟著預約設定存：預設問人數、沒有自訂欄位；存了讀得回來、對話層讀得到；沒通過檢查回 400；
// 只送 enabled（舊的前端）時預約欄位不變
func TestBookingFieldsSetting(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("GET", "/api/booking-setting", "", cookie); !jsonEqual(t, w.Body.String(), `{"enabled":true,"people_enabled":true,"fields":[]}`) {
		t.Errorf("預設：%s", w.Body)
	}
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	body := `{"enabled":true,"people_enabled":false,"fields":[{"key":"service","name":"服務項目","description":"要預約的服務","required":true},` +
		`{"key":"designer","name":"指定設計師","description":"指定的設計師","required":false}]}`
	if w := f.do("PUT", "/api/booking-setting", body, cookie); w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), body) {
		t.Fatalf("存：%d %s", w.Code, w.Body)
	}
	c, err := merchant.Load(context.Background(), f.db.SQL())
	if b := merchant.EffectiveBookingFields(c); err != nil || b.PeopleEnabled || len(b.Custom) != 2 || b.Custom[0].Key != "service" || !b.Custom[0].Required {
		t.Errorf("存了之後：%+v %v", c.BookingFields, err)
	}
	for _, bad := range []struct{ body, want string }{
		{`{"enabled":true,"people_enabled":true,"fields":[{"key":"people","name":"人","description":"d","required":true}]}`, "內建欄位"},
		{`{"enabled":true,"people_enabled":true,"fields":[{"key":"a","name":"A","description":"d"},{"key":"a","name":"B","description":"d"}]}`, "代號重複"},
		{`{"enabled":true,"people_enabled":true,"fields":[{"key":"a","name":"","description":"d"}]}`, "名稱是必填"},
		{`{"enabled":true,"people_enabled":true,"fields":[` + strings.Repeat(`{"key":"a","name":"A","description":"d"},`, 5) + `{"key":"b","name":"B","description":"d"}]}`, "最多 5 個"},
	} {
		if w := f.do("PUT", "/api/booking-setting", bad.body, cookie); w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), bad.want) {
			t.Errorf("%s：%d %s", bad.body, w.Code, w.Body)
		}
	}
	if w := f.do("PUT", "/api/booking-setting", `{"enabled":false}`, cookie); w.Code != http.StatusOK ||
		!jsonEqual(t, w.Body.String(), strings.Replace(body, `"enabled":true`, `"enabled":false`, 1)) {
		t.Errorf("只送 enabled：%d %s", w.Code, w.Body)
	}
	if w := f.do("PUT", "/api/booking-setting", `{"enabled":true,"people_enabled":true,"fields":[]}`, cookie); w.Code != http.StatusOK {
		t.Fatalf("改回預設：%d %s", w.Code, w.Body)
	}
	if c, _ := merchant.Load(context.Background(), f.db.SQL()); c.BookingFields != nil {
		t.Errorf("改回預設後應該當成沒設定：%+v", c.BookingFields)
	}
}

// 轉人工設定：要先登入；還沒存過時是預設值（總開關開、找不到答案時轉、累計 3 次、全天、不自動結束），附上綁定碼、還沒綁定
func TestHandoffSettingDefaults(t *testing.T) {
	f := newFixture(t)
	for _, c := range []struct{ method, path string }{{"GET", "/api/handoff-setting"}, {"PUT", "/api/handoff-setting"}, {"POST", "/api/handoff-setting/bind-code"}} {
		if w := f.do(c.method, c.path, "{}"); w.Code != http.StatusUnauthorized {
			t.Errorf("%s %s 沒登入：%d", c.method, c.path, w.Code)
		}
	}
	cookie := f.login()
	w := f.do("GET", "/api/handoff-setting", "", cookie)
	s, _ := f.db.HandoffSetting(context.Background())
	want := `{"enabled":true,"on_no_answer":true,"streak_limit":3,"service_hours":[],"auto_end_minutes":30,"bind_code":"` + s.BindCode + `","bound":false}`
	if w.Code != http.StatusOK || len(s.BindCode) != 10 || !jsonEqual(t, w.Body.String(), want) {
		t.Errorf("預設：%d %s", w.Code, w.Body)
	}
}

const handoffBody = `{"enabled":false,"on_no_answer":false,"streak_limit":5,` +
	`"service_hours":[{"weekday":1,"start":"09:00","end":"12:00"},{"weekday":1,"start":"12:00","end":"18:30"},{"weekday":0,"start":"00:00","end":"23:59"}],` +
	`"auto_end_minutes":1440}`

// 重新產生綁定碼：回新的設定（碼換了），舊碼不能再綁定、新碼可以；不必先有商家設定
func TestRegenerateBindCode(t *testing.T) {
	ctx := context.Background()
	f := newFixture(t)
	cookie := f.login()
	old, _ := f.db.HandoffSetting(ctx)
	w := f.do("POST", "/api/handoff-setting/bind-code", "{}", cookie)
	s, _ := f.db.HandoffSetting(ctx)
	if w.Code != http.StatusOK || s.BindCode == old.BindCode || !strings.Contains(w.Body.String(), `"bind_code":"`+s.BindCode+`"`) {
		t.Fatalf("重新產生：%d %s", w.Code, w.Body)
	}
	if ok, _ := f.db.BindStaffChat(ctx, old.BindCode, -100); ok {
		t.Error("舊碼還能用")
	}
	if ok, _ := f.db.BindStaffChat(ctx, s.BindCode, -100); !ok {
		t.Error("新碼不能用")
	}
}

// 轉人工設定的每個欄位都檢查：不合格回 400 和錯在哪裡，什麼都不存
func TestHandoffSettingValidation(t *testing.T) {
	ctx := context.Background()
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	body := func(streak, autoEnd int, hours string) string {
		return fmt.Sprintf(`{"enabled":false,"on_no_answer":true,"streak_limit":%d,"service_hours":[%s],"auto_end_minutes":%d}`, streak, hours, autoEnd)
	}
	for _, c := range []struct{ body, want string }{
		{body(0, 0, ""), "累計次數要是 1～10"},
		{body(11, 0, ""), "累計次數要是 1～10"},
		{body(3, -1, ""), "自動結束要是 0～1440 分鐘"},
		{body(3, 1441, ""), "自動結束要是 0～1440 分鐘"},
		{body(3, 0, `{"weekday":7,"start":"09:00","end":"12:00"}`), "第 1 個服務時段的星期要是 0～6"},
		{body(3, 0, `{"weekday":-1,"start":"09:00","end":"12:00"}`), "第 1 個服務時段的星期要是 0～6"},
		{body(3, 0, `{"weekday":1,"start":"9:00","end":"12:00"}`), "第 1 個服務時段的開始時間要是 HH:MM"},
		{body(3, 0, `{"weekday":1,"start":"09:00","end":"25:00"}`), "第 1 個服務時段的結束時間要是 HH:MM"},
		{body(3, 0, `{"weekday":1,"start":"09:60","end":"12:00"}`), "第 1 個服務時段的開始時間要是 HH:MM"},
		{body(3, 0, `{"weekday":1,"start":"09:00","end":"09:00"}`), "第 1 個服務時段的開始和結束不能相同"},
		{body(3, 0, `{"weekday":1,"start":"09:00","end":"12:00"},{"weekday":1,"start":"11:00","end":"13:00"}`), "第 2 個服務時段和第 1 個重疊"},
		// 跨過午夜的時段延伸到隔天：週一 22:00～02:00 和週二 01:00～03:00 重疊；週六跨到週日也算
		{body(3, 0, `{"weekday":1,"start":"22:00","end":"02:00"},{"weekday":2,"start":"01:00","end":"03:00"}`), "第 2 個服務時段和第 1 個重疊"},
		{body(3, 0, `{"weekday":0,"start":"00:30","end":"02:00"},{"weekday":6,"start":"23:00","end":"01:00"}`), "第 2 個服務時段和第 1 個重疊"},
		{`{"enabled":true,"on_no_answer":true,"streak_limit":3,"service_hours":[],"auto_end_minutes":0,"bind_code":"AAAAAAAAAA"}`, "格式不正確"},
		{`{"enabled":true,"on_no_answer":true,"streak_limit":3,"service_hours":[],"auto_end_minutes":0,"prompt_text":"x"}`, "格式不正確"},
	} {
		w := f.do("PUT", "/api/handoff-setting", c.body, cookie)
		if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), c.want) {
			t.Errorf("%s：%d %s", c.body, w.Code, w.Body)
		}
	}
	s, _ := f.db.HandoffSetting(ctx)
	c, _ := merchant.Load(ctx, f.db.SQL())
	if s.StreakLimit != 3 || s.AutoEndMinutes != 30 || len(s.ServiceHours) != 0 || !merchant.FeatureEnabled(c, "handoff") {
		t.Errorf("不該存進去：%+v %+v", s, c.Features)
	}
}

// 跨過午夜的服務時段（開始晚於結束）可以存，緊接著的下一段不算重疊
func TestHandoffSettingOvernight(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	body := `{"enabled":true,"on_no_answer":true,"streak_limit":3,"auto_end_minutes":0,"service_hours":[` +
		`{"weekday":1,"start":"18:00","end":"02:00"},{"weekday":2,"start":"02:00","end":"05:00"},{"weekday":6,"start":"20:00","end":"01:00"}]}`
	if w := f.do("PUT", "/api/handoff-setting", body, cookie); w.Code != http.StatusOK {
		t.Fatalf("存跨日時段：%d %s", w.Code, w.Body)
	}
	s, _ := f.db.HandoffSetting(context.Background())
	if len(s.ServiceHours) != 3 || s.ServiceHours[0].Start != "18:00" || s.ServiceHours[0].End != "02:00" {
		t.Errorf("讀回來：%+v", s.ServiceHours)
	}
}

// 存轉人工設定：GET 讀得回來、總開關存到功能開關 handoff；不動提示文字（prompt_text）和綁定的群組
func TestPutHandoffSetting(t *testing.T) {
	ctx := context.Background()
	f := newFixture(t)
	cookie := f.login()
	if w := f.do("PUT", "/api/merchant", infoBody, cookie); w.Code != http.StatusOK {
		t.Fatalf("存商家資訊：%d %s", w.Code, w.Body)
	}
	s, _ := f.db.HandoffSetting(ctx)
	s.PromptText = "請稍等店員"
	if err := f.db.SaveHandoffSetting(ctx, s); err != nil {
		t.Fatal(err)
	}
	if ok, _ := f.db.BindStaffChat(ctx, s.BindCode, -100); !ok {
		t.Fatal("綁定失敗")
	}
	s, _ = f.db.HandoffSetting(ctx)
	w := f.do("PUT", "/api/handoff-setting", handoffBody, cookie)
	want := handoffBody[:len(handoffBody)-1] + `,"bind_code":"` + s.BindCode + `","bound":true}`
	if w.Code != http.StatusOK || !jsonEqual(t, w.Body.String(), want) {
		t.Fatalf("存：%d %s", w.Code, w.Body)
	}
	if w := f.do("GET", "/api/handoff-setting", "", cookie); !jsonEqual(t, w.Body.String(), want) {
		t.Errorf("讀回來：%s", w.Body)
	}
	after, _ := f.db.HandoffSetting(ctx)
	if after.PromptText != "請稍等店員" || after.StaffChatID != -100 {
		t.Errorf("提示文字或群組被改掉：%+v", after)
	}
	c, err := merchant.Load(ctx, f.db.SQL())
	if err != nil || merchant.FeatureEnabled(c, "handoff") || !merchant.FeatureEnabled(c, "booking") {
		t.Errorf("功能開關：%+v %v", c.Features, err)
	}
}
