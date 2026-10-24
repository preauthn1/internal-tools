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

// JackpotNoLinkHandler Jackpot相關處理
type JackpotNoLinkHandler struct {
	Service service.JackpotNoLinkService
}

// NewJackpotNoLinkHandler create JackpotNoLinkHandler
func NewJackpotNoLinkHandler(m Mode) *JackpotNoLinkHandler {
	h := new(JackpotNoLinkHandler)

	switch m {
	case DB:
		h.Service = new(service.DBJackpotNoLinkService)
	case API:
		h.Service = new(service.APIJackpotNoLinkService)
	}

	return h
}

// IsHook 是否Mapping到此Handler
func (h *JackpotNoLinkHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID) || sclsid == int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID)
}

// OnMessage 接收訊息
func (h *JackpotNoLinkHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkRecordID:
		return h.Record(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取Jackpot資訊
func (h *JackpotNoLinkHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_JackpotRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_JackpotRead)

	if example.invalid == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	info := example.invalid

	records, err := h.Service.GetJackpotNoLink(rpc.GetCer(), info.GetGameID(), info.GetFlagID(), info.GetTableNo(), info.GetLevel())
	if err != nil {
		mlog.Error("Read Jackpot Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_JackpotRecord)
	data.Records = records

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Record 紀錄Jackpot資訊
func (h *JackpotNoLinkHandler) Record(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_JackpotRecord{}
		return mm
	})

	if rmm == nil {
		mlog.Error("Record Jackpot Error. Parse Message Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_JackpotRecord)

	if mm.Records == nil || len(mm.Records) == 0 {
		mlog.Error("Record Jackpot Error.")
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	var records []*netproto.Slot_JackpotData
	var err error

	for _, info := range mm.Records {
		records, err = h.Service.UpdateJackpotNoLink(rpc.GetCer(), info.GetGameID(), info.GetFlagID(), info.GetTableNo(), info.GetLevel(), info.GetContribution(), info.GetTotalContribution(), info.GetAwardTotal())
		if err != nil {
			mlog.Error("Record Jackpot Error. %v", err)
			return []*msg.Message{dbutil.GetDBRetMsg(0)}
		}
	}

	data := new(netproto.Slot_MSG_JackpotRecord)
	data.Records = records

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_JackpotNoLinkReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}
