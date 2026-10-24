package svrmgm

import (
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// RefreshVerCMD 刷新版本號命令
type RefreshVerCMD struct {
	msgServer  *msgserver.MsgServer
	serverlist *ServerList
}

// NewRefreshVerCmd 創建平台版本號刷新命令
func NewRefreshVerCmd(msgServer *msgserver.MsgServer, serverlist *ServerList) *RefreshVerCMD {
	refreshVerCMD := new(RefreshVerCMD)
	refreshVerCMD.msgServer = msgServer
	refreshVerCMD.serverlist = serverlist

	return refreshVerCMD
}

// Name ?
func (verCmd *RefreshVerCMD) Name() string {
	return "refreshver"
}

// Help ?
func (verCmd *RefreshVerCMD) Help() string {
	return "強制刷新平台版本號.usage:refreshver"
}

// Exec ?
func (verCmd *RefreshVerCMD) Exec(args []string) string {
	refreshVerMsg := new(msg.Message)
	refreshVerMsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	refreshVerMsg.SClassID = int32(netproto.NotifyServerClassID_RefreshVersionID)

	refreshVerMsg.MsgData = nil

	hsLst := verCmd.serverlist.GetHallServer()
	if hsLst != nil {
		for _, hs := range hsLst {
			example.invalid("向大厅%v发送更新版本号消息", hs)
			verCmd.msgServer.SendMessage(hs.ConnInfo, refreshVerMsg)
		}
	}

	gsLst := verCmd.serverlist.GetGameServer()
	if gsLst != nil {
		for _, gs := range gsLst {
			example.invalid("向遊戲%v发送更新版本号消息", gs)
			verCmd.msgServer.SendMessage(gs.ConnInfo, refreshVerMsg)
		}
	}

	return "刷新版本号成功"
}
