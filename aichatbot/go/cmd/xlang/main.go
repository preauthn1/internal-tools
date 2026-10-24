// xlang 是跨語言抽查：拿日文、英文訊息查 eval/examples.jsonl 裡的中文案例，印出最像的前幾句和分數。
//
//	go -C go run ./cmd/xlang [-k 3] [-examples ../eval/examples.jsonl] [訊息…]
//
// 沒給訊息就用內建的抽查清單（意思和某一類中文案例相同的日文、英文句子，都不在案例庫和題組裡）。
// 模型路徑見 internal/embed 的 ModelDir（BOT_EMBED_MODEL）。
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"

	"aichatbot/internal/embed"
)

// 抽查清單：訊息 → 應該找到的意圖
var samples = []struct{ text, intent string }{
	{"予約したいです", "booking"},
	{"I want to book a table", "booking"},
	{"キャンセルでお願いします", "cancel"},
	{"never mind, cancel it", "cancel"},
	{"駐車場はありますか", "store_question"},
	{"What time do you close?", "store_question"},
	{"ありがとうございます", "thanks"},
	{"thanks a lot", "thanks"},
	{"またね", "goodbye"},
	{"see you next time", "goodbye"},
	{"あなたはロボットですか", "identity"},
	{"are you a bot?", "identity"},
	{"返事はまだですか", "follow_up"},
	{"any update on my question?", "follow_up"},
	{"最悪なサービス", "complaint"},
	{"this bot is so dumb", "complaint"},
	{"株価は上がりますか", "off_topic"},
	{"write me a poem", "off_topic"},
}

type example struct {
	Text   string `json:"text"`
	Intent string `json:"intent"`
}

// isChinese：有漢字、沒有假名（和 Swift 版 scriptOf 的分法一樣）
func isChinese(s string) bool {
	han := false
	for _, r := range s {
		if unicode.In(r, unicode.Hiragana, unicode.Katakana) {
			return false
		}
		han = han || unicode.Is(unicode.Han, r)
	}
	return han
}

func main() {
	k := flag.Int("k", 3, "印出前幾名")
	path := flag.String("examples", "../eval/examples.jsonl", "案例庫")
	flag.Parse()

	f, err := os.Open(*path)
	if err != nil {
		log.Fatal(err)
	}
	var chinese []example
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		var e example
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			log.Fatal(err)
		}
		if isChinese(e.Text) {
			chinese = append(chinese, e)
		}
	}
	f.Close()

	queries := samples
	if flag.NArg() > 0 {
		queries = nil
		for _, a := range flag.Args() {
			queries = append(queries, struct{ text, intent string }{a, ""})
		}
	}

	m, err := embed.Load(embed.ModelDir())
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()
	texts := make([]string, len(chinese))
	for i, e := range chinese {
		texts[i] = e.Text
	}
	cv, err := m.Embed(texts)
	if err != nil {
		log.Fatal(err)
	}

	hits := 0
	for _, q := range queries {
		v, err := m.Embed([]string{q.text})
		if err != nil {
			log.Fatal(err)
		}
		idx := make([]int, len(cv))
		score := make([]float64, len(cv))
		for i := range cv {
			idx[i], score[i] = i, embed.Cosine(v[0], cv[i])
		}
		sort.Slice(idx, func(a, b int) bool { return score[idx[a]] > score[idx[b]] })
		mark := ""
		if q.intent != "" {
			mark = "✘"
			if chinese[idx[0]].Intent == q.intent {
				mark, hits = "✔", hits+1
			}
			mark += " 應為 " + q.intent
		}
		fmt.Printf("%s　%s\n", q.text, mark)
		for _, i := range idx[:min(*k, len(idx))] {
			fmt.Printf("    %.3f %s → %s\n", score[i], chinese[i].Text, chinese[i].Intent)
		}
	}
	if queries[0].intent != "" {
		fmt.Printf("\n第一名意圖相同：%d/%d（中文案例 %d 句）\n", hits, len(queries), len(chinese))
	}
}
