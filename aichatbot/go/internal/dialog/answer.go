package dialog

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
)

// MARK: 一般對話（ADR 0002 的混合式）：挑 FAQ → 挑不到才挑知識（選擇題）→ AI 寫回覆 → 程式驗證原文片段

// Trace 是一步 AI 判斷或程式檢查的結果，平台層寫進對話紀錄的 traces。
// Step：intent（分類）、faq（挑 FAQ）、translate（翻譯 FAQ 答案或自訂意圖的文字）、pick（挑知識）、ask（挑不到知識時確認是不是在問店家）、
// answer（寫回覆）、verify（驗證）、reask（預約中再問欄位）、custom（自訂意圖）、chat（寒暄的種類）；
// Result：意圖、挑到的 FAQ 代號（或 price_calc、unclear、not_a_question、none）、翻譯結果（ok、failed 或擋下的原因）、
// 挑到的段落（+ 分隔，沒挑到是 none）、chat、unclear、off_topic 或 question、answerable（yes、no，AI 失敗是 ai_failed、提示詞超過上下文是 too_long、被安全過濾擋下是 blocked_by_safety）、驗證結果（ok 或擋下的原因）、再問的欄位、
// 自訂意圖的動作（fixed_text、handoff、knowledge）、寒暄回的固定回覆（opener、greeting、thanks、goodbye、general）
type Trace struct {
	Step   string            `json:"step"`
	Result string            `json:"result"`
	Detail map[string]string `json:"detail,omitempty"`
}

func (t *turn) trace(step, result string, detail map[string]string) {
	t.traces = append(t.traces, Trace{Step: step, Result: result, Detail: detail})
}

// failure 是 AI 失敗時 trace 記的原因：提示詞超過本機模型上下文是 model.FailureContextExceeded（#65），
// 被安全過濾擋下是 model.FailureBlockedBySafety（#68），其他是 ai_failed
func failure(u model.Usage) string {
	if u.Failure != "" {
		return u.Failure
	}
	return "ai_failed"
}

// knowledgeItem 是挑知識清單的一項（K1…）：label 是標題，body 是寫回覆時給 AI 的全文
type knowledgeItem struct {
	id, label, body string
}

// general 是一般對話（ADR 0002 的混合式）：先挑 FAQ，挑到就回商家寫好的答案；挑不到而且是問商家，才挑知識段落請 AI 寫回覆。
// 商家沒有知識段落也沒有 FAQ、開場白（#56，預約中除外）、或只有符號（#63）時當成寒暄，不呼叫 AI。多件、多種商品的價格（程式認得的中文說法，或挑 FAQ 時選了算價）當成答不出來的問商家（#52）。
// 挑不到知識時先確認是不是在問店家（checkQuestion），和店家無關的照無關請求處理（計入連續無關次數）。答不出來時：預約中再問一次正在問的欄位；寒暄依種類回固定回覆（chatKind，#54），看不懂回看不懂，問商家的走 noAnswer（轉人工或回看不懂）
func (t *turn) general(text string) Reply {
	kind := "chat"
	fallback := func(reason string) Reply {
		switch kind {
		case "unclear":
			t.unanswered("unclear")
		case "question":
			t.unanswered("unanswerable")
		}
		// 預約中（正在問某個欄位）：再問一次那個欄位，預約接得下去（照 Swift 版）；同一個欄位問太多次就轉人工（#53）
		if f := t.customer.AskingField; f != "" {
			t.trace("reask", f, nil)
			return t.askField(f, "ask_"+f)
		}
		switch kind {
		case "chat":
			// 寒暄依種類回對應的固定回覆（#54），回覆種類維持 general，種類記在 trace
			sub := chatKind(text)
			t.trace("chat", sub, nil)
			return Reply{Key: "general", Text: t.text(sub, nil)}
		case "unclear":
			return t.reply("unclear")
		}
		return t.noAnswer(reason + "：" + text)
	}
	// 開場白還沒說要問什麼：不挑 FAQ、不確認是不是在問店家，直接回「請說」（#56）。
	// 預約中照舊再問欄位；有事在等店員時分類會問 AI，AI 判成一般對話的開場白一樣回「請說」。
	// 只有符號（＝＝、…）也不挑 FAQ：回「收到」，預約中再問欄位（#63）
	if (len(t.cfg.Knowledge) == 0 && len(t.cfg.FAQ) == 0) || (t.customer.AskingField == "" && isOpener(text)) || symbolsOnly(text) {
		return fallback("")
	}
	var f *merchant.FAQ
	if multiItemPrice(text) {
		// 常見的中文多件問法不呼叫挑 FAQ：AI 常挑成單價的 FAQ（#52）
		t.trace("faq", "price_calc", map[string]string{"rule": "multi_item"})
		kind = "price_calc"
	} else {
		kind, f = t.pickFAQ(text) // fallback 看的是這裡決定的 kind
	}
	if f != nil {
		if r, ok := t.faqReply(text, *f); ok {
			return r
		}
	}
	if kind == "price_calc" {
		// 多件、多種商品的價格要算總價：照設計當成答不出來的問商家（#44，Go 版不由程式算總價），
		// 不挑知識段落，免得 AI 只回一件的價格（#52）
		kind = "question"
		return fallback("需要計算總價")
	}
	if kind != "question" {
		return fallback("")
	}
	picked := t.pickKnowledge(text, t.cfg.Knowledge)
	if len(picked) == 0 {
		kind = t.checkQuestion(text) // 不是在問店家時 fallback 改回寒暄或看不懂，不轉人工
		if kind == "off_topic" {
			// 和店家無關：照無關請求處理（預約中也不再問欄位、不記待補問題），接上這則訊息進來時的次數
			t.customer.OffTopicStreak = t.streakBefore
			return t.offTopic("off_topic", text)
		}
		return fallback("找不到相關資料")
	}
	r, ok := t.writeAnswer(text, picked)
	if !ok {
		return fallback("資料裡找不到答案")
	}
	return r
}

// pickFAQ 照 Swift 版（swift-final 的 StoreQuestion.swift、docs/研究總結.md 第 12 節）挑一題 FAQ：選項附上答案內容，
// 另外有算價、看不懂、不是在問、都不能回答。回傳寒暄（chat）、看不懂（unclear）、要算總價（price_calc）或問商家（question）和挑到的 FAQ（沒挑到是 nil）。
// 挑到 FAQ、但 ask 抄的是那題的答案時，先用 checkQuestion 確認是在問店家才回 FAQ（#51）。AI 失敗時當成問商家、沒挑到
func (t *turn) pickFAQ(text string) (string, *merchant.FAQ) {
	faqs := t.cfg.FAQ
	detail := map[string]string{}
	// FAQ 太多時全部送出會超過本機模型上下文（#65）：先用句子向量挑出最像的幾題；沒有向量時照舊送全部
	if len(faqs) > merchant.FAQFilterAbove {
		if v := t.vector(text); v == nil {
			t.bot.log().Warn("FAQ 超過門檻題數但算不出訊息的向量，挑 FAQ 照舊送全部的題目", "faq", len(faqs))
		} else if top, ok := t.bot.similarFAQ(v, faqs, merchant.FAQFilterTop); ok {
			detail["filtered"] = fmt.Sprintf("%d/%d", len(top), len(faqs))
			faqs = top
		}
	}
	r, usage := t.bot.Model.Respond(t.ctx, FAQRequest(text, faqs))
	if r == nil {
		detail["error"] = failure(usage)
		t.trace("faq", "none", detail)
		return "question", nil
	}
	picked, _ := r.String("faq")
	kind, _ := r.String("kind")
	ask, _ := r.String("ask")
	detail["kind"], detail["ask"] = kind, ask
	t.trace("faq", picked, detail)
	i := slices.IndexFunc(faqs, func(f merchant.FAQ) bool { return f.Key == picked })
	if i >= 0 && askFromAnswer(ask, text, faqs[i].Answer) && t.checkQuestion(text) != "question" {
		// 簡短應答（知道了、OK 了解）沒有要問的事，AI 把選項的答案抄進 ask、挑了那題（#51）：不回 FAQ。
		// 真問句也會這樣抄（營業時間、Are you open on Monday?），所以只在這時多確認一次是不是在問店家；
		// 確認不是在問的（寒暄、看不懂、無關）一律當成寒暄，不轉人工
		return "chat", nil
	}
	switch {
	case i >= 0:
		return "question", &faqs[i]
	case picked == "price_calc":
		kind = "price_calc"
	case picked == "unclear":
		kind = "unclear"
	case picked == "not_a_question" || kind == "chat":
		kind = "chat"
	default:
		kind = "question"
	}
	return kind, nil
}

// FAQRequest 是挑 FAQ 送給 AI 的題目：faqs 每題附上問法和答案，另外有算價、看不懂、不是在問、都不能回答（量測工具 ctxprobe 也用它）
func FAQRequest(text string, faqs []merchant.FAQ) model.Request {
	var ids, options []string
	for _, f := range faqs {
		ids = append(ids, f.Key)
		// 附上答案內容：只看問法時 AI 會因主題相近而選（大杯、去冰→menu），看到答案才能判斷「答得出這件事嗎」
		options = append(options, merchant.FAQOption(f))
	}
	options = append(options,
		"price_calc：問兩杯以上、或同時問好幾樣飲品的價格，要算總價（不管有沒有說「總共」，例如「兩杯 X 多少」「A 和 B 各一杯多少」）；只問一杯的價格不算（自取、外送一杯也不算）",
		"unclear：看不懂顧客在問什麼（太短、沒頭沒尾、指代不明，例如「中」「那這個咧」「然後咧」）",
		"not_a_question：顧客不是在問問題（自我介紹、打招呼、回應、讓我想想、要客服扮演別人）",
		"none：以上都無法回答（店家資料沒有這項資訊，例如吸菸區、禮盒、有幾家店）")
	fields := []model.Field{
		// 先逐字抄出顧客問的那件事再選：只看整句時 AI 會被主題帶走（帶環保杯→自取優惠、停車在哪→地址）。欄位依序產生，ask 要放在 faq 前面
		{Kind: model.KindText, Name: "ask", Description: "顧客想知道的那件事，從訊息逐字照抄（不要改寫）"},
		{Kind: model.KindChoice, Name: "faq", Description: "選項的答案有正好講到上面那件事的代號", Options: slices.Concat(ids, []string{"price_calc", "unclear", "not_a_question", "none"})},
		// 寒暄還是問商家放在最後：不影響挑 FAQ（放在最前面時「好」被挑成營業時間，#45）；挑不到 FAQ 時決定要不要挑知識段落
		{Kind: model.KindChoice, Name: "kind", Description: "chat=沒有要問店家的事：打招呼、道謝、道別、回應、閒聊、表情符號、自我介紹、說自己在考慮、說誰推薦的\nquestion=想知道店家的某項資訊（營業時間、價格、停車、wifi、能不能做某件事）", Options: []string{"chat", "question"}},
	}
	return model.Request{
		Name: "FAQ", Text: text, Fields: fields,
		Instructions: "從選項中挑出能回答顧客問題的一個。顧客可能使用任何語言。" +
			"只有當選項的內容正好回答顧客問的那件事才選它；只是主題相關（同樣是飲料、折扣、到店、付款）不算，要選 none。" +
			"顧客問兩杯以上、或好幾樣飲品的價格時選 price_calc，不要選只寫一杯價格的選項。" +
			"問飲品的成分、甜度、配料或產地時，只寫價格的選項不算，要選 none。" +
			"不確定就選 none。\n選項：\n" + strings.Join(options, "\n"),
	}
}

var (
	// 第一組是「各」（各一杯＝好幾樣），第二組是「第」（第二杯是問優惠，不是數量），第三組是數量，第四組是量詞，
	// 第五組是「個」後面接的人數、時間（兩個人、兩個大人、三個月不是商品數量）
	priceQtyPattern  = regexp.MustCompile(`(各)?\s*(第)?([0-9０-９]+|[一二兩三四五六七八九十]+)\s*([杯份瓶件個本組盒包條張支雙套])(人|大人|小孩|小朋友|月|小時|晚上|禮拜|星期)?`)
	priceWordPattern = regexp.MustCompile(`多少(?:[^時分天人久]|$)|幾元|幾塊|價錢|價格|打折`) // 多少時間、多少人不是問價格
	oneQty           = []string{"1", "１", "一"}
)

// multiItemPrice：顧客用中文問多件（數量大於 1）、或好幾樣商品（各一件）的價格（#52，#64 補上飲料以外的量詞）。
// 挑 FAQ 時 AI 常把這種問法挑成單價的 FAQ，程式先擋下常見說法；其他語言靠挑 FAQ 的 price_calc 選項
func multiItemPrice(text string) bool {
	if !priceWordPattern.MatchString(text) {
		return false
	}
	for _, m := range priceQtyPattern.FindAllStringSubmatch(text, -1) {
		if m[2] == "" && !(m[4] == "個" && m[5] != "") && (m[1] != "" || !slices.Contains(oneQty, m[3])) {
			return true
		}
	}
	return false
}

// askFromAnswer：挑 FAQ 時 ask 不是從顧客的訊息抄的，而是抄了 FAQ 的答案
func askFromAnswer(ask, text, answer string) bool {
	ask = strings.TrimSpace(ask)
	return ask != "" && !strings.Contains(text, ask) && strings.Contains(answer, ask)
}

// checkQuestion 在問商家卻答不出來、要轉人工前，再確認一次顧客是不是在問店家的資訊：挑 FAQ 時 AI 會把閒聊（我想一下、😂）
// 判成問商家（#45）。挑到 FAQ、但 ask 抄的是答案時也用它確認（#51）。只在這兩條路上多呼叫一次，其他回答不受影響。回傳 chat、unclear、off_topic（和店家無關、但講得清楚，#48）或 question；
// AI 失敗時當成問商家（照舊轉人工）
func (t *turn) checkQuestion(text string) string {
	r, usage := t.bot.Model.Respond(t.ctx, model.Request{
		Name: "Ask", Text: text,
		Fields: []model.Field{{Kind: model.KindChoice, Name: "kind",
			Description: "asks_about_store=有講出想知道店家的哪件事（營業時間、價格、停車、wifi、老闆、有沒有某樣東西、能不能做某件事），不管店家資料有沒有\n" +
				"small_talk=沒有要問店家的事：打招呼、道謝、回應（表示收到、明白或同意的話，例如收到囉、我懂、曉得了、got it）、閒聊、表情符號、自我介紹、說自己在考慮（我想一下、等等）、說誰推薦的\n" +
				"off_topic=和店家無關、但講得清楚的請求或問題（時事、常識、算數、寫詩、寫程式、翻譯）\n" +
				"too_vague=只有「然後咧」「那這個咧」這種，沒講出是哪件事（和店家無關但講得清楚的不算）",
			Options: []string{"asks_about_store", "small_talk", "off_topic", "too_vague"}}},
		Instructions: "判斷顧客這句話是不是在問店家的資訊。顧客可能使用任何語言。只是表示收到、知道了的回應不是在問。",
	})
	if r == nil {
		t.trace("ask", "question", map[string]string{"error": failure(usage)})
		return "question"
	}
	got, _ := r.String("kind")
	kind := map[string]string{"small_talk": "chat", "too_vague": "unclear", "off_topic": "off_topic"}[got]
	if kind == "" {
		kind = "question"
	}
	t.trace("ask", kind, nil)
	return kind
}

// faqReply 是挑到的 FAQ 的回覆：顧客用商家慣用語言（或看不出語言）就回答案原文；其他語言翻成顧客的語言，
// 譯文要通過數字、時間、電話、網址、語言和長度的驗證（FAQ 答案就是原文片段）。ok 是 false 表示翻譯失敗或驗證沒過
func (t *turn) faqReply(text string, f merchant.FAQ) (Reply, bool) {
	translated, ok := t.translateVerified(text, f.Answer)
	if !ok {
		return Reply{}, false
	}
	return Reply{Key: "answer", Text: translated}, true
}

// translateVerified 把商家用慣用語言寫的 source 翻成顧客的語言（顧客用慣用語言或看不出語言時就是原文），
// 譯文要通過數字、時間、電話、網址、語言和長度的驗證（source 就是原文片段）。ok 是 false 表示翻譯失敗或驗證沒過
func (t *turn) translateVerified(message, source string) (string, bool) {
	from := t.cfg.Merchant.PreferredLanguage
	if t.customer.Lang == "" || t.customer.Lang == from {
		return source, true
	}
	translated, err := merchant.Translate(t.ctx, t.bot.Model, from, t.customer.Lang, source)
	if err != nil {
		var detail map[string]string
		if errors.Is(err, merchant.ErrTranslateBlocked) {
			detail = map[string]string{"error": model.FailureBlockedBySafety}
		}
		t.trace("translate", "failed", detail)
		return "", false
	}
	result := verifyAnswer(answerCheck{Source: source, Quote: source, Reply: translated, Message: message, Lang: t.customer.Lang})
	if result != "" {
		t.trace("translate", result, map[string]string{"reply": translated})
		t.verifyFailed = true
		return "", false
	}
	t.trace("translate", "ok", nil)
	return translated, true
}

var pickSlots = []string{"pick1", "pick2", "pick3"}

// pickKnowledge 從 sections 回傳挑到的知識段落（最多 3 段，去掉重複）。AI 失敗時當成沒挑到
func (t *turn) pickKnowledge(text string, sections []merchant.Section) []knowledgeItem {
	items := knowledgeItems(sections)
	ids := []string{}
	lines := []string{}
	for _, it := range items {
		ids, lines = append(ids, it.id), append(lines, it.id+" "+it.label)
	}
	if len(items) == 0 {
		return nil
	}
	options := slices.Concat(ids, []string{"none"})
	fields := []model.Field{
		// 先逐字抄出顧客問的那件事再選（Swift 版挑 FAQ 的研究：只看整句時 AI 會被主題帶走）
		{Kind: model.KindText, Name: "ask", Description: "顧客想知道的那件事，從訊息逐字照抄（不要改寫）"},
		{Kind: model.KindChoice, Name: "pick1", Description: "最能回答上面那件事的資料編號；沒有就選 none", Options: options},
		{Kind: model.KindChoice, Name: "pick2", Description: "回答還需要的另一段資料編號；不需要就選 none", Options: options},
		{Kind: model.KindChoice, Name: "pick3", Description: "回答還需要的另一段資料編號；不需要就選 none", Options: options},
	}
	r, usage := t.bot.Model.Respond(t.ctx, model.Request{
		Name: "Pick", Text: text, Fields: fields,
		Instructions: "從資料清單挑出能回答顧客這句話的資料，最多 3 段。顧客可能使用任何語言。" +
			"只有資料的內容正好回答顧客問的那件事才選它；只是主題相關不算。不確定就選 none。\n資料清單：\n" + strings.Join(lines, "\n"),
	})
	if r == nil {
		t.trace("pick", "none", map[string]string{"error": failure(usage)})
		return nil
	}
	var picked []knowledgeItem
	var got []string
	for _, slot := range pickSlots {
		id, _ := r.String(slot)
		i := slices.Index(ids, id)
		if i < 0 || slices.Contains(got, id) {
			continue
		}
		picked, got = append(picked, items[i]), append(got, id)
	}
	result := strings.Join(got, "+")
	if result == "" {
		result = "none"
	}
	ask, _ := r.String("ask")
	t.trace("pick", result, map[string]string{"ask": ask})
	return picked
}

// writeAnswer 請 AI 讀挑到的段落、用顧客的語言回覆，再驗證原文片段；ok 是 false 表示 AI 失敗、答不出來或驗證沒過
func (t *turn) writeAnswer(text string, picked []knowledgeItem) (Reply, bool) {
	var ids []string
	for _, it := range picked {
		ids = append(ids, it.id)
	}
	r, usage := t.bot.Model.Respond(t.ctx, answerRequest(t.cfg.Merchant.Name, t.customer.Lang, text, picked))
	if r == nil {
		t.trace("answer", failure(usage), nil)
		return Reply{}, false
	}
	answerable, _ := r.String("answerable")
	source, _ := r.String("source")
	quote, _ := r.String("quote")
	reply, _ := r.String("reply")
	reply = strings.TrimSpace(reply)
	t.trace("answer", answerable, map[string]string{"source": source, "quote": quote, "reply": reply})
	if answerable != "yes" {
		return Reply{}, false
	}
	sourceText := ""
	if i := slices.Index(ids, source); i >= 0 {
		sourceText = picked[i].body
	}
	result := verifyAnswer(answerCheck{Source: sourceText, Quote: quote, Reply: reply, Message: text, Lang: t.customer.Lang})
	if result != "" {
		t.trace("verify", result, nil)
		t.verifyFailed = true
		return Reply{}, false
	}
	t.trace("verify", "ok", nil)
	return Reply{Key: "answer", Text: reply}, true
}

// AnswerRequest 是寫回覆送給 AI 的題目：sections 依序編成 K1、K2…，全文附在指示裡（量測工具 ctxprobe 用它）
func AnswerRequest(store, lang, text string, sections []merchant.Section) model.Request {
	return answerRequest(store, lang, text, knowledgeItems(sections))
}

// knowledgeItems 把知識段落依序編成 K1、K2…
func knowledgeItems(sections []merchant.Section) []knowledgeItem {
	var items []knowledgeItem
	for i, s := range sections {
		items = append(items, knowledgeItem{id: fmt.Sprintf("K%d", i+1), label: s.Title, body: s.Title + "\n" + s.Body})
	}
	return items
}

// answerRequest 是寫回覆的題目：store 是店名，lang 是顧客的語言代碼
func answerRequest(store, lang, text string, picked []knowledgeItem) model.Request {
	var ids, docs []string
	for _, it := range picked {
		ids = append(ids, it.id)
		docs = append(docs, "["+it.id+"] "+it.body)
	}
	langName := merchant.LanguageName(lang)
	fields := []model.Field{
		// 先抄證據再寫回覆（欄位依序產生）：逼 AI 先找到根據
		{Kind: model.KindText, Name: "quote", Description: "從資料逐字抄出回答根據的那一句，照資料原本的語言抄（不要改寫、不要翻成顧客的語言）；資料沒有答案就填空字串"},
		{Kind: model.KindChoice, Name: "source", Description: "上面那句出自哪一段資料", Options: ids},
		{Kind: model.KindChoice, Name: "answerable", Description: "資料能回答顧客的問題選 yes；資料沒有寫、要推測才能回答選 no", Options: []string{"yes", "no"}},
		{Kind: model.KindText, Name: "reply", Description: "給顧客的回覆，用" + langName + "，簡短親切；只用資料裡的資訊，數字、時間照資料寫"},
	}
	return model.Request{
		Name: "Answer", Text: text, Fields: fields,
		Instructions: "你是「" + store + "」的客服。只能根據以下資料回答顧客，用" + langName + "回覆。" +
			"資料沒有寫的事不要回答、不要推測，answerable 選 no。\n資料：\n" + strings.Join(docs, "\n\n"),
	}
}
