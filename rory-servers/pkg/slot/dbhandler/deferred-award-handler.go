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

// DeferredAwardHandler DeferredAward相關處理
type DeferredAwardHandler struct {
	Service service.DeferredAwardService
}

// NewDeferredAwardHandler create DeferredAwardHandler
func NewDeferredAwardHandler(m Mode) *DeferredAwardHandler {
	h := new(DeferredAwardHandler)

	switch m {
	case DB:
		h.Service = new(service.DBDeferredAwardService)
	case API:
		h.Service = new(service.APIDeferredAwardService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *DeferredAwardHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID)
}

// OnMessage 接收訊息
func (h *DeferredAwardHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_DeferredAwardRecordID:
		return h.Record(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_DeferredAwardPopID:
		return h.Pop(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_DeferredAwardResetID:
		return h.Reset(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取DeferredAward資訊
func (h *DeferredAwardHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_DeferredAwardRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_DeferredAwardRead)

	records, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetAwardType())
	if err != nil {
		mlog.Error("Read DeferredAward Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_DeferredAwardRecord)
	data.UserID = proto.Int32(mm.GetUserID())
	data.GameID = proto.Int32(mm.GetGameID())
	data.FlagID = proto.Int32(mm.GetFlagID())
	data.Records = records

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_DeferredAwardReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 紀錄DeferredAward資訊
func (h *DeferredAwardHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_DeferredAwardRecord{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Record DeferredAward Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_DeferredAwardRecord)

	if mm.Records == nil || len(mm.Records) == 0 {
		mlog.Error("Record DeferredAward Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	var err error

	for _, info := range mm.Records {
		err = h.Service.Record(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), info.GetAwardType(), info.GetPatternID(), info.GetOriginID(), info.GetData())
		if err != nil {
			mlog.Error("Record DeferredAward Error. %v", err)
			return []*msg.Message{dbutil.GetDBRetMsg(0)}
		}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

// Pop 重置DeferredAward資訊
func (h *DeferredAwardHandler) Pop(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_DeferredAwardRead{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Pop DeferredAward Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_DeferredAwardRead)

	if record, err := h.Service.Pop(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetAwardType()); err != nil {
		mlog.Error("Pop DeferredAward Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	} else if record == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

// Reset 重置DeferredAward資訊
func (h *DeferredAwardHandler) Reset(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_DeferredAwardRead{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Reset DeferredAward Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_DeferredAwardRead)

	var exsit bool = true

	for exsit {
		if record, err := h.Service.Pop(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetAwardType()); err != nil {
			mlog.Error("Record DeferredAward Error. %v", err)
			return []*msg.Message{dbutil.GetDBRetMsg(0)}
		} else if record == nil {
			exsit = false
		}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
