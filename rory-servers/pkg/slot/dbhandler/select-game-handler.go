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

// SelectGameHandler 寫入及讀取選擇遊戲機會
type SelectGameHandler struct {
	Service service.SelectionService
}

// NewSelectGameHandler create SelectGameHandler
func NewSelectGameHandler(m Mode) *SelectGameHandler {
	h := new(SelectGameHandler)

	switch m {
	case DB:
		h.Service = new(service.DBSelectionService)
	case API:
		h.Service = new(service.APISelectionService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *SelectGameHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_SelectGameRecordID)
}

// OnMessage 接收訊息
func (h *SelectGameHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_SelectGameReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_SelectGameRecordID:
		return h.Record(ctx, clt, m)
	}

	return []*msg.Message{}
}

// ReadSelectGame 查詢選擇遊戲機會
func (h *SelectGameHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_SelectGameRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_SelectGameRead)

	info, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID())

	if err != nil {
		mlog.Error("Read SelectGame Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_SelectGameRecord)
	example.invalid = info

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_SelectGameReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 更新選擇遊戲機會
func (h *SelectGameHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_SelectGameRecord{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_SelectGameRecord)

	if example.invalid == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if mm.GetInfo().GetTimes() == 0 {
		err := h.Service.Reset(rpc.GetCer(), mm.GetInfo().GetUserID(), mm.GetInfo().GetGameID(), mm.GetInfo().GetFlagID())
		if err != nil {
			mlog.Error("Record SelectGame Error. %v", err)
			return []*msg.Message{dbutil.GetDBRetMsg(0)}
		}

		return []*msg.Message{dbutil.GetDBRetMsg(1)}
	}

	err := h.Service.Record(rpc.GetCer(), mm.GetInfo().GetUserID(), mm.GetInfo().GetGameID(), mm.GetInfo().GetFlagID(), mm.GetInfo().GetPatternID(), mm.GetInfo().GetOriginID(), mm.GetInfo().GetTimes(), mm.GetInfo().GetData())
	if err != nil {
		mlog.Error("Record SelectGame Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
