package merchant

// 預約欄位（CONTEXT.md「預約欄位」，#67）：日期、時間、姓名、電話固定要問；人數可以關掉，另外可以加自訂欄位

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// BookingFields 是預約欄位的設定。Config 裡是 nil 時用預設（問人數、沒有自訂欄位），和 #67 以前相同
type BookingFields struct {
	PeopleEnabled bool           `json:"people_enabled"`
	Custom        []BookingField `json:"custom"`
}

// BookingField 是一個自訂欄位。Key 是代號（回覆種類 ask_<代號>、意圖案例狀態 booking:<代號>、預約的 extra 都用它），
// Description 寫進給 AI 抽取的欄位說明；選填（Required = false）的不追問，顧客有講才記下來
type BookingField struct {
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// MaxBookingFields 是最多幾個自訂欄位：每多一個，給本機模型抽取的欄位就多一個，太多時小模型抽得慢也抽不準
const MaxBookingFields = 5

// bookingFieldKey 是自訂欄位的代號：小寫英文開頭，小寫英文、數字、底線，最多 20 個字
var bookingFieldKey = regexp.MustCompile(`^[a-z][a-z0-9_]{0,19}$`)

// builtinBookingKeys 是內建欄位和給 AI 抽取用的欄位名稱，自訂欄位的代號不能和它們一樣；
// 另外 ask_<代號> 不能和內建的固定回覆（ask_time_midnight、ask_name_soft）撞名
var builtinBookingKeys = []string{"date", "time", "people", "name", "phone",
	"date_quote", "date_month", "date_day", "time_quote", "time_hour", "time_minute", "people_quote",
	"time_midnight", "name_soft"}

// UnmarshalJSON：匯入檔沒寫 people_enabled 時當成問人數（沒寫不該默默關掉）；不認得的欄位照樣報錯（ParseJSON 的檢查不會傳進這裡）
func (b *BookingFields) UnmarshalJSON(data []byte) error {
	type bookingFields BookingFields
	v := bookingFields{PeopleEnabled: true}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&v); err != nil {
		return err
	}
	*b = BookingFields(v)
	return nil
}

// validateBookingFields 是 Validate 對預約欄位的檢查
func validateBookingFields(c Config) []string {
	if c.BookingFields == nil {
		return nil
	}
	var errs []string
	list := c.BookingFields.Custom
	if len(list) > MaxBookingFields {
		errs = append(errs, fmt.Sprintf("自訂欄位最多 %d 個", MaxBookingFields))
	}
	keys := map[string]bool{}
	for i, f := range list {
		switch {
		case !bookingFieldKey.MatchString(f.Key):
			errs = append(errs, "自訂欄位的代號不正確："+f.Key+"（小寫英文開頭，只能用小寫英文、數字、底線）")
		case slices.Contains(builtinBookingKeys, f.Key):
			errs = append(errs, "自訂欄位的代號不能用內建欄位："+f.Key)
		case keys[f.Key]:
			errs = append(errs, "自訂欄位的代號重複："+f.Key)
		}
		keys[f.Key] = true
		for _, x := range []struct{ name, v string }{{"名稱", f.Name}, {"說明", f.Description}} {
			if strings.TrimSpace(x.v) == "" {
				errs = append(errs, fmt.Sprintf("第 %d 個自訂欄位的%s是必填", i+1, x.name))
			}
		}
	}
	return errs
}

// EffectiveBookingFields 是實際用的預約欄位：沒設定時是預設（問人數、沒有自訂欄位）
func EffectiveBookingFields(c Config) BookingFields {
	if c.BookingFields == nil {
		return BookingFields{PeopleEnabled: true, Custom: []BookingField{}}
	}
	b := *c.BookingFields
	if b.Custom == nil {
		b.Custom = []BookingField{}
	}
	return b
}

// isDefault：和沒設定時相同（存檔、讀出時都當成沒設定，匯出檔不多出 booking_fields）
func (b *BookingFields) isDefault() bool {
	return b == nil || (b.PeopleEnabled && len(b.Custom) == 0)
}

// SaveBookingFields 只換掉預約欄位（後台的預約設定頁）；名稱、說明去掉前後空白。
// 沒通過檢查回 *InvalidError，還沒有商家設定時回 ErrInfoFirst
func SaveBookingFields(ctx context.Context, db *sql.DB, b BookingFields) error {
	custom := make([]BookingField, len(b.Custom))
	for i, f := range b.Custom {
		f.Name, f.Description = strings.TrimSpace(f.Name), strings.TrimSpace(f.Description)
		custom[i] = f
	}
	b.Custom = custom
	return saveBlock(ctx, db, false, func(c *Config) error { c.BookingFields = &b; return nil }, func(tx *sql.Tx, c Config) error {
		return replaceBookingFields(ctx, tx, c.BookingFields)
	})
}

// replaceBookingFields 用 b 取代預約欄位；預設的設定不存（兩張表都清空）
func replaceBookingFields(ctx context.Context, tx *sql.Tx, b *BookingFields) error {
	for _, table := range []string{"booking_setting", "booking_field"} {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+table); err != nil {
			return err
		}
	}
	if b.isDefault() {
		return nil
	}
	if _, err := tx.ExecContext(ctx, "INSERT INTO booking_setting (id, people_enabled) VALUES (1, ?)", b.PeopleEnabled); err != nil {
		return err
	}
	for i, f := range b.Custom {
		if _, err := tx.ExecContext(ctx, "INSERT INTO booking_field (position, key, name, description, required) VALUES (?, ?, ?, ?, ?)",
			i, f.Key, f.Name, f.Description, f.Required); err != nil {
			return err
		}
	}
	return nil
}

// loadBookingFields 讀出預約欄位；沒存過（或存的是預設）回 nil
func loadBookingFields(ctx context.Context, db querier) (*BookingFields, error) {
	b := &BookingFields{PeopleEnabled: true}
	err := db.QueryRowContext(ctx, "SELECT people_enabled FROM booking_setting WHERE id = 1").Scan(&b.PeopleEnabled)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if b.Custom, err = query(ctx, db, "SELECT key, name, description, required FROM booking_field ORDER BY position, id", func(r *sql.Rows) (f BookingField, err error) {
		return f, r.Scan(&f.Key, &f.Name, &f.Description, &f.Required)
	}); err != nil {
		return nil, err
	}
	if b.isDefault() {
		return nil, nil
	}
	return b, nil
}
