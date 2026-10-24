package model

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"
)

// AI 回答快取：量測和驗收重跑題組時，同一個請求直接用存好的回答。照搬 Swift 版（tag swift-final 的 Sources/BotCore/AICache.swift）。
// BOT_AI_CACHE=1 才開，正式伺服器不用。包在 Queued 外面：命中時不用排隊。
// 檔名是請求的 SHA256，內容只有各欄位的值，不存明文輸入；失敗或逾時不存

// CacheStats 是命中和實際呼叫的次數（Apple、Haiku 兩個快取共用一份）
type CacheStats struct {
	hits, misses atomic.Int64
}

// Summary 是量測結束時印的那一行（eval/ 的腳本會讀）
func (s *CacheStats) Summary() string {
	return fmt.Sprintf("AI 快取：命中 %d 次、實際呼叫 %d 次", s.hits.Load(), s.misses.Load())
}

type Cache struct {
	base     ChatModel
	dir      string
	identity func(context.Context) (Identity, error)
	stats    *CacheStats
}

// NewCache：dir 是已經分好命名空間的資料夾；identity 回傳模型名稱和版本，放進 key
func NewCache(base ChatModel, dir string, identity func(context.Context) (Identity, error), stats *CacheStats) *Cache {
	return &Cache{base: base, dir: dir, identity: identity, stats: stats}
}

func (c *Cache) Respond(ctx context.Context, r Request) (*Reply, Usage) {
	id, err := c.identity(ctx)
	if err != nil {
		// 問不到版本（模型服務沒開）就不用快取，直接呼叫
		c.stats.misses.Add(1)
		return c.base.Respond(ctx, r)
	}
	file := c.path(CacheKey(r, id))
	if data, err := os.ReadFile(file); err == nil {
		var values map[string]Value
		if json.Unmarshal(data, &values) == nil {
			c.stats.hits.Add(1)
			return &Reply{Values: values}, Usage{}
		}
	}
	c.stats.misses.Add(1)
	reply, usage := c.base.Respond(ctx, r)
	if reply != nil {
		if data, err := json.Marshal(reply.Values); err == nil {
			writeAtomic(file, data)
		}
	}
	return reply, usage
}

func (c *Cache) path(key string) string {
	return filepath.Join(c.dir, key[:2], key+".json")
}

// 先寫暫存檔再改名，同時跑的程序不會讀到寫一半的檔案；寫不進去就算了（下次重新呼叫）
func writeAtomic(file string, data []byte) {
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(file), ".tmp-*")
	if err != nil {
		return
	}
	_, err = tmp.Write(data)
	if cerr := tmp.Close(); err == nil {
		err = cerr
	}
	if err == nil {
		err = os.Rename(tmp.Name(), file)
	}
	if err != nil {
		os.Remove(tmp.Name())
	}
}

// CacheKey：送給模型的全部內容＋模型名稱與版本＋系統版本，編成 JSON 陣列再算 SHA256，欄位串接不會互相碰撞。
// 開頭的 go-v1 和 Swift 版（v1）分開，格式改了就換版本
func CacheKey(r Request, id Identity) string {
	fields := make([][]string, 0, len(r.Fields))
	for _, f := range r.Fields {
		fields = append(fields, append([]string{string(f.Kind), f.Name, f.Description}, f.Options...))
	}
	data, _ := json.Marshal([]any{"go-v1", id.Model, id.Version, id.OS, r.Name, r.Instructions, r.Text, fields})
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}
