package admin

// 後台設定和檔案匯入的比對（#39）：從空的資料庫只透過後台 API 做出一份商家設定，
// 匯出後應該和直接匯入（merchant.Save）的相同；go/cmd/adminparity 用它產生設定檔給 eval/snapshot.py 逐題比對

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// BuildViaAPI 在空的資料庫 db 上，照後台頁面的操作順序打 API 做出 cfg：商家資訊、語言、FAQ、預約和轉人工設定、
// 可以改的固定回覆（先送慣用語言讓後台翻譯，再逐語言修正成 cfg 的文字）、自訂意圖、商家新增的意圖案例。
// 產品內建的部分（內建案例、不開放改的固定回覆）商家改不了，存好商家資訊後直接寫進資料庫。
// 翻譯用照抄原文的假模型（第二次送出會換成 cfg 的文字），不需要本機模型；evalDir 是題庫（新增案例時檢查重複）
func BuildViaAPI(ctx context.Context, db *store.DB, evalDir string, cfg merchant.Config) error {
	const username, password = [REDACTED]", "parity-password"
	if err := CreateUser(ctx, db, username, password); err != nil {
		return err
	}
	srv := &Server{DB: db, EvalDir: evalDir, Chat: echoModel{}, Now: time.Now, Log: slog.New(slog.DiscardHandler)}
	c := &apiClient{ctx: ctx, h: srv.Handler()}
	if err := c.login(username, password); err != nil {
		return err
	}
	m := cfg.Merchant
	if err := c.do("PUT", "/api/merchant", example.invalid{Name: m.Name, Intro: m.Intro, Contact: m.Contact, Knowledge: cfg.Knowledge}, nil); err != nil {
		return err
	}
	if err := writeBuiltins(ctx, db, cfg); err != nil {
		return err
	}
	if err := c.do("PUT", "/api/languages", merchant.LanguageSettings{Preferred: m.PreferredLanguage, Languages: cfg.Languages}, nil); err != nil {
		return err
	}
	// FAQ 代號照 cfg 送：後台前端存檔時也會帶回原本的代號
	if err := c.do("PUT", "/api/faq", faqBody{cfg.FAQ}, nil); err != nil {
		return err
	}
	if err := c.do("PUT", "/api/booking-setting", bookingSettingOf(cfg), nil); err != nil {
		return err
	}
	var handoff handoffSetting
	if err := c.do("GET", "/api/handoff-setting", nil, &handoff); err != nil {
		return err
	}
	handoff.Enabled = merchant.FeatureEnabled(cfg, "handoff")
	if err := c.do("PUT", "/api/handoff-setting", handoff, nil); err != nil {
		return err
	}
	for _, b := range editableRepliesIn(cfg) {
		path := "/api/fixed-replies/" + b.Key
		if err := c.do("PUT", path, map[string]any{"texts": map[string][]string{m.PreferredLanguage: b.Texts[m.PreferredLanguage]}}, nil); err != nil {
			return err
		}
		if err := c.do("PUT", path, map[string]any{"texts": b.Texts}, nil); err != nil {
			return err
		}
	}
	if err := c.do("PUT", "/api/custom-intents", customIntentsBody{CustomIntents: cfg.CustomIntents}, nil); err != nil {
		return err
	}
	for _, e := range cfg.Examples {
		if !e.Builtin {
			if err := c.do("POST", "/api/intent-examples", exampleRequest{Text: e.Text, Intent: e.Intent}, nil); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeBuiltins 把 cfg 的內建案例和不開放改的固定回覆寫進資料庫，商家設定的部分照資料庫裡的寫回。
// 用 merchant.Save 整份寫回，所以要在第一次存商家資訊之後、其他 API 之前呼叫，不會蓋掉、也不會跳過後台存的設定
func writeBuiltins(ctx context.Context, db *store.DB, cfg merchant.Config) error {
	c, err := merchant.Load(ctx, db.SQL())
	if err != nil {
		return err
	}
	for _, e := range cfg.Examples {
		if e.Builtin {
			c.Examples = append(c.Examples, e)
		}
	}
	for _, r := range cfg.Replies {
		if _, ok := merchant.EditableReplyByKey(r.Key); !ok {
			c.Replies = append(c.Replies, r)
		}
	}
	_, err = merchant.Save(ctx, db.SQL(), c)
	return err
}

// editableRepliesIn 是 cfg 裡存的、後台可以改的固定回覆（照 merchant.EditableReplies 的順序），語言 → 句子。
// 不用 merchant.RepliesOf：它會補預設文字、把 * 展開成每種語言，送出去就不是 cfg 原本的內容
func editableRepliesIn(cfg merchant.Config) []merchant.ReplyBlock {
	var out []merchant.ReplyBlock
	for _, e := range merchant.EditableReplies {
		b := merchant.ReplyBlock{Key: e.Key, Texts: map[string][]string{}}
		for _, r := range cfg.Replies {
			if r.Key == e.Key {
				b.Texts[r.Lang] = strings.Split(r.Text, "|")
			}
		}
		if len(b.Texts) > 0 {
			out = append(out, b)
		}
	}
	return out
}

// echoModel 是照抄原文的假翻譯
type echoModel struct{}

func (echoModel) Respond(ctx context.Context, r model.Request) (*model.Reply, model.Usage) {
	return &model.Reply{Values: map[string]model.Value{"translation": model.Text(r.Text)}}, model.Usage{}
}

// apiClient 直接呼叫後台的 handler（不開連接埠），帶著登入後的 session cookie
type apiClient struct {
	ctx    context.Context
	h      http.Handler
	cookie *http.Cookie
}

func (c *apiClient) login(username, password string) error {
	w, err := c.send("POST", "/api/auth/login", map[string]string{"username": username, "password": password})
	if err != nil {
		return err
	}
	cookies := w.Result().Cookies()
	i := slices.IndexFunc(cookies, func(x *http.Cookie) bool { return x.Name == SessionCookie })
	if i < 0 {
		return fmt.Errorf("登入後沒有 session cookie")
	}
	c.cookie = cookies[i]
	return nil
}

// do 送一個請求；body 不是 nil 時用 JSON 送出，out 不是 nil 時讀回應
func (c *apiClient) do(method, path string, body, out any) error {
	w, err := c.send(method, path, body)
	if err != nil || out == nil {
		return err
	}
	return json.Unmarshal(w.Body.Bytes(), out)
}

// send 送一個請求；不是 2xx 時回錯誤（含後台回的訊息）
func (c *apiClient) send(method, path string, body any) (*httptest.ResponseRecorder, error) {
	var buf bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return nil, err
		}
	}
	req := httptest.NewRequestWithContext(c.ctx, method, path, &buf)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.cookie != nil {
		req.AddCookie(c.cookie)
	}
	w := httptest.NewRecorder()
	c.h.ServeHTTP(w, req)
	if w.Code/100 != 2 {
		return nil, fmt.Errorf("%s %s：%d %s", method, path, w.Code, strings.TrimSpace(w.Body.String()))
	}
	return w, nil
}
