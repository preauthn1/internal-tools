package roomlogic

import (
	"os"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/module/logicmodule"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/serverhandler"
)

// GameServerReg 遊戲服務器註冊
type GameServerReg struct {
	clientAddr   string
	notifyAddr   string
	httpAddr     string
	serverid     int32
	rpcmsgqueue  *rpcmsg.Queue
	room         *gamelogic.Room
	serverStatus logicmodule.ServerStatusProvider

	gameServerRegRet *netproto.GameServerStartRet
	regHallHandler   *serverhandler.HallServerRegHandler
}

// NewGameServerReg 創建遊戲服務器註冊處理器對象
func NewGameServerReg(clientAddr string, notifyAddr string, httpAddr string, serverid int32, dbmsgqueue *rpcmsg.Queue, room *gamelogic.Room, serverStatus logicmodule.ServerStatusProvider) *GameServerReg {
	gsreg := new(GameServerReg)
	gsreg.clientAddr = clientAddr
	gsreg.notifyAddr = notifyAddr
	gsreg.httpAddr = httpAddr
	gsreg.serverid = serverid
	gsreg.rpcmsgqueue = dbmsgqueue
	gsreg.room = room
	gsreg.serverStatus = serverStatus

	return gsreg
}

// SetRegHallHandler 將 regHallHandler 配置給 遊戲註冊器使用
func (hsreg *GameServerReg) SetRegHallHandler(regHallHandler *serverhandler.HallServerRegHandler) {
	hsreg.regHallHandler = regHallHandler
}

// IsHook 是否匹配
func (hsreg *GameServerReg) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_DBServer) && sclsid == int32(netproto.DBServerClassID_GameServerStartRetID)
}

// OnRPCMessage 處理訊息
func (hsreg *GameServerReg) OnRPCMessage(m *msg.Message, ctx *module.MouduleCtx) {
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.GameServerStartRet{}
		return mm
	})

	mm := rmm.(*netproto.GameServerStartRet)

	hsreg.gameServerRegRet = mm

	example.invalid("收到 RPC Message Queue 返回註冊信息. Code: %d, Message: %s, Status: %d, Andriod[%s] IOS[%s]", mm.GetCode(), mm.GetMessage(), mm.GetStatus(), mm.GetAndroidVersion(), mm.GetIOSVersion())

	if mm.GetCode() == -1 {
		mlog.Error("向 RPC Message Queue 註冊失敗. %s", mm.GetMessage())
		os.Exit(1)
	} else {
		example.invalid("向 RPC Message Queue 註冊成功")
		//收到 RPC Message Queue 注冊返回後即時向HallServer更新客戶端資料
		gv := new(netproto.GameVerInfo)
		gv.AndroidVersion = hsreg.gameServerRegRet.AndroidVersion
		gv.IOSVersion = hsreg.gameServerRegRet.IOSVersion
		hsreg.regHallHandler.RefreshGameVer(gv)
	}

	// 服務器未開放
	if mm.GetStatus() == 0 {
		hsreg.serverStatus.SetServerStatus(conf.ServerStatusRecover, mm.GetStatusMessage())
	}
}

// SendDBServerStartMsg 通知服務器啟動
func (hsreg *GameServerReg) SendDBServerStartMsg() {
	svrmsg := new(msg.Message)
	svrmsg.BClassID = int32(netproto.MessageBClassID_DBServer)
	svrmsg.SClassID = int32(netproto.DBServerClassID_GameServerStartID)

	svrdata := new(netproto.GameServerStart)
	svrdata.ClientAddr = proto.String(hsreg.clientAddr)
	svrdata.NotifyAddr = proto.String(hsreg.notifyAddr)
	svrdata.ServerID = proto.Int32(hsreg.serverid)
	svrdata.GameID = proto.Int32(hsreg.room.GameInfo.GameID)
	svrdata.LoginMoney = proto.Int32(hsreg.room.LoginMoney)
	svrdata.ServerName = proto.String(hsreg.room.ServerName)
	svrdata.ShowName = proto.String(hsreg.room.Flag)
	svrdata.HttpAddr = proto.String(hsreg.httpAddr)

	svrmsg.MsgData = svrdata

	hsreg.rpcmsgqueue.PushMessage(netproto.CreateEmptyRPCInfo(hsreg.serverid), false, svrmsg)

}

// GameVerProvider ?
type GameVerProvider interface {
	GetVer() *netproto.GameVerInfo
}

// RPCGameVerProvider ?
type RPCGameVerProvider struct {
	gsreg *GameServerReg
}

// NewRPCGameVerProvider ?
func NewRPCGameVerProvider(gsreg *GameServerReg) *RPCGameVerProvider {
	gvp := new(RPCGameVerProvider)
	gvp.gsreg = gsreg

	return gvp
}

// GetVer ?
func (gsv *RPCGameVerProvider) GetVer() *netproto.GameVerInfo {
	gv := new(netproto.GameVerInfo)
	if gsv.gsreg.gameServerRegRet != nil {
		gv.AndroidVersion = gsv.gsreg.gameServerRegRet.AndroidVersion
		gv.IOSVersion = gsv.gsreg.gameServerRegRet.IOSVersion
	}

	return gv
}
