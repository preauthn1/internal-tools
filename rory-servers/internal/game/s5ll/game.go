package s5ll

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

	ChipWinMoney = logic.SlotLogicBase.CalcGameAward(logic.GamePayResult, logic.GameBetMoney, logic.GameBetRate)

	if logic.IsFreeSpin() {
		HistoryType = slot.HistoryTypeFREE

		logic.FreeSpin.WinCount++
		logic.FreeSpin.WinMoney += ChipWinMoney
	} else {
		HistoryType = slot.HistoryTypeMAIN
	}

	logic.GameWinMoney = ChipWinMoney

	logic.GameResultManager.ReelMode()
	logic.AddSlotGameHistoryLogFull(logic.GameUser, logic.GameUser.Score, logic.GamePayMoney, ChipWinMoney, HistoryType, slot.LogTypeREEL, logic.GameResultManager.ToSymbolString())
	logic.GameResultManager.RowMode()

	var winSelection bool

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
		case slot.PayTypeSelection:

			if winSelection {
				continue
			}

			selection := logic.WinSelection(result.Rule)

			if selection == nil {
				continue
			}

			logic.SelectChances = append(logic.SelectChances, selection)

			if logic.IsFreeSpin() {
				logic.GameLog.AddGameStep("因目前正於免費遊戲中, 將於免費遊戲結束後開始選擇, 剩餘 %d 次", len(logic.SelectChances))
				logic.SendSelectChanceMessage(logic.GameUser)
			}

			logic.RecordSelectGameTimes(logic.GameUser)

			winSelection = true
		}

	}
}
