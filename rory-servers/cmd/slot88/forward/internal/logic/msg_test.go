package logic

import (
	"testing"

	"example.invalid/ServerTeam/Servers/cmd/slot88/forward/internal/encrypt"
)

func TestDecode(t *testing.T) {

	aes, _ := encrypt.NewAES("fa7064003411ac6587bdeb60441e5f45")

	b := aes.Decrypt([]byte("e36181f37befe354f8d66ea4793ec5a1cb920b551d074d78ad64eece9ed5d60c07900804bf5a5cf4688ff8c67a7436d067914e4a2c6ed910a9710bc8fe97eae9567eba63d965e35bf8a63863be0ac5d9"))

	t.Log(string(b))
}
