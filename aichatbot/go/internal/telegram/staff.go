package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"aichatbot/internal/dialog"
)

// MARK: 員工群組（照搬 Swift 版 Server.swift 的 handleStaff、handleCallback、notifyStaff）。
// 群組用綁定碼綁定（存在資料庫的 handoff_setting），只認綁定的那個群組

// handleGroup：「/bind 綁定碼」（或只傳綁定碼）綁定這個群組；綁定的群組裡員工「回覆」通知，內容轉給顧客；其他忽略
func (b *Bot) handleGroup(ctx context.Context, update int64, m *Message) {
	if m.MigrateToChatID != 0 {
		if ok, err := b.Store.MigrateStaffChat(ctx, m.Chat.ID, m.MigrateToChatID); err != nil {
			b.Log.Error("員工群組改綁新的 chat id 失敗", "update", update, "error", err.Error())
		} else if ok {
			example.invalid("員工群組升級成超級群組，已改綁新的 chat id", "update", update, "chat", m.MigrateToChatID)
		}
		return
	}
	if code, ok := bindCode(m.Text); ok {
		bound, err := b.Store.BindStaffChat(ctx, code, m.Chat.ID)
		if err != nil {
			b.Log.Error("綁定員工群組失敗", "update", update, "error", err.Error())
			return
		}
		if bound {
			example.invalid("已綁定員工群組", "update", update, "chat", m.Chat.ID)
			_ = b.API.SendMessage(ctx, m.Chat.ID, "✅ 已綁定成員工群組：之後轉人工、預約和顧客傳的圖片都會通知到這裡，「回覆」通知就能和顧客對話")
			return
		}
	}
	hs, err := b.Store.HandoffSetting(ctx)
	if err != nil {
		b.Log.Error("讀轉人工設定失敗", "update", update, "error", err.Error())
		return
	}
	if hs.StaffChatID == 0 || m.Chat.ID != hs.StaffChatID {
		example.invalid("忽略不是員工群組的訊息", "update", update, "chat", m.Chat.ID)
		return
	}
	if m.ReplyToMessage == nil {
		example.invalid("員工群組訊息（不是回覆通知），略過", "update", update)
		return
	}
	notice := m.ReplyToMessage.MessageID
	t, ok, err := b.Store.TicketByGroupMessage(ctx, notice)
	if err != nil {
		b.Log.Error("找員工回覆的待處理事項失敗", "update", update, "error", err.Error())
		return
	}
	if !ok {
		example.invalid("員工回覆的不是待處理事項的通知，略過", "update", update)
		return
	}
	reply := func(text string) {
		_, _ = b.API.SendMessageID(ctx, m.Chat.ID, text, SendOptions{ReplyTo: m.MessageID})
	}
	if strings.TrimSpace(m.Text) == "" {
		reply("目前只能轉送文字")
		return
	}
	out, err := b.Dialog.StaffReply(ctx, t.ID, m.Text)
	if err != nil {
		b.Log.Error("員工回覆失敗", "update", update, "ticket", t.ID, "error", err.Error())
		reply("⚠️ 傳給顧客失敗")
		return
	}
	if !b.sendTo(ctx, out) {
		reply("⚠️ 傳給顧客失敗")
		return
	}
	example.invalid("員工回覆已轉給顧客", "update", update, "ticket", t.ID)
	reply("✅ 已轉給顧客")
}

// sendTo 把員工那邊的操作要給顧客的話送出去；送不出去回 false（Client 已經記錄）
func (b *Bot) sendTo(ctx context.Context, o dialog.Outgoing) bool {
	chat, err := strconv.ParseInt(o.ChatID, 10, 64)
	if err != nil {
		b.Log.Error("顧客的 chat id 不正確", "chat", o.ChatID)
		return false
	}
	return b.API.SendMessage(ctx, chat, o.Text) == nil
}

// bindCode 取出訊息裡的綁定碼：「/bind 碼」「/bind@機器人名稱 碼」或整則只有一個詞
func bindCode(text string) (string, bool) {
	f := strings.Fields(text)
	switch {
	case len(f) == 2 && (f[0] == "/bind" || strings.HasPrefix(f[0], "/bind@")):
		return f[1], true
	case len(f) == 1 && !strings.HasPrefix(f[0], "/"):
		return f[0], true
	}
	return "", false
}

var kindLabels = map[string]string{
	"handoff": "轉人工", "booking": "預約", "urgent": "緊急", "claim": "聲稱", "unsupported": "不支援的語言", "media": "非文字訊息", "question": "問題",
}

var speakers = map[string]string{"customer": "👤 顧客：", "bot": "🤖 機器人：", "staff": "🧑‍💼 員工："}

// clip：對話太長只留開頭（通知裡一則一行）
func clip(text string, max int) string {
	line := []rune(strings.ReplaceAll(text, "\n", " "))
	if len(line) > max {
		return string(line[:max]) + "…"
	}
	return string(line)
}

// NotifyStaff 送一則文字到員工群組（費用超過上限這類通知）；還沒綁定群組時只記日誌
func (b *Bot) NotifyStaff(ctx context.Context, text string) {
	hs, err := b.Store.HandoffSetting(ctx)
	if err != nil {
		b.Log.Error("讀轉人工設定失敗，未通知員工群組", "error", err.Error())
		return
	}
	if hs.StaffChatID == 0 {
		example.invalid("沒有員工群組，未通知", "text", text)
		return
	}
	_ = b.API.SendMessage(ctx, hs.StaffChatID, text)
}

// notifyTicket 把待處理事項通知員工群組：類型和通知文字、顧客名字（聊天中問到的名字，沒有就用 Telegram 名稱）＋@username＋使用者 ID。
// 轉人工附最近 5 則對話和「結束轉人工」按鈕；預約附「確認／拒絕」按鈕；其他附原話。
// 送出後記下通知的訊息編號（員工「回覆」通知時用它找回事項）。
// 回傳員工群組和通知：還沒綁定群組時 Chat 是 0；送不出去時 ID 是 0
func (b *Bot) notifyTicket(ctx context.Context, ticketID int64, from *User) staffNotice {
	hs, err := b.Store.HandoffSetting(ctx)
	if err != nil {
		b.Log.Error("讀轉人工設定失敗，未通知員工群組", "ticket", ticketID, "error", err.Error())
		return staffNotice{}
	}
	if hs.StaffChatID == 0 {
		example.invalid("沒有員工群組，未通知", "ticket", ticketID)
		return staffNotice{}
	}
	n := staffNotice{Chat: hs.StaffChatID}
	t, err := b.Store.Ticket(ctx, ticketID)
	if err != nil {
		b.Log.Error("讀待處理事項失敗，未通知員工群組", "ticket", ticketID, "error", err.Error())
		return n
	}
	name := strings.TrimSpace(from.FirstName + " " + from.LastName)
	if c, err := b.Dialog.LoadCustomer(ctx, t.Platform, t.UserID); err == nil && c.Name != "" {
		name = c.Name
	}
	if name == "" {
		name = "（沒有名字）"
	}
	who := "👤 " + name
	if from.Username != "" {
		who += " @" + from.Username
	}
	label := kindLabels[t.Kind]
	if label == "" {
		label = t.Kind
	}
	lines := []string{fmt.Sprintf("#%d【%s】%s", t.ID, label, t.Text), fmt.Sprintf("%s（ID %s）", who, t.UserID)}
	var opts SendOptions
	switch t.Kind {
	case "handoff":
		lines = append(lines, "🗂 最近對話：")
		recent, err := b.Store.RecentLogs(ctx, t.Platform, t.UserID, 5)
		if err != nil {
			b.Log.Error("讀最近對話失敗", "ticket", t.ID, "error", err.Error())
		}
		for _, e := range recent {
			lines = append(lines, speakers[e.Direction]+clip(e.Text, 100))
		}
		lines = append(lines, "↩️ 用「回覆」這則訊息和顧客對話；處理完按「結束轉人工」，AI 恢復回覆")
		opts.Buttons = [][]Button{{{Text: "結束轉人工", Data: fmt.Sprintf("end:%d", t.ID)}}}
	default:
		if t.Original != "" && !strings.Contains(t.Text, t.Original) {
			lines = append(lines, "💬 原話："+t.Original)
		}
		if t.Kind == "booking" {
			lines = append(lines, "👇 按按鈕確認或拒絕；用「回覆」這則訊息可以補充說明（例如拒絕的原因）給顧客")
			opts.Buttons = [][]Button{{{Text: "✅ 確認", Data: fmt.Sprintf("confirm:%d", t.ID)}, {Text: "❌ 拒絕", Data: fmt.Sprintf("reject:%d", t.ID)}}}
		} else {
			lines = append(lines, "↩️ 用「回覆」這則訊息，內容會轉給顧客")
		}
	}
	n.Text = strings.Join(lines, "\n")
	if n.ID, err = b.API.SendMessageID(ctx, hs.StaffChatID, n.Text, opts); err != nil {
		b.Log.Error("通知員工群組失敗", "ticket", t.ID, "error", err.Error())
		return n
	}
	if err := b.Store.SetTicketGroupMessage(ctx, t.ID, n.ID); err != nil {
		b.Log.Error("記下通知的訊息編號失敗", "ticket", t.ID, "error", err.Error())
	}
	example.invalid("已通知員工群組", "ticket", t.ID, "kind", t.Kind)
	return n
}

// staffNotice 是送到員工群組 Chat 的一則通知
type staffNotice struct {
	Chat int64
	ID   int64
	Text string
}

// handleCallback：員工群組通知上的按鈕。end:<事項編號> 結束轉人工；confirm:<編號>、reject:<編號> 確認或拒絕預約。
// 處理完把通知改成「原通知＋結果和操作的員工」（按鈕拿掉）；已處理過的只回提示。不是員工群組的按鈕只回應 Telegram，不做任何事
func (b *Bot) handleCallback(ctx context.Context, update int64, q *CallbackQuery) {
	hs, err := b.Store.HandoffSetting(ctx)
	if err != nil || hs.StaffChatID == 0 || q.Message == nil || q.Message.Chat.ID != hs.StaffChatID {
		example.invalid("不是員工群組的按鈕回呼，略過", "update", update)
		_ = b.API.AnswerCallbackQuery(ctx, q.ID, "")
		return
	}
	action, idText, _ := strings.Cut(q.Data, ":")
	id, err := strconv.ParseInt(idText, 10, 64)
	if err != nil {
		action = ""
	}
	staff := strings.TrimSpace(q.From.FirstName + " " + q.From.LastName)
	if staff == "" {
		staff = fmt.Sprintf("ID %d", q.From.ID)
	}
	if q.From.Username != "" {
		staff += " @" + q.From.Username
	}
	notice := q.Message
	edit := func(result string) {
		_ = b.API.EditMessageText(ctx, notice.Chat.ID, notice.MessageID, notice.Text+"\n\n"+result)
	}
	switch action {
	case "end":
		already, err := b.Dialog.EndHandoff(ctx, id)
		switch {
		case errors.Is(err, dialog.ErrTicketNotFound):
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "找不到這件事項")
		case err != nil:
			b.Log.Error("結束轉人工失敗", "update", update, "ticket", id, "error", err.Error())
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "處理失敗，請再按一次")
		case already:
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "已處理")
		default:
			example.invalid("結束轉人工", "update", update, "ticket", id)
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "已恢復 AI 回覆")
			edit("✅ 已結束轉人工（操作：" + staff + "），AI 恢復回覆")
		}
	case "confirm", "reject":
		decide, result, done := b.Dialog.ConfirmBooking, "✅ 已確認", "已確認"
		if action == "reject" {
			decide, result, done = b.Dialog.RejectBooking, "❌ 已拒絕", "已拒絕"
		}
		out, err := decide(ctx, id, staff)
		switch {
		case errors.Is(err, dialog.ErrAlreadyHandled):
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "已處理")
		case errors.Is(err, dialog.ErrTicketNotFound):
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "找不到這筆預約")
		case err != nil:
			b.Log.Error("處理預約失敗", "update", update, "booking", id, "error", err.Error())
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, "處理失敗，請再按一次")
		default:
			result += "（操作：" + staff + "）"
			if !b.sendTo(ctx, out) {
				result += "\n⚠️ 傳給顧客失敗"
			}
			example.invalid("員工"+done+"預約", "update", update, "booking", id)
			edit(result)
			_ = b.API.AnswerCallbackQuery(ctx, q.ID, done)
		}
	default:
		example.invalid("看不懂的按鈕回呼，略過", "update", update)
		_ = b.API.AnswerCallbackQuery(ctx, q.ID, "無法處理")
	}
}
