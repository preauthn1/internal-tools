package dialog

import (
	"math/rand/v2"
	"strings"

	"aichatbot/internal/merchant"
)

// 內建的固定回覆（照搬 Swift 版的 eval/replies.json，拿掉測試店家的店名和營業時間）。
// 可以用 | 分隔多句（隨機挑一句）；{name} 這類是變數，{merchant} 是商家名稱
var builtinReplies = map[string]map[string]string{
	"general":           {"zh": "收到～", "en": "Got it~", "ja": "承知しました～"},
	"goodbye":           {"zh": "謝謝您，期待再為您服務～", "en": "Thank you, hope to see you again soon~", "ja": "ありがとうございました、またのご利用をお待ちしております～"},
	"thanks":            {"zh": "不客氣～有需要隨時找我喔！", "en": "You're welcome~ Let me know if you need anything else!", "ja": "どういたしまして～また何かあればお気軽にどうぞ！"},
	"opener":            {"zh": "在喔～有什麼可以幫您？", "en": "I'm here~ How can I help you?", "ja": "はい、おります～どのようなご用件でしょうか？"},
	"unclear":           {"zh": "不太確定您的意思，可以再說詳細一點嗎？", "en": "Sorry, I'm not quite sure what you mean. Could you tell me a bit more?", "ja": "すみません、もう少し詳しく教えていただけますか？"},
	"complaint":         {"zh": "抱歉讓您不滿意 🙏 需要的話我可以幫您轉給店員喔！", "en": "Sorry for the trouble 🙏 I can connect you with our staff if you'd like!", "ja": "ご不便をおかけしてすみません 🙏 必要ならスタッフにおつなぎします！"},
	"claim":             {"zh": "這部分我沒辦法確認，已幫您記錄下來，會請店員回覆您喔！", "en": "I can't confirm that, but I've noted it down and our staff will get back to you!", "ja": "こちらでは確認できないため、記録してスタッフからご連絡しますね！"},
	"off_topic":         {"zh": "這個我不太清楚耶～", "en": "Sorry, I'm not sure about that~", "ja": "すみません、それはちょっと分かりません～"},
	"waiting":           {"zh": "已經通知店員了，店員看到後會盡快回覆您，請稍等喔～", "en": "I've already notified our staff — they'll get back to you as soon as they can. Thanks for waiting~", "ja": "スタッフにはすでに連絡済みです。確認でき次第ご返信しますので、少々お待ちください～"},
	"nudge":             {"zh": "我再幫您催一下店員，造成不便很抱歉 🙏", "en": "I'll nudge our staff again — sorry for the wait 🙏", "ja": "スタッフにもう一度催促しますね。お待たせして申し訳ありません 🙏"},
	"handoff":           {"zh": "我先幫您轉給店員，請稍候喔～", "en": "Let me connect you with our staff, please wait a moment~", "ja": "スタッフにおつなぎしますので、少々お待ちください～"},
	"unsupported":       {"zh": "抱歉，目前只能用{languages}為您服務，已幫您通知店員～", "en": "Sorry, we currently support {languages} only. Our staff will get back to you.", "ja": "申し訳ありません、現在は{languages}のみ対応しております。"},
	"staff_reply":       {"zh": "💬 店員回覆：{text}", "en": "💬 Reply from our staff: {text}", "ja": "💬 スタッフより：{text}"},
	"text_only":         {"zh": "抱歉，目前只看得懂文字喔，已幫您通知店員～", "en": "Sorry, I can only read text messages for now. I've let our staff know~", "ja": "申し訳ありません、現在はテキストしか読めません。スタッフにお伝えしました～"},
	"off_hours":         {"zh": "現在不是店員服務時間，店員上班後會盡快回覆您喔！", "en": "Our staff are off right now. They'll get back to you as soon as they're back!", "ja": "ただいまスタッフの対応時間外です。対応時間になりしだいご返信します！"},
	"booking_off":       {"zh": "目前沒有開放線上預約，可以透過這個方式聯絡我們：{contact}", "en": "Online booking isn't available right now. You can reach us here: {contact}", "ja": "現在オンライン予約は受け付けておりません。こちらからご連絡ください：{contact}"},
	"handoff_off":       {"zh": "目前沒有真人客服，可以透過這個方式聯絡我們：{contact}", "en": "We don't have live staff chat right now. You can reach us here: {contact}", "ja": "現在スタッフによる対応は行っておりません。こちらからご連絡ください：{contact}"},
	"ask_date":          {"zh": "想預約哪一天呢？", "en": "Which date would you like to book?", "ja": "ご希望の日にちはいつですか？"},
	"ask_time":          {"zh": "想約幾點呢？", "en": "What time would you like?", "ja": "何時がご希望ですか？"},
	"ask_time_midnight": {"zh": "晚上 12 點是隔天的凌晨 0 點，想約哪一天的幾點呢？可以連同日期說「幾月幾號凌晨 0 點」這樣～", "en": "Midnight is 0:00 of the next day. Which date and time would you like? Please tell me the date too, like \"[month/day] at 0:00\".", "ja": "夜12時は翌日の0時になります。何日の何時がご希望ですか？「○月○日 0時」のように日付も一緒にお知らせください。"},
	"ask_people":        {"zh": "請問幾位呢？", "en": "How many people?", "ja": "何名様ですか？"},
	"ask_name":          {"zh": "請問預約人大名？", "en": "May I have your name for the reservation?", "ja": "ご予約のお名前をお願いします。"},
	"ask_phone":         {"zh": "請留一支聯絡電話～", "en": "Could you leave a contact phone number?", "ja": "ご連絡先の電話番号をお願いします～"},
	"err_date":          {"zh": "日期看不太懂，可以再說一次嗎？", "en": "Sorry, I didn't catch the date. Could you say it again?", "ja": "日付がよく分かりませんでした。もう一度お願いします。"},
	"err_past":          {"zh": "這天已經過了喔，請選今天以後的日期～", "en": "That date has passed. Please pick a date from today onward~", "ja": "その日は過ぎています。今日以降の日付をお選びください～"},
	"err_time":          {"zh": "時間看不太懂，可以再說一次嗎？", "en": "Sorry, I didn't catch the time. Could you say it again?", "ja": "時間がよく分かりませんでした。もう一度お願いします。"},
	"err_people":        {"zh": "請問幾位呢？", "en": "How many people?", "ja": "何名様ですか？"},
	"err_phone":         {"zh": "電話號碼好像不完整，可以再確認一下嗎？", "en": "The phone number looks incomplete. Could you double-check it?", "ja": "電話番号が不完全なようです。もう一度ご確認いただけますか？"},
	"people_unit":       {"zh": " 位", "en": " people", "ja": "名"},
	"booking_done":      {"zh": "收到您的預約申請：{summary}。店員確認後會再通知您喔！", "en": "Got your reservation request: {summary}. We'll let you know once our staff confirms it!", "ja": "ご予約リクエストを承りました：{summary}。スタッフが確認後ご連絡します！"},
	"booking_confirmed": {"zh": "您的預約已確認：{summary}。期待您的光臨！", "en": "Your reservation is confirmed: {summary}. We look forward to seeing you!", "ja": "ご予約が確定しました：{summary}。ご来店をお待ちしております！"},
	"booking_rejected":  {"zh": "很抱歉，這次的預約無法安排：{summary}。", "en": "Sorry, we can't accommodate this reservation: {summary}.", "ja": "申し訳ございません、今回のご予約はお受けできませんでした：{summary}。"},
	"booking_cancel":    {"zh": "好的，已取消這次預約～", "en": "OK, I've cancelled this reservation request~", "ja": "かしこまりました、今回のご予約は取り消しました～"},
	"bot_identity":      {"zh": "我是{merchant}的 AI 客服～可以回答店家的問題，或協助預約喔！", "en": "I'm the AI assistant of {merchant}~ I can answer questions about us or help with reservations!", "ja": "{merchant}の AI アシスタントです～お店についてのご質問やご予約のお手伝いができます！"},
	"you_are":           {"zh": "您是{name}～", "en": "You're {name}~ ", "ja": "{name}様ですね～"},
	"ask_name_soft":     {"zh": "還不知道怎麼稱呼您呢？", "en": " May I have your name?", "ja": "お名前を伺ってもよろしいですか？"},
	"self_intro":        {"zh": "{name}您好～很高興認識您！有什麼可以幫您的嗎？", "en": "Nice to meet you, {name}! How can I help you today?", "ja": "{name}様、はじめまして！何かお手伝いできることはありますか？"},
}

func init() {
	// 打招呼的內建文字就是後台的預設文字（商家改過的在 cfg.Replies，textIn 會優先用）
	greeting, _ := merchant.EditableReplyByKey("greeting")
	builtinReplies["greeting"] = greeting.Default
}

// text 是固定回覆的文字，用這位顧客的語言
func (t *turn) text(key string, vars map[string]string) string {
	return t.textIn(key, t.customer.Lang, vars)
}

// textIn 是固定回覆的文字：商家設定的（資料庫）優先，沒有才用內建的。
// 沒有 lang 這種語言就用 *，再沒有就用商家的慣用語言
func (t *turn) textIn(key, lang string, vars map[string]string) string {
	table := map[string]string{}
	for _, r := range t.cfg.Replies {
		if r.Key == key {
			table[r.Lang] = r.Text
		}
	}
	if len(table) == 0 {
		table = builtinReplies[key]
	}
	if len(table) == 0 {
		table = t.askCustomField(key)
	}
	s, ok := table[lang]
	if !ok {
		s, ok = table["*"]
	}
	if !ok {
		s = table[t.cfg.Merchant.PreferredLanguage]
	}
	options := strings.Split(s, "|")
	s = options[rand.IntN(len(options))]
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{"+k+"}", v)
	}
	return strings.ReplaceAll(s, "{merchant}", t.cfg.Merchant.Name)
}

// askCustomFieldTemplates 是追問自訂欄位（ask_<代號>）的預設句子，{field} 換成欄位名稱；
// 匯入檔的固定回覆寫了 ask_<代號> 就用那個（後台的固定回覆頁沒有開放改）
var askCustomFieldTemplates = map[string]string{"zh": "請問{field}是？", "en": "Could you tell me your {field}?", "ja": "{field}をお知らせください。"}

// askCustomField：key 是追問自訂欄位（ask_<代號>）時的預設句子；不是就回 nil
func (t *turn) askCustomField(key string) map[string]string {
	field, ok := strings.CutPrefix(key, "ask_")
	spec := bookingSpecOf(t.cfg)
	if !ok || !spec.isCustom(field) {
		return nil
	}
	out := map[string]string{}
	for lang, s := range askCustomFieldTemplates {
		out[lang] = strings.ReplaceAll(s, "{field}", spec.names[field])
	}
	return out
}

func (t *turn) reply(key string) Reply {
	return Reply{Key: key, Text: t.text(key, nil)}
}
