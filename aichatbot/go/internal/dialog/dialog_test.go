package dialog

import (
	"context"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
	"aichatbot/internal/store"
)

// 今天 2026-10-07（週三）
var wednesday = time.Date(2026, 10, 7, 12, 0, 0, 0, time.Local)

// testMerchant 是測試用的商家設定：預約、轉人工都開，顧客語言中英日
func testMerchant() merchant.Config {
	c := merchant.Default()
	c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = "森日茶屋", "賣茶的店", "電話 02-1234-5678"
	return c
}

type env struct {
	t     *testing.T
	db    *store.DB
	model *model.Scripted
	bot   *Bot
}

// newEnv：temp 資料庫存好商家設定、綁好員工群組（沒綁群組不會轉人工，#55），模型照腳本回答
func newEnv(t *testing.T, c merchant.Config, script map[string][]map[string]model.Value) *env {
	t.Helper()
	ctx := context.Background()
	db, err := store.Open(filepath.Join(t.TempDir(), "bot.sqlite"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	if _, err := merchant.Save(ctx, db.SQL(), c); err != nil {
		t.Fatal(err)
	}
	s, err := db.HandoffSetting(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ok, err := db.BindStaffChat(ctx, s.BindCode, -100); !ok || err != nil {
		t.Fatalf("綁員工群組：%v %v", ok, err)
	}
	m := model.NewScripted(script)
	return &env{t: t, db: db, model: m, bot: &Bot{DB: db, Model: m, Now: func() time.Time { return wednesday }}}
}

// say：顧客 telegram:U1 說一句話
func (e *env) say(text string) Reply {
	e.t.Helper()
	r, err := e.bot.Receive(context.Background(), Message{Platform: "telegram", UserID: "U1", ChatID: "C1", Text: text})
	if err != nil {
		e.t.Fatal(err)
	}
	return r
}

func intent(code string) map[string]model.Value {
	return map[string]model.Value{"intent": model.Text(code), "self_name": model.Text("")}
}

// 一次選擇題：11 個內建意圖，同一次呼叫抽名字
func TestClassifyIsOneChoiceOfElevenIntents(t *testing.T) {
	e := newEnv(t, testMerchant(), nil)
	e.say("你好")
	reqs := e.model.Requests()
	if len(reqs) != 1 || reqs[0].Name != "Intent" || reqs[0].Text != "你好" {
		t.Fatalf("requests：%+v", reqs)
	}
	var names []string
	var options []string
	for _, f := range reqs[0].Fields {
		names = append(names, f.Name)
		if f.Name == "intent" {
			options = f.Options
		}
	}
	if strings.Join(names, ",") != "want,intent,self_name" {
		t.Errorf("欄位：%v", names)
	}
	want := "booking,cancel,follow_up,claim,identity,self_intro,human_request,complaint,off_topic,unclear,general"
	if strings.Join(options, ",") != want {
		t.Errorf("選項：%v", options)
	}
}

// AI 失敗時當成看不懂：請顧客說清楚，比轉給員工安全
func TestAIFailureIsUnclear(t *testing.T) {
	e := newEnv(t, testMerchant(), nil)
	if r := e.say("你們營業到幾點呀"); r.Key != "unclear" || r.Text != "不太確定您的意思，可以再說詳細一點嗎？" {
		t.Errorf("got %+v", r)
	}
}

func named(code, name string) map[string]model.Value {
	return map[string]model.Value{"intent": model.Text(code), "self_name": model.Text(name)}
}

// 自我介紹記住名字（存在資料庫的顧客狀態），之後問「我是誰」回得出來
func TestSelfIntroRemembersName(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named("self_intro", "小康"), intent("identity")}})
	if r := e.say("我是小康"); r.Key != "self_intro" || !strings.Contains(r.Text, "小康您好") {
		t.Errorf("自我介紹：%+v", r)
	}
	if r := e.say("我是誰"); r.Key != "identity" || !strings.HasPrefix(r.Text, "您是小康～") || !strings.Contains(r.Text, "森日茶屋") {
		t.Errorf("問身分：%+v", r)
	}
}

// 名字的防呆：問身分、追問、聲稱、無關、抱怨、想找真人的句子不是自我介紹；名字不能就是整句話；
// 單一個字只在自我介紹時採用；AI 編的名字（不在原文裡）不收
func TestNameGuards(t *testing.T) {
	for _, c := range []struct{ text, intent, name string }{
		{"我是誰", "identity", "我是誰"},
		{"小美呢", "follow_up", "小美"},
		{"店員小美說打五折", "claim", "小美"},
		{"幫我查周杰倫的歌", "off_topic", "周杰倫"},
		{"你很笨欸阿明", "complaint", "阿明"},
		{"我要找真人小美", "human_request", "小美"},
		{"中", "general", "中"},
		{"我是誰", "general", "我"},
		{"你好", "general", "小康"},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named(c.intent, c.name)}})
		e.say(c.text)
		got, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "" {
			t.Errorf("%s（%s）：記成名字 %q", c.text, c.intent, got.Name)
		}
	}
	// 自我介紹只說姓：補成「X先生／小姐」
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named("self_intro", "王")}})
	if r := e.say("我姓王"); r.Key != "self_intro" || !strings.Contains(r.Text, "王先生／小姐") {
		t.Errorf("只說姓：%+v", r)
	}
}

// 名字要有出處：出現在自我介紹的句型裡。問句開頭的地名、商品名不是名字（「台灣最高的山是哪座」曾回「台灣您好」）
func TestNameNeedsIntroPhrase(t *testing.T) {
	for _, c := range []struct{ text, intent, name string }{
		{"台灣最高的山是哪座", "general", "台灣"},
		{"台灣最高的山是哪座", "self_intro", "台灣"},
		{"森日抹茶好喝嗎", "general", "森日"},
		{"這裡是大安區那間嗎", "general", "大安區"},
		{"ここは東京駅ですか", "general", "東京駅"},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named(c.intent, c.name)}})
		if r := e.say(c.text); r.Key == "self_intro" || strings.Contains(r.Text, c.name+"您好") {
			t.Errorf("%s（%s）：回了自我介紹 %+v", c.text, c.intent, r)
		}
		got, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "" {
			t.Errorf("%s（%s）：記成名字 %q", c.text, c.intent, got.Name)
		}
	}
}

// 正常的自我介紹（中、英、日）照樣記得名字
func TestIntroPhrasesRememberName(t *testing.T) {
	for _, c := range []struct{ text, intent, name string }{
		{"嗨我叫阿哲", "self_intro", "阿哲"},
		{"敝姓林，想問一下", "self_intro", "林"},
		{"My name is Amy", "general", "Amy"},
		{"call me Ben", "self_intro", "Ben"},
		{"田中と申します", "self_intro", "田中"},
		{"佐藤です、よろしく", "general", "佐藤"},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named(c.intent, c.name)}})
		if r := e.say(c.text); r.Key != "self_intro" {
			t.Errorf("%s（%s）：%+v", c.text, c.intent, r)
		}
		got, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.HasPrefix(got.Name, c.name) {
			t.Errorf("%s（%s）：記成名字 %q", c.text, c.intent, got.Name)
		}
	}
}

// #61 重放：AI 把「我是親」「叫我親」「我叫做親」分成問身分，曾一直回「還不知道怎麼稱呼您呢」。
// 句子就是自我介紹句型時改成自我介紹；「叫我」「我叫」的單字名照原樣叫（顧客說「叫我親」就是要被這樣叫），
// 「我是X」照舊補成「X先生／小姐」
func TestIdentityLikeIntroIsSelfIntro(t *testing.T) {
	for _, c := range []struct{ text, want string }{
		{"我是親", "親先生／小姐您好"},
		{"叫我親", "親您好"},
		{"我叫做親", "親您好"},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity")}})
		if r := e.say(c.text); r.Key != "self_intro" || !strings.HasPrefix(r.Text, c.want) {
			t.Errorf("%s：%+v", c.text, r)
		}
	}
}

// 同一位顧客只追問一次稱呼（#61）：第二次問身分時只回客服的身分
func TestIdentityAsksNameOnce(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity"), intent("identity")}})
	if r := e.say("你是誰"); r.Key != "identity" || !strings.Contains(r.Text, "怎麼稱呼您") {
		t.Errorf("第一次：%+v", r)
	}
	if r := e.say("你是機器人嗎"); r.Key != "identity" || strings.Contains(r.Text, "怎麼稱呼您") || !strings.Contains(r.Text, "森日茶屋") {
		t.Errorf("第二次：%+v", r)
	}
}

// 剛問過稱呼、還不知道名字時，自我介紹句型的短句直接當成報名字，不問 AI（#61）
func TestAskedNameIntroSkipsAI(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity"), intent("unclear")}})
	e.say("你是哪位")
	if r := e.say("我叫做親啦"); r.Key != "self_intro" || !strings.HasPrefix(r.Text, "親您好") {
		t.Errorf("報名字：%+v", r)
	}
	if n := len(e.model.Requests()); n != 1 {
		t.Errorf("問了 AI %d 次", n)
	}
	e = newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity"), intent("unclear")}})
	e.say("你是哪位")
	if r := e.say("喊我萱就行"); r.Key != "self_intro" || !strings.HasPrefix(r.Text, "萱您好") {
		t.Errorf("喊我：%+v", r)
	}
}

// 只有問稱呼的下一句才不問 AI；中間隔了別的話，照常問 AI 分類（#61）
func TestAskedNameOnlyNextTurn(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity"), intent("complaint"), intent("unclear")}})
	e.say("你是哪位")
	e.say("你很笨")
	if r := e.say("我叫外送"); r.Key == "self_intro" {
		t.Errorf("隔一句還當成報名字：%+v", r)
	}
	if n := len(e.model.Requests()); n != 3 {
		t.Errorf("問了 AI %d 次", n)
	}
}

// 剛問過稱呼時，問句、提到別人、不是名字的句子照常問 AI，不記成名字（#61）
func TestAskedNameNotIntro(t *testing.T) {
	for _, c := range []struct{ text, intent string }{
		{"我是誰", "identity"},
		{"我叫什麼", "identity"},
		{"我是顧客", "general"},
		{"我是來訂位的", "booking"},
		{"叫我去找店員喔", "complaint"},
		{"我朋友叫小美", "general"},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("identity"), intent(c.intent)}})
		e.say("你是哪位")
		e.say(c.text)
		if reqs := e.model.Requests(); len(reqs) < 2 || reqs[1].Name != "Intent" || reqs[1].Text != c.text {
			t.Errorf("%s：沒有問 AI 分類 %+v", c.text, reqs)
		}
		got, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "" {
			t.Errorf("%s：記成名字 %q", c.text, got.Name)
		}
	}
}

// 角色稱呼（顧客、客服、店員、customer…）不是名字
func TestRoleWordsNotNames(t *testing.T) {
	for _, c := range []struct{ text, intent, name string }{
		{"我是顧客", "self_intro", "顧客"},
		{"I'm a customer", "general", "customer"},
		{"我是店員", "self_intro", "店員"},
		{"我是老闆娘", "self_intro", "老闆娘"},
		{"我是會員", "identity", ""}, // #61：AI 分成問身分、程式從句型抽名字時也不收
		{"我是學生", "identity", ""},
		{"我是新客", "identity", ""},
	} {
		e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named(c.intent, c.name)}})
		if r := e.say(c.text); r.Key == "self_intro" {
			t.Errorf("%s（%s）：回了自我介紹 %+v", c.text, c.intent, r)
		}
		got, err := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
		if err != nil {
			t.Fatal(err)
		}
		if got.Name != "" {
			t.Errorf("%s（%s）：記成名字 %q", c.text, c.intent, got.Name)
		}
	}
}

// 一般對話順便介紹了名字（「嗨，我是 Kevin」常被分成一般對話）：當成自我介紹；已經知道名字就照一般對話回
func TestGeneralWithNewNameIsSelfIntro(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {named("general", "Kevin"), named("general", "Bob")}})
	if r := e.say("Hi, I am Kevin"); r.Key != "self_intro" || !strings.Contains(r.Text, "Kevin") {
		t.Errorf("got %+v", r)
	}
	if r := e.say("Hi, I am Bob"); r.Key != "general" {
		t.Errorf("已經知道名字：%+v", r)
	}
}

// 只有表情、標點（沒有文字和數字）不用問 AI
func TestSymbolsOnlySkipsAI(t *testing.T) {
	e := newEnv(t, testMerchant(), nil)
	if r := e.say("😂😂！"); r.Key != "general" {
		t.Errorf("got %+v", r)
	}
	if n := len(e.model.Requests()); n != 0 {
		t.Errorf("呼叫了 AI %d 次", n)
	}
}

// pending：還沒處理的待處理事項文字
func (e *env) pending() []string {
	e.t.Helper()
	ts, err := e.db.OpenTickets(context.Background())
	if err != nil {
		e.t.Fatal(err)
	}
	var out []string
	for _, t := range ts {
		out = append(out, t.Kind+"："+t.Text)
	}
	return out
}

// 聲稱、無關累計 3 次轉人工；轉人工後 AI 不回覆、不呼叫 AI
func TestClaimAndOffTopicStreakHandsOff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("off_topic"), intent("claim"), intent("off_topic")}})
	for i, want := range []string{"off_topic", "claim", "handoff"} {
		if r := e.say("第幾句"); r.Key != want {
			t.Errorf("第 %d 句：%+v，應為 %s", i+1, r, want)
		}
	}
	if r := e.say("抹茶拿鐵多少錢"); !r.Silent {
		t.Errorf("轉人工後還回覆：%+v", r)
	}
	if n := len(e.model.Requests()); n != 3 {
		t.Errorf("呼叫 AI %d 次", n)
	}
	got := strings.Join(e.pending(), "\n")
	if !strings.Contains(got, "claim：❓ 顧客聲稱／更正：第幾句") || !strings.Contains(got, "handoff：🙋 轉人工：連續 3 次無關或聲稱的訊息") {
		t.Errorf("待處理：%s", got)
	}
}

// 中間有一般對話或預約就重新計算
func TestStreakResetsOnGeneral(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("off_topic"), intent("off_topic"), intent("general"), intent("off_topic")}})
	var keys []string
	for range 4 {
		keys = append(keys, e.say("一句話").Key)
	}
	if got := strings.Join(keys, ","); got != "off_topic,off_topic,general,off_topic" {
		t.Errorf("got %s", got)
	}
}

// 追問：沒有在等員工時是看不懂；在等員工時回請稍等，第 3 次起回催促，並開緊急通知引用原本那張
func TestFollowUpNudgesOnThird(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {
		intent("follow_up"), intent("claim"), intent("follow_up"), intent("follow_up"), intent("follow_up"), intent("follow_up"),
	}})
	var keys []string
	for _, text := range []string{"有消息了嗎", "上次說打五折", "有消息了嗎", "還要多久", "好了沒", "快點"} {
		keys = append(keys, e.say(text).Key)
	}
	if got := strings.Join(keys, ","); got != "unclear,claim,waiting,waiting,nudge,nudge" {
		t.Errorf("got %s", got)
	}
	p := e.pending()
	if len(p) != 2 || p[1] != "urgent：🔥 緊急：顧客已追問 3 次，請盡快回覆\n📌 原本的事項 #1：上次說打五折" {
		t.Errorf("待處理：%q", p)
	}
}

// setFeature 改功能開關（booking、handoff）
func setFeature(c merchant.Config, name string, enabled bool) merchant.Config {
	for i := range c.Features {
		if c.Features[i].Name == name {
			c.Features[i].Enabled = enabled
		}
	}
	return c
}

// 想找真人：轉人工開著就轉人工
func TestHumanRequestHandsOff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
	if r := e.say("我要找真人"); r.Key != "handoff" {
		t.Errorf("got %+v", r)
	}
	if r := e.say("在嗎"); !r.Silent {
		t.Errorf("轉人工後還回覆：%+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 轉人工：顧客要求找真人" {
		t.Errorf("待處理：%q", p)
	}
}

// 只說「轉」「轉人工」：內建案例庫（eval/examples.jsonl）標成想找真人，不問 AI 就轉人工（#63）
func TestShortTransferHandsOff(t *testing.T) {
	c, err := merchant.FromEvalDir("../../../eval")
	if err != nil {
		t.Fatal(err)
	}
	for _, text := range []string{"轉", "轉人工"} {
		e := newEnv(t, c, nil)
		if r := e.say(text); r.Key != "handoff" || !e.customer().Human || len(e.model.Requests()) != 0 {
			t.Errorf("%s：%+v %+v", text, r, e.model.Requests())
		}
	}
}

// 想找真人：轉人工關著就回固定回覆 handoff_off（帶入商家的聯絡方式，商家改過的優先），不轉人工
func TestHumanRequestWithHandoffOff(t *testing.T) {
	e := newEnv(t, setFeature(testMerchant(), "handoff", false), map[string][]map[string]model.Value{"Intent": {intent("human_request"), intent("human_request")}})
	if r := e.say("我要找真人"); r.Key != "handoff_off" || r.Text != "目前沒有真人客服，可以透過這個方式聯絡我們：電話 02-1234-5678" {
		t.Errorf("got %+v", r)
	}
	if r := e.say("真人"); r.Silent {
		t.Errorf("沒有轉人工，AI 要繼續回覆")
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
	c := setFeature(testMerchant(), "handoff", false)
	c.Replies = []merchant.Reply{{Key: "handoff_off", Lang: "zh", Text: "請來電 {contact}"}}
	e = newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
	if r := e.say("我要找真人"); r.Key != "handoff_off" || r.Text != "請來電 電話 02-1234-5678" {
		t.Errorf("商家改過：%+v", r)
	}
}

// 連續抱怨（#62）：上一句已經道歉過，再抱怨就轉人工，不重複回同一句道歉
func TestRepeatedComplaintHandsOff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("complaint"), intent("complaint")}})
	first := e.say("你一直轉店員")
	if first.Key != "complaint" {
		t.Errorf("第一次：%+v", first)
	}
	if r := e.say("你在回答什麼"); r.Key != "handoff" || r.Text == first.Text {
		t.Errorf("第二次：%+v", r)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "handoff：🙋 轉人工：顧客連續抱怨" {
		t.Errorf("待處理：%q", p)
	}
}

// 連續抱怨、轉人工關著（#62）：第二次回 handoff_off（商家的聯絡方式），不再道歉
func TestRepeatedComplaintWithHandoffOff(t *testing.T) {
	e := newEnv(t, setFeature(testMerchant(), "handoff", false), map[string][]map[string]model.Value{"Intent": {intent("complaint"), intent("complaint")}})
	first := e.say("你一直轉店員")
	if first.Key != "complaint" {
		t.Errorf("第一次：%+v", first)
	}
	if r := e.say("你在回答什麼"); r.Key != "handoff_off" || r.Text == first.Text {
		t.Errorf("第二次：%+v", r)
	}
}

// 只有道歉的下一句才算連續抱怨；中間隔了別的話，再抱怨照樣先道歉（#62）
func TestComplaintAfterOtherTurnApologizes(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("complaint"), intent("unclear"), intent("complaint")}})
	e.say("你一直轉店員")
	e.say("嗯")
	if r := e.say("你在回答什麼"); r.Key != "complaint" {
		t.Errorf("隔一句就轉人工：%+v", r)
	}
}

// 商家設定的固定回覆（資料庫）優先：先找顧客的語言，再找 *；都沒有才用內建的
func TestMerchantRepliesOverrideBuiltin(t *testing.T) {
	c := testMerchant()
	c.Replies = []merchant.Reply{{Key: "bot_identity", Lang: "zh", Text: "我是{merchant}的小茶"}, {Key: "complaint", Lang: "*", Text: "抱歉"}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("identity"), intent("complaint"), intent("unclear")}})
	if r := e.say("你是誰"); r.Text != "我是森日茶屋的小茶還不知道怎麼稱呼您呢？" {
		t.Errorf("identity：%+v", r)
	}
	if r := e.say("爛"); r.Text != "抱歉" {
		t.Errorf("complaint：%+v", r)
	}
	if r := e.say("嗯？"); r.Text != "不太確定您的意思，可以再說詳細一點嗎？" {
		t.Errorf("unclear：%+v", r)
	}
}

// 顧客換語言就跟著換；表情、標點沿用上一句的語言
func TestReplyInCustomerLanguage(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general"), intent("general"), intent("general")}})
	for _, c := range []struct{ text, want string }{{"Thanks", "You're welcome~ Let me know if you need anything else!"}, {"👍", "Got it~"}, {"ありがとう", "どういたしまして～また何かあればお気軽にどうぞ！"}, {"謝謝", "不客氣～有需要隨時找我喔！"}} {
		if r := e.say(c.text); r.Text != c.want {
			t.Errorf("%s：%+v，應為 %s", c.text, r, c.want)
		}
	}
}

// 沒啟用的語言：開待處理事項給員工，回「目前只支援…」（列出啟用的語言），不呼叫 AI；之後改用支援的語言就恢復
func TestUnsupportedLanguage(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general")}})
	r := e.say("안녕하세요")
	if r.Key != "unsupported" || !strings.Contains(r.Text, "中文、English、日本語") {
		t.Errorf("got %+v", r)
	}
	if n := len(e.model.Requests()); n != 0 {
		t.Errorf("呼叫了 AI %d 次", n)
	}
	if p := e.pending(); len(p) != 1 || p[0] != "unsupported：🌐 不支援的語言（ko）：안녕하세요" {
		t.Errorf("待處理：%q", p)
	}
	if r := e.say("你們幾點開"); r.Key != "general" {
		t.Errorf("改用中文：%+v", r)
	}
	// 商家沒勾英文：英文也不支援
	c := testMerchant()
	c.Languages = []merchant.Language{{Code: "zh", Enabled: true}, {Code: "en", Enabled: false}}
	e = newEnv(t, c, nil)
	if r := e.say("What time do you close?"); r.Key != "unsupported" || !strings.Contains(r.Text, "中文") || strings.Contains(r.Text, "English") {
		t.Errorf("沒勾英文：%+v", r)
	}
}

// 分類時把目前的狀態告訴 AI：在等員工時「有消息了嗎」就是追問；預約中剛問了哪個欄位
func TestClassifyTellsState(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("claim"), intent("booking")}})
	e.say("上次說打五折")
	c, _ := e.bot.LoadCustomer(context.Background(), "telegram", "U1")
	c.AskingField = "people"
	e.bot.saveCustomer(context.Background(), "telegram", "U1", c)
	e.say("四個")
	reqs := e.model.Requests()
	if strings.Contains(reqs[0].Instructions, "轉給店員") {
		t.Errorf("一開始沒有狀態：%s", reqs[0].Instructions)
	}
	if !strings.Contains(reqs[1].Instructions, "目前有一件事已轉給店員處理、尚未回覆。目前正在幫顧客預約，剛剛問了顧客的人數。") {
		t.Errorf("狀態：%s", reqs[1].Instructions)
	}
}

// fakeEmbedder：每句話的向量照表回答（表裡沒有的是零向量），記下算過哪些句子
type fakeEmbedder struct {
	vectors  map[string][]float32
	mu       sync.Mutex
	embedded []string
}

func (f *fakeEmbedder) Embed(texts []string) ([][]float32, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	out := make([][]float32, len(texts))
	for i, s := range texts {
		f.embedded = append(f.embedded, s)
		out[i] = f.vectors[s]
		if out[i] == nil {
			out[i] = []float32{0, 0, 0}
		}
	}
	return out, nil
}

// 意圖案例庫裡最像的 K 句附在分類的指示裡；案例的向量算過就記住，不重算
func TestClassifyShowsSimilarExamples(t *testing.T) {
	c := testMerchant()
	c.Examples = []merchant.Example{
		{Text: "你們開到幾點", Intent: "general", Builtin: true},
		{Text: "我要訂位", Intent: "booking", Builtin: true},
		{Text: "什麼爛客服", Intent: "complaint"},
	}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general"), intent("general")}})
	emb := &fakeEmbedder{vectors: map[string][]float32{
		"幾點打烊": {1, 0, 0}, "你們開到幾點": {0.9, 0.1, 0}, "我要訂位": {0.5, 0.5, 0}, "什麼爛客服": {0, 0, 1},
	}}
	e.bot.Embedder, e.bot.K = emb, 2
	e.say("幾點打烊")
	e.say("幾點打烊")
	ins := e.model.Requests()[0].Instructions
	if !strings.Contains(ins, "\n參考範例（類似訊息的正確分類）：\n「你們開到幾點」 → general\n「我要訂位」 → booking") || strings.Contains(ins, "什麼爛客服") {
		t.Errorf("指示：%s", ins)
	}
	if got := strings.Join(emb.embedded, ","); got != "幾點打烊,你們開到幾點,我要訂位,什麼爛客服,幾點打烊" {
		t.Errorf("算過的句子：%s", got)
	}
}

// 案例庫裡有完全相同的句子（比對規則同 check_overlap，狀態也相同）：直接用標好的意圖，不問 AI
func TestExactExampleSkipsAI(t *testing.T) {
	c := testMerchant()
	c.Examples = []merchant.Example{{Text: "店員回了嗎", Intent: "follow_up", State: "ticket"}, {Text: "什麼爛客服", Intent: "complaint"}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general")}})
	if r := e.say("什麼爛客服！"); r.Key != "complaint" {
		t.Errorf("相同句子：%+v", r)
	}
	if n := len(e.model.Requests()); n != 0 {
		t.Errorf("呼叫了 AI %d 次", n)
	}
	// 狀態不同（沒有在等員工）：照常問 AI
	if r := e.say("店員回了嗎"); r.Key != "general" || len(e.model.Requests()) != 1 {
		t.Errorf("狀態不同：%+v", r)
	}
}

// 相似案例附上分類當下的狀態
func TestSimilarExamplesShowState(t *testing.T) {
	c := testMerchant()
	c.Examples = []merchant.Example{{Text: "店員回了嗎", Intent: "follow_up", State: "ticket"}, {Text: "五位", Intent: "booking", State: "booking:people"}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general")}})
	e.bot.Embedder = &fakeEmbedder{vectors: map[string][]float32{}}
	e.say("好了沒")
	ins := e.model.Requests()[0].Instructions
	if !strings.Contains(ins, "「店員回了嗎」（有事在等店員時） → follow_up\n「五位」（預約中問人數時） → booking") {
		t.Errorf("指示：%s", ins)
	}
}

// 只做分類（分類題組用）：和正式流程同一個分類，狀態由呼叫的人給，不存任何東西
func TestClassifyOnly(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("follow_up"), intent("booking")}})
	got, err := e.bot.Classify(context.Background(), "有消息了嗎", Customer{OpenTicket: true})
	if err != nil || got != "follow_up" {
		t.Errorf("got %q %v", got, err)
	}
	if ins := e.model.Requests()[0].Instructions; !strings.Contains(ins, "目前有一件事已轉給店員處理") {
		t.Errorf("指示：%s", ins)
	}
	// 極短訊息規則也一樣套用
	if got, _ := e.bot.Classify(context.Background(), "中", Customer{}); got != "unclear" {
		t.Errorf("極短訊息：%q", got)
	}
	if p := e.pending(); len(p) != 0 {
		t.Errorf("待處理：%q", p)
	}
}

// 啟動時先算好案例庫的向量（純 Go 的 e5 一句約 0.2 秒，案例庫兩百句要等 40 秒，不能讓第一位顧客等）
func TestWarmUpEmbedsExamples(t *testing.T) {
	c := testMerchant()
	c.Examples = []merchant.Example{{Text: "你們開到幾點", Intent: "general"}, {Text: "我要訂位", Intent: "booking"}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general")}})
	emb := &fakeEmbedder{vectors: map[string][]float32{}}
	e.bot.Embedder = emb
	if err := e.bot.WarmUp(context.Background()); err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(emb.embedded, ","); got != "你們開到幾點,我要訂位" {
		t.Errorf("WarmUp 算的：%s", got)
	}
	e.say("幾點打烊")
	if got := strings.Join(emb.embedded, ","); got != "你們開到幾點,我要訂位,幾點打烊" {
		t.Errorf("顧客訊息只算訊息本身：%s", got)
	}
}

// FAQ 超過門檻題數時，啟動時也先算好 FAQ 問法的向量（挑 FAQ 要先篩選，#65）；沒超過不算
func TestWarmUpEmbedsManyFAQ(t *testing.T) {
	for _, n := range []int{merchant.FAQFilterAbove, merchant.FAQFilterAbove + 1} {
		e := newEnv(t, manyFAQ(n), map[string][]map[string]model.Value{})
		emb := &fakeEmbedder{vectors: map[string][]float32{}}
		e.bot.Embedder = emb
		if err := e.bot.WarmUp(context.Background()); err != nil {
			t.Fatal(err)
		}
		want := 0
		if n > merchant.FAQFilterAbove {
			want = n + 1 // 每題的主問法，加上最後一題的其他問法
		}
		if !slices.Contains(emb.embedded, "幾點打烊") == (want > 0) || len(emb.embedded) != want {
			t.Errorf("%d 題 FAQ：WarmUp 算了 %d 句：%v", n, len(emb.embedded), emb.embedded)
		}
	}
}

// 分類的選項和後台的意圖白話名稱表（merchant.Intents）是同一份、同樣的順序
func TestIntentsMatchMerchantTable(t *testing.T) {
	var a, b []string
	for _, in := range intents {
		a = append(a, in.code)
	}
	for _, in := range merchant.Intents {
		b = append(b, in.Code)
	}
	if strings.Join(a, ",") != strings.Join(b, ",") {
		t.Errorf("分類 %v，後台 %v", a, b)
	}
}

func TestGeneralRepliesReceived(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general")}})
	r := e.say("你們幾點開")
	if r.Key != "general" || r.Text != "收到～" {
		t.Errorf("got %+v", r)
	}
}

// 寒暄依種類回對應的固定回覆（#54）：打招呼回 greeting，回覆種類還是 general
func TestGeneralChatReplies(t *testing.T) {
	cases := []struct{ text, want string }{
		{"你好", "您好～有什麼可以幫您的嗎？"},
		{"hi", "Hi there~ How can I help you?"},
		{"this is nice", "Got it~"}, // 英文要整個字，this 不算 hi
		{"謝謝", "不客氣～有需要隨時找我喔！"},
		{"Thank you!", "You're welcome~ Let me know if you need anything else!"},
		{"ありがとう", "どういたしまして～また何かあればお気軽にどうぞ！"},
		{"掰掰", "謝謝您，期待再為您服務～"},
		{"謝謝 再見", "謝謝您，期待再為您服務～"}, // 道別優先於道謝
		{"bye", "Thank you, hope to see you again soon~"},
		{"またね", "ありがとうございました、またのご利用をお待ちしております～"},
		{"收到", "收到～"},
		{"你在瞎掰吧", "收到～"}, // 瞎掰不是道別
	}
	var script []map[string]model.Value
	for range cases {
		script = append(script, intent("general"))
	}
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": script})
	for _, c := range cases {
		if r := e.say(c.text); r.Key != "general" || r.Text != c.want {
			t.Errorf("%s：%+v，應為 %s", c.text, r, c.want)
		}
	}
}

// 平台層的 bot 指令（/start）：回打招呼，用商家慣用語言；不呼叫模型、不拿指令判斷顧客語言
func TestGreetUsesPreferredLanguage(t *testing.T) {
	c := testMerchant()
	c.Merchant.PreferredLanguage = "en"
	e := newEnv(t, c, nil)
	r, err := e.bot.Greet(context.Background(), Message{Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "/start"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Key != "greeting" || r.Text != "Hi there~ How can I help you?" || len(e.model.Requests()) != 0 {
		t.Errorf("got %+v %+v", r, e.model.Requests())
	}
}

// 老顧客送 bot 指令：用已經知道的語言回
func TestGreetUsesKnownLanguage(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("general")}})
	e.say("Good evening")
	r, err := e.bot.Greet(context.Background(), Message{Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "/start"})
	if err != nil {
		t.Fatal(err)
	}
	if r.Key != "greeting" || r.Text != "Hi there~ How can I help you?" || e.customer().Lang != "en" {
		t.Errorf("got %+v，語言 %s", r, e.customer().Lang)
	}
}

// bot 指令是 /英數字（可以接 @bot名稱、再接參數）；斜線後面接空白或中文的照常送進對話層
func TestIsBotCommand(t *testing.T) {
	for text, want := range map[string]bool{
		"/start": true, " /start ": true, "/start@sunchabot": true, "/help": true, "/start abc": true,
		"/ 請問": false, "/你好": false, "你好 /start": false, "/": false, "1/2": false, "/2": false, "/3 人": false, // 指令要英文字母開頭，數字可能是在選選項或回人數
	} {
		if got := IsBotCommand(text); got != want {
			t.Errorf("%q：%v，應為 %v", text, got, want)
		}
	}
}

// 轉人工中收到 bot 指令：和其他訊息一樣不回覆
func TestGreetSilentDuringHandoff(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("human_request")}})
	e.say("我要找真人")
	r, err := e.bot.Greet(context.Background(), Message{Platform: "telegram", UserID: "U1", ChatID: "C1", Text: "/start"})
	if err != nil {
		t.Fatal(err)
	}
	if !r.Silent || r.Key != "silent" {
		t.Errorf("got %+v", r)
	}
}

// 開場白（#56）：沒有狀態時不問 AI、不挑 FAQ、不確認是不是在問店家，回「請說」類的固定回覆
func TestOpenerRepliesWithoutAI(t *testing.T) {
	c := testMerchant()
	c.Knowledge = []merchant.Section{{Title: "營業時間", Body: "每天 10:00～20:00"}}
	e := newEnv(t, c, nil)
	cases := []struct{ text, want string }{
		{"有人在嗎", "在喔～有什麼可以幫您？"},
		{"我想問你個事情", "在喔～有什麼可以幫您？"},
		{"請問一下～", "在喔～有什麼可以幫您？"},
		{"在嗎?", "在喔～有什麼可以幫您？"},
		{"嗨嗨，有人在嗎", "在喔～有什麼可以幫您？"},
		{"可以問你一個問題嗎", "在喔～有什麼可以幫您？"},
		{"請問有人在嗎", "在喔～有什麼可以幫您？"},
		{"我想請問一下", "在喔～有什麼可以幫您？"},
		{"想問你一件事情", "在喔～有什麼可以幫您？"},
		{"Is anyone there?", "I'm here~ How can I help you?"},
	}
	for _, c := range cases {
		if r := e.say(c.text); r.Key != "general" || r.Text != c.want {
			t.Errorf("%s：%+v", c.text, r)
		}
	}
	if reqs := e.model.Requests(); len(reqs) != 0 {
		t.Errorf("不該呼叫模型：%+v", reqs)
	}
}

// 只有符號、標點、顏文字是一般對話，不問 AI（#63）；有文字或數字的照常問 AI
func TestClassifySymbolsOnly(t *testing.T) {
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": {intent("booking"), intent("general")}})
	for _, c := range []Customer{{}, {OpenTicket: true}, {AskingField: "people"}} {
		for _, text := range []string{"＝＝", "…", "^^", "?!", "🙂"} {
			if got, _ := e.bot.Classify(context.Background(), text, c); got != "general" {
				t.Errorf("%s（%+v）：%s", text, c, got)
			}
		}
	}
	if n := len(e.model.Requests()); n != 0 {
		t.Fatalf("呼叫模型 %d 次", n)
	}
	for _, text := range []string{"４", "ok"} {
		e.bot.Classify(context.Background(), text, Customer{AskingField: "people"})
	}
	if n := len(e.model.Requests()); n != 2 {
		t.Errorf("有文字或數字要問 AI：呼叫 %d 次", n)
	}
}

// 有講出要問什麼的、只有打招呼的不是開場白，照常問 AI；有事在等店員時「有人在嗎」是在追問，也照常問 AI
func TestNotOpenerAsksAI(t *testing.T) {
	texts := []string{"我想問營業時間", "請問幾點開", "你好", "我想問你事情你就轉店員"}
	var script []map[string]model.Value
	for range texts {
		script = append(script, intent("general"))
	}
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": append(script, intent("follow_up"))})
	for _, text := range texts {
		if r := e.say(text); r.Text == "在喔～有什麼可以幫您？" {
			t.Errorf("%s 不是開場白：%+v", text, r)
		}
	}
	if n := len(e.model.Requests()); n != len(texts) {
		t.Fatalf("呼叫模型 %d 次，應為 %d", n, len(texts))
	}
	if got, _ := e.bot.Classify(context.Background(), "有人在嗎", Customer{OpenTicket: true}); got != "follow_up" {
		t.Errorf("等店員時：%q", got)
	}
}

// 抱怨被轉給店員是抱怨，不是想找真人（#56）：AI 分成想找真人時程式改判；真的要求找真人照舊
func TestTransferComplaintIsComplaint(t *testing.T) {
	cases := []struct{ text, want string }{
		{"我想問你事情你就轉店員", "complaint"},
		{"怎麼又轉人工", "complaint"},
		{"每次問一下就把我丟給店員", "complaint"},
		{"幹嘛一直叫我找真人", "complaint"},
		{"你們只會叫我去找店員喔", "complaint"},
		{"怎麼我一問問題就轉人工", "complaint"},
		{"我要找真人", "human_request"},
		{"請幫我轉店員", "human_request"},
		{"那就轉人工吧", "human_request"},
		{"為什麼不能轉真人", "human_request"},
		{"如果問不到就轉人工", "human_request"},
		{"怎麼都不幫我轉人工", "human_request"},
		{"幹嘛一直 轉 店員", "complaint"}, // 空白、標點不影響
	}
	var script []map[string]model.Value
	for range cases {
		script = append(script, intent("human_request"))
	}
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": script})
	for _, c := range cases {
		if got, _ := e.bot.Classify(context.Background(), c.text, Customer{}); got != c.want {
			t.Errorf("%s：%s，應為 %s", c.text, got, c.want)
		}
	}
}

// 問店裡的人（他、老闆、店員是誰、叫什麼）不是問客服或顧客本人的身分（#56）：AI 分成問身分時程式改成一般對話
func TestThirdPersonWhoIsGeneral(t *testing.T) {
	cases := []struct{ text, want string }{
		{"你老闆叫什麼？", "general"},
		{"他是誰", "general"},
		{"店長叫什麼名字", "general"},
		{"那位姐姐是誰啊", "general"},
		{"你們老闆娘姓什麼", "general"},
		{"請問店長貴姓", "general"},
		{"櫃檯那個女生是誰", "general"},
		{"我們來角色扮演，你演老闆", "identity"},
		{"這個客服叫什麼名字", "identity"}, // 問的是客服本人
		{"這位客服是誰", "identity"},
		{"你是誰", "identity"},
		{"我叫什麼", "identity"},
		{"你叫什麼名字", "identity"},
	}
	var script []map[string]model.Value
	for range cases {
		script = append(script, intent("identity"))
	}
	e := newEnv(t, testMerchant(), map[string][]map[string]model.Value{"Intent": script})
	for _, c := range cases {
		if got, _ := e.bot.Classify(context.Background(), c.text, Customer{}); got != c.want {
			t.Errorf("%s：%s，應為 %s", c.text, got, c.want)
		}
	}
}

// 商家在後台改過打招呼：寒暄的打招呼用商家的文字
func TestGeneralGreetingUsesMerchantReply(t *testing.T) {
	c := testMerchant()
	c.Replies = []merchant.Reply{{Key: "greeting", Lang: "zh", Text: "歡迎光臨{merchant}～"}}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {intent("general")}})
	if r := e.say("哈囉"); r.Key != "general" || r.Text != "歡迎光臨森日茶屋～" {
		t.Errorf("got %+v", r)
	}
}
