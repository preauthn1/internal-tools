package admin

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

func openDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

// 帳號 3～32 個英數字（可以用 _、-）、密碼至少 10 個字；帳號重複時失敗，不覆蓋原本的
func TestCreateUser(t *testing.T) {
	ctx := context.Background()
	db := openDB(t)
	if err := CreateUser(ctx, db, "shop_owner-1", "correct horse"); err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct{ username, password string }{
		{"ab", "correct horse"},
		{"這是中文", "correct horse"},
		{"has space", "correct horse"},
		{"a23456789012345678901234567890123", "correct horse"},
		{"owner2", "short"},
		{"owner2", "密碼只有九個字喔喔"},
	} {
		if err := CreateUser(ctx, db, c.username, c.password); err == nil {
			t.Errorf("%q／%q 應該不能建立", c.username, c.password)
		}
	}
	if err := CreateUser(ctx, db, "owner2", "密碼剛好有十個字喔喔"); err != nil {
		t.Errorf("十個字的密碼：%v", err)
	}
	if err := CreateUser(ctx, db, "shop_owner-1", "another password"); !errors.Is(err, ErrUserExists) {
		t.Errorf("重複的帳號：got %v, want ErrUserExists", err)
	}
}

// fixture：資料庫裡有一個管理員 owner，時間可以往後推
type fixture struct {
	t   *testing.T
	db  *store.DB
	srv *Server
	h   http.Handler
	now time.Time
}

const password = [REDACTED] horse"

func newFixture(t *testing.T) *fixture {
	t.Helper()
	f := &fixture{t: t, db: openDB(t), now: time.Date(2026, 10, 9, 12, 0, 0, 0, time.UTC)}
	if err := CreateUser(context.Background(), f.db, "owner", password); err != nil {
		t.Fatal(err)
	}
	f.srv = &Server{DB: f.db, Dir: t.TempDir(), Now: func() time.Time { return f.now }, Log: slog.New(slog.DiscardHandler)}
	f.h = f.srv.Handler()
	return f
}

// do 送一個請求；body 不是空字串時帶 Content-Type: application/json
func (f *fixture) do(method, path, body string, cookies ...*http.Cookie) *httptest.ResponseRecorder {
	f.t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	for _, c := range cookies {
		req.AddCookie(c)
	}
	w := httptest.NewRecorder()
	f.h.ServeHTTP(w, req)
	return w
}

func loginBody(username, password string) string {
	b, _ := json.Marshal(map[string]string{"username": username, "password": password})
	return string(b)
}

// login 登入成功，回傳 session cookie
func (f *fixture) login() *http.Cookie {
	f.t.Helper()
	w := f.do("POST", "/api/auth/login", loginBody("owner", password))
	if w.Code != http.StatusOK {
		f.t.Fatalf("登入：%d %s", w.Code, w.Body)
	}
	for _, c := range w.Result().Cookies() {
		if c.Name == SessionCookie {
			return c
		}
	}
	f.t.Fatal("登入後沒有 session cookie")
	return nil
}

// 登入成功：回帳號，session cookie 設 HttpOnly、Secure、SameSite=Strict、Path=/，7 天
func TestLogin(t *testing.T) {
	f := newFixture(t)
	w := f.do("POST", "/api/auth/login", loginBody("owner", password))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `"username":"owner"`) {
		t.Fatalf("登入：%d %s", w.Code, w.Body)
	}
	cookies := w.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("cookie：%v", cookies)
	}
	c := cookies[0]
	if c.Name != SessionCookie || len(c.Value) < 40 || !c.HttpOnly || !c.Secure || c.SameSite != http.SameSiteStrictMode || c.Path != "/" || c.MaxAge != 7*24*60*60 {
		t.Errorf("cookie 屬性不對：%+v", c)
	}
}

// 密碼錯、帳號不存在：一律 401「帳號或密碼錯誤」，不給 cookie
func TestLoginFails(t *testing.T) {
	f := newFixture(t)
	for _, body := range []string{loginBody("owner", "wrong password"), loginBody("nobody", password)} {
		w := f.do("POST", "/api/auth/login", body)
		if w.Code != http.StatusUnauthorized || !strings.Contains(w.Body.String(), "帳號或密碼錯誤") || len(w.Result().Cookies()) != 0 {
			t.Errorf("%s：%d %s %v", body, w.Code, w.Body, w.Result().Cookies())
		}
	}
}

// 目前使用者：沒登入（沒有 cookie、cookie 是亂填的）回 401；登入後回帳號和介面語言（還沒有商家設定時是繁中）
func TestMe(t *testing.T) {
	f := newFixture(t)
	for _, cookies := range [][]*http.Cookie{nil, {{Name: SessionCookie, Value: "made-up"}}} {
		if w := f.do("GET", "/api/auth/me", "", cookies...); w.Code != http.StatusUnauthorized {
			t.Errorf("沒登入：%d %s", w.Code, w.Body)
		}
	}
	w := f.do("GET", "/api/auth/me", "", f.login())
	if w.Code != http.StatusOK {
		t.Fatalf("登入後：%d %s", w.Code, w.Body)
	}
	var me struct{ Username, Language string }
	json.Unmarshal(w.Body.Bytes(), &me)
	if me.Username != "owner" || me.Language != "zh-TW" {
		t.Errorf("目前使用者：%s", w.Body)
	}
}

// 有商家設定時，介面語言跟著商家慣用語言（顧客語言代碼 zh、en、ja 對到介面的 zh-TW、en-US、ja-JP）
func TestMeFollowsPreferredLanguage(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	for _, c := range []struct{ preferred, want string }{{"ja", "ja-JP"}, {"en", "en-US"}, {"zh", "zh-TW"}} {
		cfg := merchant.Config{
			Merchant:  merchant.Merchant{Name: "森日茶屋", Intro: "茶館", Contact: "02-1234-5678", PreferredLanguage: c.preferred},
			Languages: []merchant.Language{{Code: c.preferred, Enabled: true}},
		}
		if _, err := merchant.Save(context.Background(), f.db.SQL(), cfg); err != nil {
			t.Fatal(err)
		}
		w := f.do("GET", "/api/auth/me", "", cookie)
		if !strings.Contains(w.Body.String(), `"language":"`+c.want+`"`) {
			t.Errorf("慣用語言 %s：%s", c.preferred, w.Body)
		}
	}
}

// 登出：清掉 cookie，原本的 session 之後就不能用
func TestLogout(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	w := f.do("POST", "/api/auth/logout", "{}", cookie)
	if w.Code != http.StatusNoContent {
		t.Fatalf("登出：%d %s", w.Code, w.Body)
	}
	cleared := w.Result().Cookies()
	if len(cleared) != 1 || cleared[0].Name != SessionCookie || cleared[0].MaxAge >= 0 {
		t.Errorf("登出沒清 cookie：%v", cleared)
	}
	if w := f.do("GET", "/api/auth/me", "", cookie); w.Code != http.StatusUnauthorized {
		t.Errorf("登出後還能用：%d", w.Code)
	}
}

// session 從登入算起 7 天過期，中間有使用也不延長
func TestSessionExpires(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	f.now = f.now.Add(SessionTTL - time.Second)
	if w := f.do("GET", "/api/auth/me", "", cookie); w.Code != http.StatusOK {
		t.Fatalf("快到期前：%d", w.Code)
	}
	f.now = f.now.Add(time.Second)
	if w := f.do("GET", "/api/auth/me", "", cookie); w.Code != http.StatusUnauthorized {
		t.Errorf("到期後：%d", w.Code)
	}
}

// 防 CSRF：會改資料的請求不是 application/json 就回 415（登入也一樣），不會登入、不會登出
func TestWritesRequireJSON(t *testing.T) {
	f := newFixture(t)
	cookie := f.login()
	for _, c := range []struct{ path, contentType string }{
		{"/api/auth/login", "application/x-www-form-urlencoded"},
		{"/api/auth/login", "text/plain"},
		{"/api/auth/logout", ""},
		{"/api/auth/logout", "multipart/form-data; boundary=x"},
	} {
		req := httptest.NewRequest("POST", c.path, strings.NewReader(loginBody("owner", password)))
		if c.contentType != "" {
			req.Header.Set("Content-Type", c.contentType)
		}
		req.AddCookie(cookie)
		w := httptest.NewRecorder()
		f.h.ServeHTTP(w, req)
		if w.Code != http.StatusUnsupportedMediaType || len(w.Result().Cookies()) != 0 {
			t.Errorf("%s %q：%d %v", c.path, c.contentType, w.Code, w.Result().Cookies())
		}
	}
	if w := f.do("GET", "/api/auth/me", "", cookie); w.Code != http.StatusOK {
		t.Errorf("被擋下的登出不該生效：%d", w.Code)
	}
	// charset 參數可以有
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(loginBody("owner", password)))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	w := httptest.NewRecorder()
	f.h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("application/json; charset=utf-8：%d", w.Code)
	}
}

// 沒有的 API：沒登入 401，登入後 404，都是 JSON（不會掉到後台頁面）
func TestUnknownAPI(t *testing.T) {
	f := newFixture(t)
	if w := f.do("GET", "/api/nothing", ""); w.Code != http.StatusUnauthorized {
		t.Errorf("沒登入：%d", w.Code)
	}
	w := f.do("GET", "/api/nothing", "", f.login())
	if w.Code != http.StatusNotFound || !strings.HasPrefix(w.Header().Get("Content-Type"), "application/json") {
		t.Errorf("登入後：%d %s %s", w.Code, w.Header().Get("Content-Type"), w.Body)
	}
}

// 後台頁面不用登入：有的檔案照給，其他非 API 路徑都給 index.html（前端路由用 history 模式），index.html 不快取
func TestServesPages(t *testing.T) {
	f := newFixture(t)
	os.MkdirAll(filepath.Join(f.srv.Dir, "js"), 0o755)
	os.WriteFile(filepath.Join(f.srv.Dir, "index.html"), []byte("<html>後台</html>"), 0o644)
	os.WriteFile(filepath.Join(f.srv.Dir, "js", "app.js"), []byte("console.log(1)"), 0o644)
	for _, c := range []struct{ path, want string }{
		{"/", "<html>後台</html>"},
		{"/js/app.js", "console.log(1)"},
		{"/dashboard", "<html>後台</html>"},
		{"/auth/login?redirect=%2Fdashboard", "<html>後台</html>"},
	} {
		w := f.do("GET", c.path, "")
		if w.Code != http.StatusOK || w.Body.String() != c.want {
			t.Errorf("%s：%d %q", c.path, w.Code, w.Body)
		}
		if c.want == "<html>後台</html>" && w.Header().Get("Cache-Control") != "no-cache" {
			t.Errorf("%s：index.html 要設 no-cache，got %q", c.path, w.Header().Get("Cache-Control"))
		}
	}
}

// 還沒建置前端：頁面回 503 和建置的指令
func TestPagesNotBuilt(t *testing.T) {
	f := newFixture(t)
	f.srv.Dir = filepath.Join(t.TempDir(), "missing")
	w := f.do("GET", "/", "")
	if w.Code != http.StatusServiceUnavailable || !strings.Contains(w.Body.String(), "corepack pnpm build:antd") {
		t.Errorf("%d %s", w.Code, w.Body)
	}
}

// 設定：預設只聽本機 8081、頁面在 web/apps/web-antd/dist、題庫在 eval；BOT_ADMIN_PORT、BOT_ADMIN_DIR、BOT_EVAL_DIR 可以改
func TestConfigFromEnv(t *testing.T) {
	env := map[string]string{}
	getenv := func(k string) string { return env[k] }
	c, err := ConfigFromEnv(getenv)
	if err != nil || c != (Config{Host: "0.0.0.0", Port: 8081, Dir: "web/apps/web-antd/dist", EvalDir: "eval"}) {
		t.Errorf("預設：%+v %v", c, err)
	}
	env["BOT_ADMIN_PORT"], env["BOT_ADMIN_DIR"], env["BOT_EVAL_DIR"] = "9000", "/srv/admin", "/srv/eval"
	if c, err := ConfigFromEnv(getenv); err != nil || c.Port != 9000 || c.Dir != "/srv/admin" || c.EvalDir != "/srv/eval" || c.Host != "0.0.0.0" {
		t.Errorf("有設：%+v %v", c, err)
	}
	env["BOT_ADMIN_PORT"] = "abc"
	if _, err := ConfigFromEnv(getenv); err == nil {
		t.Error("BOT_ADMIN_PORT 不是數字應該失敗")
	}
}
