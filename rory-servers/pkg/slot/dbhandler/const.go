package dbhandler

// Mode Service模式
type Mode uint32

// Modes
const (
	// DB 讀取 SP
	DB Mode = iota

	// 讀取 API
	API
)
