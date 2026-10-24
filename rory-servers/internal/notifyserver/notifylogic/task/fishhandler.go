package task

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
)

// FishHandler 捕魚公告
type FishHandler struct {
	ctx *RunCtx
}

// NewFishHandler ?
func NewFishHandler(ctx *RunCtx) *FishHandler {
	fh := new(FishHandler)
	fh.ctx = ctx
	return fh
}

// IsHook ?
func (svrreg *FishHandler) IsHook(bclsid int32, sclsid int32) bool {
	return bclsid == int32(netproto.MessageBClassID_NotifyServer) && sclsid == int32(netproto.NotifyServerClassID_AnnFishRegID)
}

// OnMessage ?
func (svrreg *FishHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	mmbyte := m.MsgData.([]byte)
	mm := &netproto.BY_NotifyFishMsg{}
	err := proto.Unmarshal(mmbyte, mm)

	if err != nil {
		ctx.Error("解析捕魚notifyserver註冊消息錯誤,%v", err)
		ctx.Server.CloseConn(clt)
		return nil
	}

	annmsg := new(msg.Message)
	annmsg.BClassID = int32(netproto.MessageBClassID_NotifyServer)
	annmsg.SClassID = int32(netproto.NotifyServerClassID_AnnFishID)
	annmsg.MsgData = mm
	svrreg.ctx.SendHallMessage(annmsg)

	return nil
}
