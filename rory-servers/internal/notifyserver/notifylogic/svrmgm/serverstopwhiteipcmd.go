package svrmgm

import (
	"strings"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// ServerStopWhiteIPCmd 停服IP白名單
type ServerStopWhiteIPCmd struct {
	msgServer  *msgserver.MsgServer
	serverlist *ServerList
}

// NewServerStopWhiteIPCmd 創建停服IP白名單命令
func NewServerStopWhiteIPCmd(msgServer *msgserver.MsgServer, serverlist *ServerList) *ServerStopWhiteIPCmd {
	cmd := new(ServerStopWhiteIPCmd)
	cmd.msgServer = msgServer
	cmd.serverlist = serverlist

	return cmd
}

// Name ?
func (wip *ServerStopWhiteIPCmd) Name() string {
	return "whiteip"
}

// Help /
func (wip *ServerStopWhiteIPCmd) Help() string {
	return "停服白名單.usage:whiteip 0.0.0.0,0.0.0.0 ..."
}

// Exec /
func (wip *ServerStopWhiteIPCmd) Exec(args []string) string {
	if len(args) == 0 {
		return "error,seeref usage"
	}

	ipArr := strings.Split(args[0], ",")

	if len(ipArr) == 0 {
		return "error.白名單IP地址不正確"
	}

	wipMsg := new(msg.Message)
	wipMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	wipMsg.SClassID = int32(netproto.NotifyServerClassID_ServerStopWhiteIPListID)

	wipData := new(netproto.ServerStopWhiteIPList)
	for _, ips := range ipArr {
		if ips == "" {
			continue
		}

		wipData.IP = append(wipData.IP, ips)
	}
	wipMsg.MsgData = wipData

	hslst := wip.serverlist.GetAllServer()
	if hslst != nil {
		for _, hs := range hslst {
			example.invalid("通知服務器%v刷新IP白名單", hs)
			wip.msgServer.SendMessage(hs.ConnInfo, wipMsg)
		}
	}

	return "刷新IP白名單成功"
}
