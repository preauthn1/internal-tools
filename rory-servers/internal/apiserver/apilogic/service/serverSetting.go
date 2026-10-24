package service

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// ServerSetting ?
type ServerSetting struct {
	*Base
}

// NewServerSetting ?
func NewServerSetting(ctx *module.MouduleCtx) *ServerSetting {
	s := &ServerSetting{}
	s.Base = NewBase(ctx)
	return s
}

// GetRobots ?
func (s *ServerSetting) GetRobots(rpcInfo *netproto.RPCInfo, query map[string]interface{}) (bool, *model.RobotListBody) {
	api := s.buildBackendAPI("/system/server/robots", query)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.RobotListBody{}
	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}

// PostReRegister ?
func (s *ServerSetting) PostReRegister(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/system/server/re-register", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// PostStartProcess ?
func (s *ServerSetting) PostStartProcess(rpcInfo *netproto.RPCInfo, params map[string]interface{}) (bool, *model.ServerStartProcessBody) {
	api := s.buildBackendAPI("/system/server/start-process", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.ServerStartProcessBody{}
	groupRequest.AddTask(httputil.POST, api, params, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}
