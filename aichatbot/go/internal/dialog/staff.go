package dialog

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"time"

	"aichatbot/internal/store"
)

// MARK: 員工群組那邊的操作和轉人工設定（照搬 Swift 版 Bot.swift 的 staffReply、confirmBooking、rejectBooking、endHandoff）

// ErrTicketNotFound：沒有這張待處理事項（或這筆預約）
var ErrTicketNotFound = errors.New("找不到這件事項")

// Outgoing 是員工那邊的操作要送給顧客的一則訊息
type Outgoing struct {
	Platform string
	UserID   string
	ChatID   string
	Key      string // 固定回覆的 key
	Text     string
}

// StaffReply：員工回覆了待處理事項 ticketID。回傳要送給顧客的話（員工的原文加上「店員回覆」，不經過 AI），
// 這件事（含緊急通知）改成已處理、追問次數歸零，對話紀錄記一筆員工的。轉人工不會因此結束
func (b *Bot) StaffReply(ctx context.Context, ticketID int64, text string) (Outgoing, error) {
	t, err := b.ticket(ctx, ticketID)
	if err != nil {
		return Outgoing{}, err
	}
	defer b.lockCustomer(t.Platform, t.UserID)()
	if err := b.DB.ReplyTicket(ctx, t.ID); err != nil {
		return Outgoing{}, err
	}
	var out Outgoing
	err = b.updateCustomer(ctx, t.Platform, t.UserID, func(tt *turn) {
		tt.customer.OpenTicket, tt.customer.FollowUps, tt.customer.LastStaffReply = false, 0, tt.now
		out = Outgoing{Platform: t.Platform, UserID: t.UserID, ChatID: t.ChatID, Key: "staff_reply",
			Text: tt.text("staff_reply", map[string]string{"text": text})}
	})
	if err != nil {
		return Outgoing{}, err
	}
	b.logStaff(ctx, out)
	// 轉人工後員工的回覆預填成待補問題的答案（原話，不含「店員回覆」）；寫不進去只記錯誤
	if err := b.DB.AddPendingStaffAnswer(ctx, t, text); err != nil {
		b.log().Error("寫入待補問題的員工回覆失敗", "error", err.Error())
	}
	return out, nil
}

// ErrAlreadyHandled：預約已經確認或拒絕過了
var ErrAlreadyHandled = store.ErrAlreadyHandled

// ConfirmBooking 確認預約申請 id（by 是操作的員工），回傳要送給顧客的固定回覆
func (b *Bot) ConfirmBooking(ctx context.Context, id int64, by string) (Outgoing, error) {
	return b.decideBooking(ctx, id, "confirmed", by)
}

// RejectBooking 拒絕預約申請 id。原因由員工回覆那則預約通知另外轉給顧客
func (b *Bot) RejectBooking(ctx context.Context, id int64, by string) (Outgoing, error) {
	return b.decideBooking(ctx, id, "rejected", by)
}

// decideBooking：只有待確認的預約能處理，已確認或已拒絕的回 ErrAlreadyHandled；回覆用顧客自己的語言
func (b *Bot) decideBooking(ctx context.Context, id int64, status, by string) (Outgoing, error) {
	bk, err := b.DB.Booking(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return Outgoing{}, ErrTicketNotFound
	}
	if err != nil {
		return Outgoing{}, err
	}
	defer b.lockCustomer(bk.Platform, bk.UserID)()
	if err := b.DB.SetBookingStatus(ctx, id, status, by, b.now()); err != nil {
		return Outgoing{}, err
	}
	key := "booking_" + status
	out := Outgoing{Platform: bk.Platform, UserID: bk.UserID, ChatID: bk.ChatID, Key: key}
	err = b.updateCustomer(ctx, bk.Platform, bk.UserID, func(t *turn) {
		out.Text = t.text(key, map[string]string{"summary": bk.Summary})
	})
	if err != nil {
		return Outgoing{}, err
	}
	b.logStaff(ctx, out)
	return out, nil
}

// EndHandoff：員工按了轉人工通知上的「結束轉人工」，之後恢復 AI 回覆（不另外通知顧客）。
// 顧客已經不在轉人工中（重複按、按到舊通知）時回傳 already = true，不改任何狀態
func (b *Bot) EndHandoff(ctx context.Context, ticketID int64) (already bool, err error) {
	t, err := b.ticket(ctx, ticketID)
	if err != nil {
		return false, err
	}
	// 按到舊通知（之後又轉人工了一次）不結束新的那次；這次改版前就在轉人工中的顧客沒有記編號
	return b.endHandoffOf(ctx, t.Platform, t.UserID, func(c Customer) (string, bool) {
		return t.ChatID, c.HandoffTicket == 0 || c.HandoffTicket == t.ID
	})
}

// EndCustomerHandoff：後台在對話紀錄按了「結束轉人工」（#55），效果同 EndHandoff，以顧客定位。
// 顧客不在轉人工中時回傳 already = true
func (b *Bot) EndCustomerHandoff(ctx context.Context, platform, userID string) (already bool, err error) {
	return b.endHandoffOf(ctx, platform, userID, func(c Customer) (string, bool) {
		// 對話紀錄記在轉人工那張事項的聊天室；舊資料沒記事項編號時留空
		var chatID string
		if c.HandoffTicket != 0 {
			if t, err := b.DB.Ticket(ctx, c.HandoffTicket); err == nil {
				chatID = t.ChatID
			}
		}
		return chatID, true
	})
}

// endHandoffOf 鎖住顧客，在轉人工中而且 target 回傳 ok 時結束轉人工、對話紀錄記一筆（記在 target 回傳的聊天室）
func (b *Bot) endHandoffOf(ctx context.Context, platform, userID string, target func(Customer) (chatID string, ok bool)) (already bool, err error) {
	defer b.lockCustomer(platform, userID)()
	c, err := b.LoadCustomer(ctx, platform, userID)
	if err != nil {
		return false, err
	}
	if !c.Human {
		return true, nil
	}
	chatID, ok := target(c)
	if !ok {
		return true, nil
	}
	err = b.updateCustomer(ctx, platform, userID, func(tt *turn) { tt.endHandoff() })
	if err != nil {
		return false, err
	}
	b.logStaff(ctx, Outgoing{Platform: platform, UserID: userID, ChatID: chatID, Key: "end_handoff", Text: "（結束轉人工）"})
	return false, nil
}

// autoEndDue：有設自動結束，而且員工最後一次回覆（沒回覆過就從轉人工起算）已經超過設定的分鐘數
func (t *turn) autoEndDue() bool {
	n := t.handoffCfg.AutoEndMinutes
	if n <= 0 {
		return false
	}
	since := t.customer.HumanSince
	if t.customer.LastStaffReply.After(since) {
		since = t.customer.LastStaffReply
	}
	// 這次改版前就在轉人工中的顧客沒有時間，不自動結束
	if since.IsZero() {
		return false
	}
	return !t.now.Before(since.Add(time.Duration(n) * time.Minute))
}

// endHandoff 結束轉人工：狀態歸零，這位顧客還沒處理的轉人工事項改成已處理
func (t *turn) endHandoff() {
	c := &t.customer
	c.Human, c.OffTopicStreak, c.FollowUps, c.OpenTicket = false, 0, 0, false
	if err := t.bot.DB.CloseHandoffTickets(t.ctx, t.msg.Platform, t.msg.UserID); err != nil && t.err == nil {
		t.err = err
	}
}

func (b *Bot) ticket(ctx context.Context, id int64) (store.Ticket, error) {
	t, err := b.DB.Ticket(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return t, ErrTicketNotFound
	}
	return t, err
}

// updateCustomer 讀出顧客狀態、改完存回；change 拿到的 turn 只能用來改狀態和查固定回覆
func (b *Bot) updateCustomer(ctx context.Context, platform, userID string, change func(*turn)) error {
	cfg, err := b.config(ctx)
	if err != nil {
		return err
	}
	c, err := b.LoadCustomer(ctx, platform, userID)
	if err != nil {
		return err
	}
	if c.Lang == "" {
		c.Lang = cfg.Merchant.PreferredLanguage
	}
	t := &turn{bot: b, ctx: ctx, cfg: cfg, customer: c, msg: Message{Platform: platform, UserID: userID}, now: b.now()}
	change(t)
	if t.err != nil {
		return t.err
	}
	return b.saveCustomer(ctx, platform, userID, t.customer)
}

// logStaff 記一筆員工那邊送給顧客的對話紀錄。狀態已經存好了，寫不進去只記錯誤（照樣把話送給顧客）
func (b *Bot) logStaff(ctx context.Context, o Outgoing) {
	err := b.DB.AppendLog(ctx, store.LogEntry{Time: b.now(), Platform: o.Platform, UserID: o.UserID, ChatID: o.ChatID,
		Direction: "staff", Text: o.Text, ReplyKey: o.Key})
	if err != nil {
		b.log().Error("寫入員工的對話紀錄失敗", "error", err.Error())
	}
}

// keyedMutex 是每位顧客一把鎖：顧客的訊息（Receive）和員工那邊的操作（回覆、按鈕）讀出狀態到存回之間不能插隊
type keyedMutex struct {
	mu    sync.Mutex
	locks map[string]*keyedLock
}

type keyedLock struct {
	mu    sync.Mutex
	users int
}

// lockCustomer 鎖住這位顧客，回傳解鎖的函式
func (b *Bot) lockCustomer(platform, userID string) func() {
	k := &b.customers
	key := platform + ":" + userID
	k.mu.Lock()
	if k.locks == nil {
		k.locks = map[string]*keyedLock{}
	}
	l := k.locks[key]
	if l == nil {
		l = &keyedLock{}
		k.locks[key] = l
	}
	l.users++
	k.mu.Unlock()
	l.mu.Lock()
	return func() {
		l.mu.Unlock()
		k.mu.Lock()
		if l.users--; l.users == 0 {
			delete(k.locks, key)
		}
		k.mu.Unlock()
	}
}

// 員工服務時段照台北時間算（和費用上限的月份一樣用固定時區，開發容器不用裝時區資料）
var taipei = time.FixedZone("Asia/Taipei", 8*60*60)

// inServiceHours：now 在員工服務時段內；沒設時段是全天。結束不晚於開始的時段跨夜（22:00～02:00 涵蓋隔天凌晨）；
// 時間看不懂的時段略過
func inServiceHours(hours []store.ServiceHours, now time.Time) bool {
	if len(hours) == 0 {
		return true
	}
	now = now.In(taipei)
	day, clock := int(now.Weekday()), now.Hour()*60+now.Minute()
	for _, h := range hours {
		start, ok1 := minutesOf(h.Start)
		end, ok2 := minutesOf(h.End)
		switch {
		case !ok1 || !ok2:
		case start < end:
			if h.Weekday == day && start <= clock && clock < end {
				return true
			}
		default:
			if h.Weekday == day && clock >= start || (h.Weekday+1)%7 == day && clock < end {
				return true
			}
		}
	}
	return false
}

// minutesOf 把 "9:00"、"09:00" 換成從 0 點起算的分鐘數
func minutesOf(s string) (int, bool) {
	t, err := time.Parse("15:04", s)
	if err != nil {
		return 0, false
	}
	return t.Hour()*60 + t.Minute(), true
}
