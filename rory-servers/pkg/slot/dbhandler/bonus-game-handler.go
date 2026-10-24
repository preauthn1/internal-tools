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

// BonusGameHandler BonusGame相關處理
type BonusGameHandler struct {
	Service service.BonusGameService
}

// NewBonusGameHandler create BonusGameHandler
func NewBonusGameHandler(m Mode) *BonusGameHandler {
	h := new(BonusGameHandler)

	switch m {
	case DB:
		h.Service = new(service.DBBonusGameService)
	case API:
		h.Service = new(service.APIBonusGameService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *BonusGameHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_BonusGameResetID)
}

// OnMessage 接收訊息
func (h *BonusGameHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_BonusGameReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_BonusGameRecordID:
		return h.Record(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_BonusGameResetID:
		return h.Reset(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取BonusGame資訊
func (h *BonusGameHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_BonusGameRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_BonusGameRead)

	records, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID())
	if err != nil {
		mlog.Error("Read BonusGame Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_BonusGameRecord)
	data.Detail = records

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 紀錄BonusGame資訊
func (h *BonusGameHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_BonusGameRecord{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Record BonusGame Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_BonusGameRecord)

	if mm.Detail == nil || len(mm.Detail) == 0 {
		mlog.Warn("Record BonusGame Error. Empty Detail")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	var err error

	for _, info := range mm.Detail {
		err = h.Service.Record(rpc.GetCer(), info.GetUserID(), info.GetGameID(), info.GetFlagID(), info.GetPatternID(), info.GetOriginID(), info.GetTurnID(), info.GetLocation(), info.GetBetMoney(), info.GetBetRate(), info.GetData())
		if err != nil {
			mlog.Error("Record BonusGame Error. %v", err)
			return []*msg.Message{dbutil.GetDBRetMsg(0)}
		}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

// Reset 重置BonusGame資訊
func (h *BonusGameHandler) Reset(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_BonusGameRead{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Reset BonusGame Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_BonusGameRead)

	if err := h.Service.Reset(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetOriginID()); err != nil {
		mlog.Error("Record BonusGame Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
