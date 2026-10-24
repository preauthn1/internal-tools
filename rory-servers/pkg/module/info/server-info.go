package info

import (
	"time"

	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/pkg/conf"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/msgserver"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// ServerInfoHandler 處理查詢Server資訊
type ServerInfoHandler struct {
	msgServer *msgserver.MsgServer
}

// NewServerInfoHandler Create ServerInfoHandler
func NewServerInfoHandler() *ServerInfoHandler {
	return &ServerInfoHandler{}
}

// IsHook 是否需要處理
func (bmh *ServerInfoHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Common) && sclsid == int32(netproto.PlatformCommonClassID_ServerInfoID)
}

// OnMessage 收到訊息
func (bmh *ServerInfoHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	mm := new(msg.Message)
	mm.BClassID = int32(netproto.MessageBClassID_Common)
	mm.SClassID = int32(netproto.PlatformCommonClassID_ServerInfoID)

	now := time.Now()
	zone, offset := now.Zone()

	mm.MsgData = &netproto.ServerInfo{
		Version:    proto.String(conf.GitTag),
		TimeZone:   proto.String(zone),
		TimeOffset: proto.Int64(int64(offset)),
	}

	return []*msg.Message{mm}
}
