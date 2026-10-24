package msghandler

import (
	"bytes"
	"fmt"
	"strconv"

	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GsMgmCmd 遊戲服務器管理命令
type GsMgmCmd struct {
	gameserverlist *GameServerList
}

// NewGsMgmCmd 創建遊戲服務器管理器命令
func NewGsMgmCmd(gslist *GameServerList) *GsMgmCmd {
	gsmcmd := new(GsMgmCmd)
	gsmcmd.gameserverlist = gslist

	return gsmcmd
}

// Name ?
func (cmd *GsMgmCmd) Name() string {
	return "gameserverlist"
}

// Help ?
func (cmd *GsMgmCmd) Help() string {
	return "遊戲服務器管理.usage gameserverlist/gameserverlist area areaid/gameserverlist del serverid/gameserverlist find gameid (flag)"
}

// Exec ?
func (cmd *GsMgmCmd) Exec(args []string) string {
	arglen := len(args)

	//列表gameserver
	if arglen == 0 {
		var totaluserCount int32 = 0
		var totalRobotCount int32 = 0

		var buf bytes.Buffer
		cmd.gameserverlist.RangeGameServer(func(svrinfo *GameServerInfo) {
			s := svrinfo.ServerRegInfo
			buf.WriteString(fmt.Sprintf("gameid:%d,gamename:%s,serverid:%d,servername:%s,hostandbackup:%v,addr:%s:%d,flag:%s,maxodd:%d,loginmoney:%d,usercount:%d,robotcount:%d,androidver:%s,iosver:%s\r\n", s.GetGameID(), s.GetGameName(), s.GetServerID(), s.GetServerName(), s.GetHostAndBackup(), s.GetIPAddress(), s.GetPort(), s.GetFlag(), s.GetMaxOdd(), s.GetLoginMoney(), s.GetUserCount(), s.GetRobotCount(), s.GetAndroidVersion(), s.GetIOSVersion()))
			totaluserCount += svrinfo.ServerRegInfo.GetUserCount()
			totalRobotCount += svrinfo.ServerRegInfo.GetRobotCount()
		})
		buf.WriteString(fmt.Sprintf("當前服務器總玩家數:%d,總機器人數:%d", totaluserCount, totalRobotCount))

		return buf.String()
	}

	if arglen == 2 && args[1] == "del" {
		ids := args[2]

		if !util.IsNumber(ids) {
			return "error,服務器ID必須為數字"
		}

		id, _ := strconv.Atoi(ids)

		cmd.gameserverlist.RemoveGameServer(int32(id))
	}

	if arglen == 2 && args[0] == "area" {
		s := args[1]

		if !util.IsNumber(s) {
			return "error,服務器ID必須為數字"
		}

		areaID, _ := strconv.Atoi(s)
		var totaluserCount int32 = 0
		var totalRobotCount int32 = 0

		var buf bytes.Buffer
		cmd.gameserverlist.RangeGameServerByAreaID(int64(areaID), func(svrinfo *GameServerInfo) {
			s := svrinfo.ServerRegInfo
			buf.WriteString(fmt.Sprintf("gameid:%d,gamename:%s,serverid:%d,servername:%s,hostandbackup:%v,addr:%s:%d,flag:%s,maxodd:%d,loginmoney:%d,usercount:%d,robotcount:%d,androidver:%s,iosver:%s\r\n", s.GetGameID(), s.GetGameName(), s.GetServerID(), s.GetServerName(), s.GetHostAndBackup(), s.GetIPAddress(), s.GetPort(), s.GetFlag(), s.GetMaxOdd(), s.GetLoginMoney(), s.GetUserCount(), s.GetRobotCount(), s.GetAndroidVersion(), s.GetIOSVersion()))
			totaluserCount += svrinfo.ServerRegInfo.GetUserCount()
			totalRobotCount += svrinfo.ServerRegInfo.GetRobotCount()
		})
		buf.WriteString(fmt.Sprintf("當前服務器總玩家數:%d,總機器人數:%d", totaluserCount, totalRobotCount))

		return buf.String()
	}

	if arglen >= 1 && args[0] == "find" {
		if arglen == 2 {
			gids := args[1]
			if !util.IsNumber(gids) {
				return "error,游戏ID必须为数字"
			}

			id, _ := strconv.Atoi(gids)

			var totaluserCount int32 = 0
			var totalRobotCount int32 = 0

			var buf bytes.Buffer
			cmd.gameserverlist.RangeGameServer(func(svrinfo *GameServerInfo) {
				s := svrinfo.ServerRegInfo
				if s.GetGameID() != int32(id) {
					return
				}

				buf.WriteString(fmt.Sprintf("gameid:%d,gamename:%s,serverid:%d,servername:%s,addr:%s:%d,flag:%s,loginmoney:%d,usercount:%d,robotcount:%d\r\n", s.GetGameID(), s.GetGameName(), s.GetServerID(), s.GetServerName(), s.GetIPAddress(), s.GetPort(), s.GetFlag(), s.GetLoginMoney(), s.GetUserCount(), s.GetRobotCount()))
				totaluserCount += svrinfo.ServerRegInfo.GetUserCount()
				totalRobotCount += svrinfo.ServerRegInfo.GetRobotCount()
			})
			buf.WriteString(fmt.Sprintf("游戏总玩家数:%d,总机器人数:%d", totaluserCount, totalRobotCount))

			return buf.String()

		} else if arglen == 3 {
			gids := args[1]
			flag := args[2]
			if !util.IsNumber(gids) {
				return "error,游戏ID必须为数字"
			}

			id, _ := strconv.Atoi(gids)

			var totaluserCount int32 = 0
			var totalRobotCount int32 = 0

			var buf bytes.Buffer
			cmd.gameserverlist.RangeGameServer(func(svrinfo *GameServerInfo) {
				s := svrinfo.ServerRegInfo
				if s.GetGameID() != int32(id) || s.GetFlag() != flag {
					return
				}

				buf.WriteString(fmt.Sprintf("gameid:%d,gamename:%s,serverid:%d,servername:%s,addr:%s:%d,flag:%s,loginmoney:%d,usercount:%d,robotcount:%d\r\n", s.GetGameID(), s.GetGameName(), s.GetServerID(), s.GetServerName(), s.GetIPAddress(), s.GetPort(), s.GetFlag(), s.GetLoginMoney(), s.GetUserCount(), s.GetRobotCount()))
				totaluserCount += svrinfo.ServerRegInfo.GetUserCount()
				totalRobotCount += svrinfo.ServerRegInfo.GetRobotCount()
			})
			buf.WriteString(fmt.Sprintf("游戏总玩家数:%d,总机器人数:%d", totaluserCount, totalRobotCount))

			return buf.String()
		}
	}

	return "error,see usage ref"

}
