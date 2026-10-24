package msghandler

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
)

// ServerAllocStrategy 服務器分配策略
type ServerAllocStrategy interface {
	GetServerStatus(usercount int32) string                                                                               // 獲取服務器狀態描述
	GetServerAddr(serverlist *GameServerList, gameid int32, flag string) []*netproto.GameServerAddr                       // 分配服務器
	GetServerAddrByAreaID(serverlist *GameServerList, areaid int64, gameid int32, flag string) []*netproto.GameServerAddr // 分配服務器
}

// DefaultServerAllocStrategy 默認的服務器分配策略
type DefaultServerAllocStrategy struct {
	Config   *hallutil.HallServerConfig
	Strategy []*hallutil.ServerAllocStrategyList
}

// NewDefaultServerAllocStrategy 創建默認的服務器分配策略
func NewDefaultServerAllocStrategy(config *hallutil.HallServerConfig) *DefaultServerAllocStrategy {
	dftallocStrategy := new(DefaultServerAllocStrategy)

	dftallocStrategy.Config = config
	dftallocStrategy.Strategy = config.StrategyList

	return dftallocStrategy
}

// GetServerStatus 未知
func (dftstrategy *DefaultServerAllocStrategy) GetServerStatus(usercount int32) string {
	return "未知"
}

// GetServerAddr 未知
func (dftstrategy *DefaultServerAllocStrategy) GetServerAddr(serverlist *GameServerList, gameid int32, flag string) []*netproto.GameServerAddr {
	addrArray := make([]*netproto.GameServerAddr, 0)
	serverlist.RangeGameServerOrderbyUserCount(func(svr *GameServerInfo) {
		if svr.ServerRegInfo.GetFlag() == flag && svr.ServerRegInfo.GetGameID() == gameid {
			var addr *netproto.GameServerAddr
			addr = StrategyLiangHao(svr, dftstrategy, gameid, flag, addr)
			addrArray = append(addrArray, addr)
		}
	})

	if addrArray != nil && len(addrArray) != 0 {
		return addrArray
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[GsMgmHandler] gameid[%v] flag[%v] 服务器分配策略找不到適合空間，請儘速洽資訊人員.", gameid, flag)).Error("HallServer")
	return nil
}

// GetServerAddrByAreaID 未知
func (dftstrategy *DefaultServerAllocStrategy) GetServerAddrByAreaID(serverlist *GameServerList, areaID int64, gameID int32, flag string) []*netproto.GameServerAddr {
	addrArray := make([]*netproto.GameServerAddr, 0)
	serverlist.RangeGameServerOrderbyUserCountByAreaID(areaID, func(svr *GameServerInfo) {
		if svr.ServerRegInfo.GetFlag() == flag && svr.ServerRegInfo.GetGameID() == gameID {
			var addr *netproto.GameServerAddr
			addr = StrategyLiangHao(svr, dftstrategy, gameID, flag, addr)
			if addr == nil {
				return
			}
			addrArray = append(addrArray, addr)
		}
	})

	if addrArray != nil && len(addrArray) != 0 {
		return addrArray
	}

	mlog.NewDefaultLogger().SetData(fmt.Sprintf("[GsMgmHandler] gameid[%v] flag[%v] 服务器分配策略找不到適合空間，請儘速洽資訊人員.", gameID, flag)).Error("HallServer")
	return nil
}

// StrategyLiangHao 分配可用空間
func StrategyLiangHao(svr *GameServerInfo, dftstrategy *DefaultServerAllocStrategy, gameid int32, flag string, addr *netproto.GameServerAddr) *netproto.GameServerAddr {
	if svr.ServerRegInfo.GetHostAndBackup() != roomlogic.HostTypeMaster {
		return nil
	}

	urcount := svr.ServerRegInfo.GetUserCount() + svr.ServerRegInfo.GetRobotCount()
	if urcount <= dftstrategy.Config.MaxUser {
		addr = new(netproto.GameServerAddr)
		addr.ServerID = svr.ServerRegInfo.ServerID
		addr.IPAddress = svr.ServerRegInfo.IPAddress
		addr.Port = svr.ServerRegInfo.Port
	}

	// for _, v := range dftstrategy.Strategy {
	// 	if _, ok := v.GameIdList[svr.ServerRegInfo.GetGameID()]; ok {
	// 		if urcount <= *v.MaxUser && svr.ServerRegInfo.GetHostAndBackup() == roomlogic.HostType_Master {
	// 			addr = new(netproto.GameServerAddr)
	// 			addr.ServerID = svr.ServerRegInfo.ServerID
	// 			addr.IPAddress = svr.ServerRegInfo.IPAddress
	// 			addr.Port = svr.ServerRegInfo.Port
	// 		}
	// 	}
	// }

	return addr
}
