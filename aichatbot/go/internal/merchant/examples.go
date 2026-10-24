package merchant

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// StoredExample 是資料庫裡的一句意圖案例；ID 給後台刪除用
type StoredExample struct {
	ID      int64  `json:"id"`
	Text    string `json:"text"`
	Intent  string `json:"intent"`
	State   string `json:"state,omitempty"`
	Builtin bool   `json:"builtin"`
}

// ErrExampleNotFound：要刪除的案例不存在
var ErrExampleNotFound = errors.New("找不到這個案例")

// MaxExampleLength 是一句案例最多幾個字
const MaxExampleLength = 200

// ErrInvalidExample：句子或意圖不符合規定
var ErrInvalidExample = errors.New("案例不正確")

// ErrDuplicateExample：案例庫裡已經有比對後相同的句子（不管是不是內建）
var ErrDuplicateExample = errors.New("案例庫裡已經有這句")

// OverlapError：新增的句子和題庫重複；Sets 是哪些題庫，Text 是哪一句（一次存好幾句時才填，例如自訂意圖的例句）
type OverlapError struct {
	Sets []string
	Text string
}

func (e *OverlapError) Error() string {
	msg := "和題庫重複（" + strings.Join(e.Sets, "、") + "），請換一種說法"
	if e.Text != "" {
		msg = "「" + e.Text + "」" + msg
	}
	return msg
}

// AddExample 新增一句商家的案例（非內建）；和題庫重複（比對規則見 Normalize）回 *OverlapError
func AddExample(ctx context.Context, db *sql.DB, tests TestSet, text, intent string) (StoredExample, error) {
	e := StoredExample{Text: strings.TrimSpace(text), Intent: intent}
	switch {
	case Normalize(e.Text) == "":
		return StoredExample{}, fmt.Errorf("%w：句子要有文字", ErrInvalidExample)
	case utf8.RuneCountInString(e.Text) > MaxExampleLength:
		return StoredExample{}, fmt.Errorf("%w：句子最多 %d 字", ErrInvalidExample, MaxExampleLength)
	}
	if sets := tests.Find(e.Text); len(sets) > 0 {
		return StoredExample{}, &OverlapError{Sets: sets}
	}
	// 檢查重複和新增放在同一個 transaction：同時送出兩次也不會存成兩句
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return StoredExample{}, err
	}
	defer tx.Rollback()
	// 意圖是內建意圖或已有的自訂意圖
	if _, ok := IntentByCode(intent); !ok {
		var n int
		if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM custom_intent WHERE code = ?", intent).Scan(&n); err != nil {
			return StoredExample{}, err
		}
		if n == 0 {
			return StoredExample{}, fmt.Errorf("%w：沒有這個意圖", ErrInvalidExample)
		}
	}
	existing, err := listExamples(ctx, tx)
	if err != nil {
		return StoredExample{}, err
	}
	// 自訂意圖的例句分類時也一起比對，一樣不能重複
	custom, err := listCustomIntents(ctx, tx)
	if err != nil {
		return StoredExample{}, err
	}
	texts := []string{}
	for _, x := range existing {
		texts = append(texts, x.Text)
	}
	for _, x := range CustomIntentExamples(custom) {
		texts = append(texts, x.Text)
	}
	for _, s := range texts {
		if Normalize(s) == Normalize(e.Text) {
			return StoredExample{}, fmt.Errorf("%w：「%s」", ErrDuplicateExample, s)
		}
	}
	res, err := tx.ExecContext(ctx, "INSERT INTO intent_example (text, intent, builtin) VALUES (?, ?, 0)", e.Text, e.Intent)
	if err != nil {
		return StoredExample{}, err
	}
	if e.ID, err = res.LastInsertId(); err != nil {
		return StoredExample{}, err
	}
	return e, tx.Commit()
}

// ErrBuiltinExample：內建案例不能刪除
var ErrBuiltinExample = errors.New("內建案例不能刪除")

// DeleteExample 刪除一句商家的案例；內建案例回 ErrBuiltinExample，不刪
func DeleteExample(ctx context.Context, db *sql.DB, id int64) error {
	var builtin bool
	err := db.QueryRowContext(ctx, "SELECT builtin FROM intent_example WHERE id = ?", id).Scan(&builtin)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrExampleNotFound
	}
	if err != nil {
		return err
	}
	if builtin {
		return ErrBuiltinExample
	}
	_, err = db.ExecContext(ctx, "DELETE FROM intent_example WHERE id = ? AND builtin = 0", id)
	return err
}

// ListExamples 列出意圖案例庫（內建和商家新增的），依新增的順序
func ListExamples(ctx context.Context, db *sql.DB) ([]StoredExample, error) {
	return listExamples(ctx, db)
}

func listExamples(ctx context.Context, db querier) ([]StoredExample, error) {
	return query(ctx, db, "SELECT id, text, intent, state, builtin FROM intent_example ORDER BY id", func(r *sql.Rows) (e StoredExample, err error) {
		return e, r.Scan(&e.ID, &e.Text, &e.Intent, &e.State, &e.Builtin)
	})
}
