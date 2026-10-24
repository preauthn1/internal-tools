package service

import (
	"example.invalid/ServerTeam/Servers/internal/dbserver/dal/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
)

// GameInterface 遊戲相關接口
type GameInterface struct {
	*Base
}

// NewGameInterface 遊戲相關接口
func NewGameInterface(ctx *module.MouduleCtx) *GameInterface {
	s := &GameInterface{}
	s.Base = NewBase(ctx)
	return s
}

// PostGameLaunch 三方遊戲啟動及登入
func(s *GameInterface) PostGameLaunch(rpcInfo *netproto.RPCInfo, gameID int32) (bool, *model.GameLaunchResponseBody){
	responseBody := &model.GameLaunchResponseBody{}
	api := s.buildGameAPI("/game/launch", nil)
	params := map[string]interface{}{
		"UserID": rpcInfo.GetUserID(),
		"Certificate": rpcInfo.GetCer(),
		"GameID": gameID,
		"IPAddress": rpcInfo.GetIPAddress(),
	}
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, responseBody)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true, responseBody
	}
	return false, responseBody
}

// PostWalletAutoTransfer 三方遊戲額度轉回
func(s *GameInterface) PostWalletAutoTransfer(rpcInfo *netproto.RPCInfo) bool{
	api := s.buildGameAPI("/wallet/auto-transfer", nil)
	params := map[string]interface{}{
		"UserID": rpcInfo.GetUserID(),
		"Certificate": rpcInfo.GetCer(),
		"IPAddress": rpcInfo.GetIPAddress(),
	}
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}