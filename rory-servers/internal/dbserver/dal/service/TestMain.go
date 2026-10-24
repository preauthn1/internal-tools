package service

import (
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"testing"
)

func newBase() *Base{
	base := &Base{}
	base.gameAPI = "0.0.0.0:3009"
	base.apiCertificate = ""
	return base
}

// TestMain 初始化測試
func TestMain(m *testing.M) {
	httputil.DebugMode(true)
	m.Run()
}
