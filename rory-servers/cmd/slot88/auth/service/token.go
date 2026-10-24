package service

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"example.invalid/dgrijalva/jwt-go"
	"example.invalid/spf13/viper"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

var tokenStore map[string]string
var mu sync.RWMutex

func init() {
	tokenStore = make(map[string]string)
}

// GenToken create jwt token
func GenToken(i Info) (string, error) {

	now := time.Now()
	expire := now.Add(viper.GetDuration("token.expire"))

	i.StandardClaims = jwt.StandardClaims{
		ExpiresAt: expire.Unix(),
		IssuedAt:  now.Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, i)

	return t.SignedString([]byte(viper.GetString("token.key")))
}

// SaveToken create token md5 key, then save to db
func SaveToken(db *sql.DB, token string) (string, error) {

	code := fmt.Sprintf("%x", md5.Sum([]byte(token)))

	// _, err := db.Exec("INSERT INTO [Slot88].[Token]([key], [token]) VALUES (?, ?)", code, token)
	// if err != nil {
	// 	return "", err
	// }

	mu.Lock()
	defer mu.Unlock()

	tokenStore[code] = token

	return code, nil
}

// GetToken get token by md5 key from db
func GetToken(db *sql.DB, code string) (*Info, error) {

	// row := db.QueryRow("SELECT [token] FROM [Slot88].[Token] WHERE [key] = ?", code)

	// var str string

	// if err := row.Scan(&str); err != nil {
	// 	return nil, err
	// }

	mu.RLock()
	defer mu.RUnlock()

	str := tokenStore[code]

	info := new(Info)
	token, err := jwt.ParseWithClaims(str, info, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(viper.GetString("token.key")), nil
	})
	if err != nil {
		return nil, err
	}

	if i, ok := token.Claims.(*Info); ok && token.Valid {
		return i, nil
	}

	return nil, jwt.NewValidationError("", jwt.ValidationErrorExpired)
}

// CleanExpireToken clean token
func CleanExpireToken(db *sql.DB) (int, error) {

	// rows, err := db.Query("SELECT [key], [token] FROM [Slot88].[Token]")
	// if err != nil {
	// 	return 0, err
	// }
	// defer rows.Close()

	// tx, err := db.Begin()
	// if err != nil {
	// 	return 0, err
	// }

	var cnt int

	mu.Lock()
	defer mu.Unlock()

	for code, token := range tokenStore {
		info := new(Info)
		_, err := jwt.ParseWithClaims(token, info, func(token *jwt.Token) (i interface{}, err error) {
			return []byte(viper.GetString("token.key")), nil
		})
		if err != nil {
			delete(tokenStore, code)
			cnt++
		}
	}

	// for rows.Next() {
	// 	var code, str string

	// 	if err := rows.Scan(&code, &str); err != nil {
	// 		tx.Rollback()
	// 		return 0, err
	// 	}

	// 	info := new(Info)
	// 	_, err := jwt.ParseWithClaims(str, info, func(token *jwt.Token) (i interface{}, err error) {
	// 		return []byte(viper.GetString("token.key")), nil
	// 	})
	// 	if err != nil {
	// 		result, _ := tx.Exec("DELETE FROM [Slot88].[Token] WHERE [key] = ? ", code)
	// 		i, _ := result.RowsAffected()
	// 		cnt += int(i)
	// 	}

	// }

	// tx.Commit()

	mlog.System("Clean Expired Token (%d)", cnt)

	return cnt, nil
}
