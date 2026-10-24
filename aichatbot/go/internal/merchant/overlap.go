package merchant

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode"
)

// Normalize 是 eval/check_overlap.py 的比對規則：轉小寫，只留字母和數字（等於 Python 的 re.sub(r"[\W_]+", "", t.lower())）
func Normalize(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			return r
		}
		return -1
	}, strings.ToLower(s))
}

// TestSet 是題庫的句子：Normalize 後的句子 → 出現在哪些題庫。案例庫的句子不能和題庫重複（案例和考題重複，分數就只是背答案）
type TestSet map[string][]string

// LoadTestSet 讀 evalDir 底下的題庫，範圍和 check_overlap.py 一樣；任一檔讀不到或格式錯就回錯誤。
// 驗證題組（holdout/）兩個檔都叫「驗證題組」：錯誤訊息會給商家和 AI 看到，不能看出是哪一個檔
func LoadTestSet(evalDir string) (TestSet, error) {
	t := TestSet{}
	for _, f := range []struct {
		name, path string
		texts      func([]byte) ([]string, error)
	}{
		{"intent_set", "intent_set.json", itemTexts},
		{"testset", "testset.json", caseQuestions},
		{"faq_set", "faq_set.json", caseQuestions},
		{"flow_tests", "flow_tests.json", turnUsers},
		// 第二家測試商家（美髮沙龍，#66）的題組
		{"salon/intent_set", "stores/salon/intent_set.json", itemTexts},
		{"salon/testset", "stores/salon/testset.json", caseQuestions},
		{"salon/flow_tests", "stores/salon/flow_tests.json", turnUsers},
		{"驗證題組", "holdout/intent_holdout.json", itemTexts},
		{"驗證題組", "holdout/flow_holdout.json", turnUsers},
	} {
		data, err := os.ReadFile(filepath.Join(evalDir, f.path))
		if err != nil {
			return nil, err
		}
		texts, err := f.texts(data)
		if err != nil {
			return nil, fmt.Errorf("%s：%w", f.path, err)
		}
		for _, s := range texts {
			if n := Normalize(s); !slices.Contains(t[n], f.name) {
				t[n] = append(t[n], f.name)
			}
		}
	}
	return t, nil
}

// Find 回 text 出現在哪些題庫（排序過）；沒有重複時是空的
func (t TestSet) Find(text string) []string {
	sets := slices.Clone(t[Normalize(text)])
	slices.Sort(sets)
	return sets
}

func itemTexts(data []byte) ([]string, error) {
	var f struct{ Items []struct{ Text string } }
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	var out []string
	for _, i := range f.Items {
		out = append(out, i.Text)
	}
	return out, nil
}

func caseQuestions(data []byte) ([]string, error) {
	var f struct {
		Cases []struct {
			Q string `json:"q"`
		}
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	var out []string
	for _, c := range f.Cases {
		out = append(out, c.Q)
	}
	return out, nil
}

func turnUsers(data []byte) ([]string, error) {
	var f struct {
		Conversations []struct{ Turns []struct{ User string } }
	}
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	var out []string
	for _, c := range f.Conversations {
		for _, t := range c.Turns {
			out = append(out, t.User)
		}
	}
	return out, nil
}
