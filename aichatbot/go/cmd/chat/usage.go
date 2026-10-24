package main

import (
	"context"
	"fmt"

	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// usageTally 包在用量紀錄外面，加總這次執行實際呼叫模型的 token 和費用（Switch 只記實際呼叫，快取命中不算）
type usageTally struct {
	model.Ledger
	input, output int
	cost          float64
}

func (u *usageTally) RecordUsage(ctx context.Context, r store.UsageRecord) error {
	u.input += r.InputTokens
	u.output += r.OutputTokens
	u.cost += r.CostUSD
	return u.Ledger.RecordUsage(ctx, r)
}

// Summary 是量測結束時印的那一行（eval/verify.py 讀；本機模型費用是 0）
func (u *usageTally) Summary() string {
	return fmt.Sprintf("AI 用量：輸入 %d token、輸出 %d token、估算 US$%.4f", u.input, u.output, u.cost)
}
