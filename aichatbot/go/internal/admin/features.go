package admin

import (
	"fmt"
	"net/http"
	"regexp"

	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

// 預約設定、轉人工設定的 API。關閉時回的文字是固定回覆 booking_off、handoff_off，在固定回覆的 API 改

// bookingSetting 是預約設定頁的內容：總開關是功能開關 booking；人數開關和自訂欄位是預約欄位（#67）。
// 送出時 people_enabled、fields 都沒帶（#67 以前的前端）就不改預約欄位
type bookingSetting struct {
	Enabled       bool                     `json:"enabled"`
	PeopleEnabled *bool                    `json:"people_enabled"`
	Fields        *[]merchant.BookingField `json:"fields"`
}

func (s *Server) getBookingSetting(w http.ResponseWriter, r *http.Request) {
	if c, ok := s.loadConfig(w, r); ok {
		writeJSON(w, http.StatusOK, bookingSettingOf(c))
	}
}

func bookingSettingOf(c merchant.Config) bookingSetting {
	b := merchant.EffectiveBookingFields(c)
	return bookingSetting{merchant.FeatureEnabled(c, "booking"), &b.PeopleEnabled, &b.Custom}
}

// handoffSetting 是轉人工設定頁的內容：總開關是功能開關 handoff；提示文字用固定回覆 handoff（不開放改 prompt_text）
type handoffSetting struct {
	Enabled        bool                 `json:"enabled"`
	OnNoAnswer     bool                 `json:"on_no_answer"`
	StreakLimit    int                  `json:"streak_limit"`
	ServiceHours   []store.ServiceHours `json:"service_hours"`
	AutoEndMinutes int                  `json:"auto_end_minutes"`
}

// clockPattern 是服務時段的時間：24 小時制、兩位數的 HH:MM
var clockPattern = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

// validate 檢查每個欄位：累計次數 1～10、自動結束 0～1440 分鐘；
// 服務時段的星期 0～6、時間 HH:MM、開始和結束不同（開始晚於結束是跨過午夜）、時段之間不能重疊（跨日的也算）
func (h handoffSetting) validate() error {
	var errs []string
	if h.StreakLimit < 1 || h.StreakLimit > 10 {
		errs = append(errs, "累計次數要是 1～10")
	}
	if h.AutoEndMinutes < 0 || h.AutoEndMinutes > 1440 {
		errs = append(errs, "自動結束要是 0～1440 分鐘（0 是不自動結束）")
	}
	valid := make([]bool, len(h.ServiceHours))
	for i, x := range h.ServiceHours {
		n := len(errs)
		if x.Weekday < 0 || x.Weekday > 6 {
			errs = append(errs, fmt.Sprintf("第 %d 個服務時段的星期要是 0～6", i+1))
		}
		for _, c := range []struct{ name, v string }{{"開始", x.Start}, {"結束", x.End}} {
			if !clockPattern.MatchString(c.v) {
				errs = append(errs, fmt.Sprintf("第 %d 個服務時段的%s時間要是 HH:MM", i+1, c.name))
			}
		}
		if len(errs) == n && x.Start == x.End { // 開始晚於結束是跨過午夜（對話層的 inServiceHours 也這樣算）
			errs = append(errs, fmt.Sprintf("第 %d 個服務時段的開始和結束不能相同", i+1))
		}
		valid[i] = len(errs) == n
		for j := range i {
			if valid[i] && valid[j] && overlaps(weekSpan(x), weekSpan(h.ServiceHours[j])) {
				errs = append(errs, fmt.Sprintf("第 %d 個服務時段和第 %d 個重疊", i+1, j+1))
			}
		}
	}
	if len(errs) > 0 {
		return &merchant.InvalidError{Problems: errs}
	}
	return nil
}

const weekMinutes = 7 * 24 * 60

// weekSpan 把服務時段換成一週裡的分鐘區間 [開始, 結束)，跨過午夜的結束算到隔天（可能超過一週，overlaps 會處理）
func weekSpan(x store.ServiceHours) [2]int {
	start := x.Weekday*24*60 + clockMinutes(x.Start)
	end := x.Weekday*24*60 + clockMinutes(x.End)
	if end <= start {
		end += 24 * 60
	}
	return [2]int{start, end}
}

// overlaps 比較兩個一週裡的區間，週六跨到週日的時段要和下一週的開頭比
func overlaps(a, b [2]int) bool {
	for _, shift := range []int{-weekMinutes, 0, weekMinutes} {
		if a[0] < b[1]+shift && b[0]+shift < a[1] {
			return true
		}
	}
	return false
}

// clockMinutes 把已經檢查過的 HH:MM 換成從 0 點起算的分鐘數
func clockMinutes(s string) int {
	return int(s[0]-'0')*600 + int(s[1]-'0')*60 + int(s[3]-'0')*10 + int(s[4]-'0')
}

// handoffSettingResponse 另外附上只讀的綁定碼和有沒有綁好員工群組（不給 chat id）
type handoffSettingResponse struct {
	handoffSetting
	BindCode string `json:"bind_code"`
	Bound    bool   `json:"bound"`
}

func (s *Server) getHandoffSetting(w http.ResponseWriter, r *http.Request) {
	c, ok := s.loadConfig(w, r)
	if !ok {
		return
	}
	h, err := s.DB.HandoffSetting(r.Context())
	if err != nil {
		s.internalError(w, err)
		return
	}
	hours := h.ServiceHours
	if hours == nil {
		hours = []store.ServiceHours{}
	}
	writeJSON(w, http.StatusOK, handoffSettingResponse{
		handoffSetting{merchant.FeatureEnabled(c, "handoff"), h.OnNoAnswer, h.StreakLimit, hours, h.AutoEndMinutes},
		h.BindCode, h.StaffChatID != 0,
	})
}

// putHandoffSetting 存總開關（功能開關 handoff）和 handoff_setting；提示文字、綁定碼、綁定的群組不動。
// 兩張表不在同一個 transaction：先存功能開關（還沒有商家設定時在這一步回 400，什麼都沒存），再存 handoff_setting
func (s *Server) putHandoffSetting(w http.ResponseWriter, r *http.Request) {
	var body handoffSetting
	if !decodeBody(w, r, &body) {
		return
	}
	if err := body.validate(); err != nil {
		s.saved(w, r, err, nil)
		return
	}
	if err := merchant.SaveFeature(r.Context(), s.DB.SQL(), "handoff", body.Enabled); err != nil {
		s.saved(w, r, err, nil)
		return
	}
	h, err := s.DB.HandoffSetting(r.Context())
	if err == nil {
		h.OnNoAnswer, h.StreakLimit, h.ServiceHours, h.AutoEndMinutes = body.OnNoAnswer, body.StreakLimit, body.ServiceHours, body.AutoEndMinutes
		err = s.DB.SaveHandoffSetting(r.Context(), h)
	}
	s.saved(w, r, err, s.getHandoffSetting)
}

// regenerateBindCode 換一組新的綁定碼（舊碼不能再用），回新的轉人工設定；已經綁定的群組不變
func (s *Server) regenerateBindCode(w http.ResponseWriter, r *http.Request) {
	if _, err := s.DB.RegenerateBindCode(r.Context()); err != nil {
		s.internalError(w, err)
		return
	}
	s.getHandoffSetting(w, r)
}

// putBookingSetting 存總開關和預約欄位（沒帶的一邊照舊）。不在同一個 transaction：先存預約欄位（沒通過檢查時什麼都沒存），再存總開關
func (s *Server) putBookingSetting(w http.ResponseWriter, r *http.Request) {
	var body bookingSetting
	if !decodeBody(w, r, &body) {
		return
	}
	if body.PeopleEnabled != nil || body.Fields != nil {
		c, ok := s.loadConfig(w, r)
		if !ok {
			return
		}
		b := merchant.EffectiveBookingFields(c)
		if body.PeopleEnabled != nil {
			b.PeopleEnabled = *body.PeopleEnabled
		}
		if body.Fields != nil {
			b.Custom = *body.Fields
		}
		if err := merchant.SaveBookingFields(r.Context(), s.DB.SQL(), b); err != nil {
			s.saved(w, r, err, nil)
			return
		}
	}
	s.saved(w, r, merchant.SaveFeature(r.Context(), s.DB.SQL(), "booking", body.Enabled), s.getBookingSetting)
}
