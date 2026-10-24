package by

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/internal/monitorserver/clientframework"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// UserRecoverRetHandler 場景接收
type UserRecoverRetHandler struct{}

// NewUserRecoverRetHandler ?
func NewUserRecoverRetHandler() *UserRecoverRetHandler {
	return new(UserRecoverRetHandler)
}

// IsHook ?
func (by *UserRecoverRetHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Game) && sclsid == int32(netproto.BY_GameMessageClassID_UserRecoverRetID)
}

// OnMessage ?
func (by *UserRecoverRetHandler) OnMessage(client *clientframework.CtxBase, m *msg.Message) {
	urRecoverRet := &netproto.BY_UserRecoverRet{}
	byteMsgData := m.MsgData.([]byte)
	proto.Unmarshal(byteMsgData, urRecoverRet)

	user := urRecoverRet.GetUserStatus()
	_ = user
}

// SetCallback ?
func (by *UserRecoverRetHandler) SetCallback(cb func(m *msg.Message)) {}
