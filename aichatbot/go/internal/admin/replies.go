package admin

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"aichatbot/internal/merchant"
)

// 固定回覆的 API：只開放 merchant.EditableReplies 那 5 則

type fixedReply struct {
	Key    string              `json:"key"`
	Labels map[string]string   `json:"labels"`
	Texts  map[string][]string `json:"texts"`
}

type fixedRepliesResponse struct {
	Preferred string       `json:"preferred_language"`
	Languages []string     `json:"languages"` // 啟用的語言，照設定的順序
	Replies   []fixedReply `json:"replies"`
}

func (s *Server) getFixedReplies(w http.ResponseWriter, r *http.Request) {
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	resp := fixedRepliesResponse{Preferred: c.Merchant.PreferredLanguage, Languages: enabledLanguages(c), Replies: []fixedReply{}}
	for _, b := range merchant.RepliesOf(c) {
		e, _ := merchant.EditableReplyByKey(b.Key)
		resp.Replies = append(resp.Replies, fixedReply{b.Key, e.Labels, b.Texts})
	}
	writeJSON(w, http.StatusOK, resp)
}

// putFixedReply 存一則固定回覆。慣用語言版本有改時，每一句用 AI 翻成其他啟用的語言（取代送來的）；
// 沒改時照送來的存（商家修正翻譯），只有空著的語言才翻。
// 翻譯在 transaction 外面先做好（AI 很慢，不能鎖住只有一條連線的資料庫）；任何一句翻不出來就整次不存，回 502
func (s *Server) putFixedReply(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Texts map[string][]string `json:"texts"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	key := r.PathValue("key")
	if _, ok := merchant.EditableReplyByKey(key); !ok {
		writeError(w, http.StatusNotFound, "沒有這則可以改的固定回覆："+key)
		return
	}
	c, err := merchant.Load(r.Context(), s.DB.SQL())
	if errors.Is(err, merchant.ErrNotConfigured) {
		err = merchant.ErrInfoFirst
	}
	if err == nil {
		// 翻譯前先檢查送來的內容，沒過的不必等 AI
		err = merchant.CheckReply(c, key, body.Texts)
	}
	if err != nil {
		s.saved(w, r, err, nil)
		return
	}
	var current merchant.ReplyBlock
	for _, b := range merchant.RepliesOf(c) {
		if b.Key == key {
			current = b
		}
	}
	pref := c.Merchant.PreferredLanguage
	changed := !slices.Equal(trimAll(body.Texts[pref]), current.Texts[pref])
	texts := map[string][]string{pref: body.Texts[pref]}
	for _, lang := range enabledLanguages(c) {
		if lang == pref {
			continue
		}
		if !changed && len(body.Texts[lang]) > 0 {
			texts[lang] = body.Texts[lang]
			continue
		}
		var out []string
		for _, sentence := range trimAll(body.Texts[pref]) {
			t, err := merchant.Translate(r.Context(), s.Chat, pref, lang, sentence)
			if err != nil {
				writeError(w, http.StatusBadGateway, err.Error())
				return
			}
			out = append(out, t)
		}
		texts[lang] = out
	}
	s.saved(w, r, merchant.SaveReply(r.Context(), s.DB.SQL(), key, texts), s.getFixedReplies)
}

func trimAll(list []string) []string {
	out := make([]string, len(list))
	for i, s := range list {
		out[i] = strings.TrimSpace(s)
	}
	return out
}

func enabledLanguages(c merchant.Config) []string {
	out := []string{}
	for _, l := range c.Languages {
		if l.Enabled {
			out = append(out, l.Code)
		}
	}
	return out
}
