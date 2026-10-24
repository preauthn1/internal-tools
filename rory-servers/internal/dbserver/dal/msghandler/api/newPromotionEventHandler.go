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

// NewPromotionEventHandler 幸運輪盤相關接口
type NewPromotionEventHandler struct {
}

// IsHook 是否Mapping到此Handler
func (h *NewPromotionEventHandler) IsHook(bclsid int32, sclsid int32) bool {
	if bclsid != int32(netproto.MessageBClassID_Hall) {
		return false
	}

	return sclsid == int32(netproto.HallMsgClassID_NewPromotionEventsInfoID) || sclsid == int32(netproto.HallMsgClassID_NewPromotionEventsListID) || sclsid == int32(netproto.HallMsgClassID_NewPromotionEventsLuckyRouletteRequestID)
}

// OnMessage 接收訊息
func (h *NewPromotionEventHandler) OnMessage(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {

	switch netproto.HallMsgClassID(m.SClassID) {
	case netproto.HallMsgClassID_NewPromotionEventsInfoID:
		return h.GetEventInfo(ctx, clt, m)
	case netproto.HallMsgClassID_NewPromotionEventsListID:
		return h.GetEventList(ctx, clt, m)
	case netproto.HallMsgClassID_NewPromotionEventsLuckyRouletteRequestID:
		return h.PostLuckyRouletteRequest(ctx, clt, m)
	}

	return []*msg.Message{}
}

// GetEventInfo 獲取幸運輪盤活動資訊
func (h *NewPromotionEventHandler) GetEventInfo(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NewPromotionEventInfoRequest{}
		return mm
	})
	mm := rmm.(*netproto.NewPromotionEventInfoRequest)
	if ok, body := service.PromotionEventService.GetEventInfo(rpcInfo, mm.GetEventID()); ok {
		msgData := &netproto.NewPromotionEventInfo{}
		{
			basicInfo := body.Data.BasicInfo
			msgData.BasicInfo = &netproto.NewPromotionEvent{
				EventID:   proto.Int32(basicInfo.EventID),
				EventName: proto.String(basicInfo.EventName),
				EventType: proto.Int32(basicInfo.EventType),
				StartTime: proto.String(basicInfo.StartTime.String()),
				EndTime:   proto.String(basicInfo.EndTime.String()),
			}
		}
		{
			luckyRouletteInfo := body.Data.LuckyRouletteInfo
			msgData.LuckyRouletteInfo = &netproto.NewPromotionLuckyRouletteInfo{
				TicketCount: proto.Int32(luckyRouletteInfo.TicketCount),
				EventMsg:    proto.String(luckyRouletteInfo.EventMsg),
				TipsMsg:     proto.String(luckyRouletteInfo.TipsMsg),
				BonusItems:  []*netproto.NewPromotionBonusItem{},
			}
			for _, item := range luckyRouletteInfo.BonusItems {
				msgData.GetLuckyRouletteInfo().BonusItems = append(msgData.GetLuckyRouletteInfo().BonusItems, &netproto.NewPromotionBonusItem{
					BonusID:   proto.Int32(item.BonusID),
					BonusName: proto.String(item.BonusName),
					Amount:    proto.Int32(item.Amount),
					BonusType: proto.Int32(item.BonusType),
					ImageURL:  proto.String(item.ImageURL),
				})
			}
		}
		retMsg := new(msg.Message)
		retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
		retMsg.SClassID = int32(netproto.HallMsgClassID_NewPromotionEventsInfoRetID)
		retMsg.MsgData = msgData
		return []*msg.Message{retMsg}
	}
	return nil
}

// GetEventList 獲取幸運輪盤活動列表
func (h *NewPromotionEventHandler) GetEventList(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, _ := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		return nil
	})

	msgData := &netproto.NewPromotionEventList{}
	msgData.Events = []*netproto.NewPromotionEvent{}
	if ok, body := service.PromotionEventService.GetEventList(rpcInfo); ok {
		for _, event := range body.Data.Events {
			msgData.Events = append(msgData.Events, &netproto.NewPromotionEvent{
				EventID:   proto.Int32(event.EventID),
				EventName: proto.String(event.EventName),
				EventType: proto.Int32(event.EventType),
				StartTime: proto.String(event.StartTime.String()),
				EndTime:   proto.String(event.EndTime.String()),
			})
		}
	}
	retMsg := new(msg.Message)
	retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
	retMsg.SClassID = int32(netproto.HallMsgClassID_NewPromotionEventsListRetID)
	retMsg.MsgData = msgData
	return []*msg.Message{retMsg}
}

// PostLuckyRouletteRequest 抽獎
func (h *NewPromotionEventHandler) PostLuckyRouletteRequest(ctx *module.MouduleCtx, clt *network.ConnInfo, m *msg.Message) []*msg.Message {
	rpcInfo, rmm := rpcmsg.ParseRPCMessage(m, func() proto.Message {
		mm := &netproto.NewPromotionEventInfoRequest{}
		return mm
	})
	mm := rmm.(*netproto.NewPromotionEventInfoRequest)
	params := map[string]interface{}{
		"eventID":     mm.GetEventID(),
		"userID":      rpcInfo.GetUserID(),
		"certificate": rpcInfo.GetCer(),
	}
	if ok, body := service.PromotionEventService.PostLuckyRouletteRequest(rpcInfo, params); ok {
		data := body.Data
		msgData := &netproto.NewPromotionLuckyRouletteData{
			IsSuccess:        proto.Bool(data.IsSuccess),
			ResultMsg:        proto.String(data.ResultMsg),
			BonusID:          proto.Int32(data.BonusID),
			LotteryRemaining: proto.Int32(data.LotteryRemaining),
		}
		retMsg := new(msg.Message)
		retMsg.BClassID = int32(netproto.MessageBClassID_Hall)
		retMsg.SClassID = int32(netproto.HallMsgClassID_NewPromotionEventsLuckyRouletteRequestRetID)
		retMsg.MsgData = msgData
		return []*msg.Message{retMsg}
	}
	return nil
}
