package dialog

import (
	"context"
	"testing"

	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// pendingQuestions 是待處理的待補問題，每筆都讀出每次出現
func (e *env) pendingQuestions() []store.PendingQuestion {
	e.t.Helper()
	list, err := e.db.PendingQuestions(context.Background(), "open")
	if err != nil {
		e.t.Fatal(err)
	}
	for i, q := range list {
		if list[i], err = e.db.PendingQuestion(context.Background(), q.ID); err != nil {
			e.t.Fatal(err)
		}
	}
	return list
}

// AI 說答不出來（answerable: no）：記成待補問題（來源 unanswerable），轉人工開的事項記在那次出現上；回覆照舊
func TestUnanswerableBecomesPendingQuestion(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("", "K1")},
		"Answer": {answer("K1", "", "no", "")},
	})
	r := e.say("停車場在哪")
	if r.Key != "handoff" {
		t.Fatalf("回覆：%+v", r)
	}
	qs := e.pendingQuestions()
	if len(qs) != 1 || qs[0].Text != "停車場在哪" || qs[0].Lang != "zh" || qs[0].Count != 1 {
		t.Fatalf("待補問題：%+v", qs)
	}
	if h := qs[0].Hits[0]; h.Source != "unanswerable" || h.UserID != "U1" || h.ChatID != "C1" || h.TicketID != r.Tickets[0] || !h.Time.Equal(wednesday) {
		t.Errorf("出現：%+v（事項 %v）", h, r.Tickets)
	}
}

// 分類為看不懂（或挑 FAQ 時看不懂）：來源是 unclear；AI 失敗被當成看不懂的、寒暄的都不記
func TestUnclearBecomesPendingQuestion(t *testing.T) {
	for _, c := range []struct {
		name   string
		script map[string][]map[string]model.Value
		want   int
	}{
		{"分類為看不懂", map[string][]map[string]model.Value{"Intent": {intent("unclear")}}, 1},
		{"挑 FAQ 時看不懂", map[string][]map[string]model.Value{"Intent": {intent("general")}, "FAQ": {faqPick("unclear", "question")}}, 1},
		{"分類 AI 失敗", nil, 0},
		{"寒暄", map[string][]map[string]model.Value{"Intent": {intent("general")}, "FAQ": {faqPick("not_a_question", "chat")}}, 0},
	} {
		e := newEnv(t, knowledgeMerchant(), c.script)
		if r := e.say("那個呢"); r.Key != "unclear" && r.Key != "general" {
			t.Errorf("%s：回覆 %+v", c.name, r)
		}
		qs := e.pendingQuestions()
		if len(qs) != c.want || c.want == 1 && qs[0].Hits[0].Source != "unclear" {
			t.Errorf("%s：待補問題 %+v", c.name, qs)
		}
	}
}

// 合併計次：有 Embedder 時語意向量夠像（cosine ≥ 門檻）就是相同的問題；沒有 Embedder 時正規化後文字相同才合併
func TestPendingQuestionsMergeSimilar(t *testing.T) {
	unclear := func() map[string][]map[string]model.Value {
		return map[string][]map[string]model.Value{"Intent": {intent("unclear"), intent("unclear"), intent("unclear")}}
	}
	counts := func(qs []store.PendingQuestion) map[string]int {
		out := map[string]int{}
		for _, q := range qs {
			out[q.Text] = q.Count
		}
		return out
	}
	e := newEnv(t, knowledgeMerchant(), unclear())
	e.bot.Embedder = &fakeEmbedder{vectors: map[string][]float32{
		"停車場在哪": {1, 0, 0}, "哪裡可以停車": {0.95, 0.31, 0}, "可以帶狗嗎": {0.6, 0.8, 0},
	}}
	for _, s := range []string{"停車場在哪", "哪裡可以停車", "可以帶狗嗎"} {
		e.say(s)
	}
	if got := counts(e.pendingQuestions()); len(got) != 2 || got["停車場在哪"] != 2 || got["可以帶狗嗎"] != 1 {
		t.Errorf("有 Embedder：%v", got)
	}

	e = newEnv(t, knowledgeMerchant(), unclear())
	for _, s := range []string{"停車場在哪", "停車場 在哪？", "停車場在哪裡"} {
		e.say(s)
	}
	if got := counts(e.pendingQuestions()); len(got) != 2 || got["停車場在哪"] != 2 || got["停車場在哪裡"] != 1 {
		t.Errorf("沒有 Embedder：%v", got)
	}
}

// 顧客用的不是商家慣用語言：新的待補問題在背景（Bot.Async）用翻譯顧客訊息的提示詞翻成慣用語言；合併進去的不再翻；
// 沒設 Async 就不翻（後台詳細頁再補）
func TestPendingQuestionTranslated(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent":            {intent("unclear"), intent("unclear")},
		"TranslateCustomer": {{"translation": model.Text("哪裡可以停車？")}},
	})
	var later []func()
	e.bot.Async = func(f func()) { later = append(later, f) }
	e.say("Where can I park?")
	e.say("Where can I park")
	if len(later) != 1 {
		t.Fatalf("背景工作 %d 件", len(later))
	}
	later[0]()
	if qs := e.pendingQuestions(); len(qs) != 1 || qs[0].Lang != "en" || qs[0].Translated != "哪裡可以停車？" {
		t.Errorf("待補問題：%+v", qs)
	}
	if reqs := e.model.Requests(); reqs[len(reqs)-1].Name != "TranslateCustomer" || reqs[len(reqs)-1].Text != "Where can I park?" {
		t.Errorf("翻譯：%+v", reqs[len(reqs)-1])
	}

	e = newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{"Intent": {intent("unclear")}})
	e.say("Where can I park?")
	if reqs := e.model.Requests(); len(reqs) != 1 {
		t.Errorf("沒設 Async 也翻譯了：%+v", reqs)
	}
}

// 轉人工後員工在群組回覆那張事項（或對應的緊急通知）：員工的原話預填成待補問題的答案，回覆多次就接起來
func TestStaffReplyFillsPendingAnswer(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("", "K1")},
		"Answer": {answer("K1", "", "no", "")},
	})
	r := e.say("停車場在哪")
	ctx := context.Background()
	if _, err := e.bot.StaffReply(ctx, r.Tickets[0], "對面有付費停車場"); err != nil {
		t.Fatal(err)
	}
	urgent, err := e.db.AddTicket(ctx, store.Ticket{Kind: "urgent", Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "🔥", RelatedTo: r.Tickets[0]}, wednesday)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := e.bot.StaffReply(ctx, urgent, "一小時 40 元"); err != nil {
		t.Fatal(err)
	}
	if qs := e.pendingQuestions(); len(qs) != 1 || qs[0].StaffAnswer != "對面有付費停車場\n一小時 40 元" {
		t.Errorf("待補問題：%+v", qs)
	}
}

// 挑到 FAQ 但翻譯 AI 失敗（不是答不出來），之後挑不到知識：不記
func TestTranslateFailureNotPending(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("card", "question")},
		"Pick":   {pick("")},
	})
	e.say("Can I pay by card?")
	if qs := e.pendingQuestions(); len(qs) != 0 {
		t.Errorf("待補問題：%+v", qs)
	}
}

// 原文片段驗證沒過：來源是 verify_failed
func TestVerifyFailedBecomesPendingQuestion(t *testing.T) {
	e := newEnv(t, knowledgeMerchant(), map[string][]map[string]model.Value{
		"Intent": {intent("general")},
		"FAQ":    {faqPick("none", "question")},
		"Pick":   {pick("", "K2")},
		"Answer": {answer("K2", "桂花烏龍 110 元", "yes", "桂花烏龍一杯 120 元")},
	})
	e.say("桂花烏龍多少錢")
	if qs := e.pendingQuestions(); len(qs) != 1 || qs[0].Hits[0].Source != "verify_failed" {
		t.Errorf("待補問題：%+v", qs)
	}
}
