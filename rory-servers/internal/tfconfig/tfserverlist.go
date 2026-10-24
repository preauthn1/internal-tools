package tfconfig

import (
	"sync"
	"time"
)

// TFServerList 轉發服務器列表
type TFServerList struct {
	innerServerLst map[string]*TFServerInfo // 服務器列表
	locker         *sync.RWMutex            // 鎖
}

// NewTFserverList create TFserverList
func NewTFserverList() *TFServerList {
	tfs := &TFServerList{}
	tfs.innerServerLst = make(map[string]*TFServerInfo)
	tfs.locker = &sync.RWMutex{}

	return tfs
}

// GetServerList 取得服務器列表
func (tfs *TFServerList) GetServerList() map[string]TFServerInfo {
	tfs.locker.RLock()
	defer tfs.locker.RUnlock()

	data := make(map[string]TFServerInfo)

	for k, v := range tfs.innerServerLst {
		data[k] = *v
	}

	return data
}

// UpdateTFServer 更新服務器訊息
func (tfs *TFServerList) UpdateTFServer(ip string, cn int) {
	tfs.locker.Lock()
	defer tfs.locker.Unlock()

	etf, ok := tfs.innerServerLst[ip]
	if ok {
		etf.CurrentConn = cn
		etf.IPAddress = ip
		etf.LastRegTime = time.Now()
	} else {
		tfs.innerServerLst[ip] = &TFServerInfo{IPAddress: ip, CurrentConn: cn, LastRegTime: time.Now()}
	}
}

// TFServerManager 管理TFServer
var TFServerManager *TFServerList

func init() {
	TFServerManager = NewTFserverList()
}
