package dialog

import (
	"context"
	"maps"
	"reflect"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
)

type values = map[string]model.Value

func txt(s string) model.Value { return model.Text(s) }
func num(n int) model.Value    { return model.Integer(n) }

// 照搬 Swift 版「假模型：預約對話走到開出預約申請」（Go 版沒有挑 FAQ 這一步）
func TestBookingConversationOpensRequest(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("booking"), intent("booking")},
		"Booking": {{},
			{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("下午2點"), "time_hour": num(14),
				"people": num(4), "people_quote": txt("4位")},
			{"name": txt("小康"), "phone": txt("0912345678")}},
	})
	if r := e.say("我想訂位喔"); r.Key != "ask_date" {
		t.Errorf("第 1 句：%+v", r)
	}
	if r := e.say("10/21 下午2點 4位"); r.Key != "ask_name" {
		t.Errorf("第 2 句：%+v", r)
	}
	r := e.say("小康 0912345678")
	summary := "2026-10-21 14:00, 4 位, 小康, 0912345678"
	if r.Key != "booking_done" || !strings.Contains(r.Text, summary) {
		t.Errorf("第 3 句：%+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "booking：📅 預約申請："+summary {
		t.Errorf("待處理：%q", p)
	}
	b, err := e.db.Booking(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if b.Date != "2026-10-21" || b.Time != "14:00" || b.People != "4" || b.Name != "小康" || b.Phone != "0912345678" || b.Status != "pending" || b.Summary != summary {
		t.Errorf("預約：%+v", b)
	}
	c, _ := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
	if c.AskingField != "" || len(c.Booking) != 0 {
		t.Errorf("預約完要清掉：%+v", c)
	}
	var names []string
	for _, r := range e.model.Requests() {
		names = append(names, r.Name)
	}
	if got := strings.Join(names, ","); got != "Intent,Booking,Intent,Booking,Intent,Booking" {
		t.Errorf("AI 呼叫：%s", got)
	}
}

// #42：問姓名、電話那輪 AI 把人數抽成 1、出處抄成「王小明」，蓋掉第一句的 4 人
func TestBookingKeepsPeopleWhenLaterTurnMisextracts(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("booking"), intent("booking")},
		"Booking": {
			{"date_quote": txt("明天"), "time_quote": txt("晚上7點"), "time_hour": num(19), "people": num(4), "people_quote": txt("4位")},
			{"name": txt("王小明"), "people": num(1), "people_quote": txt("王小明")},
			{"phone": txt("0912345678"), "people": num(1), "date_quote": txt("0912345678")}},
	})
	if r := e.say("我想訂位，明天晚上7點，4位"); r.Key != "ask_name" {
		t.Errorf("第 1 句：%+v", r)
	}
	if r := e.say("王小明"); r.Key != "ask_phone" {
		t.Errorf("第 2 句：%+v", r)
	}
	r := e.say("0912345678")
	if summary := "2026-10-08 19:00, 4 位, 王小明, 0912345678"; r.Key != "booking_done" || !strings.Contains(r.Text, summary) {
		t.Errorf("第 3 句：%+v", r)
	}
	if b, err := e.db.Booking(context.Background(), 1); err != nil || b.People != "4" {
		t.Errorf("預約：%+v %v", b, err)
	}
}

// 已經有人數、這輪不是在問人數：要有人數的出處才改，整句碰巧有同一個數字不算
func TestBookingKeepsPeopleWithoutQuote(t *testing.T) {
	e := midBooking(t, "zh", "name", "booking", values{"name": txt("王小明"), "people": num(1)})
	if r := e.say("王小明，可能會晚1小時到"); r.Key != "ask_phone" || e.customer().Booking["people"] != "4" {
		t.Errorf("%+v %+v", r, e.customer())
	}
	e = midBooking(t, "zh", "name", "booking", values{"name": txt("王小明"), "people": num(5), "people_quote": txt("改5位")})
	if r := e.say("王小明，改5位"); r.Key != "ask_phone" || e.customer().Booking["people"] != "5" {
		t.Errorf("有出處要改：%+v %+v", r, e.customer())
	}
}

// 人數的各種寫法：出處由程式換算成人數，要和 AI 填的值一致才收，對不上就再問人數
func TestBookingPeopleWordings(t *testing.T) {
	for _, c := range []struct {
		text, quote string
		people      int
		want        string
	}{
		{"想訂四個人的位子", "四個人", 4, "4"},
		{"想訂四個人的位子", "四個人", 1, ""},
		{"想訂位 2大1小", "2大1小", 3, "3"},
		// #53：AI 只填了大人的 2，大人小孩的寫法由程式加總
		{"想訂位 2大1小", "2大1小", 2, "3"},
		// AI 沒填人數也沒抄出處：程式從整句讀
		{"下週六晚上6點，2大1小，想預約", "", 0, "3"},
		{"想訂位，兩大一小", "", 0, "3"},
		{"想訂位，大人兩位小孩一位", "", 0, "3"},
		{"想訂位，大人兩位小孩一位", "大人兩位", 2, "3"},
		{"a table for 4 please", "for 4", 4, "4"},
		{"我想訂位，4位", "4位", 1, ""},
		{"我想訂位，4位", "", 1, ""},
		{"我想訂位，明天晚上7點，4位", "", 7, ""},
		{"我想訂位，明天晚上7點，4位", "晚上7點", 7, ""},
	} {
		_, b := bookOnce(t, c.text, values{"people": num(c.people), "people_quote": txt(c.quote)})
		if b["people"] != c.want {
			t.Errorf("%s／AI 填 %d：%q，應為 %q", c.text, c.people, b["people"], c.want)
		}
	}
}

// bookOnce：今天 2026-10-07（週三），一句話抽預約欄位；回傳回覆種類和存下的預約欄位
func bookOnce(t *testing.T, text string, booking values) (string, map[string]string) {
	t.Helper()
	e := newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("booking")}, "Booking": {booking}})
	r := e.say(text)
	c, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
	if err != nil {
		t.Fatal(err)
	}
	return r.Key, c.Booking
}

func sameBooking(got, want map[string]string) bool {
	if len(got) != len(want) {
		return false
	}
	for k, v := range want {
		if got[k] != v {
			return false
		}
	}
	return true
}

// 照搬 Swift 版的日期、時間一致性檢查（10/19 原本回公休，ADR 0003 之後照收）
func TestBookingDateAndTimeConsistency(t *testing.T) {
	for _, c := range []struct {
		name, text string
		ai         values
		key        string
		want       map[string]string
	}{
		{"下週二用程式換算，不用 AI 填的月日", "我想預約下週二下午2點 4位",
			values{"date_quote": txt("下週二"), "date_month": num(10), "date_day": num(10), "time_quote": txt("下午2點"), "time_hour": num(14), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-13", "time": "14:00", "people": "4"}},
		{"換算不了的相對說法不填日期", "下下週三下午2點",
			values{"date_quote": txt("下下週三"), "date_month": num(10), "date_day": num(21), "time_quote": txt("下午2點"), "time_hour": num(14)},
			"ask_date", map[string]string{"time": "14:00"}},
		{"AI 只抄到「下週三」", "我想訂下下週三下午兩點",
			values{"date_quote": txt("下週三"), "date_month": num(10), "date_day": num(14), "time_quote": txt("下午兩點"), "time_hour": num(14)},
			"ask_date", map[string]string{"time": "14:00"}},
		{"AI 只抄到「週三」", "我想訂下下週三下午兩點",
			values{"date_quote": txt("週三"), "date_month": num(10), "date_day": num(14), "time_quote": txt("下午兩點"), "time_hour": num(14)},
			"ask_date", map[string]string{"time": "14:00"}},
		{"寫了月日：AI 填 0 也讀得到", "我想訂 10/19 下午兩點",
			values{"date_quote": txt("10/19"), "date_month": num(0), "date_day": num(0), "time_quote": txt("下午兩點"), "time_hour": num(14)},
			"ask_people", map[string]string{"date": "2026-10-19", "time": "14:00"}},
		{"日期和時間出處重疊", "我想訂 10/19 下午兩點",
			values{"date_quote": txt("10/19 下午兩點"), "time_quote": txt("10/19 下午兩點"), "time_hour": num(14)},
			"ask_people", map[string]string{"date": "2026-10-19", "time": "14:00"}},
		{"沒有出處", "我想訂 10/19 下午兩點",
			values{"time_hour": num(14)},
			"ask_time", map[string]string{"date": "2026-10-19"}},
		{"AI 填錯的日也不用", "10/21 下午2點 4位",
			values{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(12), "time_quote": txt("下午2點"), "time_hour": num(14), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}},
		{"寫了兩個日期，AI 填的不是其中之一", "10/21 或 10/22 下午2點",
			values{"date_quote": txt("10/21 或 10/22"), "date_month": num(10), "date_day": num(12), "time_quote": txt("下午2點"), "time_hour": num(14)},
			"ask_date", map[string]string{"time": "14:00"}},
		{"寫了兩個日期，AI 填的是其中之一", "10/21 或 10/22 下午2點",
			values{"date_quote": txt("10/21 或 10/22"), "date_month": num(10), "date_day": num(22), "time_quote": txt("下午2點"), "time_hour": num(14)},
			"ask_people", map[string]string{"date": "2026-10-22", "time": "14:00"}},
		{"「21號」AI 填的月日和出處不一致", "21號 下午2點 4位",
			values{"date_quote": txt("21號"), "date_month": num(10), "date_day": num(12), "time_quote": txt("下午2點"), "time_hour": num(14), "people": num(4), "people_quote": txt("4位")},
			"ask_date", map[string]string{"time": "14:00", "people": "4"}},
		{"時和出處不一致：用出處讀出的時", "10/21 下午2點 4位",
			values{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("下午2點"), "time_hour": num(18), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}},
		{"出處讀不出時（7:30）：AI 填的時和原文不一致就不收", "10/21 7:30 4位",
			values{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("7:30"), "time_hour": num(8), "time_minute": num(30), "people": num(4), "people_quote": txt("4位")},
			"ask_time", map[string]string{"date": "2026-10-21", "people": "4"}},
		{"兩點半的分沒有數字", "明天兩點半 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("兩點半"), "time_hour": num(14), "time_minute": num(30), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "14:30", "people": "4"}},
		{"人數和出處的數字不一致", "10/21 下午2點 25 個人",
			values{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("下午2點"), "time_hour": num(14), "people": num(10), "people_quote": txt("25 個人")},
			"ask_people", map[string]string{"date": "2026-10-21", "time": "14:00"}},
		{"AI 把下午兩點填成 2：原文寫了下午就加 12", "明天下午兩點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("下午兩點"), "time_hour": num(2), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "14:00", "people": "4"}},
		{"7pm 填成 7", "tomorrow 7pm for 2",
			values{"date_quote": txt("tomorrow"), "time_quote": txt("7pm"), "time_hour": num(7), "people": num(2), "people_quote": txt("2")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "19:00", "people": "2"}},
		{"早上十點維持 10 點", "明天早上10點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("早上10點"), "time_hour": num(10), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "10:00", "people": "4"}},
		{"#46 晚上七點填成 6：用出處讀出的時", "下週三晚上七點兩個人，可以幫我訂嗎",
			values{"date_quote": txt("下週三"), "time_quote": txt("晚上七點"), "time_hour": num(6), "people": num(2), "people_quote": txt("兩個人")},
			"ask_name", map[string]string{"date": "2026-10-14", "time": "19:00", "people": "2"}},
		{"#46 實際的 AI 輸出：出處只抄「七點」、填 18，晚上看出處前面的字", "下週三晚上七點兩個人，可以幫我訂嗎",
			values{"date_quote": txt("下週三"), "time_quote": txt("七點"), "time_hour": num(18), "people": num(2), "people_quote": txt("兩個人")},
			"ask_name", map[string]string{"date": "2026-10-14", "time": "19:00", "people": "2"}},
		{"晚上七點填 7", "下週三晚上七點兩個人，可以幫我訂嗎",
			values{"date_quote": txt("下週三"), "time_quote": txt("晚上七點"), "time_hour": num(7), "people": num(2), "people_quote": txt("兩個人")},
			"ask_name", map[string]string{"date": "2026-10-14", "time": "19:00", "people": "2"}},
		{"晚上七點填 19", "下週三晚上七點兩個人，可以幫我訂嗎",
			values{"date_quote": txt("下週三"), "time_quote": txt("晚上七點"), "time_hour": num(19), "people": num(2), "people_quote": txt("兩個人")},
			"ask_name", map[string]string{"date": "2026-10-14", "time": "19:00", "people": "2"}},
		{"下午三點半", "明天下午三點半 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("下午三點半"), "time_hour": num(3), "time_minute": num(30), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "15:30", "people": "4"}},
		{"中午十二點不加 12", "明天中午十二點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午十二點"), "time_hour": num(12), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "12:00", "people": "4"}},
		{"七點十五沒寫分：分照 AI 填的，要在出處裡", "明天晚上七點十五 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("晚上七點十五"), "time_hour": num(19), "time_minute": num(15), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "19:15", "people": "4"}},
		{"早上十點 AI 填 22：原文寫了早上就換回上午", "明天早上十點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("早上十點"), "time_hour": num(22), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "10:00", "people": "4"}},
		{"上午、下午都寫了：AI 填 14 不換回上午", "明天上午不行，下午2點 4位",
			values{"date_quote": txt("明天"), "time_hour": num(14), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "14:00", "people": "4"}},
		{"出處連星期一起抄：「週五七點」不能讀成 57 點", "下週五七點 3位",
			values{"date_quote": txt("下週五"), "time_quote": txt("下週五七點"), "time_hour": num(19), "people": num(3), "people_quote": txt("3位")},
			"ask_name", map[string]string{"date": "2026-10-16", "time": "19:00", "people": "3"}},
		{"出處只抄「七點」：晚上和時間中間隔了字", "明天晚上 大概七點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("七點"), "time_hour": num(18), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "19:00", "people": "4"}},
		{"整句寫了兩個時間：不用出處讀出的時，AI 填的對不上就再問", "明天想訂八點，上次七點去太晚 2位",
			values{"date_quote": txt("明天"), "time_quote": txt("七點"), "time_hour": num(8), "people": num(2), "people_quote": txt("2位")},
			"ask_time", map[string]string{"date": "2026-10-08", "people": "2"}},
		{"七點的分填成 7 不收", "明天晚上七點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("晚上七點"), "time_hour": num(19), "time_minute": num(7), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "19:00", "people": "4"}},
		{"中午十二點 AI 填 0", "明天中午十二點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午十二點"), "time_hour": num(0), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "12:00", "people": "4"}},
		{"早上十點不加 12", "明天早上十點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("早上十點"), "time_hour": num(10), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "10:00", "people": "4"}},
		{"傍晚六點填 6", "明天傍晚六點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("傍晚六點"), "time_hour": num(6), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "18:00", "people": "4"}},
		{"早上十點填成 11：用出處讀出的時", "明天早上十點 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("早上十點"), "time_hour": num(11), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "10:00", "people": "4"}},
		{"十點十分的分由程式讀", "明天早上十點十分 4位",
			values{"date_quote": txt("明天"), "time_quote": txt("早上十點十分"), "time_hour": num(10), "people": num(4), "people_quote": txt("4位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "10:10", "people": "4"}},
		{"#49 中午一點填 1：加 12", "我想訂明天中午一點，3位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午一點"), "time_hour": num(1), "people": num(3), "people_quote": txt("3位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "13:00", "people": "3"}},
		{"#49 今晚七點填 7：加 12", "今晚七點 2位",
			values{"date_quote": txt("今晚"), "time_quote": txt("今晚七點"), "time_hour": num(7), "people": num(2), "people_quote": txt("2位")},
			"ask_name", map[string]string{"date": "2026-10-07", "time": "19:00", "people": "2"}},
		{"#49 今天晚上七點：出處只抄「七點」", "今天晚上七點 2位",
			values{"date_quote": txt("今天"), "time_quote": txt("七點"), "time_hour": num(7), "people": num(2), "people_quote": txt("2位")},
			"ask_name", map[string]string{"date": "2026-10-07", "time": "19:00", "people": "2"}},
		{"#49 中午一點填 13", "我想訂明天中午一點，3位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午一點"), "time_hour": num(13), "people": num(3), "people_quote": txt("3位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "13:00", "people": "3"}},
		{"#49 中午十一點半照字面", "明天中午十一點半 3位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午十一點半"), "time_hour": num(11), "time_minute": num(30), "people": num(3), "people_quote": txt("3位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "11:30", "people": "3"}},
		{"#49 半夜三點填 3：不加 12", "明天半夜三點 2位",
			values{"date_quote": txt("明天"), "time_quote": txt("半夜三點"), "time_hour": num(3), "people": num(2), "people_quote": txt("2位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "03:00", "people": "2"}},
		{"#49 晚上十二點是跨日寫法：不收，再問一次", "明天晚上十二點 2位",
			values{"date_quote": txt("明天"), "time_quote": txt("晚上十二點"), "time_hour": num(12), "people": num(2), "people_quote": txt("2位")},
			"ask_time_midnight", map[string]string{"date": "2026-10-08", "people": "2"}},
		{"#49 今晚 12 個人沒說時間：不是晚上十二點", "今晚12個人有位嗎",
			values{"date_quote": txt("今晚"), "time_hour": num(0), "people": num(12), "people_quote": txt("12個人")},
			"ask_time", map[string]string{"date": "2026-10-07", "people": "12"}},
		{"#49 中午12點、後面才說晚上：照收 12:00", "明天中午12點，晚上不行，2位",
			values{"date_quote": txt("明天"), "time_quote": txt("中午12點"), "time_hour": num(12), "people": num(2), "people_quote": txt("2位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "12:00", "people": "2"}},
		{"#49 12am 是跨日寫法", "tomorrow 12am for 2",
			values{"date_quote": txt("tomorrow"), "time_quote": txt("12am"), "time_hour": num(0), "people": num(2), "people_quote": txt("for 2")},
			"ask_time_midnight", map[string]string{"date": "2026-10-08", "people": "2"}},
		{"#49 深夜十一點是 23 點", "明天深夜十一點 2位",
			values{"date_quote": txt("明天"), "time_quote": txt("深夜十一點"), "time_hour": num(11), "people": num(2), "people_quote": txt("2位")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "23:00", "people": "2"}},
		{"#49 6 in the afternoon 是 18 點", "tomorrow 6 in the afternoon for 2",
			values{"date_quote": txt("tomorrow"), "time_quote": txt("6 in the afternoon"), "time_hour": num(6), "people": num(2), "people_quote": txt("for 2")},
			"ask_name", map[string]string{"date": "2026-10-08", "time": "18:00", "people": "2"}},
		{"已經過了的日期", "10/5 下午2點",
			values{"date_quote": txt("10/5"), "time_quote": txt("下午2點"), "time_hour": num(14)},
			"err_past", map[string]string{"time": "14:00"}},
	} {
		key, b := bookOnce(t, c.text, c.ai)
		if key != c.key || !sameBooking(b, c.want) {
			t.Errorf("%s：%s %v，應為 %s %v", c.name, key, b, c.key, c.want)
		}
	}
}

// #53：問人數時顧客報的是姓名、電話：照樣收下、回頭問人數；人數問了 2 次還是沒有就轉人工（附已收到的預約資料），不再一直問
func TestBookingHandsOffAfterAskingSameFieldTwice(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("booking"), intent("booking"), intent("booking"), intent("booking")},
		"Booking": {
			{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("晚上6點"), "time_hour": num(18)},
			{"name": txt("黃志明")},
			{"phone": txt("0955667788")},
			{},
			{}},
	})
	// 有收到新的預約資料（姓名、電話）的那一輪不算（使用者決定，aichatbot!71）：顧客在配合，只是沒先回人數
	for i, s := range []string{"10/21 晚上6點想預約", "黃志明", "0955667788"} {
		if r := e.say(s); r.Key != "ask_people" {
			t.Errorf("第 %d 句：%+v", i+1, r)
		}
	}
	if c := e.customer(); c.Booking["name"] != "黃志明" || c.Booking["phone"] != "0955667788" {
		t.Errorf("姓名、電話要收下：%+v", c)
	}
	// 沒有進展的輪次才累計：再問一次（第 2 次），還是沒有就轉人工
	if r := e.say("嗯"); r.Key != "ask_people" {
		t.Errorf("第 4 句：%+v", r)
	}
	if r := e.say("嗯嗯"); r.Key != "handoff" {
		t.Errorf("第 5 句：%+v", r)
	}
	c := e.customer()
	if !c.Human || c.AskingField != "" || len(c.Booking) != 0 {
		t.Errorf("轉人工、清掉預約：%+v", c)
	}
	want := "handoff：🙋 轉人工：預約資料收不齊（人數）：日期 2026-10-21、時間 18:00、預約人姓名 黃志明、電話 0955667788"
	if p := e.pending(); len(p) != 1 || p[0] != want {
		t.Errorf("待處理：%q", p)
	}
}

// #53：問人數時先報姓名，再問一次人數就給了：照常問電話、完成預約，不轉人工
func TestBookingNameBeforePeopleStillCompletes(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("booking"), intent("booking"), intent("booking")},
		"Booking": {
			{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21), "time_quote": txt("晚上6點"), "time_hour": num(18)},
			{"name": txt("黃志明")},
			{"people": num(3), "people_quote": txt("3位")},
			{"phone": txt("0955667788")}},
	})
	var keys []string
	for _, s := range []string{"10/21 晚上6點想預約", "黃志明", "3位", "0955667788"} {
		keys = append(keys, e.say(s).Key)
	}
	if got := strings.Join(keys, ","); got != "ask_people,ask_people,ask_phone,booking_done" {
		t.Errorf("got %s", got)
	}
}

// #53：預約中說了別的話、一般對話再問欄位（reask）也算一次：時間問了 2 次還是沒有就轉人工
func TestBookingReaskCountsTowardHandoff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent":  {intent("booking"), intent("general"), intent("general")},
		"Booking": {{"date_quote": txt("10/21"), "date_month": num(10), "date_day": num(21)}, {}, {}},
	})
	var keys []string
	for _, s := range []string{"10/21 想預約", "嗯嗯", "嗯嗯"} {
		keys = append(keys, e.say(s).Key)
	}
	if got := strings.Join(keys, ","); got != "ask_time,ask_time,handoff" {
		t.Errorf("got %s", got)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 轉人工：預約資料收不齊（時間）：日期 2026-10-21" {
		t.Errorf("待處理：%q", p)
	}
}

// #53：欄位格式不對（電話位數不對）也算問了一次：電話連續錯就轉人工，不一直回 err_phone
func TestBookingFormatErrorsCountTowardHandoff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent":  {intent("booking"), intent("booking"), intent("booking")},
		"Booking": {{"phone": txt("0912")}, {"phone": txt("0913")}, {"phone": txt("0914")}},
	})
	c := Customer{Lang: "zh", AskingField: "phone", AskCount: 1, Booking: map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康"}}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", c); err != nil {
		t.Fatal(err)
	}
	var keys []string
	for _, s := range []string{"0912", "0913"} {
		keys = append(keys, e.say(s).Key)
	}
	if got := strings.Join(keys, ","); got != "err_phone,handoff" {
		t.Errorf("got %s", got)
	}
}

// #53：轉人工關閉時回 handoff_off（附聯絡方式），預約狀態一樣清掉，不再重問
func TestBookingStuckWithHandoffOff(t *testing.T) {
	e := newEnv(t, setFeature(testMerchant(), "handoff", false), map[string][]values{
		"Intent":  {intent("booking")},
		"Booking": {{}},
	})
	c := Customer{Lang: "zh", AskingField: "people", AskCount: 2, Booking: map[string]string{"date": "2026-10-21", "time": "14:00"}}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", c); err != nil {
		t.Fatal(err)
	}
	if r := e.say("還在想"); r.Key != "handoff_off" || !strings.Contains(r.Text, "電話 02-1234-5678") {
		t.Errorf("%+v", r)
	}
	if c := e.customer(); c.Human || c.AskingField != "" || c.AskCount != 0 || len(c.Booking) != 0 || len(e.pending()) != 0 {
		t.Errorf("%+v %q", c, e.pending())
	}
}

// #53：取消預約時連問了幾次一起清掉
func TestBookingCancelClearsAskCount(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("cancel")}})
	c := Customer{Lang: "zh", AskingField: "people", AskCount: 2, Booking: map[string]string{"date": "2026-10-21"}}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", c); err != nil {
		t.Fatal(err)
	}
	if r := e.say("不訂了"); r.Key != "booking_cancel" {
		t.Errorf("%+v", r)
	}
	if c := e.customer(); c.AskingField != "" || c.AskCount != 0 || len(c.Booking) != 0 {
		t.Errorf("%+v", c)
	}
}

// #49：晚上十二點再問一次；顧客連同日期寫明 0 點就收成 00:00（ADR 0003 不檢查營業時間）
func TestBookingMidnightAsksAgain(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("booking")},
		"Booking": {
			{"date_quote": txt("明天"), "time_quote": txt("半夜十二點"), "time_hour": num(0), "people": num(2), "people_quote": txt("2位")},
			{"date_quote": txt("10/9"), "date_month": num(10), "date_day": num(9), "time_quote": txt("凌晨0點"), "time_hour": num(0)}},
	})
	if r := e.say("明天半夜十二點 2位"); r.Key != "ask_time_midnight" {
		t.Errorf("第 1 句：%+v", r)
	}
	if r := e.say("10/9 凌晨0點"); r.Key != "ask_name" {
		t.Errorf("第 2 句：%+v", r)
	}
	if b := e.customer().Booking; b["date"] != "2026-10-09" || b["time"] != "00:00" {
		t.Errorf("預約：%v", b)
	}
}

// #49：問時間時回沒有時間寫法的話（報名字），AI 自己填的時不收、再問一次；短答「七點」「7pm」「十點半」照收
func TestBookingAskingTimeNeedsTimeWording(t *testing.T) {
	for _, c := range []struct {
		text string
		ai   values
		key  string
		time string
	}{
		{"林小華", values{"name": txt("林小華"), "time_hour": num(14)}, "ask_time", ""},
		{"林小華 0912345678", values{"name": txt("林小華"), "phone": txt("0912345678"), "time_hour": num(14)}, "ask_time", ""},
		{"I am John", values{"name": txt("John"), "time_hour": num(14)}, "ask_time", ""},
		{"七點", values{"time_quote": txt("七點"), "time_hour": num(7)}, "ask_name", "07:00"},
		{"7pm", values{"time_hour": num(7)}, "ask_name", "19:00"},
		{"十點半", values{"time_hour": num(10), "time_minute": num(30)}, "ask_name", "10:30"},
	} {
		e := newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("booking")}, "Booking": {c.ai}})
		cu := Customer{Lang: "zh", AskingField: "time", Booking: map[string]string{"date": "2026-10-21", "people": "4"}}
		if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", cu); err != nil {
			t.Fatal(err)
		}
		r := e.say(c.text)
		if got := e.customer().Booking["time"]; r.Key != c.key || got != c.time {
			t.Errorf("%s：%s %q，應為 %s %q", c.text, r.Key, got, c.key, c.time)
		}
	}
}

// midBooking：預約問到 asking 欄位（日期、時間、人數已經有了）的顧客，這句分類成 intent、AI 抽到 booking
func midBooking(t *testing.T, lang, asking, intentCode string, booking values) *env {
	t.Helper()
	e := newEnv(t, testMerchant(), map[string][]values{"Intent": {intent(intentCode)}, "Booking": {booking}})
	c := Customer{Lang: lang, AskingField: asking, Booking: map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}}
	if asking == "phone" {
		c.Booking["name"] = "小康"
	}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", c); err != nil {
		t.Fatal(err)
	}
	return e
}

func (e *env) customer() Customer {
	e.t.Helper()
	c, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
	if err != nil {
		e.t.Fatal(err)
	}
	return c
}

// Swift 版 #17：預約問欄位時回外文名字、電話，不切換語言、不當成不支援的語言；改說完整的外文句子才切換
func TestShortAnswerKeepsLanguageWhileBooking(t *testing.T) {
	for _, name := range []string{"Bill", "たなか", "김민수"} {
		e := midBooking(t, "zh", "name", "booking", values{"name": txt(name)})
		r := e.say(name)
		if c := e.customer(); c.Lang != "zh" || c.Booking["name"] != name {
			t.Errorf("%s：%+v", name, c)
		}
		if r.Key != "ask_phone" || r.Text != "請留一支聯絡電話～" {
			t.Errorf("%s：%+v", name, r)
		}
		if p := e.pending(); len(p) != 0 {
			t.Errorf("%s：待處理 %q", name, p)
		}
	}
	e := midBooking(t, "zh", "phone", "booking", values{"phone": txt("0912-345-678")})
	if r := e.say("0912-345-678"); r.Key != "booking_done" || e.customer().Lang != "zh" {
		t.Errorf("電話：%+v %+v", r, e.customer())
	}
	e = midBooking(t, "ja", "name", "booking", values{"name": txt("田中")})
	if e.say("田中"); e.customer().Lang != "ja" {
		t.Errorf("日文顧客回漢字名字：%+v", e.customer())
	}
	e = midBooking(t, "zh", "name", "general", nil)
	if e.say("What time do you close today?"); e.customer().Lang != "en" {
		t.Errorf("完整英文句子：%+v", e.customer())
	}
	// 沒有在預約時，短的外文照舊切換
	e = newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("general")}})
	if e.say("Thanks"); e.customer().Lang != "en" {
		t.Errorf("沒在預約：%+v", e.customer())
	}
}

// 照搬 Swift 版 #17 第二版：問姓名時 AI 抽的名字不能用，短答整句就是名字；
// 只在分類成預約、自我介紹時用，整句漢字也不用（「不知道」「好」曾被收成名字）
func TestShortAnswerAsName(t *testing.T) {
	for _, c := range []struct {
		text, intent, aiName, key, name string
	}{
		{"たなか", "booking", "田中", "ask_phone", "たなか"},
		{"たなか", "booking", "Tanaka", "ask_phone", "たなか"},
		{"たなか", "booking", "", "ask_phone", "たなか"},
		{"たなか", "self_intro", "", "ask_phone", "たなか"},
		{"Bill", "general", "Bill", "ask_phone", "Bill"},
		{"不知道", "unclear", "", "unclear", ""},
		{"不知道", "booking", "", "ask_name", ""},
		{"好", "general", "", "ask_name", ""}, // 一般對話答不出來時再問一次（F33）
		{"Okay", "general", "", "ask_name", ""},
		{"0923456781", "booking", "", "ask_name", ""},
		{"What time do you close today?", "booking", "", "ask_name", ""},
		// 角色稱呼不是名字
		{"顧客", "booking", "顧客", "ask_name", ""},
		{"customer", "booking", "", "ask_name", ""},
	} {
		e := midBooking(t, "zh", "name", c.intent, values{"name": txt(c.aiName)})
		r := e.say(c.text)
		got := e.customer()
		if r.Key != c.key || got.Booking["name"] != c.name {
			t.Errorf("%s／%s／AI 抽到 %q：%s %q，應為 %s %q", c.text, c.intent, c.aiName, r.Key, got.Booking["name"], c.key, c.name)
		}
		if c.name == "" && got.AskingField != "name" {
			t.Errorf("%s／%s：還要繼續問姓名，現在是 %q", c.text, c.intent, got.AskingField)
		}
	}
}

// 問姓名時只有短答可以不靠自我介紹句型：改問完整的問句，句首的地名不能當成名字（#50）
func TestAskingNameQuestionIsNotName(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]values{"Intent": {named("general", "台灣")}, "Booking": {{}}})
	c := Customer{Lang: "zh", AskingField: "name", Booking: map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4"}}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", c); err != nil {
		t.Fatal(err)
	}
	e.say("台灣有幾間分店呀")
	if got := e.customer(); got.Name != "" || got.Booking["name"] != "" {
		t.Errorf("記成名字：%q／預約姓名 %q", got.Name, got.Booking["name"])
	}
}

// 預約中取消：清掉預約；沒有在預約時說取消當成一般對話
func TestCancel(t *testing.T) {
	e := midBooking(t, "zh", "time", "cancel", nil)
	if r := e.say("算了取消"); r.Key != "booking_cancel" {
		t.Errorf("預約中：%+v", r)
	}
	if c := e.customer(); c.AskingField != "" || len(c.Booking) != 0 {
		t.Errorf("要清掉預約：%+v", c)
	}
	e = newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("cancel")}})
	if r := e.say("不用了"); r.Key != "general" {
		t.Errorf("沒在預約：%+v", r)
	}
	if n := len(e.model.Requests()); n != 1 {
		t.Errorf("AI 呼叫 %d 次", n)
	}
}

// 預約中聲稱、追問、想找真人：不抽預約欄位，跳出去處理，預約狀態保留
func TestBookingInterruptedByClaim(t *testing.T) {
	e := midBooking(t, "zh", "people", "claim", values{"people": num(5)})
	if r := e.say("上次說可以打五折"); r.Key != "claim" {
		t.Errorf("got %+v", r)
	}
	if c := e.customer(); c.AskingField != "people" || c.Booking["people"] != "4" {
		t.Errorf("預約狀態：%+v", c)
	}
	for _, r := range e.model.Requests() {
		if r.Name == "Booking" {
			t.Errorf("不該抽預約欄位")
		}
	}
}

// 預約關閉：回固定回覆 booking_off（帶入商家的聯絡方式；商家改過的優先），不抽預約欄位、不轉人工
func TestBookingSwitchOff(t *testing.T) {
	script := func() map[string][]values {
		return map[string][]values{"Intent": {intent("booking")}, "Booking": {{"people": num(4), "people_quote": txt("4位")}}}
	}
	e := newEnv(t, setFeature(testMerchant(), "booking", false), script())
	if r := e.say("我想訂位 4位"); r.Key != "booking_off" || r.Text != "目前沒有開放線上預約，可以透過這個方式聯絡我們：電話 02-1234-5678" {
		t.Errorf("預設文字：%+v", r)
	}
	if c := e.customer(); c.AskingField != "" || len(c.Booking) != 0 || c.Human {
		t.Errorf("預設文字：%+v", c)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("預設文字：待處理 %q", p)
	}
	if n := len(e.model.Requests()); n != 1 {
		t.Errorf("預設文字：AI 呼叫 %d 次", n)
	}

	c := setFeature(testMerchant(), "booking", false)
	c.Replies = []merchant.Reply{{Key: "booking_off", Lang: "zh", Text: "目前不開放線上預約，請來電：{contact}"}}
	e = newEnv(t, c, script())
	if r := e.say("我想訂位 4位"); r.Key != "booking_off" || r.Text != "目前不開放線上預約，請來電：電話 02-1234-5678" {
		t.Errorf("商家改過：%+v", r)
	}
}

// 極短訊息規則：1～2 個字（純英文 3 個字母以內）被分成預約時改成看不懂；預約中不套用（顧客可能只回「4」）
func TestTinyMessageIsNotBooking(t *testing.T) {
	for _, text := range []string{"中", "阿", "abc"} {
		e := newEnv(t, testMerchant(), map[string][]values{"Intent": {intent("booking")}})
		if r := e.say(text); r.Key != "unclear" {
			t.Errorf("%s：%+v", text, r)
		}
	}
	e := midBooking(t, "zh", "people", "booking", values{"people": num(5)})
	if r := e.say("5"); r.Key != "ask_name" || e.customer().Booking["people"] != "5" {
		t.Errorf("預約中回人數：%+v %+v", r, e.customer())
	}
}

// #67：沒設定預約欄位時，送給 AI 的預約抽取和改版前逐字相同（問人數那輪的提示、欄位與說明）
func TestBookingPromptUnchangedWithoutBookingFields(t *testing.T) {
	e := midBooking(t, "zh", "people", "booking", values{"people": num(5)})
	e.say("5")
	var req model.Request
	for _, r := range e.model.Requests() {
		if r.Name == "Booking" {
			req = r
		}
	}
	instructions := "從顧客訊息抽取預約資訊。今天是 2026-10-07。客服剛剛問了顧客的人數，這句話很可能就是在回答它。" +
		"只說相對日期（明天、下週二…）時照抄原文，月、日填 0，程式會換算；有寫月日就照填。顧客可能使用任何語言。"
	if req.Instructions != instructions {
		t.Errorf("Instructions：%s", req.Instructions)
	}
	fields := []model.Field{
		{Kind: model.KindText, Name: "date_quote", Description: "訊息中提到日期的原文，逐字照抄；未提及則為空字串"},
		{Kind: model.KindInteger, Name: "date_month", Description: "預約日期的月；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"},
		{Kind: model.KindInteger, Name: "date_day", Description: "預約日期的日；有寫月日就照填；只說相對日期（明天、下週二…）或未提及則為 0"},
		{Kind: model.KindText, Name: "time_quote", Description: "訊息中提到時間的原文，逐字照抄；未提及則為空字串"},
		{Kind: model.KindInteger, Name: "time_hour", Description: "時，24 小時制（下午兩點=14）；未提及則為 0"},
		{Kind: model.KindInteger, Name: "time_minute", Description: "分（兩點半=30）；未提及則為 0"},
		{Kind: model.KindInteger, Name: "people", Description: "人數；未提及則為 0"},
		{Kind: model.KindText, Name: "people_quote", Description: "訊息中提到人數的原文，逐字照抄；未提及則為空字串"},
		{Kind: model.KindText, Name: "name", Description: "預約人姓名；未提及則為空字串"},
		{Kind: model.KindText, Name: "phone", Description: "電話號碼；未提及則為空字串"},
	}
	if !reflect.DeepEqual(req.Fields, fields) {
		t.Errorf("Fields：%+v", req.Fields)
	}
}

// salonMerchant 是不問人數、有服務項目（必填）和指定設計師（選填）的商家
func salonMerchant() merchant.Config {
	c := testMerchant()
	c.BookingFields = &merchant.BookingFields{PeopleEnabled: false, Custom: []merchant.BookingField{
		{Key: "service", Name: "服務項目", Description: "要預約的服務", Required: true},
		{Key: "designer", Name: "指定設計師", Description: "指定的設計師名字"},
	}}
	return c
}

// #67：人數關掉、有必填的服務項目：不問人數，日期時間之後問服務項目（問的時候短答整句就是答案），不存人數
func TestBookingWithCustomFields(t *testing.T) {
	e := newEnv(t, salonMerchant(), map[string][]values{
		"Intent": {intent("booking"), intent("general"), intent("booking"), intent("booking")},
		"Booking": {
			{"date_quote": txt("10/9"), "date_month": num(10), "date_day": num(9), "time_quote": txt("下午三點"), "time_hour": num(15), "people": num(3)},
			{"service": txt("")},
			{"name": txt("林小姐")},
			{"phone": txt("0911222333")}},
	})
	if r := e.say("我想約 10/9 下午三點"); r.Key != "ask_service" || r.Text != "請問服務項目是？" {
		t.Errorf("第 1 句：%+v", r)
	}
	if r := e.say("剪髮加護髮"); r.Key != "ask_name" {
		t.Errorf("第 2 句：%+v %+v", r, e.customer())
	}
	e.say("林小姐")
	r := e.say("0911222333")
	summary := "2026-10-09 15:00, 剪髮加護髮, 林小姐, 0911222333"
	if r.Key != "booking_done" || !strings.Contains(r.Text, summary) {
		t.Errorf("第 4 句：%+v", r)
	}
	b, err := e.db.Booking(context.Background(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if b.People != "" || b.Extra["service"] != "剪髮加護髮" || b.Summary != summary {
		t.Errorf("預約：%+v", b)
	}
	var names []string
	for _, f := range e.model.Requests()[1].Fields {
		names = append(names, f.Name)
	}
	if got := strings.Join(names, ","); got != "date_quote,date_month,date_day,time_quote,time_hour,time_minute,name,phone,service,designer" {
		t.Errorf("抽取欄位：%s", got)
	}
}

// #67：預約途中商家關掉人數：之前收到的人數不存進預約
func TestBookingDropsPeopleWhenDisabled(t *testing.T) {
	c := testMerchant()
	c.BookingFields = &merchant.BookingFields{PeopleEnabled: false}
	e := newEnv(t, c, map[string][]values{"Intent": {intent("booking")}, "Booking": {{"phone": txt("0912345678")}}})
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "phone",
		Booking: map[string]string{"date": "2026-10-21", "time": "14:00", "people": "4", "name": "小康"}}); err != nil {
		t.Fatal(err)
	}
	if r := e.say("0912345678"); r.Key != "booking_done" {
		t.Fatalf("%+v", r)
	}
	if b, err := e.db.Booking(context.Background(), 1); err != nil || b.People != "" || b.Summary != "2026-10-21 14:00, 小康, 0912345678" {
		t.Errorf("預約：%+v %v", b, err)
	}
}

// #67：自訂欄位要出現在原文裡才收（AI 自己編的不收）；選填的顧客有講才記、沒講不追問；問自訂欄位時反問不當成答案
func TestBookingCustomFieldsNeedQuote(t *testing.T) {
	salon := func(text string, booking values) (string, map[string]string) {
		e := newEnv(t, salonMerchant(), map[string][]values{"Intent": {intent("booking")}, "Booking": {booking}})
		return e.say(text).Key, e.customer().Booking
	}
	when := values{"date_quote": txt("10/10"), "date_month": num(10), "date_day": num(10), "time_quote": txt("早上十一點"), "time_hour": num(11)}
	with := func(extra values) values {
		out := maps.Clone(when)
		maps.Copy(out, extra)
		return out
	}
	key, b := salon("幫我約 10/10 早上十一點染髮，指定小葵", with(values{"service": txt("染髮"), "designer": txt("小葵")}))
	if key != "ask_name" || b["service"] != "染髮" || b["designer"] != "小葵" {
		t.Errorf("有講服務和設計師：%s %v", key, b)
	}
	key, b = salon("幫我約 10/10 早上十一點", with(values{"service": txt("剪髮"), "designer": txt("Mia")}))
	if key != "ask_service" || b["service"] != "" || b["designer"] != "" {
		t.Errorf("AI 自己編的：%s %v", key, b)
	}
	key, b = salon("幫我約 10/10 早上十一點燙髮", with(values{"service": txt("燙髮")}))
	if key != "ask_name" || b["designer"] != "" {
		t.Errorf("沒講設計師不追問：%s %v", key, b)
	}

	e := newEnv(t, salonMerchant(), map[string][]values{"Intent": {intent("general")}, "Booking": {{}}})
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "service",
		Booking: map[string]string{"date": "2026-10-10", "time": "11:00"}}); err != nil {
		t.Fatal(err)
	}
	e.say("有哪些可以選？")
	if c := e.customer(); c.Booking["service"] != "" {
		t.Errorf("反問被當成服務項目：%+v", c)
	}
}

// #67：問自訂欄位時，分類的提示寫出欄位名稱；意圖案例的狀態 booking:<代號> 也用欄位名稱說明
func TestClassifyHintNamesCustomField(t *testing.T) {
	c := salonMerchant()
	c.Examples = append(c.Examples, merchant.Example{Text: "剪瀏海", Intent: "booking", State: "booking:service"})
	e := newEnv(t, c, map[string][]values{"Intent": {intent("booking")}, "Booking": {{}}})
	e.bot.Embedder = &fakeEmbedder{vectors: map[string][]float32{"剪瀏海": {1, 0, 0}, "剪瀏海和洗髮": {1, 0, 0}}}
	if err := e.bot.saveCustomer(context.Background(), "telegram", "U1", Customer{Lang: "zh", AskingField: "service",
		Booking: map[string]string{"date": "2026-10-10", "time": "11:00"}}); err != nil {
		t.Fatal(err)
	}
	e.say("剪瀏海和洗髮")
	ins := e.model.Requests()[0].Instructions
	if !strings.Contains(ins, "剛剛問了顧客的服務項目。") || !strings.Contains(ins, "「剪瀏海」（預約中問服務項目時） → booking") {
		t.Errorf("分類的提示：%s", ins)
	}
}

// ADR 0003：不檢查時段、公休日和人數上限（10/19 是週一、晚上十點、25 人照收），由員工確認或拒絕
func TestBookingSkipsHoursClosedDayAndPeopleLimit(t *testing.T) {
	key, b := bookOnce(t, "我想訂 10/19 晚上10點 25 個人", values{"date_quote": txt("10/19"), "date_month": num(10), "date_day": num(19),
		"time_quote": txt("晚上10點"), "time_hour": num(22), "people": num(25), "people_quote": txt("25 個人")})
	if key != "ask_name" || !sameBooking(b, map[string]string{"date": "2026-10-19", "time": "22:00", "people": "25"}) {
		t.Errorf("%s：%v", key, b)
	}
}
