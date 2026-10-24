package dialog

import (
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"
)

// MARK: 名字（由分類那次呼叫一起抽出，這裡只做驗證）

// verifiedName：名字必須真的出現在訊息裡（防止 AI 編名字）；只有一個漢字（只說姓）時補成「X先生／小姐」
func verifiedName(raw, text string) string {
	name := strings.TrimSpace(raw)
	core := nameCore(name)
	if core == "" || !containsFold(text, core) || isRoleWord(core) {
		return ""
	}
	if r, _ := utf8.DecodeRuneInString(name); utf8.RuneCountInString(name) == 1 && unicode.Is(unicode.Han, r) {
		return name + "先生／小姐"
	}
	return name
}

// nameCore：去掉稱謂的名字本體
func nameCore(name string) string {
	for _, title := range []string{"先生／小姐", "先生", "小姐", "女士", "太太", "様", "さん", "Mr. ", "Ms. "} {
		name = strings.ReplaceAll(name, title, "")
	}
	return strings.TrimSpace(name)
}

// roleWords：角色稱呼不是名字（「我是顧客」「I'm the customer」；「我是會員」「我是學生」，#61）
var roleWords = []string{"顧客", "客人", "客戶", "客服", "老闆", "老闆娘", "店員", "店長", "服務生", "員工",
	"會員", "學生", "新客", "熟客", "常客", "老客人", "男生", "女生", "外國人",
	"customer", "client", "staff", "owner", "manager", "boss", "お客様", "お客さん", "スタッフ"}

// isRoleWord：整個名字（去掉稱謂、不分大小寫）就是角色稱呼
func isRoleWord(name string) bool {
	core := strings.ToLower(nameCore(name))
	for _, w := range []string{"the ", "a ", "an "} {
		core = strings.TrimPrefix(core, w)
	}
	return slices.Contains(roleWords, core)
}

// 自我介紹的句型：名字前面接這些（我是小康、My name is Kevin、私は田中），或後面接這些（田中と申します、田中です）。
// 「這裡是」「this is」多半在問店家（「這裡是大安區那間嗎」），不算；「ですか」是問句（「ここは東京駅ですか」），也不算
var (
	introPrefixes = []string{"我是", "我叫", "叫我", "喊我", "叫做", "名叫", "敝姓", "我姓", "本人", "名字是", "名字叫",
		"my name is", "my name's", "i am", "i'm", "i’m", "call me", "私は", "わたしは", "僕は", "俺は", "名前は"}
	introSuffixes = []string{"と申します", "といいます", "と言います", "です", "だよ"}
)

// introducedIn：名字出現在自我介紹的句型裡（只要有一處就算）。
// 句子開頭的地名、店名、商品名、問句裡的名詞不算（「台灣最高的山是哪座」曾被抽成名字「台灣」）
func introducedIn(name, text string) bool {
	core := strings.ToLower(nameCore(name))
	lower := strings.ToLower(text)
	if core == "" {
		return false
	}
	trim := func(r rune) bool { return unicode.IsSpace(r) || unicode.IsPunct(r) }
	for i := 0; ; {
		j := strings.Index(lower[i:], core)
		if j < 0 {
			return false
		}
		start, end := i+j, i+j+len(core)
		before := strings.TrimRightFunc(lower[:start], trim)
		after := strings.TrimLeftFunc(lower[end:], unicode.IsSpace)
		for _, p := range introPrefixes {
			if strings.HasSuffix(before, p) {
				return true
			}
		}
		for _, s := range introSuffixes {
			if rest, ok := strings.CutPrefix(after, s); ok && !strings.HasPrefix(rest, "か") {
				return true
			}
		}
		i = end
	}
}

// 整句就是自我介紹的句型（#61）：「叫我」「我叫」這類是顧客明確要被這樣叫，單字名不補「先生／小姐」。
// 長的放前面（「我叫做」在「我叫」前、「我的名字是」在「我是」前）；新增句型也要加進 introPrefixes
var (
	introPatterns = []struct {
		prefix string
		callMe bool
	}{
		{"我的名字是", false}, {"我的名字叫", true}, {"我叫做", true}, {"我叫", true}, {"叫我", true}, {"喊我", true}, {"我是", false},
	}
	introTails = []string{"就可以了", "就可以", "就好了", "就好", "就行", "即可"}
	// 名字裡不會有的字：疑問詞（我是誰、我叫什麼）、代名詞、虛詞和動詞（我是來訂位的、叫我去找店員）
	notNameRunes = "誰啥哪嗎呢幾麼你我他她您的了來要想在有不問說是去找幫等"
)

// introName：整句（去掉頭尾的標點、語助詞）就是「我是X」「叫我X」「我叫做X就好」這類自我介紹時回傳名字，否則是空字串。
// X 是 1～4 個漢字；句中出現的不算（「我朋友叫小美」不是顧客的名字）
func introName(text string) string {
	core := strings.TrimFunc(text, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) || strings.ContainsRune("啦喔哦呀啊～~", r)
	})
	for _, p := range introPatterns {
		rest, ok := strings.CutPrefix(core, p.prefix)
		if !ok {
			continue
		}
		for _, tail := range introTails {
			rest = strings.TrimSuffix(rest, tail)
		}
		n := utf8.RuneCountInString(rest)
		if n < 1 || n > 4 || strings.ContainsFunc(rest, func(r rune) bool { return !unicode.Is(unicode.Han, r) || strings.ContainsRune(notNameRunes, r) }) {
			return ""
		}
		if p.callMe && !isRoleWord(rest) {
			return rest
		}
		return verifiedName(rest, text) // 單字名補「先生／小姐」、角色稱呼不收
	}
	return ""
}

// containsFold：不分大小寫的包含
func containsFold(s, sub string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(sub))
}
