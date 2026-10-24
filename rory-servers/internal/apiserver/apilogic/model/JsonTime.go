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

// JSONTime ?
type JSONTime struct {
	time.Time
}

// MarshalJSON ?
func (jt JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", jt.UTC().Format(time.RFC3339Nano))), nil
	// return []byte(fmt.Sprintf("\"%s\"", time.Time(jt.Time).Format(jsonTimeFormart))), nil
}

// UnmarshalJSON ?
func (jt *JSONTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(fmt.Sprintf("\"%s\"", jsonTimeFormart), string(data))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	jt.Time = t.Local()
	return nil
}

func (jt *JSONTime) String() string {
	return jt.Format(jsonTimeFormart)
	// return util.FormatTimeStr(jt.Time)
}
