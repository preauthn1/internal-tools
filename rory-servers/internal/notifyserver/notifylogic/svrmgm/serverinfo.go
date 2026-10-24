package svrmgm

import (
	"fmt"
	"sync"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// ServerInfo 服務器信息
type ServerInfo struct {
	ConnInfo       *network.ConnInfo // 網絡連接
	ServerType     int               // 服務器類型
	ServerID       int               // 服務器ID
	GameID         int               // 遊戲ID,大廳則為0或者-1
	HostAndBackup  int32
	GameServerAddr string
	FlagID         string
}

func (svr *ServerInfo) String() string {
	return fmt.Sprintf("serverid :%d,servertype:%d,gameid:%d,conn:%v,flagid:%v", svr.ServerID, svr.ServerType, svr.GameID, svr.ConnInfo, svr.FlagID)
}

// NewServerInfo ?
func NewServerInfo(conninfo *network.ConnInfo, serverType int, serverID int, gameid int, hostAndBackup int32, gameServerAddr string, flagID string) *ServerInfo {
	svrinfo := new(ServerInfo)
	svrinfo.ConnInfo = conninfo
	svrinfo.ServerType = serverType
	svrinfo.ServerID = serverID
	svrinfo.GameID = gameid
	svrinfo.HostAndBackup = hostAndBackup
	svrinfo.GameServerAddr = gameServerAddr
	svrinfo.FlagID = flagID

	return svrinfo
}

// ServerList 服務器列表
type ServerList struct {
	servers map[string]*ServerInfo // 服務器信息
	svrLock *sync.RWMutex          // 服務器訪問鎖
}

// NewServerList 創建服務器集合
func NewServerList() *ServerList {
	svrlst := new(ServerList)
	svrlst.servers = make(map[string]*ServerInfo)
	svrlst.svrLock = &sync.RWMutex{}

	return svrlst
}

// RegServer 註冊服務器
func (svrlst *ServerList) RegServer(conninfo *network.ConnInfo, svrType int, svrID int, gameID int, hostAndBackup int32, gameServerAddr string, flagID string) *ServerInfo {
	sign := svrlst.getServerSign(svrType, svrID)
	svrinfo := NewServerInfo(conninfo, svrType, svrID, gameID, hostAndBackup, gameServerAddr, flagID)

	svrlst.svrLock.Lock()
	svrlst.servers[sign] = svrinfo
	svrlst.svrLock.Unlock()

	return svrinfo
}

// 獲取服務器標識
func (svrlst *ServerList) getServerSign(svrtype int, svrid int) string {
	sign := fmt.Sprintf("%d_%d", svrtype, svrid)

	return sign
}

// GetHallServer 獲取所有大廳服務器
func (svrlst *ServerList) GetHallServer() []*ServerInfo {
	return svrlst.GetServerByServerType(conf.ServerTypeHallServer)
}

// GetGameServer 獲取所有遊戲服務器
func (svrlst *ServerList) GetGameServer() []*ServerInfo {
	return svrlst.GetServerByServerType(conf.ServerTypeGameServer)
}

// GetServerByServerType 獲取所有大廳服務器
func (svrlst *ServerList) GetServerByServerType(serverType int) []*ServerInfo {
	svrlst.svrLock.Lock()
	defer svrlst.svrLock.Unlock()

	hsvrs := make([]*ServerInfo, 0)

	for _, svr := range svrlst.servers {
		if svr.ServerType == serverType {
			hsvrs = append(hsvrs, svr)
		}
	}

	return hsvrs
}

// GetAllServer ?
func (svrlst *ServerList) GetAllServer() []*ServerInfo {
	svrlst.svrLock.Lock()
	defer svrlst.svrLock.Unlock()

	svrs := make([]*ServerInfo, 0)

	for _, svr := range svrlst.servers {
		svrs = append(svrs, svr)
	}

	return svrs
}

// GetServerByGameID 根據遊戲ID獲取服務器列表
func (svrlst *ServerList) GetServerByGameID(gameid int) []*ServerInfo {
	svrlst.svrLock.Lock()
	defer svrlst.svrLock.Unlock()

	svrs := make([]*ServerInfo, 0)

	for _, svr := range svrlst.servers {
		if svr.GameID == gameid {
			svrs = append(svrs, svr)
		}
	}

	return svrs
}

// FindServer 查找服務器信息
func (svrlst *ServerList) FindServer(serverType int, serverID int) *ServerInfo {
	svrlst.svrLock.Lock()
	defer svrlst.svrLock.Unlock()

	for _, svr := range svrlst.servers {
		if svr.ServerID == serverID && svr.ServerType == serverType {
			return svr
		}
	}

	return nil
}

// GetServerByConn 根據連接獲取服務器信息
func (svrlst *ServerList) GetServerByConn(conninfo *network.ConnInfo) *ServerInfo {
	svrlst.svrLock.Lock()
	defer svrlst.svrLock.Unlock()

	for _, svr := range svrlst.servers {
		if svr.ConnInfo == conninfo {
			return svr
		}
	}

	return nil
}

// RemoveServerByConn 根據連接移除服務器信息
func (svrlst *ServerList) RemoveServerByConn(conninfo *network.ConnInfo) {
	svr := svrlst.GetServerByConn(conninfo)

	if svr == nil {
		return
	}

	sign := svrlst.getServerSign(svr.ServerType, svr.ServerID)

	svrlst.svrLock.Lock()
	delete(svrlst.servers, sign)
	svrlst.svrLock.Unlock()
}
