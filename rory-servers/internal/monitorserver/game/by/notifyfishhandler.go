package by

import (
	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// NotifyFishHandler ?
type NotifyFishHandler struct{}

// NewNotifyFishHandler ?
func NewNotifyFishHandler() *NotifyFishHandler {
	return new(NotifyFishHandler)
}

// IsHook ?
func (by *NotifyFishHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.BY_GameMessageClassID_BYNotifyFish)
}

// OnMessage ?
func (by *NotifyFishHandler) OnMessage(client *clientframework.CtxBase, m *msg.Message) {}

// SetCallback ?
func (by *NotifyFishHandler) SetCallback(cb func(m *msg.Message)) {}
