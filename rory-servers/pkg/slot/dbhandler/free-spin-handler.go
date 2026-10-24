package dbhandler

import (
	"example.invalid/golang/protobuf/proto"

	dbutil "example.invalid/ServerTeam/Servers/pkg/db/dbutil"
	"example.invalid/ServerTeam/Servers/pkg/mlog"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/network"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
	"example.invalid/ServerTeam/Servers/pkg/slot/dbhandler/service"
)

// FreeSpinHandler 寫入及讀取免費遊戲
type FreeSpinHandler struct {
	Service service.FreeSpinService
}

// NewFreeSpinHandler create FreeSpinHandler
func NewFreeSpinHandler(m Mode) *FreeSpinHandler {
	h := new(FreeSpinHandler)

	switch m {
	case DB:
		h.Service = new(service.DBFreeSpinService)
	case API:
		h.Service = new(service.APIFreeSpinService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *FreeSpinHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_FreeGameRecordID)
}

// OnMessage 接收訊息
func (h *FreeSpinHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_FreeGameReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_FreeGameRecordID:
		return h.Record(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 查詢免費遊戲
func (h *FreeSpinHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_UserGameQuery{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_UserGameQuery)

	info, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID())
	if err != nil {
		mlog.Error("Read FreeSpin Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_FreeGameRecord)
	example.invalid = info

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_FreeGameReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 更新免費遊戲
func (h *FreeSpinHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_FreeGameRecord{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_FreeGameRecord)

	if example.invalid == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	err := h.Service.Record(rpc.GetCer(), example.invalid.GetUserID(), example.invalid.GetGameID(), example.invalid.GetFlagID(), example.invalid.GetPatternID(), example.invalid.GetOriginID(), example.invalid.GetUserBet(), example.invalid.GetChooseRate(), example.invalid.GetFreeSpinGroup(), example.invalid.GetTimes(), example.invalid.GetWinMoney(), example.invalid.GetWinTimes(), example.invalid.GetData())
	if err != nil {
		mlog.Error("Record FreeSpin Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
