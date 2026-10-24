package telegram

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Dedup 記住處理過的 update_id（Telegram 沒收到 200 會重送）：第一次看到回傳 true 並記下來
type Dedup interface {
	MarkUpdate(ctx context.Context, updateID int64, now time.Time) (bool, error)
}

// Webhook 收 Telegram 的 webhook：secret 不符回 401；同一個 update 只處理一次；先回 200，再在背景處理。
// 同一位顧客的 update 依收到的順序一則一則處理（狀態機讀出狀態、等 AI、存回去之間不能插隊），不同顧客同時處理
type Webhook struct {
	Secret  string
	Dedup   Dedup
	Process func(context.Context, Update)
	Log     *slog.Logger
	// Context 是背景處理用的 context（請求結束後還要繼續跑，不能用請求的）；nil 用 context.Background()
	Context context.Context

	wg    sync.WaitGroup
	queue keyedQueue
}

func (h *Webhook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "", http.StatusMethodNotAllowed)
		return
	}
	header := r.Header.Get("X-Telegram-Bot-Api-Secret-Token")
	// 比較時間不隨內容變化，避免用回應時間猜 secret
	if header == "" || subtle.ConstantTimeCompare([]byte(header), []byte(h.Secret)) != 1 {
		h.Log.Warn("webhook secret 不符，回 401")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(http.MaxBytesReader(w, r.Body, 1<<20))
	var u Update
	// 讀不懂的也回 200，不然 Telegram 會一直重送
	if err != nil || json.Unmarshal(body, &u) != nil || u.UpdateID == 0 {
		h.Log.Warn("webhook 內容無法解析，略過")
		w.WriteHeader(http.StatusOK)
		return
	}
	first, err := h.Dedup.MarkUpdate(r.Context(), u.UpdateID, time.Now())
	if err != nil {
		// 寫不進資料庫時當成第一次（寧可多處理一次，不要漏掉訊息）
		h.Log.Error("記錄 update 失敗，照常處理", "update", u.UpdateID, "error", err.Error())
		first = true
	}
	if !first {
		example.invalid("重複的 update，略過", "update", u.UpdateID)
		w.WriteHeader(http.StatusOK)
		return
	}
	ctx := h.Context
	if ctx == nil {
		ctx = context.Background()
	}
	h.wg.Add(1)
	h.queue.run(customerKey(u), func() {
		defer h.wg.Done()
		h.Process(ctx, u)
	})
	w.WriteHeader(http.StatusOK)
}

// customerKey：同一位顧客（私訊的發送者）的 update 依序處理；群組訊息和按鈕回呼（員工那邊）全部排成一列依序處理，
// 和顧客之間由對話層的每位顧客一把鎖互斥；沒有發送者的各自處理
func customerKey(u Update) string {
	switch m := u.Message; {
	case u.CallbackQuery != nil, m != nil && m.Chat.Type != "private":
		return "staff"
	case m != nil && m.From != nil:
		return "user:" + strconv.FormatInt(m.From.ID, 10)
	}
	return "update:" + strconv.FormatInt(u.UpdateID, 10)
}

// keyedQueue 讓同一個 key 的工作照登記的順序一個接一個跑，不同 key 同時跑（照搬 Swift 版 #11 的「依顧客加鎖」）。
// 登記在收到 webhook 的當下同步完成，所以順序就是收到的順序；每個 key 記著最後一件工作的完成訊號，跑完的 key 會清掉
type keyedQueue struct {
	mu    sync.Mutex
	tails map[string]chan struct{}
}

func (q *keyedQueue) run(key string, fn func()) {
	q.mu.Lock()
	if q.tails == nil {
		q.tails = map[string]chan struct{}{}
	}
	prev, done := q.tails[key], make(chan struct{})
	q.tails[key] = done
	q.mu.Unlock()
	go func() {
		if prev != nil {
			<-prev
		}
		fn()
		q.mu.Lock()
		if q.tails[key] == done {
			delete(q.tails, key)
		}
		q.mu.Unlock()
		close(done)
	}()
}

// Wait 等背景處理全部結束（關機和測試用）
func (h *Webhook) Wait() { h.wg.Wait() }
