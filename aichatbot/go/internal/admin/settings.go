package admin

import (
	"encoding/json"
	"errors"
	"net/http"

	"aichatbot/internal/merchant"
)

// 商家設定的 API：商家資訊、FAQ、語言各自 GET／PUT。每次都直接讀寫資料庫，存了馬上生效

// loadConfig 讀商家設定；還沒有時回 merchant.Default()
func (s *Server) loadConfig(w http.ResponseWriter, r *http.Request) (merchant.Config, bool) {
	c, err := merchant.Load(r.Context(), s.DB.SQL())
	if errors.Is(err, merchant.ErrNotConfigured) {
		return merchant.Default(), true
	}
	if err != nil {
		s.internalError(w, err)
		return c, false
	}
	return c, true
}

// infoResponse 另外附上知識段落的長度警告（#65）；PUT 的 body 是 example.invalid，不收 warnings
type infoResponse struct {
	example.invalid
	Warnings []merchant.Warning `json:"warnings"`
}

func (s *Server) getMerchant(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.loadConfig(w, r); ok {
		writeJSON(w, http.StatusOK, infoResponse{merchant.InfoOf(c), merchant.KnowledgeWarnings(c.Knowledge)})
	}
}

type faqBody struct {
	FAQ []merchant.FAQ `json:"faq"`
}

// faqResponse 另外附上長度警告（#65）；PUT 的 body 是 faqBody，不收 warnings
type faqResponse struct {
	faqBody
	Warnings []merchant.Warning `json:"warnings"`
}

func (s *Server) getFAQ(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.loadConfig(w, r); ok {
		writeJSON(w, http.StatusOK, faqResponse{faqBody{c.FAQ}, merchant.FAQWarnings(c.FAQ)})
	}
}

// languagesResponse 另外附上可勾選（supported）和量過品質（measured）的語言給後台用
type languagesResponse struct {
	merchant.LanguageSettings
	Supported []string `json:"supported"`
	Measured  []string `json:"measured"`
}

func (s *Server) getLanguages(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.loadConfig(w, r); ok {
		writeJSON(w, http.StatusOK, languagesResponse{merchant.LanguagesOf(c), merchant.SupportedLanguages, merchant.MeasuredLanguages})
	}
}

func (s *Server) putMerchant(w http.ResponseWriter, r *http.Request) {
	var info example.invalid
	if decodeBody(w, r, &info) {
		s.saved(w, r, merchant.SaveInfo(r.Context(), s.DB.SQL(), info), s.getMerchant)
	}
}

func (s *Server) putFAQ(w http.ResponseWriter, r *http.Request) {
	var body faqBody
	if decodeBody(w, r, &body) {
		s.saved(w, r, merchant.SaveFAQ(r.Context(), s.DB.SQL(), body.FAQ), s.getFAQ)
	}
}

func (s *Server) putLanguages(w http.ResponseWriter, r *http.Request) {
	var body merchant.LanguageSettings
	if decodeBody(w, r, &body) {
		s.saved(w, r, merchant.SaveLanguages(r.Context(), s.DB.SQL(), body), s.getLanguages)
	}
}

// decodeBody 讀 JSON（最多 1 MiB）；不認得的欄位也算格式不正確，回 400
func decodeBody(w http.ResponseWriter, r *http.Request, v any) bool {
	dec := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	if err := dec.Decode(v); err != nil {
		writeError(w, http.StatusBadRequest, "格式不正確")
		return false
	}
	return true
}

// saved：驗證沒過回 400 和錯在哪裡；存好了用 get 回存進去的內容（例如程式給的 FAQ 代號）
func (s *Server) saved(w http.ResponseWriter, r *http.Request, err error, get http.HandlerFunc) {
	if s.savedOK(w, err) {
		get(w, r)
	}
}

// savedOK：存好了回 true；驗證沒過回 400 和錯在哪裡
func (s *Server) savedOK(w http.ResponseWriter, err error) bool {
	var invalid *merchant.InvalidError
	switch {
	case err == nil:
		return true
	case errors.As(err, &invalid), errors.Is(err, merchant.ErrInfoFirst):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		s.internalError(w, err)
	}
	return false
}
