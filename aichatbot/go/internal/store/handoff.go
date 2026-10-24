package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
)

// HandoffSetting 是轉人工設定和綁定的員工群組（v7 的 handoff_setting）。總開關是功能開關 handoff（feature_switch）
type HandoffSetting struct {
	StaffChatID    int64          // 綁定的員工群組；0 是還沒綁定
	BindCode       string         // 員工在群組裡輸入這組碼就綁定那個群組
	OnNoAnswer     bool           // 找不到答案時轉人工
	StreakLimit    int            // 聲稱／無關累計幾次轉人工
	ServiceHours   []ServiceHours // 員工服務時段（台北時間）；空的是全天
	AutoEndMinutes int            // 員工最後一次回覆後幾分鐘恢復 AI 回覆；0 是不自動結束，預設 30
	PromptText     string         // 轉人工時回顧客的話；空的用固定回覆 handoff
}

// ServiceHours 是某個星期幾的服務時段：Weekday 0 是星期日；Start、End 是 "09:00" 這種 24 小時制，End 不含
type ServiceHours struct {
	Weekday int    `json:"weekday"`
	Start   string `json:"start"`
	End     string `json:"end"`
}

// DefaultHandoffSetting 是還沒存過時的設定（綁定碼另外產生）
func DefaultHandoffSetting() HandoffSetting {
	return HandoffSetting{OnNoAnswer: true, StreakLimit: 3, AutoEndMinutes: 30}
}

// HandoffSetting 讀轉人工設定；還沒有那一列時用預設值並產生綁定碼存起來
func (d *DB) HandoffSetting(ctx context.Context) (HandoffSetting, error) {
	var s HandoffSetting
	var chat sql.NullInt64
	var hours string
	err := d.sql.QueryRowContext(ctx, `
		SELECT staff_chat_id, bind_code, on_no_answer, streak_limit, service_hours, auto_end_minutes, prompt_text
		FROM handoff_setting WHERE id = 1`).
		Scan(&chat, &s.BindCode, &s.OnNoAnswer, &s.StreakLimit, &hours, &s.AutoEndMinutes, &s.PromptText)
	if errors.Is(err, sql.ErrNoRows) {
		s = DefaultHandoffSetting()
		s.BindCode = newBindCode()
		// 兩邊同時建立時以先寫進去的為準
		if _, err := d.sql.ExecContext(ctx, `
			INSERT OR IGNORE INTO handoff_setting (id, staff_chat_id, bind_code, on_no_answer, streak_limit, service_hours, auto_end_minutes, prompt_text)
			VALUES (1, NULL, ?, ?, ?, '[]', ?, ?)`, s.BindCode, s.OnNoAnswer, s.StreakLimit, s.AutoEndMinutes, s.PromptText); err != nil {
			return s, err
		}
		return d.HandoffSetting(ctx)
	}
	if err != nil {
		return s, err
	}
	s.StaffChatID = chat.Int64
	return s, json.Unmarshal([]byte(hours), &s.ServiceHours)
}

// SaveHandoffSetting 存轉人工設定；綁定的群組和綁定碼不改（那兩個由 BindStaffChat 改）
func (d *DB) SaveHandoffSetting(ctx context.Context, s HandoffSetting) error {
	if _, err := d.HandoffSetting(ctx); err != nil { // 確定那一列存在
		return err
	}
	hours := s.ServiceHours
	if hours == nil {
		hours = []ServiceHours{}
	}
	data, err := json.Marshal(hours)
	if err != nil {
		return err
	}
	_, err = d.sql.ExecContext(ctx, `
		UPDATE handoff_setting SET on_no_answer = ?, streak_limit = ?, service_hours = ?, auto_end_minutes = ?, prompt_text = ?
		WHERE id = 1`, s.OnNoAnswer, s.StreakLimit, string(data), s.AutoEndMinutes, s.PromptText)
	return err
}

// BindStaffChat：code 是目前的綁定碼才把 chatID 綁成員工群組（蓋掉之前綁的），並換一組新碼（舊碼不能再用）。碼不對回 false
func (d *DB) BindStaffChat(ctx context.Context, code string, chatID int64) (bool, error) {
	if _, err := d.HandoffSetting(ctx); err != nil {
		return false, err
	}
	res, err := d.sql.ExecContext(ctx, "UPDATE handoff_setting SET staff_chat_id = ?, bind_code = ? WHERE id = 1 AND bind_code = ?",
		chatID, newBindCode(), code)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

// RegenerateBindCode 換一組新的綁定碼（舊碼不能再用），回傳新碼；已經綁定的群組不變
func (d *DB) RegenerateBindCode(ctx context.Context) (string, error) {
	if _, err := d.HandoffSetting(ctx); err != nil { // 確定那一列存在
		return "", err
	}
	code := newBindCode()
	_, err := d.sql.ExecContext(ctx, "UPDATE handoff_setting SET bind_code = ? WHERE id = 1", code)
	return code, err
}

// MigrateStaffChat：綁定的群組 from 升級成超級群組、chat id 變成 to 時改綁 to；from 不是綁定的群組時不改
func (d *DB) MigrateStaffChat(ctx context.Context, from, to int64) (bool, error) {
	res, err := d.sql.ExecContext(ctx, "UPDATE handoff_setting SET staff_chat_id = ? WHERE id = 1 AND staff_chat_id = ?", to, from)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	return n == 1, err
}

// 綁定碼：10 個字，不用容易看錯的 0、O、1、I、L
const bindAlphabet = "23456789ABCDEFGHJKMNPQRSTUVWXYZ"

func newBindCode() string {
	b := make([]byte, 10)
	rand.Read(b)
	for i := range b {
		b[i] = bindAlphabet[int(b[i])%len(bindAlphabet)]
	}
	return string(b)
}
