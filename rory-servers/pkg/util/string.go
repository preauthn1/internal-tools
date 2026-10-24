package util

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TrimAllSpace 刪除字符串中的所有空格符
func TrimAllSpace(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\t", "", -1)

	return str
}

// TrimEnter 刪除換行符
func TrimEnter(str string) string {
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)

	return str
}

// IsNumber 是否正式
func IsNumber(str string) bool {
	regx := regexp.MustCompile("^\\d+$")
	return regx.Match([]byte(str))
}

// IsNullOrEmpty 是否空字符
func IsNullOrEmpty(str string) bool {
	str = TrimAllSpace(str)
	str = TrimEnter(str)

	return str == ""
}

// FormatTimeStr 格式化事件為標準的字符串
func FormatTimeStr(t time.Time) string {
	return t.Format(time.RFC3339)
}

// StrToInt64 字串轉數字
func StrToInt64(s string) int64 {
	i32, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return i32
}

// StrToInt32 ?
func StrToInt32(s string) int32 {
	return int32(StrToInt64(s))
}

// PadRight ?
func PadRight(str, pad string, lenght int) string {
	for {
		str += pad
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

// PadLeft ?
func PadLeft(str, pad string, lenght int) string {
	for {
		str = pad + str
		if len(str) > lenght {
			return str[0:lenght]
		}
	}
}

// AesDecryptString target 解密
func AesDecryptString(ciphertext string, pkey string, piv string) (string, error) {
	key := []byte(pkey)
	cipherhex, _ := hex.DecodeString(ciphertext)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := []byte(piv)
	//ciphertext = ciphertext[aes.BlockSize:]

	// CBC mode always works in whole blocks.
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	// 這個驗證是copy mode.CryptBlocks 內部panic 的判斷, 拿來外面先校驗.
	if len(cipherhex)%mode.BlockSize() != 0 {
		return "", errors.New("input not full blocks")
		//panic("crypto/cipher: input not full blocks")
	}

	// CryptBlocks can work in-place if the two arguments are the same.
	mode.CryptBlocks(cipherhex, cipherhex)
	cipherhex, _ = Unpad(cipherhex)
	// If the original plaintext lengths are not a multiple of the block
	// size, padding would have to be added when encrypting, which would be
	// removed at this point. For an example, see
	// https://example.invalid/html/rfc5246#section-0.0.0.0. However, it's
	// critical to note that ciphertexts must be authenticated (i.e. by
	// using crypto/hmac) before being decrypted in order to avoid creating
	// a padding oracle.
	return string(cipherhex), nil
}

// AesEncryptString ?
func AesEncryptString(text string, pkey string, piv string) (string, error) {
	key := []byte(pkey)
	plaintext := []byte(text)
	// auto pad
	plaintext = Pad(plaintext)
	// CBC mode works on blocks so plaintexts may need to be padded to the
	// next whole block. For an example of such padding, see
	// https://example.invalid/html/rfc5246#section-0.0.0.0. Here we'll
	// assume that the plaintext is already of the correct length.
	if len(plaintext)%aes.BlockSize != 0 {
		return "", errors.New("plaintext is not a multiple of the block size")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, len(plaintext))
	iv := []byte(piv)
	//if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	//	return "", err
	//}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)

	// It's important to remember that ciphertexts must be authenticated
	// (i.e. by using crypto/hmac) as well as being encrypted in order to
	// be secure.
	return hex.EncodeToString(ciphertext), nil
}

// Unpad ?
func Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpad error. This could happen when incorrect encryption key is used")
	}

	return src[:(length - unpadding)], nil
}

// Pad ?
var Pad = func(src []byte) []byte {
	padding := aes.BlockSize - len(src)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}
