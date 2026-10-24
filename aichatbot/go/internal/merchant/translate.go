package merchant

// 固定回覆的翻譯：商家改慣用語言版本，存檔時用 AI 翻成其他啟用的語言。
// 顧客訊息翻成商家慣用語言（後台看的）用另一套提示詞（TranslateCustomer）

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"aichatbot/internal/model"
)

// languageNames 是語言代碼在翻譯指示裡的名稱（SupportedLanguages 每一種都要有）
var languageNames = map[string]string{
	"zh": "繁體中文", "en": "English", "ja": "日本語", "ko": "한국어", "th": "ภาษาไทย", "vi": "Tiếng Việt",
	"id": "Bahasa Indonesia", "ms": "Bahasa Melayu", "fr": "Français", "de": "Deutsch", "es": "Español",
}

// LanguageName 是語言代碼在給 AI 的指示裡的名稱（翻譯、一般對話寫回覆）
func LanguageName(code string) string {
	return languageNames[code]
}

// ErrTranslate：AI 沒翻出來、譯文含 |，或譯文裡的 {…} 變數和原文不一樣
var ErrTranslate = errors.New("翻譯失敗，請稍後再試")

// ErrTranslateBlocked：AI 沒翻出來是因為被安全過濾擋下（#68）；errors.Is 也認得是 ErrTranslate。
// 同一句每次都會被擋（greedy），訊息不叫人稍後再試
var ErrTranslateBlocked error = translateBlocked{}

type translateBlocked struct{}

func (translateBlocked) Error() string        { return "翻譯被安全過濾擋下，請改寫原文再試" }
func (translateBlocked) Is(target error) bool { return target == ErrTranslate }

// TranslateCustomer 把一則顧客傳給店家的訊息從 from 翻成 to（語言代碼），給後台看（對話紀錄的翻譯按鈕、待補問題的譯文）：
// 照原意翻，不像 Translate 那樣檢查 | 和 {…} 變數（那是給固定回覆的）。AI 沒翻出來回 ErrTranslate
func TranslateCustomer(ctx context.Context, chat model.ChatModel, from, to, text string) (string, error) {
	out, err := translate(ctx, chat, "TranslateCustomer", fmt.Sprintf("把這則顧客傳給店家的訊息從%s翻成%s，照原意翻：不要加語氣詞、客套話或表情符號，"+
		"原文沒有的不要加，原文有的表情符號照留。只輸出譯文。", languageNames[from], languageNames[to]), to, text)
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", ErrTranslate
	}
	return out, nil
}

// Translate 把一句店家客服的回覆從 from 翻成 to（語言代碼）
func Translate(ctx context.Context, chat model.ChatModel, from, to, text string) (string, error) {
	out, err := translate(ctx, chat, "Translate", fmt.Sprintf("把這句店家客服的回覆從%s翻成%s。語氣和原文一致（親切、口語），表情符號照留；"+
		"{name} 這類大括號變數原樣保留、不要翻譯。只輸出譯文。", languageNames[from], languageNames[to]), to, text)
	if err != nil {
		return "", err
	}
	// | 是存的時候分隔多句用的，譯文裡有就存不進去
	if out == "" || strings.Contains(out, "|") || !slices.Equal(variables(out), variables(text)) {
		return "", ErrTranslate
	}
	return out, nil
}

// 模型服務看題目 name（Translate、TranslateCustomer）決定用寬鬆的安全過濾（Swift AppleModel.permissiveRequests），改名要兩邊一起改。
// translate 用題目 name 和指示把 text 翻成 to，回去掉前後空白的譯文；AI 沒翻出來是空字串，被安全過濾擋下時另外回 ErrTranslateBlocked
func translate(ctx context.Context, chat model.ChatModel, name, instructions, to, text string) (string, error) {
	reply, usage := chat.Respond(ctx, model.Request{
		Name:         name,
		Text:         text,
		Instructions: instructions,
		Fields:       []model.Field{{Kind: model.KindText, Name: "translation", Description: languageNames[to] + "的譯文"}},
	})
	if reply == nil {
		if usage.Failure == model.FailureBlockedBySafety {
			return "", ErrTranslateBlocked
		}
		return "", nil
	}
	out, _ := reply.String("translation")
	return strings.TrimSpace(out), nil
}

var variablePattern = regexp.MustCompile(`\{[^{}]*\}`)

// variables 是句子裡用到的 {…} 變數（排序、去掉重複）
func variables(s string) []string {
	vs := variablePattern.FindAllString(s, -1)
	slices.Sort(vs)
	return slices.Compact(vs)
}
