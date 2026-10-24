package blackip

import (
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// LoadBlackIPFunc 加載黑名單函數
type LoadBlackIPFunc func(func(map[string]string))

// RPCBlackIPList rpc IP黑名單
type RPCBlackIPList struct {
	blackIPList     map[string]string
	loadBlackIPFunc LoadBlackIPFunc
	blockIPFunc     func(ip string, sec int, memo string)
	loadInterval    int

	innerLocker *sync.RWMutex
}

// NewRPCBlackIPList 創建RPC IP黑名單
func NewRPCBlackIPList(loadInterval int, loadIPFunc LoadBlackIPFunc, blockIPFunc func(ip string, sec int, memo string)) *RPCBlackIPList {
	rpciplist := new(RPCBlackIPList)
	rpciplist.loadBlackIPFunc = loadIPFunc
	rpciplist.blockIPFunc = blockIPFunc
	rpciplist.innerLocker = &sync.RWMutex{}
	rpciplist.loadInterval = loadInterval
	rpciplist.blackIPList = make(map[string]string)

	return rpciplist
}

// StartLoadIPList 開始
func (rpciplst *RPCBlackIPList) StartLoadIPList() {
	go func() {
		defer util.PanicErrStack()

		for {
			rpciplst.loadBlackIPFunc(func(ips map[string]string) {
				if ips == nil {
					ips = make(map[string]string)
				}

				rpciplst.innerLocker.Lock()
				rpciplst.blackIPList = ips
				rpciplst.innerLocker.Unlock()
			})

			time.Sleep(time.Duration(rpciplst.loadInterval) * time.Second)
		}
	}()
}

// BlockIP 封鎖ip
func (rpciplst *RPCBlackIPList) BlockIP(ip string, sec int, memo string) {
	rpciplst.blockIPFunc(ip, sec, memo)

	rpciplst.innerLocker.Lock()
	rpciplst.blackIPList[ip] = ip
	rpciplst.innerLocker.Unlock()

}

// CheckIP 檢查IP
func (rpciplst *RPCBlackIPList) CheckIP(ip string) bool {
	rpciplst.innerLocker.Lock()
	defer rpciplst.innerLocker.Unlock()

	if _, ok := rpciplst.blackIPList[ip]; ok {
		return true
	}

	return false
}

// Clear 清空IP
func (rpciplst *RPCBlackIPList) Clear() {
	rpciplst.innerLocker.Lock()
	defer rpciplst.innerLocker.Unlock()

	rpciplst.blackIPList = make(map[string]string)
}
