package model

import (
	"bytes"
	"encoding/json"
)

// Schema 把欄位定義轉成 JSON schema（Haiku 的 structured outputs 和 Swift 模型服務共用）。
// properties 照欄位順序輸出（map 會打亂順序，所以自己組）；required 也照欄位順序，模型服務靠它還原順序
func Schema(r Request) json.RawMessage {
	var b bytes.Buffer
	b.WriteString(`{"type":"object","properties":{`)
	required := make([]string, 0, len(r.Fields))
	for i, f := range r.Fields {
		if i > 0 {
			b.WriteByte(',')
		}
		prop := struct {
			Type        string   `json:"type"`
			Description string   `json:"description"`
			Enum        []string `json:"enum,omitempty"`
		}{Type: "string", Description: f.Description}
		switch f.Kind {
		case KindChoice:
			prop.Enum = f.Options
		case KindInteger:
			prop.Type = "integer"
		}
		writeJSON(&b, f.Name)
		b.WriteByte(':')
		writeJSON(&b, prop)
		required = append(required, f.Name)
	}
	b.WriteString(`},"required":`)
	writeJSON(&b, required)
	b.WriteString(`,"additionalProperties":false}`)
	return b.Bytes()
}

// 不跳脫 <、>、&（指示和選項裡可能有）
func writeJSON(b *bytes.Buffer, v any) {
	enc := json.NewEncoder(b)
	enc.SetEscapeHTML(false)
	enc.Encode(v)
	b.Truncate(b.Len() - 1) // Encode 會多一個換行
}
