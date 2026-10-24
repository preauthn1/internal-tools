package bb

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

	var ChipWinMoney int64
	var HistoryType int32

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE
		logic.MultiMode = slot.MultiModeEachWildMulti
		// 免費遊戲檢查是否有龍加倍
		multi, d := logic.CalcMultiMode(slot.MultiModeEachWildMulti, logic.FreeSpin.Origin.GetMultiplierValue())
		ChipWinMoney = logic.SlotLogicBase.CalcGameAwardByMultiType(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate, multi, slot.MultiTypeLinkWildTarget, d)
		logic.GameMultiDetail = d

		if ChipWinMoney > 0 {
			logic.FreeSpin.WinCount++
			logic.FreeSpin.WinMoney += ChipWinMoney
		}
	} else {
		HistoryType = slot.HistoryTypeMAIN
		logic.MultiMode = slot.MultiModeNone
		// 一般遊戲錢照算
		ChipWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)
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
