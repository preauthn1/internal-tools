// embed 把案例庫和分類題組的句子先算成 e5 向量，輸出 JSON（句子 → 向量），給 Swift 版的 BOT_EMBED_FILE 讀（技術驗證用）。
//
//	go -C go run ./cmd/embed -o /tmp/e5.json [-eval ../eval] [額外的句子…]
//
// 一句一句算（不批次），stderr 印出載入模型的時間和平均一句幾毫秒（和正式環境一則訊息算一次的情況一樣）。
// 模型路徑見 internal/embed 的 ModelDir（BOT_EMBED_MODEL）。
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"aichatbot/internal/embed"
)

func main() {
	evalDir := flag.String("eval", "../eval", "eval/ 資料夾（讀 examples.jsonl、intent_set.json）")
	out := flag.String("o", "", "輸出的 JSON 檔")
	flag.Parse()
	if *out == "" {
		log.Fatal("要用 -o 指定輸出檔")
	}

	texts, err := sentences(*evalDir)
	if err != nil {
		log.Fatal(err)
	}
	// 額外的句子：分類前程式改寫過的訊息（例如相對日期換算成「10/27」），BOT_DEBUG=1 時會印「向量檔沒有這句」
	texts = append(texts, flag.Args()...)

	dir := embed.ModelDir()
	info, err := os.Stat(filepath.Join(dir, "model.onnx"))
	if err != nil {
		log.Fatal(err)
	}
	start := time.Now()
	m, err := embed.Load(dir)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()
	loaded := time.Since(start)

	vectors := make(map[string][]float32, len(texts))
	start = time.Now()
	for _, t := range texts {
		v, err := m.Embed([]string{t})
		if err != nil {
			log.Fatal(err)
		}
		vectors[t] = v[0]
	}
	per := time.Since(start) / time.Duration(len(texts))

	data, err := json.Marshal(vectors)
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(*out, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(os.Stderr, "後端 %s，model.onnx %.0f MB，載入模型 %v；%d 句，平均一句 %v；寫入 %s\n",
		embed.Backend, float64(info.Size())/1e6, loaded.Round(time.Millisecond), len(texts), per.Round(100*time.Microsecond), *out)
}

// sentences 是案例庫（examples.jsonl）和分類題組（intent_set.json）的所有句子，去掉重複
func sentences(evalDir string) ([]string, error) {
	seen := map[string]bool{}
	var texts []string
	add := func(t string) {
		if !seen[t] {
			seen[t] = true
			texts = append(texts, t)
		}
	}

	f, err := os.Open(filepath.Join(evalDir, "examples.jsonl"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var e struct{ Text string }
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, err
		}
		add(e.Text)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(evalDir, "intent_set.json"))
	if err != nil {
		return nil, err
	}
	var set struct{ Items []struct{ Text string } }
	if err := json.Unmarshal(data, &set); err != nil {
		return nil, err
	}
	for _, it := range set.Items {
		add(it.Text)
	}
	return texts, nil
}
