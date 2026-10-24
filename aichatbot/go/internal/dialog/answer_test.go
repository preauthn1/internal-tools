package dialog

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// knowledgeMerchant 是有知識段落和 FAQ 的商家
func knowledgeMerchant() merchant.Config {
	c := testMerchant()
	c.Knowledge = []merchant.Section{
		{Title: "營業時間", Body: "11:00–20:00 營業，週一公休。"},
		{Title: "菜單與價格", Body: "桂花烏龍 110 元\n抹茶拿鐵 150 元"},
	}
	c.FAQ = []merchant.FAQ{{Key: "card", Question: "可以刷卡嗎", AltQuestions: []string{"收信用卡嗎"}, Answer: "信用卡需消費滿 200 元。"}}
	return c
}

func pick(ask string, picks ...string) map[string]model.Value {
	v := map[string]model.Value{"ask": model.Text(ask)}
	for i := range 3 {
		p := "none"
		if i < len(picks) {
			p = picks[i]
		}
		v[[]string{"pick1", "pick2", "pick3"}[i]] = model.Text(p)
	}
	return v
}

func answer(source, quote, answerable, reply string) map[string]model.Value {
	return map[string]model.Value{"quote": model.Text(quote), "source": model.Text(source), "answerable": model.Text(answerable), "reply": model.Text(reply)}
}

func faqPick(faq, kind string) map[string]model.Value {
	return faqPickAsk("", faq, kind)
}

func faqPickAsk(ask, faq, kind string) map[string]model.Value {
	return map[string]model.Value{"ask": model.Text(ask), "faq": model.Text(faq), "kind": model.Text(kind)}
}

// cardCopied 是挑到刷卡的 FAQ、ask 抄的是那題的答案（#51）
var cardCopied = faqPickAsk("信用卡需消費滿 200 元。", "card", "question")

func steps(r Reply) string {
	var out []string
	for _, tr := range r.Traces {
		out = append(out, tr.Step+":"+tr.Result)
	}
	return strings.Join(out, ",")
}

// 挑到 FAQ、顧客用商家慣用語言：回店家寫好的答案原文，不呼叫寫回覆；選項附上答案內容
func TestGeneralFAQAnswerVerbatim(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("card", "question")},
	})
	r := e.say("能刷卡嗎")
	if r.Key != "answer" || r.Text != "信用卡需消費滿 200 元。" {
		t.Fatalf("got %+v", r)
	}
	reqs := e.model.Requests()
	if len(reqs) != 2 || reqs[1].Name != "FAQ" {
		t.Fatalf("requests：%+v", reqs)
	}
	if !strings.Contains(reqs[1].Instructions, "card：可以刷卡嗎/收信用卡嗎（答案：信用卡需消費滿 200 元。）") {
		t.Errorf("挑 FAQ 的選項：%s", reqs[1].Instructions)
	}
	if got := steps(r); got != "intent:general,faq:card" {
		t.Errorf("traces：%s", got)
	}
}

// 挑到 FAQ、顧客用其他語言：把答案翻成顧客的語言，譯文通過數字驗證就回譯文，不呼叫寫回覆
func TestGeneralFAQTranslated(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent":    {intent("general")},
		"FAQ":       {faqPick("card", "question")},
		"Translate": {{"translation": model.Text("Credit cards are accepted for orders of 200 or more.")}},
	})
	r := e.say("Can I pay by card?")
	if r.Key != "answer" || r.Text != "Credit cards are accepted for orders of 200 or more." {
		t.Fatalf("got %+v", r)
	}
	reqs := e.model.Requests()
	if len(reqs) != 3 || reqs[2].Name != "Translate" || reqs[2].Text != "信用卡需消費滿 200 元。" || !strings.Contains(reqs[2].Instructions, "English") {
		t.Fatalf("requests：%+v", reqs)
	}
	if got := steps(r); got != "intent:general,faq:card,translate:ok" {
		t.Errorf("traces：%s", got)
	}
}

// 譯文驗證沒過（多出原文沒有的數字）：改走挑知識段落、AI 寫回覆；挑到 FAQ 就是問商家（kind 填 chat 也一樣）
func TestGeneralFAQTranslationBlocked(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent":    {intent("general")},
		"FAQ":       {faqPick("card", "chat")},
		"Translate": {{"translation": model.Text("Credit cards need a minimum of 300.")}},
		"Pick":      {pick("card", "K2")},
		"Answer":    {answer("K2", "抹茶拿鐵 150 元", "yes", "A matcha latte is 150.")},
	})
	r := e.say("Can I pay by card?")
	if r.Key != "answer" || r.Text != "A matcha latte is 150." {
		t.Fatalf("got %+v", r)
	}
	if got := steps(r); got != "intent:general,faq:card,translate:number_not_in_quote:300,pick:K2,answer:yes,verify:ok" {
		t.Errorf("traces：%s", got)
	}
}

// manyFAQ 是 n 題 FAQ（f1…fn）；最後一題的其他問法是「幾點打烊」
func manyFAQ(n int) merchant.Config {
	c := knowledgeMerchant()
	c.FAQ = nil
	for i := 1; i <= n; i++ {
		f := merchant.FAQ{Key: fmt.Sprintf("f%d", i), Question: fmt.Sprintf("問題%d", i), Answer: fmt.Sprintf("答案%d", i)}
		if i == n {
			f.AltQuestions = []string{"幾點打烊"}
		}
		c.FAQ = append(c.FAQ, f)
	}
	return c
}

// FAQ 超過門檻題數（#65）：先用句子向量挑出最像的 N 題才交給 AI，最像的那題（排在最後）在裡面，挑到它就回它的答案
func TestPickFAQFiltersManyFAQ(t *testing.T) {
	n := merchant.FAQFilterAbove + 1
	last := fmt.Sprintf("f%d", n)
	e := newEnv(t, manyFAQ(n), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick(last, "question")},
	})
	// 最後一題（其他問法）最像；f26 以後的主問法有點像（一樣像的照原本順序取）；f1…f25 完全不像
	vectors := map[string][]float32{"你們開到幾點": {1, 0, 0}, "幾點打烊": {0.9, 0.1, 0}}
	for i := 26; i <= n; i++ {
		vectors[fmt.Sprintf("問題%d", i)] = []float32{0.5, 0.5, 0}
	}
	e.bot.Embedder = &fakeEmbedder{vectors: vectors}
	r := e.say("你們開到幾點")
	if r.Key != "answer" || r.Text != fmt.Sprintf("答案%d", n) {
		t.Fatalf("got %+v", r)
	}
	req := e.model.Requests()[1]
	var want []string
	for i := 26; i < 26+merchant.FAQFilterTop-1; i++ {
		want = append(want, fmt.Sprintf("f%d", i))
	}
	want = append(want, last, "price_calc", "unclear", "not_a_question", "none")
	if options := req.Fields[1].Options; !slices.Equal(options, want) {
		t.Errorf("挑 FAQ 的選項（%d 個）：%v", len(options), options)
	}
	if strings.Count(req.Instructions, "（答案：") != merchant.FAQFilterTop {
		t.Errorf("指示裡的 FAQ 題數不對：%s", req.Instructions)
	}
	if tr := r.Traces[1]; tr.Step != "faq" || tr.Detail["filtered"] != fmt.Sprintf("%d/%d", merchant.FAQFilterTop, n) {
		t.Errorf("trace：%+v", tr)
	}
}

// FAQ 超過門檻題數、但沒有 Embedder（算不出向量）：照舊送全部的題目
func TestPickFAQNoEmbedderSendsAll(t *testing.T) {
	n := merchant.FAQFilterAbove + 1
	e := newEnv(t, manyFAQ(n), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
	})
	e.bot.Log = slog.New(slog.DiscardHandler)
	e.say("你們開到幾點")
	if got := len(e.model.Requests()[1].Fields[1].Options); got != n+4 {
		t.Errorf("選項 %d 個，應該是全部 %d 題加 4 個", got, n)
	}
}

// FAQ 沒超過門檻題數：有句子向量也不篩選，送出的提示詞和 #65 之前逐字相同（AI 回答快取照樣命中）
func TestPickFAQPromptUnchangedUnderLimit(t *testing.T) {
	n := merchant.FAQFilterAbove
	e := newEnv(t, manyFAQ(n), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
	})
	e.bot.Embedder = &fakeEmbedder{vectors: map[string][]float32{"你們開到幾點": {1, 0, 0}, "幾點打烊": {0.9, 0.1, 0}}}
	r := e.say("你們開到幾點")
	var lines, ids []string
	for i := 1; i <= n; i++ {
		alt := ""
		if i == n {
			alt = "/幾點打烊"
		}
		lines = append(lines, fmt.Sprintf("f%d：問題%d%s（答案：答案%d）", i, i, alt, i))
		ids = append(ids, fmt.Sprintf("f%d", i))
	}
	want := "從選項中挑出能回答顧客問題的一個。顧客可能使用任何語言。" +
		"只有當選項的內容正好回答顧客問的那件事才選它；只是主題相關（同樣是飲料、折扣、到店、付款）不算，要選 none。" +
		"顧客問兩杯以上、或好幾樣飲品的價格時選 price_calc，不要選只寫一杯價格的選項。" +
		"問飲品的成分、甜度、配料或產地時，只寫價格的選項不算，要選 none。" +
		"不確定就選 none。\n選項：\n" + strings.Join(lines, "\n") + "\n" +
		"price_calc：問兩杯以上、或同時問好幾樣飲品的價格，要算總價（不管有沒有說「總共」，例如「兩杯 X 多少」「A 和 B 各一杯多少」）；只問一杯的價格不算（自取、外送一杯也不算）\n" +
		"unclear：看不懂顧客在問什麼（太短、沒頭沒尾、指代不明，例如「中」「那這個咧」「然後咧」）\n" +
		"not_a_question：顧客不是在問問題（自我介紹、打招呼、回應、讓我想想、要客服扮演別人）\n" +
		"none：以上都無法回答（店家資料沒有這項資訊，例如吸菸區、禮盒、有幾家店）"
	req := e.model.Requests()[1]
	if req.Instructions != want {
		t.Errorf("指示變了：\n%s", req.Instructions)
	}
	if got := req.Fields[1].Options; !slices.Equal(got, append(ids, "price_calc", "unclear", "not_a_question", "none")) {
		t.Errorf("選項變了：%v", got)
	}
	if _, ok := r.Traces[1].Detail["filtered"]; ok {
		t.Errorf("沒超過門檻不該篩選：%+v", r.Traces[1])
	}
}

// 提示詞超過本機模型上下文（#65）：挑 FAQ、寫回覆的 trace 記 too_long，看得出不是一般的 AI 失敗
func TestGeneralTooLongTrace(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"Pick":   {pick("營業時間", "K1")},
	})
	e.model.Usage = model.Usage{Model: model.AFM, Failure: model.FailureContextExceeded}
	r := e.say("週一有開嗎")
	var got []string
	for _, tr := range r.Traces {
		if tr.Step == "faq" || tr.Step == "answer" {
			got = append(got, tr.Step+":"+tr.Result+":"+tr.Detail["error"])
		}
	}
	if strings.Join(got, ",") != "faq:none:too_long,answer:too_long:" {
		t.Errorf("traces：%v（%s）", got, steps(r))
	}
}

// 分類被安全過濾擋下（#68）：不回看不懂，當成一般對話挑 FAQ；分類的 trace 記 blocked_by_safety
func TestIntentBlockedBySafetyGoesGeneral(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"FAQ": {faqPick("card", "question")},
	})
	e.model.Usage = model.Usage{Model: model.AFM, Failure: model.FailureBlockedBySafety}
	r := e.say("能刷卡嗎")
	if r.Key != "answer" || r.Text != "信用卡需消費滿 200 元。" {
		t.Fatalf("got %+v（%s）", r, steps(r))
	}
	if tr := r.Traces[0]; tr.Step != "intent" || tr.Result != "general" || tr.Detail["error"] != model.FailureBlockedBySafety {
		t.Errorf("分類的 trace：%+v", tr)
	}
}

// 分類和一般對話每一步都被安全過濾擋下（#68）：照「找不到答案」轉人工，每一步的 trace 都記 blocked_by_safety
func TestAllBlockedBySafetyHandsOff(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), nil)
	e.model.Usage = model.Usage{Model: model.AFM, Failure: model.FailureBlockedBySafety}
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	r := e.say("指定店長要加錢嗎")
	if r.Key != "handoff" {
		t.Fatalf("got %+v（%s）", r, steps(r))
	}
	var got []string
	for _, tr := range r.Traces {
		got = append(got, tr.Step+":"+tr.Detail["error"])
	}
	if strings.Join(got, ",") != "intent:blocked_by_safety,faq:blocked_by_safety,pick:blocked_by_safety,ask:blocked_by_safety" {
		t.Errorf("traces：%v", got)
	}
}

// 挑到的 FAQ 翻譯時被安全過濾擋下（#68）：translate 的 trace 記 blocked_by_safety
func TestTranslateBlockedBySafetyTrace(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("card", "question")},
	})
	e.model.Usage = model.Usage{Model: model.AFM, Failure: model.FailureBlockedBySafety}
	r := e.say("Can I pay by card?")
	i := slices.IndexFunc(r.Traces, func(tr Trace) bool { return tr.Step == "translate" })
	if i < 0 || r.Traces[i].Result != "failed" || r.Traces[i].Detail["error"] != model.FailureBlockedBySafety {
		t.Errorf("traces：%+v", r.Traces)
	}
}

// F33：預約問姓名時回「好」（寒暄）或問了答不出來的事：再問一次姓名，不轉人工、預約不動
func TestGeneralDuringBookingReasks(t *testing.T) {
	for _, c := range []struct {
		p   map[string]model.Value
		ask []map[string]model.Value
	}{
		{faqPick("not_a_question", "chat"), nil},
		{faqPick("none", "question"), nil},
		{cardCopied, []map[string]model.Value{{"kind": model.Text("small_talk")}}}, // #51
		{faqPick("price_calc", "question"), nil},                                   // #52
	} {
		p := c.p
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent":  {intent("general")},
			"Booking": {{}},
			"FAQ":     {p},
			"Pick":    {pick("")},
			"Ask":     c.ask,
		})
		e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
		booking := map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}
		if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "name", Booking: booking}); err != nil {
			t.Fatal(err)
		}
		r := e.say("好")
		c := e.customer()
		if r.Key != "ask_name" || c.Human || len(e.pending()) != 0 || c.AskingField != "name" || c.Booking["date"] != "2026-10-21" {
			t.Errorf("%v：%+v %+v %q", p["faq"], r, c, e.pending())
		}
		if last := r.Traces[len(r.Traces)-1]; last.Step != "reask" || last.Result != "name" {
			t.Errorf("traces：%s", steps(r))
		}
	}
	// 商家沒有知識段落也沒有 FAQ：一樣再問
	e := midBooking(t, "zh", "name", "general", values{})
	if r := e.say("好"); r.Key != "ask_name" {
		t.Errorf("沒有知識：%+v", r)
	}
}

// 挑 FAQ 的欄位依序產生：寒暄還是問商家放在 faq 後面。放在最前面時，真模型把 F33 的「好」挑成營業時間的 FAQ（#45）
func TestGeneralFAQFieldOrder(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("not_a_question", "chat")},
	})
	e.say("好")
	reqs := e.model.Requests()
	if len(reqs) != 2 || reqs[1].Name != "FAQ" {
		t.Fatalf("requests：%+v", reqs)
	}
	var names []string
	for _, f := range reqs[1].Fields {
		names = append(names, f.Name)
	}
	if got := strings.Join(names, ","); got != "ask,faq,kind" {
		t.Errorf("挑 FAQ 的欄位：%s", got)
	}
}

// #51：簡短應答（知道了、OK 了解）挑到 FAQ 時，真模型的 ask 抄的是選項的答案、不是顧客的訊息：不回 FAQ，當成寒暄回固定回覆，不挑知識段落
func TestGeneralFAQAskCopiedFromAnswer(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {cardCopied},
		"Ask":    {{"kind": model.Text("small_talk")}},
	})
	r := e.say("知道了")
	if r.Key != "general" || r.Text != "收到～" || len(e.model.Requests()) != 3 {
		t.Errorf("got %+v %+v", r, e.model.Requests())
	}
	if got := steps(r); got != "intent:general,faq:card,ask:chat,chat:general" {
		t.Errorf("traces：%s", got)
	}
	// 確認時真模型把「知道了」「OK 了解」判成在問店家：提示詞要講明回應、表示收到不是在問。
	// 只列例句時「瞭解」還是被判成在問（#53），回應要講出是哪一類話（不列題庫裡的例句）
	ask := e.model.Requests()[2]
	if ask.Name != "Ask" || !strings.Contains(ask.Instructions, "表示收到、知道了的回應不是在問") || !strings.Contains(ask.Fields[0].Description, "回應（表示收到、明白或同意的話，例如收到囉、我懂、曉得了、got it）") {
		t.Errorf("確認的提示詞：%+v", ask)
	}
}

// #51：真問句也會把答案抄進 ask（營業時間、Are you open on Monday?）：確認是在問店家，或確認時 AI 失敗，都照舊回 FAQ 答案
func TestGeneralFAQAskCopiedButQuestion(t *testing.T) {
	for _, check := range [][]map[string]model.Value{{{"kind": model.Text("asks_about_store")}}, nil} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {cardCopied},
			"Ask":    check,
		})
		r := e.say("刷卡")
		if r.Key != "answer" || r.Text != "信用卡需消費滿 200 元。" {
			t.Errorf("%v：%+v", check, r)
		}
		if got := steps(r); got != "intent:general,faq:card,ask:question" {
			t.Errorf("%v traces：%s", check, got)
		}
	}
}

// #51：ask 是從顧客的訊息抄的（就算答案裡也有這幾個字）：照舊回 FAQ 答案
func TestGeneralFAQAskFromMessage(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPickAsk("信用卡", "card", "question")},
	})
	if r := e.say("信用卡可以嗎"); r.Key != "answer" || r.Text != "信用卡需消費滿 200 元。" {
		t.Errorf("got %+v", r)
	}
}

// 挑不到 FAQ：挑知識段落（清單只列知識段落）→ 寫回覆 → 驗證通過，回 AI 寫的回覆；寫回覆只給挑到的段落；每步都記進 traces
func TestGeneralAnswerPasses(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("", "K1")},
		"Answer": {answer("K1", "11:00–20:00 營業", "yes", "我們 11 點開門喔～")},
	})
	r := e.say("你們幾點開")
	if r.Key != "answer" || r.Text != "我們 11 點開門喔～" {
		t.Fatalf("got %+v", r)
	}
	reqs := e.model.Requests()
	if len(reqs) != 4 || reqs[2].Name != "Pick" || reqs[3].Name != "Answer" {
		t.Fatalf("requests：%+v", reqs)
	}
	if !strings.Contains(reqs[2].Instructions, "K1 營業時間") || strings.Contains(reqs[2].Instructions, "可以刷卡嗎") {
		t.Errorf("挑知識的清單：%s", reqs[2].Instructions)
	}
	if !strings.Contains(reqs[3].Instructions, "11:00–20:00 營業，週一公休。") || strings.Contains(reqs[3].Instructions, "桂花烏龍") {
		t.Errorf("寫回覆的資料：%s", reqs[3].Instructions)
	}
	if got := steps(r); got != "intent:general,faq:none,pick:K1,answer:yes,verify:ok" {
		t.Errorf("traces：%s", got)
	}
}

// 驗證沒過（回覆多出原文沒有的數字）：問商家的轉人工；traces 記下擋下的原因
func TestGeneralBlockedByVerify(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("", "K2", "K2")},
		"Answer": {answer("K2", "桂花烏龍 110 元", "yes", "桂花烏龍一杯 120 元")},
	})
	r := e.say("桂花烏龍多少錢")
	if r.Key != "handoff" || !e.customer().Human {
		t.Errorf("問商家：%+v", r)
	}
	if got := steps(r); got != "intent:general,faq:none,pick:K2,answer:yes,verify:number_not_in_quote:120" {
		t.Errorf("traces（重複的段落去掉）：%s", got)
	}
}

// 寒暄（挑 FAQ 時選 not_a_question 或 kind 是 chat）：回固定回覆，不挑知識段落；看不懂回看不懂；
// 商家沒有知識段落也沒有 FAQ：不呼叫挑 FAQ
func TestGeneralSkipsAI(t *testing.T) {
	for _, p := range []map[string]model.Value{faqPick("not_a_question", "question"), faqPick("none", "chat")} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {p},
		})
		if r := e.say("謝謝你"); r.Key != "general" || r.Text != "不客氣～有需要隨時找我喔！" || len(e.model.Requests()) != 2 || e.customer().Human {
			t.Errorf("寒暄：%+v %+v", r, e.model.Requests())
		}
	}
	for _, kind := range []string{"question", "chat"} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {faqPick("unclear", kind)},
		})
		if r := e.say("那個呢"); r.Key != "unclear" || len(e.model.Requests()) != 2 || e.customer().Human {
			t.Errorf("看不懂（kind %s）：%+v", kind, r)
		}
	}
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general")}})
	if r := e.say("謝謝你"); r.Key != "general" || len(e.model.Requests()) != 1 {
		t.Errorf("沒有知識：%+v %+v", r, e.model.Requests())
	}
}

// 只有符號、標點、顏文字（沒有文字和數字）是應答：回「收到」，不挑 FAQ、不呼叫 AI（#63，「＝＝」曾被挑成營業時間）
func TestSymbolsOnlyIsAcknowledgement(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), nil)
	for _, text := range []string{"＝＝", "…", "＝ ＝", "^^"} {
		if r := e.say(text); r.Key != "general" || r.Text != "收到～" {
			t.Errorf("%s：%+v", text, r)
		}
	}
	if reqs := e.model.Requests(); len(reqs) != 0 {
		t.Errorf("不該呼叫模型：%+v", reqs)
	}
}

// 預約中只回符號：再問一次正在問的欄位，不呼叫 AI（#63）
func TestSymbolsOnlyDuringBookingReasks(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), nil)
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "people", Booking: map[string]string{"date": "2026-10-21", "time": "14:00"}}); err != nil {
		t.Fatal(err)
	}
	if r := e.say("…"); r.Key != "ask_people" || e.customer().AskingField != "people" || len(e.model.Requests()) != 0 {
		t.Errorf("got %+v %+v", r, e.model.Requests())
	}
}

// #45：挑 FAQ 把閒聊判成問商家、挑不到知識時，轉人工前再確認一次：閒聊回寒暄、沒講出是哪件事回看不懂，不轉人工、不記待處理事項
func TestGeneralChatMistakenForQuestion(t *testing.T) {
	for _, c := range []struct{ text, check, key, kind string }{
		{"我想一下", "small_talk", "general", "chat,chat:general"},
		{"我朋友叫小美，她推薦你們", "small_talk", "general", "chat,chat:general"},
		{"還有嗎", "too_vague", "unclear", "unclear"},
	} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {faqPick("none", "question")},
			"Pick":   {pick(c.text)},
			"Ask":    {{"kind": model.Text(c.check)}},
		})
		e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
		r := e.say(c.text)
		if r.Key != c.key || e.customer().Human || len(e.pending()) != 0 {
			t.Errorf("%s：%+v %q", c.text, r, e.pending())
		}
		if got := steps(r); !strings.HasSuffix(got, "faq:none,pick:none,ask:"+c.kind) {
			t.Errorf("%s traces：%s", c.text, got)
		}
	}
}

// 問店家、資料裡沒有（ADR 0002）：確認是在問店家，或確認時 AI 失敗，都照舊轉人工，待處理事項記下原因
func TestGeneralStoreQuestionWithoutDataHandsOff(t *testing.T) {
	for _, check := range [][]map[string]model.Value{{{"kind": model.Text("asks_about_store")}}, nil} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {faqPick("none", "question")},
			"Pick":   {pick("有停車位嗎")},
			"Ask":    check,
		})
		e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
		r := e.say("有停車位嗎")
		if r.Key != "handoff" || !e.customer().Human || len(e.pending()) != 1 || e.pending()[0] != "handoff：🙋 轉人工：找不到相關資料：有停車位嗎" {
			t.Errorf("%v：%+v %q", check, r, e.pending())
		}
		if got := steps(r); got != "intent:general,faq:none,pick:none,ask:question" {
			t.Errorf("%v traces：%s", check, got)
		}
	}
}

// #56：問店裡的人被 AI 分成問身分時改成一般對話，照問店家處理：資料沒有就轉人工，不回客服的自我介紹
func TestThirdPersonWhoWithoutDataHandsOff(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("identity")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("老闆叫什麼")},
		"Ask":    {{"kind": model.Text("asks_about_store")}},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	r := e.say("你老闆叫什麼？")
	if r.Key != "handoff" || strings.Contains(r.Text, "AI 客服") {
		t.Errorf("got %+v", r)
	}
	if got := steps(r); got != "intent:general,faq:none,pick:none,ask:question" {
		t.Errorf("traces：%s", got)
	}
}

// #52：挑 FAQ 時選了算價（程式認不得的說法，例如英文）：照設計直接轉人工（#44），不挑知識段落、不讓 AI 寫出單價
func TestGeneralPriceCalcHandsOff(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPickAsk("How much for two matcha lattes", "price_calc", "question")},
		"Pick":   {pick("How much for two matcha lattes", "K2")},
		"Answer": {answer("K2", "抹茶拿鐵 150 元", "yes", "A matcha latte is 150.")},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	r := e.say("How much for two matcha lattes?")
	if r.Key != "handoff" || !e.customer().Human || len(e.pending()) != 1 || e.pending()[0] != "handoff：🙋 轉人工：需要計算總價：How much for two matcha lattes?" {
		t.Errorf("%+v %q", r, e.pending())
	}
	if got := steps(r); got != "intent:general,faq:price_calc" {
		t.Errorf("traces：%s", got)
	}
}

// #52：AI 把「A 和 B 各一杯多少」挑成單價的 FAQ 時，程式認得多品項的價格問法，照樣轉人工、不呼叫挑 FAQ
func TestGeneralMultiItemPriceRuleHandsOff(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("card", "question")},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	r := e.say("紅玉和抹茶拿鐵各一杯多少")
	if r.Key != "handoff" || len(e.pending()) != 1 || e.pending()[0] != "handoff：🙋 轉人工：需要計算總價：紅玉和抹茶拿鐵各一杯多少" {
		t.Errorf("%+v %q", r, e.pending())
	}
	if got := steps(r); got != "intent:general,faq:price_calc" || r.Traces[1].Detail["rule"] != "multi_item" {
		t.Errorf("traces：%s %v", got, r.Traces)
	}
}

// #52、#64：多件（數量大於 1）、多品項的價格問法才算；只問一件、沒問價格、不是商品數量的不算
func TestMultiItemPrice(t *testing.T) {
	for text, want := range map[string]bool{
		"外帶兩杯抹茶拿鐵多少":    true,
		"3杯紅玉多少錢":       true,
		"買十二杯桂花烏龍要多少錢":  true,
		"紅玉跟抹茶各一杯多少":    true,
		"兩杯桂花烏龍自取有打折嗎":  true,
		"抹茶拿鐵多少":        false,
		"紅玉一杯多少":        false,
		"一杯多少":          false,
		"兩杯抹茶拿鐵，名字是王小明": false,
		"兩點開門嗎":         false,
		"包場低消多少":        false,
		"兩個人多少錢":        false,
		"抹茶拿鐵第二杯多少":     false,
		"訂兩杯要等多少時間":     false,
		"滿五杯有折扣嗎":       false,
		// #64：飲料以外的量詞
		"兩本筆記本多少":     true,
		"三件 T 恤多少錢":   true,
		"書和筆各一本多少":    true,
		"一件多少":        false,
		"第二件多少":       false,
		"兩個蛋糕多少錢":     true,
		"三個月多少錢":      false,
		"兩個小時多少錢":     false,
		"兩個大人一個小孩多少錢": false,
		"兩個小朋友多少錢":    false,
		"兩個晚上多少錢":     false,
	} {
		if got := multiItemPrice(text); got != want {
			t.Errorf("%s：got %v", text, got)
		}
	}
}

// #48：確認時判成和店家無關（講得清楚的常識、時事）：回無關的固定回覆、計入連續無關次數，不轉人工、不記待處理事項
func TestGeneralOffTopicCheckCountsStreak(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("台灣最高的山是哪座")},
		"Ask":    {{"kind": model.Text("off_topic")}},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	r := e.say("台灣最高的山是哪座")
	c := e.customer()
	if r.Key != "off_topic" || c.Human || len(e.pending()) != 0 || c.OffTopicStreak != 1 {
		t.Errorf("%+v %+v %q", r, c, e.pending())
	}
	if got := steps(r); got != "intent:general,faq:none,pick:none,ask:off_topic" {
		t.Errorf("traces：%s", got)
	}
}

// #48：前兩句分類成無關，第三句分類成一般對話、確認時判成無關：接上前兩句的次數，第三句轉人工
func TestGeneralOffTopicCheckStreakHandsOff(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("off_topic"), intent("off_topic"), intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("地球離太陽多遠")},
		"Ask":    {{"kind": model.Text("off_topic")}},
	})
	var keys []string
	for _, text := range []string{"幫我寫一首詩", "推薦一部電影", "地球離太陽多遠"} {
		keys = append(keys, e.say(text).Key)
	}
	if got := strings.Join(keys, ","); got != "off_topic,off_topic,handoff" || !e.customer().Human {
		t.Errorf("got %s", got)
	}
	if got := strings.Join(e.pending(), "\n"); got != "handoff：🙋 轉人工：連續 3 次無關或聲稱的訊息" {
		t.Errorf("待處理：%s", got)
	}
}

// #48：預約問欄位時確認判成無關：和意圖直接分成無關一樣回無關的固定回覆、計次，不再問欄位、預約不動
func TestGeneralOffTopicCheckDuringBooking(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent":  {intent("general")},
		"Booking": {{}},
		"FAQ":     {faqPick("none", "question")},
		"Pick":    {pick("地球離太陽多遠")},
		"Ask":     {{"kind": model.Text("off_topic")}},
	})
	booking := map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "name", Booking: booking, OffTopicStreak: 1}); err != nil {
		t.Fatal(err)
	}
	r := e.say("地球離太陽多遠")
	c := e.customer()
	if r.Key != "off_topic" || c.OffTopicStreak != 2 || c.AskingField != "name" || c.Booking["date"] != "2026-10-21" {
		t.Errorf("%+v %+v", r, c)
	}
}

// #48：確認時判成無關以外的結果（閒聊）：連續無關次數照樣歸零
func TestGeneralChatCheckResetsStreak(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("off_topic"), intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("我想一下")},
		"Ask":    {{"kind": model.Text("small_talk")}},
	})
	e.say("幫我寫一首詩")
	if r := e.say("我想一下"); r.Key != "general" || e.customer().OffTopicStreak != 0 {
		t.Errorf("%+v %+v", r, e.customer())
	}
}

// 沒有進行中的預約時說取消、沒有名字時被分成自我介紹：都當成一般對話
func TestCancelAndSelfIntroFallBackToGeneral(t *testing.T) {
	for _, code := range []string{"cancel", "self_intro"} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent(code)},
			"FAQ":    {faqPick("none", "question")},
			"Pick":   {pick("", "K1")},
			"Answer": {answer("K1", "週一公休", "yes", "我們週一公休喔")},
		})
		if r := e.say("那週一呢"); r.Key != "answer" {
			t.Errorf("%s：%+v", code, r)
		}
	}
}

// AI 說答不出來（answerable: no）而且是問商家：照「找不到答案時轉人工」的設定，開著轉人工、關著回看不懂；traces 記下 answer:no
func TestGeneralUnanswerableQuestion(t *testing.T) {
	for _, on := range []bool{true, false} {
		e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
			"Intent": {intent("general")},
			"FAQ":    {faqPick("none", "question")},
			"Pick":   {pick("", "K1")},
			"Answer": {answer("K1", "", "no", "")},
		})
		e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = on })
		r := e.say("停車場在哪")
		if on && (r.Key != "handoff" || !e.customer().Human || e.pending()[0] != "handoff：🙋 轉人工：資料裡找不到答案：停車場在哪") {
			t.Errorf("開著：%+v %q", r, e.pending())
		}
		if !on && (r.Key != "unclear" || e.customer().Human || len(e.pending()) != 0) {
			t.Errorf("關著：%+v %q", r, e.pending())
		}
		if last := r.Traces[len(r.Traces)-1]; last.Step != "answer" || last.Result != "no" {
			t.Errorf("traces：%+v", r.Traces)
		}
	}
}
