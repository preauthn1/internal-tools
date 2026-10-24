package gamelogic

import (
	"net/http"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// SuperToolBrGameServiceBase 百人遊戲超端服務基類
type SuperToolBrGameServiceBase struct {
	Room             *Room
	tfManager        *event.TimerFuncManager
	DBMsgQueue       *rpcmsg.Queue
	UpdateTicker     *event.TimerFunc
	superToolHandler ISuperToolHandler
	RefreshInterval  int32
}

// NewBrGameSuperToolServiceBase 百人遊戲超端服務基類
func NewBrGameSuperToolServiceBase(room *Room, tfm *event.TimerFuncManager, dbMsgQueue *rpcmsg.Queue) ISuperToolService {
	service := new(SuperToolBrGameServiceBase)
	service.Room = room
	service.tfManager = tfm
	service.DBMsgQueue = dbMsgQueue
	service.superToolHandler = NewBrGameSuperToolHandlerBase(service)
	service.RefreshInterval = 30
	return service
}

// Start ?
func (service *SuperToolBrGameServiceBase) Start() {
	service.RefreshVIPUser()
	service.UpdateTicker = service.tfManager.NewNamedTicker(30, service.RefreshVIPUser, "timer_refresh_vip_user")
	service.superToolHandler.SysPoolStart() //啟動時要一次系統水池
}

// Stop ?
func (service *SuperToolBrGameServiceBase) Stop() {
	service.UpdateTicker.Stop()
}

// GetSuperToolHandler ?
func (service *SuperToolBrGameServiceBase) GetSuperToolHandler() ISuperToolHandler {
	return service.superToolHandler
}

// GetBrGameSuperToolHandler ?
func (service *SuperToolBrGameServiceBase) GetBrGameSuperToolHandler() *SuperToolBrGameHandlerBase {
	return service.superToolHandler.(*SuperToolBrGameHandlerBase)
}

// RefreshVIPUser ?
func (service *SuperToolBrGameServiceBase) RefreshVIPUser() {
	gameid := service.Room.GameInfo.GameID
	serverid := service.Room.ServerID
	pmsg := new(msg.Message)
	pmsg.BClassID = int32(netproto.MessageBClassID_Game)
	pmsg.SClassID = int32(netproto.GameMessageClassID_SuperToolGetVipUserID)
	pdata := new(netproto.SuperToolGetVipUserRequest)
	pdata.GameID = proto.Int32(gameid)
	pmsg.MsgData = pdata

	rpcinfo := netproto.CreateEmptyRPCInfo(serverid)
	// 設置RPC消息桌號
	service.DBMsgQueue.PushMessage(rpcinfo, false, pmsg)
}

// GetRoomServerID ?
func (service *SuperToolBrGameServiceBase) GetRoomServerID() int32 {
	return service.Room.ServerID
}

// PushRPCMessage ?
func (service *SuperToolBrGameServiceBase) PushRPCMessage(rpcinfo *netproto.RPCInfo, safe bool, msg *msg.Message) error {
	return service.DBMsgQueue.PushMessage(rpcinfo, safe, msg)
}

// HTTPService ?
func (service *SuperToolBrGameServiceBase) HTTPService(w http.ResponseWriter, req *http.Request) {}

// GetGameID 獲取遊戲ID
func (service *SuperToolBrGameServiceBase) GetGameID() int32 {
	return service.Room.GameInfo.GameID
}

// GetFlag 獲取房間等級
func (service *SuperToolBrGameServiceBase) GetFlag() string {
	return service.Room.Flag
}
