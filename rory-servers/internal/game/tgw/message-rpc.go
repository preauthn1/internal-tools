package tgw

import (
	"example.invalid/golang/protobuf/proto"

	"example.invalid/ServerTeam/Servers/pkg/msg"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/rpcmsg"
)

// OnRecvRPCMessageFunc 接收DB Server回傳的Message訊息
func (logic *Logic) OnRecvRPCMessageFunc(m *msg.Message) {
	switch m.SClassID {
	// 從DB獲得的遊戲歷程紀錄
	case int32(netproto.Slot_GameMessageClassID_Slot_GameLogReadRetID):
		logic.recvSlotGameLog(m)
	case int32(netproto.Slot_GameMessageClassID_Slot_BonusGameReadRetID):
		logic.RecvPickGame(m)
	}

}

// recvSlotGameLog 接收db server回傳的SlotGameLog資訊，直接傳給Client
func (logic *Logic) recvSlotGameLog(m *msg.Message) {
	if logic.GameUser == nil {
		return
	}

	// 解析出來的detail
	_, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.Slot_MSG_GameLogRecord{}
		return mm
	})

	detail := rmm.(*netproto.Slot_MSG_GameLogRecord)

	if detail.GetUserID() != logic.GameUser.UserID || detail.GetGameID() != logic.TableManager.Room.GameInfo.GameID {
		return
	}

	var historys []*netproto.Slot_SendHistoryLog
	for _, row := range detail.GetLogs() {
		h := new(netproto.Slot_SendHistoryLog)

		h.RowNum = proto.Int32(row.GetRowNum())
		h.CreateTime = proto.String(row.GetCreateTime())
		h.OrderID = proto.Int64(row.GetOrderID())
		h.BetMoney = proto.Int64(row.GetBetMoney())
		h.Profit = proto.Int64(row.GetProfit())

		historys = append(historys, h)
	}
	logic.SendGameLogMessage(logic.GameUser, detail.GetUserID(), detail.GetGameID(), detail.GetPageCount(), historys)
}
