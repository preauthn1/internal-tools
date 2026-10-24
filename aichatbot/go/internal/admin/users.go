// Package admin 是後台：管理員帳號、登入 session、後台 API，以及提供建好的後台頁面（web/apps/web-antd/dist）。
// 只有一種管理員角色；除了登入，/api/ 底下一律要先登入。
package admin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"
	"unicode/utf8"

	"example.invalid/x/crypto/bcrypt"

	"aichatbot/internal/store"
)

// MinPasswordLength 是密碼最少的字數（以字元算，不是 byte）
const MinPasswordLength = 10

var usernamePattern = regexp.MustCompile(`^[A-Za-z0-9_-]{3,32}$`)

// ErrUserExists：帳號已經有了（create-admin 不覆蓋原本的帳號）
var ErrUserExists = errors.New("帳號已經存在")

// CreateUser 建立管理員；密碼存 bcrypt（預設 cost）
func CreateUser(ctx context.Context, db *store.DB, username, password string) error {
	if !usernamePattern.MatchString(username) {
		return errors.New("帳號要是 3～32 個英數字（可以用 _、-）")
	}
	if utf8.RuneCountInString(password) < MinPasswordLength {
		return fmt.Errorf("密碼至少要 %d 個字", MinPasswordLength)
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	res, err := db.SQL().ExecContext(ctx, `
		INSERT INTO admin_user (username, password_hash, created_at) VALUES (?, ?, ?)
		ON CONFLICT (username) DO NOTHING`, username, string(hash), store.Timestamp(time.Now()))
	if err != nil {
		return err
	}
	if n, err := res.RowsAffected(); err != nil {
		return err
	} else if n == 0 {
		return ErrUserExists
	}
	return nil
}

// User 是登入的管理員
type User struct {
	ID       int64
	Username string
}

// errBadLogin：帳號不存在或密碼錯（不分開說，不透露帳號存不存在）
var errBadLogin = errors.New("帳號或密碼錯誤")

// dummyHash 給帳號不存在時比對用，讓回應時間和密碼錯一樣
var dummyHash, _ = bcrypt.GenerateFromPassword([]byte("aichatbot dummy password"), bcrypt.DefaultCost)

// authenticate 檢查帳號密碼，錯了回 errBadLogin
func authenticate(ctx context.Context, db *store.DB, username, password string) (User, error) {
	u := User{Username: username}
	var hash string
	err := db.SQL().QueryRowContext(ctx, "SELECT id, password_hash FROM admin_user WHERE username = ?", username).Scan(&u.ID, &hash)
	if errors.Is(err, sql.ErrNoRows) {
		bcrypt.CompareHashAndPassword(dummyHash, []byte(password))
		return User{}, errBadLogin
	}
	if err != nil {
		return User{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) != nil {
		return User{}, errBadLogin
	}
	return u, nil
}
