package dialog

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// MARK: 相對日期：「明天」「下週二」由程式依今天換算（AI 自己算曾把週三說的「下週二」算成週六）
// AI 只抄日期的原文出處，這裡判斷是不是相對說法、換算成哪一天。純函式，不碰 AI。照搬 Swift 版的 RelativeDate.swift

type dateKind int

const (
	dateNotRelative dateKind = iota // 不是相對說法（10/21、10月21日…）：照絕對日期處理
	dateResolved                    // 換算出的日期（當天 0 點）
	dateUnresolved                  // 是相對說法但換算不了（下下週、週末、下個月…）：不要填，改問顧客
)

func mustAll(patterns ...string) []*regexp.Regexp {
	out := make([]*regexp.Regexp, len(patterns))
	for i, p := range patterns {
		out[i] = regexp.MustCompile(p)
	}
	return out
}

func anyMatch(res []*regexp.Regexp, s string) bool {
	for _, re := range res {
		if re.MatchString(s) {
			return true
		}
	}
	return false
}

// 換算不了的相對說法（先比對，免得「下下週二」被當成「下週二」、「大後天」被當成「後天」）
var unresolvedPatterns = mustAll(
	`下\s*下\s*(週|周|禮拜|礼拜|星期|個|个)`, `大\s*(後|后)\s*天`, `再\s*(來|来)\s*週`,
	`(週|周)\s*末`, `weekend`, `土日`,
	`(下|這|这|本|上)\s*(個|个)?\s*月`, `next\s+month`, `来月`, `來月`, `月底`, `月末`,
	`(過|过)\s*(幾|几)\s*天`, `再\s*(過|过)`, `(天|日|週|周|星期|禮拜|礼拜|月)\s*(以)?(後|后)`,
	`\bin\s+\S+\s+(days?|weeks?)`, `(days?|weeks?)\s+(later|from\s+now)`, `week\s+after\s+next`,
	`改天`, `某天`, `哪天`, `\blater\b`, `\bsomeday\b`,
)

// 今天＋N 天的說法（「後天」要在「明天」之前比對；「day after tomorrow」要在「tomorrow」之前）
var dayWords = []struct {
	re     *regexp.Regexp
	offset int
}{
	{regexp.MustCompile(`day\s+after\s+tomorrow`), 2}, {regexp.MustCompile(`明後日`), 2}, {regexp.MustCompile(`あさって`), 2},
	{regexp.MustCompile(`後天`), 2}, {regexp.MustCompile(`后天`), 2},
	{regexp.MustCompile(`明天`), 1}, {regexp.MustCompile(`明日`), 1}, {regexp.MustCompile(`明晚`), 1},
	{regexp.MustCompile(`tomorrow`), 1}, {regexp.MustCompile(`あした`), 1},
	{regexp.MustCompile(`今天`), 0}, {regexp.MustCompile(`今日`), 0}, {regexp.MustCompile(`本日`), 0},
	{regexp.MustCompile(`今晚`), 0}, {regexp.MustCompile(`today`), 0}, {regexp.MustCompile(`tonight`), 0},
}

// 星期幾（週一＝0 … 週日＝6）
var (
	chineseWeekday  = regexp.MustCompile(`(週|周|星期|禮拜|礼拜)\s*([一二三四五六日天])`)
	japaneseWeekday = regexp.MustCompile(`([月火水木金土日])曜`)
	englishWeekday  = regexp.MustCompile(`\b(monday|tuesday|wednesday|thursday|friday|saturday|sunday|mon|tues|tue|wed|thurs|thur|thu|fri|sat|sun)\b`)
	weekdayKeys     = [][]string{{"一", "月", "mon"}, {"二", "火", "tue"}, {"三", "水", "wed"}, {"四", "木", "thu"}, {"五", "金", "fri"}, {"六", "土", "sat"}, {"日", "天", "sun"}}
)

func weekdayIn(text string) (int, bool) {
	for _, re := range []*regexp.Regexp{chineseWeekday, japaneseWeekday, englishWeekday} {
		m := re.FindString(text)
		if m == "" {
			continue
		}
		// 中文的「週」「星期」本身不含星期幾，只看最後一個字；日文看「曜」前面那個字；英文看開頭
		r := []rune(m)
		var key string
		switch re {
		case englishWeekday:
			key = string(r[:3])
		case japaneseWeekday:
			key = string(r[:1])
		default:
			key = string(r[len(r)-1:])
		}
		for i, names := range weekdayKeys {
			for _, n := range names {
				if n == key {
					return i, true
				}
			}
		}
	}
	return 0, false
}

var (
	nextWeek     = mustAll(`下\s*(個|个)?\s*(週|周|禮拜|礼拜|星期)`, `next\s+week`, `\bnext\s+(mon|tue|wed|thu|fri|sat|sun)`, `来週`, `來週`, `らいしゅう`)
	thisWeek     = mustAll(`(這|这|今)\s*(個|个)?\s*(週|周|禮拜|礼拜|星期)`, `本\s*(週|周)`, `this\s+week`, `\bthis\s+(mon|tue|wed|thu|fri|sat|sun)`, `今週`)
	writtenMonth = regexp.MustCompile(`\d+\s*[/\-.月]\s*\d+|\d+\s*(日|號|号)`)
)

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// readDate 讀日期的原文出處：相對說法依 today 換算。「下週」是下一個星期一開始的那一週；只說「週X」是今天起最近的那一天（含今天）
func readDate(quote string, today time.Time) (time.Time, dateKind) {
	q := strings.ToLower(quote)
	// 寫了月日（10/21（三）、10月21日星期三）就是絕對日期，星期幾只是附註
	if writtenMonth.MatchString(q) {
		return time.Time{}, dateNotRelative
	}
	if anyMatch(unresolvedPatterns, q) {
		return time.Time{}, dateUnresolved
	}
	start := startOfDay(today)
	plus := func(days int) (time.Time, dateKind) { return start.AddDate(0, 0, days), dateResolved }
	for _, w := range dayWords {
		if w.re.MatchString(q) {
			return plus(w.offset)
		}
	}
	todayIndex := (int(start.Weekday()) + 6) % 7 // Weekday：0=週日
	isNext, isThis := anyMatch(nextWeek, q), anyMatch(thisWeek, q)
	target, ok := weekdayIn(q)
	switch {
	case !ok && (isNext || isThis):
		return time.Time{}, dateUnresolved
	case !ok:
		return time.Time{}, dateNotRelative
	case isNext:
		return plus(7 - todayIndex + target) // 下一個星期一＋星期幾
	case isThis:
		return plus(target - todayIndex) // 本週一＋星期幾（可能已經過了，交給「已經過了」的檢查）
	}
	return plus((target - todayIndex + 7) % 7)
}

// partOfUnresolved：出處是不是原句裡某個換算不了的說法的一部分（AI 曾把「下下週三」只抄成「下週三」，被當成下週三）。
// 只看和出處重疊的說法，句子別處的「晚點再過去」不影響
func partOfUnresolved(quote, text string) bool {
	t := strings.ToLower(text)
	quotes := regexp.MustCompile(regexp.QuoteMeta(strings.ToLower(quote))).FindAllStringIndex(t, -1)
	for _, re := range unresolvedPatterns {
		for _, u := range re.FindAllStringIndex(t, -1) {
			for _, q := range quotes {
				if q[0] < u[1] && u[0] < q[1] {
					return true
				}
			}
		}
	}
	return false
}

// 明確寫出的月日（10/19、10-19、10.19、10月19日）：前後不能接著數字（電話、年份），後面不能是人數或時間（4-6人、2.5點）
var (
	monthDayPattern = regexp.MustCompile(`(\d{1,2})\s*(?:[/\-.]|月)\s*(\d{1,2})`)
	unitAfter       = regexp.MustCompile(`^\s*(個|个|位|人|點|点|時|时|分)`)
)

func isDigit(b byte) bool { return b >= '0' && b <= '9' }

// explicitMonthDays 讀出訊息裡明確寫的月日（不靠 AI 填的整數；AI 曾把「10/19」的月日填成 0）。只留合法範圍，重複的去掉
func explicitMonthDays(text string) [][2]int {
	var found [][2]int
	// Go 的 regexp 沒有 lookbehind：被前後的數字擋掉時，從下一個字元重新找（「2026/10/19」才找得到 10/19）
	for pos := 0; pos < len(text); {
		m := monthDayPattern.FindStringSubmatchIndex(text[pos:])
		if m == nil {
			break
		}
		base := pos
		start, end := base+m[0], base+m[1]
		pos = start + 1
		if (start > 0 && isDigit(text[start-1])) || (end < len(text) && isDigit(text[end])) || unitAfter.MatchString(text[end:]) {
			continue
		}
		month, _ := strconv.Atoi(text[base+m[2] : base+m[3]])
		d, _ := strconv.Atoi(text[base+m[4] : base+m[5]])
		md := [2]int{month, d}
		if month < 1 || month > 12 || d < 1 || d > 31 || containsMD(found, md) {
			continue
		}
		found = append(found, md)
		pos = end
	}
	return found
}

func containsMD(list [][2]int, md [2]int) bool {
	for _, x := range list {
		if x == md {
			return true
		}
	}
	return false
}

// absoluteDate：只有月、日的絕對日期補上年份：今年；這天在今天往前 7 天之前就算明年（7 天內維持今年，回「這天已經過了」）。
// 不存在的日期（2/30）照樣組成字串，交給日期檢查回「看不懂」
func absoluteDate(month, d int, today time.Time) string {
	year := today.Year()
	t := time.Date(year, time.Month(month), d, 0, 0, 0, 0, today.Location())
	valid := t.Month() == time.Month(month) && t.Day() == d
	if valid && t.Before(startOfDay(today).AddDate(0, 0, -7)) {
		year++
	}
	return fmt.Sprintf("%04d-%02d-%02d", year, month, d)
}
