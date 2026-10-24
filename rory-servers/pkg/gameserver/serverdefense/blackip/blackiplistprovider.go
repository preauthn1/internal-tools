package blackip

// ListProvider ip黑名單
type ListProvider interface {
	// 加入黑名單
	BlockIP(ip string, sec int, memo string)

	// 檢查IP是否被攔截
	CheckIP(ip string) bool
}
