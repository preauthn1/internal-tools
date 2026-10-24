package by

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// SettingHandler ?
type SettingHandler struct{}

// NewBySettingHandler ?
func NewBySettingHandler() *SettingHandler {
	return new(SettingHandler)
}

// IsHook ?
func (by *SettingHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.BY_GameMessageClassID_BYSettingID)
}

// OnMessage ?
func (by *SettingHandler) OnMessage(client *clientframework.CtxBase, m *msg.Message) {
	setting := &netproto.BY_Setting{}
	byteMsgData := m.MsgData.([]byte)
	proto.Unmarshal(byteMsgData, setting)
}

// SetCallback ?
func (by *SettingHandler) SetCallback(cb func(m *msg.Message)) {}
