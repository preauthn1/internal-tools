package dialog

import (
	"context"
	"reflect"
	"strings"
	"testing"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

func giftBox() merchant.CustomIntent {
	return merchant.CustomIntent{
		Code: "custom_1", Name: "禮盒", Description: "詢問禮盒或送禮", Examples: []string{"有賣禮盒嗎", "想買來送人", "過年送禮推薦"},
		Action: merchant.ActionFixedText, Text: "禮盒請洽門市，電話 02-1234-5678。",
	}
}

// customMerchant 是有一個自訂意圖的商家
func customMerchant(ci merchant.CustomIntent) merchant.Config {
	c := knowledgeMerchant()
	c.Knowledge = append(c.Knowledge, merchant.Section{Title: "包場", Body: "包場請提前一週預約，最多 20 位。"})
	c.CustomIntents = []merchant.CustomIntent{ci}
	return c
}

// 沒有自訂意圖：送給 AI 的分類欄位和原本一模一樣（題庫分數不受影響）；
// 有自訂意圖：放在 general 前面，說明附上例句
func TestClassifyOptionsWithCustomIntents(t *testing.T) {
	e := newEnv(t, testMerchant(), nil)
	e.say("你好")
	if got := e.model.Requests()[0].Fields; !reflect.DeepEqual(got, intentFields) {
		t.Errorf("沒有自訂意圖時欄位變了：%+v", got)
	}

	e = newEnv(t, customMerchant(giftBox()), nil)
	e.say("你好")
	var choice model.Field
	for _, f := range e.model.Requests()[0].Fields {
		if f.Name == "intent" {
			choice = f
		}
	}
	want := "booking,cancel,follow_up,claim,identity,self_intro,human_request,complaint,off_topic,unclear,custom_1,general"
	if strings.Join(choice.Options, ",") != want {
		t.Errorf("選項：%v", choice.Options)
	}
	if !strings.Contains(choice.Description, "\ncustom_1=詢問禮盒或送禮。例：有賣禮盒嗎、想買來送人、過年送禮推薦\ngeneral=") {
		t.Errorf("說明：%s", choice.Description)
	}
}

// 回固定文字：顧客用商家慣用語言就回原文，不再呼叫 AI；traces 記下 custom 步驟
func TestCustomFixedText(t *testing.T) {
	e := newEnv(t, customMerchant(giftBox()), map[string][]map[string]model.Value{"Intent": {intent("custom_1")}})
	r := e.say("中秋有禮盒嗎")
	if r.Key != "custom:custom_1" || r.Text != "禮盒請洽門市，電話 02-1234-5678。" || len(e.model.Requests()) != 1 {
		t.Errorf("got %+v", r)
	}
	if got := steps(r); got != "intent:custom_1,custom:fixed_text" {
		t.Errorf("traces：%s", got)
	}
}

// 回固定文字、顧客用其他語言：翻成顧客的語言，譯文通過驗證就回譯文；沒通過（多出原文沒有的數字）就回原文
func TestCustomFixedTextTranslated(t *testing.T) {
	for _, c := range []struct{ translation, want, trace string }{
		{"Please contact the store for gift boxes, tel 02-1234-5678.", "Please contact the store for gift boxes, tel 02-1234-5678.", "translate:ok"},
		{"Gift boxes are 500, call 02-1234-5678.", "禮盒請洽門市，電話 02-1234-5678。", "translate:number_not_in_quote:500"},
	} {
		e := newEnv(t, customMerchant(giftBox()), map[string][]map[string]model.Value{
			"Intent":    {intent("custom_1")},
			"Translate": {{"translation": model.Text(c.translation)}},
		})
		r := e.say("Do you sell gift boxes?")
		if r.Key != "custom:custom_1" || r.Text != c.want {
			t.Errorf("got %+v", r)
		}
		if got := steps(r); got != "intent:custom_1,custom:fixed_text,"+c.trace {
			t.Errorf("traces：%s", got)
		}
	}
}

func refund() merchant.CustomIntent {
	return merchant.CustomIntent{
		Code: "custom_2", Name: "退貨", Description: "要退貨或退款", Examples: []string{"我要退貨", "可以退款嗎", "買錯了想退"},
		Action: merchant.ActionHandoff, Text: "退貨由店員為您處理，請稍候。",
	}
}

// 轉人工：開轉人工的事項，回商家填的提示（沒填用轉人工的預設提示）；轉人工關閉時照原本的處理（回固定回覆 handoff_off，不轉）
func TestCustomHandoff(t *testing.T) {
	e := newEnv(t, customMerchant(refund()), map[string][]map[string]model.Value{"Intent": {intent("custom_2")}})
	r := e.say("杯子破了要退")
	if r.Key != "custom:custom_2" || r.Text != "退貨由店員為您處理，請稍候。" || !e.customer().Human {
		t.Errorf("got %+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 轉人工：自訂意圖「退貨」：杯子破了要退" {
		t.Errorf("待處理：%q", p)
	}

	noPrompt := refund()
	noPrompt.Text = ""
	e = newEnv(t, customMerchant(noPrompt), map[string][]map[string]model.Value{"Intent": {intent("custom_2")}})
	if r := e.say("杯子破了要退"); r.Text != "我先幫您轉給店員，請稍候喔～" || !e.customer().Human {
		t.Errorf("沒填提示：%+v", r)
	}

	e = newEnv(t, setFeature(customMerchant(refund()), "handoff", false), map[string][]map[string]model.Value{"Intent": {intent("custom_2")}})
	if r := e.say("杯子破了要退"); r.Key != "handoff_off" || !strings.HasSuffix(r.Text, "：電話 02-1234-5678") || e.customer().Human || len(e.pending()) != 0 {
		t.Errorf("轉人工關閉：%+v %q", r, e.pending())
	}

	// 沒綁員工群組（#55）：同轉人工關閉，但照樣開事項給後台看
	e = newEnv(t, customMerchant(refund()), map[string][]map[string]model.Value{"Intent": {intent("custom_2")}})
	e.unbindStaffChat()
	if r := e.say("杯子破了要退"); r.Key != "handoff_off" || !strings.HasSuffix(r.Text, "：電話 02-1234-5678") || e.customer().Human {
		t.Errorf("沒綁群組：%+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 想找真人（沒有員工群組，未轉人工）：自訂意圖「退貨」：杯子破了要退" {
		t.Errorf("沒綁群組的待處理：%q", p)
	}
}

func party() merchant.CustomIntent {
	return merchant.CustomIntent{
		Code: "custom_3", Name: "包場", Description: "詢問包場", Examples: []string{"可以包場嗎", "想包下整間店", "包場怎麼收費"},
		Action: merchant.ActionKnowledge, Knowledge: []string{"包場"},
	}
}

// 只根據指定知識回答：不挑 FAQ，挑知識的清單只有指定的段落，寫回覆後驗證原文片段
func TestCustomKnowledgeAnswers(t *testing.T) {
	e := newEnv(t, customMerchant(party()), map[string][]map[string]model.Value{
		"Intent": {intent("custom_3")},
		"Pick":   {pick("包場", "K1")},
		"Answer": {answer("K1", "最多 20 位", "yes", "包場最多 20 位喔")},
	})
	r := e.say("包場可以坐幾個人")
	if r.Key != "custom:custom_3" || r.Text != "包場最多 20 位喔" {
		t.Fatalf("got %+v", r)
	}
	reqs := e.model.Requests()
	if len(reqs) != 3 || reqs[1].Name != "Pick" {
		t.Fatalf("requests：%+v", reqs)
	}
	if !strings.Contains(reqs[1].Instructions, "K1 包場") || strings.Contains(reqs[1].Instructions, "營業時間") || strings.Contains(reqs[1].Instructions, "菜單") {
		t.Errorf("挑知識的清單：%s", reqs[1].Instructions)
	}
	if got := steps(r); got != "intent:custom_3,custom:knowledge,pick:K1,answer:yes,verify:ok" {
		t.Errorf("traces：%s", got)
	}
}

// 只根據指定知識回答、挑不到或驗證沒過：走 noAnswer（找不到答案時轉人工開著就轉人工，關著回看不懂）
func TestCustomKnowledgeNoAnswer(t *testing.T) {
	e := newEnv(t, customMerchant(party()), map[string][]map[string]model.Value{
		"Intent": {intent("custom_3")},
		"Pick":   {pick("")},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = true })
	if r := e.say("包場有送蛋糕嗎"); r.Key != "handoff" || e.pending()[0] != "handoff：🙋 轉人工：找不到相關資料：包場有送蛋糕嗎" {
		t.Errorf("挑不到：%+v %q", r, e.pending())
	}
	e = newEnv(t, customMerchant(party()), map[string][]map[string]model.Value{
		"Intent": {intent("custom_3")},
		"Pick":   {pick("", "K1")},
		"Answer": {answer("K1", "最多 20 位", "yes", "包場最多 30 位")},
	})
	e.setHandoff(func(s *store.HandoffSetting) { s.OnNoAnswer = false })
	if r := e.say("包場可以坐幾個人"); r.Key != "unclear" || e.customer().Human {
		t.Errorf("驗證沒過：%+v", r)
	}
}

// 自訂意圖的例句也算意圖案例：完全相同的句子直接沿用，相似的附在分類的指示裡
func TestCustomExamplesJoinExampleLookup(t *testing.T) {
	c := customMerchant(giftBox())
	c.Examples = []merchant.Example{{Text: "我要訂位", Intent: "booking", Builtin: true}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general")}})
	if got, _ := e.bot.Classify(context.Background(), "有賣禮盒嗎！", Customer{}); got != "custom_1" || len(e.model.Requests()) != 0 {
		t.Errorf("完全相同：%q，呼叫 AI %d 次", got, len(e.model.Requests()))
	}
	e.bot.Embedder, e.bot.K = &fakeEmbedder{vectors: map[string][]float32{"禮物": {1, 0, 0}, "想買來送人": {1, 0, 0}}}, 1
	e.bot.Classify(context.Background(), "禮物", Customer{})
	if ins := e.model.Requests()[0].Instructions; !strings.Contains(ins, "「想買來送人」 → custom_1") {
		t.Errorf("相似案例：%s", ins)
	}
}
