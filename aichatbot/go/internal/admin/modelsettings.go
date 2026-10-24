package admin

import (
	"net/http"

	"aichatbot/internal/model"
)

// ModelSwitch 是模型開關（model.Switch）：後台讀目前的設定、改了馬上生效
type ModelSwitch interface {
	Settings() (model string, budget float64)
	HaikuAvailable() bool
	Apply(model string, budget float64) error
}

// modelSettings 是模型設定的 API：只附有沒有設定 API key，不回 key 本身
type modelSettings struct {
	Model          string  `json:"model"`
	MonthlyBudget  float64 `json:"monthly_budget_usd"`
	HaikuAvailable bool    `json:"haiku_available"`
}

func (s *Server) getModel(w http.ResponseWriter, r *http.Request) {
	m, budget := s.Model.Settings()
	writeJSON(w, http.StatusOK, modelSettings{m, budget, s.Model.HaikuAvailable()})
}

// putModel 先檢查、存進資料庫，存好了才換掉模型開關（存不進去時兩邊都不變）。
// 同時有兩個請求時一個一個來，資料庫和模型開關最後一定是同一個設定
func (s *Server) putModel(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Model         string  `json:"model"`
		MonthlyBudget float64 `json:"monthly_budget_usd"`
	}
	if !decodeBody(w, r, &body) {
		return
	}
	if err := model.CheckSetting(body.Model, body.MonthlyBudget, s.Model.HaikuAvailable()); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.modelMu.Lock()
	defer s.modelMu.Unlock()
	if err := s.DB.SaveModelSetting(r.Context(), body.Model, body.MonthlyBudget); err != nil {
		s.internalError(w, err)
		return
	}
	if err := s.Model.Apply(body.Model, body.MonthlyBudget); err != nil {
		s.internalError(w, err)
		return
	}
	s.getModel(w, r)
}
