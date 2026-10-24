package service

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
)

// SuperToolPool ?
type SuperToolPool struct {
	*Base
}

// NewSuperToolPool ?
func NewSuperToolPool(ctx *module.MouduleCtx) *SuperToolPool {
	s := &SuperToolPool{}
	s.Base = NewBase(ctx)
	return s
}

// GetWinLoss 獲取系統水池資料
func (s *SuperToolPool) GetWinLoss(query map[string]interface{}) (bool, *model.SuperToolPoolBody) {
	api := s.buildBackendAPI("/super-tool/pool/win-loss", query)
	header := s.buildHeader("")
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.SuperToolPoolBody{}

	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostWinLoss 更新系統水池資料
func (s *SuperToolPool) PostWinLoss(params map[string]interface{}) (bool, *model.SuperToolPoolBody) {
	api := s.buildBackendAPI("/super-tool/pool/win-loss", nil)
	header := s.buildHeader("")
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.SuperToolPoolBody{}

	groupRequest.AddTask(httputil.POST, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}
