package portrouter

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
)

// LocalConfigMapLoader 基於本地配置文件方式的端口映射加載器
type LocalConfigMapLoader struct {
	portMap     PortMap
	cfgFileName string
}

// NewLocalConfigMapLoader 創建基於本地配置文件的端口映射加載器
func NewLocalConfigMapLoader(cfgFileName string) *LocalConfigMapLoader {
	pl := new(LocalConfigMapLoader)
	pl.cfgFileName = cfgFileName
	pl.loadConfigData()

	return pl
}

// 加載配置數據
func (mp *LocalConfigMapLoader) loadConfigData() {
	mp.portMap = make(map[int]*router.RouteTarget)

	cfgdata := util.ReadXMLFile(mp.cfgFileName)
	if cfgdata != nil {
		for _, item := range cfgdata {
			//源端口
			sSourcePort := item.Key
			//目的地址
			sTarget := item.Value

			sport, err := strconv.Atoi(sSourcePort)
			if err != nil {
				mlog.Error("解析源端口發生錯誤,%v", err)
				continue
			}

			sArrTarget := strings.Split(sTarget, ":")
			if len(sArrTarget) != 2 {
				mlog.Error("解析目標地址錯誤,必須以:分隔")
				continue
			}

			sTargetIP := sArrTarget[0]
			sTargetPort := sArrTarget[1]
			dport, err := strconv.Atoi(sTargetPort)
			if err != nil {
				mlog.Error("解析目標地址錯誤,端口號不正確.%v", err)
				continue
			}

			rt := new(router.RouteTarget)
			rt.TargetPort = dport
			rt.TargetIP = sTargetIP

			mp.portMap[sport] = rt
		}
	}
}

// LoadPortMap 取得PortMap
func (mp *LocalConfigMapLoader) LoadPortMap() map[int]*router.RouteTarget {
	return mp.portMap
}

// AutoRefreshLocalConfigMapLoader 自動刷新的本地配置映射加載器
type AutoRefreshLocalConfigMapLoader struct {
	locker      *sync.RWMutex
	innerLoader *LocalConfigMapLoader
}

// NewAutoRefreshLocalConfigMapLoader 創建自動刷新的本地映射配置加載器
func NewAutoRefreshLocalConfigMapLoader(cfgFileName string, loadInterval int) *AutoRefreshLocalConfigMapLoader {
	ml := new(AutoRefreshLocalConfigMapLoader)
	ml.locker = &sync.RWMutex{}
	ml.innerLoader = NewLocalConfigMapLoader(cfgFileName)
	ml.innerLoader.loadConfigData()

	go func() {
		time.Sleep(time.Duration(loadInterval) * time.Second)

		ml.Refresh()
	}()

	return ml
}

// Refresh 重新讀取
func (mp *AutoRefreshLocalConfigMapLoader) Refresh() {
	mp.locker.Lock()
	mp.innerLoader.loadConfigData()
	mp.locker.Unlock()
	example.invalid("刷新端口映射配置文件完成...")
}

// LoadPortMap 取得PortMap
func (mp *AutoRefreshLocalConfigMapLoader) LoadPortMap() map[int]*router.RouteTarget {
	mp.locker.RLock()
	defer mp.locker.RUnlock()

	return mp.innerLoader.portMap
}
