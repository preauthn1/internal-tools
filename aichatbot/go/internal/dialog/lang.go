package dialog

import (
	"strings"
	"unicode"
)

// MARK: 語言（程式判斷，不用 AI）
// Swift 版用 Apple 的 NLLanguageRecognizer；Go 版沒有，改成先看文字系統，拉丁字母再看常用字

// DetectLanguage 回傳 ISO 639-1 代碼；判斷不了（表情、標點、太短）回空字串，沿用上一句的語言。
// 後台判斷對話紀錄要不要翻譯也用它
func DetectLanguage(text string) string {
	// 先看文字系統：假名→日文、韓文字母→韓文、漢字→中文（「可以用LINE Pay嗎」不會被判成英文）
	has := func(tables ...*unicode.RangeTable) bool {
		return strings.ContainsFunc(text, func(r rune) bool { return unicode.In(r, tables...) })
	}
	switch {
	case has(unicode.Hiragana, unicode.Katakana):
		return "ja"
	case has(unicode.Hangul):
		return "ko"
	case has(unicode.Han):
		return "zh"
	case has(unicode.Thai):
		return "th"
	case has(unicode.Cyrillic):
		return "ru"
	case has(unicode.Arabic):
		return "ar"
	}
	letters := 0
	for _, r := range text {
		if unicode.IsLetter(r) {
			letters++
		}
	}
	if letters < 2 {
		return ""
	}
	lower := strings.ToLower(text)
	// 越南文特有的字母（ă、đ、ơ、ư 和帶聲調的母音）
	if strings.ContainsFunc(lower, func(r rune) bool {
		return (r >= 0x1EA0 && r <= 0x1EF9) || strings.ContainsRune("ăđơư", r)
	}) {
		return "vi"
	}
	best, bestScore, en := "en", 0, 0
	words := strings.FieldsFunc(lower, func(r rune) bool { return !unicode.IsLetter(r) && r != '\'' })
	for _, lang := range latinOrder {
		score := 0
		for _, w := range words {
			if latinWords[lang][w] {
				score++
			}
		}
		if lang == "en" {
			en = score
		}
		if score > bestScore {
			best, bestScore = lang, score
		}
	}
	// 短句容易誤判（Swift 版「Hi, I am Kevin」曾被判成印尼文）：有英文常用字就當英文
	if best != "en" && len([]rune(text)) < 40 && en > 0 {
		return "en"
	}
	if bestScore > 0 {
		return best
	}
	// 沒有常用字：看特殊字母，都沒有就當英文
	for lang, chars := range latinLetters {
		if strings.ContainsAny(lower, chars) {
			return lang
		}
	}
	return "en"
}

var latinOrder = []string{"en", "es", "fr", "de", "pt", "it", "id", "ms"}

func wordSet(s string) map[string]bool {
	m := map[string]bool{}
	for _, w := range strings.Fields(s) {
		m[w] = true
	}
	return m
}

// 各語言最常用的字（判斷拉丁字母的語言用）
var latinWords = map[string]map[string]bool{
	"en": wordSet("the a an i you we he she it they is are am was be do does did can could will would what when where how why who which this that my your our to of in on at for with and or not no yes please thanks thank hi hello hey bye book table time open close have has me want like there any"),
	"es": wordSet("el la los las un una es son está estoy qué cuándo dónde cómo por para con y o no sí gracias hola quiero mesa tiene hay de del que"),
	"fr": wordSet("le la les un une est sont je tu vous nous il elle quoi quand où comment pour avec et ou non oui merci bonjour veux table avez du des que c'est"),
	"de": wordSet("der die das ein eine ist sind ich du sie wir was wann wo wie warum für mit und oder nicht nein ja danke hallo möchte tisch haben gibt es"),
	"pt": wordSet("o os um uma é são estou você nós quando onde como por para com e ou não sim obrigado obrigada olá quero mesa tem há"),
	"it": wordSet("il lo gli un una è sono io tu voi noi cosa quando dove come perché per con e o non sì grazie ciao voglio tavolo avete c'è"),
	"id": wordSet("saya anda kamu kami apa kapan di mana bagaimana kenapa untuk dengan dan atau tidak ya terima kasih halo mau meja ada buka tutup jam berapa ini itu"),
	"ms": wordSet("saya awak anda kami apa bila di mana bagaimana kenapa untuk dengan dan atau tidak ya terima kasih helo mahu meja ada buka tutup pukul berapa ini itu"),
}

// 沒有常用字時看特殊字母
var latinLetters = map[string]string{"es": "ñ¿¡", "de": "ßäöü", "fr": "çèêëœ", "pt": "ãõ"}

// languageNames 是各語言自己的名稱（「目前只支援…」用）
var languageNames = map[string]string{
	"zh": "中文", "en": "English", "ja": "日本語", "ko": "한국어", "th": "ไทย", "vi": "Tiếng Việt",
	"id": "Bahasa Indonesia", "ms": "Bahasa Melayu", "fr": "Français", "de": "Deutsch", "es": "Español",
}
