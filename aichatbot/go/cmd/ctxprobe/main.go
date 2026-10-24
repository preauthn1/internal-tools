// ctxprobe 量本機模型上下文（8192 tokens，輸入加輸出共用）的上限（#65）：FAQ 題數對挑 FAQ、知識段落字數對寫回覆的
// token 數、每則秒數和成功或失敗，印成 markdown 表格。用森日茶屋（eval/）的資料當底，FAQ 和段落不夠時補上合成的內容；
// 送出的題目和對話層一模一樣（dialog.FAQRequest、dialog.AnswerRequest）。
//
//	go -C go run ./cmd/ctxprobe [-faq 20,50,100,200] [-chars 1000,3000,6000,8000] [-url http://0.0.0.0:8765]
//
// 直接呼叫模型服務、不開 AI 回答快取；量測前拿 ~/Library/Caches/aichatbot/measure.lock（和 eval/verify.py 同一把鎖，排隊等其他量測）。
// 失敗的原因到模型服務的 log（~/Library/Logs/aichatbot/model-server.log）對照時間查
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"aichatbot/internal/dialog"
	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
)

// 挑 FAQ 用的問題：前兩題的答案在森日茶屋原本的 FAQ 裡，第三題在合成的 FAQ 裡（題數夠多時）
var faqQuestions = []string{"你們幾點開門？", "可以刷卡嗎", "第 7 號禮盒裡面有什麼"}

// 寫回覆用的問題：答案在第一段開頭
var answerQuestions = []string{"週一有開嗎？", "最多可以幾個人訂位"}

func main() {
	faqCounts := flag.String("faq", "20,50,100,200", "FAQ 題數（逗號分隔）")
	charCounts := flag.String("chars", "1000,3000,6000,8000", "知識段落總字數（逗號分隔，分成 3 段）")
	url := flag.String("url", "http://0.0.0.0:8765", "模型服務")
	evalDir := flag.String("eval", "../eval", "森日茶屋的資料")
	flag.Parse()

	cfg, err := merchant.FromEvalDir(*evalDir)
	if err != nil {
		log.Fatal(err)
	}
	lockMeasurement()
	m := model.NewApple(*url)

	fmt.Println("## 挑 FAQ（所有題目的問法和答案一次送出）")
	fmt.Println()
	fmt.Println("| FAQ 題數 | 選項字數 | 問題 | input tokens | output tokens | 秒數 | 結果 |")
	fmt.Println("|---:|---:|---|---:|---:|---:|---|")
	for _, n := range ints(*faqCounts) {
		faqs := syntheticFAQ(cfg.FAQ, n)
		for _, q := range faqQuestions {
			req := dialog.FAQRequest(q, faqs)
			row(fmt.Sprint(n), len([]rune(req.Instructions)), q, m, req)
		}
	}
	fmt.Println()
	fmt.Println("## 寫回覆（挑到的 3 段知識段落全文送出）")
	fmt.Println()
	fmt.Println("| 段落總字數 | 指示字數 | 問題 | input tokens | output tokens | 秒數 | 結果 |")
	fmt.Println("|---:|---:|---|---:|---:|---:|---|")
	for _, n := range ints(*charCounts) {
		sections := syntheticSections(cfg.Knowledge, n)
		for _, q := range answerQuestions {
			req := dialog.AnswerRequest(cfg.Merchant.Name, "zh", q, sections)
			row(fmt.Sprint(n), len([]rune(req.Instructions)), q, m, req)
		}
	}
}

func row(size string, chars int, q string, m *model.Apple, req model.Request) {
	start := time.Now()
	r, usage := m.Respond(context.Background(), req)
	seconds := time.Since(start).Seconds()
	result := "成功"
	if r == nil {
		result = "失敗（" + start.Format("15:04:05") + "）"
		if usage.Failure != "" {
			result = usage.Failure + "（" + start.Format("15:04:05") + "）"
		}
	} else if v, ok := r.String("faq"); ok {
		result += "：" + v
	} else if v, ok := r.String("answerable"); ok {
		result += "：" + v
	}
	fmt.Printf("| %s | %d | %s | %d | %d | %.1f | %s |\n", size, chars, q, usage.InputTokens, usage.OutputTokens, seconds, result)
}

func ints(s string) []int {
	var out []int
	for _, p := range strings.Split(s, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			log.Fatalf("要是逗號分隔的整數：%s", s)
		}
		out = append(out, n)
	}
	return out
}

// syntheticFAQ 是 n 題 FAQ：先放森日茶屋原本的，不夠的補上禮盒題（問法、答案長度和原本的差不多）
func syntheticFAQ(base []merchant.FAQ, n int) []merchant.FAQ {
	out := append([]merchant.FAQ(nil), base[:min(n, len(base))]...)
	for i := 1; len(out) < n; i++ {
		out = append(out, merchant.FAQ{
			Key:          fmt.Sprintf("gift%d", i),
			Question:     fmt.Sprintf("第 %d 號禮盒有什麼", i),
			AltQuestions: []string{fmt.Sprintf("%d 號禮盒內容", i), fmt.Sprintf("gift box %d", i), fmt.Sprintf("ギフト%d番", i)},
			Answer:       fmt.Sprintf("第 %d 號禮盒有紅玉紅茶茶包 %d 包、桂花烏龍茶包 %d 包，附提袋，售價 %d 元，需三天前預訂～", i, 5+i%5, 3+i%4, 380+i*20),
		})
	}
	return out
}

// syntheticSections 是 3 段、總共約 chars 字的知識段落：每段先放森日茶屋原本的段落，再補上店家介紹類的文字到要的長度
func syntheticSections(base []merchant.Section, chars int) []merchant.Section {
	filler := []rune("本店茶葉皆向南投魚池鄉的契作茶園直接採購，每批茶葉入店前都會試泡確認香氣與回甘，不合格的批次一律退回。" +
		"店內座位採木質桌椅，靠窗座位可以看到街景，平日下午較安靜，適合閱讀或小型聚會。")
	out := make([]merchant.Section, 3)
	per := chars / 3
	for i := range out {
		s := base[i%len(base)]
		body := []rune(s.Body)
		for j := 0; len(body) < per; j++ {
			body = append(body, filler[j%len(filler)])
		}
		out[i] = merchant.Section{Title: s.Title, Body: string(body)}
	}
	return out
}

// lockMeasurement 拿量測鎖（和 eval/verify.py 的 lock_measurement 一樣）：有別的量測在跑就排隊等，程式結束時自動放掉
func lockMeasurement() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Join(home, "Library", "Caches", "aichatbot", "measure.lock")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Fatal(err)
	}
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0o644)
	if err != nil {
		log.Fatal(err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		who, _ := os.ReadFile(path)
		fmt.Fprintf(os.Stderr, "⏳ 等其他量測結束（%s）…\n", strings.TrimSpace(string(who)))
		if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX); err != nil {
			log.Fatal(err)
		}
	}
	// 寫上是誰在量（排隊的人看得到）；寫不進去不影響鎖
	if err := f.Truncate(0); err == nil {
		f.WriteAt([]byte(fmt.Sprintf("pid %d，ctxprobe，%s 開始", os.Getpid(), time.Now().Format("15:04:05"))), 0)
	}
	// f 不關：鎖在程式結束時才放掉
}
