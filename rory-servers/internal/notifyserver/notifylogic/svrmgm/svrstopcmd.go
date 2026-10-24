package svrmgm

import (
	"strconv"
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// SvrStopCmd 停服
type SvrStopCmd struct {
	msgServer  *msgserver.MsgServer
	serverlist *ServerList
}

// NewSvrStopCmd 創建停服命令對象
func NewSvrStopCmd(msgSvr *msgserver.MsgServer, serverlist *ServerList) *SvrStopCmd {
	sscmd := new(SvrStopCmd)
	sscmd.msgServer = msgSvr
	sscmd.serverlist = serverlist

	return sscmd
}

// Name ?
func (stopCmd *SvrStopCmd) Name() string {
	return "stopserver"
}

// Help ?
// server game start 9 21
// server game stop 9 21
func (stopCmd *SvrStopCmd) Help() string {
	return "停服.usage:stopserver type(0 全服 1 遊戲 2 房間) status(0 停服 1 开服) gameid serverid msg"
}

// Exec ?
func (stopCmd *SvrStopCmd) Exec(args []string) string {
	if len(args) < 3 {
		return "error,seeref usage"
	}

	//0 全服 1 遊戲 2 房間
	sstype := args[0]
	//0 停服 1 开服
	ssStatus := args[1]

	// 檢查狀態值
	if ssStatus != "0" && ssStatus != "1" {
		return "error,seeref usage"
	}

	status, _ := strconv.Atoi(ssStatus)

	//构造停服消息
	svrstopMsg := new(msg.Message)
	svrstopMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	svrstopMsg.SClassID = int32(netproto.NotifyServerClassID_ServerStopID)

	svrstopData := new(netproto.ServerStop)
	svrstopData.Type = proto.Int32(0)
	svrstopData.Status = proto.Int32(int32(status))

	// 全服
	if sstype == "0" {

		stopmsg := strings.Join(args[2:], " ")
		svrstopData.Message = proto.String(stopmsg)
		svrstopMsg.MsgData = svrstopData

		stopCmd.msgServer.SendAllMessage(svrstopMsg)

	} else if sstype == "1" { // 遊戲
		if len(args) < 4 {
			return "error,seeref usage"
		}

		ssGameid := args[2]
		stopmsg := strings.Join(args[3:], " ")

		gameid, err := strconv.Atoi(ssGameid)

		if err != nil {
			return "error,seeref usage"
		}
		// 設置遊戲ID
		svrstopData.GameID = proto.Int32(int32(gameid))
		// 停服消息
		svrstopData.Message = proto.String(stopmsg)
		svrstopMsg.MsgData = svrstopData

		svrs := stopCmd.serverlist.GetServerByGameID(gameid)

		if svrs != nil && len(svrs) > 0 {
			for _, svr := range svrs {
				stopCmd.msgServer.SendMessage(svr.ConnInfo, svrstopMsg)
			}
		}

	} else if sstype == "2" { // 房間
		if len(args) < 5 {
			return "error,seeref usage"
		}

		// 遊戲ID
		ssGameid := args[2]
		// 服務器ID
		ssServerid := args[3]
		// 停服消息
		stopmsg := strings.Join(args[4:], " ")

		serverid, err := strconv.Atoi(ssServerid)
		if err != nil {
			return "error,seeref usage"
		}

		gameid, err := strconv.Atoi(ssGameid)

		if err != nil {
			return "error,seeref usage"
		}

		//设置游戏ID
		svrstopData.GameID = proto.Int32(int32(gameid))
		//设置服务器ID
		svrstopData.ServerID = proto.Int32(int32(serverid))
		//停服消息
		svrstopData.Message = proto.String(stopmsg)
		svrstopMsg.MsgData = svrstopData

		svr := stopCmd.serverlist.FindServer(conf.ServerTypeGameServer, serverid)

		if svr != nil {
			stopCmd.msgServer.SendMessage(svr.ConnInfo, svrstopMsg)
		}
	}

	return "server stop msg sended..."
}
