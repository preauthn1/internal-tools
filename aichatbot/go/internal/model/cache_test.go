package model

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func fixedIdentity(id Identity) func(context.Context) (Identity, error) {
	return func(context.Context) (Identity, error) { return id, nil }
}

var haikuID = Identity{Model: "claude-haiku-5-5", Version: HaikuVersion}

// 命中：同一個請求第二次直接用存好的回答，不呼叫模型、不回報用量
func TestCacheHit(t *testing.T) {
	ctx := context.Background()
	base := NewScripted(map[string][]map[string]Value{"Intent": {{"intent": Text("booking"), "people": Integer(4)}}})
	stats := &CacheStats{}
	c := NewCache(base, t.TempDir(), fixedIdentity(haikuID), stats)

	first, _ := c.Respond(ctx, intentRequest)
	second, usage := c.Respond(ctx, intentRequest)
	if first == nil || second == nil {
		t.Fatalf("got %v %v", first, second)
	}
	if s, _ := second.String("intent"); s != "booking" {
		t.Errorf("intent：got %q", s)
	}
	if n, _ := second.Int("people"); n != 4 {
		t.Errorf("people：got %d", n)
	}
	if usage != (Usage{}) {
		t.Errorf("命中不該回報用量：%+v", usage)
	}
	if len(base.Requests()) != 1 {
		t.Errorf("模型被呼叫了 %d 次", len(base.Requests()))
	}
	if got := stats.Summary(); got != "AI 快取：命中 1 次、實際呼叫 1 次" {
		t.Errorf("Summary：got %q", got)
	}
}

// 檔名和內容都不含顧客的話和指示：只存欄位的值
func TestCacheStoresOnlyValues(t *testing.T) {
	dir := t.TempDir()
	base := NewScripted(map[string][]map[string]Value{"Intent": {{"intent": Text("booking")}}})
	NewCache(base, dir, fixedIdentity(haikuID), &CacheStats{}).Respond(context.Background(), intentRequest)
	var files []string
	filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if len(files) != 1 {
		t.Fatalf("檔案：%v", files)
	}
	data, _ := os.ReadFile(files[0])
	if string(data) != `{"intent":"booking"}` || strings.Contains(files[0], "訂位") {
		t.Errorf("%s：%s", files[0], data)
	}
}

// 失敗（nil）不存：下次重新呼叫
func TestCacheFailureNotStored(t *testing.T) {
	ctx := context.Background()
	base := NewScripted(map[string][]map[string]Value{"Intent": {nil, {"intent": Text("other")}}})
	stats := &CacheStats{}
	c := NewCache(base, t.TempDir(), fixedIdentity(haikuID), stats)
	if r, _ := c.Respond(ctx, intentRequest); r != nil {
		t.Fatalf("第一次應該失敗：%+v", r)
	}
	r, _ := c.Respond(ctx, intentRequest)
	if s, _ := r.String("intent"); s != "other" {
		t.Errorf("第二次要重新呼叫：got %+v", r)
	}
	if got := stats.Summary(); got != "AI 快取：命中 0 次、實際呼叫 2 次" {
		t.Errorf("Summary：got %q", got)
	}
}

// 未命中：指示、欄位、模型名稱、版本、系統版本任一個改了都重新呼叫
func TestCacheMissWhenRequestOrModelChanges(t *testing.T) {
	ctx := context.Background()
	dir := t.TempDir()
	instr := intentRequest
	instr.Instructions += "。"
	options := intentRequest
	options.Fields = []Field{{Kind: KindChoice, Name: "intent", Description: "意圖", Options: []string{"booking", "cancel"}}, intentRequest.Fields[1]}
	cases := []struct {
		what string
		r    Request
		id   Identity
	}{
		{"指示", instr, haikuID},
		{"選項", options, haikuID},
		{"模型名稱", intentRequest, Identity{Model: "claude-haiku-4-5", Version: HaikuVersion}},
		{"版本", intentRequest, Identity{Model: haikuID.Model, Version: "haiku/t0/max500/timeout15"}},
		{"系統版本", intentRequest, Identity{Model: haikuID.Model, Version: HaikuVersion, OS: "Version 26.2"}},
	}
	seed := NewScripted(map[string][]map[string]Value{"Intent": {{"intent": Text("booking")}}})
	NewCache(seed, dir, fixedIdentity(haikuID), &CacheStats{}).Respond(ctx, intentRequest)
	for _, c := range cases {
		base := NewScripted(map[string][]map[string]Value{"Intent": {{"intent": Text("other")}}})
		r, _ := NewCache(base, dir, fixedIdentity(c.id), &CacheStats{}).Respond(ctx, c.r)
		if s, _ := r.String("intent"); s != "other" || len(base.Requests()) != 1 {
			t.Errorf("%s改了應該重新呼叫：got %+v", c.what, r)
		}
	}
}
