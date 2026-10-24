package svrmgm

import (
	"fmt"
	"strconv"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/roomlogic"
)

// SvrMasterSlaveCmd 切換Master與Slave
type SvrMasterSlaveCmd struct {
	msgServer  *msgserver.MsgServer
	serverlist *ServerList
}

// NewSvrMasterSlaveCmd 創建master slave命令對象
func NewSvrMasterSlaveCmd(msgSvr *msgserver.MsgServer, serverlist *ServerList) *SvrMasterSlaveCmd {
	mscmd := new(SvrMasterSlaveCmd)
	mscmd.msgServer = msgSvr
	mscmd.serverlist = serverlist

	return mscmd
}

// Name ?
func (msCmd *SvrMasterSlaveCmd) Name() string {
	return "masterslave"
}

// Help ?
func (msCmd *SvrMasterSlaveCmd) Help() string {
	return "切換AB服. usage: masterslave change [gameid] | 遊戲列表：masterslave gamelist | 切正式：masterslave up [serverid] | 切備援：masterslave down [serverid] | 一次全 up：masterslave allup | 一次全 down：masterslave alldown"
}

// Exec ?
func (msCmd *SvrMasterSlaveCmd) Exec(args []string) string {
	operation := args[0]

	switch operation {
	case "gamelist":
		gsLst := msCmd.serverlist.GetGameServer()
		s := ""
		for _, v := range gsLst {
			if v != nil {
				s += fmt.Sprintf("gamelist GameId[%v] ServerId[%v] FlagId [%v] HostAndBackup[%v] Addr[%v] \n", v.GameID, v.ServerID, v.FlagID, v.HostAndBackup, v.GameServerAddr)
			}
		}
		return s
	case "up":
		iServerID, err := strconv.Atoi(args[1])
		if err != nil {
			return "error, " + err.Error()
		}

		return msCmd.masterslaveUpDown(iServerID, true)
	case "down":
		iServerID, err := strconv.Atoi(args[1])
		if err != nil {
			return "error, " + err.Error()
		}

		return msCmd.masterslaveUpDown(iServerID, false)
	case "change":
		masterGameID, err := strconv.Atoi(args[1])

		if err != nil {
			return "error, master gameid error"
		}
		mlog.Debug("[masterslave] masterslave change GameID[%d]", masterGameID)

		gsLst := msCmd.serverlist.GetGameServer()

		masterServers := make([]*ServerInfo, 0)
		slaveServers := make([]*ServerInfo, 0)
		for _, v := range gsLst {
			if v.GameID == masterGameID {
				if v.HostAndBackup == roomlogic.HostTypeMaster {
					masterServers = append(masterServers, v)
				} else {
					slaveServers = append(slaveServers, v)
				}
			}
		}
		// 我們允許只要同場次有服務器存在就好, 不需要存同場次的服務器數量一定要一樣
		// 從 master 來找 flagid 在 slave 是否存在, 不存在則不允許切換
		validateMS := true
		for _, ms := range masterServers {
			hasSlave := false
			for _, ss := range slaveServers {
				if ms.FlagID == ss.FlagID {
					hasSlave = true
				}
			}
			if hasSlave == false {
				validateMS = false
			}
		}
		if validateMS == false {
			return "error, serverid not match 失敗, 正式與備援沒有對應相同數量"
		}

		//切換
		for _, masterServer := range masterServers {
			var slaveServer *ServerInfo
			for _, ss := range slaveServers { //找出對應場次的slave server做切換
				if masterServer.FlagID == ss.FlagID {
					slaveServer = ss
					svrHostAndBackUpMsg := new(msg.Message)
					svrHostAndBackUpMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
					svrHostAndBackUpMsg.SClassID = int32(netproto.NotifyServerClassID_ServerHostAndBackUp)
					svrHostData := new(netproto.HostAndBackUp)
					svrHostData.Master = new(netproto.HostInfo)
					svrHostData.Slave = new(netproto.HostInfo)
					svrHostData.Master.GameID = proto.Int32(int32(masterGameID))
					svrHostData.Master.ServerID = proto.Int32(int32(slaveServer.ServerID)) //切換
					svrHostData.Slave.GameID = proto.Int32(int32(masterGameID))
					svrHostData.Slave.ServerID = proto.Int32(int32(masterServer.ServerID))      //切換
					svrHostData.Slave.GameServerAddr = proto.String(slaveServer.GameServerAddr) //告訴正式服，備援服的地址
					mlog.System("master slave change, master is %d, slave is %d, client go to gameserveraddr %s",
						slaveServer.ServerID, masterServer.ServerID, slaveServer.GameServerAddr)
					svrHostAndBackUpMsg.MsgData = svrHostData
					msCmd.msgServer.SendMessage(masterServer.ConnInfo, svrHostAndBackUpMsg)
					time.Sleep(200 * time.Millisecond)
					msCmd.msgServer.SendMessage(slaveServer.ConnInfo, svrHostAndBackUpMsg)
					msCmd.stopGameServerRoom(svrHostData.Slave)
					masterServer.HostAndBackup = roomlogic.HostTypeSlave
				}
			}
		}
		for _, ss := range slaveServers {
			ss.HostAndBackup = roomlogic.HostTypeMaster
		}

		return "master slave change done."
	case "allup": // gameserver 一次全 up
		return msCmd.masterslaveAllUpDown(true)
	case "alldown": // gameserver 一次全 down
		return msCmd.masterslaveAllUpDown(false)
	default:
		return "masterslave 參數不正確"
	}

	return "masterslave 參數不正確"
}

//關閉遊戲房間
func (msCmd *SvrMasterSlaveCmd) stopGameServerRoom(slaveInfo *netproto.HostInfo) {
	svrstopMsg := new(msg.Message)
	svrstopMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	svrstopMsg.SClassID = int32(netproto.NotifyServerClassID_ServerStopID)

	svrstopData := new(netproto.ServerStop)
	svrstopData.Type = proto.Int32(0)
	svrstopData.Status = proto.Int32(conf.ServerStatusStoping)

	svrstopData.GameID = slaveInfo.GameID
	svrstopData.ServerID = slaveInfo.ServerID
	svrstopData.Message = proto.String("-")
	svrstopMsg.MsgData = svrstopData
	mlog.Debug("[masterslave] stopGameServerRoom, slaveInfo[%d, %d], svrstopData[%d, %d]",
		*slaveInfo.GameID, *slaveInfo.ServerID, *svrstopData.GameID, *svrstopData.ServerID)

	svrSlave := msCmd.serverlist.FindServer(conf.ServerTypeGameServer, int(slaveInfo.GetServerID()))
	if svrSlave != nil {
		mlog.Debug("[masterslave] svrSlave[%d, %d]", svrSlave.GameID, svrSlave.ServerID)
		msCmd.msgServer.SendMessage(svrSlave.ConnInfo, svrstopMsg)
	}
}

// 開啟遊戲房間
func (msCmd *SvrMasterSlaveCmd) startGameServerRoom(masterInfo *netproto.HostInfo) {
	svrstopMsg := new(msg.Message)
	svrstopMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	svrstopMsg.SClassID = int32(netproto.NotifyServerClassID_ServerStopID)

	svrstopData := new(netproto.ServerStop)
	svrstopData.Type = proto.Int32(0)
	svrstopData.Status = proto.Int32(conf.ServerStatusRuning)

	svrstopData.GameID = masterInfo.GameID
	svrstopData.ServerID = masterInfo.ServerID
	svrstopData.Message = proto.String("-")
	svrstopMsg.MsgData = svrstopData
	mlog.Debug("[masterslave] startGameServerRoom, masterInfo[%d, %d], svrstopData[%d, %d]",
		*masterInfo.GameID, *masterInfo.ServerID, *svrstopData.GameID, *svrstopData.ServerID)

	svrMaster := msCmd.serverlist.FindServer(conf.ServerTypeGameServer, int(masterInfo.GetServerID()))
	if svrMaster != nil {
		mlog.Debug("[masterslave] svrMaster[%d, %d]", svrMaster.GameID, svrMaster.ServerID)
		msCmd.msgServer.SendMessage(svrMaster.ConnInfo, svrstopMsg)
	}
}

// masterslave up or down
func (msCmd *SvrMasterSlaveCmd) masterslaveUpDown(iServerID int, bMaster bool) string {
	var strRetMsg string = "masterslave 參數不正確"

	svr := msCmd.serverlist.FindServer(conf.ServerTypeGameServer, iServerID)
	if svr == nil {
		return fmt.Sprintf("error, server %d not exists.", iServerID)
	}

	if bMaster {
		svr.HostAndBackup = roomlogic.HostTypeMaster
	} else {
		svr.HostAndBackup = roomlogic.HostTypeSlave
	}
	svrHostAndBackUpMsg := new(msg.Message)
	svrHostAndBackUpMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	svrHostAndBackUpMsg.SClassID = int32(netproto.NotifyServerClassID_ServerHostAndBackUp)
	svrHostData := new(netproto.HostAndBackUp)
	svrHostData.Master = new(netproto.HostInfo)
	svrHostData.Slave = new(netproto.HostInfo)
	switch svr.HostAndBackup {
	case roomlogic.HostTypeMaster:
		svrHostData.Master.GameID = proto.Int32(int32(svr.GameID))
		svrHostData.Master.ServerID = proto.Int32(int32(svr.ServerID))
		svrHostData.Slave.GameID = proto.Int32(int32(0))
		svrHostData.Slave.ServerID = proto.Int32(int32(0))
	case roomlogic.HostTypeSlave:
		svrHostData.Master.GameID = proto.Int32(int32(0))
		svrHostData.Master.ServerID = proto.Int32(int32(0))
		svrHostData.Slave.GameID = proto.Int32(int32(svr.GameID))
		svrHostData.Slave.ServerID = proto.Int32(int32(svr.ServerID))
	}
	svrHostAndBackUpMsg.MsgData = svrHostData
	msCmd.msgServer.SendMessage(svr.ConnInfo, svrHostAndBackUpMsg)
	time.Sleep(200 * time.Millisecond)

	switch svr.HostAndBackup {
	case roomlogic.HostTypeMaster:
		msCmd.startGameServerRoom(svrHostData.Master)
		mlog.Debug("[masterslave] masterslave up GameID[%d], ServerID[%d]", *svrHostData.Master.GameID, *svrHostData.Master.ServerID)
		strRetMsg = fmt.Sprintf("server %d up done.", iServerID)
	case roomlogic.HostTypeSlave:
		msCmd.stopGameServerRoom(svrHostData.Slave)
		mlog.Debug("[masterslave] masterslave down GameID[%d], ServerID[%d]", *svrHostData.Slave.GameID, *svrHostData.Slave.ServerID)
		strRetMsg = fmt.Sprintf("server %d down done.", iServerID)
	}

	return strRetMsg
}

// masterslave all up or down
func (msCmd *SvrMasterSlaveCmd) masterslaveAllUpDown(bMaster bool) string {
	var strRetMsg string

	if bMaster {
		strRetMsg = "masterslave allup result:\n"
	} else {
		strRetMsg = "masterslave alldown result:\n"
	}

	aryPGameList := msCmd.serverlist.GetGameServer()
	for _, pServer := range aryPGameList {
		if pServer == nil {
			continue
		}

		msCmd.masterslaveUpDown(pServer.ServerID, bMaster)
	}

	aryPResult := msCmd.serverlist.GetGameServer()
	for _, pServer := range aryPResult {
		if pServer == nil {
			continue
		}

		strRetMsg += fmt.Sprintf("gamelist GameId[%v] ServerId[%v] FlagId [%v] HostAndBackup[%v] Addr[%v] \n",
			pServer.GameID, pServer.ServerID, pServer.FlagID, pServer.HostAndBackup, pServer.GameServerAddr)
	}

	return strRetMsg
}
