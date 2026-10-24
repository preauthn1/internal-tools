package service

import "testing"

// TestMain ?
func TestMain(m *testing.M) {
	InitTestService()
	m.Run()
}
