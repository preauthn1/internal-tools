package dialog

import (
	"regexp"
	"strings"

	"aichatbot/internal/merchant"
)

// botCommandPattern 是 Telegram 的 bot 指令：/英文字母開頭的英數字，可以接 @bot名稱，後面接空白和參數
// （「/2」「/3 人」是顧客在選選項或回人數，不算指令）
var botCommandPattern = regexp.MustCompile(`^/[A-Za-z][A-Za-z0-9_]{0,31}(@[A-Za-z0-9_]+)?(\s|$)`)

// IsBotCommand：顧客送的是 bot 指令（/start、/help）。平台層攔下來改呼叫 Greet，不送進分類（#54）
func IsBotCommand(text string) bool {
	return botCommandPattern.MatchString(strings.TrimSpace(text))
}

// openerPattern 是開場白（#56）：還沒說要問什麼，只是先叫人（有人在嗎、我想問你個事情、請問一下）。
// 比對 merchant.Normalize 後的整句（去掉標點、空白、表情，英文小寫），前面可以先打招呼，可以連著好幾句（請問有人在嗎），後面可以接語助詞。
// 「老闆在嗎」可能是要找特定的人，不算
var openerPattern = regexp.MustCompile(`^(你好|您好|哈囉|嗨|hi|hello)*` +
	`(請問一下|請問|我?(想|可以)?(跟|和)?(你|妳|您)?(請問|問|請教)(你|妳|您)?(一下|個?事情?|一?件事情?|一?個問題)?|有人在嗎?|有人嗎|在嗎|在不在|` +
	`(is)?any(one|body)there|quickquestion|ihaveaquestion|caniaskyousomething|すみません|ちょっと(聞きたいんですが|いいですか)|質問があります)+` +
	`(喔|哦|耶|欸|呢|啊|呀|嗎)*$`)

// isOpener：整句是開場白（#56）
func isOpener(text string) bool {
	return openerPattern.MatchString(merchant.Normalize(text))
}

// symbolsOnly：整句只有符號、標點、顏文字，沒有任何文字或數字（＝＝、…、^^、🙂），當成應答（#63）
func symbolsOnly(text string) bool {
	return strings.TrimSpace(text) != "" && merchant.Normalize(text) == ""
}

// MARK: 寒暄的種類（#54）：一般對話已經確定是寒暄時，用常見的短詞判斷是打招呼還是應答，回對應的固定回覆

// 英文詞要整個字（this 不算 hi）
var (
	goodbyePattern  = regexp.MustCompile(`(?i)掰掰|掰囉|掰啦|拜拜|再見|晚安|\b(bye|goodbye|see you|good night)\b|さようなら|またね|じゃあね|おやすみ`)
	thanksPattern   = regexp.MustCompile(`(?i)謝|感謝|感恩|\b(thanks?|thx)\b|ありがと`)
	greetingPattern = regexp.MustCompile(`(?i)你好|您好|哈囉|嗨|安安|早安|午安|\b(hi|hello|hey|good (morning|afternoon|evening))\b|こんにちは|こんばんは|おはよう`)
)

// chatKind 是寒暄回哪一則固定回覆：opener（開場白，#56）、goodbye（道別）、thanks（道謝）、greeting（打招呼），其他是 general（收到）。
// 同一句符合好幾種時，優先順序是開場白、道別、道謝、打招呼（「你好，請問一下」是開場白，「謝謝掰掰」是道別）。之後要改成由 AI 判斷，只動這個函式
func chatKind(text string) string {
	switch {
	case isOpener(text):
		return "opener"
	case goodbyePattern.MatchString(text):
		return "goodbye"
	case thanksPattern.MatchString(text):
		return "thanks"
	case greetingPattern.MatchString(text):
		return "greeting"
	}
	return "general"
}
