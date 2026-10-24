package service

import (
	"fmt"
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// PromotionEvent 活動相關接口
type PromotionEvent struct {
	*Base
}

// NewPromotionEvent 活動相關接口
func NewPromotionEvent(ctx *module.MouduleCtx) *PromotionEvent {
	s := &PromotionEvent{}
	s.Base = NewBase(ctx)
	return s
}

// GetEventList 獲取活動列表
func(s *PromotionEvent) GetEventList(rpcInfo *netproto.RPCInfo) (bool, *model.GetEventListBody){
	responseBody := &model.GetEventListBody{}
	api := s.buildGameAPI("/events", map[string]interface{}{
		"userID": rpcInfo.GetUserID(),
		"certificate": rpcInfo.GetCer(),
	})
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.GET, api, nil, header, responseBody)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, responseBody
	}
	return false, nil
}

// GetEventInfo 獲取活動資訊
func(s *PromotionEvent) GetEventInfo(rpcInfo *netproto.RPCInfo, eventID int32) (bool, *model.GetEventInfoBody){
	responseBody := &model.GetEventInfoBody{}
	api := s.buildGameAPI(fmt.Sprintf("/events/%d", eventID), map[string]interface{}{
		"userID": rpcInfo.GetUserID(),
		"certificate": rpcInfo.GetCer(),
	})
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.GET, api, nil, header, responseBody)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, responseBody
	}
	return false, nil
}

// PostLuckyRouletteRequest 參加抽獎
func(s *PromotionEvent) PostLuckyRouletteRequest(rpcInfo *netproto.RPCInfo, params map[string]interface{}) (bool, *model.PostLuckyRouletteResponseBody){
	responseBody := &model.PostLuckyRouletteResponseBody{}
	api := s.buildGameAPI("/events/lucky-roulette/request", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, responseBody)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, responseBody
	}
	return false, nil
}