package telegram

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/store"
)

// Bot 處理一則 update：私訊文字交給對話層，先送「正在輸入」再回對話層的回覆，進出各寫一筆對話紀錄；其他一律忽略。
// 轉人工中對話層不回覆，也不顯示正在輸入
type Bot struct {
	API    *Client
	Store  *store.DB
	Log    *slog.Logger
	Now    func() time.Time // nil 用 time.Now
	Dialog *dialog.Bot
}

func (b *Bot) now() time.Time {
	if b.Now != nil {
		return b.Now()
	}
	return time.Now()
}

func (b *Bot) Process(ctx context.Context, u Update) {
	m := u.Message
	switch {
	case u.CallbackQuery != nil:
		b.handleCallback(ctx, u.UpdateID, u.CallbackQuery)
		return
	case m == nil:
		example.invalid("忽略不支援的 update", "update", u.UpdateID)
		return
	case m.Chat.Type == "group" || m.Chat.Type == "supergroup":
		b.handleGroup(ctx, u.UpdateID, m)
		return
	case m.Chat.Type != "private":
		example.invalid("忽略非私訊", "update", u.UpdateID, "chat_type", m.Chat.Type)
		return
	case m.From == nil:
		example.invalid("私訊沒有發送者，略過", "update", u.UpdateID)
		return
	}
	entry := store.LogEntry{
		Platform: "telegram",
		UserID:   strconv.FormatInt(m.From.ID, 10),
		ChatID:   strconv.FormatInt(m.Chat.ID, 10),
	}
	if kind := m.mediaKind(); kind != "" {
		b.handleMedia(ctx, u.UpdateID, m, entry, kind)
		return
	}
	text := m.Text
	// 貼圖當閒聊：交給對話層的是貼圖的表情（只有符號，不呼叫 AI）
	if text == "" && m.Sticker != nil {
		text = m.Sticker.Emoji
		if text == "" {
			text = "🙂"
		}
	}
	if text == "" {
		example.invalid("私訊不是文字，略過", "update", u.UpdateID)
		return
	}
	b.appendLog(ctx, entry, "customer", text, "")

	if c, err := b.Dialog.LoadCustomer(ctx, entry.Platform, entry.UserID); err == nil && !c.Human {
		_ = b.API.SendChatAction(ctx, m.Chat.ID) // 失敗只影響顯示，Client 已經記錄
	}
	receive := b.Dialog.Receive
	if dialog.IsBotCommand(text) {
		receive = b.Dialog.Greet // /start、/help 不送進分類，回打招呼（#54）
	}
	r, err := receive(ctx, dialog.Message{Platform: entry.Platform, UserID: entry.UserID, ChatID: entry.ChatID, Text: text})
	if err != nil {
		b.Log.Error("對話層處理失敗", "update", u.UpdateID, "error", err.Error())
		return
	}
	b.reply(ctx, u.UpdateID, m, entry, r)
	for _, id := range r.Tickets {
		b.notifyTicket(ctx, id, m.From)
	}
}

// handleMedia：圖片、語音、檔案看不懂，回「目前只看得懂文字」並開單；原始訊息轉送（copyMessage）到員工群組，回覆在通知下面，
// 員工回覆轉送的那則也能回給顧客。轉送失敗時在通知後面註明
func (b *Bot) handleMedia(ctx context.Context, update int64, m *Message, entry store.LogEntry, kind string) {
	b.appendLog(ctx, entry, "customer", strings.TrimSpace("["+kind+"] "+m.Caption), "")
	r, err := b.Dialog.ReceiveMedia(ctx, dialog.Message{Platform: entry.Platform, UserID: entry.UserID, ChatID: entry.ChatID, Text: m.Caption}, kind)
	if err != nil {
		b.Log.Error("對話層處理失敗", "update", update, "error", err.Error())
		return
	}
	b.reply(ctx, update, m, entry, r)
	for _, id := range r.Tickets {
		if n := b.notifyTicket(ctx, id, m.From); n.Chat != 0 {
			b.forward(ctx, id, m, n)
		}
	}
}

// forward 把顧客的原始訊息複製到員工群組、回覆在通知下面（通知送不出去時照樣轉送，員工回覆它也能回給顧客）
func (b *Bot) forward(ctx context.Context, ticketID int64, m *Message, n staffNotice) {
	copied, err := b.API.CopyMessage(ctx, n.Chat, m.Chat.ID, m.MessageID, n.ID)
	if err != nil {
		if n.ID != 0 {
			_ = b.API.EditMessageText(ctx, n.Chat, n.ID, n.Text+"\n⚠️ 原始訊息轉送失敗，請直接私訊顧客確認")
		}
		return
	}
	if err := b.Store.SetTicketForwardedMessage(ctx, ticketID, copied); err != nil {
		b.Log.Error("記下轉送的訊息編號失敗", "ticket", ticketID, "error", err.Error())
	}
}

// reply 把對話層的回覆送給顧客並記一筆對話紀錄（附 AI 判斷）；轉人工中不回覆，紀錄照寫一筆空白的。
// 送不出去時紀錄照寫，標成送出失敗
func (b *Bot) reply(ctx context.Context, update int64, m *Message, entry store.LogEntry, r dialog.Reply) {
	if len(r.Traces) > 0 {
		if data, err := json.Marshal(r.Traces); err == nil {
			entry.Traces = data
		}
	}
	if r.Silent {
		b.appendLog(ctx, entry, "bot", "", r.Key)
		example.invalid("轉人工中，不回覆", "update", update)
		return
	}
	err := b.API.SendMessage(ctx, m.Chat.ID, r.Text)
	entry.SendFailed = err != nil
	b.appendLog(ctx, entry, "bot", r.Text, r.Key)
	if err != nil {
		b.Log.Error("回覆送出失敗", "update", update, "error", err.Error())
		return
	}
	example.invalid("已回覆私訊", "update", update, "reply", r.Key)
}

func (b *Bot) appendLog(ctx context.Context, e store.LogEntry, direction, text, replyKey string) {
	e.Time, e.Direction, e.Text, e.ReplyKey = b.now(), direction, text, replyKey
	if err := b.Store.AppendLog(ctx, e); err != nil {
		b.Log.Error("寫入對話紀錄失敗", "error", err.Error())
	}
}
