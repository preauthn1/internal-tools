package encrypt

import "testing"

func TestAES(t *testing.T) {

	token := "13f7c57379fe5e3aec0747483f87f867"

	aes, err := NewAES(token)
	if err != nil {
		t.Error(err)
		return
	}

	b := aes.Encrypt([]byte(token))

	t.Logf("%s\n", b)
}
