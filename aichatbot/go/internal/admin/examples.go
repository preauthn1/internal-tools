package admin

import (
	"encoding/json"
	"errors"
	"net/http"
	"slices"
	"strconv"

	"aichatbot/internal/merchant"
)

// intents 回意圖代號和三語白話名稱（後台不顯示代號）；自訂意圖接在內建意圖後面，名稱只有商家取的一種
func (s *Server) intents(w http.ResponseWriter, r *http.Request) {
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	list := slices.Clone(merchant.Intents)
	for _, ci := range c.CustomIntents {
		list = append(list, merchant.Intent{Code: ci.Code, Labels: map[string]string{"zh-TW": ci.Name, "en-US": ci.Name, "ja-JP": ci.Name}})
	}
	writeJSON(w, http.StatusOK, list)
}

// listExamples 列出意圖案例庫
func (s *Server) listExamples(w http.ResponseWriter, r *http.Request) {
	examples, err := merchant.ListExamples(r.Context(), s.DB.SQL())
	if err != nil {
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, examples)
}

// exampleRequest 是新增案例的 body
type exampleRequest struct {
	Text   string `json:"text"`
	Intent string `json:"intent"`
}

// addExample 新增商家的案例：body 是 {text, intent}
func (s *Server) addExample(w http.ResponseWriter, r *http.Request) {
	var req exampleRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "格式不正確")
		return
	}
	if e, ok := s.saveExample(w, r, req); ok {
		writeJSON(w, http.StatusCreated, e)
	}
}

// saveExample 新增一句商家的案例（檢查和題庫、案例庫重複）；存不進去時回錯誤，ok 是 false
func (s *Server) saveExample(w http.ResponseWriter, r *http.Request, req exampleRequest) (merchant.StoredExample, bool) {
	// 每次都重新讀題庫：題庫很小，改了題庫也不用重啟
	tests, err := merchant.LoadTestSet(s.EvalDir)
	if err != nil {
		s.Log.Error("題庫讀不到", "dir", s.EvalDir, "error", err.Error())
		writeError(w, http.StatusInternalServerError, "題庫讀不到，無法檢查重複")
		return merchant.StoredExample{}, false
	}
	e, err := merchant.AddExample(r.Context(), s.DB.SQL(), tests, req.Text, req.Intent)
	var overlap *merchant.OverlapError
	switch {
	case err == nil:
		return e, true
	case errors.Is(err, merchant.ErrInvalidExample):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.As(err, &overlap), errors.Is(err, merchant.ErrDuplicateExample):
		writeError(w, http.StatusConflict, err.Error())
	default:
		s.internalError(w, err)
	}
	return e, false
}

// deleteExample 刪除商家的案例
func (s *Server) deleteExample(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil {
		writeError(w, http.StatusNotFound, merchant.ErrExampleNotFound.Error())
		return
	}
	switch err := merchant.DeleteExample(r.Context(), s.DB.SQL(), id); {
	case err == nil:
		w.WriteHeader(http.StatusNoContent)
	case errors.Is(err, merchant.ErrExampleNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, merchant.ErrBuiltinExample):
		writeError(w, http.StatusForbidden, err.Error())
	default:
		s.internalError(w, err)
	}
}

// exampleMethodNotAllowed：案例不能修改（要改就刪掉再新增），內建案例也不能改；
// 要明確註冊，不然會被 /api/ 接走回 404
func exampleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Allow", "DELETE")
	writeError(w, http.StatusMethodNotAllowed, "案例不能修改，請刪掉再新增")
}
