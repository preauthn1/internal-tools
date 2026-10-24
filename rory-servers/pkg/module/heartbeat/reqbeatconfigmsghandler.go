package heartbeat

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// TimeDuration
const (
	BeatInterval     int32 = 3  // 心跳包時間間隔
	ConnLoseTimeSpan int32 = 10 // 心跳超時時間
)

// ReqBeatConfigMsgHandler 請求心跳包配置消息處理器
type ReqBeatConfigMsgHandler struct{}

// NewReqBeatConfigMsgHandler 創建請求心跳包配置處理器
func NewReqBeatConfigMsgHandler() *ReqBeatConfigMsgHandler {
	reqh := new(ReqBeatConfigMsgHandler)

	return reqh
}

// IsHook 是否處理
func (bh *ReqBeatConfigMsgHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Common) && sclsid == int32(netproto.PlatformCommonClassID_RequestHeartBeatConfigID)
}

// OnMessage 處理訊息
func (bh *ReqBeatConfigMsgHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	example.invalid("接收到請求心跳配置.")

	bcmsg := new(msg.Message)
	bcmsg.BClassID = int32(netproto.MessageBClassID_Common)
	bcmsg.SClassID = int32(netproto.PlatformCommonClassID_HeartBeatConfigID)

	bcmsgData := new(netproto.HeartBeatConfig)
	bcmsgData.BeatInterval = proto.Int32(BeatInterval)
	bcmsgData.ConnLoseTimeSpan = proto.Int32(ConnLoseTimeSpan)
	bcmsg.MsgData = bcmsgData

	return []*msg.Message{bcmsg}
}
