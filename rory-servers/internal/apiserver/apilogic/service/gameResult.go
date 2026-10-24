package service

import (
	"example.invalid/ServerTeam/Servers/internal/apiserver/apilogic/model"
	"example.invalid/ServerTeam/Servers/pkg/httputil"
	"example.invalid/ServerTeam/Servers/pkg/module"
	"example.invalid/ServerTeam/Servers/pkg/netproto"
	"example.invalid/ServerTeam/Servers/pkg/util"
)

// GameResult ?
type GameResult struct {
	*Base
}

// NewGameResult ?
func NewGameResult(ctx *module.MouduleCtx) *GameResult {
	s := &GameResult{}
	s.Base = NewBase(ctx)
	return s
}

// PostPerGameLog 寫入每局遊戲紀錄
func (s *GameResult) PostPerGameLog(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/result/per-game-log", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// PostSettlement ?
func (s *GameResult) PostSettlement(rpcInfo *netproto.RPCInfo, params map[string]interface{}) bool {
	api := s.buildBackendAPI("/game/result/settlement", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)
	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}
	return false
}

// PostSettlementV2 ?
func (s *GameResult) PostSettlementV2(accessToken string, playerID, serverID int32, tradeID, updateAmt int64, ipAddr string, transactionNo string, betAmt, validAmt, payoutAmt, taxAmt int64) (ok bool, afterAmt int64) {
	api := s.buildBackendAPI("/v2/game/result/settlement", nil)
	header := s.buildHeader(accessToken)
	var groupRequest = new(httputil.WaitGroupRequest)

	params := map[string]interface{}{
		"playerID":     playerID,
		"serverID":     serverID,
		"tradeID":      tradeID,
		"updateAmount": util.CentToMoney(updateAmt),
		"ipAddress":    ipAddr,
		"betInfo": map[string]interface{}{
			"transactionNo": transactionNo,
			"betAmount":     util.CentToMoney(betAmt),
			"validAmount":   util.CentToMoney(validAmt),
			"payoutAmount":  util.CentToMoney(payoutAmt),
			"taxAmount":     util.CentToMoney(taxAmt),
		},
	}

	res := new(model.MapDataBody)

	groupRequest.AddTask(httputil.POST, api, params, header, res)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		if res.Code != 0 {
			return false, 0
		}

		afterAmount := res.Data["afterAmount"].(float64) * 100

		return true, int64(afterAmount)
	}

	return false, 0
}

// PostPerGameLogV2 寫入每局遊戲紀錄
func (s *GameResult) PostPerGameLogV2(rpcInfo *netproto.RPCInfo, deallog *netproto.DealLog) bool {
	api := s.buildBackendAPI("/v2/game/result/per-game-log", nil)
	header := s.buildHeader(rpcInfo.GetCer())
	var groupRequest = new(httputil.WaitGroupRequest)

	var details []map[string]interface{}

	for _, row := range deallog.GetUserGameDetails() {

		if row.GetTradeID() == 0 {
			continue
		}

		d := make(map[string]interface{})

		d["playerID"] = row.GetUserID()
		d["tradeID"] = row.GetTradeID()
		d["currentAmount"] = util.CentToMoney(row.GetCurrentScore())
		d["updateAmount"] = util.CentToMoney(row.GetUpdateScore())
		d["betAmount"] = util.CentToMoney(row.GetBetScore())
		d["payoutAmount"] = util.CentToMoney(row.GetWinScore())
		d["validAmount"] = util.CentToMoney(row.GetValidScore())
		d["taxAmount"] = util.CentToMoney(row.GetTaxAmount())
		d["transactionNo"] = row.GetTransactionNo()

		details = append(details, d)
	}

	params := map[string]interface{}{
		"perGameLog": map[string]interface{}{
			"gameID":             deallog.GetGameID(),
			"serverID":           deallog.GetServerID(),
			"gameData":           deallog.GetGameData(),
			"taxAmount":          deallog.GetTaxAmount(),
			"gameBeginTime":      deallog.GetGameBeginTime(),
			"gameEndTime":        deallog.GetGameEndTime(),
			"userGameDetailLogs": details,
		},
	}

	groupRequest.AddTask(httputil.POST, api, params, header, nil)
	groupRequest.Send()

	if groupRequest.IsAllDone() {
		return true
	}

	return false
}
