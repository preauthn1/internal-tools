// Package store 是 SQLite 資料庫：顧客狀態、已處理的 Telegram update、對話紀錄、待處理事項、預約、AI 用量和模型設定、
// 員工群組和轉人工設定（handoff.go）、待補問題（pending.go），
// 商家設定（表在這裡建，讀寫在 internal/merchant），以及後台的管理員帳號 admin_user、登入 session admin_session
// （表在這裡建，讀寫在 internal/admin）。
// 資料表照搬 Swift 版（tag swift-final 的 Sources/Storage/Database.swift），欄位改成 snake_case。
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "example.invalid/sqlite" // 純 Go 的 SQLite，不用 cgo
)

// 對話紀錄保留 90 天、update 紀錄保留 7 天（Cleanup 刪掉更舊的）
const (
	LogRetention    = 90 * 24 * time.Hour
	UpdateRetention = 7 * 24 * time.Hour
)

type DB struct {
	sql *sql.DB
}

// Open 開啟（沒有就建立）資料庫檔案，上層目錄不存在就建立；開檔後跑 migration
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	dsn := "file:" + path + "?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)"
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}
	// 一次只用一條連線：寫入不會互搶鎖，伺服器的量也不需要更多
	db.SetMaxOpenConns(1)
	d := &DB{sql: db}
	if err := d.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return d, nil
}

func (d *DB) Close() error { return d.sql.Close() }

// SQL 給同一個資料庫裡其他表的讀寫用（商家設定）
func (d *DB) SQL() *sql.DB { return d.sql }

// migrations[i] 把 user_version 從 i 升到 i+1；只能往後加，不能改已經有的
var migrations = []string{
	// v1：照搬 Swift 版的 v1～v3
	`
	-- 顧客狀態：整個存成 JSON
	CREATE TABLE customer_state (
		platform   TEXT NOT NULL,
		user_id    TEXT NOT NULL,
		data       TEXT NOT NULL,
		updated_at TEXT NOT NULL,
		PRIMARY KEY (platform, user_id)
	);
	CREATE TABLE processed_update (
		update_id  INTEGER PRIMARY KEY,
		created_at TEXT NOT NULL
	);
	-- 對話紀錄：direction 是 customer、bot、staff；traces 是 AI 判斷的 JSON 陣列（顧客、員工那筆是空陣列）
	CREATE TABLE conversation_log (
		id        INTEGER PRIMARY KEY AUTOINCREMENT,
		time      TEXT NOT NULL,
		platform  TEXT NOT NULL,
		user_id   TEXT NOT NULL,
		chat_id   TEXT NOT NULL,
		direction TEXT NOT NULL,
		text      TEXT NOT NULL,
		reply_key TEXT,
		traces    TEXT NOT NULL
	);
	CREATE INDEX conversation_log_time ON conversation_log (time);
	CREATE INDEX conversation_log_customer ON conversation_log (platform, user_id, time);
	-- 待處理事項：group_message_id 是員工群組裡那則通知的編號（員工「回覆」通知時用它找回來）
	CREATE TABLE ticket (
		id               INTEGER PRIMARY KEY AUTOINCREMENT,
		kind             TEXT NOT NULL,
		platform         TEXT NOT NULL,
		user_id          TEXT NOT NULL,
		chat_id          TEXT NOT NULL,
		original         TEXT NOT NULL,
		text             TEXT NOT NULL,
		status           TEXT NOT NULL,   -- open、replied
		related_to       INTEGER,         -- 緊急通知：對應的原本那張
		group_message_id INTEGER,
		created_at       TEXT NOT NULL
	);
	CREATE INDEX ticket_group_message ON ticket (group_message_id);
	-- 預約：編號和開出來的待處理事項相同；handled_by 是確認或拒絕的員工
	CREATE TABLE booking (
		id         INTEGER PRIMARY KEY,
		platform   TEXT NOT NULL,
		user_id    TEXT NOT NULL,
		chat_id    TEXT NOT NULL,
		date       TEXT,
		time       TEXT,
		people     TEXT,
		name       TEXT,
		phone      TEXT,
		summary    TEXT NOT NULL,
		status     TEXT NOT NULL,         -- pending、confirmed、rejected
		handled_by TEXT,
		handled_at TEXT,
		created_at TEXT NOT NULL
	);
	`,
	// v2：商家設定（一個部署只有一個商家，見 ADR 0003）；position 是後台顯示的順序
	`
	CREATE TABLE merchant (
		id                 INTEGER PRIMARY KEY CHECK (id = 1),
		name               TEXT NOT NULL,
		intro              TEXT NOT NULL,
		contact            TEXT NOT NULL,
		preferred_language TEXT NOT NULL
	);
	CREATE TABLE knowledge_section (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		position INTEGER NOT NULL,
		title    TEXT NOT NULL,
		body     TEXT NOT NULL
	);
	-- alt_questions：同一題的其他問法（JSON 字串陣列）
	CREATE TABLE faq (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		position      INTEGER NOT NULL,
		key           TEXT NOT NULL UNIQUE,
		question      TEXT NOT NULL,
		alt_questions TEXT NOT NULL,
		answer        TEXT NOT NULL
	);
	-- builtin：內建案例（隨程式提供，商家不能改）
	CREATE TABLE intent_example (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		text    TEXT NOT NULL,
		intent  TEXT NOT NULL,
		builtin INTEGER NOT NULL
	);
	-- 功能開關：booking、handoff；off_text 是關閉時回的文字（v9 拿掉，改用固定回覆 booking_off、handoff_off）
	CREATE TABLE feature_switch (
		name     TEXT PRIMARY KEY,
		enabled  INTEGER NOT NULL,
		off_text TEXT NOT NULL
	);
	-- 固定回覆：lang 是 * 時所有語言共用；text 可以用 | 分隔多句、含 {name} 這類變數
	CREATE TABLE fixed_reply (
		key  TEXT NOT NULL,
		lang TEXT NOT NULL,
		text TEXT NOT NULL,
		PRIMARY KEY (key, lang)
	);
	-- 顧客語言
	CREATE TABLE language (
		code     TEXT PRIMARY KEY,
		position INTEGER NOT NULL,
		enabled  INTEGER NOT NULL
	);
	`,
	// v3：AI 用量（每次呼叫一筆）：request 是題目名稱（不存顧客的話）；本機模型的 token 和費用記 0；ok 是有沒有讀到回答
	`
	CREATE TABLE ai_usage (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		time          TEXT NOT NULL,
		model         TEXT NOT NULL,
		request       TEXT NOT NULL,
		input_tokens  INTEGER NOT NULL,
		output_tokens INTEGER NOT NULL,
		cost_usd      REAL NOT NULL,
		ok            INTEGER NOT NULL
	);
	CREATE INDEX ai_usage_time ON ai_usage (time);
	`,
	// v4：後台管理員（只有一種角色）和登入 session；密碼存 bcrypt，session 只存 token 的 SHA-256（hex）
	`
	CREATE TABLE admin_user (
		id            INTEGER PRIMARY KEY AUTOINCREMENT,
		username      TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL,
		created_at    TEXT NOT NULL
	);
	CREATE TABLE admin_session (
		token_hash TEXT PRIMARY KEY,
		user_id    INTEGER NOT NULL REFERENCES admin_user (id) ON DELETE CASCADE,
		created_at TEXT NOT NULL,
		expires_at TEXT NOT NULL
	);
	CREATE INDEX admin_session_expires ON admin_session (expires_at);
	`,
	// v5：模型設定（後台存的）：model 是 afm 或 haiku，monthly_budget 是 Haiku 每月費用上限（美元，0 表示不限）。
	// 沒有這一列時用環境變數的初始值
	`
	CREATE TABLE model_setting (
		id             INTEGER PRIMARY KEY CHECK (id = 1),
		model          TEXT NOT NULL,
		monthly_budget REAL NOT NULL
	);
	`,
	// v6：內建意圖改成 11 個（#26）：打招呼、道謝、道別、閒聊、問商家合併成一般對話（general），商家新增的案例也一起改；
	// 意圖案例加上分類當下的對話狀態 state（ticket、booking:欄位；空字串是沒有特別的狀態）
	`
	UPDATE intent_example SET intent = 'general' WHERE intent IN ('greeting', 'thanks', 'goodbye', 'small_talk', 'store_question');
	ALTER TABLE intent_example ADD COLUMN state TEXT NOT NULL DEFAULT '';
	`,
	// v7：員工群組和轉人工設定（#28）。handoff_setting 只有一列：staff_chat_id 是用綁定碼綁定的員工群組（NULL 是還沒綁定），
	// service_hours 是員工服務時段的 JSON 陣列（空陣列是全天），auto_end_minutes 是員工最後一次回覆後幾分鐘恢復 AI（0 是不自動結束）。
	// 轉人工總開關用 v2 的 feature_switch（handoff）。
	// ticket.forwarded_message_id：顧客傳的圖片、語音、檔案轉送（copyMessage）到員工群組的那則，員工回覆它也能找回這張事項
	`
	CREATE TABLE handoff_setting (
		id               INTEGER PRIMARY KEY CHECK (id = 1),
		staff_chat_id    INTEGER,
		bind_code        TEXT NOT NULL,
		on_no_answer     INTEGER NOT NULL,
		streak_limit     INTEGER NOT NULL,
		service_hours    TEXT NOT NULL,
		auto_end_minutes INTEGER NOT NULL,
		prompt_text      TEXT NOT NULL
	);
	ALTER TABLE ticket ADD COLUMN forwarded_message_id INTEGER;
	CREATE INDEX ticket_forwarded_message ON ticket (forwarded_message_id);
	`,
	// v8：自訂意圖（#31）：code 是程式給的 custom_N；examples、knowledge（指定的知識段落標題）是 JSON 字串陣列；
	// text 是回固定文字的文字或轉人工的提示
	`
	CREATE TABLE custom_intent (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		position    INTEGER NOT NULL,
		code        TEXT NOT NULL UNIQUE,
		name        TEXT NOT NULL,
		description TEXT NOT NULL,
		examples    TEXT NOT NULL,
		action      TEXT NOT NULL CHECK (action IN ('fixed_text', 'handoff', 'knowledge')),
		text        TEXT NOT NULL,
		knowledge   TEXT NOT NULL
	);
	`,
	// v9：功能開關拿掉 off_text（#37）：關閉時改回固定回覆 booking_off、handoff_off。
	// 有填的 off_text 搬到那則固定回覆的慣用語言版本（那個語言已經有的不蓋掉；還沒有商家那一列就不搬）
	`
	INSERT OR IGNORE INTO fixed_reply (key, lang, text)
	SELECT f.name || '_off', m.preferred_language, TRIM(f.off_text)
	FROM feature_switch f, merchant m
	WHERE f.name IN ('booking', 'handoff') AND TRIM(f.off_text) != '';
	ALTER TABLE feature_switch DROP COLUMN off_text;
	`,
	// v10：待補問題（#30）：機器人答不出來的訊息，相同問題合併計次。text、lang 是第一次出現的原話和語言，
	// translated 是翻成商家慣用語言的版本（還沒翻是空字串），vector 是 text 的語意向量（little-endian float32；沒有是 NULL），
	// staff_answer 是轉人工後員工在群組回覆的原話（預填答案用）。
	// status：open（待處理）、ignored（忽略，照樣計次）、resolved（已轉成 resolution：faq、knowledge、example）。
	// pending_question_hit 是每次出現：source 是 unanswerable、verify_failed、unclear；ticket_id 是這次轉人工開的事項
	`
	CREATE TABLE pending_question (
		id           INTEGER PRIMARY KEY AUTOINCREMENT,
		text         TEXT NOT NULL,
		lang         TEXT NOT NULL,
		translated   TEXT NOT NULL,
		vector       BLOB,
		count        INTEGER NOT NULL,
		status       TEXT NOT NULL CHECK (status IN ('open', 'ignored', 'resolved')),
		resolution   TEXT NOT NULL,
		staff_answer TEXT NOT NULL,
		first_seen   TEXT NOT NULL,
		last_seen    TEXT NOT NULL
	);
	CREATE INDEX pending_question_status ON pending_question (status, count);
	CREATE TABLE pending_question_hit (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		question_id INTEGER NOT NULL REFERENCES pending_question (id) ON DELETE CASCADE,
		time        TEXT NOT NULL,
		platform    TEXT NOT NULL,
		user_id     TEXT NOT NULL,
		chat_id     TEXT NOT NULL,
		text        TEXT NOT NULL,
		lang        TEXT NOT NULL,
		source      TEXT NOT NULL,
		ticket_id   INTEGER
	);
	CREATE INDEX pending_question_hit_question ON pending_question_hit (question_id);
	CREATE INDEX pending_question_hit_ticket ON pending_question_hit (ticket_id);
	`,
	// v11：對話紀錄的 send_failed（#43）：機器人的回覆送不出去是 1（照樣記下回覆和 AI 判斷），其他是 0
	`
	ALTER TABLE conversation_log ADD COLUMN send_failed INTEGER NOT NULL DEFAULT 0;
	`,
	// v12：自動結束轉人工的預設改成 30 分鐘（#55）。沒有記「存過沒有」，所以全部欄位都還是舊預設值的那一列當成沒存過
	`
	UPDATE handoff_setting SET auto_end_minutes = 30
	WHERE auto_end_minutes = 0 AND on_no_answer = 1 AND streak_limit = 3 AND service_hours = '[]' AND prompt_text = '';
	`,
	// v13：預約欄位（#67）。booking_setting 只有一列（沒有這一列是預設：問人數）；booking_field 是自訂欄位，
	// key 是代號（ask_<key>、booking:<key>、booking.extra 的鍵）。都沒有資料時和 #67 以前相同。
	// booking.extra 是自訂欄位的值（JSON 物件，代號 → 值；沒有是 NULL）
	`
	CREATE TABLE booking_setting (
		id             INTEGER PRIMARY KEY CHECK (id = 1),
		people_enabled INTEGER NOT NULL
	);
	CREATE TABLE booking_field (
		id          INTEGER PRIMARY KEY AUTOINCREMENT,
		position    INTEGER NOT NULL,
		key         TEXT NOT NULL UNIQUE,
		name        TEXT NOT NULL,
		description TEXT NOT NULL,
		required    INTEGER NOT NULL
	);
	ALTER TABLE booking ADD COLUMN extra TEXT;
	`,
}

func (d *DB) migrate(ctx context.Context) error {
	var version int
	if err := d.sql.QueryRowContext(ctx, "PRAGMA user_version").Scan(&version); err != nil {
		return err
	}
	if version > len(migrations) {
		return fmt.Errorf("資料庫版本 %d 比程式新（程式只到 %d）", version, len(migrations))
	}
	for v := version; v < len(migrations); v++ {
		tx, err := d.sql.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, migrations[v]); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration v%d：%w", v+1, err)
		}
		if _, err := tx.ExecContext(ctx, fmt.Sprintf("PRAGMA user_version = %d", v+1)); err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Commit(); err != nil {
			return err
		}
	}
	return nil
}

// 時間一律存 UTC、固定寬度，字串比大小就是時間先後
const timeFormat = "2006-01-02T15:04:05.000Z"

func Timestamp(t time.Time) string { return t.UTC().Format(timeFormat) }

func parseTimestamp(s string) (time.Time, error) { return time.Parse(timeFormat, s) }

// MarkUpdate 第一次看到這個 update 回傳 true 並記下來；重複的回傳 false
func (d *DB) MarkUpdate(ctx context.Context, updateID int64, now time.Time) (bool, error) {
	res, err := d.sql.ExecContext(ctx, "INSERT OR IGNORE INTO processed_update (update_id, created_at) VALUES (?, ?)", updateID, Timestamp(now))
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

// LoadCustomer 讀出顧客狀態的 JSON；沒有這位顧客時回 nil
func (d *DB) LoadCustomer(ctx context.Context, platform, userID string) ([]byte, error) {
	var data string
	err := d.sql.QueryRowContext(ctx, "SELECT data FROM customer_state WHERE platform = ? AND user_id = ?", platform, userID).Scan(&data)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

// SaveCustomer 整份存回顧客狀態的 JSON
func (d *DB) SaveCustomer(ctx context.Context, platform, userID string, data []byte, now time.Time) error {
	_, err := d.sql.ExecContext(ctx, `
		INSERT INTO customer_state (platform, user_id, data, updated_at) VALUES (?, ?, ?, ?)
		ON CONFLICT (platform, user_id) DO UPDATE SET data = excluded.data, updated_at = excluded.updated_at`,
		platform, userID, string(data), Timestamp(now))
	return err
}

// Ticket 是待處理事項（給員工處理的事）。Kind：question、booking、urgent、unsupported、claim、handoff
type Ticket struct {
	ID        int64
	Kind      string
	Platform  string
	UserID    string
	ChatID    string
	Original  string // 顧客的原話
	Text      string // 給員工看的通知文字
	RelatedTo int64  // 緊急通知：對應的原本那張；0 是沒有
}

// AddTicket 開一張待處理事項（狀態 open），回傳編號
func (d *DB) AddTicket(ctx context.Context, t Ticket, now time.Time) (int64, error) {
	var related any
	if t.RelatedTo != 0 {
		related = t.RelatedTo
	}
	res, err := d.sql.ExecContext(ctx, `
		INSERT INTO ticket (kind, platform, user_id, chat_id, original, text, status, related_to, created_at)
		VALUES (?, ?, ?, ?, ?, ?, 'open', ?, ?)`,
		t.Kind, t.Platform, t.UserID, t.ChatID, t.Original, t.Text, related, Timestamp(now))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// OpenTickets 依開立的順序列出還沒處理的待處理事項
func (d *DB) OpenTickets(ctx context.Context) ([]Ticket, error) {
	rows, err := d.sql.QueryContext(ctx, "SELECT "+ticketColumns+" FROM ticket WHERE status = 'open' ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Ticket
	for rows.Next() {
		t, err := scanTicket(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

const ticketColumns = "id, kind, platform, user_id, chat_id, original, text, COALESCE(related_to, 0)"

func scanTicket(row interface{ Scan(...any) error }) (Ticket, error) {
	var t Ticket
	err := row.Scan(&t.ID, &t.Kind, &t.Platform, &t.UserID, &t.ChatID, &t.Original, &t.Text, &t.RelatedTo)
	return t, err
}

// Ticket 讀一張待處理事項；沒有這張回 sql.ErrNoRows
func (d *DB) Ticket(ctx context.Context, id int64) (Ticket, error) {
	return scanTicket(d.sql.QueryRowContext(ctx, "SELECT "+ticketColumns+" FROM ticket WHERE id = ?", id))
}

// ReplyTicket：員工回覆了這張，這張、它對應的原本那張和同一件事的緊急通知都改成已處理
func (d *DB) ReplyTicket(ctx context.Context, id int64) error {
	_, err := d.sql.ExecContext(ctx, `
		WITH root AS (SELECT COALESCE(related_to, id) AS id FROM ticket WHERE id = ?)
		UPDATE ticket SET status = 'replied'
		WHERE status = 'open' AND (id = ? OR id = (SELECT id FROM root) OR related_to = (SELECT id FROM root))`, id, id)
	return err
}

// CloseHandoffTickets：結束轉人工時，這位顧客還沒處理的轉人工事項改成已處理
func (d *DB) CloseHandoffTickets(ctx context.Context, platform, userID string) error {
	_, err := d.sql.ExecContext(ctx, "UPDATE ticket SET status = 'replied' WHERE status = 'open' AND kind = 'handoff' AND platform = ? AND user_id = ?",
		platform, userID)
	return err
}

// SetTicketGroupMessage 記下員工群組裡那則通知的編號
func (d *DB) SetTicketGroupMessage(ctx context.Context, id, messageID int64) error {
	_, err := d.sql.ExecContext(ctx, "UPDATE ticket SET group_message_id = ? WHERE id = ?", messageID, id)
	return err
}

// SetTicketForwardedMessage 記下轉送到員工群組的原始訊息（圖片、語音、檔案）的編號
func (d *DB) SetTicketForwardedMessage(ctx context.Context, id, messageID int64) error {
	_, err := d.sql.ExecContext(ctx, "UPDATE ticket SET forwarded_message_id = ? WHERE id = ?", messageID, id)
	return err
}

// TicketByGroupMessage 用員工群組裡的訊息編號（通知或轉送的原始訊息）找待處理事項；找不到時 ok 是 false
func (d *DB) TicketByGroupMessage(ctx context.Context, messageID int64) (Ticket, bool, error) {
	t, err := scanTicket(d.sql.QueryRowContext(ctx, "SELECT "+ticketColumns+
		" FROM ticket WHERE group_message_id = ? OR forwarded_message_id = ? ORDER BY id DESC LIMIT 1", messageID, messageID))
	if errors.Is(err, sql.ErrNoRows) {
		return t, false, nil
	}
	return t, err == nil, err
}

// Booking 是預約申請：編號和開出來的待處理事項相同；Status 是 pending、confirmed、rejected
type Booking struct {
	ID       int64
	Platform string
	UserID   string
	ChatID   string
	Date     string
	Time     string
	People   string
	Name     string
	Phone    string
	Extra    map[string]string // 自訂欄位的值（#67）：代號 → 值；沒有自訂欄位是 nil
	Summary  string
	Status   string
}

// AddBooking 存一筆待確認的預約申請（ID 用預約申請那張待處理事項的編號）
func (d *DB) AddBooking(ctx context.Context, b Booking, now time.Time) error {
	extra, err := extraJSON(b.Extra)
	if err != nil {
		return err
	}
	_, err = d.sql.ExecContext(ctx, `
		INSERT INTO booking (id, platform, user_id, chat_id, date, time, people, name, phone, extra, summary, status, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 'pending', ?)`,
		b.ID, b.Platform, b.UserID, b.ChatID, b.Date, b.Time, b.People, b.Name, b.Phone, extra, b.Summary, Timestamp(now))
	return err
}

// extraJSON 把自訂欄位的值存成 JSON；沒有值存 NULL
func extraJSON(extra map[string]string) (any, error) {
	if len(extra) == 0 {
		return nil, nil
	}
	b, err := json.Marshal(extra)
	return string(b), err
}

// parseExtra 讀回自訂欄位的值；NULL（沒有自訂欄位、#67 以前的預約）是 nil
func parseExtra(s string) (map[string]string, error) {
	if s == "" {
		return nil, nil
	}
	var m map[string]string
	return m, json.Unmarshal([]byte(s), &m)
}

// Booking 讀一筆預約申請
func (d *DB) Booking(ctx context.Context, id int64) (Booking, error) {
	b := Booking{ID: id}
	var extra string
	err := d.sql.QueryRowContext(ctx, `
		SELECT platform, user_id, chat_id, COALESCE(date, ''), COALESCE(time, ''), COALESCE(people, ''), COALESCE(name, ''), COALESCE(phone, ''),
			COALESCE(extra, ''), summary, status
		FROM booking WHERE id = ?`, id).
		Scan(&b.Platform, &b.UserID, &b.ChatID, &b.Date, &b.Time, &b.People, &b.Name, &b.Phone, &extra, &b.Summary, &b.Status)
	if err != nil {
		return b, err
	}
	b.Extra, err = parseExtra(extra)
	return b, err
}

// BookingRow 是後台預約列表的一筆
type BookingRow struct {
	ID        int64             `json:"id"`
	Platform  string            `json:"platform"`
	UserID    string            `json:"user_id"`
	Date      string            `json:"date"`
	Time      string            `json:"time"`
	People    string            `json:"people"`
	Name      string            `json:"name"`
	Phone     string            `json:"phone"`
	Extra     map[string]string `json:"extra"` // 自訂欄位的值（#67），沒有是 null
	Summary   string            `json:"summary"`
	Status    string            `json:"status"`
	HandledBy string            `json:"handled_by"`
	HandledAt *time.Time        `json:"handled_at"` // 還沒處理是 nil
	CreatedAt time.Time         `json:"created_at"`
}

// BookingFilter：Status 空字串是全部；From、To 是預約日期（YYYY-MM-DD，含頭尾），空字串是不限
type BookingFilter struct {
	Status, From, To string
}

// maxBookings 是預約列表最多列幾筆
const maxBookings = 500

// Bookings 列出預約，新申請的在前（最多 maxBookings 筆）
func (d *DB) Bookings(ctx context.Context, f BookingFilter) ([]BookingRow, error) {
	rows, err := d.sql.QueryContext(ctx, `
		SELECT id, platform, user_id, COALESCE(date, ''), COALESCE(time, ''), COALESCE(people, ''), COALESCE(name, ''), COALESCE(phone, ''),
			COALESCE(extra, ''), summary, status, COALESCE(handled_by, ''), COALESCE(handled_at, ''), created_at
		FROM booking
		WHERE (? = '' OR status = ?) AND (? = '' OR date >= ?) AND (? = '' OR date <= ?)
		ORDER BY created_at DESC, id DESC LIMIT ?`,
		f.Status, f.Status, f.From, f.From, f.To, f.To, maxBookings)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []BookingRow{}
	for rows.Next() {
		var b BookingRow
		var extra, handledAt, createdAt string
		if err := rows.Scan(&b.ID, &b.Platform, &b.UserID, &b.Date, &b.Time, &b.People, &b.Name, &b.Phone,
			&extra, &b.Summary, &b.Status, &b.HandledBy, &handledAt, &createdAt); err != nil {
			return nil, err
		}
		if b.Extra, err = parseExtra(extra); err != nil {
			return nil, err
		}
		if b.CreatedAt, err = parseTimestamp(createdAt); err != nil {
			return nil, err
		}
		if handledAt != "" {
			t, err := parseTimestamp(handledAt)
			if err != nil {
				return nil, err
			}
			b.HandledAt = &t
		}
		out = append(out, b)
	}
	return out, rows.Err()
}

// ErrAlreadyHandled：預約已經確認或拒絕過了
var ErrAlreadyHandled = errors.New("已經處理過了")

// SetBookingStatus 確認（confirmed）或拒絕（rejected）待確認的預約，記下操作的員工，預約申請那張事項改成已處理。
// 沒有這筆回 sql.ErrNoRows；已經處理過回 ErrAlreadyHandled（狀態不變）
func (d *DB) SetBookingStatus(ctx context.Context, id int64, status, by string, now time.Time) error {
	tx, err := d.sql.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.ExecContext(ctx, "UPDATE booking SET status = ?, handled_by = ?, handled_at = ? WHERE id = ? AND status = 'pending'",
		status, by, Timestamp(now), id)
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		// 沒改到：有這筆就是處理過了，沒有這筆回 sql.ErrNoRows
		if err := tx.QueryRowContext(ctx, "SELECT 1 FROM booking WHERE id = ?", id).Scan(new(int)); err != nil {
			return err
		}
		return ErrAlreadyHandled
	}
	if _, err := tx.ExecContext(ctx, "UPDATE ticket SET status = 'replied' WHERE id = ?", id); err != nil {
		return err
	}
	return tx.Commit()
}

// LogEntry 是對話紀錄的一筆（一則訊息）
type LogEntry struct {
	ID        int64 // 讀出來時才有；AppendLog 不看
	Time      time.Time
	Platform  string
	UserID    string
	ChatID    string
	Direction string // customer、bot、staff
	Text      string
	ReplyKey  string          // 固定回覆的 key；空字串存成 NULL
	Traces    json.RawMessage // AI 判斷；nil 存成 []
	// SendFailed 是機器人的回覆沒有送到顧客（例如 Telegram 送出失敗）
	SendFailed bool
}

func (d *DB) AppendLog(ctx context.Context, e LogEntry) error {
	traces := string(e.Traces)
	if traces == "" {
		traces = "[]"
	}
	var replyKey any
	if e.ReplyKey != "" {
		replyKey = e.ReplyKey
	}
	_, err := d.sql.ExecContext(ctx, `
		INSERT INTO conversation_log (time, platform, user_id, chat_id, direction, text, reply_key, traces, send_failed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		Timestamp(e.Time), e.Platform, e.UserID, e.ChatID, e.Direction, e.Text, replyKey, traces, e.SendFailed)
	return err
}

// Logs 依寫入順序讀出某位顧客的對話紀錄
func (d *DB) Logs(ctx context.Context, platform, userID string) ([]LogEntry, error) {
	return d.queryLogs(ctx, `
		SELECT id, time, platform, user_id, chat_id, direction, text, COALESCE(reply_key, ''), traces, send_failed
		FROM conversation_log WHERE platform = ? AND user_id = ? ORDER BY id`, platform, userID)
}

// LogEntry 讀一則對話紀錄；沒有這則（或已經清掉）回 sql.ErrNoRows
func (d *DB) LogEntry(ctx context.Context, id int64) (LogEntry, error) {
	logs, err := d.queryLogs(ctx, `
		SELECT id, time, platform, user_id, chat_id, direction, text, COALESCE(reply_key, ''), traces, send_failed
		FROM conversation_log WHERE id = ?`, id)
	if err != nil {
		return LogEntry{}, err
	}
	if len(logs) == 0 {
		return LogEntry{}, sql.ErrNoRows
	}
	return logs[0], nil
}

// RecentLogs 是某位顧客最近的 n 則對話（舊的在前），略過轉人工中沒回覆的空白紀錄：轉人工通知附給員工看
func (d *DB) RecentLogs(ctx context.Context, platform, userID string, n int) ([]LogEntry, error) {
	return d.queryLogs(ctx, `
		SELECT id, time, platform, user_id, chat_id, direction, text, reply_key, traces, send_failed FROM (
			SELECT time, platform, user_id, chat_id, direction, text, COALESCE(reply_key, '') AS reply_key, traces, send_failed, id
			FROM conversation_log WHERE platform = ? AND user_id = ? AND text != '' ORDER BY id DESC LIMIT ?
		) ORDER BY id`, platform, userID, n)
}

func (d *DB) queryLogs(ctx context.Context, query string, args ...any) ([]LogEntry, error) {
	rows, err := d.sql.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []LogEntry
	for rows.Next() {
		var e LogEntry
		var t, traces string
		if err := rows.Scan(&e.ID, &t, &e.Platform, &e.UserID, &e.ChatID, &e.Direction, &e.Text, &e.ReplyKey, &traces, &e.SendFailed); err != nil {
			return nil, err
		}
		if e.Time, err = parseTimestamp(t); err != nil {
			return nil, err
		}
		e.Traces = json.RawMessage(traces)
		out = append(out, e)
	}
	return out, rows.Err()
}

// CustomerSummary 是後台對話紀錄的一位顧客：最後一則有文字的訊息、最後的時間和訊息數
type CustomerSummary struct {
	Platform string    `json:"platform"`
	UserID   string    `json:"user_id"`
	LastText string    `json:"last_text"`
	LastTime time.Time `json:"last_time"`
	Messages int       `json:"messages"`
}

// maxCustomers 是對話紀錄的顧客列表最多列幾位
const maxCustomers = 200

// Customers 列出有對話紀錄的顧客，最近有對話的在前（最多 maxCustomers 位）。
// query 不是空字串時只列顧客編號或任一則訊息含有它的顧客
func (d *DB) Customers(ctx context.Context, query string) ([]CustomerSummary, error) {
	like := "%" + strings.NewReplacer(`\`, `\\`, "%", `\%`, "_", `\_`).Replace(query) + "%"
	rows, err := d.sql.QueryContext(ctx, `
		SELECT platform, user_id, MAX(time), COUNT(*), COALESCE((
			SELECT text FROM conversation_log t WHERE t.platform = l.platform AND t.user_id = l.user_id AND t.text != '' ORDER BY t.id DESC LIMIT 1
		), '')
		FROM conversation_log l
		GROUP BY platform, user_id
		HAVING ? = '' OR user_id LIKE ? ESCAPE '\' OR SUM(text LIKE ? ESCAPE '\') > 0
		ORDER BY MAX(time) DESC, MAX(id) DESC LIMIT ?`,
		query, like, like, maxCustomers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []CustomerSummary{}
	for rows.Next() {
		var c CustomerSummary
		var last string
		if err := rows.Scan(&c.Platform, &c.UserID, &last, &c.Messages, &c.LastText); err != nil {
			return nil, err
		}
		if c.LastTime, err = parseTimestamp(last); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UsageRecord 是一次 AI 呼叫的用量
type UsageRecord struct {
	Time         time.Time
	Model        string
	Request      string // 題目名稱（Intent、Pick…）
	InputTokens  int
	OutputTokens int
	CostUSD      float64
	OK           bool
}

func (d *DB) RecordUsage(ctx context.Context, u UsageRecord) error {
	_, err := d.sql.ExecContext(ctx, `
		INSERT INTO ai_usage (time, model, request, input_tokens, output_tokens, cost_usd, ok)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		Timestamp(u.Time), u.Model, u.Request, u.InputTokens, u.OutputTokens, u.CostUSD, u.OK)
	return err
}

// CostBetween 加總 [from, to) 之間的費用（美元）
func (d *DB) CostBetween(ctx context.Context, from, to time.Time) (float64, error) {
	var cost float64
	err := d.sql.QueryRowContext(ctx, "SELECT COALESCE(SUM(cost_usd), 0) FROM ai_usage WHERE time >= ? AND time < ?",
		Timestamp(from), Timestamp(to)).Scan(&cost)
	return cost, err
}

// UsageSummary 是一個模型在一段時間的用量：呼叫次數（含失敗）、失敗次數、token 和費用（美元）
type UsageSummary struct {
	Model        string  `json:"model"`
	Calls        int     `json:"calls"`
	Failed       int     `json:"failed"`
	InputTokens  int     `json:"input_tokens"`
	OutputTokens int     `json:"output_tokens"`
	CostUSD      float64 `json:"cost_usd"`
}

// UsageBetween 依模型加總 [from, to) 之間的用量，費用高的在前（費用相同時呼叫多的在前）
func (d *DB) UsageBetween(ctx context.Context, from, to time.Time) ([]UsageSummary, error) {
	rows, err := d.sql.QueryContext(ctx, `
		SELECT model, COUNT(*), SUM(ok = 0), SUM(input_tokens), SUM(output_tokens), SUM(cost_usd)
		FROM ai_usage WHERE time >= ? AND time < ?
		GROUP BY model ORDER BY SUM(cost_usd) DESC, COUNT(*) DESC, model`,
		Timestamp(from), Timestamp(to))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []UsageSummary{}
	for rows.Next() {
		var u UsageSummary
		if err := rows.Scan(&u.Model, &u.Calls, &u.Failed, &u.InputTokens, &u.OutputTokens, &u.CostUSD); err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}

// ModelSetting 讀後台存的模型設定；還沒存過時 ok 是 false
func (d *DB) ModelSetting(ctx context.Context) (model string, budget float64, ok bool, err error) {
	err = d.sql.QueryRowContext(ctx, "SELECT model, monthly_budget FROM model_setting WHERE id = 1").Scan(&model, &budget)
	if errors.Is(err, sql.ErrNoRows) {
		return "", 0, false, nil
	}
	return model, budget, err == nil, err
}

// SaveModelSetting 存（或取代）模型設定；值要先檢查過（model.Switch.Apply）
func (d *DB) SaveModelSetting(ctx context.Context, model string, budget float64) error {
	_, err := d.sql.ExecContext(ctx, "INSERT OR REPLACE INTO model_setting (id, model, monthly_budget) VALUES (1, ?, ?)", model, budget)
	return err
}

// Cleanup 刪掉超過保留天數的對話紀錄、update 紀錄和過期的管理員 session，回傳各刪了幾筆
func (d *DB) Cleanup(ctx context.Context, now time.Time) (logs, updates, sessions int64, err error) {
	res, err := d.sql.ExecContext(ctx, "DELETE FROM conversation_log WHERE time < ?", Timestamp(now.Add(-LogRetention)))
	if err != nil {
		return 0, 0, 0, err
	}
	if logs, err = res.RowsAffected(); err != nil {
		return 0, 0, 0, err
	}
	res, err = d.sql.ExecContext(ctx, "DELETE FROM processed_update WHERE created_at < ?", Timestamp(now.Add(-UpdateRetention)))
	if err != nil {
		return logs, 0, 0, err
	}
	if updates, err = res.RowsAffected(); err != nil {
		return logs, 0, 0, err
	}
	res, err = d.sql.ExecContext(ctx, "DELETE FROM admin_session WHERE expires_at <= ?", Timestamp(now))
	if err != nil {
		return logs, updates, 0, err
	}
	sessions, err = res.RowsAffected()
	return logs, updates, sessions, err
}
