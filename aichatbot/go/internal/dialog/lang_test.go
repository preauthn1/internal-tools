package dialog

import "testing"

// 語言判斷：先看文字系統，拉丁字母看常用字；短句有英文常用字就當英文；表情、標點判斷不了
func TestDetectLanguage(t *testing.T) {
	for text, want := range map[string]string{
		"你們幾點開":                          "zh",
		"可以用LINE Pay嗎":                   "zh",
		"抹茶ラテはいくらですか":                    "ja",
		"田中です":                           "ja",
		"영업시간이 어떻게 돼요?":                  "ko",
		"Hi, I am Kevin":                 "en",
		"Sorry, what time do you close?": "en",
		"Kevin":                          "en",
		"Bill":                           "en",
		"Okay":                           "en",
		"Bonjour, avez-vous une table pour ce soir ?":   "fr",
		"Hola, quiero reservar una mesa para mañana":    "es",
		"Guten Tag, ich möchte einen Tisch reservieren": "de",
		"Saya mau pesan meja untuk besok":               "id",
		"Xin chào, cho tôi đặt bàn":                     "vi",
		"สวัสดีครับ":                                    "th",
		"😂😂":                                            "",
		"！？":                                            "",
		"4":                                             "",
		"0912345678":                                    "",
	} {
		if got := DetectLanguage(text); got != want {
			t.Errorf("%s：%q，應為 %q", text, got, want)
		}
	}
}
