package tfconfig

import (
	"fmt"
	"time"
)

// TFServerInfo TFServer 資訊
type TFServerInfo struct {
	IPAddress   string    // 轉發服務器IP位置
	CurrentConn int       // 當前連接數
	LastRegTime time.Time // 最後一次註冊時間
}

func (tfs *TFServerInfo) String() string {
	return fmt.Sprintf("%s 連接數:%d %v", tfs.IPAddress, tfs.CurrentConn, tfs.LastRegTime)
}
