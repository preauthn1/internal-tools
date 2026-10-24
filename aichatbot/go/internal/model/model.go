// Package model 是模型層：給顧客的話、指示和欄位定義，回傳各欄位的值（失敗或逾時回 nil）。
// 業務程式只透過 ChatModel 呼叫 AI。實作有 Haiku、Apple 本機模型（呼叫 Swift 模型服務）和照腳本回答的假模型；
// 排隊、AI 回答快取、費用記錄和月上限都是包在外面的一層。做法照搬 Swift 版（tag swift-final 的 Sources/BotCore/Model.swift）
package model

import (
	"context"
	"encoding/json"
	"fmt"
)

type FieldKind string

const (
	KindChoice  FieldKind = "choice"
	KindText    FieldKind = "text"
	KindInteger FieldKind = "integer"
)

// Field 是一個要模型填的欄位；Options 只有選擇題用
type Field struct {
	Kind        FieldKind
	Name        string
	Description string
	Options     []string
}

type Request struct {
	Name         string // 題目名稱（Intent、Pick、Booking…）
	Text         string // 顧客的話
	Instructions string
	Fields       []Field // 欄位依序產生，順序會影響結果
}

// Value 是欄位的值：文字或整數；JSON 直接寫成字串或數字（和 Swift 版的 FieldValue 一樣）
type Value struct {
	s     string
	n     int
	isInt bool
}

func Text(s string) Value { return Value{s: s} }
func Integer(n int) Value { return Value{n: n, isInt: true} }

func (v Value) MarshalJSON() ([]byte, error) {
	if v.isInt {
		return json.Marshal(v.n)
	}
	return json.Marshal(v.s)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	var n int
	if err := json.Unmarshal(data, &n); err == nil {
		*v = Integer(n)
		return nil
	}
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("欄位的值要是字串或整數：%s", data)
	}
	*v = Text(s)
	return nil
}

// Reply 是各欄位的值；讀不到的欄位不在裡面
type Reply struct {
	Values map[string]Value
}

func (r *Reply) String(name string) (string, bool) {
	v, ok := r.Values[name]
	if !ok || v.isInt {
		return "", false
	}
	return v.s, true
}

func (r *Reply) Int(name string) (int, bool) {
	v, ok := r.Values[name]
	if !ok || !v.isInt {
		return 0, false
	}
	return v.n, true
}

// parseValues 每個欄位各自讀，讀不到或型別不符的略過（其他欄位照用，和 Swift 版一樣）；整個不是 JSON 物件時回 false
func parseValues(fields []Field, data []byte) (map[string]Value, bool) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, false
	}
	values := map[string]Value{}
	for _, f := range fields {
		v, ok := raw[f.Name]
		if !ok {
			continue
		}
		if f.Kind == KindInteger {
			var n int
			if json.Unmarshal(v, &n) == nil {
				values[f.Name] = Integer(n)
			}
		} else {
			var s string
			if json.Unmarshal(v, &s) == nil {
				values[f.Name] = Text(s)
			}
		}
	}
	return values, true
}

// Usage 是一次呼叫用掉的 token；Model 是空字串表示沒有實際呼叫模型（快取命中、連不上）
type Usage struct {
	Model        string
	InputTokens  int
	OutputTokens int
	Failure      string // 失敗的原因，只有失敗時才有值：FailureContextExceeded、FailureBlockedBySafety；其他失敗（逾時、連不上、本機模型其他錯誤）是空字串
}

// FailureContextExceeded：提示詞超過本機模型的上下文（8192 tokens，輸入加輸出共用，#65）
const FailureContextExceeded = "too_long"

// FailureBlockedBySafety：被 Apple 本機模型的安全過濾擋下（May contain unsafe content，#68）
const FailureBlockedBySafety = "blocked_by_safety"

type ChatModel interface {
	// Respond 失敗或逾時回 nil；Usage 不管成功失敗都照實回報
	Respond(ctx context.Context, r Request) (*Reply, Usage)
}
