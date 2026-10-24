package msghandler

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/pkg/command"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"

	"fmt"
	"sort"
	"sync"
	"time"
)

// GameServerInfo 遊戲服務器信息
type GameServerInfo struct {
	ServerRegInfo *netproto.GameServerRegInfo // gameserver註冊信息
	ConnInfo      *network.ConnInfo           // 連接信息
	LastRegTime   time.Time                   // 最後一次註冊時間
}

// GameServerList 遊戲服務器列表
type GameServerList struct {
	innerGameServerList map[int32]*GameServerInfo
	svrInfoLock         *sync.RWMutex
	hallVerProvider     hallutil.HallVerProvider
}

// NewGameServerList 創建遊戲服務器列表對象
func NewGameServerList() *GameServerList {
	glist := new(GameServerList)
	glist.innerGameServerList = make(map[int32]*GameServerInfo)
	glist.svrInfoLock = &sync.RWMutex{}

	// 添加命令
	command.AddCommand(NewGsMgmCmd(glist))

	return glist
}

// RangeGameServer 遍歷所有遊戲服務器信息
func (glist *GameServerList) RangeGameServer(f func(svrinfo *GameServerInfo)) {
	glist.svrInfoLock.RLock()
	defer glist.svrInfoLock.RUnlock()

	var keys []int
	for k := range glist.innerGameServerList {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, s := range keys {
		f(glist.innerGameServerList[int32(s)])
	}
}

// RangeGameServerOrderbyUserCount 遍歷所有遊戲服務器信息, 使用者少至多排序
func (glist *GameServerList) RangeGameServerOrderbyUserCount(f func(svrinfo *GameServerInfo)) {
	glist.svrInfoLock.RLock()
	defer glist.svrInfoLock.RUnlock()

	var keys []int
	for k := range glist.innerGameServerList {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for skey, s := range keys {
		list1 := glist.innerGameServerList[int32(s)]
		list1AllUserCount := list1.ServerRegInfo.GetUserCount() + list1.ServerRegInfo.GetRobotCount()
		for ckey, c := range keys {
			list2 := glist.innerGameServerList[int32(c)]
			list2AllUserCount := list2.ServerRegInfo.GetUserCount() + list2.ServerRegInfo.GetRobotCount()

			if list2AllUserCount <= list1AllUserCount {
				keys[skey], keys[ckey] = keys[ckey], keys[skey]
			}
		}
	}

	for _, s := range keys {
		if glist.innerGameServerList[int32(s)].ServerRegInfo.GetHostAndBackup() == roomlogic.HostTypeMaster {
			f(glist.innerGameServerList[int32(s)])
		}
	}

}

// GameServerIsRunChangeToSlave 檢查GameServer的Master是否在執行中, 若有則要將目前要註冊的GameServer改為Slave
func (glist *GameServerList) GameServerIsRunChangeToSlave(gsreginfo *netproto.GameServerRegInfo) []*msg.Message {
	msgarr := make([]*msg.Message, 0)
	sameServerID := false
	glist.RangeGameServer(func(svrinfo *GameServerInfo) {
		if svrinfo.ServerRegInfo.GetGameID() == gsreginfo.GetGameID() &&
			svrinfo.ServerRegInfo.GetServerID() == gsreginfo.GetServerID() &&
			svrinfo.ServerRegInfo.GetHostAndBackup() == gsreginfo.GetHostAndBackup() {
			sameServerID = true
		}
	})

	if sameServerID {
		return nil
	}

	glist.RangeGameServer(func(svrinfo *GameServerInfo) {
		if svrinfo.ServerRegInfo.GetGameID() == gsreginfo.GetGameID() &&
			svrinfo.ServerRegInfo.GetFlag() == gsreginfo.GetFlag() &&
			gsreginfo.GetHostAndBackup() == roomlogic.HostTypeMaster {

			svrHostAndBackUpMsg := new(msg.Message)
			svrHostAndBackUpMsg.BClassID = int32(netproto.MessageBClassID_ServerMgm)
			svrHostAndBackUpMsg.SClassID = int32(netproto.NotifyServerClassID_ServerHostAndBackUp)

			svrHostData := new(netproto.HostAndBackUp)
			svrHostData.Master = new(netproto.HostInfo)
			svrHostData.Slave = new(netproto.HostInfo)

			svrHostData.Master.GameID = proto.Int32(int32(svrinfo.ServerRegInfo.GetGameID()))
			svrHostData.Slave.GameID = proto.Int32(int32(svrinfo.ServerRegInfo.GetGameID()))

			if svrinfo.ServerRegInfo.GetHostAndBackup() == int32(0) {
				svrHostData.Master.ServerID = proto.Int32(int32(svrinfo.ServerRegInfo.GetServerID()))
				svrHostData.Slave.ServerID = proto.Int32(int32(gsreginfo.GetServerID()))
			} else if gsreginfo.GetHostAndBackup() == int32(0) {
				svrHostData.Master.ServerID = proto.Int32(int32(gsreginfo.GetServerID()))
				svrHostData.Slave.ServerID = proto.Int32(int32(svrinfo.ServerRegInfo.GetServerID()))
			}

			mlog.Debug("修改slave gsHost[%v] master[%v] slave[%v]", gsreginfo.GetHostAndBackup(),
				svrHostData.Master.GetServerID(), svrHostData.Slave.GetServerID())

			svrHostAndBackUpMsg.MsgData = svrHostData

			msgarr = append(msgarr, svrHostAndBackUpMsg)
		}
	})

	return msgarr
}

// SetGameServer 設置遊戲服務器信息
func (glist *GameServerList) SetGameServer(gsreginfo *netproto.GameServerRegInfo, conninfo *network.ConnInfo) {
	glist.svrInfoLock.Lock()
	defer glist.svrInfoLock.Unlock()

	gsinfo := new(GameServerInfo)
	gsinfo.ServerRegInfo = gsreginfo
	gsinfo.ConnInfo = conninfo
	gsinfo.LastRegTime = time.Now()

	glist.innerGameServerList[gsreginfo.GetServerID()] = gsinfo
}

// RemoveGameServer 移除遊戲服務器信息
func (glist *GameServerList) RemoveGameServer(serverid int32) {
	if _, ok := glist.innerGameServerList[serverid]; ok {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[GameServerList]移除遊戲服務器信息 serverid[%v]", serverid)).System("HallServer壓測時")
		delete(glist.innerGameServerList, serverid)
	}
}

// RemoveGameServerByConn 根據連接信息移除服務器
func (glist *GameServerList) RemoveGameServerByConn(conninfo *network.ConnInfo) {
	glist.svrInfoLock.Lock()
	defer glist.svrInfoLock.Unlock()

	svrid := make([]int32, 0)

	for id, s := range glist.innerGameServerList {
		if s.ConnInfo == conninfo {
			svrid = append(svrid, int32(id))
		}
	}

	for _, id := range svrid {
		mlog.NewDefaultLogger().SetData(fmt.Sprintf("[GameServerList]移除遊戲服務器信息 serverid[%v]", int32(id))).Error("HallServer壓測時")
		delete(glist.innerGameServerList, int32(id))
	}
}

// GetGameServerByID 根據服務器ID獲取服務器信息
func (glist *GameServerList) GetGameServerByID(serverid int32) *netproto.GameServerRegInfo {
	glist.svrInfoLock.RLock()
	defer glist.svrInfoLock.RUnlock()

	gsinfo := glist.innerGameServerList[serverid]

	if gsinfo == nil {
		return nil
	}

	return gsinfo.ServerRegInfo
}

// SetHallVerProvider ?
func (glist *GameServerList) SetHallVerProvider(hallVerProvider hallutil.HallVerProvider) {
	glist.hallVerProvider = hallVerProvider
}

func (glist *GameServerList) getGameListByAreaID(areaID int64) []*GameServerInfo {
	platformGameServerList := glist.hallVerProvider.GetGameServerListByArea(areaID) // 取得平台可用的服務器列表
	gsList := make([]*GameServerInfo, 0)

	for _, gs := range glist.innerGameServerList {
		for _, v := range platformGameServerList {
			if v.GetStatus() && gs.ServerRegInfo.GetGameID() == int32(v.GetGameID()) && gs.ServerRegInfo.GetServerID() == int32(v.GetServerID()) {
				//存在平台表裡且狀態為true的
				gsList = append(gsList, gs)
				break
			}
		}
	}
	return gsList
}

// RangeGameServerByAreaID 遍歷所有遊戲服務器信息
func (glist *GameServerList) RangeGameServerByAreaID(areaID int64, f func(svrinfo *GameServerInfo)) {
	glist.svrInfoLock.RLock()
	defer glist.svrInfoLock.RUnlock()
	gsList := glist.getGameListByAreaID(areaID)

	for _, gs := range gsList {
		f(gs)
	}
}

// RangeGameServerOrderbyUserCountByAreaID 遍歷所有遊戲服務器信息, 使用者少至多排序
func (glist *GameServerList) RangeGameServerOrderbyUserCountByAreaID(areaID int64, f func(svrinfo *GameServerInfo)) {
	glist.svrInfoLock.RLock()
	defer glist.svrInfoLock.RUnlock()

	gsList := glist.getGameListByAreaID(areaID)

	sort.Slice(gsList, func(i, j int) bool {
		iCount := gsList[i].ServerRegInfo.GetUserCount() + gsList[i].ServerRegInfo.GetRobotCount()
		jCount := gsList[j].ServerRegInfo.GetUserCount() + gsList[j].ServerRegInfo.GetRobotCount()
		return iCount < jCount
	})

	for _, gs := range gsList {
		if gs.ServerRegInfo.GetHostAndBackup() == roomlogic.HostTypeMaster {
			f(gs)
		}
	}
}
