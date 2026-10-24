package service

import (
	"fmt"
	"testing"
)

func TestGamePlay_PlaceBet(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	// var cer string = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTkwOSwiTG9naW5OYW1lIjoiZGVtby0xNjEyMDk4Njg0NjEwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjEyOTYyNjg5LCJSbmRJRCI6IjljMGUxZWUzLWRhZDktNDY3Yy04OGJlLWU5ZjY3NWExNmY4ZiIsIlNlc3Npb25Ub2tlbiI6ImVhMzUxZWJhLTIyMmMtNGI5Zi04MzU2LTIyMzMwOGFiYWI3MyIsIlNpZ25LZXkiOiIyNmQyNDNhNWQzODE0OTkyNTZiMzYzOGM5OTYyOTAyMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj+y/j+cZglfyueeAlO6OxOVhT5mlcBNYUjPJHeg60m/udJBYAz9rUFGBfa+t5rSIS5vjhjHWk5sXzOq6rv9cjzo4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+xwkImZSqPczNB3l7uxmVIS3JDJ/47U7Xl01XIll17nXxiCtuHw4ruATjCazNaKSf1hho6vFbQckz/HoH0l+HgD8TWQRDz9s30tn59nuG2wsyxI4YU9tq+oVz0O3hicsWvSuyoMkAM0l5+7pNqyjuqgkh4ceHqHLvIiRB9yr/eL4wr7M/poYPK+XnTKCqHddlVbU0F/V5XFaVN9FcLc0Em8YV9Eyh6Kaq3+JbASZcdSSnoFYBs89G74KaqjlA4HDb6KvdszgaotF7hEBi8F53nyiAT1DOpVEk+SECSfzEHI794MaLtMAeE+rtw72hrpMsEE+/JGqq2tacREC1f3g5Di"
	pid, cer := GetTestAccessToken()
	fmt.Println(pid, cer)

	var servID int32 = 102100410
	var tradeID int64 = 20191200000095331
	var betAmount int64 = 0
	var ipAddr string = "0.0.0.0"
	var memo string = "test"

	ok, afterAmount, transactionNo := GamePlayService.PostPlaceBet(cer, servID, tradeID, betAmount, ipAddr, memo)

	fmt.Println(ok, afterAmount, transactionNo)
}

func TestGamePlay_GetWallet(t *testing.T) {

	if testing.Short() {
		t.Skip()
	}

	var cer string = "eyJhbGciOiJBRVMiLCJ0eXAiOiJKV1QifQ==.eyJMb2dpblVzZXJJRCI6MTAxNTkwOSwiTG9naW5OYW1lIjoiZGVtby0xNjEyMDk4Njg0NjEwIiwiTWVyY2hhbnRJRCI6MTc1MDEzLCJMYW5nIjoiRU4tVVMiLCJUb2tlbkV4cGlyZWRUaW1lIjoxNjEyOTYyNjg5LCJSbmRJRCI6IjljMGUxZWUzLWRhZDktNDY3Yy04OGJlLWU5ZjY3NWExNmY4ZiIsIlNlc3Npb25Ub2tlbiI6ImVhMzUxZWJhLTIyMmMtNGI5Zi04MzU2LTIyMzMwOGFiYWI3MyIsIlNpZ25LZXkiOiIyNmQyNDNhNWQzODE0OTkyNTZiMzYzOGM5OTYyOTAyMyIsIlRva2VuVmVyc2lvbiI6MTAwMH0=.Ob3fFSoNyH1NPXWFYDOOUe+bTDFIZ1dC3fPEXelRt9A+T/s1+7siu8S0oOir8NZj+y/j+cZglfyueeAlO6OxOVhT5mlcBNYUjPJHeg60m/udJBYAz9rUFGBfa+t5rSIS5vjhjHWk5sXzOq6rv9cjzo4ljFliNkOusCJOjR0dFfpW7xlW5B7PyaX0XqT8jOzFfTWXLbg/sD/YD2X6cG/Op2jCFajdZ3/pI0VFhYvnD+xwkImZSqPczNB3l7uxmVIS3JDJ/47U7Xl01XIll17nXxiCtuHw4ruATjCazNaKSf1hho6vFbQckz/HoH0l+HgD8TWQRDz9s30tn59nuG2wsyxI4YU9tq+oVz0O3hicsWvSuyoMkAM0l5+7pNqyjuqgkh4ceHqHLvIiRB9yr/eL4wr7M/poYPK+XnTKCqHddlVbU0F/V5XFaVN9FcLc0Em8YV9Eyh6Kaq3+JbASZcdSSnoFYBs89G74KaqjlA4HDb6KvdszgaotF7hEBi8F53nyiAT1DOpVEk+SECSfzEHI794MaLtMAeE+rtw72hrpMsEE+/JGqq2tacREC1f3g5Di"

	ok, amount, merchantWalletMode := GamePlayService.GetWallet(cer)

	fmt.Println(ok, amount, merchantWalletMode)
}
