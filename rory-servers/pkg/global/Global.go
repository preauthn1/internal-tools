package global

import "fmt"

// MoneyString 將金額轉成 client 顯示的格式字串
func MoneyString(i64Money int64) string {
	strMoney := ""

	// 整數部分
	strSign := ""
	if i64Money < 0 {
		strSign = "-"
	}
	i64Integer := i64Money / 100

	// 小數部分
	i64Decimal := i64Money % 100
	if i64Decimal < 0 {
		i64Decimal = -i64Decimal
	}

	if i64Integer == 0 {
		strMoney = strSign + fmt.Sprintf("%d", i64Integer) + "." + fmt.Sprintf("%02d", i64Decimal)
	} else {
		strMoney = fmt.Sprintf("%d", i64Integer) + "." + fmt.Sprintf("%02d", i64Decimal)
	}

	return strMoney
}
