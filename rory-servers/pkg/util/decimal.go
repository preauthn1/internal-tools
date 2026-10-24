package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Decimal ?
func Decimal(value float64, n int) float32 {
	var buffer bytes.Buffer
	buffer.WriteString("%.")
	buffer.WriteString(strconv.Itoa(n))
	buffer.WriteString("f")

	value, _ = strconv.ParseFloat(fmt.Sprintf(buffer.String(), value), 64)
	return float32(value)
}

// MoneyToCent ?
func MoneyToCent(money json.Number) int64 {
	foo := strings.Split(money.String(), ".")
	if len(foo) == 1 {
		v, _ := strconv.ParseInt(foo[0]+"00", 10, 64)
		return v
	}
	if len(foo) == 2 {
		v, _ := strconv.ParseInt(foo[0]+PadRight(foo[1], "0", 2), 10, 64)
		return v
	}
	panic("json.Number is not validate")
}

// CentToMoney ?
func CentToMoney(cent int64) json.Number {
	return json.Number(fmt.Sprintf("%.2f", float64(cent)/100))
}

// ABSInt64 回傳絕對值
func ABSInt64(n int64) int64 {
	if n < 0 {
		return -n
	}
	return n
}
