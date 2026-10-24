package roomlogic

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/player"
	"example.invalid/ServerTeam/Servers/pkg/roomutil"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

//用 new *gamelogic.Room 透過RPC 更新 master slave, 參考NewNotifyRefreshVerHandler

// RoomMasterSlaveHandler ?
type RoomMasterSlaveHandler struct {
	room        *RoomLogic
	rpcmsgqueue *rpcmsg.Queue
	gameid      int32
	serverid    int32
}

// NewRoomMasterSlaveHandler ?
func NewRoomMasterSlaveHandler(room *RoomLogic, rpcmsgqueue *rpcmsg.Queue, gameid int32, serverid int32) *RoomMasterSlaveHandler {
	roomHandler := new(RoomMasterSlaveHandler)
	roomHandler.room = room
	roomHandler.rpcmsgqueue = rpcmsgqueue
	roomHandler.gameid = gameid
	roomHandler.serverid = serverid

	return roomHandler
}

// IsHook ?
func (roomHandler *RoomMasterSlaveHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_ServerHostAndBackUp)
}

// OnNotifyMessage ?
func (roomHandler *RoomMasterSlaveHandler) OnNotifyMessage(m *msg.Message, ctx *module.MouduleCtx, playerMang *player.Manager) {
	defer util.PanicErrStack()

	mmbyte := m.MsgData.([]byte)
	mm := &netproto.HostAndBackUp{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		mlog.Error("解析錯誤 HostAndBackUp")
	}
	mlog.Debug("[masterslave] roomHandler[GameID: %d][ServerID: %d], mm[GameID: %d][ServerID: %d]",
		roomHandler.room.room.GameInfo.GameID, roomHandler.room.room.ServerID, *mm.Master.GameID, *mm.Master.ServerID)

	if roomHandler.room.room.ServerID == *mm.Master.ServerID &&
		roomHandler.room.room.GameInfo.GameID == *mm.Master.GameID {
		roomHandler.room.UpdateHostAndBackUp(HostTypeMaster, "")
	} else {
		//mlog.System("收到room房間切換AB服 serverid[%v] ServerAddr[%v]", roomHandler.room.room.ServerID, mm.Slave.GetGameServerAddr())
		roomHandler.room.UpdateHostAndBackUp(HostTypeSlave, mm.Slave.GetGameServerAddr())
	}
	mlog.Debug("[masterslave] roomHandler.HostAndBackup[%d]", roomHandler.room.room.HostAndBackup)

	roomHandler.room.regHallHandler.RegHallServer()

	roomHandler.stopGameServer(mm)
}

func (roomHandler *RoomMasterSlaveHandler) stopGameServer(mm *netproto.HostAndBackUp) {
	mlog.Debug("[masterslave] ===== stopGameServer() =====")
	freeusers := make([]*player.Player, 0)
	roomHandler.room.gameServer.GameServerRunCtx.PlayerManager.Range(func(pr *player.Player) {
		mlog.Debug("[masterslave] pr.User[%d], IsGuest[%t], IsOnline[%t]", pr.UserID, pr.IsGuest(), pr.IsOnline)
		if pr.IsGuest() {
			freeusers = append(freeusers, pr)
		} else {
			ur := pr.UserData.(*gamelogic.GameUser)
			mlog.Debug("[masterslave] ur.User[%d], IsFree[%t], IsSit[%t], IsReady[%t], IsQueue[%t]",
				ur.UserID, ur.IsFree(), ur.IsSit(), ur.IsReady(), ur.IsQueue())
			if ur.IsFree() || ur.IsSit() || ur.IsReady() || ur.IsQueue() {
				freeusers = append(freeusers, pr)
			}
		}
	})
	mlog.Debug("[masterslave] ----- freeusers[%v] -----", freeusers)
	mlog.Debug("[masterslave] freeusers.length[%d]", len(freeusers))
	for _, pr := range freeusers {
		mlog.Debug("[masterslave] UserID[%d], mm.Slave.GetGameServerAddr()[%s], pr.Conn.ConnID()[%s]",
			pr.UserID, mm.Slave.GetGameServerAddr(), pr.Conn.ConnID())
		kickUser(roomHandler.room.gameServer.GameServerRunCtx.MsgServer, mm.Slave.GetGameServerAddr(), pr.Conn.ConnID())
	}
}

func kickUser(ctx *msgserver.MsgServer, str string, connid string) {
	example.invalid("停服檢查邏輯,%s 斷開客戶端%s", str, connid)

	var m *msg.Message
	if str != "" && str != "-" {
		//m, _ = roomutil.GetTipMessageWithParam(str, roomutil.MsgTypeID_RoomKick, 6)
		m, _ = roomutil.GetTipMessageWithParam("-", roomutil.MsgTypeIDRoomKick, 4)
	} else {
		m, _ = roomutil.GetTipMessageWithParam("-", roomutil.MsgTypeIDRoomKick, 4)
	}

	ctx.SendMessageByConnID(connid, m)

	go func() {
		time.Sleep(time.Duration(1) * time.Second)

		ctx.CloseConnByConnID(connid)
	}()
}
