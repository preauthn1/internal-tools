package service

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
)

// SystemSetting ?
type SystemSetting struct {
	*Base
}

// NewSystemSetting ?
func NewSystemSetting(ctx *module.MouduleCtx) *SystemSetting {
	s := &SystemSetting{}
	s.Base = NewBase(ctx)
	return s
}

// GetDenyIP ?
func (s *SystemSetting) GetDenyIP() (bool, *model.LoadBlockIPListBody) {
	api := s.buildBackendAPI("/system/setting/deny-ip", nil)
	header := s.buildHeader("")
	var groupRequest = new(httputil.WaitGroupRequest)
	obj := &model.LoadBlockIPListBody{}

	groupRequest.AddTask(httputil.GET, api, nil, header, obj)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, obj
	}
	return false, obj
}
