package store

import (
	"context"
	"database/sql"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

// sameText 是測試用的「相同問題」判斷：文字完全相同
func sameText(text string) func(string, []float32) bool {
	return func(other string, _ []float32) bool { return other == text }
}

// 相同的問題併成一筆、次數加一，每次出現都記下來；不同的問題另開一筆
func TestRecordPendingMerges(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	t0 := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	record := func(text string, at time.Time) (int64, bool) {
		t.Helper()
		id, merged, err := db.RecordPending(ctx, PendingHit{Time: at, Platform: "telegram", UserID: "u1", ChatID: "c1", Text: text, Lang: "zh", Source: "unanswerable"},
			nil, sameText(text))
		if err != nil {
			t.Fatal(err)
		}
		return id, merged
	}
	a, merged := record("有停車場嗎", t0)
	if merged {
		t.Error("第一次就說併進去了")
	}
	b, merged := record("有停車場嗎", t0.Add(time.Hour))
	if !merged || b != a {
		t.Errorf("相同的問題沒有合併：%d %d %v", a, b, merged)
	}
	if c, merged := record("可以帶狗嗎", t0.Add(2*time.Hour)); merged || c == a {
		t.Errorf("不同的問題被合併了：%d", c)
	}
	list, err := db.PendingQuestions(ctx, "open")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 2 || list[0].ID != a || list[0].Count != 2 || list[1].Count != 1 || !list[0].LastSeen.Equal(t0.Add(time.Hour)) {
		t.Fatalf("列表（依次數排序）：%+v", list)
	}
	q, err := db.PendingQuestion(ctx, a)
	if err != nil {
		t.Fatal(err)
	}
	if q.Text != "有停車場嗎" || q.Lang != "zh" || len(q.Hits) != 2 || q.Hits[0].Source != "unanswerable" || q.Hits[1].UserID != "u1" {
		t.Errorf("詳細：%+v", q)
	}
}

// 已忽略的照樣併進去計次（不重新顯示）；已轉成 FAQ 的不再合併，之後又出現就開一筆新的。比對時拿得到存起來的向量
func TestRecordPendingByStatus(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	vec := []float32{0.6, 0.8}
	sameVector := func(_ string, v []float32) bool { return len(v) == 2 && v[0] == 0.6 && v[1] == 0.8 }
	record := func(text string) (int64, bool) {
		t.Helper()
		id, merged, err := db.RecordPending(ctx, PendingHit{Time: now, Platform: "telegram", UserID: "u1", Text: text, Source: "unclear"}, vec, sameVector)
		if err != nil {
			t.Fatal(err)
		}
		return id, merged
	}
	a, _ := record("停車場？")
	if err := db.TransitionPending(ctx, a, "open", "ignored", ""); err != nil {
		t.Fatal(err)
	}
	if b, merged := record("哪裡停車"); !merged || b != a {
		t.Fatalf("已忽略的沒有合併：%d %v", b, merged)
	}
	if q, _ := db.PendingQuestion(ctx, a); q.Status != "ignored" || q.Count != 2 {
		t.Errorf("已忽略的：%+v", q)
	}
	if list, _ := db.PendingQuestions(ctx, "open"); len(list) != 0 {
		t.Errorf("已忽略的出現在待處理：%+v", list)
	}
	if err := db.TransitionPending(ctx, a, "ignored", "resolved", "faq"); err != nil {
		t.Fatal(err)
	}
	if c, merged := record("停車場在哪"); merged || c == a {
		t.Errorf("已轉成 FAQ 的被合併了：%d", c)
	}
	if q, _ := db.PendingQuestion(ctx, a); q.Status != "resolved" || q.Resolution != "faq" {
		t.Errorf("已轉成 FAQ 的：%+v", q)
	}
}

// TransitionPending 只在狀態還是 from 時才改（同時按兩次只有一次成功）；轉成 FAQ 後員工再回覆也不再接到預填答案
func TestTransitionPending(t *testing.T) {
	ctx := context.Background()
	db := open(t, filepath.Join(t.TempDir(), "bot.sqlite"))
	now := time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)
	id, _, err := db.RecordPending(ctx, PendingHit{Time: now, Platform: "telegram", UserID: "u1", Text: "停車場？", Source: "unanswerable", TicketID: 7}, nil, sameText("停車場？"))
	if err != nil {
		t.Fatal(err)
	}
	if err := db.TransitionPending(ctx, id, "open", "resolved", "faq"); err != nil {
		t.Fatal(err)
	}
	if err := db.TransitionPending(ctx, id, "open", "resolved", "faq"); !errors.Is(err, ErrAlreadyHandled) {
		t.Errorf("第二次：%v", err)
	}
	if err := db.TransitionPending(ctx, 999, "open", "ignored", ""); !errors.Is(err, sql.ErrNoRows) {
		t.Errorf("沒有這筆：%v", err)
	}
	if err := db.AddPendingStaffAnswer(ctx, Ticket{ID: 7}, "對面有停車場"); err != nil {
		t.Fatal(err)
	}
	if q, _ := db.PendingQuestion(ctx, id); q.StaffAnswer != "" {
		t.Errorf("已處理的還接上員工回覆：%q", q.StaffAnswer)
	}
}
