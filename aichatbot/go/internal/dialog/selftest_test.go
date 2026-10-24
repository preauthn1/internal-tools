package dialog

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"aichatbot/internal/merchant"
	"aichatbot/internal/model"
)

// 自我測試：每一句例句（案例庫＋自訂意圖）各分類一次，分類時排除這句本身（不然完全相同的句子會直接沿用），
// 帶著例句的狀態；列出被分到別的意圖的例句，AI 失敗的另外列
func TestSelfTestFindsConflicts(t *testing.T) {
	c := customMerchant(giftBox())
	c.Examples = []merchant.Example{
		{Text: "我要訂位", Intent: "booking", Builtin: true},
		{Text: "店員回了嗎", Intent: "follow_up", State: "ticket"},
	}
	e := newEnv(t, c, map[string][]map[string]model.Value{"Intent": {
		intent("booking"), intent("follow_up"), intent("custom_1"), intent("general"), nil,
	}})
	var progress []int
	report, err := e.bot.SelfTest(context.Background(), func(done, total int) { progress = append(progress, done) })
	if err != nil {
		t.Fatal(err)
	}
	reqs := e.model.Requests()
	var texts []string
	for _, r := range reqs {
		texts = append(texts, r.Text)
	}
	if got := strings.Join(texts, ","); got != "我要訂位,店員回了嗎,有賣禮盒嗎,想買來送人,過年送禮推薦" {
		t.Fatalf("分類的句子：%s", got)
	}
	if !strings.Contains(reqs[1].Instructions, "目前有一件事已轉給店員處理") {
		t.Errorf("沒帶狀態：%s", reqs[1].Instructions)
	}
	// 自訂意圖的選項說明附的例句也要排除受測的那句
	for _, f := range reqs[2].Fields {
		if f.Name == "intent" && (strings.Contains(f.Description, "有賣禮盒嗎") || !strings.Contains(f.Description, "例：想買來送人、過年送禮推薦")) {
			t.Errorf("選項說明：%s", f.Description)
		}
	}
	if report.Total != 5 || len(report.Conflicts) != 1 || len(report.Failed) != 1 {
		t.Fatalf("報告：%+v", report)
	}
	if m := report.Conflicts[0]; m.Text != "想買來送人" || m.Expected != "custom_1" || m.Got != "general" {
		t.Errorf("衝突：%+v", m)
	}
	if f := report.Failed[0]; f.Text != "過年送禮推薦" || f.Expected != "custom_1" {
		t.Errorf("AI 失敗：%+v", f)
	}
	if got := len(progress); got != 5 || progress[4] != 5 {
		t.Errorf("進度：%v", progress)
	}
}

// slowModel 每次呼叫花 delay；Intent 照 intents（句子 → 意圖）回答，其他題目回 nil（AI 失敗）。
// 記下每次呼叫的句子和開始時間；started 收到每次呼叫開始的句子；句子是 holdText 時等 hold 關掉才回答
type slowModel struct {
	delay    time.Duration
	intents  map[string]string
	started  chan string
	holdText string
	hold     chan struct{}

	mu    sync.Mutex
	calls []slowCall
}

type slowCall struct {
	name, text string
	start      time.Time
}

func (m *slowModel) Respond(ctx context.Context, r model.Request) (*model.Reply, model.Usage) {
	m.mu.Lock()
	m.calls = append(m.calls, slowCall{r.Name, r.Text, time.Now()})
	m.mu.Unlock()
	select {
	case m.started <- r.Text:
	default:
	}
	time.Sleep(m.delay)
	if r.Text == m.holdText {
		<-m.hold
	}
	if code, ok := m.intents[r.Text]; ok && r.Name == "Intent" {
		return &model.Reply{Values: intent(code)}, model.Usage{}
	}
	return nil, model.Usage{}
}

func (m *slowModel) Calls() []slowCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Clone(m.calls)
}

// 顧客訊息優先（#58）：自我測試進行中顧客傳訊息，顧客最多等自我測試正在跑的那一個呼叫；
// 顧客的回合（分類、抽預約資料兩次呼叫）結束前自我測試不再送新的呼叫；自我測試最後照樣跑完、報告正確
func TestSelfTestYieldsToCustomers(t *testing.T) {
	c := testMerchant()
	intents := map[string]string{"明天兩位": "booking"}
	c.Examples = nil
	for i := range 8 {
		text := fmt.Sprintf("例句%d", i)
		c.Examples = append(c.Examples, merchant.Example{Text: text, Intent: "general"})
		intents[text] = "general"
	}
	intents["例句7"] = "booking" // 一句衝突
	e := newEnv(t, c, nil)
	slow := &slowModel{delay: 30 * time.Millisecond, intents: intents, started: make(chan string, 16)} // 放得下每次呼叫，不會漏掉「例句1」
	e.bot.Model = model.NewQueued(slow)

	type result struct {
		report SelfTestReport
		err    error
	}
	done := make(chan result, 1)
	var progress []int
	go func() {
		r, err := e.bot.SelfTest(context.Background(), func(done, total int) { progress = append(progress, done) })
		done <- result{r, err}
	}()
	for text := range slow.started {
		if text == "例句1" {
			break
		}
	}
	sent := time.Now()
	e.say("明天兩位")
	turnEnd := time.Now()
	res := <-done
	if res.err != nil {
		t.Fatal(res.err)
	}

	var order []string
	var during int
	for _, call := range slow.Calls() {
		order = append(order, call.name+":"+call.text)
		if call.text == "明天兩位" || !call.start.After(sent) {
			continue
		}
		if call.start.Before(turnEnd) {
			during++
		}
	}
	// 顧客只等送出訊息時正在跑的「例句1」；顧客的兩次呼叫中間、回合結束前，自我測試都沒有新的呼叫
	customerFirst := slices.Index(order, "Intent:明天兩位")
	if customerFirst < 0 || customerFirst+1 == len(order) || order[customerFirst+1] != "Booking:明天兩位" {
		t.Fatalf("顧客的呼叫中間被插隊：%v", order)
	}
	if before := customerFirst - slices.Index(order, "Intent:例句1") - 1; before > 0 || during > 0 {
		t.Errorf("顧客等了自我測試 %d 個呼叫、回合中自我測試送了 %d 個：%v", before, during, order)
	}
	if r := res.report; r.Total != 8 || len(r.Conflicts) != 1 || r.Conflicts[0].Text != "例句7" || len(r.Failed) != 0 {
		t.Errorf("報告：%+v", r)
	}
	if len(progress) != 8 || progress[7] != 8 {
		t.Errorf("進度：%v", progress)
	}
}

// 等顧客回合結束時取消（重跑、取消）：自我測試馬上停下來回 ctx.Err()，不送任何呼叫
func TestSelfTestCanceledWhileWaitingForCustomer(t *testing.T) {
	c := testMerchant()
	c.Examples = []merchant.Example{{Text: "例句", Intent: "general"}}
	e := newEnv(t, c, nil)
	slow := &slowModel{started: make(chan string, 1), holdText: "明天兩位", hold: make(chan struct{})}
	e.bot.Model = slow
	customerDone := make(chan struct{})
	go func() {
		e.say("明天兩位")
		close(customerDone)
	}()
	<-slow.started

	ctx, cancel := context.WithCancel(context.Background())
	errs := make(chan error, 1)
	go func() {
		_, err := e.bot.SelfTest(ctx, nil)
		errs <- err
	}()
	time.Sleep(50 * time.Millisecond) // 讓自我測試走到等顧客那一步
	cancel()
	select {
	case err := <-errs:
		if !errors.Is(err, context.Canceled) {
			t.Errorf("err：%v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("取消後自我測試沒有停下來")
	}
	close(slow.hold)
	<-customerDone
	for _, call := range slow.Calls() {
		if call.text == "例句" {
			t.Errorf("取消了還送出自我測試的呼叫：%+v", slow.Calls())
		}
	}
}
