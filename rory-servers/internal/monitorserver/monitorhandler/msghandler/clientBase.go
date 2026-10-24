package msghandler

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"
	"example.invalid/ServerTeam/Servers/internal/monitorserver/game/by"
	"example.invalid/ServerTeam/Servers/internal/monitorserver/monitorutil"

	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// ?
var (
	SelfIP string
	locker = &sync.RWMutex{}
)

// ClientBase ?
type ClientBase struct {
	ConnClt          *msgserver.TCPMsgClient
	TCPAddr          *net.TCPAddr
	StopSig          chan bool
	MsgParser        *msg.ParserGroup
	MsgPacker        *network.HeadMsgPacker
	Timeout          int
	MsgHandler       *clientframework.HookMsgHandlerGroup
	IPReported       bool
	Running          bool
	EnableHeartBeat  bool
	Token            string
	LoginGameRoomRet *netproto.LoginGameRoomRet
}

// Start ?
func (cb *ClientBase) Start() error {
	evtchan := msgserver.NewMsgEvtChan(100)
	cb.Running = true
	cb.MsgPacker = network.NewHeadMsgPacker(conf.MSGPACKERBUF, true)
	cb.registerMsgParser()
	//cb.RegisterMsgHandler()

	tcpclt := msgserver.NewTCPMsgClient(cb.MsgPacker, evtchan, cb.MsgParser, false)
	err := tcpclt.Connect(cb.TCPAddr.String()) //先連線一次

	if err != nil {
		mlog.Error("%v", err)
		return err
	}
	cb.ConnClt = tcpclt
	cb.EnableHeartBeat = false

	go cb.tcpEventChanFunc(evtchan)

	go func() {
		defer util.PanicErrStack()
		for {
			time.Sleep(time.Millisecond * time.Duration(1))
			select {
			case <-cb.StopSig:
				mlog.System("終止自动重连服务器%s", cb.TCPAddr.String())
				return
			default:
			}
			if !tcpclt.IsConnected() && cb.Running {
				mlog.System("开始自动重连服务器%s", cb.TCPAddr.String())
				tcpclt.Connect(cb.TCPAddr.String())
			}
		}
	}()

	if !cb.WaitConnect() {
		errmsg := fmt.Sprintf("無法連上服務器%s, IsConnected %v, IpReport %v, Running %v", cb.TCPAddr.String(), tcpclt.IsConnected(), cb.IPReported, cb.Running)
		mlog.Error(errmsg)
		return errors.New(errmsg)
	}
	return nil
}

// Stop ?
func (cb *ClientBase) Stop() {
	cb.Running = false
	cb.StopSig <- true
	if cb.ConnClt != nil {
		cb.ConnClt.CloseAndStop()
	}
}

// RegisterMsgHandler 註冊所有消息handler
func (cb *ClientBase) RegisterMsgHandler() {
	cb.MsgHandler = clientframework.NewHookMsgHandlerGroup()
	mh := cb.MsgHandler
	mh.AddHandler(NewHeartBeatConfigHandler())
	mh.AddHandler(NewTipMessageHandler())

	// 遊戲消息註冊
	mh.AddHandler(by.NewBySettingHandler())
	mh.AddHandler(NewGameStartHandler())
	mh.AddHandler(by.NewUserRecoverRetHandler())
	mh.AddHandler(by.NewFrameS2CHandler())
}

func (cb *ClientBase) registerMsgParser() {
	cb.MsgParser = msg.NewParserGroup()
	msgparser := cb.MsgParser
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_TipMessageID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_RequestHeartBeatConfigID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_HeartBeatConfigID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_HeartBeatID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_HeartBeatReturnID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_IPReportID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_ConnTokenID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_NotifyServer), int32(netproto.NotifyServerClassID_AnnID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestHallVerID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_HallVerInfoID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GuestLoginID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_LoginRetID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestServerListID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_ServerListDataID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_RequestGameServerAddrID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Hall), int32(netproto.HallMsgClassID_GameServerAddrInfoID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomRetID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_LoginRoomID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.GameRoomClassID_GameVerInfoID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.GameRoomClassID_GameReadyID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_GameRoom), int32(netproto.GameRoomClassID_GameReadyID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_LookOnUserNotifyID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_BYSettingID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.GameMessageClassID_GameStartID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_UserRecoverRetID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_FrameS2CID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_FrameC2SID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_BYNotifyFish), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_NotifyServer), int32(netproto.NotifyServerClassID_AttrChangeID), msg.NewProtoMsgParser())
	msgparser.RegisterParser(int32(netproto.MessageBClassID_Game), int32(netproto.BY_GameMessageClassID_BYLeaveTip), msg.NewProtoMsgParser())
}

func (cb *ClientBase) ipReport() error {
	localIP := SelfIP
	//localIP = "0.0.0.0"
	msgData := new(netproto.IPReport)
	msgData.IPAddress = proto.String(localIP)
	msgByte, _ := proto.Marshal(msgData)
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_IPReportID), msgByte)
	//mlog.System("-----==============-----上報IP : %s", localIP)
	cb.ConnClt.SendMessage(msg)
	return nil
}

// SetToken ?
func (cb *ClientBase) SetToken(token string) {
	cb.Token = [REDACTED]
}

func (cb *ClientBase) tokenReport() error {
	msgData := new(netproto.ConnTokenMessage)
	msgData.Token = [REDACTED]
	msgByte, _ := proto.Marshal(msgData)
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_ConnTokenID), msgByte)
	//mlog.System("-----==============-----上報Token [%v]", cb.Token)
	cb.ConnClt.SendMessage(msg)
	return nil
}

// RequestHeartbeatConfig ?
func (cb *ClientBase) RequestHeartbeatConfig() error {
	// heartbeat
	msg := msg.NewMessage(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_RequestHeartBeatConfigID), nil)
	cb.SendMessage(msg)
	return nil
}

func (cb *ClientBase) tcpEventChanFunc(tcpevt *msgserver.MsgEvtChan) {
	defer util.PanicErrStack()
	limiter := time.Tick(time.Duration(1) * time.Millisecond)
	for {
		<-limiter
		select {
		case clt := <-tcpevt.ConnectedSig:
			mlog.System("connect tcpserver,%s", clt.GetConn().RemoteAddr().String())
			if err := cb.tokenReport(); err != nil {
				mlog.Error("發送token驗證發生錯誤: %s", err.Error())
			}
			if err := cb.ipReport(); err != nil {
				mlog.Error("上報IP發生錯誤: %s", err.Error())
			}
			cb.IPReported = true
			cb.RequestHeartbeatConfig()
		case clt := <-tcpevt.DisconnectedSig:
			mlog.Error("disconnect tcpserver,%v", clt.GetConn().RemoteAddr().String())
			cb.IPReported = false
		case <-tcpevt.DataSig:

		case msgdata := <-tcpevt.MsgSig:
			ctx := clientframework.NewCtxBase(cb.ConnClt, cb.Running, cb.EnableHeartBeat)
			cb.MsgHandler.OnMessage(ctx, msgdata)
		case <-time.After(time.Duration(30) * time.Second):
			mlog.Warn("消息处理线程退出")
			return
		}
	}
}

// SendMessage ?
func (cb *ClientBase) SendMessage(msg *msg.Message) error {
	if !cb.WaitConnect() {
		mlog.Error("發送消息等待連線超時, %v", msg)
		return fmt.Errorf("發送消息等待連線超時, %v", msg)
	}
	example.invalid("發送消息, %v", msg)
	return cb.ConnClt.SendMessage(msg)
}

// WaitConnect ?
func (cb *ClientBase) WaitConnect() bool {
	get := make(chan bool, 1)

	go func() {
		for {
			if cb.ConnClt.IsConnected() && cb.IPReported {
				get <- true
				return
			}
			if !cb.Running {
				get <- false
				return
			}
			time.Sleep(time.Duration(1) * time.Millisecond)
		}
	}()

	select {
	case r := <-get:
		return r
	case <-time.After(time.Duration(cb.Timeout) * time.Second):
		return false
		// the read from ch has timed out
	}
}

func selfIP() (localip string, err error) {
	conf := monitorutil.GetMonitorServerConfig("monitor")

	getip := conf
	localip, err = monitorutil.Get(getip.GetIP)
	if err != nil {
		return "", err
	}

	return localip, nil
}
