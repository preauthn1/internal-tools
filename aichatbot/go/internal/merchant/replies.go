package merchant

// 後台可以改的固定回覆：商家改慣用語言版本，其他語言由 AI 翻譯（translate.go），後台可以逐語言修正

import (
	"context"
	"database/sql"
	"fmt"
	"maps"
	"slices"
	"strings"
)

// EditableReply 是後台可以改的一則固定回覆：Labels 是後台顯示的白話名稱（key 是後台介面語言），
// Default 是還沒設定時的預設文字（不含店名，key 是顧客語言）
type EditableReply struct {
	Key     string            `json:"key"`
	Labels  map[string]string `json:"labels"`
	Default map[string]string `json:"-"`
}

// EditableReplies 是後台可以改的固定回覆，順序就是後台顯示的順序。其他固定回覆（預約流程的問句、錯誤提示…）不開放改
var EditableReplies = []EditableReply{
	{"greeting", map[string]string{"zh-TW": "打招呼", "en-US": "Greeting", "ja-JP": "あいさつ"},
		map[string]string{"zh": "您好～有什麼可以幫您的嗎？", "en": "Hi there~ How can I help you?", "ja": "こんにちは～何かお手伝いできることはありますか？"}},
	{"booking_off", map[string]string{"zh-TW": "預約關閉", "en-US": "Booking turned off", "ja-JP": "予約受付停止中"},
		map[string]string{"zh": "目前沒有開放線上預約，可以透過這個方式聯絡我們：{contact}", "en": "Online booking isn't available right now. You can reach us here: {contact}", "ja": "現在オンライン予約は受け付けておりません。こちらからご連絡ください：{contact}"}},
	{"handoff", map[string]string{"zh-TW": "轉人工提示", "en-US": "Handing off to staff", "ja-JP": "スタッフへの引き継ぎ"},
		map[string]string{"zh": "我先幫您轉給店員，請稍候喔～", "en": "Let me connect you with our staff, please wait a moment~", "ja": "スタッフにおつなぎしますので、少々お待ちください～"}},
	{"handoff_off", map[string]string{"zh-TW": "轉人工關閉", "en-US": "Staff chat turned off", "ja-JP": "スタッフ対応停止中"},
		map[string]string{"zh": "目前沒有真人客服，可以透過這個方式聯絡我們：{contact}", "en": "We don't have live staff chat right now. You can reach us here: {contact}", "ja": "現在スタッフによる対応は行っておりません。こちらからご連絡ください：{contact}"}},
	{"off_hours", map[string]string{"zh-TW": "時段外提示", "en-US": "Outside staff hours", "ja-JP": "対応時間外のお知らせ"},
		map[string]string{"zh": "現在不是店員服務時間，店員上班後會盡快回覆您喔！", "en": "Our staff are off right now. They'll get back to you as soon as they're back!", "ja": "ただいまスタッフの対応時間外です。対応時間になりしだいご返信します！"}},
}

// ReplyBlock 是一則固定回覆：Texts 是語言代碼 → 句子（存的時候用 | 接起來，機器人隨機挑一句）
type ReplyBlock struct {
	Key   string              `json:"key"`
	Texts map[string][]string `json:"texts"`
}

// RepliesOf 取出後台可以改的固定回覆，只列啟用的語言。lang 是 * 的那筆每種語言都有；
// 慣用語言還沒設定時用預設文字（只有中英日有），其他語言空白（存檔時翻譯）
func RepliesOf(c Config) []ReplyBlock {
	stored := map[string]map[string]string{}
	for _, r := range c.Replies {
		if stored[r.Key] == nil {
			stored[r.Key] = map[string]string{}
		}
		stored[r.Key][r.Lang] = r.Text
	}
	out := make([]ReplyBlock, 0, len(EditableReplies))
	for _, e := range EditableReplies {
		b := ReplyBlock{Key: e.Key, Texts: map[string][]string{}}
		texts := stored[e.Key]
		for _, l := range c.Languages {
			if !l.Enabled {
				continue
			}
			text, found := texts[l.Code]
			if !found {
				text, found = texts["*"]
			}
			if !found && l.Code == c.Merchant.PreferredLanguage {
				text, found = e.Default[l.Code]
			}
			b.Texts[l.Code] = splitReply(text, found)
		}
		out = append(out, b)
	}
	return out
}

// SaveReply 只換掉一則固定回覆：texts 是語言代碼 → 句子（去掉前後空白，用 | 接起來）。
// 只換掉 texts 裡的語言；其他語言（沒啟用的）和 * 那筆不動（沒有自己那筆的語言還是用 *）。
// 只能改 EditableReplies、只能存啟用的語言；慣用語言至少一句；句子不能是空白、不能含 |
func SaveReply(ctx context.Context, db *sql.DB, key string, texts map[string][]string) error {
	return saveBlock(ctx, db, false, func(c *Config) error {
		if err := CheckReply(*c, key, texts); err != nil {
			return err
		}
		replies := []Reply{}
		for _, r := range c.Replies {
			if _, replaced := texts[r.Lang]; r.Key != key || !replaced {
				replies = append(replies, r)
			}
		}
		for lang, sentences := range texts {
			trimmed := make([]string, len(sentences))
			for i, s := range sentences {
				trimmed[i] = strings.TrimSpace(s)
			}
			replies = append(replies, Reply{Key: key, Lang: lang, Text: strings.Join(trimmed, "|")})
		}
		c.Replies = replies
		return nil
	}, func(tx *sql.Tx, c Config) error {
		if _, err := tx.ExecContext(ctx, "DELETE FROM fixed_reply WHERE key = ?", key); err != nil {
			return err
		}
		for _, r := range c.Replies {
			if r.Key != key {
				continue
			}
			if _, err := tx.ExecContext(ctx, "INSERT INTO fixed_reply (key, lang, text) VALUES (?, ?, ?)", r.Key, r.Lang, r.Text); err != nil {
				return err
			}
		}
		return nil
	})
}

// EditableReplyByKey 找 key 對應的可以改的固定回覆
func EditableReplyByKey(key string) (EditableReply, bool) {
	for _, e := range EditableReplies {
		if e.Key == key {
			return e, true
		}
	}
	return EditableReply{}, false
}

// CheckReply 是 SaveReply 存檔前的檢查（後台翻譯前也先檢查一次，沒過的不必等 AI）
func CheckReply(c Config, key string, texts map[string][]string) error {
	if _, ok := EditableReplyByKey(key); !ok {
		return &InvalidError{[]string{"不能改的固定回覆：" + key}}
	}
	var errs []string
	enabled := map[string]bool{}
	for _, l := range c.Languages {
		enabled[l.Code] = l.Enabled
	}
	langs := slices.Sorted(maps.Keys(texts))
	for _, lang := range langs {
		if !enabled[lang] {
			errs = append(errs, fmt.Sprintf("語言 %s 沒有啟用", lang))
			continue
		}
		for i, s := range texts[lang] {
			if strings.TrimSpace(s) == "" {
				errs = append(errs, fmt.Sprintf("%s 的第 %d 句是空白", lang, i+1))
			} else if strings.Contains(s, "|") {
				errs = append(errs, fmt.Sprintf("%s 的第 %d 句不能有 |", lang, i+1))
			}
		}
	}
	if p := c.Merchant.PreferredLanguage; len(texts[p]) == 0 {
		errs = append(errs, fmt.Sprintf("慣用語言（%s）至少要一句", p))
	}
	if len(errs) > 0 {
		return &InvalidError{errs}
	}
	return nil
}

// splitReply 把 | 分隔的多句拆開；沒有這個語言時是空的
func splitReply(text string, found bool) []string {
	if !found {
		return []string{}
	}
	return strings.Split(text, "|")
}
