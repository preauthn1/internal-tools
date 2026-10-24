package model

import (
	"context"
	"sync"
)

// Scripted 是假模型：依題目名稱照腳本依序回答，腳本沒寫到的題目回 nil（模擬 AI 失敗）。測試用，不必等本機模型
type Scripted struct {
	Usage Usage // 每次呼叫回報的用量（費用測試用）

	mu       sync.Mutex
	script   map[string][]map[string]Value
	requests []Request
}

func NewScripted(script map[string][]map[string]Value) *Scripted {
	return &Scripted{script: script}
}

// Requests 是收到的請求（依呼叫順序）
func (m *Scripted) Requests() []Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]Request(nil), m.requests...)
}

func (m *Scripted) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	m.mu.Lock()
	m.requests = append(m.requests, r)
	queue := m.script[r.Name]
	var answer map[string]Value
	if len(queue) > 0 {
		answer, m.script[r.Name] = queue[0], queue[1:]
	}
	m.mu.Unlock()
	if answer == nil {
		return nil, m.Usage
	}
	return &Reply{Values: answer}, m.Usage
}
