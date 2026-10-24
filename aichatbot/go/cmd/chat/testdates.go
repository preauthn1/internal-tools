package main

// 題庫的相對日期：把 {{date:+1:M/D}} 這類標記換成以執行當天換算的實際日期，規則和 eval/testdates.py 一致：
// {{date:±N:格式}}：基準日＝「今天 + 12 天」當天或之後的第一個星期一，值＝基準日 + 位移天數。
// {{today:±N:格式}}：今天 + N 天。
// {{weekday:next|this|near:1-7:格式}}：1＝週一，next＝下週X、this＝這週X、near＝只說週X（和對話層的相對日期換算相同）。

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

const dateFormats = `(M/D|MM-DD|M月D日|M月D號|YYYY-MM-DD)`

var (
	dateMarker    = regexp.MustCompile(`\{\{date:([+-]\d+):` + dateFormats + `\}\}`)
	todayMarker   = regexp.MustCompile(`\{\{today:([+-]\d+):` + dateFormats + `\}\}`)
	weekdayMarker = regexp.MustCompile(`\{\{weekday:(next|this|near):([1-7]):` + dateFormats + `\}\}`)
)

// now 是量測和對話用的「現在」：設了 BOT_TODAY=YYYY-MM-DD 就用那天、時間照現在
func now() time.Time {
	t := time.Now()
	if s := os.Getenv("BOT_TODAY"); s != "" {
		if d, err := time.ParseInLocation("2006-01-02", s, time.Local); err == nil {
			return time.Date(d.Year(), d.Month(), d.Day(), t.Hour(), t.Minute(), t.Second(), 0, time.Local)
		}
	}
	return t
}

func formatDate(d time.Time, f string) string {
	switch f {
	case "M/D":
		return fmt.Sprintf("%d/%d", d.Month(), d.Day())
	case "MM-DD":
		return fmt.Sprintf("%02d-%02d", d.Month(), d.Day())
	case "M月D日":
		return fmt.Sprintf("%d月%d日", d.Month(), d.Day())
	case "M月D號":
		return fmt.Sprintf("%d月%d號", d.Month(), d.Day())
	}
	return d.Format("2006-01-02")
}

func resolveTestDates(text string, today time.Time) string {
	day := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())
	start := day.AddDate(0, 0, 12)
	base := start.AddDate(0, 0, (8-int(start.Weekday()))%7) // Weekday：0=週日、1=週一
	weekdayIndex := (int(day.Weekday()) + 6) % 7            // 週一＝0
	replace := func(re *regexp.Regexp, s string, value func(groups []string) time.Time) string {
		return re.ReplaceAllStringFunc(s, func(m string) string {
			g := re.FindStringSubmatch(m)[1:]
			return formatDate(value(g), g[len(g)-1])
		})
	}
	out := replace(dateMarker, text, func(g []string) time.Time { n, _ := strconv.Atoi(g[0]); return base.AddDate(0, 0, n) })
	out = replace(todayMarker, out, func(g []string) time.Time { n, _ := strconv.Atoi(g[0]); return day.AddDate(0, 0, n) })
	return replace(weekdayMarker, out, func(g []string) time.Time {
		target, _ := strconv.Atoi(g[1])
		target--
		offset := (target - weekdayIndex + 7) % 7
		switch g[0] {
		case "next":
			offset = 7 - weekdayIndex + target
		case "this":
			offset = target - weekdayIndex
		}
		return day.AddDate(0, 0, offset)
	})
}
