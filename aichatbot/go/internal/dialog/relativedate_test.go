package dialog

import (
	"fmt"
	"testing"
	"time"
)

// 照搬 Swift 版 bot-tests 的相對日期題目

func day(m, d int) time.Time { return time.Date(2026, time.Month(m), d, 12, 0, 0, 0, time.Local) }

func read(quote string, today time.Time) string {
	d, kind := readDate(quote, today)
	switch kind {
	case dateResolved:
		return fmt.Sprintf("%d/%d", d.Month(), d.Day())
	case dateUnresolved:
		return "換算不了"
	default:
		return "不是相對"
	}
}

func checkDates(t *testing.T, today time.Time, cases [][2]string) {
	t.Helper()
	for _, c := range cases {
		if got := read(c[0], today); got != c[1] {
			t.Errorf("%s：%s，應為 %s", c[0], got, c[1])
		}
	}
}

func TestReadDateOnMonday(t *testing.T) {
	checkDates(t, day(10, 5), [][2]string{
		{"今天", "10/5"}, {"明天", "10/6"}, {"後天", "10/7"}, {"tomorrow", "10/6"}, {"明日", "10/6"}, {"明後日", "10/7"},
		{"the day after tomorrow", "10/7"}, {"這週五", "10/9"}, {"這禮拜日", "10/11"}, {"this Friday", "10/9"},
		{"今週の金曜日", "10/9"}, {"下週二", "10/13"}, {"下禮拜一", "10/12"}, {"next Tuesday", "10/13"},
		{"来週の火曜日", "10/13"}, {"週二", "10/6"}, {"星期三", "10/7"}, {"禮拜天", "10/11"}, {"Wednesday", "10/7"},
		{"水曜日", "10/7"}, {"週一", "10/5"},
	})
}

func TestReadDateOnWednesday(t *testing.T) {
	checkDates(t, day(10, 7), [][2]string{
		{"今天", "10/7"}, {"明天", "10/8"}, {"後天", "10/9"}, {"Tomorrow", "10/8"}, {"這週五", "10/9"}, {"這週一", "10/5"},
		{"this Friday", "10/9"}, {"下週二", "10/13"}, {"下週二下午兩點", "10/13"}, {"下禮拜一", "10/12"}, {"下星期三", "10/14"},
		{"next Tuesday", "10/13"}, {"来週の火曜日", "10/13"}, {"星期三", "10/7"}, {"週二", "10/13"}, {"禮拜天", "10/11"},
		{"Wednesday", "10/7"}, {"水曜日", "10/7"}, {"tue", "10/13"},
	})
}

func TestReadDateOnSunday(t *testing.T) {
	checkDates(t, day(10, 11), [][2]string{
		{"明天", "10/12"}, {"後天", "10/13"}, {"tomorrow", "10/12"}, {"這週日", "10/11"}, {"這週一", "10/5"},
		{"下週一", "10/12"}, {"下週二", "10/13"}, {"next Sunday", "10/18"}, {"来週の日曜日", "10/18"},
		{"週一", "10/12"}, {"星期天", "10/11"}, {"日曜日", "10/11"}, {"Saturday", "10/17"},
	})
}

func TestReadDateUnresolvedAndAbsolute(t *testing.T) {
	var cases [][2]string
	for _, q := range []string{"下下週二", "大後天", "週末", "這週末", "this weekend", "next week", "下週", "下個月", "再過三天", "3天後", "in 3 days", "来月",
		"下下禮拜三", "下下星期三", "下下周三", "下 下週三", "我想訂下下週三下午兩點"} {
		cases = append(cases, [2]string{q, "換算不了"})
	}
	for _, q := range []string{"10/21", "10月21日", "10/21（三）", "10月21日星期三", "晚上"} {
		cases = append(cases, [2]string{q, "不是相對"})
	}
	checkDates(t, day(10, 7), cases)
}

func TestPartOfUnresolved(t *testing.T) {
	for _, c := range []struct {
		quote, text string
		want        bool
	}{
		{"下週三", "我想訂下下週三下午兩點", true},
		{"週三", "我想訂下下週三下午兩點", true},
		{"next Wednesday", "the week after next Wednesday", true},
		{"明天", "明天晚上7點，我們晚點再過去", false},
		{"下週三", "我想訂下週三下午兩點", false},
	} {
		if got := partOfUnresolved(c.quote, c.text); got != c.want {
			t.Errorf("%s／%s：%v", c.quote, c.text, got)
		}
	}
}

func TestExplicitMonthDays(t *testing.T) {
	for _, c := range [][2]string{
		{"我想訂 10/19 下午兩點", "10/19"}, {"10月19日", "10/19"}, {"10月19號晚上", "10/19"}, {"10/21（三）", "10/21"},
		{"10-19", "10/19"}, {"10.19", "10/19"}, {"10/19 14:00", "10/19"}, {"2026/10/19", "10/19"},
		{"10/19 或 10/20，10/19 優先", "10/19,10/20"}, {"電話 0912-345-678", ""}, {"14:00", ""}, {"4-6人", ""},
		{"下午兩點", ""}, {"13/40", ""},
	} {
		got := ""
		for i, md := range explicitMonthDays(c[0]) {
			if i > 0 {
				got += ","
			}
			got += fmt.Sprintf("%d/%d", md[0], md[1])
		}
		if got != c[1] {
			t.Errorf("%s：%s，應為 %s", c[0], got, c[1])
		}
	}
}

func TestAbsoluteDate(t *testing.T) {
	for _, c := range []struct {
		m, d int
		want string
	}{{10, 21, "2026-10-21"}, {10, 1, "2026-10-01"}, {9, 1, "2027-09-01"}, {1, 5, "2027-01-05"}} {
		if got := absoluteDate(c.m, c.d, day(10, 7)); got != c.want {
			t.Errorf("%d/%d：%s，應為 %s", c.m, c.d, got, c.want)
		}
	}
}
