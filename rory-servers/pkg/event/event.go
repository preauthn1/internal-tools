package event

import (
	"fmt"
)

// Args 事件參數
type Args map[string]interface{}

// NewEventArgs 創建事件參數
func NewEventArgs() Args {
	args := make(map[string]interface{})

	return args
}

// Ret 事件結果
type Ret struct {
	Ret   map[string]interface{}
	Error interface{}
}

// NewEventRet 創建事件結果
func NewEventRet() *Ret {
	ret := &Ret{
		Ret:   make(map[string]interface{}),
		Error: nil}

	return ret
}

// HasError 結果是否有錯誤
func (ret *Ret) HasError() bool {
	return ret.Error != nil
}

// IsEmpty 是否空的結果
func (ret *Ret) IsEmpty() bool {
	return ret.Ret == nil || len(ret.Ret) == 0
}

func (ret *Ret) mustNotEmpty() {
	if ret.Ret == nil {
		panic("結果集為空")
	}
}

// GetString 從結果中獲取字符串
func (ret *Ret) GetString(key string) string {
	ret.mustNotEmpty()

	return ret.Ret[key].(string)
}

// GetInt 從結果中獲取int
func (ret *Ret) GetInt(key string) int {
	ret.mustNotEmpty()

	return ret.Ret[key].(int)
}

// GetInt64 從結果中獲取int64
func (ret *Ret) GetInt64(key string) int64 {
	ret.mustNotEmpty()

	return ret.Ret[key].(int64)
}

// GetFloat32 從結果中獲取float32
func (ret *Ret) GetFloat32(key string) float32 {
	ret.mustNotEmpty()

	return ret.Ret[key].(float32)
}

// GetFloat64 從結果中獲取float64
func (ret *Ret) GetFloat64(key string) float64 {
	ret.mustNotEmpty()

	return ret.Ret[key].(float64)
}

// GetBool 獲取bool值
func (ret *Ret) GetBool(key string) bool {
	ret.mustNotEmpty()

	return ret.Ret[key].(bool)
}

// SetResult 設置結果
func (ret *Ret) SetResult(key string, value interface{}) {
	ret.mustNotEmpty()

	ret.Ret[key] = value
}

// Single ...
const (
	SINGLERESULTKEYMAME = "result"
	SINGLEARGSKEYNAME   = "arg"
)

// SetSingleResult 設置單一結果，用於返回結果只有一個值的情況
func (ret *Ret) SetSingleResult(value interface{}) {
	ret.SetResult(SINGLERESULTKEYMAME, value)
}

// GetSingleResult 獲取單一結果
func (ret *Ret) GetSingleResult() interface{} {
	return ret.Ret[SINGLERESULTKEYMAME]
}

// ToString 結果轉變為字符串
func (ret *Ret) ToString() string {
	if ret.HasError() {
		return fmt.Sprintf("%v", ret.Error)
	}

	retstr := ""
	for _, v := range ret.Ret {
		retstr += v.(string) + "."
	}

	return retstr
}

// Event 事件
type Event struct {
	Name string // 事件名稱
	Args Args   // 事件參數
}

// NewEvent 創建事件對象
func NewEvent(name string, args Args) *Event {
	evt := new(Event)
	evt.Name = name
	evt.Args = args

	return evt
}

// GetSingleArgValue 取得單一平均值
func (evt *Event) GetSingleArgValue() interface{} {
	return evt.Args[SINGLEARGSKEYNAME]
}

// GetString 獲取string參數
func (evt *Event) GetString(key string) string {
	return evt.Args[key].(string)
}

// GetInt 獲取int參數
func (evt *Event) GetInt(key string) int {
	return evt.Args[key].(int)
}

// GetInt64 獲取int64參數
func (evt *Event) GetInt64(key string) int64 {
	return evt.Args[key].(int64)
}

// GetFloat32 獲取float32參數
func (evt *Event) GetFloat32(key string) float32 {
	return evt.Args[key].(float32)
}

// GetFloat64 獲取float64參數
func (evt *Event) GetFloat64(key string) float64 {
	return evt.Args[key].(float64)
}

// GetBool 獲取bool參數
func (evt *Event) GetBool(key string) bool {
	return evt.Args[key].(bool)
}
