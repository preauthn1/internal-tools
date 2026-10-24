package admin

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// SessionCookie 是 session cookie 的名稱
const SessionCookie = "aichatbot_session"

// Server 是後台的 HTTP 服務：/api/ 是後台 API，其他路徑是 Dir 裡建好的頁面
type Server struct {
	DB      *store.DB
	Dir     string // 建好的後台頁面（web/apps/web-antd/dist）
	EvalDir string // 題庫（eval/）：新增意圖案例時檢查不能和題庫重複
	Model   ModelSwitch
	Chat    model.ChatModel // 翻譯固定回覆用（目前開著的模型）
	Now     func() time.Time
	Log     *slog.Logger

	SelfTester SelfTester   // 分類自我測試；nil 時不跑
	Handoff    HandoffEnder // 後台結束轉人工（和顧客訊息同一個 dialog.Bot，共用顧客鎖）；nil 時回 503

	modelMu  sync.Mutex // 改模型設定時存資料庫和換模型開關要一起做完
	selfTest selfTestRunner
}

func (s *Server) Handler() http.Handler {
	// 要先登入的 API
	api := http.NewServeMux()
	api.HandleFunc("GET /api/auth/me", example.invalid)
	api.HandleFunc("POST /api/auth/logout", s.logout)
	api.HandleFunc("GET /api/merchant", s.getMerchant)
	api.HandleFunc("GET /api/faq", s.getFAQ)
	api.HandleFunc("GET /api/languages", s.getLanguages)
	api.HandleFunc("PUT /api/merchant", s.putMerchant)
	api.HandleFunc("PUT /api/faq", s.putFAQ)
	api.HandleFunc("PUT /api/languages", s.putLanguages)
	api.HandleFunc("GET /api/intents", s.intents)
	api.HandleFunc("GET /api/intent-examples", s.listExamples)
	api.HandleFunc("POST /api/intent-examples", s.addExample)
	api.HandleFunc("DELETE /api/intent-examples/{id}", s.deleteExample)
	api.HandleFunc("/api/intent-examples/{id}", exampleMethodNotAllowed)
	api.HandleFunc("GET /api/custom-intents", s.getCustomIntents)
	api.HandleFunc("PUT /api/custom-intents", s.putCustomIntents)
	api.HandleFunc("GET /api/custom-intents/self-test", s.getSelfTest)
	api.HandleFunc("POST /api/custom-intents/self-test", s.postSelfTest)
	api.HandleFunc("GET /api/fixed-replies", s.getFixedReplies)
	api.HandleFunc("PUT /api/fixed-replies/{key}", s.putFixedReply)
	api.HandleFunc("GET /api/booking-setting", s.getBookingSetting)
	api.HandleFunc("PUT /api/booking-setting", s.putBookingSetting)
	api.HandleFunc("GET /api/handoff-setting", s.getHandoffSetting)
	api.HandleFunc("PUT /api/handoff-setting", s.putHandoffSetting)
	api.HandleFunc("POST /api/handoff-setting/bind-code", s.regenerateBindCode)
	api.HandleFunc("GET /api/pending-questions", s.listPending)
	api.HandleFunc("GET /api/pending-questions/{id}", s.getPending)
	api.HandleFunc("POST /api/pending-questions/{id}/faq", s.pendingFAQ())
	api.HandleFunc("POST /api/pending-questions/{id}/knowledge", s.pendingKnowledge())
	api.HandleFunc("POST /api/pending-questions/{id}/example", s.pendingExample())
	api.HandleFunc("POST /api/pending-questions/{id}/ignore", s.pendingIgnore())
	api.HandleFunc("GET /api/bookings", s.listBookings)
	api.HandleFunc("GET /api/conversations", s.listCustomers)
	api.HandleFunc("GET /api/conversations/{platform}/{userID}", s.getConversation)
	api.HandleFunc("POST /api/conversations/logs/{id}/translate", s.translateLog)
	api.HandleFunc("POST /api/conversations/{platform}/{userID}/end-handoff", s.endHandoff)
	api.HandleFunc("GET /api/usage", s.getUsage)
	api.HandleFunc("GET /api/model", s.getModel)
	api.HandleFunc("PUT /api/model", s.putModel)
	api.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) { writeError(w, http.StatusNotFound, "沒有這個 API") })

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth/login", s.login)
	mux.Handle("/api/", s.requireLogin(api))
	mux.HandleFunc("/", s.pages) // 不能寫成 GET /：跟 /api/ 衝突
	return requireJSON(mux)
}

// pages 給建好的後台頁面：Dir 裡有的檔案照給，其他路徑都給 index.html（前端路由用 history 模式）
func (s *Server) pages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	index, err := os.Open(filepath.Join(s.Dir, "index.html"))
	if err != nil {
		http.Error(w, "後台頁面還沒建置：cd web && corepack pnpm install --frozen-lockfile && corepack pnpm build:antd", http.StatusServiceUnavailable)
		return
	}
	defer index.Close()
	if name := path.Clean(r.URL.Path); name != "/" {
		// http.Dir 擋掉 .. 跑出資料夾的路徑
		if f, err := http.Dir(s.Dir).Open(name); err == nil {
			st, err := f.Stat()
			f.Close()
			if err == nil && !st.IsDir() {
				http.FileServer(http.Dir(s.Dir)).ServeHTTP(w, r)
				return
			}
		}
	}
	st, err := index.Stat()
	if err != nil {
		s.internalError(w, err)
		return
	}
	// index.html 每次都要重新確認，前端更新後才拿得到新的檔名
	w.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(w, r, "index.html", st.ModTime(), index)
}

// requireJSON 防 CSRF：除了 SameSite=Strict，會改資料的請求一律要是 application/json（瀏覽器的表單送不出這種請求）
func requireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions:
		default:
			if mt, _, err := mime.ParseMediaType(r.Header.Get("Content-Type")); err != nil || mt != "application/json" {
				writeError(w, http.StatusUnsupportedMediaType, "要用 application/json")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

type userKey struct{}

// requireLogin：沒有有效的 session 回 401
func (s *Server) requireLogin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie(SessionCookie)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "請先登入")
			return
		}
		u, err := sessionUser(r.Context(), s.DB, c.Value, s.Now())
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, "請先登入")
			return
		}
		if err != nil {
			s.internalError(w, err)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userKey{}, u)))
	})
}

func currentUser(r *http.Request) User { return r.Context().Value(userKey{}).(User) }

// uiLanguages：顧客語言代碼 → 後台介面語言；不在裡面的用 defaultUILanguage
var uiLanguages = map[string]string{"zh": "zh-TW", "en": "en-US", "ja": "ja-JP"}

const defaultUILanguage = "zh-TW"

// me 回目前使用者：帳號和介面語言（跟著商家慣用語言，還沒有商家設定時用繁中）
func (s *Server) me(w http.ResponseWriter, r *http.Request) {
	lang := defaultUILanguage
	c, err := merchant.Load(r.Context(), s.DB.SQL())
	switch {
	case err == nil:
		if l, ok := uiLanguages[c.Merchant.PreferredLanguage]; ok {
			lang = l
		}
	case !errors.Is(err, merchant.ErrNotConfigured):
		s.internalError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"username": currentUser(r).Username, "language": lang})
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<16)).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "格式不正確")
		return
	}
	u, err := authenticate(r.Context(), s.DB, req.Username, req.Password)
	if errors.Is(err, errBadLogin) {
		writeError(w, http.StatusUnauthorized, err.Error())
		return
	}
	if err != nil {
		s.internalError(w, err)
		return
	}
	token, err := newSession(r.Context(), s.DB, u.ID, s.Now())
	if err != nil {
		s.internalError(w, err)
		return
	}
	setSessionCookie(w, token, int(SessionTTL/time.Second))
	writeJSON(w, http.StatusOK, map[string]string{"username": u.Username})
}

// logout 刪掉這個 session、清掉 cookie
func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	c, _ := r.Cookie(SessionCookie) // requireLogin 已經確認過
	if err := deleteSession(r.Context(), s.DB, c.Value); err != nil {
		s.internalError(w, err)
		return
	}
	setSessionCookie(w, "", -1)
	w.WriteHeader(http.StatusNoContent)
}

// setSessionCookie：maxAge 是秒數，負的是清掉。
// 一律設 Secure（Safari 在 http://localhost 會丟掉，本機開發用 Chrome／Firefox）：
// cloudflared 轉進來的正式流量也是從 loopback 連進來，分不出本機開發，不能看連線來源拿掉 Secure
func setSessionCookie(w http.ResponseWriter, token string, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name: SessionCookie, Value: token, Path: "/", MaxAge: maxAge,
		HttpOnly: true, Secure: true, SameSite: http.SameSiteStrictMode,
	})
}

func (s *Server) internalError(w http.ResponseWriter, err error) {
	s.Log.Error("後台 API 失敗", "error", err.Error())
	writeError(w, http.StatusInternalServerError, "伺服器錯誤")
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

// writeError 回 {"error": 訊息}
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
