package service

import (
	"example.invalid/dgrijalva/jwt-go"
)

// AccessToken 登入資訊物件
type AccessToken struct {
	LoginUserID      int64
	LoginName        string
	MerchantID       int64
	Lang             string
	TokenExpiredTime int64
	RndID            string
	SignKey          string
	TokenVersion     int64
	jwt.StandardClaims
}

// {
// 	"LoginUserID": 1013025,
// 	"LoginName": "demo-1598584773155",
// 	"MerchantID": 175013,
// 	"Lang": "EN-US",
// 	"TokenExpiredTime": 1599448780,
// 	"RndID": "b276499f-4ace-4dad-936e-04ee4618527d",
// 	"SignKey": "012d8bc4f0a6e6ccbf017b95d839d0a7",
// 	"TokenVersion": 1000
//   }

// ParseAccessToken 解析 jwt AccessToken
func ParseAccessToken(cer string) *AccessToken {
	if len(cer) == 0 {
		return nil
	}
	at := new(AccessToken)
	token, _ := jwt.ParseWithClaims(cer, at, nil)
	if token == nil {
		return nil
	}
	if i, ok := token.Claims.(*AccessToken); ok {
		return i
	}
	return nil
}
