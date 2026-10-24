package merchant

import (
	"slices"
	"sort"
	"strings"
	"unicode/utf8"
)

// 本機模型的上下文是 8192 tokens（輸入加輸出共用），超過就答不出來（#65）。門檻依 go/cmd/ctxprobe 的量測訂：
// 森日茶屋的 FAQ 加上合成題（每題選項約 100 字）80 題 7308 tokens 成功、90 題 8208 tokens 失敗；
// 寫回覆的 3 段知識段落合計 9000 字 7847 tokens 成功（12.9 秒，接近 15 秒逾時）、10000 字失敗

const (
	// FAQFilterAbove：FAQ 超過這個題數時，挑 FAQ 先用句子向量挑出最像的 FAQFilterTop 題再交給 AI（失敗點約 85 題的六成，50 題約 5 秒）
	FAQFilterAbove = 50
	// FAQFilterTop：篩選後交給 AI 的題數（20 題約 2000 tokens，答案長一點也放得下）
	FAQFilterTop = 20
	// FAQCharsLimit：挑 FAQ 一次送出的選項（問法加答案）字數上限，超過時後台警告（失敗點約 8800 字的七成）
	FAQCharsLimit = 6000
	// KnowledgeCharsLimit：寫回覆一次最多送 3 段知識段落全文，最長 3 段合計（標題加內文）的字數上限（失敗點約 9500 字的六成多；
	// 6000 字時一題已經答成 no，每則約 8 秒）
	KnowledgeCharsLimit = 6000
	// KnowledgeSectionLimit：單段知識段落的字數上限（合計上限的一半：被挑到時另外兩段只剩一半的空間）
	KnowledgeSectionLimit = 3000
	// knowledgePicks：寫回覆一次最多送幾段（dialog 的挑知識最多挑 3 段）
	knowledgePicks = 3
)

// Warning 是後台的長度警告（#65）：不擋儲存，提醒商家內容太長時本機模型會答不出來。後台依 Kind 顯示說明：
// faq_filtered：FAQ 有 Count 題、超過門檻，每則訊息先挑最像的 Top 題給 AI；
// faq_too_long：一次送出的 FAQ 選項最多 Chars 字、超過上限 Limit，Targets 是最長的幾題的問題；
// knowledge_section_long：Targets 這幾段知識段落單段就超過上限 Limit；
// knowledge_too_long：最長 3 段（Targets）合計 Chars 字、超過上限 Limit
type Warning struct {
	Kind    string   `json:"kind"`
	Count   int      `json:"count,omitempty"`
	Top     int      `json:"top,omitempty"`
	Chars   int      `json:"chars,omitempty"`
	Limit   int      `json:"limit,omitempty"`
	Targets []string `json:"targets,omitempty"`
}

// FAQOption 是挑 FAQ 時一題的選項：代號、問法（/ 分隔）和答案
func FAQOption(f FAQ) string {
	return f.Key + "：" + strings.Join(slices.Concat([]string{f.Question}, f.AltQuestions), "/") + "（答案：" + f.Answer + "）"
}

// FAQWarnings 是 FAQ 頁的警告；沒有警告時是空陣列。超過門檻題數時，一次送出的是最像的幾題，以最長的那幾題估最多送出的字數
func FAQWarnings(faq []FAQ) []Warning {
	out := []Warning{}
	sent := len(faq)
	if sent > FAQFilterAbove {
		out = append(out, Warning{Kind: "faq_filtered", Count: len(faq), Top: FAQFilterTop})
		sent = FAQFilterTop
	}
	idx, chars := longest(len(faq), func(i int) int { return utf8.RuneCountInString(FAQOption(faq[i])) })
	total := 0
	for _, i := range idx[:sent] {
		total += chars[i]
	}
	if total > FAQCharsLimit {
		var targets []string
		for _, i := range idx[:min(3, sent)] {
			targets = append(targets, faq[i].Question)
		}
		out = append(out, Warning{Kind: "faq_too_long", Chars: total, Limit: FAQCharsLimit, Targets: targets})
	}
	return out
}

// KnowledgeWarnings 是商家資訊頁（知識段落）的警告；沒有警告時是空陣列。字數算標題加內文
func KnowledgeWarnings(sections []Section) []Warning {
	out := []Warning{}
	idx, chars := longest(len(sections), func(i int) int {
		return utf8.RuneCountInString(sections[i].Title) + utf8.RuneCountInString(sections[i].Body)
	})
	var long []string
	for _, i := range idx {
		if chars[i] > KnowledgeSectionLimit {
			long = append(long, sections[i].Title)
		}
	}
	if len(long) > 0 {
		out = append(out, Warning{Kind: "knowledge_section_long", Limit: KnowledgeSectionLimit, Targets: long})
	}
	total := 0
	var targets []string
	for _, i := range idx[:min(knowledgePicks, len(idx))] {
		total += chars[i]
		targets = append(targets, sections[i].Title)
	}
	if total > KnowledgeCharsLimit {
		out = append(out, Warning{Kind: "knowledge_too_long", Chars: total, Limit: KnowledgeCharsLimit, Targets: targets})
	}
	return out
}

// longest 回傳 0…n-1 依字數由長到短排（一樣長的照原本順序）和每項的字數
func longest(n int, length func(i int) int) (idx, chars []int) {
	idx, chars = make([]int, n), make([]int, n)
	for i := range n {
		idx[i], chars[i] = i, length(i)
	}
	sort.SliceStable(idx, func(a, b int) bool { return chars[idx[a]] > chars[idx[b]] })
	return idx, chars
}
