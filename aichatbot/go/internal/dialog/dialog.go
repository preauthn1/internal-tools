// Package dialog 是對話層：一則顧客訊息進來，分類意圖、依「狀態 × 意圖」分派，回傳要回顧客的話。
// 照搬 Swift 版（tag swift-final 的 Sources/BotCore：Bot、Handle、Intent、Booking、RelativeDate、State），
// 改成 11 個內建意圖（CONTEXT.md），預約不檢查時段、公休日和人數（ADR 0003）。
package dialog

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// Message 是平台無關的顧客訊息
type Message struct {
	Platform string // telegram、terminal…
	UserID   string
	ChatID   string // 回覆送到哪個聊天室
	Text     string
}

// Reply 是要回顧客的話；Key 是回覆種類（固定回覆的 key）。Silent 是不回覆（轉人工中）。
// Tickets 是這則訊息開的待處理事項（依開的順序），平台層據此通知員工群組；Traces 是 AI 判斷，平台層寫進對話紀錄
type Reply struct {
	Key     string
	Text    string
	Silent  bool
	Tickets []int64
	Traces  []Trace
}

// Embedder 算句子向量（internal/embed 的 e5 模型）；找意圖案例庫裡和顧客訊息最像的句子
type Embedder interface {
	Embed(texts []string) ([][]float32, error)
}

// DefaultK 是分類時附幾句最像的意圖案例（Swift 版也是 6）
const DefaultK = 6

// Bot 處理顧客訊息；商家設定每次從資料庫讀（後台改了馬上生效）
type Bot struct {
	DB       *store.DB
	Model    model.ChatModel
	Embedder Embedder // nil：分類時不附相似案例
	K        int      // 附幾句相似案例；0 用 DefaultK
	Now      func() time.Time
	Log      *slog.Logger // nil 用 slog.Default()
	// Async 在背景執行 f（待補問題翻成商家慣用語言）；nil 時不在背景翻（後台詳細頁再補），
	// 終端機和量測不設，免得和下一題搶本機模型
	Async func(f func())

	mu        sync.Mutex
	vectors   map[string][]float32 // 案例句子 → 向量（案例庫會變，用句子當 key）
	customers keyedMutex           // 每位顧客一把鎖

	busyMu sync.Mutex
	busy   int           // 進行中的顧客回合數
	idle   chan struct{} // busy > 0 時才有，歸零時 close
}

// beginTurn 記一個進行中的顧客回合，回傳結束時呼叫的函式
func (b *Bot) beginTurn() (end func()) {
	b.busyMu.Lock()
	defer b.busyMu.Unlock()
	if b.busy++; b.busy == 1 {
		b.idle = make(chan struct{})
	}
	return func() {
		b.busyMu.Lock()
		defer b.busyMu.Unlock()
		if b.busy--; b.busy == 0 {
			close(b.idle)
			b.idle = nil
		}
	}
}

// waitCustomersIdle 等到沒有進行中的顧客回合；ctx 取消時回 ctx.Err()
func (b *Bot) waitCustomersIdle(ctx context.Context) error {
	for {
		b.busyMu.Lock()
		idle := b.idle
		b.busyMu.Unlock()
		if idle == nil {
			return ctx.Err()
		}
		select {
		case <-idle:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (b *Bot) log() *slog.Logger {
	if b.Log != nil {
		return b.Log
	}
	return slog.Default()
}

func (b *Bot) now() time.Time {
	if b.Now != nil {
		return b.Now()
	}
	return time.Now()
}

// Receive 處理一則訊息：讀出顧客狀態、處理、整份存回。
// 同一位顧客的訊息要由呼叫的人依序送進來（順序由呼叫的人決定）；和員工那邊的操作之間由每位顧客一把鎖互斥
func (b *Bot) Receive(ctx context.Context, m Message) (Reply, error) {
	return b.receive(ctx, m, func(t *turn) Reply { return t.handle(m.Text) })
}

// Greet 處理平台層的 bot 指令（/start、/help，見 IsBotCommand）：不分類、不拿指令判斷顧客語言，
// 回打招呼的固定回覆（新顧客用商家慣用語言，老顧客用已知的語言）；轉人工中不回覆
func (b *Bot) Greet(ctx context.Context, m Message) (Reply, error) {
	return b.receive(ctx, m, func(t *turn) Reply {
		if t.customer.Human {
			return Reply{Key: "silent", Silent: true}
		}
		return t.reply("greeting")
	})
}

// ReceiveMedia 處理顧客傳的圖片、語音、檔案（kind 是給員工看的種類，m.Text 是說明文字）：
// 看不懂，回「目前只看得懂文字」並開一張在等員工的事項（平台層把原始訊息轉送到員工群組）；轉人工中不回覆，照樣開單
func (b *Bot) ReceiveMedia(ctx context.Context, m Message, kind string) (Reply, error) {
	return b.receive(ctx, m, func(t *turn) Reply {
		notice := "📎 顧客傳了" + kind
		if m.Text != "" {
			notice += "：" + m.Text
		}
		t.msg.Text = strings.TrimSpace("[" + kind + "] " + m.Text)
		t.openTicket("media", notice, true, 0)
		if t.customer.Human {
			return Reply{Key: "silent", Silent: true}
		}
		return t.reply("text_only")
	})
}

func (b *Bot) receive(ctx context.Context, m Message, handle func(*turn) Reply) (Reply, error) {
	// 分類自我測試看到有顧客回合在處理就先暫停，把本機模型讓給顧客（#58）；排隊等顧客鎖的也算
	defer b.beginTurn()()
	defer b.lockCustomer(m.Platform, m.UserID)()
	cfg, err := b.config(ctx)
	if err != nil {
		return Reply{}, err
	}
	c, err := b.LoadCustomer(ctx, m.Platform, m.UserID)
	if err != nil {
		return Reply{}, err
	}
	if c.Lang == "" {
		c.Lang = cfg.Merchant.PreferredLanguage
	}
	hs, err := b.DB.HandoffSetting(ctx)
	if err != nil {
		return Reply{}, err
	}
	t := &turn{bot: b, ctx: ctx, cfg: cfg, handoffCfg: hs, customer: c, msg: m, now: b.now()}
	if t.customer.Human && t.autoEndDue() {
		t.endHandoff()
	}
	r := handle(t)
	if t.err != nil {
		return Reply{}, t.err
	}
	r.Tickets, r.Traces = t.opened, t.traces
	if err := b.saveCustomer(ctx, m.Platform, m.UserID, t.customer); err != nil {
		return Reply{}, err
	}
	t.recordPending()
	return r, nil
}

// Classify 只做分類（分類題組用）：和 Receive 同一個分類，顧客狀態由呼叫的人給，不存任何東西
func (b *Bot) Classify(ctx context.Context, text string, c Customer) (string, error) {
	cfg, err := b.config(ctx)
	if err != nil {
		return "", err
	}
	t := &turn{bot: b, ctx: ctx, cfg: cfg, customer: c, now: b.now()}
	intent, _ := t.classify(text)
	return intent, nil
}

// config 讀商家設定；還沒有設定時用預設值（商家資訊空白）
func (b *Bot) config(ctx context.Context) (merchant.Config, error) {
	cfg, err := merchant.Load(ctx, b.DB.SQL())
	if errors.Is(err, merchant.ErrNotConfigured) {
		return merchant.Default(), nil
	}
	return cfg, err
}

// LoadCustomer 讀出顧客狀態；沒有這位顧客時是空的
func (b *Bot) LoadCustomer(ctx context.Context, platform, userID string) (Customer, error) {
	var c Customer
	data, err := b.DB.LoadCustomer(ctx, platform, userID)
	if err != nil || data == nil {
		return c, err
	}
	return c, json.Unmarshal(data, &c)
}

// ResetCustomer 把顧客狀態清成新顧客（終端機的 /reset）
func (b *Bot) ResetCustomer(ctx context.Context, platform, userID string) error {
	return b.saveCustomer(ctx, platform, userID, Customer{})
}

func (b *Bot) saveCustomer(ctx context.Context, platform, userID string, c Customer) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}
	return b.DB.SaveCustomer(ctx, platform, userID, data, b.now())
}

// turn 是處理一則訊息時的脈絡：這位顧客的狀態（處理完整份存回）和商家設定
type turn struct {
	bot        *Bot
	ctx        context.Context
	cfg        merchant.Config
	handoffCfg store.HandoffSetting
	customer   Customer
	msg        Message
	now        time.Time
	err        error   // 寫資料庫失敗（第一個）；有錯就不存顧客狀態
	opened     []int64 // 這則訊息開的待處理事項
	traces     []Trace // 這則訊息的 AI 判斷
	filled     int     // 這則訊息進來時已收到幾個預約欄位（askField 判斷這一輪有沒有進展）

	missed       string // 這則訊息答不出來：待補問題的來源（unanswerable、verify_failed、unclear）；答得出來是空字串
	verifyFailed bool   // 這則訊息有原文片段驗證（或譯文驗證）沒過
	streakBefore int    // 這則訊息進來時的連續無關次數（意圖是一般對話時會先歸零，確認判成無關時還原，#48）

	vectors map[string][]float32 // 這則訊息裡算過的語意向量（見 vector）
}

// openTicket 開待處理事項；waiting 時同時記下「有事在等員工」（追問的次數從這裡開始算）。回傳編號（失敗是 0）
func (t *turn) openTicket(kind, text string, waiting bool, relatedTo int64) int64 {
	id, err := t.bot.DB.AddTicket(t.ctx, store.Ticket{
		Kind: kind, Platform: t.msg.Platform, UserID: t.msg.UserID, ChatID: t.msg.ChatID,
		Original: t.msg.Text, Text: text, RelatedTo: relatedTo,
	}, t.now)
	if err != nil && t.err == nil {
		t.err = err
	}
	if id != 0 {
		t.opened = append(t.opened, id)
	}
	if waiting && !t.customer.OpenTicket {
		t.customer.OpenTicket, t.customer.FollowUps = true, 0
	}
	return id
}
