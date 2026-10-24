package merchant

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"aichatbot/internal/store"
)

const evalDir = "../../../eval"

func openDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// importExport 匯入到新的資料庫再匯出
func importExport(t *testing.T, c Config) []byte {
	t.Helper()
	ctx := context.Background()
	db := openDB(t)
	if _, err := Save(ctx, db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	loaded, err := Load(ctx, db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	out, err := Export(loaded)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// 從 eval/ 匯入後匯出成 A，A 匯入新的資料庫再匯出成 B：A 和 B 逐字相同
func TestEvalRoundTrip(t *testing.T) {
	c, err := ReadFile(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	a := importExport(t, c)

	path := filepath.Join(t.TempDir(), "merchant.json")
	if err := os.WriteFile(path, a, 0o644); err != nil {
		t.Fatal(err)
	}
	fromA, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	b := importExport(t, fromA)
	if !bytes.Equal(a, b) {
		t.Errorf("再匯入匯出後內容不同\nA:\n%s\nB:\n%s", a, b)
	}
}

// 匯入的內容和原檔的筆數、key 一致
func TestEvalImportMatchesFiles(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}

	faqLines := 0
	data, _ := os.ReadFile(filepath.Join(evalDir, "faq.txt"))
	for _, l := range strings.Split(string(data), "\n") {
		if l = strings.TrimSpace(l); l != "" && !strings.HasPrefix(l, "#") {
			faqLines++
		}
	}
	if len(c.FAQ) != faqLines || faqLines == 0 {
		t.Errorf("FAQ %d 題，faq.txt 有 %d 行", len(c.FAQ), faqLines)
	}
	hours := c.FAQ[0]
	if hours.Key != "hours" || hours.Question != "營業時間" || !contains(hours.AltQuestions, "定休日") || !strings.Contains(hours.Answer, "週一公休") {
		t.Errorf("第一題 FAQ 不對：%+v", hours)
	}

	data, _ = os.ReadFile(filepath.Join(evalDir, "replies.json"))
	var all map[string]json.RawMessage
	if err := json.Unmarshal(data, &all); err != nil {
		t.Fatal(err)
	}
	delete(all, "說明")
	raw := map[string]map[string]string{}
	for k, v := range all {
		var m map[string]string
		if err := json.Unmarshal(v, &m); err != nil {
			t.Fatal(err)
		}
		if !usesStructuredVars(m) { // 見 TestEvalImportSkipsStructuredReplies
			raw[k] = m
		}
	}
	got := map[string]map[string]string{}
	for _, r := range c.Replies {
		if got[r.Key] == nil {
			got[r.Key] = map[string]string{}
		}
		got[r.Key][r.Lang] = r.Text
	}
	if !equalReplies(raw, got) {
		t.Errorf("固定回覆和 replies.json 不同：%v vs %v", keys(got), keys(raw))
	}

	examples := 0
	data, _ = os.ReadFile(filepath.Join(evalDir, "examples.jsonl"))
	for _, l := range strings.Split(string(data), "\n") {
		if strings.TrimSpace(l) != "" {
			examples++
		}
	}
	if len(c.Examples) != examples || examples == 0 {
		t.Errorf("意圖案例 %d 筆，examples.jsonl 有 %d 行", len(c.Examples), examples)
	}

	if c.Merchant.Name != "森日茶屋" || c.Merchant.Intro == "" || !strings.Contains(c.Merchant.Contact, "復興南路") {
		t.Errorf("商家資訊不對：%+v", c.Merchant)
	}
	var knowledge strings.Builder
	for _, s := range c.Knowledge {
		knowledge.WriteString(s.Title + "\n" + s.Body + "\n")
	}
	for _, want := range []string{"週一公休", "11:00–20:00", "日月潭紅玉紅茶 120 元", "滿 200 元", "滿 300 元打 9 折", "最多 20 位"} {
		if !strings.Contains(knowledge.String(), want) {
			t.Errorf("知識段落少了「%s」：\n%s", want, knowledge.String())
		}
	}
}

// 用到結構化規則變數（營業時間、人數上限、自取門檻）的固定回覆不匯入：ADR 0003 拿掉了這些欄位，對話層用內建的說法
func TestEvalImportSkipsStructuredReplies(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	keys := map[string]bool{}
	for _, r := range c.Replies {
		keys[r.Key] = true
		for _, v := range []string{"{open}", "{close}", "{max}", "{min}"} {
			if strings.Contains(r.Text, v) {
				t.Errorf("%s（%s）用到 %s", r.Key, r.Lang, v)
			}
		}
	}
	if keys["ask_time"] || keys["err_hours"] || keys["err_too_many"] || !keys["ask_date"] || !keys["bot_identity"] {
		t.Errorf("匯入的固定回覆：%v", keys)
	}
}

// 關掉資料庫再開，商家設定還在
func TestReopenKeepsConfig(t *testing.T) {
	ctx := context.Background()
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := Load(ctx, db.SQL()); err != ErrNotConfigured {
		t.Errorf("空的資料庫應該回 ErrNotConfigured，got %v", err)
	}
	if _, err := Save(ctx, db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	before, _ := Load(ctx, db.SQL())
	db.Close()

	db, err = store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	after, err := Load(ctx, db.SQL())
	if err != nil {
		t.Fatal(err)
	}
	a, _ := Export(before)
	b, _ := Export(after)
	if !bytes.Equal(a, b) {
		t.Error("重開後商家設定不同")
	}
}

func TestValidate(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	c.Merchant.Contact = " "
	c.FAQ = append(c.FAQ, c.FAQ[0])
	err = c.Validate()
	if err == nil || !strings.Contains(err.Error(), "聯絡方式") || !strings.Contains(err.Error(), "FAQ 代號重複") {
		t.Errorf("應該擋下沒填聯絡方式和重複的 FAQ：%v", err)
	}
	if _, err := Save(context.Background(), openDB(t).SQL(), c); err == nil {
		t.Error("沒通過檢查的設定不該存進去")
	}
}

func TestParseJSONRejectsUnknownFields(t *testing.T) {
	if _, err := ParseJSON([]byte(`{"merchant": {"nmae": "x"}}`)); err == nil {
		t.Error("打錯的欄位名稱應該報錯")
	}
}

// 舊的備份（#37 以前）功能開關有 off_text：匯入時忽略，其他欄位照讀；匯出不再有 off_text
func TestParseJSONIgnoresOldOffText(t *testing.T) {
	c, err := ParseJSON([]byte(`{"features": [{"name": "booking", "enabled": false, "off_text": "請來電"}]}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Features) != 1 || c.Features[0] != (Feature{Name: "booking", Enabled: false}) {
		t.Errorf("功能開關：%+v", c.Features)
	}
	out, err := Export(c)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "off_text") {
		t.Errorf("匯出還有 off_text：%s", out)
	}
	if _, err := ParseJSON([]byte(`{"features": [{"name": "booking", "enabeld": false}]}`)); err == nil {
		t.Error("功能開關打錯的欄位名稱應該報錯")
	}
}

// #67：預約欄位（人數開關、自訂欄位）匯入、存進資料庫、再匯出不變；沒設定的舊檔匯出也不會多出 booking_fields
func TestBookingFieldsRoundTrip(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	if out := importExport(t, c); strings.Contains(string(out), "booking_fields") {
		t.Errorf("沒設定預約欄位不該匯出 booking_fields")
	}
	c.BookingFields = &BookingFields{PeopleEnabled: false, Custom: []BookingField{
		{Key: "service", Name: "服務項目", Description: "要預約的服務", Required: true},
		{Key: "designer", Name: "指定設計師", Description: "指定的設計師名字"},
	}}
	out := importExport(t, c)
	back, err := ParseJSON(out)
	if err != nil {
		t.Fatal(err)
	}
	if b := back.BookingFields; b == nil || b.PeopleEnabled || len(b.Custom) != 2 || b.Custom[0] != c.BookingFields.Custom[0] || b.Custom[1] != c.BookingFields.Custom[1] {
		t.Errorf("預約欄位：%+v", back.BookingFields)
	}
}

// #67：自訂欄位最多 5 個；代號是小寫英文開頭的英數底線、不能重複、不能是內建欄位；名稱、說明必填
func TestValidateBookingFields(t *testing.T) {
	c, err := FromEvalDir(evalDir)
	if err != nil {
		t.Fatal(err)
	}
	f := func(key string) BookingField { return BookingField{Key: key, Name: key, Description: "說明"} }
	for _, x := range []struct {
		custom []BookingField
		want   string
	}{
		{[]BookingField{f("a"), f("b"), f("c"), f("d"), f("e"), f("g")}, "自訂欄位最多 5 個"},
		{[]BookingField{f("Service")}, "代號不正確：Service"},
		{[]BookingField{f("1st")}, "代號不正確：1st"},
		{[]BookingField{f("service"), f("service")}, "代號重複：service"},
		{[]BookingField{f("people")}, "代號不能用內建欄位：people"},
		{[]BookingField{f("name_soft")}, "代號不能用內建欄位：name_soft"}, // ask_name_soft 是內建的固定回覆
		{[]BookingField{{Key: "service", Name: " ", Description: "說明"}}, "第 1 個自訂欄位的名稱是必填"},
		{[]BookingField{{Key: "service", Name: "服務", Description: ""}}, "第 1 個自訂欄位的說明是必填"},
	} {
		c.BookingFields = &BookingFields{PeopleEnabled: true, Custom: x.custom}
		if err := c.Validate(); err == nil || !strings.Contains(err.Error(), x.want) {
			t.Errorf("%+v：%v，應有「%s」", x.custom, err, x.want)
		}
	}
	c.BookingFields = &BookingFields{Custom: []BookingField{f("a"), f("b"), f("c"), f("d"), f("service_2")}}
	if err := c.Validate(); err != nil {
		t.Errorf("5 個應該可以：%v", err)
	}
}

// 匯入檔的 booking_fields 沒寫 people_enabled 時照預設問人數（不會默默關掉）；打錯的欄位名稱照樣報錯
func TestParseBookingFieldsDefaultsPeople(t *testing.T) {
	c, err := ParseJSON([]byte(`{"booking_fields": {"custom": [{"key": "service", "name": "服務", "description": "d", "required": true}]}}`))
	if err != nil {
		t.Fatal(err)
	}
	if b := c.BookingFields; b == nil || !b.PeopleEnabled || len(b.Custom) != 1 {
		t.Errorf("預約欄位：%+v", c.BookingFields)
	}
	if _, err := ParseJSON([]byte(`{"booking_fields": {"people_enable": false}}`)); err == nil {
		t.Error("打錯的欄位名稱應該報錯")
	}
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

func equalReplies(a, b map[string]map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, langs := range a {
		if len(langs) != len(b[k]) {
			return false
		}
		for l, text := range langs {
			if b[k][l] != text {
				return false
			}
		}
	}
	return true
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
