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

// UserAssignedModeHandler UserAssignedMode 相關處理
type UserAssignedModeHandler struct {
	Service service.UserAssignedModeService
}

// NewUserAssignedModeHandler create UserAssignedModeHandler
func NewUserAssignedModeHandler(m Mode) *UserAssignedModeHandler {
	h := new(UserAssignedModeHandler)

	switch m {
	case DB:
		h.Service = new(service.DBUserAssignedModeService)
	case API:
		h.Service = new(service.APIUserAssignedModeService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *UserAssignedModeHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeRecordID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeDeleteID)
}

// OnMessage 接收訊息
func (h *UserAssignedModeHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_UserAssignedModeRecordID:
		return h.Record(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_UserAssignedModeDeleteID:
		return h.Delete(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取UserAssignedMode資訊
func (h *UserAssignedModeHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_UserGameQuery{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_UserGameQuery)

	record, err := h.Service.Read(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID())
	if err != nil {
		mlog.Error("Read UserAssignedMode Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_UserAssignedModeRecord)
	data.UserID = proto.Int32(mm.GetUserID())
	data.GameID = proto.Int32(mm.GetGameID())
	data.FlagID = proto.Int32(mm.GetFlagID())
	data.Records = record

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_UserAssignedModeReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 紀錄 UserAssignedMode 資訊
func (h *UserAssignedModeHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_UserAssignedModeRecord{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Record UserAssignedMode Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_UserAssignedModeRecord)

	if mm.Records == nil {
		mlog.Error("Record UserAssignedMode Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	if err := h.Service.Record(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID(), mm.GetRecords().GetMode(), mm.GetRecords().GetExpireTime()); err != nil {
		mlog.Error("Record UserAssignedMode Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

// Delete 刪除 UserAssignedMode 資訊
func (h *UserAssignedModeHandler) Delete(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_UserGameQuery{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Delete UserAssignedMode Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_UserGameQuery)

	if err := h.Service.Delete(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetFlagID()); err != nil {
		mlog.Error("Delete UserAssignedMode Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}
