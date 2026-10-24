package dbhandler

import (
	"fmt"

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

// GameLogHandler 遊戲歷程相關處理
type GameLogHandler struct {
	Service  service.GameLogService
	BatchLog []*netproto.Slot_MSG_GameLogAdd
}

// NewGameLogHandler create GameLogHandler
func NewGameLogHandler(m Mode) *GameLogHandler {
	h := new(GameLogHandler)

	switch m {
	case DB:
		h.Service = new(service.DBGameLogService)
	case API:
		h.Service = new(service.APIGameLogService)
	}

	h.BatchLog = make([]*netproto.Slot_MSG_GameLogAdd, 0)

	return h
}

// IsHook 是否Mapping到此Handler
func (h *GameLogHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Game) {
		return false
	}

	return sclsid == int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_GameLogAddID) ||
		sclsid == int32(netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID)
}

// OnMessage 接收訊息
func (h *GameLogHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.Slot_GameMessageClassID(m.SClassID) {
	case netproto.Slot_GameMessageClassID_Slot_GameLogReadID:
		return h.Read(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_GameLogAddID:
		return h.Add(ctx, clt, m)
	case netproto.Slot_GameMessageClassID_Slot_GameLogBatchAddID:
		return h.BatchAdd(ctx, clt, m)
	}

	return []*msg.Message{}
}

// Read 讀取遊戲歷程資訊
func (h *GameLogHandler) Read(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_GameLogRead{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_GameLogRead)

	rowcnt, records, err := h.Service.Search(rpc.GetCer(), mm.GetUserID(), mm.GetGameID(), mm.GetPageIndex(), mm.GetPageSize(), mm.GetStartTime(), mm.GetEndTime())
	if err != nil {
		mlog.Error("Read GameLog Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	data := new(netproto.Slot_MSG_GameLogRecord)
	data.UserID = proto.Int32(mm.GetUserID())
	data.GameID = proto.Int32(mm.GetGameID())
	data.PageCount = proto.Int32(rowcnt)
	data.Logs = records

	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Game)
	retMsg.SClassID = int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID)
	retMsg.MsgData = data

	return []*msg.Message{retMsg}
}

// Add 新增遊戲歷程資訊
func (h *GameLogHandler) Add(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_GameLogAdd{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_GameLogAdd)

	err := h.Service.Add(rpc.GetCer(), mm)
	if err != nil {
		mlog.Error("Add GameLog Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

// BatchAdd 新增遊戲歷程資訊(多筆)
func (h *GameLogHandler) BatchAdd(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	rpc, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_BatchGameLogAdd{}
		return mm
	})

	if rmm == nil {
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	mm := rmm.(*netproto.Slot_MSG_BatchGameLogAdd)

	err := h.Service.BatchAdd(rpc.GetCer(), mm.GetLogs())

	if err != nil {
		mlog.Error("BatchAdd GameLog Error. %v", err)
		return []*msg.Message{dbutil.GetDBRetMsg(0)}
	}

	return []*msg.Message{dbutil.GetDBRetMsg(1)}
}

func (h *GameLogHandler) batchLogToString(logs []*netproto.Slot_MSG_GameLogAdd) string {
	if len(logs) > 0 {
		var datalist string
		for _, v := range logs {
			datalist += fmt.Sprintf("%v;%v;%v;%v;%v;%v;%v;%v;%v;%v;%v!",
				v.GetCreateTime(), v.GetGameID(), v.GetServerID(), v.GetUserID(), v.GetCurrentAmount(), v.GetBetMoney(), v.GetWinMoney(), v.GetFinalScore(), v.GetGameType(), v.GetLogType(), v.GetGameData())
		}
		datalist = datalist[:len(datalist)-1]
		// fmt.Println(datalist)
		return datalist
	}
	return ""
}
