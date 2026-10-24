package merchant

// Intent 是一個意圖：Code 是分類用的代號，Labels 是後台顯示的白話名稱（key 是後台介面語言 zh-TW、en-US、ja-JP）
type Intent struct {
	Code   string            `json:"code"`
	Labels map[string]string `json:"labels"`
}

// Intents 是 11 個內建意圖（CONTEXT.md），順序就是後台顯示的順序，和對話層（internal/dialog）的分類選項相同。
// general（一般對話）是 Swift 版的打招呼、道謝、道別、閒聊、問商家合併；human_request（想找真人）是 Go 版新增的
var Intents = []Intent{
	{"booking", map[string]string{"zh-TW": "想預約", "en-US": "Wants to book", "ja-JP": "予約したい"}},
	{"cancel", map[string]string{"zh-TW": "想取消預約", "en-US": "Wants to cancel a booking", "ja-JP": "予約を取り消したい"}},
	{"follow_up", map[string]string{"zh-TW": "追問進度", "en-US": "Asking for an update", "ja-JP": "進み具合を聞いている"}},
	{"claim", map[string]string{"zh-TW": "聲稱（例如說店家答應過）", "en-US": "Claims something (e.g. staff promised it)", "ja-JP": "主張している（店が約束したなど）"}},
	{"identity", map[string]string{"zh-TW": "問身分（例如是不是真人）", "en-US": "Asks who they're talking to (e.g. a bot?)", "ja-JP": "相手が誰か聞いている（人間かどうかなど）"}},
	{"self_intro", map[string]string{"zh-TW": "自我介紹", "en-US": "Introduces themselves", "ja-JP": "自己紹介"}},
	{"human_request", map[string]string{"zh-TW": "想找真人", "en-US": "Wants a human", "ja-JP": "人と話したい"}},
	{"complaint", map[string]string{"zh-TW": "抱怨", "en-US": "Complaint", "ja-JP": "苦情"}},
	{"off_topic", map[string]string{"zh-TW": "和店家無關", "en-US": "Not about the shop", "ja-JP": "お店と関係ない"}},
	{"unclear", map[string]string{"zh-TW": "看不懂", "en-US": "Unclear", "ja-JP": "意味がわからない"}},
	{"general", map[string]string{"zh-TW": "一般對話（打招呼、道謝、閒聊、問店家）", "en-US": "General (greetings, thanks, chat, questions about the shop)", "ja-JP": "一般的な会話（あいさつ、お礼、雑談、お店への質問）"}},
}

// legacyIntents 是 Swift 版的舊代號，Go 版合併成一般對話（eval/examples.jsonl 還是舊代號，Swift 版在用）
var legacyIntents = map[string]string{
	"greeting": "general", "thanks": "general", "goodbye": "general", "small_talk": "general", "store_question": "general",
}

// IntentByCode 找代號對應的意圖
func IntentByCode(code string) (Intent, bool) {
	for _, in := range Intents {
		if in.Code == code {
			return in, true
		}
	}
	return Intent{}, false
}
