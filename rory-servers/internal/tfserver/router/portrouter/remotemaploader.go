package portrouter

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/util"

	"example.invalid/ServerTeam/Servers/internal/tfserver/router"
)

// RemoteMapLoader 遠程端口映射加載器
type RemoteMapLoader struct {
	remoteAddr string
	pmData     PortMap
}

// NewRemoteMapLoader 創建遠程端口映射加載器對象
func NewRemoteMapLoader(remoteAddr string) *RemoteMapLoader {
	rml := new(RemoteMapLoader)
	rml.remoteAddr = remoteAddr
	rml.pmData = rml.loadMapData()

	return rml
}

func (rml *RemoteMapLoader) loadMapData() PortMap {

	resp, err := http.Get(rml.remoteAddr)
	defer func() {
		if resp != nil && resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		mlog.Error("加載遠程端口映射配置文件錯誤%v", err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		mlog.Error("讀取遠程端口映射配置文件錯誤.%v", err)
		return nil
	}

	cfgdata := util.ReadXMLMap(body)
	pm := make(map[int]*router.RouteTarget)

	for skey, sval := range cfgdata {
		//源端口
		sSourcePort := skey
		//目的地址
		sTarget := sval

		sport, err := strconv.Atoi(sSourcePort)
		if err != nil {
			mlog.Error("解析源端口發生錯誤,%s,%v", sSourcePort, err)
			continue
		}

		sArrTarget := strings.Split(sTarget, ":")
		if len(sArrTarget) != 2 {
			mlog.Error("解析目標地址錯誤,必須以:分隔.%s", sTarget)
			continue
		}

		sTargetIP := sArrTarget[0]
		sTargetPort := sArrTarget[1]
		dport, err := strconv.Atoi(sTargetPort)
		if err != nil {
			mlog.Error("解析目标地址错误,端口号不正确.%s,%v", sTarget, err)
			continue
		}

		rt := new(router.RouteTarget)
		rt.TargetPort = dport
		rt.TargetIP = sTargetIP

		pm[sport] = rt
	}

	return pm
}

// LoadPortMap 取得PortMap
func (rml *RemoteMapLoader) LoadPortMap() PortMap {
	return rml.pmData
}

// AutoRefreshRemoteMapLoader 自動刷新的遠程端口映射加載器
type AutoRefreshRemoteMapLoader struct {
	innerRemoteMapLoader *RemoteMapLoader
	locker               *sync.RWMutex
	portMap              PortMap
}

// NewAutoRefreshRemoteConfigMapLoader 創建自動刷新的遠程端口映射加載器對象
func NewAutoRefreshRemoteConfigMapLoader(addr string, refreshInterval int) *AutoRefreshRemoteMapLoader {
	arml := new(AutoRefreshRemoteMapLoader)
	arml.innerRemoteMapLoader = NewRemoteMapLoader(addr)
	arml.locker = &sync.RWMutex{}

	arml.portMap = arml.innerRemoteMapLoader.LoadPortMap()
	example.invalid("刷新遠程端口映射完成...%v", arml.portMap)
	// 定時重新加載
	go func() {
		for {
			time.Sleep(time.Duration(refreshInterval) * time.Second)

			arml.Refresh()
		}
	}()

	return arml
}

// Refresh 重新讀取
func (arml *AutoRefreshRemoteMapLoader) Refresh() {
	arml.locker.Lock()
	mp := arml.innerRemoteMapLoader.loadMapData()
	arml.locker.Unlock()

	if mp != nil && len(mp) > 0 {
		arml.portMap = mp
		example.invalid("刷新遠程端口映射完成...%v", mp)
	}
}

// LoadPortMap 取得PortMap
func (arml *AutoRefreshRemoteMapLoader) LoadPortMap() PortMap {
	arml.locker.Lock()
	defer arml.locker.Unlock()

	return arml.portMap
}
