package dialog

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

// MARK: 分派：狀態 × 意圖 → 怎麼回

func (t *turn) handle(text string) Reply {
	t.filled = filledFields(t.customer.Booking)
	if t.customer.Human {
		return Reply{Key: "silent", Silent: true}
	}
	complained := t.customer.Complained
	t.customer.Complained = false // 只有回了道歉的下一句算
	// 預約問欄位時的短答（「Bill」「たなか」「김민수」）不切換語言，完整句子才跟著換
	if lang := DetectLanguage(text); lang != "" && (t.customer.AskingField == "" || !isShortAnswer(text)) {
		t.customer.Lang = lang
	}
	if enabled := t.languages(); !slices.Contains(enabled, t.customer.Lang) {
		t.openTicket("unsupported", fmt.Sprintf("🌐 不支援的語言（%s）：%s", t.customer.Lang, text), true, 0)
		return Reply{Key: "unsupported", Text: t.unsupported(enabled)}
	}
	// 只有表情、標點（沒有文字和數字）不用問 AI
	intent, name := "general", ""
	if merchant.Normalize(text) != "" {
		var failed string
		intent, name, failed = t.classifyWith(text)
		var detail map[string]string
		if failed != "" {
			detail = map[string]string{"error": failed}
		}
		t.trace("intent", intent, detail)
	}
	t.customer.AskingName = false // 只有問稱呼的下一句算

	// 第一次知道名字就記住（之後不自動覆蓋，避免把朋友的名字當成顧客）
	justIntroduced := false
	if t.customer.Name == "" && name != "" {
		t.customer.Name, justIntroduced = name, true
	}

	// 預約中：顧客多半是在回答剛問的欄位；明確是別的事（追問、聲稱、想找真人）才跳出去處理
	// （身分類不跳出：正在問名字時只回「小康」，常被分成身分）
	if t.customer.AskingField != "" {
		if intent == "cancel" {
			t.customer.endBooking()
			return t.reply("booking_cancel")
		}
		// 只有符號（…、^^）沒有可抽的欄位：不呼叫 AI 抽取，由一般對話再問一次欄位（#63）
		if !slices.Contains([]string{"follow_up", "claim", "human_request"}, intent) && !symbolsOnly(text) {
			if r, ok := t.bookingStep(text, intent != "booking", intent); ok {
				return r
			}
		}
	}

	t.streakBefore = t.customer.OffTopicStreak
	if intent == "general" || intent == "booking" {
		t.customer.OffTopicStreak = 0
	}
	if ci, ok := t.customIntent(intent); ok {
		t.customer.OffTopicStreak = 0
		return t.custom(ci, text)
	}

	switch intent {
	case "booking":
		// 預約關閉：回固定回覆 booking_off（帶入商家的聯絡方式），不轉人工
		if !t.feature("booking") {
			return Reply{Key: "booking_off", Text: t.text("booking_off", map[string]string{"contact": t.cfg.Merchant.Contact})}
		}
		if r, ok := t.bookingStep(text, true, intent); ok {
			return r
		}
		return t.bookingNext(nil)
	case "cancel":
		return t.general(text) // 沒有進行中的預約，當成一般對話
	case "follow_up":
		if !t.customer.OpenTicket {
			return t.reply("unclear")
		}
		t.customer.FollowUps++
		if t.customer.FollowUps == 3 {
			// 引用還在等員工的那張（最近開的、不是緊急通知），員工回覆緊急通知時對應回那張
			msg := fmt.Sprintf("🔥 緊急：顧客已追問 %d 次，請盡快回覆", t.customer.FollowUps)
			var related int64
			if w := t.waitingTicket(); w != nil {
				related = w.ID
				original := w.Original
				if original == "" {
					original = w.Text
				}
				msg += fmt.Sprintf("\n📌 原本的事項 #%d：%s", w.ID, original)
			}
			t.openTicket("urgent", msg, false, related)
		}
		if t.customer.FollowUps >= 3 {
			return t.reply("nudge")
		}
		return t.reply("waiting")
	case "claim", "off_topic":
		return t.offTopic(intent, text)
	case "human_request":
		return t.handoff("顧客要求找真人")
	case "identity":
		if n := t.customer.Name; n != "" {
			return Reply{Key: "identity", Text: t.text("you_are", map[string]string{"name": n}) + t.text("bot_identity", nil)}
		}
		// 同一位顧客只追問一次稱呼，不連續回同一句（#61）
		if t.customer.AskedName {
			return Reply{Key: "identity", Text: t.text("bot_identity", nil)}
		}
		t.customer.AskedName, t.customer.AskingName = true, true
		return Reply{Key: "identity", Text: t.text("bot_identity", nil) + t.text("ask_name_soft", nil)}
	case "self_intro":
		if t.customer.Name == "" {
			return t.general(text)
		}
		return t.selfIntro()
	case "complaint":
		// 上一句已經道歉過，再抱怨就轉人工，不重複回同一句道歉（#62）
		if complained {
			return t.handoff("顧客連續抱怨")
		}
		t.customer.Complained = true
		return t.reply("complaint")
	case "unclear":
		t.unanswered("unclear")
		return t.reply("unclear")
	default:
		// 這句話介紹了名字（「我是小康」常被分成一般對話）：當成自我介紹
		if justIntroduced {
			return t.selfIntro()
		}
		return t.general(text)
	}
}

// offTopic 是聲稱（claim）或無關請求（off_topic）：計入連續次數，到上限轉人工，否則回固定回覆。
// 意圖分類和一般對話的確認（checkQuestion 判成無關）共用
func (t *turn) offTopic(intent, text string) Reply {
	t.customer.OffTopicStreak++
	if intent == "claim" {
		t.openTicket("claim", "❓ 顧客聲稱／更正："+text, true, 0)
	}
	limit := t.handoffCfg.StreakLimit
	if limit <= 0 { // 還沒有後台檢查，存錯時用預設值
		limit = store.DefaultHandoffSetting().StreakLimit
	}
	if t.customer.OffTopicStreak >= limit {
		return t.handoff(fmt.Sprintf("連續 %d 次無關或聲稱的訊息", t.customer.OffTopicStreak))
	}
	return t.reply(intent)
}

// languages 是商家啟用的顧客語言
func (t *turn) languages() []string {
	var out []string
	for _, l := range t.cfg.Languages {
		if l.Enabled {
			out = append(out, l.Code)
		}
	}
	return out
}

// unsupported 是「目前只支援…」：每種啟用的語言各一行（內建只有中英日的說法，都沒啟用時用英文）
func (t *turn) unsupported(enabled []string) string {
	var names, lines []string
	for _, code := range enabled {
		names = append(names, languageNames[code])
	}
	vars := map[string]string{"languages": strings.Join(names, "、")}
	for _, code := range enabled {
		if _, ok := builtinReplies["unsupported"][code]; ok {
			lines = append(lines, t.textIn("unsupported", code, vars))
		}
	}
	if len(lines) == 0 {
		lines = append(lines, t.textIn("unsupported", "en", vars))
	}
	return strings.Join(lines, "\n")
}

// feature 是功能開關（booking、handoff）有沒有開；設定裡沒有就當成開著
func (t *turn) feature(name string) bool { return merchant.FeatureEnabled(t.cfg, name) }

// waitingTicket 是這位顧客最近開的、還沒處理的待處理事項（不含緊急通知）；沒有是 nil
func (t *turn) waitingTicket() *store.Ticket {
	ts, err := t.bot.DB.OpenTickets(t.ctx)
	if err != nil {
		if t.err == nil {
			t.err = err
		}
		return nil
	}
	for i := len(ts) - 1; i >= 0; i-- {
		if x := ts[i]; x.Platform == t.msg.Platform && x.UserID == t.msg.UserID && x.Kind != "urgent" {
			return &x
		}
	}
	return nil
}

// hasOpenHandoffTicket：這位顧客有還沒處理的轉人工事項
func (t *turn) hasOpenHandoffTicket() bool {
	ts, err := t.bot.DB.OpenTickets(t.ctx)
	if err != nil {
		if t.err == nil {
			t.err = err
		}
		return false
	}
	return slices.ContainsFunc(ts, func(x store.Ticket) bool {
		return x.Platform == t.msg.Platform && x.UserID == t.msg.UserID && x.Kind == "handoff"
	})
}

// handoff 轉人工：之後 AI 不回覆這位顧客，開一張轉人工的待處理事項給員工。
// 轉人工關閉、或還沒綁員工群組時不轉，回固定回覆 handoff_off（帶入商家的聯絡方式）
func (t *turn) handoff(reason string) Reply { return t.handoffWithPrompt(reason, "") }

// handoffAvailable：轉人工開著而且綁了員工群組。沒綁群組時沒人收得到通知、按不了結束，照轉人工關閉處理（#55）
func (t *turn) handoffAvailable() bool { return t.feature("handoff") && t.handoffCfg.StaffChatID != 0 }

// handoffWithPrompt 是 handoff，prompt 不是空的就取代轉人工的提示（自訂意圖的提示）
func (t *turn) handoffWithPrompt(reason, prompt string) Reply {
	if !t.handoffAvailable() {
		// 轉人工開著、只是沒綁群組：照樣開一張事項，後台看得到；沒有人會回，不算在等員工。
		// 已經有一張還沒處理的就不再開（之後每句無關、答不出來都會走到這裡）
		if t.feature("handoff") && !t.hasOpenHandoffTicket() {
			t.openTicket("handoff", "🙋 想找真人（沒有員工群組，未轉人工）："+reason, false, 0)
		}
		return Reply{Key: "handoff_off", Text: t.text("handoff_off", map[string]string{"contact": t.cfg.Merchant.Contact})}
	}
	t.customer.Human, t.customer.HumanSince, t.customer.LastStaffReply = true, t.now, time.Time{}
	t.customer.HandoffTicket = t.openTicket("handoff", "🙋 轉人工："+reason, false, 0)
	r := t.reply("handoff")
	if p := strings.TrimSpace(t.handoffCfg.PromptText); p != "" {
		r.Text = p
	}
	if prompt != "" {
		r.Text = prompt
	}
	// 員工服務時段外照樣轉人工，告訴顧客上班後回覆
	if !inServiceHours(t.handoffCfg.ServiceHours, t.now) {
		r.Text += "\n" + t.text("off_hours", nil)
	}
	return r
}

// noAnswer：機器人答不出來時（一般對話挑不到知識、AI 答不出來或驗證沒過），
// 轉人工設定的「找不到答案時轉人工」開著就轉人工，關著回看不懂
func (t *turn) noAnswer(reason string) Reply {
	t.unanswered("unanswerable")
	if t.handoffCfg.OnNoAnswer {
		return t.handoff(reason)
	}
	return t.reply("unclear")
}

func (t *turn) selfIntro() Reply {
	return Reply{Key: "self_intro", Text: t.text("self_intro", map[string]string{"name": t.customer.Name})}
}
