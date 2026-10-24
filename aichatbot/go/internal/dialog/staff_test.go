package dialog

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// setHandoff 改轉人工設定
func (e *env) setHandoff(change func(*store.HandoffSetting)) {
	e.t.Helper()
	s, err := e.db.HandoffSetting(context.Background())
	if err != nil {
		e.t.Fatal(err)
	}
	change(&s)
	if err := e.db.SaveHandoffSetting(context.Background(), s); err != nil {
		e.t.Fatal(err)
	}
}

// unbindStaffChat 解除員工群組的綁定（newEnv 預設綁好）
func (e *env) unbindStaffChat() {
	e.t.Helper()
	if _, err := e.db.SQL().Exec("UPDATE handoff_setting SET staff_chat_id = NULL WHERE id = 1"); err != nil {
		e.t.Fatal(err)
	}
}

// 沒綁員工群組時沒有人能接（#55）：不進轉人工，回 handoff_off（附聯絡方式），照樣開一張事項給後台看；下一句 AI 照常回
func TestHandoffWithoutStaffChat(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("general")}})
	e.unbindStaffChat()
	r := e.say("我要找真人")
	if r.Key != "handoff_off" || !strings.Contains(r.Text, "電話 02-1234-5678") || fmt.Sprint(r.Tickets) != "[1]" {
		t.Errorf("got %+v", r)
	}
	if c := e.customer(); c.Human || c.OpenTicket {
		t.Errorf("顧客狀態：%+v", c)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 想找真人（沒有員工群組，未轉人工）：顧客要求找真人" {
		t.Errorf("待處理：%q", p)
	}
	if r := e.say("你好"); r.Silent || r.Key != "general" {
		t.Errorf("AI 沒照常回：%+v", r)
	}
}

// 沒綁員工群組時，這位顧客已經有一張還沒處理的轉人工事項就不再開（每句都轉不了人工，不能每句開一張）
func TestHandoffWithoutStaffChatOpensOneTicket(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("human_request")}})
	e.unbindStaffChat()
	e.say("我要找真人")
	if r := e.say("真人呢"); r.Key != "handoff_off" || len(r.Tickets) != 0 {
		t.Errorf("got %+v", r)
	}
	if p := e.pending(); len(p) != 1 {
		t.Errorf("待處理：%q", p)
	}
}

// 聲稱／無關累計幾次轉人工照設定
func TestStreakLimitFromSetting(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("off_topic"), intent("claim")}})
	e.setHandoff(func(s *store.HandoffSetting) { s.StreakLimit = 2 })
	var keys []string
	for range 2 {
		keys = append(keys, e.say("一句話").Key)
	}
	if got := strings.Join(keys, ","); got != "off_topic,handoff" {
		t.Errorf("got %s", got)
	}
}

// 累計次數存成 0 以下（還沒有後台檢查）時用預設的 3 次，不會第一句就轉人工
func TestStreakLimitZeroUsesDefault(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("off_topic")}})
	e.setHandoff(func(s *store.HandoffSetting) { s.StreakLimit = 0 })
	if r := e.say("一句話"); r.Key != "off_topic" {
		t.Errorf("got %+v", r)
	}
}

// 員工服務時段：時間沒補零（9:00）照樣比對；結束早於開始是跨夜（22:00～02:00 涵蓋隔天凌晨）
func TestServiceHoursParsing(t *testing.T) {
	at := func(day, hour int) time.Time { return time.Date(2026, 10, day, hour, 30, 0, 0, taipei) } // 10/7 是週三
	for _, c := range []struct {
		h    store.ServiceHours
		now  time.Time
		want bool
	}{
		{store.ServiceHours{Weekday: 3, Start: "9:00", End: "18:00"}, at(7, 10), true},
		{store.ServiceHours{Weekday: 3, Start: "9:00", End: "18:00"}, at(7, 8), false},
		{store.ServiceHours{Weekday: 3, Start: "22:00", End: "02:00"}, at(7, 23), true},
		{store.ServiceHours{Weekday: 3, Start: "22:00", End: "02:00"}, at(8, 1), true},
		{store.ServiceHours{Weekday: 3, Start: "22:00", End: "02:00"}, at(8, 3), false},
		{store.ServiceHours{Weekday: 6, Start: "22:00", End: "02:00"}, at(11, 1), true}, // 週六跨到週日
	} {
		if got := inServiceHours([]store.ServiceHours{c.h}, c.now); got != c.want {
			t.Errorf("%+v %s：got %v", c.h, c.now.Format("Mon 15:04"), got)
		}
	}
}

// 提示文字有填就用它，沒填用固定回覆 handoff
func TestHandoffPromptText(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
	e.setHandoff(func(s *store.HandoffSetting) { s.PromptText = "專人馬上來" })
	if r := e.say("我要找真人"); r.Key != "handoff" || r.Text != "專人馬上來" {
		t.Errorf("got %+v", r)
	}
}

// 員工服務時段（台北時間，依星期幾）：時段外照樣轉人工，回覆後面加上時段外提示；沒設時段是全天
func TestHandoffServiceHours(t *testing.T) {
	// 2026-10-07 週三 12:00（台北）
	noon := time.Date(2026, 10, 7, 4, 0, 0, 0, time.UTC)
	for _, c := range []struct {
		hours    []store.ServiceHours
		offHours bool
	}{
		{nil, false},
		{[]store.ServiceHours{{Weekday: 3, Start: "09:00", End: "18:00"}}, false},
		{[]store.ServiceHours{{Weekday: 3, Start: "09:00", End: "12:00"}}, true},
		{[]store.ServiceHours{{Weekday: 4, Start: "00:00", End: "23:59"}}, true},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
		e.bot.Now = func() time.Time { return noon }
		e.setHandoff(func(s *store.HandoffSetting) { s.ServiceHours = c.hours })
		r := e.say("我要找真人")
		want := "我先幫您轉給店員，請稍候喔～"
		if c.offHours {
			want += "\n現在不是店員服務時間，店員上班後會盡快回覆您喔！"
		}
		if r.Key != "handoff" || r.Text != want {
			t.Errorf("%+v：got %+v", c.hours, r)
		}
	}
}

// 員工回覆待處理事項：回傳要送給顧客的話（加上「店員回覆」），事項結束、追問次數歸零，對話紀錄記一筆員工的
func TestStaffReply(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("claim"), intent("follow_up")}})
	e.say("上次說打五折")
	e.say("有消息了嗎")
	out, err := e.bot.StaffReply(ctx, 1, "有的，結帳時說一聲")
	if err != nil {
		t.Fatal(err)
	}
	if out.Platform != "telegram" || out.UserID != "U1" || out.ChatID != "C1" || out.Key != "staff_reply" || out.Text != "💬 店員回覆：有的，結帳時說一聲" {
		t.Errorf("送給顧客：%+v", out)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
	if c := e.customer(); c.OpenTicket || c.FollowUps != 0 {
		t.Errorf("顧客狀態：%+v", c)
	}
	logs, _ := e.db.Logs(ctx, "telegram", "U1")
	if len(logs) != 1 || logs[0].Direction != "staff" || logs[0].Text != out.Text || logs[0].ReplyKey != "staff_reply" || logs[0].ChatID != "C1" {
		t.Errorf("對話紀錄：%+v", logs)
	}
	if _, err := e.bot.StaffReply(ctx, 99, "x"); !errors.Is(err, ErrTicketNotFound) {
		t.Errorf("沒有這張：%v", err)
	}
}

// 員工回覆緊急通知：對應回原本那張，兩張都結束
func TestStaffReplyToUrgentClosesOriginal(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("claim"), intent("follow_up"), intent("follow_up"), intent("follow_up")}})
	for _, s := range []string{"上次說打五折", "有消息了嗎", "還要多久", "好了沒"} {
		e.say(s)
	}
	if p := e.pending(); len(p) != 2 {
		t.Fatalf("待處理：%q", p)
	}
	if _, err := e.bot.StaffReply(context.Background(), 2, "有的"); err != nil {
		t.Fatal(err)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
}

// 結束轉人工：AI 恢復回覆、轉人工事項結束；重複按回傳已經結束，不改狀態
func TestEndHandoff(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("general")}})
	e.say("我要找真人")
	already, err := e.bot.EndHandoff(ctx, 1)
	if err != nil || already {
		t.Fatalf("already=%v err=%v", already, err)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
	if r := e.say("你好"); r.Silent || r.Key != "general" {
		t.Errorf("AI 沒恢復回覆：%+v", r)
	}
	if already, err := e.bot.EndHandoff(ctx, 1); err != nil || !already {
		t.Errorf("重複按：already=%v err=%v", already, err)
	}
	if _, err := e.bot.EndHandoff(ctx, 99); !errors.Is(err, ErrTicketNotFound) {
		t.Errorf("沒有這張：%v", err)
	}
}

// 後台依顧客結束轉人工（#55）：效果同群組的按鈕，對話紀錄記一筆；不在轉人工中回傳已經結束
func TestEndCustomerHandoff(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("general")}})
	e.say("我要找真人")
	already, err := e.bot.EndCustomerHandoff(ctx, "telegram", "U1")
	if err != nil || already {
		t.Fatalf("already=%v err=%v", already, err)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
	logs, _ := e.db.Logs(ctx, "telegram", "U1")
	if l := logs[len(logs)-1]; l.Direction != "staff" || l.Text != "（結束轉人工）" || l.ReplyKey != "end_handoff" || l.ChatID != "C1" {
		t.Errorf("對話紀錄：%+v", l)
	}
	if r := e.say("你好"); r.Silent || r.Key != "general" {
		t.Errorf("AI 沒恢復回覆：%+v", r)
	}
	if already, err := e.bot.EndCustomerHandoff(ctx, "telegram", "U1"); err != nil || !already {
		t.Errorf("不在轉人工中：already=%v err=%v", already, err)
	}
}

// 這則訊息開了哪些待處理事項都回傳給平台層（之後送員工群組通知）
func TestReplyListsOpenedTickets(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general"), intent("claim"), intent("claim")}})
	e.setHandoff(func(s *store.HandoffSetting) { s.StreakLimit = 2 })
	var got []string
	for range 3 {
		got = append(got, fmt.Sprint(e.say("一句話").Tickets))
	}
	if strings.Join(got, ",") != "[],[1],[2 3]" {
		t.Errorf("開的事項：%v", got)
	}
}

// 圖片、語音、檔案：回「目前只看得懂文字」並開一張在等員工的事項（說明文字放進通知）；轉人工中不回覆，照樣開單
func TestMediaOpensTicket(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
	r, err := e.bot.ReceiveMedia(ctx, Message{Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "這個有貨嗎"}, "圖片")
	if err != nil {
		t.Fatal(err)
	}
	if r.Key != "text_only" || r.Text != "抱歉，目前只看得懂文字喔，已幫您通知店員～" || fmt.Sprint(r.Tickets) != "[1]" {
		t.Errorf("got %+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "media：📎 顧客傳了圖片：這個有貨嗎" {
		t.Errorf("待處理：%q", p)
	}
	if !e.customer().OpenTicket {
		t.Error("沒記下有事在等員工")
	}
	e.say("我要找真人")
	r, err = e.bot.ReceiveMedia(ctx, Message{Platform: "telegram", UserID: "U1", ChatID: "C1"}, "語音")
	if err != nil || !r.Silent || fmt.Sprint(r.Tickets) != "[3]" {
		t.Errorf("轉人工中：%+v %v", r, err)
	}
}

// addBooking 直接存一筆待確認的預約申請（顧客 telegram:U1），回傳編號
func (e *env) addBooking() int64 {
	e.t.Helper()
	ctx := context.Background()
	id, err := e.db.AddTicket(ctx, store.Ticket{Kind: "booking", Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "📅 預約申請：S"}, wednesday)
	if err != nil {
		e.t.Fatal(err)
	}
	if err := e.db.AddBooking(ctx, store.Booking{ID: id, Platform: "telegram", UserID: "U1", ChatID: "C1", Summary: "10/21 14:00, 4 位"}, wednesday); err != nil {
		e.t.Fatal(err)
	}
	return id
}

// 確認或拒絕預約：回傳要送給顧客的固定回覆（含預約內容），記下操作的員工；已處理過的回 ErrAlreadyHandled
func TestConfirmAndRejectBooking(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), nil)
	id := e.addBooking()
	out, err := e.bot.ConfirmBooking(ctx, id, "小明")
	if err != nil {
		t.Fatal(err)
	}
	if out.ChatID != "C1" || out.Key != "booking_confirmed" || out.Text != "您的預約已確認：10/21 14:00, 4 位。期待您的光臨！" {
		t.Errorf("確認：%+v", out)
	}
	if _, err := e.bot.RejectBooking(ctx, id, "小華"); !errors.Is(err, ErrAlreadyHandled) {
		t.Errorf("重複處理：%v", err)
	}
	if b, _ := e.db.Booking(ctx, id); b.Status != "confirmed" {
		t.Errorf("狀態：%s", b.Status)
	}

	id = e.addBooking()
	out, err = e.bot.RejectBooking(ctx, id, "小華")
	if err != nil || out.Key != "booking_rejected" || out.Text != "很抱歉，這次的預約無法安排：10/21 14:00, 4 位。" {
		t.Errorf("拒絕：%+v %v", out, err)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
	logs, _ := e.db.Logs(ctx, "telegram", "U1")
	if len(logs) != 2 || logs[0].Direction != "staff" || logs[1].ReplyKey != "booking_rejected" {
		t.Errorf("對話紀錄：%+v", logs)
	}
	if _, err := e.bot.ConfirmBooking(ctx, 99, "小明"); !errors.Is(err, ErrTicketNotFound) {
		t.Errorf("沒有這筆：%v", err)
	}
}

// #67：有自訂欄位、人數關掉時，員工群組的預約通知和確認、拒絕的回覆都照設定的欄位（有服務和設計師，沒有人數）
func TestStaffBookingMessagesFollowBookingFields(t *testing.T) {
	ctx := context.Background()
	script := func() map[string][]values {
		return map[string][]values{
			"Intent": {intent("booking")},
			"Booking": {{"date_quote": txt("10/9"), "date_month": num(10), "date_day": num(9), "time_quote": txt("下午三點"), "time_hour": num(15),
				"service": txt("染髮"), "designer": txt("小葵"), "name": txt("林小姐"), "phone": txt("0911222333")}},
		}
	}
	summary := "2026-10-09 15:00, 染髮, 小葵, 林小姐, 0911222333"
	e := newEnv(t, salonMerchant(), script())
	if r := e.say("10/9 下午三點染髮指定小葵，林小姐 0911222333"); r.Key != "booking_done" {
		t.Fatalf("預約：%+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "booking：📅 預約申請："+summary {
		t.Errorf("員工群組的通知：%q", p)
	}
	if out, err := e.bot.ConfirmBooking(ctx, 1, "小明"); err != nil || out.Text != "您的預約已確認："+summary+"。期待您的光臨！" {
		t.Errorf("確認：%+v %v", out, err)
	}

	e = newEnv(t, salonMerchant(), script())
	e.say("10/9 下午三點染髮指定小葵，林小姐 0911222333")
	if out, err := e.bot.RejectBooking(ctx, 1, "小華"); err != nil || out.Text != "很抱歉，這次的預約無法安排："+summary+"。" {
		t.Errorf("拒絕：%+v %v", out, err)
	}
}

// 舊通知上的「結束轉人工」不會結束之後另一次轉人工
func TestOldEndHandoffButtonIgnored(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("human_request")}})
	e.say("我要找真人")
	e.bot.EndHandoff(ctx, 1)
	e.say("再找一次真人")
	if already, err := e.bot.EndHandoff(ctx, 1); err != nil || !already {
		t.Errorf("舊按鈕：already=%v err=%v", already, err)
	}
	if !e.customer().Human {
		t.Error("舊按鈕結束了新的轉人工")
	}
}

// 轉人工中收到任何訊息（文字、想找真人、預約中的回答、圖片、bot 指令）都不回覆，也不呼叫 AI（#57）
func TestHandoffMakesNoAICalls(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general")}})
	data, _ := json.Marshal(Customer{Human: true, HumanSince: wednesday, AskingField: "date", Booking: map[string]string{"people": "4"}})
	e.db.SaveCustomer(ctx, "telegram", "U1", data, wednesday)
	m := Message{Platform: "telegram", UserID: "U1", ChatID: "C1"}
	for _, text := range []string{"你好", "我要找真人", "明天下午兩點", "停車場在哪"} {
		if r := e.say(text); !r.Silent {
			t.Errorf("%s：%+v", text, r)
		}
	}
	if r, err := e.bot.ReceiveMedia(ctx, m, "圖片"); err != nil || !r.Silent {
		t.Errorf("圖片：%+v %v", r, err)
	}
	m.Text = "/start"
	if r, err := e.bot.Greet(ctx, m); err != nil || !r.Silent {
		t.Errorf("bot 指令：%+v %v", r, err)
	}
	for _, req := range e.model.Requests() {
		t.Errorf("轉人工中呼叫了 AI：%s %q", req.Name, req.Text)
	}
}

// 這次改版前就在轉人工中的顧客（沒有轉人工時間）：開了自動結束也不會一傳訊息就恢復
func TestAutoEndWithoutHandoffTime(t *testing.T) {
	ctx := context.Background()
	e := newEnv(t, testMerchant(), nil)
	e.setHandoff(func(s *store.HandoffSetting) { s.AutoEndMinutes = 30 })
	data, _ := json.Marshal(Customer{Human: true})
	e.db.SaveCustomer(ctx, "telegram", "U1", data, wednesday)
	if r := e.say("在嗎"); !r.Silent {
		t.Errorf("沒有轉人工時間就恢復了：%+v", r)
	}
}

// 自動結束：員工最後一次回覆（沒回覆過就從轉人工起算）後 N 分鐘，顧客再傳訊息時恢復 AI 回覆；0 是不自動結束
func TestHandoffAutoEnd(t *testing.T) {
	ctx := context.Background()
	for _, minutes := range []int{30, 0} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("general")}})
		e.setHandoff(func(s *store.HandoffSetting) { s.AutoEndMinutes = minutes })
		now := wednesday
		e.bot.Now = func() time.Time { return now }
		e.say("我要找真人")
		now = wednesday.Add(10 * time.Minute)
		if _, err := e.bot.StaffReply(ctx, 1, "您好"); err != nil {
			t.Fatal(err)
		}
		now = wednesday.Add(35 * time.Minute)
		if r := e.say("還在嗎"); !r.Silent {
			t.Errorf("%d 分鐘：員工回覆後 25 分鐘就恢復了：%+v", minutes, r)
		}
		now = wednesday.Add(41 * time.Minute)
		r := e.say("你好")
		if minutes == 30 && (r.Silent || r.Key != "general") {
			t.Errorf("員工回覆後 31 分鐘還沒恢復：%+v", r)
		}
		if minutes == 0 && !r.Silent {
			t.Errorf("不自動結束卻恢復了：%+v", r)
		}
	}
	// 員工一直沒回覆：從轉人工起算
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("general")}})
	e.setHandoff(func(s *store.HandoffSetting) { s.AutoEndMinutes = 30 })
	now := wednesday
	e.bot.Now = func() time.Time { return now }
	e.say("我要找真人")
	now = wednesday.Add(31 * time.Minute)
	if r := e.say("你好"); r.Silent {
		t.Errorf("員工沒回覆，轉人工 31 分鐘後還沒恢復")
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("轉人工事項沒結束：%q", p)
	}
}
