package merchant

import (
	"reflect"
	"strings"
	"testing"
)

// 知識段落（#65）：寫回覆一次最多送 3 段全文。單段太長、最長 3 段合計太長各一則警告，附上是哪幾段；沒超過時是空陣列
func TestKnowledgeWarnings(t *testing.T) {
	s := func(title string, n int) Section { return Section{Title: title, Body: strings.Repeat("字", n)} }
	sections := []Section{s("短", 10), s("菜單", 3100), s("交通", 2000), s("包場", 1500), s("歷史", 2500)}
	// 字數算標題加內文：菜單 3102、歷史 2502、交通 2002
	want := []Warning{
		{Kind: "knowledge_section_long", Limit: KnowledgeSectionLimit, Targets: []string{"菜單"}},
		{Kind: "knowledge_too_long", Chars: 3102 + 2502 + 2002, Limit: KnowledgeCharsLimit, Targets: []string{"菜單", "歷史", "交通"}},
	}
	if got := KnowledgeWarnings(sections); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v", got)
	}
	if got := KnowledgeWarnings(sections[2:4]); len(got) != 0 || got == nil {
		t.Errorf("沒超過：got %#v", got)
	}
}

// FAQ 的選項（問法加答案）合計太長（#65）：警告附上字數、上限和最長的幾題；沒超過時是空陣列
func TestFAQWarningsTooLong(t *testing.T) {
	long := strings.Repeat("長", 2100)
	faq := []FAQ{{Key: "k1", Question: "q1", Answer: long}, {Key: "k2", Question: "q2", Answer: "短"}, {Key: "k3", Question: "q3", Answer: long}, {Key: "k4", Question: "q4", Answer: long}}
	// 每題長的選項是「k1：q1（答案：…）」2110 字，短的「k2：q2（答案：短）」11 字
	want := []Warning{{Kind: "faq_too_long", Chars: 2110*3 + 11, Limit: FAQCharsLimit, Targets: []string{"q1", "q3", "q4"}}}
	if got := FAQWarnings(faq); !reflect.DeepEqual(got, want) {
		t.Errorf("got %+v", got)
	}
	if got := FAQWarnings(faq[:2]); len(got) != 0 || got == nil {
		t.Errorf("沒超過：got %#v", got)
	}
}
