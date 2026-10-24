package faf

import (
	"example.invalid/ServerTeam/Servers/pkg/slot"
)

// Spin 玩家轉輪軸
func (logic *Logic) Spin(symbolGroupID int32) {
	// 產生遊戲結果
	manager, _ := slot.NewGameResultManager(logic.DefaultReelNum, logic.DefaultRowNum)
	manager.GenerateGameResult(symbolGroupID, logic.Proto.SlotSymbolDistribution)

	logic.GameResultManager = manager

	logic.LogGameResult(logic.GameResultManager)

	logic.CalcGameResult()

	logic.CalcGameAward()
}

// CalcGameResult 計算遊戲結果
func (logic *Logic) CalcGameResult() {
	logic.GamePayResult = logic.SlotLogicBase.CalcGameResult(logic.GameResultManager)
}

// CalcGameAward 計算遊戲獎勵
func (logic *Logic) CalcGameAward() {

	var ChipWinMoney int64 = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
	var HistoryType int32

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE

		if ChipWinMoney > 0 {
			logic.FreeSpin.WinCount++
		}
		if logic.FreeSpin.Origin.GetMultiplierType() == "level" {
			multi := logic.CalcMultiLevel(logic.FreeSpin.Origin.GetMultiplierValue(), logic.FreeSpin.WinCount)
			ChipWinMoney = ChipWinMoney * multi

			logic.GameLog.AddGameStep("免費遊戲 第 %d 獲勝, 贏得獎金乘 %d 倍, 總贏得: %d", logic.FreeSpin.WinCount, multi, ChipWinMoney)
		}
		logic.FreeSpin.WinMoney += ChipWinMoney
	} else {
		HistoryType = slot.HistoryTypeMAIN
	}

	logic.GameWinMoney = ChipWinMoney

	logic.GameResultManager.ReelMode()
	logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	logic.GameResultManager.RowMode()

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

			continue

		case slot.PayTypeAddTimes:

			if logic.FreeSpin == nil {
				continue
			}

			logic.WinFreeTimes(result.Rule, logic.FreeSpin)

			logic.UpdateFreeSpin(logic.GameUser)

			continue
		}

	}
}
