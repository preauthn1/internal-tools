package msghandler

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/hallserver/halllogic/hallutil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// GetHallVerHandler ?
type GetHallVerHandler struct {
	verprovider hallutil.HallVerProvider
}

// NewGetHallVerHandler ?
func NewGetHallVerHandler(verprovider hallutil.HallVerProvider) *GetHallVerHandler {
	hallverHandler := new(GetHallVerHandler)
	hallverHandler.verprovider = verprovider

	return hallverHandler
}

// OnMessage ?
func (h *GetHallVerHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	areaid := int64(0)

	dd := m.MsgData.([]byte)
	mm := &netproto.HallVerReq{}

	err := proto.Unmarshal(dd, mm)
	if err == nil {
		areaid = mm.GetAreaId()
	}

	vermsg := new(msg.Message)
	vermsg.BClassID = int32(netproto.MessageBClassID_Hall)
	vermsg.SClassID = int32(netproto.HallMsgClassID_HallVerInfoID)

	ver := h.verprovider.GetVer(ctx, clt, areaid)
	vermsg.MsgData = ver

	if ver == nil {
		return nil
	}

	return []*msg.Message{vermsg}
}

// IsHook ?
func (h *GetHallVerHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_Hall) && sclsid == int32(netproto.HallMsgClassID_RequestHallVerID)
}
