package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// 只能看的頁面的 API：預約列表（確認、拒絕在員工群組）、對話紀錄（附 AI 判斷、可以翻譯一則）、本月 AI 費用

// bookingStatuses 是預約列表能篩選的狀態；空字串是全部
var bookingStatuses = map[string]bool{"": true, "pending": true, "confirmed": true, "rejected": true}

// listBookings 列出預約，新申請的在前：?status=pending、confirmed、rejected；?from=、?to= 是預約日期（YYYY-MM-DD，含頭尾）
func (s *Server) listBookings(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	f := store.BookingFilter{Status: q.Get("status"), From: q.Get("from"), To: q.Get("to")}
	if !bookingStatuses[f.Status] {
		writeError(w, http.StatusBadRequest, "狀態只能是 pending、confirmed、rejected")
		return
	}
	for _, d := range []string{f.From, f.To} {
		if _, err := time.Parse(time.DateOnly, d); d != "" && err != nil {
			writeError(w, http.StatusBadRequest, "日期要是 YYYY-MM-DD")
			return
		}
	}
	list, err := s.DB.Bookings(r.Context(), f)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// trace 是一則機器人回覆的一步 AI 判斷（dialog.Trace）：白話說明在後台頁面照 step、result 對照
type trace struct {
	Step   string            `json:"step"`
	Result string            `json:"result"`
	Detail map[string]string `json:"detail,omitempty"`
}

type message struct {
	ID           int64     `json:"id"`
	Time         time.Time `json:"time"`
	Direction    string    `json:"direction"` // customer、bot、staff
	Text         string    `json:"text"`
	ReplyKey     string    `json:"reply_key"`
	Traces       []trace   `json:"traces"`
	Translatable bool      `json:"translatable"` // 不是商家慣用語言的顧客、機器人訊息才有翻譯按鈕
	SendFailed   bool      `json:"send_failed"`  // 機器人的回覆沒有送到顧客
}

type conversation struct {
	Platform  string    `json:"platform"`
	UserID    string    `json:"user_id"`
	Lang      string    `json:"lang"`  // 顧客狀態記的語言
	Human     bool      `json:"human"` // 轉人工中
	Preferred string    `json:"preferred_language"`
	Messages  []message `json:"messages"`
}

// getConversation 回一位顧客的對話紀錄（舊的在前），每則附 AI 判斷；沒有紀錄（或已經過了保留天數）回 404
func (s *Server) getConversation(w http.ResponseWriter, r *http.Request) {
	platform, userID := r.PathValue("platform"), r.PathValue("userID")
	logs, err := s.DB.Logs(r.Context(), platform, userID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	if len(logs) == 0 {
		writeError(w, http.StatusNotFound, "找不到這位顧客的對話紀錄")
		return
	}
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	customer, err := s.customer(r, platform, userID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	lang := customer.Lang
	resp := conversation{Platform: platform, UserID: userID, Lang: lang, Human: customer.Human, Preferred: c.Merchant.PreferredLanguage, Messages: []message{}}
	for _, e := range logs {
		m := message{ID: e.ID, Time: e.Time, Direction: e.Direction, Text: e.Text, ReplyKey: e.ReplyKey, Traces: []trace{}, SendFailed: e.SendFailed}
		if err := json.Unmarshal(e.Traces, &m.Traces); err != nil {
			s.internalError(w, err)
			return
		}
		from := messageLanguage(e, lang)
		m.Translatable = from != "" && from != resp.Preferred
		resp.Messages = append(resp.Messages, m)
	}
	writeJSON(w, http.StatusOK, resp)
}

// customer 是顧客狀態；沒有這位顧客是空的狀態
func (s *Server) customer(r *http.Request, platform, userID string) (dialog.Customer, error) {
	return (&dialog.Bot{DB: s.DB}).LoadCustomer(r.Context(), platform, userID)
}

// customerLang 是顧客狀態記的語言；沒有這位顧客或還沒判斷出語言是空字串
func (s *Server) customerLang(r *http.Request, platform, userID string) (string, error) {
	c, err := s.customer(r, platform, userID)
	return c.Lang, err
}

// HandoffEnder 依顧客結束轉人工（dialog.Bot）
type HandoffEnder interface {
	EndCustomerHandoff(ctx context.Context, platform, userID string) (already bool, err error)
}

// endHandoff 是對話內容頁的「結束轉人工」（#55）：效果同員工群組的按鈕；顧客不在轉人工中回 409
func (s *Server) endHandoff(w http.ResponseWriter, r *http.Request) {
	var body struct{}
	if !decodeBody(w, r, &body) {
		return
	}
	if s.Handoff == nil {
		writeError(w, http.StatusServiceUnavailable, "沒有接上對話層，無法結束轉人工")
		return
	}
	already, err := s.Handoff.EndCustomerHandoff(r.Context(), r.PathValue("platform"), r.PathValue("userID"))
	if err != nil {
		s.internalError(w, err)
		return
	}
	if already {
		writeError(w, http.StatusConflict, "這位顧客沒有在轉人工中")
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

// messageLanguage 是一則訊息的語言：員工的訊息、空白訊息是空字串（不翻譯）；
// 其他用對話層判斷語言的規則，判斷不出來（表情、太短）用顧客的語言
func messageLanguage(e store.LogEntry, customerLang string) string {
	if e.Direction == "staff" || strings.TrimSpace(e.Text) == "" {
		return ""
	}
	if lang := dialog.DetectLanguage(e.Text); lang != "" {
		return lang
	}
	return customerLang
}

// translateLog 把一則對話翻成商家慣用語言，回 {translation}：每按一次呼叫一次目前的模型，譯文不存起來。
// 顧客的訊息照原意翻（merchant.TranslateCustomer），機器人的回覆用翻店家回覆的提示詞
func (s *Server) translateLog(w http.ResponseWriter, r *http.Request) {
	var body struct{}
	if !decodeBody(w, r, &body) {
		return
	}
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, errLogNotFound)
		return
	}
	e, err := s.DB.LogEntry(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, errLogNotFound)
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	lang, err := s.customerLang(r, e.Platform, e.UserID)
	if err != nil {
		s.internalError(w, err)
		return
	}
	from, to := messageLanguage(e, lang), c.Merchant.PreferredLanguage
	if from == "" || from == to {
		writeError(w, http.StatusBadRequest, "這則已經是慣用語言，不用翻譯")
		return
	}
	if s.Chat == nil {
		writeError(w, http.StatusServiceUnavailable, "目前沒有可用的模型")
		return
	}
	translate := merchant.Translate
	if e.Direction == "customer" {
		translate = merchant.TranslateCustomer
	}
	translated, err := translate(r.Context(), s.Chat, from, to, e.Text)
	if err != nil {
		s.Log.Warn("對話紀錄翻譯失敗", "id", id, "error", err.Error())
		writeError(w, http.StatusBadGateway, "翻譯失敗，請再試一次")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"translation": translated})
}

const errLogNotFound = "找不到這則對話（可能已經過了保留天數）"

type usage struct {
	From          time.Time            `json:"from"`
	To            time.Time            `json:"to"`
	TotalCost     float64              `json:"total_cost_usd"`
	TotalCalls    int                  `json:"total_calls"`
	ByModel       []store.UsageSummary `json:"by_model"`
	Model         string               `json:"model"` // 目前的模型開關：afm、haiku
	MonthlyBudget float64              `json:"monthly_budget_usd"`
}

// getUsage 回本月（台北時間，和費用上限同一套）的 AI 費用和呼叫次數，附目前的模型和每月上限
func (s *Server) getUsage(w http.ResponseWriter, r *http.Request) {
	from, to := model.MonthRange(s.Now())
	list, err := s.DB.UsageBetween(r.Context(), from, to)
	if err != nil {
		s.internalError(w, err)
		return
	}
	u := usage{From: from, To: to, ByModel: list}
	u.Model, u.MonthlyBudget = s.Model.Settings()
	for _, m := range list {
		u.TotalCost += m.CostUSD
		u.TotalCalls += m.Calls
	}
	writeJSON(w, http.StatusOK, u)
}

// listCustomers 是對話紀錄的顧客列表，最近有對話的在前；?q= 比對顧客編號或訊息內容
func (s *Server) listCustomers(w http.ResponseWriter, r *http.Request) {
	list, err := s.DB.Customers(r.Context(), strings.TrimSpace(r.URL.Query().Get("q")))
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}
