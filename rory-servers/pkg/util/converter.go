package util

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Converter utilites

// ToJSON convert interface{} -> json
func ToJSON(m interface{}) string {
	js, err := json.Marshal(m)
	if err != nil {
		// return "", err
		return ""
	}
	return string(js)
}

// ToMap convert json -> map[string]interface{}
func ToMap(b []byte) map[string]interface{} {
	var err error
	var jsonMap = make(map[string]interface{})
	err = json.Unmarshal(b, &jsonMap)
	if err != nil {
		return nil
	}

	return jsonMap
}

// ToMapMap convert json -> map[string]map[string]interface{}
func ToMapMap(b []byte) map[string]map[string]interface{} {
	var err error
	var jsonMap = make(map[string]map[string]interface{})
	err = json.Unmarshal(b, &jsonMap)
	if err != nil {
		return nil
	}

	return jsonMap
}

// ToInterface convert json -> interface{}
func ToInterface(b []byte) interface{} {
	var err error
	var data interface{}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil
	}
	return data
}

// ToInt convert interface{} -> int, if failed return default value
func ToInt(m interface{}, defaultValue ...int) int {
	val, err := strconv.ParseInt(fmt.Sprintf("%v", m), 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return int(val)
}

// ToInt32 convert interface{} -> int32, if failed return default value
func ToInt32(m interface{}, defaultValue ...int32) int32 {
	val, err := strconv.ParseInt(fmt.Sprintf("%v", m), 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return int32(val)
}

// ToInt64 convert interface{} -> int64, if failed return default value
func ToInt64(m interface{}, defaultValue ...int64) int64 {
	val, err := strconv.ParseInt(fmt.Sprintf("%v", m), 10, 64)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		}
		return 0
	}

	return val
}

// ToString convert interface{} -> string, if failed return default value
func ToString(m interface{}, defaultValue ...string) string {
	if val, ok := m.(string); ok {
		return val
	}
	if len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return ""
}
