package serverhandler

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	gshandler "example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/roomutil"
)

// ?
const (
	TimerNameServerStop = "roomserver_serverhandler_serverstophandler_ticker"
)

// ServerStopHandler 停服服務器處理器
type ServerStopHandler struct {
	serverStatus        logicmodule.ServerStatusProvider
	kickUserTimeSpan    int32
	cgx                 *gshandler.GameServerCtx
	svrStopNotifyTicker *event.TimerFunc
	room                *gamelogic.Room
}

// NewServerStopHandler 創建停服處理器對象
func NewServerStopHandler(serverStatus logicmodule.ServerStatusProvider, kickUserTimeSpan int32, room *gamelogic.Room) *ServerStopHandler {
	sth := new(ServerStopHandler)
	sth.kickUserTimeSpan = kickUserTimeSpan
	sth.serverStatus = serverStatus
	sth.room = room

	return sth
}

// Name ?
func (stm *ServerStopHandler) Name() string {
	return "gameserverstop"
}

// Stop ?
func (stm *ServerStopHandler) Stop() {
	if stm.svrStopNotifyTicker != nil {
		stm.svrStopNotifyTicker.Stop()
	}
}

// Start ?
func (stm *ServerStopHandler) Start(ctx *gshandler.GameServerCtx) {
	stm.cgx = ctx
	stm.svrStopNotifyTicker = ctx.TimerFuncManager.NewNamedTicker(stm.kickUserTimeSpan, func() {
		if stm.serverStatus.IsServerRuning() {
			return
		}

		// 廣播停服消息
		//stm.boardCastServerStopMsg(ctx)

		//example.invalid("停服檢查空閒用戶邏輯開始運行...")
		freeusers := make([]*player.Player, 0)
		ctx.PlayerManager.Range(func(pr *player.Player) {
			if pr.IsGuest() {
				freeusers = append(freeusers, pr)
			} else {
				ur := pr.UserData.(*gamelogic.GameUser)
				if ur.IsFree() || ur.IsSit() || ur.IsReady() || ur.IsQueue() {
					freeusers = append(freeusers, pr)
				}
			}
		})

		for _, pr := range freeusers {
			stm.kickUser(ctx, "", pr.Conn.ConnID())
		}

	}, TimerNameServerStop)

}

func (stm *ServerStopHandler) kickUser(ctx *gshandler.GameServerCtx, str string, connid string) {
	mlog.Warn("停服檢查邏輯,%s 斷開客戶端%s", str, connid)

	var m *msg.Message
	if stm.room.NewGameServerAddr != "" && stm.room.NewGameServerAddr != "-" {
		//m, _ = roomutil.GetTipMessageWithParam(stm.room.NewGameServerAddr , roomutil.MsgTypeID_RoomKick, 6)
		m, _ = roomutil.GetTipMessageWithParam("-", roomutil.MsgTypeIDRoomKick, 4)
	} else {
		m, _ = roomutil.GetTipMessageWithParam("-", roomutil.MsgTypeIDRoomKick, 4)
	}

	ctx.MsgServer.SendMessageByConnID(connid, m)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		ctx.MsgServer.CloseConnByConnID(connid)
	}()
}

// 廣播停服消息
func (stm *ServerStopHandler) boardCastServerStopMsg(ctx *gshandler.GameServerCtx) {
	tipmsg := new(msg.Message)
	tipmsg.BClassID = int32(netproto.MessageBClassID_Common)
	tipmsg.SClassID = int32(netproto.PlatformCommonClassID_TipMessageID)

	tipdata := new(netproto.TipMessage)
	tipdata.MsgType = proto.Int32(3)
	tipdata.Message = proto.String(stm.serverStatus.GetStopServerTipMsg())
	tipmsg.MsgData = tipdata

	ctx.MsgServer.SendAllMessage(tipmsg)
}
