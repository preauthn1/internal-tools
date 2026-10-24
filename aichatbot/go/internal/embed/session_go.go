//go:build !ORT

package embed

import "example.invalid/knights-analytics/hugot"

// 預設：hugot 的純 Go 後端，不用 cgo、不用裝系統函式庫（M4 一句約 210ms）
const Backend = "go"

func newSession() (*hugot.Session, error) { return hugot.NewGoSession() }
