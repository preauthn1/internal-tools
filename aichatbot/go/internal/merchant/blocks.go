package merchant

// 後台依區塊（商家資訊、FAQ、語言）個別存取商家設定

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// SupportedLanguages 是可以勾選的顧客語言（ISO 639-1），後台照這個順序列出；要加語言改這一行
var SupportedLanguages = []string{"zh", "en", "ja", "ko", "th", "vi", "id", "ms", "fr", "de", "es"}

// MeasuredLanguages 是有量過回覆品質的語言；其他語言後台會提示「品質沒量過」
var MeasuredLanguages = []string{"zh", "en", "ja"}

// Default 是還沒有商家設定時的起點：商家資訊空白，預約、轉人工都開，慣用語言中文、顧客語言中英日（ADR 0003）
func Default() Config {
	return Config{
		Merchant:  Merchant{PreferredLanguage: preferredLanguage},
		Knowledge: []Section{},
		FAQ:       []FAQ{},
		Examples:  []Example{},
		Features:  append([]Feature(nil), defaultFeatures...),
		Replies:   []Reply{},
		Languages: append([]Language(nil), defaultLanguages...),

		CustomIntents: []CustomIntent{},
	}
}

// Info 是商家資訊區塊：名稱、介紹、聯絡方式和知識段落
type Info struct {
	Name      string    `json:"name"`
	Intro     string    `json:"intro"`
	Contact   string    `json:"contact"`
	Knowledge []Section `json:"knowledge"`
}

// InfoOf 取出商家資訊區塊
func InfoOf(c Config) Info {
	return Info{Name: c.Merchant.Name, Intro: c.Merchant.Intro, Contact: c.Merchant.Contact, Knowledge: c.Knowledge}
}

// SaveInfo 只換掉商家資訊區塊。還沒有商家設定時，其他區塊用 Default() 建立
func SaveInfo(ctx context.Context, db *sql.DB, info Info) error {
	return saveBlock(ctx, db, true, func(c *Config) error {
		c.Merchant.Name, c.Merchant.Intro, c.Merchant.Contact = info.Name, info.Intro, info.Contact
		c.Knowledge = info.Knowledge
		if c.Knowledge == nil {
			c.Knowledge = []Section{}
		}
		return nil
	}, func(tx *sql.Tx, c Config) error {
		if err := writeMerchant(ctx, tx, c.Merchant); err != nil {
			return err
		}
		return replaceKnowledge(ctx, tx, c.Knowledge)
	})
}

// SaveFAQ 只換掉 FAQ，順序照 faq 的順序。代號是空的就給一個還沒用過的 faq-N（後台不讓商家看到代號）；
// 其他問法去掉前後空白和空白的項目
func SaveFAQ(ctx context.Context, db *sql.DB, faq []FAQ) error {
	out := make([]FAQ, len(faq))
	used := map[string]bool{}
	for i, f := range faq {
		f.Key = strings.TrimSpace(f.Key)
		used[f.Key] = true
		f.AltQuestions = trimAlt(f.AltQuestions)
		out[i] = f
	}
	for i := range out {
		if out[i].Key == "" {
			out[i].Key = unusedFAQKey(used)
		}
	}
	return saveBlock(ctx, db, false, func(c *Config) error { c.FAQ = out; return nil }, func(tx *sql.Tx, c Config) error {
		return replaceFAQ(ctx, tx, c.FAQ)
	})
}

// AddFAQ 在 FAQ 最後新增一題（不像 SaveFAQ 整份換掉，同時有人在後台改 FAQ 也不會蓋掉），回傳程式給的代號；
// 問題、答案去掉前後空白，其他問法去掉空白的項目
func AddFAQ(ctx context.Context, db *sql.DB, f FAQ) (string, error) {
	f.Question, f.Answer, f.AltQuestions = strings.TrimSpace(f.Question), strings.TrimSpace(f.Answer), trimAlt(f.AltQuestions)
	err := saveBlock(ctx, db, false, func(c *Config) error {
		used := map[string]bool{}
		for _, x := range c.FAQ {
			used[x.Key] = true
		}
		f.Key = unusedFAQKey(used)
		c.FAQ = append(c.FAQ, f)
		return nil
	}, func(tx *sql.Tx, c Config) error {
		return insertFAQ(ctx, tx, len(c.FAQ)-1, f)
	})
	return f.Key, err
}

// AppendKnowledge 補一段知識：title 是現有段落的標題（去掉前後空白比對）就把 body 接在那段後面（換行分隔），
// 否則在最後新增一段；其他段落不動
func AppendKnowledge(ctx context.Context, db *sql.DB, title, body string) error {
	title, body = strings.TrimSpace(title), strings.TrimSpace(body)
	if body == "" {
		return &InvalidError{[]string{"知識的內容是必填"}}
	}
	return saveBlock(ctx, db, false, func(c *Config) error {
		for i, s := range c.Knowledge {
			if strings.TrimSpace(s.Title) == title {
				c.Knowledge[i].Body = strings.TrimRight(s.Body, "\n") + "\n" + body
				return nil
			}
		}
		c.Knowledge = append(c.Knowledge, Section{Title: title, Body: body})
		return nil
	}, func(tx *sql.Tx, c Config) error {
		return replaceKnowledge(ctx, tx, c.Knowledge)
	})
}

// trimAlt 去掉其他問法的前後空白和空白的項目
func trimAlt(alt []string) []string {
	out := []string{}
	for _, q := range alt {
		if q = strings.TrimSpace(q); q != "" {
			out = append(out, q)
		}
	}
	return out
}

// unusedFAQKey 給一個還沒用過的 faq-N（記進 used）
func unusedFAQKey(used map[string]bool) string {
	n := 1
	for used[fmt.Sprintf("faq-%d", n)] {
		n++
	}
	key := fmt.Sprintf("faq-%d", n)
	used[key] = true
	return key
}

// LanguageSettings 是語言區塊：慣用語言和顧客語言（順序照 Languages）
type LanguageSettings struct {
	Preferred string     `json:"preferred_language"`
	Languages []Language `json:"languages"`
}

// LanguagesOf 取出語言區塊
func LanguagesOf(c Config) LanguageSettings {
	return LanguageSettings{Preferred: c.Merchant.PreferredLanguage, Languages: c.Languages}
}

// SaveLanguages 只換掉慣用語言和顧客語言
func SaveLanguages(ctx context.Context, db *sql.DB, s LanguageSettings) error {
	return saveBlock(ctx, db, false, func(c *Config) error {
		c.Merchant.PreferredLanguage = s.Preferred
		c.Languages = s.Languages
		if c.Languages == nil {
			c.Languages = []Language{}
		}
		return nil
	}, func(tx *sql.Tx, c Config) error {
		if err := writeMerchant(ctx, tx, c.Merchant); err != nil {
			return err
		}
		return replaceLanguages(ctx, tx, c.Languages)
	})
}

// FeatureEnabled 是功能開關（booking、handoff）有沒有開；設定裡沒有就當成開著（和對話層一致）
func FeatureEnabled(c Config, name string) bool {
	for _, f := range c.Features {
		if f.Name == name {
			return f.Enabled
		}
	}
	return true
}

// SaveFeature 只改一個功能開關（booking、handoff），沒有那一列就新增
func SaveFeature(ctx context.Context, db *sql.DB, name string, enabled bool) error {
	if name != "booking" && name != "handoff" {
		return &InvalidError{[]string{"不認得的功能開關：" + name}}
	}
	return saveBlock(ctx, db, false, func(c *Config) error { return nil }, func(tx *sql.Tx, c Config) error {
		_, err := tx.ExecContext(ctx, "INSERT OR REPLACE INTO feature_switch (name, enabled) VALUES (?, ?)", name, enabled)
		return err
	})
}

// ErrInfoFirst：還沒有商家設定時，要先存商家資訊（名稱、介紹、聯絡方式必填），FAQ、語言才能存
var ErrInfoFirst = errors.New("請先填商家資訊")

// saveBlock 在同一個 transaction 裡讀出整份設定、套用 apply（區塊自己的檢查沒過時回錯誤）、用整份設定檢查，
// 再用 write 只寫這個區塊的表。還沒有商家設定時：firstOK（商家資訊）從 Default() 開始，整份寫進去；其他區塊回 ErrInfoFirst
func saveBlock(ctx context.Context, db *sql.DB, firstOK bool, apply func(*Config) error, write func(*sql.Tx, Config) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	c, err := load(ctx, tx)
	first := errors.Is(err, ErrNotConfigured)
	if first && !firstOK {
		return ErrInfoFirst
	} else if first {
		c = Default()
	} else if err != nil {
		return err
	}
	if err := apply(&c); err != nil {
		return err
	}
	if err := c.Validate(); err != nil {
		return err
	}
	if first {
		_, err = writeAll(ctx, tx, c)
	} else {
		err = write(tx, c)
	}
	if err != nil {
		return err
	}
	return tx.Commit()
}

// writeMerchant 寫入（或取代）商家那一列
func writeMerchant(ctx context.Context, tx *sql.Tx, m Merchant) error {
	_, err := tx.ExecContext(ctx, "INSERT OR REPLACE INTO merchant (id, name, intro, contact, preferred_language) VALUES (1, ?, ?, ?, ?)",
		m.Name, m.Intro, m.Contact, m.PreferredLanguage)
	return err
}

func replaceKnowledge(ctx context.Context, tx *sql.Tx, sections []Section) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM knowledge_section"); err != nil {
		return err
	}
	for i, s := range sections {
		if _, err := tx.ExecContext(ctx, "INSERT INTO knowledge_section (position, title, body) VALUES (?, ?, ?)", i, s.Title, s.Body); err != nil {
			return err
		}
	}
	return nil
}

func replaceFAQ(ctx context.Context, tx *sql.Tx, faq []FAQ) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM faq"); err != nil {
		return err
	}
	for i, f := range faq {
		if err := insertFAQ(ctx, tx, i, f); err != nil {
			return err
		}
	}
	return nil
}

func insertFAQ(ctx context.Context, tx *sql.Tx, position int, f FAQ) error {
	alt := f.AltQuestions
	if alt == nil {
		alt = []string{}
	}
	altJSON, err := json.Marshal(alt)
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, "INSERT INTO faq (position, key, question, alt_questions, answer) VALUES (?, ?, ?, ?, ?)",
		position, f.Key, f.Question, string(altJSON), f.Answer)
	return err
}

func replaceLanguages(ctx context.Context, tx *sql.Tx, langs []Language) error {
	if _, err := tx.ExecContext(ctx, "DELETE FROM language"); err != nil {
		return err
	}
	for i, l := range langs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO language (code, position, enabled) VALUES (?, ?, ?)", l.Code, i, l.Enabled); err != nil {
			return err
		}
	}
	return nil
}
