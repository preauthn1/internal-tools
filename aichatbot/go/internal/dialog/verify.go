package dialog

import (
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"aichatbot/internal/merchant"
)

// MARK: 原文片段驗證（ADR 0002）：AI 寫的回覆只有通過程式檢查才送出，不做第二次 AI 檢查

// answerCheck 是一次寫回覆要驗證的內容：出處全文、AI 抄的原文片段、AI 的回覆、顧客訊息和顧客的語言
type answerCheck struct {
	Source, Quote, Reply, Message, Lang string
}

// verifyAnswer 回傳擋下的原因；通過是空字串
func verifyAnswer(c answerCheck) string {
	switch {
	case strings.TrimSpace(c.Reply) == "":
		return "empty_reply"
	case foldText(c.Quote) == "":
		return "empty_quote"
	case utf8.RuneCountInString(c.Reply) > maxReplyLength:
		return "too_long"
	case strings.ContainsAny(c.Reply, "{}`"):
		// 本機模型偶爾把輸出格式（」}```json{）混進回覆；顧客的回覆不會有大括號和反引號
		return "malformed_reply"
	}
	if lang := ReplyLanguage(c.Reply, c.Lang); lang != "" && lang != c.Lang {
		return "language_mismatch:" + lang
	}
	if !strings.Contains(foldText(c.Source), foldText(c.Quote)) {
		return "quote_not_in_source"
	}
	allowed := facts(c.Quote + "\n" + c.Message)
	for _, f := range factList(c.Reply, false) {
		if !allowed[f] {
			return "number_not_in_quote:" + f
		}
	}
	return ""
}

// maxReplyLength 是 AI 回覆的長度上限（字數，英文也以字元計）
const maxReplyLength = 300

// ReplyLanguage 是回覆的語言（判斷不了是空字串；回答題庫也用它）。顧客用拉丁字母的語言、回覆也大多是拉丁字母時不看漢字和假名：
// 英文回覆裡附上中文品名（Green Tea（綠茶））不算換了語言
func ReplyLanguage(reply, want string) string {
	if !slices.Contains([]string{"zh", "ja", "ko", "th", "ru", "ar"}, want) {
		cjk := func(r rune) bool {
			return unicode.In(r, unicode.Han, unicode.Hiragana, unicode.Katakana, unicode.Hangul)
		}
		latin, other := 0, 0
		for _, r := range reply {
			switch {
			case cjk(r):
				other++
			case unicode.In(r, unicode.Latin):
				latin++
			}
		}
		if latin > other {
			reply = strings.Map(func(r rune) rune {
				if cjk(r) {
					return -1
				}
				return r
			}, reply)
		}
	}
	return DetectLanguage(reply)
}

var (
	numberRe    = regexp.MustCompile(`\d+(?:\.\d+)?`)
	thousandsRe = regexp.MustCompile(`(\d),(\d{3})\b`)
	urlRe       = regexp.MustCompile(`(?i)(?:https?://|www\.)[a-z0-9\-._~:/?#\[\]@!$&'()*+,;=%]+`)
	// 電話以 0 或 + 開頭（02-1234-5678、0912 345 678、+886…）；2026-10-09 這種日期不算
	phoneRe = regexp.MustCompile(`(?:\+|\b0)\d*(?:[- ]\d+)+|\b0\d{7,}`)
	// 時間：11:00、11 點、11 點半、11時、8 pm；前面的下午、晚上或後面的 pm 換成 24 小時制。
	// 只有小時、沒有單位的數字不是時間（第 4 組和第 5 組都空時照原樣留給下面當數字）
	timeRe = regexp.MustCompile(`(下午|晚上)?\s*(\d+)(?::(\d{2}))?\s*(點半|点半|點|点|時|时)?(\s*(?i:[ap]\.?m)\b\.?)?`)
	// 折扣：「9 折」（付九成）、日文的「1割」（折掉一成）
	discountRe = regexp.MustCompile(`(\d+(?:\.\d)?)\s*(折|割)`)
	// 回覆裡不當成數字的常用詞（快一點、十分抱歉、萬一、千萬別）
	numeralWords = []string{"一點", "一点", "十分", "萬一", "万一", "千萬", "千万"}
)

// factList 是一段文字裡的數字、時間（統一成「時:分」）、電話、網址，統一全形半形、中文數字、千分位後依出現順序列出。
// loose 是原文片段和顧客訊息（可以引用的一方）：每串中文數字都換；回覆只換「十一」這種兩個字以上、
// 或後面接單位（點、元、折…）的，避免「一起」「週一」「快一點」被當成數字
func factList(s string, loose bool) []string {
	var out []string
	s = urlRe.ReplaceAllStringFunc(halfWidth(s), func(m string) string {
		out = append(out, strings.ToLower(strings.TrimRight(m, ".,;:!?)'")))
		return " "
	})
	// 電話：去掉 - 和空白後 8 位以上；不到 8 位的照原樣留給下面當數字
	s = phoneRe.ReplaceAllStringFunc(s, func(m string) string {
		digits := strings.NewReplacer("-", "", " ", "").Replace(m)
		if len(strings.TrimPrefix(digits, "+")) < 8 {
			return m
		}
		out = append(out, digits)
		return " "
	})
	s = chineseNumerals(s, loose)
	for thousandsRe.MatchString(s) { // 1,234,567 要拿掉兩次
		s = thousandsRe.ReplaceAllString(s, "$1$2")
	}
	// 折扣：「9 折」是付原價的九成，和數字 9 分開記（譯成 9% off、9割引 是錯的）
	s = discountRe.ReplaceAllStringFunc(s, func(m string) string {
		g := discountRe.FindStringSubmatch(m)
		out = append(out, g[1]+g[2])
		return " "
	})
	s = timeRe.ReplaceAllStringFunc(s, func(m string) string {
		g := timeRe.FindStringSubmatch(m)
		h, _ := strconv.Atoi(g[2])
		pm := strings.HasPrefix(strings.ToLower(strings.TrimSpace(g[5])), "p")
		if (g[3] == "" && g[4] == "" && g[5] == "") || h > 24 {
			return m
		}
		if (g[1] != "" || pm) && h < 12 {
			h += 12
		}
		minutes := "00"
		switch {
		case g[3] != "":
			minutes = g[3]
		case strings.HasSuffix(g[4], "半"):
			minutes = "30"
		}
		out = append(out, fmt.Sprintf("%d:%s", h, minutes))
		return " "
	})
	for _, n := range numberRe.FindAllString(s, -1) {
		if trimmed := strings.TrimLeft(n, "0"); trimmed != "" && trimmed[0] != '.' {
			n = trimmed // 09 和 9 相同
		}
		out = append(out, n)
	}
	return out
}

// facts 是可以引用的數字、時間…：時間另外接受只寫小時（20:00 → 20）；
// 折扣另外接受折掉的百分比和成數（9 折 → 10% off、1割；85 折 → 15% off），不放行單獨的數字 9、1
func facts(s string) map[string]bool {
	set := map[string]bool{}
	for _, f := range factList(s, true) {
		set[f] = true
		if h, _, ok := strings.Cut(f, ":"); ok {
			set[h] = true
		}
		if n, ok := strings.CutSuffix(f, "折"); ok {
			// 換成百分比的整數：9 → 90、8.8 → 88、85 → 85
			whole, tenth, _ := strings.Cut(n, ".")
			if tenth == "" {
				tenth = "0"
			}
			if len(whole) == 1 {
				whole += tenth
			}
			if rate, err := strconv.Atoi(whole); err == nil && rate < 100 {
				set[strconv.Itoa(100-rate)] = true
				if (100-rate)%10 == 0 {
					set[strconv.Itoa((100-rate)/10)+"割"] = true
				}
			}
		}
	}
	return set
}

const (
	chineseDigits = "零〇一二兩三四五六七八九"
	chineseUnits  = "十百千萬"
	// 回覆裡單一個中文數字後面接這些字才當成數字
	numeralUnits = "點点元塊块折位名號号月時时分歲岁%"
)

// chineseNumerals 把中文數字換成阿拉伯數字（十一 → 11、一百零五 → 105、二〇二六 → 2026）
func chineseNumerals(s string, loose bool) string {
	rs := []rune(s)
	var b strings.Builder
	for i := 0; i < len(rs); {
		j := i
		for j < len(rs) && strings.ContainsRune(chineseDigits+chineseUnits, rs[j]) {
			j++
		}
		if j == i {
			b.WriteRune(rs[i])
			i++
			continue
		}
		run := rs[i:j]
		word := !loose && slices.ContainsFunc(numeralWords, func(w string) bool { return strings.HasPrefix(string(rs[i:]), w) })
		if !word && (loose || len(run) >= 2 || (j < len(rs) && strings.ContainsRune(numeralUnits, rs[j]))) {
			b.WriteString(strconv.Itoa(parseChineseNumber(run)))
		} else {
			b.WriteString(string(run))
		}
		i = j
	}
	return b.String()
}

func parseChineseNumber(rs []rune) int {
	digit := func(r rune) int {
		if r == '兩' {
			return 2
		}
		if r == '〇' {
			return 0
		}
		return slices.Index([]rune("零一二三四五六七八九"), r)
	}
	// 沒有十百千萬：逐位數（二〇二六）
	if !strings.ContainsFunc(string(rs), func(r rune) bool { return strings.ContainsRune(chineseUnits, r) }) {
		n := 0
		for _, r := range rs {
			n = n*10 + digit(r)
		}
		return n
	}
	total, section, cur := 0, 0, 0
	for _, r := range rs {
		switch r {
		case '十', '百', '千':
			unit := map[rune]int{'十': 10, '百': 100, '千': 1000}[r]
			if cur == 0 && r == '十' {
				cur = 1 // 十一 = 11
			}
			section, cur = section+cur*unit, 0
		case '萬':
			total, section, cur = total+(section+cur)*10000, 0, 0
		default:
			cur = digit(r)
		}
	}
	return total + section + cur
}

// foldText 是比對原文片段用的寫法：全形轉半形，只留文字和數字（空白、標點不算）
func foldText(s string) string {
	return merchant.Normalize(halfWidth(s))
}

// halfWidth 把全形英數、符號和全形空白轉成半形
func halfWidth(s string) string {
	return strings.Map(func(r rune) rune {
		switch {
		case r == '　':
			return ' '
		case r >= '！' && r <= '～':
			return r - 0xFEE0
		}
		return r
	}, s)
}
