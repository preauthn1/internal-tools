package admin

import (
	"errors"
	"net/http"

	"aichatbot/internal/merchant"
)

// 自訂意圖的 API：整份 GET／PUT；內建意圖不能改，沒有對應的 API

// customIntentsBody：GET 另外附上知識段落的標題（只根據指定知識回答時選段落用）
type customIntentsBody struct {
	CustomIntents   []merchant.CustomIntent `json:"custom_intents"`
	KnowledgeTitles []string                `json:"knowledge_titles,omitempty"`
}

func (s *Server) getCustomIntents(w http.ResponseWriter, r *http.Request) {
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	titles := []string{}
	for _, k := range c.Knowledge {
		titles = append(titles, k.Title)
	}
	writeJSON(w, http.StatusOK, customIntentsBody{c.CustomIntents, titles})
}

// putCustomIntents 存自訂意圖；例句和題庫重複回 409。存好後在背景跑分類自我測試
func (s *Server) putCustomIntents(w http.ResponseWriter, r *http.Request) {
	var body customIntentsBody
	if !decodeBody(w, r, &body) {
		return
	}
	// 每次都重新讀題庫（同新增案例）
	tests, err := merchant.LoadTestSet(s.EvalDir)
	if err != nil {
		s.Log.Error("題庫讀不到", "dir", s.EvalDir, "error", err.Error())
		writeError(w, http.StatusInternalServerError, "題庫讀不到，無法檢查重複")
		return
	}
	err = merchant.SaveCustomIntents(r.Context(), s.DB.SQL(), tests, body.CustomIntents)
	var overlap *merchant.OverlapError
	if errors.As(err, &overlap) {
		writeError(w, http.StatusConflict, err.Error())
		return
	}
	if err == nil {
		s.startSelfTest()
	}
	s.saved(w, r, err, s.getCustomIntents)
}
