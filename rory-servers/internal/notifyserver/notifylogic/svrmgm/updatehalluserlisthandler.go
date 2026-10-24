package svrmgm

import (
	"sync"

	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// UpdateHallUserListHandler 服務器註冊處理器
type UpdateHallUserListHandler struct {
	hallUserList *sync.Map
}

// NewUpdateHallUserListHandler 創建服務器註冊處理器
func NewUpdateHallUserListHandler(hallUserList *sync.Map) *UpdateHallUserListHandler {
	handler := new(UpdateHallUserListHandler)
	handler.hallUserList = hallUserList
	return handler
}

// IsHook ?
func (h *UpdateHallUserListHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_UpdateHallUserListID)
}

// OnMessage ?
func (h *UpdateHallUserListHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	mmbyte := m.MsgData.([]byte)

	mm := &netproto.HallUserList{}
	err := proto.Unmarshal(mmbyte, mm)
	if err != nil {
		ctx.Error("Analyze the notifyserver received the lobby update online user error,%v", err)
		ctx.Server.CloseConn(clt)
		return nil
	}
	h.hallUserList.Store(mm.GetServerID(), mm)
	ctx.Debug("收到大廳更新UpdateHallUserList, len: %d", len(mm.UserIDList))
	return nil
}
