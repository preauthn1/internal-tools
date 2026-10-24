package util

import (
	"encoding/json"
	"fmt"
	"example.invalid/stretchr/testify/assert"
	"testing"
)

func TestMoneyToCent(t *testing.T) {

	assert.Equal(t, MoneyToCent(json.Number("0.01")), int64(1), "should be equal 1")
	assert.Equal(t, MoneyToCent(json.Number("0.11")), int64(11), "should be equal 11")
	assert.Equal(t, MoneyToCent(json.Number("1.11")), int64(111), "should be equal 111")
	assert.Equal(t, MoneyToCent(json.Number("0.009")), int64(0), "should be equal 1")
	assert.Equal(t, MoneyToCent(json.Number("0.1")), int64(10), "should be equal 10")
	assert.Equal(t, MoneyToCent(json.Number("0.29")), int64(29), "should be equal 29")
	assert.Equal(t, MoneyToCent(json.Number("0.58")), int64(58), "should be equal 58")
}

func TestCentToMoney(t *testing.T) {

	assert.Equal(t, CentToMoney(1).String(), json.Number("0.01").String(), "should be equal 0.01")
	assert.Equal(t, CentToMoney(11).String(), json.Number("0.11").String(), "should be equal 0.11")
	assert.Equal(t, CentToMoney(111).String(), json.Number("1.11").String(), "should be equal 1.11")
	assert.Equal(t, CentToMoney(0).String(), json.Number("0.00").String(), "should be equal 0.00")
	assert.Equal(t, CentToMoney(110).String(), json.Number("1.10").String(), "should be equal 1.10")
	assert.Equal(t, CentToMoney(10).String(), json.Number("0.10").String(), "should be equal 0.10")
	assert.Equal(t, CentToMoney(99).String(), json.Number("0.99").String(), "should be equal 0.99")
}

// MoneyToCent
func TestMoneyToCent_001To100(t *testing.T) {
	for i := 0; i <= 100; i++ {
		var origin = float64(0.01) * float64(i)

		var int64Cent = MoneyToCent(json.Number(fmt.Sprintf("%v", origin)))
		var int64Value = int64(i)

		assert.Equal(t, int64Cent, int64Value, fmt.Sprintf("MoneyToCent [%v] %v, %v, %v, ", i, origin, int64Cent, int64Value))
	}
}

func TestCentToMoney_1To100(t *testing.T) {
	for i := 0; i <= 100; i++ {
		var origin = int64(i)

		var float64Value = float64(0.01) * float64(i)
		var value = json.Number(fmt.Sprintf("%.2f", float64Value)).String()
		var money = CentToMoney(origin).String()

		assert.Equal(t, value, money, fmt.Sprintf("CentToMoney [%v] %v, %s, %s, ", i, origin, value, money))
	}
}
