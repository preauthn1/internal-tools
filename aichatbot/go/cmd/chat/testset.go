package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"aichatbot/internal/dialog"
)

// testCase 是回答題庫（eval/testset.json、eval/faq_set.json）的一題，格式見題庫的「說明」
type testCase struct {
	ID      string   `json:"id"`
	Q       string   `json:"q"`
	Expect  string   `json:"expect"` // answer、no_answer、either
	All     []string `json:"all"`
	Any     []string `json:"any"`
	None    []string `json:"none"`
	Lang    string   `json:"lang"`
	Holdout bool     `json:"holdout"` // 驗證題：不跑（AGENTS.md：驗證題組不能看）
	Group   string   `json:"group"`   // trap＝誘答題、xlang＝跨語言題（eval/testset_extra.json），另外統計
}

// testsetEval 跑回答題庫（BOT_TESTSET，逗號分隔多個檔）：每題一位新顧客走完整的 Receive，
// 回覆種類是 answer 才算回答（FAQ 原文、譯文和 AI 寫的回覆都是）。印出各 expect 的答對數、答錯的（回答了但不對）、
// 該答沒答的（需要計算的 calc-* 另外列，ADR 0003 的已知代價）和被擋下的（驗證或譯文驗證沒過，原因取自 traces）。
// resultPath 不是空字串時，每題的結果另外寫成 JSON（eval/snapshot.py 逐題比對用），並印出「結果：路徑」。
// holdoutOnly（HOLDOUT=1）只跑驗證題，給 eval/snapshot.py 比對驗證題組（輸出只寫進不給迴圈看的檔案）
func testsetEval(ctx context.Context, bot *dialog.Bot, paths []string, out io.Writer, resultPath string, holdoutOnly bool) error {
	var cases []testCase
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var set struct {
			Cases []testCase `json:"cases"`
		}
		if err := json.Unmarshal([]byte(resolveTestDates(string(data), now())), &set); err != nil {
			return fmt.Errorf("%s：%w", path, err)
		}
		cases = append(cases, set.Cases...)
	}
	if len(cases) == 0 {
		return errors.New(strings.Join(paths, ",") + "：沒有題目")
	}
	type tally struct{ ok, total int }
	byExpect := map[string]*tally{}
	var all tally
	answered, skipped := 0, 0
	var wrong, missed, calc, blocked []string
	calcAnswered, blockedCases := 0, 0 // 一題可能擋下兩次（譯文、AI 寫的回覆），題數另外算
	traps, fabricated := 0, 0          // 誘答題：回答了而且判錯就是編造
	var xlang tally                    // 跨語言題
	type caseResult struct {
		ID     string `json:"id"`
		Result string `json:"result"` // ok、wrong（回答了但不對）、missed（沒回答）
		Key    string `json:"key"`    // 回覆種類
		Pick   string `json:"pick"`   // 挑到的 FAQ（沒挑是空字串）
	}
	results := []caseResult{}
	var elapsed time.Duration // 每題 Receive 的時間加總（開快取時含命中）
	for _, c := range cases {
		if c.Holdout != holdoutOnly {
			if c.Holdout {
				skipped++
			}
			continue
		}
		start := time.Now()
		r, err := bot.Receive(ctx, dialog.Message{Platform: "testset", UserID: c.ID, ChatID: c.ID, Text: c.Q})
		if err != nil {
			return err
		}
		elapsed += time.Since(start)
		var steps []string // 一般對話以外的意圖、挑 FAQ、翻譯、挑知識、寫回覆、驗證的結果
		aiReply := ""
		before := len(blocked)
		res := caseResult{ID: c.ID, Key: r.Key}
		for _, tr := range r.Traces {
			if tr.Step != "intent" || tr.Result != "general" {
				steps = append(steps, tr.Step+":"+tr.Result)
			}
			if tr.Step == "faq" {
				res.Pick = tr.Result
			}
			if tr.Step == "answer" || tr.Step == "translate" {
				aiReply = tr.Detail["reply"]
			}
			if (tr.Step == "verify" || tr.Step == "translate") && tr.Result != "ok" && tr.Result != "failed" {
				blocked = append(blocked, fmt.Sprintf("%s %s → %s（%s:%s）", c.ID, c.Q, aiReply, tr.Step, tr.Result))
			}
		}
		if len(blocked) > before {
			blockedCases++
		}
		isAnswer := r.Key == "answer"
		if isAnswer {
			answered++
		}
		ok, why := judge(c, isAnswer, r.Text)
		switch {
		case ok:
			res.Result = "ok"
		case isAnswer:
			res.Result = "wrong"
		default:
			res.Result = "missed"
		}
		results = append(results, res)
		if byExpect[c.Expect] == nil {
			byExpect[c.Expect] = &tally{}
		}
		byExpect[c.Expect].total++
		all.total++
		switch c.Group {
		case "trap":
			traps++
			if isAnswer && !ok {
				fabricated++
			}
		case "xlang":
			xlang.total++
			if ok {
				xlang.ok++
			}
		}
		if ok {
			byExpect[c.Expect].ok++
			all.ok++
			continue
		}
		trail := ""
		if len(steps) > 0 {
			trail = "（" + strings.Join(steps, " ") + "）"
		}
		switch {
		case strings.HasPrefix(c.ID, "calc-"):
			// 需要計算（兩杯多少、這一筆能不能刷卡）：Go 版沒有程式算價，沒答或只回通則（一杯的價格）都是已知代價
			calc = append(calc, fmt.Sprintf("%s %s　%s → [%s] %s%s", c.ID, c.Q, why, r.Key, r.Text, trail))
			if isAnswer {
				calcAnswered++
			}
		case isAnswer:
			wrong = append(wrong, fmt.Sprintf("%s %s　%s → [%s] %s%s", c.ID, c.Q, why, r.Key, r.Text, trail))
		default:
			missed = append(missed, fmt.Sprintf("%s %s → [%s] %s%s", c.ID, c.Q, r.Key, r.Text, trail))
		}
	}
	fmt.Fprintln(out, "========== 各類（expect）==========")
	for _, e := range []string{"answer", "no_answer", "either"} {
		if t := byExpect[e]; t != nil {
			fmt.Fprintf(out, "%-12s %d/%d\n", e, t.ok, t.total)
		}
	}
	for _, sec := range []struct {
		title string
		lines []string
	}{
		{"答錯的（回答了但不對）", wrong},
		{"該答沒答", missed},
		{"需要計算（ADR 0003，沒答或只回通則）", calc},
		{"被擋下的（驗證沒過）", blocked},
	} {
		fmt.Fprintf(out, "\n========== %s==========\n", sec.title)
		for _, l := range sec.lines {
			fmt.Fprintln(out, l)
		}
	}
	fmt.Fprintf(out, "\n總計：%d/%d（%.0f%%）　回答 %d 題　答錯 %d 題　該答沒答 %d 題　需要計算 %d 題（其中回答了 %d 題）　被擋下 %d 題　略過驗證題 %d 題\n",
		all.ok, all.total, float64(all.ok)/float64(max(all.total, 1))*100, answered, len(wrong), len(missed), len(calc), calcAnswered, blockedCases, skipped)
	if traps > 0 {
		fmt.Fprintf(out, "誘答題 編造 %d 題（共 %d 題）\n", fabricated, traps)
	}
	if xlang.total > 0 {
		fmt.Fprintf(out, "跨語言 答對 %d/%d\n", xlang.ok, xlang.total)
	}
	fmt.Fprintf(out, "每則平均 %.2f 秒（%d 則）\n", elapsed.Seconds()/float64(max(all.total, 1)), all.total)
	if resultPath != "" {
		data, err := json.Marshal(results)
		if err != nil {
			return err
		}
		if err := os.WriteFile(resultPath, append(data, '\n'), 0o644); err != nil {
			return err
		}
		fmt.Fprintln(out, "結果："+resultPath)
	}
	return nil
}

// judge 照 Swift 版 eval/run_testset.swift 的評分：回答了就檢查 all、any、none 和語言
func judge(c testCase, answered bool, reply string) (bool, string) {
	switch {
	case c.Expect == "no_answer" && answered:
		return false, "應該不回答，卻回答了"
	case c.Expect == "answer" && !answered:
		return false, "應回答"
	case !answered:
		return true, ""
	}
	for _, s := range c.All {
		if !strings.Contains(reply, s) {
			return false, "缺少「" + s + "」"
		}
	}
	if len(c.Any) > 0 && !containsAny(reply, c.Any) {
		return false, "沒有提到 " + strings.Join(c.Any, "/")
	}
	for _, s := range c.None {
		if strings.Contains(reply, s) {
			return false, "不該出現「" + s + "」"
		}
	}
	if c.Lang != "" {
		if lang := dialog.ReplyLanguage(reply, c.Lang); lang != c.Lang {
			return false, "語言應為 " + c.Lang + "，實際是 " + lang
		}
	}
	return true, ""
}

func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
