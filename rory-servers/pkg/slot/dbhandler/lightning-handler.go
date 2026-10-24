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

// LightningHandler Jackpot相關處理
type LightningHandler struct {
	Service service.LightningService
}

// NewLightningHandler create LightningHandler
func NewLightningHandler(m Mode) *LightningHandler {
	h := new(LightningHandler)

	switch m {
	case DB:
		h.Service = new(service.DBLightningService)
	case API:
		h.Service = new(service.APILightningService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *LightningHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_LightningReadID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_LightningRecordID)
}

// OnMessage 接收訊息
func (h *LightningHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_LightningReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_LightningRecordID:
		return h.Record(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取Lightning資訊
func (h *LightningHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_LightningRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_LightningRead)

	patternID, _data, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID())
	if err != nil {
		mlog.Error("Read Lightning Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_LightningRecord)
	data.UserID = proto.Int32(mm.GetUserID())
	data.GameID = proto.Int32(mm.GetGameID())
	data.FlagID = proto.Int32(mm.GetFlagID())
	data.PatternID = proto.Int32(patternID)
	data.Data = proto.String(_data)

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_LightningReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 紀錄Lightning資訊
func (h *LightningHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_LightningRecord{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Record Lightning Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_LightningRecord)

	if err := h.Service.Record(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetPatternID(), mm.GetData()); err != nil {
		mlog.Error("Record Lightning Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
