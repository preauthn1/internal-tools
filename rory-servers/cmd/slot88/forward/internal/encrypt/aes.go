package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
)

// AES slot88 專用加解密工具
type AES struct {
	cipher cipher.Block
	token  string
}

// NewAES 建立 AES 加解密工具,
func NewAES(token string) (*AES, error) {
	key, err := hex.DecodeString(token)
	if err != nil {
		return nil, err
	}

	cipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aes := new(AES)
	aes.cipher = cipher
	aes.token = [REDACTED]

	return aes, nil
}

// Encrypt 加密
func (aes *AES) Encrypt(b []byte) []byte {

	defer func() {
		if r := recover(); r != nil {
			mlog.Error("AES Encrypt Error, %v", r)
		}
	}()

	b = PKCS7Encode(b)

	dst := make([]byte, len(b))
	for bs, be := 0, aes.cipher.BlockSize(); bs < len(b); bs, be = bs+aes.cipher.BlockSize(), be+aes.cipher.BlockSize() {
		aes.cipher.Encrypt(dst[bs:be], b[bs:be])
	}

	_dst := hex.EncodeToString(dst)

	return []byte(_dst)
}

// Decrypt 解密
func (aes *AES) Decrypt(b []byte) []byte {

	defer func() {
		if r := recover(); r != nil {
			mlog.Error("AES Decrypt Error, %v", r)
		}
	}()

	b, _ = hex.DecodeString(string(b))

	dst := make([]byte, len(b))
	for bs, be := 0, aes.cipher.BlockSize(); bs < len(b); bs, be = bs+aes.cipher.BlockSize(), be+aes.cipher.BlockSize() {
		aes.cipher.Decrypt(dst[bs:be], b[bs:be])
	}

	dst = PKCS7Decode(dst)

	return []byte(dst)
}

// Token 取得原Token
func (aes *AES) Token() string {
	return aes.token
}
