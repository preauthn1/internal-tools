package merchant

// 自訂意圖（CONTEXT.md「意圖」）：商家另外新增的意圖，每個綁一種動作

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// 自訂意圖的動作
const (
	ActionFixedText = "fixed_text" // 回固定文字（Text）
	ActionHandoff   = "handoff"    // 轉人工，回 Text 當提示（空的用轉人工的預設提示）
	ActionKnowledge = "knowledge"  // 只根據 Knowledge 指定的知識段落回答
)

// CustomIntent 是一個自訂意圖。Code 由程式給（custom_1～custom_10），後台只顯示 Name；
// Description 和 Examples 寫進分類的選項說明，Examples 也加進意圖案例庫一起比對
type CustomIntent struct {
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Examples    []string `json:"examples"`
	Action      string   `json:"action"`
	Text        string   `json:"text"`
	Knowledge   []string `json:"knowledge"` // 知識段落的標題
}

// MaxCustomIntents 是最多幾個自訂意圖；MinCustomExamples 是每個至少幾句例句
const (
	MaxCustomIntents  = 10
	MinCustomExamples = 3
)

// customCode 是自訂意圖的代號：custom_1～custom_10
var customCode = regexp.MustCompile(`^custom_([1-9]|10)$`)

// cleanExamples 是去掉前後空白和空白項目的例句
func (ci CustomIntent) cleanExamples() []string {
	var out []string
	for _, s := range ci.Examples {
		if s = strings.TrimSpace(s); s != "" {
			out = append(out, s)
		}
	}
	return out
}

// validateCustomIntents 是 Validate 對自訂意圖的檢查（題庫重複在後台存檔時另外檢查，見 SaveCustomIntents）
func validateCustomIntents(c Config) []string {
	var errs []string
	if len(c.CustomIntents) > MaxCustomIntents {
		errs = append(errs, fmt.Sprintf("自訂意圖最多 %d 個", MaxCustomIntents))
	}
	names := map[string]bool{}
	for _, in := range Intents {
		for _, l := range in.Labels {
			names[l] = true
		}
	}
	titles := map[string]bool{}
	for _, s := range c.Knowledge {
		titles[strings.TrimSpace(s.Title)] = true
	}
	library := map[string]string{} // 案例庫 Normalize 後的句子 → 錯誤訊息裡的說法
	for _, e := range c.Examples {
		if e.Builtin {
			library[Normalize(e.Text)] = "內建案例"
		} else {
			library[Normalize(e.Text)] = "案例庫的句子"
		}
	}
	seen := map[string]bool{}  // 自訂意圖的例句（Normalize 後）
	codes := map[string]bool{} // 代號（空的由程式給）
	for i, ci := range c.CustomIntents {
		// 代號會蓋過同代號的內建意圖，不能讓匯入檔或 API 隨便給
		if ci.Code != "" && (!customCode.MatchString(ci.Code) || codes[ci.Code]) {
			errs = append(errs, "自訂意圖的代號不正確："+ci.Code)
		}
		codes[ci.Code] = true
		name := strings.TrimSpace(ci.Name)
		for _, x := range []struct{ name, v string }{{"名稱", ci.Name}, {"說明", ci.Description}} {
			if strings.TrimSpace(x.v) == "" {
				errs = append(errs, fmt.Sprintf("第 %d 個自訂意圖的%s是必填", i+1, x.name))
			}
		}
		if name != "" && names[name] {
			errs = append(errs, "自訂意圖的名稱重複："+name)
		}
		names[name] = true
		label := "「" + name + "」"
		if name == "" {
			label = fmt.Sprintf("第 %d 個自訂意圖", i+1)
		}
		list := ci.cleanExamples()
		if len(list) < MinCustomExamples {
			errs = append(errs, fmt.Sprintf("%s至少要 %d 句例句", label, MinCustomExamples))
		}
		for _, s := range list {
			key := Normalize(s)
			switch {
			case key == "":
				errs = append(errs, label+"的例句要有文字："+s)
			case seen[key]:
				errs = append(errs, label+"的例句重複："+s)
			case library[key] != "":
				errs = append(errs, label+"的例句和"+library[key]+"重複："+s)
			}
			seen[key] = true
		}
		switch ci.Action {
		case ActionFixedText:
			if strings.TrimSpace(ci.Text) == "" {
				errs = append(errs, label+"要填回覆的文字")
			}
		case ActionKnowledge:
			if len(ci.Knowledge) == 0 {
				errs = append(errs, label+"至少要指定一段知識")
			}
			for _, title := range ci.Knowledge {
				if !titles[strings.TrimSpace(title)] {
					errs = append(errs, label+"指定的知識段落不存在："+title)
				}
			}
		case ActionHandoff:
		default:
			errs = append(errs, label+"的動作不正確")
		}
	}
	return errs
}

// CustomIntentExamples 把自訂意圖的例句轉成意圖案例（非內建、沒有狀態），分類時和案例庫一起比對
func CustomIntentExamples(list []CustomIntent) []Example {
	var out []Example
	for _, ci := range list {
		for _, s := range ci.cleanExamples() {
			out = append(out, Example{Text: s, Intent: ci.Code})
		}
	}
	return out
}

// SaveCustomIntents 只換掉自訂意圖，順序照 list 的順序；例句和題庫重複回 *OverlapError（第一句重複的），
// 沒通過檢查回 *InvalidError。還沒有商家設定時回 ErrInfoFirst
func SaveCustomIntents(ctx context.Context, db *sql.DB, tests TestSet, list []CustomIntent) error {
	for _, ci := range list {
		for _, s := range ci.cleanExamples() {
			if sets := tests.Find(s); len(sets) > 0 {
				return &OverlapError{Sets: sets, Text: s}
			}
		}
	}
	return saveBlock(ctx, db, false, func(c *Config) error { c.CustomIntents = list; return nil }, func(tx *sql.Tx, c Config) error {
		return replaceCustomIntents(ctx, tx, c.CustomIntents)
	})
}

// withCodes 給還沒有代號的自訂意圖一個還沒用過的 custom_N，已有的代號保留；名稱和例句去掉前後空白和空白的例句
func withCodes(list []CustomIntent) []CustomIntent {
	out := make([]CustomIntent, len(list))
	used := map[string]bool{}
	for i, ci := range list {
		used[ci.Code] = true
		ci.Name, ci.Examples = strings.TrimSpace(ci.Name), ci.cleanExamples()
		out[i] = ci
	}
	n := 1
	for i := range out {
		if out[i].Code != "" {
			continue
		}
		for used[fmt.Sprintf("custom_%d", n)] {
			n++
		}
		out[i].Code = fmt.Sprintf("custom_%d", n)
		used[out[i].Code] = true
	}
	return out
}

// replaceCustomIntents 用 list 取代自訂意圖（沒有代號的給新代號）。沒有留下來的代號（已經刪掉的自訂意圖），
// 屬於它的商家案例一起刪掉；要在給新代號之前算，不然新的自訂意圖拿到刪掉的那個代號時，會接收舊的案例
func replaceCustomIntents(ctx context.Context, tx *sql.Tx, list []CustomIntent) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM custom_intent"); err != nil {
		return err
	}
	kept := []string{}
	for _, ci := range list {
		if ci.Code != "" {
			kept = append(kept, ci.Code)
		}
	}
	list = withCodes(list)
	codesJSON, err := jsonList(kept)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, "DELETE FROM intent_example WHERE builtin = 0 AND intent LIKE 'custom\\_%' ESCAPE '\\' AND intent NOT IN (SELECT value FROM json_each(?))", codesJSON); err != nil {
		return err
	}
	for i, ci := range list {
		examples, err := jsonList(ci.Examples)
		if err != nil {
			return err
		}
		knowledge, err := jsonList(ci.Knowledge)
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO custom_intent (position, code, name, description, examples, action, text, knowledge) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
			i, ci.Code, ci.Name, ci.Description, examples, ci.Action, ci.Text, knowledge); err != nil {
			return err
		}
	}
	return nil
}

// jsonList 把字串陣列存成 JSON（nil 存成 []）
func jsonList(list []string) (string, error) {
	if list == nil {
		list = []string{}
	}
	b, err := json.Marshal(list)
	return string(b), err
}

// listCustomIntents 讀出自訂意圖，依存的順序
func listCustomIntents(ctx context.Context, db querier) ([]CustomIntent, error) {
	return query(ctx, db, "SELECT code, name, description, examples, action, text, knowledge FROM custom_intent ORDER BY position, id", func(r *sql.Rows) (ci CustomIntent, err error) {
		var examples, knowledge string
		if err = r.Scan(&ci.Code, &ci.Name, &ci.Description, &examples, &ci.Action, &ci.Text, &knowledge); err != nil {
			return ci, err
		}
		if err = json.Unmarshal([]byte(examples), &ci.Examples); err != nil {
			return ci, err
		}
		return ci, json.Unmarshal([]byte(knowledge), &ci.Knowledge)
	})
}
