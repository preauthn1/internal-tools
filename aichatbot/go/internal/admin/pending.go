package admin

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

// 待補問題（CONTEXT.md）的 API：列表依次數排序、詳細頁附每次出現前後的對話；
// 動作：轉成 FAQ、加進知識段落、加成意圖案例、忽略

// pendingContext 是詳細頁每次出現附上前、後各幾則對話
const pendingContext = 10

// listPending 列出待補問題（?status=open、ignored、resolved，預設 open），次數多的在前
func (s *Server) listPending(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	if status == "" {
		status = "open"
	}
	list, err := s.DB.PendingQuestions(r.Context(), status)
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, list)
}

// logLine 是詳細頁的一則對話
type logLine struct {
	Time       time.Time `json:"time"`
	Direction  string    `json:"direction"` // customer、bot、staff
	Text       string    `json:"text"`
	SendFailed bool      `json:"send_failed"` // 機器人的回覆沒有送到顧客
}

type pendingHitView struct {
	store.PendingHit
	Conversation []logLine `json:"conversation"`
}

type pendingDetail struct {
	store.PendingQuestion
	Hits []pendingHitView `json:"hits"`
}

// getPending 回一筆待補問題、每次出現和那位顧客前後的對話（對話紀錄 90 天後清掉，之後只剩原話）
func (s *Server) getPending(w http.ResponseWriter, r *http.Request) {
	q, ok := s.pendingQuestion(w, r)
	if !ok {
		return
	}
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	// 對話層在背景翻成商家慣用語言；那時翻不出來的這裡補翻（翻不出來照樣回，下次打開再試）
	if to := c.Merchant.PreferredLanguage; q.Translated == "" && q.Lang != "" && q.Lang != to && s.Chat != nil {
		translated, err := merchant.TranslateCustomer(r.Context(), s.Chat, q.Lang, to, q.Text)
		if err == nil {
			err = s.DB.SetPendingTranslation(r.Context(), q.ID, translated)
		}
		if err != nil {
			s.Log.Warn("待補問題翻譯失敗", "id", q.ID, "error", err.Error())
		} else {
			q.Translated = translated
		}
	}
	d := pendingDetail{PendingQuestion: q, Hits: []pendingHitView{}}
	for _, h := range q.Hits {
		logs, err := s.DB.LogsAround(r.Context(), h.Platform, h.UserID, h.Time, pendingContext)
		if err != nil {
			s.internalError(w, err)
			return
		}
		v := pendingHitView{PendingHit: h, Conversation: []logLine{}}
		for _, e := range logs {
			v.Conversation = append(v.Conversation, logLine{e.Time, e.Direction, e.Text, e.SendFailed})
		}
		d.Hits = append(d.Hits, v)
	}
	writeJSON(w, http.StatusOK, d)
}

// pendingQuestion 讀網址上 {id} 那筆待補問題；沒有就回 404
func (s *Server) pendingQuestion(w http.ResponseWriter, r *http.Request) (store.PendingQuestion, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, errPendingNotFound)
		return store.PendingQuestion{}, false
	}
	q, err := s.DB.PendingQuestion(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, errPendingNotFound)
		return q, false
	}
	if err != nil {
		s.internalError(w, err)
		return q, false
	}
	return q, true
}

const errPendingNotFound = "找不到這個待補問題"

// pendingAction 是四種動作共用的流程：讀 body，先把待補問題改成 status（resolution 是轉成什麼；待處理、已忽略都可以轉，
// 同時按兩次時第二次回 409，不會新增兩次），再做 apply；apply 回 false（已經回了錯誤）時把狀態改回去。成功回這筆待補問題
func pendingAction[T any](s *Server, status, resolution string, apply func(w http.ResponseWriter, r *http.Request, body T) bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q, ok := s.pendingQuestion(w, r)
		if !ok {
			return
		}
		var body T
		if !decodeBody(w, r, &body) {
			return
		}
		err := errAlreadyResolved
		if q.Status != "resolved" {
			err = s.DB.TransitionPending(r.Context(), q.ID, q.Status, status, resolution)
		}
		switch {
		case errors.Is(err, errAlreadyResolved), errors.Is(err, store.ErrAlreadyHandled):
			writeError(w, http.StatusConflict, errAlreadyResolved.Error())
			return
		case err != nil:
			s.internalError(w, err)
			return
		}
		if !apply(w, r, body) {
			if err := s.DB.TransitionPending(r.Context(), q.ID, status, q.Status, q.Resolution); err != nil {
				s.Log.Error("待補問題的狀態改不回去", "id", q.ID, "error", err.Error())
			}
			return
		}
		s.getPending(w, r)
	}
}

var errAlreadyResolved = errors.New("這個待補問題已經處理過了")

// pendingFAQ 轉成 FAQ：body 是 {question, alt_questions, answer}（後台預填顧客原話、員工回覆），在 FAQ 最後新增一題
func (s *Server) pendingFAQ() http.HandlerFunc {
	return pendingAction(s, "resolved", "faq", func(w http.ResponseWriter, r *http.Request, f merchant.FAQ) bool {
		_, err := merchant.AddFAQ(r.Context(), s.DB.SQL(), f)
		return s.savedOK(w, err)
	})
}

// pendingKnowledge 加進知識段落：body 是 {title, body}，title 是現有段落就接在後面，否則新增一段
func (s *Server) pendingKnowledge() http.HandlerFunc {
	return pendingAction(s, "resolved", "knowledge", func(w http.ResponseWriter, r *http.Request, sec merchant.Section) bool {
		return s.savedOK(w, merchant.AppendKnowledge(r.Context(), s.DB.SQL(), sec.Title, sec.Body))
	})
}

// pendingIgnore 忽略：body 是 {}；之後相同的問題照樣計次，但不在待處理列表
func (s *Server) pendingIgnore() http.HandlerFunc {
	return pendingAction(s, "ignored", "", func(http.ResponseWriter, *http.Request, struct{}) bool { return true })
}

// pendingExample 加成意圖案例：body 是 {text, intent}，和新增案例一樣檢查題庫、案例庫重複
func (s *Server) pendingExample() http.HandlerFunc {
	return pendingAction(s, "resolved", "example", func(w http.ResponseWriter, r *http.Request, req exampleRequest) bool {
		_, ok := s.saveExample(w, r, req)
		return ok
	})
}
