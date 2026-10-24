package serverhandler

import (
	proto "example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/event"
	gshandler "example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
const (
	TimerNameServerStop = "hallserver_serverhandler_serverstophandler_ticker"
)

// ServerStopHandler 服務器狀態管理器
type ServerStopHandler struct {
	statusProvider  logicmodule.ServerStatusProvider
	svrStopTicker   *event.TimerFunc
	sendMsgTimeSpan int
}

// NewServerStopHandler 創建服務器狀態管理器
func NewServerStopHandler(statusProvider logicmodule.ServerStatusProvider, timespan int) *ServerStopHandler {
	stm := new(ServerStopHandler)
	stm.statusProvider = statusProvider
	stm.sendMsgTimeSpan = timespan

	return stm
}

// Name ?
func (stm *ServerStopHandler) Name() string {
	return "hallserverstop"
}

// Stop ?
func (stm *ServerStopHandler) Stop() {
	if stm.svrStopTicker != nil {
		stm.svrStopTicker.Stop()
	}
}

// Start ?
func (stm *ServerStopHandler) Start(ctx *gshandler.GameServerCtx) {
	defer util.PanicErrStack()

	sendTipMsgFunc := func() {
		if !stm.statusProvider.IsServerRuning() {
			example.invalid("廣播服務器停服消息,%s", stm.statusProvider.GetStopServerTipMsg())

			tipmsg := new(msg.Message)
			tipmsg.BClassID = int32(netproto.MessageBClassID_Common)
			tipmsg.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

			tipdata := new(netproto.TipMessage)
			tipdata.MsgType = proto.Int32(3)
			tipdata.Message = proto.String(stm.statusProvider.GetStopServerTipMsg())
			tipmsg.MsgData = tipdata

			ctx.MsgServer.RangeConn(func(clt *network.ConnInfo) {
				if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) && !stm.statusProvider.IsWhiteIP(clt.IPAddress) {
					ctx.MsgServer.SendMessage(clt, tipmsg)
				}
			})
		}

		// 如果服務器關閉
		if stm.statusProvider.GetServerStatus() == conf.ServerStatusClose {
			ctx.MsgServer.RangeConn(func(clt *network.ConnInfo) {
				if clt.IsServerConn(conf.TCPSERVERNAMEHALLCLIENT) && !stm.statusProvider.IsWhiteIP(clt.IPAddress) {
					go func() {
						ctx.MsgServer.CloseConn(clt)
					}()
				}
			})
		}
	}

	stm.svrStopTicker = ctx.TimerFuncManager.NewNamedTicker(int32(stm.sendMsgTimeSpan), sendTipMsgFunc, TimerNameServerStop)
}
