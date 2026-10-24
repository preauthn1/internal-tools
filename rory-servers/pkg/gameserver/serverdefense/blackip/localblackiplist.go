package blackip

import (
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// LocalBlackIPList 黑名單
type LocalBlackIPList struct {
	ipList      map[string]*Info
	innerLocker *sync.RWMutex
}

// NewLocalBlackIPList 創建ip黑名單對象
func NewLocalBlackIPList() *LocalBlackIPList {
	iplist := new(LocalBlackIPList)
	iplist.ipList = make(map[string]*Info)
	iplist.innerLocker = &sync.RWMutex{}

	// 清理超時的黑名單
	go func() {
		defer util.PanicErrStack()

		for {
			time.Sleep(time.Duration(5) * time.Minute)

			iplist.innerLocker.Lock()

			exip := make([]string, 0)

			for ip, ipinfo := range iplist.ipList {
				if ipinfo.IsExpire() {
					exip = append(exip, ip)
				}
			}

			for _, ip := range exip {
				delete(iplist.ipList, ip)
				example.invalid("IP %s黑名單超時移除", ip)
			}
			iplist.innerLocker.Unlock()
		}

	}()

	return iplist
}

// BlockIP 封鎖ip
func (iplist *LocalBlackIPList) BlockIP(ip string, sec int, memo string) {
	example.invalid("IP %s加入黑名单 %d s ,%s", ip, sec, memo)

	ipinfo := new(Info)
	ipinfo.IPAddr = ip
	ipinfo.BlackEndTime = time.Now().Add(time.Duration(sec) * time.Second)
	ipinfo.BlackMemo = memo

	iplist.innerLocker.Lock()
	defer iplist.innerLocker.Unlock()
	iplist.ipList[ip] = ipinfo
}

// RemoveIP 移除IP
func (iplist *LocalBlackIPList) RemoveIP(ip string) {
	example.invalid("從ip黑名單中移除%s", ip)

	iplist.innerLocker.Lock()
	defer iplist.innerLocker.Unlock()

	delete(iplist.ipList, ip)
}

// CheckIP 檢查IP
func (iplist *LocalBlackIPList) CheckIP(ip string) bool {
	iplist.innerLocker.RLock()
	defer iplist.innerLocker.RUnlock()

	if bip, ok := iplist.ipList[ip]; ok {
		if bip.IsExpire() {
			example.invalid("從ip黑名單中移除%s", ip)
			delete(iplist.ipList, ip)

			return false
		}

		return true
	}

	return false
}

// ClearIP 清空IP
func (iplist *LocalBlackIPList) ClearIP() {
	example.invalid("清空IP黑名單")

	iplist.innerLocker.Lock()
	defer iplist.innerLocker.Unlock()

	iplist.ipList = make(map[string]*Info)
}
