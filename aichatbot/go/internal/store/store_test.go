package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func open(t *testing.T, path string) *DB {
	t.Helper()
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestReopenKeepsData(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "sub", "bot.sqlite")
	db := open(t, path)
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	if _, err := db.MarkUpdate(ctx, 42, now); err != nil {
		t.Fatal(err)
	}
	if err := db.AppendLog(ctx, LogEntry{Time: now, Platform: "telegram", UserID: "u1", ChatID: "c1", Direction: "customer", Text: "你好"}); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db = open(t, path)
	first, err := db.MarkUpdate(ctx, 42, now)
	if err != nil {
		t.Fatal(err)
	}
	if first {
		t.Error("重開後同一個 update 又被當成第一次")
	}
	logs, err := db.Logs(ctx, "telegram", "u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 1 || logs[0].Text != "你好" || !logs[0].Time.Equal(now) || string(logs[0].Traces) != "[]" || logs[0].ReplyKey != "" {
		t.Errorf("重開後的對話紀錄不對：%+v", logs)
	}
}

// 送不出去的回覆標成 SendFailed，讀出來還在；v10 的舊紀錄升級後是沒失敗
func TestLogSendFailed(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	all := migrations
	migrations = all[:10]
	db, err := Open(path)
	migrations = all
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	if _, err := db.SQL().ExecContext(ctx, `INSERT INTO conversation_log (time, platform, user_id, chat_id, direction, text, traces)
		VALUES (?, 'telegram', 'u1', 'c1', 'bot', '舊的', '[]')`, Timestamp(now)); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db = open(t, path)
	if err := db.AppendLog(ctx, LogEntry{Time: now, Platform: "telegram", UserID: "u1", ChatID: "c1", Direction: "bot", Text: "送不出去", SendFailed: true}); err != nil {
		t.Fatal(err)
	}
	logs, err := db.Logs(ctx, "telegram", "u1")
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[0].SendFailed || !logs[1].SendFailed {
		t.Errorf("Logs：%+v", logs)
	}
	around, err := db.LogsAround(ctx, "telegram", "u1", now, 5)
	if err != nil {
		t.Fatal(err)
	}
	recent, err := db.RecentLogs(ctx, "telegram", "u1", 5)
	if err != nil {
		t.Fatal(err)
	}
	one, err := db.LogEntry(ctx, logs[1].ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(around) != 2 || !around[1].SendFailed || len(recent) != 2 || !recent[1].SendFailed || !one.SendFailed {
		t.Errorf("LogsAround %+v、RecentLogs %+v、LogEntry %+v", around, recent, one)
	}
}

func TestMarkUpdate(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Now()
	for i, want := range []bool{true, false, false} {
		got, err := db.MarkUpdate(ctx, 7, now)
		if err != nil {
			t.Fatal(err)
		}
		if got != want {
			t.Errorf("第 %d 次：got %v, want %v", i+1, got, want)
		}
	}
	if got, _ := db.MarkUpdate(ctx, 8, now); !got {
		t.Error("不同的 update 應該是第一次")
	}
}

func TestCleanup(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	day := 24 * time.Hour
	for _, age := range []time.Duration{91 * day, 89 * day} {
		e := LogEntry{Time: now.Add(-age), Platform: "telegram", UserID: "u1", ChatID: "c1", Direction: "bot", Text: age.String()}
		if err := db.AppendLog(ctx, e); err != nil {
			t.Fatal(err)
		}
	}
	db.MarkUpdate(ctx, 1, now.Add(-8*day))
	db.MarkUpdate(ctx, 2, now.Add(-6*day))
	// 管理員 session：過期的刪掉，還沒過期的留著
	if _, err := db.SQL().ExecContext(ctx, `INSERT INTO admin_user (id, username, password_hash, created_at) VALUES (1, 'admin', 'x', ?)`, Timestamp(now)); err != nil {
		t.Fatal(err)
	}
	for _, s := range []struct{ hash, expires string }{{"old", Timestamp(now.Add(-time.Second))}, {"new", Timestamp(now.Add(time.Second))}} {
		if _, err := db.SQL().ExecContext(ctx, `INSERT INTO admin_session (token_hash, user_id, created_at, expires_at) VALUES (?, 1, ?, ?)`, s.hash, Timestamp(now), s.expires); err != nil {
			t.Fatal(err)
		}
	}

	logs, updates, sessions, err := db.Cleanup(ctx, now)
	if err != nil {
		t.Fatal(err)
	}
	if logs != 1 || updates != 1 || sessions != 1 {
		t.Errorf("刪掉 %d 筆紀錄、%d 筆 update、%d 個 session，應該各 1 筆", logs, updates, sessions)
	}
	var left string
	db.SQL().QueryRowContext(ctx, "SELECT group_concat(token_hash) FROM admin_session").Scan(&left)
	if left != "new" {
		t.Errorf("留下的 session：%q", left)
	}
	leftLogs, _ := db.Logs(ctx, "telegram", "u1")
	if len(leftLogs) != 1 || leftLogs[0].Text != (89*day).String() {
		t.Errorf("留下的紀錄不對：%+v", leftLogs)
	}
	// 被清掉的 update 會再被當成第一次，沒清掉的不會
	if first, _ := db.MarkUpdate(ctx, 1, now); !first {
		t.Error("8 天前的 update 應該已經清掉")
	}
	if first, _ := db.MarkUpdate(ctx, 2, now); first {
		t.Error("6 天前的 update 不該清掉")
	}
}

// AI 用量：每次呼叫一筆；CostBetween 加總 [from, to) 之間的費用（含失敗的呼叫）
func TestUsageCost(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	oct := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)
	nov := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
	for _, u := range []UsageRecord{
		{Time: oct.Add(-time.Second), Model: "claude-haiku-5-5", Request: "Intent", InputTokens: 100, OutputTokens: 10, CostUSD: 1, OK: true},
		{Time: oct, Model: "claude-haiku-5-5", Request: "Intent", InputTokens: 100, OutputTokens: 10, CostUSD: 0.25, OK: true},
		{Time: oct.Add(time.Hour), Model: "claude-haiku-5-5", Request: "Pick", InputTokens: 50, OutputTokens: 300, CostUSD: 0.5, OK: false},
		{Time: oct.Add(2 * time.Hour), Model: "afm", Request: "Intent", OK: true},
		{Time: nov, Model: "claude-haiku-5-5", Request: "Intent", CostUSD: 4, OK: true},
	} {
		if err := db.RecordUsage(ctx, u); err != nil {
			t.Fatal(err)
		}
	}
	cost, err := db.CostBetween(ctx, oct, nov)
	if err != nil {
		t.Fatal(err)
	}
	if cost != 0.75 {
		t.Errorf("十月的費用：got %v, want 0.75", cost)
	}
	if cost, _ := db.CostBetween(ctx, nov.AddDate(0, 1, 0), nov.AddDate(0, 2, 0)); cost != 0 {
		t.Errorf("沒有紀錄的月份：got %v", cost)
	}
}

// v6：意圖案例的舊代號（打招呼、道謝、道別、閒聊、問商家）合併成 general，商家新增的也一起改；新增 state 欄位（預設空字串）
func TestMigrateMergesGeneralIntents(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	all := migrations
	migrations = all[:5]
	db, err := Open(path)
	migrations = all
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range [][]any{{"你好", "greeting", 1}, {"謝啦", "thanks", 1}, {"掰", "goodbye", 1}, {"哈哈", "small_talk", 0},
		{"幾點開", "store_question", 0}, {"我要訂位", "booking", 1}} {
		if _, err := db.SQL().ExecContext(ctx, "INSERT INTO intent_example (text, intent, builtin) VALUES (?, ?, ?)", e...); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()

	db = open(t, path)
	rows, err := db.SQL().QueryContext(ctx, "SELECT text, intent, builtin, state FROM intent_example ORDER BY id")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var text, intent, state string
		var builtin bool
		if err := rows.Scan(&text, &intent, &builtin, &state); err != nil {
			t.Fatal(err)
		}
		got = append(got, fmt.Sprintf("%s=%s/%v/%q", text, intent, builtin, state))
	}
	want := `你好=general/true/"",謝啦=general/true/"",掰=general/true/"",哈哈=general/false/"",幾點開=general/false/"",我要訂位=booking/true/""`
	if strings.Join(got, ",") != want {
		t.Errorf("升級後：%s", strings.Join(got, ","))
	}
}

// 模型設定：還沒存過時 ok 是 false；存了之後重開資料庫還在，再存一次取代原本的
func TestModelSetting(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	db := open(t, path)
	if _, _, ok, err := db.ModelSetting(ctx); err != nil || ok {
		t.Fatalf("還沒存過：ok=%v err=%v", ok, err)
	}
	if err := db.SaveModelSetting(ctx, "haiku", 5); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveModelSetting(ctx, "afm", 2.5); err != nil {
		t.Fatal(err)
	}
	db.Close()
	db = open(t, path)
	m, budget, ok, err := db.ModelSetting(ctx)
	if err != nil || !ok || m != "afm" || budget != 2.5 {
		t.Errorf("重開後：%s %v ok=%v err=%v", m, budget, ok, err)
	}
}

// 已經在 v2 的資料庫升級後，原本的資料還在、可以記用量
func TestMigrateFromV2(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	all := migrations
	migrations = all[:2]
	db, err := Open(path)
	migrations = all
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.SQL().ExecContext(ctx, "INSERT INTO language (code, position, enabled) VALUES ('zh-TW', 0, 1)"); err != nil {
		t.Fatal(err)
	}
	db.Close()

	db = open(t, path)
	var version, langs int
	db.SQL().QueryRowContext(ctx, "PRAGMA user_version").Scan(&version)
	db.SQL().QueryRowContext(ctx, "SELECT COUNT(*) FROM language").Scan(&langs)
	if version != len(migrations) || langs != 1 {
		t.Errorf("升級後：版本 %d、語言 %d 筆", version, langs)
	}
	if err := db.RecordUsage(ctx, UsageRecord{Time: time.Now(), Model: "afm", Request: "Intent", OK: true}); err != nil {
		t.Fatal(err)
	}
}

// v9 拿掉功能開關的 off_text：有填的搬到固定回覆 booking_off／handoff_off 的慣用語言版本，那個語言已經有的不蓋掉
func TestMigrateMovesOffTextToFixedReplies(t *testing.T) {
	ctx := context.Background()
	path := filepath.Join(t.TempDir(), "bot.sqlite")
	all := migrations
	migrations = all[:8]
	db, err := Open(path)
	migrations = all
	if err != nil {
		t.Fatal(err)
	}
	for _, q := range []string{
		"INSERT INTO merchant (id, name, intro, contact, preferred_language) VALUES (1, '森日茶屋', '茶館', '02', 'ja')",
		"INSERT INTO feature_switch (name, enabled, off_text) VALUES ('booking', 0, ' 請來電預約 '), ('handoff', 0, '請來電 02-8888')",
		"INSERT INTO fixed_reply (key, lang, text) VALUES ('handoff_off', 'ja', '電話してください'), ('booking_off', 'zh', '中文版')",
	} {
		if _, err := db.SQL().ExecContext(ctx, q); err != nil {
			t.Fatal(err)
		}
	}
	db.Close()

	db = open(t, path)
	rows, err := db.SQL().QueryContext(ctx, "SELECT key, lang, text FROM fixed_reply ORDER BY key, lang")
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var got []string
	for rows.Next() {
		var key, lang, text string
		if err := rows.Scan(&key, &lang, &text); err != nil {
			t.Fatal(err)
		}
		got = append(got, key+"/"+lang+"="+text)
	}
	if want := "booking_off/ja=請來電預約,booking_off/zh=中文版,handoff_off/ja=電話してください"; strings.Join(got, ",") != want {
		t.Errorf("固定回覆：%s", strings.Join(got, ","))
	}
	var n int
	db.SQL().QueryRowContext(ctx, "SELECT COUNT(*) FROM pragma_table_info('feature_switch') WHERE name = 'off_text'").Scan(&n)
	if n != 0 {
		t.Error("feature_switch 還有 off_text")
	}
	var enabled bool
	if err := db.SQL().QueryRowContext(ctx, "SELECT enabled FROM feature_switch WHERE name = 'booking'").Scan(&enabled); err != nil || enabled {
		t.Errorf("功能開關：enabled=%v err=%v", enabled, err)
	}
}

// v12 把沒存過的轉人工設定（全部欄位都是舊預設值）改成自動結束 30 分鐘；商家改過別的欄位、刻意存 0 的不動
func TestHandoffAutoEndMigration(t *testing.T) {
	ctx := context.Background()
	for _, c := range []struct {
		streak int
		want   int
	}{{3, 30}, {5, 0}} {
		path := filepath.Join(t.TempDir(), "bot.sqlite")
		all := migrations
		migrations = all[:11]
		db, err := Open(path)
		migrations = all
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.SQL().ExecContext(ctx, `INSERT INTO handoff_setting (id, staff_chat_id, bind_code, on_no_answer, streak_limit, service_hours, auto_end_minutes, prompt_text)
			VALUES (1, NULL, 'CODE', 1, ?, '[]', 0, '')`, c.streak); err != nil {
			t.Fatal(err)
		}
		db.Close()

		db = open(t, path)
		if s, err := db.HandoffSetting(ctx); err != nil || s.AutoEndMinutes != c.want {
			t.Errorf("累計 %d 次：%+v %v", c.streak, s, err)
		}
	}
}

// 還沒存過轉人工設定時用預設值（自動結束 30 分鐘，#55），並產生一組綁定碼（之後讀到的是同一組）
func TestHandoffSettingDefaults(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	s, err := db.HandoffSetting(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if s.StaffChatID != 0 || !s.OnNoAnswer || s.StreakLimit != 3 || len(s.ServiceHours) != 0 || s.AutoEndMinutes != 30 || s.PromptText != "" {
		t.Errorf("預設值：%+v", s)
	}
	if len(s.BindCode) < 8 {
		t.Errorf("綁定碼太短：%q", s.BindCode)
	}
	again, err := db.HandoffSetting(ctx)
	if err != nil || again.BindCode != s.BindCode {
		t.Errorf("第二次讀到不同的綁定碼：%q、%q（%v）", s.BindCode, again.BindCode, err)
	}
}

// 存了轉人工設定之後讀回來一樣（綁定的群組和綁定碼不受影響）
func TestSaveHandoffSetting(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	before, _ := db.HandoffSetting(ctx)
	want := HandoffSetting{
		OnNoAnswer: false, StreakLimit: 2, AutoEndMinutes: 30, PromptText: "請稍候",
		ServiceHours: []ServiceHours{{Weekday: 1, Start: "09:00", End: "18:00"}},
	}
	if err := db.SaveHandoffSetting(ctx, want); err != nil {
		t.Fatal(err)
	}
	got, err := db.HandoffSetting(ctx)
	if err != nil {
		t.Fatal(err)
	}
	want.BindCode = before.BindCode
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("讀回來：%+v，應該 %+v", got, want)
	}
}

// 綁定碼：碼不對不綁；碼對就綁定那個群組（蓋掉之前的）並換一組新碼，舊碼不能再用
func TestBindStaffChat(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	s, _ := db.HandoffSetting(ctx)
	if ok, err := db.BindStaffChat(ctx, "WRONG", -100); err != nil || ok {
		t.Fatalf("錯的碼：ok=%v err=%v", ok, err)
	}
	if ok, err := db.BindStaffChat(ctx, s.BindCode, -100); err != nil || !ok {
		t.Fatalf("對的碼：ok=%v err=%v", ok, err)
	}
	after, _ := db.HandoffSetting(ctx)
	if after.StaffChatID != -100 || after.BindCode == s.BindCode || after.BindCode == "" {
		t.Errorf("綁定後：%+v", after)
	}
	if ok, _ := db.BindStaffChat(ctx, s.BindCode, -200); ok {
		t.Error("舊碼還能用")
	}
	if ok, _ := db.BindStaffChat(ctx, after.BindCode, -200); !ok {
		t.Error("新碼不能用")
	}
	if s, _ := db.HandoffSetting(ctx); s.StaffChatID != -200 {
		t.Errorf("沒有蓋掉之前的群組：%d", s.StaffChatID)
	}
}

// 重新產生綁定碼：換一組新碼，舊碼不能再綁、新碼可以；已經綁好的群組不變
func TestRegenerateBindCode(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	s, _ := db.HandoffSetting(ctx)
	if ok, _ := db.BindStaffChat(ctx, s.BindCode, -100); !ok {
		t.Fatal("綁定失敗")
	}
	old, _ := db.HandoffSetting(ctx)
	code, err := db.RegenerateBindCode(ctx)
	if err != nil {
		t.Fatal(err)
	}
	after, _ := db.HandoffSetting(ctx)
	if code == old.BindCode || len(code) != 10 || after.BindCode != code || after.StaffChatID != -100 {
		t.Errorf("重新產生後：%q %+v", code, after)
	}
	if ok, _ := db.BindStaffChat(ctx, old.BindCode, -200); ok {
		t.Error("舊碼還能用")
	}
	if ok, _ := db.BindStaffChat(ctx, code, -200); !ok {
		t.Error("新碼不能用")
	}
}

// 員工回覆群組裡的通知或轉送的原始訊息，都能找回那張待處理事項；找不到時 ok 是 false
func TestTicketByGroupMessage(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Now()
	id, _ := db.AddTicket(ctx, Ticket{Kind: "media", Platform: "telegram", UserID: "u1", ChatID: "c1", Text: "📎"}, now)
	if err := db.SetTicketGroupMessage(ctx, id, 70); err != nil {
		t.Fatal(err)
	}
	if err := db.SetTicketForwardedMessage(ctx, id, 71); err != nil {
		t.Fatal(err)
	}
	for _, msg := range []int64{70, 71} {
		got, ok, err := db.TicketByGroupMessage(ctx, msg)
		if err != nil || !ok || got.ID != id || got.UserID != "u1" || got.Kind != "media" {
			t.Errorf("訊息 %d：%+v ok=%v err=%v", msg, got, ok, err)
		}
	}
	if _, ok, err := db.TicketByGroupMessage(ctx, 99); err != nil || ok {
		t.Errorf("不是通知：ok=%v err=%v", ok, err)
	}
}

// 預約只有待確認的能確認或拒絕，已處理過的回 ErrAlreadyHandled；確認後那張預約申請的事項也結束
func TestSetBookingStatus(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Now()
	id, _ := db.AddTicket(ctx, Ticket{Kind: "booking", Platform: "telegram", UserID: "u1", ChatID: "c1", Text: "📅"}, now)
	if err := db.AddBooking(ctx, Booking{ID: id, Platform: "telegram", UserID: "u1", ChatID: "c1", Summary: "s"}, now); err != nil {
		t.Fatal(err)
	}
	if err := db.SetBookingStatus(ctx, id, "confirmed", "小明", now); err != nil {
		t.Fatal(err)
	}
	if b, _ := db.Booking(ctx, id); b.Status != "confirmed" {
		t.Errorf("狀態：%s", b.Status)
	}
	if open, _ := db.OpenTickets(ctx); len(open) != 0 {
		t.Errorf("預約申請的事項沒結束：%+v", open)
	}
	if err := db.SetBookingStatus(ctx, id, "rejected", "小華", now); !errors.Is(err, ErrAlreadyHandled) {
		t.Errorf("重複處理：%v", err)
	}
	if err := db.SetBookingStatus(ctx, 999, "rejected", "小華", now); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("沒有這筆：%v", err)
	}
}

// #67：預約的自訂欄位存在 extra，讀一筆和後台列表都讀得回來；沒有自訂欄位的是 nil
func TestBookingExtra(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Now()
	if err := db.AddBooking(ctx, Booking{ID: 1, Platform: "telegram", UserID: "u1", ChatID: "c1", Summary: "s",
		Extra: map[string]string{"service": "染髮", "designer": "小葵"}}, now); err != nil {
		t.Fatal(err)
	}
	if err := db.AddBooking(ctx, Booking{ID: 2, Platform: "telegram", UserID: "u1", ChatID: "c1", People: "4", Summary: "s"}, now); err != nil {
		t.Fatal(err)
	}
	if b, err := db.Booking(ctx, 1); err != nil || b.Extra["service"] != "染髮" || b.Extra["designer"] != "小葵" {
		t.Errorf("讀一筆：%+v %v", b, err)
	}
	rows, err := db.Bookings(ctx, BookingFilter{})
	if err != nil || len(rows) != 2 {
		t.Fatalf("列表：%+v %v", rows, err)
	}
	if rows[0].ID != 2 || rows[0].Extra != nil || rows[1].Extra["service"] != "染髮" {
		t.Errorf("列表：%+v", rows)
	}
}

// 最近的對話：舊的在前，只取最後 n 則，略過轉人工中沒回覆的空白紀錄和別的顧客
func TestRecentLogs(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Now()
	add := func(user, dir, text string) {
		if err := db.AppendLog(ctx, LogEntry{Time: now, Platform: "telegram", UserID: user, ChatID: user, Direction: dir, Text: text}); err != nil {
			t.Fatal(err)
		}
	}
	for _, s := range []string{"1", "2", "3", "4", "5"} {
		add("u1", "customer", s)
	}
	add("u1", "bot", "")
	add("u2", "customer", "別人")
	add("u1", "customer", "6")
	got, err := db.RecentLogs(ctx, "telegram", "u1", 3)
	if err != nil {
		t.Fatal(err)
	}
	var texts []string
	for _, e := range got {
		texts = append(texts, e.Text)
	}
	if strings.Join(texts, ",") != "4,5,6" {
		t.Errorf("最近對話：%v", texts)
	}
}
