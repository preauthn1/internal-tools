package blackip

import (
	"time"
)

// Info 黑名單信息
type Info struct {
	IPAddr       string    // ip地址
	BlackEndTime time.Time // 黑名單截止時間
	BlackMemo    string    // 加入黑名單原因
}

// IsExpire 是否過期
func (bip *Info) IsExpire() bool {
	return time.Now().After(bip.BlackEndTime)
}
