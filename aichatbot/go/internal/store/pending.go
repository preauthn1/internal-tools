package store

import (
	"context"
	"database/sql"
	"encoding/binary"
	"errors"
	"math"
	"time"
)

// PendingHit 是待補問題的一次出現（v10 的 pending_question_hit）。Source：unanswerable、verify_failed、unclear
type PendingHit struct {
	Time     time.Time `json:"time"`
	Platform string    `json:"platform"`
	UserID   string    `json:"user_id"`
	ChatID   string    `json:"chat_id"`
	Text     string    `json:"text"`
	Lang     string    `json:"lang"`
	Source   string    `json:"source"`
	TicketID int64     `json:"ticket_id,omitempty"` // 這次轉人工開的事項；0 是沒有轉人工
}

// PendingQuestion 是待補問題（v10 的 pending_question）：Text、Lang 是第一次出現的原話和語言，Hits 只有 PendingQuestion 會讀
type PendingQuestion struct {
	ID          int64        `json:"id"`
	Text        string       `json:"text"`
	Lang        string       `json:"lang"`
	Translated  string       `json:"translated"`
	Count       int          `json:"count"`
	Status      string       `json:"status"`
	Resolution  string       `json:"resolution"`
	StaffAnswer string       `json:"staff_answer"`
	FirstSeen   time.Time    `json:"first_seen"`
	LastSeen    time.Time    `json:"last_seen"`
	Hits        []PendingHit `json:"hits,omitempty"`
}

// RecordPending 記一次答不出來的訊息：和還沒轉成 FAQ、知識、案例的待補問題（待處理、已忽略）逐一用 same（原話、向量）比對，
// 相同就併進去（次數加一，狀態不變），都不同就新開一筆。vector 是這次訊息的語意向量（沒有是 nil）。回傳待補問題的編號和有沒有合併
func (d *DB) RecordPending(ctx context.Context, h PendingHit, vector []float32, same func(text string, vector []float32) bool) (id int64, merged bool, err error) {
	tx, err := d.sql.BeginTx(ctx, nil)
	if err != nil {
		return 0, false, err
	}
	defer tx.Rollback()
	id, err = findPending(ctx, tx, same)
	if err != nil {
		return 0, false, err
	}
	now := Timestamp(h.Time)
	if id != 0 {
		merged = true
		_, err = tx.ExecContext(ctx, "UPDATE pending_question SET count = count + 1, last_seen = ? WHERE id = ?", now, id)
	} else {
		var res sql.Result
		res, err = tx.ExecContext(ctx, `
			INSERT INTO pending_question (text, lang, translated, vector, count, status, resolution, staff_answer, first_seen, last_seen)
			VALUES (?, ?, '', ?, 1, 'open', '', '', ?, ?)`, h.Text, h.Lang, encodeVector(vector), now, now)
		if err == nil {
			id, err = res.LastInsertId()
		}
	}
	if err != nil {
		return 0, false, err
	}
	var ticket any
	if h.TicketID != 0 {
		ticket = h.TicketID
	}
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO pending_question_hit (question_id, time, platform, user_id, chat_id, text, lang, source, ticket_id)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, id, now, h.Platform, h.UserID, h.ChatID, h.Text, h.Lang, h.Source, ticket); err != nil {
		return 0, false, err
	}
	return id, merged, tx.Commit()
}

// findPending 找第一筆（最早開的）和這次訊息相同的待處理、已忽略的待補問題；沒有是 0
func findPending(ctx context.Context, tx *sql.Tx, same func(string, []float32) bool) (int64, error) {
	rows, err := tx.QueryContext(ctx, "SELECT id, text, vector FROM pending_question WHERE status IN ('open', 'ignored') ORDER BY id")
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var text string
		var blob []byte
		if err := rows.Scan(&id, &text, &blob); err != nil {
			return 0, err
		}
		if same(text, decodeVector(blob)) {
			return id, nil
		}
	}
	return 0, rows.Err()
}

const pendingColumns = "id, text, lang, translated, count, status, resolution, staff_answer, first_seen, last_seen"

func scanPending(row interface{ Scan(...any) error }) (PendingQuestion, error) {
	var q PendingQuestion
	var first, last string
	err := row.Scan(&q.ID, &q.Text, &q.Lang, &q.Translated, &q.Count, &q.Status, &q.Resolution, &q.StaffAnswer, &first, &last)
	if err != nil {
		return q, err
	}
	if q.FirstSeen, err = parseTimestamp(first); err != nil {
		return q, err
	}
	q.LastSeen, err = parseTimestamp(last)
	return q, err
}

// PendingQuestions 列出某個狀態的待補問題，次數多的在前（次數相同時最近出現的在前）
func (d *DB) PendingQuestions(ctx context.Context, status string) ([]PendingQuestion, error) {
	rows, err := d.sql.QueryContext(ctx, "SELECT "+pendingColumns+" FROM pending_question WHERE status = ? ORDER BY count DESC, last_seen DESC, id DESC", status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []PendingQuestion{}
	for rows.Next() {
		q, err := scanPending(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, q)
	}
	return out, rows.Err()
}

// PendingQuestion 讀一筆待補問題和每次出現（依時間先後）；沒有這筆回 sql.ErrNoRows
func (d *DB) PendingQuestion(ctx context.Context, id int64) (PendingQuestion, error) {
	q, err := scanPending(d.sql.QueryRowContext(ctx, "SELECT "+pendingColumns+" FROM pending_question WHERE id = ?", id))
	if err != nil {
		return q, err
	}
	rows, err := d.sql.QueryContext(ctx, `
		SELECT time, platform, user_id, chat_id, text, lang, source, COALESCE(ticket_id, 0)
		FROM pending_question_hit WHERE question_id = ? ORDER BY time, id`, id)
	if err != nil {
		return q, err
	}
	defer rows.Close()
	q.Hits = []PendingHit{}
	for rows.Next() {
		var h PendingHit
		var t string
		if err := rows.Scan(&t, &h.Platform, &h.UserID, &h.ChatID, &h.Text, &h.Lang, &h.Source, &h.TicketID); err != nil {
			return q, err
		}
		if h.Time, err = parseTimestamp(t); err != nil {
			return q, err
		}
		q.Hits = append(q.Hits, h)
	}
	return q, rows.Err()
}

// LogsAround 是某位顧客在 at 前後的對話紀錄（舊的在前）：at 和之前的最多 n 則、之後的最多 n 則
func (d *DB) LogsAround(ctx context.Context, platform, userID string, at time.Time, n int) ([]LogEntry, error) {
	return d.queryLogs(ctx, `
		SELECT id, time, platform, user_id, chat_id, direction, text, reply_key, traces, send_failed FROM (
			SELECT * FROM (
				SELECT time, platform, user_id, chat_id, direction, text, COALESCE(reply_key, '') AS reply_key, traces, send_failed, id
				FROM conversation_log WHERE platform = ? AND user_id = ? AND time <= ? ORDER BY id DESC LIMIT ?)
			UNION ALL
			SELECT * FROM (
				SELECT time, platform, user_id, chat_id, direction, text, COALESCE(reply_key, '') AS reply_key, traces, send_failed, id
				FROM conversation_log WHERE platform = ? AND user_id = ? AND time > ? ORDER BY id LIMIT ?)
		) ORDER BY id`, platform, userID, Timestamp(at), n, platform, userID, Timestamp(at), n)
}

// TransitionPending 在待補問題的狀態還是 from 時改成 to（open、ignored、resolved）；resolution 是轉成什麼（faq、knowledge、example），
// 不是 resolved 時傳空字串。狀態已經不是 from 回 ErrAlreadyHandled（同時按兩次只有一次成功），沒有這筆回 sql.ErrNoRows
func (d *DB) TransitionPending(ctx context.Context, id int64, from, to, resolution string) error {
	err := execOne(d.sql.ExecContext(ctx, "UPDATE pending_question SET status = ?, resolution = ? WHERE id = ? AND status = ?", to, resolution, id, from))
	if errors.Is(err, sql.ErrNoRows) {
		// 沒改到：有這筆就是狀態變了，沒有這筆回 sql.ErrNoRows
		if err := d.sql.QueryRowContext(ctx, "SELECT 1 FROM pending_question WHERE id = ?", id).Scan(new(int)); err != nil {
			return err
		}
		return ErrAlreadyHandled
	}
	return err
}

// SetPendingTranslation 存待補問題翻成商家慣用語言的版本；沒有這筆回 sql.ErrNoRows
func (d *DB) SetPendingTranslation(ctx context.Context, id int64, translated string) error {
	return execOne(d.sql.ExecContext(ctx, "UPDATE pending_question SET translated = ? WHERE id = ?", translated, id))
}

// AddPendingStaffAnswer：員工回覆了待處理事項 ticket，有待補問題的那次出現轉人工開的是這張（或緊急通知對應的原本那張）時，
// 把員工的原話接在那筆待補問題的員工回覆後面（換行分隔）；已經轉成 FAQ、知識、案例的不再接
func (d *DB) AddPendingStaffAnswer(ctx context.Context, ticket Ticket, text string) error {
	_, err := d.sql.ExecContext(ctx, `
		UPDATE pending_question SET staff_answer = CASE WHEN staff_answer = '' THEN ? ELSE staff_answer || char(10) || ? END
		WHERE status != 'resolved' AND id IN (SELECT question_id FROM pending_question_hit WHERE ticket_id IN (?, ?))`, text, text, ticket.ID, ticket.RelatedTo)
	return err
}

// execOne：沒改到任何一列時回 sql.ErrNoRows
func execOne(res sql.Result, err error) error {
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err == nil && n == 0 {
		err = sql.ErrNoRows
	}
	return err
}

// encodeVector 把向量存成 little-endian float32；nil 存成 NULL
func encodeVector(v []float32) any {
	if v == nil {
		return nil
	}
	b := make([]byte, 4*len(v))
	for i, x := range v {
		binary.LittleEndian.PutUint32(b[4*i:], math.Float32bits(x))
	}
	return b
}

func decodeVector(b []byte) []float32 {
	if b == nil {
		return nil
	}
	v := make([]float32, len(b)/4)
	for i := range v {
		v[i] = math.Float32frombits(binary.LittleEndian.Uint32(b[4*i:]))
	}
	return v
}
