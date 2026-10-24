package serverhandler

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/event"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

const (
	// TimerNameNotifyReg 註冊Notify定時器名稱
	TimerNameNotifyReg = "gameserver_serverhandler_notifyreg"
)

// NotifyHandler 通知協調處理器
type NotifyHandler struct {
	nsAddr         string
	TCPMsgClient   *msgserver.TCPMsgClient
	evtchan        *msgserver.MsgEvtChan
	stopsig        chan bool
	serverID       int
	serverType     int
	gameid         int
	HostAndBackup  int32
	cgx            *GameServerCtx
	gameServerAddr string
	FlagID         string
	nfRegTicker    *event.TimerFunc
}

// NewNotifyHandler 創建通知協調處理器
func NewNotifyHandler(nsaddr string, serverid int, servertype int, gameid int, hostAndBackup int32, gameserverAddr string, flagID string) *NotifyHandler {
	ntfhandler := new(NotifyHandler)
	ntfhandler.nsAddr = nsaddr
	ntfhandler.stopsig = make(chan bool, 1)
	ntfhandler.serverID = serverid
	ntfhandler.serverType = servertype
	ntfhandler.gameid = gameid
	ntfhandler.HostAndBackup = hostAndBackup
	ntfhandler.gameServerAddr = gameserverAddr
	ntfhandler.FlagID = flagID

	return ntfhandler
}

// Start 啟動
func (nshandler *NotifyHandler) Start(cgx *GameServerCtx) {
	// 創建到NsServer的TCP連接
	msgparser := msg.NewProtoMsgParser()
	msgpacker := network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)

	nshandler.evtchan = msgserver.NewMsgEvtChan(100)
	nshandler.TCPMsgClient = msgserver.NewTCPMsgClient(msgpacker, nshandler.evtchan, msgparser, true)

	nshandler.cgx = cgx

	go nshandler.tcpEventChanFunc(nshandler.evtchan)

	err := nshandler.TCPMsgClient.Connect(nshandler.nsAddr)

	if err != nil {
		mlog.Warn("向協調服務器 %s 註冊失敗, %v", nshandler.nsAddr, err)
		return
	}

	//nshandler.nfRegTicker = cgx.TimerFuncManager.NewNamedTicker(30, nshandler.sendMsgNotify, TimerNameNotifyReg)
}

func (nshandler *NotifyHandler) sendMsgNotify() {
	data := []byte("test")
	nshandler.TCPMsgClient.SendData(data)
}

func (nshandler *NotifyHandler) tcpEventChanFunc(tcpevt *msgserver.MsgEvtChan) {

	defer util.PanicErrStack()

	for {
		select {
		case <-nshandler.stopsig:
			mlog.Warn("notifyserver 消息處理線程退出")
			return
		default:
		}

		select {
		case clt := <-tcpevt.ConnectedSig:
			example.invalid("連接協調服務器成功, %s", clt.GetConn().RemoteAddr().String())

			// 發送服務器信息
			nsregmsg := new(msg.Message)
			nsregmsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
			nsregmsg.SClassID = int32(netproto.NotifyServerClassID_ServerRegID)

			nsregdata := new(netproto.NotifyServerRegInfo)
			nsregmsg.MsgData = nsregdata
			nsregdata.ServerID = proto.Int32(int32(nshandler.serverID))
			nsregdata.ServerType = proto.Int32(int32(nshandler.serverType))
			nsregdata.HostAndBackup = proto.Int32(nshandler.HostAndBackup)
			nsregdata.GameServerAddr = proto.String(nshandler.gameServerAddr)
			nsregdata.FlagID = proto.String(nshandler.FlagID)

			if nshandler.gameid > 0 && nshandler.serverType == conf.ServerTypeGameServer {
				nsregdata.GameID = proto.Int32(int32(nshandler.gameid))
			}

			nshandler.TCPMsgClient.SendMessage(nsregmsg)

		case clt := <-tcpevt.DisconnectedSig:
			mlog.Warn("協調服務器連接斷開, %v", clt.GetConn().RemoteAddr().String())
		case <-tcpevt.DataSig:

		case msg := <-tcpevt.MsgSig:
			mlog.Debug("收到協調服務器消息: %v", msg)

			// 分發通知事件
			notifymsgevt := msgserver.NewNotifyMessageEvent(msg)
			nshandler.cgx.MsgServer.EventChannel.OnAsyncEvent(notifymsgevt, nil)
		}
	}
}

// Stop 停止
func (nshandler *NotifyHandler) Stop() {
	nshandler.stopsig <- true
}

// Name 名稱
func (nshandler *NotifyHandler) Name() string {
	return "notify"
}
