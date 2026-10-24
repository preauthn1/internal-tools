package api

import (
	"example.invalid/golang/protobuf/proto"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/service"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// ThirdPartyGameInterfaceHandler 三方遊戲接口
type ThirdPartyGameInterfaceHandler struct {
}

// IsHook 是否Mapping到此Handler
func (h *ThirdPartyGameInterfaceHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Hall) {
		return false
	}

	return sclsid == int32(netproto.HallMsgClassID_ThirdPartyGameLaunchID) || sclsid == int32(netproto.HallMsgClassID_ThirdPartyGameWalletAutoTransferID)
}

// OnMessage 接收訊息
func (h *ThirdPartyGameInterfaceHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.HallMsgClassID(m.SClassID) {
	case netproto.HallMsgClassID_ThirdPartyGameLaunchID:
		return h.PostThirdPartyGameLaunch(ctx, clt, m)
	case netproto.HallMsgClassID_ThirdPartyGameWalletAutoTransferID:
		return h.PostThirdPartyGameWalletAutoTransfer(ctx, clt, m)
	}

	return []*msg.Message{}
}

// PostThirdPartyGameLaunch 三方遊戲登入接口
func (h *ThirdPartyGameInterfaceHandler) PostThirdPartyGameLaunch(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.ThirdPartyGameLaunchRequest{}
		return mm
	})
	mm := rmm.(*netproto.ThirdPartyGameLaunchRequest)
	gameID := mm.GetGameID()
	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retMsg.SClassID = int32(netproto.HallMsgClassID_ThirdPartyGameLaunchRetID)
	data := new(netproto.ThirdPartyGameLaunchResponse)
	if ok, body := service.GameInterfaceService.PostGameLaunch(rpcInfo, gameID); ok {
		data.RetCode = proto.Int32(1)
		data.RetMessage = proto.String(body.Message)
		data.RedirectUrl = proto.String(body.Data.RedirectURL)
	} else {
		data.RetCode = proto.Int32(0)
		data.RetMessage = proto.String(body.Message)
		data.RedirectUrl = proto.String("")
	}
	retMsg.MsgData = data
	return []*msg.Message{retMsg}
}

// PostThirdPartyGameWalletAutoTransfer 三方遊戲額度轉回
func (h *ThirdPartyGameInterfaceHandler) PostThirdPartyGameWalletAutoTransfer(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, _ := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		return nil
	})
	service.GameInterfaceService.PostWalletAutoTransfer(rpcInfo)
	return nil
}
