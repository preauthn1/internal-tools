package dialog

import (
	"context"
	"slices"
	"strings"

	"aichatbot/internal/merchant"
)

// MARK: 分類自我測試（CONTEXT.md）：用所有例句（案例庫＋自訂意圖）各分類一次，找出被分到別的意圖的例句

// SelfTestItem 是自我測試的一句：Expected 是標的意圖，Got 是分類結果（AI 失敗時是空字串）
type SelfTestItem struct {
	Text     string `json:"text"`
	State    string `json:"state,omitempty"`
	Expected string `json:"expected"`
	Got      string `json:"got"`
}

// SelfTestReport：Total 是測了幾句；Conflicts 是被分到別的意圖的例句，Failed 是 AI 失敗的例句（不算衝突）
type SelfTestReport struct {
	Total     int            `json:"total"`
	Conflicts []SelfTestItem `json:"conflicts"`
	Failed    []SelfTestItem `json:"failed"`
}

// SelfTest 讀目前的商家設定，每一句例句帶著它的狀態分類一次；分類時把這句本身排除在
// 「完全相同就沿用」、相似案例和自訂意圖選項說明的例句之外（不然每一句都會分對）。progress 每分類完一句呼叫一次（可以是 nil）。
// 有顧客回合在處理時先暫停（顧客訊息優先）。ctx 取消時停下來回 ctx.Err()
func (b *Bot) SelfTest(ctx context.Context, progress func(done, total int)) (SelfTestReport, error) {
	cfg, err := b.config(ctx)
	if err != nil {
		return SelfTestReport{}, err
	}
	examples := allExamples(cfg)
	report := SelfTestReport{Total: len(examples), Conflicts: []SelfTestItem{}, Failed: []SelfTestItem{}}
	for i, e := range examples {
		// 顧客訊息優先：每句之前等顧客的回合結束（#58）
		if err := b.waitCustomersIdle(ctx); err != nil {
			return SelfTestReport{}, err
		}
		t := &turn{bot: b, ctx: ctx, cfg: withoutExample(cfg, i), customer: customerIn(e.State), now: b.now()}
		got, _, failed := t.classifyWith(e.Text)
		item := SelfTestItem{Text: e.Text, State: e.State, Expected: e.Intent, Got: got}
		switch {
		case failed != "":
			item.Got = ""
			report.Failed = append(report.Failed, item)
		case got != e.Intent:
			report.Conflicts = append(report.Conflicts, item)
		}
		if progress != nil {
			progress(i+1, len(examples))
		}
	}
	return report, nil
}

// withoutExample 是拿掉 allExamples(cfg) 第 i 句的設定：案例庫的句子從案例庫拿掉；
// 自訂意圖的例句從那個自訂意圖拿掉（分類選項的說明也會附例句）
func withoutExample(cfg merchant.Config, i int) merchant.Config {
	if i < len(cfg.Examples) {
		cfg.Examples = slices.Delete(slices.Clone(cfg.Examples), i, i+1)
		return cfg
	}
	i -= len(cfg.Examples)
	cfg.CustomIntents = slices.Clone(cfg.CustomIntents)
	for k, ci := range cfg.CustomIntents {
		if i < len(ci.Examples) {
			cfg.CustomIntents[k].Examples = slices.Delete(slices.Clone(ci.Examples), i, i+1)
			break
		}
		i -= len(ci.Examples)
	}
	return cfg
}

// customerIn 是處在意圖案例狀態（ticket、booking:欄位）的顧客
func customerIn(state string) Customer {
	var c Customer
	if field, ok := strings.CutPrefix(state, "booking:"); ok {
		c.AskingField = field
	}
	c.OpenTicket = state == "ticket"
	return c
}
