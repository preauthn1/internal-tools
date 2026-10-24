package encrypt

// PKCS7Encode 確保內容為16個byte為基數, 若不足會補齊到16的倍數, padder為還差n個
func PKCS7Encode(src []byte) (dst []byte) {
	l := len(src)

	padder := 16 - l%16

	dst = make([]byte, l+padder)

	for i := 0; i < l; i++ {
		dst[i] = src[i]
	}

	for i := 0; i < padder; i++ {
		dst[i+l] = byte(padder)
	}

	return
}

// PKCS7Decode 取最後一位當作padder, 把最後n個byte移除
func PKCS7Decode(src []byte) (dst []byte) {
	l := len(src)

	padder := int(src[len(src)-1])

	if padder <= 16 {
		l -= padder
	}

	dst = make([]byte, l)

	for i := 0; i < l; i++ {
		dst[i] = src[i]
	}

	return
}
