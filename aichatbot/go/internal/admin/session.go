package admin

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"time"

	"aichatbot/internal/store"
)

// SessionTTL：session 從登入算起 7 天，使用時不延長
const SessionTTL = 7 * 24 * time.Hour

// newSession 建一個 session，回傳給 cookie 的 token（32 bytes 隨機值）；資料庫只存 token 的 SHA-256
func newSession(ctx context.Context, db *store.DB, userID int64, now time.Time) (string, error) {
	b := make([]byte, 32)
	rand.Read(b)
	token := base64.RawURLEncoding.EncodeToString(b)
	_, err := db.SQL().ExecContext(ctx, "INSERT INTO admin_session (token_hash, user_id, created_at, expires_at) VALUES (?, ?, ?, ?)",
		hashToken(token), userID, store.Timestamp(now), store.Timestamp(now.Add(SessionTTL)))
	if err != nil {
		return "", err
	}
	return token, nil
}

// sessionUser 找出 token 對應、還沒過期的 session 的管理員；找不到回 sql.ErrNoRows
func sessionUser(ctx context.Context, db *store.DB, token string, now time.Time) (User, error) {
	var u User
	err := db.SQL().QueryRowContext(ctx, `
		SELECT u.id, u.username FROM admin_session s JOIN admin_user u ON u.id = s.user_id
		WHERE s.token_hash = ? AND s.expires_at > ?`, hashToken(token), store.Timestamp(now)).Scan(&u.ID, &u.Username)
	return u, err
}

func deleteSession(ctx context.Context, db *store.DB, token string) error {
	_, err := db.SQL().ExecContext(ctx, "DELETE FROM admin_session WHERE token_hash = ?", hashToken(token))
	return err
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}
