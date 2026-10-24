package dialog

import (
	"slices"
	"strings"

	"aichatbot/internal/merchant"
)

// MARK: 自訂意圖（CONTEXT.md「意圖」）：商家新增的意圖，各綁一種動作

// customIntent 找分類結果對應的自訂意圖
func (t *turn) customIntent(code string) (merchant.CustomIntent, bool) {
	for _, ci := range t.cfg.CustomIntents {
		if ci.Code == code {
			return ci, true
		}
	}
	return merchant.CustomIntent{}, false
}

// custom 依自訂意圖的動作回覆；Reply.Key 是 custom:代號（轉人工關閉、答不出來時是 handoff_off、noAnswer 原本的 key）
func (t *turn) custom(ci merchant.CustomIntent, text string) Reply {
	key := "custom:" + ci.Code
	t.trace("custom", ci.Action, nil)
	switch ci.Action {
	case merchant.ActionHandoff:
		reason := "自訂意圖「" + ci.Name + "」：" + text
		// 轉人工關閉或沒綁員工群組：照原本的處理（回固定回覆 handoff_off），不必翻譯提示
		if !t.handoffAvailable() {
			return t.handoff(reason)
		}
		prompt := ""
		if p := strings.TrimSpace(ci.Text); p != "" {
			prompt = t.inCustomerLanguage(text, p)
		}
		r := t.handoffWithPrompt(reason, prompt)
		r.Key = key
		return r
	case merchant.ActionKnowledge:
		// 不挑 FAQ，只從指定的段落挑（ADR 0002 的挑知識 → 寫回覆 → 驗證原文片段）
		var sections []merchant.Section
		for _, s := range t.cfg.Knowledge {
			if slices.Contains(ci.Knowledge, strings.TrimSpace(s.Title)) {
				sections = append(sections, s)
			}
		}
		picked := t.pickKnowledge(text, sections)
		if len(picked) == 0 {
			return t.noAnswer("找不到相關資料：" + text)
		}
		r, ok := t.writeAnswer(text, picked)
		if !ok {
			return t.noAnswer("資料裡找不到答案：" + text)
		}
		r.Key = key
		return r
	default:
		return Reply{Key: key, Text: t.inCustomerLanguage(text, ci.Text)}
	}
}

// inCustomerLanguage 是商家用慣用語言寫的一段文字：顧客用慣用語言（或看不出語言）就回原文；
// 其他語言翻成顧客的語言，譯文要通過驗證（同 FAQ 答案，見 translateVerified），沒過就回原文
func (t *turn) inCustomerLanguage(message, source string) string {
	if translated, ok := t.translateVerified(message, source); ok {
		return translated
	}
	return source
}
