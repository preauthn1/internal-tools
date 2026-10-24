package serverhandler

import (
	"strings"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/gamelogic"
	"example.invalid/ServerTeam/Servers/pkg/gameserver/serverhandler"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/robot/robotcomm"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	// TimerNameHallServerReg 註冊大廳定時器名稱
	TimerNameHallServerReg = "roomserver_serverhandler_hallserverreg"
)

// HallServerRegHandler 大廳註冊處理器
type HallServerRegHandler struct {
	stopSig        chan bool
	hallserverAddr []string
	tcpclients     []*msgserver.TCPMsgClient
	room           *gamelogic.Room
	robotManager   *robotcomm.RobotManager
	regInterval    int32
	regHallTicker  *event.TimerFunc
	ver            *netproto.GameVerInfo
}

// NewHallServerRegHandler 創建大廳註冊處理器
func NewHallServerRegHandler(addr string, room *gamelogic.Room, regInterval int32, robotManager *robotcomm.RobotManager) *HallServerRegHandler {
	h := new(HallServerRegHandler)
	h.stopSig = make(chan bool, 1)
	h.room = room
	h.regInterval = regInterval
	h.robotManager = robotManager

	// 解析大廳地址
	addrArr := strings.Split(addr, ",")

	for _, saddr := range addrArr {
		if util.IsNullOrEmpty(saddr) {
			continue
		}

		h.hallserverAddr = append(h.hallserverAddr, saddr)
	}

	return h
}

// Name 名稱
func (hreg *HallServerRegHandler) Name() string {
	return "hallserverreg"
}

// Start 啟動
func (hreg *HallServerRegHandler) Start(ctx *serverhandler.GameServerCtx) {
	evtchan := msgserver.NewMsgEvtChan(conf.EVTCHANSIZE)

	// 開始連接hallserver
	for _, addr := range hreg.hallserverAddr {
		msgpacker := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
		msgparser := msg.NewProtoMsgParser()
		clt := msgserver.NewTCPMsgClient(msgpacker, evtchan, msgparser, true)

		if err := clt.Connect(addr); err != nil {
			mlog.Error("向大廳 %s 註冊失敗, %v", addr, err)
		}

		hreg.tcpclients = append(hreg.tcpclients, clt)
	}

	go hreg.tcpEventChanFunc(evtchan)

	// 啟動時先註冊一次
	hreg.RegHallServer()
	// 定時註冊
	hreg.regHallTicker = ctx.TimerFuncManager.NewNamedTicker(hreg.regInterval, hreg.RegHallServer, TimerNameHallServerReg)

}

// Stop 停止
func (hreg *HallServerRegHandler) Stop() {
	hreg.stopSig <- true
	if hreg.regHallTicker != nil {
		hreg.regHallTicker.Stop()
	}
}

func (hreg *HallServerRegHandler) tcpEventChanFunc(ch *msgserver.MsgEvtChan) {
	defer util.PanicErrStack()

	for {
		select {
		case <-hreg.stopSig:
			mlog.Error("hallserver 註冊消息處理線程退出")
			break
		default:
		}

		select {
		case clt := <-ch.ConnectedSig:
			example.invalid("connect hallserver,%s", clt.GetConn().RemoteAddr().String())
		case clt := <-ch.DisconnectedSig:
			example.invalid("disconnect hallserver,%v", clt.GetConn().RemoteAddr().String())
		case <-ch.DataSig:

		case msg := <-ch.MsgSig:
			if hreg.IsHook(msg.BClassID, msg.SClassID) {
				hreg.OnMessage(msg)
			}
		}
	}
}

// RefreshGameVer 更新版號後向大廳更新資料
func (hreg *HallServerRegHandler) RefreshGameVer(ver *netproto.GameVerInfo) {
	hreg.ver = ver
	hreg.RegHallServer()
}

// RegHallServer 遊戲註冊大廳
func (hreg *HallServerRegHandler) RegHallServer() {
	defer util.PanicErrStack()

	sinfo := new(netproto.GameServerRegInfo)
	sinfo.ServerID = proto.Int32(hreg.room.ServerID)
	sinfo.GameID = proto.Int32(hreg.room.GameInfo.GameID)
	sinfo.ServerName = proto.String(hreg.room.ServerName)
	sinfo.GameName = proto.String(hreg.room.GameInfo.GameName)
	sinfo.Flag = proto.String(hreg.room.Flag)
	sinfo.GameType = proto.Int32(hreg.room.GameType)
	sinfo.IPAddress = proto.String(hreg.room.IPAddress)
	sinfo.Port = proto.Int32(hreg.room.Port)
	sinfo.UserCount = proto.Int32(int32(hreg.room.RoomUserManager.UserCount()))
	sinfo.LoginMoney = proto.Int32(hreg.room.LoginMoney)
	sinfo.BaseScore = proto.Int32(hreg.room.BaseScore)
	sinfo.RobotCount = proto.Int32(int32(hreg.robotManager.RobotCount()))
	sinfo.Single = proto.Int32(int32(hreg.room.Single))
	sinfo.HostAndBackup = proto.Int32(int32(hreg.room.HostAndBackup))
	sinfo.AndroidVersion = proto.String(hreg.ver.GetAndroidVersion())
	sinfo.IOSVersion = proto.String(hreg.ver.GetIOSVersion())
	sinfo.KickOutMoney = proto.Int32(hreg.room.KickOutMoney)
	sinfo.MaxOdd = proto.Int32(hreg.room.MaxOdd)
	data, _ := proto.Marshal(sinfo)

	m := msg.NewMessage(int32(netproto.MessageBClassID_ServerMgm), int32(netproto.ServerMgmClassID_HallServerRegister), data)

	for _, clt := range hreg.tcpclients {
		if clt.IsConnected() {
			// mlog.Debug("向大廳 %s 註冊", clt.RemoteAddr().String())
			clt.SendMessage(m)
		}
	}

}

// IsHook 是否匹配
func (hreg *HallServerRegHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_ServerMgm) && sclsid == int32(netproto.NotifyServerClassID_ServerHostAndBackUp)
}

// OnMessage 收到封包
func (hreg *HallServerRegHandler) OnMessage(m *msg.Message) {
	defer util.PanicErrStack()

	mmbyte := m.MsgData.([]byte)
	mm := &netproto.HostAndBackUp{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		mlog.Error("Hall HostAndBackUp Error : %v", err.Error())
		return
	}

	mlog.Debug("[masterslave] hreg.ServerID[%d], mm.GetServerID()[%d]", hreg.room.ServerID, mm.Slave.GetServerID())
	if hreg.room.ServerID == mm.Slave.GetServerID() {
		hreg.room.HostAndBackup = 2
	}
	mlog.Debug("[masterslave] hreg.room.HostAndBackup[%d]", hreg.room.HostAndBackup)

	hreg.RegHallServer()
}
