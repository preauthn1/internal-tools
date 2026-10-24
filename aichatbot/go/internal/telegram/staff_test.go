package telegram

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"aichatbot/internal/model"
)

// 員工群組（照搬 Swift 版 telegram-tests 的店員群組測試）：-100 是員工群組，顧客 555 私訊

const staffChat = -100

// send 送一個 update（update_id 自動編號）並等背景處理完
func (e *env) send(u map[string]any) {
	e.t.Helper()
	e.nextUpdate++
	u["update_id"] = e.nextUpdate
	body, _ := json.Marshal(u)
	e.post(testSecret, string(body))
	e.hook.Wait()
}

// groupText：群組 chat 裡使用者 from 傳的文字；replyTo 不是 0 時是回覆那則訊息
func groupText(chat, from int64, text string, replyTo int64) map[string]any {
	m := map[string]any{"message_id": 9, "chat": map[string]any{"id": chat, "type": "supergroup"},
		"from": map[string]any{"id": from, "first_name": "小明"}, "text": text}
	if replyTo != 0 {
		m["reply_to_message"] = map[string]any{"message_id": replyTo, "chat": map[string]any{"id": chat, "type": "supergroup"}}
	}
	return map[string]any{"message": m}
}

// privateMsg：顧客 user 私訊的文字
func privateMsg(user int64, text string) map[string]any {
	return map[string]any{"message_id": 5, "chat": map[string]any{"id": user, "type": "private"},
		"from": map[string]any{"id": user, "first_name": "王", "last_name": "小康", "username": "kang"}, "text": text}
}

// setFail 換掉 fail（背景處理同時在讀）
func (f *fakeAPI) setFail(fail func(method string, n int) bool) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.fail = fail
}

// params 是某個 method 的呼叫參數（依順序）
func (f *fakeAPI) params(method string) []map[string]any {
	f.mu.Lock()
	defer f.mu.Unlock()
	var out []map[string]any
	for _, c := range f.calls {
		if c.Method == method {
			out = append(out, c.Params)
		}
	}
	return out
}

// bind 把 staffChat 綁定成員工群組（不經過 webhook）
func (e *env) bind() {
	e.t.Helper()
	ctx := context.Background()
	s, err := e.db.HandoffSetting(ctx)
	if err != nil {
		e.t.Fatal(err)
	}
	if ok, err := e.db.BindStaffChat(ctx, s.BindCode, staffChat); !ok || err != nil {
		e.t.Fatalf("綁定失敗：%v", err)
	}
}

// intents 讓對話層依序把訊息分成這些意圖
func (e *env) intents(codes ...string) {
	var script []map[string]model.Value
	for _, c := range codes {
		script = append(script, map[string]model.Value{"intent": model.Text(c)})
	}
	e.dialog.Model = model.NewScripted(map[string][]map[string]model.Value{"Intent": script})
}

// toChat 是送到某個 chat 的 sendMessage 參數
func (e *env) toChat(chat int64) []map[string]any {
	var out []map[string]any
	for _, p := range e.api.params("sendMessage") {
		if p["chat_id"] == float64(chat) {
			out = append(out, p)
		}
	}
	return out
}

// buttons 是訊息上的按鈕的 callback_data
func buttons(p map[string]any) []string {
	var out []string
	markup, _ := p["reply_markup"].(map[string]any)
	rows, _ := markup["inline_keyboard"].([]any)
	for _, row := range rows {
		for _, b := range row.([]any) {
			out = append(out, b.(map[string]any)["callback_data"].(string))
		}
	}
	return out
}

// 轉人工：員工群組收到通知，有原因、最近的對話（含觸發的那句和機器人的回覆）和「結束轉人工」按鈕
func TestHandoffNotifiesStaff(t *testing.T) {
	e := newEnv(t)
	e.bind()
	e.intents("general", "human_request")
	e.send(map[string]any{"message": privateMsg(555, "你好")})
	e.send(map[string]any{"message": privateMsg(555, "我要找真人")})
	notices := e.toChat(staffChat)
	if len(notices) != 1 {
		t.Fatalf("群組通知：%v", notices)
	}
	text := notices[0]["text"].(string)
	for _, want := range []string{"#1【轉人工】🙋 轉人工：顧客要求找真人", "👤 顧客：你好", "🤖 機器人：您好～有什麼可以幫您的嗎？", "👤 顧客：我要找真人", "🤖 機器人：我先幫您轉給店員"} {
		if !strings.Contains(text, want) {
			t.Errorf("通知少了 %q：\n%s", want, text)
		}
	}
	if b := buttons(notices[0]); strings.Join(b, ",") != "end:1" {
		t.Errorf("按鈕：%v", b)
	}
}

// msgID 是 fakeAPI 回給這次呼叫的訊息編號
func msgID(p map[string]any) int64 { return int64(p["_message_id"].(float64)) }

// 員工「回覆」通知：原文轉給顧客（加上「店員回覆」）、記一筆員工的對話紀錄，群組回「已轉給顧客」；
// 回覆的不是通知（員工聊天）不理會
func TestStaffReplyForwardedToCustomer(t *testing.T) {
	e := newEnv(t)
	e.bind()
	e.intents("human_request")
	e.send(map[string]any{"message": privateMsg(555, "我要找真人")})
	notice := msgID(e.toChat(staffChat)[0])

	e.send(groupText(staffChat, 7, "隨便聊聊", 12345))
	e.send(groupText(staffChat, 7, "您好，我是店長", notice))
	toCustomer := e.toChat(555)
	if last := toCustomer[len(toCustomer)-1]; last["text"] != "💬 店員回覆：您好，我是店長" {
		t.Errorf("顧客收到：%v", last)
	}
	group := e.toChat(staffChat)
	if len(group) != 2 || !strings.Contains(group[1]["text"].(string), "已轉給顧客") {
		t.Errorf("群組：%v", group)
	}
	logs, _ := e.db.Logs(context.Background(), "telegram", "555")
	if last := logs[len(logs)-1]; last.Direction != "staff" || last.Text != "💬 店員回覆：您好，我是店長" {
		t.Errorf("對話紀錄：%+v", last)
	}
}

// callback：員工按了 chat 裡訊息 notice 上的按鈕
func callback(id string, chat, notice int64, text, data string) map[string]any {
	return map[string]any{"callback_query": map[string]any{
		"id": id, "data": data, "from": map[string]any{"id": 7, "first_name": "小明", "username": "ming"},
		"message": map[string]any{"message_id": notice, "chat": map[string]any{"id": chat, "type": "supergroup"}, "text": text},
	}}
}

// 結束轉人工：通知改成已結束（拿掉按鈕），AI 恢復回覆；重複按只回「已處理」。
// 不是員工群組的按鈕只回應 Telegram，不做任何事
func TestEndHandoffButton(t *testing.T) {
	e := newEnv(t)
	e.bind()
	e.intents("human_request", "general")
	e.send(map[string]any{"message": privateMsg(555, "我要找真人")})
	n := e.toChat(staffChat)[0]
	notice, text := msgID(n), n["text"].(string)

	e.send(callback("q0", -999, notice, text, "end:1"))
	if a := e.api.params("answerCallbackQuery"); len(a) != 1 || a[0]["callback_query_id"] != "q0" || len(e.api.params("editMessageText")) != 0 {
		t.Fatalf("非員工群組：%v", e.api.methods())
	}
	e.send(callback("q1", staffChat, notice, text, "end:1"))
	edits := e.api.params("editMessageText")
	if len(edits) != 1 || edits[0]["message_id"] != float64(notice) || !strings.HasPrefix(edits[0]["text"].(string), text) ||
		!strings.Contains(edits[0]["text"].(string), "已結束轉人工（操作：小明 @ming）") {
		t.Errorf("通知：%v", edits)
	}
	e.send(map[string]any{"message": privateMsg(555, "你好")})
	if c := e.toChat(555); c[len(c)-1]["text"] != "您好～有什麼可以幫您的嗎？" {
		t.Errorf("AI 沒恢復：%v", c[len(c)-1])
	}
	e.send(callback("q2", staffChat, notice, text, "end:1"))
	a := e.api.params("answerCallbackQuery")
	if a[len(a)-1]["text"] != "已處理" || len(e.api.params("editMessageText")) != 1 {
		t.Errorf("重複按：%v", a)
	}
}

// book：顧客 user 一句話送出預約申請，回傳群組裡那則預約通知
func (e *env) book(user int64) map[string]any {
	e.t.Helper()
	e.dialog.Model = model.NewScripted(map[string][]map[string]model.Value{
		"Intent": {{"intent": model.Text("booking")}},
		"Booking": {{"date_quote": model.Text("10/21"), "date_month": model.Integer(10), "date_day": model.Integer(21),
			"time_quote": model.Text("下午2點"), "time_hour": model.Integer(14), "people": model.Integer(4), "people_quote": model.Text("4位"),
			"name": model.Text("小康"), "phone": model.Text("0912345678")}},
	})
	e.send(map[string]any{"message": privateMsg(user, "10/21 下午2點 4位 小康 0912345678")})
	n := e.toChat(staffChat)
	if len(n) == 0 {
		e.t.Fatal("沒有預約通知")
	}
	return n[len(n)-1]
}

// 預約：通知附確認、拒絕按鈕；確認後顧客收到確認、通知改成已確認（拿掉按鈕）；重複按回「已處理」。
// 拒絕時傳給顧客失敗要在通知註明
func TestBookingButtons(t *testing.T) {
	e := newEnv(t)
	e.bind()
	n := e.book(555)
	notice, text := msgID(n), n["text"].(string)
	if b := strings.Join(buttons(n), ","); b != "confirm:1,reject:1" || !strings.Contains(text, "#1【預約】📅 預約申請：") {
		t.Fatalf("預約通知：%v %s", b, text)
	}
	e.send(callback("q1", staffChat, notice, text, "confirm:1"))
	if c := e.toChat(555); !strings.HasPrefix(c[len(c)-1]["text"].(string), "您的預約已確認：") {
		t.Errorf("顧客收到：%v", c[len(c)-1])
	}
	edits := e.api.params("editMessageText")
	if len(edits) != 1 || edits[0]["text"] != text+"\n\n✅ 已確認（操作：小明 @ming）" {
		t.Errorf("通知：%v", edits)
	}
	e.send(callback("q2", staffChat, notice, text, "reject:1"))
	a := e.api.params("answerCallbackQuery")
	if a[len(a)-1]["text"] != "已處理" || len(e.api.params("editMessageText")) != 1 {
		t.Errorf("重複按：%v", a)
	}

	n = e.book(556)
	e.api.setFail(func(method string, _ int) bool { return method == "sendMessage" })
	e.send(callback("q3", staffChat, msgID(n), n["text"].(string), "reject:2"))
	e.api.setFail(nil)
	edits = e.api.params("editMessageText")
	if got := edits[len(edits)-1]["text"].(string); !strings.HasSuffix(got, "❌ 已拒絕（操作：小明 @ming）\n⚠️ 傳給顧客失敗") {
		t.Errorf("拒絕：%s", got)
	}
}

// photo：顧客 user 傳一張照片（附說明文字）
func photo(user int64, caption string) map[string]any {
	m := privateMsg(user, "")
	delete(m, "text")
	m["message_id"], m["photo"], m["caption"] = 77, []any{map[string]any{"file_id": "f1"}}, caption
	return map[string]any{"message": m}
}

// 圖片：顧客收到「目前只看得懂文字」；群組收到通知，原始照片轉送（copyMessage）到群組並回覆在通知下面；
// 員工回覆轉送的那張照片，也能回給顧客
func TestPhotoForwardedToStaff(t *testing.T) {
	e := newEnv(t)
	e.bind()
	e.send(photo(555, "這個有貨嗎"))
	if c := e.toChat(555); len(c) != 1 || c[0]["text"] != "抱歉，目前只看得懂文字喔，已幫您通知店員～" {
		t.Errorf("顧客收到：%v", c)
	}
	n := e.toChat(staffChat)
	if len(n) != 1 || !strings.Contains(n[0]["text"].(string), "#1【非文字訊息】📎 顧客傳了圖片：這個有貨嗎") {
		t.Fatalf("群組通知：%v", n)
	}
	copies := e.api.params("copyMessage")
	if len(copies) != 1 || copies[0]["chat_id"] != float64(staffChat) || copies[0]["from_chat_id"] != float64(555) || copies[0]["message_id"] != float64(77) ||
		copies[0]["reply_parameters"].(map[string]any)["message_id"] != float64(msgID(n[0])) {
		t.Fatalf("轉送：%v", copies)
	}
	e.send(groupText(staffChat, 7, "有的喔", msgID(copies[0])))
	if c := e.toChat(555); c[len(c)-1]["text"] != "💬 店員回覆：有的喔" {
		t.Errorf("顧客收到：%v", c[len(c)-1])
	}
}

// 原始訊息轉送失敗：通知後面註明，顧客照常收到回覆
func TestForwardFailureNoted(t *testing.T) {
	e := newEnv(t)
	e.bind()
	e.api.setFail(func(method string, _ int) bool { return method == "copyMessage" })
	e.send(photo(555, ""))
	n := e.toChat(staffChat)
	edits := e.api.params("editMessageText")
	if len(n) != 1 || len(edits) != 1 || edits[0]["message_id"] != float64(msgID(n[0])) ||
		edits[0]["text"] != n[0]["text"].(string)+"\n⚠️ 原始訊息轉送失敗，請直接私訊顧客確認" {
		t.Errorf("通知：%v\n修改：%v", n, edits)
	}
	if c := e.toChat(555); len(c) != 1 {
		t.Errorf("顧客收到：%v", c)
	}
}

// 貼圖當閒聊：交給對話層的是貼圖的表情（只有符號，不呼叫 AI），不開單
func TestStickerIsSmallTalk(t *testing.T) {
	e := newEnv(t)
	e.bind()
	m := privateMsg(555, "")
	delete(m, "text")
	m["sticker"] = map[string]any{"emoji": "👍"}
	e.send(map[string]any{"message": m})
	if c := e.toChat(555); len(c) != 1 || c[0]["text"] != "收到～" {
		t.Errorf("顧客收到：%v", c)
	}
	if n := e.toChat(staffChat); len(n) != 0 {
		t.Errorf("不該通知群組：%v", n)
	}
	logs, _ := e.db.Logs(context.Background(), "telegram", "555")
	if len(logs) != 2 || logs[0].Text != "👍" {
		t.Errorf("對話紀錄：%+v", logs)
	}
}

// 還沒綁定員工群組：照常轉人工、回覆顧客，只在日誌記下沒通知
func TestUnboundStillRepliesCustomer(t *testing.T) {
	e := newEnv(t)
	e.intents("human_request")
	e.send(map[string]any{"message": privateMsg(555, "我要找真人")})
	if m := strings.Join(e.api.methods(), ","); m != "sendChatAction,sendMessage" {
		t.Errorf("呼叫的 API：%s", m)
	}
	if !strings.Contains(e.log.String(), "沒有員工群組，未通知") {
		t.Errorf("日誌：%s", e.log.String())
	}
}

// 綁定的群組升級成超級群組（chat id 會變）：改綁新的 id；別的群組升級不理會
func TestStaffGroupMigrated(t *testing.T) {
	e := newEnv(t)
	e.bind()
	migrate := func(from, to int64) map[string]any {
		return map[string]any{"message": map[string]any{"message_id": 1, "chat": map[string]any{"id": from, "type": "group"},
			"from": map[string]any{"id": 1}, "migrate_to_chat_id": to}}
	}
	e.send(migrate(-5, -1005))
	if s, _ := e.db.HandoffSetting(context.Background()); s.StaffChatID != staffChat {
		t.Errorf("別的群組升級卻改了：%d", s.StaffChatID)
	}
	e.send(migrate(staffChat, -1000100))
	if s, _ := e.db.HandoffSetting(context.Background()); s.StaffChatID != -1000100 {
		t.Errorf("沒改綁新的 id：%d", s.StaffChatID)
	}
}

// 費用超過上限這類通知直接送到員工群組；還沒綁定時只記日誌
func TestNotifyStaffText(t *testing.T) {
	e := newEnv(t)
	e.bot.NotifyStaff(context.Background(), "⚠️ 費用到上限")
	if m := e.api.methods(); len(m) != 0 {
		t.Errorf("還沒綁定：%v", m)
	}
	e.bind()
	e.bot.NotifyStaff(context.Background(), "⚠️ 費用到上限")
	if n := e.toChat(staffChat); len(n) != 1 || n[0]["text"] != "⚠️ 費用到上限" {
		t.Errorf("群組：%v", n)
	}
}

// 員工在群組傳「/bind 綁定碼」就綁定那個群組，機器人回「已綁定」；綁定後換一組新碼。錯的碼不理會
func TestBindStaffGroup(t *testing.T) {
	e := newEnv(t)
	ctx := context.Background()
	s, _ := e.db.HandoffSetting(ctx)
	e.send(groupText(staffChat, 1, "/bind WRONGCODE", 0))
	if m := e.api.methods(); len(m) != 0 {
		t.Errorf("錯的碼不該回應：%v", m)
	}
	e.send(groupText(staffChat, 1, "/bind "+s.BindCode, 0))
	after, _ := e.db.HandoffSetting(ctx)
	if after.StaffChatID != staffChat || after.BindCode == s.BindCode {
		t.Errorf("綁定後：%+v", after)
	}
	sent := e.api.params("sendMessage")
	if len(sent) != 1 || sent[0]["chat_id"] != float64(staffChat) || !strings.Contains(sent[0]["text"].(string), "已綁定") {
		t.Errorf("回應：%v", sent)
	}
}
