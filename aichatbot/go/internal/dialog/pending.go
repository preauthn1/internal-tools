package dialog

import (
	"context"
	"slices"

	"aichatbot/internal/embed"
	"aichatbot/internal/merchant"
	"aichatbot/internal/store"
)

// MARK: 待補問題（CONTEXT.md）：機器人答不出來的訊息，相同問題合併計次，管理員在後台轉成 FAQ、知識段落、意圖案例或忽略

// unanswered 記下這則訊息答不出來（處理完才寫進待補問題，見 recordPending）。source 是預設的來源，
// 這則訊息先前有原文片段驗證沒過時用 verify_failed
func (t *turn) unanswered(source string) {
	if t.missed != "" {
		return
	}
	if t.verifyFailed {
		source = "verify_failed"
	}
	t.missed = source
}

// recordPending 把答不出來的訊息寫進待補問題（回覆照舊）；有 AI 失敗的不記（不是真的答不出來）。
// 這則訊息轉人工了就把轉人工的事項記在這次出現上（員工回覆後預填答案）。寫不進去只記錯誤
func (t *turn) recordPending() {
	if t.missed == "" || slices.ContainsFunc(t.traces, aiFailed) {
		return
	}
	var ticket int64
	if t.customer.Human && slices.Contains(t.opened, t.customer.HandoffTicket) {
		ticket = t.customer.HandoffTicket
	}
	hit := store.PendingHit{Time: t.now, Platform: t.msg.Platform, UserID: t.msg.UserID, ChatID: t.msg.ChatID,
		Text: t.msg.Text, Lang: t.customer.Lang, Source: t.missed, TicketID: ticket}
	vector := t.vector(hit.Text)
	id, merged, err := t.bot.DB.RecordPending(t.ctx, hit, vector, samePending(hit.Text, vector))
	if err != nil {
		t.bot.log().Error("寫入待補問題失敗", "error", err.Error())
		return
	}
	// 新的待補問題附上翻成商家慣用語言的版本：在背景翻，不拖慢回覆；翻不出來就留空，後台詳細頁再補
	if from, to := hit.Lang, t.cfg.Merchant.PreferredLanguage; !merged && t.bot.Async != nil && from != "" && from != to {
		ctx := context.WithoutCancel(t.ctx)
		t.bot.Async(func() { t.bot.TranslatePending(ctx, id, from, to, hit.Text) })
	}
}

// TranslatePending 把待補問題 id 的原話 text 從 from 翻成 to 存起來；翻不出來只記錯誤
func (b *Bot) TranslatePending(ctx context.Context, id int64, from, to, text string) {
	translated, err := merchant.TranslateCustomer(ctx, b.Model, from, to, text)
	if err == nil {
		err = b.DB.SetPendingTranslation(ctx, id, translated)
	}
	if err != nil {
		b.log().Warn("待補問題翻譯失敗", "id", id, "error", err.Error())
	}
}

// pendingMergeThreshold：兩句的語意向量 cosine 到這裡才算同一個問題。這裡只判斷是不是「同一句問題」，
// 門檻設高、寧可不併（Swift 版的經驗：主題相近的不同問題 cosine 也常在 0.85 上下）
const pendingMergeThreshold = 0.92

// samePending 判斷已經有的待補問題（原話、向量）和這次的訊息是不是同一個問題：
// 兩邊都有向量時看 cosine，否則正規化後（轉小寫、只留字母和數字）文字相同才算
func samePending(text string, vector []float32) func(string, []float32) bool {
	return func(other string, v []float32) bool {
		if vector != nil && v != nil {
			return embed.Cosine(vector, v) >= pendingMergeThreshold
		}
		return merchant.Normalize(other) == merchant.Normalize(text)
	}
}

// aiFailed：這一步 AI 沒回答。翻譯失敗（translate:failed）也算：AI 沒回答和譯文格式不對分不出來，都不是答不出來
func aiFailed(tr Trace) bool {
	return tr.Result == "ai_failed" || tr.Detail["error"] == "ai_failed" || tr.Step == "translate" && tr.Result == "failed"
}
