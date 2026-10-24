package merchant

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// 新匯入的預設值：預約、轉人工都開；慣用語言中文，顧客語言中英日
var (
	defaultFeatures  = []Feature{{Name: "booking", Enabled: true}, {Name: "handoff", Enabled: true}}
	defaultLanguages = []Language{{Code: "zh", Enabled: true}, {Code: "en", Enabled: true}, {Code: "ja", Enabled: true}}
)

const preferredLanguage = "zh"

// FromEvalDir 讀 eval/ 的森日茶屋檔案（store.json、faq.txt、replies.json、examples.jsonl，介紹和聯絡方式參考 shop.txt），
// 轉成商家設定。store.json 的結構化規則轉成知識段落文字（ADR 0003 拿掉了結構化欄位）
func FromEvalDir(dir string) (Config, error) {
	var c Config
	st, err := readStoreJSON(filepath.Join(dir, "store.json"))
	if err != nil {
		return c, err
	}
	if c.FAQ, err = readFAQ(filepath.Join(dir, "faq.txt")); err != nil {
		return c, err
	}
	if c.Replies, err = readReplies(filepath.Join(dir, "replies.json")); err != nil {
		return c, err
	}
	if c.Examples, err = readExamples(filepath.Join(dir, "examples.jsonl")); err != nil {
		return c, err
	}
	shop, _ := os.ReadFile(filepath.Join(dir, "shop.txt")) // 沒有也沒關係
	c.Merchant = Merchant{
		Name:              st.Name,
		Intro:             st.intro(),
		Contact:           contact(string(shop), c.FAQ),
		PreferredLanguage: preferredLanguage,
	}
	c.Knowledge = st.sections()
	c.Features = append([]Feature(nil), defaultFeatures...)
	c.Languages = append([]Language(nil), defaultLanguages...)
	c.CustomIntents = []CustomIntent{}
	return c, c.Validate()
}

// evalStore 是 store.json 用得到的欄位
type evalStore struct {
	Name           string `json:"name"`
	OpenHour       int    `json:"openHour"`
	CloseHour      int    `json:"closeHour"`
	ClosedWeekdays []int  `json:"closedWeekdays"` // Swift Calendar 的星期：1=週日、2=週一…
	MaxPeople      int    `json:"maxPeople"`
	CardMinimum    int    `json:"cardMinimum"`
	PickupDiscount *struct {
		MinTotal int     `json:"minTotal"`
		Rate     float64 `json:"rate"`
	} `json:"pickupDiscount"`
	Menu []struct {
		Price int               `json:"price"`
		Names map[string]string `json:"names"`
	} `json:"menu"`
}

func readStoreJSON(path string) (evalStore, error) {
	var s evalStore
	data, err := os.ReadFile(path)
	if err != nil {
		return s, err
	}
	if err := json.Unmarshal(data, &s); err != nil {
		return s, fmt.Errorf("%s：%w", path, err)
	}
	if s.Name == "" {
		return s, fmt.Errorf("%s：缺少 name", path)
	}
	return s, nil
}

var weekdayNames = []string{"日", "一", "二", "三", "四", "五", "六"}

func (s evalStore) menuNames() []string {
	var names []string
	for _, m := range s.Menu {
		names = append(names, m.Names[preferredLanguage])
	}
	return names
}

func (s evalStore) intro() string {
	if len(s.Menu) == 0 {
		return s.Name
	}
	return fmt.Sprintf("%s，招牌飲品有%s。", s.Name, strings.Join(s.menuNames(), "、"))
}

// sections 把 store.json 的規則寫成知識段落；沒有的規則不產生段落
func (s evalStore) sections() []Section {
	var out []Section
	if s.CloseHour > s.OpenHour {
		body := fmt.Sprintf("%02d:00–%02d:00 營業", s.OpenHour, s.CloseHour)
		var closed []string
		for _, d := range s.ClosedWeekdays {
			if d >= 1 && d <= 7 {
				closed = append(closed, "週"+weekdayNames[d-1])
			}
		}
		if len(closed) > 0 {
			body += "，" + strings.Join(closed, "、") + "公休"
		}
		out = append(out, Section{Title: "營業時間", Body: body + "。"})
	}
	if len(s.Menu) > 0 {
		var lines []string
		for i, name := range s.menuNames() {
			lines = append(lines, fmt.Sprintf("%s %d 元", name, s.Menu[i].Price))
		}
		out = append(out, Section{Title: "菜單與價格", Body: strings.Join(lines, "\n")})
	}
	if s.CardMinimum > 0 {
		out = append(out, Section{Title: "付款", Body: fmt.Sprintf("信用卡需消費滿 %d 元。", s.CardMinimum)})
	}
	if d := s.PickupDiscount; d != nil && d.Rate > 0 && d.Rate < 1 {
		out = append(out, Section{Title: "自取優惠", Body: fmt.Sprintf("店內自取滿 %d 元打 %s 折。", d.MinTotal, strconv.FormatFloat(d.Rate*10, 'f', -1, 64))})
	}
	if s.MaxPeople > 0 {
		out = append(out, Section{Title: "人數上限", Body: fmt.Sprintf("店內最多 %d 位。", s.MaxPeople)})
	}
	return out
}

// contact：shop.txt 的「- 地址：」那行，沒有就用 FAQ 的 address 答案，都沒有就填佔位文字（後台再改）
func contact(shop string, faq []FAQ) string {
	for _, line := range strings.Split(shop, "\n") {
		if rest, ok := strings.CutPrefix(strings.TrimSpace(line), "- 地址："); ok && rest != "" {
			return "地址：" + rest
		}
	}
	for _, f := range faq {
		if f.Key == "address" {
			return f.Answer
		}
	}
	return "（請在後台填寫聯絡方式）"
}

// readFAQ：一行一題，「代號 | 常見問法（/ 分隔） | 中文答案 | English | 日本語」；# 開頭和空行略過。
// 只存慣用語言（中文）的答案；第一個問法當問題，其餘放進其他問法
func readFAQ(path string) ([]FAQ, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := []FAQ{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.Split(line, " | ")
		if len(parts) < 3 {
			return nil, fmt.Errorf("%s 第 %d 行：至少要有代號、問法、中文答案三欄", path, n)
		}
		var questions []string
		for _, q := range strings.Split(parts[1], "/") {
			if q = strings.TrimSpace(q); q != "" {
				questions = append(questions, q)
			}
		}
		if len(questions) == 0 {
			return nil, fmt.Errorf("%s 第 %d 行：沒有問法", path, n)
		}
		out = append(out, FAQ{
			Key:          strings.TrimSpace(parts[0]),
			Question:     questions[0],
			AltQuestions: append([]string{}, questions[1:]...),
			Answer:       strings.TrimSpace(parts[2]),
		})
	}
	return out, sc.Err()
}

// usesStructuredVars：有沒有用到 store.json 結構化規則的變數（營業時間、人數上限、自取門檻）。
// ADR 0003 拿掉了這些欄位，對話層填不進去，這種回覆不匯入（ask_time、err_hours…），對話層改用內建的說法
func usesStructuredVars(langs map[string]string) bool {
	for _, text := range langs {
		for _, v := range []string{"{open}", "{close}", "{max}", "{min}"} {
			if strings.Contains(text, v) {
				return true
			}
		}
	}
	return false
}

// readReplies：{"key": {"zh": "...", "en": "...", "*": "..."}}，「說明」不是回覆。依 key、lang 排序
func readReplies(path string) ([]Reply, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%s：%w", path, err)
	}
	out := []Reply{}
	for key, v := range raw {
		if key == "說明" {
			continue
		}
		var langs map[string]string
		if err := json.Unmarshal(v, &langs); err != nil {
			return nil, fmt.Errorf("%s 的 %s：%w", path, key, err)
		}
		if usesStructuredVars(langs) {
			continue
		}
		for lang, text := range langs {
			out = append(out, Reply{Key: key, Lang: lang, Text: text})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Key != out[j].Key {
			return out[i].Key < out[j].Key
		}
		return out[i].Lang < out[j].Lang
	})
	return out, nil
}

// readExamples：一行一個 {"text": ..., "intent": ...[, "state": ...]}，當內建案例；# 開頭和空行略過。
// Swift 版的舊代號（打招呼、道謝、道別、閒聊、問商家）改成 general
func readExamples(path string) ([]Example, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	out := []Example{}
	sc := bufio.NewScanner(bytes.NewReader(data))
	for n := 1; sc.Scan(); n++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") { // 和 Swift 版一樣略過註解
			continue
		}
		var e struct{ Text, Intent, State string }
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			return nil, fmt.Errorf("%s 第 %d 行：%w", path, n, err)
		}
		if e.Text == "" || e.Intent == "" {
			return nil, fmt.Errorf("%s 第 %d 行：要有 text 和 intent", path, n)
		}
		if code, ok := legacyIntents[e.Intent]; ok {
			e.Intent = code
		}
		out = append(out, Example{Text: e.Text, Intent: e.Intent, State: e.State, Builtin: true})
	}
	return out, sc.Err()
}

// ReadFile 讀匯入來源：資料夾當 eval/ 格式，檔案當 Export 的 JSON
func ReadFile(path string) (Config, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Config{}, err
	}
	if info.IsDir() {
		return FromEvalDir(path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	c, err := ParseJSON(data)
	if err != nil {
		return Config{}, fmt.Errorf("%s：%w", path, err)
	}
	return c, c.Validate()
}
