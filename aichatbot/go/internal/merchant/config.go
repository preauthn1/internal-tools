// Package merchant 是商家設定（CONTEXT.md「商家設定」）：存在 SQLite，可以匯出／匯入成一個 JSON 檔，
// 也能從 eval/ 的森日茶屋檔案匯入（importeval.go）。
package merchant

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Config 是整份商家設定。slice 的順序就是後台顯示的順序，匯出時照存的順序輸出
type Config struct {
	Merchant  Merchant   `json:"merchant"`
	Knowledge []Section  `json:"knowledge"`
	FAQ       []FAQ      `json:"faq"`
	Examples  []Example  `json:"intent_examples"`
	Features  []Feature  `json:"features"`
	Replies   []Reply    `json:"fixed_replies"`
	Languages []Language `json:"languages"`
	// 舊的匯出檔沒有這個欄位，讀進來是 nil（當成沒有自訂意圖）
	CustomIntents []CustomIntent `json:"custom_intents"`
	// 預約欄位（#67）：nil 是預設（問人數、沒有自訂欄位），匯出時不輸出
	BookingFields *BookingFields `json:"booking_fields,omitempty"`
}

// Merchant 是商家資訊：名稱、介紹、聯絡方式必填
type Merchant struct {
	Name              string `json:"name"`
	Intro             string `json:"intro"`
	Contact           string `json:"contact"`
	PreferredLanguage string `json:"preferred_language"`
}

// Section 是知識段落：有標題的自由文字
type Section struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// FAQ 只存慣用語言；AltQuestions 是同一題的其他問法
type FAQ struct {
	Key          string   `json:"key"`
	Question     string   `json:"question"`
	AltQuestions []string `json:"alt_questions"`
	Answer       string   `json:"answer"`
}

// Example 是意圖案例庫的一句；Builtin 是內建案例（商家不能改）。
// State 是分類當下的對話狀態：ticket（有事在等員工）、booking:欄位（預約中問到該欄位）；空的是沒有特別的狀態
type Example struct {
	Text    string `json:"text"`
	Intent  string `json:"intent"`
	State   string `json:"state,omitempty"`
	Builtin bool   `json:"builtin"`
}

// Feature 是功能開關（booking、handoff）；關閉時回的文字是固定回覆 booking_off、handoff_off
type Feature struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// UnmarshalJSON 忽略舊備份（#37 以前）的 off_text，其他不認得的欄位照樣報錯（ParseJSON 的檢查不會傳進這裡）
func (f *Feature) UnmarshalJSON(data []byte) error {
	type feature Feature
	var v struct {
		feature
		OffText string `json:"off_text"`
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return err
	}
	*f = Feature(v.feature)
	return nil
}

// Reply 是固定回覆的一種語言；Lang 是 * 時所有語言共用
type Reply struct {
	Key  string `json:"key"`
	Lang string `json:"lang"`
	Text string `json:"text"`
}

// Language 是顧客語言
type Language struct {
	Code    string `json:"code"`
	Enabled bool   `json:"enabled"`
}

// Validate 檢查必填欄位
func (c Config) Validate() error {
	var errs []string
	m := c.Merchant
	for _, f := range []struct{ name, v string }{{"商家名稱", m.Name}, {"商家介紹", m.Intro}, {"聯絡方式", m.Contact}, {"慣用語言", m.PreferredLanguage}} {
		if strings.TrimSpace(f.v) == "" {
			errs = append(errs, f.name+"是必填")
		}
	}
	found, anyEnabled, codes := false, false, map[string]bool{}
	for _, l := range c.Languages {
		found = found || (l.Code == m.PreferredLanguage && l.Enabled)
		anyEnabled = anyEnabled || l.Enabled
		if !slices.Contains(SupportedLanguages, l.Code) {
			errs = append(errs, "不認得的語言代碼："+l.Code)
		}
		if codes[l.Code] {
			errs = append(errs, "語言重複："+l.Code)
		}
		codes[l.Code] = true
	}
	if !anyEnabled {
		errs = append(errs, "至少要勾一種顧客語言")
	}
	if m.PreferredLanguage != "" && !found {
		errs = append(errs, fmt.Sprintf("慣用語言 %s 不在啟用的語言裡", m.PreferredLanguage))
	}
	// 挑知識時用標題選段落，標題重複會選不清楚
	titles := map[string]bool{}
	for i, s := range c.Knowledge {
		for _, x := range []struct{ name, v string }{{"標題", s.Title}, {"內容", s.Body}} {
			if strings.TrimSpace(x.v) == "" {
				errs = append(errs, fmt.Sprintf("第 %d 段知識的%s是必填", i+1, x.name))
			}
		}
		title := strings.TrimSpace(s.Title)
		if title != "" && titles[title] {
			errs = append(errs, "知識段落標題重複："+title)
		}
		titles[title] = true
	}
	keys := map[string]bool{}
	for i, f := range c.FAQ {
		if keys[f.Key] {
			errs = append(errs, "FAQ 代號重複："+f.Key)
		}
		keys[f.Key] = true
		for _, x := range []struct{ name, v string }{{"代號", f.Key}, {"問題", f.Question}, {"答案", f.Answer}} {
			if strings.TrimSpace(x.v) == "" {
				errs = append(errs, fmt.Sprintf("第 %d 題 FAQ 的%s是必填", i+1, x.name))
			}
		}
	}
	errs = append(errs, validateCustomIntents(c)...)
	errs = append(errs, validateBookingFields(c)...)
	if len(errs) > 0 {
		return &InvalidError{errs}
	}
	return nil
}

// InvalidError 是 Validate 沒過：每一項說明錯在哪裡（後台直接顯示給管理員）
type InvalidError struct{ Problems []string }

func (e *InvalidError) Error() string { return strings.Join(e.Problems, "；") }

// Export 輸出成縮排的 JSON（不跳脫中文和 <>&，結尾換行）
func Export(c Config) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseJSON 讀 Export 輸出的格式；不認得的欄位視為錯誤，避免打錯字的設定被默默丟掉
func ParseJSON(data []byte) (Config, error) {
	var c Config
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&c); err != nil {
		return Config{}, err
	}
	return c, nil
}

// Save 用整份設定取代資料庫裡的商家設定（同一個 transaction）。
// 意圖案例只取代內建的，商家新增的保留；回傳因為和內建案例重複而刪掉的商家新增的案例
func Save(ctx context.Context, db *sql.DB, c Config) (removed []Example, err error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if removed, err = writeAll(ctx, tx, c); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return removed, nil
}

// writeAll 用 c 取代每一張商家設定的表（意圖案例見 replaceExamples）
func writeAll(ctx context.Context, tx *sql.Tx, c Config) (removed []Example, err error) {
	for _, table := range []string{"feature_switch", "fixed_reply"} {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+table); err != nil {
			return nil, err
		}
	}
	if err := writeMerchant(ctx, tx, c.Merchant); err != nil {
		return nil, err
	}
	if err := replaceKnowledge(ctx, tx, c.Knowledge); err != nil {
		return nil, err
	}
	if err := replaceFAQ(ctx, tx, c.FAQ); err != nil {
		return nil, err
	}
	if err := replaceLanguages(ctx, tx, c.Languages); err != nil {
		return nil, err
	}
	if removed, err = replaceExamples(ctx, tx, c.Examples); err != nil {
		return nil, err
	}
	if err := replaceCustomIntents(ctx, tx, c.CustomIntents); err != nil {
		return nil, err
	}
	if err := replaceBookingFields(ctx, tx, c.BookingFields); err != nil {
		return nil, err
	}
	for _, f := range c.Features {
		if _, err := tx.ExecContext(ctx, "INSERT INTO feature_switch (name, enabled) VALUES (?, ?)", f.Name, f.Enabled); err != nil {
			return nil, err
		}
	}
	for _, r := range c.Replies {
		if _, err := tx.ExecContext(ctx, "INSERT INTO fixed_reply (key, lang, text) VALUES (?, ?, ?)", r.Key, r.Lang, r.Text); err != nil {
			return nil, fmt.Errorf("固定回覆 %s（%s）：%w", r.Key, r.Lang, err)
		}
	}
	return removed, nil
}

// replaceExamples 用 file 的內建案例取代資料庫的內建案例，商家新增的案例接在後面寫回。
// 商家新增的案例是資料庫裡的加上 file 裡非內建的（bot export 的備份），比對後相同（見 Normalize）的只留先出現的一句；
// 和內建案例相同的只留內建的，回傳被刪掉的商家新增的案例
func replaceExamples(ctx context.Context, tx *sql.Tx, file []Example) (removed []Example, err error) {
	existing, err := listExamples(ctx, tx)
	if err != nil {
		return nil, err
	}
	var rows, added []Example
	for _, e := range existing {
		if !e.Builtin {
			added = append(added, Example{Text: e.Text, Intent: e.Intent, State: e.State})
		}
	}
	seen, builtin := map[string]bool{}, map[string]bool{}
	for _, e := range file {
		if e.Builtin {
			rows = append(rows, e)
			builtin[Normalize(e.Text)] = true
		} else {
			added = append(added, e)
		}
	}
	for _, e := range added {
		key := Normalize(e.Text)
		if seen[key] {
			continue
		}
		seen[key] = true
		if builtin[key] {
			removed = append(removed, e)
		} else {
			rows = append(rows, e)
		}
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM intent_example"); err != nil {
		return nil, err
	}
	for _, e := range rows {
		if _, err := tx.ExecContext(ctx, "INSERT INTO intent_example (text, intent, state, builtin) VALUES (?, ?, ?, ?)", e.Text, e.Intent, e.State, e.Builtin); err != nil {
			return nil, err
		}
	}
	return removed, nil
}

// ErrNotConfigured：資料庫裡還沒有商家設定（還沒匯入）
var ErrNotConfigured = errors.New("還沒有商家設定，請先匯入")

// Load 讀出整份商家設定。固定回覆依 key、lang 排序，功能開關依名稱排序，其餘依存的順序。
// 不做快取：機器人和後台每次用到都直接 Load，後台存檔後馬上生效，不用重開服務
func Load(ctx context.Context, db *sql.DB) (Config, error) { return load(ctx, db) }

// querier 是 *sql.DB 或 *sql.Tx（資料庫只開一條連線，transaction 裡要用 tx 讀）
type querier interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func load(ctx context.Context, db querier) (Config, error) {
	var c Config
	m := &c.Merchant
	err := db.QueryRowContext(ctx, "SELECT name, intro, contact, preferred_language FROM merchant WHERE id = 1").
		Scan(&m.Name, &m.Intro, &m.Contact, &m.PreferredLanguage)
	if errors.Is(err, sql.ErrNoRows) {
		return Config{}, ErrNotConfigured
	}
	if err != nil {
		return Config{}, err
	}
	if c.Knowledge, err = query(ctx, db, "SELECT title, body FROM knowledge_section ORDER BY position, id", func(r *sql.Rows) (s Section, err error) {
		return s, r.Scan(&s.Title, &s.Body)
	}); err != nil {
		return Config{}, err
	}
	if c.FAQ, err = query(ctx, db, "SELECT key, question, alt_questions, answer FROM faq ORDER BY position, id", func(r *sql.Rows) (f FAQ, err error) {
		var alt string
		if err = r.Scan(&f.Key, &f.Question, &alt, &f.Answer); err != nil {
			return f, err
		}
		return f, json.Unmarshal([]byte(alt), &f.AltQuestions)
	}); err != nil {
		return Config{}, err
	}
	if c.Examples, err = query(ctx, db, "SELECT text, intent, state, builtin FROM intent_example ORDER BY id", func(r *sql.Rows) (e Example, err error) {
		return e, r.Scan(&e.Text, &e.Intent, &e.State, &e.Builtin)
	}); err != nil {
		return Config{}, err
	}
	if c.Features, err = query(ctx, db, "SELECT name, enabled FROM feature_switch ORDER BY name", func(r *sql.Rows) (f Feature, err error) {
		return f, r.Scan(&f.Name, &f.Enabled)
	}); err != nil {
		return Config{}, err
	}
	if c.Replies, err = query(ctx, db, "SELECT key, lang, text FROM fixed_reply ORDER BY key, lang", func(r *sql.Rows) (x Reply, err error) {
		return x, r.Scan(&x.Key, &x.Lang, &x.Text)
	}); err != nil {
		return Config{}, err
	}
	if c.Languages, err = query(ctx, db, "SELECT code, enabled FROM language ORDER BY position, code", func(r *sql.Rows) (l Language, err error) {
		return l, r.Scan(&l.Code, &l.Enabled)
	}); err != nil {
		return Config{}, err
	}
	if c.CustomIntents, err = listCustomIntents(ctx, db); err != nil {
		return Config{}, err
	}
	if c.BookingFields, err = loadBookingFields(ctx, db); err != nil {
		return Config{}, err
	}
	return c, nil
}

// query 讀出每一列；沒有資料時回傳空 slice（匯出成 []，不是 null）
func query[T any](ctx context.Context, db querier, q string, scan func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []T{}
	for rows.Next() {
		v, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}
