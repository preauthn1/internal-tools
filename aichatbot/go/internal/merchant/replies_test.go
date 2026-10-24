package merchant

import (
	"context"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
)

// 後台只列出可以改的 5 則固定回覆、啟用的語言；| 分隔的多句拆開；* 那筆每種語言都有；
// 還沒設定的那則，慣用語言用預設文字（不含店名）、其他語言空白
func TestRepliesOf(t *testing.T) {
	c := valid()
	c.Languages = []Language{{"zh", true}, {"en", true}, {"ko", false}, {"ja", true}}
	c.Replies = []Reply{
		{"fallback", "zh", "不會列出來"},
		{"greeting", "en", "Hi|Hello"},
		{"greeting", "zh", "您好|哈囉"},
		{"handoff", "*", "請稍候"},
	}
	got, _ := json.Marshal(RepliesOf(c))
	want := `[` +
		`{"key":"greeting","texts":{"zh":["您好","哈囉"],"en":["Hi","Hello"],"ja":[]}},` +
		`{"key":"booking_off","texts":{"zh":["目前沒有開放線上預約，可以透過這個方式聯絡我們：{contact}"],"en":[],"ja":[]}},` +
		`{"key":"handoff","texts":{"zh":["請稍候"],"en":["請稍候"],"ja":["請稍候"]}},` +
		`{"key":"handoff_off","texts":{"zh":["目前沒有真人客服，可以透過這個方式聯絡我們：{contact}"],"en":[],"ja":[]}},` +
		`{"key":"off_hours","texts":{"zh":["現在不是店員服務時間，店員上班後會盡快回覆您喔！"],"en":[],"ja":[]}}]`
	var a, b any
	json.Unmarshal(got, &a)
	json.Unmarshal([]byte(want), &b)
	ja, _ := json.Marshal(a)
	jb, _ := json.Marshal(b)
	if string(ja) != string(jb) {
		t.Errorf("got  %s\nwant %s", ja, jb)
	}

	// 這則只有沒啟用的語言有內容時，慣用語言照樣給預設文字
	c.Replies = []Reply{{"off_hours", "ko", "지금은 영업시간이 아닙니다"}}
	if zh := RepliesOf(c)[4].Texts["zh"]; !slices.Equal(zh, []string{"現在不是店員服務時間，店員上班後會盡快回覆您喔！"}) {
		t.Errorf("慣用語言：%v", zh)
	}
}

// 存一則固定回覆：多句用 | 接起來、去掉前後空白；只換掉這則送來的語言，
// 其他固定回覆、沒啟用的語言、* 那筆（沒有自己那筆的語言還要用）、其他區塊都不動
func TestSaveReplyKeepsOthers(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c := valid()
	c.Languages = []Language{{"zh", true}, {"en", true}, {"ko", false}, {"th", false}}
	c.Replies = []Reply{{"fallback", "zh", "轉給店員"}, {"handoff", "*", "請稍候"}, {"handoff", "ko", "잠시만요"}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	if err := SaveReply(ctx, db, "handoff", map[string][]string{"zh": {" 稍等喔 ", "馬上來"}, "en": {"One moment"}}); err != nil {
		t.Fatal(err)
	}
	got, err := Load(ctx, db)
	if err != nil {
		t.Fatal(err)
	}
	want := []Reply{{"fallback", "zh", "轉給店員"}, {"handoff", "*", "請稍候"}, {"handoff", "en", "One moment"}, {"handoff", "ko", "잠시만요"}, {"handoff", "zh", "稍等喔|馬上來"}}
	if !slices.Equal(got.Replies, want) {
		t.Errorf("got  %v\nwant %v", got.Replies, want)
	}
	// 之後勾選泰文時，泰文還是 * 那筆
	got.Languages[3].Enabled = true
	if th := RepliesOf(got)[2].Texts["th"]; !slices.Equal(th, []string{"請稍候"}) {
		t.Errorf("泰文：%v", th)
	}
	if got.Merchant != c.Merchant || len(got.Languages) != 4 {
		t.Errorf("其他區塊被改了：%+v", got)
	}
}

// 只能改那 4 則、只能存啟用的語言；慣用語言至少一句；句子不能是空白、不能含 |（| 是分隔多句用的）。沒過的不會存進去
func TestSaveReplyValidation(t *testing.T) {
	ctx := context.Background()
	db := openDB(t).SQL()
	c := valid()
	c.Languages = []Language{{"zh", true}, {"en", true}, {"ko", false}}
	if _, err := Save(ctx, db, c); err != nil {
		t.Fatal(err)
	}
	for _, x := range []struct {
		key   string
		texts map[string][]string
		want  string
	}{
		{"fallback", map[string][]string{"zh": {"好"}}, "不能改的固定回覆：fallback"},
		{"handoff", map[string][]string{"zh": {"好"}, "ko": {"네"}}, "語言 ko 沒有啟用"},
		{"handoff", map[string][]string{"zh": {}, "en": {"OK"}}, "慣用語言（zh）至少要一句"},
		{"handoff", map[string][]string{"en": {"OK"}}, "慣用語言（zh）至少要一句"},
		{"handoff", map[string][]string{"zh": {"好", " "}}, "zh 的第 2 句是空白"},
		{"handoff", map[string][]string{"zh": {"好|壞"}}, "zh 的第 1 句不能有 |"},
	} {
		err := SaveReply(ctx, db, x.key, x.texts)
		var invalid *InvalidError
		if !errors.As(err, &invalid) || !strings.Contains(err.Error(), x.want) {
			t.Errorf("%s %v：got %v, want %s", x.key, x.texts, err, x.want)
		}
	}
	if got, _ := Load(ctx, db); len(got.Replies) != 0 {
		t.Errorf("沒過的不該存進去：%v", got.Replies)
	}
}
