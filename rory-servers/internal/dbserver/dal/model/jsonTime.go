package model

import (
	"fmt"
	"time"
)

const (
	// jsonTimeFormart = "2006-01-02T15:04:05"
	// jsonTimeFormart = "2006-01-02T15:04:05.00"
	// jsonTimeFormart = "2006-01-02T15:04:05.000"
	jsonTimeFormart = time.RFC3339
)

// JSONTime JSON 時間型別
type JSONTime struct {
	time.Time
}

// MarshalJSON 編碼
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(jt.Time).Format(jsonTimeFormart)), nil
}

// UnmarshalJSON 解碼
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(fmt.Sprintf("\"%s\"", jsonTimeFormart), string(data))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	jt.Time = t
	return nil
}

// String 返回字串
func (jt *JSONTime) String() string {
	return jt.Format(jsonTimeFormart)
	// return util.FormatTimeStr(jt.Time)
}
