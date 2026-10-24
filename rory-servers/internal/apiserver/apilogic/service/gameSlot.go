package service

import (
	"fmt"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"

	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameSlot ?
type GameSlot struct {
	*Base
}

// NewGameSlot ?
func NewGameSlot(ctx *module.MouduleCtx) *GameSlot {
	s := &GameSlot{}
	s.Base = NewBase(ctx)
	return s
}

// GetBonus ?
func (s *GameSlot) GetBonus(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotBonusRecordReadBody) {
	api := s.buildBackendAPI("/game/slot/bonus", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotBonusRecordReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostBonus ?
func (s *GameSlot) PostBonus(rpcInfo *netproto.RPCInfo, params []map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/bonus", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)

	for i := 0; i < len(params); i++ {
		groupRequest.AddTask(httputil.POST, api, params[i], header, nil)
	}

	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// DeleteBonus ?
func (s *GameSlot) DeleteBonus(rpcInfo *netproto.RPCInfo, query map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/bonus", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)

	groupRequest.AddTask(httputil.DELETE, api, nil, header, nil)

	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// GetFreeSpin ?
func (s *GameSlot) GetFreeSpin(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotFreeSpinRecordReadBody) {
	api := s.buildBackendAPI("/game/slot/free-spin", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotFreeSpinRecordReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostFreeSpin ?
func (s *GameSlot) PostFreeSpin(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/free-spin", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// GetLog ?
func (s *GameSlot) GetLog(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.GameSlotHistoryLogBody) {
	api := s.buildBackendAPI("/game/slot/log", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.GameSlotHistoryLogBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostLog ?
func (s *GameSlot) PostLog(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/log", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// GetProb ?
func (s *GameSlot) GetProb(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotProbRecordReadBody) {
	api := s.buildBackendAPI("/game/slot/prob", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotProbRecordReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostProb ?
func (s *GameSlot) PostProb(rpcInfo *netproto.RPCInfo, params map[string]interface{}) (bool, *model.PlayerSlotProbRecordUpdateBody) {
	api := s.buildBackendAPI("/game/slot/prob", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotProbRecordUpdateBody{}
	groupRequest.AddTask(httputil.POST, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// DeleteProb ?
func (s *GameSlot) DeleteProb(rpcInfo *netproto.RPCInfo, query, params map[string]interface{}) (bool, *model.PlayerSlotProbRecordUpdateBody) {
	api := s.buildBackendAPI("/game/slot/prob", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotProbRecordUpdateBody{}
	groupRequest.AddTask(httputil.DELETE, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostJackpot ?
func (s *GameSlot) PostJackpot(rpcInfo *netproto.RPCInfo, gameID int32, params map[string]interface{}) (bool, *model.SlotJackpotBody) {
	api := s.buildBackendAPI(fmt.Sprintf("/game/%d/jackpot", gameID), nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.SlotJackpotBody{}
	groupRequest.AddTask(httputil.POST, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// GetLightning Get Lightning
func (s *GameSlot) GetLightning(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotLightningRecordReadBody) {
	api := s.buildBackendAPI("/game/slot/lightning", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotLightningRecordReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostLightning Post Lightning
func (s *GameSlot) PostLightning(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/lightning", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// GetDeferredAward Get DeferredAward
func (s *GameSlot) GetDeferredAward(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotDeferredAwardReadBody) {
	api := s.buildBackendAPI("/game/slot/deferred-award", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotDeferredAwardReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostDeferredAward Post DeferredAward
func (s *GameSlot) PostDeferredAward(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/deferred-award", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// DeleteDeferredAward Delete DeferredAward
func (s *GameSlot) DeleteDeferredAward(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotDeferredAwardDeleteBody) {
	api := s.buildBackendAPI("/game/slot/deferred-award", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotDeferredAwardDeleteBody{}
	groupRequest.AddTask(httputil.DELETE, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// GetUserAssignedMode Get UserAssignedMode
func (s *GameSlot) GetUserAssignedMode(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.PlayerSlotUserAssignedModeReadBody) {
	api := s.buildBackendAPI("/game/slot/user-assigned-mode", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.PlayerSlotUserAssignedModeReadBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostUserAssignedMode Post UserAssignedMode
func (s *GameSlot) PostUserAssignedMode(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/user-assigned-mode", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)

	obj := new(model.BodyBase)

	groupRequest.AddTask(httputil.POST, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// DeleteUserAssignedMode Delete UserAssignedMode
func (s *GameSlot) DeleteUserAssignedMode(rpcInfo *netproto.RPCInfo, query map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/slot/user-assigned-mode", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.DELETE, api, nil, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}
