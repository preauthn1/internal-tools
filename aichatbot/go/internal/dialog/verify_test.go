package dialog

import (
	"strings"
	"testing"
)

// 原文片段驗證：沒過時回擋下的原因，通過是空字串
func TestVerifyAnswer(t *testing.T) {
	const hours = "11:00–20:00 營業，週一公休。"
	const discount = "店內自取滿 300 元打 9 折喔！"
	for _, c := range []struct {
		name string
		in   answerCheck
		want string
	}{
		{"假出處：原文片段不在出處裡", answerCheck{Source: hours, Quote: "週日公休", Reply: "我們週日公休喔", Lang: "zh"}, "quote_not_in_source"},
		{"多出數字", answerCheck{Source: "桂花烏龍 110 元", Quote: "桂花烏龍 110 元", Reply: "桂花烏龍一杯 120 元", Lang: "zh"}, "number_not_in_quote:120"},
		{"顧客訊息原本就有的數字不算", answerCheck{Source: "店內最多 8 位。", Quote: "店內最多 8 位", Reply: "4 位沒問題，店內最多 8 位", Message: "4 個人可以嗎", Lang: "zh"}, ""},
		{"千分位", answerCheck{Source: "包場 1,200 元起", Quote: "包場 1,200 元起", Reply: "包場是 1200 元起", Lang: "zh"}, ""},
		{"全形數字", answerCheck{Source: "桂花烏龍 １１０ 元", Quote: "桂花烏龍 １１０ 元", Reply: "Osmanthus Oolong is 110 NTD.", Lang: "en"}, ""},
		{"跨語言：十一點對 11:00", answerCheck{Source: "每天十一點開門", Quote: "十一點開門", Reply: "We open at 11:00.", Lang: "en"}, ""},
		{"跨語言：11:00 對 11 點", answerCheck{Source: hours, Quote: "11:00–20:00 營業", Reply: "我們 11 點開門", Lang: "zh"}, ""},
		{"24 小時制對晚上 8 點", answerCheck{Source: hours, Quote: "11:00–20:00 營業", Reply: "營業到晚上八點", Lang: "zh"}, ""},
		{"24 小時制對 8 pm", answerCheck{Source: hours, Quote: "11:00–20:00 營業", Reply: "We close at 8 pm.", Lang: "en"}, ""},
		{"中文數字編的時間", answerCheck{Source: hours, Quote: "11:00–20:00 營業", Reply: "我們十二點開門", Lang: "zh"}, "number_not_in_quote:12:00"},
		{"一般用語裡的一不算數字", answerCheck{Source: hours, Quote: "週一公休", Reply: "週一公休，一起期待您光臨", Lang: "zh"}, ""},
		{"一點、十分、萬一、千萬不是數字", answerCheck{Source: hours, Quote: "週一公休", Reply: "週一公休，十分抱歉，萬一要來請早一點，千萬別白跑", Lang: "zh"}, ""},
		{"日期不是電話", answerCheck{Source: hours, Quote: "週一公休", Reply: "10 月 9 日是週一，公休喔", Message: "2026-10-09 有開嗎", Lang: "zh"}, ""},
		{"24 小時制不能當成別的數字", answerCheck{Source: hours, Quote: "11:00–20:00 營業", Reply: "最多 8 人", Lang: "zh"}, "number_not_in_quote:8"},
		{"電話寫法不同", answerCheck{Source: "電話 02-1234-5678", Quote: "電話 02-1234-5678", Reply: "請撥 02 1234 5678", Lang: "zh"}, ""},
		{"編的電話", answerCheck{Source: "電話 02-1234-5678", Quote: "電話 02-1234-5678", Reply: "請撥 02-1234-5679", Lang: "zh"}, "number_not_in_quote:0212345679"},
		{"網址", answerCheck{Source: "菜單 https://example.invalid/menu", Quote: "菜單 https://example.invalid/menu", Reply: "菜單在 https://example.invalid/menu 喔", Lang: "zh"}, ""},
		{"編的網址", answerCheck{Source: "菜單 https://example.invalid/menu", Quote: "菜單 https://example.invalid/menu", Reply: "菜單在 https://example.invalid/order 喔", Lang: "zh"}, "number_not_in_quote:https://example.invalid/order"},
		{"空的原文片段", answerCheck{Source: hours, Quote: " ", Reply: "我們週一公休喔", Lang: "zh"}, "empty_quote"},
		{"回覆混進 JSON 格式", answerCheck{Source: hours, Quote: "週一公休", Reply: "我們週一公休喔！」}```json{", Lang: "zh"}, "malformed_reply"},
		{"空的回覆", answerCheck{Source: hours, Quote: "週一公休", Reply: "", Lang: "zh"}, "empty_reply"},
		{"語言和顧客不同", answerCheck{Source: hours, Quote: "週一公休", Reply: "我們週一公休喔", Lang: "en"}, "language_mismatch:zh"},
		{"英文回覆附中文品名不算換語言", answerCheck{Source: "桂花烏龍 110 元", Quote: "桂花烏龍 110 元", Reply: "The Osmanthus Oolong (桂花烏龍) is 110 NTD.", Lang: "en"}, ""},
		{"看不出語言的短回覆不擋", answerCheck{Source: "人數 8 位", Quote: "人數 8 位", Reply: "8 👍", Lang: "en"}, ""},
		{"太長", answerCheck{Source: hours, Quote: "週一公休", Reply: strings.Repeat("週一公休喔", 61), Lang: "zh"}, "too_long"},
		{"打 9 折譯成 9% off", answerCheck{Source: discount, Quote: discount, Reply: "Orders over 300 get 9% off!", Lang: "en"}, "number_not_in_quote:9"},
		{"打 9 折譯成 9割引", answerCheck{Source: discount, Quote: discount, Reply: "300元以上で9割引です！", Lang: "ja"}, "number_not_in_quote:9割"},
		{"打 9 折不放行單獨的 1", answerCheck{Source: discount, Quote: discount, Reply: "Orders of 300 get 1 free cup!", Lang: "en"}, "number_not_in_quote:1"},
		{"打 8.8 折對 12% off", answerCheck{Source: "自取打 8.8 折", Quote: "自取打 8.8 折", Reply: "Pickup gets 12% off.", Lang: "en"}, ""},
		{"打 9 折對 10% off", answerCheck{Source: discount, Quote: discount, Reply: "Orders of 300 or more get 10% off!", Lang: "en"}, ""},
		{"打 9 折對 1割引", answerCheck{Source: discount, Quote: discount, Reply: "300元以上で1割引です！", Lang: "ja"}, ""},
		{"打九折對 9 折", answerCheck{Source: "自取打九折", Quote: "自取打九折", Reply: "自取打 9 折喔", Lang: "zh"}, ""},
		{"通過", answerCheck{Source: hours, Quote: "週一公休", Reply: "我們週一公休喔", Lang: "zh"}, ""},
	} {
		if got := verifyAnswer(c.in); got != c.want {
			t.Errorf("%s：got %q, want %q", c.name, got, c.want)
		}
	}
}
