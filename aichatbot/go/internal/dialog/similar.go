package dialog

import (
	"context"
	"sort"

	"aichatbot/internal/embed"
	"aichatbot/internal/merchant"
)

// MARK: 意圖案例庫：找和顧客訊息最像的幾句給 AI 參考（e5 向量，#23 的結論；不限文字系統，中文案例也能給日文、英文訊息參考）

// vector 是 text 的語意向量（同一則訊息裡只算一次：分類和記待補問題共用）；沒有 Embedder 或算不出來時是 nil
func (t *turn) vector(text string) []float32 {
	if t.bot.Embedder == nil {
		return nil
	}
	if v, ok := t.vectors[text]; ok {
		return v
	}
	var v []float32
	if vs, err := t.bot.Embedder.Embed([]string{text}); err != nil {
		t.bot.log().Warn("算訊息的向量失敗", "error", err.Error())
	} else {
		v = vs[0]
	}
	if t.vectors == nil {
		t.vectors = map[string][]float32{}
	}
	t.vectors[text] = v
	return v
}

// similarExamples 回傳和顧客訊息（向量 v）最像的 K 句（最像的在前）；v 是 nil、案例庫是空的或算不出向量時是空的
func (b *Bot) similarExamples(v []float32, examples []merchant.Example) []merchant.Example {
	if v == nil || len(examples) == 0 {
		return nil
	}
	if err := b.embedExamples(examples); err != nil {
		b.log().Warn("算意圖案例的向量失敗，不附相似案例", "error", err.Error())
		return nil
	}
	type scored struct {
		e     merchant.Example
		score float64
	}
	b.mu.Lock()
	list := make([]scored, len(examples))
	for i, e := range examples {
		list[i] = scored{e, embed.Cosine(v, b.vectors[e.Text])}
	}
	b.mu.Unlock()
	sort.SliceStable(list, func(i, j int) bool { return list[i].score > list[j].score })
	k := b.K
	if k == 0 {
		k = DefaultK
	}
	var out []merchant.Example
	for i := 0; i < len(list) && i < k; i++ {
		out = append(out, list[i].e)
	}
	return out
}

// WarmUp 先算好意圖案例庫的向量，FAQ 超過門檻題數時也算 FAQ 問法的向量（啟動時呼叫；純 Go 的 e5 一句約 0.2 秒，案例庫兩百句約 40 秒）。沒有 Embedder 時不做事
func (b *Bot) WarmUp(ctx context.Context) error {
	if b.Embedder == nil {
		return nil
	}
	cfg, err := b.config(ctx)
	if err != nil {
		return err
	}
	if err := b.embedExamples(allExamples(cfg)); err != nil {
		return err
	}
	// FAQ 超過門檻題數時挑 FAQ 要先篩選（#65），問法的向量也先算好
	if len(cfg.FAQ) > merchant.FAQFilterAbove {
		return b.embedTexts(faqQuestions(cfg.FAQ))
	}
	return nil
}

// faqQuestions 是每題 FAQ 的主問法和其他問法
func faqQuestions(faqs []merchant.FAQ) []string {
	var texts []string
	for _, f := range faqs {
		texts = append(texts, f.Question)
		texts = append(texts, f.AltQuestions...)
	}
	return texts
}

// similarFAQ 回傳和顧客訊息（向量 v）最像的 top 題 FAQ（#65）：每題取主問法和其他問法裡最像的分數，
// 挑出的題目照原本的 FAQ 順序排（和沒篩選時的選項排法一樣）。算不出向量時 ok 是 false
func (b *Bot) similarFAQ(v []float32, faqs []merchant.FAQ, top int) (out []merchant.FAQ, ok bool) {
	if err := b.embedTexts(faqQuestions(faqs)); err != nil {
		b.log().Warn("算 FAQ 問法的向量失敗，挑 FAQ 照舊送全部的題目", "error", err.Error())
		return nil, false
	}
	idx := make([]int, len(faqs))
	score := make([]float64, len(faqs))
	b.mu.Lock()
	for i, f := range faqs {
		idx[i], score[i] = i, embed.Cosine(v, b.vectors[f.Question])
		for _, q := range f.AltQuestions {
			score[i] = max(score[i], embed.Cosine(v, b.vectors[q]))
		}
	}
	b.mu.Unlock()
	sort.SliceStable(idx, func(i, j int) bool { return score[idx[i]] > score[idx[j]] })
	idx = idx[:min(top, len(idx))]
	sort.Ints(idx)
	for _, i := range idx {
		out = append(out, faqs[i])
	}
	return out, true
}

// embedExamples 算還沒算過的案例向量（一次算完、記住）
func (b *Bot) embedExamples(examples []merchant.Example) error {
	texts := make([]string, len(examples))
	for i, e := range examples {
		texts[i] = e.Text
	}
	return b.embedTexts(texts)
}

// embedTexts 算還沒算過的句子向量（一次算完、記住）；意圖案例和 FAQ 問法共用
func (b *Bot) embedTexts(texts []string) error {
	b.mu.Lock()
	if b.vectors == nil {
		b.vectors = map[string][]float32{}
	}
	var missing []string
	seen := map[string]bool{}
	for _, s := range texts {
		if _, ok := b.vectors[s]; !ok && !seen[s] {
			missing = append(missing, s)
			seen[s] = true
		}
	}
	b.mu.Unlock()
	if len(missing) == 0 {
		return nil
	}
	vs, err := b.Embedder.Embed(missing)
	if err != nil {
		return err
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	for i, s := range missing {
		b.vectors[s] = vs[i]
	}
	return nil
}
