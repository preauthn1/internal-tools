package xmbx

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	// 檢查需要重轉的盤面的規則
	logic.HandleAfterSpin(manager)

	logic.GameResultManager = manager

	// 紀錄開獎符號
	logic.LogGameResult(manager)

	logic.CalcGameResult()

	logic.CalcGameAward()
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult() {
	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)
}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward() {

	logic.GameWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)

	if logic.FreeSpin != nil && logic.GameWinMoney > 0 {

		multi := logic.FreeSpin.GetMultipleValue()
		win := logic.GameWinMoney
		total := win * int64(multi)

		logic.FreeSpin.WinCount++
		logic.FreeSpin.WinMoney += total
		logic.GameWinMoney = total

		// logic.GameLog.AddGameStep("免費遊戲中第 %d 次中獎, 加乘 %d 倍, 原中獎金額: %d, 加乘後金額: %d", logic.FreeSpin.WinCount, multi, win, total)
	}

	for _, result := range logic.GamePayResult {

		if !result.IsMatch || result.Rule == nil {
			continue
		}

		switch result.Rule.GetUseType() {
		case "normal":
			if logic.IsFreeSpin() {
				continue
			}
		case "freespin":
			if !logic.IsFreeSpin() {
				continue
			}
		}

		switch result.GetPayType() {
		case slot.PayTypeFreeSpin:

			logic.WinFreeSpin(result.Rule, true)
			logic.UpdateFreeSpin(logic.GameUser)

		case slot.PayTypeJackpot:

			_, win := logic.WinJackpot(result.Rule)
			result.WinMoney = win

			if logic.IsFreeSpin() {
				logic.FreeSpin.WinMoney += win
			}

			logic.GameWinMoney += win
			logic.GameWinJackpot = win

		case slot.PayTypeLightning:

			logic.LightningGame = logic.WinLightningGame(result.Rule)
			logic.HandleGetAwardLightning(logic.LightningGame)

			logic.RecordLightningData(logic.GameUser)

		}

	}

}
