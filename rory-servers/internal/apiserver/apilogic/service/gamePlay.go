package service

import (
	"net/http"

	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"

	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GamePlay ?
type GamePlay struct {
	*Base
}

// NewGamePlay ?
func NewGamePlay(ctx *module.MouduleCtx) *GamePlay {
	s := &GamePlay{}
	s.Base = NewBase(ctx)
	return s
}

// PostPlaceBet 投注
func (s *GamePlay) PostPlaceBet(accessToken string, serverID int32, tradeID, betAmount int64, ipAddr, memo string) (ok bool, afterAmount int64, transactionNo string) {
	api := s.buildBackendAPI("/game/placeBet", nil)
	header := s.buildHeader(accessToken)
	var groupRequest = new(httputil.WaitGroupRequest)

	params := map[string]interface{}{
		"serverID":  serverID,
		"tradeID":   tradeID,
		"betAmount": util.CentToMoney(betAmount),
		"ipAddress": ipAddr,
		"memo":      memo,
	}

	res := new(model.MapDataBody)

	groupRequest.AddTask(httputil.POST, api, params, header, res)

	responses := groupRequest.Send()

	if len(responses) == 0 {
		return false, 0, ""
	}

	if groupRequest.IsAllDone() {
		if res.Code != 0 {
			return false, 0, ""
		}

		afterAmount := res.Data["afterAmount"].(float64) * 100

		return true, int64(afterAmount), res.Data["transactionNo"].(string)
	}

	return false, 0, ""
}

// GetWallet 查詢玩家錢包金額
func (s *GamePlay) GetWallet(accessToken string) (ok bool, amount int64, merchantWalletMode int32) {
	api := s.buildBackendAPI("/players/wallet", nil)
	header := s.buildHeader(accessToken)
	var groupRequest = new(httputil.WaitGroupRequest)

	res := new(model.MapDataBody)

	groupRequest.AddTask(http.MethodGet, api, nil, header, res)

	responses := groupRequest.Send()

	if len(responses) == 0 {
		return false, 0, 0
	}

	if groupRequest.IsAllDone() {
		if res.Code != 0 {
			return false, 0, 0
		}

		amount := res.Data["amount"].(float64) * 100

		return true, int64(amount), int32(res.Data["merchantWalletMode"].(float64))
	}

	return false, 0, 0
}

// PostRoomIn ?
func (s *GamePlay) PostRoomIn(rpcInfo *netproto.RPCInfo, params map[string]interface{}) (bool, *model.GameRoomEnterBody, *model.PlayerWalletBody) {

	header := s.buildHeader(rpcInfo.GetCer())

	obj1 := &model.GameRoomEnterBody{}
	obj2 := &model.PlayerWalletBody{}

	{ // 先取得錢包
		api2 := s.buildBackendAPI("/players/wallet", nil)

		groupRequest := new(httputil.WaitGroupRequest)
		groupRequest.AddTask(httputil.GET, api2, nil, header, obj2)

		groupRequest.Send()

		if !groupRequest.IsAllDone() {
			return false, nil, obj2
		}
	}

	{ // 登入
		api1 := s.buildBackendAPI("/game/room/in", nil)

		groupRequest := new(httputil.WaitGroupRequest)
		groupRequest.AddTask(httputil.POST, api1, params, header, obj1)

		groupRequest.Send()

		if groupRequest.IsAllDone() {
			return true, obj1, obj2
		}

		return false, obj1, obj2
	}
}

// PostRoomOut ?
func (s *GamePlay) PostRoomOut(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/room/out", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)

	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}
