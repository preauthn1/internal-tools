package msghandler

import (
	"time"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// HeartBeatConfigHandler 找回銀行密碼
type HeartBeatConfigHandler struct {
	callback func(m *msg.Message)
	interval int32
	timeout  int32
	hbSig    chan bool // 心跳超時訊號
	stopSig  chan bool // 心跳中止訊號
}

// NewHeartBeatConfigHandler ?
func NewHeartBeatConfigHandler() *HeartBeatConfigHandler {
	reqh := new(HeartBeatConfigHandler)
	reqh.hbSig = make(chan bool, 1)
	reqh.stopSig = make(chan bool, 1)
	return reqh
}

// IsHook ?
func (rh *HeartBeatConfigHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Common) && (sclsid == int32(netproto.PlatformCommonClassID_HeartBeatConfigID) || sclsid == int32(netproto.PlatformCommonClassID_HeartBeatReturnID))
}

// OnMessage ?
func (rh *HeartBeatConfigHandler) OnMessage(client *clientframework.CtxBase, m *msg.Message) {
	msg := &netproto.HeartBeatConfig{}
	byteMsgData := m.MsgData.([]byte) // 須將interface轉[]byte
	proto.Unmarshal(byteMsgData, msg)
	if m.BClassID == int32(netproto.MessageBClassID_Common) && m.SClassID == int32(netproto.PlatformCommonClassID_HeartBeatConfigID) {
		rh.timeout = msg.GetConnLoseTimeSpan()
		rh.interval = msg.GetBeatInterval()
		go rh.heartbeatInterval(client)
		go rh.setTimeout(client)

	} else {
		rh.hbSig <- true
	}
}

// SetCallback ?
func (rh *HeartBeatConfigHandler) SetCallback(f func(m *msg.Message)) {}

func (rh *HeartBeatConfigHandler) heartbeatInterval(client *clientframework.CtxBase) {
ForEnd:
	for {
		time.Sleep(time.Duration(rh.interval) * time.Second)
		if !client.Running {
			break ForEnd
		}
		select {
		case <-rh.stopSig:
			break ForEnd
		default:
		}
		msg := msg.NewMessage(int32(netproto.MessageBClassID_Common), int32(netproto.PlatformCommonClassID_HeartBeatID), nil)
		client.ConnClt.SendMessage(msg)
	}
	example.invalid("Heartbeat 線程終止")
	return
}

func (rh *HeartBeatConfigHandler) setTimeout(client *clientframework.CtxBase) {
ForEnd:
	for {
		select {
		case <-rh.hbSig:
		case <-time.After(time.Duration(rh.timeout) * time.Second):
			if client.EnableHeartBeat {
				mlog.Error("心跳time after rh.timeout[%v] EnableHeartBeat[%v]", rh.timeout, client.EnableHeartBeat)
				break ForEnd
			}
		}
	}
	rh.stopSig <- true
	rh.HeartbeatTimeout()
	mlog.System("心跳超時")
	return
}

// HeartbeatTimeout ?
func (rh *HeartBeatConfigHandler) HeartbeatTimeout() {
	mlog.Error("心跳超時")
}
